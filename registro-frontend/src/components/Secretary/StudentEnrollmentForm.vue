<template>
  <q-card style="min-width: 800px; max-width: 90vw;">
    <q-card-section>
      <div class="text-h6">{{ t('studentsPage.newEnrollment') || 'Nuova Iscrizione Studente' }}</div>
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
          :title="t('studentsPage.personalData') || 'Dati Anagrafici'"
          icon="person"
          :done="step > 1"
        >
          <div class="row q-col-gutter-md">
             <div class="col-12 col-md-6">
                <q-input v-model="form.firstName" :label="t('common.name') || 'Nome'" outlined dense :rules="[val => !!val || (t('common.requiredField') || 'Obbligatorio')]" />
             </div>
             <div class="col-12 col-md-6">
                <q-input v-model="form.lastName" :label="t('common.surname') || 'Cognome'" outlined dense :rules="[val => !!val || (t('common.requiredField') || 'Obbligatorio')]" />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.dob" :label="t('studentsPage.birthDate') || 'Data di Nascita'" outlined dense type="date" />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.pob" :label="t('studentsPage.birthPlace') || 'Luogo di Nascita'" outlined dense />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.cf" :label="t('studentsPage.taxCode') || 'Codice Fiscale'" outlined dense maxlength="16" />
             </div>
             <div class="col-12 col-md-6">
                <q-select v-model="form.gender" :options="['M', 'F', 'X']" :label="t('studentsPage.gender') || 'Sesso'" outlined dense />
             </div>
          </div>
        </q-step>

        <q-step
          :name="2"
          :title="t('studentsPage.residenceAndContacts') || 'Residenza & Contatti'"
          icon="home"
          :done="step > 2"
        >
          <div class="row q-col-gutter-md">
             <div class="col-12">
                <q-input v-model="form.address" :label="t('common.address') || 'Indirizzo'" outlined dense />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.city" :label="t('common.city') || 'Città'" outlined dense />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.zip" :label="t('common.zip') || 'CAP'" outlined dense maxlength="5" />
             </div>
             <div class="col-12 col-md-4">
                <q-input v-model="form.phone" :label="t('login.phone') || 'Telefono'" outlined dense />
             </div>
             <div class="col-12">
                <q-input v-model="form.email" :label="t('studentsPage.personalEmail') || 'Email Personale'" outlined dense />
             </div>
          </div>
        </q-step>

        <q-step
          :name="3"
          :title="t('studentsPage.parentsGuardians') || 'Genitori / Tutori'"
          icon="family_restroom"
          :done="step > 3"
        >
          <div class="text-subtitle2 q-mb-sm">{{ t('studentsPage.primaryContact') || 'Contatto Principale (Padre/Madre/Tutore)' }}</div>
          <div class="row q-col-gutter-md">
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent1Name" :label="t('common.fullName') || 'Nome e Cognome'" outlined dense />
             </div>
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent1Email" :label="t('login.emailLabel') || 'Email'" outlined dense />
             </div>
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent1Phone" :label="t('login.phone') || 'Telefono'" outlined dense />
             </div>
          </div>
          
          <q-separator class="q-my-md" />
          
          <div class="text-subtitle2 q-mb-sm">{{ t('studentsPage.secondaryContact') || 'Secondo Contatto (Opzionale)' }}</div>
          <div class="row q-col-gutter-md">
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent2Name" :label="t('common.fullName') || 'Nome e Cognome'" outlined dense />
             </div>
             <div class="col-12 col-md-6">
                <q-input v-model="form.parent2Email" :label="t('login.emailLabel') || 'Email'" outlined dense />
             </div>
          </div>
        </q-step>

        <q-step
          :name="4"
          :title="t('studentsPage.schoolData') || 'Dati Scolastici'"
          icon="school"
        >
           <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                 <q-select v-model="form.class" :label="t('studentsPage.assignedClass') || 'Classe Assegnata'" :options="classOptions" outlined dense />
              </div>
              <div class="col-12 col-md-6">
                 <q-input v-model="form.enrollmentDate" :label="t('studentsPage.enrollmentDate') || 'Data Iscrizione'" type="date" outlined dense />
              </div>
              <div class="col-12">
                 <q-toggle v-model="form.documentsSubmitted" :label="t('studentsPage.documentsSubmitted') || 'Documenti consegnati (Carta identità, Foto)'" />
              </div>
           </div>
        </q-step>

        <template v-slot:navigation>
          <q-stepper-navigation>
            <q-btn @click="nextStep" color="primary" :label="step === 4 ? (t('studentsPage.completeEnrollment') || 'Completa Iscrizione') : (t('common.next') || 'Avanti')" />
            <q-btn v-if="step > 1" flat color="primary" @click="$refs.stepper.previous()" :label="t('common.back') || 'Indietro'" class="q-ml-sm" />
          </q-stepper-navigation>
        </template>
      </q-stepper>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const { t } = useI18n()
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
            $q.notify({ type: 'positive', message: t('studentsPage.enrollmentSuccess') || 'Iscrizione completata con successo' })
            emit('complete', { ...form })
        }, 1000)
    } else {
        stepper.value.next()
    }
}

defineExpose({
    form,
    step,
    nextStep,
    stepper
})
</script>
