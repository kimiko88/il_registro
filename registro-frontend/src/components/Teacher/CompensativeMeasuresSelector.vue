<template>
  <div class="q-gutter-y-xs">
    <div class="text-caption text-weight-bold text-slate-700 row items-center justify-between q-mb-xs">
      <span>
        <q-icon name="accessibility_new" color="primary" class="q-mr-xs" />
        Misure Compensative BES / DSA
      </span>
      <q-chip v-if="modelValue.length > 0" dense color="primary" text-color="white" size="xs">
        {{ modelValue.length }} selezionate
      </q-chip>
    </div>

    <div class="row q-gutter-xs wrap">
      <q-chip
        v-for="measure in availableMeasures"
        :key="measure.value"
        clickable
        dense
        :color="isSelected(measure.value) ? 'primary' : 'grey-3'"
        :text-color="isSelected(measure.value) ? 'white' : 'grey-9'"
        :icon="isSelected(measure.value) ? 'check_circle' : 'add'"
        class="transition-all"
        @click="toggle(measure.value)"
      >
        {{ measure.label }}
      </q-chip>
    </div>
  </div>
</template>

<script setup>

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  readOnly: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const availableMeasures = [
  { value: 'calcolatrice', label: 'Uso Calcolatrice' },
  { value: 'tempo_aggiuntivo_30', label: 'Tempo Agg. (+30%)' },
  { value: 'tempo_aggiuntivo_50', label: 'Tempo Agg. (+50%)' },
  { value: 'prova_equipollente', label: 'Prova Equipollente' },
  { value: 'sintesi_vocale', label: 'Sintesi Vocale' },
  { value: 'mappe_concettuali', label: 'Mappe Concettuali' },
  { value: 'tavola_pitagorica', label: 'Tavola Pitagorica' },
  { value: 'tabelle_formule', label: 'Tabelle / Formulario' },
  { value: 'dizionario_ortografico', label: 'Dizionario Digitale' },
  { value: 'testo_ingrandito', label: 'Testo Ingrandito / High Contrast' }
]

function isSelected(val) {
  return props.modelValue.includes(val)
}

function toggle(val) {
  if (props.readOnly) return
  const current = [...props.modelValue]
  const idx = current.indexOf(val)
  if (idx >= 0) {
    current.splice(idx, 1)
  } else {
    current.push(val)
  }
  emit('update:modelValue', current)
}
</script>
