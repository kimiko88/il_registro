<template>
  <q-dialog :model-value="modelValue" persistent @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(500px, 95vw); max-width: 95vw;" class="rounded-xl shadow-24" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-white text-slate-800'">
      <q-card-section class="row items-center justify-between" :class="$q.dark.isActive ? 'bg-grey-9 text-white' : 'bg-orange-8 text-white'">
        <div class="text-h6 text-weight-bold row items-center">
          <q-icon name="assignment" class="q-mr-sm" size="24px" />
          {{ $t('homework.assignHomework') || 'Assegna Compito' }}
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
            v-model="form.description"
            :label="($t('homework.description') || 'Descrizione Compito') + ' *'"
            type="textarea"
            outlined
            dense
            autogrow
            placeholder="Es: Pagine 145-150, esercizi 12-18"
            :rules="[v => (!!v && v.trim().length > 0) || ($t('common.requiredField') || 'Campo obbligatorio')]"
          />

          <q-input
            v-model="form.dueDate"
            type="date"
            :label="($t('homework.dueDate') || 'Data di Consegna') + ' *'"
            outlined
            dense
            :rules="[v => !!v || ($t('common.requiredField') || 'Campo obbligatorio')]"
          />

          <q-card-actions align="right" class="q-px-none q-pt-sm">
            <q-btn flat :label="$t('common.cancel') || 'Annulla'" v-close-popup color="grey-7" no-caps />
            <q-btn
              type="submit"
              color="orange-9"
              text-color="white"
              :label="$t('homework.assign') || 'Assegna'"
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
import { ref, watch } from 'vue'
import { useQuasar } from 'quasar'

const props = defineProps({
  modelValue: {
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
  saving: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'save'])
const $q = useQuasar()
const formRef = ref(null)

const form = ref({
  subject_id: null,
  description: '',
  dueDate: ''
})

watch(() => props.initialData, (val) => {
  if (val) {
    form.value = {
      subject_id: val.subject_id || null,
      description: val.description || '',
      dueDate: val.dueDate || ''
    }
  }
}, { immediate: true, deep: true })

async function handleSubmit() {
  if (formRef.value && typeof formRef.value.validate === 'function') {
    const valid = await formRef.value.validate()
    if (!valid) return
  }

  if (!form.value.description || !form.value.dueDate || !form.value.subject_id) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }

  emit('save', { ...form.value })
}
</script>
