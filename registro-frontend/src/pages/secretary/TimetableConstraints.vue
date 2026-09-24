<template>
  <q-page padding class="min-h-screen" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-slate-50 text-slate-800'">
    <!-- Top Header -->
    <div class="row items-center justify-between q-mb-lg gap-4">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none row items-center gap-2" :class="$q.dark.isActive ? 'text-white' : 'text-slate-800'">
          <q-icon name="tune" color="primary" size="36px" />
          {{ t('timetableConstraints.title') }}
        </h1>
        <p class="text-subtitle1 q-mt-xs q-mb-none" :class="$q.dark.isActive ? 'text-grey-4' : 'text-slate-500'">
          {{ t('timetableConstraints.subtitle') }}
        </p>
      </div>

      <div class="row items-center gap-2">
        <q-btn
          :label="t('timetableConstraints.newConstraintBtn')"
          icon="add"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openRoomReqModal"
        />
      </div>
    </div>

    <!-- Section 1: Subject Room Requirements -->
    <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
      <q-card-section class="row items-center justify-between">
        <div>
          <div class="text-h6 text-weight-bold row items-center gap-2">
            <q-icon name="meeting_room" color="indigo-7" />
            {{ t('timetableConstraints.sectionTitle') }}
          </div>
          <div class="text-caption text-grey-6">
            {{ t('timetableConstraints.sectionSubtitle') }}
          </div>
        </div>
        <q-btn flat dense icon="refresh" color="primary" no-caps :label="t('timetableConstraints.refresh')" @click="fetchRoomReqs" />
      </q-card-section>

      <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

      <div v-if="loadingReqs" class="row justify-center q-pa-lg">
        <q-spinner-dots size="40px" color="primary" />
      </div>

      <div v-else-if="roomReqs.length === 0" class="text-center q-pa-xl">
        <q-icon name="science" size="48px" color="grey-5" />
        <div class="text-h6 text-grey-6 q-mt-sm">{{ t('timetableConstraints.noConstraintsTitle') }}</div>
        <p class="text-grey-5">{{ t('timetableConstraints.noConstraintsDesc') }}</p>
        <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.addConstraintBtn')" icon="add" class="rounded-xl q-mt-sm" @click="openRoomReqModal" />
      </div>

      <q-table
        v-else
        flat
        :rows="roomReqs"
        :columns="roomReqColumns"
        row-key="id"
        class="no-shadow"
      >
        <template v-slot:body-cell-required_room_type="props">
          <q-td :props="props">
            <q-badge color="indigo-7" class="q-px-sm q-py-xs font-bold">
              {{ getRoomTypeLabel(props.row.required_room_type) }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-is_mandatory="props">
          <q-td :props="props" align="center">
            <q-badge :color="props.row.is_mandatory ? 'negative' : 'warning'">
              {{ props.row.is_mandatory ? t('timetableConstraints.mandatoryBadge') : t('timetableConstraints.preferentialBadge') }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" align="right">
            <q-btn flat round dense color="negative" icon="delete" @click="confirmDeleteReq(props.row)" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Modal: Add Room Requirement -->
    <q-dialog v-model="showReqModal" persistent>
      <q-card style="min-width: 440px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold">{{ t('timetableConstraints.modalTitle') }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-select
            v-model="reqForm.subject_id"
            :options="subjects"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.selectSubject')"
            :rules="[val => !!val || t('timetableConstraints.subjectRequired')]"
          />

          <q-select
            v-model="reqForm.required_room_type"
            :options="roomTypeOptions"
            option-value="value"
            option-label="label"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.selectRoomType')"
            :rules="[val => !!val || t('timetableConstraints.roomTypeRequired')]"
          />

          <q-toggle
            v-model="reqForm.is_mandatory"
            :label="t('timetableConstraints.mandatoryDesc')"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel')" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="t('timetableConstraints.saveBtn')" no-caps class="rounded-xl q-px-md font-bold" @click="saveReq" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import timetableGenService from '@/services/timetableGenService';
import api from '@/services/api';

const $q = useQuasar();
const { t } = useI18n();

const roomReqs = ref([]);
const subjects = ref([]);
const loadingReqs = ref(false);

const showReqModal = ref(false);
const reqForm = ref({
  subject_id: null,
  required_room_type: 'lab_informatica',
  is_mandatory: true
});

const roomTypeOptions = computed(() => [
  { label: t('roomsManagement.types.lab_info') || 'Laboratorio di Informatica', value: 'lab_informatica' },
  { label: t('roomsManagement.types.lab_science') || 'Laboratorio di Scienze', value: 'lab_scienze' },
  { label: t('roomsManagement.types.lab_lang') || 'Laboratorio di Lingue', value: 'lab_lingue' },
  { label: t('roomsManagement.types.lab_chemistry') || 'Laboratorio di Chimica', value: 'lab_chimica' },
  { label: t('roomsManagement.types.lab_physics') || 'Laboratorio di Fisica', value: 'lab_fisica' },
  { label: t('roomsManagement.types.art') || 'Laboratorio di Arte', value: 'lab_arte' },
  { label: t('roomsManagement.types.gym') || 'Palestra', value: 'palestra' },
  { label: t('roomsManagement.types.auditorium') || 'Aula Magna', value: 'aula_magna' },
  { label: t('roomsManagement.types.library') || 'Biblioteca', value: 'biblioteca' }
]);

const roomReqColumns = computed(() => [
  { name: 'subject_name', label: t('timetableConstraints.colSubject'), field: 'subject_name', sortable: true, align: 'left' },
  { name: 'required_room_type', label: t('timetableConstraints.colRoomType'), field: 'required_room_type', align: 'left' },
  { name: 'is_mandatory', label: t('timetableConstraints.colConstraint'), field: 'is_mandatory', align: 'center' },
  { name: 'actions', label: t('timetableConstraints.colActions'), field: 'id', align: 'right' }
]);

function getRoomTypeLabel(val) {
  const f = roomTypeOptions.value.find(o => o.value === val);
  return f ? f.label : val;
}

async function fetchRoomReqs() {
  loadingReqs.value = true;
  try {
    const res = await timetableGenService.getRoomRequirements();
    roomReqs.value = res.data || [];
  } catch {
    $q.notify({ type: 'negative', message: t('timetableConstraints.notifyFetchError') });
  } finally {
    loadingReqs.value = false;
  }
}

async function fetchSubjects() {
  try {
    const res = await api.get('/subjects');
    subjects.value = res.data || [];
  } catch (err) {
    console.error('Error fetching subjects', err);
  }
}

function openRoomReqModal() {
  reqForm.value = {
    subject_id: subjects.value[0]?.id || null,
    required_room_type: 'lab_informatica',
    is_mandatory: true
  };
  showReqModal.value = true;
}

async function saveReq() {
  if (!reqForm.value.subject_id || !reqForm.value.required_room_type) return;
  try {
    await timetableGenService.saveRoomRequirement(reqForm.value);
    $q.notify({ type: 'positive', message: t('timetableConstraints.notifySaved') });
    showReqModal.value = false;
    fetchRoomReqs();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('timetableConstraints.notifySaveError') });
  }
}

function confirmDeleteReq(row) {
  $q.dialog({
    title: t('timetableConstraints.removeTitle'),
    message: t('timetableConstraints.removeMsg', { subject: row.subject_name }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await timetableGenService.deleteRoomRequirement(row.id);
      $q.notify({ type: 'positive', message: t('timetableConstraints.notifyRemoved') });
      fetchRoomReqs();
    } catch {
      $q.notify({ type: 'negative', message: t('timetableConstraints.notifyDeleteError') });
    }
  });
}

onMounted(() => {
  fetchRoomReqs();
  fetchSubjects();
});
</script>
