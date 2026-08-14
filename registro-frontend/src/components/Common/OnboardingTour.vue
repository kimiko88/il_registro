<template>
  <!-- Onboarding Tour Modal -->
  <q-dialog
    v-model="isVisible"
    persistent
    maximized
    transition-show="fade"
    transition-hide="fade"
    class="onboarding-dialog"
  >
    <div class="onboarding-overlay" @click.self="() => {}">
      <div class="onboarding-card" role="dialog" :aria-label="t('onboarding.welcomeTitle')" aria-modal="true">

        <!-- Welcome screen (step 0) -->
        <transition name="slide-fade" mode="out-in">
          <div v-if="currentStep === 0" key="welcome" class="onboarding-welcome">
            <div class="welcome-illustration">
              <div class="illustration-ring outer-ring"></div>
              <div class="illustration-ring middle-ring"></div>
              <div class="illustration-ring inner-ring"></div>
              <div class="illustration-icon">
                <q-icon :name="roleIcon" size="64px" color="white" />
              </div>
            </div>

            <div class="welcome-content">
              <div class="role-badge">
                <q-chip color="primary" text-color="white" :icon="roleIcon" dense>
                  {{ roleLabel }}
                </q-chip>
              </div>
              <h1 class="welcome-title">{{ t('onboarding.welcomeTitle') }}</h1>
              <p class="welcome-subtitle">{{ t('onboarding.welcomeSubtitle') }}</p>

              <div class="welcome-features">
                <div
                  v-for="(step, idx) in tourSteps"
                  :key="idx"
                  class="feature-item"
                >
                  <q-icon :name="step.icon" color="primary" size="20px" class="feature-icon" />
                  <span class="feature-label">{{ step.title }}</span>
                </div>
              </div>

              <div class="welcome-actions">
                <q-btn
                  unelevated
                  color="primary"
                  size="lg"
                  :label="t('onboarding.startTour')"
                  icon="rocket_launch"
                  class="btn-start"
                  @click="currentStep = 1"
                />
                <q-btn
                  flat
                  color="grey-7"
                  :label="t('onboarding.skipTour')"
                  @click="completeTour"
                  class="btn-skip"
                />
              </div>
            </div>
          </div>

          <!-- Tour steps (step 1-N) -->
          <div v-else key="steps" class="onboarding-step">
            <!-- Progress bar -->
            <div class="step-progress">
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
              </div>
              <span class="progress-label">
                {{ t('onboarding.stepOf', { current: currentStep, total: tourSteps.length }) }}
              </span>
            </div>

            <!-- Step indicator dots -->
            <div class="step-dots">
              <button
                v-for="i in tourSteps.length"
                :key="i"
                class="dot"
                :class="{ active: i === currentStep, done: i < currentStep }"
                @click="currentStep = i"
                :aria-label="`Step ${i}`"
              >
                <q-icon v-if="i < currentStep" name="check" size="10px" color="white" />
              </button>
            </div>

            <!-- Step content -->
            <transition name="slide-content" mode="out-in">
              <div :key="currentStep" class="step-content">
                <!-- Step illustration -->
                <div class="step-illustration">
                  <div class="step-icon-bg" :class="`bg-${activeStep.color}-50`">
                    <q-icon :name="activeStep.icon" :color="activeStep.color" size="48px" />
                  </div>
                </div>

                <div class="step-text">
                  <div class="step-number">
                    <q-chip dense color="primary" text-color="white" size="sm">
                      {{ currentStep }}/{{ tourSteps.length }}
                    </q-chip>
                  </div>
                  <h2 class="step-title">{{ activeStep.title }}</h2>
                  <p class="step-desc">{{ activeStep.desc }}</p>

                  <!-- Tips if any -->
                  <div v-if="activeStep.tip" class="step-tip">
                    <q-icon name="lightbulb" color="amber" size="18px" class="q-mr-xs" />
                    <span>{{ activeStep.tip }}</span>
                  </div>
                </div>
              </div>
            </transition>

            <!-- Navigation -->
            <div class="step-navigation">
              <q-btn
                flat
                round
                icon="arrow_back"
                color="grey-7"
                :disable="currentStep <= 1"
                @click="prevStep"
                :aria-label="t('onboarding.prev')"
              />
              <q-btn
                flat
                color="grey-6"
                :label="t('onboarding.skipTour')"
                @click="completeTour"
                size="sm"
              />
              <q-btn
                v-if="currentStep < tourSteps.length"
                unelevated
                color="primary"
                :label="t('onboarding.next')"
                icon-right="arrow_forward"
                @click="nextStep"
              />
              <q-btn
                v-else
                unelevated
                color="positive"
                :label="t('onboarding.finish')"
                icon="check_circle"
                @click="completeTour"
              />
            </div>
          </div>
        </transition>

      </div>
    </div>
  </q-dialog>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'completed'])

