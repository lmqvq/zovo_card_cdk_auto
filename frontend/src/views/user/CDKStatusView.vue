<template>
  <div class="u-page lookup-page">
    <UserPageHeader hide-hero :title="t('cdkLookup.title')" />
    <div class="u-main lookup-wrap">
      <div class="lookup-subnav">
        <div class="lookup-crumb">{{ t('hero.redeemEyebrow') }} <i>/</i> <b>{{ t('cdkLookup.title') }}</b> <i>/</i> LOOKUP</div>
        <RedeemModeTabs compact />
      </div>

      <section class="lookup-hero">
        <div class="lookup-hero-copy">
          <div class="lookup-eyebrow">LOOKUP <em>· {{ t('cdkLookup.eyebrow') }}</em></div>
          <h1>{{ t('cdkLookup.title') }}</h1>
          <p class="lookup-sub">{{ t('cdkLookup.subtitle') }}</p>
          <ul class="lookup-badges">
            <li>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 6 9 17l-5-5" /></svg>
              {{ t('cdkLookup.badgeStatus') }}
            </li>
            <li>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" /><path d="m9 12 2 2 4-4" /></svg>
              {{ t('cdkLookup.badgeNoSession') }}
            </li>
            <li>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M4 6h16" /><path d="M4 12h16" /><path d="M4 18h10" /></svg>
              {{ t('cdkLookup.badgeBatch') }}
            </li>
          </ul>
        </div>
        <div class="lookup-aside">
          <span>{{ t('cdkLookup.asideNoCode') }}</span>
          <router-link to="/recharge">
            {{ t('cdkLookup.goRedeem') }}
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M7 17 17 7" /><path d="M8 7h9v9" /></svg>
          </router-link>
        </div>
      </section>

      <div class="lookup-grid">
        <div class="lookup-col">
          <section class="card lookup-card">
            <div class="lookup-modes">
              <button type="button" class="lookup-mode" :class="{ active: mode === 'single' }" @click="setMode('single')">
                {{ t('cdkLookup.modeSingle') }}
              </button>
              <button type="button" class="lookup-mode" :class="{ active: mode === 'batch' }" @click="setMode('batch')">
                {{ t('cdkLookup.modeBatch') }}
              </button>
            </div>

            <template v-if="mode === 'single'">
              <div class="lookup-head">
                <span class="lookup-micro">LOOKUP · {{ t('cdkLookup.modeSingle') }}</span>
                <h2>{{ t('cdkLookup.queryTitle') }}</h2>
              </div>

              <div class="lookup-hint">
                <span class="lookup-hint-icon">i</span>
                <p>{{ t('cdkLookup.queryHint') }}</p>
              </div>

              <div class="lookup-field">
                <div class="lookup-label"><span>{{ t('cdkLookup.cdkLabel') }}</span><span>{{ t('cdkLookup.required') }}</span></div>
                <div class="lookup-input-shell">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 9a3 3 0 0 0 0 6v3h18v-3a3 3 0 0 1 0-6V6H3z" /><path d="M13 6v2" /><path d="M13 11v2" /><path d="M13 16v2" /></svg>
                  <input
                    v-model="cdkCode"
                    class="lookup-input"
                    type="text"
                    :placeholder="t('cdkLookup.cdkPlaceholder')"
                    autocomplete="off"
                    spellcheck="false"
                    :aria-label="t('cdkLookup.cdkLabel')"
                    @keyup.enter="query"
                  />
                </div>
              </div>

              <button class="lookup-cta" type="button" :disabled="!cdkCode.trim() || querying" @click="query">
                <span v-if="!querying">{{ t('cdkLookup.queryBtn') }}</span>
                <span v-else class="lookup-loading"><span class="spinner"></span>{{ t('common.querying') }}</span>
                <svg v-if="!querying" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12h14" /><path d="m13 6 6 6-6 6" /></svg>
              </button>
              <p class="lookup-note">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
                {{ t('cdkLookup.readOnlyNote') }}
              </p>
            </template>

            <template v-else>
              <div class="lookup-head">
                <span class="lookup-micro">LOOKUP · {{ t('cdkLookup.modeBatch') }}</span>
                <h2>{{ t('cdkLookup.modeBatch') }}</h2>
              </div>

              <div class="lookup-hint">
                <span class="lookup-hint-icon">i</span>
                <p>{{ t('cdkLookup.batchHint', { max: LOOKUP_BATCH_MAX }) }}</p>
              </div>

              <div class="lookup-field">
                <div class="lookup-label"><span>{{ t('cdkLookup.batchLabel') }}</span><span>{{ t('cdkLookup.required') }}</span></div>
                <textarea
                  v-model="batchText"
                  rows="8"
                  class="lookup-textarea"
                  :placeholder="t('cdkLookup.batchPlaceholder')"
                />
                <p class="lookup-meta">
                  {{ t('cdkLookup.batchRecognized', { n: parsedCodes.length }) }}
                  <span v-if="parsedCodes.length > LOOKUP_BATCH_MAX" class="is-warn">
                    · {{ t('cdkLookup.batchTruncated', { max: LOOKUP_BATCH_MAX }) }}
                  </span>
                </p>
              </div>

              <button class="lookup-cta" type="button" :disabled="parsedCodes.length === 0 || querying" @click="queryBatch">
                <span v-if="!querying">{{ t('cdkLookup.batchQueryBtn') }}</span>
                <span v-else class="lookup-loading"><span class="spinner"></span>{{ t('common.querying') }}</span>
                <svg v-if="!querying" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12h14" /><path d="m13 6 6 6-6 6" /></svg>
              </button>
              <p class="lookup-note">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
                {{ t('cdkLookup.readOnlyNote') }}
              </p>
            </template>

            <p v-if="error" class="lookup-error">{{ error }}</p>
          </section>

          <section v-if="mode === 'single' && result" class="card lookup-result">
            <div class="lookup-result-head">
              <div>
                <span class="lookup-micro">LOOKUP RESULT</span>
                <h2>{{ t('cdkLookup.resultTitle') }}</h2>
              </div>
              <span class="lookup-status" :class="statusTone(result.status)">
                <i></i>{{ statusText(result.status) }}
              </span>
            </div>

            <div class="lookup-receipt">
              <div class="lookup-row"><span>{{ t('cdkLookup.cdkCode') }}</span><b class="mono">{{ result.cdk_code }}</b></div>
              <div class="lookup-row">
                <span>{{ t('cdkLookup.useStatus') }}</span>
                <b :class="statusTone(result.status)">{{ statusText(result.status) }}</b>
              </div>
              <div class="lookup-row">
                <span>{{ t('cdkLookup.rechargeEmail') }}</span>
                <b class="mono">{{ result.account_email || t('cdkLookup.emailEmpty') }}</b>
              </div>
              <div v-if="result.plan" class="lookup-row"><span>{{ t('cdkLookup.plan') }}</span><b>{{ result.plan }}</b></div>
              <div v-if="result.used_at" class="lookup-row"><span>{{ t('cdkLookup.usedAt') }}</span><b class="mono">{{ result.used_at }}</b></div>
            </div>

            <div v-if="result.status === 'failed'" class="lookup-message is-err">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 8v5" /><path d="M12 16h.01" /></svg>
              {{ result.message || t('cdkLookup.msgFailed') }}
            </div>
            <div v-else-if="result.used" class="lookup-message is-good">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m5 12 4.5 4.5L19 7" /></svg>
              {{ result.message || t('cdkLookup.msgUsed') }}
            </div>
            <div v-else-if="result.status === 'processing'" class="lookup-message is-info">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></svg>
              {{ result.message || t('cdkLookup.msgProcessing') }}
            </div>
            <div v-else-if="result.status === 'disabled' || result.status === 'expired'" class="lookup-message is-warn">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 8v5" /><path d="M12 16h.01" /></svg>
              {{ result.message || statusText(result.status) }}
            </div>
            <div v-else class="lookup-message is-info">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 8v.5" /><path d="M12 11v5" /></svg>
              {{ result.message || t('cdkLookup.msgUnused') }}
            </div>

            <div class="lookup-actions">
              <button v-if="result.can_resubmit" class="lookup-primary" type="button" @click="goRedeem(result.cdk_code)">
                {{ t('cdkLookup.resubmit') }}
              </button>
              <button v-else-if="result.status === 'unused'" class="lookup-primary" type="button" @click="goRedeem(result.cdk_code)">
                {{ t('cdkLookup.goRedeem') }}
              </button>
              <button class="lookup-secondary" type="button" @click="reset">{{ t('cdkLookup.queryOther') }}</button>
            </div>
          </section>

          <section v-if="mode === 'batch' && batchResults.length" class="card lookup-result lookup-batch-results">
            <div class="lookup-result-head">
              <div>
                <span class="lookup-micro">LOOKUP RESULT</span>
                <h2>{{ t('cdkLookup.resultTitle') }} <span class="lookup-count">{{ batchResults.length }}</span></h2>
              </div>
              <div class="lookup-tools">
                <button type="button" class="lookup-tool" @click="copyUsedEmails">{{ t('cdkLookup.copyEmails') }}</button>
                <button type="button" class="lookup-tool" @click="exportCsv">{{ t('cdkLookup.exportCsv') }}</button>
              </div>
            </div>

            <div class="lookup-summary">
              <span class="lookup-chip is-good">{{ t('cdkLookup.summaryUsed', { n: batchSummary.used }) }}</span>
              <span class="lookup-chip is-err">{{ t('cdkLookup.summaryFailed', { n: batchSummary.failed }) }}</span>
              <span class="lookup-chip is-muted">{{ t('cdkLookup.summaryUnused', { n: batchSummary.unused }) }}</span>
              <span class="lookup-chip is-info">{{ t('cdkLookup.summaryProcessing', { n: batchSummary.processing }) }}</span>
              <span class="lookup-chip is-warn">{{ t('cdkLookup.summaryUnknown', { n: batchSummary.unknown }) }}</span>
            </div>

            <div class="lookup-table-wrap">
              <table class="lookup-table">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>{{ t('cdkLookup.cdkCode') }}</th>
                    <th>{{ t('cdkLookup.useStatus') }}</th>
                    <th>{{ t('cdkLookup.rechargeEmail') }}</th>
                    <th>{{ t('cdkLookup.plan') }}</th>
                    <th>{{ t('cdkLookup.usedAt') }}</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(row, i) in batchResults" :key="row.cdk_code + i">
                    <td class="lookup-index">{{ i + 1 }}</td>
                    <td class="mono">{{ row.cdk_code }}</td>
                    <td><span class="lookup-status-text" :class="statusTone(row.status)">{{ statusText(row.status) }}</span></td>
                    <td class="mono">{{ row.account_email || t('cdkLookup.emailEmpty') }}</td>
                    <td>{{ row.plan || t('cdkLookup.emailEmpty') }}</td>
                    <td class="lookup-time">{{ row.used_at || t('cdkLookup.emailEmpty') }}</td>
                    <td class="lookup-table-action">
                      <button v-if="row.can_resubmit" type="button" class="lookup-row-action is-primary" @click="goRedeem(row.cdk_code)">
                        {{ t('cdkLookup.resubmit') }}
                      </button>
                      <button v-else-if="row.status === 'unused'" type="button" class="lookup-row-action" @click="goRedeem(row.cdk_code)">
                        {{ t('cdkLookup.goRedeem') }}
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>
        </div>

        <aside class="lookup-col lookup-side">
          <section class="lookup-scope">
            <div class="lookup-scope-top"><span>LOOKUP SCOPE</span><span>READ ONLY</span></div>
            <h3>{{ t('cdkLookup.scopeTitle') }}</h3>
            <p class="lookup-scope-sub">{{ t('cdkLookup.scopeSub') }}</p>
            <ul class="lookup-scope-list">
              <li>
                <span class="lookup-scope-check"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m5 12 4.5 4.5L19 7" /></svg></span>
                <div><b>{{ t('cdkLookup.useStatus') }}</b><p>{{ t('cdkLookup.scopeStatusDesc') }}</p></div>
              </li>
              <li>
                <span class="lookup-scope-check"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m5 12 4.5 4.5L19 7" /></svg></span>
                <div><b>{{ t('cdkLookup.rechargeEmail') }}</b><p>{{ t('cdkLookup.scopeEmailDesc') }}</p></div>
              </li>
              <li>
                <span class="lookup-scope-check"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m5 12 4.5 4.5L19 7" /></svg></span>
                <div><b>{{ t('cdkLookup.plan') }} / {{ t('cdkLookup.usedAt') }}</b><p>{{ t('cdkLookup.scopePlanDesc') }}</p></div>
              </li>
              <li class="is-no">
                <span class="lookup-scope-cross"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 6l12 12" /><path d="M18 6 6 18" /></svg></span>
                <div><b>Session / Token</b><p>{{ t('cdkLookup.scopeNoSessionDesc') }}</p></div>
              </li>
            </ul>
            <div class="lookup-scope-foot">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
              {{ t('cdkLookup.scopeFoot') }}
            </div>
          </section>

          <section class="card lookup-batch-info">
            <h3>{{ t('cdkLookup.modeBatch') }}</h3>
            <p class="lookup-batch-desc">{{ t('cdkLookup.batchCardDesc', { max: LOOKUP_BATCH_MAX }) }}</p>
            <ol class="lookup-batch-list">
              <li><span>01</span><div><b>{{ t('cdkLookup.batchRow1') }}</b><p>{{ t('cdkLookup.batchRow1Desc') }}</p></div></li>
              <li><span>02</span><div><b>{{ t('cdkLookup.batchRow2') }}</b><p>{{ t('cdkLookup.batchRow2Desc') }}</p></div></li>
              <li><span>03</span><div><b>{{ t('cdkLookup.batchRow3') }}</b><p>{{ t('cdkLookup.batchRow3Desc') }}</p></div></li>
            </ol>
            <p class="lookup-batch-note">{{ t('cdkLookup.batchNote') }}</p>
          </section>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import UserPageHeader from '../../components/UserPageHeader.vue'
