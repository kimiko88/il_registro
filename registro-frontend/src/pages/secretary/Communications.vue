<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4 text-weight-bold">Comunicazioni & Circolari</div>
       <q-btn color="primary" icon="add" label="Nuova Circolare" @click="showCreator = true" />
    </div>

    <div class="row q-col-gutter-md">
        <!-- Stats/Filters -->
        <div class="col-12 col-md-3">
            <q-list bordered class="bg-white rounded-borders">
                <q-item clickable v-ripple active-class="bg-blue-1 text-primary" :active="filter === 'all'" @click="filter = 'all'">
                    <q-item-section avatar><q-icon name="inbox" /></q-item-section>
                    <q-item-section>Tutte</q-item-section>
                </q-item>
                <q-item clickable v-ripple :active="filter === 'teachers'" @click="filter = 'teachers'">
                    <q-item-section avatar><q-icon name="school" /></q-item-section>
                    <q-item-section>Docenti</q-item-section>
                </q-item>
                <q-item clickable v-ripple :active="filter === 'families'" @click="filter = 'families'">
                     <q-item-section avatar><q-icon name="family_restroom" /></q-item-section>
                     <q-item-section>Famiglie</q-item-section>
                </q-item>
            </q-list>
        </div>

        <!-- List -->
        <div class="col-12 col-md-9">
            <q-card>
                <q-table
                    :rows="filteredCirculars"
                    :columns="columns"
                    row-key="id"
                    :filter="search"
                >
                    <template v-slot:top-right>
                         <q-input dense debounce="300" v-model="search" placeholder="Cerca...">
                            <template v-slot:append><q-icon name="search" /></template>
                         </q-input>
                    </template>
                    
                    <template v-slot:body-cell-recipients="props">
                        <q-td :props="props">
                            <q-chip v-if="props.row.recipients.teachers" size="sm" icon="school" label="Docenti" />
                            <q-chip v-if="props.row.recipients.parents" size="sm" icon="people" label="Genitori" />
                            <q-chip v-if="props.row.recipients.students" size="sm" icon="face" label="Studenti" />
                            <span v-if="props.row.specificClasses.length > 0" class="text-caption text-grey">({{ props.row.specificClasses.join(', ') }})</span>
                        </q-td>
                    </template>

                     <template v-slot:body-cell-actions="props">
                        <q-td :props="props" auto-width>
                            <q-btn flat round icon="visibility" color="grey-7" />
                            <q-btn flat round icon="delete" color="negative" @click="deleteCircular(props.row.id)" />
                        </q-td>
                    </template>
                </q-table>
            </q-card>
        </div>
    </div>

    <q-dialog v-model="showCreator" persistent>
        <CircularCreator @sent="onSent" @cancel="showCreator = false" />
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCommunicationsStore } from 'src/stores/communications'
import CircularCreator from 'src/components/Secretary/CircularCreator.vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const commStore = useCommunicationsStore()
const showCreator = ref(false)
const filter = ref('all')
const search = ref('')

const circulars = computed(() => commStore.communications.map(c => ({
    id: c.id,
    title: c.title,
    date: new Date(c.created_at).toLocaleDateString('it-IT'),
    recipients: c.recipients || { teachers: false, parents: false, students: false },
    specificClasses: c.specific_classes || [],
    content: c.content
})))

const columns = [
    { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true, style: 'width: 100px' },
    { name: 'title', label: 'Oggetto', field: 'title', align: 'left', sortable: true },
    { name: 'recipients', label: 'Destinatari', field: 'recipients', align: 'left' },
    { name: 'actions', label: '', align: 'right' }
]

onMounted(async () => {
    await fetchData()
})

const fetchData = async () => {
    try {
        await commStore.fetchCommunications()
    } catch (err) {
        $q.notify({ type: 'negative', message: 'Errore durante il caricamento' })
    }
}

const filteredCirculars = computed(() => {
    let res = circulars.value
    if (filter.value === 'teachers') res = res.filter(c => c.recipients.teachers)
    if (filter.value === 'families') res = res.filter(c => c.recipients.parents || c.recipients.students)
    return res
})

const onSent = async () => {
    showCreator.value = false
    await fetchData()
}

const deleteCircular = async (id) => {
    $q.dialog({
        title: 'Conferma',
        message: 'Vuoi eliminare questa circolare?',
        cancel: true,
        persistent: true
    }).onOk(async () => {
        try {
            await commStore.deleteCommunication(id)
            $q.notify({ type: 'positive', message: 'Circolare eliminata' })
            await fetchData()
        } catch (err) {
            $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione' })
        }
    })
}
</script>
