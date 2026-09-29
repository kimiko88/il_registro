<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg sticky-header">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="groups" color="primary" class="q-mr-sm" />
          Gruppi Linguistici / Articolati
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Creazione e gestione composizione gruppi classe, sezioni articolate e gruppi di livello linguistico
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn
          color="primary"
          icon="add"
          label="Nuovo Gruppo Linguistico / Articolato"
          unelevated
          class="rounded-lg text-weight-bold"
          @click="openCreateModal"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="fetchGroups" />
      </div>
    </div>

    <!-- Groups Grid / List -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="48px" />
    </div>

    <div v-else-if="groups.length === 0" class="text-center q-pa-xl text-slate-400">
      <q-icon name="groups" size="64px" class="q-mb-md opacity-40" />
      <div class="text-h6">Nessun Gruppo Linguistico / Articolato creato</div>
      <div class="text-caption q-mb-md">Crea il primo gruppo per articolare gli studenti tra più sezioni o livelli di lingua.</div>
      <q-btn color="primary" icon="add" label="Crea Nuovo Gruppo" unelevated @click="openCreateModal" />
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="group in groups" :key="group.id" class="col-12 col-md-6 col-lg-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft full-height column justify-between hover-shadow transition-all">
          <q-card-section>
            <div class="row items-center justify-between q-mb-xs">
              <q-badge color="teal" class="text-weight-bold q-px-sm q-py-xs">
                {{ group.student_count || (group.students ? group.students.length : 0) }} Alunni Associati
              </q-badge>
              <div class="text-caption text-grey-6 text-weight-medium">
                Gruppo Scuola
              </div>
            </div>

            <div class="text-h6 text-weight-bold text-slate-800 q-mt-sm">
              {{ group.name }}
            </div>

            <!-- Subject & Teacher Badges -->
            <div class="row items-center q-gutter-xs q-mt-xs q-mb-sm">
              <q-chip
                v-if="group.subject_name || getSubjectName(group.subject_id)"
                dense
                size="12px"
                icon="menu_book"
                color="blue-1"
                text-color="blue-9"
                class="q-ma-none text-weight-medium"
              >
                {{ group.subject_name || getSubjectName(group.subject_id) }}
              </q-chip>
              <q-chip
                v-if="group.teacher_name || getTeacherName(group.teacher_id)"
                dense
                size="12px"
                icon="person"
                color="purple-1"
                text-color="purple-9"
                class="q-ma-none text-weight-medium"
              >
                {{ group.teacher_name || getTeacherName(group.teacher_id) }}
              </q-chip>
            </div>

            <div class="text-body2 text-slate-600 q-mt-xs q-mb-md line-clamp-2">
              {{ group.description || 'Nessuna descrizione specificata.' }}
            </div>
          </q-card-section>

          <q-card-actions class="bg-slate-50 border-t border-slate-200 justify-between q-px-md">
            <div class="row items-center q-gutter-xs">
              <q-btn flat dense icon="edit" color="primary" label="Modifica" @click="openEditModal(group)" />
              <q-btn flat dense icon="person_add" color="secondary" label="Studenti" @click="openManageStudentsModal(group)" />
            </div>
            <div>
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDeleteGroup(group)" />
            </div>
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Create / Edit Group Modal -->
    <q-dialog v-model="showCreateModal" persistent max-width="650px">
      <q-card style="width: 650px; max-width: 95vw" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md q-px-lg">
          <div class="text-h6 text-weight-bold">
            <q-icon name="groups" class="q-mr-xs" />
            {{ editingGroupId ? 'Modifica Gruppo Linguistico / Articolato' : 'Nuovo Gruppo Linguistico / Articolato' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-lg q-gutter-y-md">
          <q-input
            v-model="createForm.name"
            label="Nome Gruppo / Articolazione *"
            placeholder="es. Gruppo Inglese Avanzato B2 (Sez. A/B)"
            outlined dense hide-bottom-space
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <!-- Docente Assegnato -->
          <q-select
            v-model="createForm.teacher_id"
            :options="teacherOptions"
            option-value="id"
            option-label="displayName"
            emit-value
            map-options
            clearable
            outlined dense
            label="Docente Assegnato"
            placeholder="Seleziona docente"
          >
            <template v-slot:prepend>
              <q-icon name="person" color="primary" />
            </template>
          </q-select>

          <!-- Materia Coinvolta -->
          <q-select
            v-model="createForm.subject_id"
            :options="subjectOptions"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            clearable
            outlined dense
            label="Materia (es. Lingua Inglese, Francese...)"
            placeholder="Seleziona materia"
          >
            <template v-slot:prepend>
              <q-icon name="menu_book" color="primary" />
            </template>
          </q-select>

          <!-- Classi Coinvolte -->
          <q-select
            v-model="createForm.class_ids"
            :options="classOptions"
            option-value="id"
            option-label="displayName"
            emit-value
            map-options
            multiple
            use-chips
            outlined dense
            label="Classi Coinvolte nell'Articolazione"
            placeholder="Seleziona una o più classi"
            @update:model-value="onClassesChanged"
          >
            <template v-slot:prepend>
              <q-icon name="school" color="primary" />
            </template>
          </q-select>

          <!-- Quick Student Selection for Involved Classes -->
          <div v-if="createForm.class_ids && createForm.class_ids.length > 0" class="border rounded-lg bg-slate-50 q-pa-sm">
            <div class="row items-center justify-between q-mb-xs">
              <div class="text-caption text-weight-bold text-slate-700">
                Studenti delle Classi Selezionate ({{ classStudents.length }} trovati, {{ createForm.student_ids.length }} selezionati)
              </div>
              <div class="row q-gutter-xs">
                <q-btn flat dense size="xs" color="primary" label="Seleziona Tutti" @click="selectAllClassStudents" />
                <q-btn flat dense size="xs" color="grey-7" label="Deseleziona" @click="createForm.student_ids = []" />
              </div>
            </div>

            <div v-if="loadingClassStudents" class="text-center q-pa-sm">
              <q-spinner size="20px" color="primary" />
            </div>
            <q-scroll-area v-else style="height: 160px;" class="bg-white rounded border q-pa-xs">
              <div v-if="classStudents.length === 0" class="text-caption text-grey-5 text-center q-pa-sm">
                Nessuno studente trovato per le classi selezionate
              </div>
              <div v-else class="row q-col-gutter-xs">
                <div v-for="st in classStudents" :key="st.id" class="col-12 col-sm-6">
                  <q-checkbox
                    v-model="createForm.student_ids"
                    :val="st.id"
                    dense
                    size="sm"
                    color="primary"
                    :label="`${st.last_name || ''} ${st.first_name || ''} (${st.class_name || ''})`"
                  />
                </div>
              </div>
            </q-scroll-area>
          </div>

          <q-input
            v-model="createForm.description"
            type="textarea"
            rows="2"
            label="Descrizione e Finalità Didattica"
            placeholder="es. Gruppo trasversale per potenziamento lingua o articolazione orario"
            outlined dense hide-bottom-space
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn
            color="primary"
            icon="save"
            :label="editingGroupId ? 'Salva Modifiche' : 'Crea Gruppo'"
            unelevated
            :loading="submitting"
            @click="saveGroup"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Manage Students Modal -->
    <q-dialog v-model="showManageModal" persistent max-width="700px">
      <q-card style="width: 700px; max-width: 95vw" class="rounded-xl" v-if="selectedGroup">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md q-px-lg">
          <div>
            <div class="text-h6 text-weight-bold">Gestione Studenti: {{ selectedGroup.name }}</div>
            <div class="text-caption opacity-90">Seleziona o deseleziona gli alunni da includere nel gruppo articolato</div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-input
            v-model="studentFilter"
            placeholder="Cerca per nome o cognome..."
            outlined dense clearable class="q-mb-md"
          >
            <template v-slot:append><q-icon name="search" /></template>
          </q-input>

          <div v-if="loadingStudents" class="text-center q-pa-lg">
            <q-spinner color="primary" size="36px" />
          </div>

          <q-scroll-area v-else style="height: 380px;" class="border rounded-lg bg-slate-50 q-pa-xs">
            <q-list separator dense>
              <q-item v-for="student in filteredAllStudents" :key="student.id" clickable @click="toggleStudentSelection(student.id)">
                <q-item-section side>
                  <q-checkbox v-model="selectedStudentIds" :val="student.id" color="primary" />
                </q-item-section>
                <q-item-section avatar>
                  <q-avatar size="32px" color="indigo-1" text-color="indigo-8" icon="person" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ student.last_name }} {{ student.first_name }}</q-item-label>
                  <q-item-label caption class="text-slate-500">
                    {{ student.class_name ? `Classe ${student.class_name}` : `Matr: ${student.id.substring(0,8)}` }}
                  </q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-scroll-area>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" icon="check" label="Salva Composizione Gruppo" unelevated :loading="submitting" @click="saveGroupStudents" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import api from '@/services/api'
