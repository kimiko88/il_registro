<template>
  <q-dialog
    v-model="isVisible"
    persistent
    maximized
    transition-show="fade"
    transition-hide="fade"
    class="onboarding-dialog"
  >
    <div class="onboarding-overlay" @keydown="handleKeydown" tabindex="-1" ref="overlayRef">

      <!-- WELCOME SCREEN -->
      <transition name="slide-up" mode="out-in">
        <div v-if="currentStep === 0" key="welcome" class="onboarding-welcome-card">

          <!-- Left: gradient hero -->
          <div class="welcome-hero" :class="`hero-${roleColor}`">
            <div class="hero-rings">
              <div class="ring r1"></div>
              <div class="ring r2"></div>
              <div class="ring r3"></div>
            </div>
            <div class="hero-avatar">
              <q-icon :name="roleIcon" size="72px" color="white" />
            </div>
            <div class="hero-label">
              <q-chip dense color="white" :text-color="rolePrimary" class="text-weight-bold">
                {{ roleLabel }}
              </q-chip>
            </div>
            <!-- Floating mini feature icons -->
            <div class="hero-floats">
              <div v-for="(step, i) in tourSteps.slice(0,4)" :key="i"
                class="float-icon"
                :style="floatIconStyle(i)"
              >
                <q-icon :name="step.icon" size="18px" color="white" />
              </div>
            </div>
          </div>

          <!-- Right: content -->
          <div class="welcome-body">
            <div class="welcome-top">
              <div class="welcome-eyebrow">{{ te('common.appName') ? t('common.appName') : 'Registro Elettronico' }}</div>
              <h1 class="welcome-title">{{ t('onboarding.welcomeTitle') }}</h1>
              <p class="welcome-subtitle">{{ t('onboarding.welcomeSubtitle') }}</p>
            </div>

            <!-- Feature grid preview -->
            <div class="feature-grid">
              <div
                v-for="(step, i) in tourSteps"
                :key="i"
                class="feature-chip"
                @click="jumpToStep(i + 1)"
              >
                <div class="feature-chip-icon" :class="`icon-${step.color}`">
                  <q-icon :name="step.icon" size="16px" />
                </div>
                <span>{{ step.title }}</span>
              </div>
            </div>

            <div class="welcome-actions">
              <q-btn
                unelevated
                color="primary"
                size="lg"
                :label="t('onboarding.startTour')"
                icon="rocket_launch"
                class="btn-cta"
                @click="goToStep(1)"
              />
              <q-btn
                flat
                color="grey-6"
                :label="t('onboarding.skipTour')"
                @click="completeTour"
                size="sm"
              />
            </div>

            <div class="welcome-hint">
              <q-icon name="keyboard" size="16px" color="grey-4" class="q-mr-xs" />
              <span>{{ te('onboarding.keyboardHint') ? t('onboarding.keyboardHint') : 'Usa ← → per navigare' }}</span>
            </div>
          </div>
        </div>

        <!-- COMPLETION SCREEN -->
        <div v-else-if="currentStep === COMPLETION_STEP" key="completion" class="onboarding-completion-card">
          <div class="completion-animation">
            <div class="confetti-ring"></div>
            <div class="completion-icon">🎉</div>
          </div>
          <h2 class="completion-title">{{ t('onboardingExtra.completionTitle') }}</h2>
          <p class="completion-desc">{{ t('onboardingExtra.completionDesc') }}</p>
          <div class="completion-actions">
            <q-btn
              unelevated
              color="primary"
              size="lg"
              :label="t('onboarding.finish')"
              icon="check_circle"
              class="btn-cta"
              @click="completeTour"
            />
            <q-btn
              flat
              color="primary"
              :label="t('onboardingExtra.openGuide')"
              icon="menu_book"
              @click="openGuideAndComplete"
            />
          </div>
        </div>

        <!-- STEP SCREENS -->
        <div v-else key="steps" class="onboarding-step-card">

          <!-- Left panel: visual -->
          <div class="step-panel-left" :class="`panel-${activeStep.color}`">
            <!-- Progress sidebar -->
            <div class="sidebar-steps">
              <button
                v-for="(s, i) in tourSteps"
                :key="i"
                class="sidebar-dot"
                :class="{
                  active: (i + 1) === currentStep,
                  done: (i + 1) < currentStep
                }"
                @click="goToStep(i + 1)"
                :title="s.title"
              >
                <q-icon v-if="(i+1) < currentStep" name="check" size="11px" color="white" />
                <span v-else class="dot-num">{{ i + 1 }}</span>
              </button>
            </div>

            <!-- Step visual -->
            <div class="step-visual">
              <transition name="icon-pop" mode="out-in">
                <div :key="currentStep" class="step-icon-wrapper">
                  <div class="step-icon-glow"></div>
                  <q-icon :name="activeStep.icon" size="72px" color="white" />
                </div>
              </transition>
              <div class="step-visual-label">{{ activeStep.title }}</div>
            </div>

            <!-- Progress bar at bottom of left panel -->
            <div class="sidebar-progress">
              <div class="sidebar-progress-fill" :style="{ height: progressPercent + '%' }"></div>
            </div>
          </div>

          <!-- Right panel: content -->
          <div class="step-panel-right">
            <!-- Top bar -->
            <div class="step-topbar">
              <div class="step-counter">
                <span class="step-num-current">{{ currentStep }}</span>
                <span class="step-num-sep">/</span>
                <span class="step-num-total">{{ tourSteps.length }}</span>
              </div>
              <q-btn
                flat
                round
                dense
                icon="close"
                color="grey-5"
                size="sm"
                @click="completeTour"
                :aria-label="t('onboarding.skipTour')"
              />
            </div>

            <!-- Step content with transition -->
            <transition name="step-slide" mode="out-in">
              <div :key="currentStep" class="step-content-area">
                <div class="step-badge" :class="`badge-${activeStep.color}`">
                  <q-icon :name="activeStep.icon" size="14px" />
                  {{ activeStep.category || '' }}
                </div>

                <h2 class="step-title">{{ activeStep.title }}</h2>
                <p class="step-desc">{{ activeStep.desc }}</p>

                <!-- Bullet list features -->
                <ul v-if="activeStep.bullets && activeStep.bullets.length" class="step-bullets">
                  <li v-for="(b, bi) in activeStep.bullets" :key="bi" class="step-bullet">
                    <span class="bullet-dot" :class="`dot-${activeStep.color}`"></span>
                    {{ b }}
                  </li>
                </ul>

                <!-- Tip if any -->
                <div v-if="activeStep.tip" class="step-tip">
                  <q-icon name="lightbulb" color="amber-7" size="16px" />
                  <span>{{ activeStep.tip }}</span>
                </div>
              </div>
            </transition>

            <!-- Navigation -->
            <div class="step-nav">
              <q-btn
                flat
                :label="t('onboarding.prev')"
                icon="arrow_back"
                color="grey-6"
                :disable="currentStep <= 1"
                @click="prevStep"
              />

              <!-- Progress dots (compact) -->
              <div class="nav-dots">
                <span
                  v-for="i in tourSteps.length"
                  :key="i"
                  class="nav-dot"
                  :class="{
                    'nav-dot-active': i === currentStep,
                    'nav-dot-done': i < currentStep
                  }"
                  @click="goToStep(i)"
                ></span>
              </div>

              <q-btn
                v-if="currentStep < tourSteps.length"
                unelevated
                :label="t('onboarding.next')"
                icon-right="arrow_forward"
                color="primary"
                @click="nextStep"
              />
              <q-btn
                v-else
                unelevated
                :label="t('onboarding.finish')"
                icon="check_circle"
                color="positive"
                @click="goToStep(COMPLETION_STEP)"
              />
            </div>
          </div>

        </div>
      </transition>

    </div>
  </q-dialog>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  modelValue: { type: Boolean, default: false }
})
const emit = defineEmits(['update:modelValue', 'completed', 'open-guide'])

