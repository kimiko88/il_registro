<template>
  <q-page class="q-pa-md bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">Comunicazioni & Bacheca</h1>
        <p class="text-subtitle2 text-slate-500 q-mt-xs q-mb-none">Circolari, avvisi scolastici e messaggi ufficiali</p>
      </div>
      <q-btn round flat icon="refresh" color="primary" @click="loadData" />
    </div>

    <!-- Main Navigation Tabs -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft q-mb-md">
      <q-tabs
        v-model="mainTab"
        dense
        class="text-slate-600 bg-slate-100 border-b border-slate-200"
        active-color="primary"
        indicator-color="primary"
        align="left"
        no-caps
      >
        <q-tab name="bacheca" icon="dashboard" label="Bacheca Scolastica" />
        <q-tab name="messages" icon="email" label="Messaggi Personali" />
      </q-tabs>
    </q-card>

    <!-- TAB 1: BACHECA -->
    <div v-if="mainTab === 'bacheca'">
      <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
        <q-card-section class="bg-slate-100 border-b border-slate-200 row items-center justify-between q-py-sm q-px-md">
          <div class="text-subtitle1 text-weight-bold text-slate-800">Avvisi e Bacheca della Scuola</div>
          <q-badge color="primary" class="q-px-sm q-py-xs text-weight-bold">
            {{ bachecaMessages.length }} Messaggi
          </q-badge>
        </q-card-section>

        <div v-if="loadingBacheca" class="text-center q-pa-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>

        <div v-else-if="bachecaMessages.length === 0" class="text-center q-pa-xl text-slate-400">
          <q-icon name="campaign" size="64px" class="q-mb-md opacity-40" />
          <div class="text-h6">Nessun avviso in bacheca</div>
          <div class="text-caption">Al momento non ci sono avvisi pubblicati per la tua scuola o classe.</div>
        </div>

        <q-list v-else separator class="rounded-lg">
          <q-expansion-item
            v-for="msg in bachecaMessages"
            :key="msg.id"
            group="bacheca"
            header-class="q-py-md hover:bg-slate-50 transition-colors"
          >
            <template v-slot:header>
              <q-item-section avatar>
                <q-avatar :color="!msg.read_at ? 'blue-1' : 'grey-2'" :text-color="!msg.read_at ? 'primary' : 'grey-7'" icon="campaign" size="44px" />
              </q-item-section>

              <q-item-section>
                <div class="row items-center q-gutter-xs q-mb-xs">
                  <span class="text-weight-bold text-slate-800 text-subtitle1">{{ msg.title || msg.subject }}</span>
                  <q-badge v-if="!msg.read_at && !msg.is_read" color="negative" class="q-px-xs text-weight-bold">
                    Nuovo
                  </q-badge>
                </div>
                <div class="text-caption text-slate-500">
                  Pubblicato il {{ formatDate(msg.created_at || msg.date) }} · Da {{ msg.sender_name || 'Segreteria' }}
                </div>
              </q-item-section>
            </template>

            <q-card flat class="bg-slate-50 border-t border-slate-200 q-pa-md">
              <div class="text-body1 text-slate-700 q-mb-md" style="white-space: pre-line">
                {{ msg.body || msg.content }}
              </div>

              <div class="row items-center justify-between border-t border-slate-200 q-pt-md">
                <div>
                  <q-chip v-if="msg.requires_signature && msg.signed_at" color="positive" text-color="white" icon="check_circle" size="sm" class="text-weight-bold">
                    Firmato per presa visione
                  </q-chip>
                  <q-btn
                    v-else-if="msg.requires_signature"
                    color="positive"
                    unelevated
                    size="sm"
                    icon="draw"
                    label="Firma per Presa Visione"
                    :loading="signingId === msg.id"
                    no-caps
                    @click="signBachecaMessage(msg.id)"
                  />
                </div>
              </div>
            </q-card>
          </q-expansion-item>
        </q-list>
      </q-card>
    </div>

    <!-- TAB 2: MESSAGGI PERSONALI -->
    <div v-else class="row q-col-gutter-lg">
      <div class="col-12 col-md-5">
        <q-card flat bordered class="column full-height rounded-xl bg-white shadow-soft" style="min-height: 70vh">
          <q-card-section class="q-pa-none">
            <q-tabs v-model="tab" dense class="text-grey bg-grey-1" active-color="primary" indicator-color="primary" align="justify">
              <q-tab name="inbox">
                <div class="row items-center no-wrap">
                  <q-icon name="mail" class="q-mr-sm" />
                  <div>In Arrivo</div>
                  <q-badge color="red" floating v-if="unreadCount">{{ unreadCount }}</q-badge>
                </div>
              </q-tab>
              <q-tab name="archive" icon="archive" label="Archivio" />
            </q-tabs>
            <q-separator />

            <q-input v-model="search" dense borderless placeholder="Cerca..." class="q-px-md">
              <template v-slot:append>
                <q-icon name="search" />
              </template>
            </q-input>
            <q-separator />
          </q-card-section>

          <q-card-section class="q-pa-none scroll col">
            <q-list separator>
              <q-item
                v-for="msg in filteredMessages"
                :key="msg.id"
                clickable
                v-ripple
                :active="selectedMessage?.id === msg.id"
                active-class="bg-blue-1 text-primary"
                @click="selectMessage(msg)"
              >
                <q-item-section avatar>
                  <q-avatar :color="msg.read ? 'grey-3' : 'blue'" text-color="white" size="md">
                    {{ msg.sender[0] }}
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label :class="{'text-weight-bold': !msg.read}">{{ msg.subject }}</q-item-label>
                  <q-item-label caption lines="1">{{ msg.sender }}</q-item-label>
                  <q-item-label caption lines="2">{{ msg.preview }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="text-caption">{{ msg.date }}</div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-7">
        <q-card flat bordered class="full-height rounded-xl bg-white shadow-soft" style="min-height: 70vh">
          <div v-if="selectedMessage">
            <q-card-section class="bg-grey-1 row items-center justify-between">
              <div class="row items-center">
                <q-avatar color="primary" text-color="white" class="q-mr-md">
                  {{ selectedMessage.sender[0] }}
                </q-avatar>
                <div>
                  <div class="text-h6">{{ selectedMessage.subject }}</div>
                  <div class="text-caption">Da: <strong>{{ selectedMessage.sender }}</strong></div>
                </div>
              </div>
              <div class="text-caption text-grey">{{ selectedMessage.fullDate }}</div>
            </q-card-section>
            <q-separator />

            <q-card-section class="q-pa-lg">
              <div class="text-body1" style="white-space: pre-line">{{ selectedMessage.body }}</div>
            </q-card-section>

            <q-separator />
            <q-card-actions align="right" class="q-pa-md row items-center justify-between">
              <div class="row items-center q-gutter-sm">
                <q-btn
                  v-if="!selectedMessage.is_signed"
                  color="positive"
                  icon="check"
                  label="Firma per Presa Visione"
                  no-caps
                  @click="signReceipt(selectedMessage)"
                />
                <q-badge v-else color="positive" class="q-pa-sm text-weight-bold" outline>
                  <q-icon name="check_circle" class="q-mr-xs" />
                  Letto e Firmato per Presa Visione
                </q-badge>
              </div>
            </q-card-actions>
          </div>

          <div v-else class="full-height flex flex-center text-grey column">
            <q-icon name="mail_outline" size="100px" color="grey-3" />
            <div class="text-h5 q-mt-md">Seleziona un messaggio</div>
          </div>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { communicationService } from 'src/services/communicationService'
import api from 'src/services/api'

const $q = useQuasar()

const mainTab = ref('bacheca')
const tab = ref('inbox')
const search = ref('')
const selectedMessage = ref(null)

const messages = ref([])
const bachecaMessages = ref([])
const loadingBacheca = ref(false)
const signingId = ref(null)

const formatDate = (d) => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY HH:mm') : ''

onMounted(() => {
  loadData()
})

function loadData() {
  fetchBacheca()
  fetchMessages()
}

async function fetchBacheca() {
  loadingBacheca.value = true
  try {
    const res = await api.get('/communications/bacheca')
    bachecaMessages.value = res.data || []
  } catch {
    bachecaMessages.value = []
  } finally {
    loadingBacheca.value = false
  }
}

async function fetchMessages() {
  try {
    const res = await communicationService.getMessages()
    messages.value = (res.data || []).map(m => ({
      id: m.id,
      sender: m.sender_name || 'Sistema',
      email: m.sender_email || '',
      subject: m.subject,
      preview: (m.body || '').substring(0, 50) + '...',
      body: m.body,
      date: new Date(m.created_at).toLocaleDateString('it-IT'),
      fullDate: new Date(m.created_at).toLocaleString('it-IT'),
      read: m.read || false,
      is_signed: m.is_signed || false,
      archived: m.archived || false
    }))
  } catch (e) {
    console.error(e)
  }
}

async function signBachecaMessage(id) {
  signingId.value = id
  try {
    await api.post(`/communications/${id}/sign`)
    $q.notify({ type: 'positive', message: 'Firma per presa visione registrata' })
    await fetchBacheca()
  } catch {
    $q.notify({ type: 'negative', message: 'Errore durante la registrazione della firma' })
  } finally {
    signingId.value = null
  }
}

const unreadCount = computed(() => messages.value.filter(m => !m.read && !m.archived).length)

const filteredMessages = computed(() => {
  return messages.value.filter(m => {
    const matchesTab = tab.value === 'inbox' ? !m.archived : m.archived
    const matchesSearch = m.subject.toLowerCase().includes(search.value.toLowerCase()) ||
                          m.sender.toLowerCase().includes(search.value.toLowerCase())
    return matchesTab && matchesSearch
  })
})

async function selectMessage(msg) {
  selectedMessage.value = msg
  if (!msg.read) {
    msg.read = true
    try {
      await communicationService.markAsRead(msg.id)
    } catch {
      msg.read = false
    }
  }
}

async function signReceipt(msg) {
  try {
    await communicationService.signMessage(msg.id)
    msg.is_signed = true
    $q.notify({ type: 'positive', message: 'Presa visione registrata con successo' })
  } catch {
    $q.notify({ type: 'negative', message: 'Errore nella registrazione della firma' })
  }
}
</script>
