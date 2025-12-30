<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
        <h1 class="text-h4 q-my-none">Book Meeting</h1>
        <ChildrenSelector />
    </div>

    <!-- Tabs: Book vs My Bookings -->
    <q-tabs v-model="tab" class="text-primary q-mb-md" align="left">
        <q-tab name="book" label="Available Slots" />
        <q-tab name="my" label="My Bookings" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
        <q-tab-panel name="book">
            <!-- Teacher List Mock -->
            <q-list bordered separator>
                <q-item-label header>Select Teacher</q-item-label>
                <q-item clickable v-ripple @click="confirmBook">
                    <q-item-section avatar>
                         <q-avatar color="primary" text-color="white">B</q-avatar>
                    </q-item-section>
                    <q-item-section>
                        <q-item-label>Prof. Bianchi</q-item-label>
                        <q-item-label caption>Mathematics</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <q-chip color="green" text-color="white" size="sm">3 Slots Available</q-chip>
                    </q-item-section>
                </q-item>
            </q-list>
        </q-tab-panel>

        <q-tab-panel name="my">
            <q-list bordered separator>
                <q-item v-for="booking in bookings" :key="booking.id">
                    <q-item-section>
                        <q-item-label>{{ booking.teacherName }} - {{ booking.subject }}</q-item-label>
                        <q-item-label caption>{{ booking.date }}</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <q-btn flat round icon="delete" color="red" @click="cancelBooking(booking.id)" />
                    </q-item-section>
                </q-item>
                 <div v-if="bookings.length === 0" class="text-center text-grey q-pa-lg">
                    No upcoming meetings.
                </div>
            </q-list>
        </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref } from 'vue';
import { useColloquiBooking } from 'src/composables/useColloquiBooking';
import ChildrenSelector from 'src/components/Parent/ChildrenSelector.vue';

const tab = ref('book');
const { bookings, bookSlot, cancelBooking } = useColloquiBooking();

const confirmBook = async () => {
    await bookSlot('mock-slot-id');
};
</script>
