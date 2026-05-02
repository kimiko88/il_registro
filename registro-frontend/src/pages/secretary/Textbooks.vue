<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center q-mb-md">
      <div class="text-h4 text-weight-bold">
        <q-icon name="auto_stories" color="primary" class="q-mr-sm" />
        Gestione Libri di Testo
      </div>
      <q-space />
      <q-btn color="primary" icon="add" label="Nuovo Libro" @click="openCreateDialog" />
    </div>

    <q-table
      :rows="textbooks"
      :columns="columns"
      row-key="id"
      flat bordered
      :loading="loading"
    >
      <template v-slot:body-cell-actions="props">
        <q-td :props="props" auto-width>
          <q-btn flat round color="negative" icon="delete" @click="confirmDelete(props.row)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Create Textbook -->
    <q-dialog v-model="showDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Nuovo Libro di Testo</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="saveTextbook" class="q-gutter-md">
            <q-input v-model="form.title" label="Titolo" outlined :rules="[val => !!val || 'Obbligatorio']" />
            <q-input v-model="form.author" label="Autore" outlined />
            <q-input v-model="form.isbn" label="ISBN" outlined />
            <q-input v-model="form.publisher" label="Editore" outlined />
            <q-input v-model.number="form.price" label="Prezzo" type="number" outlined />
            
            <div class="row justify-end q-mt-md">
              <q-btn flat label="Annulla" v-close-popup />
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
import { useQuasar } from 'quasar'
import { textbookService } from 'src/services/textbookService'

const $q = useQuasar()
const textbooks = ref([])
const loading = ref(false)
const showDialog = ref(false)

const form = reactive({
  title: '',
  author: '',
  isbn: '',
  publisher: '',
  price: 0
})

const columns = [
  { name: 'title', label: 'Titolo', field: 'title', align: 'left', sortable: true },
  { name: 'author', label: 'Autore', field: 'author', align: 'left' },
  { name: 'isbn', label: 'ISBN', field: 'isbn', align: 'left' },
  { name: 'publisher', label: 'Editore', field: 'publisher', align: 'left' },
  { name: 'price', label: 'Prezzo', field: 'price', format: val => `€${val.toFixed(2)}`, align: 'right' },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

onMounted(fetchTextbooks)

async function fetchTextbooks() {
  loading.value = true
  try {
    const res = await textbookService.getAll()
    textbooks.value = res.data || []
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  Object.assign(form, { title: '', author: '', isbn: '', publisher: '', price: 0 })
  showDialog.value = true
}

async function saveTextbook() {
  try {
    await textbookService.create(form)
    $q.notify({ type: 'positive', message: 'Libro creato con successo' })
    showDialog.value = false
    fetchTextbooks()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' })
  }
}

async function confirmDelete(row) {
  $q.dialog({
    title: 'Conferma',
    message: `Vuoi eliminare "${row.title}"?`,
    cancel: true
  }).onOk(async () => {
    await textbookService.delete(row.id)
    fetchTextbooks()
  })
}
</script>
