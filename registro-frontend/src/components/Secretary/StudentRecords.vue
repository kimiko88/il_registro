<template>
  <q-dialog v-model="show" full-width>
    <q-card class="column no-wrap rounded-xl overflow-hidden glass-card" style="min-height: 80vh">
      <q-card-section class="bg-gradient-primary text-white row items-center q-pa-lg">
        <div>
          <div class="text-h5 text-weight-bold text-outfit">{{ t('studentsPage.studentFile') || 'Scheda Studente' }}</div>
          <div class="text-subtitle2 opacity-80">{{ student?.first_name }} {{ student?.last_name }}</div>
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-tabs
        v-model="tab"
        dense
        class="text-grey-7 bg-slate-50"
        active-color="primary"
        indicator-color="primary"
        align="justify"
        narrow-indicator
      >
        <q-tab name="grades" icon="grade" :label="t('roleDashboards.tabGrades') || 'Voti'" />
        <q-tab name="attendance" icon="event" :label="t('roleDashboards.tabAttendance') || 'Assenze'" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="tab" animated class="col scroll">
        <!-- Grades Panel -->
        <q-tab-panel name="grades" class="q-pa-lg">
          <div v-if="loadingGrades" class="row justify-center q-pa-xl">
            <q-spinner-dots color="primary" size="40px" />
          </div>
          
          <div v-else-if="grades.length === 0" class="column items-center justify-center q-pa-xl text-grey-6">
            <q-icon name="history_edu" size="64px" class="q-mb-md" />
            <div class="text-h6">{{ t('gradesPage.noGrades') || 'Nessun voto registrato' }}</div>
          </div>

          <div v-else>
            <q-table
              :rows="grades"
              :columns="gradeColumns"
              row-key="id"
              flat
              bordered
              class="rounded-lg"
              :pagination="{ rowsPerPage: 10 }"
            >
              <template v-slot:body-cell-grade="props">
                <q-td :props="props">
                  <q-chip :color="getGradeColor(props.value)" text-color="white" dense class="text-weight-bold">
                    {{ props.value }}
                  </q-chip>
                </q-td>
              </template>
            </q-table>
          </div>
        </q-tab-panel>

        <!-- Attendance Panel -->
        <q-tab-panel name="attendance" class="q-pa-lg">
           <div v-if="loadingAttendance" class="row justify-center q-pa-xl">
            <q-spinner-dots color="primary" size="40px" />
          </div>
          
          <div v-else-if="attendance.length === 0" class="column items-center justify-center q-pa-xl text-grey-6">
            <q-icon name="event_available" size="64px" class="q-mb-md" />
            <div class="text-h6">{{ t('attendance.noAbsences') || 'Nessuna assenza registrata' }}</div>
          </div>

          <div v-else>
            <q-table
              :rows="attendance"
              :columns="attendanceColumns"
              row-key="id"
              flat
              bordered
              class="rounded-lg"
              :pagination="{ rowsPerPage: 10 }"
            >
              <template v-slot:body-cell-status="props">
                <q-td :props="props">
                  <q-chip :color="getAttendanceColor(props.value)" text-color="white" dense>
                    {{ props.value }}
                  </q-chip>
                </q-td>
              </template>
              <template v-slot:body-cell-justified="props">
                <q-td :props="props">
                  <q-icon 
                    :name="props.value ? 'check_circle' : 'warning'" 
                    :color="props.value ? 'positive' : 'warning'" 
                    size="24px"
                  />
                </q-td>
              </template>
            </q-table>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { gradeService } from '@/services/gradeService'
import { attendanceService } from '@/services/attendanceService'
import { date } from 'quasar'

const { t } = useI18n()

const props = defineProps({
  modelValue: Boolean,
  student: Object
})

const emit = defineEmits(['update:modelValue'])

const show = ref(false)
const tab = ref('grades')
const loadingGrades = ref(false)
const loadingAttendance = ref(false)
const grades = ref([])
const attendance = ref([])

const gradeColumns = computed(() => [
  { name: 'subject', label: t('agendaPage.subject') || 'Materia', field: row => row.subject_id, align: 'left', sortable: true },
  { name: 'grade', label: t('classRegister.tableHeaderGrade') || 'Voto', field: 'grade_value', align: 'center', sortable: true },
  { name: 'date', label: t('classRegister.dateLabel') || 'Data', field: row => date.formatDate(row.date, 'DD/MM/YYYY'), align: 'left', sortable: true },
  { name: 'type', label: t('classRegister.tableHeaderGradeType') || 'Tipo', field: 'grade_type', align: 'left' },
  { name: 'comment', label: t('common.comment') || 'Commento', field: 'description', align: 'left' }
])

const attendanceColumns = computed(() => [
  { name: 'date', label: t('classRegister.dateLabel') || 'Data', field: row => date.formatDate(row.date, 'DD/MM/YYYY'), align: 'left', sortable: true },
  { name: 'status', label: t('substitutionsPage.status') || 'Stato', field: 'status', align: 'center' },
  { name: 'justified', label: t('attendance.justified') || 'Giustificata', field: 'is_justified', align: 'center' },
  { name: 'reason', label: t('attendance.reason') || 'Motivazione', field: 'justification_reason', align: 'left' }
])

watch(() => props.modelValue, (val) => {
  show.value = val
  if (val && props.student) {
    fetchGrades()
    fetchAttendance()
  }
})

watch(show, (val) => {
  emit('update:modelValue', val)
})

const fetchGrades = async () => {
  if (!props.student?.id) return
  loadingGrades.value = true
  try {
    const res = await gradeService.getStudentGrades(props.student.id)
    grades.value = res.data || []
  } catch (e) {
    console.error('Error fetching grades:', e)
  } finally {
    loadingGrades.value = false
  }
}

const fetchAttendance = async () => {
  if (!props.student?.id) return
  loadingAttendance.value = true
  try {
    const res = await attendanceService.getChildAttendance(props.student.id)
    attendance.value = res.data || []
  } catch (e) {
    console.error('Error fetching attendance:', e)
  } finally {
    loadingAttendance.value = false
  }
}

const getGradeColor = (val) => {
  const n = parseFloat(val)
  if (n >= 6) return 'positive'
  if (n >= 5) return 'warning'
  return 'negative'
}

const getAttendanceColor = (status) => {
  if (!status) return 'grey'
  switch (String(status).toLowerCase()) {
    case 'present': return 'positive';
    case 'absent': return 'negative';
    case 'late': return 'warning';
    case 'left early': return 'info';
    default: return 'grey';
  }
}
</script>
