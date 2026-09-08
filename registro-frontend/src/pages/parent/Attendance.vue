<template>
  <q-page class="q-pa-md bg-slate-50">
    <!-- Header with Child Switcher -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          Presenze & Giustificazioni
        </h1>
        <p class="text-subtitle2 text-slate-500 q-mt-xs q-mb-none">
          Monitoraggio presenze, giustificazioni digitali e statistiche scolastiche
        </p>
      </div>

      <q-select
        v-if="childOptions.length > 1"
        v-model="selectedChildId"
        :options="childOptions"
        label="Seleziona Figlio"
        outlined
        dense
        emit-value
        map-options
        class="bg-white"
        style="min-width: 200px"
        @update:model-value="loadChildData"
      />
    </div>

    <!-- Navigation Tabs -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft q-mb-md">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-slate-600 bg-slate-100 border-b border-slate-200"
        active-color="primary"
        indicator-color="primary"
        align="left"
        no-caps
      >
        <q-tab name="list" icon="list" label="Registro Presenze" />
        <q-tab name="unjustified">
          <div class="row items-center no-wrap">
            <q-icon name="assignment_late" class="q-mr-xs" />
            <span>Da Giustificare</span>
            <q-badge v-if="unjustifiedList.length > 0" color="negative" class="q-ml-xs text-weight-bold">
              {{ unjustifiedList.length }}
            </q-badge>
          </div>
        </q-tab>
        <q-tab name="stats" icon="analytics" label="Statistiche & KPI" />
      </q-tabs>
    </q-card>

    <!-- 75% PRESENCE RISK WARNING BANNER -->
    <q-banner
      v-if="statsData && statsData.absence_percentage > 25"
      role="alert"
      aria-live="assertive"
      class="bg-negative text-white rounded-xl q-mb-md shadow-md"
    >
      <template v-slot:avatar>
        <q-icon name="warning" color="white" size="md" />
      </template>
      <div class="text-weight-bold text-subtitle1">ATTENZIONE: Rischio Non Ammissione all'Anno Scolastico!</div>
      <div class="text-body2">
        La percentuale di frequenza è inferiore al 75% minimo richiesto dalla normativa italiana per la validità dell'anno scolastico.
      </div>
    </q-banner>

    <!-- TAB 1: REGISTRO PRESENZE -->
    <div v-if="activeTab === 'list'">
      <!-- Stats KPI cards summary -->
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-12 col-sm-3">
          <q-card flat bordered class="bg-green-50 border-green-200 text-positive rounded-xl shadow-soft">
            <q-card-section class="text-center">
              <div class="text-h4 text-weight-bold">{{ presenceRate }}%</div>
              <div class="text-caption text-slate-600">Presenza Totale</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-sm-3">
          <q-card flat bordered class="bg-red-50 border-red-200 text-negative rounded-xl shadow-soft">
            <q-card-section class="text-center">
              <div class="text-h4 text-weight-bold">{{ summary.total_absences ?? 0 }}</div>
              <div class="text-caption text-slate-600">Assenze Totali</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-sm-3">
          <q-card flat bordered class="bg-amber-50 border-amber-200 text-warning rounded-xl shadow-soft">
            <q-card-section class="text-center">
              <div class="text-h4 text-weight-bold">{{ summary.total_lates ?? 0 }}</div>
              <div class="text-caption text-slate-600">Ritardi</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-12 col-sm-3">
          <q-card flat bordered :class="unjustifiedList.length > 0 ? 'bg-orange-50 border-orange-200 text-orange-8' : 'bg-slate-50 border-slate-200 text-slate-600'" class="rounded-xl shadow-soft">
            <q-card-section class="text-center">
              <div class="text-h4 text-weight-bold">{{ unjustifiedList.length }}</div>
              <div class="text-caption">Da Giustificare</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Events List -->
      <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
        <q-card-section class="bg-slate-100 border-b border-slate-200 row items-center justify-between q-py-sm q-px-md">
          <div class="text-subtitle1 text-weight-bold text-slate-800">Storico Registro Presenze</div>
        </q-card-section>

        <div v-if="loading" class="text-center q-pa-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>

        <div v-else-if="events.length === 0" class="text-center q-pa-xl text-slate-400">
          <q-icon name="check_circle" size="64px" color="positive" class="q-mb-sm opacity-50" />
          <div class="text-h6 text-slate-700">Nessun evento o assenza registrata</div>
        </div>

        <q-list v-else separator>
          <q-item v-for="ev in events" :key="ev.id || ev.date" class="q-py-md">
            <q-item-section avatar>
              <q-avatar :color="getEventBadgeColor(ev.status)" text-color="white" :icon="getEventIcon(ev.status)" size="40px" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-bold text-slate-800">{{ getStatusLabel(ev.status) }}</q-item-label>
              <q-item-label caption class="text-slate-500">Data: {{ formatDate(ev.date) }}</q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-chip v-if="ev.justified || ev.parent_justified" color="positive" text-color="white" size="xs" icon="check_circle">
                Giustificata
              </q-chip>
              <q-chip v-else color="warning" text-color="white" size="xs" icon="pending">
                Non Giustificata
              </q-chip>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>
    </div>

    <!-- TAB 2: DA GIUSTIFICARE -->
    <div v-else-if="activeTab === 'unjustified'">
      <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
        <q-card-section class="bg-slate-100 border-b border-slate-200 row items-center justify-between q-py-sm q-px-md">
          <div class="text-subtitle1 text-weight-bold text-slate-800">Assenze e Ritardi in Attesa di Giustifica</div>
          <q-badge color="negative" class="q-px-sm text-weight-bold">
            {{ unjustifiedList.length }} In Attesa
          </q-badge>
        </q-card-section>

        <div v-if="unjustifiedList.length === 0" class="text-center q-pa-xl text-slate-400">
          <q-icon name="task_alt" size="64px" color="positive" class="q-mb-md opacity-40" />
          <div class="text-h6 text-slate-700">Tutte le assenze sono state giustificate!</div>
        </div>

        <q-list v-else separator>
          <q-item v-for="att in unjustifiedList" :key="att.id" class="q-py-md">
            <q-item-section avatar>
              <q-avatar color="red-1" text-color="negative" icon="priority_high" size="40px" />
            </q-item-section>

            <q-item-section>
              <q-item-label class="text-weight-bold text-slate-800">
                {{ getStatusLabel(att.status) }} del {{ formatDate(att.date) }}
              </q-item-label>
              <q-item-label caption class="text-slate-500">
                Stato: Non Giustificata
              </q-item-label>
            </q-item-section>

            <q-item-section side>
              <q-btn
                color="primary"
                unelevated
                size="sm"
                icon="draw"
                label="Giustifica"
                no-caps
                @click="openJustifyModal(att)"
              />
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>
    </div>

    <!-- TAB 3: STATISTICHE GRAFICHE -->
    <div v-else-if="activeTab === 'stats'">
      <!-- Ministerial Absence Limit 25% Widget -->
      <AbsenceLimitWidget :student-id="selectedChildId" class="q-mb-lg" />

      <!-- KPI Row -->
      <div class="row q-col-gutter-md q-mb-lg">
        <div class="col-6 col-md-3">
          <q-card flat bordered class="bg-white rounded-xl shadow-soft text-center q-pa-sm">
            <div class="text-caption text-slate-400">Giorni Scuola</div>
            <div class="text-h4 text-weight-bold text-slate-800">{{ statsData?.total_school_days || 150 }}</div>
          </q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card flat bordered class="bg-white rounded-xl shadow-soft text-center q-pa-sm">
            <div class="text-caption text-slate-400">Presenti</div>
            <div class="text-h4 text-weight-bold text-positive">{{ statsData?.days_present ?? 0 }}</div>
          </q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card flat bordered class="bg-white rounded-xl shadow-soft text-center q-pa-sm">
            <div class="text-caption text-slate-400">Assenti</div>
            <div class="text-h4 text-weight-bold text-negative">{{ statsData?.days_absent ?? 0 }}</div>
          </q-card>
        </div>
        <div class="col-6 col-md-3">
          <q-card flat bordered class="bg-white rounded-xl shadow-soft text-center q-pa-sm">
            <div class="text-caption text-slate-400">% Presenza</div>
            <div class="text-h4 text-weight-bold" :class="getPresenceColorClass(statsPresencePercentage)">
              {{ (100 - (statsData?.absence_percentage ?? 0)).toFixed(1) }}%
            </div>
          </q-card>
        </div>
      </div>

      <!-- Charts Row -->
      <div class="row q-col-gutter-md">
        <!-- Donut Progress -->
        <div class="col-12 col-md-5">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft q-pa-md h-full">
            <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Ripartizione Presenze</div>

            <div class="column items-center justify-center space-y-4 q-py-md">
              <q-circular-progress
                show-value
                font-size="16px"
                :value="100 - (statsData?.absence_percentage ?? 0)"
                size="160px"
                :thickness="0.2"
                color="positive"
                track-color="red-2"
                class="text-weight-bold text-slate-800"
              >
                {{ (100 - (statsData?.absence_percentage ?? 0)).toFixed(1) }}%
              </q-circular-progress>

              <div class="w-full space-y-2 q-mt-md">
                <div class="row items-center justify-between text-body2">
                  <span class="row items-center"><div class="w-3 h-3 rounded bg-positive q-mr-xs" /> Presenti</span>
                  <span class="text-weight-bold">{{ statsData?.days_present ?? 0 }} giorni</span>
                </div>
                <div class="row items-center justify-between text-body2">
                  <span class="row items-center"><div class="w-3 h-3 rounded bg-blue-500 q-mr-xs" /> Assenti Giustificate</span>
                  <span class="text-weight-bold">{{ statsData?.justified ?? 0 }} giorni</span>
                </div>
                <div class="row items-center justify-between text-body2">
                  <span class="row items-center"><div class="w-3 h-3 rounded bg-negative q-mr-xs" /> Assenti Non Giustificate</span>
                  <span class="text-weight-bold">{{ statsData?.unjustified ?? 0 }} giorni</span>
                </div>
              </div>
            </div>
          </q-card>
        </div>

        <!-- Monthly Breakdown Bar Chart -->
        <div class="col-12 col-md-7">
          <q-card flat bordered class="rounded-xl bg-white shadow-soft q-pa-md h-full">
            <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Andamento Mensile Presenze</div>

            <div class="space-y-4 q-py-sm">
              <div v-for="m in statsData?.monthly_breakdown || []" :key="m.month" class="space-y-1">
                <div class="row items-center justify-between text-caption text-weight-bold text-slate-700">
                  <span>{{ m.month }}</span>
                  <span>{{ m.present }} Pres. / {{ m.absent }} Ass.</span>
                </div>
                <div class="w-full bg-slate-100 rounded-full h-4 overflow-hidden flex" :title="`Mese ${m.month}: ${m.present} presenti, ${m.absent} assenti`" :aria-label="`Mese ${m.month}: ${m.present} presenti, ${m.absent} assenti`">
                  <div class="bg-positive h-full" :style="{ width: ((m.present / (m.present + m.absent || 1)) * 100) + '%' }" />
                  <div class="bg-negative h-full" :style="{ width: ((m.absent / (m.present + m.absent || 1)) * 100) + '%' }" />
                </div>
              </div>
            </div>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Justification Dialog -->
    <q-dialog v-model="justifyModal">
      <q-card style="min-width: 360px" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div class="text-h6 text-weight-bold">Giustifica Assenza</div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-md space-y-4">
          <q-select
            v-model="justifyForm.reason"
            :options="['Malattia', 'Motivi familiari', 'Visita medica', 'Altro']"
            label="Motivo Assenza"
            outlined
            dense
          />
          <q-input
            v-model="justifyForm.notes"
            type="textarea"
            rows="3"
            label="Note aggiuntive (opzionale)"
            outlined
            dense
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn
            color="primary"
            unelevated
            label="Conferma Giustifica"
            :loading="submittingJustify"
            no-caps
            @click="submitJustify"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar, date as qdate } from 'quasar'
