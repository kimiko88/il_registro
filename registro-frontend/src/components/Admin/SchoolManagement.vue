<template>
  <div class="q-pa-md">
    <div class="row q-mb-md justify-between items-center">
      <div class="text-h6">{{ t('roleDashboards.schoolManagement') || 'Gestione Istituti' }}</div>
      <q-btn color="primary" icon="add" :label="t('common.add') || 'Aggiungi Scuola'" @click="$emit('add')" />
    </div>

    <q-list bordered separator>
      <q-item v-for="school in schools" :key="school.id">
        <q-item-section avatar>
          <q-icon name="school" color="primary" />
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ school.name }}</q-item-label>
          <q-item-label caption>{{ school.address }}</q-item-label>
        </q-item-section>
        <q-item-section side>
          <div class="row q-gutter-xs">
            <q-btn flat round icon="edit" size="sm" :aria-label="t('common.edit') || 'Modifica'" @click="$emit('edit', school)">
              <q-tooltip>{{ t('common.edit') || 'Modifica' }}</q-tooltip>
            </q-btn>
            <q-btn flat round icon="delete" color="negative" size="sm" :aria-label="t('common.delete') || 'Elimina'" @click="$emit('delete', school)">
              <q-tooltip>{{ t('common.delete') || 'Elimina' }}</q-tooltip>
            </q-btn>
          </div>
        </q-item-section>
      </q-item>
    </q-list>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
defineEmits(['add', 'edit', 'delete'])

const schools = ref([
  { id: 1, name: 'Liceo Scientifico', address: 'Via Roma 1' },
  { id: 2, name: 'Istituto Tecnico', address: 'Via Milano 2' }
])
</script>
