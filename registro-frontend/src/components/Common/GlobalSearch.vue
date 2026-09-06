<template>
  <!-- Activator overlay -->
  <teleport to="body">
    <transition name="search-fade">
      <div
        v-if="isOpen"
        class="global-search-overlay"
        @click.self="close"
        role="dialog"
        aria-modal="true"
        aria-label="Ricerca globale e comandi rapidi"
      >
        <div class="global-search-container">
          <!-- Search Input -->
          <div class="search-input-wrap row items-center q-px-md q-py-sm">
            <q-icon name="search" size="22px" color="primary" class="q-mr-sm" />
            <input
              ref="inputRef"
              v-model="query"
              class="search-native-input col"
              :placeholder="t('search.placeholder') || 'Cerca studenti, classi, circolari o azioni rapide... (Ctrl+K)'"
              autocomplete="off"
              spellcheck="false"
              @keydown.escape="close"
              @keydown.enter="selectHighlighted"
              @keydown.arrow-down.prevent="moveDown"
              @keydown.arrow-up.prevent="moveUp"
            />
            <q-btn flat round dense icon="close" size="sm" color="grey-6" :aria-label="t('common.close') || 'Chiudi'" @click="close" />
          </div>

          <q-separator />

          <!-- No input: Quick Actions & Navigation shortcuts -->
          <div v-if="!query" class="q-pa-sm">
            <div class="text-caption text-weight-bold text-grey-7 q-px-md q-pt-xs q-pb-xs letter-spacing-1 text-uppercase">
              {{ t('search.quickActions') || 'Azioni Rapide & Scorciatoie' }}
            </div>
            <q-list dense class="q-py-xs">
              <q-item
                v-for="(item, idx) in defaultItems"
                :key="item.id"
                clickable
                class="search-result-item rounded-lg q-mx-sm q-py-xs"
                :class="{ 'result-highlighted': highlightedIndex === idx }"
                @click="openResult(item)"
                @mouseenter="highlightedIndex = idx"
              >
                <q-item-section avatar min-width="36px">
                  <q-avatar :color="item.color || 'primary'" text-color="white" size="30px">
                    <q-icon :name="item.icon || 'bolt'" size="16px" />
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold text-body2">
                    {{ item.label }}
                  </q-item-label>
                  <q-item-label caption class="text-grey-6">{{ item.subtitle }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip dense outline size="xs" color="grey-6" class="text-weight-medium">
                    {{ item.action ? 'Azione' : 'Vai' }}
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>

            <div class="row justify-center q-pt-sm q-pb-xs q-gutter-xs text-caption text-grey-6">
              <q-chip dense outline color="grey-5" :label="'↑↓ ' + (t('search.navigate') || 'Naviga')" />
              <q-chip dense outline color="grey-5" :label="'Enter ' + (t('search.open') || 'Apri')" />
              <q-chip dense outline color="grey-5" :label="'Esc ' + (t('search.close') || 'Chiudi')" />
            </div>
          </div>

          <!-- Loading -->
          <div v-else-if="loading" class="q-pa-lg text-center">
            <q-spinner color="primary" size="28px" />
          </div>

          <!-- No results -->
          <div v-else-if="query && results.length === 0" class="q-pa-lg text-center text-grey-6">
            <q-icon name="search_off" size="32px" class="opacity-50 q-mb-sm" /><br />
            <span class="text-caption">{{ t('search.noResults') || 'Nessun risultato per' }} "<strong>{{ query }}</strong>"</span>
          </div>

          <!-- Results -->
          <q-scroll-area v-else style="max-height: 420px;">
            <q-list class="q-py-sm">
              <template v-for="(group, gi) in groupedResults" :key="gi">
                <q-item-label header class="text-caption text-weight-bold text-primary q-px-md q-pt-sm q-pb-xs letter-spacing-1 text-uppercase">
                  {{ group.label }}
                </q-item-label>
                <q-item
                  v-for="(item, ii) in group.items"
                  :key="item.id"
                  clickable
                  class="search-result-item rounded-lg q-mx-sm"
                  :class="{ 'result-highlighted': highlightedIndex === getGlobalIndex(gi, ii) }"
                  @click="openResult(item)"
                  @mouseenter="highlightedIndex = getGlobalIndex(gi, ii)"
                >
                  <q-item-section avatar min-width="36px">
                    <q-avatar :color="item.color || 'primary'" text-color="white" size="32px">
                      <q-icon :name="item.icon || 'search'" size="18px" />
                    </q-avatar>
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-body2">
                      <span v-html="highlight(item.label)"></span>
                    </q-item-label>
                    <q-item-label caption class="text-grey-6">{{ item.subtitle }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-icon name="arrow_forward_ios" size="14px" color="grey-4" />
                  </q-item-section>
                </q-item>
              </template>
            </q-list>
          </q-scroll-area>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar, debounce } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import api from '@/services/api'

// ── State ──────────────────────────────────────────────────────────────────
const { t } = useI18n()
const $q = useQuasar()
const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()

const isOpen = ref(false)
const query = ref('')
const loading = ref(false)
const highlightedIndex = ref(0)
const inputRef = ref(null)
const apiResults = ref([])

// ── Keyboard shortcut ───────────────────────────────────────────────────────
function onKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    if (isOpen.value) {
      close()
    } else {
      open()
    }
  }
}
onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))

