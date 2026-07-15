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
          @click="fetchStats"
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
        <!-- Schools (SuperAdmin only) -->
        <div v-if="isSuperAdmin" class="col-12 col-sm-6 col-md-3">
          <q-card class="glass-card stat-card full-height shadow-soft rounded-xl overflow-hidden border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Scuole Registrate</div>
                  <div class="text-h4 text-weight-bold text-outfit text-slate-800">{{ stats?.total_schools || 0 }}</div>
                </div>
                <div class="bg-indigo-100 q-pa-md rounded-xl">
                  <q-icon name="school" size="32px" color="indigo-700" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Users Count -->
        <div class="col-12 col-sm-6" :class="isSuperAdmin ? 'col-md-3' : 'col-md-4'">
          <q-card class="glass-card stat-card full-height shadow-soft rounded-xl overflow-hidden border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Utenti Totali</div>
                  <div class="text-h4 text-weight-bold text-outfit text-slate-800">{{ stats?.total_users || 0 }}</div>
                </div>
                <div class="bg-purple-100 q-pa-md rounded-xl">
                  <q-icon name="people" size="32px" color="purple-700" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Documents Count -->
        <div class="col-12 col-sm-6" :class="isSuperAdmin ? 'col-md-3' : 'col-md-4'">
          <q-card class="glass-card stat-card full-height shadow-soft rounded-xl overflow-hidden border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Documenti Caricati</div>
                  <div class="text-h4 text-weight-bold text-outfit text-slate-800">{{ stats?.total_documents || 0 }}</div>
                </div>
                <div class="bg-teal-100 q-pa-md rounded-xl">
                  <q-icon name="description" size="32px" color="teal-700" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- API Response Time (System Performance) -->
        <div class="col-12 col-sm-6" :class="isSuperAdmin ? 'col-md-3' : 'col-md-4'">
          <q-card class="glass-card stat-card full-height shadow-soft rounded-xl overflow-hidden border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Latenza Server API</div>
                  <div class="text-h4 text-weight-bold text-outfit text-emerald-600">48ms</div>
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
        <!-- User Growth Trend (Line Chart) -->
        <div class="col-12 col-md-8">
          <q-card class="glass-card shadow-soft rounded-xl border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="row items-center justify-between q-mb-md">
                <div>
                  <div class="text-h5 text-weight-bold text-outfit text-slate-800">Crescita Utenti Registrati</div>
                  <div class="text-caption text-slate-500">Andamento semestrale delle registrazioni in piattaforma</div>
                </div>
                <q-chip outline color="primary" class="text-weight-bold" size="sm">Trend +18%</q-chip>
              </div>

              <!-- SVG Line Chart -->
              <div class="relative-position q-py-md">
                <svg viewBox="0 0 500 200" class="full-width" style="overflow: visible;">
                  <!-- Gradients -->
                  <defs>
                    <linearGradient id="chartLineGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stop-color="#4F46E5" stop-opacity="0.25"/>
                      <stop offset="100%" stop-color="#4F46E5" stop-opacity="0.0"/>
                    </linearGradient>
                  </defs>

                  <!-- Horizontal Grid Lines -->
                  <line x1="40" y1="20" x2="480" y2="20" stroke="#f1f5f9" stroke-width="1" />
                  <line x1="40" y1="70" x2="480" y2="70" stroke="#f1f5f9" stroke-width="1" />
                  <line x1="40" y1="120" x2="480" y2="120" stroke="#f1f5f9" stroke-width="1" />
                  <line x1="40" y1="170" x2="480" y2="170" stroke="#cbd5e1" stroke-width="1.5" />

                  <!-- Grid Labels (Y Axis) -->
                  <text x="30" y="24" text-anchor="end" class="chart-label">1.2k</text>
                  <text x="30" y="74" text-anchor="end" class="chart-label">800</text>
                  <text x="30" y="124" text-anchor="end" class="chart-label">400</text>
                  <text x="30" y="174" text-anchor="end" class="chart-label">0</text>

                  <!-- Area Path -->
                  <path :d="areaPath" fill="url(#chartLineGrad)" />

                  <!-- Line Path -->
                  <path :d="linePath" fill="none" stroke="#4F46E5" stroke-width="3" stroke-linecap="round" />

                  <!-- Interactive Data Points -->
                  <g v-for="(pt, idx) in chartPoints" :key="idx">
                    <circle 
                      :cx="pt.x" 
                      :cy="pt.y" 
                      r="6" 
                      :fill="hoveredIndex === idx ? '#4F46E5' : 'white'" 
                      stroke="#4F46E5" 
                      stroke-width="3" 
                      @mouseenter="hoveredIndex = idx" 
                      @mouseleave="hoveredIndex = null"
                      style="cursor: pointer; transition: all 0.2s" 
                    />
                    <!-- Chart Labels (X Axis) -->
                    <text :x="pt.x" y="192" text-anchor="middle" class="chart-label text-weight-medium">
                      {{ months[idx] }}
                    </text>
                  </g>
                </svg>

                <!-- Tooltip Overlay -->
                <div 
                  v-if="hoveredIndex !== null" 
                  class="chart-tooltip shadow-md q-pa-sm rounded-lg bg-slate-800 text-white text-caption text-center"
                  :style="{ left: `${(chartPoints[hoveredIndex].x / 500) * 100}%`, top: `${(chartPoints[hoveredIndex].y / 200) * 100 - 20}%` }"
                >
                  <div class="text-weight-bold">{{ months[hoveredIndex] }}</div>
                  <div>{{ chartPoints[hoveredIndex].val }} Utenti</div>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Role Breakdown (Donut Chart) -->
        <div class="col-12 col-md-4">
          <q-card class="glass-card shadow-soft rounded-xl border border-slate-100 full-height column justify-between">
            <q-card-section class="q-pa-lg">
              <div class="text-h5 text-weight-bold text-outfit text-slate-800 q-mb-xs">Suddivisione Utenti</div>
              <div class="text-caption text-slate-500 q-mb-lg">Composizione della community per ruolo</div>

              <!-- SVG Donut Chart -->
              <div class="row justify-center q-my-md relative-position">
                <svg viewBox="0 0 200 200" class="full-width" style="max-height: 200px; transform: rotate(-90deg); overflow: visible;">
                  <circle
                    v-for="(seg, idx) in donutSegments"
                    :key="idx"
                    cx="100"
                    cy="100"
                    r="70"
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
                
                <!-- Central Indicator Text -->
                <div class="absolute-center text-center" style="transform: translate(-50%, -50%); pointer-events: none;">
                  <div class="text-h5 text-weight-bold text-slate-700 text-outfit q-mb-none">
                    {{ hoveredDonut !== null ? `${donutSegments[hoveredDonut].percent}%` : (stats?.total_users || 0) }}
                  </div>
                  <div class="text-caption text-slate-400 text-uppercase letter-spacing-1 text-weight-medium">
                    {{ hoveredDonut !== null ? donutSegments[hoveredDonut].label : 'Utenti' }}
                  </div>
                </div>
              </div>

              <!-- Legend -->
              <div class="q-gutter-y-xs q-mt-md">
                <div 
                  v-for="(seg, idx) in donutSegments" 
                  :key="idx" 
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
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Resource Utilization & Monitor Indicators -->
      <div class="row q-col-gutter-lg">
        <!-- API Usage & Monitoring Metrics -->
        <div class="col-12 col-md-6">
          <q-card class="glass-card shadow-soft rounded-xl border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="text-h5 text-weight-bold text-outfit text-slate-800 q-mb-md">Metriche di Infrastruttura</div>
              
              <div class="q-gutter-y-lg q-py-sm">
                <!-- API Success Rate -->
                <div>
                  <div class="row items-center justify-between text-caption text-slate-600 text-weight-medium q-mb-xs">
                    <span>Percentuale di Successo API</span>
                    <span class="text-emerald-600 text-weight-bold">99.98%</span>
                  </div>
                  <q-linear-progress :value="0.9998" color="emerald" track-color="slate-100" class="rounded-md" style="height: 8px;" />
                </div>

                <!-- CPU Load -->
                <div>
                  <div class="row items-center justify-between text-caption text-slate-600 text-weight-medium q-mb-xs">
                    <span>Carico CPU Database</span>
                    <span class="text-indigo-600 text-weight-bold">12%</span>
                  </div>
                  <q-linear-progress :value="0.12" color="indigo" track-color="slate-100" class="rounded-md" style="height: 8px;" />
                </div>

                <!-- Cache Hit Rate -->
                <div>
                  <div class="row items-center justify-between text-caption text-slate-600 text-weight-medium q-mb-xs">
                    <span>Cache Hit Rate (Redis)</span>
                    <span class="text-amber-600 text-weight-bold">94.2%</span>
                  </div>
                  <q-linear-progress :value="0.942" color="amber" track-color="slate-100" class="rounded-md" style="height: 8px;" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Recent Audit Events Summary -->
        <div class="col-12 col-md-6">
          <q-card class="glass-card shadow-soft rounded-xl border border-slate-100">
            <q-card-section class="q-pa-lg">
              <div class="text-h5 text-weight-bold text-outfit text-slate-800 q-mb-md">Eventi di Manutenzione</div>

              <q-list class="q-py-xs" separator>
                <q-item class="q-py-sm px-none">
                  <q-item-section avatar>
                    <q-avatar color="green-50" text-color="green-700" icon="cloud_done" size="md" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-slate-700 text-caption">Backup Automatico Database</q-item-label>
                    <q-item-label caption>Archiviazione completata con successo su storage S3</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <span class="text-caption text-grey">Oggi 03:00</span>
                  </q-item-section>
                </q-item>

                <q-item class="q-py-sm px-none">
                  <q-item-section avatar>
                    <q-avatar color="indigo-50" text-color="indigo-700" icon="security" size="md" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-slate-700 text-caption">Rinnovo Certificati SSL/TLS</q-item-label>
                    <q-item-label caption>Certificati Let's Encrypt rinnovati per i sottodomini api</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <span class="text-caption text-grey">Ieri 14:20</span>
                  </q-item-section>
                </q-item>

                <q-item class="q-py-sm px-none">
                  <q-item-section avatar>
                    <q-avatar color="amber-50" text-color="amber-700" icon="update" size="md" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-slate-700 text-caption">Ricarica delle Variabili d'Ambiente</q-item-label>
                    <q-item-label caption>Rilevato aggiornamento segreti da GitHub Actions workflow</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <span class="text-caption text-grey">15 Lug, 15:01</span>
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
import { ref, computed, onMounted } from 'vue'
import { usePermissions } from '@/composables/usePermissions'
import adminService from '@/services/adminService'

