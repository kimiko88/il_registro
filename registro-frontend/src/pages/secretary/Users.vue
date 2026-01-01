<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4 text-weight-bold">Gestione Utenti</div>
       <div class="q-gutter-sm">
           <q-btn label="Nuova Classe" color="primary" icon="add" outline @click="openClassDialog" />
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
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">{{ isEditing ? 'Modifica Utente' : 'Nuovo Utente' }}</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>
            
            <q-card-section>
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
                    <q-input v-model="userForm.fiscal_code" label="Codice Fiscale" outlined dense maxlength="16" />
                    
                    <q-select 
                        v-model="userForm.role" 
                        :options="roleOptions"
                        label="Ruolo"
                        outlined
                        dense
                        emit-value
                        map-options
                    />

                    <div v-if="userForm.role === 'student'">
                         <q-select
                            v-model="userForm.class_id"
                            :options="classOptions"
                            label="Classe"
                            outlined
                            dense
                            emit-value
                            map-options
                            hint="Seleziona la classe di appartenenza"
                            :loading="loadingClasses"
                         />
                    </div>
                    
                     <q-input
                         v-if="!isEditing"
                         v-model="userForm.password"
                         label="Password Provvisoria"
                         outlined dense
                         type="password"
                         :rules="[val => !!val || 'Campo obbligatorio', val => val.length >= 8 || 'Minimo 8 caratteri']"
                    />

                    <div class="row justify-end q-mt-lg">
                        <q-btn label="Annulla" flat v-close-popup color="grey" />
                        <q-btn :label="isEditing ? 'Salva' : 'Crea'" type="submit" color="primary" class="q-ml-sm" />
                    </div>
                </q-form>
            </q-card-section>
        </q-card>
    </q-dialog>

    <!-- Create Class Dialog -->
    <q-dialog v-model="showClassDialog">
        <q-card style="min-width: 400px">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Nuova Classe</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section>
                 <q-form @submit="saveClass" class="q-gutter-md">
                    <q-input v-model="classForm.name" label="Nome (es. 1A)" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
                    <q-input v-model="classForm.academic_year" label="Anno Scolastico" outlined dense />
                    
                    <div class="row justify-end">
                        <q-btn label="Annulla" flat v-close-popup color="grey" />
                        <q-btn label="Crea Classe" type="submit" color="primary" />
                    </div>
                 </q-form>
            </q-card-section>
        </q-card>
    </q-dialog>

    <!-- Import Dialog -->
    <q-dialog v-model="showImport">
        <q-card style="min-width: 400px">
             <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Importazione Massiva</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>
            <q-card-section>
                <div class="text-caption text-grey q-mb-md">Carica un file CSV o XLSX con i dati degli utenti.</div>
                <q-file v-model="importFile" label="Seleziona File" outlined dense accept=".csv, .xlsx">
                     <template v-slot:prepend>
                        <q-icon name="attach_file" />
                     </template>
                </q-file>
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
import { ref, computed, onMounted, reactive, watch } from 'vue';
import { useQuasar, exportFile } from 'quasar';
import UserTable from 'src/components/Secretary/UserTable.vue';
import { userService } from 'src/services/userService';
import adminService from 'src/services/adminService';
import { useAuthStore } from 'src/stores/auth';

const $q = useQuasar();
const authStore = useAuthStore();
const loading = ref(false);
const showUserDialog = ref(false);
const showClassDialog = ref(false);
const showImport = ref(false);
const isEditing = ref(false);
const importFile = ref(null);
const currentRoleFilter = ref('all');

const users = ref([]); 

// Classes
const classes = ref([]);
const loadingClasses = ref(false);
const classOptions = computed(() => {
    return classes.value.map(c => ({
        label: c.name || `${c.section} ${c.academic_year}`, // Fallback if name optional
        value: c.id
    }))
});

// Forms
const userForm = reactive({
    id: null,
    first_name: '',
    last_name: '',
    email: '',
    role: 'student',
    class_id: null,
    fiscal_code: '',
    password: ''
});

const classForm = reactive({
    name: '',
    section: '',
    academic_year: '2024/2025'
});