import RedeemModeTabs from '../../components/RedeemModeTabs.vue'
import { parseCdks } from '../../lib/batch-session'
import { copyToClipboard } from '../../lib/clipboard'
import { toast } from '../../lib/ui'

const { t } = useI18n({ useScope: 'global' })
const router = useRouter()

const LOOKUP_BATCH_MAX = 100

interface CDKStatusResult {
  cdk_code: string
  status: string
  used: boolean
  can_resubmit?: boolean
  account_email?: string
  plan?: string
  used_at?: string
  notes?: string
  message?: string
}

const mode = ref<'single' | 'batch'>('single')
const cdkCode = ref('')
const batchText = ref('')
const querying = ref(false)
const error = ref('')
const result = ref<CDKStatusResult | null>(null)
const batchResults = ref<CDKStatusResult[]>([])

const parsedCodes = computed(() => parseCdks(batchText.value))

const batchSummary = computed(() => {
  const s = { used: 0, failed: 0, unused: 0, processing: 0, unknown: 0 }
  for (const row of batchResults.value) {
    if (row.status === 'used') s.used++
    else if (row.status === 'failed') s.failed++
    else if (row.status === 'unused') s.unused++
    else if (row.status === 'processing') s.processing++
    else s.unknown++
  }
  return s
})

