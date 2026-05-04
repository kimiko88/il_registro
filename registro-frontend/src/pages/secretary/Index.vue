<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h3 text-weight-bold text-outfit q-my-none text-gradient-premium">
          Dashboard Segreteria
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">Benvenuto, ecco il riepilogo della gestione scolastica.</div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="row q-col-gutter-lg q-mb-xl">
      <div v-for="(stat, index) in statsCards" :key="index" class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card shadow-soft hover-card border-slate-100 overflow-hidden">
          <q-card-section class="row items-center no-wrap">
            <q-avatar 
              :color="stat.color + '-50'" 
              :text-color="stat.color + '-700'" 
              :icon="stat.icon" 
              font-size="28px" 
              rounded 
              size="56px" 
            />
            <div class="q-ml-lg">
              <div class="text-subtitle2 text-slate-500 text-uppercase letter-spacing-1">{{ stat.label }}</div>
              <div class="text-h3 text-weight-bold text-slate-800">{{ stat.value }}</div>
            </div>
          </q-card-section>
          <div :class="`bg-${stat.color}-500`" style="height: 4px; opacity: 0.8"></div>
        </q-card>
      </div>
    </div>

    <div class="row q-col-gutter-xl">
      <!-- Left Column: Tasks & Events -->
      <div class="col-12 col-md-8">
        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white overflow-hidden q-mb-lg">
           <q-card-section class="row items-center justify-between q-pa-lg">
             <div class="text-h5 text-weight-bold text-slate-800">Da Revisionare</div>
             <q-btn flat round icon="refresh" color="primary" @click="fetchData" />
           </q-card-section>
           
           <q-separator color="slate-100" />
           
           <q-list separator v-if="pendingReviews.length > 0">
             <q-item v-for="item in pendingReviews" :key="item.id" clickable v-ripple @click="router.push(`/secretary/documents?id=${item.id}`)" class="q-pa-lg">
               <q-item-section avatar>
                 <q-avatar color="orange-50" text-color="orange-700" icon="description" size="44px" />
               </q-item-section>
               <q-item-section>
                 <q-item-label class="text-subtitle1 text-weight-bold text-slate-800">{{ item.title }}</q-item-label>
                 <q-item-label caption class="text-slate-500">{{ item.author }} • {{ item.date }}</q-item-label>
               </q-item-section>
               <q-item-section side>
                 <q-chip size="sm" color="orange-100" text-color="orange-800" label="In Revisione" class="text-weight-bold rounded-lg" />
               </q-item-section>
             </q-item>
           </q-list>
           <div v-else class="q-pa-xl text-center text-slate-400">
             <q-icon name="check_circle" size="64px" color="emerald-400" class="opacity-40 q-mb-md" />
             <div class="text-h6">Ottimo lavoro!</div>
             <div>Nessun documento in attesa di revisione</div>
           </div>
        </q-card>

        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white overflow-hidden">
             <q-card-section class="q-pa-lg">
                  <div class="text-h5 text-weight-bold text-slate-800">Attività Recenti</div>
             </q-card-section>
             <q-separator color="slate-100" />
             <q-list v-if="recentEvents.length > 0">
                 <q-item v-for="event in recentEvents" :key="event.id" class="q-pa-md">
                     <q-item-section avatar>
                         <q-avatar :color="getEventColor(event.type) + '-50'" :text-color="getEventColor(event.type) + '-700'" :icon="getEventIcon(event.type)" size="40px" />
                     </q-item-section>
                     <q-item-section>
                         <q-item-label class="text-slate-800 text-weight-medium">{{ event.description }}</q-item-label>
                         <q-item-label caption class="text-slate-400">{{ event.user_name }}</q-item-label>
                     </q-item-section>
                     <q-item-section side>
                         <span class="text-slate-400 text-caption">{{ formatDate(event.created_at) }}</span>
                     </q-item-section>
                 </q-item>
             </q-list>
             <div v-else class="q-pa-xl text-center text-slate-400">Nessuna attività recente</div>
        </q-card>
      </div>

      <!-- Right Column: Quick Actions & Notifications -->
      <div class="col-12 col-md-4">
        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white overflow-hidden q-mb-lg">
          <q-card-section class="q-pa-lg">
            <div class="text-h5 text-weight-bold text-slate-800 q-mb-lg">Azioni Rapide</div>
            <div class="column q-gutter-y-md">
              <q-btn 
                color="indigo" 
                class="rounded-xl text-weight-bold shadow-sm" 
                size="lg" 
                padding="md"
                no-caps 
                unelevated
                align="left"
                icon="campaign" 
                label="Nuova Circolare" 
                to="/secretary/communications" 
              />
              <q-btn 
                color="emerald" 
                class="rounded-xl text-weight-bold shadow-sm bg-emerald-600" 
                size="lg" 
                padding="md"
                no-caps 
                unelevated
                align="left"
                icon="person_add" 
                label="Registra Utente" 
                to="/secretary/users" 
              />
              <q-btn 
                color="amber" 
                class="rounded-xl text-weight-bold shadow-sm bg-amber-600" 
                size="lg" 
                padding="md"
                no-caps 
                unelevated
                align="left"
                icon="upload_file" 
                label="Carica Documento" 
                to="/secretary/documents" 
              />
              <q-btn 
                color="slate-200"
                text-color="slate-800"
                class="rounded-xl text-weight-bold border-slate-200" 
                size="lg" 
                padding="md"
                no-caps 
                unelevated
                align="left"
                icon="settings" 
                label="Impostazioni" 
                to="/secretary/settings" 
              />
            </div>
          </q-card-section>
        </q-card>
        
        <!-- Notifications Premium -->
        <q-card class="bg-gradient-premium text-white rounded-xl shadow-soft overflow-hidden">
             <q-card-section class="q-pa-lg">
                  <div class="text-h6 text-weight-bold q-mb-md row items-center no-wrap">
                    <q-icon name="auto_awesome" class="q-mr-sm" />
                    Avvisi di Sistema
                  </div>
                  <div class="column q-gutter-y-md" v-if="announcements.length > 0">
                      <div v-for="ann in announcements" :key="ann.id" class="glass-effect q-pa-md rounded-lg border-white-10 cursor-pointer hover-scale">
                        <div class="text-weight-bold ellipsis">{{ ann.title }}</div>
                        <div class="text-caption opacity-80">{{ ann.date || 'Recente' }}</div>
                      </div>
                  </div>
                  <div v-else class="text-caption opacity-80 italic">Nessun avviso importante al momento.</div>
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
    { label: 'Pendenti', value: '0', icon: 'pending_actions', color: 'orange' },
    { label: 'Circolari', value: '0', icon: 'campaign', color: 'blue' },
    { label: 'Studenti', value: '0', icon: 'school', color: 'green' },
    { label: 'Docenti', value: '0', icon: 'work', color: 'purple' }
])

