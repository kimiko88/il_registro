<template>
  <q-card>
    <q-banner v-if="readOnly" class="bg-amber-1 text-amber-10 border-b q-pa-md">
      <template v-slot:avatar>
        <q-icon name="lock" color="amber-9" size="24px" />
      </template>
      <div class="text-weight-bold">{{ t('gradesPage.substitutionWarningTitle') || 'Modalità Supplenza - Consultazione Voti Disabilitata' }}</div>
      <div class="text-caption">{{ t('gradesPage.substitutionWarningDesc') || 'L\'inserimento e la visualizzazione del registro voti è riservata ai docenti titolari di materia della classe.' }}</div>
    </q-banner>

    <q-table
       v-else
       :rows="studentsWithGrades"
       :columns="columns"
       row-key="student_id"
       :loading="loading"
       separator="cell"
       dense
    >
       <template v-slot:header="props">
          <q-tr :props="props">
            <q-th key="name" :props="props">{{ t('classRegister.tableHeaderStudent') || 'Studente' }}</q-th>
            <q-th key="current_grade" :props="props">{{ isReligionSubject ? 'Nuovo Giudizio IRC' : (t('classRegister.tableHeaderGrade') || 'Nuovo Voto') }} ({{ date }})</q-th>
            <q-th key="history" :props="props">{{ t('common.history') || 'Storico' }}</q-th>
            <q-th key="average" :props="props">{{ isReligionSubject ? 'Ultimo Giudizio' : (t('gradesPage.average') || 'Media') }}</q-th>
            <q-th key="actions" :props="props">{{ t('common.actions') || 'Azioni' }}</q-th>
          </q-tr>
       </template>

       <template v-slot:body="props">
           <q-tr :props="props" v-if="props.row">
              <q-td key="name" :props="props">
                 <div class="text-weight-bold">{{ props.row.full_name }}</div>
                 <div class="text-caption text-grey">{{ t('classRegister.absent') || 'Assenze' }}: {{ props.row.absences || 0 }}</div>
              </q-td>

              <q-td key="current_grade" :props="props" style="width: 280px">
                 <div v-if="isReligionSubject && isStudentExemptFromReligion(props.row)" class="row items-center q-gutter-xs">
                     <q-badge
                        :color="props.row.religion_choice === 'attivita_alternativa' ? 'purple-7' : 'orange-8'"
                        text-color="white"
                        class="text-weight-bold q-py-xs q-px-sm rounded-md"
                     >
                        <q-icon :name="props.row.religion_choice === 'attivita_alternativa' ? 'swap_horiz' : 'block'" size="14px" class="q-mr-xs" />
                        {{ props.row.religion_choice === 'attivita_alternativa' ? 'Attività Alternativa' : 'Non Avvalente IRC' }}
                     </q-badge>
                     <q-tooltip>Questo studente non si avvale dell'IRC e non può ricevere valutazioni di religione</q-tooltip>
                 </div>
                 <div v-else-if="entryData[props.row.student_id]" class="row items-center no-wrap q-gutter-sm">
                     <q-select
                        v-model="entryData[props.row.student_id].value"
                       :options="currentGradeOptions"
                       dense outlined
                       placeholder="-"
                       :style="{ width: isReligionSubject ? '160px' : '100px' }"
                       @keydown.tab="focusNext(props.rowIndex)"
                       :bg-color="getGradeColor(entryData[props.row.student_id].value)"
                    />
                    <q-input
                       v-model="entryData[props.row.student_id].notes"
                      dense outlined
                      :placeholder="t('classRegister.tableHeaderGradeNotes') || 'Note...'"
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
                <!-- Se materia IRC, mostra tutti i giudizi ricevuti -->
                <div v-if="isReligionSubject" class="q-py-xs">
                   <div class="row items-center q-gutter-xs">
                      <q-badge
                         v-for="g in (props.row.grades || [])"
                         :key="g.id"
                         :color="getBadgeColor(g.description || g.grade_value)"
                         class="cursor-pointer text-weight-bold"
                         @click="editGradeDialog(g, props.row)"
                      >
                         {{ g.description || formatGrade(g.grade_value) }}
                         <q-tooltip anchor="top middle" self="bottom middle">
                            {{ formatDate(g.date) }} - {{ g.description || formatGrade(g.grade_value) }}
                         </q-tooltip>
                      </q-badge>
                      <span v-if="!props.row.grades || props.row.grades.length === 0" class="text-caption text-grey-4">-</span>
                   </div>
                </div>
                <!-- Altrimenti mostra la suddivisione standard Scritto / Orale / Pratico -->
                <div v-else class="q-py-xs">
                   <!-- Scritto -->
                   <div class="row items-center q-mb-xs" style="min-height: 24px">
                      <span class="text-caption text-grey-7 text-weight-bold q-mr-sm" style="min-width: 50px">{{ t('gradesPage.written') || 'Scritto' }}:</span>
                      <div class="row q-gutter-xs">
                         <q-badge
                            v-for="g in getGradesByType(props.row.grades, 'Written')"
                            :key="g.id"
                            :color="getBadgeColor(g.grade_value)"
                            class="cursor-pointer text-weight-bold"
                            @click="editGradeDialog(g, props.row)"
                         >
                            {{ formatGrade(g.grade_value) }}
                            <q-tooltip anchor="top middle" self="bottom middle">
                               {{ formatDate(g.date) }} - {{ g.description || t('common.noDescription') || 'Nessuna nota' }}
                            </q-tooltip>
                         </q-badge>
                         <span v-if="getGradesByType(props.row.grades, 'Written').length === 0" class="text-caption text-grey-4">-</span>
                      </div>
                   </div>
                   
                   <!-- Orale -->
                   <div class="row items-center q-mb-xs" style="min-height: 24px">
                      <span class="text-caption text-grey-7 text-weight-bold q-mr-sm" style="min-width: 50px">{{ t('gradesPage.oral') || 'Orale' }}:</span>
                      <div class="row q-gutter-xs">
                         <q-badge
                            v-for="g in getGradesByType(props.row.grades, 'Oral')"
                            :key="g.id"
                            :color="getBadgeColor(g.grade_value)"
                            class="cursor-pointer text-weight-bold"
                            @click="editGradeDialog(g, props.row)"
                         >
                            {{ formatGrade(g.grade_value) }}
                            <q-tooltip anchor="top middle" self="bottom middle">
                               {{ formatDate(g.date) }} - {{ g.description || t('common.noDescription') || 'Nessuna nota' }}
                            </q-tooltip>
                         </q-badge>
                         <span v-if="getGradesByType(props.row.grades, 'Oral').length === 0" class="text-caption text-grey-4">-</span>
                      </div>
                   </div>
                   
                   <!-- Pratico -->
                   <div class="row items-center" style="min-height: 24px">
                      <span class="text-caption text-grey-7 text-weight-bold q-mr-sm" style="min-width: 50px">{{ t('gradesPage.practical') || 'Pratico' }}:</span>
                      <div class="row q-gutter-xs">
                         <q-badge
                            v-for="g in getGradesByType(props.row.grades, 'Practical')"
                            :key="g.id"
                            :color="getBadgeColor(g.grade_value)"
                            class="cursor-pointer text-weight-bold"
                            @click="editGradeDialog(g, props.row)"
                         >
                            {{ formatGrade(g.grade_value) }}
                            <q-tooltip anchor="top middle" self="bottom middle">
                               {{ formatDate(g.date) }} - {{ g.description || 'Nessuna nota' }}
                            </q-tooltip>
                         </q-badge>
                         <span v-if="getGradesByType(props.row.grades, 'Practical').length === 0" class="text-caption text-grey-4">-</span>
                      </div>
                   </div>
                </div>
             </q-td>

             <q-td key="average" :props="props" class="text-center">
                <span :class="getAverageClass(props.row.average)">{{ props.row.average }}</span>
             </q-td>
             
             <q-td key="actions" :props="props" auto-width>
                 <q-btn flat round icon="history" color="grey-7" :aria-label="t('common.details') || 'Vedi Dettagli'">
                   <q-tooltip>{{ t('common.details') || 'Vedi Dettagli' }}</q-tooltip>
                 </q-btn>
             </q-td>
          </q-tr>
       </template>

       <template v-slot:bottom>
           <div class="row full-width justify-between items-center q-pa-sm">
               <div class="text-caption">
                  <q-icon name="offline_pin" color="green" v-if="isOnline" />
                  <q-icon name="signal_wifi_off" color="warning" v-else />
                  {{ isOnline ? (t('common.onlineSynced') || 'Online - Dati sincronizzati') : (t('common.offlineSaved') || 'Offline - Modifiche salvate in locale') }}
               </div>
               <div>
                  <q-btn label="Salva Tutti" color="primary" icon="save_alt" @click="saveAll" :disable="!hasChanges" />
               </div>
           </div>
       </template>
    </q-table>
  </q-card>

    <!-- Edit Grade Dialog -->
    <q-dialog v-model="showEditGradeDialog">
      <q-card style="min-width: 350px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6 text-weight-bold">{{ isReligionSubject ? 'Modifica Giudizio' : 'Modifica Voto' }} - {{ selectedStudentName }}</div>
        </q-card-section>
        
        <q-card-section class="q-pa-md q-gutter-md">
          <q-select
            v-model="editGradeForm.value"
            :options="currentGradeOptions"
            :label="isReligionSubject ? 'Giudizio IRC' : 'Voto'"
            outlined
            dense
            :bg-color="getGradeColor(editGradeForm.value)"
          />
          <q-input
            v-model="editGradeForm.date"
            type="date"
            label="Data Voto"
            outlined
            dense
          />
          <q-select
            v-if="!isReligionSubject"
            v-model="editGradeForm.evaluationType"
            :options="['Scritto', 'Orale', 'Pratico']"
            label="Tipo Valutazione"
            outlined
            dense
          />
          <q-input
            v-model="editGradeForm.notes"
            label="Titolo / Descrizione Voto"
            outlined
            dense
            autogrow
          />
        </q-card-section>
        
        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat label="Elimina" color="negative" @click="deleteGradeConfirm" />
          <q-space />
          <q-btn flat label="Annulla" v-close-popup color="grey-7" />
          <q-btn label="Salva" color="primary" @click="saveIndividualGradeEdit" />
        </q-card-actions>
      </q-card>
    </q-dialog>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useGradesStore } from '@/stores/grades';
