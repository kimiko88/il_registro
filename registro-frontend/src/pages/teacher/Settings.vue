<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="settings" color="primary" class="q-mr-sm" />
          {{ t('settingsPage.teacherTitle') || t('settings.title') || 'Impostazioni Docente' }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ t('settingsPage.teacherSubtitle') || 'Personalizza preferenze, sicurezza, lingua e notifiche per il tuo registro elettronico' }}
        </p>
      </div>
      <div>
        <q-btn
          color="primary"
          icon="save"
          :label="t('settingsPage.saveAll') || 'Salva Tutte le Preferenze'"
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
        <q-tab name="general" icon="language" :label="t('settingsPage.tabGeneral') || 'Lingua & Localizzazione'" />
        <q-tab name="accessibility" icon="accessibility_new" :label="t('settingsPage.tabAccessibility') || 'Accessibilità Visiva (DSA & Contrasto)'" />
        <q-tab name="security" icon="lock" :label="t('settingsPage.tabSecurity') || 'Sicurezza & Password'" />
        <q-tab name="notifications" icon="notifications" :label="t('settingsPage.tabNotifications') || 'Notifiche & Avvisi'" />
        <q-tab name="register" icon="tune" :label="t('settingsPage.tabRegister') || 'Personalizzazione Registro'" />
        <q-tab name="signature" icon="draw" :label="t('settingsPage.tabSignature') || 'Firma Digitale & PIN'" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated class="bg-transparent">
        
        <!-- Tab 1: Lingua & Localizzazione -->
        <q-tab-panel name="general" class="q-pa-md">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">{{ t('settingsPage.langAndDateFormat') || 'Lingua & Formato Data' }}</div>
          <div class="text-caption text-slate-500 q-mb-md">{{ t('settingsPage.langAndDateSub') || 'Configura le opzioni della lingua di interfaccia e visualizzazione temporale' }}</div>

          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">{{ t('settingsPage.appLanguage') || "Lingua dell'Applicazione" }}</div>
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

                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">{{ t('settingsPage.defaultDateFormat') || 'Formato Data Predefinito' }}</div>
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

                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">{{ t('settingsPage.firstDayOfWeek') || 'Primo Giorno della Settimana' }}</div>
                <q-option-group
                  v-model="generalSettings.firstDayOfWeek"
                  :options="firstDayOfWeekOptions"
                  color="primary"
                  inline
                />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">{{ t('settingsPage.timeFormatAndTimezone') || 'Fuso Orario & Formato Ora' }}</div>
                <q-select
                  v-model="generalSettings.timeFormat"
                  :options="timeFormatOptions"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="text-subtitle2 text-weight-bold q-mb-sm text-slate-700">{{ t('settingsPage.defaultTimezone') || 'Fuso Orario Predefinito' }}</div>
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
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">{{ t('settingsPage.accountSecurityAndPwd') || 'Sicurezza Account & Modifica Password' }}</div>
          <div class="text-caption text-slate-500 q-mb-md">{{ t('settingsPage.accountSecuritySub') || 'Aggiorna la tua password di accesso ed imposta la sicurezza del tuo profilo' }}</div>

          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-7">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-md">{{ t('settingsPage.editPassword') || 'Modifica Password' }}</div>

                <q-form @submit.prevent="changePassword" class="q-gutter-md">
                  <q-input
                    v-model="pwdForm.currentPassword"
                    :type="showCurrentPwd ? 'text' : 'password'"
                    :label="t('settingsPage.currentPasswordLabel') || 'Password Attuale *'"
                    outlined
                    dense
                    class="rounded-lg"
                    :rules="[val => !!val || (t('settingsPage.currentPasswordReq') || 'Inserisci la password attuale')]"
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
                    :label="t('settingsPage.newPasswordLabel') || 'Nuova Password *'"
                    outlined
                    dense
                    class="rounded-lg"
                    :rules="[
                      val => !!val || (t('settingsPage.newPasswordReq') || 'Inserisci la nuova password'),
                      val => val.length >= 8 || (t('settingsPage.passwordMin8') || 'La password deve contenere almeno 8 caratteri')
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
                      <span>{{ t('settingsPage.passwordStrength') || 'Forza Password:' }}</span>
                      <span class="text-weight-bold" :class="passwordStrengthColorText">{{ passwordStrengthLabel }}</span>
                    </div>
                    <q-linear-progress :value="passwordStrengthValue" :color="passwordStrengthColor" class="rounded-borders" style="height: 6px" />
                  </div>

                  <q-input
                    v-model="pwdForm.confirmPassword"
                    :type="showConfirmPwd ? 'text' : 'password'"
                    :label="t('settingsPage.confirmPasswordLabel') || 'Conferma Nuova Password *'"
                    outlined
                    dense
                    class="rounded-lg"
                    :rules="[
                      val => !!val || (t('settingsPage.confirmPasswordReq') || 'Conferma la nuova password'),
                      val => val === pwdForm.newPassword || (t('settingsPage.passwordMismatch') || 'Le password non coincidono')
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
                      :label="t('settingsPage.updatePasswordBtn') || 'Aggiorna Password'"
                      unelevated
                      :loading="updatingPassword"
                    />
                  </div>
                </q-form>
              </q-card>
            </div>

            <div class="col-12 col-md-5">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">{{ t('settingsPage.twoFactorTitle') || 'Autenticazione a Due Fattori (2FA)' }}</div>
                <div class="text-caption text-slate-500 q-mb-md">
                  {{ t('settingsPage.twoFactorDesc') || 'Proteggi il tuo account docente aggiungendo un livello di sicurezza extra tramite app OTP (Google Authenticator, Authy).' }}
                </div>

                <div class="row items-center justify-between q-pa-sm bg-slate-50 rounded-lg">
                  <div class="row items-center">
                    <q-icon name="phonelink_lock" color="primary" size="28px" class="q-mr-sm" />
                    <div>
                      <div class="text-weight-bold text-slate-800">{{ t('settingsPage.twoFactorStatus') || 'Stato 2FA' }}</div>
                      <div class="text-caption text-grey-6">{{ mfaEnabled ? (t('settingsPage.active') || 'Attivo') : (t('settingsPage.inactive') || 'Non Attivo') }}</div>
                    </div>
                  </div>
                  <q-badge :color="mfaEnabled ? 'positive' : 'grey-6'" :label="mfaEnabled ? (t('settingsPage.enabled') || 'Abilitato') : (t('settingsPage.disabled') || 'Disabilitato')" />
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- Tab 3: Notifiche & Avvisi -->
        <q-tab-panel name="notifications" class="q-pa-md">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">{{ t('settingsPage.teacherNotifPref') || 'Preferenze Notifiche Docente' }}</div>
          <div class="text-caption text-slate-500 q-mb-md">{{ t('settingsPage.teacherNotifSub') || 'Seleziona come e quando desideri ricevere comunicazioni ed avvisi dal sistema scolastico' }}</div>

          <q-card flat bordered class="q-pa-md rounded-xl bg-white">
            <q-list separator>
              <q-item tag="label" v-ripple>
                <q-item-section avatar>
                  <q-icon name="email" color="indigo" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ t('settingsPage.emailCircularsTitle') || 'Notifiche Email Circolari & Comunicati' }}</q-item-label>
                  <q-item-label caption>{{ t('settingsPage.emailCircularsDesc') || 'Ricevi una copia via email di ogni nuova circolare pubblicata dal Dirigente o dalla Segreteria.' }}</q-item-label>
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
                  <q-item-label class="text-weight-bold">{{ t('settingsPage.substitutionsAlertTitle') || 'Avvisi Sostituzioni Lezioni' }}</q-item-label>
                  <q-item-label caption>{{ t('settingsPage.substitutionsAlertDesc') || 'Notifica immediata quando ti viene assegnata una sostituzione o supplenza.' }}</q-item-label>
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
                  <q-item-label class="text-weight-bold">{{ t('settingsPage.colloquiBookingsTitle') || 'Prenotazioni Colloqui Famiglie' }}</q-item-label>
                  <q-item-label caption>{{ t('settingsPage.colloquiBookingsDesc') || 'Invia notifica quando un genitore prenota o annulla un appuntamento per il ricevimento.' }}</q-item-label>
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
                  <q-item-label class="text-weight-bold">{{ t('settingsPage.scrutinyDeadlinesTitle') || 'Promemoria Scrutini & Consegna Voti' }}</q-item-label>
                  <q-item-label caption>{{ t('settingsPage.scrutinyDeadlinesDesc') || 'Avvisi automatici in prossimità della chiusura del tabellone voti per il consiglio di classe.' }}</q-item-label>
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
                  <q-item-label class="text-weight-bold">{{ t('settingsPage.inAppSoundsTitle') || 'Suoni di Notifica In-App' }}</q-item-label>
                  <q-item-label caption>{{ t('settingsPage.inAppSoundsDesc') || 'Riproduci un breve segnale acustico all\'arrivo di nuove notifiche nel registro.' }}</q-item-label>
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
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">{{ t('settingsPage.registerCustomizationTitle') || 'Personalizzazione Registro & Griglia Voti' }}</div>
          <div class="text-caption text-slate-500 q-mb-md">{{ t('settingsPage.registerCustomizationSub') || "Adatta l'aspetto grafico e il comportamento predefinito del registro elettronico" }}</div>

          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">{{ t('settingsPage.landingPageTitle') || 'Pagina Iniziale Predefinita (Landing Page)' }}</div>
                <q-select
                  v-model="registerSettings.defaultLanding"
                  :options="landingPageOptions"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">{{ t('settingsPage.gradeDisplayFormat') || 'Formato Visualizzazione Voti' }}</div>
                <q-select
                  v-model="registerSettings.gradeFormat"
                  :options="gradeFormatOptions"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">{{ t('settingsPage.decimalSeparator') || 'Separatore dei Decimali' }}</div>
                <q-select
                  v-model="registerSettings.decimalSeparator"
                  :options="decimalSeparatorOptions"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg q-mb-md"
                />

                <div class="p-3 bg-slate-50 rounded-lg border border-slate-200 text-xs">
                  <div class="text-weight-bold text-slate-700 q-mb-xs">{{ t('settingsPage.livePreviewTitle') || 'Anteprima Visualizzazione Voti Live:' }}</div>
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
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">{{ t('settingsPage.layoutDisplayOptions') || 'Opzioni Visualizzazione Layout' }}</div>
                <q-list separator>
                  <q-item tag="label" v-ripple>
                    <q-item-section>
                      <q-item-label class="text-weight-bold">{{ t('settingsPage.compactGridTitle') || 'Griglia Registro Compatta' }}</q-item-label>
                      <q-item-label caption>{{ t('settingsPage.compactGridDesc') || 'Riduci il padding delle tabelle per mostrare più studenti e colonne contemporaneamente.' }}</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                      <q-toggle v-model="registerSettings.compactGrid" color="primary" />
                    </q-item-section>
                  </q-item>

                  <q-item tag="label" v-ripple>
                    <q-item-section>
                      <q-item-label class="text-weight-bold">{{ t('settingsPage.showPhotosTitle') || "Mostra Foto Studenti nell'Appello" }}</q-item-label>
                      <q-item-label caption>{{ t('settingsPage.showPhotosDesc') || "Visualizza l'avatar o la fototessera dello studente nella schermata di appello." }}</q-item-label>
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
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">{{ t('settingsPage.digitalSignatureTitle') || 'Firma Digitale & PIN Veloce Lezioni' }}</div>
          <div class="text-caption text-slate-500 q-mb-md">{{ t('settingsPage.digitalSignatureSub') || 'Gestisci il codice PIN per la firma immediata del registro di classe' }}</div>

          <div class="row q-col-gutter-lg">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md rounded-xl bg-white">
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">{{ t('settingsPage.quickPinTitle') || 'PIN Veloce Registro (4 Cifre)' }}</div>
                <div class="text-caption text-slate-500 q-mb-md">
                  {{ t('settingsPage.quickPinDesc') || "Il PIN veloce ti consente di firmare e validare l'ora di lezione sul registro di classe senza dover reinserire la password completa." }}
                </div>

                <div class="row q-col-gutter-sm items-center q-mb-md">
                  <div class="col-8">
                    <q-input
                      v-model="pinValue"
                      type="password"
                      mask="####"
                      maxlength="4"
                      :label="t('settingsPage.quickPinLabel') || 'PIN Veloce (4 cifre)'"
                      outlined
                      dense
                      class="rounded-lg"
                    />
                  </div>
                  <div class="col-4">
                    <q-btn color="primary" :label="t('settingsPage.savePinBtn') || 'Salva PIN'" unelevated class="full-width" @click="savePin" />
                  </div>
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- Tab Accessibilità & Aspetto -->
        <q-tab-panel name="accessibility" class="q-pa-md">
          <AccessibilitySettingsPanel />
        </q-tab-panel>

      </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useThemeStore } from '@/stores/theme'
