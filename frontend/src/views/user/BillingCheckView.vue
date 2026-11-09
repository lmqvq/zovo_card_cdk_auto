<template>
  <div class="u-page">
    <UserPageHeader hide-hero :title="t('redeemTabs.billing')" />
    <div class="u-main bill-wrap">
      <div class="bill-subnav">
        <div class="bill-crumb">卡密兑换 <i>/</i> <b>账单查询</b> <i>/</i> BILLING</div>
        <RedeemModeTabs compact />
      </div>

      <section class="bill-hero">
        <div class="bill-hero-copy">
          <div class="bill-eyebrow">订阅账单 <em>· BILLING LOOKUP</em></div>
          <h1>{{ t('redeemTabs.billing') }}</h1>
          <p v-if="!summary">用卡密或 Session 查询订阅状态、续费时间与历史付款账单。</p>
          <p v-else>{{ resultSubtitle }}</p>
        </div>
        <div class="bill-trust">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
          只读查询 · 不修改订阅
        </div>
      </section>

      <div class="bill-grid">
        <div class="bill-main">
          <section v-if="!summary" class="bill-card bill-query">
            <div class="bill-query-head">
              <div>
                <h2>查询订阅与账单</h2>
                <p>选择一种凭证方式开始查询</p>
              </div>
              <div class="bill-mode" role="tablist" aria-label="查询方式">
                <button type="button" :class="{ on: mode === 'cdk' }" :aria-pressed="mode === 'cdk'" @click="switchMode('cdk')">卡密查询</button>
                <button type="button" :class="{ on: mode === 'session' }" :aria-pressed="mode === 'session'" @click="switchMode('session')">Session 查询</button>
              </div>
            </div>

            <template v-if="mode === 'cdk'">
              <div class="bill-field">
                <div class="bill-field-label">
                  <label for="bill-cdk">CDK 卡密</label>
                  <span>输入后即可查询</span>
                </div>
                <div class="bill-input" :class="{ 'is-error': error && !cdkCode.trim() }">
                  <span class="bill-input-prefix">CDK</span>
                  <input
                    id="bill-cdk"
                    ref="cdkInput"
                    v-model="cdkCode"
                    autocomplete="off"
                    spellcheck="false"
                    placeholder="SXC-XXXX-XXXX-XXXX-XXXX"
                    @keyup.enter="query"
                  />
                  <button class="bill-paste" type="button" @click="pasteCdk">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="8" y="8" width="11" height="12" rx="2" /><path d="M16 8V5a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h3" /></svg>
                    粘贴
                  </button>
                </div>
                <p class="bill-hint">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 8v.5" /><path d="M12 11v5" /></svg>
                  系统会使用兑换时绑定的 Session 查询，无需再次提交 token。
                </p>
              </div>
            </template>

            <template v-else>
              <div class="bill-field">
                <div class="bill-field-label">
                  <label for="bill-session">完整 Session JSON</label>
                  <a class="bill-session-link" href="https://chatgpt.com/api/auth/session" target="_blank" rel="noopener">获取 Session ↗</a>
                </div>
                <textarea
                  id="bill-session"
                  v-model="tokenInput"
                  class="bill-textarea"
                  placeholder='{"user":{...},"accessToken":"eyJ...","sessionToken":"eyJ..."}'
                ></textarea>
                <p class="bill-hint">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 8v.5" /><path d="M12 11v5" /></svg>
                  须包含完整 sessionToken；已禁用纯 Access Token。
                </p>
              </div>
            </template>

            <div v-if="error" class="bill-alert" role="alert">{{ error }}</div>

            <button class="btn-primary bill-query-btn" type="button" :disabled="loading" @click="query">
              {{ loading ? '查询中…' : '查询订阅与账单' }}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12h14" /><path d="m13 6 6 6-6 6" /></svg>
            </button>
            <p class="bill-readonly">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
              只读查询，不会修改订阅、不会产生扣款。
            </p>
          </section>

          <section v-if="!summary" class="bill-help bill-help--wide">
            <h3>查询前准备</h3>
            <ol>
              <li><span class="bill-num">01</span><div><b>卡密查询</b><p>卡密需在兑换时使用 Session 预检，系统才能读取绑定账单。</p></div></li>
              <li><span class="bill-num">02</span><div><b>Session 查询</b><p>打开 chatgpt.com/api/auth/session，复制完整 JSON（含 sessionToken）。</p></div></li>
              <li><span class="bill-num">03</span><div><b>查询内容</b><p>订阅状态、续费时间、历史账单与账单链接。</p></div></li>
            </ol>
          </section>

          <div v-else class="bill-result">
            <div class="bill-context">
              <span class="bill-context-label">查询条件</span>
              <span class="bill-context-code">{{ queryCode }}</span>
              <span class="bill-source">{{ sourceText }}</span>
              <button class="bill-change" type="button" @click="resetQuery">更换卡密</button>
            </div>

            <section class="bill-card bill-summary">
              <div class="bill-section-head">
                <h2>订阅摘要</h2>
                <span class="bill-status" :class="statusTone">{{ statusText }}</span>
              </div>
              <div class="bill-summary-top">
                <div>
                  <div class="bill-product">{{ planText }}</div>
                  <div class="bill-account">{{ displayEmail }}</div>
                </div>
                <span class="bill-renew" :class="renewTone">{{ renewText }}</span>
              </div>
              <dl class="bill-stats">
                <div><dt>计费周期</dt><dd>{{ billingText }}</dd></div>
                <div><dt>下次续费</dt><dd class="mono">{{ renewalDateText }}</dd></div>
                <div><dt>订阅到期</dt><dd class="mono">{{ expiryDateText }}</dd></div>
                <div><dt>剩余</dt><dd>{{ remainingText }}</dd></div>
              </dl>
              <div v-if="periodStartDate" class="bill-period">
                <div class="bill-period-head"><span>当前周期</span><span>{{ periodRange }}</span></div>
                <div class="bill-bar"><i :style="{ width: periodProgress + '%' }"></i></div>
              </div>
            </section>

            <section class="bill-card bill-invoices">
              <div class="bill-section-head">
                <h2>付款账单<span class="bill-count">{{ invoices?.length || 0 }} 笔</span></h2>
              </div>
              <div v-if="invoices && invoices.length" class="bill-table">
                <div class="bill-table-head"><span>日期</span><span>项目</span><span>金额</span><span>状态</span><span>操作</span></div>
                <div v-for="inv in invoices" :key="inv.id || inv.number || inv.created" class="bill-invoice-row">
                  <span class="bill-date">{{ invoiceDate(inv) }}</span>
                  <div class="bill-item">
                    <b>{{ invoiceTitle(inv) }}</b>
                    <small v-if="invoiceNo(inv)">{{ invoiceNo(inv) }}</small>
                  </div>
                  <span class="bill-amount">{{ formatAmount(inv.total, inv.currency) }}</span>
                  <span class="bill-paid" :class="invoiceStatusTone(inv)">{{ invoiceStatus(inv) }}</span>
                  <div class="bill-links">
                    <a v-if="inv.hosted_invoice_url" :href="inv.hosted_invoice_url" target="_blank" rel="noopener">账单 ↗</a>
                    <a v-if="inv.invoice_pdf" :href="inv.invoice_pdf" target="_blank" rel="noopener">PDF ↗</a>
                    <span v-if="!inv.hosted_invoice_url && !inv.invoice_pdf">—</span>
                  </div>
                </div>
              </div>
              <p v-else class="bill-empty">暂无付款账单记录。</p>
              <p class="bill-note">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 8v.5" /><path d="M12 11v5" /></svg>
                账单数据来自 ChatGPT 账单页，可能有几分钟延迟。
              </p>
            </section>
          </div>
        </div>

        <aside class="bill-aside">
          <template v-if="!summary">
            <section class="bill-dark">
              <div class="bill-dark-label">READ ONLY</div>
              <div class="bill-dark-lock">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
              </div>
              <h3>只读查询</h3>
              <p>不会修改订阅，也不会发起任何扣款。</p>
              <ul class="bill-dark-list">
                <li>不写入账号设置</li>
                <li>不触发支付动作</li>
                <li>只读取账单数据</li>
              </ul>
            </section>

            <section class="bill-help bill-help--compact">
              <h3>查不到账单？</h3>
              <div class="bill-faq"><b>卡密没有绑定 Session</b><p>兑换时需使用 Session 凭证，或改用 Session 查询。</p></div>
              <div class="bill-faq"><b>Session 已失效</b><p>重新登录 ChatGPT，复制新的完整 JSON。</p></div>
            </section>
          </template>

          <template v-else>
            <section class="bill-dark bill-query-info">
              <div class="bill-dark-label">BILLING LOOKUP</div>
              <div class="bill-dark-status">
                <span class="bill-dark-ring">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m4 12.5 5 5L20 6.5" /></svg>
                </span>
                <div><b>查询成功</b><span>{{ queryTimeText }}</span></div>
              </div>
              <div class="bill-info-email" :class="{ 'is-muted': !hasEmail }">{{ displayEmail }}</div>
              <dl class="bill-info-rows">
                <div><dt>SOURCE</dt><dd>{{ sourceText }}</dd></div>
                <div><dt>QUERY</dt><dd>{{ queryTimeText }}</dd></div>
                <div><dt>MODE</dt><dd>READ ONLY</dd></div>
              </dl>
              <button class="bill-dark-btn" type="button" :disabled="!hasEmail" @click="copyAccount">复制账号</button>
              <button class="bill-dark-link" type="button" @click="resetQuery">重新查询 ↗</button>
            </section>

            <section class="bill-help bill-help--compact">
              <h3>账单没有显示？</h3>
              <div class="bill-faq"><b>账单可能延迟 1–5 分钟</b><p>刷新页面或稍后重新查询。</p></div>
              <div class="bill-faq"><b>Session 已失效</b><p>重新登录 ChatGPT，复制新的完整 Session。</p></div>
            </section>

            <p class="bill-aside-note">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 8v.5" /><path d="M12 11v5" /></svg>
              账单查询不会修改订阅或产生扣款，查询结果仅在本页展示。
            </p>
          </template>
        </aside>
      </div>

      <footer class="bill-foot">
        <span>LLDCard · 订阅账单查询</span>
        <span class="bill-foot-mono">BILLING LOOKUP · READ ONLY</span>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import UserPageHeader from '../../components/UserPageHeader.vue'
