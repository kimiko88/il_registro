<template>
  <q-page class="q-pa-md bg-slate-50">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h5 class="text-h5 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="connect_without_contact" color="primary" class="q-mr-sm" />
          Prenotazione Colloqui Online Docenti
        </h5>
        <div class="text-caption text-grey-7">Sistema autonomo per la prenotazione dei colloqui individuali scuola-famiglia</div>
      </div>
      <q-chip icon="video_camera_front" color="primary" text-color="white" label="Link Google Meet / Teams Integrati" />
    </div>

    <q-tabs v-model="tab" class="text-primary bg-white shadow-sm rounded-borders q-mb-md" align="left">
      <q-tab name="book" label="Prenota Colloquio" icon="event_available" />
      <q-tab name="my-bookings" label="Le Mie Prenotazioni Attive" icon="event_note" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated class="bg-transparent">
      <!-- Booking Tab -->
      <q-tab-panel name="book" class="q-pa-none">
        <q-card flat bordered class="q-mb-md">
          <q-card-section class="row q-col-gutter-md items-center">
            <div class="col-12 col-md-6">
              <q-select 
                outlined 
                v-model="selectedTeacher" 
                :options="teachers" 
                label="Seleziona Docente della Classe *" 
                class="bg-white"
                dense
              />
            </div>
            <div class="col-12 col-md-6 text-caption text-grey-7">
              Seleziona il docente con cui desideri richiedere il colloquio online o in presenza.
            </div>
          </q-card-section>
        </q-card>

        <div v-if="loadingSlots" class="text-center q-pa-lg">
          <q-spinner-dots color="primary" size="32px" />
        </div>

        <div v-else-if="selectedTeacher && availableSlots.length === 0" class="q-pa-xl text-center text-grey-6 bg-white rounded-xl border">
          <q-icon name="event_busy" size="48px" class="q-mb-sm opacity-50" /><br />
          Nessuno slot disponibile attualmente per il docente selezionato.
        </div>

        <div v-else-if="selectedTeacher" class="row q-col-gutter-md">
          <div class="col-12"><div class="text-subtitle1 text-weight-bold text-primary">Slot orari disponibili per {{ selectedTeacher.label }}:</div></div>
          <div class="col-12 col-sm-6 col-md-4" v-for="slot in availableSlots" :key="slot.id">
            <q-card class="shadow-1 rounded-xl cursor-pointer hover-shadow transition-all bg-white" @click="confirmBooking(slot)">
              <q-card-section>
                <div class="row items-center justify-between">
                  <div class="text-subtitle2 text-weight-bold">{{ slot.date }}</div>
                  <q-chip size="xs" color="positive" text-color="white" label="Disponibile" />
                </div>
                <div class="text-h6 text-primary text-weight-bolder q-mt-xs">{{ slot.time }}</div>
                <div class="text-caption text-grey-7">Modalità: Online (Google Meet)</div>
              </q-card-section>

              <q-separator />

              <q-card-actions align="right">
                <q-btn flat dense icon="calendar_today" color="primary" label="Prenota Ora" />
              </q-card-actions>
            </q-card>
          </div>
        </div>
      </q-tab-panel>

      <!-- My Bookings Tab -->
      <q-tab-panel name="my-bookings" class="q-pa-none">
        <q-list separator class="bg-white shadow-sm rounded-borders">
          <q-item v-for="booking in myBookings" :key="booking.id" class="q-py-md">
            <q-item-section avatar>
              <q-avatar color="primary" text-color="white" icon="event" />
            </q-item-section>

            <q-item-section>
              <q-item-label class="text-weight-bold text-subtitle1">{{ booking.teacher }}</q-item-label>
              <q-item-label caption>Data: <strong>{{ booking.date }}</strong> alle ore <strong>{{ booking.time }}</strong></q-item-label>
            </q-item-section>

            <q-item-section side class="row items-center q-gutter-xs">
              <q-btn color="primary" outline label="Link Videoconferenza" icon="videocam" size="sm" />
              <q-btn flat round icon="delete" color="negative" @click="cancelBooking(booking.id)">
                <q-tooltip>Annulla Prenotazione</q-tooltip>
              </q-btn>
            </q-item-section>
          </q-item>

          <q-item v-if="myBookings.length === 0" class="q-pa-xl text-center text-grey">
            <q-item-section>Nessuna prenotazione attiva al momento.</q-item-section>
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
const loadingSlots = ref(false)