const { t } = useI18n()
const $q = useQuasar()
const authStore = useAuthStore()

const isVisible = ref(false)
const currentStep = ref(0)

// Determine user role
const userRole = computed(() => {
  const role = authStore.userRole || authStore.user?.role || 'student'
  return role.toLowerCase()
})

// Role-specific icon and label
const roleConfig = {
  teacher: { icon: 'school', label: 'roles.teacher', color: 'indigo' },
  student: { icon: 'face', label: 'roles.student', color: 'teal' },
  parent: { icon: 'family_restroom', label: 'roles.parent', color: 'purple' },
  secretary: { icon: 'admin_panel_settings', label: 'roles.secretary', color: 'orange' },
  admin: { icon: 'manage_accounts', label: 'roles.admin', color: 'red' },
  superadmin: { icon: 'manage_accounts', label: 'roles.superadmin', color: 'deep-orange' }
}

const roleIcon = computed(() => roleConfig[userRole.value]?.icon || 'person')
const roleLabel = computed(() => t(roleConfig[userRole.value]?.label || 'roles.user'))

const progressPercent = computed(() => {
  if (!tourSteps.value.length) return 0
  return ((currentStep.value) / tourSteps.value.length) * 100
})

const activeStep = computed(() => tourSteps.value[currentStep.value - 1] || {})

// Define role steps using i18n
const tourSteps = computed(() => {
  const role = userRole.value
  const roleKey = ['teacher', 'student', 'parent', 'secretary', 'admin', 'superadmin'].includes(role)
    ? (role === 'superadmin' ? 'admin' : role)
    : 'student'

  const stepIcons = {
    teacher: [
      { icon: 'dashboard', color: 'indigo' },
      { icon: 'menu_book', color: 'blue' },
      { icon: 'grade', color: 'green' },
      { icon: 'event', color: 'orange' },
      { icon: 'people', color: 'purple' },
      { icon: 'settings', color: 'grey' }
    ],
    student: [
      { icon: 'dashboard', color: 'teal' },
      { icon: 'grade', color: 'green' },
      { icon: 'event_available', color: 'blue' },
      { icon: 'assignment', color: 'orange' },
      { icon: 'description', color: 'purple' },
      { icon: 'calendar_month', color: 'red' }
    ],
    parent: [
      { icon: 'dashboard', color: 'purple' },
      { icon: 'child_care', color: 'pink' },
      { icon: 'grade', color: 'green' },
      { icon: 'campaign', color: 'blue' },
      { icon: 'meeting_room', color: 'orange' },
      { icon: 'payment', color: 'teal' }
    ],
    secretary: [
      { icon: 'dashboard', color: 'orange' },
      { icon: 'groups', color: 'indigo' },
      { icon: 'verified', color: 'green' },
      { icon: 'schedule', color: 'blue' },
      { icon: 'assessment', color: 'purple' },
      { icon: 'manage_accounts', color: 'red' }
    ],
    admin: [
      { icon: 'dashboard', color: 'red' },
      { icon: 'monitor_heart', color: 'orange' },
      { icon: 'manage_accounts', color: 'blue' },
      { icon: 'analytics', color: 'green' },
      { icon: 'security', color: 'purple' },
      { icon: 'settings', color: 'grey' }
    ]
  }

  const icons = stepIcons[roleKey] || stepIcons.student

  return icons.map((ic, idx) => {
    const stepNum = idx + 1
    return {
      icon: ic.icon,
      color: ic.color,
      title: t(`onboarding.${roleKey}.step${stepNum}_title`),
      desc: t(`onboarding.${roleKey}.step${stepNum}_desc`)
    }
  })
})

function nextStep() {
  if (currentStep.value < tourSteps.value.length) {
    currentStep.value++
  }
}

