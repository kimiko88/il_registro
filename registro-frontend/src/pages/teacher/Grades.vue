<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Sticky Header for Context -->
    <q-card class="sticky-header q-mb-md shadow-2">
      <q-card-section class="row items-center justify-between wrap q-gutter-sm q-pb-none">

        <div class="text-h6 q-mr-md">Gestione Voti</div>
        <q-select
          v-model="selectedClassId"
          :options="classesStore.classes"
          option-value="id"
          option-label="label"
          label="Classe"
          dense outlined
          options-dense
          class="q-mr-md"
          style="min-width: 200px"
          emit-value
          map-options
          :loading="classesStore.loading"
        />
        <q-select
          v-model="selectedSubject"
          :options="gradesStore.subjects"
          option-label="subject_name"
          option-value="subject_id"
          emit-value
          map-options
          label="Materia"
          dense outlined
          options-dense
          style="min-width: 150px"
          class="q-mr-md"
          @update:model-value="refreshGrades"
        />
        <q-space />
        <div class="row items-center q-gutter-sm">
             <q-btn-toggle
                v-model="viewMode"
                toggle-color="primary"
                flat
                aria-label="Modalità visualizzazione"
                :options="[
                    {label: 'Registro', value: 'table'},
                    {label: 'Statistiche', value: 'stats'},
                    {label: 'Storico', value: 'history'}
                ]"
             />
             <q-btn icon="assignment" label="Nuova Verifica" color="primary" class="q-mr-sm" @click="openTestDialog" :disable="!selectedClassId || !selectedSubject || !isAssignedClass" />
             <q-btn icon="file_upload" label="Importa CSV" outline color="primary" @click="showImportDialog = true" :disable="!isAssignedClass" />
             <q-btn icon="print" flat round color="grey-8" aria-label="Stampa registro voti" @click="printReport">
               <q-tooltip>Stampa Registro Voti</q-tooltip>
             </q-btn>
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
          </div>
      </q-card-section>
    </q-card>

    <!-- Substitution Non-Titolare Warning Banner -->
    <q-banner v-if="selectedClassId && !isAssignedClass" class="bg-amber-1 text-amber-10 rounded-xl border border-amber-300 q-mb-md shadow-soft">
      <template v-slot:avatar>
        <q-icon name="lock" color="amber-9" size="28px" />
      </template>
      <div class="text-weight-bold text-subtitle1">🔒 Modalità Supplenza - Consultazione ed Inserimento Voti Disabilitati</div>
      <div class="text-caption">
        Questa classe non rientra tra le tue cattedre titolari. In qualità di docente supplente puoi firmare l'ora nel <strong>Registro Presenze</strong> ed appuntare presenze e note disciplinari, mentre l'inserimento ed la visione dei voti è riservata ai docenti titolari.
      </div>
    </q-banner>

    <!-- Main Content Area -->
    <div v-if="selectedClassId">
        
        <!-- Table View -->
        <div v-show="viewMode === 'table'">
            <div class="row q-col-gutter-md">
                <div class="col-12" :class="{'col-md-9': showRubric, 'col-md-12': !showRubric}">
                     <GradeEntry
                        :class-id="selectedClassId"
                        :subject="selectedSubject"
                        :date="filterDate"
                        :type="gradeType"
                        :read-only="!isAssignedClass"
                        @refresh="refreshGrades"
                     />
                </div>
                <div class="col-12 col-md-3" v-if="showRubric && isAssignedClass">
                     <q-card class="bg-white">
                        <q-card-section class="bg-primary text-white text-subtitle2">Rubrica Valutazione</q-card-section>
                        <q-list separator dense>
                            <q-item><q-item-section><q-item-label>10 - Eccellente</q-item-label><q-item-label caption>Comprensione completa, esposizione brillante.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>9 - Ottimo</q-item-label><q-item-label caption>Comprensione approfondita, esposizione sicura.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>8 - Buono</q-item-label><q-item-label caption>Comprensione buona, esposizione corretta.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>7 - Discreto</q-item-label><q-item-label caption>Comprensione essenziale, qualche imprecisione.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>6 - Sufficiente</q-item-label><q-item-label caption>Conoscenze basilari raggiunte.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>5 - Insufficiente</q-item-label><q-item-label caption>Conoscenze frammentarie, errori rilevanti.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>4 - Gravemente Insufficiente</q-item-label><q-item-label caption>Lacune diffuse, scarsa autonomia.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>3 - Molto Scarso</q-item-label><q-item-label caption>Gravi lacune concettuali non colmate.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>2 - Nullo / Non Svolto</q-item-label><q-item-label caption>Compito non eseguito o nullo.</q-item-label></q-item-section></q-item>
                            <q-item><q-item-section><q-item-label>1 - Non Classificabile</q-item-label><q-item-label caption>Assenza totale di contenuti o consegna in bianco.</q-item-label></q-item-section></q-item>
                        </q-list>
                     </q-card>
                </div>
            </div>
        </div>

        <!-- Stats View -->
        <div v-show="viewMode === 'stats'">
             <GradeStatistics v-if="isAssignedClass" :class-id="selectedClassId" :subject="selectedSubject" />
             <div v-else class="text-center q-pa-xl text-grey-7">Le statistiche voti sono riservate ai docenti titolari.</div>
        </div>

        <!-- History / Verifiche View -->
        <div v-show="viewMode === 'history'">
            <q-card v-if="isAssignedClass">
                <q-card-section class="row items-center justify-between">
                    <div class="text-h6 text-primary text-weight-bold">Elenco Verifiche</div>
                    <q-btn icon="refresh" flat round dense aria-label="Aggiorna elenco verifiche" @click="fetchTests" />
                </q-card-section>
                <q-separator />
                <q-card-section v-if="loadingTests" class="q-pa-md">
                    <SkeletonTable :rows="6" :cols="4" />
                </q-card-section>
                <q-card-section v-else-if="classTests.length === 0" class="text-center q-pa-xl text-slate-600">
                    <q-icon name="quiz" size="64px" color="grey-5" class="q-mb-md" />
                    <div class="text-h6 text-weight-bold">Nessuna verifica trovata</div>
                    <div class="text-caption text-grey-7 q-mb-md">Non ci sono ancora verifiche o prove registrate per questa materia.</div>
                    <q-btn icon="add" label="Crea la Prima Verifica" color="primary" unelevated no-caps @click="openTestDialog" />
                </q-card-section>
                <q-list separator v-else>
                    <q-item v-for="test in classTests" :key="test.id" class="q-py-md">
                        <q-item-section>
                            <q-item-label class="text-weight-bold text-subtitle1">{{ test.title }}</q-item-label>
                            <q-item-label caption class="row items-center q-gutter-x-sm">
                                <span>Data: {{ formatDate(test.date) }}</span>
                                <span>|</span>
                                <span>Tipo: {{ test.evaluation_type === 'Written' ? 'Scritto' : (test.evaluation_type === 'Oral' ? 'Orale' : 'Pratico') }}</span>
                            </q-item-label>
                            <q-item-label class="text-caption text-grey-8 q-mt-xs" v-if="test.teacher_notes">
                                <strong>Note Docente:</strong> {{ test.teacher_notes }}
                            </q-item-label>
                            <q-item-label class="text-caption text-grey-8" v-if="test.parent_notes">
                                <strong>Note Genitori:</strong> {{ test.parent_notes }}
                            </q-item-label>
                        </q-item-section>
                        <q-item-section side>
                            <div class="row q-gutter-sm">
                                <q-btn icon="edit" label="Modifica" outline color="primary" @click="openEditTestDialog(test)" />
                                <q-btn icon="delete" label="Elimina" outline color="negative" @click="deleteTestConfirm(test.id)" />
                            </div>
                        </q-item-section>
                    </q-item>
                </q-list>
            </q-card>
            <div v-else class="text-center q-pa-xl text-grey-7">La cronologia verifiche è riservata ai docenti titolari.</div>
        </div>

    </div>
    <div v-else class="text-center q-pa-xl column items-center">
        <q-icon name="school" size="96px" color="primary" class="q-mb-md opacity-80" />
        <div class="text-h5 text-weight-bold q-mb-xs">Seleziona una classe per iniziare</div>
        <div class="text-subtitle2 text-grey-7 q-mb-lg" style="max-width: 480px;">
          Scegli una classe dal menu in alto per accedere al registro:
        </div>
        <div class="row q-gutter-sm justify-center" v-if="classesStore.classes && classesStore.classes.length > 0">
          <q-btn
            v-for="cls in classesStore.classes"
            :key="cls.id"
            unelevated
            color="primary"
            outline
            :label="cls.name"
            icon="class"
            class="q-px-md"
            @click="selectedClassId = cls.id"
          />
        </div>
    </div>

    <!-- Import Dialog -->
    <q-dialog v-model="showImportDialog">
        <q-card style="width: min(450px, 95vw)">
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

    <!-- Create Class Test Dialog (Bulk) -->
    <q-dialog v-model="showTestDialog" persistent max-width="96vw">
      <q-card style="width: 1200px; max-width: 96vw; max-height: 92vh; overflow-y: auto;">
        <q-card-section class="bg-primary text-white row items-center">
          <div class="text-h6 text-weight-bold">Crea Nuova Verifica</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-form @submit.prevent="submitTest">
            <div class="row q-col-gutter-lg">
              <!-- Test Details -->
              <div class="col-12 col-md-4">
                <div class="text-subtitle1 q-mb-md text-weight-bold text-primary">Dettagli Verifica</div>
                <q-input
                  v-model="testForm.title"
                  label="Titolo Verifica *"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                  class="q-mb-sm"
                />
                <q-input
                  v-model="testForm.date"
                  type="date"
                  label="Data Verifica *"
                  outlined
                  dense
                  class="q-mb-sm"
                />

                <!-- Overlapping Test Alert -->
                <q-banner v-if="overlappingTestsCount >= 2" rounded dense class="bg-amber-1 text-amber-9 border border-amber-3 q-mb-sm">
                  <template v-slot:avatar>
                    <q-icon name="warning" color="amber-9" />
                  </template>
                  Attenzione: La classe ha già <strong>{{ overlappingTestsCount }} verifiche</strong> in programma il {{ testForm.date }}! (Max raccomandato: 2)
                </q-banner>
                <q-select
                  v-model="testForm.evaluationType"
                  :options="['Scritto', 'Orale', 'Pratico']"
                  label="Tipo Valutazione *"
                  outlined
                  dense
                  class="q-mb-md"
                />
                <q-input
                  v-model="testForm.teacherNotes"
                  type="textarea"
                  label="Testo per il docente (Note Interne)"
                  outlined
                  dense
                  rows="3"
                  class="q-mb-sm"
                />
                <q-input
                  v-model="testForm.parentNotes"
                  type="textarea"
                  label="Testo visualizzato dai genitori"
                  outlined
                  dense
                  rows="3"
                  class="q-mb-sm"
                />
              </div>

              <!-- Student Grades -->
              <div class="col-12 col-md-8">
                <div class="text-subtitle1 q-mb-xs text-weight-bold text-primary row items-center justify-between">
                  <div>Voti Alunni ({{ filledTestGradesCount }}/{{ testForm.grades.length }} inseriti)</div>
                  <div class="row items-center q-gutter-x-xs">
                    <q-btn icon="block" size="sm" outline color="warning" label="Segna tutti assenti" @click="markAllAbsent(testForm)" />
                  </div>
                </div>
                <q-linear-progress :value="testForm.grades.length ? filledTestGradesCount / testForm.grades.length : 0" color="primary" class="q-mb-sm" />
                
                <q-scroll-area style="height: 420px;" tabindex="0" aria-label="Lista inserimento voti alunni" class="border-grey rounded-borders q-pa-sm bg-grey-2">
                  <q-list separator>
                    <q-item v-for="(student, idx) in testForm.grades" :key="student.student_id" class="q-py-sm">
                      <q-item-section>
                        <q-item-label class="text-weight-bold">{{ student.full_name }}</q-item-label>
                      </q-item-section>
                      <q-item-section side style="width: 390px">
                        <div class="row items-center q-gutter-sm no-wrap">
                          <q-select
                            v-model="student.grade_value"
                            :options="gradeOptions"
                            emit-value
                            map-options
                            label="Voto"
                            outlined
                            dense
                            style="width: 110px"
                            :bg-color="getGradeColor(student.grade_value)"
                            placeholder="-"
                            :ref="el => setGradeInputRef(el, idx)"
                            @keydown.enter.prevent="focusNextStudent(idx)"
                          />
                          <q-input
                            v-model="student.notes"
                            label="Note personali"
                            outlined
                            dense
                            class="col"
                            placeholder="Note..."
                            @keydown.enter.prevent="focusNextStudent(idx)"
                          />
                        </div>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-scroll-area>
              </div>
            </div>

            <q-card-actions align="right" class="q-mt-md q-px-none">
              <q-btn flat label="Annulla" v-close-popup color="grey-7" />
              <q-btn type="submit" label="Salva Verifica e Voti" color="primary" :loading="loading" />
            </q-card-actions>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Edit Class Test Dialog (Bulk) -->
    <q-dialog v-model="showEditTestDialog" persistent max-width="96vw">
      <q-card style="width: 1200px; max-width: 96vw; max-height: 92vh; overflow-y: auto;">
        <q-card-section class="bg-primary text-white row items-center">
          <div class="text-h6 text-weight-bold">Modifica Verifica in Blocco</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-form @submit.prevent="submitEditTest">
            <div class="row q-col-gutter-lg">
              <!-- Test Details -->
              <div class="col-12 col-md-4 q-gutter-y-md">
                <q-input
                  v-model="editTestForm.title"
                  label="Titolo Verifica"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
                <q-input
                  v-model="editTestForm.date"
                  type="date"
                  label="Data Verifica"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
                <q-select
                  v-model="editTestForm.evaluationType"
                  :options="['Scritto', 'Orale', 'Pratico']"
                  label="Tipo Valutazione"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
                <q-input
                  v-model="editTestForm.teacherNotes"
                  type="textarea"
                  label="Note Interne (solo docente)"
                  outlined
                  dense
                  rows="3"
                />
                <q-input
                  v-model="editTestForm.parentNotes"
                  type="textarea"
                  label="Note per i Genitori (visibili in bacheca)"
                  outlined
                  dense
                  rows="3"
                />
              </div>

              <!-- Student Grades -->
              <div class="col-12 col-md-8">
                <div class="text-subtitle1 q-mb-xs text-weight-bold text-primary row items-center justify-between">
                  <div>Voti Alunni ({{ filledEditTestGradesCount }}/{{ editTestForm.grades.length }} inseriti)</div>
                  <div class="row items-center q-gutter-x-xs">
                    <q-btn icon="block" size="sm" outline color="warning" label="Segna tutti assenti" @click="markAllAbsent(editTestForm)" />
                  </div>
                </div>
                <q-linear-progress :value="editTestForm.grades.length ? filledEditTestGradesCount / editTestForm.grades.length : 0" color="primary" class="q-mb-sm" />
                
                <q-scroll-area style="height: 420px;" tabindex="0" aria-label="Lista modifica voti alunni" class="border-grey rounded-borders q-pa-sm bg-grey-2">
                  <q-list separator>
                    <q-item v-for="(student, idx) in editTestForm.grades" :key="student.student_id" class="q-py-sm">
                      <q-item-section>
                        <q-item-label class="text-weight-bold">{{ student.full_name }}</q-item-label>
                      </q-item-section>
                      <q-item-section side style="width: 390px">
                        <div class="row items-center q-gutter-sm no-wrap">
                          <q-select
                            v-model="student.grade_value"
                            :options="gradeOptions"
                            emit-value
                            map-options
                            label="Voto"
                            outlined
                            dense
                            style="width: 110px"
                            :bg-color="getGradeColor(student.grade_value)"
                            placeholder="-"
                            :ref="el => setEditGradeInputRef(el, idx)"
                            @keydown.enter.prevent="focusNextEditStudent(idx)"
                          />
                          <q-input
                            v-model="student.notes"
                            label="Note personali"
                            outlined
                            dense
                            class="col"
                            placeholder="Note..."
                            @keydown.enter.prevent="focusNextEditStudent(idx)"
                          />
                        </div>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-scroll-area>
              </div>
            </div>

            <q-card-actions align="right" class="q-mt-md q-px-none">
              <q-btn flat label="Annulla" v-close-popup color="grey-7" />
              <q-btn type="submit" label="Salva Modifiche Verifica" color="primary" :loading="loading" />
            </q-card-actions>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useClassesStore } from 'src/stores/classes';
