<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-lg">
      <div class="col">
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">Archivio Documenti</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">Gestione PDP, PFI, Certificati e Documentazione PCTO</p>
      </div>
      <div class="col-auto">
        <q-btn color="primary" icon="add" label="Nuovo Documento" class="rounded-lg q-px-md shadow-sm" @click="openCreateDialog" />
        <q-btn flat color="primary" icon="settings" class="q-ml-sm rounded-lg" @click="showTemplates = true">
          <q-tooltip>Gestione Template</q-tooltip>
        </q-btn>
      </div>
    </div>

    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-3" v-for="cat in categories" :key="cat.type">
        <q-card 
          clickable 
          v-ripple 
          class="rounded-xl shadow-soft border-slate-100 cursor-pointer transition-all hover:translate-y-[-2px]"
          :class="selectedType === cat.type ? 'border-primary border-2' : ''"
          @click="selectedType = cat.type"
        >
          <q-card-section class="row items-center no-wrap">
            <q-avatar :color="cat.color + '-50'" :text-color="cat.color + '-700'" :icon="cat.icon" size="48px" />
            <div class="q-ml-md overflow-hidden">
              <div class="text-weight-bold text-slate-800">{{ cat.label }}</div>
              <div class="text-caption text-slate-500">{{ cat.count }} documenti</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden">
      <div class="q-pa-md row items-center border-b-slate-100">
        <q-input 
          v-model="search" 
          placeholder="Cerca per titolo, studente o classe..." 
          outlined 
          dense 
          class="col-12 col-md-4"
        >
          <template v-slot:prepend>
            <q-icon name="search" color="grey-5" />
          </template>
        </q-input>
        <q-space />
        <q-btn-toggle
          v-model="filterStatus"
          flat
          toggle-color="primary"
          color="grey-7"
          :options="[
            {label: 'Tutti', value: 'all'},
            {label: 'In Bozza', value: 'draft'},
            {label: 'Da Firmare', value: 'submitted'},
            {label: 'Firmati', value: 'signed'}
          ]"
        />
      </div>

      <q-table
        :rows="filteredDocs"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        class="bg-white"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template v-slot:body-cell-type="props">
          <q-td :props="props">
            <q-chip size="sm" class="text-weight-bold" outline :color="getTypeColor(props.value)">
              {{ props.value.toUpperCase() }}
            </q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="getStatusColor(props.value)" rounded class="q-px-sm q-py-xs">
              {{ getStatusLabel(props.value) }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-right">
            <q-btn flat round dense color="blue-600" icon="visibility" @click="previewDoc(props.row)">
              <q-tooltip>Visualizza</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="indigo-600" icon="edit" @click="editDoc(props.row)" v-if="props.row.status === 'draft'">
              <q-tooltip>Modifica</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="red-600" icon="delete" @click="confirmDelete(props.row)">
              <q-tooltip>Elimina</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Create Document Dialog -->
    <q-dialog v-model="createDialog" persistent maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="bg-slate-50">
        <q-toolbar class="bg-white border-b-slate-100 q-px-lg">
          <q-btn flat round dense icon="close" v-close-popup />
          <q-toolbar-title class="text-weight-bold">
            {{ isEdit ? 'Modifica Documento' : 'Nuovo Documento' }}
          </q-toolbar-title>
          <q-btn color="primary" label="Salva e Chiudi" @click="saveDocument" :loading="saving" />
        </q-toolbar>

        <q-card-section class="q-pa-lg">
          <div class="row q-col-gutter-lg justify-center">
            <div class="col-12 col-md-4">
              <q-card class="rounded-xl shadow-soft border-slate-100">
                <q-card-section>
                  <div class="text-subtitle1 text-weight-bold q-mb-md">Informazioni Generali</div>
                  <q-form class="q-gutter-md">
                    <q-input v-model="form.title" label="Titolo Documento" outlined dense />
                    <q-select 
                      v-model="form.type" 
                      :options="['pdp', 'pfi', 'certificate', 'pcto', 'generic']" 
                      label="Tipologia" 
                      outlined 
                      dense 
                      emit-value
                    />
                    <q-select 
                      v-model="form.student_id" 
                      label="Studente (opzionale)" 
                      outlined 
                      dense 
                      use-input
                      @filter="filterStudents"
                      :options="studentOptions"
                      option-label="name"
                      option-value="id"
                      emit-value
                      map-options
                    />
                  </q-form>
                </q-card-section>
              </q-card>
            </div>
            <div class="col-12 col-md-8">
              <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden">
                <q-card-section class="bg-slate-100 border-b-slate-200 row items-center">
                  <div class="text-subtitle2 text-weight-bold text-slate-700">Contenuto Documento</div>
                  <q-space />
                  <q-btn-dropdown flat label="Usa Template" color="primary" size="sm">
                    <q-list>
                      <q-item v-for="t in templates" :key="t.id" clickable v-close-popup @click="applyTemplate(t)">
                        <q-item-section>{{ t.name }}</q-item-section>
                      </q-item>
                      <q-item v-if="templates.length === 0">
                        <q-item-section class="text-grey italic">Nessun template disponibile</q-item-section>
                      </q-item>
                    </q-list>
                  </q-btn-dropdown>
                </q-card-section>
                <q-editor 
                  v-model="form.content" 
                  min-height="400px" 
                  flat 
                  :toolbar="[
                    ['bold', 'italic', 'strike', 'underline'],
                    ['quote', 'unordered', 'ordered'],
                    ['undo', 'redo'],
                    ['viewsource']
                  ]"
                />
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Preview Dialog -->
    <q-dialog v-model="showPreview">
      <q-card style="width: 700px; max-width: 90vw;" class="rounded-xl">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold">{{ selectedDoc?.title }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg scroll" style="max-height: 70vh">
          <div v-html="selectedDoc?.content || 'Caricamento...'" class="document-content"></div>
        </q-card-section>

        <q-card-actions align="right" class="bg-slate-50">
          <q-btn flat label="Scarica PDF" icon="picture_as_pdf" color="primary" />
          <q-btn unelevated label="Chiudi" color="grey-7" v-close-popup />
        </q-card-actions>
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
const search = ref('');
const selectedType = ref('all');
const filterStatus = ref('all');
const docs = ref([]);
const templates = ref([]);
const createDialog = ref(false);
const isEdit = ref(false);
const showPreview = ref(false);
const selectedDoc = ref(null);
const studentOptions = ref([]);

