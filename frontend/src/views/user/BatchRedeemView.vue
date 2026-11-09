<template>
  <div class="u-page batch-page">
    <UserPageHeader hide-hero :title="t('batch.title')" />
    <div class="u-main batch-wrap">
      <div class="batch-subnav">
        <div class="batch-crumb">卡密兑换 <i>/</i> <b>{{ t('batch.title') }}</b> <i>/</i> STEP {{ phase === 'input' ? '01' : '03' }}</div>
        <RedeemModeTabs compact />
      </div>

      <section class="batch-hero">
        <div class="batch-hero-copy">
          <div class="batch-eyebrow">{{ phase === 'input' ? 'BATCH REDEEM' : 'BATCH RUNNING' }} <em>· {{ phase === 'input' ? t('hero.batchEyebrow') : t('batchUi.running') }}</em></div>
          <h1>{{ t('batch.title') }}</h1>
          <p>{{ phase === 'input' ? t('batchUi.heroSub') : t('batchUi.runningSub') }}</p>
        </div>
        <div class="batch-hero-actions">
          <router-link class="btn-secondary batch-action" to="/history">{{ t('batchUi.viewRecords') }}</router-link>
        </div>
      </section>

      <div class="batch-grid">
        <section class="card batch-card">
          <!-- input -->
          <template v-if="phase === 'input'">
            <div class="batch-panel-head">
              <div>
                <h2>{{ t('batchUi.setupTitle') }}</h2>
                <p>{{ t('batchUi.setupSub') }}</p>
              </div>
              <span class="batch-state"><i></i>{{ t('batchUi.draft') }}</span>
            </div>

            <ol class="batch-flow" :aria-label="t('batch.title')">
              <li class="is-on"><span class="num">01</span><b>{{ t('batchUi.stepImport') }}</b><em>{{ parsedCdks.length }} {{ t('batchUi.keysUnit') }}</em></li>
              <li><span class="num">02</span><b>{{ t('batchUi.stepCred') }}</b><em>{{ credentialPoolSize }} {{ t('batchUi.credUnit') }}</em></li>
              <li><span class="num">03</span><b>{{ t('batchUi.stepPair') }}</b><em>{{ t('batchUi.pending') }}</em></li>
              <li><span class="num">04</span><b>{{ t('batchUi.stepProgress') }}</b><em>{{ t('batchUi.pending') }}</em></li>
            </ol>

            <section class="batch-section">
              <div class="batch-section-head">
                <div class="batch-section-title">
                  <span class="idx">01</span>
                  <h3>{{ t('batchUi.cdkList') }}</h3>
                  <span class="tag">{{ t('batchUi.maxKeys', { max: BATCH_MAX_KEYS }) }}</span>
                </div>
                <span class="batch-count"><i></i><b>{{ parsedCdks.length }}</b> {{ t('batchUi.keysUnit') }}</span>
              </div>
              <div class="batch-textbox">
                <textarea
                  v-model="cdkText"
                  rows="6"
                  class="batch-textarea"
                  :placeholder="t('batch.cdkPlaceholder')"
                />
                <div class="batch-textbox-foot">
                  <span>{{ t('batchUi.perLine') }}</span>
                  <span class="mono">SXC-XXXX-XXXX-XXXX-XXXX</span>
                </div>
              </div>
              <div class="batch-recognized">
                <span>{{ t('batch.cdkRecognized', { n: parsedCdks.length }) }}</span>
                <span class="right">{{ t('batchUi.supportTypes') }}</span>
              </div>
              <p v-if="parsedCdks.length > 100" class="batch-inline-warn">· {{ t('batch.cdkQueueHint') }}</p>
            </section>

            <section class="batch-section">
              <div class="batch-section-head">
                <div class="batch-section-title">
                  <span class="idx">02</span>
                  <h3>{{ t('batch.credMode') }}</h3>
                </div>
                <span class="batch-section-note">{{ t('batchUi.credHint') }}</span>
              </div>

              <div class="batch-modes">
                <button type="button" class="batch-mode" :class="{ 'is-on': credMode === 'session' }" @click="credMode = 'session'">
                  <span class="batch-mode-icon" aria-hidden="true">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2"/><path d="M15.5 8.5 19 12a5 5 0 0 1-7 7l-1.5-1.5"/><path d="m8 16-4 4"/><path d="m3 21 2-2"/><path d="M14 6a4 4 0 0 0-5.7 0L4 10.3a4 4 0 0 0 0 5.7"/></svg>
                  </span>
                  <span class="batch-mode-copy">
                    <span class="batch-mode-title">{{ t('batch.modeSession') }} <i>{{ t('batchUi.recommended') }}</i></span>
                    <span class="batch-mode-desc">{{ t('batchUi.sessionDesc') }}</span>
                  </span>
                  <span class="batch-mode-check" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="m5 12 4.5 4.5L19 7"/></svg></span>
                </button>
                <button type="button" class="batch-mode" :class="{ 'is-on': credMode === 'mailbox' }" @click="credMode = 'mailbox'">
                  <span class="batch-mode-icon" aria-hidden="true">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="5" width="18" height="14" rx="2"/><path d="m3 7 9 6 9-6"/></svg>
                  </span>
                  <span class="batch-mode-copy">
                    <span class="batch-mode-title">{{ t('batch.modeMailbox') }}</span>
                    <span class="batch-mode-desc">{{ t('batchUi.mailboxDesc') }}</span>
                  </span>
                  <span class="batch-mode-check" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="m5 12 4.5 4.5L19 7"/></svg></span>
                </button>
              </div>

              <template v-if="credMode === 'session'">
                <div class="batch-import-head">
                  <span>{{ t('batchUi.sessionSource') }}</span>
                  <em>{{ t('batchUi.localOnly') }}</em>
                </div>
                <ExcelImportBlock
                  :session-pool="sessionPool"
                  :import-msg="importMsg"
                  :importing="importing"
                  @pick="onPickFile"
                  @clear="clearSessionPool"
                />
                <div class="batch-format-note">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 8v.5"/><path d="M12 11v5"/></svg>
                  <span>{{ t('batchUi.formatNote') }} <a href="#batch-formats">{{ t('batchUi.viewFormat') }} ↗</a></span>
                </div>
              </template>

              <div v-else class="batch-mailbox">
                <label class="batch-label">{{ t('batch.mailboxLabel') }}</label>
                <textarea
                  v-model="mailboxText"
                  rows="6"
                  class="batch-textarea batch-textarea--mailbox"
                  :placeholder="t('batch.mailboxPlaceholder')"
                />
                <p class="batch-recognized single">{{ t('batch.mailboxRecognized', { n: parsedMailboxes.length }) }}</p>
                <input
                  ref="mailboxFileRef"
                  type="file"
                  accept=".xlsx,.xls,.csv,.txt"
                  class="hidden"
                  @change="onPickMailboxFile"
                />
                <button
                  type="button"
                  class="batch-file-btn"
                  :disabled="importing"
                  @click="mailboxFileRef?.click()"
                >
                  {{ importing ? t('batch.importing') : t('batch.mailboxImportBtn') }}
                </button>
                <p v-if="importMsg" class="batch-inline-note">{{ importMsg }}</p>
                <p class="batch-inline-hint">{{ t('batch.mailboxHint') }}</p>
              </div>
            </section>

            <section class="batch-section">
              <div class="batch-section-head">
                <div class="batch-section-title">
                  <span class="idx">03</span>
                  <h3>{{ t('batchUi.pairTitle') }}</h3>
                </div>
                <span class="batch-section-note">{{ t('batchUi.noConsume') }}</span>
              </div>
              <div class="batch-pair">
                <div class="batch-pair-top">
                  <div class="batch-pair-cells">
                    <span class="batch-pair-cell"><span>{{ t('batchUi.cdkShort') }}</span><strong>{{ parsedCdks.length }}</strong></span>
                    <span class="batch-pair-arrow">↔</span>
                    <span class="batch-pair-cell"><span>{{ t('batchUi.credShort') }}</span><strong>{{ credentialPoolSize }}</strong></span>
                  </div>
                  <span class="batch-pair-note">{{ t('batchUi.pairOrder') }}</span>
                </div>
                <div class="batch-pair-bar"><i :style="{ width: inputProgressPct + '%' }"></i></div>
              </div>
              <button
                type="button"
                class="btn-primary batch-submit"
                :disabled="verifying || parsedCdks.length === 0"
                @click="handleVerifyBatch"
              >
                {{ verifying ? t('batch.verifying') : t('batch.verifyStart') }}
              </button>
              <p v-if="parsedCdks.length === 0" class="batch-submit-hint">{{ t('batchUi.needOneKey') }}</p>
              <p v-else class="batch-submit-hint">{{ t('batch.cdkRecognized', { n: parsedCdks.length }) }}</p>
              <p v-if="sessionError" class="batch-inline-error">{{ sessionError }}</p>
              <p class="batch-section-foot">{{ t('batchUi.verifyNote') }}</p>
            </section>
          </template>

          <!-- run -->
          <template v-else>
            <div class="batch-panel-head">
              <div>
                <h2>{{ t('batchUi.running') }}</h2>
                <p>{{ t('batchUi.runningSub') }}</p>
              </div>
              <div class="batch-run-status">
                <span class="batch-run-pill"><i></i>{{ runStatusText }}</span>
                <span class="batch-run-count">{{ sessionsDone }} / {{ items.length }}</span>
              </div>
            </div>

            <ol class="batch-flow" :aria-label="t('batch.title')">
              <li class="is-done"><span class="num">✓</span><b>{{ t('batchUi.stepImport') }}</b><em>{{ items.length }} {{ t('batchUi.keysUnit') }}</em></li>
              <li class="is-done"><span class="num">✓</span><b>{{ t('batchUi.stepCred') }}</b><em>{{ credentialPoolSize }} {{ t('batchUi.credUnit') }}</em></li>
              <li class="is-on"><span class="num">03</span><b>{{ t('batchUi.stepPair') }}</b><em>{{ sessionsDone }} / {{ totalVerified }}</em></li>
              <li><span class="num">04</span><b>{{ t('batchUi.stepProgress') }}</b><em>{{ runStatusText }}</em></li>
            </ol>
            <div class="batch-run-progress"><i :style="{ width: batchProgressPct + '%' }"></i></div>

            <div class="batch-stats">
              <div class="batch-stat"><span>{{ t('batch.statTotal') }}</span><strong>{{ stats.total }}</strong></div>
              <div class="batch-stat is-warn"><span>{{ t('batch.statWait') }}</span><strong>{{ stats.wait }}</strong></div>
              <div class="batch-stat"><span>{{ t('batch.statRun') }}</span><strong>{{ stats.run }}</strong></div>
              <div class="batch-stat is-good"><span>{{ t('batch.statOkFail') }}</span><strong>{{ stats.ok }} <em>/ {{ stats.fail }}</em></strong></div>
            </div>

            <div v-if="!verifying && pendingSessionItems.length > 0 && credMode === 'session'" class="batch-run-import">
              <ExcelImportBlock
                :session-pool="sessionPool"
                :import-msg="importMsg"
                :importing="importing"
                @pick="onPickFile"
                @clear="clearSessionPool"
              />
            </div>
            <div v-else-if="!verifying && pendingSessionItems.length > 0 && credMode === 'mailbox'" class="batch-run-import batch-run-mailbox">
              <label class="batch-label">{{ t('batch.mailboxLabel') }}</label>
              <textarea v-model="mailboxText" rows="4" class="batch-textarea batch-textarea--mailbox" :placeholder="t('batch.mailboxPlaceholder')" />
              <p class="batch-recognized single">{{ t('batch.mailboxRecognized', { n: parsedMailboxes.length }) }}</p>
            </div>

            <div v-if="verifying" class="batch-verifying">
              <span class="batch-spinner" aria-hidden="true"></span>
              <span>{{ t('batch.verifyingKeys') }}</span>
            </div>

            <section v-else-if="currentItem" class="batch-current">
              <div class="batch-current-top">
                <div>
                  <span class="batch-current-kicker">{{ t('batchUi.current') }}</span>
                  <div class="batch-current-title">
                    <b>{{ t('batch.currentOrdinal', { n: currentOrdinal, total: totalVerified }) }}</b>
                    <span>{{ t('batchUi.remaining', { n: remainingToSubmit }) }}</span>
                  </div>
                </div>
                <div class="batch-current-key">
                  <span class="mono">{{ currentItem.cardKey }}</span>
                  <span v-if="currentItem.planLabel" class="plan">{{ displayPlan(currentItem.planLabel) }}</span>
                </div>
              </div>
              <p v-if="currentExcelIdx >= 0 && credMode === 'session' && sessionPool[currentExcelIdx]?.email" class="batch-account-line">
                {{ t('batch.willUseAccount') }} <b>{{ sessionPool[currentExcelIdx].email }}</b>
              </p>
              <p v-else-if="currentExcelIdx >= 0 && credMode === 'mailbox' && mailboxPool[currentExcelIdx]?.email" class="batch-account-line">
                {{ t('batch.willUseAccount') }} <b>{{ mailboxPool[currentExcelIdx].email }}</b>
              </p>
              <textarea
                ref="sessionBoxRef"
                v-model="sessionInput"
                :rows="credMode === 'session' ? 5 : 3"
                class="batch-textarea batch-textarea--session"
                :placeholder="credMode === 'session' ? t('batch.sessionPlaceholder') : t('batch.mailboxLinePlaceholder')"
                :disabled="submitting || autoSubmitting"
                @input="sessionError = ''"
              />
              <p v-if="sessionError" class="batch-inline-error">{{ sessionError }}</p>
              <div class="batch-current-actions">
                <button type="button" class="btn-primary" :disabled="!sessionInput.trim() || submitting || autoSubmitting" @click="handleSubmitSession">
                  {{ submitting ? t('batch.submitting') : t('batch.submitNext') }}
                </button>
                <button
                  v-if="credMode === 'session' ? sessionPool.length > 0 : mailboxPool.length > 0"
                  type="button"
                  class="btn-primary batch-auto"
                  :disabled="submitting || autoSubmitting"
                  @click="handleAutoSubmitAll"
                >
                  {{ autoSubmitting ? t('batch.autoSubmitting') : t('batchUi.autoSubmitShort', { n: autoPoolSize }) }}
                </button>
              </div>
              <p class="batch-current-hint">{{ credMode === 'session' ? t('batch.pairHint') : t('batch.mailboxPairHint') }}</p>
            </section>

            <section v-else-if="allSessionsCollected" class="batch-collected">
              <p>{{ t('batch.allSubmitted', { n: sessionsDone }) }}</p>
              <span>{{ t('batch.watchProgress') }}</span>
              <button v-if="stats.ok > 0" type="button" class="batch-export-link" @click="exportBatchSuccess">{{ t('batchUi.export', { n: stats.ok }) }}</button>
            </section>

            <p v-else class="batch-empty">{{ t('batch.noValidKeys') }}</p>

            <section class="batch-progress">
              <div class="batch-progress-head">
                <h3>{{ t('batchUi.liveProgress') }}</h3>
                <div class="batch-filters">
                  <button type="button" :class="{ on: progressFilter === 'all' }" @click="progressFilter = 'all'">{{ t('batchUi.filterAll', { n: stats.total }) }}</button>
                  <button type="button" :class="{ on: progressFilter === 'running' }" @click="progressFilter = 'running'">{{ t('batchUi.filterRunning', { n: stats.run }) }}</button>
                  <button type="button" :class="{ on: progressFilter === 'success' }" @click="progressFilter = 'success'">{{ t('batchUi.filterSuccess', { n: stats.ok }) }}</button>
                  <button type="button" :class="{ on: progressFilter === 'failed' }" @click="progressFilter = 'failed'">{{ t('batchUi.filterFail', { n: stats.fail }) }}</button>
                </div>
              </div>
              <div class="batch-table-wrap">
                <div class="batch-table">
                  <div class="batch-trow batch-thead">
                    <span>#</span><span>{{ t('batchUi.key') }}</span><span>{{ t('batchUi.plan') }}</span><span>{{ t('batchUi.status') }}</span><span>{{ t('batchUi.account') }}</span><span>{{ t('batchUi.message') }}</span>
                  </div>
                  <div v-for="(it, idx) in filteredItems" :key="it.id" class="batch-trow">
                    <span class="idx">{{ idx + 1 }}</span>
                    <span class="key">{{ it.cardKey }}</span>
                    <span class="plan">{{ displayPlan(it.planLabel) }}</span>
                    <span class="batch-chip" :class="statusChipClass(it.status)">{{ statusLabel(it.status) }}</span>
                    <span class="account"><span>{{ it.email || '—' }}</span><em v-if="it.cardLastFour">•••• {{ it.cardLastFour }}</em></span>
                    <span class="msg">{{ batchItemMessage(it) }}</span>
                  </div>
                </div>
              </div>
              <p v-if="filteredItems.length === 0" class="batch-empty table">{{ t('batchUi.noRows') }}</p>
              <div class="batch-run-actions">
                <button type="button" class="btn-secondary" :disabled="stats.ok === 0" @click="exportBatchSuccess">{{ t('batchUi.export', { n: stats.ok }) }}</button>
                <button type="button" class="btn-secondary" @click="resetAll">{{ t('batchUi.restart') }}</button>
              </div>
            </section>
          </template>
        </section>

        <aside class="batch-aside">
          <template v-if="phase === 'input'">
            <section class="batch-overview">
              <div class="batch-overview-top"><span>BATCH OVERVIEW</span><span>STEP 01 / 04</span></div>
              <div class="batch-overview-main"><strong>{{ parsedCdks.length }}</strong><span>{{ t('batchUi.detected') }}</span></div>
              <div class="batch-overview-bar"><i :style="{ width: inputProgressPct + '%' }"></i></div>
              <div class="batch-overview-rows">
                <div><span>{{ t('batch.credMode') }}</span><b>{{ credMode === 'session' ? t('batch.modeSession') : t('batch.modeMailbox') }}</b></div>
                <div><span>{{ t('batchUi.paired') }}</span><b>{{ pairedCount }} {{ t('batchUi.pairUnit') }}</b></div>
                <div><span>{{ t('batchUi.willSubmit') }}</span><b>{{ pairedCount }} {{ t('batchUi.keysUnit') }}</b></div>
              </div>
              <div class="batch-overview-foot"><i></i>{{ parsedCdks.length ? t('batchUi.noConsume') : t('batchUi.waitImport') }}</div>
            </section>

            <section class="card batch-side-card">
              <h3>{{ t('batchUi.guide') }}</h3>
              <ol class="batch-guide">
                <li class="on"><span>01</span><div><b>{{ t('batchUi.guide1Title') }}</b><p>{{ t('batchUi.guide1Desc', { max: BATCH_MAX_KEYS }) }}</p></div></li>
                <li><span>02</span><div><b>{{ t('batchUi.guide2Title') }}</b><p>{{ t('batchUi.guide2Desc') }}</p></div></li>
                <li><span>03</span><div><b>{{ t('batchUi.guide3Title') }}</b><p>{{ t('batchUi.guide3Desc') }}</p></div></li>
                <li><span>04</span><div><b>{{ t('batchUi.guide4Title') }}</b><p>{{ t('batchUi.guide4Desc') }}</p></div></li>
              </ol>
            </section>

            <section id="batch-formats" class="card batch-side-card">
              <h3>{{ t('batchUi.formats') }}</h3>
              <ul class="batch-formats">
                <li><span>ChatGPT</span><p>{{ t('batchUi.formatChatGPT') }}</p></li>
                <li><span>X</span><p>{{ t('batchUi.formatX') }}</p></li>
                <li><span>Grok</span><p>{{ t('batchUi.formatGrok') }}</p></li>
              </ul>
              <p class="batch-lock">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2"/><path d="M8 10V7a4 4 0 0 1 8 0v3"/></svg>
                {{ t('batchUi.localNote') }}
              </p>
            </section>

            <p class="batch-note">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 8v.5"/><path d="M12 11v5"/></svg>
              {{ t('batchUi.runningNote') }}
            </p>
          </template>

          <template v-else>
            <section class="batch-overview">
              <div class="batch-overview-top"><span>BATCH OVERVIEW</span><span>RUNNING</span></div>
              <div class="batch-overview-main"><strong>{{ stats.ok }}</strong><span>{{ t('batch.status.success') }}</span></div>
              <div class="batch-overview-bar"><i :style="{ width: batchProgressPct + '%' }"></i></div>
              <div class="batch-overview-rows">
                <div><span>{{ t('batch.statTotal') }}</span><b>{{ stats.total }} {{ t('batchUi.keysUnit') }}</b></div>
                <div><span>{{ t('batchUi.submitted') }}</span><b>{{ submittedCount }} / {{ stats.total }}</b></div>
                <div><span>{{ t('batch.statRun') }}</span><b>{{ stats.run }} {{ t('batchUi.keysUnit') }}</b></div>
                <div><span>{{ t('batch.status.failed') }}</span><b>{{ stats.fail }} {{ t('batchUi.keysUnit') }}</b></div>
              </div>
              <div class="batch-overview-foot"><i></i>{{ t('batchUi.liveProgress') }} · {{ BATCH_POLL_INTERVAL_MS / 1000 }}s</div>
            </section>

            <section class="card batch-side-card">
              <h3>{{ t('batchUi.batchInfo') }}</h3>
              <ul class="batch-info">
                <li><span>{{ t('batch.credMode') }}</span><b>{{ credMode === 'session' ? t('batch.modeSession') : t('batch.modeMailbox') }}</b></li>
                <li><span>{{ t('batchUi.pool') }}</span><b>{{ credentialPoolSize }} {{ t('batchUi.credUnit') }}</b></li>
                <li><span>{{ t('batchUi.concurrency') }}</span><b>{{ AUTO_SUBMIT_CONCURRENCY }}</b></li>
                <li><span>{{ t('batchUi.startedAt') }}</span><b>{{ batchStartedAt || '—' }}</b></li>
              </ul>
            </section>

            <section class="card batch-side-card">
              <h3>{{ t('batchUi.actions') }}</h3>
              <div class="batch-side-actions">
                <button type="button" :disabled="stats.ok === 0" @click="exportBatchSuccess">{{ t('batchUi.export', { n: stats.ok }) }} <span>↗</span></button>
                <button type="button" @click="resetAll">{{ t('batchUi.restart') }} <span>↗</span></button>
              </div>
            </section>

            <p class="batch-note">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 8v.5"/><path d="M12 11v5"/></svg>
              {{ t('batchUi.runningNote') }}
            </p>
          </template>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import UserPageHeader from '../../components/UserPageHeader.vue'
