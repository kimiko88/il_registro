<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h3 text-weight-bold text-outfit q-my-none text-gradient-premium">
          {{ t('settings.title') }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">{{ t('settings.subtitle') }}</div>
      </div>
    </div>
    
    <q-card class="glass-card shadow-soft overflow-hidden">
      <q-list separator>
        <q-item-label header class="text-uppercase letter-spacing-1 text-slate-400 text-weight-bold">Generale & Configurazione</q-item-label>
        
        <!-- Lingua -->
        <q-item clickable v-ripple class="q-py-md" @click="openLanguageDialog">
          <q-item-section avatar>
            <div class="bg-indigo-50 text-indigo-600 q-pa-sm rounded-lg">
              <q-icon name="language" size="24px" />
            </div>
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-bold text-slate-700">{{ t('settings.languageLabel') }}</q-item-label>
            <q-item-label caption>{{ currentLanguageLabel }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-badge color="primary" label="Attivo" class="q-mr-sm" />
            <q-icon name="chevron_right" color="slate-300" />
          </q-item-section>
        </q-item>
  
        <!-- Sicurezza -->
        <q-item clickable v-ripple class="q-py-md" @click="openSecurityDialog">
          <q-item-section avatar>
            <div class="bg-red-50 text-red-600 q-pa-sm rounded-lg">
              <q-icon name="security" size="24px" />
            </div>
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-bold text-slate-700">{{ t('settings.securityLabel') }}</q-item-label>
            <q-item-label caption>{{ t('settings.securitySub') }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-badge :color="securitySettings.require_mfa ? 'positive' : 'warning'" :label="securitySettings.require_mfa ? t('settings.mfaActive') : t('settings.mfaOptional')" class="q-mr-sm" />
            <q-icon name="chevron_right" color="slate-300" />
          </q-item-section>
        </q-item>

        <!-- Notifiche -->
        <q-item clickable v-ripple class="q-py-md" @click="openNotificationsDialog">
          <q-item-section avatar>
            <div class="bg-amber-50 text-amber-600 q-pa-sm rounded-lg">
              <q-icon name="notifications" size="24px" />
            </div>
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-bold text-slate-700">{{ t('settings.notificationsLabel') }}</q-item-label>
            <q-item-label caption>{{ t('settings.notificationsSub') }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-badge color="info" :label="activeNotificationsCount + ' attive'" class="q-mr-sm" />
            <q-icon name="chevron_right" color="slate-300" />
          </q-item-section>
        </q-item>
      </q-list>
    </q-card>

    <!-- Dialog Sicurezza -->
    <q-dialog v-model="securityDialog" persistent>
      <q-card style="min-width: 450px; max-width: 600px;" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center">
          <div class="text-h6"><q-icon name="security" class="q-mr-sm" />Configurazione Sicurezza</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg">
          <div class="text-subtitle2 text-grey-8 q-mb-md">Gestisci le politiche di accesso e protezione account per l'intero sistema.</div>

          <q-list separator>
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('settings.mfaTitle') }}</q-item-label>
                <q-item-label caption>{{ t('settings.mfaCaption') }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="securitySettings.require_mfa" color="primary" />
              </q-item-section>
            </q-item>

            <q-item class="q-py-md">
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('settings.minPasswordLength') }}</q-item-label>
                <q-item-label caption>Numero minimo di caratteri richiesti per la creazione o modifica delle password.</q-item-label>
              </q-item-section>
              <q-item-section side style="width: 140px;">
                <q-select v-model="securitySettings.min_password_length" :options="passwordLengthOptions" dense outlined emit-value map-options />
              </q-item-section>
            </q-item>

            <q-item class="q-py-md">
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('settings.sessionTimeout') }}</q-item-label>
                <q-item-label caption>Tempo trascorso il quale la sessione utente inattiva viene disconnessa automaticamente.</q-item-label>
              </q-item-section>
              <q-item-section side style="width: 140px;">
                <q-select v-model="securitySettings.session_timeout_hours" :options="sessionTimeoutOptions" dense outlined emit-value map-options />
              </q-item-section>
            </q-item>

            <q-item class="q-py-md">
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('settings.maxLoginAttempts') }}</q-item-label>
                <q-item-label caption>Soglia di errori di login oltre la quale l'account viene temporaneamente sospeso.</q-item-label>
              </q-item-section>
              <q-item-section side style="width: 140px;">
                <q-select v-model="securitySettings.max_login_attempts" :options="loginAttemptsOptions" dense outlined emit-value map-options />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-1">
          <q-btn flat :label="t('common.cancel')" color="grey-7" v-close-popup />
          <q-btn color="primary" :label="t('common.save')" icon="save" :loading="savingSecurity" @click="saveSecuritySettings" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Notifiche -->
    <q-dialog v-model="notificationsDialog" persistent>
      <q-card style="min-width: 450px; max-width: 600px;" class="rounded-xl">
        <q-card-section class="bg-amber-8 text-white row items-center">
          <div class="text-h6"><q-icon name="notifications" class="q-mr-sm" />Preferenze Notifiche & Alert</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg">
          <div class="text-subtitle2 text-grey-8 q-mb-md">Personalizza i canali di messaggistica e gli avvisi automatici di sistema.</div>

          <q-list separator>
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Notifiche Assegnazione Supplenze</q-item-label>
                <q-item-label caption>Invia un avviso istantaneo ai docenti quando vengono nominati per una supplenza.</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notificationSettings.enable_substitute_notifications" color="amber-9" />
              </q-item-section>
            </q-item>

            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Email di Sistema per Comunicati & Circolari</q-item-label>
                <q-item-label caption>Spedisci resoconto via mail a famiglie e personale alla pubblicazione di nuovi avvisi.</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notificationSettings.enable_email_notifications" color="amber-9" />
              </q-item-section>
            </q-item>

            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Notifiche Push (App & Web)</q-item-label>
                <q-item-label caption>Abilita notifiche push per variazioni di orario, assenze dell'ultimo minuto e voti registrati.</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notificationSettings.enable_push_notifications" color="amber-9" />
              </q-item-section>
            </q-item>

            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">Alert Scadenze Scrutinii & Verbali</q-item-label>
                <q-item-label caption>Avviso automatico ai coordinatori prima della chiusura delle finestre di scrutinio.</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notificationSettings.enable_scrutiny_alerts" color="amber-9" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-1">
          <q-btn flat :label="t('common.cancel')" color="grey-7" v-close-popup />
          <q-btn color="amber-9" text-color="white" :label="t('common.save')" icon="save" :loading="savingNotifications" @click="saveNotificationSettings" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Lingua / i18n -->
    <q-dialog v-model="languageDialog">
      <q-card style="min-width: 450px; max-width: 550px;" class="rounded-xl">
        <q-card-section class="bg-indigo text-white row items-center">
          <div class="text-h6"><q-icon name="language" class="q-mr-sm" />Impostazioni Lingua (i18n)</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg">
          <div class="q-mb-md">
            <q-select
              v-model="selectedLanguage"
              :options="languageOptions"
              label="Seleziona Lingua Interfaccia"
              outlined
              emit-value
              map-options
            >
              <template v-slot:option="scope">
                <q-item v-bind="scope.itemProps">
                  <q-item-section avatar>
                    <q-icon :name="scope.opt.icon" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label>{{ scope.opt.label }}</q-item-label>
                  </q-item-section>
                </q-item>
              </template>
            </q-select>
          </div>

          <q-banner rounded class="bg-indigo-1 text-indigo-9 border-indigo">
            <template v-slot:avatar>
              <q-icon name="check_circle" color="indigo" />
            </template>
            <strong>Multilingua Attivo (i18n)</strong><br />
            Selezionando una lingua, l'intera interfaccia utente ed i messaggi di risposta del server si aggiorneranno istantaneamente.
          </q-banner>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-1">
          <q-btn flat :label="t('common.close')" color="grey-7" v-close-popup />
          <q-btn color="indigo" label="Applica Lingua" icon="check" @click="saveLanguage" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'