const { t, tm, te } = useI18n()
const $q = useQuasar()
const authStore = useAuthStore()

const isVisible = ref(false)
const currentStep = ref(0)
const overlayRef = ref(null)

const COMPLETION_STEP = 999

// ── Role helpers ─────────────────────────────────────────────────────────────
const userRole = computed(() => {
  const role = authStore.userRole || authStore.user?.role || 'student'
  const r = role.toLowerCase()
  if (r === 'superadmin') return 'admin'
  if (['teacher', 'student', 'parent', 'secretary', 'admin',
       'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico', 'dsga'].includes(r)) {
    return r
  }
  return 'student'
})

const roleMeta = {
  teacher:                   { icon: 'school',              color: 'indigo',      primary: 'indigo',        label: 'roles.teacher' },
  student:                   { icon: 'face',                color: 'teal',        primary: 'teal',          label: 'roles.student' },
  parent:                    { icon: 'family_restroom',      color: 'purple',      primary: 'purple',        label: 'roles.parent' },
  secretary:                 { icon: 'admin_panel_settings', color: 'orange',      primary: 'orange',        label: 'roles.secretary' },
  admin:                     { icon: 'manage_accounts',      color: 'red',         primary: 'red',           label: 'roles.admin' },
  assistente_amministrativo: { icon: 'manage_accounts',      color: 'cyan',        primary: 'cyan-8',        label: 'roles.assistente_amministrativo' },
  collaboratore_ds:          { icon: 'co_present',          color: 'amber',       primary: 'amber-9',       label: 'roles.collaboratore_ds' },
  collaboratore_scolastico:  { icon: 'door_front',          color: 'teal',        primary: 'teal-8',        label: 'roles.collaboratore_scolastico' },
  dsga:                      { icon: 'account_balance',     color: 'deep-orange', primary: 'deep-orange-8', label: 'roles.dsga' }
}

