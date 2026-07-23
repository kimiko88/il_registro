<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none">Configurazione Istituto & Feature Flags</h1>
        <div class="text-subtitle2 text-grey-7">Gestisci le regole d'istituto, le autorizzazioni e i flussi di lavoro</div>
      </div>
      <q-space />
      <q-btn
        color="primary"
        icon="save"
        label="Salva Modifiche"
        :loading="saving"
        @click="saveSettings"
      />
    </div>

    <q-card v-if="loading" class="q-pa-lg text-center shadow-2">
      <q-spinner color="primary" size="3em" />
      <div class="q-mt-sm text-grey-7">Caricamento impostazioni in corso...</div>
    </q-card>

    <div v-else class="row q-col-gutter-md">
      <!-- Section 1: Note & Disciplina -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-2 rounded-borders">
          <q-card-section class="bg-primary text-white">
            <div class="text-h6"><q-icon name="gavel" class="q-mr-sm" />Note & Disciplina</div>
            <div class="text-caption opacity-80">Regole di convalida e visibilità delle sanzioni disciplinari</div>
          </q-card-section>
          <q-card-section class="q-pa-md">
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Convalida Note Disciplinari dalla Dirigenza</q-item-label>
                <q-item-label caption>
                  Le note disciplinari richiedono l'approvazione del Dirigente/Vice Dirigente prima di essere visibili a studenti e genitori.
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.require_principal_approval_for_notes" color="secondary" />
              </q-item-section>
            </q-item>
          </q-card-section>
        </q-card>
      </div>

      <!-- Section 2: Visibilità Voti & Valutazioni -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-2 rounded-borders">
          <q-card-section class="bg-indigo-9 text-white">
            <div class="text-h6"><q-icon name="visibility" class="q-mr-sm" />Visibilità Voti & Medie</div>
            <div class="text-caption opacity-80">Permessi di consultazione per genitori e alunni</div>
          </q-card-section>
          <q-card-section class="q-pa-md">
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Mostra Voti ai Genitori</q-item-label>
                <q-item-label caption>
                  Consente ai genitori di visualizzare le singole valutazioni in tempo reale sul registro.
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.allow_parents_view_grades" color="indigo" />
              </q-item-section>
            </q-item>

            <q-separator q-my-xs />

            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Mostra Medie di Classe agli Studenti</q-item-label>
                <q-item-label caption>
                  Abilita il grafico ed il confronto della media individuale con la media di classe.
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.allow_students_view_class_averages" color="indigo" />
              </q-item-section>
            </q-item>
          </q-card-section>
        </q-card>
      </div>

      <!-- Section 3: Scrutini & Flussi -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-2 rounded-borders">
          <q-card-section class="bg-teal-8 text-white">
            <div class="text-h6"><q-icon name="assessment" class="q-mr-sm" />Scrutini & Flussi di Lavoro</div>
            <div class="text-caption opacity-80">Protezioni per le sessioni di scrutinio finale</div>
          </q-card-section>
          <q-card-section class="q-pa-md">
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Blocco Modifica Scrutinio dopo la Validazione</q-item-label>
                <q-item-label caption>
                  Una volta che lo scrutinio è validato dal Dirigente, le valutazioni finali non possono più essere verificate o modificate.
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.lock_scrutiny_editing_after_validation" color="teal" />
              </q-item-section>
            </q-item>
          </q-card-section>
        </q-card>
      </div>

      <!-- Section 4: Sicurezza & Supplenze -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-2 rounded-borders">
          <q-card-section class="bg-deep-orange-8 text-white">
            <div class="text-h6"><q-icon name="security" class="q-mr-sm" />Sicurezza & Supplenze</div>
            <div class="text-caption opacity-80">Autenticazione avanzata e notifiche notifiche per supplenti</div>
          </q-card-section>
          <q-card-section class="q-pa-md">
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Obbligo MFA per il Personale Docente/ATA</q-item-label>
                <q-item-label caption>
                  Richiede l'autenticazione a due fattori (2FA) per tutti i docenti ed il personale di segreteria.
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.require_mfa_for_staff" color="deep-orange" />
              </q-item-section>
            </q-item>

            <q-separator q-my-xs />

            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Notifiche Assegnazione Supplenze</q-item-label>
                <q-item-label caption>
                  Invia una notifica istantanea al docente quando viene assegnato come supplente per una lezione.
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.enable_substitute_notifications" color="deep-orange" />
              </q-item-section>
            </q-item>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import schoolSettingsService from '@/services/schoolSettingsService';

const $q = useQuasar();
const loading = ref(true);
const saving = ref(false);

const settings = ref({
  require_principal_approval_for_notes: false,
  allow_parents_view_grades: true,
  allow_students_view_class_averages: true,
  require_mfa_for_staff: false,
  lock_scrutiny_editing_after_validation: true,
  enable_substitute_notifications: true
});

const fetchSettings = async () => {
  loading.value = true;
  try {
    const res = await schoolSettingsService.getSettings();
    settings.value = res.data;
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel caricamento delle impostazioni: ' + (err.response?.data?.error || err.message) });
  } finally {
    loading.value = false;
  }
};

const saveSettings = async () => {
  saving.value = true;
  try {
    await schoolSettingsService.updateSettings(settings.value);
    $q.notify({ type: 'positive', message: 'Impostazioni aggiornate con successo' });
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio: ' + (err.response?.data?.error || err.message) });
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  fetchSettings();
});
</script>
