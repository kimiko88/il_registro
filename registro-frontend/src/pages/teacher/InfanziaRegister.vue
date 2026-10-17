<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('infanzia.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('infanzia.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm items-center">
        <!-- Colored Section Selector -->
        <q-select
          v-model="selectedSectionId"
          :options="sections"
          option-value="id"
          option-label="name"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          label="Sezione Infanzia"
          style="min-width: 220px"
          @update:model-value="onSectionChanged"
        >
          <template #option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section avatar>
                <div
                  class="rounded-full"
                  :style="{ backgroundColor: scope.opt.color_hex || '#3b82f6', width: '16px', height: '16px' }"
                />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ scope.opt.name }}</q-item-label>
              </q-item-section>
            </q-item>
          </template>
        </q-select>

        <q-btn
          color="primary"
          icon="add"
          :label="$t('infanzia.newSectionBtn')"
          class="rounded-lg shadow-sm"
          @click="showNewSectionDialog = true"
        />
      </div>
    </div>

    <!-- Main Tabs -->
    <q-card class="rounded-xl shadow-soft border-slate-100 bg-white">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey-7"
        active-color="primary"
        indicator-color="primary"
        align="left"
      >
        <q-tab name="attendance" icon="lunch_dining" :label="$t('infanzia.attendancesTab')" />
        <q-tab name="campi" icon="psychology" :label="$t('infanzia.campiTab')" />
        <q-tab name="diario" icon="photo_camera" :label="$t('infanzia.diarioTab')" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated>
        <!-- TAB 1: Rilevazione Presenze e Mensa Cucina -->
        <q-tab-panel name="attendance" class="q-pa-lg">
          <div class="row items-center justify-between q-mb-md">
            <div class="row items-center q-gutter-x-md">
              <q-input
                v-model="selectedDate"
                type="date"
                outlined
                dense
                label="Data del Giorno"
                @update:model-value="loadAttendances"
              />
              <q-btn
                color="positive"
                icon="save"
                label="Salva Appello &amp; Mensa"
                unelevated
                class="rounded-lg"
                @click="saveAttendances"
              />
            </div>

            <!-- Summary for school kitchen -->
            <div class="row q-gutter-x-sm">
              <q-badge color="primary" class="q-pa-sm text-subtitle2">
                Presenti: {{ totalPresentCount }}
              </q-badge>
              <q-badge color="indigo-7" class="q-pa-sm text-subtitle2">
                Pasti Totali Mensa: {{ totalMealsCount }}
              </q-badge>
              <q-badge v-if="specialDietCount > 0" color="warning" class="q-pa-sm text-subtitle2">
                Diete Speciali: {{ specialDietCount }}
              </q-badge>
            </div>
          </div>

          <q-table
            :rows="rosterAttendances"
            :columns="attendanceColumns"
            row-key="student_id"
            flat
            dense
            :pagination="{ rowsPerPage: 25 }"
          >
            <template #body-cell-present="props">
              <q-td :props="props" align="center">
                <q-toggle
                  v-model="props.row.is_present"
                  color="positive"
                  @update:model-value="val => { if (!val) props.row.has_meal = false }"
                />
              </q-td>
            </template>

            <template #body-cell-meal="props">
              <q-td :props="props" align="center">
                <q-toggle
                  v-model="props.row.has_meal"
                  color="indigo-7"
                  :disable="!props.row.is_present"
                />
              </q-td>
            </template>

            <template #body-cell-diet="props">
              <q-td :props="props">
                <q-input
                  v-model="props.row.diet_notes"
                  outlined
                  dense
                  placeholder="Es. Privo di glutine / lattosio"
                />
              </q-td>
            </template>

            <template #body-cell-notes="props">
              <q-td :props="props">
                <q-input
                  v-model="props.row.notes"
                  outlined
                  dense
                  placeholder="Note (es. entrata 9:30)"
                />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- TAB 2: Osservazioni Campi di Esperienza -->
        <q-tab-panel name="campi" class="q-pa-lg">
          <div class="row q-col-gutter-lg">
            <!-- Form: Inserisci Osservazione -->
            <div class="col-12 col-md-5">
              <q-card class="bg-slate-50 border border-slate-200 rounded-xl">
                <q-card-section>
                  <div class="text-h6 text-weight-bold text-slate-800 q-mb-md">
                    Registra Osservazione Evolutiva
                  </div>

                  <div class="q-gutter-y-md">
                    <q-select
                      v-model="newObservation.student_id"
                      :options="studentsOptions"
                      option-value="id"
                      option-label="name"
                      emit-value
                      map-options
                      outlined
                      dense
                      label="Bambino / Bambina *"
                      @update:model-value="loadStudentObservations"
                    />

                    <q-select
                      v-model="newObservation.campo_esperienza"
                      :options="campiEsperienzaOptions"
                      outlined
                      dense
                      label="Campo di Esperienza *"
                    />

                    <q-select
                      v-model="newObservation.level_descriptor"
                      :options="levelOptions"
                      outlined
                      dense
                      label="Livello di Sviluppo Raggiunto *"
                    />

                    <q-input
                      v-model="newObservation.observation_date"
                      type="date"
                      outlined
                      dense
                      label="Data Osservazione *"
                    />

                    <q-input
                      v-model="newObservation.notes"
                      outlined
                      dense
                      type="textarea"
                      rows="3"
                      label="Note descrittive sul comportamento, socializzazione ed autonomia"
                    />

                    <q-btn
                      color="primary"
                      icon="save"
                      label="Registra nel Profilo del Bambino"
                      class="full-width rounded-lg"
                      unelevated
                      :disable="!newObservation.student_id || !newObservation.campo_esperienza"
                      @click="saveObservation"
                    />
                  </div>
                </q-card-section>
              </q-card>
            </div>

            <!-- List: Storico Osservazioni Bambino -->
            <div class="col-12 col-md-7">
              <div class="text-h6 text-weight-bold text-slate-800 q-mb-md">
                Storico Campi di Esperienza
              </div>

              <div v-if="studentObservations.length === 0" class="q-pa-lg text-center bg-white rounded-xl border-dashed border-slate-200 text-slate-400">
                Seleziona un bambino per visualizzarne le osservazioni registrate.
              </div>

              <q-list separator v-else class="rounded-xl border border-slate-200 bg-white">
                <q-item v-for="obs in studentObservations" :key="obs.id" class="q-py-md">
                  <q-item-section avatar>
                    <q-avatar color="indigo-1" text-color="primary" icon="psychology" />
                  </q-item-section>

                  <q-item-section>
                    <div class="row items-center justify-between q-mb-xs">
                      <div class="text-subtitle1 text-weight-bold text-slate-800">
                        {{ obs.campo_esperienza }}
                      </div>
                      <q-badge color="positive">
                        {{ obs.level_descriptor }}
                      </q-badge>
                    </div>
                    <div class="text-caption text-slate-500 q-mb-xs">
                      Osservato il: {{ obs.observation_date }}
                    </div>
                    <div v-if="obs.notes" class="text-body2 text-slate-700 italic">
                      "{{ obs.notes }}"
                    </div>
                  </q-item-section>
                </q-item>
              </q-list>
            </div>
          </div>
        </q-tab-panel>

        <!-- TAB 3: Diario di Bordo Fotografico -->
        <q-tab-panel name="diario" class="q-pa-lg">
          <div class="row items-center justify-between q-mb-lg">
            <div>
              <div class="text-h6 text-weight-bold text-slate-800">
                Racconto delle Giornate (Visibile alle Famiglie)
              </div>
              <div class="text-caption text-slate-500">
                Condividi attività di gruppo, laboratori di pittura, narrazioni e momenti all'aperto
              </div>
            </div>

            <q-btn
              color="primary"
              icon="add_photo_alternate"
              label="Nuovo Racconto nel Diario"
              class="rounded-lg shadow-sm"
              @click="showNewStoryDialog = true"
            />
          </div>

          <div v-if="stories.length === 0" class="q-pa-xl text-center bg-white rounded-xl border-dashed border-slate-200">
            <q-icon name="menu_book" size="56px" color="slate-300" class="q-mb-sm" />
            <div class="text-h6 text-weight-bold text-slate-600">Nessun racconto pubblicato</div>
            <div class="text-body2 text-slate-400">
              I genitori della sezione potranno leggere e visualizzare le foto delle attività quotidiane.
            </div>
          </div>

          <div v-else class="row q-col-gutter-lg">
            <div v-for="s in stories" :key="s.id" class="col-12 col-md-6 col-lg-4">
              <q-card class="rounded-xl shadow-soft border-slate-100 bg-white">
                <q-card-section>
                  <div class="row items-center justify-between q-mb-xs">
                    <div class="text-h6 text-weight-bold text-slate-800">
                      {{ s.title }}
                    </div>
                    <q-badge color="primary">
                      {{ s.story_date }}
                    </q-badge>
                  </div>

                  <p class="text-body2 text-slate-700 q-mb-md">
                    {{ s.narrative_text }}
                  </p>

                  <div v-if="s.photo_urls_json" class="row q-gutter-xs q-mb-sm">
                    <span class="text-caption text-primary row items-center q-gutter-x-xs">
                      <q-icon name="photo" />
                      <span>Fotografie allegate all'attività</span>
                    </span>
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>

    <!-- Dialog: Nuova Sezione -->
    <q-dialog v-model="showNewSectionDialog" persistent>
      <q-card style="min-width: 440px" class="rounded-xl">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6 text-weight-bold text-slate-800">
            {{ $t('infanzia.newSectionBtn') }}
          </div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="newSection.name"
            outlined
            dense
            label="Nome Sezione (es. Sezione Rossi, Coccinelle) *"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <div>
            <div class="text-caption text-slate-600 q-mb-xs">Colore Distintivo Sezione:</div>
            <div class="row q-gutter-sm">
              <div
                v-for="c in ['#ef4444', '#3b82f6', '#10b981', '#f59e0b', '#8b5cf6', '#ec4899']"
                :key="c"
                class="cursor-pointer rounded-full border-2"
                :style="{
                  backgroundColor: c,
                  width: '32px',
                  height: '32px',
                  borderColor: newSection.color_hex === c ? '#000' : 'transparent'
                }"
                @click="newSection.color_hex = c"
              />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            icon="save"
            label="Crea Sezione"
            :disable="!newSection.name"
            @click="createSection"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Nuovo Racconto Diario -->
    <q-dialog v-model="showNewStoryDialog" persistent>
      <q-card style="min-width: 500px" class="rounded-xl">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6 text-weight-bold text-slate-800">
            Pubblica Racconto nel Diario
          </div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="newStory.title"
            outlined
            dense
            label="Titolo Attività *"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <q-input
            v-model="newStory.story_date"
            type="date"
            outlined
            dense
            label="Data Attività *"
          />

          <q-input
            v-model="newStory.narrative_text"
            outlined
            dense
            type="textarea"
            rows="4"
            label="Racconto dell'Esperienza per i Genitori *"
          />

          <q-input
            v-model="newStory.photos"
            outlined
            dense
            label="URL Foto (separate da virgola)"
            hint="Es. /photos/laboratorio1.jpg, /photos/giardino.jpg"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            icon="publish"
            label="Pubblica Diario"
            :disable="!newStory.title || !newStory.narrative_text"
            @click="publishStory"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const activeTab = ref('attendance')

