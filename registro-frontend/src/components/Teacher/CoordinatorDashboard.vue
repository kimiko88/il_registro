<template>
  <q-card class="bg-indigo-1">
    <q-card-section>
        <div class="text-h6 text-indigo-9">Coordinator Dashboard - {{ className }}</div>
    </q-card-section>

    <q-card-section>
        <div class="row q-col-gutter-md">
            <!-- At Risk Students -->
            <div class="col-12 col-md-6">
                <q-card flat bordered class="bg-white">
                    <q-card-section class="text-subtitle2 text-red">At Risk Students</q-card-section>
                    <q-list dense>
                        <q-item v-for="student in problems" :key="student.id">
                            <q-item-section avatar>
                                <q-icon name="warning" color="red" />
                            </q-item-section>
                            <q-item-section>
                                <q-item-label>{{ student.name }}</q-item-label>
                                <q-item-label caption>{{ student.issue }}</q-item-label>
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-card>
            </div>

            <!-- Class Stats -->
            <div class="col-12 col-md-6">
                <q-card flat bordered class="bg-white">
                     <q-card-section class="text-subtitle2">Class Overview</q-card-section>
                     <q-list dense>
                        <q-item>
                            <q-item-section>Attendance Average</q-item-section>
                            <q-item-section side>92%</q-item-section>
                        </q-item>
                        <q-item>
                            <q-item-section>Grade Average</q-item-section>
                            <q-item-section side>7.4</q-item-section>
                        </q-item>
                     </q-list>
                </q-card>
            </div>
        </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { computed } from 'vue';
import { useCoordination } from 'src/composables/useCoordination';

const props = defineProps({
    classId: String,
    className: String
});

const { getProblemStudents } = useCoordination();
const problems = computed(() => getProblemStudents(props.classId));
</script>
