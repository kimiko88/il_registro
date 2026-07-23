<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Sticky Header for Context -->
    <q-card class="sticky-header q-mb-md shadow-2 z-top">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6 q-mr-md">Gestione Voti</div>
        <q-select
          v-model="selectedClassId"
          :options="classesStore.classes"
          option-value="id"
          option-label="name"
          label="Classe"
          dense outlined
          options-dense
          class="q-mr-md"
          style="min-width: 150px"
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
                :options="[
                    {label: 'Registro', value: 'table'},
                    {label: 'Statistiche', value: 'stats'},
                    {label: 'Storico', value: 'history'}
                ]"
             />
             <q-btn icon="assignment" label="Nuova Verifica" color="primary" class="q-mr-sm" @click="openTestDialog" :disable="!selectedClassId || !selectedSubject" />
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
             <GradeStatistics :class-id="selectedClassId" :subject="selectedSubject" />
        </div>

        <!-- History / Verifiche View -->
        <div v-show="viewMode === 'history'">
            <q-card>
                <q-card-section class="row items-center justify-between">
                    <div class="text-h6 text-primary text-weight-bold">Elenco Verifiche</div>
                    <q-btn icon="refresh" flat round dense @click="fetchTests" />
                </q-card-section>
                <q-separator />
                <q-card-section v-if="loadingTests" class="text-center q-pa-xl">
                    <q-spinner color="primary" size="40px" />
                </q-card-section>
                <q-card-section v-else-if="classTests.length === 0" class="text-center text-grey q-pa-xl">
                    Nessuna verifica registrata per questa materia
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

    <!-- Create Class Test Dialog -->
    <q-dialog v-model="showTestDialog" persistent max-width="80vw">
      <q-card style="width: 1000px; max-width: 90vw;">
        <q-card-section class="bg-primary text-white row items-center">
          <div class="text-h6">Crea Nuova Verifica</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-form @submit.prevent="submitTest">
            <div class="row q-col-gutter-md">
              <!-- Test Details -->
              <div class="col-12 col-md-5">
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
              <div class="col-12 col-md-7">
                <div class="text-subtitle1 q-mb-md text-weight-bold text-primary row items-center justify-between">
                  <div>Voti Alunni</div>
                  <div class="text-caption text-grey-8">Inserisci i voti per ciascun alunno (lascia vuoto per assenti)</div>
                </div>
                
                <q-scroll-area style="height: 350px;" class="border-grey rounded-borders q-pa-sm bg-grey-2">
                  <q-list separator>
                    <q-item v-for="student in testForm.grades" :key="student.student_id" class="q-py-sm">
                      <q-item-section>
                        <q-item-label class="text-weight-bold">{{ student.full_name }}</q-item-label>
                      </q-item-section>
                      <q-item-section side style="width: 320px">
                        <div class="row items-center q-gutter-sm no-wrap">
                          <q-select
                            v-model="student.grade_value"
                            :options="gradeOptions"
                            label="Voto"
                            outlined
                            dense
                            style="width: 100px"
                            :bg-color="getGradeColor(student.grade_value)"
                            placeholder="-"
                          />
                          <q-input
                            v-model="student.notes"
                            label="Note personali"
                            outlined
                            dense
                            class="col"
                            placeholder="Note..."
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
    <q-dialog v-model="showEditTestDialog" persistent max-width="80vw">
      <q-card style="width: 1000px; max-width: 90vw;">
        <q-card-section class="bg-primary text-white row items-center">
          <div class="text-h6 text-weight-bold">Modifica Verifica in Blocco</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-form @submit.prevent="submitEditTest">
            <div class="row q-col-gutter-md">
              <!-- Test Details -->
              <div class="col-12 col-md-5 q-gutter-y-md">
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
              <div class="col-12 col-md-7">
                <div class="text-subtitle1 q-mb-md text-weight-bold text-primary row items-center justify-between">
                  <div>Voti Alunni</div>
                  <div class="text-caption text-grey-8">Modifica i voti per ciascun alunno (lascia vuoto per assenti/eliminare)</div>
                </div>
                
                <q-scroll-area style="height: 350px;" class="border-grey rounded-borders q-pa-sm bg-grey-2">
                  <q-list separator>
                    <q-item v-for="student in editTestForm.grades" :key="student.student_id" class="q-py-sm">
                      <q-item-section>
                        <q-item-label class="text-weight-bold">{{ student.full_name }}</q-item-label>
                      </q-item-section>
                      <q-item-section side style="width: 320px">
                        <div class="row items-center q-gutter-sm no-wrap">
                          <q-select
                            v-model="student.grade_value"
                            :options="gradeOptions"
                            label="Voto"
                            outlined
                            dense
                            style="width: 100px"
                            :bg-color="getGradeColor(student.grade_value)"
                            placeholder="-"
                          />
                          <q-input
                            v-model="student.notes"
                            label="Note personali"
                            outlined
                            dense
                            class="col"
                            placeholder="Note..."
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
import { ref, watch, onMounted } from 'vue';
import { useClassesStore } from 'src/stores/classes';
import { useGradesStore } from 'src/stores/grades';
import GradeEntry from 'src/components/Teacher/GradeEntry.vue';
import GradeStatistics from 'src/components/Teacher/GradeStatistics.vue';
import { gradeService } from 'src/services/gradeService';
import { useQuasar, date } from 'quasar';

const $q = useQuasar();
const classesStore = useClassesStore();
const gradesStore = useGradesStore();

const selectedClassId = ref(null);
const selectedSubject = ref(null); 
const viewMode = ref('table');
const filterDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'));
const gradeType = ref('Orale');
const showRubric = ref(false);
const offlineMode = ref(false);

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

