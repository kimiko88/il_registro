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
                    <q-btn flat round icon="school" color="secondary" tooltip="Voti/Assenze" />
                </q-td>
            </template>
        </q-table>
    </q-card>

    <!-- Enrollment Dialog -->
    <q-dialog v-model="showEnrollment" persistent maximized transition-show="slide-up" transition-hide="slide-down">
       <StudentEnrollmentForm @complete="onEnrollmentComplete" @cancel="showEnrollment = false" />
        <q-btn dense flat icon="close" v-close-popup class="absolute-top-right q-ma-sm" />
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import StudentEnrollmentForm from 'src/components/Secretary/StudentEnrollmentForm.vue'

const showEnrollment = ref(false)
const filter = ref('')
const students = ref([
    { id: 1, name: 'Rossi Mario', class: '1A', dob: '2010-05-20', parent: 'Rossi Luigi' },
    { id: 2, name: 'Bianchi Sofia', class: '2B', dob: '2009-11-12', parent: 'Bianchi Anna' }
])

const columns = [
    { name: 'name', label: 'Nome', field: 'name', align: 'left', sortable: true },
    { name: 'class', label: 'Classe', field: 'class', align: 'center', sortable: true },
    { name: 'dob', label: 'Data di Nascita', field: 'dob', align: 'center' },
    { name: 'parent', label: 'Genitore', field: 'parent', align: 'left' },
    { name: 'actions', label: 'Azioni', align: 'right' }
]

const onEnrollmentComplete = (data) => {
    showEnrollment.value = false
    // Add to list mock
    students.value.unshift({
        id: Date.now(),
        name: `${data.lastName} ${data.firstName}`,
        class: data.class,
        dob: data.dob,
        parent: data.parent1Name
    })
}

const editStudent = (row) => {
    // Logic to open edit form (could reuse enrollment form in edit mode)
}
</script>