import { useGradesStore } from 'src/stores/grades';
import GradeEntry from 'src/components/Teacher/GradeEntry.vue';
import GradeStatistics from 'src/components/Teacher/GradeStatistics.vue';
import { gradeService } from 'src/services/gradeService';
import { useQuasar, date } from 'quasar';
import SkeletonTable from '@/components/Common/SkeletonTable.vue';
import { useUndoToast } from '@/composables/useUndoToast';
import { useSchoolYearStore } from '@/stores/schoolYear';
import { ITALIAN_GRADE_OPTIONS, gradeToNumeric, formatGrade, getGradeColor } from '@/utils/gradeUtils';

const $q = useQuasar();
const { notifyWithUndo } = useUndoToast();
const classesStore = useClassesStore();
const gradesStore = useGradesStore();
const schoolYearStore = useSchoolYearStore();

const selectedClassId = ref(null);
const selectedSubject = ref(null); 
const viewMode = ref('table');
const filterDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'));
const gradeType = ref('Orale');
const showRubric = ref(false);

const gradeOptions = ITALIAN_GRADE_OPTIONS;

const isAssignedClass = computed(() => {
  if (!selectedClassId.value) return false;
  return classesStore.classes.some(c => String(c.id) === String(selectedClassId.value));
});

const showImportDialog = ref(false);
const importFile = ref(null);

