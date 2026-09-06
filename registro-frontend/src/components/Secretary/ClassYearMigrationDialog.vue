<template>
  <q-dialog :model-value="modelValue" persistent class="premium-dialog" @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(1000px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
      <!-- Header -->
      <q-card-section class="bg-indigo-9 text-white row items-center q-pa-md shrink-0">
        <div class="row items-center">
          <q-avatar color="white-20" text-color="white" icon="published_with_changes" class="q-mr-sm" size="36px" />
          <div>
            <div class="text-h6 text-weight-bold">
              {{ $t('secretaryClasses.migrationTitle') || "Migrazione Studenti e Passaggio d'Anno" }}
            </div>
            <div class="text-subtitle2 text-white/90">
              {{ $t('secretaryClasses.migrationSubtitle') || "Gestisci avanzamento classi, promossi, bocciati e diplomati" }}
            </div>
          </div>
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
      </q-card-section>

      <!-- Body / Stepper -->
      <q-card-section class="q-pa-md col overflow-y-auto">
        <!-- Step 1: Configuration -->
        <div v-if="migrationStep === 1" class="q-gutter-y-md">
          <div class="bg-indigo-50 border border-indigo-100 rounded-xl q-pa-md text-indigo-9">
            <div class="text-weight-bold flex items-center gap-2">
              <q-icon name="info" size="sm" /> {{ $t('secretaryClasses.migrationProcedure') || "Procedura di Passaggio d'Anno" }}
            </div>
            <div class="text-caption q-mt-xs">
              {{ $t('secretaryClasses.migrationProcedureDesc') || "Seleziona l'Anno Scolastico sorgente di origine (es. 2024/2025) e quello di destinazione per il nuovo anno. Potrai configurare rapidamente per ciascuno studente lo stato (Promosso, Bocciato, Diplomato o Trasferito)." }}
            </div>
          </div>

          <div class="row q-col-gutter-md q-mt-sm">
            <div class="col-12 col-md-6">
              <q-select
                v-model="migrationSourceYear"
                :options="academicYearOptions"
                :label="$t('secretaryClasses.sourceYear') || 'Anno Sorgente (Origine)'"
                outlined dense
              />
            </div>
            <div class="col-12 col-md-6">
              <q-select
                v-model="migrationTargetYear"
                :options="academicYearOptions"
                :label="$t('secretaryClasses.targetYear') || 'Anno Destinazione (Nuovo Anno)'"
                outlined dense
              />
            </div>
            <div class="col-12">
              <q-select
                v-model="migrationSourceClassId"
                :options="migrationClassOptions"
                :label="$t('secretaryClasses.filterByClass') || 'Filtra per Classe (oppure Tutte le Classi)'"
                outlined dense
                emit-value map-options
              />
            </div>
          </div>
        </div>

        <!-- Step 2: Student Outcomes & Classes Setup -->
        <div v-else-if="migrationStep === 2" class="q-gutter-y-md">
          <!-- Group Actions Toolbar -->
          <div class="row items-center justify-between bg-slate-50 border border-slate-200 rounded-xl q-pa-sm">
            <div class="text-subtitle2 text-weight-bold text-slate-700">
              {{ $t('secretaryClasses.studentsToProcess') || 'Studenti da Elaborare' }}: {{ migrationStudents.length }}
            </div>
            <div class="row q-gutter-xs">
              <q-btn size="xs" color="positive" icon="done_all" :label="$t('secretaryClasses.markAllPromoted') || 'Segna Tutti Promossi'" no-caps unelevated @click="setAllMigrationAction('promoted')" />
              <q-btn size="xs" color="purple" icon="school" :label="$t('secretaryClasses.smartDefaults') || 'Diploma 5° / Promuovi 1-4°'" no-caps unelevated @click="setSmartMigrationDefaults()" />
              <q-btn size="xs" color="negative" icon="block" :label="$t('secretaryClasses.markAllRepeater') || 'Segna Tutti Bocciati'" no-caps unelevated @click="setAllMigrationAction('repeater')" />
            </div>
          </div>

          <!-- Loading dots if fetching -->
          <div v-if="migrationLoading" class="text-center q-pa-xl">
            <q-spinner-dots color="indigo-7" size="40px" />
            <div class="text-caption text-slate-500 q-mt-sm">
              {{ $t('secretaryClasses.loadingMigration') || "Caricamento studenti e predisposizione classi dell'anno destinazione..." }}
            </div>
          </div>

          <!-- Students List with Actions -->
          <div v-else-if="migrationStudents.length === 0" class="text-center q-pa-xl text-slate-400 border rounded-xl bg-slate-50">
            {{ $t('secretaryClasses.noStudentsFound') || 'Nessuno studente trovato per i criteri selezionati.' }}
          </div>

          <q-scroll-area v-else style="height: 420px;" class="rounded-xl border border-slate-200 bg-white">
            <q-list separator dense>
              <q-item v-for="st in migrationStudents" :key="st.student_id" class="q-py-sm items-center justify-between">
                <!-- Student Info -->
                <q-item-section style="width: 250px" class="shrink-0">
                  <q-item-label class="text-weight-bold text-slate-800 ellipsis">
                    {{ st.last_name }} {{ st.first_name || st.name }}
                  </q-item-label>
                  <q-item-label caption class="text-slate-400 text-xs ellipsis">
                    {{ $t('secretaryClasses.currentClassLabel') || 'Classe Attuale' }}: <q-badge color="cyan-8" :label="st.current_class_name" class="q-ml-xs" />
                  </q-item-label>
                </q-item-section>

                <!-- Action Buttons -->
                <q-item-section class="col q-px-sm">
                  <q-btn-toggle
                    v-model="st.action"
                    dense
                    toggle-color="indigo-7"
                    size="xs"
                    no-caps
                    spread
                    :options="[
                      { label: '🟢 ' + ($t('secretaryClasses.actionPromoted') || 'Promosso/a'), value: 'promoted' },
                      { label: '🔴 ' + ($t('secretaryClasses.actionRepeater') || 'Bocciato/a'), value: 'repeater' },
                      { label: '🎓 ' + ($t('secretaryClasses.actionGraduated') || 'Diplomato/a'), value: 'graduated' },
                      { label: '🚪 ' + ($t('secretaryClasses.actionLeft') || 'Trasferito/a'), value: 'left' }
                    ]"
                    @update:model-value="onStudentActionChange(st)"
                  />
                </q-item-section>

                <!-- Target Class Selector -->
                <q-item-section style="width: 220px" class="shrink-0 q-pl-sm">
                  <q-select
                    v-if="st.action === 'promoted' || st.action === 'repeater'"
                    v-model="st.target_class_id"
                    :options="getTargetClassOptions()"
                    :label="$t('secretaryClasses.targetClassLabel') || 'Classe Destinazione'"
                    outlined dense emit-value map-options
                    style="font-size: 11px;"
                  />
                  <div v-else class="text-caption text-slate-400 text-italic text-center">
                    {{ $t('secretaryClasses.noClassUnassigned') || 'Nessuna classe (Disassociato)' }}
                  </div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-scroll-area>
        </div>

        <!-- Step 3: Confirmation Summary -->
        <div v-else-if="migrationStep === 3" class="q-gutter-y-md">
          <div class="text-subtitle1 text-weight-bold text-slate-800">
            {{ $t('secretaryClasses.migrationSummaryTitle') || 'Riepilogo Migrazione' }} ({{ migrationSourceYear }} → {{ migrationTargetYear }})
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-6 col-md-3">
              <q-card flat class="bg-emerald-50 border border-emerald-200 text-emerald-9 q-pa-md text-center rounded-xl">
                <div class="text-h4 text-weight-bold">{{ migrationSummary.promoted }}</div>
                <div class="text-caption text-weight-medium">{{ $t('secretaryClasses.promotedCount') || 'Promossi' }}</div>
              </q-card>
            </div>
            <div class="col-6 col-md-3">
              <q-card flat class="bg-rose-50 border border-rose-200 text-rose-9 q-pa-md text-center rounded-xl">
                <div class="text-h4 text-weight-bold">{{ migrationSummary.repeater }}</div>
                <div class="text-caption text-weight-medium">{{ $t('secretaryClasses.repeaterCount') || 'Bocciati' }}</div>
              </q-card>
            </div>
            <div class="col-6 col-md-3">
              <q-card flat class="bg-purple-50 border border-purple-200 text-purple-9 q-pa-md text-center rounded-xl">
                <div class="text-h4 text-weight-bold">{{ migrationSummary.graduated }}</div>
                <div class="text-caption text-weight-medium">{{ $t('secretaryClasses.graduatedCount') || 'Diplomati' }}</div>
              </q-card>
            </div>
            <div class="col-6 col-md-3">
              <q-card flat class="bg-slate-100 border border-slate-200 text-slate-7 q-pa-md text-center rounded-xl">
                <div class="text-h4 text-weight-bold">{{ migrationSummary.left }}</div>
                <div class="text-caption text-weight-medium">{{ $t('secretaryClasses.leftCount') || 'Trasferiti / Ritirati' }}</div>
              </q-card>
            </div>
          </div>
        </div>
      </q-card-section>

      <!-- Footer Actions -->
      <q-card-actions align="right" class="q-pa-md bg-slate-50 border-t border-slate-200 shrink-0">
        <q-btn v-if="migrationStep > 1" flat :label="$t('common.back') || 'Indietro'" color="slate-600" :disabled="migrationLoading" @click="migrationStep--" />
        <q-space />
        <q-btn flat :label="$t('common.cancel') || 'Annulla'" v-close-popup color="slate-500" />
        <q-btn
          v-if="migrationStep === 1"
          color="indigo-7"
          :label="$t('secretaryClasses.btnConfigureStudents') || 'Avanti: Configura Studenti'"
          icon-right="arrow_forward"
          no-caps class="rounded-lg q-px-md"
          @click="goToStep2"
        />
        <q-btn
          v-else-if="migrationStep === 2"
          color="indigo-7"
          :label="$t('secretaryClasses.btnVerifySummary') || 'Avanti: Verifica Riepilogo'"
          icon-right="arrow_forward"
          no-caps class="rounded-lg q-px-md"
          :disabled="migrationStudents.length === 0"
          @click="migrationStep = 3"
        />
        <q-btn
          v-else-if="migrationStep === 3"
          color="emerald-7"
          :label="$t('secretaryClasses.btnExecuteMigration') || 'Conferma ed Esegui Migrazione'"
          icon="check_circle"
          no-caps class="rounded-lg q-px-lg shadow-sm"
          :loading="migrationLoading"
          @click="executeMigration"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import adminService from '@/services/adminService'
