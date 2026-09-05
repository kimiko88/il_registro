<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-xl">
      <div>
        <h1 class="text-h4 text-weight-bold text-outfit q-my-none text-gradient-premium">
          Gestione PCTO
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-sm q-mb-none">Pianificazione e monitoraggio percorsi per l'orientamento</p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn color="primary" unelevated icon="add" label="Nuovo Progetto" class="rounded-lg q-px-md shadow-sm" no-caps @click="openCreateDialog" />
        <q-btn outline color="primary" icon="business" label="Aziende" class="rounded-lg q-px-md" no-caps @click="showCompanies = true" />
      </div>
    </div>

    <div class="row q-col-gutter-lg q-mb-xl">
      <div class="col-12 col-sm-6 col-md-3" v-for="stat in stats" :key="stat.label">
        <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden h-full bg-white">
          <q-card-section class="q-pa-lg">
            <div class="row items-center justify-between q-mb-md">
              <div class="text-overline text-slate-400 letter-spacing-1">{{ stat.label }}</div>
              <q-avatar :color="stat.color + '-50'" :text-color="stat.color + '-700'" :icon="stat.icon" size="40px" />
            </div>
            <div class="text-h3 text-weight-bold text-slate-800 q-my-none">{{ stat.value }}</div>
            <div class="text-caption text-slate-500 q-mt-sm">{{ stat.trendText }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden bg-white">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-slate-500 border-b border-slate-100"
        active-color="primary"
        indicator-color="primary"
        align="left"
        narrow-indicator
        no-caps
      >
        <q-tab name="active" label="Progetti Attivi" class="q-px-lg py-4" />
        <q-tab name="archived" label="Archivio Storico" class="q-px-lg py-4" />
      </q-tabs>

      <q-tab-panels v-model="activeTab" animated class="bg-transparent">
        <q-tab-panel name="active" class="q-pa-none">
          <q-table
            :rows="filteredProjects"
            :columns="columns"
            row-key="id"
            flat
            :loading="loading"
            class="bg-transparent"
            :pagination="{ rowsPerPage: 10 }"
          >
            <template v-slot:body-cell-type="props">
              <q-td :props="props">
                <q-chip 
                  :color="props.value === 'Interno' ? 'indigo-50' : 'blue-50'" 
                  :text-color="props.value === 'Interno' ? 'indigo-700' : 'blue-700'"
                  size="sm" 
                  class="text-weight-bold rounded-md"
                >
                  {{ props.value }}
                </q-chip>
              </q-td>
            </template>

            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge :color="getStatusColor(props.row)" rounded class="q-px-sm q-py-xs shadow-xs">
                  {{ getStatusLabel(props.row) }}
                </q-badge>
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props" class="text-right">
                <div class="row q-gutter-xs justify-end">
                  <q-btn flat round dense color="primary" icon="edit" @click="editProject(props.row)">
                    <q-tooltip>Modifica Progetto</q-tooltip>
                  </q-btn>
                  <q-btn flat round dense color="indigo" icon="group_add" @click="assignStudents(props.row)">
                    <q-tooltip>Assegna Studenti</q-tooltip>
                  </q-btn>
                  <q-btn flat round dense color="negative" icon="delete" @click="confirmDelete(props.row)">
                    <q-tooltip>Elimina</q-tooltip>
                  </q-btn>
                </div>
              </q-td>
            </template>
            
            <template v-slot:no-data>
              <div class="full-width q-pa-xl text-center text-slate-400">
                <q-icon name="work_outline" size="64px" class="opacity-20 q-mb-md" />
                <div class="text-h6">Nessun progetto trovato</div>
              </div>
            </template>
          </q-table>
        </q-tab-panel>
        
        <q-tab-panel name="archived" class="q-pa-none">
          <div class="q-pa-xl text-center text-slate-400">
            <q-icon name="history" size="64px" class="opacity-20 q-mb-md" />
            <div class="text-h6">L'archivio storico verrà caricato a breve</div>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Project Dialog -->
    <q-dialog v-model="dialog" persistent class="premium-dialog">
      <q-card style="width: min(600px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden shadow-24">
        <q-card-section class="bg-gradient-primary text-white q-pa-lg row items-center justify-between">
          <div class="text-h5 text-weight-bold">{{ isEdit ? 'Modifica Progetto' : 'Nuovo Progetto PCTO' }}</div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-xl">
          <q-form @submit="saveProject" class="q-gutter-y-lg">
            <q-input v-model="form.title" label="Titolo Progetto" outlined :rules="[val => !!val || 'Campo richiesto']" />
            <q-input v-model="form.description" type="textarea" label="Descrizione" outlined rows="3" />
            
            <div class="row q-col-gutter-lg">
              <div class="col-6">
                <q-select v-model="form.type" :options="['Interno', 'Esterno']" label="Tipologia" outlined />
              </div>
              <div class="col-6">
                <q-input v-model.number="form.total_hours" type="number" label="Ore Totali" outlined />
              </div>
            </div>

            <div class="row q-col-gutter-lg">
              <div class="col-6">
                <q-input v-model="form.start_date" type="date" label="Data Inizio" outlined stack-label />
              </div>
              <div class="col-6">
                <q-input v-model="form.end_date" type="date" label="Data Fine" outlined stack-label />
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
            />

            <q-input v-model="form.company_tutor_name" label="Tutor Aziendale" outlined />

            <div class="row justify-end q-mt-xl q-gutter-sm">
              <q-btn flat label="Annulla" color="slate-400" v-close-popup no-caps />
              <q-btn type="submit" color="primary" :label="isEdit ? 'Aggiorna Progetto' : 'Crea Progetto'" class="q-px-xl rounded-lg shadow-sm" no-caps :loading="saving" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Companies List Dialog -->
    <q-dialog v-model="showCompanies" maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="bg-slate-50">
        <q-toolbar class="bg-white border-b border-slate-100 q-py-md">
          <q-btn flat round dense icon="arrow_back" v-close-popup class="text-slate-400" />
          <q-toolbar-title class="text-weight-bold text-slate-800">Aziende Convenzionate</q-toolbar-title>
          <q-btn unelevated color="primary" label="Aggiungi Azienda" class="rounded-lg q-px-md" no-caps @click="showAddCompany = true" />
        </q-toolbar>

        <q-card-section class="q-pa-xl">
          <div class="row q-col-gutter-xl">
            <div class="col-12 col-sm-6 col-md-4" v-for="company in companies" :key="company.id">
              <q-card flat class="rounded-xl border-slate-100 bg-white shadow-soft h-full overflow-hidden">
                <q-card-section class="q-pa-lg">
                  <div class="row items-center no-wrap q-mb-lg">
                    <q-avatar color="blue-50" text-color="blue-700" icon="business" size="48px" />
                    <div class="q-ml-md">
                      <div class="text-subtitle1 text-weight-bold text-slate-800">{{ company.name }}</div>
                      <div class="text-caption text-slate-400">P.IVA: {{ company.vat_number }}</div>
                    </div>
                  </div>
                  <q-separator class="q-my-lg opacity-50" />
                  <div class="q-gutter-y-sm">
                    <div class="row items-center">
                      <q-icon name="place" color="slate-300" size="18px" class="q-mr-sm" />
                      <span class="text-caption text-slate-600">{{ company.address }}</span>
                    </div>
                    <div class="row items-center">
                      <q-icon name="person" color="slate-300" size="18px" class="q-mr-sm" />
                      <span class="text-caption text-slate-600">
                        {{ company.contact_person_first_name ? `${company.contact_person_first_name} ${company.contact_person_last_name}` : (company.contact_person || 'Contatto non definito') }}
                      </span>
                    </div>
                    <div v-if="company.contact_person_phone" class="row items-center">
                      <q-icon name="phone" color="slate-300" size="18px" class="q-mr-sm" />
                      <span class="text-caption text-slate-600">{{ company.contact_person_phone }}</span>
                    </div>
                    <div v-if="company.email" class="row items-center">
                      <q-icon name="email" color="slate-300" size="18px" class="q-mr-sm" />
                      <span class="text-caption text-slate-600">{{ company.email }}</span>
                    </div>
                  </div>
                </q-card-section>
                <q-card-actions align="right" class="bg-slate-50 q-pa-md">
                  <q-btn flat dense color="primary" icon="edit" label="Modifica" size="sm" no-caps />
                  <q-btn flat dense color="negative" icon="delete" label="Rimuovi" size="sm" no-caps />
                </q-card-actions>
              </q-card>
            </div>
          </div>
          
          <div v-if="companies.length === 0" class="full-width q-pa-xl text-center text-slate-400">
            <q-icon name="business_center" size="64px" class="opacity-20 q-mb-md" />
            <div class="text-h6">Nessuna azienda convenzionata</div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Add Company Dialog -->
    <q-dialog v-model="showAddCompany" class="premium-dialog">
      <q-card style="width: min(450px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden shadow-24 bg-white">
        <q-card-section class="bg-gradient-primary text-white q-pa-lg row items-center justify-between">
          <div class="text-h5 text-weight-bold">Nuova Azienda</div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>
        <q-card-section class="q-pa-xl">
          <q-form @submit="addCompany" class="q-gutter-y-lg">
            <q-input v-model="companyForm.name" label="Ragione Sociale" outlined />
            <q-input v-model="companyForm.vat_number" label="Partita IVA" outlined />
            <q-input v-model="companyForm.address" label="Indirizzo Sede" outlined />
            <div class="row q-col-gutter-md">
              <div class="col-6">
                <q-input v-model="companyForm.contact_person_first_name" label="Nome Referente" outlined />
              </div>
              <div class="col-6">
                <q-input v-model="companyForm.contact_person_last_name" label="Cognome Referente" outlined />
              </div>
            </div>
            <q-input v-model="companyForm.contact_person_phone" label="Telefono Referente" outlined />
            <q-input v-model="companyForm.email" label="Email Referente" outlined type="email" />
            <div class="row justify-end q-mt-xl q-gutter-sm">
              <q-btn flat label="Annulla" color="slate-400" v-close-popup no-caps />
              <q-btn type="submit" color="primary" label="Salva Azienda" class="q-px-xl rounded-lg shadow-sm" no-caps :loading="savingCompany" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import api from '@/services/api';

const $q = useQuasar();
const { t } = useI18n();
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
  { label: 'Progetti Totali', value: pctoStats.value.total_projects, icon: 'work', color: 'indigo', trendText: 'Dati in tempo reale' },
  { label: 'Studenti Coinvolti', value: pctoStats.value.total_students, icon: 'people', color: 'blue', trendText: 'Iscritti ai percorsi' },
  { label: 'Ore Registrate', value: pctoStats.value.total_hours.toFixed(0), icon: 'timer', color: 'emerald', trendText: 'Ore totali validate' },
  { label: 'Aziende Partner', value: pctoStats.value.active_companies, icon: 'business', color: 'orange', trendText: 'Convenzioni attive' },
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
  contact_person_first_name: '',
  contact_person_last_name: '',
  contact_person_phone: '',
  email: ''
});

