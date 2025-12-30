<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">User Registry</div>
    <UserTable :users="users" :loading="loading" @import="showImport = true" />
    
    <q-dialog v-model="showImport">
        <q-card style="min-width: 300px">
            <q-card-section>
                <div class="text-h6">Import Users</div>
            </q-card-section>
            <q-card-section>
                <q-file v-model="importFile" label="Select CSV/XLSX" />
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Cancel" v-close-popup />
                <q-btn color="primary" label="Upload" @click="handleImport" />
            </q-card-actions>
        </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import UserTable from 'src/components/Secretary/UserTable.vue';
import { useUserManagement } from 'src/composables/useUserManagement';

const { users, loading, fetchUsers, importUsers } = useUserManagement();
const showImport = ref(false);
const importFile = ref(null);

onMounted(fetchUsers);

const handleImport = async () => {
    if (importFile.value) {
        await importUsers(importFile.value);
        showImport.value = false;
    }
};
</script>
