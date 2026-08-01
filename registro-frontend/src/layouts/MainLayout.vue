<template>
  <q-layout view="hHh Lpr lFf">
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
        
        <!-- Theme Selector Menu -->
        <q-btn-dropdown
          flat
          round
          dense
          icon="palette"
          color="primary"
          class="q-mr-sm"
          key="theme-toggle"
          aria-label="Scegli il Tema Visivo"
        >
          <q-tooltip>Seleziona Tema Visivo</q-tooltip>
          <q-list style="min-width: 280px" class="q-py-xs">
            <q-item-label header class="text-weight-bold text-uppercase text-caption letter-spacing-1">
              Temi e Palette Visive
            </q-item-label>

            <q-item
              v-for="themeOption in THEMES"
              :key="themeOption.id"
              clickable
              v-close-popup
              @click="themeStore.setTheme(themeOption.id)"
              :active="themeStore.currentTheme === themeOption.id"
              active-class="bg-indigo-50 text-primary text-weight-bold"
              class="rounded-lg q-mx-xs q-mb-xs"
            >
              <q-item-section avatar>
                <q-avatar size="32px" :color="themeOption.badgeColor" text-color="white">
                  <q-icon :name="themeOption.icon" size="18px" />
                </q-avatar>
              </q-item-section>

              <q-item-section>
                <q-item-label class="text-weight-bold row items-center justify-between">
                  <span>{{ themeOption.name }}</span>
                  <q-chip
                    dense
                    size="xs"
                    color="grey-3"
                    text-color="grey-9"
                    class="text-weight-medium"
                  >
                    {{ themeOption.recommendedRole }}
                  </q-chip>
                </q-item-label>
                <q-item-label caption class="text-grey-7">
                  {{ themeOption.description }}
                </q-item-label>
              </q-item-section>

              <q-item-section side v-if="themeStore.currentTheme === themeOption.id">
                <q-icon name="check_circle" color="primary" size="20px" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>

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

        <!-- Notifications -->
        <q-btn flat round dense icon="notifications" color="primary" class="q-mr-sm" aria-label="Notifiche" @click="navigateToNotifications">
          <q-tooltip>Notifiche e Comunicazioni</q-tooltip>
        </q-btn>


        <q-btn flat round dense icon="account_circle" color="primary" aria-label="Profilo utente" @click="navigateToProfile" />

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
        <div class="q-pa-md bg-primary text-white relative-position overflow-hidden" v-if="userName" role="region" aria-label="Profilo utente">
          <div class="row items-center relative-position" style="z-index: 1">
            <q-avatar size="42px" color="white" text-color="primary" class="q-mr-md shadow-soft" aria-hidden="true">
              <q-icon name="person" size="24px" />
            </q-avatar>
            <div class="col">
              <div class="text-subtitle1 text-weight-bold no-wrap ellipsis" :aria-label="'Utente connesso: ' + userName">{{ userName }}</div>
              <div class="text-caption opacity-80 text-uppercase letter-spacing-1" :aria-label="'Ruolo: ' + roleLabel">{{ roleLabel }}</div>
            </div>
          </div>
          <!-- Decorative Circle -->
          <div class="absolute-bottom-right q-mr-n-lg q-mb-n-lg" style="width: 90px; height: 90px; border-radius: 50%; background: rgba(255,255,255,0.1)" aria-hidden="true"></div>
        </div>

        <!-- Menu Items -->
        <q-scroll-area class="col">
          <div class="q-pa-sm">
            <div class="text-overline text-grey-5 q-px-sm q-mb-xs letter-spacing-2" aria-hidden="true">MENU PRINCIPALE</div>
            <q-list dense padding class="q-gutter-y-xs" aria-label="Navigazione principale">
              <template v-for="(item, idx) in menuItems" :key="item.path || item.category || idx">
                <!-- Group Category with children -->
                <q-expansion-item
                  v-if="item.children"
                  group="menu-group"
                  :icon="item.icon"
                  :label="item.category"
                  dense
                  header-class="text-weight-bold text-slate-700 rounded-lg"
                  :default-opened="isCategoryActive(item)"
                >
                  <q-list dense class="q-pl-sm q-gutter-y-xs">
                    <q-item
                      v-for="child in item.children"
                      :key="child.path"
                      clickable
                      :to="child.path"
                      :exact="child.exact !== undefined ? child.exact : false"
                      active-class="active-menu-item"
                      class="rounded-lg transition-all"
                      :aria-label="child.label"
                    >
                      <q-item-section avatar min-width="32px">
                        <q-icon :name="child.icon" size="18px" aria-hidden="true" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">{{ child.label }}</q-item-label>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-expansion-item>

                <!-- Single Item -->
                <q-item
                  v-else
                  clickable
                  :to="item.path"
                  :exact="item.exact !== undefined ? item.exact : false"
                  active-class="active-menu-item"
                  class="rounded-lg transition-all"
                  :aria-label="item.label"
                >
                  <q-item-section avatar min-width="32px">
                    <q-icon :name="item.icon" size="20px" aria-hidden="true" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">{{ item.label }}</q-item-label>
                  </q-item-section>
                </q-item>
              </template>
            </q-list>
          </div>
        </q-scroll-area>

        <!-- Logout Button at Bottom -->
        <div class="q-pa-md border-t border-slate-100">
          <q-item
            clickable
            class="rounded-lg q-pa-md text-grey-8"
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

    <q-page-container role="main" id="main-content">
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
      <router-view v-slot="{ Component }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" :key="$route.path" />
        </transition>
      </router-view>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTeacherStore } from '@/stores/teacher'
