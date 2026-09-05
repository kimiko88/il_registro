<template>
  <transition
    appear
    enter-active-class="animated fadeInDown"
    leave-active-class="animated fadeOutUp"
  >
    <div
      v-if="!isOnline"
      class="offline-banner bg-negative text-white q-px-md q-py-xs row items-center justify-between shadow-2"
      role="alert"
      aria-live="assertive"
    >
      <div class="row items-center q-gutter-x-sm">
        <q-icon name="wifi_off" size="20px" />
        <span class="text-weight-medium text-body2">
          {{ t('offlineBanner.offlineText') || 'Nessuna connessione a Internet. Le modifiche potrebbero non essere salvate.' }}
        </span>
      </div>
      <q-badge color="white" text-color="negative" class="text-weight-bold">
        {{ t('offlineBanner.offlineBadge') || 'OFFLINE' }}
      </q-badge>
    </div>

    <div
      v-else-if="showBackOnline"
      class="offline-banner bg-positive text-white q-px-md q-py-xs row items-center justify-between shadow-2"
      role="status"
      aria-live="polite"
    >
      <div class="row items-center q-gutter-x-sm">
        <q-icon name="wifi" size="20px" />
        <span class="text-weight-medium text-body2">
          {{ t('offlineBanner.onlineText') || 'Connessione a Internet ripristinata.' }}
        </span>
      </div>
      <q-btn flat round dense icon="close" size="xs" color="white" :aria-label="t('common.close') || 'Chiudi avviso online'" @click="showBackOnline = false" />
    </div>
  </transition>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNetworkStatus } from '@/composables/useNetworkStatus'

const { t } = useI18n()
const { isOnline, wasOffline } = useNetworkStatus()

const showBackOnline = ref(false)
let timer = null

watch(isOnline, (online) => {
  if (online && wasOffline.value) {
    showBackOnline.value = true
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      showBackOnline.value = false
    }, 4000)
  }
})
</script>

<style scoped>
.offline-banner {
  position: sticky;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 10000;
  min-height: 38px;
}
</style>
