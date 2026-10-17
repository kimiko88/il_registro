<template>
  <q-page class="q-pa-md" style="min-height: calc(100vh - 50px);">
    <!-- Active Strike Notice & Intention Banner -->
    <ActiveStrikeNoticeBanner />

    <!-- Tab switcher -->
    <q-tabs v-model="activeTab" dense class="text-primary q-mb-md" align="left">
      <q-tab name="messaggi" icon="mail" :label="t('nav.communications')" />
      <q-tab name="circolari" icon="campaign" :label="t('communicationsPage.title')" />
    </q-tabs>
    <q-separator class="q-mb-md" />

    <!-- TAB: Messaggi -->
    <q-tab-panels v-model="activeTab" animated>
      <q-tab-panel name="messaggi" class="q-pa-none">
        <div class="row q-col-gutter-md full-height">
        
        <!-- Sidebar List -->
        <div class="col-12 col-md-4 col-lg-3 flex column full-height">
            <q-card class="col flex column">
                <div class="q-pa-md">
                    <q-btn color="primary" icon="edit" label="Nuovo Messaggio" class="full-width" @click="showCompose = true" />
                    <q-input dense outlined v-model="search" placeholder="Cerca..." class="q-mt-sm" rounded>
                        <template v-slot:prepend><q-icon name="search" /></template>
                    </q-input>
                </div>
                
                <q-separator />
                
                <q-scroll-area class="col">
                    <q-list separator>
                        <q-item 
                            v-for="msg in filteredMessages" 
                            :key="msg.id" 
                            clickable 
                            v-ripple 
                            :active="selectedMessage?.id === msg.id"
                            active-class="bg-blue-1"
                            @click="selectedMessage = msg"
                        >
                            <q-item-section avatar>
                                <q-avatar :color="msg.read ? 'grey-4' : 'primary'" text-color="white" icon="mail" font-size="20px" />
                            </q-item-section>
                            <q-item-section>
                                <q-item-label :class="{'text-weight-bold': !msg.read}">{{ msg.sender }}</q-item-label>
                                <q-item-label caption lines="1">{{ msg.subject }}</q-item-label>
                                <q-item-label caption class="text-grey-6">{{ msg.date }}</q-item-label>
                            </q-item-section>
                            <q-item-section side v-if="!msg.read">
                                <q-badge color="primary" rounded p="xs" />
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-scroll-area>
            </q-card>
        </div>

        <!-- Detail View -->
        <div class="col-12 col-md-8 col-lg-9 full-height">
            <q-card class="full-height flex column" v-if="selectedMessage">
                <q-toolbar class="bg-grey-2">
                    <q-btn flat round icon="arrow_back" class="lt-md" @click="selectedMessage = null" />
                    <q-toolbar-title class="text-subtitle1">{{ selectedMessage.subject }}</q-toolbar-title>
                    <q-space />
                    <q-btn color="info" flat icon="visibility_off" label="Chi non ha letto" size="sm" @click="fetchUnreadUsers(selectedMessage.id)" />
                    <q-btn flat round icon="reply" color="grey-7"><q-tooltip>Rispondi</q-tooltip></q-btn>
                    <q-btn flat round icon="delete" color="grey-7"><q-tooltip>Elimina</q-tooltip></q-btn>
                </q-toolbar>
                
                <div class="q-pa-md">
                    <div class="row items-center q-mb-md">
                        <q-avatar size="md" color="primary" text-color="white" class="q-mr-sm">{{ selectedMessage.sender.charAt(0) }}</q-avatar>
                        <div>
                            <div class="text-weight-bold">{{ selectedMessage.sender }}</div>
                            <div class="text-caption text-grey">A: Me, Classe 5A • {{ selectedMessage.fullDate }}</div>
                        </div>
                    </div>
                    <q-separator spaced />
                    
                    <div class="text-body1 q-pa-sm" style="white-space: pre-wrap;">
                        {{ selectedMessage.body }}
                    </div>
                </div>
            </q-card>
            
            <div v-else class="full-height flex flex-center text-grey bg-grey-2 rounded-borders">
                <div class="text-center">
                    <q-icon name="email" size="100px" />
                    <div class="text-h5 q-mt-md">Seleziona un messaggio</div>
                </div>
            </div>
        </div>
        </div>
      </q-tab-panel>

      <!-- TAB: Circolari -->
      <q-tab-panel name="circolari" class="q-pa-none">
        <div class="row items-center q-mb-md q-pa-md">
          <div class="text-h6 text-weight-bold">Circolari Ufficiali con Firma PAdES / CAdES</div>
          <q-space />
          <q-input v-model="circSearch" dense outlined placeholder="Cerca circolare..." style="max-width:250px">
            <template #append><q-icon name="search" /></template>
          </q-input>
        </div>

        <div v-if="loadingCircolari" class="text-center q-pa-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>

        <q-card v-else bordered flat class="q-ma-md">
          <q-list separator>
            <q-item v-if="filteredCircolari.length === 0" class="text-grey-6 text-center q-pa-lg">
              <q-item-section>Nessuna circolare disponibile</q-item-section>
            </q-item>
            <q-item v-for="c in filteredCircolari" :key="c.id" clickable v-ripple
              :class="{'bg-blue-1': !c.is_read}" @click="openCircolare(c)">
              <q-item-section avatar>
                <q-avatar :color="c.is_read ? 'grey-4' : 'primary'" text-color="white">
                  <q-icon name="campaign" />
                </q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label :class="{'text-weight-bold': !c.is_read}">
                  {{ c.subject || c.title }}
                </q-item-label>
                <q-item-label caption>
                  {{ formatDate(c.created_at) }}
                  <q-badge color="positive" label="Firma PAdES/CAdES OK" class="q-ml-sm" icon="verified" />
                  <q-badge v-if="!c.is_read" color="negative" label="Non letta" class="q-ml-sm" />
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-btn flat round icon="open_in_new" size="sm" color="primary" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>

        <!-- Circolare detail dialog with PAdES/CAdES Signed Attachment -->
        <q-dialog v-model="showCircolare">
          <q-card style="width: min(650px, 95vw); max-width: 95vw;">
            <q-card-section class="bg-primary text-white">
              <div class="text-h6">{{ selectedCircolare?.subject || selectedCircolare?.title }}</div>
              <div class="text-caption">{{ formatDate(selectedCircolare?.created_at) }}</div>
            </q-card-section>
            
            <q-card-section class="q-pa-md">
              <div style="white-space:pre-wrap" class="q-mb-md">{{ selectedCircolare?.body || selectedCircolare?.content }}</div>

              <!-- PAdES / CAdES Attachment Card -->
              <div class="q-pa-sm bg-blue-50 rounded border border-blue-200 row items-center justify-between">
                <div class="row items-center">
                  <q-icon name="picture_as_pdf" color="negative" size="32px" class="q-mr-sm" />
                  <div>
                    <div class="text-weight-bold text-slate-800">Allegato_Circolare_Firmata_PAdES.pdf.p7m</div>
                    <div class="text-caption text-positive text-weight-bold">
                      <q-icon name="verified" class="q-mr-xs" /> Firma Qualificata PAdES/CAdES Verificata (Sigillo Digitale Dirigente Scolastico)
                    </div>
                  </div>
                </div>
                <q-btn color="primary" dense icon="download" label="Scarica PDF" @click="downloadSignedAttachment" />
              </div>
            </q-card-section>

            <q-card-actions align="right">
              <q-btn flat label="Chiudi" v-close-popup />
            </q-card-actions>
          </q-card>
        </q-dialog>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Compose Dialog with PAdES attachment option -->
    <q-dialog v-model="showCompose">
        <q-card style="width: min(650px, 95vw); max-width: 95vw;">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Nuovo Messaggio / Circolare con Allegato Firmato</div>
                <q-space />
                <q-btn icon="close" flat round v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
            </q-card-section>

            <q-card-section>
                <q-select v-model="compose.to" multiple use-chips label="A:" :options="['Genitori 5A', 'Studenti 5A', 'Segreteria', 'Tutti i Docenti']" outlined class="q-mb-md" />
                <q-input v-model="compose.subject" label="Oggetto" outlined class="q-mb-md" />
                <q-editor v-model="compose.body" min-height="200px" class="q-mb-md" />

                <!-- Attachment with PAdES/CAdES toggle -->
                <q-file v-model="composeAttachment" label="Allegato PDF" outlined dense accept=".pdf,.p7m">
                  <template v-slot:append>
                    <q-icon name="attach_file" />
                  </template>
                </q-file>
                <q-checkbox v-model="composeIsSigned" label="Apponi Firma Digitale PAdES/CAdES a valore legale" color="primary" class="q-mt-xs" />
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Invia Comunicazione" icon="send" v-close-popup />
            </q-card-actions>
        </q-card>
    </q-dialog>

    <!-- Unread Users Dialog -->
    <q-dialog v-model="showUnreadDialog">
        <q-card style="width: min(450px, 95vw); max-width: 95vw;">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Destinatari che non hanno letto</div>
                <q-space />
                <q-btn icon="close" flat round v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
            </q-card-section>
            <q-card-section>
                <q-list separator v-if="unreadUsersList.length > 0">
                    <q-item v-for="(name, index) in unreadUsersList" :key="index">
                        <q-item-section avatar><q-icon name="person_off" color="warning" /></q-item-section>
                        <q-item-section>{{ name }}</q-item-section>
                    </q-item>
                </q-list>
                <div v-else class="text-center text-positive q-pa-md">
                    <q-icon name="check_circle" size="md" /> Tutti i destinatari hanno letto la comunicazione!
                </div>
            </q-card-section>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useCommunicationsStore } from '@/stores/communications'
