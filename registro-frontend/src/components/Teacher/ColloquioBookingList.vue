<template>
  <q-list bordered separator>
    <q-item v-for="slot in slots" :key="slot.id">
        <q-item-section>
            <q-item-label>{{ slot.startTime }} - {{ slot.endTime }}</q-item-label>
            <q-item-label caption v-if="slot.booked" class="text-positive">Booked</q-item-label>
            <q-item-label caption v-else class="text-grey">Available</q-item-label>
        </q-item-section>
        
        <q-item-section v-if="slot.booked">
            <q-item-label>{{ getBooking(slot.id)?.studentName }}</q-item-label>
            <q-item-label caption>Parent: {{ getBooking(slot.id)?.parentName }}</q-item-label>
        </q-item-section>

        <q-item-section side>
            <q-btn flat round icon="delete" color="red" size="sm" @click="$emit('delete', slot.id)" :disable="slot.booked" />
        </q-item-section>
    </q-item>
  </q-list>
</template>

<script setup>
import { useColloquiStore } from 'src/stores/colloqui';

const props = defineProps(['slots']);
const store = useColloquiStore();

const getBooking = (slotId) => store.bookings.find(b => b.slotId === slotId);
</script>
