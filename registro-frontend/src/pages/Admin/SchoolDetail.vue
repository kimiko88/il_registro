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
                      <div class="text-weight-bold">{{ props.row.year }} {{ props.row.section }}</div>
                    </q-td>
                  </template>
               </q-table>
             </q-tab-panel>

             <q-tab-panel name="users">
                <div class="row q-mb-md">
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
                        @update:model-value="fetchUsers"
                    />
                </div>
               <q-table
                 :rows="users"
                 :columns="userColumns"
                 row-key="id"
                 :loading="loadingUsers"
                 flat
               />
             </q-tab-panel>
           </q-tab-panels>
         </q-card>

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
    { name: 'students', label: 'Studenti', field: val => val.students_count || 0, align: 'center' }
]

// Users Data
const users = ref([])
const loadingUsers = ref(false)
const userRoleFilter = ref(null)
const userColumns = [
    { name: 'name', label: 'Nome', field: row => `${row.first_name} ${row.last_name}`, align: 'left', sortable: true },
    { name: 'email', label: 'Email', field: 'email', align: 'left', sortable: true },
    { name: 'role', label: 'Ruolo', field: 'role', align: 'left', sortable: true },
    { name: 'active', label: 'Attivo', field: 'is_active', format: val => val ? 'Sì' : 'No', align: 'center' }
]

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
