<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">Incontri Scuola-Famiglia</div>
       <q-btn color="primary" icon="add" label="Nuova Disponibilità" @click="showSlotDialog = true" />
    </div>

    <div class="row q-col-gutter-lg">
        <!-- Calendar/Slots View -->
        <div class="col-12 col-md-8">
            <q-card>
                <q-tabs v-model="tab" class="text-primary" align="left">
                    <q-tab name="upcoming" label="Prossimi Incontri" />
                    <q-tab name="slots" label="Le Mie Disponibilità" />
                </q-tabs>
                <q-separator />
                
                <q-tab-panels v-model="tab" animated>
                    <q-tab-panel name="upcoming">
                         <q-list separator>
                             <q-item v-for="meeting in meetings" :key="meeting.id">
                                 <q-item-section avatar>
                                     <q-avatar color="primary" text-color="white" icon="event" />
                                 </q-item-section>
                                 <q-item-section>
                                     <q-item-label class="text-weight-bold">{{ meeting.parentName }} ({{ meeting.studentName }})</q-item-label>
                                     <q-item-label caption>{{ meeting.date }} - {{ meeting.time }}</q-item-label>
                                 </q-item-section>
                                 <q-item-section side>
                                     <div class="row q-gutter-xs">
                                         <q-btn flat round color="green" icon="videocam" v-if="meeting.isOnline" />
                                         <q-btn flat round color="orange" icon="edit" />
                                         <q-btn flat round color="negative" icon="cancel" />
                                     </div>
                                 </q-item-section>
                             </q-item>
                             <q-item v-if="meetings.length === 0">
                                 <q-item-section class="text-center text-grey">Nessun incontro programmato</q-item-section>
                             </q-item>
                         </q-list>
                    </q-tab-panel>

                    <q-tab-panel name="slots">
                         <div class="text-subtitle2 q-mb-sm">Disponibilità configurate per i genitori</div>
                         <q-list bordered separator>
                             <q-item v-for="slot in slots" :key="slot.id">
                                 <q-item-section>
                                     <q-item-label>{{ slot.day }} - {{ slot.startTime }} / {{ slot.endTime }}</q-item-label>
                                     <q-item-label caption>Durata slot: {{ slot.duration }} min - Max {{ slot.maxBookings }} prenotazioni</q-item-label>
                                 </q-item-section>
                                 <q-item-section side>
                                     <q-toggle v-model="slot.active" color="green" />
                                 </q-item-section>
                             </q-item>
                         </q-list>
                    </q-tab-panel>
                </q-tab-panels>
            </q-card>
        </div>

        <!-- Settings sidebar -->
        <div class="col-12 col-md-4">
            <q-card class="bg-blue-1">
                <q-card-section>
                    <div class="text-h6">Impostazioni Colloqui</div>
                    <div class="q-mt-sm">
                        <q-toggle v-model="settings.onlineEnabled" label="Abilita Colloqui Online (Meet/Zoom)" />
                        <q-input v-if="settings.onlineEnabled" v-model="settings.meetLink" label="Link Riunione Predefinito" dense outlined class="q-mt-sm bg-white" />
                    </div>
                </q-card-section>
            </q-card>
            
            <q-card class="q-mt-md">
                <q-card-section>
                    <div class="text-subtitle1">Richieste in Attesa</div>
                    <q-list class="q-mt-sm">
                        <q-item class="bg-orange-1 rounded-borders q-mb-xs">
                            <q-item-section>
                                <q-item-label>Rossi Luigi</q-item-label>
                                <q-item-label caption>Richiesta urgente - 15 Gen</q-item-label>
                            </q-item-section>
                            <q-item-section side>
                                <q-btn size="sm" color="primary" label="Accetta" />
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-card-section>
            </q-card>
        </div>
    </div>

    <!-- Create Slot Dialog -->
    <q-dialog v-model="showSlotDialog">
        <q-card style="min-width: 400px">
            <q-card-section class="text-h6">Nuova Disponibilità</q-card-section>
            <q-card-section class="q-gutter-md">
                 <q-date v-model="newSlot.date" mask="YYYY-MM-DD" minimal />
                 <div class="row q-gutter-sm">
                     <q-input v-model="newSlot.start" mask="time" :rules="['time']" label="Inizio" outlined dense class="col" />
                     <q-input v-model="newSlot.end" mask="time" :rules="['time']" label="Fine" outlined dense class="col" />
                 </div>
                 <q-input v-model.number="newSlot.duration" type="number" label="Durata colloquio (min)" outlined dense />
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Crea Slot" v-close-popup />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, reactive } from 'vue'

const tab = ref('upcoming')
const showSlotDialog = ref(false)

const meetings = ref([
    { id: 1, parentName: 'Bianchi Giovanni', studentName: 'Bianchi Anna', date: '2025-01-20', time: '15:00', isOnline: true },
    { id: 2, parentName: 'Rossi Maria', studentName: 'Rossi Mario', date: '2025-01-20', time: '15:15', isOnline: false }
])

const slots = ref([
    { id: 1, day: 'Lunedi', startTime: '15:00', endTime: '17:00', duration: 15, maxBookings: 8, active: true },
    { id: 2, day: 'Mercoledi', startTime: '10:00', endTime: '11:00', duration: 15, maxBookings: 4, active: true }
])

const settings = reactive({
    onlineEnabled: true,
    meetLink: 'https://meet.google.com/abc-defg-hij'
})

const newSlot = reactive({ date: '', start: '15:00', end: '17:00', duration: 15 })
</script>
