<template>
  <div>
    <q-card v-for="(grades, subject) in gradesBySubject" :key="subject" class="q-mb-md">
      <q-expansion-item expand-separator>
        <template v-slot:header>
            <q-item-section>
                <div class="text-h6">{{ subject }}</div>
            </q-item-section>
            
            <q-item-section side>
                <div class="row items-center">
                    <span class="text-h6 text-weight-bold q-mr-sm" :class="getAverageColor(averages[subject])">{{ averages[subject] }}</span>
                    <q-icon :name="getTrendIcon(getTrend(subject))" :color="getTrendColor(getTrend(subject))" />
                </div>
            </q-item-section>
        </template>

        <q-card-section class="q-pa-none">
            <q-list separator>
                <q-item v-for="grade in grades" :key="grade.id">
                    <q-item-section avatar>
                         <q-badge :color="getGradeColor(grade.value)" class="text-subtitle1 q-pa-xs">
                             {{ grade.value }}
                         </q-badge>
                    </q-item-section>
                    <q-item-section>
                        <q-item-label>{{ grade.description }}</q-item-label>
                        <q-item-label caption>{{ grade.type }}</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <q-item-label caption>{{ grade.date }}</q-item-label>
                    </q-item-section>
                </q-item>
            </q-list>
        </q-card-section>
      </q-expansion-item>
    </q-card>
  </div>
</template>

<script setup>
defineProps(['gradesBySubject', 'averages', 'getTrend']);

const getGradeColor = (val) => {
    if (val >= 8) return 'green';
    if (val >= 6) return 'blue';
    return 'red';
};

const getAverageColor = (val) => {
    return val >= 6 ? 'text-green' : 'text-red';
};

const getTrendIcon = (trend) => {
    if (trend === 'up') return 'trending_up';
    if (trend === 'down') return 'trending_down';
    return 'trending_flat';
};

const getTrendColor = (trend) => {
    if (trend === 'up') return 'green';
    if (trend === 'down') return 'red';
    return 'grey';
};
</script>
