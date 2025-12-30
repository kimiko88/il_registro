<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center q-mb-md justify-between">
      <div>
        <h1 class="text-h4 q-my-none">Teacher Dashboard</h1>
        <div class="text-subtitle1 text-grey-8">Welcome back, {{ teacherStore.fullName }}</div>
      </div>
      <div class="text-right">
        <div class="text-caption text-grey">Today</div>
        <div class="text-h6">{{ todayDate }}</div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="dashboard-card bg-primary text-white">
          <q-card-section>
            <div class="text-subtitle2">Next Class</div>
            <div class="text-h5 q-mt-sm">5A - Math</div>
            <div class="text-caption">09:00 - 10:00 (Room 101)</div>
          </q-card-section>
          <q-icon name="schedule" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="dashboard-card bg-orange text-white">
          <q-card-section>
            <div class="text-subtitle2">Pending Actions</div>
            <div class="text-h5 q-mt-sm">3 Reviews</div>
            <div class="text-caption">2 Justifications, 1 Document</div>
          </q-card-section>
          <q-icon name="pending_actions" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="dashboard-card bg-teal text-white">
          <q-card-section>
            <div class="text-subtitle2">Upcoming Colloqui</div>
            <div class="text-h5 q-mt-sm">2 Bookings</div>
            <div class="text-caption">Tomorrow, 15:00</div>
          </q-card-section>
          <q-icon name="people" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="dashboard-card bg-deep-purple text-white">
          <q-card-section>
            <div class="text-subtitle2">Unread Messages</div>
            <div class="text-h5 q-mt-sm">{{ teacherStore.notifications.filter(n => !n.read).length }} New</div>
            <div class="text-caption">Check Communications</div>
          </q-card-section>
          <q-icon name="mail" class="card-bg-icon" />
        </q-card>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="row q-col-gutter-md">
      <!-- Quick Actions -->
      <div class="col-12 col-md-8">
        <q-card class="q-mb-md">
          <q-card-section>
            <div class="text-h6 q-mb-md">Quick Actions</div>
            <div class="row q-col-gutter-sm">
              <div class="col-6 col-sm-3">
                <q-btn
                  push color="primary" label="Grade Class" icon="grade"
                  class="full-width" size="lg" to="/teacher/grades"
                />
              </div>
              <div class="col-6 col-sm-3">
                <q-btn
                  push color="secondary" label="Attendance" icon="fact_check"
                  class="full-width" size="lg" to="/teacher/attendance"
                />
              </div>
              <div class="col-6 col-sm-3">
                <q-btn
                  push color="accent" label="Sign Doc" icon="draw"
                  class="full-width" size="lg" to="/teacher/documents"
                />
              </div>
              <div class="col-6 col-sm-3">
                <q-btn
                  push color="info" label="Message" icon="send"
                  class="full-width" size="lg" to="/teacher/communications"
                />
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Recent Activity / Notifications -->
        <q-card>
          <q-card-section>
            <div class="text-h6">Notifications & Activity</div>
          </q-card-section>
          <q-list separator>
            <q-item v-for="note in teacherStore.notifications" :key="note.id" clickable v-ripple>
              <q-item-section avatar>
                <q-icon :name="note.type === 'info' ? 'info' : 'check_circle'" :color="note.type === 'info' ? 'blue' : 'green'" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ note.title }}</q-item-label>
                <q-item-label caption>{{ note.message }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-btn flat round icon="close" size="sm" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Side Panel: Class List Shortcut -->
      <div class="col-12 col-md-4">
        <q-card class="full-height">
          <q-card-section class="bg-grey-2">
            <div class="text-h6">My Classes</div>
          </q-card-section>
          <q-list separator>
            <q-item v-for="cls in classesStore.classes" :key="cls.id" clickable @click="classesStore.selectClass(cls.id)" :active="cls.id === classesStore.selectedClassId" active-class="bg-blue-1 text-primary">
              <q-item-section avatar>
                <q-avatar color="primary" text-color="white" size="sm">{{ cls.name }}</q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ cls.id }}</q-item-label>
                <q-item-label caption>{{ cls.type }} - {{ cls.studentsCount }} students</q-item-label>
              </q-item-section>
              <q-item-section side v-if="cls.coordinator">
                <q-icon name="star" color="orange" title="Coordinator" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted, computed } from 'vue';
import { useTeacherStore } from 'src/stores/teacher';
import { useClassesStore } from 'src/stores/classes';
import { date } from 'quasar';

const teacherStore = useTeacherStore();
const classesStore = useClassesStore();

const todayDate = computed(() => date.formatDate(Date.now(), 'DD MMMM YYYY'));

onMounted(async () => {
  await Promise.all([
    teacherStore.fetchProfile(),
    teacherStore.fetchNotifications(),
    classesStore.fetchAssignedClasses()
  ]);
});
</script>

<style scoped>
.dashboard-card {
  height: 140px;
  position: relative;
  overflow: hidden;
  transition: transform 0.2s;
}
.dashboard-card:hover {
  transform: translateY(-5px);
}
.card-bg-icon {
  position: absolute;
  right: -20px;
  bottom: -20px;
  font-size: 100px;
  opacity: 0.2;
}
</style>
