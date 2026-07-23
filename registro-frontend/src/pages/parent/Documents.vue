<template>
  <q-page class="q-pa-md" role="main">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          Documenti e Comunicazioni
        </h1>
        <div class="text-subtitle1 text-slate-600 q-mt-xs">
          Fruizione e scarico documenti ufficiali per {{ childName }}
        </div>
      </div>

      <div class="row items-center q-gutter-sm">
        <q-input
          v-model="searchQuery"
          placeholder="Cerca documento..."
          dense
          outlined
          class="bg-white rounded-lg"
          style="min-width: 220px"
        >
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>
    </div>

    <!-- Category Tabs -->
    <q-card flat class="glass-card rounded-xl q-mb-lg overflow-hidden">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey-7 bg-white"
        active-color="primary"
        indicator-color="primary"
        align="left"
        narrow-indicator
      >
        <q-tab name="all" label="Tutti i Documenti" icon="folder" no-caps />
        <q-tab name="pagelle" label="Pagelle e Valutazioni" icon="assignment" no-caps />
        <q-tab name="circolari" label="Circolari e Moduli" icon="description" no-caps />
        <q-tab name="didattica" label="Piani Didattici (PDP/PEI)" icon="folder_shared" no-caps />
      </q-tabs>
    </q-card>

    <!-- Document List -->
    <q-card flat class="glass-card rounded-xl overflow-hidden">
      <q-list separator class="q-py-xs">
        <q-item
          v-for="doc in filteredDocuments"
          :key="doc.id"
          clickable
          v-ripple
          class="q-py-md hover-bg-grey"
          @click="handleDocumentClick(doc)"
        >
          <q-item-section avatar>
            <q-avatar size="44px" :color="getCategoryColor(doc.category)" text-color="white">
              <q-icon :name="getCategoryIcon(doc.category)" size="24px" />
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <div class="row items-center">
              <q-item-label class="text-weight-bold text-subtitle1 text-slate-800 q-mr-sm">
                {{ doc.title }}
              </q-item-label>
              <q-badge color="red" class="text-weight-bold" v-if="doc.isNew">Nuovo</q-badge>
            </div>
            <q-item-label caption class="text-slate-600 q-mt-xs">
              {{ doc.subtitle }} • Caricato il {{ doc.date }}
            </q-item-label>
          </q-item-section>

          <q-item-section side class="gt-xs">
            <q-chip dense size="sm" :color="doc.signed ? 'positive' : 'warning'" text-color="white">
              {{ doc.signed ? 'Firmato' : 'Da Consultare' }}
            </q-chip>
          </q-item-section>

          <q-item-section side>
            <div class="row items-center q-gutter-xs">
              <q-btn flat round color="primary" icon="visibility" @click.stop="previewDocument(doc)">
                <q-tooltip>Anteprima Documento</q-tooltip>
              </q-btn>
              <q-btn flat round color="primary" icon="download" @click.stop="downloadDocument(doc)">
                <q-tooltip>Scarica File</q-tooltip>
              </q-btn>
            </div>
          </q-item-section>
        </q-item>

        <div v-if="filteredDocuments.length === 0" class="text-center q-pa-xl text-slate-500">
          <q-icon name="folder_off" size="64px" color="grey-4" class="q-mb-md" />
          <div class="text-h6 text-weight-bold">Nessun documento trovato</div>
          <div class="text-caption">Non ci sono documenti corrispondenti alla ricerca per questa categoria.</div>
        </div>
      </q-list>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

const activeTab = ref('all')
const searchQuery = ref('')

const childName = computed(() => {
  if (!selectedChild.value) return 'tuo figlio'
  return `${selectedChild.value.first_name || selectedChild.value.firstName || ''} ${selectedChild.value.last_name || selectedChild.value.lastName || ''}`.trim()
})

const documents = ref([
  {
    id: '1',
    title: 'Pagella Primo Quadrimestre',
    subtitle: 'Valutazioni intermedie A.S. 2024/2025',
    category: 'pagelle',
    date: '15/02/2025',
    isNew: true,
    signed: true
  },
  {
    id: '2',
    title: 'Piano Didattico Personalizzato (PDP)',
    subtitle: 'Documento approvato dal consiglio di classe',
    category: 'didattica',
    date: '10/11/2024',
    isNew: false,
    signed: true
  },
  {
    id: '3',
    title: 'Circolare n. 42 - Uscita Didattica Museo',
    subtitle: 'Autorizzazione e quota di partecipazione',
    category: 'circolari',
    date: '20/01/2025',
    isNew: false,
    signed: false
  },
  {
    id: '4',
    title: 'Modulo Richiesta Permesso Uscita Anticipata',
    subtitle: 'Modulo stampabile per giustifica cartacea',
    category: 'circolari',
    date: '01/10/2024',
    isNew: false,
    signed: true
  }
])

const filteredDocuments = computed(() => {
  return documents.value.filter(doc => {
    const matchesTab = activeTab.value === 'all' || doc.category === activeTab.value
    const matchesSearch = !searchQuery.value ||
      doc.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      doc.subtitle.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchesTab && matchesSearch
  })
})

function getCategoryIcon(cat) {
  switch (cat) {
    case 'pagelle': return 'assignment'
    case 'didattica': return 'folder_shared'
    case 'circolari': return 'description'
    default: return 'insert_drive_file'
  }
}

function getCategoryColor(cat) {
  switch (cat) {
    case 'pagelle': return 'primary'
    case 'didattica': return 'orange'
    case 'circolari': return 'indigo'
    default: return 'grey-7'
  }
}

function previewDocument(doc) {
  $q.dialog({
    title: doc.title,
    message: `Visualizzazione anteprima per "${doc.title}" (${doc.subtitle}).`,
    ok: 'Chiudi'
  })
}

function downloadDocument(doc) {
  doc.isNew = false
  $q.notify({
    type: 'positive',
    message: `Download avviato: ${doc.title}.pdf`,
    icon: 'download',
    timeout: 2500
  })
}

function handleDocumentClick(doc) {
  previewDocument(doc)
}
</script>

<style scoped>
.hover-bg-grey:hover {
  background-color: rgba(99, 102, 241, 0.04);
}
</style>
