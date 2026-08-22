<template>
  <q-btn-dropdown
    v-if="store.children.length > 0"
    flat
    dense
    no-caps
    color="primary"
    icon="face"
    :label="selectedChildName"
  >
    <q-list style="min-width: 200px">
      <q-item-label header>{{ t('competenciesPage.student') || 'Figli' }}</q-item-label>
      <q-item
        v-for="child in store.children"
        :key="child.id"
        clickable
        v-close-popup
        @click="store.selectChild(child.id)"
        :active="store.selectedChildId === child.id"
      >
        <q-item-section avatar>
          <q-avatar color="primary" text-color="white" icon="person" size="sm" />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ child.firstName }} {{ child.lastName }}</q-item-label>
          <q-item-label caption v-if="child.className">{{ t('udaPage.classLabel') || 'Classe' }} {{ child.className }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
  </q-btn-dropdown>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useChildrenStore } from '@/stores/children'

const { t } = useI18n()
const store = useChildrenStore()

const selectedChildName = computed(() => {
  const c = store.selectedChild
  if (!c) return t('competenciesPage.student') || 'Seleziona Figlio'
  return `${c.firstName || ''} ${c.lastName || ''}`.trim() || c.id
})

onMounted(() => {
  if (store.children.length === 0) {
    store.fetchChildren()
  }
})
</script>
