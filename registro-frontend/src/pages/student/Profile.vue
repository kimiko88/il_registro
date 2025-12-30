<template>
  <q-page class="q-pa-md">
    <h1 class="text-h4 q-my-none q-mb-lg">Profile</h1>

    <div v-if="store.profile">
        <div class="row q-col-gutter-md">
            <div class="col-12 col-md-4 text-center">
                <q-card>
                    <q-card-section>
                        <q-avatar size="100px" class="q-mb-md">
                            <img :src="store.profile.avatar || 'https://cdn.quasar.dev/img/boy-avatar.png'" />
                        </q-avatar>
                        <div class="text-h6">{{ store.fullName }}</div>
                        <div class="text-subtitle2 text-grey">{{ store.className }}</div>
                    </q-card-section>
                </q-card>
            </div>
            
            <div class="col-12 col-md-8">
                <q-card>
                    <q-card-section>
                        <div class="text-h6 q-mb-md">Personal Info</div>
                        <q-list separator>
                            <q-item>
                                <q-item-section>Email</q-item-section>
                                <q-item-section side>{{ store.profile.email }}</q-item-section>
                            </q-item>
                            <q-item>
                                <q-item-section>Student ID</q-item-section>
                                <q-item-section side>{{ store.profile.id }}</q-item-section>
                            </q-item>
                        </q-list>
                    </q-card-section>
                    <q-card-actions align="right">
                        <q-btn flat label="Change Password" color="primary" />
                        <q-btn flat label="Privacy Settings" color="primary" />
                    </q-card-actions>
                </q-card>
            </div>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useStudentStore } from 'src/stores/student';

const store = useStudentStore();

onMounted(() => {
    if (!store.profile) store.fetchProfile();
});
</script>
