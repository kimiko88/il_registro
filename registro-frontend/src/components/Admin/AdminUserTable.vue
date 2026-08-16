<template>
  <q-table
    :title="t('roleDashboards.userManagement') || 'Amministratori'"
    :rows="store.admins"
    :columns="columns"
    row-key="id"
    v-model:pagination="store.pagination"
    :loading="store.loading"
    @request="onRequest"
  >
    <template v-slot:top-right>
      <q-btn color="primary" icon="add" :label="t('common.add') || 'Nuovo Admin'" @click="$emit('create')" />
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
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAdminStore } from '@/stores/admin';

const { t } = useI18n();
const store = useAdminStore();
const columns = computed(() => [
  { name: 'email', label: t('login.emailLabel') || 'Email', field: 'email', align: 'left' },
  { name: 'school', label: t('login.selectSchool') || 'Scuola', field: row => row.school?.name, align: 'left' },
  { name: 'actions', label: t('common.actions') || 'Azioni', align: 'right' }
]);

const emit = defineEmits(['create', 'delete', 'reset-password', 'request']);

const onRequest = (props) => {
  emit('request', props);
};
</script>
