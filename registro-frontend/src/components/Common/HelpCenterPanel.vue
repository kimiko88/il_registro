<template>
  <!-- Full-screen Help Center dialog -->
  <q-dialog
    v-model="isOpen"
    maximized
    transition-show="slide-up"
    transition-hide="slide-down"
    class="help-center-dialog"
  >
    <div class="help-center-root" :class="{ 'help-dark': $q.dark.isActive }">

      <!-- ── TOP BAR ──────────────────────────────────────────────────────── -->
      <div class="hc-topbar">
        <div class="hc-topbar-left">
          <div class="hc-logo">
            <q-icon name="help_center" size="26px" color="primary" />
          </div>
          <div>
            <div class="hc-title">{{ t('guideCenter.title') }}</div>
            <div class="hc-subtitle">{{ t('guideCenter.subtitle') }}</div>
          </div>
        </div>

        <div class="hc-topbar-center">
          <q-input
            v-model="searchQuery"
            :placeholder="t('guideCenter.search')"
            dense
            outlined
            clearable
            class="hc-search"
            bg-color="transparent"
          >
            <template #prepend>
              <q-icon name="search" color="grey-5" size="18px" />
            </template>
            <template #append v-if="!searchQuery">
              <kbd class="search-kbd">F</kbd>
            </template>
          </q-input>
        </div>

        <div class="hc-topbar-right">
          <q-btn
            flat
            :label="t('help.restartTour')"
            icon="rocket_launch"
            color="primary"
            size="sm"
            class="q-mr-sm"
            @click="restartTour"
          />
          <q-btn
            flat
            round
            dense
            icon="close"
            color="grey-6"
            @click="close"
            :aria-label="t('common.close', 'Chiudi')"
          />
        </div>
      </div>

      <!-- ── MAIN BODY ────────────────────────────────────────────────────── -->
      <div class="hc-body">

        <!-- SIDEBAR: categories -->
        <aside class="hc-sidebar">
          <div class="sidebar-role-badge">
            <q-icon :name="roleIcon" size="18px" :color="roleColor" class="q-mr-xs" />
            <span>{{ roleLabel }}</span>
          </div>

          <div class="sidebar-section-label">{{ t('help.categories') }}</div>

          <nav class="sidebar-nav" role="navigation">
            <button
              v-for="cat in categories"
              :key="cat.key"
              class="sidebar-item"
              :class="{ active: activeCategory === cat.key }"
              @click="selectCategory(cat.key)"
            >
              <div class="sidebar-item-icon" :class="`si-${cat.color}`">
                <q-icon :name="cat.icon" size="16px" />
              </div>
              <span class="sidebar-item-label">{{ cat.label }}</span>
              <span class="sidebar-item-count">{{ cat.count }}</span>
            </button>
          </nav>

          <!-- Restart Tour CTA -->
          <div class="sidebar-footer">
            <q-btn
              unelevated
              color="primary"
              icon="rocket_launch"
              :label="t('help.restartTour')"
              size="sm"
              class="full-width sidebar-tour-btn"
              @click="restartTour"
            />
            <q-btn
              flat
              color="grey-6"
              icon="support_agent"
              :label="t('help.contactSupport')"
              size="sm"
              class="full-width q-mt-xs"
              @click="goToSupport"
            />
          </div>
        </aside>

        <!-- MAIN CONTENT -->
        <main class="hc-content" ref="contentRef">

          <!-- Search results mode -->
          <template v-if="searchQuery">
            <div class="hc-section-header">
              <h2 class="hc-section-title">
                <q-icon name="search" size="20px" class="q-mr-xs" />
                "{{ searchQuery }}"
              </h2>
              <span class="hc-section-count">{{ filteredGuides.length }} {{ t('help.categories') }}</span>
            </div>

            <div v-if="filteredGuides.length === 0" class="hc-empty">
              <q-icon name="search_off" size="64px" color="grey-4" />
              <h3>{{ t('guideCenter.noResults') }} "{{ searchQuery }}"</h3>
              <p>Prova con altre parole chiave.</p>
              <q-btn flat color="primary" label="Cancella ricerca" @click="searchQuery = ''" />
            </div>

            <div class="guides-grid" v-else>
              <div
                v-for="guide in filteredGuides"
                :key="guide.key"
                class="guide-card"
                :class="{ 'guide-card-active': activeGuide?.key === guide.key }"
                @click="openGuide(guide)"
                tabindex="0"
                role="button"
                @keydown.enter="openGuide(guide)"
              >
                <div class="guide-card-icon" :class="`si-${guide.color}`">
                  <q-icon :name="guide.icon" size="22px" />
                </div>
                <div class="guide-card-body">
                  <div class="guide-card-title">{{ guide.title }}</div>
                  <div class="guide-card-desc">{{ guide.desc }}</div>
                  <div class="guide-card-meta">
                    <q-icon name="schedule" size="13px" color="grey-5" class="q-mr-xs" />
                    {{ guide.readTime }} {{ t('guideCenter.readingTime') }}
                  </div>
                </div>
                <q-icon name="chevron_right" size="18px" color="grey-4" class="guide-card-arrow" />
              </div>
            </div>
          </template>

          <!-- Category view -->
          <template v-else>
            <div class="hc-section-header">
              <div class="hc-section-header-left">
                <div class="hc-section-icon" :class="`si-${activeCategoryMeta?.color}`">
                  <q-icon :name="activeCategoryMeta?.icon || 'help'" size="20px" />
                </div>
                <div>
                  <h2 class="hc-section-title">{{ activeCategoryMeta?.label }}</h2>
                  <p class="hc-section-desc">{{ activeCategoryMeta?.description }}</p>
                </div>
              </div>
            </div>

            <!-- Guide cards grid -->
            <div class="guides-grid">
              <div
                v-for="guide in guidesInCategory"
                :key="guide.key"
                class="guide-card"
                :class="{ 'guide-card-active': activeGuide?.key === guide.key }"
                @click="openGuide(guide)"
                tabindex="0"
                role="button"
                @keydown.enter="openGuide(guide)"
              >
                <div class="guide-card-icon" :class="`si-${guide.color}`">
                  <q-icon :name="guide.icon" size="22px" />
                </div>
                <div class="guide-card-body">
                  <div class="guide-card-title">{{ guide.title }}</div>
                  <div class="guide-card-desc">{{ guide.desc }}</div>
                  <div class="guide-card-meta">
                    <q-icon name="schedule" size="13px" color="grey-5" class="q-mr-xs" />
                    {{ guide.readTime }} {{ t('guideCenter.readingTime') }}
                  </div>
                </div>
                <q-icon name="chevron_right" size="18px" color="grey-4" class="guide-card-arrow" />
              </div>
            </div>

            <!-- Quick FAQ cards -->
            <div class="hc-faq-section" v-if="faqsInCategory.length">
              <div class="hc-faq-label">
                <q-icon name="quiz" size="16px" color="primary" class="q-mr-xs" />
                Domande Frequenti
              </div>
              <div class="faq-list">
                <q-expansion-item
                  v-for="faq in faqsInCategory"
                  :key="faq.key"
                  :label="faq.question"
                  class="faq-item"
                  expand-separator
                  header-class="faq-header"
                  dense
                >
                  <div class="faq-answer">
                    <q-icon name="chevron_right" size="16px" color="primary" class="q-mr-xs" style="flex-shrink:0;margin-top:2px" />
                    <p>{{ faq.answer }}</p>
                  </div>
                </q-expansion-item>
              </div>
            </div>
          </template>
        </main>

        <!-- ARTICLE PANEL (slides in from right) -->
        <transition name="article-slide">
          <aside v-if="activeGuide" class="hc-article-panel">
            <div class="article-topbar">
              <button class="article-back" @click="activeGuide = null">
                <q-icon name="arrow_back" size="18px" />
                <span>{{ t('help.categories') }}</span>
              </button>
              <q-btn flat round dense icon="close" size="sm" color="grey-5" @click="activeGuide = null" />
            </div>

            <!-- Article header -->
            <div class="article-header" :class="`art-${activeGuide.color}`">
              <div class="article-icon">
                <q-icon :name="activeGuide.icon" size="36px" color="white" />
              </div>
              <div>
                <div class="article-category">{{ activeCategoryMeta?.label }}</div>
                <h2 class="article-title">{{ activeGuide.title }}</h2>
                <div class="article-meta">
                  <q-icon name="schedule" size="13px" color="rgba(255,255,255,0.7)" class="q-mr-xs" />
                  {{ activeGuide.readTime }} {{ t('guideCenter.readingTime') }}
                </div>
              </div>
            </div>

            <!-- Article body -->
            <div class="article-body" ref="articleBodyRef">
              <p class="article-intro">{{ activeGuide.desc }}</p>

              <!-- Render parsed content -->
              <div v-for="(block, i) in parsedContent" :key="i" class="content-block">
                <!-- Step -->
                <div v-if="block.type === 'step'" class="content-step">
                  <div class="step-num-badge" :class="`badge-${activeGuide.color}`">
                    {{ block.num }}
                  </div>
                  <div class="step-body">
                    <p>{{ block.text }}</p>
                  </div>
                </div>

                <!-- Tip -->
                <div v-else-if="block.type === 'tip'" class="content-tip">
                  <q-icon name="lightbulb" color="amber-7" size="18px" class="q-mr-sm" style="flex-shrink:0;margin-top:2px" />
                  <div>
                    <strong>{{ t('guideCenter.tip') }}:</strong>
                    {{ block.text }}
                  </div>
                </div>

                <!-- Warning -->
                <div v-else-if="block.type === 'warning'" class="content-warning">
                  <q-icon name="warning" color="orange-8" size="18px" class="q-mr-sm" style="flex-shrink:0;margin-top:2px" />
                  <div>
                    <strong>{{ t('guideCenter.warning') }}:</strong>
                    {{ block.text }}
                  </div>
                </div>

                <!-- Shortcut -->
                <div v-else-if="block.type === 'shortcut'" class="content-shortcut">
                  <q-icon name="keyboard" color="indigo" size="18px" class="q-mr-sm" style="flex-shrink:0;margin-top:2px" />
                  <div>
                    <strong>{{ t('guideCenter.shortcut') }}:</strong>
                    <kbd class="kbd-tag">{{ block.text }}</kbd>
                  </div>
                </div>

                <!-- Plain paragraph -->
                <p v-else class="content-para">{{ block.text }}</p>
              </div>
            </div>

            <!-- Article footer: related guides -->
            <div class="article-footer" v-if="relatedGuides.length">
              <div class="related-label">Guide correlate</div>
              <div class="related-list">
                <button
                  v-for="rel in relatedGuides"
                  :key="rel.key"
                  class="related-item"
                  @click="openGuide(rel)"
                >
                  <q-icon :name="rel.icon" size="16px" color="primary" class="q-mr-xs" />
                  {{ rel.title }}
                  <q-icon name="arrow_forward" size="14px" color="grey-4" class="q-ml-auto" />
                </button>
              </div>
            </div>
          </aside>
        </transition>

      </div>
    </div>
  </q-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const emit = defineEmits(['restart-tour'])
