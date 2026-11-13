<template>
  <div class="pb-2 space-y-5">
    <div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-ink">自建卡密</h1>
        <p class="text-sm text-muted mt-2">
          本站自己生成、走 Pro20x 通道的兑换码。与卡台 CDK 分开管理，兑换页按码自动分流。
        </p>
      </div>
      <div class="flex flex-wrap gap-2">
        <el-button :loading="loading" @click="load">刷新</el-button>
        <el-button type="primary" @click="genOpen = true">生成卡密</el-button>
      </div>
    </div>

    <div class="card !py-3 flex flex-wrap items-center gap-4 text-sm">
      <span class="text-muted">共 <b class="mono text-ink">{{ total }}</b> 张</span>
      <span class="text-muted">未使用 <b class="mono text-ink">{{ unused }}</b> 张</span>
      <el-input v-model="batchFilter" clearable placeholder="按批次筛选" style="width: 200px" @clear="load" @keyup.enter="load" />
      <el-select v-model="statusFilter" clearable placeholder="状态" style="width: 140px" @change="load">
        <el-option label="未使用" value="unused" />
        <el-option label="已使用" value="consumed" />
        <el-option label="已禁用" value="disabled" />
      </el-select>
      <el-button @click="load">筛选</el-button>
    </div>

    <div v-if="error" class="alert alert-error">{{ error }}</div>

    <div class="card overflow-hidden !p-0">
      <el-table :data="rows" v-loading="loading" size="small" stripe empty-text="还没有自建卡密">
        <el-table-column label="卡密" min-width="170">
          <template #default="{ row }">
            <span class="mono text-sm">{{ row.code }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="plan" label="套餐" width="100" />
        <el-table-column prop="country" label="地区" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="statusType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="batch_no" label="批次" min-width="130" />
        <el-table-column label="绑定邮箱" min-width="170">
          <template #default="{ row }">
            <span class="mono text-xs">{{ row.bound_email || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="售价" width="90">
          <template #default="{ row }">{{ row.price_minor ? (row.price_minor / 100).toFixed(2) : '—' }}</template>
        </el-table-column>
        <el-table-column label="有效期" width="150">
          <template #default="{ row }">{{ row.expires_at || '长期' }}</template>
        </el-table-column>
        <el-table-column label="创建" width="150">
          <template #default="{ row }">{{ row.created_at }}</template>
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
          <span class="text-xs text-muted ml-2">仅记账，不参与扣款</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="genOpen = false">取消</el-button>
        <el-button type="primary" :loading="generating" @click="generate">生成</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="resultOpen" title="生成结果（只显示这一次）" width="620px">
      <p class="text-sm text-muted mb-2">批次 {{ lastBatch }}，共 {{ lastCodes.length }} 张，请立即复制保存。</p>
      <div class="mono text-sm" style="max-height: 320px; overflow: auto; background: #0f172a; color: #e2e8f0; padding: 12px; border-radius: 8px">
        <div v-for="c in lastCodes" :key="c" style="line-height: 1.8">{{ c }}</div>
      </div>
      <template #footer>
        <el-button @click="copyAll">复制全部</el-button>
        <el-button type="primary" @click="resultOpen = false">我已保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { authFetch } from '../../lib/api'
import { dialog } from '../../lib/dialog'

const rows = ref<any[]>([])
const total = ref(0)
const unused = ref(0)
const loading = ref(false)
const error = ref('')
const statusFilter = ref('')
const batchFilter = ref('')

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

const form = ref({
  plan: 'plus',
  count: 1,
  country: 'PH',
  valid_days: 90,
  price_minor: 0,
})

function statusLabel(s: string) {
  if (s === 'unused' || s === '') return '未使用'
  if (s === 'consumed') return '已使用'
  if (s === 'disabled') return '已禁用'
  return s
}
function statusType(s: string) {
  if (s === 'consumed') return 'info'
  if (s === 'disabled') return 'danger'
  return 'success'
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams()
    if (statusFilter.value) params.set('status', statusFilter.value)
    if (batchFilter.value.trim()) params.set('batch', batchFilter.value.trim())
    params.set('limit', '300')
    const r = await authFetch(`/api/v1/admin/cdk/local/list?${params}`)
    const d = await r.json().catch(() => ({}))
    if (!r.ok) {
      error.value = d.error || '加载失败'
      return
    }
    rows.value = Array.isArray(d.list) ? d.list : []
    total.value = Number(d.total) || 0
    unused.value = Number(d.unused) || 0
  } catch (e: any) {
    error.value = e?.message || '网络错误'
  } finally {
    loading.value = false
  }
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

async function copyAll() {
  try {
    await navigator.clipboard.writeText(lastCodes.value.join('\n'))
    dialog.toast('已复制', 'ok')
  } catch {
    dialog.toast('复制失败，请手动选择', 'err')
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
