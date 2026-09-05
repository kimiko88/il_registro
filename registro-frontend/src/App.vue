<template>
  <router-view />

  <!-- Modal di conferma ricarica pagina (F5 / Ctrl+R) -->
  <q-dialog v-model="showReloadConfirmDialog" persistent>
    <q-card style="min-width: 380px; max-width: 480px;" class="rounded-xl">
      <q-card-section class="bg-amber-8 text-white row items-center">
        <q-icon name="security" size="28px" class="q-mr-sm" />
        <div class="text-h6 text-weight-bold">{{ t('common.confirm') || 'Conferma Ricarica Pagina' }}</div>
      </q-card-section>

      <q-card-section class="q-pa-lg text-body1 text-slate-800">
        <p class="q-mb-sm text-weight-medium">
          {{ t('settingsPage.security') || 'Sei sicuro di voler ricaricare la pagina?' }}
        </p>
        <p class="text-caption text-grey-8 q-mb-none">
          {{ t('settingsPage.subtitle') || 'Per ragioni di sicurezza i dati non salvati andranno persi.' }}
        </p>
      </q-card-section>

      <q-card-actions align="right" class="q-pa-md bg-grey-1">
        <q-btn flat :label="t('common.cancel') || 'Annulla'" color="grey-8" v-close-popup no-caps />
        <q-btn 
          color="amber-9" 
          :label="t('common.confirm') || 'Conferma e Ricarica'" 
          icon="refresh" 
          no-caps 
          unelevated 
          @click="executeReload" 
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useWebSocketStore } from '@/stores/websocket'
import { useThemeStore } from '@/stores/theme'

const { t } = useI18n()
const authStore = useAuthStore()
const wsStore = useWebSocketStore()
const themeStore = useThemeStore()

const showReloadConfirmDialog = ref(false)

function handleBeforeUnload(e) {
  if (authStore.isAuthenticated) {
    e.preventDefault()
    e.returnValue = ''
    return ''
  }
}

function handleKeyDown(e) {
  if (!authStore.isAuthenticated) return

  if (e.key === 'F5' || ((e.ctrlKey || e.metaKey) && (e.key === 'r' || e.key === 'R'))) {
    e.preventDefault()
    showReloadConfirmDialog.value = true
  }
}

function executeReload() {
  if (typeof window !== 'undefined') {
    window.removeEventListener('beforeunload', handleBeforeUnload)
    window.removeEventListener('keydown', handleKeyDown)
    showReloadConfirmDialog.value = false
    window.location.reload()
  }
}

onMounted(async () => {
  themeStore.initTheme()
  if (typeof window !== 'undefined') {
    window.addEventListener('beforeunload', handleBeforeUnload)
    window.addEventListener('keydown', handleKeyDown)
  }

  if (!authStore.token && (localStorage.getItem('user') || sessionStorage.getItem('user'))) {
    await authStore.initAuth()
  }
  if (authStore.isAuthenticated) {
    wsStore.connect()
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('beforeunload', handleBeforeUnload)
    window.removeEventListener('keydown', handleKeyDown)
  }
  try {
    wsStore.disconnect()
  } catch {
    // ws store cleanup guard
  }
})

watch(() => authStore.isAuthenticated, (val) => {
  if (val) {
    wsStore.connect()
  } else {
    wsStore.disconnect()
  }
})
</script>

<style>
@font-face {
  font-family: 'OpenDyslexic';
  src: url('https://cdn.jsdelivr.net/npm/open-dyslexic@1.0.3/otf/OpenDyslexic-Regular.otf') format('opentype');
  font-weight: normal;
  font-style: normal;
  font-display: swap;
}
@font-face {
  font-family: 'OpenDyslexic';
  src: url('https://cdn.jsdelivr.net/npm/open-dyslexic@1.0.3/otf/OpenDyslexic-Bold.otf') format('opentype');
  font-weight: bold;
  font-style: normal;
  font-display: swap;
}

/* Accessibility DSA Font — applies to all elements EXCEPT icon fonts */
body.dsa-font-active,
body.dsa-font-active *:not(.material-icons):not(.material-symbols-outlined):not(.material-symbols-rounded):not(.material-symbols-sharp):not([class*="q-icon"]):not(i.q-icon):not(.notranslate) {
  font-family: 'OpenDyslexic', 'Atkinson Hyperlegible', 'Trebuchet MS', sans-serif !important;
  letter-spacing: 0.05em !important;
  word-spacing: 0.12em !important;
  line-height: 1.65 !important;
}

