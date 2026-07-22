<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="event" color="primary" class="q-mr-sm" />
          Agenda & Calendario Didattico
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Pianificazione compiti, verifiche, ed avvisi per le tue classi
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-select
          v-model="selectedClassFilter"
          :options="classOptions"
          label="Filtra per Classe"
          outlined dense
          clearable
          emit-value
          map-options
          style="min-width: 180px"
          class="bg-white"
          @update:model-value="loadAgendaEvents"
        />
        <q-btn
          color="primary"
          unelevated
          icon="add"
          label="Nuovo Evento"
          class="rounded-lg q-px-md"
          no-caps
          @click="openCreateDialog"
        />
        <q-btn flat round icon="refresh" color="primary" :loading="agendaStore.loading" @click="loadAgendaEvents" />
      </div>
    </div>

    <!-- Main Content Layout: Left Calendar & Right Day Timeline -->
    <div class="row q-col-gutter-lg">
      <!-- Left Column: Calendar & Quick Filters -->
      <div class="col-12 col-md-5 col-lg-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden q-mb-md">
          <q-card-section class="q-pa-sm text-center">
            <q-date
              v-model="selectedDate"
              today-btn
              color="primary"
              flat
              class="full-width text-outfit"
              @update:model-value="onDateChange"
            />
          </q-card-section>
        </q-card>

        <!-- Summary Legend Card -->
        <q-card flat bordered class="rounded-xl bg-white shadow-soft q-pa-md">
          <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-sm">Legenda Tipi Evento</div>
          <div class="row q-col-gutter-sm text-caption">
            <div class="col-6 row items-center">
              <span class="inline-block w-3 h-3 rounded-full bg-blue-500 mr-2"></span> Compito
            </div>
            <div class="col-6 row items-center">
              <span class="inline-block w-3 h-3 rounded-full bg-red-500 mr-2"></span> Verifica
            </div>
            <div class="col-6 row items-center">
              <span class="inline-block w-3 h-3 rounded-full bg-amber-500 mr-2"></span> Avviso
            </div>
            <div class="col-6 row items-center">
              <span class="inline-block w-3 h-3 rounded-full bg-emerald-500 mr-2"></span> Evento / Uscita
            </div>
          </div>
        </q-card>
      </div>

      <!-- Right Column: Selected Date Timeline & Event List -->
      <div class="col-12 col-md-7 col-lg-8">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft min-h-500 overflow-hidden">
          <q-card-section class="bg-slate-100 border-b border-slate-200 row items-center justify-between q-py-sm q-px-md">
            <div class="row items-center">
              <q-icon name="today" size="24px" color="primary" class="q-mr-xs" />
              <span class="text-subtitle1 text-weight-bold text-slate-800">
                Eventi del {{ formattedSelectedDate }}
              </span>
            </div>
            <q-badge color="primary" class="q-px-sm q-py-xs text-weight-bold">
              {{ dayEvents.length }} {{ dayEvents.length === 1 ? 'evento' : 'eventi' }}
            </q-badge>
          </q-card-section>

          <!-- Loading State -->
          <div v-if="agendaStore.loading" class="text-center q-pa-xl">
            <q-spinner-dots color="primary" size="40px" />
          </div>

          <!-- Empty State -->
          <div v-else-if="dayEvents.length === 0" class="text-center q-pa-xl text-slate-400">
            <q-icon name="event_busy" size="64px" class="q-mb-md opacity-40" />
            <div class="text-h6">Nessun evento in agenda</div>
            <div class="text-caption">Non sono presenti verifiche, compiti o avvisi per la data selezionata.</div>
            <q-btn
              color="primary"
              flat
              icon="add"
              label="Aggiungi Evento"
              class="q-mt-md"
              no-caps
              @click="openCreateDialog"
            />
          </div>

          <!-- Timeline Events View -->
          <q-card-section v-else class="q-pa-md">
            <q-timeline color="primary" class="q-px-sm">
              <q-timeline-entry
                v-for="ev in dayEvents"
                :key="ev.id"
                :title="ev.title"
                :subtitle="getEventSubtitle(ev)"
                :icon="getEventIcon(ev.type)"
                :color="getEventTimelineColor(ev.type)"
                class="cursor-pointer hover:bg-slate-50 transition-colors rounded-lg q-pa-xs"
                @click="openEditDialog(ev)"
              >
                <template #title>
                  <div class="row items-center justify-between">
                    <span class="text-weight-bold text-slate-800 text-subtitle1">{{ ev.title }}</span>
                    <div class="row items-center q-gutter-xs">
                      <q-chip size="xs" :color="getTypeColor(ev.type)" text-color="white" class="text-weight-bold uppercase">
                        {{ ev.type }}
                      </q-chip>
                      <q-chip size="xs" color="indigo-1" text-color="indigo-8" class="text-weight-bold" v-if="ev.class_name || ev.class_id">
                        Classe {{ ev.class_name || ev.class_id }}
                      </q-chip>
                    </div>
                  </div>
                </template>

                <div class="text-body2 text-slate-600 q-mt-xs" v-if="ev.description">
                  {{ ev.description }}
                </div>

                <div class="row items-center justify-between q-mt-sm text-caption text-slate-400 border-t border-slate-100 pt-2">
                  <div>
                    <q-icon name="schedule" size="14px" class="q-mr-xs" />
                    {{ ev.start_time || '08:00' }} <span v-if="ev.end_time">- {{ ev.end_time }}</span>
                    <span v-if="ev.subject_name" class="q-ml-sm text-weight-medium text-slate-600">· {{ ev.subject_name }}</span>
                  </div>
                  <div class="row items-center text-slate-500">
                    <q-icon :name="ev.visible_to_students !== false ? 'visibility' : 'visibility_off'" size="14px" class="q-mr-xs" />
                    <span>{{ ev.visible_to_students !== false ? 'Visibile Studenti' : 'Riservato' }}</span>
                  </div>
                </div>
              </q-timeline-entry>
            </q-timeline>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Floating Action Button for Mobile / Quick Create -->
    <q-page-sticky position="bottom-right" :offset="[18, 18]">
      <q-btn round color="primary" icon="add" size="lg" class="shadow-lg" @click="openCreateDialog" />
    </q-page-sticky>

    <!-- Create / Edit Event Dialog -->
    <q-dialog v-model="dialogVisible">
      <q-card style="min-width: 480px; max-width: 600px" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div class="text-h6 text-weight-bold">
            <q-icon :name="isEditMode ? 'edit_calendar' : 'event_available'" class="q-mr-xs" />
            {{ isEditMode ? 'Modifica Evento Agenda' : 'Nuovo Evento Agenda' }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md space-y-4">
          <!-- Title Input -->
          <q-input
            v-model="form.title"
            label="Titolo Evento *"
            outlined
            dense
            :rules="[val => !!val || 'Il titolo è obbligatorio']"
          />

          <!-- Description Input -->
          <q-input
            v-model="form.description"
            label="Descrizione / Dettagli"
            outlined
            dense
            type="textarea"
            rows="3"
          />

          <!-- Row 1: Event Type & Class Select -->
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-select
                v-model="form.type"
                :options="typeOptions"
                label="Tipo Evento *"
                outlined dense
                emit-value
                map-options
              />
            </div>
            <div class="col-12 col-sm-6">
              <q-select
                v-model="form.class_id"
                :options="classOptions"
                label="Classe Destinataria *"
                outlined dense
                emit-value
                map-options
                :rules="[val => !!val || 'Seleziona una classe']"
              />
            </div>
          </div>

          <!-- Row 2: Date & Times -->
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-4">
              <q-input
                v-model="form.date"
                type="date"
                label="Data *"
                outlined dense
                :rules="[val => !!val || 'Data obbligatoria']"
              />
            </div>
            <div class="col-12 col-sm-4">
              <q-input
                v-model="form.start_time"
                type="time"
                label="Ora Inizio"
                outlined dense
              />
            </div>
            <div class="col-12 col-sm-4">
              <q-input
                v-model="form.end_time"
                type="time"
                label="Ora Fine"
                outlined dense
              />
            </div>
          </div>

          <!-- Visibility Toggle -->
          <div class="bg-slate-50 q-pa-sm rounded-lg border border-slate-200">
            <q-toggle
              v-model="form.visible_to_students"
              label="Visibile agli studenti e genitori"
              color="primary"
            />
          </div>
        </q-card-section>

        <q-separator />

        <!-- Card Actions -->
        <q-card-actions align="right" class="q-pa-md">
          <q-btn
            v-if="isEditMode"
            flat
            color="negative"
            icon="delete"
            label="Elimina"
            :loading="saving"
            @click="confirmDelete"
          />
          <div class="flex-1"></div>
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn
            color="primary"
            :label="isEditMode ? 'Salva Modifiche' : 'Crea Evento'"
            :loading="saving"
            @click="saveEvent"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { useAgendaStore } from '@/stores/agenda'
import { useClassesStore } from '@/stores/classes'

const $q = useQuasar()
const agendaStore = useAgendaStore()
const classesStore = useClassesStore()

const todayStr = qdate.formatDate(new Date(), 'YYYY/MM/DD')
const selectedDate = ref(todayStr)
const selectedClassFilter = ref(null)

const dialogVisible = ref(false)
const isEditMode = ref(false)
const editId = ref(null)
const saving = ref(false)

const form = reactive({
  title: '',
  description: '',
  type: 'compito',
  class_id: '',
  date: new Date().toISOString().substring(0, 10),
  start_time: '09:00',
  end_time: '10:00',
  visible_to_students: true
})

const typeOptions = [
  { label: 'Compito per casa', value: 'compito' },
  { label: 'Verifica / Incontro', value: 'verifica' },
  { label: 'Avviso di Classe', value: 'avviso' },
  { label: 'Evento / Uscita Didattica', value: 'evento' }
]

const classOptions = computed(() => {
  return classesStore.classes.map(c => ({
    label: c.name || `Classe ${c.id}`,
    value: c.id
  }))
})

const formattedSelectedDate = computed(() => {
  if (!selectedDate.value) return ''
  const normalized = selectedDate.value.replace(/\//g, '-')
  return qdate.formatDate(new Date(normalized), 'DD/MM/YYYY')
})

const isoSelectedDate = computed(() => {
  if (!selectedDate.value) return new Date().toISOString().substring(0, 10)
  return selectedDate.value.replace(/\//g, '-')
})

const dayEvents = computed(() => {
  const sel = isoSelectedDate.value
  return agendaStore.events.filter(ev => {
    const evDate = ev.date ? ev.date.substring(0, 10) : ''
    const matchesDate = evDate === sel
    const matchesClass = !selectedClassFilter.value || ev.class_id === selectedClassFilter.value
    return matchesDate && matchesClass
  })
})

onMounted(async () => {
  await classesStore.fetchAssignedClasses().catch(() => classesStore.fetchClasses())
  await loadAgendaEvents()
})

async function loadAgendaEvents() {
  const params = {}
  if (selectedClassFilter.value) {
    params.class_id = selectedClassFilter.value
  }
  await agendaStore.fetchAgenda(params).catch(() => {})
}

function onDateChange() {
  loadAgendaEvents()
}

function openCreateDialog() {
  isEditMode.value = false
  editId.value = null
  form.title = ''
  form.description = ''
  form.type = 'compito'
  form.class_id = classOptions.value[0]?.value || ''
  form.date = isoSelectedDate.value
  form.start_time = '09:00'
  form.end_time = '10:00'
  form.visible_to_students = true
  dialogVisible.value = true
}

function openEditDialog(ev) {
  isEditMode.value = true
  editId.value = ev.id
  form.title = ev.title || ''
  form.description = ev.description || ''
  form.type = ev.type || 'compito'
  form.class_id = ev.class_id || ''
  form.date = ev.date ? ev.date.substring(0, 10) : isoSelectedDate.value
  form.start_time = ev.start_time || '09:00'
  form.end_time = ev.end_time || '10:00'
  form.visible_to_students = ev.visible_to_students !== false
  dialogVisible.value = true
}

async function saveEvent() {
  if (!form.title || !form.class_id || !form.date) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }

  saving.value = true
  try {
    const payload = {
      title: form.title,
      description: form.description,
      type: form.type,
      class_id: form.class_id,
      date: form.date,
      start_time: form.start_time,
      end_time: form.end_time,
      visible_to_students: form.visible_to_students
    }

    if (isEditMode.value) {
      await agendaStore.updateEvent(editId.value, payload)
      $q.notify({ type: 'positive', message: 'Evento modificato con successo' })
    } else {
      await agendaStore.createEvent(payload)
      $q.notify({ type: 'positive', message: 'Evento creato con successo' })
    }
    dialogVisible.value = false
    await loadAgendaEvents()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante il salvataggio' })
  } finally {
    saving.value = false
  }
}

async function confirmDelete() {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: 'Sei sicuro di voler eliminare questo evento dall\'agenda?',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    saving.value = true
    try {
      await agendaStore.deleteEvent(editId.value)
      $q.notify({ type: 'positive', message: 'Evento eliminato con successo' })
      dialogVisible.value = false
      await loadAgendaEvents()
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante l\'eliminazione' })
    } finally {
      saving.value = false
    }
  })
}

function getTypeColor(type) {
  switch (type) {
    case 'compito': return 'info'
    case 'verifica': return 'negative'
    case 'avviso': return 'warning'
    case 'evento': return 'positive'
    default: return 'primary'
  }
}

function getEventTimelineColor(type) {
  switch (type) {
    case 'compito': return 'blue'
    case 'verifica': return 'red'
    case 'avviso': return 'amber'
    case 'evento': return 'emerald'
    default: return 'primary'
  }
}

function getEventIcon(type) {
  switch (type) {
    case 'compito': return 'assignment'
    case 'verifica': return 'quiz'
    case 'avviso': return 'campaign'
    case 'evento': return 'directions_bus'
    default: return 'event'
  }
}

function getEventSubtitle(ev) {
  return ev.subject_name ? `Materia: ${ev.subject_name}` : 'Evento Didattico'
}
</script>

<style scoped>
.min-h-500 {
  min-height: 500px;
}
</style>
