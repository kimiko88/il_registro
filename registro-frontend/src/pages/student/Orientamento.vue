<template>
  <q-page class="q-pa-md bg-slate-50 min-h-screen">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="explore" color="primary" class="q-mr-sm" />
          Orientamento Scolastico & Universitario
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Pianificazione 30 ore annuali, percorsi post-diploma ed eventi formativi
        </p>
      </div>
      <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="loadData" />
    </div>

    <div class="row q-col-gutter-lg">
      <!-- Main Content Tabs -->
      <div class="col-12 col-md-8">
        <q-tabs
          v-model="tab"
          dense
          class="text-slate-600 bg-white rounded-xl shadow-soft q-mb-md border border-slate-100"
          active-color="primary"
          indicator-color="primary"
          align="left"
          no-caps
        >
          <q-tab name="upcoming" icon="event" :label="t('verbaliPage.agenda') || 'Eventi Disponibili'" />
          <q-tab name="registered" icon="bookmark" :label="t('colloquiPage.bookings') || 'I Miei Eventi'" />
          <q-tab name="past" icon="history" :label="t('verbaliPage.resolutions') || 'Storico & Presenze'" />
        </q-tabs>

        <div v-if="loading" class="text-center q-pa-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>

        <q-tab-panels v-else v-model="tab" animated class="bg-transparent">
          <!-- TAB 1: UPCOMING AVAILABLE EVENTS -->
          <q-tab-panel name="upcoming" class="q-pa-none">
            <div v-if="availableEvents.length === 0" class="q-pa-xl text-center bg-white rounded-xl border border-slate-100 shadow-soft">
              <q-icon name="event_busy" size="64px" color="slate-300" class="q-mb-md" />
              <div class="text-h6 text-slate-700">Nessun nuovo evento disponibile</div>
              <div class="text-caption text-slate-500">Sei già iscritto a tutti gli eventi in programma o non ce ne sono di attivi.</div>
            </div>

            <div v-else class="row q-col-gutter-md">
              <div class="col-12 col-md-6" v-for="event in availableEvents" :key="event.id">
                <q-card flat bordered class="rounded-xl shadow-soft bg-white h-full flex column justify-between">
                  <q-card-section>
                    <div class="row items-center justify-between q-mb-xs">
                      <q-chip size="xs" color="teal-1" text-color="teal-9" class="text-weight-bold">
                        {{ event.category || 'Orientamento' }}
                      </q-chip>
                      <q-chip size="xs" color="indigo-1" text-color="indigo-9" icon="timer">
                        {{ event.hours || 2 }}h Riconosciute
                      </q-chip>
                    </div>

                    <div class="text-h6 text-weight-bold text-slate-800 q-mt-xs">{{ event.title }}</div>
                    <div class="text-caption text-slate-600 q-mt-sm">{{ event.description || 'Nessuna descrizione disponibile' }}</div>

                    <div class="q-mt-md text-caption text-slate-500 space-y-1">
                      <div class="row items-center">
                        <q-icon name="schedule" size="xs" class="q-mr-xs text-slate-400" />
                        <span>{{ formatDateTime(event.date) }}</span>
                      </div>
                      <div v-if="event.location" class="row items-center">
                        <q-icon name="place" size="xs" class="q-mr-xs text-slate-400" />
                        <span>{{ event.location }}</span>
                      </div>
                      <div v-if="event.max_attendees" class="row items-center">
                        <q-icon name="group" size="xs" class="q-mr-xs text-slate-400" />
                        <span>Posti limitati: max {{ event.max_attendees }} studenti</span>
                      </div>
                    </div>
                  </q-card-section>

                  <q-card-actions align="right" class="bg-slate-50 border-t border-slate-100 q-pa-sm">
                    <q-btn
                      unelevated
                      color="primary"
                      icon="event_available"
                      :label="t('colloquiPage.booked') || 'Iscriviti'"
                      no-caps
                      class="rounded-lg font-bold"
                      :loading="registeringId === event.id"
                      @click="register(event)"
                    />
                  </q-card-actions>
                </q-card>
              </div>
            </div>
          </q-tab-panel>

          <!-- TAB 2: REGISTERED EVENTS -->
          <q-tab-panel name="registered" class="q-pa-none">
            <div v-if="registeredParticipations.length === 0" class="q-pa-xl text-center bg-white rounded-xl border border-slate-100 shadow-soft">
              <q-icon name="bookmark_border" size="64px" color="slate-300" class="q-mb-md" />
              <div class="text-h6 text-slate-700">Nessuna iscrizione attiva</div>
              <div class="text-caption text-slate-500">Iscriviti a un evento dalla scheda "Eventi Disponibili".</div>
            </div>

            <q-list v-else bordered separator class="bg-white rounded-xl shadow-soft border border-slate-100 overflow-hidden">
              <q-item v-for="part in registeredParticipations" :key="part.id" class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="teal-1" text-color="teal-8" icon="event" size="44px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-800 text-subtitle1">
                    {{ part.event?.title || part.title || 'Evento Orientamento' }}
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    <span v-if="part.event?.location">{{ part.event.location }} &bull; </span>
                    <span>{{ formatDateTime(part.event?.date || part.registered_at) }}</span>
                    <span v-if="part.event?.hours"> &bull; {{ part.event.hours }} ore previste</span>
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip color="positive" text-color="white" icon="check" size="sm" class="font-bold">
                    {{ part.status || 'Iscritto' }}
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>
          </q-tab-panel>

          <!-- TAB 3: PAST & ATTENDED EVENTS -->
          <q-tab-panel name="past" class="q-pa-none">
            <div v-if="pastParticipations.length === 0" class="q-pa-xl text-center bg-white rounded-xl border border-slate-100 shadow-soft">
              <q-icon name="history_toggle_off" size="64px" color="slate-300" class="q-mb-md" />
              <div class="text-h6 text-slate-700">Nessun evento nello storico</div>
              <div class="text-caption text-slate-500">Le presenze confermate dai docenti compariranno qui.</div>
            </div>

            <q-list v-else bordered separator class="bg-white rounded-xl shadow-soft border border-slate-100 overflow-hidden">
              <q-item v-for="part in pastParticipations" :key="part.id" class="q-py-md">
                <q-item-section avatar>
                  <q-avatar :color="part.attended ? 'green-1' : 'grey-2'" :text-color="part.attended ? 'positive' : 'grey-7'" :icon="part.attended ? 'verified' : 'event_note'" size="44px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-800 text-subtitle1">
                    {{ part.event?.title || 'Evento Concluso' }}
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    {{ formatDateTime(part.event?.date || part.registered_at) }}
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="text-right">
                    <q-chip :color="part.attended ? 'positive' : 'grey-6'" text-color="white" size="xs" class="font-bold">
                      {{ part.attended ? 'Presenza Riconosciuta' : (part.status || 'Concluso') }}
                    </q-chip>
                    <div v-if="part.attended && part.event?.hours" class="text-caption text-weight-bold text-primary q-mt-xs">
                      +{{ part.event.hours }} ore
                    </div>
                  </div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-tab-panel>
        </q-tab-panels>
      </div>

      <!-- Sidebar -->
      <div class="col-12 col-md-4">
        <!-- Ore Orientamento Card -->
        <q-card flat bordered class="rounded-xl bg-gradient-primary text-white text-center q-mb-md shadow-soft">
          <q-card-section class="q-pa-lg">
            <div class="text-h2 text-weight-bolder">{{ totalHours }}</div>
            <div class="text-subtitle1 opacity-90">Ore di Orientamento Svolte</div>
            <q-linear-progress
              :value="Math.min(totalHours / 30, 1)"
              color="white"
              track-color="white"
              class="q-mt-md rounded-borders opacity-80"
              size="10px"
            />
            <div class="text-caption opacity-75 q-mt-xs text-right">Target ministeriale: 30 ore / anno</div>
          </q-card-section>
        </q-card>

        <!-- Preference / Career Target Card -->
        <q-card flat bordered class="rounded-xl bg-white shadow-soft q-mb-md border border-slate-100">
          <q-card-section>
            <div class="row items-center justify-between">
              <div class="text-subtitle1 text-weight-bold text-slate-800">Il Mio Profilo di Uscita</div>
              <q-btn flat round dense icon="edit" color="primary" @click="openPreferenceDialog" />
            </div>
            <div class="q-mt-sm text-caption text-slate-600">
              <div v-if="preference.preferred_track">
                <span class="text-weight-bold text-slate-700">Percorso di interesse:</span> {{ preference.preferred_track }}
              </div>
              <div v-if="preference.target_field" class="q-mt-xs">
                <span class="text-weight-bold text-slate-700">Ambito professionale:</span> {{ preference.target_field }}
              </div>
              <div v-if="preference.notes" class="q-mt-xs text-italic text-slate-500">
                "{{ preference.notes }}"
              </div>
              <div v-if="!preference.preferred_track && !preference.target_field" class="text-slate-400">
                Nessuna preferenza impostata. Clicca sull'icona per indicare il tuo percorso universitario o lavorativo.
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Risorse Utili -->
        <q-card flat bordered class="rounded-xl bg-white shadow-soft border border-slate-100">
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold text-slate-800">Risorse & Strumenti</div>
          </q-card-section>
          <q-list separator>
            <q-item clickable v-ripple href="https://www.universitaly.it/" target="_blank">
              <q-item-section avatar><q-icon name="public" color="blue" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium">Universitaly</q-item-label>
                <q-item-label caption>Portale ufficiale del Ministero dell'Università</q-item-label>
              </q-item-section>
              <q-item-section side><q-icon name="open_in_new" size="xs" /></q-item-section>
            </q-item>
            <q-item clickable v-ripple href="https://www.invalsi.it/" target="_blank">
              <q-item-section avatar><q-icon name="school" color="teal" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium">INVALSI & Competenze</q-item-label>
                <q-item-label caption>Quadro nazionale delle competenze</q-item-label>
              </q-item-section>
              <q-item-section side><q-icon name="open_in_new" size="xs" /></q-item-section>
            </q-item>
            <q-item clickable v-ripple @click="openPreferenceDialog">
              <q-item-section avatar><q-icon name="psychology" color="orange" /></q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-medium">Questionario Vocazionale</q-item-label>
                <q-item-label caption>Aggiorna i tuoi interessi e le tue scelte</q-item-label>
              </q-item-section>
              <q-item-section side><q-icon name="chevron_right" /></q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>
    </div>

    <!-- Career Guidance Preference Dialog -->
    <q-dialog v-model="showPrefDialog">
      <q-card style="min-width: 400px; max-width: 90vw;" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold">Preferenze Post-Diploma</div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <q-select
            v-model="prefForm.preferred_track"
            :options="['Università / Laurea Triennale', 'ITS Academy (Tecnici Superiori)', 'Mercato del Lavoro / Impresa', 'Formazione Post-Diploma / AFAM', 'Altro']"
            label="Percorso Preferito"
            outlined
            dense
          />
          <q-input
            v-model="prefForm.target_field"
            label="Ambito di Interesse (es. Informatica, Medicina, Economia)"
            outlined
            dense
          />
          <q-input
            v-model="prefForm.notes"
            label="Note Personali / Obiettivi"
            type="textarea"
            outlined
            dense
            rows="3"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup no-caps />
          <q-btn color="primary" label="Salva Preferenze" :loading="savingPref" no-caps @click="savePreferences" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import api from '@/services/api'

