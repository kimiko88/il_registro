<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="settings" color="primary" class="q-mr-sm" />
          {{ t('settings.title') || 'Impostazioni Docente' }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          Personalizza preferenze, sicurezza, lingua e notifiche per il tuo registro elettronico
        </p>
      </div>
      <div>
        <q-btn
          color="primary"
          icon="save"
          label="Salva Tutte le Preferenze"
          unelevated
          class="rounded-lg shadow-sm"
          :loading="savingAll"
          @click="saveAllPreferences"
        />
      </div>
    </div>

    <!-- Navigation Tabs -->
    <q-card class="glass-card shadow-soft border-slate-100 rounded-2xl q-mb-lg overflow-hidden">
      <q-tabs
        v-model="activeTab"
        dense
        active-color="primary"
        indicator-color="primary"
        align="left"
        class="text-slate-600 bg-white border-b border-slate-100"
        narrow-indicator
      >
        <q-tab name="general" icon="language" label="Lingua & Localizzazione" />
        <q-tab name="accessibility" icon="accessibility_new" label="Accessibilità Visiva (DSA & Contrasto)" />
        <q-tab name="security" icon="lock" label="Sicurezza & Password" />
        <q-tab name="notifications" icon="notifications" label="Notifiche & Avvisi" />
        <q-tab name="register" icon="tune" label="Personalizzazione Registro" />
        <q-tab name="signature" icon="draw" label="Firma Digitale & PIN" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated class="bg-transparent">
        
        <!-- Tab 1: Lingua & Localizzazione -->
        <q-tab-panel name="general" class="q-pa-md">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Lingua &amp; Formato Data</div>
          <div class="text-caption text-slate-500 q-mb-md">Configura le opzioni della lingua di interfaccia e visualizzazione temporale</div>

          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">Lingua dell'Applicazione</div>
                <q-select
                  v-model="selectedLocale"
                  :options="localeOptions"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                  @update:model-value="onLocaleChange"
                >
                  <template v-slot:option="scope">
                    <q-item v-bind="scope.itemProps">
                      <q-item-section avatar>
                        <q-avatar size="24px">{{ scope.opt.flag }}</q-avatar>
                      </q-item-section>
                      <q-item-section>
                        <q-item-label>{{ scope.opt.label }}</q-item-label>
                      </q-item-section>
                    </q-item>
                  </template>
                </q-select>

                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">Formato Data Predefinito</div>
                <q-select
                  v-model="generalSettings.dateFormat"
                  :options="[
                    { label: 'DD/MM/YYYY (Es. 08/08/2026)', value: 'DD/MM/YYYY' },
                    { label: 'YYYY-MM-DD (Es. 2026-08-08)', value: 'YYYY-MM-DD' },
                    { label: 'DD.MM.YYYY (Es. 08.08.2026)', value: 'DD.MM.YYYY' }
                  ]"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">Primo Giorno della Settimana</div>
                <q-option-group
                  v-model="generalSettings.firstDayOfWeek"
                  :options="[
                    { label: 'Lunedì', value: 1 },
                    { label: 'Domenica', value: 0 }
                  ]"
                  color="primary"
                  inline
                />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">Fuso Orario &amp; Formato Ora</div>
                <q-select
                  v-model="generalSettings.timeFormat"
                  :options="[
                    { label: '24 Ore (Es. 14:30)', value: '24h' },
                    { label: '12 Ore AM/PM (Es. 2:30 PM)', value: '12h' }
                  ]"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">Fuso Orario Predefinito</div>
                <q-select
                  v-model="generalSettings.timezone"
                  :options="[
                    { label: 'Europe/Rome (UTC+1 / UTC+2)', value: 'Europe/Rome' },
                    { label: 'UTC', value: 'UTC' }
                  ]"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- Tab 2: Sicurezza & Cambio Password -->
        <q-tab-panel name="security" class="q-pa-md">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Sicurezza Account &amp; Modifica Password</div>
          <div class="text-caption text-slate-500 q-mb-md">Aggiorna la tua password di accesso ed imposta la sicurezza del tuo profilo</div>

          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-7">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-md">Modifica Password</div>

                <q-form @submit.prevent="changePassword" class="q-gutter-md">
                  <q-input
                    v-model="pwdForm.currentPassword"
                    :type="showCurrentPwd ? 'text' : 'password'"
                    label="Password Attuale *"
                    outlined
                    dense
                    class="rounded-lg"
                    :rules="[val => !!val || 'Inserisci la password attuale']"
                  >
                    <template v-slot:append>
                      <q-icon
                        :name="showCurrentPwd ? 'visibility_off' : 'visibility'"
                        class="cursor-pointer"
                        @click="showCurrentPwd = !showCurrentPwd"
                      />
                    </template>
                  </q-input>

                  <q-input
                    v-model="pwdForm.newPassword"
                    :type="showNewPwd ? 'text' : 'password'"
                    label="Nuova Password *"
                    outlined
                    dense
                    class="rounded-lg"
                    :rules="[
                      val => !!val || 'Inserisci la nuova password',
                      val => val.length >= 8 || 'La password deve contenere almeno 8 caratteri'
                    ]"
                  >
                    <template v-slot:append>
                      <q-icon
                        :name="showNewPwd ? 'visibility_off' : 'visibility'"
                        class="cursor-pointer"
                        @click="showNewPwd = !showNewPwd"
                      />
                    </template>
                  </q-input>

                  <!-- Indicatori di Forza Password -->
                  <div v-if="pwdForm.newPassword" class="q-mb-xs">
                    <div class="row items-center justify-between text-caption text-grey-7 q-mb-xs">
                      <span>Forza Password:</span>
                      <span class="text-weight-bold" :class="passwordStrengthColorText">{{ passwordStrengthLabel }}</span>
                    </div>
                    <q-linear-progress :value="passwordStrengthValue" :color="passwordStrengthColor" class="rounded-borders" style="height: 6px" />
                  </div>

                  <q-input
                    v-model="pwdForm.confirmPassword"
                    :type="showConfirmPwd ? 'text' : 'password'"
                    label="Conferma Nuova Password *"
                    outlined
                    dense
                    class="rounded-lg"
                    :rules="[
                      val => !!val || 'Conferma la nuova password',
                      val => val === pwdForm.newPassword || 'Le password non coincidono'
                    ]"
                  >
                    <template v-slot:append>
                      <q-icon
                        :name="showConfirmPwd ? 'visibility_off' : 'visibility'"
                        class="cursor-pointer"
                        @click="showConfirmPwd = !showConfirmPwd"
                      />
                    </template>
                  </q-input>

                  <div class="row justify-end q-mt-md">
                    <q-btn
                      type="submit"
                      color="primary"
                      icon="lock_reset"
                      label="Aggiorna Password"
                      unelevated
                      :loading="updatingPassword"
                    />
                  </div>
                </q-form>
              </q-card>
            </div>

            <div class="col-12 col-md-5">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">Autenticazione a Due Fattori (2FA)</div>
                <div class="text-caption text-slate-500 q-mb-md">
                  Proteggi il tuo account docente aggiungendo un livello di sicurezza extra tramite app OTP (Google Authenticator, Authy).
                </div>

                <div class="row items-center justify-between q-pa-sm bg-slate-50 rounded-lg">
                  <div class="row items-center">
                    <q-icon name="phonelink_lock" color="primary" size="28px" class="q-mr-sm" />
                    <div>
                      <div class="text-weight-bold text-slate-800">Stato 2FA</div>
                      <div class="text-caption text-grey-6">{{ mfaEnabled ? 'Attivo' : 'Non Attivo' }}</div>
                    </div>
                  </div>
                  <q-badge :color="mfaEnabled ? 'positive' : 'grey-6'" :label="mfaEnabled ? 'Abilitato' : 'Disabilitato'" />
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- Tab 3: Notifiche & Avvisi -->
        <q-tab-panel name="notifications" class="q-pa-md">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Preferenze Notifiche Docente</div>
          <div class="text-caption text-slate-500 q-mb-md">Seleziona come e quando desideri ricevere comunicazioni ed avvisi dal sistema scolastico</div>

          <q-card flat bordered class="q-pa-md rounded-xl bg-white">
            <q-list separator>
              <q-item tag="label" v-ripple>
                <q-item-section avatar>
                  <q-icon name="email" color="indigo" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Notifiche Email Circolari &amp; Comunicati</q-item-label>
                  <q-item-label caption>Ricevi una copia via email di ogni nuova circolare pubblicata dal Dirigente o dalla Segreteria.</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.emailCirculars" color="indigo" />
                </q-item-section>
              </q-item>

              <q-item tag="label" v-ripple>
                <q-item-section avatar>
                  <q-icon name="swap_horiz" color="amber-9" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Avvisi Sostituzioni Lezioni</q-item-label>
                  <q-item-label caption>Notifica immediata quando ti viene assegnata una sostituzione o supplenza.</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.substituteAlerts" color="amber-9" />
                </q-item-section>
              </q-item>

              <q-item tag="label" v-ripple>
                <q-item-section avatar>
                  <q-icon name="event" color="teal" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Prenotazioni Colloqui Famiglie</q-item-label>
                  <q-item-label caption>Invia notifica quando un genitore prenota o annulla un appuntamento per il ricevimento.</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.parentColloqui" color="teal" />
                </q-item-section>
              </q-item>

              <q-item tag="label" v-ripple>
                <q-item-section avatar>
                  <q-icon name="analytics" color="purple" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Promemoria Scrutini &amp; Consegna Voti</q-item-label>
                  <q-item-label caption>Avvisi automatici in prossimità della chiusura del tabellone voti per il consiglio di classe.</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.scrutinyDeadlines" color="purple" />
                </q-item-section>
              </q-item>

              <q-item tag="label" v-ripple>
                <q-item-section avatar>
                  <q-icon name="volume_up" color="deep-orange" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Suoni di Notifica In-App</q-item-label>
                  <q-item-label caption>Riproduci un breve segnale acustico all'arrivo di nuove notifiche nel registro.</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.soundEnabled" color="deep-orange" />
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </q-tab-panel>

        <!-- Tab 4: Personalizzazione Registro -->
        <q-tab-panel name="register" class="q-pa-md">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Personalizzazione Registro &amp; Griglia Voti</div>
          <div class="text-caption text-slate-500 q-mb-md">Adatta l'aspetto grafico e il comportamento predefinito del registro elettronico</div>

          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">Pagina Iniziale Predefinita (Landing Page)</div>
                <q-select
                  v-model="registerSettings.defaultLanding"
                  :options="[
                    { label: 'Dashboard Generale', value: '/teacher' },
                    { label: 'Le Mie Classi', value: '/teacher/classes' },
                    { label: 'Registro Appello & Presenze', value: '/teacher/attendance' },
                    { label: 'Gestione Voti', value: '/teacher/grades' },
                    { label: 'Agenda & Compiti', value: '/teacher/agenda' }
                  ]"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">Formato Visualizzazione Voti</div>
                <q-select
                  v-model="registerSettings.gradeFormat"
                  :options="[
                    { label: 'Decimali Standard (Es. 7,5 o 7.5)', value: 'decimal' },
                    { label: 'Simboli e Frazioni (Es. 7+, 7½, 8-)', value: 'fractional' },
                    { label: 'Centesimi (Es. 75/100)', value: 'centesimal' }
                  ]"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">Separatore dei Decimali</div>
                <q-select
                  v-model="registerSettings.decimalSeparator"
                  :options="[
                    { label: 'Virgola ( , ) — Standard Italiano (es. 7,5)', value: ',' },
                    { label: 'Punto ( . ) — Standard Internazionale (es. 7.5)', value: '.' }
                  ]"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="p-3 bg-slate-50 rounded-lg border border-slate-200 text-xs">
                  <div class="text-weight-bold text-slate-700 q-mb-xs">Anteprima Visualizzazione Voti Live:</div>
                  <div class="row q-col-gutter-xs">
                    <div class="col-4 text-center">
                      <span class="text-caption text-grey-7">Voto 7.5:</span><br>
                      <q-badge color="indigo" class="text-subtitle2 px-2 py-1">{{ sampleFormattedGrade(7.5) }}</q-badge>
                    </div>
                    <div class="col-4 text-center">
                      <span class="text-caption text-grey-7">Voto 6.25:</span><br>
                      <q-badge color="teal" class="text-subtitle2 px-2 py-1">{{ sampleFormattedGrade(6.25) }}</q-badge>
                    </div>
                    <div class="col-4 text-center">
                      <span class="text-caption text-grey-7">Media 8.42:</span><br>
                      <q-badge color="purple" class="text-subtitle2 px-2 py-1">{{ sampleFormattedAverage(8.4166) }}</q-badge>
                    </div>
                  </div>
                </div>
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">Opzioni Visualizzazione Layout</div>
                <q-list separator>
                  <q-item tag="label" v-ripple>
                    <q-item-section>
                      <q-item-label class="text-weight-bold">Griglia Registro Compatta</q-item-label>
                      <q-item-label caption>Riduci il padding delle tabelle per mostrare più studenti e colonne contemporaneamente.</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                      <q-toggle v-model="registerSettings.compactGrid" color="primary" />
                    </q-item-section>
                  </q-item>

                  <q-item tag="label" v-ripple>
                    <q-item-section>
                      <q-item-label class="text-weight-bold">Mostra Foto Studenti nell'Appello</q-item-label>
                      <q-item-label caption>Visualizza l'avatar o la fototessera dello studente nella schermata di appello.</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                      <q-toggle v-model="registerSettings.showStudentPhotos" color="primary" />
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- Tab 5: Firma Digitale & PIN -->
        <q-tab-panel name="signature" class="q-pa-md">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Firma Digitale &amp; PIN Veloce Lezioni</div>
          <div class="text-caption text-slate-500 q-mb-md">Gestisci il codice PIN per la firma immediata del registro di classe</div>

          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">PIN Veloce Registro (4 Cifre)</div>
                <div class="text-caption text-slate-500 q-mb-md">
                  Il PIN veloce ti consente di firmare e validare l'ora di lezione sul registro di classe senza dover reinserire la password completa.
                </div>

                <div class="row q-col-gutter-sm items-center q-mb-md">
                  <div class="col-8">
                    <q-input
                      v-model="pinValue"
                      type="password"
                      mask="####"
                      maxlength="4"
                      label="PIN Veloce (4 cifre)"
                      outlined
                      dense
                      class="rounded-lg"
                    />
                  </div>
                  <div class="col-4">
                    <q-btn color="primary" label="Salva PIN" unelevated class="full-width" @click="savePin" />
                  </div>
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- Tab Accessibilità & Aspetto -->
        <q-tab-panel name="accessibility" class="q-pa-md">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Accessibilità Visiva &amp; Modalità di Lettura</div>
          <div class="text-caption text-slate-500 q-mb-md">Attiva font ad alta leggibilità (OpenDyslexic) per utenti DSA, contrasto elevato e temi visivi dell'interfaccia.</div>

          <div class="row q-col-gutter-lg">
            <!-- Font OpenDyslexic (DSA) -->
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white full-height shadow-sm">
                <div class="row items-center justify-between q-mb-sm">
                  <div class="row items-center">
                    <q-avatar color="indigo-50" text-color="indigo-700" icon="spellcheck" size="44px" class="q-mr-sm" />
                    <div>
                      <div class="text-subtitle1 text-weight-bold text-slate-800">Font OpenDyslexic (Alta Leggibilità DSA)</div>
                      <div class="text-caption text-slate-500">Applica il font ad alta leggibilità specifico per la dislessia ed i disturbi dell'apprendimento</div>
                    </div>
                  </div>
                  <q-toggle
                    v-model="themeStore.dsaFont"
                    color="indigo"
                    size="lg"
                    @update:model-value="themeStore.toggleDsaFont"
                  />
                </div>
                <q-separator class="q-my-sm" />
                <div class="q-pa-md bg-slate-50 rounded-lg text-slate-700 text-body2 q-mt-sm border border-slate-100" :class="{ 'dsa-font-active': themeStore.dsaFont }">
                  <span class="text-weight-bold">Anteprima Testo:</span> Il registro elettronico garantisce un'esperienza di lettura inclusiva ed accessibile per tutti gli utenti ed alunni con BES/DSA.
                </div>
              </q-card>
            </div>

            <!-- Contrasto Elevato -->
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white full-height shadow-sm">
                <div class="row items-center justify-between q-mb-sm">
                  <div class="row items-center">
                    <q-avatar color="amber-50" text-color="amber-9" icon="contrast" size="44px" class="q-mr-sm" />
                    <div>
                      <div class="text-subtitle1 text-weight-bold text-slate-800">Modalità Contrasto Elevato</div>
                      <div class="text-caption text-slate-500">Aumenta la definizione dei bordi e la nitidezza del testo (WCAG 2.1 AAA)</div>
                    </div>
                  </div>
                  <q-toggle
                    v-model="themeStore.highContrast"
                    color="amber-9"
                    size="lg"
                    @update:model-value="themeStore.toggleHighContrast"
                  />
                </div>
                <q-separator class="q-my-sm" />
                <div class="q-pa-md bg-slate-50 rounded-lg text-slate-700 text-body2 q-mt-sm border border-slate-100" :class="{ 'high-contrast-active': themeStore.highContrast }">
                  <span class="text-weight-bold">Anteprima Contrasto:</span> Pulsanti, card ed evidenziatori avranno bordi netti a 2px per la massima visibilità in ambienti molto illuminati.
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

      </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useThemeStore } from '@/stores/theme'
