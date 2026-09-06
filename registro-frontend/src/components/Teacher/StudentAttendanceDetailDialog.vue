<template>
  <q-dialog :model-value="modelValue" maximized-if-mobile @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="min-width: 340px; max-width: 520px; width: 100%">
      <q-bar class="bg-indigo-8 text-white">
        <q-icon name="person" class="q-mr-sm" />
        <span class="text-subtitle2 text-weight-bold">
          {{ student?.last_name }} {{ student?.first_name }}
        </span>
        <q-space />
        <q-btn flat round dense icon="close" v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
      </q-bar>

      <q-card-section class="q-pa-md">
        <!-- Loading -->
        <div v-if="loading" class="row justify-center q-pa-lg">
          <q-spinner color="primary" size="40px" />
        </div>

        <template v-else>
          <!-- Personal Info -->
          <div class="text-caption text-weight-bold text-grey-6 q-mb-xs text-uppercase letter-spacing-wide">
            {{ t('studentDetail.personalInfo') }}
          </div>
          <q-list bordered separator rounded class="q-mb-md">
            <q-item dense>
              <q-item-section avatar><q-icon name="badge" color="indigo" /></q-item-section>
              <q-item-section>
                <q-item-label caption>{{ t('studentDetail.fullName') }}</q-item-label>
                <q-item-label>{{ studentInfo?.last_name || student?.last_name }} {{ studentInfo?.first_name || student?.first_name }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item dense v-if="studentInfo?.fiscal_code">
              <q-item-section avatar><q-icon name="fingerprint" color="indigo" /></q-item-section>
              <q-item-section>
                <q-item-label caption>{{ t('studentDetail.fiscalCode') }}</q-item-label>
                <q-item-label class="text-mono">{{ studentInfo.fiscal_code }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item dense v-if="studentInfo?.class_name">
              <q-item-section avatar><q-icon name="class" color="indigo" /></q-item-section>
              <q-item-section>
                <q-item-label caption>{{ t('studentDetail.class') }}</q-item-label>
                <q-item-label>{{ studentInfo.class_name }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item dense v-if="studentInfo?.email">
              <q-item-section avatar><q-icon name="email" color="indigo" /></q-item-section>
              <q-item-section>
                <q-item-label caption>{{ t('studentDetail.email') }}</q-item-label>
                <q-item-label>{{ studentInfo.email }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item dense v-if="studentInfo?.phone_number">
              <q-item-section avatar><q-icon name="phone" color="indigo" /></q-item-section>
              <q-item-section>
                <q-item-label caption>{{ t('studentDetail.phone') }}</q-item-label>
                <q-item-label>{{ studentInfo.phone_number }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item dense v-if="studentInfo?.date_of_birth">
              <q-item-section avatar><q-icon name="cake" color="indigo" /></q-item-section>
              <q-item-section>
                <q-item-label caption>{{ t('studentDetail.birthDate') }}</q-item-label>
                <q-item-label>{{ formatDate(studentInfo.date_of_birth) }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>

          <!-- Attendance Summary -->
          <div class="text-caption text-weight-bold text-grey-6 q-mb-xs text-uppercase">
            {{ t('studentDetail.attendanceSummary') }}
          </div>
          <div class="row q-col-gutter-sm q-mb-md">
            <div class="col-6">
              <q-card flat bordered class="text-center q-pa-sm">
                <div class="text-h5 text-negative text-weight-bold">{{ summary?.total_absences ?? '—' }}</div>
                <div class="text-caption text-grey-7">{{ t('studentDetail.totalAbsences') }}</div>
              </q-card>
            </div>
            <div class="col-6">
              <q-card flat bordered class="text-center q-pa-sm">
                <div class="text-h5 text-warning text-weight-bold">{{ summary?.total_lates ?? '—' }}</div>
                <div class="text-caption text-grey-7">{{ t('studentDetail.lates') }}</div>
              </q-card>
            </div>
            <div class="col-6">
              <q-card flat bordered class="text-center q-pa-sm">
                <div class="text-h5 text-purple text-weight-bold">{{ summary?.total_early_exits ?? '—' }}</div>
                <div class="text-caption text-grey-7">{{ t('studentDetail.earlyExits') }}</div>
              </q-card>
            </div>
            <div class="col-6">
              <q-card flat bordered class="text-center q-pa-sm">
                <div class="text-h5 text-positive text-weight-bold">{{ summary?.justified_count ?? '—' }}</div>
                <div class="text-caption text-grey-7">{{ t('studentDetail.justified') }}</div>
              </q-card>
            </div>
          </div>

          <!-- Absence Rate + Risk -->
          <div v-if="summary" class="q-mb-sm">
            <div class="row items-center justify-between q-mb-xs">
              <span class="text-caption text-grey-7">{{ t('studentDetail.absenceRate') }}</span>
              <span
                class="text-caption text-weight-bold"
                :class="summary.absence_rate > 25 ? 'text-negative' : summary.absence_rate > 10 ? 'text-warning' : 'text-positive'"
              >
                {{ summary.absence_rate?.toFixed(1) }}%
              </span>
            </div>
            <q-linear-progress
              :value="(summary.absence_rate || 0) / 100"
              :color="summary.absence_rate > 25 ? 'negative' : summary.absence_rate > 10 ? 'warning' : 'positive'"
              rounded
              size="8px"
              class="q-mb-xs"
            />
            <q-chip
              dense
              :color="summary.risk_level === 'high' ? 'negative' : summary.risk_level === 'medium' ? 'warning' : 'positive'"
              text-color="white"
              :icon="summary.risk_level === 'high' ? 'warning' : summary.risk_level === 'medium' ? 'info' : 'check_circle'"
            >
              {{ t('studentDetail.risk') }}: {{ getRiskLabel(summary.risk_level) }}
            </q-chip>
          </div>

          <!-- Today's attendance for this student -->
          <div class="text-caption text-weight-bold text-grey-6 q-mt-md q-mb-xs text-uppercase">
            {{ t('studentDetail.todayAttendanceByHour') }}
          </div>
          <div class="row q-gutter-xs">
            <q-badge
              v-for="h in 8"
              :key="h"
              :color="getHourBadgeColor(student, h)"
              :label="String(h) + 'ª'"
              class="text-weight-bold"
              style="font-size: 11px; padding: 4px 8px"
            >
              <q-tooltip>{{ getHourLabel(student, h) }}</q-tooltip>
            </q-badge>
          </div>
        </template>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/services/api'
import attendanceService from '@/services/attendanceService'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  student: {
    type: Object,
    default: null
  },
  allTodayAttendance: {
    type: Array,
    default: () => []
  }
})

defineEmits(['update:modelValue'])

const { t } = useI18n()

const loading = ref(false)
const studentInfo = ref(null)
const summary = ref(null)

const formatDate = (isoDate) => {
  if (!isoDate) return ''
  const parts = isoDate.split('T')[0].split('-')
  if (parts.length === 3) {
    return `${parts[2]}/${parts[1]}/${parts[0]}`
  }
  return isoDate
}

const getRiskLabel = (riskLevel) => {
  if (riskLevel === 'high') return t('studentDetail.riskHigh')
  if (riskLevel === 'medium') return t('studentDetail.riskMedium')
  return t('studentDetail.riskLow')
}

const getHourBadgeColor = (student, hour) => {
  if (!student) return 'grey-3'
  const rec = (props.allTodayAttendance || []).find(
    r => r.student_id === student.id && String(r.hour) === String(hour)
  )
  if (!rec) return 'grey-3'
  switch (rec.status) {
    case 'Present': return 'positive'
    case 'OutOfClass': return 'teal'
    case 'Absent': return 'negative'
    case 'Late': return 'warning'
    case 'LeftEarly': return 'purple'
    default: return 'grey-3'
  }
}

const getHourLabel = (student, hour) => {
  if (!student) return t('studentDetail.hourNotRegistered', { hour })
  const rec = (props.allTodayAttendance || []).find(
    r => r.student_id === student.id && String(r.hour) === String(hour)
  )
  if (!rec) return t('studentDetail.hourNotRegistered', { hour })

  const statusMap = {
    Present: t('classRegister.present'),
    OutOfClass: t('classRegister.outOfClass'),
    Absent: t('classRegister.absent'),
    Late: `${t('classRegister.late')}${rec.entry_time ? ` (${rec.entry_time})` : ''}`,
    LeftEarly: `${t('classRegister.earlyExit')}${rec.exit_time ? ` (${rec.exit_time})` : ''}`
  }

  const label = statusMap[rec.status] || rec.status
  return t('studentDetail.hourStatus', { hour, status: label })
}

const fetchStudentDetails = async () => {
  if (!props.student?.id) return
  loading.value = true
  studentInfo.value = null
  summary.value = null

  try {
    const [infoRes, summaryRes] = await Promise.allSettled([
      api.get(`/users/${props.student.id}`),
      attendanceService.getStudentSummary(props.student.id)
    ])

    if (infoRes.status === 'fulfilled') {
      studentInfo.value = infoRes.value.data
    }
    if (summaryRes.status === 'fulfilled') {
      summary.value = summaryRes.value
    }
  } catch {
    // Graceful fallback
  } finally {
    loading.value = false
  }
}

watch(
  () => [props.modelValue, props.student?.id],
  ([newVal, studentId]) => {
    if (newVal && studentId) {
      fetchStudentDetails()
    }
  },
  { immediate: true }
)
</script>
