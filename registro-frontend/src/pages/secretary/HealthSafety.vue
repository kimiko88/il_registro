<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('healthSafety.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('healthSafety.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm">
        <q-btn
          v-if="activeTab === 'medications'"
          color="primary"
          icon="medical_services"
          :label="$t('healthSafety.registerMedBtn')"
          class="rounded-lg shadow-sm"
          @click="showMedDialog = true"
        />
        <q-btn
          v-if="activeTab === 'accidents'"
          color="negative"
          icon="emergency"
          label="Nuovo Infortunio a Scuola"
          class="rounded-lg shadow-sm"
          @click="showAccidentDialog = true"
        />
      </div>
    </div>

    <!-- TABS -->
    <q-card class="glass-card shadow-soft border-slate-100 overflow-hidden q-mb-xl">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey-7 bg-white"
        active-color="primary"
        indicator-color="primary"
        align="left"
      >
        <q-tab name="medications" icon="medication" :label="$t('healthSafety.medicationsTab')" />
        <q-tab name="accidents" icon="warning" :label="$t('healthSafety.accidentsTab')" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated class="bg-transparent">
        <!-- 1. FARMACI A SCUOLA -->
        <q-tab-panel name="medications" class="q-pa-md">
          <q-table
            :rows="medications"
            :columns="medColumns"
            row-key="id"
            flat
            dense
            :pagination="{ rowsPerPage: 10 }"
          >
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge :color="props.value === 'attivo' ? 'positive' : 'grey-7'">
                  {{ props.value.toUpperCase() }}
                </q-badge>
                <q-badge v-if="props.row.is_medication_expired" color="negative" class="q-ml-xs">
                  Farmaco Scaduto
                </q-badge>
                <q-badge v-if="props.row.is_prescription_due" color="warning" class="q-ml-xs">
                  Certificato Medico Scaduto
                </q-badge>
              </q-td>
            </template>
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" auto-width>
                <q-btn
                  flat
                  dense
                  color="primary"
                  icon="vaccines"
                  label="Somministra"
                  @click="openAdministerDialog(props.row)"
                />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- 2. INFORTUNI & INAIL -->
        <q-tab-panel name="accidents" class="q-pa-md">
          <div class="q-mb-md bg-red-50 border border-red-200 rounded-borders q-pa-sm text-red-9 text-caption row items-center">
            <q-icon name="priority_high" size="sm" class="q-mr-sm" />
            <span>{{ $t('healthSafety.inailAlert48h') }}</span>
          </div>

          <q-table
            :rows="accidents"
            :columns="accidentColumns"
            row-key="id"
            flat
            dense
            :pagination="{ rowsPerPage: 10 }"
          >
            <template v-slot:body-cell-inail="props">
              <q-td :props="props">
                <q-badge v-if="props.row.inail_reported" color="positive" icon="check">
                  Trasmesso INAIL ({{ props.row.inail_protocol }})
                </q-badge>
                <q-badge v-else-if="props.row.is_urgent_48h_alert" color="negative" icon="timer">
                  URGENTE: {{ Math.round(props.row.hours_remaining_inail) }}h al termine INAIL
                </q-badge>
                <span v-else class="text-grey-6">In attesa referto</span>
              </q-td>
            </template>
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" auto-width>
                <q-btn
                  v-if="!props.row.inail_reported"
                  flat
                  dense
                  color="secondary"
                  icon="send"
                  label="Invia INAIL"
                  @click="openInailDialog(props.row)"
                />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- DIALOG REGISTRA FARMACO -->
    <q-dialog v-model="showMedDialog">
      <q-card style="min-width: 450px" class="q-pa-md rounded-xl">
        <q-card-section>
          <div class="text-h6 text-weight-bold">Registra Piano Terapeutico Farmaco a Scuola</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="newMed.medication_name" outlined dense label="Nome Farmaco (es. Ventolin, Baqsimi, EpiPen)" />
          <q-input v-model="newMed.dosage" outlined dense label="Dosaggio (es. 2 puff / 1 iniezione)" />
          <q-input v-model="newMed.administration_mode" outlined dense label="Modalità Somministrazione (es. Inalatore, Autoiniettore)" />
          <q-input v-model="newMed.diagnosis_reason" outlined dense label="Diagnosi / Patologia (es. Asma, Shock anafilattico)" />
          <q-input v-model="newMed.emergency_protocol" outlined dense label="Protocollo di Emergenza e Istruzioni" />
          <q-input v-model="newMed.prescription_expiry" type="date" outlined dense label="Scadenza Certificato Medico" />
          <q-input v-model="newMed.medication_expiry_date" type="date" outlined dense label="Scadenza Confezione Farmaco in Plesso" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="$t('common.save')" @click="saveMedication" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- DIALOG SOMMINISTRA FARMACO -->
    <q-dialog v-model="showAdministerDialog">
      <q-card style="min-width: 350px" class="q-pa-md rounded-xl">
        <q-card-section>
          <div class="text-h6 text-weight-bold">Registra Somministrazione Farmaco</div>
          <div class="text-caption text-slate-500">{{ selectedMed?.medication_name }}</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="adminData.dose_given" outlined dense label="Dose somministrata" />
          <q-input v-model="adminData.notes" outlined dense label="Note sull'evento e reazione dello studente" />
          <q-toggle v-model="adminData.emergency_services_called" label="Chiamata 118 / Emergenza effettuata" color="negative" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="$t('common.save')" @click="confirmAdministration" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- DIALOG TRASMISSIONE INAIL -->
    <q-dialog v-model="showInailDialog">
      <q-card style="min-width: 380px" class="q-pa-md rounded-xl">
        <q-card-section>
          <div class="text-h6 text-weight-bold">Protocollo Denuncia Telematica INAIL</div>
          <div class="text-caption text-slate-500">{{ selectedAccident?.accident_number }}</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="inailProtocolInput" outlined dense label="Numero Protocollo / Ricevuta Telematica INAIL" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn color="negative" label="Registra Protocollo" @click="confirmInail" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const activeTab = ref('medications')
const medications = ref([])
const accidents = ref([])

const showMedDialog = ref(false)
const newMed = ref({
  student_id: '00000000-0000-0000-0000-000000000001',
  class_id: '00000000-0000-0000-0000-000000000002',
  medication_name: '',
  dosage: '',
  administration_mode: '',
  diagnosis_reason: '',
  emergency_protocol: '',
  prescription_expiry: new Date().toISOString().substring(0, 10),
  medication_expiry_date: new Date().toISOString().substring(0, 10)
})

const showAdministerDialog = ref(false)
const selectedMed = ref(null)
const adminData = ref({
  dose_given: '',
  notes: '',
  emergency_services_called: false
})

const showInailDialog = ref(false)
const selectedAccident = ref(null)
const inailProtocolInput = ref('')

const medColumns = [
  { name: 'name', label: 'Nome Farmaco', field: 'medication_name', align: 'left', sortable: true },
  { name: 'dosage', label: 'Dosaggio', field: 'dosage', align: 'left' },
  { name: 'mode', label: 'Modalità', field: 'administration_mode', align: 'left' },
  { name: 'diagnosis', label: 'Patologia', field: 'diagnosis_reason', align: 'left' },
  { name: 'expiry', label: 'Scadenza Farmaco', field: 'medication_expiry_date', align: 'center' },
  { name: 'status', label: 'Stato / Alert', field: 'status', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

const accidentColumns = [
  { name: 'number', label: 'Codice Infortunio', field: 'accident_number', align: 'left' },
  { name: 'time', label: 'Data & Ora', field: 'accident_time', align: 'left' },
  { name: 'location', label: 'Luogo', field: 'location', align: 'left' },
  { name: 'dynamics', label: 'Dinamica Evento', field: 'dynamics_description', align: 'left' },
  { name: 'injuries', label: 'Lesioni', field: 'reported_injuries', align: 'left' },
  { name: 'inail', label: 'Termine 48h / INAIL', field: 'inail', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

async function loadData() {
  try {
    const medRes = await api.get('/health-safety/medications')
    if (medRes.data?.medications) medications.value = medRes.data.medications

    const accRes = await api.get('/health-safety/accidents')
    if (accRes.data?.accidents) accidents.value = accRes.data.accidents
  } catch (err) {
    console.error('Error loading health & safety data:', err)
  }
}

async function saveMedication() {
  try {
    await api.post('/health-safety/medications', newMed.value)
    $q.notify({ type: 'positive', message: 'Piano terapeutico registrato con successo' })
    showMedDialog.value = false
    await loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel salvataggio' })
  }
}

function openAdministerDialog(med) {
  selectedMed.value = med
  adminData.value.dose_given = med.dosage
  showAdministerDialog.value = true
}

async function confirmAdministration() {
  if (!selectedMed.value) return
  try {
    await api.post(`/health-safety/medications/${selectedMed.value.id}/administer`, adminData.value)
    $q.notify({ type: 'positive', message: 'Somministrazione registrata nel log ufficiale' })
    showAdministerDialog.value = false
    await loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nella registrazione' })
  }
}

function openInailDialog(acc) {
  selectedAccident.value = acc
  inailProtocolInput.value = ''
  showInailDialog.value = true
}

async function confirmInail() {
  if (!selectedAccident.value || !inailProtocolInput.value) return
  try {
    await api.post(`/health-safety/accidents/${selectedAccident.value.id}/transmit-inail`, {
      protocol: inailProtocolInput.value
    })
    $q.notify({ type: 'positive', message: 'Trasmissione INAIL registrata' })
    showInailDialog.value = false
    await loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante la registrazione INAIL' })
  }
}

onMounted(() => {
  loadData()
})
</script>
