<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
        <div>
            <h1 class="text-h4 q-my-none">Attendance</h1>
             <div class="text-subtitle1 text-grey" v-if="selectedChild">
                for {{ selectedChild.firstName }}
            </div>
        </div>
        <ChildrenSelector />
    </div>

    <!-- Reusing Student Components -->
    <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-sm-6">
            <AttendanceChart :stats="stats" />
        </div>
        <div class="col-12 col-sm-6">
             <!-- Different view for parent? Maybe just summary -->
             <q-card class="bg-red-1 text-center q-py-md full-height" v-if="stats.absent > 5">
                <div class="text-h6 text-red">Attention Required</div>
                <div>High number of absences</div>
             </q-card>
             <q-card class="bg-green-1 text-center q-py-md full-height" v-else>
                <div class="text-h6 text-green">Good Standing</div>
                <div>Regular attendance</div>
             </q-card>
        </div>
    </div>

    <q-list bordered separator>
         <q-item-label header>Recent Activity</q-item-label>
         <q-item v-for="record in records" :key="record.date">
             <q-item-section avatar>
                 <q-icon name="circle" :color="record.status === 'Present' ? 'green' : 'red'" size="xs" />
             </q-item-section>
             <q-item-section>
                 <q-item-label>{{ record.date }}</q-item-label>
                 <q-item-label caption>{{ record.status }}</q-item-label>
             </q-item-section>
             <q-item-section side>
                 <q-badge v-if="record.justificationStatus" :label="record.justificationStatus" />
                 <span v-else-if="record.status === 'Absent'" class="text-red pointer" @click="justify(record)">Justify needed</span>
             </q-item-section>
         </q-item>
    </q-list>
  </q-page>
</template>

<script setup>
import { useChildAttendance } from 'src/composables/useChildAttendance';
import ChildrenSelector from 'src/components/Parent/ChildrenSelector.vue';
import AttendanceChart from 'src/components/Student/AttendanceChart.vue';
import { useQuasar } from 'quasar';

const { selectedChild, stats, records } = useChildAttendance();
const $q = useQuasar();

const justify = (record) => {
    $q.notify({ type: 'info', message: 'Justification flow mock' });
};
</script>
