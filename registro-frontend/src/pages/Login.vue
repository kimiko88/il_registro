<template>
  <q-page class="flex flex-center">
    <div class="glass-card q-pa-xl" style="width: 100%; max-width: 420px">
      <div class="text-center q-mb-lg">
        <h1 class="text-h4 text-weight-bold text-primary q-mb-xs" style="letter-spacing: -1px">Bentornato</h1>
        <div class="text-grey-7">Accedi per entrare nel Registro Elettronico</div>
      </div>

      <q-form aria-label="Modulo di accesso" @submit="onSubmit" class="q-gutter-y-md">
        <q-input
          v-model="email"
          label="Indirizzo Email"
          type="email"
          autocomplete="email"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          :rules="[val => !!val || 'L\'email è obbligatoria']"
          @keyup.enter="() => passwordInputRef?.focus()"
        >
          <template v-slot:prepend>
            <q-icon name="email" color="primary" />
          </template>
        </q-input>

        <q-input
          ref="passwordInputRef"
          v-model="password"
          label="Password"
          :type="showPassword ? 'text' : 'password'"
          autocomplete="current-password"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          :rules="[val => !!val || 'La password è obbligatoria']"
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
              :aria-label="showPassword ? 'Nascondi password' : 'Mostra password'"
              @click="showPassword = !showPassword"
            />
          </template>
        </q-input>
        
        <div class="row justify-between items-center q-mt-sm">
          <q-checkbox id="remember-me" v-model="rememberMe" label="Ricordami" dense size="sm" color="primary" />
        </div>

        <!-- Inline Error Alert -->
        <div v-if="errorMessage" role="alert" aria-live="assertive" class="q-mt-sm bg-red-1 text-negative q-pa-sm rounded-lg text-caption text-center row items-center justify-center">
          <q-icon name="error_outline" size="18px" class="q-mr-xs" />
          <span>{{ errorMessage }}</span>
        </div>

        <div class="q-mt-lg">
          <q-btn 
            label="Accedi" 
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
        Non hai un account? <span class="text-primary text-weight-bold cursor-pointer hover-underline" @click="openContactSecretary">Contatta la Segreteria</span>
      </div>
    </div>

    <!-- Modal Contatta la Segreteria -->
    <q-dialog v-model="showSecretaryDialog">
      <q-card style="min-width: 420px; max-width: 550px;" class="rounded-xl">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6 text-weight-bold row items-center">
            <q-icon name="contact_support" class="q-mr-sm" size="24px" /> Contatta la Segreteria
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-lg">
          <div class="text-subtitle2 text-grey-8 q-mb-md">
            Seleziona la tua scuola dal menu a tendina per visualizzare l'indirizzo email ed i recapiti della Segreteria Didattica.
          </div>

          <q-select
            v-model="selectedSchool"
            :options="schools"
            option-label="name"
            label="Seleziona la tua Scuola / Istituto *"
            outlined
            dense
            clearable
            :loading="loadingSchools"
            class="q-mb-lg"
          >
            <template v-slot:no-option>
              <q-item>
                <q-item-section class="text-grey">Nessuna scuola trovata</q-item-section>
              </q-item>
            </template>
            <template v-slot:option="scope">
              <q-item v-bind="scope.itemProps">
                <q-item-section avatar>
                  <q-icon name="school" color="primary" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ scope.opt.name }}</q-item-label>
                  <q-item-label caption v-if="scope.opt.code">Codice: {{ scope.opt.code }} • {{ scope.opt.city || 'Italia' }}</q-item-label>
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
                <span class="text-weight-medium q-mr-xs">Email Segreteria:</span>
                <a :href="'mailto:' + selectedSchool.email" class="text-primary text-weight-bold text-decoration-none">
                  {{ selectedSchool.email }}
                </a>
              </div>

              <div v-if="selectedSchool.phone" class="row items-center q-mt-xs">
                <q-icon name="phone" color="indigo-8" size="18px" class="q-mr-sm" />
                <span class="text-weight-medium q-mr-xs">Telefono:</span>
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
                label="Invia Email"
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
                label="Copia Email"
                no-caps
                @click="copyEmail(selectedSchool.email)"
                class="col"
              />
            </div>
          </div>

          <div v-else-if="!loadingSchools" class="text-center text-grey-6 q-py-md">
            <q-icon name="arrow_upward" size="24px" class="q-mb-xs" /><br />
            Scegli un istituto dal menu in alto per visualizzare i dettagli di contatto della Segreteria.
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md bg-grey-1">
          <q-btn flat label="Chiudi" color="grey-7" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuth } from '@/composables/useAuth'
import api from '@/services/api'

const $q = useQuasar()
const email = ref('')
const password = ref('')
const passwordInputRef = ref(null)
const showPassword = ref(false)
const rememberMe = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const { login } = useAuth()
const route = useRoute()

const showSecretaryDialog = ref(false)
const loadingSchools = ref(false)
const selectedSchool = ref(null)
const schools = ref([])

const defaultSchools = [
  { id: '1', name: 'Liceo Scientifico Statale Galileo Galilei', code: 'RMPS010001', email: 'rmps010001@istruzione.it', phone: '+39 06 12345678', city: 'Roma', address: 'Via delle Fornaci 200' },
  { id: '2', name: 'Liceo Ginnasio Statale Ennio Quirino Visconti', code: 'RMPC080007', email: 'rmpc080007@istruzione.it', phone: '+39 06 6793508', city: 'Roma', address: 'Piazza del Collegio Romano 4' },
  { id: '3', name: 'Istituto d\'Istruzione Superiore Camillo Cavour', code: 'RMIS00100X', email: 'rmis00100x@istruzione.it', phone: '+39 06 4880574', city: 'Roma', address: 'Via delle Carine 1' },
  { id: '4', name: 'Liceo Scientifico e Linguistico Guglielmo Marconi', code: 'BOPS01000V', email: 'bops01000v@istruzione.it', phone: '+39 051 6142145', city: 'Bologna', address: 'Via Maria Grazia Agnesi 1' },
  { id: '5', name: 'Scuola di Prova (Ambiente Demo)', code: 'PROVA123', email: 'segreteria.prova@scuola.it', phone: '+39 06 5551234', city: 'Roma', address: 'Via delle Prove 10' }
]

onMounted(() => {
  document.title = 'Accedi — Registro Elettronico'
  if (route?.query?.reason === 'session_expired') {
    errorMessage.value = 'Sessione scaduta. Effettua nuovamente l\'accesso.'
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

function copyEmail(emailStr) {
  if (!emailStr) return
  navigator.clipboard.writeText(emailStr)
  $q.notify({
    type: 'positive',
    icon: 'content_copy',
    message: 'Indirizzo email copiato negli appunti!'
  })
}

async function onSubmit() {
  errorMessage.value = ''
  loading.value = true
  const error = await login(email.value, password.value, rememberMe.value)
  loading.value = false
  
  if (error) {
    errorMessage.value = error
  }
}
</script>

<style scoped>
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
