<template>
  <q-page padding class="min-h-screen" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-slate-50 text-slate-800'">
    <!-- Top Header -->
    <div class="row items-center justify-between q-mb-lg gap-4">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none row items-center gap-2" :class="$q.dark.isActive ? 'text-white' : 'text-slate-800'">
          <q-icon name="edit_calendar" color="primary" size="36px" />
          {{ t('roomBooking.title') }}
        </h1>
        <p class="text-subtitle1 q-mt-xs q-mb-none" :class="$q.dark.isActive ? 'text-grey-4' : 'text-slate-500'">
          {{ t('roomBooking.subtitle') }}
        </p>
      </div>

      <div>
        <q-btn
          :label="t('roomBooking.newBookingBtn')"
          icon="add"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="showBookingModal = true"
        />
      </div>
    </div>

    <!-- Active Bookings List -->
    <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
      <q-card-section class="row items-center justify-between">
        <div class="text-h6 text-weight-bold row items-center gap-2">
          <q-icon name="bookmark" color="secondary" />
          {{ t('roomBooking.myBookingsTitle') }}
        </div>
        <q-btn flat dense icon="refresh" :label="t('roomBooking.refresh')" color="primary" no-caps @click="fetchMyBookings" />
      </q-card-section>

      <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

      <div v-if="loadingMyBookings" class="row justify-center q-pa-xl">
        <q-spinner-dots size="40px" color="primary" />
      </div>

      <div v-else-if="myBookings.length === 0" class="text-center q-pa-xl">
        <q-icon name="event_busy" size="56px" color="grey-5" />
        <div class="text-h6 text-grey-6 q-mt-sm">{{ t('roomBooking.noBookingsTitle') }}</div>
        <p class="text-grey-5">{{ t('roomBooking.noBookingsDesc') }}</p>
        <q-btn color="primary" unelevated no-caps :label="t('roomBooking.bookNowBtn')" icon="add" class="rounded-xl q-mt-sm" @click="showBookingModal = true" />
      </div>

      <q-table
        v-else
        flat
        :rows="myBookings"
        :columns="bookingColumns"
        row-key="id"
        class="no-shadow"
      >
        <template v-slot:body-cell-hour_index="props">
          <q-td :props="props" align="center">
            <q-badge color="blue-7" class="q-px-sm q-py-xs font-bold">
              {{ t('roomsManagement.hourSlot', { hour: props.row.hour_index }) }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="props.row.status === 'confirmed' ? 'positive' : 'grey-6'">
              {{ props.row.status === 'confirmed' ? (t('common.confirmed') || 'Confermata') : (t('common.cancelled') || 'Cancellata') }}
            </q-badge>
            <q-badge v-if="props.row.is_recurring" color="purple-6" class="q-ml-xs">
              {{ t('common.recurring') || 'Ricorrente' }} ({{ props.row.recurrence_pattern === 'weekly' ? t('roomBooking.recurrenceWeekly') : t('roomBooking.recurrenceBiweekly') }})
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" align="right">
            <q-btn
              v-if="props.row.status === 'confirmed'"
              flat
              round
              dense
              color="negative"
              icon="cancel"
              :title="t('roomBooking.cancelBookingTitle')"
              @click="confirmCancel(props.row)"
            />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Modal: New Booking -->
    <q-dialog v-model="showBookingModal" persistent>
      <q-card style="min-width: 500px; max-width: 600px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold row items-center gap-2">
            <q-icon name="meeting_room" color="primary" />
            {{ t('roomBooking.modalTitle') }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <!-- Building Filter & Room Selection -->
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-select
                v-model="filterBuilding"
                :options="buildings"
                option-value="id"
                option-label="name"
                emit-value
                map-options
                outlined
                dense
                clearable
                :label="t('roomBooking.filterBuilding')"
                @update:model-value="onBuildingChange"
              />
            </div>
            <div class="col-12 col-sm-6">
              <q-select
                v-model="bookingForm.room_id"
                :options="filteredRooms"
                option-value="id"
                option-label="name"
                emit-value
                map-options
                outlined
                dense
                :label="t('roomBooking.selectRoomPlaceholder')"
                :rules="[val => !!val || t('roomBooking.selectRoomRequired')]"
              >
                <template v-slot:option="scope">
                  <q-item v-bind="scope.itemProps">
                    <q-item-section>
                      <q-item-label>{{ scope.opt.name }}</q-item-label>
                      <q-item-label caption>{{ scope.opt.building_name }} • {{ t('roomsManagement.capacity', { count: scope.opt.capacity }) }}</q-item-label>
                    </q-item-section>
                  </q-item>
                </template>
              </q-select>
            </div>
          </div>

          <!-- Class and Subject -->
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-select
                v-model="bookingForm.class_id"
                :options="classes"
                option-value="id"
                option-label="name"
                emit-value
                map-options
                outlined
                dense
                clearable
                :label="t('roomBooking.classLabel')"
              />
            </div>
            <div class="col-12 col-sm-6">
              <q-select
                v-model="bookingForm.subject_id"
                :options="subjects"
                option-value="id"
                option-label="name"
                emit-value
                map-options
                outlined
                dense
                clearable
                :label="t('roomBooking.subjectLabel')"
              />
            </div>
          </div>

          <!-- Date & Hour -->
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-7">
              <q-input
                v-model="bookingForm.booking_date"
                type="date"
                outlined
                dense
                :label="t('roomBooking.dateLabel')"
                :rules="[val => !!val || t('roomBooking.dateRequired')]"
              />
            </div>
            <div class="col-12 col-sm-5">
              <q-select
                v-model="bookingForm.hour_index"
                :options="hourOptions"
                emit-value
                map-options
                outlined
                dense
                :label="t('roomBooking.hourLabel')"
                :rules="[val => !!val || t('roomBooking.hourRequired')]"
              />
            </div>
          </div>

          <!-- Recurring Booking Options -->
          <q-card flat bordered class="rounded-xl q-pa-sm bg-slate-50 border-slate-200">
            <div class="row items-center justify-between">
              <span class="text-weight-bold text-subtitle2">{{ t('roomBooking.recurringBooking') }}</span>
              <q-toggle v-model="bookingForm.is_recurring" color="purple-6" />
            </div>

            <div v-if="bookingForm.is_recurring" class="q-mt-sm q-gutter-sm">
              <q-select
                v-model="bookingForm.recurrence_pattern"
                :options="[
                  { label: t('roomBooking.recurrenceWeekly'), value: 'weekly' },
                  { label: t('roomBooking.recurrenceBiweekly'), value: 'biweekly' }
                ]"
                emit-value
                map-options
                outlined
                dense
                :label="t('roomBooking.recurrenceFrequency')"
              />

              <q-input
                v-model="bookingForm.recurring_until"
                type="date"
                outlined
                dense
                :label="t('roomBooking.untilDate')"
                :rules="[val => !!val || t('roomBooking.endDateRequired')]"
              />
            </div>
          </q-card>

          <q-input v-model="bookingForm.notes" type="textarea" rows="2" :label="t('roomBooking.notesOptional')" outlined dense />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel')" no-caps v-close-popup />
          <q-btn
            unelevated
            color="primary"
            :label="t('roomBooking.confirmBookingBtn')"
            no-caps
            class="rounded-xl q-px-md font-bold"
            :loading="submittingBooking"
            @click="submitBooking"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import roomsService from '@/services/roomsService';
import api from '@/services/api';

const $q = useQuasar();
const { t } = useI18n();

const myBookings = ref([]);
const loadingMyBookings = ref(false);
const submittingBooking = ref(false);

const rooms = ref([]);
const buildings = ref([]);
const classes = ref([]);
const subjects = ref([]);

const filterBuilding = ref(null);
const showBookingModal = ref(false);

const todayStr = new Date().toISOString().split('T')[0];
const oneMonthLaterStr = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];