import api from '@/services/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  academicYearOptions: {
    type: Array,
    default: () => []
  },
  currentYearStr: {
    type: String,
    default: ''
  },
  classes: {
    type: Array,
    default: () => []
  },
  schoolId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'migrated'])

const $q = useQuasar()
const { t } = useI18n()

const migrationStep = ref(1)
const migrationLoading = ref(false)
const migrationSourceYear = ref(props.currentYearStr)
const migrationTargetYear = ref('')
const migrationSourceClassId = ref(null)

const targetYearClasses = ref([])
const migrationStudents = ref([])

// When opened, reset wizard state
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    migrationStep.value = 1
    migrationSourceYear.value = props.currentYearStr
    migrationSourceClassId.value = null
    migrationStudents.value = []
    targetYearClasses.value = []

    if (props.academicYearOptions && props.academicYearOptions.length > 0) {
      const idx = props.academicYearOptions.indexOf(props.currentYearStr)
      if (idx !== -1 && idx < props.academicYearOptions.length - 1) {
        migrationTargetYear.value = props.academicYearOptions[idx + 1]
      } else {
        const parts = (props.currentYearStr || '').split('/')
        if (parts.length === 2) {
          const start = parseInt(parts[0], 10)
          migrationTargetYear.value = `${start + 1}/${start + 2}`
        } else {
          migrationTargetYear.value = props.academicYearOptions[0] || ''
        }
      }
    }
  }
}, { immediate: true })

