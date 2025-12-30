<template>
  <q-table
    title="Schools"
    :rows="store.schools"
    :columns="columns"
    row-key="id"
    v-model:pagination="store.pagination"
    :loading="store.loading"
    @request="onRequest"
    binary-state-sort
  >
    <template v-slot:top-right>
      <q-btn color="primary" icon="add" label="New School" @click="$emit('create')" />
      <q-space />
      <q-input dense debounce="300" v-model="filter" placeholder="Search">
        <template v-slot:append>
          <q-icon name="search" />
        </template>
      </q-input>
    </template>

    <template v-slot:body-cell-actions="props">
      <q-td :props="props">
        <q-btn flat round color="primary" icon="edit" @click="$emit('edit', props.row)" />
        <q-btn flat round color="negative" icon="delete" @click="$emit('delete', props.row.id)" />
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { ref } from 'vue';
import { useSchoolStore } from 'src/stores/schools';

const store = useSchoolStore();
const filter = ref('');

const columns = [
  { name: 'name', label: 'Name', align: 'left', field: 'name', sortable: true },
  { name: 'address', label: 'Address', align: 'left', field: 'address' },
  { name: 'email', label: 'Email', align: 'left', field: 'email' },
  { name: 'actions', label: 'Actions', align: 'right' }
];

const emit = defineEmits(['create', 'edit', 'delete', 'request']);

const onRequest = (props) => {
  emit('request', props);
};
</script>
