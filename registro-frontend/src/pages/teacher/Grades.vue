<template>
  <q-page class="q-pa-md" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-grey-1 text-dark'">
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
          :options="availableSubjectOptions"
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
                    {label: 'Griglia', value: 'matrix'},
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
      <q-card-section class="q-pt-sm q-pb-sm" :class="$q.dark.isActive ? 'bg-dark' : 'bg-grey-1'" v-if="viewMode === 'table'">
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
    <q-banner v-if="selectedClassId && !isAssignedClass" class="rounded-xl border q-mb-md shadow-soft" :class="$q.dark.isActive ? 'bg-amber-10 text-amber-1 border-amber-8' : 'bg-amber-1 text-amber-10 border-amber-300'">
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
                     <q-card :class="$q.dark.isActive ? 'bg-dark border border-grey-8' : 'bg-white'">
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

        <!-- Matrix View -->
        <div v-show="viewMode === 'matrix'">
          <GradeMatrixGrid
            v-if="isAssignedClass && selectedClassId && selectedSubject"
            :students-list="matrixStudentsList"
            :subject-id="String(selectedSubject)"
            :class-id="String(selectedClassId)"
            @saved="refreshGrades"
          />
          <div v-else-if="!isAssignedClass" class="text-center q-pa-xl text-grey-7">
            L'inserimento voti in griglia è riservato ai docenti titolari.
          </div>
          <div v-else class="text-center q-pa-xl text-grey-7">
            Seleziona una classe e una materia per visualizzare la griglia valutazioni.
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
                <q-card-section v-else-if="classTests.length === 0" class="text-center q-pa-xl" :class="$q.dark.isActive ? 'text-grey-4' : 'text-slate-600'">
                    <q-icon name="quiz" size="64px" color="grey-5" class="q-mb-md" />
                    <div class="text-h6 text-weight-bold">{{ t('gradesPage.noTestsFound') }}</div>
                    <div class="text-caption q-mb-md" :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-7'">{{ t('gradesPage.noTestsFoundDesc') }}</div>
                    <q-btn icon="add" :label="t('gradesPage.createFirstTest')" color="primary" unelevated no-caps @click="openTestDialog" />
                </q-card-section>
                <q-list separator v-else>
                    <q-item v-for="test in classTests" :key="test.id" class="q-py-md">
                        <q-item-section>
                            <q-item-label class="text-weight-bold text-subtitle1">{{ test.title }}</q-item-label>
                            <q-item-label caption class="row items-center q-gutter-x-sm">
                                <span>{{ t('gradesPage.date') }}: {{ formatDate(test.date) }}</span>
                                <span>|</span>
                                <span>{{ t('gradesPage.gradeType') }}: {{ test.evaluation_type === 'Written' ? 'Scritto' : (test.evaluation_type === 'Oral' ? 'Orale' : 'Pratico') }}</span>
                            </q-item-label>
                            <q-item-label class="text-caption q-mt-xs" :class="$q.dark.isActive ? 'text-grey-3' : 'text-grey-8'" v-if="test.teacher_notes">
                                <strong>Note Docente:</strong> {{ test.teacher_notes }}
                            </q-item-label>
                            <q-item-label class="text-caption" :class="$q.dark.isActive ? 'text-grey-3' : 'text-grey-8'" v-if="test.parent_notes">
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
            <div v-else class="text-center q-pa-xl" :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-7'">La cronologia verifiche è riservata ai docenti titolari.</div>
        </div>

    </div>
    <div v-else class="text-center q-pa-xl column items-center">
        <q-icon name="school" size="96px" color="primary" class="q-mb-md opacity-80" />
        <div class="text-h5 text-weight-bold q-mb-xs">{{ t('gradesPage.selectClassPrompt') }}</div>
        <div class="text-subtitle2 q-mb-lg" :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-7'" style="max-width: 480px;">
          {{ t('gradesPage.selectClassSub') }}
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
    <!-- Class Test Bulk Dialog (Create / Edit) -->
    <ClassTestBulkDialog
      v-model="showTestDialog"
      :is-edit="isEditTest"
      :test-data="testForm"
      :class-id="selectedClassId"
      :subject-id="selectedSubject"
      :overlapping-tests-count="overlappingTestsCount"
      :grade-options="gradeOptions"
      @saved="refreshGrades"
    />

  </q-page>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useClassesStore } from '@/stores/classes';
