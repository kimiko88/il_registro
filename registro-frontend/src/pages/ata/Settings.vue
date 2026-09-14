<template>
  <q-page padding class="bg-slate-50 ata-settings-page">
    <!-- Top Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge :color="roleBadgeColor" text-color="white" class="q-px-sm q-py-xs text-weight-bold rounded-borders shadow-xs">
            <q-icon name="badge" size="14px" class="q-mr-xs" />
            {{ roleDisplayName }}
          </q-badge>
          <q-badge outline color="primary" class="q-px-sm q-py-xs text-weight-bold">
            {{ t('ataDashboard.badgeAta') || 'Personale ATA' }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bold text-slate-900 q-my-none flex items-center">
          <q-icon name="manage_accounts" color="primary" class="q-mr-sm" />
          {{ t('settings.title') || 'Impostazioni & Profilo Personale' }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('settings.subtitle') || 'Gestione credenziali di accesso, sicurezza, preferenze operative e accessibilità visiva.' }}
        </p>
      </div>

      <div class="row items-center q-gutter-sm">
        <q-btn
          flat
          color="slate-700"
          icon="arrow_back"
          :label="t('common.back') || 'Torna al Pannello ATA'"
          to="/ata"
          no-caps
          class="rounded-lg border-slate-200"
        />
        <q-btn
          v-if="activeTab === 'notifications'"
          color="primary"
          icon="save"
          :label="t('common.save') || 'Salva Notifiche'"
          unelevated
          class="rounded-lg shadow-sm"
          :loading="savingNotifs"
          @click="saveNotificationSettings"
        />
        <q-btn
          v-if="activeTab === 'profile'"
          color="primary"
          icon="save"
          :label="t('common.save') || 'Salva Dati Profilo'"
          unelevated
          class="rounded-lg shadow-sm"
          :loading="savingProfile"
          @click="saveProfileData"
        />
      </div>
    </div>

    <!-- Main Navigation Card with Tabs -->
    <q-card class="rounded-2xl shadow-soft border border-slate-200 overflow-hidden bg-white q-mb-xl">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-slate-600 bg-white border-b border-slate-100"
        active-color="primary"
        indicator-color="primary"
        align="left"
        narrow-indicator
        no-caps
      >
        <q-tab name="profile" icon="account_circle" :label="t('nav.profile') || 'Profilo & Anagrafica'" class="q-px-lg py-3 text-weight-bold" />
        <q-tab name="security" icon="lock" :label="t('settingsPage.tabSecurity') || 'Sicurezza & Password'" class="q-px-lg py-3 text-weight-bold" />
        <q-tab name="notifications" icon="notifications" :label="t('settingsPage.tabNotifications') || 'Notifiche & Avvisi ATA'" class="q-px-lg py-3 text-weight-bold" />
        <q-tab name="accessibility" icon="accessibility_new" :label="t('settingsPage.tabAccessibility') || 'Accessibilità & Lingua'" class="q-px-lg py-3 text-weight-bold" />
      </q-tabs>

      <q-tab-panels v-model="activeTab" animated class="bg-transparent">

        <!-- TAB 1: PROFILO & ANAGRAFICA -->
        <q-tab-panel name="profile" class="q-pa-lg">
          <div class="row q-col-gutter-lg">
            <!-- Left Info Card -->
            <div class="col-12 col-md-4">
              <q-card flat bordered class="rounded-xl p-4 bg-slate-50 border-slate-200 full-height flex column justify-between">
                <div class="text-center q-pt-md">
                  <q-avatar size="96px" :color="roleBadgeColor" text-color="white" class="shadow-md q-mb-md">
                    <q-icon name="person" size="56px" />
                  </q-avatar>
                  <div class="text-h6 text-weight-bold text-slate-800">{{ userName }}</div>
                  <div class="text-caption text-weight-medium text-slate-500 q-mb-sm">{{ userEmail }}</div>
                  <q-badge :color="roleBadgeColor" text-color="white" class="q-px-md q-py-xs text-caption text-weight-bold rounded-borders">
                    {{ roleDisplayName }}
                  </q-badge>
                </div>

                <q-separator class="q-my-md" />

                <div class="q-gutter-y-sm text-caption text-slate-600">
                  <div class="row items-center justify-between">
                    <span class="text-slate-400">ID Utente:</span>
                    <span class="text-weight-bold font-mono">{{ userId }}</span>
                  </div>
                  <div class="row items-center justify-between" v-if="userBadgeCode">
                    <span class="text-slate-400">Codice Badge:</span>
                    <span class="text-weight-bold font-mono text-primary">{{ userBadgeCode }}</span>
                  </div>
                  <div class="row items-center justify-between">
                    <span class="text-slate-400">Stato Account:</span>
                    <q-badge color="positive" label="Attivo & Confermato" />
                  </div>
                  <div class="row items-center justify-between">
                    <span class="text-slate-400">Inquadramento:</span>
                    <span class="text-weight-bold text-slate-700">{{ isDsga ? 'Area Funzionari ed EQ' : 'Personale ATA' }}</span>
                  </div>
                  <div class="row items-center justify-between">
                    <span class="text-slate-400">Sede Servizio:</span>
                    <span class="text-weight-bold text-slate-700">Sede Centrale</span>
                  </div>
                </div>

                <div class="q-mt-lg p-3 bg-white rounded-lg border border-slate-200">
                  <div class="row items-center text-caption text-slate-600">
                    <q-icon name="verified_user" color="positive" size="20px" class="q-mr-xs" />
                    <span>Accesso conforme alle norme SPID / CIE e AgID.</span>
                  </div>
                </div>
              </q-card>
            </div>

            <!-- Right Detail Form -->
            <div class="col-12 col-md-8">
              <q-card flat bordered class="rounded-xl p-5 bg-white border-slate-200">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
                  Dati di Servizio e Recapiti Istituzionali
                </div>
                <div class="text-caption text-slate-500 q-mb-lg">
                  Visualizza e aggiorna i tuoi riferimenti per l'orario e le comunicazioni interne.
                </div>

                <div class="row q-col-gutter-md">
                  <div class="col-12 col-sm-6">
                    <q-input
                      v-model="profileForm.firstName"
                      label="Nome"
                      outlined
                      dense
                      readonly
                      bg-color="grey-1"
                      class="rounded-lg"
                    >
                      <template v-slot:append>
                        <q-icon name="lock" size="18px" color="slate-400">
                          <q-tooltip>Modificabile solo dall'amministratore di sistema</q-tooltip>
                        </q-icon>
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12 col-sm-6">
                    <q-input
                      v-model="profileForm.lastName"
                      label="Cognome"
                      outlined
                      dense
                      readonly
                      bg-color="grey-1"
                      class="rounded-lg"
                    >
                      <template v-slot:append>
                        <q-icon name="lock" size="18px" color="slate-400">
                          <q-tooltip>Modificabile solo dall'amministratore di sistema</q-tooltip>
                        </q-icon>
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12 col-sm-6">
                    <q-input
                      v-model="profileForm.email"
                      label="Email Istituzionale"
                      outlined
                      dense
                      readonly
                      bg-color="grey-1"
                      class="rounded-lg"
                    >
                      <template v-slot:append>
                        <q-icon name="email" size="18px" color="slate-400" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12 col-sm-6">
                    <q-input
                      v-model="profileForm.phone"
                      label="Telefono Diretto / Interno Ufficio"
                      placeholder="Es. 06 123456 - Int. 102"
                      outlined
                      dense
                      class="rounded-lg"
                    >
                      <template v-slot:prepend>
                        <q-icon name="phone" size="18px" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12 col-sm-6">
                    <q-input
                      v-model="profileForm.office"
                      label="Ufficio / Reparto Assegnato"
                      placeholder="Es. Ufficio DSGA / Segreteria del Personale"
                      outlined
                      dense
                      class="rounded-lg"
                    >
                      <template v-slot:prepend>
                        <q-icon name="business" size="18px" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12 col-sm-6">
                    <q-input
                      v-model="profileForm.workingHours"
                      label="Orario Tipico di Servizio"
                      placeholder="Es. Lun-Ven 08:00 - 15:12 (36h)"
                      outlined
                      dense
                      class="rounded-lg"
                    >
                      <template v-slot:prepend>
                        <q-icon name="schedule" size="18px" color="primary" />
                      </template>
                    </q-input>
                  </div>

                  <div class="col-12">
                    <q-input
                      v-model="profileForm.notes"
                      label="Note Operative Personali / Ricevimento al Pubblico"
                      placeholder="Informazioni o orari per il ricevimento dell'utenza e del personale scolastico..."
                      type="textarea"
                      rows="3"
                      outlined
                      class="rounded-lg"
                    />
                  </div>
                </div>

                <div class="row justify-end q-mt-lg">
                  <q-btn
                    color="primary"
                    icon="save"
                    label="Salva Dati Profilo"
                    unelevated
                    class="rounded-lg shadow-sm q-px-lg"
                    :loading="savingProfile"
                    @click="saveProfileData"
                  />
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- TAB 2: SICUREZZA & CAMBIO PASSWORD -->
        <q-tab-panel name="security" class="q-pa-lg">
          <div class="row q-col-gutter-lg">
            <!-- Left: Password Form -->
            <div class="col-12 col-md-7">
              <q-card flat bordered class="rounded-xl p-5 bg-white border-slate-200">
                <div class="row items-center q-mb-md">
                  <q-avatar color="indigo-50" text-color="indigo-700" icon="lock_reset" size="44px" class="q-mr-md" />
                  <div>
                    <div class="text-h6 text-weight-bold text-slate-800">
                      {{ t('settingsPage.editPassword') || 'Modifica Password Personale' }}
                    </div>
                    <div class="text-caption text-slate-500">
                      {{ t('settingsPage.accountSecuritySub') || 'Aggiorna la tua password di accesso al registro elettronico' }}
                    </div>
                  </div>
                </div>

                <q-form @submit.prevent="changePassword" class="q-gutter-y-md q-mt-sm">
                  <!-- Password Attuale -->
                  <q-input
                    v-model="pwdForm.currentPassword"
                    :type="showCurrentPwd ? 'text' : 'password'"
                    :label="t('settingsPage.currentPasswordLabel') || 'Password Attuale *'"
                    outlined
                    dense
                    autocomplete="current-password"
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

                  <!-- Nuova Password -->
                  <q-input
                    v-model="pwdForm.newPassword"
                    :type="showNewPwd ? 'text' : 'password'"
                    label="Nuova Password * (min. 10 caratteri)"
                    outlined
                    dense
                    autocomplete="new-password"
                    class="rounded-lg"
                    :rules="[
                      val => !!val || (t('settingsPage.newPasswordReq') || 'Inserisci la nuova password'),
                      val => val.length >= 10 || 'La password deve contenere almeno 10 caratteri'
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

                  <!-- Password Strength Indicator -->
                  <div v-if="pwdForm.newPassword" class="q-mb-xs p-3 bg-slate-50 rounded-lg border border-slate-200">
                    <div class="row items-center justify-between text-caption text-slate-700 q-mb-xs">
                      <span class="text-weight-medium">Forza Password:</span>
                      <span class="text-weight-bold" :class="passwordStrengthColorText">{{ passwordStrengthLabel }}</span>
                    </div>
                    <q-linear-progress
                      :value="passwordStrengthValue"
                      :color="passwordStrengthColor"
                      class="rounded-borders"
                      style="height: 6px"
                    />
                    <div class="row q-gutter-x-md text-xs text-slate-500 q-mt-xs">
                      <span :class="{ 'text-positive text-weight-bold': pwdForm.newPassword.length >= 10 }">
                        ✓ Min. 10 car.
                      </span>
                      <span :class="{ 'text-positive text-weight-bold': /[A-Z]/.test(pwdForm.newPassword) }">
                        ✓ Lettera maiuscola
                      </span>
                      <span :class="{ 'text-positive text-weight-bold': /[0-9]/.test(pwdForm.newPassword) }">
                        ✓ Numero
                      </span>
                      <span :class="{ 'text-positive text-weight-bold': /[^A-Za-z0-9]/.test(pwdForm.newPassword) }">
                        ✓ Simbolo speciale
                      </span>
                    </div>
                  </div>

                  <!-- Conferma Nuova Password -->
                  <q-input
                    v-model="pwdForm.confirmPassword"
                    :type="showConfirmPwd ? 'text' : 'password'"
                    :label="t('settingsPage.confirmPasswordLabel') || 'Conferma Nuova Password *'"
                    outlined
                    dense
                    autocomplete="new-password"
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

                  <div class="row items-center justify-between q-mt-lg">
                    <q-btn
                      flat
                      color="slate-500"
                      label="Cancella"
                      no-caps
                      @click="resetPwdForm"
                    />
                    <q-btn
                      type="submit"
                      color="primary"
                      icon="lock_reset"
                      :label="t('settingsPage.updatePasswordBtn') || 'Aggiorna Password'"
                      unelevated
                      :loading="updatingPassword"
                      class="rounded-lg shadow-sm q-px-xl"
                    />
                  </div>
                </q-form>
              </q-card>
            </div>

            <!-- Right: Security Policies & Info -->
            <div class="col-12 col-md-5">
              <div class="q-gutter-y-md">
                <!-- AgID / ACN Security Card -->
                <q-card flat bordered class="rounded-xl p-4 bg-white border-slate-200">
                  <div class="row items-center q-mb-sm">
                    <q-avatar color="teal-50" text-color="teal-8" icon="shield" size="36px" class="q-mr-sm" />
                    <div class="text-subtitle2 text-weight-bold text-slate-800">Requisiti di Sicurezza AgID</div>
                  </div>
                  <div class="text-caption text-slate-600 q-mb-sm">
                    In conformità con le linee guida AgID per i sistemi della Pubblica Amministrazione Scolastica:
                  </div>
                  <ul class="text-caption text-slate-600 q-pl-md q-my-none space-y-1">
                    <li>Lunghezza minima obbligatoria di 10 caratteri</li>
                    <li>Utilizzo di maiuscole, minuscole, cifre e caratteri speciali</li>
                    <li>Non riutilizzare password già impiegate in precedenza</li>
                    <li>Non comunicare mai le credenziali a colleghi o terzi</li>
                  </ul>
                </q-card>

                <!-- 2FA Info Card -->
                <q-card flat bordered class="rounded-xl p-4 bg-white border-slate-200">
                  <div class="row items-center justify-between q-mb-sm">
                    <div class="row items-center">
                      <q-avatar color="blue-50" text-color="blue-700" icon="phonelink_lock" size="36px" class="q-mr-sm" />
                      <div>
                        <div class="text-subtitle2 text-weight-bold text-slate-800">Autenticazione a Due Fattori</div>
                        <div class="text-xs text-slate-500">SPID / CIE / OTP</div>
                      </div>
                    </div>
                    <q-badge color="positive" label="Attivo" />
                  </div>
                  <div class="text-caption text-slate-600">
                    Il tuo profilo è abilitato per l'accesso sicuro con identità digitale. Eventuali variazioni sulle modalità MFA possono essere richieste all'Amministratore di Sistema.
                  </div>
                </q-card>
              </div>
            </div>
          </div>
        </q-tab-panel>

        <!-- TAB 3: NOTIFICHE & AVVISI ATA -->
        <q-tab-panel name="notifications" class="q-pa-lg">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
            Preferenze Notifiche & Avvisi ATA
          </div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Configura le notifiche e gli avvisi di servizio specifici per le mansioni amministrative e contabili.
          </div>

          <q-card flat bordered class="rounded-xl bg-white border-slate-200 overflow-hidden">
            <q-list separator>
              <!-- 1. Richieste Ferie e Permessi da vistare (DSGA) -->
              <q-item tag="label" v-ripple class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="indigo-50" text-color="indigo-700" icon="fact_check" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-800">
                    Richieste Ferie & Permessi da Vistare
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    Ricevi un avviso immediato via email quando un assistente o collaboratore inoltra una richiesta di ferie, permesso o recupero.
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.leaveVistoAlerts" color="indigo" />
                </q-item-section>
              </q-item>

              <!-- 2. Scadenze Flussi SIDI -->
              <q-item tag="label" v-ripple class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="teal-50" text-color="teal-700" icon="cloud_sync" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-800">
                    Promemoria Scadenze Flussi SIDI & Bilancio
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    Avvisi automatici in prossimità delle scadenze per la trasmissione dei flussi contabili e del personale al Ministero (MIM).
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.sidiDeadlines" color="teal" />
                </q-item-section>
              </q-item>

              <!-- 3. Rilevazione Scioperi -->
              <q-item tag="label" v-ripple class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="amber-50" text-color="amber-800" icon="campaign" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-800">
                    Rilevazione Preventiva Scioperi
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    Notifiche per l'apertura di nuove schede di rilevazione sciopero e promemoria compilazione adesione del personale.
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.strikeAlerts" color="amber-8" />
                </q-item-section>
              </q-item>

              <!-- 4. Circolari & Comunicazioni Interne -->
              <q-item tag="label" v-ripple class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="purple-50" text-color="purple-700" icon="mail" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-800">
                    Circolari & Comunicati di Segreteria
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    Avviso immediato alla pubblicazione di nuove disposizioni dirigenziali o circolari con richiesta di firma per presa visione.
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.circularsAlerts" color="purple" />
                </q-item-section>
              </q-item>

              <!-- 5. Suoni In-App -->
              <q-item tag="label" v-ripple class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="rose-50" text-color="rose-700" icon="volume_up" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold text-slate-800">
                    Suoni di Notifica In-App
                  </q-item-label>
                  <q-item-label caption class="text-slate-500">
                    Riproduci un breve segnale acustico all'arrivo di nuove notifiche operative in tempo reale.
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="notificationSettings.inAppSound" color="rose-7" />
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>

          <div class="row justify-end q-mt-lg">
            <q-btn
              color="primary"
              icon="save"
              label="Salva Preferenze Notifiche"
              unelevated
              class="rounded-lg shadow-sm q-px-xl"
              :loading="savingNotifs"
              @click="saveNotificationSettings"
            />
          </div>
        </q-tab-panel>

        <!-- TAB 4: ACCESSIBILITÀ & LINGUA -->
        <q-tab-panel name="accessibility" class="q-pa-lg">
          <!-- Language Selector -->
          <q-card flat bordered class="rounded-xl p-4 bg-white border-slate-200 q-mb-lg">
            <div class="row items-center justify-between">
              <div class="row items-center">
                <q-avatar color="indigo-50" text-color="indigo-700" icon="language" size="40px" class="q-mr-md" />
                <div>
                  <div class="text-subtitle1 text-weight-bold text-slate-800">
                    {{ t('settingsPage.appLanguage') || "Lingua dell'Interfaccia" }}
                  </div>
                  <div class="text-caption text-slate-500">
                    Seleziona la lingua di visualizzazione per l'intero registro elettronico
                  </div>
                </div>
              </div>

              <div style="min-width: 220px;">
                <q-select
                  v-model="selectedLocale"
                  :options="localeOptions"
                  emit-value
                  map-options
                  outlined
                  dense
                  class="rounded-lg"
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
              </div>
            </div>
          </q-card>

          <!-- Reusable Visual Accessibility Panel -->
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
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { authService } from '@/services/authService'
import { userService } from '@/services/userService'
import AccessibilitySettingsPanel from '@/components/Common/AccessibilitySettingsPanel.vue'
import { SUPPORTED_LOCALES, applyLocale, normalizeLocale } from '@/utils/locale'

