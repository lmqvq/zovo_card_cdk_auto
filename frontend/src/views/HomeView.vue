<template>
  <div class="u-page home-page">
    <header class="home-topbar">
      <div class="home-topbar-inner">
        <router-link to="/" class="home-brand" :aria-label="t('common.backHome')">
          <span class="home-logo" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M6 3h12l4 6-10 12L2 9z" />
              <path d="M11 3 8 9l4 12 4-12-3-6" />
              <path d="M2 9h20" />
            </svg>
          </span>
          <span class="home-brand-copy">
            <strong>{{ brand.name || t('home.brand') }}</strong>
            <span class="home-brand-divider" aria-hidden="true">|</span>
            <small>{{ brand.sub || t('home.brandSub') }}</small>
          </span>
        </router-link>

        <nav class="home-nav" :aria-label="t('home.nav.label')">
          <a href="#operation-guide" class="home-nav-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <circle cx="12" cy="12" r="9" />
              <path d="M9.6 9a2.5 2.5 0 1 1 4.1 1.9c-1 .8-1.7 1.2-1.7 2.6" />
              <path d="M12 17h.01" />
            </svg>
            <span>{{ t('home.nav.help') }}</span>
          </a>
          <span class="home-status" :title="t('home.nav.statusHint')">
            <span class="home-status-dot" aria-hidden="true"></span>
            <span>{{ t('home.nav.status') }}</span>
          </span>
          <LanguageToggle />
          <ThemeToggle />
        </nav>
      </div>
    </header>

    <section class="home-hero" aria-labelledby="home-title">
      <div class="home-hero-copy">
        <span class="home-hero-eyebrow">{{ t('home.heroEyebrow') }}</span>
        <h1 id="home-title" class="home-hero-title">{{ t('home.heroTitle') }}</h1>
        <p class="home-hero-subtitle">{{ t('home.heroSub') }}</p>
      </div>
      <div class="home-hero-art" aria-hidden="true">
        <div class="home-art-ribbon home-art-ribbon--back"></div>
        <div class="home-art-ribbon home-art-ribbon--front"></div>
        <div class="home-art-badge">
          <svg viewBox="0 0 80 80" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 24h40l10 16-30 28L10 40z" />
            <path d="M36 24 28 40l12 28 12-28-8-16" />
            <path d="M10 40h60" />
          </svg>
        </div>
      </div>
    </section>

    <main class="u-main home-content max-w-6xl">
      <section class="home-quick-redeem" aria-labelledby="quick-redeem-title">
        <div class="home-quick-intro">
          <span class="home-quick-icon" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 7.5A2.5 2.5 0 0 0 6.5 5h11A2.5 2.5 0 0 0 20 7.5v9a2.5 2.5 0 0 0-2.5 2.5h-11A2.5 2.5 0 0 0 4 16.5z" />
              <path d="M8 5v14M16 5v14" />
            </svg>
          </span>
          <div>
            <h2 id="quick-redeem-title">{{ t('home.quickRedeem.title') }}</h2>
            <p>{{ t('home.quickRedeem.hint') }}</p>
          </div>
        </div>
        <form class="home-quick-form" @submit.prevent="goRedeem">
          <label for="quick-cdk" class="sr-only">{{ t('home.quickRedeem.label') }}</label>
          <input
            id="quick-cdk"
            v-model="quickCode"
            class="home-code-input mono"
            :placeholder="t('home.quickRedeem.placeholder')"
            autocomplete="off"
            spellcheck="false"
          />
          <button type="submit" class="home-quick-button">
            <span>{{ t('home.quickRedeem.cta') }}</span>
            <span aria-hidden="true">→</span>
          </button>
        </form>
      </section>

      <section class="home-trust-row" :aria-label="t('home.trust.label')">
        <div v-for="item in trustItems" :key="item.key" class="home-trust-item">
          <span class="home-trust-icon" v-html="item.icon" aria-hidden="true"></span>
          <span>
            <strong>{{ item.title }}</strong>
            <small>{{ item.description }}</small>
          </span>
        </div>
      </section>

      <section class="home-services-grid" :aria-label="t('home.servicesLabel')">
        <router-link
          v-for="svc in services"
          :key="svc.to"
          :to="svc.to"
          class="home-service-card card-hover"
        >
          <span class="home-service-icon" v-html="svc.icon" aria-hidden="true"></span>
          <span class="home-service-body">
            <span class="home-service-heading">
              <span>
                <strong>{{ svc.title }}</strong>
                <small>{{ svc.en }}</small>
              </span>
              <span class="home-service-arrow" aria-hidden="true">›</span>
            </span>
            <span class="home-service-summary">{{ svc.points.join(' · ') }}</span>
          </span>
        </router-link>
      </section>

      <div class="home-guide-link-wrap">
        <a
          class="home-guide-link"
          href="https://docs.qq.com/doc/p/930d2d6ba228c404b7fbfb15c971c079b6001d03?nlc=1"
          target="_blank"
          rel="noopener noreferrer"
          :aria-label="t('home.visualGuide.label')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <circle cx="5" cy="12" r="2" />
            <circle cx="19" cy="6" r="2" />
            <circle cx="19" cy="18" r="2" />
            <path d="m7 11 10-4M7 13l10 4" />
          </svg>
          <span>{{ t('home.visualGuide.label') }}</span>
          <span aria-hidden="true">→</span>
        </a>
        <span>{{ t('home.visualGuide.hint') }}</span>
      </div>

      <section id="operation-guide" class="operation-guide" aria-labelledby="operation-guide-title">
        <div class="home-guide-panel card">
          <div class="home-guide-heading">
            <h2 id="operation-guide-title">{{ t('home.flowTitle') }}</h2>
            <p>{{ t('home.flowSub') }}</p>
          </div>

          <div class="operation-flow-tabs" role="tablist" :aria-label="t('home.flowTitle')">
            <button
              v-for="flow in flowItems"
              :id="'flow-tab-' + flow.key"
              :key="flow.key"
              type="button"
              role="tab"
              class="operation-flow-tab"
              :class="{ 'is-active': activeFlow === flow.key }"
              :aria-selected="activeFlow === flow.key"
              :aria-controls="'flow-panel-' + flow.key"
              @click="activeFlow = flow.key"
            >
              {{ flow.title }}
            </button>
          </div>

          <div
            v-if="activeFlowItem"
            :id="'flow-panel-' + activeFlowItem.key"
            class="operation-flow-panel"
            role="tabpanel"
            :aria-labelledby="'flow-tab-' + activeFlowItem.key"
          >
            <div class="operation-flow-panel-heading">
              <div>
                <h3>{{ activeFlowItem.title }}</h3>
                <p>{{ activeFlowItem.en }}</p>
              </div>
              <span class="operation-flow-badge">{{ activeFlowItem.badge }}</span>
            </div>

            <ol class="operation-flow-steps">
              <li v-for="(step, index) in activeFlowItem.steps" :key="step.title" class="operation-flow-step">
                <div class="operation-flow-step-head">
                  <span class="operation-flow-step-dot">{{ index + 1 }}</span>
                  <span v-if="index < activeFlowItem.steps.length - 1" class="operation-flow-step-line" aria-hidden="true"></span>
                </div>
                <h4>{{ step.title }}</h4>
                <p>{{ step.detail }}</p>
              </li>
            </ol>

            <div class="operation-flow-footer">
              <div>
                <strong>{{ t('home.flowReady') }}</strong>
                <span class="operation-flow-prepare">{{ activeFlowItem.prepare }}</span>
              </div>
              <router-link :to="activeFlowItem.to" class="btn-primary home-flow-cta">
                {{ activeFlowItem.cta }} <span aria-hidden="true">→</span>
              </router-link>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import LanguageToggle from '../components/LanguageToggle.vue'
