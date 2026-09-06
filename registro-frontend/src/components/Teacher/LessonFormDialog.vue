<template>
  <q-dialog :model-value="modelValue" persistent @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(550px, 95vw); max-width: 95vw;" class="rounded-xl shadow-24" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-white text-slate-800'">
      <q-card-section class="row items-center justify-between" :class="$q.dark.isActive ? 'bg-grey-9 text-white' : 'bg-indigo-7 text-white'">
        <div class="text-h6 text-weight-bold row items-center">
          <q-icon name="menu_book" class="q-mr-sm" size="24px" />
          {{ isEditing ? ($t('classRegister.editLesson') || 'Modifica Lezione') : ($t('classRegister.newLesson') || 'Nuova Lezione') }}
        </div>
        <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-gutter-y-md q-pt-md">
        <q-form ref="formRef" @submit.prevent="handleSubmit" class="q-gutter-y-md">
          <q-select
            v-model="form.subject_id"
            :options="availableSubjectOptions"
            option-value="subject_id"
            option-label="subject_name"
            emit-value
            map-options
            :label="($t('classRegister.subject') || 'Materia') + ' *'"
            outlined
            dense
            :rules="[v => !!v || ($t('common.requiredField') || 'Campo obbligatorio')]"
          />

          <q-input
            v-model="form.date"
            type="date"
            :label="($t('classRegister.dateLabel') || 'Data') + ' *'"
            outlined
            dense
            :rules="[v => !!v || ($t('common.requiredField') || 'Campo obbligatorio')]"
          />

          <q-input
            v-model="form.topic"
            :label="($t('classRegister.lessonTopic') || 'Argomento della Lezione') + ' *'"
            outlined
            dense
            placeholder="Es: Le equazioni di secondo grado o Supplenza"
            :rules="[v => (!!v && v.trim().length > 0) || ($t('common.requiredField') || 'Campo obbligatorio')]"
          />

          <div class="row q-col-gutter-sm">
            <div class="col-12 col-md-6">
              <q-input
                v-model.number="form.hour"
                type="number"
                :label="($t('classRegister.hour') || 'Ora Lezione') + ' *'"
                outlined
                dense
                min="1"
                max="10"
                :rules="[v => (!!v && v >= 1 && v <= 10) || ($t('common.requiredField') || 'Campo obbligatorio')]"
              />
            </div>
            <div class="col-12 col-md-6">
              <q-input
                v-model.number="form.duration"
                type="number"
                :label="($t('classRegister.duration') || 'Durata (ore)') + ' *'"
                outlined
                dense
                min="1"
                max="5"
                :rules="[v => (!!v && v >= 1 && v <= 5) || ($t('common.requiredField') || 'Campo obbligatorio')]"
              />
            </div>
          </div>

          <q-select
            v-model="form.type"
            :options="lessonTypeOptions"
            :label="($t('classRegister.lessonType') || 'Tipo di Lezione') + ' *'"
            outlined
            dense
            :rules="[v => !!v || ($t('common.requiredField') || 'Campo obbligatorio')]"
          />

          <q-select
            v-model="form.activity_type"
            :options="activityTypeOptions"
            option-value="value"
            option-label="label"
            emit-value
            map-options
            :label="$t('classRegister.activityType') || 'Tipologia Attività'"
            outlined
            dense
          >
            <template v-slot:option="scope">
              <q-item v-bind="scope.itemProps">
                <q-item-section avatar>
                  <q-icon :name="scope.opt.icon" :color="scope.opt.color" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ scope.opt.label }}</q-item-label>
                  <q-item-label caption>{{ scope.opt.caption }}</q-item-label>
                </q-item-section>
              </q-item>
            </template>
          </q-select>

          <q-banner
            v-if="isPctoOrOrientamento"
            class="rounded-borders q-mt-xs"
            :class="$q.dark.isActive ? 'bg-blue-10 text-blue-1' : 'bg-blue-1 text-blue-9'"
            dense
          >
            <template v-slot:avatar><q-icon name="info" color="blue-7" /></template>
            Per le attività {{ getActivityTypeLabel(form.activity_type) }} non è possibile inserire valutazioni agli studenti.
          </q-banner>

          <q-toggle
            v-model="form.is_co_teaching"
            label="Compresenza (docente co-presente in classe)"
            color="deep-purple"
            icon="people"
          />

          <q-input
            v-model="form.notes"
            label="Note / Osservazioni interne"
            type="textarea"
            outlined
            dense
            autogrow
            placeholder="(Opzionale)"
          />

          <q-toggle
            v-if="!isEditing && !isPctoOrOrientamento"
            v-model="assignHomeworkToo"
            label="Assegna anche un compito per questa lezione"
            color="orange"
          />

          <template v-if="assignHomeworkToo && !isEditing">
            <q-input
              v-model="form.homeworkDesc"
              label="Descrizione Compito *"
              type="textarea"
              outlined
              dense
              autogrow
              placeholder="Es: Pagine 145-150, esercizi 12-18"
              :rules="[v => !assignHomeworkToo || (!!v && v.trim().length > 0) || ($t('common.requiredField') || 'Campo obbligatorio')]"
            />
            <q-input
              v-model="form.homeworkDue"
              type="date"
              label="Data Consegna *"
              outlined
              dense
              :rules="[v => !assignHomeworkToo || !!v || ($t('common.requiredField') || 'Campo obbligatorio')]"
            />
          </template>

          <q-card-actions align="right" class="q-px-none q-pt-sm">
            <q-btn flat :label="$t('common.cancel') || 'Annulla'" v-close-popup color="grey-7" no-caps />
            <q-btn
              type="submit"
              color="primary"
              :label="isEditing ? ($t('common.update') || 'Aggiorna') : ($t('common.save') || 'Salva')"
              :loading="saving"
              no-caps
              class="rounded-lg q-px-md font-bold"
            />
          </q-card-actions>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useQuasar } from 'quasar'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  isEditing: {
    type: Boolean,
    default: false
  },
  initialData: {
    type: Object,
    default: () => ({})
  },
  availableSubjectOptions: {
    type: Array,
    default: () => []
  },
  activityTypeOptions: {
    type: Array,
    default: () => []
  },
  saving: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'save'])
