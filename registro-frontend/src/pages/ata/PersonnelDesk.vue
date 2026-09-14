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
          {{ t('personnelDesk.workflowPhasesTitle') }}
        </div>
        <div class="row items-center justify-between workflow-steps-grid">
          <div class="workflow-step-pill" :class="{ 'step-active': statusFilter === 'submitted' }">
            <div class="step-num bg-blue-100 text-blue-800">1</div>
            <div>
              <div class="text-weight-bold text-caption text-slate-900">{{ t('personnelDesk.step1Title') }}</div>
              <div class="text-caption text-slate-400">{{ t('personnelDesk.step1Subtitle') }}</div>
            </div>
          </div>
          <q-icon name="arrow_forward" color="grey-5" size="18px" class="gt-xs" />
          <div class="workflow-step-pill" :class="{ 'step-active': statusFilter === 'aa_review' }">
            <div class="step-num bg-cyan-100 text-cyan-800">2</div>
            <div>
              <div class="text-weight-bold text-caption text-slate-900">{{ t('personnelDesk.step2Title') }}</div>
              <div class="text-caption text-slate-400">{{ t('personnelDesk.step2Subtitle') }}</div>
            </div>
          </div>
          <q-icon name="arrow_forward" color="grey-5" size="18px" class="gt-xs" />
          <div class="workflow-step-pill" :class="{ 'step-active': statusFilter === 'dsga_review' }">
            <div class="step-num bg-teal-100 text-teal-800">3</div>
            <div>
              <div class="text-weight-bold text-caption text-slate-900">{{ t('personnelDesk.step3Title') }}</div>
              <div class="text-caption text-slate-400">{{ t('personnelDesk.step3Subtitle') }}</div>
            </div>
          </div>
          <q-icon name="arrow_forward" color="grey-5" size="18px" class="gt-xs" />
          <div class="workflow-step-pill" :class="{ 'step-active': statusFilter === 'ds_review' }">
            <div class="step-num bg-purple-100 text-purple-800">4</div>
            <div>
              <div class="text-weight-bold text-caption text-slate-900">{{ t('personnelDesk.step4Title') }}</div>
              <div class="text-caption text-slate-400">{{ t('personnelDesk.step4Subtitle') }}</div>
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
        :options="filterOptions"
        @update:model-value="loadRequests"
      />
    </div>

    <!-- Requests Cards Grid -->
    <div v-if="loading" class="q-pa-xl text-center">
      <q-spinner-dots color="primary" size="40px" />
    </div>
    <div v-else-if="requests.length === 0" class="q-pa-xl text-center text-slate-400 bg-white rounded-2xl shadow-sm border border-slate-200">
      <q-icon name="inbox" size="48px" class="q-mb-sm" />
      <div class="text-body1 text-weight-medium">{{ t('personnelDesk.noRequests') }}</div>
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
              {{ req.applicant_name || t('personnelDesk.applicantDefault') }}
            </div>
            <div class="text-caption text-slate-500 q-mb-sm">
              {{ t('personnelDesk.roleLabel') }} <strong>{{ req.applicant_role ? (t('roles.' + req.applicant_role) || req.applicant_role) : t('personnelDesk.roleDefault') }}</strong>
            </div>

            <!-- Period & Details -->
            <div class="period-box q-pa-sm rounded-borders bg-slate-50 border border-slate-200 q-mb-sm">
              <div class="row items-center text-slate-700 text-caption text-weight-bold">
                <q-icon name="event" size="16px" class="q-mr-xs text-primary" />
                {{ req.start_date }} ➔ {{ req.end_date }}
              </div>
              <div class="text-caption text-slate-600 q-mt-xs">
                {{ t('personnelDesk.durationLabel') }} <strong>{{ req.days ? `${req.days} ${t('personnelDesk.daysUnit')}` : `${req.hours} ${t('personnelDesk.hoursUnit')}` }}</strong>
                <span v-if="req.sub_category"> • {{ req.sub_category }}</span>
              </div>
            </div>

            <!-- Description -->
            <div class="text-body2 text-slate-700 q-mb-sm">
              {{ req.description }}
            </div>

            <!-- Workflow Notes Progress -->
            <div v-if="req.aa_note" class="text-caption text-cyan-900 bg-cyan-50 q-pa-xs rounded-borders q-mb-xs">
              <strong>{{ t('personnelDesk.aaNoteLabel') }}</strong> {{ req.aa_note }}
            </div>
            <div v-if="req.dsga_note" class="text-caption text-teal-900 bg-teal-50 q-pa-xs rounded-borders q-mb-xs">
              <strong>{{ t('personnelDesk.dsgaNoteLabel') }}</strong> {{ req.dsga_note }}
            </div>
            <div v-if="req.ds_decree_num" class="text-caption text-purple-900 bg-purple-50 q-pa-xs rounded-borders">
              <strong>{{ t('personnelDesk.dsDecreeLabel') }}</strong> {{ req.ds_decree_num }} — {{ req.ds_note }}
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
                :label="t('personnelDesk.sendDraftBtn')"
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
                :label="t('personnelDesk.aaReviewBtn')"
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
                :label="t('personnelDesk.dsgaSignBtn')"
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
                :label="t('personnelDesk.dsApproveBtn')"
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
            {{ t('personnelDesk.dialogNewTitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.categoryLabel') }}</label>
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
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.subCategoryLabel') }}</label>
            <q-input v-model="createForm.sub_category" outlined dense :placeholder="t('personnelDesk.subCategoryPlaceholder')" />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.startDateLabel') }}</label>
              <q-input v-model="createForm.start_date" type="date" outlined dense />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.endDateLabel') }}</label>
              <q-input v-model="createForm.end_date" type="date" outlined dense />
            </div>
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.daysLabel') }}</label>
              <q-input v-model.number="createForm.days" type="number" step="0.5" outlined dense />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.hoursLabel') }}</label>
              <q-input v-model.number="createForm.hours" type="number" step="0.5" outlined dense />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.descriptionLabel') }}</label>
            <q-input v-model="createForm.description" outlined dense type="textarea" rows="3" :placeholder="t('personnelDesk.descriptionPlaceholder')" />
          </div>

          <q-checkbox v-model="createForm.submit_now" :label="t('personnelDesk.submitNowLabel')" />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('personnelDesk.cancelBtn')" color="grey-7" v-close-popup />
          <q-btn
            color="indigo-8"
            :label="t('personnelDesk.saveRequestBtn')"
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
            {{ t('personnelDesk.dialogAATitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div class="text-body2 text-slate-700">
            {{ t('personnelDesk.applicant') }} <strong>{{ activeReq?.applicant_name }}</strong><br />
            {{ t('personnelDesk.type') }} <strong>{{ getCategoryLabel(activeReq?.category) }}</strong> ({{ activeReq?.start_date }} ➔ {{ activeReq?.end_date }})
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.aaNotesLabel') }}</label>
            <q-input v-model="aaForm.note" outlined dense type="textarea" rows="3" :placeholder="t('personnelDesk.aaNotesPlaceholder')" />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.aaOutcomeLabel') }}</label>
            <q-btn-toggle
              v-model="aaForm.approve"
              toggle-color="primary"
              dense
              rounded
              :options="[
                { label: t('personnelDesk.aaFavorable'), value: true },
                { label: t('personnelDesk.aaUnfavorable'), value: false }
              ]"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('personnelDesk.cancelBtn')" color="grey-7" v-close-popup />
          <q-btn
            color="cyan-8"
            :label="t('personnelDesk.confirmAABtn')"
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
            {{ t('personnelDesk.dialogDSGATitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div class="text-body2 text-slate-700">
            {{ t('personnelDesk.applicant') }} <strong>{{ activeReq?.applicant_name }}</strong><br />
            {{ t('personnelDesk.type') }} <strong>{{ getCategoryLabel(activeReq?.category) }}</strong>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.dsgaNotesLabel') }}</label>
            <q-input v-model="dsgaForm.note" outlined dense type="textarea" rows="3" :placeholder="t('personnelDesk.dsgaNotesPlaceholder')" />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.dsgaOutcomeLabel') }}</label>
            <q-btn-toggle
              v-model="dsgaForm.approve"
              toggle-color="teal-8"
              dense
              rounded
              :options="[
                { label: t('personnelDesk.dsgaFavorable'), value: true },
                { label: t('personnelDesk.dsgaUnfavorable'), value: false }
              ]"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('personnelDesk.cancelBtn')" color="grey-7" v-close-popup />
          <q-btn
            color="teal-8"
            :label="t('personnelDesk.confirmDSGABtn')"
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
            {{ t('personnelDesk.dialogDSTitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div class="text-body2 text-slate-700">
            {{ t('personnelDesk.applicant') }} <strong>{{ activeReq?.applicant_name }}</strong><br />
            {{ t('personnelDesk.type') }} <strong>{{ getCategoryLabel(activeReq?.category) }}</strong>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.dsDecreeNumLabel') }}</label>
            <q-input v-model="dsForm.decree_num" outlined dense :placeholder="t('personnelDesk.dsDecreePlaceholder')" />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.dsNotesLabel') }}</label>
            <q-input v-model="dsForm.note" outlined dense type="textarea" rows="3" :placeholder="t('personnelDesk.dsNotesPlaceholder')" />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('personnelDesk.dsOutcomeLabel') }}</label>
            <q-btn-toggle
              v-model="dsForm.approve"
              toggle-color="purple-8"
              dense
              rounded
              :options="[
                { label: t('personnelDesk.dsApprove'), value: true },
                { label: t('personnelDesk.dsReject'), value: false }
              ]"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('personnelDesk.cancelBtn')" color="grey-7" v-close-popup />
          <q-btn
            color="purple-8"
            :label="t('personnelDesk.confirmDSBtn')"
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
  return t('roles.' + r) || r
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

const filterOptions = computed(() => [
  { label: t('personnelDesk.filterAll'), value: '' },
  { label: t('personnelDesk.filterSubmitted'), value: 'submitted' },
  { label: t('personnelDesk.filterDSGAReview'), value: 'dsga_review' },
  { label: t('personnelDesk.filterDSReview'), value: 'ds_review' },
  { label: t('personnelDesk.filterApproved'), value: 'approved' },
  { label: t('personnelDesk.filterRejected'), value: 'rejected' }
])

const categoryOptions = computed(() => [
  { label: t('personnelDesk.catFerie'), value: 'ferie' },
  { label: t('personnelDesk.catPermessoBreve'), value: 'permesso_breve' },
  { label: t('personnelDesk.catMalattia'), value: 'malattia' },
  { label: t('personnelDesk.catCongedoParentale'), value: 'congedo_parentale' },
  { label: t('personnelDesk.catPermessoStudio'), value: 'permesso_studio' },
  { label: t('personnelDesk.catAspettativa'), value: 'aspettativa' },
  { label: t('personnelDesk.catAltro'), value: 'altro' }
])

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
    $q.notify({ type: 'negative', message: t('personnelDesk.notifyLoadError') })
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
    $q.notify({ type: 'warning', message: t('personnelDesk.notifyDescRequired') })
    return
  }
  saving.value = true
  try {
    await personnelDeskService.createRequest(createForm.value)
    $q.notify({ type: 'positive', message: t('personnelDesk.notifySavedSuccess') })
    createDialog.value = false
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('personnelDesk.notifySaveError') })
  } finally {
    saving.value = false
  }
}

