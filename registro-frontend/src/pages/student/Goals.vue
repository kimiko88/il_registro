<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="stars" color="amber" class="q-mr-sm" />
          Obiettivi & Badge Studente
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Traccia i tuoi obiettivi accademici personali, accumula punti e sblocca badge
        </p>
      </div>
      <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="loadData" />
    </div>

    <!-- Stats Summary Cards -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Punti Totali</div>
              <div class="text-h4 text-weight-bold text-amber-7 q-mt-xs">{{ totalPoints }} pt</div>
            </div>
            <q-avatar color="amber-1" text-color="amber-8" icon="emoji_events" size="48px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Completati</div>
              <div class="text-h4 text-weight-bold text-positive q-mt-xs">{{ completedCount }}</div>
            </div>
            <q-avatar color="green-1" text-color="positive" icon="check_circle" size="48px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">In Corso</div>
              <div class="text-h4 text-weight-bold text-primary q-mt-xs">{{ inProgressCount }}</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="trending_up" size="48px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Badge Sbloccati</div>
              <div class="text-h4 text-weight-bold text-purple-7 q-mt-xs">{{ unlockedBadgesCount }}/{{ badges.length }}</div>
            </div>
            <q-avatar color="purple-1" text-color="purple-7" icon="military_tech" size="48px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Section Tabs -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
      <q-tabs
        v-model="tab"
        dense
        class="text-slate-600 bg-slate-100 border-b border-slate-200"
        active-color="primary"
        indicator-color="primary"
        align="left"
        no-caps
      >
        <q-tab name="goals" icon="flag" label="I Miei Obiettivi" />
        <q-tab name="badges" icon="military_tech" label="Badge Sbloccati" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="tab" animated class="bg-white">
        <!-- TAB: GOALS -->
        <q-tab-panel name="goals" class="q-pa-lg">
          <div class="row items-center justify-between q-mb-md">
            <div class="row q-gutter-xs">
              <q-btn
                v-for="f in ['all', 'in_progress', 'completed']" :key="f"
                :flat="filterStatus !== f"
                :unelevated="filterStatus === f"
                :color="filterStatus === f ? 'primary' : 'grey-7'"
                dense no-caps class="q-px-sm"
                @click="filterStatus = f"
              >
                {{ f === 'all' ? 'Tutti' : f === 'in_progress' ? 'In Corso' : 'Completati' }}
              </q-btn>
            </div>
          </div>

          <div v-if="loading" class="text-center q-pa-xl">
            <q-spinner-dots color="primary" size="40px" />
          </div>

          <div v-else-if="filteredGoals.length === 0" class="text-center q-pa-xl text-slate-400">
            <q-icon name="emoji_events" size="64px" class="q-mb-md opacity-40" />
            <div class="text-h6">Nessun obiettivo trovato</div>
            <div class="text-caption">I tuoi docenti ti assegneranno traguardi da raggiungere durante l'anno.</div>
          </div>

          <div v-else class="row q-col-gutter-md">
            <div v-for="g in filteredGoals" :key="g.id" class="col-12 col-md-6 col-lg-4">
              <q-card bordered flat class="rounded-xl h-full flex column justify-between">
                <q-card-section>
                  <div class="row items-center justify-between q-mb-xs">
                    <q-chip size="xs" :color="categoryColor(g.category)" text-color="white" class="text-weight-bold">
                      {{ g.category || 'Generale' }}
                    </q-chip>
                    <q-badge :color="g.status === 'completed' ? 'positive' : 'info'">
                      {{ g.status === 'completed' ? 'Completato' : 'In Corso' }}
                    </q-badge>
                  </div>

                  <div class="text-subtitle1 text-weight-bold text-slate-800 q-mt-xs">{{ g.title }}</div>
                  <div class="text-caption text-slate-600 q-mt-xs">{{ g.description }}</div>
                </q-card-section>

                <q-card-section class="bg-slate-50 q-py-sm">
                  <div class="row items-center justify-between text-caption text-slate-600 q-mb-xs">
                    <span>Progresso</span>
                    <span class="text-weight-bold text-amber-8">+{{ g.points || 10 }} Punti</span>
                  </div>
                  <q-linear-progress
                    :value="g.status === 'completed' ? 1.0 : (g.progress || 0.4)"
                    :color="g.status === 'completed' ? 'positive' : 'primary'"
                    rounded size="8px"
                  />

                  <div class="row items-center justify-between q-mt-sm" v-if="g.status !== 'completed'">
                    <span class="text-caption text-slate-400" v-if="g.target_date">Entro {{ formatDate(g.target_date) }}</span>
                    <q-btn
                      flat size="xs" color="positive" icon="check" label="Segna Completato"
                      :loading="updatingId === g.id"
                      @click="markCompleted(g.id)"
                    />
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- TAB: BADGES -->
        <q-tab-panel name="badges" class="q-pa-lg">
          <div class="row q-col-gutter-md">
            <div v-for="b in badges" :key="b.name" class="col-12 col-sm-6 col-md-4 col-lg-3">
              <q-card
                flat bordered
                class="text-center q-pa-md rounded-xl"
                :class="b.unlocked ? 'bg-amber-50 border-amber-200' : 'bg-slate-100 border-slate-200 opacity-60'"
              >
                <div class="text-h2 q-mb-xs">{{ b.icon }}</div>
                <div class="text-subtitle1 text-weight-bold" :class="b.unlocked ? 'text-amber-9' : 'text-slate-500'">
                  {{ b.name }}
                </div>
                <div class="text-caption text-slate-600 q-mt-xs">{{ b.desc }}</div>
                <div class="q-mt-sm">
                  <q-chip v-if="b.unlocked" size="sm" color="amber-8" text-color="white" icon="verified">
                    Sbloccato!
                  </q-chip>
                  <q-chip v-else size="sm" color="grey" text-color="white" icon="lock">
                    Bloccato
                  </q-chip>
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { studentGoalService } from '@/services/studentGoalService'
import { useAuthStore } from '@/stores/auth'
import api from 'src/services/api'

