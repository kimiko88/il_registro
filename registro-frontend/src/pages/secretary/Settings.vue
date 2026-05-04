<template>
  <q-page padding class="bg-grey-1">
    <div class="row items-center q-mb-lg">
      <div class="col">
        <h1 class="text-h4 text-weight-bold text-outfit q-my-none text-gradient-premium">
          Impostazioni Scuola
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">Configurazione istituto e calendario scolastico</div>
      </div>
    </div>

    <q-card class="shadow-1 rounded-lg overflow-hidden">
        <q-tabs
            v-model="tab"
            dense
            class="bg-white text-grey-7"
            active-color="primary"
            indicator-color="primary"
            align="left"
            narrow-indicator
        >
            <q-tab name="general" label="Generale" icon="settings" class="q-px-lg" />
            <q-tab name="calendar" label="Calendario Scolastico" icon="calendar_today" class="q-px-lg" />
            <q-tab name="hours" label="Orari Uffici" icon="schedule" class="q-px-lg" />
        </q-tabs>

        <q-separator />

        <q-tab-panels v-model="tab" animated class="bg-white">
            <!-- General Settings -->
            <q-tab-panel name="general" class="q-pa-xl">
                <div class="text-h6 text-weight-bold q-mb-xl row items-center">
                  <q-icon name="business" color="primary" size="md" class="q-mr-sm" />
                  Dati Istituto
                </div>
                <div class="row q-col-gutter-lg">
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.schoolName" label="Nome Istituto" outlined bg-color="white" />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.schoolCode" label="Codice Meccanografico" outlined bg-color="white" />
                    </div>
                    <div class="col-12 col-md-8">
                        <q-input v-model="settings.address" label="Indirizzo" outlined bg-color="white" />
                    </div>
                    <div class="col-12 col-md-4">
                        <q-input v-model="settings.email" label="Email Segreteria" outlined bg-color="white" />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.pec" label="PEC" outlined bg-color="white" />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.phone" label="Telefono" outlined bg-color="white" />
                    </div>
                </div>
                <div class="row q-mt-xl">
                    <q-btn color="primary" label="Salva Modifiche" size="lg" padding="md xl" no-caps class="rounded-lg shadow-soft" @click="saveSettings" />
                </div>
            </q-tab-panel>

            <!-- Calendar -->
            <q-tab-panel name="calendar" class="q-pa-xl">
                <div class="row items-center justify-between q-mb-xl">
                  <div class="text-h6 text-weight-bold row items-center">
                    <q-icon name="event" color="primary" size="md" class="q-mr-sm" />
                    Anno Scolastico 2024/2025
                  </div>
                  <q-btn color="primary" outline label="Aggiungi Chiusura" icon="add" no-caps @click="openHolidayDialog()" />
                </div>
                
                <div class="q-pa-lg bg-indigo-50 rounded-lg q-mb-xl border-indigo-100 border-1">
                     <div class="text-subtitle2 text-indigo-9 text-weight-bold q-mb-md">Periodi Valutazione</div>
                     <div class="row q-col-gutter-lg items-center">
                         <div class="col-12 col-sm-6">
                             <q-input v-model="settings.term1End" type="date" outlined label="Fine I Quadrimestre" bg-color="white" stack-label />
                         </div>
                         <div class="col-12 col-sm-6">
                             <q-input v-model="settings.term2End" type="date" outlined label="Fine II Quadrimestre" bg-color="white" stack-label />
                         </div>
                     </div>
                </div>

                <div class="text-subtitle2 text-weight-bold q-mb-md">Festività & Chiusure</div>
                <q-list bordered separator class="rounded-lg bg-white overflow-hidden">
                     <q-item v-for="(holiday, index) in holidays" :key="index" class="q-py-md">
                         <q-item-section avatar>
                            <q-avatar color="orange-1" text-color="orange-9" icon="celebration" size="md" />
                         </q-item-section>
                         <q-item-section>
                            <q-item-label class="text-weight-bold">{{ holiday.title }}</q-item-label>
                            <q-item-label caption>{{ formatDateRange(holiday.start, holiday.end) }}</q-item-label>
                         </q-item-section>
                         <q-item-section side>
                            <div class="row q-gutter-xs">
                              <q-btn flat round icon="edit" color="primary" size="sm" @click="openHolidayDialog(index)" />
                              <q-btn flat round icon="delete" color="negative" size="sm" @click="confirmDeleteHoliday(index)" />
                            </div>
                         </q-item-section>
                     </q-item>
                     <q-item v-if="holidays.length === 0" class="q-pa-xl text-center text-grey">
                        <q-item-section>
                          <q-icon name="calendar_today" size="40px" class="q-mb-sm" />
                          <div>Nessuna festività configurata</div>
                        </q-item-section>
                     </q-item>
                </q-list>
                
                <div class="row q-mt-xl">
                    <q-btn color="primary" label="Salva Calendario" size="lg" padding="md xl" no-caps class="rounded-lg shadow-soft" @click="saveHolidays" :loading="savingHolidays" />
                </div>
            </q-tab-panel>
            
            <!-- Hours -->
            <q-tab-panel name="hours" class="q-pa-xl">
                 <div class="row items-center justify-between q-mb-xl">
                   <div class="text-h6 text-weight-bold row items-center">
                      <q-icon name="schedule" color="primary" size="md" class="q-mr-sm" />
                      Orari Ricevimento Segreteria
                   </div>
                   <q-btn color="primary" outline label="Aggiungi Orario" icon="add" no-caps @click="addHourRow" />
                 </div>

                 <div class="row q-col-gutter-lg q-mt-sm">
                     <div class="col-12 col-md-6" v-for="(hour, index) in officeHours" :key="index">
                         <q-card bordered flat class="rounded-lg bg-grey-1">
                             <q-card-section class="row q-col-gutter-sm items-center">
                                 <div class="col-4">
                                     <q-select 
                                        v-model="hour.day" 
                                        :options="['Lunedì', 'Martedì', 'Mercoledì', 'Giovedì', 'Venerdì', 'Sabato', 'Domenica']" 
                                        label="Giorno" 
                                        outlined 
                                        dense 
                                        bg-color="white" 
                                     />
                                 </div>
                                 <div class="col">
                                     <q-input outlined dense v-model="hour.start" label="Dalle" type="time" bg-color="white" />
                                 </div>
                                 <div class="col">
                                     <q-input outlined dense v-model="hour.end" label="Alle" type="time" bg-color="white" />
                                 </div>
                                 <div class="col-auto">
                                     <q-btn flat round icon="delete" color="negative" size="sm" @click="officeHours.splice(index, 1)" />
                                 </div>
                             </q-card-section>
                         </q-card>
                     </div>
                 </div>

                 <div v-if="officeHours.length === 0" class="q-pa-xl text-center text-grey">
                    <q-icon name="schedule" size="40px" class="q-mb-sm" />
                    <div>Nessun orario configurato</div>
                 </div>

                 <div class="q-mt-xl row">
                     <q-btn color="primary" label="Salva Orari" size="lg" padding="md xl" no-caps class="rounded-lg shadow-soft" @click="saveOfficeHours" :loading="savingHours" />
                 </div>
            </q-tab-panel>
        </q-tab-panels>
    </q-card>

    <!-- Holiday Dialog -->
    <q-dialog v-model="holidayDialog.show" persistent>
      <q-card style="min-width: 400px" class="rounded-lg">
        <q-card-section class="bg-primary text-white q-pa-md">
          <div class="text-h6 text-weight-bold">{{ holidayDialog.editIndex !== null ? 'Modifica Chiusura' : 'Nuova Chiusura' }}</div>
        </q-card-section>

        <q-card-section class="q-pa-lg q-gutter-y-md">
          <q-input v-model="holidayForm.title" label="Descrizione (es. Vacanze Natale)" outlined />
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input v-model="holidayForm.start" type="date" label="Inizio" outlined stack-label />
            </div>
            <div class="col-6">
              <q-input v-model="holidayForm.end" type="date" label="Fine" outlined stack-label />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Annulla" v-close-popup color="grey-7" />
          <q-btn :label="holidayDialog.editIndex !== null ? 'Aggiorna' : 'Aggiungi'" color="primary" @click="saveHolidayToList" :disable="!holidayForm.title || !holidayForm.start || !holidayForm.end" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import adminService from 'src/services/adminService'

