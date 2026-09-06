<template>
  <q-dialog :model-value="modelValue" @update:model-value="val => emit('update:modelValue', val)">
    <q-card style="width: min(600px, 95vw)" class="rounded-xl overflow-hidden shadow-24 border-slate-300">
      <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
        <div class="text-h6 text-weight-bold">
          <q-icon :name="isEditMode ? (isCurrentAuthor ? 'edit_calendar' : 'info') : 'event_available'" class="q-mr-xs" />
          {{ !isEditMode ? t('agendaPage.newEvent') : (isCurrentAuthor ? t('agendaPage.saveChanges') : t('agendaPage.event')) }}
        </div>
        <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
      </q-card-section>

      <!-- Read-Only Notice for Non-Authors -->
      <div v-if="isEditMode && !isCurrentAuthor" class="bg-amber-50 border-b border-amber-200 q-pa-sm text-amber-900 text-caption row items-center gap-2">
        <q-icon name="info" size="18px" color="amber-9" />
        <span>Creato da <strong>{{ currentEventTeacherName }}</strong>. Solo l'autore dell'evento può apportare modifiche o eliminarlo.</span>
      </div>

      <q-form ref="formRef" @submit="saveEvent" greedy>
        <q-card-section class="q-pa-md">
          <!-- Title Input -->
          <q-input
            v-model="form.title"
            :label="t('agendaPage.eventTitle')"
            outlined
            dense
            tabindex="1"
            class="q-mb-md"
            :readonly="isEditMode && !isCurrentAuthor"
            :rules="[val => (!!val && val.trim().length > 0) || t('common.requiredField') || 'Il titolo è obbligatorio']"
          />

          <!-- Description Input -->
          <q-input
            v-model="form.description"
            :label="t('agendaPage.eventDescription')"
            outlined
            dense
            type="textarea"
            rows="3"
            tabindex="2"
            class="q-mb-md"
            :readonly="isEditMode && !isCurrentAuthor"
          />

          <!-- Row 1: Event Type & Class Select -->
          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-12 col-sm-6">
              <q-select
                v-model="form.type"
                :options="typeOptions"
                :label="t('agendaPage.eventType')"
                outlined dense
                emit-value
                map-options
                tabindex="3"
                :disable="isEditMode && !isCurrentAuthor"
              />
            </div>
            <div class="col-12 col-sm-6">
              <q-select
                v-model="form.class_id"
                :options="availableClassOptions"
                :label="t('agendaPage.targetClass')"
                outlined dense
                emit-value
                map-options
                tabindex="4"
                :disable="isEditMode && !isCurrentAuthor"
                :rules="[val => !!val || t('common.requiredField') || 'Seleziona una classe']"
              />
            </div>
          </div>

          <!-- All-Day Toggle -->
          <div class="row items-center justify-between bg-blue-50/60 q-pa-sm rounded-lg border border-blue-200 q-mb-md">
            <div>
              <div class="text-subtitle2 text-slate-800 font-semibold">📌 {{ t('agendaPage.allDay') }}</div>
              <div class="text-caption text-slate-500">L'evento impegna l'intera giornata senza orario specifico</div>
            </div>
            <q-toggle v-model="form.all_day" color="primary" :disable="isEditMode && !isCurrentAuthor" />
          </div>

          <!-- Row 2: Date & Time Inputs -->
          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-12" :class="form.all_day ? 'col-sm-12' : 'col-sm-4'">
              <q-input
                v-model="form.date"
                type="date"
                :label="t('agendaPage.eventDate')"
                outlined dense
                tabindex="5"
                :readonly="isEditMode && !isCurrentAuthor"
                :rules="[val => !!val || t('common.requiredField') || 'La data è obbligatoria']"
              />
            </div>
            <div v-if="!form.all_day" class="col-6 col-sm-4">
              <q-input
                v-model="form.start_time"
                type="time"
                :label="t('agendaPage.startTime')"
                outlined dense
                tabindex="6"
                :readonly="isEditMode && !isCurrentAuthor"
                :rules="[val => !!val || t('common.requiredField') || 'Ora inizio obbligatoria']"
              />
            </div>
            <div v-if="!form.all_day" class="col-6 col-sm-4">
              <q-input
                v-model="form.end_time"
                type="time"
                :label="t('agendaPage.endTime')"
                outlined dense
                tabindex="7"
                :readonly="isEditMode && !isCurrentAuthor"
                :rules="[
                  val => !!val || t('common.requiredField') || 'Ora fine obbligatoria',
                  val => !form.start_time || val > form.start_time || 'L\'ora di fine deve essere successiva all\'ora di inizio'
                ]"
              />
            </div>
          </div>

          <!-- Toggle: Visible to Students -->
          <div class="row items-center justify-between bg-slate-50 q-pa-md rounded-lg border border-slate-200">
            <div>
              <div class="text-subtitle2 text-slate-700">Visibile agli Studenti & Genitori</div>
              <div class="text-caption text-slate-500">Se disattivato, l'evento sarà visibile solo ai docenti della classe</div>
            </div>
            <q-toggle v-model="form.visible_to_students" color="primary" :disable="isEditMode && !isCurrentAuthor" />
          </div>
        </q-card-section>

        <q-card-actions align="between" class="q-pa-md bg-slate-50 border-t border-slate-100">
          <q-btn
            v-if="isEditMode && isCurrentAuthor"
            color="negative"
            flat
            icon="delete"
            :label="t('agendaPage.deleteEvent') || t('common.delete')"
            no-caps
            @click="confirmDelete"
          />
          <div v-else />

          <div class="row q-gutter-sm">
            <q-btn flat :label="t('common.close') || 'Chiudi'" no-caps v-close-popup />
            <q-btn
              v-if="!isEditMode || isCurrentAuthor"
              type="submit"
              color="primary"
              unelevated
              :label="isEditMode ? t('agendaPage.saveChanges') : t('agendaPage.createEvent')"
              :loading="saving"
              no-caps
            />
          </div>
        </q-card-actions>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useAgendaStore } from '@/stores/agenda'
