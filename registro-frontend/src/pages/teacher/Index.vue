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

    <!-- ⚡ Quick Sign Banner (Firma Ora Corrente con 1 Click) -->
    <q-card class="q-mb-xl shadow-lg rounded-xl overflow-hidden text-white" style="background: linear-gradient(135deg, #1e1b4b 0%, #3730a3 100%); border: 1px solid rgba(255,255,255,0.15)">
      <q-card-section class="q-pa-lg">
        <div class="row items-center justify-between wrap gap-y-3">
          <div class="row items-center q-gutter-x-md">
            <q-avatar size="52px" color="white" text-color="indigo-9" class="shadow-soft">
              <q-icon name="edit_calendar" size="28px" />
            </q-avatar>
            <div>
              <div class="text-caption opacity-80 text-uppercase letter-spacing-1 text-weight-bold">
                Lezione in Corso — {{ currentHourLabel }}
              </div>
              <div class="text-h5 text-weight-bold">
                {{ activeLesson ? `${activeLesson.class_name} — ${activeLesson.subject_name}` : 'Seleziona una classe per firmare l\'appello' }}
              </div>
            </div>
          </div>

          <div class="row items-center q-gutter-x-sm">
            <q-btn
              color="positive"
              icon="draw"
              label="FIRMA ORA E REGISTRA PRESENZE"
              size="lg"
              unelevated
              class="rounded-lg text-weight-bolder shadow-md"
              @click="quickSignLesson"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Stats Cards -->
    <div class="row q-col-gutter-lg q-mb-xl">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card
          class="dashboard-card glass-card bg-indigo-600 text-white shadow-soft overflow-hidden cursor-pointer"
          @click="$router.push('/teacher/agenda')"
        >
          <q-card-section>
            <div class="text-caption opacity-80 text-uppercase letter-spacing-1">Prossima Lezione</div>
            <div class="text-h4 text-weight-bold q-mt-sm">{{ nextLesson?.class_name || nextLesson?.class_id || 'Nessuna' }}</div>
            <div class="text-caption q-mt-xs">{{ nextLesson?.subject_name || nextLesson?.subject_id || '-' }}</div>
          </q-card-section>
          <q-icon name="schedule" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card
          class="dashboard-card glass-card bg-orange-600 text-white shadow-soft overflow-hidden cursor-pointer"
          @click="$router.push('/teacher/attendance')"
        >
          <q-card-section>
            <div class="text-caption opacity-80 text-uppercase letter-spacing-1">Da Fare</div>
            <div class="text-h4 text-weight-bold q-mt-sm">{{ teacherStore.pendingJustifications }} Revisioni</div>
            <div class="text-caption q-mt-xs">Giustificazioni in sospeso</div>
          </q-card-section>
          <q-icon name="pending_actions" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card
          class="dashboard-card glass-card bg-emerald-600 text-white shadow-soft overflow-hidden cursor-pointer"
          @click="$router.push('/teacher/colloqui')"
        >
          <q-card-section>
            <div class="text-caption opacity-80 text-uppercase letter-spacing-1">Colloqui</div>
            <div class="text-h4 text-weight-bold q-mt-sm">{{ teacherStore.upcomingColloqui }} Prenotazioni</div>
            <div class="text-caption q-mt-xs">Controlla l'agenda</div>
          </q-card-section>
          <q-icon name="people" class="card-bg-icon" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card
          class="dashboard-card glass-card bg-violet-600 text-white shadow-soft overflow-hidden cursor-pointer"
          @click="$router.push('/teacher/communications')"
        >
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
            <div class="text-h6 q-mb-md">Azioni Rapide</div>
            <div class="row q-col-gutter-sm">
              <div class="col-6 col-sm-3">
                <q-btn
                  push color="primary" label="Registra Voti" icon="grade"
                  class="full-width" size="lg" to="/teacher/grades"
                />
              </div>
              <div class="col-6 col-sm-3">
                <q-btn
                  push color="secondary" label="Presenze" icon="fact_check"
                  class="full-width" size="lg" to="/teacher/attendance"
                />
              </div>
              <div class="col-6 col-sm-3">
                <q-btn
                  push color="accent" label="Firma Doc" icon="draw"
                  class="full-width" size="lg" to="/teacher/documents"
                />
              </div>
              <div class="col-6 col-sm-3">
                <q-btn
                  push color="info" label="Messaggi" icon="send"
                  class="full-width" size="lg" to="/teacher/communications"
                />
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Notifications -->
        <q-card>
          <q-card-section>
            <div class="text-h6">Notifiche & Attività</div>
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
                <q-btn flat round icon="close" size="sm" @click.stop="teacherStore.dismissNotification?.(note.id)" />
              </q-item-section>
            </q-item>
            <q-item v-if="!teacherStore.notifications || teacherStore.notifications.length === 0">
              <q-item-section class="text-center text-grey-5 q-py-lg">
                <q-icon name="notifications_none" size="32px" class="q-mb-sm" />
                <div>Nessuna notifica recente</div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Side Panel: Class List -->
      <div class="col-12 col-md-4">
        <q-card class="full-height">
          <q-card-section class="bg-grey-2">
            <div class="text-h6">Le Mie Classi</div>
          </q-card-section>
          <q-list separator>
            <q-item v-for="cls in classesStore.classes" :key="cls.id" clickable @click="classesStore.selectClass(cls.id)" :active="cls.id === classesStore.selectedClassId" active-class="bg-blue-1 text-primary">
              <q-item-section avatar>
                <q-avatar color="primary" text-color="white" size="sm">{{ cls.name?.charAt(0) || '?' }}</q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ cls.name || '' }}{{ cls.section || '' }}{{ cls.articolazione ? ' - ' + cls.articolazione : '' }}</q-item-label>
                <q-item-label caption>{{ cls.academic_year }}</q-item-label>
              </q-item-section>
              <q-item-section side v-if="cls.coordinator_id">
                <q-icon name="star" color="orange" title="Coordinatore" />
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
import { useRouter } from 'vue-router';
import { useTeacherStore } from 'src/stores/teacher';
import { useClassesStore } from 'src/stores/classes';

