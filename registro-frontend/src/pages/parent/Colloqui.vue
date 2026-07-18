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
import { ref, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import api from '@/services/api'
import { useChildrenStore } from '@/stores/children'

const $q = useQuasar()
const childrenStore = useChildrenStore()

const tab = ref('book')
const selectedTeacher = ref(null)
const teachers = ref([])
const availableSlots = ref([])
const myBookings = ref([])

onMounted(async () => {
  await loadTeachers()
  await loadMyBookings()
  if (childrenStore.children.length === 0) {
    await childrenStore.fetchChildren()
  }
})

// Reload slots when selected teacher changes
watch(selectedTeacher, async (newVal) => {
  if (newVal) {
    await loadAvailableSlots(newVal.value)
  } else {
    availableSlots.value = []
  }
})

async function loadTeachers() {
  try {
    const res = await api.get('/teachers')
    teachers.value = (res.data || []).map(t => ({
      label: `Prof. ${t.last_name || ''} ${t.first_name || ''} (${t.qualification || 'Docente'})`,
      value: t.id
    }))
  } catch (err) {
    console.error('Failed to load teachers:', err)
  }
}

async function loadAvailableSlots(teacherId) {
  try {
    const res = await api.get('/colloqui/available-slots', {
      params: { teacher_id: teacherId }
    })
    availableSlots.value = (res.data || []).map(s => ({
      id: s.id,
      date: s.date,
      time: s.time_range,
      teacher_id: s.teacher_id
    }))
  } catch (err) {
    console.error('Failed to load available slots:', err)
    $q.notify({ type: 'negative', message: 'Errore caricamento slot' })
  }
}

async function loadMyBookings() {
  try {
    const res = await api.get('/colloqui/my-bookings')
    myBookings.value = (res.data || []).map(b => {
      // Find teacher name from loaded list or default
      const teacher = teachers.value.find(t => t.value === b.slot_info.teacher_id)
      return {
        id: b.id,
        teacher: teacher ? teacher.label : `Docente (Ref: ${b.slot_info.teacher_id.slice(0, 8)})`,
        date: b.slot_info.date,
        time: b.slot_info.time_range
      }
    })
  } catch (err) {
    console.error('Failed to load my bookings:', err)
  }
}

function confirmBooking(slot) {
  if (!childrenStore.selectedChildId) {
    $q.notify({ type: 'warning', message: 'Seleziona prima uno studente/figlio' })
    return
  }
  $q.dialog({
    title: 'Conferma Prenotazione',
    message: `Vuoi prenotare il colloquio con ${selectedTeacher.value.label} per il ${slot.date} alle ${slot.time}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await api.post('/colloqui/book', {
        slot_id: slot.id,
        student_id: childrenStore.selectedChildId
      })
      $q.notify({ type: 'positive', message: 'Prenotazione confermata!' })
      await loadMyBookings()
      tab.value = 'my-bookings'
      // Refresh available slots for selected teacher
      if (selectedTeacher.value) {
        await loadAvailableSlots(selectedTeacher.value.value)
      }
    } catch (err) {
      console.error('Failed to book slot:', err)
      $q.notify({ type: 'negative', message: 'Errore prenotazione' })
    }
  })
}

function cancelBooking(id) {
  $q.dialog({
    title: 'Annulla',
    message: 'Sei sicuro di voler annullare?',
    cancel: true
  }).onOk(async () => {
    try {
      await api.patch(`/colloqui/bookings/${id}/cancel`)
      $q.notify({ type: 'info', message: 'Prenotazione annullata' })
      await loadMyBookings()
      if (selectedTeacher.value) {
        await loadAvailableSlots(selectedTeacher.value.value)
      }
    } catch (err) {
      console.error('Failed to cancel booking:', err)
      $q.notify({ type: 'negative', message: 'Errore annullamento' })
    }
  })
}
</script>
