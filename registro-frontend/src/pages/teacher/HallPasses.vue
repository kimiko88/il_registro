<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('securityPasses.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('securityPasses.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-select
          v-model="selectedClassId"
          :options="classesList"
          option-value="id"
          option-label="name"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          label="Classe Selezionata"
          style="min-width: 180px"
          @update:model-value="onClassChanged"
        />

        <q-btn
          color="primary"
          icon="timer"
          :label="$t('securityPasses.issuePassBtn')"
          class="rounded-lg shadow-sm"
          :disable="!selectedClassId"
          @click="openIssueDialog"
        />
      </div>
    </div>

    <!-- Active Passes In Progress -->
    <div class="q-mb-xl">
      <div class="row items-center justify-between q-mb-md">
        <div class="text-h6 text-weight-bold text-slate-800 row items-center q-gutter-x-sm">
          <q-icon name="directions_walk" color="primary" size="md" />
          <span>Studenti Attualmente Fuori Aula ({{ activePasses.length }})</span>
        </div>
        <q-btn flat round dense icon="refresh" color="primary" @click="loadPasses">
          <q-tooltip>Aggiorna</q-tooltip>
        </q-btn>
      </div>

      <div v-if="activePasses.length === 0" class="q-pa-lg text-center bg-white rounded-xl border-dashed border-slate-200">
        <q-icon name="check_circle" size="48px" color="positive" class="q-mb-sm" />
        <div class="text-subtitle1 text-weight-medium text-slate-700">Tutti gli studenti sono in aula</div>
        <div class="text-caption text-slate-400">Nessun pass fuori classe attivo per questa classe</div>
      </div>

      <div v-else class="row q-col-gutter-md">
        <div v-for="pass in activePasses" :key="pass.id" class="col-12 col-md-6 col-lg-4">
          <q-card
            class="shadow-soft rounded-xl transition-all"
            :class="[
              getElapsedMinutes(pass.departed_at) > 15
                ? 'border-2 border-negative bg-red-50/40'
                : 'border-slate-100 bg-white'
            ]"
          >
            <q-card-section>
              <div class="row items-center justify-between q-mb-xs">
                <div class="text-subtitle1 text-weight-bold text-slate-800">
                  {{ pass.student_name || 'Studente' }}
                </div>
                <q-badge
                  :color="getPassTypeColor(pass.pass_type)"
                  class="q-pa-xs text-weight-bold"
                >
                  {{ getPassTypeLabel(pass.pass_type) }}
                </q-badge>
              </div>

              <div class="text-caption text-slate-500 q-mb-sm">
                {{ $t('securityPasses.departedAt') }}: {{ formatTime(pass.departed_at) }}
              </div>

              <!-- Timer and Anomaly Banner -->
              <div class="row items-center justify-between q-pa-sm rounded-lg q-mb-md" :class="getElapsedMinutes(pass.departed_at) > 15 ? 'bg-red-100 text-red-900' : 'bg-slate-100 text-slate-700'">
                <div class="row items-center q-gutter-x-xs">
                  <q-icon :name="getElapsedMinutes(pass.departed_at) > 15 ? 'warning' : 'schedule'" size="sm" :color="getElapsedMinutes(pass.departed_at) > 15 ? 'negative' : 'primary'" />
                  <span class="text-weight-bold text-body2">
                    {{ getElapsedMinutes(pass.departed_at) }} min
                  </span>
                </div>
                <div v-if="getElapsedMinutes(pass.departed_at) > 15" class="text-caption text-weight-bold text-negative">
                  &gt; 15m ALERT!
                </div>
              </div>

              <div v-if="pass.notes" class="text-caption text-slate-600 q-mb-md italic">
                "{{ pass.notes }}"
              </div>

              <q-btn
                color="positive"
                icon="assignment_turned_in"
                :label="$t('securityPasses.markReturned')"
                class="full-width rounded-lg"
                unelevated
                @click="markReturned(pass.id)"
              />
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <!-- History Table -->
    <q-card class="glass-card shadow-soft border-slate-100 rounded-xl">
      <q-card-section>
        <div class="text-h6 text-weight-bold text-slate-800 q-mb-md row items-center q-gutter-x-sm">
          <q-icon name="history" color="slate-600" size="sm" />
          <span>Storico Passaggi Fuori Aula</span>
        </div>

        <q-table
          :rows="historyPasses"
          :columns="historyColumns"
          row-key="id"
          flat
          dense
          :pagination="{ rowsPerPage: 15 }"
          no-data-label="Nessun passaggio registrato di recente"
        >
          <template #body-cell-pass_type="props">
            <q-td :props="props">
              <q-badge :color="getPassTypeColor(props.value)" class="text-weight-medium">
                {{ getPassTypeLabel(props.value) }}
              </q-badge>
            </q-td>
          </template>

          <template #body-cell-duration="props">
            <q-td :props="props" align="center">
              <span v-if="props.row.returned_at" class="text-weight-bold">
                {{ calculateDuration(props.row.departed_at, props.row.returned_at) }} min
              </span>
              <q-badge v-else color="warning" label="In corso" />
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog: Issue New Pass -->
    <q-dialog v-model="showIssueDialog" persistent>
      <q-card style="min-width: 440px" class="rounded-xl">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6 text-weight-bold text-slate-800">
            {{ $t('securityPasses.issuePassBtn') }}
          </div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <q-select
            v-model="newPass.student_id"
            :options="studentsInClass"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            label="Seleziona Alunno"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <q-select
            v-model="newPass.pass_type"
            :options="passTypeOptions"
            emit-value
            map-options
            outlined
            dense
            :label="$t('securityPasses.passType')"
          />

          <q-input
            v-model="newPass.notes"
            outlined
            dense
            type="textarea"
            rows="2"
            label="Note opzionali (es. autorizzato dal docente)"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            icon="check"
            label="Rilascia Pass"
            :disable="!newPass.student_id"
            @click="issuePass"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'

