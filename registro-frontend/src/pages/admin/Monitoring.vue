<template>
  <q-page class="q-pa-md">
    <div class="text-h4 text-weight-bold q-mb-md">Monitoraggio Sistema</div>
    
    <div class="row q-gutter-md">
        <!-- Placeholder content -->
        <div class="col-12">
            <q-banner class="bg-info text-white rounded-borders">
                <template v-slot:avatar>
                    <q-icon name="info" />
                </template>
                Funzionalità di monitoraggio avanzato in arrivo nella prossima versione.
            </q-banner>
        </div>

        <div class="col-12 col-md-4">
             <q-card>
                <q-card-section>
                    <div class="text-h6">Server Status</div>
                    <div class="text-positive text-weight-bold">ONLINE</div>
                </q-card-section>
            </q-card>
        </div>
    </div>
  </q-page>
</template>
