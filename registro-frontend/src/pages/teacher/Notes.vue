<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="assignment_late" color="primary" class="q-mr-sm" />
          Note & Annotazioni Disciplinari
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Gestione delle annotazioni scolastiche e note riservate al coordinatore
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn
          color="primary"
          unelevated
          icon="add"
          label="Nuova Nota"
          class="rounded-lg q-px-md"
          no-caps
          @click="openDialog"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="notesStore.loading" @click="loadNotes" />
      </div>
    </div>

    <!-- Filters Section Card -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft q-pa-md q-mb-lg">
      <div class="row q-col-gutter-md items-center">
        <div class="col-12 col-sm-4">
          <q-select
            v-model="selectedClassId"
            :options="classOptions"
            label="Seleziona Classe *"
            outlined dense
            emit-value map-options
            @update:model-value="onClassChange"
          />
        </div>
        <div class="col-12 col-sm-4">
          <q-select
            v-model="selectedStudentId"
            :options="studentOptions"
            label="Filtra per Studente"
            outlined dense
            clearable
            emit-value map-options
            @update:model-value="loadNotes"
          />
        </div>
        <div class="col-12 col-sm-4 text-right">
          <q-badge color="indigo-1" text-color="indigo-8" class="q-px-md q-py-xs text-weight-bold">
            Totale Note: {{ filteredNotes.length }}
          </q-badge>
        </div>
      </div>
    </q-card>

    <!-- Notes List Section -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden min-h-400">
      <div v-if="notesStore.loading" class="text-center q-pa-xl">
        <q-spinner-dots color="primary" size="40px" />
      </div>

      <div v-else-if="filteredNotes.length === 0" class="text-center q-pa-xl text-slate-400">
        <q-icon name="assignment_turned_in" size="64px" class="q-mb-md opacity-40" />
        <div class="text-h6">Nessuna nota trovata</div>
        <div class="text-caption">Seleziona una classe o aggiungi una nuova annotazione disciplinare.</div>
      </div>

      <q-list v-else separator class="rounded-lg">
        <q-item v-for="n in filteredNotes" :key="n.id" class="q-py-md">
          <q-item-section avatar>
            <q-avatar :color="n.is_reserved ? 'red-1' : 'grey-2'" :text-color="n.is_reserved ? 'negative' : 'grey-8'" size="44px">
              <q-icon :name="n.is_reserved ? 'lock' : 'assignment'" />
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <div class="row items-center justify-between q-mb-xs">
              <div class="row items-center q-gutter-xs">
                <span class="text-weight-bold text-slate-800 text-subtitle1">{{ n.teacher_name || 'Docente' }}</span>
                <span class="text-caption text-slate-400">· {{ formatDate(n.date) }}</span>
              </div>
              <div class="row items-center q-gutter-xs">
                <q-chip
                  v-if="n.type === 'disciplinary'"
                  size="xs"
                  :color="n.is_approved ? 'positive' : 'warning'"
                  text-color="white"
                  class="text-weight-bold"
                  :icon="n.is_approved ? 'check_circle' : 'hourglass_empty'"
                >
                  {{ n.is_approved ? 'Approvata Admin' : 'In attesa di approvazione Admin' }}
                </q-chip>
                <q-chip
                  size="xs"
                  :color="n.is_viewed_by_parent ? 'info' : 'grey-5'"
                  text-color="white"
                  class="text-weight-bold"
                  :icon="n.is_viewed_by_parent ? 'visibility' : 'visibility_off'"
                >
                  {{ n.is_viewed_by_parent ? 'Letta dal genitore' : 'Non ancora letta' }}
                </q-chip>
                <q-chip
                  size="xs"
                  :color="n.is_reserved ? 'negative' : 'grey-7'"
                  text-color="white"
                  class="text-weight-bold"
                  :icon="n.is_reserved ? 'lock' : 'public'"
                >
                  {{ n.is_reserved ? 'Riservata Coordinatore' : 'Standard' }}
                </q-chip>
                <q-chip size="xs" :color="getTypeColor(n.type)" text-color="white" class="text-weight-bold uppercase">
                  {{ n.type }}
                </q-chip>
              </div>
            </div>

            <div class="text-body2 text-slate-700 q-mt-xs">{{ n.note }}</div>

            <div class="text-caption text-slate-400 q-mt-sm" v-if="n.subject_name">
              Materia: <span class="text-weight-medium text-slate-600">{{ n.subject_name }}</span>
            </div>
          </q-item-section>

          <q-item-section side>
            <q-btn flat round icon="delete" color="negative" size="sm" @click="confirmDelete(n.id)" />
          </q-item-section>
        </q-item>
      </q-list>
    </q-card>

    <!-- Create Note Dialog -->
    <q-dialog v-model="dialogVisible">
      <q-card style="min-width: 450px; max-width: 550px" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div class="text-h6 text-weight-bold">
            <q-icon name="edit_note" class="q-mr-xs" />
            Nuova Nota Disciplinare / Annotazione
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md space-y-4">
          <!-- Class & Student Select -->
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-select
                v-model="form.class_id"
                :options="classOptions"
                label="Classe *"
                outlined dense
                emit-value map-options
                @update:model-value="onDialogClassChange"
              />
            </div>
            <div class="col-12 col-sm-6">
              <q-select
                v-model="form.student_id"
                :options="dialogStudentOptions"
                label="Studente *"
                outlined dense
                emit-value map-options
              />
            </div>
          </div>

          <!-- Note Type & Date -->
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-select
                v-model="form.type"
                :options="typeOptions"
                label="Tipo Nota *"
                outlined dense
                emit-value map-options
              />
            </div>
            <div class="col-12 col-sm-6">
              <q-input
                v-model="form.date"
                type="date"
                label="Data *"
                outlined dense
              />
            </div>
          </div>

          <!-- Note Body Textarea -->
          <q-input
            v-model="form.note"
            label="Testo della Nota *"
            outlined
            type="textarea"
            rows="4"
            :rules="[val => !!val || 'Inserisci il testo della nota']"
          />

          <!-- Reserved Toggle -->
          <div class="bg-red-50 border border-red-200 q-pa-sm rounded-lg">
            <q-toggle
              v-model="form.is_reserved"
              label="Nota Riservata (visibile SOLO a coordinatore e dirigenza/admin)"
              color="negative"
            />
          </div>
        </q-card-section>

        <q-separator />

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Salva Nota" :loading="saving" @click="saveNote" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { useNotesStore } from '@/stores/notes'
import { useClassesStore } from '@/stores/classes'
import api from 'src/services/api'

