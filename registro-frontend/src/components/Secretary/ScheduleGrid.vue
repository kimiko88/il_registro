<template>
  <div class="schedule-grid-container">
    <!-- Top summary & Save Action Bar -->
    <div class="row items-center justify-between bg-slate-100 p-3 rounded-xl border border-slate-300 q-mb-md gap-2">
      <div class="row items-center gap-2">
        <q-icon name="schedule" color="primary" size="22px" />
        <span class="text-subtitle2 text-slate-800 font-bold">Riepilogo Ore Settimanali:</span>
        <q-badge color="primary" class="text-bold q-px-sm">{{ totalHours }} ore/settimana</q-badge>
      </div>
      <q-btn
        label="Salva Orario Settimanale"
        color="primary"
        unelevated
        icon="save"
        class="rounded-lg q-px-md shadow-xs"
        no-caps
        :loading="loading"
        @click="$emit('save', gridEntries)"
      />
    </div>

    <!-- Timetable Grid Table -->
    <div class="grid-scroll border-2 border-slate-400 rounded-xl overflow-hidden shadow-xs bg-white">
      <table class="timetable-grid">
        <thead>
          <tr class="bg-slate-100 border-b-2 border-slate-400">
            <th class="hour-col border-r-2 border-slate-400 py-3 text-xs font-bold text-slate-700">Ora</th>
            <th
              v-for="day in days"
              :key="day.value"
              class="day-col border-r border-slate-300 py-3 text-xs font-bold text-slate-700 uppercase tracking-wider"
            >
              {{ day.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="hour in 8"
            :key="hour"
            class="border-b border-slate-300"
          >
            <td class="hour-cell font-bold bg-slate-100 border-r-2 border-slate-400 text-center text-slate-700 text-xs">
              {{ hour }}ª ora
            </td>
            <td
              v-for="day in 6"
              :key="day"
              class="schedule-cell border-r border-slate-300 relative p-1 cursor-pointer transition-colors"
              :class="{
                'bg-blue-50/80 hover:bg-blue-100/90': getCell(day, hour),
                'bg-white hover:bg-slate-100': !getCell(day, hour)
              }"
              @click="editCell(day, hour)"
            >
              <div v-if="getCell(day, hour)" class="cell-content p-2 rounded-lg bg-white border border-blue-300 shadow-2xs row items-center justify-between">
                <div>
                  <div class="text-caption font-bold text-blue-900 leading-tight">{{ getCell(day, hour).subject_name }}</div>
                  <div class="text-[10px] text-slate-600 truncate mt-0.5">{{ getCell(day, hour).teacher_name || 'Docente N/D' }}</div>
                  <div v-if="getCell(day, hour).room" class="text-[10px] text-slate-500 font-semibold mt-0.5">Aula: {{ getCell(day, hour).room }}</div>
                </div>
                <q-btn
                  flat round dense icon="close"
                  size="xs"
                  color="negative"
                  :aria-label="t('common.remove') || 'Rimuovi'"
                  @click.stop="removeCell(day, hour)"
                />
              </div>
              <div v-else class="add-placeholder flex items-center justify-center h-full min-h-[46px] text-slate-400 hover:text-primary transition-colors">
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
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div class="text-subtitle1 font-bold">
            <q-icon name="edit_calendar" class="q-mr-xs" />
            {{ editingCell ? `${days[editingCell.day - 1]?.label} - ${editingCell.hour}ª Ora` : 'Assegna Ora' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-md space-y-3">
          <div v-if="assignmentOptions.length === 0" class="bg-amber-50 border border-amber-300 p-3 rounded-lg text-amber-900 text-caption">
            ⚠️ Nessuna materia/cattedra assegnata a questa classe. Per poter comporre l'orario, prima aggiungi le materie e i docenti nella sezione Cattedre.
          </div>

          <q-select
            v-model="selectedAssignment"
            :options="assignmentOptions"
            label="Materia & Docente *"
            outlined
            dense
            emit-value
            map-options
            class="bg-white"
          />

          <q-input
            v-model="room"
            label="Aula (es. Lab 2, Aula Magna)"
            outlined
            dense
            class="bg-white"
          />
        </q-card-section>

        <q-card-actions align="between" class="bg-slate-100 q-px-md q-py-sm border-t">
          <q-btn
            v-if="editingCell && getCell(editingCell.day, editingCell.hour)"
            flat
            color="negative"
            :label="t('common.remove') || 'Rimuovi'"
            icon="delete"
            no-caps
            @click="removeCell(editingCell.day, editingCell.hour); cellDialogVisible = false"
          />
          <div v-else></div>
          <div class="row q-gutter-sm">
            <q-btn flat :label="t('common.cancel') || 'Annulla'" color="grey-7" no-caps v-close-popup />
            <q-btn unelevated color="primary" :label="t('common.confirm') || 'Conferma'" no-caps @click="applyCell" />
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
  assignments: { type: Array, default: () => [] },
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
const selectedAssignment = ref(null)
const room = ref('')
const cellDialogVisible = ref(false)

// Initialize grid from props
watch(() => props.initialSchedule, (val) => {
  if (Array.isArray(val)) {
    gridEntries.value = val.map(e => ({
      day_of_week: e.day_of_week,
      hour_index: e.hour_index,
      subject_id: e.subject_id,
      subject_name: e.subject_name,
      teacher_id: e.teacher_id,
      teacher_name: e.teacher_name,
      room: e.room || ''
    }))
  }
}, { immediate: true })

const assignmentOptions = computed(() => {
  return props.assignments.map(a => ({
    label: `${a.subject_name} (${a.teacher_name || 'Docente N/D'})`,
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
    const match = assignmentOptions.value.find(o => o.subject_id === existing.subject_id && o.teacher_id === existing.teacher_id)
    selectedAssignment.value = match ? match.value : null
    room.value = existing.room || ''
  } else {
    selectedAssignment.value = assignmentOptions.value[0]?.value || null
    room.value = ''
  }
  cellDialogVisible.value = true
}

const applyCell = () => {
  if (!editingCell.value) return

  // Remove existing entry for this cell
  gridEntries.value = gridEntries.value.filter(e => 
    !(e.day_of_week === editingCell.value.day && e.hour_index === editingCell.value.hour)
  )

  if (selectedAssignment.value) {
    const selectedOption = assignmentOptions.value.find(o => o.value === selectedAssignment.value)
    if (selectedOption) {
      gridEntries.value.push({
        day_of_week: editingCell.value.day,
        hour_index: editingCell.value.hour,
        subject_id: selectedOption.subject_id,
        subject_name: selectedOption.subject_name,
        teacher_id: selectedOption.teacher_id,
        teacher_name: selectedOption.teacher_name,
        room: room.value
      })
    }
  }
  
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
.schedule-grid-container {
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
