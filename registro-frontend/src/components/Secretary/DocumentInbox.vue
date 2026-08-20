<template>
  <q-table
    :rows="store.inbox"
    :columns="columns"
    row-key="id"
    v-model:pagination="store.pagination"
    :loading="store.loading"
    @request="onRequest"
    selection="multiple"
    v-model:selected="selected"
    :filter="filter"
  >
    <template v-slot:top>
        <div class="text-h6 q-mr-md">{{ t('documentsPage.inbox') || 'Documenti in Arrivo' }}</div>
        
        <q-btn-toggle
            v-model="typeFilter"
            push
            glossy
            toggle-color="primary"
            class="q-mr-md"
            :options="[
                {label: t('common.all') || 'Tutti', value: 'all'},
                {label: 'PDP', value: 'PDP'},
                {label: 'PFI', value: 'PFI'},
                {label: 'PCTO', value: 'PCTO'},
                {label: t('documentsPage.certificates') || 'Certificati', value: 'Certificate'}
            ]"
        />
        
        <q-space />
        <q-input dense debounce="300" v-model="filter" :placeholder="t('documentsPage.searchPlaceholder') || 'Cerca documento...'">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
    </template>

    <!-- Batch Actions -->
    <template v-slot:top-row v-if="selected.length > 0">
        <q-tr>
            <q-td colspan="100%">
                 <div class="row items-center q-gutter-sm bg-blue-1 q-pa-sm rounded-borders">
                   <span class="text-weight-bold text-primary">{{ selected.length }} {{ t('common.selected') || 'selezionati' }}</span>
                   <q-space />
                   <q-btn size="sm" color="positive" icon="check_circle" :label="t('documentsPage.approveSelected') || 'Approva Selezionati'" @click="batchApprove" />
                   <q-btn size="sm" color="orange" icon="archive" :label="t('documentsPage.archive') || 'Archivia'" @click="batchArchive" />
                 </div>
            </q-td>
        </q-tr>
    </template>

    <template v-slot:body-cell-favorite="props">
        <q-td :props="props" auto-width>
            <q-btn flat round dense :icon="props.row.favorite ? 'star' : 'star_border'" :color="props.row.favorite ? 'amber' : 'grey'" :aria-label="t('documentsPage.favorite') || 'Preferito'" @click="toggleFavorite(props.row)" />
        </q-td>
    </template>

    <template v-slot:body-cell-status="props">
        <q-td :props="props">
            <q-badge :color="getStatusColor(props.value)" :label="props.value" />
        </q-td>
    </template>

    <template v-slot:body-cell-actions="props">
      <q-td :props="props" auto-width>
        <q-btn flat round icon="visibility" :aria-label="t('common.preview') || 'Anteprima'" @click="$emit('preview', props.row)">
            <q-tooltip>{{ t('common.preview') || 'Anteprima' }}</q-tooltip>
        </q-btn>
        <q-btn flat round color="primary" icon="rate_review" :aria-label="t('common.review') || 'Revisiona'" @click="$emit('review', props.row)" v-if="props.row.status === 'Pending'">
            <q-tooltip>{{ t('common.review') || 'Revisiona' }}</q-tooltip>
        </q-btn>
        <q-btn flat round color="grey" icon="archive" :aria-label="t('documentsPage.archive') || 'Archivia'" @click="$emit('archive', props.row)">
            <q-tooltip>{{ t('documentsPage.archive') || 'Archivia' }}</q-tooltip>
        </q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useDocumentsStore } from '@/stores/documents';
import { useQuasar } from 'quasar';

const { t } = useI18n();
const store = useDocumentsStore();
const $q = useQuasar();
const filter = ref('');
const typeFilter = ref('all');
const selected = ref([]);

const columns = computed(() => [
  { name: 'favorite', label: '', field: 'favorite', align: 'center', sortable: true },
  { name: 'title', label: t('documentsPage.titleLabel') || 'Titolo', field: 'title', align: 'left', sortable: true },
  { name: 'type', label: t('classRegister.tableHeaderGradeType') || 'Tipo', field: 'type', align: 'left', sortable: true },
  { name: 'student', label: t('competenciesPage.student') || 'Studente', field: 'student', align: 'left' },
  { name: 'status', label: t('substitutionsPage.status') || 'Stato', field: 'status', align: 'center', sortable: true },
  { name: 'date', label: t('classRegister.dateLabel') || 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'actions', label: t('common.actions') || 'Azioni', align: 'right' }
]);

const emit = defineEmits(['preview', 'review', 'request', 'archive']);

const onRequest = (props) => emit('request', props);

const getStatusColor = (status) => {
    if (status === 'Pending') return 'orange';
    if (status === 'Approved') return 'green';
    if (status === 'Rejected') return 'red';
    if (status === 'Draft') return 'grey';
    return 'blue';
};

const toggleFavorite = (row) => {
    row.favorite = !row.favorite;
    $q.notify({ 
        message: row.favorite ? (t('documentsPage.favoriteAdded') || 'Aggiunto ai preferiti') : (t('documentsPage.favoriteRemoved') || 'Rimosso dai preferiti'), 
        icon: 'star', 
        color: row.favorite ? 'amber' : 'grey-8' 
    });
};

const batchApprove = () => {
    $q.dialog({
        title: t('documentsPage.batchApproveTitle') || 'Approvazione Multipla',
        message: t('documentsPage.batchApproveConfirm', { count: selected.value.length }) || `Approvare ${selected.value.length} documenti?`,
        cancel: true
    }).onOk(() => {
        selected.value.forEach(d => d.status = 'Approved');
        selected.value = [];
        $q.notify({ type: 'positive', message: t('documentsPage.documentsApproved') || 'Documenti approvati' });
    })
};

const batchArchive = () => {
     $q.notify({ message: t('documentsPage.documentsArchived') || 'Documenti archiviati', icon: 'archive' });
     selected.value = [];
}

defineExpose({
    getStatusColor,
    toggleFavorite,
    batchApprove,
    batchArchive,
    selected,
    filter,
    typeFilter
})
</script>
