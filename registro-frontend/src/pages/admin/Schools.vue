<template>
  <q-page padding>
    <div class="row q-mb-md">
      <div class="text-h4">{{ $t('nav.schools') || 'Gestione Scuole' }}</div>
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
import { useI18n } from 'vue-i18n';
import SchoolTable from '@/components/Admin/SchoolTable.vue';
import SchoolForm from '@/components/Admin/SchoolForm.vue';
import { useSchoolManagement } from '@/composables/useSchoolManagement';

const { t } = useI18n();
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