function statusText(status: string) {
  const key = `cdkLookup.status.${status}`
  const label = t(key)
  return label === key ? status : label
}

function statusTone(status: string) {
  if (status === 'used') return 'is-good'
  if (status === 'failed') return 'is-err'
  if (status === 'processing') return 'is-info'
  if (status === 'disabled' || status === 'expired' || status === 'unknown') return 'is-warn'
  return 'is-muted'
}

function goRedeem(code: string) {
  router.push({ path: '/recharge', query: { cdk: code } })
}

async function query() {
  const code = cdkCode.value.trim()
  if (!code) return
  querying.value = true
  error.value = ''
  result.value = null
  try {
    const r = await fetch(`/api/v1/lookup/cdk?code=${encodeURIComponent(code)}`)
    const data = await r.json().catch(() => ({}))
    if (!r.ok) {
      error.value = data.error || data.message || t('cdkLookup.errNotFound')
      return
    }
    result.value = data as CDKStatusResult
  } catch {
    error.value = t('cdkLookup.errNetwork')
  } finally {
    querying.value = false
  }
}

async function queryBatch() {
  const codes = parsedCodes.value.slice(0, LOOKUP_BATCH_MAX)
  if (!codes.length) {
    error.value = t('cdkLookup.batchEmpty')
    return
  }
  querying.value = true
  error.value = ''
  batchResults.value = []
  try {
    const r = await fetch('/api/v1/lookup/cdk/batch', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ codes }),
    })
    const data = await r.json().catch(() => ({}))
    if (!r.ok) {
      error.value = data.error || data.message || t('cdkLookup.errNetwork')
      return
    }
    batchResults.value = Array.isArray(data.results) ? data.results : []
  } catch {
    error.value = t('cdkLookup.errNetwork')
  } finally {
    querying.value = false
  }
}