const $q = useQuasar()
const authStore = useAuthStore()
const loading = ref(false)
const updatingId = ref(null)
const tab = ref('goals')
const filterStatus = ref('all')
const goals = ref([])

const badges = ref([
  { icon: '🏆', name: 'Primo Traguardo', desc: 'Completa il tuo primo obiettivo scolastico', unlocked: true },
  { icon: '📚', name: 'Lettore Instancabile', desc: 'Leggi e firma 10 circolari di classe', unlocked: true },
  { icon: '⭐', name: 'Media Eccellente', desc: 'Mantieni una media generale pari o superiore a 8', unlocked: false },
  { icon: '⏰', name: 'Presenza Perfetta', desc: 'Nessuna assenza per 30 giorni consecutivi', unlocked: true },
  { icon: '🎯', name: 'Punto di Svolta', desc: 'Raggiungi 100 punti obiettivo totali', unlocked: false },
  { icon: '🧪', name: 'Pratico Provvetto', desc: 'Ottieni un voto d\'eccellenza in laboratorio', unlocked: false }
])

const totalPoints = computed(() => {
  return goals.value
    .filter(g => g.status === 'completed')
    .reduce((acc, g) => acc + (g.points || 10), 0)
})

const completedCount = computed(() => goals.value.filter(g => g.status === 'completed').length)
const inProgressCount = computed(() => goals.value.filter(g => g.status !== 'completed').length)
const unlockedBadgesCount = computed(() => badges.value.filter(b => b.unlocked).length)

const filteredGoals = computed(() => {
  if (filterStatus.value === 'all') return goals.value
  if (filterStatus.value === 'completed') return goals.value.filter(g => g.status === 'completed')
  return goals.value.filter(g => g.status !== 'completed')
})

const formatDate = (d) => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY') : ''
const categoryColor = (cat) => ({
  Studio: 'primary',
  Presenza: 'positive',
  Voti: 'purple',
  Comportamento: 'deep-orange'
}[cat] || 'indigo')

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const studentId = authStore.user?.student_id || authStore.user?.id
    if (studentId) {
      const res = await studentGoalService.listByStudent(studentId)
      goals.value = res.data || []
    }
  } catch (e) {
    console.error('Errore caricamento obiettivi:', e)
  } finally {
    loading.value = false
  }
}

async function markCompleted(goalId) {
  updatingId.value = goalId
  try {
    await api.patch(`/student-goals/${goalId}/status`, { status: 'completed' })
    $q.notify({ type: 'positive', message: 'Obiettivo completato! Complimenti!' })
    await loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante l\'aggiornamento dell\'obiettivo' })
  } finally {
    updatingId.value = null
  }
}
</script>
