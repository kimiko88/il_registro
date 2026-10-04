<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-lg justify-between">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">{{ t('textbooksPage.title') }}</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ t('textbooksAie.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm">
        <q-btn
          v-if="currentTab === 'catalog'"
          color="primary"
          icon="add"
          :label="t('textbooksPage.newBook')"
          class="rounded-lg q-px-md shadow-soft"
          @click="openCreateDialog"
        />
        <q-btn
          v-if="currentTab === 'aie'"
          outline
          color="primary"
          icon="cloud_upload"
          :label="t('textbooksAie.importAie')"
          class="rounded-lg q-px-md"
          @click="showImportDialog = true"
        />
        <q-btn
          v-if="currentTab === 'aie' && selectedClass"
          color="secondary"
          icon="file_download"
          :label="t('textbooksAie.exportAie')"
          class="rounded-lg q-px-md"
          @click="exportAIE"
        />
      </div>
    </div>

    <!-- Navigation Tabs -->
    <q-tabs
      v-model="currentTab"
      dense
      class="text-grey-7 bg-white rounded-t-lg shadow-sm"
      active-color="primary"
      indicator-color="primary"
      align="left"
    >
      <q-tab name="catalog" icon="menu_book" :label="t('textbooksAie.catalogTab')" />
      <q-tab name="aie" icon="verified" :label="t('textbooksAie.aieTab')" />
    </q-tabs>

    <q-separator />

    <q-tab-panels v-model="currentTab" animated class="bg-transparent q-mt-md">
      <!-- PANEL 1: CATALOGO LIBRI ISTITUTO -->
      <q-tab-panel name="catalog" class="q-pa-none">
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
                  €{{ (props.value || 0).toFixed(2) }}
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
      </q-tab-panel>

      <!-- PANEL 2: ADOZIONI AIE & TETTO SPESA -->
      <q-tab-panel name="aie" class="q-pa-none">
        <!-- Class Selector & Controls -->
        <div class="row q-col-gutter-md q-mb-md items-center">
          <div class="col-12 col-md-4">
            <q-select
              v-model="selectedClass"
              :options="classList"
              option-value="id"
              option-label="name"
              emit-value
              map-options
              outlined
              dense
              bg-color="white"
              :label="t('textbooksAie.selectClass')"
              @update:model-value="loadClassAdoptions"
            >
              <template v-slot:prepend>
                <q-icon name="class" color="primary" />
              </template>
            </q-select>
          </div>
          <div class="col-12 col-md-8 text-right">
            <q-btn
              v-if="selectedClass"
              color="primary"
              icon="add"
              :label="t('textbooksAie.addAdoption')"
              dense
              class="q-px-md"
              @click="openAdoptionDialog"
            />
          </div>
        </div>

        <!-- SPENDING LIMIT MONITOR CARD -->
        <q-card v-if="selectedClass && spendingReport" class="q-mb-md glass-card shadow-soft border-slate-100 overflow-hidden">
          <q-card-section class="q-pa-md">
            <div class="row items-center justify-between q-col-gutter-md">
              <div class="col-12 col-md-4 row items-center no-wrap">
                <div class="q-mr-md relative-position">
                  <q-circular-progress
                    :value="Math.min(spendingReport.percentage || 0, 100)"
                    size="70px"
                    :thickness="0.22"
                    :color="getStatusColor(spendingReport.status)"
                    track-color="grey-3"
                    class="q-ma-none"
                  >
                    <span class="text-caption text-weight-bold">{{ Math.round(spendingReport.percentage || 0) }}%</span>
                  </q-circular-progress>
                </div>
                <div>
                  <div class="text-subtitle2 text-weight-bold text-slate-800">{{ t('textbooksAie.spendingControl') }}</div>
                  <q-badge :color="getStatusColor(spendingReport.status)" class="text-caption q-mt-xs q-py-xs q-px-sm">
                    {{ getStatusLabel(spendingReport.status) }}
                  </q-badge>
                </div>
              </div>

              <div class="col-12 col-md-8">
                <div class="row text-center q-col-gutter-sm">
                  <div class="col-3">
                    <div class="text-caption text-grey-6">{{ t('textbooksAie.actualSpending') }}</div>
                    <div class="text-weight-bold text-h6 text-slate-800">€{{ spendingReport.total_spending.toFixed(2) }}</div>
                  </div>
                  <div class="col-3">
                    <div class="text-caption text-grey-6">{{ t('textbooksAie.limit') }}</div>
                    <div class="text-weight-bold text-h6 text-primary">€{{ spendingReport.spending_limit.toFixed(2) }}</div>
                  </div>
                  <div class="col-3">
                    <div class="text-caption text-grey-6">{{ t('textbooksAie.tolerance') }}</div>
                    <div class="text-weight-bold text-h6 text-amber-9">€{{ spendingReport.tolerance_threshold.toFixed(2) }}</div>
                  </div>
                  <div class="col-3">
                    <div class="text-caption text-grey-6">{{ t('textbooksAie.deviation') }}</div>
                    <div :class="['text-weight-bold text-h6', spendingReport.difference > 0 ? 'text-negative' : 'text-positive']">
                      {{ spendingReport.difference > 0 ? '+' : '' }}€{{ spendingReport.difference.toFixed(2) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- CLASS ADOPTIONS TABLE -->
        <q-card v-if="selectedClass" class="glass-card shadow-soft border-slate-100 overflow-hidden">
          <q-table
            :rows="classAdoptions"
            :columns="adoptionColumns"
            row-key="id"
            flat
            :loading="adoptionsLoading"
            class="bg-transparent"
            :pagination="{ rowsPerPage: 15 }"
          >
            <template v-slot:body-cell-adoption_type="props">
              <q-td :props="props">
                <q-badge v-if="props.value === 'nuova_adozione'" color="blue-7" :label="t('textbooksAie.newAdoption')" />
                <q-badge v-else-if="props.value === 'scorrimento'" color="purple-7" :label="t('textbooksAie.scorrimento')" />
                <q-badge v-else-if="props.value === 'consigliato'" color="orange-8" :label="t('textbooksAie.recommended')" />
                <span v-else>{{ props.value }}</span>
              </q-td>
            </template>

            <template v-slot:body-cell-is_already_owned="props">
              <q-td :props="props">
                <q-chip v-if="props.value" color="green-1" text-color="green-9" dense icon="check_circle">
                  {{ t('textbooksAie.alreadyOwned') }}
                </q-chip>
                <span v-else class="text-grey-5">—</span>
              </q-td>
            </template>

            <template v-slot:body-cell-price="props">
              <q-td :props="props">
                <span class="text-weight-bold">€{{ (props.value || 0).toFixed(2) }}</span>
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props" auto-width>
                <q-btn flat round dense color="negative" icon="delete" @click="deleteAdoption(props.row.id)">
                  <q-tooltip>{{ t('textbooksAie.removeAdoption') }}</q-tooltip>
                </q-btn>
              </q-td>
            </template>

            <template v-slot:no-data>
              <div class="full-width column flex-center q-pa-xl text-slate-400">
                <q-icon name="menu_book" size="60px" class="opacity-20" />
                <div class="text-subtitle1 q-mt-sm">{{ t('textbooksAie.noAdoptions') }}</div>
              </div>
            </template>
          </q-table>
        </q-card>

        <div v-else class="text-center q-pa-xl text-grey-6 bg-white rounded-lg border border-dashed">
          <q-icon name="class" size="48px" class="text-grey-4" />
          <div class="text-h6 q-mt-sm">{{ t('textbooksAie.selectClass') }}</div>
        </div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Import AIE CSV -->
    <q-dialog v-model="showImportDialog">
      <q-card style="width: 500px; max-width: 90vw;" class="q-pa-md">
        <q-card-section class="row items-center">
          <div class="text-h6 text-weight-bold">{{ t('textbooksAie.importModalTitle') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <p class="text-caption text-grey-7">
            {{ t('textbooksAie.importNotes') }}
          </p>
          <q-file
            v-model="importFile"
            :label="t('textbooksAie.selectFile')"
            outlined
            accept=".csv, .txt"
          >
            <template v-slot:prepend>
              <q-icon name="attach_file" />
            </template>
          </q-file>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('textbooksAie.importAie')" :loading="importing" :disable="!importFile" @click="handleImportAIE" />
        </q-card-actions>
      </q-card>
    </q-dialog>


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
                <q-input 
                  v-model="form.author" 
                  :label="t('textbooksPage.colAuthor')" 
                  outlined 
                  :rules="[val => !!val || t('common.requiredField')]" 
                />
              </div>
              <div class="col-12 col-md-6">
                <q-input 
                  v-model="form.subject" 
                  :label="t('textbooksPage.colSubject')" 
                  outlined 
                  :rules="[val => !!val || t('common.requiredField')]" 
                />
              </div>
            </div>

            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <q-input 
                  v-model="form.isbn" 
                  :label="t('textbooksPage.colIsbn')" 
                  outlined 
                  :rules="[val => !!val || t('common.requiredField')]" 
                />
              </div>
              <div class="col-12 col-md-6">
                <q-input 
                  v-model="form.publisher" 
                  :label="t('textbooksPage.colPublisher')" 
                  outlined 
                  :rules="[val => !!val || t('common.requiredField')]" 
                />
              </div>
            </div>

            <q-input 
              v-model.number="form.price" 
              type="number" 
              step="0.01" 
              :label="t('textbooksPage.colPrice')" 
              outlined 
              prefix="€"
              :rules="[val => val >= 0 || t('textbooksPage.errNegativePrice')]" 
            />

            <div class="row justify-end q-gutter-md q-pt-md">
              <q-btn flat :label="t('common.cancel')" v-close-popup color="grey-7" />
              <q-btn 
                color="primary" 
                type="submit" 
                :label="t('common.save')" 
                :loading="saving" 
                class="rounded-lg q-px-lg shadow-soft" 
              />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { textbookService } from '@/services/textbookService'
import { useClassesStore } from '@/stores/classes'

const { t } = useI18n()
const $q = useQuasar()
const classesStore = useClassesStore()

const currentTab = ref('catalog')
const textbooks = ref([])
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const isEdit = ref(false)
const selectedId = ref(null)

// AIE State
const selectedClass = ref(null)
const classList = ref([])
const spendingReport = ref(null)
const classAdoptions = ref([])
const adoptionsLoading = ref(false)
const showImportDialog = ref(false)
const importFile = ref(null)
const importing = ref(false)

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

const adoptionColumns = computed(() => [
  { name: 'subject_name', label: t('common.category') || 'Materia', field: 'subject_name', align: 'left' },
  { name: 'book_title', label: t('textbooksPage.colTitle'), field: 'book_title', align: 'left' },
  { name: 'authors', label: t('textbooksPage.colAuthor'), field: 'authors', align: 'left' },
  { name: 'isbn', label: t('textbooksPage.colIsbn'), field: 'isbn', align: 'left' },
  { name: 'adoption_type', label: t('textbooksAie.newAdoption'), field: 'adoption_type', align: 'center' },
  { name: 'is_already_owned', label: t('textbooksAie.alreadyOwned'), field: 'is_already_owned', align: 'center' },
  { name: 'price', label: t('textbooksPage.colPrice'), field: 'price', align: 'right' },
  { name: 'actions', label: t('common.actions'), align: 'center' }
])

onMounted(async () => {
  await fetchTextbooks()
  await fetchClasses()
})

async function fetchClasses() {
  try {
    await classesStore.fetchClasses()
    classList.value = classesStore.classes || []
    if (classList.value.length > 0 && !selectedClass.value) {
      selectedClass.value = classList.value[0].id
      await loadClassAdoptions(selectedClass.value)
    }
  } catch (err) {
    console.error(err)
  }
}

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

async function loadClassAdoptions(classId) {
  if (!classId) return
  adoptionsLoading.value = true
  try {
    const [reportRes, adoptionsRes] = await Promise.all([
      textbookService.getSpendingReport(classId).catch(() => ({ data: null })),
      textbookService.listClassAdoptions(classId).catch(() => ({ data: [] }))
    ])
    spendingReport.value = reportRes.data
    classAdoptions.value = adoptionsRes.data || []
  } catch (err) {
    console.error(err)
  } finally {
    adoptionsLoading.value = false
  }
}

function getStatusColor(status) {
  if (status === 'WITHIN_LIMIT') return 'positive'
  if (status === 'WARNING_TOLERANCE') return 'warning'
  return 'negative'
}

function getStatusLabel(status) {
  if (status === 'WITHIN_LIMIT') return t('textbooksAie.statusWithinLimit')
  if (status === 'WARNING_TOLERANCE') return t('textbooksAie.statusWithinTolerance')
  return t('textbooksAie.statusExceeded')
}

async function handleImportAIE() {
  if (!importFile.value) return
  importing.value = true
  try {
    const fd = new FormData()
    fd.append('file', importFile.value)
    const res = await textbookService.importAIE(fd)
    $q.notify({ type: 'positive', message: res.data?.message || t('textbooksAie.importSuccess') })
    showImportDialog.value = false
    importFile.value = null
  } catch (err) {
    $q.notify({ type: 'negative', message: t('textbooksAie.importError') })
  } finally {
    importing.value = false
  }
}

async function exportAIE() {
  if (!selectedClass.value) return
  try {
    const res = await textbookService.exportClassAIE(selectedClass.value)
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `adozioni_aie_${selectedClass.value}.txt`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    $q.notify({ type: 'positive', message: t('textbooksAie.exportSuccess') })
  } catch (err) {
    $q.notify({ type: 'negative', message: t('textbooksAie.exportError') })
  }
}

async function deleteAdoption(id) {
  try {
    await textbookService.deleteAdoption(id)
    $q.notify({ type: 'positive', message: t('textbooksAie.adoptionDeleted') })
    if (selectedClass.value) {
      await loadClassAdoptions(selectedClass.value)
    }
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
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
