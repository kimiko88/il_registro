<template>
  <q-page class="q-pa-md bg-slate-50">
    <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">
      Presenze: {{ selectedChild?.firstName || '...' }}
    </div>

    <!-- Attendance Stats -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-4">
        <q-card class="bg-blue-1 text-primary shadow-sm" bordered>
          <q-card-section class="text-center">
            <div class="text-h4 text-weight-bold">92%</div>
            <div class="text-subtitle2">Presenza Totale</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-orange-1 text-orange shadow-sm" bordered>
          <q-card-section class="text-center">
             <div class="text-h4 text-weight-bold">3</div>
             <div class="text-subtitle2">Assenze da Giustificare</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Absences List -->
    <q-list bordered class="bg-white rounded-borders shadow-sm">
      <q-item-label header class="text-weight-bold">Ultimi Eventi</q-item-label>

      <q-item v-for="event in events" :key="event.id">
        <q-item-section avatar>
          <q-icon 
            :name="event.status === 'absent' ? 'cancel' : 'schedule'" 
            :color="event.status === 'absent' ? 'negative' : 'warning'" 
          />
        </q-item-section>
        
        <q-item-section>
          <q-item-label>{{ event.date }} - {{ event.status === 'absent' ? 'Assenza' : 'Ritardo' }}</q-item-label>
          <q-item-label caption v-if="event.is_justified" class="text-positive">Giustificata</q-item-label>
          <q-item-label caption v-else class="text-negative">Da Giustificare</q-item-label>
        </q-item-section>

        <q-item-section side v-if="!event.is_justified">
          <q-btn outline color="primary" size="sm" label="Giustifica" @click="openJustifyDialog(event)" />
        </q-item-section>
      </q-item>
    </q-list>

    <!-- Justify Dialog -->
    <q-dialog v-model="justifyDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Giustifica Assenza</div>
          <div class="text-caption">del {{ selectedEvent?.date }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-select v-model="reason" :options="['Salute', 'Motivi Familiari', 'Altro']" label="Motivazione" />
          <q-input v-model="notes" label="Note (Opzionale)" type="textarea" autogrow />
          <div class="q-mt-md text-caption text-grey">
            Cliccando "Conferma", dichiari di essere il genitore/tutore legale.
          </div>
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn flat label="Conferma" @click="submitJustification" />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { attendanceService } from 'src/services/attendanceService'

const $q = useQuasar()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

const justifyDialog = ref(false)
const selectedEvent = ref(null)
const reason = ref('Salute')
const notes = ref('')

const events = ref([])

onMounted(() => {
    if (selectedChild.value) {
        fetchAttendance()
    }
})

watch(selectedChild, (val) => {
    if (val) fetchAttendance()
})

const fetchAttendance = async () => {
    try {
        const res = await attendanceService.getChildAttendance(selectedChild.value.id)
        // Backend returns generic attendance list. Filter for absence/late if needed or show all?
        // Mock UI showed "Ultimi Eventi" (Absences/Lays).
        // Let's filter for non-present statuses for the list.
        events.value = (res.data || []).filter(e => e.status !== 'present')
    } catch (e) {
        console.error(e)
    }
}

function openJustifyDialog(event) {
  selectedEvent.value = event
  reason.value = 'Salute'
  notes.value = ''
  justifyDialog.value = true
}

async function submitJustification() {
  try {
      // API call to justify
      await attendanceService.justify(selectedEvent.value.id, {
          student_id: selectedChild.value.id,
          start_date: selectedEvent.value.date,
          end_date: selectedEvent.value.date,
          reason: reason.value + (notes.value ? ` - ${notes.value}` : '')
      })
      
      $q.notify({ type: 'positive', message: 'Giustificazione inviata con successo' })
      justifyDialog.value = false
      // Refetch to update status
      fetchAttendance() 
  } catch (e) {
      $q.notify({ type: 'negative', message: 'Errore invio giustificazione' })
  }
}
</script>
