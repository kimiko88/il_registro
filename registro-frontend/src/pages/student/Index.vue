<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center q-mb-md justify-between">
      <div>
        <h1 class="text-h4 q-my-none">Bentornato, {{ studentStore.profile?.firstName }}</h1>
        <div class="text-subtitle1 text-grey-8">{{ studentStore.className || 'Classe 5A' }}</div>
      </div>
      <q-btn round flat icon="notifications" color="grey-8">
        <q-badge color="red" floating v-if="studentStore.notifications.length">{{ studentStore.notifications.length }}</q-badge>
        <q-menu>
             <q-list style="min-width: 300px">
                 <q-item-label header>Notifiche</q-item-label>
                 <q-item v-for="n in studentStore.notifications" :key="n.id" clickable v-close-popup>
                     <q-item-section avatar><q-icon :name="n.icon" :color="n.color" /></q-item-section>
                     <q-item-section>
                         <q-item-label>{{ n.title }}</q-item-label>
                         <q-item-label caption>{{ n.time }}</q-item-label>
                     </q-item-section>
                 </q-item>
                 <q-item v-if="!studentStore.notifications.length">
                     <q-item-section class="text-center text-grey">Nessuna nuova notifica</q-item-section>
                 </q-item>
             </q-list>
        </q-menu>
      </q-btn>
    </div>

    <!-- Quick Stats Cards with Trends -->
    <div class="row q-col-gutter-md q-mb-lg">
      <!-- Average Grade -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-2">
            <q-card-section class="row items-center no-wrap">
                <div class="col">
                    <div class="text-h6 text-weight-bold">{{ averageGrade }}</div>
                    <div class="text-caption text-grey">Media Voti</div>
                </div>
                <div class="col-auto">
                    <q-icon name="trending_up" color="green" size="md" v-if="Number(averageGrade) >= 6" />
                    <q-icon name="trending_down" color="red" size="md" v-else />
                </div>
            </q-card-section>
            <q-linear-progress :value="Number(averageGrade)/10" color="primary" class="q-mt-none" />
        </q-card>
      </div>

      <!-- Attendance -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-2">
            <q-card-section class="row items-center no-wrap">
                <div class="col">
                    <div class="text-h6 text-weight-bold">{{ attendanceRate }}%</div>
                    <div class="text-caption text-grey">Presenze</div>
                </div>
                <div class="col-auto">
                    <q-circular-progress
                      show-value
                      :value="attendanceRate"
                      size="40px"
                      :thickness="0.2"
                      color="green"
                      track-color="grey-3"
                      class="q-ma-none"
                    >
                        <q-icon name="check" size="12px" />
                    </q-circular-progress>
                </div>
            </q-card-section>
        </q-card>
      </div>

      <!-- PCTO Hours -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-2">
            <q-card-section class="row items-center no-wrap">
                <div class="col">
                    <div class="text-h6 text-weight-bold">{{ pctoHours }}h</div>
                    <div class="text-caption text-grey">Ore PCTO</div>
                </div>
                <div class="col-auto">
                    <q-icon name="work" color="orange" size="md" />
                </div>
            </q-card-section>
            <q-linear-progress :value="pctoHours/100" color="orange" class="q-mt-none" />
        </q-card>
      </div>

       <!-- Unread Messages -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="bg-white text-dark shadow-2 cursor-pointer" @click="$router.push('/student/communications')">
            <q-card-section class="row items-center no-wrap">
                <div class="col">
                    <div class="text-h6 text-weight-bold">{{ unreadMessages }}</div>
                    <div class="text-caption text-grey">Messaggi Nuovi</div>
                </div>
                <div class="col-auto">
                    <q-icon name="mail" color="blue" size="md" />
                </div>
            </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="row q-col-gutter-lg">
      <!-- Recent Activity / Grades -->
      <div class="col-12 col-md-8">
        <q-card class="q-mb-md">
          <q-card-section class="row items-center justify-between">
            <div class="text-h6">Ultimi Voti</div>
            <q-btn flat label="Vedi Tutti" color="primary" to="/student/grades" size="sm" />
          </q-card-section>
          <q-separator />
          
          <q-list separator>
             <q-item v-for="grade in recentGrades" :key="grade.id">
                <q-item-section avatar>
                     <q-avatar size="md" :color="getGradeColor(grade.value)" text-color="white">{{ grade.value }}</q-avatar>
                </q-item-section>
                <q-item-section>
                    <q-item-label class="text-weight-bold">{{ grade.subject }}</q-item-label>
                    <q-item-label caption>{{ grade.date }} - {{ grade.type }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                    <div class="text-caption">{{ grade.description }}</div>
                </q-item-section>
             </q-item>
             <q-item v-if="!recentGrades.length">
                 <q-item-section class="text-center text-grey q-py-lg">Nessun voto recente</q-item-section>
             </q-item>
          </q-list>
        </q-card>
        
        <!-- Agenda/Upcoming -->
        <q-card>
            <q-card-section class="text-h6">In Arrivo</q-card-section>
            <q-list>
                <q-item v-for="event in upcomingEvents" :key="event.id">
                    <q-item-section avatar>
                        <q-icon :name="event.icon" :color="event.color" />
                    </q-item-section>
                    <q-item-section>
                        <q-item-label>{{ event.title }}</q-item-label>
                        <q-item-label caption>{{ event.date }} {{ event.time }}</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <q-chip :color="event.color" text-color="white" size="sm">{{ event.tag }}</q-chip>
                    </q-item-section>
                </q-item>
            </q-list>
        </q-card>
      </div>

      <!-- Quick Actions Sidebar -->
      <div class="col-12 col-md-4">
        <q-card class="bg-primary text-white q-mb-md">
            <q-card-section>
                <div class="text-h6">Accesso Rapido</div>
            </q-card-section>
            <q-card-section class="q-pt-none">
                <div class="row q-col-gutter-sm">
                    <div class="col-12">
                        <q-btn unelevated color="white" text-color="primary" class="full-width" icon="fact_check" label="Giustifica Assenza" to="/student/attendance" align="left" />
                    </div>
                     <div class="col-12">
                        <q-btn unelevated color="white" text-color="primary" class="full-width" icon="download" label="Scarica Pagella" to="/student/documents" align="left" />
                    </div>
                     <div class="col-12">
                        <q-btn unelevated color="white" text-color="primary" class="full-width" icon="edit_calendar" label="Registro Elettronico" to="/student/grades" align="left" />
                    </div>
                </div>
            </q-card-section>
        </q-card>

        <!-- Profile/System Status -->
         <q-card>
             <q-card-section>
                 <div class="text-subtitle1">Stato Sistema</div>
                 <q-item>
                     <q-item-section avatar><q-icon name="wifi" color="green" /></q-item-section>
                     <q-item-section>
                         <q-item-label>Online</q-item-label>
                         <q-item-label caption>Sincronizzato adesso</q-item-label>
                     </q-item-section>
                 </q-item>
             </q-card-section>
         </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted, computed, ref } from 'vue';
