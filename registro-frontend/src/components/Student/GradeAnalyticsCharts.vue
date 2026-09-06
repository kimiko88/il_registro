<template>
  <q-card class="glass-card shadow-soft rounded-2xl overflow-hidden grade-analytics-card">
    <!-- Header -->
    <q-card-section class="row items-center justify-between q-pb-none">
      <div>
        <div class="text-h6 text-weight-bold text-slate-800 row items-center">
          <q-icon name="insights" color="primary" class="q-mr-sm" size="26px" />
          <span>{{ title || $t('roleDashboards.analyticsTitle') || 'Analisi Rendimento & Trend' }}</span>
        </div>
        <div class="text-caption text-slate-500">
          {{ subtitle || $t('roleDashboards.analyticsSubtitle') || 'Evoluzione cronologica dei voti e regolarità presenze' }}
        </div>
      </div>
      <div v-if="trendStats.count > 0" class="row items-center q-gutter-xs">
        <q-badge
          :color="trendDelta >= 0 ? 'positive' : 'negative'"
          class="q-px-sm q-py-xs text-caption text-weight-bold rounded-lg"
        >
          <q-icon :name="trendDelta >= 0 ? 'trending_up' : 'trending_down'" class="q-mr-xs" />
          {{ trendDelta >= 0 ? `+${trendDelta.toFixed(1)}` : trendDelta.toFixed(1) }}
        </q-badge>
        <span class="text-caption text-slate-500">{{ $t('roleDashboards.trendLabel') || 'Trend Recente' }}</span>
      </div>
    </q-card-section>

    <!-- Main Analytics Body -->
    <q-card-section class="q-pt-md">
      <div class="row q-col-gutter-lg items-center">
        <!-- SVG Sparkline Curve Chart -->
        <div class="col-12 col-md-8">
          <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-xs row items-center justify-between">
            <span>{{ $t('roleDashboards.gradeEvolution') || 'Andamento Cronologico Voti' }}</span>
            <span class="text-caption text-slate-400">Min 2.0 • Sufficienza 6.0 • Max 10.0</span>
          </div>

          <div class="svg-chart-container relative-position">
            <svg
              v-if="chartPoints.length >= 2"
              viewBox="0 0 500 160"
              class="w-full h-auto sparkline-svg"
              preserveAspectRatio="none"
              role="img"
              :aria-label="$t('roleDashboards.sparklineAria') || 'Grafico andamento cronologico voti'"
            >
              <defs>
                <!-- Area gradient -->
                <linearGradient id="curveGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#4f46e5" stop-opacity="0.35" />
                  <stop offset="60%" stop-color="#6366f1" stop-opacity="0.12" />
                  <stop offset="100%" stop-color="#a5b4fc" stop-opacity="0.0" />
                </linearGradient>

                <!-- Passing threshold line dash pattern -->
                <linearGradient id="strokeGradient" x1="0" y1="0" x2="1" y2="0">
                  <stop offset="0%" stop-color="#6366f1" />
                  <stop offset="100%" stop-color="#3b82f6" />
                </linearGradient>
              </defs>

              <!-- Background Grid Lines -->
              <!-- Max 10 -->
              <line x1="35" y1="20" x2="480" y2="20" stroke="#e2e8f0" stroke-width="1" stroke-dasharray="2,2" />
              <text x="25" y="24" font-size="10" fill="#94a3b8" text-anchor="end">10</text>

              <!-- Passing 6.0 guideline -->
              <line x1="35" y1="75" x2="480" y2="75" stroke="#10b981" stroke-width="1.5" stroke-dasharray="4,4" stroke-opacity="0.75" />
              <text x="25" y="78" font-size="10" fill="#10b981" font-weight="bold" text-anchor="end">6.0</text>

              <!-- Min 2 -->
              <line x1="35" y1="130" x2="480" y2="130" stroke="#e2e8f0" stroke-width="1" stroke-dasharray="2,2" />
              <text x="25" y="134" font-size="10" fill="#94a3b8" text-anchor="end">2</text>

              <!-- Area Fill Under Curve -->
              <path :d="areaPath" fill="url(#curveGradient)" />

              <!-- Main Trend Curve -->
              <path
                :d="linePath"
                fill="none"
                stroke="url(#strokeGradient)"
                stroke-width="3"
                stroke-linecap="round"
                stroke-linejoin="round"
              />

              <!-- Interactive Data Points -->
              <g v-for="(pt, idx) in chartPoints" :key="idx">
                <circle
                  :cx="pt.x"
                  :cy="pt.y"
                  r="6"
                  :fill="pt.grade >= 6 ? '#10b981' : '#ef4444'"
                  stroke="#ffffff"
                  stroke-width="2.5"
                  class="chart-dot cursor-pointer"
                  @mouseenter="hoveredPoint = pt"
                  @mouseleave="hoveredPoint = null"
                />
              </g>

              <!-- Native SVG Tooltip overlay on hover -->
              <g v-if="hoveredPoint" :transform="`translate(${hoveredPoint.x}, ${Math.max(30, hoveredPoint.y - 12)})`">
                <rect
                  :x="hoveredPoint.x > 420 ? -110 : (hoveredPoint.x < 80 ? 0 : -55)"
                  y="-34"
                  width="110"
                  height="30"
                  rx="6"
                  fill="#0f172a"
                  fill-opacity="0.9"
                />
                <text
                  :x="hoveredPoint.x > 420 ? -55 : (hoveredPoint.x < 80 ? 55 : 0)"
                  y="-20"
                  text-anchor="middle"
                  fill="#ffffff"
                  font-size="11"
                  font-weight="bold"
                >
                  {{ hoveredPoint.grade }} • {{ hoveredPoint.subject }}
                </text>
                <text
                  :x="hoveredPoint.x > 420 ? -55 : (hoveredPoint.x < 80 ? 55 : 0)"
                  y="-9"
                  text-anchor="middle"
                  fill="#94a3b8"
                  font-size="9"
                >
                  {{ hoveredPoint.formattedDate }}
                </text>
              </g>
            </svg>

            <!-- Fallback when < 2 data points -->
            <div v-else class="flex flex-center q-pa-lg text-slate-400 bg-slate-50 rounded-xl" style="height: 160px">
              <q-icon name="query_stats" size="36px" class="q-mr-sm opacity-60" />
              <span>{{ $t('roleDashboards.insufficientData') || 'Inserisci almeno due voti per visualizzare il trend' }}</span>
            </div>
          </div>
        </div>

        <!-- Attendance Radial Ring & Summary -->
        <div class="col-12 col-md-4">
          <div class="p-4 bg-slate-50/80 rounded-2xl border border-slate-100 flex flex-center column">
            <div class="text-caption text-weight-bold text-slate-600 text-uppercase letter-spacing-1 q-mb-sm">
              {{ $t('roleDashboards.presenceRate') || 'Tasso di Presenza' }}
            </div>

            <!-- Native SVG Radial Gauge -->
            <div class="relative-position flex flex-center" style="width: 100px; height: 100px">
              <svg viewBox="0 0 100 100" width="100" height="100">
                <!-- Track -->
                <circle
                  cx="50"
                  cy="50"
                  r="40"
                  fill="none"
                  stroke="#e2e8f0"
                  stroke-width="8"
                />
                <!-- Progress ring -->
                <circle
                  cx="50"
                  cy="50"
                  r="40"
                  fill="none"
                  :stroke="radialColor"
                  stroke-width="8"
                  stroke-linecap="round"
                  :stroke-dasharray="circumference"
                  :stroke-dashoffset="radialOffset"
                  transform="rotate(-90 50 50)"
                  style="transition: stroke-dashoffset 0.8s ease-in-out;"
                />
              </svg>
              <div class="absolute text-center">
                <div class="text-h6 text-weight-bolder text-slate-800 line-height-tight">{{ attendanceRate }}%</div>
                <div class="text-caption text-slate-400" style="font-size: 9px">Validità</div>
              </div>
            </div>

            <!-- Threshold indicator -->
            <div class="row items-center justify-between full-width q-mt-sm text-caption">
              <span class="text-slate-500">Min. MIUR (75%)</span>
              <q-badge :color="attendanceRate >= 75 ? 'positive' : 'negative'" class="q-px-xs">
                {{ attendanceRate >= 75 ? 'Regolare' : 'A Rischio' }}
              </q-badge>
            </div>
          </div>
        </div>
      </div>

      <!-- Distribution Stacked Bar -->
      <div class="q-mt-lg">
        <div class="row items-center justify-between text-caption text-slate-600 q-mb-xs">
          <span class="text-weight-bold">{{ $t('roleDashboards.gradeDistribution') || 'Distribuzione Voti' }}</span>
          <span>{{ distribution.total }} {{ $t('roleDashboards.totalGradesEvaluated') || 'valutazioni' }}</span>
        </div>

        <!-- Segmented bar -->
        <div class="distribution-bar row no-wrap overflow-hidden rounded-borders bg-slate-100" style="height: 16px">
          <div
            v-if="distribution.failPct > 0"
            :style="{ width: `${distribution.failPct}%` }"
            class="bg-rose-500 text-white flex flex-center text-caption font-bold"
            :title="`Insufficienti (<6): ${distribution.failCount} (${distribution.failPct}%)`"
          >
            <span v-if="distribution.failPct >= 12" style="font-size: 10px">{{ distribution.failPct }}%</span>
          </div>
          <div
            v-if="distribution.passPct > 0"
            :style="{ width: `${distribution.passPct}%` }"
            class="bg-amber-500 text-white flex flex-center text-caption font-bold"
            :title="`Sufficienti (6-7.5): ${distribution.passCount} (${distribution.passPct}%)`"
          >
            <span v-if="distribution.passPct >= 12" style="font-size: 10px">{{ distribution.passPct }}%</span>
          </div>
          <div
            v-if="distribution.goodPct > 0"
            :style="{ width: `${distribution.goodPct}%` }"
            class="bg-emerald-500 text-white flex flex-center text-caption font-bold"
            :title="`Ottimi (8-10): ${distribution.goodCount} (${distribution.goodPct}%)`"
          >
            <span v-if="distribution.goodPct >= 12" style="font-size: 10px">{{ distribution.goodPct }}%</span>
          </div>
        </div>

        <!-- Distribution Legend Pills -->
        <div class="row q-col-gutter-sm q-mt-xs">
          <div class="col-4">
            <div class="p-2 rounded-lg bg-rose-50 border border-rose-100 row items-center justify-between">
              <div class="row items-center">
                <span class="legend-dot bg-rose-500 q-mr-xs"></span>
                <span class="text-caption text-rose-900 text-weight-medium">&lt; 6 (Insuff.)</span>
              </div>
              <span class="text-weight-bold text-rose-800 text-caption">{{ distribution.failCount }}</span>
            </div>
          </div>
          <div class="col-4">
            <div class="p-2 rounded-lg bg-amber-50 border border-amber-100 row items-center justify-between">
              <div class="row items-center">
                <span class="legend-dot bg-amber-500 q-mr-xs"></span>
                <span class="text-caption text-amber-900 text-weight-medium">6 - 7.5 (Suff.)</span>
              </div>
              <span class="text-weight-bold text-amber-800 text-caption">{{ distribution.passCount }}</span>
            </div>
          </div>
          <div class="col-4">
            <div class="p-2 rounded-lg bg-emerald-50 border border-emerald-100 row items-center justify-between">
              <div class="row items-center">
                <span class="legend-dot bg-emerald-500 q-mr-xs"></span>
                <span class="text-caption text-emerald-900 text-weight-medium">8 - 10 (Ottimo)</span>
              </div>
              <span class="text-weight-bold text-emerald-800 text-caption">{{ distribution.goodCount }}</span>
            </div>
          </div>
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  grades: {
    type: Array,
    default: () => []
  },
  attendanceRate: {
    type: Number,
    default: 100
  },
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  }
})

