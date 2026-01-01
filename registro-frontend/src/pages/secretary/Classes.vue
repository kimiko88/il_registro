<template>
  <q-page class="q-pa-md">
    <q-table
      title="Gestione Classi"
      :rows="classesStore.classes"
      :columns="columns"
      row-key="id"
      :filter="filter"
      :loading="classesStore.loading"
      :pagination.sync="pagination"
    >
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filter" placeholder="Cerca">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
        <q-btn color="primary" icon="add" label="Nuova Classe" class="q-ml-md" @click="openDialog()" />
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn flat round color="primary" icon="edit" @click="openDialog(props.row)" />
          <q-btn flat round color="negative" icon="delete" @click="confirmDelete(props.row)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Create/Edit -->
    <q-dialog v-model="showDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">{{ isEdit ? 'Modifica Classe' : 'Nuova Classe' }}</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="saveClass" class="q-gutter-md">
            <q-input v-model="form.name" label="Nome (es. 1, 5, I, V)" outlined :rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model="form.section" label="Sezione (es. A, B)" outlined :rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model="form.academic_year" label="Anno Accademico (es. 2024/2025)" outlined :rules="[val => !!val || 'Campo obbligatorio']" />
            <!-- Coordinator Selection could go here if we fetch teachers -->
            
            <div class="row justify-end">
              <q-btn flat label="Annulla" color="primary" v-close-popup />
              <q-btn type="submit" label="Salva" color="primary" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useClassesStore } from '@/stores/classes'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const classesStore = useClassesStore()

const filter = ref('')
const showDialog = ref(false)
const isEdit = ref(false)
const pagination = ref({ rowsPerPage: 10 })

const form = reactive({
  id: null,
  name: '',
  section: '',
  academic_year: '2024/2025', // Default
  coordinator_id: ''
})

const columns = [
  { name: 'name', label: 'Nome', align: 'left', field: 'name', sortable: true },
  { name: 'section', label: 'Sezione', align: 'left', field: 'section', sortable: true },
  { name: 'academic_year', label: 'Anno Accademico', align: 'center', field: 'academic_year', sortable: true },
  { name: 'coordinator', label: 'Coordinatore', align: 'left', field: row => row.coordinator_id || 'N/A' }, // Can enhance with Teacher name
  { name: 'actions', label: 'Azioni', align: 'right' }
]

onMounted(() => {
  classesStore.fetchClasses()
})

const openDialog = (row = null) => {
  if (row) {
    isEdit.value = true
    form.id = row.id
    form.name = row.name
    form.section = row.section
    form.academic_year = row.academic_year
    form.coordinator_id = row.coordinator_id
  } else {
    isEdit.value = false
    form.id = null
    form.name = ''
    form.section = ''
    form.academic_year = '2024/2025'
    form.coordinator_id = ''
  }
  showDialog.value = true
}

const saveClass = async () => {
  try {
    if (isEdit.value) {
      await classesStore.updateClass(form.id, { ...form })
      $q.notify({ type: 'positive', message: 'Classe aggiornata' })
    } else {
      await classesStore.createClass({ ...form })
      $q.notify({ type: 'positive', message: 'Classe creata' })
    }
    showDialog.value = false
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel salvataggio' })
  }
}

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Conferma',
    message: `Vuoi eliminare la classe ${row.name}${row.section}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await classesStore.deleteClass(row.id)
      $q.notify({ type: 'positive', message: 'Classe eliminata' })
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Errore nell\'eliminazione' })
    }
  })
}
</script>
