<template>
  <q-card class="glass-card shadow-soft rounded-2xl overflow-hidden absence-limit-card" role="region" :aria-label="t('parent.absenceLimit.ariaLabel') || 'Riepilogo Assenze e Monitoraggio Limite 25% DPR 122/2009'">
    <!-- Header -->
    <q-card-section class="row items-center justify-between q-pb-none">
      <div>
        <div class="text-h6 text-weight-bold text-slate-800 row items-center">
          <q-icon name="schedule" :color="statusColor" class="q-mr-sm" size="26px" />
          <span>{{ t('parent.absenceLimit.title') || 'Limite Assenze 25% (D.P.R. 122/2009)' }}</span>
        </div>
        <div class="text-caption text-slate-500 q-mt-xs">
          {{ t('parent.absenceLimit.subtitle') || 'Monitoraggio soglia validità anno scolastico (min. 75% frequenza)' }}
        </div>
      </div>

      <div class="row items-center q-gutter-xs">
        <q-badge
          :color="statusColor"
          class="q-px-sm q-py-xs text-caption text-weight-bold rounded-lg"
        >
          <q-icon :name="statusIcon" class="q-mr-xs" />
          {{ statusLabel }}
        </q-badge>

        <!-- Annual Hours Preset Menu -->
        <q-btn flat round dense icon="tune" color="slate-600" :aria-label="t('parent.absenceLimit.configHours') || 'Configura Monte Ore Annuo'">
          <q-tooltip>{{ t('parent.absenceLimit.configHoursTooltip') || 'Configura monte ore annuo della scuola' }}</q-tooltip>
          <q-menu auto-close>
            <q-list dense style="min-width: 180px">
              <q-item-label header class="text-caption text-weight-bold">Monte Ore Annuale</q-item-label>
              <q-item
                v-for="hours in annualHoursPresets"
                :key="hours.value"
                clickable
                :active="currentAnnualHours === hours.value"
                active-class="bg-blue-1 text-primary text-weight-bold"
                @click="currentAnnualHours = hours.value"
              >
                <q-item-section>{{ hours.label }} ({{ hours.value }}h)</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </div>
    </q-card-section>

    <!-- Loading State -->
    <q-card-section v-if="loading" class="q-py-lg">
      <div class="row q-col-gutter-md">
        <div class="col-12"><q-skeleton type="rect" height="18px" class="rounded-borders" /></div>
        <div v-for="n in 4" :key="n" class="col-6 col-md-3">
          <q-skeleton type="rect" height="60px" class="rounded-borders" />
        </div>
      </div>
    </q-card-section>

    <!-- Main Content -->
    <q-card-section v-else class="q-pt-md">
      <!-- Progress Bar with Critical 25% Marker -->
      <div class="q-mb-md">
        <div class="row items-center justify-between text-caption text-slate-600 q-mb-xs">
          <span class="text-weight-medium">
            {{ t('parent.absenceLimit.usedPercentage') || 'Assenze effettuate' }}:
            <strong :class="`text-${statusColor}`">{{ absencePercentOfTotal.toFixed(1) }}%</strong>
            <span class="text-slate-400"> ({{ totalAbsenceHours.toFixed(1) }}h su {{ currentAnnualHours }}h)</span>
          </span>
          <span class="text-weight-bold" :class="remainingHours <= 20 ? 'text-negative' : 'text-slate-700'">
            {{ t('parent.absenceLimit.remainingHours') || 'Ore residue consentite' }}: {{ Math.max(0, remainingHours).toFixed(1) }}h
          </span>
        </div>

        <!-- Progress bar container with 25% threshold pin -->
        <div class="relative-position q-py-xs">
          <q-linear-progress
            :value="progressRatio"
            :color="statusColor"
            track-color="slate-200"
            size="14px"
            rounded
            class="absence-progress-bar"
          />
          <!-- 25% Limit Marker Line -->
          <div
            class="limit-marker-pin absolute"
            style="left: 25%; top: 0; bottom: 0;"
            :title="t('parent.absenceLimit.legalLimitTooltip') || 'Soglia limite di legge: 25% assenze max'"
          >
            <div class="marker-line bg-negative"></div>
            <div class="marker-label text-negative text-weight-bolder">25%</div>
          </div>
        </div>

        <div class="row items-center justify-between text-caption text-slate-400 q-mt-xs">
          <span>0h (0%)</span>
          <span class="text-negative text-weight-bold">
            <q-icon name="warning" size="xs" />
            Max {{ maxAllowedAbsenceHours }}h (25%)
          </span>
          <span>{{ currentAnnualHours }}h (100%)</span>
        </div>
      </div>

      <!-- 4 Summary KPI Tiles -->
      <div class="row q-col-gutter-sm q-mb-md">
        <!-- Metric 1: Ore Assenza Effettuate -->
        <div class="col-6 col-sm-3">
          <div class="kpi-mini-card bg-slate-50 border border-slate-200 rounded-xl q-pa-sm text-center">
            <div class="text-caption text-slate-500">{{ t('parent.absenceLimit.totalHours') || 'Ore di Assenza' }}</div>
            <div class="text-h6 text-weight-bold q-mt-xs" :class="`text-${statusColor}`">
              {{ totalAbsenceHours.toFixed(1) }}h
            </div>
            <div class="text-caption text-slate-400 text-xs">
              ~{{ totalDaysEquivalent }} {{ t('parent.absenceLimit.daysEquiv') || 'giorni equiv.' }}
            </div>
          </div>
        </div>

        <!-- Metric 2: Ore Residue -->
        <div class="col-6 col-sm-3">
          <div class="kpi-mini-card bg-slate-50 border border-slate-200 rounded-xl q-pa-sm text-center">
            <div class="text-caption text-slate-500">{{ t('parent.absenceLimit.remainingAllowance') || 'Ore Rimanenti' }}</div>
            <div class="text-h6 text-weight-bold q-mt-xs" :class="remainingHours <= 30 ? 'text-negative' : 'text-positive'">
              {{ Math.max(0, remainingHours).toFixed(1) }}h
            </div>
            <div class="text-caption text-slate-400 text-xs">
              {{ quotaUsedPercent.toFixed(0) }}% {{ t('parent.absenceLimit.quotaUsed') || 'della quota usata' }}
            </div>
          </div>
        </div>

        <!-- Metric 3: Ritardi & Uscite Anticipate -->
        <div class="col-6 col-sm-3">
          <div class="kpi-mini-card bg-slate-50 border border-slate-200 rounded-xl q-pa-sm text-center">
            <div class="text-caption text-slate-500">{{ t('parent.absenceLimit.latesAndExits') || 'Ritardi / Uscite' }}</div>
            <div class="text-h6 text-weight-bold text-amber-800 q-mt-xs">
              {{ totalLates }} / {{ totalEarlyExits }}
            </div>
            <div class="text-caption text-slate-400 text-xs">
              {{ t('parent.absenceLimit.partialImpact') || 'incidono sul monte ore' }}
            </div>
          </div>
        </div>

        <!-- Metric 4: Da Giustificare -->
        <div class="col-6 col-sm-3">
          <div class="kpi-mini-card bg-slate-50 border border-slate-200 rounded-xl q-pa-sm text-center">
            <div class="text-caption text-slate-500">{{ t('parent.absenceLimit.unjustified') || 'Da Giustificare' }}</div>
            <div class="text-h6 text-weight-bold q-mt-xs" :class="unjustifiedCount > 0 ? 'text-negative' : 'text-positive'">
              {{ unjustifiedCount }}
            </div>
            <router-link to="/parent/attendance" class="text-caption text-primary text-weight-medium text-xs text-decoration-none">
              {{ t('parent.absenceLimit.viewAll') || 'Apri registro →' }}
            </router-link>
          </div>
        </div>
      </div>

      <!-- Critical Alert Banner if Above or Close to 25% -->
      <q-banner
        v-if="absencePercentOfTotal >= 20"
        rounded
        dense
        :class="absencePercentOfTotal >= 25 ? 'bg-red-1 text-negative border-red-3' : 'bg-amber-1 text-amber-9 border-amber-3'"
        class="border q-mb-sm"
      >
        <template v-slot:avatar>
          <q-icon :name="absencePercentOfTotal >= 25 ? 'error' : 'warning'" :color="absencePercentOfTotal >= 25 ? 'negative' : 'warning'" size="sm" />
        </template>
        <div class="text-weight-bold text-caption">
          {{ absencePercentOfTotal >= 25
            ? (t('parent.absenceLimit.alertExceeded') || 'ATTENZIONE: Superato il 25% di assenze (limite di validità dell\'anno)!')
            : (t('parent.absenceLimit.alertApproaching') || 'AVVISO: Le assenze si stanno avvicinando alla soglia critica del 25%.')
          }}
        </div>
        <div class="text-caption text-xs q-mt-xs">
          {{ t('parent.absenceLimit.alertDetail') || 'Si raccomanda di verificare con il coordinatore di classe eventuali deroghe per motivi di salute documentati.' }}
        </div>
      </q-banner>

      <!-- Legal Context Accordion -->
      <q-expansion-item
        dense
        dense-toggle
        expand-separator
        icon="gavel"
        :label="t('parent.absenceLimit.legalTitle') || 'Riferimento Normativo (Art. 14 DPR 122/2009)'"
        header-class="text-caption text-slate-600 q-px-xs"
        class="bg-slate-50 rounded-lg"
      >
        <div class="q-pa-sm text-caption text-slate-600 text-xs">
          <p class="q-mb-xs">
            <strong>D.P.R. 22 giugno 2009, n. 122, art. 14 comma 7:</strong>
            <em>"Ai fini della validità dell'anno scolastico... per poter accedere alla valutazione finale in sede di scrutinio, è necessaria la frequenza di almeno tre quarti dell'orario annuale personalizzato."</em>
          </p>
          <p class="q-mb-none text-slate-500">
            Il Collegio Docenti può deliberare motivate e documentate deroghe per assenze continuative (gravi motivi di salute, terapie, partecipazione ad attività sportive agonistiche certificate).
          </p>
        </div>
      </q-expansion-item>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import attendanceService from '@/services/attendanceService'

