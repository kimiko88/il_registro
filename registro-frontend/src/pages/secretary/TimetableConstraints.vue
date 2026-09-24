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
      <q-tab name="teacher_prefs" icon="person_search" :label="t('timetableConstraints.tabTeacherPrefs')" class="q-py-md text-weight-bold" />
      <q-tab name="spezzoni" icon="hourglass_empty" :label="t('timetableConstraints.tabSpezzoni')" class="q-py-md text-weight-bold" />
    </q-tabs>

    <!-- TAB 1: Subject Room Requirements (Aule e Laboratori) -->
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
    </div>

    <!-- TAB 2: Teacher Preferences & Desiderata -->
    <div v-show="currentTab === 'teacher_prefs'">
      <!-- Teacher Selection Card -->
      <q-card flat bordered class="rounded-2xl shadow-sm q-mb-lg q-pa-md" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <div class="row items-center q-col-gutter-md">
          <div class="col-12 col-md-6">
            <div class="text-subtitle1 text-weight-bold row items-center gap-2 q-mb-xs">
              <q-icon name="person" color="primary" />
              {{ t('timetableConstraints.teacherSelectLabel') }}
            </div>
            <q-select
              v-model="selectedTeacher"
              :options="teacherOptions"
              option-value="id"
              option-label="label"
              outlined
              dense
              use-input
              input-debounce="200"
              @filter="filterTeachers"
              :placeholder="t('timetableConstraints.teacherSelectPlaceholder')"
              @update:model-value="onTeacherChanged"
            >
              <template v-slot:no-option>
                <q-item>
                  <q-item-section class="text-grey">Nessun docente trovato</q-item-section>
                </q-item>
              </template>
            </q-select>
          </div>

          <div class="col-12 col-md-6 row items-center justify-end gap-2 q-pt-md" v-if="selectedTeacher">
            <!-- Quick Day Off Button -->
            <q-btn-dropdown
              color="indigo-7"
              unelevated
              no-caps
              icon="event_busy"
              :label="t('timetableConstraints.quickDayOff')"
              class="rounded-xl text-weight-bold"
            >
              <q-list>
                <q-item
                  v-for="d in days"
                  :key="'dayoff-' + d.index"
                  clickable
                  v-close-popup
                  @click="setQuickDayOff(d.index, d.label)"
                >
                  <q-item-section avatar>
                    <q-icon name="block" color="negative" size="xs" />
                  </q-item-section>
                  <q-item-section>{{ d.label }}</q-item-section>
                </q-item>
              </q-list>
            </q-btn-dropdown>

            <!-- Reset to Neutral Button -->
            <q-btn
              flat
              dense
              icon="restart_alt"
              color="grey-7"
              no-caps
              :label="t('timetableConstraints.resetNeutral')"
              @click="resetTeacherGridToNeutral"
            />

            <!-- Save Teacher Prefs Button -->
            <q-btn
              color="positive"
              unelevated
              no-caps
              icon="save"
              :loading="savingTeacherPrefs"
              :label="t('timetableConstraints.saveTeacherPrefs')"
              class="rounded-xl text-weight-bold shadow-sm"
              @click="saveCurrentTeacherPreferences"
            />
          </div>
        </div>
      </q-card>

      <!-- Preferences Grid & Controls -->
      <div v-if="!selectedTeacher" class="text-center q-pa-xl bg-white rounded-2xl border border-slate-200">
        <q-icon name="touch_app" size="48px" color="grey-5" />
        <div class="text-h6 text-grey-6 q-mt-sm">{{ t('timetableConstraints.selectTeacherFirst') }}</div>
      </div>

      <div v-else>
        <!-- Legend and Info -->
        <q-card flat bordered class="rounded-2xl q-pa-md q-mb-lg shadow-sm" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
          <div class="row items-center justify-between gap-4">
            <div class="row items-center gap-4 text-sm font-medium">
              <span class="text-grey-6">{{ t('schedulePreferences.legend') }}</span>
              <div class="row items-center gap-1">
                <span class="inline-block w-4 h-4 rounded-full bg-emerald-500"></span>
                <span>🟢 {{ t('schedulePreferences.preferred') }}</span>
              </div>
              <div class="row items-center gap-1">
                <span class="inline-block w-4 h-4 rounded-full bg-slate-300"></span>
                <span>⚪ {{ t('schedulePreferences.neutral') }}</span>
              </div>
              <div class="row items-center gap-1">
                <span class="inline-block w-4 h-4 rounded-full bg-rose-500"></span>
                <span>🔴 {{ t('schedulePreferences.unavailable') }}</span>
              </div>
            </div>

            <div class="text-caption text-grey-6">
              {{ t('schedulePreferences.clickToChange') }}
            </div>
          </div>
        </q-card>

        <!-- Interactive Grid -->
        <q-card flat bordered class="rounded-2xl p-4 shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
          <div v-if="loadingTeacherPrefs" class="row justify-center q-pa-xl">
            <q-spinner-dots size="48px" color="primary" />
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full border-collapse">
              <thead>
                <tr class="bg-slate-100 text-slate-700">
                  <th class="p-3 border border-slate-200 text-left w-24">{{ t('schedulePreferences.hourCol') }}</th>
                  <th
                    v-for="d in days"
                    :key="d.index"
                    class="p-3 border border-slate-200 text-center font-bold text-base"
                  >
                    {{ d.label }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="h in 8" :key="h">
                  <td class="p-3 border border-slate-200 font-bold text-center bg-slate-50 text-slate-800">
                    {{ t('schedulePreferences.hourSlot', { hour: h }) }}
                  </td>
                  <td
                    v-for="d in days"
                    :key="d.index + '-' + h"
                    class="p-2 border border-slate-200 text-center cursor-pointer transition-colors select-none"
                    :class="getCellBgClass(d.index, h)"
                    @click="cyclePreference(d.index, h)"
                  >
                    <div class="py-2 px-1 rounded-lg">
                      <div class="font-bold text-xs" :class="getCellTextClass(d.index, h)">
                        {{ getCellLabel(d.index, h) }}
                      </div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </q-card>
      </div>
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
            outlined
            dense
            type="textarea"
            rows="2"
            :label="t('timetableConstraints.spezzoneNotes')"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel')" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="t('timetableConstraints.saveBtn')" no-caps class="rounded-xl q-px-md font-bold" @click="saveSpezzone" />
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

const $q = useQuasar();
const { t } = useI18n();

const currentTab = ref('rooms');

// ----------------- TAB 1: Room Requirements -----------------
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

// ----------------- TAB 2: Teacher Preferences & Desiderata -----------------
const allTeachers = ref([]);
const teacherOptions = ref([]);
const selectedTeacher = ref(null);
const loadingTeacherPrefs = ref(false);
const savingTeacherPrefs = ref(false);

const days = computed(() => [
  { index: 1, label: t('schedulePreferences.days.monday') || 'Lunedì' },
  { index: 2, label: t('schedulePreferences.days.tuesday') || 'Martedì' },
  { index: 3, label: t('schedulePreferences.days.wednesday') || 'Mercoledì' },
  { index: 4, label: t('schedulePreferences.days.thursday') || 'Giovedì' },
  { index: 5, label: t('schedulePreferences.days.friday') || 'Venerdì' },
  { index: 6, label: t('schedulePreferences.days.saturday') || 'Sabato' }
]);

const teacherGrid = ref({});

function initTeacherGrid() {
  const g = {};
  for (let d = 1; d <= 6; d++) {
    g[d] = {};
    for (let h = 1; h <= 8; h++) {
      g[d][h] = 'neutral';
    }
  }
  teacherGrid.value = g;
}

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
    teacherOptions.value = allTeachers.value;
  } catch (e) {
    console.error('Failed fetching teachers', e);
  }
}