const { t } = useI18n()
const $q = useQuasar()
const tab = ref('upcoming')
const loading = ref(false)
const registeringId = ref(null)

const events = ref([])
const myParticipations = ref([])
const preference = ref({})

const showPrefDialog = ref(false)
const savingPref = ref(false)
const prefForm = ref({
  preferred_track: '',
  target_field: '',
  notes: ''
})

const registeredEventIds = computed(() => {
  return new Set(myParticipations.value.map(p => p.event_id))
})

const availableEvents = computed(() => {
  return events.value.filter(e => !registeredEventIds.value.has(e.id))
})

const registeredParticipations = computed(() => {
  return myParticipations.value.filter(p => !p.attended && p.status !== 'NoShow')
})

const pastParticipations = computed(() => {
  return myParticipations.value.filter(p => p.attended || p.status === 'Attended' || p.status === 'NoShow')
})

const totalHours = computed(() => {
  return myParticipations.value
    .filter(p => p.attended || p.status === 'Attended')
    .reduce((acc, p) => acc + (p.event?.hours || 0), 0)
})

const formatDateTime = (d) => {
  if (!d) return ''
  return new Date(d).toLocaleDateString('it-IT', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const [eventsRes, myRes, prefRes] = await Promise.allSettled([
      api.get('/orientamento/events'),
      api.get('/orientamento/my-events'),
      api.get('/orientamento/preference')
    ])

    if (eventsRes.status === 'fulfilled') {
      events.value = eventsRes.value.data || []
    }
    if (myRes.status === 'fulfilled') {
      myParticipations.value = myRes.value.data || []
    }
    if (prefRes.status === 'fulfilled') {
      preference.value = prefRes.value.data || {}
    }
  } catch (err) {
    console.error('Error fetching orientamento data:', err)
  } finally {
    loading.value = false
  }
}

async function register(event) {
  $q.dialog({
    title: 'Conferma Iscrizione',
    message: `Vuoi iscriverti all'evento di orientamento "${event.title}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    registeringId.value = event.id
    try {
      await api.post('/orientamento/register', { event_id: event.id })
      $q.notify({ type: 'positive', message: 'Iscrizione effettuata con successo!' })
      await loadData()
      tab.value = 'registered'
    } catch (err) {
      $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore durante l\'iscrizione all\'evento' })
    } finally {
      registeringId.value = null
    }
  })
}

function openPreferenceDialog() {
  prefForm.value = {
    preferred_track: preference.value.preferred_track || '',
    target_field: preference.value.target_field || '',
    notes: preference.value.notes || ''
  }
  showPrefDialog.value = true
}

async function savePreferences() {
  savingPref.value = true
  try {
    await api.post('/orientamento/preference', prefForm.value)
    $q.notify({ type: 'positive', message: 'Preferenze salvate con successo!' })
    preference.value = { ...prefForm.value }
    showPrefDialog.value = false
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel salvataggio delle preferenze' })
  } finally {
    savingPref.value = false
  }
}
</script>