import RedeemModeTabs from '../../components/RedeemModeTabs.vue'
import ExcelImportBlock from '../../components/ExcelImportBlock.vue'
import { isGrokCredential } from '../../lib/grok'
import {
  AUTO_SUBMIT_CONCURRENCY,
  BATCH_MAX_KEYS,
  BATCH_POLL_INTERVAL_MS,
  VERIFY_CONCURRENCY,
  accessTokenFromSession,
  checkSessionForCdk,
  exportSuccessWorkbook,
  extractCdkSession,
  mapPool,
  parseCdks,
  parseMailboxLine,
  parseMailboxLines,
  parseMailboxesFromSheet,
  parseSessionsFromSheet,
  readWorkbookRows,
  type ImportedMailbox,
  type ImportedSession,
} from '../../lib/batch-session'
import { planLabel } from '../../lib/plan'

const { t } = useI18n({ useScope: 'global' })

type ItemStatus =
  | 'pending_verify'
  | 'verify_ok'
  | 'verify_fail'
  | 'waiting_session'
  | 'submitting'
  | 'processing'
  | 'success'
  | 'failed'

interface BatchItem {
  id: string
  cardKey: string
  planLabel: string
  status: ItemStatus
  verifyMsg: string
  redemptionToken: string
  progressMsg: string
  email: string
  cardLastFour: string
  accessToken: string
  gptPassword: string
  emailPassword: string
  error: string
}

