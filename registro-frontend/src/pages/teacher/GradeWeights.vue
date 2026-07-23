<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center q-mb-md">
      <q-icon name="balance" size="28px" color="primary" class="q-mr-sm" />
      <div class="text-h5 text-weight-bold">Configurazione Pesi Voti</div>
      <q-space />
      <q-btn icon="add" color="primary" label="Aggiungi Peso" @click="openDialog()" />
    </div>

    <q-card class="q-mb-md shadow-1">
      <q-card-section class="bg-blue-1">
        <q-icon name="info" color="primary" class="q-mr-sm" />
        <span class="text-caption">
          I pesi configurati qui vengono applicati come <strong>default</strong> nel calcolo delle medie ponderate.
          Se un voto ha un peso specifico impostato dal docente, quello ha la precedenza.
          <br>Scala: 0 = escluso dalla media, 1 = peso normale, 2 = doppio peso.
        </span>
      </q-card-section>
    </q-card>

    <!-- Filtri -->
    <q-card class="q-mb-md" bordered>
      <q-card-section class="row q-col-gutter-md">
        <div class="col-12 col-md-4">
          <q-select
            v-model="filterSubjectId"
            :options="subjectOptions"
            option-value="id"
            option-label="name"
            emit-value map-options
            label="Filtra per Materia"
            dense outlined clearable
            @update:model-value="loadConfigs"
          />
        </div>
        <div class="col-12 col-md-4">
          <q-select
            v-model="filterClassId"
            :options="classOptions"
            option-value="id"
            option-label="name"
            emit-value map-options
            label="Filtra per Classe"
            dense outlined clearable
            @update:model-value="loadConfigs"
          />
        </div>
      </q-card-section>
    </q-card>

    <!-- Loading -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="50px" />
    </div>

    <!-- Table -->
    <q-card v-else bordered>
      <q-table
        :rows="configs"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        no-data-label="Nessuna configurazione peso trovata. Clicca 'Aggiungi Peso' per iniziare."
      >
        <template #body-cell-grade_category="{ row }">
          <q-td>
            <q-badge :color="categoryColor(row.grade_category)" :label="categoryLabel(row.grade_category)" />
          </q-td>
        </template>
        <template #body-cell-evaluation_type="{ row }">
          <q-td>{{ row.evaluation_type || 'Tutti i tipi' }}</q-td>
        </template>
        <template #body-cell-weight="{ row }">
          <q-td>
            <q-chip :color="weightColor(row.weight)" text-color="white" dense>
              {{ row.weight }}×
            </q-chip>
          </q-td>
        </template>
        <template #body-cell-actions="{ row }">
          <q-td class="text-right">
            <q-btn flat round icon="edit" color="primary" size="sm" @click="openDialog(row)" />
            <q-btn flat round icon="delete" color="negative" size="sm" @click="deleteConfig(row.id)" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog -->
    <q-dialog v-model="dialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ editingId ? 'Modifica Peso' : 'Nuovo Peso' }}</div>
        </q-card-section>

        <q-card-section class="q-gutter-md">
          <q-select
            v-model="form.grade_category"
            :options="categoryOptions"
            option-value="value"
            option-label="label"
            emit-value map-options
            label="Categoria *"
            dense outlined
          />
          <q-select
            v-model="form.evaluation_type"
            :options="[{label: 'Tutti i tipi', value: null}, {label: 'Scritto', value: 'Written'}, {label: 'Orale', value: 'Oral'}, {label: 'Pratico', value: 'Practical'}]"
            option-value="value"
            option-label="label"
            emit-value map-options
            label="Tipo di Valutazione"
            dense outlined
          />
          <q-select
            v-model="form.subject_id"
            :options="subjectOptions"
            option-value="id"
            option-label="name"
            emit-value map-options
            label="Materia (opzionale)"
            dense outlined clearable
          />
          <q-select
            v-model="form.class_id"
            :options="classOptions"
            option-value="id"
            option-label="name"
            emit-value map-options
            label="Classe (opzionale)"
            dense outlined clearable
          />
          <q-input
            v-model.number="form.weight"
            type="number"
            label="Peso *"
            dense outlined
            :rules="[v => (v >= 0 && v <= 10) || 'Valore tra 0 e 10']"
            hint="0=escluso, 1=normale, 2=doppio peso"
            step="0.25" min="0" max="10"
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Salva" :loading="saving" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import api from 'src/services/api'

