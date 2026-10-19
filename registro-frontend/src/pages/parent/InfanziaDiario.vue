<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          Diario di Bordo &amp; Infanzia
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          Esperienze educative, laboratori, presenze e mensa per i più piccoli (3-6 anni)
        </p>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-select
          v-model="selectedChildId"
          :options="childrenList"
          option-value="id"
          option-label="name"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          label="Seleziona Figlio/a"
          style="min-width: 200px"
          @update:model-value="onChildChanged"
        />
      </div>
    </div>

    <!-- Quick Status Cards: Today's Presence & Meal -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-6">
        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white">
          <q-card-section class="row items-center q-gutter-x-md">
            <q-avatar
              :color="todayStatus.is_present ? 'positive' : 'grey-5'"
              text-color="white"
              icon="school"
              size="lg"
            />
            <div>
              <div class="text-caption text-slate-500">Stato Presenza Oggi</div>
              <div class="text-h6 text-weight-bold" :class="todayStatus.is_present ? 'text-positive' : 'text-slate-600'">
                {{ todayStatus.is_present ? 'Regolarmente a Scuola' : 'Assente / Non rilevato' }}
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-6">
        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white">
          <q-card-section class="row items-center q-gutter-x-md">
            <q-avatar
              :color="todayStatus.has_meal ? 'indigo-7' : 'grey-5'"
              text-color="white"
              icon="restaurant"
              size="lg"
            />
            <div>
              <div class="text-caption text-slate-500">Servizio Mensa Scolastica</div>
              <div class="text-h6 text-weight-bold" :class="todayStatus.has_meal ? 'text-indigo-8' : 'text-slate-600'">
                {{ todayStatus.has_meal ? 'Pasto Prenotato in Cucina' : 'Non usufruisce del pasto' }}
              </div>
              <div v-if="todayStatus.diet_notes" class="text-caption text-warning text-weight-bold">
                Dieta: {{ todayStatus.diet_notes }}
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Content Tabs: Diario Stories and Campi di Esperienza -->
    <q-card class="rounded-xl shadow-soft border-slate-100 bg-white">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey-7"
        active-color="primary"
        indicator-color="primary"
        align="left"
      >
        <q-tab name="diario" icon="photo_library" label="Diario Fotografico &amp; Attività" />
        <q-tab name="campi" icon="psychology" label="Percorso nei Campi di Esperienza" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated>
        <!-- Stories Feed -->
        <q-tab-panel name="diario" class="q-pa-lg">
          <div v-if="stories.length === 0" class="q-pa-xl text-center text-slate-400">
            <q-icon name="menu_book" size="64px" color="slate-300" class="q-mb-sm" />
            <div class="text-h6 text-weight-bold text-slate-600">Nessuna attività pubblicata di recente</div>
            <div class="text-body2">Le maestre pubblicheranno qui i racconti delle giornate e le foto dei laboratori.</div>
          </div>

          <div v-else class="row q-col-gutter-lg">
            <div v-for="s in stories" :key="s.id" class="col-12 col-md-6">
              <q-card class="rounded-xl shadow-soft border-slate-100 bg-white full-height">
                <q-card-section>
                  <div class="row items-center justify-between q-mb-sm">
                    <div class="text-h6 text-weight-bold text-slate-800">
                      {{ s.title }}
                    </div>
                    <q-badge color="primary" class="q-pa-xs">
                      {{ s.story_date }}
                    </q-badge>
                  </div>

                  <p class="text-body1 text-slate-700 q-mb-md leading-relaxed">
                    {{ s.narrative_text }}
                  </p>

                  <div v-if="getPhotos(s.photo_urls_json).length > 0" class="row q-col-gutter-sm">
                    <div
                      v-for="(photo, idx) in getPhotos(s.photo_urls_json)"
                      :key="idx"
                      class="col-6"
                    >
                      <q-img
                        :src="photo"
                        ratio="1.33"
                        class="rounded-lg shadow-sm border border-slate-200"
                        placeholder-src="https://placehold.co/400x300/e2e8f0/64748b?text=Attività+Infanzia"
                      />
                    </div>
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- Campi di Esperienza Observations -->
        <q-tab-panel name="campi" class="q-pa-lg">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">
            Osservazioni Evolutive (Indicazioni Nazionali Scuola dell'Infanzia)
          </div>
          <p class="text-caption text-slate-500 q-mb-lg">
            La valutazione alla scuola dell'infanzia ha carattere formativo ed è basata sull'osservazione sistematica dei progressi del bambino nei cinque campi di esperienza.
          </p>

          <div v-if="observations.length === 0" class="q-pa-lg text-center bg-slate-50 rounded-xl border border-slate-200 text-slate-400">
            Nessuna osservazione registrata al momento per questo alunno.
          </div>

          <div v-else class="row q-col-gutter-md">
            <div v-for="obs in observations" :key="obs.id" class="col-12 col-md-6">
              <q-card class="rounded-xl border border-slate-200 bg-white shadow-soft">
                <q-card-section>
                  <div class="row items-center justify-between q-mb-xs">
                    <div class="text-subtitle1 text-weight-bold text-slate-800">
                      {{ obs.campo_esperienza }}
                    </div>
                    <q-badge color="positive">
                      {{ obs.level_descriptor }}
                    </q-badge>
                  </div>
                  <div class="text-caption text-slate-500 q-mb-sm">
                    Data: {{ obs.observation_date }}
                  </div>
                  <div v-if="obs.notes" class="text-body2 text-slate-700 italic">
                    "{{ obs.notes }}"
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'

const activeTab = ref('diario')
const childrenList = ref([])
const selectedChildId = ref(null)

const todayStatus = ref({
  is_present: false,
  has_meal: false,
  diet_notes: ''
})

const stories = ref([])
const observations = ref([])

function getPhotos(json) {
  if (!json) return []
  try {
    const parsed = JSON.parse(json)
    return Array.isArray(parsed) ? parsed : []
  } catch (e) {
    return []
  }
}

async function loadChildren() {
  try {
    const res = await api.get('/users/me/children')
    const raw = res.data?.children || res.data || []
    childrenList.value = raw.map(c => ({
      id: c.id,
      name: `${c.first_name || ''} ${c.last_name || ''}`.trim() || c.username
    }))
    if (childrenList.value.length > 0) {
      selectedChildId.value = childrenList.value[0].id
      await onChildChanged(selectedChildId.value)
    }
  } catch (err) {
    console.error('Errore caricamento figli:', err)
  }
}

async function onChildChanged(childId) {
  if (!childId) return
  await Promise.all([
    loadObservations(childId),
    loadTodayAttendance(),
    loadStories()
  ])
}

async function loadObservations(childId) {
  try {
    const res = await api.get('/infanzia/observations/by-student', {
      params: { student_id: childId }
    })
    observations.value = res.data?.observations || []
  } catch (err) {
    console.error('Errore caricamento osservazioni:', err)
  }
}

async function loadTodayAttendance() {
  try {
    const res = await api.get('/infanzia/sections')
    const sections = res.data?.sections || []
    if (sections.length > 0) {
      const today = new Date().toISOString().substring(0, 10)
      const attRes = await api.get(`/infanzia/sections/${sections[0].id}/attendances`, {
        params: { date: today }
      })
      const list = attRes.data?.attendances || []
      const found = list.find(a => a.student_id === selectedChildId.value)
      if (found) {
        todayStatus.value = found
      } else {
        todayStatus.value = { is_present: true, has_meal: true, diet_notes: '' }
      }
    }
  } catch (err) {
    console.error('Errore verifica presenze:', err)
  }
}

async function loadStories() {
  try {
    const res = await api.get('/infanzia/sections')
    const sections = res.data?.sections || []
    if (sections.length > 0) {
      const storyRes = await api.get(`/infanzia/sections/${sections[0].id}/diario`)
      stories.value = storyRes.data?.stories || []
    }
  } catch (err) {
    console.error('Errore caricamento storie:', err)
  }
}

onMounted(() => {
  loadChildren()
})
</script>
