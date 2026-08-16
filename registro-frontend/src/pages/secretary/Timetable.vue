<template>
  <q-page padding class="bg-slate-50 min-h-screen">
    <!-- Top Header -->
    <div class="row items-center justify-between q-mb-lg gap-4">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center gap-2">
          <q-icon name="schedule" color="primary" size="36px" />
          {{ t('timetablePage.title') || 'Orario Scolastico & Cattedre' }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('timetablePage.subtitle') || 'Gestisci l\'orario delle lezioni per classe o consulta l\'orario settimanale dei singoli docenti.' }}
        </p>
      </div>

      <div class="row items-center gap-3">
        <q-btn-toggle
          v-model="viewMode"
          toggle-color="primary"
          flat
          dense
          no-caps
          class="bg-slate-200 rounded-xl q-pa-xs border border-slate-300 shadow-xs"
          :options="[
            { label: t('timetablePage.classSchedule') || 'Orario per Classe', value: 'class', icon: 'groups' },
            { label: t('timetablePage.teacherSchedule') || 'Orario per Docente', value: 'teacher', icon: 'person' }
          ]"
          @update:model-value="onViewModeChange"
        />

        <q-btn
          v-if="viewMode === 'class' && selectedClass"
          :label="isEditing ? (t('timetablePage.viewMode') || 'Vista Lettura') : (t('timetablePage.editSchedule') || 'Modifica Orario')"
          :icon="isEditing ? 'visibility' : 'edit_calendar'"
          :color="isEditing ? 'secondary' : 'primary'"
          unelevated
          no-caps
          class="rounded-xl q-px-md shadow-xs font-bold"
          @click="isEditing = !isEditing"
        />

        <q-btn
          v-if="viewMode === 'teacher' && selectedTeacher"
          :label="isTeacherEditing ? (t('timetablePage.viewMode') || 'Vista Lettura') : (t('timetablePage.editTeacherSchedule') || 'Modifica Orario Docente')"
          :icon="isTeacherEditing ? 'visibility' : 'edit_calendar'"
          :color="isTeacherEditing ? 'secondary' : 'positive'"
          unelevated
          no-caps
          class="rounded-xl q-px-md shadow-xs font-bold"
          @click="isTeacherEditing = !isTeacherEditing"
        />

        <q-btn
          v-if="viewMode === 'class' && selectedClass"
          :label="t('timetablePage.subjects') || 'Cattedre / Materie'"
          icon="menu_book"
          color="indigo-7"
          outline
          no-caps
          class="rounded-xl q-px-md shadow-xs"
          @click="openSubjectsDialog"
        />
      </div>
    </div>

    <!-- Filters Bar -->
    <q-card flat bordered class="rounded-2xl bg-white q-pa-md q-mb-lg shadow-sm">
      <div class="row items-center q-col-gutter-md">
        <!-- Class Selector (Mode = Class) -->
        <div v-if="viewMode === 'class'" class="col-12 col-sm-6 col-md-4">
          <q-select
            v-model="selectedClass"
            :options="classOptions"
            option-value="id"
            option-label="label"
            emit-value map-options
            label="Seleziona Classe *"
            outlined
            dense
            bg-color="white"
            class="rounded-lg"
          >
            <template v-slot:prepend>
              <q-icon name="room" color="primary" />
            </template>
          </q-select>
        </div>

        <!-- Teacher Selector (Mode = Teacher) -->
        <div v-if="viewMode === 'teacher'" class="col-12 col-sm-6 col-md-4">
          <q-select
            v-model="selectedTeacher"
            :options="teacherOptions"
            option-value="id"
            option-label="label"
            emit-value map-options
            label="Seleziona Docente *"
            outlined
            dense
            bg-color="white"
            class="rounded-lg"
          >
            <template v-slot:prepend>
              <q-icon name="person" color="primary" />
            </template>
          </q-select>
        </div>

        <div class="col-auto flex items-center gap-2">
          <q-badge v-if="viewMode === 'class' && currentClassInfo" color="blue-1" text-color="blue-9" class="q-pa-xs px-3 text-caption font-bold rounded-lg border border-blue-200">
            Anno Scolastico: {{ currentClassInfo.academic_year || '2025/2026' }}
          </q-badge>
          <q-badge v-if="viewMode === 'teacher' && selectedTeacher" color="emerald-1" text-color="emerald-9" class="q-pa-xs px-3 text-caption font-bold rounded-lg border border-emerald-300">
            Totale Ore Insegnamento: {{ teacherTotalHours }} ore/settimana
          </q-badge>
        </div>
      </div>
    </q-card>

    <!-- Main Content Area -->
    <div v-if="loading" class="text-center q-pa-xl bg-white rounded-2xl border border-slate-200 shadow-sm">
      <q-spinner-dots color="primary" size="60px" />
      <div class="text-slate-500 q-mt-md">Caricamento orario scolastico...</div>
    </div>

    <!-- Mode: CLASS -->
    <div v-else-if="viewMode === 'class'">
      <div v-if="!selectedClass" class="text-center q-pa-xl bg-white rounded-2xl border border-slate-200 shadow-sm text-slate-500">
        <q-icon name="touch_app" size="72px" class="q-mb-md opacity-30" />
        <div class="text-h6 text-slate-700">Seleziona una Classe</div>
        <div class="text-caption text-slate-400">Scegli una classe dal menu in alto per visualizzare o modificare l'orario delle lezioni.</div>
      </div>

      <div v-else-if="isEditing">
        <q-card flat bordered class="rounded-2xl bg-white q-pa-lg shadow-sm">
          <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md row items-center gap-2">
            <q-icon name="edit_calendar" color="primary" />
            Composizione Orario Settimanale - {{ currentClassInfo?.label || 'Classe' }}
          </div>
          <ScheduleGrid
            :assignments="classAssignments"
            :initial-schedule="scheduleEntries"
            :loading="saving"
            @save="onSaveSchedule"
          />
        </q-card>
      </div>

      <div v-else-if="scheduleEntries.length === 0" class="text-center q-pa-xl bg-white rounded-2xl border border-slate-200 shadow-sm text-slate-500">
        <q-icon name="event_busy" size="72px" class="q-mb-md opacity-30 text-amber-500" />
        <div class="text-h6 text-slate-700">Orario non ancora configurato</div>
        <div class="text-caption text-slate-400 q-mb-lg">Non risulta un orario scolastico salvato per questa classe. Puoi configurarlo ora.</div>
        <q-btn label="Configura Orario Ora" color="primary" icon="edit_calendar" no-caps class="rounded-xl q-px-lg shadow-xs" @click="isEditing = true" />
      </div>

      <!-- Class Timetable Read-Only Table -->
      <q-card v-else flat bordered class="rounded-2xl bg-white overflow-hidden shadow-sm">
        <div class="grid-scroll">
          <table class="timetable-grid">
            <thead>
              <tr class="bg-slate-100 border-b border-slate-300">
                <th class="hour-col py-3 text-center text-slate-700 font-bold text-xs uppercase tracking-wider border-r border-slate-300">Ora</th>
                <th v-for="day in days" :key="day.value" class="day-col py-3 text-center text-slate-700 font-bold text-xs uppercase tracking-wider border-r border-slate-200">
                  {{ day.label }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="hour in 8" :key="hour" class="border-b border-slate-200">
                <td class="hour-cell font-bold text-center bg-slate-50 text-slate-700 border-r border-slate-300 py-2">{{ hour }}ª ora</td>
                <td 
                  v-for="day in 6" 
                  :key="day" 
                  class="schedule-cell p-2 border-r border-slate-200"
                  :class="{ 'has-content': getCell(day, hour) }"
                >
                  <div v-if="getCell(day, hour)" class="cell-content p-2 rounded-xl bg-indigo-50/80 border border-indigo-200 shadow-2xs">
                    <div class="text-subtitle2 text-weight-bold text-indigo-900 leading-tight">{{ getCell(day, hour).subject_name }}</div>
                    <div class="text-caption text-weight-bold text-slate-700 mt-0.5">{{ getCell(day, hour).teacher_name || 'Docente non assegnato' }}</div>
                    <div v-if="getCell(day, hour).room" class="text-caption text-slate-500 mt-0.5">
                      <q-icon name="room" size="xs" class="q-mr-xs" />Aula: {{ getCell(day, hour).room }}
                    </div>
                  </div>
                  <div v-else class="empty-cell text-center text-slate-300 text-caption">-</div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </q-card>
    </div>

    <!-- Mode: TEACHER -->
    <div v-else-if="viewMode === 'teacher'">
      <div v-if="!selectedTeacher" class="text-center q-pa-xl bg-white rounded-2xl border border-slate-200 shadow-sm text-slate-500">
        <q-icon name="badge" size="72px" class="q-mb-md opacity-30" />
        <div class="text-h6 text-slate-700">Seleziona un Docente</div>
        <div class="text-caption text-slate-400">Scegli un docente dal menu in alto per visualizzare o modificare il suo orario completo di insegnamento.</div>
      </div>

      <div v-else-if="isTeacherEditing">
        <q-card flat bordered class="rounded-2xl bg-white q-pa-lg shadow-sm">
          <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md row items-center gap-2">
            <q-icon name="edit_calendar" color="positive" />
            Composizione Orario Docente - {{ teacherOptions.find(t => t.id === selectedTeacher)?.label }}
          </div>
          <TeacherScheduleGrid
            :classes="classOptions"
            :subjects="subjectOptions"
            :initial-schedule="teacherScheduleEntries"
            :loading="saving"
            @save="onSaveTeacherSchedule"
          />
        </q-card>
      </div>

      <div v-else-if="teacherScheduleEntries.length === 0" class="text-center q-pa-xl bg-white rounded-2xl border border-slate-200 shadow-sm text-slate-500">
        <q-icon name="event_busy" size="72px" class="q-mb-md opacity-30 text-amber-500" />
        <div class="text-h6 text-slate-700">Nessuna lezione a orario</div>
        <div class="text-caption text-slate-400 q-mb-lg">Non risultano lezioni assegnate a questo docente negli orari delle classi. Puoi configurarle ora.</div>
        <q-btn label="Configura Orario Docente Ora" color="positive" icon="edit_calendar" no-caps class="rounded-xl q-px-lg shadow-xs font-bold" @click="isTeacherEditing = true" />
      </div>

      <!-- Teacher Timetable Read-Only Table -->
      <q-card v-else flat bordered class="rounded-2xl bg-white overflow-hidden shadow-sm">
        <div class="grid-scroll">
          <table class="timetable-grid">
            <thead>
              <tr class="bg-slate-100 border-b border-slate-300">
                <th class="hour-col py-3 text-center text-slate-700 font-bold text-xs uppercase tracking-wider border-r border-slate-300">Ora</th>
                <th v-for="day in days" :key="day.value" class="day-col py-3 text-center text-slate-700 font-bold text-xs uppercase tracking-wider border-r border-slate-200">
                  {{ day.label }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="hour in 8" :key="hour" class="border-b border-slate-200">
                <td class="hour-cell font-bold text-center bg-slate-50 text-slate-700 border-r border-slate-300 py-2">{{ hour }}ª ora</td>
                <td 
                  v-for="day in 6" 
                  :key="day" 
                  class="schedule-cell p-2 border-r border-slate-200"
                  :class="{ 'has-content': getTeacherCell(day, hour) }"
                >
                  <div v-if="getTeacherCell(day, hour)" class="cell-content p-2 rounded-xl bg-emerald-50/80 border border-emerald-200 shadow-2xs">
                    <div class="text-subtitle2 text-weight-bold text-emerald-900 leading-tight">{{ getTeacherCell(day, hour).subject_name }}</div>
                    <div class="text-caption text-weight-bold text-slate-700 mt-0.5">Classe: {{ getTeacherCellClassName(getTeacherCell(day, hour)) }}</div>
                    <div v-if="getTeacherCell(day, hour).room" class="text-caption text-slate-500 mt-0.5">
                      <q-icon name="room" size="xs" class="q-mr-xs" />Aula: {{ getTeacherCell(day, hour).room }}
                    </div>
                  </div>
                  <div v-else class="empty-cell text-center text-slate-300 text-caption">-</div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </q-card>
    </div>

    <!-- Modal Gestione Cattedre / Materie della classe -->
    <q-dialog v-model="showSubjectsDialog">
      <q-card style="width: min(900px, 95vw); max-height: 90vh;" class="rounded-2xl overflow-hidden shadow-24 bg-white column">
        <q-card-section class="bg-indigo-7 text-white row items-center justify-between q-pa-lg col-auto">
          <div class="text-h6 text-weight-bold">Cattedre &amp; Materie - {{ currentClassInfo?.label }}</div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg scroll col" style="flex: 1; overflow-y: auto;">
          <div class="row q-col-gutter-lg">
            <!-- Left: Existing Assignments -->
            <div class="col-12 col-md-7">
              <div class="text-subtitle2 font-bold text-slate-700 q-mb-sm">Materie e Docenti Assegnati</div>
              <q-table
                :rows="classAssignments"
                :columns="assignmentColumns"
                row-key="id"
                flat
                bordered
                dense
                class="rounded-xl"
              >
                <template #body-cell-actions="props">
                  <q-td :props="props" auto-width>
                    <q-btn flat round dense color="negative" icon="delete" size="sm" @click="removeAssignment(props.row.id)" />
                  </q-td>
                </template>
              </q-table>
            </div>

            <!-- Right: Add New Assignment -->
            <div class="col-12 col-md-5">
              <q-card flat class="rounded-xl bg-slate-50 q-pa-md border border-slate-200">
                <div class="text-subtitle2 font-bold text-slate-800 q-mb-md">Assegna Cattedra a Classe</div>
                <q-form @submit="addAssignment" class="q-gutter-y-md">
                  <q-select
                    v-model="assignForm.subject_id"
                    :options="subjectOptions"
                    label="Materia *"
                    outlined
                    dense
                    emit-value
                    map-options
                    :rules="[val => !!val || 'Seleziona materia']"
                  />
                  <q-select
                    v-model="assignForm.teacher_id"
                    :options="teacherOptions"
                    label="Docente"
                    outlined
                    dense
                    emit-value
                    map-options
                  />
                  <q-input
                    v-model.number="assignForm.hours_per_week"
                    label="Ore Settimanali"
                    type="number"
                    outlined
                    dense
                    min="1"
                  />
                  <q-btn type="submit" label="Assegna Cattedra" color="primary" class="full-width rounded-lg q-py-sm shadow-xs" no-caps />
                </q-form>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch, computed, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import adminService from '@/services/adminService'
import ScheduleGrid from '@/components/Secretary/ScheduleGrid.vue'
import TeacherScheduleGrid from '@/components/Secretary/TeacherScheduleGrid.vue'

const $q = useQuasar()
const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)
const isEditing = ref(false)
const isTeacherEditing = ref(false)
const viewMode = ref('class')
const selectedClass = ref(null)
const selectedTeacher = ref(null)

const classOptions = ref([])
const teacherOptions = ref([])
const subjectOptions = ref([])

const scheduleEntries = ref([])
const teacherScheduleEntries = ref([])
const classAssignments = ref([])
const showSubjectsDialog = ref(false)

const assignForm = reactive({
  subject_id: null,
  teacher_id: null,
  hours_per_week: 2
})

const days = computed(() => [
  { label: t('timetablePage.monday') || 'Lunedì', value: 1 },
  { label: t('timetablePage.tuesday') || 'Martedì', value: 2 },
  { label: t('timetablePage.wednesday') || 'Mercoledì', value: 3 },
  { label: t('timetablePage.thursday') || 'Giovedì', value: 4 },
  { label: t('timetablePage.friday') || 'Venerdì', value: 5 },
  { label: t('timetablePage.saturday') || 'Sabato', value: 6 }
])

const assignmentColumns = computed(() => [
  { name: 'subject_name', label: t('agendaPage.subject') || 'Materia', field: 'subject_name', align: 'left' },
  { name: 'teacher_name', label: t('agendaPage.teacher') || 'Docente', field: 'teacher_name', align: 'left' },
  { name: 'hours_per_week', label: t('timetablePage.hoursPerWeek') || 'Ore/Sett.', field: 'hours_per_week', align: 'center' },
  { name: 'actions', label: '', field: 'actions', align: 'right' }
])

const currentClassInfo = computed(() => {
  return classOptions.value.find(c => c.id === selectedClass.value)
})

const teacherTotalHours = computed(() => {
  return teacherScheduleEntries.value.length
})

onMounted(async () => {
  loading.value = true
  await Promise.all([
    fetchClasses(),
    fetchTeachers(),
    fetchSubjects()
  ])
  if (classOptions.value.length > 0) {
    selectedClass.value = classOptions.value[0].id
  }
  loading.value = false
})

const fetchClasses = async () => {
  try {
    const res = await adminService.getClasses()
    const raw = res.data || []
    classOptions.value = raw.map(c => ({
      id: c.id,
      label: `Classe ${c.name}${c.section} (${c.academic_year || '2025/2026'})`,
      name: c.name,
      section: c.section,
      academic_year: c.academic_year
    }))
  } catch (e) {
    console.error('Failed fetching classes', e)
  }
}

const fetchTeachers = async () => {
  try {
    const res = await adminService.getTeachersList()
    const raw = res.data || []
    teacherOptions.value = raw.map(t => ({
      id: t.id || t.user_id,
      user_id: t.user_id,
      label: `${t.first_name || ''} ${t.last_name || ''}`.trim() || t.email
    }))
  } catch (e) {
    console.error('Failed fetching teachers', e)
  }
}

const fetchSubjects = async () => {
  try {
    const res = await adminService.getSubjects()
    const raw = res.data || []
    subjectOptions.value = raw.map(s => ({
      value: s.id,
      label: s.name
    }))
  } catch (e) {
    console.error('Failed fetching subjects', e)
  }
}

const onViewModeChange = async (val) => {
  loading.value = true
  isEditing.value = false
  isTeacherEditing.value = false
  if (val === 'class' && selectedClass.value) {
    await Promise.all([fetchClassSchedule(), fetchClassAssignments()])
  } else if (val === 'teacher' && selectedTeacher.value) {
    await fetchTeacherSchedule()
  }
  loading.value = false
}

watch(selectedClass, async (newVal) => {
  if (newVal && viewMode.value === 'class') {
    loading.value = true
    await Promise.all([fetchClassSchedule(), fetchClassAssignments()])
    loading.value = false
  }
})

watch(selectedTeacher, async (newVal) => {
  if (newVal && viewMode.value === 'teacher') {
    loading.value = true
    await fetchTeacherSchedule()
    loading.value = false
  }
})

const fetchClassSchedule = async () => {
  try {
    const res = await adminService.getClassSchedule(selectedClass.value)
    scheduleEntries.value = res.data || []
  } catch (e) {
    console.error(e)
    scheduleEntries.value = []
  }
}

const fetchClassAssignments = async () => {
  try {
    const res = await adminService.getClassSubjects(selectedClass.value)
    classAssignments.value = res.data || []
  } catch (e) {
    console.error(e)
    classAssignments.value = []
  }
}

const fetchTeacherSchedule = async () => {
  if (!selectedTeacher.value) return
  try {
    const teacherObj = teacherOptions.value.find(t => t.id === selectedTeacher.value)
    const targetId = teacherObj?.user_id || selectedTeacher.value
    const res = await adminService.getTeacherSchedule(targetId).catch(() => null)
    if (res && res.data) {
      teacherScheduleEntries.value = res.data
    } else {
      teacherScheduleEntries.value = []
    }
  } catch (e) {
    console.error('Failed fetching teacher schedule', e)
    teacherScheduleEntries.value = []
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
    $q.notify({ type: 'positive', message: 'Orario scolastico salvato con successo!' })
    await fetchClassSchedule()
    isEditing.value = false
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante il salvataggio dell\'orario' })
  } finally {
    saving.value = false
  }
}

