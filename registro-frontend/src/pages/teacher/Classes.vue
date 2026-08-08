<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h4">Le Mie Classi</div>
      <q-select
        v-model="selectedClass"
        :options="classesStore.classes"
        option-label="label"
        label="Seleziona Classe"
        outlined
        dense
        options-dense
        bg-color="white"
        style="min-width: 280px"
        @update:model-value="selectClass"
      >
        <template v-slot:option="scope">
          <q-item v-bind="scope.itemProps">
            <q-item-section avatar>
              <q-avatar color="primary" text-color="white" size="28px">{{ scope.opt.name }}</q-avatar>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ scope.opt.label || scope.opt.displayName || scope.opt.name }}</q-item-label>
            </q-item-section>
            <q-item-section side v-if="scope.opt.isCoordinator">
              <q-badge color="amber-9" label="COORD" />
            </q-item-section>
          </q-item>
        </template>
      </q-select>
    </div>

    <div class="row q-col-gutter-lg">
      <!-- Class List Sidebar -->
      <div class="col-12 col-md-3">
        <q-list bordered class="bg-white rounded-borders">
          <q-item-label header class="text-weight-bold bg-grey-2">Elenco Classi</q-item-label>
          <q-item 
            v-for="cls in classesStore.classes" 
            :key="cls.id" 
            clickable 
            v-ripple
            :active="selectedClass?.id === cls.id"
            active-class="bg-blue-1 text-primary"
            @click="selectClass(cls)"
          >
            <q-item-section avatar>
              <q-avatar color="primary" text-color="white">{{ cls.name }}</q-avatar>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ getClassLabel(cls) }}</q-item-label>
              <q-item-label caption>{{ cls.students ?? cls.students_count ?? 0 }} Studenti</q-item-label>
            </q-item-section>
            <q-item-section side v-if="cls.isCoordinator">
              <q-icon name="star" color="orange"><q-tooltip>Coordinatore</q-tooltip></q-icon>
            </q-item-section>
          </q-item>
        </q-list>
      </div>

      <!-- Class Detail View -->
      <div class="col-12 col-md-9" v-if="selectedClass">
        <q-card>
          <q-toolbar class="bg-primary text-white">
            <q-toolbar-title>Classe {{ getClassLabel(selectedClass) }}</q-toolbar-title>
            <q-tabs v-model="tab" shrink stretch>
              <q-tab name="students" label="Studenti" />
              <q-tab name="notes" label="Note di Classe" />
              <q-tab name="grades" label="Riepilogo Voti" />
              <q-tab name="dashboard" label="Coordinatore" v-if="selectedClass.isCoordinator" />
            </q-tabs>
          </q-toolbar>

          <q-separator />

          <q-tab-panels v-model="tab" animated>
            <!-- Student List -->
            <q-tab-panel name="students">
              <div class="row items-center justify-between q-mb-md">
                <div class="text-h6">Elenco Studenti</div>
                <q-input dense outlined placeholder="Cerca studente..." v-model="search" rounded>
                  <template v-slot:append><q-icon name="search" /></template>
                </q-input>
              </div>
              <q-table :rows="filteredStudents" :columns="columns" flat bordered row-key="id" :loading="loadingStudents">
                <template v-slot:body-cell-actions="props">
                  <q-td :props="props">
                    <q-btn round flat dense icon="note_add" color="grey-7" @click="openNoteDialog(props.row)">
                      <q-tooltip>Aggiungi Nota</q-tooltip>
                    </q-btn>
                  </q-td>
                </template>
              </q-table>
            </q-tab-panel>

            <!-- Notes List -->
            <q-tab-panel name="notes">
              <div class="row items-center justify-between q-mb-md">
                <div class="text-h6">Note di Classe</div>
                <q-select
                  v-model="filterNoteType"
                  :options="[
                    {label:'Tutte', value:''}, 
                    {label:'Nota Generica', value:'generic'}, 
                    {label:'Richiamo Compiti', value:'homework'}, 
                    {label:'Nota Comportamentale', value:'behavior'}, 
                    {label:'Nota Disciplinare', value:'disciplinary'}
                  ]"
                  option-value="value"
                  option-label="label"
                  emit-value map-options
                  label="Filtra per Tipo"
                  dense outlined
                  style="min-width: 180px"
                />
              </div>

              <q-table :rows="filteredNotes" :columns="noteColumns" flat bordered row-key="id" :loading="loadingNotes">
                <template v-slot:body-cell-actions="props">
                  <q-td :props="props" class="q-gutter-xs">
                    <q-btn round flat dense icon="edit" color="primary" @click="editNote(props.row)">
                      <q-tooltip>Modifica Nota</q-tooltip>
                    </q-btn>
                    <q-btn round flat dense icon="delete" color="negative" @click="deleteNote(props.row.id)">
                      <q-tooltip>Elimina Nota</q-tooltip>
                    </q-btn>
                  </q-td>
                </template>
              </q-table>
            </q-tab-panel>

            <!-- Grades Summary -->
            <q-tab-panel name="grades">
              <div class="text-center text-grey text-h6 q-pa-xl">
                <q-icon name="analytics" size="64px" class="q-mb-sm" />
                <div>Riepilogo Voti e Media Classe</div>
                <q-btn label="Vedi Dettagli Voti" color="primary" flat class="q-mt-sm" to="/teacher/grades" />
              </div>
            </q-tab-panel>

            <!-- Coordinator Dashboard -->
            <q-tab-panel name="dashboard">
              <div class="row q-col-gutter-md">
                <div class="col-12 col-md-6">
                  <q-card bordered flat class="bg-red-1">
                    <q-card-section>
                      <div class="text-subtitle1 text-red-9 text-weight-bold">Studenti a Rischio</div>
                      <q-list dense>
                        <q-item><q-item-label>• Nessuna segnalazione critica</q-item-label></q-item>
                      </q-list>
                    </q-card-section>
                  </q-card>
                </div>
                <div class="col-12 col-md-6">
                  <q-card bordered flat class="bg-yellow-1">
                    <q-card-section>
                      <div class="text-subtitle1 text-orange-9 text-weight-bold">Scadenze Coordinamento</div>
                      <q-list dense>
                        <q-item><q-item-label>• Preparazione Consigli di Classe</q-item-label></q-item>
                      </q-list>
                    </q-card-section>
                  </q-card>
                </div>
              </div>
            </q-tab-panel>
          </q-tab-panels>
        </q-card>
      </div>

      <div class="col-12 col-md-9 text-center text-grey" v-else>
        <div class="q-mt-xl">
          <q-icon name="school" size="100px" />
          <div class="text-h5">Seleziona una classe per gestire</div>
        </div>
      </div>
    </div>

    <!-- Note Dialog -->
    <NoteDialog
      v-if="selectedClass" 
      v-model="showNoteDialog"
      :student="selectedStudentForNote"
      :class-id="String(selectedClass.id)" 
      :note-to-edit="selectedNoteForEdit"
      @saved="onNoteSaved"
    />
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import NoteDialog from 'src/components/Teacher/NoteDialog.vue'
import { useClassesStore } from '@/stores/classes'
import notesService from '@/services/notesService'
import api from '@/services/api'