const $q = useQuasar()
const { locale, t } = useI18n()

// Dialog states
const securityDialog = ref(false)
const notificationsDialog = ref(false)
const languageDialog = ref(false)

const savingSecurity = ref(false)
const savingNotifications = ref(false)

// Lingua
const selectedLanguage = ref(locale.value || 'it-IT')
const languageOptions = [
  { label: 'Italiano (Italia)', value: 'it-IT', icon: 'flag' },
  { label: 'English (United States)', value: 'en-US', icon: 'language' },
  { label: 'Deutsch (Deutschland)', value: 'de-DE', icon: 'language' },
  { label: 'Français (France)', value: 'fr-FR', icon: 'language' },
  { label: 'Español (España)', value: 'es-ES', icon: 'language' },
  { label: 'Русский (Россия)', value: 'ru-RU', icon: 'language' },
  { label: 'Українська (Україна)', value: 'uk-UA', icon: 'language' },
  { label: 'العربية (السعودية)', value: 'ar-SA', icon: 'language' },
  { label: '中文 (简体)', value: 'zh-CN', icon: 'language' }
]

const currentLanguageLabel = computed(() => {
  const opt = languageOptions.find(o => o.value === locale.value)
  return opt ? opt.label : 'Italiano (Italia)'
})

// Settings state loaded from localStorage / defaults
const securitySettings = reactive({
  require_mfa: false,
  min_password_length: 8,
  session_timeout_hours: 24,
  max_login_attempts: 5
})

