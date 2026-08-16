<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card style="min-width: 400px">
      <q-card-section>
        <div class="text-h6">{{ t('classRegister.tableHeaderGrade') || 'Inserisci Voto' }}</div>
        <div class="text-subtitle2">{{ studentName }}</div>
      </q-card-section>

      <q-card-section>
        <q-form class="q-gutter-md">
          <q-select
            v-model="form.type"
            :options="typeOptions"
            :label="t('classRegister.tableHeaderGradeType') || 'Tipo Voto'"
            outlined
            dense
          />
          <q-input
            v-model.number="form.value"
            type="number"
            :label="t('classRegister.tableHeaderGrade') || 'Voto (1-10)'"
            outlined
            dense
            :rules="[val => val >= 1 && val <= 10 || 'Voto non valido']"
          />
          <q-input
            v-model="form.date"
            type="date"
            :label="t('classRegister.dateLabel') || 'Data'"
            outlined
            dense
          />
          <q-input
            v-model="form.description"
            :label="t('classRegister.tableHeaderGradeNotes') || 'Descrizione / Argomento'"
            outlined
            dense
          />
        </q-form>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="t('common.cancel') || 'Annulla'" color="primary" v-close-popup />
        <q-btn flat :label="t('common.save') || 'Salva'" color="primary" @click="save" :loading="submitting" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useGradeEntry } from '@/composables/useGradeEntry';
import { date } from 'quasar';

const { t } = useI18n();

const typeOptions = computed(() => [
  { label: t('agendaPage.test') || 'Scritto', value: 'Written' },
  { label: t('agendaPage.oralTest') || 'Orale', value: 'Oral' },
  { label: t('classRegister.activityLab') || 'Pratico', value: 'Practical' }
])

const props = defineProps({
  modelValue: Boolean,
  studentId: String,
  studentName: String
});

const emit = defineEmits(['update:modelValue', 'saved']);

const isOpen = ref(false);
watch(() => props.modelValue, (val) => isOpen.value = val);
watch(isOpen, (val) => emit('update:modelValue', val));

const form = ref({
  type: 'Written',
  value: null,
  date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
  description: ''
});

const { submitGrade, submitting } = useGradeEntry();

const save = async () => {
  const success = await submitGrade({
    studentId: props.studentId,
    ...form.value
  });
  
  if (success) {
    isOpen.value = false;
    emit('saved');
    // Reset form
    form.value.value = null;
    form.value.description = '';
  }
};
</script>
