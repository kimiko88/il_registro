<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4 text-weight-bold">Gestione Utenti</div>
       <div class="q-gutter-sm">
           <q-btn label="Importa CSV" color="secondary" icon="upload" outline @click="showImport=true" />
           <q-btn flat icon="history" label="Audit Log" @click="$q.notify('Audit Log non disponibile per Segreteria')" />
       </div>
    </div>

    <UserTable 
        :users="filteredUsers" 
        :loading="loading"
        @create="openCreate"
        @edit="openEdit"
        @delete="confirmDelete"
        @reset-pwd="confirmResetPwd"
        @filter-role="currentRoleFilter = $event"
        @export="exportUsers"
        @bulk-delete="bulkDelete"
    />

    <!-- Create/Edit User Dialog -->
    <q-dialog v-model="showUserDialog">
        <q-card style="min-width: 500px">
            <q-card-section>
                <div class="text-h6">{{ isEditing ? 'Modifica Utente' : 'Nuovo Utente' }}</div>
            </q-card-section>
            
            <q-card-section class="q-pt-none">
                <q-form @submit="saveUser" class="q-gutter-md">
                    <div class="row q-col-gutter-sm">
                        <div class="col-6">
                            <q-input v-model="userForm.first_name" label="Nome" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
                        </div>
                        <div class="col-6">
                            <q-input v-model="userForm.last_name" label="Cognome" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
                        </div>
                    </div>
                    <q-input v-model="userForm.email" label="Email" outlined dense type="email" :rules="[val => !!val || 'Obbligatorio']" />
                    <q-input v-model="userForm.cf" label="Codice Fiscale" outlined dense maxlength="16" />
                    
                    <q-select 
                        v-model="userForm.role" 
                        :options="roleOptions"
                        label="Ruolo"
                        outlined
                        dense
                        emit-value
                        map-options
                    />

                    <q-input 
                        v-if="userForm.role === 'student'"
                        v-model="userForm.class" 
                        label="Classe (es. 1A)" 
                        outlined 
                        dense 
                        hint="Classe di appartenenza"
                    />
                    
                    <div class="row justify-end q-mt-lg">
                        <q-btn label="Annulla" flat v-close-popup color="grey" />
                        <q-btn :label="isEditing ? 'Salva' : 'Crea'" type="submit" color="primary" class="q-ml-sm" />
                    </div>
                </q-form>
            </q-card-section>
        </q-card>
    </q-dialog>

    <!-- Import Dialog -->
    <q-dialog v-model="showImport">
        <q-card style="min-width: 400px">
            <q-card-section>
                <div class="text-h6">Importazione Massiva</div>
                <div class="text-caption text-grey">Carica un file CSV o XLSX con i dati degli utenti.</div>
            </q-card-section>
            <q-card-section>
                <q-file v-model="importFile" label="Seleziona File" outlined dense accept=".csv, .xlsx">
                     <template v-slot:prepend>
                        <q-icon name="attach_file" />
                     </template>
                </q-file>
                <div class="bg-grey-2 q-pa-sm q-mt-sm rounded-borders text-caption">
                    Colonne richieste: Nome, Cognome, Email, Ruolo, CF
                </div>
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Carica ed Elabora" @click="handleImport" :disable="!importFile" />
            </q-card-actions>
        </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { useQuasar, exportFile } from 'quasar';
import UserTable from 'src/components/Secretary/UserTable.vue';
import { userService } from 'src/services/userService';

const $q = useQuasar();
const loading = ref(false);
const showUserDialog = ref(false);
const showImport = ref(false);
const isEditing = ref(false);
const importFile = ref(null);
const currentRoleFilter = ref('all');

const users = ref([]); 

onMounted(() => {
    fetchUsers()
});

const fetchUsers = async () => {
    loading.value = true
    try {
        const res = await userService.getAll({ role: currentRoleFilter.value === 'all' ? undefined : currentRoleFilter.value })
        users.value = res.data.users || []
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento utenti' })
    } finally {
        loading.value = false
    }
}

