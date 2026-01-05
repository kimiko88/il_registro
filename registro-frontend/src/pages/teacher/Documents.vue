<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">Documenti Didattici</div>
       <q-btn color="primary" icon="note_add" label="Nuovo Documento" @click="showCreateDialog = true" />
    </div>

    <!-- Filters -->
    <div class="row q-gutter-md q-mb-md">
        <q-select dense outlined v-model="filter.class" :options="classes" label="Classe" style="min-width: 120px" bg-color="white" />
        <q-select dense outlined v-model="filter.type" :options="['Tutti', 'PDP', 'PFI', 'PCTO', 'Programmazione']" label="Tipo" style="min-width: 150px" bg-color="white" />
    </div>

    <!-- Documents Grid -->
    <div class="row q-col-gutter-md">
        <div class="col-12 col-md-4 col-lg-3" v-for="doc in documents" :key="doc.id">
            <q-card class="hover-shadow">
                <q-card-section>
                    <div class="row items-center justify-between">
                         <q-chip :color="getTypeColor(doc.type)" text-color="white" size="sm">{{ doc.type }}</q-chip>
                         <q-btn flat round icon="more_vert" size="sm">
                             <q-menu>
                                 <q-list style="min-width: 100px">
                                     <q-item clickable v-close-popup @click="editDocument(doc)"><q-item-section>Modifica</q-item-section></q-item>
                                     <q-item clickable v-close-popup><q-item-section>PDF / Stampa</q-item-section></q-item>
                                     <q-item clickable v-close-popup class="text-primary" @click="startSigning(doc)"><q-item-section>Firma Documento</q-item-section></q-item>
                                     <q-item clickable v-close-popup class="text-negative"><q-item-section>Elimina</q-item-section></q-item>
                                 </q-list>
                             </q-menu>
                         </q-btn>
                    </div>
                    <div class="text-h6 q-mt-sm ellipsis">{{ doc.title }}</div>
                    <div class="text-subtitle2 text-grey-8">{{ doc.studentName || 'Classe Intera' }}</div>
                    <div class="text-caption text-grey q-mt-xs">Aggiornato: {{ doc.date }}</div>
                </q-card-section>
                
                <q-separator />
                
                <q-card-actions align="right">
                    <q-chip v-if="doc.status === 'draft'" icon="edit" label="Bozza" size="xs" color="grey-3" />
                    <q-chip v-if="doc.status === 'published'" icon="check" label="Pubblicato" size="xs" color="green-1" text-color="green" />
                    <q-btn flat color="primary" label="Apri" @click="openDocument(doc)" />
                </q-card-actions>
            </q-card>
        </div>
    </div>

    <!-- Create Dialog -->
    <q-dialog v-model="showCreateDialog" maximized transition-show="slide-up" transition-hide="slide-down">
        <q-card class="bg-grey-1">
            <q-toolbar class="bg-primary text-white">
                <q-btn flat round dense icon="close" v-close-popup />
                <q-toolbar-title>Editor Documento</q-toolbar-title>
                <q-btn flat label="Salva Bozza" icon="save" class="q-mr-sm" />
                <q-btn flat label="Pubblica / Invia" icon="send" />
            </q-toolbar>

            <q-card-section class="q-pa-md">
                <div class="row q-col-gutter-lg justify-center">
                    <div class="col-12 col-md-8">
                        <q-card class="q-pa-md">
                            <!-- Meta Data -->
                            <div class="row q-col-gutter-md q-mb-md">
                                <div class="col-12 col-md-4">
                                    <q-select outlined v-model="newDoc.type" :options="['PDP', 'PFI', 'Relazione', 'Verbale']" label="Tipo Documento" />
                                </div>
                                <div class="col-12 col-md-4">
                                    <q-select outlined v-model="newDoc.student" :options="['Rossi Mario', 'Bianchi Anna']" label="Studente (Opzionale)" />
                                </div>
                                <div class="col-12 col-md-4">
                                    <q-input outlined v-model="newDoc.title" label="Titolo" />
                                </div>
                            </div>

                            <q-separator class="q-mb-md" />

                            <!-- Editor -->
                            <div class="text-subtitle2 q-mb-sm">Contenuto</div>
                            <q-editor
                                v-model="newDoc.content"
                                min-height="400px"
                                :toolbar="[
                                    ['bold', 'italic', 'strike', 'underline'],
                                    ['token', 'hr', 'link', 'custom_btn'],
                                    ['print', 'fullscreen'],
                                    ['quote', 'unordered', 'ordered', 'outdent', 'indent'],
                                    ['undo', 'redo']
                                ]"
                            />
                        </q-card>
                    </div>
                    
                    <!-- Sidebar Helpers -->
                    <div class="col-12 col-md-3">
                        <q-card class="q-pa-md q-mb-md">
                           <div class="text-subtitle1 text-weight-bold">Template & Aiuto</div>
                           <q-list dense bordered class="rounded-borders q-mt-sm">
                               <q-item clickable v-ripple @click="loadTemplate('pdp')"><q-item-section>Template PDP Standard</q-item-section></q-item>
                               <q-item clickable v-ripple @click="loadTemplate('bes')"><q-item-section>Template BES</q-item-section></q-item>
                           </q-list>
                        </q-card>
                        
                        <q-card class="q-pa-md">
                           <div class="text-subtitle1 text-weight-bold">Storico Versioni</div>
                           <q-timeline color="secondary" class="q-mt-sm">
                               <q-timeline-entry title="Bozza creata" subtitle="Oggi, 10:00" />
                           </q-timeline>
                        </q-card>
                    </div>
                </div>
            </q-card-section>
        </q-card>
    </q-dialog>

  </q-page>

  <!-- Signing Dialog -->
  <q-dialog v-model="showSignDialog">
      <q-card style="min-width: 350px">
          <q-card-section>
              <div class="text-h6">Firma Digitale Reale</div>
              <div class="text-caption">Stai firmando: {{ docToSign?.title }}</div>
          </q-card-section>

          <q-card-section class="q-pt-none">
              <q-input v-model="signPin" label="Inserisci PIN Firma (Es. 1234)" type="password" outlined autofocus @keyup.enter="confirmSign" />
              <div class="text-caption text-grey q-mt-sm">
                  Nota: Questa azione renderà il documento immutabile.
              </div>
          </q-card-section>

          <q-card-actions align="right" class="text-primary">
              <q-btn flat label="Annulla" v-close-popup />
              <q-btn flat label="Firma Ora" @click="confirmSign" :loading="docsStore.loading" />
          </q-card-actions>
      </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useDocumentsStore } from 'src/stores/documents'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const docsStore = useDocumentsStore()
