<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none">📅 {{ $t('agendaPage.studentTitle') }}</h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none">{{ $t('agendaPage.studentSubtitle') }}</p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn flat round icon="refresh" color="primary" :loading="agendaStore.loading" @click="loadAgenda" />
        <q-btn-toggle
          v-model="filterType"
          toggle-color="primary"
          :options="[
            { label: $t('agendaPage.all'), value: 'all' },
            { label: $t('agendaPage.homework'), value: 'compito' },
            { label: $t('agendaPage.test'), value: 'verifica' },
            { label: $t('agendaPage.oralTest'), value: 'interrogazione' },
            { label: $t('agendaPage.other'), value: 'altro' }
          ]"
        />
      </div>
    </div>

    <div v-if="agendaStore.loading" class="row q-col-gutter-md q-my-md">
      <div v-for="n in 3" :key="n" class="col-12 col-md-6 col-lg-4">
        <q-card flat bordered class="q-pa-md rounded-borders">
          <div class="row items-center justify-between q-mb-md">
            <q-skeleton type="QChip" width="80px" />
            <q-skeleton type="text" width="120px" />
          </div>
          <q-skeleton type="text" class="text-h6 q-mb-sm" />
          <q-skeleton type="text" width="60%" class="q-mb-xs" />
          <q-skeleton type="text" width="40%" />
        </q-card>
      </div>
    </div>

    <div v-else-if="filteredHomeworks.length === 0" class="text-center q-my-xl q-pa-lg bg-grey-1 rounded-lg bordered">
      <q-icon name="event_busy" size="48px" color="grey-6" />
      <div class="text-h6 text-grey-7 q-mt-sm">{{ $t('agendaPage.noEvents') }}</div>
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="item in filteredHomeworks" :key="item.id" class="col-12 col-md-6 col-lg-4">
        <q-card flat bordered :class="{ 'bg-grey-2': isItemCompleted(item.id) }">
          <q-card-section class="row items-center justify-between">
            <q-chip
              :color="getBadgeColor(item.type)"
              text-color="white"
              size="sm"
            >
              {{ (item.type || 'compito').toUpperCase() }}
            </q-chip>
            <div class="text-caption text-grey">{{ $t('studentAgenda.dueDate') }} {{ formatDate(item.due_date || item.date) }}</div>
          </q-card-section>

          <q-card-section>
            <div class="text-h6 text-weight-bold" :class="{ 'text-strike text-grey': isItemCompleted(item.id) }">
              {{ item.description || item.title }}
            </div>
            <div class="text-subtitle2 text-grey-7" v-if="item.teacher_name">{{ $t('studentAgenda.teacher') }} {{ item.teacher_name }}</div>
            <div class="text-caption text-grey-6" v-if="item.subject_name">{{ $t('studentAgenda.subject') }} {{ item.subject_name }}</div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="between">
            <q-checkbox
              :model-value="isItemCompleted(item.id)"
              @update:model-value="(val) => toggleCompleted(item.id, val)"
              :label="$t('studentAgenda.markAsDone')"
              color="positive"
            />
            <q-btn flat round icon="event" color="grey" />
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAgendaStore } from 'src/stores/agenda'

const agendaStore = useAgendaStore()
const filterType = ref('all')
const completedMap = ref({})

const loadCompletedFromStorage = () => {
  try {
    const stored = localStorage.getItem('agenda_completed_tasks')
    if (stored) {
      completedMap.value = JSON.parse(stored)
    }
  } catch (e) {
    console.error('Error loading agenda completed state:', e)
  }
}

const isItemCompleted = (id) => {
  return !!completedMap.value[id]
}

const toggleCompleted = async (id, val) => {
  completedMap.value[id] = val
  try {
    localStorage.setItem('agenda_completed_tasks', JSON.stringify(completedMap.value))
    await agendaStore.updateEvent(id, { is_completed: val }).catch(() => {})
  } catch (e) {
    console.error('Error persisting completion:', e)
  }
}

const loadAgenda = async () => {
  try {
    await agendaStore.fetchAgenda()
  } catch (e) {
    console.error('Error loading agenda:', e)
  }
}

const getBadgeColor = (type) => {
  const t = (type || '').toLowerCase()
  if (t === 'verifica') return 'negative'
  if (t === 'interrogazione') return 'warning'
  if (t === 'compito') return 'primary'
  return 'secondary'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  return d.toLocaleDateString('it-IT', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

const homeworks = computed(() => {
  return (agendaStore.events || []).map(e => ({
    id: e.id,
    type: e.event_type || e.type || 'compito',
    description: e.description || e.title || '',
    due_date: e.date || e.due_date || '',
    teacher_name: e.teacher_name || e.created_by_name || '',
    subject_name: e.subject_name || ''
  }))
})

const filteredHomeworks = computed(() => {
  if (filterType.value === 'all') return homeworks.value
  return homeworks.value.filter(h => (h.type || '').toLowerCase() === filterType.value)
})

onMounted(() => {
  loadCompletedFromStorage()
  loadAgenda()
})
</script>

<style scoped>
.bordered {
  border: 1px solid rgba(0, 0, 0, 0.08);
}
</style>
