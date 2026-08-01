<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="schedule" color="primary" class="q-mr-sm" />
          Orario Scolastico
        </div>
        <div class="text-caption text-grey">Seleziona una classe per visualizzare il suo orario delle lezioni</div>
      </div>
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md shadow-1">
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

    <!-- Timetable Grid -->
    <q-card v-if="loading" class="text-center q-pa-xl shadow-1">
      <q-spinner-dots color="primary" size="60px" />
    </q-card>

    <div v-else>
      <q-card v-if="!selectedClass" class="text-center q-pa-xl text-grey shadow-1">
        <q-icon name="touch_app" size="80px" class="q-mb-md" />
        <div class="text-h6">Seleziona una classe</div>
        <div class="text-caption">Scegli una classe dall'elenco per visualizzare l'orario delle lezioni.</div>
      </q-card>

      <q-card v-else-if="scheduleEntries.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="event_busy" size="80px" class="q-mb-md" />
        <div class="text-h6">Orario non ancora configurato</div>
        <div class="text-caption">La segreteria non ha ancora configurato l'orario per questa classe.</div>
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
import adminService from 'src/services/adminService'

const $q = useQuasar()
const classesStore = useClassesStore()

const loading = ref(false)
const selectedClass = ref(null)
const classOptions = ref([])
const scheduleEntries = ref([])

const days = [
  { label: 'Lunedì', value: 1 },
  { label: 'Martedì', value: 2 },
  { label: 'Mercoledì', value: 3 },
  { label: 'Giovedì', value: 4 },
  { label: 'Venerdì', value: 5 },
  { label: 'Sabato', value: 6 }
]

onMounted(async () => {
  await classesStore.fetchAssignedClasses()
  classOptions.value = classesStore.classes || []
  if (classOptions.value.length > 0) {
    selectedClass.value = classOptions.value[0].id
  }
})

watch(selectedClass, async (newVal) => {
  if (newVal) {
    loading.value = true
    await fetchSchedule()
    loading.value = false
  } else {
    scheduleEntries.value = []
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

const getCell = (day, hour) => {
  return scheduleEntries.value.find(e => e.day_of_week === day && e.hour_index === hour)
}
</script>

<style scoped>
.grid-scroll {
  overflow-x: auto;
}
.timetable-grid {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  min-width: 800px;
}
.timetable-grid th, .timetable-grid td {
  border: 1px solid rgba(0,0,0,0.06);
  padding: 12px;
  text-align: center;
}
.timetable-grid th {
  background-color: var(--q-primary);
  color: white;
  font-weight: 600;
}
.hour-col {
  width: 80px;
}
.day-col {
  width: 15%;
}
.hour-cell {
  background-color: #f8fafc;
  color: #64748b;
  font-size: 1.1rem;
}
.schedule-cell {
  height: 90px;
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
