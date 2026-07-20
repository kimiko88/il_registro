<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none">📅 Agenda & Compiti</h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none">Calendario scadenze, verifiche, ed esercitazioni con tracciamento completamento</p>
      </div>
      <q-btn-toggle
        v-model="filterType"
        toggle-color="primary"
        :options="[
          { label: 'Tutti', value: 'all' },
          { label: 'Compiti', value: 'compito' },
          { label: 'Verifiche', value: 'verifica' }
        ]"
      />
    </div>

    <div class="row q-col-gutter-md">
      <div v-for="item in filteredHomeworks" :key="item.id" class="col-12 col-md-6 col-lg-4">
        <q-card flat bordered :class="{ 'bg-grey-2': item.completed }">
          <q-card-section class="row items-center justify-between">
            <q-chip
              :color="item.type === 'verifica' ? 'negative' : 'primary'"
              text-color="white"
              size="sm"
            >
              {{ (item.type || 'compito').toUpperCase() }}
            </q-chip>
            <div class="text-caption text-grey">Scadenza: {{ item.due_date }}</div>
          </q-card-section>

          <q-card-section>
            <div class="text-h6 text-weight-bold" :class="{ 'text-strike text-grey': item.completed }">
              {{ item.description }}
            </div>
            <div class="text-subtitle2 text-grey-7">Docente: {{ item.teacher_name }}</div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="between">
            <q-checkbox
              v-model="item.completed"
              label="Segna come fatto"
              color="positive"
            />
            <q-btn flat round icon="event" color="grey" />
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import lessonService from '@/services/lessonService'

const filterType = ref('all')
const homeworks = ref([
  { id: '1', type: 'compito', description: 'Pagina 394 esercizi 1-5', due_date: '2026-07-22', teacher_name: 'Prof. Severus Piton', completed: false },
  { id: '2', type: 'verifica', description: 'Verifica Scritta sulle Pozioni', due_date: '2026-07-25', teacher_name: 'Prof. Severus Piton', completed: false }
])

const filteredHomeworks = computed(() => {
  if (filterType.value === 'all') return homeworks.value
  return homeworks.value.filter(h => h.type === filterType.value)
})
</script>
