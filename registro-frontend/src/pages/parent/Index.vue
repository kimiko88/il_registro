<template>
  <q-page class="q-pa-md bg-slate-50">
    <!-- Header with Child Switcher -->
    <div class="row items-center justify-between q-mb-xl">
      <div>
        <h1 class="text-h3 text-weight-bold text-outfit bg-clip-text text-transparent bg-gradient-premium q-my-none" style="display: inline-block;">
          Bentornato, Genitore
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">Panoramica delle attività per i tuoi figli</div>
      </div>
      <div v-if="children.length > 0">
        <q-btn-dropdown
          color="indigo-600"
          unelevated
          no-caps
          class="rounded-xl shadow-soft q-px-md"
          :label="selectedChild ? `${selectedChild.firstName} ${selectedChild.lastName}` : 'Seleziona Figlio'"
          icon="face"
        >
          <q-list>
            <q-item
              v-for="child in children"
              :key="child.id"
              clickable
              v-close-popup
              @click="selectChild(child.id)"
              :active="selectedChildId === child.id"
              active-class="bg-blue-1 text-primary"
            >
              <q-item-section avatar>
                <q-avatar size="sm" color="primary" text-color="white">{{ child.firstName.charAt(0) }}</q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ child.firstName }} {{ child.lastName }}</q-item-label>
                <q-item-label caption>{{ child.class }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="row justify-center q-pa-lg">
      <q-spinner color="primary" size="3em" />
    </div>

    <!-- Dashboard Content -->
    <div v-else-if="selectedChild" class="row q-col-gutter-md">
      
      <!-- Quick Stats -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card shadow-soft full-height overflow-hidden">
          <q-card-section>
            <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">Media Voti</div>
            <div class="text-h3 text-weight-bold text-indigo-600 q-mt-sm">7.8</div>
            <div class="row items-center q-mt-sm">
              <q-icon name="trending_up" color="positive" class="q-mr-xs" />
              <span class="text-positive text-caption text-weight-medium">+0.2 vs mese scorso</span>
            </div>
          </q-card-section>
          <q-icon name="grade" class="card-bg-icon text-indigo-100" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card shadow-soft full-height overflow-hidden">
          <q-card-section>
            <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">Assenze</div>
            <div class="text-h3 text-weight-bold text-orange-600 q-mt-sm">3</div>
            <div class="row items-center q-mt-sm">
              <span class="text-caption text-slate-500">Ultima: 12/12/2024</span>
            </div>
          </q-card-section>
          <q-icon name="how_to_reg" class="card-bg-icon text-orange-100" />
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card shadow-soft full-height overflow-hidden">
          <q-card-section>
            <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">Prossimo Colloquio</div>
            <div class="text-h5 text-weight-bold q-mt-sm">Nessuno</div>
            <q-btn flat dense no-caps color="indigo-600" label="Prenota ora" to="/parent/colloqui" class="q-mt-sm rounded-lg" />
          </q-card-section>
          <q-icon name="event" class="card-bg-icon text-slate-100" />
        </q-card>
      </div>

       <div class="col-12 col-sm-6 col-md-3">
        <q-card class="glass-card stat-card shadow-soft full-height overflow-hidden">
          <q-card-section>
            <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">Avvisi</div>
            <div class="text-h3 text-weight-bold text-rose-600 q-mt-sm">2</div>
            <div class="text-caption text-slate-500 q-mt-sm text-weight-medium">Da leggere</div>
          </q-card-section>
          <q-icon name="notifications_active" class="card-bg-icon text-rose-100" />
        </q-card>
      </div>

      <!-- Recent Activities / Grades -->
      <div class="col-12 col-md-8">
        <q-card class="shadow-sm rounded-lg">
          <q-card-section class="row items-center justify-between">
            <div class="text-h6 text-slate-800">Ultimi Voti</div>
            <q-btn flat no-caps color="primary" label="Vedi tutti" to="/parent/grades" />
          </q-card-section>
          <q-separator />
          <q-list separator>
            <q-item v-for="n in 3" :key="n">
              <q-item-section>
                <q-item-label class="text-weight-medium">Matematica</q-item-label>
                <q-item-label caption>Verifica scritta</q-item-label>
              </q-item-section>
              <q-item-section side>
                <div class="row items-center">
                   <q-badge :color="n === 1 ? 'negative' : 'positive'" class="text-subtitle1 q-pa-xs">
                     {{ n === 1 ? '5.0' : '8.5' }}
                   </q-badge>
                   <div class="text-caption text-grey q-ml-md">20 Dic</div>
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Upcoming Events -->
      <div class="col-12 col-md-4">
        <q-card class="shadow-sm rounded-lg full-height">
          <q-card-section>
            <div class="text-h6 text-slate-800 q-mb-sm">Prossimi Eventi</div>
            <q-timeline color="primary" layout="dense">
              <q-timeline-entry title="Vacanza Invernale" subtitle="23 Dic - 7 Gen" icon="school" />
              <q-timeline-entry title="Consiglio di Classe" subtitle="15 Gen" icon="people" color="orange" />
            </q-timeline>
          </q-card-section>
        </q-card>
      </div>

    </div>

    <!-- Empty State -->
    <div v-else class="text-center q-pa-xl">
      <q-icon name="family_restroom" size="4em" color="grey-4" />
      <div class="text-h6 text-grey-6 q-mt-sm">Nessun figlio associato</div>
      <p class="text-grey-5">Contatta la segreteria se ritieni ci sia un errore.</p>
    </div>

    <!-- Quick Actions (FAB on Mobile) -->
    <q-page-sticky position="bottom-right" :offset="[18, 18]" class="lt-md">
      <q-fab icon="add" direction="up" color="primary">
        <q-fab-action color="orange" icon="edit_calendar" label="Giustifica" to="/parent/attendance" />
        <q-fab-action color="secondary" icon="event" label="Colloquio" to="/parent/colloqui" />
      </q-fab>
    </q-page-sticky>

  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'

const parentStore = useParentStore()
const { children, selectedChild, selectedChildId, loading } = storeToRefs(parentStore)
const { fetchChildren, selectChild } = parentStore

onMounted(() => {
  if (children.value.length === 0) {
    fetchChildren()
  }
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
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.card-bg-icon {
  position: absolute;
  right: -10px;
  bottom: -10px;
  font-size: 80px;
  opacity: 0.5;
  z-index: 0;
}

.truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