const props = defineProps({
  studentId: {
    type: String,
    required: false,
    default: ''
  },
  annualHours: {
    type: Number,
    default: 990
  },
  initialSummary: {
    type: Object,
    default: () => null
  }
})

const { t } = useI18n()

const loading = ref(false)
const currentAnnualHours = ref(props.annualHours || 990)
const summaryData = ref(props.initialSummary || null)
const rawRecords = ref([])

const annualHoursPresets = [
  { label: 'Standard (30h/settimana)', value: 990 },
  { label: 'Tecnico/Professionale (32h/settimana)', value: 1056 },
  { label: 'Tempo Pieno / Liceo Artistico (34h/settimana)', value: 1122 },
  { label: 'Scuola Secondaria I Grado (27h/settimana)', value: 891 }
]

const loadAttendanceData = async () => {
  if (!props.studentId) return

  loading.value = true
  try {
    const [sumRes, attRes] = await Promise.allSettled([
      attendanceService.getChildAttendanceSummary(props.studentId),
      attendanceService.getChildAttendance(props.studentId)
    ])

    if (sumRes.status === 'fulfilled' && sumRes.value?.data) {
      summaryData.value = sumRes.value.data
    }

    if (attRes.status === 'fulfilled' && attRes.value?.data) {
      rawRecords.value = Array.isArray(attRes.value.data) ? attRes.value.data : (attRes.value.data.records || [])
    }
  } catch (err) {
    console.warn('Could not load absence limit details', err)
  } finally {
    loading.value = false
  }
}