import RedeemModeTabs from '../../components/RedeemModeTabs.vue'
import { planLabel as displayPlanLabel } from '../../lib/plan'
import { copyToClipboard } from '../../lib/clipboard'
import { dialog } from '../../lib/dialog'

const { t } = useI18n({ useScope: 'global' })
const route = useRoute()

const mode = ref<'cdk' | 'session'>('cdk')
const cdkCode = ref('')
const tokenInput = ref('')
const cdkInput = ref<HTMLInputElement | null>(null)
const loading = ref(false)
const error = ref('')
const summary = ref<any>(null)
const invoices = ref<any[] | null>(null)
const lastQueryAt = ref('')
const lastQueryMode = ref<'cdk' | 'session'>('cdk')
const lastCdkCode = ref('')

function boolValue(value: any): boolean | null {
  if (value === true || value === 'true' || value === 1 || value === '1') return true
  if (value === false || value === 'false' || value === 0 || value === '0') return false
  return null
}

const planRaw = computed(() => String(summary.value?.plan_type || summary.value?.subscription_plan || '').trim())
const planText = computed(() => {
  if (!planRaw.value) return '未知套餐'
  if (planRaw.value.toLowerCase() === 'free') return '免费版'
  const label = displayPlanLabel(planRaw.value)
  const raw = planRaw.value.toLowerCase()
  const isChatGPT = /^(chatgpt|plus|pro|team)/.test(raw)
    || label === 'Plus'
    || label.startsWith('Pro')
    || label.startsWith('Team')
  return isChatGPT && !label.startsWith('ChatGPT') ? `ChatGPT ${label}` : label
})