function prevStep() {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

function completeTour() {
  const role = userRole.value || 'user'
  localStorage.setItem(`onboarding_done_${role}`, 'true')
  isVisible.value = false
  emit('update:modelValue', false)
  emit('completed')
  $q.notify({
    type: 'positive',
    icon: 'check_circle',
    message: t('onboarding.tourCompleted'),
    caption: t('onboarding.tourCompletedMsg'),
    timeout: 4000,
    position: 'bottom-right'
  })
}

function startTour() {
  currentStep.value = 0
  isVisible.value = true
  emit('update:modelValue', true)
}

function checkAndStartTour() {
  const role = userRole.value || 'user'
  const done = localStorage.getItem(`onboarding_done_${role}`)
  if (!done) {
    // Small delay for page to render
    setTimeout(() => {
      startTour()
    }, 800)
  }
}

// Watch modelValue from parent
watch(() => props.modelValue, (val) => {
  isVisible.value = val
  if (val) currentStep.value = 0
})

watch(isVisible, (val) => {
  if (!val) emit('update:modelValue', false)
})

onMounted(() => {
  checkAndStartTour()
})

// Expose startTour for external calls (e.g. FAB button)
defineExpose({ startTour })
</script>

<style scoped>
.onboarding-dialog {
  backdrop-filter: blur(8px);
}

.onboarding-overlay {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  min-width: 100vw;
  background: rgba(0, 0, 0, 0.6);
  padding: 16px;
}

.onboarding-card {
  background: white;
  border-radius: 24px;
  width: 100%;
  max-width: 560px;
  box-shadow: 0 32px 80px rgba(0, 0, 0, 0.25);
  overflow: hidden;
  position: relative;
}

/* Dark mode */
:global(.q-dark) .onboarding-card {
  background: #1e1e2e;
  color: #e0e0e0;
}

/* Welcome screen */
.onboarding-welcome {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0;
}

.welcome-illustration {
  width: 100%;
  height: 200px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.illustration-ring {
  position: absolute;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.outer-ring {
  width: 260px;
  height: 260px;
}

.middle-ring {
  width: 180px;
  height: 180px;
}

.inner-ring {
  width: 110px;
  height: 110px;
}

.illustration-icon {
  position: relative;
  z-index: 2;
  background: rgba(255,255,255,0.15);
  border-radius: 50%;
  width: 96px;
  height: 96px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(4px);
  border: 2px solid rgba(255,255,255,0.3);
}

.welcome-content {
  padding: 28px 32px 32px;
  text-align: center;
  width: 100%;
}

.role-badge {
  margin-bottom: 12px;
}

.welcome-title {
  font-size: 1.6rem;
  font-weight: 800;
  color: #1a1a2e;
  margin: 0 0 10px;
  line-height: 1.2;
}

:global(.q-dark) .welcome-title {
  color: #e0e0e0;
}

.welcome-subtitle {
  font-size: 0.97rem;
  color: #6b7280;
  margin: 0 0 24px;
  line-height: 1.6;
}

.welcome-features {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-bottom: 28px;
  text-align: left;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: #f8f9ff;
  border-radius: 10px;
  font-size: 0.85rem;
  font-weight: 500;
  color: #374151;
}

:global(.q-dark) .feature-item {
  background: #2d2d3f;
  color: #c0c0d0;
}

.feature-icon {
  flex-shrink: 0;
}

.welcome-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
}

.btn-start {
  width: 100%;
  border-radius: 12px;
  font-weight: 700;
  font-size: 1rem;
  padding: 14px 24px;
}

.btn-skip {
  font-size: 0.875rem;
}

/* Step content */
.onboarding-step {
  padding: 28px 32px 32px;
}

.step-progress {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.progress-bar {
  flex: 1;
  height: 6px;
  background: #e5e7eb;
  border-radius: 999px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border-radius: 999px;
  transition: width 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.progress-label {
  font-size: 0.8rem;
  color: #9ca3af;
  white-space: nowrap;
  font-weight: 500;
}

.step-dots {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 24px;
}

.dot {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid #e5e7eb;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  font-size: 0;
  padding: 0;
}

:global(.q-dark) .dot {
  background: #2d2d3f;
  border-color: #4a4a5e;
}

.dot.active {
  border-color: #667eea;
  background: #667eea;
  transform: scale(1.2);
}

.dot.done {
  border-color: #667eea;
  background: #667eea;
}

.step-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.step-illustration {
  margin-bottom: 24px;
}

.step-icon-bg {
  width: 100px;
  height: 100px;
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.step-number {
  margin-bottom: 12px;
}

.step-title {
  font-size: 1.4rem;
  font-weight: 800;
  color: #1a1a2e;
  margin: 0 0 12px;
  line-height: 1.25;
}

:global(.q-dark) .step-title {
  color: #e0e0e0;
}

.step-desc {
  font-size: 0.97rem;
  color: #6b7280;
  line-height: 1.7;
  margin: 0 0 16px;
  max-width: 420px;
}

.step-tip {
  display: inline-flex;
  align-items: flex-start;
  gap: 6px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 10px;
  padding: 10px 14px;
  font-size: 0.875rem;
  color: #92400e;
  text-align: left;
}

:global(.q-dark) .step-tip {
  background: #2d2a1e;
  border-color: #92400e;
  color: #fde68a;
}

.step-navigation {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 28px;
  padding-top: 20px;
  border-top: 1px solid #f3f4f6;
}

:global(.q-dark) .step-navigation {
  border-color: #3a3a4e;
}

/* Transitions */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease;
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.slide-content-enter-active,
.slide-content-leave-active {
  transition: all 0.25s ease;
}

.slide-content-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.slide-content-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* RTL support */
:global([dir='rtl']) .step-navigation {
  flex-direction: row-reverse;
}

:global([dir='rtl']) .welcome-features {
  text-align: right;
}
</style>
