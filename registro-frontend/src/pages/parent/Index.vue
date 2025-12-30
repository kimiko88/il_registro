<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h5">Parent Dashboard</div>
       <ChildrenSelector />
    </div>

    <!-- Child Summary -->
    <div v-if="childrenStore.selectedChild" class="row q-col-gutter-md q-mb-lg">
        <div class="col-12">
            <q-card class="bg-primary text-white">
                <q-card-section class="row items-center">
                    <q-avatar size="60px" class="q-mr-md">
                        <img :src="childrenStore.selectedChild.avatar" />
                    </q-avatar>
                    <div>
                        <div class="text-h6">{{ childrenStore.selectedChild.firstName }} {{ childrenStore.selectedChild.lastName }}</div>
                        <div class="text-subtitle2">{{ childrenStore.selectedChild.className }} • {{ childrenStore.selectedChild.school }}</div>
                    </div>
                </q-card-section>
            </q-card>
        </div>

        <!-- Quick Stats for Selected Child -->
        <div class="col-6 col-sm-4">
             <q-card class="text-center q-py-sm cursor-pointer" @click="$router.push('/parent/grades')">
                 <div class="text-h6 text-primary">Grades</div>
                 <q-icon name="analytics" size="md" color="primary" />
             </q-card>
        </div>
        <div class="col-6 col-sm-4">
             <q-card class="text-center q-py-sm cursor-pointer" @click="$router.push('/parent/attendance')">
                 <div class="text-h6 text-orange">Attendance</div>
                 <q-icon name="schedule" size="md" color="orange" />
             </q-card>
        </div>
         <div class="col-12 col-sm-4">
             <q-card class="text-center q-py-sm cursor-pointer" @click="$router.push('/parent/colloqui')">
                 <div class="text-h6 text-purple">Book Meeting</div>
                 <q-icon name="calendar_today" size="md" color="purple" />
             </q-card>
        </div>
    </div>

    <!-- Recent Activity / Notifications -->
    <div class="text-h6 q-mb-sm">Notifications</div>
    <q-list bordered separator class="bg-white rounded-borders">
        <q-item v-for="notif in parentStore.notifications" :key="notif.id">
            <q-item-section avatar>
                <q-icon :name="notif.type === 'warning' ? 'warning' : 'info'" :color="notif.type === 'warning' ? 'orange' : 'blue'" />
            </q-item-section>
            <q-item-section>
                <q-item-label>{{ notif.title }}</q-item-label>
                <q-item-label caption>{{ notif.message }}</q-item-label>
            </q-item-section>
            <q-item-section side>
                <q-item-label caption>{{ notif.date }}</q-item-label>
            </q-item-section>
        </q-item>
    </q-list>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useParentStore } from 'src/stores/parent';
import { useChildrenStore } from 'src/stores/children';
import ChildrenSelector from 'src/components/Parent/ChildrenSelector.vue';

const parentStore = useParentStore();
const childrenStore = useChildrenStore();

onMounted(() => {
    parentStore.fetchProfile();
    parentStore.fetchNotifications();
    childrenStore.fetchChildren();
});
</script>
