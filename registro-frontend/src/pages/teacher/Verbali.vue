<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none">Verbali Consiglio di Classe (Ver.Di 2.0)</h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none">Gestione, stesura e firma digitale con marca temporale ed IP audit</p>
      </div>
      <q-btn color="primary" icon="add" label="Nuovo Verbale" @click="showDialog = true" unelevated />
    </div>

    <q-card flat bordered>
      <q-table
        :rows="verbali"
        :columns="columns"
        row-key="id"
        :loading="loading"
        flat
      >
        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn
              color="secondary"
              label="Firma con IP"
              icon="draw"
              size="sm"
              unelevated
              @click="signVerbale(props.row.id)"
            />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <q-dialog v-model="showDialog">
      <q-card style="min-width: 500px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Crea Nuovo Verbale</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <q-form @submit="createVerbale" class="q-gutter-md">
            <q-input v-model="form.class_id" label="ID Classe" outlined dense :rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model="form.title" label="Titolo Verbale" outlined dense :rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model="form.content" type="textarea" label="Contenuto / OdG" outlined dense rows="5" />

            <div class="row justify-end q-mt-md">
              <q-btn label="Annulla" flat v-close-popup />
              <q-btn label="Salva Verbale" color="primary" type="submit" unelevated :loading="submitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import verbaliService from '@/services/verbaliService'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const verbali = ref([])
const loading = ref(false)
const showDialog = ref(false)
const submitting = ref(false)

const form = ref({
  class_id: '',
  title: '',
  content: ''
})

const columns = [
  { name: 'title', label: 'Titolo', field: 'title', align: 'left', sortable: true },
  { name: 'class_id', label: 'Classe', field: 'class_id', align: 'center' },
  { name: 'signed_by', label: 'Firme Raccolte', field: r => (r.signed_by || []).length, align: 'center' },
  { name: 'created_at', label: 'Data', field: r => new Date(r.created_at).toLocaleDateString(), align: 'right' },
  { name: 'actions', label: 'Azioni', align: 'center' }
]

const loadVerbali = async () => {
  loading.value = true
  try {
    const res = await verbaliService.getVerbali()
    verbali.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore caricamento verbali' })
  } finally {
    loading.value = false
  }
}

const createVerbale = async () => {
  submitting.value = true
  try {
    await verbaliService.createVerbale(form.value)
    $q.notify({ type: 'positive', message: 'Verbale creato con successo' })
    showDialog.value = false
    loadVerbali()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore creazione verbale' })
  } finally {
    submitting.value = false
  }
}

const signVerbale = async (id) => {
  try {
    await verbaliService.signVerbale(id)
    $q.notify({ type: 'positive', message: 'Firma registrata con audit IP!' })
    loadVerbali()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore firma verbale' })
  }
}

onMounted(() => {
  loadVerbali()
})
</script>
