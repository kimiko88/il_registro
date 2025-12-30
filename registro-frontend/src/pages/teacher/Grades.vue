<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div class="row items-center">
        <h1 class="text-h4 q-my-none q-mr-md">Grades</h1>
        <ClassSelector />
      </div>
      <div>
        <q-select
          v-model="selectedSubject"
          :options="gradesStore.subjects"
          label="Subject"
          dense
          outlined
          class="inline-block"
          style="min-width: 150px"
        />
      </div>
    </div>

    <!-- Stats -->
    <GradeStatistics v-if="classesStore.selectedClass" />

    <!-- Main Table -->
    <GradeTable 
      :students="students" 
      :loading="loading"
      @add-grade="openGradeForm"
    />

    <!-- Entry Dialog -->
    <GradeForm
      v-model="showAddDialog"
      :student-id="selectedStudent?.id"
      :student-name="selectedStudent?.name"
      @saved="refreshGrades"
    />
  </q-page>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import { useClassesStore } from 'src/stores/classes';
import { useGradesStore } from 'src/stores/grades';
import ClassSelector from 'src/components/Teacher/ClassSelector.vue';
import GradeTable from 'src/components/Teacher/GradeTable.vue';
import GradeStatistics from 'src/components/Teacher/GradeStatistics.vue';
import GradeForm from 'src/components/Teacher/GradeForm.vue';

const classesStore = useClassesStore();
const gradesStore = useGradesStore();

const selectedSubject = ref(gradesStore.subjects[0]);
const showAddDialog = ref(false);
const selectedStudent = ref(null);

// Computeds
const students = computed(() => {
  // In real app, this comes from classesStore.currentClassDetails.students
  // Mocking for now since classesStore mock doesn't have students list yet
  if (!classesStore.selectedClassId) return [];
  return [
    { id: 's1', name: 'Giuseppe Verdi' },
    { id: 's2', name: 'Mario Rossi' },
    { id: 's3', name: 'Sofia Bianchi' }
  ];
});

const loading = computed(() => classesStore.loading || gradesStore.loading);

// Methods
const openGradeForm = (student) => {
  selectedStudent.value = student;
  showAddDialog.value = true;
};

const refreshGrades = () => {
  if (classesStore.selectedClassId && selectedSubject.value) {
    gradesStore.fetchGrades(classesStore.selectedClassId, selectedSubject.value);
  }
};

// Watchers
watch(() => classesStore.selectedClassId, refreshGrades);
watch(selectedSubject, refreshGrades);

onMounted(() => {
  if (classesStore.classes.length === 0) classesStore.fetchAssignedClasses();
  refreshGrades();
});
</script>

<style scoped>
.inline-block {
  display: inline-block;
}
</style>
