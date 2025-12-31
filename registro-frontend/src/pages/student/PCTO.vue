<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">PCTO - Percorsi Competenze</div>
       <q-btn label="Registra Ore" color="secondary" icon="schedule" @click="showHoursDialog = true" />
    </div>

    <!-- Progress Card -->
    <q-card class="bg-primary text-white q-mb-lg">
        <q-card-section>
            <div class="row items-center">
                <div class="col-12 col-md-3 text-center">
                    <q-knob
                        show-value
                        :min="0"
                        :max="90"
                        :model-value="totalHours"
                        size="90px"
                        :thickness="0.2"
                        color="secondary"
                        track-color="blue-8"
                        class="text-white q-ma-md"
                    >
                        {{ totalHours }}h
                    </q-knob>
                    <div class="text-caption">Su 90h Richieste</div>
                </div>
                <div class="col-12 col-md-9">
                    <div class="text-h6">Stato Complessivo</div>
                    <div>Hai completato il <strong>{{ percentComplete }}%</strong> delle ore previste per il triennio.</div>
                    <q-linear-progress :value="percentComplete / 100" color="secondary" class="q-mt-sm" track-color="blue-8" size="10px" rounded />
                </div>
            </div>
        </q-card-section>
    </q-card>

    <!-- Active Projects -->
    <div class="text-h6 q-mb-sm">Progetti Attivi</div>
    <div class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
            <q-card bordered flat>
                <q-item>
                    <q-item-section avatar>
                        <q-avatar color="orange-1" text-color="orange" icon="business" />
                    </q-item-section>
                    <q-item-section>
                        <q-item-label class="text-weight-bold">Stage presso TechCorp SRL</q-item-label>
                        <q-item-label caption>Tutor Aziendale: Rossi Marco</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <q-badge color="green">In Corso</q-badge>
                    </q-item-section>
                </q-item>
                <q-separator />
                <q-card-section>
                    <div class="row justify-between text-caption text-grey">
                        <span>Periodo: 10 Gen - 20 Feb</span>
                        <span>40/80 Ore Svolte</span>
                    </div>
                    <q-btn flat class="full-width q-mt-sm" label="Vedi Dettagli / Diario" color="primary" />
                </q-card-section>
            </q-card>
        </div>
    </div>

    <!-- Hours Log Dialog -->
    <q-dialog v-model="showHoursDialog">
        <q-card style="min-width: 400px">
            <q-card-section class="text-h6">Registra Attività</q-card-section>
            <q-card-section class="q-gutter-md">
                <q-select outlined v-model="newLog.project" :options="['Stage TechCorp', 'Corso Sicurezza']" label="Progetto" />
                <q-input outlined v-model="newLog.date" type="date" label="Data" />
                <div class="row q-gutter-sm">
                    <q-input outlined v-model="newLog.hours" type="number" label="Ore" class="col" />
                    <q-input outlined v-model="newLog.minutes" type="number" label="Minuti" class="col" />
                </div>
                <q-input outlined v-model="newLog.description" type="textarea" label="Descrizione Attività" />
            </q-card-section>
            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Salva" @click="saveLog" />
            </q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const totalHours = ref(45)
const percentComplete = computed(() => Math.round((totalHours.value / 90) * 100))
const showHoursDialog = ref(false)

const newLog = ref({ project: '', date: '', hours: 0, minutes: 0, description: '' })

const saveLog = () => {
    // API Call
    $q.notify({ type: 'positive', message: 'Ore registrate con successo' })
    showHoursDialog.value = false
}
</script>
