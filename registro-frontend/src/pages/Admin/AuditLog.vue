<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="col">
        <div class="text-h4 text-weight-bold">Audit Logs</div>
        <div class="text-subtitle1 text-grey-7">Monitoraggio attività di sistema</div>
      </div>
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-md">
           <!-- Date Range Filters can come here later -->
          <div class="col-12 col-md-4">
             <q-select
              v-model="filters.action"
              :options="actionOptions"
              label="Tipo Azione"
              dense
              outlined
              clearable
              emit-value
              map-options
              @update:model-value="fetchLogs"
            />
          </div>
           <div class="col-12 col-md-4">
            <q-btn
              outline
              color="primary"
              icon="refresh"
              label="Aggiorna"
              @click="fetchLogs"
              :loading="loading"
              class="full-width"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Logs Table -->
    <q-card>
      <q-table
        :rows="logs"
        :columns="columns"
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        @request="onRequest"
      >
        <template v-slot:body="props">
          <q-tr :props="props">
            <q-td key="created_at" :props="props">
              {{ new Date(props.row.created_at).toLocaleString() }}
            </q-td>
            <q-td key="admin_name" :props="props">
              <div class="text-weight-bold">{{ props.row.admin_name }}</div>
              <div class="text-caption text-grey">{{ props.row.admin_id }}</div>
            </q-td>
            <q-td key="action_type" :props="props">
              <q-chip :color="getActionColor(props.row.action_type)" text-color="white" size="sm">
                {{ props.row.action_type.toUpperCase() }}
              </q-chip>
            </q-td>
            <q-td key="target" :props="props">
              {{ props.row.target }}
              <span v-if="props.row.target_id" class="text-grey-6 text-caption">
                (#{{ props.row.target_id.substring(0,8) }})
              </span>
            </q-td>
            <q-td key="school_name" :props="props">
              {{ props.row.school_name || '-' }}
            </q-td>
            <q-td key="details" :props="props">
              {{ props.row.details }}
            </q-td>
          </q-tr>
        </template>
      </q-table>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import adminService from '@/services/adminService'

const $q = useQuasar()
const logs = ref([])
const loading = ref(false)

const filters = reactive({
  action: null
})

const pagination = ref({
  page: 1,
  rowsPerPage: 20,
  rowsNumber: 0
})

const columns = [
  { name: 'created_at', label: 'Data/Ora', align: 'left', field: 'created_at', sortable: true },
  { name: 'admin_name', label: 'Admin', align: 'left', field: 'admin_name' },
  { name: 'action_type', label: 'Azione', align: 'center', field: 'action_type' },
  { name: 'target', label: 'Target', align: 'left', field: 'target' },
  { name: 'school_name', label: 'Scuola', align: 'left', field: 'school_name' },
  { name: 'details', label: 'Dettagli', align: 'left', field: 'details' }
]

const actionOptions = [
  { label: 'Tutti', value: null },
  { label: 'Create', value: 'create' },
  { label: 'Update', value: 'update' },
  { label: 'Delete', value: 'delete' },
  { label: 'Login', value: 'login' }
]

const getActionColor = (action) => {
  switch(action) {
    case 'create': return 'positive'
    case 'update': return 'warning'
    case 'delete': return 'negative'
    case 'login': return 'info'
    default: return 'grey'
  }
}

const fetchLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.rowsPerPage,
      action: filters.action || undefined
    }
    const response = await adminService.getAuditLogs(params)
    logs.value = response.data.items
    pagination.value.rowsNumber = response.data.total
  } catch (error) {
    $q.notify({
      type: 'negative',
      message: 'Errore caricamento logs',
      caption: error.message
    })
  } finally {
    loading.value = false
  }
}

const onRequest = (props) => {
  pagination.value = props.pagination
  fetchLogs()
}

onMounted(() => {
  fetchLogs()
})
</script>
