<template>
  <q-card class="glass-card shadow-soft rounded-xl overflow-hidden border border-slate-100">
    <q-card-section class="bg-gradient-to-r from-red-50 to-amber-50 q-pa-lg border-b border-slate-100">
      <div class="row items-center justify-between">
        <div>
          <div class="text-h6 text-weight-bold text-slate-800 row items-center">
            <q-icon name="warning" color="negative" class="q-mr-sm" size="24px" />
            {{ t('dropoutRisk.title') || 'Cruscotto Rischio Dispersione Scolastica (Early Warning)' }}
          </div>
          <div class="text-caption text-slate-600 q-mt-xs">
            {{ t('dropoutRisk.subtitle') || 'Monitoraggio preventivo dispersione ai sensi dell\'art. 14 DPR 122/2009' }}
          </div>
        </div>
        <div class="row items-center q-gutter-sm">
          <q-btn
            outline
            color="primary"
            icon="download"
            :label="t('dropoutRisk.exportPlan') || 'Esporta Piano di Supporto (CSV)'"
            class="rounded-lg bg-white"
            @click="exportSupportPlan"
            :loading="exporting"
            :disable="loading || items.length === 0"
          />
          <q-btn
            flat
            round
            dense
            icon="refresh"
            color="grey-7"
            @click="fetchData"
            :loading="loading"
          >
            <q-tooltip>Aggiorna</q-tooltip>
          </q-btn>
        </div>
      </div>

      <!-- Quick KPI Chips -->
      <div class="row q-gutter-sm q-mt-md">
        <q-chip dense color="red-1" text-color="negative" icon="dangerous">
          <strong>{{ countCritical }}</strong> &nbsp;{{ t('dropoutRisk.riskCritical') || 'Critico' }} (&ge; 25% assenze)
        </q-chip>
        <q-chip dense color="orange-1" text-color="deep-orange" icon="warning">
          <strong>{{ countHigh }}</strong> &nbsp;{{ t('dropoutRisk.riskHigh') || 'Alto' }}
        </q-chip>
        <q-chip dense color="amber-1" text-color="amber-9" icon="report_problem">
          <strong>{{ countMedium }}</strong> &nbsp;{{ t('dropoutRisk.riskMedium') || 'Moderato' }}
        </q-chip>
        <q-chip dense color="green-1" text-color="positive" icon="check_circle">
          <strong>{{ countLow }}</strong> &nbsp;{{ t('dropoutRisk.riskLow') || 'Basso' }}
        </q-chip>
      </div>
    </q-card-section>

    <!-- Filters Bar -->
    <q-card-section class="q-pa-md bg-white border-b border-slate-100">
      <div class="row q-col-gutter-md items-center">
        <div class="col-12 col-sm-6 col-md-4">
          <q-select
            v-model="filterClass"
            :options="classOptions"
            option-label="label"
            option-value="value"
            emit-value
            map-options
            dense
            outlined
            clearable
            :label="t('dropoutRisk.filterClass') || 'Filtra per Classe'"
            @update:model-value="fetchData"
          >
            <template v-slot:prepend>
              <q-icon name="class" />
            </template>
          </q-select>
        </div>
        <div class="col-12 col-sm-6 col-md-4">
          <q-select
            v-model="filterRisk"
            :options="riskOptions"
            emit-value
            map-options
            dense
            outlined
            clearable
            :label="t('dropoutRisk.filterRisk') || 'Filtra per Rischio'"
            @update:model-value="fetchData"
          >
            <template v-slot:prepend>
              <q-icon name="filter_alt" />
            </template>
          </q-select>
        </div>
      </div>
    </q-card-section>

    <!-- Table -->
    <q-card-section class="q-pa-none">
      <q-table
        :rows="items"
        :columns="columns"
        row-key="student_id"
        :loading="loading"
        flat
        :pagination="pagination"
        :no-data-label="t('dropoutRisk.noStudentsAtRisk') || 'Nessuno studente a rischio rilevato'"
      >
        <!-- Student Slot -->
        <template v-slot:body-cell-student="props">
          <q-td :props="props">
            <div class="text-weight-bold text-slate-800">
              {{ props.row.last_name }} {{ props.row.first_name }}
            </div>
            <div class="text-caption text-slate-500">ID: {{ props.row.student_id.substring(0, 8) }}...</div>
          </q-td>
        </template>

        <!-- Absence Rate Slot -->
        <template v-slot:body-cell-absence_rate="props">
          <q-td :props="props">
            <q-badge
              :color="getAbsenceBadgeColor(props.row.absence_rate)"
              class="text-caption font-bold q-pa-xs"
            >
              {{ props.row.absence_rate }}%
            </q-badge>
            <div class="text-caption text-slate-500">
              {{ props.row.absence_hours }}h su {{ props.row.total_hours }}h
            </div>
            <q-tooltip v-if="props.row.absence_rate >= 25">
              Superata soglia critica 25% DPR 122/2009 (Rischio non ammissione)
            </q-tooltip>
          </q-td>
        </template>

        <!-- Failing Subjects Slot -->
        <template v-slot:body-cell-failing_subjects="props">
          <q-td :props="props">
            <div v-if="props.row.failing_subjects_count > 0" class="row items-center q-gutter-xs">
              <q-badge color="deep-orange" class="font-bold">
                {{ props.row.failing_subjects_count }} materie &lt; 5.0
              </q-badge>
              <q-tooltip>
                {{ props.row.failing_subjects.join(', ') }}
              </q-tooltip>
            </div>
            <span v-else class="text-slate-400 text-caption">Nessuna</span>
          </q-td>
        </template>

        <!-- Delays Slot -->
        <template v-slot:body-cell-delays="props">
          <q-td :props="props">
            <span class="text-slate-700">
              {{ props.row.lates_count }} ritardi / {{ props.row.early_exits_count }} uscite
            </span>
          </q-td>
        </template>

        <!-- Risk Score Slot -->
        <template v-slot:body-cell-risk_level="props">
          <q-td :props="props">
            <q-chip
              dense
              :color="getRiskLevelColor(props.row.risk_level)"
              text-color="white"
              class="font-bold text-caption"
            >
              {{ props.row.risk_level }} ({{ props.row.risk_score }})
            </q-chip>
          </q-td>
        </template>

        <!-- Action / Support Plan Slot -->
        <template v-slot:body-cell-recommended_action="props">
          <q-td :props="props" style="max-width: 320px; white-space: normal;">
            <div class="text-caption text-slate-700">
              {{ props.row.recommended_action }}
            </div>
          </q-td>
        </template>
      </q-table>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from '@/services/api'
