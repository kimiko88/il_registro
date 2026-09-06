<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card class="keyboard-shortcuts-card q-pa-md" style="width: min(700px, 95vw); max-width: 95vw; border-radius: 16px;">
      <q-card-section class="row items-center justify-between q-pb-none">
        <div class="text-h6 text-weight-bold flex items-center gap-sm text-primary">
          <q-icon name="keyboard" size="28px" />
          {{ t('a11y.shortcutsTitle') || 'Scorciatoie da Tastiera & Accessibilità' }}
        </div>
        <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-pt-sm">
        <div class="text-caption text-grey-7 q-mb-md">
          {{ t('a11y.shortcutsSubtitle') || 'Usa la combinazione di tasti per navigare rapidamente nel registro senza utilizzare il mouse.' }}
        </div>

        <q-input
          v-model="filterQuery"
          dense
          outlined
          :placeholder="t('a11y.shortcutsSearch') || 'Cerca comando o scorciatoia...'"
          class="q-mb-md"
          clearable
        >
          <template #prepend>
            <q-icon name="search" />
          </template>
        </q-input>

        <div class="shortcuts-list q-gutter-y-sm" style="max-height: 400px; overflow-y: auto;">
          <div v-for="(group, idx) in filteredGroups" :key="idx" class="q-mb-md">
            <div class="text-subtitle2 text-weight-bold text-grey-8 q-mb-xs">
              {{ group.category }}
            </div>
            <q-list bordered separator class="rounded-borders bg-grey-1">
              <q-item v-for="(item, itemIdx) in group.items" :key="itemIdx" class="items-center">
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ item.description }}</q-item-label>
                  <q-item-label caption v-if="item.hint">{{ item.hint }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="flex gap-xs items-center">
                    <kbd v-for="key in item.keys" :key="key" class="shortcut-key">{{ key }}</kbd>
                  </div>
                </q-item-section>
              </q-item>
            </q-list>
          </div>
        </div>
      </q-card-section>

      <q-card-actions align="right" class="q-pt-none">
        <q-btn flat :label="(t('common.close') || 'Chiudi') + ' (Esc)'" color="primary" v-close-popup no-caps />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useThemeStore } from 'src/stores/theme'

const { t, locale } = useI18n()
const themeStore = useThemeStore()

const isOpen = computed({
  get: () => themeStore.keyboardShortcutsHelpOpen,
  set: (val) => themeStore.toggleKeyboardShortcutsHelp(val)
})

const filterQuery = ref('')

const shortcutGroups = computed(() => {
  const _ = locale && typeof locale === 'object' ? locale.value : locale
  return [
    {
      category: t('a11y.shortcutsCatGlobal') || 'Navigazione Globale (Ruolo Docente & Segreteria)',
    items: [
      { description: t('a11y.shortcutDashboard') || 'Vai alla Dashboard', keys: ['Alt', '1'] },
      { description: t('a11y.shortcutGrades') || 'Vai al Registro Voti', keys: ['Alt', 'V'] },
      { description: t('a11y.shortcutAttendance') || 'Vai al Registro Presenze', keys: ['Alt', 'P'] },
      { description: t('a11y.shortcutAgenda') || 'Vai all\'Agenda & Compiti', keys: ['Alt', 'A'] },
      { description: t('a11y.shortcutSearch') || 'Attiva Ricerca Globale', keys: ['Alt', 'S'] }
    ]
  },
  {
    category: t('a11y.shortcutsCatA11y') || 'Funzionalità di Accessibilità (A11y)',
    items: [
      { description: t('a11y.shortcutHelp') || 'Mostra / Nascondi Guida Scorciatoie', keys: ['?'] },
      { description: t('a11y.shortcutRuler') || 'Attiva / Disattiva Righello di Lettura', keys: ['Alt', 'R'] },
      { description: t('a11y.shortcutRulerMove') || 'Sposta Righello Su / Giù', keys: ['Alt', '↑ / ↓'] },
      { description: t('a11y.shortcutTts') || 'Attiva / Disattiva Lettura Vocale (TTS)', keys: ['Alt', 'T'] },
      { description: t('a11y.shortcutSkip') || 'Salta al Contenuto Principale (Skip to Main)', keys: [t('a11y.keyTabStart') || 'Tab (Inizio Pagina)'] },
      { description: t('a11y.shortcutEsc') || 'Chiudi Modale / Annulla Operazione', keys: ['Esc'] }
    ]
  }
]})

const filteredGroups = computed(() => {
  const groups = shortcutGroups.value
  if (!filterQuery.value) return groups
  const q = filterQuery.value.toLowerCase()
  return groups
    .map(g => ({
      ...g,
      items: g.items.filter(i =>
        i.description.toLowerCase().includes(q) ||
        i.keys.some(k => k.toLowerCase().includes(q))
      )
    }))
    .filter(g => g.items.length > 0)
})
</script>

<style scoped>
.shortcut-key {
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  border-bottom: 2px solid #94a3b8;
  border-radius: 6px;
  padding: 3px 8px;
  font-family: monospace;
  font-size: 0.85rem;
  font-weight: 700;
  color: #1e293b;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}
</style>