const onSaveTeacherSchedule = async (entries) => {
  if (!selectedTeacher.value) return
  saving.value = true
  try {
    const teacherObj = teacherOptions.value.find(t => t.id === selectedTeacher.value)
    const targetId = teacherObj?.user_id || selectedTeacher.value
    const formattedEntries = entries.map(e => ({
      day_of_week: e.day_of_week,
      hour_index: e.hour_index,
      class_id: e.class_id,
      subject_id: e.subject_id,
      room: e.room || ''
    }))
    await adminService.saveTeacherSchedule(targetId, { entries: formattedEntries })
    $q.notify({ type: 'positive', message: 'Orario docente salvato con successo!' })
    await fetchTeacherSchedule()
    isTeacherEditing.value = false
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante il salvataggio dell\'orario docente' })
  } finally {
    saving.value = false
  }
}

const getCell = (day, hour) => {
  return scheduleEntries.value.find(e => e.day_of_week === day && e.hour_index === hour)
}

const getTeacherCell = (day, hour) => {
  return teacherScheduleEntries.value.find(e => e.day_of_week === day && e.hour_index === hour)
}

const getTeacherCellClassName = (cell) => {
  if (!cell) return ''
  if (cell.class_name) return cell.class_name
  const found = classOptions.value.find(c => c.id === cell.class_id)
  return found ? found.label : (cell.class_id ? `Classe ${cell.class_id.substring(0, 5)}` : '')
}

const openSubjectsDialog = () => {
  showSubjectsDialog.value = true
}

const addAssignment = async () => {
  if (!assignForm.subject_id || !selectedClass.value) return
  try {
    await adminService.assignSubjectToClass(selectedClass.value, assignForm)
    $q.notify({ type: 'positive', message: 'Cattedra assegnata correttamente!' })
    await fetchClassAssignments()
    assignForm.subject_id = null
    assignForm.teacher_id = null
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante l\'assegnazione della cattedra' })
  }
}

const removeAssignment = async (assignmentId) => {
  try {
    await adminService.removeClassSubject(selectedClass.value, assignmentId)
    $q.notify({ type: 'positive', message: 'Cattedra rimossa' })
    await fetchClassAssignments()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante la rimozione della cattedra' })
  }
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
  min-width: 750px;
}

.timetable-grid th, .timetable-grid td {
  vertical-align: middle;
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
  width: 80px;
  min-width: 80px;
}

.schedule-cell {
  height: 80px;
}

.has-content {
  background-color: rgba(248, 250, 252, 0.6);
}
</style>