function filterTeachers(val, update) {
  if (!val) {
    update(() => {
      teacherOptions.value = allTeachers.value;
    });
    return;
  }
  update(() => {
    const needle = val.toLowerCase();
    teacherOptions.value = allTeachers.value.filter(v => v.label.toLowerCase().includes(needle) || (v.email && v.email.toLowerCase().includes(needle)));
  });
}

async function onTeacherChanged(teacher) {
  if (!teacher) return;
  loadingTeacherPrefs.value = true;
  initTeacherGrid();
  try {
    const res = await timetableGenService.getPreferences({ teacher_id: teacher.id });
    const prefs = res.data || [];
    for (const p of prefs) {
      if (teacherGrid.value[p.day_of_week] && teacherGrid.value[p.day_of_week][p.hour_index] !== undefined) {
        teacherGrid.value[p.day_of_week][p.hour_index] = p.preference_type;
      }
    }
  } catch (err) {
    console.warn('Could not load teacher preferences', err);
  } finally {
    loadingTeacherPrefs.value = false;
  }
}

function cyclePreference(day, hour) {
  const current = teacherGrid.value[day]?.[hour] || 'neutral';
  let next = 'neutral';
  if (current === 'neutral') next = 'preferred';
  else if (current === 'preferred') next = 'unavailable';
  else if (current === 'unavailable') next = 'neutral';

  if (!teacherGrid.value[day]) teacherGrid.value[day] = {};
  teacherGrid.value[day][hour] = next;
}

