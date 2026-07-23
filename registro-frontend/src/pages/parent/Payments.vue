<template>
  <q-page class="q-pa-md" role="main">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          Pagamenti Scolastici & PagoPA
        </h1>
        <div class="text-subtitle1 text-slate-600 q-mt-xs">
          Gestione delle quote scolastiche, contributi e ricevute di pagamento
        </div>
      </div>
      <q-chip color="primary" text-color="white" class="text-weight-bold q-py-md q-px-md shadow-sm">
        Totale da Pagare: € {{ totalPending }}
      </q-chip>
    </div>

    <!-- MOCK NOTICE BANNER -->
    <q-banner class="bg-amber-1 text-amber-10 rounded-xl q-mb-md border border-amber-200">
      <template v-slot:avatar>
        <q-icon name="science" color="amber-8" />
      </template>
      <strong>Modalità Dimostrativa / Simulazione</strong> — I pagamenti visualizzati in questa sezione sono ad uso dimostrativo. Nessun addebito bancario reale verrà effettuato.
    </q-banner>

    <div class="row q-col-gutter-lg">
      <div class="col-12 col-md-8">
        <q-card flat class="glass-card rounded-xl overflow-hidden">
          <q-tabs v-model="tab" dense class="text-grey-7 bg-white" active-color="primary" indicator-color="primary" align="justify">
            <q-tab name="pending" label="Da Pagare" icon="payment" class="text-amber-9" />
            <q-tab name="history" label="Storico & Ricevute" icon="history" />
          </q-tabs>
          <q-separator />

          <q-tab-panels v-model="tab" animated>
            <q-tab-panel name="pending" class="q-pa-none">
              <q-list separator>
                <q-item v-for="item in pendingItems" :key="item.id" class="q-py-md">
                  <q-item-section avatar>
                    <q-avatar color="amber-1" text-color="amber-9" icon="receipt" size="44px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-subtitle1 text-slate-800">{{ item.title }}</q-item-label>
                    <q-item-label caption class="text-slate-600">Scadenza: {{ item.dueDate }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <div class="text-h6 text-weight-bold text-primary">€ {{ item.amount }}</div>
                    <q-btn label="Paga con PagoPA" color="primary" unelevated no-caps size="sm" class="q-mt-xs rounded-lg" @click="pay(item)" />
                  </q-item-section>
                </q-item>
                <q-item v-if="pendingItems.length === 0" class="text-center q-pa-xl">
                  <q-item-section class="text-center text-positive">
                    <q-icon name="check_circle" size="56px" class="q-mb-sm" />
                    <div class="text-h6 text-weight-bold">Tutti i pagamenti sono in regola!</div>
                    <div class="text-caption text-slate-500">Nessun contributo o tassa scolastica in sospeso.</div>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-tab-panel>

            <q-tab-panel name="history" class="q-pa-none">
              <q-list separator>
                <q-item v-for="item in historyItems" :key="item.id" class="q-py-md">
                  <q-item-section avatar>
                    <q-avatar color="emerald-1" text-color="emerald-7" icon="verified" size="44px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-slate-800">{{ item.title }}</q-item-label>
                    <q-item-label caption class="text-slate-600">Pagato il: {{ item.paidDate }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <div class="text-weight-bold text-subtitle1">€ {{ item.amount }}</div>
                    <q-btn flat round icon="download" color="primary" size="sm" aria-label="Scarica ricevuta" @click="downloadReceipt(item)">
                      <q-tooltip>Scarica Ricevuta PDF</q-tooltip>
                    </q-btn>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-tab-panel>
          </q-tab-panels>
        </q-card>
      </div>

      <div class="col-12 col-md-4">
        <q-card flat class="glass-card rounded-xl q-pa-md">
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold text-slate-800 row items-center">
              <q-icon name="info" color="primary" class="q-mr-xs" />
              Info Pagamenti PagoPA
            </div>
            <p class="text-body2 text-slate-600 q-mt-sm">
              I pagamenti avvengono in modo sicuro tramite la piattaforma nazionale PagoPA.
              Le ricevute hanno piena validità fiscale e sono sempre scaricabili dalla sezione Storico.
            </p>
            <q-btn flat label="Guida PagoPA" color="primary" class="full-width rounded-lg q-mt-xs" no-caps @click="openPagoPAGuide" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Post Payment Summary Dialog -->
    <q-dialog v-model="showReceiptDialog">
      <q-card style="width: min(480px, 95vw)" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-positive text-white text-center q-pa-lg">
          <q-icon name="check_circle" size="64px" class="q-mb-xs" />
          <div class="text-h5 text-weight-bold">Pagamento Completato!</div>
          <div class="text-subtitle2 opacity-90">Operazione registrata con successo</div>
        </q-card-section>

        <q-card-section class="q-pa-md" v-if="completedPayment">
          <q-list class="q-gutter-y-xs">
            <q-item>
              <q-item-section>
                <q-item-label caption>Causale / Oggetto</q-item-label>
                <q-item-label class="text-weight-bold text-subtitle1">{{ completedPayment.title }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label caption>Importo Pagato</q-item-label>
                <q-item-label class="text-weight-bold text-h6 text-primary">€ {{ completedPayment.amount }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label caption>Data e Ora Transazione</q-item-label>
                <q-item-label class="text-weight-medium">{{ completedPayment.paidDate }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-separator />

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Chiudi" v-close-popup no-caps />
          <q-btn color="primary" icon="download" label="Scarica Ricevuta" unelevated no-caps class="rounded-lg" @click="downloadReceipt(completedPayment)" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tab = ref('pending')
const showReceiptDialog = ref(false)
const completedPayment = ref(null)

const totalPending = computed(() => {
  return pendingItems.value.reduce((acc, item) => acc + Number(item.amount), 0).toFixed(2)
})

const openPagoPAGuide = () => {
  $q.dialog({
    title: 'Guida ai Pagamenti PagoPA',
    message: 'PagoPA è la piattaforma nazionale che permette di effettuare pagamenti verso la Pubblica Amministrazione in modo semplice e sicuro. Puoi pagare online con carta di credito, conto corrente o app di pagamento, oppure sul territorio presso tabaccherie, ricevitorie e banconote abilitate.',
    ok: { label: 'Ho Capito', color: 'primary' }
  })
}

const pendingItems = ref([
    { id: 1, title: 'Assicurazione Scolastica Integrativa', dueDate: '30/01/2025', amount: '8.50' },
    { id: 2, title: 'Gita Scolastica - Firenze', dueDate: '15/02/2025', amount: '45.00' }
])

const historyItems = ref([
    { id: 101, title: 'Contributo Volontario', paidDate: '10/09/2024', amount: '120.00' }
])

const pay = (item) => {
    $q.loading.show({ message: 'Connessione al nodo PagoPA in corso...' })
    setTimeout(() => {
        $q.loading.hide()
        const paidItem = { ...item, paidDate: new Date().toLocaleString('it-IT') }
        pendingItems.value = pendingItems.value.filter(i => i.id !== item.id)
        historyItems.value.unshift(paidItem)
        completedPayment.value = paidItem
        showReceiptDialog.value = true
    }, 1500)
}

const downloadReceipt = (item) => {
    $q.notify({
        type: 'positive',
        message: `Ricevuta scaricata: ${item.title}.pdf`,
        icon: 'download',
        timeout: 2500
    })
}
</script>

<style scoped>
</style>
