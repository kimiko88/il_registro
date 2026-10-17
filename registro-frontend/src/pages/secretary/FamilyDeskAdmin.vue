<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">
          <q-icon name="assignment" class="q-mr-sm" />{{ t('familyDesk.adminTitle') }}
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          {{ t('familyDesk.adminSubtitle') }}
        </p>
      </div>
      <q-btn flat icon="refresh" :label="t('common.refresh')" @click="loadData" :loading="loading" />
    </div>

    <!-- Requests Table -->
    <q-card flat bordered class="bg-white">
      <q-table
        :rows="requests"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        :no-data-label="t('common.noData')"
      >
        <template #body-cell-request_type="props">
          <q-td :props="props">
            <q-badge color="indigo-7" :label="formatRequestType(props.row.request_type)" />
          </q-td>
        </template>

        <template #body-cell-status="props">
          <q-td :props="props">
            <q-chip :color="getStatusColor(props.row.status)" text-color="white" size="sm">
              {{ formatStatus(props.row.status) }}
            </q-chip>
          </q-td>
        </template>

        <template #body-cell-details="props">
          <q-td :props="props">
            <div v-if="props.row.request_type === 'delega_ritiro'" class="text-caption">
              <strong>{{ t('familyDesk.delegateName') }}:</strong> {{ props.row.form_data?.first_name }} {{ props.row.form_data?.last_name }} ({{ props.row.form_data?.relationship }})
              <br />
              <span class="text-grey-7">CF: {{ props.row.form_data?.tax_code }} - Tel: {{ props.row.form_data?.phone }}</span>
            </div>
            <div v-else-if="props.row.request_type === 'uscita_autonoma_under14'" class="text-caption text-positive">
              <q-icon name="check_circle" class="q-mr-xs" /> {{ t('familyDesk.disclaimerUnder14') }}
            </div>
            <div v-else class="text-caption">
              {{ props.row.form_data?.notes || '-' }}
            </div>
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" align="right">
            <div v-if="props.row.status === 'submitted' || props.row.status === 'in_istruttoria'" class="row q-gutter-xs justify-end">
              <q-btn
                size="sm"
                color="positive"
                :label="t('familyDesk.assignProtocol')"
                unelevated
                no-caps
                @click="openApproveDialog(props.row)"
              />
              <q-btn
                size="sm"
                color="negative"
                flat
                :label="t('common.reject')"
                no-caps
                @click="openRejectDialog(props.row)"
              />
            </div>
            <div v-else-if="props.row.status === 'approved' && props.row.request_type === 'delega_ritiro'">
              <q-chip size="xs" color="teal-1" text-color="teal-9" icon="sync">
                {{ t('familyDesk.syncPortineriaChip') }}
              </q-chip>
            </div>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Approve Dialog -->
    <q-dialog v-model="showApproveDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-positive text-white">
          <div class="text-h6">{{ t('familyDesk.assignProtocol') }}</div>
        </q-card-section>
        <q-card-section class="q-pt-md">
          <p class="text-caption">{{ t('agidProtocol.stampTooltip') }}</p>
          <q-input v-model="approvalForm.protocol_number" :label="t('agidProtocol.protocolNumber') + ' *'" outlined dense autofocus />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="positive" :label="t('common.confirm')" unelevated :loading="actionLoading" @click="confirmApproval" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Reject Dialog -->
    <q-dialog v-model="showRejectDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-negative text-white">
          <div class="text-h6">{{ t('common.reject') }}</div>
        </q-card-section>
        <q-card-section class="q-pt-md">
          <q-input v-model="rejectForm.rejection_reason" :label="t('familyDesk.notes') + ' *'" type="textarea" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="negative" :label="t('common.confirm')" unelevated :loading="actionLoading" @click="confirmRejection" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { familyDeskService } from '@/services/familyDeskService'

