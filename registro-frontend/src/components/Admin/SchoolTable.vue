<template>
  <q-table
    :title="t('roleDashboards.schoolManagement') || 'Istituti Scolastici'"
    :rows="store.schools"
    :columns="columns"
    :filter="filter"
    row-key="id"
    v-model:pagination="store.pagination"
    :loading="store.loading"
    @request="onRequest"
    binary-state-sort
  >
    <template v-slot:top-right>
      <q-btn color="primary" icon="add" :label="t('common.add') || 'Nuova Scuola'" @click="$emit('create')" />
      <q-space />
      <q-input dense debounce="300" v-model="filter" :placeholder="t('common.search') || 'Cerca'">
        <template v-slot:append>
          <q-icon name="search" />
        </template>
      </q-input>
    </template>

    <template v-slot:body-cell-actions="props">
      <q-td :props="props">
        <q-btn flat round color="primary" icon="edit" :aria-label="t('common.edit') || 'Modifica'" @click="$emit('edit', props.row)">
          <q-tooltip>{{ t('common.edit') || 'Modifica' }}</q-tooltip>
        </q-btn>
        <q-btn flat round color="negative" icon="delete" :aria-label="t('common.delete') || 'Elimina'" @click="$emit('delete', props.row.id)">
          <q-tooltip>{{ t('common.delete') || 'Elimina' }}</q-tooltip>
        </q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useSchoolStore } from '@/stores/schools';

const { t } = useI18n();
const store = useSchoolStore();
const filter = ref('');

const columns = computed(() => [
  { name: 'name', label: t('common.name') || 'Nome', align: 'left', field: 'name', sortable: true },
  { name: 'address', label: t('common.address') || 'Indirizzo', align: 'left', field: 'address' },
  { name: 'email', label: t('login.emailLabel') || 'Email', align: 'left', field: 'email' },
  { name: 'actions', label: t('common.actions') || 'Azioni', align: 'right' }
]);

const emit = defineEmits(['create', 'edit', 'delete', 'request']);

const onRequest = (props) => {
  emit('request', props);
};
</script>
