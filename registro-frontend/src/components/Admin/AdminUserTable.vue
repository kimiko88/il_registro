<template>
  <q-table
    title="Admin Users"
    :rows="store.admins"
    :columns="columns"
    row-key="id"
    v-model:pagination="store.pagination"
    :loading="store.loading"
    @request="onRequest"
  >
    <template v-slot:top-right>
      <q-btn color="primary" icon="add" label="New Admin" @click="$emit('create')" />
    </template>
    <template v-slot:body-cell-actions="props">
      <q-td :props="props">
        <q-btn flat round icon="lock_reset" @click="$emit('reset-password', props.row.id)" />
        <q-btn flat round color="negative" icon="delete" @click="$emit('delete', props.row.id)" />
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { useAdminStore } from 'src/stores/admin';

const store = useAdminStore();
const columns = [
  { name: 'email', label: 'Email', field: 'email', align: 'left' },
  { name: 'school', label: 'School', field: row => row.school?.name, align: 'left' },
  { name: 'actions', label: 'Actions', align: 'right' }
];

const emit = defineEmits(['create', 'delete', 'reset-password', 'request']);

const onRequest = (props) => {
  emit('request', props);
};
</script>
