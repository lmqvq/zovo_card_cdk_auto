<template>
  <div class="u-page">
    <UserPageHeader v-if="step !== 1 && step !== 4" :eyebrow="t('hero.redeemEyebrow')" :title="t('hero.redeemTitle')" :subtitle="t('hero.redeemSub')" badges />
    <UserPageHeader v-else hide-hero :title="t('hero.redeemTitle')" />
    <div class="u-main" :class="(step === 1 || step === 4) ? 'result-wrap' : 'max-w-3xl'">
      <template v-if="step === 2 || step === 3">
        <RedeemModeTabs />
        <!-- redeem-flow v2: no public fee reference -->
        <!-- steps -->
        <ol class="stepper" :aria-label="t('hero.redeemTitle')">
          <li v-for="(s, i) in steps" :key="s" :class="{ 'is-done': step > i + 1, 'is-current': step === i + 1 }" :aria-current="step === i + 1 ? 'step' : undefined">
            <span class="stepper-dot">{{ i + 1 }}</span>
            <span class="stepper-label">{{ s }}</span>
          </li>
        </ol>
      </template>
      <template v-else-if="step === 1">
        <div class="result-subnav">
          <div class="result-crumb">{{ t('hero.redeemEyebrow') }} <i>/</i> <b>输入卡密</b> <i>/</i> STEP 01</div>
          <RedeemModeTabs compact />
        </div>
        <section class="result-hero cdk-hero">
          <div class="result-hero-copy">
            <div class="result-eyebrow">REDEEM <em>· 输入卡密</em></div>
            <h1>{{ t('hero.redeemTitle') }}</h1>
            <p>{{ t('hero.redeemSub') }}</p>
            <ul class="cdk-badges">
              <li>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" /><path d="m9 12 2 2 4-4" /></svg>
                {{ t('hero.badgeOfficial') }}
              </li>
              <li>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 6 9 17l-5-5" /></svg>
                {{ t('hero.badgeSelf') }}
              </li>
              <li>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></svg>
                {{ t('hero.badgeTrack') }}
              </li>
            </ul>
          </div>
          <div class="cdk-hero-aside">
            <span>已经兑换过？</span>
            <router-link to="/history">
              查询卡密状态
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M7 17 17 7" /><path d="M8 7h9v9" /></svg>
            </router-link>
          </div>
        </section>
      </template>
      <template v-else>
        <div class="result-subnav">
          <div class="result-crumb">卡密兑换 <i>/</i> <b>兑换结果</b> <i>/</i> STEP 04</div>
          <RedeemModeTabs compact />
        </div>
        <section class="result-hero">
          <div class="result-hero-copy">
            <div class="result-eyebrow">REDEEM RESULT <em>· {{ resultToneText }}</em></div>
            <h1>{{ resultTitle }}</h1>
            <p v-if="resultSubtitle">{{ resultSubtitle }}</p>
          </div>
          <div class="result-actions">
            <router-link class="btn-secondary result-action" to="/billing">查询账单</router-link>
            <button class="btn-primary result-action" @click="resetAll">
              再兑一张
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12h14" /><path d="m13 6 6 6-6 6" /></svg>
            </button>
          </div>
        </section>
      </template>

      <!-- 1 preview -->
      <div v-show="step === 1" class="cdk-grid">
        <div class="cdk-col">
          <section class="card cdk-form-card">
            <div class="cdk-form-head">
              <div>
                <span class="cdk-micro">STEP 01 · 卡密校验</span>
                <h2>输入 CDK</h2>
              </div>
              <span class="cdk-form-badge"><i></i>仅校验卡密</span>
            </div>

            <div class="cdk-field">
              <div class="cdk-field-label"><span>CDK 卡密</span><span>必填</span></div>
              <div class="cdk-input-shell">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 9a3 3 0 0 0 0 6v3h18v-3a3 3 0 0 1 0-6V6H3z" /><path d="M13 6v2" /><path d="M13 11v2" /><path d="M13 16v2" /></svg>
                <input
                  v-model="code"
                  class="cdk-input"
                  type="text"
                  autocomplete="off"
                  autocapitalize="characters"
                  spellcheck="false"
                  placeholder="ZC-XXXX-XXXX-XXXX-XXXX"
                  aria-label="CDK 卡密"
                  @keyup.enter="doPreview"
                />
              </div>
              <div class="cdk-hint">请粘贴完整卡密，包含所有分段。</div>
            </div>

            <div v-if="error" class="alert alert-error cdk-error">{{ error }}</div>
            <div v-if="previewInfo" class="cdk-preview">
              <div class="cdk-preview-row"><span>套餐</span><b>{{ targetPlan ? planLabel(targetPlan) : '—' }}</b></div>
              <p>{{ previewInfo.message || '码有效，可继续' }}</p>
            </div>

            <button class="btn-primary cdk-cta" type="button" :disabled="busy" @click="doPreview">
              {{ busy ? '校验中…' : '预览 / 下一步' }}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M5 12h14" /><path d="m13 6 6 6-6 6" /></svg>
            </button>
            <p class="cdk-note">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
              点击预览仅校验卡密有效性，不会扣款或提交账号。
            </p>
          </section>

          <section class="card cdk-tips-card">
            <h3>兑换前确认</h3>
            <ol class="cdk-tips-list">
              <li><span class="cdk-tips-num">01</span><div><b>卡密未使用</b><p>同一张卡密仅能兑换一次。</p></div></li>
              <li><span class="cdk-tips-num">02</span><div><b>账号可登录</b><p>后续按页面提示填写登录态或邮箱密码。</p></div></li>
              <li><span class="cdk-tips-num">03</span><div><b>套餐未满足</b><p>账号已有更高套餐时，可能无法重复购买。</p></div></li>
            </ol>
          </section>
        </div>

        <aside class="cdk-col cdk-aside">
          <section class="cdk-flow-card">
            <div class="cdk-flow-top"><span>REDEEM FLOW</span><span>4 STEPS</span></div>
            <h3>兑换流程</h3>
            <p class="cdk-flow-sub">从卡密到开通，四步完成。</p>
            <ol class="cdk-flow-steps">
              <li class="is-active">
                <span class="cdk-flow-num">01</span>
                <div><b>{{ steps[0] }}</b><p>输入卡密并校验套餐信息</p></div>
                <em>当前</em>
              </li>
              <li>
                <span class="cdk-flow-num">02</span>
                <div><b>{{ steps[1] }}</b><p>按页面提示填写登录态</p></div>
              </li>
              <li>
                <span class="cdk-flow-num">03</span>
                <div><b>{{ steps[2] }}</b><p>核对账号与目标套餐</p></div>
              </li>
              <li>
                <span class="cdk-flow-num">04</span>
                <div><b>{{ steps[3] }}</b><p>自动处理，凭证可查</p></div>
              </li>
            </ol>
          </section>

          <section class="cdk-security">
            <h3>安全提示</h3>
            <p>卡密仅用于本次校验与兑换；兑换结果绑定当前设备，请勿在公共设备保存登录态。</p>
          </section>
        </aside>
      </div>

      <!-- 2 preflight -->
      <div v-show="step === 2" class="card space-y-4">
        <h2 class="text-xl font-bold text-ink">{{ isGrok ? t('grok.credential') : isXPremiumPlan(targetPlan) ? t('xPremium.credential') : 'ChatGPT 凭证' }}</h2>
        <div v-if="!isXPremiumPlan(targetPlan) && !isGrok" class="flex gap-2">
          <button type="button" class="btn-secondary !py-1" :class="{ 'ring-2': credMode === 'session' }" @click="credMode = 'session'">Session</button>
        </div>
        <template v-if="credMode === 'session'">
          <p v-if="isGrok" class="text-sm text-muted">{{ t('grok.hint') }}</p>
          <p v-else-if="isXPremiumPlan(targetPlan)" class="text-sm text-muted">{{ t('xPremium.hint') }}</p>
          <p v-else class="text-sm text-muted">打开
            <a class="app-link" href="https://chatgpt.com/api/auth/session" target="_blank" rel="noopener">chatgpt.com/api/auth/session</a>
            复制<strong>完整 JSON</strong>（必须含 <code>sessionToken</code>）。已禁用纯 Access Token。
          </p>
          <input
            ref="sessionFileInput"
            class="sr-only"
            type="file"
            :accept="SESSION_FILE_ACCEPT"
            @change="onSessionFileChange"
          />
          <div
            class="session-dropzone"
            :class="{ 'is-dragging': isSessionDragActive }"
            @dragenter.prevent="onSessionDragEnter"
            @dragover.prevent="onSessionDragOver"
            @dragleave.prevent="onSessionDragLeave"
            @drop.prevent="onSessionFileDrop"
          >
            <textarea
              v-model="sessionRaw"
              class="input h-36 font-mono text-xs"
              :placeholder="isGrok ? t('grok.placeholder') : isXPremiumPlan(targetPlan) ? t('xPremium.placeholder') : 'Session JSON (sessionToken)'"
              @input="sessionFileName = ''"
            />
            <div v-if="isSessionDragActive" class="session-drop-overlay" aria-hidden="true">
              {{ t('recharge.fileDropActive') }}
            </div>
          </div>
          <div class="session-file-actions">
            <button
              type="button"
              class="btn-secondary session-file-button"
              :disabled="busy || readingSessionFile"
              @click="openSessionFilePicker"
            >
              <el-icon><Document /></el-icon>
              {{ readingSessionFile ? t('recharge.fileReading') : t('recharge.fileButton') }}
            </button>
            <span v-if="sessionFileName" class="text-xs text-muted break-all">
              {{ t('recharge.fileLoaded', { name: sessionFileName }) }}
            </span>
          </div>
          <p class="text-xs text-muted">{{ t('recharge.fileHint') }}</p>
          <p class="text-xs text-subtle">{{ t('recharge.filePrivacy') }}</p>
        </template>
        <template v-else>
          <input v-model="email" class="input" placeholder="email@outlook.com" />
          <input v-model="password" class="input" type="password" placeholder="邮箱密码" />
        </template>
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div class="flex gap-3">
          <button class="btn-secondary flex-1" @click="step = 1">上一步</button>
          <button class="btn-primary flex-1" :disabled="busy || readingSessionFile" @click="doPreflight">{{ busy ? '预检中…' : '预检' }}</button>
        </div>
      </div>

      <!-- 3 redeem -->
      <div v-show="step === 3" class="card space-y-4">
        <h2 class="text-xl font-bold text-ink">确认兑换</h2>
        <p class="text-sm text-muted">将提交兑换请求；结果不确定或 review 时请勿重复提交，请轮询结果。</p>

        <!-- 与卡台直充预检一致：邮箱 + 订阅事实 -->
        <div v-if="account.checked" class="rounded-xl border p-4 space-y-3" style="border-color: var(--brd); background: var(--surface-2, var(--soft))">
          <div class="flex flex-wrap items-center gap-2">
            <strong class="text-ink text-base">{{ account.email || '账号已验证' }}</strong>
            <el-tag size="small" :type="subscriptionStatusTag">{{ subscriptionStatusText }}</el-tag>
          </div>
          <dl class="account-facts">
            <div>
              <dt>目标套餐</dt>
              <dd>{{ planLabel(targetPlan) }}</dd>
            </div>
            <div>
              <dt>当前套餐</dt>
              <dd>{{ planLabel(isGrok && account.currentPlan ? 'grok_' + String(account.currentPlan).replace(/^grok_/i, '') : account.currentPlan) }}</dd>
            </div>
            <div v-if="isGrok && account.xSubscriptionType">
              <dt>X 订阅等级</dt>
              <dd>{{ account.xSubscriptionType }}</dd>
            </div>
            <div>
              <dt>订阅状态</dt>
              <dd>{{ subscriptionStatusText }}</dd>
            </div>
            <div>
              <dt>套餐到期</dt>
              <dd>{{ account.subscriptionActiveUntil ? fmtTime(account.subscriptionActiveUntil) : '上游未提供' }}</dd>
            </div>
            <div>
              <dt>剩余</dt>
              <dd>{{ account.subscriptionActiveUntil ? subscriptionRemaining : '上游未提供' }}</dd>
            </div>
            <div>
              <dt>自动续费</dt>
              <dd :class="account.subscriptionWillRenew === true ? 'text-warn' : account.subscriptionWillRenew === false ? 'text-good' : ''">
                {{ renewalStatusText }}
              </dd>
            </div>
            <div>
              <dt>最近付款时间</dt>
              <dd>{{ lastPaymentAtText }}</dd>
            </div>
            <div>
              <dt>最近付款</dt>
              <dd>{{ lastPaymentText }}</dd>
            </div>
            <div>
              <dt>支付方式</dt>
              <dd>{{ paymentMethodText }}</dd>
            </div>
          </dl>
        </div>
        <div v-else class="rounded-xl bg-soft p-3 text-sm text-muted">
          预检未返回账号摘要，仍可尝试兑换（以卡台校验为准）。
        </div>

        <div v-if="grokBlockingReason" class="alert" style="background: var(--warn-soft, #fef3c7); color: var(--warn, #b45309); border-color: var(--warn, #d97706)">
          {{ grokBlockingReason }}
        </div>

        <div v-if="alreadySatisfied && !grokBlockingReason && !needsSubscriptionRecovery && !recoveryPending" class="alert" style="background: var(--warn-soft, #fef3c7); color: var(--warn, #b45309); border-color: var(--warn, #d97706)">
          {{ alreadySatisfiedHint }}
        </div>

        <div v-if="needsSubscriptionRecovery || recoveryPending" class="alert" style="background: var(--warn-soft, #fef3c7); color: var(--warn, #b45309); border-color: var(--warn, #d97706)">
          <p>{{ t(recoveryPending ? 'grace.pending' : 'grace.hint') }}</p>
          <button v-if="needsSubscriptionRecovery && !recoveryPending" class="btn-secondary mt-3" :disabled="busy || recoveringSubscription || !preflightToken" @click="recoverGraceSubscription">
            {{ t(recoveringSubscription ? 'grace.processing' : 'grace.action') }}
          </button>
          <button v-else class="btn-secondary mt-3" :disabled="busy || recoveringSubscription" @click="doPreflight">{{ t('grace.recheck') }}</button>
        </div>

        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div class="flex gap-3">
          <button class="btn-secondary flex-1" @click="step = 2">上一步</button>
          <button class="btn-primary flex-1" :disabled="busy || recoveringSubscription || recoveryPending || needsSubscriptionRecovery || alreadySatisfied || !canRedeem" @click="doRedeem">
            {{ busy ? '提交中…' : alreadySatisfied ? '当前套餐已满足' : '兑换' }}
          </button>
        </div>
      </div>

      <!-- 4 result -->
      <div v-show="step === 4" class="result-grid">
        <section class="card result-card">
          <header class="result-head">
            <div class="result-tone" :class="'is-' + resultTone">
              <svg v-if="resultTone === 'good'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m4 12.5 5 5L20 6.5" /></svg>
              <svg v-else-if="resultTone === 'bad'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 6l12 12" /><path d="M18 6 6 18" /></svg>
              <svg v-else-if="resultTone === 'warn'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 8v5" /><path d="M12 17h.01" /><path d="M10.3 3.9 2.5 17.4A2 2 0 0 0 4.2 20h15.6a2 2 0 0 0 1.7-2.6L13.7 3.9a2 2 0 0 0-3.4 0z" /></svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></svg>
            </div>
            <div class="result-head-copy">
              <h2>{{ resultHeadTitle }}</h2>
              <p>
                <span class="result-email">{{ displayResultEmail || '—' }}</span>
                <template v-if="resultLastAt"> · {{ resultHeadVerb }} {{ resultCompletedText }}</template>
              </p>
            </div>
            <div class="result-head-meta">
              <span class="result-pill" :class="'is-' + resultTone"><i></i>{{ statusBadge }}</span>
              <small>{{ resultMetaText }}</small>
            </div>
          </header>

          <div v-if="error" class="alert alert-error result-alert">{{ error }}</div>
          <div v-if="isTerminal(resultStatus) && resultStatus !== 'completed'" class="alert alert-error result-alert">
            兑换未成功。若状态为 review / pending 请勿重复提交；可联系发码方或稍后用同一设备再查结果。
          </div>

          <ol class="result-rail" aria-label="处理阶段">
            <li v-for="(s, i) in resultStages" :key="s.key" class="result-stage" :class="'is-' + s.state">
              <span class="result-dot">
                <svg v-if="s.state === 'done'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m5 12 4.5 4.5L19 7" /></svg>
                <svg v-else-if="s.state === 'bad'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 6l12 12" /><path d="M18 6 6 18" /></svg>
                <svg v-else-if="s.state === 'warn'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 8v5" /><path d="M12 17h.01" /></svg>
                <span v-else-if="s.state === 'active'" class="result-pulse"></span>
                <template v-else>{{ i + 1 }}</template>
              </span>
              <b>{{ s.label }}</b>
              <time v-if="s.time">{{ s.time }}</time>
            </li>
          </ol>

          <section class="result-section">
            <div class="result-section-head">
              <h3>兑换详情</h3>
              <span>仅本设备可见</span>
            </div>
            <dl class="result-facts">
              <div class="result-fact result-fact--email">
                <dt>兑换账号</dt>
                <dd class="is-mono">
                  {{ displayResultEmail || '—' }}
                  <button v-if="displayResultEmail" type="button" class="result-copy" aria-label="复制账号" @click="copyText(displayResultEmail, '账号已复制')">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="9" y="9" width="11" height="11" rx="2" /><path d="M5 15V5a2 2 0 0 1 2-2h10" /></svg>
                  </button>
                </dd>
              </div>
              <div class="result-fact">
                <dt>开通产品</dt>
                <dd>{{ resultPlanText }}</dd>
              </div>
              <div class="result-fact">
                <dt>{{ resultTimeLabel }}</dt>
                <dd class="is-mono">{{ resultCompletedText }}</dd>
              </div>
              <div class="result-fact result-fact--code">
                <dt>CDK 卡密</dt>
                <dd class="is-mono">
                  {{ resultCdkCode || '—' }}
                  <button v-if="resultCdkCode" type="button" class="result-copy" aria-label="复制 CDK" @click="copyText(resultCdkCode, 'CDK 已复制')">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="9" y="9" width="11" height="11" rx="2" /><path d="M5 15V5a2 2 0 0 1 2-2h10" /></svg>
                  </button>
                </dd>
              </div>
            </dl>
          </section>

          <section class="result-section">
            <div class="result-section-head">
              <h3>处理明细</h3>
              <button v-if="timeline.length" type="button" class="result-link" @click="showAllEvents = !showAllEvents">
                {{ showAllEvents ? '收起全部事件' : `展开全部 ${timeline.length} 条事件` }}
              </button>
            </div>
            <ol v-if="resultRows.length" class="result-events">
              <li v-for="row in resultRows" :key="row.key" class="result-event" :class="'is-' + row.state">
                <time>{{ row.time || '—' }}</time>
                <span class="result-node" aria-hidden="true"></span>
                <div class="result-event-body">
                  <div class="result-event-top">
                    <b>{{ row.label }}</b>
                    <span class="result-chip" :class="'is-' + row.state">{{ row.chip }}</span>
                    <span v-if="row.duration" class="result-dur">{{ row.duration }}</span>
                  </div>
                  <p>{{ row.message }}</p>
                </div>
              </li>
            </ol>
            <p v-else class="result-empty">{{ polling ? '已提交，等待卡台返回步骤明细…' : '暂无处理明细，可稍后重新查询。' }}</p>

            <ol v-if="showAllEvents && timeline.length" class="result-raw">
              <li v-for="(ev, idx) in timeline" :key="ev.id || idx">
                <span class="result-raw-time">{{ fmtShortTime(ev.created_at) }}</span>
                <b>{{ stepLabel(ev.step) }}</b>
                <span class="result-chip" :class="categoryChipClass(ev.category)">{{ categoryLabel(ev.category) }}</span>
                <span class="result-raw-msg">{{ ev.public_message || ev.to_status || '—' }}</span>
              </li>
            </ol>
          </section>

          <div class="result-foot">
            <span>{{ resultFootText }}</span>
            <button v-if="!polling" type="button" class="result-link" @click="startPoll">重新查询结果</button>
            <span v-else class="result-foot-live"><i></i>自动刷新中</span>
          </div>
        </section>

        <aside class="result-aside">
          <section class="result-receipt">
            <div class="receipt-top">
              <span>REDEEM RECEIPT</span>
              <span>CDK {{ resultCdkCode || '—' }}</span>
            </div>
            <div class="receipt-status">
              <span class="receipt-ring" :class="'is-' + resultTone">
                <svg v-if="resultTone === 'good'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m4 12.5 5 5L20 6.5" /></svg>
                <svg v-else-if="resultTone === 'bad'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 6l12 12" /><path d="M18 6 6 18" /></svg>
                <svg v-else-if="resultTone === 'warn'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 8v5" /><path d="M12 17h.01" /></svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></svg>
              </span>
              <div>
                <b>{{ receiptStatusText }}</b>
                <span>{{ statusBadge }} · {{ receiptDate }}</span>
              </div>
            </div>
            <div class="receipt-product">{{ resultPlanText }}</div>
            <div class="receipt-sub">本次兑换凭证</div>
            <div class="receipt-perf" aria-hidden="true"></div>
            <dl class="receipt-rows">
              <div class="receipt-row"><dt>ACCOUNT</dt><dd>{{ displayResultEmail || '—' }}</dd></div>
              <div class="receipt-row"><dt>PAID</dt><dd>{{ paymentTimeText }}</dd></div>
            </dl>
            <button type="button" class="receipt-btn" @click="copyReceipt">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true"><path d="M8 8V5a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-3" /><rect x="3" y="8" width="13" height="13" rx="2" /></svg>
              复制凭证信息
            </button>
            <router-link class="receipt-link" to="/billing">查询这笔账单 ↗</router-link>
            <div class="receipt-foot">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><rect x="4" y="10" width="16" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
              仅本设备可查看，请勿公开分享
            </div>
          </section>

          <section class="result-next">
            <h3>接下来</h3>
            <ol>
              <li>
                <span class="result-next-num">01</span>
                <div><b>登录 ChatGPT 确认兑换结果</b><p>权益或支付方式通常立即生效，未显示请重新登录。</p></div>
              </li>
              <li>
                <span class="result-next-num">02</span>
                <div><b>保存本次兑换凭证</b><p>复制订单信息，售后核对时更快。</p></div>
                <button type="button" class="result-mini" @click="copyReceipt">复制</button>
              </li>
              <li>
                <span class="result-next-num">03</span>
                <div><b>超过 10 分钟未生效</b><p>凭 CDK 卡密联系客服，无需重复提交。</p></div>
              </li>
            </ol>
          </section>

          <p class="result-device-note">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M12 8v.5" /><path d="M12 11v5" /></svg>
            兑换结果与本设备绑定。更换设备或清理浏览器数据后，可在「卡密查询」中重新检索。
          </p>
        </aside>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { isXPremiumPlan, xPremiumCredential } from '../../lib/x-premium'