const sections = ref([])
const selectedSectionId = ref(null)
const selectedDate = ref(new Date().toISOString().substring(0, 10))

const rosterAttendances = ref([])
const studentsOptions = ref([])
const studentObservations = ref([])
const stories = ref([])

const showNewSectionDialog = ref(false)
const newSection = ref({
  name: '',
  color_hex: '#ef4444',
  academic_year: '2025/2026'
})

const campiEsperienzaOptions = [
  'Il sé e l\'altro',
  'Il corpo e il movimento',
  'Immagini, suoni, colori',
  'I discorsi e le parole',
  'La conoscenza del mondo'
]

const levelOptions = [
  'Iniziale',
  'In via di prima acquisizione',
  'In evoluzione',
  'Consolidato'
]

const newObservation = ref({
  student_id: '',
  campo_esperienza: 'Il sé e l\'altro',
  level_descriptor: 'In evoluzione',
  observation_date: new Date().toISOString().substring(0, 10),
  notes: ''
})

const showNewStoryDialog = ref(false)
const newStory = ref({
  title: '',
  story_date: new Date().toISOString().substring(0, 10),
  narrative_text: '',
  photos: ''
})

const attendanceColumns = [
  { name: 'student_id', label: 'Bambino/a', field: 'student_id', align: 'left', sortable: true },
  { name: 'present', label: 'Presente', align: 'center' },
  { name: 'meal', label: 'Pasto Mensa', align: 'center' },
  { name: 'diet', label: 'Dieta Speciale', align: 'left' },
  { name: 'notes', label: 'Note', align: 'left' }
]

