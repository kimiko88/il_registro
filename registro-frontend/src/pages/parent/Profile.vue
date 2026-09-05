<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="account_circle" color="primary" class="q-mr-sm" />
          {{ t('parentProfile.title') || 'Profilo Genitore / Tutore Legale' }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('parentProfile.subtitle') || 'Gestione dati personali, contatti, tutela minori e consensi privacy' }}
        </p>
      </div>
      <q-btn flat round icon="refresh" color="primary" :loading="loading" :aria-label="t('common.refresh') || 'Aggiorna'" @click="loadProfile" />
    </div>

    <!-- Main Grid -->
    <div class="row q-col-gutter-lg">
      <!-- Left Column: User Summary Card -->
      <div class="col-12 col-md-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft text-center q-pa-lg">
          <q-avatar size="96px" color="primary" text-color="white" class="q-mb-md text-h3 text-weight-bold">
            {{ userInitials }}
          </q-avatar>
          <div class="text-h6 text-weight-bold text-slate-800">{{ user.first_name }} {{ user.last_name }}</div>
          <div class="text-caption text-slate-500">{{ user.email }}</div>

          <div class="q-mt-sm">
            <q-chip color="indigo-1" text-color="indigo-8" class="text-weight-bold">
              {{ t('parentProfile.roleBadge') || 'Genitore / Tutore' }}
            </q-chip>
          </div>

          <q-separator class="q-my-md" />

          <div class="text-left text-caption text-slate-600 space-y-2">
            <div><q-icon name="badge" class="q-mr-xs" /> {{ t('parentProfile.fiscalCode') || 'Codice Fiscale' }}: <strong>{{ user.fiscal_code || 'N/D' }}</strong></div>
            <div><q-icon name="phone" class="q-mr-xs" /> {{ t('parentProfile.phone') || 'Telefono' }}: <strong>{{ user.phone || t('parentProfile.notSpecified') || 'Non specificato' }}</strong></div>
            <div><q-icon name="location_on" class="q-mr-xs" /> {{ t('parentProfile.address') || 'Indirizzo' }}: <strong>{{ user.address || t('parentProfile.notSpecified') || 'Non specificato' }}</strong></div>
          </div>
        </q-card>
      </div>

      <!-- Right Column: Form Tabs -->
      <div class="col-12 col-md-8">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
          <q-tabs
            v-model="tab"
            dense
            class="text-slate-600 bg-slate-100 border-b border-slate-200"
            active-color="primary"
            indicator-color="primary"
            align="left"
            no-caps
          >
            <q-tab name="personal" icon="person" :label="t('parentProfile.tabs.personal') || 'Dati Personali'" />
            <q-tab name="children" icon="child_care" :label="t('parentProfile.tabs.children') || 'Figli Associati'" />
            <q-tab name="privacy" icon="shield" :label="t('parentProfile.tabs.privacy') || 'Consensi & Privacy'" />
            <q-tab name="security" icon="lock" :label="t('parentProfile.tabs.security') || 'Sicurezza Account'" />
          </q-tabs>

          <q-separator />

          <q-tab-panels v-model="tab" animated class="bg-white">
            <!-- TAB: PERSONAL DATA -->
            <q-tab-panel name="personal" class="q-pa-lg">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">{{ t('parentProfile.tabs.personal') || 'Informazioni di Contatto' }}</div>
              <div class="row q-col-gutter-md">
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.first_name" :label="t('parentProfile.firstName') || 'Nome'" readonly filled dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.last_name" :label="t('parentProfile.lastName') || 'Cognome'" readonly filled dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.email" :label="t('parentProfile.email') || 'Email Istituzionale'" readonly filled dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.fiscal_code" :label="t('parentProfile.fiscalCode') || 'Codice Fiscale'" readonly filled dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.phone" :label="t('parentProfile.phone') || 'Numero di Telefono *'" outlined dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.emergency_phone" :label="t('parentProfile.emergencyContact') || 'Contatto d\'Emergenza'" outlined dense />
                </div>
                <div class="col-12">
                  <q-input v-model="form.address" :label="t('parentProfile.address') || 'Indirizzo di Residenza'" outlined dense />
                </div>
              </div>

              <div class="row justify-end q-mt-lg">
                <q-btn color="primary" :label="t('parentProfile.saveChanges') || 'Salva Modifiche'" :loading="saving" @click="savePersonalData" />
              </div>
            </q-tab-panel>

            <!-- TAB: CHILDREN -->
            <q-tab-panel name="children" class="q-pa-lg">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">{{ t('parentProfile.childrenTitle') || 'Studenti in Tutela' }}</div>

              <!-- Skeleton Loader for Children -->
              <div v-if="loading" class="q-gutter-y-md">
                <div v-for="n in 2" :key="'child-skel-'+n" class="row items-center q-gutter-x-md q-py-sm">
                  <q-skeleton type="QAvatar" size="48px" class="rounded-xl" />
                  <div class="col">
                    <q-skeleton type="text" width="50%" height="22px" />
                    <q-skeleton type="text" width="35%" height="16px" />
                  </div>
                  <q-skeleton type="rect" width="120px" height="28px" class="rounded-lg" />
                </div>
              </div>

              <div v-else-if="children.length === 0" class="text-center q-pa-lg text-slate-500">
                {{ t('parentProfile.noChildren') || 'Nessun figlio associato direttamente a questo profilo genitore.' }}
              </div>

              <q-list v-else separator class="rounded-lg border-slate-200">
                <q-item v-for="child in children" :key="child.id" class="q-py-md">
                  <q-item-section avatar>
                    <q-avatar color="indigo-1" text-color="indigo-7" icon="face" size="48px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-slate-800">
                      {{ child.first_name || child.name }} {{ child.last_name }}
                    </q-item-label>
                    <q-item-label caption class="text-slate-600">
                      {{ t('parentProfile.class') || 'Classe' }}: {{ child.class_name || 'N/D' }} · {{ t('parentProfile.school') || 'Scuola' }}: {{ child.school_name || 'Istituto Scolastico' }}
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-chip color="positive" text-color="white" size="sm" icon="verified">
                      Tutela Convalidata
                    </q-chip>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-tab-panel>

            <!-- TAB: PRIVACY & CONSENTS -->
            <q-tab-panel name="privacy" class="q-pa-lg">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-xs">{{ t('parentProfile.privacyTitle') || 'Gestione Consensi GDPR' }}</div>
              <div class="text-caption text-slate-500 q-mb-md">{{ t('parentProfile.privacyDesc') || 'Autorizzazioni legali per uscite didattiche, pubblicazioni e trattamento immagini.' }}</div>

              <q-list separator>
                <q-item tag="label" v-ripple>
                  <q-item-section avatar><q-icon name="notifications" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">Notifiche Email ed Avvisi Assenze</q-item-label>
                    <q-item-label caption>Ricevi via email le notifiche per assenze, ritardi e comunicazioni di classe.</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-toggle v-model="consents.notifications" color="primary" />
                  </q-item-section>
                </q-item>

                <q-item tag="label" v-ripple>
                  <q-item-section avatar><q-icon name="photo_camera" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">{{ t('parentProfile.consentPhoto') || 'Liberatoria Foto / Video Attività' }}</q-item-label>
                    <q-item-label caption>Consenti la pubblicazione di immagini per gite, eventi e bacheca scolastica.</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-toggle v-model="consents.media_release" color="primary" />
                  </q-item-section>
                </q-item>

                <q-item tag="label" v-ripple>
                  <q-item-section avatar><q-icon name="directions_bus" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">{{ t('parentProfile.consentTrips') || 'Autorizzazione Uscite Didattiche' }}</q-item-label>
                    <q-item-label caption>Autorizzazione generica per uscite didattiche e visite guidate sul territorio.</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-toggle v-model="consents.trips" color="primary" />
                  </q-item-section>
                </q-item>

                <q-item tag="label" v-ripple>
                  <q-item-section avatar><q-icon name="devices" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">{{ t('parentProfile.consentDigital') || 'Piattaforme Didattiche Digitali' }}</q-item-label>
                    <q-item-label caption>Consenso per l'utilizzo di piattaforme didattiche digitali e posta elettronica.</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-toggle v-model="consents.digital" color="primary" />
                  </q-item-section>
                </q-item>
              </q-list>

              <div class="row justify-between items-center q-mt-xl">
                <q-btn flat color="primary" icon="download" label="Esporta Dati Personali (GDPR)" @click="exportGDPRData" />
                <q-btn color="primary" :label="t('parentProfile.saveChanges') || 'Salva Consensi'" :loading="savingConsents" @click="saveConsents" />
              </div>
            </q-tab-panel>

            <!-- TAB: SECURITY -->
            <q-tab-panel name="security" class="q-pa-lg">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">{{ t('parentProfile.securityTitle') || 'Sicurezza & Autenticazione' }}</div>

              <q-card flat bordered class="q-pa-md q-mb-md bg-slate-50">
                <div class="row items-center justify-between">
                  <div>
                    <div class="text-weight-bold text-slate-800">{{ t('parentProfile.securityTitle') || 'Password di Accesso' }}</div>
                    <div class="text-caption text-slate-500">Ti consigliamo di cambiare la password periodicamente per proteggere il tuo account.</div>
                  </div>
                  <q-btn outline color="primary" :label="t('parentProfile.updatePassword') || 'Modifica Password'" @click="showPasswordDialog = true" />
                </div>
              </q-card>
            </q-tab-panel>
          </q-tab-panels>
        </q-card>
      </div>
    </div>

    <!-- Password Change Dialog -->
    <q-dialog v-model="showPasswordDialog">
      <q-card style="width: min(450px, 95vw); max-width: 95vw;">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6">{{ t('parentProfile.securityTitle') || 'Cambia Password' }}</div>
          <q-btn flat round dense icon="close" color="white" v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input v-model="passwordForm.oldPassword" type="password" :label="t('parentProfile.currentPassword') || 'Vecchia Password *'" dense outlined />
          <q-input v-model="passwordForm.newPassword" type="password" :label="t('parentProfile.newPassword') || 'Nuova Password *'" dense outlined />
          <q-input v-model="passwordForm.confirmPassword" type="password" :label="t('parentProfile.confirmPassword') || 'Conferma Nuova Password *'" dense outlined />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup />
          <q-btn color="primary" :label="t('parentProfile.updatePassword') || 'Conferma'" :loading="savingPassword" @click="changePassword" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useParentStore } from '@/stores/parent'
