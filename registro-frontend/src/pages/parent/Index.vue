<template>
  <q-page class="q-pa-md" role="main">
    <!-- Header with Child Switcher -->
    <div class="row items-center justify-between q-mb-xl">
      <div>
        <h1 class="text-h3 text-weight-bold text-outfit parent-heading q-my-none" style="display: inline-block;">
          {{ $t('roleDashboards.parentPanel', { name: parentName }) }}
        </h1>
        <div class="text-subtitle1 text-slate-600 q-mt-sm" aria-live="polite">{{ $t('roleDashboards.parentSub') }}</div>
      </div>
      <div v-if="children.length > 0">
        <q-btn-dropdown
          color="primary"
          unelevated
          no-caps
          class="rounded-xl shadow-soft q-px-md"
          :label="selectedChild ? `${selectedChild.first_name || selectedChild.firstName} ${selectedChild.last_name || selectedChild.lastName}` : $t('roleDashboards.selectChild')"
          icon="face"
          :aria-label="$t('parentAria.selectChild')"
        >
          <q-list role="listbox" :aria-label="$t('parentAria.childrenList')">
            <q-item
              v-for="child in children"
              :key="child.id"
              clickable
              v-close-popup
              @click="selectChild(child.id)"
              :active="selectedChildId === child.id"
              active-class="bg-blue-1 text-primary"
              role="option"
              :aria-selected="selectedChildId === child.id"
            >
              <q-item-section avatar>
                <q-avatar size="sm" color="primary" text-color="white" :aria-label="$t('parentAria.childInitial', { name: child.firstName || child.first_name })">{{ (child.firstName || child.first_name || '?').charAt(0) }}</q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ child.firstName || child.first_name }} {{ child.lastName || child.last_name }}</q-item-label>
                <q-item-label caption>{{ child.class }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </div>
    </div>

    <!-- Loading State with Skeletons & Spinner -->
    <div v-if="loading" class="row justify-center q-pa-lg" role="status" :aria-label="$t('parentAria.loading')">
      <q-spinner color="primary" size="3em" class="q-mb-md" />
      <div class="row q-col-gutter-md full-width">
        <div v-for="n in 4" :key="n" class="col-12 col-sm-6 col-md-3">
          <q-skeleton type="rect" height="120px" class="rounded-borders" />
        </div>
      </div>
    </div>

    <!-- Dashboard Content -->
    <div v-else-if="selectedChild" class="row q-col-gutter-md">

      <!-- Quick Stats -->
      <div class="col-12 col-sm-6 col-md-3" role="region" :aria-label="$t('parentAria.averageGrade')">
        <q-card class="glass-card stat-card shadow-soft full-height overflow-hidden">
          <q-card-section>
            <div class="text-caption text-slate-600 text-uppercase letter-spacing-1" style="font-size: 12px">{{ $t('roleDashboards.averageGrade') }}</div>
            <div class="text-h3 text-weight-bold text-indigo-700 q-mt-sm" :aria-label="`${$t('parentAria.averageGrade')}: ${averageGrade}`">{{ averageGrade }}</div>
            <div class="row items-center q-mt-sm">
              <q-icon name="trending_up" color="positive" class="q-mr-xs" aria-hidden="true" />
              <span class="text-positive text-caption text-weight-medium">{{ $t('roleDashboards.generalTrend') }}</span>
            </div>
          </q-card-section>
          <q-icon name="grade" class="card-bg-icon text-indigo-100" aria-hidden="true" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3" role="region" :aria-label="$t('roleDashboards.totalAbsences')">
        <q-card class="glass-card stat-card shadow-soft full-height overflow-hidden">
          <q-card-section>
            <div class="text-caption text-slate-600 text-uppercase letter-spacing-1" style="font-size: 12px">{{ $t('roleDashboards.totalAbsences') }}</div>
            <div class="text-h3 text-weight-bold text-orange-700 q-mt-sm" :aria-label="`${$t('roleDashboards.totalAbsences')}: ${totalAbsences}`">{{ totalAbsences }}</div>
            <div class="row items-center q-mt-sm">
              <span class="text-caption text-slate-600">{{ $t('roleDashboards.currentYear') }}</span>
            </div>
          </q-card-section>
          <q-icon name="how_to_reg" class="card-bg-icon text-orange-100" aria-hidden="true" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3" role="region" :aria-label="$t('colloquiPage.title') || 'Prossimo Colloquio'">
        <q-card class="glass-card stat-card shadow-soft full-height overflow-hidden">
          <q-card-section>
            <div class="text-caption text-slate-600 text-uppercase letter-spacing-1" style="font-size: 12px">{{ $t('colloquiPage.title') || 'Prossimo Colloquio' }}</div>
            <div v-if="nextColloquio" class="text-h6 text-weight-bold text-slate-800 q-mt-sm">{{ formatDate(nextColloquio.date) }}</div>
            <div v-else class="text-h6 text-weight-bold text-slate-500 q-mt-sm">-</div>
            <q-btn flat dense no-caps color="primary" :label="$t('colloquiPage.booked') || 'Prenota ora'" to="/parent/colloqui" class="q-mt-sm rounded-lg" :aria-label="$t('parentAria.bookMeeting')" />
          </q-card-section>
          <q-icon name="event" class="card-bg-icon text-slate-100" aria-hidden="true" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3" role="region" :aria-label="$t('communicationsPage.title')">
        <q-card class="glass-card stat-card shadow-soft full-height overflow-hidden">
          <q-card-section>
            <div class="text-caption text-slate-600 text-uppercase letter-spacing-1" style="font-size: 12px">{{ $t('communicationsPage.title') }}</div>
            <div class="text-h3 text-weight-bold text-rose-700 q-mt-sm" :aria-label="`${unreadCount} avvisi`">{{ unreadCount }}</div>
            <div class="text-caption text-slate-600 q-mt-sm text-weight-medium">Da leggere</div>
          </q-card-section>
          <q-icon name="notifications_active" class="card-bg-icon text-rose-100" aria-hidden="true" />
        </q-card>
      </div>

      <!-- Analytics Charts: Grade Trend & Presence -->
      <div class="col-12 q-mb-md">
        <GradeAnalyticsCharts
          :grades="allChildGrades"
          :attendance-rate="childAttendanceRate"
        />
      </div>

      <!-- Absence Limit 25% Monitoring Widget (Art. 14 DPR 122/2009) -->
      <div class="col-12 q-mb-md">
        <AbsenceLimitWidget
          :student-id="selectedChildId"
          :student-name="selectedChild ? `${selectedChild.first_name || selectedChild.firstName} ${selectedChild.last_name || selectedChild.lastName}` : ''"
        />
      </div>

      <!-- Recent Grades -->
      <div class="col-12 col-md-8">
        <q-card class="shadow-sm rounded-lg" role="region" :aria-label="$t('gradesPage.title')">
          <q-card-section class="row items-center justify-between">
            <div class="text-h6 text-slate-800">{{ $t('gradesPage.title') }}</div>
            <q-btn flat no-caps color="primary" :label="$t('common.viewAll') || 'Vedi tutti'" to="/parent/grades" :aria-label="$t('parentAria.viewAllGrades')" />
          </q-card-section>
          <q-separator />
          <q-list separator v-if="recentGrades.length > 0" role="list" :aria-label="$t('gradesPage.title')">
            <q-item v-for="grade in recentGrades" :key="grade.id" role="listitem">
              <q-item-section>
                <q-item-label class="text-weight-medium text-slate-800">{{ grade.subject_name || grade.subject_id }}</q-item-label>
                <q-item-label caption class="text-slate-600">{{ grade.grade_type }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <div class="row items-center">
                  <q-badge
                    :color="grade.grade_value === -1 ? 'grey' : (grade.grade_value >= 6 ? 'positive' : 'negative')"
                    class="text-subtitle1 q-pa-xs"
                    :aria-label="`Voto: ${grade.grade_value === -1 ? 'Assente' : grade.grade_value}`"
                  >
                    {{ grade.grade_value === -1 ? 'A' : grade.grade_value }}
                  </q-badge>
                  <div class="text-caption text-slate-600 q-ml-md">{{ new Date(grade.date).toLocaleDateString('it-IT') }}</div>
                </div>
              </q-item-section>
            </q-item>
          </q-list>
          <div v-else class="q-pa-lg text-center text-slate-500" role="status">{{ $t('parentAria.noRecentGrades') }}</div>
        </q-card>
      </div>

      <!-- Upcoming Events -->
      <div class="col-12 col-md-4">
        <q-card class="shadow-sm rounded-lg full-height">
          <q-card-section>
            <div class="text-h6 text-slate-800 q-mb-sm">{{ $t('parentAria.upcomingEvents') }}</div>
            <div v-if="upcomingTests.length === 0" class="text-center text-grey q-pa-md">
              <q-icon name="event_available" size="2em" color="grey-4" class="q-mb-sm" />
              <div class="text-caption">{{ $t('parentAria.noUpcomingEvents') }}</div>
            </div>
            <q-timeline v-else color="primary" layout="dense">
              <q-timeline-entry
                v-for="test in upcomingTests"
                :key="test.id"
                :title="test.title"
                :subtitle="formatDate(test.date)"
                :icon="evalTypeIcon(test.evaluation_type)"
                :color="evalTypeColor(test.evaluation_type)"
              >
                <div v-if="test.parent_notes" class="text-caption text-grey-7">
                  {{ test.parent_notes }}
                </div>
              </q-timeline-entry>
            </q-timeline>
          </q-card-section>
        </q-card>
      </div>

    </div>

    <!-- Empty State -->
    <div v-else class="text-center q-pa-xl" role="status">
      <q-icon name="family_restroom" size="4em" color="grey-6" aria-hidden="true" />
      <div class="text-h6 text-slate-600 q-mt-sm">{{ $t('parentAria.noChildAssociated') }}</div>
      <p class="text-slate-500">{{ $t('parentAria.contactSecretary') }}</p>
    </div>

    <!-- Quick Actions (FAB on Mobile) -->
    <q-page-sticky position="bottom-right" :offset="[18, 18]" class="lt-md">
      <q-fab icon="add" direction="up" color="primary" :aria-label="$t('parentAria.quickActions')">
        <q-fab-action color="orange" icon="edit_calendar" label="Giustifica" to="/parent/attendance" :aria-label="$t('parentAria.justifyAbsence')" />
        <q-fab-action color="secondary" icon="event" label="Colloquio" to="/parent/colloqui" :aria-label="$t('parentAria.bookMeeting')" />
      </q-fab>
    </q-page-sticky>

  </q-page>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useParentStore } from '@/stores/parent'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { gradeService } from '@/services/gradeService'
import { attendanceService } from '@/services/attendanceService'
import { communicationService } from '@/services/communicationService'
import { colloquiService } from '@/services/colloquiService'
import GradeAnalyticsCharts from '@/components/Student/GradeAnalyticsCharts.vue'
import AbsenceLimitWidget from '@/components/Parent/AbsenceLimitWidget.vue'



const parentStore = useParentStore()
const authStore = useAuthStore()
const { children, selectedChild, selectedChildId, loading } = storeToRefs(parentStore)
const { fetchChildren, selectChild } = parentStore

const parentName = computed(() => authStore.user?.first_name || authStore.user?.name || 'Genitore')

const averageGrade = ref('-')
const totalAbsences = ref(0)
const childAttendanceRate = ref(100)
const unreadCount = ref(0)
const recentGrades = ref([])
const allChildGrades = ref([])
const upcomingTests = ref([])
const nextColloquio = ref(null)
const dataLoading = ref(false)

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('it-IT', { day: '2-digit', month: 'short', year: 'numeric' })
}

