<template>
  <q-dialog
    :model-value="modelValue"
    persistent
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <q-card style="width: min(750px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden bg-white shadow-24">
      <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
        <div class="text-h6 text-weight-bold">
          <q-icon name="upload_file" class="q-mr-xs" />
          {{ t('usersPage.importTitle') || 'Importazione Utenti da CSV' }}
        </div>
        <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-pa-none">
        <q-stepper v-model="importStep" animated color="primary" flat>
          <!-- Step 1: Configurazione -->
          <q-step :name="1" :title="t('usersPage.stepConfig') || 'Configurazione'" icon="settings" :done="importStep > 1">
            <div class="space-y-4 q-pa-md">
              <q-select
                v-model="importRole"
                :options="roleOptions"
                :label="t('usersPage.userTypeToImport') || 'Tipo Utenti da Importare'"
                outlined
                dense
                emit-value
                map-options
              />

              <q-select
                v-if="importRole === 'student'"
                v-model="importClassId"
                :options="classOptions"
                :label="t('usersPage.targetClass') || 'Classe Destinazione'"
                outlined
                dense
                emit-value
                map-options
              />

              <q-toggle
                v-model="sendWelcomeEmail"
                :label="t('usersPage.sendWelcomeEmail') || 'Invia email di benvenuto con credenziali temporanee'"
              />

              <div class="q-pt-sm">
                <q-btn
                  flat
                  color="primary"
                  icon="download"
                  :label="t('usersPage.downloadTemplate') || 'Scarica Template CSV Esempio'"
                  no-caps
                  @click="downloadCsvTemplate"
                />
              </div>
            </div>

            <div class="row justify-end q-pa-md border-t border-slate-200">
              <q-btn color="primary" :label="t('common.next') || 'Avanti'" no-caps @click="importStep = 2" />
            </div>
          </q-step>

          <!-- Step 2: Upload File & Preview -->
          <q-step :name="2" :title="t('usersPage.stepUpload') || 'Caricamento File'" icon="cloud_upload" :done="importStep > 2">
            <div class="q-pa-md space-y-4">
              <q-file
                v-model="importFile"
                :label="t('usersPage.selectCsvFile') || 'Seleziona File CSV (.csv max 10MB)'"
                outlined
                dense
                accept=".csv"
                @update:model-value="onCsvFileSelected"
              >
                <template v-slot:prepend>
                  <q-icon name="attach_file" />
                </template>
              </q-file>

              <div v-if="previewRows.length > 0" class="q-mt-md">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-xs">
                  {{ t('usersPage.previewRowsTitle') || 'Anteprima prime 5 righe' }}
                </div>
                <q-table
                  flat
                  dense
                  bordered
                  :rows="previewRows"
                  hide-pagination
                  :pagination="{ rowsPerPage: 0 }"
                />
              </div>
            </div>

            <div class="row justify-between q-pa-md border-t border-slate-200">
              <q-btn flat :label="t('common.back') || 'Indietro'" color="slate-600" no-caps @click="importStep = 1" />
              <q-btn
                color="primary"
                :label="t('usersPage.runImport') || 'Esegui Importazione'"
                :loading="importing"
                :disable="!importFile"
                no-caps
                @click="runBulkImport"
              />
            </div>
          </q-step>

          <!-- Step 3: Risultato -->
          <q-step :name="3" :title="t('usersPage.stepResult') || 'Risultato'" icon="check_circle">
            <div class="q-pa-md space-y-4">
              <q-banner class="bg-positive text-white rounded-lg">
                <template v-slot:avatar>
                  <q-icon name="check_circle" size="md" />
                </template>
                <div class="text-weight-bold">
                  {{ t('usersPage.importSuccessCount', { count: importResult?.imported || 0 }) || `${importResult?.imported || 0} Utenti Importati Con Successo!` }}
                </div>
                <div v-if="importResult?.skipped > 0" class="text-caption">
                  {{ t('usersPage.skippedRowsCount', { count: importResult.skipped }) || `${importResult.skipped} righe saltate o non valide.` }}
                </div>
              </q-banner>

              <div v-if="importResult?.errors && importResult.errors.length > 0">
                <div class="text-subtitle2 text-weight-bold text-negative q-mb-xs">
                  {{ t('usersPage.errorsEncountered') || 'Errori Riscontrati' }}
                </div>
                <q-list bordered separator class="rounded-lg bg-red-50">
                  <q-item v-for="(err, idx) in importResult.errors" :key="idx" dense>
                    <q-item-section avatar>
                      <q-icon name="warning" color="negative" size="xs" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="text-caption text-negative">
                        {{ t('usersPage.rowError', { row: err.row, reason: err.reason }) || `Riga ${err.row}: ${err.reason}` }}
                      </q-item-label>
                    </q-item-section>
                  </q-item>
                </q-list>
              </div>
            </div>

            <div class="row justify-between q-pa-md border-t border-slate-200">
              <q-btn
                v-if="importResult?.errors?.length > 0"
                flat
                color="negative"
                icon="download"
                :label="t('usersPage.downloadErrorReport') || 'Scarica Report Errori'"
                no-caps
                @click="downloadErrorReport"
              />
              <div v-else />

              <q-btn
                color="primary"
                :label="t('usersPage.closeAndUpdate') || 'Chiudi e Aggiorna Lista'"
                no-caps
                @click="closeImportModal"
              />
            </div>
          </q-step>
        </q-stepper>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar, exportFile } from 'quasar';
