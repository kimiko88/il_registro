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
          :key="'drawer-toggle'"
        />

        <q-toolbar-title class="text-weight-bold text-primary">
          Registro Elettronico
        </q-toolbar-title>

        <q-space />
        
        <!-- Dark Mode Toggle -->
        <q-btn flat round dense :icon="$q.dark.isActive ? 'light_mode' : 'dark_mode'" @click="$q.dark.toggle()" color="primary" class="q-mr-sm" :key="'dark-toggle'">
           <q-tooltip>Toggle Dark Mode</q-tooltip>
        </q-btn>

        <!-- Fullscreen Toggle -->
        <q-btn 
          v-if="$q.fullscreen"
          flat 
          round 
          dense 
          :icon="$q.fullscreen.isActive ? 'fullscreen_exit' : 'fullscreen'" 
          @click="$q.fullscreen.toggle()" 
          color="primary" 
          class="q-mr-sm"
          :key="'fullscreen-toggle'"
        >
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
      <div class="column full-height no-wrap">
        <!-- User Profile Section -->
        <div class="q-pa-lg bg-gradient-premium text-white relative-position overflow-hidden" v-if="userName">
          <div class="row items-center q-mb-sm relative-position" style="z-index: 1">
            <q-avatar size="56px" color="white" text-color="primary" class="q-mr-md shadow-soft">
              <q-icon name="person" size="32px" />
            </q-avatar>
            <div class="col">
              <div class="text-h6 text-weight-bold no-wrap">{{ userName }}</div>
              <div class="text-caption opacity-80 text-uppercase letter-spacing-1">{{ roleLabel }}</div>
            </div>
          </div>
          <!-- Decorative Circle -->
          <div class="absolute-bottom-right q-mr-n-lg q-mb-n-lg" style="width: 120px; height: 120px; border-radius: 50%; background: rgba(255,255,255,0.1)"></div>
        </div>

        <!-- Menu Items -->
        <q-scroll-area class="col">
          <div class="q-pa-md">
            <div class="text-overline text-grey-5 q-px-md q-mb-sm letter-spacing-2">MENU PRINCIPALE</div>
            <q-list padding class="q-gutter-y-xs">
              <q-item 
                v-for="item in menuItems"
                :key="item.path"
                clickable 
                :to="item.path"
                :exact="item.exact"
                active-class="bg-indigo-50 text-indigo-700 active-menu-item"
                class="rounded-lg q-mx-sm transition-all"
              >
                <q-item-section avatar>
                  <q-icon :name="item.icon" size="22px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ item.label }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </div>
        </q-scroll-area>

        <!-- Logout Button at Bottom -->
        <div class="q-pa-md border-t border-slate-100">
          <q-item
            clickable
            class="rounded-lg bg-red-50 text-negative q-pa-md"
            @click="handleLogout"
            :disable="loggingOut"
          >
            <q-item-section avatar>
              <q-icon name="logout" size="20px" />
            </q-item-section>
            <q-item-section class="text-weight-bold">
              Esci
            </q-item-section>
            <q-item-section side v-if="loggingOut">
              <q-spinner size="20px" />
            </q-item-section>
          </q-item>
        </div>
      </div>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
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
  if (!userRole.value) return 'Utente'
  const roleLabels = {
    admin: 'Amministratore',
    secretary: 'Segretario',
    teacher: 'Docente',
    student: 'Studente',
    parent: 'Genitore'
  }
  return roleLabels[userRole.value] || userRole.value
})

// Get menu items based on role
const menuItems = ref([])
watch(userRole, (newRole) => {
  if (newRole) {
    menuItems.value = useMenuItems(newRole)
  } else {
    menuItems.value = []
  }
}, { immediate: true })

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
.bg-gradient-premium {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
}

.opacity-80 {
  opacity: 0.8;
}

.letter-spacing-1 {
    letter-spacing: 1px;
}

.letter-spacing-2 {
    letter-spacing: 2px;
}

.transition-all {
    transition: all 0.3s ease;
}

.active-menu-item {
    position: relative;
    box-shadow: inset 4px 0 0 #4f46e5;
}

.active-menu-item::after {
    content: '';
    position: absolute;
    right: 8px;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: #4f46e5;
}
</style>