const router = useRouter()
const { t } = useI18n()
const $q = useQuasar()
const authStore = useAuthStore()

const isOpen = ref(false)
const searchQuery = ref('')
const activeCategory = ref('')
const activeGuide = ref(null)
const contentRef = ref(null)
const articleBodyRef = ref(null)

// ── Role ──────────────────────────────────────────────────────────────────────
const userRole = computed(() => {
  const role = authStore.userRole || authStore.user?.role || 'student'
  const r = role.toLowerCase()
  if (r === 'superadmin') return 'admin'
  return ['teacher', 'student', 'parent', 'secretary', 'admin'].includes(r) ? r : 'student'
})

const roleMeta = {
  teacher:   { icon: 'school',              color: 'indigo',  label: 'roles.teacher' },
  student:   { icon: 'face',                color: 'teal',    label: 'roles.student' },
  parent:    { icon: 'family_restroom',      color: 'purple',  label: 'roles.parent' },
  secretary: { icon: 'admin_panel_settings', color: 'orange',  label: 'roles.secretary' },
  admin:     { icon: 'manage_accounts',      color: 'red',     label: 'roles.admin' }
}
const roleIcon  = computed(() => roleMeta[userRole.value]?.icon  || 'person')
const roleColor = computed(() => roleMeta[userRole.value]?.color || 'primary')
const roleLabel = computed(() => t(roleMeta[userRole.value]?.label || 'roles.user'))