import { useAttendanceStore } from '@/stores/attendance'
import api from '@/services/api'
import AbsenceLimitWidget from '@/components/Parent/AbsenceLimitWidget.vue'

const $q = useQuasar()
const { t } = useI18n()
const attendanceStore = useAttendanceStore()

const activeTab = ref('list')
const loading = ref(false)
const children = ref([])
const selectedChildId = ref(null)

const events = ref([])
const unjustifiedList = ref([])
const statsData = ref(null)
const summary = ref({ total_absences: 0, total_lates: 0 })

const justifyModal = ref(false)
const selectedAttendance = ref(null)
const submittingJustify = ref(false)
const justifyForm = ref({
  reason: 'Malattia',
  notes: ''
})

const childOptions = computed(() => {
  return children.value.map(c => ({
    label: `${c.first_name || c.name} ${c.last_name || ''}`,
    value: c.user_id || c.id  // use user_id; attendance.student_id = users.id
  }))
})

const presenceRate = computed(() => {
  if (!statsData.value) return 92
  const pct = 100 - (statsData.value.absence_percentage || 0)
  return pct.toFixed(1)
})

const statsPresencePercentage = computed(() => {
  return parseFloat(presenceRate.value)
})

onMounted(async () => {
  await fetchChildren()
  await loadChildData()
})

