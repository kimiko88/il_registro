<template>
  <div class="q-pa-md q-gutter-sm">
    <q-editor
      v-model="internalContent"
      min-height="15rem"
      :toolbar="[
        ['bold', 'italic', 'strike', 'underline'],
        ['quote', 'unordered', 'ordered', 'outdent', 'indent'],
        ['undo', 'redo']
      ]"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps(['modelValue']);
const emit = defineEmits(['update:modelValue']);

const internalContent = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
});
</script>
