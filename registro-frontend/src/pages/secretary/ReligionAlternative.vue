<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">
          <q-icon name="balance" class="q-mr-sm" />{{ t('religionAlternative.title') }}
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          {{ t('religionAlternative.subtitle') }}
        </p>
      </div>
      <div class="row q-gutter-sm">
        <q-btn flat icon="refresh" :label="t('common.refresh')" @click="loadData" :loading="loading" />
        <q-btn color="primary" icon="add" :label="t('religionAlternative.title')" unelevated no-caps @click="openOptionDialog" />
      </div>
    </div>

    <!-- Summary cards -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-3">
        <q-card flat bordered class="bg-blue-1 text-primary">
          <q-card-section>
            <div class="text-subtitle2">{{ t('religionAlternative.optionIrc') }}</div>
            <div class="text-h5 text-weight-bold">{{ countByOption('irc') }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-3">
        <q-card flat bordered class="bg-teal-1 text-teal-9">
          <q-card-section>
            <div class="text-subtitle2">{{ t('religionAlternative.optionA') }}</div>
            <div class="text-h5 text-weight-bold">{{ countByOption('materia_alternativa') }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-3">
        <q-card flat bordered class="bg-amber-1 text-amber-9">
          <q-card-section>
            <div class="text-subtitle2">{{ t('religionAlternative.optionB') }} / {{ t('religionAlternative.optionC') }}</div>
            <div class="text-h5 text-weight-bold">{{ countByOption('studio_assistito') + countByOption('studio_libero') }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-3">
        <q-card flat bordered class="bg-purple-1 text-purple-9">
          <q-card-section>
            <div class="text-subtitle2">{{ t('religionAlternative.optionD') }}</div>
            <div class="text-h5 text-weight-bold">{{ countByOption('uscita_scuola') }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Table -->
    <q-card flat bordered class="bg-white">
      <q-tabs v-model="tabFilter" dense class="text-grey" active-color="primary" indicator-color="primary" align="left">
        <q-tab name="" :label="t('common.all')" />
        <q-tab name="materia_alternativa" :label="t('religionAlternative.groupsTab')" />
        <q-tab name="uscita_scuola" :label="t('religionAlternative.optionD')" />
        <q-tab name="evaluations" :label="t('religionAlternative.evaluationsTab')" />
      </q-tabs>
      <q-separator />

      <div v-if="tabFilter !== 'evaluations'">
        <q-table
          :rows="filteredOptions"
          :columns="optionColumns"
          row-key="id"
          flat
          :loading="loading"
          :no-data-label="t('common.noData')"
        >
          <template #body-cell-option_type="props">
            <q-td :props="props">
              <q-badge :color="getOptionBadgeColor(props.row.option_type)" :label="formatOptionType(props.row.option_type)" />
            </q-td>
          </template>

          <template #body-cell-notes="props">
            <q-td :props="props">
              <span class="text-caption text-grey-8">{{ props.row.notes || '-' }}</span>
            </q-td>
          </template>
        </q-table>
      </div>

      <!-- Evaluations view -->
      <div v-else class="q-pa-md">
        <div class="row justify-between items-center q-mb-sm">
          <div class="text-subtitle2 text-weight-bold">{{ t('religionAlternative.judgmentLevels') }}</div>
          <q-btn color="secondary" size="sm" icon="add" :label="t('religionAlternative.judgment')" unelevated no-caps @click="openEvaluationDialog" />
        </div>
        <q-table
          :rows="evaluations"
          :columns="evalColumns"
          row-key="id"
          flat
          :loading="loading"
          :no-data-label="t('common.noData')"
        >
          <template #body-cell-judgment_level="props">
            <q-td :props="props">
              <q-chip :color="getJudgmentColor(props.row.judgment_level)" text-color="white" size="sm" class="text-weight-bold">
                {{ props.row.judgment_level.toUpperCase() }}
              </q-chip>
            </q-td>
          </template>
        </q-table>
      </div>
    </q-card>

    <!-- Save Option Dialog -->
    <q-dialog v-model="showOptionDialog" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ t('religionAlternative.title') }}</div>
        </q-card-section>
        <q-card-section class="q-pt-md">
          <q-input v-model="optionForm.student_id" :label="t('religionAlternative.student') + ' *'" outlined dense class="q-mb-sm" />
          <q-select
            v-model="optionForm.option_type"
            :options="optionChoices"
            emit-value
            map-options
            :label="t('religionAlternative.title') + ' *'"
            outlined
            dense
            class="q-mb-sm"
          />
          <q-input v-model="optionForm.notes" :label="t('common.notes')" outlined dense type="textarea" rows="2" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('common.save')" unelevated :loading="submitting" @click="saveOption" />
        </q-card-actions>
      </q-card>
    </q-dialog>


    <!-- Save Evaluation Dialog -->
    <q-dialog v-model="showEvaluationDialog" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="bg-secondary text-white">
          <div class="text-h6">{{ t('religionAlternative.evaluationsTab') }}</div>
        </q-card-section>
        <q-card-section class="q-pt-md">
          <q-input v-model="evalForm.student_id" :label="t('religionAlternative.student') + ' *'" outlined dense class="q-mb-sm" />
          <q-select
            v-model="evalForm.subject_kind"
            :options="evalSubjectChoices"
            emit-value
            map-options
            :label="t('common.details') + ' *'"
            outlined
            dense
            class="q-mb-sm"
          />
          <q-select
            v-model="evalForm.period"
            :options="periodChoices"
            emit-value
            map-options
            :label="t('common.period') + ' *'"
            outlined
            dense
            class="q-mb-sm"
          />
          <q-select
            v-model="evalForm.judgment_level"
            :options="judgmentChoices"
            emit-value
            map-options
            :label="t('religionAlternative.judgmentLevels') + ' *'"
            outlined
            dense
            class="q-mb-sm"
          />
          <q-input v-model="evalForm.descriptive_notes" :label="t('religionAlternative.notes')" outlined dense type="textarea" rows="2" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="secondary" :label="t('common.save')" unelevated :loading="submitting" @click="saveEvaluation" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { religionAlternativeService } from '@/services/religionAlternativeService'

export default defineComponent({
  name: 'ReligionAlternative',
  setup() {
    const { t } = useI18n()
    const options = ref([])
    const evaluations = ref([])
    const loading = ref(false)
    const submitting = ref(false)
    const tabFilter = ref('')
    const showOptionDialog = ref(false)
    const showEvaluationDialog = ref(false)

    const optionChoices = computed(() => [
      { label: t('religionAlternative.optionIrc'), value: 'irc' },
      { label: t('religionAlternative.optionA'), value: 'materia_alternativa' },
      { label: t('religionAlternative.optionB'), value: 'studio_assistito' },
      { label: t('religionAlternative.optionC'), value: 'studio_libero' },
      { label: t('religionAlternative.optionD'), value: 'uscita_scuola' }
    ])

    const evalSubjectChoices = computed(() => [
      { label: t('religionAlternative.optionIrc'), value: 'irc' },
      { label: t('religionAlternative.optionA'), value: 'materia_alternativa' }
    ])

    const periodChoices = computed(() => [
      { label: 'Q1 - ' + t('common.period') + ' 1', value: 'q1' },
      { label: 'Q2 - ' + t('common.period') + ' 2', value: 'q2' },
      { label: t('middleSchoolExam.finalHonor') || 'Finale', value: 'finale' }
    ])

    const judgmentChoices = computed(() => [
      { label: t('religionAlternative.excellent'), value: 'ottimo' },
      { label: t('religionAlternative.distinct'), value: 'distinto' },
      { label: t('religionAlternative.good'), value: 'buono' },
      { label: t('religionAlternative.sufficient'), value: 'sufficiente' },
      { label: t('religionAlternative.insufficient'), value: 'non_sufficiente' }
    ])

    const optionForm = ref({
      student_id: '',
      option_type: 'materia_alternativa',
      notes: ''
    })

    const evalForm = ref({
      student_id: '',
      subject_kind: 'materia_alternativa',
      period: 'q1',
      judgment_level: 'ottimo',
      descriptive_notes: ''
    })

    const optionColumns = computed(() => [
      { name: 'student_id', label: t('religionAlternative.student'), field: 'student_id', align: 'left' },
      { name: 'option_type', label: t('common.details'), field: 'option_type', align: 'left' },
      { name: 'notes', label: t('common.notes'), field: 'notes', align: 'left' },
      { name: 'chosen_at', label: t('common.history'), field: 'chosen_at', align: 'right', format: val => val ? new Date(val).toLocaleDateString() : '' }
    ])

    const evalColumns = computed(() => [
      { name: 'student_id', label: t('religionAlternative.student'), field: 'student_id', align: 'left' },
      { name: 'subject_kind', label: t('common.details'), field: 'subject_kind', align: 'left', format: val => val === 'irc' ? 'IRC' : t('religionAlternative.optionA') },
      { name: 'period', label: t('common.period'), field: 'period', align: 'center', format: val => val.toUpperCase() },
      { name: 'judgment_level', label: t('religionAlternative.judgment'), field: 'judgment_level', align: 'center' },
      { name: 'descriptive_notes', label: t('religionAlternative.notes'), field: 'descriptive_notes', align: 'left' }
    ])

    const countByOption = (type) => {
      return options.value.filter(o => o.option_type === type).length
    }

    const filteredOptions = computed(() => {
      if (!tabFilter.value) return options.value
      return options.value.filter(o => o.option_type === tabFilter.value)
    })

    const loadData = async () => {
      loading.value = true
      try {
        const optRes = await religionAlternativeService.getOptions()
        options.value = optRes.data || []
        const evalRes = await religionAlternativeService.getEvaluations()
        evaluations.value = evalRes.data || []
      } catch (e) {
        console.error('Errore caricamento scelte religione:', e)
      } finally {
        loading.value = false
      }
    }

    const openOptionDialog = () => {
      optionForm.value = {
        student_id: '',
        option_type: 'materia_alternativa',
        notes: ''
      }
      showOptionDialog.value = true
    }

    const saveOption = async () => {
      submitting.value = true
      try {
        await religionAlternativeService.saveOption(optionForm.value)
        showOptionDialog.value = false
        await loadData()
      } catch (e) {
        console.error('Errore salvataggio opzione:', e)
      } finally {
        submitting.value = false
      }
    }

    const openEvaluationDialog = () => {
      evalForm.value = {
        student_id: '',
        subject_kind: 'materia_alternativa',
        period: 'q1',
        judgment_level: 'ottimo',
        descriptive_notes: ''
      }
      showEvaluationDialog.value = true
    }

    const saveEvaluation = async () => {
      submitting.value = true
      try {
        await religionAlternativeService.saveEvaluation(evalForm.value)
        showEvaluationDialog.value = false
        await loadData()
      } catch (e) {
        console.error('Errore salvataggio giudizio:', e)
      } finally {
        submitting.value = false
      }
    }

    const formatOptionType = (opt) => {
      const match = optionChoices.value.find(c => c.value === opt)
      return match ? match.label : opt
    }

    const getOptionBadgeColor = (opt) => {
      switch (opt) {
        case 'irc': return 'primary'
        case 'materia_alternativa': return 'teal'
        case 'studio_assistito': return 'amber-9'
        case 'studio_libero': return 'orange-8'
        case 'uscita_scuola': return 'purple'
        default: return 'grey'
      }
    }

    const getJudgmentColor = (lvl) => {
      switch (lvl) {
        case 'ottimo': return 'positive'
        case 'distinto': return 'teal'
        case 'buono': return 'blue'
        case 'sufficiente': return 'amber-8'
        case 'non_sufficiente': return 'negative'
        default: return 'grey'
      }
    }

    onMounted(loadData)

    return {
      t,
      options,
      evaluations,
      loading,
      submitting,
      tabFilter,
      showOptionDialog,
      showEvaluationDialog,
      optionChoices,
      evalSubjectChoices,
      periodChoices,
      judgmentChoices,
      optionForm,
      evalForm,
      optionColumns,
      evalColumns,
      countByOption,
      filteredOptions,
      loadData,
      openOptionDialog,
      saveOption,
      openEvaluationDialog,
      saveEvaluation,
      formatOptionType,
      getOptionBadgeColor,
      getJudgmentColor
    }
  }
})
</script>
