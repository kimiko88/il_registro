<template>
  <transition
    appear
    enter-active-class="animated fadeInDown"
    leave-active-class="animated fadeOutUp"
  >
    <div
      v-if="needRefresh && !userDismissed"
      class="pwa-update-banner bg-primary text-white q-px-md q-py-xs row items-center justify-between shadow-3"
      role="alert"
      aria-live="polite"
    >
      <div class="row items-center q-gutter-x-sm">
        <q-icon name="system_update" size="20px" aria-hidden="true" />
        <span class="text-weight-medium text-body2">
          {{ t('pwa.updateAvailable') || 'È disponibile una nuova versione del Registro Elettronico.' }}
        </span>
      </div>
      <div class="row items-center q-gutter-x-xs">
        <q-btn
          flat
          dense
          no-caps
          color="white"
          class="text-weight-bold q-px-sm"
          icon="refresh"
          :label="t('pwa.reloadNow') || 'Aggiorna Ora'"
          @click="updateAndReload"
        />
        <q-btn
          flat
          round
          dense
          icon="close"
          size="xs"
          color="white"
          @click="dismiss"
          :aria-label="t('common.close') || 'Chiudi'"
        />
      </div>
    </div>
  </transition>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { usePwaUpdate } from '@/composables/usePwaUpdate'

const { t } = useI18n()
const { needRefresh, userDismissed, updateAndReload, dismiss } = usePwaUpdate()
</script>

<style scoped>
.pwa-update-banner {
  position: sticky;
  top: 0;
  z-index: 6500;
  width: 100%;
}
</style>
