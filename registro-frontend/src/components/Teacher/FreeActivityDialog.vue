<template>
  <q-dialog :model-value="modelValue" persistent @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(550px, 95vw); max-width: 95vw;" class="rounded-xl shadow-24" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-white text-slate-800'">
      <q-card-section class="row items-center justify-between" :class="$q.dark.isActive ? 'bg-grey-9 text-white' : 'bg-teal-8 text-white'">
        <div class="text-h6 text-weight-bold row items-center">
          <q-icon name="event_busy" class="q-mr-sm" size="24px" />
          {{ isEditing ? 'Modifica Attività Libera' : 'Nuova Attività Libera' }}
        </div>
        <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
      </q-card-section>

      <q-card-section class="q-gutter-y-md q-pt-md">
        <q-banner class="rounded-borders" :class="$q.dark.isActive ? 'bg-teal-10 text-teal-1' : 'bg-teal-1 text-teal-9'" dense>
          <template v-slot:avatar><q-icon name="info" color="teal-7" /></template>
          Usa questo form per registrare ore in cui non sei in classe: classe in gita, disponibilità, riunione, formazione, etc.
        </q-banner>

        <q-form ref="formRef" @submit.prevent="handleSubmit" class="q-gutter-y-md">
          <q-input
            v-model="form.date"
            type="date"
            :label="($t('classRegister.dateLabel') || 'Data') + ' *'"
            outlined
            dense
            :rules="[v => !!v || ($t('common.requiredField') || 'Campo obbligatorio')]"
          />

          <div class="row q-col-gutter-sm">
            <div class="col-12 col-md-6">
              <q-select
                v-model="form.start_hour"
                :options="[1,2,3,4,5,6,7,8]"
                label="Ora di Inizio *"
                outlined
                dense
                :rules="[v => !!v || ($t('common.requiredField') || 'Campo obbligatorio')]"
              >
                <template v-slot:prepend><q-icon name="schedule" /></template>
              </q-select>
            </div>
            <div class="col-12 col-md-6">
              <q-select
                v-model="form.duration"
                :options="[1,2,3,4,5,6]"
                label="Durata (ore) *"
                outlined
                dense
                :rules="[v => !!v || ($t('common.requiredField') || 'Campo obbligatorio')]"
              >
                <template v-slot:prepend><q-icon name="timelapse" /></template>
              </q-select>
            </div>
          </div>

          <q-select
            v-model="form.activity_type"
            :options="freeActivityTypeOptions"
            option-value="value"
            option-label="label"
            emit-value
            map-options
            label="Tipo di Attività *"
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

          <q-input
            v-model="form.description"
            label="Descrizione *"
            outlined
            dense
            autogrow
            placeholder="Es: Classe in gita a Roma — ore a disposizione, riunione dipartimento, etc."
            :rules="[v => (!!v && v.trim().length > 0) || ($t('common.requiredField') || 'Campo obbligatorio')]"
          />

          <q-input
            v-model="form.notes"
            label="Note aggiuntive"
            type="textarea"
            outlined
            dense
            autogrow
            placeholder="(Opzionale)"
          />

          <q-card-actions align="right" class="q-px-none q-pt-sm">
            <q-btn flat :label="$t('common.cancel') || 'Annulla'" v-close-popup color="grey-7" no-caps />
            <q-btn
              type="submit"
              color="teal"
              text-color="white"
              :label="isEditing ? ($t('common.update') || 'Aggiorna') : ($t('common.save') || 'Salva Attività')"
              :loading="saving"
              icon="save"
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
  isEditing: {
    type: Boolean,
    default: false
  },
  initialData: {
    type: Object,
    default: () => ({})
  },
  freeActivityTypeOptions: {
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
  date: '',
  start_hour: 1,
  duration: 1,
  activity_type: 'disposition',
  description: '',
  notes: ''
})

watch(() => props.initialData, (val) => {
  if (val) {
    form.value = {
      date: val.date || '',
      start_hour: val.start_hour || 1,
      duration: val.duration || 1,
      activity_type: val.activity_type || 'disposition',
      description: val.description || '',
      notes: val.notes || ''
    }
  }
}, { immediate: true, deep: true })

async function handleSubmit() {
  if (formRef.value && typeof formRef.value.validate === 'function') {
    const valid = await formRef.value.validate()
    if (!valid) return
  }

  if (!form.value.date || !form.value.description) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }

  emit('save', { ...form.value })
}
</script>