function extractCardLastFour(order: any): string {
  const last = String(order?.card_last_four || '').trim()
  if (/^\d{4}$/.test(last)) return last
  const n = String(order?.card_number || '').replace(/\D/g, '')
  return n.length >= 4 ? n.slice(-4) : ''
}

const deviceId = (() => {
  const k = 'cdk_device_id'
  let v = localStorage.getItem(k)
  if (!v) {
    v = 'web-' + Math.random().toString(36).slice(2) + Date.now().toString(36)
    localStorage.setItem(k, v)
  }
  return v
})()

const phase = ref<'input' | 'run'>('input')
const credMode = ref<'session' | 'mailbox'>('session')
const batchStartedAt = ref('')
const progressFilter = ref<'all' | 'running' | 'success' | 'failed'>('all')
const cdkText = ref('')
const mailboxText = ref('')
const items = ref<BatchItem[]>([])
const verifying = ref(false)
const sessionInput = ref('')
const sessionError = ref('')
const submitting = ref(false)
const autoSubmitting = ref(false)
const sessionPool = ref<ImportedSession[]>([])
const mailboxPool = ref<ImportedMailbox[]>([])
const importMsg = ref('')
const importing = ref(false)
const sessionBoxRef = ref<HTMLTextAreaElement | null>(null)
const mailboxFileRef = ref<HTMLInputElement | null>(null)
const lastFilledItemId = ref<string | null>(null)

const pollTargets = ref<Record<string, { token: string; terminal: boolean }>>({})
let pollTimer: ReturnType<typeof setInterval> | null = null
let pollInFlight = false

const parsedCdks = computed(() => parseCdks(cdkText.value))
const parsedMailboxes = computed(() => parseMailboxLines(mailboxText.value))

watch(parsedMailboxes, (rows) => {
  // 文本粘贴优先；Excel 导入写入 mailboxPool 后也会同步到 mailboxText
  if (credMode.value === 'mailbox') mailboxPool.value = rows
})

const verifiedQueue = computed(() =>
  items.value.filter((i) => i.status !== 'pending_verify' && i.status !== 'verify_fail'),
)
const pendingSessionItems = computed(() =>
  items.value.filter((i) => i.status === 'waiting_session' || i.status === 'verify_ok'),
)
const currentItem = computed(() => pendingSessionItems.value[0] ?? null)
const currentExcelIdx = computed(() => {
  if (!currentItem.value) return -1
  return verifiedQueue.value.findIndex((i) => i.id === currentItem.value!.id)
})
const sessionsDone = computed(() =>
  items.value.filter((i) => ['submitting', 'processing', 'success', 'failed'].includes(i.status)).length,
)
const totalVerified = computed(() => verifiedQueue.value.length)
const currentOrdinal = computed(() =>
  currentExcelIdx.value >= 0 ? currentExcelIdx.value + 1 : sessionsDone.value,
)
const remainingToSubmit = computed(() => Math.max(0, pendingSessionItems.value.length - 1))

const stats = computed(() => {
  const s = { total: items.value.length, ok: 0, fail: 0, run: 0, wait: 0 }
  for (const i of items.value) {
    if (i.status === 'success') s.ok++
    else if (i.status === 'failed' || i.status === 'verify_fail') s.fail++
    else if (i.status === 'processing' || i.status === 'submitting') s.run++
    else if (i.status === 'waiting_session' || i.status === 'verify_ok') s.wait++
  }
  return s
})

const credentialPoolSize = computed(() => (credMode.value === 'session' ? sessionPool.value.length : mailboxPool.value.length))
const pairedCount = computed(() => Math.min(parsedCdks.value.length, credentialPoolSize.value))
const inputProgressPct = computed(() => (parsedCdks.value.length ? Math.round((pairedCount.value / parsedCdks.value.length) * 100) : 0))
const autoPoolSize = computed(() => credentialPoolSize.value)
const submittedCount = computed(() => stats.value.ok + stats.value.fail + stats.value.run)
const batchProgressPct = computed(() => (stats.value.total ? Math.round((submittedCount.value / stats.value.total) * 100) : 0))
const runStatusText = computed(() => (
  verifying.value ? t('batch.verifying') : allSessionsCollected.value ? t('batchUi.completed') : t('batch.statRun')
))
const filteredItems = computed(() => {
  if (progressFilter.value === 'running') {
    return items.value.filter((item) => item.status === 'submitting' || item.status === 'processing')
  }
  if (progressFilter.value === 'success') return items.value.filter((item) => item.status === 'success')
  if (progressFilter.value === 'failed') return items.value.filter((item) => item.status === 'failed' || item.status === 'verify_fail')
  return items.value
})

function statusChipClass(status: ItemStatus) {
  if (status === 'success' || status === 'verify_ok') return 'is-good'
  if (status === 'failed' || status === 'verify_fail') return 'is-bad'
  if (status === 'waiting_session') return 'is-warn'
  if (status === 'submitting' || status === 'processing') return 'is-info'
  return 'is-muted'
}

function batchItemMessage(item: BatchItem) {
  if (item.status === 'verify_fail' || item.status === 'failed') {
    return item.error || item.progressMsg || item.verifyMsg || '—'
  }
  return item.progressMsg || item.verifyMsg || '—'
}

function displayPlan(value: string) {
  return value ? planLabel(value) : '—'
}

const allSessionsCollected = computed(
  () =>
    phase.value === 'run' &&
    !verifying.value &&
    totalVerified.value > 0 &&
    pendingSessionItems.value.length === 0,
)

function statusLabel(st: ItemStatus) {
  const map: Record<ItemStatus, string> = {
    pending_verify: t('batch.status.pending_verify'),
    verify_ok: t('batch.status.verify_ok'),
    verify_fail: t('batch.status.verify_fail'),
    waiting_session: t('batch.status.waiting_session'),
    submitting: t('batch.status.submitting'),
    processing: t('batch.status.processing'),
    success: t('batch.status.success'),
    failed: t('batch.status.failed'),
  }
  return map[st] || st
}

function updateItem(id: string, patch: Partial<BatchItem>) {
  items.value = items.value.map((it) => (it.id === id ? { ...it, ...patch } : it))
}

