<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="email" color="primary" class="q-mr-sm" />
          Comunicazioni Istituzionali
        </div>
        <div class="text-caption text-grey">Circolari, avvisi e comunicazioni scolastiche</div>
      </div>
      <q-btn round flat icon="refresh" @click="fetchMessages" />
    </div>

    <!-- Child Warning -->
    <q-card v-if="!selectedChild" class="text-center q-pa-lg bg-warning text-white q-mb-md">
      Seleziona un figlio dal menu in alto per visualizzare le sue comunicazioni
    </q-card>

    <div v-else class="row q-col-gutter-lg">
      <!-- Message List -->
      <div class="col-12 col-md-5">
        <q-card class="column full-height" style="min-height: 70vh">
          <q-card-section class="q-pa-none">
            <q-tabs v-model="tab" dense class="text-grey bg-grey-1" active-color="primary" indicator-color="primary" align="justify">
              <q-tab name="inbox">
                <div class="row items-center no-wrap">
                  <q-icon name="mail" class="q-mr-sm" />
                  <div>In Arrivo</div>
                  <q-badge color="red" floating v-if="unreadCount">{{ unreadCount }}</q-badge>
                </div>
              </q-tab>
              <q-tab name="circolari">
                <div class="row items-center no-wrap">
                  <q-icon name="campaign" class="q-mr-sm" />
                  <div>Circolari</div>
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

      <!-- Detail View -->
      <div class="col-12 col-md-7">
        <q-card class="full-height" style="min-height: 70vh">
          <div v-if="selectedMessage">
            <q-card-section class="bg-grey-1 row items-center justify-between">
              <div class="row items-center">
                <q-avatar color="primary" text-color="white" class="q-mr-md">
                  {{ selectedMessage.sender[0] }}
                </q-avatar>
                <div>
                  <div class="text-h6">{{ selectedMessage.subject }}</div>
                  <div class="text-caption">Da: <strong>{{ selectedMessage.sender }}</strong> &lt;{{ selectedMessage.email }}&gt;</div>
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
                  Firmato per Presa Visione
                </q-badge>
              </div>
              <q-btn flat icon="archive" label="Archivia" color="warning" />
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
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { storeToRefs } from 'pinia'
import { useParentStore } from 'src/stores/parent'
import { communicationService } from 'src/services/communicationService'

const $q = useQuasar()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

const tab = ref('inbox')
const search = ref('')
const selectedMessage = ref(null)
const messages = ref([])
const loading = ref(false)

onMounted(async () => {
  if (selectedChild.value) {
    await fetchMessages()
  }
})

watch(selectedChild, async (newVal) => {
  if (newVal) {
    selectedMessage.value = null
    await fetchMessages()
  } else {
    messages.value = []
    selectedMessage.value = null
  }
})

const fetchMessages = async () => {
  loading.value = true
  try {
    const res = await communicationService.getMessages()
    messages.value = (res.data || []).map(m => ({
      id: m.id,
      sender: m.sender_name || 'Sistema / Segreteria',
      email: m.sender_email || '',
      subject: m.subject,
      preview: m.body.substring(0, 50) + '...',
      body: m.body,
      date: new Date(m.created_at).toLocaleDateString('it-IT'),
      fullDate: new Date(m.created_at).toLocaleString('it-IT'),
      read: m.read || false,
      is_signed: m.is_signed || false,
      archived: m.archived || false
    }))
  } catch (e) {
    $q.notify({ message: 'Errore nel caricamento delle comunicazioni', color: 'negative' })
    console.error(e)
  } finally {
    loading.value = false
  }
}

const unreadCount = computed(() => messages.value.filter(m => !m.read && !m.archived).length)

const filteredMessages = computed(() => {
  return messages.value.filter(m => {
    let matchesTab = false;
    if (tab.value === 'inbox') matchesTab = !m.archived && m.type !== 'circolare';
    else if (tab.value === 'circolari') matchesTab = !m.archived && (m.type === 'circolare' || m.is_circolare);
    else if (tab.value === 'archive') matchesTab = m.archived;
    
    const matchesSearch = m.subject.toLowerCase().includes(search.value.toLowerCase()) || 
                          m.sender.toLowerCase().includes(search.value.toLowerCase());
    return matchesTab && matchesSearch;
  })
})

const selectMessage = (msg) => {
  selectedMessage.value = msg
  if (!msg.read) {
    msg.read = true
  }
}

const signReceipt = async (msg) => {
  try {
    await communicationService.signMessage(msg.id)
    msg.is_signed = true
    $q.notify({ type: 'positive', message: 'Firma registrata con successo' })
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante la firma della comunicazione' })
  }
}
</script>