import { userService } from '@/services/userService'
import { useAuthStore } from '@/stores/auth'

const $q = useQuasar()
const { t, locale } = useI18n()
const themeStore = useThemeStore()
const authStore = useAuthStore()

const activeTab = ref('general')
const savingAll = ref(false)

// Lingua e Localizzazione
const selectedLocale = ref(locale.value || 'it')
const localeOptions = [
  { label: 'Italiano', value: 'it', flag: '🇮🇹' },
  { label: 'English', value: 'en', flag: '🇬🇧' },
  { label: 'Español', value: 'es', flag: '🇪🇸' },
  { label: 'Français', value: 'fr', flag: '🇫🇷' },
  { label: 'Deutsch', value: 'de', flag: '🇩🇪' },
  { label: 'Русский', value: 'ru', flag: '🇷🇺' },
  { label: 'Українська', value: 'uk', flag: '🇺🇦' },
  { label: 'العربية', value: 'ar', flag: '🇸🇦' },
  { label: '中文 (简体)', value: 'zh', flag: '🇨🇳' }
]

const generalSettings = reactive({
  dateFormat: 'DD/MM/YYYY',
  firstDayOfWeek: 1,
  timeFormat: '24h',
  timezone: 'Europe/Rome'
})

const onLocaleChange = (newLoc) => {
  if (newLoc) {
    locale.value = newLoc
    localStorage.setItem('user_locale', newLoc)
    $q.notify({
      type: 'positive',
      message: `Lingua impostata su: ${newLoc.toUpperCase()}`,
      position: 'top',
      timeout: 1500
    })
  }
}

