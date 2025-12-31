<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-md">Le Mie Classi</div>

    <div class="row q-col-gutter-lg">
        <!-- Class List Sidebar -->
        <div class="col-12 col-md-3">
             <q-list bordered class="bg-white rounded-borders">
                 <q-item-label header class="text-weight-bold bg-grey-2">Elenco Classi</q-item-label>
                 <q-item 
                    v-for="cls in classes" 
                    :key="cls.id" 
                    clickable 
                    v-ripple
                    :active="selectedClass?.id === cls.id"
                    active-class="bg-blue-1 text-primary"
                    @click="selectClass(cls)"
                 >
                     <q-item-section avatar>
                         <q-avatar color="primary" text-color="white">{{ cls.name }}</q-avatar>
                     </q-item-section>
                     <q-item-section>
                         <q-item-label>{{ cls.name }}</q-item-label>
                         <q-item-label caption>{{ cls.students }} Studenti</q-item-label>
                     </q-item-section>
                     <q-item-section side v-if="cls.isCoordinator">
                         <q-icon name="star" color="orange"><q-tooltip>Coordinatore</q-tooltip></q-icon>
                     </q-item-section>
                 </q-item>
             </q-list>
        </div>

        <!-- Class Detail View -->
        <div class="col-12 col-md-9" v-if="selectedClass">
             <q-card>
                 <q-toolbar class="bg-primary text-white">
                     <q-toolbar-title>Classe {{ selectedClass.name }}</q-toolbar-title>
                     <q-tabs v-model="tab" shrink stretch>
                         <q-tab name="students" label="Studenti" />
                         <q-tab name="grades" label="Riepilogo Voti" />
                         <q-tab name="dashboard" label="Coordinatore" v-if="selectedClass.isCoordinator" />
                     </q-tabs>
                 </q-toolbar>

                 <q-separator />

                 <q-tab-panels v-model="tab" animated>
                     <!-- Student List -->
                     <q-tab-panel name="students">
                         <div class="row items-center justify-between q-mb-md">
                             <div class="text-h6">Elenco Studenti</div>
                             <q-input dense outlined placeholder="Cerca studente..." v-model="search" rounded>
                                 <template v-slot:append><q-icon name="search" /></template>
                             </q-input>
                         </div>
                         <q-table :rows="filteredStudents" :columns="columns" flat bordered row-key="id" />
                     </q-tab-panel>

                     <!-- Grades Summary -->
                     <q-tab-panel name="grades">
                         <div class="text-center text-grey text-h6 q-pa-xl">
                             <q-icon name="analytics" size="64px" class="q-mb-sm" />
                             <div>Grafici andamento classe (Mock)</div>
                             <q-btn label="Vedi Dettagli Voti" color="primary" flat class="q-mt-sm" to="/teacher/grades" />
                         </div>
                     </q-tab-panel>

                     <!-- Coordinator Dashboard -->
                     <q-tab-panel name="dashboard">
                         <div class="row q-col-gutter-md">
                             <div class="col-12 col-md-6">
                                 <q-card bordered flat class="bg-red-1">
                                     <q-card-section>
                                         <div class="text-subtitle1 text-red-9 text-weight-bold">Studenti a Rischio</div>
                                         <q-list dense>
                                             <q-item><q-item-label>• Bianchi Anna (Media: 5.2)</q-item-label></q-item>
                                             <q-item><q-item-label>• Verdi Paolo (Assenze: 25%)</q-item-label></q-item>
                                         </q-list>
                                     </q-card-section>
                                 </q-card>
                             </div>
                             <div class="col-12 col-md-6">
                                 <q-card bordered flat class="bg-yellow-1">
                                     <q-card-section>
                                         <div class="text-subtitle1 text-orange-9 text-weight-bold">Scadenze Coordinamento</div>
                                         <q-list dense>
                                             <q-item><q-item-label>• Approvazione PDP (3 mancanti)</q-item-label></q-item>
                                             <q-item><q-item-label>• Preparazione Consigli di Classe</q-item-label></q-item>
                                         </q-list>
                                     </q-card-section>
                                 </q-card>
                             </div>
                         </div>
                     </q-tab-panel>
                 </q-tab-panels>
             </q-card>
        </div>
        <div class="col-12 col-md-9 text-center text-grey" v-else>
             <div class="q-mt-xl">
                 <q-icon name="school" size="100px" />
                 <div class="text-h5">Seleziona una classe per gestire</div>
             </div>
        </div>
    </div>

  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'

const tab = ref('students')
const search = ref('')
const selectedClass = ref(null)

const classes = ref([
    { id: 1, name: '5A', students: 22, isCoordinator: true },
    { id: 2, name: '4B', students: 19, isCoordinator: false },
    { id: 3, name: '3C', students: 25, isCoordinator: false }
])

const studentsMock = [
    { id: 1, name: 'Rossi Mario', email: 'mario.rossi@school.it', avg: 7.5 },
    { id: 2, name: 'Bianchi Anna', email: 'anna.bianchi@school.it', avg: 5.2 },
    { id: 3, name: 'Verdi Paolo', email: 'paolo.verdi@school.it', avg: 6.0 },
]

const columns = [
    { name: 'name', label: 'Nome', field: 'name', align: 'left', sortable: true },
    { name: 'email', label: 'Email', field: 'email', align: 'left' },
    { name: 'avg', label: 'Media Attuale', field: 'avg', align: 'center', sortable: true }
]

const filteredStudents = computed(() => {
    if (!search.value) return studentsMock
    return studentsMock.filter(s => s.name.toLowerCase().includes(search.value.toLowerCase()))
})

const selectClass = (cls) => {
    selectedClass.value = cls
    tab.value = 'students'
}
</script>
