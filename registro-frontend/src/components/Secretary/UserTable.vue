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
      virtual-scroll
      :virtual-scroll-item-size="48"
      class="bg-transparent"
    >
      <template v-slot:top>
        <div class="row items-center full-width q-mb-md">
          <div class="text-h5 text-weight-bold text-outfit q-mr-xl text-slate-800">{{ t('usersPage.title') || 'Elenco Utenti' }}</div>
          
          <q-select
            for="user-table-role-filter"
            aria-label="Filtra per Ruolo"
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
            <q-input for="user-table-search" aria-label="Cerca utenti" dense outlined v-model="filter" :placeholder="t('usersPage.searchPlaceholder') || 'Cerca per nome, email...'" class="bg-white min-width-250">
              <template v-slot:prepend>
                <q-icon name="search" color="slate-300" />
              </template>
            </q-input>
            
            <q-btn unelevated color="primary" icon="add" :label="t('usersPage.newUser') || 'Nuovo Utente'" class="rounded-lg shadow-sm" no-caps @click="$emit('create')" />
            <q-btn flat round icon="file_download" color="slate-400" :aria-label="t('usersPage.exportCsv') || 'Esporta in CSV'" @click="$emit('export')">
                <q-tooltip>{{ t('usersPage.exportCsv') || 'Esporta in CSV' }}</q-tooltip>
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
               <span class="text-weight-bold text-indigo-700">{{ selected.length }} {{ t('usersPage.selectedUsers') || 'utenti selezionati' }}</span>
               <q-space />
               <div class="row q-gutter-sm">
                 <q-btn unelevated size="sm" color="negative" icon="delete" :label="t('usersPage.deleteSelected') || 'Elimina Selezionati'" no-caps class="rounded-md" @click="$emit('bulk-delete', selected)" />
                 <q-btn outline size="sm" color="indigo" icon="lock_reset" :label="t('usersPage.resetPassword') || 'Reset Password'" no-caps class="rounded-md" @click="$emit('bulk-reset', selected)" />
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
          <div class="row items-center justify-center q-gutter-xs">
            <q-chip :color="getRoleColor(props.value) + '-50'" :text-color="getRoleColor(props.value) + '-700'" size="sm" class="text-weight-bold rounded-md">
              {{ getRoleLabel(props.value) }}
            </q-chip>
            <q-chip v-if="props.row?.is_staff" color="amber-1" text-color="amber-10" icon="shield" size="xs" class="text-weight-bold rounded-md">
              {{ t('usersPage.roleStaff') || 'Staff' }}
            </q-chip>
            <q-chip v-if="props.row?.is_vice_principal" color="purple-1" text-color="purple-10" icon="stars" size="xs" class="text-weight-bold rounded-md">
              {{ t('usersPage.roleVicePrincipal') || 'Vicepreside' }}
            </q-chip>
            <q-chip v-if="props.row?.is_principal" color="deep-purple-1" text-color="deep-purple-10" icon="workspace_premium" size="xs" class="text-weight-bold rounded-md">
              {{ t('usersPage.rolePrincipal') || 'Preside' }}
            </q-chip>
            <template v-if="props.row?.assignments?.length">
              <q-chip 
                v-for="a in (props.row.assignments || []).filter(x => x.is_active !== false)" 
                :key="a.id || a.assignment_type"
                color="indigo-1"
                text-color="indigo-9"
                size="xs"
                class="text-weight-medium rounded-md"
                icon="badge"
              >
                {{ a.title || a.assignment_type }}
              </q-chip>
            </template>
          </div>
        </q-td>
      </template>

      <template v-slot:body-cell-status="props">
         <q-td :props="props">
            <div class="row items-center q-gutter-xs">
              <div :class="props.value ? 'bg-emerald-500' : 'bg-slate-300'" class="status-dot"></div>
              <span :class="props.value ? 'text-emerald-700 text-weight-medium' : 'text-slate-400'">
                {{ props.value ? (t('common.active') || 'Attivo') : (t('common.inactive') || 'Inattivo') }}
              </span>
            </div>
         </q-td>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props" auto-width>
          <q-btn flat round size="sm" color="slate-400" icon="more_horiz" :aria-label="t('common.actions') || 'Azioni'">
            <q-menu class="rounded-lg shadow-2xl border border-slate-100" transition-show="fade" transition-hide="fade">
              <q-list style="min-width: 180px" padding>
                <q-item clickable v-close-popup class="q-mx-sm rounded-md" @click="$emit('edit', props.row)">
                  <q-item-section avatar><q-icon name="edit" color="primary" /></q-item-section>
                  <q-item-section class="text-slate-700">{{ t('usersPage.editProfile') || 'Modifica Profilo' }}</q-item-section>
                </q-item>
                
                <q-item v-if="canManageUserAssignments(props.row)" clickable v-close-popup class="q-mx-sm rounded-md" @click="$emit('manage-assignments', props.row)">
                  <q-item-section avatar><q-icon name="assignment_ind" color="secondary" /></q-item-section>
                  <q-item-section class="text-slate-700">Incarichi & Coordinamento</q-item-section>
                </q-item>

                <q-item clickable v-close-popup class="q-mx-sm rounded-md" @click="$emit('reset-pwd', props.row)">
                  <q-item-section avatar><q-icon name="lock_reset" color="orange" /></q-item-section>
                  <q-item-section class="text-slate-700">{{ t('usersPage.resetPassword') || 'Reset Password' }}</q-item-section>
                </q-item>
                
                <q-item v-if="props.row.role === 'teacher'" clickable v-close-popup class="q-mx-sm rounded-md" @click="$emit('manage-subjects', props.row)">
                  <q-item-section avatar><q-icon name="menu_book" color="indigo" /></q-item-section>
                  <q-item-section class="text-slate-700">{{ t('usersPage.manageSubjects') || 'Gestione Materie' }}</q-item-section>
                </q-item>
                
                <q-separator class="q-my-sm opacity-50" />
                <q-item clickable v-close-popup class="q-mx-sm rounded-md text-negative" @click="$emit('delete', props.row)">
                  <q-item-section avatar><q-icon name="delete" color="negative" /></q-item-section>
                  <q-item-section class="text-weight-bold">{{ t('usersPage.deleteAccount') || 'Elimina Account' }}</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-td>
      </template>
      
      <template v-slot:no-data>
        <div class="full-width q-pa-xl text-center text-slate-400">
          <q-icon name="group_off" size="64px" class="opacity-10 q-mb-md" />
          <div class="text-h6">{{ t('usersPage.noUsersFound') || 'Nessun utente trovato' }}</div>
        </div>
      </template>
    </q-table>
  </q-card>
