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
        aria-label="Ricerca globale"
      >
        <div class="global-search-container">
          <!-- Search Input -->
          <div class="search-input-wrap row items-center q-px-md q-py-sm">
            <q-icon name="search" size="22px" color="primary" class="q-mr-sm" />
            <input
              ref="inputRef"
              v-model="query"
              class="search-native-input col"
              :placeholder="t('search.placeholder')"
              autocomplete="off"
              spellcheck="false"
              @keydown.escape="close"
              @keydown.enter="selectHighlighted"
              @keydown.arrow-down.prevent="moveDown"
              @keydown.arrow-up.prevent="moveUp"
            />
            <q-btn flat round dense icon="close" size="sm" color="grey-6" @click="close" />
          </div>

          <q-separator />

          <!-- No input hint -->
          <div v-if="!query" class="search-hint q-pa-lg text-center text-grey-6">
            <q-icon name="keyboard" size="32px" class="q-mb-sm opacity-50" /><br />
            <span class="text-caption">{{ t('search.hint') }}</span>
            <div class="row justify-center q-mt-md q-gutter-sm">
              <q-chip dense outline color="grey-5" :label="'↑↓ ' + t('search.navigate')" />
              <q-chip dense outline color="grey-5" :label="'Enter ' + t('search.open')" />
              <q-chip dense outline color="grey-5" :label="'Esc ' + t('search.close')" />
            </div>
          </div>

          <!-- Loading -->
          <div v-else-if="loading" class="q-pa-lg text-center">
            <q-spinner color="primary" size="28px" />
          </div>

          <!-- No results -->
          <div v-else-if="query && results.length === 0" class="q-pa-lg text-center text-grey-6">
            <q-icon name="search_off" size="32px" class="opacity-50 q-mb-sm" /><br />
            <span class="text-caption">{{ t('search.noResults') }} "<strong>{{ query }}</strong>"</span>
          </div>

          <!-- Results -->
          <q-scroll-area v-else style="max-height: 420px;">
            <q-list class="q-py-sm">
              <template v-for="(group, gi) in groupedResults" :key="gi">
                <q-item-label header class="text-caption text-weight-bold text-primary q-px-md q-pt-sm q-pb-xs letter-spacing-1">
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
                  <q-item-section avatar>
                    <q-avatar :color="item.color || 'primary'" text-color="white" size="34px">
                      <q-icon :name="item.icon || 'search'" size="18px" />
                    </q-avatar>
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">
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
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import { debounce } from 'quasar'
import { useI18n } from 'vue-i18n'

// ── State ──────────────────────────────────────────────────────────────────
const { t } = useI18n()
const isOpen = ref(false)
const query = ref('')
const loading = ref(false)
const highlightedIndex = ref(0)
const inputRef = ref(null)

const apiResults = ref([])

const router = useRouter()
const authStore = useAuthStore()

// ── Keyboard shortcut ───────────────────────────────────────────────────────
function onKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    open()
  }
}
onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))

function open() {
  isOpen.value = true
  query.value = ''
  highlightedIndex.value = 0
  nextTick(() => inputRef.value?.focus())
}
function close() {
  isOpen.value = false
  query.value = ''
  apiResults.value = []
}

// ── Search logic ────────────────────────────────────────────────────────────
const role = computed(() => authStore.userRole)

// Static menu items that match the current role
const menuItems = computed(() => {
  const base = {
    teacher: [
      { id: 'm-1', label: 'Gestione Voti', icon: 'grade', color: 'indigo', route: '/teacher/grades', subtitle: 'Menu → Valutazione' },
      { id: 'm-2', label: 'Presenze e Registro', icon: 'how_to_reg', color: 'teal', route: '/teacher/attendance', subtitle: 'Menu → Didattica' },
      { id: 'm-3', label: 'Agenda e Impegni', icon: 'edit_calendar', color: 'orange', route: '/teacher/agenda', subtitle: 'Menu → Organizzazione' },
      { id: 'm-4', label: 'Scrutinio e Pagelle', icon: 'analytics', color: 'purple', route: '/teacher/scrutiny', subtitle: 'Menu → Didattica' },
      { id: 'm-5', label: 'Colloqui con Genitori', icon: 'event', color: 'pink', route: '/teacher/colloqui', subtitle: 'Menu → Organizzazione' },
      { id: 'm-6', label: 'Registro Supplenze e Sostituzioni', icon: 'swap_horiz', color: 'cyan', route: '/teacher/substitutions', subtitle: 'Menu → Organizzazione' },
      { id: 'm-7', label: 'Programmazione UdA', icon: 'auto_stories', color: 'deep-orange', route: '/teacher/uda', subtitle: 'Menu → Curricolo' },
      { id: 'm-8', label: 'Valutazione Competenze', icon: 'stars', color: 'blue-8', route: '/teacher/competencies', subtitle: 'Menu → Curricolo' },
    ],
    secretary: [
      { id: 'm-1', label: 'Gestione Studenti', icon: 'school', color: 'indigo', route: '/secretary/students', subtitle: 'Menu → Anagrafiche' },
      { id: 'm-2', label: 'Gestione Utenti e Segreteria', icon: 'people', color: 'teal', route: '/secretary/users', subtitle: 'Menu → Anagrafiche' },
      { id: 'm-3', label: 'Gestione Classi', icon: 'room', color: 'orange', route: '/secretary/classes', subtitle: 'Menu → Anagrafiche' },
      { id: 'm-4', label: 'Documenti e Fascicoli', icon: 'description', color: 'purple', route: '/secretary/documents', subtitle: 'Menu → Atti' },
      { id: 'm-5', label: 'Comunicazioni e Circolari', icon: 'email', color: 'pink', route: '/secretary/communications', subtitle: 'Menu → Servizi' },
    ],
    admin: [
      { id: 'm-1', label: 'Gestione Utenti', icon: 'people', color: 'indigo', route: '/admin/users', subtitle: 'Menu → Admin' },
      { id: 'm-2', label: 'Impostazioni Sistema', icon: 'settings', color: 'grey-8', route: '/admin/settings', subtitle: 'Menu → Admin' },
      { id: 'm-3', label: 'Analytics e Log', icon: 'analytics', color: 'green-8', route: '/admin/analytics', subtitle: 'Menu → Admin' },
    ]
  }
  return base[role.value] || base.teacher || []
})

