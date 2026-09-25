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
          v-if="currentTab === 'rooms'"
          :label="t('timetableConstraints.newConstraintBtn')"
          icon="add"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openRoomReqModal"
        />
        <q-btn
          v-else-if="currentTab === 'groups'"
          :label="t('timetableConstraints.newGroupBtn')"
          icon="add"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openGroupModal"
        />
        <q-btn
          v-else-if="currentTab === 'spezzoni'"
          :label="t('timetableConstraints.newSpezzoneBtn')"
          icon="add"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openSpezzoneModal"
        />
      </div>
    </div>

    <!-- Policy Banner: Lab Constraints Only, Desiderata Excluded -->
    <q-banner rounded class="bg-indigo-50 border border-indigo-200 text-indigo-950 q-mb-lg rounded-xl shadow-xs">
      <template v-slot:avatar>
        <q-icon name="info" color="indigo-7" size="28px" />
      </template>
      <div class="text-body2 font-medium">
        <strong>Regole di Generazione Orario:</strong> I desiderata personali dei docenti sono esclusi dall'algoritmo. Vengono applicati solo i vincoli di laboratorio con il relativo numero di ore scelte, la copertura completa delle ore per ogni materia della classe, le cattedre non ancora assegnate e la compresenza in più classi per gruppi linguistici / articolati.
      </div>
    </q-banner>

    <!-- Navigation Tabs -->
    <q-tabs
      v-model="currentTab"
      dense
      align="left"
      active-color="primary"
      indicator-color="primary"
      class="text-grey-7 q-mb-lg bg-white rounded-xl shadow-sm overflow-hidden"
      :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'"
    >
      <q-tab name="rooms" icon="meeting_room" :label="t('timetableConstraints.tabRooms')" class="q-py-md text-weight-bold" />
      <q-tab name="groups" icon="groups" :label="t('timetableConstraints.tabGroups')" class="q-py-md text-weight-bold" />
      <q-tab name="spezzoni" icon="hourglass_empty" :label="t('timetableConstraints.tabSpezzoni')" class="q-py-md text-weight-bold" />
    </q-tabs>

    <!-- TAB 1: Subject Room Requirements (Aule e Laboratori con scelta ore) -->
    <div v-show="currentTab === 'rooms'">
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
          <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.addConstraintBtn')" icon="add" class="rounded-xl q-mt-sm font-bold" @click="openRoomReqModal" />
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

          <template v-slot:body-cell-lab_hours="props">
            <q-td :props="props" align="center">
              <q-badge color="purple-8" class="q-px-sm q-py-xs font-bold">
                {{ props.row.lab_hours || 1 }} ore / sett.
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
    </div>

    <!-- TAB 2: Gruppi Linguistici & Classi Articolate (Cattedre Congiunte) -->
    <div v-show="currentTab === 'groups'">
      <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <q-card-section class="row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold row items-center gap-2">
              <q-icon name="groups" color="indigo-7" />
              {{ t('timetableConstraints.groupsSectionTitle') }}
            </div>
            <div class="text-caption text-grey-6">
              {{ t('timetableConstraints.groupsSectionSubtitle') }}
            </div>
          </div>
          <div class="row items-center gap-2">
            <q-btn flat dense icon="refresh" color="primary" no-caps :label="t('timetableConstraints.refresh')" @click="fetchConstraints" />
            <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.newGroupBtn')" icon="add" class="rounded-xl font-bold" @click="openGroupModal" />
          </div>
        </q-card-section>

        <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

        <div v-if="loadingConstraints" class="row justify-center q-pa-lg">
          <q-spinner-dots size="40px" color="primary" />
        </div>

        <div v-else-if="associatedGroupsList.length === 0" class="text-center q-pa-xl">
          <q-icon name="groups" size="48px" color="grey-5" />
          <div class="text-h6 text-grey-6 q-mt-sm">{{ t('timetableConstraints.noGroupsTitle') }}</div>
          <p class="text-grey-5">{{ t('timetableConstraints.noGroupsDesc') }}</p>
          <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.newGroupBtn')" icon="add" class="rounded-xl q-mt-sm font-bold" @click="openGroupModal" />
        </div>

        <q-table
          v-else
          flat
          :rows="associatedGroupsList"
          :columns="groupColumns"
          row-key="id"
          class="no-shadow"
        >
          <template v-slot:body-cell-name="props">
            <q-td :props="props">
              <span class="text-weight-bold text-slate-800">{{ props.row.name }}</span>
            </q-td>
          </template>

          <template v-slot:body-cell-subject_name="props">
            <q-td :props="props">
              <q-badge color="indigo-8" class="q-px-sm q-py-xs font-bold">
                {{ props.row.subject_name }}
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-teacher_name="props">
            <q-td :props="props">
              <q-chip size="sm" icon="person" color="blue-1" text-color="blue-9" class="font-bold">
                {{ props.row.teacher_name }}
              </q-chip>
            </q-td>
          </template>

          <template v-slot:body-cell-classes="props">
            <q-td :props="props">
              <div class="row gap-1">
                <q-badge v-for="cName in props.row.class_names" :key="cName" color="teal-8" class="q-px-xs">
                  {{ cName }}
                </q-badge>
              </div>
            </q-td>
          </template>

          <template v-slot:body-cell-hours_per_week="props">
            <q-td :props="props" align="center">
              <q-badge color="purple-9" class="q-px-sm q-py-xs text-weight-bold">
                {{ props.row.hours_per_week }}h / sett.
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-actions="props">
            <q-td :props="props" align="right">
              <q-btn flat round dense color="negative" icon="delete" @click="confirmDeleteGroup(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card>
    </div>

    <!-- TAB 3: Spezzoni Orari & Cattedre Vacanti per Classe -->
    <div v-show="currentTab === 'spezzoni'">
      <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <q-card-section class="row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold row items-center gap-2">
              <q-icon name="hourglass_empty" color="amber-8" />
              {{ t('timetableConstraints.spezzoniSectionTitle') }}
            </div>
            <div class="text-caption text-grey-6">
              {{ t('timetableConstraints.spezzoniSectionSubtitle') }}
            </div>
          </div>
          <div class="row items-center gap-2">
            <q-btn flat dense icon="refresh" color="primary" no-caps :label="t('timetableConstraints.refresh')" @click="fetchConstraints" />
            <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.newSpezzoneBtn')" icon="add" class="rounded-xl font-bold" @click="openSpezzoneModal" />
          </div>
        </q-card-section>

        <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

        <div v-if="loadingConstraints" class="row justify-center q-pa-lg">
          <q-spinner-dots size="40px" color="primary" />
        </div>

        <div v-else-if="spezzoniList.length === 0" class="text-center q-pa-xl">
          <q-icon name="work_outline" size="48px" color="grey-5" />
          <div class="text-h6 text-grey-6 q-mt-sm">{{ t('timetableConstraints.noSpezzoniTitle') }}</div>
          <p class="text-grey-5">{{ t('timetableConstraints.noSpezzoniDesc') }}</p>
          <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.newSpezzoneBtn')" icon="add" class="rounded-xl q-mt-sm font-bold" @click="openSpezzoneModal" />
        </div>

        <q-table
          v-else
          flat
          :rows="spezzoniList"
          :columns="spezzoniColumns"
          row-key="id"
          class="no-shadow"
        >
          <template v-slot:body-cell-class_name="props">
            <q-td :props="props">
              <q-badge color="indigo-8" class="q-px-sm q-py-xs font-bold text-sm">
                {{ props.row.class_name }}
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-subject_name="props">
            <q-td :props="props">
              <span class="text-weight-bold text-slate-800">{{ props.row.subject_name }}</span>
            </q-td>
          </template>

          <template v-slot:body-cell-hours_per_week="props">
            <q-td :props="props" align="center">
              <q-badge color="amber-9" class="q-px-sm q-py-xs text-weight-bold">
                {{ props.row.hours_per_week }}h / sett.
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-placeholder_name="props">
            <q-td :props="props">
              <q-chip size="sm" icon="person_outline" color="blue-1" text-color="blue-9" class="font-bold">
                {{ props.row.placeholder_name }}
              </q-chip>
            </q-td>
          </template>

          <template v-slot:body-cell-notes="props">
            <q-td :props="props">
              <span class="text-caption text-grey-7">{{ props.row.notes || '-' }}</span>
            </q-td>
          </template>

          <template v-slot:body-cell-actions="props">
            <q-td :props="props" align="right">
              <q-btn flat round dense color="negative" icon="delete" @click="confirmDeleteSpezzone(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card>
    </div>

    <!-- Modal: Add Room Requirement with Lab Hours -->
    <q-dialog v-model="showReqModal" persistent>
      <q-card style="min-width: 460px;" class="rounded-2xl q-pa-sm">
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
          />

          <q-select
            v-model="reqForm.required_room_type"
            :options="roomTypeOptions"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.selectRoomType')"
          />

          <q-input
            v-model.number="reqForm.lab_hours"
            type="number"
            min="1"
            max="12"
            outlined
            dense
            :label="t('timetableConstraints.labHoursLabel')"
            :hint="t('timetableConstraints.labHoursHint')"
          />

          <q-toggle
            v-model="reqForm.is_mandatory"
            :label="t('timetableConstraints.mandatoryDesc')"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="t('timetableConstraints.saveBtn')" no-caps class="rounded-xl q-px-md font-bold" @click="saveReq" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Add Gruppo Linguistico / Articolato -->
    <q-dialog v-model="showGroupModal" persistent>
      <q-card style="min-width: 520px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold row items-center gap-2">
            <q-icon name="groups" color="primary" />
            {{ t('timetableConstraints.groupModalTitle') }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-input
            v-model="groupForm.name"
            outlined
            dense
            :label="t('timetableConstraints.groupName')"
            placeholder="es. Gruppo Spagnolo 3A-3B"
          />

          <q-select
            v-model="groupForm.subject_id"
            :options="subjects"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.groupSubject')"
          />

          <q-select
            v-model="groupForm.teacher_id"
            :options="allTeachers"
            option-value="id"
            option-label="label"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.groupTeacher')"
          />

          <q-select
            v-model="groupForm.class_ids"
            :options="classes"
            option-value="id"
            option-label="label"
            emit-value
            map-options
            multiple
            use-chips
            outlined
            dense
            :label="t('timetableConstraints.groupClasses')"
          />

          <q-input
            v-model.number="groupForm.hours_per_week"
            type="number"
            min="1"
            max="10"
            outlined
            dense
            :label="t('timetableConstraints.groupHours')"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="t('timetableConstraints.saveBtn') || 'Salva Gruppo'" no-caps class="rounded-xl q-px-md font-bold" @click="saveGroup" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Add Spezzone Orario per Classe -->
    <q-dialog v-model="showSpezzoneModal" persistent>
      <q-card style="min-width: 500px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold row items-center gap-2">
            <q-icon name="add_circle_outline" color="primary" />
            {{ t('timetableConstraints.spezzoneModalTitle') }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-select
            v-model="spezzoneForm.class_id"
            :options="classes"
            option-value="id"
            option-label="label"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.spezzoneClass')"
          />

          <q-select
            v-model="spezzoneForm.subject_id"
            :options="subjects"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.spezzoneSubject')"
            @update:model-value="onSpezzoneSubjectChanged"
          />

          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-input
                v-model.number="spezzoneForm.hours_per_week"
                type="number"
                min="1"
                max="18"
                outlined
                dense
                :label="t('timetableConstraints.spezzoneHours')"
              />
            </div>
            <div class="col-12 col-sm-6">
              <q-input
                v-model="spezzoneForm.placeholder_name"
                outlined
                dense
                :label="t('timetableConstraints.spezzonePlaceholder')"
              />
            </div>
          </div>

          <q-input
            v-model="spezzoneForm.notes"
            type="textarea"
            rows="2"
            outlined
            dense
            :label="t('timetableConstraints.spezzoneNotes')"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="t('timetableConstraints.saveBtn') || 'Salva'" no-caps class="rounded-xl q-px-md font-bold" @click="saveSpezzone" />
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
import adminService from '@/services/adminService';
import api from '@/services/api';

const { t } = useI18n();
const $q = useQuasar();

const currentTab = ref('rooms');

// ----------------- TAB 1: Subject Room Requirements -----------------
const loadingReqs = ref(false);
const roomReqs = ref([]);
const subjects = ref([]);

const showReqModal = ref(false);
const reqForm = ref({
  subject_id: null,
  required_room_type: 'lab_informatica',
  lab_hours: 2,
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
  { name: 'lab_hours', label: t('timetableConstraints.colLabHours'), field: 'lab_hours', align: 'center' },
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
    lab_hours: 2,
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

// ----------------- TAB 2: Gruppi Linguistici & Articolati -----------------
const allTeachers = ref([]);
const associatedGroupsList = ref([]);
const showGroupModal = ref(false);
const groupForm = ref({
  name: '',
  subject_id: null,
  teacher_id: null,
  class_ids: [],
  hours_per_week: 3
});

const groupColumns = computed(() => [
  { name: 'name', label: t('timetableConstraints.colGroup'), field: 'name', sortable: true, align: 'left' },
  { name: 'subject_name', label: t('timetableConstraints.colSubject'), field: 'subject_name', align: 'left' },
  { name: 'teacher_name', label: t('timetableConstraints.colTeacher'), field: 'teacher_name', align: 'left' },
  { name: 'classes', label: t('timetableConstraints.colClasses'), field: 'class_names', align: 'left' },
  { name: 'hours_per_week', label: t('timetableConstraints.colHours'), field: 'hours_per_week', align: 'center' },
  { name: 'actions', label: t('timetableConstraints.colActions'), field: 'id', align: 'right' }
]);

async function fetchTeachers() {
  try {
    const res = await adminService.getTeachersList();
    const raw = res.data || [];
    allTeachers.value = raw.map(t => ({
      id: t.id || t.user_id,
      user_id: t.user_id,
      label: `${t.last_name || ''} ${t.first_name || ''}`.trim() || t.email,
      email: t.email
    })).sort((a, b) => a.label.localeCompare(b.label));
  } catch (e) {
    console.error('Failed fetching teachers', e);
  }
}

function openGroupModal() {
  groupForm.value = {
    name: '',
    subject_id: subjects.value[0]?.id || null,
    teacher_id: allTeachers.value[0]?.id || null,
    class_ids: [],
    hours_per_week: 3
  };
  showGroupModal.value = true;
}

async function saveGroup() {
  if (!groupForm.value.name || !groupForm.value.subject_id || !groupForm.value.teacher_id || groupForm.value.class_ids.length < 2) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi e seleziona almeno 2 classi per il gruppo linguistico/articolato' });
    return;
  }

  const subj = subjects.value.find(s => s.id === groupForm.value.subject_id);
  const teacher = allTeachers.value.find(t => t.id === groupForm.value.teacher_id);
  const classNames = groupForm.value.class_ids.map(cid => {
    const c = classes.value.find(cl => cl.id === cid);
    return c ? c.label : cid;
  });

  const payload = {
    constraint_type: 'associated_group',
    target_type: 'group',
    is_hard: true,
    priority: 10,
    is_active: true,
    parameters: {
      name: groupForm.value.name,
      subject_id: groupForm.value.subject_id,
      subject_name: subj?.name || 'Materia',
      teacher_id: groupForm.value.teacher_id,
      teacher_name: teacher?.label || 'Docente',
      class_ids: groupForm.value.class_ids,
      class_names: classNames,
      hours_per_week: groupForm.value.hours_per_week || 3
    }
  };

  try {
    await timetableGenService.saveConstraint(payload);
    $q.notify({ type: 'positive', message: 'Gruppo linguistico / articolato configurato con successo' });
    showGroupModal.value = false;
    fetchConstraints();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore nel salvataggio del gruppo' });
  }
}

function confirmDeleteGroup(row) {
  $q.dialog({
    title: t('timetableConstraints.removeGroupTitle'),
    message: t('timetableConstraints.removeGroupMsg', { name: row.name }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await timetableGenService.deleteConstraint(row.id);
      $q.notify({ type: 'positive', message: 'Gruppo associato rimosso' });
      fetchConstraints();
    } catch {
      $q.notify({ type: 'negative', message: 'Errore nella rimozione del gruppo' });
    }
  });
}

// ----------------- TAB 3: Spezzoni Orari & Cattedre Vacanti -----------------
const loadingConstraints = ref(false);
const constraints = ref([]);
const spezzoniList = ref([]);
const classes = ref([]);

const showSpezzoneModal = ref(false);
const spezzoneForm = ref({
  class_id: null,
  subject_id: null,
  hours_per_week: 2,
  placeholder_name: 'Docente da Nominare (Cattedra non assegnata)',
  notes: ''
});

const spezzoniColumns = computed(() => [
  { name: 'class_name', label: t('timetableConstraints.colClass'), field: 'class_name', sortable: true, align: 'left' },
  { name: 'subject_name', label: t('timetableConstraints.colSubject'), field: 'subject_name', sortable: true, align: 'left' },
  { name: 'hours_per_week', label: t('timetableConstraints.colHours'), field: 'hours_per_week', align: 'center' },
  { name: 'placeholder_name', label: t('timetableConstraints.colPlaceholder'), field: 'placeholder_name', align: 'left' },
  { name: 'notes', label: t('timetableConstraints.colNotes'), field: 'notes', align: 'left' },
  { name: 'actions', label: t('timetableConstraints.colActions'), field: 'id', align: 'right' }
]);

async function fetchClasses() {
  try {
    const res = await adminService.getClasses();
    const raw = res.data || [];
    classes.value = raw.map(c => ({
      id: c.id,
      label: `${c.year || ''}${c.section || ''} ${c.name || ''}`.trim() || c.id
    })).sort((a, b) => a.label.localeCompare(b.label));
  } catch (err) {
    console.error('Error fetching classes', err);
  }
}

async function fetchConstraints() {
  loadingConstraints.value = true;
  try {
    const res = await timetableGenService.getConstraints();
    const raw = res.data || [];
    constraints.value = raw;

    spezzoniList.value = raw
      .filter(c => c.constraint_type === 'unassigned_hours')
      .map(c => {
        const p = typeof c.parameters === 'string' ? JSON.parse(c.parameters) : (c.parameters || {});
        const cls = classes.value.find(cl => cl.id === c.target_id);
        return {
          id: c.id,
          class_id: c.target_id,
          class_name: cls ? cls.label : (c.target_id?.substring(0, 8) || '-'),
          subject_id: p.subject_id,
          subject_name: p.subject_name || '-',
          hours_per_week: p.hours_per_week || 2,
          placeholder_name: p.placeholder_name || 'Docente da Nominare',
          notes: p.notes || ''
        };
      });

    associatedGroupsList.value = raw
      .filter(c => c.constraint_type === 'associated_group')
      .map(c => {
        const p = typeof c.parameters === 'string' ? JSON.parse(c.parameters) : (c.parameters || {});
        return {
          id: c.id,
          name: p.name || 'Gruppo Senza Nome',
          subject_id: p.subject_id,
          subject_name: p.subject_name || 'Materia',
          teacher_id: p.teacher_id,
          teacher_name: p.teacher_name || 'Docente',
          class_ids: p.class_ids || [],
          class_names: p.class_names || [],
          hours_per_week: p.hours_per_week || 3
        };
      });
  } catch (err) {
    console.error('Error fetching constraints', err);
  } finally {
    loadingConstraints.value = false;
  }
}

function openSpezzoneModal() {
  spezzoneForm.value = {
    class_id: classes.value[0]?.id || null,
    subject_id: subjects.value[0]?.id || null,
    hours_per_week: 2,
    placeholder_name: `Docente da Nominare (${subjects.value[0]?.name || ''})`,
    notes: ''
  };
  showSpezzoneModal.value = true;
}

function onSpezzoneSubjectChanged(subjId) {
  const subj = subjects.value.find(s => s.id === subjId);
  if (subj) {
    spezzoneForm.value.placeholder_name = `Docente da Nominare (${subj.name})`;
  }
}

async function saveSpezzone() {
  if (!spezzoneForm.value.class_id || !spezzoneForm.value.subject_id) return;
  const subj = subjects.value.find(s => s.id === spezzoneForm.value.subject_id);

  const payload = {
    constraint_type: 'unassigned_hours',
    target_type: 'class',
    target_id: spezzoneForm.value.class_id,
    is_hard: true,
    priority: 10,
    is_active: true,
    parameters: {
      subject_id: spezzoneForm.value.subject_id,
      subject_name: subj?.name || 'Materia',
      hours_per_week: spezzoneForm.value.hours_per_week,
      placeholder_name: spezzoneForm.value.placeholder_name || 'Docente da Nominare',
      notes: spezzoneForm.value.notes || ''
    }
  };

  try {
    await timetableGenService.saveConstraint(payload);
    $q.notify({ type: 'positive', message: t('timetableConstraints.spezzoneSaved') });
    showSpezzoneModal.value = false;
    fetchConstraints();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('timetableConstraints.notifySaveError') });
  }
}

function confirmDeleteSpezzone(row) {
  $q.dialog({
    title: t('timetableConstraints.removeSpezzoneTitle'),
    message: t('timetableConstraints.removeSpezzoneMsg', { hours: row.hours_per_week, subject: row.subject_name, class: row.class_name }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await timetableGenService.deleteConstraint(row.id);
      $q.notify({ type: 'positive', message: t('timetableConstraints.spezzoneRemoved') });
      fetchConstraints();
    } catch {
      $q.notify({ type: 'negative', message: t('timetableConstraints.notifyDeleteError') });
    }
  });
}

onMounted(async () => {
  await Promise.all([
    fetchRoomReqs(),
    fetchSubjects(),
    fetchTeachers(),
    fetchClasses()
  ]);
  await fetchConstraints();
});
</script>