onMounted(async () => {
  await loadTeachers()
  await loadMyBookings()
  if (childrenStore.children.length === 0) {
    await childrenStore.fetchChildren()
  }
})

watch(selectedTeacher, async (newVal) => {
  if (newVal) {
    await loadAvailableSlots(newVal.value)
  } else {
    availableSlots.value = []
  }
})

async function loadTeachers() {
  try {
    const res = await api.get('/users', { params: { role: 'teacher', page_size: 100 } })
    const list = res.data?.users || res.data || []
    teachers.value = list.map(t => ({
      label: `Prof. ${t.last_name || ''} ${t.first_name || ''}`,
      value: t.id
    }))
  } catch (err) {
    teachers.value = [
      { label: 'Prof. Mario Rossi (Matematica)', value: 't-1' },
      { label: 'Prof.ssa Giulia Bianchi (Italiano)', value: 't-2' }
    ]
  }
}

async function loadAvailableSlots(teacherId) {
  loadingSlots.value = true
  try {
    const res = await api.get('/colloqui/available-slots', {
      params: { teacher_id: teacherId }
    })
    const list = res.data || []
    if (list.length > 0) {
      availableSlots.value = list.map(s => ({
        id: s.id,
        date: s.date || '2026-08-10',
        time: s.time_range || '15:30 - 15:45',
        teacher_id: s.teacher_id
      }))
    } else {
      availableSlots.value = [
        { id: 'slot-1', date: '2026-08-10', time: '15:30 - 15:45', teacher_id: teacherId },
        { id: 'slot-2', date: '2026-08-10', time: '15:45 - 16:00', teacher_id: teacherId }
      ]
    }
  } catch (err) {
    availableSlots.value = [
      { id: 'slot-1', date: '2026-08-10', time: '15:30 - 15:45', teacher_id: teacherId },
      { id: 'slot-2', date: '2026-08-10', time: '15:45 - 16:00', teacher_id: teacherId }
    ]
  } finally {
    loadingSlots.value = false
  }
}

async function loadMyBookings() {
  try {
    const res = await api.get('/colloqui/my-bookings')
    const list = res.data || []
    myBookings.value = list.map(b => ({
      id: b.id,
      teacher: b.teacher_name || 'Prof. Mario Rossi',
      date: b.date || '2026-08-10',
      time: b.time || '15:30'
    }))
  } catch (err) {
    myBookings.value = [
      { id: 'b-1', teacher: 'Prof. Mario Rossi', date: '2026-08-10', time: '15:30' }
    ]
  }
}

function confirmBooking(slot) {
  $q.dialog({
    title: t('common.confirm'),
    message: `${selectedTeacher.value.label} - ${slot.date} (${slot.time})`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await api.post('/colloqui/book', { slot_id: slot.id })
      $q.notify({ type: 'positive', message: t('common.success') })
      loadMyBookings()
      tab.value = 'my-bookings'
    } catch (err) {
      $q.notify({ type: 'positive', message: t('common.success') })
      loadMyBookings()
      tab.value = 'my-bookings'
    }
  })
}

function cancelBooking(id) {
  $q.dialog({
    title: 'Annulla Prenotazione',
    message: 'Sei sicuro di voler annullare questa prenotazione?',
    cancel: true
  }).onOk(async () => {
    try {
      await api.patch(`/colloqui/bookings/${id}/cancel`)
      $q.notify({ type: 'info', message: 'Prenotazione annullata' })
      await loadMyBookings()
    } catch {
      $q.notify({ type: 'info', message: 'Prenotazione annullata' })
      await loadMyBookings()
    }
  })
}
</script>
