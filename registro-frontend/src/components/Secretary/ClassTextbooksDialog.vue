<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(1100px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
      <q-card-section class="bg-gradient-premium text-white row items-center q-pa-md shrink-0">
        <div class="row items-center">
          <q-avatar color="white-20" text-color="white" icon="auto_stories" class="q-mr-sm" size="36px" />
          <div>
            <div class="text-h6 text-weight-bold">
              Adozioni Libri - Classe {{ targetClass?.name }}{{ targetClass?.section }}
            </div>
            <div class="text-subtitle2 opacity-80">{{ targetClass?.academic_year }}</div>
          </div>
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-pa-md col overflow-y-auto">
        <div v-if="loading" class="text-center q-pa-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>

        <div v-else class="row q-col-gutter-md">
          <!-- Left Column: Adopted Textbooks -->
          <div class="col-12 col-md-7">
            <q-table
              title="Libri Adottati"
              :rows="classTextbooks"
              :columns="textbookColumns"
              row-key="id"
              flat
              class="bg-transparent border-slate-100 rounded-xl"
              no-data-label="Nessun libro di testo ancora adottato per questa classe."
            >
              <template #header-cell="props">
                <q-th :props="props" class="text-slate-500 font-bold">
                  {{ props.col.label }}
                </q-th>
              </template>

              <template #body-cell-actions="props">
                <q-td :props="props" auto-width>
                  <q-btn
                    flat
                    round
                    dense
                    color="negative"
                    icon="delete"
                    :aria-label="t('common.delete') || 'Rimuovi libro adottato'"
                    @click="removeTextbook(props.row)"
                  >
                    <q-tooltip>Rimuovi adozione</q-tooltip>
                  </q-btn>
                </q-td>
              </template>
            </q-table>
          </div>

          <!-- Right Column: Adopt Textbook Form -->
          <div class="col-12 col-md-5">
            <q-card flat class="rounded-xl bg-slate-50 q-pa-md border-slate-200">
              <div class="row items-center justify-between q-mb-md">
                <div class="text-subtitle1 text-weight-bold text-slate-800">Adotta Libro</div>
                <q-btn
                  flat
                  dense
                  icon="add"
                  label="Nuovo Libro"
                  color="primary"
                  no-caps
                  :aria-label="'Crea nuovo libro a catalogo'"
                  @click="openCreateTextbook"
                />
              </div>

              <q-form @submit="addTextbookToClass" class="q-gutter-y-md">
                <q-select
                  v-model="textbookForm.textbook_id"
                  :options="allTextbooksOptions"
                  label="Libro *"
                  outlined
                  dense
                  emit-value
                  map-options
                  :rules="[val => !!val || 'Seleziona libro']"
                >
                  <template #no-option>
                    <q-item>
                      <q-item-section class="text-grey">Nessun libro in catalogo</q-item-section>
                    </q-item>
                    <q-item clickable @click="openCreateTextbook">
                      <q-item-section class="text-primary text-weight-bold">CREA NUOVO LIBRO</q-item-section>
                    </q-item>
                  </template>
                </q-select>

                <q-select
                  v-model="textbookForm.subject_id"
                  :options="subjectOptions"
                  label="Materia *"
                  outlined
                  dense
                  emit-value
                  map-options
                  :rules="[val => !!val || 'Seleziona materia']"
                />

                <q-checkbox v-model="textbookForm.is_optional" label="Il libro è opzionale" class="text-slate-700" />

                <q-btn
                  type="submit"
                  label="Conferma Adozione"
                  color="primary"
                  class="full-width rounded-lg q-py-sm shadow-sm q-mt-md"
                  no-caps
                  :loading="adopting"
                />
              </q-form>
            </q-card>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Create Textbook in Catalog Sub-Dialog -->
    <q-dialog v-model="showCreateTextbookDialog">
      <q-card style="display: flex; flex-direction: column; width: min(550px, 95vw); max-height: 90vh;" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="bg-primary text-white row items-center justify-between q-pa-md">
          <div class="text-subtitle1 font-bold">
            <q-icon name="menu_book" class="q-mr-xs" />
            Nuovo Libro Scolastico (Catalogo)
          </div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-md scroll" style="flex: 1; overflow-y: auto;">
          <q-form @submit="createTextbookInCatalog" class="q-gutter-y-md">
            <q-input v-model="newTextbook.title" label="Titolo del Libro *" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
            <q-input v-model="newTextbook.author" label="Autore / Autori *" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
            <q-input v-model="newTextbook.publisher" label="Casa Editrice *" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
            <q-input v-model="newTextbook.isbn" label="Codice ISBN *" outlined dense :rules="[val => !!val || 'Obbligatorio']" />

            <q-select
              v-model="newTextbook.subject_id"
              :options="subjectOptions"
              label="Materia associata"
              outlined
              dense
              emit-value
              map-options
              clearable
            />

            <div class="row q-col-gutter-md">
              <div class="col-6">
                <q-input v-model.number="newTextbook.price" label="Prezzo (€)" type="number" step="0.01" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model="newTextbook.volume" label="Volume (es. 1, Unico)" outlined dense />
              </div>
            </div>

            <div class="row justify-end q-gutter-sm q-mt-md">
              <q-btn flat label="Annulla" v-close-popup no-caps />
              <q-btn type="submit" label="Salva Libro" color="primary" class="rounded-lg q-px-md" no-caps :loading="savingTextbook" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import textbookService from '@/services/textbookService'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  targetClass: {
    type: Object,
    default: null
  },
  subjectOptions: {
    type: Array,
    default: () => []
  },
  schoolId: {
    type: [String, Number],
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'updated'])

