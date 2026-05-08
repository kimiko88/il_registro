<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-xl">
       <div>
         <h1 class="text-h4 text-weight-bold text-outfit q-my-none text-gradient-premium">
           Gestione Utenti
         </h1>
         <p class="text-subtitle1 text-slate-500 q-mt-sm q-mb-none">Amministrazione profili studenti, docenti e staff</p>
       </div>
       <div class="row q-gutter-sm items-center">
           <q-select
               v-if="isSuperAdmin"
               v-model="filterSchoolId"
               :options="schoolOptions"
               label="Filtra per Scuola"
               option-label="name"
               option-value="id"
               emit-value
               map-options
               dense
               outlined
               clearable
               bg-color="white"
               style="min-width: 250px"
           />
           <q-btn unelevated label="Nuova Classe" color="primary" icon="add" class="rounded-lg shadow-sm" no-caps @click="openClassDialog" />
           <q-btn outline label="Importa CSV" color="primary" icon="upload" class="rounded-lg" no-caps @click="showImport=true" />
           <q-btn flat icon="history" label="Log Attività" class="rounded-lg text-slate-400" no-caps @click="$q.notify('Audit Log non disponibile per Segreteria')" />
       </div>
    </div>

    <UserTable 
        :users="users" 
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
    <q-dialog v-model="showUserDialog" class="premium-dialog">
        <q-card style="min-width: 600px" class="rounded-xl overflow-hidden shadow-24 bg-white">
            <q-card-section class="bg-gradient-primary text-white q-pa-lg row items-center justify-between">
                <div class="text-h5 text-weight-bold">{{ isEditing ? 'Modifica Profilo' : 'Crea Nuovo Profilo' }}</div>
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>
            
            <q-card-section class="q-pa-xl">
                <q-form @submit="saveUser" class="q-gutter-y-lg">
                    <div class="row q-col-gutter-lg">
                        <div class="col-6">
                            <q-input v-model="userForm.first_name" label="Nome" outlined :rules="[val => !!val || 'Campo richiesto']" />
                        </div>
                        <div class="col-6">
                            <q-input v-model="userForm.last_name" label="Cognome" outlined :rules="[val => !!val || 'Campo richiesto']" />
                        </div>
                    </div>
                    <q-input v-model="userForm.email" label="Email Istituzionale" outlined type="email" :rules="[val => !!val || 'Inserire un email valida']" />
                    <q-input v-model="userForm.fiscal_code" label="Codice Fiscale" outlined maxlength="16" class="uppercase-input" />
                    
                    <div class="row q-col-gutter-lg">
                        <div :class="isSuperAdmin ? 'col-6' : 'col-12'">
                            <q-select 
                                v-model="userForm.role" 
                                :options="roleOptions"
                                label="Ruolo"
                                outlined
                                emit-value
                                map-options
                            />
                        </div>
                        <div v-if="isSuperAdmin" class="col-6">
                            <q-select 
                                v-model="userForm.school_id" 
                                :options="schoolOptions"
                                label="Scuola"
                                outlined
                                emit-value
                                map-options
                                option-label="name"
                                option-value="id"
                                :rules="[val => !!val || 'Selezionare una scuola']"
                                @update:model-value="fetchClassesForSchool"
                            />
                        </div>
                    </div>

                    <div v-if="userForm.role === 'student'" class="bg-indigo-50 q-pa-lg rounded-xl border border-indigo-100">
                         <div class="text-subtitle2 text-indigo-700 q-mb-md">Dettagli Studente</div>
                         <q-select
                            v-model="userForm.class_id"
                            :options="classOptions"
                            label="Classe di appartenenza"
                            outlined
                            emit-value
                            map-options
                            hint="Lo studente verrà inserito automaticamente nel registro di questa classe"
                            :loading="loadingClasses"
                            bg-color="white"
                         />
                    </div>
                    
                     <q-input
                         v-if="!isEditing"
                         v-model="userForm.password"
                         label="Password Iniziale"
                         outlined
                         type="password"
                         :rules="[val => !!val || 'Campo obbligatorio', val => val.length >= 8 || 'La password deve contenere almeno 8 caratteri']"
                    />

                    <div class="row justify-end q-mt-xl q-gutter-sm">
                        <q-btn label="Annulla" flat v-close-popup color="slate-400" no-caps />
                        <q-btn :label="isEditing ? 'Salva Modifiche' : 'Crea Profilo'" type="submit" color="primary" class="q-px-xl rounded-lg shadow-sm" no-caps />
                    </div>
                </q-form>
            </q-card-section>
        </q-card>
    </q-dialog>

    <!-- Create Class Dialog -->
    <q-dialog v-model="showClassDialog" class="premium-dialog">
        <q-card style="min-width: 450px" class="rounded-xl overflow-hidden shadow-24 bg-white">
            <q-card-section class="bg-gradient-primary text-white q-pa-lg row items-center justify-between">
                <div class="text-h5 text-weight-bold">Nuova Classe</div>
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>
            <q-card-section class="q-pa-xl">
                 <q-form @submit="saveClass" class="q-gutter-y-lg">
                    <q-input v-model="classForm.name" label="Nome Classe (es. 1A, 5B)" outlined :rules="[val => !!val || 'Obbligatorio']" />
                    <q-input v-model="classForm.academic_year" label="Anno Scolastico" outlined :placeholder="currentYearStr" />
                    
                    <div class="row justify-end q-mt-xl q-gutter-sm">
                        <q-btn label="Annulla" flat v-close-popup color="slate-400" no-caps />
                        <q-btn label="Crea Classe" type="submit" color="primary" class="q-px-xl rounded-lg shadow-sm" no-caps />
                    </div>
                 </q-form>
            </q-card-section>
        </q-card>
    </q-dialog>

    <!-- Import Dialog -->
    <q-dialog v-model="showImport" class="premium-dialog">
        <q-card style="min-width: 450px" class="rounded-xl overflow-hidden shadow-24 bg-white">
             <q-card-section class="bg-gradient-primary text-white q-pa-lg row items-center justify-between">
                <div class="text-h5 text-weight-bold">Importazione Massiva</div>
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>
            <q-card-section class="q-pa-xl">
                <div class="text-body1 text-slate-500 q-mb-lg">Seleziona un file CSV o Excel contenente l'elenco degli utenti da importare.</div>
                <q-file v-model="importFile" label="Scegli file..." outlined counter class="rounded-lg">
                     <template #prepend>
                        <q-icon name="cloud_upload" color="primary" />
                     </template>
                </q-file>
                
                <div class="row justify-end q-mt-xl q-gutter-sm">
                  <q-btn flat label="Annulla" color="slate-400" v-close-popup no-caps />
                  <q-btn color="primary" label="Avvia Importazione" class="q-px-xl rounded-lg shadow-sm" no-caps @click="handleImport" :disable="!importFile" />
                </div>
            </q-card-section>
        </q-card>
    </q-dialog>

    <!-- Teacher Subjects Dialog -->
    <q-dialog v-model="showSubjectsDialog" full-width>
        <q-card class="bg-slate-50 rounded-xl overflow-hidden shadow-24">
            <q-card-section class="bg-white border-b border-slate-100 q-pa-lg row items-center justify-between">
                <div class="text-h5 text-weight-bold text-slate-800">Materie Docente: <span class="text-primary">{{ currentTeacherName }}</span></div>
                <q-btn icon="close" flat round dense v-close-popup color="slate-400" />
            </q-card-section>

            <q-card-section class="q-pa-xl">
                <div class="row q-col-gutter-xl">
                    <div class="col-md-7 col-12">
                        <q-card flat class="rounded-xl border border-slate-100 bg-white overflow-hidden shadow-soft">
                            <q-item-label header class="text-weight-bold text-uppercase text-xs letter-spacing-1 q-pa-lg bg-slate-50 border-b border-slate-100">Materie Abilitate</q-item-label>
                            <q-list separator>
                                <q-item v-if="teacherSubjects.length === 0" class="q-pa-xl text-center">
                                    <q-item-section>
                                      <q-icon name="menu_book" size="48px" class="opacity-10 q-mb-sm" />
                                      <div class="text-slate-400 italic">Nessuna materia assegnata a questo docente</div>
                                    </q-item-section>
                                </q-item>
                                <q-item v-for="ts in teacherSubjects" :key="ts.id" class="q-pa-md">
                                    <q-item-section avatar>
                                      <q-avatar color="indigo-50" text-color="indigo-700" icon="book" size="32px" />
                                    </q-item-section>
                                    <q-item-section>
                                        <q-item-label class="text-weight-medium">{{ ts.subject_name }}</q-item-label>
                                    </q-item-section>
                                    <q-item-section side>
                                        <q-btn flat round icon="delete" color="negative" size="sm" @click="removeTeacherSubject(ts.subject_id)">
                                          <q-tooltip>Rimuovi Abilitazione</q-tooltip>
                                        </q-btn>
                                    </q-item-section>
                                </q-item>
                            </q-list>
                        </q-card>
                    </div>
                    <div class="col-md-5 col-12">
                         <q-card flat class="rounded-xl border border-slate-100 bg-white q-pa-xl shadow-soft">
                            <div class="text-h6 text-weight-bold q-mb-lg">Aggiungi Abilitazione</div>
                            <q-select
                                v-model="selectedSubjectToAdd"
                                :options="availableSubjects"
                                label="Seleziona Materia"
                                outlined
                                emit-value
                                map-options
                                class="q-mb-xl"
                            />
                            <q-btn label="Assegna Materia" color="primary" unelevated class="full-width rounded-lg q-py-md shadow-sm" no-caps @click="addTeacherSubject" :disable="!selectedSubjectToAdd" />
                         </q-card>
                    </div>
                </div>
            </q-card-section>
        </q-card>
    </q-dialog>

    <!-- Reset Password Dialog -->
    <q-dialog v-model="showResetPwdDialog" class="premium-dialog">
        <q-card style="min-width: 400px" class="rounded-xl overflow-hidden shadow-24 bg-white">
            <q-card-section class="bg-gradient-primary text-white q-pa-lg">
                <div class="text-h5 text-weight-bold">Reset Password</div>
                <div class="text-subtitle1 opacity-80">{{ resetTargetName }}</div>
            </q-card-section>
            <q-card-section class="q-pa-xl">
                <q-input v-model="newPassword" label="Nuova Password" type="password" outlined />
                <div class="text-caption text-slate-400 q-mt-sm">Inserire una password sicura di almeno 8 caratteri.</div>
                
                <div class="row justify-end q-mt-xl q-gutter-sm">
                  <q-btn flat label="Annulla" color="slate-400" v-close-popup no-caps />
                  <q-btn color="primary" label="Aggiorna Password" class="q-px-xl rounded-lg shadow-sm" no-caps @click="handleResetPwd" :disable="newPassword.length < 8" />
                </div>
            </q-card-section>
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
import { usePermissions } from 'src/composables/usePermissions';

