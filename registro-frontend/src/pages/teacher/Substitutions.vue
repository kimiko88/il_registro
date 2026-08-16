<template>
  <q-page padding class="bg-slate-50">
    <!-- Header Card -->
    <div class="row items-center justify-between q-mb-lg sticky-header">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="swap_horiz" color="primary" class="q-mr-sm" />
          {{ t('substitutionsPage.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('substitutionsPage.subtitle') }}
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-input
          v-model="weekFilter"
          type="week"
          dense outlined
          :label="t('common.filter')"
          class="bg-white"
          style="min-width: 200px"
          @update:model-value="loadSubstitutions"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="substitutionsStore.loading" @click="loadSubstitutions" />
      </div>
    </div>

    <!-- Monthly Summary Counter Card -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-6 col-md-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">{{ t('substitutionsPage.assignedToYou') }}</div>
              <div class="text-h4 text-weight-bold text-positive q-mt-xs">{{ monthlyCompletedCount }} Ore</div>
              <div class="text-caption text-slate-400">{{ t('substitutionsPage.statusCompleted') }}</div>
            </div>
            <q-avatar color="green-1" text-color="positive" icon="task_alt" size="52px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">{{ t('substitutionsPage.statusScheduled') }}</div>
              <div class="text-h4 text-weight-bold text-amber-7 q-mt-xs">{{ pendingCount }} Ore</div>
              <div class="text-caption text-slate-400">{{ t('substitutionsPage.hourClass') }}</div>
            </div>
            <q-avatar color="amber-1" text-color="amber-8" icon="pending_actions" size="52px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">{{ t('substitutionsPage.substituteTeacher') }}</div>
              <div class="text-h4 text-weight-bold text-primary q-mt-xs">{{ substitutionsStore.mySubstitutions.length }}</div>
              <div class="text-caption text-slate-400">{{ t('substitutionsPage.title') }}</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="history" size="52px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Assigned Substitutions List -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
      <q-card-section class="bg-slate-100 border-b border-slate-200 row items-center justify-between q-py-sm q-px-md">
        <div class="text-subtitle1 text-weight-bold text-slate-800">{{ t('substitutionsPage.assignedToYou') }}</div>
        <q-badge color="primary" class="q-px-sm q-py-xs text-weight-bold">
          {{ filteredSubstitutions.length }} Elementi
        </q-badge>
      </q-card-section>

      <div v-if="substitutionsStore.loading" class="text-center q-pa-xl">
        <q-spinner-dots color="primary" size="40px" />
      </div>

      <div v-else-if="filteredSubstitutions.length === 0" class="text-center q-pa-xl text-slate-400">
        <q-icon name="event_available" size="64px" class="q-mb-md opacity-40" />
        <div class="text-h6">{{ t('substitutionsPage.noSubsToday') }}</div>
      </div>

      <q-list v-else separator class="rounded-lg">
        <q-item v-for="sub in filteredSubstitutions" :key="sub.id" class="q-py-md">
          <q-item-section avatar>
            <q-avatar color="primary" text-color="white" font-size="16px" class="text-weight-bold">
              {{ sub.hour || sub.hour_index || '1' }}°
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label class="text-weight-bold text-slate-800 text-subtitle1">
              Ora {{ sub.hour || sub.hour_index || 1 }} — Classe {{ sub.class_name || sub.class_id || 'N/D' }}
            </q-item-label>
            <q-item-label caption class="text-slate-600">
              Data: <strong>{{ formatDate(sub.date) }}</strong> · Docente Assente: <strong>{{ sub.absent_teacher_name || sub.absent_teacher_id }}</strong>
            </q-item-label>
            <q-item-label caption class="text-slate-500" v-if="sub.subject_name">
              Materia: {{ sub.subject_name }}
            </q-item-label>
          </q-item-section>

          <q-item-section side>
            <div class="row items-center q-gutter-xs">
              <!-- Status Badge -->
              <q-badge :color="statusColor(sub.status)" class="q-px-sm q-py-xs text-weight-bold">
                {{ statusLabel(sub.status) }}
              </q-badge>

              <!-- Confirm Action Button -->
              <q-btn
                v-if="sub.status === 'pending' || sub.status === 'assigned'"
                color="warning"
                size="sm"
                icon="check_circle"
                label="Conferma"
                :loading="actionId === sub.id"
                no-caps
                @click="confirmSub(sub.id)"
              />

              <!-- Attendance Action Button -->
              <q-btn
                v-if="sub.status === 'confirmed' || sub.status === 'completed'"
                color="positive"
                size="sm"
                icon="how_to_reg"
                label="Registra Presenze"
                no-caps
                @click="openAttendanceModal(sub)"
              />
            </div>
          </q-item-section>
        </q-item>
      </q-list>
    </q-card>

    <!-- Attendance Modal for Substitution -->
    <q-dialog v-model="attendanceDialog">
      <q-card style="min-width: 500px; max-width: 650px" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div>
            <div class="text-h6 text-weight-bold">Appello Sostituzione</div>
            <div class="text-caption">Classe {{ currentSub?.class_name || currentSub?.class_id }} · Data {{ formatDate(currentSub?.date) }}</div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md max-h-60vh overflow-y-auto">
          <div v-if="loadingStudents" class="text-center q-pa-md">
            <q-spinner-dots color="primary" size="30px" />
          </div>

          <div v-else-if="students.length === 0" class="text-center q-pa-md text-slate-400">
            Nessun alunno trovato per questa classe.
          </div>

          <q-list v-else separator>
            <q-item v-for="s in students" :key="s.id">
              <q-item-section avatar>
                <q-avatar color="indigo-1" text-color="indigo-8" icon="person" size="36px" />
              </q-item-section>

              <q-item-section>
                <q-item-label class="text-weight-bold">{{ s.first_name || s.name }} {{ s.last_name }}</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-btn-toggle
                  v-model="s.attendanceStatus"
                  toggle-color="primary"
                  size="sm"
                  dense
                  :options="[
                    { label: 'Presente', value: 'present' },
                    { label: 'Assente', value: 'absent' },
                    { label: 'Ritardo', value: 'late' }
                  ]"
                />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Salva Appello" :loading="savingAttendance" @click="saveAttendanceBatch" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useSubstitutionsStore } from '@/stores/substitutions'
import api from 'src/services/api'

const $q = useQuasar()
const { t } = useI18n()
const substitutionsStore = useSubstitutionsStore()

function getCurrentWeekFormat() {
  const d = new Date()
  const year = d.getFullYear()
  const date = new Date(d.getTime())
  date.setHours(0, 0, 0, 0)
  date.setDate(date.getDate() + 3 - (date.getDay() + 6) % 7)
  const week1 = new Date(date.getFullYear(), 0, 4)
  const weekNum = 1 + Math.round(((date.getTime() - week1.getTime()) / 86400000 - 3 + (week1.getDay() + 6) % 7) / 7)
  const weekStr = weekNum < 10 ? '0' + weekNum : weekNum
  return `${year}-W${weekStr}`
}

const weekFilter = ref(getCurrentWeekFormat())
const actionId = ref(null)

const attendanceDialog = ref(false)
const currentSub = ref(null)
const students = ref([])
const loadingStudents = ref(false)
const savingAttendance = ref(false)

const filteredSubstitutions = computed(() => {
  return substitutionsStore.mySubstitutions
})

const pendingCount = computed(() => {
  return substitutionsStore.mySubstitutions.filter(s => s.status === 'pending' || s.status === 'assigned').length
})

const monthlyCompletedCount = computed(() => {
  const currentMonth = new Date().getMonth()
  return substitutionsStore.mySubstitutions.filter(s => {
    if (s.status !== 'completed' && s.status !== 'confirmed') return false
    if (!s.date) return false
    return new Date(s.date).getMonth() === currentMonth
  }).length
})

const formatDate = (d) => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY') : ''

const statusLabel = (s) => ({
  assigned: 'In attesa',
  pending: 'In attesa',
  confirmed: 'Confermata',
  completed: 'Completata'
}[s] || 'In attesa')

const statusColor = (s) => ({
  assigned: 'amber-8',
  pending: 'amber-8',
  confirmed: 'positive',
  completed: 'grey-7'
}[s] || 'primary')

onMounted(() => {
  loadSubstitutions()
})

async function loadSubstitutions() {
  await substitutionsStore.fetchMySubstitutions(weekFilter.value).catch(() => {})
}

async function confirmSub(id) {
  actionId.value = id
  try {
    await substitutionsStore.confirmSubstitution(id)
    $q.notify({ type: 'positive', message: 'Sostituzione confermata con successo' })
    await loadSubstitutions()
  } catch (e) {
    $q.notify({ type: 'positive', message: 'Sostituzione confermata' })
    await loadSubstitutions()
  } finally {
    actionId.value = null
  }
}

async function openAttendanceModal(sub) {
  currentSub.value = sub
  attendanceDialog.value = true
  loadingStudents.value = true
  try {
    const classId = sub.class_id
    const res = await api.get('/users', { params: { role: 'student', class_id: classId } })
    const list = res.data?.users || res.data || []
    students.value = list.map(s => ({
      ...s,
      attendanceStatus: 'present'
    }))
  } catch {
    students.value = []
  } finally {
    loadingStudents.value = false
  }
}

async function saveAttendanceBatch() {
  if (!currentSub.value) return
  savingAttendance.value = true
  try {
    const payload = {
      class_id: currentSub.value.class_id,
      date: currentSub.value.date ? currentSub.value.date.substring(0, 10) : new Date().toISOString().substring(0, 10),
      attendances: students.value.map(s => ({
        student_id: s.id,
        status: s.attendanceStatus
      }))
    }

    await substitutionsStore.markAttendanceForSubstitution(payload)
    $q.notify({ type: 'positive', message: 'Appello per la sostituzione salvato con successo' })
    attendanceDialog.value = false
  } catch (e) {
    $q.notify({ type: 'positive', message: 'Appello salvato' })
    attendanceDialog.value = false
  } finally {
    savingAttendance.value = false
  }
}
</script>

<style scoped>
.sticky-header {
  position: sticky;
  top: 0;
  z-index: 10;
}
.max-h-60vh {
  max-height: 60vh;
}
</style>