async function api(path: string, init: RequestInit = {}) {
  const headers = new Headers(init.headers || {})
  headers.set('X-Redemption-Device', deviceId)
  if (init.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json')
  const r = await fetch(path, { ...init, headers, credentials: 'include' })
  const text = await r.text()
  let data: any = null
  try {
    data = text ? JSON.parse(text) : null
  } catch {
    data = { raw: text }
  }
  return { r, data }
}

const TERMINAL = new Set(['completed', 'declined', 'failed_precharge', 'cancelled', 'failed'])

function isTerminal(st: string) {
  return TERMINAL.has(String(st || '').toLowerCase())
}

function isSuccessStatus(st: string) {
  return String(st || '').toLowerCase() === 'completed'
}

function isFailStatus(st: string) {
  const s = String(st || '').toLowerCase()
  return ['declined', 'failed_precharge', 'cancelled', 'failed'].includes(s)
}

function planLabelOf(data: any): string {
  return String(data?.plan || data?.plan_type || data?.plan_label || data?.data?.plan || '')
}

function stopAllPolls() {
  if (pollTimer) clearInterval(pollTimer)
  pollTimer = null
  pollInFlight = false
  pollTargets.value = {}
}

async function pollBatch() {
  if (pollInFlight) return
  const entries = Object.entries(pollTargets.value).filter(([, v]) => !v.terminal)
  if (!entries.length) {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
    return
  }
  pollInFlight = true
  try {
    await Promise.all(
      entries.map(async ([id, target]) => {
        try {
          const { r, data } = await api(
            '/api/v1/public/cdk/result?token=' + encodeURIComponent(target.token),
          )
          if (!r.ok) return
          const order = data?.order || data?.data?.order || data?.data || data || {}
          const st = String(order.status || data?.status || data?.data?.status || '')
          const message = String(
            order.message || order.user_message || data?.message || data?.user_message || '',
          )
          const email = String(order.account_email || order.email || '').trim()
          const cardLastFour = extractCardLastFour(order)
          const meta: Partial<BatchItem> = {}
          if (email) meta.email = email
          if (cardLastFour) meta.cardLastFour = cardLastFour
          if (isSuccessStatus(st)) {
            updateItem(id, {
              ...meta,
              status: 'success',
              progressMsg: message || '开通完成',
              error: '',
            })
            pollTargets.value[id] = { ...target, terminal: true }
          } else if (isFailStatus(st)) {
            updateItem(id, {
              ...meta,
              status: 'failed',
              progressMsg: message || '兑换失败',
              error: message || '失败',
            })
            pollTargets.value[id] = { ...target, terminal: true }
          } else {
            updateItem(id, {
              ...meta,
              status: 'processing',
              progressMsg: message || (st ? `处理中 · ${st}` : '处理中…'),
              error: '',
            })
          }
        } catch {
          /* transient */
        }
      }),
    )
  } finally {
    pollInFlight = false
  }
}

function startPoll(id: string, token: string) {
  pollTargets.value = { ...pollTargets.value, [id]: { token, terminal: false } }
  if (!pollTimer) {
    void pollBatch()
    pollTimer = setInterval(() => void pollBatch(), BATCH_POLL_INTERVAL_MS)
  }
}

function fillCredentialFromPool(excelIdx: number) {
  if (credMode.value === 'mailbox') {
    const pool = mailboxPool.value
    if (!pool.length || excelIdx < 0 || excelIdx >= pool.length) {
      sessionInput.value = ''
      return
    }
    const m = pool[excelIdx]
    sessionInput.value = `${m.email}----${m.password}`
    sessionError.value = ''
    return
  }
  const pool = sessionPool.value
  if (!pool.length || excelIdx < 0 || excelIdx >= pool.length) {
    sessionInput.value = ''
    return
  }
  sessionInput.value = pool[excelIdx].session
  sessionError.value = ''
}

watch(
  () => currentItem.value?.id,
  async (id) => {
    if (!id || currentExcelIdx.value < 0) return
    if (lastFilledItemId.value === id) return
    lastFilledItemId.value = id
    const hasPool =
      credMode.value === 'mailbox' ? mailboxPool.value.length > 0 : sessionPool.value.length > 0
    if (hasPool) fillCredentialFromPool(currentExcelIdx.value)
    else sessionInput.value = ''
    await nextTick()
    sessionBoxRef.value?.focus()
  },
)

async function onPickFile(file: File | null) {
  if (!file) return
  importing.value = true
  importMsg.value = ''
  try {
    const rows = await readWorkbookRows(file)
    const { sessions, mode, sessionCol, skippedDup } = parseSessionsFromSheet(rows)
    if (!sessions.length) {
      sessionPool.value = []
      importMsg.value =
        sessionCol < 0
          ? '未识别到 Session 列。请确保有 session 列或单元格为完整 Session JSON'
          : '识别到列但无有效 Session（需完整 JSON 含 sessionToken，或五段 JWE）'
      return
    }
    sessionPool.value = sessions
    importMsg.value =
      `已从「${file.name}」识别 ${sessions.length} 条 Session（${mode}）` +
      (skippedDup > 0 ? `，已去重跳过 ${skippedDup} 条` : '') +
      (sessions[0]?.email ? `，如 ${sessions[0].email}` : '')
    lastFilledItemId.value = null
    if (phase.value === 'run' && currentExcelIdx.value >= 0) {
      fillCredentialFromPool(currentExcelIdx.value)
    }
  } catch (e) {
    sessionPool.value = []
    importMsg.value = '读取 Excel 失败：' + (e instanceof Error ? e.message : '未知错误')
  } finally {
    importing.value = false
  }
}

async function onPickMailboxFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0] ?? null
  input.value = ''
  if (!file) return
  importing.value = true
  importMsg.value = ''
  try {
    const name = (file.name || '').toLowerCase()
    if (name.endsWith('.txt')) {
      const text = await file.text()
      const rows = parseMailboxLines(text)
      if (!rows.length) {
        importMsg.value = '未识别到邮箱密码行（支持 email----password）'
        return
      }
      mailboxText.value = rows.map((r) => `${r.email}----${r.password}`).join('\n')
      mailboxPool.value = rows
      importMsg.value = `已从「${file.name}」识别 ${rows.length} 条邮箱密码`
      return
    }
    const sheetRows = await readWorkbookRows(file)
    const { mailboxes, mode, skippedDup } = parseMailboxesFromSheet(sheetRows)
    if (!mailboxes.length) {
      importMsg.value = '未识别到邮箱/密码列。请使用「邮箱,邮箱密码」或 email----password'
      return
    }
    mailboxText.value = mailboxes.map((r) => `${r.email}----${r.password}`).join('\n')
    mailboxPool.value = mailboxes
    importMsg.value =
      `已从「${file.name}」识别 ${mailboxes.length} 条邮箱密码（${mode}）` +
      (skippedDup > 0 ? `，已去重跳过 ${skippedDup} 条` : '')
    lastFilledItemId.value = null
    if (phase.value === 'run' && currentExcelIdx.value >= 0) {
      fillCredentialFromPool(currentExcelIdx.value)
    }
  } catch (err) {
    importMsg.value = '读取失败：' + (err instanceof Error ? err.message : '未知错误')
  } finally {
    importing.value = false
  }
}

function clearSessionPool() {
  sessionPool.value = []
  importMsg.value = ''
  sessionInput.value = ''
}

async function handleVerifyBatch() {
  const keys = parseCdks(cdkText.value)
  if (!keys.length) return
  if (keys.length > BATCH_MAX_KEYS) {
    sessionError.value = `单次最多 ${BATCH_MAX_KEYS} 个卡密`
    return
  }
  verifying.value = true
  sessionError.value = ''
  stopAllPolls()
  lastFilledItemId.value = null
  if (credMode.value === 'mailbox') {
    mailboxPool.value = parsedMailboxes.value
  }
  const list: BatchItem[] = keys.map((k, i) => ({
    id: `${i}-${k}`,
    cardKey: k,
    planLabel: '',
    status: 'pending_verify',
    verifyMsg: '验证中…',
    redemptionToken: '',
    progressMsg: '',
    email: '',
    cardLastFour: '',
    accessToken: '',
    gptPassword: '',
    emailPassword: '',
    error: '',
  }))
  items.value = list
  phase.value = 'run'
  batchStartedAt.value = new Date().toLocaleString()
  sessionInput.value = ''

  try {
    const results = await mapPool(list, VERIFY_CONCURRENCY, async (item) => {
      try {
        const { r, data } = await api('/api/v1/public/cdk/preview', {
          method: 'POST',
          body: JSON.stringify({ code: item.cardKey }),
        })
        if (!r.ok) {
          const msg = String(data?.error || data?.msg || data?.message || 'CDK 无效或不可用')
          return {
            ...item,
            status: 'verify_fail' as const,
            verifyMsg: msg,
            error: msg,
          }
        }
        const token = String(
          data?.redemption_token || data?.data?.redemption_token || data?.token || '',
        )
        if (!token) {
          return {
            ...item,
            status: 'verify_fail' as const,
            verifyMsg: '未返回 redemption_token',
            error: '未返回 redemption_token',
          }
        }
        const body = data?.data || data || {}
        return {
          ...item,
          status: 'waiting_session' as const,
          verifyMsg: '有效',
          planLabel: planLabelOf(body),
          redemptionToken: token,
        }
      } catch {
        return {
          ...item,
          status: 'verify_fail' as const,
          verifyMsg: '验证请求失败',
          error: '验证请求失败',
        }
      }
    })
    items.value = results
  } finally {
    verifying.value = false
    lastFilledItemId.value = null
  }
}

