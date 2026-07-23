<template>
  <q-layout view="lHh Lpr lFf" class="bg-slate-50">
    <q-header class="glass-effect text-slate-900 q-py-xs" :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'" role="banner">
      <q-toolbar role="navigation" aria-label="Barra di navigazione principale">
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Apri/chiudi menu di navigazione"
          :aria-expanded="leftDrawerOpen"
          color="primary"
          @click="toggleLeftDrawer"
          :key="'drawer-toggle'"
        />

        <q-toolbar-title class="text-weight-bold text-primary">
          <span role="heading" aria-level="1">Registro Elettronico</span>
        </q-toolbar-title>

        <q-space />
        
        <!-- Dark Mode Toggle -->
        <q-btn flat round dense :icon="$q.dark.isActive ? 'light_mode' : 'dark_mode'" @click="$q.dark.toggle()" color="primary" class="q-mr-sm" :key="'dark-toggle'" :aria-label="$q.dark.isActive ? 'Attiva modalità chiara' : 'Attiva modalità scura'">
           <q-tooltip>Attiva/Disattiva Modalità Scura</q-tooltip>
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
          :aria-label="$q.fullscreen.isActive ? 'Esci da schermo intero' : 'Vai a schermo intero'"
        >
           <q-tooltip>Attiva/Disattiva Schermo Intero</q-tooltip>
        </q-btn>

        <div class="text-caption text-grey-6 q-mr-sm" aria-hidden="true">v0.0.1</div>
        <q-btn flat round dense icon="account_circle" color="primary" aria-label="Profilo utente" @click="$router.push('/profile')" />
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
      :width="260"
      role="navigation"
      aria-label="Menu laterale di navigazione"
    >
      <div class="column full-height no-wrap">
        <!-- User Profile Section -->
        <div class="q-pa-lg bg-gradient-premium text-white relative-position overflow-hidden" v-if="userName" role="region" aria-label="Profilo utente">
          <div class="row items-center q-mb-sm relative-position" style="z-index: 1">
            <q-avatar size="56px" color="white" text-color="primary" class="q-mr-md shadow-soft" aria-hidden="true">
              <q-icon name="person" size="32px" />
            </q-avatar>
            <div class="col">
              <div class="text-h6 text-weight-bold no-wrap" :aria-label="'Utente connesso: ' + userName">{{ userName }}</div>
              <div class="text-caption opacity-80 text-uppercase letter-spacing-1" :aria-label="'Ruolo: ' + roleLabel">{{ roleLabel }}</div>
            </div>
          </div>
          <!-- Decorative Circle -->
          <div class="absolute-bottom-right q-mr-n-lg q-mb-n-lg" style="width: 120px; height: 120px; border-radius: 50%; background: rgba(255,255,255,0.1)" aria-hidden="true"></div>
        </div>

        <!-- Menu Items -->
        <q-scroll-area class="col">
          <div class="q-pa-md">
            <div class="text-overline text-grey-5 q-px-md q-mb-sm letter-spacing-2" aria-hidden="true">MENU PRINCIPALE</div>
            <q-list padding class="q-gutter-y-xs" role="menubar" aria-label="Navigazione principale">
              <q-item 
                v-for="item in menuItems"
                :key="item.path"
                clickable 
                :to="item.path"
                :exact="item.exact"
                active-class="bg-indigo-50 text-indigo-700 active-menu-item"
                class="rounded-lg q-mx-sm transition-all"
                role="menuitem"
                :aria-label="item.label"
              >
                <q-item-section avatar>
                  <q-icon :name="item.icon" size="22px" aria-hidden="true" />
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
            role="button"
            aria-label="Esci dall'applicazione"
            :aria-busy="loggingOut"
          >
            <q-item-section avatar>
              <q-icon name="logout" size="20px" aria-hidden="true" />
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

    <q-page-container role="main">
      <!-- Dynamic Breadcrumb Navigation Header -->
      <div v-if="breadcrumbs.length > 0" class="q-px-md q-pt-md">
        <q-breadcrumbs class="text-caption text-grey-7" active-color="primary">
          <template v-slot:separator>
            <q-icon size="1.2em" name="chevron_right" color="grey-5" />
          </template>
          <q-breadcrumbs-el icon="home" to="/dashboard" label="Dashboard" />
          <q-breadcrumbs-el
            v-for="(crumb, idx) in breadcrumbs"
            :key="idx"
            :label="crumb.label"
            :to="crumb.path"
            :icon="crumb.icon"
          />
        </q-breadcrumbs>
      </div>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAuth } from '@/composables/useAuth'
import { useMenuItems } from '@/composables/useMenuItems'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const route = useRoute()
const $q = useQuasar()

// Dynamic Breadcrumbs
const breadcrumbs = computed(() => {
  if (!route.path || route.path === '/dashboard' || route.path === '/') return []
  const items = []
  
  const routeNamesMap = {
    '/teacher/grades': { label: 'Gestione Voti', icon: 'grade' },
    '/teacher/attendance': { label: 'Appello e Presenze', icon: 'how_to_reg' },
    '/teacher/timetable': { label: 'Orario Lezioni', icon: 'schedule' },
    '/teacher/didactics': { label: 'Materiale Didattico', icon: 'folder' },
    '/teacher/coordinator': { label: 'Area Coordinatore', icon: 'star' },
    '/teacher/groups': { label: 'Gruppi Linguistici', icon: 'groups' },
    '/teacher/rubrics': { label: 'Rubriche di Valutazione', icon: 'rule' },
    '/teacher/scrutiny': { label: 'Scrutini', icon: 'assessment' },
    '/student/grades': { label: 'I Miei Voti', icon: 'grade' },
    '/student/attendance': { label: 'Le Mie Presenze', icon: 'event_available' },
    '/student/homework': { label: 'Compiti', icon: 'assignment' },
    '/student/timetable': { label: 'Orario', icon: 'schedule' },
    '/parent/grades': { label: 'Voti Figlio', icon: 'grade' },
    '/parent/attendance': { label: 'Presenze e Giustifiche', icon: 'fact_check' },
    '/admin/school-settings': { label: 'Impostazioni Scuola', icon: 'settings' },
    '/admin/users': { label: 'Gestione Utenti', icon: 'people' }
  }

  const current = routeNamesMap[route.path] || { label: route.meta?.title || route.name || 'Pagina', icon: 'chevron_right' }

  if (route.path.startsWith('/teacher/')) {
    items.push({ label: 'Docente', icon: 'school' })
  } else if (route.path.startsWith('/student/')) {
    items.push({ label: 'Studente', icon: 'person' })
  } else if (route.path.startsWith('/parent/')) {
    items.push({ label: 'Genitore', icon: 'family_restroom' })
  } else if (route.path.startsWith('/admin/')) {
    items.push({ label: 'Amministrazione', icon: 'admin_panel_settings' })
  }

  items.push(current)
  return items
})

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
