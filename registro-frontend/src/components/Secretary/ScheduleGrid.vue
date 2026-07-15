<template>
  <div class="schedule-grid-container">
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-8">
        <div class="grid-scroll">
          <table class="timetable-grid">
            <thead>
              <tr>
                <th class="hour-col">Ora</th>
                <th v-for="day in days" :key="day.value" class="day-col">
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
                  @click="editCell(day, hour)"
                >
                  <div v-if="getCell(day, hour)" class="cell-content">
                    <div class="text-caption text-weight-bold text-primary">{{ getCell(day, hour).subject_name }}</div>
                    <div class="text-caption opacity-70">{{ getCell(day, hour).teacher_name }}</div>
                    <q-btn 
                      flat round dense icon="close" 
                      size="xs" 
                      class="remove-btn" 
                      @click.stop="removeCell(day, hour)" 
                    />
                  </div>
                  <div v-else class="add-placeholder">
                    <q-icon name="add" size="xs" color="grey-4" />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="col-12 col-md-4">
        <q-card flat class="rounded-2xl border-slate-200 bg-slate-50 q-pa-lg sticky-top">
          <div class="text-subtitle1 text-weight-bold q-mb-md">
            {{ editingCell ? `Modifica: ${days[editingCell.day-1].label}, ${editingCell.hour}ª ora` : 'Seleziona una cella' }}
          </div>
          
          <div v-if="editingCell">
            <q-select
              v-model="selectedAssignment"
              :options="assignmentOptions"
              label="Materia & Docente"
              outlined
              dense
              class="q-mb-md bg-white"
              clearable
            />
            
            <q-input
              v-model="room"
              label="Aula (opzionale)"
              outlined
              dense
              class="q-mb-lg bg-white"
            />

            <div class="row q-gutter-sm">
              <q-btn label="Applica" color="primary" class="col rounded-lg" @click="applyCell" />
              <q-btn label="Chiudi" flat color="slate-400" class="col" @click="editingCell = null" />
            </div>
          </div>
          <div v-else class="text-center q-pa-xl opacity-50">
            <q-icon name="touch_app" size="md" class="q-mb-md" />
            <p>Clicca su una cella della griglia per assegnare una materia</p>
          </div>

          <q-separator class="q-my-xl" />

          <div class="row justify-between items-center q-mb-md">
            <div class="text-weight-bold">Riepilogo Ore</div>
            <q-badge color="indigo">{{ totalHours }} ore/sett</q-badge>
          </div>
          
          <div class="row justify-end q-mt-xl">
            <q-btn 
              label="Salva Orario" 
              color="primary" 
              class="full-width q-py-md rounded-xl shadow-soft" 
              icon="save" 
              :loading="loading"
              @click="$emit('save', gridEntries)" 
            />
          </div>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  assignments: { type: Array, default: () => [] },
  initialSchedule: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false }
})

// eslint-disable-next-line no-unused-vars
const emit = defineEmits(['save'])

const days = [
  { label: 'Lunedì', value: 1 },
  { label: 'Martedì', value: 2 },
  { label: 'Mercoledì', value: 3 },
  { label: 'Giovedì', value: 4 },
  { label: 'Venerdì', value: 5 },
  { label: 'Sabato', value: 6 }
]

const gridEntries = ref([])
const editingCell = ref(null)
const selectedAssignment = ref(null)
const room = ref('')

// Initialize grid from props
watch(() => props.initialSchedule, (val) => {
  gridEntries.value = val.map(e => ({
    day_of_week: e.day_of_week,
    hour_index: e.hour_index,
    subject_id: e.subject_id,
    subject_name: e.subject_name,
    teacher_id: e.teacher_id,
    teacher_name: e.teacher_name,
    room: e.room || ''
  }))
}, { immediate: true })

const assignmentOptions = computed(() => {
  return props.assignments.map(a => ({
    label: `${a.subject_name} (${a.teacher_name || 'N/A'})`,
    value: a.id,
    subject_id: a.subject_id,
    subject_name: a.subject_name,
    teacher_id: a.teacher_id,
    teacher_name: a.teacher_name
  }))
})

const getCell = (day, hour) => {
  return gridEntries.value.find(e => e.day_of_week === day && e.hour_index === hour)
}

const editCell = (day, hour) => {
  editingCell.value = { day, hour }
  const existing = getCell(day, hour)
  if (existing) {
    selectedAssignment.value = assignmentOptions.value.find(o => o.subject_id === existing.subject_id && o.teacher_id === existing.teacher_id)
    room.value = existing.room
  } else {
    selectedAssignment.value = null
    room.value = ''
  }
}

const applyCell = () => {
  if (!editingCell.value) return
  
  // Remove existing
  gridEntries.value = gridEntries.value.filter(e => 
    !(e.day_of_week === editingCell.value.day && e.hour_index === editingCell.value.hour)
  )

  if (selectedAssignment.value) {
    gridEntries.value.push({
      day_of_week: editingCell.value.day,
      hour_index: editingCell.value.hour,
      subject_id: selectedAssignment.value.subject_id,
      subject_name: selectedAssignment.value.subject_name,
      teacher_id: selectedAssignment.value.teacher_id,
      teacher_name: selectedAssignment.value.teacher_name,
      room: room.value
    })
  }
  
  editingCell.value = null
}

const removeCell = (day, hour) => {
  gridEntries.value = gridEntries.value.filter(e => 
    !(e.day_of_week === day && e.hour_index === hour)
  )
}

const totalHours = computed(() => gridEntries.value.length)

</script>

<style scoped>
.schedule-grid-container {
  width: 100%;
}

.grid-scroll {
  overflow-x: auto;
  background: white;
  border-radius: 1rem;
  border: 1px solid #f1f5f9;
}

.timetable-grid {
  width: 100%;
  border-collapse: collapse;
  min-width: 800px;
}

.timetable-grid th, .timetable-grid td {
  border: 1px solid #f1f5f9;
  padding: 0.75rem;
  text-align: center;
}

.timetable-grid thead th {
  background: #f8fafc;
  color: #64748b;
  font-weight: 700;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.05em;
}

.hour-col { width: 60px; }
.day-col { width: calc((100% - 60px) / 6); }

.schedule-cell {
  height: 80px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.schedule-cell:hover {
  background-color: #f1f5f9;
}

.schedule-cell.has-content {
  background-color: #f0f7ff;
}

.cell-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.remove-btn {
  position: absolute;
  top: 2px;
  right: 2px;
  opacity: 0;
  transition: opacity 0.2s;
}

.schedule-cell:hover .remove-btn {
  opacity: 1;
}

.add-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  opacity: 0;
}

.schedule-cell:hover .add-placeholder {
  opacity: 1;
}

.sticky-top {
  position: sticky;
  top: 20px;
}

.shadow-soft {
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.rounded-2xl {
  border-radius: 1.5rem;
}
</style>
