<template>
  <q-badge
    :color="badgeColor"
    :text-color="badgeTextColor"
    :dense="dense"
    class="grade-badge shadow-xs text-weight-bold"
    :class="[
      dense ? 'q-px-xs q-py-none text-caption' : 'q-px-sm q-py-xs text-body2',
      statusClass
    ]"
    role="status"
    :aria-label="computedAriaLabel"
  >
    <q-icon
      v-if="showIcon && iconName"
      :name="iconName"
      :size="dense ? '11px' : '14px'"
      class="q-mr-xs"
      aria-hidden="true"
    />

    <span class="grade-value">{{ displayValue }}</span>

    <span
      v-if="showType && displayTypeShort"
      class="q-ml-xs text-caption opacity-80"
      aria-hidden="true"
    >
      ({{ displayTypeShort }})
    </span>

    <q-tooltip
      v-if="tooltip && (hasTooltipContent || isAbsent)"
      anchor="top middle"
      self="bottom middle"
      class="bg-dark text-white text-caption shadow-4 rounded-borders"
    >
      <div class="column q-gutter-xs">
        <div class="text-weight-bold">
          {{ isAbsent ? (t('gradesPage.absent') || 'Assente alla prova') : `${t('classRegister.tableHeaderGrade') || 'Voto'}: ${displayValue}` }}
        </div>
        <div v-if="displaySubject">
          {{ t('gradesPage.subject') || 'Materia' }}: {{ displaySubject }}
        </div>
        <div v-if="displayTypeLabel">
          {{ t('classRegister.tableHeaderGradeType') || 'Tipo' }}: {{ displayTypeLabel }}
        </div>
        <div v-if="displayDate">
          {{ t('classRegister.dateLabel') || 'Data' }}: {{ displayDate }}
        </div>
        <div v-if="displayWeight && displayWeight !== 1">
          {{ t('gradesPage.weight') || 'Peso' }}: {{ displayWeight }}x
        </div>
        <div v-if="displayNotes" class="text-italic">
          {{ displayNotes }}
        </div>
      </div>
    </q-tooltip>
  </q-badge>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  value: {
    type: [Number, String, Object],
    default: null
  },
  dense: {
    type: Boolean,
    default: false
  },
  showIcon: {
    type: Boolean,
    default: true
  },
  showType: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: ''
  },
  weight: {
    type: Number,
    default: null
  },
  date: {
    type: String,
    default: ''
  },
  subject: {
    type: String,
    default: ''
  },
  notes: {
    type: String,
    default: ''
  },
  tooltip: {
    type: Boolean,
    default: true
  },
  ariaLabel: {
    type: String,
    default: ''
  }
})

let t = (key, fallback) => (typeof fallback === 'string' ? fallback : key)
try {
  const i18n = useI18n()
  if (i18n && i18n.t) {
    t = i18n.t
  }
} catch {
  // Test fallback
}

const rawGradeValue = computed(() => {
  if (props.value !== null && typeof props.value === 'object') {
    return props.value.grade_value ?? props.value.value ?? null
  }
  return props.value
})

const isAbsent = computed(() => {
  const v = rawGradeValue.value
  return v === -1 || v === '-1' || v === 'A' || v === 'a' || v === 'Assente'
})

const numericValue = computed(() => {
  if (isAbsent.value) return null
  const v = rawGradeValue.value
  if (v === null || v === undefined || v === '') return null
  const num = Number(v)
  return isNaN(num) ? null : num
})

const isSufficient = computed(() => {
  if (numericValue.value === null) return null
  return numericValue.value >= 6
})

const displayValue = computed(() => {
  if (isAbsent.value) return 'A'
  if (numericValue.value !== null) {
    // Se intero, mostra intero; altrimenti max 2 decimali senza zeri inutili
    return Number.isInteger(numericValue.value)
      ? String(numericValue.value)
      : String(Math.round(numericValue.value * 100) / 100)
  }
  return rawGradeValue.value !== null && rawGradeValue.value !== undefined ? String(rawGradeValue.value) : '—'
})

