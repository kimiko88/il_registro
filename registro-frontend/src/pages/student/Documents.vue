<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-md">I Miei Documenti</div>

    <q-card>
        <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify">
            <q-tab name="reports" icon="school" label="Pagelle" />
            <q-tab name="plans" icon="assignment" label="PDP / PFI" />
            <q-tab name="certs" icon="verified" label="Certificati" />
        </q-tabs>
        <q-separator />

        <q-tab-panels v-model="tab" animated>
            <q-tab-panel name="reports">
                 <q-list separator>
                     <q-item v-for="doc in reports" :key="doc.id">
                         <q-item-section avatar>
                             <q-icon name="picture_as_pdf" color="red" size="md" />
                         </q-item-section>
                         <q-item-section>
                             <q-item-label>{{ doc.title }}</q-item-label>
                             <q-item-label caption>{{ doc.date }} - A.S. {{ doc.year }}</q-item-label>
                         </q-item-section>
                         <q-item-section side>
                             <q-btn flat round icon="download" color="primary" @click="download(doc)">
                                 <q-tooltip>Scarica PDF</q-tooltip>
                             </q-btn>
                             <q-btn flat round icon="visibility" color="grey" @click="preview(doc)">
                                 <q-tooltip>Anteprima</q-tooltip>
                             </q-btn>
                         </q-item-section>
                     </q-item>
                 </q-list>
            </q-tab-panel>

            <q-tab-panel name="plans">
                 <q-list separator>
                     <q-item v-for="doc in plans" :key="doc.id">
                         <q-item-section avatar>
                             <q-icon name="description" color="blue" size="md" />
                         </q-item-section>
                         <q-item-section>
                             <q-item-label>{{ doc.title }}</q-item-label>
                             <q-item-label caption>Firmato il {{ doc.date }}</q-item-label>
                         </q-item-section>
                          <q-item-section side>
                             <q-chip v-if="doc.signed" color="green" text-color="white" icon="verified" size="sm">Firmato</q-chip>
                             <q-btn v-else color="orange" label="Firma Ora" size="sm" />
                          </q-item-section>
                         <q-item-section side>
                             <q-btn flat round icon="download" color="primary" />
                         </q-item-section>
                     </q-item>
                 </q-list>
            </q-tab-panel>

             <q-tab-panel name="certs">
                 <div class="text-caption text-grey q-mb-md">Certificati di Frequenza e Iscrizione generati automaticamente.</div>
                 <q-item clickable v-ripple>
                     <q-item-section avatar><q-icon name="badge" color="purple" /></q-item-section>
                     <q-item-section>Certificato di Iscrizione</q-item-section>
                     <q-item-section side><q-btn outline label="Genera" size="sm" /></q-item-section>
                 </q-item>
             </q-tab-panel>
        </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tab = ref('reports')

const reports = ref([
    { id: 1, title: 'Pagella 1° Quadrimestre', date: '10/02/2025', year: '2024/25' },
    { id: 2, title: 'Pagella Finale', date: '12/06/2024', year: '2023/24' }
])

const plans = ref([
    { id: 101, title: 'Piano Didattico Personalizzato (PDP)', date: '15/11/2024', signed: true }
])

const download = (_doc) => {
    $q.notify({ type: 'progress', message: 'Download in corso...' })
}

const preview = (doc) => {
    $q.dialog({ title: doc.title, message: 'Anteprima non disponibile offline.' })
}
</script>