import { useAuthStore } from '@/stores/auth';
import { useQuasar } from 'quasar';
import { gradeService } from '@/services/gradeService';
import {
  ITALIAN_GRADE_OPTIONS,
  RELIGION_JUDGMENT_OPTIONS,
  RELIGION_JUDGMENT_MAP,
  gradeToNumeric,
  formatGrade,
  getGradeColor
} from '@/utils/gradeUtils';

const { t } = useI18n();

const props = defineProps({
  classId: String,
  subject: String,
  isReligion: { type: Boolean, default: false },
  date: String,
  type: String,
  readOnly: { type: Boolean, default: false }
});

const emit = defineEmits(['refresh']);
const $q = useQuasar();
const gradesStore = useGradesStore();
const authStore = useAuthStore();

const loading = ref(false);
const isOnline = ref(navigator.onLine);
const entryData = ref({});
const initialSnapshot = ref({});

const isReligionSubject = computed(() => {
    if (gradesStore.grades?.is_religion_subject) return true;
    const currentSub = (gradesStore.subjects || []).find(s =>
        String(s.id) === String(props.subject) ||
        String(s.subject_id) === String(props.subject)
    );
    if (currentSub?.is_religion) return true;
    if (currentSub?.subject_name && currentSub.subject_name.toLowerCase().includes('religione')) return true;
    return false;
});

