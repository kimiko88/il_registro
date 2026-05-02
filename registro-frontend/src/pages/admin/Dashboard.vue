<template>
  <q-page class="q-pa-md">
    <!-- Header -->
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h3 text-weight-bold text-outfit q-my-none bg-clip-text text-transparent bg-gradient-premium" style="display: inline-block;">
          Dashboard Admin
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">
          {{ isSuperAdmin ? 'Benvenuto nel pannello di controllo globale' : 'Gestione centralizzata della tua scuola' }}
        </div>
      </div>
      <div class="col-auto">
        <q-btn
          v-if="isSuperAdmin"
          unelevated
          color="white"
          text-color="primary"
          icon="refresh"
          label="Sincronizza Dati"
          class="rounded-lg shadow-soft q-px-md"
          @click="fetchDashboardStats"
          :loading="loading"
        />
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !stats" class="row q-gutter-md">
      <div v-for="i in 4" :key="i" class="col-12 col-md-3">
        <q-skeleton height="120px" />
      </div>
    </div>

    <!-- Stats Cards -->
    <div v-else-if="stats" class="row q-col-gutter-lg q-mb-xl">
      <!-- Total Schools (SuperAdmin only) -->
      <div v-if="isSuperAdmin" class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card full-height">
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="col">
                <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Scuole</div>
                <div class="text-h4 text-weight-bold text-outfit">{{ stats.total_schools }}</div>
              </div>
              <div class="col-auto">
                <div class="bg-indigo-100 q-pa-md rounded-xl">
                    <q-icon name="school" size="32px" color="indigo-700" />
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Users -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card full-height">
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="col">
                <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Utenti</div>
                <div class="text-h4 text-weight-bold text-outfit">{{ stats.total_users }}</div>
              </div>
              <div class="col-auto">
                <div class="bg-cyan-100 q-pa-md rounded-xl">
                    <q-icon name="people" size="32px" color="cyan-700" />
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Students -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card full-height">
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="col">
                <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Studenti</div>
                <div class="text-h4 text-weight-bold text-outfit">{{ stats.total_students }}</div>
              </div>
              <div class="col-auto">
                <div class="bg-amber-100 q-pa-md rounded-xl">
                    <q-icon name="face" size="32px" color="amber-700" />
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Teachers -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card full-height">
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="col">
                <div class="text-caption text-slate-500 text-uppercase letter-spacing-1 q-mb-xs">Docenti</div>
                <div class="text-h4 text-weight-bold text-outfit">{{ stats.total_teachers }}</div>
              </div>
              <div class="col-auto">
                <div class="bg-purple-100 q-pa-md rounded-xl">
                    <q-icon name="supervisor_account" size="32px" color="purple-700" />
                </div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <div v-if="stats" class="row q-col-gutter-xl">
      <!-- Recent Events -->
      <div class="col-12 col-md-8">
        <q-card class="glass-card shadow-soft overflow-hidden">
          <q-card-section class="q-pa-lg">
            <div class="row items-center justify-between q-mb-lg">
                <div class="text-h5 text-weight-bold text-outfit">Attività Recenti</div>
                <q-btn flat color="primary" label="Vedi Audit Log" to="/admin/audit-logs" no-caps v-if="isSuperAdmin" />
            </div>
            <q-list v-if="stats.recent_events && stats.recent_events.length > 0" separator>
              <q-item v-for="event in stats.recent_events" :key="event.id">
                <q-item-section avatar>
                  <q-avatar :color="getEventColor(event.type)" text-color="white">
                    <q-icon :name="getEventIcon(event.type)" />
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ event.description }}</q-item-label>
                  <q-item-label caption>
                    {{ event.user_name }}
                    <span v-if="event.school_name"> • {{ event.school_name }}</span>
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-item-label caption>{{ formatDate(event.created_at) }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
            <div v-else class="text-center q-pa-md text-grey-6">
              Nessun evento recente
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Health Status (SuperAdmin only) -->
      <div v-if="isSuperAdmin && stats.health_status" class="col-12 col-md-4">
        <q-card class="glass-card shadow-soft q-mb-lg">
          <q-card-section class="q-pa-lg">
            <div class="text-h5 text-weight-bold text-outfit q-mb-lg">Monitoraggio</div>
            <div class="q-gutter-sm">
              <div class="health-item">
                <div class="row items-center">
                  <div class="col">Database</div>
                  <div class="col-auto">
                    <q-badge
                      :color="getHealthColor(stats.health_status.database.status)"
                      :label="stats.health_status.database.status"
                    />
                  </div>
                </div>
              </div>
              <div class="health-item">
                <div class="row items-center">
                  <div class="col">Storage</div>
                  <div class="col-auto">
                    <q-badge
                      :color="getHealthColor(stats.health_status.storage.status)"
                      :label="stats.health_status.storage.status"
                    />
                  </div>
                </div>
              </div>
              <div class="health-item">
                <div class="row items-center">
                  <div class="col">API</div>
                  <div class="col-auto">
                    <q-badge
                      :color="getHealthColor(stats.health_status.api.status)"
                      :label="stats.health_status.api.status"
                    />
                  </div>
                </div>
              </div>
              <q-separator class="q-my-md" />
              <div class="text-center">
                <q-badge
                  :color="getHealthColor(stats.health_status.overall_status)"
                  :label="`Stato: ${stats.health_status.overall_status}`"
                  size="lg"
                />
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Active Users Card -->
        <q-card class="q-mt-md">
          <q-card-section>
            <div class="text-h6 q-mb-sm">Utenti Attivi (24h)</div>
            <div class="text-h4 text-center text-primary q-my-md">
              {{ stats.active_users_24h }}
            </div>
            <div class="text-caption text-center text-grey-7">
              Utenti attivi nelle ultime 24 ore
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Quick Actions (Admin only) -->
      <div v-if="!isSuperAdmin" class="col-12 col-md-4">
        <q-card>
          <q-card-section>
            <div class="text-h6 q-mb-md">Azioni Rapide</div>
            <q-list>
              <q-item clickable v-ripple @click="goToSchool">
                <q-item-section avatar>
                  <q-icon name="school" color="primary" />
                </q-item-section>
                <q-item-section>Gestisci Scuola</q-item-section>
              </q-item>
              <q-item clickable v-ripple @click="goToAnalytics">
                <q-item-section avatar>
                  <q-icon name="analytics" color="primary" />
                </q-item-section>
                <q-item-section>Visualizza Analytics</q-item-section>
              </q-item>
              <q-item clickable v-ripple @click="goToSettings">
                <q-item-section avatar>
                  <q-icon name="settings" color="primary" />
                </q-item-section>
                <q-item-section>Impostazioni</q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="row">
      <div class="col-12">
        <q-banner class="bg-negative text-white">
          <template v-slot:avatar>
            <q-icon name="error" />
          </template>
          Errore nel caricamento dei dati: {{ error }}
          <template v-slot:action>
            <q-btn flat label="Riprova" @click="fetchDashboardStats" />
          </template>
        </q-banner>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { usePermissions } from '@/composables/usePermissions'
import adminService from '@/services/adminService'

const router = useRouter()
const $q = useQuasar()
const { isSuperAdmin } = usePermissions()

const loading = ref(false)
const stats = ref(null)
const error = ref(null)

const fetchDashboardStats = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await adminService.getDashboardStats()
    stats.value = response.data
  } catch (err) {
    error.value = err.response?.data?.message || err.message
    $q.notify({
      type: 'negative',
      message: 'Errore nel caricamento della dashboard',
      caption: error.value
    })
  } finally {
    loading.value = false
  }
}