import ThemeToggle from '../components/ThemeToggle.vue'
import { siteBrand } from '../theme'

const { t } = useI18n({ useScope: 'global' })
const router = useRouter()
const brand = siteBrand
const quickCode = ref('')

type FlowKey = 'submit' | 'batch' | 'lookup' | 'billing'

const flowSpecs: Array<{ key: FlowKey; service: string; flow: FlowKey; to: string }> = [
  { key: 'submit', service: 'recharge', flow: 'submit', to: '/recharge' },
  { key: 'batch', service: 'batch', flow: 'batch', to: '/batch' },
  { key: 'lookup', service: 'lookup', flow: 'lookup', to: '/history' },
  { key: 'billing', service: 'billing', flow: 'billing', to: '/billing' },
]

const activeFlow = ref<FlowKey>('submit')

function goRedeem() {
  const code = quickCode.value.trim()
  router.push({ path: '/recharge', query: code ? { code } : undefined })
}

const trustItems = computed(() => [
  {
    key: 'safe',
    title: t('home.trust.safe.title'),
    description: t('home.trust.safe.description'),
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3 20 6v5c0 5.4-3.4 8.8-8 10-4.6-1.2-8-4.6-8-10V6z"/><path d="m8.5 12 2.2 2.2 4.8-5"/></svg>',
  },
  {
    key: 'track',
    title: t('home.trust.track.title'),
    description: t('home.trust.track.description'),
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="m8 4-4 4 4 4"/><path d="M4 8h9a6 6 0 0 1 6 6v4"/><path d="m16 18 3 3 3-3"/><path d="M16 4h.01M20 4h.01"/></svg>',
  },
  {
    key: 'official',
    title: t('home.trust.official.title'),
    description: t('home.trust.official.description'),
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3 2.2 2.1 3-.2.9 2.9 2.6 1.5-1.5 2.6.2 3-2.9.9L15 19l-3 2-3-2-.5-3.2-2.9-.9.2-3-1.5-2.6L7 7.8l.9-2.9 3 .2z"/><path d="m9 12 2 2 4-4"/></svg>',
  },
])

