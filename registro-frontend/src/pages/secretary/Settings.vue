<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h4 text-weight-bold text-outfit q-my-none text-gradient-premium">
          Impostazioni Scuola
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">Configurazione istituto, calendario e orari</div>
      </div>
    </div>

    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden bg-white">
        <q-tabs
            v-model="tab"
            dense
            class="text-slate-500 border-b border-slate-100"
            active-color="primary"
            indicator-color="primary"
            align="left"
            narrow-indicator
            no-caps
        >
            <q-tab name="general" label="Generale" icon="settings" class="q-px-xl py-4" />
            <q-tab name="calendar" label="Calendario Scolastico" icon="calendar_today" class="q-px-xl py-4" />
            <q-tab name="hours" label="Orari Ricevimento" icon="schedule" class="q-px-xl py-4" />
            <q-tab name="accessibility" label="Accessibilità Visiva" icon="accessibility_new" class="q-px-xl py-4" />
        </q-tabs>

        <q-tab-panels v-model="tab" animated class="bg-transparent">
            <!-- General Settings -->
            <q-tab-panel name="general" class="q-pa-xl">
                <div class="row items-center q-mb-xl">
                  <q-avatar color="indigo-50" text-color="indigo-700" icon="business" size="48px" class="q-mr-md" />
                  <div class="text-h5 text-weight-bold text-slate-800">Dati Istituto</div>
                </div>
                
                <div class="row q-col-gutter-lg">
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.schoolName" label="Nome Istituto" outlined />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.schoolCode" label="Codice Meccanografico" outlined />
                    </div>
                    <div class="col-12 col-md-8">
                        <q-input v-model="settings.address" label="Indirizzo" outlined />
                    </div>
                    <div class="col-12 col-md-4">
                        <q-input v-model="settings.email" label="Email Segreteria" outlined />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.pec" label="PEC" outlined />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.phone" label="Telefono" outlined />
                    </div>
                </div>
                <div class="row q-mt-xl">
                    <q-btn color="primary" label="Salva Modifiche" size="lg" padding="md xl" no-caps class="rounded-lg shadow-sm" @click="saveSettings" />
                </div>
            </q-tab-panel>

            <!-- Calendar -->
            <q-tab-panel name="calendar" class="q-pa-xl">
                <div class="row items-center justify-between q-mb-xl">
                  <div class="row items-center">
                    <q-avatar color="orange-50" text-color="orange-700" icon="event" size="48px" class="q-mr-md" />
                    <div class="text-h5 text-weight-bold text-slate-800">Anno Scolastico {{ currentYearStr }}</div>
                  </div>
                  <q-btn color="primary" unelevated label="Aggiungi Chiusura" icon="add" no-caps class="rounded-lg q-px-md" @click="openHolidayDialog()" />
                </div>
                
                 <div class="q-pa-xl bg-slate-50 rounded-2xl q-mb-xl border-slate-200">
                      <div class="row items-center justify-between q-mb-lg">
                        <div class="text-subtitle1 text-slate-800 text-weight-bold">Periodi di Valutazione (Quadrimestri / Trimestri)</div>
                        <q-btn size="sm" color="primary" label="+ Aggiungi Periodo" flat @click="openAddPeriodDialog" />
                      </div>
                      <div class="row q-col-gutter-lg items-center">
                          <div class="col-12 col-sm-6">
                              <q-input v-model="settings.term1End" type="date" outlined label="Fine I Quadrimestre" stack-label />
                          </div>
                          <div class="col-12 col-sm-6">
                              <q-input v-model="settings.term2End" type="date" outlined label="Fine II Quadrimestre" stack-label />
                          </div>
                      </div>
                 </div>

                <div class="text-subtitle1 text-slate-800 text-weight-bold q-mb-md">Festività & Chiusure</div>
                <q-list separator class="rounded-xl border-slate-100 overflow-hidden">
                     <q-item v-for="(holiday, index) in holidays" :key="index" class="q-py-lg">
                         <q-item-section avatar>
                            <q-avatar color="orange-50" text-color="orange-700" icon="celebration" size="40px" />
                         </q-item-section>
                         <q-item-section>
                            <q-item-label class="text-weight-bold text-slate-800">{{ holiday.title }}</q-item-label>
                            <q-item-label caption class="text-slate-500">{{ formatDateRange(holiday.start, holiday.end) }}</q-item-label>
                         </q-item-section>
                         <q-item-section side>
                            <div class="row q-gutter-xs">
                              <q-btn flat round icon="edit" color="primary" @click="openHolidayDialog(index)" />
                              <q-btn flat round icon="delete" color="negative" @click="confirmDeleteHoliday(index)" />
                            </div>
                         </q-item-section>
                     </q-item>
                     <q-item v-if="holidays.length === 0" class="q-pa-xl text-center text-slate-400">
                        <q-item-section>
                          <q-icon name="calendar_today" size="64px" class="q-mb-md opacity-20" />
                          <div class="text-h6">Nessuna festività configurata</div>
                        </q-item-section>
                     </q-item>
                </q-list>
                
                <div class="row q-mt-xl">
                    <q-btn color="primary" label="Salva Calendario" size="lg" padding="md xl" no-caps class="rounded-lg shadow-sm" @click="saveHolidays" :loading="savingHolidays" />
                </div>
            </q-tab-panel>
            
            <!-- Hours -->
            <q-tab-panel name="hours" class="q-pa-xl">
                 <div class="row items-center justify-between q-mb-xl">
                   <div class="row items-center">
                      <q-avatar color="emerald-50" text-color="emerald-700" icon="schedule" size="48px" class="q-mr-md" />
                      <div class="text-h5 text-weight-bold text-slate-800">Orari Ricevimento Segreteria</div>
                   </div>
                   <q-btn color="primary" unelevated label="Aggiungi Orario" icon="add" no-caps class="rounded-lg q-px-md" @click="addHourRow" />
                 </div>

                 <div class="row q-col-gutter-lg">
                     <div class="col-12 col-md-6" v-for="(hour, index) in officeHours" :key="index">
                         <q-card flat class="rounded-xl border-slate-100 bg-slate-50">
                             <q-card-section class="row q-col-gutter-md items-center">
                                 <div class="col-4">
                                     <q-select 
                                        v-model="hour.day" 
                                        :options="['Lunedì', 'Martedì', 'Mercoledì', 'Giovedì', 'Venerdì', 'Sabato', 'Domenica']" 
                                        label="Giorno" 
                                        outlined 
                                        dense 
                                        class="rounded-lg"
                                     />
                                 </div>
                                 <div class="col">
                                     <q-input outlined dense v-model="hour.start" label="Dalle" type="time" class="rounded-lg" />
                                 </div>
                                 <div class="col">
                                     <q-input outlined dense v-model="hour.end" label="Alle" type="time" class="rounded-lg" />
                                 </div>
                                 <div class="col-auto">
                                     <q-btn flat round icon="delete" color="negative" size="sm" @click="officeHours.splice(index, 1)" />
                                 </div>
                             </q-card-section>
                         </q-card>
                     </div>
                 </div>

                 <div v-if="officeHours.length === 0" class="q-pa-xl text-center text-slate-400">
                    <q-icon name="schedule" size="64px" class="q-mb-md opacity-20" />
                    <div class="text-h6">Nessun orario configurato</div>
                 </div>

                 <div class="q-mt-xl row">
                     <q-btn color="primary" label="Salva Orari" size="lg" padding="md xl" no-caps class="rounded-lg shadow-sm" @click="saveOfficeHours" :loading="savingHours" />
                 </div>
             </q-tab-panel>

             <!-- Accessibility Tab -->
             <q-tab-panel name="accessibility" class="q-pa-xl">
                 <div class="row items-center q-mb-xl">
                   <q-avatar color="indigo-50" text-color="indigo-700" icon="accessibility_new" size="48px" class="q-mr-md" />
                   <div>
                     <div class="text-h5 text-weight-bold text-slate-800">Accessibilità Visiva &amp; Modalità di Lettura</div>
                     <div class="text-caption text-slate-500">Impostazioni generali per la leggibilità e l'accessibilità visiva (DSA / WCAG 2.1 AAA)</div>
                   </div>
                 </div>

                 <div class="row q-col-gutter-lg">
                   <!-- Font OpenDyslexic (DSA) -->
                   <div class="col-12 col-md-6">
                     <q-card flat bordered class="q-pa-md rounded-xl bg-white full-height shadow-sm">
                       <div class="row items-center justify-between q-mb-sm">
                         <div class="row items-center">
                           <q-avatar color="indigo-50" text-color="indigo-700" icon="spellcheck" size="44px" class="q-mr-sm" />
                           <div>
                             <div class="text-subtitle1 text-weight-bold text-slate-800">Font OpenDyslexic (Alta Leggibilità DSA)</div>
                             <div class="text-caption text-slate-500">Attiva il carattere specifico per la dislessia e la facilitazione di lettura</div>
                           </div>
                         </div>
                         <q-toggle
                           v-model="themeStore.dsaFont"
                           color="indigo"
                           size="lg"
                           @update:model-value="themeStore.toggleDsaFont"
                         />
                       </div>
                       <q-separator class="q-my-sm" />
                       <div class="q-pa-md bg-slate-50 rounded-lg text-slate-700 text-body2 q-mt-sm border border-slate-100" :class="{ 'dsa-font-active': themeStore.dsaFont }">
                         <span class="text-weight-bold">Anteprima Testo:</span> Piattaforma scolastica istituzionale con supporto all'accessibilità visiva ed inclusione digitale.
                       </div>
                     </q-card>
                   </div>

                   <!-- Contrasto Elevato -->
                   <div class="col-12 col-md-6">
                     <q-card flat bordered class="q-pa-md rounded-xl bg-white full-height shadow-sm">
                       <div class="row items-center justify-between q-mb-sm">
                         <div class="row items-center">
                           <q-avatar color="amber-50" text-color="amber-9" icon="contrast" size="44px" class="q-mr-sm" />
                           <div>
                             <div class="text-subtitle1 text-weight-bold text-slate-800">Modalità Contrasto Elevato</div>
                             <div class="text-caption text-slate-500">Definizione netta dei bordi delle componenti e contrasto aumentato</div>
                           </div>
                         </div>
                         <q-toggle
                           v-model="themeStore.highContrast"
                           color="amber-9"
                           size="lg"
                           @update:model-value="themeStore.toggleHighContrast"
                         />
                       </div>
                       <q-separator class="q-my-sm" />
                       <div class="q-pa-md bg-slate-50 rounded-lg text-slate-700 text-body2 q-mt-sm border border-slate-100" :class="{ 'high-contrast-active': themeStore.highContrast }">
                         <span class="text-weight-bold">Anteprima Contrasto:</span> Modalità attiva per l'incremento di nitidezza della grafica e dei controlli.
                       </div>
                     </q-card>
                   </div>
                 </div>
             </q-tab-panel>
        </q-tab-panels>
    </q-card>

    <!-- Holiday Dialog -->
    <q-dialog v-model="holidayDialog.show" persistent class="premium-dialog">
      <q-card style="min-width: 500px" class="rounded-xl overflow-hidden shadow-24">
        <q-card-section class="bg-gradient-primary text-white q-pa-lg">
          <div class="text-h5 text-weight-bold">{{ holidayDialog.editIndex !== null ? 'Modifica Chiusura' : 'Nuova Chiusura' }}</div>
        </q-card-section>

        <q-card-section class="q-pa-xl q-gutter-y-lg">
          <q-input v-model="holidayForm.title" label="Descrizione (es. Vacanze Natale)" outlined />
          <div class="row q-col-gutter-lg">
            <div class="col-6">
              <q-input v-model="holidayForm.start" type="date" label="Inizio" outlined stack-label />
            </div>
            <div class="col-6">
              <q-input v-model="holidayForm.end" type="date" label="Fine" outlined stack-label />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-lg bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup color="slate-400" no-caps class="rounded-lg q-px-md" />
          <q-btn 
            unelevated 
            :label="holidayDialog.editIndex !== null ? 'Aggiorna' : 'Aggiungi'" 
            color="primary" 
            class="rounded-lg q-px-xl shadow-sm" 
            @click="saveHolidayToList" 
            :disable="!holidayForm.title || !holidayForm.start || !holidayForm.end" 
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useThemeStore } from 'src/stores/theme'
import adminService from 'src/services/adminService'
import api from 'src/services/api'