const hoveredPoint = ref(null)

// Normalize grades list: filter out absences or non-numeric values
const validGrades = computed(() => {
  if (!Array.isArray(props.grades)) return []
  return props.grades
    .map(g => {
      const val = g.grade_value ?? g.value ?? g.grade
      const numVal = parseFloat(val)
      return {
        ...g,
        numGrade: isNaN(numVal) ? null : numVal,
        dateObj: g.date ? new Date(g.date) : new Date()
      }
    })
    .filter(g => g.numGrade !== null && g.numGrade > 0 && g.numGrade <= 10)
    .sort((a, b) => a.dateObj - b.dateObj)
})

// Trend stats: compare last 3 grades average vs preceding grades
const trendStats = computed(() => {
  const list = validGrades.value
  const count = list.length
  if (count === 0) return { count: 0, delta: 0, average: 0 }

  const totalSum = list.reduce((acc, g) => acc + g.numGrade, 0)
  const average = totalSum / count

  if (count < 4) {
    const delta = list[count - 1].numGrade - list[0].numGrade
    return { count, delta, average }
  }

  const recent = list.slice(-3)
  const previous = list.slice(0, -3)

  const recentAvg = recent.reduce((acc, g) => acc + g.numGrade, 0) / recent.length
  const prevAvg = previous.reduce((acc, g) => acc + g.numGrade, 0) / previous.length

  return {
    count,
    delta: recentAvg - prevAvg,
    average
  }
})