const $q = useQuasar()
const { t, te, locale } = useI18n()
const authStore = useAuthStore()
const { user, userRole, userName } = storeToRefs(authStore)

const activeTab = ref('profile')

// User details
const userId = computed(() => user.value?.id || '—')
const userEmail = computed(() => user.value?.email || '—')
const userBadgeCode = computed(() => user.value?.badge_code || '')
const isDsga = computed(() => (userRole.value || '').toLowerCase() === 'dsga')

const roleDisplayName = computed(() => {
  const role = (userRole.value || '').toLowerCase()
  if (te('roles.' + role)) {
    return t('roles.' + role)
  }
  const map = {
    dsga: 'DSGA (Direttore dei Servizi Generali e Amministrativi)',
    collaboratore_ds: 'Collaboratore del Dirigente Scolastico',
    assistente_amministrativo: 'Assistente Amministrativo',
    collaboratore_scolastico: 'Collaboratore Scolastico',
    assistente_alunni: 'Assistente Amministrativo - Area Alunni',
    assistente_personale: 'Assistente Amministrativo - Area Personale',
    assistente_contabilita: 'Assistente Amministrativo - Area Contabile',
    assistente_protocollo: 'Assistente Amministrativo - Protocollo',
    assistente_sportello: 'Assistente Amministrativo - Sportello Utenza',
    assistente_tecnico: 'Assistente Tecnico',
    responsabile_servizio: 'Responsabile di Servizio',
    principal: 'Dirigente Scolastico',
    vice_principal: 'Collaboratore Vicario',
    admin: 'Amministratore',
    superadmin: 'Super Amministratore'
  }
  return map[role] || role
})

