<template>
  <q-page class="q-pa-md">
    <q-table
      title="Gestione Classi"
      :rows="classesStore.classes"
      :columns="columns"
      row-key="id"
      :filter="filter"
      :loading="classesStore.loading"
      :pagination.sync="pagination"
    >
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filter" placeholder="Cerca">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
        <q-btn color="primary" icon="add" label="Nuova Classe" class="q-ml-md" @click="openDialog()" />
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn flat round color="secondary" icon="menu_book" @click="openAssignmentsDialog(props.row)">
            <q-tooltip>Gestione Materie & Docenti</q-tooltip>
          </q-btn>
          <q-btn flat round color="primary" icon="edit" @click="openDialog(props.row)" />
          <q-btn flat round color="negative" icon="delete" @click="confirmDelete(props.row)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Create/Edit Class -->
    <q-dialog v-model="showDialog">
      <q-card style="min-width: 400px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ isEdit ? 'Modifica Classe' : 'Nuova Classe' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <q-form @submit="saveClass" class="q-gutter-md">
            <q-input v-model="form.name" label="Nome (es. 1A, 5B)" outlined :rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model="form.section" label="Sezione (es. A, B)" outlined :rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model="form.academic_year" label="Anno Accademico (es. 2024/2025)" outlined :rules="[val => !!val || 'Campo obbligatorio']" />
            
            <div class="row justify-end">
              <q-btn flat label="Annulla" color="primary" v-close-popup />
              <q-btn type="submit" label="Salva" color="primary" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Assignments (Cattedre) -->
    <q-dialog v-model="showAssignmentsDialog" full-width>
      <q-card>
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Gestione Materie - Classe {{ currentClass?.name }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
            <div class="row q-col-gutter-md">
                <!-- Left: List -->
                <div class="col-12 col-md-8">
                    <q-table
                        title="Programmazione Didattica"
                        :rows="assignments"
                        :columns="assignmentsColumns"
                        row-key="id"
                        flat
                        bordered
                    >
                        <template v-slot:body-cell-actions="props">
                            <q-td :props="props" auto-width>
                                <q-btn flat round dense color="negative" icon="delete" @click="removeAssignment(props.row)" />
                            </q-td>
                        </template>
                    </q-table>
                </div>
                
                <!-- Right: Add Form -->
                <div class="col-12 col-md-4">
                    <q-card flat bordered class="q-pa-md">
                        <div class="text-subtitle1 q-mb-md">Assegna Materia</div>
                        <q-form @submit="addAssignment" class="q-gutter-md">
                             <q-select
                                v-model="assignForm.subject_id"
                                :options="subjectOptions"
                                label="Materia"
                                outlined
                                emit-value
                                map-options
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
                                emit-value
                                map-options
                             />

                             <q-input
                                v-model.number="assignForm.hours_per_week"
                                label="Ore Settimanali"
                                type="number"
                                outlined
                                min="1"
                             />

                             <q-btn type="submit" label="Assegna" color="primary" class="full-width" />
                        </q-form>
                    </q-card>
                </div>
            </div>
        </q-card-section>
      </q-card>
    </q-dialog>
    
    <!-- Quick Create Subject Dialog -->
    <q-dialog v-model="showSubjectDialog">
        <q-card style="min-width: 300px">
            <q-card-section>
                <div class="text-h6">Nuova Materia</div>
            </q-card-section>
            <q-card-section>
                <q-input v-model="newSubjectName" label="Nome Materia" outlined autofocus @keyup.enter="createSubject" />
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn flat label="Crea" color="primary" @click="createSubject" />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { useClassesStore } from '@/stores/classes'
import { useAuthStore } from '@/stores/auth'
import adminService from '@/services/adminService'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const classesStore = useClassesStore()
const authStore = useAuthStore()

const filter = ref('')
const showDialog = ref(false)
const isEdit = ref(false)
const pagination = ref({ rowsPerPage: 10 })

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

const form = reactive({
  id: null,
  name: '',
  section: '',
  academic_year: '2024/2025', // Default
  coordinator_id: ''
})

const columns = [
  { name: 'name', label: 'Nome', align: 'left', field: 'name', sortable: true },
  { name: 'section', label: 'Sezione', align: 'left', field: 'section', sortable: true },
  { name: 'academic_year', label: 'Anno Accademico', align: 'center', field: 'academic_year', sortable: true },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

const assignmentsColumns = [
    { name: 'subject', label: 'Materia', field: 'subject_name', align: 'left', sortable: true },
    { name: 'teacher', label: 'Docente', field: row => row.teacher_name || 'N/A', align: 'left' },
    { name: 'hours', label: 'Ore', field: 'hours_per_week', align: 'center' },
    { name: 'actions', label: 'Azioni', align: 'right' }
]

const subjectOptions = computed(() => subjects.value.map(s => ({ label: s.name, value: s.id })))
const teacherOptions = computed(() => teachers.value.map(t => ({ label: `${t.last_name} ${t.first_name}`, value: t.id })))

// Watch subject selection
import { watch } from 'vue' // Ensure import

watch(() => assignForm.subject_id, async (newVal) => {
    assignForm.teacher_id = null
    if (newVal) {
        // Fetch teachers qualified for this subject
        try {
            const res = await adminService.getTeachersList(authStore.user.school_id, newVal)
            teachers.value = res.data || [] // Backend returns []Teacher
        } catch(e) {
            console.error("Error filtering teachers", e)
        }
    } else {
        // Reset to all teachers
        fetchSchoolData()
    }
})

onMounted(() => {
  if (authStore.user?.school_id) {
    classesStore.fetchClasses({ school_id: authStore.user.school_id })
    fetchSchoolData()
  }
})

const fetchSchoolData = async () => {
    try {
        const sRes = await adminService.getSubjects(authStore.user.school_id)
        subjects.value = sRes.data || []
        
        // Load default full list (users?role=teacher returns Users, but we want Teachers entity ideally)
        // adminService.getSchoolUsers returns Users (with role).
        // adminService.getTeachersList returns Teachers (with qual).
        // Let's switch to getTeachersList for consistency.
        const tRes = await adminService.getTeachersList(authStore.user.school_id)
        teachers.value = tRes.data || []
    } catch(e) {
        console.error("Error loading school data", e)
    }
}

const openDialog = (row = null) => {
  // ... existing code ...
  if (row) {
    isEdit.value = true
    form.id = row.id
    form.name = row.name
    form.section = row.section
    form.academic_year = row.academic_year
    form.coordinator_id = row.coordinator_id
  } else {
    isEdit.value = false
    form.id = null
    form.name = ''
    form.section = ''
    form.academic_year = '2024/2025'
    form.coordinator_id = ''
  }
  showDialog.value = true
}

const saveClass = async () => {
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
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel salvataggio' })
  }
}

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Conferma',
    message: `Vuoi eliminare la classe ${row.name}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await classesStore.deleteClass(row.id)
      $q.notify({ type: 'positive', message: 'Classe eliminata' })
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
    if (subjects.value.length === 0) fetchSchoolData() // Retry if empty
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
        // Reset form slightly?
        // assignForm.subject_id = null
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore assegnazione' })
    }
}

const removeAssignment = async (row) => {
    try {
        await adminService.removeClassSubject(row.class_id, row.id) // row.id is assignment ID
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
        fetchSchoolData() // Refresh list
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore creazione materia' })
    }
}
</script>
