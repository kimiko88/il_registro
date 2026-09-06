<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">{{ t('schedulerPage.title') || 'Schedule Generator' }}</div>

    <div class="row q-col-gutter-md">
      <!-- Configuration Panel -->
      <div class="col-12 col-md-4">
        <q-card>
          <q-card-section>
            <div class="text-h6">{{ t('schedulerPage.constraints') || 'Constraints' }}</div>
          </q-card-section>

          <q-card-section>
            <q-input
              v-model="config.startDate"
              type="date"
              :label="t('schedulerPage.startDate') || 'Start Date'"
              filled
              class="q-mb-sm"
            />
            <div class="row q-col-gutter-sm">
              <div class="col">
                <q-input
                  v-model.number="config.startHour"
                  type="number"
                  :label="t('schedulerPage.startHour') || 'Start Hour'"
                  filled
                />
              </div>
              <div class="col">
                <q-input
                  v-model.number="config.endHour"
                  type="number"
                  :label="t('schedulerPage.endHour') || 'End Hour'"
                  filled
                />
              </div>
            </div>
            
            <q-select
              v-model="config.days"
              multiple
              :options="dayOptions"
              :label="t('schedulerPage.schoolDays') || 'School Days'"
              filled
              class="q-mt-sm"
              emit-value
              map-options
            />
          </q-card-section>

          <q-card-actions align="right">
            <q-btn
              :label="t('schedulerPage.generate') || 'Generate'"
              color="primary"
              @click="handleGenerate"
              :loading="store.isLoading"
            />
          </q-card-actions>
        </q-card>
        
        <!-- Conflicts Panel -->
        <q-card class="q-mt-md" v-if="store.conflicts.length > 0">
          <q-card-section class="bg-negative text-white">
            <div class="text-h6">
              {{ t('schedulerPage.conflictsDetected') || 'Conflicts Detected' }} ({{ store.conflicts.length }})
            </div>
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
            <div class="text-h6">{{ t('schedulerPage.preview') || 'Preview' }}</div>
          </q-card-section>
          
          <q-card-section>
            <div v-if="store.generatedSchedule.length === 0" class="text-center text-grey q-pa-lg">
              {{ t('schedulerPage.noScheduleYet') || 'No schedule generated yet. Configure constraints and click Generate.' }}
            </div>
            
            <div v-else class="schedule-grid">
               <!-- Simplified list for MVP visualization -->
               <q-table
                 :title="t('schedulerPage.generatedSlots') || 'Generated Slots'"
                 :rows="store.generatedSchedule"
                 :columns="columns"
                 row-key="id"
                 flat
                 bordered
               />
            </div>
          </q-card-section>
          
          <q-card-actions align="right" v-if="store.generatedSchedule.length > 0">
             <q-btn
               :label="t('schedulerPage.validate') || 'Validate'"
               flat
               color="warning"
               @click="handleValidate"
               :loading="store.isLoading"
             />
             <q-btn
               :label="t('schedulerPage.saveSchedule') || 'Save Schedule'"
               color="positive"
               icon="save"
             />
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSchedulingStore } from '@/stores/scheduling'

const { t } = useI18n()
const store = useSchedulingStore()

const config = ref({
  startDate: new Date().toISOString().split('T')[0],
  startHour: 8,
  endHour: 14,
  days: [1, 2, 3, 4, 5]
})

const dayOptions = computed(() => [
  { label: t('schedulerPage.mon') || 'Monday', value: 1 },
  { label: t('schedulerPage.tue') || 'Tuesday', value: 2 },
  { label: t('schedulerPage.wed') || 'Wednesday', value: 3 },
  { label: t('schedulerPage.thu') || 'Thursday', value: 4 },
  { label: t('schedulerPage.fri') || 'Friday', value: 5 },
  { label: t('schedulerPage.sat') || 'Saturday', value: 6 }
])

const columns = computed(() => [
  { name: 'day', label: t('schedulerPage.day') || 'Day', field: 'day', sortable: true },
  { name: 'hour', label: t('schedulerPage.hour') || 'Hour', field: 'hour', sortable: true },
  { name: 'class', label: t('schedulerPage.class') || 'Class', field: 'class_id', sortable: true },
  { name: 'teacher', label: t('schedulerPage.teacher') || 'Teacher', field: 'teacher_id', sortable: true },
  { name: 'subject', label: t('schedulerPage.subject') || 'Subject', field: 'subject_id', sortable: true },
])

async function handleGenerate() {
  await store.generateSchedule(config.value)
}

async function handleValidate() {
  await store.validateSchedule(store.generatedSchedule)
}
</script>