const roleBadgeColor = computed(() => {
  const role = (userRole.value || '').toLowerCase()
  const map = {
    dsga: 'teal-8',
    collaboratore_ds: 'deep-orange-7',
    assistente_amministrativo: 'cyan-8',
    collaboratore_scolastico: 'amber-9',
    assistente_alunni: 'indigo-8',
    assistente_personale: 'purple-8',
    assistente_contabilita: 'green-8',
    assistente_protocollo: 'blue-8',
    assistente_sportello: 'orange-8',
    assistente_tecnico: 'blue-grey-8',
    responsabile_servizio: 'deep-purple-8'
  }
  return map[role] || 'primary'
})

// Profile Form
const savingProfile = ref(false)
const profileForm = reactive({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  office: '',
  workingHours: '',
  notes: ''
})

const initProfileForm = () => {
  if (user.value) {
    profileForm.firstName = user.value.first_name || user.value.firstName || ''
    profileForm.lastName = user.value.last_name || user.value.lastName || ''
    profileForm.email = user.value.email || ''
  }
  const savedExtra = localStorage.getItem(`ata_profile_extra_${user.value?.id || 'default'}`)
  if (savedExtra) {
    try {
      const parsed = JSON.parse(savedExtra)
      profileForm.phone = parsed.phone || ''
      profileForm.office = parsed.office || ''
      profileForm.workingHours = parsed.workingHours || ''
      profileForm.notes = parsed.notes || ''
    } catch (e) {
      console.warn('Could not parse saved ATA profile details', e)
    }
  } else {
    // Defaults based on role
    if (isDsga.value) {
      profileForm.office = 'Direzione Servizi Generali e Amministrativi'
      profileForm.workingHours = 'Lun-Ven 08:00 - 15:12 (36h)'
    }
  }
}

