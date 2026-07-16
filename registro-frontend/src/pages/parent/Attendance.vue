<template>
  <q-page class="q-pa-md bg-slate-50">
    <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">
      Presenze: {{ selectedChild?.first_name || '...' }}
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="60px" />
      <div class="text-grey q-mt-md">Caricamento presenze...</div>
    </div>

    <div v-else>
      <!-- Attendance Stats from DB -->
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-12 col-md-3">
          <q-card class="bg-green-1 text-positive shadow-sm" bordered>
            <q-card-section class="text-center">
              <div class="text-h4 text-weight-bold">{{ presenceRate }}%</div>
              <div class="text-subtitle2">Presenza Totale</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-md-3">
          <q-card class="bg-red-1 text-negative shadow-sm" bordered>
            <q-card-section class="text-center">
              <div class="text-h4 text-weight-bold">{{ summary.total_absences }}</div>
              <div class="text-subtitle2">Assenze Totali</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-md-3">
          <q-card class="bg-orange-1 text-orange shadow-sm" bordered>
            <q-card-section class="text-center">
              <div class="text-h4 text-weight-bold">{{ summary.total_lates }}</div>
              <div class="text-subtitle2">Ritardi</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-md-3">
          <q-card :class="unjustifiedCount > 0 ? 'bg-deep-orange-1 text-deep-orange' : 'bg-grey-2 text-grey'" class="shadow-sm" bordered>
            <q-card-section class="text-center">
              <div class="text-h4 text-weight-bold">{{ unjustifiedCount }}</div>
              <div class="text-subtitle2">Da Giustificare</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Risk badge -->
      <div v-if="summary.risk_level && summary.risk_level !== 'Normal'" class="q-mb-md">
        <q-banner
          :class="summary.risk_level === 'Critical' ? 'bg-negative text-white' : 'bg-warning text-white'"
          rounded dense
        >
          <template v-slot:avatar>
            <q-icon :name="summary.risk_level === 'Critical' ? 'error' : 'warning'" />
          </template>
          <span v-if="summary.risk_level === 'Critical'">
            ⚠️ Livello assenze critico ({{ summary.absence_rate?.toFixed(1) }}%). Contattare la scuola.
          </span>
          <span v-else>
            Livello assenze elevato ({{ summary.absence_rate?.toFixed(1) }}%). Monitorare la situazione.
          </span>
        </q-banner>
      </div>

      <!-- No events -->
      <q-card v-if="events.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="event_available" color="positive" size="80px" class="q-mb-md" />
        <div class="text-h6 text-positive">Nessuna assenza o ritardo negli ultimi 30 giorni</div>
        <div class="text-caption">Ottimo! Continua così.</div>
      </q-card>

      <!-- Events List -->
      <q-list v-else bordered class="bg-white rounded-borders shadow-sm">
        <q-item-label header class="text-weight-bold">Ultimi 30 giorni — Assenze e Ritardi</q-item-label>

        <template v-for="event in events" :key="event.id">
          <q-item>
            <q-item-section avatar>
              <q-icon
                :name="getIconName(event.status)"
                :color="getIconColor(event.status)"
                size="md"
              />
            </q-item-section>

            <q-item-section>
              <q-item-label class="text-weight-medium">
                {{ formatDate(event.date) }} — {{ event.displayStatus }}
              </q-item-label>
              <q-item-label caption v-if="event.notes" class="text-grey-7">
                <q-icon name="notes" size="xs" class="q-mr-xs" />{{ event.notes }}
              </q-item-label>
              <q-item-label caption>
                <q-badge
                  v-if="event.is_justified"
                  color="positive"
                  label="Giustificata"
                />
                <q-badge
                  v-else
                  color="negative"
                  label="Da Giustificare"
                />
              </q-item-label>
            </q-item-section>

            <q-item-section side v-if="!event.is_justified">
              <q-btn outline color="primary" size="sm" label="Giustifica" @click="openJustifyDialog(event)" />
            </q-item-section>
          </q-item>
          <q-separator spaced inset />
        </template>
      </q-list>
    </div>

    <!-- Justify Dialog -->
    <q-dialog v-model="justifyDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Giustifica Assenza</div>
          <div class="text-caption text-grey">del {{ selectedEvent?.date }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-select v-model="justifyReason" :options="['Salute', 'Motivi Familiari', 'Altro']" label="Motivazione" />
          <q-input v-model="justifyNotes" label="Note (Opzionale)" type="textarea" autogrow class="q-mt-sm" />
          <div class="q-mt-md text-caption text-grey">
            Cliccando "Conferma", dichiari di essere il genitore/tutore legale.
          </div>
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn flat label="Conferma" :loading="submitting" @click="submitJustification" />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'
import { useQuasar, date as qdate } from 'quasar'
import { attendanceService } from 'src/services/attendanceService'

const $q = useQuasar()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

const loading = ref(false)
const submitting = ref(false)
const justifyDialog = ref(false)
const selectedEvent = ref(null)
const justifyReason = ref('Salute')
const justifyNotes = ref('')

const events = ref([])
const summary = ref({
    total_absences: 0,
    total_lates: 0,
    total_early_exits: 0,
    justified_count: 0,
    absence_rate: 0,
    risk_level: 'Normal',
    pending_requests: 0
})

// Computed
const unjustifiedCount = computed(() =>
    events.value.filter(e => !e.is_justified).length
)

const presenceRate = computed(() => {
    const rate = summary.value.absence_rate
    if (!rate && rate !== 0) return '—'
    return (100 - rate).toFixed(1)
})

onMounted(() => {
    if (selectedChild.value) {
        fetchAll()
    }
})

watch(selectedChild, (val) => {
    if (val) fetchAll()
})

const fetchAll = async () => {
    loading.value = true
    try {
        await Promise.all([fetchAttendance(), fetchSummary()])
    } finally {
        loading.value = false
    }
}

const fetchAttendance = async () => {
    try {
        const res = await attendanceService.getChildAttendance(selectedChild.value.id)
        const records = Array.isArray(res.data) ? res.data : (res.data?.records || [])
        events.value = records
            .filter(e => {
                const s = String(e.status).toLowerCase()
                return s !== 'present' && s !== 'exempt'
            })
            .map(e => {
                let timeDetail = ''
                const s = String(e.status).toLowerCase()
                if (s === 'late') {
                    timeDetail = e.entry_time ? ` (Ingresso: ${e.entry_time})` : ' (Ritardo)'
                } else if (s === 'leftearly' || s === 'early') {
                    timeDetail = e.exit_time ? ` (Uscita: ${e.exit_time})` : ' (Uscita anticipata)'
                }
                return {
                    id: e.id,
                    date: e.date,
                    status: e.status,
                    displayStatus: mapStatus(e.status) + timeDetail,
                    is_justified: e.is_justified,
                    notes: e.notes || ''
                }
            })
    } catch (e) {
        console.error('Errore caricamento presenze:', e)
    }
}

const fetchSummary = async () => {
    try {
        const res = await attendanceService.getChildAttendanceSummary(selectedChild.value.id)
        if (res.data) {
            summary.value = res.data
        }
    } catch (e) {
        console.error('Errore caricamento riepilogo:', e)
    }
}

function formatDate(d) {
    if (!d) return ''
    try {
        return qdate.formatDate(new Date(d), 'DD/MM/YYYY')
    } catch {
        return d
    }
}

function mapStatus(status) {
    const s = String(status).toLowerCase()
    if (s === 'absent') return 'Assenza'
    if (s === 'late') return 'Ritardo'
    if (s === 'leftearly' || s === 'early') return 'Uscita Anticipata'
    return status
}

function openJustifyDialog(event) {
    selectedEvent.value = event
    justifyReason.value = 'Salute'
    justifyNotes.value = ''
    justifyDialog.value = true
}

async function submitJustification() {
    submitting.value = true
    try {
        await attendanceService.justify(selectedEvent.value.id, {
            student_id: selectedChild.value.id,
            start_date: selectedEvent.value.date,
            end_date: selectedEvent.value.date,
            reason: justifyReason.value + (justifyNotes.value ? ` - ${justifyNotes.value}` : '')
        })
        $q.notify({ type: 'positive', message: 'Giustificazione inviata con successo' })
        justifyDialog.value = false
        await fetchAll()
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore invio giustificazione' })
    } finally {
        submitting.value = false
    }
}

function getIconName(status) {
    const s = String(status).toLowerCase()
    if (s === 'absent') return 'cancel'
    if (s === 'late') return 'schedule'
    if (s === 'leftearly' || s === 'early') return 'logout'
    return 'help'
}

function getIconColor(status) {
    const s = String(status).toLowerCase()
    if (s === 'absent') return 'negative'
    if (s === 'late') return 'warning'
    if (s === 'leftearly' || s === 'early') return 'blue'
    return 'grey'
}
</script>
