<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="groups" color="primary" class="q-mr-sm" />
          {{ t('verbaliPage.title') }}
        </div>
        <div class="text-caption text-grey">{{ t('verbaliPage.subtitle') }}</div>
      </div>
      <q-btn round flat icon="refresh" :loading="loading" @click="loadMeetings" />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="50px" />
    </div>

    <!-- Empty -->
    <q-card v-else-if="meetings.length === 0" class="text-center q-pa-xl text-grey shadow-1">
      <q-icon name="event_busy" size="60px" class="q-mb-md" />
      <div class="text-h6">{{ t('verbaliPage.noVerbali') }}</div>
    </q-card>

    <!-- Meetings list -->
    <div v-else class="row q-col-gutter-md">
      <div v-for="m in meetings" :key="m.id" class="col-12 col-md-6">
        <q-card bordered flat class="full-height flex column">
          <q-card-section>
            <div class="row items-center justify-between q-mb-xs">
              <q-badge v-if="m.is_mandatory" color="negative" label="Obbligatoria" />
              <q-badge v-else color="info" label="Incontro" />
              <span class="text-caption text-grey">{{ formatDate(m.meeting_date) }}</span>
            </div>
            <div class="text-h6 text-weight-bold text-primary">{{ m.title }}</div>
            <div v-if="m.location" class="text-caption text-grey-8 q-mb-xs">
              <q-icon name="place" size="xs" /> {{ m.location }}
            </div>
            <div class="text-body2 text-grey-9 q-mt-sm">{{ m.description }}</div>
          </q-card-section>

          <q-space />

          <q-card-section class="bg-grey-2 q-py-sm row items-center justify-between">
            <div class="text-caption text-grey-8">
              <span v-if="m.max_participants">
                {{ m.registrations_count }}/{{ m.max_participants }}
              </span>
              <span v-else>{{ m.registrations_count }}</span>
              <div v-if="m.registration_deadline" class="text-caption text-negative">
                {{ formatDate(m.registration_deadline) }}
              </div>
            </div>

            <div>
              <q-btn
                v-if="m.is_registered"
                outline color="negative" :label="t('common.cancel')" size="sm"
                :loading="processingId === m.id"
                @click="toggleRegistration(m)"
              />
              <q-btn
                v-else
                color="primary" :label="t('common.save')" size="sm"
                :loading="processingId === m.id"
                :disable="isFull(m) || isExpired(m)"
                @click="toggleRegistration(m)"
              />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const loading = ref(false)
const processingId = ref(null)
const meetings = ref([])

onMounted(() => {
  loadMeetings()
})

const formatDate = (d) => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY HH:mm') : ''
const isExpired = (m) => m.registration_deadline && new Date() > new Date(m.registration_deadline)
const isFull = (m) => m.max_participants && m.registrations_count >= m.max_participants

async function loadMeetings() {
  loading.value = true
  try {
    const res = await api.get('/general-meetings')
    meetings.value = res.data || []
  } catch (e) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    loading.value = false
  }
}

async function toggleRegistration(m) {
  processingId.value = m.id
  try {
    if (m.is_registered) {
      await api.post(`/general-meetings/${m.id}/unregister`)
      m.is_registered = false
      m.registrations_count--
      $q.notify({ type: 'info', message: t('common.success') })
    } else {
      await api.post(`/general-meetings/${m.id}/register`)
      m.is_registered = true
      m.registrations_count++
      $q.notify({ type: 'positive', message: t('common.success') })
    }
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || t('common.error') })
  } finally {
    processingId.value = null
  }
}
</script>
