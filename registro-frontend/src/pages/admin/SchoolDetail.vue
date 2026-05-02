<template>
  <q-page class="q-pa-md">
    <div v-if="loading" class="flex flex-center" style="height: 400px">
      <q-spinner color="primary" size="3em" />
    </div>

    <div v-else-if="school">
      <!-- Header -->
      <div class="row items-center q-mb-lg">
        <q-btn flat round icon="arrow_back" @click="$router.back()" class="q-mr-sm" />
        <div>
          <div class="text-h4 text-weight-bold">{{ school.name }}</div>
          <div class="text-subtitle1 text-grey-7">{{ school.code }} - {{ school.city }} ({{ school.province }})</div>
        </div>
        <q-space />
        <q-chip :color="school.is_active ? 'positive' : 'negative'" text-color="white">
          {{ school.is_active ? 'Attiva' : 'Disattiva' }}
        </q-chip>
      </div>

      <!-- Info Cards -->
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-12 col-md-8">
          <q-card class="h-100">
            <q-card-section>
              <div class="text-h6 q-mb-md">Dettagli Contatto</div>
              <q-list separator>
                <q-item>
                  <q-item-section avatar><q-icon name="place" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label caption>Indirizzo</q-item-label>
                    <q-item-label>{{ school.address }}, {{ school.zip_code }} {{ school.city }}</q-item-label>
                  </q-item-section>
                </q-item>
                <q-item>
                  <q-item-section avatar><q-icon name="phone" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label caption>Telefono</q-item-label>
                    <q-item-label>{{ school.phone || 'N/D' }}</q-item-label>
                  </q-item-section>
                </q-item>
                <q-item>
                  <q-item-section avatar><q-icon name="email" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label caption>Email</q-item-label>
                    <q-item-label>{{ school.email || 'N/D' }}</q-item-label>
                  </q-item-section>
                </q-item>
                <q-item>
                  <q-item-section avatar><q-icon name="language" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label caption>Sito Web</q-item-label>
                    <q-item-label>
                      <a v-if="school.website" :href="school.website" target="_blank">{{ school.website }}</a>
                      <span v-else>N/D</span>
                    </q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-md-4">
          <q-card class="h-100">
            <q-card-section>
              <div class="text-h6 q-mb-md">Statistiche</div>
              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-card bordered flat class="text-center q-pa-sm">
                    <div class="text-h4 text-primary text-weight-bold">{{ school.student_count || 0 }}</div>
                    <div class="text-caption text-grey">Studenti</div>
                  </q-card>
                </div>
                <div class="col-6">
                  <q-card bordered flat class="text-center q-pa-sm">
                    <div class="text-h4 text-secondary text-weight-bold">{{ school.teacher_count || 0 }}</div>
                    <div class="text-caption text-grey">Docenti</div>
                  </q-card>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
      
       <!-- Additional sections -->
       <q-card>
             <div class="row items-center justify-between q-pa-sm">
                <q-tabs
                    v-model="tab"
                    dense
                    class="text-grey"
                    active-color="primary"
                    indicator-color="primary"
                    align="justify"
                    narrow-indicator
                >
                    <q-tab name="classes" label="Classi" />
                    <q-tab name="users" label="Utenti" />
                </q-tabs>
                <q-btn 
                    v-if="tab === 'classes'"
                    color="primary" 
                    icon="add" 
                    label="Nuova Classe" 
                    size="sm" 
                    unelevated 
                    @click="openClassDialog"
                />
             </div>

           <q-separator />

           <q-tab-panels v-model="tab" animated>
             <q-tab-panel name="classes">
               <q-table
                 :rows="classes"
                 :columns="classColumns"
                 row-key="id"
                 :loading="loadingClasses"
                 flat
               >
                   <template v-slot:body-cell-section="props">
                     <q-td :props="props">
                       <div class="text-weight-bold">{{ props.row.name }}</div>
                     </q-td>
                   </template>
                   <template v-slot:body-cell-actions="props">
                     <q-td :props="props" auto-width>
                        <q-btn flat round size="sm" icon="edit" color="primary" @click="editClass(props.row)" />
                        <q-btn flat round size="sm" icon="delete" color="negative" @click="confirmDeleteClass(props.row)" />
                     </q-td>
                   </template>
               </q-table>
             </q-tab-panel>

             <q-tab-panel name="users">
                <div class="row items-center justify-between q-mb-md">
                    <q-option-group
                        v-model="userRoleFilter"
                        :options="[
                            { label: 'Tutti', value: null },
                            { label: 'Docenti', value: 'teacher' },
                            { label: 'Studenti', value: 'student' },
                            { label: 'Genitori', value: 'parent' }
                        ]"
                        color="primary"
                        inline
                    />
                    <q-btn 
                        color="primary" 
                        icon="person_add" 
                        label="Nuovo Utente" 
                        size="sm" 
                        unelevated 
                        @click="openUserDialog"
                    />
                </div>
               <q-table
                 :rows="users"
                 :columns="userColumns"
                 row-key="id"
                 :loading="loadingUsers"
                 flat
               >
                   <template v-slot:body-cell-active="props">
                     <q-td :props="props" class="text-center">
                       <q-badge :color="props.row.is_active ? 'positive' : 'grey'">
                         {{ props.row.is_active ? 'Sì' : 'No' }}
                       </q-badge>
                     </q-td>
                   </template>
                   <template v-slot:body-cell-actions="props">
                     <q-td :props="props" auto-width>
                        <q-btn flat round size="sm" icon="edit" color="primary" @click="editUser(props.row)" />
                        <q-btn flat round size="sm" icon="delete" color="negative" @click="confirmDeleteUser(props.row)" />
                     </q-td>
                   </template>
               </q-table>
             </q-tab-panel>
           </q-tab-panels>
         </q-card>

         <!-- Class Dialog -->
         <q-dialog v-model="showClassDialog">
            <q-card style="min-width: 400px">
                <q-card-section>
                    <div class="text-h6">{{ editingClass ? 'Modifica Classe' : 'Nuova Classe' }}</div>
                </q-card-section>
                <q-card-section>
                    <q-form @submit="saveClass" class="q-gutter-md">
                        <q-input 
                            v-model="classForm.name" 
                            label="Nome (es. 1A)" 
                            outlined dense 
                            :rules="[val => !!val || 'Campo obbligatorio']" 
                        />
                         <q-input 
                            v-model="classForm.section" 
                            label="Sezione (es. A)" 
                            outlined dense 
                        />
                        <q-input 
                            v-model="classForm.academic_year" 
                            label="Anno Scolastico (es. 2024/2025)" 
                            outlined dense 
                            :rules="[val => !!val || 'Campo obbligatorio']" 
                        />
                        <div class="row justify-end q-gutter-sm">
                            <q-btn flat label="Annulla" color="grey" v-close-popup />
                            <q-btn type="submit" :label="editingClass ? 'Salva' : 'Crea'" color="primary" />
                        </div>
                    </q-form>
                </q-card-section>
            </q-card>
         </q-dialog>

         <!-- User Dialog -->
         <q-dialog v-model="showUserDialog">
            <q-card style="min-width: 500px">
                <q-card-section>
                    <div class="text-h6">{{ editingUser ? 'Modifica Utente' : 'Nuovo Utente' }}</div>
                </q-card-section>
                <q-card-section>
                    <q-form @submit="saveUser" class="q-gutter-y-md">
                        <div class="row q-col-gutter-md">
                             <div class="col-6">
                                <q-input 
                                    v-model="userForm.first_name" 
                                    label="Nome" 
                                    outlined dense 
                                    :rules="[val => !!val || 'Campo obbligatorio']" 
                                />
                             </div>
                             <div class="col-6">
                                <q-input 
                                    v-model="userForm.last_name" 
                                    label="Cognome" 
                                    outlined dense 
                                    :rules="[val => !!val || 'Campo obbligatorio']" 
                                />
                             </div>
                        </div>
                        <q-input 
                            v-model="userForm.email" 
                            label="Email" 
                            type="email"
                            outlined dense 
                            :rules="[val => !!val || 'Campo obbligatorio']" 
                        />
                         <q-select 
                            v-if="!editingUser"
                            v-model="userForm.role" 
                            :options="roleOptions"
                            label="Ruolo" 
                            outlined dense 
                            emit-value map-options
                            :rules="[val => !!val || 'Campo obbligatorio']" 
                        />
                         <q-input 
                            v-model="userForm.fiscal_code" 
                            label="Codice Fiscale" 
                            outlined dense 
                            :rules="[val => !!val || 'Campo obbligatorio', val => val.length === 16 || 'Deve essere 16 caratteri']" 
                            uppercase
                        />
                        <q-input
                             v-if="!editingUser"
                             v-model="userForm.password"
                             label="Password Provvisoria"
                             outlined dense
                             type="password"
                             :rules="[val => !!val || 'Campo obbligatorio', val => val.length >= 8 || 'Minimo 8 caratteri']"
                        />
                        
                        <div class="row justify-end q-gutter-sm">
                            <q-btn flat label="Annulla" color="grey" v-close-popup />
                            <q-btn type="submit" :label="editingUser ? 'Salva' : 'Crea'" color="primary" />
                        </div>
                    </q-form>
                </q-card-section>
            </q-card>
         </q-dialog>

    </div>
    
    <div v-else class="text-center q-pa-xl">
        <q-icon name="error_outline" size="4em" color="negative" />
        <div class="text-h5 q-mt-md">Scuola non trovata</div>
        <q-btn label="Torna indietro" color="primary" flat @click="$router.back()" class="q-mt-md" />
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import adminService from '@/services/adminService'