const showTestDialog = ref(false);
const loading = ref(false);
const testForm = ref({
    title: '',
    date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
    evaluationType: 'Scritto',
    teacherNotes: '',
    parentNotes: '',
    grades: []
});

const filledTestGradesCount = computed(() => {
    if (!testForm.value || !testForm.value.grades) return 0;
    return testForm.value.grades.filter(g => g.grade_value !== null && g.grade_value !== undefined && g.grade_value !== '').length;
});

const filledEditTestGradesCount = computed(() => {
    if (!editTestForm.value || !editTestForm.value.grades) return 0;
    return editTestForm.value.grades.filter(g => g.grade_value !== null && g.grade_value !== undefined && g.grade_value !== '').length;
});

const overlappingTestsCount = computed(() => {
    if (!testForm.value?.date || !classTests.value) return 0;
    return classTests.value.filter(t => t.date && t.date.startsWith(testForm.value.date)).length;
});

const openTestDialog = () => {
    if (!isAssignedClass.value) {
        $q.notify({ type: 'warning', message: 'Non hai i permessi per inserire verifiche in una classe non tua' });
        return;
    }
    if (!gradesStore.grades || !gradesStore.grades.students || gradesStore.grades.students.length === 0) {
        $q.notify({ type: 'warning', message: 'Nessun alunno caricato per questa classe' });
        return;
    }
    testForm.value = {
        title: '',
        date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
        evaluationType: 'Scritto',
        teacherNotes: '',
        parentNotes: '',
        grades: gradesStore.grades.students.map(s => ({
            student_id: s.student_id,
            full_name: s.full_name,
            grade_value: null,
            notes: ''
        }))
    };
    gradeInputRefs.value = [];
    showTestDialog.value = true;
};

