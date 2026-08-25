<template>
  <div class="attendance-marker">
    <q-btn-group spread>
      <q-btn 
        color="positive" 
        :label="t('classRegister.present') || 'Presente'" 
        :outline="modelValue !== 'present'"
        @click="update('present')"
        data-test="btn-present"
      />
      <q-btn 
        color="warning" 
        :label="t('classRegister.late') || 'Ritardo'" 
        :outline="modelValue !== 'late'"
        @click="update('late')"
        data-test="btn-late"
      />
      <q-btn 
        color="negative" 
        :label="t('classRegister.absent') || 'Assente'" 
        :outline="modelValue !== 'absent'"
        @click="update('absent')"
        data-test="btn-absent"
      />
    </q-btn-group>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

defineProps({
  modelValue: {
    type: String,
    default: 'present'
  }
})

const emit = defineEmits(['update:modelValue'])

function update(status) {
  emit('update:modelValue', status)
}
</script>