const evalTypeIcon = (type) => {
  const icons = {
    'scritto': 'edit_note',
    'orale': 'record_voice_over',
    'pratico': 'build',
    'test': 'quiz',
  }
  return icons[type?.toLowerCase()] || 'assignment'
}

const evalTypeColor = (type) => {
  const colors = {
    'scritto': 'blue',
    'orale': 'orange',
    'pratico': 'green',
    'test': 'purple',
  }
  return colors[type?.toLowerCase()] || 'primary'
}

const fetchChildData = async () => {
    if (!selectedChildId.value) return

    dataLoading.value = true
    try {
        // Fetch Grades
        const gradesRes = await gradeService.getChildGrades(selectedChildId.value)
        const allGrades = []
        if (gradesRes.data && gradesRes.data.semesters) {
            gradesRes.data.semesters.forEach(s => {
                if (s.grades) allGrades.push(...s.grades)
            })
        }

        const validGrades = allGrades.filter(g => g.grade_value >= 0)
        if (validGrades.length > 0) {
            const sum = validGrades.reduce((acc, g) => acc + g.grade_value, 0)
            averageGrade.value = (sum / validGrades.length).toFixed(1)
        } else {
            averageGrade.value = '-'
        }
        allGrades.sort((a, b) => new Date(b.date) - new Date(a.date))
        allChildGrades.value = allGrades
        recentGrades.value = allGrades.slice(0, 5)

        // Fetch Attendance
        const attRes = await attendanceService.getChildAttendance(selectedChildId.value)
        if (attRes.data) {
            const records = Array.isArray(attRes.data) ? attRes.data : (attRes.data.records || [])
            const absences = records.filter(r => r.status === 'absent').length
            totalAbsences.value = absences
            if (records.length > 0) {
                childAttendanceRate.value = Math.round(((records.length - absences) / records.length) * 100)
            } else {
                childAttendanceRate.value = 100
            }
        }

        // Fetch Unread Communications
        try {
            const commsRes = await communicationService.getMessages()
            const messages = commsRes.data || []
            unreadCount.value = messages.filter(m => !m.is_read).length
        } catch (e) {
            unreadCount.value = 0
        }

        // Fetch next colloquio
        try {
            const colloquiRes = await colloquiService.getBookedSlots()
            const slots = Array.isArray(colloquiRes.data) ? colloquiRes.data : (colloquiRes.data?.items || [])
            const now = new Date()
            const upcoming = slots
                .filter(s => new Date(s.date) > now)
                .sort((a, b) => new Date(a.date) - new Date(b.date))
            nextColloquio.value = upcoming[0] || null
        } catch (e) {
            nextColloquio.value = null
        }

        // Fetch Upcoming Tests
        const child = selectedChild.value
        if (child && child.class_id) {
            try {
                const testsRes = await gradeService.getUpcomingTestsForClass(child.class_id)
                upcomingTests.value = Array.isArray(testsRes.data) ? testsRes.data : []
            } catch (e) {
                console.warn('Could not fetch upcoming tests:', e)
                upcomingTests.value = []
            }
        } else {
            upcomingTests.value = []
        }
    } catch (e) {
        console.error('Error fetching child data', e)
    } finally {
        dataLoading.value = false
    }
}

watch(selectedChildId, () => {
    fetchChildData()
})

onMounted(async () => {
  if (children.value.length === 0) {
    await fetchChildren()
  }
  fetchChildData()
})
</script>

<style scoped>
.parent-heading {
  color: #312e81;
}

.body--dark .parent-heading {
  color: #a5b4fc;
}

.letter-spacing-1 {
    letter-spacing: 1px;
}

.stat-card {
  position: relative;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.card-bg-icon {
  position: absolute;
  right: -10px;
  bottom: -10px;
  font-size: 80px;
  opacity: 0.5;
  z-index: 0;
}

.truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