/* Ensure Material Icons always keep their font regardless of DSA mode */
body.dsa-font-active .material-icons,
body.dsa-font-active .material-symbols-outlined,
body.dsa-font-active .material-symbols-rounded,
body.dsa-font-active .material-symbols-sharp,
body.dsa-font-active i.q-icon,
body.dsa-font-active .q-icon,
body.dsa-font-active .notranslate {
  font-family: 'Material Icons', 'Material Symbols Outlined', 'Material Symbols Rounded', 'Material Symbols Sharp' !important;
  letter-spacing: normal !important;
  word-spacing: normal !important;
  line-height: 1 !important;
  font-feature-settings: 'liga' !important;
  -webkit-font-feature-settings: 'liga' !important;
  text-rendering: optimizeLegibility !important;
}

/* Additional Font Families */
body.font-family-lexend,
body.font-family-lexend * {
  font-family: 'Lexend', sans-serif !important;
}

body.font-family-fredoka,
body.font-family-fredoka * {
  font-family: 'Fredoka', cursive, sans-serif !important;
}

body.font-family-roboto,
body.font-family-roboto * {
  font-family: 'Roboto', sans-serif !important;
}

/* Font Size Scaling */
body.font-size-large {
  font-size: 110% !important;
}

body.font-size-xlarge {
  font-size: 122% !important;
}

/* High Contrast Mode */
body.high-contrast-active {
  filter: contrast(135%);
}
body.high-contrast-active .q-card,
body.high-contrast-active .q-btn,
body.high-contrast-active .q-table {
  border: 2px solid #0f172a !important;
}

/* Visible Focus Indicator (WCAG 2.4.7) */
:focus-visible,
body.focus-visible-high :focus-visible {
  outline: 3px solid #f59e0b !important;
  outline-offset: 3px !important;
  box-shadow: 0 0 0 5px rgba(245, 158, 11, 0.3) !important;
}

/* Text Spacing Overrides (WCAG 1.4.12) */
body.line-height-relaxed {
  line-height: 1.8 !important;
}
body.line-height-loose {
  line-height: 2.1 !important;
}
body.letter-spacing-wide {
  letter-spacing: 0.08em !important;
}
body.letter-spacing-wider {
  letter-spacing: 0.16em !important;
}
body.word-spacing-wide {
  word-spacing: 0.18em !important;
}
body.word-spacing-wider {
  word-spacing: 0.32em !important;
}

/* Colorblind Modes */
body.colorblind-protanopia {
  filter: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg"><filter id="protanopia"><feColorMatrix type="matrix" values="0.567, 0.433, 0, 0, 0 0.558, 0.442, 0, 0, 0 0, 0.242, 0.758, 0, 0 0, 0, 0, 1, 0"/></filter></svg>#protanopia');
}
body.colorblind-deuteranopia {
  filter: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg"><filter id="deuteranopia"><feColorMatrix type="matrix" values="0.625, 0.375, 0, 0, 0 0.7, 0.3, 0, 0, 0 0, 0.3, 0.7, 0, 0 0, 0, 0, 1, 0"/></filter></svg>#deuteranopia');
}
body.colorblind-tritanopia {
  filter: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg"><filter id="tritanopia"><feColorMatrix type="matrix" values="0.95, 0.05, 0, 0, 0 0, 0.433, 0.567, 0, 0 0, 0.475, 0.525, 0, 0 0, 0, 0, 1, 0"/></filter></svg>#tritanopia');
}
body.colorblind-monochrome {
  filter: grayscale(100%);
}

/* Advanced OLED High Contrast / Amber / Green */
body.contrast-oled-amber {
  background-color: #000000 !important;
  color: #fbbf24 !important;
}
body.contrast-oled-amber .q-card,
body.contrast-oled-amber .q-table,
body.contrast-oled-amber .q-header,
body.contrast-oled-amber .q-drawer {
  background-color: #000000 !important;
  color: #fbbf24 !important;
  border-color: #d97706 !important;
}

body.contrast-oled-green {
  background-color: #000000 !important;
  color: #4ade80 !important;
}
body.contrast-oled-green .q-card,
body.contrast-oled-green .q-table,
body.contrast-oled-green .q-header,
body.contrast-oled-green .q-drawer {
  background-color: #000000 !important;
  color: #4ade80 !important;
  border-color: #16a34a !important;
}

body.contrast-inverted {
  filter: invert(100%) hue-rotate(180deg);
}

/* Visual Grade Indicators */
.grade-badge-fail::before {
  content: '▼ ';
  font-size: 0.8em;
  margin-right: 2px;
}
.grade-badge-pass::before {
  content: '✓ ';
  font-size: 0.8em;
  margin-right: 2px;
}

/* Focus Mode (ADHD & DSA Clean Reading Layout) */
body.focus-mode-active .q-drawer {
  display: none !important;
}
body.focus-mode-active .q-page-container {
  padding-left: 0 !important;
  max-width: 920px !important;
  margin: 0 auto !important;
  line-height: 1.85 !important;
}
body.focus-mode-active .q-page {
  padding: 24px !important;
  font-size: 1.1rem !important;
}
</style>

