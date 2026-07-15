<template>
  <q-page class="q-pa-md" style="height: calc(100vh - 50px);"> <!-- Fill height minus header -->
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

    <!-- Compose Dialog -->
    <q-dialog v-model="showCompose">
        <q-card style="min-width: 600px">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Nuovo Messaggio</div>
                <q-space />
                <q-btn icon="close" flat round v-close-popup />
            </q-card-section>

            <q-card-section>
                <q-select v-model="compose.to" multiple use-chips label="A:" :options="['Genitori 5A', 'Studenti 5A', 'Segreteria']" outlined class="q-mb-md" />
                <q-input v-model="compose.subject" label="Oggetto" outlined class="q-mb-md" />
                <q-editor v-model="compose.body" min-height="200px" />
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Invia" icon="send" v-close-popup />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCommunicationsStore } from 'src/stores/communications'



const store = useCommunicationsStore()
const showCompose = ref(false)
const selectedMessage = ref(null)
const search = ref('')

onMounted(() => {
    store.fetchCommunications()
})

const filteredMessages = computed(() => {
    if (!search.value) return store.communications
    return store.communications.filter(m => 
        (m.sender_name || '').toLowerCase().includes(search.value.toLowerCase()) || 
        (m.subject || '').toLowerCase().includes(search.value.toLowerCase())
    )
})

const compose = ref({
    recipients: [], // Changed from to
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