// ── Guide definitions ─────────────────────────────────────────────────────────
const GUIDE_DEFS = {
  teacher: {
    categories: [
      { key: 'dashboard',       label: 'Dashboard',         icon: 'dashboard',       color: 'indigo', description: 'Tutto sulla tua area di lavoro principale' },
      { key: 'attendance',      label: 'Presenze',          icon: 'event_available', color: 'blue',   description: 'Registra e gestisci le presenze' },
      { key: 'grades',          label: 'Voti',              icon: 'grade',           color: 'green',  description: 'Inserimento e analisi delle valutazioni' },
      { key: 'agenda',          label: 'Agenda',            icon: 'event',           color: 'orange', description: 'Pianifica compiti, verifiche e attività' },
      { key: 'communications',  label: 'Comunicazioni',     icon: 'campaign',        color: 'purple', description: 'Messaggi a studenti e genitori' },
      { key: 'meetings',        label: 'Colloqui',          icon: 'people',          color: 'pink',   description: 'Gestione ricevimento e colloqui' },
      { key: 'scrutiny',        label: 'Scrutinio',         icon: 'gavel',           color: 'amber',  description: 'Scrutini intermedi e finali del CdC' },
      { key: 'pdp',             label: 'PDP / PEI',         icon: 'psychology',      color: 'teal',   description: 'Piani didattici personalizzati per BES/DSA' }
    ],
    guides: {
      dashboard:      { key: 'dashboard',       icon: 'dashboard',       color: 'indigo', readTime: 3 },
      attendance:     { key: 'attendance',      icon: 'event_available', color: 'blue',   readTime: 4 },
      grades:         { key: 'grades',          icon: 'grade',           color: 'green',  readTime: 5 },
      agenda:         { key: 'agenda',          icon: 'event',           color: 'orange', readTime: 4 },
      communications: { key: 'communications',  icon: 'campaign',        color: 'purple', readTime: 3 },
      meetings:       { key: 'meetings',        icon: 'people',          color: 'pink',   readTime: 4 },
      scrutiny:       { key: 'scrutiny',        icon: 'gavel',           color: 'amber',  readTime: 6 },
      pdp:            { key: 'pdp',             icon: 'psychology',      color: 'teal',   readTime: 5 }
    }
  },
  student: {
    categories: [
      { key: 'dashboard',  label: 'Dashboard',   icon: 'dashboard',       color: 'teal',   description: 'La tua pagina principale' },
      { key: 'grades',     label: 'Voti',        icon: 'grade',           color: 'green',  description: 'Consulta e analizza i tuoi voti' },
      { key: 'homework',   label: 'Compiti',     icon: 'assignment',      color: 'orange', description: 'Gestisci compiti e scadenze' },
      { key: 'attendance', label: 'Presenze',    icon: 'event_available', color: 'blue',   description: 'Traccia assenze e ritardi' },
      { key: 'documents',  label: 'Documenti',   icon: 'description',     color: 'purple', description: 'Pagella e documenti rilasciati' },
      { key: 'simulator',  label: 'Simulatore',  icon: 'calculate',       color: 'indigo', description: 'Simula medie e voti target' },
      { key: 'pcto',       label: 'PCTO',        icon: 'work',            color: 'brown',  description: 'Ore alternanza scuola-lavoro' }
    ],
    guides: {
      dashboard:  { key: 'dashboard',  icon: 'dashboard',       color: 'teal',   readTime: 2 },
      grades:     { key: 'grades',     icon: 'grade',           color: 'green',  readTime: 3 },
      homework:   { key: 'homework',   icon: 'assignment',      color: 'orange', readTime: 3 },
      attendance: { key: 'attendance', icon: 'event_available', color: 'blue',   readTime: 3 },
      documents:  { key: 'documents',  icon: 'description',     color: 'purple', readTime: 3 },
      simulator:  { key: 'simulator',  icon: 'calculate',       color: 'indigo', readTime: 4 },
      pcto:       { key: 'pcto',       icon: 'work',            color: 'brown',  readTime: 4 }
    }
  },
  parent: {
    categories: [
      { key: 'monitoring',     label: 'Monitoraggio',     icon: 'monitor_heart', color: 'purple', description: 'Tieni d\'occhio i progressi dei tuoi figli' },
      { key: 'meetings',       label: 'Colloqui',         icon: 'people',        color: 'orange', description: 'Prenota e gestisci i colloqui' },
      { key: 'communications', label: 'Comunicazioni',    icon: 'campaign',      color: 'blue',   description: 'Circolari e comunicazioni scuola' },
      { key: 'pagopa',         label: 'Pagamenti PagoPA', icon: 'payments',      color: 'green',  description: 'Gestione rette e contributi' },
      { key: 'documents',      label: 'Documenti',        icon: 'folder',        color: 'indigo', description: 'Pagelle e modulistica' },
      { key: 'justifications', label: 'Giustificazioni',  icon: 'fact_check',    color: 'teal',   description: 'Giustifica assenze e ritardi' }
    ],
    guides: {
      monitoring:     { key: 'monitoring',     icon: 'monitor_heart', color: 'purple', readTime: 4 },
      meetings:       { key: 'meetings',       icon: 'people',        color: 'orange', readTime: 3 },
      communications: { key: 'communications', icon: 'campaign',      color: 'blue',   readTime: 3 },
      pagopa:         { key: 'pagopa',         icon: 'payments',      color: 'green',  readTime: 4 },
      documents:      { key: 'documents',      icon: 'folder',        color: 'indigo', readTime: 3 },
      justifications: { key: 'justifications', icon: 'fact_check',    color: 'teal',   readTime: 3 }
    }
  },
  secretary: {
    categories: [
      { key: 'students',       label: 'Studenti',     icon: 'group',        color: 'orange', description: 'Gestione anagrafica studenti' },
      { key: 'classes',        label: 'Classi',        icon: 'school',       color: 'blue',   description: 'Organizzazione classi ed elenchi' },
      { key: 'certificates',   label: 'Certificati',  icon: 'verified',     color: 'green',  description: 'Generazione documenti ufficiali' },
      { key: 'timetable',      label: 'Orario & Sost.',icon: 'schedule',     color: 'purple', description: 'Orario scolastico e sostituzioni' },
      { key: 'communications', label: 'Comunicazioni', icon: 'campaign',     color: 'indigo', description: 'Circolari e bacheca digitale' },
      { key: 'reports',        label: 'Report',       icon: 'assessment',   color: 'teal',   description: 'Estrazioni dati ministeriali' }
    ],
    guides: {
      students:       { key: 'students',       icon: 'group',        color: 'orange', readTime: 5 },
      classes:        { key: 'classes',        icon: 'school',       color: 'blue',   readTime: 4 },
      certificates:   { key: 'certificates',   icon: 'verified',     color: 'green',  readTime: 4 },
      timetable:      { key: 'timetable',      icon: 'schedule',     color: 'purple', readTime: 5 },
      communications: { key: 'communications', icon: 'campaign',     color: 'indigo', readTime: 4 },
      reports:        { key: 'reports',        icon: 'assessment',   color: 'teal',   readTime: 5 }
    }
  },
  admin: {
    categories: [
      { key: 'monitoring',   label: 'Monitoraggio', icon: 'monitor_heart',   color: 'red',         description: 'Stato e performance del sistema' },
      { key: 'users',        label: 'Utenti',       icon: 'manage_accounts', color: 'blue',        description: 'Gestione account e permessi' },
      { key: 'schools',      label: 'Scuole',       icon: 'domain',          color: 'amber',       description: 'Gestione istituti e plessi' },
      { key: 'security',     label: 'Sicurezza',    icon: 'security',        color: 'deep-orange', description: 'Policy password e 2FA' },
      { key: 'analytics',    label: 'Analytics',    icon: 'insights',        color: 'indigo',      description: 'Statistiche e dispersione scolastica' },
      { key: 'integrations', label: 'Integrazioni', icon: 'extension',       color: 'teal',        description: 'Classroom, Teams, SSO' },
      { key: 'audit',        label: 'Audit Log',    icon: 'policy',          color: 'purple',      description: 'Tracciabilità e log immutabili' }
    ],
    guides: {
      monitoring:   { key: 'monitoring',   icon: 'monitor_heart',   color: 'red',         readTime: 4 },
      users:        { key: 'users',        icon: 'manage_accounts', color: 'blue',        readTime: 5 },
      schools:      { key: 'schools',      icon: 'domain',          color: 'amber',       readTime: 4 },
      security:     { key: 'security',     icon: 'security',        color: 'deep-orange', readTime: 4 },
      analytics:    { key: 'analytics',    icon: 'insights',        color: 'indigo',      readTime: 5 },
      integrations: { key: 'integrations', icon: 'extension',       color: 'teal',        readTime: 5 },
      audit:        { key: 'audit',        icon: 'policy',          color: 'purple',      readTime: 4 }
    }
  }
}

