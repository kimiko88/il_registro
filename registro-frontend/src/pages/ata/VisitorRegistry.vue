<template>
  <q-page class="q-pa-md q-pa-lg-xl visitor-registry-page">
    <!-- Hero Header -->
    <div class="row items-center justify-between q-mb-lg gap-md">
      <div class="col-12 col-md-7">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge color="amber-9" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="meeting_room" size="14px" class="q-mr-xs" />
            {{ t('visitorRegistry.badge') || 'PORTINERIA & VIGILANZA • COLLABORATORE SCOLASTICO' }}
          </q-badge>
          <q-badge outline color="primary" class="q-px-sm q-py-xs text-weight-bold">
            {{ selectedDateFormatted }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none flex items-center">
          <q-icon name="door_front" color="amber-9" class="q-mr-sm" size="36px" />
          {{ t('visitorRegistry.title') || 'Registro Visitatori & Portineria' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          {{ t('visitorRegistry.subtitle') || 'Controllo accessi esterni, uscite anticipate studenti con delega e segnalazioni anomalie strutturali' }}
        </div>
      </div>

      <!-- Controls & Actions -->
      <div class="col-12 col-md-5 row items-center justify-end q-gutter-sm">
        <q-btn
          v-if="activeTab === 'visitors'"
          color="amber-9"
          icon="person_add"
          :label="t('visitorRegistry.newVisitorBtn') || 'Registra Visitatore'"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openNewVisitorDialog"
        />
        <q-btn
          v-else-if="activeTab === 'early_exits'"
          color="indigo-7"
          icon="exit_to_app"
          :label="t('visitorRegistry.newEarlyExitBtn') || 'Nuova Uscita Studente'"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openNewEarlyExitDialog"
        />
        <q-btn
          v-else-if="activeTab === 'maintenance'"
          color="deep-orange-7"
          icon="handyman"
          :label="t('visitorRegistry.newReportBtn') || 'Segnala Guasto'"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openNewMaintenanceDialog"
        />
        <q-btn
          flat
          round
          dense
          color="primary"
          icon="refresh"
          :loading="loading"
          @click="loadCurrentTabData"
        >
          <q-tooltip>{{ t('common.refresh') || 'Aggiorna' }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <q-tabs
      v-model="activeTab"
      dense
      class="bg-white rounded-xl shadow-1 text-slate-600 q-mb-lg"
      active-color="primary"
      indicator-color="primary"
      align="left"
      narrow-indicator
      @update:model-value="onTabChange"
    >
      <q-tab name="visitors" icon="badge" :label="t('visitorRegistry.tabVisitors') || 'Visitatori Esterni'" />
      <q-tab name="early_exits" icon="school" :label="t('visitorRegistry.tabEarlyExits') || 'Uscite Anticipate Studenti'" />
      <q-tab name="maintenance" icon="build" :label="t('visitorRegistry.tabMaintenance') || 'Segnalazioni Manutenzione'" />
    </q-tabs>

    <!-- Tab Panels -->
    <q-tab-panels v-model="activeTab" animated class="bg-transparent">
      <!-- PANEL 1: VISITATORI ESTERNI -->
      <q-tab-panel name="visitors" class="q-pa-none">
        <!-- Stats Summary -->
        <div class="row q-col-gutter-md q-mb-lg">
          <div class="col-12 col-sm-4">
            <q-card class="stat-card bg-emerald-50 border-emerald-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-emerald-100 text-emerald-800">
                  <q-icon name="location_on" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-emerald-800 text-weight-bold">{{ t('visitorRegistry.statPresent') || 'Attualmente in Sede' }}</div>
                  <div class="text-h5 text-weight-bolder text-emerald-900">{{ currentVisitors.length }}</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 col-sm-4">
            <q-card class="stat-card bg-blue-50 border-blue-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-blue-100 text-blue-800">
                  <q-icon name="logout" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-blue-800 text-weight-bold">{{ t('visitorRegistry.statExited') || 'Usciti Oggi' }}</div>
                  <div class="text-h5 text-weight-bolder text-blue-900">{{ exitedVisitors.length }}</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 col-sm-4">
            <q-card class="stat-card bg-amber-50 border-amber-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-amber-100 text-amber-800">
                  <q-icon name="groups" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-amber-800 text-weight-bold">{{ t('visitorRegistry.statTotal') || 'Totale Ingressi Oggi' }}</div>
                  <div class="text-h5 text-weight-bolder text-amber-900">{{ visitorsList.length }}</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>

        <!-- Table: Visitatori -->
        <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white">
          <q-table
            :rows="visitorsList"
            :columns="visitorColumns"
            row-key="id"
            :loading="loading"
            no-data-label="Nessun visitatore registrato per oggi"
            class="no-shadow"
            :pagination="{ rowsPerPage: 15 }"
          >
            <template v-slot:body-cell-purpose="props">
              <q-td :props="props">
                <q-chip :color="getPurposeColor(props.row.purpose)" text-color="white" size="sm" class="text-weight-bold">
                  {{ getPurposeLabel(props.row.purpose) }}
                </q-chip>
              </q-td>
            </template>

            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge v-if="!props.row.exit_time" color="positive" class="q-px-sm q-py-xs">
                  <q-icon name="check_circle" size="12px" class="q-mr-xs" /> In sede
                </q-badge>
                <q-badge v-else color="grey-6" class="q-px-sm q-py-xs">
                  Uscito ore {{ formatTime(props.row.exit_time) }}
                </q-badge>
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props" align="right">
                <q-btn
                  v-if="!props.row.exit_time"
                  color="negative"
                  flat
                  dense
                  icon="logout"
                  :label="t('visitorRegistry.exitAction') || 'Registra Uscita'"
                  no-caps
                  size="sm"
                  class="text-weight-bold"
                  @click="recordExit(props.row)"
                />
                <span v-else class="text-caption text-slate-400">—</span>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- PANEL 2: USCITE ANTICIPATE STUDENTI -->
      <q-tab-panel name="early_exits" class="q-pa-none">
        <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white">
          <q-table
            :rows="earlyExitsList"
            :columns="earlyExitColumns"
            row-key="id"
            :loading="loading"
            no-data-label="Nessuna uscita anticipata registrata oggi"
            class="no-shadow"
            :pagination="{ rowsPerPage: 15 }"
          >
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge v-if="!props.row.return_time" color="warning" text-color="dark" class="q-px-sm q-py-xs">
                  <q-icon name="directions_walk" size="12px" class="q-mr-xs" /> Fuori sede
                </q-badge>
                <q-badge v-else color="positive" class="q-px-sm q-py-xs">
                  Rientrato ore {{ formatTime(props.row.return_time) }}
                </q-badge>
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props" align="right">
                <q-btn
                  v-if="!props.row.return_time"
                  color="primary"
                  flat
                  dense
                  icon="keyboard_return"
                  :label="t('visitorRegistry.returnAction') || 'Segna Rientro'"
                  no-caps
                  size="sm"
                  class="text-weight-bold"
                  @click="recordReturn(props.row)"
                />
                <span v-else class="text-caption text-slate-400">—</span>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- PANEL 3: SEGNALAZIONI GUASTI -->
      <q-tab-panel name="maintenance" class="q-pa-none">
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-sm-6 col-md-3">
            <q-select
              v-model="maintenanceStatusFilter"
              :options="[
                { label: 'Tutte le segnalazioni', value: '' },
                { label: 'Aperte', value: 'aperto' },
                { label: 'In Lavorazione', value: 'in_lavorazione' },
                { label: 'Chiuse / Risolte', value: 'chiuso' }
              ]"
              emit-value
              map-options
              outlined
              dense
              bg-color="white"
              @update:model-value="loadMaintenance"
            />
          </div>
        </div>

        <div class="row q-col-gutter-md">
          <div v-for="rep in maintenanceList" :key="rep.id" class="col-12 col-md-6 col-lg-4">
            <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white full-height column justify-between">
              <q-card-section>
                <div class="row items-center justify-between q-mb-sm">
                  <q-badge :color="getPriorityColor(rep.priority)" class="text-uppercase text-weight-bold q-px-sm">
                    {{ rep.priority }}
                  </q-badge>
                  <q-badge :color="getMaintenanceStatusColor(rep.status)">
                    {{ getMaintenanceStatusLabel(rep.status) }}
                  </q-badge>
                </div>

                <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs flex items-center">
                  <q-icon name="place" color="primary" size="20px" class="q-mr-xs" />
                  {{ rep.location }}
                </div>
                <div class="text-caption text-slate-500 q-mb-sm text-capitalize">
                  Categoria: <strong>{{ rep.category }}</strong>
                </div>
                <div class="text-body2 text-slate-700 q-mb-md">
                  {{ rep.description }}
                </div>
              </q-card-section>

              <q-card-section class="border-t border-slate-100 row items-center justify-between q-pt-sm">
                <div class="text-caption text-slate-400">
                  {{ formatDate(rep.created_at) }}
                </div>

                <!-- Status updater -->
                <q-btn-dropdown
                  flat
                  dense
                  no-caps
                  color="primary"
                  :label="t('visitorRegistry.updateStatus') || 'Aggiorna Stato'"
                  size="sm"
                >
                  <q-list dense>
                    <q-item clickable v-close-popup @click="updateMaintStatus(rep, 'aperto')">
                      <q-item-section>Aperto</q-item-section>
                    </q-item>
                    <q-item clickable v-close-popup @click="updateMaintStatus(rep, 'in_lavorazione')">
                      <q-item-section>In Lavorazione</q-item-section>
                    </q-item>
                    <q-item clickable v-close-popup @click="updateMaintStatus(rep, 'chiuso')">
                      <q-item-section>Chiuso / Risolto</q-item-section>
                    </q-item>
                  </q-list>
                </q-btn-dropdown>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog: Nuovo Visitatore -->
    <q-dialog v-model="visitorDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="badge" color="amber-9" class="q-mr-sm" size="24px" />
            {{ t('visitorRegistry.dialogVisitorTitle') || 'Registra Ingresso Visitatore' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.nameLabel') || 'Nome & Cognome Visitatore *' }}</label>
            <q-input v-model="visitorForm.name" outlined dense placeholder="Mario Rossi" />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.docLabel') || 'Documento Identità' }}</label>
              <q-input v-model="visitorForm.document_id" outlined dense placeholder="CI / Patente / Passaporto" />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.badgeLabel') || 'N° Badge Assegnato' }}</label>
              <q-input v-model="visitorForm.badge_number" outlined dense placeholder="es. 04" />
            </div>
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.purposeLabel') || 'Motivo Visita *' }}</label>
              <q-select
                v-model="visitorForm.purpose"
                :options="purposeOptions"
                emit-value
                map-options
                outlined
                dense
              />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.hostLabel') || 'Persona / Ufficio da incontrare' }}</label>
              <q-input v-model="visitorForm.host_name" outlined dense placeholder="es. Segreteria, Prof. Bianchi" />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.notesLabel') || 'Note' }}</label>
            <q-input v-model="visitorForm.notes" outlined dense type="textarea" rows="2" />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" color="grey-7" v-close-popup />
          <q-btn
            color="amber-9"
            :label="t('visitorRegistry.saveVisitor') || 'Registra Ingresso'"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="saving"
            @click="submitVisitor"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Nuova Uscita Anticipata -->
    <q-dialog v-model="earlyExitDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="school" color="indigo-7" class="q-mr-sm" size="24px" />
            {{ t('visitorRegistry.dialogEarlyExitTitle') || 'Registra Uscita Anticipata Studente' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.studentLabel') || 'Studente *' }}</label>
            <q-select
              v-model="earlyExitForm.student_id"
              :options="studentOptions"
              option-value="id"
              option-label="name"
              emit-value
              map-options
              outlined
              dense
              use-input
              @filter="filterStudents"
              placeholder="Cerca studente..."
            />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-7">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.delegateeLabel') || 'Nome Delegato al ritiro *' }}</label>
              <q-input v-model="earlyExitForm.delegatee_name" outlined dense placeholder="Genitore o delegato" />
            </div>
            <div class="col-5">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.relLabel') || 'Relazione / Titolo' }}</label>
              <q-input v-model="earlyExitForm.delegate_rel" outlined dense placeholder="Madre / Padre / Tutore" />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.reasonLabel') || 'Motivazione' }}</label>
            <q-select
              v-model="earlyExitForm.reason_code"
              :options="[
                { label: 'Visita Medica', value: 'visita_medica' },
                { label: 'Motivi Familiari', value: 'motivi_familiari' },
                { label: 'Malessere Improvviso a Scuola', value: 'malessere' },
                { label: 'Altro', value: 'altro' }
              ]"
              emit-value
              map-options
              outlined
              dense
            />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.notesLabel') || 'Note' }}</label>
            <q-input v-model="earlyExitForm.notes" outlined dense type="textarea" rows="2" />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" color="grey-7" v-close-popup />
          <q-btn
            color="indigo-7"
            :label="t('visitorRegistry.saveEarlyExit') || 'Registra Uscita'"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="saving"
            @click="submitEarlyExit"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Segnalazione Manutenzione -->
    <q-dialog v-model="maintenanceDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="build" color="deep-orange-7" class="q-mr-sm" size="24px" />
            {{ t('visitorRegistry.dialogMaintTitle') || 'Segnala Guasto / Anomalia' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.locationLabel') || 'Luogo / Plesso / Aula *' }}</label>
            <q-input v-model="maintenanceForm.location" outlined dense placeholder="es. Bagni Piano 1 ala est, Aula 12" />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.categoryLabel') || 'Categoria *' }}</label>
              <q-select
                v-model="maintenanceForm.category"
                :options="[
                  { label: 'Elettrico', value: 'elettrico' },
                  { label: 'Idraulico', value: 'idraulico' },
                  { label: 'Strutturale / Infissi', value: 'strutturale' },
                  { label: 'Pulizia / Igiene', value: 'pulizia' },
                  { label: 'Informatica / LIM', value: 'informatica' },
                  { label: 'Altro', value: 'altro' }
                ]"
                emit-value
                map-options
                outlined
                dense
              />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.priorityLabel') || 'Priorità *' }}</label>
              <q-select
                v-model="maintenanceForm.priority"
                :options="[
                  { label: 'Bassa', value: 'bassa' },
                  { label: 'Media', value: 'media' },
                  { label: 'Alta', value: 'alta' },
                  { label: 'Urgente / Pericolo', value: 'urgente' }
                ]"
                emit-value
                map-options
                outlined
                dense
              />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.descLabel') || 'Descrizione Dettagliata *' }}</label>
            <q-input
              v-model="maintenanceForm.description"
              outlined
              dense
              type="textarea"
              rows="3"
              placeholder="Descrivere il guasto o l'anomalia rilevata"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" color="grey-7" v-close-popup />
          <q-btn
            color="deep-orange-7"
            :label="t('visitorRegistry.saveMaint') || 'Invia Segnalazione'"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="saving"
            @click="submitMaintenance"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import visitorService from '@/services/visitorService'
import userService from '@/services/userService'

const $q = useQuasar()
const { t, te, locale } = useI18n()

const activeTab = ref('visitors')
const loading = ref(false)
const saving = ref(false)

const todayStr = new Date().toISOString().substring(0, 10)
const selectedDateFormatted = computed(() => {
  const d = new Date()
  return d.toLocaleDateString(locale.value || 'it-IT', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
})

// Tab 1: Visitors
const visitorsList = ref([])
const visitorDialog = ref(false)
const visitorForm = ref({
  name: '',
  document_id: '',
  purpose: 'parent',
  host_name: '',
  badge_number: '',
  notes: ''
})

const purposeOptions = [
  { label: 'Genitore / Tutore', value: 'parent' },
  { label: 'Fornitore / Tecnico Esterno', value: 'supplier' },
  { label: 'Ente / Istituzione / ASL', value: 'institution' },
  { label: 'Altro', value: 'other' }
]

const currentVisitors = computed(() => visitorsList.value.filter(v => !v.exit_time))
const exitedVisitors = computed(() => visitorsList.value.filter(v => v.exit_time))

const visitorColumns = [
  { name: 'name', label: 'Nome Visitatore', field: 'name', align: 'left', sortable: true },
  { name: 'document_id', label: 'Documento', field: 'document_id', align: 'left' },
  { name: 'purpose', label: 'Motivo', field: 'purpose', align: 'left' },
  { name: 'host_name', label: 'Referente', field: 'host_name', align: 'left' },
  { name: 'badge_number', label: 'Badge', field: 'badge_number', align: 'center' },
  { name: 'entry_time', label: 'Ingresso', field: row => formatTime(row.entry_time), align: 'center' },
  { name: 'status', label: 'Stato', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

// Tab 2: Early Exits
const earlyExitsList = ref([])
const earlyExitDialog = ref(false)
const earlyExitForm = ref({
  student_id: '',
  delegatee_name: '',
  delegate_rel: '',
  reason_code: 'visita_medica',
  notes: ''
})
const studentList = ref([])
const studentOptions = ref([])

const earlyExitColumns = [
  { name: 'student_name', label: 'Studente', field: 'student_name', align: 'left', sortable: true },
  { name: 'class_name', label: 'Classe', field: 'class_name', align: 'center' },
  { name: 'delegatee_name', label: 'Ritirato da', field: 'delegatee_name', align: 'left' },
  { name: 'delegate_rel', label: 'Grado parentela', field: 'delegate_rel', align: 'left' },
  { name: 'exit_time', label: 'Ora Uscita', field: row => formatTime(row.exit_time), align: 'center' },
  { name: 'status', label: 'Stato', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

// Tab 3: Maintenance
const maintenanceList = ref([])
const maintenanceDialog = ref(false)
const maintenanceStatusFilter = ref('')
const maintenanceForm = ref({
  location: '',
  category: 'elettrico',
  priority: 'media',
  description: ''
})

function onTabChange(val) {
  loadCurrentTabData()
}

async function loadCurrentTabData() {
  if (activeTab.value === 'visitors') {
    await loadVisitors()
  } else if (activeTab.value === 'early_exits') {
    await loadEarlyExits()
  } else if (activeTab.value === 'maintenance') {
    await loadMaintenance()
  }
}

async function loadVisitors() {
  loading.value = true
  try {
    const res = await visitorService.listTodayVisitors(todayStr)
    visitorsList.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore caricamento visitatori' })
  } finally {
    loading.value = false
  }
}

async function loadEarlyExits() {
  loading.value = true
  try {
    const res = await visitorService.listTodayEarlyExits(todayStr)
    earlyExitsList.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore caricamento uscite anticipate' })
  } finally {
    loading.value = false
  }
}

async function loadMaintenance() {
  loading.value = true
  try {
    const res = await visitorService.listMaintenanceReports(maintenanceStatusFilter.value)
    maintenanceList.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore caricamento segnalazioni' })
  } finally {
    loading.value = false
  }
}

function openNewVisitorDialog() {
  visitorForm.value = {
    name: '',
    document_id: '',
    purpose: 'parent',
    host_name: '',
    badge_number: '',
    notes: ''
  }
  visitorDialog.value = true
}

async function submitVisitor() {
  if (!visitorForm.value.name) {
    $q.notify({ type: 'warning', message: 'Nome visitatore obbligatorio' })
    return
  }
  saving.value = true
  try {
    await visitorService.registerVisitor(visitorForm.value)
    $q.notify({ type: 'positive', message: 'Visitatore registrato con successo' })
    visitorDialog.value = false
    await loadVisitors()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore registrazione visitatore' })
  } finally {
    saving.value = false
  }
}

async function recordExit(v) {
  try {
    await visitorService.recordVisitorExit(v.id, 'Uscita registrata da portineria')
    $q.notify({ type: 'positive', message: `Uscita registrata per ${v.name}` })
    await loadVisitors()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore registrazione uscita' })
  }
}

function openNewEarlyExitDialog() {
  earlyExitForm.value = {
    student_id: '',
    delegatee_name: '',
    delegate_rel: 'Genitore',
    reason_code: 'visita_medica',
    notes: ''
  }
  earlyExitDialog.value = true
}

async function submitEarlyExit() {
  if (!earlyExitForm.value.student_id || !earlyExitForm.value.delegatee_name) {
    $q.notify({ type: 'warning', message: 'Compilare studente e nome del delegato' })
    return
  }
  saving.value = true
  try {
    await visitorService.recordEarlyExit(earlyExitForm.value)
    $q.notify({ type: 'positive', message: 'Uscita anticipata registrata' })
    earlyExitDialog.value = false
    await loadEarlyExits()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore registrazione uscita anticipata' })
  } finally {
    saving.value = false
  }
}

async function recordReturn(exit) {
  try {
    await visitorService.recordStudentReturn(exit.id, 'Rientro in aula')
    $q.notify({ type: 'positive', message: 'Rientro dello studente registrato' })
    await loadEarlyExits()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore rientro studente' })
  }
}

function openNewMaintenanceDialog() {
  maintenanceForm.value = {
    location: '',
    category: 'elettrico',
    priority: 'media',
    description: ''
  }
  maintenanceDialog.value = true
}

async function submitMaintenance() {
  if (!maintenanceForm.value.location || !maintenanceForm.value.description) {
    $q.notify({ type: 'warning', message: 'Luogo e descrizione sono obbligatori' })
    return
  }
  saving.value = true
  try {
    await visitorService.createMaintenanceReport(maintenanceForm.value)
    $q.notify({ type: 'positive', message: 'Segnalazione registrata con successo' })
    maintenanceDialog.value = false
    await loadMaintenance()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore registrazione guasto' })
  } finally {
    saving.value = false
  }
}

async function updateMaintStatus(rep, status) {
  try {
    await visitorService.updateMaintenanceStatus(rep.id, { status })
    $q.notify({ type: 'positive', message: 'Stato segnalazione aggiornato' })
    await loadMaintenance()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore aggiornamento stato' })
  }
}

function filterStudents(val, update) {
  if (val === '') {
    update(() => {
      studentOptions.value = studentList.value
    })
    return
  }
  update(() => {
    const needle = val.toLowerCase()
    studentOptions.value = studentList.value.filter(s => s.name.toLowerCase().includes(needle))
  })
}

function formatTime(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleTimeString(locale.value || 'it-IT', { hour: '2-digit', minute: '2-digit' })
}

function formatDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleDateString(locale.value || 'it-IT', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}

function getPurposeLabel(p) {
  switch (p) {
    case 'parent': return 'Genitore'
    case 'supplier': return 'Fornitore'
    case 'institution': return 'Istituzione'
    default: return 'Altro'
  }
}

function getPurposeColor(p) {
  switch (p) {
    case 'parent': return 'blue-7'
    case 'supplier': return 'amber-9'
    case 'institution': return 'purple-7'
    default: return 'grey-7'
  }
}

function getPriorityColor(p) {
  switch (p) {
    case 'urgente': return 'negative'
    case 'alta': return 'deep-orange-7'
    case 'media': return 'amber-8'
    case 'bassa': return 'blue-grey-6'
    default: return 'grey-6'
  }
}

function getMaintenanceStatusColor(st) {
  switch (st) {
    case 'aperto': return 'negative'
    case 'in_lavorazione': return 'warning'
    case 'chiuso': return 'positive'
    default: return 'grey-6'
  }
}

function getMaintenanceStatusLabel(st) {
  switch (st) {
    case 'aperto': return 'Aperto'
    case 'in_lavorazione': return 'In corso'
    case 'chiuso': return 'Risolto'
    default: return st
  }
}

onMounted(async () => {
  await loadVisitors()
  try {
    const sRes = await userService.getUsers({ role: 'student' })
    const raw = sRes.data?.users || sRes.data || []
    studentList.value = raw.map(u => ({
      id: u.id,
      name: `${u.last_name || ''} ${u.first_name || ''}`.trim() || u.email
    }))
    studentOptions.value = studentList.value
  } catch {
    // Non blocking
  }
})
</script>

<style scoped>
.visitor-registry-page {
  background: #f8fafc;
  min-height: 100vh;
}

.stat-card {
  border-radius: 16px;
  border-width: 1px;
  border-style: solid;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
