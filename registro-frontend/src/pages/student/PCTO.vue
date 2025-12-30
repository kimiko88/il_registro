<template>
  <q-page class="q-pa-md">
    <h1 class="text-h4 q-my-none q-mb-md">PCTO Dashboard</h1>
    
    <!-- Summary -->
    <q-card class="bg-primary text-white q-mb-lg">
        <q-card-section>
            <div class="text-h2 text-center">{{ store.totalHours }}h</div>
            <div class="text-subtitle1 text-center">Total Hours Completed</div>
        </q-card-section>
    </q-card>

    <h2 class="text-h5">My Projects</h2>
    <q-list bordered separator>
        <q-item v-for="proj in store.projects" :key="proj.id">
            <q-item-section>
                <q-item-label>{{ proj.title }}</q-item-label>
                <q-item-label caption>{{ proj.company }}</q-item-label>
            </q-item-section>
            <q-item-section side>
                <q-chip :color="proj.status === 'Completed' ? 'green' : 'orange'" text-color="white" size="sm">
                    {{ proj.hours }}h • {{ proj.status }}
                </q-chip>
            </q-item-section>
        </q-item>
    </q-list>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { usePCTOStore } from 'src/stores/pcto';

const store = usePCTOStore();

onMounted(() => {
    store.fetchProjects();
});
</script>