const openTestDialog = () => {
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
    showTestDialog.value = true;
};

const submitTest = async () => {
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

        await gradesStore.createClassTest(payload);
        $q.notify({
            type: 'positive',
            message: 'Verifica e voti salvati con successo!'
        });
        showTestDialog.value = false;
        refreshGrades();
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

watch(selectedClassId, async (newVal) => {
    if (newVal) {
        await gradesStore.fetchClassSubjects(newVal);
        if (gradesStore.subjects.length > 0) {
            selectedSubject.value = gradesStore.subjects[0].subject_id;
        } else {
            selectedSubject.value = null;
        }
    }
    refreshGrades();
});

const refreshGrades = async () => {
    if (!selectedClassId.value) return;
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
    if (!selectedClassId.value || !selectedSubject.value) return;
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
    const studentGrades = gradesStore.grades?.students || [];
    const gradesList = studentGrades.map(s => {
        const grade = s.grades?.find(g => g.test_id === test.id);
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
    showEditTestDialog.value = true;
};

const submitEditTest = async () => {
    loading.value = true;
    try {
        const payload = {
            title: editTestForm.value.title,
            date: editTestForm.value.date,
            teacher_notes: editTestForm.value.teacherNotes,
            parent_notes: editTestForm.value.parentNotes,
            evaluation_type: editTestForm.value.evaluationType === 'Scritto' ? 'Written' : (editTestForm.value.evaluationType === 'Orale' ? 'Oral' : 'Practical'),
            grades: editTestForm.value.grades.map(g => ({
                student_id: g.student_id,
                grade_value: (g.grade_value !== null && g.grade_value !== undefined && g.grade_value !== '') ? gradeToNumeric(g.grade_value) : null,
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

const deleteTestConfirm = async (testId) => {
    $q.dialog({
        title: 'Elimina Verifica',
        message: 'Sei sicuro di voler eliminare questa verifica e TUTTI i voti ad essa collegati? L\'operazione non è reversibile.',
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            await gradesStore.deleteClassTest(testId);
            $q.notify({ type: 'positive', message: 'Verifica eliminata con successo!' });
            await refreshGrades();
        } catch (err) {
            console.error(err);
            $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione della verifica' });
        }
    });
};

const formatDate = (dateStr) => {
    if (!dateStr) return '';
    try {
        const d = new Date(dateStr);
        const day = String(d.getDate()).padStart(2, '0');
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const year = d.getFullYear();
        return `${day}/${month}/${year}`;
    } catch (e) {
        return dateStr;
    }
};

function formatGrade(val) {
    if (val === undefined || val === null || val === '-') return null;
    const num = Number(val);
    if (isNaN(num)) return val;
    if (num === -1) return 'A';
    
    const integerPart = Math.floor(num);
    const decimalPart = num - integerPart;
    
    if (Math.abs(decimalPart - 0.5) < 0.01) {
        return `${integerPart}½`;
    }
    if (Math.abs(decimalPart - 0.25) < 0.01) {
        return `${integerPart}+`;
    }
    if (Math.abs(decimalPart - 0.75) < 0.01) {
        return `${integerPart + 1}-`;
    }
    if (Math.abs(decimalPart) < 0.01) {
        return `${integerPart}`;
    }
    return String(num).replace('.', ',');
}

const printReport = () => {
    window.print();
};

const processImport = async () => {
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

onMounted(async () => {
    await classesStore.fetchAssignedClasses();
    if (classesStore.classes.length > 0) {
        selectedClassId.value = classesStore.classes[0].id;
    }
});
const gradeOptions = [
  'A',
  '10', '10-', '9½', '9+', '9', '9-', '8½', '8+', '8', '8-', '7½', '7+', '7', '7-', '6½', '6+', '6', '6-', '5½', '5+', '5', '5-', '4½', '4+', '4', '4-', '3½', '3+', '3', '3-', '2½', '2+', '2', '2-', '1½', '1+', '1', '1-', '0½', '0'
];

function gradeToNumeric(gradeStr) {
    if (gradeStr === undefined || gradeStr === null || gradeStr === '') return null;
    const clean = String(gradeStr).trim().toUpperCase();
    if (clean === 'A') return -1;
    
    if (clean.includes('/')) {
        const parts = clean.split('/');
        if (parts.length === 2) {
            const n1 = parseFloat(parts[0]);
            const n2 = parseFloat(parts[1]);
            if (!isNaN(n1) && !isNaN(n2)) {
                return (n1 + n2) / 2;
            }
        }
    }
    
    if (clean.endsWith('1/2') || clean.endsWith('½')) {
        const base = parseFloat(clean.replace('1/2', '').replace('½', '').trim());
        if (!isNaN(base)) return base + 0.5;
    }
    
    if (clean.endsWith('+')) {
        const base = parseFloat(clean.slice(0, -1).trim());
        if (!isNaN(base)) return base + 0.25;
    }
    
    if (clean.endsWith('-')) {
        const base = parseFloat(clean.slice(0, -1).trim());
        if (!isNaN(base)) return base - 0.25;
    }
    
    const val = parseFloat(clean.replace(',', '.'));
    return isNaN(val) ? 0 : val;
}

const getGradeColor = (val) => {
    if (val === undefined || val === null || val === '') return '';
    const numeric = gradeToNumeric(val);
    if (numeric === -1) return 'bg-red-1';
    return numeric < 6 ? 'bg-red-1' : 'bg-green-1';
};
</script>

<style scoped>
.sticky-header {
    position: sticky;
    top: 50px; /* Adjust based on navbar height */
    z-index: 100;
}
.z-top { z-index: 1000; }
</style>
