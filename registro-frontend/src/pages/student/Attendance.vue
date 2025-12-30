<template>
  <q-page class="q-pa-md">
    <h1 class="text-h4 q-my-none q-mb-md">Attendance</h1>

    <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-sm-6">
            <AttendanceChart :stats="stats" />
        </div>
        <div class="col-12 col-sm-6">
            <q-card class="bg-orange-1 text-center q-py-md full-height">
                <div class="text-h6">{{ stats.absent }} Absences</div>
                <div class="text-caption q-mt-sm">Justify promptly to avoid issues</div>
            </q-card>
        </div>
    </div>

    <!-- History -->
    <q-list bordered separator class="bg-white rounded-borders">
        <q-item-label header>Recent Activity</q-item-label>
        <q-item v-for="record in records" :key="record.date">
             <q-item-section avatar>
                 <q-icon :name="getStatusIcon(record.status)" :color="getStatusColor(record.status)" />
             </q-item-section>
             <q-item-section>
                 <q-item-label>{{ formatDate(record.date) }}</q-item-label>
                 <q-item-label caption>{{ record.status }} <span v-if="record.time">({{ record.time }})</span></q-item-label>
             </q-item-section>
             <q-item-section side>
                 <q-btn 
                    v-if="record.status === 'Absent' && !record.justificationStatus" 
                    dense flat label="Justify" color="primary" 
                    @click="openJustify(record.date)"
                 />
                 <q-badge v-else-if="record.justificationStatus" color="orange" :label="record.justificationStatus" />
             </q-item-section>
        </q-item>
    </q-list>
    
    <JustificationForm v-model="showJustify" :date="selectedDate" @submit="handleJustify" />
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAttendanceStore } from 'src/stores/attendance';
import { useMyAttendance } from 'src/composables/useMyAttendance';
import { date } from 'quasar';
import AttendanceChart from 'src/components/Student/AttendanceChart.vue';
import JustificationForm from 'src/components/Student/JustificationForm.vue';

const store = useAttendanceStore();
const { stats, records, requestJustification } = useMyAttendance();

const showJustify = ref(false);
const selectedDate = ref('');

const openJustify = (date) => {
    selectedDate.value = date;
    showJustify.value = true;
};

const handleJustify = async (data) => {
    await requestJustification(selectedDate.value, data.reason);
};

const getStatusIcon = (status) => {
    if (status === 'Present') return 'check_circle';
    if (status === 'Absent') return 'cancel';
    return 'schedule';
};

const getStatusColor = (status) => {
    if (status === 'Present') return 'green';
    if (status === 'Absent') return 'red';
    return 'orange';
};

const formatDate = (val) => date.formatDate(val, 'dddd, DD MMMM');

onMounted(() => {
    store.fetchMyAttendance();
});
</script>
