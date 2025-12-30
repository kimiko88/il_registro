<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-lg">Admin Dashboard</div>
    
    <q-tabs
      v-model="tab"
      dense
      class="text-grey"
      active-color="primary"
      indicator-color="primary"
      align="justify"
      narrow-indicator
    >
      <q-tab name="users" label="Users" />
      <q-tab name="schools" label="Schools" />
    </q-tabs>

    <q-separator />

    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="users">
        <UserManagement />
      </q-tab-panel>

      <q-tab-panel name="schools">
        <SchoolManagement />
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import UserManagement from '@/components/Admin/UserManagement.vue'
import SchoolManagement from '@/components/Admin/SchoolManagement.vue'

const tab = ref('users')
</script>
