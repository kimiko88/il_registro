<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('competenciesDM14.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('competenciesDM14.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-badge
          v-if="currentCertId"
          :color="getStatusBadgeColor(currentStatus)"
          class="q-pa-xs text-weight-bold text-caption"
        >
          Stato: {{ currentStatus.toUpperCase() }}
        </q-badge>
        <q-btn
          color="primary"
          icon="save"
          label="Salva Certificazione"
          class="rounded-lg shadow-sm"
          :loading="saving"
          :disable="!selectedStudentId || !selectedClassId"
          @click="saveCertification"
        />
        <q-btn
          v-if="currentCertId"
          color="positive"
          icon="verified"
          :label="$t('competenciesDM14.signByPrincipal')"
          class="rounded-lg shadow-sm"
          :loading="signing"
          @click="signCertification"
        />
      </div>
    </div>

    <!-- SELEZIONE CLASSE & STUDENTE REALI DAL DATABASE -->
    <q-card class="glass-card shadow-soft border-slate-100 q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-md items-center">
          <!-- Selezione Classe -->
          <div class="col-12 col-md-3">
            <q-select
              v-model="selectedClassId"
              :options="classOptions"
              option-value="id"
              option-label="label"
              emit-value
              map-options
              outlined
              dense
              label="Classe (Database)"
              :loading="loadingClasses"
              @update:model-value="onClassChanged"
            >
              <template #no-option>
                <q-item>
                  <q-item-section class="text-grey">Nessuna classe trovata</q-item-section>
                </q-item>
              </template>
            </q-select>
          </div>

          <!-- Selezione Studente -->
          <div class="col-12 col-md-3">
            <q-select
              v-model="selectedStudentId"
              :options="studentOptions"
              option-value="id"
              option-label="full_name"
              emit-value
              map-options
              outlined
              dense
              label="Alunno / Alunna (Database)"
              :loading="loadingStudents"
              :disable="!selectedClassId || studentOptions.length === 0"
              @update:model-value="onStudentChanged"
            >
              <template #no-option>
                <q-item>
                  <q-item-section class="text-grey">
                    {{ selectedClassId ? 'Nessuno studente in classe' : 'Seleziona prima una classe' }}
                  </q-item-section>
                </q-item>
              </template>
            </q-select>
          </div>

          <!-- Ordine di Scuola -->
          <div class="col-12 col-md-3">
            <q-select
              v-model="selectedTier"
              :options="[
                { label: $t('competenciesDM14.tierPrimaria'), value: 'primaria' },
                { label: $t('competenciesDM14.tierPrimoCiclo'), value: 'primo_ciclo' }
              ]"
              emit-value
              map-options
              outlined
              dense
              label="Ordine di Scuola (D.M. 14/2024)"
              @update:model-value="onTierChanged"
            />
          </div>

          <!-- Data Delibera Consiglio -->
          <div class="col-12 col-md-3">
            <q-input
              v-model="councilDate"
              type="date"
              outlined
              dense
              :label="$t('competenciesDM14.deliberationDate')"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- TABELLA COMPETENZE EUROPEE & DIGCOMP 2.2 -->
    <q-card class="glass-card shadow-soft border-slate-100 q-mb-lg">
      <q-card-section>
        <div class="row items-center justify-between q-mb-md">
          <div class="text-h6 text-weight-bold row items-center q-gutter-x-sm">
            <q-icon name="stars" color="primary" />
            <span>{{ $t('competenciesDM14.rubricTitle') }}</span>
          </div>

          <div v-if="selectedStudent" class="text-subtitle2 text-primary text-weight-bold">
            Compilazione per: {{ selectedStudent.full_name }}
          </div>
        </div>

        <div v-if="!selectedStudentId" class="text-center q-pa-lg text-slate-400">
          <q-icon name="person_search" size="48px" class="q-mb-sm opacity-50" />
          <div class="text-subtitle1">Seleziona una classe e uno studente per visualizzare e compilare le competenze dal database.</div>
        </div>

        <div v-else>
          <q-list separator>
            <q-item v-for="(comp, idx) in competenceList" :key="idx" class="q-py-md">
              <q-item-section>
                <q-item-label class="text-weight-bold text-slate-800 text-subtitle1">
                  {{ comp.name }}
                </q-item-label>
                <q-item-label caption class="text-slate-500">
                  {{ comp.description }}
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-btn-toggle
                  v-model="evaluations[comp.key]"
                  toggle-color="primary"
                  color="grey-2"
                  text-color="grey-8"
                  flat
                  dense
                  :options="[
                    { label: $t('competenciesDM14.levelIniziale'), value: 'iniziale' },
                    { label: $t('competenciesDM14.levelBase'), value: 'base' },
                    { label: $t('competenciesDM14.levelIntermedio'), value: 'intermedio' },
                    { label: $t('competenciesDM14.levelAvanzato'), value: 'avanzato' }
                  ]"
                />
              </q-item-section>
            </q-item>
          </q-list>

          <div class="q-mt-lg">
            <div class="text-subtitle2 text-weight-bold q-mb-xs">{{ $t('competenciesDM14.considerations') }}</div>
            <q-input
              v-model="finalConsiderations"
              type="textarea"
              outlined
              rows="3"
              placeholder="Inserisci la sintesi del consiglio di classe sul percorso dello studente..."
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- ELENCO CERTIFICAZIONI SALVATE NEL DATABASE PER LA CLASSE -->
    <q-card class="glass-card shadow-soft border-slate-100 q-mb-xl">
      <q-card-section>
        <div class="row items-center justify-between q-mb-md">
          <div class="text-h6 text-weight-bold row items-center q-gutter-x-sm">
            <q-icon name="list_alt" color="secondary" />
            <span>Certificazioni della Classe Registrate nel Database</span>
          </div>
          <q-btn flat round icon="refresh" color="secondary" :loading="loadingCertifications" @click="loadClassCertifications" />
        </div>

        <q-table
          :rows="classCertifications"
          :columns="certColumns"
          row-key="id"
          flat
          bordered
          :loading="loadingCertifications"
          no-data-label="Nessuna certificazione presente nel database per questa classe"
        >
          <template #body-cell-status="props">
            <q-td :props="props">
              <q-badge :color="getStatusBadgeColor(props.value)" class="q-pa-xs text-weight-bold">
                {{ (props.value || 'bozza').toUpperCase() }}
              </q-badge>
            </q-td>
          </template>

          <template #body-cell-actions="props">
            <q-td :props="props">
              <q-btn
                flat
                dense
                round
                color="primary"
                icon="edit"
                @click="selectCertRow(props.row)"
              >
                <q-tooltip>Modifica / Compila</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()