// Password Form
const pwdForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const showCurrentPwd = ref(false)
const showNewPwd = ref(false)
const showConfirmPwd = ref(false)
const updatingPassword = ref(false)
const mfaEnabled = ref(false)

const passwordStrengthValue = computed(() => {
  const p = pwdForm.newPassword
  if (!p) return 0
  let score = 0
  if (p.length >= 8) score += 0.25
  if (/[A-Z]/.test(p)) score += 0.25
  if (/[0-9]/.test(p)) score += 0.25
  if (/[^A-Za-z0-9]/.test(p)) score += 0.25
  return score
})

const passwordStrengthLabel = computed(() => {
  const v = passwordStrengthValue.value
  if (v <= 0.25) return 'Debole'
  if (v <= 0.5) return 'Media'
  if (v <= 0.75) return 'Buona'
  return 'Forte'
})

const passwordStrengthColor = computed(() => {
  const v = passwordStrengthValue.value
  if (v <= 0.25) return 'negative'
  if (v <= 0.5) return 'warning'
  if (v <= 0.75) return 'info'
  return 'positive'
})

const passwordStrengthColorText = computed(() => {
  const v = passwordStrengthValue.value
  if (v <= 0.25) return 'text-negative'
  if (v <= 0.5) return 'text-warning'
  if (v <= 0.75) return 'text-info'
  return 'text-positive'
})

