<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
        <h1 class="text-h4 q-my-none">My Children</h1>
        <q-btn label="Add Child" color="primary" icon="add" />
    </div>

    <ChildrenList :children="children" @remove="removeChild" />

    <div class="q-mt-lg text-grey text-caption text-center">
        Use the "Add Child" button to link another student account to your profile using their Access Code.
    </div>
  </q-page>
</template>

<script setup>
import { useChildrenManagement } from 'src/composables/useChildrenManagement';
import ChildrenList from 'src/components/Parent/ChildrenList.vue';

const { children, removeChild } = useChildrenManagement();
</script>
