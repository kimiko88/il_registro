<template>
  <div class="dashboard-quick-actions q-mb-xl">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-subtitle1 text-weight-bold text-dark">{{ title }}</div>
      <slot name="header-actions" />
    </div>

    <div class="row q-col-gutter-md">
      <div
        v-for="action in visibleActions"
        :key="action.key"
        class="col-6 col-sm-4 col-md-3 col-lg-2"
      >
        <q-card
          flat
          bordered
          clickable
          class="quick-action-card cursor-pointer rounded-xl transition-all full-height"
          :class="[
            action.color ? `quick-action-card--${action.color}` : 'quick-action-card--primary',
            { 'quick-action-card--alert': action.badge > 0 && action.alertOnBadge }
          ]"
          role="button"
          :aria-label="action.label"
          @click="handleAction(action)"
          @keydown.enter.prevent="handleAction(action)"
          @keydown.space.prevent="handleAction(action)"
          tabindex="0"
        >
          <q-card-section class="column items-center justify-center text-center q-pa-md" style="min-height: 96px;">
            <!-- Badge contatore -->
            <div class="relative-position q-mb-sm">
              <q-icon :name="action.icon" size="28px" :class="`text-${action.color || 'primary'}-600`" aria-hidden="true" />
              <q-badge
                v-if="action.badge > 0"
                :label="action.badge > 99 ? '99+' : String(action.badge)"
                :color="action.alertOnBadge ? 'negative' : (action.color || 'primary')"
                floating
                rounded
                class="quick-action-badge"
              />
            </div>
            <div class="text-caption text-weight-bold text-dark ellipsis-2-lines" style="line-height: 1.3;">
              {{ action.label }}
            </div>
            <div v-if="action.subtitle" class="text-caption text-grey-6 q-mt-xs ellipsis">
              {{ action.subtitle }}
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  /** Titolo della sezione (es. "Azioni Rapide") */
  title: {
    type: String,
    default: 'Azioni Rapide'
  },
  /**
   * Lista delle azioni da mostrare.
   * Ogni azione: { key, label, icon, color?, path?, action?, badge?, alertOnBadge?, subtitle?, hidden? }
   * - key: identificatore univoco
   * - label: testo breve
   * - icon: Material icon name
   * - color: colore Quasar (primary, secondary, positive, negative, warning, info, indigo, teal, amber…)
   * - path: percorso router (alternativo ad action)
   * - action: funzione da eseguire al click (alternativo a path)
   * - badge: numero da mostrare come badge (0 = nessun badge)
   * - alertOnBadge: se true il badge viene in rosso (negative)
   * - subtitle: testo piccolo sotto la label
   * - hidden: se true la card non viene mostrata
   */
  actions: {
    type: Array,
    default: () => []
  },
  /** Numero massimo di azioni da mostrare (default: tutte) */
  maxVisible: {
    type: Number,
    default: Infinity
  }
})

const emit = defineEmits(['action'])
const router = useRouter()

const visibleActions = computed(() =>
  props.actions
    .filter(a => !a.hidden)
    .slice(0, props.maxVisible)
)

function handleAction(action) {
  if (action.action) {
    action.action()
  } else if (action.path || action.route) {
    const target = action.path || action.route
    if (router.currentRoute?.value?.path !== target) {
      router.push(target).catch(() => {})
    }
  }
  emit('action', action)
}
</script>

<style scoped>
.quick-action-card {
  transition: transform 0.15s ease, box-shadow 0.15s ease, background-color 0.15s ease;
  border-color: rgba(0, 0, 0, 0.08);
}

.quick-action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.10);
}

.quick-action-card:focus-visible {
  outline: 2px solid var(--q-primary);
  outline-offset: 2px;
}

/* Colori per tipologia */
.quick-action-card--primary { background: rgba(99, 102, 241, 0.04); }
.quick-action-card--primary:hover { background: rgba(99, 102, 241, 0.10); }

.quick-action-card--secondary { background: rgba(100, 116, 139, 0.04); }
.quick-action-card--secondary:hover { background: rgba(100, 116, 139, 0.10); }

.quick-action-card--positive { background: rgba(34, 197, 94, 0.04); }
.quick-action-card--positive:hover { background: rgba(34, 197, 94, 0.10); }

.quick-action-card--warning { background: rgba(245, 158, 11, 0.05); }
.quick-action-card--warning:hover { background: rgba(245, 158, 11, 0.12); }

.quick-action-card--negative { background: rgba(239, 68, 68, 0.04); }
.quick-action-card--negative:hover { background: rgba(239, 68, 68, 0.10); }

.quick-action-card--indigo { background: rgba(99, 102, 241, 0.06); }
.quick-action-card--indigo:hover { background: rgba(99, 102, 241, 0.13); }

.quick-action-card--teal { background: rgba(20, 184, 166, 0.04); }
.quick-action-card--teal:hover { background: rgba(20, 184, 166, 0.10); }

.quick-action-card--amber { background: rgba(245, 158, 11, 0.05); }
.quick-action-card--amber:hover { background: rgba(245, 158, 11, 0.12); }

/* Alert (badge con valore > 0 e alertOnBadge) */
.quick-action-card--alert {
  border-color: rgba(239, 68, 68, 0.3);
  background: rgba(239, 68, 68, 0.04);
}

.quick-action-badge {
  font-size: 10px;
  min-width: 18px;
  height: 18px;
}
</style>