export default defineComponent({
  name: 'FamilyDeskAdmin',
  setup() {
    const { t } = useI18n()
    const requests = ref([])
    const loading = ref(false)
    const actionLoading = ref(false)
    const showApproveDialog = ref(false)
    const showRejectDialog = ref(false)
    const selectedRequest = ref(null)

    const approvalForm = ref({ protocol_number: '' })
    const rejectForm = ref({ rejection_reason: '' })

    const columns = computed(() => [
      { name: 'created_at', label: t('familyDesk.date'), field: 'created_at', align: 'left', format: val => val ? new Date(val).toLocaleDateString() : '' },
      { name: 'request_type', label: t('familyDesk.requestType'), field: 'request_type', align: 'left' },
      { name: 'details', label: t('common.details'), field: 'form_data', align: 'left' },
      { name: 'status', label: t('common.details'), field: 'status', align: 'center' },
      { name: 'protocol', label: t('agidProtocol.protocolNumber'), field: 'protocol_number', align: 'left' },
      { name: 'actions', label: t('common.actions'), align: 'right' }
    ])

    const loadData = async () => {
      loading.value = true
      try {
        const res = await familyDeskService.getRequests()
        requests.value = res.data || []
      } catch (e) {
        console.error('Errore caricamento pratiche:', e)
      } finally {
        loading.value = false
      }
    }

    const openApproveDialog = (req) => {
      selectedRequest.value = req
      approvalForm.value.protocol_number = `PROT-${new Date().getFullYear()}-${Math.floor(1000 + Math.random() * 9000)}`
      showApproveDialog.value = true
    }

    const openRejectDialog = (req) => {
      selectedRequest.value = req
      rejectForm.value.rejection_reason = ''
      showRejectDialog.value = true
    }

    const confirmApproval = async () => {
      if (!selectedRequest.value) return
      actionLoading.value = true
      try {
        await familyDeskService.reviewRequest(selectedRequest.value.id, {
          status: 'approved',
          protocol_number: approvalForm.value.protocol_number
        })
        showApproveDialog.value = false
        await loadData()
      } catch (e) {
        console.error('Errore approvazione:', e)
      } finally {
        actionLoading.value = false
      }
    }

    const confirmRejection = async () => {
      if (!selectedRequest.value) return
      actionLoading.value = true
      try {
        await familyDeskService.reviewRequest(selectedRequest.value.id, {
          status: 'rejected',
          rejection_reason: rejectForm.value.rejection_reason
        })
        showRejectDialog.value = false
        await loadData()
      } catch (e) {
        console.error('Errore rifiuto:', e)
      } finally {
        actionLoading.value = false
      }
    }

    const formatRequestType = (type) => {
      switch (type) {
        case 'delega_ritiro': return t('familyDesk.delegaRitiro')
        case 'uscita_autonoma_under14': return t('familyDesk.uscitaAutonoma')
        case 'certificato_iscrizione_frequenza': return t('familyDesk.certificato')
        case 'nulla_osta': return t('familyDesk.altro')
        default: return type
      }
    }

    const formatStatus = (s) => {
      switch (s) {
        case 'submitted': return t('familyDesk.statusPending')
        case 'in_istruttoria': return t('familyDesk.statusPending')
        case 'approved': return t('familyDesk.statusApproved')
        case 'rejected': return t('familyDesk.statusRejected')
        default: return s
      }
    }

    const getStatusColor = (s) => {
      switch (s) {
        case 'approved': return 'positive'
        case 'rejected': return 'negative'
        case 'in_istruttoria': return 'warning'
        default: return 'primary'
      }
    }

    onMounted(loadData)

    return {
      t,
      requests,
      loading,
      actionLoading,
      columns,
      showApproveDialog,
      showRejectDialog,
      approvalForm,
      rejectForm,
      loadData,
      openApproveDialog,
      openRejectDialog,
      confirmApproval,
      confirmRejection,
      formatRequestType,
      formatStatus,
      getStatusColor
    }
  }
})
</script>
