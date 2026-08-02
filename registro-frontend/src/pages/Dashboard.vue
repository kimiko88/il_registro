<template>
  <q-page class="q-pa-md q-pa-lg-xl">
    <!-- Hero Section -->
    <div class="row items-center q-mb-xl">
      <div class="col-12 col-md-8">
        <h1 class="text-h4 text-sm-h3 text-weight-bold text-outfit q-my-none text-primary">
          {{ greeting }}, {{ user?.first_name || 'Utente' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">
          Bentornato! Ecco il riepilogo delle attività scolastiche di oggi.
        </div>
      </div>
      <div class="col-12 col-md-4 text-right gt-sm">
        <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">Data Odierna</div>
        <div class="text-h6 text-outfit text-weight-bold text-slate-700">{{ today }}</div>
      </div>
    </div>

    <!-- Quick Stats -->
    <div class="row q-col-gutter-lg q-mb-xl" :aria-busy="loadingData" aria-live="polite">
      <div class="col-12 col-sm-6 col-md-3" v-for="(stat, index) in stats" :key="index">
        <q-card :aria-label="`${stat.label}: ${stat.value}`" class="glass-card stat-card full-height">
          <q-card-section class="row items-center no-wrap" v-if="!loadingData">
            <div :class="`bg-${stat.color}-100 text-${stat.color}-700 q-pa-md rounded-xl q-mr-md`">
              <q-icon :name="stat.icon" size="28px" />
            </div>
            <div>
              <div class="text-h5 text-weight-bold text-outfit">{{ stat.value }}</div>
              <div class="text-caption text-slate-500 text-uppercase letter-spacing-1" style="font-size: 12px">{{ stat.label }}</div>
            </div>
          </q-card-section>
          <q-card-section class="row items-center no-wrap" v-else>
            <q-skeleton type="QAvatar" size="48px" class="q-mr-md" />
            <div class="col">
              <q-skeleton type="text" width="60%" />
              <q-skeleton type="text" width="40%" />
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
            <q-btn flat round dense icon="more_horiz" color="grey-7" aria-label="Opzioni e filtro attività" />
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
                <q-chip size="sm" :color="getLessonStatusColor(entry)" :text-color="getLessonStatusTextColor(entry)">
                  {{ getLessonStatus(entry) }}
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
        <q-card class="no-shadow glass-card q-mb-md" style="border-left: 4px solid var(--q-primary);">
          <q-card-section>
            <q-skeleton v-if="loadingAnnouncements" type="text" :lines="3" />
            <template v-else>
              <div class="text-subtitle2 text-primary q-mb-xs">
                {{ latestAnnouncement ? latestAnnouncement.type.toUpperCase() : 'COMUNICAZIONE' }}
              </div>
              <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">
                {{ latestAnnouncement ? latestAnnouncement.subject : 'Benvenuto nel Registro' }}
              </div>
              <div class="text-body2 text-slate-600 opacity-80">
                {{ latestAnnouncement ? latestAnnouncement.body : 'Le comunicazioni ufficiali e gli annunci per l\'anno scolastico corrente saranno mostrati in questa sezione.' }}
              </div>
            </template>
          </q-card-section>
        </q-card>

        <q-card class="no-shadow bordered-card">
          <q-card-section>
            <div class="text-h6 text-weight-bold text-dark q-mb-md">Azioni Rapide</div>
            <div class="row q-col-gutter-sm">
              <div class="col-6" v-for="action in actions" :key="action.key">
                <q-btn 
                  outline 
                  class="full-width text-dark" 
                  style="border-color: #e2e8f0; border-radius: 12px; height: 80px"
                  no-caps
                  :loading="navigatingAction === action.key"
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
import api from '@/services/api'

const router = useRouter()
const authStore = useAuthStore()
const { user, userRole } = storeToRefs(authStore)
const $q = useQuasar()

const realStats = ref([])
const recentEvents = ref([])
const announcements = ref([])
const todaySchedule = ref([])
const loadingData = ref(false)
const loadingAnnouncements = ref(false)
const navigatingAction = ref(null)

const currentRole = computed(() => userRole.value || user.value?.role || authStore.userRole || authStore.user?.role || 'student')

const today = computed(() =>
  new Date().toLocaleDateString('it-IT', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
)

const isDashboardAdmin = computed(() => {
  return currentRole.value === 'secretary' || currentRole.value === 'admin' || currentRole.value === 'superadmin'
})

const latestAnnouncement = computed(() => {
  if (announcements.value && announcements.value.length > 0) {
    return announcements.value[0]
  }
  return null
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Buongiorno'
  if (hour < 18) return 'Buon pomeriggio'
  return 'Buonasera'
})

const stats = computed(() => {
  if (realStats.value && realStats.value.length > 0) {
    return realStats.value
  }
  // Nessun dato disponibile — mostra zeri invece di valori inventati
  return [
    { label: 'Dati', value: '-', icon: 'info', color: 'grey' }
  ]
})

const getLessonStatus = (entry) => {
  const now = new Date()
  const currentMinutes = now.getHours() * 60 + now.getMinutes()
  const hourIndex = Math.max(1, entry?.hour_index || 1)
  // Assumiamo che la 1ª ora scolastica inizi alle 8:00
  const startMinutes = (8 * 60) + ((hourIndex - 1) * 60)
  const endMinutes = startMinutes + 60
  if (currentMinutes >= endMinutes) return 'Completata'
  if (currentMinutes >= startMinutes) return 'In corso'
  return 'Pianificata'
}

const getLessonStatusColor = (entry) => {
  const status = getLessonStatus(entry)
  if (status === 'Completata') return 'grey-3'
  if (status === 'In corso') return 'positive'
  return 'blue-1'
}

const getLessonStatusTextColor = (entry) => {
  const status = getLessonStatus(entry)
  if (status === 'Completata') return 'grey-7'
  if (status === 'In corso') return 'white'
  return 'primary'
}

const fetchDashboardData = async () => {
    loadingData.value = true
    try {
        const role = currentRole.value
        const data = await dashboardService.getDashboardStats(role)
        if (data) {
            if (role === 'admin' || role === 'superadmin') {
                realStats.value = [
                    { label: 'Totale Scuole', value: data.total_schools ?? '0', icon: 'school', color: 'indigo' },
                    { label: 'Utenti Attivi', value: data.total_users ?? '0', icon: 'people', color: 'cyan' },
                    { label: 'Attivi 24h', value: data.active_users_24h ?? '0', icon: 'event', color: 'amber' },
                    { label: 'Doc. Pending', value: data.pending_documents_count ?? '0', icon: 'assignment', color: 'red' }
                ]
                recentEvents.value = data.recent_events || []
            } else if (role === 'secretary') {
                realStats.value = [
                    { label: 'Studenti', value: data.total_students ?? '0', icon: 'school', color: 'indigo' },
                    { label: 'Docenti', value: data.total_teachers ?? '0', icon: 'people', color: 'cyan' },
                    { label: 'Documenti', value: data.total_documents ?? '0', icon: 'description', color: 'amber' },
                    { label: 'Richieste', value: data.pending_documents_count ?? '0', icon: 'assignment', color: 'red' }
                ]
                recentEvents.value = data.recent_events || []
            } else if (role === 'teacher') {
                realStats.value = [
                    { label: 'Le Mie Classi', value: data.classes_count ?? '0', icon: 'class', color: 'indigo' },
                    { label: 'Studenti', value: data.students_count ?? '0', icon: 'school', color: 'cyan' },
                    { label: 'Lezioni Oggi', value: data.lessons_today_count ?? '0', icon: 'event', color: 'amber' },
                    { label: 'Voti da inserire', value: data.grades_pending_count ?? '0', icon: 'grade', color: 'red' }
                ]
            } else if (role === 'student') {
                realStats.value = [
                    { label: 'Media Voti', value: data.average_grade ?? '-', icon: 'grade', color: 'indigo' },
                    { label: 'Presenze', value: data.attendance_rate != null ? data.attendance_rate + '%' : '-', icon: 'how_to_reg', color: 'cyan' },
                    { label: 'Compiti', value: data.homework_count ?? '0', icon: 'assignment', color: 'amber' },
                    { label: 'Documenti', value: data.documents_count ?? '0', icon: 'description', color: 'purple' }
                ]
            } else if (role === 'parent') {
                realStats.value = [
                    { label: 'I Miei Figli', value: data.children_count ?? '0', icon: 'family_restroom', color: 'indigo' },
                    { label: 'Colloqui', value: data.upcoming_colloqui ?? '0', icon: 'event', color: 'cyan' },
                    { label: 'Comunicazioni', value: data.unread_communications ?? '0', icon: 'email', color: 'amber' },
                    { label: 'Documenti', value: data.documents_count ?? '0', icon: 'description', color: 'purple' }
                ]
            }
        }
    } catch (e) {
        console.error('Error fetching dashboard data', e)
        realStats.value = []
    } finally {
        loadingData.value = false
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
    loadingAnnouncements.value = true
    try {
        const response = await communicationService.getMessages()
        announcements.value = response.data || []
    } catch (e) {
        console.error('Error fetching announcements:', e)
    } finally {
        loadingAnnouncements.value = false
    }
}

const fetchTodaySchedule = async () => {
    if (isDashboardAdmin.value) return

    try {
        const todayDay = new Date().getDay()
        if (todayDay === 0) {
            todaySchedule.value = []
            return
        }

        if (userRole.value === 'teacher') {
            const res = await api.get('/timetables/my-schedule')
            const allEntries = res.data || []
            const filtered = allEntries.filter(e => e.day_of_week === todayDay)
            filtered.sort((a, b) => a.hour_index - b.hour_index)
            todaySchedule.value = filtered
            return
        }

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
        }

        if (classId) {
            const res = await adminService.getClassSchedule(classId)
            const allEntries = res.data || []

            const filtered = allEntries.filter(e => e.day_of_week === todayDay)
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

const handleActionClick = async (action) => {
    navigatingAction.value = action.key
    try {
        if (action.key === 'communications' || action.key === 'email') {
            if (userRole.value === 'teacher') {
                await router.push('/teacher/communications')
            } else if (userRole.value === 'secretary') {
                await router.push('/secretary/communications')
            } else if (userRole.value === 'student') {
                await router.push('/student/communications')
            } else if (userRole.value === 'parent') {
                await router.push('/parent/communications')
            } else {
                await router.push('/admin/users')
            }
        } else if (action.key === 'settings') {
            if (userRole.value === 'admin' || userRole.value === 'superadmin') {
                await router.push('/admin/settings')
            } else if (userRole.value === 'secretary') {
                await router.push('/secretary/settings')
            } else {
                await router.push(userRole.value === 'teacher' ? '/teacher' : `/${userRole.value}/profile`)
            }
        } else if (action.key === 'grades') {
            if (userRole.value === 'teacher') {
                await router.push('/teacher/grades')
            } else if (userRole.value === 'student') {
                await router.push('/student/grades')
            } else if (userRole.value === 'parent') {
                await router.push('/parent/grades')
            } else {
                $q.notify({
                    type: 'info',
                    message: 'Funzionalità disponibile per docenti, studenti e genitori.'
                })
            }
        }
    } finally {
        navigatingAction.value = null
    }
}

const allActions = [
  { label: 'Comunicazioni', icon: 'campaign', key: 'communications', roles: ['teacher', 'secretary', 'student', 'parent', 'admin', 'superadmin'] },
  { label: 'Invia Email', icon: 'mail', key: 'email', roles: ['teacher', 'secretary', 'admin', 'superadmin'] },
  { label: 'Stampa Voti', icon: 'print', key: 'grades', roles: ['teacher', 'student', 'parent'] },
  { label: 'Impostazioni', icon: 'settings', key: 'settings', roles: ['teacher', 'secretary', 'student', 'parent', 'admin', 'superadmin'] }
]

const actions = computed(() => {
  return allActions.filter(a => !a.roles || a.roles.includes(currentRole.value))
})
</script>

<style scoped>
.letter-spacing-1 {
    letter-spacing: 1px;
}

.bordered-card {
  border: 1px solid rgba(0,0,0,0.05);
  border-radius: 12px;
}

.hover-bg-grey:hover {
  background-color: rgba(79, 70, 229, 0.05);
}
</style>