function open() {
  isOpen.value = true
  query.value = ''
  highlightedIndex.value = 0
  apiResults.value = []
  nextTick(() => inputRef.value?.focus())
}
function close() {
  isOpen.value = false
  query.value = ''
  apiResults.value = []
}

// ── Role & Context ─────────────────────────────────────────────────────────
const role = computed(() => (authStore.userRole || authStore.user?.role || 'user').toLowerCase())

// Quick Actions
const quickActions = computed(() => {
  return [
    {
      id: 'qa-dark',
      label: $q.dark.isActive ? (t('layout.lightMode') || 'Passa a Modalità Chiara') : (t('layout.darkMode') || 'Passa a Modalità Scura'),
      subtitle: t('layout.themeAriaLabel') || 'Personalizzazione Aspetto',
      icon: $q.dark.isActive ? 'light_mode' : 'dark_mode',
      color: 'amber-9',
      group: t('search.quickActions') || 'Azioni Rapide',
      action: () => $q.dark.toggle()
    },
    {
      id: 'qa-contrast',
      label: t('layout.highContrast') || 'Contrasto Elevato',
      subtitle: t('layout.highContrastDesc') || 'Accessibilità Visiva WCAG 2.2',
      icon: 'contrast',
      color: 'blue-grey-8',
      group: t('search.quickActions') || 'Azioni Rapide',
      action: () => themeStore.toggleHighContrast()
    },
    {
      id: 'qa-shortcuts',
      label: t('a11y.shortcutsTitle') || 'Scorciatoie da Tastiera',
      subtitle: 'Guida comandi rapidi (?)',
      icon: 'keyboard',
      color: 'purple',
      group: t('search.quickActions') || 'Azioni Rapide',
      action: () => themeStore.toggleKeyboardShortcutsHelp(true)
    }
  ]
})

