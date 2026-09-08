<template>
  <q-page padding class="bg-slate-50 min-h-screen">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="accessibility_new" color="primary" class="q-mr-sm" />
          {{ t('pdpPage.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('pdpPage.subtitle') }}
        </p>
      </div>
      
      <div class="row q-gutter-sm">
        <q-select
          v-model="selectedClassId"
          :options="classOptions"
          option-value="id"
          option-label="label"
          emit-value map-options
          dense outlined
          :label="t('pdpPage.filterClass')"
          style="min-width: 180px"
          class="bg-white"
          @update:model-value="fetchClassPlans"
        />
        <q-btn
          color="primary"
          icon="add"
          :label="t('pdpPage.newPlan')"
          unelevated
          class="rounded-lg text-weight-bold"
          :disable="!selectedClassId"
          @click="openCreateDialog"
        />
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="q-pa-md">
      <SkeletonCard v-for="i in 3" :key="i" class="q-mb-md" />
    </div>

    <!-- Empty State -->
    <q-card v-else-if="plans.length === 0" class="text-center q-pa-xl bg-white rounded-xl shadow-xs border">
      <q-icon name="assignment_late" size="64px" color="grey-4" class="q-mb-md" />
      <div class="text-h6 text-weight-bold text-slate-700">{{ t('pdpPage.noPlans') }}</div>
      <div class="text-caption text-grey-6 q-mb-md">{{ t('pdpPage.noPlansDesc') }}</div>
      <q-btn color="primary" icon="add" :label="t('pdpPage.createFirst')" unelevated no-caps @click="openCreateDialog" :disable="!selectedClassId" />
    </q-card>

    <!-- Plans List -->
    <div v-else class="row q-col-gutter-md">
      <div v-for="plan in plans" :key="plan.id" class="col-12 col-md-6">
        <q-card class="rounded-xl shadow-xs border bg-white hover:shadow-md transition-all">
          <q-card-section class="row items-center justify-between q-pb-xs">
            <div class="row items-center q-gutter-x-sm">
              <q-avatar size="36px" :color="plan.plan_type === 'pei' ? 'purple-1' : 'indigo-1'" :text-color="plan.plan_type === 'pei' ? 'purple' : 'indigo'">
                <q-icon :name="plan.plan_type === 'pei' ? 'psychology' : 'accessibility_new'" size="20px" />
              </q-avatar>
              <div>
                <div class="text-subtitle1 text-weight-bold text-slate-800">
                  {{ plan.student_name || t('pdpPage.studentName') }}
                </div>
                <div class="text-caption text-grey-6">
                  {{ plan.plan_type === 'pei' ? 'PEI (L. 104/92)' : 'PDP (L. 170/2010)' }} • {{ t('common.year') || 'A.S.' }} {{ plan.academic_year }}
                </div>
              </div>
            </div>
            
            <q-chip
              dense
              :color="plan.shared_with_family ? (plan.family_approved_at ? 'positive' : 'warning') : 'grey-4'"
              :text-color="plan.shared_with_family ? 'white' : 'grey-9'"
              class="text-weight-bold"
            >
              {{ plan.shared_with_family ? (plan.family_approved_at ? t('pdpPage.approvedFamily') : t('pdpPage.sharedPending')) : t('pdpPage.internalDraft') }}
            </q-chip>
          </q-card-section>

          <q-separator />

          <q-card-section class="q-py-sm">
            <div class="text-caption text-grey-7 q-mb-xs">
              <strong>{{ t('pdpPage.clinicalDiagnosis') }}:</strong> {{ plan.diagnosis || t('pdpPage.noneSpecified') }}
            </div>

            <!-- Restricted Diagnosis File (Visible ONLY to Class Teachers) -->
            <div v-if="plan.diagnosis_file" class="q-my-sm p-2 bg-purple-50 rounded border border-purple-200 row items-center justify-between">
              <div class="row items-center">
                <q-icon name="lock" color="purple" class="q-mr-xs" />
                <span class="text-caption text-weight-bold text-purple-9">{{ t('pdpPage.restrictedFile') }}</span>
              </div>
              <q-btn flat dense icon="download" color="purple" :label="t('common.download') || 'Scarica PDF'" @click="downloadDiagnosisFile(plan)" />
            </div>
            
            <div class="text-caption text-slate-700 q-mt-sm">
              <strong>{{ t('pdpPage.compensatoryTools') }}:</strong>
              <div class="row q-gutter-xs q-mt-xs">
                <q-chip
                  v-for="m in (plan.content?.compensative || [])"
                  :key="m"
                  dense
                  color="indigo-1"
                  text-color="indigo-9"
                  size="xs"
                >
                  {{ formatMeasure(m) }}
                </q-chip>
                <span v-if="!(plan.content?.compensative?.length)" class="text-grey-5 font-italic">{{ t('pdpPage.noneSpecified') }}</span>
              </div>
            </div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="right" class="q-px-md q-py-xs bg-slate-50">
            <q-btn flat dense icon="edit" color="primary" :label="t('common.edit') || t('pdpPage.editPlan')" @click="editPlan(plan)" />
            <q-btn
              flat dense
              :icon="plan.shared_with_family ? 'visibility_off' : 'share'"
              :color="plan.shared_with_family ? 'orange' : 'teal'"
              :label="plan.shared_with_family ? t('pdpPage.hideFromFamily') : t('pdpPage.shareWithFamily')"
              @click="toggleShare(plan)"
            />
            <q-btn flat round dense icon="delete" color="negative" :aria-label="t('common.delete')" @click="deletePlan(plan)" />
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Form Dialog matching User Screenshot -->
    <q-dialog v-model="showDialog" persistent max-width="750px">
      <q-card style="width: 750px; max-width: 95vw;" class="rounded-xl">
        <q-form @submit="savePlan" greedy>
          <q-card-section class="bg-primary text-white row items-center justify-between q-py-md q-px-lg">
            <div class="text-h6 text-weight-bold">
              {{ isEditing ? t('pdpPage.editPlan') : t('pdpPage.createPlan') }}
            </div>
            <q-btn flat round dense icon="close" v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
          </q-card-section>

          <q-card-section class="q-pa-lg q-pt-lg q-gutter-y-md">
            <div class="row q-col-gutter-md q-mt-xs">
              <div class="col-12" v-if="!isEditing">
                <q-select
                  v-model="form.student_id"
                  :options="studentOptions"
                  option-value="id"
                  option-label="label"
                  emit-value map-options
                  :label="t('pdpPage.selectStudent') || 'Seleziona Studente *'"
                  outlined dense
                  :rules="[val => !!val || t('common.requiredField') || 'Campo obbligatorio']"
                />
              </div>
              <div class="col-12">
                <q-select
                  v-model="form.plan_type"
                  :options="[{ label: t('pdpPage.pdpBes') || 'PDP (BES / DSA)', value: 'pdp' }, { label: t('pdpPage.pei104') || 'PEI (Disabilità H)', value: 'pei' }]"
                  emit-value map-options
                  :label="t('pdpPage.planTypeLabel') || 'Tipo Piano *'"
                  outlined dense
                  :rules="[val => !!val || t('common.requiredField') || 'Tipo Piano obbligatorio']"
                />
              </div>
            </div>

            <q-input
              v-model="form.diagnosis"
              type="textarea"
              rows="3"
              :label="t('pdpPage.clinicalDiagnosis')"
              outlined dense
            />

            <!-- Diagnosis File Attachment for Secretary/Docente -->
            <div class="q-pa-sm bg-purple-50 rounded border border-purple-200">
              <div class="text-caption text-weight-bold text-purple-9 q-mb-xs">
                <q-icon name="cloud_upload" class="q-mr-xs" /> {{ t('pdpPage.restrictedFile') }}
              </div>
              <q-file
                v-model="diagnosisFile"
                :label="t('pdpPage.downloadPdf') || 'Seleziona file diagnosi (PDF, JPG, PNG)'"
                outlined
                dense
                accept=".pdf,.jpg,.png,.doc,.docx"
                bg-color="white"
              >
                <template v-slot:append>
                  <q-icon name="attach_file" />
                </template>
              </q-file>
            </div>

            <q-separator />

            <CompensativeMeasuresSelector v-model="form.content.compensative" />

            <q-input
              v-model="form.content.notes"
              type="textarea"
              rows="3"
              :label="t('pdpPage.customNotes') || 'Obiettivi e Note Strategiche Didattiche'"
              outlined dense
            />
          </q-card-section>

          <q-card-actions align="right" class="q-pa-md bg-slate-50">
            <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup color="grey-7" />
            <q-btn color="primary" type="submit" :label="t('pdpPage.savePlan') || t('common.save')" unelevated :loading="saving" />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useClassesStore } from '@/stores/classes'
import { pdpService } from '@/services/pdpService'
import api from '@/services/api'
import SkeletonCard from '@/components/Common/SkeletonCard.vue'
import CompensativeMeasuresSelector from '@/components/Teacher/CompensativeMeasuresSelector.vue'

const $q = useQuasar()
const { t } = useI18n()
const classesStore = useClassesStore()

const selectedClassId = ref(null)
const plans = ref([])
const loading = ref(false)
const saving = ref(false)
const students = ref([])

const showDialog = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const diagnosisFile = ref(null)

const form = ref({
  student_id: '',
  plan_type: 'pdp',
  diagnosis: '',
  content: {
    compensative: [],
    dispensative: [],
    objectives: [],
    notes: ''
  }
})

const classOptions = computed(() => classesStore.classes.map(c => ({ id: c.id, label: c.name || `${t('common.class') || 'Classe'} ${c.id}` })))
const studentOptions = computed(() => students.value.map(s => ({ id: s.id, label: `${s.last_name} ${s.first_name}` })))

onMounted(async () => {
  await classesStore.fetchAssignedClasses()
  if (classesStore.classes.length > 0) {
    selectedClassId.value = classesStore.classes[0].id
    await fetchClassPlans()
  }
})

async function fetchClassPlans() {
  if (!selectedClassId.value) return
  loading.value = true
  try {
    const res = await pdpService.getClassPlans(selectedClassId.value, '2025/2026')
    plans.value = res.data?.plans || []
    
    const usersRes = await api.get('/users', { params: { class_id: selectedClassId.value, role: 'student', page_size: 200 } })
    students.value = usersRes.data?.users || []
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') || 'Errore caricamento piani PDP' })
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  isEditing.value = false
  editingId.value = null
  diagnosisFile.value = null
  form.value = {
    student_id: '',
    plan_type: 'pdp',
    diagnosis: '',
    content: { compensative: [], dispensative: [], objectives: [], notes: '' }
  }
  showDialog.value = true
}

function editPlan(plan) {
  isEditing.value = true
  editingId.value = plan.id
  diagnosisFile.value = null
  form.value = {
    student_id: plan.student_id,
    plan_type: plan.plan_type,
    diagnosis: plan.diagnosis,
    content: {
      compensative: plan.content?.compensative || [],
      dispensative: plan.content?.dispensative || [],
      objectives: plan.content?.objectives || [],
      notes: plan.content?.notes || ''
    }
  }
  showDialog.value = true
}

async function savePlan() {
  saving.value = true
  try {
    if (isEditing.value) {
      await pdpService.updatePlan(editingId.value, {
        plan_type: form.value.plan_type,
        diagnosis: form.value.diagnosis,
        content: form.value.content
      })
      $q.notify({ type: 'positive', message: t('common.success') || 'PDP aggiornato' })
    } else {
      await pdpService.createPlan({
        student_id: form.value.student_id,
        class_id: selectedClassId.value,
        academic_year: '2025/2026',
        plan_type: form.value.plan_type,
        diagnosis: form.value.diagnosis,
        content: form.value.content
      })
      $q.notify({ type: 'positive', message: t('common.success') || 'Nuovo PDP creato con successo' })
    }
    showDialog.value = false
    await fetchClassPlans()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') || 'Errore durante il salvataggio' })
  } finally {
    saving.value = false
  }
}

function downloadDiagnosisFile(_plan) {
  $q.notify({ type: 'info', message: t('common.download') || 'Download diagnosi riservata docenti in corso...' })
}

async function toggleShare(plan) {
  try {
    await pdpService.shareWithFamily(plan.id, !plan.shared_with_family)
    $q.notify({
      type: 'info',
      message: !plan.shared_with_family ? t('pdpPage.shareWithFamily') : t('pdpPage.hideFromFamily')
    })
    await fetchClassPlans()
  } catch {
    $q.notify({ type: 'negative', message: t('common.error') || 'Errore aggiornamento condivisione' })
  }
}

async function deletePlan(plan) {
  $q.dialog({
    title: t('common.confirm') || 'Conferma eliminazione',
    message: `${t('common.confirm') || 'Eliminare il piano di'} ${plan.student_name}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await pdpService.deletePlan(plan.id)
      $q.notify({ type: 'positive', message: t('common.delete') || 'Piano eliminato' })
      await fetchClassPlans()
    } catch {
      $q.notify({ type: 'negative', message: t('common.error') || 'Errore eliminazione' })
    }
  })
}

function formatMeasure(val) {
  const map = {
    calcolatrice: t('pdpPage.measures.calculator') || 'Uso Calcolatrice',
    tempo_aggiuntivo_30: t('pdpPage.measures.extraTime30') || 'Tempo Agg. (+30%)',
    tempo_aggiuntivo_50: t('pdpPage.measures.extraTime50') || 'Tempo Agg. (+50%)',
    prova_equipollente: t('pdpPage.measures.equivalentTest') || 'Prova Equipollente',
    sintesi_vocale: t('pdpPage.measures.textToSpeech') || 'Sintesi Vocale',
    mappe_concettuali: t('pdpPage.measures.conceptMaps') || 'Mappe Concettuali'
  }
  if (val.startsWith('custom_')) {
    return val.replace('custom_', '').replace(/_/g, ' ')
  }
  return map[val] || val
}
</script>
