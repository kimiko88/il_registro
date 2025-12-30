<template>
  <q-page class="q-pa-md">
    <h1 class="text-h4 q-my-none q-mb-md">My Documents</h1>
    <DocumentList :documents="documents" @download="downloadFile" v-if="!loading" />
    <div v-if="loading" class="text-center q-pa-lg">
        <q-spinner color="primary" size="2em" />
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useDocumentsStore } from 'src/stores/documents';
import { useMyDocuments } from 'src/composables/useMyDocuments';
import DocumentList from 'src/components/Student/DocumentList.vue';

const store = useDocumentsStore();
const { documents, loading, downloadFile } = useMyDocuments();

onMounted(() => {
    store.fetchMyFiles();
});
</script>