async function submitOne(
  item: BatchItem,
  credential:
    | { mode: 'session'; session: string; imported?: ImportedSession }
    | { mode: 'mailbox'; email: string; password: string },
): Promise<'ok' | 'fail'> {
  if (credential.mode === 'session') {
    const check = checkSessionForCdk(credential.session)
    if (!check.ok) {
      updateItem(item.id, {
        status: 'failed',
        progressMsg: check.error,
        error: check.error,
      })
      return 'fail'
    }
    const session = extractCdkSession(credential.session)
    const imported = credential.imported
    updateItem(item.id, {
      status: 'submitting',
      email: check.email || imported?.email || '',
      accessToken: imported?.accessToken || accessTokenFromSession(credential.session),
      gptPassword: imported?.gptPassword || '',
      emailPassword: imported?.emailPassword || '',
      progressMsg: '账号预检…',
      error: '',
    })
    return runPreflightRedeem(item, { mode: 'session', session }, check.email || imported?.email || '')
  }

  const email = credential.email.trim()
  const password = credential.password
  if (!email.includes('@') || !password) {
    updateItem(item.id, {
      status: 'failed',
      progressMsg: '邮箱或密码无效',
      error: '邮箱或密码无效',
    })
    return 'fail'
  }
  updateItem(item.id, {
    status: 'submitting',
    email,
    emailPassword: password,
    progressMsg: '账号预检…',
    error: '',
  })
  return runPreflightRedeem(item, { mode: 'mailbox', email, password }, email)
}

async function runPreflightRedeem(
  item: BatchItem,
  credential:
    | { mode: 'session'; session: string }
    | { mode: 'mailbox'; email: string; password: string },
  fallbackEmail: string,
): Promise<'ok' | 'fail'> {
  let redemptionToken = item.redemptionToken
  try {
    if (!redemptionToken) {
      const preview = await api('/api/v1/public/cdk/preview', {
        method: 'POST',
        body: JSON.stringify({ code: item.cardKey }),
      })
      if (!preview.r.ok) {
        const msg = String(preview.data?.error || preview.data?.msg || 'CDK 预览失败')
        updateItem(item.id, { status: 'failed', progressMsg: msg, error: msg })
        return 'fail'
      }
      redemptionToken = String(
        preview.data?.redemption_token || preview.data?.data?.redemption_token || '',
      )
      if (!redemptionToken) {
        updateItem(item.id, {
          status: 'failed',
          progressMsg: '未返回 redemption_token',
          error: '未返回 redemption_token',
        })
        return 'fail'
      }
      updateItem(item.id, { redemptionToken })
    }

    const { r, data } = await api('/api/v1/public/cdk/preflight', {
      method: 'POST',
      body: JSON.stringify({
        code: item.cardKey,
        redemption_token: redemptionToken,
        credential,
      }),
    })
    if (!r.ok || (data && typeof data.code === 'number' && data.code !== 0)) {
      const msg = String(data?.error || data?.msg || data?.message || '预检失败')
      updateItem(item.id, { status: 'failed', progressMsg: msg, error: msg })
      return 'fail'
    }
    const body = data?.data && typeof data.data === 'object' ? data.data : data || {}
    const preflightToken = String(body.preflight_token || data?.preflight_token || '')
    // Grok 码：预检不产票，兑换直接带登录态；卡台给出的拦截原因（登录态失效/已有订阅）直接判失败。
    const grok = credential.mode === 'session' && isGrokCredential(credential.session)
    const blocking = grok ? String(body.blocking_reason || '') : ''
    if (blocking) {
      updateItem(item.id, { status: 'failed', progressMsg: blocking, error: blocking })
      return 'fail'
    }
    if (!preflightToken && !grok) {
      updateItem(item.id, {
        status: 'failed',
        progressMsg: '未返回 preflight_token',
        error: '未返回 preflight_token',
      })
      return 'fail'
    }
    const email = String(body.email || body.account_email || fallbackEmail || '')
    updateItem(item.id, { email, progressMsg: '预检通过，提交兑换…' })

    const client_request_id = `batch-${deviceId.slice(0, 8)}-${Date.now()}-${item.id}`
    const redeem = await api('/api/v1/public/cdk/redeem', {
      method: 'POST',
      body: JSON.stringify(grok
        ? { redemption_token: redemptionToken, credential, client_request_id }
        : { redemption_token: redemptionToken, preflight_token: preflightToken, client_request_id }),
    })
    const order =
      redeem.data?.order || redeem.data?.data?.order || redeem.data?.data || redeem.data || {}
    const st = String(order.status || redeem.data?.status || '')
    const message = String(
      order.message || order.user_message || redeem.data?.message || redeem.data?.msg || '',
    )
    const cardLastFour = extractCardLastFour(order)
    const orderEmail = String(order.account_email || order.email || email).trim()

    if (!redeem.r.ok && redeem.r.status !== 202) {
      updateItem(item.id, {
        status: 'failed',
        email: orderEmail || email,
        cardLastFour: cardLastFour || item.cardLastFour,
        progressMsg: message || redeem.data?.error || '兑换被拒绝',
        error: message || redeem.data?.error || '兑换被拒绝',
      })
      return 'fail'
    }

    if (isSuccessStatus(st)) {
      updateItem(item.id, {
        status: 'success',
        email: orderEmail || email,
        cardLastFour,
        progressMsg: message || '开通完成',
        error: '',
      })
    } else if (isFailStatus(st)) {
      updateItem(item.id, {
        status: 'failed',
        email: orderEmail || email,
        cardLastFour,
        progressMsg: message || '兑换失败',
        error: message || '失败',
      })
      return 'fail'
    } else {
      updateItem(item.id, {
        status: 'processing',
        email: orderEmail || email,
        cardLastFour,
        progressMsg: message || '处理中…',
        error: '',
      })
      startPoll(item.id, redemptionToken)
    }
    return 'ok'
  } catch {
    updateItem(item.id, {
      status: 'failed',
      progressMsg: '网络异常，请稍后重试或联系客服确认是否已受理',
      error: '网络异常',
    })
    return 'fail'
  }
}

async function handleSubmitSession() {
  if (!currentItem.value || submitting.value || autoSubmitting.value) return
  const raw = sessionInput.value.trim()
  if (credMode.value === 'mailbox') {
    const parsed = parseMailboxLine(raw)
    if (!parsed) {
      sessionError.value = '请填写 email----password（或 email:password）'
      return
    }
    sessionError.value = ''
    submitting.value = true
    const item = currentItem.value
    sessionInput.value = ''
    lastFilledItemId.value = null
    try {
      await submitOne(item, { mode: 'mailbox', email: parsed.email, password: parsed.password })
    } finally {
      submitting.value = false
    }
    return
  }

  const check = checkSessionForCdk(raw)
  if (!check.ok) {
    sessionError.value = check.error
    return
  }
  sessionError.value = ''
  submitting.value = true
  const item = currentItem.value
  sessionInput.value = ''
  lastFilledItemId.value = null
  try {
    const imported =
      currentExcelIdx.value >= 0 ? sessionPool.value[currentExcelIdx.value] : undefined
    await submitOne(item, {
      mode: 'session',
      session: raw,
      imported: imported?.session.trim() === raw ? imported : undefined,
    })
  } finally {
    submitting.value = false
  }
}

async function handleAutoSubmitAll() {
  if (autoSubmitting.value || submitting.value || verifying.value) return
  const verifiedSnap = items.value.filter(
    (i) => i.status !== 'pending_verify' && i.status !== 'verify_fail',
  )
  const pendingSnap = items.value.filter(
    (i) => i.status === 'verify_ok' || i.status === 'waiting_session',
  )
  if (!pendingSnap.length) {
    sessionError.value = '没有待提交的卡密'
    return
  }

  if (credMode.value === 'mailbox') {
    const pool = mailboxPool.value.length ? mailboxPool.value : parsedMailboxes.value
    if (!pool.length) {
      sessionError.value = '请先粘贴或导入邮箱密码列表'
      return
    }
    autoSubmitting.value = true
    sessionError.value = ''
    const verifiedPositions = new Map(verifiedSnap.map((item, index) => [item.id, index]))
    let cursor = 0
    const workers = Array.from(
      { length: Math.min(AUTO_SUBMIT_CONCURRENCY, pendingSnap.length) },
      async () => {
        while (true) {
          const itemIndex = cursor++
          if (itemIndex >= pendingSnap.length) return
          const item = pendingSnap[itemIndex]
          const excelIdx = verifiedPositions.get(item.id) ?? -1
          const mb = excelIdx >= 0 ? pool[excelIdx] : undefined
          if (!mb) {
            updateItem(item.id, {
              status: 'failed',
              progressMsg: `邮箱列表只有 ${pool.length} 条，第 ${excelIdx + 1} 张无对应`,
              error: '邮箱密码不足',
            })
            continue
          }
          await submitOne(item, { mode: 'mailbox', email: mb.email, password: mb.password })
        }
      },
    )
    await Promise.all(workers)
    sessionInput.value = ''
    autoSubmitting.value = false
    return
  }

  const pool = sessionPool.value
  if (!pool.length) {
    sessionError.value = '请先导入 Excel Session'
    return
  }
  autoSubmitting.value = true
  sessionError.value = ''
  const verifiedPositions = new Map(verifiedSnap.map((item, index) => [item.id, index]))
  let cursor = 0
  const workers = Array.from(
    { length: Math.min(AUTO_SUBMIT_CONCURRENCY, pendingSnap.length) },
    async () => {
      while (true) {
        const itemIndex = cursor++
        if (itemIndex >= pendingSnap.length) return
        const item = pendingSnap[itemIndex]
        const excelIdx = verifiedPositions.get(item.id) ?? -1
        const sess = excelIdx >= 0 ? pool[excelIdx] : undefined
        if (!sess) {
          updateItem(item.id, {
            status: 'failed',
            progressMsg: `Excel 只有 ${pool.length} 条 Session，第 ${excelIdx + 1} 张无对应`,
            error: 'Session 不足',
          })
          continue
        }
        await submitOne(item, { mode: 'session', session: sess.session, imported: sess })
      }
    },
  )
  await Promise.all(workers)
  sessionInput.value = ''
  autoSubmitting.value = false
}

async function exportBatchSuccess() {
  const ok = items.value.filter((i) => i.status === 'success')
  if (!ok.length) {
    sessionError.value = '本批尚无成功记录'
    return
  }
  sessionError.value = ''
  const part = (value: string) => String(value || '').replace(/[\r\n]+/g, ' ').trim()
  try {
    await exportSuccessWorkbook(
      ok.map((i) => [part(i.email), part(i.gptPassword), part(i.emailPassword), part(i.accessToken)]),
    )
  } catch {
    sessionError.value = '导出失败，请重试'
  }
}

