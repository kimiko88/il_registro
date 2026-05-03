<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">Comunicazioni</div>
       <q-btn round flat icon="refresh" @click="studentStore.fetchNotifications()" />
    </div>

    <div class="row q-col-gutter-lg">
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
                                <q-icon name="attachment" size="xs" color="grey" v-if="msg.hasAttachment" />
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-card-section>
            </q-card>
        </div>

        <!-- content -->
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

                    <q-card-section v-if="selectedMessage.attachments && selectedMessage.attachments.length">
                        <div class="text-subtitle2 q-mb-sm">Allegati ({{ selectedMessage.attachments.length }})</div>
                        <div class="row q-gutter-sm">
                            <q-chip 
                                v-for="att in selectedMessage.attachments" 
                                :key="att" 
                                clickable 
                                @click="downloadAttachment(att)"
                                color="grey-3" 
                                text-color="black" 
                                icon="attachment"
                            >
                                {{ att }}
                            </q-chip>
                        </div>
                    </q-card-section>
                    
                    <q-separator />
                    <q-card-actions align="right">
                         <q-btn flat icon="archive" label="Archivia" color="warning" />
                         <q-btn flat icon="reply" label="Rispondi" color="primary" disable>
                             <q-tooltip>Risposta disabilitata per studenti</q-tooltip>
                         </q-btn>
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
import { useStudentStore } from 'src/stores/student'
import { useQuasar } from 'quasar'
import { communicationService } from 'src/services/communicationService'

const $q = useQuasar()
const studentStore = useStudentStore()

const tab = ref('inbox')
const search = ref('')
const selectedMessage = ref(null)
const messages = ref([])
const loading = ref(false)

const fetchMessages = async () => {
    loading.value = true
    try {
        const res = await communicationService.getMessages()
        // Map backend messages to frontend format if necessary
        // Backend Message might have: id, sender_id, recipient_id, subject, body, created_at, read
        messages.value = (res.data || []).map(m => ({
            id: m.id,
            sender: m.sender_name || 'Sistema',
            email: m.sender_email || '',
            subject: m.subject,
            preview: m.body.substring(0, 50) + '...',
            body: m.body,
            date: new Date(m.created_at).toLocaleDateString('it-IT'),
            fullDate: new Date(m.created_at).toLocaleString('it-IT'),
            read: m.read || false,
            hasAttachment: false, // Update if backend supports attachments
            attachments: [],
            archived: m.archived || false
        }))
    } catch (e) {
        $q.notify({ message: 'Errore nel caricamento dei messaggi', color: 'negative' })
        console.error(e)
    } finally {
        loading.value = false
    }
}

const unreadCount = computed(() => messages.value.filter(m => !m.read && !m.archived).length)

const filteredMessages = computed(() => {
    return messages.value.filter(m => {
        const matchesTab = tab.value === 'inbox' ? !m.archived : m.archived;
        const matchesSearch = m.subject.toLowerCase().includes(search.value.toLowerCase()) || 
                              m.sender.toLowerCase().includes(search.value.toLowerCase());
        return matchesTab && matchesSearch;
    })
})

const selectMessage = (msg) => {
    selectedMessage.value = msg
    if (!msg.read) {
        msg.read = true
        // api.markAsRead(msg.id)
    }
}

const downloadAttachment = (name) => {
    $q.notify({ message: `Download ${name}`, color: 'primary' })
}

onMounted(() => {
    fetchMessages()
})
})
</script>