// Debounced API search
const doSearch = debounce(async (q) => {
  if (!q || q.length < 2) { apiResults.value = []; loading.value = false; return }
  loading.value = true
  try {
    const res = await api.get('/search/global', { params: { q, limit: 12 } })
    apiResults.value = res.data?.results || []
  } catch {
    apiResults.value = []
  } finally {
    loading.value = false
  }
}, 280)

watch(query, (q) => {
  loading.value = !!q && q.length >= 2
  doSearch(q)
  highlightedIndex.value = 0
})

// ── Grouped results ─────────────────────────────────────────────────────────
const results = computed(() => {
  const q = query.value.toLowerCase()
  if (!q) return []

  const items = []

  // Transform backend API search items
  for (const r of apiResults.value) {
    const typeLabel = r.type === 'communication' ? 'Bacheca & Circolari' : (r.type === 'lesson' ? 'Lezioni & Argomenti' : (r.type === 'class' ? 'Classi' : 'Utenti & Anagrafiche'))
    const iconName = r.type === 'communication' ? 'email' : (r.type === 'lesson' ? 'menu_book' : (r.type === 'class' ? 'room' : (r.type === 'student' ? 'school' : 'person')))
    const iconColor = r.type === 'communication' ? 'pink' : (r.type === 'lesson' ? 'indigo' : (r.type === 'class' ? 'orange' : 'teal'))
    const targetRoute = r.type === 'communication' ? '/teacher/communications' : (r.type === 'lesson' ? '/teacher/lessons' : (r.type === 'class' ? '/secretary/classes' : '/secretary/users'))

    items.push({
      id: r.id,
      label: r.title || r.name,
      subtitle: r.description || r.subtitle || r.type,
      group: typeLabel,
      icon: iconName,
      color: iconColor,
      route: targetRoute
    })
  }

  // Filter menu items locally
  const menuMatches = menuItems.value.filter(m =>
    m.label.toLowerCase().includes(q) || m.subtitle?.toLowerCase().includes(q)
  )
  items.push(...menuMatches)

  return items
})

const groupedResults = computed(() => {
  const groups = {}
  for (const item of results.value) {
    const g = item.group || (item.route?.startsWith('/teacher') ? 'Voce di Menu' : 'Risultati')
    if (!groups[g]) groups[g] = []
    groups[g].push(item)
  }
  return Object.entries(groups).map(([label, items]) => ({ label, items }))
})

// ── Navigation helpers ───────────────────────────────────────────────────────
let _totalFlat = []
watch(groupedResults, (gr) => {
  _totalFlat = gr.flatMap(g => g.items)
})

function getGlobalIndex(gi, ii) {
  let idx = 0
  for (let i = 0; i < gi; i++) idx += groupedResults.value[i].items.length
  return idx + ii
}

function moveDown() {
  highlightedIndex.value = Math.min(highlightedIndex.value + 1, (_totalFlat.length || 1) - 1)
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
  if (item.route) {
    router.push(item.route)
  } else if (item.url) {
    window.open(item.url, '_blank')
  }
}

// ── Highlight matching text ──────────────────────────────────────────────────
function highlight(text) {
  if (!query.value || !text) return text || ''
  const escaped = query.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return text.replace(new RegExp(`(${escaped})`, 'gi'), '<mark class="search-highlight">$1</mark>')
}

// Expose open() so MainLayout can call it
defineExpose({ open, close })
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
  background: rgba(99, 102, 241, 0.08) !important;
}
.body--dark .result-highlighted {
  background: rgba(99, 102, 241, 0.18) !important;
}

.search-hint {
  min-height: 160px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
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
