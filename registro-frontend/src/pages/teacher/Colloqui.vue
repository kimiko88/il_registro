<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 font-bold text-slate-800 q-my-none">Gestione Colloqui</h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Pianifica le tue disponibilità ed accetta le prenotazioni dei genitori
        </p>
      </div>
      <q-btn color="indigo" icon="add" label="Nuova Disponibilità" rounded @click="showSlotDialog = true" />
    </div>

    <!-- Main Content -->
    <div class="row q-col-gutter-lg">
        <div class="col-12 col-md-8">
            <q-card flat bordered class="rounded-xl shadow-sm overflow-hidden">
                <q-tabs v-model="tab" class="text-indigo bg-indigo-50/50" active-color="indigo" indicator-color="indigo" align="left">
                    <q-tab name="meetings" label="Incontri Programmati" icon="event" />
                    <q-tab name="slots" label="Le tue Disponibilità" icon="schedule" />
                </q-tabs>

                <q-separator />

                <q-tab-panels v-model="tab" animated>
                    <q-tab-panel name="meetings" class="q-pa-none">
                         <q-list separator>
                             <q-item v-for="meeting in meetings" :key="meeting.id" class="q-py-md">
                                 <q-item-section avatar>
                                     <q-avatar color="indigo-1" text-color="indigo" icon="person" />
                                 </q-item-section>
                                 <q-item-section>
                                     <q-item-label class="text-weight-bold text-slate-700">
                                       {{ meeting.parent_name || 'Genitore' }} 
                                       <span class="text-weight-light text-slate-500">per</span> 
                                       {{ meeting.student_name || 'Studente' }}
                                     </q-item-label>
                                     <q-item-label v-if="meeting.slot_info" caption class="text-slate-500">
                                       {{ formatDate(meeting.slot_info.date) }} @ {{ meeting.slot_info.time_range }}
                                     </q-item-label>
                                 </q-item-section>
                                 <q-item-section side>
                                     <div class="row items-center q-gutter-xs">
                                         <q-btn 
                                           v-if="meeting.status === 'Pending'" 
                                           dense 
                                           unelevated
                                           color="positive" 
                                           icon="check" 
                                           label="Conferma" 
                                           no-caps 
                                           class="q-mr-xs rounded-pill"
                                           @click="confirmMeeting(meeting)" 
                                         />
                                         <q-btn 
                                           v-if="meeting.meet_link || settings.meetLink" 
                                           dense 
                                           unelevated
                                           color="primary" 
                                           icon="videocam" 
                                           label="Entra nella riunione" 
                                           no-caps 
                                           class="q-mr-xs rounded-pill"
                                           :href="meeting.meet_link || settings.meetLink"
                                           target="_blank"
                                         />
                                         <q-badge :color="getStatusColor(meeting.status)" rounded class="q-mr-sm">
                                           {{ formatStatusLabel(meeting.status) }}
                                         </q-badge>
                                         <q-btn flat round dense color="negative" icon="cancel" @click="confirmCancelBooking(meeting)" />
                                     </div>
                                 </q-item-section>
                             </q-item>
                             <q-item v-if="meetings.length === 0" class="q-py-xl">
                                 <q-item-section class="text-center text-slate-400">
                                   <q-icon name="calendar_today" size="48px" class="q-mb-sm" />
                                   <div>Nessun incontro programmato</div>
                                 </q-item-section>
                             </q-item>
                         </q-list>
                    </q-tab-panel>

                    <q-tab-panel name="slots" class="q-pa-none">
                         <q-list separator>
                             <q-item v-for="slot in slots" :key="slot.id" class="q-py-md">
                                 <q-item-section>
                                     <q-item-label class="text-weight-medium text-slate-700">
                                       {{ formatDate(slot.date) }} - {{ slot.time_range }}
                                     </q-item-label>
                                     <q-item-label caption>
                                       Tipo: {{ getSlotTypeLabel(slot.type) }} • Prenotazioni: {{ slot.current_bookings || (slot.available ? 0 : 1) }}/{{ slot.max_bookings || 1 }}
                                     </q-item-label>
                                 </q-item-section>
                                 <q-item-section side>
                                     <div class="row items-center q-gutter-sm">
                                       <q-btn 
                                         flat 
                                         round 
                                         dense 
                                         :color="slot.available ? 'negative' : 'grey'" 
                                         icon="delete" 
                                         @click="handleDeleteSlot(slot)" 
                                       >
                                         <q-tooltip>Elimina o Annulla Slot</q-tooltip>
                                       </q-btn>
                                     </div>
                                 </q-item-section>
                             </q-item>
                             <q-item v-if="slots.length === 0" class="q-py-xl">
                                 <q-item-section class="text-center text-slate-400">
                                   <q-icon name="history" size="48px" class="q-mb-sm" />
                                   <div>Nessuna disponibilità configurata</div>
                                 </q-item-section>
                             </q-item>
                         </q-list>
                    </q-tab-panel>
                </q-tab-panels>
            </q-card>
        </div>

        <!-- Settings sidebar -->
        <div class="col-12 col-md-4">
            <q-card flat bordered class="rounded-xl shadow-sm bg-indigo-50/50">
                <q-card-section>
                    <div class="text-h6 text-indigo-900 font-bold">Impostazioni Colloqui</div>
                    <div class="q-mt-sm">
                        <q-toggle v-model="settings.onlineEnabled" label="Abilita Colloqui Online (Meet/Zoom)" color="indigo" />
                        <q-input 
                          v-if="settings.onlineEnabled" 
                          v-model="settings.meetLink" 
                          label="Link Riunione Predefinito" 
                          dense 
                          outlined 
                          rounded
                          class="q-mt-sm bg-white" 
                        />
                        <q-btn color="indigo" label="Salva Impostazioni" dense icon="save" class="q-mt-md full-width" @click="saveSettings" />
                    </div>
                </q-card-section>
            </q-card>
            
            <q-card flat bordered class="q-mt-md rounded-xl shadow-sm overflow-hidden">
                <q-card-section class="bg-amber-50">
                    <div class="text-subtitle1 text-amber-900 font-bold">Informazioni Utili</div>
                    <div class="text-caption text-amber-800 q-mt-xs">
                      Gli slot senza prenotazioni vengono eliminati permanentemente. 
                      Quelli già prenotati verranno contrassegnati come "Annullati".
                    </div>
                </q-card-section>
            </q-card>
        </div>
    </div>

    <!-- Create Slot Dialog -->
    <q-dialog v-model="showSlotDialog" persistent>
        <q-card style="min-width: 450px" class="rounded-xl">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Crea Disponibilità</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section class="q-pa-md">
                 <div class="text-subtitle2 q-mb-xs">Seleziona uno o più giorni:</div>
                 <q-date 
                   v-model="newSlot.dates" 
                   mask="YYYY-MM-DD" 
                   multiple 
                   flat 
                   bordered
                   class="full-width q-mb-md" 
                   minimal
                 />
                 
                 <div class="row q-col-gutter-md">
                     <div class="col-6">
                       <q-input v-model="newSlot.start" mask="time" label="Dalle ore" outlined dense rounded>
                         <template v-slot:append>
                           <q-icon name="access_time" class="cursor-pointer">
                             <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                               <q-time v-model="newSlot.start">
                                 <div class="row items-center justify-end">
                                   <q-btn v-close-popup label="Chiudi" color="primary" flat />
                                 </div>
                               </q-time>
                             </q-popup-proxy>
                           </q-icon>
                         </template>
                       </q-input>
                     </div>
                     <div class="col-6">
                       <q-input v-model="newSlot.end" mask="time" label="Alle ore" outlined dense rounded>
                         <template v-slot:append>
                           <q-icon name="access_time" class="cursor-pointer">
                             <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                               <q-time v-model="newSlot.end">
                                 <div class="row items-center justify-end">
                                   <q-btn v-close-popup label="Chiudi" color="primary" flat />
                                 </div>
                               </q-time>
                             </q-popup-proxy>
                           </q-icon>
                         </template>
                       </q-input>
                     </div>
                 </div>

                 <div class="row q-col-gutter-md q-mt-xs">
                   <div class="col-6">
                      <q-input 
                        v-model.number="newSlot.duration" 
                        type="number" 
                        label="Durata slot (min)" 
                        outlined 
                        dense 
                        rounded
                        hint="0 per slot unico"
                      />
                   </div>
                   <div class="col-6">
                      <q-select 
                        v-model="newSlot.type" 
                        :options="[
                          { label: 'Individuale', value: 'Individual' },
                          { label: 'Generale', value: 'General' },
                          { label: 'Assemblea', value: 'Assembly' }
                        ]" 
                        label="Tipo" 
                        outlined 
                        dense 
                        rounded
                        emit-value
                        map-options
                      />
                   </div>
                 </div>

                 <q-input v-model="newSlot.location" label="Aula / Link Online" outlined dense rounded class="q-mt-md" />

                 <div class="q-mt-md bg-slate-50 q-pa-sm rounded-borders">
                    <div class="row items-center justify-between">
                       <div class="text-subtitle2">Ripeti settimanalmente</div>
                       <q-toggle v-model="newSlot.isRecurring" color="primary" />
                    </div>
                    <q-slide-transition>
                      <div v-if="newSlot.isRecurring" class="q-mt-sm">
                        <q-input v-model="newSlot.recurringUntil" label="Fino al" outlined dense rounded mask="date">
                          <template v-slot:append>
                            <q-icon name="event" class="cursor-pointer">
                              <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                                <q-date v-model="newSlot.recurringUntil">
                                  <div class="row items-center justify-end">
                                    <q-btn v-close-popup label="Chiudi" color="primary" flat />
                                  </div>
                                </q-date>
                              </q-popup-proxy>
                            </q-icon>
                          </template>
                        </q-input>
                      </div>
                    </q-slide-transition>
                 </div>
            </q-card-section>

            <q-card-actions align="right" class="q-pa-md">
                <q-btn flat label="Annulla" v-close-popup rounded />
                <q-btn color="indigo" label="Genera Disponibilità" rounded @click="saveSlots" :loading="savingSlots" />
            </q-card-actions>
        </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import api from 'src/services/api'

