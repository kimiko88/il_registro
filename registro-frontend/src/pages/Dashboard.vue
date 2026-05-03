<template>
  <q-page class="q-pa-md q-pa-lg-xl">
    <!-- Hero Section -->
    <div class="row items-center q-mb-xl">
      <div class="col-12 col-md-8">
        <h1 class="text-h3 text-weight-bold text-outfit q-my-none text-gradient-premium">
          {{ greeting }}, {{ user?.first_name || 'Utente' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">
          Bentornato! Ecco il riepilogo delle attività scolastiche di oggi.
        </div>
      </div>
      <div class="col-12 col-md-4 text-right gt-sm">
        <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">Data Odierna</div>
        <div class="text-h6 text-outfit text-weight-bold text-slate-700">{{ new Date().toLocaleDateString('it-IT', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}</div>
      </div>
    </div>

    <!-- Quick Stats -->
    <div class="row q-col-gutter-lg q-mb-xl">
      <div class="col-12 col-sm-6 col-md-3" v-for="(stat, index) in stats" :key="index">
        <q-card class="glass-card stat-card full-height">
          <q-card-section class="row items-center no-wrap">
            <div :class="`bg-${stat.color}-100 text-${stat.color}-700 q-pa-md rounded-xl q-mr-md`">
              <q-icon :name="stat.icon" size="28px" />
            </div>
            <div>
              <div class="text-h5 text-weight-bold text-outfit">{{ stat.value }}</div>
              <div class="text-caption text-slate-500 text-uppercase letter-spacing-1" style="font-size: 10px">{{ stat.label }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="row q-col-gutter-lg">
      <!-- Recent Activity / Schedule -->
      <div class="col-12 col-md-8">
        <q-card class="no-shadow bordered-card full-height">
          <q-card-section class="row items-center justify-between">
            <div class="text-h6 text-weight-bold text-dark">Today's Schedule</div>
            <q-btn flat round dense icon="more_horiz" color="grey-7" />
          </q-card-section>
          
          <q-list class="q-px-sm">
            <q-item v-for="n in 3" :key="n" class="q-mb-sm rounded-lg hover-bg-grey">
              <q-item-section avatar>
                <div class="text-center bg-grey-2 rounded-lg q-pa-sm" style="min-width: 50px">
                  <div class="text-weight-bold text-primary">0{{ 8 + n }}:00</div>
                </div>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Mathematics - Class 3A</q-item-label>
                <q-item-label caption>Room 102 • Lecture Hall</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-chip size="sm" :color="n === 1 ? 'primary' : 'grey-3'" :text-color="n === 1 ? 'white' : 'grey-8'">
                  {{ n === 1 ? 'Ongoing' : 'Upcoming' }}
                </q-chip>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Quick Actions / Notifications -->
      <div class="col-12 col-md-4">
        <q-card class="no-shadow bg-primary text-white q-mb-md" style="background: linear-gradient(135deg, #4F46E5 0%, #3B82F6 100%);">
          <q-card-section>
            <div class="text-subtitle2 text-blue-1 q-mb-xs">ANNOUNCEMENT</div>
            <div class="text-h6 text-weight-bold q-mb-sm">School Meeting</div>
            <div class="text-body2 text-blue-1 opacity-80">
              There will be a staff meeting today at 2 PM in the main auditorium.
            </div>
          </q-card-section>
        </q-card>

        <q-card class="no-shadow bordered-card">
          <q-card-section>
            <div class="text-h6 text-weight-bold text-dark q-mb-md">Quick Actions</div>
            <div class="row q-col-gutter-sm">
              <div class="col-6" v-for="action in actions" :key="action.label">
                <q-btn 
                  outline 
                  class="full-width text-dark" 
                  style="border-color: #e2e8f0; border-radius: 12px; height: 80px"
                  no-caps
                >
                  <div class="column items-center">
                    <q-icon :name="action.icon" color="primary" size="sm" class="q-mb-xs" />
                    <div class="text-caption text-weight-medium">{{ action.label }}</div>
                  </div>
                </q-btn>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { ref, computed, onMounted } from 'vue'
import dashboardService from 'src/services/dashboardService'

const authStore = useAuthStore()
const { user, userName, userRole } = storeToRefs(authStore)

const realStats = ref([])

// Greeting based on time of day
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Buongiorno'
  if (hour < 18) return 'Buon pomeriggio'
  return 'Buonasera'
})

// Role-specific stats
const stats = computed(() => {
  if (realStats.value && realStats.value.length > 0) {
      return realStats.value
  }

  // Fallback to placeholders if no real data
  const roleStats = {
    admin: [
      { label: 'Totale Scuole', value: '12', icon: 'school', color: 'indigo' },
      { label: 'Utenti Attivi', value: '1,245', icon: 'people', color: 'cyan' },
      { label: 'Eventi Oggi', value: '4', icon: 'event', color: 'amber' },
      { label: 'Report Pending', value: '12', icon: 'assignment', color: 'red' }
    ],
    teacher: [
      { label: 'Le Mie Classi', value: '5', icon: 'class', color: 'indigo' },
      { label: 'Studenti', value: '120', icon: 'school', color: 'cyan' },
      { label: 'Lezioni Oggi', value: '4', icon: 'event', color: 'amber' },
      { label: 'Voti da inserire', value: '8', icon: 'grade', color: 'red' }
    ],
    student: [
      { label: 'Media Voti', value: '7.5', icon: 'grade', color: 'indigo' },
      { label: 'Presenze', value: '95%', icon: 'how_to_reg', color: 'cyan' },
      { label: 'Compiti', value: '3', icon: 'assignment', color: 'amber' },
      { label: 'Documenti', value: '12', icon: 'description', color: 'purple' }
    ],
    parent: [
      { label: 'I Miei Figli', value: '2', icon: 'family_restroom', color: 'indigo' },
      { label: 'Colloqui', value: '1', icon: 'event', color: 'cyan' },
      { label: 'Comunicazioni', value: '3', icon: 'email', color: 'amber' },
      { label: 'Documenti', value: '8', icon: 'description', color: 'purple' }
    ],
    secretary: [
      { label: 'Studenti', value: '450', icon: 'school', color: 'indigo' },
      { label: 'Docenti', value: '45', icon: 'people', color: 'cyan' },
      { label: 'Documenti', value: '24', icon: 'description', color: 'amber' },
      { label: 'Richieste', value: '7', icon: 'assignment', color: 'red' }
    ]
  }
  
  return roleStats[userRole.value] || roleStats.student
})

const fetchDashboardData = async () => {
    try {
        const data = await dashboardService.getDashboardStats(userRole.value)
        if (data && data.stats) {
            // Map backend stats to the format expected by the component
            // For admin, it might return { schools_count: X, users_count: Y, ... }
            if (userRole.value === 'admin' || userRole.value === 'superadmin') {
                realStats.value = [
                    { label: 'Totale Scuole', value: data.schools_count || '0', icon: 'school', color: 'indigo' },
                    { label: 'Utenti Attivi', value: data.users_count || '0', icon: 'people', color: 'cyan' },
                    { label: 'Eventi Oggi', value: '0', icon: 'event', color: 'amber' },
                    { label: 'Report Pending', value: '0', icon: 'assignment', color: 'red' }
                ]
            }
        }
    } catch (e) {
        console.error("Error fetching dashboard data", e)
    }
}

onMounted(() => {
    fetchDashboardData()
})

const actions = [
  { label: 'Nuovo Evento', icon: 'add_circle' },
  { label: 'Invia Email', icon: 'mail' },
  { label: 'Stampa Voti', icon: 'print' },
  { label: 'Impostazioni', icon: 'settings' }
]
</script>

<style scoped>
.letter-spacing-1 {
    letter-spacing: 1px;
}

.stat-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.bordered-card {
  border: 1px solid rgba(0,0,0,0.05);
  border-radius: 20px;
}

.hover-bg-grey:hover {
  background-color: rgba(79, 70, 229, 0.05);
}
</style>