const $q = useQuasar()

const loading = ref(false)
const saving = ref(false)
const dialog = ref(false)
const editingId = ref(null)

const configs = ref([])
const subjectOptions = ref([])
const classOptions = ref([])

const filterSubjectId = ref(null)
const filterClassId = ref(null)

const defaultForm = () => ({
  grade_category: 'summative',
  evaluation_type: null,
  subject_id: null,
  class_id: null,
  weight: 1.0
})

const form = ref(defaultForm())

const columns = [
  { name: 'grade_category', label: 'Categoria', field: 'grade_category', align: 'left', sortable: true },
  { name: 'evaluation_type', label: 'Tipo Valutazione', field: 'evaluation_type', align: 'left' },
  { name: 'subject_id', label: 'Materia', field: 'subject_id', align: 'left',
    format: (v) => subjectOptions.value.find(s => s.id === v)?.name || (v ? v.slice(0, 8) + '…' : 'Tutte') },
  { name: 'class_id', label: 'Classe', field: 'class_id', align: 'left',
    format: (v) => classOptions.value.find(c => c.id === v)?.name || (v ? v.slice(0, 8) + '…' : 'Tutte') },
  { name: 'weight', label: 'Peso', field: 'weight', align: 'center', sortable: true },
  { name: 'actions', label: 'Azioni', field: 'actions', align: 'right' }
]

const categoryOptions = [
  { value: 'formative', label: 'Formativa (in itinere)' },
  { value: 'summative', label: 'Sommativa (verifiche finali)' },
  { value: 'practical', label: 'Pratica (laboratorio)' }
]

const categoryLabel = (cat) => categoryOptions.find(c => c.value === cat)?.label || cat
const categoryColor = (cat) => ({ formative: 'blue', summative: 'deep-purple', practical: 'teal' }[cat] || 'grey')
const weightColor = (w) => w === 0 ? 'grey' : w > 1.5 ? 'deep-orange' : w < 0.8 ? 'blue-grey' : 'primary'

onMounted(async () => {
  await Promise.all([loadConfigs(), loadSubjects(), loadClasses()])
})

async function loadConfigs() {
  loading.value = true
  try {
    const params = {}
    if (filterSubjectId.value) params.subject_id = filterSubjectId.value
    if (filterClassId.value) params.class_id = filterClassId.value
    const res = await api.get('/grades/weight-configs', { params })
    configs.value = res.data || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore caricamento configurazioni' })
  } finally {
    loading.value = false
  }
}

async function loadSubjects() {
  try {
    const res = await api.get('/subjects')
    subjectOptions.value = res.data || []
  } catch { /* ignore */ }
}

async function loadClasses() {
  try {
    const res = await api.get('/classes')
    classOptions.value = res.data || []
  } catch { /* ignore */ }
}

function openDialog(row = null) {
  if (row) {
    editingId.value = row.id
    form.value = {
      grade_category: row.grade_category,
      evaluation_type: row.evaluation_type || null,
      subject_id: row.subject_id || null,
      class_id: row.class_id || null,
      weight: row.weight
    }
  } else {
    editingId.value = null
    form.value = defaultForm()
  }
  dialog.value = true
}

async function save() {
  if (!form.value.grade_category || form.value.weight == null) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }
  if (form.value.weight < 0 || form.value.weight > 10) {
    $q.notify({ type: 'warning', message: 'Il peso deve essere compreso tra 0 e 10' })
    return
  }
  saving.value = true
  try {
    const payload = { ...form.value }
    if (!payload.evaluation_type) payload.evaluation_type = null
    if (!payload.subject_id) payload.subject_id = null
    if (!payload.class_id) payload.class_id = null

    await api.put('/grades/weight-configs', payload)
    $q.notify({ type: 'positive', message: 'Configurazione salvata' })
    dialog.value = false
    await loadConfigs()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore salvataggio' })
  } finally {
    saving.value = false
  }
}

async function deleteConfig(id) {
  $q.dialog({
    title: 'Conferma eliminazione',
    message: 'Eliminare questa configurazione peso?',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await api.delete(`/grades/weight-configs/${id}`)
      $q.notify({ type: 'positive', message: 'Configurazione eliminata' })
      await loadConfigs()
    } catch (e) {
      $q.notify({ type: 'negative', message: 'Errore eliminazione' })
    }
  })
}
</script>
