<template>
  <q-dialog v-model="isOpen">
    <q-card style="min-width: 300px">
        <q-card-section>
            <div class="text-h6">{{ t('attendance.justifyAbsence') || 'Giustifica Assenza' }}</div>
            <div class="text-caption">{{ t('gradesPage.date') }}: {{ date }}</div>
        </q-card-section>

        <q-card-section>
            <q-form @submit="onSubmit">
                <q-select 
                    v-model="reason" 
                    :options="reasonOptions" 
                    :label="t('notesPage.descriptionLabel') || 'Motivazione'" 
                    emit-value
                    map-options
                    filled 
                    class="q-mb-md" 
                />
                <q-input 
                    v-model="notes" 
                    :label="t('gradesPage.notes') || 'Note'" 
                    filled 
                    type="textarea" 
                    rows="3" 
                />
                <div class="row justify-end q-mt-md">
                    <q-btn flat :label="t('common.cancel')" v-close-popup />
                    <q-btn type="submit" :label="t('common.save')" color="primary" />
                </div>
            </q-form>
        </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const props = defineProps(['modelValue', 'date'])
const emit = defineEmits(['update:modelValue', 'submit'])

const isOpen = ref(props.modelValue)
const reason = ref('Motivi di Salute')
const notes = ref('')

const reasonOptions = computed(() => [
  { label: t('attendance.healthReasons') || 'Motivi di Salute', value: 'Motivi di Salute' },
  { label: t('attendance.familyReasons') || 'Motivi Familiari', value: 'Motivi Familiari' },
  { label: t('common.other') || 'Altro', value: 'Altro' }
])

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val) {
    reason.value = 'Motivi di Salute'
    notes.value = ''
  }
})
watch(isOpen, (val) => emit('update:modelValue', val))

const onSubmit = () => {
    emit('submit', { reason: reason.value, notes: notes.value })
    isOpen.value = false
}
</script>
