<template>
  <q-page class="accessibility-statement-page q-pa-lg">
    <div class="max-w-4xl mx-auto">
      <!-- Header -->
      <div class="q-mb-xl text-center">
        <q-avatar size="64px" color="primary" text-color="white" icon="accessibility_new" class="q-mb-md shadow-2" />
        <h1 class="text-h4 text-weight-bold text-primary q-my-none">
          {{ t('a11y.statementTitle') || 'Dichiarazione di Accessibilità' }}
        </h1>
        <div class="text-subtitle1 text-grey-7 q-mt-sm">
          {{ t('a11y.statementSubtitle') || 'Conforme alle Linee Guida AgID, Legge 4/2004 (Legge Stanca) e Direttiva UE 2016/2102' }}
        </div>
        <div class="text-caption text-grey-6 q-mt-xs">
          {{ t('a11y.statementUpdate') || 'Ultimo aggiornamento: 27 Agosto 2026 • Livello di conformità: WCAG 2.2 AA / AAA' }}
        </div>
      </div>

      <!-- Quick Summary Cards -->
      <div class="row q-col-gutter-md q-mb-xl">
        <div class="col-12 col-sm-6 col-md-3">
          <q-card flat bordered class="rounded-xl text-center q-pa-md bg-white">
            <q-icon name="check_circle" color="positive" size="36px" />
            <div class="text-subtitle2 text-weight-bold q-mt-xs">{{ t('a11y.complianceStatus') || 'Stato Conformità' }}</div>
            <div class="text-caption text-grey-8">{{ t('a11y.complianceStatusVal') || 'Pienamente Conforme WCAG 2.2 AA' }}</div>
          </q-card>
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <q-card flat bordered class="rounded-xl text-center q-pa-md bg-white">
            <q-icon name="record_voice_over" color="primary" size="36px" />
            <div class="text-subtitle2 text-weight-bold q-mt-xs">{{ t('a11y.ttsCardTitle') || 'Text-to-Speech' }}</div>
            <div class="text-caption text-grey-8">{{ t('a11y.ttsCardVal') || 'Sintesi Vocale Web Speech API' }}</div>
          </q-card>
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <q-card flat bordered class="rounded-xl text-center q-pa-md bg-white">
            <q-icon name="font_download" color="accent" size="36px" />
            <div class="text-subtitle2 text-weight-bold q-mt-xs">{{ t('a11y.dsaSupportTitle') || 'Supporto DSA' }}</div>
            <div class="text-caption text-grey-8">{{ t('a11y.dsaSupportVal') || 'OpenDyslexic & Righello Lettura' }}</div>
          </q-card>
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <q-card flat bordered class="rounded-xl text-center q-pa-md bg-white">
            <q-icon name="palette" color="secondary" size="36px" />
            <div class="text-subtitle2 text-weight-bold q-mt-xs">{{ t('a11y.colorblindTitle') || 'Daltonismo & OLED' }}</div>
            <div class="text-caption text-grey-8">{{ t('a11y.colorblindVal') || 'Filtri ottici e contrasti 21:1' }}</div>
          </q-card>
        </div>
      </div>

      <!-- Main Content Sections -->
      <div class="q-gutter-y-lg">
        <!-- 1. Impegno per l'accessibilità -->
        <q-card flat bordered class="rounded-xl q-pa-lg bg-white shadow-soft">
          <h2 class="text-h6 text-weight-bold text-primary flex items-center gap-sm q-my-none">
            <q-icon name="school" size="24px" />
            {{ t('a11y.sec1Title') || "1. Impegno per l'Inclusione Scolastica Digitale" }}
          </h2>
          <div class="text-body1 text-grey-9 q-mt-md leading-relaxed">
            {{ t('a11y.sec1Text') || "L'Istituzione Scolastica si impegna a rendere il proprio Registro Elettronico e i servizi digitali accessibili a tutti gli utenti (studenti, docenti, genitori, personale ATA), in conformità alla Legge 9 gennaio 2004, n. 4, al Decreto Legislativo 10 agosto 2018, n. 106 e alle Linee Guida AgID sull'Accessibilità degli strumenti informatici." }}
          </div>
        </q-card>

        <!-- 2. Tecnologie e Misure Adottate -->
        <q-card flat bordered class="rounded-xl q-pa-lg bg-white shadow-soft">
          <h2 class="text-h6 text-weight-bold text-primary flex items-center gap-sm q-my-none">
            <q-icon name="settings_accessibility" size="24px" />
            {{ t('a11y.sec2Title') || '2. Misure e Tecnologie di Accessibilità Integrate' }}
          </h2>
          <div class="text-body1 text-grey-9 q-mt-md">
            {{ t('a11y.sec2Intro') || "L'applicazione include funzionalità di accessibilità avanzate e native attivabili dal pannello rapido o nelle Impostazioni:" }}
            <ul class="q-pl-md q-gutter-y-xs q-mt-sm">
              <li>{{ t('a11y.sec2Tts') }}</li>
              <li>{{ t('a11y.sec2Dsa') }}</li>
              <li>{{ t('a11y.sec2Ruler') }}</li>
              <li>{{ t('a11y.sec2Color') }}</li>
              <li>{{ t('a11y.sec2Keyboard') }}</li>
              <li>{{ t('a11y.sec2ScreenReader') }}</li>
            </ul>
          </div>
        </q-card>

        <!-- 3. Meccanismo di Feedback e Segnalazione -->
        <q-card flat bordered class="rounded-xl q-pa-lg bg-white shadow-soft" id="feedback-section">
          <h2 class="text-h6 text-weight-bold text-primary flex items-center gap-sm q-my-none">
            <q-icon name="feedback" size="24px" />
            {{ t('a11y.sec3Title') || '3. Meccanismo di Feedback & Segnalazione Barriere Digitali' }}
          </h2>
          <div class="text-body1 text-grey-9 q-mt-md">
            {{ t('a11y.sec3Intro') || 'Se riscontri qualsiasi difficoltà di accesso o desideri richiedere informazioni e contenuti in un formato alternativo accessibile, puoi inviare una segnalazione compilando il modulo sottostante:' }}
            
            <!-- Feedback Confirmation Banner -->
            <q-banner v-if="submittedProtocol" rounded class="bg-positive text-white q-mt-md shadow-1">
              <template v-slot:avatar>
                <q-icon name="verified" size="32px" />
              </template>
              <div class="text-subtitle1 text-weight-bold">{{ t('a11y.feedbackSuccessTitle') || 'Segnalazione Registrata Correttamente' }}</div>
              <div>
                {{ t('a11y.feedbackSuccessText', { protocol: submittedProtocol }) || `La tua segnalazione è stata protocollata con codice ${submittedProtocol} ed è stata presa in carico dal Responsabile della Transizione Digitale (RTD). Riceverai riscontro all'indirizzo email indicato entro i termini di legge (30 giorni).` }}
              </div>
              <template v-slot:action>
                <q-btn flat color="white" :label="t('a11y.newFeedbackBtn') || 'Nuova Segnalazione'" @click="submittedProtocol = null" />
              </template>
            </q-banner>
          </div>

          <q-form @submit="submitFeedback" class="q-mt-lg q-gutter-y-md">
            <div class="row q-col-gutter-md">
              <div class="col-12 col-sm-6">
                <q-input
                  v-model="feedbackForm.name"
                  :label="t('a11y.fullNameLabel') || 'Nome e Cognome *'"
                  outlined
                  dense
                  :rules="[val => !!val || t('a11y.requiredField') || 'Campo obbligatorio']"
                />
              </div>
              <div class="col-12 col-sm-6">
                <q-input
                  v-model="feedbackForm.email"
                  :label="t('a11y.emailLabel') || 'Email di Contatto *'"
                  type="email"
                  outlined
                  dense
                  :rules="[val => !!val || t('a11y.requiredField') || 'Campo obbligatorio']"
                />
              </div>
            </div>

            <q-select
              v-model="feedbackForm.barrierType"
              :options="barrierOptions"
              :label="t('a11y.barrierTypeLabel') || 'Tipologia di Barriera / Segnalazione *'"
              outlined
              dense
              emit-value
              map-options
            />

            <q-input
              v-model="feedbackForm.description"
              type="textarea"
              :label="t('a11y.descLabel') || 'Descrizione dettagliata della problematica riscontrata e pagina/sezione interessata *'"
              outlined
              rows="4"
              :rules="[val => !!val || t('a11y.describeIssue') || 'Descrivi la problematica riscontrata']"
            />

            <div class="flex justify-end">
              <q-btn
                type="submit"
                color="primary"
                icon="send"
                :label="t('a11y.submitBtn') || 'Invia Segnalazione di Accessibilità'"
                :loading="sendingFeedback"
                unelevated
                no-caps
                rounded
                class="q-px-lg"
              />
            </div>
          </q-form>
        </q-card>

        <!-- 4. Procedura di Attuazione AgID -->
        <q-card flat bordered class="rounded-xl q-pa-lg bg-white shadow-soft">
          <h2 class="text-h6 text-weight-bold text-primary flex items-center gap-sm q-my-none">
            <q-icon name="gavel" size="24px" />
            {{ t('a11y.sec4Title') || '4. Procedura di Attuazione (Difensore Civico per il Digitale)' }}
          </h2>
          <div class="text-body1 text-grey-9 q-mt-md">
            {{ t('a11y.sec4Text') || "In caso di risposta insoddisfacente o mancata risposta entro 30 giorni dalla notifica o dalla richiesta di informazioni accessibili, il cittadino può inoltrare una segnalazione al Difensore Civico per il Digitale tramite l'apposito servizio istituito presso l'AgID all'indirizzo:" }}
            <div class="q-mt-sm">
              <a href="https://www.agid.gov.it" target="_blank" rel="noopener" class="text-primary text-weight-bold">
                {{ t('a11y.agidLinkText') || "https://www.agid.gov.it - Agenzia per l'Italia Digitale" }}
              </a>
            </div>
          </div>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { accessibilityService } from '@/services/accessibilityService'