import { useClassesStore } from '@/stores/classes'

const $q = useQuasar()
const { t } = useI18n()
const classesStore = useClassesStore()

const items = ref([])
const loading = ref(false)
const exporting = ref(false)
const filterClass = ref(null)
const filterRisk = ref(null)

const pagination = ref({
  sortBy: 'risk_score',
  descending: true,
  page: 1,
  rowsPerPage: 10
})

const classOptions = computed(() => {
  return (classesStore.classes || []).map(c => ({
    label: c.name + (c.section ? ` ${c.section}` : ''),
    value: c.id
  }))
})

const riskOptions = computed(() => [
  { label: 'Tutti i Livelli', value: '' },
  { label: t('dropoutRisk.riskCritical') || 'Critico', value: 'Critico' },
  { label: t('dropoutRisk.riskHigh') || 'Alto', value: 'Alto' },
  { label: t('dropoutRisk.riskMedium') || 'Moderato', value: 'Moderato' },
  { label: t('dropoutRisk.riskLow') || 'Basso', value: 'Basso' }
])

const columns = computed(() => [
  { name: 'student', label: t('dropoutRisk.student') || 'Studente', align: 'left', field: 'last_name', sortable: true },
  { name: 'class_name', label: t('dropoutRisk.class') || 'Classe', align: 'left', field: 'class_name', sortable: true },
  { name: 'absence_rate', label: t('dropoutRisk.absenceRate') || 'Tasso Assenze', align: 'left', field: 'absence_rate', sortable: true },
  { name: 'failing_subjects', label: t('dropoutRisk.failingSubjects') || 'Materie Insufficienti', align: 'left', field: 'failing_subjects_count', sortable: true },
  { name: 'delays', label: t('dropoutRisk.anomalousDelays') || 'Ritardi & Uscite', align: 'left', field: 'lates_count', sortable: true },
  { name: 'risk_level', label: t('dropoutRisk.riskLevel') || 'Livello Rischio', align: 'left', field: 'risk_score', sortable: true },
  { name: 'recommended_action', label: 'Piano di Supporto Raccomandato', align: 'left', field: 'recommended_action' }
])

const countCritical = computed(() => items.value.filter(i => i.risk_level === 'Critico').length)
const countHigh = computed(() => items.value.filter(i => i.risk_level === 'Alto').length)
const countMedium = computed(() => items.value.filter(i => i.risk_level === 'Moderato').length)
const countLow = computed(() => items.value.filter(i => i.risk_level === 'Basso').length)

const getAbsenceBadgeColor = (rate) => {
  if (rate >= 25) return 'negative'
  if (rate >= 20) return 'deep-orange'
  if (rate >= 15) return 'amber-9'
  return 'positive'
}

const getRiskLevelColor = (level) => {
  switch (level) {
    case 'Critico': return 'negative'
    case 'Alto': return 'deep-orange'
    case 'Moderato': return 'amber-9'
    default: return 'positive'
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {}
    if (filterClass.value) params.class_id = filterClass.value
    if (filterRisk.value) params.risk_level = filterRisk.value

    const res = await api.get('/reports/dropout-risk', { params })
    items.value = res.data?.data || []
  } catch (err) {
    console.error('Error fetching dropout risk:', err)
    $q.notify({ type: 'negative', message: 'Errore nel caricamento del rischio dispersione scolastica' })
  } finally {
    loading.value = false
  }
}

const exportSupportPlan = async () => {
  exporting.value = true
  try {
    const params = {}
    if (filterClass.value) params.class_id = filterClass.value
    if (filterRisk.value) params.risk_level = filterRisk.value

    const response = await api.get('/reports/dropout-risk/export', {
      params,
      responseType: 'blob'
    })

    const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8;' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `piano_supporto_dispersione_${new Date().toISOString().slice(0, 10)}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    $q.notify({ type: 'positive', message: t('dropoutRisk.exportSuccess') || 'Piano di supporto scaricato con successo' })
  } catch (err) {
    console.error('Export error:', err)
    $q.notify({ type: 'negative', message: 'Errore durante l\'esportazione del piano di supporto' })
  } finally {
    exporting.value = false
  }
}

onMounted(async () => {
  if (!classesStore.classes || classesStore.classes.length === 0) {
    if (typeof classesStore.fetchClasses === 'function') {
      try {
        await classesStore.fetchClasses()
      } catch (err) {
        // Silently ignore if classes cannot be loaded
        void err
      }
    }
  }
  await fetchData()
})
</script>