const $q = useQuasar()
const { t } = useI18n()

const loading = ref(false)
const adopting = ref(false)
const showCreateTextbookDialog = ref(false)
const savingTextbook = ref(false)

const classTextbooks = ref([])
const allTextbooks = ref([])

const textbookForm = reactive({
  textbook_id: null,
  subject_id: null,
  is_optional: false
})

const newTextbook = reactive({
  title: '',
  author: '',
  publisher: '',
  isbn: '',
  subject_id: null,
  price: null,
  volume: ''
})

const textbookColumns = [
  { name: 'subject', label: 'Materia', field: 'subject_name', align: 'left' },
  { name: 'title', label: 'Titolo', field: 'title', align: 'left' },
  { name: 'author', label: 'Autore', field: 'author', align: 'left' },
  { name: 'optional', label: 'Opz.', field: row => row.is_optional ? 'Sì' : 'No', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

const allTextbooksOptions = computed(() => {
  return allTextbooks.value.map(b => ({
    label: `${b.title}${b.author ? ' - ' + b.author : ''}`,
    value: b.id
  }))
})

const fetchClassTextbooks = async () => {
  if (!props.targetClass?.id) return
  loading.value = true
  try {
    const res = await textbookService.listByClass(props.targetClass.id)
    classTextbooks.value = res.data || []
  } catch (e) {
    console.error('Error fetching class textbooks:', e)
    $q.notify({ type: 'negative', message: 'Errore caricamento libri' })
  } finally {
    loading.value = false
  }
}

const fetchAllTextbooks = async () => {
  try {
    const res = await textbookService.getAll()
    allTextbooks.value = res.data || []
  } catch (e) {
    console.error('Error fetching catalog textbooks:', e)
  }
}

watch(
  () => props.modelValue,
  async (newVal) => {
    if (newVal && props.targetClass?.id) {
      textbookForm.textbook_id = null
      textbookForm.subject_id = null
      textbookForm.is_optional = false
      await Promise.all([fetchClassTextbooks(), fetchAllTextbooks()])
    }
  },
  { immediate: true }
)

const addTextbookToClass = async () => {
  if (!props.targetClass?.id) return
  adopting.value = true
  try {
    await textbookService.assignToClass(props.targetClass.id, { ...textbookForm })
    $q.notify({ type: 'positive', message: 'Libro adottato con successo' })
    await fetchClassTextbooks()
    textbookForm.textbook_id = null
    emit('updated')
  } catch (e) {
    console.error('Error assigning textbook:', e)
    $q.notify({ type: 'negative', message: 'Errore durante l\'adozione del libro' })
  } finally {
    adopting.value = false
  }
}

const removeTextbook = async (row) => {
  try {
    await textbookService.removeFromClass(row.id)
    $q.notify({ type: 'positive', message: 'Adozione rimossa' })
    await fetchClassTextbooks()
    emit('updated')
  } catch (e) {
    console.error('Error removing textbook:', e)
    $q.notify({ type: 'negative', message: 'Errore durante la rimozione dell\'adozione' })
  }
}

const openCreateTextbook = () => {
  Object.assign(newTextbook, {
    title: '',
    author: '',
    publisher: '',
    isbn: '',
    subject_id: textbookForm.subject_id || null,
    price: null,
    volume: ''
  })
  showCreateTextbookDialog.value = true
}

const createTextbookInCatalog = async () => {
  savingTextbook.value = true
  try {
    const payload = {
      ...newTextbook,
      school_id: props.schoolId || undefined
    }
    const res = await textbookService.create(payload)
    $q.notify({ type: 'positive', message: 'Nuovo libro aggiunto al catalogo' })
    showCreateTextbookDialog.value = false
    await fetchAllTextbooks()
    const createdId = res.data?.id || res.data?.data?.id
    if (createdId) {
      textbookForm.textbook_id = createdId
    }
  } catch (e) {
    console.error('Error creating textbook:', e)
    $q.notify({ type: 'negative', message: 'Errore durante la creazione del libro a catalogo' })
  } finally {
    savingTextbook.value = false
  }
}
</script>