import { userService } from '@/services/userService'
import { useAuthStore } from '@/stores/auth'
import { SUPPORTED_LOCALES, applyLocale, normalizeLocale } from '@/utils/locale'
import AccessibilitySettingsPanel from '@/components/Common/AccessibilitySettingsPanel.vue'

const $q = useQuasar()
const { t, te, locale } = useI18n()
const themeStore = useThemeStore()
const authStore = useAuthStore()

const activeTab = ref('general')
const savingAll = ref(false)

// Lingua e Localizzazione
const selectedLocale = ref(normalizeLocale(locale.value))
const localeOptions = SUPPORTED_LOCALES

watch(locale, (newLoc) => {
  selectedLocale.value = normalizeLocale(newLoc)
})

const firstDayOfWeekOptions = computed(() => [
  { label: t('settingsPage.monday') || 'Lunedì', value: 1 },
  { label: t('settingsPage.sunday') || 'Domenica', value: 0 }
])

const timeFormatOptions = computed(() => [
  { label: t('settingsPage.timeFormat24') || '24 Ore (Es. 14:30)', value: '24h' },
  { label: t('settingsPage.timeFormat12') || '12 Ore AM/PM (Es. 2:30 PM)', value: '12h' }
])

const landingPageOptions = computed(() => [
  { label: t('settingsPage.landingDashboard') || 'Dashboard Generale', value: '/teacher' },
  { label: t('settingsPage.landingClasses') || 'Le Mie Classi', value: '/teacher/classes' },
  { label: t('settingsPage.landingAttendance') || 'Registro Appello & Presenze', value: '/teacher/attendance' },
  { label: t('settingsPage.landingGrades') || 'Gestione Voti', value: '/teacher/grades' },
  { label: t('settingsPage.landingAgenda') || 'Agenda & Compiti', value: '/teacher/agenda' }
])