// FAQ definitions per ruolo (from help.* i18n keys)
const FAQ_CATEGORY_MAP = {
  teacher:   { dashboard: [1,2], attendance: [3,4], grades: [5,6], agenda: [7], communications: [8], meetings: [9], scrutiny: [10], pdp: [] },
  student:   { dashboard: [1], grades: [2,3], homework: [4,5], attendance: [6], documents: [7], simulator: [8,9], pcto: [10] },
  parent:    { monitoring: [1,2], meetings: [3,4], communications: [5], pagopa: [6,7], documents: [8], justifications: [9,10] },
  secretary: { students: [1,2], classes: [3], certificates: [4,5], timetable: [6,7], communications: [8], reports: [9,10] },
  admin:     { monitoring: [1], users: [2,3], schools: [4], security: [5,6], analytics: [7], integrations: [8], audit: [9,10] }
}

const categories = computed(() => {
  const role = userRole.value
  const def = GUIDE_DEFS[role] || GUIDE_DEFS.student
  return def.categories.map(cat => {
    const guideKeys = Object.keys(def.guides).filter(k => k === cat.key)
    const faqIdxs = (FAQ_CATEGORY_MAP[role] || {})[cat.key] || []
    return { ...cat, count: guideKeys.length + faqIdxs.length }
  })
})