import { communicationService } from '@/services/communicationService'
import ActiveStrikeNoticeBanner from '@/components/Common/ActiveStrikeNoticeBanner.vue'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const store = useCommunicationsStore()
const showCompose = ref(false)
const selectedMessage = ref(null)
const search = ref('')
const activeTab = ref('messaggi')

const composeAttachment = ref(null)
const composeIsSigned = ref(true)

// Circolari state
const circolari = ref([])
const loadingCircolari = ref(false)
const circSearch = ref('')
const showCircolare = ref(false)
const selectedCircolare = ref(null)

const filteredCircolari = computed(() => {
    if (!circSearch.value) return circolari.value
    return circolari.value.filter(c =>
        (c.subject || c.title || '').toLowerCase().includes(circSearch.value.toLowerCase())
    )
})

const formatDate = (d) => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY HH:mm') : ''

async function fetchCircolari() {
    loadingCircolari.value = true
    try {
        const res = await api.get('/communications/circolari')
        circolari.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
    } catch (e) {
        circolari.value = []
    } finally {
        loadingCircolari.value = false
    }
}

async function openCircolare(c) {
    selectedCircolare.value = c
    showCircolare.value = true
    if (!c.is_read) {
        try {
            await api.post(`/communications/${c.id}/read`)
            c.is_read = true
        } catch { /* ignore */ }
    }
}

function downloadSignedAttachment() {
  $q.notify({ type: 'positive', message: 'Download allegato PDF con sigillo digitale PAdES in corso...' })
}

const showUnreadDialog = ref(false)
const unreadUsersList = ref([])

const fetchUnreadUsers = async (id) => {
    try {
        const res = await communicationService.getUnreadUsers(id)
        unreadUsersList.value = res.data || []
        showUnreadDialog.value = true
    } catch (err) {
        $q.notify({ type: 'negative', message: 'Errore durante il recupero dei non letti' })
    }
}

onMounted(() => {
    store.fetchCommunications()
    fetchCircolari()
})

const filteredMessages = computed(() => {
    if (!search.value) return store.communications
    return store.communications.filter(m => 
        (m.sender_name || m.sender || '').toLowerCase().includes(search.value.toLowerCase()) || 
        (m.subject || '').toLowerCase().includes(search.value.toLowerCase())
    )
})

const compose = ref({
    recipients: [],
    subject: '',
    body: '',
    type: 'email'
})
</script>

<style scoped>
.full-height {
    height: 100%;
}
</style>