const gradeFormatOptions = computed(() => [
  { label: t('settingsPage.gradeFormatDecimal') || 'Decimali Standard (Es. 7,5 o 7.5)', value: 'decimal' },
  { label: t('settingsPage.gradeFormatFractional') || 'Simboli e Frazioni (Es. 7+, 7½, 8-)', value: 'fractional' },
  { label: t('settingsPage.gradeFormatCentesimal') || 'Centesimi (Es. 75/100)', value: 'centesimal' }
])

const decimalSeparatorOptions = computed(() => [
  { label: t('settingsPage.commaSeparator') || 'Virgola ( , ) — Standard Italiano (es. 7,5)', value: ',' },
  { label: t('settingsPage.dotSeparator') || 'Punto ( . ) — Standard Internazionale (es. 7.5)', value: '.' }
])

const generalSettings = reactive({
  dateFormat: 'DD/MM/YYYY',
  firstDayOfWeek: 1,
  timeFormat: '24h',
  timezone: 'Europe/Rome'
})

const onLocaleChange = (newLoc) => {
  if (newLoc) {
    applyLocale(newLoc, { locale }, $q)
    $q.notify({
      type: 'positive',
      icon: 'language',
      message: t('notifications.languageChanged'),
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
  if (v <= 0.25) return t('settingsPage.strengthWeak') || 'Debole'
  if (v <= 0.5) return t('settingsPage.strengthMedium') || 'Media'
  if (v <= 0.75) return t('settingsPage.strengthGood') || 'Buona'
  return t('settingsPage.strengthStrong') || 'Forte'
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
    $q.notify({ type: 'negative', message: t('errors.passwordMismatch') })
    return
  }
  const userId = authStore.user?.id
  if (!userId) {
    $q.notify({ type: 'negative', message: t('errors.sessionInvalid') })
    return
  }
  updatingPassword.value = true
  try {
    await userService.changePassword(userId, pwdForm.currentPassword, pwdForm.newPassword)
    $q.notify({ type: 'positive', message: t('notifications.passwordUpdated') })
    pwdForm.currentPassword = ''
    pwdForm.newPassword = ''
    pwdForm.confirmPassword = ''
  } catch (err) {
    const errData = err.response?.data
    const code = errData?.code || errData?.error
    const msg = code && te(`errors.${code}`) ? t(`errors.${code}`) : (errData?.message || errData?.error || t('errors.serverError'))
    $q.notify({ type: 'negative', message: msg })
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