import { useClassesStore } from '@/stores/classes'
import { useThemeStore, THEMES } from '@/stores/theme'
import { useAuth } from '@/composables/useAuth'
import { useMenuItems } from '@/composables/useMenuItems'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const route = useRoute()
const router = useRouter()
const $q = useQuasar()
const themeStore = useThemeStore()

const navigateToNotifications = () => {
  const role = userRole.value
  if (role === 'teacher') {
    router.push('/teacher/communications')
  } else if (role === 'student') {
    router.push('/student/communications')
  } else if (role === 'parent') {
    router.push('/parent/communications')
  } else {
    router.push('/secretary/communications')
  }
}

const navigateToProfile = () => {
  const role = userRole.value
  if (role === 'student') {
    router.push('/student/profile')
  } else if (role === 'parent') {
    router.push('/parent/profile')
  } else if (role === 'admin' || role === 'superadmin') {
    router.push('/admin/settings')
  } else if (role === 'secretary') {
    router.push('/secretary/settings')
  } else {
    router.push('/dashboard')
  }
}


onMounted(() => {
  themeStore.initTheme()
})



// Dynamic Breadcrumbs
const breadcrumbs = computed(() => {
  if (!route.path || route.path === '/dashboard' || route.path === '/') return []
  const items = []
  
  const routeNamesMap = {
    '/dashboard': { label: 'Dashboard', icon: 'dashboard' },
    '/profile': { label: 'Profilo Utente', icon: 'person' },
    '/teacher/grades': { label: 'Gestione Voti', icon: 'grade' },
    '/teacher/attendance': { label: 'Appello e Presenze', icon: 'how_to_reg' },
    '/teacher/timetable': { label: 'Orario Lezioni', icon: 'schedule' },
    '/teacher/didactics': { label: 'Materiale Didattico', icon: 'folder' },
    '/teacher/coordinator': { label: 'Area Coordinatore', icon: 'star' },
    '/teacher/groups': { label: 'Gruppi Linguistici', icon: 'groups' },
    '/teacher/rubrics': { label: 'Rubriche di Valutazione', icon: 'rule' },
    '/teacher/scrutiny': { label: 'Scrutini', icon: 'assessment' },
    '/teacher/communications': { label: 'Comunicazioni', icon: 'campaign' },
    '/teacher/agenda': { label: 'Agenda e Registo', icon: 'event' },
    '/teacher/classes': { label: 'Le Mie Classi', icon: 'class' },
    '/teacher/colloqui': { label: 'Colloqui e Incontri', icon: 'people' },
    '/teacher/documents': { label: 'Documenti', icon: 'description' },
    '/teacher/grade-weights': { label: 'Pesi Voti', icon: 'balance' },
    '/teacher/verbali': { label: 'Verbali', icon: 'gavel' },
    '/teacher/substitutions': { label: 'Sostituzioni', icon: 'swap_horiz' },
    '/student/grades': { label: 'I Miei Voti', icon: 'grade' },
    '/student/attendance': { label: 'Le Mie Presenze', icon: 'event_available' },
    '/student/homework': { label: 'Compiti', icon: 'assignment' },
    '/student/timetable': { label: 'Orario', icon: 'schedule' },
    '/parent/grades': { label: 'Voti Figlio', icon: 'grade' },
    '/parent/attendance': { label: 'Presenze e Giustifiche', icon: 'fact_check' },
    '/parent/communications': { label: 'Comunicazioni', icon: 'campaign' },
    '/parent/documents': { label: 'Documentazione', icon: 'folder_shared' },
    '/parent/payments': { label: 'Pagamenti', icon: 'payments' },
    '/parent/colloqui': { label: 'Incontri e Colloqui', icon: 'forum' },
    '/parent/meetings': { label: 'Riunioni', icon: 'groups' },
    '/parent/notes': { label: 'Note Disciplinari', icon: 'report' },
    '/parent/report-card': { label: 'Pagella Online', icon: 'assignment' },
    '/admin/school-settings': { label: 'Impostazioni Scuola', icon: 'settings' },
    '/admin/users': { label: 'Gestione Utenti', icon: 'people' },
    '/admin/analytics': { label: 'Analisi e Statistiche', icon: 'analytics' }
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
const { user, userName, userRole } = storeToRefs(authStore)
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

const teacherStore = useTeacherStore()
const classesStore = useClassesStore()

const isTeacherCoordinator = computed(() => {
  if (userRole.value !== 'teacher') return true
  if (teacherStore.isCoordinator) return true
  const currentUserId = user.value?.id
  if (currentUserId && classesStore.classes.some(c => c.coordinator_id === currentUserId)) {
    return true
  }
  return false
})

// Get menu items based on role
const menuItems = ref([])
watch([userRole, isTeacherCoordinator], ([newRole, isCoord]) => {
  if (!newRole) {
    menuItems.value = []
    return
  }
  let items = useMenuItems(newRole)
  if (newRole === 'teacher' && !isCoord) {
    items = items.map(cat => {
      if (!cat.children) return cat
      return {
        ...cat,
        children: cat.children.filter(child => !child.coordinatorOnly)
      }
    }).filter(cat => !cat.children || cat.children.length > 0)
  }
  menuItems.value = items
}, { immediate: true })

watch(userRole, (newRole) => {
  if (newRole === 'teacher') {
    classesStore.fetchAssignedClasses()
  }
}, { immediate: true })

const isCategoryActive = (category) => {
  if (!category || !category.children) return false
  return category.children.some(child => route.path === child.path || (child.path !== '/' && route.path.startsWith(child.path)))
}

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
</style>