const pendingReviews = ref([])
const recentEvents = ref([])
const announcements = ref([])
const loading = ref(false)

const fetchData = async () => {
    loading.value = true
    try {
        const statsRes = await adminService.getDashboardStats()
        if (statsRes.data) {
            statsCards.value[0].value = statsRes.data.pending_documents_count || '0'
            statsCards.value[1].value = statsRes.data.announcements_count || '0'
            statsCards.value[2].value = statsRes.data.total_students || '0'
            statsCards.value[3].value = statsRes.data.total_teachers || '0'
            recentEvents.value = statsRes.data.recent_events || []
        }

        const docRes = await documentService.getInbox({ status: 'pending' })
        if (docRes.data && docRes.data.items) {
            pendingReviews.value = docRes.data.items.slice(0, 5).map(d => ({
                id: d.id,
                title: d.title,
                author: d.author || 'Docente',
                date: new Date(d.created_at).toLocaleDateString('it-IT')
            }))
        }

        const { useCommunicationsStore } = await import('src/stores/communications')
        const commStore = useCommunicationsStore()
        await commStore.fetchCommunications()
        announcements.value = commStore.communications.slice(0, 3)
    } catch (e) {
        console.error("Error fetching dashboard data", e)
    } finally {
        loading.value = false
    }
}

const getEventIcon = (type) => {
    const icons = { create: 'add_circle', update: 'edit', delete: 'delete', login: 'login' }
    return icons[type] || 'event'
}

const getEventColor = (type) => {
    const colors = { create: 'emerald', update: 'indigo', delete: 'negative', login: 'primary' }
    return colors[type] || 'slate'
}

const formatDate = (dateString) => {
    const date = new Date(dateString)
    const diff = new Date() - date
    if (diff < 3600000) return `${Math.floor(diff / 60000)} min fa`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h fa`
    return date.toLocaleDateString('it-IT', { day: '2-digit', month: 'short' })
}

onMounted(fetchData)
</script>

<style scoped>
.letter-spacing-1 { letter-spacing: 0.05em; }
.opacity-40 { opacity: 0.4; }
.opacity-80 { opacity: 0.8; }
.border-white-10 { border: 1px solid rgba(255, 255, 255, 0.1); }

.hover-card {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.hover-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04) !important;
}

.hover-scale {
    transition: transform 0.2s ease;
}
.hover-scale:hover {
    transform: scale(1.02);
}

.bg-emerald-600 { background-color: #059669; }
.bg-amber-600 { background-color: #d97706; }
.text-emerald-700 { color: #047857; }
.bg-emerald-50 { background-color: #ecfdf5; }

.body--dark .bg-emerald-600 { background-color: #10b981; }
.body--dark .bg-amber-600 { background-color: #fbbf24; }
.body--dark .bg-emerald-600 { background-color: #10b981; }
.body--dark .bg-amber-600 { background-color: #fbbf24; }
</style>
