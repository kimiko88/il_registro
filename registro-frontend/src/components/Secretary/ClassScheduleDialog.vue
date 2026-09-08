<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(1200px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
      <q-card-section class="bg-primary text-white row items-center q-pa-md shrink-0">
        <div class="row items-center">
          <q-avatar color="white-20" text-color="white" icon="calendar_today" class="q-mr-sm" size="36px" />
          <div>
            <div class="text-h6 text-weight-bold">
              {{ $t('secretaryClasses.scheduleTitle') || 'Orario Settimanale - Classe' }} {{ currentClass?.name }}{{ currentClass?.section }}
            </div>
            <div class="text-subtitle2 text-white/80">{{ currentClass?.academic_year }}</div>
          </div>
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-pa-md col overflow-y-auto">
        <ScheduleGrid 
          :assignments="assignments"
          :initial-schedule="currentSchedule"
          :loading="loading"
          @save="$emit('save', $event)"
        />
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import ScheduleGrid from '@/components/Secretary/ScheduleGrid.vue'

defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  currentClass: {
    type: Object,
    default: null
  },
  assignments: {
    type: Array,
    default: () => []
  },
  currentSchedule: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['update:modelValue', 'save'])
</script>

<style scoped>
.rounded-xl { border-radius: 1rem; }
</style>