const changePassword = async () => {
  if (pwdForm.newPassword !== pwdForm.confirmPassword) {
    $q.notify({ type: 'negative', message: 'Le password non coincidono' })
    return
  }
  const userId = authStore.user?.id
  if (!userId) {
    $q.notify({ type: 'negative', message: 'Sessione non valida, effettua nuovamente il login' })
    return
  }
  updatingPassword.value = true
  try {
    await userService.changePassword(userId, pwdForm.currentPassword, pwdForm.newPassword)
    $q.notify({ type: 'positive', message: 'Password aggiornata con successo!' })
    pwdForm.currentPassword = ''
    pwdForm.newPassword = ''
    pwdForm.confirmPassword = ''
  } catch (err) {
    const errorMsg = err.response?.data?.error || err.response?.data?.message || 'Errore durante la modifica della password'
    $q.notify({ type: 'negative', message: errorMsg })
  } finally {
    updatingPassword.value = false
  }
}

// Notifications
const notificationSettings = reactive({
  emailCirculars: true,
  substituteAlerts: true,
  parentColloqui: true,
  scrutinyDeadlines: true,
  soundEnabled: false
})

// Register preferences
const registerSettings = reactive({
  defaultLanding: '/teacher',
  gradeFormat: 'decimal',
  decimalSeparator: ',',
  compactGrid: false,
  showStudentPhotos: true
})

