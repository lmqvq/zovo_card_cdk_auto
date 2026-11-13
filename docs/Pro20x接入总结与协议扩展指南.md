# 卡密系统架构与「支付协议接入」指南

> 版本基线：v1.4.35.24
> 适用对象：后续接 Grok / Claude / Devin / 其他卡台时的复用模板

---

## 一、系统全景

```
                          ┌─────────────────────────────────────────────┐
                          │              用户浏览器 (站点域名)              │
                          │  兑换页 RechargeView  4 步：                 │
                          │  ① 输入卡密 → ② 粘贴 Session → ③ 确认 → ④ 进度 │
                          └───────────────────┬─────────────────────────┘
                                              │ HTTPS（Cloudflare 前置）
                                              ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                          Go 后端 (Gin + SQLite)                              │
│                                                                              │
│  /api/v1/public/cdk/preview    ──► 判断码来源 ──┬──► 本地码 → 本地预览       │
│  /api/v1/public/cdk/preflight                     └──► 卡台码 → 转发 zovo   │
│  /api/v1/public/cdk/redeem                                                   │
│  /api/v1/public/cdk/result                                                   │
│  /api/v1/public/cdk/result-by-code                                           │
│                                                                              │
│  /api/v1/admin/cdk/local/*   ──► 自建码管理（生成/列表/状态/备注）           │
│  /api/v1/admin/pro20x/*      ──► Pro20x 通道配置与连通性自检                 │
└───────────────┬──────────────────────────────────────┬───────────────────────┘
                │                                      │
    ┌───────────▼───────────┐              ┌───────────▼────────────────────┐
    │  通道 A：zovo 卡台     │              │  通道 B：Pro20x Automate        │
    │  （卡台直充，历史链路） │              │  （自建码专用，本期新增）        │
    │                       │              │                                │
    │  1. 卡台开卡并注资      │              │  1. POST /account-check        │
    │  2. 卡台提交 Stripe    │              │     → 轮询 poll（校验账号）     │
    │  3. 卡台绑定到账号      │              │  2. POST /tasks（提交充值）     │
    │  4. 本站监控删卡        │              │  3. GET 任务（轮询进度）        │
    └───────────┬───────────┘              └───────────┬────────────────────┘
                │                                      │
                ▼                                      ▼
        ┌───────────────┐                     ┌────────────────────┐
        │ ZovoCard 卡台  │                     │ Pro20x → ZovoCard  │
        │ 余额/卡池      │                     │ 用卡策略选卡注资    │
        └───────────────┘                     └────────────────────┘
                                                       │
                                                       ▼
                                            ┌─────────────────────┐
                                            │ ChatGPT billing     │
                                            │ Stripe 支付 + 删卡  │
                                            └─────────────────────┘
```

---

## 二、兑换链路（详细时序）

### 2.1 分流判据

只有一个判据：**卡密在本地表里的 `source` 字段**。

```go
// handler/cdk_local_redeem.go
func isLocalCDKCode(code string) bool {
    _, ok, _ := db.LookupLocalCDK(code)   // SELECT ... WHERE code=? AND source='local'
    return ok
}
```

- `source = 'local'` → 走 Pro20x
- `source = 'zovo'` 或查不到 → 原样转发卡台（老链路一个字节没改）

### 2.2 通道 A：zovo 卡台（历史链路，未改动）

```
用户                本站后端                  zovo 卡台
 │  POST preview      │                          │
 │───────────────────►│ GET/POST /openapi/v1/... │
 │                    │─────────────────────────►│
 │   redemption_token │◄─────────────────────────│
 │◄───────────────────│                          │
 │  POST preflight    │                          │
 │───────────────────►│ 转发 + device 头         │
 │   preflight_token  │◄─────────────────────────│
 │◄───────────────────│                          │
 │  POST redeem       │                          │
 │───────────────────►│ 注入选卡策略/坏卡黑名单   │
 │                    │─────────────────────────►│ 开卡→注资→Stripe
 │   {id,status}      │◄─────────────────────────│
 │◄───────────────────│                          │
 │  GET result (轮询)  │                          │
 │───────────────────►│─────────────────────────►│
 │   {order,events}   │◄─────────────────────────│
 │◄───────────────────│                          │
```

### 2.3 通道 B：Pro20x（本期新增）

