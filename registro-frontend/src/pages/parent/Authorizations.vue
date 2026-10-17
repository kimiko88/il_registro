<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">{{ t('dualSign.title') }}</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ t('dualSign.subtitle') }}
        </p>
      </div>
      <div>
        <q-btn flat round icon="refresh" color="primary" @click="fetchAuthorizations">
          <q-tooltip>{{ t('common.refresh') }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Filter Tabs -->
    <q-tabs
      v-model="activeFilter"
      dense
      class="text-grey-7 bg-white rounded-t-lg shadow-sm q-mb-md"
      active-color="primary"
      indicator-color="primary"
      align="left"
    >
      <q-tab name="all" :label="t('dualSign.all')" />
      <q-tab name="to_sign" :label="t('dualSign.toSign')" icon="edit_document" />
      <q-tab name="waiting_other" :label="t('dualSign.waitingOther')" icon="hourglass_empty" />
      <q-tab name="completed" :label="t('dualSign.completed')" icon="task_alt" />
      <q-tab name="rejected" :label="t('dualSign.rejected')" icon="cancel" />
    </q-tabs>

    <!-- Authorizations Cards List -->
    <div v-if="loading" class="text-center q-pa-xl">
      <q-spinner color="primary" size="40px" />
      <div class="text-grey-6 q-mt-md">{{ t('dualSign.loading') }}</div>
    </div>

    <div v-else-if="filteredAuthorizations.length === 0" class="text-center q-pa-xl bg-white rounded-xl shadow-soft border-slate-100">
      <q-icon name="check_circle" size="64px" color="positive" class="opacity-30" />
      <div class="text-h6 text-slate-700 q-mt-md">{{ t('dualSign.noAuth') }}</div>
      <p class="text-slate-400">{{ t('dualSign.noAuthSubtitle') }}</p>
    </div>

    <div v-else class="row q-col-gutter-md">
      <div
        v-for="item in filteredAuthorizations"
        :key="item.id"
        class="col-12 col-md-6"
      >
        <q-card class="glass-card shadow-soft border-slate-100 rounded-xl overflow-hidden h-full flex flex-col justify-between">
          <q-card-section>
            <div class="row items-center justify-between no-wrap q-mb-xs">
              <q-badge :color="getTypeColor(item.document_type)" class="text-caption q-py-xs q-px-sm">
                {{ getTypeLabel(item.document_type) }}
              </q-badge>
              <q-badge :color="getStatusColor(item.status)" class="text-caption q-py-xs q-px-sm text-weight-bold">
                {{ getStatusLabel(item.status) }}
              </q-badge>
            </div>

            <div class="text-h6 text-weight-bold text-slate-800 q-mt-sm">{{ item.title }}</div>

            <div class="row q-gutter-x-md q-mt-sm text-caption text-grey-7">
              <div v-if="item.deadline" class="row items-center">
                <q-icon name="event" size="16px" class="q-mr-xs text-negative" />
                {{ t('dualSign.deadline') }}: {{ formatDate(item.deadline) }}
              </div>
              <div class="row items-center">
                <q-icon name="history" size="16px" class="q-mr-xs text-primary" />
                {{ t('dualSign.requestedOn') }}: {{ formatDate(item.created_at) }}
              </div>
            </div>

            <!-- Dual signatures indicators -->
            <div class="q-mt-md q-pa-sm bg-slate-50 rounded-lg border border-slate-200">
              <div class="text-caption text-weight-bold text-slate-700 q-mb-xs">{{ t('dualSign.parentSignaturesStatus') }}:</div>
              <div class="row q-col-gutter-sm">
                <div class="col-6 row items-center">
                  <q-icon
                    :name="item.parent1_signed_at ? 'check_circle' : 'radio_button_unchecked'"
                    :color="item.parent1_signed_at ? 'positive' : 'grey-5'"
                    size="18px"
                    class="q-mr-xs"
                  />
                  <span class="text-caption">{{ t('dualSign.parent1') }} {{ item.parent1_signed_at ? '(' + t('dualSign.signed') + ')' : '(' + t('dualSign.waiting') + ')' }}</span>
                </div>
                <div class="col-6 row items-center">
                  <q-icon
                    :name="item.parent2_signed_at ? 'check_circle' : 'radio_button_unchecked'"
                    :color="item.parent2_signed_at ? 'positive' : 'grey-5'"
                    size="18px"
                    class="q-mr-xs"
                  />
                  <span class="text-caption">{{ t('dualSign.parent2') }} {{ item.parent2_signed_at ? '(' + t('dualSign.signed') + ')' : '(' + t('dualSign.waiting') + ')' }}</span>
                </div>
              </div>
            </div>

            <div v-if="item.rejection_reason" class="q-mt-sm q-pa-sm bg-red-50 text-negative text-caption rounded-lg">
              <strong>{{ t('dualSign.rejected') }}:</strong> {{ item.rejection_reason }}
            </div>
          </q-card-section>

          <q-card-actions align="right" class="q-pa-md bg-white border-t border-slate-100">
            <template v-if="canCurrentUserSign(item)">
              <q-btn
                flat
                color="negative"
                :label="t('dualSign.rejectBtn')"
                icon="close"
                class="rounded-lg"
                @click="openRejectDialog(item)"
              />
              <q-btn
                color="primary"
                :label="t('dualSign.signBtn')"
                icon="fingerprint"
                class="rounded-lg shadow-sm q-px-md"
                @click="openSignDialog(item)"
              />
            </template>
            <template v-else-if="item.status === 'pending_second'">
              <div class="row items-center text-amber-9 text-caption">
                <q-icon name="schedule" size="18px" class="q-mr-xs" />
                {{ t('dualSign.waitingOther') }}
              </div>
            </template>
            <template v-else-if="item.status === 'completed'">
              <div class="row items-center text-positive text-caption text-weight-bold">
                <q-icon name="verified" size="18px" class="q-mr-xs" />
                {{ t('dualSign.statusCompleted') }}
              </div>
            </template>
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- PIN Signature Dialog -->
    <q-dialog v-model="showSignDialog">
      <q-card style="width: 400px; max-width: 90vw;" class="q-pa-md rounded-xl">
        <q-card-section class="row items-center">
          <q-avatar icon="fingerprint" color="primary" text-color="white" class="q-mr-sm" />
          <div class="text-h6 text-weight-bold">{{ t('dualSign.pinTitle') }}</div>
        </q-card-section>

        <q-card-section>
          <p class="text-caption text-grey-7">
            {{ t('dualSign.pinHelp') }}: <strong>{{ activeItem?.title }}</strong>
          </p>

          <q-input
            v-model="pinInput"
            :label="t('dualSign.pinPlaceholder')"
            type="password"
            maxlength="6"
            outlined
            dense
            autofocus
            :rules="[val => (val && val.length >= 4) || t('common.requiredField')]"
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            :label="t('dualSign.signBtn')"
            :loading="signing"
            :disable="!pinInput || pinInput.length < 4"
            @click="submitSignature"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>


    <!-- Reject Dialog -->
    <q-dialog v-model="showRejectDialog">
      <q-card style="width: 450px; max-width: 90vw;" class="q-pa-md rounded-xl">
        <q-card-section class="row items-center">
          <q-avatar icon="warning" color="negative" text-color="white" class="q-mr-sm" />
          <div class="text-h6 text-weight-bold">{{ t('dualSign.rejectBtn') }}</div>
        </q-card-section>

        <q-card-section>
          <p class="text-caption text-grey-7">
            {{ activeItem?.title }}
          </p>

          <q-input
            v-model="rejectReasonInput"
            :label="t('common.optionalNotes')"
            type="textarea"
            outlined
            rows="3"
            :rules="[val => !!val || t('common.requiredField')]"
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn
            color="negative"
            :label="t('dualSign.rejectBtn')"
            :loading="rejecting"
            :disable="!rejectReasonInput"
            @click="submitRejection"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { parentDualSignService } from '@/services/parentDualSignService'

const $q = useQuasar()
const { t } = useI18n()

const loading = ref(false)
const authorizations = ref([])
const activeFilter = ref('all')

const showSignDialog = ref(false)
const showRejectDialog = ref(false)
const activeItem = ref(null)
const pinInput = ref('')
const rejectReasonInput = ref('')
const signing = ref(false)
const rejecting = ref(false)

onMounted(fetchAuthorizations)

async function fetchAuthorizations() {
  loading.value = true
  try {
    const res = await parentDualSignService.listAuthorizations()
    authorizations.value = res.data || []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const filteredAuthorizations = computed(() => {
  if (activeFilter.value === 'to_sign') {
    return authorizations.value.filter(a => a.status === 'pending_first')
  }
  if (activeFilter.value === 'waiting_other') {
    return authorizations.value.filter(a => a.status === 'pending_second')
  }
  if (activeFilter.value === 'completed') {
    return authorizations.value.filter(a => a.status === 'completed')
  }
  if (activeFilter.value === 'rejected') {
    return authorizations.value.filter(a => a.status === 'rejected')
  }
  return authorizations.value
})

function canCurrentUserSign(item) {
  if (item.status === 'pending_first') return true
  return false
}

function getStatusColor(status) {
  switch (status) {
    case 'completed': return 'positive'
    case 'pending_second': return 'amber-9'
    case 'pending_first': return 'primary'
    case 'rejected': return 'negative'
    default: return 'grey-7'
  }
}

function getStatusLabel(status) {
  switch (status) {
    case 'completed': return t('dualSign.statusCompleted')
    case 'pending_second': return t('dualSign.statusPendingSecond')
    case 'pending_first': return t('dualSign.statusPendingFirst')
    case 'rejected': return t('dualSign.statusRejected')
    default: return status
  }
}

function getTypeColor(type) {
  switch (type) {
    case 'trip_consent': return 'deep-purple'
    case 'pdp_approval': return 'teal'
    case 'religion_change': return 'indigo'
    case 'enrollment': return 'blue'
    default: return 'primary'
  }
}

function getTypeLabel(type) {
  switch (type) {
    case 'trip_consent': return t('common.category')
    case 'pdp_approval': return 'PDP'
    case 'religion_change': return t('religionAlternative.title')
    case 'enrollment': return t('common.name')
    default: return t('common.details')
  }
}

function formatDate(isoStr) {
  if (!isoStr) return ''
  const d = new Date(isoStr)
  return d.toLocaleDateString(undefined, { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function openSignDialog(item) {
  activeItem.value = item
  pinInput.value = ''
  showSignDialog.value = true
}

function openRejectDialog(item) {
  activeItem.value = item
  rejectReasonInput.value = ''
  showRejectDialog.value = true
}

async function submitSignature() {
  if (!activeItem.value || !pinInput.value) return
  signing.value = true
  try {
    await parentDualSignService.signAuthorization(activeItem.value.id, pinInput.value)
    $q.notify({ type: 'positive', message: t('dualSign.signedSuccess') })
    showSignDialog.value = false
    await fetchAuthorizations()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    signing.value = false
  }
}

async function submitRejection() {
  if (!activeItem.value || !rejectReasonInput.value) return
  rejecting.value = true
  try {
    await parentDualSignService.rejectAuthorization(activeItem.value.id, rejectReasonInput.value)
    $q.notify({ type: 'info', message: t('dualSign.rejectedSuccess') })
    showRejectDialog.value = false
    await fetchAuthorizations()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    rejecting.value = false
  }
}

</script>

<style scoped>
.opacity-30 {
  opacity: 0.3;
}
</style>
