<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-md">Colloqui Management</div>

    <div class="row q-col-gutter-lg">
        <!-- Calendar / List -->
        <div class="col-12 col-md-8">
            <q-card>
                <q-toolbar class="bg-primary text-white">
                    <q-toolbar-title>Schedule</q-toolbar-title>
                    <q-btn flat icon="chevron_left" />
                    <div class="text-subtitle1 q-mx-sm">January 2025</div>
                    <q-btn flat icon="chevron_right" />
                </q-toolbar>
                <q-card-section>
                    <div class="text-h6 q-mb-sm">Slots for {{ selectedDate }}</div>
                    <ColloquioBookingList :slots="store.slots" @delete="store.deleteSlot" />
                </q-card-section>
            </q-card>
        </div>

        <!-- Slot Manager -->
        <div class="col-12 col-md-4">
            <ColloquioSlotManager />
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useColloquiStore } from 'src/stores/colloqui';
import ColloquioSlotManager from 'src/components/Teacher/ColloquioSlotManager.vue';
import ColloquioBookingList from 'src/components/Teacher/ColloquioBookingList.vue';
import { date } from 'quasar';

const store = useColloquiStore();
const selectedDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'));

onMounted(() => {
    store.fetchSlots();
});
</script>
