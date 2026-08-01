<template>
  <div class="q-pa-md">
    <!-- Controls -->
    <div class="row items-center q-gutter-md q-mb-md">
      <q-select
        v-model="selectedClassId"
        :options="classOptions"
        option-value="id"
        option-label="label"
        emit-value map-options
        label="Classe"
        dense outlined
        style="min-width: 220px"
      />
      <q-input v-model="dateVal" type="date" label="Data" dense outlined style="max-width: 160px" />
      <q-space />
      <q-btn color="primary" icon="save" label="Salva Presenze" :loading="saving" @click="saveAll" :disable="!hasChanges" />
    </div>

    <q-banner v-if="!selectedClassId" class="bg-blue-1 q-mb-md">
      <template #avatar><q-icon name="info" color="primary" /></template>
      Seleziona una classe per registrare le presenze.
    </q-banner>

    <q-table
      v-if="selectedClassId"
      :rows="students"
      :columns="columns"
      row-key="id"
      flat bordered
      :loading="loading"
      hide-bottom
    >
      <template v-slot:body="props">
        <q-tr :props="props">
          <!-- Student Name -->
          <q-td key="name" :props="props">
            <div class="text-weight-bold" v-if="props.row">{{ props.row.name }}</div>
          </q-td>

          <!-- Status Toggle -->
          <q-td key="status" :props="props" style="width:280px">
            <q-btn-toggle
              v-model="attendanceMap[props.row.id].status"
              toggle-color="primary"
              flat dense
              :options="[
                { label: 'Presente', value: 'present', icon: 'check_circle' },
                { label: 'Assente', value: 'absent', icon: 'cancel' },
                { label: 'Ritardo', value: 'late', icon: 'schedule' },
                { label: 'Uscita', value: 'early_exit', icon: 'exit_to_app' }
              ]"
            />
          </q-td>

          <!-- Minutes Late (only for 'late') -->
          <q-td key="minutes" :props="props" style="width:140px">
            <q-input
              v-if="attendanceMap[props.row.id].status === 'late'"
              v-model.number="attendanceMap[props.row.id].minutesLate"
              type="number"
              label="Min. ritardo"
              dense outlined
              min="1" max="120"
              style="max-width:120px"
            />
            <!-- Exit time for early_exit -->
            <q-input
              v-else-if="attendanceMap[props.row.id].status === 'early_exit'"
              v-model="attendanceMap[props.row.id].exitTime"
              type="time"
              label="Ora uscita"
              dense outlined
              style="max-width:120px"
            />
            <span v-else class="text-grey-5">—</span>
          </q-td>

          <!-- Already Justified badge -->
          <q-td key="justified" :props="props" class="text-center">
            <template v-if="attendanceMap[props.row.id].status !== 'present'">
              <q-chip
                v-if="attendanceMap[props.row.id].isJustified"
                color="positive" text-color="white" icon="check" size="sm" dense
              >Giustificato</q-chip>
              <q-btn
                v-else
                flat dense size="sm" icon="fact_check" color="grey"
                label="Giustifica"
                @click="openJustify(props.row)"
              />
            </template>
          </q-td>
        </q-tr>
      </template>
    </q-table>

    <!-- Justify Dialog -->
    <q-dialog v-model="justifyDialog" persistent>
      <q-card style="min-width: 380px">
        <q-card-section>
          <div class="text-h6">Giustifica Assenza</div>
          <div class="text-caption text-grey" v-if="justifyTarget">{{ justifyTarget.name }}</div>
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-select
            v-model="justifyReason"
            :options="['Motivi di Salute', 'Motivi Familiari', 'Visita Medica', 'Altro']"
            label="Motivazione *"
            outlined dense
          />
          <q-input v-model="justifyNotes" label="Note" type="textarea" outlined dense autogrow />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Giustifica" @click="submitJustify" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, watch, computed, onMounted } from 'vue'
import { useQuasar, date } from 'quasar'
import { useClassesStore } from 'src/stores/classes'
import { attendanceService } from 'src/services/attendanceService'

const $q = useQuasar()
const classesStore = useClassesStore()

