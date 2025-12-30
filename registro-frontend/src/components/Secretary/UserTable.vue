<template>
  <q-table :rows="users" :columns="columns" :loading="loading" title="Users">
    <template v-slot:top-right>
        <q-btn color="primary" label="Import" icon="upload" @click="$emit('import')" />
        <q-btn flat icon="add" label="New User" @click="$emit('create')" />
    </template>
  </q-table>
</template>

<script setup>
const props = defineProps(['users', 'loading']);
const columns = [
    { name: 'name', label: 'Name', field: 'name' },
    { name: 'role', label: 'Role', field: 'role' },
    { name: 'email', label: 'Email', field: 'email' }
];
</script>