watch(() => props.studentId, (newId) => {
  if (newId) {
    loadAttendanceData()
  }
}, { immediate: true })

watch(() => props.annualHours, (newH) => {
  if (newH) currentAnnualHours.value = newH
})

// Calculations
const maxAllowedAbsenceHours = computed(() => {
  return Math.floor(currentAnnualHours.value * 0.25)
})

const totalAbsenceHours = computed(() => {
  // If raw hourly records are present, calculate exact absent hours
  if (rawRecords.value && rawRecords.value.length > 0) {
    let hours = 0
    rawRecords.value.forEach(r => {
      const status = (r.status || '').toLowerCase()
      if (status === 'absent') {
        hours += 1
      } else if (status === 'late' || status === 'early_exit' || status === 'left_early') {
        hours += 0.5 // Average half hour impact for partial presence
      }
    })
    return hours
  }

  // Fallback to summaryData daily count multiplied by standard 5.5 hours/day
  if (summaryData.value) {
    const absentDays = summaryData.value.total_absences || 0
    const lates = summaryData.value.total_lates || 0
    const earlyExits = summaryData.value.total_early_exits || 0
    return (absentDays * 5.5) + (lates * 0.5) + (earlyExits * 0.5)
  }

  return 0
})

const totalDaysEquivalent = computed(() => {
  if (summaryData.value?.total_absences) {
    return summaryData.value.total_absences
  }
  return (totalAbsenceHours.value / 5.5).toFixed(1)
})