// Static menu items matching current role
const menuItems = computed(() => {
  const map = {
    teacher: [
      { id: 'm-1', label: t('nav.grades') || 'Gestione Voti', icon: 'grade', color: 'indigo', route: '/teacher/grades', subtitle: 'Valutazione e verifiche' },
      { id: 'm-2', label: t('nav.attendance') || 'Presenze e Registro', icon: 'how_to_reg', color: 'teal', route: '/teacher/attendance', subtitle: 'Appello giornaliero e lezioni' },
      { id: 'm-3', label: t('nav.agenda') || 'Agenda e Compiti', icon: 'edit_calendar', color: 'orange-8', route: '/teacher/agenda', subtitle: 'Pianificazione didattica' },
      { id: 'm-4', label: t('nav.scrutiny') || 'Scrutinio e Pagelle', icon: 'analytics', color: 'purple', route: '/teacher/scrutiny', subtitle: 'Valutazioni periodiche' },
      { id: 'm-5', label: t('nav.colloqui') || 'Colloqui con Genitori', icon: 'event', color: 'pink', route: '/teacher/colloqui', subtitle: 'Ricevimento famiglie' },
      { id: 'm-6', label: t('nav.communications') || 'Circolari e Comunicazioni', icon: 'email', color: 'blue-8', route: '/teacher/communications', subtitle: 'Bacheca avvisi' }
    ],
    student: [
      { id: 'm-s1', label: t('nav.myGrades') || 'I Miei Voti', icon: 'grade', color: 'indigo', route: '/student/grades', subtitle: 'Riepilogo valutazioni' },
      { id: 'm-s2', label: t('nav.attendance') || 'Presenze e Assenze', icon: 'event_available', color: 'teal', route: '/student/attendance', subtitle: 'Storico orario e giustifiche' },
      { id: 'm-s3', label: t('nav.homework') || 'Compiti e Lezioni', icon: 'assignment', color: 'orange-8', route: '/student/homework', subtitle: 'Scadenze e studio' },
      { id: 'm-s4', label: t('nav.communications') || 'Circolari e Avvisi', icon: 'email', color: 'blue-8', route: '/student/communications', subtitle: 'Comunicazioni scuola' },
      { id: 'm-s5', label: t('nav.reportCard') || 'Pagella e Documenti', icon: 'description', color: 'purple', route: '/student/report-card', subtitle: 'Scheda di valutazione' }
    ],
    parent: [
      { id: 'm-p1', label: t('nav.childrenGrades') || 'Voti dei Figli', icon: 'grade', color: 'indigo', route: '/parent/grades', subtitle: 'Riepilogo valutazioni' },
      { id: 'm-p2', label: t('nav.justifications') || 'Giustifiche e Assenze', icon: 'verified', color: 'teal', route: '/parent/justifications', subtitle: 'Giustifica assenze e ritardi' },
      { id: 'm-p3', label: t('nav.colloqui') || 'Prenota Colloqui', icon: 'event', color: 'pink', route: '/parent/colloqui', subtitle: 'Incontra i docenti' },
      { id: 'm-p4', label: t('nav.communications') || 'Circolari e Avvisi', icon: 'email', color: 'blue-8', route: '/parent/communications', subtitle: 'Bacheca della scuola' }
    ],
    secretary: [
      { id: 'm-sec1', label: t('nav.students') || 'Gestione Studenti', icon: 'school', color: 'indigo', route: '/secretary/students', subtitle: 'Anagrafiche e iscrizioni' },
      { id: 'm-sec2', label: t('nav.classes') || 'Gestione Classi', icon: 'room', color: 'teal', route: '/secretary/classes', subtitle: 'Sezioni e coordinamento' },
      { id: 'm-sec3', label: t('nav.users') || 'Gestione Utenti', icon: 'people', color: 'orange-8', route: '/secretary/users', subtitle: 'Docenti, ATA e account' },
      { id: 'm-sec4', label: t('nav.documents') || 'Protocollo e Documenti', icon: 'description', color: 'purple', route: '/secretary/documents', subtitle: 'Atti e certificati' },
      { id: 'm-sec5', label: t('nav.timetable') || 'Orario Scolastico', icon: 'schedule', color: 'blue-8', route: '/secretary/timetable', subtitle: 'Gestione orari lezioni' }
    ],
    admin: [
      { id: 'm-a1', label: t('nav.adminUsers') || 'Gestione Utenti', icon: 'people', color: 'indigo', route: '/admin/users', subtitle: 'Amministrazione account' },
      { id: 'm-a2', label: t('nav.adminSettings') || 'Impostazioni Sistema', icon: 'settings', color: 'grey-8', route: '/admin/settings', subtitle: 'Configurazione registro' },
      { id: 'm-a3', label: t('nav.auditLog') || 'Registro Audit e Log', icon: 'security', color: 'green-8', route: '/admin/audit-log', subtitle: 'Tracciamento attività' }
    ]
  }
  return map[role.value] || map.teacher
})

