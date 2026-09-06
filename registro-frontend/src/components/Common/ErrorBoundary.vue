<template>
  <div v-if="hasError" class="error-boundary-wrapper q-pa-md flex flex-center" role="alert" aria-live="assertive">
    <q-card
      class="rounded-xl shadow-lg border full-width max-w-xl text-center q-pa-lg"
      :class="$q.dark.isActive ? 'bg-grey-9 text-white border-grey-8' : 'bg-white text-slate-800 border-slate-200'"
      style="max-width: 540px;"
    >
      <q-card-section class="q-pb-none">
        <q-avatar size="64px" color="negative" text-color="white" icon="error_outline" class="q-mb-md shadow-2" />
        <h2 class="text-h6 text-weight-bold q-my-none">
          {{ title || (t('common.somethingWentWrong') || 'Si è verificato un errore imprevisto') }}
        </h2>
        <p class="text-body2 q-mt-sm q-mb-none" :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-7'">
          {{ subtitle || (t('common.errorBoundaryHelp') || 'Questa sezione dell\'applicazione ha riscontrato un problema imprevisto. Puoi riprovare o tornare alla pagina principale.') }}
        </p>
      </q-card-section>

      <q-card-section v-if="errorMessage" class="q-py-md">
        <q-expansion-item
          dense
          dense-toggle
          expand-separator
          icon="code"
          :label="t('common.technicalDetails') || 'Dettagli tecnici'"
          header-class="text-caption text-grey-6"
          class="rounded-borders text-left"
          :class="$q.dark.isActive ? 'bg-grey-10 text-grey-3' : 'bg-grey-1 text-grey-9'"
        >
          <pre class="q-pa-sm text-caption text-negative overflow-auto" style="max-height: 150px; font-family: monospace; font-size: 11px;">{{ errorMessage }}</pre>
        </q-expansion-item>
      </q-card-section>

      <q-card-actions align="center" class="q-gutter-sm q-pt-sm">
        <q-btn
          data-test="retry-btn"
          color="primary"
          icon="refresh"
          unelevated
          no-caps
          class="rounded-lg text-weight-bold"
          :label="t('common.retry') || 'Riprova'"
          @click="resetError"
        />
        <q-btn
          data-test="home-btn"
          outline
          color="grey-7"
          icon="home"
          no-caps
          class="rounded-lg"
          :label="t('common.backToHome') || 'Torna alla Dashboard'"
          @click="goHome"
        />
      </q-card-actions>
    </q-card>
  </div>
  <slot v-else :key="retryKey" />
</template>

<script setup>
import { ref, onErrorCaptured } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useErrorStore } from '@/stores/error'

defineProps({
  title: {
    type: String,
    default: null
  },
  subtitle: {
    type: String,
    default: null
  }
})

const emit = defineEmits(['error-captured'])

const router = useRouter()
const { t } = useI18n()
const $q = useQuasar()
const errorStore = useErrorStore()

const hasError = ref(false)
const errorMessage = ref('')
const retryKey = ref(0)

onErrorCaptured((err, instance, info) => {
  hasError.value = true
  errorMessage.value = err?.stack || err?.message || String(err)

  try {
    errorStore.reportError(err, t('common.somethingWentWrong') || 'Errore imprevisto nel componente')
  } catch (_e) {
    // Store not ready or error in error reporter
  }

  emit('error-captured', { error: err, info })

  // Prevent error from propagating further up
  return false
})

function resetError() {
  hasError.value = false
  errorMessage.value = ''
  retryKey.value++
}

function goHome() {
  resetError()
  if (router && typeof router.push === 'function') {
    router.push('/')
  } else if (typeof window !== 'undefined') {
    window.location.href = '/'
  }
}

defineExpose({
  resetError,
  hasError
})
</script>
