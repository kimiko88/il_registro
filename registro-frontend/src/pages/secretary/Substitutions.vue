<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg sticky-header">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="swap_horiz" color="primary" class="q-mr-sm" />
          Gestione Sostituzioni Docenti
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Pianificazione e assegnazione delle ore di supplenza per i docenti assenti
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-input
          v-model="selectedDate"
          type="date"
          dense outlined
          label="Data Sostituzioni"
          class="bg-white"
          style="min-width: 180px"
          @update:model-value="fetchSubstitutions"
        />
        <q-btn
          color="primary"
          icon="add"
          label="Nuova Sostituzione"
          no-caps
          class="rounded-lg shadow-sm"
          @click="openCreateDialog"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="fetchSubstitutions" />
      </div>
    </div>

    <!-- Summary Cards -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Da Assegnare</div>
              <div class="text-h4 text-weight-bold text-amber-7 q-mt-xs">{{ unassignedCount }}</div>
              <div class="text-caption text-slate-400">Richieste in attesa di sostituto</div>
            </div>
            <q-avatar color="amber-1" text-color="amber-9" icon="pending" size="48px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Assegnate</div>
              <div class="text-h4 text-weight-bold text-primary q-mt-xs">{{ assignedCount }}</div>
              <div class="text-caption text-slate-400">Supplenze coperte</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="how_to_reg" size="48px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Totale Giornaliero</div>
              <div class="text-h4 text-weight-bold text-slate-800 q-mt-xs">{{ substitutions.length }}</div>
              <div class="text-caption text-slate-400">Supplenze in totale per la data</div>
            </div>
            <q-avatar color="slate-1" text-color="slate-7" icon="event_available" size="48px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Table of Substitutions -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
      <q-table
        :rows="substitutions"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        class="bg-transparent"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template v-slot:body-cell-hour="props">
          <q-td :props="props">
            <q-badge color="indigo-7" class="q-px-sm q-py-xs text-weight-bold">
              {{ props.value }}ª ora
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-chip
              dense
              :color="statusColor(props.value)"
              text-color="white"
              class="text-weight-bold"
            >
              {{ statusLabel(props.value) }}
            </q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" auto-width>
            <div class="row q-gutter-xs">
              <q-btn
                color="primary"
                size="sm"
                icon="person_search"
                label="Assegna Docente"
                no-caps
                class="rounded-lg"
                @click="openAssignDialog(props.row)"
              />
            </div>
          </q-td>
        </template>

        <template v-slot:no-data>
          <div class="full-width column flex-center q-pa-xl text-slate-400">
            <q-icon name="swap_horiz" size="64px" class="opacity-30" />
            <div class="text-h6 q-mt-md">Nessuna sostituzione per la data selezionata</div>
            <p class="text-caption">Clicca su "Nuova Sostituzione" per creare una richiesta di supplenza.</p>
          </div>
        </template>
      </q-table>
    </q-card>

    <!-- Create Substitution Dialog -->
    <q-dialog v-model="showCreateDialog" persistent>
      <q-card style="min-width: 500px" class="rounded-xl overflow-hidden shadow-24 bg-white">
        <q-card-section class="bg-primary text-white q-pa-lg row items-center justify-between">
          <div class="text-h6 text-weight-bold">
            <q-icon name="add_circle" class="q-mr-xs" />
            Nuova Richiesta Sostituzione
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg space-y-4">
          <q-form @submit="createSubstitution" class="q-gutter-y-md">
            <q-select
              v-model="createForm.absent_teacher_id"
              :options="teacherOptions"
              label="Docente Assente *"
              outlined dense emit-value map-options
              :rules="[val => !!val || 'Selezionare il docente assente']"
            />

            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <q-input
                  v-model="createForm.date"
                  type="date"
                  label="Data *"
                  outlined dense
                  :rules="[val => !!val || 'Obbligatorio']"
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="createForm.hour_index"
                  :options="[1, 2, 3, 4, 5, 6, 7, 8]"
                  label="Ora di Lezione *"
                  outlined dense
                  :rules="[val => !!val || 'Obbligatorio']"
                />
              </div>
            </div>

            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="createForm.class_id"
                  :options="classOptions"
                  label="Classe *"
                  outlined dense emit-value map-options
                  :rules="[val => !!val || 'Obbligatorio']"
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="createForm.subject_id"
                  :options="subjectOptions"
                  label="Materia *"
                  outlined dense emit-value map-options
                  :rules="[val => !!val || 'Obbligatorio']"
                />
              </div>
            </div>

            <q-input
              v-model="createForm.notes"
              label="Note / Motivo Assenza (opzionale)"
              outlined dense
              placeholder="Es. Malattia, Permesso, Corso di formazione"
            />

            <div class="row justify-end q-mt-lg q-gutter-sm">
              <q-btn flat label="Annulla" v-close-popup no-caps />
              <q-btn type="submit" label="Crea Sostituzione" color="primary" class="rounded-lg q-px-lg" no-caps :loading="saving" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Assign Substitute Dialog -->
    <q-dialog v-model="showAssignDialog" persistent>
      <q-card style="min-width: 550px" class="rounded-xl overflow-hidden shadow-24 bg-white">
        <q-card-section class="bg-primary text-white q-pa-lg row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold">Assegna Docente Sostituto</div>
            <div class="text-caption" v-if="selectedSub">
              Ora {{ selectedSub.hour }}ª · Classe {{ selectedSub.class_name || selectedSub.class_id }} · Data {{ selectedSub.date ? selectedSub.date.substring(0, 10) : '' }}
            </div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg">
          <div class="text-subtitle2 text-weight-bold text-slate-800 q-mb-sm">
            Docenti Consigliati / Disponibili:
          </div>

          <div v-if="loadingRecs" class="text-center q-pa-md">
            <q-spinner-dots color="primary" size="30px" />
          </div>

          <q-list v-else-if="recommendedTeachers.length > 0" separator class="border rounded-lg bg-slate-50 q-mb-md">
            <q-item
              v-for="rec in recommendedTeachers"
              :key="rec.teacher_id"
              tag="label"
              v-ripple
              class="cursor-pointer"
            >
              <q-item-section side>
                <q-radio v-model="selectedSubstituteId" :val="rec.teacher_id" />
              </q-item-section>

              <q-item-section>
                <q-item-label class="text-weight-bold">{{ rec.teacher_name }}</q-item-label>
                <q-item-label caption class="text-slate-600">
                  {{ rec.reason || rec.status || 'Disponibile per la classe' }}
                </q-item-label>
              </q-item-section>

              <q-item-section side v-if="rec.is_free">
                <q-badge color="positive" label="Ora Libera" />
              </q-item-section>
            </q-item>
          </q-list>

          <div class="q-mb-md">
            <q-select
              v-model="selectedSubstituteId"
              :options="teacherOptions"
              label="Oppure seleziona un altro docente"
              outlined dense emit-value map-options
            />
          </div>

          <div class="row justify-end q-gutter-sm">
            <q-btn flat label="Annulla" v-close-popup no-caps />
            <q-btn
              label="Conferma Assegnazione"
              color="primary"
              class="rounded-lg q-px-lg"
              no-caps
              :disabled="!selectedSubstituteId"
              :loading="saving"
              @click="assignSubstitute"
            />
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useQuasar } from 'quasar'
import { substitutionService } from 'src/services/substitutionService'
import adminService from 'src/services/adminService'
import api from 'src/services/api'