const trendDelta = computed(() => trendStats.value.delta)

// SVG Sparkline Math
// Coordinate space: viewBox="0 0 500 160"
// Plot bounds: x: 35 -> 480 (width: 445), y: 20 (grade 10) -> 130 (grade 2) (height: 110)
const chartPoints = computed(() => {
  const list = validGrades.value
  if (list.length === 0) return []

  const padLeft = 35
  const chartWidth = 445
  const padTop = 20
  const chartHeight = 110
  const n = list.length

  return list.map((g, i) => {
    const x = n === 1 ? padLeft + chartWidth / 2 : padLeft + (i / (n - 1)) * chartWidth
    // Clamp grade between 2 and 10 for coordinates
    const clampedGrade = Math.max(2, Math.min(10, g.numGrade))
    const y = padTop + chartHeight * ((10 - clampedGrade) / 8)

    return {
      x: Math.round(x * 10) / 10,
      y: Math.round(y * 10) / 10,
      grade: g.numGrade,
      subject: g.subject || g.subject_name || 'Valutazione',
      formattedDate: g.dateObj.toLocaleDateString('it-IT', { day: '2-digit', month: 'short' })
    }
  })
})

// Bezier smoothing path generator
const linePath = computed(() => {
  const pts = chartPoints.value
  if (pts.length === 0) return ''
  if (pts.length === 1) return `M ${pts[0].x} ${pts[0].y}`

  let d = `M ${pts[0].x} ${pts[0].y}`
  for (let i = 0; i < pts.length - 1; i++) {
    const p0 = pts[i]
    const p1 = pts[i + 1]
    const cp1x = Math.round((p0.x + (p1.x - p0.x) / 2) * 10) / 10
    const cp1y = p0.y
    const cp2x = cp1x
    const cp2y = p1.y
    d += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${p1.x} ${p1.y}`
  }
  return d
})

// Gradient Area Path under curve
const areaPath = computed(() => {
  const pts = chartPoints.value
  if (pts.length < 2) return ''
  const bottomY = 145
  const last = pts[pts.length - 1]
  const first = pts[0]
  return `${linePath.value} L ${last.x} ${bottomY} L ${first.x} ${bottomY} Z`
})

// Attendance Radial Gauge (radius = 40)
const circumference = 2 * Math.PI * 40 // ~251.327
const radialOffset = computed(() => {
  const rate = Math.max(0, Math.min(100, props.attendanceRate))
  return circumference * (1 - rate / 100)
})

const radialColor = computed(() => {
  if (props.attendanceRate >= 90) return '#10b981' // emerald
  if (props.attendanceRate >= 75) return '#f59e0b' // amber
  return '#ef4444' // rose
})

// Grade distribution breakdown (<6, 6-7.5, 8-10)
const distribution = computed(() => {
  const list = validGrades.value
  const total = list.length
  if (total === 0) {
    return { total: 0, failCount: 0, failPct: 0, passCount: 0, passPct: 0, goodCount: 0, goodPct: 0 }
  }

  let failCount = 0
  let passCount = 0
  let goodCount = 0

  list.forEach(g => {
    if (g.numGrade < 6) failCount++
    else if (g.numGrade <= 7.5) passCount++
    else goodCount++
  })

  return {
    total,
    failCount,
    failPct: Math.round((failCount / total) * 100),
    passCount,
    passPct: Math.round((passCount / total) * 100),
    goodCount,
    goodPct: Math.round((goodCount / total) * 100)
  }
})
</script>

<style scoped>
.sparkline-svg {
  filter: drop-shadow(0 4px 6px -1px rgba(99, 102, 241, 0.1));
}

.chart-dot {
  transition: r 0.2s cubic-bezier(0.4, 0, 0.2, 1), stroke-width 0.2s ease;
}

.chart-dot:hover {
  r: 8;
  stroke-width: 3.5;
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 9999px;
  display: inline-block;
}

.letter-spacing-1 {
  letter-spacing: 0.05em;
}

.line-height-tight {
  line-height: 1.1;
}
</style>
