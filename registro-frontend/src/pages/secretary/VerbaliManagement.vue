<template>
  <q-page class="q-pa-md q-gutter-y-md">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-xs">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none text-primary">
          <q-icon name="account_balance" class="q-mr-sm" />{{ t('verbaliManagement.title') }}
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none">
          {{ t('verbaliManagement.subtitle') }}
        </p>
      </div>
      <div class="q-gutter-sm">
        <q-btn
          color="primary"
          icon="add_circle"
          :label="t('verbaliManagement.newTemplateBtn')"
          unelevated
          @click="openCreateTemplateDialog"
        />
        <q-btn
          color="secondary"
          icon="event"
          :label="t('verbaliManagement.newMeetingBtn')"
          unelevated
          @click="openCreateMeetingDialog"
        />
      </div>
    </div>

    <!-- Banner Istituzionale di Riservatezza -->
    <q-banner rounded class="bg-indigo-1 text-indigo-10 q-mb-md">
      <template v-slot:avatar>
        <q-icon name="policy" color="primary" size="md" />
      </template>
      <div class="text-weight-bold text-body2">
        {{ t('verbaliManagement.privacyBannerTitle') }}
      </div>
      <div class="text-caption text-grey-8">
        {{ t('verbaliManagement.privacyBannerText') }}
      </div>
    </q-banner>

    <!-- Tabs Principali -->
    <q-card flat bordered class="rounded-borders">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey-7 bg-grey-1"
        active-color="primary"
        indicator-color="primary"
        align="left"
      >
        <q-tab name="verbali" icon="verified" :label="t('verbaliManagement.tabVerbali')" />
        <q-tab name="templates" icon="dashboard_customize" :label="t('verbaliManagement.tabTemplates')" />
        <q-tab name="meetings" icon="calendar_month" :label="t('verbaliManagement.tabMeetings')" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated>
        <!-- Tab 1: Verbali Ufficiali Firmati -->
        <q-tab-panel name="verbali" class="q-pa-md">
          <div class="row items-center justify-between q-mb-md">
            <div class="text-subtitle1 text-weight-bold text-grey-8">
              {{ t('verbaliManagement.archiveTitle', { count: verbaliList.length }) }}
            </div>
            <q-input
              v-model="searchVerbali"
              dense
              outlined
              :placeholder="t('verbaliManagement.searchPlaceholder')"
              class="q-w-md"
            >
              <template v-slot:append>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>

          <q-table
            :rows="filteredVerbali"
            :columns="verbaliColumns"
            row-key="id"
            :loading="loadingVerbali"
            flat
            bordered
            :no-data-label="t('verbaliManagement.noVerbali')"
          >
            <!-- Badge Stato -->
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge color="positive" class="q-pa-xs text-weight-bold">
                  <q-icon name="lock" class="q-mr-xs" /> {{ t('verbaliManagement.statusSigned') }}
                </q-badge>
              </q-td>
            </template>

            <!-- Firme Apposte -->
            <template v-slot:body-cell-signatures="props">
              <q-td :props="props">
                <q-btn
                  flat
                  dense
                  size="sm"
                  color="primary"
                  icon="history_edu"
                  :label="t('verbaliManagement.viewSignaturesBtn')"
                  @click="openViewSignatures(props.row)"
                />
              </q-td>
            </template>

            <!-- Azioni -->
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" class="q-gutter-xs">
                <q-btn
                  flat
                  round
                  dense
                  color="primary"
                  icon="visibility"
                  @click="openViewVerbale(props.row)"
                >
                  <q-tooltip>{{ t('verbaliManagement.readVerbaleTooltip') }}</q-tooltip>
                </q-btn>
                <q-btn
                  flat
                  round
                  dense
                  color="negative"
                  icon="picture_as_pdf"
                  @click="downloadPdf(props.row)"
                >
                  <q-tooltip>{{ t('verbaliManagement.downloadPdfTooltip') }}</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- Tab 2: Gestione Modelli & Ordini del Giorno (Dirigenza) -->
        <q-tab-panel name="templates" class="q-pa-md">
          <div class="row items-center justify-between q-mb-md">
            <div>
              <div class="text-subtitle1 text-weight-bold text-grey-8">
                {{ t('verbaliManagement.templatesTitle') }}
              </div>
              <div class="text-caption text-grey-6">
                {{ t('verbaliManagement.templatesSubtitle') }}
              </div>
            </div>
            <q-btn
              color="primary"
              icon="add"
              :label="t('verbaliManagement.createTemplateBtn')"
              unelevated
              @click="openCreateTemplateDialog"
            />
          </div>

          <q-table
            :rows="templatesList"
            :columns="templatesColumns"
            row-key="id"
            :loading="loadingTemplates"
            flat
            bordered
            :no-data-label="t('verbaliManagement.noTemplates')"
          >
            <!-- Badge Tipologia -->
            <template v-slot:body-cell-meeting_type="props">
              <q-td :props="props">
                <q-badge :color="getTypeColor(props.row.meeting_type)">
                  {{ formatMeetingType(props.row.meeting_type) }}
                </q-badge>
              </q-td>
            </template>

            <!-- Anteprima ODG -->
            <template v-slot:body-cell-agenda="props">
              <q-td :props="props">
                <div class="ellipsis-2-lines text-caption" style="max-width: 300px;">
                  {{ props.row.default_agenda }}
                </div>
              </q-td>
            </template>

            <!-- Azioni Modello -->
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" class="q-gutter-xs">
                <q-btn
                  flat
                  round
                  dense
                  color="primary"
                  icon="edit"
                  @click="openEditTemplate(props.row)"
                >
                  <q-tooltip>{{ t('verbaliManagement.editTemplateTooltip') }}</q-tooltip>
                </q-btn>
                <q-btn
                  flat
                  round
                  dense
                  color="negative"
                  icon="delete"
                  @click="deleteTemplate(props.row)"
                >
                  <q-tooltip>{{ t('verbaliManagement.deleteTemplateTooltip') }}</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- Tab 3: Calendario Riunioni d'Istituto -->
        <q-tab-panel name="meetings" class="q-pa-md">
          <div class="row items-center justify-between q-mb-md">
            <div class="text-subtitle1 text-weight-bold text-grey-8">
              {{ t('verbaliManagement.meetingsTitle') }}
            </div>
            <q-btn
              color="secondary"
              icon="add"
              :label="t('verbaliManagement.scheduleMeetingBtn')"
              unelevated
              @click="openCreateMeetingDialog"
            />
          </div>

          <q-table
            :rows="meetingsList"
            :columns="meetingsColumns"
            row-key="id"
            :loading="loadingMeetings"
            flat
            bordered
            :no-data-label="t('verbaliManagement.noMeetings')"
          >
            <template v-slot:body-cell-meeting_type="props">
              <q-td :props="props">
                <q-badge :color="getTypeColor(props.row.meeting_type)">
                  {{ formatMeetingType(props.row.meeting_type) }}
                </q-badge>
              </q-td>
            </template>

            <template v-slot:body-cell-agenda="props">
              <q-td :props="props">
                <div class="ellipsis-2-lines text-body2" style="max-width: 320px;">
                  {{ props.row.agenda || t('verbaliManagement.noAgendaSpecified') }}
                </div>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog 1: Creazione / Modifica Modello ODG (Dirigenza) -->
    <q-dialog v-model="showTemplateDialog">
      <q-card style="width: min(750px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">
            <q-icon :name="isEditingTemplate ? 'edit' : 'add_circle'" class="q-mr-sm" />
            {{ isEditingTemplate ? t('verbaliManagement.dialogEditTitle') : t('verbaliManagement.dialogNewTitle') }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-form @submit="submitTemplate" class="q-gutter-y-sm">
            <q-input
              v-model="templateForm.title"
              :label="t('verbaliManagement.formTitleLabel')"
              outlined
              dense
              :rules="[val => !!val || t('common.requiredField')]"
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="templateForm.meeting_type"
                  :options="meetingTypeOptions"
                  emit-value
                  map-options
                  :label="t('verbaliManagement.formTypeLabel')"
                  outlined
                  dense
                />
              </div>
              <div class="col-12 col-md-6">
                <q-input
                  v-model="templateForm.description"
                  :label="t('verbaliManagement.formDescLabel')"
                  outlined
                  dense
                />
              </div>
            </div>

            <q-input
              v-model="templateForm.default_agenda"
              type="textarea"
              :label="t('verbaliManagement.formAgendaLabel')"
              outlined
              dense
              rows="5"
              :hint="t('verbaliManagement.formAgendaHint')"
              :rules="[val => !!val || t('common.requiredField')]"
            />

            <q-input
              v-model="templateForm.template_content"
              type="textarea"
              :label="t('verbaliManagement.formContentLabel')"
              outlined
              dense
              rows="8"
              :hint="t('verbaliManagement.formContentHint')"
              :rules="[val => !!val || t('common.requiredField')]"
            />

            <div class="row justify-end q-mt-md q-gutter-sm">
              <q-btn :label="t('common.cancel')" flat v-close-popup />
              <q-btn :label="t('verbaliManagement.saveTemplateBtn')" color="primary" type="submit" unelevated :loading="submitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog 2: Creazione Riunione Istituzionale -->
    <q-dialog v-model="showMeetingDialog">
      <q-card style="width: min(650px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">
            <q-icon name="event_available" class="q-mr-sm" />{{ t('verbaliManagement.dialogMeetingTitle') }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-form @submit="submitMeeting" class="q-gutter-y-sm">
            <q-input
              v-model="meetingForm.title"
              :label="t('verbaliManagement.formMeetingTitle')"
              outlined
              dense
              :rules="[val => !!val || t('common.requiredField')]"
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-4">
                <q-input
                  v-model="meetingForm.date"
                  type="date"
                  :label="t('verbaliManagement.formDate')"
                  outlined
                  dense
                  :rules="[val => !!val || t('common.requiredField')]"
                />
              </div>
              <div class="col-6 col-md-4">
                <q-input
                  v-model="meetingForm.start_time"
                  type="time"
                  :label="t('verbaliManagement.formStartTime')"
                  outlined
                  dense
                  :rules="[val => !!val || t('common.requiredField')]"
                />
              </div>
              <div class="col-6 col-md-4">
                <q-input
                  v-model="meetingForm.end_time"
                  type="time"
                  :label="t('verbaliManagement.formEndTime')"
                  outlined
                  dense
                  :rules="[val => !!val || t('common.requiredField')]"
                />
              </div>
            </div>

            <q-select
              v-model="meetingForm.meeting_type"
              :options="meetingTypeOptions"
              emit-value
              map-options
              :label="t('verbaliManagement.formMeetingType')"
              outlined
              dense
            />

            <q-input
              v-model="meetingForm.agenda"
              type="textarea"
              :label="t('verbaliManagement.formAgenda')"
              outlined
              dense
              rows="4"
              :rules="[val => !!val || t('common.requiredField')]"
            />

            <div class="row justify-end q-mt-md q-gutter-sm">
              <q-btn :label="t('common.cancel')" flat v-close-popup />
              <q-btn :label="t('verbaliManagement.createMeetingBtn')" color="primary" type="submit" unelevated :loading="submitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog 3: Visualizzazione Verbale Ufficiale & Firme IP -->
    <q-dialog v-model="showViewDialog">
      <q-card style="width: min(700px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold">{{ currentVerbale?.title }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div class="q-mb-md">
            <q-badge color="positive" class="q-pa-xs text-subtitle2">
              <q-icon name="verified" class="q-mr-xs" /> {{ t('verbaliManagement.viewOfficialBadge') }}
            </q-badge>
          </div>

          <q-card flat bordered class="bg-grey-1 q-pa-md q-mb-md">
            <div class="text-subtitle2 text-grey-8 q-mb-xs">{{ t('verbaliManagement.viewVerbaleText') }}</div>
            <div style="white-space: pre-wrap; font-family: monospace;" class="text-body2">
              {{ currentVerbale?.content }}
            </div>
          </q-card>

          <!-- Firme Apposte -->
          <div class="text-subtitle1 text-weight-bold q-mb-sm">
            <q-icon name="history_edu" class="q-mr-xs text-primary" />{{ t('verbaliManagement.viewSignaturesHeader') }}
          </div>

          <q-list bordered separator class="rounded-borders">
            <q-item v-for="sig in verbaleSignatures" :key="sig.id">
              <q-item-section avatar>
                <q-avatar icon="verified_user" color="positive" text-color="white" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ sig.user_name || t('verbaliManagement.defaultSigner') }}</q-item-label>
                <q-item-label caption>
                  {{ t('verbaliManagement.signedAt', { date: new Date(sig.signed_at).toLocaleString(locale) }) }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-badge color="grey-8">{{ t('verbaliManagement.ipLabel', { ip: sig.ip_address || t('verbaliManagement.ipRegistered') }) }}</q-badge>
              </q-item-section>
            </q-item>
          </q-list>

          <div class="row justify-end q-mt-md q-gutter-sm">
            <q-btn
              color="negative"
              icon="picture_as_pdf"
              :label="t('verbaliManagement.downloadPdfBtn')"
              unelevated
              @click="downloadPdf(currentVerbale)"
            />
            <q-btn :label="t('common.close')" flat v-close-popup />
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import verbaliService from '@/services/verbaliService'

const $q = useQuasar()
const { t, locale } = useI18n()

// Tabs
const activeTab = ref('verbali')
const verbaliList = ref([])
const templatesList = ref([])
const meetingsList = ref([])

const searchVerbali = ref('')
const loadingVerbali = ref(false)
const loadingTemplates = ref(false)
const loadingMeetings = ref(false)
const submitting = ref(false)

// Modali
const showTemplateDialog = ref(false)
const showMeetingDialog = ref(false)
const showViewDialog = ref(false)

const isEditingTemplate = ref(false)
const currentVerbale = ref(null)
const verbaleSignatures = ref([])

// Form Modello
const templateForm = ref({
  id: '',
  title: '',
  meeting_type: 'consiglio_classe',
  description: '',
  default_agenda: '',
  template_content: ''
})

// Form Riunione
const meetingForm = ref({
  title: '',
  date: new Date().toISOString().split('T')[0],
  start_time: '15:00',
  end_time: '16:30',
  meeting_type: 'collegio_docenti',
  agenda: ''
})

// Opzioni tipologia riunione
const meetingTypeOptions = computed(() => [
  { label: t('verbali.types.consiglio_classe'), value: 'consiglio_classe' },
  { label: t('verbali.types.collegio_docenti'), value: 'collegio_docenti' },
  { label: t('verbali.types.dipartimento'), value: 'dipartimento' },
  { label: t('verbali.types.generale'), value: 'generale' }
])

// Colonne Verbali Ufficiali
const verbaliColumns = computed(() => [
  { name: 'title', label: t('verbaliManagement.colTitle'), field: 'title', align: 'left', sortable: true },
  { name: 'status', label: t('verbaliManagement.colStatus'), field: 'status', align: 'center' },
  { name: 'created_at', label: t('verbaliManagement.colArchivedDate'), field: r => new Date(r.created_at).toLocaleDateString(locale.value), align: 'center', sortable: true },
  { name: 'signatures', label: t('verbaliManagement.colSignatures'), align: 'center' },
  { name: 'actions', label: t('common.actions'), align: 'center' }
])

// Colonne Modelli
const templatesColumns = computed(() => [
  { name: 'title', label: t('verbaliManagement.colTemplateTitle'), field: 'title', align: 'left', sortable: true },
  { name: 'meeting_type', label: t('verbaliManagement.colMeetingType'), field: 'meeting_type', align: 'center' },
  { name: 'agenda', label: t('verbaliManagement.colAgendaPreview'), align: 'left' },
  { name: 'created_at', label: t('verbaliManagement.colCreatedAt'), field: r => new Date(r.created_at).toLocaleDateString(locale.value), align: 'center' },
  { name: 'actions', label: t('common.actions'), align: 'center' }
])

// Colonne Riunioni
const meetingsColumns = computed(() => [
  { name: 'title', label: t('verbaliManagement.colMeetingTitle'), field: 'title', align: 'left', sortable: true },
  { name: 'meeting_type', label: t('verbaliManagement.colMeetingType'), field: 'meeting_type', align: 'center' },
  { name: 'date', label: t('verbaliManagement.colDate'), field: r => new Date(r.date).toLocaleDateString(locale.value), align: 'center', sortable: true },
  { name: 'time', label: t('verbaliManagement.colTime'), field: r => `${r.start_time} - ${r.end_time}`, align: 'center' },
  { name: 'agenda', label: t('verbaliManagement.colAgenda'), align: 'left' }
])

// Filtro Verbali
const filteredVerbali = computed(() => {
  if (!searchVerbali.value) return verbaliList.value
  const q = searchVerbali.value.toLowerCase()
  return verbaliList.value.filter(v => (v.title || '').toLowerCase().includes(q))
})

// Helpers visivi
const formatMeetingType = (type) => {
  switch (type) {
    case 'consiglio_classe': return t('verbali.types.consiglio_classe')
    case 'collegio_docenti': return t('verbali.types.collegio_docenti')
    case 'dipartimento': return t('verbali.types.dipartimento')
    default: return t('verbali.types.generale')
  }
}

const getTypeColor = (type) => {
  switch (type) {
    case 'collegio_docenti': return 'purple-8'
    case 'dipartimento': return 'teal-8'
    default: return 'indigo-8'
  }
}

// Caricamento Dati
const loadVerbali = async () => {
  loadingVerbali.value = true
  try {
    const res = await verbaliService.getAllVerbali()
    verbaliList.value = res.data || []
  } catch {
    $q.notify({ type: 'negative', message: t('verbaliManagement.notifyLoadVerbaliError') })
  } finally {
    loadingVerbali.value = false
  }
}

const loadTemplates = async () => {
  loadingTemplates.value = true
  try {
    const res = await verbaliService.getTemplates()
    templatesList.value = res.data || []
  } catch {
    $q.notify({ type: 'negative', message: t('verbaliManagement.notifyLoadTemplatesError') })
  } finally {
    loadingTemplates.value = false
  }
}

const loadMeetings = async () => {
  loadingMeetings.value = true
  try {
    const res = await verbaliService.listMeetings()
    meetingsList.value = res.data || []
  } catch {
    $q.notify({ type: 'negative', message: t('verbaliManagement.notifyLoadMeetingsError') })
  } finally {
    loadingMeetings.value = false
  }
}

// Gestione Modelli
const openCreateTemplateDialog = () => {
  isEditingTemplate.value = false
  templateForm.value = {
    id: '',
    title: '',
    meeting_type: 'consiglio_classe',
    description: '',
    default_agenda: '',
    template_content: ''
  }
  showTemplateDialog.value = true
}

const openEditTemplate = (tpl) => {
  isEditingTemplate.value = true
  templateForm.value = {
    id: tpl.id,
    title: tpl.title,
    meeting_type: tpl.meeting_type,
    description: tpl.description,
    default_agenda: tpl.default_agenda,
    template_content: tpl.template_content
  }
  showTemplateDialog.value = true
}

const submitTemplate = async () => {
  submitting.value = true
  try {
    if (isEditingTemplate.value) {
      await verbaliService.updateTemplate(templateForm.value.id, templateForm.value)
      $q.notify({ type: 'positive', message: t('verbaliManagement.notifyTemplateUpdated') })
    } else {
      await verbaliService.createTemplate(templateForm.value)
      $q.notify({ type: 'positive', message: t('verbaliManagement.notifyTemplateCreated') })
    }
    showTemplateDialog.value = false
    await loadTemplates()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('verbaliManagement.notifyTemplateSaveError') })
  } finally {
    submitting.value = false
  }
}

