<template>
  <q-page padding class="bg-slate-50">
    <!-- Header with Month Selector -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="calendar_month" color="primary" class="q-mr-sm" />
          Calendario Scolastico
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Festività, eventi scolastici, verifiche e scadenze importanti
        </p>
      </div>

      <!-- Month Navigation Controls -->
      <div class="row items-center bg-white rounded-xl shadow-sm border border-slate-200 q-pa-xs">
        <q-btn flat round icon="chevron_left" color="primary" @click="prevMonth" />
        <div class="text-subtitle1 text-weight-bold text-slate-800 q-px-md min-w-160 text-center">
          {{ monthName }} {{ currentYear }}
        </div>
        <q-btn flat round icon="chevron_right" color="primary" @click="nextMonth" />
      </div>
    </div>

    <!-- Calendar Grid Card -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden q-mb-lg">
      <div v-if="calendarStore.loading" class="text-center q-pa-xl">
        <q-spinner-dots color="primary" size="40px" />
      </div>

      <div v-else class="q-pa-md">
        <!-- Weekday Headers -->
        <div class="grid grid-cols-7 gap-2 text-center q-mb-sm">
          <div v-for="day in weekDays" :key="day" class="text-weight-bold text-slate-500 text-caption uppercase q-py-xs">
            {{ day }}
          </div>
        </div>

        <!-- Days Grid -->
        <div class="grid grid-cols-7 gap-2">
          <!-- Empty Leading Cells for Offset -->
          <div v-for="blank in leadingBlanks" :key="'blank-' + blank" class="min-h-100 bg-slate-50/50 rounded-xl border border-dashed border-slate-200" />

          <!-- Days of Month -->
          <div
            v-for="dateNum in daysInMonth"
            :key="'date-' + dateNum"
            class="min-h-100 p-2 rounded-xl border transition-all cursor-pointer flex column justify-between"
            :class="getDayCellClass(dateNum)"
            @click="selectDay(dateNum)"
          >
            <!-- Day Number -->
            <div class="row items-center justify-between">
              <span class="text-weight-bold text-body2" :class="isToday(dateNum) ? 'text-white bg-primary q-px-xs rounded' : 'text-slate-700'">
                {{ dateNum }}
              </span>
              <q-badge v-if="getEventsForDay(dateNum).length > 0" color="grey-3" text-color="slate-8" size="xs">
                {{ getEventsForDay(dateNum).length }}
              </q-badge>
            </div>

            <!-- Event Indicators inside Cell -->
            <div class="q-mt-xs space-y-1">
              <div
                v-for="ev in getEventsForDay(dateNum).slice(0, 2)"
                :key="ev.id || ev.title"
                class="text-caption rounded q-px-xs q-py-none text-truncate text-weight-medium text-white"
                :class="getEventBadgeClass(ev.type)"
              >
                {{ ev.title }}
              </div>
              <div v-if="getEventsForDay(dateNum).length > 2" class="text-caption text-slate-400 text-weight-bold">
                +{{ getEventsForDay(dateNum).length - 2 }} altri
              </div>
            </div>
          </div>
        </div>
      </div>
    </q-card>

    <!-- Legend Footer Card -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft q-pa-md">
      <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">Legenda Tipi Evento</div>
      <div class="row q-col-gutter-md">
        <div class="col-6 col-sm-3 row items-center">
          <div class="w-4 h-4 rounded bg-red-500 q-mr-sm" />
          <span class="text-body2 text-slate-600">Festività / Sospensione</span>
        </div>
        <div class="col-6 col-sm-3 row items-center">
          <div class="w-4 h-4 rounded bg-blue-500 q-mr-sm" />
          <span class="text-body2 text-slate-600">Evento Scolastico</span>
        </div>
        <div class="col-6 col-sm-3 row items-center">
          <div class="w-4 h-4 rounded bg-amber-500 q-mr-sm" />
          <span class="text-body2 text-slate-600">Verifica Programmata</span>
        </div>
        <div class="col-6 col-sm-3 row items-center">
          <div class="w-4 h-4 rounded bg-emerald-500 q-mr-sm" />
          <span class="text-body2 text-slate-600">Scadenza / Consegna</span>
        </div>
      </div>
    </q-card>

    <!-- Day Details Dialog -->
    <q-dialog v-model="detailsDialog">
      <q-card style="min-width: 380px; max-width: 500px" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div class="text-h6 text-weight-bold">
            Eventi del {{ selectedDateFormatted }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md">
          <div v-if="selectedDayEvents.length === 0" class="text-center q-pa-md text-slate-400">
            Nessun evento registrato per questo giorno.
          </div>

          <q-list v-else separator>
            <q-item v-for="ev in selectedDayEvents" :key="ev.id || ev.title" class="q-py-md">
              <q-item-section avatar>
                <q-avatar :class="getEventBadgeClass(ev.type)" text-color="white" icon="event" size="40px" />
              </q-item-section>

              <q-item-section>
                <q-item-label class="text-weight-bold text-slate-800">{{ ev.title }}</q-item-label>
                <q-item-label caption class="text-slate-600" v-if="ev.description">
                  {{ ev.description }}
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-chip size="xs" :class="getEventBadgeClass(ev.type)" text-color="white" class="text-weight-bold uppercase">
                  {{ getTypeName(ev.type) }}
                </q-chip>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Chiudi" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { date as qdate } from 'quasar'
import { useSchoolCalendarStore } from '@/stores/schoolCalendar'

const { t } = useI18n()
const calendarStore = useSchoolCalendarStore()
const currentDate = ref(new Date())
const detailsDialog = ref(false)
const selectedDayNum = ref(null)

const currentMonth = computed(() => currentDate.value.getMonth())
const currentYear = computed(() => currentDate.value.getFullYear())

const monthNames = [
  'Gennaio', 'Febbraio', 'Marzo', 'Aprile', 'Maggio', 'Giugno',
  'Luglio', 'Agosto', 'Settembre', 'Ottobre', 'Novembre', 'Dicembre'
]
const monthName = computed(() => monthNames[currentMonth.value].toUpperCase())

const weekDays = ['Lun', 'Mar', 'Mer', 'Gio', 'Ven', 'Sab', 'Dom']

const daysInMonth = computed(() => {
  return new Date(currentYear.value, currentMonth.value + 1, 0).getDate()
})

const leadingBlanks = computed(() => {
  const firstDayIndex = new Date(currentYear.value, currentMonth.value, 1).getDay()
  // Monday is index 0 in Italian calendar: convert Sunday=0 to 6
  return firstDayIndex === 0 ? 6 : firstDayIndex - 1
})

onMounted(() => {
  loadEvents()
})

async function loadEvents() {
  const m = currentMonth.value + 1
  await calendarStore.fetchCalendar(currentYear.value.toString(), m.toString()).catch(() => {})
}

function prevMonth() {
  currentDate.value = new Date(currentYear.value, currentMonth.value - 1, 1)
  loadEvents()
}

function nextMonth() {
  currentDate.value = new Date(currentYear.value, currentMonth.value + 1, 1)
  loadEvents()
}

function getFormattedDateStr(dayNum) {
  const m = (currentMonth.value + 1).toString().padStart(2, '0')
  const d = dayNum.toString().padStart(2, '0')
  return `${currentYear.value}-${m}-${d}`
}

function getEventsForDay(dayNum) {
  const targetDateStr = getFormattedDateStr(dayNum)
  return calendarStore.events.filter(e => e.date === targetDateStr)
}

function isToday(dayNum) {
  const today = new Date()
  return (
    today.getDate() === dayNum &&
    today.getMonth() === currentMonth.value &&
    today.getFullYear() === currentYear.value
  )
}

function getDayCellClass(dayNum) {
  const events = getEventsForDay(dayNum)
  if (events.length === 0) return 'bg-white border-slate-200 hover:border-primary'

  const hasHoliday = events.some(e => e.type === 'holiday')
  if (hasHoliday) return 'bg-red-50/70 border-red-300'

  const hasExam = events.some(e => e.type === 'exam')
  if (hasExam) return 'bg-amber-50/70 border-amber-300'

  const hasDeadline = events.some(e => e.type === 'deadline')
  if (hasDeadline) return 'bg-emerald-50/70 border-emerald-300'

  return 'bg-blue-50/70 border-blue-300'
}

function getEventBadgeClass(type) {
  switch (type) {
    case 'holiday': return 'bg-red-500'
    case 'exam': return 'bg-amber-500'
    case 'deadline': return 'bg-emerald-500'
    case 'event': default: return 'bg-blue-500'
  }
}

function getTypeName(type) {
  switch (type) {
    case 'holiday': return 'Festività'
    case 'exam': return 'Verifica'
    case 'deadline': return 'Scadenza'
    case 'event': default: return 'Evento'
  }
}

function selectDay(dayNum) {
  selectedDayNum.value = dayNum
  detailsDialog.value = true
}

const selectedDayEvents = computed(() => {
  if (!selectedDayNum.value) return []
  return getEventsForDay(selectedDayNum.value)
})

const selectedDateFormatted = computed(() => {
  if (!selectedDayNum.value) return ''
  const d = new Date(currentYear.value, currentMonth.value, selectedDayNum.value)
  return qdate.formatDate(d, 'DD/MM/YYYY')
})
</script>

<style scoped>
.min-h-100 {
  min-height: 100px;
}
.min-w-160 {
  min-width: 160px;
}
.grid-cols-7 {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
}
.text-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
