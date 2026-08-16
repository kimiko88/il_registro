<template>
  <q-dialog v-model="isOpen">
    <q-card style="min-width: 300px">
        <q-card-section>
            <div class="text-h6">{{ t('common.save') }}</div>
            <div class="text-caption">{{ t('gradesPage.date') }}: {{ date }}</div>
        </q-card-section>

        <q-card-section>
            <q-form @submit="onSubmit">
                <q-select 
                    v-model="reason" 
                    :options="['Motivi di Salute', 'Motivi Familiari', 'Altro']" 
                    :label="t('notesPage.descriptionLabel')" 
                    filled 
                    class="q-mb-md" 
                />
                <q-input 
                    v-model="notes" 
                    :label="t('gradesPage.notes')" 
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
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const props = defineProps(['modelValue', 'date'])
const emit = defineEmits(['update:modelValue', 'submit'])

const isOpen = ref(props.modelValue)
const reason = ref('Motivi di Salute')
const notes = ref('')

watch(() => props.modelValue, (val) => isOpen.value = val)
watch(isOpen, (val) => emit('update:modelValue', val))

const onSubmit = () => {
    emit('submit', { reason: reason.value, notes: notes.value })
    isOpen.value = false
}
</script>
