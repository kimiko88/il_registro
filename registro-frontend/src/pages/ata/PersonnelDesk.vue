<template>
  <q-page class="q-pa-md q-pa-lg-xl personnel-desk-page">
    <!-- Hero Header -->
    <div class="row items-center justify-between q-mb-lg gap-md">
      <div class="col-12 col-md-7">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge color="indigo-8" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="assignment" size="14px" class="q-mr-xs" />
            {{ t('personnelDesk.badge') || 'SPORTELLO DIGITALE • WORKFLOW 4 LIVELLI' }}
          </q-badge>
          <q-badge outline color="primary" class="q-px-sm q-py-xs text-weight-bold">
            {{ userRoleLabel }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none flex items-center">
          <q-icon name="forward_to_inbox" color="indigo-8" class="q-mr-sm" size="36px" />
          {{ t('personnelDesk.title') || 'Sportello Digitale Personale' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          {{ t('personnelDesk.subtitle') || 'Workflow autorizzativo istanze dipendenti: Dipendente ➔ Istruttoria AA ➔ Visto DSGA ➔ Decreto DS' }}
        </div>
      </div>

      <!-- Controls & Actions -->
      <div class="col-12 col-md-5 row items-center justify-end q-gutter-sm">
        <q-btn
          color="indigo-8"
          icon="add_circle"
          :label="t('personnelDesk.newRequestBtn') || 'Nuova Domanda'"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openNewRequestDialog"
        />
        <q-btn
          flat
          round
          dense
          color="primary"
          icon="refresh"
          :loading="loading"
          @click="loadRequests"
        >
          <q-tooltip>{{ t('common.refresh') || 'Aggiorna' }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Workflow Legend Banner -->
    <q-card class="rounded-2xl shadow-sm border border-indigo-100 bg-white q-mb-lg">
      <q-card-section class="q-pa-md">
        <div class="text-caption text-slate-500 text-weight-bold text-uppercase q-mb-sm">
          Fasi del Workflow Istruttorio
        </div>
        <div class="row items-center justify-between workflow-steps-grid">
          <div class="workflow-step-pill" :class="{ 'step-active': statusFilter === 'submitted' }">
            <div class="step-num bg-blue-100 text-blue-800">1</div>
            <div>
              <div class="text-weight-bold text-caption text-slate-900">Istanza Inviata</div>
              <div class="text-caption text-slate-400">Dipendente</div>
            </div>
          </div>
          <q-icon name="arrow_forward" color="grey-5" size="18px" class="gt-xs" />
          <div class="workflow-step-pill" :class="{ 'step-active': statusFilter === 'aa_review' }">
            <div class="step-num bg-cyan-100 text-cyan-800">2</div>
            <div>
              <div class="text-weight-bold text-caption text-slate-900">Istruttoria AA</div>
              <div class="text-caption text-slate-400">Ass. Amministrativo</div>
            </div>
          </div>
          <q-icon name="arrow_forward" color="grey-5" size="18px" class="gt-xs" />
          <div class="workflow-step-pill" :class="{ 'step-active': statusFilter === 'dsga_review' }">
            <div class="step-num bg-teal-100 text-teal-800">3</div>
            <div>
              <div class="text-weight-bold text-caption text-slate-900">Visto di Regolarità</div>
              <div class="text-caption text-slate-400">DSGA</div>
            </div>
          </div>
          <q-icon name="arrow_forward" color="grey-5" size="18px" class="gt-xs" />
          <div class="workflow-step-pill" :class="{ 'step-active': statusFilter === 'ds_review' }">
            <div class="step-num bg-purple-100 text-purple-800">4</div>
            <div>
              <div class="text-weight-bold text-caption text-slate-900">Decreto / Provvedimento</div>
              <div class="text-caption text-slate-400">Dirigente Scolastico</div>
            </div>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Filters Bar -->
    <div class="row items-center justify-between q-mb-md">
      <q-btn-toggle
        v-model="statusFilter"
        toggle-color="primary"
        no-caps
        dense
        rounded
        class="bg-white shadow-1"
        :options="[
          { label: 'Tutte', value: '' },
          { label: 'Da Istruire (AA)', value: 'submitted' },
          { label: 'Attesa Visto (DSGA)', value: 'dsga_review' },
          { label: 'Attesa Decreto (DS)', value: 'ds_review' },
          { label: 'Approvate', value: 'approved' },
          { label: 'Respinte', value: 'rejected' }
        ]"
        @update:model-value="loadRequests"
      />
    </div>

    <!-- Requests Cards Grid -->
    <div v-if="loading" class="q-pa-xl text-center">
      <q-spinner-dots color="primary" size="40px" />
    </div>
    <div v-else-if="requests.length === 0" class="q-pa-xl text-center text-slate-400 bg-white rounded-2xl shadow-sm border border-slate-200">
      <q-icon name="inbox" size="48px" class="q-mb-sm" />
      <div class="text-body1 text-weight-medium">Nessuna richiesta presente per questo filtro</div>
    </div>
    <div v-else class="row q-col-gutter-lg">
      <div v-for="req in requests" :key="req.id" class="col-12 col-md-6 col-lg-4">
        <q-card class="request-card rounded-2xl shadow-sm border border-slate-200 bg-white full-height column justify-between">
          <q-card-section>
            <!-- Card Header -->
            <div class="row items-center justify-between q-mb-sm">
              <q-badge :color="getCategoryColor(req.category)" class="text-weight-bold q-px-sm">
                {{ getCategoryLabel(req.category) }}
              </q-badge>
              <q-badge :color="getStatusBadgeColor(req.status)">
                {{ getStatusBadgeLabel(req.status) }}
              </q-badge>
            </div>

            <!-- Applicant Info -->
            <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs">
              {{ req.applicant_name || 'Dipendente' }}
            </div>
            <div class="text-caption text-slate-500 q-mb-sm">
              Ruolo: <strong>{{ req.applicant_role || 'Docente/ATA' }}</strong>
            </div>

            <!-- Period & Details -->
            <div class="period-box q-pa-sm rounded-borders bg-slate-50 border border-slate-200 q-mb-sm">
              <div class="row items-center text-slate-700 text-caption text-weight-bold">
                <q-icon name="event" size="16px" class="q-mr-xs text-primary" />
                {{ req.start_date }} ➔ {{ req.end_date }}
              </div>
              <div class="text-caption text-slate-600 q-mt-xs">
                Durata: <strong>{{ req.days ? `${req.days} giorni` : `${req.hours} ore` }}</strong>
                <span v-if="req.sub_category"> • {{ req.sub_category }}</span>
              </div>
            </div>

            <!-- Description -->
            <div class="text-body2 text-slate-700 q-mb-sm">
              {{ req.description }}
            </div>

            <!-- Workflow Notes Progress -->
            <div v-if="req.aa_note" class="text-caption text-cyan-900 bg-cyan-50 q-pa-xs rounded-borders q-mb-xs">
              <strong>Istruttoria AA:</strong> {{ req.aa_note }}
            </div>
            <div v-if="req.dsga_note" class="text-caption text-teal-900 bg-teal-50 q-pa-xs rounded-borders q-mb-xs">
              <strong>Visto DSGA:</strong> {{ req.dsga_note }}
            </div>
            <div v-if="req.ds_decree_num" class="text-caption text-purple-900 bg-purple-50 q-pa-xs rounded-borders">
              <strong>Decreto DS:</strong> {{ req.ds_decree_num }} — {{ req.ds_note }}
            </div>
          </q-card-section>

          <!-- Action Buttons / Workflow Triggers -->
          <q-card-section class="border-t border-slate-100 row items-center justify-between q-pt-sm">
            <div class="text-caption text-slate-400">
              {{ formatDate(req.created_at) }}
            </div>

            <div class="row items-center q-gutter-xs">
              <!-- Dipendente: invia bozza -->
              <q-btn
                v-if="req.status === 'draft' && req.applicant_id === user?.id"
                color="primary"
                dense
                rounded
                no-caps
                label="Invia Pratica"
                size="sm"
                class="q-px-sm"
                @click="submitDraft(req)"
              />

              <!-- AA: Istruttoria -->
              <q-btn
                v-if="(req.status === 'submitted' || req.status === 'aa_review') && canAAReview"
                color="cyan-8"
                dense
                rounded
                no-caps
                icon="fact_check"
                label="Istruisci AA"
                size="sm"
                class="q-px-sm text-weight-bold"
                @click="openAAReviewDialog(req)"
              />

              <!-- DSGA: Visto -->
              <q-btn
                v-if="req.status === 'dsga_review' && canDSGASign"
                color="teal-8"
                dense
                rounded
                no-caps
                icon="draw"
                label="Apponi Visto"
                size="sm"
                class="q-px-sm text-weight-bold"
                @click="openDSGASignDialog(req)"
              />

              <!-- DS: Provvedimento / Decreto -->
              <q-btn
                v-if="req.status === 'ds_review' && canDSApprove"
                color="purple-8"
                dense
                rounded
                no-caps
                icon="gavel"
                label="Emana Decreto"
                size="sm"
                class="q-px-sm text-weight-bold"
                @click="openDSApproveDialog(req)"
              />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Dialog: Nuova Domanda -->
    <q-dialog v-model="createDialog" persistent>
      <q-card style="min-width: 520px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="post_add" color="indigo-8" class="q-mr-sm" size="24px" />
            Nuova Istanza / Domanda Personale
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Tipologia Istanza *</label>
            <q-select
              v-model="createForm.category"
              :options="categoryOptions"
              emit-value
              map-options
              outlined
              dense
            />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Sotto-categoria / Norma (opzionale)</label>
            <q-input v-model="createForm.sub_category" outlined dense placeholder="es. Legge 104/92 art. 33, Permesso sindacale" />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Data Inizio *</label>
              <q-input v-model="createForm.start_date" type="date" outlined dense />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Data Fine *</label>
              <q-input v-model="createForm.end_date" type="date" outlined dense />
            </div>
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Numero Giorni</label>
              <q-input v-model.number="createForm.days" type="number" step="0.5" outlined dense />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Ore (se a ore)</label>
              <q-input v-model.number="createForm.hours" type="number" step="0.5" outlined dense />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Motivazione / Relazione *</label>
            <q-input v-model="createForm.description" outlined dense type="textarea" rows="3" placeholder="Specificare i dettagli dell'istanza..." />
          </div>

          <q-checkbox v-model="createForm.submit_now" label="Invia immediatamente per l'istruttoria (altrimenti salva in bozza)" />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat label="Annulla" color="grey-7" v-close-popup />
          <q-btn
            color="indigo-8"
            label="Salva Istanza"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="saving"
            @click="submitCreate"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Istruttoria AA -->
    <q-dialog v-model="aaDialog" persistent>
      <q-card style="min-width: 460px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="fact_check" color="cyan-8" class="q-mr-sm" size="24px" />
            Istruttoria Pratica — Assistente Amministrativo
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div class="text-body2 text-slate-700">
            Richiedente: <strong>{{ activeReq?.applicant_name }}</strong><br />
            Tipologia: <strong>{{ activeReq?.category }}</strong> ({{ activeReq?.start_date }} ➔ {{ activeReq?.end_date }})
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Note di Istruttoria (verifica monte ore, documentazione)</label>
            <q-input v-model="aaForm.note" outlined dense type="textarea" rows="3" placeholder="Documentazione verificata con esito positivo..." />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Esito Istruttoria</label>
            <q-btn-toggle
              v-model="aaForm.approve"
              toggle-color="primary"
              dense
              rounded
              :options="[
                { label: 'Favorevole (Inoltra a DSGA)', value: true },
                { label: 'Non Conforme (Respingi)', value: false }
              ]"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat label="Annulla" color="grey-7" v-close-popup />
          <q-btn
            color="cyan-8"
            label="Conferma Istruttoria"
            no-caps
            rounded
            :loading="saving"
            @click="submitAAReview"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Visto DSGA -->
    <q-dialog v-model="dsgaDialog" persistent>
      <q-card style="min-width: 460px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="draw" color="teal-8" class="q-mr-sm" size="24px" />
            Visto di Regolarità Contabile / Organizzativa — DSGA
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div class="text-body2 text-slate-700">
            Richiedente: <strong>{{ activeReq?.applicant_name }}</strong><br />
            Tipologia: <strong>{{ activeReq?.category }}</strong>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Parere / Visto DSGA</label>
            <q-input v-model="dsgaForm.note" outlined dense type="textarea" rows="3" placeholder="Visto favorevole di regolarità contabile ed organizzativa..." />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Decisione Visto</label>
            <q-btn-toggle
              v-model="dsgaForm.approve"
              toggle-color="teal-8"
              dense
              rounded
              :options="[
                { label: 'Visto Favorevole (Inoltra a DS)', value: true },
                { label: 'Parere Contrario (Respingi)', value: false }
              ]"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat label="Annulla" color="grey-7" v-close-popup />
          <q-btn
            color="teal-8"
            label="Apponi Visto"
            no-caps
            rounded
            :loading="saving"
            @click="submitDSGASign"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Decreto DS -->
    <q-dialog v-model="dsDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="gavel" color="purple-8" class="q-mr-sm" size="24px" />
            Provvedimento Finale — Dirigente Scolastico
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div class="text-body2 text-slate-700">
            Richiedente: <strong>{{ activeReq?.applicant_name }}</strong><br />
            Tipologia: <strong>{{ activeReq?.category }}</strong>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Numero Decreto / Provvedimento (se approvato)</label>
            <q-input v-model="dsForm.decree_num" outlined dense placeholder="es. DECR-2026/089" />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Note / Motivazione Provvedimento</label>
            <q-input v-model="dsForm.note" outlined dense type="textarea" rows="3" placeholder="Si decreta l'approvazione dell'istanza..." />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Esito Finale</label>
            <q-btn-toggle
              v-model="dsForm.approve"
              toggle-color="purple-8"
              dense
              rounded
              :options="[
                { label: 'Emana Decreto (Approva)', value: true },
                { label: 'Rigetta Istanza', value: false }
              ]"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat label="Annulla" color="grey-7" v-close-popup />
          <q-btn
            color="purple-8"
            label="Conferma Provvedimento"
            no-caps
            rounded
            :loading="saving"
            @click="submitDSApprove"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import personnelDeskService from '@/services/personnelDeskService'

