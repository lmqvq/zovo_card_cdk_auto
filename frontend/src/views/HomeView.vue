<template>
  <div class="u-page">
    <UserPageHeader :eyebrow="t('hero.redeemEyebrow')" :title="t('home.heroTitle')" :subtitle="t('home.heroSub')" badges />
    <div class="u-main max-w-6xl">
      <!-- Main Services -->
      <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6 animate-slideInUp">
        <router-link
          v-for="svc in services"
          :key="svc.to"
          :to="svc.to"
          class="card card-hover group svc-card"
        >
          <span class="grid h-12 w-12 place-items-center rounded-2xl" style="background: var(--primary-soft); color: var(--primary)" v-html="svc.icon" />
          <h3 class="mt-4 text-xl font-bold text-ink">{{ svc.title }}</h3>
          <p class="mt-1 text-sm text-subtle">{{ svc.en }}</p>
          <div class="mt-5 space-y-2 text-sm text-muted">
            <div v-for="line in svc.points" :key="line" class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 shrink-0" style="color: var(--primary)" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12" /></svg>
              <span>{{ line }}</span>
            </div>
          </div>
          <div class="mt-6 text-sm font-medium app-link">{{ svc.cta }} →</div>
        </router-link>
      </div>

      <div class="mt-10 flex flex-col items-center gap-2">
        <a
          class="btn-secondary inline-flex items-center gap-2 rounded-full !px-5 !py-2 text-sm"
          href="https://docs.qq.com/doc/p/930d2d6ba228c404b7fbfb15c971c079b6001d03?nlc=1"
          target="_blank"
          rel="noopener noreferrer"
          :aria-label="t('home.visualGuide.label')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <circle cx="5" cy="12" r="2" />
            <circle cx="19" cy="6" r="2" />
            <circle cx="19" cy="18" r="2" />
            <path d="m7 11 10-4M7 13l10 4" />
          </svg>
          <span>{{ t('home.visualGuide.label') }}</span>
          <span aria-hidden="true">→</span>
        </a>
        <span class="text-xs text-subtle">{{ t('home.visualGuide.hint') }}</span>
      </div>

      <!-- Flow Section -->
      <section id="operation-guide" class="operation-guide mt-16 border-t bd pt-12" aria-labelledby="operation-guide-title">
        <div class="card !p-5 sm:!p-8">
          <div class="text-center">
            <h3 id="operation-guide-title" class="text-2xl font-bold text-ink">{{ t('home.flowTitle') }}</h3>
            <p class="mt-2 text-sm text-muted">{{ t('home.flowSub') }}</p>
          </div>

          <div class="operation-flow-tabs mt-7 grid grid-cols-2 gap-1 rounded-2xl p-1 sm:grid-cols-4" role="tablist" :aria-label="t('home.flowTitle')">
            <button
              v-for="flow in flowItems"
              :id="`flow-tab-${flow.key}`"
              :key="flow.key"
              type="button"
              role="tab"
              class="operation-flow-tab"
              :class="{ 'is-active': activeFlow === flow.key }"
              :aria-selected="activeFlow === flow.key"
              :aria-controls="`flow-panel-${flow.key}`"
              @click="activeFlow = flow.key"
            >
              {{ flow.title }}
            </button>
          </div>

          <div
            v-if="activeFlowItem"
            :id="`flow-panel-${activeFlowItem.key}`"
            class="operation-flow-panel mt-5 rounded-3xl border bd p-5 sm:p-8"
            role="tabpanel"
            :aria-labelledby="`flow-tab-${activeFlowItem.key}`"
          >
            <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div>
                <h4 class="text-xl font-bold text-ink">{{ activeFlowItem.title }}</h4>
                <p class="mt-1 text-sm text-subtle">{{ activeFlowItem.en }}</p>
              </div>
              <span class="operation-flow-badge">{{ activeFlowItem.badge }}</span>
            </div>

            <ol class="operation-flow-steps mt-8 grid gap-7 md:grid-cols-3 md:gap-5">
              <li v-for="(step, index) in activeFlowItem.steps" :key="step.title" class="operation-flow-step">
                <div class="operation-flow-step-head">
                  <span class="operation-flow-step-dot">{{ index + 1 }}</span>
                  <span v-if="index < activeFlowItem.steps.length - 1" class="operation-flow-step-line" aria-hidden="true" />
                </div>
                <h5 class="mt-4 font-bold text-ink">{{ step.title }}</h5>
                <p class="mt-2 text-sm leading-6 text-muted">{{ step.detail }}</p>
              </li>
            </ol>

            <div class="operation-flow-footer mt-8 flex flex-col gap-5 border-t bd pt-5 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-3">
                <span class="text-sm font-semibold text-ink">{{ t('home.flowReady') }}</span>
                <span class="operation-flow-prepare">{{ activeFlowItem.prepare }}</span>
              </div>
              <router-link :to="activeFlowItem.to" class="btn-primary shrink-0">
                {{ activeFlowItem.cta }} <span aria-hidden="true">→</span>
              </router-link>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import UserPageHeader from '../components/UserPageHeader.vue'