const totalLates = computed(() => {
  if (summaryData.value?.total_lates !== undefined) {
    return summaryData.value.total_lates
  }
  return rawRecords.value.filter(r => (r.status || '').toLowerCase() === 'late').length
})

const totalEarlyExits = computed(() => {
  if (summaryData.value?.total_early_exits !== undefined) {
    return summaryData.value.total_early_exits
  }
  return rawRecords.value.filter(r => ['early_exit', 'left_early'].includes((r.status || '').toLowerCase())).length
})

const unjustifiedCount = computed(() => {
  if (rawRecords.value && rawRecords.value.length > 0) {
    return rawRecords.value.filter(r => (r.status || '').toLowerCase() === 'absent' && !r.is_justified && !r.justified).length
  }
  if (summaryData.value) {
    const total = summaryData.value.total_absences || 0
    const justified = summaryData.value.justified_count || 0
    return Math.max(0, total - justified)
  }
  return 0
})

const absencePercentOfTotal = computed(() => {
  if (currentAnnualHours.value <= 0) return 0
  return Math.min(100, (totalAbsenceHours.value / currentAnnualHours.value) * 100)
})

const quotaUsedPercent = computed(() => {
  if (maxAllowedAbsenceHours.value <= 0) return 0
  return Math.min(100, (totalAbsenceHours.value / maxAllowedAbsenceHours.value) * 100)
})

const remainingHours = computed(() => {
  return maxAllowedAbsenceHours.value - totalAbsenceHours.value
})

const progressRatio = computed(() => {
  return Math.min(1, absencePercentOfTotal.value / 100)
})

// Status styling
const statusColor = computed(() => {
  if (absencePercentOfTotal.value >= 25) return 'negative'
  if (absencePercentOfTotal.value >= 20) return 'deep-orange'
  if (absencePercentOfTotal.value >= 15) return 'amber-9'
  return 'positive'
})

const statusIcon = computed(() => {
  if (absencePercentOfTotal.value >= 25) return 'error'
  if (absencePercentOfTotal.value >= 20) return 'warning'
  if (absencePercentOfTotal.value >= 15) return 'info'
  return 'check_circle'
})

const statusLabel = computed(() => {
  if (absencePercentOfTotal.value >= 25) return t('parent.absenceLimit.statusCritical') || 'Soglia Superata (>25%)'
  if (absencePercentOfTotal.value >= 20) return t('parent.absenceLimit.statusDanger') || 'Rischio Elevato (≥20%)'
  if (absencePercentOfTotal.value >= 15) return t('parent.absenceLimit.statusWarning') || 'Attenzione (≥15%)'
  return t('parent.absenceLimit.statusOk') || 'Regolare (<15%)'
})

onMounted(() => {
  if (props.studentId) {
    loadAttendanceData()
  }
})
</script>

<style scoped>
.absence-limit-card {
  background: white;
  border: 1px solid #e2e8f0;
}

.limit-marker-pin {
  width: 2px;
  transform: translateX(-50%);
  pointer-events: none;
}

.marker-line {
  width: 2px;
  height: 14px;
  border-radius: 1px;
}

.marker-label {
  font-size: 10px;
  position: absolute;
  top: 15px;
  left: 50%;
  transform: translateX(-50%);
}

.kpi-mini-card {
  transition: all 0.2s ease;
}

.kpi-mini-card:hover {
  border-color: #cbd5e1;
  background-color: #f1f5f9;
}

.text-xs {
  font-size: 0.75rem;
}
</style>
