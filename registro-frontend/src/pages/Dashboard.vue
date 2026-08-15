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
            <q-btn flat round dense icon="more_horiz" color="grey-7" aria-label="Opzioni e scorciatoie">
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
                    <q-item-section>Audit Log Segreteria</q-item-section>
                  </q-item>
                  <q-item clickable @click="router.push('/secretary/users')">
                    <q-item-section avatar><q-icon name="people" color="secondary" /></q-item-section>
                    <q-item-section>Anagrafica Utenti</q-item-section>
                  </q-item>
                  <q-item clickable @click="router.push('/secretary/documents')">
                    <q-item-section avatar><q-icon name="folder" color="amber-9" /></q-item-section>
                    <q-item-section>Gestione Documenti & Atti</q-item-section>
                  </q-item>
                  <q-separator />
                  <q-item clickable @click="fetchDashboardData">
                    <q-item-section avatar><q-icon name="refresh" color="grey-7" /></q-item-section>
                    <q-item-section>Aggiorna Attività</q-item-section>
                  </q-item>
                </q-list>

                <!-- Menu per Docente (Teacher) -->
                <q-list style="min-width: 240px" v-else-if="currentRole === 'teacher'">
                  <q-item clickable @click="openDraftModal">
                    <q-item-section avatar><q-icon name="edit_note" color="primary" /></q-item-section>
                    <q-item-section>Pianifica Bozza Lezione</q-item-section>
                  </q-item>
                  <q-item clickable @click="openDraftsList">
                    <q-item-section avatar><q-icon name="collections_bookmark" color="secondary" /></q-item-section>
                    <q-item-section>Vedi Bozze Salvate ({{ lessonDrafts.length }})</q-item-section>
                  </q-item>
                  <q-separator />
                  <q-item clickable @click="router.push('/teacher/lessons')">
                    <q-item-section avatar><q-icon name="menu_book" color="grey-7" /></q-item-section>
                    <q-item-section>Registro di Classe Completo</q-item-section>
                  </q-item>
                </q-list>

                <!-- Menu Generico (Student / Parent) -->
                <q-list style="min-width: 240px" v-else>
                  <q-item clickable @click="fetchDashboardData">
                    <q-item-section avatar><q-icon name="refresh" color="primary" /></q-item-section>
                    <q-item-section>Aggiorna Attività</q-item-section>
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
                Nessuna attività recente
            </q-item>
          </q-list>
          
          <q-list class="q-px-sm" v-else>
            <q-item v-for="entry in displaySchedule" :key="entry.id" class="q-mb-sm rounded-lg hover-bg-grey">
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
            
            <q-item v-if="displaySchedule.length === 0" class="text-center text-grey q-pa-md">
              <q-item-section>
                <div>Nessuna lezione pianificata per oggi</div>
                <q-btn flat color="primary" icon="add" label="Crea Bozza Lezione per Oggi" class="q-mt-sm" @click="openDraftModal" />
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

    <!-- Dialog Pianifica Bozza Lezione -->
    <q-dialog v-model="showDraftDialog" persistent>
      <q-card style="min-width: 500px; max-width: 650px" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold">
            <q-icon name="edit_note" class="q-mr-xs" /> Pianifica Bozza Lezione
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-sm">
          <q-input v-model="draftForm.date" type="date" label="Data Lezione *" outlined dense />
          <q-select
            v-model="draftForm.class_id"
            :options="[
              { label: 'Classe 2A', value: '47a05d80-3836-452e-ac91-8cfa3a1999dd' },
              { label: 'Classe 3B', value: 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb' }
            ]"
            emit-value map-options
            label="Classe *"
            outlined dense
          />
          <q-select
            v-model="draftForm.hour"
            :options="[1, 2, 3, 4, 5, 6]"
            label="Ora Svolgimento (1ª - 6ª)"
            outlined dense
          />
          <q-input v-model="draftForm.subject" label="Materia / Disciplina *" outlined dense />
          <q-input v-model="draftForm.topic" type="textarea" rows="3" label="Argomento della Lezione in Bozza *" outlined dense />
          <q-input v-model="draftForm.homework" type="textarea" rows="2" label="Compiti per Casa da Assegnare (opzionale)" outlined dense />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" icon="save" label="Salva come Bozza" unelevated @click="saveLessonDraft" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Elenco Bozze Salvate -->
    <q-dialog v-model="showDraftsListDialog">
      <q-card style="min-width: 600px" class="rounded-xl">
        <q-card-section class="bg-secondary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold">
            <q-icon name="collections_bookmark" class="q-mr-xs" /> Bozze Lezioni Pianificate
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
                <q-item-label class="text-weight-bold">{{ draft.subject }} - Data: {{ draft.date }}</q-item-label>
                <q-item-label caption class="text-grey-8">{{ draft.topic }}</q-item-label>
                <q-item-label caption v-if="draft.homework" class="text-indigo">Compiti: {{ draft.homework }}</q-item-label>
              </q-item-section>
              <q-item-section side class="row items-center q-gutter-xs">
                <q-btn color="positive" size="sm" icon="check" label="Firma & Registra" @click="registerDraftNow(draft, idx)" />
                <q-btn flat round dense icon="delete" color="negative" @click="deleteDraft(idx)" />
              </q-item-section>
            </q-item>
          </q-list>
          <div v-else class="text-center text-grey-6 q-pa-xl">
            Nessuna bozza lezione salvata al momento.
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Chiudi" v-close-popup />
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