async function copyUsedEmails() {
  const emails = Array.from(
    new Set(
      batchResults.value
        .map((row) => String(row.account_email || '').trim())
        .filter(Boolean),
    ),
  )
  if (!emails.length) {
    toast(t('cdkLookup.emailEmpty'), 'warn')
    return
  }
  const ok = await copyToClipboard(emails.join('\n'))
  toast(ok ? t('cdkLookup.copied') : t('cdkLookup.copyFail'), ok ? 'ok' : 'err')
}

function exportCsv() {
  const header = ['cdk_code', 'status', 'can_resubmit', 'account_email', 'plan', 'used_at', 'message']
  const lines = [header.join(',')]
  for (const row of batchResults.value) {
    const cells = [
      row.cdk_code,
      row.status,
      row.can_resubmit ? '1' : '0',
      row.account_email || '',
      row.plan || '',
      row.used_at || '',
      row.message || '',
    ].map((v) => `"${String(v).replace(/"/g, '""')}"`)
    lines.push(cells.join(','))
  }
  const blob = new Blob([`\uFEFF${lines.join('\n')}`], { type: 'text/csv;charset=utf-8' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `cdk-lookup-${new Date().toISOString().slice(0, 10)}.csv`
  a.click()
  URL.revokeObjectURL(a.href)
}

function setMode(next: 'single' | 'batch') {
  mode.value = next
  error.value = ''
}

function reset() {
  result.value = null
  error.value = ''
  cdkCode.value = ''
}
</script>

<style scoped>
.lookup-wrap {
  width: 100%;
  max-width: 1180px;
  margin: 0 auto;
  padding: 0 24px 48px;
}
.lookup-subnav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 16px 0;
  border-bottom: 1px solid var(--brd);
}
.lookup-crumb {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  white-space: nowrap;
}
.lookup-crumb b { color: var(--ink); font-weight: 600; }
.lookup-crumb i { color: var(--ink-4); font-style: normal; }

.lookup-hero {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 28px;
  padding: 24px 0 26px;
}
.lookup-hero-copy { min-width: 0; }
.lookup-eyebrow {
  display: flex;
  align-items: center;
  gap: 9px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
}
.lookup-eyebrow::before { content: ''; width: 18px; height: 1px; background: var(--ink-3); }
.lookup-eyebrow em { color: var(--ink); font-style: normal; letter-spacing: 0.16em; }
.lookup-hero h1 {
  margin: 12px 0 0;
  color: var(--ink);
  font-size: 54px;
  font-weight: 800;
  letter-spacing: -0.05em;
  line-height: 1.08;
}
.lookup-sub {
  max-width: 620px;
  margin: 14px 0 0;
  color: var(--ink-2);
  font-size: 13.5px;
  line-height: 1.75;
}
.lookup-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 18px 0 0;
  padding: 0;
  list-style: none;
}
.lookup-badges li {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 10px;
  border: 1px solid var(--brd);
  border-radius: 999px;
  background: var(--surface);
  color: var(--ink-2);
  font-size: 10.5px;
}
.lookup-badges svg { width: 11px; height: 11px; color: var(--ink-3); }
.lookup-aside {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  flex: none;
  padding-bottom: 3px;
  text-align: right;
}
.lookup-aside span { color: var(--ink-3); font-size: 10.5px; }
.lookup-aside a {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 5px;
  color: var(--ink);
  font-size: 11.5px;
  font-weight: 600;
  text-decoration: none;
}
.lookup-aside a svg { width: 11px; height: 11px; }

