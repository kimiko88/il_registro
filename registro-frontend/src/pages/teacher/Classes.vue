<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <h1 class="text-h4 q-my-none">Class Management</h1>
       <ClassSelector />
    </div>

    <div v-if="classesStore.selectedClass">
        <!-- Coordinator Panel -->
        <CoordinatorDashboard 
            v-if="classesStore.selectedClass.coordinator" 
            :class-id="classesStore.selectedClassId"
            :class-name="classesStore.selectedClass.name"
            class="q-mb-lg"
        />
        
        <!-- Tabs -->
        <q-tabs v-model="tab" align="left" class="text-primary q-mb-md">
            <q-tab name="students" label="Students" icon="people" />
            <q-tab name="info" label="Info & Groups" icon="info" />
        </q-tabs>

        <q-tab-panels v-model="tab" animated>
            <q-tab-panel name="students" class="q-pa-none">
                <StudentsList 
                    :students="classesStore.currentClassDetails?.students" 
                    v-if="classesStore.currentClassDetails"
                />
                <q-spinner v-else color="primary" size="3em" class="q-ma-md" />
            </q-tab-panel>

            <q-tab-panel name="info">
                <div class="text-h6">Class Information</div>
                <div class="row q-my-md">
                    <div class="col-12 col-md-6">
                        <q-list bordered separator>
                            <q-item>
                                <q-item-section>
                                    <q-item-label caption>Classroom</q-item-label>
                                    <q-item-label>Room 101 (First Floor)</q-item-label>
                                </q-item-section>
                            </q-item>
                            <q-item>
                                <q-item-section>
                                    <q-item-label caption>Representatives</q-item-label>
                                    <q-item-label>{{ classesStore.currentClassDetails?.representatives.join(', ') }}</q-item-label>
                                </q-item-section>
                            </q-item>
                        </q-list>
                    </div>
                </div>
            </q-tab-panel>
        </q-tab-panels>
    </div>
    <div v-else class="text-center text-grey q-pa-lg">
        Select a class
    </div>
  </q-page>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useClassesStore } from 'src/stores/classes';
import ClassSelector from 'src/components/Teacher/ClassSelector.vue';
import StudentsList from 'src/components/Teacher/StudentsList.vue';
import CoordinatorDashboard from 'src/components/Teacher/CoordinatorDashboard.vue';

const classesStore = useClassesStore();
const tab = ref('students');

watch(() => classesStore.selectedClassId, (newId) => {
    if (newId) classesStore.fetchClassDetails(newId);
});

onMounted(() => {
    if (classesStore.classes.length === 0) classesStore.fetchAssignedClasses();
    if (classesStore.selectedClassId) classesStore.fetchClassDetails(classesStore.selectedClassId);
});
</script>
