<template>
  <q-dialog v-model="visible" persistent>
    <q-card style="min-width: 400px">
      <q-card-section>
        <div class="text-h6">Nuova Nota</div>
        <div class="text-subtitle2" v-if="student">Studente: {{ student.name }}</div>
      </q-card-section>

      <q-card-section>
        <q-form ref="form" @submit="onSubmit" class="q-gutter-md">
          <q-select
            v-model="noteData.type"
            :options="typeOptions"
            label="Tipo Nota *"
            outlined
            emit-value
            map-options
            :rules="[val => !!val || 'Seleziona un tipo']"
          />

          <q-input
            v-model="noteData.note"
            label="Contenuto *"
            type="textarea"
            outlined
            autogrow
            :rules="[val => !!val || 'Scrivi il contenuto']"
          />

          <q-input
             v-model="noteData.date"
             type="date"
             label="Data"
             outlined
             :rules="[val => !!val || 'Data obbligatoria']"
          />

          <div class="row justify-end q-gutter-sm q-mt-md">
            <q-btn flat label="Annulla" color="grey" v-close-popup />
            <q-btn type="submit" label="Salva" color="primary" :loading="loading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useQuasar } from 'quasar'
import notesService from '@/services/notesService'

const props = defineProps({
  modelValue: Boolean,
  student: Object,
  classId: String,
  noteToEdit: Object
})

const emit = defineEmits(['update:modelValue', 'saved'])

const $q = useQuasar()
const loading = ref(false)

const noteData = reactive({
  type: 'generic',
  note: '',
  date: new Date().toISOString().split('T')[0]
})

const typeOptions = [
  { label: 'Nota Generica', value: 'generic' },
  { label: 'Richiamo Compiti', value: 'homework' },
  { label: 'Nota Comportamentale', value: 'behavior' },
  { label: 'Nota Disciplinare', value: 'disciplinary' }
]

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