const notificationSettings = reactive({
  enable_substitute_notifications: true,
  enable_email_notifications: true,
  enable_push_notifications: true,
  enable_scrutiny_alerts: true
})

// Options
const passwordLengthOptions = [
  { label: '8 Caratteri (Minimo Standard)', value: 8 },
  { label: '10 Caratteri', value: 10 },
  { label: '12 Caratteri (Raccomandato)', value: 12 },
  { label: '16 Caratteri (Alta Sicurezza)', value: 16 }
]

const sessionTimeoutOptions = [
  { label: '8 Ore (Turno)', value: 8 },
  { label: '12 Ore', value: 12 },
  { label: '24 Ore (Standard)', value: 24 },
  { label: '48 Ore', value: 48 }
]

const loginAttemptsOptions = [
  { label: '3 Tentativi', value: 3 },
  { label: '5 Tentativi (Consigliato)', value: 5 },
  { label: '10 Tentativi', value: 10 }
]

const activeNotificationsCount = computed(() => {
  let count = 0
  if (notificationSettings.enable_substitute_notifications) count++
  if (notificationSettings.enable_email_notifications) count++
  if (notificationSettings.enable_push_notifications) count++
  if (notificationSettings.enable_scrutiny_alerts) count++
  return count
})

const loadSettings = () => {
  try {
    const savedSec = localStorage.getItem('superadmin_security_settings')
    if (savedSec) Object.assign(securitySettings, JSON.parse(savedSec))

    const savedNotif = localStorage.getItem('superadmin_notification_settings')
    if (savedNotif) Object.assign(notificationSettings, JSON.parse(savedNotif))

    const savedLang = localStorage.getItem('superadmin_language')
    if (savedLang) {
      selectedLanguage.value = savedLang
      locale.value = savedLang
    }
  } catch (e) {
    console.warn('Errore lettura impostazioni salvate:', e)
  }
}

const openSecurityDialog = () => {
  securityDialog.value = true
}

const openNotificationsDialog = () => {
  notificationsDialog.value = true
}

const openLanguageDialog = () => {
  languageDialog.value = true
}

const saveSecuritySettings = async () => {
  savingSecurity.value = true
  try {
    localStorage.setItem('superadmin_security_settings', JSON.stringify(securitySettings))
    $q.notify({
      type: 'positive',
      icon: 'security',
      message: 'Impostazioni di Sicurezza aggiornate con successo!'
    })
    securityDialog.value = false
  } catch (e) {
    $q.notify({
      type: 'negative',
      message: 'Errore nel salvataggio impostazioni di sicurezza'
    })
  } finally {
    savingSecurity.value = false
  }
}

const saveNotificationSettings = async () => {
  savingNotifications.value = true
  try {
    localStorage.setItem('superadmin_notification_settings', JSON.stringify(notificationSettings))
    $q.notify({
      type: 'positive',
      icon: 'notifications',
      message: 'Preferenze Notifiche salvate con successo!'
    })
    notificationsDialog.value = false
  } catch (e) {
    $q.notify({
      type: 'negative',
      message: 'Errore nel salvataggio preferenze notifiche'
    })
  } finally {
    savingNotifications.value = false
  }
}

const saveLanguage = () => {
  locale.value = selectedLanguage.value
  localStorage.setItem('superadmin_language', selectedLanguage.value)
  $q.notify({
    type: 'positive',
    icon: 'language',
    message: `Lingua impostata su: ${currentLanguageLabel.value}`
  })
  languageDialog.value = false
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.letter-spacing-1 {
    letter-spacing: 1px;
}
.border-indigo {
    border: 1px solid #c7d2fe;
}
</style>
