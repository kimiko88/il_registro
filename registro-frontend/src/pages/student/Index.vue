<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center q-mb-md justify-between">
      <div>
        <h1 class="text-h4 q-my-none">{{ $t('roleDashboards.studentPanel', { name: studentStore.profile?.first_name || '' }) }}</h1>
        <div class="text-subtitle1 text-grey-8">{{ studentStore.className }}</div>
      </div>
      <q-btn round flat icon="notifications" color="grey-8">
        <q-badge color="red" floating v-if="studentStore.notifications.length">{{ studentStore.notifications.length }}</q-badge>
        <q-menu>
             <q-list style="min-width: 300px">
                 <q-item-label header>{{ $t('dashboardPage.notifications') || 'Notifiche' }}</q-item-label>
                 <q-item v-for="n in studentStore.notifications" :key="n.id" clickable v-close-popup>
                     <q-item-section avatar><q-icon :name="n.icon || 'notifications'" :color="n.color || 'primary'" /></q-item-section>
                     <q-item-section>
                         <q-item-label>{{ n.title }}</q-item-label>
                         <q-item-label caption>{{ n.time || n.created_at }}</q-item-label>
                     </q-item-section>
                 </q-item>
                 <q-item v-if="!studentStore.notifications.length">
                     <q-item-section class="text-center text-grey">{{ $t('dashboardPage.noNotifications') }}</q-item-section>
                 </q-item>
             </q-list>
        </q-menu>
      </q-btn>
    </div>

    <!-- Quick Stats Cards -->
    <div class="row q-col-gutter-lg q-mb-xl">
      <!-- Average Grade -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card full-height">
            <q-card-section class="row items-center no-wrap">
                <div class="col">
                    <div class="text-h3 text-weight-bold text-outfit">{{ averageGrade }}</div>
                    <div class="text-caption text-slate-500 text-uppercase letter-spacing-1">{{ $t('roleDashboards.averageGrade') }}</div>
                </div>
                <div class="col-auto">
                    <div :class="Number(averageGrade) >= 6 ? 'bg-green-100' : 'bg-red-100'" class="q-pa-md rounded-xl">
                        <q-icon :name="Number(averageGrade) >= 6 ? 'trending_up' : 'trending_down'" :color="Number(averageGrade) >= 6 ? 'green-7' : 'red-7'" size="32px" />
                    </div>
                </div>
            </q-card-section>
            <q-linear-progress :value="Number(averageGrade)/10" :color="Number(averageGrade) >= 6 ? 'green' : 'red'" class="q-mt-none" size="4px" />
        </q-card>
      </div>

      <!-- Attendance -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card full-height">
            <q-card-section class="row items-center no-wrap">
                <div class="col">
                    <div class="text-h3 text-weight-bold text-outfit">{{ attendanceRate }}%</div>
                    <div class="text-caption text-slate-500 text-uppercase letter-spacing-1">{{ $t('roleDashboards.attendanceRate') }}</div>
                </div>
                <div class="col-auto">
                    <q-circular-progress
                      show-value
                      :value="attendanceRate"
                      size="56px"
                      :thickness="0.15"
                      color="indigo"
                      track-color="indigo-100"
                      class="q-ma-none text-indigo-700 text-weight-bold"
                    >
                        <div style="font-size: 10px">{{ attendanceRate }}%</div>
                    </q-circular-progress>
                </div>
            </q-card-section>
        </q-card>
      </div>

      <!-- PCTO Hours -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card full-height">
            <q-card-section class="row items-center no-wrap">
                <div class="col">
                    <div class="text-h3 text-weight-bold text-outfit">{{ pctoHours }}h</div>
                    <div class="text-caption text-slate-500 text-uppercase letter-spacing-1">{{ $t('roleDashboards.pctoHours') }}</div>
                </div>
                <div class="col-auto">
                    <div class="bg-orange-100 q-pa-md rounded-xl">
                        <q-icon name="work" color="orange-7" size="32px" />
                    </div>
                </div>
            </q-card-section>
            <q-linear-progress :value="pctoHours/150" color="orange" class="q-mt-none" size="4px" />
        </q-card>
      </div>

       <!-- Unread Messages -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card full-height cursor-pointer" @click="$router.push('/student/communications')">
            <q-card-section class="row items-center no-wrap">
                <div class="col">
                    <div class="text-h3 text-weight-bold text-outfit">{{ unreadMessages }}</div>
                    <div class="text-caption text-slate-500 text-uppercase letter-spacing-1">{{ $t('nav.communications') || 'Messaggi' }}</div>
                </div>
                <div class="col-auto">
                    <div class="bg-blue-100 q-pa-md rounded-xl">
                        <q-icon name="mail" color="blue-7" size="32px" />
                    </div>
                </div>
            </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="row q-col-gutter-lg">
      <!-- Recent Grades -->
      <div class="col-12 col-md-8">
        <q-card class="glass-card shadow-soft q-mb-lg overflow-hidden">
          <q-card-section class="row items-center justify-between q-pa-lg">
            <div class="text-h5 text-weight-bold text-outfit">{{ $t('gradesPage.recentGrades') || 'Ultimi Voti' }}</div>
            <q-btn flat :label="$t('common.viewAll') || 'Vedi Tutti'" color="primary" to="/student/grades" no-caps />
          </q-card-section>
          <q-separator color="white" style="opacity: 0.1" />
          
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
                 <q-item-section class="text-center text-grey q-py-lg">{{ $t('gradesPage.noGrades') || 'Nessun voto recente' }}</q-item-section>
             </q-item>
          </q-list>
        </q-card>
        
        <!-- Upcoming Events & Homework from Agenda -->
        <q-card class="glass-card shadow-soft overflow-hidden">
            <q-card-section class="row items-center justify-between q-pa-lg">
                <div>
                    <div class="text-h5 text-weight-bold text-outfit row items-center">
                        <q-icon name="assignment" color="primary" class="q-mr-sm" size="24px" />
                        <span>{{ $t('agendaPage.dueHomework') || 'Compiti & Verifiche in Arrivo' }}</span>
                    </div>
                    <div class="text-caption text-slate-500 q-mt-xs">{{ $t('agendaPage.organizeStudy') || 'Organizza le tue prossime scadenze di studio' }}</div>
                </div>
                <q-btn flat :label="$t('common.viewAll') || 'Vedi Tutti'" color="primary" to="/student/homework" no-caps />
            </q-card-section>
            <q-separator color="white" style="opacity: 0.1" />

            <q-list separator>
                <q-item v-for="event in upcomingEvents" :key="event.id" class="q-py-md cursor-pointer hover:bg-slate-50" @click="$router.push('/student/homework')">
                    <q-item-section avatar>
                        <q-avatar :color="event.color ? `${event.color}-1` : 'primary-1'" :text-color="event.color || 'primary'" :icon="event.icon || 'assignment'" size="42px" />
                    </q-item-section>
                    <q-item-section>
                        <q-item-label class="text-weight-bold text-slate-800 text-subtitle1">{{ event.title }}</q-item-label>
                        <q-item-label caption class="text-slate-500">
                            <span>{{ $t('agendaPage.due') || 'Scadenza' }}: {{ formatEventDate(event.start_date || event.date) }}</span>
                            <span v-if="event.start_time || event.time"> • Ore {{ event.start_time || event.time }}</span>
                            <span v-if="event.subject_name"> • {{ event.subject_name }}</span>
                        </q-item-label>
                    </q-item-section>
                    <q-item-section side v-if="event.type">
                        <q-chip :color="event.color || 'primary'" text-color="white" size="sm" class="text-weight-bold uppercase">{{ event.type }}</q-chip>
                    </q-item-section>
                </q-item>
                <q-item v-if="!upcomingEvents.length">
                    <q-item-section class="text-center text-slate-500 q-py-xl">
                        <q-icon name="task_alt" size="48px" color="positive" class="q-mb-sm opacity-80" />
                        <div class="text-subtitle1 text-weight-bold">{{ $t('agendaPage.noPendingHomework') || 'Nessun compito o verifica in arrivo' }}</div>
                        <div class="text-caption">{{ $t('agendaPage.allCaughtUp') || 'Sei in pari con tutte le attività!' }}</div>
                    </q-item-section>
                </q-item>
            </q-list>
        </q-card>
      </div>

      <!-- Quick Actions Sidebar -->
      <div class="col-12 col-md-4">
        <q-card class="bg-primary text-white q-mb-md">
            <q-card-section>
                <div class="text-h6">{{ $t('common.actions') || 'Accesso Rapido' }}</div>
            </q-card-section>
            <q-card-section class="q-pt-none">
                <div class="row q-col-gutter-sm">
                    <div class="col-12">
                        <q-btn unelevated color="white" text-color="primary" class="full-width" icon="fact_check" :label="$t('nav.myAttendance') || 'Presenze & Assenze'" to="/student/attendance" align="left" />
                    </div>
                     <div class="col-12">
                        <q-btn unelevated color="white" text-color="primary" class="full-width" icon="download" :label="$t('documentsPage.title') || 'Fascicolo Documentale & Atti'" to="/student/documents" align="left" />
                    </div>
                     <div class="col-12">
                        <q-btn unelevated color="white" text-color="primary" class="full-width" icon="edit_calendar" :label="$t('nav.myGrades') || 'I Miei Voti'" to="/student/grades" align="left" />
                    </div>
                </div>
            </q-card-section>
        </q-card>

        <q-card>
             <q-card-section>
                 <div class="text-subtitle1">{{ $t('dashboardPage.systemStatus') || 'Stato Sistema' }}</div>
                 <q-item>
                     <q-item-section avatar><q-icon name="wifi" color="green" /></q-item-section>
                     <q-item-section>
                         <q-item-label>{{ $t('dashboardPage.online') || 'Online' }}</q-item-label>
                         <q-item-label caption>{{ $t('dashboardPage.syncedNow') || 'Sincronizzato adesso' }}</q-item-label>
                     </q-item-section>
                 </q-item>
             </q-card-section>
         </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useStudentStore } from 'src/stores/student';