const $q = useQuasar()
const tab = ref('general')
const loading = ref(false)
const savingHolidays = ref(false)
const savingHours = ref(false)

const settings = reactive({
    schoolName: 'Istituto Comprensivo "Alessandro Volta"',
    schoolCode: 'RMPC123456',
    address: 'Via Roma 1, 00100 Roma',
    email: 'segreteria@scuola.it',
    pec: 'scuola@pec.it',
    phone: '06 12345678',
    term1End: '2025-01-31',
    term2End: '2025-06-08'
})

const holidays = ref([])
const officeHours = ref([])

const holidayDialog = reactive({
  show: false,
  editIndex: null
})

const holidayForm = reactive({
  title: '',
  start: '',
  end: ''
})

const fetchHolidays = async () => {
  try {
    const res = await adminService.getSchoolSetting('school_calendar_holidays')
    if (res.data && res.data.value) {
      holidays.value = JSON.parse(res.data.value)
    }
  } catch (e) {
    console.error("Error fetching holidays", e)
  }
}

const fetchOfficeHours = async () => {
  try {
    const res = await adminService.getSchoolSetting('school_office_hours')
    if (res.data && res.data.value) {
      officeHours.value = JSON.parse(res.data.value)
    } else {
      // Default fallback
      officeHours.value = [
        { day: 'Lunedì', start: '08:00', end: '14:00' },
        { day: 'Martedì', start: '08:00', end: '14:00' },
        { day: 'Mercoledì', start: '08:00', end: '14:00' },
        { day: 'Giovedì', start: '08:00', end: '14:00' },
        { day: 'Venerdì', start: '08:00', end: '14:00' }
      ]
    }
  } catch (e) {
    console.error("Error fetching office hours", e)
  }
}

