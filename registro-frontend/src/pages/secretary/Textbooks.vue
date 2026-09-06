<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-lg">
      <div class="col">
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">{{ t('textbooksPage.title') }}</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">{{ t('textbooksPage.subtitle') }}</p>
      </div>
      <div class="col-auto">
        <q-btn 
          color="primary" 
          icon="add" 
          :label="t('textbooksPage.newBook')" 
          class="rounded-lg q-px-md shadow-soft" 
          @click="openCreateDialog" 
        />
      </div>
    </div>

    <q-card class="glass-card shadow-soft border-slate-100 overflow-hidden">
      <q-table
        :rows="textbooks"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        class="bg-transparent"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template v-slot:header="props">
          <q-tr :props="props" class="bg-slate-50 text-slate-700">
            <q-th v-for="col in props.cols" :key="col.name" :props="props" class="text-weight-bold text-uppercase">
              {{ col.label }}
            </q-th>
          </q-tr>
        </template>

        <template v-slot:body-cell-title="props">
          <q-td :props="props">
            <div class="row items-center no-wrap">
              <q-avatar color="indigo-50" text-color="indigo-700" icon="book" size="32px" class="q-mr-sm" />
              <div class="text-weight-bold text-slate-800">{{ props.value }}</div>
            </div>
          </q-td>
        </template>

        <template v-slot:body-cell-price="props">
          <q-td :props="props">
            <q-chip outline color="primary" text-color="primary" dense class="text-weight-bold">
              €{{ props.value.toFixed(2) }}
            </q-chip>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" auto-width>
            <div class="row q-gutter-xs">
              <q-btn flat round dense color="primary" icon="edit" @click="openEditDialog(props.row)">
                <q-tooltip>{{ t('common.edit') }}</q-tooltip>
              </q-btn>
              <q-btn flat round dense color="negative" icon="delete" @click="confirmDelete(props.row)">
                <q-tooltip>{{ t('common.delete') }}</q-tooltip>
              </q-btn>
            </div>
          </q-td>
        </template>

        <template v-slot:no-data>
          <div class="full-width column flex-center q-pa-xl text-slate-400">
            <q-icon name="auto_stories" size="80px" class="opacity-20" />
            <div class="text-h6 q-mt-md">{{ t('textbooksPage.noDataTitle') }}</div>
            <p>{{ t('textbooksPage.noDataSubtitle') }}</p>
          </div>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Create/Edit Textbook -->
    <q-dialog v-model="showDialog" persistent class="premium-dialog">
      <q-card style="display: flex; flex-direction: column; width: 550px; max-width: 95vw; max-height: 90vh;" class="glass-card overflow-hidden bg-white">
        <q-card-section class="bg-gradient-primary text-white q-pa-lg row items-center">
          <div class="text-h5 text-weight-bold text-outfit">
            {{ isEdit ? t('textbooksPage.editBook') : t('textbooksPage.createTitle') }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-xl scroll" style="flex: 1; overflow-y: auto;">

          <q-form @submit="saveTextbook" class="q-gutter-y-lg">
            <q-input 
              v-model="form.title" 
              :label="t('textbooksPage.colTitle')" 
              outlined 
              :placeholder="t('textbooksPage.placeholderTitle')"
              :rules="[val => !!val || t('common.requiredField')]" 
            />
            
            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <q-input v-model="form.author" :label="t('textbooksPage.colAuthor')" outlined :placeholder="t('textbooksPage.placeholderAuthor')" />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.subject" :label="t('textbooksPage.colSubject')" outlined :placeholder="t('textbooksPage.placeholderSubject')" />
              </div>
            </div>

            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <q-input v-model="form.publisher" :label="t('textbooksPage.colPublisher')" outlined :placeholder="t('textbooksPage.placeholderPublisher')" />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.isbn" :label="t('textbooksPage.colIsbn')" outlined :placeholder="t('textbooksPage.placeholderIsbn')" />
              </div>
            </div>

            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-4">
                <q-input v-model.number="form.price" :label="t('textbooksPage.colPrice') + ' (€)'" type="number" step="0.01" outlined />
              </div>
            </div>
            
            <div class="row justify-end q-mt-xl q-gutter-sm">
              <q-btn flat :label="t('common.cancel')" v-close-popup class="rounded-lg" />
              <q-btn 
                type="submit" 
                :label="isEdit ? (t('common.update') || 'Aggiorna') : (t('common.create') || 'Crea Libro')" 
                color="primary" 
                class="q-px-xl rounded-lg shadow-sm" 
                :loading="saving"
              />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { textbookService } from '@/services/textbookService'

const $q = useQuasar()
const { t } = useI18n()
const textbooks = ref([])
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const isEdit = ref(false)
const selectedId = ref(null)

const form = reactive({
  title: '',
  author: '',
  subject: '',
  isbn: '',
  publisher: '',
  price: 0
})

const columns = computed(() => [
  { name: 'title', label: t('textbooksPage.colTitle'), field: 'title', align: 'left', sortable: true },
  { name: 'subject', label: t('textbooksPage.colSubject'), field: 'subject', align: 'left', sortable: true },
  { name: 'author', label: t('textbooksPage.colAuthor'), field: 'author', align: 'left', sortable: true },
  { name: 'isbn', label: t('textbooksPage.colIsbn'), field: 'isbn', align: 'left' },
  { name: 'publisher', label: t('textbooksPage.colPublisher'), field: 'publisher', align: 'left', sortable: true },
  { name: 'price', label: t('textbooksPage.colPrice'), field: 'price', align: 'right', sortable: true },
  { name: 'actions', label: t('textbooksPage.colActions'), align: 'center' }
])

onMounted(fetchTextbooks)

async function fetchTextbooks() {
  loading.value = true
  try {
    const res = await textbookService.getAll()
    textbooks.value = res.data || []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  isEdit.value = false
  selectedId.value = null
  Object.assign(form, { title: '', author: '', subject: '', isbn: '', publisher: '', price: 0 })
  showDialog.value = true
}

function openEditDialog(row) {
  isEdit.value = true
  selectedId.value = row.id
  Object.assign(form, { 
    title: row.title, 
    author: row.author,
    subject: row.subject || '',
    isbn: row.isbn, 
    publisher: row.publisher, 
    price: row.price 
  })
  showDialog.value = true
}

async function saveTextbook() {
  saving.value = true
  try {
    if (isEdit.value) {
      await textbookService.update(selectedId.value, form)
      $q.notify({ type: 'positive', message: t('textbooksPage.updateSuccess') })
    } else {
      await textbookService.create(form)
      $q.notify({ type: 'positive', message: t('textbooksPage.saveSuccess') })
    }
    showDialog.value = false
    fetchTextbooks()
  } catch (e) {
    $q.notify({ type: 'negative', message: t('textbooksPage.saveError') })
  } finally {
    saving.value = false
  }
}

async function confirmDelete(row) {
  $q.dialog({
    title: t('textbooksPage.deleteConfirmTitle'),
    message: t('textbooksPage.deleteConfirmMsg', { title: row.title }),
    cancel: true,
    persistent: true,
    ok: {
      color: 'negative',
      label: t('common.delete'),
      flat: false
    }
  }).onOk(async () => {
    try {
      await textbookService.delete(row.id)
      $q.notify({ type: 'positive', message: t('textbooksPage.deleteSuccess') })
      fetchTextbooks()
    } catch (err) {
      $q.notify({ type: 'negative', message: t('textbooksPage.deleteError') })
    }
  })
}
</script>

<style scoped>
.opacity-20 {
  opacity: 0.2;
}
</style>
