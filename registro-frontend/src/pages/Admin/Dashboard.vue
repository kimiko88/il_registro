<template>
  <q-page class="q-pa-md">
    <!-- Header -->
    <div class="row items-center q-mb-lg">
      <div class="col">
        <div class="text-h4 text-weight-bold">Dashboard Admin</div>
        <div class="text-subtitle1 text-grey-7">
          {{ isSuperAdmin ? 'Panoramica Globale' : 'La Mia Scuola' }}
        </div>
      </div>
      <div class="col-auto">
        <q-btn
          v-if="isSuperAdmin"
          color="primary"
          icon="refresh"
          label="Aggiorna"
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
    <div v-else-if="stats" class="row q-gutter-md q-mb-lg">
      <!-- Total Schools (SuperAdmin only) -->
      <div v-if="isSuperAdmin" class="col-12 col-md-3">
        <q-card class="stat-card">
          <q-card-section>
            <div class="row items-center">
              <div class="col">
                <div class="text-h6 text-weight-bold">{{ stats.total_schools }}</div>
                <div class="text-caption text-grey-7">Totale Scuole</div>
              </div>
              <div class="col-auto">
                <q-icon name="school" size="48px" color="indigo" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Users -->
      <div class="col-12 col-md-3">
        <q-card class="stat-card">
          <q-card-section>
            <div class="row items-center">
              <div class="col">
                <div class="text-h6 text-weight-bold">{{ stats.total_users }}</div>
                <div class="text-caption text-grey-7">Totale Utenti</div>
              </div>
              <div class="col-auto">
                <q-icon name="people" size="48px" color="cyan" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Students -->
      <div class="col-12 col-md-3">
        <q-card class="stat-card">
          <q-card-section>
            <div class="row items-center">
              <div class="col">
                <div class="text-h6 text-weight-bold">{{ stats.total_students }}</div>
                <div class="text-caption text-grey-7">Studenti</div>
              </div>
              <div class="col-auto">
                <q-icon name="school" size="48px" color="amber" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Total Teachers -->
      <div class="col-12 col-md-3">
        <q-card class="stat-card">
          <q-card-section>
            <div class="row items-center">
              <div class="col">
                <div class="text-h6 text-weight-bold">{{ stats.total_teachers }}</div>
                <div class="text-caption text-grey-7">Docenti</div>
              </div>
              <div class="col-auto">
                <q-icon name="person" size="48px" color="purple" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <div v-if="stats" class="row q-col-gutter-md">
      <!-- Recent Events -->
      <div class="col-12 col-md-8">
        <q-card>
          <q-card-section>
            <div class="text-h6 q-mb-md">Eventi Recenti</div>
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
        <q-card>
          <q-card-section>
            <div class="text-h6 q-mb-md">Stato Sistema</div>
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
  router.push('/settings')
}

onMounted(() => {
  fetchDashboardStats()
})
</script>

<style scoped>
.stat-card {
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.health-item {
  padding: 8px 0;
}
</style>
