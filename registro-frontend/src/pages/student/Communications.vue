<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-md">Le Mie Comunicazioni</div>

    <q-card>
        <q-toolbar class="bg-grey-2">
            <q-input dense outlined v-model="search" placeholder="Cerca..." class="full-width" bg-color="white">
                <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
        </q-toolbar>
        
        <q-separator />

        <q-card-section v-if="loading" class="text-center">
            <q-spinner color="primary" size="3em" />
        </q-card-section>

        <q-list separator v-else>
            <q-item 
                v-for="msg in filteredMessages" 
                :key="msg.id" 
                clickable 
                v-ripple 
                @click="openMessage(msg)"
                :class="{'bg-blue-1': !msg.read_at}"
            >
                <q-item-section avatar>
                    <q-avatar color="primary" text-color="white" icon="campaign" />
                </q-item-section>
                
                <q-item-section>
                    <q-item-label class="text-weight-bold">{{ msg.subject }}</q-item-label>
                    <q-item-label caption lines="1">{{ msg.sender_id }}</q-item-label>
                </q-item-section>

                <q-item-section side>
                    <q-item-label caption>{{ formatDate(msg.created_at) }}</q-item-label>
                    <q-chip v-if="msg.type === 'circolare'" size="xs" color="orange" text-color="white">Circolare</q-chip>
                </q-item-section>
            </q-item>
            
            <q-item v-if="filteredMessages.length === 0">
                 <q-item-section class="text-center text-grey q-pa-lg">Nessun messaggio trovato</q-item-section>
            </q-item>
        </q-list>
    </q-card>
    
    <!-- Message Detail Dialog -->
    <q-dialog v-model="showMessage" transition-show="scale" transition-hide="scale">
        <q-card style="min-width: 500px">
            <q-card-section class="bg-primary text-white">
                <div class="text-h6">{{ selectedMessage?.subject }}</div>
                <div class="text-subtitle2">{{ selectedMessage?.sender_id }} - {{ formatDate(selectedMessage?.created_at) }}</div>
            </q-card-section>

            <q-card-section class="q-pt-md">
                <div class="text-body1" style="white-space: pre-wrap;">{{ selectedMessage?.body }}</div>
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Chiudi" v-close-popup />
                <q-btn flat label="Scarica Allegato" icon="download" color="primary" v-if="selectedMessage?.hasAttachment" />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCommunicationsStore } from 'src/stores/communications'
import { date } from 'quasar'

const store = useCommunicationsStore()
const search = ref('')
const showMessage = ref(false)
const selectedMessage = ref(null)

onMounted(() => {
    store.fetchMessages()
})

const loading = computed(() => store.loading)

const filteredMessages = computed(() => {
    if(!search.value) return store.messages
    return store.messages.filter(m => m.subject.toLowerCase().includes(search.value.toLowerCase()))
})

const openMessage = (msg) => {
    selectedMessage.value = msg
    showMessage.value = true
    // Mark as read API call would go here
}

const formatDate = (val) => {
    return date.formatDate(val, 'DD/MM/YYYY HH:mm')
}
</script>
