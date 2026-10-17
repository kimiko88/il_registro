<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">
          <q-icon name="support_agent" class="q-mr-sm" />{{ t('helpDesk.bookingTitle') }}
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          {{ t('helpDesk.bookingSubtitle') }}
        </p>
      </div>
      <q-btn flat icon="refresh" :label="t('common.refresh')" @click="loadSlots" :loading="loading" />
    </div>

    <!-- Slots List -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots size="40px" color="primary" />
    </div>

    <div v-else-if="slots.length === 0" class="q-pa-xl text-center text-grey-6">
      <q-icon name="event_busy" size="64px" class="q-mb-md text-grey-4" />
      <div class="text-subtitle1">{{ t('common.noData') }}</div>
      <div class="text-caption">{{ t('helpDesk.bookingSubtitle') }}</div>
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="slot in slots" :key="slot.id" class="col-12 col-md-6">
        <q-card flat bordered class="bg-white">
          <q-card-section class="row items-center justify-between bg-blue-grey-1 q-py-sm">
            <div class="text-subtitle2 text-weight-bold text-primary">
              <q-icon name="menu_book" class="q-mr-xs" /> {{ slot.subject_name || t('helpDesk.subject') }}
            </div>
            <q-badge
              :color="slot.status === 'open' ? 'positive' : 'negative'"
              :label="slot.status === 'open' ? t('common.active') : t('common.inactive')"
            />
          </q-card-section>

          <q-card-section>
            <div class="text-subtitle1 text-weight-bold">{{ slot.teacher_name }}</div>
            <div class="text-caption text-grey-7 q-mt-xs">
              <q-icon name="event" class="q-mr-xs" /> {{ t('helpDesk.dateTime') }}: <strong>{{ formatDate(slot.slot_date) }} ({{ slot.start_time }} - {{ slot.end_time }})</strong>
            </div>
            <div class="text-caption text-grey-7 q-mt-xs">
              <q-icon name="group" class="q-mr-xs" /> {{ t('helpDesk.remainingSeats') }}: <strong>{{ slot.bookings_count || 0 }} / {{ slot.max_capacity }}</strong>
            </div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="right">
            <q-btn
              color="primary"
              :label="t('helpDesk.bookSlot')"
              unelevated
              no-caps
              :disable="slot.status !== 'open'"
              @click="openBookingDialog(slot)"
            />
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Booking Dialog -->
    <q-dialog v-model="showBookingDialog" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ t('helpDesk.bookSlot') }}</div>
        </q-card-section>
        <q-card-section class="q-pt-md">
          <div class="text-subtitle2 text-weight-bold q-mb-xs">{{ t('helpDesk.teacher') }}: {{ selectedSlot?.teacher_name }}</div>
          <div class="text-caption text-grey-7 q-mb-md">{{ t('helpDesk.subject') }}: {{ selectedSlot?.subject_name }} | {{ t('helpDesk.dateTime') }}: {{ selectedSlot?.slot_date }} ({{ selectedSlot?.start_time }})</div>

          <q-input
            v-model="bookingForm.topic_description"
            :label="t('helpDesk.topic') + ' *'"
            type="textarea"
            rows="3"
            outlined
            dense
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('helpDesk.bookSlot')" unelevated :loading="submitting" @click="confirmBooking" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { helpDeskService } from '@/services/helpDeskService'

export default defineComponent({
  name: 'HelpDeskBooking',
  setup() {
    const { t } = useI18n()
    const slots = ref([])
    const loading = ref(false)
    const submitting = ref(false)
    const showBookingDialog = ref(false)
    const selectedSlot = ref(null)

    const bookingForm = ref({
      topic_description: ''
    })

    const loadSlots = async () => {
      loading.value = true
      try {
        const res = await helpDeskService.getSlots()
        slots.value = res.data || []
      } catch (e) {
        console.error('Errore caricamento slot help desk:', e)
      } finally {
        loading.value = false
      }
    }

    const openBookingDialog = (slot) => {
      selectedSlot.value = slot
      bookingForm.value.topic_description = ''
      showBookingDialog.value = true
    }

    const confirmBooking = async () => {
      if (!selectedSlot.value) return
      submitting.value = true
      try {
        await helpDeskService.bookSlot(selectedSlot.value.id, {
          topic_description: bookingForm.value.topic_description
        })
        showBookingDialog.value = false
        await loadSlots()
      } catch (e) {
        console.error('Errore prenotazione slot:', e)
      } finally {
        submitting.value = false
      }
    }

    const formatDate = (d) => {
      if (!d) return ''
      return new Date(d).toLocaleDateString(undefined, { weekday: 'short', day: '2-digit', month: '2-digit', year: 'numeric' })
    }

    onMounted(loadSlots)

    return {
      t,
      slots,
      loading,
      submitting,
      showBookingDialog,
      selectedSlot,
      bookingForm,
      loadSlots,
      openBookingDialog,
      confirmBooking,
      formatDate
    }
  }
})
</script>
