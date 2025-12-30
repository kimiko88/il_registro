<template>
  <q-card style="min-width: 400px">
    <q-card-section>
      <div class="text-h6">{{ isEdit ? 'Edit School' : 'New School' }}</div>
    </q-card-section>

    <q-card-section class="q-pt-none">
      <q-form @submit="onSubmit" class="q-gutter-md">
        <q-input v-model="form.name" label="Name" :rules="[val => !!val || 'Required']" />
        <q-input v-model="form.code" label="Min. Code" :rules="[val => !!val || 'Required']" />
        <q-input v-model="form.address" label="Address" />
        <q-input v-model="form.email" label="Email" type="email" />
        <q-input v-model="form.phone" label="Phone" />

        <div align="right">
          <q-btn flat label="Cancel" color="primary" v-close-popup />
          <q-btn label="Save" type="submit" color="primary" />
        </div>
      </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, onMounted } from 'vue';

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
