<template>
  <q-btn-dropdown
    flat
    dense
    no-caps
    color="primary"
    icon="face"
    :label="selectedChild ? selectedChild.name : 'Seleziona Figlio'"
  >
    <q-list style="min-width: 200px">
      <q-item-label header>Figli Collegati</q-item-label>
      <q-item
        v-for="child in children"
        :key="child.id"
        clickable
        v-close-popup
        @click="switchChild(child)"
        :active="selectedChild && selectedChild.id === child.id"
      >
        <q-item-section avatar>
          <q-avatar color="primary" text-color="white" icon="person" size="sm" />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ child.name }}</q-item-label>
          <q-item-label caption>Classe {{ child.class_name }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
  </q-btn-dropdown>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import userService from '@/services/userService'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const children = ref([])
const selectedChild = ref(null)

const loadChildren = async () => {
  try {
    const res = await userService.getChildren()
    children.value = res.data || []
    if (children.value.length > 0) {
      selectedChild.value = children.value[0]
    }
  } catch (err) {
    // Demo fallback
    children.value = [
      { id: 'std-1', name: 'Mario Rossi', class_name: '3A' },
      { id: 'std-2', name: 'Lucia Rossi', class_name: '1B' }
    ]
    selectedChild.value = children.value[0]
  }
}

const switchChild = async (child) => {
  selectedChild.value = child
  try {
    await userService.switchChild(child.id)
    $q.notify({ type: 'positive', message: `Profilo cambiato in ${child.name}` })
    window.location.reload()
  } catch (err) {
    $q.notify({ type: 'info', message: `Attivo profilo ${child.name}` })
  }
}

onMounted(() => {
  loadChildren()
})
</script>