const roleIcon    = computed(() => roleMeta[userRole.value]?.icon    || 'person')
const roleColor   = computed(() => roleMeta[userRole.value]?.color   || 'indigo')
const rolePrimary = computed(() => roleMeta[userRole.value]?.primary || 'indigo')
const roleLabel   = computed(() => t(roleMeta[userRole.value]?.label || 'roles.user'))

// ── Steps definition ─────────────────────────────────────────────────────────
const STEP_DEFS = {
  teacher: [
    { icon:'dashboard',       color:'indigo'  },
    { icon:'menu_book',       color:'blue'    },
    { icon:'grade',           color:'green'   },
    { icon:'event',           color:'orange'  },
    { icon:'people',          color:'purple'  },
    { icon:'settings',        color:'grey'    },
    { icon:'draw',            color:'cyan'    },
    { icon:'workspace_premium', color:'amber' }
  ],
  student: [
    { icon:'dashboard',       color:'teal'    },
    { icon:'grade',           color:'green'   },
    { icon:'event_available', color:'blue'    },
    { icon:'assignment',      color:'orange'  },
    { icon:'description',     color:'purple'  },
    { icon:'calendar_month',  color:'red'     },
    { icon:'work',            color:'cyan'    },
    { icon:'campaign',        color:'pink'    }
  ],
  parent: [
    { icon:'dashboard',       color:'purple'  },
    { icon:'child_care',      color:'pink'    },
    { icon:'grade',           color:'green'   },
    { icon:'campaign',        color:'blue'    },
    { icon:'meeting_room',    color:'orange'  },
    { icon:'payment',         color:'teal'    },
    { icon:'verified_user',   color:'indigo'  },
    { icon:'analytics',       color:'red'     }
  ],
  secretary: [
    { icon:'dashboard',       color:'orange'  },
    { icon:'groups',          color:'indigo'  },
    { icon:'verified',        color:'green'   },
    { icon:'schedule',        color:'blue'    },
    { icon:'assessment',      color:'purple'  },
    { icon:'manage_accounts', color:'red'     },
    { icon:'badge',           color:'cyan'    },
    { icon:'newspaper',       color:'teal'    }
  ],
  admin: [
    { icon:'dashboard',       color:'red'     },
    { icon:'monitor_heart',   color:'orange'  },
    { icon:'manage_accounts', color:'blue'    },
    { icon:'analytics',       color:'green'   },
    { icon:'security',        color:'purple'  },
    { icon:'settings',        color:'grey'    },
    { icon:'corporate_fare',  color:'indigo'  },
    { icon:'api',             color:'cyan'    }
  ],
  assistente_amministrativo: [
    { icon: 'dashboard',        color: 'cyan' },
    { icon: 'calendar_month',   color: 'blue' },
    { icon: 'forward_to_inbox', color: 'purple' },
    { icon: 'cloud_sync',       color: 'indigo' },
    { icon: 'gavel',            color: 'teal' }
  ],
  collaboratore_ds: [
    { icon: 'dashboard',        color: 'amber' },
    { icon: 'bolt',             color: 'red' },
    { icon: 'swap_horiz',       color: 'blue' },
    { icon: 'co_present',       color: 'orange' },
    { icon: 'gavel',            color: 'purple' }
  ],
  collaboratore_scolastico: [
    { icon: 'dashboard',        color: 'teal' },
    { icon: 'door_front',       color: 'blue' },
    { icon: 'logout',           color: 'orange' },
    { icon: 'build',            color: 'red' },
    { icon: 'badge',            color: 'green' }
  ],
  dsga: [
    { icon: 'dashboard',        color: 'deep-orange' },
    { icon: 'assessment',       color: 'blue' },
    { icon: 'verified',         color: 'green' },
    { icon: 'cloud_sync',       color: 'indigo' },
    { icon: 'gavel',            color: 'purple' }
  ]
}