const activeState = computed(() => boolValue(summary.value?.has_active_subscription))
const statusText = computed(() => (activeState.value === true ? '有效' : activeState.value === false ? '已到期' : '未知'))
const statusTone = computed(() => (activeState.value === true ? 'is-good' : activeState.value === false ? 'is-warn' : 'is-muted'))

const renewState = computed(() => boolValue(summary.value?.will_renew))
const renewText = computed(() => (renewState.value === true ? '自动续费 已开启' : renewState.value === false ? '自动续费 已关闭' : '自动续费 未知'))
const renewTone = computed(() => (renewState.value === true ? 'is-on' : renewState.value === false ? 'is-off' : 'is-muted'))

const displayEmail = computed(() => String(summary.value?.email || '').trim() || '账号未返回')
const hasEmail = computed(() => !!String(summary.value?.email || '').trim())

function toDate(value: any): Date | null {
  if (value === null || value === undefined || value === '') return null
  if (value instanceof Date) return Number.isNaN(value.getTime()) ? null : value
  const text = String(value).trim()
  if (!text) return null
  if (/^\d{10,13}$/.test(text)) {
    const number = Number(text)
    const date = new Date(number > 1e12 ? number : number * 1000)
    return Number.isNaN(date.getTime()) ? null : date
  }
  const date = new Date(text)
  return Number.isNaN(date.getTime()) ? null : date
}

function pad2(value: number) {
  return String(value).padStart(2, '0')
}

function formatDateValue(value: any) {
  const date = toDate(value)
  if (!date) return value ? String(value) : '—'
  return `${date.getFullYear()}/${pad2(date.getMonth() + 1)}/${pad2(date.getDate())}`
}

function formatDateTime(value: any) {
  const date = toDate(value)
  if (!date) return '—'
  return `${formatDateValue(date)} ${pad2(date.getHours())}:${pad2(date.getMinutes())}`
}

function formatAmount(total: any, currency: any) {
  if (total === null || total === undefined || total === '') return '—'
  const amount = Number(total)
  if (!Number.isFinite(amount)) return '—'
  const cur = String(currency || 'usd').toUpperCase()
  try {
    return new Intl.NumberFormat('zh-CN', { style: 'currency', currency: cur, minimumFractionDigits: 2 }).format(amount / 100)
  } catch {
    return `${(amount / 100).toFixed(2)} ${cur}`
  }
}

function periodLabel() {
  const value = String(summary.value?.billing_period || '').toLowerCase()
  if (value.includes('year') || value === 'annual') return '年'
  if (value.includes('month')) return '月'
  return value
}

const latestInvoice = computed(() => (invoices.value && invoices.value.length ? invoices.value[0] : null))
const billingText = computed(() => {
  const invoice = latestInvoice.value
  const amount = invoice ? formatAmount(invoice.total, invoice.currency) : ''
  const period = periodLabel()
  if (amount && amount !== '—') return period ? `${amount} / ${period}` : amount
  const currency = summary.value?.billing_currency
  if (currency) return `${String(currency).toUpperCase()}${period ? ` / ${period}` : ''}`
  return period || '—'
})

const renewalRaw = computed(() => summary.value?.renews_at || summary.value?.active_until || summary.value?.expires_at || '')
const expiryRaw = computed(() => summary.value?.active_until || summary.value?.expires_at || summary.value?.renews_at || '')
const renewalDateText = computed(() => formatDateValue(renewalRaw.value))
const expiryDateText = computed(() => formatDateValue(expiryRaw.value))

const periodEndDate = computed(() => toDate(expiryRaw.value))
const periodStartDate = computed(() => {
  const end = periodEndDate.value
  if (!end) return null
  const start = new Date(end)
  const value = String(summary.value?.billing_period || '').toLowerCase()
  if (value.includes('year') || value === 'annual') start.setFullYear(start.getFullYear() - 1)
  else if (value.includes('month')) start.setMonth(start.getMonth() - 1)
  else start.setDate(start.getDate() - 30)
  return start
})

const periodRange = computed(() => {
  const start = periodStartDate.value
  const end = periodEndDate.value
  return start && end ? `${formatDateValue(start)} — ${formatDateValue(end)}` : ''
})