.lookup-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  align-items: start;
  gap: 24px;
}
.lookup-col { display: flex; flex-direction: column; gap: 16px; min-width: 0; }
.lookup-side { min-width: 0; }

.lookup-card { padding: 26px 28px 22px; }
.lookup-modes {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3px;
  padding: 4px;
  border: 1px solid var(--brd);
  border-radius: 999px;
  background: var(--soft);
}
.lookup-mode {
  height: 34px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: var(--ink-2);
  font-family: var(--font-sans);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}
.lookup-mode.active { background: var(--ink); color: var(--bg); }
.lookup-head { margin-top: 18px; }
.lookup-micro {
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 8.5px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}
.lookup-head h2 {
  margin: 7px 0 0;
  color: var(--ink);
  font-size: 23px;
  font-weight: 800;
  letter-spacing: -0.025em;
}
.lookup-hint {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-top: 14px;
  padding: 12px 14px;
  border: 1px solid var(--brd);
  border-radius: 12px;
  background: var(--surface-2);
}
.lookup-hint-icon {
  display: grid;
  place-items: center;
  flex: none;
  width: 22px;
  height: 22px;
  border: 1px solid var(--brd-2);
  border-radius: 50%;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 11px;
}
.lookup-hint p { margin: 0; color: var(--ink-2); font-size: 11.5px; line-height: 1.65; }
.lookup-field { margin-top: 18px; }
.lookup-label {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
  color: var(--ink-2);
  font-size: 11px;
}
.lookup-label span:last-child { color: var(--ink-3); font-size: 9.5px; }
.lookup-input-shell {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 54px;
  padding: 0 14px;
  border: 1px solid var(--brd);
  border-radius: 12px;
  background: var(--surface-2);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
.lookup-input-shell:focus-within { border-color: var(--ink-3); box-shadow: 0 0 0 3px var(--primary-soft); }
.lookup-input-shell > svg { width: 17px; height: 17px; flex: none; color: var(--ink-3); }
.lookup-input {
  flex: 1;
  min-width: 0;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 13px;
}
.lookup-input::placeholder { color: var(--ink-4); }
.lookup-textarea {
  display: block;
  width: 100%;
  min-height: 160px;
  padding: 12px 14px;
  resize: vertical;
  border: 1px solid var(--brd);
  border-radius: 12px;
  outline: 0;
  background: var(--surface-2);
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 12px;
  line-height: 1.6;
}
.lookup-textarea:focus { border-color: var(--ink-3); box-shadow: 0 0 0 3px var(--primary-soft); }
.lookup-meta { margin: 7px 0 0; color: var(--ink-3); font-size: 10px; }
.lookup-meta .is-warn { color: var(--warn); }
.lookup-cta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  height: 48px;
  margin-top: 20px;
  border: 0;
  border-radius: 999px;
  background: var(--ink);
  color: var(--bg);
  font-family: var(--font-sans);
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 8px 18px -10px rgba(0, 0, 0, 0.55);
}
.lookup-cta svg { width: 13px; height: 13px; }
.lookup-cta:disabled { opacity: 0.45; cursor: not-allowed; box-shadow: none; }
.lookup-loading { display: inline-flex; align-items: center; gap: 8px; }
.lookup-note {
  display: flex;
  align-items: flex-start;
  gap: 7px;
  margin: 10px 0 0;
  color: var(--ink-3);
  font-size: 10.5px;
  line-height: 1.55;
}
.lookup-note svg { width: 12px; height: 12px; flex: none; margin-top: 1px; }
.lookup-error { margin: 12px 0 0; color: var(--err); font-size: 11.5px; line-height: 1.5; }