```
用户            本站后端                       Pro20x                zovo 卡台
 │ preview         │                              │                      │
 │────────────────►│ 查本地码：状态/有效期/套餐映射 │                      │
 │  RT-token       │ 插入 pro20x_orders(preview)   │                      │
 │◄────────────────│                              │                      │
 │ preflight       │                              │                      │
 │────────────────►│ sessionAsJSON(session)        │                      │
 │                 │─────────────────────────────►│ POST /account-check  │
 │                 │                              │ 202 {job,poll}       │
 │                 │◄─────────────────────────────│                      │
 │                 │ GET poll（3 秒一次）          │                      │
 │                 │─────────────────────────────►│                      │
 │                 │◄─────────────────────────────│ {email,tier,active}  │
 │                 │ 已订阅 → 直接拒绝             │                      │
 │                 │ 存 session + status=preflight │                      │
 │  {preflight_token, email, currentPlan, ...}    │                      │
 │◄────────────────│                              │                      │
 │ redeem          │                              │                      │
 │────────────────►│ ClaimPro20xOrder(preflight→redeeming)                │
 │                 │ ClaimLocalCDK(unused→consuming)                      │
 │                 │ building: planId+country+authSession                 │
 │                 │ +cancelAutoRenew+deletePaymentCard                    │
 │                 │─────────────────────────────►│ POST /tasks          │
 │                 │                              │─────────────────────►│ 选卡/开卡/注资
 │                 │◄─────────────────────────────│ {task.id,queued}     │
 │                 │ ConsumeLocalCDK(→consumed)   │                      │
 │                 │ status=running + task_id     │                      │
 │  {id,status}    │                              │                      │
 │◄────────────────│                              │─────────────────────►│ Stripe 支付
 │ result (轮询)    │                              │                      │ 删卡核验
 │────────────────►│ 3 秒节流 → GetTask            │                      │
 │                 │─────────────────────────────►│                      │
 │  {order,events} │◄─────────────────────────────│ {status,steps,cleanup}│
 │◄────────────────│ mapLocalStatus → 前端认的      │                      │
```

---

## 三、自建码状态机（防二次兑换的核心）

```
                     preview
                        │
                        ▼
              status = 'preview'  ──────────────┐
                        │ preflight             │ 用户重新预览
                        ▼                       │（复用同一 token）
              status = 'preflight'  ◄───────────┘
                        │
                        │ ClaimPro20xOrder (条件 UPDATE)
                        ▼
              status = 'redeeming'   ← 并发第二个请求在这里被 409
                        │
                        │ CreateTask 成功
                        ▼
              status = 'running'  → 轮询回写 → completed / failed

卡密侧（独立状态机）：
     unused ──ClaimLocalCDK──► consuming ──ConsumeLocalCDK──► consumed
                                  │
                                  └──ReleaseLocalCDK──► unused（提交失败回滚）
```

**两道原子锁的真实 SQL**：

```sql
-- 会话锁
UPDATE pro20x_orders SET status='redeeming'
WHERE redeem_token=? AND status='preflight';
-- RowsAffected=1 才继续，否则 409

-- 卡密锁
UPDATE cardplatform_cdk_codes SET status='consuming'
WHERE code=? AND source='local' AND TRIM(status) IN ('','unused');
-- RowsAffected=1 才继续，否则 409
```

**服务重启兜底**：

```sql
-- 启动时把上次进程残留的 consuming 放回 unused
UPDATE cardplatform_cdk_codes SET status='unused'
WHERE source='local' AND status='consuming';
```

---

## 四、数据表

### cardplatform_cdk_codes（统一码表）

| 字段 | 说明 |
|---|---|
| upstream_id | 卡台侧 ID，自建码固定 0 |
| code | 完整卡密（UNIQUE） |
| plan | `plus` / `pro_5x` / `pro_20x` / `credit250` ... |
| payment_country | PH / EG / JP / US / CL / SG |
| status | unused / consuming / consumed / disabled |
| **source** | **`zovo` 或 `local`（分流依据）** |
| expires_at | 仅自建码，空=长期 |
| batch_no | 批次号，生成时自动 `B<时间戳>` |
| agent_id | 预留代理归属 |
| price_minor | 售价（分），仅记账 |
| bound_email | 兑换成功后写入的账号邮箱 |
| note | 备注（列表内点击编辑） |

### pro20x_orders（自建码兑换单）

| 字段 | 说明 |
|---|---|
| redeem_token | `RT-<32hex>`，前端轮询凭据 |
| cdk_code | 关联卡密 |
| client_order_id | `FB-<token 后缀>`，传给 Pro20x 做幂等 |
| task_id | Pro20x 任务 ID（`AUTO-...`） |
| status | preview / preflight / redeeming / running / succeeded / failed |
| stage | Pro20x 原始 stage |
| account_email | account-check 得到的邮箱 |
| error_code | Pro20x 错误码 |
| cleanup_status | 删卡核验结果（deleted_verified） |
| cost_points | 积分成本（1） |

---

## 五、Pro20x 接入全过程复盘（可复用作模板）

### 5.1 接入顺序（9 步）