// Default items shown when search box is empty
const defaultItems = computed(() => {
  return [...quickActions.value, ...menuItems.value.slice(0, 4)]
})

// ── Search logic ────────────────────────────────────────────────────────────
function getRouteForResult(item) {
  const currentRole = role.value
  const type = (item.type || '').toLowerCase()

  if (type === 'communication') {
    if (currentRole === 'student') return '/student/communications'
    if (currentRole === 'parent') return '/parent/communications'
    if (currentRole === 'secretary') return '/secretary/communications'
    return '/teacher/communications'
  }
  if (type === 'lesson') {
    return currentRole === 'student' ? '/student/homework' : '/teacher/lessons'
  }
  if (type === 'class') {
    if (currentRole === 'secretary' || currentRole === 'admin') return '/secretary/classes'
    return '/teacher/attendance'
  }
  if (type === 'student') {
    if (currentRole === 'secretary') return '/secretary/students'
    if (currentRole === 'teacher') return '/teacher/grades'
    if (currentRole === 'parent') return '/parent/grades'
    return '/student/grades'
  }
  if (type === 'teacher') {
    if (currentRole === 'parent') return '/parent/colloqui'
    return '/secretary/users'
  }
  if (type === 'document') {
    return currentRole === 'secretary' ? '/secretary/documents' : '/teacher/documents'
  }
  return '/dashboard'
}

// Debounced API search: uses GET /search for all authenticated users, and /search/global only for superadmin
const doSearch = debounce(async (q) => {
  if (!q || q.length < 2) {
    apiResults.value = []
    loading.value = false
    return
  }
  loading.value = true
  try {
    const endpoint = role.value === 'superadmin' ? '/search/global' : '/search'
    const res = await api.get(endpoint, { params: { q, limit: 15 } })
    apiResults.value = res.data?.results || []
  } catch (_err) {
    apiResults.value = []
  } finally {
    loading.value = false
  }
}, 250)

watch(query, (q) => {
  loading.value = !!q && q.length >= 2
  doSearch(q)
  highlightedIndex.value = 0
})

// Grouped results
const results = computed(() => {
  const q = (query.value || '').toLowerCase().trim()
  if (!q) return []

  const items = []

  // 1. Matched Quick Actions
  const matchedActions = quickActions.value.filter(a =>
    a.label.toLowerCase().includes(q) || a.subtitle.toLowerCase().includes(q)
  )
  items.push(...matchedActions)

  // 2. Matched Menu Items
  const matchedMenu = menuItems.value.filter(m =>
    m.label.toLowerCase().includes(q) || m.subtitle?.toLowerCase().includes(q)
  ).map(m => ({
    ...m,
    group: t('search.navigation') || 'Navigazione & Pagine'
  }))
  items.push(...matchedMenu)

  // 3. Transform backend API search items
  for (const r of apiResults.value) {
    const groupName = r.type === 'communication'
      ? (t('nav.communications') || 'Circolari & Avvisi')
      : (r.type === 'lesson'
        ? (t('nav.lessons') || 'Lezioni & Didattica')
        : (r.type === 'class'
          ? (t('nav.classes') || 'Classi & Sezioni')
          : (r.type === 'student'
            ? (t('common.student') || 'Studenti')
            : (t('common.users') || 'Utenti & Anagrafiche'))))

    const iconName = r.type === 'communication'
      ? 'email'
      : (r.type === 'lesson'
        ? 'menu_book'
        : (r.type === 'class'
          ? 'room'
          : (r.type === 'student'
            ? 'school'
            : 'person')))

    const iconColor = r.type === 'communication'
      ? 'pink'
      : (r.type === 'lesson'
        ? 'indigo'
        : (r.type === 'class'
          ? 'orange-8'
          : (r.type === 'student'
            ? 'teal'
            : 'blue-8')))

    items.push({
      id: r.id || `api-${Math.random()}`,
      label: r.title || r.name,
      subtitle: r.description || r.subtitle || r.type,
      group: groupName,
      icon: iconName,
      color: iconColor,
      route: getRouteForResult(r)
    })
  }

  return items
})