const saveProfileData = () => {
  savingProfile.value = true
  try {
    const extraData = {
      phone: profileForm.phone,
      office: profileForm.office,
      workingHours: profileForm.workingHours,
      notes: profileForm.notes
    }
    localStorage.setItem(`ata_profile_extra_${user.value?.id || 'default'}`, JSON.stringify(extraData))
    $q.notify({
      type: 'positive',
      message: 'Dati del profilo salvati con successo',
      position: 'top',
      timeout: 1800
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: 'Errore durante il salvataggio dei dati profilo'
    })
  } finally {
    savingProfile.value = false
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

const resetPwdForm = () => {
  pwdForm.currentPassword = ''
  pwdForm.newPassword = ''
  pwdForm.confirmPassword = ''
}

const passwordStrengthValue = computed(() => {
  const p = pwdForm.newPassword
  if (!p) return 0
  let score = 0
  if (p.length >= 10) score += 0.25
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
  return 'Molto Forte'
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
    $q.notify({
      type: 'negative',
      message: t('settingsPage.passwordMismatch') || 'Le nuove password non coincidono'
    })
    return
  }

  if (pwdForm.newPassword.length < 10) {
    $q.notify({
      type: 'warning',
      message: 'La password deve contenere almeno 10 caratteri'
    })
    return
  }

  updatingPassword.value = true
  try {
    // Attempt authService first
    try {
      await authService.changePassword(pwdForm.currentPassword, pwdForm.newPassword)
    } catch (authErr) {
      // If authService endpoint 404s or needs userId, fallback to userService
      const userIdVal = user.value?.id
      if (userIdVal && authErr.response?.status === 404) {
        await userService.changePassword(userIdVal, pwdForm.currentPassword, pwdForm.newPassword)
      } else {
        throw authErr
      }
    }

    $q.notify({
      type: 'positive',
      icon: 'check_circle',
      message: t('notifications.passwordUpdated') || 'Password aggiornata con successo!',
      position: 'top',
      timeout: 2500
    })
    resetPwdForm()
  } catch (err) {
    const errData = err.response?.data
    const code = errData?.code || errData?.error
    const msg = code && te(`errors.${code}`)
      ? t(`errors.${code}`)
      : (errData?.message || errData?.error || t('errors.serverError') || 'Errore durante la modifica della password')
    $q.notify({
      type: 'negative',
      message: msg,
      position: 'top'
    })
  } finally {
    updatingPassword.value = false
  }
}

// Notifications Settings
const savingNotifs = ref(false)
const notificationSettings = reactive({
  leaveVistoAlerts: true,
  sidiDeadlines: true,
  strikeAlerts: true,
  circularsAlerts: true,
  inAppSound: false
})

const loadNotificationSettings = () => {
  const saved = localStorage.getItem('ata_notification_settings')
  if (saved) {
    try {
      Object.assign(notificationSettings, JSON.parse(saved))
    } catch (e) {
      console.warn('Could not parse ATA notification settings', e)
    }
  }
}

const saveNotificationSettings = () => {
  savingNotifs.value = true
  try {
    localStorage.setItem('ata_notification_settings', JSON.stringify(notificationSettings))
    $q.notify({
      type: 'positive',
      message: 'Preferenze di notifica salvate con successo',
      position: 'top',
      timeout: 1800
    })
  } catch (e) {
    $q.notify({
      type: 'negative',
      message: 'Errore durante il salvataggio delle notifiche'
    })
  } finally {
    savingNotifs.value = false
  }
}

// Language Settings
const selectedLocale = ref(normalizeLocale(locale.value))
const localeOptions = SUPPORTED_LOCALES

watch(() => (typeof locale === 'object' && locale && 'value' in locale ? locale.value : locale), (newLoc) => {
  if (newLoc) selectedLocale.value = normalizeLocale(newLoc)
})

const onLocaleChange = (newLoc) => {
  if (newLoc) {
    applyLocale(newLoc, { locale }, $q)
    $q.notify({
      type: 'positive',
      icon: 'language',
      message: t('notifications.languageChanged') || 'Lingua aggiornata con successo',
      position: 'top',
      timeout: 1500
    })
  }
}

onMounted(() => {
  initProfileForm()
  loadNotificationSettings()
})
</script>

<style scoped>
.ata-settings-page {
  min-height: 100vh;
}
.shadow-soft {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}
.py-3 {
  padding-top: 0.75rem;
  padding-bottom: 0.75rem;
}
</style>