const sampleFormattedGrade = (val) => {
  const sep = registerSettings.decimalSeparator || ','
  const fmt = registerSettings.gradeFormat || 'decimal'
  if (fmt === 'fractional') {
    if (val === 7.5) return '7½'
    if (val === 6.25) return '6+'
    return String(val)
  }
  if (fmt === 'centesimal') {
    return `${(val * 10).toFixed(0)}/100`
  }
  return String(val).replace('.', sep)
}

const sampleFormattedAverage = (val) => {
  const sep = registerSettings.decimalSeparator || ','
  return Number(val).toFixed(2).replace('.', sep)
}

// PIN
const pinValue = ref('')

const savePin = () => {
  if (!pinValue.value || pinValue.value.length !== 4) {
    $q.notify({ type: 'warning', message: 'Il PIN deve essere di esattamente 4 cifre' })
    return
  }
  localStorage.setItem('teacher_quick_pin', pinValue.value)
  $q.notify({ type: 'positive', message: 'PIN Veloce salvato correttamente!' })
}

// Load saved settings
onMounted(() => {
  const savedLocale = localStorage.getItem('user_locale')
  if (savedLocale) {
    selectedLocale.value = savedLocale
    locale.value = savedLocale
  }

  const savedNotif = localStorage.getItem('teacher_notification_settings')
  if (savedNotif) {
    try { Object.assign(notificationSettings, JSON.parse(savedNotif)) } catch (e) { console.warn('Could not parse notification settings', e) }
  }

  const savedReg = localStorage.getItem('teacher_register_settings')
  if (savedReg) {
    try { Object.assign(registerSettings, JSON.parse(savedReg)) } catch (e) { console.warn('Could not parse register settings', e) }
  }
  const savedSep = localStorage.getItem('user_decimal_separator')
  if (savedSep) {
    registerSettings.decimalSeparator = savedSep
  }

  const savedPin = localStorage.getItem('teacher_quick_pin')
  if (savedPin) pinValue.value = savedPin
})

const saveAllPreferences = () => {
  savingAll.value = true
  try {
    localStorage.setItem('teacher_notification_settings', JSON.stringify(notificationSettings))
    localStorage.setItem('teacher_register_settings', JSON.stringify(registerSettings))
    localStorage.setItem('teacher_general_settings', JSON.stringify(generalSettings))
    localStorage.setItem('user_decimal_separator', registerSettings.decimalSeparator || ',')
    $q.notify({
      type: 'positive',
      message: 'Tutte le preferenze sono state salvate con successo!',
      position: 'top',
      timeout: 2000
    })
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio delle impostazioni' })
  } finally {
    savingAll.value = false
  }
}
</script>

<style scoped>
.glass-card {
  backdrop-filter: blur(12px);
  background: rgba(255, 255, 255, 0.95);
}
.shadow-soft {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}
</style>