const $q = useQuasar()

const formRef = ref(null)
const assignHomeworkToo = ref(false)

const lessonTypeOptions = [
  'Frontale',
  'Supplenza',
  'Laboratorio',
  'Verifica',
  'Discussione',
  'Lavoro di gruppo',
  'Altro'
]

const form = ref({
  date: '',
  hour: 1,
  duration: 1,
  topic: '',
  type: 'Frontale',
  activity_type: 'standard',
  is_co_teaching: false,
  notes: '',
  subject_id: null,
  homeworkDesc: '',
  homeworkDue: ''
})

watch(() => props.initialData, (val) => {
  if (val) {
    form.value = {
      date: val.date || '',
      hour: val.hour || 1,
      duration: val.duration || 1,
      topic: val.topic || '',
      type: val.type || 'Frontale',
      activity_type: val.activity_type || 'standard',
      is_co_teaching: !!val.is_co_teaching,
      notes: val.notes || '',
      subject_id: val.subject_id || null,
      homeworkDesc: val.homeworkDesc || '',
      homeworkDue: val.homeworkDue || ''
    }
  }
}, { immediate: true, deep: true })

const isPctoOrOrientamento = computed(() =>
  form.value.activity_type === 'pcto' ||
  form.value.activity_type === 'orientamento' ||
  form.value.activity_type === 'pcto_orientamento'
)

function getActivityTypeLabel(type) {
  const opt = props.activityTypeOptions.find(o => o.value === type)
  return opt ? opt.label : type
}

async function handleSubmit() {
  if (formRef.value && typeof formRef.value.validate === 'function') {
    const valid = await formRef.value.validate()
    if (!valid) return
  }

  if (!form.value.topic || !form.value.date || !form.value.type || !form.value.subject_id) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }

  emit('save', {
    ...form.value,
    assignHomework: assignHomeworkToo.value
  })
}
</script>