import { grokCredential, isGrokPlan } from '../../lib/grok'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import UserPageHeader from '../../components/UserPageHeader.vue'
import RedeemModeTabs from '../../components/RedeemModeTabs.vue'
import { dialog } from '../../lib/dialog'
import { copyToClipboard } from '../../lib/clipboard'
import { planLabel, planSatisfied as isSatisfied } from '../../lib/plan'
import {
  SESSION_FILE_ACCEPT,
  SessionFileError,
  readSessionFile,
} from '../../lib/session-file'

const { t } = useI18n({ useScope: 'global' })
const route = useRoute()
const steps = computed(() => [t('hero.steps.preview'), t('hero.steps.cred'), t('hero.steps.redeem'), t('hero.steps.result')])
const step = ref(1)
const busy = ref(false)
const recoveringSubscription = ref(false)
const recoveryPending = ref(false)
let preflightSequence = 0
const error = ref('')
const code = ref('')
const previewInfo = ref<any>(null)
const redemptionToken = ref('')
const preflightToken = ref('')
// Grok 码不产预检票：预检通过后兑换直接带这份登录态（卡台按码上的产品分流）。
const grokSession = ref('')
const grokBlockingReason = ref('')
const credMode = ref<'session' | 'mailbox'>('session')
const sessionRaw = ref('')
const sessionFileInput = ref<HTMLInputElement | null>(null)
const sessionFileName = ref('')
const readingSessionFile = ref(false)
const isSessionDragActive = ref(false)
let sessionDragDepth = 0
const email = ref('')
const password = ref('')
/** 预检成功后的账号/订阅摘要（与卡台 GPT 直充 preflight 字段对齐） */
const account = ref({
  checked: false,
  email: '',
  currentPlan: '',
  xSubscriptionType: '', // Grok：X 订阅等级(Premium/…)，首月半价资格来源
  subscriptionHasActive: null as boolean | null,
  subscriptionActiveUntil: '',
  canPurchaseAt: '',
  subscriptionWillRenew: null as boolean | null,
  subscriptionRecoveryRequired: false,
  subscriptionIsDelinquent: null as boolean | null,
  lastPayment: null as any,
  paymentMethod: null as any,
})
const resultStatus = ref('')
const resultStage = ref('')
const resultMessage = ref('')
const resultEmail = ref('')
const resultCardLastFour = ref('')
const resultBody = ref<any>(null)
const timeline = ref<any[]>([])
const polling = ref(false)
const showAllEvents = ref(false)
let pollTimer: any = null
const displayResultEmail = computed(() => resultEmail.value || account.value.email || '')

