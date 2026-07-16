<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">Registro Studenti</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">Gestione iscrizioni, classi e associazioni familiari</p>
      </div>
      <q-btn color="primary" icon="person_add" label="Nuova Iscrizione" class="rounded-lg shadow-sm" @click="openEnrollment" />
    </div>

    <!-- Student List -->
    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden">
      <q-table
        :rows="students"
        :columns="columns"
        :filter="filter"
        :loading="loading"
        row-key="id"
        flat
        class="bg-white"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template v-slot:top-right>
          <q-input dense debounce="300" v-model="filter" placeholder="Cerca studente..." outlined class="bg-white">
            <template v-slot:append>
              <q-icon name="search" color="grey-5" />
            </template>
          </q-input>
        </template>
        
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-right">
            <q-btn flat round dense icon="edit" color="blue-600" @click="editStudent(props.row)">
              <q-tooltip>Modifica Studente</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="family_restroom" color="indigo-600" @click="openGuardians(props.row)">
              <q-tooltip>Gestione Genitori</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="school" color="emerald-600" @click="openRecords(props.row)">
              <q-tooltip>Voti e Assenze</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <StudentRecords v-model="showRecords" :student="selectedStudentForRecords" />

    <!-- Enrollment / Edit Dialog -->
    <q-dialog v-model="showUserDialog" persistent>
      <q-card style="min-width: 500px" class="rounded-xl shadow-2xl">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold">{{ isEditing ? 'Modifica Studente' : 'Nuova Iscrizione' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit="saveStudent" class="q-gutter-md">
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-input v-model="userForm.first_name" label="Nome" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
              </div>
              <div class="col-6">
                <q-input v-model="userForm.last_name" label="Cognome" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
              </div>
            </div>
            
            <q-input v-model="userForm.email" label="Email" outlined dense type="email" :rules="[val => !!val || 'Obbligatorio']" />
            <q-input v-model="userForm.fiscal_code" label="Codice Fiscale" outlined dense maxlength="16" />
            
            <q-select
              v-model="userForm.class_id"
              :options="classOptions"
              label="Classe di Appartenenza"
              outlined
              dense
              emit-value
              map-options
              :loading="loadingClasses"
            />

            <q-input
              v-if="!isEditing"
              v-model="userForm.password"
              label="Password Provvisoria"
              outlined dense
              type="password"
              :rules="[val => !!val || 'Campo obbligatorio', val => val.length >= 8 || 'Minimo 8 caratteri']"
            />

            <div class="row justify-end q-mt-lg">
              <q-btn label="Annulla" flat v-close-popup color="grey-7" class="q-mr-sm" />
              <q-btn :label="isEditing ? 'Aggiorna' : 'Crea Studente'" type="submit" color="primary" class="q-px-lg rounded-md" :loading="saving" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Guardians Management Dialog -->
    <q-dialog v-model="showGuardiansDialog">
      <q-card style="min-width: 600px; max-width: 90vw;" class="rounded-xl shadow-2xl overflow-hidden">
        <q-card-section class="bg-indigo-600 text-white row items-center">
          <div class="text-h6 text-weight-bold">Genitori / Tutori</div>
          <q-space />
          <div class="text-subtitle2">{{ selectedStudent?.last_name }} {{ selectedStudent?.first_name }}</div>
          <q-btn icon="close" flat round dense v-close-popup class="q-ml-md" />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div class="row q-col-gutter-md">
            <!-- Current Guardians -->
            <div class="col-12 col-md-7">
              <div class="text-subtitle2 text-slate-700 q-mb-sm">Associazioni Attuali</div>
              <q-list bordered separator class="rounded-lg bg-white overflow-hidden">
                <q-item v-if="guardians.length === 0">
                  <q-item-section class="text-grey italic text-center q-pa-lg">
                    Nessun genitore associato a questo studente
                  </q-item-section>
                </q-item>
                <q-item v-for="g in guardians" :key="g.id">
                  <q-item-section avatar>
                    <q-avatar color="indigo-50" text-color="indigo-700" icon="person" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">{{ g.first_name }} {{ g.last_name }}</q-item-label>
                    <q-item-label caption>{{ g.email }}</q-item-label>
                    <q-item-label caption class="text-indigo-600 text-weight-medium">{{ g.relationship_type || 'Genitore' }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-btn flat round dense icon="link_off" color="negative" @click="removeGuardian(g)">
                      <q-tooltip>Rimuovi Associazione</q-tooltip>
                    </q-btn>
                  </q-item-section>
                </q-item>
              </q-list>
            </div>

            <!-- Add New Guardian -->
            <div class="col-12 col-md-5">
              <div class="text-subtitle2 text-slate-700 q-mb-sm">Aggiungi Genitore</div>
              <q-card flat bordered class="rounded-lg bg-white q-pa-none overflow-hidden">
                <q-card-section class="bg-slate-50 q-pa-md border-b">
                   <q-input
                    v-model="parentSearchText"
                    placeholder="Cerca per nome/cognome..."
                    outlined
                    dense
                    debounce="300"
                    @update:model-value="searchParents"
                    class="bg-white"
                  >
                    <template v-slot:append>
                      <q-icon name="search" />
                    </template>
                  </q-input>
                </q-card-section>

                <q-scroll-area style="height: 300px;">
                  <q-list separator>
                    <q-item v-if="parentOptions.length === 0" class="text-center q-pa-md">
                      <q-item-section class="text-grey">Nessun genitore trovato</q-item-section>
                    </q-item>
                    <q-item 
                      v-for="p in parentOptions" 
                      :key="p.id" 
                      clickable 
                      @click="parentSearch = p"
                      :active="parentSearch?.id === p.id"
                      active-class="bg-indigo-50 text-indigo-700"
                    >
                      <q-item-section avatar>
                        <q-avatar size="32px" color="indigo-100" text-color="indigo-700">
                          {{ p.last_name.charAt(0) }}
                        </q-avatar>
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">{{ p.last_name }} {{ p.first_name }}</q-item-label>
                        <q-item-label caption>{{ p.email }}</q-item-label>
                      </q-item-section>
                      <q-item-section side v-if="parentSearch?.id === p.id">
                        <q-icon name="check_circle" color="indigo" />
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-scroll-area>

                <q-card-section class="q-pa-md border-t bg-slate-50">
                  <q-select
                    v-model="newGuardianRel"
                    :options="['Padre', 'Madre', 'Tutore Legale', 'Delegato']"
                    label="Relazione"
                    outlined
                    dense
                    class="bg-white q-mb-md"
                  />

                  <q-btn 
                    label="Associa Selezionato" 
                    color="indigo" 
                    unelevated
                    class="full-width rounded-md text-weight-bold" 
                    @click="addGuardian" 
                    :disable="!parentSearch || !newGuardianRel"
                    :loading="addingGuardian"
                  />
                  
                  <div class="text-caption text-grey-6 q-mt-md text-center">
                    Se il genitore non è registrato, crealo prima in <router-link to="/secretary/users" class="text-indigo-600 font-bold">Gestione Utenti</router-link>
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { userService } from 'src/services/userService'
import adminService from 'src/services/adminService'
import { useAuthStore } from 'src/stores/auth'
import { useQuasar } from 'quasar'
import StudentRecords from 'src/components/Secretary/StudentRecords.vue'

const $q = useQuasar()
const authStore = useAuthStore()

const loading = ref(false)
const saving = ref(false)
const filter = ref('')
const students = ref([])
const classes = ref([])
const loadingClasses = ref(false)

// User Form Dialog
const showUserDialog = ref(false)
const isEditing = ref(false)
const userForm = reactive({
    id: null,
    first_name: '',
    last_name: '',
    email: '',
    fiscal_code: '',
    role: 'student',
    class_id: null,
    school_id: '',
    password: ''
})

// Guardians Dialog
const showGuardiansDialog = ref(false)
const selectedStudent = ref(null)
const guardians = ref([])
const parentSearch = ref(null)
const parentSearchText = ref('')
const parentOptions = ref([])
const newGuardianRel = ref('Padre')
const addingGuardian = ref(false)

// Student Records Dialog
const showRecords = ref(false)
const selectedStudentForRecords = ref(null)

const openRecords = (student) => {
    selectedStudentForRecords.value = student
    showRecords.value = true
}

const columns = [
    { name: 'name', label: 'Nome', field: row => `${row.last_name} ${row.first_name}`, align: 'left', sortable: true },
    { name: 'class', label: 'Classe', field: row => row.class_name || row.ClassName || '-', align: 'center', sortable: true },
    { name: 'email', label: 'Email', field: 'email', align: 'left' },
    { name: 'fiscal_code', label: 'Codice Fiscale', field: 'fiscal_code', align: 'left' },
    { name: 'actions', label: 'Azioni', align: 'right' }
]

const classOptions = computed(() => {
    return classes.value.map(c => {
        let label = `${c.name || ''}${c.section || ''}`.trim()
        if (c.articolazione) {
            label += ` - ${c.articolazione}`
        }
        if (c.academic_year) {
            label += ` (${c.academic_year})`
        }
        return {
            label: label || `Classe ${c.id.substring(0, 8)}`,
            value: c.id
        }
    })
})

onMounted(() => {
    fetchStudents()
    fetchClasses()
})

const fetchStudents = async () => {
    loading.value = true
    try {
        const res = await userService.getAll({ role: 'student' })
        students.value = res.data.users || []
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento studenti' })
    } finally {
        loading.value = false
    }
}

const fetchClasses = async () => {
    if (!authStore.user?.school_id) return
    loadingClasses.value = true
    try {
        const res = await adminService.getSchoolClasses(authStore.user.school_id)
        classes.value = res.data
    } catch (e) {
        console.error("Error loading classes", e)
    } finally {
        loadingClasses.value = false
    }
}

const openEnrollment = () => {
    isEditing.value = false
    Object.assign(userForm, {
        id: null,
        first_name: '',
        last_name: '',
        email: '',
        fiscal_code: '',
        role: 'student',
        class_id: null,
        school_id: authStore.user.school_id,
        password: ''
    })
    fetchClasses()
    showUserDialog.value = true
}

const editStudent = (row) => {
    isEditing.value = true
    Object.assign(userForm, row)
    userForm.class_id = row.ClassID || row.class_id
    fetchClasses()
    showUserDialog.value = true
}

const saveStudent = async () => {
    saving.value = true
    try {
        const payload = { ...userForm }
        if (isEditing.value) {
            delete payload.password
            await userService.update(userForm.id, payload)
            $q.notify({ type: 'positive', message: 'Studente aggiornato' })
        } else {
            await userService.create(payload)
            $q.notify({ type: 'positive', message: 'Studente iscritto correttamente' })
        }
        showUserDialog.value = false
        fetchStudents()
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' })
    } finally {
        saving.value = false
    }
}

// Guardian Logic
const openGuardians = async (student) => {
    selectedStudent.value = student
    guardians.value = []
    parentSearch.value = null
    parentSearchText.value = ''
    showGuardiansDialog.value = true
    fetchGuardians()
    searchParents('') // Load initial list
}

const fetchGuardians = async () => {
    try {
        const res = await userService.getGuardians(selectedStudent.value.id)
        guardians.value = res.data || []
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento genitori' })
    }
}

const searchParents = async (val) => {
    try {
        const params = { role: 'parent', page_size: 50 }
        if (val && val.length >= 1) {
            params.q = val
        }
        
        const res = await userService.getAll(params)
        parentOptions.value = (res.data.users || [])
    } catch (e) {
        console.error(e)
    }
}

const addGuardian = async () => {
    if (!parentSearch.value) return
    addingGuardian.value = true
    try {
        await userService.addGuardian(selectedStudent.value.id, {
            parent_user_id: parentSearch.value.id,
            relationship_type: newGuardianRel.value
        })
        $q.notify({ type: 'positive', message: 'Genitore associato con successo' })
        parentSearch.value = null
        fetchGuardians()
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore associazione genitore' })
    } finally {
        addingGuardian.value = false
    }
}

const removeGuardian = (guardian) => {
    $q.dialog({
        title: 'Rimuovi Associazione',
        message: `Vuoi rimuovere l'associazione con ${guardian.first_name} ${guardian.last_name}?`,
        cancel: true,
        persistent: true,
        ok: { color: 'negative', label: 'Rimuovi' }
    }).onOk(async () => {
        try {
            await userService.removeGuardian(selectedStudent.value.id, guardian.id)
            $q.notify({ type: 'positive', message: 'Associazione rimossa' })
            fetchGuardians()
        } catch (e) {
            $q.notify({ type: 'negative', message: 'Errore rimozione associazione' })
        }
    })
}

defineExpose({
    showUserDialog,
    isEditing,
    userForm,
    fetchStudents,
    editStudent,
    saveStudent,
    openEnrollment
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
