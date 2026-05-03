<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center q-mb-xl justify-between">
      <div>
        <h1 class="text-h3 text-weight-bold text-outfit q-my-none text-gradient-premium">
          Pannello Docente
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">Benvenuto, Prof. {{ teacherStore.fullName }}</div>
      </div>
      <div class="text-right">
        <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">Oggi</div>
        <div class="text-h6 text-outfit text-weight-bold text-slate-700">{{ todayDate }}</div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="row q-col-gutter-lg q-mb-xl">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="dashboard-card glass-card bg-indigo-600 text-white shadow-soft overflow-hidden">
          <q-card-section>
            <div class="text-caption opacity-80 text-uppercase letter-spacing-1">Prossima Lezione</div>
            <div class="text-h4 text-weight-bold q-mt-sm">{{ nextLesson?.class_id || 'Nessuna' }}</div>
            <div class="text-caption q-mt-xs">{{ nextLesson?.subject_id || '-' }}</div>
          </q-card-section>
          <q-icon name="schedule" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="dashboard-card glass-card bg-orange-600 text-white shadow-soft overflow-hidden">
          <q-card-section>
            <div class="text-caption opacity-80 text-uppercase letter-spacing-1">Da Fare</div>
            <div class="text-h4 text-weight-bold q-mt-sm">{{ pendingTasksCount }} Revisioni</div>
            <div class="text-caption q-mt-xs">Giustificazioni in sospeso</div>
          </q-card-section>
          <q-icon name="pending_actions" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="dashboard-card glass-card bg-emerald-600 text-white shadow-soft overflow-hidden">
          <q-card-section>
            <div class="text-caption opacity-80 text-uppercase letter-spacing-1">Colloqui</div>
            <div class="text-h4 text-weight-bold q-mt-sm">{{ upcomingColloquiCount }} Prenotazioni</div>
            <div class="text-caption q-mt-xs">Controlla l'agenda</div>
          </q-card-section>
          <q-icon name="people" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="dashboard-card glass-card bg-violet-600 text-white shadow-soft overflow-hidden">
          <q-card-section>
            <div class="text-caption opacity-80 text-uppercase letter-spacing-1">Messaggi</div>
            <div class="text-h4 text-weight-bold q-mt-sm">{{ unreadMessagesCount }} Nuovi</div>
            <div class="text-caption q-mt-xs">Comunicazioni interne</div>
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

const nextLesson = computed(() => teacherStore.profile?.next_lesson || null);
const pendingTasksCount = computed(() => teacherStore.notifications.filter(n => n.type === 'action').length);
const upcomingColloquiCount = computed(() => teacherStore.profile?.colloqui_count || 0);
const unreadMessagesCount = computed(() => teacherStore.notifications.filter(n => !n.read).length);

onMounted(async () => {
  await Promise.all([
    teacherStore.fetchProfile(),
    teacherStore.fetchNotifications(),
    classesStore.fetchAssignedClasses()
  ]);
});
</script>

<style scoped>
.letter-spacing-1 {
    letter-spacing: 1px;
}

.dashboard-card {
  height: 160px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none !important;
}

.dashboard-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2);
}

.card-bg-icon {
  position: absolute;
  right: -20px;
  bottom: -20px;
  font-size: 120px;
  opacity: 0.15;
}

.bg-indigo-600 { background: linear-gradient(135deg, #4f46e5 0%, #6366f1 100%); }
.bg-orange-600 { background: linear-gradient(135deg, #ea580c 0%, #f97316 100%); }
.bg-emerald-600 { background: linear-gradient(135deg, #059669 0%, #10b981 100%); }
.bg-violet-600 { background: linear-gradient(135deg, #7c3aed 0%, #8b5cf6 100%); }
</style>
