<template>
  <q-page class="q-pa-md" @keydown.ctrl.k.prevent="focusSearch">
    <!-- Header -->
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h3 text-weight-bold text-outfit bg-clip-text text-transparent bg-gradient-premium q-my-none" style="display: inline-block;">
          {{ isSuperAdmin ? 'Gestione Scuole' : 'La Mia Scuola' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">
          {{ pagination.rowsNumber || 0 }} istituti registrati nel sistema
          <span v-if="selected.length > 0" class="text-indigo-600 text-weight-bold q-ml-md bg-indigo-50 q-px-sm rounded-lg">
            {{ selected.length }} selezionate
          </span>
        </div>
      </div>
      <div class="col-auto q-gutter-sm">
        <q-btn
            v-if="selected.length > 0 && canDeleteSchools"
            color="negative"
            icon="delete"
            label="Elimina Selezione"
            unelevated
            class="rounded-lg shadow-soft"
            @click="deleteSelected"
        />
        <q-btn
            color="white"
            text-color="grey-7"
            icon="file_download"
            label="Esporta"
            unelevated
            class="rounded-lg shadow-soft"
            @click="exportTable"
        />
        <q-btn
          v-if="canCreateSchools"
          color="primary"
          icon="add"
          label="Nuova Scuola"
          unelevated
          class="rounded-lg shadow-soft q-px-lg"
          @click="showCreateDialog = true"
        />
      </div>
    </div>

    <!-- Filters -->
    <q-card class="glass-card q-mb-xl shadow-soft">
      <q-card-section class="q-pa-lg">
        <div class="row q-col-gutter-lg">
          <div class="col-12 col-md-6">
            <q-input
              ref="searchInput"
              v-model="filters.search"
              placeholder="Cerca scuola... (Ctrl+K)"
              outlined
              bg-color="white"
              clearable
              @update:model-value="debouncedFetch"
            >
              <template v-slot:prepend>
                <q-icon name="search" color="primary" />
              </template>
            </q-input>
          </div>
          <div class="col-12 col-md-3">
            <q-select
              v-model="filters.status"
              :options="statusOptions"
              label="Stato Istituto"
              outlined
              bg-color="white"
              clearable
              emit-value
              map-options
              @update:model-value="fetchSchools"
            />
          </div>
          <div class="col-12 col-md-3">
            <q-btn
              unelevated
              color="indigo-50"
              text-color="indigo-700"
              icon="refresh"
              label="Aggiorna"
              @click="fetchSchools"
              :loading="loading"
              class="full-width rounded-lg h-full"
              style="height: 56px"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Schools Table -->
    <q-card class="glass-card shadow-soft overflow-hidden">
      <q-table
        v-model:selected="selected"
        :rows="schools"
        :columns="columns"
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        @request="onRequest"
        :selection="isSuperAdmin ? 'multiple' : 'none'"
        binary-state-sort
        flat
        class="bg-transparent"
      >
        <template v-slot:body-cell-name="props">
          <q-td :props="props">
            <div class="text-weight-medium">{{ props.row.name }}</div>
            <div class="text-caption text-grey-7">{{ props.row.code }}</div>
          </q-td>
        </template>

        <template v-slot:body-cell-location="props">
          <q-td :props="props">
            <div>{{ props.row.city }}, {{ props.row.province }}</div>
            <div class="text-caption text-grey-7">{{ props.row.address }}</div>
          </q-td>
        </template>

        <template v-slot:body-cell-students="props">
          <q-td :props="props">
            <q-badge color="amber">{{ props.row.student_count || 0 }}</q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-teachers="props">
          <q-td :props="props">
            <q-badge color="purple">{{ props.row.teacher_count || 0 }}</q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-is_active="props">
          <q-td :props="props">
            <q-badge :color="props.row.is_active ? 'positive' : 'negative'">
              {{ props.row.is_active ? 'Attiva' : 'Disattiva' }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn
              flat
              dense
              round
              icon="visibility"
              @click="viewSchool(props.row)"
            >
              <q-tooltip>Visualizza</q-tooltip>
            </q-btn>
            <q-btn
              v-if="canEditSchool(props.row.id)"
              flat
              dense
              round
              icon="edit"
              @click="editSchool(props.row)"
            >
              <q-tooltip>Modifica</q-tooltip>
            </q-btn>
            <q-btn
              v-if="canDeleteSchools"
              flat
              dense
              round
              icon="delete"
              color="negative"
              @click="confirmDelete(props.row)"
            >
              <q-tooltip>Elimina</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Create/Edit Dialog -->
    <q-dialog v-model="showCreateDialog" persistent>
      <q-card style="min-width: 600px">
        <q-card-section>
          <div class="text-h6">{{ editingSchool ? 'Modifica Scuola' : 'Nuova Scuola' }}</div>
        </q-card-section>

        <q-card-section>
          <div class="q-gutter-y-md">
            <q-input
              v-model="schoolForm.name"
              label="Nome Scuola *"
              outlined
              :rules="[val => !!val || 'Campo obbligatorio']"
            />
            <q-input
              v-model="schoolForm.code"
              label="Codice Meccanografico *"
              outlined
              :rules="[val => !!val || 'Campo obbligatorio']"
            />
            <q-input
              v-model="schoolForm.address"
              label="Indirizzo *"
              outlined
              :rules="[val => !!val || 'Campo obbligatorio']"
            />
            <div class="row q-col-gutter-md">
              <div class="col-8">
                <q-input
                  v-model="schoolForm.city"
                  label="Città *"
                  outlined
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
              <div class="col-4">
                <q-input
                  v-model="schoolForm.province"
                  label="Provincia *"
                  outlined
                  maxlength="2"
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
            </div>
            <q-input
              v-model="schoolForm.zip_code"
              label="CAP *"
              outlined
              :rules="[val => !!val || 'Campo obbligatorio']"
            />
            <q-input
              v-model="schoolForm.phone"
              label="Telefono"
              outlined
            />
            <q-input
              v-model="schoolForm.email"
              label="Email"
              type="email"
              outlined
            />
            <q-input
              v-model="schoolForm.website"
              label="Sito Web"
              outlined
            />
            <q-toggle
              v-if="editingSchool"
              v-model="schoolForm.is_active"
              label="Scuola Attiva"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Annulla" color="grey" @click="closeDialog" />
          <q-btn
            label="Salva"
            color="primary"
            @click="saveSchool"
            :loading="saving"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar, exportFile, debounce } from 'quasar'
import { usePermissions } from '@/composables/usePermissions'
import adminService from '@/services/adminService'

const router = useRouter()
const $q = useQuasar()
const {
  isSuperAdmin,
  canCreateSchools,
  canDeleteSchools,
  canEditSchool
} = usePermissions()

const schools = ref([])
const selected = ref([])
const loading = ref(false)
const saving = ref(false)
const showCreateDialog = ref(false)
const editingSchool = ref(null)
const searchInput = ref(null)

const filters = reactive({
  search: '',
  status: null
})

const pagination = ref({
  page: 1,
  rowsPerPage: 20,
  rowsNumber: 0
})

const schoolForm = reactive({
  name: '',
  code: '',
  address: '',
  city: '',
  province: '',
  zip_code: '',
  phone: '',
  email: '',
  website: '',
  is_active: true
})

const columns = [
  {
    name: 'name',
    label: 'Scuola',
    align: 'left',
    field: 'name',
    sortable: true
  },
  {
    name: 'location',
    label: 'Località',
    align: 'left',
    field: 'city',
    sortable: true
  },
  {
    name: 'students',
    label: 'Studenti',
    align: 'center',
    field: 'student_count',
    sortable: true
  },
  {
    name: 'teachers',
    label: 'Docenti',
    align: 'center',
    field: 'teacher_count',
    sortable: true
  },
  {
    name: 'is_active',
    label: 'Stato',
    align: 'center',
    field: 'is_active',
    sortable: true
  },
  {
    name: 'actions',
    label: 'Azioni',
    align: 'center'
  }
]

const statusOptions = [
  { label: 'Tutte', value: null },
  { label: 'Attive', value: 'active' },
  { label: 'Disattive', value: 'inactive' }
]

const fetchSchools = async () => {
  loading.value = true
  
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.rowsPerPage,
      search: filters.search || undefined,
      status: filters.status || undefined
    }

    const response = await adminService.getSchools(params)
    schools.value = response.data.items
    pagination.value.rowsNumber = response.data.total
  } catch (error) {
    $q.notify({
      type: 'negative',
      message: 'Errore nel caricamento delle scuole',
      caption: error.response?.data?.message || error.message
    })
  } finally {
    loading.value = false
  }
}