const $q = useQuasar()

const tab = ref('meetings')
const showSlotDialog = ref(false)
const savingSlots = ref(false)

const settings = reactive({
    onlineEnabled: false,
    meetLink: ''
})

const newSlot = reactive({
    dates: [],
    start: '09:00',
    end: '12:00',
    duration: 15,
    type: 'Individual',
    location: '',
    isRecurring: false,
    recurringUntil: ''
})

const meetings = ref([])
const slots = ref([])

const getSlotTypeLabel = (type) => {
  switch (type) {
    case 'Individual': return 'Individuale'
    case 'General': return 'Generale'
    case 'Assembly': return 'Assemblea'
    default: return type || 'Individuale'
  }
}

const loadData = async () => {
    try {
        const [slotsRes, meetingsRes] = await Promise.all([
            api.get('/colloqui/slots/my'),
            api.get('/colloqui/bookings')
        ])
        slots.value = slotsRes.data || []
        meetings.value = meetingsRes.data || []
    } catch (err) {
        console.error('Failed to load colloquio data:', err)
        $q.notify({ color: 'negative', message: 'Errore nel caricamento dei dati' })
    }
}

const confirmMeeting = async (meeting) => {
  try {
    await api.patch(`/colloqui/bookings/${meeting.id}/confirm`)
    $q.notify({ color: 'positive', message: 'Incontro confermato' })
    loadData()
  } catch (err) {
    $q.notify({ color: 'negative', message: 'Errore durante la conferma' })
  }
}

