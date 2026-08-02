<template>
  <q-page class="q-pa-lg bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-xl">
      <div>
        <h1 class="text-h3 text-weight-bold text-outfit q-my-none text-gradient-premium">
          Analytics di Sistema
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-sm q-mb-none">
          {{ isSuperAdmin ? 'Metriche globali della piattaforma e statistiche generali' : 'Reportistica ed utilizzo delle risorse per la tua scuola' }}
        </p>
      </div>
      <div class="col-auto">
        <q-btn
          unelevated
          color="white"
          text-color="primary"
          icon="refresh"
          label="Aggiorna Dati"
          class="rounded-lg shadow-soft q-px-md"
          @click="fetchAll"
          :loading="loading"
        />
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="loading && !stats" class="row q-col-gutter-lg">
      <div v-for="i in 4" :key="i" class="col-12 col-md-3">
        <q-skeleton height="120px" class="rounded-xl" />
      </div>
      <div class="col-12 col-md-8">
        <q-skeleton height="350px" class="rounded-xl" />
      </div>
      <div class="col-12 col-md-4">
        <q-skeleton height="350px" class="rounded-xl" />
      </div>
    </div>

    <div v-else class="q-gutter-y-xl">
      <!-- KPI Overview Cards -->
      <div class="row q-col-gutter-lg">
        <div v-if="isSuperAdmin" class="col-12 col-sm-6 col-md-3">
          <q-card class="glass-card stat-card full-height shadow-soft rounded-xl overflow-hidden border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Scuole Registrate</div>
                  <div class="text-h4 text-weight-bold text-outfit text-slate-800">{{ stats?.total_schools ?? '-' }}</div>
                </div>
                <div class="bg-indigo-100 q-pa-md rounded-xl">
                  <q-icon name="school" size="32px" color="indigo-700" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-sm-6" :class="isSuperAdmin ? 'col-md-3' : 'col-md-4'">
          <q-card class="glass-card stat-card full-height shadow-soft rounded-xl overflow-hidden border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Utenti Totali</div>
                  <div class="text-h4 text-weight-bold text-outfit text-slate-800">{{ stats?.total_users ?? '-' }}</div>
                </div>
                <div class="bg-purple-100 q-pa-md rounded-xl">
                  <q-icon name="people" size="32px" color="purple-700" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-sm-6" :class="isSuperAdmin ? 'col-md-3' : 'col-md-4'">
          <q-card class="glass-card stat-card full-height shadow-soft rounded-xl overflow-hidden border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Documenti Caricati</div>
                  <div class="text-h4 text-weight-bold text-outfit text-slate-800">{{ stats?.total_documents ?? '-' }}</div>
                </div>
                <div class="bg-teal-100 q-pa-md rounded-xl">
                  <q-icon name="description" size="32px" color="teal-700" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- API Latency from real metrics -->
        <div class="col-12 col-sm-6" :class="isSuperAdmin ? 'col-md-3' : 'col-md-4'">
          <q-card class="glass-card stat-card full-height shadow-soft rounded-xl overflow-hidden border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Latenza Server API</div>
                  <div class="text-h4 text-weight-bold text-outfit" :class="latencyColor">
                    {{ systemMetrics?.api_latency_ms != null ? systemMetrics.api_latency_ms + 'ms' : '-' }}
                  </div>
                </div>
                <div class="bg-emerald-100 q-pa-md rounded-xl">
                  <q-icon name="bolt" size="32px" color="emerald-700" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Main Analytics Graphs -->
      <div class="row q-col-gutter-lg">
        <!-- User Growth Trend -->
        <div class="col-12 col-md-8">
          <q-card class="glass-card shadow-soft rounded-xl border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between q-mb-md">
                <div>
                  <div class="text-h5 text-weight-bold text-outfit text-slate-800">Crescita Utenti Registrati</div>
                  <div class="text-caption text-slate-500">Andamento semestrale delle registrazioni in piattaforma</div>
                </div>
                <q-chip v-if="userGrowthTrend !== null" outline color="primary" class="text-weight-bold" size="sm">
                  Trend {{ userGrowthTrend >= 0 ? '+' : '' }}{{ userGrowthTrend }}%
                </q-chip>
              </div>

              <div v-if="chartPoints.length === 0" class="text-center text-slate-400 q-pa-xl">
                <q-icon name="bar_chart" size="48px" class="q-mb-sm" />
                <div>Dati storici non disponibili</div>
              </div>

              <div v-else class="relative-position q-py-md">
                <svg viewBox="0 0 500 200" class="full-width" style="overflow: visible;" role="img" aria-labelledby="chart-title chart-desc">
                  <title id="chart-title">Crescita Utenti Registrati</title>
                  <desc id="chart-desc">Andamento semestrale delle registrazioni degli utenti in piattaforma</desc>
                  <defs>
                    <linearGradient id="chartLineGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stop-color="#4F46E5" stop-opacity="0.25"/>
                      <stop offset="100%" stop-color="#4F46E5" stop-opacity="0.0"/>
                    </linearGradient>
                  </defs>
                  <line x1="40" y1="20" x2="480" y2="20" stroke="#f1f5f9" stroke-width="1" />
                  <line x1="40" y1="70" x2="480" y2="70" stroke="#f1f5f9" stroke-width="1" />
                  <line x1="40" y1="120" x2="480" y2="120" stroke="#f1f5f9" stroke-width="1" />
                  <line x1="40" y1="170" x2="480" y2="170" stroke="#cbd5e1" stroke-width="1.5" />
                  <text x="30" y="24" text-anchor="end" class="chart-label">{{ chartYMax }}</text>
                  <text x="30" y="74" text-anchor="end" class="chart-label">{{ Math.round(chartYMax * 0.67) }}</text>
                  <text x="30" y="124" text-anchor="end" class="chart-label">{{ Math.round(chartYMax * 0.33) }}</text>
                  <text x="30" y="174" text-anchor="end" class="chart-label">0</text>
                  <path :d="areaPath" fill="url(#chartLineGrad)" />
                  <path :d="linePath" fill="none" stroke="#4F46E5" stroke-width="3" stroke-linecap="round" />
                  <g v-for="(pt, idx) in chartPoints" :key="idx">
                    <circle
                      :cx="pt.x" :cy="pt.y" r="6"
                      :fill="hoveredIndex === idx ? '#4F46E5' : 'white'"
                      stroke="#4F46E5" stroke-width="3"
                      @mouseenter="hoveredIndex = idx"
                      @mouseleave="hoveredIndex = null"
                      style="cursor: pointer; transition: all 0.2s"
                    />
                    <text :x="pt.x" y="192" text-anchor="middle" class="chart-label text-weight-medium">{{ pt.label }}</text>
                  </g>
                </svg>
                <div
                  v-if="hoveredIndex !== null"
                  class="chart-tooltip shadow-md q-pa-sm rounded-lg bg-slate-800 text-white text-caption text-center"
                  :style="{ left: `${(chartPoints[hoveredIndex].x / 500) * 100}%`, top: `${(chartPoints[hoveredIndex].y / 200) * 100 - 20}%` }"
                >
                  <div class="text-weight-bold">{{ chartPoints[hoveredIndex].label }}</div>
                  <div>{{ chartPoints[hoveredIndex].val }} Utenti</div>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Role Breakdown Donut -->
        <div class="col-12 col-md-4">
          <q-card class="glass-card shadow-soft rounded-xl border border-slate-100 full-height column justify-between">
            <q-card-section class="q-pa-lg">
              <div class="text-h5 text-weight-bold text-outfit text-slate-800 q-mb-xs">Suddivisione Utenti</div>
              <div class="text-caption text-slate-500 q-mb-lg">Composizione della community per ruolo</div>

              <div v-if="donutSegments.length === 0 || stats?.total_users === 0" class="text-center text-slate-400 q-pa-md">
                <q-icon name="pie_chart" size="48px" class="q-mb-sm" />
                <div class="text-caption">Nessun dato disponibile</div>
              </div>

              <div v-else>
                <div class="row justify-center q-my-md relative-position">
                  <svg viewBox="0 0 200 200" class="full-width" style="max-height: 200px; transform: rotate(-90deg); overflow: visible;">
                    <circle
                      v-for="(seg, idx) in donutSegments"
                      :key="idx"
                      cx="100" cy="100" r="70"
                      fill="transparent"
                      :stroke="seg.color"
                      :stroke-width="hoveredDonut === idx ? '24' : '18'"
                      :stroke-dasharray="seg.dashArray"
                      :stroke-dashoffset="seg.dashOffset"
                      @mouseenter="hoveredDonut = idx"
                      @mouseleave="hoveredDonut = null"
                      style="cursor: pointer; transition: all 0.3s; transform-origin: center;"
                    />
                  </svg>
                  <div class="absolute-center text-center" style="transform: translate(-50%, -50%); pointer-events: none;">
                    <div class="text-h5 text-weight-bold text-slate-700 text-outfit q-mb-none">
                      {{ hoveredDonut !== null ? `${donutSegments[hoveredDonut].percent}%` : (stats?.total_users ?? 0) }}
                    </div>
                    <div class="text-caption text-slate-400 text-uppercase letter-spacing-1 text-weight-medium">
                      {{ hoveredDonut !== null ? donutSegments[hoveredDonut].label : 'Utenti' }}
                    </div>
                  </div>
                </div>
                <div class="q-gutter-y-xs q-mt-md">
                  <div
                    v-for="(seg, idx) in donutSegments" :key="idx"
                    class="row items-center justify-between q-py-xs q-px-sm rounded-lg hover-bg-grey cursor-pointer"
                    @mouseenter="hoveredDonut = idx"
                    @mouseleave="hoveredDonut = null"
                    :class="{ 'bg-slate-100': hoveredDonut === idx }"
                  >
                    <div class="row items-center q-gutter-x-sm">
                      <div class="rounded-circle" :style="{ width: '12px', height: '12px', backgroundColor: seg.color, borderRadius: '50%' }"></div>
                      <span class="text-weight-medium text-slate-700">{{ seg.label }}</span>
                    </div>
                    <div class="text-weight-bold text-slate-600">
                      {{ seg.value }} <span class="text-weight-regular text-caption text-slate-400">({{ seg.percent }}%)</span>
                    </div>
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Infrastructure Metrics from API -->
      <div class="row q-col-gutter-lg">
        <div class="col-12 col-md-6">
          <q-card class="glass-card shadow-soft rounded-xl border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="text-h5 text-weight-bold text-outfit text-slate-800 q-mb-md">Metriche di Infrastruttura</div>

              <div v-if="!systemMetrics" class="text-center text-slate-400 q-pa-md">
                <q-icon name="monitoring" size="48px" class="q-mb-sm" />
                <div class="text-caption">Metriche non disponibili — endpoint <code>/admin/system/metrics</code> richiesto</div>
              </div>

              <div v-else class="q-gutter-y-lg q-py-sm">
                <div>
                  <div class="row items-center justify-between text-caption text-slate-600 text-weight-medium q-mb-xs">
                    <span>Percentuale di Successo API</span>
                    <span :class="systemMetrics.api_success_rate >= 99 ? 'text-emerald-600' : 'text-orange-600'" class="text-weight-bold">
                      {{ systemMetrics.api_success_rate != null ? systemMetrics.api_success_rate + '%' : '-' }}
                    </span>
                  </div>
                  <q-linear-progress :value="(systemMetrics.api_success_rate ?? 0) / 100" color="emerald" track-color="slate-100" class="rounded-md" style="height: 8px;" />
                </div>
                <div>
                  <div class="row items-center justify-between text-caption text-slate-600 text-weight-medium q-mb-xs">
                    <span>Carico CPU Database</span>
                    <span :class="(systemMetrics.db_cpu_percent ?? 0) > 80 ? 'text-red-600' : 'text-indigo-600'" class="text-weight-bold">
                      {{ systemMetrics.db_cpu_percent != null ? systemMetrics.db_cpu_percent + '%' : '-' }}
                    </span>
                  </div>
                  <q-linear-progress :value="(systemMetrics.db_cpu_percent ?? 0) / 100" color="indigo" track-color="slate-100" class="rounded-md" style="height: 8px;" />
                </div>
                <div>
                  <div class="row items-center justify-between text-caption text-slate-600 text-weight-medium q-mb-xs">
                    <span>Cache Hit Rate (Redis)</span>
                    <span :class="(systemMetrics.cache_hit_rate ?? 0) >= 90 ? 'text-amber-600' : 'text-red-600'" class="text-weight-bold">
                      {{ systemMetrics.cache_hit_rate != null ? systemMetrics.cache_hit_rate + '%' : '-' }}
                    </span>
                  </div>
                  <q-linear-progress :value="(systemMetrics.cache_hit_rate ?? 0) / 100" color="amber" track-color="slate-100" class="rounded-md" style="height: 8px;" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Recent Audit Events from API -->
        <div class="col-12 col-md-6">
          <q-card class="glass-card shadow-soft rounded-xl border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="text-h5 text-weight-bold text-outfit text-slate-800 q-mb-md">Ultimi Eventi di Sistema</div>
              <div v-if="recentAuditEvents.length === 0" class="text-center text-slate-400 q-pa-md">
                <q-icon name="history" size="48px" class="q-mb-sm" />
                <div class="text-caption">Nessun evento recente</div>
              </div>
              <q-list v-else class="q-py-xs" separator>
                <q-item v-for="event in recentAuditEvents" :key="event.id" class="q-py-sm px-none">
                  <q-item-section avatar>
                    <q-avatar :color="getEventBgColor(event.type)" :text-color="getEventTextColor(event.type)" :icon="getEventIcon(event.type)" size="md" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-slate-700 text-caption">{{ event.description }}</q-item-label>
                    <q-item-label caption>{{ event.user_name || 'Sistema' }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <span class="text-caption text-grey">{{ formatDate(event.created_at) }}</span>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useQuasar } from 'quasar'
