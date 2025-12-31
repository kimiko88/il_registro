<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">Le Mie Presenze</div>
       <q-btn label="Giustifica Assenze" color="primary" icon="assignment_turned_in" @click="showModal = true" />
    </div>

    <!-- Counters -->
    <div class="row q-col-gutter-md q-mb-md">
         <div class="col-6 col-sm-3">
             <q-card class="text-center q-pa-sm">
                 <div class="text-caption text-grey">Assenze</div>
                 <div class="text-h4 text-red">{{ stats.absences }}</div>
             </q-card>
         </div>
         <div class="col-6 col-sm-3">
             <q-card class="text-center q-pa-sm">
                 <div class="text-caption text-grey">Ritardi</div>
                 <div class="text-h4 text-orange">{{ stats.delays }}</div>
             </q-card>
         </div>
         <div class="col-6 col-sm-3">
             <q-card class="text-center q-pa-sm">
                 <div class="text-caption text-grey">Uscite Anticipate</div>
                 <div class="text-h4 text-blue">{{ stats.early }}</div>
             </q-card>
         </div>
         <div class="col-6 col-sm-3">
             <q-card class="text-center q-pa-sm">
                 <div class="text-caption text-grey">Totale Ore Perse</div>
                 <div class="text-h4">{{ stats.hoursLost }}</div>
             </q-card>
         </div>
    </div>

    <!-- Calendar / List Toggle -->
    <q-card>
        <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify">
            <q-tab name="list" icon="list" label="Elenco Eventi" />
            <q-tab name="calendar" icon="event" label="Calendario" />
        </q-tabs>

        <q-separator />

        <q-tab-panels v-model="tab" animated>
            <q-tab-panel name="list">
                <q-table
                    :rows="events"
                    :columns="columns"
                    row-key="id"
                    flat
                >
                    <template v-slot:body-cell-type="props">
                        <q-td :props="props">
                            <q-badge :color="getTypeColor(props.value)" :label="props.value" />
                        </q-td>
                    </template>
                    <template v-slot:body-cell-justified="props">
                        <q-td :props="props">
                             <q-icon v-if="props.value" name="check_circle" color="green" />
                             <q-icon v-else name="cancel" color="red" />
                        </q-td>
                    </template>
                </q-table>
            </q-tab-panel>

            <q-tab-panel name="calendar">
                <div class="text-center text-grey q-pa-xl">
                    <q-icon name="event_note" size="64px" />
                    <div>Vista Calendario (Placeholder)</div>
                </div>
            </q-tab-panel>
        </q-tab-panels>
    </q-card>

    <!-- Justification Request Dialog -->
    <q-dialog v-model="showModal">
        <q-card style="min-width: 400px">
            <q-card-section class="text-h6">Richiedi Giustificazione</q-card-section>
            
            <q-card-section>
                <q-select 
                    v-model="justification.event" 
                    :options="unjustifiedEvents" 
                    option-label="label"
                    label="Seleziona Assenza/Ritardo" 
                    outlined
                    class="q-mb-md"
                />
                
                <q-select
                    v-model="justification.reason"
                    :options="['Motivi di Salute', 'Motivi Familiari', 'Visita Medica', 'Traffico/Trasporti']"
                    label="Motivazione"
                    outlined
                    class="q-mb-md"
                />

                <q-input v-model="justification.note" type="textarea" label="Note Aggiuntive" outlined />
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Invia Richiesta" @click="submitJustification" />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tab = ref('list')
const showModal = ref(false)

const events = ref([
    { id: 1, date: '2025-01-10', type: 'Assenza', justified: false, hours: 5 },
    { id: 2, date: '2025-01-12', type: 'Ritardo', justified: true, hours: 1 },
    { id: 3, date: '2025-01-15', type: 'Assenza', justified: false, hours: 5 },
])

const stats = computed(() => ({
    absences: events.value.filter(e => e.type === 'Assenza').length,
    delays: events.value.filter(e => e.type === 'Ritardo').length,
    early: events.value.filter(e => e.type === 'Uscita Anticipata').length,
    hoursLost: events.value.reduce((acc, curr) => acc + curr.hours, 0)
}))

const columns = [
    { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
    { name: 'type', label: 'Evento', field: 'type', align: 'left' },
    { name: 'hours', label: 'Ore', field: 'hours', align: 'center' },
    { name: 'justified', label: 'Giustificata', field: 'justified', align: 'center' }
]

const unjustifiedEvents = computed(() => 
    events.value.filter(e => !e.justified).map(e => ({
        label: `${e.date} - ${e.type}`,
        value: e.id,
        ...e
    }))
)

// Form
const justification = ref({ event: null, reason: 'Motivi di Salute', note: '' })

const getTypeColor = (type) => {
    return type === 'Assenza' ? 'red' : (type === 'Ritardo' ? 'orange' : 'blue')
}

const submitJustification = () => {
    if(!justification.value.event) return
    
    $q.loading.show()
    setTimeout(() => {
        $q.loading.hide()
        showModal.value = false
        $q.notify({ type: 'positive', message: 'Richiesta inviata al coordinatore.' })
        // Optimistic update
        const evt = events.value.find(e => e.id === justification.value.event.value)
        if(evt) evt.justified = true // Simulating immediate 'pending' state visual
        
        justification.value = { event: null, reason: 'Motivi di Salute', note: '' }
    }, 1000)
}
</script>
