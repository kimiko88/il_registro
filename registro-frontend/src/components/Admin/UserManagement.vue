<template>
  <div class="q-pa-md">
    <div class="row q-mb-md justify-between items-center">
      <div class="text-h6">{{ t('roleDashboards.userManagement') || 'Gestione Utenti' }}</div>
      <q-btn color="primary" icon="add" :label="t('common.add') || 'Aggiungi Utente'" @click="openAddUserDialog" />
    </div>

    <q-table
      :rows="users || []"
      :columns="columns"
      row-key="id"
      :loading="loading"
    >
      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn flat round color="primary" icon="edit" size="sm" :aria-label="t('common.edit') || 'Modifica'" @click="editUser(props.row)">
            <q-tooltip>{{ t('common.edit') || 'Modifica' }}</q-tooltip>
          </q-btn>
          <q-btn flat round color="negative" icon="delete" size="sm" :aria-label="t('common.delete') || 'Elimina'" @click="confirmDelete(props.row)">
            <q-tooltip>{{ t('common.delete') || 'Elimina' }}</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <confirm-dialog
      v-model="showConfirm"
      :message="t('common.confirmDelete') || 'Sei sicuro di voler eliminare questo utente?'"
      @confirm="deleteUser"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useApi } from '@/composables/useApi'
import { useQuasar } from 'quasar'
import adminService from '@/services/adminService'
import ConfirmDialog from '@/components/Common/ConfirmDialog.vue'

const { t } = useI18n()
const $q = useQuasar()
const { data: users, loading, fetch } = useApi('/users')
const showConfirm = ref(false)
const selectedUser = ref(null)

const columns = computed(() => [
  { name: 'id', label: 'ID', field: 'id', sortable: true },
  { name: 'firstName', label: t('common.name') || 'Nome', field: 'first_name', sortable: true },
  { name: 'lastName', label: t('common.surname') || 'Cognome', field: 'last_name', sortable: true },
  { name: 'email', label: t('login.emailLabel') || 'Email', field: 'email', sortable: true },
  { name: 'role', label: t('common.role') || 'Ruolo', field: 'role', sortable: true },
  { name: 'actions', label: t('common.actions') || 'Azioni', field: 'actions' }
])

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
    $q.notify({ type: 'positive', message: t('admin.userDeleted') || 'Utente eliminato con successo' })
    fetch() // Reload list
  } catch (error) {
    $q.notify({ type: 'negative', message: t('admin.userDeleteError') || 'Errore durante l\'eliminazione dell\'utente' })
  }
}
</script>