const getEventColor = (type) => {
  const colors = {
    create: 'positive',
    update: 'info',
    delete: 'negative',
    login: 'primary'
  }
  return colors[type] || 'grey'
}

const getEventIcon = (type) => {
  const icons = {
    create: 'add_circle',
    update: 'edit',
    delete: 'delete',
    login: 'login'
  }
  return icons[type] || 'event'
}

const getHealthColor = (status) => {
  const colors = {
    healthy: 'positive',
    warning: 'warning',
    error: 'negative',
    unknown: 'grey'
  }
  return colors[status] || 'grey'
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date
  
  // Less than 1 hour
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes} min fa`
  }
  
  // Less than 1 day
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000)
    return `${hours}h fa`
  }
  
  // Format as date
  return date.toLocaleDateString('it-IT', {
    day: '2-digit',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const goToSchool = () => {
  router.push('/admin/schools')
}

const goToAnalytics = () => {
  router.push('/admin/analytics')
}

const goToSettings = () => {
  router.push('/admin/settings')
}

onMounted(() => {
  fetchDashboardStats()
})
</script>

<style scoped>
.bg-clip-text {
    -webkit-background-clip: text;
    background-clip: text;
}

.bg-gradient-premium {
    background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
}

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

.health-item {
  padding: 12px 0;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}
.health-item:last-child {
    border-bottom: none;
}
</style>