const activeCategoryMeta = computed(() =>
  categories.value.find(c => c.key === activeCategory.value)
)

// Set default category on open
watch(isOpen, val => {
  if (val && !activeCategory.value) {
    activeCategory.value = categories.value[0]?.key || ''
  }
})

const guidesInCategory = computed(() => {
  const role = userRole.value
  const def = GUIDE_DEFS[role] || GUIDE_DEFS.student
  const catKey = activeCategory.value
  const guideKeys = Object.keys(def.guides).filter(k => k === catKey)
  return guideKeys.map(k => {
    const g = def.guides[k]
    const titleKey = `guideCenter.${role}.${k}.title`
    const descKey  = `guideCenter.${role}.${k}.desc`
    return {
      key:      `${role}-${k}`,
      sectionKey: k,
      icon:     g.icon,
      color:    g.color,
      title:    t(titleKey),
      desc:     t(descKey),
      content:  t(`guideCenter.${role}.${k}.content`),
      readTime: g.readTime
    }
  })
})

const faqsInCategory = computed(() => {
  const role = userRole.value
  const catKey = activeCategory.value
  const idxs = (FAQ_CATEGORY_MAP[role] || {})[catKey] || []
  return idxs.map(i => ({
    key: `faq-${role}-${i}`,
    question: t(`help.${role}.q${i}`),
    answer:   t(`help.${role}.a${i}`)
  })).filter(f => f.question && !f.question.startsWith('help.'))
})

// All guides (for search)
const allGuides = computed(() => {
  const role = userRole.value
  const def = GUIDE_DEFS[role] || GUIDE_DEFS.student
  return Object.keys(def.guides).map(k => {
    const g = def.guides[k]
    return {
      key:        `${role}-${k}`,
      sectionKey: k,
      icon:       g.icon,
      color:      g.color,
      title:      t(`guideCenter.${role}.${k}.title`),
      desc:       t(`guideCenter.${role}.${k}.desc`),
      content:    t(`guideCenter.${role}.${k}.content`),
      readTime:   g.readTime
    }
  })
})

const filteredGuides = computed(() => {
  if (!searchQuery.value) return allGuides.value
  const q = searchQuery.value.toLowerCase()
  return allGuides.value.filter(g =>
    g.title.toLowerCase().includes(q) ||
    g.desc.toLowerCase().includes(q)  ||
    g.content.toLowerCase().includes(q)
  )
})

const relatedGuides = computed(() => {
  if (!activeGuide.value) return []
  return allGuides.value
    .filter(g => g.key !== activeGuide.value.key)
    .slice(0, 3)
})

