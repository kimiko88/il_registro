<template>
  <q-page class="q-pa-md q-gutter-y-md">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-xs">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none text-primary">
          <q-icon name="gavel" class="q-mr-sm" />{{ t('verbali.title') }}
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none">
          {{ t('verbali.subtitle') }}
        </p>
      </div>
      <div class="q-gutter-sm">
        <q-btn
          color="primary"
          icon="add"
          :label="t('verbali.newMeetingVerbaleBtn')"
          unelevated
          @click="openCreateMeetingDialog"
        />
      </div>
    </div>

    <!-- Banner informativo sul ciclo di vita del verbale -->
    <q-banner rounded class="bg-blue-1 text-blue-10 q-mb-md">
      <template v-slot:avatar>
        <q-icon name="security" color="primary" size="md" />
      </template>
      <div class="text-weight-bold text-body2">
        {{ t('verbali.lifecycleBannerTitle') }}
      </div>
      <div class="text-caption text-grey-8">
        {{ t('verbali.lifecycleBannerText') }}
      </div>
    </q-banner>

    <!-- Schede: Verbali vs Riunioni -->
    <q-card flat bordered class="rounded-borders">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey-7 bg-grey-1"
        active-color="primary"
        indicator-color="primary"
        align="left"
      >
        <q-tab name="verbali" icon="description" :label="t('verbali.tabVerbali')" />
        <q-tab name="meetings" icon="event" :label="t('verbali.tabMeetings')" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated>
        <!-- Tab 1: Verbali -->
        <q-tab-panel name="verbali" class="q-pa-none">
          <q-table
            :rows="verbaliList"
            :columns="verbaliColumns"
            row-key="id"
            :loading="loadingVerbali"
            flat
            :no-data-label="t('verbali.noVerbali')"
          >
            <!-- Badge Stato / Ciclo di vita -->
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge
                  v-if="props.row.is_signed"
                  color="positive"
                  text-color="white"
                  class="q-pa-xs text-weight-bold"
                >
                  <q-icon name="verified" class="q-mr-xs" /> {{ t('verbali.statusSigned') }}
                </q-badge>
                <q-badge
                  v-else
                  color="orange-9"
                  text-color="white"
                  class="q-pa-xs text-weight-bold"
                >
                  <q-icon name="edit_note" class="q-mr-xs" /> {{ t('verbali.statusDraft') }}
                </q-badge>
              </q-td>
            </template>

            <!-- Permesso di Modifica -->
            <template v-slot:body-cell-editing="props">
              <q-td :props="props">
                <span v-if="props.row.is_signed" class="text-grey-7 text-caption">
                  <q-icon name="lock" class="q-mr-xs" /> {{ t('verbali.lockedForAll') }}
                </span>
                <q-badge
                  v-else-if="props.row.can_edit"
                  color="indigo-7"
                  text-color="white"
                  class="text-caption"
                >
                  <q-icon name="edit" class="q-mr-xs" /> {{ t('verbali.youCanEdit') }}
                </q-badge>
                <span v-else class="text-grey-6 text-caption">
                  <q-icon name="visibility" class="q-mr-xs" /> {{ t('verbali.readOnlyReserved') }}
                </span>
              </q-td>
            </template>

            <!-- Firme / Audit -->
            <template v-slot:body-cell-signatures="props">
              <q-td :props="props">
                <q-btn
                  flat
                  dense
                  size="sm"
                  color="primary"
                  icon="history_edu"
                  :label="t('verbali.viewSignatures')"
                  @click="viewSignatures(props.row)"
                />
              </q-td>
            </template>

            <!-- Azioni -->
            <template v-slot:body-cell-actions="props">
              <q-td :props="props" class="q-gutter-xs">
                <!-- Visualizza -->
                <q-btn
                  flat
                  round
                  dense
                  color="grey-8"
                  icon="visibility"
                  @click="openViewVerbale(props.row)"
                >
                  <q-tooltip>{{ t('verbali.viewDocTooltip') }}</q-tooltip>
                </q-btn>

                <!-- Modifica Bozza (Abilitato solo se non firmato e utente è redattore autorizzato) -->
                <q-btn
                  v-if="!props.row.is_signed && props.row.can_edit"
                  flat
                  round
                  dense
                  color="primary"
                  icon="edit"
                  @click="openEditVerbale(props.row)"
                >
                  <q-tooltip>{{ t('verbali.editDraftTooltip') }}</q-tooltip>
                </q-btn>

                <!-- Firma Digitale con IP (Se non ancora firmato da questo utente) -->
                <q-btn
                  v-if="!props.row.is_signed && !props.row.is_signed_by_me"
                  color="secondary"
                  size="sm"
                  unelevated
                  icon="draw"
                  :label="t('verbali.signIpBtn')"
                  @click="signVerbale(props.row.id)"
                >
                  <q-tooltip>{{ t('verbali.signIpTooltip') }}</q-tooltip>
                </q-btn>

                <!-- Download PDF Ufficiale -->
                <q-btn
                  flat
                  round
                  dense
                  color="negative"
                  icon="picture_as_pdf"
                  @click="downloadPdf(props.row)"
                >
                  <q-tooltip>{{ t('verbali.downloadPdfTooltip') }}</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- Tab 2: Riunioni & Convocazioni (ODG) -->
        <q-tab-panel name="meetings" class="q-pa-none">
          <q-table
            :rows="meetingsList"
            :columns="meetingsColumns"
            row-key="id"
            :loading="loadingMeetings"
            flat
            :no-data-label="t('verbali.noMeetings')"
          >
            <template v-slot:body-cell-agenda="props">
              <q-td :props="props">
                <div class="ellipsis-2-lines text-body2" style="max-width: 320px;">
                  {{ props.row.agenda || t('verbali.noAgenda') }}
                </div>
                <q-btn
                  v-if="props.row.agenda"
                  flat
                  dense
                  size="xs"
                  color="primary"
                  :label="t('verbali.expandAgenda')"
                  @click="showAgendaDetail(props.row)"
                />
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props">
                <q-btn
                  color="primary"
                  size="sm"
                  unelevated
                  icon="post_add"
                  :label="t('verbali.writeVerbaleBtn')"
                  @click="openCreateVerbaleForMeeting(props.row)"
                />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog 1: Nuova Riunione (con selezione Modello ODG della Dirigente) -->
    <q-dialog v-model="showMeetingDialog">
      <q-card style="width: min(650px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">
            <q-icon name="event_available" class="q-mr-sm" />{{ t('verbali.dialogMeetingTitle') }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <!-- Selezione Modello Predefinito della Dirigente -->
          <div class="q-mb-md">
            <q-select
              v-model="selectedTemplate"
              :options="templatesList"
              option-label="title"
              :label="t('verbali.selectTemplateLabel')"
              outlined
              dense
              clearable
              @update:model-value="onTemplateSelected"
            >
              <template v-slot:prepend>
                <q-icon name="auto_awesome" color="amber-9" />
              </template>
              <template v-slot:option="scope">
                <q-item v-bind="scope.itemProps">
                  <q-item-section>
                    <q-item-label class="text-weight-bold">{{ scope.opt.title }}</q-item-label>
                    <q-item-label caption>{{ scope.opt.description }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-badge color="indigo">{{ formatMeetingType(scope.opt.meeting_type) }}</q-badge>
                  </q-item-section>
                </q-item>
              </template>
            </q-select>
          </div>

          <q-form @submit="submitMeeting" class="q-gutter-y-sm">
            <q-input
              v-model="meetingForm.title"
              :label="t('verbali.formMeetingTitle')"
              outlined
              dense
              :rules="[val => !!val || t('common.requiredField')]"
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-4">
                <q-input
                  v-model="meetingForm.date"
                  type="date"
                  :label="t('verbali.formDate')"
                  outlined
                  dense
                  :rules="[val => !!val || t('common.requiredField')]"
                />
              </div>
              <div class="col-6 col-md-4">
                <q-input
                  v-model="meetingForm.start_time"
                  type="time"
                  :label="t('verbali.formStartTime')"
                  outlined
                  dense
                  :rules="[val => !!val || t('common.requiredField')]"
                />
              </div>
              <div class="col-6 col-md-4">
                <q-input
                  v-model="meetingForm.end_time"
                  type="time"
                  :label="t('verbali.formEndTime')"
                  outlined
                  dense
                  :rules="[val => !!val || t('common.requiredField')]"
                />
              </div>
            </div>

            <q-input
              v-model="meetingForm.class_id"
              :label="t('verbali.formClassId')"
              outlined
              dense
              :hint="t('verbali.formClassIdHint')"
            />

            <q-input
              v-model="meetingForm.agenda"
              type="textarea"
              :label="t('verbali.formAgenda')"
              outlined
              dense
              rows="4"
              :rules="[val => !!val || t('verbali.formAgendaHint')]"
            />

            <div class="row justify-end q-mt-md q-gutter-sm">
              <q-btn :label="t('common.cancel')" flat v-close-popup />
              <q-btn :label="t('verbali.saveMeetingBtn')" color="primary" type="submit" unelevated :loading="submitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog 2: Redazione o Modifica Bozza Verbale -->
    <q-dialog v-model="showVerbaleDialog">
      <q-card style="width: min(750px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">
            <q-icon :name="isEditingVerbale ? 'edit' : 'post_add'" class="q-mr-sm" />
            {{ isEditingVerbale ? t('verbali.dialogEditVerbaleTitle') : t('verbali.dialogNewVerbaleTitle') }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div class="bg-amber-1 text-amber-10 q-pa-sm rounded-borders q-mb-md text-caption">
            <q-icon name="info" class="q-mr-xs" />
            {{ t('verbali.draftWarning') }}
          </div>

          <q-form @submit="submitVerbale" class="q-gutter-y-sm">
            <q-input
              v-model="verbaleForm.title"
              :label="t('verbali.formVerbaleTitle')"
              outlined
              dense
              :rules="[val => !!val || t('common.requiredField')]"
            />

            <q-input
              v-model="verbaleForm.content"
              type="textarea"
              :label="t('verbali.formVerbaleContent')"
              outlined
              dense
              rows="10"
              :rules="[val => !!val || t('verbali.formVerbaleContentHint')]"
            />

            <div class="row justify-end q-mt-md q-gutter-sm">
              <q-btn :label="t('common.cancel')" flat v-close-popup />
              <q-btn :label="t('verbali.saveDraftBtn')" color="primary" type="submit" unelevated :loading="submitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog 3: Visualizzazione Verbale / Firme Audit -->
    <q-dialog v-model="showViewDialog">
      <q-card style="width: min(700px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold">{{ currentVerbale?.title }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div class="q-mb-md">
            <q-badge
              :color="currentVerbale?.is_signed ? 'positive' : 'orange-9'"
              class="q-pa-xs text-subtitle2"
            >
              <q-icon :name="currentVerbale?.is_signed ? 'verified' : 'pending_actions'" class="q-mr-xs" />
              {{ currentVerbale?.is_signed ? t('verbali.viewOfficialSigned') : t('verbali.viewDraftReserved') }}
            </q-badge>
          </div>

          <q-card flat bordered class="bg-grey-1 q-pa-md q-mb-md">
            <div class="text-subtitle2 text-grey-8 q-mb-xs">{{ t('verbali.viewContentHeader') }}</div>
            <div style="white-space: pre-wrap; font-family: monospace;" class="text-body2">
              {{ currentVerbale?.content }}
            </div>
          </q-card>

          <!-- Firme Apposte -->
          <div class="text-subtitle1 text-weight-bold q-mb-sm">
            <q-icon name="history_edu" class="q-mr-xs text-primary" />{{ t('verbali.viewSignaturesHeader') }}
          </div>

          <div v-if="verbaleSignatures.length === 0" class="text-grey-6 text-caption">
            {{ t('verbali.noSignaturesYet') }}
          </div>
          <q-list v-else bordered separator class="rounded-borders">
            <q-item v-for="sig in verbaleSignatures" :key="sig.id">
              <q-item-section avatar>
                <q-avatar icon="verified_user" color="positive" text-color="white" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ sig.user_name || t('verbali.defaultSigner') }}</q-item-label>
                <q-item-label caption>
                  {{ t('verbali.signedAt', { date: new Date(sig.signed_at).toLocaleString(locale) }) }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-badge color="grey-8">{{ t('verbali.ipLabel', { ip: sig.ip_address || t('verbaliManagement.ipRegistered') }) }}</q-badge>
              </q-item-section>
            </q-item>
          </q-list>

          <div class="row justify-end q-mt-md q-gutter-sm">
            <q-btn
              v-if="currentVerbale?.is_signed"
              color="negative"
              icon="picture_as_pdf"
              :label="t('verbali.downloadPdfBtn')"
              unelevated
              @click="downloadPdf(currentVerbale)"
            />
            <q-btn :label="t('common.close')" flat v-close-popup />
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog 4: ODG Espanso Riunione -->
    <q-dialog v-model="showAgendaDialog">
      <q-card style="width: min(550px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center q-pb-none bg-grey-2">
          <div class="text-subtitle1 text-weight-bold">
            <q-icon name="list_alt" class="q-mr-xs text-primary" />{{ t('verbali.dialogAgendaTitle') }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-pa-md">
          <div class="text-h6 text-primary q-mb-sm">{{ currentMeeting?.title }}</div>
          <div class="text-caption text-grey-7 q-mb-md">
            {{ t('verbali.agendaDateTime', { date: currentMeeting?.date ? new Date(currentMeeting.date).toLocaleDateString(locale) : '', time: `${currentMeeting?.start_time || ''} - ${currentMeeting?.end_time || ''}` }) }}
          </div>
          <div class="bg-grey-1 q-pa-md rounded-borders" style="white-space: pre-wrap;">
            {{ currentMeeting?.agenda }}
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

// Stato navigazione
const activeTab = ref('verbali')
const verbaliList = ref([])
const meetingsList = ref([])
const templatesList = ref([])

const loadingVerbali = ref(false)
const loadingMeetings = ref(false)
const submitting = ref(false)

// Modali
const showMeetingDialog = ref(false)
const showVerbaleDialog = ref(false)
const showViewDialog = ref(false)
const showAgendaDialog = ref(false)

const isEditingVerbale = ref(false)
const selectedTemplate = ref(null)
const currentVerbale = ref(null)
const currentMeeting = ref(null)
const verbaleSignatures = ref([])

// Form Riunione
const meetingForm = ref({
  title: '',
  date: new Date().toISOString().split('T')[0],
  start_time: '15:00',
  end_time: '16:30',
  class_id: '',
  meeting_type: 'consiglio_classe',
  agenda: ''
})

// Form Verbale
const verbaleForm = ref({
  id: '',
  meeting_id: '',
  title: '',
  content: ''
})

// Colonne Tabella Verbali
const verbaliColumns = computed(() => [
  { name: 'title', label: t('verbali.colTitle'), field: 'title', align: 'left', sortable: true },
  { name: 'status', label: t('verbali.colStatus'), field: 'status', align: 'center' },
  { name: 'editing', label: t('verbali.colPermissions'), field: 'can_edit', align: 'center' },
  { name: 'created_at', label: t('verbali.colCreatedDate'), field: r => new Date(r.created_at).toLocaleDateString(locale.value), align: 'center' },
  { name: 'signatures', label: t('verbali.colSignatures'), align: 'center' },
  { name: 'actions', label: t('common.actions'), align: 'center' }
])

// Colonne Tabella Riunioni
const meetingsColumns = computed(() => [
  { name: 'title', label: t('verbali.colMeetingTitle'), field: 'title', align: 'left', sortable: true },
  { name: 'date', label: t('verbali.colDate'), field: r => new Date(r.date).toLocaleDateString(locale.value), align: 'center', sortable: true },
  { name: 'time', label: t('verbali.colTime'), field: r => `${r.start_time} - ${r.end_time}`, align: 'center' },
  { name: 'agenda', label: t('verbali.colAgenda'), align: 'left' },
  { name: 'actions', label: t('common.actions'), align: 'center' }
])

const formatMeetingType = (type) => {
  switch (type) {
    case 'consiglio_classe': return t('verbali.types.consiglio_classe')
    case 'collegio_docenti': return t('verbali.types.collegio_docenti')
    case 'dipartimento': return t('verbali.types.dipartimento')
    default: return t('verbali.types.generale')
  }
}

// Caricamento Dati
const loadVerbali = async () => {
  loadingVerbali.value = true
  try {
    const res = await verbaliService.getAllVerbali()
    verbaliList.value = res.data || []
  } catch {
    $q.notify({ type: 'negative', message: t('verbali.notifyLoadVerbaliError') })
  } finally {
    loadingVerbali.value = false
  }
}

const loadMeetings = async () => {
  loadingMeetings.value = true
  try {
    const res = await verbaliService.listMeetings()
    meetingsList.value = res.data || []
  } catch {
    $q.notify({ type: 'negative', message: t('verbali.notifyLoadMeetingsError') })
  } finally {
    loadingMeetings.value = false
  }
}

const loadTemplates = async () => {
  try {
    const res = await verbaliService.getTemplates()
    templatesList.value = res.data || []
  } catch {
    console.error('Errore caricamento modelli')
  }
}

// Azioni Modello
const onTemplateSelected = (tpl) => {
  if (!tpl) return
  meetingForm.value.title = tpl.title
  meetingForm.value.agenda = tpl.default_agenda
  meetingForm.value.meeting_type = tpl.meeting_type
}

// Creazione Riunione
const openCreateMeetingDialog = () => {
  selectedTemplate.value = null
  meetingForm.value = {
    title: '',
    date: new Date().toISOString().split('T')[0],
    start_time: '15:00',
    end_time: '16:30',
    class_id: '',
    meeting_type: 'consiglio_classe',
    agenda: ''
  }
  showMeetingDialog.value = true
}

const submitMeeting = async () => {
  submitting.value = true
  try {
    await verbaliService.createMeeting(meetingForm.value)
    $q.notify({ type: 'positive', message: t('verbali.notifyMeetingSaved') })
    showMeetingDialog.value = false
    await loadMeetings()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('verbali.notifyMeetingSaveError') })
  } finally {
    submitting.value = false
  }
}

// Redazione / Modifica Verbale
const openCreateVerbaleForMeeting = (meeting) => {
  isEditingVerbale.value = false
  verbaleForm.value = {
    id: '',
    meeting_id: meeting.id,
    title: t('verbali.meetingVerbalePrefix', { title: meeting.title }),
    content: selectedTemplate.value?.template_content || t('verbali.meetingVerbaleInitialContent', { date: meeting.date, agenda: meeting.agenda })
  }
  showVerbaleDialog.value = true
}

const openEditVerbale = (verbale) => {
  isEditingVerbale.value = true
  verbaleForm.value = {
    id: verbale.id,
    meeting_id: verbale.meeting_id,
    title: verbale.title,
    content: verbale.content
  }
  showVerbaleDialog.value = true
}

const submitVerbale = async () => {
  submitting.value = true
  try {
    if (isEditingVerbale.value) {
      await verbaliService.updateVerbale(verbaleForm.value.id, {
        title: verbaleForm.value.title,
        content: verbaleForm.value.content
      })
      $q.notify({ type: 'positive', message: t('verbali.notifyDraftUpdated') })
    } else {
      await verbaliService.createVerbale({
        meeting_id: verbaleForm.value.meeting_id,
        title: verbaleForm.value.title,
        content: verbaleForm.value.content
      })
      $q.notify({ type: 'positive', message: t('verbali.notifyDraftCreated') })
    }
    showVerbaleDialog.value = false
    await loadVerbali()
  } catch (err) {
    const msg = err.response?.data?.error || t('verbali.notifyVerbaleSaveError')
    $q.notify({ type: 'negative', message: msg })
  } finally {
    submitting.value = false
  }
}

// Firma Digitale con Audit IP
const signVerbale = async (id) => {
  $q.dialog({
    title: t('verbali.confirmSignTitle'),
    message: t('verbali.confirmSignMessage'),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await verbaliService.signVerbale(id)
      $q.notify({ type: 'positive', message: t('verbali.notifySignedSuccess') })
      await loadVerbali()
    } catch (err) {
      $q.notify({ type: 'negative', message: err.response?.data?.error || t('verbali.notifySignError') })
    }
  })
}

// Visualizzazione
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

const viewSignatures = async (verbale) => {
  await openViewVerbale(verbale)
}

const showAgendaDetail = (meeting) => {
  currentMeeting.value = meeting
  showAgendaDialog.value = true
}

// Download PDF
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
    $q.notify({ type: 'negative', message: t('verbali.notifyPdfError') })
  }
}

onMounted(() => {
  loadVerbali()
  loadMeetings()
  loadTemplates()
})
</script>
