<template>
  <q-card class="q-mb-md">
    <q-card-section>
      <div class="text-h6">Class Statistics</div>
    </q-card-section>
    <q-card-section class="row q-col-gutter-md">
      <div class="col-12 col-md-4 text-center">
        <div class="text-caption text-grey">Class Average</div>
        <div class="text-h4 text-primary">{{ average }}</div>
      </div>
      <div class="col-12 col-md-4 text-center">
        <div class="text-caption text-grey">Tests Graded</div>
        <div class="text-h4">{{ totalGrades }}</div>
      </div>
      <div class="col-12 col-md-4">
        <!-- Placeholder for a mini chart using Quasar linear progress or just bars -->
        <div class="text-caption text-grey q-mb-sm">Distribution</div>
        <div class="row items-center q-mb-xs">
          <span class="text-xs q-mr-xs">9-10</span>
          <q-linear-progress :value="0.2" color="green" class="col" />
        </div>
        <div class="row items-center q-mb-xs">
          <span class="text-xs q-mr-xs">6-8</span>
          <q-linear-progress :value="0.5" color="blue" class="col" />
        </div>
        <div class="row items-center">
          <span class="text-xs q-mr-xs">&lt;6</span>
          <q-linear-progress :value="0.3" color="red" class="col" />
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { computed } from 'vue';
import { useGradesStore } from 'src/stores/grades';

const gradesStore = useGradesStore();

const average = computed(() => gradesStore.classAverage);
const totalGrades = computed(() => {
    if (!gradesStore.grades || !gradesStore.grades.students) return 0;
    return gradesStore.grades.students.reduce((acc, s) => acc + s.grades.length, 0);
});
</script>