import { useUsersStore } from '@/stores/users';

defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  classOptions: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['update:modelValue', 'imported']);

const $q = useQuasar();
const { t } = useI18n();
const usersStore = useUsersStore();

const importStep = ref(1);
const importRole = ref('student');
const importClassId = ref(null);
const sendWelcomeEmail = ref(false);
const importFile = ref(null);
const previewRows = ref([]);
const importResult = ref(null);
const importing = ref(false);

const roleOptions = computed(() => [
  { label: t('usersPage.roleStudents') || 'Studenti', value: 'student' },
  { label: t('usersPage.roleTeachers') || 'Docenti', value: 'teacher' },
  { label: t('usersPage.roleParents') || 'Genitori', value: 'parent' }
]);

const downloadCsvTemplate = () => {
  let templateHeader = 'cognome,nome,codice_fiscale,data_nascita,email,classe\nRossi,Mario,RSSMRA00A01H501Z,2008-01-01,mario@example.com,1A';
  if (importRole.value === 'teacher') {
    templateHeader = 'cognome,nome,codice_fiscale,email,materie\nBianchi,Anna,BNCNNA80B41H501X,anna@example.com,"Matematica,Fisica"';
  } else if (importRole.value === 'parent') {
    templateHeader = 'cognome,nome,codice_fiscale,email\nVerdi,Giuseppe,VRDGPP75A01H501Y,giuseppe@example.com';
  }
  exportFile(`template_import_${importRole.value}.csv`, templateHeader, 'text/csv');
};

const onCsvFileSelected = (file) => {
  if (!file) {
    previewRows.value = [];
    return;
  }
  const reader = new FileReader();
  reader.onload = (e) => {
    const text = e.target?.result || '';
    const lines = text.split(/\r?\n/).filter(l => l.trim() !== '');
    previewRows.value = lines.slice(1, 6).map((line, idx) => {
      const parts = line.split(',');
      return {
        riga: idx + 2,
        cognome: parts[0] || '',
        nome: parts[1] || '',
        codice_fiscale: parts[2] || '',
        email: parts[3] || parts[4] || ''
      };
    });
  };
  reader.readAsText(file);
};

const runBulkImport = async () => {
  if (!importFile.value) return;
  importing.value = true;
  try {
    const formData = new FormData();
    formData.append('file', importFile.value);
    formData.append('role', importRole.value);
    if (importClassId.value) formData.append('class_id', importClassId.value);
    formData.append('send_welcome_email', sendWelcomeEmail.value.toString());

    importResult.value = await usersStore.bulkImportUsers(formData);
    importStep.value = 3;
  } catch (e) {
    console.error('Error during CSV import:', e);
    $q.notify({
      type: 'negative',
      message: t('usersPage.importError') || 'Errore durante l\'importazione CSV'
    });
  } finally {
    importing.value = false;
  }
};

const downloadErrorReport = () => {
  if (!importResult.value?.errors) return;
  const errContent = 'Riga,Motivo Errore\n' + importResult.value.errors.map(e => `${e.row},"${e.reason}"`).join('\n');
  exportFile('report_errori_import.csv', errContent, 'text/csv');
};

const closeImportModal = () => {
  emit('update:modelValue', false);
  importStep.value = 1;
  importFile.value = null;
  previewRows.value = [];
  importResult.value = null;
  emit('imported');
};
</script>
