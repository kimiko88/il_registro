<template>
  <div class="teacher-schedule-grid-container">
    <!-- Top summary & Save Action Bar -->
    <div class="row items-center justify-between bg-emerald-50 p-3 rounded-xl border border-emerald-200 q-mb-md gap-2">
      <div class="row items-center gap-2">
        <q-icon name="schedule" color="emerald" size="22px" class="text-emerald-700" />
        <span class="text-subtitle2 text-emerald-950 font-bold">Totale Ore Docente:</span>
        <q-badge color="emerald-7" class="text-bold q-px-sm bg-emerald-600 text-white">{{ totalHours }} ore/settimana</q-badge>
      </div>
      <q-btn
        label="Salva Orario Docente"
        color="positive"
        unelevated
        icon="save"
        class="rounded-lg q-px-md shadow-xs font-bold"
        no-caps
        :loading="loading"
        @click="$emit('save', gridEntries)"
      />
    </div>

    <!-- Timetable Grid Table -->
    <div class="grid-scroll border border-slate-300 rounded-xl overflow-hidden shadow-sm bg-white">
      <table class="timetable-grid">
        <thead>
          <tr class="bg-slate-100 border-b border-slate-300">
            <th class="hour-col border-r border-slate-300 py-3 text-xs font-bold text-slate-700 text-center">Ora</th>
            <th
              v-for="day in days"
              :key="day.value"
              class="day-col border-r border-slate-200 py-3 text-xs font-bold text-slate-700 uppercase tracking-wider text-center"
            >
              {{ day.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="hour in 8"
            :key="hour"
            class="border-b border-slate-200"
          >
            <td class="hour-cell font-bold bg-slate-50 border-r border-slate-300 text-center text-slate-700 text-xs py-2">
              {{ hour }}ª ora
            </td>
            <td
              v-for="day in 6"
              :key="day"
              class="schedule-cell border-r border-slate-200 relative p-1.5 cursor-pointer transition-colors"
              :class="{
                'bg-emerald-50/70 hover:bg-emerald-100/80': getCell(day, hour),
                'bg-white hover:bg-slate-50': !getCell(day, hour)
              }"
              @click="editCell(day, hour)"
            >
              <div v-if="getCell(day, hour)" class="cell-content p-2 rounded-lg bg-white border border-emerald-300 shadow-2xs row items-center justify-between gap-1">
                <div class="overflow-hidden">
                  <div class="text-xs font-bold text-emerald-950 truncate">{{ getCell(day, hour).subject_name }}</div>
                  <div class="text-[11px] font-semibold text-emerald-700 truncate mt-0.5">Classe: {{ getCellClassName(getCell(day, hour)) }}</div>
                  <div v-if="getCell(day, hour).room" class="text-[10px] text-slate-500 font-medium mt-0.5 truncate">Aula: {{ getCell(day, hour).room }}</div>
                </div>
                <q-btn
                  flat round dense icon="close"
                  size="xs"
                  color="negative"
                  @click.stop="removeCell(day, hour)"
                />
              </div>
              <div v-else class="add-placeholder flex items-center justify-center h-full min-h-[46px] text-slate-400 hover:text-emerald-700 transition-colors">
                <span class="text-caption italic font-medium">+ Assegna</span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Cell Edit Dialog Pop-up -->
    <q-dialog v-model="cellDialogVisible">
      <q-card style="width: min(450px, 90vw)" class="rounded-xl overflow-hidden shadow-24 border-slate-300">
        <q-card-section class="bg-emerald-700 text-white row items-center justify-between q-py-md">
          <div class="text-subtitle1 font-bold row items-center gap-2">
            <q-icon name="edit_calendar" />
            {{ editingCell ? `${days[editingCell.day - 1].label} - ${editingCell.hour}ª Ora` : 'Assegna Ora Docente' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md space-y-3">
          <q-select
            v-model="selectedClassId"
            :options="formattedClasses"
            label="Classe *"
            outlined
            dense
            emit-value
            map-options
            class="bg-white"
          />

          <q-select
            v-model="selectedSubjectId"
            :options="formattedSubjects"
            label="Materia *"
            outlined
            dense
            emit-value
            map-options
            class="bg-white"
          />

          <q-input
            v-model="room"
            label="Aula (opzionale)"
            outlined
            dense
            placeholder="Es. Lab Informatica, Aula 2A"
            class="bg-white"
          />
        </q-card-section>

        <q-card-actions align="between" class="q-pa-md bg-slate-50 border-t border-slate-100">
          <q-btn
            v-if="editingCell && getCell(editingCell.day, editingCell.hour)"
            label="Rimuovi"
            color="negative"
            flat
            no-caps
            @click="removeCell(editingCell.day, editingCell.hour); cellDialogVisible = false"
          />
          <div v-else />

          <div class="row q-gutter-sm">
            <q-btn label="Annulla" flat no-caps v-close-popup />
            <q-btn label="Conferma" color="positive" unelevated class="rounded-lg q-px-md font-bold" no-caps @click="applyCell" />
          </div>
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const props = defineProps({
  classes: { type: Array, default: () => [] },
  subjects: { type: Array, default: () => [] },
  initialSchedule: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false }
})

defineEmits(['save'])

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
const selectedClassId = ref(null)
const selectedSubjectId = ref(null)
const room = ref('')
const cellDialogVisible = ref(false)

watch(() => props.initialSchedule, (val) => {
  if (Array.isArray(val)) {
    gridEntries.value = val.map(e => ({
      day_of_week: e.day_of_week,
      hour_index: e.hour_index,
      class_id: e.class_id,
      class_name: e.class_name,
      subject_id: e.subject_id,
      subject_name: e.subject_name,
      room: e.room || ''
    }))
  }
}, { immediate: true })

const formattedClasses = computed(() => {
  return props.classes.map(c => ({
    label: c.label || `Classe ${c.name || ''}${c.section || ''}`,
    value: c.id || c.value
  }))
})

const formattedSubjects = computed(() => {
  return props.subjects.map(s => ({
    label: s.label || s.name,
    value: s.value || s.id
  }))
})

const getCell = (day, hour) => {
  return gridEntries.value.find(e => e.day_of_week === day && e.hour_index === hour)
}

const getCellClassName = (cell) => {
  if (!cell) return ''
  if (cell.class_name) return cell.class_name
  const found = props.classes.find(c => (c.id || c.value) === cell.class_id)
  return found ? (found.label || `${found.name || ''}${found.section || ''}`) : (cell.class_id ? cell.class_id.substring(0, 5) : '')
}

const editCell = (day, hour) => {
  editingCell.value = { day, hour }
  const existing = getCell(day, hour)
  if (existing) {
    selectedClassId.value = existing.class_id
    selectedSubjectId.value = existing.subject_id
    room.value = existing.room || ''
  } else {
    selectedClassId.value = formattedClasses.value[0]?.value || null
    selectedSubjectId.value = formattedSubjects.value[0]?.value || null
    room.value = ''
  }
  cellDialogVisible.value = true
}

const applyCell = () => {
  if (!editingCell.value || !selectedClassId.value || !selectedSubjectId.value) return

  // Remove existing entry for this cell
  gridEntries.value = gridEntries.value.filter(e => 
    !(e.day_of_week === editingCell.value.day && e.hour_index === editingCell.value.hour)
  )

  const selectedClassObj = props.classes.find(c => (c.id || c.value) === selectedClassId.value)
  const selectedSubjectObj = props.subjects.find(s => (s.value || s.id) === selectedSubjectId.value)

  gridEntries.value.push({
    day_of_week: editingCell.value.day,
    hour_index: editingCell.value.hour,
    class_id: selectedClassId.value,
    class_name: selectedClassObj ? (selectedClassObj.label || `${selectedClassObj.name || ''}${selectedClassObj.section || ''}`) : '',
    subject_id: selectedSubjectId.value,
    subject_name: selectedSubjectObj ? (selectedSubjectObj.label || selectedSubjectObj.name) : '',
    room: room.value
  })
  
  cellDialogVisible.value = false
}

const removeCell = (day, hour) => {
  gridEntries.value = gridEntries.value.filter(e => 
    !(e.day_of_week === day && e.hour_index === hour)
  )
}

const totalHours = computed(() => gridEntries.value.length)
</script>

<style scoped>
.teacher-schedule-grid-container {
  width: 100%;
}

.grid-scroll {
  width: 100%;
  overflow-x: auto;
}

.timetable-grid {
  width: 100%;
  border-collapse: collapse;
  min-width: 750px;
}

.hour-col { width: 80px; min-width: 80px; }
.day-col { width: calc((100% - 80px) / 6); min-width: 110px; }

.schedule-cell {
  height: 64px;
}
</style>
