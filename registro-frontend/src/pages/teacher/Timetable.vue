<template>
  <q-page class="q-pa-md bg-slate-50">
    <!-- Header & Summary Cards -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold text-slate-800 row items-center gap-2">
          <q-icon name="schedule" color="primary" size="md" />
          {{ t('timetablePage.title') }}
        </div>
        <div class="text-caption text-slate-500">
          {{ t('timetablePage.subtitle') }}
        </div>
      </div>
      <div class="row items-center gap-2">
        <q-btn
          :label="t('timetablePage.printTimetable')"
          icon="print"
          color="primary"
          outline
          no-caps
          class="rounded-lg"
          @click="printSchedule"
        />
        <q-btn
          icon="refresh"
          color="slate"
          flat
          round
          @click="fetchMySchedule"
        >
          <q-tooltip>Aggiorna Orario</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Summary Stats Chips -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card flat class="bg-white rounded-xl border border-slate-200 q-pa-sm">
          <div class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">{{ t('timetablePage.title') }}</div>
              <div class="text-h6 text-weight-bold text-primary">{{ myScheduleEntries.length }} {{ t('timetablePage.hour') }}</div>
            </div>
            <q-avatar color="indigo-1" text-color="primary" icon="access_time" />
          </div>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat class="bg-white rounded-xl border border-slate-200 q-pa-sm">
          <div class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">{{ t('nav.classes') }}</div>
              <div class="text-h6 text-weight-bold text-emerald-600">{{ uniqueClassesCount }} {{ t('nav.classes') }}</div>
            </div>
            <q-avatar color="emerald-1" text-color="positive" icon="groups" />
          </div>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat class="bg-white rounded-xl border border-slate-200 q-pa-sm">
          <div class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">{{ t('udaPage.subject') }}</div>
              <div class="text-h6 text-weight-bold text-amber-600">{{ uniqueSubjectsCount }}</div>
            </div>
            <q-avatar color="amber-1" text-color="warning" icon="menu_book" />
          </div>
        </q-card>
      </div>
    </div>

    <!-- Timetable Grid -->
    <q-card v-if="loading" class="text-center q-pa-xl shadow-1 rounded-xl">
      <q-spinner-dots color="primary" size="60px" />
    </q-card>

    <q-card v-else-if="myScheduleEntries.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1 rounded-xl bg-white">
      <q-icon name="event_busy" size="80px" color="slate-400" class="q-mb-md" />
      <div class="text-h6 text-slate-700">{{ t('timetablePage.freeSlot') }}</div>
    </q-card>

    <q-card v-else class="shadow-soft rounded-xl overflow-hidden bg-white border border-slate-200">
      <div class="grid-scroll">
        <table class="timetable-grid">
          <thead>
            <tr>
              <th class="hour-col text-outfit">{{ t('timetablePage.hour') }}</th>
              <th v-for="day in days" :key="day.value" class="day-col text-outfit">
                {{ day.label }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="hour in 8" :key="hour">
              <td class="hour-cell text-weight-bold">{{ hour }}ª {{ t('timetablePage.hour') }}</td>
              <td 
                v-for="day in days" 
                :key="day.value" 
                class="timetable-cell"
              >
                <div v-if="getMyCell(day.value, hour)" class="cell-content bg-indigo-50 border-indigo-200">
                  <div class="text-weight-bold text-primary">{{ getMyCell(day.value, hour).subject_name }}</div>
                  <div class="text-caption text-grey-8">{{ getMyCellName(getMyCell(day.value, hour)) }}</div>
                  <div v-if="getMyCell(day.value, hour).room" class="text-caption text-grey-6">{{ t('timetablePage.classroom') }}: {{ getMyCell(day.value, hour).room }}</div>
                </div>
                <div v-else class="empty-cell text-grey-4 text-caption">
                  {{ t('timetablePage.freeSlot') }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useClassesStore } from '@/stores/classes'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const classesStore = useClassesStore()

const loading = ref(false)
const myScheduleEntries = ref([])

const days = computed(() => [
  { label: t('timetablePage.monday'), value: 1 },
  { label: t('timetablePage.tuesday'), value: 2 },
  { label: t('timetablePage.wednesday'), value: 3 },
  { label: t('timetablePage.thursday'), value: 4 },
  { label: t('timetablePage.friday'), value: 5 },
  { label: t('timetablePage.saturday'), value: 6 }
])

onMounted(async () => {
  loading.value = true
  await Promise.all([
    fetchMySchedule(),
    classesStore.fetchAssignedClasses()
  ])
  loading.value = false
})

const fetchMySchedule = async () => {
  try {
    const res = await api.get('/timetables/my-schedule')
    myScheduleEntries.value = res.data || []
  } catch (e) {
    console.error('Failed fetching my schedule', e)
    myScheduleEntries.value = []
    $q.notify({ type: 'negative', message: 'Impossibile caricare l\'orario docente' })
  }
}

const uniqueClassesCount = computed(() => {
  const set = new Set(myScheduleEntries.value.map(e => e.class_id || e.class_name))
  return set.size
})

const uniqueSubjectsCount = computed(() => {
  const set = new Set(myScheduleEntries.value.map(e => e.subject_id || e.subject_name))
  return set.size
})

const getMyCell = (day, hour) => {
  return myScheduleEntries.value.find(e => e.day_of_week === day && e.hour_index === hour)
}

const getMyCellName = (cell) => {
  if (!cell) return ''
  if (cell.class_name) return cell.class_name
  const found = (classesStore.classes || []).find(c => c.id === cell.class_id)
  if (found) return found.label || `${found.name}${found.section}`
  return cell.class_id ? `Classe ${cell.class_id.substring(0, 4)}` : ''
}

const printSchedule = () => {
  window.print()
}
</script>

<style scoped>
.grid-scroll {
  width: 100%;
  overflow-x: auto;
}
.timetable-grid {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  min-width: 750px;
}
.timetable-grid th, .timetable-grid td {
  border: 1px solid rgba(226, 232, 240, 0.8);
  padding: 12px 8px;
  text-align: center;
  vertical-align: middle;
}
.timetable-grid th {
  background-color: #4f46e5;
  color: white;
  font-weight: 600;
  font-size: 0.95rem;
}
.hour-col {
  width: 90px;
  min-width: 90px;
}
.day-col {
  width: calc((100% - 90px) / 6);
  min-width: 120px;
}
.hour-cell {
  background-color: #f8fafc;
  color: #475569;
  font-size: 0.9rem;
  width: 90px;
  min-width: 90px;
}
.schedule-cell {
  height: 85px;
  vertical-align: middle;
  transition: background-color 0.2s;
}
.schedule-cell.has-content {
  background-color: #f0f7ff;
}
.cell-content {
  padding: 4px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.empty-cell {
  color: #cbd5e1;
  font-size: 1.1rem;
}
</style>