import { useSchoolYearStore } from '@/stores/schoolYear'

const classesStore = useClassesStore()
const schoolYearStore = useSchoolYearStore()
const $q = useQuasar()

function getClassLabel(cls) {
  if (!cls) return ''
  if (cls.label) return cls.label
  if (cls.displayName) return cls.displayName
  let nameText = cls.name || ''
  if (cls.section && !nameText.endsWith(cls.section)) {
    nameText += cls.section
  }
  if (cls.articolazione) {
    nameText += ` - ${cls.articolazione}`
  }
  return nameText
}

const tab = ref('students')
const search = ref('')
const selectedClass = ref(null)
const students = ref([])
const loadingStudents = ref(false)

const showNoteDialog = ref(false)
const selectedStudentForNote = ref(null)
const selectedNoteForEdit = ref(null)

const notes = ref([])
const loadingNotes = ref(false)
const filterNoteType = ref('')

onMounted(async () => {
  await classesStore.fetchAssignedClasses(schoolYearStore.selectedSchoolYear)
  if (classesStore.classes.length > 0) {
    selectClass(classesStore.classes[0])
  }
})

watch(() => schoolYearStore.selectedSchoolYear, async (newYear) => {
  await classesStore.fetchAssignedClasses(newYear)
  if (classesStore.classes.length > 0) {
    selectClass(classesStore.classes[0])
  } else {
    selectedClass.value = null
  }
})

