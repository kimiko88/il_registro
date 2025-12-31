<template>
  <q-card>
    <q-table
      :rows="users"
      :columns="columns"
      :loading="loading"
      :filter="filter"
      row-key="id"
      selection="multiple"
      v-model:selected="selected"
    >
      <template v-slot:top>
        <div class="text-h6 q-mr-md">Lista Utenti</div>
        
        <q-btn-toggle
          v-model="roleFilter"
          push
          glossy
          toggle-color="primary"
          :options="[
            {label: 'Tutti', value: 'all'},
            {label: 'Studenti', value: 'student'},
            {label: 'Docenti', value: 'teacher'},
            {label: 'Genitori', value: 'parent'},
            {label: 'Personale', value: 'staff'}
          ]"
          class="q-mr-md"
          @update:model-value="$emit('filter-role', $event)"
        />

        <q-space />
        
        <q-input dense debounce="300" v-model="filter" placeholder="Cerca...">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
        
        <q-btn color="primary" icon="add" label="Nuovo" class="q-ml-md" @click="$emit('create')" />
        <q-btn flat round icon="cloud_download" color="primary" @click="$emit('export')" tooltip="Export CSV" />
      </template>

      <!-- Bulk Actions -->
      <template v-slot:top-row v-if="selected.length > 0">
         <q-tr>
           <q-td colspan="100%">
             <div class="row items-center q-gutter-sm bg-blue-1 q-pa-sm rounded-borders">
               <span class="text-weight-bold text-primary">{{ selected.length }} selezionati</span>
               <q-space />
               <q-btn size="sm" color="negative" icon="delete" label="Elimina Massa" @click="$emit('bulk-delete', selected)" />
               <q-btn size="sm" color="warning" text-color="dark" icon="lock_reset" label="Reset Pwd Massa" @click="$emit('bulk-reset', selected)" />
             </div>
           </q-td>
         </q-tr>
      </template>

      <!-- Custom Body -->
      <template v-slot:body-cell-role="props">
        <q-td :props="props">
          <q-chip :color="getRoleColor(props.value)" text-color="white" size="sm">
            {{ getRoleLabel(props.value) }}
          </q-chip>
        </q-td>
      </template>

      <template v-slot:body-cell-status="props">
         <q-td :props="props">
            <q-badge :color="props.row.active ? 'positive' : 'grey'">
                {{ props.row.active ? 'Attivo' : 'Inattivo' }}
            </q-badge>
         </q-td>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props" auto-width>
          <q-btn flat round size="sm" color="grey-7" icon="more_vert">
            <q-menu>
              <q-list style="min-width: 150px">
                <q-item clickable v-close-popup @click="$emit('edit', props.row)">
                  <q-item-section avatar><q-icon name="edit" /></q-item-section>
                  <q-item-section>Modifica</q-item-section>
                </q-item>
                <q-item clickable v-close-popup @click="$emit('reset-pwd', props.row)">
                  <q-item-section avatar><q-icon name="lock_reset" /></q-item-section>
                  <q-item-section>Reset Password</q-item-section>
                </q-item>
                <q-separator />
                <q-item clickable v-close-popup class="text-negative" @click="$emit('delete', props.row)">
                  <q-item-section avatar><q-icon name="delete" /></q-item-section>
                  <q-item-section>Elimina</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-td>
      </template>
    </q-table>
  </q-card>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps(['users', 'loading']);
const emit = defineEmits(['create', 'edit', 'delete', 'reset-pwd', 'filter-role', 'export', 'bulk-delete', 'bulk-reset', 'import'])

const filter = ref('')
const roleFilter = ref('all')
const selected = ref([])

const columns = [
    { name: 'name', label: 'Nome Completo', field: row => `${row.last_name} ${row.first_name}`, sortable: true, align: 'left' },
    { name: 'email', label: 'Email', field: 'email', sortable: true, align: 'left' },
    { name: 'role', label: 'Ruolo', field: 'role', sortable: true, align: 'center' },
    { name: 'class', label: 'Classe', field: row => row.class || '-', align: 'center' },
    { name: 'status', label: 'Stato', field: 'active', align: 'center' },
    { name: 'actions', label: 'Azioni', align: 'right' }
];

const getRoleColor = (role) => {
    switch(role) {
        case 'student': return 'green'
        case 'teacher': return 'purple'
        case 'parent': return 'orange'
        case 'staff': return 'blue-grey'
        default: return 'grey'
    }
}

const getRoleLabel = (role) => {
    switch(role) {
        case 'student': return 'Studente'
        case 'teacher': return 'Docente'
        case 'parent': return 'Genitore'
        case 'staff': return 'Personale'
        default: return role
    }
}
</script>
