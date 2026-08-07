<template>
  <div class="q-gutter-y-xs">
    <div class="text-caption text-weight-bold text-slate-700 row items-center justify-between q-mb-xs">
      <span class="row items-center">
        <q-icon name="accessibility_new" color="primary" class="q-mr-xs" />
        Misure Compensative BES / DSA
      </span>
      <q-chip v-if="modelValue.length > 0" dense color="primary" text-color="white" size="xs">
        {{ modelValue.length }} selezionate
      </q-chip>
    </div>

    <!-- Chips list -->
    <div class="row q-gutter-xs wrap q-mb-sm">
      <q-chip
        v-for="measure in allMeasures"
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

    <!-- Custom Measure Add Section for Teachers & Secretary -->
    <div v-if="!readOnly" class="row items-center q-gutter-xs q-mt-xs">
      <q-input
        v-model="customInput"
        dense
        outlined
        placeholder="+ Aggiungi misura personalizzata (es. Software per mappe)"
        class="col"
        @keydown.enter.prevent="addCustomMeasure"
      />
      <q-btn
        color="secondary"
        icon="add"
        label="Aggiungi"
        dense
        unelevated
        no-caps
        @click="addCustomMeasure"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

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

const customInput = ref('')

const defaultMeasures = [
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

const customMeasures = ref([])

const allMeasures = computed(() => [
  ...defaultMeasures,
  ...customMeasures.value
])

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

function addCustomMeasure() {
  const val = customInput.value.trim()
  if (!val) return
  const slug = 'custom_' + val.toLowerCase().replace(/\s+/g, '_')
  
  if (!customMeasures.value.some(m => m.value === slug)) {
    customMeasures.value.push({ value: slug, label: val })
  }
  
  toggle(slug)
  customInput.value = ''
}
</script>
