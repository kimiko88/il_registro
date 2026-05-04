<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">Gestione Classi</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">Pianificazione classi, cattedre e adozioni libri</p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-select
          v-model="selectedYear"
          :options="academicYearOptions"
          label="Anno Accademico"
          outlined
          dense
          class="bg-white rounded-lg"
          style="min-width: 150px"
          @update:model-value="onYearChange"
        />
        <q-btn color="primary" icon="add" label="Nuova Classe" class="rounded-lg shadow-sm" @click="openDialog()" />
      </div>
    </div>

    <!-- Classes List -->
    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden">
      <q-table
        :rows="classesStore.classes"
        :columns="columns"
        :filter="filter"
        :loading="classesStore.loading"
        row-key="id"
        flat
        class="bg-white"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template v-slot:top-right>
          <q-input dense debounce="300" v-model="filter" placeholder="Cerca classe..." outlined class="bg-white">
            <template v-slot:append>
              <q-icon name="search" color="grey-5" />
            </template>
          </q-input>
        </template>
        
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-right">
            <q-btn flat round dense icon="menu_book" color="indigo-600" @click="openAssignmentsDialog(props.row)">
              <q-tooltip>Gestione Materie & Docenti</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="auto_stories" color="emerald-600" @click="openTextbooksDialog(props.row)">
              <q-tooltip>Libri di Testo</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="edit" color="blue-600" @click="openDialog(props.row)">
              <q-tooltip>Modifica Classe</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="delete" color="red-600" @click="confirmDelete(props.row)">
              <q-tooltip>Elimina Classe</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Create/Edit Class -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 450px" class="rounded-xl shadow-2xl">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold">{{ isEdit ? 'Modifica Classe' : 'Nuova Classe' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit="saveClass" class="q-gutter-md">
            <div class="row q-col-gutter-sm">
              <div class="col-8">
                <q-input v-model="form.name" label="Nome (es. 1, 5)" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
              </div>
              <div class="col-4">
                <q-input v-model="form.section" label="Sezione (es. A, B)" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
              </div>
            </div>
            
            <q-select
              v-model="form.academic_year"
              :options="academicYearOptions"
              label="Anno Accademico"
              outlined
              dense
              :rules="[val => !!val || 'Obbligatorio']"
            />
            
            <q-select
              v-model="form.coordinator_id"
              :options="teacherUserOptions"
              label="Coordinatore di Classe"
              outlined
              dense
              emit-value
              map-options
              clearable
            />
            
            <div class="row justify-end q-mt-lg">
              <q-btn label="Annulla" flat v-close-popup color="grey-7" class="q-mr-sm" />
              <q-btn :label="isEdit ? 'Aggiorna' : 'Crea Classe'" type="submit" color="primary" class="q-px-lg rounded-md" :loading="saving" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Assignments Dialog -->
    <q-dialog v-model="showAssignmentsDialog" full-width>
      <q-card class="rounded-xl overflow-hidden shadow-2xl">
        <q-card-section class="bg-indigo-600 text-white row items-center">
          <div class="text-h6 text-weight-bold">Cattedre - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
          <q-space />
          <div class="text-subtitle2">{{ currentClass?.academic_year }}</div>
          <q-btn icon="close" flat round dense v-close-popup class="q-ml-md" />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-8">
              <q-table
                title="Programmazione Didattica"
                :rows="assignments"
                :columns="assignmentsColumns"
                row-key="id"
                flat
                bordered
                class="rounded-lg"
              >
                <template v-slot:body-cell-actions="props">
                  <q-td :props="props" auto-width>
                    <q-btn flat round dense color="negative" icon="delete" @click="removeAssignment(props.row)" />
                  </q-td>
                </template>
              </q-table>
            </div>
            
            <div class="col-12 col-md-4">
              <q-card flat bordered class="rounded-lg bg-slate-50 q-pa-md">
                <div class="text-subtitle1 text-weight-bold q-mb-md">Assegna Materia</div>
                <q-form @submit="addAssignment" class="q-gutter-md">
                  <q-select
                    v-model="assignForm.subject_id"
                    :options="subjectOptions"
                    label="Materia"
                    outlined
                    dense
                    emit-value
                    map-options
                    class="bg-white"
                    :rules="[val => !!val || 'Seleziona materia']"
                  >
                    <template v-slot:no-option>
                      <q-item>
                        <q-item-section class="text-grey">Nessuna materia trovata</q-item-section>
                      </q-item>
                      <q-item clickable @click="openCreateSubject">
                        <q-item-section class="text-primary text-weight-bold">AGGIUNGI NUOVA MATERIA</q-item-section>
                      </q-item>
                    </template>
                  </q-select>

                  <q-select
                    v-model="assignForm.teacher_id"
                    :options="teacherOptions"
                    label="Docente"
                    outlined
                    dense
                    emit-value
                    map-options
                    class="bg-white"
                  />

                  <q-input
                    v-model.number="assignForm.hours_per_week"
                    label="Ore Settimanali"
                    type="number"
                    outlined
                    dense
                    min="1"
                    class="bg-white"
                  />

                  <q-btn type="submit" label="Assegna" color="indigo-600" class="full-width rounded-md" />
                </q-form>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Textbooks Dialog -->
    <q-dialog v-model="showTextbooksDialog" full-width>
      <q-card class="rounded-xl overflow-hidden shadow-2xl">
        <q-card-section class="bg-emerald-600 text-white row items-center">
          <div class="text-h6 text-weight-bold">Adozioni Libri - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
          <q-space />
          <div class="text-subtitle2">{{ currentClass?.academic_year }}</div>
          <q-btn icon="close" flat round dense v-close-popup class="q-ml-md" />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-8">
              <q-table
                title="Libri Adottati"
                :rows="classTextbooks"
                :columns="textbookColumns"
                row-key="id"
                flat
                bordered
                class="rounded-lg"
              >
                <template v-slot:body-cell-actions="props">
                  <q-td :props="props" auto-width>
                    <q-btn flat round dense color="negative" icon="delete" @click="removeTextbook(props.row)" />
                  </q-td>
                </template>
              </q-table>
            </div>
            
            <div class="col-12 col-md-4">
              <q-card flat bordered class="rounded-lg bg-slate-50 q-pa-md">
                <div class="text-subtitle1 text-weight-bold q-mb-md">Adotta Libro</div>
                <q-form @submit="addTextbookToClass" class="q-gutter-md">
                  <q-select
                    v-model="textbookForm.textbook_id"
                    :options="allTextbooksOptions"
                    label="Libro"
                    outlined
                    dense
                    emit-value
                    map-options
                    class="bg-white"
                    :rules="[val => !!val || 'Seleziona libro']"
                  />
                  <q-select
                    v-model="textbookForm.subject_id"
                    :options="subjectOptions"
                    label="Materia"
                    outlined
                    dense
                    emit-value
                    map-options
                    class="bg-white"
                    :rules="[val => !!val || 'Seleziona materia']"
                  />
                  <q-checkbox v-model="textbookForm.is_optional" label="Opzionale" />
                  <q-btn type="submit" label="Aggiungi" color="emerald-600" class="full-width rounded-md" />
                </q-form>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
    
    <!-- Quick Create Subject Dialog -->
    <q-dialog v-model="showSubjectDialog">
      <q-card style="min-width: 350px" class="rounded-xl">
        <q-card-section>
          <div class="text-h6 text-weight-bold">Nuova Materia</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="newSubjectName" label="Nome Materia" outlined dense autofocus @keyup.enter="createSubject" />
        </q-card-section>
        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Annulla" v-close-popup color="grey-7" />
          <q-btn label="Crea" color="primary" class="rounded-md q-px-md" @click="createSubject" />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, onMounted, reactive, computed, watch } from 'vue'
import { useClassesStore } from 'src/stores/classes'
import { useAuthStore } from 'src/stores/auth'
import adminService from 'src/services/adminService'
import { useQuasar } from 'quasar'
import { textbookService } from 'src/services/textbookService'

const $q = useQuasar()
const classesStore = useClassesStore()
const authStore = useAuthStore()

const filter = ref('')
const showDialog = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const selectedYear = ref('2024/2025')

const academicYearOptions = [
  '2023/2024',
  '2024/2025',
  '2025/2026'
]

// Assignments State
const showAssignmentsDialog = ref(false)
const showSubjectDialog = ref(false)
const currentClass = ref(null)
const assignments = ref([])
const subjects = ref([])
const teachers = ref([])
const newSubjectName = ref('')

const assignForm = reactive({
    subject_id: null,
    teacher_id: null,
    hours_per_week: 1
})

const showTextbooksDialog = ref(false)
const classTextbooks = ref([])
const allTextbooks = ref([])
const textbookForm = reactive({
    textbook_id: null,
    subject_id: null,
    is_optional: false
})

const form = reactive({
  id: null,
  name: '',
  section: '',
  academic_year: '2024/2025',
  coordinator_id: ''
})

const columns = [
  { name: 'name', label: 'Classe', align: 'left', field: row => `${row.name}${row.section}`, sortable: true },
  { name: 'academic_year', label: 'Anno Accademico', align: 'center', field: 'academic_year', sortable: true },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

const assignmentsColumns = [
    { name: 'subject', label: 'Materia', field: 'subject_name', align: 'left', sortable: true },
    { name: 'teacher', label: 'Docente', field: row => row.teacher_name || 'N/A', align: 'left' },
    { name: 'hours', label: 'Ore/Sett', field: 'hours_per_week', align: 'center' },
    { name: 'actions', label: 'Azioni', align: 'right' }
]

const textbookColumns = [
    { name: 'subject', label: 'Materia', field: 'subject_name', align: 'left' },
    { name: 'title', label: 'Titolo', field: 'title', align: 'left' },
    { name: 'author', label: 'Autore', field: 'author', align: 'left' },
    { name: 'optional', label: 'Opz.', field: row => row.is_optional ? 'Sì' : 'No', align: 'center' },
    { name: 'actions', label: 'Azioni', align: 'right' }
]

const subjectOptions = computed(() => subjects.value.map(s => ({ label: s.name, value: s.id })))
const teacherOptions = computed(() => teachers.value.map(t => ({ label: `${t.last_name} ${t.first_name}`, value: t.id })))
const teacherUserOptions = computed(() => teachers.value.map(t => ({ label: `${t.last_name} ${t.first_name}`, value: t.user_id })))
const allTextbooksOptions = computed(() => allTextbooks.value.map(b => ({ label: b.title, value: b.id })))

onMounted(() => {
  if (authStore.user?.school_id) {
    refreshClasses()
    fetchSchoolData()
  }
})

const refreshClasses = () => {
  classesStore.fetchClasses({ 
    school_id: authStore.user.school_id,
    academic_year: selectedYear.value
  })
}

// Watch subject selection to filter teachers
watch(() => assignForm.subject_id, async (newVal) => {
    assignForm.teacher_id = null
    if (newVal) {
        try {
            const res = await adminService.getTeachersList(authStore.user.school_id, newVal)
            teachers.value = res.data || []
        } catch(e) {
            console.error("Error filtering teachers", e)
        }
    } else {
        // Reset to all teachers
        fetchSchoolData()
    }
})

const onYearChange = () => {
  refreshClasses()
}

const fetchSchoolData = async () => {
    try {
        const [sRes, tRes, bRes] = await Promise.all([
            adminService.getSubjects(authStore.user.school_id),
            adminService.getTeachersList(authStore.user.school_id),
            textbookService.getAll()
        ])
        subjects.value = sRes.data || []
        teachers.value = tRes.data || []
        allTextbooks.value = bRes.data || []
    } catch(e) {
        console.error("Error loading school data", e)
    }
}

const openDialog = (row = null) => {
  if (row) {
    isEdit.value = true
    Object.assign(form, row)
  } else {
    isEdit.value = false
    Object.assign(form, {
      id: null,
      name: '',
      section: '',
      academic_year: selectedYear.value,
      coordinator_id: ''
    })
  }
  showDialog.value = true
}

const saveClass = async () => {
  saving.value = true
  try {
    const payload = { ...form, school_id: authStore.user.school_id };
    if (isEdit.value) {
      await classesStore.updateClass(form.id, payload)
      $q.notify({ type: 'positive', message: 'Classe aggiornata' })
    } else {
      await classesStore.createClass(payload)
      $q.notify({ type: 'positive', message: 'Classe creata' })
    }
    showDialog.value = false
    refreshClasses()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel salvataggio' })
  } finally {
    saving.value = false
  }
}

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: `Vuoi eliminare la classe ${row.name}${row.section}? Tutte le associazioni verranno rimosse.`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: 'Elimina' }
  }).onOk(async () => {
    try {
      await classesStore.deleteClass(row.id)
      $q.notify({ type: 'positive', message: 'Classe eliminata' })
      refreshClasses()
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Errore nell\'eliminazione' })
    }
  })
}