const categories = [
  { type: 'all', label: 'Tutti i Documenti', icon: 'folder', color: 'slate', count: 0 },
  { type: 'pdp', label: 'PDP (BES/DSA)', icon: 'auto_awesome', color: 'indigo', count: 0 },
  { type: 'pfi', label: 'PFI (Istruzione)', icon: 'assignment_ind', color: 'blue', count: 0 },
  { type: 'certificate', label: 'Certificati', icon: 'verified', color: 'emerald', count: 0 },
];

const form = ref({
  id: null,
  title: '',
  type: 'generic',
  student_id: null,
  content: '',
  change_log: 'Aggiornamento manuale'
});

const columns = [
  { name: 'title', label: 'Titolo', field: 'title', align: 'left', sortable: true },
  { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
  { name: 'student', label: 'Studente', field: 'student_id', align: 'left' },
  { name: 'updated', label: 'Ultima Modifica', field: row => new Date(row.updated_at).toLocaleDateString(), align: 'left' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'right' }
];

const filteredDocs = computed(() => {
  let filtered = docs.value;
  if (selectedType.value !== 'all') {
    filtered = filtered.filter(d => d.type === selectedType.value);
  }
  if (filterStatus.value !== 'all') {
    filtered = filtered.filter(d => d.status === filterStatus.value);
  }
  if (search.value) {
    const s = search.value.toLowerCase();
    filtered = filtered.filter(d => 
      d.title.toLowerCase().includes(s) || 
      (d.student_id && d.student_id.toLowerCase().includes(s))
    );
  }
  return filtered;
});

onMounted(async () => {
  await fetchData();
  updateCategoryCounts();
});

const fetchData = async () => {
  loading.value = true;
  try {
    const [dRes, tRes] = await Promise.all([
      api.get('/documents'),
      api.get('/documents/template')
    ]);
    docs.value = dRes.data || [];
    templates.value = tRes.data || [];
    updateCategoryCounts();
  } catch (err) {
    console.error(err);
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento dei documenti' });
  } finally {
    loading.value = false;
  }
};

const updateCategoryCounts = () => {
  categories[0].count = docs.value.length;
  categories[1].count = docs.value.filter(d => d.type === 'pdp').length;
  categories[2].count = docs.value.filter(d => d.type === 'pfi').length;
  categories[3].count = docs.value.filter(d => d.type === 'certificate').length;
};

const openCreateDialog = () => {
  isEdit.value = false;
  form.value = {
    id: null, title: '', type: 'generic', student_id: null, content: '', change_log: 'Creazione'
  };
  createDialog.value = true;
};

const editDoc = (row) => {
  isEdit.value = true;
  form.value = { ...row };
  createDialog.value = true;
};

const previewDoc = async (row) => {
  selectedDoc.value = row;
  showPreview.value = true;
  try {
    const res = await api.get(`/documents/${row.id}`);
    selectedDoc.value = res.data;
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Impossibile caricare il contenuto' });
  }
};

const saveDocument = async () => {
  saving.value = true;
  try {
    if (isEdit.value) {
      await api.patch(`/documents/${form.value.id}`, form.value);
      $q.notify({ type: 'positive', message: 'Documento aggiornato' });
    } else {
      await api.post('/documents', form.value);
      $q.notify({ type: 'positive', message: 'Documento creato' });
    }
    createDialog.value = false;
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
    message: `Sei sicuro di voler eliminare il documento "${row.title}"?`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: 'Elimina' }
  }).onOk(async () => {
    try {
      await api.delete(`/documents/${row.id}`);
      $q.notify({ type: 'positive', message: 'Documento eliminato' });
      await fetchData();
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione' });
    }
  });
};