const PROGRESS_KEY = 'cdk_redeem_progress_v1'
const nowTick = ref(Date.now())
let nowTimer: any = null

const deviceId = (() => {
  const k = 'cdk_device_id'
  let v = localStorage.getItem(k)
  if (!v) {
    v = 'web-' + Math.random().toString(36).slice(2) + Date.now().toString(36)
    localStorage.setItem(k, v)
  }
  return v
})()

function saveProgress() {
  try {
    const payload = {
      step: step.value,
      code: code.value,
      redemptionToken: redemptionToken.value,
      preflightToken: preflightToken.value,
      previewInfo: previewInfo.value,
      account: account.value,
      resultStatus: resultStatus.value,
      resultStage: resultStage.value,
      resultMessage: resultMessage.value,
      resultEmail: resultEmail.value,
      resultCardLastFour: resultCardLastFour.value,
      resultBody: resultBody.value,
      timeline: timeline.value,
      recoveryPending: recoveryPending.value,
      savedAt: Date.now(),
    }
    sessionStorage.setItem(PROGRESS_KEY, JSON.stringify(payload))
  } catch {
    /* ignore quota */
  }
}

function loadProgress(): boolean {
  try {
    const raw = sessionStorage.getItem(PROGRESS_KEY)
    if (!raw) return false
    const p = JSON.parse(raw)
    if (!p || typeof p !== 'object') return false
    // 超过 7 天丢弃
    if (p.savedAt && Date.now() - Number(p.savedAt) > 7 * 24 * 3600 * 1000) {
      sessionStorage.removeItem(PROGRESS_KEY)
      return false
    }
    if (p.code) code.value = String(p.code)
    if (p.redemptionToken) redemptionToken.value = String(p.redemptionToken)
    if (p.preflightToken) preflightToken.value = String(p.preflightToken)
    if (p.previewInfo) previewInfo.value = p.previewInfo
    if (p.account && typeof p.account === 'object') {
      account.value = {
        checked: !!p.account.checked,
        email: String(p.account.email || ''),
        currentPlan: String(p.account.currentPlan || ''),
        subscriptionHasActive:
          typeof p.account.subscriptionHasActive === 'boolean' ? p.account.subscriptionHasActive : null,
        subscriptionActiveUntil: String(p.account.subscriptionActiveUntil || ''),
        canPurchaseAt: String(p.account.canPurchaseAt || ''),
        subscriptionWillRenew:
          typeof p.account.subscriptionWillRenew === 'boolean' ? p.account.subscriptionWillRenew : null,
        subscriptionRecoveryRequired: p.account.subscriptionRecoveryRequired === true,
        subscriptionIsDelinquent: typeof p.account.subscriptionIsDelinquent === 'boolean' ? p.account.subscriptionIsDelinquent : null,
        lastPayment: p.account.lastPayment || null,
        paymentMethod: p.account.paymentMethod || null,
      }
    }
    recoveryPending.value = p.recoveryPending === true
    if (p.resultStatus) resultStatus.value = String(p.resultStatus)
    if (p.resultStage) resultStage.value = String(p.resultStage)
    if (p.resultMessage) resultMessage.value = String(p.resultMessage)
    if (p.resultEmail) resultEmail.value = String(p.resultEmail)
    if (p.resultCardLastFour) resultCardLastFour.value = String(p.resultCardLastFour)
    if (p.resultBody) resultBody.value = p.resultBody
    if (Array.isArray(p.timeline)) timeline.value = p.timeline
    const s = Number(p.step) || 1
    // 已提交兑换：恢复结果页
    if (redemptionToken.value && s >= 4) {
      step.value = 4
      return true
    }
    // 预检完成：恢复确认页（含订阅摘要）
    if (preflightToken.value && s === 3) {
      step.value = 3
      return true
    }
    if (s >= 1 && s <= 4) {
      step.value = s
      return s > 1
    }
  } catch {
    /* ignore */
  }
  return false
}

function clearProgress() {
  try {
    sessionStorage.removeItem(PROGRESS_KEY)
  } catch {
    /* ignore */
  }
}

watch(
  [step, code, redemptionToken, preflightToken, previewInfo, account, recoveryPending, resultStatus, resultStage, resultMessage, resultBody, timeline],
  () => saveProgress(),
  { deep: true },
)

const targetPlan = computed(() => {
  const plan = String(previewInfo.value?.plan || previewInfo.value?.plan_type || '').toLowerCase()
  // 新版卡台预览带 product：Grok / X 补上产品前缀再判，裸 premium_monthly 两家都有（2026-10-19 Grok 改名）。
  const product = String(previewInfo.value?.product || '').toLowerCase()
  if (plan && product === 'grok' && !plan.startsWith('grok_')) return 'grok_' + plan
  if (plan && product === 'x' && !plan.startsWith('x_')) return 'x_' + plan
  return plan
})
const isGrok = computed(() => isGrokPlan(targetPlan.value) && !isXPremiumPlan(targetPlan.value))
// GPT / X：必须有预检票；Grok：预检通过且没有阻塞项即可。
const canRedeem = computed(() => isGrok.value ? !!grokSession.value && !grokBlockingReason.value : !!preflightToken.value)

const subscriptionStatusText = computed(() => {
  if (account.value.subscriptionHasActive === true) return '有效'
  if (
    account.value.subscriptionHasActive === false
    || String(account.value.currentPlan).toLowerCase() === 'free'
  ) {
    return '已到期 / 免费版'
  }
  return '上游未提供'
})

const subscriptionStatusTag = computed(() => {
  if (account.value.subscriptionHasActive === true) return 'success'
  if (account.value.subscriptionHasActive === false) return 'info'
  return 'info'
})

const renewalStatusText = computed(() => {
  if (account.value.subscriptionWillRenew === true) return '到期后自动续费'
  if (account.value.subscriptionWillRenew === false) return '到期后不续费'
  return '上游未提供'
})

const subscriptionRemaining = computed(() =>
  remainingTime(account.value.subscriptionActiveUntil, nowTick.value),
)

const lastPaymentAtText = computed(() => {
  const lp = account.value.lastPayment
  if (!lp) return '上游未提供'
  const at = lp.paidAt || lp.paid_at || lp.created || lp.created_at
  return at ? fmtTime(at) : '上游未提供'
})

const lastPaymentText = computed(() => {
  const lp = account.value.lastPayment
  if (!lp) return '上游未提供'
  const amount = formatPaymentAmount(lp.amountMinor ?? lp.amount_minor, lp.currency)
  const status = String(lp.status || '').toLowerCase() === 'paid' ? '已支付' : (lp.status || '上游未提供')
  return amount ? `${amount} · ${status}` : status
})

const paymentMethodText = computed(() => {
  const method = account.value.paymentMethod
  if (!method) return '上游未提供'
  const brand = String(method.brand || method.type || '').toUpperCase()
  const last4 = method.last4 || method.last_4
  const label = [brand, last4 ? `**** ${last4}` : ''].filter(Boolean).join(' ')
  const expMonth = method.expMonth || method.exp_month
  const expYear = method.expYear || method.exp_year
  const expiry = expMonth && expYear
    ? `卡片到期 ${String(expMonth).padStart(2, '0')}/${expYear}`
    : ''
  const isDefault = method.isDefault ?? method.is_default
  return [label, expiry, isDefault ? '默认支付方式' : ''].filter(Boolean).join(' · ') || '上游未提供'
})

const needsSubscriptionRecovery = computed(() => account.value.checked &&
  (account.value.subscriptionRecoveryRequired || account.value.subscriptionIsDelinquent === true))

function invalidatePreflight() {
  preflightSequence++
  preflightToken.value = ''
  grokSession.value = ''
  grokBlockingReason.value = ''
  recoveryPending.value = false
  clearAccount()
}
watch([code, sessionRaw, email, password, credMode], invalidatePreflight, { flush: 'sync' })

function hasFileDrag(event: DragEvent): boolean {
  return Array.from(event.dataTransfer?.types || []).includes('Files')
}

function onSessionDragEnter(event: DragEvent) {
  if (!hasFileDrag(event)) return
  sessionDragDepth += 1
  isSessionDragActive.value = true
}

function onSessionDragOver(event: DragEvent) {
  if (hasFileDrag(event)) isSessionDragActive.value = true
}

function onSessionDragLeave() {
  sessionDragDepth = Math.max(0, sessionDragDepth - 1)
  if (sessionDragDepth === 0) isSessionDragActive.value = false
}

function resetSessionDrag() {
  sessionDragDepth = 0
  isSessionDragActive.value = false
}

function openSessionFilePicker() {
  if (busy.value || readingSessionFile.value) return
  sessionFileInput.value?.click()
}

function fileReadErrorMessage(error: unknown): string {
  if (!(error instanceof SessionFileError)) return t('recharge.fileErrors.read')
  const messages: Record<SessionFileError['code'], string> = {
    too_large: t('recharge.fileErrors.tooLarge'),
    unsupported: t('recharge.fileErrors.unsupported'),
    empty: t('recharge.fileErrors.empty'),
    invalid_docx: t('recharge.fileErrors.invalidDocx'),
    read_failed: t('recharge.fileErrors.read'),
  }
  return messages[error.code]
}

async function loadSessionFile(file: File | null) {
  if (!file) return
  readingSessionFile.value = true
  error.value = ''
  sessionFileName.value = ''
  try {
    sessionRaw.value = await readSessionFile(file)
    sessionFileName.value = file.name
  } catch (cause) {
    error.value = fileReadErrorMessage(cause)
  } finally {
    readingSessionFile.value = false
  }
}

function onSessionFileChange(event: Event) {
  const input = event.target as HTMLInputElement | null
  const file = input?.files?.[0] || null
  if (input) input.value = ''
  void loadSessionFile(file)
}

function onSessionFileDrop(event: DragEvent) {
  resetSessionDrag()
  const files = event.dataTransfer?.files
  if (!files?.length) return
  if (files.length > 1) {
    error.value = t('recharge.fileErrors.multiple')
    return
  }
  void loadSessionFile(files[0])
}

const alreadySatisfied = computed(() => {
  if (!account.value.checked || !targetPlan.value) return false
  if (account.value.subscriptionHasActive === false) return false
  return planSatisfied(account.value.currentPlan, targetPlan.value)
})

const alreadySatisfiedHint = computed(() => {
  const plan = planLabel(targetPlan.value)
  const until = account.value.canPurchaseAt || account.value.subscriptionActiveUntil
  if (until && account.value.subscriptionWillRenew === true) {
    return `账号已有 ${plan} 或更高套餐，当前周期至 ${fmtTime(until)} 并会自动续费，取消并到期前不能重复购买。`
  }
  if (until && account.value.subscriptionWillRenew === false) {
    return `账号已有 ${plan} 或更高套餐，最早可在 ${fmtTime(until)} 后再次购买。`
  }
  if (until) {
    return `账号已有 ${plan} 或更高套餐，当前周期至 ${fmtTime(until)}，暂不能重复购买。`
  }
  return `账号已有 ${plan} 或更高套餐，本次不能重复购买。`
})

// 档位判定（可读名 / 是否已满足）抽到 lib/plan.ts：那是纯函数、有单测钉着，
// 且「绑卡档不能被判成已满足」这条与卡台后端同源。planLabel 直接复用导入的实现；
// planSatisfied 在这里只做一层薄封装，把「当前码的 plan_flow」喂给纯函数。
function planSatisfied(currentPlan: string, requestedPlan: string) {
  return isSatisfied(currentPlan, requestedPlan, previewInfo.value?.plan_flow)
}

function remainingTime(value: string, timestamp = Date.now()) {
  const expiresAt = Date.parse(value || '')
  if (!Number.isFinite(expiresAt)) return '—'
  const minutes = Math.max(0, Math.ceil((expiresAt - timestamp) / 60000))
  if (minutes <= 0) return '已到期'
  const days = Math.floor(minutes / 1440)
  const hours = Math.floor((minutes % 1440) / 60)
  const restMinutes = minutes % 60
  if (days > 0) return `${days} 天 ${hours} 小时`
  if (hours > 0) return `${hours} 小时 ${restMinutes} 分钟`
  return `${restMinutes} 分钟`
}

