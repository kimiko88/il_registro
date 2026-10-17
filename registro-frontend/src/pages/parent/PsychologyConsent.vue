<template>
  <q-page class="q-pa-md q-pa-lg-xl bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div class="row items-center q-gutter-x-sm">
        <q-avatar color="pink-1" text-color="pink-8" icon="psychology" size="44px" />
        <div>
          <h1 class="text-h5 text-weight-bolder text-slate-900 q-my-none">
            {{ $t('parentPsychology.title') || 'Sportello Psicologico (CIC) & Consenso Genitoriale' }}
          </h1>
          <div class="text-caption text-slate-500">
            {{ $t('parentPsychology.subtitle') || 'Consenso informato obbligatorio per l\'accesso dei figli allo sportello d\'ascolto (L. 56/1989)' }}
          </div>
        </div>
      </div>
    </div>

    <!-- Child Selector -->
    <q-card flat bordered class="q-pa-md bg-white rounded-borders shadow-sm q-mb-lg">
      <div class="row items-center justify-between">
        <div class="row items-center q-gutter-x-md">
          <q-icon name="face" size="28px" color="primary" />
          <div class="text-subtitle2 text-weight-bold text-slate-800">
            {{ $t('parentPsychology.selectChild') || 'Seleziona Figlio/a:' }}
          </div>
        </div>
        <div style="min-width: 250px">
          <q-select
            v-model="selectedChild"
            :options="childrenList"
            option-label="name"
            outlined
            dense
          />
        </div>
      </div>
    </q-card>

    <div class="row q-col-gutter-lg">
      <!-- Consent Card -->
      <div class="col-12 col-md-7">
        <q-card flat bordered class="q-pa-lg bg-white rounded-borders shadow-sm full-height">
          <div class="row items-center justify-between q-mb-sm">
            <div class="text-h6 text-weight-bold text-slate-800">
              ✍️ {{ $t('parentPsychology.cardConsentTitle') || 'Consenso Informato Obbligatorio (A.S. 2025/2026)' }}
            </div>
            <q-badge
              :color="selectedChild?.has_consent ? 'positive' : 'warning'"
              class="q-pa-xs text-weight-bold"
              :label="selectedChild?.has_consent ? ($t('parentPsychology.statusSigned') || 'Consenso Rilasciato') : ($t('parentPsychology.statusPending') || 'In Attesa di Firma')"
            />
          </div>

          <div class="text-caption text-slate-500 q-mb-md">
            {{ $t('parentPsychology.legalNotice') || 'Ai sensi dell\'art. 316 c.c., della Legge 56/1989 e delle linee di indirizzo MIM, l\'accesso al servizio di consulenza psicologica per studenti minori è subordinato al consenso informato dei genitori.' }}
          </div>

          <!-- Document Explanation Box -->
          <div class="q-pa-md bg-slate-50 rounded-borders border border-slate-200 text-body2 text-slate-700 q-mb-md">
            <p class="q-mb-xs"><strong>Finalità del servizio:</strong> Il Centro di Informazione e Consulenza (CIC) offre colloqui di ascolto e sostegno finalizzati alla promozione del benessere relazionale, motivazionale e psicologico a scuola.</p>
            <p class="q-mb-xs"><strong>Segreto professionale:</strong> I contenuti dei colloqui sono protetti da segreto professionale (art. 622 c.p.). Lo psicologo non redige diagnosi né intraprende percorsi di psicoterapia senza ulteriore consenso specifico.</p>
            <p class="q-mb-none"><strong>Validità:</strong> Il consenso è valido per l'intero anno scolastico 2025/2026 ed è revocabile in qualsiasi momento per iscritto.</p>
          </div>

          <div v-if="selectedChild?.has_consent" class="q-pa-md bg-emerald-50 border-emerald-200 rounded-borders text-emerald-900 q-mb-md row items-center">
            <q-icon name="check_circle" size="24px" class="q-mr-sm" />
            <div>
              <div class="text-weight-bold">Consenso Informato regolarmente registrato</div>
              <div class="text-caption">Sottoscritto digitalmente dal genitore. {{ selectedChild?.name }} può accedere liberamente agli slot su richiesta.</div>
            </div>
          </div>

          <div v-else class="q-mb-md">
            <q-checkbox
              v-model="acceptTerms"
              color="primary"
              :label="$t('parentPsychology.termsCheckbox') || 'Dichiaro di aver letto l\'informativa e autorizzo mio/a figlio/a a fruire dei colloqui dello sportello d\'ascolto psicologico.'"
            />
          </div>

          <q-btn
            v-if="!selectedChild?.has_consent"
            color="primary"
            icon="draw"
            :label="$t('parentPsychology.signConsentBtn') || 'Firma e Rilascia Consenso Informato'"
            :loading="loading"
            :disable="!acceptTerms"
            @click="signConsent"
            class="full-width q-py-sm text-weight-bold"
            unelevated
          />
        </q-card>
      </div>

      <!-- Parent Consultation Request -->
      <div class="col-12 col-md-5">
        <q-card flat bordered class="q-pa-lg bg-white rounded-borders shadow-sm full-height">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
            💬 {{ $t('parentPsychology.parentConsultTitle') || 'Richiesta Colloquio Genitore-Psicologo' }}
          </div>
          <div class="text-caption text-slate-500 q-mb-md">
            Lo sportello è aperto anche ai genitori per consulenze pedagogiche e relazionali.
          </div>

          <div class="q-mb-md">
            <label class="text-caption text-weight-bold text-slate-700 block q-mb-xs">Argomento generale</label>
            <q-select
              v-model="parentBookingTopic"
              :options="['Orientamento e motivazione scolastica', 'Difficoltà relazionali o emotive', 'Supporto alla genitorialità', 'Altro (riservato)']"
              outlined
              dense
            />
          </div>

          <div class="q-mb-lg">
            <label class="text-caption text-weight-bold text-slate-700 block q-mb-xs">Fascia oraria preferita</label>
            <q-input
              v-model="parentTimePref"
              outlined
              dense
              placeholder="Es. Martedì pomeriggio dopo le 16:00"
            />
          </div>

          <q-btn
            color="secondary"
            icon="send"
            :label="$t('parentPsychology.sendRequestBtn') || 'Invia Richiesta di Contatto'"
            :loading="loadingReq"
            @click="submitParentRequest"
            class="full-width q-py-sm text-weight-bold"
            unelevated
          />
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
import { psychologyService } from '@/services/psychologyService'

