<template>
  <router-view />

  <!-- Modal di conferma ricarica pagina (F5 / Ctrl+R) -->
  <q-dialog v-model="showReloadConfirmDialog" persistent>
    <q-card style="min-width: 380px; max-width: 480px;" class="rounded-xl">
      <q-card-section class="bg-amber-8 text-white row items-center">
        <q-icon name="security" size="28px" class="q-mr-sm" />
        <div class="text-h6 text-weight-bold">Conferma Ricarica Pagina</div>
      </q-card-section>

      <q-card-section class="q-pa-lg text-body1 text-slate-800">
        <p class="q-mb-sm text-weight-medium">
          Sei sicuro di voler ricaricare la pagina?
        </p>
        <p class="text-caption text-grey-8 q-mb-none">
          Per ragioni di sicurezza i dati di sessione sensibili non vengono memorizzati in modo permanente sul browser. 
          Ricaricando la pagina i dati non salvati andranno persi e potresti dover effettuare nuovamente l'accesso.
        </p>
      </q-card-section>

      <q-card-actions align="right" class="q-pa-md bg-grey-1">
        <q-btn flat label="Annulla (Rimani qui)" color="grey-8" v-close-popup no-caps />
        <q-btn 
          color="amber-9" 
          label="Conferma e Ricarica" 
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
import { useAuthStore } from 'src/stores/auth'
import { useWebSocketStore } from 'src/stores/websocket'
import { useThemeStore } from 'src/stores/theme'

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
  window.removeEventListener('beforeunload', handleBeforeUnload)
  showReloadConfirmDialog.value = false
  window.location.reload()
}

onMounted(async () => {
  themeStore.initTheme()
  window.addEventListener('beforeunload', handleBeforeUnload)
  window.addEventListener('keydown', handleKeyDown)

  if (!authStore.token && (localStorage.getItem('user') || sessionStorage.getItem('user'))) {
    await authStore.initAuth()
  }
  if (authStore.isAuthenticated) {
    wsStore.connect()
  }
})

onUnmounted(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
  window.removeEventListener('keydown', handleKeyDown)
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
@import url('https://fonts.googleapis.com/css2?family=Atkinson+Hyperlegible:ital,wght@0,400;0,700;1,400;1,700&family=Fredoka:wght@400;600;700&family=Lexend:wght@300;400;500;600;700&family=Roboto:wght@400;500;700&display=swap');

@font-face {
  font-family: 'OpenDyslexic';
  src: url('https://cdn.jsdelivr.net/npm/open-dyslexic@1.0.3/otf/OpenDyslexic-Regular.otf') format('opentype');
  font-weight: normal;
  font-style: normal;
}
@font-face {
  font-family: 'OpenDyslexic';
  src: url('https://cdn.jsdelivr.net/npm/open-dyslexic@1.0.3/otf/OpenDyslexic-Bold.otf') format('opentype');
  font-weight: bold;
  font-style: normal;
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
</style>