const selectedClassId = ref(null)
const selectedStudentId = ref(null)
const selectedTier = ref('primo_ciclo')
const councilDate = ref(new Date().toISOString().substring(0, 10))
const finalConsiderations = ref('')
const currentCertId = ref(null)
const currentStatus = ref('bozza')

const classOptions = ref([])
const studentOptions = ref([])
const classCertifications = ref([])

const loadingClasses = ref(false)
const loadingStudents = ref(false)
const loadingCertifications = ref(false)
const saving = ref(false)
const signing = ref(false)

const competenceList = [
  { key: 'lingua_madre', name: '1. Comunicazione nella madrelingua o lingua di scolarizzazione', description: 'Capacità di esprimere e interpretare concetti, pensieri e fatti in forma sia orale che scritta.' },
  { key: 'lingue_straniere', name: '2. Comunicazione nelle lingue straniere', description: 'Comprensione e produzione in almeno due lingue comunitarie (inglese e seconda lingua).' },
  { key: 'stem', name: '3. Competenze matematiche, scientifiche e tecnologiche (STEM)', description: 'Abilità di sviluppare e applicare il pensiero matematico e il metodo scientifico per spiegare il mondo.' },
  { key: 'digcomp', name: '4. Competenze digitali (Quadro DigComp 2.2)', description: 'Uso sicuro, critico e responsabile delle tecnologie digitali e dell\'intelligenza artificiale.' },
  { key: 'imparare_imparare', name: '5. Imparare a imparare (Competenza personale e sociale)', description: 'Capacità di perseverare nell\'apprendimento e organizzare il proprio studio autonomo.' },
  { key: 'civica_cittadinanza', name: '6. Competenze sociali e civiche (Educazione Civica)', description: 'Partecipazione attiva e democratica, rispetto della Costituzione e sostenibilità.' },
  { key: 'imprenditorialita', name: '7. Spirito di iniziativa e imprenditorialità', description: 'Capacità di tradurre idee in azioni attraverso la creatività, l\'innovazione e l\'assunzione di rischi calcolati.' },
  { key: 'consapevolezza_culturale', name: '8. Consapevolezza ed espressione culturale', description: 'Apprezzamento dell\'importanza dell\'espressione creativa di idee ed emozioni attraverso arti, musica e teatro.' }
]

