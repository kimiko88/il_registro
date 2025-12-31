<template>
  <q-layout view="lHh Lpr lFf" class="bg-slate-50">
    <q-header class="glass-effect text-slate-900 q-py-xs" :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          color="primary"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title class="text-weight-bold text-primary">
          Registro Elettronico
        </q-toolbar-title>

        <q-space />
        
        <!-- Dark Mode Toggle -->
        <q-btn flat round dense :icon="$q.dark.isActive ? 'light_mode' : 'dark_mode'" @click="$q.dark.toggle()" color="primary" class="q-mr-sm">
           <q-tooltip>Toggle Dark Mode</q-tooltip>
        </q-btn>

        <!-- Fullscreen Toggle -->
        <q-btn flat round dense :icon="$q.fullscreen.isActive ? 'fullscreen_exit' : 'fullscreen'" @click="$q.fullscreen.toggle()" color="primary" class="q-mr-sm">
           <q-tooltip>Toggle Fullscreen</q-tooltip>
        </q-btn>

        <div class="text-caption text-grey-6 q-mr-sm">v0.0.1</div>
        <q-btn flat round dense icon="account_circle" color="primary" />
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
      :width="260"
    >
      <!-- User Profile Section -->
      <div class="q-pa-md bg-gradient-primary text-white">
        <div class="row items-center q-mb-sm">
          <q-avatar size="48px" color="white" text-color="primary" class="q-mr-md">
            <q-icon name="person" size="28px" />
          </q-avatar>
          <div class="col">
            <div class="text-weight-bold">{{ userName }}</div>
            <div class="text-caption opacity-80">{{ roleLabel }}</div>
          </div>
        </div>
      </div>

      <!-- Menu Items -->
      <div class="q-pa-md">
        <div class="text-overline text-grey-6 q-mb-sm">MENU</div>
        <q-list padding class="rounded-borders">
          <q-item 
            v-for="item in menuItems"
            :key="item.path"
            clickable 
            v-ripple
            :to="item.path"
            :exact="item.exact"
            active-class="bg-primary text-white rounded-lg shadow-soft"
          >
            <q-item-section avatar>
              <q-icon :name="item.icon" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-medium">{{ item.label }}</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </div>

      <!-- Logout Button at Bottom -->
      <div class="absolute-bottom q-pa-md">
        <q-btn
          outline
          color="negative"
          icon="logout"
          label="Esci"
          class="full-width"
          @click="handleLogout"
          :loading="loggingOut"
          no-caps
        />
      </div>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useAuth } from '@/composables/useAuth'
import { useMenuItems } from '@/composables/useMenuItems'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const authStore = useAuthStore()
const { userName, userRole } = storeToRefs(authStore)
const { logout } = useAuth()

const leftDrawerOpen = ref(false)
const loggingOut = ref(false)

// Get role label for display
const roleLabel = computed(() => {
  const roleLabels = {
    admin: 'Amministratore',
    secretary: 'Segretario',
    teacher: 'Docente',
    student: 'Studente',
    parent: 'Genitore'
  }
  return roleLabels[userRole.value] || 'Utente'
})

// Get menu items based on role
const menuItems = computed(() => {
  return useMenuItems(userRole.value)
})

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

async function handleLogout() {
  loggingOut.value = true
  try {
    await logout()
    $q.notify({
      type: 'positive',
      message: 'Logout effettuato con successo',
      position: 'top',
      timeout: 2000
    })
  } catch (error) {
    $q.notify({
      type: 'negative',
      message: 'Errore durante il logout',
      position: 'top',
      timeout: 3000
    })
  } finally {
    loggingOut.value = false
  }
}
</script>

<style scoped>
.bg-gradient-primary {
  background: linear-gradient(135deg, #4F46E5 0%, #3B82F6 100%);
}

.opacity-80 {
  opacity: 0.8;
}
</style>
