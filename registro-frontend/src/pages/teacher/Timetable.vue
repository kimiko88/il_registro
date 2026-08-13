<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="schedule" color="primary" class="q-mr-sm" />
          Orario Scolastico
        </div>
        <div class="text-caption text-grey">Visualizza il tuo orario di insegnamento o gestisci l'orario per classe</div>
      </div>
      <div class="row items-center gap-2">
        <q-btn-toggle
          v-model="viewMode"
          toggle-color="primary"
          flat
          dense
          no-caps
          class="bg-slate-200 rounded-lg q-pa-xs"
          :options="[
            { label: 'Il Mio Orario', value: 'my_schedule', icon: 'person' },
            { label: 'Orario per Classe', value: 'class', icon: 'groups' }
          ]"
          @update:model-value="onViewModeChange"
        />
        <q-btn
          v-if="viewMode === 'class' && selectedClass"
          :label="isEditing ? 'Vista Lettura' : 'Modifica Orario'"
          :icon="isEditing ? 'visibility' : 'edit'"
          :color="isEditing ? 'secondary' : 'primary'"
          outline
          no-caps
          class="rounded-lg"
          @click="isEditing = !isEditing"
        />
      </div>
    </div>

    <!-- Filters for Class Mode -->
    <q-card v-if="viewMode === 'class'" class="q-mb-md shadow-1">
      <q-card-section class="row items-center q-gutter-md q-py-sm">
        <q-select
          v-model="selectedClass"
          :options="classOptions"
          option-value="id"
          option-label="label"
          emit-value map-options
          label="Classe"
          dense outlined
          style="min-width:220px"
        />
      </q-card-section>
    </q-card>

    <!-- Timetable Grid / Editor -->
    <q-card v-if="loading" class="text-center q-pa-xl shadow-1">
      <q-spinner-dots color="primary" size="60px" />
    </q-card>

    <div v-else-if="viewMode === 'my_schedule'">
      <q-card v-if="myScheduleEntries.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="event_busy" size="80px" class="q-mb-md" />
        <div class="text-h6">Nessun orario docente trovato</div>
        <div class="text-caption q-mb-md">Non risultano lezioni assegnate al tuo profilo. Seleziona "Orario per Classe" per configurarle.</div>
        <q-btn label="Configura Orario Classe" color="primary" icon="edit_calendar" no-caps class="rounded-lg" @click="viewMode = 'class'; isEditing = true" />
      </q-card>

      <q-card v-else class="shadow-soft overflow-hidden">
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
                <td class="hour-cell text-weight-bold">{{ hour }}ª</td>
                <td 
                  v-for="day in 6" 
                  :key="day" 
                  class="schedule-cell"
                  :class="{ 'has-content': getMyCell(day, hour) }"
                >
                  <div v-if="getMyCell(day, hour)" class="cell-content">
                    <div class="text-subtitle2 text-weight-bold text-primary">{{ getMyCell(day, hour).subject_name }}</div>
                    <div class="text-caption text-weight-bold text-slate-700">{{ getMyCellName(getMyCell(day, hour)) }}</div>
                    <div v-if="getMyCell(day, hour).room" class="text-caption text-grey-6 text-weight-medium">
                      <q-icon name="room" size="xs" class="q-mr-xs" />Aula: {{ getMyCell(day, hour).room }}
                    </div>
                  </div>
                  <div v-else class="empty-cell">-</div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </q-card>
    </div>

    <div v-else>
      <q-card v-if="!selectedClass" class="text-center q-pa-xl text-grey shadow-1">
        <q-icon name="touch_app" size="80px" class="q-mb-md" />
        <div class="text-h6">Seleziona una classe</div>
        <div class="text-caption">Scegli una classe dall'elenco per visualizzare o inserire l'orario delle lezioni.</div>
      </q-card>

      <div v-else-if="isEditing">
        <ScheduleGrid
          :assignments="classAssignments"
          :initial-schedule="scheduleEntries"
          :loading="saving"
          @save="onSaveSchedule"
        />
      </div>

      <q-card v-else-if="scheduleEntries.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="event_busy" size="80px" class="q-mb-md" />
        <div class="text-h6">Orario non ancora configurato</div>
        <div class="text-caption q-mb-md">Non risulta un orario scolastico salvato per questa classe. Puoi inserirlo ora.</div>
        <q-btn label="Configura Orario Ora" color="primary" icon="edit_calendar" no-caps class="rounded-lg" @click="isEditing = true" />
      </q-card>

      <q-card v-else class="shadow-soft overflow-hidden">
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
                <td class="hour-cell text-weight-bold">{{ hour }}ª</td>
                <td 
                  v-for="day in 6" 
                  :key="day" 
                  class="schedule-cell"
                  :class="{ 'has-content': getCell(day, hour) }"
                >
                  <div v-if="getCell(day, hour)" class="cell-content">
                    <div class="text-subtitle2 text-weight-bold text-primary">{{ getCell(day, hour).subject_name }}</div>
                    <div class="text-caption text-grey-7">{{ getCell(day, hour).teacher_name }}</div>
                    <div v-if="getCell(day, hour).room" class="text-caption text-grey-6 text-weight-medium">
                      <q-icon name="room" size="xs" class="q-mr-xs" />Aula: {{ getCell(day, hour).room }}
                    </div>
                  </div>
                  <div v-else class="empty-cell">-</div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </q-card>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useClassesStore } from 'src/stores/classes'