const route = useRoute()
const $q = useQuasar()

const school = ref(null)
const loading = ref(true)
const tab = ref('classes')

// Classes Data
const classes = ref([])
const loadingClasses = ref(false)
const classColumns = [
    { name: 'section', label: 'Classe', field: 'section', align: 'left', sortable: true },
    { name: 'year', label: 'Anno', field: 'academic_year', align: 'left', sortable: true },
    { name: 'students', label: 'Studenti', field: val => val.students_count || 0, align: 'center' },
    { name: 'actions', label: 'Azioni', align: 'center' }
]

const showClassDialog = ref(false)
const editingClass = ref(null)
const classForm = ref({ name: '', section: '', academic_year: '' })

const openClassDialog = () => {
    editingClass.value = null
    classForm.value = { name: '', section: '', academic_year: '2024/2025', school_id: school.value.id }
    showClassDialog.value = true
}

const editClass = (row) => {
    editingClass.value = row
    classForm.value = { ...row }
    showClassDialog.value = true
}

const saveClass = async () => {
    try {
        if (editingClass.value) {
            await adminService.updateClass(editingClass.value.id, classForm.value)
            $q.notify({ type: 'positive', message: 'Classe aggiornata' })
        } else {
            await adminService.createClass(classForm.value)
             $q.notify({ type: 'positive', message: 'Classe creata' })
        }
        showClassDialog.value = false
        fetchClasses(school.value.id)
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore salvataggio classe', caption: e.message })
    }
}

