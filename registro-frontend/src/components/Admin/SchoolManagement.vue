<template>
  <div class="q-pa-md">
    <div class="row q-mb-md justify-between items-center">
      <div class="text-h6">School Management</div>
      <q-btn color="primary" icon="add" label="Add School" />
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
            <q-btn flat round icon="edit" size="sm" />
            <q-btn flat round icon="delete" color="negative" size="sm" />
          </div>
        </q-item-section>
      </q-item>
    </q-list>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const schools = ref([
  { id: 1, name: 'Liceo Scientifico', address: 'Via Roma 1' },
  { id: 2, name: 'Istituto Tecnico', address: 'Via Milano 2' }
])
</script>
