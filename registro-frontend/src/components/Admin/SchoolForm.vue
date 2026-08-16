<template>
  <q-card style="min-width: 400px">
    <q-card-section>
      <div class="text-h6">{{ isEdit ? (t('common.edit') || 'Modifica Scuola') : (t('common.add') || 'Nuova Scuola') }}</div>
    </q-card-section>

    <q-card-section class="q-pt-none">
      <q-form @submit="onSubmit" class="q-gutter-md">
        <q-input v-model="form.name" :label="t('common.name') || 'Nome'" :rules="[val => !!val || (t('common.requiredField') || 'Campo obbligatorio')]" />
        <q-input v-model="form.code" :label="t('common.code') || 'Codice Meccanografico'" :rules="[val => !!val || (t('common.requiredField') || 'Campo obbligatorio')]" />
        <q-input v-model="form.address" :label="t('common.address') || 'Indirizzo'" />
        <q-input v-model="form.email" :label="t('login.emailLabel') || 'Email'" type="email" />
        <q-input v-model="form.phone" :label="t('login.phone') || 'Telefono'" />

        <div align="right">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" color="primary" v-close-popup />
          <q-btn :label="t('common.save') || 'Salva'" type="submit" color="primary" />
        </div>
      </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const props = defineProps({
  school: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['submit']);
const isEdit = !!props.school;

const form = ref({
  name: '',
  code: '',
  address: '',
  email: '',
  phone: ''
});

onMounted(() => {
  if (props.school) {
    form.value = { ...props.school };
  }
});

const onSubmit = () => {
  emit('submit', form.value);
};
</script>