const confirmDeleteClass = (row) => {
    $q.dialog({
        title: 'Elimina Classe',
        message: `Sei sicuro di voler eliminare la classe ${row.name}?`,
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            await adminService.deleteClass(row.id)
            $q.notify({ type: 'positive', message: 'Classe eliminata' })
            fetchClasses(school.value.id)
        } catch (e) {
            $q.notify({ type: 'negative', message: 'Errore eliminazione classe' })
        }
    })
}

// Users Data
const users = ref([])
const loadingUsers = ref(false)
const userRoleFilter = ref(null)
const userColumns = [
    { name: 'name', label: 'Nome', field: row => `${row.first_name} ${row.last_name}`, align: 'left', sortable: true },
    { name: 'email', label: 'Email', field: 'email', align: 'left', sortable: true },
    { name: 'role', label: 'Ruolo', field: 'role', align: 'left', sortable: true },
    { name: 'active', label: 'Attivo', field: 'is_active', align: 'center' },
    { name: 'actions', label: 'Azioni', align: 'center' }
]

const showUserDialog = ref(false)
const editingUser = ref(null)
const userForm = ref({ first_name: '', last_name: '', email: '', role: 'student', fiscal_code: '', password: '' })
const roleOptions = [
    { label: 'Studente', value: 'student' },
    { label: 'Docente', value: 'teacher' },
    { label: 'Genitore', value: 'parent' },
    { label: 'Segreteria', value: 'secretary' },
    { label: 'Preside', value: 'principal' }
]

