<template>
  <q-page class="q-pa-md">
    <div class="row justify-between items-center q-mb-md">
        <h1 class="text-h4 q-my-none">My Grades</h1>
        <q-btn round flat icon="file_download" color="primary" @click="exportReport">
             <q-tooltip>Export PDF</q-tooltip>
        </q-btn>
    </div>

    <!-- Filters -->
    <div class="row q-gutter-md q-mb-md">
        <q-select v-model="semester" :options="['Semester 1', 'Semester 2']" dense outlined class="col" label="Period" />
        <q-select v-model="filterType" :options="['All', 'Written', 'Oral', 'Lab']" dense outlined class="col" label="Type" />
    </div>

    <!-- Chart -->
    <GradeChart :averages="averages" v-if="!loading" />

    <!-- List -->
    <GradesList 
        :grades-by-subject="gradesBySubject" 
        :averages="averages"
        :get-trend="getTrend"
        v-if="!loading"
    />

    <div v-if="loading" class="flex flex-center q-pa-xl">
        <q-spinner color="primary" size="3em" />
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useGradesStore } from 'src/stores/grades';
import { useMyGrades } from 'src/composables/useMyGrades';
import { useQuasar } from 'quasar';
import GradesList from 'src/components/Student/GradesList.vue';
import GradeChart from 'src/components/Student/GradeChart.vue';

const store = useGradesStore();
const { gradesBySubject, averages, getTrend, loading } = useMyGrades();
const $q = useQuasar();

const semester = ref('Semester 1');
const filterType = ref('All');

const exportReport = () => {
    $q.notify({ type: 'positive', message: 'Generating PDF...' });
    setTimeout(() => {
        $q.notify({ type: 'info', message: 'Download started' });
    }, 1000);
};

onMounted(() => {
    store.fetchMyGrades('me');
});
</script>