const migrationClassOptions = computed(() => {
  return [
    { label: t('secretaryClasses.allSourceClasses') || "Tutte le classi dell'anno sorgente", value: null },
    ...props.classes.map(c => ({
      label: `${c.name}${c.section} (${c.academic_year})`,
      value: c.id
    }))
  ]
})

const sortAlphabetically = (arr) => {
  return [...arr].sort((a, b) => {
    const nameA = `${a.last_name || ''} ${a.first_name || ''}`.trim()
    const nameB = `${b.last_name || ''} ${b.first_name || ''}`.trim()
    return nameA.localeCompare(nameB)
  })
}

const autoCreateMissingTargetClasses = async () => {
  const existingNames = new Set(targetYearClasses.value.map(c => `${c.name}${c.section}`))
  const sourceClasses = migrationSourceClassId.value
    ? props.classes.filter(c => c.id === migrationSourceClassId.value)
    : props.classes

  for (const sc of sourceClasses) {
    const section = sc.section || ''
    const rawName = sc.name || ''
    const match = rawName.match(/^(\d+)(.*)$/)
    const gradeNum = match ? parseInt(match[1], 10) : null
    const restName = match ? match[2] : rawName

    // Repeater class target
    const repeaterTargetName = `${sc.name}${section}`
    if (!existingNames.has(repeaterTargetName)) {
      try {
        const created = await adminService.createClass({
          name: sc.name,
          section: sc.section,
          articolazione: sc.articolazione || '',
          academic_year: migrationTargetYear.value
        })
        targetYearClasses.value.push(created.data)
        existingNames.add(repeaterTargetName)
      } catch (e) {
        console.error('Failed auto-creating repeater target class', e)
      }
    }

    // Promoted class target
    if (gradeNum && gradeNum < 5) {
      const nextGradeName = `${gradeNum + 1}${restName}`
      const promotedTargetName = `${nextGradeName}${section}`
      if (!existingNames.has(promotedTargetName)) {
        try {
          const created = await adminService.createClass({
            name: nextGradeName,
            section: sc.section,
            articolazione: sc.articolazione || '',
            academic_year: migrationTargetYear.value
          })
          targetYearClasses.value.push(created.data)
          existingNames.add(promotedTargetName)
        } catch (e) {
          console.error('Failed auto-creating promoted target class', e)
        }
      }
    }
  }
}