onMounted(() => {
    fetchUsers()
    fetchClasses()
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

const fetchClasses = async () => {
    if (!authStore.user?.school_id) return;
    loadingClasses.value = true;
    try {
        const res = await adminService.getSchoolClasses(authStore.user.school_id);
        classes.value = res.data;
    } catch (e) {
        console.error("Error loading classes", e);
    } finally {
        loadingClasses.value = false;
    }
}

// Watch filter change
watch(currentRoleFilter, () => {
    fetchUsers()
})

const filteredUsers = computed(() => users.value);

const roleOptions = [
    { label: 'Studente', value: 'student' },
    { label: 'Docente', value: 'teacher' },
    { label: 'Genitore', value: 'parent' },
    { label: 'Personale ATA', value: 'staff' }
];

const openCreate = () => {
    isEditing.value = false;
    // Reset form
    userForm.id = null;
    userForm.first_name = '';
    userForm.last_name = '';
    userForm.email = '';
    userForm.role = 'student';
    userForm.class_id = null;
    userForm.fiscal_code = '';
    userForm.password = '';
    
    showUserDialog.value = true;
};

const openEdit = (user) => {
    isEditing.value = true;
    Object.assign(userForm, user);
    // Explicitly set class_id if missing (though should come from API now)
    if (user.ClassID) userForm.class_id = user.ClassID; // Ensure casing matches API DTO
    // Note: API returns snake_case usually, check naming. 
    // Go struct: ClassID `json:"class_id"` -> response has class_id.
    // user object from API should have class_id.
    // The table might be flattening it.
    
    showUserDialog.value = true;
};

const openClassDialog = () => {
    classForm.name = '';
    classForm.section = '';
    classForm.academic_year = '2024/2025';
    showClassDialog.value = true;
};

const saveUser = async () => {
    try {
        const payload = { ...userForm, school_id: authStore.user.school_id };
        // If editing, don't send empty password
        if (isEditing.value) {
            delete payload.password; 
            await userService.update(userForm.id, payload);
             $q.notify({ type: 'positive', message: 'Utente aggiornato' });
        } else {
             await userService.create(payload);
             $q.notify({ type: 'positive', message: 'Utente creato' });
        }
        showUserDialog.value = false;
        fetchUsers();
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore salvataggio utente', caption: e.message });
    }
};

const saveClass = async () => {
    try {
        // Need to split name into section if needed or just send name
        // The API CreateClassRequest expects Name, AcademicYear, SchoolID
        // We'll map name to name (e.g. "1A") and maybe section to "A"?
        // Simpler: Just send name.
        const payload = {
            name: classForm.name,
            academic_year: classForm.academic_year,
            school_id: authStore.user.school_id,
            section: classForm.name.replace(/[0-9]/g, '') // Rough guess
        };
        await adminService.createClass(payload);
        $q.notify({ type: 'positive', message: 'Classe creata con successo' });
        showClassDialog.value = false;
        fetchClasses();
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore creazione classe', caption: e.message });
    }
};

const confirmDelete = (user) => {
    $q.dialog({
        title: 'Conferma eliminazione',
        message: `Vuoi davvero eliminare ${user.first_name} ${user.last_name}?`,
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            await userService.delete(user.id);
            $q.notify({ type: 'positive', message: 'Utente eliminato' });
            fetchUsers();
        } catch(e) {
            $q.notify({ type: 'negative', message: 'Errore eliminazione' });
        }
    });
};

const confirmResetPwd = (user) => {
    $q.dialog({
        title: 'Reset Password',
        message: `Inviare link di reset password a ${user.email}?`,
        cancel: true
    }).onOk(async () => {
         try {
            await userService.resetPassword(user.id);
            $q.notify({ type: 'positive', message: 'Link inviato con successo' });
        } catch(e) {
             $q.notify({ type: 'negative', message: 'Errore reset password' });
        }
    });
};

const bulkDelete = (selected) => {
     // TODO: Implement bulk delete API
     $q.notify({ type: 'warning', message: 'Funzionalità non ancora implementata nel backend' });
};

const handleImport = async () => {
    // TODO: Implement real import via service
     $q.notify({ type: 'warning', message: 'Funzionalità mock per demo' });
     showImport.value = false;
};

const exportUsers = () => {
    const content = [
        'ID,Nome,Cognome,Email,Ruolo,Classe',
        ...filteredUsers.value.map(u => `${u.id},${u.first_name},${u.last_name},${u.email},${u.role},${u.class_name || ''}`)
    ].join('\r\n');

    const status = exportFile(
        'users-export.csv',
        content,
        'text/csv'
    );
    if (!status) $q.notify({ type: 'negative', message: 'Export fallito' });
};
</script>