// Assignments Logic
const openAssignmentsDialog = async (row) => {
    currentClass.value = row
    showAssignmentsDialog.value = true
    await fetchAssignments(row.id)
}

const fetchAssignments = async (classId) => {
    try {
        const res = await adminService.getClassSubjects(classId)
        assignments.value = res.data || []
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento materie' })
    }
}

const addAssignment = async () => {
    if (!currentClass.value) return
    try {
        await adminService.assignSubjectToClass(currentClass.value.id, assignForm)
        $q.notify({ type: 'positive', message: 'Materia assegnata' })
        fetchAssignments(currentClass.value.id)
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore assegnazione' })
    }
}

const removeAssignment = async (row) => {
    try {
        await adminService.removeClassSubject(row.class_id, row.id)
        $q.notify({ type: 'positive', message: 'Assegnazione rimossa' })
        fetchAssignments(currentClass.value.id)
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore rimozione' })
    }
}

const openCreateSubject = () => {
    showSubjectDialog.value = true
}

const createSubject = async () => {
    if (!newSubjectName.value) return
    try {
        const payload = {
            name: newSubjectName.value,
            school_id: authStore.user.school_id
        }
        await adminService.createSubject(payload)
        $q.notify({ type: 'positive', message: 'Materia creata' })
        showSubjectDialog.value = false
        newSubjectName.value = ''
        fetchSchoolData()
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore creazione materia' })
    }
}

