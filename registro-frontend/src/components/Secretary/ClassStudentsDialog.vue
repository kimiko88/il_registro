<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(1100px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
      <q-card-section class="bg-cyan-8 text-white row items-center q-pa-md shrink-0">
        <div class="row items-center">
          <q-avatar color="white-20" text-color="white" icon="groups" class="q-mr-sm" size="36px" />
          <div>
            <div class="text-h6 text-weight-bold">
              {{ t('secretaryClasses.studentsTitle') }} - {{ t('secretaryClasses.classLabel') }} {{ targetClass?.name }}{{ targetClass?.section }}
            </div>
            <div class="text-subtitle2 text-white/90">
              {{ classStudents.length }} {{ t('secretaryClasses.enrolledStudents') }} · {{ t('secretaryClasses.academicYear') }} {{ targetClass?.academic_year }}
            </div>
          </div>
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-pa-md col overflow-y-auto">
        <div v-if="loadingStudents" class="text-center q-pa-xl">
          <q-spinner-dots color="cyan-8" size="40px" />
        </div>

        <div v-else class="row q-col-gutter-md">
          <!-- Left Column: Currently Enrolled Students -->
          <div class="col-12 col-md-7">
            <div class="row items-center justify-between q-mb-sm">
              <div class="text-subtitle1 text-weight-bold text-slate-800">
                {{ t('secretaryClasses.studentsInClass') }} ({{ classStudents.length }})
              </div>
              <q-input
                v-model="studentSearchFilter"
                :placeholder="t('secretaryClasses.filterStudent')"
                dense outlined
                class="bg-white rounded-lg min-width-200"
              >
                <template #append>
                  <q-icon name="search" size="xs" />
                </template>
              </q-input>
            </div>

            <div v-if="filteredClassStudents.length === 0" class="bg-slate-50 border border-slate-200 rounded-xl p-6 text-center text-slate-400">
              <q-icon name="person_off" size="48px" class="q-mb-xs opacity-40" />
              <div class="text-subtitle2">{{ t('secretaryClasses.noStudentsInClass') }}</div>
              <div class="text-caption">{{ t('secretaryClasses.useRightPanel') }}</div>
            </div>

            <q-scroll-area style="height: 440px;" class="rounded-xl border border-slate-200 bg-white">
              <q-list separator>
                <q-item v-for="student in filteredClassStudents" :key="student.id" class="q-py-sm items-center justify-between">
                  <q-item-section avatar>
                    <q-avatar color="cyan-1" text-color="cyan-9" icon="person" font-size="20px" />
                  </q-item-section>

                  <q-item-section class="col overflow-hidden q-pr-sm">
                    <q-item-label class="text-weight-bold text-slate-800 ellipsis">
                      {{ student.last_name || '' }} {{ student.first_name || student.name || '' }}
                    </q-item-label>
                    <q-item-label caption class="text-slate-500 ellipsis">
                      {{ student.email }} {{ student.enrollment_number ? '· Matr: ' + student.enrollment_number : '' }}
                    </q-item-label>
                  </q-item-section>

                  <q-item-section side class="shrink-0">
                    <q-btn flat round dense color="negative" icon="person_remove" :aria-label="t('secretaryClasses.removeFromClass') || 'Rimuovi studente dalla classe'" @click="removeStudentFromClass(student)">
                      <q-tooltip>{{ t('secretaryClasses.removeFromClass') }}</q-tooltip>
                    </q-btn>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-scroll-area>
          </div>

          <!-- Right Column: Add / Transfer Student -->
          <div class="col-12 col-md-5">
            <q-card flat class="rounded-xl bg-slate-50 q-pa-md border border-slate-200">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">
                {{ t('secretaryClasses.assignOrTransfer') }}
              </div>

              <div class="q-gutter-y-md">
                <q-select
                  v-model="selectedStudentToAssign"
                  :options="assignableStudentOptions"
                  :label="t('secretaryClasses.selectStudent')"
                  outlined dense
                  emit-value map-options
                  filterable
                >
                  <template #no-option>
                    <q-item>
                      <q-item-section class="text-grey">{{ t('secretaryClasses.allStudentsAssigned') }}</q-item-section>
                    </q-item>
                  </template>
                </q-select>

                <q-btn
                  color="cyan-8"
                  :label="t('secretaryClasses.assignToClass')"
                  icon="person_add"
                  class="full-width rounded-lg q-py-sm shadow-sm"
                  no-caps
                  :disabled="!selectedStudentToAssign"
                  @click="assignStudentToClass(selectedStudentToAssign)"
                />
              </div>

              <q-separator class="q-my-md" />

              <!-- Unassigned Students Quick Pick -->
              <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-xs">
                {{ t('secretaryClasses.unassignedStudents') }} ({{ unassignedStudents.length }})
              </div>
              <div v-if="unassignedStudents.length === 0" class="text-caption text-slate-400">
                {{ t('secretaryClasses.allStudentsHaveClass') }}
              </div>
              <q-list v-else separator class="rounded-lg border border-slate-200 bg-white overflow-y-auto overflow-x-hidden" style="max-height: 200px;">
                <q-item v-for="uSt in unassignedStudents" :key="uSt.id" class="q-py-xs items-center justify-between">
                  <q-item-section class="col overflow-hidden q-pr-xs">
                    <q-item-label class="text-caption text-weight-bold text-slate-800 ellipsis">
                      {{ uSt.last_name }} {{ uSt.first_name || uSt.name }}
                    </q-item-label>
                    <q-item-label caption class="text-slate-400 text-xs ellipsis">
                      {{ uSt.email }}
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side class="shrink-0">
                    <q-btn flat round size="sm" color="cyan-8" icon="person_add" :aria-label="t('secretaryClasses.assignToClass') || 'Assegna studente alla classe'" @click="assignStudentToClass(uSt.id)">
                      <q-tooltip>{{ t('secretaryClasses.assignToClass') }}</q-tooltip>
                    </q-btn>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import api from '@/services/api'
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
  schoolId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'refresh'])