const deleteTemplate = (tpl) => {
  $q.dialog({
    title: t('verbaliManagement.deleteDialogTitle'),
    message: t('verbaliManagement.deleteDialogMessage', { title: tpl.title }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await verbaliService.deleteTemplate(tpl.id)
      $q.notify({ type: 'positive', message: t('verbaliManagement.notifyTemplateDeleted') })
      await loadTemplates()
    } catch {
      $q.notify({ type: 'negative', message: t('verbaliManagement.notifyTemplateDeleteError') })
    }
  })
}

// Gestione Riunioni
const openCreateMeetingDialog = () => {
  meetingForm.value = {
    title: '',
    date: new Date().toISOString().split('T')[0],
    start_time: '15:00',
    end_time: '16:30',
    meeting_type: 'collegio_docenti',
    agenda: ''
  }
  showMeetingDialog.value = true
}

const submitMeeting = async () => {
  submitting.value = true
  try {
    await verbaliService.createMeeting(meetingForm.value)
    $q.notify({ type: 'positive', message: t('verbaliManagement.notifyMeetingCreated') })
    showMeetingDialog.value = false
    await loadMeetings()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('verbaliManagement.notifyMeetingCreateError') })
  } finally {
    submitting.value = false
  }
}

// Visualizzazione e Download PDF
const openViewVerbale = async (verbale) => {
  currentVerbale.value = verbale
  try {
    const res = await verbaliService.getSignatures(verbale.id)
    verbaleSignatures.value = res.data || []
  } catch {
    verbaleSignatures.value = []
  }
  showViewDialog.value = true
}

const openViewSignatures = async (verbale) => {
  await openViewVerbale(verbale)
}

const downloadPdf = async (verbale) => {
  try {
    const res = await verbaliService.exportPdf(verbale.id)
    const blob = new Blob([res.data], { type: 'application/pdf' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `verbale_${verbale.id.substring(0, 8)}.pdf`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
  } catch {
    $q.notify({ type: 'negative', message: t('verbaliManagement.notifyPdfDownloadError') })
  }
}

onMounted(() => {
  loadVerbali()
  loadTemplates()
  loadMeetings()
})
</script>
