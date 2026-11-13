<template>
  <div class="pb-2 space-y-5">
    <!-- 标题 -->
    <div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-ink">自建卡密</h1>
        <p class="text-sm text-muted mt-2">
          本站生成、走 Pro20x 通道的兑换码。兑换页按码来源自动分流，与卡台 CDK 互不影响。
        </p>
      </div>
      <div class="flex flex-wrap gap-2">
        <el-button :loading="loading" @click="load">刷新</el-button>
      </div>
    </div>

    <!-- 统计 -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <div class="card !py-3">
        <div class="text-xs text-muted">总数</div>
        <div class="text-2xl font-bold text-ink mono">{{ stats.total }}</div>
      </div>
      <div class="card !py-3">
        <div class="text-xs text-muted">未使用</div>
        <div class="text-2xl font-bold text-ink mono">{{ stats.unused }}</div>
      </div>
      <div class="card !py-3">
        <div class="text-xs text-muted">已使用</div>
        <div class="text-2xl font-bold text-ink mono">{{ stats.consumed }}</div>
      </div>
      <div class="card !py-3">
        <div class="text-xs text-muted">已禁用 / 兑换中</div>
        <div class="text-2xl font-bold text-ink mono">{{ stats.disabled + stats.consuming }}</div>
      </div>
    </div>

    <!-- 生成面板 -->
    <div class="card !p-0 overflow-hidden">
      <div
        class="flex items-center justify-between px-4 py-3 cursor-pointer select-none"
        @click="genOpen = !genOpen"
      >
        <div class="text-base font-semibold text-ink">
          生成卡密
          <span class="text-xs text-muted ml-2">当前选择：{{ planLabel(form.plan) }} × {{ form.count }}</span>
        </div>
        <span class="text-sm text-muted">{{ genOpen ? '收起' : '展开' }}</span>
      </div>

      <div v-show="genOpen" class="border-t px-4 py-4 space-y-4">
        <!-- 套餐卡片 -->
        <div class="grid grid-cols-2 lg:grid-cols-3 gap-3">
          <div
            v-for="p in planOptions"
            :key="p.value"
            class="rounded-lg border px-3 py-3 cursor-pointer transition"
            :class="form.plan === p.value ? 'border-2 border-ink bg-[var(--surface-2,#f8fafc)]' : 'border'"
            @click="form.plan = p.value"
          >
            <div class="flex items-center justify-between">
              <span class="text-sm font-semibold text-ink">{{ p.label }}</span>
              <span class="text-xs text-muted">{{ p.hint }}</span>
            </div>
            <div class="mt-1 text-xs text-muted">Pro20x · 1 积分</div>
          </div>
        </div>

        <!-- 数量 + 地区 -->
        <div class="flex flex-wrap items-center gap-3 text-sm">
          <span class="text-muted">数量</span>
          <el-input-number v-model="form.count" :min="1" :max="500" style="width: 120px" />
          <el-button v-for="n in [1, 10, 50, 100, 200]" :key="n" size="small" @click="form.count = n">{{ n }}</el-button>
          <span class="text-muted ml-3">充值地区</span>
          <el-select v-model="form.country" style="width: 180px">
            <el-option v-for="c in countryOptions" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
          <span class="text-muted ml-3">有效期(天)</span>
          <el-input-number v-model="form.valid_days" :min="0" :max="3650" style="width: 120px" />
          <span class="text-xs text-muted">0 = 长期</span>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <el-button type="primary" :loading="generating" @click="generate">
            生成 {{ form.count }} 张 {{ planLabel(form.plan) }}
          </el-button>
          <span class="text-xs text-muted">生成后完整卡密只显示一次，请立即复制或导出</span>
        </div>
      </div>
    </div>

    <!-- 列表 -->
    <div class="card !p-0 overflow-hidden">
      <div class="flex flex-wrap items-center justify-between gap-3 px-4 py-3 border-b">
        <div class="text-base font-semibold text-ink">卡密列表</div>
        <div class="text-xs text-muted">共 {{ filteredRows.length }} 条（最多显示最近 300 条）</div>
      </div>

      <div class="flex flex-wrap items-center gap-2 px-4 py-3 text-sm">
        <el-input v-model="q" clearable placeholder="搜索卡密 / 批次 / 邮箱" style="width: 260px" />
        <el-select v-model="statusFilter" clearable placeholder="状态" style="width: 140px">
          <el-option label="未使用" value="unused" />
          <el-option label="已使用" value="consumed" />
          <el-option label="已禁用" value="disabled" />
          <el-option label="兑换中" value="consuming" />
        </el-select>
        <el-button @click="load">查询</el-button>
        <el-button @click="copyAllCodes">复制全部未使用</el-button>
      </div>

      <div v-if="error" class="px-4 pb-3">
        <div class="alert alert-error">{{ error }}</div>
      </div>

      <el-table :data="filteredRows" v-loading="loading" size="small" stripe empty-text="还没有自建卡密">
        <el-table-column label="卡密" min-width="250">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <span class="mono text-sm">{{ row.code }}</span>
              <el-button link type="primary" size="small" @click="copyOne(row.code)">复制</el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="套餐" width="110">
          <template #default="{ row }">{{ planLabel(row.plan) }}</template>
        </el-table-column>
        <el-table-column label="区域" width="140">
          <template #default="{ row }">{{ countryLabel(row.country) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="statusType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="batch_no" label="批次" min-width="130" />
        <el-table-column label="绑定邮箱" min-width="180">
          <template #default="{ row }">
            <span class="mono text-xs">{{ row.bound_email || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="备注" min-width="180">
          <template #default="{ row }">
            <span class="text-xs cursor-pointer" :class="row.note ? 'text-ink' : 'text-muted'" @click="openNote(row)">
              {{ row.note || '添加备注…' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="创建" width="150">
          <template #default="{ row }">{{ row.created_at }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="copyOne(row.code)">复制</el-button>
            <el-button
              v-if="row.status === 'unused'"
              link
              type="danger"
              size="small"
              @click="setStatus(row.code, 'disabled')"
            >禁用</el-button>
            <el-button
              v-else-if="row.status === 'disabled'"
              link
              type="primary"
              size="small"
              @click="setStatus(row.code, 'unused')"
            >启用</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 备注编辑 -->
    <el-dialog v-model="noteOpen" :title="`备注 ${noteTarget?.batch_no || ''}`" width="520px">
      <p class="text-sm text-muted mb-3">为该 CDK 填写备注（留空并确定可清空）</p>
      <el-input
        v-model="noteDraft"
        type="textarea"
        :rows="2"
        maxlength="200"
        show-word-limit
        placeholder="最多 200 字，例如批次/客户/用途"
      />
      <template #footer>
        <el-button @click="noteOpen = false">取消</el-button>
        <el-button type="primary" @click="saveNote">保存</el-button>
      </template>
    </el-dialog>

    <!-- 生成结果 -->
    <el-dialog v-model="resultOpen" title="生成结果" width="620px">
      <p class="text-sm text-muted mb-2">批次 {{ lastBatch }}，共 {{ lastCodes.length }} 张，可直接复制或稍后在列表查看。</p>
      <div class="mono text-sm" style="max-height: 320px; overflow: auto; background: #0f172a; color: #e2e8f0; padding: 12px; border-radius: 8px">
        <div v-for="c in lastCodes" :key="c" style="line-height: 1.8">{{ c }}</div>
      </div>
      <template #footer>
        <el-button @click="copyAll">复制全部</el-button>
        <el-button type="primary" @click="resultOpen = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { authFetch } from '../../lib/api'
import { dialog } from '../../lib/dialog'

const rows = ref<any[]>([])
const loading = ref(false)
const error = ref('')
const q = ref('')
const statusFilter = ref('')

const genOpen = ref(true)
const generating = ref(false)
const resultOpen = ref(false)
const lastCodes = ref<string[]>([])
const lastBatch = ref('')
const noteDraft = ref('')
const noteOpen = ref(false)
const noteTarget = ref<any>(null)

function openNote(row: any) {
  noteTarget.value = row
  noteDraft.value = row.note || ''
  noteOpen.value = true
}

const planOptions = [
  { value: 'plus', label: 'Plus', hint: '直充' },
  { value: 'pro_5x', label: 'Pro 5x', hint: '直充' },
  { value: 'pro_20x', label: 'Pro 20x', hint: '直付' },
  { value: 'credit250', label: 'Codex 点数 250', hint: '点数' },
  { value: 'credit500', label: 'Codex 点数 500', hint: '点数' },
  { value: 'credit1000', label: 'Codex 点数 1000', hint: '点数' },
]

const countryOptions = [
  { value: 'PH', label: '菲律宾 (PH)' },
  { value: 'EG', label: '埃及 (EG)' },
  { value: 'JP', label: '日本 (JP)' },
  { value: 'US', label: '美区 (US)' },
  { value: 'CL', label: '智利 (CL)' },
  { value: 'SG', label: '新加坡 (SG)' },
]

const form = ref({ plan: 'plus', count: 1, country: 'PH', valid_days: 0 })

const planNames: Record<string, string> = {
  plus: 'Plus',
  pro_5x: 'Pro 5x',
  pro_20x: 'Pro 20x',
  credit250: 'Codex 点数 250',
  credit500: 'Codex 点数 500',
  credit1000: 'Codex 点数 1000',
}
const countryNames: Record<string, string> = {
  PH: '菲律宾',
  EG: '埃及',
  JP: '日本',
  US: '美区',
  CL: '智利',
  SG: '新加坡',
}

function planLabel(p: string) {
  return planNames[p] || p || '—'
}
function countryLabel(c: string) {
  if (!c) return '—'
  return `${countryNames[c] || c} (${c})`
}

const stats = computed(() => {
  const s = { total: 0, unused: 0, consumed: 0, disabled: 0, consuming: 0 }
  for (const r of rows.value) {
    s.total++
    const st = String(r.status || '')
    if (st === 'consumed') s.consumed++
    else if (st === 'disabled') s.disabled++
    else if (st === 'consuming') s.consuming++
    else s.unused++
  }
  return s
})

const filteredRows = computed(() => {
  let list = rows.value
  if (statusFilter.value) {
    list = list.filter((r) => (String(r.status || '') || 'unused') === statusFilter.value)
  }
  const kw = q.value.trim().toLowerCase()
  if (kw) {
    list = list.filter((r) =>
      String(r.code || '').toLowerCase().includes(kw) ||
      String(r.batch_no || '').toLowerCase().includes(kw) ||
      String(r.bound_email || '').toLowerCase().includes(kw),
    )
  }
  return list
})

function statusLabel(s: string) {
  if (s === 'unused' || s === '') return '未使用'
  if (s === 'consumed') return '已使用'
  if (s === 'disabled') return '已禁用'
  if (s === 'consuming') return '兑换中'
  return s
}
function statusType(s: string) {
  if (s === 'consumed') return 'info'
  if (s === 'disabled') return 'danger'
  if (s === 'consuming') return 'warning'
  return 'success'
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const r = await authFetch('/api/v1/admin/cdk/local/list?limit=300')
    const d = await r.json().catch(() => ({}))
    if (!r.ok) {
      error.value = d.error || '加载失败'
      return
    }
    rows.value = Array.isArray(d.list) ? d.list : []
  } catch (e: any) {
    error.value = e?.message || '网络错误'
  } finally {
    loading.value = false
  }
}

async function generate() {
  generating.value = true
  try {
    const body = {
      plan: form.value.plan,
      count: form.value.count,
      country: form.value.country,
      valid_days: form.value.valid_days,
      price_minor: 0,
    }
    const r = await authFetch('/api/v1/admin/cdk/local/generate', {
      method: 'POST',
      body: JSON.stringify(body),
    })
    const d = await r.json().catch(() => ({}))
    if (!r.ok) {
      dialog.toast(d.error || '生成失败', 'err')
      return
    }
    lastCodes.value = (d.codes || []).map((c: any) => c.code)
    lastBatch.value = d.batch || ''
    resultOpen.value = true
    await load()
  } catch (e: any) {
    dialog.toast(e?.message || '网络错误', 'err')
  } finally {
    generating.value = false
  }
}

async function saveNote() {
  const row = noteTarget.value
  if (!row) return
  try {
    const r = await authFetch('/api/v1/admin/cdk/local/note', {
      method: 'PUT',
      body: JSON.stringify({ code: row.code, note: noteDraft.value }),
    })
    const d = await r.json().catch(() => ({}))
    if (!r.ok) {
      dialog.toast(d.error || '保存失败', 'err')
      return
    }
    row.note = noteDraft.value
    noteOpen.value = false
    dialog.toast('备注已保存', 'ok')
  } catch (e: any) {
    dialog.toast(e?.message || '网络错误', 'err')
  }
}

async function setStatus(code: string, status: string) {  try {
    const r = await authFetch('/api/v1/admin/cdk/local/status', {
      method: 'PUT',
      body: JSON.stringify({ code, status }),
    })
    const d = await r.json().catch(() => ({}))
    if (!r.ok) {
      dialog.toast(d.error || '操作失败', 'err')
      return
    }
    dialog.toast(status === 'disabled' ? '已禁用' : '已启用', 'ok')
    await load()
  } catch (e: any) {
    dialog.toast(e?.message || '网络错误', 'err')
  }
}

async function copyOne(code: string) {
  try {
    await navigator.clipboard.writeText(code)
    dialog.toast('已复制', 'ok')
  } catch {
    dialog.toast('复制失败', 'err')
  }
}

async function copyAll() {
  try {
    await navigator.clipboard.writeText(lastCodes.value.join('\n'))
    dialog.toast('已复制', 'ok')
  } catch {
    dialog.toast('复制失败', 'err')
  }
}

async function copyAllCodes() {
  const list = filteredRows.value.filter((r) => (r.status || 'unused') === 'unused').map((r) => r.code)
  if (list.length === 0) {
    dialog.toast('没有未使用的卡密', 'warn')
    return
  }
  try {
    await navigator.clipboard.writeText(list.join('\n'))
    dialog.toast(`已复制 ${list.length} 张未使用卡密`, 'ok')
  } catch {
    dialog.toast('复制失败', 'err')
  }
}

onMounted(load)
</script>

<style scoped>
.mono {
  font-variant-numeric: tabular-nums;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}
</style>
