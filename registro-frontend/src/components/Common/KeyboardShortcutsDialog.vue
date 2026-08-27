<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card class="keyboard-shortcuts-card q-pa-md" style="min-width: 550px; max-width: 700px; border-radius: 16px;">
      <q-card-section class="row items-center justify-between q-pb-none">
        <div class="text-h6 text-weight-bold flex items-center gap-sm text-primary">
          <q-icon name="keyboard" size="28px" />
          {{ $t('a11y.shortcutsTitle') || 'Scorciatoie da Tastiera & Accessibilità' }}
        </div>
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section class="q-pt-sm">
        <div class="text-caption text-grey-7 q-mb-md">
          {{ $t('a11y.shortcutsSubtitle') || 'Usa la combinazione di tasti per navigare rapidamente nel registro senza utilizzare il mouse.' }}
        </div>

        <q-input
          v-model="filterQuery"
          dense
          outlined
          placeholder="Cerca comando o scorciatoia..."
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
        <q-btn flat label="Chiudi (Esc)" color="primary" v-close-popup no-caps />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useThemeStore } from 'src/stores/theme'

const themeStore = useThemeStore()

const isOpen = computed({
  get: () => themeStore.keyboardShortcutsHelpOpen,
  set: (val) => themeStore.toggleKeyboardShortcutsHelp(val)
})

const filterQuery = ref('')

const shortcutGroups = [
  {
    category: 'Navigazione Globale (Ruolo Docente & Segreteria)',
    items: [
      { description: 'Vai alla Dashboard', keys: ['Alt', '1'] },
      { description: 'Vai al Registro Voti', keys: ['Alt', 'V'] },
      { description: 'Vai al Registro Presenze', keys: ['Alt', 'P'] },
      { description: 'Vai all\'Agenda & Compiti', keys: ['Alt', 'A'] },
      { description: 'Attiva Ricerca Globale', keys: ['Alt', 'S'] }
    ]
  },
  {
    category: 'Funzionalità di Accessibilità (A11y)',
    items: [
      { description: 'Mostra / Nascondi Guida Scorciatoie', keys: ['?'] },
      { description: 'Attiva / Disattiva Righello di Lettura', keys: ['Alt', 'R'] },
      { description: 'Sposta Righello Su / Giù', keys: ['Alt', '↑ / ↓'] },
      { description: 'Attiva / Disattiva Lettura Vocale (TTS)', keys: ['Alt', 'T'] },
      { description: 'Salta al Contenuto Principale (Skip to Main)', keys: ['Tab (Inizio Pagina)'] },
      { description: 'Chiudi Modale / Annulla Operazione', keys: ['Esc'] }
    ]
  }
]

const filteredGroups = computed(() => {
  if (!filterQuery.value) return shortcutGroups
  const q = filterQuery.value.toLowerCase()
  return shortcutGroups
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