// ── Content parser ─────────────────────────────────────────────────────────────
const parsedContent = computed(() => {
  if (!activeGuide.value?.content) return []
  const raw = activeGuide.value.content
  const lines = raw.split('\n').filter(l => l.trim())
  return lines.map(line => {
    const trimmed = line.trim()
    // Detect "Passaggio N:" pattern
    const stepMatch = trimmed.match(/^Passaggio\s+(\d+):\s*(.+)$/i)
    if (stepMatch) return { type: 'step', num: stepMatch[1], text: stepMatch[2] }
    // Detect "Suggerimento:"
    const tipMatch = trimmed.match(/^Suggerimento:\s*(.+)$/i)
    if (tipMatch) return { type: 'tip', text: tipMatch[1] }
    // Detect "Attenzione:"
    const warnMatch = trimmed.match(/^Attenzione:\s*(.+)$/i)
    if (warnMatch) return { type: 'warning', text: warnMatch[1] }
    // Detect "Scorciatoia:"
    const shcMatch = trimmed.match(/^Scorciatoia:\s*(.+)$/i)
    if (shcMatch) return { type: 'shortcut', text: shcMatch[1] }
    // Default: paragraph (skip "Step N:" prefixes for EN locale)
    const enStepMatch = trimmed.match(/^Step\s+(\d+):\s*(.+)$/i)
    if (enStepMatch) return { type: 'step', num: enStepMatch[1], text: enStepMatch[2] }
    return { type: 'para', text: trimmed }
  })
})

// ── Actions ───────────────────────────────────────────────────────────────────
function selectCategory(key) {
  activeCategory.value = key
  activeGuide.value = null
  searchQuery.value = ''
  nextTick(() => contentRef.value?.scrollTo({ top: 0, behavior: 'smooth' }))
}

function openGuide(guide) {
  activeGuide.value = guide
  nextTick(() => articleBodyRef.value?.scrollTo({ top: 0, behavior: 'smooth' }))
}

function restartTour() {
  close()
  emit('restart-tour')
}

function goToSupport() {
  close()
  router.push('/support')
}

function open(catKey) {
  isOpen.value = true
  if (catKey) {
    nextTick(() => { activeCategory.value = catKey })
  }
}

function close() {
  isOpen.value = false
  activeGuide.value = null
}

defineExpose({ open, close })
</script>

<style scoped>
/* ── ROOT ────────────────────────────────────────────────────────── */
.help-center-dialog { backdrop-filter: blur(4px); }

.help-center-root {
  display: flex;
  flex-direction: column;
  height: 100dvh;
  background: #f8f9ff;
  font-family: inherit;
}
.help-dark { background: #16162a; }

/* ── TOPBAR ──────────────────────────────────────────────────────── */
.hc-topbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 24px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
  flex-shrink: 0;
  z-index: 10;
}
.help-dark .hc-topbar {
  background: #1e1e2e;
  border-color: #3a3a4e;
}

