<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="text-h4 q-mb-md">Pagamenti Scolastici</div>

    <div class="row q-col-gutter-lg">
        <div class="col-12 col-md-8">
            <q-card>
                <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify">
                    <q-tab name="pending" label="Da Pagare" icon="payment" class="text-orange" />
                    <q-tab name="history" label="Storico" icon="history" />
                </q-tabs>
                <q-separator />

                <q-tab-panels v-model="tab" animated>
                    <q-tab-panel name="pending">
                         <q-list separator>
                             <q-item v-for="item in pendingItems" :key="item.id">
                                 <q-item-section avatar>
                                     <q-icon name="receipt" color="orange" />
                                 </q-item-section>
                                 <q-item-section>
                                     <q-item-label class="text-weight-bold">{{ item.title }}</q-item-label>
                                     <q-item-label caption>Scadenza: {{ item.dueDate }}</q-item-label>
                                 </q-item-section>
                                 <q-item-section side>
                                     <div class="text-h6 text-primary">€ {{ item.amount }}</div>
                                     <q-btn label="Paga con PagoPA" color="primary" unelevated size="sm" class="q-mt-xs" @click="pay(item)" />
                                 </q-item-section>
                             </q-item>
                             <q-item v-if="pendingItems.length === 0">
                                 <q-item-section class="text-center text-green q-pa-lg">
                                     <q-icon name="check_circle" size="48px" />
                                     <div>Nessun pagamento in sospeso.</div>
                                 </q-item-section>
                             </q-item>
                         </q-list>
                    </q-tab-panel>

                    <q-tab-panel name="history">
                         <q-list separator>
                             <q-item v-for="item in historyItems" :key="item.id">
                                 <q-item-section avatar>
                                     <q-icon name="check_circle" color="green" />
                                 </q-item-section>
                                 <q-item-section>
                                     <q-item-label>{{ item.title }}</q-item-label>
                                     <q-item-label caption>Pagato il: {{ item.paidDate }}</q-item-label>
                                 </q-item-section>
                                 <q-item-section side>
                                     <div class="text-weight-bold">€ {{ item.amount }}</div>
                                     <q-btn flat round icon="download" color="grey" size="sm" />
                                 </q-item-section>
                             </q-item>
                         </q-list>
                    </q-tab-panel>
                </q-tab-panels>
            </q-card>
        </div>

        <div class="col-12 col-md-4">
            <q-card class="bg-blue-1">
                <q-card-section>
                    <div class="text-subtitle1 text-weight-bold">Info Pagamenti</div>
                    <p class="text-caption q-mt-sm">
                        I pagamenti avvengono tramite la piattaforma PagoPA.
                        Le ricevute sono scaricabili dalla sezione Storico ed hanno validità fiscale.
                    </p>
                    <q-btn flat label="Guida PagoPA" color="primary" class="full-width" />
                </q-card-section>
            </q-card>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tab = ref('pending')

const pendingItems = ref([
    { id: 1, title: 'Assicurazione Scolastica Integrativa', dueDate: '30/01/2025', amount: '8.50' },
    { id: 2, title: 'Gita Scolastica - Firenze', dueDate: '15/02/2025', amount: '45.00' }
])

const historyItems = ref([
    { id: 101, title: 'Contributo Volontario', paidDate: '10/09/2024', amount: '120.00' }
])

const pay = (item) => {
    $q.loading.show({ message: 'Connessione a PagoPA...' })
    setTimeout(() => {
        $q.loading.hide()
        $q.notify({ type: 'positive', message: 'Pagamento simulato con successo!' })
        pendingItems.value = pendingItems.value.filter(i => i.id !== item.id)
        historyItems.value.unshift({ ...item, paidDate: 'Oggi' })
    }, 2000)
}
</script>