const submitTest = async () => {
    if (!isAssignedClass.value || loading.value) return;
    loading.value = true;
    try {
        const payload = {
            class_id: selectedClassId.value,
            subject_id: selectedSubject.value,
            title: testForm.value.title,
            date: testForm.value.date,
            teacher_notes: testForm.value.teacherNotes,
            parent_notes: testForm.value.parentNotes,
            evaluation_type: testForm.value.evaluationType === 'Scritto' ? 'Written' : (testForm.value.evaluationType === 'Orale' ? 'Oral' : 'Practical'),
            grades: testForm.value.grades
                .filter(g => g.grade_value !== null && g.grade_value !== undefined && g.grade_value !== '')
                .map(g => ({
                    student_id: g.student_id,
                    grade_value: gradeToNumeric(g.grade_value),
                    notes: g.notes
                }))
        };

        await gradeService.createTestWithGrades(payload);
        showTestDialog.value = false;
        await refreshGrades();
        
        $q.notify({
            type: 'positive',
            message: `✓ Verifica "${testForm.value.title}" e ${payload.grades.length} voti salvati con successo!`,
            icon: 'check_circle',
            position: 'bottom-right'
        });
    } catch (err) {
        console.error(err);
        $q.notify({
            type: 'negative',
            message: 'Errore nel salvataggio della verifica'
        });
    } finally {
        loading.value = false;
    }
};

