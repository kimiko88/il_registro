<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">{{ title }}</div>
    </q-card-section>
    <q-card-section>
      <q-form @submit="onSubmit" class="q-gutter-md">
        <slot></slot>
        <div align="right">
          <q-btn label="Save" type="submit" color="primary" :loading="loading" />
        </div>
      </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
defineProps({
  title: String,
  loading: Boolean
});
const emit = defineEmits(['submit']);
const onSubmit = () => emit('submit');
</script>
