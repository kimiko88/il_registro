<template>
  <q-card>
    <q-card-section>
        <div class="text-h6">{{ t('colloquiPage.createSlotsTitle') || 'Crea Disponibilità Colloqui' }}</div>
    </q-card-section>
    <q-card-section>
        <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input v-model="form.date" type="date" :label="t('classRegister.dateLabel') || 'Data'" filled />
            <div class="row q-gutter-md">
                <q-input v-model="form.startTime" type="time" :label="t('colloquiPage.startTime') || 'Ora Inizio'" filled class="col" />
                <q-input v-model="form.endTime" type="time" :label="t('colloquiPage.endTime') || 'Ora Fine'" filled class="col" />
            </div>
            <div class="row q-gutter-md">
                <q-input v-model.number="form.duration" type="number" :label="t('colloquiPage.duration') || 'Durata (min)'" filled class="col" />
                <q-input v-model.number="form.break" type="number" :label="t('colloquiPage.break') || 'Pausa (min)'" filled class="col" />
            </div>
            <q-btn type="submit" :label="t('colloquiPage.generateSlots') || 'Genera Slot'" color="primary" :loading="generating" />
        </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useColloquiScheduling } from '@/composables/useColloquiScheduling';
import { date } from 'quasar';

const { t } = useI18n();
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