watch(() => schoolYearStore.selectedSchoolYear, async (newSY) => {
    await classesStore.fetchAssignedClasses(newSY);
    if (classesStore.classes.length > 0) {
        selectedClassId.value = classesStore.classes[0].id;
    } else {
        selectedClassId.value = null;
        selectedSubject.value = null;
    }
}, { immediate: true });

let classChangeReqId = 0;
watch(selectedClassId, async (newVal) => {
    const currentReq = ++classChangeReqId;
    selectedSubject.value = null;
    if (newVal) {
        await gradesStore.fetchClassSubjects(newVal);
        if (currentReq !== classChangeReqId) return;
        if (gradesStore.subjects && gradesStore.subjects.length > 0) {
            selectedSubject.value = gradesStore.subjects[0].subject_id;
        } else {
            selectedSubject.value = null;
        }
        await refreshGrades();
    }
});

const refreshGrades = async () => {
    if (!selectedClassId.value || !isAssignedClass.value) return;
    await gradesStore.fetchGrades(selectedClassId.value, selectedSubject.value);
    if (viewMode.value === 'history') {
        await fetchTests();
    }
};

const classTests = ref([]);
const loadingTests = ref(false);
const showEditTestDialog = ref(false);
const editTestForm = ref({
    id: '',
    title: '',
    date: '',
    evaluationType: '',
    teacherNotes: '',
    parentNotes: '',
    grades: []
});