const periodProgress = computed(() => {
  const start = periodStartDate.value
  const end = periodEndDate.value
  if (!start || !end) return 0
  const total = end.getTime() - start.getTime()
  if (total <= 0) return 100
  const elapsed = Date.now() - start.getTime()
  return Math.min(100, Math.max(8, Math.round((elapsed / total) * 100)))
})

const remainingText = computed(() => {
  const end = periodEndDate.value
  if (!end) return '—'
  const days = Math.ceil((end.getTime() - Date.now()) / 86400000)
  if (days < 0) return '已到期'
  if (days === 0) return '今天'
  return `${days} 天`
})

const resultSubtitle = computed(() => {
  const count = invoices.value?.length || 0
  if (statusText.value === '有效') {
    return count ? `当前订阅有效，已找到 ${count} 笔付款记录。` : '当前订阅有效，暂无付款账单记录。'
  }
  if (statusText.value === '已到期') {
    return count ? `当前订阅已到期，已找到 ${count} 笔付款记录。` : '当前订阅已到期。'
  }
  return count ? `已查询到订阅信息与 ${count} 笔付款记录。` : '已查询到订阅信息。'
})

const queryCode = computed(() => (lastQueryMode.value === 'cdk' ? (lastCdkCode.value || '—') : 'Session 凭证'))
const sourceText = computed(() => (lastQueryMode.value === 'cdk' ? '卡密绑定' : 'Session'))
const queryTimeText = computed(() => (lastQueryAt.value ? formatDateTime(lastQueryAt.value) : '—'))

function invoiceDate(inv: any) {
  return formatDateValue(inv?.created)
}

function invoiceNo(inv: any) {
  return String(inv?.number || inv?.invoice_number || inv?.id || '').trim()
}

function invoiceTitle(inv: any) {
  const description = String(inv?.description || '').trim()
  if (description) return description
  const period = periodLabel()
  return `${planText.value}${period ? ` · ${period}付` : ''}`
}

function invoiceStatus(inv: any) {
  return boolValue(inv?.paid) === true ? '已支付' : (inv?.status || '未支付')
}

function invoiceStatusTone(inv: any) {
  return boolValue(inv?.paid) === true ? 'is-good' : 'is-warn'
}

function switchMode(next: 'cdk' | 'session') {
  if (mode.value === next) return
  mode.value = next
  error.value = ''
  void nextTick(() => {
    if (next === 'cdk') cdkInput.value?.focus()
  })
}

async function pasteCdk() {
  try {
    const text = await navigator.clipboard?.readText()
    if (!text?.trim()) {
      error.value = '剪贴板为空，请手动粘贴。'
      return
    }
    cdkCode.value = text.trim()
    error.value = ''
  } catch {
    error.value = '无法读取剪贴板，请手动粘贴。'
  }
}

async function copyAccount() {
  if (!hasEmail.value) return
  const ok = await copyToClipboard(String(summary.value?.email || ''))
  dialog.toast(ok ? '账号已复制' : '复制失败，请手动选择文本', ok ? 'ok' : 'warn')
}

async function query() {
  if (loading.value) return
  const cdk = cdkCode.value.trim()
  const token = tokenInput.value.trim()
  if (mode.value === 'cdk' && !cdk) {
    error.value = '请输入 CDK 卡密。'
    cdkInput.value?.focus()
    return
  }
  if (mode.value === 'session' && !token) {
    error.value = '请粘贴完整 Session JSON（须含 sessionToken）。'
    return
  }

  error.value = ''
  loading.value = true
  lastQueryMode.value = mode.value
  lastCdkCode.value = mode.value === 'cdk' ? cdk : ''
  lastQueryAt.value = ''
  summary.value = null
  invoices.value = null

  try {
    const body = mode.value === 'cdk'
      ? { cdk_code: cdk }
      : { token_input: token }
    const response = await fetch('/api/v1/public/billing/check', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      error.value = data.error || '查询失败'
      return
    }
    summary.value = data.summary || {}
    invoices.value = Array.isArray(data.invoices) ? data.invoices : []
    lastQueryAt.value = new Date().toISOString()
  } catch (cause: any) {
    error.value = cause?.message || '网络错误'
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  summary.value = null
  invoices.value = null
  error.value = ''
  lastQueryAt.value = ''
  mode.value = lastQueryMode.value
  void nextTick(() => {
    if (mode.value === 'cdk') cdkInput.value?.focus()
  })
}

onMounted(() => {
  const cdk = String(route.query.cdk || '').trim()
  if (cdk) {
    cdkCode.value = cdk
    mode.value = 'cdk'
    void query()
  }
})
</script>

<style scoped>
.bill-wrap {
  max-width: 72rem;
  margin-left: auto;
  margin-right: auto;
  padding: 0 24px;
}
.bill-subnav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 8px 0 18px;
  border-bottom: 1px solid var(--brd);
}
.bill-crumb {
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  white-space: nowrap;
}
.bill-crumb b { color: var(--ink); font-weight: 600; }
.bill-crumb i { margin: 0 8px; font-style: normal; opacity: 0.65; }

