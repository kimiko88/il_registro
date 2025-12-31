<template>
  <q-page class="q-pa-md">
    <!-- Header -->
    <div class="row items-center q-mb-md">
      <div class="col">
        <div class="text-h4 text-weight-bold">Gestione Admin Scuole</div>
        <div class="text-subtitle1 text-grey-7">
          Gestisci gli amministratori delle singole scuole
        </div>
      </div>
      <div class="col-auto">
        <q-btn
          color="primary"
          icon="add"
          label="Nuovo Admin"
          @click="openCreateDialog"
        />
      </div>
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <q-input
              v-model="filters.search"
              placeholder="Cerca per nome o email..."
              dense
              outlined
              clearable
              @update:model-value="debouncedFetch"
            >
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
          <div class="col-12 col-md-3">
             <q-select
              v-model="filters.school_id"
              :options="schoolOptions"
              label="Filtra per Scuola"
              option-label="name"
              option-value="id"
              emit-value
              map-options
              dense
              outlined
              clearable
              @update:model-value="fetchAdmins"
            />
          </div>
          <div class="col-12 col-md-3">
            <q-btn
              outline
              color="primary"
              icon="refresh"
              label="Aggiorna"
              @click="fetchAdmins"
              :loading="loading"
              class="full-width"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Admins Table -->
    <q-card>
      <q-table
        :rows="admins"
        :columns="columns"
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        @request="onRequest"
        binary-state-sort
      >
        <template v-slot:body-cell-user="props">
          <q-td :props="props">
            <div class="text-weight-medium">{{ props.row.first_name }} {{ props.row.last_name }}</div>
            <div class="text-caption text-grey-7">{{ props.row.email }}</div>
          </q-td>
        </template>

        <template v-slot:body-cell-school="props">
          <q-td :props="props">
             <div v-if="props.row.school_name">
                <q-icon name="school" size="xs" color="primary" class="q-mr-xs" />
                {{ props.row.school_name }}
             </div>
             <div v-else class="text-grey-5 text-italic">Non assegnata</div>
          </q-td>
        </template>

        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="props.row.is_active ? 'positive' : 'negative'">
              {{ props.row.is_active ? 'Attivo' : 'Disattivo' }}
            </q-badge>
          </q-td>
        </template>
        
        <template v-slot:body-cell-last_login="props">
          <q-td :props="props">
            {{ formatDate(props.row.last_login_at) }}
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn-dropdown flat dense round icon="more_vert">
              <q-list>
                <q-item clickable v-close-popup @click="editAdmin(props.row)">
                  <q-item-section avatar>
                    <q-icon name="edit" color="primary" />
                  </q-item-section>
                  <q-item-section>Modifica</q-item-section>
                </q-item>
                
                <q-item clickable v-close-popup @click="viewActivity(props.row)">
                  <q-item-section avatar>
                    <q-icon name="history" color="info" />
                  </q-item-section>
                  <q-item-section>Log Attività</q-item-section>
                </q-item>

                <q-separator />

                <q-item clickable v-close-popup @click="confirmResetPassword(props.row)">
                  <q-item-section avatar>
                    <q-icon name="lock_reset" color="warning" />
                  </q-item-section>
                  <q-item-section>Reset Password</q-item-section>
                </q-item>

                <q-item clickable v-close-popup @click="confirmDelete(props.row)">
                  <q-item-section avatar>
                    <q-icon name="delete" color="negative" />
                  </q-item-section>
                  <q-item-section>Elimina</q-item-section>
                </q-item>
              </q-list>
            </q-btn-dropdown>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Create/Edit Dialog -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ editingAdmin ? 'Modifica Admin' : 'Nuovo Admin' }}</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="saveAdmin" class="q-gutter-md">
            <div class="row q-col-gutter-md">
              <div class="col-6">
                <q-input
                  v-model="adminForm.first_name"
                  label="Nome *"
                  outlined
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
              <div class="col-6">
                <q-input
                  v-model="adminForm.last_name"
                  label="Cognome *"
                  outlined
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
            </div>

            <q-input
              v-model="adminForm.email"
              label="Email *"
              type="email"
              outlined
              :rules="[val => !!val || 'Campo obbligatorio', val => /.+@.+\..+/.test(val) || 'Email non valida']"
            />
            
            <q-select
              v-model="adminForm.school_id"
              :options="schoolOptions"
              label="Scuola Assegnata *"
              option-label="name"
              option-value="id"
              emit-value
              map-options
              outlined
              :rules="[val => !!val || 'Campo obbligatorio']"
            >
              <template v-slot:no-option>
                <q-item>
                  <q-item-section class="text-grey">
                    Nessuna scuola trovata
                  </q-item-section>
                </q-item>
              </template>
            </q-select>

            <q-toggle
              v-if="editingAdmin"
              v-model="adminForm.is_active"
              label="Account Attivo"
            />
            
            <div v-if="!editingAdmin" class="text-caption text-grey-7 q-mt-sm">
              <q-icon name="info" /> 
              Verrà inviata una email all'utente per impostare la password.
            </div>

            <div class="row justify-end q-gutter-sm q-mt-lg">
              <q-btn flat label="Annulla" color="grey" v-close-popup />
              <q-btn type="submit" :label="editingAdmin ? 'Salva' : 'Crea'" color="primary" :loading="saving" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Activity Log Dialog -->
     <q-dialog v-model="showActivityDialog">
      <q-card style="min-width: 600px; max-width: 900px">
        <q-card-section class="row items-center">
          <div class="text-h6">Log Attività: {{ selectedAdminForActivity?.first_name }} {{ selectedAdminForActivity?.last_name }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-none">
          <q-table
            :rows="activityLogs"
            :columns="activityColumns"
            row-key="id"
            flat
          >
             <template v-slot:body-cell-created_at="props">
                <q-td :props="props">
                  {{ formatDateTime(props.row.created_at) }}
                </q-td>
              </template>
          </q-table>
        </q-card-section>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar, debounce } from 'quasar'
