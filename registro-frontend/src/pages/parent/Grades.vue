<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
        <div>
            <h1 class="text-h4 q-my-none">Grades</h1>
            <div class="text-subtitle1 text-grey" v-if="selectedChild">
                for {{ selectedChild.firstName }}
            </div>
        </div>
        <ChildrenSelector />
    </div>

    <!-- Reusing Student Components for consistency -->
    <GradeChart :averages="averages" v-if="!loading" />
    <GradesList 
        :grades-by-subject="gradesBySubject" 
        :averages="averages"
        :get-trend="getTrend"
        v-if="!loading"
    />
  </q-page>
</template>

<script setup>
import { useChildGrades } from 'src/composables/useChildGrades';
import ChildrenSelector from 'src/components/Parent/ChildrenSelector.vue';
import GradesList from 'src/components/Student/GradesList.vue'; // Reuse!
import GradeChart from 'src/components/Student/GradeChart.vue'; // Reuse!

const { selectedChild, gradesBySubject, averages, getTrend, loading } = useChildGrades();
</script>
