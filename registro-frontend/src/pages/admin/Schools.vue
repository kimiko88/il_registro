<template>
  <q-page padding>
    <div class="row q-mb-md">
      <div class="text-h4">School Management</div>
    </div>

    <SchoolTable 
      @create="openCreate"
      @edit="openEdit"
      @delete="confirmDelete"
      @request="onRequest"
    />

    <q-dialog v-model="showDialog">
      <SchoolForm :school="editingSchool" @submit="submitSchool" />
    </q-dialog>
  </q-page>
</template>

<script setup>
import SchoolTable from 'src/components/Admin/SchoolTable.vue';
import SchoolForm from 'src/components/Admin/SchoolForm.vue';
import { useSchoolManagement } from 'src/composables/useSchoolManagement';

const { 
  showDialog, 
  editingSchool, 
  openCreate, 
  openEdit, 
  submitSchool, 
  confirmDelete, 
  onRequest 
} = useSchoolManagement();
</script>
