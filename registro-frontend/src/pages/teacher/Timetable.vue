<template>
  <q-page class="q-pa-md bg-slate-50">
    <!-- Header & Summary Cards -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold text-slate-800 row items-center gap-2">
          <q-icon name="schedule" color="primary" size="md" />
          Il Mio Orario di Insegnamento
        </div>
        <div class="text-caption text-slate-500">
          Orario settimanale delle lezioni assegnate con materie, classi e aule
        </div>
      </div>
      <div class="row items-center gap-2">
        <q-btn
          label="Stampa Orario"
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
              <div class="text-caption text-slate-500">Ore Settimanali</div>
              <div class="text-h6 text-weight-bold text-primary">{{ myScheduleEntries.length }} Ore</div>
            </div>
            <q-avatar color="indigo-1" text-color="primary" icon="access_time" />
          </div>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat class="bg-white rounded-xl border border-slate-200 q-pa-sm">
          <div class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">Classi Assegnate</div>
              <div class="text-h6 text-weight-bold text-emerald-600">{{ uniqueClassesCount }} Classi</div>
            </div>
            <q-avatar color="emerald-1" text-color="positive" icon="groups" />
          </div>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat class="bg-white rounded-xl border border-slate-200 q-pa-sm">
          <div class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">Materie Insegnate</div>
              <div class="text-h6 text-weight-bold text-amber-600">{{ uniqueSubjectsCount }} Materie</div>
            </div>
            <q-avatar color="amber-1" text-color="warning" icon="menu_book" />
          </div>
        </q-card>
      </div>
    </div>

    <!-- Timetable Grid -->
    <q-card v-if="loading" class="text-center q-pa-xl shadow-1 rounded-xl">
      <q-spinner-dots color="primary" size="60px" />
      <div class="text-caption text-slate-500 q-mt-sm">Caricamento il tuo orario...</div>
    </q-card>

    <q-card v-else-if="myScheduleEntries.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1 rounded-xl bg-white">
      <q-icon name="event_busy" size="80px" color="slate-400" class="q-mb-md" />
      <div class="text-h6 text-slate-700">Nessuna lezione in orario</div>
      <div class="text-caption text-slate-500 max-w-md mx-auto q-mt-xs">
        Non risultano ancora ore di lezione o cattedre assegnate al tuo profilo per questo anno scolastico. Contatta la Segreteria per l'assegnazione dell'orario.
      </div>
    </q-card>

    <q-card v-else class="shadow-soft rounded-xl overflow-hidden bg-white border border-slate-200">
      <div class="grid-scroll">
        <table class="timetable-grid">
          <thead>
            <tr>
              <th class="hour-col text-outfit">Ora</th>
              <th v-for="day in days" :key="day.value" class="day-col text-outfit">
                {{ day.label }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="hour in 8" :key="hour">
              <td class="hour-cell text-weight-bold">{{ hour }}ª ora</td>
              <td 
                v-for="day in 6" 
                :key="day" 
                class="schedule-cell"
                :class="{ 'has-content': getMyCell(day, hour) }"
              >
                <div v-if="getMyCell(day, hour)" class="cell-content">
                  <div class="text-subtitle2 text-weight-bold text-primary line-clamp-1">
                    {{ getMyCell(day, hour).subject_name }}
                  </div>
                  <div class="q-mt-xs">
                    <q-badge color="indigo-1" text-color="indigo-9" class="text-weight-bold px-2 py-1 rounded-md text-caption">
                      {{ getMyCellName(getMyCell(day, hour)) }}
                    </q-badge>
                  </div>
                  <div v-if="getMyCell(day, hour).room" class="text-caption text-slate-500 text-weight-medium q-mt-xs row items-center justify-center">
                    <q-icon name="room" size="xs" color="slate-400" class="q-mr-xs" />
                    Aula: {{ getMyCell(day, hour).room }}
                  </div>
                </div>
                <div v-else class="empty-cell">-</div>
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
import { useClassesStore } from 'src/stores/classes'
import api from 'src/services/api'

const $q = useQuasar()
const classesStore = useClassesStore()

const loading = ref(false)
const myScheduleEntries = ref([])

const days = [
  { label: 'Lunedì', value: 1 },
  { label: 'Martedì', value: 2 },
  { label: 'Mercoledì', value: 3 },
  { label: 'Giovedì', value: 4 },
  { label: 'Venerdì', value: 5 },
  { label: 'Sabato', value: 6 }
]

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