import adminService from '@/services/adminService'

const $q = useQuasar()

// State
const admins = ref([])
const schools = ref([])
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const editingAdmin = ref(null)
const schoolOptions = ref([])

// Activity Log State
const showActivityDialog = ref(false)
const activityLogs = ref([])
const selectedAdminForActivity = ref(null)

const filters = reactive({
  search: '',
  school_id: null
})

const pagination = ref({
  page: 1,
  rowsPerPage: 20,
  rowsNumber: 0
})

const adminForm = reactive({
  first_name: '',
  last_name: '',
  email: '',
  school_id: null,
  role: 'admin', // Fixed to admin
  is_active: true
})

// Columns
const columns = [
  { name: 'user', label: 'Utente', align: 'left', field: 'email', sortable: true },
  { name: 'school', label: 'Scuola', align: 'left', field: 'school_name', sortable: true },
  { name: 'status', label: 'Stato', align: 'center', field: 'is_active', sortable: true },
  { name: 'last_login', label: 'Ultimo Accesso', align: 'left', field: 'last_login_at', sortable: true },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

const activityColumns = [
  { name: 'action', label: 'Azione', align: 'left', field: 'action_type' },
  { name: 'target', label: 'Target', align: 'left', field: 'target_entity' },
  { name: 'school', label: 'Scuola', align: 'left', field: 'school_name' },
  { name: 'details', label: 'Dettagli', align: 'left', field: 'details' },
  { name: 'created_at', label: 'Data', align: 'left', field: 'created_at', sortable: true }
]

// Methods
const fetchAdmins = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.rowsPerPage,
      search: filters.search || undefined,
      school_id: filters.school_id || undefined
    }
    const response = await adminService.getAdmins(params)
    admins.value = response.data.items || []
    pagination.value.rowsNumber = response.data.total
  } catch (error) {
    $q.notify({ type: 'negative', message: 'Errore caricamento admin', caption: error.message })
  } finally {
    loading.value = false
  }
}

const fetchSchools = async () => {
  try {
    // Get simple list for dropdown
    const response = await adminService.getSchools({ page_size: 100, status: 'active' })
    schoolOptions.value = response.data.items || []
  } catch (error) {
    console.error('Failed to load schools', error)
  }
}

const debouncedFetch = debounce(fetchAdmins, 500)

const onRequest = (props) => {
  pagination.value = props.pagination
  fetchAdmins()
}

// Create/Edit
const openCreateDialog = () => {
  editingAdmin.value = null
  Object.assign(adminForm, {
    first_name: '', last_name: '', email: '', school_id: null, role: 'admin', is_active: true
  })
  showDialog.value = true
}

const editAdmin = (row) => {
  editingAdmin.value = row
  Object.assign(adminForm, {
    first_name: row.first_name,
    last_name: row.last_name,
    email: row.email,
    school_id: row.school_id,
    role: 'admin',
    is_active: row.is_active
  })
  showDialog.value = true
}

const saveAdmin = async () => {
  saving.value = true
  try {
    if (editingAdmin.value) {
      await adminService.updateAdmin(editingAdmin.value.id, adminForm)
      $q.notify({ type: 'positive', message: 'Admin aggiornato' })
    } else {
      await adminService.createAdmin(adminForm)
      $q.notify({ type: 'positive', message: 'Admin creato con successo' })
    }
    showDialog.value = false
    fetchAdmins()
  } catch (error) {
    $q.notify({ type: 'negative', message: 'Errore salvataggio', caption: error.response?.data?.message || error.message })
  } finally {
    saving.value = false
  }
}

// Actions
const confirmDelete = (row) => {
  $q.dialog({
    title: 'Conferma eliminazione',
    message: `Sei sicuro di voler eliminare l'admin ${row.first_name} ${row.last_name}?`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: 'Elimina' }
  }).onOk(async () => {
    try {
      await adminService.deleteAdmin(row.id)
      $q.notify({ type: 'positive', message: 'Admin eliminato' })
      fetchAdmins()
    } catch (error) {
      $q.notify({ type: 'negative', message: 'Errore eliminazione' })
    }
  })
}

const confirmResetPassword = (row) => {
   $q.dialog({
    title: 'Reset Password',
    message: `Inviare email di reset password a ${row.email}?`,
    cancel: true
  }).onOk(async () => {
    try {
      await adminService.resetAdminPassword(row.id)
      $q.notify({ type: 'positive', message: 'Email di reset inviata' })
    } catch (error) {
      $q.notify({ type: 'negative', message: 'Errore reset password' })
    }
  })
}

const viewActivity = async (row) => {
  selectedAdminForActivity.value = row
  showActivityDialog.value = true
  try {
    const response = await adminService.getAdminActivity(row.id)
    activityLogs.value = response.data || []
  } catch (error) {
     $q.notify({ type: 'negative', message: 'Errore caricamento log' })
  }
}

// Utils
const formatDate = (str) => {
  if (!str) return 'Mai'
  return new Date(str).toLocaleDateString('it-IT') + ' ' + new Date(str).toLocaleTimeString('it-IT', {hour: '2-digit', minute:'2-digit'})
}

const formatDateTime = (str) => formatDate(str)

// Init
onMounted(() => {
  fetchSchools()
  fetchAdmins()
})
</script>
