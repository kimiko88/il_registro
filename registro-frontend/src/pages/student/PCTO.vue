<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-lg">
       <div class="text-h4">PCTO - Percorsi Trasversali</div>
       <q-chip color="orange" text-color="white" icon="timer">Totale: 30 / 90 Ore</q-chip>
    </div>

    <!-- Progress Bar -->
    <q-card class="q-mb-lg">
        <q-card-section>
            <div class="text-subtitle1 q-mb-sm">Avanzamento Monte Ore Triennio</div>
            <q-linear-progress size="25px" :value="0.33" color="primary" stripe rounded>
                <div class="absolute-full flex flex-center">
                    <q-badge color="white" text-color="primary" label="33%" />
                </div>
            </q-linear-progress>
        </q-card-section>
    </q-card>

    <div class="row q-col-gutter-lg">
        <!-- Projects List -->
        <div class="col-12 col-md-8">
            <div class="text-h5 q-mb-md">Progetti Attivi</div>
            <q-list bordered class="rounded-borders bg-white">
                <q-expansion-item
                    expand-separator
                    icon="business"
                    label="Stage presso Tech Solutions Srl"
                    caption="Sviluppo Web - 40 Ore"
                    header-class="text-primary"
                    default-opened
                >
                    <q-card>
                        <q-card-section>
                            <div class="row q-col-gutter-md">
                                <div class="col-12 col-md-6">
                                    <div class="text-weight-bold">Tutor Aziendale:</div>
                                    <div>Mario Rossi (m.rossi@tech.it)</div>
                                </div>
                                <div class="col-12 col-md-6">
                                    <div class="text-weight-bold">Periodo:</div>
                                    <div>10 Giu 2024 - 30 Giu 2024</div>
                                </div>
                            </div>
                            
                            <q-separator class="q-my-md" />
                            
                            <div class="text-h6 q-mb-sm">Diario di Bordo</div>
                            <q-timeline color="secondary">
                                <q-timeline-entry title="Inizio Stage" subtitle="10 Giu 2024">
                                    Accoglienza e assegnazione postazione.
                                </q-timeline-entry>
                                <q-timeline-entry title="Primo Task" subtitle="12 Giu 2024">
                                    Configurazione ambiente di sviluppo Vue.js.
                                </q-timeline-entry>
                            </q-timeline>
                        </q-card-section>
                        <q-card-actions align="right">
                             <q-btn flat icon="edit" label="Compila Diario" color="primary" />
                             <q-btn flat icon="cloud_upload" label="Carica Documenti" color="primary" />
                        </q-card-actions>
                    </q-card>
                </q-expansion-item>

                <q-expansion-item
                    expand-separator
                    icon="school"
                    label="Corso Sicurezza sul Lavoro"
                    caption="Formazione Obbligatoria - 8 Ore"
                >
                    <q-card>
                        <q-card-section>
                            Attestato conseguito il 15 Set 2023.
                        </q-card-section>
                         <q-card-actions align="right">
                             <q-btn flat icon="download" label="Scarica Attestato" color="green" />
                        </q-card-actions>
                    </q-card>
                </q-expansion-item>
            </q-list>
        </div>

        <!-- Documentation -->
        <div class="col-12 col-md-4">
             <q-card>
                <q-card-section>
                    <div class="text-h6">Documentazione</div>
                </q-card-section>
                <q-list separator>
                    <q-item clickable v-ripple>
                        <q-item-section avatar><q-icon name="description" color="grey" /></q-item-section>
                        <q-item-section>Convenzione Stage</q-item-section>
                        <q-item-section side><q-btn flat round icon="download" size="sm" /></q-item-section>
                    </q-item>
                    <q-item clickable v-ripple>
                        <q-item-section avatar><q-icon name="description" color="grey" /></q-item-section>
                        <q-item-section>Patto Formativo</q-item-section>
                        <q-item-section side><q-btn flat round icon="download" size="sm" /></q-item-section>
                    </q-item>
                    <q-item clickable v-ripple>
                        <q-item-section avatar><q-icon name="assignment" color="grey" /></q-item-section>
                        <q-item-section>Foglio Presenze</q-item-section>
                        <q-item-section side><q-btn flat round icon="print" size="sm" /></q-item-section>
                    </q-item>
                </q-list>
             </q-card>
        </div>
    </div>
  </q-page>
</template>

<script setup>
</script>
