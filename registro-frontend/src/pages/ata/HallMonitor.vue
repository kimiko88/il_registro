<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          Tabellone Pass Uscita Aule (Personale ATA)
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          Monitoraggio continuo corridoi e servizi igienici di plesso in tempo reale
        </p>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-btn
          color="primary"
          icon="refresh"
          label="Aggiorna Ora"
          unelevated
          class="rounded-lg"
          @click="loadActivePasses"
        />
      </div>
    </div>

    <!-- Overview Stats Cards -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card shadow-soft border-slate-100 rounded-xl">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">Studenti Fuori Aula</div>
              <div class="text-h4 text-weight-bolder text-primary">{{ activePasses.length }}</div>
            </div>
            <q-avatar color="primary" text-color="white" icon="directions_walk" size="lg" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card shadow-soft border-slate-100 rounded-xl">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">Alert Superamento 15m</div>
              <div class="text-h4 text-weight-bolder text-negative">{{ overdueCount }}</div>
            </div>
            <q-avatar color="negative" text-color="white" icon="warning" size="lg" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card shadow-soft border-slate-100 rounded-xl">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">Nei Servizi Igienici</div>
              <div class="text-h4 text-weight-bolder text-indigo-7">{{ bathroomCount }}</div>
            </div>
            <q-avatar color="indigo-7" text-color="white" icon="wc" size="lg" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card shadow-soft border-slate-100 rounded-xl">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-slate-500">In Infermeria / Uffici</div>
              <div class="text-h4 text-weight-bolder text-warning">{{ otherCount }}</div>
            </div>
            <q-avatar color="warning" text-color="white" icon="local_hospital" size="lg" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Overdue Alert Banner if any student >15min -->
    <div v-if="overdueCount > 0" class="q-mb-lg">
      <q-banner dense rounded class="bg-red-600 text-white shadow-md">
        <template #avatar>
          <q-icon name="warning" size="md" color="white" />
        </template>
        <div class="text-weight-bolder text-subtitle1">
          ATTENZIONE: {{ overdueCount }} studenti risultano fuori dall'aula da più di 15 minuti!
        </div>
        <div class="text-caption text-red-100">
          I collaboratori scolastici sono pregati di verificare i corridoi e i servizi del piano.
        </div>
      </q-banner>
    </div>

    <!-- Active Passes Grid -->
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h6 text-weight-bold text-slate-800">
        Elenco Uscite Attive
      </div>
      <q-input
        v-model="searchQuery"
        dense
        outlined
        bg-color="white"
        placeholder="Cerca per studente..."
        style="width: 260px"
      >
        <template #append>
          <q-icon name="search" />
        </template>
      </q-input>
    </div>

    <div v-if="filteredPasses.length === 0" class="q-pa-xl text-center bg-white rounded-xl border-dashed border-slate-200">
      <q-icon name="verified_user" size="56px" color="positive" class="q-mb-sm" />
      <div class="text-h6 text-weight-bold text-slate-700">Nessun alunno nei corridoi</div>
      <div class="text-body2 text-slate-400">Tutte le aule sono regolari con studenti al loro posto</div>
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="pass in filteredPasses" :key="pass.id" class="col-12 col-md-6 col-lg-4">
        <q-card
          class="rounded-xl shadow-soft transition-all"
          :class="[
            getElapsedMinutes(pass.departed_at) > 15
              ? 'border-2 border-negative bg-red-50/50 shadow-md'
              : 'border-slate-100 bg-white'
          ]"
        >
          <q-card-section>
            <div class="row items-center justify-between q-mb-xs">
              <div>
                <div class="text-h6 text-weight-bold text-slate-800">
                  {{ pass.student_name || 'Alunno' }}
                </div>
                <div class="text-caption text-slate-500">
                  Classe: <strong class="text-slate-700">{{ pass.class_id }}</strong>
                </div>
              </div>
              <q-badge
                :color="getPassTypeColor(pass.pass_type)"
                class="q-pa-xs text-weight-bold"
              >
                {{ getPassTypeLabel(pass.pass_type) }}
              </q-badge>
            </div>

            <!-- Elapsed Time Display -->
            <div
              class="row items-center justify-between q-pa-sm rounded-lg q-my-sm"
              :class="getElapsedMinutes(pass.departed_at) > 15 ? 'bg-red-100 text-negative' : 'bg-slate-100 text-slate-800'"
            >
              <div class="row items-center q-gutter-x-xs">
                <q-icon name="schedule" size="sm" />
                <span class="text-weight-bold">
                  Uscito alle {{ formatTime(pass.departed_at) }}
                </span>
              </div>
              <div class="text-weight-bolder text-subtitle2">
                {{ getElapsedMinutes(pass.departed_at) }} min fa
              </div>
            </div>

            <div v-if="pass.notes" class="text-caption text-slate-600 q-mb-md italic">
              "{{ pass.notes }}"
            </div>

            <q-btn
              color="positive"
              icon="check"
              label="Accompagnato / Rientrato"
              class="full-width rounded-lg"
              unelevated
              @click="markReturned(pass.id)"
            />
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const activePasses = ref([])
const searchQuery = ref('')
const nowTime = ref(Date.now())
let refreshInterval = null

const overdueCount = computed(() => {
  return activePasses.value.filter(p => getElapsedMinutes(p.departed_at) > 15).length
})

const bathroomCount = computed(() => {
  return activePasses.value.filter(p => p.pass_type === 'bathroom').length
})

const otherCount = computed(() => {
  return activePasses.value.filter(p => p.pass_type !== 'bathroom').length
})

const filteredPasses = computed(() => {
  if (!searchQuery.value.trim()) return activePasses.value
  const q = searchQuery.value.toLowerCase()
  return activePasses.value.filter(p =>
    (p.student_name && p.student_name.toLowerCase().includes(q)) ||
    (p.class_id && p.class_id.toLowerCase().includes(q)) ||
    (p.pass_type && p.pass_type.toLowerCase().includes(q))
  )
})

function getPassTypeLabel(type) {
  switch (type) {
    case 'bathroom': return 'Bagno'
    case 'infirmary': return 'Infermeria'
    case 'office': return 'Presidenza/Segreteria'
    case 'staff': return 'Personale ATA'
    default: return type
  }
}

function getPassTypeColor(type) {
  switch (type) {
    case 'bathroom': return 'primary'
    case 'infirmary': return 'negative'
    case 'office': return 'warning'
    case 'staff': return 'secondary'
    default: return 'slate-500'
  }
}

function formatTime(isoStr) {
  if (!isoStr) return '-'
  const d = new Date(isoStr)
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function getElapsedMinutes(departedAt) {
  if (!departedAt) return 0
  const dep = new Date(departedAt).getTime()
  const diffMs = nowTime.value - dep
  return Math.max(0, Math.floor(diffMs / 60000))
}

async function loadActivePasses() {
  try {
    const res = await api.get('/security-passes/hall-passes/active')
    activePasses.value = res.data?.passes || []
  } catch (err) {
    console.error('Errore caricamento pass attivi:', err)
  }
}

async function markReturned(passId) {
  try {
    await api.post(`/security-passes/hall-passes/${passId}/return`)
    $q.notify({ type: 'positive', message: 'Rientro confermato' })
    await loadActivePasses()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore registrazione rientro' })
  }
}

onMounted(() => {
  loadActivePasses()
  refreshInterval = setInterval(() => {
    nowTime.value = Date.now()
    loadActivePasses()
  }, 15000)
})

onBeforeUnmount(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>