function formatPaymentAmount(value: unknown, currency: unknown) {
  if (value === null || value === undefined || value === '') return ''
  const amount = Number(value)
  const code = String(currency || '').toUpperCase()
  if (!Number.isFinite(amount)) return ''
  let digits = 2
  if (code) {
    try {
      digits =
        new Intl.NumberFormat(undefined, { style: 'currency', currency: code }).resolvedOptions()
          .maximumFractionDigits ?? 2
    } catch {
      /* use two decimals */
    }
  }
  return `${code ? `${code} ` : ''}${(amount / 10 ** digits).toLocaleString(undefined, {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  })}`
}

function applyAccountFromPreflight(data: any) {
  const body = data?.data && typeof data.data === 'object' ? data.data : data || {}
  account.value = {
    checked: true,
    email: String(body.email || body.account_email || ''),
    currentPlan: String(body.currentPlan || body.current_plan || body.plan || ''),
    xSubscriptionType: String(body.x_subscription_type || ''),
    subscriptionHasActive:
      typeof body.subscription_has_active === 'boolean' ? body.subscription_has_active : null,
    subscriptionActiveUntil: String(body.subscription_active_until || ''),
    canPurchaseAt: String(body.can_purchase_at || body.subscription_active_until || ''),
    subscriptionWillRenew:
      typeof body.subscription_will_renew === 'boolean' ? body.subscription_will_renew : null,
    subscriptionRecoveryRequired: body.subscription_recovery_required === true,
    // 用户明确确认后处理宽限期订阅，再读取真实状态。
    subscriptionIsDelinquent:
      typeof body.subscription_is_delinquent === 'boolean' ? body.subscription_is_delinquent : null,
    lastPayment: body.last_payment || null,
    paymentMethod: body.payment_method || null,
  }
}

function clearAccount() {
  account.value = {
    checked: false,
    email: '',
    currentPlan: '',
    subscriptionHasActive: null,
    subscriptionActiveUntil: '',
    canPurchaseAt: '',
    subscriptionWillRenew: null,
    subscriptionIsDelinquent: null,
    subscriptionRecoveryRequired: false,
    lastPayment: null,
    paymentMethod: null,
  }
}

const TERMINAL = new Set(['completed', 'declined', 'failed_precharge', 'cancelled', 'failed'])

function isTerminal(st: string) {
  return TERMINAL.has(String(st || '').toLowerCase())
}

function stepLabel(stepKey: string) {
  const map: Record<string, string> = {
    queued: '排队受理',
    credential_check: '凭证校验',
    pricing: '计价',
    checkout: '开卡/绑卡',
    payment: '支付扣款',
    subscription: '订阅生效',
    invoice: '账单',
    renewal: '续费处理',
    reconcile: '对账确认',
    completed: '完成',
  }
  return map[stepKey] || stepKey || '处理中'
}

function fmtTime(v: any) {
  if (!v) return ''
  try {
    const d = new Date(v)
    if (Number.isNaN(d.getTime())) return String(v)
    return d.toLocaleString()
  } catch {
    return String(v)
  }
}

/* ── 兑换结果页展示层：只做派生，不改兑换/轮询逻辑 ─────────────────── */
type ResultTone = 'good' | 'warn' | 'bad' | 'info'
type StageState = 'done' | 'active' | 'warn' | 'bad' | 'idle'

const RESULT_STAGE_DEFS: { key: string; label: string; steps: string[] }[] = [
  { key: 'verify', label: '凭证校验', steps: ['queued', 'credential_check'] },
  { key: 'card', label: '开卡 · 绑卡', steps: ['checkout'] },
  { key: 'pay', label: '支付扣款', steps: ['payment', 'invoice'] },
  { key: 'reconcile', label: '对账确认', steps: ['reconcile'] },
  { key: 'activate', label: '订阅生效', steps: ['subscription', 'renewal'] },
]

function pad2(n: number) {
  return String(n).padStart(2, '0')
}

function fmtShortTime(v: any) {
  if (!v) return ''
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return String(v)
  return `${pad2(d.getHours())}:${pad2(d.getMinutes())}:${pad2(d.getSeconds())}`
}

function fmtDotDate(v: any) {
  if (!v) return '—'
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return '—'
  return `${d.getFullYear()}.${pad2(d.getMonth() + 1)}.${pad2(d.getDate())}`
}

function fmtDateSeconds(v: any) {
  if (!v) return '—'
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return String(v)
  return `${d.getFullYear()}/${pad2(d.getMonth() + 1)}/${pad2(d.getDate())} ${pad2(d.getHours())}:${pad2(d.getMinutes())}:${pad2(d.getSeconds())}`
}

function firstValue(...values: any[]) {
  for (const value of values) {
    if (value !== undefined && value !== null && String(value).trim() !== '') return String(value).trim()
  }
  return ''
}

function fmtDuration(seconds: number) {
  if (!Number.isFinite(seconds) || seconds <= 0) return ''
  if (seconds < 60) return `${Math.round(seconds)} 秒`
  const minutes = Math.floor(seconds / 60)
  const rest = Math.round(seconds % 60)
  return `${minutes} 分 ${rest} 秒`
}

function stageIndexFromStage(stage: string) {
  const value = String(stage || '').toLowerCase()
  if (value.includes('subscription') || value.includes('activate') || value.includes('paid') || value.includes('invoice')) return 4
  if (value.includes('reconcile') || value.includes('complete')) return 3
  if (value.includes('payment') || value.includes('checkout') || value.includes('spend') || value.includes('dispatch')) return 2
  if (value.includes('card') || value.includes('fund') || value.includes('await')) return 1
  return 0
}

function stageDuration(events: any[]) {
  if (events.length < 2) return ''
  const start = Date.parse(events[0]?.created_at || '')
  const end = Date.parse(events[events.length - 1]?.created_at || '')
  if (!Number.isFinite(start) || !Number.isFinite(end) || end <= start) return ''
  const seconds = (end - start) / 1000
  if (seconds < 1) return '<1s'
  if (seconds < 60) return `${Math.round(seconds)}s`
  return `${Math.floor(seconds / 60)}m ${Math.round(seconds % 60)}s`
}

function stageStateLabel(state: StageState) {
  return {
    done: '成功',
    active: '处理中',
    warn: '审核中',
    bad: '失败',
    idle: '等待',
  }[state]
}

function stageStateMessage(state: StageState) {
  return {
    done: '已完成',
    active: '正在处理，等待卡台返回结果…',
    warn: '等待审核结果，请勿重复提交。',
    bad: '该阶段未通过，请查看下方原因。',
    idle: '',
  }[state]
}

/** 完成阶段优先用成功事件的可读消息；接口消息仍是“正在准备…”时回落到阶段结论。 */
function stageDoneMessage(key: string, events: any[]) {
  const success = [...events].reverse().find((event) => {
    const category = String(event?.category || '').toLowerCase()
    return ['success', 'completed', 'done'].includes(category) && event?.public_message
  })
  const message = String(success?.public_message || '').trim()
  if (message && !/(正在|等待|准备|处理中|pending|running)/i.test(message)) return message
  return {
    verify: '凭证校验通过。',
    card: '支付卡片已绑定，资金确认到账。',
    pay: '支付已完成。',
    reconcile: '订单金额与卡台流水一致。',
    activate: '套餐已核验生效。',
  }[key] || '已完成'
}

function categoryLabel(value: any) {
  const category = String(value || '').toLowerCase()
  if (['success', 'completed', 'done'].includes(category)) return '成功'
  if (['failed', 'error'].includes(category)) return '失败'
  if (category === 'review') return '审核中'
  if (category === 'warning') return '需关注'
  if (category === 'running') return '处理中'
  if (category === 'queued') return '排队中'
  if (category === 'pending') return '等待中'
  return value || '处理中'
}

function categoryChipClass(value: any) {
  const category = String(value || '').toLowerCase()
  if (['success', 'completed', 'done'].includes(category)) return 'is-good'
  if (['failed', 'error'].includes(category)) return 'is-bad'
  if (['review', 'warning'].includes(category)) return 'is-warn'
  if (['running', 'queued'].includes(category)) return 'is-info'
  return 'is-muted'
}

const resultTone = computed<ResultTone>(() => {
  const status = String(resultStatus.value || '').toLowerCase()
  if (status === 'completed') return 'good'
  if (['review', 'card_open_review', 'card_recharge_review'].includes(status)) return 'warn'
  if (['declined', 'failed_precharge', 'cancelled', 'failed', 'error'].includes(status)) return 'bad'
  return 'info'
})

const resultToneText = computed(() => ({
  good: '已完成',
  warn: '审核中',
  bad: '未完成',
  info: '处理中',
}[resultTone.value]))

const statusBadge = computed(() => {
  const status = String(resultStatus.value || '').toLowerCase()
  const labels: Record<string, string> = {
    completed: 'COMPLETED',
    review: 'REVIEW',
    card_open_review: 'CARD REVIEW',
    card_recharge_review: 'RECHARGE REVIEW',
    pending: 'PENDING',
    queued: 'QUEUED',
    running: 'RUNNING',
    declined: 'DECLINED',
    cancelled: 'CANCELLED',
    failed_precharge: 'FAILED PRECHARGE',
    failed: 'FAILED',
    error: 'ERROR',
  }
  return labels[status] || (status ? status.replace(/_/g, ' ').toUpperCase() : 'PENDING')
})

const resultOrder = computed<any>(() => {
  const body = resultBody.value
  return body?.order || body?.data?.order || {}
})

const resultPlanRaw = computed(() => {
  const plan = firstValue(
    resultOrder.value.plan,
    resultOrder.value.plan_type,
    resultBody.value?.plan,
    resultBody.value?.data?.plan,
    targetPlan.value,
  )
  return plan.toLowerCase()
})

const resultPlanText = computed(() => {
  const plan = resultPlanRaw.value
  if (!plan) return '—'
  const label = planLabel(plan)
  const product = String(resultOrder.value.product || previewInfo.value?.product || '').toLowerCase()
  const isChatGPT = product === 'gpt'
    || (!product && (plan.includes('plus') || plan.includes('pro') || plan.includes('team') || plan.includes('chatgpt')))
  if (isChatGPT && !label.startsWith('ChatGPT')) return `ChatGPT ${label}`
  return label
})

const resultCdkCode = computed(() => firstValue(
  resultBody.value?.cdk_code,
  resultBody.value?.data?.cdk_code,
  code.value,
))

const resultLastAt = computed(() => {
  const lastEvent = timeline.value[timeline.value.length - 1]
  return firstValue(
    lastEvent?.created_at,
    resultOrder.value.completed_at,
    resultOrder.value.updated_at,
    resultBody.value?.completed_at,
  )
})

const resultCompletedText = computed(() => (resultLastAt.value ? fmtDateSeconds(resultLastAt.value) : '—'))
const receiptDate = computed(() => fmtDotDate(resultLastAt.value))
const resultTimeLabel = computed(() => (String(resultStatus.value || '').toLowerCase() === 'completed' ? '完成时间' : '最近更新'))

const paymentEvent = computed(() => {
  const payments = timeline.value.filter((event) => String(event?.step || '').toLowerCase() === 'payment')
  return payments.filter((event) => ['success', 'completed', 'done'].includes(String(event?.category || '').toLowerCase())).pop()
    || payments[payments.length - 1]
    || null
})

const paymentTimeText = computed(() => {
  if (isCardAttach.value) return '无需扣款'
  if (resultTone.value === 'good') {
    const at = paymentEvent.value?.created_at || resultLastAt.value
    return at ? `${fmtShortTime(at)} · Stripe` : '已完成'
  }
  const category = String(paymentEvent.value?.category || '').toLowerCase()
  if (paymentEvent.value?.created_at) {
    if (['failed', 'error'].includes(category)) return '支付失败'
    return '待确认'
  }
  if (resultTone.value === 'bad') return '未扣款'
  return '待扣款'
})

const isCardAttach = computed(() => {
  const flow = String(previewInfo.value?.plan_flow || '').toLowerCase()
  if (flow) return flow === 'card_attach'
  return /renew/.test(resultPlanRaw.value)
})

const resultTitle = computed(() => ({
  good: '兑换完成',
  warn: '兑换审核中',
  bad: '兑换未完成',
  info: resultStatus.value ? '兑换处理中' : '兑换进度',
}[resultTone.value]))

