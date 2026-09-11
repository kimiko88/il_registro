<template>
  <q-page class="q-pa-md q-gutter-y-md">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-xs">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none text-primary">
          <q-icon name="account_balance" class="q-mr-sm" />Verbali & Modelli Riunioni
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none">
          Pannello di controllo della Dirigenza: consultazione verbali ufficiali firmati, calendario riunioni e gestione modelli ODG
        </p>
      </div>
      <div class="q-gutter-sm">
        <q-btn
          color="primary"
          icon="add_circle"
          label="Nuovo Modello ODG"
          unelevated
          @click="openCreateTemplateDialog"
        />
        <q-btn
          color="secondary"
          icon="event"
          label="Conferenza / Riunione"
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
        Informativa Ministeriale sulla Riservatezza dei Consigli di Classe:
      </div>
      <div class="text-caption text-grey-8">
        I verbali in fase di bozza sono di esclusiva competenza del docente coordinatore e del segretario verbalista e non compaiono in questo pannello.
        <strong>Vengono pubblicati ed archiviati in sola lettura per la Dirigente Scolastica esclusivamente a seguito dell'avvenuta firma digitale</strong>, con tracciamento certificato di data, ora e indirizzo IP.
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
        <q-tab name="verbali" icon="verified" label="Verbali Ufficiali d'Istituto (Solo Firmati)" />
        <q-tab name="templates" icon="dashboard_customize" label="Modelli Verbali & Ordini del Giorno (ODG)" />
        <q-tab name="meetings" icon="calendar_month" label="Calendario Riunioni d'Istituto" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated>
        <!-- Tab 1: Verbali Ufficiali Firmati -->
        <q-tab-panel name="verbali" class="q-pa-md">
          <div class="row items-center justify-between q-mb-md">
            <div class="text-subtitle1 text-weight-bold text-grey-8">
              Archivio Ufficiale Documenti Vidimati ({{ verbaliList.length }} verbali firmati)
            </div>
            <q-input
              v-model="searchVerbali"
              dense
              outlined
              placeholder="Cerca verbale per titolo..."
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
            no-data-label="Nessun verbale firmato attualmente presente negli archivi"
          >
            <!-- Badge Stato -->
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge color="positive" class="q-pa-xs text-weight-bold">
                  <q-icon name="lock" class="q-mr-xs" /> Ufficiale / Firmato
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
                  label="Vedi Firme IP"
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
                  <q-tooltip>Leggi Verbale Ufficiale</q-tooltip>
                </q-btn>
                <q-btn
                  flat
                  round
                  dense
                  color="negative"
                  icon="picture_as_pdf"
                  @click="downloadPdf(props.row)"
                >
                  <q-tooltip>Scarica PDF Ufficiale</q-tooltip>
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
                Modelli di Verbale e Ordini del Giorno Standardizzati
              </div>
              <div class="text-caption text-grey-6">
                I modelli creati dalla Dirigente compaiono ai docenti durante la programmazione di riunioni e consigli
              </div>
            </div>
            <q-btn
              color="primary"
              icon="add"
              label="Crea Nuovo Modello"
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
            no-data-label="Nessun modello configurato"
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
                  <q-tooltip>Modifica Modello</q-tooltip>
                </q-btn>
                <q-btn
                  flat
                  round
                  dense
                  color="negative"
                  icon="delete"
                  @click="deleteTemplate(props.row)"
                >
                  <q-tooltip>Elimina Modello</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- Tab 3: Calendario Riunioni d'Istituto -->
        <q-tab-panel name="meetings" class="q-pa-md">
          <div class="row items-center justify-between q-mb-md">
            <div class="text-subtitle1 text-weight-bold text-grey-8">
              Calendario Istituzionale delle Riunioni e Convocazioni
            </div>
            <q-btn
              color="secondary"
              icon="add"
              label="Programma Nuova Riunione"
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
            no-data-label="Nessuna riunione registrata a calendario"
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
                  {{ props.row.agenda || 'Nessun ODG specificato' }}
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
            {{ isEditingTemplate ? 'Modifica Modello Verbale & ODG' : 'Nuovo Modello Verbale & Ordine del Giorno' }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-form @submit="submitTemplate" class="q-gutter-y-sm">
            <q-input
              v-model="templateForm.title"
              label="Titolo Modello (es. Consiglio di Classe Intermedio, Collegio Docenti PTOF) *"
              outlined
              dense
              :rules="[val => !!val || 'Campo obbligatorio']"
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="templateForm.meeting_type"
                  :options="meetingTypeOptions"
                  emit-value
                  map-options
                  label="Tipologia Riunione *"
                  outlined
                  dense
                />
              </div>
              <div class="col-12 col-md-6">
                <q-input
                  v-model="templateForm.description"
                  label="Descrizione / Ambito d'uso"
                  outlined
                  dense
                />
              </div>
            </div>

            <q-input
              v-model="templateForm.default_agenda"
              type="textarea"
              label="Ordine del Giorno Predefinito (ODG numerato) *"
              outlined
              dense
              rows="5"
              hint="Inserisci i punti all'ODG separati da a capo"
              :rules="[val => !!val || 'L\'Ordine del Giorno è obbligatorio']"
            />

            <q-input
              v-model="templateForm.template_content"
              type="textarea"
              label="Schema / Bozza Struttura Verbale (con segnaposto) *"
              outlined
              dense
              rows="8"
              hint="Puoi usare segnaposto come {{classe}}, {{data}}, {{odg}}, {{presidente}}, {{segretario}}"
              :rules="[val => !!val || 'La bozza di struttura è obbligatoria']"
            />

            <div class="row justify-end q-mt-md q-gutter-sm">
              <q-btn label="Annulla" flat v-close-popup />
              <q-btn label="Salva Modello" color="primary" type="submit" unelevated :loading="submitting" />
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
            <q-icon name="event_available" class="q-mr-sm" />Programma Riunione d'Istituto
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-form @submit="submitMeeting" class="q-gutter-y-sm">
            <q-input
              v-model="meetingForm.title"
              label="Titolo Riunione *"
              outlined
              dense
              :rules="[val => !!val || 'Campo obbligatorio']"
            />

            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-4">
                <q-input
                  v-model="meetingForm.date"
                  type="date"
                  label="Data *"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
              <div class="col-6 col-md-4">
                <q-input
                  v-model="meetingForm.start_time"
                  type="time"
                  label="Ora Inizio *"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
              <div class="col-6 col-md-4">
                <q-input
                  v-model="meetingForm.end_time"
                  type="time"
                  label="Ora Fine *"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
            </div>

            <q-select
              v-model="meetingForm.meeting_type"
              :options="meetingTypeOptions"
              emit-value
              map-options
              label="Tipologia Riunione *"
              outlined
              dense
            />

            <q-input
              v-model="meetingForm.agenda"
              type="textarea"
              label="Ordine del Giorno (ODG) *"
              outlined
              dense
              rows="4"
              :rules="[val => !!val || 'Campo obbligatorio']"
            />

            <div class="row justify-end q-mt-md q-gutter-sm">
              <q-btn label="Annulla" flat v-close-popup />
              <q-btn label="Crea Riunione" color="primary" type="submit" unelevated :loading="submitting" />
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
              <q-icon name="verified" class="q-mr-xs" /> Atto Ufficiale Vidimato e Archiviato
            </q-badge>
          </div>

          <q-card flat bordered class="bg-grey-1 q-pa-md q-mb-md">
            <div class="text-subtitle2 text-grey-8 q-mb-xs">Testo del Verbale:</div>
            <div style="white-space: pre-wrap; font-family: monospace;" class="text-body2">
              {{ currentVerbale?.content }}
            </div>
          </q-card>

          <!-- Firme Apposte -->
          <div class="text-subtitle1 text-weight-bold q-mb-sm">
            <q-icon name="history_edu" class="q-mr-xs text-primary" />Firme Digitali e Certificazione IP:
          </div>

          <q-list bordered separator class="rounded-borders">
            <q-item v-for="sig in verbaleSignatures" :key="sig.id">
              <q-item-section avatar>
                <q-avatar icon="verified_user" color="positive" text-color="white" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ sig.user_name || 'Docente Firmatario' }}</q-item-label>
                <q-item-label caption>
                  Firmato il: {{ new Date(sig.signed_at).toLocaleString() }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-badge color="grey-8">IP: {{ sig.ip_address || 'Registrato' }}</q-badge>
              </q-item-section>
            </q-item>
          </q-list>

          <div class="row justify-end q-mt-md q-gutter-sm">
            <q-btn
              color="negative"
              icon="picture_as_pdf"
              label="Scarica PDF Ufficiale"
              unelevated
              @click="downloadPdf(currentVerbale)"
            />
            <q-btn label="Chiudi" flat v-close-popup />
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import verbaliService from '@/services/verbaliService'

const $q = useQuasar()

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
const meetingTypeOptions = [
  { label: 'Consiglio di Classe', value: 'consiglio_classe' },
  { label: 'Collegio dei Docenti', value: 'collegio_docenti' },
  { label: 'Dipartimento Disciplinare', value: 'dipartimento' },
  { label: 'Generale / Altro', value: 'generale' }
]

// Colonne Verbali Ufficiali
const verbaliColumns = computed(() => [
  { name: 'title', label: 'Titolo Verbale Ufficiale', field: 'title', align: 'left', sortable: true },
  { name: 'status', label: 'Stato Atto', field: 'status', align: 'center' },
  { name: 'created_at', label: 'Data Archiviazione', field: r => new Date(r.created_at).toLocaleDateString(), align: 'center', sortable: true },
  { name: 'signatures', label: 'Firme Digitali', align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'center' }
])

// Colonne Modelli
const templatesColumns = computed(() => [
  { name: 'title', label: 'Titolo Modello', field: 'title', align: 'left', sortable: true },
  { name: 'meeting_type', label: 'Tipologia Riunione', field: 'meeting_type', align: 'center' },
  { name: 'agenda', label: 'Anteprima ODG Predefinito', align: 'left' },
  { name: 'created_at', label: 'Data Creazione', field: r => new Date(r.created_at).toLocaleDateString(), align: 'center' },
  { name: 'actions', label: 'Azioni', align: 'center' }
])

// Colonne Riunioni
const meetingsColumns = computed(() => [
  { name: 'title', label: 'Titolo Riunione', field: 'title', align: 'left', sortable: true },
  { name: 'meeting_type', label: 'Tipologia', field: 'meeting_type', align: 'center' },
  { name: 'date', label: 'Data', field: r => new Date(r.date).toLocaleDateString(), align: 'center', sortable: true },
  { name: 'time', label: 'Orario', field: r => `${r.start_time} - ${r.end_time}`, align: 'center' },
  { name: 'agenda', label: 'Ordine del Giorno (ODG)', align: 'left' }
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
    case 'consiglio_classe': return 'Consiglio di Classe'
    case 'collegio_docenti': return 'Collegio dei Docenti'
    case 'dipartimento': return 'Dipartimento'
    default: return 'Generale'
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
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento dei verbali ufficiali' })
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
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento dei modelli' })
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
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento del calendario riunioni' })
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
      $q.notify({ type: 'positive', message: 'Modello aggiornato con successo' })
    } else {
      await verbaliService.createTemplate(templateForm.value)
      $q.notify({ type: 'positive', message: 'Nuovo modello creato con successo' })
    }
    showTemplateDialog.value = false
    await loadTemplates()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore salvataggio modello' })
  } finally {
    submitting.value = false
  }
}

const deleteTemplate = (tpl) => {
  $q.dialog({
    title: 'Elimina Modello',
    message: `Sei sicuro di voler eliminare il modello "${tpl.title}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await verbaliService.deleteTemplate(tpl.id)
      $q.notify({ type: 'positive', message: 'Modello eliminato' })
      await loadTemplates()
    } catch {
      $q.notify({ type: 'negative', message: 'Errore eliminazione modello' })
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
    $q.notify({ type: 'positive', message: 'Riunione registrata nel calendario' })
    showMeetingDialog.value = false
    await loadMeetings()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore creazione riunione' })
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
    $q.notify({ type: 'negative', message: 'Errore durante il download del PDF' })
  }
}

onMounted(() => {
  loadVerbali()
  loadTemplates()
  loadMeetings()
})
</script>