const { t } = useI18n()
const $q = useQuasar()

const classStudents = ref([])
const allSchoolStudents = ref([])
const loadingStudents = ref(false)
const selectedStudentToAssign = ref(null)
const studentSearchFilter = ref('')

watch(() => props.modelValue, async (open) => {
  if (open && props.targetClass) {
    selectedStudentToAssign.value = null
    studentSearchFilter.value = ''
    await fetchStudentsForClass()
  }
})

const sortAlphabetically = (arr) => {
  return [...arr].sort((a, b) => {
    const lastA = (a.last_name || '').toLowerCase()
    const lastB = (b.last_name || '').toLowerCase()
    if (lastA !== lastB) return lastA.localeCompare(lastB, 'it')
    const firstA = (a.first_name || a.name || '').toLowerCase()
    const firstB = (b.first_name || b.name || '').toLowerCase()
    return firstA.localeCompare(firstB, 'it')
  })
}

const filteredClassStudents = computed(() => {
  let list = classStudents.value
  if (studentSearchFilter.value) {
    const q = studentSearchFilter.value.toLowerCase().trim()
    list = list.filter(s =>
      (s.first_name || s.name || '').toLowerCase().includes(q) ||
      (s.last_name || '').toLowerCase().includes(q) ||
      (s.email || '').toLowerCase().includes(q)
    )
  }
  return sortAlphabetically(list)
})

const unassignedStudents = computed(() => {
  const list = allSchoolStudents.value.filter(s => !s.class_id || s.class_id === '')
  return sortAlphabetically(list)
})

const assignableStudentOptions = computed(() => {
  const sorted = sortAlphabetically(allSchoolStudents.value.filter(s => s.class_id !== props.targetClass?.id))
  return sorted.map(s => {
    const classNameStr = s.class_name ? ` (${t('secretaryClasses.currentlyIn')} ${s.class_name})` : ` (${t('secretaryClasses.noClass')})`
    return {
      label: `${s.last_name || ''} ${s.first_name || s.name || ''}${classNameStr}`.trim(),
      value: s.id
    }
  })
})

const fetchStudentsForClass = async () => {
  if (!props.targetClass) return
  loadingStudents.value = true
  try {
    const [classRes, allRes] = await Promise.all([
      api.get('/users', { params: { role: 'student', class_id: props.targetClass.id, page_size: 500 } }),
      api.get('/users', { params: { role: 'student', school_id: props.schoolId, page_size: 500 } })
    ])
    classStudents.value = classRes.data?.users || classRes.data || []
    allSchoolStudents.value = allRes.data?.users || allRes.data || []
  } catch {
    $q.notify({ type: 'negative', message: t('secretaryClasses.loadStudentsError') })
  } finally {
    loadingStudents.value = false
  }
}

const assignStudentToClass = async (studentId) => {
  if (!studentId || !props.targetClass) return
  try {
    await adminService.updateUser(studentId, { class_id: props.targetClass.id })
    $q.notify({ type: 'positive', message: t('secretaryClasses.studentAssigned') })
    selectedStudentToAssign.value = null
    await fetchStudentsForClass()
    emit('refresh')
  } catch {
    $q.notify({ type: 'negative', message: t('secretaryClasses.assignError') })
  }
}

const removeStudentFromClass = (student) => {
  $q.dialog({
    title: t('secretaryClasses.removeTitle'),
    message: `${t('secretaryClasses.removeConfirm')} ${student.first_name || ''} ${student.last_name || ''} ${t('secretaryClasses.fromClass')} ${props.targetClass?.name || ''}${props.targetClass?.section || ''}?`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: t('common.remove') }
  }).onOk(async () => {
    try {
      await adminService.updateUser(student.id, { class_id: '' })
      $q.notify({ type: 'positive', message: t('secretaryClasses.studentRemoved') })
      await fetchStudentsForClass()
      emit('refresh')
    } catch {
      $q.notify({ type: 'negative', message: t('secretaryClasses.removeError') })
    }
  })
}
</script>
