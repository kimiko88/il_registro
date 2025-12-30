<template>
  <q-page class="flex flex-center">
    <div class="glass-card q-pa-xl" style="width: 100%; max-width: 420px">
      <div class="text-center q-mb-lg">
        <div class="text-h4 text-weight-bold text-primary q-mb-xs" style="letter-spacing: -1px">Welcome Back</div>
        <div class="text-grey-7">Sign in to continue to your dashboard</div>
      </div>

      <q-form @submit="onSubmit" class="q-gutter-y-md">
        <q-input
          v-model="email"
          label="Email Address"
          type="email"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          :rules="[val => !!val || 'Email is required']"
        >
          <template v-slot:prepend>
            <q-icon name="email" color="primary" />
          </template>
        </q-input>

        <q-input
          v-model="password"
          label="Password"
          type="password"
          outlined
          dense
          bg-color="white"
          class="rounded-input"
          :rules="[val => !!val || 'Password is required']"
        >
          <template v-slot:prepend>
            <q-icon name="lock" color="primary" />
          </template>
        </q-input>
        
        <div class="row justify-between items-center q-mt-sm">
          <q-checkbox v-model="rememberMe" label="Remember me" dense size="sm" color="primary" />
          <a href="#" class="text-caption text-primary text-weight-medium" style="text-decoration: none">Forgot password?</a>
        </div>

        <div class="q-mt-lg">
          <q-btn 
            label="Sign In" 
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
        Don't have an account? <span class="text-primary text-weight-medium cursor-pointer">Contact Administration</span>
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
const rememberMe = ref(false)
const loading = ref(false)
const { login } = useAuth()

async function onSubmit() {
  loading.value = true
  const error = await login(email.value, password.value)
  loading.value = false
  
  if (error) {
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
