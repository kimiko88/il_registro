<template>
  <!-- Floating Action Button for Help -->
  <div class="help-fab-wrapper" :class="{ expanded: isMenuOpen }">
    <!-- Expanded menu items -->
    <transition name="fab-menu">
      <div v-if="isMenuOpen" class="fab-menu">
        <!-- Contact Support -->
        <div class="fab-menu-item" @click="goToSupport">
          <span class="fab-menu-label">{{ t('help.contactSupport') }}</span>
          <button class="fab-action-btn bg-orange">
            <q-icon name="support_agent" size="20px" color="white" />
          </button>
        </div>

        <!-- Open Guide (Drawer) -->
        <div class="fab-menu-item" @click="openHelp">
          <span class="fab-menu-label">{{ t('help.openHelp') }}</span>
          <button class="fab-action-btn bg-purple">
            <q-icon name="menu_book" size="20px" color="white" />
          </button>
        </div>

        <!-- Restart Tour -->
        <div class="fab-menu-item" @click="restartTour">
          <span class="fab-menu-label">{{ t('help.restartTour') }}</span>
          <button class="fab-action-btn bg-primary">
            <q-icon name="rocket_launch" size="20px" color="white" />
          </button>
        </div>
      </div>
    </transition>

    <!-- Main FAB button -->
    <button
      class="help-fab-btn"
      :class="{ open: isMenuOpen }"
      @click="toggleMenu"
      :aria-label="t('help.fabTooltip')"
      :aria-expanded="isMenuOpen"
    >
      <q-icon :name="isMenuOpen ? 'close' : 'help_outline'" size="26px" color="white" />
      <q-tooltip anchor="top left" self="bottom left" :offset="[0, 8]">
        {{ t('help.fabTooltip') }}
      </q-tooltip>
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const emit = defineEmits(['open-help', 'restart-tour'])
const router = useRouter()
const { t } = useI18n()

const isMenuOpen = ref(false)

function toggleMenu() {
  isMenuOpen.value = !isMenuOpen.value
}

function openHelp() {
  isMenuOpen.value = false
  emit('open-help')
}

function restartTour() {
  isMenuOpen.value = false
  emit('restart-tour')
}

function goToSupport() {
  isMenuOpen.value = false
  router.push('/support')
}
</script>

<style scoped>
.help-fab-wrapper {
  position: fixed;
  bottom: 28px;
  right: 28px;
  z-index: 5000;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 10px;
}

/* RTL support */
:global([dir='rtl']) .help-fab-wrapper {
  right: auto;
  left: 28px;
  align-items: flex-start;
}

.help-fab-btn {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  outline: none;
}

.help-fab-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 8px 28px rgba(102, 126, 234, 0.65);
}

.help-fab-btn.open {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  transform: rotate(0deg);
  box-shadow: 0 6px 20px rgba(239, 68, 68, 0.5);
}

/* FAB Menu items */
.fab-menu {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-end;
}

:global([dir='rtl']) .fab-menu {
  align-items: flex-start;
}

.fab-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

:global([dir='rtl']) .fab-menu-item {
  flex-direction: row-reverse;
}

.fab-menu-label {
  background: white;
  color: #1a1a2e;
  font-size: 0.82rem;
  font-weight: 600;
  padding: 7px 14px;
  border-radius: 20px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.15);
  white-space: nowrap;
  cursor: pointer;
}

:global(.q-dark) .fab-menu-label {
  background: #2d2d3f;
  color: #e0e0e0;
  box-shadow: 0 2px 12px rgba(0,0,0,0.35);
}

.fab-action-btn {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 14px rgba(0,0,0,0.2);
  transition: transform 0.2s;
  flex-shrink: 0;
}

.fab-action-btn:hover {
  transform: scale(1.1);
}

.bg-primary { background: #667eea; }
.bg-purple { background: #7c3aed; }
.bg-orange { background: #f97316; }

/* Transition */
.fab-menu-enter-active {
  animation: fabMenuIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.fab-menu-leave-active {
  animation: fabMenuOut 0.2s ease forwards;
}

@keyframes fabMenuIn {
  from {
    opacity: 0;
    transform: translateY(12px) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes fabMenuOut {
  from {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
  to {
    opacity: 0;
    transform: translateY(8px) scale(0.9);
  }
}
</style>