const services = computed(() => [
  {
    to: '/recharge',
    title: t('home.services.recharge.title'),
    en: t('home.services.recharge.en'),
    cta: t('home.services.recharge.cta'),
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6"/><path d="m9 14 2 2 4-4"/></svg>',
    points: [t('home.services.recharge.p1'), t('home.services.recharge.p2'), t('home.services.recharge.p3')],
  },
  {
    to: '/batch',
    title: t('home.services.batch.title'),
    en: t('home.services.batch.en'),
    cta: t('home.services.batch.cta'),
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>',
    points: [t('home.services.batch.p1'), t('home.services.batch.p2'), t('home.services.batch.p3')],
  },
  {
    to: '/history',
    title: t('home.services.lookup.title'),
    en: t('home.services.lookup.en'),
    cta: t('home.services.lookup.cta'),
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>',
    points: [t('home.services.lookup.p1'), t('home.services.lookup.p2'), t('home.services.lookup.p3')],
  },
  {
    to: '/billing',
    title: t('home.services.billing.title'),
    en: t('home.services.billing.en'),
    cta: t('home.services.billing.cta'),
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M6 20V14M12 20V4M18 20v-9"/></svg>',
    points: [t('home.services.billing.p1'), t('home.services.billing.p2'), t('home.services.billing.p3')],
  },
])

const flowItems = computed(() => flowSpecs.map((spec) => ({
  key: spec.key,
  to: spec.to,
  title: t('home.services.' + spec.service + '.title'),
  en: t('home.services.' + spec.service + '.en'),
  cta: t('home.services.' + spec.service + '.cta'),
  badge: t('home.flowBadge.' + spec.key),
  prepare: t('home.flowPrepare.' + spec.key),
  steps: [
    { title: t('home.services.' + spec.service + '.p1'), detail: t('home.flows.' + spec.flow + '.s1') },
    { title: t('home.services.' + spec.service + '.p2'), detail: t('home.flows.' + spec.flow + '.s2') },
    { title: t('home.services.' + spec.service + '.p3'), detail: t('home.flows.' + spec.flow + '.s3') },
    { title: t('home.flows.' + spec.flow + '.s4Title'), detail: t('home.flows.' + spec.flow + '.s4') },
  ],
})))

