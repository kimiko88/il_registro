<template>
  <q-page padding class="q-pa-lg bg-slate-50 min-h-screen">
    <!-- Header with Back navigation -->
    <div class="row items-center justify-between q-mb-lg">
      <div class="row items-center q-gutter-x-md">
        <q-btn
          flat
          round
          icon="arrow_back"
          color="primary"
          to="/secretary/students"
          :aria-label="t('common.back') || 'Torna agli studenti'"
        >
          <q-tooltip>{{ t('common.back') || 'Torna indietro' }}</q-tooltip>
        </q-btn>
        <div>
          <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
            {{ t('fascicolo.title') || 'Fascicolo Digitale Studente' }}
          </h1>
          <p class="text-subtitle2 text-slate-500 q-mb-none">
            {{ t('fascicolo.subtitle') || 'Profilo anagrafico e carriera scolastica aggregata' }}
          </p>
        </div>
      </div>
    </div>

    <!-- Skeleton Loading State -->
    <div v-if="loading" class="q-gutter-y-md">
      <q-skeleton type="rect" height="140px" class="rounded-xl shadow-sm" />
      <q-skeleton type="rect" height="180px" class="rounded-xl shadow-sm" />
    </div>

    <!-- Empty State -->
    <q-banner v-else-if="!fascicolo" class="bg-amber-1 text-amber-10 rounded-xl q-pa-md border border-amber-300">
      <template v-slot:avatar>
        <q-icon name="info" color="amber-9" />
      </template>
      <div class="text-weight-bold">{{ t('fascicolo.notFoundTitle') || 'Nessun fascicolo trovato' }}</div>
      <div>{{ t('fascicolo.notFoundDesc') || 'Non è stato possibile reperire la scheda per lo studente selezionato.' }}</div>
    </q-banner>

    <!-- Main Content -->
    <div v-else class="q-gutter-y-lg">
      <!-- Dati Anagrafici -->
      <q-card flat bordered class="rounded-xl bg-white shadow-sm overflow-hidden">
        <q-card-section class="bg-indigo-50 text-indigo-9 q-py-sm row items-center">
          <q-icon name="badge" size="22px" class="q-mr-sm" />
          <div class="text-subtitle1 text-weight-bold">
            {{ t('fascicolo.personalInfo') || 'Dati Anagrafici Studente' }}
          </div>
        </q-card-section>

        <q-separator />

        <q-card-section class="q-pa-md">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-sm-4">
              <div class="text-caption text-grey-7">{{ t('common.fullName') || 'Nome Completo' }}</div>
              <div class="text-body1 text-weight-bold text-slate-900">
                {{ fascicolo.student?.first_name }} {{ fascicolo.student?.last_name }}
              </div>
            </div>
            <div class="col-12 col-sm-4">
              <div class="text-caption text-grey-7">Email</div>
              <div class="text-body1 text-slate-900">
                {{ fascicolo.student?.email || 'N/D' }}
              </div>
            </div>
            <div class="col-12 col-sm-4">
              <div class="text-caption text-grey-7">{{ t('fascicolo.fiscalCode') || 'Codice Fiscale' }}</div>
              <div class="text-body1 text-weight-medium text-slate-900 font-mono">
                {{ fascicolo.student?.fiscal_code || 'N/D' }}
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Genitori / Tutori Legali -->
      <q-card flat bordered class="rounded-xl bg-white shadow-sm overflow-hidden">
        <q-card-section class="bg-blue-50 text-blue-9 q-py-sm row items-center">
          <q-icon name="family_restroom" size="22px" class="q-mr-sm" />
          <div class="text-subtitle1 text-weight-bold">
            {{ t('fascicolo.guardians') || 'Genitori / Tutori Legali' }}
          </div>
        </q-card-section>

        <q-separator />

        <q-card-section class="q-pa-md">
          <div v-if="!fascicolo.guardians || fascicolo.guardians.length === 0" class="text-grey-7 q-py-sm">
            <q-icon name="person_off" size="18px" class="q-mr-xs" />
            {{ t('fascicolo.noGuardians') || 'Nessun genitore o tutore legalmente registrato.' }}
          </div>

          <div v-else class="row q-col-gutter-md">
            <div v-for="g in fascicolo.guardians" :key="g.id" class="col-12 col-sm-6">
              <q-card flat bordered class="rounded-lg bg-slate-50 q-pa-sm">
                <div class="row items-center q-gutter-x-sm">
                  <q-avatar size="36px" color="primary" text-color="white" icon="person" />
                  <div>
                    <div class="text-weight-bold text-slate-900">{{ g.first_name }} {{ g.last_name }}</div>
                    <div class="text-caption text-grey-7">{{ g.email }}</div>
                    <div class="text-caption text-grey-6">{{ g.phone ? ('Tel: ' + g.phone) : 'Recapito non specificato' }}</div>
                  </div>
                </div>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { userService } from '@/services/userService'

const route = useRoute()
const $q = useQuasar()
const { t } = useI18n()
const loading = ref(true)
const fascicolo = ref(null)

onMounted(async () => {
  try {
    const studentId = route.params.id
    if (studentId) {
      const res = await userService.getStudentFascicolo(studentId)
      fascicolo.value = res.data
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: t('fascicolo.errorLoading') || 'Errore durante il recupero del fascicolo',
      caption: err?.userMessage || err?.message
    })
  } finally {
    loading.value = false
  }
})
</script>
