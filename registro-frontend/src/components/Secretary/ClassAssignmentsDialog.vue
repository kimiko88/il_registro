<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(1100px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
      <q-card-section class="bg-gradient-primary text-white row items-center q-pa-md shrink-0">
        <div class="row items-center">
          <q-avatar color="white-20" text-color="white" icon="menu_book" class="q-mr-sm" size="36px" />
          <div>
            <div class="text-h6 text-weight-bold">{{ t('secretaryClasses.assignmentsTitle') }} - {{ t('secretaryClasses.classLabel') }} {{ targetClass?.name }}{{ targetClass?.section }}</div>
            <div class="text-subtitle2 opacity-80">{{ targetClass?.academic_year }}</div>
          </div>
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-pa-md col overflow-y-auto">
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-7">
            <q-table
              :title="t('secretaryClasses.didacticProgramming')"
              :rows="assignments"
              :columns="assignmentsColumns"
              row-key="id"
              flat
              class="bg-transparent border-slate-100 rounded-xl"
            >
              <template #header-cell="props">
                <q-th :props="props" class="text-slate-500 font-bold">
                  {{ props.col.label }}
                </q-th>
              </template>

              <template #body-cell-actions="props">
                <q-td :props="props" auto-width>
                  <q-btn flat round dense color="negative" icon="delete" :aria-label="t('common.delete') || 'Rimuovi assegnazione materia'" @click="removeAssignment(props.row)" />
                </q-td>
              </template>
            </q-table>
          </div>

          <div class="col-12 col-md-5">
            <q-card flat class="rounded-xl bg-slate-50 q-pa-md border-slate-200">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">{{ t('secretaryClasses.assignSubject') }}</div>
              <q-form @submit="addAssignment" class="q-gutter-y-md">
                <q-select
                  v-model="assignForm.subject_id"
                  :options="subjectOptions"
                  :label="t('secretaryClasses.subject') + ' *'"
                  outlined dense
                  emit-value map-options
                  :rules="[val => !!val || t('secretaryClasses.selectSubject')]"
                >
                  <template #no-option>
                    <q-item>
                      <q-item-section class="text-grey">{{ t('secretaryClasses.noSubjectsFound') }}</q-item-section>
                    </q-item>
                    <q-item clickable @click="openCreateSubject">
                      <q-item-section class="text-primary text-weight-bold">{{ t('secretaryClasses.addNewSubject') }}</q-item-section>
                    </q-item>
                  </template>
                </q-select>

                <q-select
                  v-model="assignForm.teacher_id"
                  :options="teacherOptions"
                  :label="t('secretaryClasses.teacher')"
                  outlined dense
                  emit-value map-options
                />

                <q-input
                  v-model.number="assignForm.hours_per_week"
                  :label="t('secretaryClasses.hoursPerWeek')"
                  type="number"
                  outlined dense
                  min="1"
                />

                <q-btn type="submit" :label="t('secretaryClasses.assignChair')" color="primary" class="full-width rounded-lg q-py-sm shadow-sm q-mt-md" no-caps />
              </q-form>
            </q-card>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>

  <!-- Quick Create Subject Dialog -->
  <q-dialog v-model="showSubjectDialog">
    <q-card style="min-width: 350px" class="rounded-xl shadow-24 bg-white">
      <q-card-section class="bg-gradient-primary text-white row items-center q-pa-md">
        <div class="text-h6 text-weight-bold">{{ t('secretaryClasses.newSubject') }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-card-section class="q-pa-lg">
        <q-input
          v-model="newSubjectName"
          :label="t('secretaryClasses.subjectName')"
          outlined autofocus
          :rules="[val => !!val || t('common.required')]"
        />
        <div class="row justify-end q-mt-md">
          <q-btn :label="t('secretaryClasses.createSubject')" color="primary" @click="createSubject" />
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import adminService from '@/services/adminService'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  targetClass: {
    type: Object,
    default: null
  },
  subjectOptions: {
    type: Array,
    default: () => []
  },
  teacherOptions: {
    type: Array,
    default: () => []
  },
  schoolId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const $q = useQuasar()

const assignments = ref([])
const showSubjectDialog = ref(false)
const newSubjectName = ref('')

const assignForm = reactive({
  subject_id: null,
  teacher_id: null,
  hours_per_week: 1
})

const assignmentsColumns = computed(() => [
  { name: 'subject', label: t('secretaryClasses.subject'), field: 'subject_name', align: 'left', sortable: true },
  { name: 'teacher', label: t('secretaryClasses.teacher'), field: row => row.teacher_name || 'N/A', align: 'left' },
  { name: 'hours', label: t('secretaryClasses.hours'), field: 'hours_per_week', align: 'center' },
  { name: 'actions', label: t('common.actions'), align: 'right' }
])

watch(() => props.modelValue, async (open) => {
  if (open && props.targetClass) {
    await fetchAssignments(props.targetClass.id)
  }
})

const fetchAssignments = async (classId) => {
  try {
    const res = await adminService.getClassSubjects(classId)
    assignments.value = res.data || []
  } catch {
    $q.notify({ type: 'negative', message: t('secretaryClasses.loadAssignmentsError') })
  }
}

const addAssignment = async () => {
  if (!props.targetClass) return
  try {
    await adminService.assignSubjectToClass(props.targetClass.id, assignForm)
    $q.notify({ type: 'positive', message: t('secretaryClasses.subjectAssigned') })
    await fetchAssignments(props.targetClass.id)
  } catch {
    $q.notify({ type: 'negative', message: t('secretaryClasses.assignError') })
  }
}

const removeAssignment = async (row) => {
  try {
    await adminService.removeSubjectFromClass(props.targetClass.id, row.id)
    $q.notify({ type: 'positive', message: t('secretaryClasses.subjectRemoved') })
    await fetchAssignments(props.targetClass.id)
  } catch {
    $q.notify({ type: 'negative', message: t('secretaryClasses.removeError') })
  }
}

const openCreateSubject = () => {
  newSubjectName.value = ''
  showSubjectDialog.value = true
}

const createSubject = async () => {
  if (!newSubjectName.value) return
  try {
    await adminService.createSubject({ name: newSubjectName.value, school_id: props.schoolId })
    $q.notify({ type: 'positive', message: t('secretaryClasses.subjectCreated') })
    showSubjectDialog.value = false
    emit('subjects-updated')
  } catch {
    $q.notify({ type: 'negative', message: t('secretaryClasses.createSubjectError') })
  }
}
</script>
