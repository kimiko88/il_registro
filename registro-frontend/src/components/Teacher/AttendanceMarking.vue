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
        :label="t('udaPage.classLabel')"
        dense outlined
        style="min-width: 220px"
      />
      <q-input v-model="dateVal" type="date" :label="t('gradesPage.date')" dense outlined style="max-width: 160px" />
      <q-select
        v-model="lessonType"
        :options="[
          { label: 'Lezione Standard', value: 'standard', icon: 'school' },
          { label: 'Attività PCTO', value: 'pcto', icon: 'work' },
          { label: 'Orientamento', value: 'orientamento', icon: 'explore' }
        ]"
        emit-value map-options
        label="Tipo Attività"
        dense outlined
        style="min-width: 180px"
      />
      <q-space />
      <q-btn color="primary" icon="save" :label="t('common.save')" :loading="saving" @click="saveAll" :disable="!hasChanges" />
    </div>

    <q-banner v-if="selectedClassId && lessonType !== 'standard'" class="bg-indigo-1 text-indigo-10 rounded-xl q-mb-md border border-indigo-200">
      <template #avatar><q-icon :name="lessonType === 'pcto' ? 'work' : 'explore'" color="indigo" size="24px" /></template>
      <div class="text-weight-bold">
        Attività di {{ lessonType === 'pcto' ? 'PCTO' : 'Orientamento Scolastico' }}
      </div>
    </q-banner>

    <q-banner v-if="!selectedClassId" class="bg-blue-1 q-mb-md">
      <template #avatar><q-icon name="info" color="primary" /></template>
      {{ t('competenciesPage.selectClassPrompt') }}
    </q-banner>

    <q-table
      v-if="selectedClassId"
      :rows="students"
      :columns="columns"
      row-key="id"
      flat bordered
      :loading="loading"
      :pagination="{ rowsPerPage: 50 }"
      hide-bottom
    >
      <template #body="props">
        <q-tr :props="props">
          <q-td key="name" :props="props">
            <div class="row items-center no-wrap">
              <q-avatar size="32px" color="grey-3" text-color="grey-8" class="q-mr-sm">
                {{ props.row.name ? props.row.name.charAt(0) : 'S' }}
              </q-avatar>
              <div>
                <div class="text-weight-bold">{{ props.row.name }}</div>
              </div>
            </div>
          </q-td>

          <q-td key="status" :props="props">
            <q-btn-toggle
              v-model="attendanceMap[props.row.id].status"
              toggle-color="primary"
              dense unelevated
              :options="[
                { label: 'P', value: 'present', slot: 'present' },
                { label: 'A', value: 'absent', slot: 'absent' },
                { label: 'R', value: 'late', slot: 'late' },
                { label: 'U', value: 'early_exit', slot: 'early' }
              ]"
              @update:model-value="onStatusChange(props.row.id)"
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
          <q-input v-model="justifyNotes" :label="t('gradesPage.notes')" type="textarea" outlined dense autogrow />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" label="Giustifica" @click="submitJustify" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, watch, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar, date } from 'quasar'
import { useClassesStore } from '@/stores/classes'
import { attendanceService } from '@/services/attendanceService'

const $q = useQuasar()
const { t } = useI18n()
const classesStore = useClassesStore()

const selectedClassId = ref(null)
const today = date.formatDate(Date.now(), 'YYYY-MM-DD')
const dateVal = ref(today)
const lessonType = ref('standard')

const saving = ref(false)
const loading = ref(false)

const classOptions = ref([])
const students = ref([])
const attendanceMap = reactive({})
const initialSnapshot = ref('')

const justifyDialog = ref(false)
const justifyTarget = ref(null)
const justifyReason = ref('Motivi di Salute')
const justifyNotes = ref('')

const columns = computed(() => [
  { name: 'name', label: t('competenciesPage.student'), field: 'name', align: 'left' },
  { name: 'status', label: t('substitutionsPage.status'), field: 'status', align: 'left' },
  { name: 'minutes', label: t('udaPage.detailTitle'), field: 'minutesLate', align: 'left' },
  { name: 'justified', label: t('verbaliPage.resolutions'), field: 'isJustified', align: 'center' }
])

onMounted(async () => {
  await classesStore.fetchClasses()
  classOptions.value = (classesStore.classes || []).map(c => ({
    id: c.id,
    label: `${c.year}${c.section} ${c.school_level || ''}`.trim()
  }))
  if (classOptions.value.length > 0) {
    selectedClassId.value = classOptions.value[0].id
  }
})

watch([selectedClassId, dateVal], () => {
  if (selectedClassId.value) {
    fetchStudentsAndAttendance()
  }
})

const snapshotCurrent = () => JSON.stringify(attendanceMap)

const hasChanges = computed(() => {
  if (!initialSnapshot.value) return false
  return snapshotCurrent() !== initialSnapshot.value
})

const onStatusChange = (studentId) => {
  const att = attendanceMap[studentId]
  if (att.status !== 'late') att.minutesLate = null
  if (att.status !== 'early_exit') att.exitTime = null
}

const fetchStudentsAndAttendance = async () => {
  loading.value = true
  try {
    const res = await attendanceService.getAttendance(selectedClassId.value, dateVal.value)
    const rawData = res.data || []
    students.value = rawData.map(s => ({ id: s.student_id, name: `${s.first_name} ${s.last_name}` }))
    
    Object.keys(attendanceMap).forEach(k => delete attendanceMap[k])
    rawData.forEach(s => {
      attendanceMap[s.student_id] = {
        attendanceId: s.attendance_id || null,
        status: s.status || 'present',
        minutesLate: s.minutes_late || null,
        exitTime: s.exit_time || null,
        isJustified: Boolean(s.is_justified)
      }
    })
    initialSnapshot.value = snapshotCurrent()
  } catch (e) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    loading.value = false
  }
}

const saveAll = async () => {
  saving.value = true
  try {
    const records = Object.entries(attendanceMap).map(([studentId, data]) => ({
      student_id: studentId,
      status: data.status,
      minutes_late: data.minutesLate,
      exit_time: data.exitTime
    }))
    await attendanceService.recordBulk({
      class_id: selectedClassId.value,
      date: dateVal.value,
      lesson_type: lessonType.value,
      records
    })
    $q.notify({ type: 'positive', message: t('common.success') })
    initialSnapshot.value = snapshotCurrent()
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
    $q.notify({ type: 'positive', message: t('common.success') })
    justifyDialog.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: t('common.error') })
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
