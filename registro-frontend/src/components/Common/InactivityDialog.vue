<template>
  <q-dialog v-model="showWarningDialog" persistent role="alertdialog" aria-modal="true">
    <q-card style="min-width: 380px; max-width: 460px;" class="rounded-xl shadow-24 overflow-hidden">
      <q-card-section class="bg-amber-8 text-white row items-center q-pa-md">
        <q-icon name="timer" size="28px" class="q-mr-sm" />
        <div class="text-h6 text-weight-bold">
          {{ t('inactivity.warningTitle') || 'Sessione in scadenza' }}
        </div>
      </q-card-section>

      <q-card-section class="q-pa-lg text-center">
        <p class="text-body1 text-slate-800 q-mb-md">
          {{ t('inactivity.warningMessage') || 'Non abbiamo rilevato alcuna attività da diversi minuti. Per ragioni di sicurezza, la sessione verrà chiusa a breve.' }}
        </p>

        <div class="column items-center q-my-md">
          <q-circular-progress
            show-value
            :value="progressValue"
            size="72px"
            :thickness="0.22"
            color="amber-8"
            track-color="amber-2"
            class="text-amber-9 text-weight-bolder text-h5"
          >
            {{ secondsRemaining }}s
          </q-circular-progress>
          <div class="text-caption text-grey-7 q-mt-xs">
            {{ t('inactivity.timeRemaining') || 'Secondi rimanenti' }}
          </div>
        </div>
      </q-card-section>

      <q-card-actions align="between" class="q-pa-md bg-slate-50">
        <q-btn
          flat
          color="grey-8"
          :label="t('inactivity.logoutNow') || 'Disconnetti Ora'"
          no-caps
          @click="logoutNow"
        />
        <q-btn
          unelevated
          color="primary"
          icon="lock_open"
          :label="t('inactivity.stayConnected') || 'Rimani Connesso'"
          class="rounded-lg text-weight-bold"
          no-caps
          @click="stayLoggedIn"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useInactivityTimer } from '@/composables/useInactivityTimer'

const { t } = useI18n()
const {
  showWarningDialog,
  secondsRemaining,
  stayLoggedIn,
  logoutNow
} = useInactivityTimer()

const progressValue = computed(() => {
  return Math.min(100, Math.max(0, (secondsRemaining.value / 60) * 100))
})
</script>
