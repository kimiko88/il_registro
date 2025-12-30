<template>
  <q-page class="q-pa-md">
    <h1 class="text-h4 q-my-none q-mb-md">Orientamento</h1>
    
    <div class="text-subtitle1 q-mb-md">Upcoming Events</div>
    <div class="row q-col-gutter-md">
        <div class="col-12 col-md-6" v-for="evt in store.events" :key="evt.id">
            <q-card>
                <q-card-section>
                    <div class="text-h6">{{ evt.title }}</div>
                    <div class="text-subtitle2 text-primary">{{ evt.date }} • {{ evt.hours }}h</div>
                    <div class="text-caption text-grey">{{ evt.location }}</div>
                </q-card-section>
                <q-card-actions align="right">
                    <q-btn flat label="Register" color="primary" />
                </q-card-actions>
            </q-card>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useOrientamentoStore } from 'src/stores/orientamento';

const store = useOrientamentoStore();

onMounted(() => {
    store.fetchEvents();
});
</script>
