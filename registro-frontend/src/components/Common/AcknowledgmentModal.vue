<template>
  <q-dialog v-model="show" persistent backdrop-filter="blur(4px)">
    <q-card style="width: 550px; max-width: 90vw;" class="rounded-xl border border-red-2 shadow-2xl">
      <q-card-section class="bg-red-7 text-white row items-center q-pb-md">
        <q-avatar icon="priority_high" color="white" text-color="red-7" class="q-mr-sm" size="40px" />
        <div>
          <div class="text-h6 text-weight-bold">{{ t('communicationsPage.title') }}</div>
          <div class="text-caption opacity-90">{{ t('communicationsPage.requiresAck') }}</div>
        </div>
      </q-card-section>

      <q-card-section class="q-pa-lg text-slate-800" v-if="communication">
        <div class="text-subtitle1 text-weight-bold text-slate-900 q-mb-xs">
          {{ communication.subject }}
        </div>
        <div class="text-caption text-slate-500 q-mb-md">
          {{ t('communicationsPage.publishDate') }}: {{ formatDate(communication.created_at) }}
        </div>

        <div class="text-body2 bg-slate-50 border q-pa-md rounded-lg text-slate-700 q-mb-md leading-relaxed" style="max-height: 220px; overflow-y: auto;">
          {{ communication.body }}
        </div>

        <q-banner dense rounded class="bg-amber-50 text-amber-9 border border-amber-2 q-mb-none">
          <template v-slot:avatar>
            <q-icon name="info" color="amber-9" />
          </template>
          {{ t('communicationsPage.requiresAck') }}
        </q-banner>
      </q-card-section>

      <q-card-actions align="right" class="q-pa-md bg-slate-50 border-t">
        <q-btn
          color="red-7"
          icon="check_circle"
          :label="t('communicationsPage.ackButton')"
          unelevated
          class="rounded-lg text-weight-bold full-width"
          size="lg"
          :loading="loading"
          @click="confirmAck"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import api from '@/services/api'

const $q = useQuasar()
const { t, locale } = useI18n()
const show = ref(false)
const communication = ref(null)
const loading = ref(false)

onMounted(async () => {
  try {
    const unreadRes = await api.get('/communications/bacheca')
    const urgent = (unreadRes.data || []).find(c => c.requires_acknowledgment && !c.acknowledged)
    if (urgent) {
      communication.value = urgent
      show.value = true
    }
  } catch (err) {
    // Ignore error silently
  }
})

async function confirmAck() {
  if (!communication.value) return
  loading.value = true
  try {
    await api.post(`/communications/${communication.value.id}/ack`)
    $q.notify({ type: 'positive', message: t('communicationsPage.ackConfirmed') })
    show.value = false
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    loading.value = false
  }
}

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(locale.value || 'it-IT', { day: '2-digit', month: 'long', year: 'numeric' })
}
</script>