import { usePermissions } from '@/composables/usePermissions'
import adminService from '@/services/adminService'
import api from '@/services/api'

const $q = useQuasar()
const { isSuperAdmin } = usePermissions()

const loading = ref(false)
const stats = ref(null)
const systemMetrics = ref(null)
const userGrowthHistory = ref([])   // [{ label: 'Gen', value: 120 }, ...]
const recentAuditEvents = ref([])
const hoveredIndex = ref(null)
const hoveredDonut = ref(null)
const autoRefresh = ref(false)
let refreshInterval = null

// ── Fetch stats ──────────────────────────────────────────────
const fetchStats = async () => {
  try {
    const response = await adminService.getDashboardStats()
    stats.value = response.data
    recentAuditEvents.value = response.data?.recent_events || []
  } catch (e) {
    console.warn('Dashboard stats failed:', e)
    $q.notify({ type: 'warning', message: 'Impossibile caricare le statistiche generali' })
  }
}

// ── Fetch system metrics (latency, CPU, cache) ────────────────
const fetchSystemMetrics = async () => {
  try {
    const res = await api.get('/admin/system/metrics')
    systemMetrics.value = res.data
  } catch (e) {
    console.warn('System metrics endpoint not available:', e?.response?.status)
    systemMetrics.value = null
  }
}

