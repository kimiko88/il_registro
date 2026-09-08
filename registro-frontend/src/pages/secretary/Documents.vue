<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h4 text-weight-bold text-outfit q-my-none text-gradient-premium">Archivio Documenti</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">Gestione PDP, PFI, Certificati e Documentazione PCTO</p>
      </div>
      <div class="col-auto row q-gutter-sm">
        <q-btn unelevated color="primary" icon="add" label="Nuovo Documento" class="rounded-lg q-px-lg shadow-sm" no-caps @click="openCreateDialog" />
        <q-btn outline color="primary" icon="settings" label="Modelli" class="rounded-lg q-px-md" no-caps @click="showTemplates = true">
          <q-tooltip>Gestione Modelli Documenti</q-tooltip>
        </q-btn>
      </div>
    </div>

    <div class="row q-col-gutter-lg q-mb-xl">
      <div class="col-12 col-md-3" v-for="cat in categories" :key="cat.type">
        <q-card 
          clickable 
          v-ripple 
          class="rounded-xl shadow-soft border-slate-100 cursor-pointer transition-all hover:translate-y-[-4px]"
          :class="selectedType === cat.type ? 'border-primary ring-2 ring-primary ring-opacity-10' : 'bg-white'"
          @click="selectedType = cat.type"
        >
          <q-card-section class="row items-center no-wrap q-pa-lg">
            <q-avatar :color="cat.color + '-50'" :text-color="cat.color + '-700'" :icon="cat.icon" size="56px" class="rounded-lg" />
            <div class="q-ml-lg overflow-hidden">
              <div class="text-h6 text-weight-bold text-slate-800 line-height-tight">{{ cat.label }}</div>
              <div class="text-subtitle2 text-slate-400">{{ cat.count }} file archiviati</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden bg-white">
      <div class="q-pa-lg row items-center border-b border-slate-100 bg-slate-50/50">
        <q-input 
          v-model="search" 
          placeholder="Cerca per titolo, studente o classe..." 
          outlined 
          dense 
          class="col-12 col-md-4 bg-white"
        >
          <template v-slot:prepend>
            <q-icon name="search" color="slate-400" />
          </template>
        </q-input>
        <q-space />
        <q-btn-toggle
          v-model="filterStatus"
          unelevated
          toggle-color="indigo-50"
          toggle-text-color="indigo-700"
          color="white"
          text-color="slate-400"
          class="rounded-lg border border-slate-100"
          no-caps
          :options="[
            {label: 'Tutti', value: 'all'},
            {label: 'In Bozza', value: 'draft'},
            {label: 'Inviati', value: 'submitted'},
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
        class="bg-transparent"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template v-slot:body-cell-type="props">
          <q-td :props="props">
            <q-chip size="sm" class="text-weight-bold rounded-md" :color="getTypeColor(props.value) + '-50'" :text-color="getTypeColor(props.value) + '-700'">
              {{ props.value.toUpperCase() }}
            </q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="getStatusColor(props.value)" rounded class="q-px-md q-py-xs text-weight-bold">
              {{ getStatusLabel(props.value) }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-right">
            <q-btn flat round dense color="slate-400" icon="visibility" @click="previewDoc(props.row)">
              <q-tooltip>Visualizza</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="primary" icon="edit" @click="editDoc(props.row)" v-if="props.row.status === 'draft'">
              <q-tooltip>Modifica</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="negative" icon="delete_outline" @click="confirmDelete(props.row)">
              <q-tooltip>Elimina</q-tooltip>
            </q-btn>
          </q-td>
        </template>
        
        <template v-slot:no-data>
          <div class="full-width q-pa-xl text-center text-slate-400">
            <q-icon name="description" size="64px" class="opacity-10 q-mb-md" />
            <div class="text-h6">Nessun documento trovato</div>
          </div>
        </template>
      </q-table>
    </q-card>

    <!-- Create Document Dialog -->
    <q-dialog v-model="createDialog" persistent maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="bg-slate-50 column no-wrap">
        <q-toolbar class="bg-white border-b border-slate-100 q-px-xl q-py-md">
          <q-btn flat round dense icon="close" v-close-popup color="slate-400" :aria-label="$t('common.close') || 'Chiudi'" />
          <q-toolbar-title class="text-weight-bold text-slate-800 text-outfit">
            {{ isEdit ? 'Modifica Documento' : 'Redazione Nuovo Documento' }}
          </q-toolbar-title>
          <q-btn unelevated color="primary" label="Salva Documento" class="rounded-lg q-px-lg shadow-sm" no-caps @click="saveDocument" :loading="saving" />
        </q-toolbar>

        <q-card-section class="col q-pa-xl scroll">
          <div class="row q-col-gutter-xl justify-center max-w-7xl mx-auto">
            <div class="col-12 col-md-4">
              <q-card flat class="rounded-xl border border-slate-100 bg-white shadow-soft">
                <q-card-section class="q-pa-xl">
                  <div class="text-h6 text-weight-bold text-slate-800 q-mb-xl">Informazioni Documento</div>
                  <q-form class="q-gutter-y-lg">
                    <q-input v-model="form.title" label="Titolo Documento" outlined class="rounded-lg" />
                    <q-select 
                      v-model="form.type" 
                      :options="typeOptions" 
                      label="Tipologia" 
                      outlined 
                      emit-value
                      map-options
                      class="rounded-lg"
                    />
                    <q-select 
                      v-model="form.student_id" 
                      label="Studente di riferimento" 
                      outlined 
                      use-input
                      @filter="filterStudents"
                      :options="studentOptions"
                      option-label="name"
                      option-value="id"
                      emit-value
                      map-options
                      class="rounded-lg"
                    />
                  </q-form>
                </q-card-section>
              </q-card>
              
              <q-card flat class="rounded-xl border border-slate-100 bg-indigo-50 q-mt-xl shadow-soft">
                <q-card-section class="q-pa-lg">
                  <div class="row items-center q-gutter-sm q-mb-md">
                    <q-icon name="info" color="indigo" />
                    <div class="text-subtitle2 text-indigo-700 font-medium">Suggerimento</div>
                  </div>
                  <div class="text-body2 text-indigo-600">
                    Utilizza i modelli predefiniti per velocizzare la compilazione. Puoi personalizzare i modelli cliccando sul tasto "Modelli" in alto.
                  </div>
                </q-card-section>
              </q-card>
            </div>
            
            <div class="col-12 col-md-8">
              <q-card flat class="rounded-xl border border-slate-100 bg-white shadow-soft overflow-hidden">
                <q-card-section class="bg-slate-50 border-b border-slate-100 q-pa-lg row items-center justify-between">
                  <div class="text-subtitle1 text-weight-bold text-slate-700">Editor Contenuto</div>
                  <q-btn-dropdown unelevated label="Applica Modello" color="indigo-50" text-color="indigo-700" icon="auto_awesome" class="rounded-lg no-caps">
                    <q-list padding class="rounded-lg">
                      <q-item v-for="t in templates" :key="t.id" clickable v-close-popup class="q-mx-sm rounded-md" @click="applyTemplate(t)">
                        <q-item-section>
                          <q-item-label class="text-weight-medium">{{ t.name }}</q-item-label>
                          <q-item-label caption>{{ getTypeLabel(t.type) }}</q-item-label>
                        </q-item-section>
                      </q-item>
                      <q-item v-if="templates.length === 0" class="q-pa-md text-center text-slate-400">
                        <q-item-section>Nessun modello disponibile</q-item-section>
                      </q-item>
                    </q-list>
                  </q-btn-dropdown>
                </q-card-section>
                <q-editor 
                  v-model="form.content" 
                  min-height="600px" 
                  flat 
                  class="q-pa-lg"
                  :toolbar="[
                    ['bold', 'italic', 'strike', 'underline'],
                    ['quote', 'unordered', 'ordered'],
                    ['undo', 'redo'],
                    ['viewsource', 'fullscreen']
                  ]"
                />
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Preview Dialog -->
    <q-dialog v-model="showPreview" class="premium-dialog">
      <q-card style="width: 850px; max-width: 95vw;" class="rounded-2xl overflow-hidden shadow-24 bg-white">
        <q-card-section class="row items-center q-pa-xl border-b border-slate-100">
          <div class="text-h5 text-weight-bold text-slate-800 text-outfit">{{ selectedDoc?.title }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup color="slate-400" :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-xl scroll bg-slate-50" style="max-height: 75vh">
          <div class="document-paper shadow-lg rounded-sm q-pa-xl bg-white mx-auto" style="max-width: 800px">
            <div v-html="sanitizedPreviewContent" class="document-content-html"></div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-lg bg-white border-t border-slate-100">
          <q-btn flat label="Chiudi" color="slate-400" v-close-popup no-caps />
          <q-btn unelevated label="Esporta PDF" icon="picture_as_pdf" color="primary" class="rounded-lg q-px-lg shadow-sm" no-caps />
        </q-card-actions>
      </q-card>
    </q-dialog>
    <TemplateManager v-model="showTemplates" @templates-updated="fetchData" />
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import api from '@/services/api';
import TemplateManager from '@/components/Secretary/TemplateManager.vue';
import { sanitizeHTMLContent } from '@/utils/sanitize';

const $q = useQuasar();
const { t } = useI18n();
const loading = ref(false);
const saving = ref(false);
const search = ref('');
const selectedType = ref('all');
const filterStatus = ref('all');
const docs = ref([]);
const templates = ref([]);
const createDialog = ref(false);
const showTemplates = ref(false);
const isEdit = ref(false);
const showPreview = ref(false);
const selectedDoc = ref(null);
const studentOptions = ref([]);

const typeOptions = [
  { label: 'PDP (BES/DSA)', value: 'pdp' },
  { label: 'PFI (Istruzione)', value: 'pfi' },
  { label: 'Certificato', value: 'certificate' },
  { label: 'Documentazione PCTO', value: 'pcto' },
  { label: 'Altro / Generico', value: 'generic' }
];

const categories = [
  { type: 'all', label: 'Archivio Totale', icon: 'folder_copy', color: 'slate', count: 0 },
  { type: 'pdp', label: 'PDP (BES/DSA)', icon: 'description', color: 'indigo', count: 0 },
  { type: 'pfi', label: 'Progetti PFI', icon: 'assignment', color: 'blue', count: 0 },
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
  { name: 'title', label: 'Titolo Documento', field: 'title', align: 'left', sortable: true, classes: 'text-weight-bold text-slate-800' },
  { name: 'type', label: 'Tipologia', field: 'type', align: 'left' },
  { name: 'student', label: 'Studente', field: 'student_id', align: 'left', classes: 'text-slate-500' },
  { name: 'updated', label: 'Ultima Modifica', field: row => new Date(row.updated_at).toLocaleDateString('it-IT'), align: 'left', classes: 'text-slate-500' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center' },
  { name: 'actions', label: '', align: 'right' }
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

const sanitizedPreviewContent = computed(() => {
  return sanitizeHTMLContent(selectedDoc.value?.content || 'Caricamento...');
});

onMounted(async () => {
  await fetchData();
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
    $q.notify({ type: 'negative', message: 'Errore nel recupero della documentazione' });
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
    id: null, title: '', type: 'generic', student_id: null, content: '', change_log: 'Creazione documento'
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
    $q.notify({ type: 'negative', message: 'Errore nel caricamento del contenuto' });
  }
};

const saveDocument = async () => {
  saving.value = true;
  try {
    if (isEdit.value) {
      await api.patch(`/documents/${form.value.id}`, form.value);
      $q.notify({ type: 'positive', message: 'Documento aggiornato correttamente' });
    } else {
      await api.post('/documents', form.value);
      $q.notify({ type: 'positive', message: 'Nuovo documento archiviato' });
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
    title: 'Eliminazione Documento',
    message: `Sei sicuro di voler eliminare definitivamente il documento "${row.title}"?`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: 'Elimina Ora', flat: false }
  }).onOk(async () => {
    try {
      await api.delete(`/documents/${row.id}`);
      $q.notify({ type: 'positive', message: 'Documento rimosso' });
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

const getTypeLabel = (type) => {
  switch (type) {
    case 'circular': return 'Circolare';
    case 'certificate': return 'Certificato';
    case 'report': return 'Verbale/Pagella';
    case 'other': return 'Altro';
    default: return type;
  }
};

const getStatusLabel = (status) => {
  switch (status) {
    case 'draft': return 'Bozza';
    case 'submitted': return 'Inviato';
    case 'signed': return 'Firmato';
    case 'rejected': return 'Rifiutato';
    default: return status;
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case 'draft': return 'slate-400';
    case 'submitted': return 'orange-500';
    case 'signed': return 'emerald-500';
    case 'rejected': return 'rose-500';
    default: return 'slate-300';
  }
};
</script>

<style scoped>
.line-height-tight { line-height: 1.25; }
.max-w-7xl { max-width: 80rem; }
.mx-auto { margin-left: auto; margin-right: auto; }
.opacity-10 { opacity: 0.1; }
.document-paper {
  aspect-ratio: 1 / 1.414;
  min-height: 800px;
}
.document-content-html {
  line-height: 1.6;
  font-family: 'Times New Roman', Times, serif;
  color: #334155;
}
.document-content-html :deep(h1) { font-size: 1.5rem; margin-bottom: 1rem; }
.document-content-html :deep(p) { margin-bottom: 0.75rem; }
</style>
