<template>
  <q-page class="q-pa-md parent-booking-page">
    <!-- Header Banner -->
    <div class="page-header q-mb-lg flex justify-between items-center wrap gap-md">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none flex items-center gap-sm">
          <q-icon name="confirmation_number" color="primary" />
          {{ $t('generalMeeting.parentTitle') }}
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none q-mt-xs">
          {{ $t('generalMeeting.parentSubtitle') }}
        </p>
      </div>
    </div>

    <!-- Active Parent Tickets Section -->
    <div v-if="myTickets.length > 0" class="q-mb-lg">
      <div class="text-h6 text-weight-bold text-primary q-mb-sm flex items-center gap-xs">
        <q-icon name="local_activity" />
        {{ $t('generalMeeting.myActiveTickets') }}
      </div>
      <div class="row q-col-gutter-md">
        <div v-for="ticket in myTickets" :key="ticket.id" class="col-12 col-md-6 col-lg-4">
          <q-card class="ticket-card shadow-2 rounded-borders">
            <q-card-section class="bg-primary text-white flex justify-between items-center">
              <div>
                <div class="text-caption text-blue-2">{{ ticket.meeting_title || $t('generalMeeting.generalMeetingFallback') }}</div>
                <div class="text-h4 text-weight-bolder">#{{ ticket.ticket_number }}</div>
              </div>
              <q-badge
                :color="ticket.status === 'in_colloquio' ? 'positive' : 'white'"
                :text-color="ticket.status === 'in_colloquio' ? 'white' : 'primary'"
                class="text-weight-bold"
                :label="$t(`generalMeeting.ticketStatus.${ticket.status}`)"
              />
            </q-card-section>

            <q-card-section>
              <div class="text-subtitle1 text-weight-bold">{{ ticket.teacher_name }}</div>
              <div class="text-caption text-grey-7">{{ $t('generalMeeting.station') }} {{ ticket.room_or_table || $t('generalMeeting.stationFallback') }}</div>
              <div class="text-caption text-primary q-mt-xs">{{ $t('generalMeeting.estimatedCallTime') }} <strong>{{ ticket.scheduled_time }}</strong></div>
              <div class="text-caption text-grey-8 q-mt-xs">{{ $t('generalMeeting.studentLabel') }} {{ ticket.student_name }}</div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Skeleton Loader while loading meetings -->
    <div v-if="loading && meetings.length === 0" class="row q-col-gutter-md q-mb-lg">
      <div v-for="n in 3" :key="n" class="col-12 col-md-4">
        <q-skeleton type="rect" height="160px" class="rounded-borders" />
      </div>
    </div>

    <!-- Available General Meetings & Teachers to Book -->
    <div class="text-h6 text-weight-bold text-primary q-mb-sm flex items-center gap-xs">
      <q-icon name="event_available" />
      {{ $t('generalMeeting.availableEvents') }}
    </div>

    <div v-if="meetings.length === 0 && !loading" class="q-pa-lg text-center text-grey-6 text-italic bg-white rounded-borders shadow-1">
      {{ $t('generalMeeting.noEventsAvailable') }}
    </div>

    <div v-for="meeting in meetings" :key="meeting.id" class="q-mb-lg">
      <q-card class="shadow-2 rounded-borders">
        <q-card-section class="bg-blue-1 flex justify-between items-center wrap gap-sm">
          <div>
            <div class="text-h6 text-weight-bold text-primary">{{ meeting.title }}</div>
            <div class="text-caption text-grey-8">
              {{ $t('generalMeeting.dateLabel') }} <strong>{{ meeting.event_date }}</strong> | {{ $t('generalMeeting.timeRangeLabel') }} <strong>{{ meeting.start_time }} - {{ meeting.end_time }}</strong> | {{ $t('generalMeeting.slotDurationLabel') }} <strong>{{ meeting.slot_duration_minutes }} min</strong>
            </div>
          </div>
          <q-badge color="primary" :label="meeting.location_type === 'in_presenza' ? $t('generalMeeting.inPerson') : $t('generalMeeting.online')" />
        </q-card-section>

        <q-separator />

        <!-- Teachers Grid in this Meeting -->
        <q-card-section>
          <div class="text-subtitle2 text-weight-bold text-grey-8 q-mb-sm">{{ $t('generalMeeting.selectTeacher') }}:</div>
          <div class="row q-col-gutter-sm">
            <div
              v-for="slot in meeting.teacher_slots || []"
              :key="slot.id"
              class="col-12 col-md-4"
            >
              <q-card flat bordered class="rounded-borders hover-elevate">
                <q-card-section>
                  <div class="text-weight-bold text-primary">{{ slot.teacher_name }}</div>
                  <div class="text-caption text-grey-7">{{ slot.subject_name || $t('generalMeeting.classTeacher') }}</div>
                  <div class="text-caption text-grey-6 q-mt-xs">{{ $t('generalMeeting.station') }} {{ slot.room_or_table || $t('generalMeeting.stationFallback') }}</div>
                  <div class="text-caption text-secondary q-mt-xs">{{ $t('generalMeeting.booked') }} {{ slot.booked_count }} / {{ slot.max_bookings }}</div>
                </q-card-section>
                <q-card-actions align="right">
                  <q-btn
                    color="primary"
                    icon="add_circle_outline"
                    :label="$t('generalMeeting.bookTicketBtn')"
                    dense
                    unelevated
                    no-caps
                    @click="openBookingModal(meeting, slot)"
                  />
                </q-card-actions>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- Booking Confirmation Dialog -->
    <q-dialog v-model="bookingDialog" persistent>
      <q-card style="width: min(500px, 95vw); max-width: 95vw;" class="rounded-borders">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-primary">{{ $t('generalMeeting.confirmBookingTitle') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-form @submit.prevent="submitBookingTicket" class="q-gutter-md">
            <div class="bg-grey-2 q-pa-sm rounded-borders">
              <div class="text-weight-bold text-primary">{{ activeBookingSlot?.teacher_name }}</div>
              <div class="text-caption">{{ activeBookingMeeting?.title }} ({{ activeBookingMeeting?.event_date }})</div>
            </div>

            <q-select
              v-model="bookingForm.student_id"
              :options="childrenOptions"
              emit-value
              map-options
              :label="$t('generalMeeting.form.selectStudent')"
              outlined
              dense
              required
            />

            <q-input
              v-model="bookingForm.notes"
              type="textarea"
              rows="2"
              :label="$t('generalMeeting.form.notes')"
              :placeholder="$t('generalMeeting.notesPlaceholder')"
              outlined
              dense
            />

            <div class="flex justify-end q-mt-md gap-sm">
              <q-btn flat :label="$t('common.cancel')" v-close-popup no-caps />
              <q-btn color="primary" :label="$t('generalMeeting.confirmTicket')" type="submit" :loading="bookingLoading" unelevated no-caps />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import colloquiService from 'src/services/colloquiService'