const totalPresentCount = computed(() => rosterAttendances.value.filter(a => a.is_present).length)
const totalMealsCount = computed(() => rosterAttendances.value.filter(a => a.has_meal).length)
const specialDietCount = computed(() => rosterAttendances.value.filter(a => a.diet_notes && a.diet_notes.trim()).length)

async function loadSections() {
  try {
    const res = await api.get('/infanzia/sections')
    sections.value = res.data?.sections || []
    if (sections.value.length > 0 && !selectedSectionId.value) {
      selectedSectionId.value = sections.value[0].id
      await onSectionChanged(selectedSectionId.value)
    }
  } catch (err) {
    console.error('Errore caricamento sezioni infanzia:', err)
  }
}

async function onSectionChanged(sectionId) {
  if (!sectionId) return
  await Promise.all([
    loadAttendances(),
    loadSectionStories()
  ])
}

async function loadAttendances() {
  if (!selectedSectionId.value) return
  try {
    const res = await api.get(`/infanzia/sections/${selectedSectionId.value}/attendances`, {
      params: { date: selectedDate.value }
    })
    const atts = res.data?.attendances || []
    if (atts.length > 0) {
      rosterAttendances.value = atts
    } else {
      // Default children list for section
      rosterAttendances.value = [
        { student_id: 'Alunno 1', is_present: true, has_meal: true, diet_notes: '', notes: '' },
        { student_id: 'Alunno 2', is_present: true, has_meal: true, diet_notes: '', notes: '' },
        { student_id: 'Alunno 3', is_present: true, has_meal: false, diet_notes: '', notes: '' }
      ]
    }
    studentsOptions.value = rosterAttendances.value.map(a => ({ id: a.student_id, name: a.student_id }))
  } catch (err) {
    console.error('Errore caricamento presenze:', err)
  }
}