// Textbooks Logic
const openTextbooksDialog = async (row) => {
    currentClass.value = row
    showTextbooksDialog.value = true
    fetchClassTextbooks(row.id)
}

const fetchClassTextbooks = async (classId) => {
    try {
        const res = await textbookService.listByClass(classId)
        classTextbooks.value = res.data || []
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento libri' })
    }
}

const addTextbookToClass = async () => {
    try {
        await textbookService.assignToClass(currentClass.value.id, textbookForm)
        $q.notify({ type: 'positive', message: 'Libro adottato' })
        fetchClassTextbooks(currentClass.value.id)
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore adozione libro' })
    }
}

const removeTextbook = async (row) => {
    try {
        await textbookService.removeFromClass(row.id)
        $q.notify({ type: 'positive', message: 'Adozione rimossa' })
        fetchClassTextbooks(currentClass.value.id)
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore rimozione' })
    }
}

defineExpose({
    showDialog,
    isEdit,
    form,
    showAssignmentsDialog,
    showSubjectDialog,
    currentClass,
    assignments,
    assignForm,
    newSubjectName,
    openDialog,
    saveClass,
    confirmDelete,
    openAssignmentsDialog,
    addAssignment,
    removeAssignment,
    openCreateSubject,
    createSubject
})
</script>

<style scoped>
.rounded-xl { border-radius: 1rem; }
.rounded-lg { border-radius: 0.75rem; }
.rounded-md { border-radius: 0.5rem; }
.shadow-soft { box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05); }
.border-slate-100 { border: 1px solid #f1f5f9; }
.bg-slate-50 { background-color: #f8fafc; }
</style>