const groupedResults = computed(() => {
  const groups = {}
  for (const item of results.value) {
    const g = item.group || (t('search.results') || 'Risultati')
    if (!groups[g]) groups[g] = []
    groups[g].push(item)
  }
  return Object.entries(groups).map(([label, items]) => ({ label, items }))
})

// ── Navigation helpers ───────────────────────────────────────────────────────
let _totalFlat = []
watch([groupedResults, defaultItems, query], () => {
  if (!query.value) {
    _totalFlat = defaultItems.value
  } else {
    _totalFlat = groupedResults.value.flatMap(g => g.items)
  }
}, { immediate: true })

function getGlobalIndex(gi, ii) {
  let idx = 0
  for (let i = 0; i < gi; i++) idx += groupedResults.value[i].items.length
  return idx + ii
}

function moveDown() {
  const maxIdx = (_totalFlat.length || 1) - 1
  highlightedIndex.value = Math.min(highlightedIndex.value + 1, maxIdx)
}
function moveUp() {
  highlightedIndex.value = Math.max(highlightedIndex.value - 1, 0)
}
function selectHighlighted() {
  const item = _totalFlat[highlightedIndex.value]
  if (item) openResult(item)
}

function openResult(item) {
  close()
  if (typeof item.action === 'function') {
    item.action()
  } else if (item.route) {
    router.push(item.route)
  } else if (item.url && /^https?:\/\//i.test(item.url)) {
    window.open(item.url, '_blank', 'noopener,noreferrer')
  }
}

// ── Highlight matching text ──────────────────────────────────────────────────
function escapeHtml(str) {
  return String(str || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

function highlight(text) {
  if (!text) return ''
  const safeText = escapeHtml(text)
  const q = (query.value || '').trim()
  if (!q) return safeText
  const escaped = escapeHtml(q).replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return safeText.replace(new RegExp(`(${escaped})`, 'gi'), '<mark class="search-highlight">$1</mark>')
}

// Expose open() so MainLayout can call it
defineExpose({ open, close, isOpen, query })
</script>

<style scoped>
.global-search-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(4px);
  z-index: 9999;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 80px;
}

.global-search-container {
  width: 100%;
  max-width: 640px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.25);
  overflow: hidden;
  animation: search-pop 0.18s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.body--dark .global-search-container {
  background: #1e293b;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.5);
}

.search-input-wrap {
  min-height: 56px;
}

.search-native-input {
  border: none;
  outline: none;
  background: transparent;
  font-size: 16px;
  color: inherit;
  width: 100%;
}
.body--dark .search-native-input {
  color: #e2e8f0;
}
.search-native-input::placeholder {
  color: #94a3b8;
}

.search-result-item {
  transition: background 0.12s;
  border-radius: 8px;
}

.result-highlighted {
  background: rgba(99, 102, 241, 0.1) !important;
}
.body--dark .result-highlighted {
  background: rgba(99, 102, 241, 0.22) !important;
}

:deep(.search-highlight) {
  background: #fef08a;
  color: #1e293b;
  border-radius: 2px;
  padding: 0 2px;
}
.body--dark :deep(.search-highlight) {
  background: #854d0e;
  color: #fef3c7;
}

@keyframes search-pop {
  from { opacity: 0; transform: scale(0.96) translateY(-12px); }
  to   { opacity: 1; transform: scale(1) translateY(0); }
}

.search-fade-enter-active,
.search-fade-leave-active {
  transition: opacity 0.15s ease;
}
.search-fade-enter-from,
.search-fade-leave-to {
  opacity: 0;
}
</style>
