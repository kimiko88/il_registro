<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Sticky Header for Context -->
    <q-card class="sticky-header q-mb-md shadow-2 z-top">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6 q-mr-md">Gestione Voti</div>
        <q-select
          v-model="classesStore.selectedClassId"
          :options="classesStore.classes"
          option-value="id"
          option-label="name"
          label="Classe"
          dense outlined
          options-dense
          class="q-mr-md"
          style="min-width: 120px"
          emit-value
          map-options
        />
        <q-select
          v-model="selectedSubject"
          :options="gradesStore.subjects"
          label="Materia"
          dense outlined
          options-dense
          style="min-width: 150px"
          class="q-mr-md"
        />
        <q-space />
        <div class="row items-center q-gutter-sm">
             <q-btn-toggle
                v-model="viewMode"
                toggle-color="primary"
                flat
                :options="[
                    {label: 'Registro', value: 'table'},
                    {label: 'Statistiche', value: 'stats'},
                    {label: 'Storico', value: 'history'}
                ]"
             />
             <q-btn icon="file_upload" label="Importa CSV" outline color="primary" @click="showImportDialog = true" />
             <q-btn icon="print" flat round color="grey-8" @click="printReport" />
        </div>
      </q-card-section>
      <q-separator />
      <!-- Global Controls / Filters -->
      <q-card-section class="q-pt-sm q-pb-sm bg-grey-1" v-if="viewMode === 'table'">
          <div class="row items-center q-gutter-md">
             <q-input dense v-model="filterDate" type="date" label="Data Voto" outlined style="max-width: 150px" />
             <q-select
                dense
                v-model="gradeType"
                :options="['Orale', 'Scritto', 'Pratico']"
                label="Tipo Voto"
                outlined
                style="min-width: 120px"
             />
             <q-toggle v-model="showRubric" label="Mostra Rubrica" left-label dense />
             <q-toggle v-model="offlineMode" label="Offline Mode" color="amber" dense />
          </div>
      </q-card-section>
    </q-card>

    <!-- Main Content Area -->
    <div v-if="classesStore.selectedClassId">
        
        <!-- Table View -->
        <div v-show="viewMode === 'table'">
            <div class="row q-col-gutter-md">
                <div class="col-12" :class="{'col-md-9': showRubric, 'col-md-12': !showRubric}">
                     <GradeEntry
                        :class-id="classesStore.selectedClassId"
                        :subject="selectedSubject"
                        :date="filterDate"
                        :type="gradeType"
                        @refresh="refreshGrades"
                     />
                </div>
                <div class="col-12 col-md-3" v-if="showRubric">
                     <q-card class="bg-white">
                        <q-card-section class="bg-primary text-white text-subtitle2">Rubrica Valutazione</q-card-section>
                        <q-list separator dense>
                            <q-item><q-item-section><q-item-label>10 - Eccellente</q-item-label><q-item-label caption>Comprensione completa, esposizione brillante.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>9 - Ottimo</q-item-label><q-item-label caption>Comprensione approfondita, esposizione sicura.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>8 - Buono</q-item-label><q-item-label caption>Comprensione buona, esposizione corretta.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>7 - Discreto</q-item-label><q-item-label caption>Comprensione essenziale, qualche imprecisione.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>6 - Sufficiente</q-item-label><q-item-label caption>Conoscenze basilari raggiunte.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>5 - Insufficiente</q-item-label><q-item-label caption>Conoscenze frammentarie, errori gravi.</q-item-label></q-item-section></q-item>
                        </q-list>
                     </q-card>
                </div>
            </div>
        </div>

        <!-- Stats View -->
        <div v-show="viewMode === 'stats'">
             <GradeStatistics :class-id="classesStore.selectedClassId" :subject="selectedSubject" />
        </div>

    </div>
    <div v-else class="text-center q-pa-xl text-grey-6">
        <q-icon name="class" size="100px" />
        <div class="text-h5">Seleziona una classe per iniziare</div>
    </div>

    <!-- Import Dialog -->
    <q-dialog v-model="showImportDialog">
        <q-card style="min-width: 400px">
            <q-card-section class="text-h6">Importa Voti (CSV)</q-card-section>
            <q-card-section>
               <q-file outlined v-model="importFile" label="Seleziona file CSV" accept=".csv" />
               <div class="text-caption q-mt-sm">Formato: StudenteID, Voto, Data, Tipo, Note</div>
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Importa" @click="processImport" />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useClassesStore } from 'src/stores/classes';
import { useGradesStore } from 'src/stores/grades';
import GradeEntry from 'src/components/Teacher/GradeEntry.vue';
import GradeStatistics from 'src/components/Teacher/GradeStatistics.vue';
import { useQuasar, date } from 'quasar';

const $q = useQuasar();
const classesStore = useClassesStore();
const gradesStore = useGradesStore();

const selectedSubject = ref('Matematica');
const viewMode = ref('table');
const filterDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'));
const gradeType = ref('Orale');
const showRubric = ref(false);
const offlineMode = ref(false);

const showImportDialog = ref(false);
const importFile = ref(null);

watch(() => classesStore.selectedClassId, () => {
    refreshGrades();
});

const refreshGrades = async () => {
    if (!classesStore.selectedClassId) return;
    await gradesStore.fetchGrades(classesStore.selectedClassId, selectedSubject.value);
};

const printReport = () => {
    window.print();
};

const processImport = () => {
    $q.loading.show();
    setTimeout(() => {
        $q.loading.hide();
        showImportDialog.value = false;
        $q.notify({type: 'positive', message: 'Voti importati con successo (simulato)'});
        refreshGrades();
    }, 1500);
};

onMounted(() => {
    classesStore.fetchAssignedClasses();
});
</script>

<style scoped>
.sticky-header {
    position: sticky;
    top: 50px; /* Adjust based on navbar height */
    z-index: 100;
}
.z-top { z-index: 1000; }
</style>