async function submitDraft(req) {
  try {
    await personnelDeskService.submitRequest(req.id)
    $q.notify({ type: 'positive', message: t('personnelDesk.notifySubmittedSuccess') })
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('personnelDesk.notifySubmitError') })
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
    $q.notify({ type: 'positive', message: t('personnelDesk.notifyAACompleted') })
    aaDialog.value = false
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('personnelDesk.notifyAAError') })
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
    $q.notify({ type: 'positive', message: t('personnelDesk.notifyDSGASuccess') })
    dsgaDialog.value = false
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('personnelDesk.notifyDSGAError') })
  } finally {
    saving.value = false
  }
}

function openDSApproveDialog(req) {
  activeReq.value = req
  dsForm.value = {
    decree_num: `DECR-${new Date().getFullYear()}/${Math.floor(Math.random() * 900 + 100)}`,
    note: t('personnelDesk.dsNotesPlaceholder'),
    approve: true
  }
  dsDialog.value = true
}

async function submitDSApprove() {
  saving.value = true
  try {
    await personnelDeskService.dsApprove(activeReq.value.id, dsForm.value)
    $q.notify({ type: 'positive', message: t('personnelDesk.notifyDSSuccess') })
    dsDialog.value = false
    await loadRequests()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('personnelDesk.notifyDSError') })
  } finally {
    saving.value = false
  }
}

function getCategoryLabel(c) {
  switch (c) {
    case 'ferie': return t('personnelDesk.catFerie')
    case 'permesso_breve': return t('personnelDesk.catPermessoBreve')
    case 'malattia': return t('personnelDesk.catMalattia')
    case 'congedo_parentale': return t('personnelDesk.catCongedoParentale')
    case 'permesso_studio': return t('personnelDesk.catPermessoStudio')
    case 'aspettativa': return t('personnelDesk.catAspettativa')
    case 'altro': return t('personnelDesk.catAltro')
    default: return c
  }
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
    case 'draft': return t('personnelDesk.statusDraft')
    case 'submitted': return t('personnelDesk.statusSubmitted')
    case 'aa_review': return t('personnelDesk.statusAAReview')
    case 'dsga_review': return t('personnelDesk.statusDSGAReview')
    case 'ds_review': return t('personnelDesk.statusDSReview')
    case 'approved': return t('personnelDesk.statusApproved')
    case 'rejected': return t('personnelDesk.statusRejected')
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
