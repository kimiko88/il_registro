<template>
  <q-page class="flex flex-center">
    <div class="glass-card login-card relative-position" style="width: 100%; max-width: 420px">
      <!-- Language Selector in Top Right -->
      <div class="absolute-top-right q-pa-md">
        <q-btn-dropdown
          flat
          dense
          no-caps
          color="primary"
          icon="language"
          :label="currentLangCode"
          class="rounded-lg text-weight-bold"
        >
          <q-list style="min-width: 150px">
            <q-item
              v-for="lang in languageOptions"
              :key="lang.value"
              clickable
              v-close-popup
              @click="changeLanguage(lang.value)"
              :active="locale === lang.value"
              active-class="bg-indigo-50 text-primary"
            >
              <q-item-section avatar style="min-width: 32px">
                <q-icon :name="lang.icon" size="18px" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ lang.label }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </div>

      <div class="text-center q-mb-lg q-mt-sm">
        <h1 class="text-h4 text-weight-bold text-primary q-mb-xs" style="letter-spacing: -1px">
          {{ t('login.welcomeBack') }}
        </h1>
        <div class="text-grey-7">{{ t('login.subtitle') }}</div>
      </div>

      <q-form aria-label="Modulo di accesso" @submit="onSubmit" class="q-gutter-y-md">
        <q-input
          v-model="email"
          :label="t('login.emailLabel')"
          type="email"
          autocomplete="email"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          :rules="[
            val => !!val || t('login.emailRequired'),
            val => /.+@.+\..+/.test(val) || t('login.emailInvalid')
          ]"
          @keyup.enter="() => passwordInputRef?.focus()"
        >
          <template v-slot:prepend>
            <q-icon name="email" color="primary" />
          </template>
        </q-input>

        <q-input
          ref="passwordInputRef"
          v-model="password"
          :label="t('login.passwordLabel')"
          :type="showPassword ? 'text' : 'password'"
          autocomplete="current-password"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          :rules="[val => !!val || t('login.passwordRequired')]"
          @keyup.enter="onSubmit"
        >
          <template v-slot:prepend>
            <q-icon name="lock" color="primary" />
          </template>
          <template v-slot:append>
            <q-btn
              flat
              round
              dense
              :icon="showPassword ? 'visibility_off' : 'visibility'"
              :aria-label="showPassword ? t('login.hidePassword') : t('login.showPassword')"
              @click="showPassword = !showPassword"
            />
          </template>
        </q-input>
        
        <div class="row justify-between items-center q-mt-sm">
          <q-checkbox id="remember-me" v-model="rememberMe" :label="t('login.rememberMe')" dense size="sm" color="primary" />
        </div>

        <!-- Inline Error Alert (Reattivo al cambio lingua) -->
        <div v-if="errorMessage" role="alert" aria-live="assertive" class="q-mt-sm bg-red-1 text-negative q-pa-sm rounded-lg text-caption text-center row items-center justify-center">
          <q-icon name="error_outline" size="18px" class="q-mr-xs" />
          <span>{{ errorMessage }}</span>
        </div>

        <div class="q-mt-lg">
          <q-btn 
            :label="t('login.submit')" 
            type="submit" 
            color="primary" 
            class="full-width q-py-sm shadow-soft" 
            :loading="loading" 
            no-caps
            unelevated
          />
        </div>
      </q-form>
      
      <div class="text-center q-mt-lg text-caption text-grey-6">
        {{ t('login.noAccount') }} <span class="text-primary text-weight-bold cursor-pointer hover-underline" @click="openContactSecretary">{{ t('login.contactSecretary') }}</span>
      </div>
    </div>

    <!-- Modal Contatta la Segreteria -->
    <q-dialog v-model="showSecretaryDialog">
      <q-card style="min-width: 320px; max-width: 550px; width: 100%;" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold row items-center">
            <q-icon name="contact_support" class="q-mr-sm" size="24px" /> {{ t('login.contactTitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg">
          <div class="text-subtitle2 text-grey-8 q-mb-md">
            {{ t('login.contactSubtitle') }}
          </div>

          <q-select
            v-model="selectedSchool"
            :options="schools"
            option-label="name"
            :label="t('login.selectSchool') + ' *'"
            outlined
            dense
            clearable
            :loading="loadingSchools"
            class="q-mb-lg"
          >
            <template v-slot:no-option>
              <q-item>
                <q-item-section class="text-grey">{{ t('login.noSchoolFound') }}</q-item-section>
              </q-item>
            </template>
            <template v-slot:option="scope">
              <q-item v-bind="scope.itemProps">
                <q-item-section avatar>
                  <q-icon name="school" color="primary" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ scope.opt.name }}</q-item-label>
                  <q-item-label caption v-if="scope.opt.code">{{ t('login.codeLabel') || 'Codice' }}: {{ scope.opt.code }} • {{ scope.opt.city || t('login.defaultCountry') || 'Italia' }}</q-item-label>
                </q-item-section>
              </q-item>
            </template>
          </q-select>

          <!-- Dettagli Segreteria Scuola Selezionata -->
          <div v-if="selectedSchool" class="bg-indigo-50 border-indigo q-pa-md rounded-xl">
            <div class="row items-center q-mb-sm">
              <q-icon name="apartment" color="primary" size="20px" class="q-mr-xs" />
              <div class="text-weight-bold text-slate-800 text-subtitle1">{{ selectedSchool.name }}</div>
            </div>

            <q-separator class="q-my-sm" />

            <div class="q-gutter-y-xs text-body2">
              <div class="row items-center">
                <q-icon name="email" color="indigo-8" size="18px" class="q-mr-sm" />
                <span class="text-weight-medium q-mr-xs">{{ t('login.emailSegreteria') }}:</span>
                <a :href="'mailto:' + selectedSchool.email" class="text-primary text-weight-bold text-decoration-none">
                  {{ selectedSchool.email }}
                </a>
              </div>

              <div v-if="selectedSchool.phone" class="row items-center q-mt-xs">
                <q-icon name="phone" color="indigo-8" size="18px" class="q-mr-sm" />
                <span class="text-weight-medium q-mr-xs">{{ t('login.phone') }}:</span>
                <a :href="'tel:' + selectedSchool.phone" class="text-slate-700 text-decoration-none">
                  {{ selectedSchool.phone }}
                </a>
              </div>

              <div v-if="selectedSchool.address" class="row items-center q-mt-xs text-grey-7">
                <q-icon name="place" color="indigo-8" size="18px" class="q-mr-sm" />
                <span>{{ selectedSchool.address }} {{ selectedSchool.city ? ' - ' + selectedSchool.city : '' }}</span>
              </div>
            </div>

            <div class="row q-gutter-sm q-mt-md">
              <q-btn
                color="primary"
                icon="send"
                :label="t('login.sendEmail')"
                no-caps
                unelevated
                type="a"
                :href="'mailto:' + selectedSchool.email + '?subject=Richiesta%20Creazione%20Account%20Registro%20Elettronico'"
                target="_blank"
                class="col"
              />
              <q-btn
                outline
                color="primary"
                icon="content_copy"
                :label="t('login.copyEmail')"
                no-caps
                @click="copyEmail(selectedSchool.email)"
                class="col"
              />
            </div>
          </div>

          <div v-else-if="!loadingSchools" class="text-center text-grey-6 q-py-md">
            <q-icon name="arrow_upward" size="24px" class="q-mb-xs" /><br />
            {{ t('login.chooseSchoolPrompt') }}
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-1">
          <q-btn flat :label="t('common.close')" color="grey-7" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useAuth } from '@/composables/useAuth'