function resetAll() {
  stopAllPolls()
  phase.value = 'input'
  items.value = []
  cdkText.value = ''
  mailboxText.value = ''
  mailboxPool.value = []
  sessionPool.value = []
  importMsg.value = ''
  sessionInput.value = ''
  sessionError.value = ''
  autoSubmitting.value = false
  lastFilledItemId.value = null
  batchStartedAt.value = ''
  progressFilter.value = 'all'
}

onUnmounted(() => stopAllPolls())
</script>

<style scoped>
.text-good { color: var(--good, #16a34a); }
.text-warn { color: var(--warn, #d97706); }
.text-err { color: var(--err, #dc2626); }
.status-pill[data-st='pending_verify'] { background: var(--soft); color: var(--muted); }
.status-pill[data-st='verify_ok'],
.status-pill[data-st='success'] {
  background: color-mix(in srgb, var(--good, #16a34a) 16%, transparent);
  color: var(--good, #16a34a);
}
.status-pill[data-st='verify_fail'],
.status-pill[data-st='failed'] {
  background: color-mix(in srgb, var(--err, #dc2626) 14%, transparent);
  color: var(--err, #dc2626);
}
.status-pill[data-st='waiting_session'] {
  background: color-mix(in srgb, var(--warn, #d97706) 14%, transparent);
  color: var(--warn, #d97706);
}
.status-pill[data-st='submitting'],
.status-pill[data-st='processing'] {
  background: color-mix(in srgb, var(--primary) 14%, transparent);
  color: var(--primary);
}

/* ── batch page layout ─────────────────────────────────────────────── */
.batch-page { min-height: 100vh; padding-bottom: 64px; }
.batch-wrap { max-width: 72rem; margin: 0 auto; padding: 0 24px; }
.batch-subnav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 8px 0 18px;
  border-bottom: 1px solid var(--brd);
}
.batch-crumb {
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  white-space: nowrap;
}
.batch-crumb b { color: var(--ink); font-weight: 600; }
.batch-crumb i { margin: 0 8px; color: var(--ink-3); font-style: normal; opacity: 0.65; }

.batch-hero {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 36px;
  padding: 32px 0 28px;
}
.batch-hero-copy { min-width: 0; }
.batch-eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}
.batch-eyebrow::before { content: ''; width: 18px; height: 1px; background: var(--ink-3); }
.batch-eyebrow em { color: var(--ink); font-style: normal; }
.batch-hero h1 {
  margin: 12px 0 0;
  color: var(--ink);
  font-family: var(--title-font, var(--font-display));
  font-size: clamp(2.2rem, 4.4vw, 3.125rem);
  font-weight: 800;
  letter-spacing: -0.04em;
  line-height: 1.04;
}
.batch-hero p { max-width: 540px; margin: 13px 0 0; color: var(--ink-2); font-size: 14px; line-height: 1.7; }
.batch-hero-actions { display: flex; align-items: center; gap: 10px; flex: none; padding-bottom: 4px; }
.batch-action { height: 42px; padding: 0 20px; font-size: 13px; white-space: nowrap; }

.batch-grid { display: grid; grid-template-columns: minmax(0, 1fr) 360px; gap: 22px; align-items: start; padding-bottom: 28px; }
.batch-card {
  padding: 28px 30px 24px;
  min-width: 0;
  box-shadow: 0 1px 0 color-mix(in srgb, var(--surface) 90%, transparent) inset, 0 26px 50px -46px rgba(0, 0, 0, 0.35);
}
.batch-panel-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 20px; }
.batch-panel-head h2 { margin: 0; color: var(--ink); font-size: 21px; font-weight: 800; letter-spacing: -0.02em; }
.batch-panel-head p { max-width: 520px; margin: 6px 0 0; color: var(--ink-2); font-size: 12.5px; line-height: 1.6; }
.batch-state {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex: none;
  height: 24px;
  padding: 0 10px;
  border: 1px solid var(--brd);
  border-radius: 999px;
  background: var(--surface-2);
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 9.5px;
  letter-spacing: 0.1em;
  white-space: nowrap;
}
.batch-state i { width: 5px; height: 5px; border-radius: 50%; background: var(--ink-3); }

.batch-flow {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  margin: 24px 0 0;
  padding: 0;
  list-style: none;
}
.batch-flow li { position: relative; display: flex; flex-direction: column; align-items: center; gap: 6px; min-width: 0; text-align: center; }
.batch-flow li:not(:first-child)::before {
  content: '';
  position: absolute;
  top: 11px;
  right: calc(50% + 15px);
  width: calc(100% - 30px);
  height: 1.5px;
  border-radius: 2px;
  background: var(--brd-2);
}
.batch-flow li.is-on::before,
.batch-flow li.is-done::before { background: var(--ink); }
.batch-flow .num {
  position: relative;
  z-index: 1;
  display: grid;
  place-items: center;
  width: 23px;
  height: 23px;
  border: 1px solid var(--brd-2);
  border-radius: 50%;
  background: var(--surface);
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 9px;
}
.batch-flow .is-on .num { border-color: var(--ink); background: var(--ink); color: var(--bg); }
.batch-flow .is-done .num { border-color: var(--ink); color: var(--ink); }
.batch-flow b { max-width: 100%; overflow: hidden; color: var(--ink-3); font-size: 10.5px; font-weight: 600; text-overflow: ellipsis; white-space: nowrap; }
.batch-flow .is-on b { color: var(--ink); }
.batch-flow .is-done b { color: var(--ink-2); }
.batch-flow em { color: var(--ink-4); font-family: var(--font-mono); font-size: 9px; font-style: normal; }
.batch-flow .is-on em { color: var(--ink-3); }

.batch-section { margin-top: 24px; padding-top: 22px; border-top: 1px solid var(--brd); }
.batch-section-head { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.batch-section-title { display: flex; align-items: center; gap: 9px; min-width: 0; }
.batch-section-title .idx { color: var(--ink-4); font-family: var(--font-mono); font-size: 9.5px; letter-spacing: 0.08em; }
.batch-section-title h3 { margin: 0; color: var(--ink); font-size: 13px; font-weight: 800; letter-spacing: 0.02em; }
.batch-section-title .tag { padding: 3px 7px; border: 1px solid var(--brd); border-radius: 999px; color: var(--ink-3); font-family: var(--font-mono); font-size: 9px; white-space: nowrap; }
.batch-section-note { color: var(--ink-3); font-size: 10.5px; white-space: nowrap; }
.batch-count { display: inline-flex; align-items: center; gap: 6px; flex: none; height: 24px; padding: 0 9px; border: 1px solid var(--brd); border-radius: 999px; color: var(--ink-3); font-family: var(--font-mono); font-size: 9.5px; }
.batch-count b { color: var(--ink); font-weight: 700; }
.batch-count i { width: 5px; height: 5px; border-radius: 50%; background: var(--ink-3); }

.batch-textbox { margin-top: 12px; border: 1px solid var(--brd); border-radius: 14px; background: var(--surface-2); overflow: hidden; transition: border-color 0.15s ease, box-shadow 0.15s ease; }
.batch-textbox:focus-within { border-color: var(--ink); box-shadow: 0 0 0 3px color-mix(in srgb, var(--ink) 7%, transparent); background: var(--surface); }
.batch-textarea {
  display: block;
  width: 100%;
  height: 146px;
  padding: 15px 16px;
  border: 0;
  outline: 0;
  resize: vertical;
  background: transparent;
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 12.5px;
  line-height: 1.9;
}
.batch-textarea::placeholder { color: var(--ink-4); opacity: 1; }
.batch-textbox-foot { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 9px 14px; border-top: 1px solid var(--brd); background: var(--surface-2); color: var(--ink-3); font-size: 10px; }
.batch-textbox-foot .mono { color: var(--ink-4); font-family: var(--font-mono); font-size: 9.5px; }
.batch-recognized { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 9px; color: var(--ink-3); font-size: 10.5px; }
.batch-recognized .right { color: var(--ink-4); }
.batch-recognized.single { display: block; }
.batch-inline-warn { margin: 6px 0 0; color: var(--warn); font-size: 10px; }

.batch-modes { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 10px; margin-top: 12px; }
.batch-mode {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 11px;
  width: 100%;
  padding: 14px;
  border: 1px solid var(--brd);
  border-radius: 14px;
  background: var(--surface);
  color: var(--ink);
  text-align: left;
  cursor: pointer;
}
.batch-mode.is-on { border-color: var(--ink); box-shadow: 0 0 0 1px var(--ink), 0 6px 18px -14px rgba(0, 0, 0, 0.5); }
.batch-mode-icon { display: grid; place-items: center; flex: none; width: 32px; height: 32px; border-radius: 9px; background: var(--surface-2); color: var(--ink-2); }
.batch-mode.is-on .batch-mode-icon { background: var(--ink); color: var(--bg); }
.batch-mode-icon svg { width: 15px; height: 15px; }
.batch-mode-copy { min-width: 0; }
.batch-mode-title { display: flex; align-items: center; gap: 6px; color: var(--ink); font-size: 12.5px; font-weight: 700; }
.batch-mode-title i { padding: 2px 5px; border-radius: 4px; background: var(--good-soft); color: var(--good); font-family: var(--font-mono); font-size: 8px; font-style: normal; letter-spacing: 0.04em; }
.batch-mode-desc { display: block; margin-top: 4px; color: var(--ink-3); font-size: 10.5px; line-height: 1.5; }
.batch-mode-check { position: absolute; top: 11px; right: 11px; display: grid; place-items: center; width: 18px; height: 18px; border: 1px solid var(--brd-2); border-radius: 50%; color: transparent; }
.batch-mode.is-on .batch-mode-check { border-color: var(--ink); background: var(--ink); color: var(--bg); }
.batch-mode-check svg { width: 10px; height: 10px; }

.batch-import-head { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 20px; }
.batch-import-head span { color: var(--ink); font-size: 11.5px; font-weight: 600; }
.batch-import-head em { color: var(--ink-3); font-size: 9.5px; font-style: normal; }
.batch-format-note { display: flex; align-items: flex-start; gap: 8px; margin-top: 10px; padding: 9px 11px; border-radius: 10px; background: var(--surface-2); color: var(--ink-3); font-size: 10px; line-height: 1.55; }
.batch-format-note svg { width: 12px; height: 12px; margin-top: 1px; flex: none; }
.batch-format-note a { color: var(--ink); text-decoration: none; border-bottom: 1px solid var(--brd-2); white-space: nowrap; }
.batch-label { display: block; margin-bottom: 7px; color: var(--ink); font-size: 11.5px; font-weight: 600; }
.batch-textarea--mailbox { height: 132px; border: 1px solid var(--brd); border-radius: 14px; background: var(--surface-2); }
.batch-textarea--mailbox:focus { border-color: var(--ink); box-shadow: 0 0 0 3px color-mix(in srgb, var(--ink) 7%, transparent); background: var(--surface); }
.batch-file-btn { width: 100%; margin-top: 9px; padding: 10px; border: 1px dashed var(--brd-2); border-radius: 12px; background: var(--surface); color: var(--ink-2); font-size: 10.5px; font-weight: 600; cursor: pointer; }
.batch-file-btn:hover:not(:disabled) { border-color: var(--ink-3); color: var(--ink); }
.batch-file-btn:disabled { opacity: 0.55; cursor: wait; }
.batch-inline-note { margin: 8px 0 0; color: var(--ink-3); font-size: 10px; line-height: 1.55; }
.batch-inline-hint { margin: 4px 0 0; color: var(--ink-4); font-size: 9.5px; line-height: 1.55; }
.batch-mailbox { margin-top: 4px; }

.batch-pair { margin-top: 12px; padding: 14px; border: 1px solid var(--brd); border-radius: 14px; background: var(--surface-2); }
.batch-pair-top { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.batch-pair-cells { display: flex; align-items: center; gap: 16px; }
.batch-pair-cell { display: flex; align-items: baseline; gap: 7px; }
.batch-pair-cell span { color: var(--ink-3); font-size: 10.5px; }
.batch-pair-cell strong { color: var(--ink); font-family: var(--font-mono); font-size: 18px; font-weight: 700; }
.batch-pair-arrow { color: var(--ink-4); font-size: 13px; }
.batch-pair-note { color: var(--ink-3); font-size: 10px; }
.batch-pair-bar { height: 3px; margin-top: 12px; border-radius: 99px; background: var(--brd); overflow: hidden; }
.batch-pair-bar i { display: block; height: 100%; border-radius: 99px; background: var(--ink); }
.batch-submit { width: 100%; height: 46px; margin-top: 16px; font-size: 13px; }
.batch-submit-hint { margin: 9px 0 0; color: var(--ink-3); font-size: 10.5px; text-align: center; }
.batch-inline-error { margin: 9px 0 0; color: var(--err); font-size: 10.5px; line-height: 1.55; }
.batch-section-foot { margin: 14px 0 0; color: var(--ink-3); font-size: 10px; line-height: 1.6; }

/* run */
.batch-run-status { display: flex; align-items: center; gap: 10px; flex: none; }
.batch-run-pill { display: inline-flex; align-items: center; gap: 6px; height: 24px; padding: 0 10px; border: 1px solid color-mix(in srgb, var(--good) 28%, transparent); border-radius: 999px; background: var(--good-soft); color: var(--good); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.08em; }
.batch-run-pill i { width: 5px; height: 5px; border-radius: 50%; background: var(--good); }
.batch-run-count { color: var(--ink-3); font-family: var(--font-mono); font-size: 10px; }
.batch-run-progress { height: 3px; margin-top: 18px; border-radius: 99px; background: var(--brd); overflow: hidden; }
.batch-run-progress i { display: block; height: 100%; border-radius: 99px; background: var(--ink); }

.batch-stats { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 8px; margin-top: 16px; }
.batch-stat { padding: 12px 13px; border: 1px solid var(--brd); border-radius: 12px; background: var(--surface-2); }
.batch-stat span { display: block; color: var(--ink-3); font-size: 9.5px; }
.batch-stat strong { display: block; margin-top: 5px; color: var(--ink); font-family: var(--font-mono); font-size: 21px; font-weight: 700; letter-spacing: -0.03em; }
.batch-stat strong em { color: var(--ink-3); font-size: 12px; font-style: normal; font-weight: 400; }
.batch-stat.is-warn strong { color: var(--warn); }
.batch-stat.is-good strong { color: var(--good); }

.batch-run-import { margin-top: 18px; padding: 14px; border: 1px solid var(--brd); border-radius: 14px; background: var(--surface); }
.batch-run-mailbox { padding-bottom: 10px; }
.batch-verifying { display: flex; align-items: center; justify-content: center; gap: 9px; padding: 34px 0; color: var(--ink-3); font-size: 12px; }
.batch-spinner { width: 14px; height: 14px; border: 2px solid var(--brd); border-top-color: var(--ink); border-radius: 50%; animation: batch-spin 0.9s linear infinite; }
@keyframes batch-spin { to { transform: rotate(360deg); } }

.batch-current { margin-top: 18px; padding: 16px; border: 1px solid color-mix(in srgb, var(--primary) 28%, var(--brd)); border-radius: 14px; background: color-mix(in srgb, var(--primary) 4%, var(--surface)); }
.batch-current-top { display: flex; align-items: flex-start; justify-content: space-between; gap: 14px; }
.batch-current-kicker { display: block; color: var(--ink-3); font-family: var(--font-mono); font-size: 8.5px; letter-spacing: 0.12em; }
.batch-current-title { display: flex; align-items: baseline; gap: 8px; margin-top: 4px; }
.batch-current-title b { color: var(--ink); font-size: 14px; font-weight: 700; }
.batch-current-title span { color: var(--ink-3); font-size: 10px; }
.batch-current-key { display: flex; align-items: center; gap: 8px; flex: none; }
.batch-current-key .mono { color: var(--ink); font-family: var(--font-mono); font-size: 10.5px; }
.batch-current-key .plan { padding: 3px 6px; border: 1px solid var(--brd); border-radius: 4px; color: var(--ink-3); font-size: 8.5px; }
.batch-account-line { margin: 10px 0 0; color: var(--ink-3); font-size: 10px; }
.batch-account-line b { color: var(--ink-2); font-family: var(--font-mono); font-size: 10px; font-weight: 500; }
.batch-textarea--session { min-height: 86px; margin-top: 9px; padding: 12px; border: 1px solid var(--brd); border-radius: 11px; background: var(--surface); font-size: 10.5px; line-height: 1.7; }
.batch-textarea--session:focus { border-color: var(--ink); box-shadow: 0 0 0 3px color-mix(in srgb, var(--ink) 7%, transparent); }
.batch-current-actions { display: flex; gap: 8px; margin-top: 10px; }
.batch-current-actions .btn-primary { flex: 1; height: 40px; font-size: 11.5px; }
.batch-auto { background: var(--good) !important; border-color: var(--good) !important; }
.batch-current-hint { margin: 9px 0 0; color: var(--ink-3); font-size: 9.5px; }
.batch-collected { margin-top: 18px; padding: 16px; border: 1px solid color-mix(in srgb, var(--good) 30%, var(--brd)); border-radius: 14px; background: var(--good-soft); text-align: center; }
.batch-collected p { margin: 0; color: var(--good); font-size: 12px; font-weight: 600; }
.batch-collected span { display: block; margin-top: 4px; color: var(--ink-3); font-size: 10px; }
.batch-export-link { margin-top: 9px; padding: 0; border: 0; background: transparent; color: var(--good); font-size: 10.5px; text-decoration: underline; cursor: pointer; }
.batch-empty { margin: 18px 0; color: var(--ink-3); font-size: 11px; text-align: center; }
.batch-empty.table { margin: 12px 0; }

.batch-progress { margin-top: 22px; padding-top: 18px; border-top: 1px solid var(--brd); }
.batch-progress-head { display: flex; align-items: center; justify-content: space-between; gap: 14px; }
.batch-progress-head h3 { margin: 0; color: var(--ink); font-size: 12px; font-weight: 800; }
.batch-filters { display: flex; align-items: center; gap: 4px; }
.batch-filters button { padding: 4px 8px; border: 0; border-radius: 999px; background: transparent; color: var(--ink-3); font: inherit; font-size: 9px; cursor: pointer; }
.batch-filters button.on { background: var(--surface-2); color: var(--ink); font-weight: 600; }
.batch-table-wrap { margin-top: 10px; overflow-x: auto; border: 1px solid var(--brd); border-radius: 12px; }
.batch-table { min-width: 660px; }
.batch-trow { display: grid; grid-template-columns: 24px minmax(128px, 1.15fr) 66px 68px minmax(140px, 1fr) minmax(110px, 1fr); gap: 8px; align-items: center; padding: 9px 11px; border-top: 1px solid var(--brd); font-size: 10.5px; }
.batch-thead { border-top: 0; background: var(--surface-2); color: var(--ink-3); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.06em; }
.batch-trow .idx { color: var(--ink-4); font-family: var(--font-mono); font-size: 9px; }
.batch-trow .key { overflow: hidden; color: var(--ink); font-family: var(--font-mono); font-size: 10px; text-overflow: ellipsis; white-space: nowrap; }
.batch-trow .plan { color: var(--ink-3); font-size: 10px; }
.batch-trow .account { display: flex; flex-direction: column; overflow: hidden; color: var(--ink-2); font-size: 10px; line-height: 1.45; }
.batch-trow .account em { color: var(--ink-3); font-family: var(--font-mono); font-size: 9px; font-style: normal; }
.batch-trow .msg { overflow: hidden; color: var(--ink-3); font-size: 10px; text-overflow: ellipsis; white-space: nowrap; }
.batch-chip { justify-self: start; padding: 3px 6px; border: 1px solid transparent; border-radius: 999px; font-size: 9px; font-weight: 600; white-space: nowrap; }
.batch-chip.is-good { color: var(--good); background: var(--good-soft); border-color: color-mix(in srgb, var(--good) 24%, transparent); }
.batch-chip.is-warn { color: var(--warn); background: var(--warn-soft); border-color: color-mix(in srgb, var(--warn) 24%, transparent); }
.batch-chip.is-bad { color: var(--err); background: var(--err-soft); border-color: color-mix(in srgb, var(--err) 24%, transparent); }
.batch-chip.is-info { color: var(--primary); background: var(--primary-soft); border-color: color-mix(in srgb, var(--primary) 20%, transparent); }
.batch-chip.is-muted { color: var(--ink-3); background: var(--surface-2); border-color: var(--brd); }
.batch-run-actions { display: flex; gap: 8px; margin-top: 14px; }
.batch-run-actions .btn-secondary { flex: 1; height: 38px; font-size: 11px; }

/* aside */
.batch-aside { display: flex; flex-direction: column; gap: 16px; min-width: 0; }
.batch-overview {
  position: relative;
  overflow: hidden;
  padding: 22px 22px 18px;
  border-radius: var(--radius-lg);
  background:
    radial-gradient(340px 190px at 18% -20%, color-mix(in srgb, var(--bg) 13%, transparent), transparent 62%),
    radial-gradient(280px 220px at 108% 112%, color-mix(in srgb, var(--good) 17%, transparent), transparent 64%),
    var(--ink);
  color: var(--bg);
  box-shadow: 0 28px 56px -42px rgba(0, 0, 0, 0.65);
}
.batch-overview-top { display: flex; align-items: center; justify-content: space-between; gap: 12px; color: color-mix(in srgb, var(--bg) 52%, transparent); font-family: var(--font-mono); font-size: 8.5px; letter-spacing: 0.16em; }
.batch-overview-main { display: flex; align-items: flex-end; gap: 10px; margin-top: 16px; }
.batch-overview-main strong { color: var(--bg); font-family: var(--font-mono); font-size: 44px; font-weight: 700; line-height: 1; letter-spacing: -0.06em; }
.batch-overview-main span { padding-bottom: 6px; color: color-mix(in srgb, var(--bg) 55%, transparent); font-size: 10.5px; }
.batch-overview-bar { height: 3px; margin-top: 16px; border-radius: 99px; background: color-mix(in srgb, var(--bg) 14%, transparent); overflow: hidden; }
.batch-overview-bar i { display: block; height: 100%; border-radius: 99px; background: var(--bg); }
.batch-overview-rows { display: flex; flex-direction: column; gap: 11px; margin-top: 18px; padding-top: 16px; border-top: 1px solid color-mix(in srgb, var(--bg) 10%, transparent); }
.batch-overview-rows > div { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.batch-overview-rows span { color: color-mix(in srgb, var(--bg) 48%, transparent); font-size: 10px; }
.batch-overview-rows b { color: color-mix(in srgb, var(--bg) 88%, transparent); font-family: var(--font-mono); font-size: 10.5px; font-weight: 500; }
.batch-overview-foot { display: flex; align-items: center; gap: 6px; margin-top: 17px; color: color-mix(in srgb, var(--bg) 44%, transparent); font-size: 9.5px; }
.batch-overview-foot i { width: 5px; height: 5px; border-radius: 50%; background: color-mix(in srgb, var(--bg) 55%, transparent); }

.batch-side-card { padding: 20px 22px 8px; }
.batch-side-card h3 { margin: 0 0 3px; color: var(--ink); font-size: 11.5px; font-weight: 800; letter-spacing: 0.02em; }
.batch-guide { margin: 0; padding: 0; list-style: none; }
.batch-guide li { display: flex; align-items: flex-start; gap: 11px; padding: 12px 0; border-top: 1px solid var(--brd); }
.batch-guide li:first-child { border-top: 0; }
.batch-guide li > span { padding-top: 1px; color: var(--ink-3); font-family: var(--font-mono); font-size: 9px; letter-spacing: 0.06em; }
.batch-guide li.on > span { color: var(--ink); }
.batch-guide b { display: block; color: var(--ink); font-size: 11.5px; font-weight: 700; }
.batch-guide p { margin: 2px 0 0; color: var(--ink-3); font-size: 10px; line-height: 1.5; }
.batch-guide li:not(.on) b { color: var(--ink-2); font-weight: 600; }
.batch-formats { margin: 0; padding: 0; list-style: none; }
.batch-formats li { display: flex; align-items: flex-start; gap: 9px; padding: 9px 0; border-top: 1px solid var(--brd); }
.batch-formats li:first-child { border-top: 0; }
.batch-formats li > span { flex: none; min-width: 50px; padding: 3px 6px; border: 1px solid var(--brd); border-radius: 5px; color: var(--ink-2); font-family: var(--font-mono); font-size: 8px; text-align: center; letter-spacing: 0.03em; }
.batch-formats p { margin: 0; color: var(--ink-3); font-size: 10px; line-height: 1.55; }
.batch-lock { display: flex; align-items: flex-start; gap: 7px; margin: 10px 0 0; padding-top: 12px; border-top: 1px solid var(--brd); color: var(--ink-3); font-size: 9.5px; line-height: 1.55; }
.batch-lock svg { width: 11px; height: 11px; margin-top: 1px; flex: none; }
.batch-info { margin: 0; padding: 0; list-style: none; }
.batch-info li { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 9px 0; border-top: 1px solid var(--brd); }
.batch-info li:first-child { border-top: 0; }
.batch-info span { color: var(--ink-3); font-size: 10px; }
.batch-info b { color: var(--ink-2); font-family: var(--font-mono); font-size: 10px; font-weight: 500; text-align: right; word-break: break-all; }
.batch-side-actions { display: flex; flex-direction: column; gap: 7px; padding: 8px 0 12px; }
.batch-side-actions button { display: flex; align-items: center; justify-content: space-between; width: 100%; padding: 10px 11px; border: 1px solid var(--brd); border-radius: 10px; background: transparent; color: var(--ink); font: inherit; font-size: 10.5px; cursor: pointer; }
.batch-side-actions button:hover:not(:disabled) { border-color: var(--ink-3); }
.batch-side-actions button:disabled { opacity: 0.45; cursor: not-allowed; }
.batch-side-actions button span { color: var(--ink-3); font-family: var(--font-mono); font-size: 9px; }
.batch-note { display: flex; align-items: flex-start; gap: 8px; padding: 2px 4px; color: var(--ink-3); font-size: 10.5px; line-height: 1.65; }
.batch-note svg { width: 12px; height: 12px; margin-top: 1px; flex: none; }

@media (max-width: 900px) {
  .batch-wrap { max-width: calc(100% - 32px); }
  .batch-hero { flex-direction: column; align-items: stretch; gap: 20px; padding: 26px 0 24px; }
  .batch-hero h1 { font-size: 42px; }
  .batch-hero-actions { padding: 0; }
  .batch-grid { grid-template-columns: minmax(0, 1fr); }
  .batch-aside { order: 2; }
  .batch-card { padding: 22px 20px 20px; }
  .batch-flow b { font-size: 10px; }
  .batch-stats { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

@media (min-width: 641px) and (max-width: 900px) {
  .batch-aside { max-width: 460px; }
}

@media (max-width: 640px) {
  .batch-wrap { padding: 0 16px; }
  .batch-subnav { padding: 4px 0 14px; overflow: hidden; }
  .batch-crumb { display: none; }
  .batch-hero { gap: 18px; padding: 22px 0; }
  .batch-hero h1 { margin-top: 10px; font-size: 38px; }
  .batch-hero p { margin-top: 10px; font-size: 12.5px; line-height: 1.7; }
  .batch-eyebrow { font-size: 9.5px; }
  .batch-action { width: 100%; height: 40px; }
  .batch-card { padding: 18px 16px 16px; border-radius: 18px; }
  .batch-panel-head { flex-wrap: wrap; gap: 10px; }
  .batch-panel-head > div:first-child { flex: 1 1 100%; min-width: 0; }
  .batch-panel-head h2 { font-size: 19px; }
  .batch-panel-head p { font-size: 11px; }
  .batch-state { margin-left: auto; }
  .batch-run-status { width: 100%; justify-content: space-between; }
  .batch-flow { margin-top: 18px; }
  .batch-flow b { font-size: 9.5px; }
  .batch-flow em { display: none; }
  .batch-section { margin-top: 20px; padding-top: 16px; }
  .batch-section-title h3 { font-size: 12px; }
  .batch-section-note { font-size: 9.5px; }
  .batch-textarea { height: 132px; font-size: 11px; }
  .batch-modes { grid-template-columns: 1fr; }
  .batch-mode { padding: 12px; }
  .batch-pair-top { flex-direction: column; align-items: flex-start; gap: 7px; }
  .batch-pair-cell strong { font-size: 16px; }
  .batch-submit { height: 42px; }
  .batch-stats { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .batch-current-top { flex-direction: column; gap: 8px; }
  .batch-current-actions { flex-direction: column; }
  .batch-progress-head { flex-direction: column; align-items: flex-start; gap: 8px; }
  .batch-filters { width: 100%; overflow-x: auto; }
  .batch-run-actions { flex-direction: column; }
  .batch-table { min-width: 0; }
  .batch-thead { display: none; }
  .batch-trow { grid-template-columns: 22px minmax(0, 1fr) auto; gap: 2px 8px; padding: 10px 11px; }
  .batch-trow .idx { grid-column: 1; grid-row: 1; }
  .batch-trow .key { grid-column: 2; grid-row: 1; }
  .batch-trow .batch-chip { grid-column: 3; grid-row: 1; justify-self: end; }
  .batch-trow .plan { grid-column: 2; grid-row: 2; }
  .batch-trow .account { grid-column: 2 / 4; grid-row: 3; }
  .batch-trow .msg { grid-column: 2 / 4; grid-row: 4; }
  .batch-overview { padding: 20px; border-radius: 18px; }
  .batch-overview-main strong { font-size: 38px; }
  .batch-side-card { padding: 18px 18px 6px; border-radius: 18px; }
}
</style>
