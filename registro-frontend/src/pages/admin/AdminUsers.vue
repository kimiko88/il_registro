<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">Admin Users</div>
    <AdminUserTable 
      @create="openCreate" 
      @delete="confirmDelete"
      @reset-password="resetPassword"
      @request="onRequest"
    />
    <!-- Dialogs omitted for brevity -->
  </q-page>
</template>

<script setup>
import AdminUserTable from 'src/components/Admin/AdminUserTable.vue';
import { useAdminUsers } from 'src/composables/useAdminUsers';

const { openCreate, confirmDelete, resetPassword, onRequest } = useAdminUsers();
</script>
