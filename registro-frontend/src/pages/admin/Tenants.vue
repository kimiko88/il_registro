<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none">{{ t('tenantsPage.title') || 'Gestione Tenant & Scuole' }}</h1>
        <p class="text-subtitle1 text-grey-7 q-mb-none">{{ t('tenantsPage.subtitle') || 'Pannello Superadmin per il provisioning di scuole e quote di risorse' }}</p>
      </div>
      <q-btn color="primary" icon="add" :label="t('tenantsPage.newTenant') || 'Nuovo Tenant'" @click="showDialog = true" unelevated />
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
              {{ props.value === 'active' ? (t('tenantsPage.active') || 'ATTIVO') : (t('tenantsPage.inactive') || 'DISATTIVO') }}
            </q-chip>
          </q-td>
        </template>
        <template v-slot:body-cell-quota="props">
          <q-td :props="props">
            <div class="text-caption">
              <div>{{ t('tenantsPage.maxStudents') || 'Studenti max' }}: <strong>{{ props.row.quota.max_students }}</strong></div>
              <div>{{ t('tenantsPage.maxTeachers') || 'Docenti max' }}: <strong>{{ props.row.quota.max_teachers }}</strong></div>
              <div>{{ t('tenantsPage.maxStorage') || 'Storage' }}: <strong>{{ props.row.quota.max_storage_mb }} MB</strong></div>
            </div>
          </q-td>
        </template>
        <template v-slot:no-data>
          <div class="full-width row flex-center text-grey q-gutter-sm q-py-lg">
            <q-icon size="2em" name="business" />
            <span>Nessun tenant configurato</span>
          </div>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Nuovo Tenant -->
    <q-dialog v-model="showDialog">
      <q-card style="width: min(500px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden shadow-24">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ t('tenantsPage.createTenant') || 'Crea Nuovo Tenant' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup aria-label="Chiudi" />
        </q-card-section>

        <q-card-section>
          <q-form @submit="createTenant" class="q-gutter-md">
            <q-input v-model="form.name" :label="t('tenantsPage.schoolName') || 'Nome Scuola / Istituto'" outlined dense rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model="form.code" :label="t('tenantsPage.tenantCode') || 'Codice Meccanografico / Tenant'" outlined dense rules="[val => !!val || 'Campo obbligatorio']" />
            <q-input v-model.number="form.max_students" type="number" :label="t('tenantsPage.maxStudents') || 'Max Studenti'" outlined dense />
            <q-input v-model.number="form.max_teachers" type="number" :label="t('tenantsPage.maxTeachers') || 'Max Docenti'" outlined dense />
            <q-input v-model.number="form.max_storage_mb" type="number" :label="t('tenantsPage.maxStorage') || 'Max Storage (MB)'" outlined dense />

            <div class="row justify-end q-mt-md">
              <q-btn :label="t('tenantsPage.cancel') || 'Annulla'" flat v-close-popup />
              <q-btn :label="t('tenantsPage.submit') || 'Crea Tenant'" color="primary" type="submit" unelevated :loading="submitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import tenantsService from '@/services/tenantsService'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const { t } = useI18n()
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

const columns = computed(() => [
  { name: 'name', label: t('tenantsPage.colName') || 'Nome Scuola', field: 'name', align: 'left', sortable: true },
  { name: 'code', label: t('tenantsPage.colCode') || 'Codice Tenant', field: 'code', align: 'left', sortable: true },
  { name: 'status', label: t('tenantsPage.colStatus') || 'Stato', field: 'status', align: 'center' },
  { name: 'quota', label: t('tenantsPage.colQuota') || 'Quote Assegnate', field: 'quota', align: 'left' },
  { name: 'created_at', label: t('tenantsPage.colCreatedAt') || 'Data Creazione', field: r => new Date(r.created_at).toLocaleDateString(), align: 'right' }
])

const loadTenants = async () => {
  loading.value = true
  try {
    const res = await tenantsService.getTenants()
    tenants.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: t('tenantsPage.errorLoading') || 'Errore caricamento tenant' })
  } finally {
    loading.value = false
  }
}

const createTenant = async () => {
  submitting.value = true
  try {
    await tenantsService.createTenant(form.value)
    $q.notify({ type: 'positive', message: t('tenantsPage.successCreated') || 'Tenant creato con successo!' })
    showDialog.value = false
    loadTenants()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('tenantsPage.errorCreating') || 'Errore creazione tenant' })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadTenants()
})
</script>