const debouncedFetch = debounce(fetchSchools, 500)

const onRequest = (props) => {
  pagination.value = props.pagination
  fetchSchools()
}

const viewSchool = (school) => {
  router.push(`/admin/schools/${school.id}`)
}

const editSchool = (school) => {
  editingSchool.value = school
  Object.assign(schoolForm, school)
  showCreateDialog.value = true
}

const confirmDelete = (school) => {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: `Sei sicuro di voler eliminare la scuola "${school.name}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await adminService.deleteSchool(school.id)
      $q.notify({
        type: 'positive',
        message: 'Scuola eliminata con successo'
      })
      fetchSchools()
    } catch (error) {
      $q.notify({
        type: 'negative',
        message: 'Errore nell\'eliminazione della scuola',
        caption: error.response?.data?.message || error.message
      })
    }
  })
}

// Bulk Delete
const deleteSelected = () => {
    $q.dialog({
        title: 'Conferma Eliminazione Multipla',
        message: `Vuoi eliminare ${selected.value.length} scuole selezionate?`,
        cancel: true,
        persistent: true,
        ok: { color: 'negative', label: 'Elimina Tutti' }
    }).onOk(async () => {
        // Mock bulk delete - iterate one by one for now
        for (const s of selected.value) {
            try {
                await adminService.deleteSchool(s.id)
            } catch (e) {
                console.error('Error removing school', s.name, e)
            }
        }
        selected.value = []
        fetchSchools()
        $q.notify({ type: 'positive', message: 'Elementi eliminati' })
    })
}

// Export CSV
function wrapCsvValue (val, formatFn) {
  let formatted = formatFn !== void 0
    ? formatFn(val)
    : val

  formatted = formatted === void 0 || formatted === null
    ? ''
    : String(formatted)

  formatted = formatted.split('"').join('""')
  return `"${formatted}"`
}

const exportTable = () => {
  const content = [columns.map(col => wrapCsvValue(col.label))].concat(
    schools.value.map(row => columns.map(col => wrapCsvValue(
      typeof col.field === 'function'
        ? col.field(row)
        : row[col.field === void 0 ? col.name : col.field],
      col.format
    )).join(','))
  ).join('\r\n')

  const status = exportFile(
    'scuole-export.csv',
    content,
    'text/csv'
  )

  if (status !== true) {
    $q.notify({
      message: 'Browser denied file download...',
      color: 'negative',
      icon: 'warning'
    })
  }
}

const focusSearch = () => {
    searchInput.value?.focus()
}

const saveSchool = async () => {
  saving.value = true
  
  try {
    if (editingSchool.value) {
      await adminService.updateSchool(editingSchool.value.id, { ...schoolForm })
      $q.notify({
        type: 'positive',
        message: 'Scuola aggiornata con successo'
      })
    } else {
      await adminService.createSchool({ ...schoolForm })
      $q.notify({
        type: 'positive',
        message: 'Scuola creata con successo'
      })
    }
    
    closeDialog()
    fetchSchools()
  } catch (error) {
    $q.notify({
      type: 'negative',
      message: 'Errore nel salvataggio',
      caption: error.response?.data?.message || error.message
    })
  } finally {
    saving.value = false
  }
}

const closeDialog = () => {
  showCreateDialog.value = false
  editingSchool.value = null
  Object.assign(schoolForm, {
    name: '',
    code: '',
    address: '',
    city: '',
    province: '',
    zip_code: '',
    phone: '',
    email: '',
    website: '',
    is_active: true
  })
}

const openCreate = () => {
  editingSchool.value = null
  Object.assign(schoolForm, {
    name: '',
    code: '',
    address: '',
    city: '',
    province: '',
    zip_code: '',
    phone: '',
    email: '',
    website: '',
    is_active: true
  })
  showCreateDialog.value = true
}

defineExpose({
    fetchSchools,
    openCreate,
    editSchool,
    saveSchool,
    confirmDelete,
    deleteSelected,
    exportTable,
    schoolForm,
    showCreateDialog,
    editingSchool,
    filters,
    selected
})

onMounted(() => {
  fetchSchools()
})
</script>


