<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-lg">
      <div class="col">
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">Gestione PCTO</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">Pianificazione e monitoraggio Percorsi per le Competenze Trasversali e l'Orientamento</p>
      </div>
      <div class="col-auto">
        <q-btn color="primary" icon="add" label="Nuovo Progetto" class="rounded-lg q-px-md shadow-sm" @click="openCreateDialog" />
        <q-btn outline color="primary" icon="business" label="Aziende" class="q-ml-sm rounded-lg" @click="showCompanies = true" />
      </div>
    </div>

    <div class="row q-col-gutter-lg q-mb-lg">
      <div class="col-12 col-md-3" v-for="stat in stats" :key="stat.label">
        <q-card class="rounded-xl shadow-soft border-slate-100 h-full overflow-hidden">
          <div class="absolute-top-right q-pa-md opacity-10">
            <q-icon :name="stat.icon" size="64px" />
          </div>
          <q-card-section>
            <div class="text-overline text-slate-400">{{ stat.label }}</div>
            <div class="text-h3 text-weight-bold text-slate-800 q-my-sm">{{ stat.value }}</div>
            <div class="row items-center">
              <q-icon :name="stat.trend > 0 ? 'trending_up' : 'trending_flat'" :color="stat.trend > 0 ? 'positive' : 'grey-5'" class="q-mr-xs" />
              <span :class="stat.trend > 0 ? 'text-positive' : 'text-grey-5'" class="text-caption font-medium">{{ stat.trendText }}</span>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey-7 bg-white"
        active-color="primary"
        indicator-color="primary"
        align="left"
        narrow-indicator
      >
        <q-tab name="active" label="Progetti Attivi" />
        <q-tab name="archived" label="Archivio" />
      </q-tabs>

      <q-separator />

      <q-table
        :rows="filteredProjects"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        class="bg-white"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template v-slot:body-cell-type="props">
          <q-td :props="props">
            <q-chip 
              :color="props.value === 'Interno' ? 'indigo-50' : 'blue-50'" 
              :text-color="props.value === 'Interno' ? 'indigo-700' : 'blue-700'"
              size="sm" 
              class="text-weight-medium rounded-md"
            >
              {{ props.value }}
            </q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="getStatusColor(props.row)" rounded class="q-px-sm q-py-xs">
              {{ getStatusLabel(props.row) }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-right">
            <q-btn flat round dense color="blue-600" icon="edit" @click="editProject(props.row)">
              <q-tooltip>Modifica</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="indigo-600" icon="group_add" @click="assignStudents(props.row)">
              <q-tooltip>Assegna Studenti</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="red-600" icon="delete" @click="confirmDelete(props.row)">
              <q-tooltip>Elimina</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Project Dialog -->
    <q-dialog v-model="dialog" persistent>
      <q-card style="min-width: 500px" class="rounded-xl shadow-2xl">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold">{{ isEdit ? 'Modifica Progetto' : 'Nuovo Progetto PCTO' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit="saveProject" class="q-gutter-md">
            <q-input v-model="form.title" label="Titolo Progetto" outlined dense :rules="[val => !!val || 'Campo richiesto']" />
            <q-input v-model="form.description" type="textarea" label="Descrizione" outlined dense />
            
            <div class="row q-col-gutter-md">
              <div class="col-6">
                <q-select v-model="form.type" :options="['Interno', 'Esterno']" label="Tipologia" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model.number="form.total_hours" type="number" label="Ore Totali" outlined dense />
              </div>
            </div>

            <div class="row q-col-gutter-md">
              <div class="col-6">
                <q-input v-model="form.start_date" type="date" label="Data Inizio" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model="form.end_date" type="date" label="Data Fine" outlined dense />
              </div>
            </div>

            <q-select 
              v-if="form.type === 'Esterno'"
              v-model="form.company_id" 
              :options="companies" 
              option-label="name"
              option-value="id"
              emit-value
              map-options
              label="Azienda Ospitante" 
              outlined 
              dense 
            />

            <q-input v-model="form.company_tutor_name" label="Tutor Aziendale" outlined dense />

            <div class="row justify-end q-mt-lg">
              <q-btn flat label="Annulla" color="grey-7" v-close-popup class="q-mr-sm" />
              <q-btn type="submit" color="primary" :label="isEdit ? 'Aggiorna' : 'Crea'" class="q-px-lg rounded-md" :loading="saving" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Companies List Dialog -->
    <q-dialog v-model="showCompanies" maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="bg-slate-50">
        <q-toolbar class="bg-white border-b-slate-100">
          <q-btn flat round dense icon="arrow_back" v-close-popup />
          <q-toolbar-title class="text-weight-bold">Gestione Aziende Convenzionate</q-toolbar-title>
          <q-btn color="primary" label="Aggiungi Azienda" @click="showAddCompany = true" />
        </q-toolbar>

        <q-card-section>
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-4" v-for="company in companies" :key="company.id">
              <q-card class="rounded-xl shadow-soft border-slate-100 h-full">
                <q-card-section>
                  <div class="row items-center no-wrap">
                    <q-avatar color="blue-50" text-color="blue-700" icon="business" />
                    <div class="q-ml-md">
                      <div class="text-subtitle1 text-weight-bold">{{ company.name }}</div>
                      <div class="text-caption text-grey-6">P.IVA: {{ company.vat_number }}</div>
                    </div>
                  </div>
                  <q-separator class="q-my-md opacity-50" />
                  <div class="q-gutter-xs">
                    <div class="row items-center">
                      <q-icon name="place" color="grey-5" size="16px" class="q-mr-sm" />
                      <span class="text-caption text-slate-700">{{ company.address }}</span>
                    </div>
                    <div class="row items-center">
                      <q-icon name="person" color="grey-5" size="16px" class="q-mr-sm" />
                      <span class="text-caption text-slate-700">{{ company.contact_person || 'N/D' }}</span>
                    </div>
                  </div>
                </q-card-section>
                <q-card-actions align="right">
                  <q-btn flat dense color="primary" icon="edit" label="Modifica" size="sm" />
                  <q-btn flat dense color="negative" icon="delete" label="Rimuovi" size="sm" />
                </q-card-actions>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Add Company Dialog -->
    <q-dialog v-model="showAddCompany">
      <q-card style="min-width: 400px" class="rounded-xl">
        <q-card-section class="text-h6 text-weight-bold">Aggiungi Nuova Azienda</q-card-section>
        <q-card-section>
          <q-form @submit="addCompany" class="q-gutter-md">
            <q-input v-model="companyForm.name" label="Ragione Sociale" outlined dense />
            <q-input v-model="companyForm.vat_number" label="P.IVA" outlined dense />
            <q-input v-model="companyForm.address" label="Indirizzo Sede" outlined dense />
            <q-input v-model="companyForm.contact_person" label="Persona di Riferimento" outlined dense />
            <q-input v-model="companyForm.email" label="Email Contatto" outlined dense />
            <div class="row justify-end q-mt-lg">
              <q-btn flat label="Annulla" v-close-popup />
              <q-btn type="submit" color="primary" label="Salva Azienda" class="q-ml-sm" :loading="savingCompany" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import api from 'src/services/api';

const $q = useQuasar();
const loading = ref(false);
const saving = ref(false);
const activeTab = ref('active');
const projects = ref([]);
const companies = ref([]);
const dialog = ref(false);
const isEdit = ref(false);
const showCompanies = ref(false);
const showAddCompany = ref(false);
const savingCompany = ref(false);

const pctoStats = ref({
  total_projects: 0,
  total_students: 0,
  total_hours: 0,
  active_companies: 0
});

const stats = computed(() => [
  { label: 'Progetti Totali', value: pctoStats.value.total_projects, icon: 'work', trend: 0, trendText: 'Dati in tempo reale' },
  { label: 'Studenti Coinvolti', value: pctoStats.value.total_students, icon: 'people', trend: 0, trendText: 'Iscritti ai percorsi' },
  { label: 'Ore Registrate', value: pctoStats.value.total_hours.toFixed(0), icon: 'timer', trend: 0, trendText: 'Ore totali validate' },
  { label: 'Aziende Partner', value: pctoStats.value.active_companies, icon: 'business', trend: 0, trendText: 'Convenzioni attive' },
]);

const form = ref({
  id: null,
  title: '',
  description: '',
  type: 'Esterno',
  start_date: '',
  end_date: '',
  total_hours: 40,
  company_id: null,
  company_tutor_name: ''
});

const companyForm = ref({
  name: '',
  vat_number: '',
  address: '',
  contact_person: '',
  email: ''
});

const columns = [
  { name: 'title', label: 'Progetto', field: 'title', align: 'left', sortable: true },
  { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
  { name: 'hours', label: 'Ore', field: 'total_hours', align: 'center' },
  { name: 'dates', label: 'Periodo', field: row => `${row.start_date.split('T')[0]} - ${row.end_date.split('T')[0]}`, align: 'left' },
  { name: 'status', label: 'Stato', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'right' }
];

const filteredProjects = computed(() => {
  const now = new Date();
  if (activeTab.value === 'active') {
    return projects.value.filter(p => new Date(p.end_date) >= now);
  }
  return projects.value.filter(p => new Date(p.end_date) < now);
});

onMounted(async () => {
  await fetchData();
});

const fetchData = async () => {
  loading.value = true;
  try {
    const [pRes, cRes, sRes] = await Promise.all([
      api.get('/pcto/projects'),
      api.get('/pcto/companies'),
      api.get('/pcto/stats')
    ]);
    projects.value = pRes.data || [];
    companies.value = cRes.data || [];
    pctoStats.value = sRes.data || pctoStats.value;
  } catch (err) {
    console.error(err);
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento dei dati' });
  } finally {
    loading.value = false;
  }
};

const openCreateDialog = () => {
  isEdit.value = false;
  form.value = {
    id: null, title: '', description: '', type: 'Esterno',
    start_date: '', end_date: '', total_hours: 40,
    company_id: null, company_tutor_name: ''
  };
  dialog.value = true;
};

const editProject = (row) => {
  isEdit.value = true;
  form.value = { 
    ...row,
    start_date: row.start_date.split('T')[0],
    end_date: row.end_date.split('T')[0]
  };
  dialog.value = true;
};

const saveProject = async () => {
  saving.value = true;
  try {
    if (isEdit.value) {
      await api.put(`/pcto/projects/${form.value.id}`, form.value);
      $q.notify({ type: 'positive', message: 'Progetto aggiornato' });
    } else {
      await api.post('/pcto/projects', form.value);
      $q.notify({ type: 'positive', message: 'Progetto creato' });
    }
    dialog.value = false;
    await fetchData();
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' });
  } finally {
    saving.value = false;
  }
};

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: `Sei sicuro di voler eliminare il progetto "${row.title}"? Questa azione è irreversibile.`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: 'Elimina' }
  }).onOk(async () => {
    try {
      await api.delete(`/pcto/projects/${row.id}`);
      $q.notify({ type: 'positive', message: 'Progetto eliminato' });
      await fetchData();
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione' });
    }
  });
};

