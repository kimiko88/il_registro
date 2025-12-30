<template>
  <q-page class="q-pa-md">
    <h1 class="text-h4 q-my-none q-mb-md">Profile & Settings</h1>
    
    <div v-if="store.profile">
        <q-card class="q-mb-md">
            <q-card-section>
                <div class="row items-center">
                    <q-avatar size="70px" class="q-mr-md">
                        <img :src="store.profile.avatar || 'https://cdn.quasar.dev/img/avatar2.jpg'" />
                    </q-avatar>
                    <div>
                        <div class="text-h6">{{ store.fullName }}</div>
                        <div class="text-subtitle2">{{ store.profile.email }}</div>
                    </div>
                </div>
            </q-card-section>
             <q-separator />
            <q-card-actions>
                <q-btn flat label="Edit Contact" color="primary" />
            </q-card-actions>
        </q-card>

         <q-list bordered separator class="bg-white rounded-borders">
            <q-item-label header>Preferences</q-item-label>
            <q-item>
                <q-item-section>
                    <q-item-label>Email Notifications</q-item-label>
                </q-item-section>
                <q-item-section side>
                    <q-toggle v-model="emailNotif" />
                </q-item-section>
            </q-item>
             <q-item>
                <q-item-section>
                    <q-item-label>SMS Alerts</q-item-label>
                </q-item-section>
                <q-item-section side>
                    <q-toggle v-model="smsNotif" />
                </q-item-section>
            </q-item>
         </q-list>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useParentStore } from 'src/stores/parent';

const store = useParentStore();
const emailNotif = ref(true);
const smsNotif = ref(false);

onMounted(() => {
    store.fetchProfile();
});
</script>