const saveSlots = async () => {
    if (!newSlot.dates || newSlot.dates.length === 0) {
        $q.notify({ color: 'warning', message: 'Seleziona almeno un giorno' })
        return
    }
    
    savingSlots.value = true
    try {
        const payload = {
            dates: newSlot.dates,
            start_time: newSlot.start,
            end_time: newSlot.end,
            duration: newSlot.duration,
            type: newSlot.type,
            location: newSlot.location,
            is_recurring: newSlot.isRecurring,
            recurring_until: newSlot.isRecurring ? newSlot.recurringUntil : null
        }

        await api.post('/colloqui/slots', payload)
        $q.notify({ color: 'positive', message: 'Disponibilità generate con successo' })
        showSlotDialog.value = false
        
        // Reset form
        newSlot.dates = []
        newSlot.start = '09:00'
        newSlot.end = '12:00'
        newSlot.duration = 15
        newSlot.type = 'Individual'
        newSlot.location = ''
        newSlot.isRecurring = false
        newSlot.recurringUntil = ''

        loadData()
    } catch (err) {
        $q.notify({ color: 'negative', message: 'Errore durante la generazione delle disponibilità' })
    } finally {
        savingSlots.value = false
    }
}

const handleDeleteSlot = async (slot) => {
    $q.dialog({
        title: 'Elimina Disponibilità',
        message: slot.available 
            ? 'Vuoi eliminare questa disponibilità?' 
            : 'Questa disponibilità ha già una prenotazione. Verrà annullata e il genitore notificato. Continuare?',
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            await api.delete(`/colloqui/slots/${slot.id}`)
            $q.notify({ color: 'positive', message: 'Disponibilità eliminata' })
            loadData()
        } catch (err) {
            $q.notify({ color: 'negative', message: 'Errore durante l\'eliminazione' })
        }
    })
}