const selectedClassId = ref(null)
const today = date.formatDate(Date.now(), 'YYYY-MM-DD')
const dateVal = ref(today)

const saving = ref(false)
const loading = ref(false)

const classOptions = ref([])
const students = ref([])
// Map: studentId -> { status, minutesLate, exitTime, isJustified }
const attendanceMap = reactive({})
const initialSnapshot = ref('')

const justifyDialog = ref(false)
const justifyTarget = ref(null)
const justifyReason = ref('Motivi di Salute')
const justifyNotes = ref('')

const columns = [
  { name: 'name', label: 'Studente', field: 'name', align: 'left' },
  { name: 'status', label: 'Stato', field: 'status', align: 'left' },
  { name: 'minutes', label: 'Dettagli', field: 'minutesLate', align: 'left' },
  { name: 'justified', label: 'Giustificazione', field: 'isJustified', align: 'center' }
]

onMounted(async () => {
  await classesStore.fetchAssignedClasses()
  classOptions.value = classesStore.classOptions || []
  if (classOptions.value.length > 0) {
    selectedClassId.value = classOptions.value[0].id
  }
})

watch(selectedClassId, () => {
  if (selectedClassId.value) fetchStudentsAndAttendance()
})

watch(dateVal, () => {
  if (selectedClassId.value) fetchStudentsAndAttendance()
})

const fetchStudentsAndAttendance = async () => {
  loading.value = true
  try {
    // Fetch today's attendance for the class
    const res = await attendanceService.getByClass(selectedClassId.value, dateVal.value)
    const records = res.data || []

    // Initialize attendance map
    records.forEach(r => {
      attendanceMap[r.student_id] = {
        attendanceId: r.id || null,
        status: r.status || 'present',
        minutesLate: r.minutes_late || 0,
        exitTime: r.exit_time || '',
        isJustified: r.is_justified || false
      }
    })

    // Build student list + map from existing records
    const fetched = records.map(r => ({
      id: r.student_id,
      name: r.student_name || r.student_id
    }))
    students.value = fetched
    initialSnapshot.value = JSON.stringify(attendanceMap)
  } catch (e) {
    console.error(e)
    students.value = []
  } finally {
    loading.value = false
  }
}

const hasChanges = computed(() => JSON.stringify(attendanceMap) !== initialSnapshot.value)

const saveAll = async () => {
  saving.value = true
  try {
    const records = students.value.map(s => ({
      student_id: s.id,
      date: dateVal.value,
      status: attendanceMap[s.id]?.status || 'present',
      minutes_late: attendanceMap[s.id]?.status === 'late' ? (attendanceMap[s.id]?.minutesLate || 0) : 0,
      exit_time: attendanceMap[s.id]?.status === 'early_exit' ? (attendanceMap[s.id]?.exitTime || '') : null
    }))

    await attendanceService.markAttendance({
      class_id: selectedClassId.value,
      date: dateVal.value,
      records
    })

    $q.notify({ type: 'positive', message: 'Presenze salvate con successo' })
    initialSnapshot.value = JSON.stringify(attendanceMap)
  } catch (e) {
    $q.notify({ type: 'negative', message: `Errore: ${e.response?.data?.error || e.message}` })
  } finally {
    saving.value = false
  }
}

const openJustify = (student) => {
  justifyTarget.value = student
  justifyReason.value = 'Motivi di Salute'
  justifyNotes.value = ''
  justifyDialog.value = true
}

const submitJustify = async () => {
  try {
    const att = attendanceMap[justifyTarget.value.id]
    if (att?.attendanceId) {
      await attendanceService.justify(att.attendanceId, {
        student_id: justifyTarget.value.id,
        start_date: dateVal.value,
        end_date: dateVal.value,
        reason: justifyReason.value + (justifyNotes.value ? ` - ${justifyNotes.value}` : '')
      })
      att.isJustified = true
    }
    $q.notify({ type: 'positive', message: 'Giustificazione registrata' })
    justifyDialog.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore nella giustificazione' })
  }
}

defineExpose({
    selectedClassId,
    dateVal,
    students,
    attendanceMap,
    fetchStudentsAndAttendance,
    saveAll,
    openJustify,
    submitJustify
})
</script>
