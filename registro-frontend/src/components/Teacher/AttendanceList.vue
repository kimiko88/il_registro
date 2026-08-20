<template>
  <q-list separator bordered class="rounded-borders">
    <q-item v-for="record in localRecords" :key="record.studentId" class="q-py-md">
      <q-item-section avatar>
        <q-avatar color="primary" text-color="white">{{ record.name?.charAt(0) || '?' }}</q-avatar>
      </q-item-section>

      <q-item-section>
        <q-item-label class="text-weight-bold">{{ record.name }}</q-item-label>
        <q-item-label caption>
          {{ getStatusLabel(record.status) }} <span v-if="record.time">{{ t('attendance.atTime') || 'alle' }} {{ record.time }}</span>
        </q-item-label>
      </q-item-section>

      <q-item-section>
        <div class="row q-gutter-sm justify-end">
             <q-btn-group unelevated>
                <q-btn 
                    :color="record.status === 'Present' ? 'green' : 'grey-3'"
                    :text-color="record.status === 'Present' ? 'white' : 'black'"
                    label="P" 
                    dense size="sm"
                    :aria-label="t('attendance.present') || 'Presente'"
                    @click="updateStatus(record, 'Present')"
                >
                    <q-tooltip>{{ t('attendance.present') || 'Presente' }}</q-tooltip>
                </q-btn>
                <q-btn 
                    :color="record.status === 'Absent' ? 'red' : 'grey-3'"
                    :text-color="record.status === 'Absent' ? 'white' : 'black'"
                    label="A" 
                    dense size="sm"
                    :aria-label="t('attendance.absent') || 'Assente'"
                    @click="updateStatus(record, 'Absent')"
                >
                    <q-tooltip>{{ t('attendance.absent') || 'Assente' }}</q-tooltip>
                </q-btn>
                <q-btn 
                    :color="record.status === 'Late' ? 'orange' : 'grey-3'"
                    :text-color="record.status === 'Late' ? 'white' : 'black'"
                    label="L" 
                    dense size="sm"
                    :aria-label="t('attendance.late') || 'In Ritardo'"
                    @click="promptLate(record)"
                >
                    <q-tooltip>{{ t('attendance.late') || 'In Ritardo' }}</q-tooltip>
                </q-btn>
             </q-btn-group>
             <q-btn flat round icon="edit_note" size="sm" :aria-label="t('classRegister.addNote') || 'Aggiungi Nota'" @click="editNote(record)" :color="record.notes ? 'primary' : 'grey'">
                <q-tooltip>{{ t('classRegister.addNote') || 'Aggiungi Nota' }}</q-tooltip>
             </q-btn>
        </div>
      </q-item-section>
    </q-item>
  </q-list>

  <!-- Note Dialog -->
  <q-dialog v-model="showNoteDialog">
    <q-card style="min-width: 300px">
        <q-card-section class="text-h6">{{ t('classRegister.addNote') || 'Aggiungi Nota' }}</q-card-section>
        <q-card-section>
            <q-input v-model="currentNote" autofocus dense outlined />
        </q-card-section>
        <q-card-actions align="right">
            <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup />
            <q-btn flat :label="t('common.save') || 'Salva'" color="primary" @click="saveNote" />
        </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Late Dialog -->
  <q-dialog v-model="showLateDialog">
    <q-card style="min-width: 300px">
        <q-card-section class="text-h6">{{ t('classRegister.tableHeaderEntryTime') || 'Ora Ingresso Ritardo' }}</q-card-section>
        <q-card-section>
            <q-input v-model="currentTime" type="time" filled />
        </q-card-section>
        <q-card-actions align="right">
            <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup />
            <q-btn flat :label="t('common.confirm') || 'Conferma'" color="orange" @click="confirmLate" />
        </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAttendanceStore } from '@/stores/attendance';

const { t } = useI18n();
const attendanceStore = useAttendanceStore();
const localRecords = computed(() => attendanceStore.records);

const showNoteDialog = ref(false);
const showLateDialog = ref(false);
const activeRecord = ref(null);
const currentNote = ref('');
const currentTime = ref('08:10');

const updateStatus = (record, status) => {
    record.status = status;
    if (status === 'Present') {
        record.time = '';
    }
};

const promptLate = (record) => {
    activeRecord.value = record;
    currentTime.value = '08:10';
    showLateDialog.value = true;
};

const confirmLate = () => {
    if (activeRecord.value) {
        activeRecord.value.status = 'Late';
        activeRecord.value.time = currentTime.value;
    }
    showLateDialog.value = false;
};

const editNote = (record) => {
    activeRecord.value = record;
    currentNote.value = record.notes;
    showNoteDialog.value = true;
};

const getStatusLabel = (status) => {
    switch (status) {
        case 'Present': return t('attendance.present') || 'Presente';
        case 'Absent': return t('attendance.absent') || 'Assente';
        case 'Late': return t('attendance.late') || 'In Ritardo';
        default: return status;
    }
};

const saveNote = () => {
    if (activeRecord.value) {
        activeRecord.value.notes = currentNote.value;
    }
    showNoteDialog.value = false;
};
</script>