const tourSteps = computed(() => {
  const role = userRole.value
  const defs = STEP_DEFS[role] || STEP_DEFS.student
  return defs.map((d, i) => {
    const n = i + 1
    const titleKey = n <= 6
      ? `onboarding.${role}.step${n}_title`
      : `onboardingExtra.${role}.step${n}_title`
    const descKey = n <= 6
      ? `onboarding.${role}.step${n}_desc`
      : `onboardingExtra.${role}.step${n}_desc`
    const bulletsKey = `onboardingExtra.${role}.step${n}_bullets`

    const rawBullets = tm(bulletsKey)
    const bullets = Array.isArray(rawBullets) ? rawBullets : []

    return {
      icon:    d.icon,
      color:   d.color,
      title:   t(titleKey),
      desc:    t(descKey),
      bullets
    }
  })
})

const activeStep     = computed(() => tourSteps.value[currentStep.value - 1] || {})
const progressPercent = computed(() => {
  const total = tourSteps.value.length
  if (!total) return 0
  return Math.min(100, (currentStep.value / total) * 100)
})

// ── Navigation ────────────────────────────────────────────────────────────────
function goToStep(n) {
  currentStep.value = n
  focusOverlay()
}
function nextStep() {
  if (currentStep.value < tourSteps.value.length) currentStep.value++
  else goToStep(COMPLETION_STEP)
  focusOverlay()
}
function prevStep() {
  if (currentStep.value > 1) currentStep.value--
  focusOverlay()
}
function jumpToStep(n) { goToStep(n) }

function handleKeydown(e) {
  if (e.key === 'ArrowRight' || e.key === 'ArrowDown')  { e.preventDefault(); nextStep() }
  if (e.key === 'ArrowLeft'  || e.key === 'ArrowUp')    { e.preventDefault(); prevStep() }
  if (e.key === 'Escape')                                { completeTour() }
}