import { useGradesStore } from '@/stores/grades';
import { useAuthStore } from '@/stores/auth';
import GradeEntry from '@/components/Teacher/GradeEntry.vue';
import GradeStatistics from '@/components/Teacher/GradeStatistics.vue';
import GradeMatrixGrid from '@/components/Teacher/GradeMatrixGrid.vue';
import ClassTestBulkDialog from '@/components/Teacher/ClassTestBulkDialog.vue';
import { gradeService } from '@/services/gradeService';
import { useQuasar, date } from 'quasar';
import SkeletonTable from '@/components/Common/SkeletonTable.vue';
import { useUndoToast } from '@/composables/useUndoToast';
import { useSchoolYearStore } from '@/stores/schoolYear';
import { ITALIAN_GRADE_OPTIONS, formatGrade } from '@/utils/gradeUtils';

const $q = useQuasar();
const { t } = useI18n();
useUndoToast();
const classesStore = useClassesStore();
const gradesStore = useGradesStore();
const authStore = useAuthStore();
const schoolYearStore = useSchoolYearStore();

const selectedClassId = ref(null);
const selectedSubject = ref(null);

const isCivicaSubject = (s) => {
  const name = (s.subject_name || s.name || '').toLowerCase();
  return name.includes('civica') || name.includes('educazione civica') || name.includes('ed. civica');
};

const isAssignedToCurrentTeacher = (s, user) => {
  if (!user) return false;
  const currentUserId = String(user.id || '');
  const teacherId = user.teacher_id ? String(user.teacher_id) : '';
  const sTeacherId = s.teacher_id ? String(s.teacher_id) : '';
  const sTeacherUserId = s.teacher_user_id ? String(s.teacher_user_id) : '';

  if (sTeacherId && (sTeacherId === currentUserId || (teacherId && sTeacherId === teacherId))) {
    return true;
  }
  if (sTeacherUserId && (sTeacherUserId === currentUserId || (teacherId && sTeacherUserId === teacherId))) {
    return true;
  }
  if (s.teacher_name && user.last_name) {
    const tName = s.teacher_name.toLowerCase();
    const uLast = user.last_name.toLowerCase();
    const uFirst = (user.first_name || '').toLowerCase();
    if (tName.includes(uLast) && (!uFirst || tName.includes(uFirst))) {
      return true;
    }
  }
  return false;
};

const availableSubjectOptions = computed(() => {
  const allSubjects = gradesStore.subjects || [];
  const user = authStore.user;
  if (!user || ['admin', 'superadmin', 'secretary'].includes(user.role)) {
    return allSubjects;
  }

  // Mostra ESCLUSIVAMENTE le materie assegnate dalla segreteria al docente loggato + Educazione Civica
  return allSubjects.filter(s => {
    return isAssignedToCurrentTeacher(s, user) || isCivicaSubject(s);
  });
});
const viewMode = ref('table');
const filterDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'));
const gradeType = ref('Orale');
const showRubric = ref(false);

const gradeOptions = ITALIAN_GRADE_OPTIONS;

const isAssignedClass = computed(() => {
  if (!selectedClassId.value) return false
  const cls = classesStore.classes.find(c => String(c.id) === String(selectedClassId.value))
  if (!cls) return false
  // If the backend provides an is_owner / is_assigned field, use it.
  // Absence of the field (undefined) is treated as owned (backward-compatible).
  return cls.is_owner !== false
});

