<template>
  <q-page class="q-pa-md bg-slate-50">
    <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">Il Mio Profilo</div>
    
    <div class="row q-col-gutter-lg">
      <div class="col-12 col-md-6">
        <q-card class="shadow-sm">
          <q-card-section>
            <div class="text-h6">Dati Personali</div>
            <q-separator class="q-my-sm" />
            <q-input v-model="form.name" label="Nome" readonly filled class="q-mb-sm" />
            <q-input v-model="form.email" label="Email" readonly filled class="q-mb-sm" />
            <q-input v-model="form.phone" label="Telefono" outlined />
            <q-btn label="Salva Modifiche" color="primary" class="q-mt-sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-6">
        <q-card class="shadow-sm">
          <q-card-section>
            <div class="text-h6 text-negative">Area Sicurezza</div>
            <q-separator class="q-my-sm" />
            <q-btn outline color="primary" label="Cambia Password" class="full-width q-mb-sm" />
            <q-btn flat color="grey" label="Esporta Dati GDPR" class="full-width" />
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { reactive } from 'vue'

const form = reactive({
  name: 'Giuseppe Verdi',
  email: 'g.verdi@example.com',
  phone: '+39 333 1234567'
})
</script>