.hc-topbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}
.hc-logo {
  width: 42px; height: 42px;
  background: linear-gradient(135deg, #667eea15, #764ba215);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #667eea22;
}
.hc-title {
  font-size: 1rem;
  font-weight: 800;
  color: #1a1a2e;
}
.help-dark .hc-title { color: #e0e0e0; }
.hc-subtitle {
  font-size: 0.75rem;
  color: #9ca3af;
}

.hc-topbar-center { flex: 1; max-width: 480px; }
.hc-search { border-radius: 10px; }

.search-kbd {
  font-size: 0.65rem;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  padding: 2px 6px;
  color: #6b7280;
  font-family: monospace;
}

.hc-topbar-right {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

/* ── BODY ────────────────────────────────────────────────────────── */
.hc-body {
  display: flex;
  flex: 1;
  overflow: hidden;
  position: relative;
}

/* ── SIDEBAR ─────────────────────────────────────────────────────── */
.hc-sidebar {
  width: 240px;
  flex-shrink: 0;
  background: white;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  padding: 16px 12px;
}
.help-dark .hc-sidebar {
  background: #1e1e2e;
  border-color: #3a3a4e;
}

.sidebar-role-badge {
  display: flex;
  align-items: center;
  font-size: 0.8rem;
  font-weight: 700;
  color: #6b7280;
  padding: 6px 8px;
  margin-bottom: 16px;
  background: #f8f9ff;
  border-radius: 8px;
}
.help-dark .sidebar-role-badge { background: #252535; color: #9ca3af; }

.sidebar-section-label {
  font-size: 0.68rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: #d1d5db;
  padding: 0 8px;
  margin-bottom: 8px;
}

.sidebar-nav { display: flex; flex-direction: column; gap: 3px; flex: 1; }

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 10px;
  border: none;
  background: transparent;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
  width: 100%;
}
.sidebar-item:hover { background: #f3f4f6; }
.help-dark .sidebar-item:hover { background: #252535; }
.sidebar-item.active {
  background: #eef2ff;
  border: 1px solid #c7d2fe;
}
.help-dark .sidebar-item.active { background: #252550; border-color: #667eea44; }

.sidebar-item-icon {
  width: 30px; height: 30px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.sidebar-item-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  flex: 1;
}
.help-dark .sidebar-item-label { color: #c0c0d0; }
.sidebar-item.active .sidebar-item-label { color: #4338ca; font-weight: 700; }

.sidebar-item-count {
  font-size: 0.72rem;
  background: #f3f4f6;
  color: #9ca3af;
  border-radius: 999px;
  padding: 1px 7px;
  min-width: 22px;
  text-align: center;
}
.sidebar-item.active .sidebar-item-count { background: #c7d2fe; color: #4338ca; }

.sidebar-footer {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid #f3f4f6;
}
.help-dark .sidebar-footer { border-color: #3a3a4e; }
.sidebar-tour-btn { border-radius: 10px !important; font-weight: 700 !important; }

/* ── MAIN CONTENT ────────────────────────────────────────────────── */
.hc-content {
  flex: 1;
  overflow-y: auto;
  padding: 28px 32px;
}

.hc-section-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
}
.hc-section-header-left { display: flex; align-items: flex-start; gap: 14px; }
.hc-section-icon {
  width: 44px; height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.hc-section-title {
  font-size: 1.3rem;
  font-weight: 800;
  color: #1a1a2e;
  margin: 0 0 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.help-dark .hc-section-title { color: #e0e0e0; }
.hc-section-desc {
  font-size: 0.88rem;
  color: #6b7280;
  margin: 0;
}
.hc-section-count {
  font-size: 0.8rem;
  color: #9ca3af;
  white-space: nowrap;
  padding-top: 6px;
}

/* Guide cards */
.guides-grid { display: grid; grid-template-columns: 1fr; gap: 12px; margin-bottom: 28px; }

.guide-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: white;
  border: 1.5px solid #e5e7eb;
  border-radius: 14px;
  padding: 16px 20px;
  cursor: pointer;
  transition: all 0.2s;
}
.help-dark .guide-card { background: #1e1e2e; border-color: #3a3a4e; }
.guide-card:hover {
  border-color: #667eea55;
  transform: translateX(4px);
  box-shadow: 0 4px 16px rgba(102,126,234,0.1);
}
.guide-card-active {
  border-color: #667eea;
  background: #f8f9ff;
}
.help-dark .guide-card-active { background: #20203a; }

.guide-card-icon {
  width: 48px; height: 48px;
  border-radius: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.guide-card-body { flex: 1; }
.guide-card-title {
  font-size: 0.97rem;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 4px;
}
.help-dark .guide-card-title { color: #e0e0e0; }
.guide-card-desc {
  font-size: 0.82rem;
  color: #6b7280;
  line-height: 1.5;
  margin-bottom: 6px;
}
.guide-card-meta {
  display: flex;
  align-items: center;
  font-size: 0.75rem;
  color: #d1d5db;
}
.guide-card-arrow { flex-shrink: 0; }

/* FAQ section */
.hc-faq-section { margin-top: 8px; }
.hc-faq-label {
  display: flex;
  align-items: center;
  font-size: 0.82rem;
  font-weight: 700;
  color: #6b7280;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 1px;
}
.faq-list { display: flex; flex-direction: column; gap: 4px; }
.faq-item {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  overflow: hidden;
}
.help-dark .faq-item { background: #1e1e2e; border-color: #3a3a4e; }
.faq-answer {
  display: flex;
  align-items: flex-start;
  padding: 12px 16px 16px;
  font-size: 0.88rem;
  color: #6b7280;
  line-height: 1.7;
  background: #f8f9ff;
}
.help-dark .faq-answer { background: #252535; }
.faq-answer p { margin: 0; }

/* Empty state */
.hc-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 64px 16px;
  gap: 12px;
  color: #9ca3af;
  text-align: center;
}
.hc-empty h3 { margin: 0; font-weight: 700; color: #4b5563; }
.hc-empty p { margin: 0; color: #9ca3af; }

/* ── ARTICLE PANEL ───────────────────────────────────────────────── */
.hc-article-panel {
  width: 440px;
  flex-shrink: 0;
  background: white;
  border-left: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.help-dark .hc-article-panel { background: #1e1e2e; border-color: #3a3a4e; }

.article-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #f3f4f6;
}
.help-dark .article-topbar { border-color: #3a3a4e; }

.article-back {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.83rem;
  font-weight: 500;
  color: #667eea;
  border: none;
  background: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}
.article-back:hover { background: #eef2ff; }

/* Article header gradient */
.article-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 24px 20px;
  color: white;
}
.art-indigo  { background: linear-gradient(135deg, #667eea, #4338ca); }
.art-teal    { background: linear-gradient(135deg, #0d9488, #0f766e); }
.art-purple  { background: linear-gradient(135deg, #a855f7, #7c3aed); }
.art-orange  { background: linear-gradient(135deg, #f97316, #ea580c); }
.art-red     { background: linear-gradient(135deg, #ef4444, #dc2626); }
.art-green   { background: linear-gradient(135deg, #22c55e, #16a34a); }
.art-blue    { background: linear-gradient(135deg, #3b82f6, #2563eb); }
.art-grey    { background: linear-gradient(135deg, #9ca3af, #6b7280); }
.art-cyan    { background: linear-gradient(135deg, #22d3ee, #0891b2); }
.art-pink    { background: linear-gradient(135deg, #ec4899, #db2777); }
.art-amber   { background: linear-gradient(135deg, #f59e0b, #d97706); }

.article-icon {
  background: rgba(255,255,255,0.2);
  border-radius: 14px;
  width: 64px; height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid rgba(255,255,255,0.25);
}
.article-category {
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 1.5px;
  text-transform: uppercase;
  opacity: 0.7;
  margin-bottom: 4px;
}
.article-title {
  font-size: 1.2rem;
  font-weight: 900;
  margin: 0 0 6px;
  line-height: 1.2;
}
.article-meta {
  display: flex;
  align-items: center;
  font-size: 0.75rem;
  opacity: 0.7;
}

/* Article body */
.article-body {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
.article-intro {
  font-size: 0.92rem;
  color: #374151;
  line-height: 1.7;
  margin: 0 0 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f3f4f6;
}
.help-dark .article-intro { color: #c0c0d0; border-color: #3a3a4e; }

.content-block { margin-bottom: 14px; }

.content-step {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  background: #f8f9ff;
  border-radius: 10px;
}
.help-dark .content-step { background: #252535; }
.step-num-badge {
  width: 28px; height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 800;
  flex-shrink: 0;
  color: white;
}
.badge-indigo  { background: #667eea; }
.badge-teal    { background: #0d9488; }
.badge-purple  { background: #a855f7; }
.badge-orange  { background: #f97316; }
.badge-red     { background: #ef4444; }
.badge-green   { background: #22c55e; }
.badge-blue    { background: #3b82f6; }
.badge-grey    { background: #9ca3af; }
.badge-cyan    { background: #22d3ee; }
.badge-pink    { background: #ec4899; }
.badge-amber   { background: #f59e0b; }

.step-body { font-size: 0.875rem; color: #374151; line-height: 1.6; }
.help-dark .step-body { color: #c0c0d0; }
.step-body p { margin: 0; }

.content-tip {
  display: flex;
  align-items: flex-start;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 10px;
  padding: 12px 14px;
  font-size: 0.875rem;
  color: #92400e;
  line-height: 1.6;
}
.help-dark .content-tip { background: #2d2a1e; border-color: #92400e44; color: #fde68a; }

.content-warning {
  display: flex;
  align-items: flex-start;
  background: #fff7ed;
  border: 1px solid #fed7aa;
  border-radius: 10px;
  padding: 12px 14px;
  font-size: 0.875rem;
  color: #9a3412;
  line-height: 1.6;
}
.help-dark .content-warning { background: #2d1e0f; border-color: #9a341244; color: #fed7aa; }

.content-shortcut {
  display: flex;
  align-items: center;
  background: #eef2ff;
  border: 1px solid #c7d2fe;
  border-radius: 10px;
  padding: 10px 14px;
  font-size: 0.875rem;
  color: #4338ca;
}
.help-dark .content-shortcut { background: #1e1e40; border-color: #4338ca44; color: #a5b4fc; }

.kbd-tag {
  display: inline-block;
  margin-left: 8px;
  font-family: monospace;
  font-size: 0.8rem;
  background: white;
  border: 1px solid #c7d2fe;
  border-radius: 4px;
  padding: 2px 8px;
  color: #4338ca;
}

.content-para {
  font-size: 0.875rem;
  color: #6b7280;
  line-height: 1.7;
  margin: 0;
}
.help-dark .content-para { color: #9ca3af; }

/* Related guides footer */
.article-footer {
  border-top: 1px solid #f3f4f6;
  padding: 16px 20px;
}
.help-dark .article-footer { border-color: #3a3a4e; }
.related-label {
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #d1d5db;
  margin-bottom: 10px;
}
.related-list { display: flex; flex-direction: column; gap: 6px; }
.related-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  background: #f8f9ff;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-size: 0.83rem;
  font-weight: 500;
  color: #374151;
  width: 100%;
  transition: background 0.2s;
  text-align: left;
}
.help-dark .related-item { background: #252535; color: #c0c0d0; }
.related-item:hover { background: #eef2ff; }

/* ── Shared icon colors ───────────────────────────────────────────── */
.si-indigo  { background: #eef2ff; color: #4338ca; }
.si-teal    { background: #f0fdfa; color: #0d9488; }
.si-purple  { background: #faf5ff; color: #7c3aed; }
.si-orange  { background: #fff7ed; color: #ea580c; }
.si-red     { background: #fef2f2; color: #dc2626; }
.si-green   { background: #f0fdf4; color: #16a34a; }
.si-blue    { background: #eff6ff; color: #2563eb; }
.si-grey    { background: #f9fafb; color: #6b7280; }
.si-cyan    { background: #ecfeff; color: #0891b2; }
.si-pink    { background: #fdf2f8; color: #db2777; }
.si-amber   { background: #fffbeb; color: #d97706; }

/* ── Transition: article panel ───────────────────────────────────── */
.article-slide-enter-active { transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
.article-slide-leave-active { transition: all 0.2s ease; }
.article-slide-enter-from   { opacity: 0; transform: translateX(32px); }
.article-slide-leave-to     { opacity: 0; transform: translateX(32px); }

/* ── Responsive ──────────────────────────────────────────────────── */
@media (max-width: 900px) {
  .hc-article-panel { width: 100%; position: absolute; inset: 0; z-index: 20; }
}
@media (max-width: 640px) {
  .hc-sidebar { display: none; }
  .hc-topbar-center { display: none; }
  .hc-content { padding: 16px; }
}
</style>