const resultSubtitle = computed(() => {
  if (resultTone.value === 'warn') return '订单已提交，正在审核。请勿重复提交，结果会自动刷新。'
  if (resultTone.value === 'bad') return '本次兑换未完成，请根据下方原因处理。'
  if (resultTone.value === 'info') {
    return resultStatus.value ? '正在处理，结果会自动刷新；请保持页面打开。' : '正在获取最新状态…'
  }
  return ''
})

const resultHeadTitle = computed(() => {
  if (resultTone.value === 'good') {
    const plan = resultPlanText.value === '—' ? '本次兑换' : resultPlanText.value
    return isCardAttach.value ? `${plan} 已完成` : `${plan} 已开通`
  }
  if (resultTone.value === 'warn') return '订单正在审核'
  if (resultTone.value === 'bad') return '本次兑换未完成'
  return '订单处理中'
})

const resultHeadVerb = computed(() => (resultTone.value === 'good' ? '兑换于' : '更新于'))

const resultDurationSeconds = computed(() => {
  if (!isTerminal(resultStatus.value)) return 0
  const times = timeline.value
    .map((event) => Date.parse(event?.created_at || ''))
    .filter((time) => Number.isFinite(time))
  if (!times.length) return 0
  const start = Math.min(...times)
  const end = Math.max(...times)
  return Math.max(0, Math.round((end - start) / 1000))
})

const resultStages = computed(() => {
  const groups = RESULT_STAGE_DEFS.map((def) => ({ key: def.key, label: def.label, events: [] as any[] }))
  for (const event of timeline.value) {
    const step = String(event?.step || '').toLowerCase()
    const index = RESULT_STAGE_DEFS.findIndex((def) => def.steps.includes(step))
    if (index >= 0) groups[index].events.push(event)
  }

  const status = String(resultStatus.value || '').toLowerCase()
  const completed = status === 'completed'
  const failed = ['declined', 'failed_precharge', 'cancelled', 'failed', 'error'].includes(status)
  const warning = ['review', 'card_open_review', 'card_recharge_review'].includes(status)

  let activeIndex = -1
  if (completed) {
    activeIndex = groups.length - 1
  } else if (timeline.value.length) {
    groups.forEach((group, index) => {
      if (group.events.length) activeIndex = index
    })
  } else if (status) {
    activeIndex = stageIndexFromStage(resultStage.value)
  }

  return groups.map((group, index) => {
    let state: StageState
    if (completed) state = 'done'
    else if (index < activeIndex) state = 'done'
    else if (index > activeIndex) state = 'idle'
    else if (failed) state = 'bad'
    else if (warning) state = 'warn'
    else state = 'active'

    if (!completed && index === activeIndex && group.events.some((event) => ['failed', 'error'].includes(String(event?.category || '').toLowerCase()))) {
      state = 'bad'
    }

    const timeSource = group.events[group.events.length - 1]

    return {
      key: group.key,
      label: group.label,
      state,
      events: group.events,
      count: group.events.length,
      time: timeSource
        ? fmtShortTime(timeSource.created_at)
        : (completed && index === groups.length - 1 ? fmtShortTime(resultLastAt.value) : ''),
    }
  })
})

const stageDoneCount = computed(() => resultStages.value.filter((stage) => stage.state === 'done').length)

const resultMetaText = computed(() => {
  const parts: string[] = []
  if (resultDurationSeconds.value > 0) parts.push(`耗时 ${fmtDuration(resultDurationSeconds.value)}`)
  parts.push(resultTone.value === 'good' ? `${stageDoneCount.value}/5 阶段` : `${stageDoneCount.value}/5 已完成`)
  if (polling.value) parts.push('实时刷新中')
  return parts.join(' · ')
})

const resultRows = computed(() => resultStages.value
  .filter((stage) => stage.count > 0 || ['active', 'warn', 'bad'].includes(stage.state))
  .map((stage) => {
    const latestMessage = [...stage.events].reverse().find((event) => event?.public_message)?.public_message
    return {
      key: stage.key,
      label: stage.label,
      state: stage.state,
      time: stage.time,
      chip: stageStateLabel(stage.state),
      message: stage.state === 'done'
        ? stageDoneMessage(stage.key, stage.events)
        : (latestMessage || stageStateMessage(stage.state)),
      duration: stageDuration(stage.events),
    }
  }))

const receiptStatusText = computed(() => {
  if (resultTone.value === 'good') return isCardAttach.value ? '已完成' : '已开通'
  return ({
    good: '已开通',
    warn: '审核中',
    bad: '未完成',
    info: '处理中',
  }[resultTone.value])
})

const resultFootText = computed(() => {
  if (resultTone.value === 'good') return '兑换已完成 · 结果已保存至本设备'
  if (resultTone.value === 'bad') return '兑换已结束 · 可凭 CDK 卡密联系客服'
  if (resultTone.value === 'warn') return '审核中 · 每 3 秒自动刷新'
  return '处理中 · 每 3 秒自动刷新'
})

async function copyText(value: string, okMessage: string) {
  if (!value) return
  const ok = await copyToClipboard(value)
  dialog.toast(ok ? okMessage : '复制失败，请手动选择文本', ok ? 'ok' : 'warn')
}

function copyReceipt() {
  const lines = [
    'LLDCard 兑换凭证',
    `状态：${receiptStatusText.value}`,
    `产品：${resultPlanText.value}`,
    `账号：${displayResultEmail.value || '—'}`,
    `CDK 卡密：${resultCdkCode.value || '—'}`,
    `完成时间：${resultCompletedText.value}`,
  ]
  void copyText(lines.join('\n'), '凭证信息已复制')
}

function extractCardLastFour(order: any): string {
  const last = String(order?.card_last_four || '').trim()
  if (/^\d{4}$/.test(last)) return last
  const n = String(order?.card_number || '').replace(/\D/g, '')
  return n.length >= 4 ? n.slice(-4) : ''
}

function applyResultPayload(data: any) {
  resultBody.value = data
  // 卡台公开结构：{ order: {status,stage,message,account_email,card_last_four}, events: [] }
  // 兼容顶层扁平 / data 包裹
  const order = data?.order || data?.data?.order || data?.data || data || {}
  const st =
    order.status ||
    data?.status ||
    data?.data?.status ||
    ''
  const stage = order.stage || data?.stage || data?.data?.stage || ''
  const message =
    order.message ||
    order.user_message ||
    data?.message ||
    data?.user_message ||
    data?.data?.message ||
    ''
  resultStatus.value = st
  resultStage.value = stage
  resultMessage.value = message
  const email = String(order.account_email || order.email || data?.account_email || '').trim()
  if (email) resultEmail.value = email
  const last4 = extractCardLastFour(order)
  if (last4) resultCardLastFour.value = last4

  let events = data?.events || data?.data?.events || order.events || []
  if (!Array.isArray(events)) events = []
  timeline.value = events.slice().sort((a: any, b: any) => {
    const ta = new Date(a.created_at || 0).getTime()
    const tb = new Date(b.created_at || 0).getTime()
    return ta - tb
  })
}

async function api(path: string, init: RequestInit = {}) {
  const headers = new Headers(init.headers || {})
  headers.set('X-Redemption-Device', deviceId)
  if (init.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json')
  const r = await fetch(path, { ...init, headers, credentials: 'include' })
  const text = await r.text()
  let data: any = null
  try { data = text ? JSON.parse(text) : null } catch { data = { raw: text } }
  return { r, data }
}

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
  if (nowTimer) clearInterval(nowTimer)
})

/** 完整 Session：必须含 sessionToken（JWE）；禁止纯 AT */
function extractSession(raw: string): string {
  const s = raw.trim()
  if (!s) return ''
  // 裸 JWE session-token
  if (!s.startsWith('{') && s.split('.').length >= 5) return s
  // 三段 JWT = 纯 AT → 拒绝
  if (!s.startsWith('{') && s.startsWith('eyJ') && s.split('.').length === 3) return ''
  if (s.startsWith('{')) {
    try {
      const o = JSON.parse(s)
      const st = String(o.sessionToken || o.session_token || o.token?.sessionToken || '').trim()
      if (!st) return ''
      // 回传完整 JSON，便于预检绑定 session
      return s
    } catch {
      return ''
    }
  }
  return s.length > 40 ? s : ''
}

async function tryResumeByCode(cdk: string): Promise<boolean> {
  const { r, data } = await api(
    '/api/v1/public/cdk/result-by-code?code=' + encodeURIComponent(cdk),
  )
  if (!r.ok) return false
  const tok = data?.redemption_token || data?.data?.redemption_token || ''
  if (tok) redemptionToken.value = tok
  applyResultPayload(data)
  if (!resultStatus.value) resultStatus.value = data?.status || data?.order?.status || 'pending'
  step.value = 4
  startPoll()
  return true
}

async function doPreview() {
  if (busy.value || recoveringSubscription.value) return
  invalidatePreflight()
  error.value = ''
  previewInfo.value = null
  if (!code.value.trim()) {
    error.value = '请输入 CDK'
    return
  }
  busy.value = true
  try {
    const cdk = code.value.trim()
    const { r, data } = await api('/api/v1/public/cdk/preview', {
      method: 'POST',
      body: JSON.stringify({ code: cdk }),
    })
    if (!r.ok) {
      const msg = data?.error || data?.msg || data?.message || 'CDK 无效或不可用'
      // 已兑换：尝试用本站绑定恢复进度，而不是卡在第一步
      if (/已兑换|已使用|used|redeemed|consumed|已消耗/i.test(String(msg))) {
        const ok = await tryResumeByCode(cdk)
        if (ok) {
          error.value = ''
          return
        }
      }
      error.value = msg
      return
    }
    // 兼容多种返回结构
    redemptionToken.value = data.redemption_token || data.data?.redemption_token || data.token || ''
    previewInfo.value = data.data || data
    if (isXPremiumPlan(targetPlan.value) || isGrokPlan(targetPlan.value)) credMode.value = 'session'
    if (!redemptionToken.value) {
      // 有的实现把 token 放在顶层其它字段
      error.value = '未返回 redemption_token，请检查卡台 Base 配置'
      return
    }
    step.value = 2
  } finally {
    busy.value = false
  }
}

async function doPreflight() {
  if (busy.value || recoveringSubscription.value || readingSessionFile.value) return
  invalidatePreflight()
  const sequence = preflightSequence
  error.value = ''
  busy.value = true
  try {
    let credential: any
    if (credMode.value === 'session') {
      const session = isGrok.value ? grokCredential(sessionRaw.value)
        : isXPremiumPlan(targetPlan.value) ? xPremiumCredential(sessionRaw.value) : extractSession(sessionRaw.value)
      if (!session) {
        error.value = isGrok.value ? t('grok.invalid')
          : isXPremiumPlan(targetPlan.value) ? t('xPremium.invalid') : '请粘贴完整 Session JSON（必须含 sessionToken），不能只用 Access Token'
        return
      }
      credential = { mode: 'session', session }
    } else {
      if (!email.value || !password.value) {
        error.value = '请填写邮箱与密码'
        return
      }
      credential = { mode: 'mailbox', email: email.value.trim(), password: password.value }
    }
    const { r, data } = await api('/api/v1/public/cdk/preflight', {
      method: 'POST',
      body: JSON.stringify({
        code: code.value.trim(),
        redemption_token: redemptionToken.value,
        credential,
      }),
    })
    if (sequence !== preflightSequence) return
    // 卡台 envelope: { code:0, msg, data:{ email, currentPlan, subscription_*, preflight_token } }
    if (!r.ok || (data && typeof data.code === 'number' && data.code !== 0)) {
      error.value = data?.error || data?.msg || data?.message || '预检失败'
      return
    }
    const body = data?.data && typeof data.data === 'object' ? data.data : data || {}
    if (isGrok.value) {
      // Grok：预检只查账号，没有 preflight_token。blocking_reason 是卡台的拦截判据（登录态失效 / 已有订阅等）。
      grokSession.value = credential.session
      grokBlockingReason.value = String(body.blocking_reason || '')
      applyAccountFromPreflight(data)
      step.value = 3
      return
    }
    preflightToken.value = body.preflight_token || data?.preflight_token || ''
    if (!preflightToken.value) {
      error.value = '未返回 preflight_token'
      return
    }
    applyAccountFromPreflight(data)
    step.value = 3
  } finally {
    busy.value = false
  }
}

