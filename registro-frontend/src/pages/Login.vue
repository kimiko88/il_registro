<template>
  <q-page class="flex flex-center">
    <q-card style="width: 400px">
      <q-card-section>
        <div class="text-h6">Login</div>
      </q-card-section>

      <q-card-section>
        <q-form @submit="onSubmit">
          <q-input
            v-model="email"
            label="Email"
            type="email"
            :rules="[val => !!val || 'Email is required']"
          />
          <q-input
            v-model="password"
            label="Password"
            type="password"
            :rules="[val => !!val || 'Password is required']"
          />
          
          <div class="row justify-end q-mt-md">
            <q-btn label="Login" type="submit" color="primary" :loading="loading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { Notify } from 'quasar'

const email = ref('')
const password = ref('')
const loading = ref(false)
const { login } = useAuth()

async function onSubmit() {
  loading.value = true
  const error = await login(email.value, password.value)
  loading.value = false
  
  if (error) {
    Notify.create({
      type: 'negative',
      message: error
    })
  }
}
</script>
