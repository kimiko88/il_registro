<template>
  <q-page class="q-pa-md bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <div class="text-h4 text-weight-bold text-slate-800">
          <q-icon name="help_center" color="primary" class="q-mr-sm" />
          Centro Supporto & FAQ
        </div>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          Benvenuto nel portale di aiuto. Qui trovi risposte rapide per il ruolo di 
          <q-badge color="primary" class="text-subtitle2 q-px-sm text-capitalize">{{ currentRoleLabel }}</q-badge>
        </div>
      </div>
      
      <!-- Connection Status Badge -->
      <div>
        <q-chip :color="isOnline ? 'green-1' : 'red-1'" :text-color="isOnline ? 'positive' : 'negative'" icon="wifi" class="text-weight-bold">
          {{ isOnline ? 'Online' : 'Offline - Modalità locale' }}
        </q-chip>
      </div>
    </div>

    <!-- Offline Banner -->
    <q-banner v-if="!isOnline" class="bg-warning text-white rounded-borders q-mb-lg shadow-sm" dense>
      <template v-slot:avatar>
        <q-icon name="cloud_off" />
      </template>
      Sei offline. Le FAQ caricate in cache rimangono accessibili. I messaggi di supporto inviati ora verranno messi in coda e spediti appena tornerai online.
    </q-banner>

    <div class="row q-col-gutter-lg">
      <!-- FAQ Section -->
      <div class="col-12 col-md-8">
        <q-card class="shadow-sm rounded-lg q-pa-md bg-white">
          <!-- Search and Category Filters -->
          <div class="row q-col-gutter-sm items-center q-mb-md">
            <div class="col-12 col-sm-6">
              <q-input 
                v-model="searchQuery" 
                placeholder="Cerca tra le domande frequenti..." 
                outlined 
                dense
                clearable
              >
                <template v-slot:append>
                  <q-icon name="search" />
                </template>
              </q-input>
            </div>
            <div class="col-12 col-sm-6 text-right">
              <q-btn-toggle
                v-model="selectedCategory"
                toggle-color="primary"
                flat
                dense
                no-caps
                :options="categoryOptions"
              />
            </div>
          </div>

          <!-- FAQ Accordion -->
          <q-list bordered class="rounded-borders separator">
            <template v-if="filteredFaqs.length > 0">
              <q-expansion-item
                v-for="(faq, index) in filteredFaqs"
                :key="index"
                :label="faq.question"
                header-class="text-weight-bold text-slate-800"
                group="faq-accordion"
              >
                <q-card class="bg-slate-50">
                  <q-card-section class="text-slate-600 text-body2 line-height-relaxed">
                    {{ faq.answer }}
                  </q-card-section>
                </q-card>
              </q-expansion-item>
            </template>
            <div v-else class="text-center q-pa-xl text-grey-5">
              <q-icon name="sentiment_dissatisfied" size="3em" class="q-mb-sm" />
              <div>Nessuna risposta trovata per "{{ searchQuery }}"</div>
            </div>
          </q-list>
        </q-card>
      </div>

      <!-- Contact Support & Offline Queue -->
      <div class="col-12 col-md-4">
        <q-card class="shadow-sm rounded-lg bg-white q-pa-md">
          <div class="text-h6 text-slate-800 q-mb-md">
            <q-icon name="mail" color="primary" class="q-mr-xs" />
            Contatta Assistenza
          </div>

          <q-form @submit.prevent="handleSubmitTicket" class="q-gutter-y-md">
            <q-input
              v-model="ticket.subject"
              label="Oggetto"
              outlined
              dense
              required
              :rules="[val => !!val || 'Oggetto richiesto']"
            />
            <q-input
              v-model="ticket.message"
              label="Messaggio"
              type="textarea"
              outlined
              dense
              autogrow
              required
              :rules="[val => !!val || 'Messaggio richiesto']"
            />
            
            <q-btn 
              type="submit" 
              color="primary" 
              unelevated 
              no-caps 
              class="full-width rounded-md"
              :loading="submitting"
            >
              {{ isOnline ? 'Invia Richiesta' : 'Salva in coda offline' }}
            </q-btn>
          </q-form>

          <!-- Queue Indicator -->
          <div v-if="queuedTickets.length > 0" class="q-mt-lg">
            <div class="text-subtitle2 text-grey-7 q-mb-sm">
              <q-icon name="hourglass_empty" color="warning" />
              Richieste in coda offline ({{ queuedTickets.length }})
            </div>
            <q-list bordered dense class="rounded-borders">
              <q-item v-for="(t, idx) in queuedTickets" :key="idx">
                <q-item-section>
                  <q-item-label class="text-weight-bold text-caption">{{ t.subject }}</q-item-label>
                  <q-item-label caption class="ellipsis">{{ t.message }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge color="warning">In attesa</q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </div>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from 'src/stores/auth'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const authStore = useAuthStore()
const { userRole } = storeToRefs(authStore)

const isOnline = ref(navigator.onLine)
const searchQuery = ref('')
const selectedCategory = ref('all')
const submitting = ref(false)
const queuedTickets = ref([])

const ticket = ref({
  subject: '',
  message: ''
})

const currentRoleLabel = computed(() => {
  const labels = {
    superadmin: 'Super Admin',
    admin: 'Amministratore',
    secretary: 'Segretario',
    teacher: 'Docente',
    student: 'Studente',
    parent: 'Genitore'
  }
  return labels[userRole.value] || userRole.value || 'Utente'
})

// FAQ database mapped by role
const faqsByRole = {
  superadmin: [
    { category: 'scuole', question: 'Come posso creare una nuova scuola?', answer: 'Vai nella sezione "Gestione Scuole", clicca su "Nuova Scuola", inserisci i dettagli (Nome, Indirizzo, Codice Meccanografico) e conferma.' },
    { category: 'sistema', question: 'Cosa indica il pannello Monitoraggio?', answer: 'Mostra i tempi di risposta delle API del database e lo stato di salute generale dei microservizi collegati in tempo reale.' },
    { category: 'sistema', question: 'Come faccio a visualizzare gli Audit Logs?', answer: 'Tutte le azioni significative degli utenti vengono registrate cronologicamente nella scheda "Audit Logs" per conformità GDPR.' },
    { category: 'sicurezza', question: 'Come posso configurare la policy password globale o l\'obbligo 2FA?', answer: 'In "Impostazioni Sistema" seleziona "Sicurezza", quindi configura la lunghezza minima password, i giorni di scadenza e il toggle per il 2FA obbligatorio per lo staff.' }
  ],
  admin: [
    { category: 'scuole', question: 'Come configuro i plessi e le sezioni della scuola?', answer: 'Dalla sezione "La mia scuola" puoi configurare i plessi (Sede centrale, succursali), i parametri di rete e le sezioni associate al tuo istituto scolastico.' },
    { category: 'utenti', question: 'Come posso inserire o importare massivamente i docenti e gli studenti?', answer: 'Vai in "Gestione Utenti" → "Import Massivo". Puoi caricare un file Excel o CSV strutturato per importare centinaia di account contemporaneamente.' },
    { category: 'analytics', question: 'Come funziona l\'analisi predittiva della dispersione scolastica?', answer: 'Nel pannello "Analytics & BI" l\'algoritmo identifica in automatico gli studenti che superano il 20% di assenze o hanno una media inferiore al 6.0.' }
  ],
  secretary: [
    { category: 'classi', question: 'Come posso iscrivere uno studente ad una classe?', answer: 'Seleziona "Studenti" dal menu principale, clicca sullo studente desiderato e usa il selettore classe per impostare la classe corrente.' },
    { category: 'didattica', question: 'Come importare i libri di testo?', answer: 'Accedi alla sezione "Libri di Testo" e compila le schede dei volumi consigliati o adottati per ciascuna classe, inserendo ISBN, editore e prezzo.' },
    { category: 'comunicazioni', question: 'Come inviare una circolare alle famiglie con richiesta di presa d\'atto?', answer: 'Vai su "Comunicazioni", clicca su "Nuovo Messaggio", seleziona il target "Genitori", spunta "Obbligo presa d\'atto" e pubblica la circolare.' },
    { category: 'certificati', question: 'Come genero e stampo un certificato con timbro digitale?', answer: 'In "Certificati" seleziona lo studente, il tipo di certificato (iscrizione, frequenza, esiti) e clicca "Genera PDF". Il documento viene firmato digitalmente dal sistema.' }
  ],
  teacher: [
    { category: 'registro', question: 'Come si firma il registro elettronico giornaliero?', answer: 'Vai nella sezione "Registro Classe" (Lesson Planner), seleziona la classe e la materia dell\'ora corrente, digita l\'argomento della lezione e clicca su "Firma".' },
    { category: 'voti', question: 'Come uso l\'inserimento voti rapido Matrix View a tastiera?', answer: 'Accedi a "Voti" e passa alla modalità Matrix. Puoi navigare con TAB/FRECCE, digitare il voto e premere INVIO per salvare all\'istante.' },
    { category: 'scrutinio', question: 'Come compilo la proposta di voto e il giudizio per lo Scrutinio?', answer: 'Nella sezione "Scrutinio" seleziona la classe ed inserisci il voto proposto, la condotta e il giudizio di livello. Se sei coordinatore, verifica le proposte dei colleghi e chiudi il tabellone.' },
    { category: 'pdp', question: 'Come applico le misure compensative e dispensative PDP per studenti BES/DSA?', answer: 'Dalla scheda "PDP / PEI" del Consiglio di Classe imposta le misure attive. Verranno mostrate con icone dedicate nel registro durante le valutazioni.' },
    { category: 'colloqui', question: 'Come imposto la mia disponibilità oraria per i colloqui?', answer: 'Accedi alla sezione "Colloqui" e definisci le fasce orarie settimanali in cui i genitori possono prenotarsi.' }
  ],
  student: [
    { category: 'voti', question: 'Dove vedo la mia media scolastica e come uso il Simulatore?', answer: 'La media dei voti è mostrata sulla Dashboard e in "I Miei Voti". Usa il "Simulatore Media" per calcolare il voto che ti serve nella prossima verifica.' },
    { category: 'didattica', question: 'Come scaricare il materiale didattico?', answer: 'Vai nella sezione "Materiale Didattico". Qui troverai tutti i file inseriti dai tuoi professori divisi per materia.' },
    { category: 'compiti', question: 'Come segno un compito come completato?', answer: 'Nella sezione "Compiti" puoi vedere la lista dei compiti assegnati e spuntare la casella di completamento.' },
    { category: 'pcto', question: 'Dove consulto il monte ore PCTO accumulato?', answer: 'Nella sezione "PCTO & Portfolio" vedi il totale delle ore svolte presso le strutture convenzionate, i progetti e le valutazioni dei tutor.' }
  ],
  parent: [
    { category: 'assenze', question: 'Come giustificare un\'assenza o un ritardo online?', answer: 'Vai nella sezione "Presenze", individua l\'assenza con la dicitura "Da Giustificare" e clicca su "Giustifica", inserendo motivo e PIN/OTP.' },
    { category: 'pagopa', question: 'Come pago un avviso di contributo o gita scolastica tramite PagoPA?', answer: 'Accedi a "Pagamenti PagoPA", individua l\'avviso e clicca "Paga con PagoPA" per pagare direttamente online o scaricare il bollettino QR.' },
    { category: 'colloqui', question: 'Come posso prenotare un colloquio con un docente?', answer: 'Accedi alla sezione "Colloqui", seleziona il docente dall\'elenco a discesa e seleziona uno degli slot disponibili per confermare.' },
    { category: 'documenti', question: 'Dove trovo la pagella scolastica di mio figlio?', answer: 'La pagella quadrimestrale e tutti i documenti ufficiali firmati sono reperibili e scaricabili nella sezione "Documenti".' }
  ]
}

const defaultFaqs = [
  { category: 'generale', question: 'Come posso abilitare l\'autenticazione a due fattori (MFA)?', answer: 'Puoi configurare la sicurezza dell\'account, inclusa la password e l\'autenticazione MFA, visitando la sezione "Profilo".' },
  { category: 'generale', question: 'Cosa faccio in caso di smarrimento credenziali?', answer: 'Contatta l\'amministratore del sistema o la segreteria didattica della tua scuola per richiedere un reset della password.' }
]

// Computed list of FAQs for current role
const roleFaqs = computed(() => {
  const specific = faqsByRole[userRole.value] || []
  return [...specific, ...defaultFaqs]
})

// Filter option configurations
const categoryOptions = computed(() => {
  const base = [
    { label: 'Tutte', value: 'all' },
    { label: 'Generale', value: 'generale' }
  ]
  const uniqueCats = new Set(roleFaqs.value.map(f => f.category))
  uniqueCats.forEach(cat => {
    if (cat !== 'generale') {
      base.push({ label: cat.charAt(0).toUpperCase() + cat.slice(1), value: cat })
    }
  })
  return base
})

const filteredFaqs = computed(() => {
  return roleFaqs.value.filter(faq => {
    const matchesSearch = searchQuery.value
      ? faq.question.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
        faq.answer.toLowerCase().includes(searchQuery.value.toLowerCase())
      : true
    const matchesCategory = selectedCategory.value === 'all'
      ? true
      : faq.category === selectedCategory.value
    return matchesSearch && matchesCategory
  })
})

onMounted(() => {
  window.addEventListener('online', handleOnlineStatus)
  window.addEventListener('offline', handleOfflineStatus)
  loadQueuedTickets()
  if (isOnline.value) {
    syncQueuedTickets()
  }
})

onUnmounted(() => {
  window.removeEventListener('online', handleOnlineStatus)
  window.removeEventListener('offline', handleOfflineStatus)
})

function handleOnlineStatus() {
  isOnline.value = true
  syncQueuedTickets()
}

function handleOfflineStatus() {
  isOnline.value = false
}

// Load tickets stored in localStorage
function loadQueuedTickets() {
  const saved = localStorage.getItem('queued_support_tickets')
  if (saved) {
    queuedTickets.value = JSON.parse(saved)
  }
}

// Save ticket in queue
function queueTicket(t) {
  queuedTickets.value.push(t)
  localStorage.setItem('queued_support_tickets', JSON.stringify(queuedTickets.value))
}

// Submit ticket
async function handleSubmitTicket() {
  submitting.value = true
  const payload = {
    subject: ticket.value.subject,
    message: ticket.value.message,
    role: userRole.value,
    timestamp: new Date().toISOString()
  }

  if (!isOnline.value) {
    queueTicket(payload)
    $q.notify({
      type: 'warning',
      message: 'Messaggio salvato offline. Verrà trasmesso non appena ci sarà connessione.',
      position: 'bottom-right'
    })
    ticket.value.subject = ''
    ticket.value.message = ''
    submitting.value = false
    return
  }

  // Simulate API post (with timeout)
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    $q.notify({
      type: 'positive',
      message: 'Richiesta di supporto inviata con successo!',
      position: 'bottom-right'
    })
    ticket.value.subject = ''
    ticket.value.message = ''
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: 'Errore durante l\'invio della richiesta.',
      position: 'bottom-right'
    })
  } finally {
    submitting.value = false
  }
}

// Sync tickets queued offline when connection is restored
async function syncQueuedTickets() {
  if (queuedTickets.value.length === 0) return

  $q.notify({
    type: 'info',
    message: 'Riconnessione rilevata. Invio dei messaggi di supporto in coda...',
    position: 'bottom-right'
  })

  // Send queued tickets
  for (const _t of queuedTickets.value) {
    try {
      // Simulate API post
      await new Promise(resolve => setTimeout(resolve, 800))
    } catch (e) {
      console.error('Failed to sync ticket:', e)
    }
  }

  queuedTickets.value = []
  localStorage.removeItem('queued_support_tickets')
  $q.notify({
    type: 'positive',
    message: 'Tutte le richieste in coda sono state inviate con successo!',
    position: 'bottom-right'
  })
}
</script>

<style scoped>
.line-height-relaxed {
  line-height: 1.6;
}
</style>