import adminService from '@/services/adminService'
import { useAuthStore } from '@/stores/auth'

const $q = useQuasar()
const { t } = useI18n()
const authStore = useAuthStore()

const groups = ref([])
const loading = ref(false)
const submitting = ref(false)

const showCreateModal = ref(false)
const showManageModal = ref(false)
const selectedGroup = ref(null)
const editingGroupId = ref(null)

const teachers = ref([])
const subjects = ref([])
const classes = ref([])

const allStudents = ref([])
const classStudents = ref([])
const loadingStudents = ref(false)
const loadingClassStudents = ref(false)
const studentFilter = ref('')
const selectedStudentIds = ref([])

const createForm = ref({
  name: '',
  subject_id: null,
  teacher_id: null,
  class_ids: [],
  student_ids: [],
  description: ''
})

const teacherOptions = computed(() => {
  return teachers.value.map(t => ({
    id: t.id,
    displayName: `${t.last_name || ''} ${t.first_name || ''}`.trim() || t.name || t.email || 'Docente'
  }))
})

const subjectOptions = computed(() => {
  return subjects.value.map(s => ({
    id: s.id,
    name: s.name
  }))
})

const classOptions = computed(() => {
  return classes.value.map(c => ({
    id: c.id,
    displayName: c.name || `${c.year_number || ''}${c.section || ''}`.trim() || 'Classe'
  }))
})

