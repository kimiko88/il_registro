<template>
  <q-page class="q-pa-md" role="main">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          {{ t('documentsPage.parentTitle') }}
        </h1>
        <div class="text-subtitle1 text-slate-600 q-mt-xs">
          {{ t('documentsPage.parentSubtitle', { name: childName }) }}
        </div>
      </div>

      <div class="row items-center q-gutter-sm">
        <q-input
          v-model="searchQuery"
          :placeholder="t('documentsPage.searchPlaceholderDoc')"
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
        <q-tab name="all" :label="t('documentsPage.tabs.all')" icon="folder" no-caps />
        <q-tab name="pagelle" :label="t('documentsPage.tabs.reportCards')" icon="assignment" no-caps />
        <q-tab name="circolari" :label="t('documentsPage.tabs.circulars')" icon="description" no-caps />
        <q-tab name="didattica" :label="t('documentsPage.tabs.plans')" icon="folder_shared" no-caps />
      </q-tabs>
    </q-card>

    <!-- Document List -->
    <q-card flat class="glass-card rounded-xl overflow-hidden">
      <div v-if="loading" class="text-center q-pa-xl">
        <q-spinner-dots color="primary" size="40px" />
        <div class="text-caption text-grey-7 q-mt-sm">{{ t('documentsPage.loadingDatabase') }}</div>
      </div>

      <q-list v-else separator class="q-py-xs">
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
              <q-badge color="red" class="text-weight-bold" v-if="doc.isNew">{{ t('documentsPage.isNew') }}</q-badge>
            </div>
            <q-item-label caption class="text-slate-600 q-mt-xs">
              {{ doc.subtitle }} <template v-if="doc.date">{{ t('documentsPage.uploadedOn', { date: doc.date }) }}</template>
            </q-item-label>
          </q-item-section>

          <q-item-section side class="gt-xs">
            <q-chip dense size="sm" :color="doc.signed ? 'positive' : 'warning'" text-color="white">
              {{ doc.signed ? t('documentsPage.signed') : t('documentsPage.toConsult') }}
            </q-chip>
          </q-item-section>

          <q-item-section side>
            <div class="row items-center q-gutter-xs">
              <q-btn flat round color="primary" icon="visibility" @click.stop="previewDocument(doc)">
                <q-tooltip>{{ t('documentsPage.previewTooltip') }}</q-tooltip>
              </q-btn>
              <q-btn flat round color="primary" icon="download" @click.stop="downloadDocument(doc)">
                <q-tooltip>{{ t('documentsPage.downloadTooltip') }}</q-tooltip>
              </q-btn>
            </div>
          </q-item-section>
        </q-item>

        <div v-if="filteredDocuments.length === 0" class="text-center q-pa-xl text-slate-500">
          <q-icon name="folder_off" size="64px" color="grey-4" class="q-mb-md" />
          <div class="text-h6 text-weight-bold">{{ t('documentsPage.noDocsFound') }}</div>
          <div class="text-caption">{{ t('documentsPage.noDocsDesc') }}</div>
        </div>
      </q-list>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import communicationService from '@/services/communicationService'
import pdpService from '@/services/pdpService'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

const activeTab = ref('all')
const searchQuery = ref('')
const loading = ref(true)
const documents = ref([])

const childName = computed(() => {
  if (!selectedChild.value) return t('documentsPage.yourChild')
  return `${selectedChild.value.first_name || selectedChild.value.firstName || ''} ${selectedChild.value.last_name || selectedChild.value.lastName || ''}`.trim()
})

async function loadDocuments() {
  loading.value = true
  const list = []
  const studentId = selectedChild.value?.id || selectedChild.value?.user_id

  // 1. Communications / Circolari
  try {
    const commsRes = await communicationService.getMessages()
    const comms = Array.isArray(commsRes.data) ? commsRes.data : (commsRes.data?.items || [])
    comms.forEach(c => {
      list.push({
        id: `comm-${c.id}`,
        title: c.title || c.subject || t('documentsPage.commsFallbackTitle'),
        subtitle: c.content ? c.content.substring(0, 80) + '...' : t('documentsPage.commsFallbackSubtitle'),
        category: 'circolari',
        date: c.created_at ? new Date(c.created_at).toLocaleDateString('it-IT') : '',
        isNew: !c.is_read,
        signed: !!c.is_signed,
        raw: c
      })
    })
  } catch (e) {
    console.warn('Could not load communications for documents view:', e)
  }

  // 2. PDP / PEI Plans
  if (studentId) {
    try {
      const pdpRes = await pdpService.getStudentPlans(studentId, '2025/2026')
      const plans = pdpRes.data?.plans || []
      plans.forEach(p => {
        list.push({
          id: `pdp-${p.id}`,
          title: t('documentsPage.pdpTitle', { title: p.title || 'PDP' }),
          subtitle: p.diagnosis ? t('documentsPage.diagnosisPrefix', { diagnosis: p.diagnosis }) : t('documentsPage.pdpFallbackSubtitle'),
          category: 'didattica',
          date: p.created_at ? new Date(p.created_at).toLocaleDateString('it-IT') : '',
          isNew: false,
          signed: !!p.family_approved_at,
          raw: p
        })
      })
    } catch (e) {
      console.warn('Could not load PDP plans for documents view:', e)
    }

    // 3. Fascicolo Documents
    try {
      const fascRes = await api.get(`/students/${studentId}/fascicolo`)
      const docs = fascRes.data?.documents || fascRes.data || []
      if (Array.isArray(docs)) {
        docs.forEach(d => {
          const isPagella = (d.category || d.title || '').toLowerCase().includes('pagell')
          list.push({
            id: `doc-${d.id}`,
            title: d.title || d.file_name || t('documentsPage.fascicoloFallbackTitle'),
            subtitle: d.category || t('documentsPage.fascicoloFallbackSubtitle'),
            category: isPagella ? 'pagelle' : 'circolari',
            date: d.created_at ? new Date(d.created_at).toLocaleDateString('it-IT') : '',
            isNew: false,
            signed: true,
            raw: d
          })
        })
      }
    } catch (e) {
      console.warn('Could not load student fascicolo for documents view:', e)
    }
  }

  documents.value = list
  loading.value = false
}

onMounted(() => {
  loadDocuments()
})

watch(selectedChild, () => {
  loadDocuments()
})

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
    message: t('documentsPage.previewDialogMessage', { subtitle: doc.subtitle, date: doc.date || 'N/D', status: doc.signed ? t('documentsPage.signedApproved') : t('documentsPage.toConsult') }),
    ok: t('common.close')
  })
}

function downloadDocument(doc) {
  doc.isNew = false
  $q.notify({
    type: 'positive',
    message: t('documentsPage.downloadStarted', { title: doc.title }),
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
