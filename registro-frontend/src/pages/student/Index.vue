<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center q-mb-md justify-between">
      <div>
        <h1 class="text-h4 q-my-none">Hi, {{ studentStore.profile?.firstName }}</h1>
        <div class="text-subtitle1 text-grey-8">{{ studentStore.className }}</div>
      </div>
      <q-avatar size="md">
        <img :src="studentStore.profile?.avatar || 'https://cdn.quasar.dev/img/boy-avatar.png'" />
      </q-avatar>
    </div>

    <!-- Quick Stats -->
    <div class="row q-col-gutter-sm q-mb-lg">
      <div class="col-6 col-sm-3">
        <q-card class="bg-primary text-white text-center q-py-sm">
          <div class="text-h6">7.8</div>
          <div class="text-caption">Average</div>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card class="bg-green text-white text-center q-py-sm">
          <div class="text-h6">92%</div>
          <div class="text-caption">Attendance</div>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card class="bg-orange text-white text-center q-py-sm">
          <div class="text-h6">30h</div>
          <div class="text-caption">PCTO</div>
        </q-card>
      </div>
      <div class="col-6 col-sm-3">
        <q-card class="bg-purple text-white text-center q-py-sm">
          <div class="text-h6">{{ studentStore.notifications.length }}</div>
          <div class="text-caption">Alerts</div>
        </q-card>
      </div>
    </div>

    <!-- Charts / Recent Activity -->
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-8">
        <q-card class="q-mb-md">
          <q-card-section>
            <div class="text-h6">Recent Grades</div>
          </q-card-section>
          <q-list separator>
             <q-item v-for="n in 3" :key="n">
                <q-item-section>
                    <q-item-label>Mathematics</q-item-label>
                    <q-item-label caption>Algebra Test</q-item-label>
                </q-item-section>
                <q-item-section side>
                    <q-badge color="green" label="8.0" />
                </q-item-section>
             </q-item>
          </q-list>
          <q-card-actions align="right">
            <q-btn flat label="View All" color="primary" to="/student/grades" />
          </q-card-actions>
        </q-card>
        
        <q-card>
            <q-card-section>
                <div class="text-h6">Upcoming</div>
            </q-card-section>
            <q-item>
                <q-item-section avatar>
                    <q-icon name="event" color="orange" />
                </q-item-section>
                <q-item-section>
                    <q-item-label>Physics Final</q-item-label>
                    <q-item-label caption>Tomorrow, 09:00</q-item-label>
                </q-item-section>
            </q-item>
        </q-card>
      </div>

      <!-- Quick Links -->
      <div class="col-12 col-md-4">
        <div class="row q-col-gutter-sm">
            <div class="col-6">
                <q-btn push color="white" text-color="black" icon="school" label="Grades" class="full-width q-py-md" to="/student/grades" />
            </div>
            <div class="col-6">
                <q-btn push color="white" text-color="black" icon="fact_check" label="Attend" class="full-width q-py-md" to="/student/attendance" />
            </div>
            <div class="col-6">
                <q-btn push color="white" text-color="black" icon="folder" label="Docs" class="full-width q-py-md" to="/student/documents" />
            </div>
            <div class="col-6">
                <q-btn push color="white" text-color="black" icon="mail" label="Msgs" class="full-width q-py-md" to="/student/communications" />
            </div>
            <div class="col-12">
                <q-btn push color="accent" text-color="white" icon="work" label="PCTO Dashboard" class="full-width q-py-sm" to="/student/pcto" />
            </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useStudentStore } from 'src/stores/student';

const studentStore = useStudentStore();

onMounted(() => {
    studentStore.fetchProfile();
    studentStore.fetchNotifications();
});
</script>