const showCreateDialog = ref(false)
const filter = reactive({ class: '5A', type: 'Tutti' })
const classes = ['5A', '4B']

const documents = ref([
    { id: 1, title: 'PDP - Rossi Mario', type: 'PDP', studentName: 'Rossi Mario', date: '10/01/2025', status: 'draft' },
    { id: 2, title: 'Programmazione Annuale Storia', type: 'Programmazione', studentName: '', date: '01/09/2024', status: 'published' }
])

const newDoc = reactive({
    type: 'PDP',
    student: null,
    title: '',
    content: '<b>Piano Didattico Personalizzato</b><br><br>...'
})

const getTypeColor = (type) => {
    if (type === 'PDP') return 'purple';
    if (type === 'PFI') return 'orange';
    return 'blue-grey';
}

const loadTemplate = (tpl) => {
    newDoc.content = tpl === 'pdp' 
        ? '<h2>Piano Didattico Personalizzato</h2><p>Dati generali...</p>' 
        : '<h2>Bisogni Educativi Speciali</h2><p>Diagnosi...</p>'
}

const editDocument = (doc) => {
    newDoc.title = doc.title;
    // Load content mock
    showCreateDialog.value = true;
}

const openDocument = (doc) => {
    editDocument(doc);
}

// Signing Logic
const showSignDialog = ref(false)
const docToSign = ref(null)
const signPin = ref('')

const startSigning = (doc) => {
    docToSign.value = doc
    signPin.value = ''
    showSignDialog.value = true
}

const confirmSign = async () => {
    if (!docToSign.value) return
    try {
        await docsStore.signDocument(docToSign.value.id.toString(), signPin.value)
        $q.notify({ type: 'positive', message: 'Documento firmato con successo!' })
        showSignDialog.value = false
        // Refresh?
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore durante la firma: PIN errato' })
    }
}
</script>

<style scoped>
.hover-shadow:hover {
    box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}
</style>