const badgeColor = computed(() => {
  if (isAbsent.value) return 'grey-7'
  if (isSufficient.value === null) return 'grey-5'
  return isSufficient.value ? 'positive' : 'negative'
})

const badgeTextColor = computed(() => 'white')

const statusClass = computed(() => {
  if (isAbsent.value) return 'grade-absent'
  if (isSufficient.value === null) return 'grade-unknown'
  return isSufficient.value ? 'grade-sufficient' : 'grade-insufficient'
})

const iconName = computed(() => {
  if (isAbsent.value) return 'event_busy'
  if (isSufficient.value === null) return 'help_outline'
  return isSufficient.value ? 'check_circle' : 'priority_high'
})

// Metadata helpers (possono provenire dall'oggetto o da props esplicite)
const displayType = computed(() => {
  if (props.type) return props.type
  if (props.value && typeof props.value === 'object') {
    return props.value.evaluation_type || props.value.grade_type || props.value.type || ''
  }
  return ''
})

const displayTypeShort = computed(() => {
  const type = String(displayType.value).toLowerCase()
  if (type.includes('writ') || type.includes('scrit')) return 'S'
  if (type.includes('oral')) return 'O'
  if (type.includes('prac') || type.includes('prat')) return 'P'
  return ''
})

const displayTypeLabel = computed(() => {
  const type = String(displayType.value).toLowerCase()
  if (type.includes('writ') || type.includes('scrit')) return t('gradesPage.types.written', 'Scritto')
  if (type.includes('oral')) return t('gradesPage.types.oral', 'Orale')
  if (type.includes('prac') || type.includes('prat')) return t('gradesPage.types.practical', 'Pratico')
  return displayType.value
})

const displaySubject = computed(() => {
  if (props.subject) return props.subject
  if (props.value && typeof props.value === 'object') {
    return props.value.subject_name || props.value.subject || ''
  }
  return ''
})

const displayDate = computed(() => {
  if (props.date) return props.date
  if (props.value && typeof props.value === 'object') {
    const rawDate = props.value.date
    if (rawDate && typeof rawDate === 'string') {
      return rawDate.split('T')[0]
    }
  }
  return ''
})

const displayWeight = computed(() => {
  if (props.weight !== null && props.weight !== undefined) return props.weight
  if (props.value && typeof props.value === 'object' && props.value.weight !== undefined) {
    return props.value.weight
  }
  return null
})

const displayNotes = computed(() => {
  if (props.notes) return props.notes
  if (props.value && typeof props.value === 'object') {
    return props.value.notes || props.value.description || ''
  }
  return ''
})

const hasTooltipContent = computed(() => {
  return Boolean(
    displaySubject.value ||
    displayTypeLabel.value ||
    displayDate.value ||
    (displayWeight.value && displayWeight.value !== 1) ||
    displayNotes.value
  )
})

const computedAriaLabel = computed(() => {
  if (props.ariaLabel) return props.ariaLabel
  if (isAbsent.value) {
    return `${t('classRegister.tableHeaderGrade') || 'Voto'}: ${t('gradesPage.absent') || 'Assente'}`
  }
  const qual = isSufficient.value === null
    ? ''
    : isSufficient.value
      ? ` - ${t('gradesPage.sufficient') || 'Sufficiente'}`
      : ` - ${t('gradesPage.insufficient') || 'Insufficiente'}`
  const typeStr = displayTypeLabel.value ? `, ${displayTypeLabel.value}` : ''
  return `${t('classRegister.tableHeaderGrade') || 'Voto'}: ${displayValue.value}${qual}${typeStr}`
})
</script>

<style scoped>
.grade-badge {
  display: inline-flex;
  align-items: center;
  border-radius: 9999px;
  font-family: inherit;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
  user-select: none;
}

.grade-badge:hover {
  transform: translateY(-1px);
}

.grade-sufficient {
  background-color: var(--q-positive, #21ba45);
}

.grade-insufficient {
  background-color: var(--q-negative, #c10015);
}

.grade-absent {
  background-color: #5c6bc0;
}
</style>
