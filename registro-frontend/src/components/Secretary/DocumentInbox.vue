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
        <div class="text-h6 q-mr-md">Documenti in Arrivo</div>
        
        <q-btn-toggle
            v-model="typeFilter"
            push
            glossy
            toggle-color="primary"
            class="q-mr-md"
            :options="[
                {label: 'Tutti', value: 'all'},
                {label: 'PDP', value: 'PDP'},
                {label: 'PFI', value: 'PFI'},
                {label: 'PCTO', value: 'PCTO'},
                {label: 'Certificati', value: 'Certificate'}
            ]"
        />
        
        <q-space />
        <q-input dense debounce="300" v-model="filter" placeholder="Cerca documento...">
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
                   <span class="text-weight-bold text-primary">{{ selected.length }} selezionati</span>
                   <q-space />
                   <q-btn size="sm" color="positive" icon="check_circle" label="Approva Selezionati" @click="batchApprove" />
                   <q-btn size="sm" color="orange" icon="archive" label="Archivia" @click="batchArchive" />
                 </div>
            </q-td>
        </q-tr>
    </template>

    <template v-slot:body-cell-favorite="props">
        <q-td :props="props" auto-width>
            <q-btn flat round dense :icon="props.row.favorite ? 'star' : 'star_border'" :color="props.row.favorite ? 'amber' : 'grey'" @click="toggleFavorite(props.row)" />
        </q-td>
    </template>

    <template v-slot:body-cell-status="props">
        <q-td :props="props">
            <q-badge :color="getStatusColor(props.value)" :label="props.value" />
        </q-td>
    </template>

    <template v-slot:body-cell-actions="props">
      <q-td :props="props" auto-width>
        <q-btn flat round icon="visibility" @click="$emit('preview', props.row)" tooltip="Anteprima">
            <q-tooltip>Anteprima</q-tooltip>
        </q-btn>
        <q-btn flat round color="primary" icon="rate_review" @click="$emit('review', props.row)" v-if="props.row.status === 'Pending'">
            <q-tooltip>Revisiona</q-tooltip>
        </q-btn>
        <q-btn flat round color="grey" icon="archive" @click="$emit('archive', props.row)" />
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { ref } from 'vue';
import { useDocumentsStore } from 'src/stores/documents';
import { useQuasar } from 'quasar';

const store = useDocumentsStore();
const $q = useQuasar();
const filter = ref('');
const typeFilter = ref('all');
const selected = ref([]);

const columns = [
  { name: 'favorite', label: '', field: 'favorite', align: 'center', sortable: true },
  { name: 'title', label: 'Titolo', field: 'title', align: 'left', sortable: true },
  { name: 'type', label: 'Tipo', field: 'type', align: 'left', sortable: true },
  { name: 'student', label: 'Studente', field: 'student', align: 'left' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center', sortable: true },
  { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'actions', label: 'Azioni', align: 'right' }
];

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
    // Should call store action
    row.favorite = !row.favorite;
    $q.notify({ 
        message: row.favorite ? 'Aggiunto ai preferiti' : 'Rimosso dai preferiti', 
        icon: 'star', 
        color: row.favorite ? 'amber' : 'grey-8' 
    });
};

const batchApprove = () => {
    $q.dialog({
        title: 'Approvazione Multipla',
        message: `Approvare ${selected.value.length} documenti?`,
        cancel: true
    }).onOk(() => {
        // Mock store call
        selected.value.forEach(d => d.status = 'Approved');
        selected.value = [];
        $q.notify({ type: 'positive', message: 'Documenti approvati' });
    })
};

const batchArchive = () => {
     $q.notify({ message: 'Documenti archiviati', icon: 'archive' });
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