| 步 | 动作 | 产出 | 本次耗时 |
|---|---|---|---|
| 1 | 读第三方 API 文档，拉 `/config` 确认真实能力 | 套餐表、地区、积分成本 | 30 分钟 |
| 2 | 写只读客户端（GetConfig / AccountCheck） | `internal/pro20x/client.go` | 2 小时 |
| 3 | 写套餐映射 + 请求组装 | `plan.go` + `BuildCreateRequest` | 1 小时 |
| 4 | 建兑换单表 + 状态字段 | `db/pro20x_order.go` | 1 小时 |
| 5 | 码表加 `source` + 迁移 | `db/cdk_local.go` + 幂等迁移 | 1 小时 |
| 6 | 四个公开接口分流 | `handler/cdk_local_redeem.go` | 3 小时 |
| 7 | 状态映射（第三方 → 本站） | `mapLocalStatus` | 30 分钟 |
| 8 | 防重状态机 + 原子抢占 | Claim/Release/Consume | 2 小时 |
| 9 | 配置接口 + 错误中文化 + 前端页面 | settings / friendlyError / 页面 | 4 小时 |

**核心心法：先只读、再写、最后接前端。** 每一步都能单独验证，出问题范围可控。

### 5.2 踩过的 5 个坑（一定要提前防）

#### 坑 1：轮询地址被拼成双前缀

```
Pro20x 返回 poll = "/api/v3/automate/account-check/jobs/ptj_xxx"
天真实现      = BaseURL + poll
              = "https://pro20x.com/api/v3/automate" + "/api/v3/automate/..."
              → 404 → 前端只看到 502「预检失败」
```

**解法**：URL 解析要区分三种形态。

```go
if strings.HasPrefix(p, "http")   → 原样
if strings.HasPrefix(p, "/api/")  → scheme://host + p   // 站点绝对路径
else                              → BaseURL + p          // 接口相对路径
```

**这个坑最隐蔽**：第三方返回 202 成功，但轮询 404，日志里只有一句 `http=404 msg=`。

#### 坑 2：Session 传字符串 vs 对象

第三方字段叫 `authSessionJson`，字面看像字符串。实际传字符串时对方可能解析不出邮箱：

```
{"success":false,"code":"protocol_tool_email_missing","message":"JSON 中没有识别到账号邮箱"}
```

**解法**：先 `json.Unmarshal` 成对象再传。

```go
func sessionAsJSON(raw string) any {
    var v any
    if json.Unmarshal([]byte(raw), &v) == nil && v != nil { return v }
    return raw
}
```

#### 坑 3：Cloudflare 按 UA 拦

同一接口，不同 User-Agent 的结果：

```
Go-http-client/1.1        → 202 ✅
Chrome/131                → 202 ✅
Python-urllib/3.12        → 403 error code 1010 ❌
```

**教训**：**排查时不要用 Python urllib 测第三方接口**，会被 CF 拦，误判成"接口挂了"。用 curl 或 Go。

#### 坑 4：任务终态字段不统一

同一个 API 文档里 `success` 和 `succeeded` 混用，导致轮询永远等不到终态、测试超时 180 秒。

**解法**：终态判断写成集合。

```go
func IsTerminal(s string) bool {
    switch strings.ToLower(s) {
    case "success", "succeeded", "failed", "error", "cancelled", "manual_review":
        return true
    }
    return false
}
```

#### 坑 5：前端字段名对不上（最容易被忽略）

后端返回：

```json
{ "preflight_token": "...", "account": { "email": "...", "plan_label": "Plus" } }
```

前端读的却是：

```js
body.email            // 顶层
body.currentPlan      // 驼峰
body.subscription_has_active
```

结果步骤 3 页面全空，但接口 200，没人报错。

**解法**：接新协议时，**先把前端的解析函数完整读一遍**，按它的字段名铺平返回，嵌套结构额外保留一份。

```go
c.JSON(200, gin.H{
    "preflight_token": token,
    "email": res.Email,                    // 前端读这个
    "currentPlan": acc.PlanLabel,          // 前端读这个
    "subscription_has_active": acc.HasActiveSubscription,  // 前端读这个
    "account": gin.H{...},                 // 保留一份结构化
})
```

### 5.3 新协议接入检查清单

接其他协议（Grok / Claude / Devin / 新卡台）时，照着勾：

**配置层**
- [ ] `site_settings` 加 `<协议>_enabled` / `_api_base` / `_api_key`
- [ ] `LoadConfig()` 读库，未配置时用默认 BaseURL
- [ ] 后台加 `GET/PUT /admin/<协议>/settings` + `GET /admin/<协议>/test`
- [ ] API Key 在返回时脱敏（只回 `configured` 和 hint）

**客户端层**
- [ ] `resolveURL` 处理三种路径形态（完整 URL / `/api/` 绝对路径 / 相对路径）
- [ ] `IsTerminal` 覆盖所有终态变体
- [ ] 错误结构统一解析成 `APIError{Code, Message, HTTP}`
- [ ] 请求带 `Idempotency-Key`（用 `client_order_id`）
- [ ] 超时分层：单次请求 30~60s，整体轮询 90~180s