// ── Fetch user growth history ─────────────────────────────────
const fetchUserGrowth = async () => {
  try {
    const res = await api.get('/admin/analytics/user-growth')
    userGrowthHistory.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
  } catch (e) {
    console.warn('User growth endpoint not available:', e?.response?.status)
    userGrowthHistory.value = []
  }
}

const fetchAll = async () => {
  loading.value = true
  try {
    await Promise.allSettled([fetchStats(), fetchSystemMetrics(), fetchUserGrowth()])
  } catch (err) {
    console.error('Error loading analytics:', err)
    $q.notify({ type: 'negative', message: 'Errore nel caricamento delle metriche' })
  } finally {
    loading.value = false
  }
}

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    refreshInterval = setInterval(fetchAll, 30000)
    $q.notify({ type: 'info', message: 'Aggiornamento automatico attivato (ogni 30s)', timeout: 2000 })
  } else {
    if (refreshInterval) clearInterval(refreshInterval)
    refreshInterval = null
    $q.notify({ type: 'info', message: 'Aggiornamento automatico disattivato', timeout: 2000 })
  }
}

onMounted(() => {
  fetchAll()
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})

// ── Chart computed ────────────────────────────────────────────
const chartYMax = computed(() => {
  if (!userGrowthHistory.value.length) return 0
  return Math.ceil(Math.max(...userGrowthHistory.value.map(d => d.value)) * 1.15)
})