async function saveAttendances() {
  if (!selectedSectionId.value) return
  try {
    await api.post(`/infanzia/sections/${selectedSectionId.value}/attendances`, {
      date: selectedDate.value,
      attendances: rosterAttendances.value
    })
    $q.notify({ type: 'positive', message: 'Presenze e pasti registrati per la cucina!' })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore salvataggio presenze' })
  }
}

async function loadStudentObservations() {
  if (!newObservation.value.student_id) return
  try {
    const res = await api.get('/infanzia/observations/by-student', {
      params: { student_id: newObservation.value.student_id }
    })
    studentObservations.value = res.data?.observations || []
  } catch (err) {
    console.error('Errore osservazioni:', err)
  }
}

async function saveObservation() {
  try {
    await api.post('/infanzia/observations', newObservation.value)
    $q.notify({ type: 'positive', message: 'Osservazione registrata nel campo di esperienza' })
    await loadStudentObservations()
    newObservation.value.notes = ''
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore salvataggio osservazione' })
  }
}

async function loadSectionStories() {
  if (!selectedSectionId.value) return
  try {
    const res = await api.get(`/infanzia/sections/${selectedSectionId.value}/diario`)
    stories.value = res.data?.stories || []
  } catch (err) {
    console.error('Errore storie diario:', err)
  }
}

async function publishStory() {
  if (!selectedSectionId.value) return
  try {
    const photosArray = newStory.value.photos.split(',').map(s => s.trim()).filter(Boolean)
    await api.post(`/infanzia/sections/${selectedSectionId.value}/diario`, {
      title: newStory.value.title,
      story_date: newStory.value.story_date,
      narrative_text: newStory.value.narrative_text,
      photo_urls_json: JSON.stringify(photosArray)
    })
    $q.notify({ type: 'positive', message: 'Racconto pubblicato nel diario delle famiglie!' })
    showNewStoryDialog.value = false
    newStory.value = {
      title: '',
      story_date: new Date().toISOString().substring(0, 10),
      narrative_text: '',
      photos: ''
    }
    await loadSectionStories()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore pubblicazione racconto' })
  }
}

async function createSection() {
  try {
    await api.post('/infanzia/sections', newSection.value)
    $q.notify({ type: 'positive', message: 'Nuova sezione creata!' })
    showNewSectionDialog.value = false
    newSection.value = { name: '', color_hex: '#ef4444', academic_year: '2025/2026' }
    await loadSections()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore creazione sezione' })
  }
}

onMounted(() => {
  loadSections()
})
</script>
