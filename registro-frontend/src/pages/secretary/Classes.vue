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
          class="rounded-lg min-width-150"
          @update:model-value="onYearChange"
        />
        <q-btn color="primary" icon="add" label="Nuova Classe" class="rounded-lg shadow-sm" @click="openDialog()" />
      </div>
    </div>

    <!-- Classes List -->
    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden bg-white">
      <q-table
        :rows="classesStore.classes"
        :columns="columns"
        :filter="filter"
        :loading="classesStore.loading"
        row-key="id"
        flat
        class="bg-transparent"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template #top-right>
          <q-input dense debounce="300" v-model="filter" placeholder="Cerca classe..." outlined>
            <template #append>
              <q-icon name="search" color="grey-5" />
            </template>
          </q-input>
        </template>
        
        <template #header-cell="props">
          <q-th :props="props" class="text-slate-500 font-bold">
            {{ props.col.label }}
          </q-th>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" class="text-right">
            <q-btn flat round dense icon="menu_book" color="indigo" @click="openAssignmentsDialog(props.row)">
              <q-tooltip>Gestione Materie & Docenti</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="auto_stories" color="emerald" @click="openTextbooksDialog(props.row)">
              <q-tooltip>Libri di Testo</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="calendar_today" color="orange" @click="openScheduleDialog(props.row)">
              <q-tooltip>Orario Settimanale</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="edit" color="primary" @click="openDialog(props.row)">
              <q-tooltip>Modifica Classe</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row)">
              <q-tooltip>Elimina Classe</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Create/Edit Class -->
    <q-dialog v-model="showDialog" persistent class="premium-dialog">
      <q-card style="min-width: 450px" class="rounded-xl overflow-hidden shadow-24">
        <q-card-section class="bg-gradient-primary text-white row items-center q-pa-lg">
          <div class="text-h6 text-weight-bold">{{ isEdit ? 'Modifica Classe' : 'Nuova Classe' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-xl">
          <q-form @submit="saveClass" class="q-gutter-y-lg">
            <div class="row q-col-gutter-lg">
              <div class="col-8">
                <q-input v-model="form.name" label="Nome (es. 1, 5)" outlined :rules="[val => !!val || 'Obbligatorio']" />
              </div>
              <div class="col-4">
                <q-input v-model="form.section" label="Sezione (es. A, B)" outlined :rules="[val => !!val || 'Obbligatorio']" />
              </div>
            </div>
            
            <q-input v-model="form.articolazione" label="Articolazione (es. Informatica, Telecomunicazioni - Opzionale)" outlined />
            
            <q-select
              v-model="form.academic_year"
              :options="academicYearOptions"
              label="Anno Accademico"
              outlined
              :rules="[val => !!val || 'Obbligatorio']"
            />
            
            <q-select
              v-model="form.coordinator_id"
              :options="teacherUserOptions"
              label="Coordinatore di Classe"
              outlined
              emit-value
              map-options
              clearable
            />
            
            <div class="row justify-end q-mt-xl q-gutter-sm">
              <q-btn label="Annulla" flat v-close-popup color="slate-400" />
              <q-btn :label="isEdit ? 'Aggiorna' : 'Crea Classe'" type="submit" color="primary" class="q-px-xl rounded-lg shadow-sm" :loading="saving" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Assignments Dialog -->
    <q-dialog v-model="showAssignmentsDialog">
      <q-card style="width: min(1100px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
        <q-card-section class="bg-gradient-primary text-white row items-center q-pa-md shrink-0">
          <div class="row items-center">
            <q-avatar color="white-20" text-color="white" icon="menu_book" class="q-mr-sm" size="36px" />
            <div>
              <div class="text-h6 text-weight-bold">Cattedre - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
              <div class="text-subtitle2 opacity-80">{{ currentClass?.academic_year }}</div>
            </div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md col overflow-y-auto">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-7">
              <q-table
                title="Programmazione Didattica"
                :rows="assignments"
                :columns="assignmentsColumns"
                row-key="id"
                flat
                class="bg-transparent border-slate-100 rounded-xl"
              >
                <template #header-cell="props">
                  <q-th :props="props" class="text-slate-500 font-bold">
                    {{ props.col.label }}
                  </q-th>
                </template>

                <template #body-cell-actions="props">
                  <q-td :props="props" auto-width>
                    <q-btn flat round dense color="negative" icon="delete" @click="removeAssignment(props.row)" />
                  </q-td>
                </template>
              </q-table>
            </div>
            
            <div class="col-12 col-md-5">
              <q-card flat class="rounded-xl bg-slate-50 q-pa-md border-slate-200">
                <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Assegna Materia</div>
                <q-form @submit="addAssignment" class="q-gutter-y-md">
                  <q-select
                    v-model="assignForm.subject_id"
                    :options="subjectOptions"
                    label="Materia *"
                    outlined
                    dense
                    emit-value
                    map-options
                    :rules="[val => !!val || 'Seleziona materia']"
                  >
                    <template #no-option>
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
                  />

                  <q-input
                    v-model.number="assignForm.hours_per_week"
                    label="Ore Settimanali"
                    type="number"
                    outlined
                    dense
                    min="1"
                  />

                  <q-btn type="submit" label="Assegna Cattedra" color="primary" class="full-width rounded-lg q-py-sm shadow-sm q-mt-md" no-caps />
                </q-form>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Textbooks Dialog -->
    <q-dialog v-model="showTextbooksDialog">
      <q-card style="width: min(1100px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
        <q-card-section class="bg-gradient-premium text-white row items-center q-pa-md shrink-0">
          <div class="row items-center">
            <q-avatar color="white-20" text-color="white" icon="auto_stories" class="q-mr-sm" size="36px" />
            <div>
              <div class="text-h6 text-weight-bold">Adozioni Libri - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
              <div class="text-subtitle2 opacity-80">{{ currentClass?.academic_year }}</div>
            </div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md col overflow-y-auto">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-7">
              <q-table
                title="Libri Adottati"
                :rows="classTextbooks"
                :columns="textbookColumns"
                row-key="id"
                flat
                class="bg-transparent border-slate-100 rounded-xl"
              >
                <template #header-cell="props">
                  <q-th :props="props" class="text-slate-500 font-bold">
                    {{ props.col.label }}
                  </q-th>
                </template>

                <template #body-cell-actions="props">
                  <q-td :props="props" auto-width>
                    <q-btn flat round dense color="negative" icon="delete" @click="removeTextbook(props.row)" />
                  </q-td>
                </template>
              </q-table>
            </div>
            
            <div class="col-12 col-md-5">
              <q-card flat class="rounded-xl bg-slate-50 q-pa-md border-slate-200">
                <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Adotta Libro</div>
                <q-form @submit="addTextbookToClass" class="q-gutter-y-md">
                  <q-select
                    v-model="textbookForm.textbook_id"
                    :options="allTextbooksOptions"
                    label="Libro *"
                    outlined
                    dense
                    emit-value
                    map-options
                    :rules="[val => !!val || 'Seleziona libro']"
                  />
                  <q-select
                    v-model="textbookForm.subject_id"
                    :options="subjectOptions"
                    label="Materia *"
                    outlined
                    dense
                    emit-value
                    map-options
                    :rules="[val => !!val || 'Seleziona materia']"
                  />
                  <q-checkbox v-model="textbookForm.is_optional" label="Il libro è opzionale" class="text-slate-700" />
                  <q-btn type="submit" label="Conferma Adozione" color="primary" class="full-width rounded-lg q-py-sm shadow-sm q-mt-md" no-caps />
                </q-form>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Schedule Dialog -->
    <q-dialog v-model="showScheduleDialog">
      <q-card style="width: min(1200px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
        <q-card-section class="bg-gradient-warning text-white row items-center q-pa-md shrink-0">
          <div class="row items-center">
            <q-avatar color="white-20" text-color="white" icon="calendar_today" class="q-mr-sm" size="36px" />
            <div>
              <div class="text-h6 text-weight-bold">Orario Settimanale - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
              <div class="text-subtitle2 opacity-80">{{ currentClass?.academic_year }}</div>
            </div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md col overflow-y-auto">
          <ScheduleGrid 
            :assignments="assignments"
            :initial-schedule="currentSchedule"
            :loading="scheduleLoading"
            @save="saveSchedule"
          />
        </q-card-section>
      </q-card>
    </q-dialog>
    
    <!-- Quick Create Subject Dialog -->
    <q-dialog v-model="showSubjectDialog">
      <q-card style="min-width: 350px" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="q-pa-lg">
          <div class="text-h6 text-weight-bold text-slate-800">Nuova Materia</div>
        </q-card-section>
        <q-card-section class="q-px-lg q-pb-lg">
          <q-input v-model="newSubjectName" label="Nome Materia" outlined autofocus @keyup.enter="createSubject" />
        </q-card-section>
        <q-card-actions align="right" class="q-pa-lg bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup color="slate-400" />
          <q-btn label="Crea" color="primary" class="rounded-lg q-px-lg" @click="createSubject" />
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
import ScheduleGrid from 'src/components/Secretary/ScheduleGrid.vue'

const $q = useQuasar()
const classesStore = useClassesStore()
const authStore = useAuthStore()

const filter = ref('')
const showDialog = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const getCurrentAcademicYear = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth() + 1; // 1-12
  // In Italy, the academic year usually starts in September
  if (month >= 9) { 
    return `${year}/${year + 1}`;
  } else {
    return `${year - 1}/${year}`;
  }
}

const currentYearStr = getCurrentAcademicYear();
const selectedYear = ref(currentYearStr);

const currentStart = parseInt(currentYearStr.split('/')[0]);
const academicYearOptions = [
  `${currentStart - 1}/${currentStart}`,
  currentYearStr,
  `${currentStart + 1}/${currentStart + 2}`
]

// Assignments State
const showAssignmentsDialog = ref(false)
const showSubjectDialog = ref(false)
const currentClass = ref(null)
const assignments = ref([])
const subjects = ref([])
const teachers = ref([])
const newSubjectName = ref('')

const showScheduleDialog = ref(false)
const scheduleLoading = ref(false)
const currentSchedule = ref([])

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
  articolazione: '',
  academic_year: currentYearStr,
  coordinator_id: ''
})

const columns = [
  { name: 'name', label: 'Classe', align: 'left', field: row => `${row.name || ''}${row.section || ''}${row.articolazione ? ' - ' + row.articolazione : ''}`, sortable: true },
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
      articolazione: '',
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


// Schedule Logic
const openScheduleDialog = async (row) => {
    currentClass.value = row
    showScheduleDialog.value = true
    scheduleLoading.value = true
    try {
        // We need both assignments (for options) and the current schedule
        await fetchAssignments(row.id)
        const res = await adminService.getClassSchedule(row.id)
        currentSchedule.value = res.data || []
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento orario' })
    } finally {
        scheduleLoading.value = false
    }
}

const saveSchedule = async (entries) => {
    if (!currentClass.value) return
    scheduleLoading.value = true
    try {
        await adminService.saveClassSchedule(currentClass.value.id, { entries })
        $q.notify({ type: 'positive', message: 'Orario salvato con successo' })
        showScheduleDialog.value = false
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore durante il salvataggio dell\'orario' })
    } finally {
        scheduleLoading.value = false
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