const userGrowthTrend = computed(() => {
  const data = userGrowthHistory.value
  if (data.length < 2) return null
  const first = data[0].value
  const last = data[data.length - 1].value
  if (!first) return null
  return Math.round(((last - first) / first) * 100)
})

const chartPoints = computed(() => {
  if (!userGrowthHistory.value.length) return []
  const data = userGrowthHistory.value
  const width = 440
  const startX = 40
  const startY = 170
  const height = 150
  const maxVal = chartYMax.value || 1
  return data.map((d, idx) => ({
    x: startX + (idx * (width / Math.max(data.length - 1, 1))),
    y: startY - ((d.value / maxVal) * height),
    val: d.value,
    label: d.label
  }))
})

// Smooth cubic Bezier curve calculation
const smoothPath = (pts) => {
  if (!pts || !pts.length) return ''
  if (pts.length === 1) return `M ${pts[0].x} ${pts[0].y}`
  let d = `M ${pts[0].x.toFixed(1)} ${pts[0].y.toFixed(1)}`
  for (let i = 0; i < pts.length - 1; i++) {
    const p0 = pts[i === 0 ? i : i - 1]
    const p1 = pts[i]
    const p2 = pts[i + 1]
    const p3 = pts[i + 2 < pts.length ? i + 2 : i + 1]
    const cp1x = p1.x + (p2.x - p0.x) / 6
    const cp1y = p1.y + (p2.y - p0.y) / 6
    const cp2x = p2.x - (p3.x - p1.x) / 6
    const cp2y = p2.y - (p3.y - p1.y) / 6
    d += ` C ${cp1x.toFixed(1)},${cp1y.toFixed(1)} ${cp2x.toFixed(1)},${cp2y.toFixed(1)} ${p2.x.toFixed(1)},${p2.y.toFixed(1)}`
  }
  return d
}

