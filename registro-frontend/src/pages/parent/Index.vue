<template>
  <q-page class="q-pa-md bg-slate-50">
    <!-- Header with Child Switcher -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold text-slate-800">Benvenuto, Genitore</div>
        <div class="text-subtitle2 text-slate-500">Panoramica studente</div>
      </div>
      <div v-if="children.length > 0">
        <q-btn-dropdown
          color="primary"
          outline
          no-caps
          rounded
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
      <div class="col-12 col-md-3">
        <q-card class="shadow-sm rounded-lg full-height">
          <q-card-section>
            <div class="text-overline text-slate-500">Media Voti</div>
            <div class="text-h4 text-weight-bold text-primary">7.8</div>
            <div class="row items-center q-mt-sm">
              <q-icon name="trending_up" color="positive" class="q-mr-xs" />
              <span class="text-positive text-caption">+0.2 vs mese scorso</span>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-3">
        <q-card class="shadow-sm rounded-lg full-height">
          <q-card-section>
            <div class="text-overline text-slate-500">Assenze</div>
            <div class="text-h4 text-weight-bold text-orange">3</div>
            <div class="row items-center q-mt-sm">
              <span class="text-caption text-grey">Ultima: 12/12/2024</span>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-3">
        <q-card class="shadow-sm rounded-lg full-height">
          <q-card-section>
            <div class="text-overline text-slate-500">Prossimo Colloquio</div>
            <div class="text-h5 text-weight-bold truncate">Nessuno</div>
            <q-btn flat dense no-caps color="primary" label="Prenota ora" to="/parent/colloqui" class="q-mt-xs" />
          </q-card-section>
        </q-card>
      </div>

       <div class="col-12 col-md-3">
        <q-card class="shadow-sm rounded-lg full-height">
          <q-card-section>
            <div class="text-overline text-slate-500">Avvisi</div>
            <div class="text-h4 text-weight-bold text-red">2</div>
            <div class="text-caption text-grey q-mt-sm">Da leggere</div>
          </q-card-section>
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