import { gradeService } from 'src/services/gradeService'
import { attendanceService } from 'src/services/attendanceService'
import { pctoService } from 'src/services/pctoService'
import { communicationService } from 'src/services/communicationService'
import dashboardService from 'src/services/dashboardService'
import adminService from 'src/services/adminService'
import api from 'src/services/api'

const { t, locale: currentLocale } = useI18n();
const studentStore = useStudentStore();

const averageGrade = ref('-')
const attendanceRate = ref(100)
const pctoHours = ref(0)
const unreadMessages = ref(0)
const subjects = ref([])

const recentGrades = ref([])
const upcomingEvents = ref([])

const getGradeColor = (val) => {
    if (val === 'A' || val === -1) return 'grey-6';
    const v = parseFloat(val);
    if (v >= 8) return 'green-6';
    if (v >= 6) return 'orange-6';
    return 'red-6';
}

const formatEventDate = (dateStr) => {
    if (!dateStr) return ''
    return new Date(dateStr).toLocaleDateString(currentLocale.value || 'it-IT', { day: '2-digit', month: 'short' })
}

const getSubjectName = (id) => {
    const s = subjects.value.find(s => s.id === id)
    return s ? s.name : id
}

onMounted(async () => {
    await studentStore.fetchProfile();
    studentStore.fetchNotifications();
    await fetchSubjects();
    fetchDashboardData();
});