import api from 'src/services/api'

const { t } = useI18n()
const $q = useQuasar()

const meetings = ref([])
const myTickets = ref([])
const childrenOptions = ref([])
const loading = ref(false)

const bookingDialog = ref(false)
const bookingLoading = ref(false)
const activeBookingMeeting = ref(null)
const activeBookingSlot = ref(null)

const bookingForm = ref({
  meeting_id: '',
  teacher_id: '',
  student_id: '',
  notes: ''
})

function extractList(response) {
  if (!response) return []
  const data = response.data !== undefined ? response.data : response
  if (Array.isArray(data)) return data
  if (data && Array.isArray(data.meetings)) return data.meetings
  if (data && Array.isArray(data.children)) return data.children
  if (data && Array.isArray(data.tickets)) return data.tickets
  return []
}

async function loadData() {
  loading.value = true
  try {
    const [meetingsRes, childrenRes] = await Promise.allSettled([
      colloquiService.listGeneralMeetings(),
      api.get('/users/me/children')
    ])

    const rawMeetings = meetingsRes.status === 'fulfilled' ? extractList(meetingsRes.value) : []
    const detailedMeetings = await Promise.all(
      rawMeetings.map(async m => {
        try {
          const detail = await colloquiService.getGeneralMeeting(m.id)
          return detail.data || m
        } catch {
          return m
        }
      })
    )
    meetings.value = detailedMeetings

    const rawChildren = childrenRes.status === 'fulfilled' ? extractList(childrenRes.value) : []
    childrenOptions.value = rawChildren.map(c => ({
      label: `${c.first_name || ''} ${c.last_name || ''}`.trim() || c.id,
      value: c.id
    }))

    // Load active tickets for parent
    if (detailedMeetings.length > 0) {
      const ticketsPromises = detailedMeetings.map(m => colloquiService.listQueueTickets(m.id))
      const ticketsResponses = await Promise.allSettled(ticketsPromises)
      myTickets.value = ticketsResponses
        .filter(r => r.status === 'fulfilled')
        .flatMap(r => extractList(r.value))
    }
  } catch (err) {
    console.error('Error loading parent general meeting data', err)
  } finally {
    loading.value = false
  }
}

function openBookingModal(meeting, slot) {
  activeBookingMeeting.value = meeting
  activeBookingSlot.value = slot
  bookingForm.value = {
    meeting_id: meeting.id,
    teacher_id: slot.teacher_id,
    student_id: childrenOptions.value.length > 0 ? childrenOptions.value[0].value : '',
    notes: ''
  }
  bookingDialog.value = true
}

async function submitBookingTicket() {
  bookingLoading.value = true
  try {
    const res = await colloquiService.bookQueueTicket(bookingForm.value)
    const ticket = res.data
    $q.notify({
      type: 'positive',
      message: t('generalMeeting.ticketBooked', { number: ticket.ticket_number, time: ticket.scheduled_time })
    })
    bookingDialog.value = false
    loadData()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('generalMeeting.ticketBookingError') })
  } finally {
    bookingLoading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.parent-booking-page {
  max-width: 1200px;
  margin: 0 auto;
}
.ticket-card {
  border-left: 5px solid var(--q-primary);
}
.hover-elevate {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.hover-elevate:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}
</style>