const $q = useQuasar()
const { t, locale } = useI18n()
const authStore = useAuthStore()
const { user, userRole } = storeToRefs(authStore)

const requests = ref([])
const loading = ref(false)
const saving = ref(false)
const statusFilter = ref('')

const activeReq = ref(null)
const createDialog = ref(false)
const aaDialog = ref(false)
const dsgaDialog = ref(false)
const dsDialog = ref(false)

const userRoleLabel = computed(() => {
  const r = (userRole.value || '').toLowerCase()
  return r === 'dsga' ? 'DSGA' : r === 'assistente_amministrativo' ? 'Assistente Amministrativo' : r === 'collaboratore_ds' ? 'Collaboratore DS' : r === 'collaboratore_scolastico' ? 'Collaboratore Scolastico' : r
})

const canAAReview = computed(() => {
  const r = (userRole.value || '').toLowerCase()
  return r === 'assistente_amministrativo' || r === 'secretary' || r === 'admin' || r === 'superadmin'
})

const canDSGASign = computed(() => {
  const r = (userRole.value || '').toLowerCase()
  return r === 'dsga' || r === 'admin' || r === 'superadmin'
})

const canDSApprove = computed(() => {
  const r = (userRole.value || '').toLowerCase()
  return r === 'principal' || r === 'vice_principal' || r === 'admin' || r === 'superadmin'
})

