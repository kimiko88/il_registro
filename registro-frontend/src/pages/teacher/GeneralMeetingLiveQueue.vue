<template>
  <q-page class="q-pa-md meeting-queue-page">
    <!-- Header Banner -->
    <div class="page-header q-mb-lg flex justify-between items-center wrap gap-md">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none flex items-center gap-sm">
          <q-icon name="meeting_room" color="primary" />
          {{ $t('generalMeeting.teacherQueueTitle') }}
        </h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none q-mt-xs">
          {{ $t('generalMeeting.teacherQueueSubtitle') }}
        </p>
      </div>

      <div class="flex gap-sm items-center">
        <q-select
          v-model="selectedMeetingId"
          :options="meetingOptions"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          style="min-width: 250px;"
          :label="$t('generalMeeting.selectMeeting')"
          @update:model-value="loadQueueTickets"
        />
        <q-btn
          color="primary"
          icon="refresh"
          round
          flat
          :loading="loadingTickets"
          @click="loadQueueTickets"
        >
          <q-tooltip>{{ $t('common.refresh') }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Active Ticket Spotlight Monitor -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-4">
        <q-card class="current-ticket-card shadow-3 rounded-borders bg-primary text-white">
          <q-card-section>
            <div class="text-caption text-uppercase text-blue-2">{{ $t('generalMeeting.currentlyInMeeting') }}</div>
            <div v-if="currentMeetingTicket" class="q-mt-sm">
              <div class="text-h3 text-weight-bolder">#{{ currentMeetingTicket.ticket_number }}</div>
              <div class="text-h6 q-mt-xs">{{ currentMeetingTicket.parent_name }}</div>
              <div class="text-caption text-blue-1">{{ $t('generalMeeting.studentLabel') }} {{ currentMeetingTicket.student_name }}</div>
              <div class="text-caption text-blue-2 q-mt-xs">{{ $t('generalMeeting.timeLabel') }} {{ currentMeetingTicket.scheduled_time }}</div>
            </div>
            <div v-else class="q-py-md text-center text-blue-2 text-italic">
              {{ $t('generalMeeting.noActiveMeeting') }}
            </div>
          </q-card-section>
          <q-card-actions align="right" class="bg-primary-dark">
            <q-btn
              v-if="currentMeetingTicket"
              flat
              color="white"
              icon="check_circle"
              :label="$t('generalMeeting.completeMeeting')"
              @click="setTicketStatus(currentMeetingTicket.id, 'concluso')"
            />
          </q-card-actions>
        </q-card>
      </div>

      <div class="col-12 col-md-4">
        <q-card class="next-ticket-card shadow-2 rounded-borders">
          <q-card-section>
            <div class="text-caption text-uppercase text-grey-7">{{ $t('generalMeeting.nextInLine') }}</div>
            <div v-if="nextTicket" class="q-mt-sm">
              <div class="text-h3 text-weight-bolder text-secondary">#{{ nextTicket.ticket_number }}</div>
              <div class="text-h6 q-mt-xs">{{ nextTicket.parent_name }}</div>
              <div class="text-caption text-grey-7">{{ $t('generalMeeting.studentLabel') }} {{ nextTicket.student_name }}</div>
              <div class="text-caption text-primary q-mt-xs">{{ $t('generalMeeting.estimatedTime') }} {{ nextTicket.scheduled_time }}</div>
            </div>
            <div v-else class="q-py-md text-center text-grey-6 text-italic">
              {{ $t('generalMeeting.queueEmpty') }}
            </div>
          </q-card-section>
          <q-card-actions align="right">
            <q-btn
              v-if="nextTicket"
              color="secondary"
              icon="campaign"
              :label="$t('generalMeeting.callNext')"
              unelevated
              no-caps
              @click="setTicketStatus(nextTicket.id, 'in_colloquio')"
            />
          </q-card-actions>
        </q-card>
      </div>

      <div class="col-12 col-md-4">
        <q-card class="shadow-1 rounded-borders">
          <q-card-section>
            <div class="text-caption text-uppercase text-grey-7">{{ $t('generalMeeting.queueSummary') }}</div>
            <div class="row q-mt-sm text-center">
              <div class="col-4">
                <div class="text-h5 text-weight-bold text-primary">{{ tickets.length }}</div>
                <div class="text-caption text-grey-7">{{ $t('generalMeeting.summary.total') }}</div>
              </div>
              <div class="col-4">
                <div class="text-h5 text-weight-bold text-positive">{{ completedTicketsCount }}</div>
                <div class="text-caption text-grey-7">{{ $t('generalMeeting.summary.completed') }}</div>
              </div>
              <div class="col-4">
                <div class="text-h5 text-weight-bold text-warning">{{ waitingTicketsCount }}</div>
                <div class="text-caption text-grey-7">{{ $t('generalMeeting.summary.waiting') }}</div>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Skeleton Loader when loading tickets initially -->
    <div v-if="loadingTickets && tickets.length === 0" class="q-gutter-sm q-mb-md">
      <q-skeleton type="rect" height="52px" class="rounded-borders" />
      <q-skeleton type="rect" height="52px" class="rounded-borders" />
      <q-skeleton type="rect" height="52px" class="rounded-borders" />
    </div>

    <!-- Coda Biglietti Ricevimento -->
    <q-card class="shadow-2 rounded-borders">
      <q-table
        :rows="tickets"
        :columns="ticketColumns"
        row-key="id"
        :loading="loadingTickets"
        :no-data-label="$t('generalMeeting.noQueue')"
        flat
      >
        <template #body-cell-ticket_number="props">
          <q-td :props="props" align="center">
            <q-badge color="primary" class="text-subtitle2 text-weight-bold" :label="`#${props.value}`" />
          </q-td>
        </template>

        <template #body-cell-status="props">
          <q-td :props="props" align="center">
            <q-badge
              :color="getTicketBadgeColor(props.value)"
              :label="$t(`generalMeeting.ticketStatus.${props.value}`)"
              rounded
            />
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" align="right">
            <div class="flex justify-end gap-xs">
              <q-btn
                v-if="props.row.status === 'prenotato'"
                dense
                round
                flat
                icon="notifications_active"
                color="primary"
                @click="setTicketStatus(props.row.id, 'chiamato')"
              >
                <q-tooltip>{{ $t('generalMeeting.callParent') }}</q-tooltip>
              </q-btn>
              <q-btn
                v-if="props.row.status === 'chiamato' || props.row.status === 'prenotato'"
                dense
                round
                flat
                icon="play_arrow"
                color="secondary"
                @click="setTicketStatus(props.row.id, 'in_colloquio')"
              >
                <q-tooltip>{{ $t('generalMeeting.startMeeting') }}</q-tooltip>
              </q-btn>
              <q-btn
                v-if="props.row.status === 'in_colloquio'"
                dense
                round
                flat
                icon="done"
                color="positive"
                @click="setTicketStatus(props.row.id, 'concluso')"
              >
                <q-tooltip>{{ $t('generalMeeting.completeMeeting') }}</q-tooltip>
              </q-btn>
              <q-btn
                v-if="props.row.status === 'prenotato' || props.row.status === 'chiamato'"
                dense
                round
                flat
                icon="person_off"
                color="negative"
                @click="setTicketStatus(props.row.id, 'assente')"
              >
                <q-tooltip>{{ $t('generalMeeting.markAbsent') }}</q-tooltip>
              </q-btn>
            </div>
          </q-td>
        </template>
      </q-table>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import colloquiService from 'src/services/colloquiService'

const { t } = useI18n()
const $q = useQuasar()

const selectedMeetingId = ref('')
const meetingOptions = ref([])
const tickets = ref([])
const loadingTickets = ref(false)

const currentMeetingTicket = computed(() => {
  return tickets.value.find(t => t.status === 'in_colloquio')
})

const nextTicket = computed(() => {
  return tickets.value.find(t => t.status === 'prenotato' || t.status === 'chiamato')
})

const completedTicketsCount = computed(() => {
  return tickets.value.filter(t => t.status === 'concluso').length
})

const waitingTicketsCount = computed(() => {
  return tickets.value.filter(t => t.status === 'prenotato' || t.status === 'chiamato').length
})

const ticketColumns = computed(() => [
  { name: 'ticket_number', label: t('generalMeeting.columns.ticket'), field: 'ticket_number', align: 'center', sortable: true },
  { name: 'scheduled_time', label: t('generalMeeting.columns.scheduledTime'), field: 'scheduled_time', align: 'center' },
  { name: 'parent_name', label: t('generalMeeting.columns.parent'), field: 'parent_name', align: 'left', sortable: true },
  { name: 'student_name', label: t('generalMeeting.columns.student'), field: 'student_name', align: 'left' },
  { name: 'status', label: t('generalMeeting.columns.status'), field: 'status', align: 'center' },
  { name: 'notes', label: t('generalMeeting.columns.notes'), field: 'notes', align: 'left' },
  { name: 'actions', label: t('generalMeeting.columns.actions'), field: 'actions', align: 'right' }
])

function getTicketBadgeColor(status) {
  switch (status) {
    case 'in_colloquio': return 'positive'
    case 'chiamato': return 'warning'
    case 'concluso': return 'grey-6'
    case 'assente': return 'negative'
    default: return 'primary'
  }
}

function extractList(response) {
  if (!response) return []
  const data = response.data !== undefined ? response.data : response
  if (Array.isArray(data)) return data
  if (data && Array.isArray(data.meetings)) return data.meetings
  if (data && Array.isArray(data.tickets)) return data.tickets
  return []
}

async function loadMeetings() {
  try {
    const res = await colloquiService.listGeneralMeetings()
    const list = extractList(res)
    meetingOptions.value = list.map(m => ({ label: `${m.title} (${m.event_date})`, value: m.id }))
    if (meetingOptions.value.length > 0) {
      selectedMeetingId.value = meetingOptions.value[0].value
      loadQueueTickets()
    }
  } catch (err) {
    console.error('Error loading general meetings', err)
  }
}

async function loadQueueTickets() {
  if (!selectedMeetingId.value) return
  loadingTickets.value = true
  try {
    const res = await colloquiService.listQueueTickets(selectedMeetingId.value)
    tickets.value = extractList(res)
  } catch (err) {
    console.error('Error loading queue tickets', err)
  } finally {
    loadingTickets.value = false
  }
}

async function setTicketStatus(ticketId, status) {
  try {
    await colloquiService.updateTicketStatus(ticketId, { status })
    $q.notify({ type: 'positive', message: t('generalMeeting.statusUpdated') })
    loadQueueTickets()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('generalMeeting.statusError') })
  }
}

onMounted(() => {
  loadMeetings()
})
</script>

<style scoped>
.meeting-queue-page {
  max-width: 1300px;
  margin: 0 auto;
}
.current-ticket-card {
  border-left: 6px solid #4caf50;
}
.next-ticket-card {
  border-left: 6px solid var(--q-secondary);
}
</style>
