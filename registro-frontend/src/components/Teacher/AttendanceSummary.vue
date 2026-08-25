<template>
  <div class="row q-col-gutter-sm q-mb-md">
    <div class="col-4">
        <q-card class="bg-green-1 text-green-9 text-center">
            <q-card-section>
                <div class="text-h6">{{ store.presentCount }}</div>
                <div class="text-caption">{{ t('classRegister.present') }}</div>
            </q-card-section>
        </q-card>
    </div>
    <div class="col-4">
        <q-card class="bg-red-1 text-red-9 text-center">
            <q-card-section>
                <div class="text-h6">{{ store.absentCount }}</div>
                <div class="text-caption">{{ t('classRegister.absent') }}</div>
            </q-card-section>
        </q-card>
    </div>
    <div class="col-4">
        <q-card class="bg-orange-1 text-orange-9 text-center">
            <q-card-section>
                <div class="text-h6">{{ store.lateCount }}</div>
                <div class="text-caption">{{ t('classRegister.late') }}</div>
            </q-card-section>
        </q-card>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import { useAttendanceStore } from '@/stores/attendance';

const { t } = useI18n();
const store = useAttendanceStore();
</script>
