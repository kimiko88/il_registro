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
            <div class="text-h6 text-weight-bold text-dark">
              {{ isDashboardAdmin ? 'Attività Recenti' : 'Lezioni di Oggi' }}
            </div>
            <q-btn flat round dense icon="more_horiz" color="grey-7" />
          </q-card-section>
          
          <q-list class="q-px-sm" v-if="isDashboardAdmin">
            <q-item v-for="event in recentEvents" :key="event.id" class="q-mb-sm rounded-lg hover-bg-grey">
              <q-item-section avatar>
                <div class="text-center bg-grey-2 rounded-lg q-pa-sm" style="min-width: 50px">
                  <q-icon :name="getEventIcon(event.type)" :color="getEventColor(event.type)" size="sm" />
                </div>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ event.description }}</q-item-label>
                <q-item-label caption>{{ event.user_name }} • {{ event.school_name || 'Sistema' }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <div class="text-caption text-grey-6">{{ formatDate(event.created_at) }}</div>
              </q-item-section>
            </q-item>
            <q-item v-if="recentEvents.length === 0" class="text-center text-grey q-pa-md">
                Nessuna attività recente
            </q-item>
          </q-list>
          
          <q-list class="q-px-sm" v-else>
            <q-item v-for="entry in todaySchedule" :key="entry.id" class="q-mb-sm rounded-lg hover-bg-grey">
              <q-item-section avatar>
                <div class="text-center bg-grey-2 rounded-lg q-pa-sm" style="min-width: 50px">
                  <div class="text-weight-bold text-primary">{{ entry.hour_index }}ª Ora</div>
                </div>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ entry.subject_name }}</q-item-label>
                <q-item-label caption>
                  {{ entry.teacher_name }} <span v-if="entry.room">• Aula {{ entry.room }}</span>
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-chip size="sm" color="grey-3" text-color="grey-8">
                  Pianificata
                </q-chip>
              </q-item-section>
            </q-item>
            <q-item v-if="todaySchedule.length === 0" class="text-center text-grey q-pa-md">
              <q-item-section>
                Nessuna lezione pianificata per oggi
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Quick Actions / Notifications -->
      <div class="col-12 col-md-4">
        <q-card class="no-shadow bg-primary text-white q-mb-md" style="background: linear-gradient(135deg, #4F46E5 0%, #3B82F6 100%);">
          <q-card-section>
            <div class="text-subtitle2 text-blue-1 q-mb-xs">
              {{ latestAnnouncement ? latestAnnouncement.type.toUpperCase() : 'COMUNICAZIONE' }}
            </div>
            <div class="text-h6 text-weight-bold q-mb-sm">
              {{ latestAnnouncement ? latestAnnouncement.subject : 'Benvenuto nel Registro' }}
            </div>
            <div class="text-body2 text-blue-1 opacity-80">
              {{ latestAnnouncement ? latestAnnouncement.body : 'Le comunicazioni ufficiali e gli annunci per l\'anno scolastico corrente saranno mostrati in questa sezione.' }}
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
                  @click="handleActionClick(action)"
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
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useStudentStore } from 'src/stores/student'
import { useParentStore } from 'src/stores/parent'
import { useClassesStore } from 'src/stores/classes'
import adminService from 'src/services/adminService'
import dashboardService from 'src/services/dashboardService'
import { communicationService } from '@/services/communicationService'

const authStore = useAuthStore()
const { user, userRole } = storeToRefs(authStore)
const router = useRouter()
const $q = useQuasar()

const realStats = ref([])
const recentEvents = ref([])
const announcements = ref([])
const todaySchedule = ref([])

const isDashboardAdmin = computed(() => {
  return userRole.value === 'secretary' || userRole.value === 'admin' || userRole.value === 'superadmin'
})