const $q = useQuasar()
const { t } = useI18n()
const sendingFeedback = ref(false)
const submittedProtocol = ref(null)

const feedbackForm = ref({
  name: '',
  email: '',
  barrierType: 'contrast',
  description: ''
})

const barrierOptions = computed(() => [
  { label: t('a11y.barrierContrast') || 'Contrasto visivo o leggibilità del testo', value: 'contrast' },
  { label: t('a11y.barrierKeyboard') || 'Difficoltà di navigazione con tastiera o screen reader', value: 'keyboard' },
  { label: t('a11y.barrierTts') || 'Problemi con la sintesi vocale (TTS) o lettura compiti/circolari', value: 'tts' },
  { label: t('a11y.barrierDocuments') || 'Contenuto non accessibile in formato PDF o documento', value: 'documents' },
  { label: t('a11y.barrierOther') || 'Altra barriera digitale o suggerimento di miglioramento', value: 'other' }
])

async function submitFeedback() {
  if (!feedbackForm.value.name || !feedbackForm.value.email || !feedbackForm.value.description) {
    $q.notify({
      type: 'warning',
      message: t('a11y.fillRequired') || 'Compila tutti i campi obbligatori'
    })
    return
  }

  sendingFeedback.value = true
  try {
    const res = await accessibilityService.submitFeedback(feedbackForm.value)
    submittedProtocol.value = res.protocol_number || res.id
    $q.notify({
      type: 'positive',
      icon: 'check_circle',
      message: t('a11y.feedbackSentProtocol', { protocol: submittedProtocol.value }) || `Grazie! Segnalazione inviata con successo. Protocollo: ${submittedProtocol.value}`,
      timeout: 7000
    })
    feedbackForm.value = {
      name: '',
      email: '',
      barrierType: 'contrast',
      description: ''
    }
  } catch (error) {
    const errorMsg = error.response?.data?.error || error.message || t('a11y.feedbackError') || 'Errore durante l\'invio della segnalazione'
    $q.notify({
      type: 'negative',
      icon: 'error',
      message: errorMsg,
      timeout: 6000
    })
  } finally {
    sendingFeedback.value = false
  }
}
</script>

<style scoped>
.accessibility-statement-page {
  background: #f8fafc;
  min-height: 100vh;
}
.max-w-4xl {
  max-width: 900px;
}
.shadow-soft {
  box-shadow: 0 4px 20px -2px rgba(0, 0, 0, 0.05);
}
.leading-relaxed {
  line-height: 1.7;
}
</style>