const bookingForm = ref({
  room_id: null,
  class_id: null,
  subject_id: null,
  booking_date: todayStr,
  hour_index: 1,
  notes: '',
  is_recurring: false,
  recurrence_pattern: 'weekly',
  recurring_until: oneMonthLaterStr
});

const hourOptions = computed(() => Array.from({ length: 8 }, (_, i) => ({
  label: t('roomsManagement.hourSlot', { hour: i + 1 }),
  value: i + 1
})));

const bookingColumns = computed(() => [
  { name: 'booking_date', label: t('roomBooking.colDate'), field: 'booking_date', sortable: true, align: 'left' },
  { name: 'hour_index', label: t('roomBooking.colHour'), field: 'hour_index', align: 'center', format: v => t('roomsManagement.hourSlot', { hour: v }) },
  { name: 'room_name', label: t('roomBooking.colRoom'), field: 'room_name', sortable: true, align: 'left' },
  { name: 'building', label: t('roomBooking.colBuilding'), field: 'building_name', align: 'left' },
  { name: 'class_name', label: t('common.class') || 'Classe', field: 'class_name', align: 'left' },
  { name: 'status', label: t('roomBooking.colStatus'), field: 'status', align: 'center' },
  { name: 'actions', label: t('roomBooking.colActions'), field: 'id', align: 'right' }
]);

