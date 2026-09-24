<template>
  <q-page padding class="min-h-screen" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-slate-50 text-slate-800'">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg gap-4">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none row items-center gap-2" :class="$q.dark.isActive ? 'text-white' : 'text-slate-800'">
          <q-icon name="meeting_room" color="primary" size="36px" />
          {{ t('roomsManagement.title') }}
        </h1>
        <p class="text-subtitle1 q-mt-xs q-mb-none" :class="$q.dark.isActive ? 'text-grey-4' : 'text-slate-500'">
          {{ t('roomsManagement.subtitle') }}
        </p>
      </div>

      <div class="row items-center gap-2">
        <q-btn
          v-if="activeTab === 'rooms'"
          :label="t('roomsManagement.newRoomBtn')"
          icon="add_circle"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openRoomModal()"
        />
        <q-btn
          v-if="activeTab === 'buildings'"
          :label="t('roomsManagement.newBuildingBtn')"
          icon="add_business"
          color="secondary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openBuildingModal()"
        />
      </div>
    </div>

    <!-- Navigation Tabs -->
    <q-tabs
      v-model="activeTab"
      dense
      class="text-grey-7 q-mb-md rounded-xl"
      active-color="primary"
      indicator-color="primary"
      align="left"
      narrow-indicator
      :class="$q.dark.isActive ? 'bg-grey-9 text-white' : 'bg-white shadow-xs'"
    >
      <q-tab name="rooms" icon="door_front" :label="t('roomsManagement.tabRooms')" no-caps class="text-weight-bold q-py-sm" />
      <q-tab name="buildings" icon="apartment" :label="t('roomsManagement.tabBuildings')" no-caps class="text-weight-bold q-py-sm" />
      <q-tab name="bookings" icon="event_note" :label="t('roomsManagement.tabBookings')" no-caps class="text-weight-bold q-py-sm" />
    </q-tabs>

    <!-- Tab 1: Rooms & Labs -->
    <div v-if="activeTab === 'rooms'">
      <!-- Filter Bar -->
      <q-card flat bordered class="rounded-2xl q-pa-md q-mb-lg shadow-sm" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <div class="row items-center q-col-gutter-md">
          <div class="col-12 col-sm-4">
            <q-select
              v-model="selectedBuildingFilter"
              :options="buildingFilterOptions"
              option-value="id"
              option-label="name"
              emit-value
              map-options
              outlined
              dense
              clearable
              :label="t('roomsManagement.filterBuilding')"
              class="rounded-lg"
              @update:model-value="fetchRooms"
            >
              <template v-slot:prepend><q-icon name="apartment" color="primary" /></template>
            </q-select>
          </div>
          <div class="col-12 col-sm-4">
            <q-select
              v-model="selectedRoomTypeFilter"
              :options="roomTypeOptions"
              option-value="value"
              option-label="label"
              emit-value
              map-options
              outlined
              dense
              clearable
              :label="t('roomsManagement.filterRoomType')"
              class="rounded-lg"
              @update:model-value="fetchRooms"
            >
              <template v-slot:prepend><q-icon name="category" color="primary" /></template>
            </q-select>
          </div>
          <div class="col-12 col-sm-4 row items-center justify-end">
            <q-btn flat icon="refresh" :label="t('roomsManagement.refresh')" no-caps color="primary" @click="fetchRooms" />
          </div>
        </div>
      </q-card>

      <!-- Rooms Grid -->
      <div v-if="loadingRooms" class="row justify-center q-pa-xl">
        <q-spinner-dots size="48px" color="primary" />
      </div>

      <div v-else-if="rooms.length === 0" class="text-center q-pa-xl">
        <q-icon name="meeting_room" size="64px" color="grey-5" />
        <div class="text-h6 text-grey-6 q-mt-sm">{{ t('roomsManagement.noRoomsTitle') }}</div>
        <p class="text-grey-5">{{ t('roomsManagement.noRoomsDesc') }}</p>
        <q-btn color="primary" unelevated no-caps :label="t('roomsManagement.createRoomBtn')" icon="add" class="rounded-xl q-mt-sm" @click="openRoomModal()" />
      </div>

      <div v-else class="row q-col-gutter-md">
        <div v-for="rm in rooms" :key="rm.id" class="col-12 col-sm-6 col-md-4">
          <q-card flat bordered class="rounded-2xl h-full shadow-xs hover:shadow-md transition-shadow" :class="$q.dark.isActive ? 'bg-grey-9 border-grey-8' : 'bg-white border-slate-200'">
            <q-card-section>
              <div class="row items-center justify-between q-mb-xs">
                <q-badge :color="getRoomTypeBadgeColor(rm.room_type)" class="q-px-sm q-py-xs rounded-md text-weight-bold">
                  {{ getRoomTypeLabel(rm.room_type) }}
                </q-badge>
                <q-chip v-if="!rm.is_active" size="xs" color="negative" text-color="white" :label="t('roomsManagement.notActive')" />
                <q-chip v-else-if="rm.requires_booking" size="xs" color="info" text-color="white" :label="t('roomsManagement.requiresBooking')" />
                <q-chip v-else size="xs" color="positive" text-color="white" :label="t('roomsManagement.freeAccess')" />
              </div>

              <div class="text-h6 text-weight-bold q-mt-xs">{{ rm.name }}</div>

              <div class="text-caption text-grey-6 row items-center gap-1 q-mt-xs">
                <q-icon name="apartment" size="16px" />
                <span>{{ rm.building_name || t('roomsManagement.noBuildingAssociated') }}</span>
              </div>

              <div class="text-caption text-grey-6 row items-center gap-1 q-mt-xs">
                <q-icon name="people" size="16px" />
                <span>{{ t('roomsManagement.capacity', { count: rm.capacity }) }}</span>
              </div>

              <!-- Equipment Chips -->
              <div v-if="parseEquipment(rm.equipment).length > 0" class="q-mt-sm row gap-1 items-center">
                <q-chip
                  v-for="(eq, idx) in parseEquipment(rm.equipment)"
                  :key="idx"
                  size="xs"
                  dense
                  color="grey-3"
                  text-color="dark"
                  icon="build"
                >
                  {{ eq }}
                </q-chip>
              </div>

              <div v-if="rm.notes" class="text-caption text-grey-5 q-mt-sm italic">
                {{ rm.notes }}
              </div>
            </q-card-section>

            <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

            <q-card-actions align="between" class="q-px-md q-py-sm">
              <q-btn
                flat
                dense
                no-caps
                size="sm"
                color="secondary"
                icon="calendar_month"
                :label="t('roomsManagement.checkAvailability')"
                @click="openAvailabilityModal(rm)"
              />
              <div>
                <q-btn flat round dense size="sm" icon="edit" color="primary" :title="t('roomsManagement.edit')" @click="openRoomModal(rm)" />
                <q-btn flat round dense size="sm" icon="delete" color="negative" :title="t('roomsManagement.delete')" @click="confirmDeleteRoom(rm)" />
              </div>
            </q-card-actions>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Tab 2: Buildings (Plessi) -->
    <div v-if="activeTab === 'buildings'">
      <div v-if="loadingBuildings" class="row justify-center q-pa-xl">
        <q-spinner-dots size="48px" color="secondary" />
      </div>

      <div v-else-if="buildings.length === 0" class="text-center q-pa-xl">
        <q-icon name="apartment" size="64px" color="grey-5" />
        <div class="text-h6 text-grey-6 q-mt-sm">{{ t('roomsManagement.noBuildingsTitle') }}</div>
        <p class="text-grey-5">{{ t('roomsManagement.noBuildingsDesc') }}</p>
        <q-btn color="secondary" unelevated no-caps :label="t('roomsManagement.createBuildingBtn')" icon="add" class="rounded-xl q-mt-sm" @click="openBuildingModal()" />
      </div>

      <div v-else class="row q-col-gutter-md">
        <div v-for="b in buildings" :key="b.id" class="col-12 col-sm-6 col-md-4">
          <q-card flat bordered class="rounded-2xl h-full shadow-xs" :class="$q.dark.isActive ? 'bg-grey-9 border-grey-8' : 'bg-white border-slate-200'">
            <q-card-section>
              <div class="row items-center justify-between">
                <div class="text-h6 text-weight-bold">{{ b.name }}</div>
                <q-badge :color="b.is_active ? 'positive' : 'grey-6'">
                  {{ b.is_active ? t('roomsManagement.active') : t('roomsManagement.deactivated') }}
                </q-badge>
              </div>

              <div v-if="b.address" class="text-body2 text-grey-6 row items-center gap-1 q-mt-sm">
                <q-icon name="location_on" size="18px" color="primary" />
                <span>{{ b.address }}</span>
              </div>

              <div v-if="b.notes" class="text-caption text-grey-5 q-mt-sm">
                {{ b.notes }}
              </div>
            </q-card-section>

            <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

            <q-card-actions align="right" class="q-px-md q-py-sm">
              <q-btn flat round dense size="sm" icon="edit" color="primary" :title="t('roomsManagement.edit')" @click="openBuildingModal(b)" />
              <q-btn flat round dense size="sm" icon="delete" color="negative" :title="t('roomsManagement.delete')" @click="confirmDeleteBuilding(b)" />
            </q-card-actions>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Tab 3: Bookings Registry -->
    <div v-if="activeTab === 'bookings'">
      <q-card flat bordered class="rounded-2xl q-pa-md q-mb-lg shadow-sm" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <div class="row items-center q-col-gutter-md">
          <div class="col-12 col-sm-4">
            <q-input v-model="bookingFromDate" type="date" :label="t('roomsManagement.fromDate')" outlined dense class="rounded-lg" @update:model-value="fetchBookings" />
          </div>
          <div class="col-12 col-sm-4">
            <q-input v-model="bookingToDate" type="date" :label="t('roomsManagement.toDate')" outlined dense class="rounded-lg" @update:model-value="fetchBookings" />
          </div>
          <div class="col-12 col-sm-4 row items-center justify-end">
            <q-btn flat icon="refresh" :label="t('roomsManagement.refresh')" no-caps color="primary" @click="fetchBookings" />
          </div>
        </div>
      </q-card>

      <q-table
        flat
        bordered
        :rows="bookings"
        :columns="bookingColumns"
        row-key="id"
        :loading="loadingBookings"
        class="rounded-2xl shadow-xs"
        :no-data-label="t('roomsManagement.noBookings')"
      >
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="props.row.status === 'confirmed' ? 'positive' : 'grey-6'">
              {{ props.row.status === 'confirmed' ? (t('common.confirmed') || 'Confermata') : (t('common.cancelled') || 'Cancellata') }}
            </q-badge>
            <q-badge v-if="props.row.is_recurring" color="purple-6" class="q-ml-xs">
              {{ t('common.recurring') || 'Ricorrente' }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" align="right">
            <q-btn
              v-if="props.row.status === 'confirmed'"
              flat
              dense
              round
              color="negative"
              icon="cancel"
              :title="t('roomsManagement.cancelBooking')"
              @click="confirmCancelBooking(props.row)"
            />
          </q-td>
        </template>
      </q-table>
    </div>

    <!-- Modal: Create / Edit Room -->
    <q-dialog v-model="showRoomModal" persistent>
      <q-card style="min-width: 480px; max-width: 600px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold">
            {{ editingRoom ? t('roomsManagement.modalEditRoom') : t('roomsManagement.modalNewRoom') }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-input v-model="roomForm.name" :label="t('roomsManagement.roomName')" outlined dense :rules="[val => !!val || t('roomsManagement.nameRequired')]" />

          <q-select
            v-model="roomForm.building_id"
            :options="buildings"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            clearable
            :label="t('roomsManagement.roomBuildingSelect')"
          />

          <q-select
            v-model="roomForm.room_type"
            :options="roomTypeOptions"
            option-value="value"
            option-label="label"
            emit-value
            map-options
            outlined
            dense
            :label="t('roomsManagement.roomTypeLabel')"
          />

          <q-input v-model.number="roomForm.capacity" type="number" :label="t('roomsManagement.roomCapacity')" outlined dense min="1" />

          <div>
            <div class="text-caption text-grey-7 q-mb-xs">{{ t('roomsManagement.roomEquipment') }}</div>
            <q-select
              v-model="roomForm.equipment"
              use-input
              use-chips
              multiple
              hide-dropdown-icon
              new-value-mode="add-unique"
              outlined
              dense
              :label="t('roomsManagement.equipmentHint2')"
            />
          </div>

          <div class="row items-center justify-between">
            <q-toggle v-model="roomForm.requires_booking" :label="t('roomsManagement.roomRequiresBooking')" />
            <q-toggle v-if="editingRoom" v-model="roomForm.is_active" :label="t('roomsManagement.roomActive')" />
          </div>

          <q-input v-model="roomForm.notes" type="textarea" rows="2" :label="t('roomsManagement.roomNotes')" outlined dense />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel')" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="editingRoom ? (t('common.saveChanges') || 'Salva Modifiche') : t('roomsManagement.createRoomBtn')" no-caps class="rounded-xl q-px-md font-bold" @click="saveRoom" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Create / Edit Building -->
    <q-dialog v-model="showBuildingModal" persistent>
      <q-card style="min-width: 420px; max-width: 500px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold">
            {{ editingBuilding ? t('roomsManagement.modalEditBuilding') : t('roomsManagement.modalNewBuilding') }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-input v-model="buildingForm.name" :label="t('roomsManagement.buildingName')" outlined dense :rules="[val => !!val || t('roomsManagement.nameRequired')]" />
          <q-input v-model="buildingForm.address" :label="t('roomsManagement.buildingAddress')" outlined dense />
          <q-input v-model="buildingForm.notes" type="textarea" rows="2" :label="t('roomsManagement.buildingNotes')" outlined dense />
          <q-toggle v-if="editingBuilding" v-model="buildingForm.is_active" :label="t('roomsManagement.buildingActive')" />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel')" no-caps v-close-popup />
          <q-btn unelevated color="secondary" :label="editingBuilding ? (t('common.saveChanges') || 'Salva Modifiche') : t('roomsManagement.createBuildingBtn')" no-caps class="rounded-xl q-px-md font-bold" @click="saveBuilding" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Room Weekly Availability Grid -->
    <q-dialog v-model="showAvailabilityModal" max-width="850px">
      <q-card class="rounded-2xl q-pa-md" style="min-width: 750px;">
        <q-card-section class="row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold">{{ inspectingRoom?.name }}</div>
            <div class="text-caption text-grey-6">{{ inspectingRoom?.building_name }} • {{ t('roomsManagement.capacity', { count: inspectingRoom?.capacity }) }}</div>
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <div v-if="loadingAvailability" class="row justify-center q-pa-lg">
            <q-spinner-dots size="40px" color="primary" />
          </div>
          <div v-else class="overflow-x-auto">
            <table class="w-full border-collapse text-sm">
              <thead>
                <tr class="bg-slate-100 text-slate-700">
                  <th class="p-2 border border-slate-200">{{ t('roomsManagement.colHour') }}</th>
                  <th v-for="day in availabilityDays" :key="day.date" class="p-2 border border-slate-200 text-center">
                    {{ formatDayHeader(day.date) }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="hour in 8" :key="hour">
                  <td class="p-2 border border-slate-200 font-bold text-center bg-slate-50">{{ t('roomsManagement.hourSlot', { hour }) }}</td>
                  <td
                    v-for="day in availabilityDays"
                    :key="day.date + '-' + hour"
                    class="p-2 border border-slate-200 text-center"
                    :class="getSlotClass(getSlotData(day, hour))"
                  >
                    <div v-if="getSlotData(day, hour)?.is_available" class="text-positive text-weight-bold">
                      <q-icon name="check_circle" size="14px" /> {{ t('roomsManagement.free') }}
                    </div>
                    <div v-else class="text-negative text-xs">
                      <div class="font-bold">{{ getSlotData(day, hour)?.booking?.teacher_name ? t('roomsManagement.bookedBy', { name: getSlotData(day, hour)?.booking?.teacher_name }) : (t('common.occupied') || 'Occupata') }}</div>
                      <div class="text-grey-7">{{ getSlotData(day, hour)?.booking?.class_name }}</div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import roomsService from '@/services/roomsService';

const $q = useQuasar();
const { t, locale } = useI18n();

const activeTab = ref('rooms');

// Data
const buildings = ref([]);
const rooms = ref([]);
const bookings = ref([]);

const loadingRooms = ref(false);
const loadingBuildings = ref(false);
const loadingBookings = ref(false);

// Filters
const selectedBuildingFilter = ref(null);
const selectedRoomTypeFilter = ref(null);

// Date filters for bookings (current week)
const today = new Date();
const nextMonth = new Date(today.getTime() + 30 * 24 * 60 * 60 * 1000);
const bookingFromDate = ref(today.toISOString().split('T')[0]);
const bookingToDate = ref(nextMonth.toISOString().split('T')[0]);

// Modals
const showRoomModal = ref(false);
const editingRoom = ref(null);
const roomForm = ref({
  name: '',
  building_id: null,
  room_type: 'classroom',
  capacity: 30,
  equipment: [],
  requires_booking: true,
  is_active: true,
  notes: ''
});

const showBuildingModal = ref(false);
const editingBuilding = ref(null);
const buildingForm = ref({
  name: '',
  address: '',
  notes: '',
  is_active: true
});

// Availability Modal
const showAvailabilityModal = ref(false);
const inspectingRoom = ref(null);
const availabilityDays = ref([]);
const loadingAvailability = ref(false);

const roomTypeOptions = computed(() => [
  { label: t('roomsManagement.types.standard') || 'Aula Normale', value: 'classroom' },
  { label: t('roomsManagement.types.lab_info') || 'Laboratorio di Informatica', value: 'lab_informatica' },
  { label: t('roomsManagement.types.lab_science') || 'Laboratorio di Scienze', value: 'lab_scienze' },
  { label: t('roomsManagement.types.lab_lang') || 'Laboratorio di Lingue', value: 'lab_lingue' },
  { label: t('roomsManagement.types.lab_chemistry') || 'Laboratorio di Chimica', value: 'lab_chimica' },
  { label: t('roomsManagement.types.lab_physics') || 'Laboratorio di Fisica', value: 'lab_fisica' },
  { label: t('roomsManagement.types.art') || 'Laboratorio di Arte', value: 'lab_arte' },
  { label: t('roomsManagement.types.gym') || 'Palestra', value: 'palestra' },
  { label: t('roomsManagement.types.auditorium') || 'Aula Magna', value: 'aula_magna' },
  { label: t('roomsManagement.types.library') || 'Biblioteca', value: 'biblioteca' },
  { label: t('roomsManagement.types.other') || 'Altro', value: 'altro' }
]);

const buildingFilterOptions = computed(() => {
  return buildings.value.map(b => ({ id: b.id, name: b.name }));
});

const bookingColumns = computed(() => [
  { name: 'date', label: t('roomsManagement.colDate'), field: 'booking_date', sortable: true, align: 'left' },
  { name: 'hour', label: t('roomsManagement.colHour'), field: 'hour_index', align: 'center', format: v => t('roomsManagement.hourSlot', { hour: v }) },
  { name: 'room', label: t('roomsManagement.colRoom'), field: 'room_name', sortable: true, align: 'left' },
  { name: 'building', label: t('roomsManagement.colBuilding'), field: 'building_name', align: 'left' },
  { name: 'teacher', label: t('roomsManagement.colTeacher'), field: 'teacher_name', sortable: true, align: 'left' },
  { name: 'class', label: t('common.class') || 'Classe', field: 'class_name', align: 'left' },
  { name: 'status', label: t('common.status') || 'Stato', field: 'status', align: 'center' },
  { name: 'actions', label: t('roomsManagement.colActions'), field: 'id', align: 'right' }
]);

// Fetchers
async function fetchBuildings() {
  loadingBuildings.value = true;
  try {
    const res = await roomsService.getBuildings();
    buildings.value = res.data || [];
  } catch (err) {
    $q.notify({ type: 'negative', message: t('roomsManagement.notifyFetchBuildingsError') });
  } finally {
    loadingBuildings.value = false;
  }
}

async function fetchRooms() {
  loadingRooms.value = true;
  try {
    const params = {};
    if (selectedBuildingFilter.value) params.building_id = selectedBuildingFilter.value;
    if (selectedRoomTypeFilter.value) params.room_type = selectedRoomTypeFilter.value;
    const res = await roomsService.getRooms(params);
    rooms.value = res.data || [];
  } catch (err) {
    $q.notify({ type: 'negative', message: t('roomsManagement.notifyFetchRoomsError') });
  } finally {
    loadingRooms.value = false;
  }
}

async function fetchBookings() {
  loadingBookings.value = true;
  try {
    const res = await roomsService.getBookings({
      from: bookingFromDate.value,
      to: bookingToDate.value
    });
    bookings.value = res.data || [];
  } catch (err) {
    $q.notify({ type: 'negative', message: t('roomsManagement.notifyFetchBookingsError') });
  } finally {
    loadingBookings.value = false;
  }
}

// Helpers
function getRoomTypeLabel(type) {
  const found = roomTypeOptions.value.find(o => o.value === type);
  return found ? found.label : type;
}

function getRoomTypeBadgeColor(type) {
  if (type === 'palestra') return 'orange-8';
  if (type?.startsWith('lab_')) return 'indigo-7';
  if (type === 'aula_magna') return 'deep-purple-7';
  return 'teal-7';
}

function parseEquipment(eq) {
  if (!eq) return [];
  if (Array.isArray(eq)) return eq;
  try {
    return JSON.parse(eq);
  } catch {
    return [];
  }
}

// Actions Building
function openBuildingModal(building = null) {
  editingBuilding.value = building;
  if (building) {
    buildingForm.value = { ...building };
  } else {
    buildingForm.value = { name: '', address: '', notes: '', is_active: true };
  }
  showBuildingModal.value = true;
}

async function saveBuilding() {
  if (!buildingForm.value.name) return;
  try {
    if (editingBuilding.value) {
      await roomsService.updateBuilding(editingBuilding.value.id, buildingForm.value);
      $q.notify({ type: 'positive', message: t('roomsManagement.notifyBuildingSaved') });
    } else {
      await roomsService.createBuilding(buildingForm.value);
      $q.notify({ type: 'positive', message: t('roomsManagement.notifyBuildingSaved') });
    }
    showBuildingModal.value = false;
    fetchBuildings();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore nel salvataggio' });
  }
}

function confirmDeleteBuilding(b) {
  $q.dialog({
    title: t('roomsManagement.deleteBuildingTitle'),
    message: t('roomsManagement.deleteBuildingConfirm', { name: b.name }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await roomsService.deleteBuilding(b.id);
      $q.notify({ type: 'positive', message: t('roomsManagement.notifyBuildingDeleted') });
      fetchBuildings();
    } catch {
      $q.notify({ type: 'negative', message: t('roomsManagement.notifyBuildingError') });
    }
  });
}

// Actions Room
function openRoomModal(room = null) {
  editingRoom.value = room;
  if (room) {
    roomForm.value = {
      ...room,
      equipment: parseEquipment(room.equipment)
    };
  } else {
    roomForm.value = {
      name: '',
      building_id: buildings.value[0]?.id || null,
      room_type: 'classroom',
      capacity: 30,
      equipment: [],
      requires_booking: true,
      is_active: true,
      notes: ''
    };
  }
  showRoomModal.value = true;
}

async function saveRoom() {
  if (!roomForm.value.name) return;
  try {
    if (editingRoom.value) {
      await roomsService.updateRoom(editingRoom.value.id, roomForm.value);
      $q.notify({ type: 'positive', message: t('roomsManagement.notifyRoomSaved') });
    } else {
      await roomsService.createRoom(roomForm.value);
      $q.notify({ type: 'positive', message: t('roomsManagement.notifyRoomSaved') });
    }
    showRoomModal.value = false;
    fetchRooms();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore nel salvataggio' });
  }
}

function confirmDeleteRoom(rm) {
  $q.dialog({
    title: t('roomsManagement.deleteRoomTitle'),
    message: t('roomsManagement.deleteRoomConfirm', { name: rm.name }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await roomsService.deleteRoom(rm.id);
      $q.notify({ type: 'positive', message: t('roomsManagement.notifyRoomDeleted') });
      fetchRooms();
    } catch {
      $q.notify({ type: 'negative', message: t('roomsManagement.notifyRoomError') });
    }
  });
}

// Availability
async function openAvailabilityModal(room) {
  inspectingRoom.value = room;
  showAvailabilityModal.value = true;
  loadingAvailability.value = true;

  const start = new Date();
  const end = new Date(start.getTime() + 5 * 24 * 60 * 60 * 1000);
  const fromStr = start.toISOString().split('T')[0];
  const toStr = end.toISOString().split('T')[0];

  try {
    const res = await roomsService.getRoomAvailability(room.id, fromStr, toStr);
    availabilityDays.value = res.data || [];
  } catch {
    $q.notify({ type: 'negative', message: t('roomsManagement.notifyAvailabilityError') });
  } finally {
    loadingAvailability.value = false;
  }
}

function formatDayHeader(dateStr) {
  const d = new Date(dateStr);
  return d.toLocaleDateString(locale.value || 'it-IT', { weekday: 'short', day: 'numeric', month: 'short' });
}

function getSlotData(day, hourIndex) {
  if (!day || !day.slots) return null;
  return day.slots.find(s => s.hour_index === hourIndex);
}

function getSlotClass(slot) {
  if (!slot) return '';
  return slot.is_available ? 'bg-emerald-50' : 'bg-rose-50';
}

// Cancel Booking
function confirmCancelBooking(booking) {
  $q.dialog({
    title: t('roomsManagement.cancelBookingTitle'),
    message: booking.is_recurring
      ? t('roomsManagement.cancelRecurringPrompt')
      : t('roomsManagement.cancelConfirmPrompt', { date: booking.booking_date, room: booking.room_name }),
    options: booking.is_recurring ? {
      type: 'radio',
      model: 'single',
      items: [
        { label: t('roomsManagement.cancelSingleOnly'), value: 'single' },
        { label: t('roomsManagement.cancelWholeSeries'), value: 'series' }
      ]
    } : undefined,
    cancel: true,
    persistent: true
  }).onOk(async (choice) => {
    const cancelSeries = choice === 'series';
    try {
      await roomsService.cancelBooking(booking.id, cancelSeries);
      $q.notify({ type: 'positive', message: t('roomsManagement.notifyBookingCancelled') });
      fetchBookings();
    } catch {
      $q.notify({ type: 'negative', message: t('roomsManagement.notifyBookingCancelError') });
    }
  });
}

onMounted(() => {
  fetchBuildings();
  fetchRooms();
  fetchBookings();
});
</script>