function setQuickDayOff(dayIndex, dayLabel) {
  if (!teacherGrid.value[dayIndex]) teacherGrid.value[dayIndex] = {};
  for (let h = 1; h <= 8; h++) {
    teacherGrid.value[dayIndex][h] = 'unavailable';
  }
  $q.notify({
    type: 'info',
    message: t('timetableConstraints.dayOffSet', { day: dayLabel }) || `Giorno libero impostato su ${dayLabel}`
  });
}

function resetTeacherGridToNeutral() {
  initTeacherGrid();
}

function getCellBgClass(day, hour) {
  const state = teacherGrid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return 'bg-emerald-100 hover:bg-emerald-200';
  if (state === 'unavailable') return 'bg-rose-100 hover:bg-rose-200';
  return 'bg-white hover:bg-slate-100';
}

function getCellTextClass(day, hour) {
  const state = teacherGrid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return 'text-emerald-800';
  if (state === 'unavailable') return 'text-rose-800';
  return 'text-slate-400';
}

function getCellLabel(day, hour) {
  const state = teacherGrid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return t('schedulePreferences.cellPreferred') || 'Preferita';
  if (state === 'unavailable') return t('schedulePreferences.cellUnavailable') || 'Indisponibile';
  return t('schedulePreferences.cellNeutral') || 'Neutra';
}

async function saveCurrentTeacherPreferences() {
  if (!selectedTeacher.value) return;
  savingTeacherPrefs.value = true;
  try {
    const entries = [];
    for (let d = 1; d <= 6; d++) {
      for (let h = 1; h <= 8; h++) {
        const type = teacherGrid.value[d]?.[h] || 'neutral';
        entries.push({
          day_of_week: d,
          hour_index: h,
          preference_type: type
        });
      }
    }
    await timetableGenService.savePreferences({ preferences: entries }, { teacher_id: selectedTeacher.value.id });
    $q.notify({
      type: 'positive',
      message: t('timetableConstraints.saveTeacherSuccess') || 'Preferenze orario del docente salvate con successo'
    });
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('timetableConstraints.notifySaveError')
    });
  } finally {
    savingTeacherPrefs.value = false;
  }
}

// ----------------- TAB 3: Spezzoni Orari per Classe -----------------
const classes = ref([]);
const constraints = ref([]);
const loadingConstraints = ref(false);
const showSpezzoneModal = ref(false);

const spezzoneForm = ref({
  class_id: null,
  subject_id: null,
  hours_per_week: 2,
  placeholder_name: 'Docente da Nominare',
  notes: ''
});

const spezzoniList = computed(() => {
  return constraints.value
    .filter(c => c.constraint_type === 'unassigned_hours')
    .map(c => {
      let p = {};
      try {
        p = typeof c.parameters === 'string' ? JSON.parse(c.parameters) : c.parameters;
      } catch {
        p = {};
      }
      const matchedClass = classes.value.find(cl => cl.id === c.target_id);
      return {
        id: c.id,
        class_id: c.target_id,
        class_name: p.class_name || matchedClass?.label || 'Classe',
        subject_id: p.subject_id,
        subject_name: p.subject_name || 'Materia',
        hours_per_week: p.hours_per_week || 2,
        placeholder_name: p.placeholder_name || 'Docente da Nominare',
        notes: p.notes || ''
      };
    });
});