const router = useRouter();
const teacherStore = useTeacherStore();
const classesStore = useClassesStore();

const todayDate = computed(() => new Date().toLocaleDateString('it-IT', { day: 'numeric', month: 'long', year: 'numeric' }));

const nextLesson = computed(() => teacherStore.profile?.next_lesson || null);
const unreadMessagesCount = computed(() => teacherStore.notifications.filter(n => !n.read).length);

// Calculate current hour based on local time
const currentHourNumber = computed(() => {
  const hour = new Date().getHours();
  if (hour < 9) return 1;
  if (hour === 9) return 2;
  if (hour === 10) return 3;
  if (hour === 11) return 4;
  if (hour === 12) return 5;
  if (hour === 13) return 6;
  return 1;
});

const currentHourLabel = computed(() => `${currentHourNumber.value}ª Ora (${new Date().toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit' })})`);

const activeLesson = computed(() => {
  if (nextLesson.value) {
    return {
      class_id: nextLesson.value.class_id,
      class_name: nextLesson.value.class_name || `Classe ${nextLesson.value.class_id}`,
      subject_name: nextLesson.value.subject_name || 'Materia'
    };
  }
  if (classesStore.classes.length > 0) {
    const firstCls = classesStore.classes[0];
    return {
      class_id: firstCls.id,
      class_name: firstCls.name || `Classe ${firstCls.id}`,
      subject_name: 'Lezione In Corso'
    };
  }
  return null;
});

function quickSignLesson() {
  const classId = activeLesson.value?.class_id || classesStore.classes[0]?.id;
  if (classId) {
    router.push({
      path: '/teacher/attendance',
      query: { class_id: classId, hour: currentHourNumber.value }
    });
  } else {
    router.push('/teacher/attendance');
  }
}

onMounted(async () => {
  await Promise.all([
    teacherStore.fetchProfile(),
    teacherStore.fetchNotifications(),
    teacherStore.fetchPendingJustifications(),
    teacherStore.fetchUpcomingColloqui(),
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