async function recoverGraceSubscription() {
  if (busy.value || recoveringSubscription.value || !needsSubscriptionRecovery.value || !preflightToken.value) return
  const token = preflightToken.value, redemption = redemptionToken.value, sequence = preflightSequence
  recoveringSubscription.value = true
  try {
    if (!await dialog.confirm(t('grace.confirm'), { title: t('grace.title'), okText: t('grace.confirmAction'), cancelText: t('grace.dismiss'), danger: true })) return
    if (sequence !== preflightSequence || token !== preflightToken.value || redemption !== redemptionToken.value) return
    preflightToken.value = ''
    recoveryPending.value = true
    try {
      const { r, data } = await api('/api/v1/public/cdk/recover-subscription', {
        method: 'POST', body: JSON.stringify({ redemption_token: redemption, preflight_token: token, confirmed: true }),
      })
      if (sequence !== preflightSequence || redemption !== redemptionToken.value) return
      if (!r.ok || (typeof data?.code === 'number' && data.code !== 0)) throw new Error('recovery_unconfirmed')
      const result = data?.data || data
      if (result?.preflight) {
        applyAccountFromPreflight(result.preflight)
        preflightToken.value = result.preflight.preflight_token || ''
      }
      recoveryPending.value = result?.status !== 'cleared' || !preflightToken.value
      error.value = ''
      dialog.toast(t(recoveryPending.value ? 'grace.pending' : 'grace.done'), recoveryPending.value ? 'warn' : 'ok')
    } catch {
      if (sequence === preflightSequence) error.value = t('grace.refreshFailed')
    }
  } finally { recoveringSubscription.value = false }
}

async function doRedeem() {
  if (busy.value || recoveringSubscription.value || recoveryPending.value || needsSubscriptionRecovery.value || alreadySatisfied.value || !canRedeem.value) return
  error.value = ''
  busy.value = true
  try {
    const client_request_id = 'web-' + deviceId.slice(0, 8) + '-' + Date.now()
    const { r, data } = await api('/api/v1/public/cdk/redeem', {
      method: 'POST',
      body: JSON.stringify(isGrok.value
        ? { redemption_token: redemptionToken.value, credential: { mode: 'session', session: grokSession.value }, client_request_id }
        : { redemption_token: redemptionToken.value, preflight_token: preflightToken.value, client_request_id }),
    })
    applyResultPayload(data)
    if (!r.ok && r.status !== 202) {
      error.value = data?.error || data?.msg || data?.message || '兑换被拒绝'
    }
    if (!resultStatus.value) {
      resultStatus.value = r.ok || r.status === 202 ? 'queued' : 'error'
    }
    step.value = 4
    startPoll()
  } finally {
    busy.value = false
  }
}

function startPoll() {
  if (pollTimer) clearInterval(pollTimer)
  if (!redemptionToken.value && !code.value.trim()) {
    polling.value = false
    return
  }
  polling.value = true
  const tick = async () => {
    try {
      let r: Response
      let data: any
      if (redemptionToken.value) {
        ;({ r, data } = await api(
          '/api/v1/public/cdk/result?token=' + encodeURIComponent(redemptionToken.value),
        ))
      } else {
        ;({ r, data } = await api(
          '/api/v1/public/cdk/result-by-code?code=' + encodeURIComponent(code.value.trim()),
        ))
        if (r.ok && data?.redemption_token) {
          redemptionToken.value = data.redemption_token
        }
      }
      if (r.ok) {
        applyResultPayload(data)
        saveProgress()
        if (isTerminal(resultStatus.value)) {
          polling.value = false
          if (pollTimer) clearInterval(pollTimer)
        }
      }
    } catch {
      /* ignore transient network */
    }
  }
  tick()
  pollTimer = setInterval(tick, 3000)
}

onMounted(() => {
  nowTimer = setInterval(() => {
    nowTick.value = Date.now()
  }, 30000)
  const q = String(route.query.cdk || route.query.code || '').trim()
  if (loadProgress()) {
    if (q && code.value.trim() !== q) {
      resetAll()
      code.value = q
    } else if (step.value === 4 && (redemptionToken.value || code.value)) {
      startPoll()
    }
  } else if (q) {
    code.value = q
  }
})

function resetAll() {
  invalidatePreflight()
  resetSessionDrag()
  if (pollTimer) clearInterval(pollTimer)
  clearProgress()
  step.value = 1
  error.value = ''
  code.value = ''
  previewInfo.value = null
  redemptionToken.value = ''
  preflightToken.value = ''
  clearAccount()
  resultBody.value = null
  resultStatus.value = ''
  resultStage.value = ''
  resultMessage.value = ''
  resultEmail.value = ''
  resultCardLastFour.value = ''
  timeline.value = []
  polling.value = false
  showAllEvents.value = false
  sessionFileName.value = ''
}
</script>

