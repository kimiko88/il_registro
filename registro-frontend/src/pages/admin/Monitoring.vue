<template>
  <q-page class="q-pa-lg bg-slate-50">
    <div class="row items-center justify-between q-mb-xl">
      <div>
        <h1 class="text-h3 text-weight-bold text-outfit q-my-none text-gradient-premium">Monitoraggio Sistema</h1>
        <p class="text-subtitle1 text-slate-500 q-mt-sm q-mb-none">Stato in tempo reale dell'infrastruttura</p>
      </div>
      <q-btn unelevated color="white" text-color="primary" icon="refresh" label="Aggiorna" class="rounded-lg shadow-soft q-px-md" @click="fetchHealth" :loading="loading" />
    </div>

    <!-- Loading -->
    <div v-if="loading && !health" class="row q-col-gutter-lg">
      <div v-for="i in 4" :key="i" class="col-12 col-md-3">
        <q-skeleton height="120px" class="rounded-xl" />
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="fetchError" class="text-center q-pa-xl">
      <q-icon name="wifi_off" size="64px" color="red-4" class="q-mb-md" />
      <div class="text-h6 text-slate-600">Impossibile raggiungere l'endpoint di monitoraggio</div>
      <div class="text-caption text-slate-400 q-mt-sm">Endpoint richiesto: <code>GET /admin/system/health</code></div>
      <q-btn class="q-mt-lg" color="primary" label="Riprova" @click="fetchHealth" unelevated />
    </div>

    <div v-else-if="health">
      <!-- Status Banner -->
      <q-banner
        :class="overallStatus === 'healthy' ? 'bg-emerald-50 text-emerald-800' : overallStatus === 'degraded' ? 'bg-amber-50 text-amber-800' : 'bg-red-50 text-red-800'"
        class="rounded-xl q-mb-xl shadow-soft"
      >
        <template v-slot:avatar>
          <q-icon :name="overallStatus === 'healthy' ? 'check_circle' : overallStatus === 'degraded' ? 'warning' : 'error'" size="32px" />
        </template>
        <div class="text-h6 text-weight-bold">
          {{ overallStatus === 'healthy' ? 'Tutti i servizi operativi' : overallStatus === 'degraded' ? 'Servizi degradati' : 'Servizi critici offline' }}
        </div>
        <div class="text-caption">Ultimo aggiornamento: {{ lastChecked }}</div>
      </q-banner>

      <!-- Service Cards -->
      <div class="row q-col-gutter-lg q-mb-xl">
        <div v-for="svc in services" :key="svc.name" class="col-12 col-sm-6 col-md-3">
          <q-card class="glass-card shadow-soft rounded-xl overflow-hidden">
            <q-card-section class="row items-center no-wrap q-pa-lg">
              <div :class="`bg-${svc.color}-100 q-pa-md rounded-xl q-mr-md`">
                <q-icon :name="svc.icon" :color="svc.color + '-700'" size="28px" />
              </div>
              <div>
                <div class="text-caption text-slate-500 text-uppercase letter-spacing-1">{{ svc.label }}</div>
                <q-chip
                  :color="svc.status === 'ok' ? 'positive' : svc.status === 'degraded' ? 'warning' : 'negative'"
                  text-color="white"
                  size="sm"
                  class="text-weight-bold q-mt-xs"
                >
                  {{ svc.status === 'ok' ? 'ONLINE' : svc.status === 'degraded' ? 'DEGRADATO' : 'OFFLINE' }}
                </q-chip>
              </div>
            </q-card-section>
            <div v-if="svc.detail" class="q-px-lg q-pb-md text-caption text-slate-500">{{ svc.detail }}</div>
          </q-card>
        </div>
      </div>

      <!-- Detailed Metrics -->
      <div class="row q-col-gutter-lg">
        <div class="col-12 col-md-6">
          <q-card class="glass-card shadow-soft rounded-xl">
            <q-card-section class="q-pa-lg">
              <div class="text-h6 text-weight-bold text-slate-800 q-mb-md">Risorse Server</div>
              <div class="q-gutter-y-md">
                <div v-for="metric in resourceMetrics" :key="metric.label">
                  <div class="row items-center justify-between text-caption text-slate-600 text-weight-medium q-mb-xs">
                    <span>{{ metric.label }}</span>
                    <span class="text-weight-bold" :class="metric.valueClass">{{ metric.display }}</span>
                  </div>
                  <q-linear-progress :value="metric.ratio" :color="metric.color" track-color="slate-100" class="rounded-md" style="height: 8px;" />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <div class="col-12 col-md-6">
          <q-card class="glass-card shadow-soft rounded-xl">
            <q-card-section class="q-pa-lg">
              <div class="text-h6 text-weight-bold text-slate-800 q-mb-md">Info Versione</div>
              <q-list dense>
                <q-item v-for="info in versionInfo" :key="info.label">
                  <q-item-section>
                    <q-item-label caption class="text-slate-500">{{ info.label }}</q-item-label>
                    <q-item-label class="text-weight-medium text-slate-800">{{ info.value }}</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Empty state if endpoint not yet implemented -->
    <div v-else class="text-center q-pa-xl">
      <q-icon name="sensors" size="64px" color="grey-4" class="q-mb-md" />
      <div class="text-h6 text-slate-600">Nessun dato disponibile</div>
      <div class="text-caption text-slate-400 q-mt-sm">Implementa <code>GET /admin/system/health</code> nel backend</div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import api from '@/services/api'