const $q = useQuasar()
const { t } = useI18n()

const selectedClassId = ref(null)
const classesList = ref([])
const studentsInClass = ref([])

const activePasses = ref([])
const historyPasses = ref([])
const nowTime = ref(Date.now())
let timerInterval = null

const showIssueDialog = ref(false)
const newPass = ref({
  student_id: '',
  pass_type: 'bathroom',
  notes: ''
})

const passTypeOptions = [
  { label: 'Servizi Igienici (Bagno)', value: 'bathroom' },
  { label: 'Infermeria / Primo Soccorso', value: 'infirmary' },
  { label: 'Segreteria / Dirigenza', value: 'office' },
  { label: 'Collaboratori Scolastici ATA', value: 'staff' }
]

const historyColumns = [
  { name: 'student_name', label: 'Studente', field: 'student_name', align: 'left', sortable: true },
  { name: 'pass_type', label: 'Destinazione', field: 'pass_type', align: 'left' },
  { name: 'departed_at', label: 'Ora Uscita', field: row => formatTime(row.departed_at), align: 'center' },
  { name: 'returned_at', label: 'Ora Rientro', field: row => row.returned_at ? formatTime(row.returned_at) : '-', align: 'center' },
  { name: 'duration', label: 'Durata', align: 'center' },
  { name: 'notes', label: 'Note', field: 'notes', align: 'left' }
]

function getPassTypeLabel(type) {
  switch (type) {
    case 'bathroom': return t('securityPasses.typeBathroom') || 'Bagno'
    case 'infirmary': return t('securityPasses.typeInfirmary') || 'Infermeria'
    case 'office': return t('securityPasses.typeOffice') || 'Presidenza/Segreteria'
    case 'staff': return t('securityPasses.typeStaff') || 'Personale ATA'
    default: return type
  }
}

function getPassTypeColor(type) {
  switch (type) {
    case 'bathroom': return 'primary'
    case 'infirmary': return 'negative'
    case 'office': return 'warning'
    case 'staff': return 'secondary'
    default: return 'slate-500'
  }
}

function formatTime(isoStr) {
  if (!isoStr) return '-'
  const d = new Date(isoStr)
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function getElapsedMinutes(departedAt) {
  if (!departedAt) return 0
  const dep = new Date(departedAt).getTime()
  const diffMs = nowTime.value - dep
  return Math.max(0, Math.floor(diffMs / 60000))
}

function calculateDuration(departedAt, returnedAt) {
  if (!departedAt || !returnedAt) return 0
  const start = new Date(departedAt).getTime()
  const end = new Date(returnedAt).getTime()
  return Math.max(1, Math.round((end - start) / 60000))
}

async function loadClasses() {
  try {
    const res = await api.get('/classes')
    classesList.value = (res.data?.classes || res.data || []).map(c => ({
      id: c.id,
      name: c.name || `${c.year || ''} ${c.section || ''}`.trim()
    }))
    if (classesList.value.length > 0 && !selectedClassId.value) {
      selectedClassId.value = classesList.value[0].id
      await onClassChanged(selectedClassId.value)
    }
  } catch (err) {
    console.error('Errore caricamento classi:', err)
  }
}

async function onClassChanged(classId) {
  if (!classId) return
  await Promise.all([loadStudents(classId), loadPasses()])
}

async function loadStudents(classId) {
  try {
    const res = await api.get(`/classes/${classId}/students`)
    const raw = res.data?.students || res.data || []
    studentsInClass.value = raw.map(s => ({
      id: s.id,
      name: `${s.last_name || ''} ${s.first_name || ''}`.trim() || s.name || s.username
    }))
  } catch (err) {
    console.error('Errore caricamento studenti classe:', err)
  }
}

async function loadPasses() {
  if (!selectedClassId.value) return
  try {
    const res = await api.get('/security-passes/hall-passes/by-class', {
      params: { class_id: selectedClassId.value, limit: 50 }
    })
    const passes = res.data?.passes || []
    activePasses.value = passes.filter(p => !p.returned_at)
    historyPasses.value = passes
  } catch (err) {
    console.error('Errore caricamento pass:', err)
  }
}

function openIssueDialog() {
  newPass.value = {
    student_id: '',
    pass_type: 'bathroom',
    notes: ''
  }
  showIssueDialog.value = true
}

async function issuePass() {
  if (!newPass.value.student_id) return
  try {
    await api.post('/security-passes/hall-passes', {
      student_id: newPass.value.student_id,
      class_id: selectedClassId.value,
      pass_type: newPass.value.pass_type,
      notes: newPass.value.notes
    })
    $q.notify({ type: 'positive', message: 'Pass di uscita autorizzato' })
    showIssueDialog.value = false
    await loadPasses()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore emissione pass' })
  }
}

async function markReturned(passId) {
  try {
    await api.post(`/security-passes/hall-passes/${passId}/return`)
    $q.notify({ type: 'positive', message: 'Studente rientrato in classe' })
    await loadPasses()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante la registrazione del rientro' })
  }
}

onMounted(async () => {
  await loadClasses()
  timerInterval = setInterval(() => {
    nowTime.value = Date.now()
  }, 10000)
})

onBeforeUnmount(() => {
  if (timerInterval) clearInterval(timerInterval)
})
</script>
