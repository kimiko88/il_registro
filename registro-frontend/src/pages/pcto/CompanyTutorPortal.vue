<template>
  <q-layout view="lHh Lpr lFf" class="bg-grey-2">
    <q-header elevated class="bg-teal-8 text-white">
      <q-toolbar>
        <q-icon name="business" size="md" class="q-mr-sm" />
        <q-toolbar-title class="text-subtitle1 text-weight-bold">
          {{ t('pctoCompanyTutor.portalTitle') }} | {{ tutor?.company_name || t('pctoCompanyTutor.company') }}
        </q-toolbar-title>
        <div v-if="tutor" class="text-caption">
          Tutor: {{ tutor.tutor_first_name }} {{ tutor.tutor_last_name }}
        </div>
      </q-toolbar>
    </q-header>

    <q-page-container>
      <q-page class="q-pa-md" style="max-width: 900px; margin: 0 auto;">
        <!-- Loading or Error -->
        <div v-if="loadingSession" class="text-center q-pa-xl">
          <q-spinner-dots size="50px" color="teal" />
          <div class="text-caption text-grey-7 q-mt-md">{{ t('common.loading') }}</div>
        </div>

        <div v-else-if="authError" class="q-pa-md">
          <q-banner rounded class="bg-red-1 text-negative">
            <template #avatar><q-icon name="error" color="negative" /></template>
            <div class="text-weight-bold">{{ t('common.error') }}</div>
            <div class="text-caption">{{ authError }}</div>
          </q-banner>
        </div>

        <!-- Main Tutor Workspace -->
        <div v-else>
          <q-banner rounded class="bg-teal-1 text-teal-9 q-mb-md">
            <template #avatar><q-icon name="info" color="teal" /></template>
            {{ t('pctoCompanyTutor.portalSubtitle') }}
          </q-banner>

          <div class="text-h6 text-weight-bold text-blue-grey-9 q-mb-sm">
            <q-icon name="school" class="q-mr-xs" /> {{ t('pctoCompanyTutor.student') }}
          </div>

          <q-card v-for="student in students" :key="student.student_id" flat bordered class="q-mb-md bg-white">
            <q-card-section class="row items-center justify-between">
              <div>
                <div class="text-subtitle1 text-weight-bold text-teal-10">{{ student.student_name }}</div>
                <div class="text-caption text-grey-7">{{ t('common.description') }}: {{ student.project_title }}</div>
                <div class="text-caption text-grey-8 q-mt-xs">
                  {{ t('pctoCompanyTutor.hours') }}: <strong>{{ student.completed_hours }}</strong> / {{ student.total_hours }}
                </div>
              </div>
              <div class="column q-gutter-xs">
                <q-btn
                  color="teal-7"
                  icon="check"
                  :label="t('pctoCompanyTutor.validateHours')"
                  size="sm"
                  unelevated
                  no-caps
                  @click="openTimesheetDialog(student)"
                />
                <q-btn
                  :color="student.is_evaluated ? 'positive' : 'indigo-7'"
                  :icon="student.is_evaluated ? 'verified' : 'rate_review'"
                  :label="student.is_evaluated ? t('pctoCompanyTutor.evaluationSubmittedSuccess') : t('pctoCompanyTutor.submitEvaluation')"
                  size="sm"
                  :flat="student.is_evaluated"
                  :unelevated="!student.is_evaluated"
                  no-caps
                  :disable="student.is_evaluated"
                  @click="openEvalDialog(student)"
                />
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Timesheet Verification Dialog -->
        <q-dialog v-model="showTimesheetDialog" persistent>
          <q-card style="min-width: 400px">
            <q-card-section class="bg-teal-8 text-white">
              <div class="text-h6">{{ t('pctoCompanyTutor.hoursValidation') }}</div>
            </q-card-section>
            <q-card-section class="q-pt-md">
              <div class="text-subtitle2 q-mb-xs">{{ t('pctoCompanyTutor.student') }}: {{ activeStudent?.student_name }}</div>
              <q-input v-model="timesheetForm.activity_date" type="date" :label="t('pctoCompanyTutor.date') + ' *'" outlined dense class="q-mb-sm" />
              <q-input v-model.number="timesheetForm.hours_declared" type="number" :label="t('pctoCompanyTutor.hours') + ' *'" outlined dense class="q-mb-sm" />
              <q-input v-model.number="timesheetForm.hours_approved" type="number" :label="t('pctoCompanyTutor.validateHours') + ' *'" outlined dense class="q-mb-sm" />
              <q-input v-model="timesheetForm.tutor_notes" :label="t('pctoCompanyTutor.activities')" outlined dense type="textarea" rows="2" />
            </q-card-section>
            <q-card-actions align="right">
              <q-btn flat :label="t('common.cancel')" v-close-popup />
              <q-btn color="teal-8" :label="t('pctoCompanyTutor.validateHours')" unelevated :loading="submitting" @click="confirmTimesheet" />
            </q-card-actions>
          </q-card>
        </q-dialog>

        <!-- Competencies Evaluation Dialog -->
        <q-dialog v-model="showEvalDialog" persistent>
          <q-card style="min-width: 500px">
            <q-card-section class="bg-indigo-8 text-white">
              <div class="text-h6">{{ t('pctoCompanyTutor.rubricTitle') }}</div>
            </q-card-section>
            <q-card-section class="q-pt-md">
              <div class="text-subtitle2 q-mb-sm">{{ t('pctoCompanyTutor.student') }}: {{ activeStudent?.student_name }}</div>

              <div class="q-mb-md">
                <div class="text-caption text-weight-bold">{{ t('pctoCompanyTutor.punctuality') }}:</div>
                <q-rating v-model="evalForm.reliability_level" max="5" size="2em" color="teal" icon="star_border" icon-selected="star" />
              </div>

              <div class="q-mb-md">
                <div class="text-caption text-weight-bold">{{ t('pctoCompanyTutor.technicalSkills') }}:</div>
                <q-rating v-model="evalForm.technical_skills" max="5" size="2em" color="teal" icon="star_border" icon-selected="star" />
              </div>

              <div class="q-mb-md">
                <div class="text-caption text-weight-bold">{{ t('pctoCompanyTutor.teamwork') }}:</div>
                <q-rating v-model="evalForm.teamwork_skills" max="5" size="2em" color="teal" icon="star_border" icon-selected="star" />
              </div>

              <q-input
                v-model="evalForm.final_feedback"
                type="textarea"
                :label="t('pctoCompanyTutor.finalComments') + ' *'"
                outlined
                dense
                rows="3"
              />
            </q-card-section>
            <q-card-actions align="right">
              <q-btn flat :label="t('common.cancel')" v-close-popup />
              <q-btn color="indigo-8" :label="t('pctoCompanyTutor.submitEvaluation')" unelevated :loading="submitting" @click="confirmEvaluation" />
            </q-card-actions>
          </q-card>
        </q-dialog>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { pctoCompanyTutorService } from '@/services/pctoCompanyTutorService'