const confirmCancelBooking = (meeting) => {
    $q.dialog({
        title: 'Annulla Incontro',
        message: 'Vuoi annullare questo incontro con il genitore?',
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            await api.patch(`/colloqui/bookings/${meeting.id}/cancel`)
            $q.notify({ color: 'positive', message: 'Incontro annullato' })
            loadData()
        } catch (err) {
            $q.notify({ color: 'negative', message: 'Errore durante l\'annullamento' })
        }
    })
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('it-IT', { weekday: 'short', day: 'numeric', month: 'short' })
  } catch {
    return dateStr
  }
}

const formatStatusLabel = (status) => {
  switch (status) {
    case 'Confirmed': return 'Confermato'
    case 'Cancelled': return 'Annullato'
    case 'Completed': return 'Completato'
    case 'Pending': return 'In Attesa'
    case 'Booked': return 'Prenotato'
    default: return status || 'N/D'
  }
}

const getStatusColor = (status) => {
  switch (status) {
    case 'Confirmed': return 'green'
    case 'Pending': return 'amber'
    case 'Cancelled': return 'red'
    case 'Completed': return 'blue'
    default: return 'grey'
  }
}

const saveSettings = () => {
  try {
    localStorage.setItem('teacher_colloqui_settings', JSON.stringify(settings))
    $q.notify({ color: 'positive', message: 'Impostazioni salvate con successo' })
  } catch (e) {
    $q.notify({ color: 'negative', message: 'Errore durante il salvataggio' })
  }
}

onMounted(() => {
  const saved = localStorage.getItem('teacher_colloqui_settings')
  if (saved) {
    try {
      const parsed = JSON.parse(saved)
      settings.onlineEnabled = !!parsed.onlineEnabled
      settings.meetLink = parsed.meetLink || ''
    } catch (e) {
      console.warn('Failed to parse settings:', e)
    }
  }
  loadData()
})
</script>

<style scoped>
.rounded-xl {
  border-radius: 1rem;
}
.rounded-pill {
  border-radius: 9999px;
}
.shadow-sm {
  box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
}
</style>
