<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="text-h4 text-weight-bold text-dark q-mb-md">Dashboard Segreteria</div>

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
             <q-item v-for="item in pendingReviews" :key="item.id" clickable v-ripple @click="$router.push(`/secretary/documents?id=${item.id}`)">
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
             <q-list>
                 <q-item>
                     <q-item-section avatar>
                         <q-icon name="person_add" color="secondary" />
                     </q-item-section>
                     <q-item-section>
                         <q-item-label>Nuovo studente iscritto</q-item-label>
                         <q-item-label caption>Mario Rossi (1A)</q-item-label>
                     </q-item-section>
                     <q-item-section side>
                         <span class="text-grey-6 text-caption">10 min fa</span>
                     </q-item-section>
                 </q-item>
                 <q-item>
                     <q-item-section avatar>
                         <q-icon name="campaign" color="primary" />
                     </q-item-section>
                     <q-item-section>
                         <q-item-label>Circolare inviata</q-item-label>
                         <q-item-label caption>Sciopero Docenti</q-item-label>
                     </q-item-section>
                     <q-item-section side>
                         <span class="text-grey-6 text-caption">1 ora fa</span>
                     </q-item-section>
                 </q-item>
             </q-list>
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
        <q-card class="bg-indigo-1 text-indigo-9">
             <q-card-section>
                 <div class="text-weight-bold q-mb-sm">Avvisi Importanti</div>
                 <ul class="q-pl-md q-mb-none">
                     <li>Scadenza iscrizioni: 20 Gennaio</li>
                     <li>Consegna PDP entro Venerdì</li>
                 </ul>
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

const $q = useQuasar()
const router = useRouter()

const statsCards = ref([
    { label: 'Documenti Pendenti', value: '12', icon: 'pending_actions', color: 'orange' },
    { label: 'Richieste Utenti', value: '5', icon: 'person_search', color: 'blue' },
    { label: 'Totale Studenti', value: '1,240', icon: 'school', color: 'green' },
    { label: 'Totale Docenti', value: '128', icon: 'work', color: 'purple' }
])

const pendingReviews = ref([
    { id: 1, title: 'PDP - Giulia Verdi', author: 'Prof. Bianchi', date: 'Oggi' },
    { id: 2, title: 'Certificato Medico - Rossi', author: 'Segreteria', date: 'Ieri' }
])

const fetchData = async () => {
    // Mock fetch
    $q.loading.show()
    setTimeout(() => {
        $q.loading.hide()
        // Here we would call stores to get real data
        $q.notify({ type: 'positive', message: 'Dati aggiornati' })
    }, 500)
}

onMounted(() => {
    // fetchData()
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