const fetchTests = async () => {
    if (!selectedClassId.value || !selectedSubject.value || !isAssignedClass.value) return;
    loadingTests.value = true;
    try {
        const response = await gradeService.getClassTests(selectedClassId.value, selectedSubject.value);
        classTests.value = response.data || [];
    } catch (err) {
        console.error("Error fetching class tests:", err);
    } finally {
        loadingTests.value = false;
    }
};

watch(viewMode, (newVal) => {
    if (newVal === 'history') {
        fetchTests();
    }
});

const openEditTestDialog = (test) => {
    if (!isAssignedClass.value) return;
    const studentGrades = gradesStore.grades?.students || [];
    const gradesList = studentGrades.map(s => {
        const grade = s.grades?.find(g => g.test_id === test.id || g.testId === test.id);
        return {
            student_id: s.student_id,
            full_name: s.full_name,
            grade_id: grade ? grade.id : null,
            grade_value: grade ? formatGrade(grade.grade_value) : null,
            notes: grade ? grade.description : ''
        };
    });

    editTestForm.value = {
        id: test.id,
        title: test.title,
        date: test.date ? test.date.split('T')[0] : '',
        evaluationType: test.evaluation_type === 'Written' ? 'Scritto' : (test.evaluation_type === 'Oral' ? 'Orale' : 'Pratico'),
        teacherNotes: test.teacher_notes || '',
        parentNotes: test.parent_notes || '',
        grades: gradesList
    };
    editGradeInputRefs.value = [];
    showEditTestDialog.value = true;
};

const submitEditTest = async () => {
    if (!isAssignedClass.value) return;
    loading.value = true;
    try {
        const payload = {
            title: editTestForm.value.title,
            date: editTestForm.value.date,
            teacher_notes: editTestForm.value.teacherNotes,
            parent_notes: editTestForm.value.parentNotes,
            evaluation_type: editTestForm.value.evaluationType === 'Scritto' ? 'Written' : (editTestForm.value.evaluationType === 'Orale' ? 'Oral' : 'Practical'),
            grades: editTestForm.value.grades
                .filter(g => g.grade_value !== null && g.grade_value !== undefined && g.grade_value !== '')
                .map(g => ({
                    student_id: g.student_id,
                    grade_id: g.grade_id ?? null,
                    grade_value: gradeToNumeric(g.grade_value),
                    notes: g.notes
                }))
        };

        await gradesStore.updateClassTest(editTestForm.value.id, payload);
        $q.notify({ type: 'positive', message: 'Verifica modificata con successo!' });
        showEditTestDialog.value = false;
        await refreshGrades();
    } catch (err) {
        console.error(err);
        $q.notify({ type: 'negative', message: 'Errore durante la modifica della verifica' });
    } finally {
        loading.value = false;
    }
};

