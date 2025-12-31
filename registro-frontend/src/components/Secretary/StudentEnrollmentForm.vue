<template>
  <q-card style="min-width: 800px; max-width: 90vw;">
    <q-card-section>
      <div class="text-h6">Nuova Iscrizione Studente</div>
    </q-card-section>

    <q-card-section class="q-pt-none">
      <q-stepper
        v-model="step"
        ref="stepper"
        color="primary"
        animated
        flat
      >
        <q-step
          :name="1"
          title="Dati Anagrafici"
          icon="person"
          :done="step > 1"
        >
          <div class="row q-col-gutter-md">
             <div class="col-12 col-md-6">
                <q-input v-model="form.firstName" label="Nome" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
             </div>
             <div class="col-12 col-md-6">
                <q-input v-model="form.lastName" label="Cognome" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.dob" label="Data di Nascita" outlined dense type="date" />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.pob" label="Luogo di Nascita" outlined dense />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.cf" label="Codice Fiscale" outlined dense maxlength="16" />
             </div>
             <div class="col-12 col-md-6">
                <q-select v-model="form.gender" :options="['M', 'F', 'X']" label="Sesso" outlined dense />
             </div>
          </div>
        </q-step>

        <q-step
          :name="2"
          title="Residenza & Contatti"
          icon="home"
          :done="step > 2"
        >
          <div class="row q-col-gutter-md">
             <div class="col-12">
                <q-input v-model="form.address" label="Indirizzo" outlined dense />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.city" label="Città" outlined dense />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.zip" label="CAP" outlined dense maxlength="5" />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.phone" label="Telefono" outlined dense />
             </div>
             <div class="col-12">
                <q-input v-model="form.email" label="Email Personale" outlined dense />
             </div>
          </div>
        </q-step>

        <q-step
          :name="3"
          title="Genitori / Tutori"
          icon="family_restroom"
          :done="step > 3"
        >
          <div class="text-subtitle2 q-mb-sm">Contatto Principale (Padre/Madre/Tutore)</div>
          <div class="row q-col-gutter-md">
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent1Name" label="Nome e Cognome" outlined dense />
             </div>
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent1Email" label="Email" outlined dense />
             </div>
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent1Phone" label="Telefono" outlined dense />
             </div>
          </div>
          
          <q-separator class="q-my-md" />
          
          <div class="text-subtitle2 q-mb-sm">Secondo Contatto (Opzionale)</div>
          <div class="row q-col-gutter-md">
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent2Name" label="Nome e Cognome" outlined dense />
             </div>
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent2Email" label="Email" outlined dense />
             </div>
          </div>
        </q-step>

        <q-step
          :name="4"
          title="Dati Scolastici"
          icon="school"
        >
           <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                 <q-select v-model="form.class" label="Classe Assegnata" :options="classOptions" outlined dense />
              </div>
              <div class="col-12 col-md-6">
                 <q-input v-model="form.enrollmentDate" label="Data Iscrizione" type="date" outlined dense />
              </div>
              <div class="col-12">
                 <q-toggle v-model="form.documentsSubmitted" label="Documenti consegnati (Carta identità, Foto)" />
              </div>
           </div>
        </q-step>

        <template v-slot:navigation>
          <q-stepper-navigation>
            <q-btn @click="nextStep" color="primary" :label="step === 4 ? 'Completa Iscrizione' : 'Avanti'" />
            <q-btn v-if="step > 1" flat color="primary" @click="$refs.stepper.previous()" label="Indietro" class="q-ml-sm" />
          </q-stepper-navigation>
        </template>
      </q-stepper>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const step = ref(1)
const stepper = ref(null)

const emit = defineEmits(['complete', 'cancel'])

const form = reactive({
    firstName: '',
    lastName: '',
    dob: '',
    pob: '',
    cf: '',
    gender: '',
    address: '',
    city: '',
    zip: '',
    phone: '',
    email: '',
    parent1Name: '',
    parent1Email: '',
    parent1Phone: '',
    parent2Name: '',
    parent2Email: '',
    class: '',
    enrollmentDate: new Date().toISOString().split('T')[0],
    documentsSubmitted: false
})

const classOptions = ['1A', '1B', '2A', '2B', '3A', '3B', '4A', '5A'] // Mock

const nextStep = () => {
    if (step.value === 4) {
        // Submit
        $q.loading.show()
        setTimeout(() => {
            $q.loading.hide()
            $q.notify({ type: 'positive', message: 'Iscrizione completata con successo' })
            emit('complete', { ...form })
        }, 1000)
    } else {
        stepper.value.next()
    }
}
</script>