const $q = useQuasar()
const notesStore = useNotesStore()
const classesStore = useClassesStore()

const selectedClassId = ref(null)
const selectedStudentId = ref(null)
const classStudents = ref([])
const dialogStudents = ref([])

const dialogVisible = ref(false)
const saving = ref(false)

const form = reactive({
  class_id: '',
  student_id: '',
  type: 'generic',
  note: '',
  date: new Date().toISOString().substring(0, 10),
  is_reserved: false,
  target_role: 'all'
})

const typeOptions = [
  { label: 'Generica', value: 'generic' },
  { label: 'Mancanza Compiti', value: 'homework' },
  { label: 'Comportamento', value: 'behavior' },
  { label: 'Nota Disciplinare Solenne', value: 'disciplinary' }
]

const classOptions = computed(() => {
  return classesStore.classes.map(c => ({
    label: c.label || c.displayName || c.name || `Classe ${c.id}`,
    value: c.id
  }))
})

const studentOptions = computed(() => {
  return classStudents.value.map(s => ({
    label: `${s.first_name || s.name} ${s.last_name || ''}`,
    value: s.id
  }))
})

const dialogStudentOptions = computed(() => {
  return dialogStudents.value.map(s => ({
    label: `${s.first_name || s.name} ${s.last_name || ''}`,
    value: s.id
  }))
})

const filteredNotes = computed(() => {
  return notesStore.notes.filter(n => {
    const matchesClass = !selectedClassId.value || n.class_id === selectedClassId.value
    const matchesStudent = !selectedStudentId.value || n.student_id === selectedStudentId.value
    return matchesClass && matchesStudent
  })
})

const formatDate = (d) => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY') : ''

function getTypeColor(type) {
  switch (type) {
    case 'generic': return 'grey-7'
    case 'homework': return 'amber-8'
    case 'behavior': return 'deep-orange'
    case 'disciplinary': return 'negative'
    default: return 'primary'
  }
}

onMounted(async () => {
  await classesStore.fetchAssignedClasses().catch(() => classesStore.fetchClasses())
  if (classesStore.classes.length > 0) {
    selectedClassId.value = classesStore.classes[0].id
    await onClassChange()
  } else {
    await loadNotes()
  }
})

async function onClassChange() {
  if (!selectedClassId.value) {
    classStudents.value = []
    selectedStudentId.value = null
    await loadNotes()
    return
  }
  try {
    const res = await api.get(`/users`, { params: { role: 'student', class_id: selectedClassId.value } })
    classStudents.value = res.data?.users || res.data || []
  } catch {
    classStudents.value = []
  }
  await loadNotes()
}

async function loadNotes() {
  const params = {}
  if (selectedClassId.value) params.class_id = selectedClassId.value
  if (selectedStudentId.value) params.student_id = selectedStudentId.value
  await notesStore.fetchNotes(params).catch(() => {})
}

async function onDialogClassChange() {
  if (!form.class_id) {
    dialogStudents.value = []
    return
  }
  try {
    const res = await api.get(`/users`, { params: { role: 'student', class_id: form.class_id } })
    dialogStudents.value = res.data?.users || res.data || []
    if (dialogStudents.value.length > 0) {
      form.student_id = dialogStudents.value[0].id
    }
  } catch {
    dialogStudents.value = []
  }
}

async function openDialog() {
  form.class_id = selectedClassId.value || classOptions.value[0]?.value || ''
  form.type = 'generic'
  form.note = ''
  form.date = new Date().toISOString().substring(0, 10)
  form.is_reserved = false
  form.target_role = 'all'
  form.student_id = ''
  dialogVisible.value = true
  await onDialogClassChange()
}

async function saveNote() {
  if (!form.note || !form.class_id || !form.student_id) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }
  saving.value = true
  try {
    await notesStore.addNote({
      student_id: form.student_id,
      class_id: form.class_id,
      type: form.type,
      note: form.note,
      date: form.date,
      is_reserved: form.is_reserved,
      target_role: form.is_reserved ? 'coordinator' : 'all'
    })
    $q.notify({ type: 'positive', message: 'Nota salvata con successo' })
    dialogVisible.value = false
    await loadNotes()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante il salvataggio' })
  } finally {
    saving.value = false
  }
}

async function confirmDelete(id) {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: 'Sei sicuro di voler eliminare questa nota?',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await notesStore.deleteNote(id)
      $q.notify({ type: 'positive', message: 'Nota eliminata' })
      await loadNotes()
    } catch {
      $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione' })
    }
  })
}
</script>

<style scoped>
.min-h-400 {
  min-height: 400px;
}
</style>