async function fetchChildren() {
  try {
    const res = await api.get('/users/me/children')
    children.value = res.data || []
    if (children.value.length > 0) {
      // Use user_id for attendance calls since attendance.student_id = users.id
      selectedChildId.value = children.value[0].user_id || children.value[0].id
    }
  } catch {
    children.value = []
  }
}

async function loadChildData() {
  if (!selectedChildId.value) return
  loading.value = true
  try {
    const [attRes, unjRes, statsRes] = await Promise.all([
      api.get(`/attendance/child-attendance/${selectedChildId.value}`).catch(() => ({ data: [] })),
      attendanceStore.fetchUnjustified(selectedChildId.value),
      attendanceStore.fetchChildStats(selectedChildId.value)
    ])

    events.value = attRes.data || []
    unjustifiedList.value = unjRes || []
    statsData.value = statsRes || null

    summary.value = {
      total_absences: events.value.filter(e => e.status === 'Absent').length,
      total_lates: events.value.filter(e => e.status === 'Late').length
    }
  } finally {
    loading.value = false
  }
}

function formatDate(d) {
  return d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY') : ''
}

function getStatusLabel(s) {
  switch (s) {
    case 'Absent': return 'Assenza'
    case 'Late': return 'Ritardo'
    case 'LeftEarly': return 'Uscita Anticipata'
    case 'Present': return 'Presente'
    default: return s
  }
}

