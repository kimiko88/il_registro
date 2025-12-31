<template>
  <q-page class="q-pa-md bg-grey-1">
    <div v-if="!selectedChild" class="text-center q-pa-xl text-grey">
        <div class="text-h6">Seleziona un figlio dalla Dashboard.</div>
    </div>
    <div v-else>
        <div class="row items-center justify-between q-mb-md">
           <div class="text-h4">Colloqui - {{ selectedChild.name }}</div>
           <q-btn label="Prenota Colloquio" color="primary" icon="event" @click="showBookingDialog = true" />
        </div>

        <div class="row q-col-gutter-lg">
            <!-- Upcoming Appointments -->
            <div class="col-12 col-md-8">
                <q-card>
                    <q-card-section class="text-h6">I Miei Appuntamenti</q-card-section>
                    <q-separator />
                    <q-list separator>
                        <q-item v-for="booking in bookings" :key="booking.id">
                            <q-item-section avatar>
                                <q-avatar color="primary" text-color="white" icon="school" />
                            </q-item-section>
                            <q-item-section>
                                <q-item-label class="text-weight-bold">{{ booking.teacher }} ({{ booking.subject }})</q-item-label>
                                <q-item-label caption>
                                    <q-icon name="event" /> {{ booking.date }} - {{ booking.time }}
                                    <span v-if="booking.isOnline" class="text-blue q-ml-sm">(Online)</span>
                                </q-item-label>
                            </q-item-section>
                            <q-item-section side>
                                <div class="row q-gutter-xs">
                                    <q-btn flat round icon="videocam" color="green" v-if="booking.isOnline" type="a" :href="booking.link" target="_blank">
                                        <q-tooltip>Link Meet</q-tooltip>
                                    </q-btn>
                                    <q-btn flat round icon="cancel" color="red" @click="cancelBooking(booking)">
                                        <q-tooltip>Annulla</q-tooltip>
                                    </q-btn>
                                </div>
                            </q-item-section>
                        </q-item>
                        <q-item v-if="bookings.length === 0">
                            <q-item-section class="text-center text-grey">Nessun appuntamento in programma.</q-item-section>
                        </q-item>
                    </q-list>
                </q-card>
            </div>

            <!-- Teachers List Sidebar -->
            <div class="col-12 col-md-4">
                 <q-card class="bg-blue-1">
                     <q-card-section>
                         <div class="text-subtitle1 text-weight-bold">Docenti della Classe</div>
                         <q-list dense class="q-mt-sm">
                             <q-item v-for="teacher in teachers" :key="teacher.id">
                                 <q-item-section avatar><q-avatar icon="person" size="sm" color="blue-2" text-color="blue-9" /></q-item-section>
                                 <q-item-section>
                                     <q-item-label>{{ teacher.name }}</q-item-label>
                                     <q-item-label caption>{{ teacher.subject }}</q-item-label>
                                 </q-item-section>
                             </q-item>
                         </q-list>
                     </q-card-section>
                 </q-card>
            </div>
        </div>
    </div>

    <!-- Booking Dialog -->
    <q-dialog v-model="showBookingDialog">
        <q-card style="min-width: 500px">
            <q-card-section class="text-h6">Nuova Prenotazione</q-card-section>
            <q-card-section>
                <q-stepper v-model="step" vertical color="primary" animated header-nav>
                    <q-step :name="1" title="Seleziona Docente" icon="person" :done="step > 1">
                        <q-select v-model="newBooking.teacher" :options="teachers" option-label="name" label="Docente" outlined />
                        <q-stepper-navigation>
                            <q-btn @click="step = 2" color="primary" label="Avanti" :disable="!newBooking.teacher" />
                        </q-stepper-navigation>
                    </q-step>

                    <q-step :name="2" title="Scegli Orario" icon="access_time" :done="step > 2">
                        <div class="q-gutter-sm row">
                            <q-btn 
                                v-for="slot in availableSlots" 
                                :key="slot.id" 
                                :label="slot.label" 
                                :color="newBooking.slot === slot ? 'primary' : 'grey-3'" 
                                :text-color="newBooking.slot === slot ? 'white' : 'black'"
                                @click="newBooking.slot = slot"
                            />
                        </div>
                        <q-stepper-navigation>
                            <q-btn @click="step = 3" color="primary" label="Avanti" :disable="!newBooking.slot" />
                            <q-btn flat @click="step = 1" label="Indietro" />
                        </q-stepper-navigation>
                    </q-step>

                    <q-step :name="3" title="Conferma" icon="check">
                        <div v-if="newBooking.teacher && newBooking.slot">
                            Confermi la prenotazione con <strong>{{ newBooking.teacher.name }}</strong> per il <strong>{{ newBooking.slot.label }}</strong>?
                        </div>
                        <q-stepper-navigation>
                            <q-btn color="green" label="Conferma Prenotazione" @click="confirmBooking" />
                            <q-btn flat @click="step = 2" label="Indietro" />
                        </q-stepper-navigation>
                    </q-step>
                </q-stepper>
            </q-card-section>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useParentStore } from 'src/stores/parent'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const parentStore = useParentStore()
const selectedChild = computed(() => parentStore.selectedChild)

const showBookingDialog = ref(false)
const step = ref(1)

const bookings = ref([
    { id: 101, teacher: 'Prof. Bianchi', subject: 'Matematica', date: '21/01/2025', time: '16:00', isOnline: true, link: '#' }
])

const teachers = [
    { id: 1, name: 'Prof. Bianchi', subject: 'Matematica' },
    { id: 2, name: 'Prof.ssa Verdi', subject: 'Italiano' }
]

const availableSlots = [
    { id: 1, label: 'Lun 25/01 - 15:00' },
    { id: 2, label: 'Lun 25/01 - 15:15' },
    { id: 3, label: 'Mer 27/01 - 10:00' }
]

const newBooking = ref({ teacher: null, slot: null })

const confirmBooking = () => {
    bookings.value.push({
        id: Date.now(),
        teacher: newBooking.value.teacher.name,
        subject: newBooking.value.teacher.subject,
        date: newBooking.value.slot.label.split(' - ')[0],
        time: newBooking.value.slot.label.split(' - ')[1],
        isOnline: true,
        link: '#'
    })
    showBookingDialog.value = false
    step.value = 1
    newBooking.value = { teacher: null, slot: null }
    $q.notify({ type: 'positive', message: 'Prenotazione confermata!' })
}

const cancelBooking = (b) => {
    bookings.value = bookings.value.filter(x => x.id !== b.id)
    $q.notify({ type: 'warning', message: 'Prenotazione cancellata' })
}
</script>
