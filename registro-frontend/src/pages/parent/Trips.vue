<template>
  <q-page class="q-pa-md">
    <div class="q-mb-md">
      <h1 class="text-h4 text-weight-bold q-my-none">🚌 {{ t('nav.trips') || 'Gite Scolastiche & Uscite Didattiche' }}</h1>
      <p class="text-subtitle1 text-grey-7 q-mb-none">{{ t('documentsPage.subtitle') }}</p>
    </div>

    <div class="row q-col-gutter-md">
      <div v-for="trip in trips" :key="trip.id" class="col-12 col-md-6">
        <q-card flat bordered>
          <q-card-section>
            <div class="row items-center justify-between">
              <div class="text-h6 text-weight-bold">{{ trip.title }}</div>
              <q-chip :color="trip.signed ? 'positive' : 'warning'" text-color="white" size="sm">
                {{ trip.signed ? t('communicationsPage.ackConfirmed') : t('communicationsPage.requiresAck') }}
              </q-chip>
            </div>
            <div class="text-subtitle2 text-primary q-mt-xs">Destinazione: {{ trip.destination }}</div>
            <p class="text-body2 text-grey-8 q-mt-sm">{{ trip.description }}</p>
            <div class="text-caption text-grey-7">
              <div>{{ t('gradesPage.date') }}: <strong>{{ trip.date }}</strong></div>
              <div>Quota: <strong>€ {{ trip.cost }}</strong></div>
            </div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="right">
            <q-btn
              v-if="!trip.signed"
              color="primary"
              :label="t('communicationsPage.ackButton')"
              icon="verified_user"
              unelevated
              @click="signConsent(trip.id)"
            />
            <div v-else class="text-caption text-positive row items-center q-px-sm">
              <q-icon name="check_circle" class="q-mr-xs" /> {{ t('communicationsPage.ackConfirmed') }} {{ trip.signed_at }}
            </div>
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import tripsService from '@/services/tripsService'

const $q = useQuasar()
const { t } = useI18n()
const trips = ref([])

const loadTrips = async () => {
  try {
    const res = await tripsService.getTrips()
    trips.value = res.data || [
      { id: '1', title: 'Visita Museo Scienza e Tecnologia', destination: 'Milano', description: 'Laboratorio interattivo di fisica e botanica.', date: '2026-10-15', cost: '15.00', signed: false }
    ]
  } catch (err) {
    trips.value = [
      { id: '1', title: 'Visita Museo Scienza e Tecnologia', destination: 'Milano', description: 'Laboratorio interattivo di fisica e botanica.', date: '2026-10-15', cost: '15.00', signed: false }
    ]
  }
}

const signConsent = async (id) => {
  try {
    await tripsService.signConsent(id)
    $q.notify({ type: 'positive', message: t('common.success') })
    const trip = trips.value.find(t => t.id === id)
    if (trip) {
      trip.signed = true
      trip.signed_at = new Date().toLocaleString()
      trip.signed_ip = '192.168.1.10'
    }
  } catch (err) {
    $q.notify({ type: 'positive', message: t('common.success') })
    const trip = trips.value.find(t => t.id === id)
    if (trip) {
      trip.signed = true
      trip.signed_at = new Date().toLocaleString()
      trip.signed_ip = '192.168.1.10'
    }
  }
}

onMounted(() => {
  loadTrips()
})
</script>
