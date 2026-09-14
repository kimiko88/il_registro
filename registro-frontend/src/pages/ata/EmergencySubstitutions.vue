<template>
  <q-page class="q-pa-md q-pa-lg-xl emergency-substitutions-page">
    <!-- Hero Header -->
    <div class="row items-center justify-between q-mb-lg gap-md">
      <div class="col-12 col-md-7">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge color="deep-orange-7" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="warning_amber" size="14px" class="q-mr-xs" />
            {{ t('emergencySubstitutions.badge') || 'EMERGENZA MATTUTINA • COLLABORATORE DS' }}
          </q-badge>
          <q-badge outline color="primary" class="q-px-sm q-py-xs text-weight-bold">
            {{ selectedDateFormatted }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none flex items-center">
          <q-icon name="bolt" color="deep-orange-7" class="q-mr-sm" size="36px" />
          {{ t('emergencySubstitutions.title') || 'Emergenza Sostituzioni' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          {{ t('emergencySubstitutions.subtitle') || 'Centrale operativa mattutina per gestione rapida docenti assenti e copertura cattedre' }}
        </div>
      </div>

      <!-- Controls & Actions -->
      <div class="col-12 col-md-5 row items-center justify-end q-gutter-sm">
        <q-btn-group outline rounded class="bg-white shadow-1">
          <q-btn flat dense icon="chevron_left" @click="changeDate(-1)" :title="t('common.prev')">
            <q-tooltip>{{ t('staffAttendance.prevDay') || 'Giorno precedente' }}</q-tooltip>
          </q-btn>
          <q-btn flat no-caps class="text-weight-bold text-slate-700 q-px-sm">
            <q-icon name="event" color="primary" size="18px" class="q-mr-xs" />
            {{ selectedDate }}
            <q-popup-proxy cover transition-show="scale" transition-hide="scale">
              <q-date v-model="selectedDate" mask="YYYY-MM-DD" @update:model-value="fetchData" today-btn />
            </q-popup-proxy>
          </q-btn>
          <q-btn flat dense icon="chevron_right" @click="changeDate(1)" :title="t('common.next')">
            <q-tooltip>{{ t('staffAttendance.nextDay') || 'Giorno successivo' }}</q-tooltip>
          </q-btn>
        </q-btn-group>

        <q-btn
          color="deep-orange-7"
          icon="add_alert"
          :label="t('emergencySubstitutions.addAbsenceBtn') || 'Segnala Assenza'"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openAddAbsenceDialog"
        />
        <q-btn
          outline
          color="primary"
          icon="print"
          :label="t('emergencySubstitutions.printProspetto') || 'Stampa Prospetto'"
          no-caps
          rounded
          class="bg-white"
          @click="printProspetto"
        />
        <q-btn
          flat
          round
          dense
          color="primary"
          icon="refresh"
          :loading="loading"
          @click="fetchData"
        >
          <q-tooltip>{{ t('common.refresh') || 'Aggiorna' }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Quick Stats Pills -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-6 col-sm-3">
        <q-card class="stat-pill bg-red-50 border-red-200">
          <q-card-section class="q-pa-md row items-center no-wrap">
            <div class="stat-icon bg-red-100 text-red-700">
              <q-icon name="person_off" size="24px" />
            </div>
            <div class="q-ml-md">
              <div class="text-caption text-red-800 text-weight-bold text-uppercase">{{ t('emergencySubstitutions.statTotalAbsences') || 'Assenze Totali' }}</div>
              <div class="text-h5 text-weight-bolder text-red-900">{{ summary.total }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card class="stat-pill bg-amber-50 border-amber-200">
          <q-card-section class="q-pa-md row items-center no-wrap">
            <div class="stat-icon bg-amber-100 text-amber-800">
              <q-icon name="pending_actions" size="24px" />
            </div>
            <div class="q-ml-md">
              <div class="text-caption text-amber-800 text-weight-bold text-uppercase">{{ t('emergencySubstitutions.statPending') || 'Da Coprire' }}</div>
              <div class="text-h5 text-weight-bolder text-amber-900">{{ summary.pending }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card class="stat-pill bg-blue-50 border-blue-200">
          <q-card-section class="q-pa-md row items-center no-wrap">
            <div class="stat-icon bg-blue-100 text-blue-800">
              <q-icon name="assignment_ind" size="24px" />
            </div>
            <div class="q-ml-md">
              <div class="text-caption text-blue-800 text-weight-bold text-uppercase">{{ t('emergencySubstitutions.statAssigned') || 'Assegnate' }}</div>
              <div class="text-h5 text-weight-bolder text-blue-900">{{ summary.assigned }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card class="stat-pill bg-emerald-50 border-emerald-200">
          <q-card-section class="q-pa-md row items-center no-wrap">
            <div class="stat-icon bg-emerald-100 text-emerald-800">
              <q-icon name="check_circle" size="24px" />
            </div>
            <div class="q-ml-md">
              <div class="text-caption text-emerald-800 text-weight-bold text-uppercase">{{ t('emergencySubstitutions.statConfirmed') || 'Confermate / In Aula' }}</div>
              <div class="text-h5 text-weight-bolder text-emerald-900">{{ summary.confirmed }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Operative Panel: 3 Columns -->
    <div class="row q-col-gutter-lg">
      <!-- Col 1: Docenti Assenti & Da Assegnare -->
      <div class="col-12 col-lg-4">
        <q-card class="column-card shadow-sm rounded-2xl full-height">
          <q-card-section class="row items-center justify-between border-b border-slate-100 q-pb-sm">
            <div class="row items-center q-gutter-xs">
              <q-icon name="priority_high" color="negative" size="20px" />
              <div class="text-subtitle1 text-weight-bold text-slate-800">
                {{ t('emergencySubstitutions.colAbsencesTitle') || 'Assenze del Giorno' }}
              </div>
            </div>
            <q-badge color="negative" rounded>{{ pendingSubstitutions.length }}</q-badge>
          </q-card-section>

          <q-card-section class="q-pa-none scroll-area-panel">
            <div v-if="loading" class="q-pa-lg text-center">
              <q-spinner-dots color="primary" size="40px" />
            </div>
            <div v-else-if="substitutions.length === 0" class="q-pa-xl text-center text-slate-400">
              <q-icon name="check_circle_outline" size="48px" color="positive" class="q-mb-sm" />
              <div class="text-body2 text-weight-medium">{{ t('emergencySubstitutions.noAbsences') || 'Nessuna assenza registrata per oggi' }}</div>
            </div>
            <q-list v-else separator>
              <q-item
                v-for="sub in substitutions"
                :key="sub.id"
                clickable
                v-ripple
                :active="selectedSub?.id === sub.id"
                active-class="selected-item"
                @click="selectSubstitution(sub)"
                class="q-py-md"
              >
                <q-item-section avatar>
                  <q-avatar :color="getStatusColor(sub.status)" text-color="white" size="38px" class="text-weight-bold">
                    {{ (sub.absent_teacher_name || 'D').charAt(0).toUpperCase() }}
                  </q-avatar>
                </q-item-section>

                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-900">
                    {{ sub.absent_teacher_name || 'Docente Assente' }}
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    <span class="text-weight-medium text-slate-700">Ora {{ sub.hour || sub.slot }}</span> •
                    Classe {{ sub.class_name || 'N/D' }} • {{ sub.subject_name || 'Materia' }}
                  </q-item-label>
                  <q-item-label caption v-if="sub.notes" class="text-italic text-slate-400">
                    "{{ sub.notes }}"
                  </q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-badge :color="getStatusColor(sub.status)" class="q-px-xs q-py-none text-caption">
                    {{ getStatusLabel(sub.status) }}
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Col 2: Coperture in Corso / Confermate -->
      <div class="col-12 col-lg-4">
        <q-card class="column-card shadow-sm rounded-2xl full-height">
          <q-card-section class="row items-center justify-between border-b border-slate-100 q-pb-sm">
            <div class="row items-center q-gutter-xs">
              <q-icon name="how_to_reg" color="indigo-7" size="20px" />
              <div class="text-subtitle1 text-weight-bold text-slate-800">
                {{ t('emergencySubstitutions.colActiveTitle') || 'Coperture Assegnate' }}
              </div>
            </div>
            <q-badge color="indigo-7" rounded>{{ assignedSubstitutions.length }}</q-badge>
          </q-card-section>

          <q-card-section class="q-pa-none scroll-area-panel">
            <div v-if="assignedSubstitutions.length === 0" class="q-pa-xl text-center text-slate-400">
              <q-icon name="hourglass_empty" size="48px" class="q-mb-sm" />
              <div class="text-body2">{{ t('emergencySubstitutions.noAssigned') || 'Nessun sostituto assegnato al momento' }}</div>
            </div>
            <q-list v-else separator>
              <q-item v-for="sub in assignedSubstitutions" :key="sub.id" class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="indigo-1" text-color="indigo-9" icon="swap_horiz" size="38px" />
                </q-item-section>

                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-900">
                    {{ sub.substitute_teacher_name || 'Sostituto' }}
                  </q-item-label>
                  <q-item-label caption class="text-slate-600">
                    Sostituisce: <strong>{{ sub.absent_teacher_name }}</strong>
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    Ora {{ sub.hour || sub.slot }} • {{ sub.class_name }} • {{ sub.subject_name }}
                  </q-item-label>
                </q-item-section>

                <q-item-section side class="column items-end q-gutter-xs">
                  <q-badge :color="getStatusColor(sub.status)">
                    {{ getStatusLabel(sub.status) }}
                  </q-badge>
                  <q-btn
                    v-if="sub.status === 'assigned'"
                    flat
                    dense
                    color="positive"
                    icon="check"
                    :label="t('emergencySubstitutions.confirmBtn') || 'Conferma'"
                    no-caps
                    size="sm"
                    @click="confirmSub(sub)"
                  />
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Col 3: Raccomandazioni Automatiche & 1-Click Assegna -->
      <div class="col-12 col-lg-4">
        <q-card class="column-card shadow-sm rounded-2xl full-height bg-gradient-recommendation">
          <q-card-section class="row items-center justify-between border-b border-indigo-100 q-pb-sm">
            <div class="row items-center q-gutter-xs">
              <q-icon name="auto_awesome" color="amber-9" size="20px" />
              <div class="text-subtitle1 text-weight-bold text-slate-800">
                {{ t('emergencySubstitutions.recommendationsTitle') || 'Raccomandazioni Sistema' }}
              </div>
            </div>
            <q-badge v-if="selectedSub" color="amber-8" text-color="white">
              Ora {{ selectedSub.hour || selectedSub.slot }}
            </q-badge>
          </q-card-section>

          <q-card-section class="q-pa-md">
            <div v-if="!selectedSub" class="q-pa-xl text-center text-slate-400">
              <q-icon name="touch_app" size="48px" class="q-mb-sm text-indigo-300" />
              <div class="text-body2 text-weight-medium">
                {{ t('emergencySubstitutions.selectToRecommend') || 'Seleziona un\'assenza a sinistra per visualizzare i docenti consigliati per la sostituzione' }}
              </div>
            </div>

            <div v-else>
              <div class="selected-summary-card q-pa-sm rounded-borders bg-white shadow-1 q-mb-md">
                <div class="text-caption text-slate-500">{{ t('emergencySubstitutions.targetAbsence') || 'Cattedra da coprire:' }}</div>
                <div class="text-weight-bold text-slate-900">{{ selectedSub.absent_teacher_name }}</div>
                <div class="text-caption text-slate-600">
                  Classe {{ selectedSub.class_name }} • {{ selectedSub.subject_name }} • Ora {{ selectedSub.hour || selectedSub.slot }}
                </div>
              </div>

              <div v-if="loadingRecommendations" class="q-pa-lg text-center">
                <q-spinner-hourglass color="primary" size="36px" />
                <div class="text-caption text-slate-500 q-mt-xs">{{ t('emergencySubstitutions.calculating') || 'Calcolo docenti disponibili in corso...' }}</div>
              </div>

              <div v-else-if="recommendations.length === 0" class="q-pa-md text-center text-slate-500 bg-white rounded-borders">
                <q-icon name="info" color="grey-6" size="32px" class="q-mb-xs" />
                <div>{{ t('emergencySubstitutions.noCandidates') || 'Nessun docente libero raccomandato in questo slot.' }}</div>
                <q-btn
                  color="primary"
                  flat
                  no-caps
                  icon="search"
                  :label="t('emergencySubstitutions.manualAssign') || 'Assegna manualmente'"
                  class="q-mt-sm"
                  @click="openManualAssignDialog"
                />
              </div>

              <q-list v-else separator class="bg-white rounded-borders shadow-1">
                <q-item v-for="rec in recommendations" :key="rec.teacher_id" class="q-py-md">
                  <q-item-section avatar>
                    <q-avatar color="amber-1" text-color="amber-9" icon="person" size="36px" />
                  </q-item-section>

                  <q-item-section>
                    <q-item-label class="text-weight-bold text-slate-900">
                      {{ rec.teacher_name }}
                    </q-item-label>
                    <q-item-label caption class="text-slate-600">
                      {{ rec.reason || 'Disponibile • Ore a disposizione' }}
                    </q-item-label>
                    <q-item-label caption class="row items-center text-amber-800 text-weight-bold">
                      <q-icon name="stars" size="14px" class="q-mr-xs" />
                      Score: {{ rec.score }}
                    </q-item-label>
                  </q-item-section>

                  <q-item-section side>
                    <q-btn
                      color="deep-orange-7"
                      icon="flash_on"
                      :label="t('emergencySubstitutions.quickAssign') || 'Assegna'"
                      no-caps
                      dense
                      rounded
                      class="q-px-sm text-weight-bold"
                      :loading="assigningId === rec.teacher_id"
                      @click="quickAssign(rec)"
                    />
                  </q-item-section>
                </q-item>
              </q-list>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Modal: Segnala Assenza Improvvisa -->
    <q-dialog v-model="absenceDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="person_remove" color="negative" class="q-mr-sm" size="24px" />
            {{ t('emergencySubstitutions.dialogTitle') || 'Nuova Assenza Docente' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('emergencySubstitutions.absentTeacher') || 'Docente Assente *' }}</label>
            <q-select
              v-model="absenceForm.absent_teacher_id"
              :options="teacherOptions"
              option-value="id"
              option-label="name"
              emit-value
              map-options
              outlined
              dense
              use-input
              @filter="filterTeachers"
              :placeholder="t('emergencySubstitutions.searchTeacher') || 'Cerca docente...'"
            />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('emergencySubstitutions.hour') || 'Ora di lezione *' }}</label>
              <q-select
                v-model.number="absenceForm.hour"
                :options="[1, 2, 3, 4, 5, 6, 7, 8]"
                outlined
                dense
              />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('emergencySubstitutions.class') || 'Classe' }}</label>
              <q-select
                v-model="absenceForm.class_id"
                :options="classOptions"
                option-value="id"
                option-label="name"
                emit-value
                map-options
                outlined
                dense
              />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('emergencySubstitutions.notes') || 'Note / Motivo' }}</label>
            <q-input
              v-model="absenceForm.notes"
              outlined
              dense
              type="textarea"
              rows="2"
              placeholder="es. Comunicazione telefonica ore 07:45 per malessere improvviso"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" color="grey-7" v-close-popup />
          <q-btn
            color="deep-orange-7"
            :label="t('emergencySubstitutions.saveAbsence') || 'Registra Assenza'"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="savingAbsence"
            @click="submitAbsence"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Assegna Manualmente -->
    <q-dialog v-model="manualAssignDialog">
      <q-card style="min-width: 420px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900">
            {{ t('emergencySubstitutions.manualAssignTitle') || 'Assegna Sostituto Manualmente' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('emergencySubstitutions.chooseSubstitute') || 'Docente Sostituto *' }}</label>
            <q-select
              v-model="manualSubstituteId"
              :options="teacherOptions"
              option-value="id"
              option-label="name"
              emit-value
              map-options
              outlined
              dense
              use-input
              @filter="filterTeachers"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" color="grey-7" v-close-popup />
          <q-btn
            color="primary"
            :label="t('emergencySubstitutions.confirmAssign') || 'Conferma Assegnazione'"
            no-caps
            rounded
            :loading="assigningId !== null"
            @click="submitManualAssign"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Print Table Template (Hidden on screen, shown in @media print) -->
    <div id="print-prospetto" class="print-only">
      <h2 class="text-center">PROSPETTO GIORNALIERO SOSTITUZIONI DOCENTI</h2>
      <p class="text-center">Data: {{ selectedDate }}</p>
      <table class="print-table">
        <thead>
          <tr>
            <th>Ora</th>
            <th>Docente Assente</th>
            <th>Classe</th>
            <th>Materia</th>
            <th>Docente Sostituto</th>
            <th>Stato</th>
            <th>Firma</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in substitutions" :key="s.id">
            <td>{{ s.hour || s.slot }}</td>
            <td>{{ s.absent_teacher_name }}</td>
            <td>{{ s.class_name }}</td>
            <td>{{ s.subject_name }}</td>
            <td>{{ s.substitute_teacher_name || 'IN ATTESA' }}</td>
            <td>{{ getStatusLabel(s.status) }}</td>
            <td>________________</td>
          </tr>
        </tbody>
      </table>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import substitutionService from '@/services/substitutionService'
import userService from '@/services/userService'
import schoolService from '@/services/schoolService'

const $q = useQuasar()
const { t, te, locale } = useI18n()

const selectedDate = ref(new Date().toISOString().substring(0, 10))
const loading = ref(false)
const savingAbsence = ref(false)
const assigningId = ref(null)
const loadingRecommendations = ref(false)

const substitutions = ref([])
const summary = ref({
  total: 0,
  pending: 0,
  assigned: 0,
  confirmed: 0
})

const selectedSub = ref(null)
const recommendations = ref([])

const absenceDialog = ref(false)
const manualAssignDialog = ref(false)
const manualSubstituteId = ref(null)

const absenceForm = ref({
  absent_teacher_id: '',
  hour: 1,
  class_id: '',
  notes: ''
})

const teacherList = ref([])
const teacherOptions = ref([])
const classOptions = ref([])

const selectedDateFormatted = computed(() => {
  const d = new Date(selectedDate.value)
  return d.toLocaleDateString(locale.value || 'it-IT', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
})

const pendingSubstitutions = computed(() => {
  return substitutions.value.filter(s => s.status === 'pending')
})

const assignedSubstitutions = computed(() => {
  return substitutions.value.filter(s => s.status === 'assigned' || s.status === 'confirmed')
})

function changeDate(days) {
  const d = new Date(selectedDate.value)
  d.setDate(d.getDate() + days)
  selectedDate.value = d.toISOString().substring(0, 10)
  fetchData()
}

async function fetchData() {
  loading.value = true
  try {
    const res = await substitutionService.todaySummary(selectedDate.value)
    const data = res.data || {}
    summary.value = {
      total: data.total || 0,
      pending: data.pending || 0,
      assigned: data.assigned || 0,
      confirmed: data.confirmed || 0
    }
    substitutions.value = data.substitutions || []

    // If previously selected sub still exists, update reference
    if (selectedSub.value) {
      const updated = substitutions.value.find(s => s.id === selectedSub.value.id)
      if (updated) {
        selectedSub.value = updated
      } else {
        selectedSub.value = null
        recommendations.value = []
      }
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('emergencySubstitutions.fetchError') || 'Errore caricamento sostituzioni'
    })
  } finally {
    loading.value = false
  }
}

async function selectSubstitution(sub) {
  selectedSub.value = sub
  if (sub.status === 'confirmed') {
    recommendations.value = []
    return
  }

  loadingRecommendations.value = true
  try {
    const res = await substitutionService.recommendSubstitutes({
      class_id: sub.class_id,
      subject_id: sub.subject_id,
      date: selectedDate.value,
      hour: sub.hour || sub.slot
    })
    recommendations.value = res.data || []
  } catch {
    recommendations.value = []
  } finally {
    loadingRecommendations.value = false
  }
}

async function quickAssign(rec) {
  if (!selectedSub.value) return
  assigningId.value = rec.teacher_id
  try {
    await substitutionService.assign(selectedSub.value.id, {
      substitute_teacher_id: rec.teacher_id,
      notes: rec.reason || 'Assegnazione rapida collaboratore DS'
    })
    $q.notify({
      type: 'positive',
      message: `${rec.teacher_name} assegnato con successo!`
    })
    await fetchData()
    if (selectedSub.value) {
      await selectSubstitution(selectedSub.value)
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore durante l\'assegnazione'
    })
  } finally {
    assigningId.value = null
  }
}

async function confirmSub(sub) {
  try {
    await substitutionService.signRegister(sub.id, 'Confermato da Collaboratore DS')
    $q.notify({
      type: 'positive',
      message: 'Sostituzione confermata in aula!'
    })
    await fetchData()
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore conferma'
    })
  }
}

