<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-lg">Teacher Dashboard</div>
    
    <q-tabs
      v-model="tab"
      dense
      class="text-grey"
      active-color="primary"
      indicator-color="primary"
      align="justify"
      narrow-indicator
    >
      <q-tab name="grades" label="Grades" />
      <q-tab name="attendance" label="Attendance" />
    </q-tabs>

    <q-separator />

    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="grades">
        <GradeEntry />
      </q-tab-panel>

      <q-tab-panel name="attendance">
        <AttendanceMarking />
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import GradeEntry from '@/components/Teacher/GradeEntry.vue'
import AttendanceMarking from '@/components/Teacher/AttendanceMarking.vue'

const tab = ref('grades')
</script>