// Evaluations state is initialized empty - NO MOCK VALUES
const evaluations = ref({})

const certColumns = [
  { name: 'student_name', label: 'Alunno / Alunna', field: 'student_full_name', align: 'left', sortable: true },
  { name: 'school_tier', label: 'Ordine', field: 'school_tier', align: 'center' },
  { name: 'deliberation_date', label: 'Data Delibera', field: 'council_deliberation_date', align: 'center' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

const selectedStudent = computed(() => {
  return studentOptions.value.find(s => s.id === selectedStudentId.value) || null
})

function getStatusBadgeColor(status) {
  switch (status) {
    case 'firmata_ds':
    case 'approvata':
      return 'positive'
    case 'bozza':
      return 'warning'
    default:
      return 'grey-6'
  }
}

async function loadClasses() {
  loadingClasses.value = true
  try {
    const res = await api.get('/classes')
    const raw = res.data?.classes || res.data || []
    classOptions.value = Array.isArray(raw)
      ? raw.map(c => ({
          id: c.id,
          label: c.name || c.section || `Classe ${c.id}`,
          school_level: c.school_level || c.order || 'primo_ciclo'
        }))
      : []

    if (classOptions.value.length > 0 && !selectedClassId.value) {
      selectedClassId.value = classOptions.value[0].id
      await onClassChanged(selectedClassId.value)
    }
  } catch (err) {
    console.error('Error loading classes:', err)
  } finally {
    loadingClasses.value = false
  }
}

async function onClassChanged(classId) {
  selectedStudentId.value = null
  studentOptions.value = []
  evaluations.value = {}
  finalConsiderations.value = ''
  currentCertId.value = null

  if (!classId) return

  const cls = classOptions.value.find(c => c.id === classId)
  if (cls && cls.school_level === 'primaria') {
    selectedTier.value = 'primaria'
  } else {
    selectedTier.value = 'primo_ciclo'
  }

  await Promise.all([
    loadStudents(classId),
    loadClassCertifications()
  ])
}

async function loadStudents(classId) {
  loadingStudents.value = true
  try {
    const res = await api.get(`/classes/${classId}/students`)
    const raw = res.data?.students || res.data || []
    studentOptions.value = Array.isArray(raw)
      ? raw.map(s => ({
          id: s.id || s.user_id,
          full_name: s.first_name && s.last_name ? `${s.last_name} ${s.first_name}` : (s.name || s.email || s.id)
        }))
      : []

    if (studentOptions.value.length > 0) {
      selectedStudentId.value = studentOptions.value[0].id
      await loadStudentCertification(selectedStudentId.value)
    }
  } catch (err) {
    console.error('Error loading class students:', err)
    studentOptions.value = []
  } finally {
    loadingStudents.value = false
  }
}

async function onStudentChanged(studentId) {
  if (!studentId) {
    evaluations.value = {}
    finalConsiderations.value = ''
    currentCertId.value = null
    currentStatus.value = 'bozza'
    return
  }
  await loadStudentCertification(studentId)
}

function onTierChanged() {
  if (selectedStudentId.value) {
    loadStudentCertification(selectedStudentId.value)
  }
}

async function loadStudentCertification(studentId) {
  try {
    const res = await api.get('/competencies-dm14/certifications/by-student', {
      params: {
        student_id: studentId,
        school_year: '2025/2026',
        school_tier: selectedTier.value
      }
    })

    if (res.data && res.data.id) {
      const cert = res.data
      currentCertId.value = cert.id
      currentStatus.value = cert.status || 'bozza'
      finalConsiderations.value = cert.final_considerations || ''
      if (cert.council_deliberation_date) {
        councilDate.value = cert.council_deliberation_date
      }

      try {
        evaluations.value = typeof cert.rubric_evaluations_json === 'string'
          ? JSON.parse(cert.rubric_evaluations_json)
          : (cert.rubric_evaluations_json || {})
      } catch {
        evaluations.value = {}
      }
    } else {
      resetCertificationForm()
    }
  } catch {
    // 404 means no certification stored in DB yet for this student
    resetCertificationForm()
  }
}

function resetCertificationForm() {
  currentCertId.value = null
  currentStatus.value = 'bozza'
  finalConsiderations.value = ''
  evaluations.value = {}
}

async function loadClassCertifications() {
  if (!selectedClassId.value) return
  loadingCertifications.value = true
  try {
    const res = await api.get('/competencies-dm14/certifications/by-class', {
      params: {
        class_id: selectedClassId.value,
        school_year: '2025/2026'
      }
    })
    const raw = res.data?.certifications || res.data || []
    classCertifications.value = Array.isArray(raw) ? raw : []
  } catch (err) {
    console.error('Error loading class certifications:', err)
    classCertifications.value = []
  } finally {
    loadingCertifications.value = false
  }
}

function selectCertRow(row) {
  if (row.student_id) {
    selectedStudentId.value = row.student_id
    if (row.school_tier) selectedTier.value = row.school_tier
    loadStudentCertification(row.student_id)
  }
}

async function saveCertification() {
  if (!selectedStudentId.value || !selectedClassId.value) {
    $q.notify({ type: 'warning', message: 'Seleziona una classe e uno studente prima di salvare.' })
    return
  }

  saving.value = true
  try {
    const payload = {
      student_id: selectedStudentId.value,
      class_id: selectedClassId.value,
      school_tier: selectedTier.value,
      school_year: '2025/2026',
      rubric_evaluations_json: JSON.stringify(evaluations.value),
      final_considerations: finalConsiderations.value,
      council_deliberation_date: councilDate.value,
      status: currentStatus.value || 'bozza'
    }

    const res = await api.post('/competencies-dm14/certifications', payload)
    if (res.data && res.data.id) {
      currentCertId.value = res.data.id
    }
    $q.notify({ type: 'positive', message: 'Certificazione salvata nel database con successo' })
    await loadClassCertifications()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore durante il salvataggio nel database' })
  } finally {
    saving.value = false
  }
}

async function signCertification() {
  if (!currentCertId.value) return
  signing.value = true
  try {
    await api.post(`/competencies-dm14/certifications/${currentCertId.value}/sign`)
    currentStatus.value = 'firmata_ds'
    $q.notify({ type: 'positive', message: 'Certificazione firmata digitalmente con FEQ del Dirigente' })
    await loadClassCertifications()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore durante la firma digitale' })
  } finally {
    signing.value = false
  }
}

onMounted(() => {
  loadClasses()
})
</script>
