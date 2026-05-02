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
        @reset-pwd="openResetPwd"
        @filter-role="currentRoleFilter = $event"
        @export="exportUsers"
        @bulk-delete="bulkDelete"
        @manage-subjects="openManageSubjects"
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

    <!-- Teacher Subjects Dialog -->
    <q-dialog v-model="showSubjectsDialog" full-width>
        <q-card>
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Materie Docente - {{ currentTeacherName }}</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section>
                <div class="row q-col-gutter-md">
                    <div class="col-md-7 col-12">
                        <q-list bordered separator>
                            <q-item-label header>Materie Abilitate</q-item-label>
                            <q-item v-if="teacherSubjects.length === 0">
                                <q-item-section class="text-grey italic">Nessuna materia assegnata</q-item-section>
                            </q-item>
                            <q-item v-for="ts in teacherSubjects" :key="ts.id">
                                <q-item-section>
                                    <q-item-label>{{ ts.subject_name }}</q-item-label>
                                </q-item-section>
                                <q-item-section side>
                                    <q-btn flat round icon="delete" color="negative" @click="removeTeacherSubject(ts.subject_id)" />
                                </q-item-section>
                            </q-item>
                        </q-list>
                    </div>
                    <div class="col-md-5 col-12">
                         <q-card flat bordered class="q-pa-md">
                            <div class="text-subtitle2 q-mb-sm">Aggiungi Abilitazione</div>
                            <q-select
                                v-model="selectedSubjectToAdd"
                                :options="availableSubjects"
                                label="Materia"
                                outlined
                                dense
                                emit-value
                                map-options
                            />
                            <q-btn label="Aggiungi" color="primary" class="full-width q-mt-md" @click="addTeacherSubject" :disable="!selectedSubjectToAdd" />
                         </q-card>
                    </div>
                </div>
            </q-card-section>
        </q-card>
    </q-dialog>

    <!-- Reset Password Dialog -->
    <q-dialog v-model="showResetPwdDialog">
        <q-card style="min-width: 350px">
            <q-card-section>
                <div class="text-h6">Reset Password</div>
                <div class="text-subtitle2">{{ resetTargetName }}</div>
            </q-card-section>
            <q-card-section>
                <q-input v-model="newPassword" label="Nuova Password" type="password" outlined dense />
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Reset" @click="handleResetPwd" :disable="newPassword.length < 6" />
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

const showResetPwdDialog = ref(false);
const resetTargetId = ref(null);
const resetTargetName = ref('');
const newPassword = ref('');

const openResetPwd = (user) => {
    resetTargetId.value = user.id;
    resetTargetName.value = `${user.first_name} ${user.last_name}`;
    newPassword.value = '';
    showResetPwdDialog.value = true;
};

const handleResetPwd = async () => {
    try {
        await userService.forceResetPassword(resetTargetId.value, newPassword.value);
        $q.notify({ type: 'positive', message: 'Password aggiornata con successo' });
        showResetPwdDialog.value = false;
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore durante il reset della password' });
    }
};

const bulkDelete = (selected) => {
    $q.dialog({
        title: 'Eliminazione Massiva',
        message: `Sei sicuro di voler eliminare ${selected.length} utenti?`,
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            const ids = selected.map(u => u.id);
            await userService.bulkDelete(ids);
            $q.notify({ type: 'positive', message: `${selected.length} utenti eliminati` });
            fetchUsers();
        } catch(e) {
            $q.notify({ type: 'negative', message: 'Errore eliminazione massiva' });
        }
    });
};

const handleImport = async () => {
    if (!importFile.value) return;
    try {
        const res = await userService.bulkImport(importFile.value);
        $q.notify({ type: 'positive', message: `Importati ${res.data.created} utenti su ${res.data.total}` });
        showImport.value = false;
        importFile.value = null;
        fetchUsers();
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore importazione file' });
    }
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

// Teacher Subjects Logic
const showSubjectsDialog = ref(false)
const teacherSubjects = ref([])
const availableSubjects = ref([])
const selectedSubjectToAdd = ref(null)
const currentTeacherId = ref(null)
const currentTeacherName = ref('')
const teachersCache = ref([]) // Cache teachers list to map user_id -> teacher_id

const openManageSubjects = async (user) => {
    currentTeacherName.value = `${user.last_name} ${user.first_name}`
    currentTeacherId.value = null
    teacherSubjects.value = []
    selectedSubjectToAdd.value = null
    showSubjectsDialog.value = true
    
    // 1. Find Teacher ID
    try {
        if (teachersCache.value.length === 0) {
            const tRes = await adminService.getTeachersList(authStore.user.school_id)
            teachersCache.value = tRes.data || []
        }
        const teacher = teachersCache.value.find(t => t.user_id === user.id)
        if (!teacher) {
            $q.notify({ type: 'warning', message: 'Profilo docente non trovato. Assicurati che sia stato creato.' })
            showSubjectsDialog.value = false
            return
        }
        currentTeacherId.value = teacher.id
        
        // 2. Load Subjects (Assigned & Available)
        await loadTeacherSubjects()
        if (availableSubjects.value.length === 0) {
            const sRes = await adminService.getSubjects(authStore.user.school_id)
            availableSubjects.value = (sRes.data || []).map(s => ({ label: s.name, value: s.id }))
        }

    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento dati docente' })
    }
}

const loadTeacherSubjects = async () => {
    if (!currentTeacherId.value) return
    const res = await adminService.getTeacherSubjects(currentTeacherId.value)
    teacherSubjects.value = res.data || []
}

const addTeacherSubject = async () => {
    if (!selectedSubjectToAdd.value || !currentTeacherId.value) return
    try {
        await adminService.assignSubjectToTeacher(currentTeacherId.value, selectedSubjectToAdd.value)
        $q.notify({ type: 'positive', message: 'Materia aggiunta' })
        loadTeacherSubjects()
        selectedSubjectToAdd.value = null
    } catch(e) {
         $q.notify({ type: 'negative', message: 'Errore assegnazione materia' })
    }
}

const removeTeacherSubject = async (subjectId) => {
    try {
        await adminService.removeTeacherSubject(currentTeacherId.value, subjectId)
        loadTeacherSubjects()
        $q.notify({ type: 'positive', message: 'Materia rimossa' })
    } catch(e) {
         $q.notify({ type: 'negative', message: 'Errore rimozione materia' })
    }
}
</script>
