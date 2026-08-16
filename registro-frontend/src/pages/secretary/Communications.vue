<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-xl">
      <div>
        <h1 class="text-h4 text-weight-bold text-outfit q-my-none text-gradient-premium">
          Comunicazioni & Circolari
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-sm q-mb-none">Gestione flussi informativi per docenti e famiglie</p>
      </div>
      <q-btn color="primary" unelevated icon="add" label="Nuova Circolare" class="rounded-lg q-px-md shadow-sm" no-caps @click="showCreator = true" />
    </div>

    <div class="row q-col-gutter-xl">
        <!-- Stats/Filters -->
        <div class="col-12 col-md-3">
            <q-card flat class="rounded-xl border-slate-100 bg-white overflow-hidden shadow-soft">
              <q-list padding class="text-slate-600">
                  <q-item-label header class="text-weight-bold text-uppercase text-xs letter-spacing-1 text-slate-400">Filtra per</q-item-label>
                  <q-item clickable v-ripple active-class="bg-indigo-50 text-indigo-700 text-weight-bold" :active="filter === 'all'" @click="filter = 'all'" class="q-mx-sm rounded-lg">
                      <q-item-section avatar><q-icon name="inbox" /></q-item-section>
                      <q-item-section>Tutte le circolari</q-item-section>
                  </q-item>
                  <q-item clickable v-ripple active-class="bg-indigo-50 text-indigo-700 text-weight-bold" :active="filter === 'teachers'" @click="filter = 'teachers'" class="q-mx-sm rounded-lg">
                      <q-item-section avatar><q-icon name="school" /></q-item-section>
                      <q-item-section>Solo Docenti</q-item-section>
                  </q-item>
                  <q-item clickable v-ripple active-class="bg-indigo-50 text-indigo-700 text-weight-bold" :active="filter === 'families'" @click="filter = 'families'" class="q-mx-sm rounded-lg">
                       <q-item-section avatar><q-icon name="family_restroom" /></q-item-section>
                       <q-item-section>Solo Famiglie</q-item-section>
                  </q-item>
              </q-list>
            </q-card>
        </div>

        <!-- List -->
        <div class="col-12 col-md-9">
            <q-card flat class="rounded-xl border-slate-100 bg-white shadow-soft overflow-hidden">
                <q-table
                    :rows="filteredCirculars"
                    :columns="columns"
                    row-key="id"
                    :filter="search"
                    flat
                    class="bg-transparent"
                >
                    <template v-slot:top-right>
                         <q-input dense debounce="300" v-model="search" placeholder="Cerca comunicazione..." outlined class="bg-white min-width-250">
                            <template v-slot:append><q-icon name="search" color="slate-400" /></template>
                         </q-input>
                    </template>
                    
                    <template v-slot:body-cell-recipients="props">
                        <q-td :props="props">
                            <div class="row q-gutter-xs">
                              <q-chip v-if="props.row.recipients.teachers" size="sm" icon="school" label="Docenti" class="bg-indigo-50 text-indigo-700" />
                              <q-chip v-if="props.row.recipients.parents" size="sm" icon="people" label="Genitori" class="bg-orange-50 text-orange-700" />
                              <q-chip v-if="props.row.recipients.students" size="sm" icon="face" label="Studenti" class="bg-emerald-50 text-emerald-700" />
                            </div>
                            <div v-if="props.row.specificClasses.length > 0" class="text-caption text-slate-400 q-mt-xs">
                              Classi: {{ props.row.specificClasses.join(', ') }}
                            </div>
                        </q-td>
                    </template>

                     <template v-slot:body-cell-actions="props">
                        <q-td :props="props" auto-width>
                            <div class="row q-gutter-sm justify-end">
                              <q-btn flat round icon="visibility" color="slate-400" size="sm">
                                <q-tooltip>Visualizza</q-tooltip>
                              </q-btn>
                              <q-btn flat round icon="delete" color="negative" size="sm" @click="deleteCircular(props.row.id)">
                                <q-tooltip>Elimina</q-tooltip>
                              </q-btn>
                            </div>
                        </q-td>
                    </template>
                    
                    <template v-slot:no-data>
                      <div class="full-width q-pa-xl text-center text-slate-400">
                        <q-icon name="mail_outline" size="64px" class="opacity-20 q-mb-md" />
                        <div class="text-h6">Nessuna comunicazione trovata</div>
                      </div>
                    </template>
                </q-table>
            </q-card>
        </div>
    </div>

    <q-dialog v-model="showCreator" persistent class="premium-dialog">
        <CircularCreator @sent="onSent" @cancel="showCreator = false" />
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCommunicationsStore } from '@/stores/communications'
import CircularCreator from '@/components/Secretary/CircularCreator.vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const { t } = useI18n()
const commStore = useCommunicationsStore()
const showCreator = ref(false)
const filter = ref('all')
const search = ref('')

const circulars = computed(() => commStore.communications.map(c => ({
    id: c.id,
    title: c.title,
    date: new Date(c.created_at).toLocaleDateString('it-IT', { day: '2-digit', month: 'short', year: 'numeric' }),
    recipients: c.recipients || { teachers: false, parents: false, students: false },
    specificClasses: c.specific_classes || [],
    content: c.content
})))

const columns = [
    { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true, style: 'width: 120px' },
    { name: 'title', label: 'Oggetto', field: 'title', align: 'left', sortable: true, classes: 'text-weight-bold text-slate-800' },
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
        title: 'Conferma Eliminazione',
        message: 'Sei sicuro di voler eliminare questa circolare? L\'azione è irreversibile.',
        cancel: true,
        persistent: true,
        ok: {
          color: 'negative',
          label: 'Elimina',
          flat: false
        }
    }).onOk(async () => {
        try {
            await commStore.deleteCommunication(id)
            $q.notify({ type: 'positive', message: 'Circolare eliminata con successo' })
            await fetchData()
        } catch (err) {
            $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione' })
        }
    })
}
</script>

<style scoped>
.min-width-250 { min-width: 250px; }
.letter-spacing-1 { letter-spacing: 1px; }
.opacity-20 { opacity: 0.2; }
</style>


