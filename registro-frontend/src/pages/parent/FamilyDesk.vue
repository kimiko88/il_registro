<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">
          <q-icon name="description" class="q-mr-sm" />{{ t('familyDesk.title') }}
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          {{ t('familyDesk.subtitle') }}
        </p>
      </div>
      <q-btn
        color="primary"
        icon="add"
        :label="t('familyDesk.newRequest')"
        unelevated
        no-caps
        @click="openNewRequestDialog"
      />
    </div>

    <!-- Active Delegates Section -->
    <q-card flat bordered class="q-mb-md bg-white">
      <q-card-section class="q-py-sm bg-blue-grey-1 row items-center justify-between">
        <div class="text-subtitle2 text-weight-bold text-blue-grey-9">
          <q-icon name="badge" class="q-mr-xs" /> {{ t('familyDesk.activeDelegates') }}
        </div>
        <q-badge color="teal" :label="`${delegates.length} ${t('common.active')}`" />
      </q-card-section>
      <q-separator />
      <q-card-section class="q-pa-none">
        <div v-if="delegates.length === 0" class="q-pa-md text-caption text-grey-6 text-center">
          {{ t('common.noData') }}
        </div>
        <q-list v-else separator>
          <q-item v-for="del in delegates" :key="del.id">
            <q-item-section avatar>
              <q-avatar color="primary" text-color="white" icon="person" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-bold">{{ del.first_name }} {{ del.last_name }} ({{ del.relationship }})</q-item-label>
              <q-item-label caption>CF: {{ del.tax_code }} | Tel: {{ del.phone }} | Doc: {{ del.id_card_details }}</q-item-label>
            </q-item-section>
            <q-item-section side>
              <q-chip color="positive" text-color="white" size="sm" icon="check_circle">{{ t('common.active') }}</q-chip>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>

    <!-- Requests Tabs & List -->
    <q-card flat bordered class="bg-white">
      <q-tabs v-model="statusFilter" dense class="text-grey" active-color="primary" indicator-color="primary" align="left">
        <q-tab name="" :label="t('common.all')" />
        <q-tab name="submitted" :label="t('familyDesk.statusPending')" />
        <q-tab name="approved" :label="t('familyDesk.statusApproved')" />
        <q-tab name="rejected" :label="t('familyDesk.statusRejected')" />
      </q-tabs>
      <q-separator />

      <q-table
        :rows="filteredRequests"
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

        <template #body-cell-protocol="props">
          <q-td :props="props">
            <span class="text-weight-bold text-primary">{{ props.row.protocol_number || '-' }}</span>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- New Request Dialog -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 550px; max-width: 90vw;">
        <q-card-section class="row items-center justify-between bg-primary text-white">
          <div class="text-h6">{{ t('familyDesk.newRequest') }}</div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-select
            v-model="form.request_type"
            :options="requestTypeOptions"
            emit-value
            map-options
            :label="t('familyDesk.requestType') + ' *'"
            outlined
            dense
            class="q-mb-md"
          />

          <!-- Form Delega Ritiro -->
          <div v-if="form.request_type === 'delega_ritiro'">
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-input v-model="form.formData.first_name" :label="t('common.name') + ' *'" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model="form.formData.last_name" :label="t('common.fullName') + ' *'" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model="form.formData.tax_code" :label="t('familyDesk.delegateCf') + ' *'" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model="form.formData.relationship" :label="t('familyDesk.delegateName') + ' *'" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model="form.formData.phone" :label="t('familyDesk.delegatePhone') + ' *'" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model="form.formData.id_card_details" :label="t('familyDesk.delegateDoc') + ' *'" outlined dense />
              </div>
            </div>
          </div>

          <!-- Form Uscita Autonoma Under 14 -->
          <div v-else-if="form.request_type === 'uscita_autonoma_under14'">
            <q-banner rounded class="bg-amber-1 text-amber-9 q-mb-md text-caption">
              {{ t('familyDesk.disclaimerUnder14') }}
            </q-banner>
            <q-checkbox
              v-model="form.formData.consent_given"
              :label="t('familyDesk.uscitaAutonoma')"
            />
            <q-checkbox
              v-model="form.formData.exempt_school_rep"
              :label="t('familyDesk.disclaimerUnder14')"
            />
          </div>

          <!-- Generic notes for other requests -->
          <div v-else>
            <q-input
              v-model="form.formData.notes"
              type="textarea"
              :label="t('familyDesk.notes') + ' *'"
              outlined
              dense
              rows="3"
            />
          </div>
        </q-card-section>

        <q-separator />

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('common.save')" unelevated :loading="submitting" @click="submitRequest" />
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
  name: 'FamilyDesk',
  setup() {
    const { t } = useI18n()
    const requests = ref([])
    const delegates = ref([])
    const loading = ref(false)
    const submitting = ref(false)
    const statusFilter = ref('')
    const showDialog = ref(false)

    const form = ref({
      student_id: 'student-default',
      request_type: 'delega_ritiro',
      formData: {
        first_name: '',
        last_name: '',
        tax_code: '',
        relationship: '',
        phone: '',
        id_card_details: '',
        consent_given: false,
        exempt_school_rep: false,
        notes: ''
      }
    })

    const requestTypeOptions = computed(() => [
      { label: t('familyDesk.delegaRitiro'), value: 'delega_ritiro' },
      { label: t('familyDesk.uscitaAutonoma'), value: 'uscita_autonoma_under14' },
      { label: t('familyDesk.certificato'), value: 'certificato_iscrizione_frequenza' },
      { label: t('familyDesk.accessoAtti'), value: 'accesso_atti' },
      { label: t('familyDesk.altro'), value: 'altro' }
    ])

    const columns = computed(() => [
      { name: 'created_at', label: t('familyDesk.date'), field: 'created_at', align: 'left', format: val => val ? new Date(val).toLocaleDateString() : '' },
      { name: 'request_type', label: t('familyDesk.requestType'), field: 'request_type', align: 'left' },
      { name: 'status', label: t('common.details'), field: 'status', align: 'center' },
      { name: 'protocol', label: t('agidProtocol.protocolNumber'), field: 'protocol_number', align: 'left' }
    ])

    const filteredRequests = computed(() => {
      if (!statusFilter.value) return requests.value
      return requests.value.filter(r => r.status === statusFilter.value)
    })

    const loadData = async () => {
      loading.value = true
      try {
        const reqRes = await familyDeskService.getRequests()
        requests.value = reqRes.data || []
        const delRes = await familyDeskService.getDelegates()
        delegates.value = delRes.data || []
      } catch (e) {
        console.error('Errore nel caricamento sportello famiglie:', e)
      } finally {
        loading.value = false
      }
    }

    const openNewRequestDialog = () => {
      form.value.formData = {
        first_name: '',
        last_name: '',
        tax_code: '',
        relationship: '',
        phone: '',
        id_card_details: '',
        consent_given: false,
        exempt_school_rep: false,
        notes: ''
      }
      showDialog.value = true
    }

    const submitRequest = async () => {
      submitting.value = true
      try {
        await familyDeskService.submitRequest({
          student_id: form.value.student_id,
          request_type: form.value.request_type,
          form_data: form.value.formData
        })
        showDialog.value = false
        await loadData()
      } catch (e) {
        console.error('Errore invio istanza:', e)
      } finally {
        submitting.value = false
      }
    }

    const formatRequestType = (type) => {
      const match = requestTypeOptions.value.find(o => o.value === type)
      return match ? match.label : type
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
      delegates,
      loading,
      submitting,
      statusFilter,
      showDialog,
      form,
      requestTypeOptions,
      columns,
      filteredRequests,
      openNewRequestDialog,
      submitRequest,
      formatRequestType,
      formatStatus,
      getStatusColor
    }
  }
})
</script>