const loading  = ref(false)
const health   = ref(null)
const fetchError = ref(false)
const lastChecked = ref('-')
let pollInterval = null

const fetchHealth = async () => {
  loading.value = true
  fetchError.value = false
  try {
    const res = await api.get('/admin/system/health')
    health.value = res.data
    lastChecked.value = new Date().toLocaleTimeString('it-IT')
  } catch (e) {
    console.error('Health check failed:', e)
    fetchError.value = true
    health.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchHealth()
  // auto-refresh ogni 60s
  pollInterval = setInterval(fetchHealth, 60000)
})
onUnmounted(() => clearInterval(pollInterval))

const overallStatus = computed(() => health.value?.status || 'unknown')

const services = computed(() => {
  if (!health.value) return []
  const s = health.value.services || {}
  return [
    { name: 'api',      label: 'API Server',    icon: 'cloud',          color: 'indigo',  status: s.api      || 'unknown', detail: health.value.api_version ? `v${health.value.api_version}` : null },
    { name: 'db',       label: 'Database',      icon: 'storage',        color: 'blue',    status: s.database || 'unknown', detail: s.db_ping_ms != null ? `Ping: ${s.db_ping_ms}ms` : null },
    { name: 'redis',    label: 'Cache Redis',   icon: 'memory',         color: 'red',     status: s.redis    || 'unknown', detail: s.redis_ping_ms != null ? `Ping: ${s.redis_ping_ms}ms` : null },
    { name: 'storage',  label: 'Storage',       icon: 'folder',         color: 'amber',   status: s.storage  || 'unknown', detail: null }
  ]
})

const resourceMetrics = computed(() => {
  const m = health.value?.metrics || {}
  return [
    {
      label: 'CPU',
      display: m.cpu_percent != null ? m.cpu_percent + '%' : '-',
      ratio: (m.cpu_percent ?? 0) / 100,
      color: (m.cpu_percent ?? 0) > 80 ? 'negative' : 'indigo',
      valueClass: (m.cpu_percent ?? 0) > 80 ? 'text-red-600' : 'text-indigo-600'
    },
    {
      label: 'RAM',
      display: m.memory_percent != null ? m.memory_percent + '%' : '-',
      ratio: (m.memory_percent ?? 0) / 100,
      color: (m.memory_percent ?? 0) > 85 ? 'negative' : 'teal',
      valueClass: (m.memory_percent ?? 0) > 85 ? 'text-red-600' : 'text-teal-600'
    },
    {
      label: 'Disco',
      display: m.disk_percent != null ? m.disk_percent + '%' : '-',
      ratio: (m.disk_percent ?? 0) / 100,
      color: (m.disk_percent ?? 0) > 90 ? 'negative' : 'amber',
      valueClass: (m.disk_percent ?? 0) > 90 ? 'text-red-600' : 'text-amber-600'
    },
    {
      label: 'Latenza API',
      display: m.api_latency_ms != null ? m.api_latency_ms + 'ms' : '-',
      ratio: Math.min((m.api_latency_ms ?? 0) / 1000, 1),
      color: (m.api_latency_ms ?? 0) > 500 ? 'negative' : 'emerald',
      valueClass: (m.api_latency_ms ?? 0) > 500 ? 'text-red-600' : 'text-emerald-600'
    }
  ]
})

const versionInfo = computed(() => {
  const h = health.value || {}
  return [
    { label: 'Versione API',      value: h.api_version      || '-' },
    { label: 'Versione DB',       value: h.db_version       || '-' },
    { label: 'Ambiente',          value: h.environment      || '-' },
    { label: 'Uptime',            value: h.uptime           || '-' },
    { label: 'Ultimo Deploy',     value: h.last_deploy      ? new Date(h.last_deploy).toLocaleString('it-IT') : '-' }
  ]
})
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
.letter-spacing-1 { letter-spacing: 1px; }
</style>
