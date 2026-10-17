<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">
          <q-icon name="co_present" class="q-mr-sm" />{{ t('helpDesk.mgmtTitle') }}
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          {{ t('helpDesk.mgmtSubtitle') }}
        </p>
      </div>
      <q-btn color="primary" icon="add" :label="t('helpDesk.createSlot')" unelevated no-caps @click="showNewSlotDialog = true" />
    </div>

    <!-- Slots Management -->
    <q-card flat bordered class="bg-white">
      <q-table
        :rows="slots"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        :no-data-label="t('common.noData')"
      >
        <template #body-cell-status="props">
          <q-td :props="props">
            <q-chip :color="getStatusColor(props.row.status)" text-color="white" size="sm">
              {{ formatStatus(props.row.status) }}
            </q-chip>
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" align="right">
            <div class="row q-gutter-xs justify-end">
              <q-btn
                size="sm"
                color="primary"
                flat
                icon="people"
                :label="t('helpDesk.rollCall')"
                no-caps
                @click="openBookingsModal(props.row)"
              />
              <q-btn
                v-if="props.row.status !== 'completed'"
                size="sm"
                color="positive"
                unelevated
                icon="check"
                :label="t('helpDesk.completeSlot')"
                no-caps
                @click="completeSlot(props.row)"
              />
            </div>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Bookings and Attendance Dialog -->
    <q-dialog v-model="showBookingsDialog" persistent>
      <q-card style="min-width: 550px">
        <q-card-section class="bg-primary text-white row items-center justify-between">
          <div class="text-h6">{{ t('helpDesk.rollCall') }}</div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>
        <q-card-section class="q-pa-none">
          <div v-if="bookings.length === 0" class="q-pa-md text-caption text-grey-6 text-center">
            {{ t('common.noData') }}
          </div>
          <q-list v-else separator>
            <q-item v-for="b in bookings" :key="b.id">
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ b.student_name || t('religionAlternative.student') }}</q-item-label>
                <q-item-label caption class="text-grey-8">{{ t('helpDesk.topic') }}: {{ b.topic_description }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <div class="row q-gutter-xs">
                  <q-btn
                    size="sm"
                    :color="b.status === 'attended' ? 'positive' : 'grey-4'"
                    :text-color="b.status === 'attended' ? 'white' : 'dark'"
                    :label="t('helpDesk.markAttended')"
                    unelevated
                    no-caps
                    @click="setAttendance(b, 'attended')"
                  />
                  <q-btn
                    size="sm"
                    :color="b.status === 'absent' ? 'negative' : 'grey-4'"
                    :text-color="b.status === 'absent' ? 'white' : 'dark'"
                    :label="t('helpDesk.markAbsent')"
                    unelevated
                    no-caps
                    @click="setAttendance(b, 'absent')"
                  />
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- New Slot Dialog -->
    <q-dialog v-model="showNewSlotDialog" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ t('helpDesk.createSlot') }}</div>
        </q-card-section>
        <q-card-section class="q-pt-md">
          <q-input v-model="newSlotForm.subject_id" :label="t('helpDesk.subject') + ' *'" outlined dense class="q-mb-sm" />
          <q-input v-model="newSlotForm.slot_date" type="date" :label="t('pctoCompanyTutor.date') + ' *'" outlined dense class="q-mb-sm" />
          <div class="row q-col-gutter-sm q-mb-sm">
            <div class="col-6">
              <q-input v-model="newSlotForm.start_time" :label="t('helpDesk.dateTime') + ' (Start) *'" outlined dense placeholder="15:00" />
            </div>
            <div class="col-6">
              <q-input v-model="newSlotForm.end_time" :label="t('helpDesk.dateTime') + ' (End) *'" outlined dense placeholder="16:30" />
            </div>
          </div>
          <q-input v-model.number="newSlotForm.max_capacity" type="number" :label="t('helpDesk.maxCapacity') + ' *'" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('common.save')" unelevated :loading="submitting" @click="createSlot" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { helpDeskService } from '@/services/helpDeskService'

export default defineComponent({
  name: 'HelpDeskManagement',
  setup() {
    const { t } = useI18n()
    const slots = ref([])
    const bookings = ref([])
    const loading = ref(false)
    const submitting = ref(false)
    const showNewSlotDialog = ref(false)
    const showBookingsDialog = ref(false)
    const activeSlot = ref(null)

    const newSlotForm = ref({
      subject_id: 'default-math',
      slot_date: new Date().toISOString().substring(0, 10),
      start_time: '15:00',
      end_time: '16:30',
      max_capacity: 4
    })

    const columns = computed(() => [
      { name: 'slot_date', label: t('pctoCompanyTutor.date'), field: 'slot_date', align: 'left' },
      { name: 'hours', label: t('helpDesk.dateTime'), align: 'left', field: row => `${row.start_time} - ${row.end_time}` },
      { name: 'subject', label: t('helpDesk.subject'), field: 'subject_name', align: 'left' },
      { name: 'bookings', label: t('helpDesk.remainingSeats'), align: 'center', field: row => `${row.bookings_count || 0} / ${row.max_capacity}` },
      { name: 'status', label: t('common.details'), field: 'status', align: 'center' },
      { name: 'actions', label: t('common.actions'), align: 'right' }
    ])

    const loadSlots = async () => {
      loading.value = true
      try {
        const res = await helpDeskService.getSlots()
        slots.value = res.data || []
      } catch (e) {
        console.error('Errore caricamento slot docente:', e)
      } finally {
        loading.value = false
      }
    }

    const openBookingsModal = async (slot) => {
      activeSlot.value = slot
      try {
        const res = await helpDeskService.getBookings(slot.id)
        bookings.value = res.data || []
        showBookingsDialog.value = true
      } catch (e) {
        console.error('Errore prenotazioni:', e)
      }
    }

    const setAttendance = async (booking, status) => {
      try {
        await helpDeskService.markAttendance(booking.id, status)
        booking.status = status
      } catch (e) {
        console.error('Errore presenza:', e)
      }
    }

    const completeSlot = async (slot) => {
      try {
        await helpDeskService.completeSlot(slot.id)
        slot.status = 'completed'
      } catch (e) {
        console.error('Errore chiusura slot:', e)
      }
    }

    const createSlot = async () => {
      submitting.value = true
      try {
        await helpDeskService.createSlot(newSlotForm.value)
        showNewSlotDialog.value = false
        await loadSlots()
      } catch (e) {
        console.error('Errore creazione slot:', e)
      } finally {
        submitting.value = false
      }
    }

    const formatStatus = (s) => {
      switch (s) {
        case 'open': return t('common.active')
        case 'fully_booked': return t('helpDesk.remainingSeats')
        case 'completed': return t('helpDesk.slotCompletedSuccess')
        case 'cancelled': return t('common.cancel')
        default: return s
      }
    }

    const getStatusColor = (s) => {
      switch (s) {
        case 'open': return 'positive'
        case 'fully_booked': return 'warning'
        case 'completed': return 'teal-8'
        default: return 'negative'
      }
    }

    onMounted(loadSlots)

    return {
      t,
      slots,
      bookings,
      loading,
      submitting,
      showNewSlotDialog,
      showBookingsDialog,
      newSlotForm,
      columns,
      loadSlots,
      openBookingsModal,
      setAttendance,
      completeSlot,
      createSlot,
      formatStatus,
      getStatusColor
    }
  }
})
</script>
