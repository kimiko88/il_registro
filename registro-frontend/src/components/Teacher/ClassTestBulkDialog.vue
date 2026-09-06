<template>
  <q-dialog
    :model-value="modelValue"
    persistent
    max-width="96vw"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <q-card style="width: 1200px; max-width: 96vw; max-height: 92vh; overflow-y: auto;">
      <q-card-section class="bg-primary text-white row items-center">
        <div class="text-h6 text-weight-bold">
          {{ isEdit ? t('gradesPage.editTestTitle') : t('gradesPage.createTestTitle') }}
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-pa-md">
        <q-form ref="formRef" @submit.prevent="handleSubmit">
          <div class="row q-col-gutter-lg">
            <!-- Test Details (Left Column) -->
            <div class="col-12 col-md-4">
              <div class="text-subtitle1 q-mb-md text-weight-bold text-primary">
                {{ t('gradesPage.testDetails') }}
              </div>

              <q-input
                v-model="localForm.title"
                :label="(t('gradesPage.testTitle') || 'Titolo Verifica') + ' *'"
                outlined
                dense
                :rules="[val => (!!val && val.trim().length > 0) || (t('common.requiredField') || 'Campo obbligatorio')]"
                class="q-mb-sm"
              />

              <q-input
                v-model="localForm.date"
                type="date"
                :label="(t('gradesPage.testDate') || 'Data Verifica') + ' *'"
                outlined
                dense
                :rules="[val => !!val || (t('common.requiredField') || 'Campo obbligatorio')]"
                class="q-mb-sm"
              />

              <!-- Overlapping Test Alert -->
              <q-banner
                v-if="overlappingTestsCount >= 2 && !isEdit"
                rounded
                dense
                class="bg-amber-1 text-amber-9 border border-amber-3 q-mb-sm"
              >
                <template v-slot:avatar>
                  <q-icon name="warning" color="amber-9" />
                </template>
                {{ t('gradesPage.overlappingWarning', { count: overlappingTestsCount, date: localForm.date }) }}
              </q-banner>

              <q-select
                v-model="localForm.evaluationType"
                :options="evaluationTypeOptions"
                :label="(t('gradesPage.evalType') || 'Tipo Valutazione') + ' *'"
                outlined
                dense
                :rules="[val => !!val || (t('common.requiredField') || 'Campo obbligatorio')]"
                class="q-mb-md"
              />

              <q-input
                v-model="localForm.teacherNotes"
                type="textarea"
                :label="t('gradesPage.teacherNotes')"
                outlined
                dense
                rows="3"
                class="q-mb-sm"
              />

              <q-input
                v-model="localForm.parentNotes"
                type="textarea"
                :label="t('gradesPage.parentNotes')"
                outlined
                dense
                rows="3"
                class="q-mb-sm"
              />
            </div>

            <!-- Student Grades (Right Column) -->
            <div class="col-12 col-md-8">
              <div class="text-subtitle1 q-mb-xs text-weight-bold text-primary row items-center justify-between">
                <div>
                  {{ t('gradesPage.studentGrades') }} ({{ filledGradesCount }}/{{ localForm.grades.length }} {{ t('gradesPage.insertedCount') }})
                </div>
                <div class="row items-center q-gutter-x-xs">
                  <q-btn
                    icon="block"
                    size="sm"
                    outline
                    color="warning"
                    :label="t('gradesPage.markAllAbsent')"
                    @click="markAllAbsent"
                  />
                </div>
              </div>

              <q-linear-progress
                :value="localForm.grades.length ? filledGradesCount / localForm.grades.length : 0"
                color="primary"
                class="q-mb-sm"
              />

              <q-scroll-area
                style="height: 420px;"
                tabindex="0"
                :aria-label="t('gradesPage.studentGrades')"
                class="rounded-borders q-pa-sm border-default"
                :class="$q.dark.isActive ? 'bg-grey-9 text-white' : 'bg-grey-2 text-dark'"
              >
                <q-list separator>
                  <q-item
                    v-for="(student, idx) in localForm.grades"
                    :key="student.student_id"
                    class="q-py-sm"
                  >
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
                          :label="t('gradesPage.cols.value') || 'Voto'"
                          outlined
                          dense
                          style="width: 110px"
                          :bg-color="getGradeColor(student.grade_value)"
                          placeholder="-"
                          :rules="[
                            val => val === null || val === undefined || val === '' ||
                              (!isNaN(gradeToNumeric(val)) && gradeToNumeric(val) >= 1 && gradeToNumeric(val) <= 10) ||
                              (t('gradesPage.gradeRuleError') || '1-10')
                          ]"
                          lazy-rules
                          hide-bottom-space
                          :ref="el => setGradeInputRef(el, idx)"
                          @keydown.enter.prevent="focusNextStudent(idx)"
                        />
                        <q-input
                          v-model="student.notes"
                          :label="t('gradesPage.notes') || 'Note personali'"
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
            <q-btn
              flat
              :label="t('common.cancel') || 'Annulla'"
              v-close-popup
              color="grey-7"
            />
            <q-btn
              type="submit"
              :label="isEdit ? (t('gradesPage.saveEditTest') || 'Salva Modifiche Verifica') : (t('gradesPage.saveTestAndGrades') || 'Salva Verifica e Voti')"
              color="primary"
              :loading="submitting || loading"
            />
          </q-card-actions>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar, date } from 'quasar';