const { t } = useI18n({ useScope: 'global' })
type FlowKey = 'submit' | 'batch' | 'lookup' | 'billing'

const flowSpecs: Array<{ key: FlowKey; service: string; flow: FlowKey; to: string }> = [
  { key: 'submit', service: 'recharge', flow: 'submit', to: '/recharge' },
  { key: 'batch', service: 'batch', flow: 'batch', to: '/batch' },
  { key: 'lookup', service: 'lookup', flow: 'lookup', to: '/history' },
  { key: 'billing', service: 'billing', flow: 'billing', to: '/billing' },
]

const activeFlow = ref<FlowKey>('submit')

const services = computed(() => [
  {
    to: '/recharge',
    title: t('home.services.recharge.title'),
    en: t('home.services.recharge.en'),
    cta: t('home.services.recharge.cta'),
    icon: '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6"/><path d="M9 14l2 2 4-4"/></svg>',
    points: [t('home.services.recharge.p1'), t('home.services.recharge.p2'), t('home.services.recharge.p3')],
  },
  {
    to: '/batch',
    title: t('home.services.batch.title'),
    en: t('home.services.batch.en'),
    cta: t('home.services.batch.cta'),
    icon: '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>',
    points: [t('home.services.batch.p1'), t('home.services.batch.p2'), t('home.services.batch.p3')],
  },
  {
    to: '/history',
    title: t('home.services.lookup.title'),
    en: t('home.services.lookup.en'),
    cta: t('home.services.lookup.cta'),
    icon: '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>',
    points: [t('home.services.lookup.p1'), t('home.services.lookup.p2'), t('home.services.lookup.p3')],
  },
  {
    to: '/billing',
    title: t('home.services.billing.title'),
    en: t('home.services.billing.en'),
    cta: t('home.services.billing.cta'),
    icon: '<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>',
    points: [t('home.services.billing.p1'), t('home.services.billing.p2'), t('home.services.billing.p3')],
  },
])

const flowItems = computed(() => flowSpecs.map((spec) => ({
  key: spec.key,
  to: spec.to,
  title: t(`home.services.${spec.service}.title`),
  en: t(`home.services.${spec.service}.en`),
  cta: t(`home.services.${spec.service}.cta`),
  badge: t(`home.flowBadge.${spec.key}`),
  prepare: t(`home.flowPrepare.${spec.key}`),
  steps: [
    { title: t(`home.services.${spec.service}.p1`), detail: t(`home.flows.${spec.flow}.s1`) },
    { title: t(`home.services.${spec.service}.p2`), detail: t(`home.flows.${spec.flow}.s2`) },
    { title: t(`home.services.${spec.service}.p3`), detail: t(`home.flows.${spec.flow}.s3`) },
  ],
})))

const activeFlowItem = computed(() => flowItems.value.find((flow) => flow.key === activeFlow.value) ?? flowItems.value[0])
</script>

<style scoped>
.operation-flow-tabs {
  background: var(--surface-2);
  border: 1px solid var(--brd);
}

.operation-flow-tab {
  min-height: 42px;
  padding: 0.65rem 0.75rem;
  border-radius: 0.9rem;
  color: var(--ink-2);
  font-size: 0.875rem;
  font-weight: 600;
  transition: background-color 0.18s ease, color 0.18s ease, box-shadow 0.18s ease;
}

.operation-flow-tab:hover {
  color: var(--ink);
  background: color-mix(in srgb, var(--primary) 8%, transparent);
}

.operation-flow-tab.is-active {
  color: var(--primary-on);
  background: var(--primary);
  box-shadow: 0 7px 18px -10px var(--primary);
}

.operation-flow-panel {
  background: var(--surface-solid, var(--surface));
}

.operation-flow-badge,
.operation-flow-prepare {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  border: 1px solid var(--brd);
  border-radius: 999px;
  color: var(--ink-2);
  background: var(--surface-2);
  font-size: 0.75rem;
  line-height: 1.25;
}

.operation-flow-badge {
  padding: 0.4rem 0.7rem;
}

.operation-flow-prepare {
  padding: 0.5rem 0.75rem;
}

.operation-flow-step-head {
  display: flex;
  align-items: center;
  min-height: 32px;
}

.operation-flow-step-dot {
  position: relative;
  z-index: 1;
  display: grid;
  width: 32px;
  height: 32px;
  flex: none;
  place-items: center;
  border: 2px solid var(--primary);
  border-radius: 50%;
  color: var(--primary-on);
  background: var(--primary);
  font-size: 0.8rem;
  font-weight: 700;
}

.operation-flow-step-line {
  height: 2px;
  flex: 1;
  margin-left: 0.75rem;
  border-radius: 999px;
  background: var(--brd-2);
}

@media (max-width: 767px) {
  .operation-flow-step:not(:last-child) {
    padding-bottom: 0.25rem;
    border-bottom: 1px dashed var(--brd);
  }
}
</style>