const filteredAllStudents = computed(() => {
  if (!studentFilter.value) return allStudents.value
  const q = studentFilter.value.toLowerCase()
  return allStudents.value.filter(s =>
    (s.first_name || '').toLowerCase().includes(q) ||
    (s.last_name || '').toLowerCase().includes(q) ||
    (s.class_name || '').toLowerCase().includes(q)
  )
})

const getSubjectName = (subjectId) => {
  if (!subjectId) return ''
  const s = subjects.value.find(sub => sub.id === subjectId)
  return s ? s.name : ''
}

const getTeacherName = (teacherId) => {
  if (!teacherId) return ''
  const t = teachers.value.find(tch => tch.id === teacherId)
  if (!t) return ''
  return `${t.last_name || ''} ${t.first_name || ''}`.trim() || t.name || t.email || ''
}

const loadMetadata = async () => {
  try {
    const [tRes, cRes, sRes] = await Promise.allSettled([
      adminService.getTeachersList(),
      adminService.getClasses(),
      api.get('/subjects')
    ])

    if (tRes.status === 'fulfilled') {
      teachers.value = tRes.value.data || []
    }
    if (cRes.status === 'fulfilled') {
      classes.value = cRes.value.data || []
    }
    if (sRes.status === 'fulfilled') {
      subjects.value = sRes.value.data || []
    }
  } catch (err) {
    console.error('Error loading metadata:', err)
  }
}