import { useAuthStore } from '@/stores/auth'
import { useClassesStore } from '@/stores/classes'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  event: {
    type: Object,
    default: null
  },
  classOptions: {
    type: Array,
    default: () => []
  },
  initialDate: {
    type: String,
    default: ''
  },
  initialHour: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'saved', 'deleted'])

const $q = useQuasar()
const { t } = useI18n()
const agendaStore = useAgendaStore()
const authStore = useAuthStore()
const classesStore = useClassesStore()

const formRef = ref(null)
const saving = ref(false)

const form = reactive({
  title: '',
  description: '',
  type: 'compito',
  all_day: false,
  class_id: '',
  date: '',
  start_time: '09:00',
  end_time: '10:00',
  visible_to_students: true
})

const isEditMode = computed(() => !!props.event?.id)
const currentEventTeacherId = computed(() => props.event?.teacher_id || authStore.user?.id)
const currentEventTeacherName = computed(() => props.event?.teacher_name || `${authStore.user?.first_name || ''} ${authStore.user?.last_name || ''}`)

const isCurrentAuthor = computed(() => {
  if (!isEditMode.value) return true
  const currentUserId = authStore.user?.id
  const eventTeacherId = currentEventTeacherId.value
  if (!eventTeacherId || !currentUserId) return false
  return eventTeacherId === currentUserId
})

const availableClassOptions = computed(() => {
  if (props.classOptions && props.classOptions.length > 0) {
    return props.classOptions
  }
  return classesStore.classes.map(c => {
    let nameText = c.name || `Classe ${c.id}`
    if (c.section && !nameText.endsWith(c.section)) nameText += c.section
    if (c.articolazione) nameText += ` - ${c.articolazione}`
    return { label: nameText, value: c.id }
  })
})

const typeOptions = computed(() => [
  { label: t('agendaPage.homework'), value: 'compito' },
  { label: t('agendaPage.test'), value: 'verifica' },
  { label: t('agendaPage.oralTest'), value: 'interrogazione' },
  { label: t('communicationsPage.title') || 'Avviso di Classe', value: 'avviso' },
  { label: t('agendaPage.event'), value: 'evento' },
  { label: t('agendaPage.other'), value: 'altro' }
])

function formatYMD(d) {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function syncForm() {
  if (props.event) {
    const ev = props.event
    form.title = ev.title || ''
    form.description = ev.description || ''
    form.type = ev.type || 'compito'
    form.all_day = !!ev.all_day
    form.class_id = ev.class_id || ''
    form.date = ev.date ? ev.date.substring(0, 10) : formatYMD(new Date())
    form.start_time = ev.start_time || '09:00'
    form.end_time = ev.end_time || '10:00'
    form.visible_to_students = ev.visible_to_students !== false
  } else {
    form.title = ''
    form.description = ''
    form.type = 'compito'
    form.all_day = false
    form.class_id = availableClassOptions.value[0]?.value || ''
    form.date = props.initialDate || formatYMD(new Date())
    form.start_time = props.initialHour || '09:00'
    if (props.initialHour) {
      const h = parseInt(props.initialHour.split(':')[0], 10)
      const nextH = (h + 1 < 24) ? String(h + 1).padStart(2, '0') : '23'
      form.end_time = `${nextH}:00`
    } else {
      form.end_time = '10:00'
    }
    form.visible_to_students = true
  }
}

watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    syncForm()
  }
}, { immediate: true })

watch(() => props.event, () => {
  if (props.modelValue) {
    syncForm()
  }
})

async function saveEvent() {
  if (!isCurrentAuthor.value) return

  saving.value = true
  try {
    const payload = {
      title: form.title,
      description: form.description,
      type: form.type,
      all_day: form.all_day,
      class_id: form.class_id,
      date: form.date,
      start_time: form.all_day ? '00:00' : form.start_time,
      end_time: form.all_day ? '23:59' : form.end_time,
      visible_to_students: form.visible_to_students
    }

    if (isEditMode.value) {
      await agendaStore.updateEvent(props.event.id, payload)
      $q.notify({ type: 'positive', message: t('agendaPage.eventUpdated') })
    } else {
      await agendaStore.createEvent(payload)
      $q.notify({ type: 'positive', message: t('agendaPage.eventCreated') })
    }
    emit('saved')
    emit('update:modelValue', false)
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || t('common.error') || 'Errore durante il salvataggio' })
  } finally {
    saving.value = false
  }
}

async function confirmDelete() {
  if (!isCurrentAuthor.value || !props.event?.id) return

  $q.dialog({
    title: t('agendaPage.confirmDeleteTitle'),
    message: t('agendaPage.confirmDeleteMsg'),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    saving.value = true
    try {
      await agendaStore.deleteEvent(props.event.id)
      $q.notify({ type: 'positive', message: t('agendaPage.eventDeleted') })
      emit('deleted', props.event.id)
      emit('update:modelValue', false)
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || t('common.error') || 'Errore durante l\'eliminazione' })
    } finally {
      saving.value = false
    }
  })
}
</script>