const currentGradeOptions = computed(() => {
    if (isReligionSubject.value) {
        return RELIGION_JUDGMENT_OPTIONS;
    }
    return ITALIAN_GRADE_OPTIONS;
});

const isStudentExemptFromReligion = (row) => {
    if (!isReligionSubject.value) return false;
    return row.religion_choice === 'non_avvalente' || row.religion_choice === 'attivita_alternativa';
};

const studentsWithGrades = computed(() => {
    if (props.readOnly) return [];
    if (gradesStore.grades && gradesStore.grades.students) {
        return gradesStore.grades.students.map(s => {
            let sum = 0;
            let count = 0;
            let lastJudgment = '';
            if (s.grades) {
                s.grades.forEach(g => {
                    if (isReligionSubject.value) {
                        if (g.description) lastJudgment = g.description;
                        count++;
                    } else if (typeof g.grade_value === 'number' && g.grade_value > 0) {
                        sum += g.grade_value;
                        count++;
                    }
                });
            }
            const avg = isReligionSubject.value
                ? (count > 0 ? (lastJudgment || `${count} valutaz.`) : '-')
                : (count > 0 ? (sum / count).toFixed(1) : '-');
            return {
                ...s,
                average: avg
            };
        });
    }
    return [];
});

const columns = computed(() => [
    { name: 'name', label: t('competenciesPage.student') || 'Studente', align: 'left' },
    { name: 'current_grade', label: t('classRegister.tableHeaderGrade') || 'Voto', align: 'left' },
    { name: 'history', label: t('common.history') || 'Storico', align: 'left' },
    { name: 'average', label: t('gradesPage.average') || 'Media', align: 'center' },
    { name: 'actions', label: '', align: 'right' }
]);

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