const gradeInputRefs = ref([]);
const editGradeInputRefs = ref([]);

const setGradeInputRef = (el, idx) => {
    if (el) gradeInputRefs.value[idx] = el;
};

const setEditGradeInputRef = (el, idx) => {
    if (el) editGradeInputRefs.value[idx] = el;
};

const focusNextStudent = (idx) => {
    if (gradeInputRefs.value && gradeInputRefs.value[idx + 1]) {
        const next = gradeInputRefs.value[idx + 1];
        if (next.focus) next.focus();
    }
};

const focusNextEditStudent = (idx) => {
    if (editGradeInputRefs.value && editGradeInputRefs.value[idx + 1]) {
        const next = editGradeInputRefs.value[idx + 1];
        if (next.focus) next.focus();
    }
};

const markAllAbsent = (formObj) => {
    if (!formObj || !formObj.grades) return;
    formObj.grades.forEach(g => {
        g.grade_value = null;
        g.notes = 'Assente';
    });
    $q.notify({ type: 'info', message: 'Tutti gli alunni segnati come assenti', timeout: 1500 });
};

const deleteTestConfirm = async (testOrId) => {
    if (!isAssignedClass.value) return;
    let test = typeof testOrId === 'object' ? testOrId : classTests.value.find(t => t.id === testOrId);
    const testTitle = test ? (test.title || 'Verifica') : 'questa verifica';
    const totalGrades = test && test.grade_count ? test.grade_count : (test && test.grades ? test.grades.length : 'tutti i');

    $q.dialog({
        title: 'Conferma Eliminazione Verifica',
        message: `Sei sicuro di voler eliminare la verifica "${testTitle}"? Verranno eliminati permanentemente ${totalGrades} voti collegati. L'operazione non è reversibile.`,
        cancel: { label: 'Annulla', flat: true },
        ok: { label: 'Elimina', color: 'negative' },
        persistent: true
    }).onOk(async () => {
        const id = test ? test.id : testOrId;
        try {
            await gradesStore.deleteClassTest(id);
            $q.notify({ type: 'positive', message: 'Verifica eliminata con successo!' });
            await refreshGrades();
            if (viewMode.value === 'history') await fetchTests();
        } catch (err) {
            console.error(err);
            $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione della verifica' });
        }
    });
};

const formatDate = (dateStr) => {
    if (!dateStr) return '';
    try {
        const parts = dateStr.split('T')[0].split('-');
        if (parts.length === 3) {
            return `${parts[2].padStart(2, '0')}/${parts[1].padStart(2, '0')}/${parts[0]}`;
        }
        const d = new Date(dateStr);
        const day = String(d.getDate()).padStart(2, '0');
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const year = d.getFullYear();
        return `${day}/${month}/${year}`;
    } catch (e) {
        return dateStr;
    }
};



const printReport = () => {
    showTestDialog.value = false;
    showEditTestDialog.value = false;
    showImportDialog.value = false;
    setTimeout(() => {
        window.print();
    }, 100);
};

const processImport = async () => {
    if (!isAssignedClass.value) return;
    if (!importFile.value) {
        $q.notify({ type: 'warning', message: 'Seleziona un file da importare' });
        return;
    }
    $q.loading.show();
    try {
        const formData = new FormData();
        formData.append('file', importFile.value);
        if (selectedClassId.value) formData.append('class_id', selectedClassId.value);
        if (selectedSubject.value) formData.append('subject_id', selectedSubject.value);

        await gradeService.bulkImport(formData);
        showImportDialog.value = false;
        importFile.value = null;
        $q.notify({ type: 'positive', message: 'Voti importati con successo' });
        refreshGrades();
    } catch (err) {
        console.error('Error importing grades:', err);
        $q.notify({ type: 'negative', message: 'Errore durante l\'importazione dei voti' });
    } finally {
        $q.loading.hide();
    }
};
</script>
