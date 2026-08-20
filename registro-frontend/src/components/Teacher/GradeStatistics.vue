<template>
  <q-card class="q-mb-md">
    <q-card-section>
      <div class="text-h6">{{ t('gradesPage.statisticsTitle') || 'Statistiche Voti' }}</div>
    </q-card-section>
    <q-card-section class="row q-col-gutter-md">
      <div class="col-12 col-md-4 text-center">
        <div class="text-caption text-grey">{{ t('roleDashboards.gradeAverage') || 'Media Voti' }}</div>
        <div class="text-h4 text-primary">{{ average }}</div>
      </div>
      <div class="col-12 col-md-4 text-center">
        <div class="text-caption text-grey">{{ t('gradesPage.totalGrades') || 'Voti Assegnati' }}</div>
        <div class="text-h4">{{ totalGrades }}</div>
      </div>
      <div class="col-12 col-md-4">
        <!-- Placeholder for a mini chart using Quasar linear progress or just bars -->
        <div class="text-caption text-grey q-mb-sm">{{ t('gradesPage.distribution') || 'Distribuzione' }}</div>
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
import { useI18n } from 'vue-i18n';
import { useGradesStore } from '@/stores/grades';

const { t } = useI18n();

const props = defineProps({
  selectedSemester: {
    type: Number,
    default: 0
  }
});

const gradesStore = useGradesStore();

const average = computed(() => {
    const avgFn = gradesStore.classAverage;
    return typeof avgFn === 'function' ? avgFn(props.selectedSemester) : avgFn;
});
const totalGrades = computed(() => {
    if (!gradesStore.grades || !gradesStore.grades.students) return 0;
    return gradesStore.grades.students.reduce((acc, s) => {
        const studentGrades = s.grades || [];
        const filtered = props.selectedSemester > 0
            ? studentGrades.filter(g => !g.semester || Number(g.semester) === Number(props.selectedSemester))
            : studentGrades;
        return acc + filtered.length;
    }, 0);
});
</script>
