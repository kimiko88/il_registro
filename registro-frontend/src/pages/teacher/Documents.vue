<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div class="row items-center">
        <h1 class="text-h4 q-my-none q-mr-md">Documents</h1>
        <ClassSelector />
      </div>
      <div>
        <q-btn color="primary" label="New Document" icon="add" @click="showTemplateDialog = true" />
      </div>
    </div>

    <!-- Document List -->
    <q-table
        :rows="documentsStore.inbox"
        :columns="columns"
        row-key="id"
        flat bordered
        :loading="documentsStore.loading"
    >
        <template v-slot:body-cell-status="props">
            <q-td :props="props">
                <q-badge :color="getStatusColor(props.value)">{{ props.value }}</q-badge>
            </q-td>
        </template>
        <template v-slot:body-cell-actions="props">
            <q-td :props="props">
                <q-btn flat round icon="edit" size="sm" color="primary" />
                <q-btn flat round icon="print" size="sm" color="grey" />
            </q-td>
        </template>
    </q-table>

    <!-- Dialogs -->
    <DocumentTemplate v-model="showTemplateDialog" @selected="onTemplateSelected" />
    
    <q-dialog v-model="showEditor" persistent maximized>
        <q-card>
            <q-bar class="bg-primary text-white">
                <div class="text-h6">Edit Document</div>
                <q-space />
                <q-btn flat round icon="close" v-close-popup />
            </q-bar>
            
            <q-card-section>
                <div class="text-subtitle1 q-mb-sm">Student: {{ selectedStudent?.name || 'Select Student' }}</div>
                <DocumentEditor v-model="editorContent" />
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Save Draft" @click="save" :loading="isSaving" />
                <q-btn color="primary" label="Publish" />
            </q-card-actions>
        </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useDocumentsStore } from 'src/stores/documents';
import { useDocumentCreation } from 'src/composables/useDocumentCreation';
import ClassSelector from 'src/components/Teacher/ClassSelector.vue';
import DocumentTemplate from 'src/components/Teacher/DocumentTemplate.vue';
import DocumentEditor from 'src/components/Teacher/DocumentEditor.vue';

const documentsStore = useDocumentsStore();
const { createDraft, isSaving } = useDocumentCreation();

const showTemplateDialog = ref(false);
const showEditor = ref(false);
const editorContent = ref('');
const selectedStudent = ref(null); // Mock

const columns = [
    { name: 'title', label: 'Title', field: 'title', align: 'left' },
    { name: 'type', label: 'Type', field: 'type', align: 'left' },
    { name: 'status', label: 'Status', field: 'status', align: 'left' },
    { name: 'date', label: 'Date', field: 'date', align: 'left' },
    { name: 'actions', label: 'Actions', align: 'right' }
];

const getStatusColor = (status) => {
    return status === 'Approved' ? 'green' : (status === 'Draft' ? 'grey' : 'orange');
};

const onTemplateSelected = (tpl) => {
    editorContent.value = tpl.content;
    showEditor.value = true;
};

const save = async () => {
    await createDraft('t1', 's1', editorContent.value);
    showEditor.value = false;
    documentsStore.fetchMyDocuments();
};

onMounted(() => {
    documentsStore.fetchMyDocuments();
});
</script>