const addCompany = async () => {
  savingCompany.value = true;
  try {
    await api.post('/pcto/companies', companyForm.value);
    $q.notify({ type: 'positive', message: 'Azienda aggiunta' });
    showAddCompany.value = false;
    await fetchData();
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' });
  } finally {
    savingCompany.value = false;
  }
};

const getStatusLabel = (row) => {
  const now = new Date();
  const start = new Date(row.start_date);
  const end = new Date(row.end_date);
  if (now < start) return 'In arrivo';
  if (now > end) return 'Completato';
  return 'In corso';
};

const getStatusColor = (row) => {
  const status = getStatusLabel(row);
  if (status === 'In corso') return 'positive';
  if (status === 'Completato') return 'grey-7';
  return 'orange-6';
};

const assignStudents = (row) => {
  // Navigation or specific dialog for student assignments
  $q.notify({ message: 'Funzionalità di assegnazione in fase di sviluppo' });
};
</script>

<style scoped>
.rounded-xl { border-radius: 1rem; }
.rounded-lg { border-radius: 0.75rem; }
.rounded-md { border-radius: 0.5rem; }
.shadow-soft { box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05); }
.border-slate-100 { border: 1px solid #f1f5f9; }
.bg-slate-50 { background-color: #f8fafc; }
</style>
