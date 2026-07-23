<template>
  <q-page class="flex flex-center">
    <div class="glass-card q-pa-xl" style="width: 100%; max-width: 420px">
      <div class="text-center q-mb-lg">
        <div class="text-h4 text-weight-bold text-primary q-mb-xs" style="letter-spacing: -1px">Bentornato</div>
        <div class="text-grey-7">Accedi per entrare nel Registro Elettronico</div>
      </div>

      <q-form @submit="onSubmit" class="q-gutter-y-md">
        <q-input
          v-model="email"
          label="Indirizzo Email"
          type="email"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          :rules="[val => !!val || 'L\'email è obbligatoria']"
        >
          <template v-slot:prepend>
            <q-icon name="email" color="primary" />
          </template>
        </q-input>

        <q-input
          v-model="password"
          label="Password"
          :type="showPassword ? 'text' : 'password'"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          :rules="[val => !!val || 'La password è obbligatoria']"
        >
          <template v-slot:prepend>
            <q-icon name="lock" color="primary" />
          </template>
          <template v-slot:append>
            <q-icon
              :name="showPassword ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="showPassword = !showPassword"
              aria-label="Mostra o nascondi password"
            />
          </template>
        </q-input>
        
        <div class="row justify-between items-center q-mt-sm">
          <q-checkbox v-model="rememberMe" label="Ricordami" dense size="sm" color="primary" />
        </div>

        <!-- Inline Error Alert -->
        <div v-if="errorMessage" class="q-mt-sm bg-red-1 text-negative q-pa-sm rounded-lg text-caption text-center row items-center justify-center">
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
        Non hai un account? <span class="text-primary text-weight-medium cursor-pointer">Contatta la Segreteria</span>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { Notify } from 'quasar'

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const rememberMe = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const { login } = useAuth()

async function onSubmit() {
  errorMessage.value = ''
  loading.value = true
  const error = await login(email.value, password.value)
  loading.value = false
  
  if (error) {
    errorMessage.value = error
    Notify.create({
      type: 'negative',
      message: error,
      position: 'top',
      timeout: 3000
    })
  }
}
</script>

<style scoped>
.rounded-input :deep(.q-field__control) {
  border-radius: 12px;
}
</style>