const goToStep2 = async () => {
  migrationLoading.value = true
  migrationStep.value = 2
  try {
    const targetClassesRes = await adminService.getSchoolClasses(props.schoolId, migrationTargetYear.value)
    targetYearClasses.value = targetClassesRes.data || []

    await autoCreateMissingTargetClasses()

    let sourceClassesToProcess = []
    if (migrationSourceClassId.value) {
      sourceClassesToProcess = props.classes.filter(c => c.id === migrationSourceClassId.value)
    } else {
      sourceClassesToProcess = props.classes.filter(c => c.academic_year === migrationSourceYear.value)
    }

    const studentPromises = sourceClassesToProcess.map(c =>
      api.get('/users', { params: { role: 'student', class_id: c.id, page_size: 500 } }).then(res => ({
        class: c,
        students: res.data?.users || res.data || []
      }))
    )

    const classResults = await Promise.all(studentPromises)
    const items = []

    for (const cr of classResults) {
      const cName = `${cr.class.name}${cr.class.section}`
      const rawName = cr.class.name || ''
      const match = rawName.match(/^(\d+)(.*)$/)
      const gradeNum = match ? parseInt(match[1], 10) : null
      const restName = match ? match[2] : rawName

      for (const st of cr.students) {
        let action = 'promoted'
        if (gradeNum >= 5) {
          action = 'graduated'
        }

        let targetClassId = null
        if (action === 'promoted' && gradeNum && gradeNum < 5) {
          const nextName = `${gradeNum + 1}${restName}`
          const foundTarget = targetYearClasses.value.find(tc => tc.name === nextName && tc.section === cr.class.section)
          if (foundTarget) targetClassId = foundTarget.id
        }

        items.push({
          student_id: st.id,
          first_name: st.first_name || st.name || '',
          last_name: st.last_name || '',
          email: st.email || '',
          current_class_id: cr.class.id,
          current_class_name: cName,
          current_grade: gradeNum,
          current_section: cr.class.section,
          action: action,
          target_class_id: targetClassId
        })
      }
    }

    migrationStudents.value = sortAlphabetically(items)
  } catch (e) {
    console.error('Error setting up migration step 2', e)
    $q.notify({ type: 'negative', message: t('secretaryClasses.migrationLoadError') || 'Errore nel caricamento dati per la migrazione' })
  } finally {
    migrationLoading.value = false
  }
}