const getBadgeColor = (val) => {
    const numeric = typeof val === 'string' ? gradeToNumeric(val) : val;
    if (numeric === null || numeric === undefined || numeric === -1) return 'grey';
    return numeric < 6 ? 'red' : 'green';
};

const getAverageClass = (val) => {
    if (!val || val === '-') return 'text-grey';
    const numeric = parseFloat(val);
    if (isNaN(numeric)) return 'text-grey';
    return numeric < 6 ? 'text-red text-weight-bold' : 'text-green text-weight-bold';
};

const saveLine = async (id) => {
    if (props.readOnly) return;
    const data = entryData.value[id];
    if (data.value === null || data.value === undefined || data.value === '') return;

    try {
        let payloadType = "numeric";
        let gradeVal = gradeToNumeric(data.value);
        let desc = data.notes || '';

        if (isReligionSubject.value) {
            payloadType = "judgment";
            gradeVal = RELIGION_JUDGMENT_MAP[data.value] ?? (gradeToNumeric(data.value) || 6);
            desc = data.value + (data.notes ? ' - ' + data.notes : '');
        } else if (typeof data.value === 'string' && isNaN(Number(data.value))) {
            payloadType = "judgment";
            desc = data.value + (data.notes ? ' - ' + data.notes : '');
        }

        let mappedType = null;
        if (!isReligionSubject.value) {
            if (props.type === "Scritto") mappedType = "Written";
            if (props.type === "Orale") mappedType = "Oral";
            if (props.type === "Pratico") mappedType = "Practical";
        }

        await gradeService.saveGrade({
            student_id: id,
            subject_id: props.subject,
            grade_value: gradeVal,
            grade_type: payloadType,
            semester: 1,
            description: desc,
            date: props.date,
            is_published: true,
            grade_category: "formative",
            evaluation_type: mappedType
        });

        $q.notify({ message: `Voto ${data.value} salvato`, color: 'positive', timeout: 500 });
        initialSnapshot.value[id] = JSON.parse(JSON.stringify(data));
        emit('refresh');
    } catch (err) {
        $q.notify({ type: 'negative', message: t('common.error') });
    }
};

const saveAll = async () => {
    if (props.readOnly) return;
    loading.value = true;
    try {
        for (const s of studentsWithGrades.value) {
            if (isDirty(s.student_id)) {
                await saveLine(s.student_id);
            }
        }
        $q.notify({ type: 'positive', message: t('common.success') });
    } catch (err) {
        $q.notify({ type: 'negative', message: t('common.error') });
    } finally {
        loading.value = false;
    }
};

const focusNext = (_index) => {};

const handleOnline = () => { isOnline.value = true };
const handleOffline = () => { isOnline.value = false };

onMounted(() => {
    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);
    initData();
});

onUnmounted(() => {
  window.removeEventListener('online', handleOnline);
  window.removeEventListener('offline', handleOffline);
});

