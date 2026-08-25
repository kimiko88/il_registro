<template>
  <q-page class="q-pa-md">
    <div class="q-mb-md">
      <h1 class="text-h4 text-weight-bold q-my-none">🚌 {{ t('nav.trips') }}</h1>
      <p class="text-subtitle1 text-grey-7 q-mb-none">{{ t('tripsPage.subtitle') }}</p>
    </div>

    <!-- Loading spinner -->
    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Empty state (no mock fallback) -->
    <div v-else-if="trips.length === 0" class="text-center q-my-xl q-pa-lg bg-grey-1 rounded-borders">
      <q-icon name="directions_bus" size="4rem" color="grey-5" />
      <div class="text-h6 text-grey-7 q-mt-md">{{ t('common.noData') }}</div>
      <p class="text-caption text-grey-6 q-mb-none">{{ t('tripsPage.noTrips') }}</p>
    </div>

    <!-- Real Trips list -->
    <div v-else class="row q-col-gutter-md">
      <div v-for="trip in trips" :key="trip.id" class="col-12 col-md-6">
        <q-card flat bordered class="full-height column justify-between shadow-1">
          <q-card-section>
            <div class="row items-center justify-between no-wrap">
              <div class="text-h6 text-weight-bold ellipsis">{{ trip.title }}</div>
              <q-chip
                :color="isConsentGiven(trip) ? 'positive' : 'warning'"
                text-color="white"
                size="sm"
                class="q-ml-sm"
              >
                {{ isConsentGiven(trip) ? t('communicationsPage.ackConfirmed') : t('communicationsPage.requiresAck') }}
              </q-chip>
            </div>

            <div class="text-subtitle2 text-primary q-mt-xs">
              <q-icon name="place" size="xs" class="q-mr-xs" />
              {{ t('tripsPage.destination') }}: <strong>{{ trip.destination }}</strong>
            </div>

            <p v-if="trip.description" class="text-body2 text-grey-8 q-mt-sm">{{ trip.description }}</p>

            <div class="text-caption text-grey-7 q-mt-md">
              <div v-if="trip.departure_date">
                <q-icon name="event" size="xs" class="q-mr-xs" />
                {{ t('tripsPage.departureDate') }}: <strong>{{ formatDate(trip.departure_date) }}</strong>
                <span v-if="trip.return_date && trip.return_date !== trip.departure_date">
                  - {{ formatDate(trip.return_date) }}
                </span>
              </div>
              <div v-if="trip.accompanying_teachers" class="q-mt-xs">
                <q-icon name="school" size="xs" class="q-mr-xs" />
                {{ t('tripsPage.teachers') }}: <strong>{{ trip.accompanying_teachers }}</strong>
              </div>
            </div>
          </q-card-section>

          <div>
            <q-separator />
            <q-card-actions align="right" class="q-pa-md">
              <q-btn
                v-if="!isConsentGiven(trip)"
                color="primary"
                :label="t('communicationsPage.ackButton')"
                icon="verified_user"
                unelevated
                :loading="submittingId === trip.id"
                @click="signConsent(trip.id)"
              />
              <div v-else class="text-caption text-positive row items-center q-px-sm">
                <q-icon name="check_circle" class="q-mr-xs" />
                {{ t('communicationsPage.ackConfirmed') }}
              </div>
            </q-card-actions>
          </div>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import tripsService from '@/services/tripsService'
import { useChildrenStore } from '@/stores/children'

const $q = useQuasar()
const { t } = useI18n()
const childrenStore = useChildrenStore()

const trips = ref([])
const loading = ref(false)
const submittingId = ref(null)

const isConsentGiven = (trip) => {
  return trip.consent_status === 'granted' || trip.consent_status === 'Consented' || trip.signed === true
}

const formatDate = (val) => {
  if (!val) return ''
  try {
    const d = new Date(val)
    return isNaN(d.getTime()) ? val : d.toLocaleDateString('it-IT')
  } catch {
    return val
  }
}

const loadTrips = async () => {
  loading.value = true
  try {
    const studentId = childrenStore.selectedChildId || childrenStore.children[0]?.id
    const params = studentId ? { student_id: studentId } : {}
    const res = await tripsService.getTrips(params)
    trips.value = Array.isArray(res.data) ? res.data : (res.data?.trips || [])
  } catch (err) {
    trips.value = []
    console.error('Error fetching trips from database:', err)
  } finally {
    loading.value = false
  }
}

const signConsent = async (id) => {
  submittingId.value = id
  try {
    const studentId = childrenStore.selectedChildId || childrenStore.children[0]?.id || childrenStore.children[0]?.userId
    await tripsService.signConsent(id, {
      student_id: studentId,
      status: 'granted'
    })
    $q.notify({ type: 'positive', message: t('common.success') })
    const trip = trips.value.find(t => t.id === id)
    if (trip) {
      trip.consent_status = 'granted'
      trip.signed = true
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('common.error')
    })
  } finally {
    submittingId.value = null
  }
}

watch(
  () => childrenStore.selectedChildId,
  () => {
    loadTrips()
  }
)

onMounted(async () => {
  if (childrenStore.children.length === 0) {
    await childrenStore.fetchChildren()
  }
  await loadTrips()
})
</script>