const onStudentActionChange = (st) => {
  if (st.action === 'promoted') {
    if (st.current_grade && st.current_grade < 5) {
      const nextName = `${st.current_grade + 1}`
      const foundTarget = targetYearClasses.value.find(tc => tc.name.startsWith(nextName) && tc.section === st.current_section)
      if (foundTarget) st.target_class_id = foundTarget.id
    } else {
      st.target_class_id = null
    }
  } else if (st.action === 'repeater') {
    const foundTarget = targetYearClasses.value.find(tc => tc.name.startsWith(`${st.current_grade}`) && tc.section === st.current_section)
    if (foundTarget) st.target_class_id = foundTarget.id
  } else {
    st.target_class_id = null
  }
}

const setAllMigrationAction = (action) => {
  migrationStudents.value.forEach(st => {
    st.action = action
    onStudentActionChange(st)
  })
}

const setSmartMigrationDefaults = () => {
  migrationStudents.value.forEach(st => {
    if (st.current_grade >= 5) {
      st.action = 'graduated'
    } else {
      st.action = 'promoted'
    }
    onStudentActionChange(st)
  })
}

const getTargetClassOptions = () => {
  return targetYearClasses.value.map(c => ({
    label: `${c.name}${c.section}${c.articolazione ? ' (' + c.articolazione + ')' : ''}`,
    value: c.id
  }))
}

const migrationSummary = computed(() => {
  const summary = { promoted: 0, repeater: 0, graduated: 0, left: 0 }
  migrationStudents.value.forEach(st => {
    if (summary[st.action] !== undefined) {
      summary[st.action]++
    }
  })
  return summary
})

const executeMigration = async () => {
  migrationLoading.value = true
  try {
    const payload = {
      source_academic_year: migrationSourceYear.value,
      target_academic_year: migrationTargetYear.value,
      migrations: migrationStudents.value.map(st => ({
        student_id: st.student_id,
        action: st.action,
        target_class_id: (st.action === 'promoted' || st.action === 'repeater') ? st.target_class_id : null
      }))
    }

    await api.post('/classes/migrate-students', payload)

    $q.notify({
      type: 'positive',
      message: t('secretaryClasses.migrationSuccess') || 'Migrazione anno scolastico completata con successo!',
      caption: `${migrationSummary.value.promoted} ${t('secretaryClasses.promotedCount') || 'Promossi'}, ${migrationSummary.value.repeater} ${t('secretaryClasses.repeaterCount') || 'Bocciati'}, ${migrationSummary.value.graduated} ${t('secretaryClasses.graduatedCount') || 'Diplomati'}`
    })

    emit('update:modelValue', false)
    emit('migrated')
  } catch (e) {
    console.error('Migration failed', e)
    $q.notify({
      type: 'negative',
      message: t('secretaryClasses.migrationError') || 'Errore durante l\'esecuzione della migrazione',
      caption: e.response?.data?.error || e.message
    })
  } finally {
    migrationLoading.value = false
  }
}
</script>

<style scoped>
.rounded-xl { border-radius: 1rem; }
.rounded-lg { border-radius: 0.75rem; }
.border-slate-200 { border: 1px solid #e2e8f0; }
.bg-slate-50 { background-color: #f8fafc; }
</style>
