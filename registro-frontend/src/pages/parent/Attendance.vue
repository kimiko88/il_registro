<template>
  <q-page class="q-pa-md bg-grey-1">
    <div v-if="!selectedChild" class="text-center q-pa-xl text-grey">
        <div class="text-h6">Seleziona un figlio dalla Dashboard.</div>
    </div>
    <div v-else>
        <div class="text-h4 q-mb-md">Presenze di {{ selectedChild.name }}</div>

        <!-- Pending Justifications -->
        <q-card class="bg-orange-1 q-mb-md" v-if="pendingJustifications.length > 0">
            <q-card-section>
                <div class="text-h6 text-orange-9">Assenze da Giustificare</div>
                <div class="text-caption text-grey-8">Le seguenti assenze richiedono la tua giustificazione.</div>
            </q-card-section>
            <q-list separator>
                <q-item v-for="evt in pendingJustifications" :key="evt.id" class="bg-white">
                    <q-item-section avatar>
                         <q-icon name="warning" color="orange" />
                    </q-item-section>
                    <q-item-section>
                        <q-item-label class="text-weight-bold">{{ evt.date }} - {{ evt.type }}</q-item-label>
                        <q-item-label caption v-if="evt.studentRequest">Richiesta Studente: "{{ evt.studentRequest }}"</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <q-btn label="Giustifica" color="primary" @click="justify(evt)" />
                    </q-item-section>
                </q-item>
            </q-list>
        </q-card>

        <!-- Attendance History -->
        <q-card>
            <q-card-section class="text-h6">Storico Eventi</q-card-section>
            <q-table
                :rows="events"
                :columns="columns"
                flat
            >
                 <template v-slot:body-cell-justified="props">
                     <q-td :props="props">
                         <q-chip v-if="props.value" color="green" text-color="white" icon="check" size="sm">Giustificata</q-chip>
                         <q-chip v-else color="red" text-color="white" icon="close" size="sm">Non Giustificata</q-chip>
                     </q-td>
                 </template>
            </q-table>
        </q-card>
    </div>

    <!-- Justify Dialog -->
    <q-dialog v-model="showJustifyDialog">
        <q-card style="min-width: 400px">
            <q-card-section class="text-h6">Giustifica Assenza</q-card-section>
            <q-card-section>
                <div class="text-subtitle1">{{ currentEvent?.date }}</div>
                <q-radio v-model="justificationType" val="Salute" label="Motivi di Salute" />
                <q-radio v-model="justificationType" val="Famiglia" label="Motivi Familiari" />
                <q-input v-model="justificationNote" label="Note (Opzionale)" outlined type="textarea" class="q-mt-sm" />
                <q-checkbox v-model="signed" label="Dichiaro di aver preso visione (Firma Digitale)" class="q-mt-md" />
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Conferma" :disable="!signed" @click="confirmJustify" />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useParentStore } from 'src/stores/parent'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const parentStore = useParentStore()
const selectedChild = computed(() => parentStore.selectedChild)

const showJustifyDialog = ref(false)
const currentEvent = ref(null)
const justificationType = ref('Salute')
const justificationNote = ref('')
const signed = ref(false)

// Mock
const events = ref([
    { id: 1, date: '2025-01-20', type: 'Assenza', justified: false, studentRequest: 'Mal di testa' },
    { id: 2, date: '2025-01-05', type: 'Ritardo', justified: true }
])

const pendingJustifications = computed(() => events.value.filter(e => !e.justified))

const columns = [
    { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
    { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
    { name: 'justified', label: 'Stato', field: 'justified', align: 'center' }
]

const justify = (evt) => {
    currentEvent.value = evt
    justificationNote.value = evt.studentRequest || ''
    signed.value = false
    showJustifyDialog.value = true
}

const confirmJustify = () => {
    // API Call
    currentEvent.value.justified = true
    showJustifyDialog.value = false
    $q.notify({ type: 'positive', message: 'Giustificazione inviata' })
}
</script>