const spezzoniColumns = computed(() => [
  { name: 'class_name', label: t('timetableConstraints.colClass'), field: 'class_name', sortable: true, align: 'left' },
  { name: 'subject_name', label: t('timetableConstraints.colSubject'), field: 'subject_name', sortable: true, align: 'left' },
  { name: 'hours_per_week', label: t('timetableConstraints.colHours'), field: 'hours_per_week', align: 'center', sortable: true },
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
      label: `Classe ${c.name}${c.section || ''}`,
      name: c.name,
      section: c.section
    }));
  } catch (e) {
    console.error('Failed fetching classes', e);
  }
}

async function fetchConstraints() {
  loadingConstraints.value = true;
  try {
    const res = await timetableGenService.getConstraints();
    constraints.value = res.data || [];
  } catch (err) {
    console.warn('Could not load constraints', err);
  } finally {
    loadingConstraints.value = false;
  }
}

function openSpezzoneModal() {
  spezzoneForm.value = {
    class_id: classes.value[0]?.id || null,
    subject_id: subjects.value[0]?.id || null,
    hours_per_week: 2,
    placeholder_name: 'Docente da Nominare',
    notes: ''
  };
  onSpezzoneSubjectChanged(spezzoneForm.value.subject_id);
  showSpezzoneModal.value = true;
}

function onSpezzoneSubjectChanged(subjectId) {
  const subj = subjects.value.find(s => s.id === subjectId);
  if (subj) {
    spezzoneForm.value.placeholder_name = `Docente da Nominare (${subj.name})`;
  }
}

async function saveSpezzone() {
  if (!spezzoneForm.value.class_id || !spezzoneForm.value.subject_id || spezzoneForm.value.hours_per_week <= 0) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori della classe e della materia' });
    return;
  }
  const matchedClass = classes.value.find(c => c.id === spezzoneForm.value.class_id);
  const matchedSubject = subjects.value.find(s => s.id === spezzoneForm.value.subject_id);

  const payload = {
    constraint_type: 'unassigned_hours',
    target_type: 'class',
    target_id: spezzoneForm.value.class_id,
    parameters: {
      class_id: spezzoneForm.value.class_id,
      class_name: matchedClass?.label || 'Classe',
      subject_id: spezzoneForm.value.subject_id,
      subject_name: matchedSubject?.name || 'Materia',
      hours_per_week: spezzoneForm.value.hours_per_week,
      placeholder_name: spezzoneForm.value.placeholder_name || 'Docente da Nominare',
      notes: spezzoneForm.value.notes
    },
    is_hard: true,
    priority: 1,
    is_active: true
  };

  try {
    await timetableGenService.saveConstraint(payload);
    $q.notify({
      type: 'positive',
      message: t('timetableConstraints.spezzoneSaved') || 'Spezzone orario configurato con successo'
    });
    showSpezzoneModal.value = false;
    fetchConstraints();
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('timetableConstraints.notifySaveError')
    });
  }
}

function confirmDeleteSpezzone(row) {
  $q.dialog({
    title: t('timetableConstraints.removeSpezzoneTitle') || 'Rimuovi Spezzone',
    message: t('timetableConstraints.removeSpezzoneMsg', { hours: row.hours_per_week, subject: row.subject_name, class: row.class_name }) ||
             `Vuoi eliminare lo spezzone di ${row.hours_per_week} ore di ${row.subject_name} per la ${row.class_name}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await timetableGenService.deleteConstraint(row.id);
      $q.notify({
        type: 'positive',
        message: t('timetableConstraints.spezzoneRemoved') || 'Spezzone orario rimosso'
      });
      fetchConstraints();
    } catch {
      $q.notify({
        type: 'negative',
        message: t('timetableConstraints.notifyDeleteError')
      });
    }
  });
}

onMounted(() => {
  fetchRoomReqs();
  fetchSubjects();
  fetchTeachers();
  fetchClasses();
  fetchConstraints();
});
</script>
