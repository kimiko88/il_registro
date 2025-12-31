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
                            <q-btn flat round icon="delete" color="negative" />
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
import { ref, computed } from 'vue'
import CircularCreator from 'src/components/Secretary/CircularCreator.vue'

const showCreator = ref(false)
const filter = ref('all')
const search = ref('')

const circulars = ref([
    { id: 1, title: 'Convocazione Collegio Docenti', date: '10/01/2025', recipients: { teachers: true }, specificClasses: [] },
    { id: 2, title: 'Sciopero 20 Gennaio', date: '12/01/2025', recipients: { parents: true, students: true }, specificClasses: [] },
    { id: 3, title: 'Uscita Anticipata 1A', date: '15/01/2025', recipients: { parents: true }, specificClasses: ['1A'] }
])

const columns = [
    { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true, style: 'width: 100px' },
    { name: 'title', label: 'Oggetto', field: 'title', align: 'left', sortable: true },
    { name: 'recipients', label: 'Destinatari', field: 'recipients', align: 'left' },
    { name: 'actions', label: '', align: 'right' }
]

const filteredCirculars = computed(() => {
    let res = circulars.value
    if (filter.value === 'teachers') res = res.filter(c => c.recipients.teachers)
    if (filter.value === 'families') res = res.filter(c => c.recipients.parents || c.recipients.students)
    return res
})

const onSent = (newCircular) => {
    showCreator.value = false
    circulars.value.unshift({
        id: Date.now(),
        title: newCircular.title,
        date: newCircular.date,
        recipients: newCircular.recipients,
        specificClasses: newCircular.specificClasses,
        content: newCircular.content
    })
}
</script>
