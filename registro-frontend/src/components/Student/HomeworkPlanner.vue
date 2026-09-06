<template>
  <q-card class="glass-card shadow-soft rounded-xl border border-slate-100 overflow-hidden">
    <q-card-section class="q-pa-lg">
      <!-- Header -->
      <div class="row items-center justify-between q-mb-md">
        <div class="row items-center q-gutter-sm">
          <q-avatar color="indigo-1" text-color="indigo" icon="task_alt" size="md" />
          <div>
            <div class="text-h6 text-weight-bold text-slate-800">
              {{ t('student.planner.title') || 'Diario Digitale & To-Do Compiti' }}
            </div>
            <div class="text-caption text-slate-500">
              {{ t('student.planner.subtitle') || 'Organizza e spunta i compiti completati sincronizzati con il registro' }}
            </div>
          </div>
        </div>

        <div class="row items-center q-gutter-xs">
          <q-btn
            flat round dense
            icon="refresh"
            color="primary"
            :loading="loading"
            @click="fetchHomeworks"
          >
            <q-tooltip>Aggiorna compiti</q-tooltip>
          </q-btn>
        </div>
      </div>

      <!-- Progress Bar & Stats -->
      <div class="bg-indigo-50 border border-indigo-100 rounded-xl q-pa-md q-mb-md">
        <div class="row items-center justify-between q-mb-xs">
          <span class="text-caption font-semibold text-indigo-900">
            Avanzamento Studio: {{ completedCount }} di {{ totalTasks }} compiti completati
          </span>
          <span class="text-caption font-bold text-indigo-700">{{ completionPercentage }}%</span>
        </div>
        <q-linear-progress
          :value="completionPercentage / 100"
          color="indigo"
          track-color="indigo-200"
          class="rounded-sm"
          size="8px"
        />
        <div v-if="completionPercentage === 100 && totalTasks > 0" class="q-mt-sm text-caption text-positive font-bold row items-center">
          <q-icon name="celebration" size="16px" class="q-mr-xs" />
          Tutti i compiti completati! Ottimo lavoro! 🚀
        </div>
      </div>

      <!-- Filters Row -->
      <div class="row q-col-gutter-xs q-mb-md items-center">
        <div class="col-12 col-sm-auto row q-gutter-xs">
          <q-btn
            v-for="f in filterOptions"
            :key="f.key"
            size="sm"
            dense
            :unelevated="activeFilter === f.key"
            :outline="activeFilter !== f.key"
            :color="activeFilter === f.key ? 'indigo' : 'grey-7'"
            class="rounded-lg q-px-sm font-medium"
            @click="activeFilter = f.key"
          >
            {{ f.label }}
          </q-btn>
        </div>

        <div class="col-12 col-sm">
          <q-select
            v-model="selectedSubject"
            :options="subjectFilterOptions"
            dense outlined
            label="Filtra Materia"
            bg-color="white"
            clearable
            class="text-caption"
          />
        </div>
      </div>

      <!-- Tasks List -->
      <div v-if="loading && tasks.length === 0" class="text-center q-py-lg">
        <q-spinner-dots color="indigo" size="36px" />
      </div>

      <div v-else-if="filteredTasks.length === 0" class="text-center q-py-xl text-slate-400">
        <q-icon name="done_all" size="48px" color="positive" class="q-mb-sm opacity-60" />
        <div class="text-subtitle1 font-semibold text-slate-700">Nessun compito da completare!</div>
        <div class="text-caption">Tutti i compiti per il filtro selezionato sono stati spuntati.</div>
      </div>

      <q-list v-else separator class="rounded-xl border border-slate-200 bg-white overflow-hidden">
        <q-item
          v-for="task in filteredTasks"
          :key="task.id"
          tag="label"
          v-ripple
          class="q-py-md transition-bg"
          :class="task.is_completed ? 'bg-slate-50 opacity-75' : ''"
        >
          <q-item-section side top>
            <q-checkbox
              v-model="task.is_completed"
              color="positive"
              @update:model-value="val => toggleTaskCompletion(task, val)"
            />
          </q-item-section>

          <q-item-section>
            <div class="row items-center q-gutter-x-sm">
              <span
                class="text-subtitle2 font-semibold"
                :class="task.is_completed ? 'text-strike text-slate-400' : 'text-slate-800'"
              >
                {{ task.title || task.description }}
              </span>
              <q-badge
                v-if="task.subject_name"
                color="indigo-1"
                text-color="indigo-9"
                class="q-px-xs rounded-borders text-xs font-semibold"
              >
                {{ task.subject_name }}
              </q-badge>
            </div>

            <div v-if="task.description && task.title" class="text-caption text-slate-500 q-mt-xs">
              {{ task.description }}
            </div>

            <div class="row items-center q-gutter-x-sm text-caption q-mt-xs">
              <span :class="getDueDateClass(task.date || task.due_date)">
                <q-icon name="schedule" size="14px" class="q-mr-xs" />
                Scadenza: {{ formatDueDate(task.date || task.due_date) }}
              </span>
              <span v-if="task.teacher_name" class="text-slate-400">
                · Assegnato da: {{ task.teacher_name }}
              </span>
            </div>
          </q-item-section>

          <q-item-section side>
            <q-badge
              :color="task.is_completed ? 'positive' : getDueBadgeColor(task.date || task.due_date)"
              class="q-px-xs text-xs font-bold"
            >
              {{ task.is_completed ? 'Fatto ✓' : getDueBadgeLabel(task.date || task.due_date) }}
            </q-badge>
          </q-item-section>
        </q-item>
      </q-list>

      <!-- Personal Study Note Section -->
      <div class="q-mt-lg border-t pt-md">
        <div class="row items-center justify-between q-mb-sm">
          <span class="text-caption font-bold text-slate-700 row items-center">
            <q-icon name="edit_note" color="indigo" size="18px" class="q-mr-xs" />
            Promemoria Personale di Studio
          </span>
          <q-btn
            v-if="personalNote"
            flat dense size="xs"
            color="negative"
            label="Cancella"
            @click="clearPersonalNote"
          />
        </div>
        <q-input
          v-model="personalNote"
          type="textarea"
          rows="2"
          dense outlined
          placeholder="Aggiungi un promemoria (es. Ripassare formule di fisica per domani, portare dizionario)..."
          bg-color="white"
          @update:model-value="savePersonalNote"
        />
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import api from '@/services/api'