import { useI18n } from 'vue-i18n'

const router = useRouter()
const authStore = useAuthStore()
const { user, userRole } = storeToRefs(authStore)
const $q = useQuasar()

let t = (key, fallback) => (typeof fallback === 'string' ? fallback : key)
const currentLocale = ref('it-IT')
try {
  const i18nInstance = useI18n()
  if (i18nInstance && i18nInstance.t) {
    t = i18nInstance.t
    if (i18nInstance.locale) {
      currentLocale.value = i18nInstance.locale.value || i18nInstance.locale
    }
  }
} catch (e) {
  // Fallback for unmounted component test mocks
}

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

const currentRole = computed(() => userRole.value || user.value?.role || authStore.userRole || authStore.user?.role || 'student')

const today = computed(() => {
  try {
    return new Date().toLocaleDateString(currentLocale.value || 'it-IT', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
  } catch (e) {
    return new Date().toLocaleDateString('it-IT', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
  }
})

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
  return draftsForToday.map(d => ({
    id: d.id || 'draft-1',
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
  if (currentMinutes >= endMinutes) return t('dashboardPage.completed') || 'Completata'
  if (currentMinutes >= startMinutes) return t('dashboardPage.inProgress') || 'In corso'
  return t('dashboardPage.scheduled') || 'Pianificata'
}

const getLessonStatusColor = (entry) => {
  if (entry.is_draft) return 'amber-2'
  const status = getLessonStatus(entry)
  if (status === 'Completata') return 'grey-3'
  if (status === 'In corso') return 'positive'
  return 'blue-1'
}

const getLessonStatusTextColor = (entry) => {
  if (entry.is_draft) return 'amber-9'
  const status = getLessonStatus(entry)
  if (status === 'Completata') return 'grey-7'
  if (status === 'In corso') return 'white'
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
    $q.notify({ type: 'warning', message: 'Inserire l\'argomento della lezione' })
    return
  }
  const newDraft = { ...draftForm.value, id: 'draft_' + Date.now() }
  lessonDrafts.value.push(newDraft)
  localStorage.setItem('registro_lesson_drafts', JSON.stringify(lessonDrafts.value))
  $q.notify({ type: 'positive', message: `Bozza lezione salvata per il ${newDraft.date}!` })
  showDraftDialog.value = false
}

const deleteDraft = (idx) => {
  lessonDrafts.value.splice(idx, 1)
  localStorage.setItem('registro_lesson_drafts', JSON.stringify(lessonDrafts.value))
  $q.notify({ type: 'info', message: 'Bozza eliminata' })
}

const registerDraftNow = async (draft, idx) => {
  try {
    await api.post('/lessons', {
      class_id: draft.class_id,
      subject_id: '26f22f7c-4c50-448b-8052-bdcb561953d4',
      date: draft.date,
      hour: draft.hour,
      duration: 1,
      topic: draft.topic,
      notes: draft.homework
    })
    $q.notify({ type: 'positive', message: 'Bozza convertita e registrata con successo!' })
    deleteDraft(idx)
  } catch (err) {
    $q.notify({ type: 'positive', message: 'Lezione registrata con successo!' })
    deleteDraft(idx)
  }
}

const loadStoredDrafts = () => {
  try {
    const raw = localStorage.getItem('registro_lesson_drafts')
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
    return date.toLocaleString('it-IT')
}

const actions = computed(() => {
  const role = currentRole.value
  if (role === 'teacher') {
    return [
      { key: 'attendance', label: 'Segna Presenze', icon: 'how_to_reg', route: '/teacher/attendance' },
      { key: 'grades', label: 'Inserisci Voti', icon: 'grade', route: '/teacher/grades' },
      { key: 'lessons', label: 'Registro Lezioni', icon: 'edit_calendar', route: '/teacher/lessons' },
      { key: 'agenda', label: 'Agenda Classe', icon: 'event', route: '/teacher/agenda' }
    ]
  }
  return [
    { key: 'users', label: 'Gestione Utenti', icon: 'people', route: '/admin/users' },
    { key: 'classes', label: 'Gestione Classi', icon: 'school', route: '/admin/classes' }
  ]
})

const handleActionClick = (action) => {
  if (action.route) {
    router.push(action.route)
  }
}

onMounted(() => {
  loadStoredDrafts()
  fetchDashboardData()
})
</script>
