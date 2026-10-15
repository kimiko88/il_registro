<template>
  <q-page class="q-pa-md q-pa-lg-xl bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <div class="row items-center q-gutter-x-sm">
          <q-avatar color="pink-1" text-color="pink-8" icon="psychology" size="44px" />
          <div>
            <h1 class="text-h5 text-weight-bolder text-slate-900 q-my-none">
              {{ $t('psychologyDesk.title') || 'Sportello d\'Ascolto & Consulenza Psicologica (CIC)' }}
            </h1>
            <div class="text-caption text-slate-500">
              {{ $t('psychologyDesk.subtitle') || 'Servizio di supporto psicologico scolastico anonimo e riservato (L. 56/1989 & MIM-CNOP)' }}
            </div>
          </div>
        </div>
      </div>

      <div class="row items-center q-gutter-x-sm">
        <q-chip outline color="pink-8" icon="verified_user">
          {{ hasConsent ? ($t('psychologyDesk.consentOk') || 'Consenso Genitoriale Attivo') : ($t('psychologyDesk.consentPending') || 'Consenso in Attesa') }}
        </q-chip>
      </div>
    </div>

    <!-- Legal & Confidentiality Guarantee Banner -->
    <q-card flat bordered class="q-pa-md bg-white border-pink-100 rounded-borders q-mb-lg shadow-sm">
      <div class="row items-center q-gutter-x-md">
        <q-icon name="lock" size="32px" color="pink-7" />
        <div class="col">
          <div class="text-subtitle2 text-weight-bold text-pink-9">
            {{ $t('psychologyDesk.privacyTitle') || 'Tutela Assoluta della Riservatezza & Segreto Professionale' }}
          </div>
          <div class="text-caption text-slate-600">
            {{ $t('psychologyDesk.privacyDesc') || 'I colloqui con lo psicologo scolastico sono coperti da segreto professionale ai sensi dell\'art. 622 c.p. e della Legge 56/1989. Nessun docente, compagno o terzo ha accesso ai contenuti discussi.' }}
          </div>
        </div>
      </div>
    </q-card>

    <!-- Consent Warning (if not signed by parents) -->
    <q-card v-if="!hasConsent" flat bordered class="q-pa-md bg-amber-50 border-amber-200 rounded-borders q-mb-lg">
      <div class="row items-center q-gutter-x-md">
        <q-icon name="warning" size="28px" color="amber-9" />
        <div class="col">
          <div class="text-subtitle2 text-weight-bold text-amber-9">
            {{ $t('psychologyDesk.consentRequiredTitle') || 'Consenso Informato Genitoriale Obbligatorio' }}
          </div>
          <div class="text-caption text-amber-8">
            {{ $t('psychologyDesk.consentRequiredDesc') || 'Per accedere ai colloqui individuali con lo psicologo, per gli studenti minorenni è necessaria la firma preventiva del consenso informato da parte dei genitori sul portale famiglie.' }}
          </div>
        </div>
        <q-btn flat dense color="amber-10" icon="info" :label="$t('psychologyDesk.howToSign') || 'Come fare'" @click="showInfoDialog = true" />
      </div>
    </q-card>

    <div class="row q-col-gutter-lg">
      <!-- Booking Form -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="q-pa-lg bg-white rounded-borders shadow-sm full-height">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
            📅 {{ $t('psychologyDesk.bookTitle') || 'Prenota un Colloquio Riservato' }}
          </div>
          <div class="text-caption text-slate-500 q-mb-md">
            {{ $t('psychologyDesk.bookSubtitle') || 'Scegli uno slot orario tra quelli disponibili. Riceverai un codice anonimo di conferma.' }}
          </div>

          <div class="q-mb-md">
            <label class="text-caption text-weight-bold text-slate-700 block q-mb-xs">
              {{ $t('psychologyDesk.selectSlot') || 'Slot Orario Disponibile' }}
            </label>
            <q-select
              v-model="bookingForm.slotTime"
              :options="availableSlots"
              outlined
              dense
              icon="schedule"
              :disable="!hasConsent"
              :placeholder="$t('psychologyDesk.placeholderSlot') || 'Seleziona giorno e orario'"
            />
          </div>

          <div class="q-mb-md">
            <label class="text-caption text-weight-bold text-slate-700 block q-mb-xs">
              {{ $t('psychologyDesk.topic') || 'Area Generale del Colloquio (Opzionale & Anonima)' }}
            </label>
            <q-select
              v-model="bookingForm.topic"
              :options="topicOptions"
              outlined
              dense
              :disable="!hasConsent"
            />
          </div>

          <div class="q-mb-lg">
            <label class="text-caption text-weight-bold text-slate-700 block q-mb-xs">
              {{ $t('psychologyDesk.notesOptional') || 'Eventuali note o preferenze orarie (facoltative)' }}
            </label>
            <q-input
              v-model="bookingForm.notes"
              type="textarea"
              rows="2"
              outlined
              dense
              :disable="!hasConsent"
              placeholder="Es. Preferisco un orario durante l'ora di buco o a ricreazione"
            />
          </div>

          <q-btn
            color="primary"
            icon="event_available"
            :label="$t('psychologyDesk.submitBookingBtn') || 'Conferma Prenotazione Anonima'"
            :loading="loading"
            :disable="!hasConsent || !bookingForm.slotTime"
            @click="submitBooking"
            class="full-width q-py-sm text-weight-bold"
            unelevated
          />
        </q-card>
      </div>

      <!-- My Active Bookings -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="q-pa-lg bg-white rounded-borders shadow-sm full-height">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
            📋 {{ $t('psychologyDesk.myBookingsTitle') || 'I Miei Appuntamenti Prenotati' }}
          </div>
          <div class="text-caption text-slate-500 q-mb-md">
            {{ $t('psychologyDesk.myBookingsSubtitle') || 'Consulta gli appuntamenti confermati con lo psicologo d\'istituto.' }}
          </div>

          <div v-if="myBookings.length === 0" class="text-center q-pa-xl text-slate-400">
            <q-icon name="event_busy" size="48px" />
            <div class="q-mt-sm">{{ $t('psychologyDesk.noBookings') || 'Nessun colloquio programmato al momento.' }}</div>
          </div>

          <q-list v-else separator class="rounded-borders border border-slate-100">
            <q-item v-for="(b, idx) in myBookings" :key="idx" class="q-py-md">
              <q-item-section avatar>
                <q-avatar color="pink-1" text-color="pink-8" icon="calendar_today" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold text-slate-900">{{ b.slotTime }}</q-item-label>
                <q-item-label caption class="text-slate-500">
                  Codice Anonimo: <span class="font-mono text-weight-bold text-primary">{{ b.code }}</span> • {{ b.topic }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-badge color="positive" :label="$t('psychologyDesk.confirmed') || 'Confermato'" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>
    </div>

    <!-- Info Dialog -->
    <q-dialog v-model="showInfoDialog">
      <q-card style="min-width: 360px">
        <q-card-section class="row items-center">
          <div class="text-h6 text-weight-bold">Consenso Informato Genitoriale</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="text-body2 text-slate-700">
          Ai sensi della normativa scolastica e dell'accordo MIM-CNOP, prima del primo colloquio con lo psicologo scolastico è obbligatorio che un genitore o tutore legale acceda alla sezione <strong>Sportello Psicologico (CIC)</strong> del proprio Registro Famiglie e sottoscriva digitalmente il consenso informato per l'anno scolastico in corso.
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Ho Capito" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useQuasar } from 'quasar'
import { psychologyService } from '@/services/psychologyService'

const $q = useQuasar()
const loading = ref(false)
const showInfoDialog = ref(false)

// In a real app this reflects student's consent status from the backend
const hasConsent = ref(true)

const availableSlots = [
  'Lunedì 12 Ottobre 2026 - 10:00 (Stanza CIC)',
  'Lunedì 12 Ottobre 2026 - 11:30 (Stanza CIC)',
  'Mercoledì 14 Ottobre 2026 - 09:30 (Stanza CIC)',
  'Mercoledì 14 Ottobre 2026 - 11:00 (Stanza CIC)',
  'Venerdì 16 Ottobre 2026 - 10:30 (Stanza CIC)'
]

const topicOptions = [
  'Benessere scolastico e motivazione',
  'Gestione dell\'ansia e stress da verifiche',
  'Relazioni con i compagni e gruppo classe',
  'Orientamento al futuro e scelte di studio',
  'Altro (riservato)'
]

const bookingForm = reactive({
  slotTime: availableSlots[0],
  topic: topicOptions[0],
  notes: ''
})

const myBookings = ref([
  {
    slotTime: 'Mercoledì 14 Ottobre 2026 - 09:30',
    topic: 'Benessere scolastico e motivazione',
    code: 'CIC-2026-A841'
  }
])

async function submitBooking() {
  loading.value = true
  try {
    await psychologyService.bookSession('psy-01', bookingForm.slotTime, 45)
    const newCode = `CIC-2026-${Math.floor(1000 + Math.random() * 9000)}`
    myBookings.value.push({
      slotTime: bookingForm.slotTime,
      topic: bookingForm.topic,
      code: newCode
    })
    $q.notify({
      type: 'positive',
      message: `Colloquio prenotato con successo! Codice anonimo: ${newCode}`
    })
  } catch (err) {
    $q.notify({
      type: 'warning',
      message: err.message || 'Prenotazione registrata in modalità riservata.'
    })
  } finally {
    loading.value = false
  }
}
</script>
