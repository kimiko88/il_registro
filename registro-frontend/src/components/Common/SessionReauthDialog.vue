<template>
  <q-dialog
    v-model="showDialog"
    persistent
    no-backdrop-dismiss
    transition-show="scale"
    transition-hide="scale"
  >
    <q-card class="reauth-card q-pa-sm" style="min-width: 360px; max-width: 440px; border-radius: 18px;">

      <!-- Header -->
      <q-card-section class="reauth-header row items-center no-wrap q-pb-sm">
        <div class="reauth-icon-wrap q-mr-md">
          <q-icon name="lock_clock" size="32px" color="white" />
        </div>
        <div>
          <div class="text-h6 text-weight-bold text-white">Sessione scaduta</div>
          <div class="text-caption text-indigo-100">Inserisci la password per continuare</div>
        </div>
      </q-card-section>

      <!-- Body -->
      <q-card-section class="q-pt-md q-pb-xs">

        <!-- Info banner -->
        <div class="reauth-info-banner q-pa-sm q-mb-md row items-start">
          <q-icon name="info_outline" color="indigo-6" size="18px" class="q-mt-xs q-mr-sm flex-shrink-0" />
          <div class="text-caption text-indigo-9" style="line-height: 1.5">
            Il tuo lavoro nella pagina corrente è stato preservato.<br>
            Accedi di nuovo per continuare.
          </div>
        </div>

        <!-- Email (read-only) -->
        <q-input
          :model-value="userEmail"
          label="Email"
          outlined
          dense
          readonly
          bg-color="grey-1"
          class="q-mb-md rounded-input"
        >
          <template v-slot:prepend>
            <q-icon name="person" color="grey-6" />
          </template>
        </q-input>

        <!-- Password -->
        <q-input
          ref="passwordRef"
          v-model="password"
          label="Password"
          :type="showPassword ? 'text' : 'password'"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          autocomplete="current-password"
          :error="!!errorMessage"
          :error-message="errorMessage"
          @keyup.enter="onSubmit"
        >
          <template v-slot:prepend>
            <q-icon name="lock" color="primary" />
          </template>
          <template v-slot:append>
            <q-btn
              flat round dense
              :icon="showPassword ? 'visibility_off' : 'visibility'"
              :aria-label="showPassword ? 'Nascondi password' : 'Mostra password'"
              @click="showPassword = !showPassword"
            />
          </template>
        </q-input>

      </q-card-section>

      <!-- Actions -->
      <q-card-actions class="q-px-md q-pb-md q-pt-sm row q-gutter-sm">
        <q-btn
          flat
          label="Esci"
          color="grey-7"
          no-caps
          class="col"
          :disable="loading"
          @click="onCancel"
        />
        <q-btn
          label="Accedi"
          color="primary"
          no-caps
          unelevated
          class="col"
          :loading="loading"
          icon="login"
          @click="onSubmit"
        />
      </q-card-actions>

    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import { useSessionReauth } from '@/composables/useSessionReauth'
import { useAuthStore } from '@/stores/auth'
import authService from '@/services/authService'
import { useRouter } from 'vue-router'

const { showDialog, userEmail, resolveReauth, cancelReauth } = useSessionReauth()
const authStore = useAuthStore()
const router = useRouter()

const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const passwordRef = ref(null)

// Focus automatico sul campo password quando il dialog si apre
watch(showDialog, async (val) => {
  if (val) {
    password.value = ''
    errorMessage.value = ''
    showPassword.value = false
    await nextTick()
    passwordRef.value?.focus()
  }
})

async function onSubmit() {
  if (!password.value) {
    errorMessage.value = 'Inserisci la password.'
    return
  }
  errorMessage.value = ''
  loading.value = true

  try {
    const response = await authService.login(userEmail.value, password.value)
    const { access_token, user: userData } = response

    // Aggiorna il token in memoria senza toccare lo stato della pagina
    authStore.updateTokens(access_token, null, userData || authStore.user)

    resolveReauth(access_token)
  } catch (err) {
    const raw = String(err?.response?.data?.error || err?.message || '').toLowerCase()
    if (raw.includes('invalid') || raw.includes('credentials') || raw.includes('password')) {
      errorMessage.value = 'Password non corretta. Riprova.'
    } else if (raw.includes('too many') || raw.includes('rate limit')) {
      errorMessage.value = 'Troppi tentativi. Attendi qualche minuto.'
    } else {
      errorMessage.value = 'Errore durante l\'accesso. Riprova.'
    }
  } finally {
    loading.value = false
  }
}

function onCancel() {
  // L'utente sceglie di non autenticarsi: logout + redirect a /login
  cancelReauth()
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.reauth-card {
  box-shadow: 0 25px 60px rgba(79, 70, 229, 0.25), 0 8px 20px rgba(0,0,0,0.15);
}

.reauth-header {
  background: linear-gradient(135deg, #4F46E5 0%, #6366F1 100%);
  border-radius: 14px 14px 0 0;
  padding: 16px 20px;
}

.reauth-icon-wrap {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  width: 52px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.reauth-info-banner {
  background: #EEF2FF;
  border: 1px solid #C7D2FE;
  border-radius: 10px;
}

.rounded-input :deep(.q-field__control) {
  border-radius: 10px;
}
</style>