function getEventBadgeColor(s) {
  switch (s) {
    case 'Absent': return 'negative'
    case 'Late': return 'warning'
    case 'LeftEarly': return 'info'
    default: return 'positive'
  }
}

function getEventIcon(s) {
  switch (s) {
    case 'Absent': return 'close'
    case 'Late': return 'schedule'
    case 'LeftEarly': return 'logout'
    default: return 'check'
  }
}

function getPresenceColorClass(pct) {
  if (pct > 90) return 'text-positive font-bold'
  if (pct >= 75) return 'text-warning font-bold'
  return 'text-negative font-bold'
}

function openJustifyModal(att) {
  selectedAttendance.value = att
  justifyForm.value = { reason: 'Malattia', notes: '' }
  justifyModal.value = true
}

async function submitJustify() {
  if (!selectedAttendance.value) return
  submittingJustify.value = true
  try {
    await attendanceStore.justifyAbsence(
      selectedChildId.value,
      selectedAttendance.value.id,
      justifyForm.value.reason,
      justifyForm.value.notes
    )
    $q.notify({ type: 'positive', message: 'Assenza giustificata con successo' })
    justifyModal.value = false
    await loadChildData()
  } catch {
    $q.notify({ type: 'negative', message: 'Errore durante l\'invio della giustifica' })
  } finally {
    submittingJustify.value = false
  }
}
</script>
