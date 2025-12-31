<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">Comunicazioni</div>
       <q-btn icon="edit" label="Scrivi Messaggio" color="primary" @click="showCompose = true" />
    </div>

    <q-card>
        <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify">
            <q-tab name="inbox" icon="mail" label="Ricevuti" />
            <q-tab name="circulars" icon="article" label="Circolari" />
        </q-tabs>
        <q-separator />

        <q-tab-panels v-model="tab" animated>
            <q-tab-panel name="inbox">
                 <q-list separator>
                     <q-item v-for="msg in messages" :key="msg.id" clickable v-ripple @click="openMessage(msg)">
                         <q-item-section avatar>
                             <q-avatar color="primary" text-color="white">{{ msg.sender.charAt(0) }}</q-avatar>
                         </q-item-section>
                         <q-item-section>
                             <q-item-label class="text-weight-bold">{{ msg.subject }}</q-item-label>
                             <q-item-label caption lines="1">{{ msg.sender }} - {{ msg.date }}</q-item-label>
                         </q-item-section>
                         <q-item-section side v-if="!msg.read">
                             <q-badge color="blue" rounded p="xs" />
                         </q-item-section>
                     </q-item>
                 </q-list>
            </q-tab-panel>

            <q-tab-panel name="circulars">
                 <q-list separator>
                     <q-item clickable v-ripple>
                         <q-item-section avatar><q-icon name="campaign" color="orange" /></q-item-section>
                         <q-item-section>
                             <q-item-label>Chiusura per Festività</q-item-label>
                             <q-item-label caption>Dirigenza - 20/12/2024</q-item-label>
                         </q-item-section>
                         <q-item-section side><q-btn flat round icon="download" /></q-item-section>
                     </q-item>
                 </q-list>
            </q-tab-panel>
        </q-tab-panels>
    </q-card>

    <!-- Compose Dialog -->
    <q-dialog v-model="showCompose">
        <q-card style="min-width: 500px">
            <q-card-section class="text-h6">Nuovo Messaggio</q-card-section>
            <q-card-section>
                <q-select v-model="compose.recipient" :options="['Coordinatore Classe', 'Segreteria', 'Presidenza']" label="Destinatario" outlined class="q-mb-md" />
                <q-input v-model="compose.subject" label="Oggetto" outlined class="q-mb-md" />
                <q-input v-model="compose.body" type="textarea" label="Messaggio" outlined />
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
import { ref } from 'vue'

const tab = ref('inbox')
const showCompose = ref(false)
const compose = ref({ recipient: null, subject: '', body: '' })

const messages = ref([
    { id: 1, sender: 'Prof. Bianchi', subject: 'Risposta: Richiesta Colloquio', date: 'Oggi 10:30', read: false },
    { id: 2, sender: 'Segreteria', subject: 'Conferma Iscrizione', date: 'Ieri', read: true }
])

const openMessage = (msg) => {
    msg.read = true
    // Logic to open detail
}
</script>
