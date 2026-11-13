<template>
  <div class="pb-2 space-y-5">
    <div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-ink">自建卡密</h1>
        <p class="text-sm text-muted mt-2">
          本站生成、走 Pro20x 通道的兑换码。兑换页按码来源自动分流，与卡台 CDK 互不影响。
        </p>
      </div>
      <div class="flex flex-wrap gap-2">
        <el-button @click="openGen">生成卡密</el-button>
        <el-button :loading="loading" @click="load">刷新</el-button>
      </div>
    </div>

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

    <div class="card !py-3 flex flex-wrap items-center gap-3 text-sm">
      <el-input v-model="q" clearable placeholder="搜索卡密 / 批次 / 邮箱" style="width: 260px" />
      <el-select v-model="statusFilter" clearable placeholder="状态" style="width: 140px">
        <el-option label="未使用" value="unused" />
        <el-option label="已使用" value="consumed" />
        <el-option label="已禁用" value="disabled" />
      </el-select>
      <el-button @click="load">筛选</el-button>
      <span class="text-xs text-muted">列表最多显示最近 300 条</span>
    </div>

    <div v-if="error" class="alert alert-error">{{ error }}</div>

    <div class="card overflow-hidden !p-0">
      <el-table :data="displayRows" v-loading="loading" size="small" stripe empty-text="还没有自建卡密">
        <el-table-column label="卡密" min-width="230">
          <template #default="{ row }">
            <span class="mono text-sm">{{ row.code }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="plan" label="套餐" width="100" />
        <el-table-column label="地区" width="120">
          <template #default="{ row }">{{ countryLabel(row.country) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="statusType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="batch_no" label="批次" min-width="120" />
        <el-table-column label="绑定邮箱" min-width="170">
          <template #default="{ row }">
            <span class="mono text-xs">{{ row.bound_email || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建" width="150">
          <template #default="{ row }">{{ row.created_at }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
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

    <el-dialog v-model="genOpen" title="生成自建卡密" width="520px">
      <el-form label-width="110px" size="small">
        <el-form-item label="套餐">
          <el-select v-model="form.plan" style="width: 100%">
            <el-option v-for="p in planOptions" :key="p.value" :label="p.label" :value="p.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="form.count" :min="1" :max="500" />
        </el-form-item>
        <el-form-item label="充值地区">
          <el-select v-model="form.country" style="width: 100%">
            <el-option v-for="c in countryOptions" :key="c" :label="c" :value="c" />
          </el-select>
        </el-form-item>
        <el-form-item label="有效期(天)">
          <el-input-number v-model="form.valid_days" :min="0" :max="3650" />
          <span class="text-xs text-muted ml-2">0 = 长期有效</span>
        </el-form-item>
        <el-form-item label="售价(分)">
          <el-input-number v-model="form.price_minor" :min="0" :max="10000000" />
          <span class="text-xs text-muted ml-2">仅记账</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="genOpen = false">取消</el-button>
        <el-button type="primary" :loading="generating" @click="generate">生成</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="resultOpen" title="生成结果" width="620px">
      <p class="text-sm text-muted mb-2">批次 {{ lastBatch }}，共 {{ lastCodes.length }} 张，可直接复制或稍后在列表中查看。</p>
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

const genOpen = ref(false)
const generating = ref(false)
const resultOpen = ref(false)
const lastCodes = ref<string[]>([])
const lastBatch = ref('')

const planOptions = [
  { value: 'plus', label: 'Plus' },
  { value: 'pro_5x', label: 'Pro 5x' },
  { value: 'pro_20x', label: 'Pro 20x' },
  { value: 'credit250', label: 'Credits 250' },
  { value: 'credit500', label: 'Credits 500' },
  { value: 'credit1000', label: 'Credits 1000' },
]
const countryOptions = ['PH', 'EG', 'JP', 'US', 'CL', 'SG']

const form = ref({ plan: 'plus', count: 1, country: 'PH', valid_days: 90, price_minor: 0 })

const stats = computed(() => {
  const s = { total: rows.value.length, unused: 0, consumed: 0, disabled: 0, consuming: 0 }
  for (const r of rows.value) {
    const st = String(r.status || '')
    if (st === 'consumed') s.consumed++
    else if (st === 'disabled') s.disabled++
    else if (st === 'consuming') s.consuming++
    else s.unused++
  }
  s.total = rows.value.length
  return s
})

const displayRows = computed(() => {
  let list = rows.value
  const kw = q.value.trim().toLowerCase()
  if (statusFilter.value) list = list.filter((r) => String(r.status || '') === statusFilter.value)
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

const countryNames: Record<string, string> = {
  PH: '菲律宾',
  EG: '埃及',
  JP: '日本',
  US: '美区',
  CL: '智利',
  SG: '新加坡',
}
function countryLabel(c: string) {
  if (!c) return '—'
  const name = countryNames[c] || c
  return `${name} (${c})`
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

function openGen() {
  genOpen.value = true
}

async function generate() {
  generating.value = true
  try {
    const r = await authFetch('/api/v1/admin/cdk/local/generate', {
      method: 'POST',
      body: JSON.stringify(form.value),
    })
    const d = await r.json().catch(() => ({}))
    if (!r.ok) {
      dialog.toast(d.error || '生成失败', 'err')
      return
    }
    lastCodes.value = (d.codes || []).map((c: any) => c.code)
    lastBatch.value = d.batch || ''
    genOpen.value = false
    resultOpen.value = true
    await load()
  } catch (e: any) {
    dialog.toast(e?.message || '网络错误', 'err')
  } finally {
    generating.value = false
  }
}

async function setStatus(code: string, status: string) {
  try {
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

onMounted(load)
</script>

<style scoped>
.mono {
  font-variant-numeric: tabular-nums;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}
</style>