import { gradeService } from '@/services/gradeService';
import { useGradesStore } from '@/stores/grades';
import { ITALIAN_GRADE_OPTIONS, gradeToNumeric, getGradeColor } from '@/utils/gradeUtils';

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  isEdit: {
    type: Boolean,
    default: false
  },
  testData: {
    type: Object,
    default: () => null
  },
  classId: {
    type: [String, Number],
    default: null
  },
  subjectId: {
    type: [String, Number],
    default: null
  },
  overlappingTestsCount: {
    type: Number,
    default: 0
  },
  gradeOptions: {
    type: Array,
    default: () => ITALIAN_GRADE_OPTIONS
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue', 'saved']);

const $q = useQuasar();
const { t } = useI18n();
const gradesStore = useGradesStore();

const formRef = ref(null);
const submitting = ref(false);
const gradeInputRefs = ref([]);

const evaluationTypeOptions = ['Scritto', 'Orale', 'Pratico'];

const localForm = ref({
  id: '',
  title: '',
  date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
  evaluationType: 'Scritto',
  teacherNotes: '',
  parentNotes: '',
  grades: []
});

const filledGradesCount = computed(() => {
  if (!localForm.value?.grades) return 0;
  return localForm.value.grades.filter(
    g => g.grade_value !== null && g.grade_value !== undefined && g.grade_value !== ''
  ).length;
});

const initForm = () => {
  if (props.isEdit && props.testData) {
    localForm.value = {
      id: props.testData.id || '',
      title: props.testData.title || '',
      date: props.testData.date ? props.testData.date.split('T')[0] : '',
      evaluationType: props.testData.evaluationType || 'Scritto',
      teacherNotes: props.testData.teacherNotes || '',
      parentNotes: props.testData.parentNotes || '',
      grades: (props.testData.grades || []).map(g => ({
        student_id: g.student_id,
        full_name: g.full_name,
        grade_id: g.grade_id ?? null,
        grade_value: g.grade_value !== undefined ? g.grade_value : null,
        notes: g.notes || ''
      }))
    };
  } else if (props.testData) {
    localForm.value = {
      id: '',
      title: props.testData.title || '',
      date: props.testData.date || date.formatDate(Date.now(), 'YYYY-MM-DD'),
      evaluationType: props.testData.evaluationType || 'Scritto',
      teacherNotes: props.testData.teacherNotes || '',
      parentNotes: props.testData.parentNotes || '',
      grades: (props.testData.grades || []).map(s => ({
        student_id: s.student_id,
        full_name: s.full_name,
        grade_id: null,
        grade_value: s.grade_value !== undefined ? s.grade_value : null,
        notes: s.notes || ''
      }))
    };
  } else {
    localForm.value = {
      id: '',
      title: '',
      date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
      evaluationType: 'Scritto',
      teacherNotes: '',
      parentNotes: '',
      grades: []
    };
  }
  gradeInputRefs.value = [];
};

watch(
  () => props.modelValue,
  (isOpen) => {
    if (isOpen) {
      initForm();
    }
  },
  { immediate: true }
);

watch(
  () => props.testData,
  () => {
    if (props.modelValue) {
      initForm();
    }
  },
  { deep: true }
);

const setGradeInputRef = (el, idx) => {
  if (el) gradeInputRefs.value[idx] = el;
};

const focusNextStudent = (idx) => {
  if (gradeInputRefs.value && gradeInputRefs.value[idx + 1]) {
    const next = gradeInputRefs.value[idx + 1];
    if (next?.focus) next.focus();
  }
};

const markAllAbsent = () => {
  if (!localForm.value?.grades) return;
  localForm.value.grades.forEach(g => {
    g.grade_value = null;
    g.notes = 'Assente';
  });
  $q.notify({
    type: 'info',
    message: t('gradesPage.allAbsentMarked') || 'Tutti gli alunni segnati come assenti',
    timeout: 1500
  });
};

const handleSubmit = async () => {
  if (submitting.value) return;

  // Validate form with Quasar form validator if available
  if (formRef.value && typeof formRef.value.validate === 'function') {
    const valid = await formRef.value.validate();
    if (!valid) return;
  }

  // Range check on filled grades (1 - 10)
  for (const g of localForm.value.grades) {
    if (g.grade_value !== null && g.grade_value !== undefined && g.grade_value !== '') {
      const numericVal = gradeToNumeric(g.grade_value);
      if (isNaN(numericVal) || numericVal < 1 || numericVal > 10) {
        $q.notify({
          type: 'warning',
          message: t('gradesPage.gradeRuleError') || 'Inserisci un voto tra 1 e 10'
        });
        return;
      }
    }
  }

  submitting.value = true;
  try {
    const rawGrades = localForm.value.grades
      .filter(g => g.grade_value !== null && g.grade_value !== undefined && g.grade_value !== '')
      .map(g => ({
        student_id: g.student_id,
        ...(props.isEdit && g.grade_id ? { grade_id: g.grade_id } : {}),
        grade_value: gradeToNumeric(g.grade_value),
        notes: g.notes || ''
      }));

    const evalTypeApi = localForm.value.evaluationType === 'Scritto'
      ? 'Written'
      : (localForm.value.evaluationType === 'Orale' ? 'Oral' : 'Practical');

    if (props.isEdit) {
      const payload = {
        title: localForm.value.title,
        date: localForm.value.date,
        teacher_notes: localForm.value.teacherNotes,
        parent_notes: localForm.value.parentNotes,
        evaluation_type: evalTypeApi,
        grades: rawGrades
      };
      await gradesStore.updateClassTest(localForm.value.id, payload);
      $q.notify({
        type: 'positive',
        message: t('gradesPage.testEditSuccess') || 'Verifica modificata con successo!'
      });
    } else {
      const payload = {
        class_id: props.classId,
        subject_id: props.subjectId,
        title: localForm.value.title,
        date: localForm.value.date,
        teacher_notes: localForm.value.teacherNotes,
        parent_notes: localForm.value.parentNotes,
        evaluation_type: evalTypeApi,
        grades: rawGrades
      };
      await gradeService.createTestWithGrades(payload);
      $q.notify({
        type: 'positive',
        message: t('gradesPage.testSavedSuccess', { title: localForm.value.title, count: rawGrades.length })
          || `✓ Verifica "${localForm.value.title}" e ${rawGrades.length} voti salvati con successo!`,
        icon: 'check_circle',
        position: 'bottom-right'
      });
    }

    emit('update:modelValue', false);
    emit('saved');
  } catch (err) {
    console.error('Error saving class test:', err);
    $q.notify({
      type: 'negative',
      message: props.isEdit
        ? (t('gradesPage.testEditError') || 'Errore durante la modifica della verifica')
        : (t('gradesPage.testSaveError') || 'Errore nel salvataggio della verifica')
    });
  } finally {
    submitting.value = false;
  }
};
</script>