const activeFlowItem = computed(() => flowItems.value.find((flow) => flow.key === activeFlow.value) ?? flowItems.value[0])
</script>

<style scoped>
.home-page {
  --home-blue: #2925a9;
  --home-blue-soft: rgba(53, 48, 200, 0.1);
  --home-hero-glow: rgba(93, 84, 232, 0.14);
  overflow: hidden;
}

.home-topbar {
  position: relative;
  z-index: 2;
  background: color-mix(in srgb, var(--surface-solid, var(--surface)) 94%, transparent);
  border-bottom: 1px solid color-mix(in srgb, var(--brd) 80%, transparent);
}

.home-topbar-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  max-width: 1180px;
  min-height: 68px;
  margin: 0 auto;
  padding: 0 28px;
}

.home-brand {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  gap: 10px;
  color: var(--ink);
}

.home-logo {
  display: grid;
  width: 38px;
  height: 38px;
  flex: none;
  place-items: center;
  border-radius: 50%;
  background: var(--logo-bg, var(--primary));
  color: var(--logo-fg, var(--primary-on));
  box-shadow: var(--logo-shadow, 0 8px 18px -10px var(--primary));
}

.home-logo svg { width: 22px; height: 22px; }

.home-brand-copy {
  display: inline-flex;
  min-width: 0;
  align-items: baseline;
  gap: 10px;
  line-height: 1.2;
}

.home-brand-copy strong {
  overflow: hidden;
  color: var(--ink);
  font-size: 17px;
  font-weight: 800;
  letter-spacing: -0.03em;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.home-brand-copy small {
  overflow: hidden;
  color: var(--ink-2);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.home-brand-divider { color: var(--brd-2); }

.home-nav {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--ink-2);
}

.home-nav-link,
.home-status {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  min-height: 36px;
  padding: 0 10px;
  border-radius: 999px;
  color: var(--ink-2);
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
}

.home-nav-link:hover { color: var(--ink); background: var(--primary-soft); }
.home-nav-link svg { width: 17px; height: 17px; }

.home-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #1aa65b;
  box-shadow: 0 0 0 3px rgba(26, 166, 91, 0.14);
}

.home-nav :deep(.btn-ghost) { min-height: 36px; color: var(--ink-2); }

.home-hero {
  position: relative;
  display: flex;
  align-items: center;
  max-width: 1180px;
  min-height: 270px;
  margin: 0 auto;
  padding: 58px 72px 48px;
  isolation: isolate;
}

.home-hero::before {
  position: absolute;
  z-index: -2;
  inset: 0 -16vw;
  background:
    radial-gradient(circle at 72% 35%, var(--home-hero-glow), transparent 28%),
    linear-gradient(180deg, color-mix(in srgb, var(--surface) 55%, transparent), transparent 80%);
  content: '';
  pointer-events: none;
}

.home-hero-copy { position: relative; z-index: 1; max-width: 620px; }

.home-hero-eyebrow {
  display: inline-flex;
  padding: 6px 16px;
  border-radius: 999px;
  color: var(--home-blue);
  background: var(--home-blue-soft);
  font-size: 13px;
  font-weight: 750;
}

.home-hero-title {
  margin-top: 18px;
  color: var(--ink);
  font-family: var(--font-display, var(--font-sans));
  font-size: clamp(2.4rem, 5.3vw, 4.45rem);
  font-weight: 850;
  letter-spacing: -0.065em;
  line-height: 1.02;
  text-wrap: balance;
}

.home-hero-subtitle {
  max-width: 520px;
  margin-top: 16px;
  color: var(--ink-2);
  font-size: clamp(1rem, 1.5vw, 1.2rem);
  line-height: 1.65;
}

.home-hero-art {
  position: absolute;
  z-index: -1;
  top: 0;
  right: -10px;
  width: 52%;
  height: 100%;
  overflow: hidden;
  pointer-events: none;
}