const $q = useQuasar()
const loading = ref(false)
const loadingReq = ref(false)
const acceptTerms = ref(false)
const parentBookingTopic = ref('Orientamento e motivazione scolastica')
const parentTimePref = ref('')

const childrenList = ref([
  { id: 'student-minor-1', name: 'Alessandro Verdi (Classe 2A)', has_consent: false },
  { id: 'student-1', name: 'Marco Rossi (Classe 5A)', has_consent: true }
])
const selectedChild = ref(childrenList.value[0])

async function signConsent() {
  if (!selectedChild.value) return
  loading.value = true
  try {
    await psychologyService.signConsent(selectedChild.value.id, '2025/2026')
    selectedChild.value.has_consent = true
    $q.notify({
      type: 'positive',
      message: `Consenso informato registrato con successo per ${selectedChild.value.name}!`
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message || 'Errore nella registrazione del consenso informato.'
    })
  } finally {
    loading.value = false
  }
}

async function submitParentRequest() {
  loadingReq.value = true
  try {
    await psychologyService.bookSession('psy-01', parentTimePref.value || 'Fascia pomeridiana', 45)
    $q.notify({
      type: 'positive',
      message: 'Richiesta di colloquio genitoriale inoltrata allo psicologo scolastico.'
    })
    parentTimePref.value = ''
  } catch (err) {
    $q.notify({
      type: 'warning',
      message: err.message || 'Richiesta inviata.'
    })
  } finally {
    loadingReq.value = false
  }
}
</script>
