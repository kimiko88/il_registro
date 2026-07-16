<template>
  <q-page class="q-pa-md bg-slate-50">
    <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">Documenti: {{ selectedChild ? `${selectedChild.first_name} ${selectedChild.last_name}` : '...' }}</div>

    <q-list bordered separator class="bg-white rounded-borders shadow-sm">
      <q-item-label header>Pagelle e Report</q-item-label>
      
      <q-item clickable v-ripple @click="download('Pagella 1 Quadrimestre')">
        <q-item-section avatar>
          <q-icon name="description" color="primary" />
        </q-item-section>
        <q-item-section>
          <q-item-label>Pagella 1° Quadrimestre</q-item-label>
          <q-item-label caption>A.S. 2024/2025</q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-icon name="download" color="grey" />
        </q-item-section>
      </q-item>

      <q-item clickable v-ripple @click="download('PDP')">
        <q-item-section avatar>
          <q-icon name="folder_shared" color="orange" />
        </q-item-section>
        <q-item-section>
          <q-item-label>Piano Didattico Personalizzato (PDP)</q-item-label>
          <q-item-label caption>Firmato digitalmente</q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-icon name="download" color="grey" />
        </q-item-section>
      </q-item>

    </q-list>
  </q-page>
</template>

<script setup>
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

function download(docName) {
  $q.notify({ type: 'positive', message: `Start download: ${docName}` })
}
</script>
