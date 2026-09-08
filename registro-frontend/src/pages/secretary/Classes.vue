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
        <q-btn
          color="indigo-7"
          icon="published_with_changes"
          label="Migrazione Anno"
          class="rounded-lg shadow-sm"
          no-caps
          @click="openMigrationWizard"
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
        virtual-scroll
        :virtual-scroll-item-size="48"
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
            <q-btn flat round dense icon="menu_book" color="indigo" :aria-label="t('secretaryClasses.manageAssignments') || 'Gestione Materie & Docenti'" @click="openAssignmentsDialog(props.row)">
              <q-tooltip>Gestione Materie & Docenti</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="groups" color="cyan-9" :aria-label="t('secretaryClasses.manageStudents') || 'Gestione Studenti della Classe'" @click="openStudentsDialog(props.row)">
              <q-tooltip>Gestione Studenti della Classe</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="auto_stories" color="emerald" :aria-label="t('routeTitles.textbooks') || 'Libri di Testo'" @click="openTextbooksDialog(props.row)">
              <q-tooltip>Libri di Testo</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="calendar_today" color="orange" :aria-label="t('routeTitles.timetable') || 'Orario Settimanale'" @click="openScheduleDialog(props.row)">
              <q-tooltip>Orario Settimanale</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="edit" color="primary" :aria-label="t('common.edit') || 'Modifica Classe'" @click="openDialog(props.row)">
              <q-tooltip>Modifica Classe</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="delete" color="negative" :aria-label="t('common.delete') || 'Elimina Classe'" @click="confirmDelete(props.row)">
              <q-tooltip>Elimina Classe</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Create/Edit Class -->
    <ClassFormDialog
      v-model="showDialog"
      :is-edit="isEdit"
      :initial-data="form"
      :academic-year-options="academicYearOptions"
      :teacher-user-options="teacherUserOptions"
      @saved="refreshClasses"
    />

    <!-- Assignments Dialog -->
    <ClassAssignmentsDialog
      v-model="showAssignmentsDialog"
      :target-class="currentClass"
      :subject-options="subjectOptions"
      :teacher-options="teacherOptions"
      :school-id="authStore.user?.school_id || ''"
      @subjects-updated="fetchSchoolData"
    />

    <!-- Textbooks Dialog -->
    <ClassTextbooksDialog
      v-model="showTextbooksDialog"
      :target-class="currentClass"
      :subject-options="subjectOptions"
      :school-id="authStore.user?.school_id || ''"
      @updated="refreshClasses"
    />

    <!-- Schedule Dialog -->
    <ClassScheduleDialog
      v-model="showScheduleDialog"
      :current-class="currentClass"
      :assignments="assignments"
      :current-schedule="currentSchedule"
      :loading="scheduleLoading"
      @save="saveSchedule"
    />

    <!-- Students Assignment Dialog -->
    <ClassStudentsDialog
      v-model="showStudentsDialog"
      :target-class="currentClass"
      :school-id="authStore.user?.school_id || ''"
      @refresh="refreshClasses"
    />
    
    <!-- Quick Create Subject Dialog -->
    <q-dialog v-model="showSubjectDialog">
      <q-card style="min-width: 350px" class="rounded-xl shadow-24 bg-white">
        <q-form @submit="createSubject">
          <q-card-section class="q-pa-lg">
            <div class="text-h6 text-weight-bold text-slate-800">Nuova Materia</div>
          </q-card-section>
          <q-card-section class="q-px-lg q-pb-lg">
            <q-input
              v-model="newSubjectName"
              label="Nome Materia"
              outlined
              autofocus
              :rules="[val => (!!val && val.trim().length > 0) || (t('common.requiredField') || 'Campo obbligatorio')]"
            />
          </q-card-section>
          <q-card-actions align="right" class="q-pa-lg bg-slate-50">
            <q-btn flat label="Annulla" v-close-popup color="slate-400" no-caps />
            <q-btn label="Crea" color="primary" type="submit" class="rounded-lg q-px-lg" no-caps />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>




    <!-- Academic Year Migration Wizard Dialog -->
    <ClassYearMigrationDialog
      v-model="showMigrationDialog"
      :academic-year-options="academicYearOptions"
      :current-year-str="selectedYear || currentYearStr"
      :classes="classesStore.classes"
      :school-id="authStore.user?.school_id || ''"
      @migrated="onMigrationComplete"
    />

  </q-page>
</template>

<script setup>
import { ref, onMounted, reactive, computed, watch } from 'vue'
import { useClassesStore } from '@/stores/classes'
import { useAuthStore } from '@/stores/auth'
import adminService from '@/services/adminService'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import ClassFormDialog from '@/components/Secretary/ClassFormDialog.vue'
import ClassAssignmentsDialog from '@/components/Secretary/ClassAssignmentsDialog.vue'
import ClassStudentsDialog from '@/components/Secretary/ClassStudentsDialog.vue'
import ClassTextbooksDialog from '@/components/Secretary/ClassTextbooksDialog.vue'
import ClassScheduleDialog from '@/components/Secretary/ClassScheduleDialog.vue'
import ClassYearMigrationDialog from '@/components/Secretary/ClassYearMigrationDialog.vue'

const $q = useQuasar()
const { t } = useI18n()
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

// Students Management State
const showStudentsDialog = ref(false)

const openStudentsDialog = (row) => {
  currentClass.value = row
  showStudentsDialog.value = true
}

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

const form = reactive({
  id: null,
  name: '',
  section: '',
  articolazione: '',
  location: '',
  academic_year: currentYearStr,
  coordinator_id: ''
})

const columns = [
  { name: 'name', label: 'Classe', align: 'left', field: row => `${row.name || ''}${row.section || ''}${row.articolazione ? ' - ' + row.articolazione : ''}`, sortable: true },
  { name: 'location', label: 'Sede', align: 'center', field: row => row.location || 'Sede Centrale', sortable: true },
  { name: 'academic_year', label: 'Anno Accademico', align: 'center', field: 'academic_year', sortable: true },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

const subjectOptions = computed(() => subjects.value.map(s => ({ label: s.name, value: s.id })))
const teacherOptions = computed(() => teachers.value.map(t => ({ label: `${t.last_name} ${t.first_name}`, value: t.id })))
const teacherUserOptions = computed(() => teachers.value.map(t => ({ label: `${t.last_name} ${t.first_name}`, value: t.user_id })))

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
        const [sRes, tRes] = await Promise.all([
            adminService.getSubjects(authStore.user.school_id),
            adminService.getTeachersList(authStore.user.school_id)
        ])
        subjects.value = sRes.data || []
        teachers.value = tRes.data || []
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
      location: '',
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
const openTextbooksDialog = (row) => {
    currentClass.value = row
    showTextbooksDialog.value = true
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

// Academic Year Migration State & Methods
const showMigrationDialog = ref(false)

const openMigrationWizard = () => {
  showMigrationDialog.value = true
}

const onMigrationComplete = async () => {
  await refreshClasses()
}

defineExpose({
    showDialog,
    isEdit,
    form,
    showAssignmentsDialog,
    showStudentsDialog,
    showTextbooksDialog,
    showSubjectDialog,
    currentClass,
    assignments,
    assignForm,
    newSubjectName,
    openDialog,
    saveClass,
    confirmDelete,
    openAssignmentsDialog,
    openStudentsDialog,
    openTextbooksDialog,
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
