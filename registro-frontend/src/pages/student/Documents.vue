<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-md">Documenti Scolastici</div>

    <div class="row q-col-gutter-lg">
        <!-- Main Documents -->
        <div class="col-12 col-md-8">
            <q-card>
                <q-tabs v-model="tab" align="left" class="text-primary">
                    <q-tab name="personal" label="Personali (PDP/PFI)" />
                    <q-tab name="class" label="Classe" />
                    <q-tab name="reports" label="Pagelle" />
                </q-tabs>
                <q-separator />
                
                <q-tab-panels v-model="tab" animated>
                    <q-tab-panel name="personal">
                         <q-list separator>
                             <q-item v-for="doc in personalDocs" :key="doc.id" clickable v-ripple @click="download(doc)">
                                 <q-item-section avatar>
                                     <q-icon name="description" color="primary" />
                                 </q-item-section>
                                 <q-item-section>
                                     <q-item-label>{{ doc.title }}</q-item-label>
                                     <q-item-label caption>{{ doc.date }} - Firmato digitalmente</q-item-label>
                                 </q-item-section>
                                 <q-item-section side>
                                     <q-btn flat round icon="download" color="grey" />
                                 </q-item-section>
                             </q-item>
                         </q-list>
                    </q-tab-panel>

                    <q-tab-panel name="class">
                         <div class="text-grey text-center q-pa-md">Nessun documento di classe condiviso di recente.</div>
                    </q-tab-panel>

                    <q-tab-panel name="reports">
                         <q-list separator>
                             <q-item clickable v-ripple>
                                 <q-item-section avatar>
                                     <q-icon name="assignment" color="secondary" />
                                 </q-item-section>
                                 <q-item-section>
                                     <q-item-label>Pagella 1° Quadrimestre 2024/25</q-item-label>
                                     <q-item-label caption>Pubblicata il 01/02/2025</q-item-label>
                                 </q-item-section>
                                 <q-item-section side>
                                     <q-btn flat round icon="visibility" color="primary" />
                                 </q-item-section>
                             </q-item>
                         </q-list>
                    </q-tab-panel>
                </q-tab-panels>
            </q-card>
        </div>

        <!-- Signature Requests -->
        <div class="col-12 col-md-4">
            <q-card class="bg-orange-1">
                <q-card-section>
                    <div class="row items-center q-mb-sm">
                        <q-icon name="warning" color="orange" size="sm" class="q-mr-sm" />
                        <div class="text-subtitle1 text-weight-bold opacity-80">Da Firmare</div>
                    </div>
                    <q-list dense>
                        <q-item class="bg-white rounded-borders q-mb-sm">
                            <q-item-section>
                                <q-item-label>Patto Educativo di Corresponsabilità</q-item-label>
                            </q-item-section>
                            <q-item-section side>
                                <q-btn size="sm" color="primary" label="Firma" />
                            </q-item-section>
                        </q-item>
                    </q-list>
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
const tab = ref('personal')

const personalDocs = [
    { id: 1, title: 'PDP - Piano Didattico Personalizzato', date: '15/10/2024' },
    { id: 2, title: 'Autorizzazione Uscita Autonoma', date: '10/09/2024' }
]

const download = (doc) => {
    $q.notify({ message: `Download ${doc.title} avviato...`, color: 'primary', icon: 'cloud_download' })
}
</script>