const filterStudents = (val, update) => {
  if (val === '') {
    update(() => {
      studentOptions.value = [];
    });
    return;
  }
  // Implement actual student search here
  update(() => {
    studentOptions.value = [];
  });
};

const applyTemplate = (t) => {
  form.value.content = t.content;
  form.value.type = t.type;
  if (!form.value.title) form.value.title = t.name;
};

const getTypeColor = (type) => {
  switch (type) {
    case 'pdp': return 'indigo';
    case 'pfi': return 'blue';
    case 'certificate': return 'emerald';
    case 'pcto': return 'orange';
    default: return 'slate';
  }
};

const getStatusLabel = (status) => {
  switch (status) {
    case 'draft': return 'In Bozza';
    case 'submitted': return 'Da Firmare';
    case 'review': return 'In Revisione';
    case 'signed': return 'Firmato';
    case 'rejected': return 'Rifiutato';
    default: return status;
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case 'draft': return 'grey-7';
    case 'submitted': return 'orange-6';
    case 'review': return 'blue-6';
    case 'signed': return 'positive';
    case 'rejected': return 'negative';
    default: return 'grey';
  }
};
</script>

<style scoped>
.rounded-xl { border-radius: 1rem; }
.rounded-lg { border-radius: 0.75rem; }
.shadow-soft { box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05); }
.border-slate-100 { border: 1px solid #f1f5f9; }
.border-b-slate-100 { border-bottom: 1px solid #f1f5f9; }
.border-b-slate-200 { border-bottom: 1px solid #e2e8f0; }
.bg-slate-50 { background-color: #f8fafc; }
.document-content {
  line-height: 1.6;
  font-family: serif;
  background: white;
  padding: 2rem;
  border: 1px solid #eee;
}
</style>