const categoryOptions = [
  { label: 'Ferie Ordinarie', value: 'ferie' },
  { label: 'Permesso Breve (ore)', value: 'permesso_breve' },
  { label: 'Malattia / Visite Specialistiche', value: 'malattia' },
  { label: 'Congedo Parentale / Maternità', value: 'congedo_parentale' },
  { label: 'Diritto allo Studio (150 ore)', value: 'permesso_studio' },
  { label: 'Aspettativa Retribuita / Non Retribuita', value: 'aspettativa' },
  { label: 'Altro Permesso CCNL', value: 'altro' }
]

const createForm = ref({
  category: 'ferie',
  sub_category: '',
  start_date: new Date().toISOString().substring(0, 10),
  end_date: new Date().toISOString().substring(0, 10),
  days: 1,
  hours: 0,
  description: '',
  submit_now: true
})

const aaForm = ref({ note: '', approve: true })
const dsgaForm = ref({ note: '', approve: true })
const dsForm = ref({ decree_num: '', note: '', approve: true })

async function loadRequests() {
  loading.value = true
  try {
    const res = await personnelDeskService.listRequests(statusFilter.value)
    requests.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore caricamento istanze' })
  } finally {
    loading.value = false
  }
}

function openNewRequestDialog() {
  createForm.value = {
    category: 'ferie',
    sub_category: '',
    start_date: new Date().toISOString().substring(0, 10),
    end_date: new Date().toISOString().substring(0, 10),
    days: 1,
    hours: 0,
    description: '',
    submit_now: true
  }
  createDialog.value = true
}