const $q = useQuasar()

const selectedDate = ref(new Date().toISOString().substring(0, 10))
const loading = ref(false)
const saving = ref(false)
const loadingRecs = ref(false)

const substitutions = ref([])
const teachers = ref([])
const classes = ref([])
const subjects = ref([])

const showCreateDialog = ref(false)
const showAssignDialog = ref(false)
const selectedSub = ref(null)
const selectedSubstituteId = ref(null)
const recommendedTeachers = ref([])

const createForm = reactive({
  absent_teacher_id: null,
  date: new Date().toISOString().substring(0, 10),
  hour_index: 1,
  class_id: null,
  subject_id: null,
  notes: ''
})

const columns = [
  { name: 'hour', label: 'Ora', field: row => row.hour || row.hour_index || 1, align: 'center', sortable: true },
  { name: 'date', label: 'Data', field: row => row.date ? row.date.substring(0, 10) : '', align: 'left' },
  { name: 'class', label: 'Classe', field: row => row.class_name || row.class_id, align: 'left', sortable: true },
  { name: 'subject', label: 'Materia', field: row => row.subject_name || row.subject_id, align: 'left' },
  { name: 'absent_teacher', label: 'Docente Assente', field: row => row.absent_teacher_name || row.absent_teacher_id, align: 'left' },
  { name: 'substitute_teacher', label: 'Sostituto Assegnato', field: row => row.substitute_teacher_name || row.substitute_teacher_id || 'Nessuno', align: 'left' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center', sortable: true },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

const unassignedCount = computed(() => substitutions.value.filter(s => s.status === 'pending' || !s.substitute_teacher_id).length)
const assignedCount = computed(() => substitutions.value.filter(s => s.status === 'assigned' || s.status === 'confirmed' || s.status === 'completed').length)

const teacherOptions = computed(() => teachers.value.map(t => ({ label: `${t.last_name || ''} ${t.first_name || t.name || ''}`, value: t.id })))
const classOptions = computed(() => classes.value.map(c => ({ label: c.label || c.name || `Classe ${c.id.substring(0,6)}`, value: c.id })))
const subjectOptions = computed(() => subjects.value.map(s => ({ label: s.name, value: s.id })))

onMounted(async () => {
  await Promise.all([
    fetchSubstitutions(),
    fetchTeachers(),
    fetchClasses(),
    fetchSubjects()
  ])
})

async function fetchSubstitutions() {
  loading.value = true
  try {
    const res = await substitutionService.listBySchool({ date: selectedDate.value })
    substitutions.value = res.data || []
  } catch (err) {
    console.error(err)
    $q.notify({ type: 'negative', message: 'Errore durante il recupero delle sostituzioni' })
  } finally {
    loading.value = false
  }
}

async function fetchTeachers() {
  try {
    const res = await api.get('/users', { params: { role: 'teacher' } })
    teachers.value = res.data?.users || res.data || []
  } catch (err) {
    console.error(err)
  }
}

async function fetchClasses() {
  try {
    const res = await adminService.getClasses()
    classes.value = res.data || []
  } catch (err) {
    console.error(err)
  }
}

async function fetchSubjects() {
  try {
    const res = await api.get('/subjects')
    subjects.value = res.data || []
  } catch (err) {
    console.error(err)
  }
}

function openCreateDialog() {
  Object.assign(createForm, {
    absent_teacher_id: teacherOptions.value[0]?.value || null,
    date: selectedDate.value,
    hour_index: 1,
    class_id: classOptions.value[0]?.value || null,
    subject_id: subjectOptions.value[0]?.value || null,
    notes: ''
  })
  showCreateDialog.value = true
}

async function createSubstitution() {
  saving.value = true
  try {
    await substitutionService.create({
      absent_teacher_id: createForm.absent_teacher_id,
      date: createForm.date,
      hour: parseInt(createForm.hour_index),
      class_id: createForm.class_id,
      subject_id: createForm.subject_id,
      notes: createForm.notes
    })
    $q.notify({ type: 'positive', message: 'Richiesta di sostituzione creata' })
    showCreateDialog.value = false
    await fetchSubstitutions()
  } catch (err) {
    console.error(err)
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore creazione sostituzione' })
  } finally {
    saving.value = false
  }
}

async function openAssignDialog(row) {
  selectedSub.value = row
  selectedSubstituteId.value = row.substitute_teacher_id || null
  showAssignDialog.value = true
  loadingRecs.value = true
  try {
    const dateStr = row.date ? row.date.substring(0, 10) : selectedDate.value
    const res = await substitutionService.recommendSubstitutes({
      class_id: row.class_id,
      subject_id: row.subject_id,
      date: dateStr,
      hour: row.hour || row.hour_index || 1
    })
    recommendedTeachers.value = res.data || []
  } catch (err) {
    console.error(err)
    recommendedTeachers.value = []
  } finally {
    loadingRecs.value = false
  }
}

async function assignSubstitute() {
  if (!selectedSub.value || !selectedSubstituteId.value) return
  saving.value = true
  try {
    await substitutionService.assign(selectedSub.value.id, {
      substitute_teacher_id: selectedSubstituteId.value
    })
    $q.notify({ type: 'positive', message: 'Docente sostituto assegnato con successo' })
    showAssignDialog.value = false
    await fetchSubstitutions()
  } catch (err) {
    console.error(err)
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore durante l\'assegnazione del sostituto' })
  } finally {
    saving.value = false
  }
}

function statusLabel(s) {
  return {
    pending: 'Da Assegnare',
    assigned: 'Assegnata',
    confirmed: 'Confermata',
    completed: 'Completata'
  }[s] || 'In attesa'
}

function statusColor(s) {
  return {
    pending: 'amber-8',
    assigned: 'info',
    confirmed: 'positive',
    completed: 'grey-7'
  }[s] || 'primary'
}
</script>

<style scoped>
.sticky-header {
  position: sticky;
  top: 0;
  z-index: 10;
}
</style>
