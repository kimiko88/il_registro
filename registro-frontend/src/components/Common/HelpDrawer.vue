<template>
  <!-- Help Drawer: slides in from the right -->
  <q-drawer
    v-model="isOpen"
    side="right"
    :width="620"
    elevated
    class="help-drawer"
    :aria-label="t('help.title')"
  >
    <!-- Drawer header -->
    <div class="help-header">
      <div class="help-header-left">
        <q-icon name="help_outline" size="28px" color="primary" class="q-mr-md" />
        <div>
          <div class="help-title">{{ t('help.title') }}</div>
          <div class="help-subtitle-text">{{ t('help.subtitle') }}</div>
        </div>
      </div>
      <q-btn flat round dense icon="close" @click="close" :aria-label="t('common.close')" />
    </div>

    <!-- Search bar -->
    <div class="help-search-bar">
      <q-input
        v-model="searchQuery"
        :placeholder="t('help.searchPlaceholder')"
        dense
        outlined
        clearable
        class="help-search-input"
        bg-color="white"
        autofocus
      >
        <template #prepend>
          <q-icon name="search" color="grey-5" />
        </template>
      </q-input>
    </div>

    <!-- Restart tour button -->
    <div class="help-tour-row">
      <q-btn
        unelevated
        color="primary"
        icon="rocket_launch"
        :label="t('help.restartTour')"
        class="full-width restart-btn"
        @click="restartTour"
      />
    </div>

    <!-- Category filters -->
    <div class="help-categories" v-if="!searchQuery">
      <button
        class="category-chip"
        :class="{ active: activeCategory === '' }"
        @click="activeCategory = ''"
      >{{ t('help.allTopics') }}</button>
      <button
        v-for="cat in categories"
        :key="cat.key"
        class="category-chip"
        :class="{ active: activeCategory === cat.key }"
        @click="activeCategory = cat.key"
      >{{ cat.label }}</button>
    </div>

    <!-- FAQ Articles -->
    <div class="help-articles">
      <div v-if="filteredArticles.length === 0" class="help-empty">
        <q-icon name="search_off" size="48px" color="grey-4" />
        <p>{{ t('help.noResults') }} "<strong>{{ searchQuery }}</strong>"</p>
      </div>

      <q-expansion-item
        v-for="article in filteredArticles"
        :key="article.key"
        :label="article.question"
        class="help-article"
        expand-separator
        header-class="help-article-header"
        :icon="article.categoryIcon"
      >
        <div class="help-article-content">
          <q-icon name="arrow_right" color="primary" size="18px" class="q-mr-xs" style="flex-shrink:0;margin-top:2px" />
          <p>{{ article.answer }}</p>
        </div>
      </q-expansion-item>
    </div>

    <!-- Footer -->
    <div class="help-footer">
      <q-btn
        flat
        color="primary"
        icon="support_agent"
        :label="t('help.contactSupport')"
        @click="goToSupport"
      />
    </div>
  </q-drawer>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const emit = defineEmits(['restart-tour'])
const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()

const isOpen = ref(false)
const searchQuery = ref('')
const activeCategory = ref('')

const userRole = computed(() => {
  const role = authStore.userRole || authStore.user?.role || 'student'
  const r = role.toLowerCase()
  if (r === 'superadmin') return 'admin'
  return ['teacher', 'student', 'parent', 'secretary', 'admin'].includes(r) ? r : 'student'
})

// Category definitions per role
const roleCategoryKeys = {
  teacher: ['cat_register', 'cat_grades', 'cat_attendance', 'cat_agenda', 'cat_settings'],
  student: ['cat_grades', 'cat_attendance', 'cat_homework', 'cat_documents', 'cat_settings'],
  parent: ['cat_monitoring', 'cat_communications', 'cat_meetings', 'cat_documents', 'cat_settings'],
  secretary: ['cat_students', 'cat_classes', 'cat_documents', 'cat_timetable', 'cat_reports'],
  admin: ['cat_monitoring', 'cat_users', 'cat_schools', 'cat_security', 'cat_analytics']
}

