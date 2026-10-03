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
                :label="t('tripsPage.authorizeWithPin') || 'Autorizza con PIN'"
                icon="verified_user"
                unelevated
                :loading="submittingId === trip.id"
                @click="openConsentDialog(trip)"
              />
              <div v-else class="column items-end q-gutter-xs">
                <div class="text-caption text-positive row items-center q-px-sm text-weight-bold">
                  <q-icon name="check_circle" class="q-mr-xs" size="18px" />
                  {{ t('tripsPage.authorizedDigitalPin') || 'Autorizzato con PIN Digitale' }}
                </div>
                <div class="row q-gutter-xs">
                  <q-badge color="teal-8" v-if="trip.dietary_notes">
                    <q-icon name="restaurant" class="q-mr-xs" size="10px" /> Dieta Speciale
                  </q-badge>
                  <q-badge color="purple-8" v-if="trip.medical_notes">
                    <q-icon name="medical_services" class="q-mr-xs" size="10px" /> Nota Medica
                  </q-badge>
                  <q-badge :color="trip.payment_status === 'paid' ? 'positive' : 'warning'" text-color="white">
                    {{ trip.payment_status === 'paid' ? 'Quota Saldata' : 'Quota in Attesa' }}
                  </q-badge>
                </div>
              </div>
            </q-card-actions>
          </div>
        </q-card>
      </div>
    </div>

    <!-- Digital Authorization PIN Modal -->
    <q-dialog v-model="consentDialog">
      <q-card style="width: min(520px, 95vw)" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="row items-center">
            <q-icon name="directions_bus" class="q-mr-sm" size="24px" />
            <div class="text-subtitle1 text-weight-bold">
              {{ t('tripsPage.authorizationTitle') || 'Autorizzazione Uscita Didattica' }}
            </div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <div class="bg-blue-50 text-blue-9 q-pa-sm rounded-borders text-caption">
            <strong>{{ selectedTripForConsent?.title }}</strong> — Destinazione: {{ selectedTripForConsent?.destination }}
          </div>

          <!-- PIN Dispositivo -->
          <q-input
            v-model="consentForm.pin"
            type="password"
            maxlength="6"
            label="PIN Dispositivo / Firma Elettronica Genitore *"
            placeholder="Es. 1234"
            outlined dense
            autofocus
            :rules="[val => !!val || 'Il PIN è obbligatorio per apporre la firma digitale di consenso']"
          >
            <template v-slot:prepend>
              <q-icon name="pin" color="primary" />
            </template>
          </q-input>

          <!-- Telefono Emergenza -->
          <q-input
            v-model="consentForm.emergency_phone"
            type="tel"
            label="Recapito Telefonico di Emergenza Reperibile *"
            placeholder="Es. +39 333 1234567"
            outlined dense
            :rules="[val => !!val || 'Indicare un recapito telefonico di emergenza']"
          >
            <template v-slot:prepend>
              <q-icon name="phone_in_talk" color="primary" />
            </template>
          </q-input>

          <!-- Note Alimentari -->
          <q-input
            v-model="consentForm.dietary_notes"
            label="Regime Alimentare / Intolleranze / Allergie"
            placeholder="Es. Celiachia (pranzo senza glutine), intolleranza lattosio..."
            outlined dense autogrow
          >
            <template v-slot:prepend>
              <q-icon name="restaurant" color="teal" />
            </template>
          </q-input>

          <!-- Note Mediche -->
          <q-input
            v-model="consentForm.medical_notes"
            label="Note Mediche / Farmaci Salvavita (Riservato ai docenti accompagnatori)"
            placeholder="Es. Necessità di portare autoiniettore adrenalina, asma da sforzo..."
            outlined dense autogrow
          >
            <template v-slot:prepend>
              <q-icon name="medical_services" color="purple" />
            </template>
          </q-input>

          <!-- Stato Quota -->
          <q-select
            v-model="consentForm.payment_status"
            :options="[
              { label: 'Quota Già Saldata (Bonifico/PagoPA)', value: 'paid' },
              { label: 'In Attesa di Versamento / Saldo', value: 'unpaid' }
            ]"
            emit-value map-options
            label="Stato Pagamento Quota Viaggio"
            outlined dense
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-1">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn
            color="primary"
            unelevated
            icon="verified_user"
            label="Firma e Autorizza"
            :loading="submittingId !== null"
            :disable="!consentForm.pin || !consentForm.emergency_phone"
            no-caps
            class="q-px-md rounded-lg font-bold"
            @click="submitConsentWithPin"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
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

const consentDialog = ref(false)
const selectedTripForConsent = ref(null)
const consentForm = reactive({
  pin: '',
  emergency_phone: '',
  dietary_notes: '',
  medical_notes: '',
  payment_status: 'paid'
})

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

const openConsentDialog = (trip) => {
  selectedTripForConsent.value = trip
  consentForm.pin = ''
  consentForm.emergency_phone = ''
  consentForm.dietary_notes = ''
  consentForm.medical_notes = ''
  consentForm.payment_status = 'paid'
  consentDialog.value = true
}

const submitConsentWithPin = async () => {
  if (!selectedTripForConsent.value) return
  const tripId = selectedTripForConsent.value.id
  submittingId.value = tripId
  try {
    const studentId = childrenStore.selectedChildId || childrenStore.children[0]?.id || childrenStore.children[0]?.userId
    await tripsService.submitConsentWithDetails(tripId, {
      student_id: studentId,
      status: 'granted',
      pin: consentForm.pin,
      emergency_phone: consentForm.emergency_phone,
      dietary_notes: consentForm.dietary_notes,
      medical_notes: consentForm.medical_notes,
      payment_status: consentForm.payment_status
    })
    $q.notify({ type: 'positive', message: 'Uscita autorizzata con successo con firma PIN!' })
    const trip = trips.value.find(t => t.id === tripId)
    if (trip) {
      trip.consent_status = 'granted'
      trip.signed = true
      trip.dietary_notes = consentForm.dietary_notes
      trip.medical_notes = consentForm.medical_notes
      trip.payment_status = consentForm.payment_status
    }
    consentDialog.value = false
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