import api from '@/services/api'
import { SUPPORTED_LOCALES, applyLocale, normalizeLocale } from '@/utils/locale'

const $q = useQuasar()
const { locale, t } = useI18n()
const email = ref('')
const password = ref('')
const passwordInputRef = ref(null)
const showPassword = ref(false)
const rememberMe = ref(false)
const loading = ref(false)

// Reattivo al cambio lingua senza ricaricare la pagina
const errorMessage = ref('')
const errorKey = ref('')
const errorParams = ref({})

watch(() => (locale && typeof locale === 'object' ? locale.value : locale), () => {
  if (errorKey.value) {
    errorMessage.value = t(errorKey.value, errorParams.value)
  }
})

const failedAttempts = ref(0)
const lockoutUntil = ref(null)
const { login, getLoginErrorKey } = useAuth()
const route = useRoute()

const showSecretaryDialog = ref(false)
const loadingSchools = ref(false)
const selectedSchool = ref(null)
const schools = ref([])

const languageOptions = SUPPORTED_LOCALES

const currentLangCode = computed(() => {
  const norm = normalizeLocale(locale.value)
  const opt = languageOptions.find(o => o.value === norm)
  return opt ? opt.code : 'IT'
})

function changeLanguage(langKey) {
  applyLocale(langKey, { locale }, $q)
  document.title = t('login.documentTitle') || `${t('login.welcomeBack')} — Registro Elettronico`
  if (errorKey.value) {
    errorMessage.value = t(errorKey.value, errorParams.value)
  }
}

