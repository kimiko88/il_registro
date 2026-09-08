<template>
  <div>
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
          <q-icon :name="outboxStore.hasPending ? 'cloud_queue' : 'wifi_off'" size="20px" />
          <span class="text-weight-medium text-body2">
            {{ outboxStore.hasPending
              ? (t('offlineBanner.pendingQueue', { count: outboxStore.pendingCount }) || `Nessuna connessione. ${outboxStore.pendingCount} modifiche salvate in locale verranno sincronizzate al ripristino.`)
              : (t('offlineBanner.offlineText') || 'Nessuna connessione a Internet. Le modifiche potrebbero non essere salvate.')
            }}
          </span>
        </div>
        <div class="row items-center q-gutter-x-xs">
          <q-badge
            v-if="outboxStore.hasPending"
            color="warning"
            text-color="dark"
            class="text-weight-bold cursor-pointer"
            role="button"
            tabindex="0"
            :aria-label="t('offlineBanner.viewQueue') || 'Visualizza la coda delle operazioni'"
            @click="showQueueDialog = true"
            @keydown.enter="showQueueDialog = true"
          >
            <q-icon name="list_alt" size="14px" class="q-mr-xs" />
            {{ outboxStore.pendingCount }} {{ t('offlineBanner.inQueue') || 'IN CODA' }}
            <q-tooltip>{{ t('offlineBanner.viewQueue') || 'Clicca per gestire le operazioni in coda' }}</q-tooltip>
          </q-badge>
          <q-badge color="white" text-color="negative" class="text-weight-bold">
            {{ t('offlineBanner.offlineBadge') || 'OFFLINE' }}
          </q-badge>
        </div>
      </div>

      <div
        v-else-if="showBackOnline"
        class="offline-banner bg-positive text-white q-px-md q-py-xs row items-center justify-between shadow-2"
        role="status"
        aria-live="polite"
      >
        <div class="row items-center q-gutter-x-sm">
          <q-icon :name="syncedCount > 0 ? 'cloud_done' : 'wifi'" size="20px" />
          <span class="text-weight-medium text-body2">
            {{ syncedCount > 0
              ? (t('offlineBanner.syncedSuccess', { count: syncedCount }) || `Connessione ripristinata: ${syncedCount} operazioni sincronizzate con successo.`)
              : (t('offlineBanner.onlineText') || 'Connessione a Internet ripristinata.')
            }}
          </span>
        </div>
        <q-btn flat round dense icon="close" size="xs" color="white" :aria-label="t('common.close') || 'Chiudi avviso online'" @click="showBackOnline = false" />
      </div>
    </transition>

    <OutboxQueueDialog v-model="showQueueDialog" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNetworkStatus } from '@/composables/useNetworkStatus'
import { useOutboxStore } from '@/stores/outbox'
import OutboxQueueDialog from '@/components/Common/OutboxQueueDialog.vue'

const { t } = useI18n()
const { isOnline, wasOffline } = useNetworkStatus()
const outboxStore = useOutboxStore()

const showQueueDialog = ref(false)
const showBackOnline = ref(false)
const syncedCount = ref(0)
let timer = null

watch(isOnline, async (online) => {
  if (online && wasOffline.value) {
    showBackOnline.value = true
    if (outboxStore.hasPending) {
      const res = await outboxStore.syncQueue()
      syncedCount.value = res?.synced || 0
    }
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      showBackOnline.value = false
      syncedCount.value = 0
    }, 5000)
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
