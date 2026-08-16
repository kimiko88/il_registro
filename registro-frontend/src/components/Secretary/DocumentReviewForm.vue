<template>
  <q-card style="min-width: 400px">
    <q-card-section>
      <div class="text-h6">{{ t('documentsPage.reviewTitle') || 'Revisiona Documento' }}</div>
      <div class="text-subtitle2">{{ doc?.title }}</div>
    </q-card-section>
    
    <q-card-section>
        <q-input 
            v-model="notes" 
            :label="t('common.notes') || 'Note di Revisione'" 
            type="textarea" 
            autogrow 
            :rules="[val => !!val || (t('common.requiredField') || 'Le note sono obbligatorie per il rifiuto')]"
        />
    </q-card-section>

    <q-card-actions align="right">
        <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup />
        <q-btn color="negative" :label="t('common.reject') || 'Rifiuta'" @click="emit('submit', 'rejected')" />
        <q-btn color="positive" :label="t('common.approve') || 'Approva'" @click="emit('submit', 'approved')" />
    </q-card-actions>
  </q-card>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const props = defineProps(['doc', 'modelValue']);
const emit = defineEmits(['update:modelValue', 'submit']);

const notes = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
});
</script>