.bill-hero {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 36px;
  padding: 34px 0 30px;
}
.bill-hero-copy { min-width: 0; }
.bill-eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}
.bill-eyebrow::before { content: ''; width: 18px; height: 1px; background: var(--ink-3); }
.bill-eyebrow em { color: var(--ink); font-style: normal; }
.bill-hero h1 {
  margin: 13px 0 0;
  color: var(--ink);
  font-family: var(--title-font, var(--font-display));
  font-size: clamp(2.2rem, 4.6vw, 3.375rem);
  font-style: var(--title-style, normal);
  font-weight: 800;
  letter-spacing: var(--title-track, -0.035em);
  line-height: 1.04;
}
.bill-hero p {
  max-width: 34rem;
  margin: 14px 0 0;
  color: var(--ink-2);
  font-size: 14px;
  line-height: 1.75;
}
.bill-trust {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  flex: none;
  margin-bottom: 4px;
  padding: 8px 13px;
  border: 1px solid var(--brd-2);
  border-radius: 999px;
  color: var(--ink-2);
  font-size: 11.5px;
  white-space: nowrap;
}
.bill-trust svg { width: 13px; height: 13px; color: var(--good); }

.bill-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 22px;
  align-items: start;
  padding-bottom: 28px;
}
.bill-main,
.bill-aside,
.bill-result {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.bill-card {
  background: var(--surface);
  border: 1px solid var(--brd);
  border-radius: 22px;
  box-shadow: 0 1px 0 color-mix(in srgb, var(--surface) 90%, transparent) inset, 0 24px 50px -46px rgba(0, 0, 0, 0.35);
}
.bill-query { padding: 28px 30px 24px; }
.bill-query-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
}
.bill-query-head h2 { margin: 1px 0 0; color: var(--ink); font-size: 22px; font-weight: 800; letter-spacing: -0.02em; }
.bill-query-head p { margin: 6px 0 0; color: var(--ink-2); font-size: 12.5px; }
.bill-mode {
  display: flex;
  flex: none;
  gap: 2px;
  padding: 4px;
  border: 1px solid var(--brd);
  border-radius: 999px;
  background: var(--surface-2);
}
.bill-mode button {
  height: 32px;
  padding: 0 16px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: var(--ink-2);
  font: inherit;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}
.bill-mode button.on { background: var(--ink); color: var(--bg); box-shadow: 0 2px 6px rgba(0, 0, 0, 0.16); }

