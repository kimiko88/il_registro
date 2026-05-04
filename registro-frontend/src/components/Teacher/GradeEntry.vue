<template>
  <q-card>
    <q-table
       :rows="studentsWithGrades"
       :columns="columns"
       row-key="student_id"
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
           <q-tr :props="props" v-if="props.row">
              <q-td key="name" :props="props">
                 <div class="text-weight-bold">{{ props.row.full_name }}</div>
                 <div class="text-caption text-grey">Assenze: {{ props.row.absences || 0 }}</div>
              </q-td>

              <q-td key="current_grade" :props="props" style="width: 250px">
                 <div class="row items-center no-wrap q-gutter-sm">
                    <q-input
                       v-model.number="entryData[props.row.student_id].value"
                      type="number"
                      dense outlined
                      placeholder="-"
                      style="width: 70px"
                      @keydown.tab="focusNext(props.rowIndex)"
                      :bg-color="getGradeColor(entryData[props.row.student_id].value)"
                      min="1" max="10" step="0.5"
                   />
                    <q-input
                       v-model="entryData[props.row.student_id].notes"
                      dense outlined
                      placeholder="Note..."
                      class="col"
                   />
                    <q-btn 
                       icon="save" round flat dense color="primary" 
                       :disable="!isDirty(props.row.student_id)"
                       @click="saveLine(props.row.student_id)" 
                    />
                </div>
             </q-td>

             <q-td key="history" :props="props">
                <div class="row q-gutter-xs">
                   <q-badge v-for="g in props.row.grades" :key="g.id" :color="getBadgeColor(g.grade_value)">
                      {{ g.grade_value }}
                      <q-tooltip>{{ g.date }} - {{ g.grade_type }}</q-tooltip>
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
import { gradeService } from 'src/services/gradeService';

const props = defineProps(['classId', 'subject', 'date', 'type']);
const emit = defineEmits(['refresh']);
const $q = useQuasar();
const gradesStore = useGradesStore();

const loading = ref(false);
const isOnline = ref(navigator.onLine);
const entryData = ref({});
const initialSnapshot = ref({});

const studentsWithGrades = computed(() => {
    // The gradesStore should have students for the selected class
    if (gradesStore.grades && gradesStore.grades.students) {
        return gradesStore.grades.students;
    }
    return [];
});

const columns = [
    { name: 'name', label: 'Studente', align: 'left' },
    { name: 'current_grade', label: 'Voto', align: 'left' },
    { name: 'history', label: 'Storico', align: 'left' },
    { name: 'average', label: 'Media', align: 'center' },
    { name: 'actions', label: '', align: 'right' }
];

const initData = () => {
    if (!studentsWithGrades.value) return;
    studentsWithGrades.value.forEach(s => {
        if (!entryData.value[s.student_id]) {
            entryData.value[s.student_id] = { value: null, notes: '' };
        }
    });
    initialSnapshot.value = JSON.parse(JSON.stringify(entryData.value));
};

watch(studentsWithGrades, initData, { immediate: true, deep: true });

const hasChanges = computed(() => {
    return JSON.stringify(entryData.value) !== JSON.stringify(initialSnapshot.value);
});

const isDirty = (id) => {
    if (!entryData.value[id] || !initialSnapshot.value[id]) return false;
    return JSON.stringify(entryData.value[id]) !== JSON.stringify(initialSnapshot.value[id]);
};

const getBadgeColor = (val) => val < 6 ? 'red' : (val >= 8 ? 'green' : 'blue');
const getGradeColor = (val) => val ? (val < 6 ? 'bg-red-1' : 'bg-green-1') : '';
const getAverageClass = (val) => val < 6 ? 'text-red text-weight-bold' : 'text-green text-weight-bold';

const saveLine = async (id) => {
    const data = entryData.value[id];
    if (!data.value) return;

    try {
        let payloadType = "numeric";
        if (typeof data.value === 'string' && isNaN(Number(data.value))) {
            payloadType = "judgment";
        }
        let mappedType = null;
        if (props.type === "Scritto") mappedType = "Written";
        if (props.type === "Orale") mappedType = "Oral";
        if (props.type === "Pratico") mappedType = "Practical";

        await gradeService.saveGrade({
            student_id: id,
            subject_id: props.subject,
            grade_value: Number(data.value),
            grade_type: payloadType,
            semester: 1, // Defaulting to 1 for MVP
            description: data.notes,
            date: props.date,
            is_published: true,
            grade_category: "formative",
            evaluation_type: mappedType
        });

        $q.notify({ message: `Voto ${data.value} salvato`, color: 'positive', timeout: 500 });
        initialSnapshot.value[id] = JSON.parse(JSON.stringify(data));
        emit('refresh');
    } catch (err) {
        $q.notify({ type: 'negative', message: 'Errore nel salvataggio' });
    }
};

const saveAll = async () => {
    loading.value = true;
    try {
        for (const s of studentsWithGrades.value) {
            if (isDirty(s.student_id)) {
                await saveLine(s.student_id);
            }
        }
        $q.notify({ type: 'positive', message: 'Tutti i voti sono stati salvati.'});
    } catch (err) {
        $q.notify({ type: 'negative', message: 'Errore durante il salvataggio multiplo' });
    } finally {
        loading.value = false;
    }
};

const focusNext = (index) => {
    // Logic to focus next input would ideally use refs map
};

window.addEventListener('online', () => isOnline.value = true);
window.addEventListener('offline', () => isOnline.value = false);

onMounted(() => {
    // Initialize if data already present
    initData();
});
</script>