const fetchSubjects = async () => {
    try {
        const schoolId = studentStore.profile?.school_id || studentStore.profile?.schoolId
        if (schoolId) {
            const res = await adminService.getSubjects(schoolId)
            subjects.value = res.data || []
        }
    } catch (e) {
        console.error('Error fetching subjects', e)
    }
}

const fetchDashboardData = async () => {
    try {
        // Fast aggregated stats from backend
        try {
            const stats = await dashboardService.getDashboardStats('student')
            if (stats) {
                if (stats.average_grade != null && stats.average_grade > 0) {
                    averageGrade.value = Number(stats.average_grade).toFixed(1)
                }
                if (stats.presence_rate != null || stats.attendance_rate != null) {
                    attendanceRate.value = Math.round(stats.presence_rate ?? stats.attendance_rate)
                }
            }
        } catch {
            // Non-blocking fallback
        }

        // Detailed Grades
        const gradesRes = await gradeService.getMyGrades()
        const allGrades = []
        if (gradesRes.data && gradesRes.data.semesters) {
            gradesRes.data.semesters.forEach(s => {
                if (s.grades) allGrades.push(...s.grades)
            })
        }
        const validGrades = allGrades.filter(g => g.grade_value > 0)
        if (validGrades.length > 0) {
            const sum = validGrades.reduce((acc, g) => acc + Number(g.grade_value), 0)
            averageGrade.value = (sum / validGrades.length).toFixed(1)
        } else if (averageGrade.value === '-') {
            averageGrade.value = '-'
        }
        
        allGrades.sort((a, b) => new Date(b.date) - new Date(a.date))
        recentGrades.value = allGrades.slice(0, 5).map(g => ({
            id: g.id,
            subject: getSubjectName(g.subject_id),
            value: g.grade_value === -1 ? 'A' : g.grade_value,
            date: new Date(g.date).toLocaleDateString(currentLocale.value || 'it-IT'),
            type: g.grade_type,
            description: g.description
        }))

        // Attendance
        const attRes = await attendanceService.getMyAttendance()
        if (attRes.data) {
             const records = Array.isArray(attRes.data) ? attRes.data : (attRes.data.records || [])
             const total = records.length
             const absences = records.filter(r => r.status === 'absent').length
             if (total > 0) {
                 attendanceRate.value = Math.round(((total - absences) / total) * 100)
             }
        }

        // PCTO
        const pctoRes = await pctoService.getMyProjects()
        if (pctoRes.data) {
             let totalHours = 0
             const projects = Array.isArray(pctoRes.data) ? pctoRes.data : []
             projects.forEach(p => {
                 totalHours += p.hours_done || 0
             })
             pctoHours.value = totalHours
        }

        // Communications
        const commsRes = await communicationService.getMessages()
        if (commsRes.data) {
            unreadMessages.value = commsRes.data.filter(m => !m.read && !m.archived).length
        }

        // Upcoming agenda events
        try {
            const agendaRes = await api.get('/agenda/events', {
                params: { from: new Date().toISOString().split('T')[0], limit: 5 }
            })
            upcomingEvents.value = Array.isArray(agendaRes.data)
                ? agendaRes.data
                : (agendaRes.data?.items || [])
        } catch (e) {
            console.warn('Could not fetch agenda events:', e)
            upcomingEvents.value = []
        }

    } catch (e) {
        console.error('Dashboard fetch error', e)
    }
}
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
</style>