.bill-field { margin-top: 26px; }
.bill-field-label { display: flex; align-items: baseline; justify-content: space-between; gap: 16px; }
.bill-field-label label { color: var(--ink); font-size: 12px; font-weight: 700; }
.bill-field-label > span { color: var(--ink-3); font-family: var(--font-mono); font-size: 9.5px; letter-spacing: 0.06em; }
.bill-input {
  display: flex;
  align-items: center;
  margin-top: 9px;
  overflow: hidden;
  border: 1px solid var(--brd);
  border-radius: 14px;
  background: var(--surface-2);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
.bill-input:focus-within { border-color: var(--ink); box-shadow: 0 0 0 3px color-mix(in srgb, var(--ink) 7%, transparent); }
.bill-input.is-error { border-color: var(--err); }
.bill-input-prefix {
  align-self: stretch;
  display: flex;
  align-items: center;
  padding: 0 13px;
  border-right: 1px solid var(--brd);
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.06em;
}
.bill-input input {
  flex: 1;
  min-width: 0;
  height: 52px;
  padding: 0 14px;
  border: 0;
  outline: none;
  background: transparent;
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 13.5px;
}
.bill-input input::placeholder { color: var(--ink-3); opacity: 0.7; }
.bill-paste {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 30px;
  margin-right: 10px;
  padding: 0 10px;
  border: 1px solid var(--brd);
  border-radius: 8px;
  background: var(--surface);
  color: var(--ink-2);
  font: inherit;
  font-size: 10.5px;
  cursor: pointer;
}
.bill-paste svg { width: 12px; height: 12px; }
.bill-paste:hover { color: var(--ink); border-color: var(--brd-2); }
.bill-textarea {
  width: 100%;
  height: 140px;
  margin-top: 9px;
  padding: 12px 14px;
  resize: vertical;
  border: 1px solid var(--brd);
  border-radius: 14px;
  outline: none;
  background: var(--surface-2);
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 11.5px;
  line-height: 1.6;
}
.bill-textarea:focus { border-color: var(--ink); box-shadow: 0 0 0 3px color-mix(in srgb, var(--ink) 7%, transparent); }
.bill-session-link { color: var(--ink-2); font-size: 10.5px; text-decoration: none; border-bottom: 1px solid var(--brd-2); padding-bottom: 1px; }
.bill-hint {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  margin: 9px 0 0;
  color: var(--ink-3);
  font-size: 11px;
  line-height: 1.55;
}
.bill-hint svg { width: 12px; height: 12px; margin-top: 2px; flex: none; }
.bill-alert {
  margin-top: 16px;
  padding: 10px 12px;
  border: 1px solid color-mix(in srgb, var(--err) 40%, transparent);
  border-radius: 10px;
  background: var(--err-soft);
  color: var(--err);
  font-size: 11.5px;
  line-height: 1.5;
}
.bill-query-btn {
  width: 100%;
  height: 48px;
  margin-top: 22px;
  padding: 0 20px;
  border-radius: 999px;
  font-size: 13px;
}
.bill-query-btn svg { width: 14px; height: 14px; }
.bill-readonly {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin: 12px 0 0;
  color: var(--ink-3);
  font-size: 10.5px;
}
.bill-readonly svg { width: 11px; height: 11px; }

.bill-help {
  border: 1px solid var(--brd);
  border-radius: 22px;
  background: var(--surface);
  padding: 22px 22px 8px;
}
.bill-help h3 { margin: 0 0 2px; color: var(--ink); font-size: 12px; font-weight: 800; letter-spacing: 0.02em; }
.bill-help ol { margin: 0; padding: 0; list-style: none; }
.bill-help li { display: flex; align-items: flex-start; gap: 12px; padding: 14px 0; border-top: 1px solid color-mix(in srgb, var(--brd) 60%, transparent); }
.bill-help li:first-child { border-top: 0; }
.bill-num { flex: none; padding-top: 2px; color: var(--ink-3); font-family: var(--font-mono); font-size: 9.5px; letter-spacing: 0.08em; }
.bill-help b { display: block; color: var(--ink); font-size: 12.5px; font-weight: 700; }
.bill-help p { margin: 3px 0 0; color: var(--ink-3); font-size: 11px; line-height: 1.55; }
.bill-help--wide { padding: 22px 24px 20px; }
.bill-help--wide ol { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 0; margin-top: 14px; }
.bill-help--wide li { padding: 0; border-top: 0; }
.bill-help--wide li + li { padding-left: 22px; border-left: 1px solid color-mix(in srgb, var(--brd) 60%, transparent); }

.bill-dark {
  position: relative;
  overflow: hidden;
  padding: 24px;
  border-radius: 22px;
  background:
    radial-gradient(360px 200px at 18% -22%, color-mix(in srgb, var(--bg) 13%, transparent), transparent 62%),
    radial-gradient(280px 200px at 108% 112%, color-mix(in srgb, var(--good) 16%, transparent), transparent 64%),
    var(--ink);
  color: var(--bg);
  box-shadow: 0 30px 60px -40px rgba(0, 0, 0, 0.65);
}
.bill-dark-label { color: color-mix(in srgb, var(--bg) 48%, transparent); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.18em; }
.bill-dark-lock {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  margin-top: 20px;
  border: 1px solid color-mix(in srgb, var(--bg) 18%, transparent);
  border-radius: 50%;
  background: color-mix(in srgb, var(--bg) 7%, transparent);
  color: var(--good);
}
.bill-dark-lock svg { width: 15px; height: 15px; }
.bill-dark h3 { margin: 16px 0 0; font-size: 20px; font-weight: 800; letter-spacing: -0.02em; }
.bill-dark > p { margin: 6px 0 0; color: color-mix(in srgb, var(--bg) 55%, transparent); font-size: 11.5px; line-height: 1.6; }
.bill-dark-list { margin: 18px 0 0; padding: 16px 0 0; border-top: 1px solid color-mix(in srgb, var(--bg) 10%, transparent); list-style: none; }
.bill-dark-list li { display: flex; align-items: center; gap: 8px; padding: 6px 0; color: color-mix(in srgb, var(--bg) 72%, transparent); font-size: 11px; }
.bill-dark-list li::before { content: ''; width: 4px; height: 4px; border-radius: 50%; background: var(--good); }

.bill-help--compact { padding: 20px 22px 8px; }
.bill-faq { padding: 13px 0; border-top: 1px solid color-mix(in srgb, var(--brd) 60%, transparent); }
.bill-faq:first-of-type { margin-top: 6px; }
.bill-faq b { display: block; color: var(--ink); font-size: 12px; font-weight: 700; }
.bill-faq p { margin: 3px 0 0; color: var(--ink-3); font-size: 10.5px; line-height: 1.55; }

.bill-context {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border: 1px solid var(--brd);
  border-radius: 14px;
  background: color-mix(in srgb, var(--surface) 70%, transparent);
}
.bill-context-label { color: var(--ink-3); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.12em; }
.bill-context-code { overflow: hidden; color: var(--ink); font-family: var(--font-mono); font-size: 11.5px; text-overflow: ellipsis; white-space: nowrap; }
.bill-source { flex: none; padding: 4px 8px; border-radius: 999px; background: var(--surface-2); color: var(--ink-2); font-size: 10px; }
.bill-change { margin-left: auto; padding: 0 0 2px; border: 0; border-bottom: 1px solid var(--brd-2); background: transparent; color: var(--ink-2); font: inherit; font-size: 10.5px; white-space: nowrap; cursor: pointer; }

.bill-summary { padding: 26px 28px 24px; }
.bill-section-head { display: flex; align-items: baseline; justify-content: space-between; gap: 16px; }
.bill-section-head h2 { margin: 0; color: var(--ink); font-size: 13px; font-weight: 800; letter-spacing: 0.02em; }
.bill-count { margin-left: 4px; color: var(--ink-3); font-family: var(--font-mono); font-size: 10px; font-weight: 400; }
.bill-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 24px;
  padding: 0 10px;
  border: 1px solid transparent;
  border-radius: 999px;
  font-size: 10.5px;
  font-weight: 600;
}
.bill-status::before { content: ''; width: 5px; height: 5px; border-radius: 50%; background: currentColor; }
.bill-status.is-good { color: var(--good); background: var(--good-soft); border-color: color-mix(in srgb, var(--good) 28%, transparent); }
.bill-status.is-warn { color: var(--warn); background: var(--warn-soft); border-color: color-mix(in srgb, var(--warn) 28%, transparent); }
.bill-status.is-muted { color: var(--ink-3); background: var(--surface-2); border-color: var(--brd); }
.bill-summary-top { display: flex; align-items: flex-start; justify-content: space-between; gap: 20px; margin-top: 18px; }
.bill-product { color: var(--ink); font-size: 27px; font-weight: 800; letter-spacing: -0.03em; line-height: 1.1; }
.bill-account { margin-top: 7px; color: var(--ink-2); font-family: var(--font-mono); font-size: 11.5px; overflow-wrap: anywhere; }
.bill-renew { flex: none; display: inline-flex; align-items: center; gap: 6px; padding: 6px 10px; border-radius: 999px; background: var(--surface-2); color: var(--ink-2); font-size: 10.5px; white-space: nowrap; }
.bill-renew::before { content: ''; width: 5px; height: 5px; border-radius: 50%; background: var(--ink-3); }
.bill-renew.is-on::before { background: var(--good); }
.bill-renew.is-off::before { background: var(--warn); }
.bill-renew.is-muted { color: var(--ink-3); }
.bill-stats { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 18px; margin: 24px 0 0; padding-top: 20px; border-top: 1px solid var(--brd); }
.bill-stats dt { color: var(--ink-3); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.1em; text-transform: uppercase; }
.bill-stats dd { margin: 6px 0 0; color: var(--ink); font-size: 14px; font-weight: 700; }
.bill-stats dd.mono { font-family: var(--font-mono); font-size: 12px; font-weight: 600; }
.bill-period { margin-top: 22px; padding-top: 16px; border-top: 1px dashed var(--brd); }
.bill-period-head { display: flex; align-items: center; justify-content: space-between; gap: 16px; color: var(--ink-3); font-size: 10.5px; }
.bill-period-head span:last-child { font-family: var(--font-mono); font-size: 10.5px; }
.bill-bar { height: 5px; margin-top: 9px; overflow: hidden; border-radius: 999px; background: var(--surface-2); }
.bill-bar i { display: block; height: 100%; border-radius: 999px; background: var(--ink); }