const defaultSchools = [
  { id: '1', name: 'Liceo Scientifico Statale Galileo Galilei', code: 'RMPS010001', email: 'rmps010001@istruzione.it', phone: '+39 06 12345678', city: 'Roma', address: 'Via delle Fornaci 200' },
  { id: '2', name: 'Liceo Ginnasio Statale Ennio Quirino Visconti', code: 'RMPC080007', email: 'rmpc080007@istruzione.it', phone: '+39 06 6793508', city: 'Roma', address: 'Piazza del Collegio Romano 4' },
  { id: '3', name: 'Istituto d\'Istruzione Superiore Camillo Cavour', code: 'RMIS00100X', email: 'rmis00100x@istruzione.it', phone: '+39 06 4880574', city: 'Roma', address: 'Via delle Carine 1' },
  { id: '4', name: 'Liceo Scientifico e Linguistico Guglielmo Marconi', code: 'BOPS01000V', email: 'bops01000v@istruzione.it', phone: '+39 051 6142145', city: 'Bologna', address: 'Via Maria Grazia Agnesi 1' },
  { id: '5', name: 'Scuola di Prova (Ambiente Demo)', code: 'PROVA123', email: 'segreteria.prova@scuola.it', phone: '+39 06 5551234', city: 'Roma', address: 'Via delle Prove 10' }
]

onMounted(() => {
  document.title = t('login.documentTitle') || `${t('login.welcomeBack')} — Registro Elettronico`
  if (route?.query?.reason === 'session_expired') {
    errorKey.value = 'login.sessionExpired'
    errorMessage.value = t('login.sessionExpired')
  }
  if (route?.path === '/register' || route?.path === '/forgot-password') {
    openContactSecretary()
  }
})

async function fetchSchools() {
  loadingSchools.value = true
  try {
    const res = await api.get('/public/schools')
    if (res.data && res.data.items && res.data.items.length > 0) {
      schools.value = res.data.items
    } else {
      schools.value = defaultSchools
    }
  } catch (e) {
    console.warn('Backend public schools not reached, using defaults', e)
    schools.value = defaultSchools
  } finally {
    loadingSchools.value = false
  }
}

function openContactSecretary() {
  showSecretaryDialog.value = true
  if (schools.value.length === 0) {
    fetchSchools()
  }
}

async function copyEmail(emailStr) {
  if (!emailStr) return
  try {
    await navigator.clipboard.writeText(emailStr)
    $q.notify({
      type: 'positive',
      icon: 'content_copy',
      message: t('login.emailCopied')
    })
  } catch {
    // Clipboard API requires HTTPS; show a fallback notification
    $q.notify({
      type: 'warning',
      icon: 'content_copy',
      message: t('login.copyFailed') || `${t('login.copyManual') || 'Copia manuale'}: ${emailStr}`
    })
  }
}

async function onSubmit() {
  // UI-level lockout after repeated failures (backend is the primary rate limiter)
  if (lockoutUntil.value && Date.now() < lockoutUntil.value) {
    const secs = Math.ceil((lockoutUntil.value - Date.now()) / 1000)
    errorKey.value = 'login.tooManyAttempts'
    errorParams.value = { secs }
    errorMessage.value = t('login.tooManyAttempts', { secs }) || `Troppi tentativi. Riprova tra ${secs}s.`
    return
  }
  errorMessage.value = ''
  errorKey.value = ''
  errorParams.value = {}
  loading.value = true
  const error = await login(email.value, password.value, rememberMe.value)
  loading.value = false
  if (error) {
    failedAttempts.value++
    if (failedAttempts.value >= 5) {
      lockoutUntil.value = Date.now() + 30_000 // 30-second UI lockout
    }
    errorMessage.value = error
    if (typeof getLoginErrorKey === 'function') {
      errorKey.value = getLoginErrorKey(error)
    }
  } else {
    failedAttempts.value = 0
    lockoutUntil.value = null
  }
}
</script>

<style scoped>
.login-card {
  padding: 3rem 2.5rem;
}
@media (max-width: 599px) {
  .login-card {
    padding: 1.5rem 1.25rem !important;
    margin: 0.5rem;
  }
}
.rounded-input :deep(.q-field__control) {
  border-radius: 12px;
}
.hover-underline:hover {
  text-decoration: underline;
}
.border-indigo {
  border: 1px solid #c7d2fe;
}
</style>