async function focusOverlay() {
  await nextTick()
  overlayRef.value?.focus()
}

// ── Completion ────────────────────────────────────────────────────────────────
function completeTour() {
  const role = userRole.value
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

function openGuideAndComplete() {
  completeTour()
  emit('open-guide')
}

// ── Public API ────────────────────────────────────────────────────────────────
function startTour() {
  currentStep.value = 0
  isVisible.value = true
  emit('update:modelValue', true)
}

function checkAndStartTour() {
  const done = localStorage.getItem(`onboarding_done_${userRole.value}`)
  if (!done) {
    setTimeout(startTour, 800)
  }
}

watch(() => props.modelValue, val => {
  isVisible.value = val
  if (val) currentStep.value = 0
})
watch(isVisible, val => {
  if (!val) emit('update:modelValue', false)
})

onMounted(checkAndStartTour)
defineExpose({ startTour })

// ── Decorative helpers ────────────────────────────────────────────────────────
function floatIconStyle(i) {
  const positions = [
    { top: '18%', left: '15%' },
    { top: '20%', right: '12%' },
    { bottom: '22%', left: '10%' },
    { bottom: '18%', right: '15%' }
  ]
  return positions[i] || {}
}
</script>

<style scoped>
/* ── Dialog overlay ─────────────────────────────────────────────── */
.onboarding-dialog { backdrop-filter: blur(10px); }

.onboarding-overlay {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  min-width: 100vw;
  background: rgba(0,0,0,0.65);
  padding: 16px;
  outline: none;
}

/* ── Shared card base ────────────────────────────────────────────── */
.onboarding-welcome-card,
.onboarding-step-card,
.onboarding-completion-card {
  background: #fff;
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 40px 100px rgba(0,0,0,0.3);
  width: 100%;
  max-width: 800px;
  min-height: 520px;
}
:global(.q-dark) .onboarding-welcome-card,
:global(.q-dark) .onboarding-step-card,
:global(.q-dark) .onboarding-completion-card {
  background: #1e1e2e;
  color: #e0e0e0;
}

/* ── WELCOME CARD ────────────────────────────────────────────────── */
.onboarding-welcome-card {
  display: flex;
  min-height: 540px;
}

/* Hero panel */
.welcome-hero {
  width: 260px;
  flex-shrink: 0;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 32px 16px;
}

.hero-indigo  { background: linear-gradient(160deg, #667eea, #4338ca); }
.hero-teal    { background: linear-gradient(160deg, #0d9488, #0f766e); }
.hero-purple  { background: linear-gradient(160deg, #a855f7, #7c3aed); }
.hero-orange  { background: linear-gradient(160deg, #f97316, #ea580c); }
.hero-red     { background: linear-gradient(160deg, #ef4444, #dc2626); }
.hero-cyan    { background: linear-gradient(160deg, #06b6d4, #0891b2); }
.hero-amber   { background: linear-gradient(160deg, #f59e0b, #d97706); }
.hero-deep-orange { background: linear-gradient(160deg, #f97316, #c2410c); }

.hero-rings { position: absolute; inset: 0; }
.ring {
  position: absolute;
  border-radius: 50%;
  border: 1.5px solid rgba(255,255,255,0.15);
  top: 50%; left: 50%;
  transform: translate(-50%, -50%);
}
.r1 { width: 320px; height: 320px; }
.r2 { width: 220px; height: 220px; }
.r3 { width: 130px; height: 130px; }

.hero-avatar {
  position: relative;
  z-index: 2;
  background: rgba(255,255,255,0.15);
  border-radius: 50%;
  width: 110px; height: 110px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid rgba(255,255,255,0.25);
  backdrop-filter: blur(4px);
  margin-bottom: 16px;
}
.hero-label { position: relative; z-index: 2; }

.hero-floats { position: absolute; inset: 0; z-index: 1; pointer-events: none; }
.float-icon {
  position: absolute;
  background: rgba(255,255,255,0.15);
  border-radius: 10px;
  width: 36px; height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(3px);
  border: 1px solid rgba(255,255,255,0.2);
  animation: float 3s ease-in-out infinite;
}
.float-icon:nth-child(2) { animation-delay: 0.5s; }
.float-icon:nth-child(3) { animation-delay: 1s; }
.float-icon:nth-child(4) { animation-delay: 1.5s; }

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50%       { transform: translateY(-6px); }
}

/* Welcome body */
.welcome-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 36px 32px 28px;
  overflow-y: auto;
}
.welcome-top { margin-bottom: 20px; }
.welcome-eyebrow {
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 2px;
  text-transform: uppercase;
  color: #9ca3af;
  margin-bottom: 6px;
}
.welcome-title {
  font-size: 1.65rem;
  font-weight: 900;
  color: #1a1a2e;
  margin: 0 0 8px;
  line-height: 1.2;
}
:global(.q-dark) .welcome-title { color: #e0e0e0; }

.welcome-subtitle {
  font-size: 0.93rem;
  color: #6b7280;
  line-height: 1.65;
  margin: 0;
}

.feature-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 24px;
  flex: 1;
}
.feature-chip {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 9px 12px;
  background: #f8f9ff;
  border-radius: 10px;
  font-size: 0.82rem;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}
:global(.q-dark) .feature-chip {
  background: #252535;
  color: #c0c0d0;
}
.feature-chip:hover {
  border-color: #667eea44;
  background: #f0f1ff;
  transform: translateY(-1px);
}
.feature-chip-icon {
  width: 28px; height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 14px;
}
.icon-indigo  { background: #eef2ff; color: #4338ca; }
.icon-teal    { background: #f0fdfa; color: #0d9488; }
.icon-purple  { background: #faf5ff; color: #7c3aed; }
.icon-orange  { background: #fff7ed; color: #ea580c; }
.icon-red     { background: #fef2f2; color: #dc2626; }
.icon-green   { background: #f0fdf4; color: #16a34a; }
.icon-blue    { background: #eff6ff; color: #2563eb; }
.icon-grey    { background: #f9fafb; color: #6b7280; }
.icon-cyan    { background: #ecfeff; color: #0891b2; }
.icon-pink    { background: #fdf2f8; color: #db2777; }
.icon-amber   { background: #fffbeb; color: #d97706; }

.welcome-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: stretch;
}
.btn-cta {
  border-radius: 12px;
  font-weight: 800;
  font-size: 1rem;
  padding: 14px;
}
.welcome-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 12px;
  font-size: 0.75rem;
  color: #d1d5db;
}

/* ── STEP CARD ───────────────────────────────────────────────────── */
.onboarding-step-card {
  display: flex;
  min-height: 520px;
}

/* Left visual panel */
.step-panel-left {
  width: 220px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 28px 16px;
}

.panel-indigo  { background: linear-gradient(160deg, #667eea, #4338ca); }
.panel-teal    { background: linear-gradient(160deg, #0d9488, #0f766e); }
.panel-purple  { background: linear-gradient(160deg, #a855f7, #7c3aed); }
.panel-orange  { background: linear-gradient(160deg, #f97316, #ea580c); }
.panel-red     { background: linear-gradient(160deg, #ef4444, #dc2626); }
.panel-green   { background: linear-gradient(160deg, #22c55e, #16a34a); }
.panel-blue    { background: linear-gradient(160deg, #3b82f6, #2563eb); }
.panel-grey    { background: linear-gradient(160deg, #9ca3af, #6b7280); }
.panel-cyan    { background: linear-gradient(160deg, #22d3ee, #0891b2); }
.panel-pink    { background: linear-gradient(160deg, #ec4899, #db2777); }
.panel-amber   { background: linear-gradient(160deg, #f59e0b, #d97706); }

/* Sidebar step dots */
.sidebar-steps {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.sidebar-dot {
  width: 24px; height: 24px;
  border-radius: 50%;
  border: 2px solid rgba(255,255,255,0.4);
  background: rgba(255,255,255,0.1);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.25s;
  font-size: 0;
  padding: 0;
  color: white;
}
.sidebar-dot.active {
  background: white;
  border-color: white;
  transform: scale(1.25);
}
.sidebar-dot.done {
  background: rgba(255,255,255,0.5);
  border-color: white;
}
.dot-num {
  font-size: 9px;
  font-weight: 700;
  color: rgba(255,255,255,0.7);
}
.sidebar-dot.active .dot-num { color: #4338ca; }

/* Step visual */
.step-visual {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  text-align: center;
}
.step-icon-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}
.step-icon-glow {
  position: absolute;
  width: 100px; height: 100px;
  background: rgba(255,255,255,0.15);
  border-radius: 50%;
  filter: blur(12px);
}
.step-visual-label {
  font-size: 0.78rem;
  font-weight: 700;
  color: rgba(255,255,255,0.75);
  text-align: center;
  padding: 0 8px;
  max-width: 160px;
}

/* Sidebar vertical progress */
.sidebar-progress {
  position: absolute;
  right: 8px;
  top: 16px;
  bottom: 16px;
  width: 3px;
  background: rgba(255,255,255,0.15);
  border-radius: 999px;
  overflow: hidden;
}
.sidebar-progress-fill {
  width: 100%;
  background: rgba(255,255,255,0.7);
  border-radius: 999px;
  transition: height 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  position: absolute;
  bottom: 0;
}

/* Right content panel */
.step-panel-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 28px 32px 24px;
  overflow-y: auto;
}

.step-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}
.step-counter {
  display: flex;
  align-items: baseline;
  gap: 3px;
}
.step-num-current {
  font-size: 1.5rem;
  font-weight: 900;
  color: #1a1a2e;
  line-height: 1;
}
:global(.q-dark) .step-num-current { color: #e0e0e0; }
.step-num-sep {
  font-size: 1rem;
  color: #d1d5db;
}
.step-num-total {
  font-size: 1rem;
  color: #9ca3af;
}

.step-content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.step-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 4px 10px;
  border-radius: 999px;
  margin-bottom: 12px;
  width: fit-content;
}
.badge-indigo  { background: #eef2ff; color: #4338ca; }
.badge-teal    { background: #f0fdfa; color: #0d9488; }
.badge-purple  { background: #faf5ff; color: #7c3aed; }
.badge-orange  { background: #fff7ed; color: #ea580c; }
.badge-red     { background: #fef2f2; color: #dc2626; }
.badge-green   { background: #f0fdf4; color: #16a34a; }
.badge-blue    { background: #eff6ff; color: #2563eb; }
.badge-grey    { background: #f9fafb; color: #4b5563; }
.badge-cyan    { background: #ecfeff; color: #0891b2; }
.badge-pink    { background: #fdf2f8; color: #db2777; }
.badge-amber   { background: #fffbeb; color: #d97706; }

.step-title {
  font-size: 1.45rem;
  font-weight: 900;
  color: #1a1a2e;
  margin: 0 0 10px;
  line-height: 1.25;
}
:global(.q-dark) .step-title { color: #e0e0e0; }

.step-desc {
  font-size: 0.93rem;
  color: #6b7280;
  line-height: 1.7;
  margin: 0 0 16px;
}

/* Bullet list */
.step-bullets {
  list-style: none;
  padding: 0;
  margin: 0 0 16px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}
.step-bullet {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.84rem;
  font-weight: 500;
  color: #374151;
}
:global(.q-dark) .step-bullet { color: #c0c0d0; }

.bullet-dot {
  width: 8px; height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot-indigo  { background: #667eea; }
.dot-teal    { background: #0d9488; }
.dot-purple  { background: #a855f7; }
.dot-orange  { background: #f97316; }
.dot-red     { background: #ef4444; }
.dot-green   { background: #22c55e; }
.dot-blue    { background: #3b82f6; }
.dot-grey    { background: #9ca3af; }
.dot-cyan    { background: #22d3ee; }
.dot-pink    { background: #ec4899; }
.dot-amber   { background: #f59e0b; }

.step-tip {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 10px;
  padding: 10px 14px;
  font-size: 0.85rem;
  color: #92400e;
  line-height: 1.5;
  margin-top: auto;
}
:global(.q-dark) .step-tip {
  background: #2d2a1e;
  border-color: #92400e;
  color: #fde68a;
}

/* Navigation */
.step-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #f3f4f6;
}
:global(.q-dark) .step-nav { border-color: #3a3a4e; }

.nav-dots {
  display: flex;
  gap: 5px;
}
.nav-dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  background: #e5e7eb;
  cursor: pointer;
  transition: all 0.25s;
}
:global(.q-dark) .nav-dot { background: #4a4a5e; }
.nav-dot-active { background: #667eea; transform: scale(1.4); }
.nav-dot-done   { background: #667eea; opacity: 0.5; }

/* ── COMPLETION CARD ─────────────────────────────────────────────── */
.onboarding-completion-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 40px;
  text-align: center;
  max-width: 480px;
  min-height: unset;
}

.completion-animation {
  position: relative;
  width: 120px; height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 28px;
}
.confetti-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 3px solid #667eea;
  animation: confetti-pulse 1.2s ease-out infinite;
}
@keyframes confetti-pulse {
  0%   { transform: scale(0.8); opacity: 1; }
  100% { transform: scale(1.6); opacity: 0; }
}
.completion-icon { font-size: 64px; line-height: 1; }

.completion-title {
  font-size: 1.7rem;
  font-weight: 900;
  color: #1a1a2e;
  margin: 0 0 12px;
}
:global(.q-dark) .completion-title { color: #e0e0e0; }

.completion-desc {
  font-size: 0.95rem;
  color: #6b7280;
  line-height: 1.7;
  margin: 0 0 28px;
}
.completion-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
  align-items: center;
}
.completion-actions .btn-cta { width: 100%; }

/* ── Transitions ─────────────────────────────────────────────────── */
.slide-up-enter-active, .slide-up-leave-active { transition: all 0.35s ease; }
.slide-up-enter-from { opacity: 0; transform: translateY(24px) scale(0.97); }
.slide-up-leave-to   { opacity: 0; transform: translateY(-16px) scale(0.97); }

.step-slide-enter-active, .step-slide-leave-active { transition: all 0.25s ease; }
.step-slide-enter-from { opacity: 0; transform: translateX(18px); }
.step-slide-leave-to   { opacity: 0; transform: translateX(-18px); }

.icon-pop-enter-active { transition: all 0.3s cubic-bezier(0.34,1.56,0.64,1); }
.icon-pop-leave-active { transition: all 0.15s ease; }
.icon-pop-enter-from   { opacity: 0; transform: scale(0.5); }
.icon-pop-leave-to     { opacity: 0; transform: scale(1.2); }

/* ── RTL ─────────────────────────────────────────────────────────── */
:global([dir='rtl']) .step-nav { flex-direction: row-reverse; }
:global([dir='rtl']) .welcome-body { text-align: right; }
:global([dir='rtl']) .step-bullets { direction: rtl; }

/* ── Responsive ──────────────────────────────────────────────────── */
@media (max-width: 640px) {
  .onboarding-welcome-card, .onboarding-step-card { flex-direction: column; }
  .welcome-hero { width: 100%; height: 180px; }
  .step-panel-left { width: 100%; height: 130px; }
  .sidebar-steps { flex-direction: row; top: auto; bottom: 10px; left: 50%; transform: translateX(-50%); }
  .sidebar-progress { display: none; }
  .step-bullets { grid-template-columns: 1fr; }
}
</style>