.bill-invoices { padding: 24px 26px 20px; }
.bill-table { margin-top: 12px; }
.bill-table-head, .bill-invoice-row { display: grid; grid-template-columns: 88px minmax(0, 1fr) 96px 72px 92px; gap: 12px; align-items: center; }
.bill-table-head { padding: 10px 12px; color: var(--ink-3); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.12em; text-transform: uppercase; }
.bill-invoice-row { padding: 16px 12px; margin: 0 -12px; border-top: 1px solid var(--brd); border-radius: 10px; transition: background 0.15s ease; }
.bill-invoice-row:hover { background: color-mix(in srgb, var(--surface-2) 60%, transparent); }
.bill-date { color: var(--ink-2); font-family: var(--font-mono); font-size: 11px; }
.bill-item b { display: block; color: var(--ink); font-size: 12.5px; font-weight: 600; }
.bill-item small { display: block; margin-top: 3px; color: var(--ink-3); font-family: var(--font-mono); font-size: 9.5px; }
.bill-amount { color: var(--ink); font-family: var(--font-mono); font-size: 13px; font-weight: 600; }
.bill-paid { justify-self: start; padding: 4px 8px; border: 1px solid transparent; border-radius: 999px; font-size: 9.5px; font-weight: 600; white-space: nowrap; }
.bill-paid.is-good { color: var(--good); background: var(--good-soft); border-color: color-mix(in srgb, var(--good) 24%, transparent); }
.bill-paid.is-warn { color: var(--warn); background: var(--warn-soft); border-color: color-mix(in srgb, var(--warn) 24%, transparent); }
.bill-links { display: flex; gap: 10px; }
.bill-links a { padding-bottom: 1px; border-bottom: 1px solid var(--brd-2); color: var(--ink-2); font-size: 10.5px; text-decoration: none; }
.bill-links span { color: var(--ink-3); font-size: 10.5px; }
.bill-empty { margin: 18px 0 0; color: var(--ink-3); font-size: 12px; }
.bill-note { display: flex; align-items: flex-start; gap: 6px; margin: 16px 0 0; padding-top: 14px; border-top: 1px dashed var(--brd); color: var(--ink-3); font-size: 10.5px; line-height: 1.55; }
.bill-note svg { width: 11px; height: 11px; margin-top: 2px; flex: none; }

