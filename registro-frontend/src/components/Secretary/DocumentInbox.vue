<template>
  <q-table
    title="Inbox"
    :rows="store.inbox"
    :columns="columns"
    row-key="id"
    v-model:pagination="store.pagination"
    :loading="store.loading"
    @request="onRequest"
  >
    <template v-slot:body-cell-status="props">
        <q-td :props="props">
            <q-badge :color="getStatusColor(props.value)" :label="props.value" />
        </q-td>
    </template>
    <template v-slot:body-cell-actions="props">
      <q-td :props="props">
        <q-btn flat round icon="visibility" @click="$emit('preview', props.row)" />
        <q-btn flat round color="primary" icon="rate_review" @click="$emit('review', props.row)" />
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { useDocumentsStore } from 'src/stores/documents';

const store = useDocumentsStore();

const columns = [
  { name: 'title', label: 'Title', field: 'title', align: 'left' },
  { name: 'type', label: 'Type', field: 'type', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'date', label: 'Date', field: 'created_at', align: 'left' },
  { name: 'actions', label: 'Actions', align: 'right' }
];

const emit = defineEmits(['preview', 'review', 'request']);

const onRequest = (props) => emit('request', props);

const getStatusColor = (status) => {
    if (status === 'Pending') return 'orange';
    if (status === 'Approved') return 'green';
    return 'grey';
};
</script>
