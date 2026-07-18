<template>
  <div class="q-pa-md">
    <div class="row q-mb-md justify-between items-center">
      <div class="text-h6">User Management</div>
      <q-btn color="primary" icon="add" label="Add User" @click="openAddUserDialog" />
    </div>

    <q-table
      :rows="users"
      :columns="columns"
      row-key="id"
      :loading="loading"
    >
      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn flat round color="primary" icon="edit" size="sm" @click="editUser(props.row)" />
          <q-btn flat round color="negative" icon="delete" size="sm" @click="confirmDelete(props.row)" />
        </q-td>
      </template>
    </q-table>

    <confirm-dialog
      v-model="showConfirm"
      message="Are you sure you want to delete this user?"
      @confirm="deleteUser"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'
import { useQuasar } from 'quasar'
import adminService from '@/services/adminService'
import ConfirmDialog from '@/components/Common/ConfirmDialog.vue'

const $q = useQuasar()
const { data: users, loading, fetch } = useApi('/users')
const showConfirm = ref(false)
const selectedUser = ref(null)

const columns = [
  { name: 'id', label: 'ID', field: 'id', sortable: true },
  { name: 'firstName', label: 'First Name', field: 'first_name', sortable: true },
  { name: 'lastName', label: 'Last Name', field: 'last_name', sortable: true },
  { name: 'email', label: 'Email', field: 'email', sortable: true },
  { name: 'role', label: 'Role', field: 'role', sortable: true },
  { name: 'actions', label: 'Actions', field: 'actions' }
]

onMounted(() => {
  fetch()
})

function openAddUserDialog() {
  // logic to open add dialog
}

function editUser(_user) {
  // logic to open edit dialog
}

function confirmDelete(user) {
  selectedUser.value = user
  showConfirm.value = true
}

async function deleteUser() {
  if (!selectedUser.value) return
  try {
    await adminService.deleteUser(selectedUser.value.id)
    $q.notify({ type: 'positive', message: 'User deleted successfully' })
    fetch() // Reload list
  } catch (error) {
    $q.notify({ type: 'negative', message: 'Error deleting user' })
  }
}
</script>