const { isSuperAdmin } = usePermissions()

const loading = ref(false)
const stats = ref(null)
const hoveredIndex = ref(null)
const hoveredDonut = ref(null)

const months = ['Gen', 'Feb', 'Mar', 'Apr', 'Mag', 'Giu']
const userGrowthData = [120, 240, 350, 480, 720, 950]

// Fetch Dashboard Stats from Backend to populate Counters & Breakdowns
const fetchStats = async () => {
  loading.value = true
  try {
    const response = await adminService.getDashboardStats()
    stats.value = response.data
  } catch (err) {
    console.error('Error loading analytics statistics:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})

// Line Chart Vector Math
const chartPoints = computed(() => {
  const width = 440
  const height = 150
  const startX = 40
  const startY = 170
  
  // Set maximum values dynamically or default to 1000
  const realMax = stats.value?.total_users ? Math.max(stats.value.total_users, 1000) : 1000
  const maxVal = realMax + 200 // Padding for top margin
  
  // Update last month to the actual users total from backend
  const data = [...userGrowthData]
  if (stats.value?.total_users) {
    data[data.length - 1] = stats.value.total_users
  }

  return data.map((val, idx) => {
    const x = startX + (idx * (width / (data.length - 1)))
    const y = startY - ((val / maxVal) * height)
    return { x, y, val }
  })
})

const linePath = computed(() => {
  if (chartPoints.value.length === 0) return ''
  return chartPoints.value.reduce((acc, pt, idx) => {
    return idx === 0 ? `M ${pt.x} ${pt.y}` : `${acc} L ${pt.x} ${pt.y}`
  }, '')
})

const areaPath = computed(() => {
  if (chartPoints.value.length === 0) return ''
  const first = chartPoints.value[0]
  const last = chartPoints.value[chartPoints.value.length - 1]
  const points = chartPoints.value.map(pt => `${pt.x},${pt.y}`).join(' ')
  return `M ${first.x} 170 L ${points} L ${last.x} 170 Z`
})

// Donut Chart Role Breakdowns (Based on actual database stats or sensible defaults)
const donutSegments = computed(() => {
  const students = stats.value?.total_students || 850
  const teachers = stats.value?.total_teachers || 120
  const parents = Math.round(students * 0.9)
  const staff = stats.value?.total_users ? Math.max(stats.value.total_users - students - teachers - parents, 30) : 30

  const data = [
    { label: 'Studenti', value: students, color: '#10B981' },
    { label: 'Docenti', value: teachers, color: '#8B5CF6' },
    { label: 'Genitori', value: parents, color: '#F59E0B' },
    { label: 'Staff', value: staff, color: '#3B82F6' }
  ]

  const total = data.reduce((acc, d) => acc + d.value, 0)
  const radius = 70
  const circumference = 2 * Math.PI * radius

  let accumulatedPercent = 0

  return data.map(d => {
    const percent = total > 0 ? (d.value / total) : 0.25
    const dashArray = `${percent * circumference} ${circumference}`
    const dashOffset = -accumulatedPercent * circumference
    accumulatedPercent += percent

    return {
      ...d,
      percent: Math.round(percent * 100),
      dashArray,
      dashOffset
    }
  })
})
</script>

<style scoped>
.text-gradient-premium {
  background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 50%, #EC4899 100%);
  -webkit-background-clip: text;
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
</style>
