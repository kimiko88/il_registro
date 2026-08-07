<template>
  <div class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h5 class="text-h5 text-weight-bold text-primary q-my-none">
          <q-icon name="menu_book" class="q-mr-sm" />
          Programmazione Didattica Annuale (UdA)
        </h5>
        <div class="text-caption text-grey-7">Pianificazione curricolare per materia, classe e competenze chiave</div>
      </div>
      <q-btn color="primary" icon="add" label="Nuova UdA" @click="openCreateDialog" class="glossy" />
    </div>

    <!-- Class Selector Filter -->
    <q-card flat bordered class="q-mb-md">
      <q-card-section class="row q-col-gutter-md items-center">
        <div class="col-12 col-md-6">
          <q-select
            v-model="selectedClassId"
            :options="classOptions"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            label="Seleziona Classe"
            outlined
            dense
            @update:model-value="fetchUdaPlans"
          />
        </div>
        <div class="col-12 col-md-6 text-right">
          <q-chip icon="insights" color="secondary" text-color="white" label="DM 742/2017 & Competenze Europee" />
        </div>
      </q-card-section>
    </q-card>

    <!-- UdA Cards List -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <div v-else-if="udaList.length === 0" class="text-center q-pa-xl text-grey-6">
      <q-icon name="assignment_late" size="48px" class="q-mb-sm" />
      <div>Nessuna Unità di Apprendimento (UdA) presente per questa classe.</div>
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="plan in udaList" :key="plan.id" class="col-12 col-md-6">
        <q-card flat bordered class="shadow-1 hover-shadow">
          <q-card-section>
            <div class="row items-center justify-between">
              <div class="text-subtitle1 text-weight-bold text-primary">{{ plan.title }}</div>
              <q-badge :color="getStatusColor(plan.status)" :label="plan.status.toUpperCase()" />
            </div>
            <div class="text-caption text-grey-7 q-mb-sm">
              Materia: <strong>{{ plan.subject_name || 'Generale' }}</strong> | Periodo: <strong>{{ plan.period }}</strong>
            </div>
            <div class="text-body2 text-grey-9 q-mb-md">
              {{ plan.description || 'Nessuna descrizione.' }}
            </div>

            <!-- Competencies Chips -->
            <div class="q-mb-sm">
              <div class="text-caption text-weight-bold">Competenze Target:</div>
              <div class="row q-gutter-xs q-mt-xs">
                <q-chip v-for="c in plan.competencies" :key="c" dense color="primary" outline size="sm">
                  {{ c }}
                </q-chip>
              </div>
            </div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="right">
            <q-btn flat dense icon="edit" color="secondary" label="Modifica" @click="openEditDialog(plan)" />
            <q-btn flat dense icon="delete" color="negative" label="Elimina" @click="deleteUda(plan.id)" />
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 500px; max-width: 700px">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6">{{ editMode ? 'Modifica UdA' : 'Nuova Unità di Apprendimento (UdA)' }}</div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-sm">
          <q-input v-model="form.title" label="Titolo dell'UdA *" outlined dense />
          <q-input v-model="form.description" label="Descrizione e Contenuti" type="textarea" outlined dense rows="3" />
          <q-select v-model="form.period" :options="['primo_quadrimestre', 'secondo_quadrimestre', 'annuale']" label="Periodo" outlined dense />

          <div class="text-caption text-weight-bold q-mt-sm">Competenze Chiave (seleziona o digita):</div>
          <q-select
            v-model="form.competencies"
            :options="['Competenza Alfabetica Funzionale', 'Competenza Multilinguistica', 'Competenza STEM', 'Competenza Digitale', 'Competenza Personale e Sociale', 'Competenza in Materia di Cittadinanza']"
            multiple
            use-chips
            outlined
            dense
          />

          <q-input v-model="form.objectives" label="Obiettivi di Apprendimento" type="textarea" outlined dense rows="2" />
          <q-input v-model="form.methodologies" label="Metodologie Didattiche (e.g. Flipped Classroom, Cooperative Learning)" outlined dense />
          <q-input v-model="form.evaluation_criteria" label="Criteri di Valutazione" outlined dense />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Salva UdA" @click="saveUda" :loading="saving" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import udaService from 'src/services/udaService'
import { useNotify } from 'src/composables/useNotify'

const notify = useNotify()

const selectedClassId = ref('47a05d80-3836-452e-ac91-8cfa3a1999dd')
const classOptions = ref([
  { id: '47a05d80-3836-452e-ac91-8cfa3a1999dd', name: 'Classe 2A' },
  { id: 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', name: 'Classe 3B' }
])

const loading = ref(false)
const saving = ref(false)
const udaList = ref([])

const showDialog = ref(false)
const editMode = ref(false)
const currentPlanId = ref(null)

const form = ref({
  class_id: '',
  subject_id: '26f22f7c-4c50-448b-8052-bdcb561953d4',
  title: '',
  description: '',
  period: 'primo_quadrimestre',
  competencies: [],
  objectives: '',
  methodologies: '',
  evaluation_criteria: ''
})

const fetchUdaPlans = async () => {
  if (!selectedClassId.value) return
  loading.value = true
  try {
    const data = await udaService.getByClass(selectedClassId.value)
    udaList.value = data || []
  } catch (err) {
    console.error('Error fetching UdA:', err)
  } finally {
    loading.value = false
  }
}

const getStatusColor = (status) => {
  if (status === 'approved') return 'positive'
  if (status === 'submitted') return 'warning'
  return 'grey-7'
}

const openCreateDialog = () => {
  editMode.value = false
  currentPlanId.value = null
  form.value = {
    class_id: selectedClassId.value,
    subject_id: '26f22f7c-4c50-448b-8052-bdcb561953d4',
    title: '',
    description: '',
    period: 'primo_quadrimestre',
    competencies: ['Competenza Digitale', 'Competenza STEM'],
    objectives: '',
    methodologies: 'Cooperative Learning',
    evaluation_criteria: 'Griglia di valutazione per competenze'
  }
  showDialog.value = true
}

const openEditDialog = (plan) => {
  editMode.value = true
  currentPlanId.value = plan.id
  form.value = { ...plan }
  showDialog.value = true
}

const saveUda = async () => {
  if (!form.value.title) {
    notify.error('Inserire il titolo dell\'UdA')
    return
  }
  saving.value = true
  try {
    form.value.class_id = selectedClassId.value
    if (editMode.value) {
      await udaService.update(currentPlanId.value, form.value)
      notify.success('UdA aggiornata con successo!')
    } else {
      await udaService.create(form.value)
      notify.success('UdA creata con successo!')
    }
    showDialog.value = false
    await fetchUdaPlans()
  } catch (err) {
    notify.error('Errore nel salvataggio dell\'UdA')
  } finally {
    saving.value = false
  }
}

const deleteUda = async (id) => {
  try {
    await udaService.delete(id)
    notify.success('UdA eliminata')
    await fetchUdaPlans()
  } catch (err) {
    notify.error('Errore durante l\'eliminazione')
  }
}

onMounted(() => {
  fetchUdaPlans()
})
</script>