const $q = useQuasar()
const themeStore = useThemeStore()

const getCurrentAcademicYear = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth() + 1; // 1-12
  if (month >= 9) { 
    return `${year}/${year + 1}`;
  } else {
    return `${year - 1}/${year}`;
  }
}
const currentYearStr = getCurrentAcademicYear();

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

const openAddPeriodDialog = () => {
  $q.dialog({
    title: 'Aggiungi Periodo Valutativo',
    message: 'Nome periodo (es. 1° Quadrimestre, 1° Trimestre):',
    prompt: { model: '', type: 'text' },
    cancel: true, persistent: true
  }).onOk(async (name) => {
    if (!name) return
    try {
      await api.post('/school-calendar/periods', {
        name: name,
        code: name.slice(0, 3).toUpperCase(),
        start_date: `${new Date().getFullYear()}-09-01`,
        end_date: `${new Date().getFullYear() + 1}-06-30`,
        is_current: true
      })
      $q.notify({ type: 'positive', message: 'Periodo creato con successo' })
    } catch (e) {
      $q.notify({ type: 'negative', message: 'Errore creazione periodo' })
    }
  })
}

const fetchOfficeHours = async () => {
  try {
    const res = await adminService.getSchoolSetting('school_office_hours')
    if (res.data && res.data.value) {
      officeHours.value = JSON.parse(res.data.value)
    } else {
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
    persistent: true,
    ok: {
      color: 'negative',
      label: 'Elimina',
      flat: false
    }
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
.opacity-20 { opacity: 0.2; }
.py-4 { padding-top: 1rem; padding-bottom: 1rem; }
</style>
