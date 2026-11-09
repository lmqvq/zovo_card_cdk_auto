<template>
  <div class="excel-import">
    <input
      ref="fileInputRef"
      type="file"
      accept=".xlsx,.xls,.csv"
      class="hidden"
      @change="onChange"
    />
    <button
      type="button"
      class="excel-drop"
      :disabled="importing"
      :aria-label="t('batchUi.chooseFile')"
      @click="fileInputRef?.click()"
    >
      <span class="excel-icon" aria-hidden="true">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M12 16V4"/><path d="m7 9 5-5 5 5"/><path d="M5 20h14"/></svg>
      </span>
      <span class="excel-copy">
        <b>{{ importing ? t('batch.importing') : sessionPool.length ? t('batch.imported', { n: sessionPool.length }) : t('batchUi.chooseFile') }}</b>
        <span>{{ t('batchUi.fileSub') }}</span>
      </span>
      <span class="excel-cta">{{ t('batchUi.chooseFileBtn') }}</span>
    </button>
    <p v-if="importMsg" class="excel-msg">{{ importMsg }}</p>
    <div v-if="sessionPool.length > 0" class="excel-list">
      <div v-for="(s, i) in sessionPool.slice(0, 8)" :key="i">
        <span>{{ i + 1 }}</span>{{ s.email || '(无邮箱)' }} · session {{ s.session.length }} 字
      </div>
      <div v-if="sessionPool.length > 8">…共 {{ sessionPool.length }} 条</div>
    </div>
    <div class="excel-foot">
      <button v-if="sessionPool.length > 0" type="button" class="excel-clear" @click="$emit('clear')">{{ t('batch.clearImport') }}</button>
      <span class="excel-hint">{{ t('batchUi.excelShort') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { ImportedSession } from '../lib/batch-session'

defineProps<{
  sessionPool: ImportedSession[]
  importMsg: string
  importing: boolean
}>()

const emit = defineEmits<{
  pick: [file: File | null]
  clear: []
}>()

const { t } = useI18n({ useScope: 'global' })
const fileInputRef = ref<HTMLInputElement | null>(null)

function onChange(e: Event) {
  const input = e.target as HTMLInputElement
  emit('pick', input.files?.[0] ?? null)
  input.value = ''
}
</script>

<style scoped>
.excel-import { margin-top: 9px; }
.excel-drop {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
  padding: 16px;
  border: 1.5px dashed var(--brd-2);
  border-radius: 14px;
  background: color-mix(in srgb, var(--surface) 78%, var(--bg));
  color: var(--ink);
  text-align: left;
  cursor: pointer;
  transition: border-color 0.15s ease, background-color 0.15s ease;
}
.excel-drop:hover:not(:disabled) {
  border-color: var(--ink-3);
  background: var(--surface);
}
.excel-drop:disabled { opacity: 0.55; cursor: wait; }
.excel-icon {
  display: grid;
  place-items: center;
  flex: none;
  width: 38px;
  height: 38px;
  border: 1px solid var(--brd);
  border-radius: 11px;
  background: var(--surface);
  color: var(--ink-2);
}
.excel-icon svg { width: 17px; height: 17px; }
.excel-copy { min-width: 0; flex: 1; }
.excel-copy b { display: block; color: var(--ink); font-size: 11.5px; font-weight: 600; }
.excel-copy span { display: block; margin-top: 3px; color: var(--ink-3); font-size: 10px; }
.excel-cta {
  flex: none;
  padding: 7px 12px;
  border: 1px solid var(--brd-2);
  border-radius: 999px;
  background: var(--surface);
  color: var(--ink);
  font-size: 10px;
  font-weight: 600;
}
.excel-msg { margin: 8px 0 0; color: var(--ink-3); font-size: 10.5px; line-height: 1.55; }
.excel-list {
  max-height: 112px;
  margin-top: 9px;
  padding: 9px 11px;
  overflow-y: auto;
  border-radius: 10px;
  background: var(--surface-2);
  color: var(--ink-3);
  font-family: var(--font-mono);
  font-size: 10px;
  line-height: 1.7;
}
.excel-list span { display: inline-block; width: 20px; color: var(--ink-4); }
.excel-foot { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; margin-top: 8px; }
.excel-clear {
  flex: none;
  padding: 0;
  border: 0;
  background: transparent;
  color: var(--err);
  font-size: 10px;
  cursor: pointer;
}
.excel-hint { color: var(--ink-3); font-size: 10px; line-height: 1.55; }
@media (max-width: 640px) {
  .excel-drop { flex-wrap: wrap; padding: 13px; }
  .excel-cta { margin-left: 52px; }
}
</style>