const latestAnnouncement = computed(() => {
  if (announcements.value && announcements.value.length > 0) {
    return announcements.value[0]
  }
  return null
})

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
        if (data) {
            // Map backend stats to the format expected by the component
            if (userRole.value === 'admin' || userRole.value === 'superadmin') {
                realStats.value = [
                    { label: 'Totale Scuole', value: data.total_schools || '0', icon: 'school', color: 'indigo' },
                    { label: 'Utenti Attivi', value: data.total_users || '0', icon: 'people', color: 'cyan' },
                    { label: 'Eventi Oggi', value: data.active_users_24h || '0', icon: 'event', color: 'amber' },
                    { label: 'Report Pending', value: data.pending_documents_count || '0', icon: 'assignment', color: 'red' }
                ]
                recentEvents.value = data.recent_events || []
            } else if (userRole.value === 'secretary') {
                realStats.value = [
                    { label: 'Studenti', value: data.total_students || '0', icon: 'school', color: 'indigo' },
                    { label: 'Docenti', value: data.total_teachers || '0', icon: 'people', color: 'cyan' },
                    { label: 'Documenti', value: data.total_documents || '0', icon: 'description', color: 'amber' },
                    { label: 'Richieste', value: data.pending_documents_count || '0', icon: 'assignment', color: 'red' }
                ]
                recentEvents.value = data.recent_events || []
            } else if (userRole.value === 'teacher') {
                realStats.value = [
                    { label: 'Le Mie Classi', value: data.classes_count || '0', icon: 'class', color: 'indigo' },
                    { label: 'Studenti', value: data.students_count || '0', icon: 'school', color: 'cyan' },
                    { label: 'Lezioni Oggi', value: data.lessons_today_count || '0', icon: 'event', color: 'amber' },
                    { label: 'Voti da inserire', value: data.grades_pending_count || '0', icon: 'grade', color: 'red' }
                ]
            }
        }
    } catch (e) {
        console.error("Error fetching dashboard data", e)
    }
}

const getEventIcon = (type) => {
    const icons = { create: 'add_circle', update: 'edit', delete: 'delete', login: 'login' }
    return icons[type] || 'event'
}

const getEventColor = (type) => {
    const colors = { create: 'positive', update: 'info', delete: 'negative', login: 'primary' }
    return colors[type] || 'grey'
}

const formatDate = (dateString) => {
    const date = new Date(dateString)
    const diff = new Date() - date
    if (diff < 3600000) return `${Math.floor(diff / 60000)} min fa`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h fa`
    return date.toLocaleDateString('it-IT', { day: '2-digit', month: 'short' })
}

const fetchAnnouncements = async () => {
    try {
        const response = await communicationService.getMessages()
        announcements.value = response.data || []
    } catch (e) {
        console.error("Error fetching announcements:", e)
    }
}

const fetchTodaySchedule = async () => {
    if (isDashboardAdmin.value) return

    try {
        let classId = null

        if (userRole.value === 'student') {
            const studentStore = useStudentStore()
            await studentStore.fetchProfile()
            classId = studentStore.profile?.class_id
        } else if (userRole.value === 'parent') {
            const parentStore = useParentStore()
            await parentStore.fetchChildren()
            if (parentStore.children && parentStore.children.length > 0) {
                classId = parentStore.children[0].class_id
            }
        } else if (userRole.value === 'teacher') {
            const classesStore = useClassesStore()
            await classesStore.fetchAssignedClasses()
            if (classesStore.classes && classesStore.classes.length > 0) {
                classId = classesStore.classes[0].id
            }
        }

        if (classId) {
            const res = await adminService.getClassSchedule(classId)
            const allEntries = res.data || []

            const todayDay = new Date().getDay() // 0 = Sunday, 1 = Monday, ...
            const targetDay = todayDay === 0 ? 1 : todayDay // Fallback to Monday if Sunday

            const filtered = allEntries.filter(e => e.day_of_week === targetDay)
            filtered.sort((a, b) => a.hour_index - b.hour_index)

            todaySchedule.value = filtered
        }
    } catch (e) {
        console.error('Error fetching today schedule:', e)
    }
}

onMounted(() => {
    fetchDashboardData()
    fetchAnnouncements()
    fetchTodaySchedule()
})

const handleActionClick = (action) => {
    if (action.label === 'Nuovo Evento' || action.label === 'Invia Email') {
        if (userRole.value === 'teacher') {
            router.push('/teacher/communications')
        } else if (userRole.value === 'secretary') {
            router.push('/secretary/communications')
        } else if (userRole.value === 'student') {
            router.push('/student/communications')
        } else if (userRole.value === 'parent') {
            router.push('/parent/communications')
        } else {
            router.push('/admin/users')
        }
    } else if (action.label === 'Impostazioni') {
        if (userRole.value === 'admin' || userRole.value === 'superadmin') {
            router.push('/admin/settings')
        } else if (userRole.value === 'secretary') {
            router.push('/secretary/settings')
        } else {
            router.push(userRole.value === 'teacher' ? '/teacher' : `/${userRole.value}/profile`)
        }
    } else if (action.label === 'Stampa Voti') {
        if (userRole.value === 'teacher') {
            router.push('/teacher/grades')
        } else if (userRole.value === 'student') {
            router.push('/student/grades')
        } else if (userRole.value === 'parent') {
            router.push('/parent/grades')
        } else {
            $q.notify({
                type: 'info',
                message: 'Funzionalità disponibile per docenti, studenti e genitori.'
            })
        }
    }
}

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