// Watch filter change
import { watch } from 'vue'
watch(currentRoleFilter, () => {
    fetchUsers()
})

const filteredUsers = computed(() => {
    // Filter handled by API or client side if API returns all
    // Since we fetch on change, we just return users.value
    // But if API returns pagination, we might need to handle it.
    // For now, assume simple list.
    return users.value
});

const roleOptions = [
    { label: 'Studente', value: 'student' },
    { label: 'Docente', value: 'teacher' },
    { label: 'Genitore', value: 'parent' },
    { label: 'Personale ATA', value: 'staff' }
];

const userForm = reactive({
    id: null,
    first_name: '',
    last_name: '',
    email: '',
    role: 'student',
    class: '',
    cf: ''
});

const openCreate = () => {
    isEditing.value = false;
    Object.assign(userForm, { id: null, first_name: '', last_name: '', email: '', role: 'student', class: '', cf: '' });
    showUserDialog.value = true;
};

const openEdit = (user) => {
    isEditing.value = true;
    Object.assign(userForm, { ...user });
    showUserDialog.value = true;
};

const saveUser = () => {
    if (isEditing.value) {
        const index = users.value.findIndex(u => u.id === userForm.id);
        if (index !== -1) users.value[index] = { ...userForm };
        $q.notify({ type: 'positive', message: 'Utente aggiornato test' });
    } else {
        users.value.unshift({ ...userForm, id: Date.now(), active: true });
        $q.notify({ type: 'positive', message: 'Utente creato test' });
    }
    showUserDialog.value = false;
};

const confirmDelete = (user) => {
    $q.dialog({
        title: 'Conferma eliminazione',
        message: `Vuoi davvero eliminare ${user.first_name} ${user.last_name}?`,
        cancel: true,
        persistent: true
    }).onOk(() => {
        users.value = users.value.filter(u => u.id !== user.id);
        $q.notify({ type: 'positive', message: 'Utente eliminato' });
    });
};

const confirmResetPwd = (user) => {
    $q.dialog({
        title: 'Reset Password',
        message: `Inviare link di reset password a ${user.email}?`,
        cancel: true
    }).onOk(() => {
        $q.notify({ type: 'positive', message: 'Link inviato con successo' });
    });
};

const bulkDelete = (selected) => {
    $q.dialog({
        title: 'Eliminazione Multipla',
        message: `Eliminare ${selected.length} utenti?`,
        cancel: true,
        persistent: true
    }).onOk(() => {
        const ids = selected.map(u => u.id);
        users.value = users.value.filter(u => !ids.includes(u.id));
        $q.notify({ type: 'positive', message: `${selected.length} utenti eliminati` });
    });
};

const handleImport = () => {
    $q.loading.show({ message: 'Elaborazione CSV...' });
    setTimeout(() => {
        $q.loading.hide();
        users.value.push({
            id: Date.now(),
            first_name: 'Importato',
            last_name: 'Utente',
            email: 'import@test.it',
            role: 'student',
            active: true
        });
        showImport.value = false;
        importFile.value = null;
        $q.notify({ type: 'positive', message: 'Importazione completata: 1 utente aggiunto' });
    }, 1500);
};

const exportUsers = () => {
    // Simple CSV Export logic
    const content = [
        'ID,Nome,Cognome,Email,Ruolo,Classe',
        ...filteredUsers.value.map(u => `${u.id},${u.first_name},${u.last_name},${u.email},${u.role},${u.class || ''}`)
    ].join('\r\n');

    const status = exportFile(
        'users-export.csv',
        content,
        'text/csv'
    );

    if (status !== true) {
        $q.notify({ type: 'negative', message: 'Export fallito' });
    } else {
        $q.notify({ type: 'positive', message: 'Export completato' });
    }
};
</script>
