<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h3 text-weight-bold text-outfit q-my-none text-gradient-premium">
          Dashboard Segreteria
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">Gestione amministrativa e scolastica</div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div v-for="(stat, index) in statsCards" :key="index" class="col-12 col-sm-6 col-md-3">
        <q-card class="shadow-1 hover-shadow transition-swing bg-white">
          <q-card-section class="row items-center no-wrap">
            <q-avatar :color="stat.color" text-color="white" :icon="stat.icon" font-size="24px" rounded size="lg" />
            <div class="q-ml-md">
              <div class="text-subtitle2 text-grey-7 text-uppercase" style="letter-spacing: 1px; font-size: 0.75rem">{{ stat.label }}</div>
              <div class="text-h4 text-weight-bolder text-dark">{{ stat.value }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Upcoming Tasks / Review Queue -->
      <div class="col-12 col-md-8">
        <q-card class="shadow-1">
           <q-card-section class="row items-center justify-between">
             <div class="text-h6 text-weight-bold">Da Revisionare</div>
             <q-btn flat round icon="refresh" color="grey" @click="fetchData" />
           </q-card-section>
           
           <q-list separator v-if="pendingReviews.length > 0">
             <q-item v-for="item in pendingReviews" :key="item.id" clickable v-ripple @click="router.push(`/secretary/documents?id=${item.id}`)">
               <q-item-section avatar>
                 <q-icon name="description" color="primary" />
               </q-item-section>
               <q-item-section>
                 <q-item-label class="text-weight-bold">{{ item.title }}</q-item-label>
                 <q-item-label caption>{{ item.author }} - {{ item.date }}</q-item-label>
               </q-item-section>
               <q-item-section side>
                 <q-chip size="sm" color="orange-1" text-color="orange-9" label="In Revisione" />
               </q-item-section>
             </q-item>
           </q-list>
           <div v-else class="q-pa-lg text-center text-grey">
             <q-icon name="check_circle" size="40px" color="positive" class="q-mb-sm" />
             <div>Nessun documento in attesa di revisione</div>
           </div>
        </q-card>

        <q-card class="shadow-1 q-mt-md">
             <q-card-section>
                  <div class="text-h6 text-weight-bold">Attività Recenti</div>
             </q-card-section>
              <q-list v-if="recentEvents.length > 0">
                  <q-item v-for="event in recentEvents" :key="event.id">
                      <q-item-section avatar>
                          <q-icon :name="getEventIcon(event.type)" :color="getEventColor(event.type)" />
                      </q-item-section>
                      <q-item-section>
                          <q-item-label>{{ event.description }}</q-item-label>
                          <q-item-label caption>{{ event.user_name }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                          <span class="text-grey-6 text-caption">{{ formatDate(event.created_at) }}</span>
                      </q-item-section>
                  </q-item>
              </q-list>
              <div v-else class="q-pa-md text-center text-grey">Nessuna attività recente</div>
        </q-card>
      </div>

      <!-- Quick Actions -->
      <div class="col-12 col-md-4">
        <q-card class="shadow-1 q-mb-md">
          <q-card-section>
            <div class="text-h6 text-weight-bold q-mb-md">Azioni Rapide</div>
            <div class="q-gutter-y-sm">
              <q-btn 
                color="primary" 
                class="full-width" 
                size="lg" 
                padding="md"
                no-caps 
                align="left"
                icon="campaign" 
                label="Nuova Circolare" 
                to="/secretary/communications" 
              />
              <q-btn 
                color="secondary" 
                class="full-width" 
                size="lg" 
                padding="md"
                no-caps 
                align="left"
                icon="person_add" 
                label="Registra Utente/Studente" 
                to="/secretary/users" 
              />
              <q-btn 
                color="accent" 
                class="full-width" 
                size="lg" 
                padding="md"
                no-caps 
                align="left"
                icon="upload_file" 
                label="Carica Documento" 
                to="/secretary/documents" 
              />
              <q-btn 
                outline
                color="grey-8" 
                class="full-width text-weight-medium" 
                size="lg" 
                padding="md"
                no-caps 
                align="left"
                icon="settings" 
                label="Impostazioni Scuola" 
                to="/secretary/settings" 
              />
            </div>
          </q-card-section>
        </q-card>
        
        <!-- Notifications small -->
        <q-card class="bg-indigo-1 text-indigo-9 shadow-soft">
             <q-card-section>
                 <div class="text-weight-bold q-mb-sm row items-center">
                    <q-icon name="campaign" class="q-mr-xs" />
                    Avvisi Importanti
                 </div>
                 <ul class="q-pl-md q-mb-none" v-if="announcements.length > 0">
                     <li v-for="ann in announcements" :key="ann.id" class="q-mb-xs">
                        {{ ann.title }}
                     </li>
                 </ul>
                 <div v-else class="text-caption text-indigo-7">Nessun avviso recente</div>
             </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router'
import adminService from 'src/services/adminService'
import documentService from 'src/services/documentService'

const $q = useQuasar()
const router = useRouter()

const statsCards = ref([
    { label: 'Documenti Pendenti', value: '0', icon: 'pending_actions', color: 'orange' },
    { label: 'Circolari', value: '0', icon: 'campaign', color: 'blue' },
    { label: 'Totale Studenti', value: '0', icon: 'school', color: 'green' },
    { label: 'Totale Docenti', value: '0', icon: 'work', color: 'purple' }
])

const pendingReviews = ref([])
const recentEvents = ref([])
const announcements = ref([])
const loading = ref(false)

const fetchData = async () => {
    loading.value = true
    try {
        // Fetch Stats
        const statsRes = await adminService.getDashboardStats()
        if (statsRes.data) {
            statsCards.value[0].value = statsRes.data.pending_documents_count || '0'
            statsCards.value[1].value = statsRes.data.announcements_count || '0'
            statsCards.value[2].value = statsRes.data.total_students || '0'
            statsCards.value[3].value = statsRes.data.total_teachers || '0'
            recentEvents.value = statsRes.data.recent_events || []
        }

        // Fetch Pending Reviews
        const docRes = await documentService.getInbox({ status: 'pending' })
        if (docRes.data && docRes.data.items) {
            pendingReviews.value = docRes.data.items.slice(0, 5).map(d => ({
                id: d.id,
                title: d.title,
                author: d.author || 'Docente',
                date: new Date(d.created_at).toLocaleDateString('it-IT')
            }))
        }

        // Fetch Communications for Announcements
        const { useCommunicationsStore } = await import('src/stores/communications')
        const commStore = useCommunicationsStore()
        await commStore.fetchCommunications()
        announcements.value = commStore.communications.slice(0, 3)
    } catch (e) {
        console.error("Error fetching secretary dashboard data", e)
    } finally {
        loading.value = false
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
    const diff = new Date() - date
    if (diff < 3600000) return `${Math.floor(diff / 60000)} min fa`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h fa`
    return date.toLocaleDateString('it-IT', { day: '2-digit', month: 'short' })
}

onMounted(() => {
    fetchData()
})
</script>

<style scoped>
.hover-shadow:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}
.transition-swing {
    transition: transform 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
}
</style>