export default defineComponent({
  name: 'CompanyTutorPortal',
  setup() {
    const { t } = useI18n()
    const route = useRoute()
    const tutor = ref(null)
    const students = ref([])
    const loadingSession = ref(true)
    const submitting = ref(false)
    const authError = ref('')
    const token = ref('')

    const showTimesheetDialog = ref(false)
    const showEvalDialog = ref(false)
    const activeStudent = ref(null)

    const timesheetForm = ref({
      activity_date: new Date().toISOString().substring(0, 10),
      hours_declared: 8,
      hours_approved: 8,
      tutor_notes: ''
    })

    const evalForm = ref({
      reliability_level: 5,
      technical_skills: 4,
      teamwork_skills: 5,
      final_feedback: ''
    })

    const loadSession = async () => {
      token.value = route.query.token || localStorage.getItem('pcto_tutor_token') || ''
      if (!token.value) {
        authError.value = 'Token non valido / Invalid token'
        loadingSession.value = false
        return
      }

      try {
        const sessionRes = await pctoCompanyTutorService.getSession(token.value)
        tutor.value = sessionRes.tutor
        localStorage.setItem('pcto_tutor_token', token.value)

        const studentsRes = await pctoCompanyTutorService.getAssignedStudents(token.value)
        students.value = studentsRes.data || []
      } catch (e) {
        authError.value = e.response?.data?.error || t('common.error')
      } finally {
        loadingSession.value = false
      }
    }

    const openTimesheetDialog = (student) => {
      activeStudent.value = student
      timesheetForm.value = {
        activity_date: new Date().toISOString().substring(0, 10),
        hours_declared: 8,
        hours_approved: 8,
        tutor_notes: ''
      }
      showTimesheetDialog.value = true
    }

    const confirmTimesheet = async () => {
      if (!activeStudent.value) return
      submitting.value = true
      try {
        await pctoCompanyTutorService.verifyTimesheet(token.value, {
          project_id: activeStudent.value.project_id,
          student_id: activeStudent.value.student_id,
          activity_date: timesheetForm.value.activity_date,
          hours_declared: timesheetForm.value.hours_declared,
          hours_approved: timesheetForm.value.hours_approved,
          tutor_notes: timesheetForm.value.tutor_notes
        })
        showTimesheetDialog.value = false
        await loadSession()
      } catch (e) {
        console.error('Errore firma ore:', e)
      } finally {
        submitting.value = false
      }
    }

    const openEvalDialog = (student) => {
      activeStudent.value = student
      evalForm.value = {
        reliability_level: 5,
        technical_skills: 4,
        teamwork_skills: 5,
        final_feedback: ''
      }
      showEvalDialog.value = true
    }

    const confirmEvaluation = async () => {
      if (!activeStudent.value) return
      submitting.value = true
      try {
        await pctoCompanyTutorService.submitEvaluation(token.value, {
          project_id: activeStudent.value.project_id,
          student_id: activeStudent.value.student_id,
          reliability_level: evalForm.value.reliability_level,
          technical_skills: evalForm.value.technical_skills,
          teamwork_skills: evalForm.value.teamwork_skills,
          final_feedback: evalForm.value.final_feedback
        })
        showEvalDialog.value = false
        await loadSession()
      } catch (e) {
        console.error('Errore invio valutazione:', e)
      } finally {
        submitting.value = false
      }
    }

    onMounted(loadSession)

    return {
      t,
      tutor,
      students,
      loadingSession,
      submitting,
      authError,
      showTimesheetDialog,
      showEvalDialog,
      activeStudent,
      timesheetForm,
      evalForm,
      openTimesheetDialog,
      confirmTimesheet,
      openEvalDialog,
      confirmEvaluation
    }
  }
})
</script>
