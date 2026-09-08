<template>
  <q-page class="q-pa-lg bg-slate-50 min-h-screen">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-ma-none">{{ t('certificatesPage.title') }}</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ t('certificatesPage.subtitle') }}
        </p>
      </div>
      <q-btn
        color="primary"
        icon="workspace_premium"
        :label="t('certificatesPage.generateNew')"
        no-caps
        class="shadow-sm rounded-lg q-px-md"
        @click="openGenerateDialog"
      />
    </div>

    <!-- Filters Bar -->
    <q-card class="bg-white rounded-xl shadow-sm q-mb-lg">
      <q-card-section class="row q-col-gutter-md items-center">
        <div class="col-12 col-md-3">
          <q-select
            v-model="filters.type"
            :options="typeFilterOptions"
            :label="t('certificatesPage.filterType')"
            outlined
            dense
            emit-value
            map-options
            @update:model-value="loadCertificates"
          />
        </div>
        <div class="col-12 col-md-4">
          <q-input
            v-model="filters.studentName"
            :label="t('certificatesPage.filterStudent')"
            outlined
            dense
            clearable
            @update:model-value="loadCertificates"
          >
            <template #append>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
        <div class="col-12 col-md-3">
          <q-input
            v-model="filters.academic_year"
            label="Anno Accademico"
            outlined
            dense
            @update:model-value="loadCertificates"
          />
        </div>
        <div class="col-12 col-md-2 text-right">
          <q-btn flat icon="refresh" :label="t('common.refresh') || 'Ricarica'" no-caps color="primary" @click="loadCertificates" />
        </div>
      </q-card-section>
    </q-card>

    <!-- Certificates Table -->
    <q-card class="bg-white rounded-xl shadow-sm">
      <q-card-section class="q-pa-none">
        <q-table
          :rows="filteredCertificates"
          :columns="columns"
          row-key="id"
          :loading="certStore.loading"
          flat
          :no-data-label="t('common.noData') || 'Nessun certificato emesso'"
        >
          <template #body-cell-type="props">
            <q-td :props="props">
              <q-chip
                :color="getTypeColor(props.value)"
                text-color="white"
                dense
                class="text-weight-bold"
              >
                {{ getTypeLabel(props.value) }}
              </q-chip>
            </q-td>
          </template>

          <template #body-cell-issued_at="props">
            <q-td :props="props">
              {{ formatDate(props.value) }}
            </q-td>
          </template>

          <template #body-cell-actions="props">
            <q-td :props="props" class="q-gutter-xs">
              <q-btn
                flat
                round
                dense
                color="primary"
                icon="picture_as_pdf"
                :title="t('common.download') || 'Scarica PDF'"
                :aria-label="t('common.download') || 'Scarica PDF'"
                @click="downloadCert(props.row.id)"
              />
              <q-btn
                flat
                round
                dense
                color="negative"
                icon="cancel"
                :title="t('common.delete') || 'Annulla Certificato'"
                :aria-label="t('common.delete') || 'Annulla Certificato'"
                @click="confirmDelete(props.row)"
              />
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog Genera Certificato (Multi-step) -->
    <q-dialog v-model="showGenerateDialog" persistent>
      <q-card style="width: min(650px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden bg-white shadow-24">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div class="text-h6 text-weight-bold row items-center">
            <q-icon name="workspace_premium" class="q-mr-xs" />
            {{ t('certificatesPage.generateNew') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-none">
          <q-stepper v-model="step" animated color="primary" flat>
            <!-- Step 1: Tipo Certificato -->
            <q-step :name="1" :title="t('certificatesPage.step1Title')" icon="bookmark" :done="step > 1">
              <div class="text-subtitle1 text-weight-bold q-mb-sm">{{ t('certificatesPage.step1Title') }}</div>
              <div class="q-gutter-sm">
                <q-card
                  v-for="opt in typeOptions"
                  :key="opt.value"
                  flat
                  bordered
                  class="cursor-pointer transition-all q-pa-md"
                  :class="form.type === opt.value ? 'bg-blue-50 border-primary shadow-xs' : 'hover:bg-slate-50'"
                  @click="form.type = opt.value"
                >
                  <div class="row items-center justify-between">
                    <div>
                      <div class="text-subtitle1 text-weight-bold text-slate-800">{{ opt.label }}</div>
                      <div class="text-caption text-slate-500">{{ opt.description }}</div>
                    </div>
                    <q-radio v-model="form.type" :val="opt.value" color="primary" />
                  </div>
                </q-card>
              </div>

              <div class="row justify-end q-mt-md">
                <q-btn color="primary" :label="t('common.next') || 'Avanti'" no-caps @click="step = 2" :disable="!form.type" />
              </div>
            </q-step>

            <!-- Step 2: Seleziona Studente -->
            <q-step :name="2" :title="t('certificatesPage.step2Title')" icon="person" :done="step > 2">
              <div class="text-subtitle1 text-weight-bold q-mb-sm">{{ t('certificatesPage.step2Title') }}</div>
              <q-select
                v-model="selectedStudent"
                :options="studentOptions"
                :label="t('certificatesPage.filterStudent')"
                outlined
                use-input
                input-debounce="300"
                @filter="filterStudents"
                class="q-mb-md"
              >
                <template #no-option>
                  <q-item>
                    <q-item-section class="text-grey">{{ t('common.noData') || 'Nessuno studente trovato' }}</q-item-section>
                  </q-item>
                </template>
              </q-select>

              <div v-if="selectedStudent" class="bg-slate-100 q-pa-md rounded-lg q-mb-md">
                <div class="text-weight-bold">{{ selectedStudent.label }}</div>
                <div class="text-caption text-slate-600">ID: {{ selectedStudent.value }}</div>
              </div>

              <div class="row justify-between q-mt-md">
                <q-btn flat :label="t('common.back') || 'Indietro'" color="slate-600" no-caps @click="step = 1" />
                <q-btn color="primary" :label="t('common.next') || 'Avanti'" no-caps @click="step = 3" :disable="!selectedStudent" />
              </div>
            </q-step>

            <!-- Step 3: Dettagli & Note -->
            <q-step :name="3" :title="t('certificatesPage.step3Title')" icon="tune" :done="step > 3">
              <div class="q-gutter-md">
                <q-input
                  v-model="form.academic_year"
                  label="Anno Accademico"
                  outlined
                  dense
                  hint="Es. 2025/2026"
                />
                <q-input
                  v-model="form.notes"
                  label="Note aggiuntive (opzionale)"
                  type="textarea"
                  rows="3"
                  outlined
                  dense
                />
              </div>

              <div class="row justify-between q-mt-lg">
                <q-btn flat :label="t('common.back') || 'Indietro'" color="slate-600" no-caps @click="step = 2" />
                <q-btn color="primary" :label="t('certificatesPage.summaryTitle') || 'Riepilogo'" no-caps @click="step = 4" />
              </div>
            </q-step>

            <!-- Step 4: Riepilogo & Generazione -->
            <q-step :name="4" :title="t('certificatesPage.step4Title')" icon="check_circle">
              <div class="bg-slate-50 q-pa-lg rounded-xl bordered q-mb-md">
                <div class="row q-col-gutter-sm">
                  <div class="col-6 text-caption text-slate-500">{{ t('certificatesPage.colType') }}:</div>
                  <div class="col-6 text-weight-bold">{{ getTypeLabel(form.type) }}</div>

                  <div class="col-6 text-caption text-slate-500">{{ t('certificatesPage.colStudent') }}:</div>
                  <div class="col-6 text-weight-bold">{{ selectedStudent?.label }}</div>

                  <div class="col-6 text-caption text-slate-500">Anno Accademico:</div>
                  <div class="col-6 text-weight-bold">{{ form.academic_year }}</div>

                  <div v-if="form.notes" class="col-6 text-caption text-slate-500">Note:</div>
                  <div v-if="form.notes" class="col-6 text-weight-bold">{{ form.notes }}</div>
                </div>
              </div>

              <div class="row justify-between q-mt-lg">
                <q-btn flat :label="t('common.back') || 'Indietro'" color="slate-600" no-caps @click="step = 3" />
                <q-btn
                  color="positive"
                  icon="picture_as_pdf"
                  :label="t('certificatesPage.generateAndDownload')"
                  no-caps
                  class="q-px-lg shadow-sm"
                  :loading="generating"
                  @click="handleGenerate"
                />
              </div>
            </q-step>
          </q-stepper>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar, date } from 'quasar';
import { useCertificatesStore } from '@/stores/certificates';
import api from '@/services/api';

const $q = useQuasar();
const { t } = useI18n();
const certStore = useCertificatesStore();

const showGenerateDialog = ref(false);
const step = ref(1);
const generating = ref(false);

const filters = reactive({
  type: 'all',
  studentName: '',
  academic_year: '2025/2026'
});

const form = reactive({
  type: 'iscrizione',
  student_id: '',
  academic_year: '2025/2026',
  notes: ''
});

const selectedStudent = ref(null);
const rawStudents = ref([]);

const typeFilterOptions = computed(() => [
  { label: t('certificatesPage.allTypes'), value: 'all' },
  { label: t('certificatesPage.types.iscrizione'), value: 'iscrizione' },
  { label: t('certificatesPage.types.frequenza'), value: 'frequenza' },
  { label: t('certificatesPage.types.promozione'), value: 'promozione' },
  { label: t('certificatesPage.types.condotta'), value: 'condotta' }
]);

const typeOptions = computed(() => [
  { value: 'iscrizione', label: t('certificatesPage.types.iscrizione'), description: t('certificatesPage.descriptions.iscrizione') },
  { value: 'frequenza', label: t('certificatesPage.types.frequenza'), description: t('certificatesPage.descriptions.frequenza') },
  { value: 'promozione', label: t('certificatesPage.types.promozione'), description: t('certificatesPage.descriptions.promozione') },
  { value: 'condotta', label: t('certificatesPage.types.condotta'), description: t('certificatesPage.descriptions.condotta') }
]);

const columns = computed(() => [
  { name: 'protocol_no', label: t('certificatesPage.colProtocol'), field: 'protocol_no', align: 'left', sortable: true },
  { name: 'student_name', label: t('certificatesPage.colStudent'), field: 'student_name', align: 'left', sortable: true },
  { name: 'class_name', label: t('certificatesPage.colClass'), field: 'class_name', align: 'left', sortable: true },
  { name: 'type', label: t('certificatesPage.colType'), field: 'type', align: 'center', sortable: true },
  { name: 'issued_at', label: t('certificatesPage.colIssuedAt'), field: 'issued_at', align: 'left', sortable: true },
  { name: 'issued_by_name', label: t('certificatesPage.colIssuedBy'), field: 'issued_by_name', align: 'left' },
  { name: 'actions', label: t('certificatesPage.colActions'), field: 'actions', align: 'center' }
]);

const filteredCertificates = computed(() => {
  let list = certStore.certificates;
  if (filters.studentName) {
    const q = filters.studentName.toLowerCase();
    list = list.filter(c => (c.student_name || '').toLowerCase().includes(q));
  }
  return list;
});

onMounted(() => {
  loadCertificates();
  loadStudents();
});

const loadCertificates = async () => {
  try {
    await certStore.fetchCertificates(filters);
  } catch (e) {
    $q.notify({ type: 'negative', message: t('certificatesPage.loadError') || 'Errore nel caricamento dei certificati' });
  }
};

const loadStudents = async () => {
  try {
    const res = await api.get('/users?role=student');
    const usersList = res.data?.users || res.data || [];
    rawStudents.value = usersList.map(u => ({
      label: `${u.last_name} ${u.first_name} (${u.class_name || 'N.D.'})`,
      value: u.id
    }));
  } catch (e) {
    console.error('Error fetching students:', e);
  }
};

const studentOptions = ref([]);
const filterStudents = (val, update) => {
  update(() => {
    if (val === '') {
      studentOptions.value = rawStudents.value.slice(0, 50);
    } else {
      const needle = val.toLowerCase();
      studentOptions.value = rawStudents.value.filter(v => v.label.toLowerCase().includes(needle)).slice(0, 50);
    }
  });
};

const openGenerateDialog = () => {
  step.value = 1;
  form.type = 'iscrizione';
  form.academic_year = '2025/2026';
  form.notes = '';
  selectedStudent.value = null;
  showGenerateDialog.value = true;
};

const handleGenerate = async () => {
  if (!selectedStudent.value) return;
  generating.value = true;
  try {
    const payload = {
      student_id: selectedStudent.value.value,
      type: form.type,
      academic_year: form.academic_year,
      notes: form.notes
    };
    const res = await certStore.generateCertificate(payload);
    $q.notify({ type: 'positive', message: t('certificatesPage.generateSuccess') || 'Certificato generato con successo!' });
    showGenerateDialog.value = false;

    if (res?.id) {
      await certStore.downloadPDF(res.id);
    }
  } catch (e) {
    $q.notify({ type: 'negative', message: t('certificatesPage.generateError') || 'Errore durante la generazione del certificato' });
  } finally {
    generating.value = false;
  }
};

const downloadCert = async (id) => {
  try {
    await certStore.downloadPDF(id);
    $q.notify({ type: 'positive', message: 'Download PDF avviato' });
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante il download del PDF' });
  }
};

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Annulla Certificato',
    message: `Sei sicuro di voler annullare il certificato ${row.protocol_no} di ${row.student_name}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await certStore.deleteCertificate(row.id);
      $q.notify({ type: 'positive', message: 'Certificato annullato' });
    } catch (e) {
      $q.notify({ type: 'negative', message: 'Errore durante l\'annullamento' });
    }
  });
};

const getTypeLabel = (tVal) => {
  return t('certificatesPage.types.' + tVal) || tVal;
};

const getTypeColor = (t) => {
  switch (t) {
    case 'iscrizione': return 'primary';
    case 'frequenza': return 'teal';
    case 'promozione': return 'positive';
    case 'condotta': return 'indigo';
    default: return 'grey-7';
  }
};

const formatDate = (d) => {
  if (!d) return '-';
  return date.formatDate(d, 'DD/MM/YYYY HH:mm');
};
</script>
