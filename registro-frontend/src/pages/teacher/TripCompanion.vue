<template>
  <q-page class="q-pa-md" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-grey-1 text-dark'">
    <!-- Top Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bolder row items-center text-primary">
          <q-icon name="directions_bus" class="q-mr-sm" size="28px" />
          {{ t('tripCompanion.title') }}
        </div>
        <div class="text-caption" :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-7'">
          {{ t('tripCompanion.subtitle') }}
        </div>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-btn
          color="secondary"
          icon="download"
          :label="t('tripCompanion.exportRoster')"
          outline dense
          class="q-px-sm rounded-lg"
          @click="exportRosterCSV"
          :disable="!selectedTrip || consents.length === 0"
        />
        <q-btn
          color="primary"
          icon="refresh"
          :label="t('common.refresh')"
          unelevated dense
          class="q-px-sm rounded-lg"
          :loading="loadingConsents"
          @click="loadConsents"
          :disable="!selectedTripId"
        />
      </div>
    </div>

    <!-- Trip Selector & Summary -->
    <q-card flat bordered class="rounded-xl q-mb-md shadow-1" :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-white'">
      <q-card-section class="row q-col-gutter-sm items-center q-py-sm">
        <div class="col-12 col-md-5">
          <q-select
            v-model="selectedTripId"
            :options="tripOptions"
            option-value="id"
            option-label="label"
            emit-value map-options
            dense outlined
            :label="t('tripCompanion.selectTrip')"
            :loading="loadingTrips"
            @update:model-value="onTripChanged"
          >
            <template v-slot:prepend>
              <q-icon name="tour" color="primary" />
            </template>
          </q-select>
        </div>

        <div class="col-12 col-md-4" v-if="selectedTrip">
          <div class="text-caption text-grey-7">
            <q-icon name="place" color="primary" size="xs" /> <strong>{{ selectedTrip.destination }}</strong>
            <span class="q-ml-sm">
              <q-icon name="event" size="xs" /> {{ formatDate(selectedTrip.departure_date) }}
            </span>
          </div>
          <div class="text-caption text-grey-6 ellipsis" v-if="selectedTrip.accompanying_teachers">
            Docenti: {{ selectedTrip.accompanying_teachers }}
          </div>
        </div>

        <div class="col-12 col-md-3 text-right" v-if="selectedTrip">
          <q-input
            v-model="studentFilter"
            dense outlined
            placeholder="Cerca studente..."
            bg-color="white"
          >
            <template v-slot:append>
              <q-icon name="search" size="xs" />
            </template>
          </q-input>
        </div>
      </q-card-section>
    </q-card>

    <!-- Stats summary badges -->
    <div class="row q-col-gutter-sm q-mb-md" v-if="selectedTrip">
      <div class="col-6 col-sm-3">
        <q-card flat bordered class="bg-green-50 border-green-200 text-center q-pa-sm rounded-xl">
          <div class="text-h5 text-weight-bolder text-positive">{{ stats.authorized }}</div>
          <div class="text-caption text-grey-8">{{ t('tripCompanion.authorizedWithPin') }}</div>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card flat bordered class="bg-amber-50 border-amber-200 text-center q-pa-sm rounded-xl">
          <div class="text-h5 text-weight-bolder text-warning">{{ stats.pending }}</div>
          <div class="text-caption text-grey-8">{{ t('tripCompanion.pendingAuth') }}</div>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card flat bordered class="bg-teal-50 border-teal-200 text-center q-pa-sm rounded-xl">
          <div class="text-h5 text-weight-bolder text-teal-8">{{ stats.paid }}</div>
          <div class="text-caption text-grey-8">{{ t('tripCompanion.paidQuota') }}</div>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card flat bordered class="bg-purple-50 border-purple-200 text-center q-pa-sm rounded-xl">
          <div class="text-h5 text-weight-bolder text-purple-8">{{ stats.dietaryOrMedical }}</div>
          <div class="text-caption text-grey-8">{{ t('tripCompanion.specialNeeds') }}</div>
        </q-card>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loadingConsents" class="text-center q-pa-xl">
      <q-spinner color="primary" size="48px" />
      <div class="text-subtitle2 text-grey-6 q-mt-md">{{ t('tripCompanion.loadingConsents') }}</div>
    </div>

    <!-- No trip selected -->
    <div v-else-if="!selectedTripId" class="text-center q-pa-xl bg-white rounded-xl shadow-1">
      <q-icon name="directions_bus" size="64px" color="grey-4" />
      <div class="text-h6 text-grey-7 q-mt-md">{{ t('tripCompanion.promptSelectTrip') }}</div>
    </div>

    <!-- Consents List Table -->
    <q-card v-else flat bordered class="rounded-xl shadow-1 overflow-hidden">
      <q-table
        :rows="filteredConsents"
        :columns="columns"
        row-key="id"
        flat
        dense
        :pagination="{ rowsPerPage: 50 }"
        class="trip-table"
      >
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-chip
              v-if="props.row.status === 'granted' || props.row.status === 'Consented'"
              color="positive"
              text-color="white"
              dense size="sm"
              icon="verified_user"
            >
              {{ t('tripCompanion.statusGranted') }}
            </q-chip>
            <q-chip
              v-else-if="props.row.status === 'denied'"
              color="negative"
              text-color="white"
              dense size="sm"
              icon="cancel"
            >
              {{ t('tripCompanion.statusDenied') }}
            </q-chip>
            <q-chip
              v-else
              color="warning"
              text-color="black"
              dense size="sm"
              icon="hourglass_empty"
            >
              {{ t('tripCompanion.statusPending') }}
            </q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-pin="props">
          <q-td :props="props" class="text-center">
            <q-icon
              v-if="props.row.pin_verified"
              name="verified"
              color="positive"
              size="20px"
            >
              <q-tooltip>{{ t('tripCompanion.pinVerifiedTooltip') }}</q-tooltip>
            </q-icon>
            <span v-else class="text-caption text-grey-5">—</span>
          </q-td>
        </template>

        <template v-slot:body-cell-payment="props">
          <q-td :props="props">
            <q-badge
              :color="props.row.payment_status === 'paid' ? 'positive' : 'warning'"
              text-color="white"
              class="q-pa-xs"
            >
              {{ props.row.payment_status === 'paid' ? 'Saldata' : 'In Sospeso' }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-emergency="props">
          <q-td :props="props">
            <a
              v-if="props.row.emergency_phone"
              :href="'tel:' + props.row.emergency_phone"
              class="text-primary text-weight-bold text-decoration-none row items-center"
            >
              <q-icon name="call" size="xs" class="q-mr-xs" />
              {{ props.row.emergency_phone }}
            </a>
            <span v-else class="text-caption text-grey-4">—</span>
          </q-td>
        </template>

        <template v-slot:body-cell-dietary="props">
          <q-td :props="props">
            <div v-if="props.row.dietary_notes" class="row items-center text-teal-9 text-caption">
              <q-icon name="restaurant" size="14px" class="q-mr-xs text-teal" />
              <span class="ellipsis" :title="props.row.dietary_notes">{{ props.row.dietary_notes }}</span>
            </div>
            <span v-else class="text-caption text-grey-4">—</span>
          </q-td>
        </template>

        <template v-slot:body-cell-medical="props">
          <q-td :props="props">
            <div v-if="props.row.medical_notes" class="row items-center text-purple-9 text-caption">
              <q-icon name="medical_services" size="14px" class="q-mr-xs text-purple" />
              <span class="ellipsis" :title="props.row.medical_notes">{{ props.row.medical_notes }}</span>
            </div>
            <span v-else class="text-caption text-grey-4">—</span>
          </q-td>
        </template>
      </q-table>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import tripsService from '@/services/tripsService'

const { t } = useI18n()
const $q = useQuasar()

const trips = ref([])
const loadingTrips = ref(false)
const selectedTripId = ref(null)

const consents = ref([])
const loadingConsents = ref(false)
const studentFilter = ref('')

const tripOptions = computed(() => {
  return trips.value.map(tr => ({
    id: tr.id,
    label: `${tr.title} — ${tr.destination} (${formatDate(tr.departure_date)})`
  }))
})

const selectedTrip = computed(() => {
  return trips.value.find(tr => tr.id === selectedTripId.value) || null
})

const columns = [
  { name: 'student_name', label: 'Alunno', field: 'student_name', align: 'left', sortable: true },
  { name: 'parent_name', label: 'Genitore Tutore', field: 'parent_name', align: 'left', sortable: true },
  { name: 'status', label: 'Stato Consenso', field: 'status', align: 'center', sortable: true },
  { name: 'pin', label: 'PIN Verificato', field: 'pin_verified', align: 'center', sortable: true },
  { name: 'payment', label: 'Quota', field: 'payment_status', align: 'center', sortable: true },
  { name: 'emergency', label: 'Telefono Emergenza', field: 'emergency_phone', align: 'left' },
  { name: 'dietary', label: 'Regime Alimentare / Intolleranze', field: 'dietary_notes', align: 'left' },
  { name: 'medical', label: 'Note Mediche / Farmaci', field: 'medical_notes', align: 'left' }
]

const filteredConsents = computed(() => {
  if (!studentFilter.value.trim()) return consents.value
  const q = studentFilter.value.toLowerCase()
  return consents.value.filter(c =>
    (c.student_name && c.student_name.toLowerCase().includes(q)) ||
    (c.parent_name && c.parent_name.toLowerCase().includes(q)) ||
    (c.dietary_notes && c.dietary_notes.toLowerCase().includes(q)) ||
    (c.medical_notes && c.medical_notes.toLowerCase().includes(q))
  )
})

const stats = computed(() => {
  let authorized = 0
  let pending = 0
  let paid = 0
  let dietaryOrMedical = 0

  for (const c of consents.value) {
    if (c.status === 'granted' || c.status === 'Consented') {
      authorized++
    } else {
      pending++
    }
    if (c.payment_status === 'paid') paid++
    if (c.dietary_notes || c.medical_notes) dietaryOrMedical++
  }

  return { authorized, pending, paid, dietaryOrMedical }
})

onMounted(async () => {
  await fetchTrips()
})

function formatDate(val) {
  if (!val) return ''
  try {
    const d = new Date(val)
    return isNaN(d.getTime()) ? val : d.toLocaleDateString('it-IT')
  } catch {
    return val
  }
}

async function fetchTrips() {
  loadingTrips.value = true
  try {
    const res = await tripsService.getTrips()
    trips.value = Array.isArray(res.data) ? res.data : (res.data?.trips || [])
    if (trips.value.length > 0) {
      selectedTripId.value = trips.value[0].id
      await loadConsents()
    }
  } catch (err) {
    console.error('Error fetching trips', err)
    trips.value = []
  } finally {
    loadingTrips.value = false
  }
}

async function onTripChanged() {
  await loadConsents()
}

async function loadConsents() {
  if (!selectedTripId.value) return
  loadingConsents.value = true
  try {
    const res = await tripsService.getConsents(selectedTripId.value)
    consents.value = Array.isArray(res.data) ? res.data : (res.data?.consents || [])
  } catch (err) {
    console.error('Error loading consents', err)
    consents.value = []
  } finally {
    loadingConsents.value = false
  }
}

function exportRosterCSV() {
  if (consents.value.length === 0) return
  let csv = 'Alunno;Genitore;Stato;PIN Verificato;Quota;Telefono Emergenza;Dieta Speciale;Note Mediche\n'
  consents.value.forEach(c => {
    csv += `"${c.student_name || ''}";"${c.parent_name || ''}";"${c.status || ''}";"${c.pin_verified ? 'SI' : 'NO'}";"${c.payment_status || 'unpaid'}";"${c.emergency_phone || ''}";"${c.dietary_notes || ''}";"${c.medical_notes || ''}"\n`
  })
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.setAttribute('download', `registro_viaggio_${selectedTrip.value?.destination || 'uscita'}.csv`)
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  window.URL.revokeObjectURL(url)
  $q.notify({ type: 'positive', message: 'Elenco partecipanti esportato con successo!' })
}
</script>

<style scoped>
.trip-table {
  background: white;
}
</style>
