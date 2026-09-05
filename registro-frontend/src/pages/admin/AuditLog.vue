<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="col">
        <div class="text-h4 text-weight-bold">{{ t('adminAudit.title') || 'Audit Logs' }}</div>
        <div class="text-subtitle1 text-grey-7">{{ t('adminAudit.subtitle') || 'Monitoraggio attività di sistema' }}</div>
      </div>
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-sm-6 col-md-4">
            <q-select
              v-model="filters.action"
              :options="actionOptions"
              :label="t('adminAudit.actionType')"
              dense
              outlined
              clearable
              emit-value
              map-options
              @update:model-value="fetchLogs"
            />
          </div>
          <div class="col-12 col-sm-6 col-md-4">
            <q-btn
              outline
              color="primary"
              icon="refresh"
              :label="t('adminAudit.refresh')"
              @click="fetchLogs"
              :loading="loading"
              class="full-width"
            />
          </div>
          <div class="col-12 col-sm-12 col-md-4">
            <q-btn
              outline
              color="secondary"
              icon="download"
              :label="t('adminAudit.export')"
              @click="exportAuditLogs"
              :disable="logs.length === 0"
              class="full-width"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Logs Table -->
    <q-card>
      <q-table
        v-model:pagination="pagination"
        :rows="logs"
        :columns="columns"
        row-key="id"
        :loading="loading"
        @request="onRequest"
      >
        <template v-slot:no-data>
          <div class="full-width row flex-center text-grey q-gutter-sm q-py-lg">
            <q-icon size="2em" name="sentiment_dissatisfied" />
            <span>{{ t('adminAudit.noLogs') }}</span>
          </div>
        </template>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useTableExport } from '@/composables/useTableExport'
import adminService from '@/services/adminService'

const $q = useQuasar()
const { t } = useI18n()
const { exportTableCsv } = useTableExport()
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

const exportAuditLogs = () => {
  exportTableCsv({
    filename: `audit_logs_${new Date().toISOString().split('T')[0]}.csv`,
    columns: columns.value,
    rows: logs.value
  })
}

const columns = computed(() => [
  { name: 'created_at', label: t('adminAudit.colDateTime'), align: 'left', field: 'created_at', sortable: true },
  { name: 'admin_name', label: t('adminAudit.colAdmin'), align: 'left', field: 'admin_name' },
  { name: 'action_type', label: t('adminAudit.colAction'), align: 'center', field: 'action_type' },
  { name: 'target', label: t('adminAudit.colTarget'), align: 'left', field: 'target' },
  { name: 'school_name', label: t('adminAudit.colSchool'), align: 'left', field: 'school_name' },
  { name: 'details', label: t('adminAudit.colDetails'), align: 'left', field: 'details' }
])

const actionOptions = computed(() => [
  { label: t('adminAudit.all') || 'Tutti', value: null },
  { label: t('adminAudit.create') || 'Create', value: 'create' },
  { label: t('adminAudit.update') || 'Update', value: 'update' },
  { label: t('adminAudit.delete') || 'Delete', value: 'delete' },
  { label: t('adminAudit.login') || 'Login', value: 'login' }
])

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
    logs.value = response.data.items || []
    const totalCount = response.data.total !== undefined ? Number(response.data.total) : logs.value.length
    pagination.value.rowsNumber = totalCount
  } catch (error) {
    $q.notify({
      type: 'negative',
      message: t('adminAudit.errorLoading') || 'Errore caricamento logs',
      caption: error.message
    })
  } finally {
    loading.value = false
  }
}

const onRequest = (props) => {
  const { page, rowsPerPage } = props.pagination
  pagination.value.page = page
  pagination.value.rowsPerPage = rowsPerPage
  fetchLogs()
}

onMounted(() => {
  fetchLogs()
})

defineExpose({
  fetchLogs,
  onRequest,
  filters,
  pagination,
  logs,
  loading,
  getActionColor
})
</script>
