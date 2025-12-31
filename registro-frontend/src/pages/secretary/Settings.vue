<template>
  <q-page padding>
    <div class="text-h4 text-weight-bold q-mb-md">Impostazioni Scuola</div>

    <q-card>
        <q-tabs
            v-model="tab"
            dense
            class="text-grey"
            active-color="primary"
            indicator-color="primary"
            align="justify"
            narrow-indicator
        >
            <q-tab name="general" label="Generale" />
            <q-tab name="calendar" label="Calendario Scolastico" />
            <q-tab name="hours" label="Orari Uffici" />
        </q-tabs>

        <q-separator />

        <q-tab-panels v-model="tab" animated>
            <!-- General Settings -->
            <q-tab-panel name="general">
                <div class="text-h6 q-mb-md">Dati Istituto</div>
                <div class="row q-col-gutter-md">
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.schoolName" label="Nome Istituto" outlined dense />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.schoolCode" label="Codice Meccanografico" outlined dense />
                    </div>
                    <div class="col-12 col-md-8">
                        <q-input v-model="settings.address" label="Indirizzo" outlined dense />
                    </div>
                    <div class="col-12 col-md-4">
                        <q-input v-model="settings.email" label="Email Segreteria" outlined dense />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.pec" label="PEC" outlined dense />
                    </div>
                    <div class="col-12 col-md-6">
                        <q-input v-model="settings.phone" label="Telefono" outlined dense />
                    </div>
                </div>
                <div class="row q-mt-md">
                    <q-btn color="primary" label="Salva Modifiche" @click="saveSettings" />
                </div>
            </q-tab-panel>

            <!-- Calendar -->
            <q-tab-panel name="calendar">
                <div class="text-h6 q-mb-md">Anno Scolastico 2024/2025</div>
                
                <div class="q-list bordered rounded-borders q-mb-md">
                     <div class="q-pa-sm bg-grey-2 text-weight-bold">Periodi Valutazione</div>
                     <div class="row q-pa-sm items-center q-gutter-md">
                         <div class="col-auto">I Quadrimestre:</div>
                         <div class="col"><q-input v-model="settings.term1End" type="date" dense outlined label="Fine" /></div>
                     </div>
                     <div class="row q-pa-sm items-center q-gutter-md">
                         <div class="col-auto">II Quadrimestre:</div>
                         <div class="col"><q-input v-model="settings.term2End" type="date" dense outlined label="Fine" /></div>
                     </div>
                </div>

                <div class="q-list bordered rounded-borders">
                     <div class="q-pa-sm bg-grey-2 text-weight-bold">Festività & Chiusure</div>
                     <!-- Simplification list -->
                     <q-item>
                         <q-item-section>Vacanze di Natale</q-item-section>
                         <q-item-section side>23 Dic - 7 Gen</q-item-section>
                         <q-item-section side><q-btn flat round icon="edit" size="sm" /></q-item-section>
                     </q-item>
                </div>
            </q-tab-panel>
            
            <!-- Hours -->
            <q-tab-panel name="hours">
                 <div class="text-h6">Orari Ricevimento Segreteria</div>
                 <div class="row q-col-gutter-md q-mt-sm">
                     <div class="col-12 col-md-6" v-for="day in ['Lun', 'Mar', 'Mer', 'Gio', 'Ven']" :key="day">
                         <div class="text-subtitle2">{{ day }}</div>
                         <div class="row q-gutter-sm">
                             <q-input outlined dense v-model="settings.hours[day].start" class="col" type="time" />
                             <q-input outlined dense v-model="settings.hours[day].end" class="col" type="time" />
                         </div>
                     </div>
                 </div>
                 <div class="q-mt-md">
                     <q-btn color="primary" label="Aggiorna Orari" @click="saveSettings" />
                 </div>
            </q-tab-panel>
        </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tab = ref('general')

const settings = reactive({
    schoolName: 'Istituto Comprensivo "Alessandro Volta"',
    schoolCode: 'RMPC123456',
    address: 'Via Roma 1, 00100 Roma',
    email: 'segreteria@scuola.it',
    pec: 'scuola@pec.it',
    phone: '06 12345678',
    term1End: '2025-01-31',
    term2End: '2025-06-08',
    hours: {
        'Lun': { start: '08:00', end: '14:00' },
        'Mar': { start: '08:00', end: '14:00' },
        'Mer': { start: '08:00', end: '14:00' },
        'Gio': { start: '08:00', end: '14:00' },
        'Ven': { start: '08:00', end: '14:00' }
    }
})

const saveSettings = () => {
    $q.loading.show()
    setTimeout(() => {
        $q.loading.hide()
        $q.notify({ type: 'positive', message: 'Impostazioni salvate' })
    }, 800)
}
</script>