const openUserDialog = () => {
    editingUser.value = null
    userForm.value = { 
        first_name: '', last_name: '', email: '', role: 'student', 
        fiscal_code: '', password: '', school_id: school.value.id 
    }
    showUserDialog.value = true
}

const editUser = (row) => {
    editingUser.value = row
    userForm.value = { ...row }
    showUserDialog.value = true
}

const saveUser = async () => {
    try {
        if (editingUser.value) {
            await adminService.updateUser(editingUser.value.id, userForm.value)
            $q.notify({ type: 'positive', message: 'Utente aggiornato' })
        } else {
            await adminService.createUser(userForm.value)
             $q.notify({ type: 'positive', message: 'Utente creato' })
        }
        showUserDialog.value = false
        fetchUsers(school.value.id)
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore salvataggio utente', caption: e.response?.data?.error || e.message })
    }
}

const confirmDeleteUser = (row) => {
    $q.dialog({
        title: 'Elimina Utente',
        message: `Sei sicuro di voler eliminare ${row.first_name} ${row.last_name}?`,
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            await adminService.deleteUser(row.id)
            $q.notify({ type: 'positive', message: 'Utente eliminato' })
            fetchUsers(school.value.id)
        } catch (e) {
            $q.notify({ type: 'negative', message: 'Errore eliminazione utente' })
        }
    })
}

const fetchSchool = async () => {
    loading.value = true
    try {
        const id = route.params.id
        const response = await adminService.getSchool(id)
        school.value = response.data
        // Load additional data
        fetchClasses(id)
        fetchUsers(id) // Pre-load or load lazy? Pre-load fine for now
    } catch (error) {
        $q.notify({
            type: 'negative',
            message: 'Errore nel caricamento della scuola',
            caption: error.response?.data?.message || error.message
        })
    } finally {
        loading.value = false
    }
}

const fetchClasses = async (schoolId) => {
    loadingClasses.value = true
    try {
        const response = await adminService.getSchoolClasses(schoolId || school.value.id)
        classes.value = response.data
    } catch (e) {
        console.error("Error fetching classes", e)
    } finally {
        loadingClasses.value = false
    }
}

const fetchUsers = async (schoolId) => {
    loadingUsers.value = true
    try {
        const id = schoolId && typeof schoolId === 'string' ? schoolId : school.value.id
        const response = await adminService.getSchoolUsers(id, userRoleFilter.value)
        // Backend returns wrapped response { users: [], ... } for /users endpoint
        users.value = response.data.users || response.data
    } catch (e) {
         console.error("Error fetching users", e)
    } finally {
        loadingUsers.value = false
    }
}

// Watch role filter to refresh users
watch(userRoleFilter, () => {
    if (school.value) fetchUsers(school.value.id)
})

onMounted(() => {
    fetchSchool()
})
</script>