.lookup-result { padding: 24px 26px 20px; }
.lookup-result-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}
.lookup-result-head h2 {
  margin: 7px 0 0;
  color: var(--ink);
  font-size: 21px;
  font-weight: 800;
  letter-spacing: -0.025em;
}
.lookup-status {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  flex: none;
  padding: 5px 9px;
  border: 1px solid var(--brd);
  border-radius: 999px;
  background: var(--soft);
  color: var(--ink-3);
  font-size: 10px;
  font-weight: 600;
}
.lookup-status i { width: 5px; height: 5px; border-radius: 50%; background: currentColor; }
.lookup-status.is-good { border-color: var(--good-line, #cde8d9); background: var(--good-soft); color: var(--good); }
.lookup-status.is-err { border-color: color-mix(in srgb, var(--err) 30%, transparent); background: color-mix(in srgb, var(--err) 10%, transparent); color: var(--err); }
.lookup-status.is-info { border-color: color-mix(in srgb, var(--primary) 30%, transparent); background: var(--primary-soft); color: var(--primary); }
.lookup-status.is-warn { border-color: color-mix(in srgb, var(--warn) 30%, transparent); background: var(--warn-soft); color: var(--warn); }
.lookup-status.is-muted { border-color: var(--brd); background: var(--soft); color: var(--ink-3); }

.lookup-receipt {
  margin-top: 16px;
  padding: 2px 16px;
  border: 1px solid var(--brd);
  border-radius: 14px;
  background: var(--surface-2);
}
.lookup-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 18px;
  padding: 11px 0;
  border-top: 1px solid var(--brd);
}
.lookup-row:first-child { border-top: 0; }
.lookup-row span { flex: none; color: var(--ink-3); font-size: 11px; }
.lookup-row b {
  color: var(--ink);
  font-size: 12.5px;
  font-weight: 600;
  text-align: right;
  overflow-wrap: anywhere;
}
.lookup-row b.mono { font-family: var(--font-mono); font-size: 11.5px; font-weight: 500; letter-spacing: -0.01em; }
.lookup-row b.is-good { color: var(--good); }
.lookup-row b.is-err { color: var(--err); }
.lookup-row b.is-info { color: var(--primary); }
.lookup-row b.is-warn { color: var(--warn); }
.lookup-row b.is-muted { color: var(--ink-3); }

.lookup-message {
  display: flex;
  align-items: flex-start;
  gap: 7px;
  margin-top: 12px;
  padding: 10px 12px;
  border-radius: 10px;
  font-size: 11px;
  line-height: 1.55;
}
.lookup-message svg { width: 13px; height: 13px; flex: none; margin-top: 1px; }
.lookup-message.is-good { background: var(--good-soft); color: var(--good); }
.lookup-message.is-err { background: color-mix(in srgb, var(--err) 10%, transparent); color: var(--err); }
.lookup-message.is-info { background: var(--primary-soft); color: var(--primary); }
.lookup-message.is-warn { background: var(--warn-soft); color: var(--warn); }

.lookup-actions { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 16px; }
.lookup-primary,
.lookup-secondary {
  height: 42px;
  padding: 0 18px;
  border-radius: 999px;
  font-family: var(--font-sans);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}
.lookup-primary { border: 0; background: var(--ink); color: var(--bg); }
.lookup-secondary { border: 1px solid var(--brd-2); background: transparent; color: var(--ink); }

.lookup-tools { display: flex; flex-wrap: wrap; gap: 8px; }
.lookup-tool {
  height: 32px;
  padding: 0 12px;
  border: 1px solid var(--brd-2);
  border-radius: 999px;
  background: transparent;
  color: var(--ink-2);
  font-family: var(--font-sans);
  font-size: 10.5px;
  font-weight: 500;
  cursor: pointer;
}
.lookup-count {
  display: inline-grid;
  place-items: center;
  min-width: 22px;
  height: 22px;
  margin-left: 4px;
  padding: 0 6px;
  border-radius: 999px;
  background: var(--soft);
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 500;
  vertical-align: 2px;
}
.lookup-summary {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 8px;
  margin-top: 14px;
}
.lookup-chip {
  padding: 9px 10px;
  border: 1px solid transparent;
  border-radius: 10px;
  text-align: center;
  font-size: 10.5px;
  font-weight: 600;
  line-height: 1.35;
}
.lookup-chip.is-good { background: var(--good-soft); color: var(--good); }
.lookup-chip.is-err { background: color-mix(in srgb, var(--err) 10%, transparent); color: var(--err); }
.lookup-chip.is-muted { background: var(--soft); color: var(--ink-3); }
.lookup-chip.is-info { background: var(--primary-soft); color: var(--primary); }
.lookup-chip.is-warn { background: var(--warn-soft); color: var(--warn); }

.lookup-table-wrap {
  margin-top: 14px;
  overflow-x: auto;
  border: 1px solid var(--brd);
  border-radius: 12px;
}
.lookup-table { width: 100%; border-collapse: collapse; font-size: 11.5px; }
.lookup-table th {
  padding: 10px 12px;
  background: var(--soft);
  color: var(--ink-3);
  font-size: 9.5px;
  font-weight: 500;
  text-align: left;
  white-space: nowrap;
}
.lookup-table td {
  padding: 10px 12px;
  border-top: 1px solid var(--brd);
  color: var(--ink-2);
  vertical-align: middle;
}
.lookup-table .mono { font-family: var(--font-mono); font-size: 10.5px; }
.lookup-index { color: var(--ink-3); font-family: var(--font-mono); font-size: 9.5px; }
.lookup-time { color: var(--ink-3); white-space: nowrap; }
.lookup-status-text { font-size: 10.5px; font-weight: 600; white-space: nowrap; }
.lookup-status-text.is-good { color: var(--good); }
.lookup-status-text.is-err { color: var(--err); }
.lookup-status-text.is-info { color: var(--primary); }
.lookup-status-text.is-warn { color: var(--warn); }
.lookup-status-text.is-muted { color: var(--ink-3); }
.lookup-table-action { text-align: right; white-space: nowrap; }
.lookup-row-action {
  height: 26px;
  padding: 0 9px;
  border: 1px solid var(--brd-2);
  border-radius: 999px;
  background: transparent;
  color: var(--ink-2);
  font-family: var(--font-sans);
  font-size: 10px;
  font-weight: 500;
  cursor: pointer;
}
.lookup-row-action.is-primary { border-color: var(--ink); background: var(--ink); color: var(--bg); }

.lookup-scope {
  padding: 24px 24px 20px;
  border-radius: 18px;
  background: #0d0d0d;
  color: #f4f4f1;
  box-shadow: 0 18px 38px -30px rgba(0, 0, 0, 0.75);
}
.lookup-scope-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: rgba(255, 255, 255, 0.38);
  font-family: var(--font-mono);
  font-size: 8.5px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}
