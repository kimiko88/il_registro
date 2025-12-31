<template>
  <q-card>
    <q-table
       :rows="studentsWithGrades"
       :columns="columns"
       row-key="id"
       :loading="loading"
       separator="cell"
       dense
    >
       <template v-slot:header="props">
          <q-tr :props="props">
            <q-th key="name" :props="props">Studente</q-th>
            <q-th key="current_grade" :props="props">Nuovo Voto ({{ date }})</q-th>
            <q-th key="history" :props="props">Storico ({{ subject }})</q-th>
            <q-th key="average" :props="props">Media</q-th>
            <q-th key="actions" :props="props">Azioni</q-th>
          </q-tr>
       </template>

       <template v-slot:body="props">
          <q-tr :props="props">
             <q-td key="name" :props="props">
                <div class="text-weight-bold">{{ props.row.name }}</div>
                <div class="text-caption text-grey">Assenze: {{ props.row.absences }}</div>
             </q-td>

             <q-td key="current_grade" :props="props" style="width: 250px">
                <div class="row items-center no-wrap q-gutter-sm">
                   <q-input
                      v-model.number="entryData[props.row.id].value"
                      type="number"
                      dense outlined
                      placeholder="-"
                      style="width: 70px"
                      @keydown.tab="focusNext(props.rowIndex)"
                      :bg-color="getGradeColor(entryData[props.row.id].value)"
                      min="1" max="10" step="0.5"
                   />
                   <q-input
                      v-model="entryData[props.row.id].notes"
                      dense outlined
                      placeholder="Note..."
                      class="col"
                   />
                   <q-btn 
                      icon="save" round flat dense color="primary" 
                      :disable="!isDirty(props.row.id)"
                      @click="saveLine(props.row.id)" 
                   />
                </div>
             </q-td>

             <q-td key="history" :props="props">
                <div class="row q-gutter-xs">
                   <q-badge v-for="g in props.row.grades" :key="g.id" :color="getBadgeColor(g.value)">
                      {{ g.value }}
                      <q-tooltip>{{ g.date }} - {{ g.type }}</q-tooltip>
                   </q-badge>
                </div>
             </q-td>

             <q-td key="average" :props="props" class="text-center">
                <span :class="getAverageClass(props.row.average)">{{ props.row.average }}</span>
             </q-td>
             
             <q-td key="actions" :props="props" auto-width>
                 <q-btn flat round icon="history" color="grey-7" tooltip="Vedi Dettagli" />
             </q-td>
          </q-tr>
       </template>

       <template v-slot:bottom>
           <div class="row full-width justify-between items-center q-pa-sm">
               <div class="text-caption">
                  <q-icon name="offline_pin" color="green" v-if="isOnline" />
                  <q-icon name="signal_wifi_off" color="warning" v-else />
                  {{ isOnline ? 'Online - Dati sincronizzati' : 'Offline - Modifiche salvate in locale' }}
               </div>
               <div>
                  <q-btn label="Salva Tutti" color="primary" icon="save_alt" @click="saveAll" :disable="!hasChanges" />
               </div>
           </div>
       </template>
    </q-table>
  </q-card>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useGradesStore } from 'src/stores/grades';
import { useQuasar } from 'quasar';

const props = defineProps(['classId', 'subject', 'date', 'type']);
const emit = defineEmits(['refresh']);
const $q = useQuasar();

const loading = ref(false);
const isOnline = ref(navigator.onLine);
const entryData = ref({});
const initialSnapshot = ref({});

// Mock data generation
const studentsWithGrades = ref([
    { id: 'S1', name: 'Rossi Mario', absences: 2, grades: [{id: 1, value: 7, date: '2025-01-10', type: 'Orale'}], average: 7.0 },
    { id: 'S2', name: 'Bianchi Anna', absences: 0, grades: [{id: 2, value: 8.5, date: '2025-01-12', type: 'Scritto'}], average: 8.5 },
    { id: 'S3', name: 'Verdi Paolo', absences: 5, grades: [{id: 3, value: 5, date: '2025-01-15', type: 'Pratico'}], average: 5.0 }
]);

const columns = [
    { name: 'name', label: 'Studente', align: 'left' },
    { name: 'current_grade', label: 'Voto', align: 'left' },
    { name: 'history', label: 'Storico', align: 'left' },
    { name: 'average', label: 'Media', align: 'center' },
    { name: 'actions', label: '', align: 'right' }
];

// Initialize entry data
const initData = () => {
    studentsWithGrades.value.forEach(s => {
        if (!entryData.value[s.id]) {
            entryData.value[s.id] = { value: null, notes: '' };
        }
    });
    initialSnapshot.value = JSON.parse(JSON.stringify(entryData.value));
};

watch(studentsWithGrades, initData, { immediate: true });

const hasChanges = computed(() => {
    return JSON.stringify(entryData.value) !== JSON.stringify(initialSnapshot.value);
});

const isDirty = (id) => {
    return JSON.stringify(entryData.value[id]) !== JSON.stringify(initialSnapshot.value[id]);
};

const getBadgeColor = (val) => val < 6 ? 'red' : (val >= 8 ? 'green' : 'blue');
const getGradeColor = (val) => val ? (val < 6 ? 'bg-red-1' : 'bg-green-1') : '';
const getAverageClass = (val) => val < 6 ? 'text-red text-weight-bold' : 'text-green text-weight-bold';

const saveLine = async (id) => {
    const data = entryData.value[id];
    if (!data.value) return;

    // Simulate save
    $q.notify({ message: `Voto ${data.value} salvato per ${id}`, color: 'positive', timeout: 500 });
    
    // Update snapshot for this item
    initialSnapshot.value[id] = JSON.parse(JSON.stringify(data));
    
    // Add to local history mock
    studentsWithGrades.value.find(s => s.id === id).grades.push({
        id: Date.now(),
        value: data.value,
        date: props.date,
        type: props.type
    });
};

const saveAll = async () => {
    loading.value = true;
    // Simulate bulk API
    setTimeout(() => {
        loading.value = false;
        $q.notify({ type: 'positive', message: 'Tutti i voti sono stati salvati.'});
        // Reset dirty state
        initialSnapshot.value = JSON.parse(JSON.stringify(entryData.value));
    }, 1000);
};

const focusNext = (index) => {
    // Logic to focus next input would ideally use refs map
};

// Check connectivity
window.addEventListener('online', () => isOnline.value = true);
window.addEventListener('offline', () => isOnline.value = false);
</script>