const { t } = useI18n()
const $q = useQuasar()

const loading = ref(false)
const tasks = ref([])
const activeFilter = ref('all')
const selectedSubject = ref(null)
const personalNote = ref(localStorage.getItem('student_study_note') || '')

const filterOptions = [
  { key: 'all', label: 'Tutti' },
  { key: 'today', label: 'Oggi' },
  { key: 'tomorrow', label: 'Domani' },
  { key: 'week', label: 'Questa Settimana' }
]

const totalTasks = computed(() => tasks.value.length)
const completedCount = computed(() => tasks.value.filter(t => t.is_completed).length)
const completionPercentage = computed(() => {
  if (totalTasks.value === 0) return 0
  return Math.round((completedCount.value / totalTasks.value) * 100)
})

const subjectFilterOptions = computed(() => {
  const subjects = new Set()
  tasks.value.forEach(t => {
    if (t.subject_name) subjects.add(t.subject_name)
  })
  return Array.from(subjects)
})

const filteredTasks = computed(() => {
  const now = new Date()
  const todayStr = now.toISOString().substring(0, 10)

  const tomorrow = new Date(now)
  tomorrow.setDate(now.getDate() + 1)
  const tomorrowStr = tomorrow.toISOString().substring(0, 10)

  const endOfWeek = new Date(now)
  endOfWeek.setDate(now.getDate() + (7 - now.getDay()))
  const endOfWeekStr = endOfWeek.toISOString().substring(0, 10)

  return tasks.value.filter(task => {
    // Subject filter
    if (selectedSubject.value && task.subject_name !== selectedSubject.value) {
      return false
    }

    const taskDate = (task.date || task.due_date || '').substring(0, 10)
    if (!taskDate) return true

    switch (activeFilter.value) {
      case 'today':
        return taskDate === todayStr
      case 'tomorrow':
        return taskDate === tomorrowStr
      case 'week':
        return taskDate >= todayStr && taskDate <= endOfWeekStr
      default:
        return true
    }
  })
})

function formatDueDate(dStr) {
  if (!dStr) return '-'
  const d = new Date(dStr)
  return d.toLocaleDateString([], { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function getDueDateClass(dStr) {
  if (!dStr) return 'text-slate-500'
  const today = new Date().toISOString().substring(0, 10)
  const taskDate = dStr.substring(0, 10)
  if (taskDate < today) return 'text-negative font-bold'
  if (taskDate === today) return 'text-amber-800 font-bold'
  return 'text-slate-600'
}

function getDueBadgeColor(dStr) {
  if (!dStr) return 'primary'
  const today = new Date().toISOString().substring(0, 10)
  const taskDate = dStr.substring(0, 10)
  if (taskDate < today) return 'negative'
  if (taskDate === today) return 'warning'
  return 'primary'
}

function getDueBadgeLabel(dStr) {
  if (!dStr) return 'In corso'
  const today = new Date().toISOString().substring(0, 10)
  const taskDate = dStr.substring(0, 10)
  if (taskDate < today) return 'Scaduto'
  if (taskDate === today) return 'Oggi'
  return 'In corso'
}

async function fetchHomeworks() {
  loading.value = true
  try {
    const res = await api.get('/agenda')
    const items = res.data || []
    tasks.value = items.map(item => ({
      ...item,
      is_completed: item.is_completed === true
    }))
  } catch (err) {
    console.error('Failed to fetch agenda homeworks', err)
  } finally {
    loading.value = false
  }
}

async function toggleTaskCompletion(task, completed) {
  try {
    if (completed) {
      await api.post(`/agenda/${task.id}/complete`)
      $q.notify({
        type: 'positive',
        message: '✓ Compito segnato come completato!',
        timeout: 1500
      })
    } else {
      await api.delete(`/agenda/${task.id}/complete`)
      $q.notify({
        type: 'info',
        message: 'Compito riaperto.',
        timeout: 1500
      })
    }
  } catch (err) {
    console.error('Failed to update task completion', err)
    // Rollback
    task.is_completed = !completed
    $q.notify({
      type: 'negative',
      message: 'Impossibile aggiornare lo stato del compito'
    })
  }
}

function savePersonalNote(val) {
  localStorage.setItem('student_study_note', val)
}

function clearPersonalNote() {
  personalNote.value = ''
  localStorage.removeItem('student_study_note')
}

onMounted(() => {
  fetchHomeworks()
})
</script>

<style scoped>
.font-medium {
  font-weight: 500;
}
.font-semibold {
  font-weight: 600;
}
.font-bold {
  font-weight: 700;
}
.text-xs {
  font-size: 0.72rem;
}
.transition-bg {
  transition: background-color 0.2s ease, opacity 0.2s ease;
}
</style>
