<template>
  <q-card>
    <q-card-section>
        <div class="text-h6">Create Availability</div>
    </q-card-section>
    <q-card-section>
        <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input v-model="form.date" type="date" label="Date" filled />
            <div class="row q-gutter-md">
                <q-input v-model="form.startTime" type="time" label="Start" filled class="col" />
                <q-input v-model="form.endTime" type="time" label="End" filled class="col" />
            </div>
            <div class="row q-gutter-md">
                <q-input v-model.number="form.duration" type="number" label="Duration (min)" filled class="col" />
                <q-input v-model.number="form.break" type="number" label="Break (min)" filled class="col" />
            </div>
            <q-btn type="submit" label="Generate Slots" color="primary" :loading="generating" />
        </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref } from 'vue';
import { useColloquiScheduling } from 'src/composables/useColloquiScheduling';
import { date } from 'quasar';

const { generateSlots, generating } = useColloquiScheduling();

const form = ref({
    date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
    startTime: '15:00',
    endTime: '17:00',
    duration: 15,
    break: 0
});

const onSubmit = async () => {
    await generateSlots(form.value);
};
</script>