const columns = [
  { name: 'name', label: 'Nome', field: row => `${row.last_name} ${row.first_name}`, align: 'left', sortable: true },
  { name: 'email', label: 'Email', field: 'email', align: 'left' },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

const noteColumns = [
  { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'student', label: 'Studente', field: row => getStudentName(row.student_id), align: 'left', sortable: true },
  { name: 'type', label: 'Tipo', field: row => formatNoteType(row.type), align: 'left' },
  { name: 'note', label: 'Contenuto', field: 'note', align: 'left' },
  { name: 'docente', label: 'Docente', field: 'teacher_name', align: 'left' },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

const filteredStudents = computed(() => {
  if (!search.value) return students.value
  const lower = search.value.toLowerCase()
  return students.value.filter(s => 
    (s.last_name && s.last_name.toLowerCase().includes(lower)) || 
    (s.first_name && s.first_name.toLowerCase().includes(lower))
  )
})

const getStudentName = (studentId) => {
  const s = students.value.find(st => st.id === studentId)
  return s ? `${s.last_name} ${s.first_name}` : `Studente`
}

const formatNoteType = (type) => {
  const map = {
    generic: 'Generica',
    homework: 'Compiti',
    behavior: 'Comportamentale',
    disciplinary: 'Disciplinare'
  }
  return map[type] || type
}

const selectClass = async (cls) => {
  selectedClass.value = cls
  tab.value = 'students'
  await fetchStudents(cls.id)
}

const fetchStudents = async (classId) => {
  loadingStudents.value = true
  try {
    const res = await api.get('/users', {
      params: { class_id: classId, role: 'student', page_size: 100 }
    })
    students.value = res.data.users || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore caricamento studenti' })
  } finally {
    loadingStudents.value = false
  }
}

const fetchNotes = async () => {
  if (!selectedClass.value) return
  loadingNotes.value = true
  try {
    const res = await notesService.getNotes({ class_id: selectedClass.value.id })
    notes.value = res.data || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore caricamento note' })
  } finally {
    loadingNotes.value = false
  }
}

watch(tab, (newTab) => {
  if (newTab === 'notes') {
    fetchNotes()
  }
})

const filteredNotes = computed(() => {
  if (!filterNoteType.value) return notes.value
  return notes.value.filter(n => n.type === filterNoteType.value)
})

const openNoteDialog = (student) => {
  selectedNoteForEdit.value = null
  selectedStudentForNote.value = student
  showNoteDialog.value = true
}

const editNote = (note) => {
  selectedNoteForEdit.value = note
  selectedStudentForNote.value = students.value.find(s => s.id === note.student_id) || { id: note.student_id }
  showNoteDialog.value = true
}

const onNoteSaved = () => {
  fetchNotes()
}

const deleteNote = (noteId) => {
  $q.dialog({
    title: 'Conferma eliminazione',
    message: 'Sei sicuro di voler eliminare questa nota?',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await notesService.deleteNote(noteId)
      $q.notify({ type: 'positive', message: 'Nota eliminata con successo' })
      fetchNotes()
    } catch (e) {
      $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione della nota' })
    }
  })
}
</script>
