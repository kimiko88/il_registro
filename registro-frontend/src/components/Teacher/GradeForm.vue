<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card style="min-width: 400px">
      <q-card-section>
        <div class="text-h6">Enter Grade</div>
        <div class="text-subtitle2">{{ studentName }}</div>
      </q-card-section>

      <q-card-section>
        <q-form class="q-gutter-md">
          <q-select
            v-model="form.type"
            :options="['Written', 'Oral', 'Practical']"
            label="Type"
            outlined
            dense
          />
          <q-input
            v-model.number="form.value"
            type="number"
            label="Grade (1-10)"
            outlined
            dense
            :rules="[val => val >= 1 && val <= 10 || 'Invalid grade']"
          />
          <q-input
            v-model="form.date"
            type="date"
            label="Date"
            outlined
            dense
          />
          <q-input
            v-model="form.description"
            label="Description / Topic"
            outlined
            dense
          />
        </q-form>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Cancel" color="primary" v-close-popup />
        <q-btn flat label="Save" color="primary" @click="save" :loading="submitting" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useGradeEntry } from 'src/composables/useGradeEntry';
import { date } from 'quasar';

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
