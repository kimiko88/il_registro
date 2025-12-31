<template>
  <q-page class="q-pa-md bg-slate-50">
    <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">Prenotazione Colloqui</div>

    <q-tabs v-model="tab" class="text-primary bg-white shadow-sm rounded-borders q-mb-md" align="left">
      <q-tab name="book" label="Prenota" icon="event" />
      <q-tab name="my-bookings" label="Le Mie Prenotazioni" icon="list" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated class="bg-transparent">
      
      <!-- Booking Tab -->
      <q-tab-panel name="book" class="q-pa-none">
        <!-- Teacher Filter -->
        <q-select 
          outlined 
          v-model="selectedTeacher" 
          :options="teachers" 
          label="Seleziona Docente" 
          class="bg-white q-mb-md" 
        />

        <div v-if="selectedTeacher" class="text-subtitle1 q-mb-sm">Slot Disponibili per {{ selectedTeacher.label }}</div>
        
        <div class="row q-col-gutter-sm">
          <div class="col-12 col-sm-6 col-md-4" v-for="slot in availableSlots" :key="slot">
             <q-card class="text-center shadow-sm cursor-pointer hover:shadow-md" @click="confirmBooking(slot)">
               <q-card-section>
                 <div class="text-h6">{{ slot.date }}</div>
                 <div class="text-primary text-weight-bold">{{ slot.time }}</div>
               </q-card-section>
             </q-card>
          </div>
        </div>
      </q-tab-panel>

      <!-- My Bookings Tab -->
      <q-tab-panel name="my-bookings" class="q-pa-none">
        <q-list separator class="bg-white shadow-sm rounded-borders">
          <q-item v-for="booking in myBookings" :key="booking.id">
            <q-item-section>
              <q-item-label class="text-weight-bold">{{ booking.teacher }}</q-item-label>
              <q-item-label caption>{{ booking.date }} alle {{ booking.time }}</q-item-label>
            </q-item-section>
            <q-item-section side>
               <q-btn flat round icon="delete" color="negative" @click="cancelBooking(booking.id)">
                 <q-tooltip>Annulla</q-tooltip>
               </q-btn>
            </q-item-section>
          </q-item>
          <q-item v-if="myBookings.length === 0">
            <q-item-section class="text-center text-grey">Nessuna prenotazione attiva</q-item-section>
          </q-item>
        </q-list>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tab = ref('book')
const selectedTeacher = ref(null)

const teachers = [
  { label: 'Prof. Verdi (Matematica)', value: 1 },
  { label: 'Prof. Bianchi (Storia)', value: 2 },
  { label: 'Prof. Neri (Inglese)', value: 3 },
]

const availableSlots = ref([
  { date: 'Lun 15 Gen', time: '16:00' },
  { date: 'Lun 15 Gen', time: '16:15' },
  { date: 'Mer 17 Gen', time: '10:00' },
])

const myBookings = ref([
  { id: 101, teacher: 'Prof. Bianchi (Storia)', date: 'Ven 12 Gen', time: '11:00' }
])

function confirmBooking(slot) {
  $q.dialog({
    title: 'Conferma Prenotazione',
    message: `Vuoi prenotare il colloquio con ${selectedTeacher.value.label} per ${slot.date} alle ${slot.time}?`,
    cancel: true,
    persistent: true
  }).onOk(() => {
    myBookings.value.push({ id: Date.now(), teacher: selectedTeacher.value.label, ...slot })
    $q.notify({ type: 'positive', message: 'Prenotazione confermata!' })
    tab.value = 'my-bookings'
  })
}

function cancelBooking(id) {
  $q.dialog({
    title: 'Annulla',
    message: 'Sei sicuro di voler annullare?',
    cancel: true
  }).onOk(() => {
    myBookings.value = myBookings.value.filter(b => b.id !== id)
    $q.notify({ type: 'info', message: 'Prenotazione annullata' })
  })
}
</script>
