<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">Schedule Generator</div>

    <div class="row q-col-gutter-md">
      <!-- Configuration Panel -->
      <div class="col-12 col-md-4">
        <q-card>
          <q-card-section>
            <div class="text-h6">Constraints</div>
          </q-card-section>

          <q-card-section>
            <q-input v-model="config.startDate" type="date" label="Start Date" filled class="q-mb-sm" />
            <div class="row q-col-gutter-sm">
              <div class="col">
                <q-input v-model.number="config.startHour" type="number" label="Start Hour" filled />
              </div>
              <div class="col">
                <q-input v-model.number="config.endHour" type="number" label="End Hour" filled />
              </div>
            </div>
            
            <q-select
              v-model="config.days"
              multiple
              :options="dayOptions"
              label="School Days"
              filled
              class="q-mt-sm"
              emit-value
              map-options
            />
          </q-card-section>

          <q-card-actions align="right">
            <q-btn label="Generate" color="primary" @click="handleGenerate" :loading="store.isLoading" />
          </q-card-actions>
        </q-card>
        
        <!-- Conflicts Panel -->
        <q-card class="q-mt-md" v-if="store.conflicts.length > 0">
          <q-card-section class="bg-negative text-white">
            <div class="text-h6">Conflicts Detected ({{ store.conflicts.length }})</div>
          </q-card-section>
          <q-list separator>
            <q-item v-for="(c, i) in store.conflicts" :key="i">
              <q-item-section avatar>
                <q-icon name="warning" color="negative" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ c.type }}</q-item-label>
                <q-item-label caption>{{ c.description }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Result Grid -->
      <div class="col-12 col-md-8">
        <q-card>
          <q-card-section>
            <div class="text-h6">Preview</div>
          </q-card-section>
          
          <q-card-section>
            <div v-if="store.generatedSchedule.length === 0" class="text-center text-grey q-pa-lg">
              No schedule generated yet. Configure constraints and click Generate.
            </div>
            
            <div v-else class="schedule-grid">
               <!-- Simplified list for MVP visualization -->
               <q-table
                 title="Generated Slots"
                 :rows="store.generatedSchedule"
                 :columns="columns"
                 row-key="id"
                 flat
                 bordered
               />
            </div>
          </q-card-section>
          
          <q-card-actions align="right" v-if="store.generatedSchedule.length > 0">
             <q-btn label="Validate" flat color="warning" @click="handleValidate" :loading="store.isLoading"/>
             <q-btn label="Save Schedule" color="positive" icon="save" />
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useSchedulingStore } from 'src/stores/scheduling'

const store = useSchedulingStore()

const config = ref({
  startDate: new Date().toISOString().split('T')[0],
  startHour: 8,
  endHour: 14,
  days: [1, 2, 3, 4, 5]
})

const dayOptions = [
  { label: 'Monday', value: 1 },
  { label: 'Tuesday', value: 2 },
  { label: 'Wednesday', value: 3 },
  { label: 'Thursday', value: 4 },
  { label: 'Friday', value: 5 },
  { label: 'Saturday', value: 6 }
]

const columns = [
  { name: 'day', label: 'Day', field: 'day', sortable: true },
  { name: 'hour', label: 'Hour', field: 'hour', sortable: true },
  { name: 'class', label: 'Class', field: 'class_id', sortable: true },
  { name: 'teacher', label: 'Teacher', field: 'teacher_id', sortable: true },
  { name: 'subject', label: 'Subject', field: 'subject_id', sortable: true },
]

async function handleGenerate() {
  await store.generateSchedule(config.value)
}

async function handleValidate() {
  await store.validateSchedule(store.generatedSchedule)
}
</script>