const showImportDialog = ref(false);
const importFile = ref(null);

const showTestDialog = ref(false);
const isEditTest = ref(false);
const testForm = ref({
    title: '',
    date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
    evaluationType: 'Scritto',
    teacherNotes: '',
    parentNotes: '',
    grades: []
});

const matrixStudentsList = computed(() => {
    return (gradesStore.grades?.students || []).map(s => ({
        id: s.student_id,
        first_name: s.first_name,
        last_name: s.last_name
    }));
});

const overlappingTestsCount = computed(() => {
    if (!testForm.value?.date || !classTests.value) return 0;
    // Use split('T')[0] for reliable comparison with ISO timestamps (e.g. '2026-08-23T00:00:00+02:00')
    return classTests.value.filter(t => t.date && t.date.split('T')[0] === testForm.value.date).length;
});

const openTestDialog = () => {
    if (!isAssignedClass.value) {
        $q.notify({ type: 'warning', message: t('gradesPage.noPermissionNotOwner') || 'Non hai i permessi per inserire verifiche in una classe non tua' });
        return;
    }
    if (!gradesStore.grades || !gradesStore.grades.students || gradesStore.grades.students.length === 0) {
        $q.notify({ type: 'warning', message: t('gradesPage.noStudentsInClass') || 'Nessun alunno caricato per questa classe' });
        return;
    }
    isEditTest.value = false;
    testForm.value = {
        title: '',
        date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
        evaluationType: 'Scritto',
        teacherNotes: '',
        parentNotes: '',
        grades: (gradesStore.grades?.students || []).map(s => ({
            student_id: s.student_id,
            full_name: s.full_name,
            grade_value: null,
            notes: ''
        }))
    };
    showTestDialog.value = true;
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
        if (availableSubjectOptions.value && availableSubjectOptions.value.length > 0) {
            selectedSubject.value = availableSubjectOptions.value[0].subject_id;
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

    isEditTest.value = true;
    testForm.value = {
        id: test.id,
        title: test.title,
        date: test.date ? test.date.split('T')[0] : '',
        evaluationType: test.evaluation_type === 'Written' ? 'Scritto' : (test.evaluation_type === 'Oral' ? 'Orale' : 'Pratico'),
        teacherNotes: test.teacher_notes || '',
        parentNotes: test.parent_notes || '',
        grades: gradesList
    };
    showTestDialog.value = true;
};

const deleteTestConfirm = async (testOrId) => {
    if (!isAssignedClass.value) return;
    let test = typeof testOrId === 'object' ? testOrId : classTests.value.find(t => t.id === testOrId);
    const testTitle = test ? (test.title || 'Verifica') : 'questa verifica';
    const totalGrades = test && test.grade_count ? test.grade_count : (test && test.grades ? test.grades.length : 'tutti i');

    $q.dialog({
        title: t('gradesPage.deleteTestConfirmTitle') || 'Conferma Eliminazione Verifica',
        message: t('gradesPage.deleteTestConfirmMsg', { title: testTitle, count: totalGrades })
          || `Sei sicuro di voler eliminare la verifica "${testTitle}"? Verranno eliminati permanentemente ${totalGrades} voti collegati. L'operazione non è reversibile.`,
        cancel: { label: t('common.cancel') || 'Annulla', flat: true },
        ok: { label: t('common.delete') || 'Elimina', color: 'negative' },
        persistent: true
    }).onOk(async () => {
        const id = test ? test.id : testOrId;
        try {
            await gradesStore.deleteClassTest(id);
            $q.notify({ type: 'positive', message: t('gradesPage.deleteTestSuccess') || 'Verifica eliminata con successo!' });
            await refreshGrades();
            if (viewMode.value === 'history') await fetchTests();
        } catch (err) {
            console.error(err);
            $q.notify({ type: 'negative', message: t('gradesPage.deleteTestError') || 'Errore durante l\'eliminazione della verifica' });
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
