<template>
  <q-table
    :title="t('classRegister.studentList') || 'Elenco Studenti'"
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
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps({
    students: { type: Array, default: () => [] }
});

const columns = computed(() => [
    { name: 'name', label: t('classRegister.tableHeaderStudent') || 'Nome Studente', align: 'left', field: 'name', sortable: true },
    { name: 'dob', label: t('classRegister.dateLabel') || 'Data di Nascita', align: 'left', field: 'dob' },
    { name: 'needs', label: t('common.notes') || 'Esigenze Specifiche', align: 'left', field: 'needs' },
    { name: 'actions', label: t('common.actions') || 'Azioni', align: 'right' }
]);
</script>
