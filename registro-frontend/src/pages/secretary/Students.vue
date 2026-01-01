<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4 text-weight-bold">Registro Studenti</div>
       <q-btn color="primary" icon="person_add" label="Nuova Iscrizione" @click="showEnrollment = true" />
    </div>

    <!-- Student List -->
    <q-card>
        <q-table
            :rows="students"
            :columns="columns"
            :filter="filter"
            :loading="loading"
            row-key="id"
        >
            <template v-slot:top-right>
                <q-input dense debounce="300" v-model="filter" placeholder="Cerca studente...">
                    <template v-slot:append>
                        <q-icon name="search" />
                    </template>
                </q-input>
            </template>
            
             <template v-slot:body-cell-actions="props">
                <q-td :props="props" auto-width>
                    <q-btn flat round icon="edit" color="primary" @click="editStudent(props.row)" />
                    <q-btn flat round icon="school" color="secondary" tooltip="Voti/Assenze" @click="$q.notify('Funzionalità voti in sviluppo')" />
                </q-td>
            </template>
        </q-table>
    </q-card>

    <!-- Generic Enrollment Dialog (Placeholder for now) -->
    <q-dialog v-model="showEnrollment" persistent maximized transition-show="slide-up" transition-hide="slide-down">
         <q-card>
            <q-toolbar class="bg-primary text-white">
                <q-toolbar-title>Nuova Iscrizione</q-toolbar-title>
                <q-btn flat round dense icon="close" v-close-popup />
            </q-toolbar>
            <q-card-section>
               <div class="text-h6">Modulo di iscrizione (WIP)</div>
               <p>Utilizzare la pagina 'Gestione Utenti' per creare nuovi account studente per ora.</p>
               <q-btn label="Vai a Gestione Utenti" to="/secretary/users" color="primary" />
            </q-card-section>
         </q-card>
    </q-dialog>

    <!-- Edit Dialog (Class Assignment) -->
    <q-dialog v-model="showEditDialog">
        <q-card style="min-width: 400px">
            <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Modifica Studente</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section>
                <q-form @submit="saveStudent" class="q-gutter-md">
                    <div class="text-subtitle1">{{ editForm.first_name }} {{ editForm.last_name }}</div>
                    
                    <q-input v-model="editForm.email" label="Email" outlined dense readonly />
                    
                    <q-select
                        v-model="editForm.class_id"
                        :options="classOptions"
                        label="Classe"
                        outlined
                        dense
                        emit-value
                        map-options
                    />

                    <div class="row justify-end">
                        <q-btn label="Annulla" flat v-close-popup color="grey" />
                        <q-btn label="Salva" type="submit" color="primary" />
                    </div>
                </q-form>
            </q-card-section>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { userService } from 'src/services/userService'
import adminService from 'src/services/adminService'
import { useAuthStore } from 'src/stores/auth'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const authStore = useAuthStore()

const showEnrollment = ref(false)
const showEditDialog = ref(false)
const filter = ref('')
const loading = ref(false)
const students = ref([])
const classes = ref([])

const editForm = reactive({
    id: null,
    first_name: '',
    last_name: '',
    email: '',
    class_id: null,
    school_id: ''
})

const columns = [
    { name: 'name', label: 'Nome', field: row => `${row.last_name} ${row.first_name}`, align: 'left', sortable: true },
    { name: 'class', label: 'Classe', field: row => row.class_name || row.ClassName || '-', align: 'center', sortable: true },
    { name: 'email', label: 'Email', field: 'email', align: 'left' },
    { name: 'actions', label: 'Azioni', align: 'right' }
]

const classOptions = computed(() => {
    return classes.value.map(c => ({
        label: c.name || `${c.section} ${c.academic_year}`,
        value: c.id
    }))
})

onMounted(() => {
    fetchStudents()
    fetchClasses()
})

const fetchStudents = async () => {
    loading.value = true
    try {
        const res = await userService.getAll({ role: 'student' })
        students.value = res.data.users || []
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento studenti' })
    } finally {
        loading.value = false
    }
}

const fetchClasses = async () => {
    if (!authStore.user?.school_id) return
    try {
        const res = await adminService.getSchoolClasses(authStore.user.school_id)
        classes.value = res.data
    } catch (e) {
        console.error("Error loading classes", e)
    }
}

const editStudent = (row) => {
    editForm.id = row.id
    editForm.first_name = row.first_name
    editForm.last_name = row.last_name
    editForm.email = row.email
    editForm.class_id = row.ClassID || row.class_id // Handle both cases
    editForm.school_id = row.school_id
    showEditDialog.value = true
}

const saveStudent = async () => {
    try {
        const payload = {
            class_id: editForm.class_id,
            school_id: editForm.school_id // Required for validation/updates
        }
        await userService.update(editForm.id, payload)
        $q.notify({ type: 'positive', message: 'Studente aggiornato' })
        showEditDialog.value = false
        fetchStudents()
    } catch (e) {
        $q.notify({ type: 'negative', message: 'Errore aggiornamento' })
    }
}

const onEnrollmentComplete = (data) => {
    showEnrollment.value = false
    fetchStudents()
}
</script>