.bill-query-info { padding: 24px 24px 20px; }
.bill-dark-status { display: flex; align-items: center; gap: 10px; margin-top: 20px; }
.bill-dark-ring { display: grid; place-items: center; flex: none; width: 34px; height: 34px; border-radius: 50%; background: var(--good); color: #fff; }
.bill-dark-ring svg { width: 15px; height: 15px; }
.bill-dark-status b { display: block; font-size: 15px; font-weight: 700; }
.bill-dark-status span { display: block; margin-top: 3px; color: color-mix(in srgb, var(--bg) 45%, transparent); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.08em; }
.bill-info-email { margin-top: 18px; padding-bottom: 16px; border-bottom: 1px dashed color-mix(in srgb, var(--bg) 18%, transparent); color: color-mix(in srgb, var(--bg) 88%, transparent); font-family: var(--font-mono); font-size: 11px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.bill-info-email.is-muted { color: color-mix(in srgb, var(--bg) 45%, transparent); }
.bill-info-rows { display: flex; flex-direction: column; gap: 12px; margin: 16px 0 0; }
.bill-info-rows > div { display: flex; align-items: baseline; justify-content: space-between; gap: 12px; }
.bill-info-rows dt { color: color-mix(in srgb, var(--bg) 42%, transparent); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.1em; }
.bill-info-rows dd { margin: 0; color: color-mix(in srgb, var(--bg) 82%, transparent); font-family: var(--font-mono); font-size: 10.5px; text-align: right; }
.bill-dark-btn { display: flex; align-items: center; justify-content: center; width: 100%; height: 39px; margin-top: 20px; border: 0; border-radius: 999px; background: var(--bg); color: var(--ink); font: inherit; font-size: 12px; font-weight: 700; cursor: pointer; }
.bill-dark-btn:disabled { opacity: 0.45; cursor: not-allowed; }
.bill-dark-link { display: block; width: max-content; margin: 13px auto 0; padding: 0 0 1px; border: 0; border-bottom: 1px solid color-mix(in srgb, var(--bg) 22%, transparent); background: transparent; color: color-mix(in srgb, var(--bg) 60%, transparent); font: inherit; font-size: 10.5px; cursor: pointer; }
.bill-aside-note { display: flex; align-items: flex-start; gap: 8px; padding: 2px 4px; color: var(--ink-3); font-size: 11px; line-height: 1.65; }
.bill-aside-note svg { width: 12px; height: 12px; margin-top: 2px; flex: none; }

.bill-foot { display: flex; align-items: center; justify-content: space-between; gap: 24px; padding: 18px 0 34px; border-top: 1px solid var(--brd); color: var(--ink-3); font-size: 10.5px; }
.bill-foot-mono { font-family: var(--font-mono); font-size: 9.5px; letter-spacing: 0.08em; }

@media (max-width: 900px) {
  .bill-grid { grid-template-columns: minmax(0, 1fr); }
  .bill-aside { order: 2; }
  .bill-stats { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .bill-hero { flex-direction: column; align-items: stretch; gap: 18px; }
  .bill-trust { align-self: flex-start; margin: 0; }
}
@media (min-width: 641px) and (max-width: 900px) {
  .bill-aside { max-width: 460px; }
}
@media (max-width: 640px) {
  .bill-wrap { padding: 0 16px; }
  .bill-subnav { padding: 4px 0 14px; }
  .bill-crumb { display: none; }
  .bill-hero { gap: 16px; padding: 22px 0; }
  .bill-hero h1 { margin-top: 10px; font-size: 38px; letter-spacing: -0.04em; }
  .bill-hero p { margin-top: 10px; font-size: 12.5px; line-height: 1.7; }
  .bill-trust { font-size: 10.5px; padding: 7px 10px; }
  .bill-query { padding: 18px 16px 16px; border-radius: 18px; }
  .bill-query-head { flex-direction: column; align-items: stretch; gap: 14px; }
  .bill-query-head h2 { font-size: 18px; }
  .bill-query-head p { font-size: 11px; }
  .bill-mode { width: 100%; }
  .bill-mode button { flex: 1; height: 30px; padding: 0; font-size: 11.5px; }
  .bill-field { margin-top: 20px; }
  .bill-input input { height: 48px; font-size: 12.5px; }
  .bill-paste { margin-right: 8px; padding: 0 8px; }
  .bill-query-btn { height: 44px; margin-top: 18px; font-size: 12.5px; }
  .bill-help { padding: 18px 18px 6px; border-radius: 18px; }
  .bill-help--wide { padding: 18px 18px 16px; }
  .bill-help--wide ol { grid-template-columns: 1fr; margin-top: 8px; }
  .bill-help--wide li { padding: 14px 0; border-top: 1px solid color-mix(in srgb, var(--brd) 60%, transparent); }
  .bill-help--wide li:first-child { border-top: 0; }
  .bill-help--wide li + li { padding-left: 0; border-left: 0; }
  .bill-context { flex-wrap: wrap; gap: 8px 10px; padding: 11px 12px; }
  .bill-context-code { max-width: 100%; font-size: 10.5px; }
  .bill-change { margin-left: auto; }
  .bill-summary { padding: 18px 16px 16px; border-radius: 18px; }
  .bill-summary-top { flex-direction: column; gap: 12px; }
  .bill-product { font-size: 23px; }
  .bill-account { font-size: 10.5px; }
  .bill-stats { gap: 16px 14px; margin-top: 18px; padding-top: 16px; }
  .bill-stats dt { font-size: 8.5px; }
  .bill-stats dd { font-size: 12.5px; }
  .bill-stats dd.mono { font-size: 10.5px; }
  .bill-period { margin-top: 18px; padding-top: 14px; }
  .bill-invoices { padding: 18px 16px 16px; border-radius: 18px; }
  .bill-table-head { display: none; }
  .bill-table { margin-top: 8px; }
  .bill-invoice-row { display: flex; flex-wrap: wrap; gap: 7px 12px; padding: 14px 0; margin: 0; border-radius: 0; }
  .bill-date { order: 1; font-size: 10px; }
  .bill-paid { order: 2; margin-left: auto; }
  .bill-item { order: 3; flex-basis: 100%; }
  .bill-amount { order: 4; font-size: 12.5px; }
  .bill-links { order: 5; margin-left: auto; }
  .bill-note { margin-top: 12px; padding-top: 12px; font-size: 10px; }
  .bill-dark { padding: 22px 20px 18px; border-radius: 18px; }
  .bill-dark h3 { font-size: 18px; }
  .bill-query-info { padding: 22px 20px 18px; }
  .bill-aside-note { font-size: 10px; }
  .bill-foot { flex-direction: column; align-items: flex-start; gap: 8px; padding-bottom: 26px; }
  .bill-foot-mono { font-size: 8.5px; }
}
</style>