async function submitCreate() {
  if (!createForm.value.description) {
    $q.notify({ type: 'warning', message: 'Descrizione / Motivazione obbligatoria' })
    return
  }
  saving.value = true
  try {
    await personnelDeskService.createRequest(createForm.value)
    $q.notify({ type: 'positive', message: 'Istanza registrata con successo!' })
    createDialog.value = false
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore salvataggio istanza' })
  } finally {
    saving.value = false
  }
}

async function submitDraft(req) {
  try {
    await personnelDeskService.submitRequest(req.id)
    $q.notify({ type: 'positive', message: 'Pratica inoltrata per l\'istruttoria' })
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore invio pratica' })
  }
}

function openAAReviewDialog(req) {
  activeReq.value = req
  aaForm.value = { note: '', approve: true }
  aaDialog.value = true
}

async function submitAAReview() {
  saving.value = true
  try {
    await personnelDeskService.aaReview(activeReq.value.id, aaForm.value)
    $q.notify({ type: 'positive', message: 'Istruttoria AA completata' })
    aaDialog.value = false
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore istruttoria' })
  } finally {
    saving.value = false
  }
}

function openDSGASignDialog(req) {
  activeReq.value = req
  dsgaForm.value = { note: '', approve: true }
  dsgaDialog.value = true
}