const categoryIcons = {
  cat_register: 'menu_book',
  cat_grades: 'grade',
  cat_attendance: 'event_available',
  cat_agenda: 'event',
  cat_settings: 'settings',
  cat_homework: 'assignment',
  cat_documents: 'description',
  cat_monitoring: 'monitor_heart',
  cat_communications: 'campaign',
  cat_meetings: 'people',
  cat_students: 'group',
  cat_classes: 'class',
  cat_timetable: 'schedule',
  cat_reports: 'assessment',
  cat_users: 'manage_accounts',
  cat_schools: 'school',
  cat_security: 'security',
  cat_analytics: 'analytics'
}

const categories = computed(() => {
  const role = userRole.value
  const keys = roleCategoryKeys[role] || []
  return keys.map(key => ({
    key,
    label: t(`help.${role}.${key}`),
    icon: categoryIcons[key] || 'help'
  }))
})

const articles = computed(() => {
  const role = userRole.value
  const result = []
  for (let i = 1; i <= 5; i++) {
    const q = t(`help.${role}.q${i}`)
    const a = t(`help.${role}.a${i}`)
    if (q && a && !q.startsWith('help.')) {
      // Find matching category (map i to category)
      const catKeys = roleCategoryKeys[role] || []
      const catIdx = Math.floor((i - 1) / 1) % catKeys.length
      const catKey = catKeys[catIdx] || ''
      result.push({
        key: `${role}-q${i}`,
        question: q,
        answer: a,
        categoryKey: catKey,
        categoryIcon: categoryIcons[catKey] || 'help_outline'
      })
    }
  }
  return result
})

const filteredArticles = computed(() => {
  let list = articles.value
  if (activeCategory.value) {
    list = list.filter(a => a.categoryKey === activeCategory.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(a =>
      a.question.toLowerCase().includes(q) ||
      a.answer.toLowerCase().includes(q)
    )
  }
  return list
})

function open() {
  isOpen.value = true
}

function close() {
  isOpen.value = false
}

function restartTour() {
  close()
  emit('restart-tour')
}

function goToSupport() {
  close()
  router.push('/support')
}

defineExpose({ open, close })
</script>

<style scoped>
.help-drawer {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.help-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 20px 16px;
  border-bottom: 1px solid #f3f4f6;
  background: linear-gradient(135deg, #667eea10, #764ba210);
}

:global(.q-dark) .help-header {
  border-color: #3a3a4e;
  background: linear-gradient(135deg, #667eea15, #764ba215);
}

.help-header-left {
  display: flex;
  align-items: center;
}

.help-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: #1a1a2e;
}

:global(.q-dark) .help-title {
  color: #e0e0e0;
}

.help-subtitle-text {
  font-size: 0.78rem;
  color: #9ca3af;
  margin-top: 1px;
}

.help-search-bar {
  padding: 14px 16px 10px;
}

.help-search-input {
  border-radius: 10px;
}

.help-tour-row {
  padding: 0 16px 12px;
}

.restart-btn {
  border-radius: 10px;
  font-weight: 700;
}

.help-categories {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 0 16px 12px;
}

.category-chip {
  padding: 6px 14px;
  border-radius: 999px;
  border: 1.5px solid #e5e7eb;
  background: white;
  font-size: 0.8rem;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

:global(.q-dark) .category-chip {
  background: #2d2d3f;
  border-color: #4a4a5e;
  color: #9ca3af;
}

.category-chip:hover {
  border-color: #667eea;
  color: #667eea;
}

.category-chip.active {
  background: #667eea;
  border-color: #667eea;
  color: white;
}

.help-articles {
  flex: 1;
  overflow-y: auto;
  padding: 0 8px;
}

.help-article {
  margin-bottom: 4px;
  border-radius: 10px;
  overflow: hidden;
}

.help-article-content {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  padding: 12px 16px 16px;
  font-size: 0.9rem;
  color: #6b7280;
  line-height: 1.7;
  background: #f8f9ff;
}

:global(.q-dark) .help-article-content {
  background: #252535;
  color: #9ca3af;
}

.help-article-content p {
  margin: 0;
}

.help-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 16px;
  color: #9ca3af;
  text-align: center;
  gap: 12px;
}

.help-footer {
  padding: 12px 16px;
  border-top: 1px solid #f3f4f6;
  display: flex;
  justify-content: center;
}

:global(.q-dark) .help-footer {
  border-color: #3a3a4e;
}
</style>