const saveHolidays = async () => {
  savingHolidays.value = true
  try {
    await adminService.updateSchoolSetting('school_calendar_holidays', JSON.stringify(holidays.value))
    $q.notify({ type: 'positive', message: 'Calendario salvato con successo' })
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio del calendario' })
  } finally {
    savingHolidays.value = false
  }
}

const saveOfficeHours = async () => {
  savingHours.value = true
  try {
    await adminService.updateSchoolSetting('school_office_hours', JSON.stringify(officeHours.value))
    $q.notify({ type: 'positive', message: 'Orari salvati con successo' })
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio degli orari' })
  } finally {
    savingHours.value = false
  }
}

const addHourRow = () => {
  officeHours.value.push({ day: 'Lunedì', start: '08:00', end: '14:00' })
}

const openHolidayDialog = (index = null) => {
  holidayDialog.editIndex = index
  if (index !== null) {
    Object.assign(holidayForm, holidays.value[index])
  } else {
    holidayForm.title = ''
    holidayForm.start = ''
    holidayForm.end = ''
  }
  holidayDialog.show = true
}

const saveHolidayToList = () => {
  if (holidayDialog.editIndex !== null) {
    holidays.value[holidayDialog.editIndex] = { ...holidayForm }
  } else {
    holidays.value.push({ ...holidayForm })
  }
  holidayDialog.show = false
}

const confirmDeleteHoliday = (index) => {
  $q.dialog({
    title: 'Conferma eliminazione',
    message: 'Sei sicuro di voler rimuovere questa chiusura?',
    cancel: true,
    persistent: true
  }).onOk(() => {
    holidays.value.splice(index, 1)
  })
}

const formatDateRange = (start, end) => {
  const s = new Date(start).toLocaleDateString('it-IT', { day: '2-digit', month: 'short' })
  const e = new Date(end).toLocaleDateString('it-IT', { day: '2-digit', month: 'short' })
  return `${s} - ${e}`
}

const saveSettings = () => {
    $q.loading.show()
    setTimeout(() => {
        $q.loading.hide()
        $q.notify({ type: 'positive', message: 'Impostazioni salvate' })
    }, 800)
}

onMounted(() => {
  fetchHolidays()
  fetchOfficeHours()
})
</script>

<style scoped>
.text-gradient-premium {
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.shadow-soft {
  box-shadow: 0 4px 15px rgba(59, 130, 246, 0.2);
}
.border-indigo-100 {
  border: 1px solid #e0e7ff;
}
</style>
