<template>
  <q-dialog v-model="visible" persistent>
    <q-card style="min-width: 400px">
      <q-card-section>
        <div class="text-h6">{{ noteToEdit ? ($t('common.edit') || 'Modifica') : ($t('classRegister.addNote') || 'Nuova Nota') }}</div>
        <div class="text-subtitle2" v-if="student">{{ $t('classRegister.tableHeaderStudent') }}: {{ student.name || `${student.last_name || ''} ${student.first_name || ''}` }}</div>
      </q-card-section>

      <q-card-section>
        <q-form ref="form" @submit="onSubmit" class="q-gutter-md">
          <q-select
            v-model="noteData.type"
            :options="typeOptions"
            :label="($t('classRegister.lessonTypeLabel') || 'Tipo Nota') + ' *'"
            outlined
            emit-value
            map-options
            :rules="[val => !!val || 'Seleziona un tipo']"
          />

          <q-input
            v-model="noteData.note"
            :label="($t('classRegister.topicLabel') || 'Contenuto') + ' *'"
            type="textarea"
            outlined
            autogrow
            :rules="[val => !!val || 'Scrivi il contenuto']"
          />

          <q-input
             v-model="noteData.date"
             type="date"
             :label="$t('classRegister.dateLabel') || 'Data'"
             outlined
             :rules="[val => !!val || 'Data obbligatoria']"
          />

          <div class="row justify-end q-gutter-sm q-mt-md">
            <q-btn flat :label="$t('common.cancel') || 'Annulla'" color="grey" v-close-popup />
            <q-btn type="submit" :label="$t('common.save') || 'Salva'" color="primary" :loading="loading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import notesService from '@/services/notesService'

const props = defineProps({
  modelValue: Boolean,
  student: Object,
  classId: String,
  noteToEdit: Object
})

const emit = defineEmits(['update:modelValue', 'saved'])

const $q = useQuasar()
let t = (key, fallback) => (typeof fallback === 'string' ? fallback : key)
try {
  const i18nInstance = useI18n()
  if (i18nInstance && i18nInstance.t) {
    t = i18nInstance.t
  }
} catch (e) {
  // Fallback for isolated unit tests without vue-i18n app plugin
}
const loading = ref(false)

const noteData = reactive({
  type: 'generic',
  note: '',
  date: new Date().toISOString().split('T')[0]
})

const typeOptions = computed(() => [
  { label: t('classRegister.activityStandard') || 'Nota Generica', value: 'generic' },
  { label: t('classRegister.assignHomework') || 'Richiamo Compiti', value: 'homework' },
  { label: t('classRegister.activityStandardCap') || 'Nota Comportamentale', value: 'behavior' },
  { label: t('classRegister.addDisciplinaryNote') || 'Nota Disciplinare', value: 'disciplinary' }
])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

watch(() => props.modelValue, (val) => {
  if (val) {
    if (props.noteToEdit) {
      noteData.type = props.noteToEdit.type
      noteData.note = props.noteToEdit.note
      noteData.date = props.noteToEdit.date
    } else {
      // Reset form on open for new note
      noteData.type = 'generic'
      noteData.note = ''
      noteData.date = new Date().toISOString().split('T')[0]
    }
  }
})

const onSubmit = async () => {
  loading.value = true
  try {
    if (props.noteToEdit) {
      await notesService.updateNote(props.noteToEdit.id, {
        type: noteData.type,
        note: noteData.note,
        date: noteData.date
      })
      $q.notify({
        type: 'positive',
        message: 'Nota modificata con successo'
      })
    } else {
      const payload = {
        student_id: props.student ? props.student.id : '',
        class_id: props.classId,
        type: noteData.type,
        note: noteData.note,
        date: noteData.date
      }
      await notesService.createNote(payload)
      $q.notify({
        type: 'positive',
        message: 'Nota salvata con successo'
      })
    }
    
    emit('saved')
    visible.value = false
  } catch (error) {
    $q.notify({
      type: 'negative',
      message: 'Errore nel salvataggio della nota',
      caption: error.response?.data?.error || error.message
    })
  } finally {
    loading.value = false
  }
}
</script>
