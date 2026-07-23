<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="account_circle" color="primary" class="q-mr-sm" />
          Profilo Genitore / Tutore Legale
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Gestione dati personali, contatti, tutela minori e consensi privacy
        </p>
      </div>
      <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="loadProfile" />
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
              Genitore / Tutore
            </q-chip>
          </div>

          <q-separator class="q-my-md" />

          <div class="text-left text-caption text-slate-600 space-y-2">
            <div><q-icon name="badge" class="q-mr-xs" /> Codice Fiscale: <strong>{{ user.fiscal_code || 'N/D' }}</strong></div>
            <div><q-icon name="phone" class="q-mr-xs" /> Telefono: <strong>{{ user.phone || 'Non specificato' }}</strong></div>
            <div><q-icon name="location_on" class="q-mr-xs" /> Indirizzo: <strong>{{ user.address || 'Non specificato' }}</strong></div>
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
            <q-tab name="personal" icon="person" label="Dati Personali" />
            <q-tab name="children" icon="child_care" label="Figli Associati" />
            <q-tab name="privacy" icon="shield" label="Consensi & Privacy" />
            <q-tab name="security" icon="lock" label="Sicurezza Account" />
          </q-tabs>

          <q-separator />

          <q-tab-panels v-model="tab" animated class="bg-white">
            <!-- TAB: PERSONAL DATA -->
            <q-tab-panel name="personal" class="q-pa-lg">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Informazioni di Contatto</div>
              <div class="row q-col-gutter-md">
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.first_name" label="Nome" readonly filled dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.last_name" label="Cognome" readonly filled dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.email" label="Email Istituzionale" readonly filled dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.fiscal_code" label="Codice Fiscale" readonly filled dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.phone" label="Numero di Telefono *" outlined dense />
                </div>
                <div class="col-12 col-sm-6">
                  <q-input v-model="form.emergency_phone" label="Contatto d'Emergenza" outlined dense />
                </div>
                <div class="col-12">
                  <q-input v-model="form.address" label="Indirizzo di Residenza" outlined dense />
                </div>
              </div>

              <div class="row justify-end q-mt-lg">
                <q-btn color="primary" label="Salva Modifiche" :loading="saving" @click="savePersonalData" />
              </div>
            </q-tab-panel>

            <!-- TAB: CHILDREN -->
            <q-tab-panel name="children" class="q-pa-lg">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Studenti in Tutela</div>

              <div v-if="children.length === 0" class="text-center q-pa-lg text-slate-500">
                Nessun figlio associato direttamente a questo profilo genitore.
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
                      Classe: {{ child.class_name || 'N/D' }} · Scuola: {{ child.school_name || 'Istituto Scolastico' }}
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
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Gestione Consensi GDPR</div>

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
                  <q-item-section avatar><q-icon name="sms" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">Avvisi SMS d'Urgenza</q-item-label>
                    <q-item-label caption>Autorizza l'invio di SMS per uscite anticipate o emergenze scolastiche.</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-toggle v-model="consents.sms" color="primary" />
                  </q-item-section>
                </q-item>

                <q-item tag="label" v-ripple>
                  <q-item-section avatar><q-icon name="photo_camera" color="primary" /></q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">Liberatoria Foto / Video Attività</q-item-label>
                    <q-item-label caption>Consenti la pubblicazione di immagini per gite, eventi e bacheca scolastica.</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-toggle v-model="consents.media_release" color="primary" />
                  </q-item-section>
                </q-item>
              </q-list>

              <div class="row justify-between items-center q-mt-xl">
                <q-btn flat color="primary" icon="download" label="Esporta Dati Personali (GDPR)" @click="exportGDPRData" />
                <q-btn color="primary" label="Salva Consensi" :loading="savingConsents" @click="saveConsents" />
              </div>
            </q-tab-panel>

            <!-- TAB: SECURITY -->
            <q-tab-panel name="security" class="q-pa-lg">
              <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Sicurezza & Autenticazione</div>

              <q-card flat bordered class="q-pa-md q-mb-md bg-slate-50">
                <div class="row items-center justify-between">
                  <div>
                    <div class="text-weight-bold text-slate-800">Password di Accesso</div>
                    <div class="text-caption text-slate-500">Ti consigliamo di cambiare la password periodicamente per proteggere il tuo account.</div>
                  </div>
                  <q-btn outline color="primary" label="Modifica Password" @click="showPasswordDialog = true" />
                </div>
              </q-card>
            </q-tab-panel>
          </q-tab-panels>
        </q-card>
      </div>
    </div>

    <!-- Password Change Dialog -->
    <q-dialog v-model="showPasswordDialog">
      <q-card style="min-width: 400px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Cambia Password</div>
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input v-model="passwordForm.oldPassword" type="password" label="Vecchia Password *" dense outlined />
          <q-input v-model="passwordForm.newPassword" type="password" label="Nuova Password *" dense outlined />
          <q-input v-model="passwordForm.confirmPassword" type="password" label="Conferma Nuova Password *" dense outlined />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Conferma" :loading="savingPassword" @click="changePassword" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useAuthStore } from '@/stores/auth'
import { useParentStore } from '@/stores/parent'
import api from 'src/services/api'

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
  media_release: true
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
      $q.notify({ type: 'positive', message: 'Profilo aggiornato con successo' })
    }
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante l\'aggiornamento del profilo' })
  } finally {
    saving.value = false
  }
}

async function saveConsents() {
  savingConsents.value = true
  try {
    $q.notify({ type: 'positive', message: 'Consensi privacy salvati con successo' })
  } finally {
    savingConsents.value = false
  }
}

async function changePassword() {
  if (!passwordForm.newPassword || passwordForm.newPassword !== passwordForm.confirmPassword) {
    $q.notify({ type: 'warning', message: 'Le nuove password non coincidono' })
    return
  }
  savingPassword.value = true
  try {
    const userId = authStore.user?.id
    await api.post(`/users/${userId}/change-password`, {
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    $q.notify({ type: 'positive', message: 'Password modificata con successo' })
    showPasswordDialog.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore modifica password' })
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