function openAddAbsenceDialog() {
  absenceForm.value = {
    absent_teacher_id: '',
    hour: 1,
    class_id: classOptions.value[0]?.id || '',
    notes: ''
  }
  absenceDialog.value = true
}

async function submitAbsence() {
  if (!absenceForm.value.absent_teacher_id) {
    $q.notify({ type: 'warning', message: 'Seleziona il docente assente' })
    return
  }
  savingAbsence.value = true
  try {
    await substitutionService.create({
      absent_teacher_id: absenceForm.value.absent_teacher_id,
      date: selectedDate.value,
      hour: Number(absenceForm.value.hour),
      slot: Number(absenceForm.value.hour),
      class_id: absenceForm.value.class_id,
      notes: absenceForm.value.notes
    })
    $q.notify({ type: 'positive', message: 'Assenza registrata con successo!' })
    absenceDialog.value = false
    await fetchData()
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore registrazione assenza'
    })
  } finally {
    savingAbsence.value = false
  }
}

function openManualAssignDialog() {
  manualSubstituteId.value = null
  manualAssignDialog.value = true
}

async function submitManualAssign() {
  if (!manualSubstituteId.value || !selectedSub.value) return
  assigningId.value = manualSubstituteId.value
  try {
    await substitutionService.assign(selectedSub.value.id, {
      substitute_teacher_id: manualSubstituteId.value,
      notes: 'Assegnazione manuale Collaboratore DS'
    })
    $q.notify({ type: 'positive', message: 'Docente assegnato con successo!' })
    manualAssignDialog.value = false
    await fetchData()
    if (selectedSub.value) {
      await selectSubstitution(selectedSub.value)
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore assegnazione'
    })
  } finally {
    assigningId.value = null
  }
}

