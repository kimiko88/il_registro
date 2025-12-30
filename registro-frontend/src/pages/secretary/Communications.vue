<template>
  <q-page padding>
    <div class="text-h4">Communications</div>
    <q-editor v-model="content" min-height="10rem" />
    <q-btn class="q-mt-md" color="primary" label="Send Circular" />
  </q-page>
</template>

<script setup>
import { ref } from 'vue';
const content = ref('');
</script>
