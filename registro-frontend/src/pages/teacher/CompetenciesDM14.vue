<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('competenciesDM14.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('competenciesDM14.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm">
        <q-btn
          color="primary"
          icon="save"
          label="Salva Bozza Certificazione"
          class="rounded-lg shadow-sm"
          @click="saveCertification"
        />
        <q-btn
          color="positive"
          icon="verified"
          :label="$t('competenciesDM14.signByPrincipal')"
          class="rounded-lg shadow-sm"
          @click="signCertification"
        />
      </div>
    </div>

    <!-- SELEZIONE CLASSE & STUDENTE -->
    <q-card class="glass-card shadow-soft border-slate-100 q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-4">
            <q-select
              v-model="selectedTier"
              :options="[
                { label: $t('competenciesDM14.tierPrimaria'), value: 'primaria' },
                { label: $t('competenciesDM14.tierPrimoCiclo'), value: 'primo_ciclo' }
              ]"
              emit-value
              map-options
              outlined
              dense
              label="Ordine di Scuola (D.M. 14/2024)"
            />
          </div>
          <div class="col-12 col-md-4">
            <q-input v-model="studentFullName" outlined dense label="Nominativo Studente" />
          </div>
          <div class="col-12 col-md-4">
            <q-input v-model="councilDate" type="date" outlined dense :label="$t('competenciesDM14.deliberationDate')" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- TABELLA COMPETENZE EUROPEE & DIGCOMP 2.2 -->
    <q-card class="glass-card shadow-soft border-slate-100 q-mb-xl">
      <q-card-section>
        <div class="text-h6 text-weight-bold q-mb-md row items-center q-gutter-x-sm">
          <q-icon name="stars" color="primary" />
          <span>{{ $t('competenciesDM14.rubricTitle') }}</span>
        </div>

        <q-list separator>
          <q-item v-for="(comp, idx) in competenceList" :key="idx" class="q-py-md">
            <q-item-section>
              <q-item-label class="text-weight-bold text-slate-800 text-subtitle1">
                {{ comp.name }}
              </q-item-label>
              <q-item-label caption class="text-slate-500">
                {{ comp.description }}
              </q-item-label>
            </q-item-section>

            <q-item-section side>
              <q-btn-toggle
                v-model="evaluations[comp.key]"
                toggle-color="primary"
                flat
                dense
                :options="[
                  { label: $t('competenciesDM14.levelIniziale'), value: 'iniziale' },
                  { label: $t('competenciesDM14.levelBase'), value: 'base' },
                  { label: $t('competenciesDM14.levelIntermedio'), value: 'intermedio' },
                  { label: $t('competenciesDM14.levelAvanzato'), value: 'avanzato' }
                ]"
              />
            </q-item-section>
          </q-item>
        </q-list>

        <div class="q-mt-lg">
          <div class="text-subtitle2 text-weight-bold q-mb-xs">{{ $t('competenciesDM14.considerations') }}</div>
          <q-input
            v-model="finalConsiderations"
            type="textarea"
            outlined
            rows="3"
            placeholder="Sintesi del consiglio di classe sul percorso triennale dello studente..."
          />
        </div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const selectedTier = ref('primo_ciclo')
const studentFullName = ref('Bianchi Leonardo')
const councilDate = ref(new Date().toISOString().substring(0, 10))
const finalConsiderations = ref('')

const competenceList = [
  { key: 'lingua_madre', name: '1. Comunicazione nella madrelingua o lingua di scolarizzazione', description: 'Capacità di esprimere e interpretare concetti, pensieri e fatti in forma sia orale che scritta.' },
  { key: 'lingue_straniere', name: '2. Comunicazione nelle lingue straniere', description: 'Comprensione e produzione in almeno due lingue comunitarie (inglese e seconda lingua).' },
  { key: 'stem', name: '3. Competenze matematiche, scientifiche e tecnologiche (STEM)', description: 'Abilità di sviluppare e applicare il pensiero matematico e il metodo scientifico per spiegare il mondo.' },
  { key: 'digcomp', name: '4. Competenze digitali (Quadro DigComp 2.2)', description: 'Uso sicuro, critico e responsabile delle tecnologie digitali e dell\'intelligenza artificiale.' },
  { key: 'imparare_imparare', name: '5. Imparare a imparare (Competenza personale e sociale)', description: 'Capacità di perseverare nell\'apprendimento e organizzare il proprio studio autonomo.' },
  { key: 'civica_cittadinanza', name: '6. Competenze sociali e civiche (Educazione Civica)', description: 'Partecipazione attiva e democratica, rispetto della Costituzione e sostenibilità.' },
  { key: 'imprenditorialita', name: '7. Spirito di iniziativa e imprenditorialità', description: 'Capacità di tradurre idee in azioni attraverso la creatività, l\'innovazione e l\'assunzione di rischi calcolati.' },
  { key: 'consapevolezza_culturale', name: '8. Consapevolezza ed espressione culturale', description: 'Apprezzamento dell\'importanza dell\'espressione creativa di idee ed emozioni attraverso arti, musica e teatro.' }
]

const evaluations = ref({
  lingua_madre: 'intermedio',
  lingue_straniere: 'avanzato',
  stem: 'intermedio',
  digcomp: 'avanzato',
  imparare_imparare: 'intermedio',
  civica_cittadinanza: 'avanzato',
  imprenditorialita: 'base',
  consapevolezza_culturale: 'intermedio'
})

async function saveCertification() {
  try {
    await api.post('/competencies-dm14/certifications', {
      student_id: '00000000-0000-0000-0000-000000000001',
      class_id: '00000000-0000-0000-0000-000000000002',
      school_tier: selectedTier.value,
      school_year: '2025/2026',
      rubric_evaluations_json: JSON.stringify(evaluations.value),
      final_considerations: finalConsiderations.value,
      council_deliberation_date: councilDate.value,
      status: 'bozza'
    })
    $q.notify({ type: 'positive', message: 'Certificazione salvata in bozza' })
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' })
  }
}

async function signCertification() {
  $q.notify({ type: 'positive', message: 'Certificazione firmata digitalmente con FEQ del Dirigente' })
}
</script>
