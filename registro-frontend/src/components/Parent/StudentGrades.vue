<template>
  <div class="q-pa-md">
    <div class="text-h6 q-mb-md">My Children</div>
    
    <div class="row q-gutter-md">
      <q-card v-for="child in children" :key="child.id" class="my-card">
        <q-card-section>
          <div class="text-h6">{{ child.name }}</div>
          <div class="text-subtitle2">Class: {{ child.class }}</div>
        </q-card-section>

        <q-separator />

        <q-card-actions vertical>
          <q-btn flat color="primary" label="View Grades" @click="viewGrades(child)" />
          <q-btn flat color="primary" label="View Attendance" @click="viewAttendance(child)" />
        </q-card-actions>
      </q-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const children = ref([
  { id: 1, name: 'Mario Rossi', class: '1A' },
  { id: 2, name: 'Anna Rossi', class: '3B' }
])

function viewGrades(child) {
  console.log('View grades for', child.name)
}

function viewAttendance(child) {
  console.log('View attendance for', child.name)
}
</script>

<style scoped>
.my-card {
  width: 100%;
  max-width: 300px;
}
</style>