import api from 'src/services/api'
import adminService from 'src/services/adminService'
import ScheduleGrid from 'src/components/Secretary/ScheduleGrid.vue'

const $q = useQuasar()
const classesStore = useClassesStore()

const loading = ref(false)
const saving = ref(false)
const isEditing = ref(false)
const viewMode = ref('my_schedule')
const selectedClass = ref(null)
const classOptions = ref([])
const scheduleEntries = ref([])
const myScheduleEntries = ref([])
const classAssignments = ref([])

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
  classOptions.value = classesStore.classes || []
  if (classOptions.value.length > 0) {
    selectedClass.value = classOptions.value[0].id
  }
  loading.value = false
})

const fetchMySchedule = async () => {
  try {
    const res = await api.get('/timetables/my-schedule')
    myScheduleEntries.value = res.data || []
  } catch (e) {
    console.error('Failed fetching my schedule', e)
    myScheduleEntries.value = []
  }
}

const onViewModeChange = async (val) => {
  if (val === 'my_schedule') {
    loading.value = true
    await fetchMySchedule()
    loading.value = false
  } else if (selectedClass.value) {
    loading.value = true
    await Promise.all([fetchSchedule(), fetchAssignments()])
    loading.value = false
  }
}

watch(selectedClass, async (newVal) => {
  if (newVal) {
    loading.value = true
    await Promise.all([fetchSchedule(), fetchAssignments()])
    loading.value = false
  } else {
    scheduleEntries.value = []
    classAssignments.value = []
  }
})

const fetchSchedule = async () => {
  try {
    const res = await adminService.getClassSchedule(selectedClass.value)
    scheduleEntries.value = res.data || []
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: 'Impossibile caricare l\'orario scolastico' })
  }
}

const fetchAssignments = async () => {
  try {
    const res = await adminService.getClassSubjects(selectedClass.value)
    classAssignments.value = res.data || []
  } catch (e) {
    console.error(e)
    classAssignments.value = []
  }
}

const onSaveSchedule = async (entries) => {
  saving.value = true
  try {
    const formattedEntries = entries.map(e => ({
      day_of_week: e.day_of_week,
      hour_index: e.hour_index,
      subject_id: e.subject_id,
      teacher_id: e.teacher_id || null,
      room: e.room || ''
    }))
    await adminService.saveClassSchedule(selectedClass.value, { entries: formattedEntries })
    $q.notify({ type: 'positive', message: 'Orario scolastico salvato con successo' })
    await fetchSchedule()
    isEditing.value = false
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante il salvataggio dell\'orario' })
  } finally {
    saving.value = false
  }
}

const getCell = (day, hour) => {
  return scheduleEntries.value.find(e => e.day_of_week === day && e.hour_index === hour)
}

const getMyCell = (day, hour) => {
  return myScheduleEntries.value.find(e => e.day_of_week === day && e.hour_index === hour)
}

const getMyCellName = (cell) => {
  if (!cell) return ''
  if (cell.class_name) return cell.class_name
  const found = classOptions.value.find(c => c.id === cell.class_id)
  return found ? (found.label || `${found.name}${found.section}`) : (cell.class_id ? `Classe ${cell.class_id.substring(0, 4)}` : '')
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
  border: 1px solid rgba(0,0,0,0.06);
  padding: 12px;
  text-align: center;
  vertical-align: middle;
}
.timetable-grid th {
  background-color: var(--q-primary);
  color: white;
  font-weight: 600;
}
.hour-col {
  width: 80px;
  min-width: 80px;
}
.day-col {
  width: calc((100% - 80px) / 6);
  min-width: 110px;
}
.hour-cell {
  background-color: #f8fafc;
  color: #64748b;
  font-size: 1rem;
  width: 80px;
  min-width: 80px;
}
.schedule-cell {
  height: 80px;
  vertical-align: middle;
  transition: background-color 0.2s;
}
.schedule-cell.has-content {
  background-color: #f0f7ff;
}
.cell-content {
  padding: 4px;
}
.empty-cell {
  color: #cbd5e1;
}
</style>