**映射层**
- [ ] `<套餐> → planId` 映射表，未映射的直接拒绝（不要猜）
- [ ] 地区白名单 + 默认地区兜底
- [ ] 必填开关（如 `deletePaymentCard`）用指针 bool，确保一定传出去

**数据层**
- [ ] 兑换单表：token / 码 / 任务 ID / 状态 / stage / 错误码 / 清理状态
- [ ] 码表加 `source`，默认老值，**幂等迁移**
- [ ] 服务启动清理残留中间态

**分流层**
- [ ] 四个公开接口按 `source` 分流，老链路零改动
- [ ] `result-by-code` 也要支持新通道
- [ ] 预览接口加 IP 限速（防撞码）

**状态层**
- [ ] 第三方 status → 本站 status/stage 映射表
- [ ] 前端 `TERMINAL` 集合要对齐（completed / failed / cancelled）
- [ ] 轮询节流（前端 2s 轮询 ≠ 后端 2s 请求第三方）

**防重层**
- [ ] 会话锁：`preflight → redeeming` 条件 UPDATE
- [ ] 卡密锁：`unused → consuming` 条件 UPDATE
- [ ] 失败双向回滚
- [ ] 并发双击返回 409，文案「兑换正在处理中，请勿重复提交」

**前端层**
- [ ] 读一遍前端解析函数，字段名对齐
- [ ] 错误码中文化（Session 失效 / 积分不足 / 风控 / 限流）
- [ ] 管理页对齐已有页面的交互（列表、备注弹窗、批量复制）

### 5.4 本次改动的文件清单

```
新增
  backend/internal/pro20x/client.go          客户端（config/account-check/tasks）
  backend/internal/pro20x/plan.go            套餐映射 + 请求组装
  backend/internal/pro20x/client_test.go     7 个单测
  backend/internal/pro20x/plan_test.go       8 个单测
  backend/internal/pro20x/url_test.go        URL 解析单测
  backend/internal/db/cdk_local.go           自建码读写 + 原子状态机
  backend/internal/db/pro20x_order.go        兑换单读写 + 会话锁
  backend/internal/handler/cdk_local.go      生成/列表/状态/备注
  backend/internal/handler/cdk_local_redeem.go  四个接口分流 + 状态映射
  backend/internal/handler/pro20x_settings.go   通道配置与自检
  frontend/src/views/admin/LocalCDKView.vue  自建卡密管理页

修改
  backend/internal/db/db.go                  建表 + 幂等迁移 + 启动清理
  backend/internal/server/http.go            路由挂载
  backend/internal/handler/cardplatform_cdk.go  四处分流入口 + 预览限速
  backend/internal/pro20x/client.go          resolveURL 修复
  frontend/src/router/index.ts               路由
  frontend/src/layouts/AdminLayout.vue       菜单
  frontend/src/element.ts                    注册 Tickets 图标
  三个 handler 测试                          内存库补 source 列
```

### 5.5 成本与收益

| 项 | 数值 |
|---|---|
| 通道 A（卡台直充） | 按卡台单卡成本计 |
| 通道 B（协议直连） | 按协议方单次任务价格计 |
| 单均节省 | 两通道差价（按实际报价核算） |
| 规模化节省 | 单均差价 × 日单量 × 30 |
| 额外收益 | Pro20x 内置 `deletePaymentCard` + 核验，替代自建删卡 |

---

## 六、常见故障速查

| 现象 | 根因 | 排查动作 |
|---|---|---|
| 预检 502 | 轮询 URL 拼错 / 第三方超时 | 看后端日志 `[local-cdk] account-check failed` 的 `http=` 和 `msg=` |
| 预检 400「账号已有订阅」 | 正常业务判断 | 换 Free 账号 |
| 预检 401 / oauth_invalid | Session 过期 | 重新导出 `/api/auth/session` |
| 兑换 409 | 并发或重复提交 | 等待或重新预览 |
| 卡密卡在 `consuming` | 进程中断 | 重启服务自动清理 |
| 码已 `consumed` 但任务失败 | 任务受理后失败 | 查 `pro20x_orders.error_code`，人工补发 |
| 前端进度条不动 | status/stage 映射缺失 | 对齐 `mapLocalStatus` 与前端 `TERMINAL` |

---

## 七、下一步待办

1. 自建码真实兑换跑通（Free 账号）
2. 代理拿货结算（余额表 + 配额）
3. 对账页加「通道」列（zovo / pro20x）与积分成本
4. 用卡策略固化：每卡次数 = 1、自动开卡产品码 = <从卡台产品列表取>
5. `note` 列在旧库的迁移验证
