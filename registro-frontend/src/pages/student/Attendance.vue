<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">Presenze e Assenze</div>
       <q-btn icon="fact_check" label="Richiedi Giustificazione" color="primary" @click="showJustifyDialog = true" />
    </div>

    <div class="row q-col-gutter-lg">
        <!-- Stats Sidebar -->
        <div class="col-12 col-md-4">
            <q-card class="text-center q-pa-md">
                <div class="text-h6">Riepilogo</div>
                <div class="q-my-md relative-position flex flex-center">
                    <q-circular-progress
                      show-value
                      font-size="20px"
                      class="q-ma-md"
                      :value="attendancePercentage"
                      size="150px"
                      :thickness="0.2"
                      color="green"
                      track-color="red-1"
                    >
                        {{ attendancePercentage }}%<br><span class="text-caption text-grey">Presente</span>
                    </q-circular-progress>
                </div>
                <div class="row justify-around">
                    <div>
                        <div class="text-h5 text-red">{{ totalAbsences }}</div>
                        <div class="text-caption">Assenze</div>
                    </div>
                    <div>
                        <div class="text-h5 text-orange">{{ totalDelays }}</div>
                        <div class="text-caption">Ritardi</div>
                    </div>
                </div>
            </q-card>

            <q-card class="q-mt-md bg-orange-1">
                <q-card-section>
                    <div class="text-subtitle2"><q-icon name="warning" /> Attenzione</div>
                    <div class="text-caption">Hai raggiunto il 15% di assenze consentite. Mettiti in regola per evitare problemi con l'anno scolastico.</div>
                </q-card-section>
            </q-card>
        </div>

        <!-- Attendance List -->
        <div class="col-12 col-md-8">
            <q-table
              title="Giornale di Classe"
              :rows="attendanceEvents"
              :columns="columns"
              row-key="id"
              :pagination="{ rowsPerPage: 10 }"
            >
                <template v-slot:body-cell-type="props">
                    <q-td :props="props">
                        <q-chip :color="getTypeColor(props.value)" text-color="white" size="sm" dense>{{ props.value }}</q-chip>
                    </q-td>
                </template>
                <template v-slot:body-cell-justified="props">
                    <q-td :props="props">
                        <q-icon v-if="props.value" name="check_circle" color="green" />
                        <q-icon v-else name="cancel" color="red" />
                    </q-td>
                </template>
            </q-table>
        </div>
    </div>

    <!-- Justification Dialog -->
    <q-dialog v-model="showJustifyDialog">
        <q-card style="min-width: 400px">
            <q-card-section class="text-h6">Nuova Giustificazione</q-card-section>
            <q-card-section>
                <q-select v-model="justification.event" :options="unjustifiedAbsences" option-label="label" label="Seleziona Assenza" outlined />
                <q-select v-model="justification.reason" :options="['Motivi di Salute', 'Motivi Familiari', 'Visita Medica']" label="Motivazione" outlined class="q-mt-md" />
                <q-input v-model="justification.notes" type="textarea" label="Note Aggiuntive" outlined class="q-mt-md" />
                
                <div class="q-mt-lg border-dashed q-pa-md text-center bg-grey-1">
                    <q-icon name="touch_app" size="md" color="grey" />
                    <div class="text-caption">Firma Digitale (PIN genitore richiesto)</div>
                </div>
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Invia Richiesta" @click="submitJustification" v-close-popup />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const showJustifyDialog = ref(false)

const attendanceEvents = ref([
    { id: 1, date: '2024-12-10', type: 'Assenza', justified: false, notes: '' },
    { id: 2, date: '2024-11-20', type: 'Ritardo', justified: true, notes: 'Entrata ore 09:00' },
    { id: 3, date: '2024-10-05', type: 'Assenza', justified: true, notes: 'Malattia' },
])

const columns = [
    { name: 'date', label: 'Data', align: 'left', field: 'date', sortable: true },
    { name: 'type', label: 'Tipo', align: 'left', field: 'type' },
    { name: 'justified', label: 'Giustificata', align: 'center', field: 'justified' },
    { name: 'notes', label: 'Note', align: 'left', field: 'notes' }
]

const totalAbsences = computed(() => attendanceEvents.value.filter(e => e.type === 'Assenza').length)
const totalDelays = computed(() => attendanceEvents.value.filter(e => e.type === 'Ritardo').length)
const totalDays = 100 // Mock total school days
const attendancePercentage = computed(() => Math.round(((totalDays - totalAbsences.value) / totalDays) * 100))

const unjustifiedAbsences = computed(() => 
    attendanceEvents.value
        .filter(e => !e.justified && e.type === 'Assenza')
        .map(e => ({ label: `${e.date} - Assenza`, value: e.id }))
)

const justification = ref({ event: null, reason: null, notes: '' })

const getTypeColor = (type) => {
    if (type === 'Assenza') return 'red';
    if (type === 'Ritardo') return 'orange';
    return 'grey';
}

const submitJustification = () => {
    $q.notify({ type: 'positive', message: 'Richiesta inviata in attesa di approvazione' })
    // In real app: call API to create JustificationRequest
}
</script>

<style scoped>
.border-dashed {
    border: 2px dashed #ccc;
    border-radius: 8px;
}
</style>