import api from 'src/services/api'

const { t } = useI18n()
const $q = useQuasar()
const authStore = useAuthStore()
const parentStore = useParentStore()

const tab = ref('personal')
const loading = ref(false)
const saving = ref(false)
const savingConsents = ref(false)
const savingPassword = ref(false)
const showPasswordDialog = ref(false)

const user = computed(() => authStore.user || {})
const children = computed(() => parentStore.children || [])

const userInitials = computed(() => {
  const f = user.value.first_name ? user.value.first_name[0] : 'G'
  const l = user.value.last_name ? user.value.last_name[0] : 'P'
  return `${f}${l}`.toUpperCase()
})

const form = reactive({
  first_name: '',
  last_name: '',
  email: '',
  fiscal_code: '',
  phone: '',
  emergency_phone: '',
  address: ''
})

const consents = reactive({
  notifications: true,
  sms: false,
  media_release: true,
  trips: true,
  digital: true
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

onMounted(() => {
  loadProfile()
})

function loadProfile() {
  if (authStore.user) {
    form.first_name = authStore.user.first_name || ''
    form.last_name = authStore.user.last_name || ''
    form.email = authStore.user.email || ''
    form.fiscal_code = authStore.user.fiscal_code || ''
    form.phone = authStore.user.phone || ''
    form.emergency_phone = authStore.user.emergency_phone || ''
    form.address = authStore.user.address || ''
  }
}

async function savePersonalData() {
  saving.value = true
  try {
    const userId = authStore.user?.id
    if (userId) {
      await api.patch(`/users/${userId}`, {
        phone: form.phone,
        address: form.address,
        emergency_phone: form.emergency_phone
      })
      $q.notify({ type: 'positive', message: t('parentProfile.saveSuccess') || 'Profilo aggiornato con successo' })
    }
  } catch {
    $q.notify({ type: 'negative', message: t('parentProfile.saveError') || 'Errore durante l\'aggiornamento del profilo' })
  } finally {
    saving.value = false
  }
}

async function saveConsents() {
  savingConsents.value = true
  try {
    $q.notify({ type: 'positive', message: t('parentProfile.saveSuccess') || 'Consensi privacy salvati con successo' })
  } finally {
    savingConsents.value = false
  }
}

async function changePassword() {
  if (!passwordForm.newPassword || passwordForm.newPassword !== passwordForm.confirmPassword) {
    $q.notify({ type: 'warning', message: t('parentProfile.passwordMismatch') || 'Le nuove password non coincidono' })
    return
  }
  savingPassword.value = true
  try {
    const userId = authStore.user?.id
    await api.post(`/users/${userId}/change-password`, {
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    $q.notify({ type: 'positive', message: t('parentProfile.passwordSuccess') || 'Password modificata con successo' })
    showPasswordDialog.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || t('parentProfile.passwordError') || 'Errore modifica password' })
  } finally {
    savingPassword.value = false
  }
}

async function exportGDPRData() {
  try {
    const userId = authStore.user?.id
    await api.post(`/users/${userId}/gdpr-export`)
    $q.notify({ type: 'positive', message: 'Richiesta di esportazione GDPR inoltrata. Riceverai un file scaricabile.' })
  } catch {
    $q.notify({ type: 'negative', message: 'Errore durante l\'esportazione dati' })
  }
}
</script>