const $q = useQuasar();
const authStore = useAuthStore();
const { isSuperAdmin } = usePermissions();
const loading = ref(false);
const showUserDialog = ref(false);
const showClassDialog = ref(false);
const showImport = ref(false);
const isEditing = ref(false);
const importFile = ref(null);
const currentRoleFilter = ref('all');
const filterSchoolId = ref(null);

const users = ref([]); 

// Classes
const classes = ref([]);
const loadingClasses = ref(false);
const classOptions = computed(() => {
    return classes.value.map(c => ({
        label: c.name || `${c.section} ${c.academic_year}`,
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
    school_id: null,
    fiscal_code: '',
    password: ''
});

const getCurrentAcademicYear = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth() + 1; // 1-12
  if (month >= 9) { 
    return `${year}/${year + 1}`;
  } else {
    return `${year - 1}/${year}`;
  }
}

const currentYearStr = getCurrentAcademicYear();

const classForm = reactive({
    name: '',
    section: '',
    academic_year: currentYearStr
});


const schools = ref([]);
const schoolOptions = computed(() => schools.value);

onMounted(() => {
    fetchUsers()
    fetchClasses()
    if (isSuperAdmin.value) {
        fetchSchools()
    }
});

const fetchUsers = async () => {
    loading.value = true
    try {
        const res = await userService.getAll({ 
            role: currentRoleFilter.value === 'all' ? undefined : currentRoleFilter.value,
            school_id: filterSchoolId.value || undefined
        })
        // FILTER: Remove 'admin' users from the list as requested
        users.value = (res.data.users || []).filter(u => u.role !== 'admin')
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore durante il caricamento degli utenti' })
    } finally {
        loading.value = false
    }
}

watch(currentRoleFilter, () => {
    fetchUsers()
})

watch(filterSchoolId, () => {
    fetchUsers()
})

const fetchClasses = async () => {
    const targetSchoolId = isSuperAdmin.value ? userForm.school_id : authStore.user?.school_id;
    if (!targetSchoolId) return;
    loadingClasses.value = true;
    try {
        const res = await adminService.getSchoolClasses(targetSchoolId);
        classes.value = res.data || [];
    } catch (e) {
        console.error("Error loading classes", e);
    } finally {
        loadingClasses.value = false;
    }
}

const fetchClassesForSchool = () => {
    classes.value = [];
    userForm.class_id = null;
    fetchClasses();
}

const fetchSchools = async () => {
    try {
        const res = await adminService.getSchools({ page_size: 100 });
        schools.value = res.data.items || [];
    } catch (e) {
        console.error("Error loading schools", e);
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
    userForm.school_id = isSuperAdmin.value ? null : authStore.user?.school_id;
    userForm.fiscal_code = '';
    userForm.password = '';
    
    showUserDialog.value = true;
};

const openEdit = (user) => {
    isEditing.value = true;
    Object.assign(userForm, user);
    if (user.class_id) userForm.class_id = user.class_id;
    if (user.school_id) userForm.school_id = user.school_id;
    if (isSuperAdmin.value && userForm.school_id) {
        fetchClasses();
    }
    showUserDialog.value = true;
};

const openClassDialog = () => {
    classForm.name = '';
    classForm.section = '';
    classForm.academic_year = getCurrentAcademicYear();
    showClassDialog.value = true;
};

const saveUser = async () => {
    try {
        const targetSchoolId = isSuperAdmin.value ? userForm.school_id : authStore.user.school_id;
        const payload = { ...userForm, school_id: targetSchoolId };
        if (isEditing.value) {
            delete payload.password; 
            await userService.update(userForm.id, payload);
             $q.notify({ type: 'positive', message: 'Profilo utente aggiornato' });
        } else {
             await userService.create(payload);
             $q.notify({ type: 'positive', message: 'Nuovo utente creato con successo' });
        }
        showUserDialog.value = false;
        fetchUsers();
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore nel salvataggio dell\'utente', caption: e.message });
    }
};

const saveClass = async () => {
    try {
        const targetSchoolId = isSuperAdmin.value ? userForm.school_id : authStore.user.school_id;
        const payload = {
            name: classForm.name,
            academic_year: classForm.academic_year,
            school_id: targetSchoolId,
            section: classForm.name.replace(/[0-9]/g, '')
        };
        await adminService.createClass(payload);
        $q.notify({ type: 'positive', message: 'Nuova classe attivata' });
        showClassDialog.value = false;
        fetchClasses();
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore nella creazione della classe', caption: e.message });
    }
};

const confirmDelete = (user) => {
    $q.dialog({
        title: 'Conferma Eliminazione',
        message: `Sei sicuro di voler eliminare definitivamente l'utente ${user.first_name} ${user.last_name}? L'azione non può essere annullata.`,
        cancel: true,
        persistent: true,
        ok: { color: 'negative', label: 'Elimina', flat: false }
    }).onOk(async () => {
        try {
            await userService.delete(user.id);
            $q.notify({ type: 'positive', message: 'Utente eliminato correttamente' });
            fetchUsers();
        } catch(e) {
            $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione' });
        }
    });
};

const showResetPwdDialog = ref(false);
const resetTargetId = ref(null);
const resetTargetName = ref('');
const newPassword = ref('');

const openResetPwd = (user) => {
    resetTargetId.value = user.id;
    resetTargetName.value = `${user.last_name} ${user.first_name}`;
    newPassword.value = '';
    showResetPwdDialog.value = true;
};

const handleResetPwd = async () => {
    try {
        await userService.forceResetPassword(resetTargetId.value, newPassword.value);
        $q.notify({ type: 'positive', message: 'Password resettata con successo' });
        showResetPwdDialog.value = false;
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore durante il reset della password' });
    }
};

const bulkDelete = (selected) => {
    $q.dialog({
        title: 'Eliminazione Massiva',
        message: `Vuoi procedere con l'eliminazione di ${selected.length} profili selezionati?`,
        cancel: true,
        persistent: true,
        ok: { color: 'negative', label: 'Elimina Tutti', flat: false }
    }).onOk(async () => {
        try {
            const ids = selected.map(u => u.id);
            await userService.bulkDelete(ids);
            $q.notify({ type: 'positive', message: `${selected.length} profili rimossi con successo` });
            fetchUsers();
        } catch(e) {
            $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione massiva' });
        }
    });
};

const handleImport = async () => {
    if (!importFile.value) return;
    try {
        const res = await userService.bulkImport(importFile.value);
        $q.notify({ type: 'positive', message: `Importazione completata: ${res.data.created} utenti creati.` });
        showImport.value = false;
        importFile.value = null;
        fetchUsers();
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore durante l\'importazione dei dati' });
    }
};

