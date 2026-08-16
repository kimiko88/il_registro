<template>
  <q-card>
    <q-card-section>
      <div class="text-h6">{{ title }}</div>
    </q-card-section>
    <q-card-section>
      <q-form @submit="onSubmit" class="q-gutter-md">
        <slot></slot>
        <div align="right">
          <q-btn :label="t('common.save') || 'Salva'" type="submit" color="primary" :loading="loading" />
        </div>
      </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps({
  title: String,
  loading: Boolean
});
const emit = defineEmits(['submit']);
const onSubmit = () => emit('submit');
</script>