.home-art-ribbon {
  position: absolute;
  right: -16%;
  width: 120%;
  height: 74px;
  transform: rotate(-19deg);
  border-radius: 999px;
  background: linear-gradient(100deg, transparent 0%, color-mix(in srgb, var(--home-blue) 22%, transparent) 45%, color-mix(in srgb, #8a8cf4 58%, transparent) 100%);
}

.home-art-ribbon--back { top: 58%; opacity: 0.32; }
.home-art-ribbon--front { top: 43%; height: 90px; opacity: 0.62; transform: rotate(-23deg); }

.home-art-badge {
  position: absolute;
  top: 19%;
  right: 21%;
  display: grid;
  width: 170px;
  height: 138px;
  place-items: center;
  transform: rotate(-27deg) skewY(3deg);
  border-radius: 28px;
  background: linear-gradient(145deg, #1a1a26 5%, #090913 72%);
  color: #ffffff;
  box-shadow: 16px 28px 36px -24px rgba(15, 17, 45, 0.75), inset 12px 10px 24px -24px #ffffff;
}

.home-art-badge::after {
  position: absolute;
  inset: 9px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 22px;
  content: '';
}

.home-art-badge svg { width: 82px; height: 82px; transform: rotate(27deg); }

.home-content { position: relative; z-index: 1; padding-bottom: 48px; }

.home-quick-redeem {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 28px;
  padding: 28px 34px;
  border: 1px solid var(--brd);
  border-radius: 22px;
  background: color-mix(in srgb, var(--surface-solid, var(--surface)) 96%, transparent);
  box-shadow: 0 22px 52px -38px rgba(35, 40, 80, 0.52);
}

.home-quick-intro { display: flex; min-width: 240px; align-items: center; gap: 16px; }
.home-quick-icon,
.home-service-icon {
  display: grid;
  flex: none;
  place-items: center;
  color: var(--home-blue);
  background: var(--home-blue-soft);
}

.home-quick-icon { width: 52px; height: 52px; border-radius: 16px; }
.home-quick-icon svg { width: 28px; height: 28px; }
.home-quick-intro h2 { color: var(--ink); font-size: 20px; font-weight: 800; letter-spacing: -0.03em; }
.home-quick-intro p { margin-top: 3px; color: var(--ink-2); font-size: 13px; }

.home-quick-form { display: flex; min-width: 0; flex: 1; align-items: center; gap: 14px; }
.home-code-input {
  min-width: 0;
  height: 52px;
  flex: 1;
  padding: 0 18px;
  border: 1px solid var(--brd-2);
  border-radius: 12px;
  color: var(--ink);
  background: var(--input-bg, var(--surface-2));
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}

.home-code-input::placeholder { color: var(--ink-3); }
.home-code-input:focus { border-color: var(--home-blue); background: var(--surface-solid, var(--surface)); box-shadow: 0 0 0 4px color-mix(in srgb, var(--home-blue) 13%, transparent); }

.home-quick-button {
  display: inline-flex;
  height: 52px;
  flex: none;
  align-items: center;
  gap: 12px;
  padding: 0 24px;
  border: 0;
  border-radius: 12px;
  color: #ffffff;
  background: linear-gradient(135deg, #3834bf, #222084);
  box-shadow: 0 15px 26px -16px rgba(44, 40, 168, 0.9);
  font-size: 15px;
  font-weight: 800;
  transition: transform 0.18s ease, box-shadow 0.18s ease, filter 0.18s ease;
}

.home-quick-button:hover { filter: brightness(1.08); box-shadow: 0 18px 30px -15px rgba(44, 40, 168, 0.95); transform: translateY(-1px); }
.home-quick-button:active { transform: translateY(0); }
.home-quick-button:focus-visible { outline: 2px solid var(--home-blue); outline-offset: 3px; }

.home-trust-row {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  margin: 26px auto 34px;
}

.home-trust-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  min-height: 54px;
  padding: 0 24px;
  border-right: 1px solid var(--brd);
}

.home-trust-item:last-child { border-right: 0; }
.home-trust-icon { width: 31px; height: 31px; color: var(--home-blue); }
.home-trust-icon svg { width: 100%; height: 100%; }
.home-trust-item strong { display: block; color: var(--ink); font-size: 15px; font-weight: 800; }
.home-trust-item small { display: block; margin-top: 2px; color: var(--ink-2); font-size: 12px; }

.home-services-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.home-service-card {
  display: flex;
  min-width: 0;
  min-height: 132px;
  align-items: flex-start;
  gap: 14px;
  padding: 20px 18px;
  border: 1px solid var(--brd);
  border-radius: 16px;
  color: inherit;
  background: var(--surface-solid, var(--surface));
  box-shadow: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.home-service-card:hover { border-color: color-mix(in srgb, var(--home-blue) 35%, var(--brd)); box-shadow: 0 18px 34px -28px rgba(38, 36, 145, 0.8); transform: translateY(-2px); }
.home-service-icon { width: 42px; height: 42px; border-radius: 13px; }
.home-service-icon svg { width: 23px; height: 23px; }
.home-service-body { display: block; min-width: 0; flex: 1; }
.home-service-heading { display: flex; align-items: flex-start; justify-content: space-between; gap: 8px; }
.home-service-heading strong { display: block; color: var(--ink); font-size: 16px; font-weight: 800; }
.home-service-heading small { display: block; margin-top: 2px; color: var(--ink-2); font-size: 11px; }
.home-service-arrow { color: var(--ink-3); font-size: 25px; line-height: 0.8; transition: color 0.2s ease, transform 0.2s ease; }
.home-service-card:hover .home-service-arrow { color: var(--home-blue); transform: translateX(2px); }
.home-service-summary { display: block; margin-top: 17px; color: var(--ink-2); font-size: 12px; line-height: 1.5; }

.home-guide-link-wrap { display: flex; flex-direction: column; align-items: center; gap: 8px; margin: 42px 0 28px; color: var(--ink-2); font-size: 12px; }
.home-guide-link { display: inline-flex; align-items: center; gap: 10px; padding: 10px 20px; border: 1px solid var(--brd-2); border-radius: 999px; color: var(--ink); background: var(--surface-solid, var(--surface)); font-size: 14px; font-weight: 750; transition: border-color 0.2s ease, color 0.2s ease, transform 0.2s ease; }
.home-guide-link:hover { border-color: var(--home-blue); color: var(--home-blue); transform: translateY(-1px); }
.home-guide-link svg { width: 17px; height: 17px; }

.operation-guide { border-top: 1px solid var(--brd); padding-top: 42px; }
.home-guide-panel { padding: 34px 34px 30px; }
.home-guide-heading { text-align: center; }
.home-guide-heading h2 { color: var(--ink); font-size: 30px; font-weight: 850; letter-spacing: -0.045em; }
.home-guide-heading p { margin-top: 7px; color: var(--ink-2); font-size: 14px; }

.operation-flow-tabs { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 4px; margin-top: 28px; padding: 4px; border: 1px solid var(--brd); border-radius: 999px; background: var(--tabs-bg, var(--surface-2)); }
.operation-flow-tab { min-height: 42px; padding: 0.65rem 0.75rem; border-radius: 999px; color: var(--ink-2); font-size: 0.875rem; font-weight: 650; transition: background-color 0.18s ease, color 0.18s ease, box-shadow 0.18s ease; }
.operation-flow-tab:hover { color: var(--ink); background: color-mix(in srgb, var(--home-blue) 8%, transparent); }
.operation-flow-tab.is-active { color: #ffffff; background: #0d0d0d; box-shadow: 0 7px 18px -10px rgba(0, 0, 0, 0.45); }
.operation-flow-panel { margin-top: 18px; padding: 28px 30px 24px; border: 1px solid var(--brd); border-radius: 20px; background: var(--surface-solid, var(--surface)); }
.operation-flow-panel-heading { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; }
.operation-flow-panel-heading h3 { color: var(--ink); font-size: 21px; font-weight: 800; }
.operation-flow-panel-heading p { margin-top: 3px; color: var(--ink-2); font-size: 12px; }
.operation-flow-badge, .operation-flow-prepare { display: inline-flex; align-items: center; width: fit-content; border: 1px solid var(--brd); border-radius: 999px; color: var(--ink-2); background: var(--surface-2); font-size: 0.75rem; line-height: 1.25; }
.operation-flow-badge { padding: 0.4rem 0.7rem; }
.operation-flow-prepare { padding: 0.5rem 0.75rem; }
.operation-flow-steps { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 18px; margin-top: 28px; list-style: none; padding: 0; }
.operation-flow-step-head { display: flex; min-height: 32px; align-items: center; }
.operation-flow-step-dot { position: relative; z-index: 1; display: grid; width: 32px; height: 32px; flex: none; place-items: center; border-radius: 50%; color: #ffffff; background: #0d0d0d; font-size: 0.8rem; font-weight: 750; }
.operation-flow-step-line { height: 2px; flex: 1; margin-left: 10px; border-radius: 999px; background: var(--brd-2); }
.operation-flow-step h4 { margin-top: 14px; color: var(--ink); font-size: 14px; font-weight: 800; }
.operation-flow-step p { margin-top: 6px; color: var(--ink-2); font-size: 12px; line-height: 1.65; }
.operation-flow-footer { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-top: 28px; padding-top: 20px; border-top: 1px solid var(--brd); }
.operation-flow-footer > div { display: flex; align-items: center; gap: 10px; color: var(--ink); font-size: 13px; }
.home-flow-cta { flex: none; padding: 0.7rem 1.1rem; }

@media (max-width: 960px) {
  .home-hero { padding-right: 38px; padding-left: 38px; }
  .home-hero-art { right: -17%; width: 60%; opacity: 0.8; }
  .home-services-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

@media (max-width: 720px) {
  .home-topbar-inner { min-height: 60px; padding: 0 16px; }
  .home-brand-copy small, .home-brand-divider, .home-nav-link, .home-status { display: none; }
  .home-nav { gap: 0; }
  .home-hero { min-height: 310px; padding: 42px 20px 38px; }
  .home-hero-copy { max-width: 78%; }
  .home-hero-title { font-size: clamp(2.3rem, 11vw, 3.45rem); }
  .home-hero-art { right: -42%; width: 90%; opacity: 0.58; }
  .home-art-badge { right: 22%; width: 128px; height: 108px; }
  .home-art-badge svg { width: 62px; height: 62px; }
  .home-content { padding: 0 16px 40px; }
  .home-quick-redeem { align-items: stretch; flex-direction: column; gap: 18px; padding: 22px 18px; }
  .home-quick-form { flex-direction: column; align-items: stretch; }
  .home-quick-button { justify-content: center; }
  .home-trust-row { gap: 12px; grid-template-columns: 1fr; margin: 22px 0 28px; }
  .home-trust-item { justify-content: flex-start; min-height: 42px; padding: 0 8px; border-right: 0; }
  .home-services-grid { grid-template-columns: 1fr; gap: 12px; }
  .home-service-card { min-height: 108px; }
  .home-guide-panel { padding: 24px 16px; }
  .home-guide-heading h2 { font-size: 26px; }
  .operation-flow-tabs { grid-template-columns: repeat(2, minmax(0, 1fr)); border-radius: 16px; }
  .operation-flow-tab { border-radius: 12px; }
  .operation-flow-panel { padding: 22px 16px; }
  .operation-flow-panel-heading { flex-direction: column; }
  .operation-flow-steps { grid-template-columns: 1fr; gap: 20px; }
  .operation-flow-step { display: grid; grid-template-columns: 42px minmax(0, 1fr); column-gap: 8px; }
  .operation-flow-step-head { grid-row: span 2; align-items: flex-start; }
  .operation-flow-step-line { display: none; }
  .operation-flow-step h4 { margin-top: 2px; }
  .operation-flow-step p { margin-top: 2px; }
  .operation-flow-footer { align-items: stretch; flex-direction: column; }
  .operation-flow-footer > div { align-items: flex-start; flex-direction: column; }
  .home-flow-cta { justify-content: center; width: 100%; }
}

@media (prefers-reduced-motion: reduce) {
  .home-service-card, .home-quick-button, .home-guide-link { transition: none; }
}

</style>