const filteredRooms = computed(() => {
  if (!filterBuilding.value) return rooms.value;
  return rooms.value.filter(r => r.building_id === filterBuilding.value);
});

async function fetchMyBookings() {
  loadingMyBookings.value = true;
  try {
    const res = await roomsService.getBookings();
    myBookings.value = res.data || [];
  } catch {
    $q.notify({ type: 'negative', message: t('roomBooking.notifyFetchError') });
  } finally {
    loadingMyBookings.value = false;
  }
}

async function loadInitialData() {
  try {
    const [bRes, rRes, cRes, sRes] = await Promise.allSettled([
      roomsService.getBuildings(),
      roomsService.getRooms({ active_only: true }),
      api.get('/classes'),
      api.get('/subjects')
    ]);

    if (bRes.status === 'fulfilled') buildings.value = bRes.value.data || [];
    if (rRes.status === 'fulfilled') rooms.value = rRes.value.data || [];
    if (cRes.status === 'fulfilled') classes.value = cRes.value.data || [];
    if (sRes.status === 'fulfilled') subjects.value = sRes.value.data || [];
  } catch (err) {
    console.error('Error loading initial data', err);
  }
}

function onBuildingChange() {
  bookingForm.value.room_id = null;
}

async function submitBooking() {
  if (!bookingForm.value.room_id || !bookingForm.value.booking_date || !bookingForm.value.hour_index) {
    $q.notify({ type: 'warning', message: t('roomBooking.notifyFillRequired') });
    return;
  }

  submittingBooking.value = true;
  try {
    const payload = {
      room_id: bookingForm.value.room_id,
      class_id: bookingForm.value.class_id,
      subject_id: bookingForm.value.subject_id,
      booking_date: bookingForm.value.booking_date,
      hour_index: bookingForm.value.hour_index,
      notes: bookingForm.value.notes,
      is_recurring: bookingForm.value.is_recurring,
      recurrence_pattern: bookingForm.value.is_recurring ? bookingForm.value.recurrence_pattern : null,
      recurring_until: bookingForm.value.is_recurring ? bookingForm.value.recurring_until : null
    };

    const res = await roomsService.createBooking(payload);
    $q.notify({
      type: 'positive',
      message: res.data?.message || t('roomBooking.notifyBookingCreated')
    });
    showBookingModal.value = false;
    fetchMyBookings();
  } catch (err) {
    const errMsg = err.response?.data?.error || 'Errore nella prenotazione';
    $q.notify({ type: 'negative', message: errMsg });
  } finally {
    submittingBooking.value = false;
  }
}

function confirmCancel(booking) {
  $q.dialog({
    title: t('roomBooking.cancelBookingTitle'),
    message: booking.is_recurring
      ? t('roomBooking.recurringCancelPrompt')
      : t('roomBooking.cancelConfirmPrompt', { date: booking.booking_date }),
    options: booking.is_recurring ? {
      type: 'radio',
      model: 'single',
      items: [
        { label: t('roomBooking.singleOnly'), value: 'single' },
        { label: t('roomBooking.wholeSeries'), value: 'series' }
      ]
    } : undefined,
    cancel: true,
    persistent: true
  }).onOk(async (choice) => {
    const cancelSeries = choice === 'series';
    try {
      await roomsService.cancelBooking(booking.id, cancelSeries);
      $q.notify({ type: 'positive', message: t('roomBooking.notifyBookingCancelled') });
      fetchMyBookings();
    } catch {
      $q.notify({ type: 'negative', message: t('roomBooking.notifyCancelError') });
    }
  });
}

onMounted(() => {
  fetchMyBookings();
  loadInitialData();
});
</script>
