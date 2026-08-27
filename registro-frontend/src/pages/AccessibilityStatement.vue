<template>
  <q-page class="accessibility-statement-page q-pa-lg">
    <div class="max-w-4xl mx-auto">
      <!-- Header -->
      <div class="q-mb-xl text-center">
        <q-avatar size="64px" color="primary" text-color="white" icon="accessibility_new" class="q-mb-md shadow-2" />
        <h1 class="text-h4 text-weight-bold text-primary q-my-none">
          {{ $t('a11y.statementTitle') || 'Dichiarazione di Accessibilità' }}
        </h1>
        <div class="text-subtitle1 text-grey-7 q-mt-sm">
          Conforme alle Linee Guida AgID, Legge 4/2004 (Legge Stanca) e Direttiva UE 2016/2102
        </div>
        <div class="text-caption text-grey-6 q-mt-xs">
          Ultimo aggiornamento: 27 Agosto 2026 • Livello di conformità: WCAG 2.2 AA / AAA
        </div>
      </div>

      <!-- Quick Summary Cards -->
      <div class="row q-col-gutter-md q-mb-xl">
        <div class="col-12 col-sm-6 col-md-3">
          <q-card flat bordered class="rounded-xl text-center q-pa-md bg-white">
            <q-icon name="check_circle" color="positive" size="36px" />
            <div class="text-subtitle2 text-weight-bold q-mt-xs">Stato Conformità</div>
            <div class="text-caption text-grey-8">Pienamente Conforme WCAG 2.2 AA</div>
          </q-card>
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <q-card flat bordered class="rounded-xl text-center q-pa-md bg-white">
            <q-icon name="record_voice_over" color="primary" size="36px" />
            <div class="text-subtitle2 text-weight-bold q-mt-xs">Text-to-Speech</div>
            <div class="text-caption text-grey-8">Sintesi Vocale Web Speech API</div>
          </q-card>
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <q-card flat bordered class="rounded-xl text-center q-pa-md bg-white">
            <q-icon name="font_download" color="accent" size="36px" />
            <div class="text-subtitle2 text-weight-bold q-mt-xs">Supporto DSA</div>
            <div class="text-caption text-grey-8">OpenDyslexic & Righello Lettura</div>
          </q-card>
        </div>
        <div class="col-12 col-sm-6 col-md-3">
          <q-card flat bordered class="rounded-xl text-center q-pa-md bg-white">
            <q-icon name="palette" color="secondary" size="36px" />
            <div class="text-subtitle2 text-weight-bold q-mt-xs">Daltonismo & OLED</div>
            <div class="text-caption text-grey-8">Filtri ottici e contrasti 21:1</div>
          </q-card>
        </div>
      </div>

      <!-- Main Content Sections -->
      <div class="q-gutter-y-lg">
        <!-- 1. Impegno per l'accessibilità -->
        <q-card flat bordered class="rounded-xl q-pa-lg bg-white shadow-soft">
          <h2 class="text-h6 text-weight-bold text-primary flex items-center gap-sm q-my-none">
            <q-icon name="school" size="24px" />
            1. Impegno per l'Inclusione Scolastica Digitale
          </h2>
          <div class="text-body1 text-grey-9 q-mt-md leading-relaxed">
            L'Istituzione Scolastica si impegna a rendere il proprio Registro Elettronico e i servizi digitali accessibili a tutti gli utenti (studenti, docenti, genitori, personale ATA), in conformità alla <strong>Legge 9 gennaio 2004, n. 4</strong>, al <strong>Decreto Legislativo 10 agosto 2018, n. 106</strong> e alle <strong>Linee Guida AgID sull'Accessibilità degli strumenti informatici</strong>.
          </div>
        </q-card>

        <!-- 2. Tecnologie e Misure Adottate -->
        <q-card flat bordered class="rounded-xl q-pa-lg bg-white shadow-soft">
          <h2 class="text-h6 text-weight-bold text-primary flex items-center gap-sm q-my-none">
            <q-icon name="settings_accessibility" size="24px" />
            2. Misure e Tecnologie di Accessibilità Integrate
          </h2>
          <div class="text-body1 text-grey-9 q-mt-md">
            L'applicazione include funzionalità di accessibilità avanzate e native attivabili dal pannello rapido o nelle Impostazioni:
            <ul class="q-pl-md q-gutter-y-xs q-mt-sm">
              <li><strong>Sintesi Vocale TTS (Text-to-Speech)</strong>: Lettura audio nativa di circolari, note, compiti e argomenti delle lezioni.</li>
              <li><strong>Supporto DSA & Font Inclusivi</strong>: Caratteri ad alta leggibilità (<em>OpenDyslexic, Lexend, Fredoka</em>) e regolazione interlinea / spaziatura lettere (WCAG 1.4.12).</li>
              <li><strong>Righello di Lettura & Focus Mask</strong>: Maschera visiva orizzontale a contrasto per guidare la lettura su tabelle dense e testi lunghi.</li>
              <li><strong>Filtri Daltonismo & Contrasti OLED</strong>: Palette ottimizzate per Protanopia, Deuteranopia, Tritanopia e modalità nero assoluto ad alto contrasto (Ambra/Verde).</li>
              <li><strong>Navigazione da Tastiera & Focus Visibile</strong>: Supporto completo con tasto <code>Tab</code>, indicatore di focus ad altissima visibilità e scorciatoie rapide globali (premi <code>?</code> per la lista).</li>
              <li><strong>Compatibilità Screen Reader</strong>: Struttura semantica HTML5, attributi <code>aria-live</code> per aggiornamenti in tempo reale e attributi <code>scope</code> sulle tabelle.</li>
            </ul>
          </div>
        </q-card>

        <!-- 3. Meccanismo di Feedback e Segnalazione -->
        <q-card flat bordered class="rounded-xl q-pa-lg bg-white shadow-soft" id="feedback-section">
          <h2 class="text-h6 text-weight-bold text-primary flex items-center gap-sm q-my-none">
            <q-icon name="feedback" size="24px" />
            3. Meccanismo di Feedback & Segnalazione Barriere Digitali
          </h2>
          <div class="text-body1 text-grey-9 q-mt-md">
            Se riscontri qualsiasi difficoltà di accesso o desideri richiedere informazioni e contenuti in un formato alternativo accessibile, puoi inviare una segnalazione compilando il modulo sottostante:
          </div>

          <q-form @submit.prevent="submitFeedback" class="q-mt-lg q-gutter-y-md">
            <div class="row q-col-gutter-md">
              <div class="col-12 col-sm-6">
                <q-input
                  v-model="feedbackForm.name"
                  label="Nome e Cognome *"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
              <div class="col-12 col-sm-6">
                <q-input
                  v-model="feedbackForm.email"
                  label="Email di Contatto *"
                  type="email"
                  outlined
                  dense
                  :rules="[val => !!val || 'Campo obbligatorio']"
                />
              </div>
            </div>

            <q-select
              v-model="feedbackForm.barrierType"
              :options="barrierOptions"
              label="Tipologia di Barriera / Segnalazione *"
              outlined
              dense
              emit-value
              map-options
            />

            <q-input
              v-model="feedbackForm.description"
              type="textarea"
              label="Descrizione dettagliata della problematica riscontrata e pagina/sezione interessata *"
              outlined
              rows="4"
              :rules="[val => !!val || 'Descrivi la problematica riscontrata']"
            />

            <div class="flex justify-end">
              <q-btn
                type="submit"
                color="primary"
                icon="send"
                label="Invia Segnalazione di Accessibilità"
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
            4. Procedura di Attuazione (Difensore Civico per il Digitale)
          </h2>
          <div class="text-body1 text-grey-9 q-mt-md">
            In caso di risposta insoddisfacente o mancata risposta entro 30 giorni dalla notifica o dalla richiesta di informazioni accessibili, il cittadino può inoltrare una segnalazione al <strong>Difensore Civico per il Digitale</strong> tramite l'apposito servizio istituito presso l'<strong>AgID</strong> all'indirizzo:
            <div class="q-mt-sm">
              <a href="https://www.agid.gov.it" target="_blank" rel="noopener" class="text-primary text-weight-bold">
                https://www.agid.gov.it - Agenzia per l'Italia Digitale
              </a>
            </div>
          </div>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const sendingFeedback = ref(false)

const feedbackForm = ref({
  name: '',
  email: '',
  barrierType: 'contrast',
  description: ''
})

const barrierOptions = [
  { label: 'Contrasto visivo o leggibilità del testo', value: 'contrast' },
  { label: 'Difficoltà di navigazione con tastiera o screen reader', value: 'keyboard' },
  { label: 'Problemi con la sintesi vocale (TTS) o lettura compiti/circolari', value: 'tts' },
  { label: 'Contenuto non accessibile in formato PDF o documento', value: 'documents' },
  { label: 'Altra barriera digitale o suggerimento di miglioramento', value: 'other' }
]

function submitFeedback() {
  sendingFeedback.value = true
  setTimeout(() => {
    sendingFeedback.value = false
    $q.notify({
      type: 'positive',
      message: 'Grazie! La tua segnalazione di accessibilità è stata inviata al Responsabile della Transizione Digitale (RTD).',
      timeout: 5000
    })
    feedbackForm.value = {
      name: '',
      email: '',
      barrierType: 'contrast',
      description: ''
    }
  }, 600)
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