const getGradesByType = (grades, type) => {
    if (!grades) return [];
    return grades.filter(g => {
        const et = String(g.evaluation_type || '').toLowerCase();
        const target = String(type).toLowerCase();
        return et === target;
    });
};

const formatDate = (dateStr) => {
    if (!dateStr) return '';
    try {
        const d = new Date(dateStr);
        const day = String(d.getDate()).padStart(2, '0');
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const year = d.getFullYear();
        return `${day}/${month}/${year}`;
    } catch (e) {
        return dateStr;
    }
};

const showEditGradeDialog = ref(false);
const selectedStudentName = ref('');
const editGradeForm = ref({
    id: '',
    value: '',
    date: '',
    evaluationType: '',
    notes: ''
});

const editGradeDialog = (grade, student) => {
    if (props.readOnly) return;
    const currentUserId = authStore.user?.id;
    if (currentUserId && grade.teacher_id && grade.teacher_id !== currentUserId && grade.created_by && grade.created_by !== currentUserId) {
        $q.notify({
            type: 'warning',
            message: `${t('gradesPage.cannotEditOtherTeacherGrade') || 'Impossibile modificare il voto: è stato inserito da un altro docente'} (${grade.teacher_name || grade.created_by_name || 'Altro Docente'})`,
            caption: t('gradesPage.sharedGradesCaption') || 'Nei voti condivisi puoi modificare solo le valutazioni inserite da te.'
        });
        return;
    }
    selectedStudentName.value = student.full_name;
    let rawDate = '';
    if (grade.date) {
        try {
            rawDate = new Date(grade.date).toISOString().split('T')[0];
        } catch (e) { rawDate = ''; }
    }
    const typeMap = { 'Written': 'Scritto', 'Oral': 'Orale', 'Practical': 'Pratico' };
    editGradeForm.value = {
        id: grade.id,
        value: formatGrade(grade.grade_value),
        date: rawDate,
        evaluationType: typeMap[grade.evaluation_type] || 'Scritto',
        notes: grade.description || ''
    };
    showEditGradeDialog.value = true;
};

const saveIndividualGradeEdit = async () => {
    if (props.readOnly) return;
    let valNumeric = gradeToNumeric(editGradeForm.value.value);
    const typeMap = { 'Scritto': 'Written', 'Orale': 'Oral', 'Pratico': 'Practical' };
    let evalType = typeMap[editGradeForm.value.evaluationType] || editGradeForm.value.evaluationType;
    let desc = editGradeForm.value.notes;
    let gradeType = undefined;

    if (isReligionSubject.value) {
        valNumeric = RELIGION_JUDGMENT_MAP[editGradeForm.value.value] ?? (gradeToNumeric(editGradeForm.value.value) || 6);
        desc = editGradeForm.value.value + (editGradeForm.value.notes && editGradeForm.value.notes !== editGradeForm.value.value ? ' - ' + editGradeForm.value.notes : '');
        gradeType = 'judgment';
        evalType = undefined;
    }

    try {
        await gradesStore.updateGrade(editGradeForm.value.id, {
            grade_value: valNumeric,
            description: desc,
            date: editGradeForm.value.date || undefined,
            evaluation_type: evalType,
            grade_type: gradeType,
            reason: 'Modifica voto singola'
        });
        $q.notify({ type: 'positive', message: isReligionSubject.value ? 'Giudizio aggiornato' : t('common.success') });
        showEditGradeDialog.value = false;
        emit('refresh');
    } catch (err) {
        $q.notify({ type: 'negative', message: t('common.error') });
    }
};

const deleteGradeConfirm = async () => {
    if (props.readOnly) return;
    $q.dialog({
        title: t('common.confirmDelete') || 'Conferma Eliminazione',
        message: t('gradesPage.confirmDeleteGrade') || 'Sei sicuro di voler eliminare questo voto permanentemente?',
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            await gradesStore.deleteGrade(editGradeForm.value.id);
            $q.notify({ type: 'positive', message: t('common.success') });
            showEditGradeDialog.value = false;
            emit('refresh');
        } catch (err) {
            $q.notify({ type: 'negative', message: t('common.error') });
        }
    });
};
</script>
