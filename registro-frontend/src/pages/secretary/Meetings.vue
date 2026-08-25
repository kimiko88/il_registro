<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="admin_panel_settings" color="primary" class="q-mr-sm" />
          Gestione Assemblee Scolastiche
        </div>
        <div class="text-caption text-grey">Crea e gestisci le assemblee e le iscrizioni dei partecipanti</div>
      </div>
      <q-btn color="primary" icon="add" label="Nuova Assemblea" @click="showCreateDialog = true" />
    </div>

    <!-- Table of meetings -->
    <q-card bordered>
      <q-table
        :rows="meetings"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        no-data-label="Nessuna assemblea registrata"
      >
        <template #body-cell-is_mandatory="{ row }">
          <q-td>
            <q-badge :color="row.is_mandatory ? 'negative' : 'grey'" :label="row.is_mandatory ? 'Sì' : 'No'" />
          </q-td>
        </template>
        <template #body-cell-registrations_count="{ row }">
          <q-td>
            {{ row.registrations_count }}{{ row.max_participants ? '/' + row.max_participants : '' }}
          </q-td>
        </template>
        <template #body-cell-actions="{ row }">
          <q-td class="text-right">
            <q-btn flat round icon="list_alt" color="primary" size="sm" @click="viewRegistrations(row)">
              <q-tooltip>Lista Iscritti</q-tooltip>
            </q-btn>
            <q-btn flat round icon="delete" color="negative" size="sm" @click="deleteMeeting(row.id)">
              <q-tooltip>Elimina</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Create Dialog -->
    <q-dialog v-model="showCreateDialog" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Nuova Assemblea / Incontro</div>
        </q-card-section>

        <q-card-section class="q-gutter-md">
          <q-input v-model="form.title" label="Titolo *" dense outlined />
          <q-input v-model="form.description" type="textarea" label="Descrizione" dense outlined rows="3" />
          <q-input v-model="form.location" label="Luogo (es. Aula Magna / Online)" dense outlined />
          <q-input v-model="form.meeting_date" type="datetime-local" label="Data e Ora Incontro *" dense outlined stack-label />
          <q-input v-model="form.registration_deadline" type="datetime-local" label="Scadenza Iscrizioni" dense outlined stack-label />
          <q-input v-model.number="form.max_participants" type="number" label="Max Partecipanti (vuoto = illimitati)" dense outlined />
          <q-toggle v-model="form.is_mandatory" label="Assemblea Obbligatoria" color="primary" />
          <q-select
            v-model="form.target_roles"
            :options="['parent', 'student', 'teacher']"
            multiple use-chips
            label="Destinatari"
            dense outlined
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Crea" :loading="saving" @click="createMeeting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Registrations List Dialog -->
    <q-dialog v-model="showRegistrationsDialog">
      <q-card style="min-width: 500px">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6">Iscritti — {{ selectedMeeting?.title }}</div>
          <q-btn flat round icon="close" color="white" v-close-popup />
        </q-card-section>
        <q-card-section>
          <div v-if="loadingRegs" class="text-center q-pa-md">
            <q-spinner-dots color="primary" size="30px" />
          </div>
          <q-list v-else-if="registrations.length > 0" separator>
            <q-item v-for="r in registrations" :key="r.id">
              <q-item-section avatar>
                <q-avatar color="primary" text-color="white" icon="person" size="sm" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ r.user_name }}</q-item-label>
                <q-item-label caption>{{ r.user_email }} · {{ r.user_role }}</q-item-label>
              </q-item-section>
              <q-item-section side class="text-caption text-grey">
                {{ formatDate(r.registered_at) }}
              </q-item-section>
            </q-item>
          </q-list>
          <div v-else class="text-center text-grey q-pa-md">Nessun iscritto al momento</div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar, date as qdate } from 'quasar'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const meetings = ref([])
const showCreateDialog = ref(false)
const showRegistrationsDialog = ref(false)
const selectedMeeting = ref(null)
const registrations = ref([])
const loadingRegs = ref(false)

const form = ref({
  title: '',
  description: '',
  location: '',
  meeting_date: '',
  registration_deadline: '',
  max_participants: null,
  is_mandatory: false,
  target_roles: ['parent', 'student']
})

const columns = [
  { name: 'title', label: 'Titolo', field: 'title', align: 'left', sortable: true },
  { name: 'meeting_date', label: 'Data/Ora', field: 'meeting_date', align: 'left', format: v => formatDate(v) },
  { name: 'location', label: 'Luogo', field: 'location', align: 'left' },
  { name: 'is_mandatory', label: 'Obbligatoria', field: 'is_mandatory', align: 'center' },
  { name: 'registrations_count', label: 'Iscritti', field: 'registrations_count', align: 'center' },
  { name: 'actions', label: 'Azioni', field: 'actions', align: 'right' }
]

const formatDate = d => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY HH:mm') : ''

onMounted(() => loadMeetings())

async function loadMeetings() {
  loading.value = true
  try {
    const res = await api.get('/general-meetings')
    meetings.value = res.data || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore caricamento assemblee' })
  } finally {
    loading.value = false
  }
}

async function createMeeting() {
  if (!form.value.title || !form.value.meeting_date) {
    $q.notify({ type: 'warning', message: 'Compila titolo e data incontro' })
    return
  }
  saving.value = true
  try {
    const payload = {
      ...form.value,
      meeting_date: new Date(form.value.meeting_date).toISOString(),
      registration_deadline: form.value.registration_deadline ? new Date(form.value.registration_deadline).toISOString() : null
    }
    await api.post('/general-meetings', payload)
    $q.notify({ type: 'positive', message: 'Assemblea creata con successo' })
    showCreateDialog.value = false
    await loadMeetings()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore creazione assemblea' })
  } finally {
    saving.value = false
  }
}

async function viewRegistrations(m) {
  selectedMeeting.value = m
  showRegistrationsDialog.value = true
  loadingRegs.value = true
  try {
    const res = await api.get(`/general-meetings/${m.id}/registrations`)
    registrations.value = res.data || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore caricamento iscritti' })
  } finally {
    loadingRegs.value = false
  }
}

async function deleteMeeting(id) {
  $q.dialog({
    title: 'Conferma eliminazione',
    message: 'Eliminare questa assemblea?',
    cancel: true, persistent: true
  }).onOk(async () => {
    try {
      await api.delete(`/general-meetings/${id}`)
      $q.notify({ type: 'positive', message: 'Assemblea eliminata' })
      await loadMeetings()
    } catch {
      $q.notify({ type: 'negative', message: 'Errore eliminazione' })
    }
  })
}
</script>