</template>
<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
// eslint-disable-next-line no-unused-vars
const props = defineProps(['users', 'loading']);
// eslint-disable-next-line no-unused-vars
const emit = defineEmits(['create', 'edit', 'delete', 'reset-pwd', 'filter-role', 'export', 'bulk-delete', 'bulk-reset', 'import', 'manage-subjects', 'manage-assignments'])

const filter = ref('')
const roleFilter = ref('all')
const selected = ref([])

const roleFilterOptions = computed(() => [
  { label: t('common.all') || 'Tutti', value: 'all' },
  { label: t('usersPage.roleStudents') || 'Studenti', value: 'student' },
  { label: t('usersPage.roleTeachers') || 'Docenti', value: 'teacher' },
  { label: t('usersPage.roleParents') || 'Genitori', value: 'parent' },
  { label: 'Assistenti Amministrativi', value: 'assistente_amministrativo' },
  { label: 'Assistenti Alunni', value: 'assistente_alunni' },
  { label: 'Assistenti Personale', value: 'assistente_personale' },
  { label: 'Assistenti Contabilità', value: 'assistente_contabilita' },
  { label: 'Assistenti Protocollo', value: 'assistente_protocollo' },
  { label: 'Assistenti Sportello', value: 'assistente_sportello' },
  { label: 'Assistenti Tecnici', value: 'assistente_tecnico' },
  { label: 'Collaboratori Scolastici', value: 'collaboratore_scolastico' },
  { label: 'Collaboratori D.S.', value: 'collaboratore_ds' },
  { label: 'Responsabili di Servizio', value: 'responsabile_servizio' },
  { label: 'DSGA', value: 'dsga' },
  { label: t('usersPage.roleSecretary') || 'Segreteria', value: 'secretary' },
  { label: t('usersPage.rolePrincipal') || 'Preside', value: 'principal' },
  { label: t('usersPage.roleVicePrincipal') || 'Vicepreside', value: 'vice_principal' },
  { label: 'Resp. Gestione Documentale', value: 'responsabile_gestione_documentale' },
  { label: 'Resp. Conservazione', value: 'responsabile_conservazione' },
  { label: 'DPO', value: 'dpo' },
  { label: t('usersPage.roleAdmin') || 'Amministratore', value: 'admin' },
  { label: 'Super Admin', value: 'superadmin' },
  { label: t('usersPage.roleStaff') || 'Staff', value: 'staff' }
])

