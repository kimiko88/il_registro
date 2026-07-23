<template>
    <q-card flat class="rounded-xl border border-slate-100 bg-white shadow-soft overflow-hidden">
    <q-table
      :rows="users"
      :columns="columns"
      :loading="loading"
      :filter="filter"
      row-key="id"
      selection="multiple"
      v-model:selected="selected"
      flat
      class="bg-transparent"
    >
      <template v-slot:top>
        <div class="row items-center full-width q-mb-md">
          <div class="text-h5 text-weight-bold text-outfit q-mr-xl text-slate-800">Elenco Utenti</div>
          
          <q-select
            v-model="roleFilter"
            :options="roleFilterOptions"
            dense
            outlined
            emit-value
            map-options
            bg-color="white"
            class="rounded-lg border border-slate-100 min-width-200"
            @update:model-value="$emit('filter-role', $event)"
          />

          <q-space />
          
          <div class="row q-gutter-sm">
            <q-input dense outlined v-model="filter" placeholder="Cerca per nome, email..." class="bg-white min-width-250">
              <template v-slot:prepend>
                <q-icon name="search" color="slate-300" />
              </template>
            </q-input>
            
            <q-btn unelevated color="primary" icon="add" label="Nuovo Utente" class="rounded-lg shadow-sm" no-caps @click="$emit('create')" />
            <q-btn flat round icon="file_download" color="slate-400" @click="$emit('export')">
                <q-tooltip>Esporta in CSV</q-tooltip>
            </q-btn>
          </div>
        </div>
      </template>

      <!-- Bulk Actions -->
      <template v-slot:top-row v-if="selected.length > 0">
         <q-tr class="bg-indigo-50 animate-fade-in">
           <q-td colspan="100%">
             <div class="row items-center q-gutter-md q-pa-sm">
               <q-icon name="check_circle" color="indigo" size="24px" />
               <span class="text-weight-bold text-indigo-700">{{ selected.length }} utenti selezionati</span>
               <q-space />
               <div class="row q-gutter-sm">
                 <q-btn unelevated size="sm" color="negative" icon="delete" label="Elimina Selezionati" no-caps class="rounded-md" @click="$emit('bulk-delete', selected)" />
                 <q-btn outline size="sm" color="indigo" icon="lock_reset" label="Reset Password" no-caps class="rounded-md" @click="$emit('bulk-reset', selected)" />
               </div>
             </div>
           </q-td>
         </q-tr>
      </template>

      <!-- Header -->
      <template v-slot:header-cell="props">
        <q-th :props="props" class="text-slate-500 font-bold">
          {{ props.col.label }}
        </q-th>
      </template>

      <!-- Custom Body -->
      <template v-slot:body-cell-role="props">
        <q-td :props="props">
          <q-chip :color="getRoleColor(props.value) + '-50'" :text-color="getRoleColor(props.value) + '-700'" size="sm" class="text-weight-bold rounded-md">
            {{ getRoleLabel(props.value) }}
          </q-chip>
        </q-td>
      </template>

      <template v-slot:body-cell-status="props">
         <q-td :props="props">
            <div class="row items-center q-gutter-xs">
              <div :class="props.value ? 'bg-emerald-500' : 'bg-slate-300'" class="status-dot"></div>
              <span :class="props.value ? 'text-emerald-700 text-weight-medium' : 'text-slate-400'">
                {{ props.value ? 'Attivo' : 'Inattivo' }}
              </span>
            </div>
         </q-td>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props" auto-width>
          <q-btn flat round size="sm" color="slate-400" icon="more_horiz">
            <q-menu class="rounded-lg shadow-2xl border border-slate-100" transition-show="fade" transition-hide="fade">
              <q-list style="min-width: 180px" padding>
                <q-item clickable v-close-popup class="q-mx-sm rounded-md" @click="$emit('edit', props.row)">
                  <q-item-section avatar><q-icon name="edit" color="primary" /></q-item-section>
                  <q-item-section class="text-slate-700">Modifica Profilo</q-item-section>
                </q-item>
                <q-item clickable v-close-popup class="q-mx-sm rounded-md" @click="$emit('reset-pwd', props.row)">
                  <q-item-section avatar><q-icon name="lock_reset" color="orange" /></q-item-section>
                  <q-item-section class="text-slate-700">Reset Password</q-item-section>
                </q-item>
                
                <q-item v-if="props.row.role === 'teacher'" clickable v-close-popup class="q-mx-sm rounded-md" @click="$emit('manage-subjects', props.row)">
                  <q-item-section avatar><q-icon name="menu_book" color="indigo" /></q-item-section>
                  <q-item-section class="text-slate-700">Gestione Materie</q-item-section>
                </q-item>
                
                <q-separator class="q-my-sm opacity-50" />
                <q-item clickable v-close-popup class="q-mx-sm rounded-md text-negative" @click="$emit('delete', props.row)">
                  <q-item-section avatar><q-icon name="delete" color="negative" /></q-item-section>
                  <q-item-section class="text-weight-bold">Elimina Account</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-td>
      </template>
      
      <template v-slot:no-data>
        <div class="full-width q-pa-xl text-center text-slate-400">
          <q-icon name="group_off" size="64px" class="opacity-10 q-mb-md" />
          <div class="text-h6">Nessun utente trovato</div>
        </div>
      </template>
    </q-table>
  </q-card>
</template>
<script setup>
import { ref } from 'vue'

// eslint-disable-next-line no-unused-vars
const props = defineProps(['users', 'loading']);
// eslint-disable-next-line no-unused-vars
const emit = defineEmits(['create', 'edit', 'delete', 'reset-pwd', 'filter-role', 'export', 'bulk-delete', 'bulk-reset', 'import', 'manage-subjects'])

const filter = ref('')
const roleFilter = ref('all')
const selected = ref([])

const roleFilterOptions = [
  {label: 'Tutti', value: 'all'},
  {label: 'Studenti', value: 'student'},
  {label: 'Docenti', value: 'teacher'},
  {label: 'Genitori', value: 'parent'},
  {label: 'Staff', value: 'staff'}
]

const columns = [
    { name: 'name', label: 'Nome Completo', field: row => `${row.last_name} ${row.first_name}`, sortable: true, align: 'left' },
    { name: 'email', label: 'Email', field: 'email', sortable: true, align: 'left' },
    { name: 'role', label: 'Ruolo', field: 'role', sortable: true, align: 'center' },
    { name: 'class', label: 'Classe', field: row => row.class_name || row.class || '-', align: 'center' },
    { name: 'status', label: 'Stato', field: row => row.is_active !== undefined ? row.is_active : row.active, align: 'center' },
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
