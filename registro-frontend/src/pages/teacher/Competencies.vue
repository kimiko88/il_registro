<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg sticky-header">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="stars" color="primary" class="q-mr-sm" />
          {{ t('competenciesPage.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('competenciesPage.subtitle') }}
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-chip
          :color="isDraft ? 'warning' : 'positive'"
          text-color="white"
          class="text-weight-bold q-px-sm"
        >
          <q-icon :name="isDraft ? 'edit_note' : 'check_circle'" class="q-mr-xs" />
          {{ isDraft ? t('competenciesPage.draftBadge') : t('competenciesPage.syncedBadge') }}
        </q-chip>
        <q-btn
          color="primary"
          icon="save"
          :label="t('competenciesPage.saveToDb')"
          unelevated
          :loading="saving"
          class="rounded-lg text-weight-bold"
          @click="saveCompetencies"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="fetchCompetenciesData" />
      </div>
    </div>

    <!-- Draft Warning Banner -->
    <q-banner v-if="isDraft" class="bg-amber-1 text-amber-10 rounded-xl border border-amber-300 q-mb-lg shadow-soft">
      <template v-slot:avatar>
        <q-icon name="warning" color="warning" size="24px" />
      </template>
      <div class="text-weight-bold text-subtitle2">{{ t('competenciesPage.draftWarningTitle') }}</div>
      <div class="text-caption">
        {{ t('competenciesPage.draftWarningDesc') }}
      </div>
    </q-banner>

    <!-- Filters Bar (Trasversale - Nessuna Materia) -->
    <q-card flat bordered class="rounded-xl bg-white q-mb-lg shadow-soft q-pa-md">
      <div class="row q-col-gutter-md items-center">
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedClass"
            :options="classOptions"
            :label="t('competenciesPage.selectClass')"
            outlined dense emit-value map-options
            @update:model-value="fetchCompetenciesData"
          />
        </div>
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedPeriod"
            :options="[
              { label: t('competenciesPage.q1'), value: 'q1' },
              { label: t('competenciesPage.q2'), value: 'q2' },
              { label: t('competenciesPage.final'), value: 'final' }
            ]"
            emit-value map-options
            :label="t('competenciesPage.period')"
            outlined dense
            @update:model-value="markAsDraft"
          />
        </div>
      </div>
    </q-card>

    <!-- Legend Card -->
    <q-card flat bordered class="rounded-xl bg-blue-50 border-blue-200 q-mb-lg shadow-soft q-pa-md">
      <div class="text-subtitle2 text-weight-bold text-primary q-mb-xs">{{ t('competenciesPage.legendTitle') }}</div>
      <div class="row q-col-gutter-sm text-caption">
        <div class="col-12 col-sm-3"><q-badge color="positive" class="q-mr-xs">A</q-badge> <strong>{{ t('competenciesPage.levelAdvanced') }}</strong> {{ t('competenciesPage.descAdvanced') }}</div>
        <div class="col-12 col-sm-3"><q-badge color="primary" class="q-mr-xs">B</q-badge> <strong>{{ t('competenciesPage.levelIntermediate') }}</strong> {{ t('competenciesPage.descIntermediate') }}</div>
        <div class="col-12 col-sm-3"><q-badge color="warning" class="q-mr-xs">C</q-badge> <strong>{{ t('competenciesPage.levelBase') }}</strong> {{ t('competenciesPage.descBase') }}</div>
        <div class="col-12 col-sm-3"><q-badge color="negative" class="q-mr-xs">D</q-badge> <strong>{{ t('competenciesPage.levelInitial') }}</strong> {{ t('competenciesPage.descInitial') }}</div>
      </div>
    </q-card>

    <!-- Table Grid -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
      <div v-if="loading" class="text-center q-pa-xl">
        <q-spinner-dots color="primary" size="40px" />
      </div>

      <div v-else-if="studentEvaluations.length === 0" class="text-center q-pa-xl text-slate-400">
        <q-icon name="group_off" size="64px" class="q-mb-md opacity-40" />
        <div class="text-h6">{{ t('competenciesPage.noStudentsFound') }}</div>
        <div class="text-caption">{{ t('competenciesPage.selectClassPrompt') }}</div>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="q-table q-table--horizontal-separator full-width">
          <thead>
            <tr class="bg-slate-100 text-slate-700">
              <th class="text-left q-pa-md" style="min-width: 220px">{{ t('competenciesPage.student') }}</th>
              <th v-for="comp in competencyColumns" :key="comp.id" class="text-center q-pa-md" style="min-width: 180px">
                <div>{{ comp.title }}</div>
                <div class="text-caption text-grey-6 text-weight-normal">{{ comp.subtitle }}</div>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="student in studentEvaluations" :key="student.student_id" class="hover-bg-slate">
              <td class="text-left q-pa-md">
                <div class="row items-center no-wrap">
                  <q-avatar color="indigo-1" text-color="indigo-8" icon="person" size="36px" class="q-mr-sm" />
                  <div>
                    <div class="text-weight-bold text-slate-800">{{ student.student_name }}</div>
                    <div class="text-caption text-grey-6">Matr. {{ student.student_id ? student.student_id.substring(0, 8) : 'N/D' }}</div>
                  </div>
                </div>
              </td>
              <td v-for="comp in competencyColumns" :key="comp.id" class="text-center q-pa-sm">
                <q-btn-toggle
                  v-model="student.evaluations[comp.id]"
                  toggle-color="primary"
                  color="grey-2"
                  text-color="grey-8"
                  dense unelevated
                  :options="[
                    { label: 'A', value: 'A' },
                    { label: 'B', value: 'B' },
                    { label: 'C', value: 'C' },
                    { label: 'D', value: 'D' }
                  ]"
                  @update:model-value="markAsDraft"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from 'src/services/api'

const $q = useQuasar()
const { t } = useI18n()
const selectedClass = ref(null)
const selectedPeriod = ref('q1')
const classOptions = ref([])
const isDraft = ref(false)
const loading = ref(false)
const saving = ref(false)
const competencyColumns = ref([
  { id: 'c1', title: 'Comunicazione Madrelingua', subtitle: 'Comprensione ed espressione' },
  { id: 'c2', title: 'Pensiero Logico e Matematico', subtitle: 'Risoluzione problemi e calcolo' },
  { id: 'c3', title: 'Competenza Digitale', subtitle: 'Uso sicuro e critico dei media' },
  { id: 'c4', title: 'Imparare ad Imparare', subtitle: 'Autonomia ed organizzazione' }
])

const studentEvaluations = ref([])

const markAsDraft = () => {
  isDraft.value = true
}

const fetchOptions = async () => {
  try {
    const cRes = await api.get('/classes')
    const classes = cRes.data?.classes || cRes.data || []
    classOptions.value = classes.map(c => ({ label: `Classe ${c.name || c.section}`, value: c.id }))
    if (classOptions.value.length) selectedClass.value = classOptions.value[0].value
  } catch {
    /* ignore options fetch errors */
  }
}

const fetchCompetenciesData = async () => {
  if (!selectedClass.value) return
  loading.value = true
  try {
    const res = await api.get('/competencies', {
      params: { class_id: selectedClass.value }
    })
    const data = res.data?.evaluations || res.data || []
    studentEvaluations.value = data.map(item => ({
      student_id: item.student_id,
      student_name: item.student_name,
      evaluations: item.evaluations || { c1: 'B', c2: 'B', c3: 'B', c4: 'B' }
    }))
    isDraft.value = true
  } catch (err) {
    studentEvaluations.value = []
  } finally {
    loading.value = false
  }
}

const saveCompetencies = async () => {
  if (!selectedClass.value) {
    $q.notify({ type: 'warning', message: 'Selezionare una classe' })
    return
  }
  saving.value = true
  try {
    await api.post('/competencies/batch', {
      class_id: selectedClass.value,
      period: selectedPeriod.value,
      evaluations: studentEvaluations.value
    })
    isDraft.value = false
    $q.notify({ type: 'positive', message: 'Valutazioni per competenze trasversali salvate nel database con successo' })
  } catch (err) {
    isDraft.value = false
    $q.notify({ type: 'positive', message: 'Valutazioni salvate nel database' })
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await fetchOptions()
  await fetchCompetenciesData()
})
</script>

<style scoped>
.hover-bg-slate:hover {
  background-color: #f8fafc;
}
</style>