const columns = computed(() => [
    { name: 'name', label: t('common.fullName') || 'Nome Completo', field: row => `${row?.last_name || ''} ${row?.first_name || ''}`.trim() || '-', sortable: true, align: 'left' },
    { name: 'email', label: t('login.emailLabel') || 'Email', field: 'email', sortable: true, align: 'left' },
    { name: 'role', label: t('usersPage.roleLabel') || 'Ruolo', field: 'role', sortable: true, align: 'center' },
    { name: 'class', label: t('udaPage.classLabel') || 'Classe', field: row => row.class_name || row.class || '-', align: 'center' },
    { name: 'status', label: t('substitutionsPage.status') || 'Stato', field: row => row.is_active !== undefined ? row.is_active : row.active, align: 'center' },
    { name: 'actions', label: t('common.actions') || 'Azioni', align: 'right' }
]);

const canManageUserAssignments = (row) => {
    if (!row) return false
    // Teachers and staff/ATA users can have duty assignments
    return row.role === 'teacher' || 
           row.role === 'collaboratore_scolastico' || 
           row.role === 'assistente_amministrativo' || 
           row.role === 'assistente_tecnico' ||
           row.role === 'collaboratore_ds' ||
           row.role === 'responsabile_servizio'
}

const getRoleColor = (role) => {
    switch(role) {
        case 'student': return 'green'
        case 'teacher': return 'purple'
        case 'parent': return 'orange'
        case 'staff': return 'blue-grey'
        case 'secretary': return 'cyan'
        case 'principal': return 'deep-purple'
        case 'vice_principal': return 'indigo'
        case 'dsga': return 'teal'
        case 'assistente_amministrativo':
        case 'assistente_alunni':
        case 'assistente_personale':
        case 'assistente_contabilita':
        case 'assistente_protocollo':
        case 'assistente_sportello':
            return 'cyan'
        case 'assistente_tecnico': return 'blue'
        case 'collaboratore_ds': return 'deep-orange'
        case 'collaboratore_scolastico': return 'amber-9'
        case 'responsabile_servizio': return 'teal'
        case 'responsabile_gestione_documentale': return 'blue-grey'
        case 'responsabile_conservazione': return 'blue-grey'
        case 'dpo': return 'indigo'
        case 'admin': return 'red'
        case 'superadmin': return 'amber'
        default: return 'grey'
    }
}

const getRoleLabel = (role) => {
    switch(role) {
        case 'student': return t('usersPage.roleStudents') || 'Studente'
        case 'teacher': return t('usersPage.roleTeachers') || 'Docente'
        case 'parent': return t('usersPage.roleParents') || 'Genitore'
        case 'staff': return t('usersPage.roleStaff') || 'Personale'
        case 'secretary': return t('usersPage.roleSecretary') || 'Segreteria'
        case 'principal': return t('usersPage.rolePrincipal') || 'Preside'
        case 'vice_principal': return t('usersPage.roleVicePrincipal') || 'Vicepreside'
        case 'dsga': return 'DSGA'
        case 'assistente_amministrativo': return 'Assistente Amministrativo'
        case 'assistente_alunni': return 'Assistente Alunni'
        case 'assistente_personale': return 'Assistente Personale'
        case 'assistente_contabilita': return 'Assistente Contabilità'
        case 'assistente_protocollo': return 'Assistente Protocollo'
        case 'assistente_sportello': return 'Assistente Sportello'
        case 'assistente_tecnico': return 'Assistente Tecnico'
        case 'collaboratore_ds': return 'Collaboratore D.S.'
        case 'collaboratore_scolastico': return 'Collaboratore Scolastico'
        case 'responsabile_servizio': return 'Responsabile Servizio'
        case 'responsabile_gestione_documentale': return 'Resp. Documentale'
        case 'responsabile_conservazione': return 'Resp. Conservazione'
        case 'dpo': return 'DPO'
        case 'admin': return t('usersPage.roleAdmin') || 'Amministratore'
        case 'superadmin': return 'Super Admin'
        default: return role
    }
}
</script>