async function submitDSGASign() {
  saving.value = true
  try {
    await personnelDeskService.dsgaSign(activeReq.value.id, dsgaForm.value)
    $q.notify({ type: 'positive', message: 'Visto DSGA registrato' })
    dsgaDialog.value = false
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore visto DSGA' })
  } finally {
    saving.value = false
  }
}

function openDSApproveDialog(req) {
  activeReq.value = req
  dsForm.value = {
    decree_num: `DECR-${new Date().getFullYear()}/${Math.floor(Math.random() * 900 + 100)}`,
    note: 'Si approva l\'istanza conformemente all\'istruttoria',
    approve: true
  }
  dsDialog.value = true
}

async function submitDSApprove() {
  saving.value = true
  try {
    await personnelDeskService.dsApprove(activeReq.value.id, dsForm.value)
    $q.notify({ type: 'positive', message: 'Provvedimento emanato con successo!' })
    dsDialog.value = false
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore provvedimento' })
  } finally {
    saving.value = false
  }
}

function getCategoryLabel(c) {
  const opt = categoryOptions.find(o => o.value === c)
  return opt ? opt.label : c
}

function getCategoryColor(c) {
  switch (c) {
    case 'ferie': return 'teal-8'
    case 'permesso_breve': return 'indigo-8'
    case 'malattia': return 'negative'
    case 'congedo_parentale': return 'purple-8'
    case 'permesso_studio': return 'cyan-8'
    case 'aspettativa': return 'amber-9'
    default: return 'blue-grey-7'
  }
}

function getStatusBadgeLabel(st) {
  switch (st) {
    case 'draft': return 'Bozza'
    case 'submitted': return '1. Inviata'
    case 'aa_review': return '2. Istruttoria AA'
    case 'dsga_review': return '3. Visto DSGA'
    case 'ds_review': return '4. Decreto DS'
    case 'approved': return 'Approvata'
    case 'rejected': return 'Respinta'
    default: return st
  }
}

function getStatusBadgeColor(st) {
  switch (st) {
    case 'draft': return 'grey-6'
    case 'submitted': return 'blue-7'
    case 'aa_review': return 'cyan-8'
    case 'dsga_review': return 'teal-8'
    case 'ds_review': return 'purple-8'
    case 'approved': return 'positive'
    case 'rejected': return 'negative'
    default: return 'grey-7'
  }
}

function formatDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleDateString(locale.value || 'it-IT', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  loadRequests()
})
</script>

<style scoped>
.personnel-desk-page {
  background: #f8fafc;
  min-height: 100vh;
}

.workflow-steps-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.workflow-step-pill {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  border-radius: 12px;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
}

.step-active {
  background: #eff6ff;
  border-color: #3b82f6;
}

.step-num {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 12px;
}

.request-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.request-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px -5px rgba(0,0,0,0.08);
}
</style>
