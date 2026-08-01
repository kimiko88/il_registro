<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none">Dashboard Coordinatore di Classe</h1>
        <div class="text-subtitle2 text-grey-7">Panoramica andamento, note disciplinari, genitori e gestione scrutini</div>
      </div>
      <q-space />
      <q-select
        v-model="selectedClassId"
        :options="classOptions"
        option-value="id"
        option-label="label"
        emit-value
        map-options
        label="Seleziona Classe Coordinata"
        outlined
        dense
        style="min-width: 250px"
        @update:model-value="onClassChange"
      />
    </div>

    <div v-if="!selectedClassId" class="q-pa-xl text-center">
      <q-icon name="co_present" size="4rem" color="grey-5" />
      <div class="text-h6 text-grey-6 q-mt-md">Non risulti coordinatore di alcuna classe per l'anno scolastico in corso.</div>
    </div>

    <div v-else>
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="left"
        narrow-indicator
      >
        <q-tab name="academic" icon="trending_up" label="Andamento Classe" />
        <q-tab name="notes" icon="report_problem" label="Note & Richiami Disciplinari" />
        <q-tab name="guardians" icon="contacts" label="Anagrafica Genitori" />
        <q-tab name="scrutiny" icon="gavel" label="Scrutinio di Classe" />
      </q-tabs>

      <q-separator class="q-mb-md" />

      <q-tab-panels v-model="activeTab" animated class="bg-transparent">
        <!-- Tab 1: Andamento Classe -->
        <q-tab-panel name="academic" class="q-pa-none">
          <q-card class="shadow-2 rounded-borders q-pa-md">
            <div class="text-h6 text-weight-bold q-mb-md">Matrice Andamento Accademico</div>
            <q-spinner v-if="loadingMatrix" color="primary" size="2em" />
            <div v-else-if="matrixData">
              <q-markup-table flat bordered dense class="rounded-borders">
                <thead>
                  <tr class="bg-primary text-white">
                    <th class="text-left">Studente</th>
                    <th v-for="sub in matrixData.subjects" :key="sub.id" class="text-center">{{ sub.name }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="stu in matrixData.students" :key="stu.student_id">
                    <td class="text-weight-bold">{{ stu.student_name }}</td>
                    <td v-for="sub in matrixData.subjects" :key="sub.id" class="text-center">
                      <q-badge
                        :color="getAverageColor(stu.subject_data[sub.id]?.average)"
                        text-color="white"
                      >
                        {{ stu.subject_data[sub.id]?.average ? stu.subject_data[sub.id].average.toFixed(1) : '-' }}
                      </q-badge>
                    </td>
                  </tr>
                </tbody>
              </q-markup-table>
            </div>
          </q-card>
        </q-tab-panel>

        <!-- Tab 2: Note & Richiami Disciplinari -->
        <q-tab-panel name="notes" class="q-pa-none">
          <q-card class="shadow-2 rounded-borders q-pa-md">
            <div class="row items-center q-mb-md">
              <div class="text-h6 text-weight-bold">Note & Annotazioni dei Docenti della Classe</div>
              <q-space />
              <q-chip color="warning" text-color="white" icon="warning">
                {{ notesList.length }} provvedimenti totali
              </q-chip>
            </div>

            <q-list separator v-if="notesList.length > 0">
              <q-item v-for="n in notesList" :key="n.id">
                <q-item-section avatar>
                  <q-avatar
                    :color="n.type === 'disciplinary' ? 'negative' : 'warning'"
                    text-color="white"
                    :icon="n.type === 'disciplinary' ? 'report' : 'error_outline'"
                  />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">
                    {{ n.note }}
                  </q-item-label>
                  <q-item-label caption>
                    Inserita da: {{ n.teacher_name || 'Docente' }} • Data: {{ n.date }}
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip :color="n.is_approved ? 'positive' : 'orange'" text-color="white" dense>
                    {{ n.is_approved ? 'Convalidata' : 'In attesa approvazione' }}
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>
            <div v-else class="text-center text-grey-6 q-pa-lg">Nessuna nota disciplinare registrata per questa classe</div>
          </q-card>
        </q-tab-panel>

        <!-- Tab 3: Anagrafica Genitori -->
        <q-tab-panel name="guardians" class="q-pa-none">
          <q-card class="shadow-2 rounded-borders q-pa-md">
            <div class="text-h6 text-weight-bold q-mb-md">Contatti Genitori & Tutori Legali</div>
            <q-markup-table flat bordered dense class="rounded-borders">
              <thead>
                <tr class="bg-indigo-9 text-white">
                  <th class="text-left">Studente</th>
                  <th class="text-left">Genitore / Tutore</th>
                  <th class="text-left">Email</th>
                  <th class="text-left">Telefono</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="g in guardiansList" :key="g.guardian_id + g.student_id">
                  <td class="text-weight-bold">{{ g.student_name }}</td>
                  <td>{{ g.last_name }} {{ g.first_name }}</td>
                  <td><a :href="'mailto:' + g.email" class="text-primary">{{ g.email }}</a></td>
                  <td>{{ g.phone || 'N/D' }}</td>
                </tr>
              </tbody>
            </q-markup-table>
          </q-card>
        </q-tab-panel>

        <!-- Tab 4: Scrutinio di Classe -->
        <q-tab-panel name="scrutiny" class="q-pa-none">
          <q-card class="shadow-2 rounded-borders q-pa-md">
            <div class="row items-center q-mb-md">
              <div>
                <div class="text-h6 text-weight-bold">Gestione Flusso Scrutinio</div>
                <div class="text-caption text-grey-7">Avvia o controlla lo stato della sessione di scrutinio</div>
              </div>
              <q-space />
              <q-btn
                color="secondary"
                icon="play_arrow"
                label="Avvia Scrutinio"
                @click="startScrutiny"
              />
            </div>
          </q-card>
        </q-tab-panel>
      </q-tab-panels>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import api from '@/services/api';
import authService from '@/services/authService';
import { scrutinyService } from '@/services/scrutinyService';
import notesService from '@/services/notesService';

const $q = useQuasar();
const activeTab = ref('academic');
const selectedClassId = ref(null);
const classOptions = ref([]);

const loadingMatrix = ref(false);
const matrixData = ref(null);
const notesList = ref([]);
const guardiansList = ref([]);

const fetchClasses = async () => {
  try {
    const userRes = await authService.getCurrentUser();
    const currentUserId = userRes?.id;
    const res = await api.get('/teacher/classes');
    const assignedClasses = res.data || [];
    const coordClasses = assignedClasses.filter(c => c.coordinator_id === currentUserId);
    classOptions.value = coordClasses.length > 0 ? coordClasses : assignedClasses;
    if (classOptions.value.length > 0) {
      selectedClassId.value = classOptions.value[0].id;
      onClassChange(selectedClassId.value);
    }
  } catch (err) {
    classOptions.value = [];
  }
};

const onClassChange = async (classId) => {
  if (!classId) return;
  fetchMatrix(classId);
  fetchNotes(classId);
  fetchGuardians(classId);
};

const fetchMatrix = async (classId) => {
  loadingMatrix.value = true;
  try {
    const res = await scrutinyService.getMatrix(classId, 1);
    matrixData.value = res.data;
  } catch (err) {
    matrixData.value = null;
  } finally {
    loadingMatrix.value = false;
  }
};

const fetchNotes = async (classId) => {
  try {
    const res = await notesService.getNotes({ class_id: classId });
    notesList.value = res.data || [];
  } catch (err) {
    notesList.value = [];
  }
};

const fetchGuardians = async (classId) => {
  try {
    const res = await api.get(`/classes/${classId}/guardians`);
    guardiansList.value = res.data || [];
  } catch (err) {
    guardiansList.value = [];
  }
};

const startScrutiny = async () => {
  try {
    await scrutinyService.startScrutiny(selectedClassId.value, 1);
    $q.notify({ type: 'positive', message: 'Sessione di scrutinio avviata con successo' });
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante l\'avvio dello scrutinio' });
  }
};

const getAverageColor = (avg) => {
  if (!avg) return 'grey';
  if (avg >= 6.0) return 'positive';
  return 'negative';
};

onMounted(() => {
  fetchClasses();
});
</script>