.lookup-scope h3 {
  margin: 20px 0 0;
  color: #f4f4f1;
  font-size: 25px;
  font-weight: 800;
  letter-spacing: -0.025em;
}
.lookup-scope-sub { margin: 8px 0 0; color: rgba(255, 255, 255, 0.38); font-size: 11.5px; line-height: 1.6; }
.lookup-scope-list { margin: 20px 0 0; padding: 0; list-style: none; }
.lookup-scope-list li {
  display: flex;
  align-items: flex-start;
  gap: 11px;
  padding: 11px 0;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}
.lookup-scope-list li:first-child { border-top: 0; }
.lookup-scope-check,
.lookup-scope-cross {
  display: grid;
  place-items: center;
  flex: none;
  width: 21px;
  height: 21px;
  border-radius: 50%;
}
.lookup-scope-check { background: rgba(87, 214, 141, 0.14); color: #57d68d; }
.lookup-scope-cross { background: rgba(248, 113, 113, 0.12); color: #f58b8b; }
.lookup-scope-check svg,
.lookup-scope-cross svg { width: 11px; height: 11px; }
.lookup-scope-list b { display: block; color: rgba(255, 255, 255, 0.82); font-size: 12.5px; font-weight: 600; }
.lookup-scope-list p { margin: 3px 0 0; color: rgba(255, 255, 255, 0.34); font-size: 10.5px; line-height: 1.5; }
.lookup-scope-list li.is-no b { color: rgba(255, 255, 255, 0.55); }
.lookup-scope-foot {
  display: flex;
  align-items: flex-start;
  gap: 7px;
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.4);
  font-size: 9.5px;
  line-height: 1.5;
}
.lookup-scope-foot svg { width: 11px; height: 11px; flex: none; margin-top: 1px; }

.lookup-batch-info { padding: 20px 22px 18px; }
.lookup-batch-info h3 { margin: 0; color: var(--ink); font-size: 12px; font-weight: 800; letter-spacing: 0.02em; }
.lookup-batch-desc { margin: 6px 0 0; color: var(--ink-3); font-size: 10.5px; line-height: 1.6; }
.lookup-batch-list { margin: 10px 0 0; padding: 0; list-style: none; }
.lookup-batch-list li {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 0;
  border-top: 1px solid var(--brd);
}
.lookup-batch-list li:first-child { border-top: 0; }
.lookup-batch-list > li > span {
  padding-top: 2px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 9px;
}
.lookup-batch-list b { display: block; color: var(--ink); font-size: 11.5px; font-weight: 700; }
.lookup-batch-list p { margin: 3px 0 0; color: var(--ink-3); font-size: 10.5px; line-height: 1.5; }
.lookup-batch-note {
  margin: 8px 0 0;
  padding: 10px 0 2px;
  border-top: 1px solid var(--brd);
  color: var(--ink-3);
  font-size: 9.5px;
  line-height: 1.5;
}

@media (max-width: 900px) {
  .lookup-wrap { max-width: calc(100% - 32px); padding-left: 0; padding-right: 0; }
  .lookup-subnav { overflow: hidden; }
  .lookup-hero { flex-direction: column; align-items: stretch; gap: 20px; padding: 24px 0; }
  .lookup-hero h1 { font-size: 42px; }
  .lookup-aside { align-items: flex-start; text-align: left; }
  .lookup-grid { grid-template-columns: minmax(0, 1fr); }
  .lookup-side { max-width: 460px; }
  .lookup-card { padding: 24px 22px 22px; }
}

@media (max-width: 640px) {
  .lookup-wrap { padding-bottom: 32px; }
  .lookup-crumb { display: none; }
  .lookup-aside { display: none; }
  .lookup-subnav { padding: 10px 0 14px; }
  .lookup-hero h1 { font-size: 38px; }
  .lookup-sub { font-size: 12.5px; line-height: 1.7; }
  .lookup-badges { gap: 6px; margin-top: 15px; }
  .lookup-badges li { padding: 5px 9px; font-size: 10px; }
  .lookup-card { padding: 20px 16px 18px; border-radius: 18px; }
  .lookup-mode { height: 32px; font-size: 11.5px; }
  .lookup-head h2 { font-size: 20px; }
  .lookup-hint { padding: 11px 12px; }
  .lookup-hint p { font-size: 11px; }
  .lookup-input-shell { height: 48px; border-radius: 12px; }
  .lookup-input { font-size: 12px; }
  .lookup-textarea { min-height: 150px; font-size: 11.5px; }
  .lookup-cta { height: 44px; font-size: 12.5px; }
  .lookup-result { padding: 18px 16px 16px; border-radius: 18px; }
  .lookup-result-head h2 { font-size: 19px; }
  .lookup-result-head { gap: 10px; }
  .lookup-tools { gap: 6px; }
  .lookup-tool { height: 30px; padding: 0 10px; font-size: 10px; }
  .lookup-summary { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .lookup-table { min-width: 700px; }
  .lookup-table th,
  .lookup-table td { padding: 9px 10px; }
  .lookup-row { gap: 12px; padding: 10px 0; }
  .lookup-row span { font-size: 10px; }
  .lookup-row b { font-size: 11.5px; }
  .lookup-row b.mono { font-size: 10px; }
  .lookup-primary,
  .lookup-secondary { height: 40px; font-size: 11.5px; }
  .lookup-scope { padding: 22px 20px 18px; border-radius: 18px; }
  .lookup-scope h3 { font-size: 23px; }
  .lookup-scope-list b { font-size: 12px; }
  .lookup-scope-list p { font-size: 10px; }
  .lookup-batch-info { padding: 18px 18px 16px; border-radius: 18px; }
}
</style>
