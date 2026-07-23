<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none">Gestione Tenant & Scuole</h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none">Pannello Superadmin per il provisioning di scuole e quote di risorse</p>
      </div>
      <q-btn color="primary" icon="add" label="Nuovo Tenant" @click="showDialog = true" unelevated />
    </div>

    <q-card flat bordered class="q-mb-lg">
      <q-table
        :rows="tenants"
        :columns="columns"
        row-key="id"
        :loading="loading"
        flat
      >
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-chip
              :color="props.value === 'active' ? 'positive' : 'negative'"
              text-color="white"
              size="sm"
            >
              {{ props.value.toUpperCase() }}
            </q-chip>
          </q-td>
        </template>
        <template v-slot:body-cell-quota="props">
          <q-td :props="props">
            <div class="text-caption">
              <div>Studenti max: <strong>{{ props.row.quota.max_students }}</strong></div>
              <div>Docenti max: <strong>{{ props.row.quota.max_teachers }}</strong></div>
              <div>Storage: <strong>{{ props.row.quota.max_storage_mb }} MB</strong></div>
            </div>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Nuovo Tenant -->
    <q-dialog v-model="showDialog">
      <q-card style="min-width: 400px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Crea Nuovo Tenant</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section>
          <q-form @submit="createTenant" class="q-gutter-md">
            <q-input v-model="form.name" label="Nome Scuola / Istituto" outlined dense rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model="form.code" label="Codice Meccanografico / Tenant" outlined dense rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model.number="form.max_students" type="number" label="Max Studenti" outlined dense />
            <q-input v-model.number="form.max_teachers" type="number" label="Max Docenti" outlined dense />
            <q-input v-model.number="form.max_storage_mb" type="number" label="Max Storage (MB)" outlined dense />

            <div class="row justify-end q-mt-md">
              <q-btn label="Annulla" flat v-close-popup />
              <q-btn label="Crea Tenant" color="primary" type="submit" unelevated :loading="submitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import tenantsService from '@/services/tenantsService'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tenants = ref([])
const loading = ref(false)
const showDialog = ref(false)
const submitting = ref(false)

const form = ref({
  name: '',
  code: '',
  max_students: 1000,
  max_teachers: 100,
  max_storage_mb: 10240
})

const columns = [
  { name: 'name', label: 'Nome Scuola', field: 'name', align: 'left', sortable: true },
  { name: 'code', label: 'Codice Tenant', field: 'code', align: 'left', sortable: true },
  { name: 'status', label: 'Stato', field: 'status', align: 'center' },
  { name: 'quota', label: 'Quote Assegnate', field: 'quota', align: 'left' },
  { name: 'created_at', label: 'Data Creazione', field: r => new Date(r.created_at).toLocaleDateString(), align: 'right' }
]

const loadTenants = async () => {
  loading.value = true
  try {
    const res = await tenantsService.getTenants()
    tenants.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore caricamento tenant' })
  } finally {
    loading.value = false
  }
}

const createTenant = async () => {
  submitting.value = true
  try {
    await tenantsService.createTenant(form.value)
    $q.notify({ type: 'positive', message: 'Tenant creato con successo!' })
    showDialog.value = false
    loadTenants()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore creazione tenant' })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadTenants()
})
</script>
