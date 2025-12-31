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

      <q-item v-for="event in mockEvents" :key="event.id">
        <q-item-section avatar>
          <q-icon 
            :name="event.type === 'absence' ? 'cancel' : 'schedule'" 
            :color="event.type === 'absence' ? 'negative' : 'warning'" 
          />
        </q-item-section>
        
        <q-item-section>
          <q-item-label>{{ event.date }} - {{ event.type === 'absence' ? 'Assenza' : 'Ritardo' }}</q-item-label>
          <q-item-label caption v-if="event.justified" class="text-positive">Giustificata</q-item-label>
          <q-item-label caption v-else class="text-negative">Da Giustificare</q-item-label>
        </q-item-section>

        <q-item-section side v-if="!event.justified">
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
import { ref } from 'vue'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

const justifyDialog = ref(false)
const selectedEvent = ref(null)
const reason = ref('Salute')
const notes = ref('')

const mockEvents = ref([
  { id: 1, date: '12/12/2024', type: 'absence', justified: false },
  { id: 2, date: '01/12/2024', type: 'delay', justified: true },
  { id: 3, date: '28/11/2024', type: 'absence', justified: true },
])

function openJustifyDialog(event) {
  selectedEvent.value = event
  reason.value = 'Salute'
  notes.value = ''
  justifyDialog.value = true
}

function submitJustification() {
  // Mock API Call
  setTimeout(() => {
    const idx = mockEvents.value.findIndex(e => e.id === selectedEvent.value.id)
    if (idx !== -1) mockEvents.value[idx].justified = true
    
    $q.notify({ type: 'positive', message: 'Giustificazione inviata con successo' })
    justifyDialog.value = false
  }, 500)
}
</script>