import { useStudentStore } from 'src/stores/student';
// Assuming useGradesStore exists or we fetch via student store
// For now, mocking specific data here or extending student store

const studentStore = useStudentStore();

// Computed/Mock Data for UI demonstration
const averageGrade = ref('7.8')
const attendanceRate = ref(92)
const pctoHours = ref(30)
const unreadMessages = ref(2)

const recentGrades = ref([
    { id: 1, subject: 'Matematica', value: '8.5', date: '20/01', type: 'Scritto', description: 'Equazioni 2° grado' },
    { id: 2, subject: 'Storia', value: '7', date: '18/01', type: 'Orale', description: 'Rivoluzione Francese' },
    { id: 3, subject: 'Inglese', value: '9', date: '15/01', type: 'Listening', description: 'Comprehension Test' }
])

const upcomingEvents = ref([
    { id: 1, title: 'Verifica Fisica', date: 'Domani', time: '09:00', icon: 'quiz', color: 'red', tag: 'Verifica' },
    { id: 2, title: 'Scadenza Progetto Info', date: 'Lun 25', time: '23:59', icon: 'timer', color: 'orange', tag: 'Scadenza' },
])

const getGradeColor = (val) => {
    const v = parseFloat(val);
    if (v >= 8) return 'green';
    if (v >= 6) return 'orange';
    return 'red';
}

onMounted(() => {
    studentStore.fetchProfile();
    studentStore.fetchNotifications();
    // In real implementation:
    // gradesStore.fetchRecentGrades()
    // attendanceStore.fetchStats()
});
</script>

<style scoped>
.q-card {
    transition: transform 0.2s;
}
.q-card:hover { 
    /* slightly lift card on hover for desktop feeling */
    transform: translateY(-2px);
}
</style>