const fetchGroups = async () => {
  loading.value = true
  try {
    const res = await api.get('/groups')
    groups.value = res.data?.groups || res.data || []
  } catch (err) {
    groups.value = []
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  editingGroupId.value = null
  createForm.value = {
    name: '',
    subject_id: null,
    teacher_id: null,
    class_ids: [],
    student_ids: [],
    description: ''
  }
  classStudents.value = []
  showCreateModal.value = true
}

const openEditModal = (group) => {
  editingGroupId.value = group.id
  createForm.value = {
    name: group.name,
    subject_id: group.subject_id || null,
    teacher_id: group.teacher_id || null,
    class_ids: group.class_ids || [],
    student_ids: (group.students || []).map(s => s.id || s.student_id),
    description: group.description || ''
  }
  showCreateModal.value = true
  if (createForm.value.class_ids.length > 0) {
    onClassesChanged(createForm.value.class_ids)
  }
}

const onClassesChanged = async (classIds) => {
  if (!classIds || classIds.length === 0) {
    classStudents.value = []
    return
  }
  loadingClassStudents.value = true
  try {
    const uRes = await api.get('/users', { params: { role: 'student', page_size: 500 } })
    const all = uRes.data?.users || uRes.data || []
    classStudents.value = all.filter(s => classIds.includes(s.class_id))
  } catch (err) {
    classStudents.value = []
  } finally {
    loadingClassStudents.value = false
  }
}

const selectAllClassStudents = () => {
  createForm.value.student_ids = classStudents.value.map(s => s.id)
}

const saveGroup = async () => {
  if (!createForm.value.name) {
    $q.notify({ type: 'warning', message: 'Inserire il nome del gruppo' })
    return
  }
  submitting.value = true
  try {
    const payload = {
      school_id: authStore.user?.school_id || '',
      name: createForm.value.name,
      subject_id: createForm.value.subject_id || null,
      teacher_id: createForm.value.teacher_id || null,
      description: createForm.value.description || '',
      student_ids: createForm.value.student_ids || []
    }

    if (editingGroupId.value) {
      await api.put(`/groups/${editingGroupId.value}`, payload)
      if (payload.student_ids.length > 0) {
        await api.post(`/groups/${editingGroupId.value}/students`, { student_ids: payload.student_ids })
      }
      $q.notify({ type: 'positive', message: 'Gruppo modificato con successo' })
    } else {
      await api.post('/groups', payload)
      $q.notify({ type: 'positive', message: 'Gruppo Linguistico / Articolato creato con successo' })
    }

    showCreateModal.value = false
    await fetchGroups()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore salvataggio gruppo' })
  } finally {
    submitting.value = false
  }
}

const openManageStudentsModal = async (group) => {
  selectedGroup.value = group
  showManageModal.value = true
  loadingStudents.value = true
  try {
    const [uRes, gRes] = await Promise.all([
      api.get('/users', { params: { role: 'student', page_size: 500 } }),
      api.get(`/groups/${group.id}`)
    ])
    allStudents.value = uRes.data?.users || uRes.data || []
    const grpData = gRes.data || {}
    const existing = grpData.students || []
    selectedStudentIds.value = existing.map(s => s.id || s.student_id)
  } catch {
    selectedStudentIds.value = []
  } finally {
    loadingStudents.value = false
  }
}

const toggleStudentSelection = (id) => {
  const idx = selectedStudentIds.value.indexOf(id)
  if (idx === -1) {
    selectedStudentIds.value.push(id)
  } else {
    selectedStudentIds.value.splice(idx, 1)
  }
}

const saveGroupStudents = async () => {
  if (!selectedGroup.value) return
  submitting.value = true
  try {
    await api.post(`/groups/${selectedGroup.value.id}/students`, {
      student_ids: selectedStudentIds.value
    })
    $q.notify({ type: 'positive', message: 'Composizione gruppo salvata con successo' })
    showManageModal.value = false
    await fetchGroups()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante l\'aggiornamento degli studenti nel gruppo' })
  } finally {
    submitting.value = false
  }
}

const confirmDeleteGroup = (group) => {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: `Sei sicuro di voler eliminare il gruppo "${group.name}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await api.delete(`/groups/${group.id}`)
      $q.notify({ type: 'positive', message: 'Gruppo eliminato' })
      await fetchGroups()
    } catch {
      $q.notify({ type: 'negative', message: 'Errore eliminazione gruppo' })
    }
  })
}

onMounted(() => {
  loadMetadata()
  fetchGroups()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.hover-shadow:hover {
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
}
</style>
