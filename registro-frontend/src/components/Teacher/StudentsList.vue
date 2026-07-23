<template>
  <q-table
    title="Students"
    :rows="students"
    :columns="columns"
    row-key="id"
    flat bordered
  >
    <template v-slot:body-cell-needs="props">
        <q-td :props="props">
            <q-badge v-if="props.value !== 'None'" color="orange">{{ props.value }}</q-badge>
            <span v-else class="text-grey">-</span>
        </q-td>
    </template>
    <template v-slot:body-cell-actions="props">
        <q-td :props="props">
            <q-btn flat round icon="info" size="sm" color="info" />
            <q-btn flat round icon="history_edu" size="sm" color="primary" />
        </q-td>
    </template>
  </q-table>
</template>

<script setup>
defineProps({
    students: { type: Array, default: () => [] }
});

const columns = [
    { name: 'name', label: 'Name', align: 'left', field: 'name', sortable: true },
    { name: 'dob', label: 'Date of Birth', align: 'left', field: 'dob' },
    { name: 'needs', label: 'Specific Needs', align: 'left', field: 'needs' },
    { name: 'actions', label: 'Actions', align: 'right' }
];
</script>
