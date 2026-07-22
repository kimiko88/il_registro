<template>
  <q-page class="q-pa-lg bg-slate-50 min-h-screen">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-ma-none">Audit Log di Sistema</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          Registro di tracciamento e sicurezza per operazioni sensibili
        </p>
      </div>
      <q-btn
        color="secondary"
        icon="download"
        label="Esporta CSV"
        no-caps
        class="shadow-sm rounded-lg q-px-md"
        @click="exportCSV"
      />
    </div>

    <!-- Filter Toolbar -->
    <q-card class="bg-white rounded-xl shadow-sm q-mb-lg">
      <q-card-section class="row q-col-gutter-md items-center">
        <div class="col-12 col-md-3">
          <q-input
            v-model="filters.actorName"
            label="Cerca per utente"
            outlined
            dense
            clearable
            @update:model-value="applyFilters"
          >
            <template #append><q-icon name="search" /></template>
          </q-input>
        </div>
        <div class="col-12 col-md-3">
          <q-select
            v-model="filters.action"
            :options="actionOptions"
            label="Tipo azione"
            outlined
            dense
            emit-value
            map-options
            @update:model-value="applyFilters"
          />
        </div>
        <div class="col-12 col-md-2">
          <q-input
            v-model="filters.from"
            type="date"
            label="Da data"
            outlined
            dense
            stack-label
            @update:model-value="applyFilters"
          />
        </div>
        <div class="col-12 col-md-2">
          <q-input
            v-model="filters.to"
            type="date"
            label="A data"
            outlined
            dense
            stack-label
            @update:model-value="applyFilters"
          />
        </div>
        <div class="col-12 col-md-2 text-right">
          <q-btn flat icon="refresh" label="Aggiorna" no-caps color="primary" @click="loadLogs" />
        </div>
      </q-card-section>
    </q-card>

    <!-- Table -->
    <q-card class="bg-white rounded-xl shadow-sm">
      <q-card-section class="q-pa-none">
        <q-table
          v-model:pagination="pagination"
          :rows="auditStore.logs"
          :columns="columns"
          row-key="id"
          :loading="auditStore.loading"
          flat
          @request="onRequest"
        >
          <template #body="props">
            <q-tr :props="props" :class="getRowClass(props.row.action)">
              <q-td key="created_at" :props="props">
                {{ formatDate(props.row.created_at) }}
              </q-td>
              <q-td key="actor_name" :props="props" class="text-weight-bold">
                {{ props.row.actor_name || props.row.actor_id || 'Sistema' }}
              </q-td>
              <q-td key="actor_role" :props="props">
                <q-chip dense size="sm" color="grey-3" text-color="slate-800" class="text-weight-bold">
                  {{ props.row.actor_role }}
                </q-chip>
              </q-td>
              <q-td key="action" :props="props">
                <q-chip
                  dense
                  size="sm"
                  :color="getActionColor(props.row.action)"
                  text-color="white"
                  class="text-weight-bold"
                >
                  {{ props.row.action }}
                </q-chip>
              </q-td>
              <q-td key="entity_type" :props="props">
                {{ props.row.entity_type }} ({{ props.row.entity_id || '-' }})
              </q-td>
              <q-td key="ip_address" :props="props" class="text-mono text-caption">
                {{ props.row.ip_address || '-' }}
              </q-td>
            </q-tr>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useQuasar, date } from 'quasar';
import { useAuditLogStore } from 'src/stores/auditLog';

const $q = useQuasar();
const auditStore = useAuditLogStore();

const filters = reactive({
  actorName: '',
  action: 'all',
  from: '',
  to: ''
});

const pagination = ref({
  page: 1,
  rowsPerPage: 20,
  rowsNumber: 0
});

const actionOptions = [
  { label: 'Tutte le azioni', value: 'all' },
  { label: 'Login', value: 'login' },
  { label: 'Voti (grade.*)', value: 'grade' },
  { label: 'Utenti (user.*)', value: 'user' },
  { label: 'Comunicazioni', value: 'communication' }
];

const columns = [
  { name: 'created_at', label: 'Data / Ora', field: 'created_at', align: 'left', sortable: true },
  { name: 'actor_name', label: 'Utente', field: 'actor_name', align: 'left', sortable: true },
  { name: 'actor_role', label: 'Ruolo', field: 'actor_role', align: 'center' },
  { name: 'action', label: 'Azione', field: 'action', align: 'center', sortable: true },
  { name: 'entity_type', label: 'Elemento', field: 'entity_type', align: 'left' },
  { name: 'ip_address', label: 'Indirizzo IP', field: 'ip_address', align: 'left' }
];

onMounted(() => {
  loadLogs();
});

const loadLogs = async () => {
  try {
    const res = await auditStore.fetchLogs(filters, pagination.value.page, pagination.value.rowsPerPage);
    pagination.value.rowsNumber = res.total;
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento degli audit log' });
  }
};

const applyFilters = () => {
  pagination.value.page = 1;
  loadLogs();
};

const onRequest = (props) => {
  pagination.value.page = props.pagination.page;
  pagination.value.rowsPerPage = props.pagination.rowsPerPage;
  loadLogs();
};

const exportCSV = async () => {
  try {
    await auditStore.exportLogs(filters);
    $q.notify({ type: 'positive', message: 'Export CSV scaricato' });
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante l\'exportation CSV' });
  }
};

const getRowClass = (action) => {
  if (!action) return '';
  const act = action.toLowerCase();
  if (act.includes('delete') || act.includes('fail')) return 'bg-red-50';
  if (act.includes('update') || act.includes('patch')) return 'bg-orange-50';
  return '';
};

const getActionColor = (action) => {
  if (!action) return 'grey-7';
  const act = action.toLowerCase();
  if (act.includes('delete') || act.includes('fail')) return 'negative';
  if (act.includes('update')) return 'warning';
  if (act.includes('login')) return 'info';
  return 'positive';
};

const formatDate = (d) => {
  if (!d) return '-';
  return date.formatDate(d, 'DD/MM/YYYY HH:mm:ss');
};
</script>
