<template>
  <q-page class="q-pa-md q-pa-lg-xl">
    <!-- Hero Section -->
    <div class="row items-center q-mb-xl">
      <div class="col-12 col-md-8">
        <h1 class="text-h4 text-sm-h3 text-weight-bold text-outfit q-my-none text-primary">
          {{ greeting }}, {{ user?.first_name || 'Utente' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">
          {{ $t('dashboardPage.welcomeSub') }}
        </div>
      </div>
      <div class="col-12 col-md-4 text-right gt-sm">
        <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">{{ $t('dashboardPage.todayDate') }}</div>
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
              {{ isDashboardAdmin ? $t('dashboardPage.recentActivity') : $t('dashboardPage.todayLessons') }}
            </div>
            
            <!-- 3 Dots Options Menu -->
            <q-btn flat round dense icon="more_horiz" color="grey-7" :aria-label="t('dashboardPage.optionsMenu')">
              <q-menu auto-close>
                <!-- Menu per SuperAdmin / Admin -->
                <q-list style="min-width: 240px" v-if="currentRole === 'admin' || currentRole === 'superadmin'">
                  <q-item clickable @click="router.push('/admin/audit-logs')">
                    <q-item-section avatar><q-icon name="fact_check" color="primary" /></q-item-section>
                    <q-item-section>{{ $t('dashboardPage.auditLogs') }}</q-item-section>
                  </q-item>
                  <q-item clickable @click="router.push('/admin/users')">
                    <q-item-section avatar><q-icon name="people" color="secondary" /></q-item-section>
                    <q-item-section>{{ $t('dashboardPage.userManagement') }}</q-item-section>
                  </q-item>
                  <q-item clickable @click="router.push('/admin/schools')">
                    <q-item-section avatar><q-icon name="school" color="indigo" /></q-item-section>
                    <q-item-section>{{ $t('dashboardPage.schoolManagement') }}</q-item-section>
                  </q-item>
                  <q-separator />
                  <q-item clickable @click="fetchDashboardData">
                    <q-item-section avatar><q-icon name="refresh" color="grey-7" /></q-item-section>
                    <q-item-section>{{ $t('dashboardPage.refreshActivity') }}</q-item-section>
                  </q-item>
                </q-list>

                <!-- Menu per Segreteria -->
                <q-list style="min-width: 240px" v-else-if="currentRole === 'secretary'">
                  <q-item clickable @click="router.push('/secretary/audit-log')">
                    <q-item-section avatar><q-icon name="fact_check" color="primary" /></q-item-section>
                    <q-item-section>{{ $t('dashboardPage.auditLogs') || 'Audit Log' }}</q-item-section>
                  </q-item>
                  <q-item clickable @click="router.push('/secretary/users')">
                    <q-item-section avatar><q-icon name="people" color="secondary" /></q-item-section>
                    <q-item-section>{{ $t('dashboardPage.userManagement') || 'Anagrafica Utenti' }}</q-item-section>
                  </q-item>
                  <q-item clickable @click="router.push('/secretary/documents')">
                    <q-item-section avatar><q-icon name="folder" color="amber-9" /></q-item-section>
                    <q-item-section>{{ $t('documentsPage.title') || 'Gestione Documenti & Atti' }}</q-item-section>
                  </q-item>
                  <q-separator />
                  <q-item clickable @click="fetchDashboardData">
                    <q-item-section avatar><q-icon name="refresh" color="grey-7" /></q-item-section>
                    <q-item-section>{{ $t('dashboardPage.refreshActivity') || 'Aggiorna Attività' }}</q-item-section>
                  </q-item>
                </q-list>

                <!-- Menu per Docente (Teacher) -->
                <q-list style="min-width: 240px" v-else-if="currentRole === 'teacher'">
                  <q-item clickable @click="openDraftModal">
                    <q-item-section avatar><q-icon name="edit_note" color="primary" /></q-item-section>
                    <q-item-section>{{ $t('udaPage.createTitle') || 'Pianifica Bozza Lezione' }}</q-item-section>
                  </q-item>
                  <q-item clickable @click="openDraftsList">
                    <q-item-section avatar><q-icon name="collections_bookmark" color="secondary" /></q-item-section>
                    <q-item-section>{{ $t('udaPage.title') || 'Bozze Salvate' }} ({{ lessonDrafts.length }})</q-item-section>
                  </q-item>
                  <q-separator />
                  <q-item clickable @click="router.push('/teacher/lessons')">
                    <q-item-section avatar><q-icon name="menu_book" color="grey-7" /></q-item-section>
                    <q-item-section>{{ $t('timetablePage.title') || 'Registro di Classe Completo' }}</q-item-section>
                  </q-item>
                </q-list>

                <!-- Menu Generico (Student / Parent) -->
                <q-list style="min-width: 240px" v-else>
                  <q-item clickable @click="fetchDashboardData">
                    <q-item-section avatar><q-icon name="refresh" color="primary" /></q-item-section>
                    <q-item-section>{{ $t('dashboardPage.refreshActivity') || 'Aggiorna Attività' }}</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
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
                {{ $t('communicationsPage.searchPlaceholder') || 'Nessuna attività recente' }}
            </q-item>
          </q-list>
          
          <q-list class="q-px-sm" v-else>
            <q-item v-for="entry in displaySchedule" :key="entry.id" class="q-mb-sm rounded-lg hover-bg-grey">
              <q-item-section avatar>
                <div class="text-center bg-grey-2 rounded-lg q-pa-sm" style="min-width: 50px">
                  <div class="text-weight-bold text-primary">{{ entry.hour_index }}ª {{ $t('timetablePage.hour') || 'Ora' }}</div>
                </div>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ entry.subject_name }}</q-item-label>
                <q-item-label caption>
                  {{ entry.teacher_name }} <span v-if="entry.room">• {{ $t('timetablePage.classroom') || 'Aula' }} {{ entry.room }}</span>
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-chip size="sm" :color="getLessonStatusColor(entry)" :text-color="getLessonStatusTextColor(entry)">
                  {{ getLessonStatus(entry) }}
                </q-chip>
              </q-item-section>
            </q-item>
            
            <q-item v-if="displaySchedule.length === 0" class="text-center text-grey q-pa-md">
              <q-item-section>
                <div>{{ $t('timetablePage.freeSlot') || 'Nessuna lezione pianificata per oggi' }}</div>
                <q-btn flat color="primary" icon="add" :label="$t('udaPage.newUda') || 'Crea Bozza Lezione per Oggi'" class="q-mt-sm" @click="openDraftModal" />
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
                {{ latestAnnouncement ? latestAnnouncement.type.toUpperCase() : ($t('communicationsPage.title') || 'COMUNICAZIONE') }}
              </div>
              <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">
                {{ latestAnnouncement ? latestAnnouncement.subject : ($t('login.welcomeBack') || 'Benvenuto nel Registro') }}
              </div>
              <div class="text-body2 text-slate-600 opacity-80">
                {{ latestAnnouncement ? latestAnnouncement.body : ($t('communicationsPage.subtitle') || 'Le comunicazioni ufficiali e gli annunci saranno mostrati in questa sezione.') }}
              </div>
            </template>
          </q-card-section>
        </q-card>

        <q-card class="no-shadow bordered-card">
          <q-card-section>
            <div class="text-h6 text-weight-bold text-dark q-mb-md">{{ $t('common.actions') || 'Azioni Rapide' }}</div>
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

    <!-- Dialog Pianifica Bozza Lezione -->
    <q-dialog v-model="showDraftDialog" persistent>
      <q-card style="min-width: 500px; max-width: 650px" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold">
            <q-icon name="edit_note" class="q-mr-xs" /> {{ $t('udaPage.createTitle') || 'Pianifica Bozza Lezione' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-sm">
          <q-input v-model="draftForm.date" type="date" :label="$t('gradesPage.date')" outlined dense />
          <q-select
            v-model="draftForm.class_id"
            :options="classesStore.classes"
            option-label="name"
            option-value="id"
            emit-value map-options
            :label="$t('udaPage.classLabel')"
            :loading="classesStore.loading"
            outlined dense
          />
          <q-select
            v-model="draftForm.hour"
            :options="[1, 2, 3, 4, 5, 6]"
            :label="$t('timetablePage.hour')"
            outlined dense
          />
          <q-input v-model="draftForm.subject" :label="$t('udaPage.subjectLabel')" outlined dense />
          <q-input v-model="draftForm.topic" type="textarea" rows="3" :label="$t('udaPage.descriptionLabel')" outlined dense />
          <q-input v-model="draftForm.homework" type="textarea" rows="2" :label="$t('didacticsPage.resourceCategory')" outlined dense />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn color="primary" icon="save" :label="$t('common.save')" unelevated @click="saveLessonDraft" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Elenco Bozze Salvate -->
    <q-dialog v-model="showDraftsListDialog">
      <q-card style="min-width: 600px" class="rounded-xl">
        <q-card-section class="bg-secondary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold">
            <q-icon name="collections_bookmark" class="q-mr-xs" /> {{ t('dashboardPage.draftsListTitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <q-list separator v-if="lessonDrafts.length > 0">
            <q-item v-for="(draft, idx) in lessonDrafts" :key="idx" class="q-py-md">
              <q-item-section avatar>
                <q-avatar color="primary" text-color="white" size="36px">{{ draft.hour }}ª</q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ draft.subject }} - {{ t('gradesPage.date') }}: {{ draft.date }}</q-item-label>
                <q-item-label caption class="text-grey-8">{{ draft.topic }}</q-item-label>
                <q-item-label caption v-if="draft.homework" class="text-indigo">{{ t('classRegister.assignHomework') }}: {{ draft.homework }}</q-item-label>
              </q-item-section>
              <q-item-section side class="row items-center q-gutter-xs">
                <q-btn color="positive" size="sm" icon="check" :label="t('dashboardPage.signAndRegister')" @click="registerDraftNow(draft, idx)" />
                <q-btn flat round dense icon="delete" color="negative" @click="deleteDraft(idx)" />
              </q-item-section>
            </q-item>
          </q-list>
          <div v-else class="text-center text-grey-6 q-pa-xl">
            {{ t('dashboardPage.noDrafts') }}
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="t('common.close')" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import dashboardService from 'src/services/dashboardService'
import api from '@/services/api'
import { useClassesStore } from '@/stores/classes'

import { useI18n } from 'vue-i18n'

const router = useRouter()
const authStore = useAuthStore()
const classesStore = useClassesStore()
const { user, userRole } = storeToRefs(authStore)
const $q = useQuasar()

const { t, locale: currentLocale } = useI18n()

const realStats = ref([])
const recentEvents = ref([])
const announcements = ref([])
const todaySchedule = ref([])
const loadingData = ref(false)
const loadingAnnouncements = ref(false)
const navigatingAction = ref(null)

const showDraftDialog = ref(false)
const showDraftsListDialog = ref(false)

const lessonDrafts = ref([])
const draftForm = ref({
  date: new Date().toISOString().substring(0, 10),
  class_id: null,
  hour: null,
  subject: '',
  topic: '',
  homework: ''
})

const currentRole = computed(() => userRole.value || user.value?.role || 'student')

const isDashboardAdmin = computed(() =>
  ['secretary', 'admin', 'superadmin', 'principal', 'vice_principal', 'system_auditor'].includes(currentRole.value)
)

const today = computed(() => {
  try {
    return new Date().toLocaleDateString(currentLocale.value || 'it-IT', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
  } catch (e) {
    return new Date().toLocaleDateString('it-IT', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
  }
})



const latestAnnouncement = computed(() => {
  if (announcements.value && announcements.value.length > 0) {
    return announcements.value[0]
  }
  return null
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return t('dashboardPage.greetingMorning') || 'Buongiorno'
  if (hour < 18) return t('dashboardPage.greetingAfternoon') || 'Buon pomeriggio'
  return t('dashboardPage.greetingEvening') || 'Buonasera'
})

const stats = computed(() => {
  if (realStats.value && realStats.value.length > 0) {
    return realStats.value
  }
  return [
    { label: 'Dati', value: '-', icon: 'info', color: 'grey' }
  ]
})

const displaySchedule = computed(() => {
  if (todaySchedule.value && todaySchedule.value.length > 0) {
    return todaySchedule.value
  }
  // Convert drafts for today into visible schedule items
  const todayStr = new Date().toISOString().substring(0, 10)
  const draftsForToday = lessonDrafts.value.filter(d => d.date === todayStr)
  return draftsForToday.map((d, index) => ({
    id: d.id || `draft-${Date.now()}-${index}`,
    hour_index: d.hour,
    subject_name: `${d.subject} (${t('dashboardPage.draft') || 'Bozza'})`,
    teacher_name: d.topic || 'Bozza preparata',
    room: 'Aula 2A',
    is_draft: true
  }))
})

const getLessonStatus = (entry) => {
  if (entry.is_draft) return t('dashboardPage.draft') || 'Bozza Pianificata'
  const now = new Date()
  const currentMinutes = now.getHours() * 60 + now.getMinutes()
  const hourIndex = Math.max(1, entry?.hour_index || 1)
  const startMinutes = (8 * 60) + ((hourIndex - 1) * 60)
  const endMinutes = startMinutes + 60
  if (currentMinutes >= endMinutes) return t('dashboardPage.completed')
  if (currentMinutes >= startMinutes) return t('dashboardPage.inProgress')
  return t('dashboardPage.scheduled')
}

// Uses internal status code instead of translated string to avoid locale-dependent comparisons
const getLessonStatusCode = (entry) => {
  if (entry.is_draft) return 'draft'
  const now = new Date()
  const currentMinutes = now.getHours() * 60 + now.getMinutes()
  const hourIndex = Math.max(1, entry?.hour_index || 1)
  const startMinutes = (8 * 60) + ((hourIndex - 1) * 60)
  const endMinutes = startMinutes + 60
  if (currentMinutes >= endMinutes) return 'completed'
  if (currentMinutes >= startMinutes) return 'inProgress'
  return 'scheduled'
}

const getLessonStatusColor = (entry) => {
  const code = getLessonStatusCode(entry)
  if (code === 'draft') return 'amber-2'
  if (code === 'completed') return 'grey-3'
  if (code === 'inProgress') return 'positive'
  return 'blue-1'
}

const getLessonStatusTextColor = (entry) => {
  const code = getLessonStatusCode(entry)
  if (code === 'draft') return 'amber-9'
  if (code === 'completed') return 'grey-7'
  if (code === 'inProgress') return 'white'
  return 'primary'
}

const openDraftModal = () => {
  draftForm.value = {
    date: new Date().toISOString().substring(0, 10),
    class_id: null,
    hour: null,
    subject: '',
    topic: '',
    homework: ''
  }
  showDraftDialog.value = true
}

const openDraftsList = () => {
  showDraftsListDialog.value = true
}

const saveLessonDraft = () => {
  if (!draftForm.value.topic) {
    $q.notify({ type: 'warning', message: t('dashboardPage.topicRequired') })
    return
  }
  const newDraft = { ...draftForm.value, id: 'draft_' + Date.now() }
  lessonDrafts.value.push(newDraft)
  sessionStorage.setItem('registro_lesson_drafts', JSON.stringify(lessonDrafts.value))
  $q.notify({ type: 'positive', message: t('dashboardPage.draftSaved', { date: newDraft.date }) })
  showDraftDialog.value = false
}

const deleteDraft = (idx) => {
  lessonDrafts.value.splice(idx, 1)
  sessionStorage.setItem('registro_lesson_drafts', JSON.stringify(lessonDrafts.value))
  $q.notify({ type: 'info', message: t('dashboardPage.deleteDraft') })
}

const registerDraftNow = async (draft, idx) => {
  try {
    await api.post('/lessons', {
      class_id: draft.class_id,
      // subject_id comes from the draft if available; the server must not require a fixed ID
      ...(draft.subject_id ? { subject_id: draft.subject_id } : {}),
      date: draft.date,
      hour: draft.hour,
      duration: 1,
      topic: draft.topic,
      notes: draft.homework
    })
    $q.notify({ type: 'positive', message: t('dashboardPage.draftRegistered') })
    deleteDraft(idx)
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
    console.error('Error registering draft lesson:', err)
  }
}

const loadStoredDrafts = () => {
  try {
    // Use sessionStorage: draft data may contain student names / homework — not appropriate for localStorage
    const raw = sessionStorage.getItem('registro_lesson_drafts')
    if (raw) {
      lessonDrafts.value = JSON.parse(raw)
    }
  } catch { /* ignore */ }
}

const fetchDashboardData = async () => {
    loadingData.value = true
    try {
        const role = currentRole.value
        const data = await dashboardService.getDashboardStats(role)
        if (data) {
            if (role === 'admin' || role === 'superadmin') {
                realStats.value = [
                    { label: t('dashboardPage.statTotalSchools'), value: data.total_schools ?? '0', icon: 'school', color: 'indigo' },
                    { label: t('dashboardPage.statActiveUsers'), value: data.total_users ?? '0', icon: 'people', color: 'cyan' },
                    { label: t('dashboardPage.statActive24h'), value: data.active_users_24h ?? '0', icon: 'event', color: 'amber' },
                    { label: t('dashboardPage.statPendingDocs'), value: data.pending_documents_count ?? '0', icon: 'assignment', color: 'red' }
                ]
                recentEvents.value = data.recent_events || []
            } else if (role === 'secretary') {
                realStats.value = [
                    { label: t('dashboardPage.statStudents'), value: data.total_students ?? '0', icon: 'school', color: 'indigo' },
                    { label: t('dashboardPage.statTeachers'), value: data.total_teachers ?? '0', icon: 'people', color: 'cyan' },
                    { label: t('dashboardPage.statDocuments'), value: data.total_documents ?? '0', icon: 'description', color: 'amber' },
                    { label: t('dashboardPage.statRequests'), value: data.pending_documents_count ?? '0', icon: 'assignment', color: 'red' }
                ]
                recentEvents.value = data.recent_events || []
            } else if (role === 'teacher') {
                realStats.value = [
                    { label: t('dashboardPage.statMyClasses'), value: data.classes_count ?? '0', icon: 'class', color: 'indigo' },
                    { label: t('dashboardPage.statStudents'), value: data.students_count ?? '0', icon: 'school', color: 'cyan' },
                    { label: t('dashboardPage.statLessonsToday'), value: data.lessons_today_count ?? '0', icon: 'event', color: 'amber' },
                    { label: t('dashboardPage.statGradesPending'), value: data.grades_pending_count ?? '0', icon: 'grade', color: 'red' }
                ]
            } else if (role === 'student') {
                realStats.value = [
                    { label: t('roleDashboards.averageGrade'), value: data.average_grade ?? '-', icon: 'grade', color: 'indigo' },
                    { label: t('roleDashboards.attendanceRate'), value: data.attendance_rate != null ? data.attendance_rate + '%' : '-', icon: 'how_to_reg', color: 'cyan' },
                    { label: t('agendaPage.homework'), value: data.homework_count ?? '0', icon: 'assignment', color: 'amber' },
                    { label: t('documentsPage.title'), value: data.documents_count ?? '0', icon: 'description', color: 'purple' }
                ]
            } else if (role === 'parent') {
                realStats.value = [
                    { label: t('nav.myChildren'), value: data.children_count ?? '0', icon: 'family_restroom', color: 'indigo' },
                    { label: t('nav.colloqui'), value: data.upcoming_colloqui ?? '0', icon: 'event', color: 'cyan' },
                    { label: t('nav.communications'), value: data.unread_communications ?? '0', icon: 'email', color: 'amber' },
                    { label: t('documentsPage.title'), value: data.documents_count ?? '0', icon: 'description', color: 'purple' }
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
    try {
        return date.toLocaleString(currentLocale.value || 'it-IT')
    } catch {
        return date.toLocaleString('it-IT')
    }
}

const actions = computed(() => {
  const role = currentRole.value
  if (role === 'teacher' || role === 'coordinator') {
    return [
      { key: 'attendance', label: t('dashboardPage.actionAttendance'), icon: 'how_to_reg', route: '/teacher/attendance' },
      { key: 'grades', label: t('dashboardPage.actionGrades'), icon: 'grade', route: '/teacher/grades' },
      { key: 'lessons', label: t('dashboardPage.actionLessons'), icon: 'edit_calendar', route: '/teacher/lessons' },
      { key: 'agenda', label: t('dashboardPage.actionAgenda'), icon: 'event', route: '/teacher/agenda' }
    ]
  }
  if (role === 'secretary' || role === 'principal' || role === 'vice_principal') {
    return [
      { key: 'users', label: t('dashboardPage.actionUsers') || 'Utenti', icon: 'people', route: '/secretary/users' },
      { key: 'classes', label: t('dashboardPage.actionClasses') || 'Classi', icon: 'school', route: '/secretary/classes' },
      { key: 'documents', label: t('documentsPage.title') || 'Documenti', icon: 'folder', route: '/secretary/documents' },
      { key: 'reports', label: t('reportsPage.title') || 'Report', icon: 'bar_chart', route: '/secretary/reports' }
    ]
  }
  if (role === 'student') {
    return [
      { key: 'grades', label: t('dashboardPage.actionGrades') || 'Voti', icon: 'grade', route: '/student/grades' },
      { key: 'attendance', label: t('dashboardPage.actionAttendance') || 'Presenze', icon: 'how_to_reg', route: '/student/attendance' },
      { key: 'homework', label: t('agendaPage.homework') || 'Compiti', icon: 'assignment', route: '/student/homework' },
      { key: 'timetable', label: t('timetablePage.title') || 'Orario', icon: 'schedule', route: '/student/timetable' }
    ]
  }
  if (role === 'parent') {
    return [
      { key: 'grades', label: t('dashboardPage.actionGrades') || 'Voti', icon: 'grade', route: '/parent/grades' },
      { key: 'attendance', label: t('dashboardPage.actionAttendance') || 'Presenze', icon: 'how_to_reg', route: '/parent/attendance' },
      { key: 'colloqui', label: t('nav.colloqui') || 'Colloqui', icon: 'event', route: '/parent/colloqui' },
      { key: 'communications', label: t('nav.communications') || 'Comunicazioni', icon: 'email', route: '/parent/communications' }
    ]
  }
  // admin / superadmin / system_auditor
  return [
    { key: 'users', label: t('dashboardPage.actionUsers') || 'Utenti', icon: 'people', route: '/admin/users' },
    { key: 'schools', label: t('dashboardPage.schoolManagement') || 'Scuole', icon: 'school', route: '/admin/schools' },
    { key: 'analytics', label: t('nav.analytics') || 'Analytics', icon: 'bar_chart', route: '/admin/analytics' },
    { key: 'settings', label: t('settings.title') || 'Impostazioni', icon: 'settings', route: '/admin/settings' }
  ]
})

const handleActionClick = async (action) => {
  if (action.route) {
    navigatingAction.value = action.key
    await router.push(action.route)
    navigatingAction.value = null
  }
}

onMounted(() => {
  loadStoredDrafts()
  fetchDashboardData()
})
</script>
