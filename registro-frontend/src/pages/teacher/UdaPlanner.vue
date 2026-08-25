<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg sticky-header">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="auto_stories" color="primary" class="q-mr-sm" />
          {{ t('udaPage.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('udaPage.subtitle') }}
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn
          color="primary"
          icon="add"
          :label="t('udaPage.newUda')"
          unelevated
          class="rounded-lg text-weight-bold"
          @click="openCreateModal"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="fetchUdaList" />
      </div>
    </div>

    <!-- Filters Bar -->
    <q-card flat bordered class="rounded-xl bg-white q-mb-lg shadow-soft q-pa-md">
      <div class="row q-col-gutter-md items-center">
        <div class="col-12 col-md-4">
          <q-select
            v-model="selectedClass"
            :options="classOptions"
            :label="t('udaPage.filterClass')"
            outlined dense emit-value map-options
            clearable
            @update:model-value="fetchUdaList"
          />
        </div>
        <div class="col-12 col-md-4">
          <q-select
            v-model="selectedSubject"
            :options="subjectOptions"
            :label="t('udaPage.filterSubject')"
            outlined dense emit-value map-options
            clearable
            @update:model-value="fetchUdaList"
          />
        </div>
        <div class="col-12 col-md-4">
          <q-input
            v-model="searchQuery"
            :placeholder="t('udaPage.searchPlaceholder')"
            outlined dense clearable
          >
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
      </div>
    </q-card>

    <!-- UdA Cards List -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="48px" />
    </div>

    <div v-else-if="filteredUdaList.length === 0" class="text-center q-pa-xl text-slate-400">
      <q-icon name="auto_stories" size="64px" class="q-mb-md opacity-40" />
      <div class="text-h6">{{ t('udaPage.noUdaFound') }}</div>
      <div class="text-caption q-mb-md">{{ t('udaPage.noUdaFoundDesc') }}</div>
      <q-btn color="primary" icon="add" :label="t('udaPage.createInDb')" unelevated @click="openCreateModal" />
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="uda in filteredUdaList" :key="uda.id" class="col-12 col-md-6 col-lg-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft full-height column justify-between hover-shadow transition-all">
          <q-card-section>
            <div class="row items-center justify-between q-mb-xs">
              <q-badge :color="getStatusColor(uda.status)" class="text-weight-bold q-px-sm q-py-xs">
                {{ getStatusLabel(uda.status) }}
              </q-badge>
              <div class="text-caption text-grey-6 text-weight-medium">
                {{ uda.hours || 10 }} {{ t('udaPage.totalHours') }}
              </div>
            </div>

            <div class="text-h6 text-weight-bold text-slate-800 q-mt-sm line-clamp-2">
              {{ uda.title }}
            </div>
            <div class="text-caption text-primary text-weight-bold q-mb-sm">
              {{ uda.subject_name || t('udaPage.subject') }} · {{ t('classRegister.classLabel', { name: uda.class_name || uda.class_id || 'N/D' }) }}
            </div>

            <div class="text-body2 text-slate-600 line-clamp-3 q-mb-md">
              {{ uda.description || t('udaPage.noDescription') }}
            </div>

            <div class="bg-slate-50 q-pa-sm rounded-lg border border-slate-200">
              <div class="text-caption text-weight-bold text-slate-700">{{ t('udaPage.targetCompetencies') }}:</div>
              <div class="text-caption text-slate-600 line-clamp-2">{{ Array.isArray(uda.competencies) ? uda.competencies.join(', ') : (uda.competencies || t('udaPage.curricularCompetencies')) }}</div>
            </div>
          </q-card-section>

          <q-card-actions class="bg-slate-50 border-t border-slate-200 justify-between q-px-md">
            <q-btn flat dense icon="visibility" color="primary" :label="t('udaPage.details')" @click="viewUda(uda)" />
            <div>
              <q-btn flat round dense icon="edit" color="secondary" @click="openEditModal(uda)" />
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDeleteUda(uda)" />
            </div>
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Create / Edit Dialog -->
    <q-dialog v-model="showModal" persistent max-width="750px">
      <q-card style="width: 750px; max-width: 95vw" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md q-px-lg">
          <div class="text-h6 text-weight-bold">
            <q-icon name="auto_stories" class="q-mr-xs" />
            {{ isEditing ? t('udaPage.editUda') : t('udaPage.createUda') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-form @submit.prevent="saveUda" class="q-pa-lg">
          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-12 col-md-8">
              <q-input
                v-model="form.title"
                :label="t('udaPage.titleLabel')"
                outlined dense hide-bottom-space
                :rules="[val => !!val || t('errors.ERR_REQUIRED_FIELDS')]"
              />
            </div>
            <div class="col-12 col-md-4">
              <q-select
                v-model="form.status"
                :options="[
                  { label: t('udaPage.draft'), value: 'draft' },
                  { label: t('udaPage.active'), value: 'active' },
                  { label: t('udaPage.completed'), value: 'completed' }
                ]"
                emit-value map-options
                :label="t('udaPage.statusLabel')"
                outlined dense hide-bottom-space
              />
            </div>
          </div>

          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-12 col-md-5">
              <q-select
                v-model="form.class_id"
                :options="classOptions"
                emit-value map-options
                :label="t('competenciesPage.selectClass')"
                outlined dense hide-bottom-space
                :rules="[val => !!val || t('errors.ERR_REQUIRED_FIELDS')]"
              />
            </div>
            <div class="col-12 col-md-5">
              <q-select
                v-model="form.subject_id"
                :options="subjectOptions"
                emit-value map-options
                :label="t('udaPage.subject') + ' *'"
                outlined dense hide-bottom-space
                :rules="[val => !!val || t('errors.ERR_REQUIRED_FIELDS')]"
              />
            </div>
            <div class="col-12 col-md-2">
              <q-input
                v-model.number="form.hours"
                type="number"
                :label="t('udaPage.duration')"
                outlined dense hide-bottom-space
              />
            </div>
          </div>

          <div class="q-mb-md">
            <q-input
              v-model="form.description"
              type="textarea"
              rows="3"
              :label="t('udaPage.descLabel')"
              outlined dense hide-bottom-space
            />
          </div>

          <div class="q-mb-md">
            <q-input
              v-model="form.competenciesText"
              type="textarea"
              rows="3"
              :label="t('udaPage.competenciesLabel')"
              outlined dense hide-bottom-space
            />
          </div>

          <div class="q-mb-sm">
            <q-input
              v-model="form.evaluation_criteria"
              type="textarea"
              rows="2"
              :label="t('udaPage.criteriaLabel')"
              outlined dense hide-bottom-space
            />
          </div>

          <q-card-actions align="right" class="q-pt-lg q-px-none q-pb-none bg-white">
            <q-btn flat :label="t('common.cancel')" v-close-popup class="text-weight-bold" />
            <q-btn color="primary" icon="save" :label="t('udaPage.saveUda')" unelevated :loading="submitting" type="submit" class="rounded-lg text-weight-bold" />
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>

    <!-- Detail Dialog -->
    <q-dialog v-model="showDetailModal">
      <q-card style="width: 650px; max-width: 95vw" class="rounded-xl" v-if="selectedUda">
        <q-card-section class="bg-slate-800 text-white row items-center justify-between q-py-md q-px-lg">
          <div>
            <div class="text-h6 text-weight-bold">{{ selectedUda.title }}</div>
            <div class="text-caption opacity-80">{{ selectedUda.subject_name }} · {{ t('classRegister.classLabel', { name: selectedUda.class_name }) }}</div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg q-gutter-y-md">
          <div>
            <div class="text-caption text-weight-bold text-slate-500 text-uppercase">{{ t('udaPage.descLabel') }}</div>
            <div class="text-body1 text-slate-800 q-mt-xs">{{ selectedUda.description || t('udaPage.noDescription') }}</div>
          </div>

          <q-separator />

          <div>
            <div class="text-caption text-weight-bold text-slate-500 text-uppercase">{{ t('udaPage.targetCompetencies') }}</div>
            <div class="text-body2 text-slate-700 q-mt-xs">{{ Array.isArray(selectedUda.competencies) ? selectedUda.competencies.join(', ') : selectedUda.competencies }}</div>
          </div>

          <q-separator />

          <div>
            <div class="text-caption text-weight-bold text-slate-500 text-uppercase">{{ t('udaPage.evaluationCriteria') }}</div>
            <div class="text-body2 text-slate-700 q-mt-xs">{{ selectedUda.evaluation_criteria || 'Criteri standard' }}</div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="t('common.close')" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from 'src/services/api'

const $q = useQuasar()
const { t } = useI18n()
const udaList = ref([])
const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const showDetailModal = ref(false)
const isEditing = ref(false)
const selectedUda = ref(null)

const selectedClass = ref(null)
const selectedSubject = ref(null)
const searchQuery = ref('')

const classOptions = ref([])
const subjectOptions = ref([])

const form = ref({
  id: null,
  title: '',
  class_id: null,
  subject_id: null,
  hours: 15,
  status: 'active',
  description: '',
  competenciesText: '',
  evaluation_criteria: ''
})

const filteredUdaList = computed(() => {
  return udaList.value.filter(uda => {
    if (selectedClass.value && uda.class_id !== selectedClass.value) return false
    if (selectedSubject.value && uda.subject_id !== selectedSubject.value) return false
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      return (uda.title || '').toLowerCase().includes(q) || (uda.description || '').toLowerCase().includes(q)
    }
    return true
  })
})

const getStatusLabel = (s) => ({ draft: 'Bozza', active: 'In Corso', completed: 'Completata' }[s] || 'In Corso')
const getStatusColor = (s) => ({ draft: 'amber-8', active: 'positive', completed: 'grey-7' }[s] || 'positive')

const fetchOptions = async () => {
  try {
    const [cRes, sRes] = await Promise.all([
      api.get('/classes'),
      api.get('/subjects')
    ])
    const classes = cRes.data?.classes || cRes.data || []
    classOptions.value = classes.map(c => ({ label: `Classe ${c.name || c.section}`, value: c.id }))
    const subjects = sRes.data?.subjects || sRes.data || []
    subjectOptions.value = subjects.map(s => ({ label: s.name, value: s.id }))
  } catch {
    /* ignore options fetch errors */
  }
}

const fetchUdaList = async () => {
  loading.value = true
  try {
    const url = selectedClass.value ? `/uda/class/${selectedClass.value}` : '/uda'
    const res = await api.get(url)
    udaList.value = res.data?.uda || res.data || []
  } catch (err) {
    udaList.value = []
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditing.value = false
  form.value = {
    id: null,
    title: '',
    class_id: classOptions.value.length ? classOptions.value[0].value : null,
    subject_id: subjectOptions.value.length ? subjectOptions.value[0].value : null,
    hours: 15,
    status: 'active',
    description: '',
    competenciesText: '',
    evaluation_criteria: ''
  }
  showModal.value = true
}

const openEditModal = (uda) => {
  isEditing.value = true
  form.value = {
    ...uda,
    competenciesText: Array.isArray(uda.competencies) ? uda.competencies.join(', ') : (uda.competencies || '')
  }
  showModal.value = true
}

const viewUda = (uda) => {
  selectedUda.value = uda
  showDetailModal.value = true
}

const saveUda = async () => {
  if (!form.value.title) {
    $q.notify({ type: 'warning', message: 'Inserire il titolo della UdA' })
    return
  }
  if (!form.value.class_id || !form.value.subject_id) {
    $q.notify({ type: 'warning', message: 'Selezionare classe e materia' })
    return
  }

  submitting.value = true
  try {
    const compArray = form.value.competenciesText ? form.value.competenciesText.split(',').map(s => s.trim()) : []
    const payload = {
      ...form.value,
      competencies: compArray
    }
    if (isEditing.value) {
      await api.put(`/uda/${form.value.id}`, payload)
    } else {
      await api.post('/uda', payload)
    }
    $q.notify({ type: 'positive', message: 'Uda salvata nel database con successo' })
    showModal.value = false
    await fetchUdaList()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore durante il salvataggio della UdA' })
  } finally {
    submitting.value = false
  }
}

const confirmDeleteUda = (uda) => {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: `Sei sicuro di voler eliminare l'UdA "${uda.title}" dal database?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await api.delete(`/uda/${uda.id}`)
      $q.notify({ type: 'positive', message: 'UdA eliminata dal database' })
      await fetchUdaList()
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Errore eliminazione UdA' })
    }
  })
}

onMounted(async () => {
  await fetchOptions()
  await fetchUdaList()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.hover-shadow:hover {
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
}
</style>
