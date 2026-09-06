<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card style="min-width: 380px; max-width: 600px; width: 95vw;" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-white'">
      <q-card-section class="row items-center justify-between q-pb-none">
        <div class="row items-center q-gutter-x-sm">
          <q-avatar icon="cloud_queue" color="primary" text-color="white" size="36px" />
          <div>
            <div class="text-h6 text-weight-bold">
              {{ t('offlineQueue.title') || 'Coda Operazioni Offline' }}
            </div>
            <div class="text-caption text-grey">
              {{ outboxStore.pendingCount }} {{ t('offlineQueue.pendingItems') || 'operazioni in attesa di sincronizzazione' }}
            </div>
          </div>
        </div>
        <q-btn flat round dense icon="close" :aria-label="t('common.close') || 'Chiudi'" @click="isOpen = false" />
      </q-card-section>

      <q-separator class="q-my-sm" />

      <q-card-section class="q-pa-none" style="max-height: 55vh; overflow-y: auto;">
        <div v-if="!outboxStore.hasPending" class="text-center q-pa-lg">
          <q-icon name="cloud_done" size="56px" color="positive" />
          <div class="text-subtitle1 text-weight-medium q-mt-sm">
            {{ t('offlineQueue.allSynced') || 'Nessuna operazione in sospeso' }}
          </div>
          <div class="text-caption text-grey">
            {{ t('offlineQueue.allSyncedSub') || 'Tutte le modifiche effettuate sono state inviate con successo al server.' }}
          </div>
        </div>

        <q-list v-else separator>
          <q-item v-for="item in outboxStore.queue" :key="item.id" class="q-py-sm">
            <q-item-section avatar>
              <q-avatar :color="getMethodColor(item.method)" text-color="white" size="30px" font-size="12px" class="text-weight-bolder">
                {{ (item.method || 'POST').toUpperCase() }}
              </q-avatar>
            </q-item-section>

            <q-item-section>
              <q-item-label class="text-weight-medium ellipsis">
                {{ item.title || item.url }}
              </q-item-label>
              <q-item-label caption class="row items-center q-gutter-x-xs">
                <span>{{ formatTimestamp(item.timestamp) }}</span>
                <span v-if="item.attempts > 0" class="text-warning">
                  • {{ item.attempts }} {{ t('offlineQueue.retries') || 'tentativi falliti' }}
                </span>
                <span v-if="item.nextRetryAt && item.nextRetryAt > Date.now()" class="text-info">
                  • {{ t('offlineQueue.backoff') || 'in attesa backoff' }}
                </span>
              </q-item-label>
            </q-item-section>

            <q-item-section side>
              <q-btn
                flat
                round
                dense
                color="negative"
                icon="delete_outline"
                size="sm"
                :aria-label="t('common.delete') || 'Elimina'"
                @click="removeItem(item.id)"
              >
                <q-tooltip>{{ t('offlineQueue.removeItemTooltip') || 'Elimina dalla coda' }}</q-tooltip>
              </q-btn>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>

      <q-separator />

      <q-card-actions align="between" class="q-pa-md">
        <q-btn
          v-if="outboxStore.hasPending"
          flat
          dense
          color="negative"
          icon="clear_all"
          :label="t('offlineQueue.clearAll') || 'Svuota Coda'"
          :disable="outboxStore.isSyncing"
          @click="confirmClearAll"
        />
        <div v-else />

        <div class="row q-gutter-sm">
          <q-btn
            v-if="outboxStore.hasPending"
            unelevated
            color="primary"
            icon="sync"
            :label="t('offlineQueue.syncNow') || 'Sincronizza Ora'"
            :loading="outboxStore.isSyncing"
            @click="handleSyncNow"
          />
          <q-btn
            flat
            :label="t('common.close') || 'Chiudi'"
            @click="isOpen = false"
          />
        </div>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Dialog } from 'quasar'
import { useOutboxStore } from '@/stores/outbox'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const outboxStore = useOutboxStore()

const isOpen = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

function getMethodColor(method) {
  const m = (method || '').toLowerCase()
  if (m === 'post') return 'positive'
  if (m === 'put' || m === 'patch') return 'warning'
  if (m === 'delete') return 'negative'
  return 'primary'
}

function formatTimestamp(ts) {
  if (!ts) return ''
  const date = new Date(ts)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function removeItem(id) {
  outboxStore.dequeue(id)
}

function confirmClearAll() {
  Dialog.create({
    title: t('offlineQueue.confirmClearTitle') || 'Svuotare la coda?',
    message: t('offlineQueue.confirmClearMessage') || 'Le modifiche salvate in locale non ancora sincronizzate andranno perse.',
    ok: {
      label: t('common.confirm') || 'Svuota',
      color: 'negative',
      flat: false
    },
    cancel: {
      label: t('common.cancel') || 'Annulla',
      flat: true
    }
  }).onOk(() => {
    outboxStore.clear()
  })
}

async function handleSyncNow() {
  await outboxStore.syncQueue()
}
</script>
