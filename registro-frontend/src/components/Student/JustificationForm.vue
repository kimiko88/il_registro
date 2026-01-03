<template>
  <q-dialog v-model="isOpen">
    <q-card style="min-width: 300px">
        <q-card-section>
            <div class="text-h6">Request Justification</div>
            <div class="text-caption">for absence on {{ date }}</div>
        </q-card-section>

        <q-card-section>
            <q-form @submit="onSubmit">
                <q-select 
                    v-model="reason" 
                    :options="['Illness', 'Family Reasons', 'Traffic', 'Other']" 
                    label="Reason" 
                    filled 
                    class="q-mb-md" 
                />
                <q-input 
                    v-model="notes" 
                    label="Additional Notes" 
                    filled 
                    type="textarea" 
                    rows="3" 
                />
                <div class="row justify-end q-mt-md">
                    <q-btn flat label="Cancel" v-close-popup />
                    <q-btn type="submit" label="Submit" color="primary" />
                </div>
            </q-form>
        </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch } from 'vue';

const props = defineProps(['modelValue', 'date']);
const emit = defineEmits(['update:modelValue', 'submit']);

const isOpen = ref(props.modelValue);
const reason = ref('Illness');
const notes = ref('');

watch(() => props.modelValue, (val) => isOpen.value = val);
watch(isOpen, (val) => emit('update:modelValue', val));

const onSubmit = () => {
    emit('submit', { reason: reason.value, notes: notes.value });
    isOpen.value = false;
};
</script>
