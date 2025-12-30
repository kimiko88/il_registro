<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">Document Inbox</div>
    
    <DocumentInbox 
        @preview="selectedPreviewDoc = $event; showPreview = true"
        @review="openReview"
        @request="store.fetchInbox"
    />

    <!-- Preview Dialog -->
    <q-dialog v-model="showPreview">
        <DocumentPreview :doc="selectedPreviewDoc" />
    </q-dialog>

    <!-- Review Dialog -->
    <q-dialog v-model="showReviewDialog">
        <DocumentReviewForm 
            :doc="selectedDoc"
            v-model="reviewNotes" 
            @submit="submitReview" 
        />
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue';
import DocumentInbox from 'src/components/Secretary/DocumentInbox.vue';
import DocumentPreview from 'src/components/Secretary/DocumentPreview.vue';
import DocumentReviewForm from 'src/components/Secretary/DocumentReviewForm.vue';
import { useDocumentProcessing } from 'src/composables/useDocumentProcessing';

const { 
    store,
    showReviewDialog,
    selectedDoc,
    reviewNotes,
    openReview,
    submitReview
} = useDocumentProcessing();

const showPreview = ref(false);
const selectedPreviewDoc = ref(null);
</script>
