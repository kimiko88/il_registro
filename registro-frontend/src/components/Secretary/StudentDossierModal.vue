<template>
  <q-dialog v-model="show" persistent full-width>
    <q-card style="max-width: 1000px; margin: auto">
      <q-card-section class="bg-primary text-white row items-center justify-between">
        <div class="row items-center">
          <q-icon name="folder_shared" size="28px" class="q-mr-sm" />
          <div class="text-h6">Fascicolo Personale Studente (Storico Multi-Anno)</div>
        </div>
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section class="q-pa-md">
        <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify">
          <q-tab name="panoramica" icon="analytics" label="Storico Multi-Anno" />
          <q-tab name="allegati" icon="attach_file" label="Documentazione & Allegati" />
          <q-tab name="bes" icon="medical_services" label="Fascicolo BES/DSA & PDP" />
        </q-tabs>

        <q-separator class="q-mb-md" />

        <q-tab-panels v-model="tab" animated>
          <!-- Panoramica Tab -->
          <q-tab-panel name="panoramica">
            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-4">
                <q-card flat bordered class="bg-grey-1">
                  <q-card-section>
                    <div class="text-subtitle2 text-weight-bold text-primary">Riepilogo Presenze Multi-Anno</div>
                    <div class="text-h4 text-weight-bolder text-secondary q-my-xs">97.8%</div>
                    <div class="text-caption text-grey-7">Assenze Totali: 12 ore (tutte giustificate)</div>
                  </q-card-section>
                </q-card>
              </div>

              <div class="col-12 col-md-4">
                <q-card flat bordered class="bg-grey-1">
                  <q-card-section>
                    <div class="text-subtitle2 text-weight-bold text-primary">Media Generale Voti</div>
                    <div class="text-h4 text-weight-bolder text-positive q-my-xs">8.3</div>
                    <div class="text-caption text-grey-7">Scrutinio Finale 1° e 2° anno: Promosso</div>
                  </q-card-section>
                </q-card>
              </div>

              <div class="col-12 col-md-4">
                <q-card flat bordered class="bg-grey-1">
                  <q-card-section>
                    <div class="text-subtitle2 text-weight-bold text-primary">Note Disciplinari</div>
                    <div class="text-h4 text-weight-bolder text-grey-8 q-my-xs">0</div>
                    <div class="text-caption text-grey-7">Nessun provvedimento disciplinare</div>
                  </q-card-section>
                </q-card>
              </div>
            </div>
          </q-tab-panel>

          <!-- Allegati Tab -->
          <q-tab-panel name="allegati">
            <div class="row items-center justify-between q-mb-sm">
              <div class="text-subtitle1 text-weight-bold">Allegati e Documenti Inseriti</div>
              <q-btn color="primary" icon="upload" label="Carica Documento" size="sm" />
            </div>

            <q-list bordered separator>
              <q-item v-for="doc in sampleDocs" :key="doc.id">
                <q-item-section avatar>
                  <q-icon name="picture_as_pdf" color="negative" size="32px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ doc.title }}</q-item-label>
                  <q-item-label caption>Categoria: {{ doc.category }} | Anno: {{ doc.academic_year }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-btn flat icon="download" color="primary" label="Download" />
                </q-item-section>
              </q-item>
            </q-list>
          </q-tab-panel>

          <!-- BES/PDP Tab -->
          <q-tab-panel name="bes">
            <div class="q-pa-sm">
              <div class="text-subtitle1 text-weight-bold text-warning">
                <q-icon name="warning" class="q-mr-xs" /> Documentazione BES / DSA Riservata
              </div>
              <div class="text-body2 text-grey-8 q-mt-xs">
                Certificato Medico Diagnostico DSA (Legge 170/2010) depositato in data 10/09/2025.
                Piano Didattico Personalizzato (PDP) approvato dal Consiglio di Classe.
              </div>
            </div>
          </q-tab-panel>
        </q-tab-panels>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Chiudi" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const show = ref(false)
const tab = ref('panoramica')

const sampleDocs = ref([
  { id: '1', title: 'Documentazione di Ingresso e Certificato di Nascita.pdf', category: 'iscrizione', academic_year: '2024/2025' },
  { id: '2', title: 'Nulla Osta da Scuola Provenienza.pdf', category: 'nulla_osta', academic_year: '2024/2025' },
  { id: '3', title: 'Pagella Storica Anno Precedente 1A.pdf', category: 'pagella_storica', academic_year: '2024/2025' }
])

defineExpose({
  open: () => { show.value = true }
})
</script>
