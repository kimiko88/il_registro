<template>
  <q-page class="q-pa-md" @keydown.ctrl.k.prevent="focusSearch">
    <!-- Header -->
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h3 text-weight-bold text-outfit text-gradient-premium q-my-none" style="display: inline-block;">
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
          @click="openCreate"
        />
      </div>
    </div>

    <!-- Filters -->
    <q-card class="glass-card q-mb-xl shadow-soft">
      <q-card-section class="q-pa-lg">
        <div class="row q-col-gutter-lg">
          <div class="col-12 col-md-5">
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
              v-model="filters.tier"
              :options="tierFilterOptions"
              label="Ordine Scolastico"
              outlined
              bg-color="white"
              clearable
              emit-value
              map-options
              @update:model-value="fetchSchools"
            />
          </div>
          <div class="col-12 col-md-2">
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
          <div class="col-12 col-md-2">
            <q-btn
              unelevated
              color="indigo-50"
              text-color="indigo-700"
              icon="refresh"
              label="Aggiorna"
              @click="fetchSchools"
              :loading="loading"
              class="full-width rounded-lg"
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
            <template v-if="props.row">
              <div class="text-weight-medium">{{ props.row.name }}</div>
              <div class="text-caption text-grey-7">{{ props.row.code }}</div>
            </template>
          </q-td>
        </template>

        <template v-slot:body-cell-location="props">
          <q-td :props="props">
            <div v-if="props.row">{{ props.row.city }}, {{ props.row.province }}</div>
            <div v-if="props.row" class="text-caption text-grey-7">{{ props.row.address }}</div>
          </q-td>
        </template>

        <template v-slot:body-cell-school_level="props">
          <q-td :props="props" align="center">
            <q-badge
              v-if="props.row"
              :color="getSchoolTierBadgeColor(props.row.school_level || props.row.type)"
              class="q-px-sm q-py-xs text-weight-bold"
            >
              {{ getSchoolTierLabel(props.row.school_level || props.row.type) }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-students="props">
          <q-td :props="props">
            <q-badge v-if="props.row" color="amber">{{ props.row.student_count || 0 }}</q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-teachers="props">
          <q-td :props="props">
            <q-badge v-if="props.row" color="purple">{{ props.row.teacher_count || 0 }}</q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-is_active="props">
          <q-td :props="props">
            <q-badge v-if="props.row" :color="props.row.is_active ? 'positive' : 'negative'">
              {{ props.row.is_active ? 'Attiva' : 'Disattiva' }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <div v-if="props.row" class="q-gutter-xs">
              <q-btn
                flat
                dense
                round
                icon="visibility"
                aria-label="Visualizza scuola"
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
                aria-label="Modifica scuola"
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
                aria-label="Elimina scuola"
                @click="confirmDelete(props.row)"
              >
                <q-tooltip>Elimina</q-tooltip>
              </q-btn>
            </div>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Create/Edit Dialog -->
    <q-dialog v-model="showCreateDialog" persistent>
      <q-card style="width: min(600px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6">{{ editingSchool ? 'Modifica Scuola' : 'Nuova Scuola' }}</div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
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
            <q-select
              v-model="schoolForm.school_level"
              :options="tierSelectOptions"
              label="Ordine Scolastico (Normativa Ministeriale) *"
              outlined
              emit-value
              map-options
              :hint="selectedTierHint"
            >
              <template v-slot:option="scope">
                <q-item v-bind="scope.itemProps">
                  <q-item-section avatar>
                    <q-icon :name="scope.opt.icon" :color="scope.opt.color" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">{{ scope.opt.label }}</q-item-label>
                    <q-item-label caption>{{ scope.opt.caption }}</q-item-label>
                  </q-item-section>
                </q-item>
              </template>
            </q-select>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useQuasar, debounce } from 'quasar'
import { usePermissions } from '@/composables/usePermissions'
import { useTableExport } from '@/composables/useTableExport'
import adminService from '@/services/adminService'

const router = useRouter()
const $q = useQuasar()
const { t } = useI18n()
const {
  isSuperAdmin,
  canCreateSchools,
  canDeleteSchools,
  canEditSchool
} = usePermissions()
const { exportTableCsv } = useTableExport()

const schools = ref([])
const selected = ref([])
const loading = ref(false)
const saving = ref(false)
const showCreateDialog = ref(false)
const editingSchool = ref(null)
const searchInput = ref(null)

const filters = reactive({
  search: '',
  status: null,
  tier: null
})

const pagination = ref({
  page: 1,
  rowsPerPage: 20,
  rowsNumber: 0
})

const schoolForm = reactive({
  name: '',
  code: '',
  school_level: 'secondaria_secondo_grado',
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
    name: 'school_level',
    label: 'Ordine',
    align: 'center',
    field: row => row.school_level || row.type || 'secondaria_secondo_grado',
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

const tierFilterOptions = computed(() => [
  { label: 'Tutti gli ordini', value: null },
  { label: t('schoolLevels.infanzia') || "Scuola dell'Infanzia", value: 'infanzia' },
  { label: t('schoolLevels.primaria') || 'Scuola Primaria', value: 'primaria' },
  { label: t('schoolLevels.secondaria_primo_grado') || 'Secondaria I Grado', value: 'secondaria_primo_grado' },
  { label: t('schoolLevels.secondaria_secondo_grado') || 'Secondaria II Grado', value: 'secondaria_secondo_grado' },
  { label: t('schoolLevels.comprensivo') || 'Istituto Comprensivo', value: 'comprensivo' },
  { label: t('schoolLevels.omnicomprensivo') || 'Istituto Omnicomprensivo', value: 'omnicomprensivo' }
])

const tierSelectOptions = computed(() => [
  {
    label: t('schoolLevels.infanzia') || "Scuola dell'Infanzia (3-6 anni)",
    value: 'infanzia',
    icon: 'child_care',
    color: 'pink-7',
    caption: t('schoolLevels.infanziaDesc') || "Campi d'Esperienza (D.M. 254/2012) e osservazioni senza voti"
  },
  {
    label: t('schoolLevels.primaria') || 'Scuola Primaria (6-11 anni)',
    value: 'primaria',
    icon: 'menu_book',
    color: 'amber-9',
    caption: t('schoolLevels.primariaDesc') || 'Obiettivi di apprendimento e 4 livelli di giudizio (O.M. 172/2020)'
  },
  {
    label: t('schoolLevels.secondaria_primo_grado') || 'Scuola Secondaria I Grado (11-14 anni)',
    value: 'secondaria_primo_grado',
    icon: 'school',
    color: 'blue-8',
    caption: t('schoolLevels.secondariaPrimoDesc') || 'Voti decimali 1-10, INVALSI e orientamento superiore'
  },
  {
    label: t('schoolLevels.secondaria_secondo_grado') || 'Scuola Secondaria II Grado (14-19 anni)',
    value: 'secondaria_secondo_grado',
    icon: 'account_balance',
    color: 'indigo-8',
    caption: t('schoolLevels.secondariaSecondoDesc') || 'Voti 1-10, crediti triennio, PCTO e Scrutinio Differito (O.M. 92/2007)'
  },
  {
    label: t('schoolLevels.comprensivo') || 'Istituto Comprensivo (Infanzia + Primaria + I Grado)',
    value: 'comprensivo',
    icon: 'hub',
    color: 'teal-8',
    caption: t('schoolLevels.comprensivoDesc') || 'Gestione integrata primo ciclo con continuità pedagogica'
  },
  {
    label: t('schoolLevels.omnicomprensivo') || 'Istituto Omnicomprensivo (Tutti gli ordini)',
    value: 'omnicomprensivo',
    icon: 'apartment',
    color: 'purple-8',
    caption: t('schoolLevels.omnicomprensivoDesc') || 'Gestione completa dall\'Infanzia alle Scuole Superiori'
  }
])

const selectedTierHint = computed(() => {
  const selected = tierSelectOptions.value.find(o => o.value === schoolForm.school_level)
  return selected ? selected.caption : ''
})

const getSchoolTierLabel = (tier) => {
  const map = {
    infanzia: t('schoolLevels.infanzia') || "Scuola dell'Infanzia",
    primaria: t('schoolLevels.primaria') || 'Scuola Primaria',
    secondaria_primo_grado: t('schoolLevels.secondaria_primo_grado') || 'Secondaria I Grado',
    secondaria_secondo_grado: t('schoolLevels.secondaria_secondo_grado') || 'Secondaria II Grado',
    comprensivo: t('schoolLevels.comprensivo') || 'Ist. Comprensivo',
    omnicomprensivo: t('schoolLevels.omnicomprensivo') || 'Ist. Omnicomprensivo'
  }
  return map[tier] || tier || 'Secondaria II Grado'
}

const getSchoolTierBadgeColor = (tier) => {
  const map = {
    infanzia: 'pink-7',
    primaria: 'amber-9',
    secondaria_primo_grado: 'blue-8',
    secondaria_secondo_grado: 'indigo-8',
    comprensivo: 'teal-8',
    omnicomprensivo: 'purple-8'
  }
  return map[tier] || 'indigo-8'
}

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
    let items = response.data.items || []
    if (filters.tier) {
      items = items.filter(s => (s.school_level || s.type) === filters.tier)
    }
    schools.value = items
    pagination.value.rowsNumber = filters.tier ? items.length : response.data.total
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
  Object.assign(schoolForm, {
    name: school.name,
    code: school.code,
    school_level: school.school_level || school.type || 'secondaria_secondo_grado',
    address: school.address,
    city: school.city,
    province: school.province,
    zip_code: school.zip_code,
    phone: school.phone,
    email: school.email,
    website: school.website,
    is_active: school.is_active
  })
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
        let successCount = 0
        let failCount = 0
        const failedNames = []

        for (const s of selected.value) {
            try {
                await adminService.deleteSchool(s.id)
                successCount++
            } catch (e) {
                console.error('Error removing school', s.name, e)
                failCount++
                failedNames.push(s.name)
            }
        }
        selected.value = []
        fetchSchools()
        if (failCount === 0) {
            $q.notify({ type: 'positive', message: `${successCount} scuole eliminate con successo` })
        } else if (successCount > 0) {
            $q.notify({
                type: 'warning',
                message: `${successCount} scuole eliminate, ${failCount} non eliminate (${failedNames.join(', ')})`
            })
        } else {
            $q.notify({ type: 'negative', message: `Impossibile eliminare le scuole selezionate (${failedNames.join(', ')})` })
        }
    })
}

// Export CSV
const exportTable = () => {
  exportTableCsv({
    filename: 'scuole-export.csv',
    columns,
    rows: schools.value
  })
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
    school_level: 'secondaria_secondo_grado',
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
    school_level: 'secondaria_secondo_grado',
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