const exportUsers = () => {
    const content = [
        'ID,Nome,Cognome,Email,Ruolo,Classe',
        ...filteredUsers.value.map(u => `${u.id},${u.first_name},${u.last_name},${u.email},${u.role},${u.class_name || ''}`)
    ].join('\r\n');

    const status = exportFile(
        'utenti_esportazione.csv',
        content,
        'text/csv'
    );
    if (!status) $q.notify({ type: 'negative', message: 'Esportazione fallita' });
};

// Teacher Subjects Logic
const showSubjectsDialog = ref(false)
const teacherSubjects = ref([])
const availableSubjects = ref([])
const selectedSubjectToAdd = ref(null)
const currentTeacherId = ref(null)
const currentTeacherName = ref('')
const teachersCache = ref([]) 

const openManageSubjects = async (user) => {
    currentTeacherName.value = `${user.last_name} ${user.first_name}`
    currentTeacherId.value = null
    teacherSubjects.value = []
    selectedSubjectToAdd.value = null
    showSubjectsDialog.value = true
    
    try {
        const targetSchoolId = user.school_id || authStore.user.school_id;
        if (teachersCache.value.length === 0) {
            const tRes = await adminService.getTeachersList(targetSchoolId)
            teachersCache.value = tRes.data || []
        }
        const teacher = teachersCache.value.find(t => t.user_id === user.id)
        if (!teacher) {
            $q.notify({ type: 'warning', message: 'Profilo docente non inizializzato correttamente.' })
            showSubjectsDialog.value = false
            return
        }
        currentTeacherId.value = teacher.id
        
        await loadTeacherSubjects()
        if (availableSubjects.value.length === 0) {
            const targetSchoolId = user.school_id || authStore.user.school_id;
            const sRes = await adminService.getSubjects(targetSchoolId)
            availableSubjects.value = (sRes.data || []).map(s => ({ label: s.name, value: s.id }))
        }

    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore nel caricamento delle materie abilitate' })
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
        $q.notify({ type: 'positive', message: 'Materia assegnata correttamente' })
        loadTeacherSubjects()
        selectedSubjectToAdd.value = null
    } catch(e) {
         $q.notify({ type: 'negative', message: 'Errore nell\'assegnazione della materia' })
    }
}

const removeTeacherSubject = async (subjectId) => {
    try {
        await adminService.removeTeacherSubject(currentTeacherId.value, subjectId)
        loadTeacherSubjects()
        $q.notify({ type: 'positive', message: 'Abilitazione rimossa' })
    } catch(e) {
         $q.notify({ type: 'negative', message: 'Errore nella rimozione della materia' })
    }
}
</script>

<style scoped>
.letter-spacing-1 { letter-spacing: 1px; }
.uppercase-input :deep(input) { text-transform: uppercase; }
.opacity-10 { opacity: 0.1; }
</style>

