<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div class="row items-center">
        <h1 class="text-h4 q-my-none q-mr-md">Attendance</h1>
        <ClassSelector />
      </div>
      <div>
        <q-input v-model="selectedDate" type="date" dense outlined class="inline-block" style="min-width: 150px" />
      </div>
    </div>

    <!-- Summary -->
    <AttendanceSummary v-if="classesStore.selectedClass" />

    <div v-if="classesStore.selectedClass">
        <!-- Toolbar -->
        <div class="row q-gutter-sm q-mb-md">
            <q-btn outline label="Mark All Present" color="green" size="sm" @click="markAll('Present')" />
            <q-space />
            <q-btn unelevated label="Save" color="primary" icon="save" :loading="isSaving" @click="save" />
        </div>

        <!-- List -->
        <AttendanceList />
    </div>
    <div v-else class="text-center q-pa-lg text-grey">
        Select a class to manage attendance
    </div>
  </q-page>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useClassesStore } from 'src/stores/classes';
import { useAttendanceStore } from 'src/stores/attendance';
import { useAttendanceMarking } from 'src/composables/useAttendanceMarking';
import { date } from 'quasar';
import ClassSelector from 'src/components/Teacher/ClassSelector.vue';
import AttendanceList from 'src/components/Teacher/AttendanceList.vue';
import AttendanceSummary from 'src/components/Teacher/AttendanceSummary.vue';

const classesStore = useClassesStore();
const attendanceStore = useAttendanceStore();
const { markAll, saveAttendance, isSaving } = useAttendanceMarking();

const selectedDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'));

const refreshData = () => {
    if (classesStore.selectedClassId) {
        attendanceStore.fetchDailyAttendance(classesStore.selectedClassId, selectedDate.value);
    }
};

const save = () => {
    if (classesStore.selectedClassId) {
        saveAttendance(classesStore.selectedClassId, selectedDate.value);
    }
};

watch(() => classesStore.selectedClassId, refreshData);
watch(selectedDate, refreshData);

onMounted(() => {
    if (classesStore.classes.length === 0) classesStore.fetchAssignedClasses();
    refreshData();
});
</script>