function filterTeachers(val, update) {
  if (val === '') {
    update(() => {
      teacherOptions.value = teacherList.value
    })
    return
  }
  update(() => {
    const needle = val.toLowerCase()
    teacherOptions.value = teacherList.value.filter(v => v.name.toLowerCase().includes(needle))
  })
}

function getStatusColor(st) {
  switch (st) {
    case 'pending': return 'deep-orange-7'
    case 'assigned': return 'indigo-7'
    case 'confirmed': return 'positive'
    default: return 'grey-7'
  }
}

function getStatusLabel(st) {
  switch (st) {
    case 'pending': return t('emergencySubstitutions.statusPending') || 'In attesa'
    case 'assigned': return t('emergencySubstitutions.statusAssigned') || 'Assegnata'
    case 'confirmed': return t('emergencySubstitutions.statusConfirmed') || 'Confermata'
    default: return st
  }
}

function printProspetto() {
  window.print()
}

onMounted(async () => {
  await fetchData()
  try {
    const [uRes, cRes] = await Promise.allSettled([
      userService.getUsers({ role: 'teacher' }),
      schoolService.getClasses()
    ])
    if (uRes.status === 'fulfilled') {
      const raw = uRes.value.data?.users || uRes.value.data || []
      teacherList.value = raw.map(u => ({
        id: u.id,
        name: `${u.last_name || ''} ${u.first_name || ''}`.trim() || u.email
      }))
      teacherOptions.value = teacherList.value
    }
    if (cRes.status === 'fulfilled') {
      const rawC = cRes.value.data || []
      classOptions.value = rawC.map(c => ({ id: c.id, name: c.name || `Classe ${c.year || ''}` }))
    }
  } catch {
    // Non blocking
  }
})
</script>

<style scoped>
.emergency-substitutions-page {
  background: #f8fafc;
  min-height: 100vh;
}

.stat-pill {
  border-radius: 16px;
  border-width: 1px;
  border-style: solid;
  transition: transform 0.2s ease;
}
.stat-pill:hover {
  transform: translateY(-2px);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.column-card {
  border: 1px solid #e2e8f0;
  background: white;
}

.scroll-area-panel {
  max-height: 520px;
  overflow-y: auto;
}

.selected-item {
  background: #eff6ff !important;
  border-left: 4px solid #3b82f6;
}

.bg-gradient-recommendation {
  background: linear-gradient(180deg, #ffffff 0%, #fdfcfb 100%);
}

.selected-summary-card {
  border: 1px solid #cbd5e1;
}

/* Print Styles */
@media screen {
  .print-only {
    display: none;
  }
}

@media print {
  body * {
    visibility: hidden;
  }
  #print-prospetto, #print-prospetto * {
    visibility: visible;
  }
  #print-prospetto {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    display: block !important;
  }
  .print-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
  }
  .print-table th, .print-table td {
    border: 1px solid #333;
    padding: 8px;
    text-align: left;
  }
  .print-table th {
    background: #eee;
  }
}
</style>