const columns = [
  { name: 'title', label: 'Progetto', field: 'title', align: 'left', sortable: true, classes: 'text-weight-bold text-slate-800' },
  { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
  { name: 'hours', label: 'Ore', field: 'total_hours', align: 'center' },
  { name: 'dates', label: 'Periodo', field: row => `${row.start_date.split('T')[0]} - ${row.end_date.split('T')[0]}`, align: 'left' },
  { name: 'status', label: 'Stato', align: 'center' },
  { name: 'actions', label: '', align: 'right' }
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
    ok: { color: 'negative', label: 'Elimina', flat: false }
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
    companyForm.value = {
      name: '',
      vat_number: '',
      address: '',
      contact_person: '',
      contact_person_first_name: '',
      contact_person_last_name: '',
      contact_person_phone: '',
      email: ''
    };
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
  $q.notify({ message: 'Funzionalità di assegnazione in fase di sviluppo' });
};
</script>

<style scoped>
.letter-spacing-1 { letter-spacing: 1px; }
.opacity-20 { opacity: 0.2; }
.py-4 { padding-top: 1rem; padding-bottom: 1rem; }
.rounded-xl { border-radius: 1rem; }
.rounded-lg { border-radius: 0.75rem; }
.rounded-md { border-radius: 0.5rem; }
.shadow-soft { box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05); }
.border-slate-100 { border: 1px solid #f1f5f9; }
.bg-slate-50 { background-color: #f8fafc; }
</style>