<style scoped>
.text-good { color: var(--good, #16a34a); }
.text-warn { color: var(--warn, #d97706); }
.border-primary { border-color: var(--primary) !important; }
.bg-primary\/10 { background: color-mix(in srgb, var(--primary) 12%, transparent); }
.border-brd { border-color: var(--brd); }
.session-dropzone {
  position: relative;
  border-radius: 0.75rem;
}
.session-dropzone.is-dragging textarea {
  border-color: var(--primary) !important;
  background: color-mix(in srgb, var(--primary) 8%, var(--surface, #fff));
}
.session-drop-overlay {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  border: 2px dashed var(--primary);
  border-radius: 0.75rem;
  background: color-mix(in srgb, var(--surface, #fff) 88%, var(--primary));
  color: var(--ink, #111827);
  font-size: 14px;
  font-weight: 600;
  pointer-events: none;
}
.session-file-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px 12px;
}
.session-file-button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.account-facts {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px 16px;
  margin: 0;
}
.account-facts > div {
  min-width: 0;
}
.account-facts dt {
  margin: 0;
  font-size: 12px;
  color: var(--ink-3);
}
.account-facts dd {
  margin: 2px 0 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--ink, #111827);
  word-break: break-all;
}
@media (max-width: 640px) {
  .account-facts { grid-template-columns: 1fr; }
}

/* ── 兑换结果页（仅 step 4，不影响前三步）──────────────────────────── */
.result-wrap {
  max-width: 72rem;
  margin-left: auto;
  margin-right: auto;
  padding: 0 24px;
}
.result-subnav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 8px 0 18px;
  border-bottom: 1px solid var(--brd);
}
.result-crumb {
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.16em;
  color: var(--ink-3);
  text-transform: uppercase;
  white-space: nowrap;
}
.result-crumb b { color: var(--ink); font-weight: 600; }
.result-crumb i { font-style: normal; margin: 0 8px; color: var(--ink-3); opacity: 0.65; }

.result-hero {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 36px;
  padding: 34px 0 30px;
}
.result-hero-copy { min-width: 0; }
.result-eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  font-family: var(--font-mono);
  font-size: 10.5px;
  letter-spacing: 0.18em;
  color: var(--ink-3);
  text-transform: uppercase;
}
.result-eyebrow::before {
  content: '';
  width: 18px;
  height: 1px;
  background: var(--ink-3);
}
.result-eyebrow em { font-style: normal; color: var(--ink); }
.result-hero h1 {
  margin: 13px 0 0;
  font-family: var(--title-font, var(--font-display));
  font-size: clamp(2.2rem, 4.6vw, 3.375rem);
  font-weight: 800;
  font-style: var(--title-style, normal);
  letter-spacing: var(--title-track, -0.035em);
  line-height: 1.04;
  color: var(--ink);
  text-wrap: balance;
}
.result-hero p {
  margin: 14px 0 0;
  max-width: 34rem;
  font-size: 14px;
  line-height: 1.75;
  color: var(--ink-2);
}
.result-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: none;
}
.result-action {
  height: 42px;
  padding: 0 20px;
  border-radius: var(--btn-radius, 999px);
  font-size: 13px;
  white-space: nowrap;
}
.result-action svg { width: 13px; height: 13px; }

.result-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 22px;
  align-items: start;
  padding-bottom: 26px;
}
.result-card {
  padding: 28px 30px 24px;
  min-width: 0;
  box-shadow: 0 1px 0 color-mix(in srgb, var(--surface) 90%, transparent) inset, 0 26px 50px -46px rgba(0, 0, 0, 0.35);
}
.result-head {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}
.result-tone {
  display: grid;
  place-items: center;
  flex: none;
  width: 50px;
  height: 50px;
  border-radius: 15px;
  border: 1px solid transparent;
}
.result-tone svg { width: 22px; height: 22px; }
.result-tone.is-good { background: var(--good-soft); color: var(--good); border-color: color-mix(in srgb, var(--good) 28%, transparent); }
.result-tone.is-warn { background: var(--warn-soft); color: var(--warn); border-color: color-mix(in srgb, var(--warn) 28%, transparent); }
.result-tone.is-bad { background: var(--err-soft); color: var(--err); border-color: color-mix(in srgb, var(--err) 28%, transparent); }
.result-tone.is-info { background: var(--primary-soft); color: var(--primary); border-color: color-mix(in srgb, var(--primary) 24%, transparent); }
.result-head-copy { min-width: 0; }
.result-head-copy h2 {
  margin: 2px 0 0;
  font-size: 22px;
  font-weight: 800;
  letter-spacing: -0.02em;
  color: var(--ink);
}
.result-head-copy p {
  margin: 6px 0 0;
  font-size: 12.5px;
  line-height: 1.65;
  color: var(--ink-2);
  overflow-wrap: anywhere;
}
.result-email {
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 11.5px;
}
.result-head-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 7px;
  margin-left: auto;
  flex: none;
}
.result-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  border: 1px solid transparent;
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.1em;
  white-space: nowrap;
}
.result-pill i { width: 5px; height: 5px; border-radius: 50%; background: currentColor; }
.result-pill.is-good { color: var(--good); background: var(--good-soft); border-color: color-mix(in srgb, var(--good) 28%, transparent); }
.result-pill.is-warn { color: var(--warn); background: var(--warn-soft); border-color: color-mix(in srgb, var(--warn) 28%, transparent); }
.result-pill.is-bad { color: var(--err); background: var(--err-soft); border-color: color-mix(in srgb, var(--err) 28%, transparent); }
.result-pill.is-info { color: var(--primary); background: var(--primary-soft); border-color: color-mix(in srgb, var(--primary) 24%, transparent); }
.result-head-meta small {
  font-family: var(--font-mono);
  font-size: 10.5px;
  color: var(--ink-3);
  letter-spacing: 0.01em;
  white-space: nowrap;
}
.result-alert { margin-top: 18px; }

.result-rail {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  margin: 26px 0 0;
  padding: 0;
  list-style: none;
}
.result-stage {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  min-width: 0;
  text-align: center;
}
.result-stage:not(:first-child)::before {
  content: '';
  position: absolute;
  top: 12px;
  right: calc(50% + 17px);
  width: calc(100% - 34px);
  height: 2px;
  border-radius: 2px;
  background: var(--brd-2);
}
.result-stage.is-done::before,
.result-stage.is-active::before,
.result-stage.is-warn::before,
.result-stage.is-bad::before { background: var(--ink); }
.result-dot {
  position: relative;
  z-index: 1;
  display: grid;
  place-items: center;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: var(--ink);
  color: var(--bg);
  font-family: var(--font-mono);
  font-size: 10px;
}
.result-dot svg { width: 12px; height: 12px; }
.result-stage.is-warn .result-dot { background: var(--warn); color: #fff; }
.result-stage.is-bad .result-dot { background: var(--err); color: #fff; }
.result-stage.is-idle .result-dot {
  background: var(--surface);
  color: var(--ink-3);
  border: 1px solid var(--brd-2);
}
.result-stage.is-done:last-child .result-dot { background: var(--good); }
.result-stage b {
  max-width: 100%;
  overflow: hidden;
  color: var(--ink);
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.01em;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.result-stage.is-idle b { color: var(--ink-3); font-weight: 500; }
.result-stage time {
  min-height: 12px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10px;
}
.result-pulse {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--bg);
  animation: result-pulse 1.8s ease-out infinite;
}
@keyframes result-pulse {
  0% { box-shadow: 0 0 0 0 color-mix(in srgb, var(--bg) 45%, transparent); }
  70% { box-shadow: 0 0 0 6px transparent; }
  100% { box-shadow: 0 0 0 0 transparent; }
}

.result-section {
  margin-top: 26px;
  padding-top: 22px;
  border-top: 1px solid var(--brd);
}
.result-section-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 16px;
}
.result-section-head h3 {
  margin: 0;
  color: var(--ink);
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.02em;
}
.result-section-head > span {
  color: var(--ink-3);
  font-size: 11px;
  white-space: nowrap;
}
.result-link {
  padding: 0 0 2px;
  border: 0;
  border-bottom: 1px solid var(--brd-2);
  background: transparent;
  color: var(--ink-2);
  font: inherit;
  font-size: 11.5px;
  cursor: pointer;
  white-space: nowrap;
}
.result-link:hover { color: var(--ink); border-color: var(--ink-3); }
.result-facts {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 19px 26px;
  margin: 18px 0 0;
}
.result-fact { min-width: 0; }
.result-fact dt {
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 9.5px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}
.result-fact dd {
  display: flex;
  align-items: center;
  gap: 2px;
  margin: 6px 0 0;
  overflow: hidden;
  color: var(--ink);
  font-size: 13.5px;
  font-weight: 600;
  line-height: 1.45;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.result-fact dd.is-mono {
  font-family: var(--font-mono);
  font-size: 12.5px;
  font-weight: 500;
  letter-spacing: -0.01em;
}
.result-copy {
  display: inline-grid;
  place-items: center;
  flex: none;
  width: 20px;
  height: 20px;
  margin-left: 3px;
  border: 0;
  background: transparent;
  color: var(--ink-3);
  cursor: pointer;
}
.result-copy svg { width: 12px; height: 12px; }
.result-copy:hover { color: var(--ink); }

.result-events {
  position: relative;
  margin: 18px 0 0;
  padding: 0;
  list-style: none;
}
.result-events::before {
  content: '';
  position: absolute;
  top: 12px;
  bottom: 16px;
  left: 72px;
  width: 1px;
  background: var(--brd);
}
.result-event {
  display: grid;
  grid-template-columns: 58px 28px minmax(0, 1fr);
  padding-bottom: 17px;
}
.result-event:last-child { padding-bottom: 0; }
.result-event time {
  padding-top: 1px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 11px;
}
.result-node {
  z-index: 1;
  width: 9px;
  height: 9px;
  margin: 4px auto 0;
  border-radius: 50%;
  background: var(--ink);
  box-shadow: 0 0 0 4px var(--surface);
}
.result-event.is-active .result-node { background: var(--primary); }
.result-event.is-warn .result-node { background: var(--warn); }
.result-event.is-bad .result-node { background: var(--err); }
.result-event.is-idle .result-node { background: var(--brd-2); }
.result-event.is-done:last-child .result-node { background: var(--good); }
.result-event-body { min-width: 0; }
.result-event-top {
  display: flex;
  align-items: center;
  gap: 8px;
}
.result-event-top b {
  color: var(--ink);
  font-size: 13px;
  font-weight: 700;
}
.result-chip {
  padding: 3px 7px;
  border: 1px solid transparent;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 600;
  line-height: 1.2;
  white-space: nowrap;
}
.result-chip.is-good { color: var(--good); background: var(--good-soft); border-color: color-mix(in srgb, var(--good) 24%, transparent); }
.result-chip.is-warn { color: var(--warn); background: var(--warn-soft); border-color: color-mix(in srgb, var(--warn) 24%, transparent); }
.result-chip.is-bad { color: var(--err); background: var(--err-soft); border-color: color-mix(in srgb, var(--err) 24%, transparent); }
.result-chip.is-info { color: var(--primary); background: var(--primary-soft); border-color: color-mix(in srgb, var(--primary) 20%, transparent); }
.result-chip.is-muted { color: var(--ink-3); background: var(--surface-2); border-color: var(--brd); }
.result-dur {
  margin-left: auto;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10px;
}
.result-event p {
  margin: 4px 0 0;
  color: var(--ink-2);
  font-size: 12.5px;
  line-height: 1.6;
}
.result-empty {
  margin: 16px 0 0;
  color: var(--ink-3);
  font-size: 12px;
}
.result-raw {
  margin: 14px 0 0;
  padding: 12px 14px;
  border-radius: 12px;
  background: var(--surface-2);
  list-style: none;
}
.result-raw li {
  display: grid;
  grid-template-columns: 64px 76px 54px minmax(0, 1fr);
  gap: 8px;
  align-items: baseline;
  padding: 5px 0;
  border-top: 1px solid var(--brd);
  color: var(--ink-2);
  font-size: 11px;
}
.result-raw li:first-child { border-top: 0; }
.result-raw-time { font-family: var(--font-mono); color: var(--ink-3); }
.result-raw li b { font-size: 11.5px; color: var(--ink); }
.result-raw .result-chip { justify-self: start; padding: 2px 6px; font-size: 9px; }
.result-raw-msg { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.result-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px dashed var(--brd);
  color: var(--ink-3);
  font-size: 11px;
}
.result-foot-live {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10px;
  white-space: nowrap;
}
.result-foot-live i {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--good);
  animation: result-pulse 1.8s ease-out infinite;
}
.result-advanced {
  margin-top: 18px;
  padding-top: 13px;
  border-top: 1px solid var(--brd);
  color: var(--ink-3);
  font-size: 11px;
}
.result-advanced-summary { cursor: pointer; user-select: none; }
.result-advanced pre {
  max-height: 220px;
  margin: 10px 0 0;
  padding: 12px;
  overflow: auto;
  border-radius: 10px;
  background: var(--surface-2);
  color: var(--ink-2);
  font-family: var(--font-mono);
  font-size: 10.5px;
  line-height: 1.55;
}

.result-aside {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}
.result-receipt {
  position: relative;
  overflow: hidden;
  padding: 24px 24px 20px;
  border-radius: var(--radius-lg);
  background:
    radial-gradient(360px 200px at 18% -22%, color-mix(in srgb, var(--bg) 14%, transparent), transparent 62%),
    radial-gradient(300px 220px at 108% 112%, color-mix(in srgb, var(--good) 16%, transparent), transparent 64%),
    var(--ink);
  color: var(--bg);
  box-shadow: 0 30px 60px -40px rgba(0, 0, 0, 0.65);
}
.receipt-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: color-mix(in srgb, var(--bg) 52%, transparent);
  font-family: var(--font-mono);
  font-size: 9px;
  letter-spacing: 0.16em;
  white-space: nowrap;
}
.receipt-top span:last-child {
  max-width: 68%;
  overflow: hidden;
  text-overflow: ellipsis;
}
.receipt-status {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 20px;
}
.receipt-ring {
  display: grid;
  place-items: center;
  flex: none;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: 1px solid color-mix(in srgb, var(--bg) 20%, transparent);
  background: color-mix(in srgb, var(--bg) 8%, transparent);
}
.receipt-ring svg { width: 15px; height: 15px; }
.receipt-ring.is-good { color: #fff; background: var(--good); border-color: transparent; }
.receipt-ring.is-warn { color: #fff; background: var(--warn); border-color: transparent; }
.receipt-ring.is-bad { color: #fff; background: var(--err); border-color: transparent; }
.receipt-ring.is-info { color: var(--primary-on, #fff); background: var(--primary); border-color: transparent; }
.receipt-status b { display: block; font-size: 15px; font-weight: 700; letter-spacing: 0.01em; }
.receipt-status span {
  display: block;
  margin-top: 3px;
  color: color-mix(in srgb, var(--bg) 46%, transparent);
  font-family: var(--font-mono);
  font-size: 9px;
  letter-spacing: 0.1em;
}
.receipt-product {
  margin-top: 18px;
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 800;
  letter-spacing: -0.03em;
  line-height: 1.1;
}
.receipt-sub {
  margin-top: 6px;
  color: color-mix(in srgb, var(--bg) 52%, transparent);
  font-size: 11.5px;
}
.receipt-perf {
  position: relative;
  height: 1px;
  margin: 22px -24px 20px;
  background: repeating-linear-gradient(90deg, color-mix(in srgb, var(--bg) 22%, transparent) 0 5px, transparent 5px 11px);
}
.receipt-perf::before,
.receipt-perf::after {
  content: '';
  position: absolute;
  top: -10px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--bg);
}
.receipt-perf::before { left: -10px; }
.receipt-perf::after { right: -10px; }
.receipt-rows { display: flex; flex-direction: column; gap: 13px; margin: 0; }
.receipt-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 14px;
}
.receipt-row dt {
  color: color-mix(in srgb, var(--bg) 46%, transparent);
  font-family: var(--font-mono);
  font-size: 9px;
  letter-spacing: 0.12em;
  white-space: nowrap;
}
.receipt-row dd {
  max-width: 72%;
  margin: 0;
  overflow: hidden;
  color: color-mix(in srgb, var(--bg) 90%, transparent);
  font-family: var(--font-mono);
  font-size: 11px;
  text-align: right;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.receipt-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  height: 39px;
  margin-top: 22px;
  border: 0;
  border-radius: 999px;
  background: var(--bg);
  color: var(--ink);
  font: inherit;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.15s ease, opacity 0.15s ease;
}
.receipt-btn:hover { transform: translateY(-1px); opacity: 0.92; }
.receipt-btn svg { width: 13px; height: 13px; }
.receipt-link {
  display: block;
  width: max-content;
  margin: 13px auto 0;
  border-bottom: 1px solid color-mix(in srgb, var(--bg) 24%, transparent);
  color: color-mix(in srgb, var(--bg) 62%, transparent);
  font-size: 11px;
  text-decoration: none;
}
.receipt-link:hover { color: var(--bg); }
.receipt-foot {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 18px;
  padding-top: 14px;
  border-top: 1px solid color-mix(in srgb, var(--bg) 10%, transparent);
  color: color-mix(in srgb, var(--bg) 42%, transparent);
  font-size: 10px;
  line-height: 1.5;
  text-align: center;
}
.receipt-foot svg { width: 10px; height: 10px; flex: none; }

.result-next {
  padding: 20px 22px 8px;
  border: 1px solid var(--brd);
  border-radius: var(--radius-lg);
  background: var(--surface);
}
.result-next h3 {
  margin: 0 0 2px;
  color: var(--ink);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.02em;
}
.result-next ol { margin: 0; padding: 0; list-style: none; }
.result-next li {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 13px 0;
  border-top: 1px solid var(--brd);
}
.result-next li:first-child { border-top: 0; }
.result-next-num {
  padding-top: 2px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 9.5px;
  letter-spacing: 0.08em;
}
.result-next b { display: block; color: var(--ink); font-size: 12.5px; font-weight: 700; }
.result-next p { margin: 3px 0 0; color: var(--ink-3); font-size: 11px; line-height: 1.55; }
.result-mini {
  flex: none;
  margin-left: auto;
  padding: 4px 9px;
  border: 1px solid var(--brd-2);
  border-radius: 999px;
  background: transparent;
  color: var(--ink-2);
  font: inherit;
  font-size: 10px;
  cursor: pointer;
}
.result-mini:hover { color: var(--ink); border-color: var(--ink-3); }
.result-device-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 2px 4px;
  color: var(--ink-3);
  font-size: 11px;
  line-height: 1.65;
}
.result-device-note svg { width: 12px; height: 12px; margin-top: 2px; flex: none; }

@media (max-width: 900px) {
  .result-grid { grid-template-columns: minmax(0, 1fr); }
  .result-aside { order: 2; }
  .result-facts { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .result-events::before { left: 62px; }
  .result-event { grid-template-columns: 48px 26px minmax(0, 1fr); }
  .result-head { flex-wrap: wrap; }
  .result-head-meta {
    margin-left: 0;
    width: 100%;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}

@media (min-width: 641px) and (max-width: 900px) {
  .result-aside { max-width: 460px; }
}

@media (max-width: 640px) {
  .result-wrap { padding: 0 16px; }
  .result-subnav { padding: 4px 0 14px; }
  .result-crumb { display: none; }
  .result-hero {
    flex-direction: column;
    align-items: stretch;
    gap: 18px;
    padding: 22px 0;
  }
  .result-hero h1 { margin-top: 10px; font-size: 38px; letter-spacing: -0.04em; }
  .result-hero p { margin-top: 10px; font-size: 12.5px; line-height: 1.7; }
  .result-eyebrow { font-size: 9.5px; }
  .result-actions { gap: 8px; }
  .result-action { flex: 1; height: 40px; padding: 0 14px; font-size: 12.5px; }
  .result-card { padding: 18px 16px 16px; border-radius: 18px; }
  .result-tone { width: 42px; height: 42px; border-radius: 12px; }
  .result-tone svg { width: 18px; height: 18px; }
  .result-head { gap: 11px; }
  .result-head-copy h2 { margin-top: 3px; font-size: 18px; }
  .result-head-copy p { margin-top: 5px; font-size: 11px; line-height: 1.55; }
  .result-head-meta { margin-top: 2px; }
  .result-head-meta small { font-size: 9.5px; }
  .result-pill { height: 22px; font-size: 9px; }
  .result-rail { margin-top: 18px; }
  .result-stage { gap: 7px; }
  .result-stage:not(:first-child)::before {
    top: 10px;
    right: calc(50% + 14px);
    width: calc(100% - 28px);
    height: 1.5px;
  }
  .result-dot { width: 22px; height: 22px; font-size: 9px; }
  .result-dot svg { width: 10px; height: 10px; }
  .result-stage b { font-size: 10px; }
  .result-stage time { font-size: 8.5px; }
  .result-section { margin-top: 20px; padding-top: 16px; }
  .result-section-head h3 { font-size: 12px; }
  .result-section-head > span { font-size: 10px; }
  .result-link { font-size: 10.5px; }
  .result-facts { gap: 15px 14px; margin-top: 15px; }
  .result-fact dt { font-size: 8.5px; letter-spacing: 0.1em; }
  .result-fact dd { margin-top: 5px; font-size: 12px; }
  .result-fact dd.is-mono { font-size: 10.5px; }
  .result-fact--email,
  .result-fact--code { grid-column: 1 / -1; }
  .result-fact--code dd { font-size: 12px; }
  .result-events { margin-top: 15px; }
  .result-events::before { left: 59px; top: 10px; bottom: 12px; }
  .result-event { grid-template-columns: 48px 22px minmax(0, 1fr); padding-bottom: 15px; }
  .result-event time { font-size: 9.5px; letter-spacing: -0.03em; }
  .result-node { width: 8px; height: 8px; margin-top: 4px; }
  .result-event-top { gap: 6px; }
  .result-event-top b { font-size: 12px; }
  .result-chip { padding: 3px 6px; font-size: 9px; }
  .result-dur { font-size: 9px; }
  .result-event p { font-size: 11.5px; line-height: 1.55; }
  .result-raw { padding: 10px 12px; }
  .result-raw li { grid-template-columns: 56px 64px 48px minmax(0, 1fr); gap: 6px; font-size: 10px; }
  .result-foot { flex-direction: column; align-items: flex-start; gap: 7px; margin-top: 18px; padding-top: 13px; font-size: 10px; }
  .result-receipt { padding: 22px 20px 18px; border-radius: 18px; }
  .receipt-perf { margin-left: -20px; margin-right: -20px; }
  .receipt-product { font-size: 26px; }
  .receipt-status { margin-top: 18px; }
  .receipt-row dd { max-width: 66%; }
  .result-next { padding: 18px 18px 6px; border-radius: 18px; }
  .result-next b { font-size: 12px; }
  .result-next p { font-size: 10.5px; }
  .result-device-note { font-size: 10px; }
}

/* ── 兑换第 1 步：输入卡密 ─────────────────────────────────────────── */
.cdk-hero-aside {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
  flex: none;
  padding-bottom: 6px;
  text-align: right;
}
.cdk-hero-aside > span { color: var(--ink-3); font-size: 11px; }
.cdk-hero-aside a {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding-bottom: 2px;
  border-bottom: 1px solid var(--brd-2);
  color: var(--ink);
  font-size: 12px;
  font-weight: 600;
  text-decoration: none;
}
.cdk-hero-aside a svg { width: 11px; height: 11px; }

.cdk-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 18px 0 0;
  padding: 0;
  list-style: none;
}
.cdk-badges li {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 11px;
  border: 1px solid var(--brd);
  border-radius: 999px;
  background: var(--surface);
  color: var(--ink-2);
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}
.cdk-badges svg { width: 13px; height: 13px; color: var(--ink-3); }

.cdk-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 22px;
  align-items: start;
  padding-bottom: 26px;
}
.cdk-col { display: flex; flex-direction: column; gap: 16px; min-width: 0; }
.cdk-aside { min-width: 0; }

.cdk-form-card { padding: 28px 30px 24px; }
.cdk-form-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; }
.cdk-micro {
  display: block;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 9.5px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}
.cdk-form-head h2 {
  margin: 8px 0 0;
  color: var(--ink);
  font-size: 24px;
  font-weight: 800;
  letter-spacing: -0.025em;
}
.cdk-form-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex: none;
  padding: 5px 9px;
  border: 1px solid var(--good-line, #cde8d9);
  border-radius: 999px;
  background: var(--good-soft);
  color: var(--good);
  font-size: 10px;
  font-weight: 600;
  white-space: nowrap;
}
.cdk-form-badge i { width: 5px; height: 5px; border-radius: 50%; background: var(--good); }

.cdk-field { margin-top: 24px; }
.cdk-field-label {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 9px;
}
.cdk-field-label span:first-child { color: var(--ink-2); font-size: 11.5px; font-weight: 600; }
.cdk-field-label span:last-child { color: var(--ink-3); font-size: 10px; }
.cdk-input-shell {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 54px;
  padding: 0 16px;
  border: 1px solid var(--brd);
  border-radius: 14px;
  background: var(--surface-2);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
.cdk-input-shell:focus-within {
  border-color: var(--ink-3);
  box-shadow: 0 0 0 3px var(--primary-soft);
}
.cdk-input-shell > svg { width: 18px; height: 18px; flex: none; color: var(--ink-3); }
.cdk-input {
  width: 100%;
  min-width: 0;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 14px;
  letter-spacing: 0.01em;
}
.cdk-input::placeholder { color: var(--ink-4); }
.cdk-hint { margin-top: 8px; color: var(--ink-3); font-size: 10.5px; line-height: 1.5; }
.cdk-error { margin-top: 14px; }
.cdk-preview {
  margin-top: 14px;
  padding: 12px 14px;
  border: 1px solid var(--brd);
  border-radius: 12px;
  background: var(--surface-2);
}
.cdk-preview-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  color: var(--ink-3);
  font-size: 11.5px;
}
.cdk-preview-row b { color: var(--ink); font-weight: 700; }
.cdk-preview p { margin: 4px 0 0; color: var(--ink-3); font-size: 10.5px; line-height: 1.5; }
.cdk-cta.btn-primary {
  width: 100%;
  height: 48px;
  padding: 0 20px;
  margin-top: 22px;
  border-radius: 999px;
  font-size: 14px;
}
.cdk-cta svg { width: 14px; height: 14px; }
.cdk-note {
  display: flex;
  align-items: center;
  gap: 7px;
  margin: 12px 0 0;
  color: var(--ink-3);
  font-size: 10.5px;
  line-height: 1.5;
}
.cdk-note svg { width: 12px; height: 12px; flex: none; }

.cdk-tips-card { padding: 20px 22px 8px; }
.cdk-tips-card h3 { margin: 0 0 2px; color: var(--ink); font-size: 12px; font-weight: 800; letter-spacing: 0.02em; }
.cdk-tips-list { margin: 0; padding: 0; list-style: none; }
.cdk-tips-list li {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 13px 0;
  border-top: 1px solid var(--brd);
}
.cdk-tips-list li:first-child { border-top: 0; }
.cdk-tips-num {
  padding-top: 2px;
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 9.5px;
  letter-spacing: 0.08em;
}
.cdk-tips-list b { display: block; color: var(--ink); font-size: 12.5px; font-weight: 700; }
.cdk-tips-list p { margin: 3px 0 0; color: var(--ink-3); font-size: 11px; line-height: 1.55; }

.cdk-flow-card {
  position: relative;
  overflow: hidden;
  padding: 24px 26px;
  border-radius: var(--radius-lg);
  background:
    radial-gradient(360px 200px at 18% -22%, color-mix(in srgb, var(--bg) 13%, transparent), transparent 62%),
    radial-gradient(300px 220px at 108% 112%, color-mix(in srgb, var(--good) 14%, transparent), transparent 64%),
    var(--ink);
  color: var(--bg);
  box-shadow: 0 30px 60px -40px rgba(0, 0, 0, 0.65);
}
.cdk-flow-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: color-mix(in srgb, var(--bg) 48%, transparent);
  font-family: var(--font-mono);
  font-size: 9px;
  letter-spacing: 0.16em;
  white-space: nowrap;
}
.cdk-flow-card h3 {
  margin: 20px 0 0;
  color: var(--bg);
  font-family: var(--font-display);
  font-size: 25px;
  font-weight: 800;
  letter-spacing: -0.025em;
}
.cdk-flow-sub { margin: 6px 0 0; color: color-mix(in srgb, var(--bg) 46%, transparent); font-size: 11.5px; }
.cdk-flow-steps { position: relative; margin: 22px 0 0; padding: 0; list-style: none; }
.cdk-flow-steps::before {
  content: '';
  position: absolute;
  top: 17px;
  bottom: 28px;
  left: 13px;
  width: 1px;
  background: color-mix(in srgb, var(--bg) 14%, transparent);
}
.cdk-flow-steps li {
  position: relative;
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr) auto;
  gap: 12px;
  align-items: flex-start;
  padding: 10px 0;
}
.cdk-flow-num {
  z-index: 1;
  display: grid;
  place-items: center;
  width: 27px;
  height: 27px;
  border: 1px solid color-mix(in srgb, var(--bg) 18%, transparent);
  border-radius: 50%;
  background: var(--ink);
  color: color-mix(in srgb, var(--bg) 40%, transparent);
  font-family: var(--font-mono);
  font-size: 9px;
}
.cdk-flow-steps b {
  display: block;
  padding-top: 3px;
  color: color-mix(in srgb, var(--bg) 56%, transparent);
  font-size: 12.5px;
  font-weight: 600;
}
.cdk-flow-steps p { margin: 3px 0 0; color: color-mix(in srgb, var(--bg) 30%, transparent); font-size: 10.5px; line-height: 1.5; }
.cdk-flow-steps li.is-active .cdk-flow-num { border-color: var(--bg); background: var(--bg); color: var(--ink); font-weight: 700; }
.cdk-flow-steps li.is-active b { color: var(--bg); }
.cdk-flow-steps li.is-active p { color: color-mix(in srgb, var(--bg) 48%, transparent); }
.cdk-flow-steps em {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding-top: 5px;
  color: var(--good);
  font-family: var(--font-mono);
  font-size: 8.5px;
  font-style: normal;
  letter-spacing: 0.08em;
}
.cdk-flow-steps em::before { content: ''; width: 4px; height: 4px; border-radius: 50%; background: var(--good); }

.cdk-security {
  padding: 18px 20px;
  border: 1px solid var(--brd);
  border-radius: 18px;
  background: var(--surface);
}
.cdk-security h3 { margin: 0; color: var(--ink); font-size: 11.5px; font-weight: 800; letter-spacing: 0.02em; }
.cdk-security p { margin: 6px 0 0; color: var(--ink-3); font-size: 10.5px; line-height: 1.65; }

@media (max-width: 900px) {
  .cdk-hero { flex-direction: column; align-items: stretch; gap: 20px; padding: 26px 0 24px; }
  .cdk-hero-aside { align-items: flex-start; padding: 0; text-align: left; }
  .cdk-grid { grid-template-columns: minmax(0, 1fr); }
  .cdk-aside { order: 2; }
  .cdk-form-card { padding: 24px 22px 22px; }
}

@media (min-width: 641px) and (max-width: 900px) {
  .cdk-aside { max-width: 460px; }
}

@media (max-width: 640px) {
  .cdk-hero { gap: 16px; padding: 22px 0; }
  .cdk-hero-aside { display: none; }
  .cdk-badges { margin-top: 15px; gap: 6px; }
  .cdk-badges li { padding: 5px 9px; font-size: 10.5px; }
  .cdk-form-card { padding: 20px 16px 18px; border-radius: 18px; }
  .cdk-form-head h2 { font-size: 21px; }
  .cdk-form-badge { padding: 4px 8px; font-size: 9px; }
  .cdk-field { margin-top: 20px; }
  .cdk-input-shell { height: 48px; padding: 0 13px; border-radius: 12px; }
  .cdk-input { font-size: 12px; }
  .cdk-cta.btn-primary { height: 44px; margin-top: 18px; font-size: 13px; }
  .cdk-tips-card { padding: 18px 18px 6px; border-radius: 18px; }
  .cdk-tips-list li { padding: 12px 0; }
  .cdk-flow-card { padding: 22px 20px; border-radius: 18px; }
  .cdk-flow-card h3 { font-size: 23px; }
  .cdk-flow-steps b { font-size: 12px; }
  .cdk-flow-steps p { font-size: 10px; }
  .cdk-security { border-radius: 16px; }
}
</style>
