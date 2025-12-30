<template>
  <q-card style="min-width: 400px">
    <q-card-section>
      <div class="text-h6">Review Document</div>
      <div class="text-subtitle2">{{ doc?.title }}</div>
    </q-card-section>
    
    <q-card-section>
        <q-input 
            v-model="notes" 
            label="Review Notes" 
            type="textarea" 
            autogrow 
            :rules="[val => !!val || 'Notes are required for rejection']"
        />
    </q-card-section>

    <q-card-actions align="right">
        <q-btn flat label="Cancel" v-close-popup />
        <q-btn color="negative" label="Reject" @click="emit('submit', 'rejected')" />
        <q-btn color="positive" label="Approve" @click="emit('submit', 'approved')" />
    </q-card-actions>
  </q-card>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps(['doc', 'modelValue']);
const emit = defineEmits(['update:modelValue', 'submit']);

const notes = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
});
</script>