const linePath = computed(() => {
  return smoothPath(chartPoints.value)
})

const areaPath = computed(() => {
  if (chartPoints.value.length < 2) return ''
  const first = chartPoints.value[0]
  const last = chartPoints.value[chartPoints.value.length - 1]
  const curve = smoothPath(chartPoints.value)
  return `${curve} L ${last.x.toFixed(1)},170 L ${first.x.toFixed(1)},170 Z`
})

const latencyColor = computed(() => {
  const ms = systemMetrics.value?.api_latency_ms
  if (ms == null) return 'text-slate-400'
  if (ms < 100) return 'text-emerald-600'
  if (ms < 500) return 'text-amber-600'
  return 'text-red-600'
})

// ── Donut chart ───────────────────────────────────────────────
const donutSegments = computed(() => {
  if (!stats.value) return []
  const students = stats.value.total_students ?? 0
  const teachers = stats.value.total_teachers ?? 0
  const parents  = stats.value.total_parents  ?? 0
  const total    = stats.value.total_users     ?? 0
  const staff    = Math.max(total - students - teachers - parents, 0)

  const data = [
    { label: 'Studenti', value: students, color: '#10B981' },
    { label: 'Docenti',  value: teachers, color: '#8B5CF6' },
    { label: 'Genitori', value: parents,  color: '#F59E0B' },
    { label: 'Staff',    value: staff,    color: '#3B82F6' }
  ].filter(d => d.value > 0)

  const totalVal = data.reduce((acc, d) => acc + d.value, 0)
  if (totalVal === 0) return []

  const radius = 70
  const circumference = 2 * Math.PI * radius
  let accumulated = 0

  return data.map(d => {
    const percent = d.value / totalVal
    const dashArray = `${percent * circumference} ${circumference}`
    const dashOffset = -accumulated * circumference
    accumulated += percent
    return { ...d, percent: Math.round(percent * 100), dashArray, dashOffset }
  })
})

// ── Audit event helpers ───────────────────────────────────────
const getEventIcon = (type) => {
  const icons = { create: 'add_circle', update: 'edit', delete: 'delete', login: 'login', backup: 'cloud_done', security: 'security' }
  return icons[type] || 'event'
}
const getEventBgColor   = (type) => ({ create: 'green-50', update: 'indigo-50', delete: 'red-50',    login: 'blue-50',   backup: 'green-50',  security: 'indigo-50' }[type] || 'grey-50')
const getEventTextColor = (type) => ({ create: 'green-700', update: 'indigo-700', delete: 'red-700', login: 'blue-700',  backup: 'green-700', security: 'indigo-700' }[type] || 'grey-700')

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const diff = new Date() - date
  if (diff < 3600000)  return `${Math.floor(diff / 60000)} min fa`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h fa`
  return date.toLocaleDateString('it-IT', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.text-gradient-premium {
  background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 50%, #EC4899 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.glass-card {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(12px);
}
.stat-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px -10px rgba(0, 0, 0, 0.08);
}
.chart-label {
  font-size: 11px;
  fill: #64748b;
  font-family: inherit;
}
.chart-tooltip {
  position: absolute;
  transform: translate(-50%, -100%);
  pointer-events: none;
  z-index: 10;
  min-width: 80px;
  transition: all 0.15s ease-out;
}
.chart-tooltip::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 50%;
  transform: translateX(-50%) rotate(45deg);
  width: 8px;
  height: 8px;
  background-color: #1e293b;
}
.hover-bg-grey:hover {
  background-color: #f8fafc;
}
.letter-spacing-1 { letter-spacing: 1px; }
</style>
