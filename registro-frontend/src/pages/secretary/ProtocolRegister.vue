<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">
          <q-icon name="receipt_long" class="q-mr-sm" />{{ t('agidProtocol.registerTitle') }}
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          {{ t('agidProtocol.registerSubtitle') }}
        </p>
      </div>
      <div class="row q-gutter-sm">
        <q-btn flat icon="refresh" :label="t('common.refresh')" @click="loadEntries" :loading="loading" />
        <q-btn color="primary" icon="add" :label="t('agidProtocol.newEntry')" unelevated no-caps @click="showNewProtocolDialog = true" />
      </div>
    </div>

    <!-- Filter by Direction -->
    <q-card flat bordered class="bg-white q-mb-md">
      <q-tabs v-model="flowFilter" dense class="text-grey" active-color="primary" indicator-color="primary" align="left">
        <q-tab name="" :label="t('common.all')" />
        <q-tab name="in" :label="t('agidProtocol.inbound')" />
        <q-tab name="out" :label="t('agidProtocol.outbound')" />
        <q-tab name="internal" :label="t('agidProtocol.internal')" />
      </q-tabs>
    </q-card>

    <!-- Protocol Table -->
    <q-card flat bordered class="bg-white">
      <q-table
        :rows="filteredEntries"
        :columns="columns"
        row-key="id"
        flat
        :loading="loading"
        :no-data-label="t('common.noData')"
      >
        <template #body-cell-protocol_number="props">
          <q-td :props="props">
            <span class="text-weight-bold text-primary font-mono">
              N. {{ String(props.row.protocol_number).padStart(7, '0') }} / {{ props.row.protocol_year }}
            </span>
          </q-td>
        </template>

        <template #body-cell-flow_direction="props">
          <q-td :props="props">
            <q-badge :color="getFlowColor(props.row.flow_direction)" :label="formatFlow(props.row.flow_direction)" />
          </q-td>
        </template>

        <template #body-cell-classification="props">
          <q-td :props="props">
            <span class="text-caption">Tit. {{ props.row.classification_title }} - Cl. {{ props.row.classification_class }}</span>
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" align="right">
            <q-btn
              size="sm"
              color="indigo-7"
              flat
              icon="code"
              :label="t('agidProtocol.exportSegnatura')"
              no-caps
              :href="`/api/v1/protocol/${props.row.id}/segnatura.xml`"
              target="_blank"
            />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- New Protocol Dialog -->
    <q-dialog v-model="showNewProtocolDialog" persistent>
      <q-card style="min-width: 550px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ t('agidProtocol.newEntry') }}</div>
        </q-card-section>
        <q-card-section class="q-pt-md">
          <div class="row q-col-gutter-sm q-mb-sm">
            <div class="col-6">
              <q-select
                v-model="newEntryForm.flow_direction"
                :options="directionOptions"
                emit-value
                map-options
                :label="t('agidProtocol.direction') + ' *'"
                outlined
                dense
              />
            </div>
            <div class="col-6">
              <q-select
                v-model="newEntryForm.classification_title"
                :options="titleOptions"
                emit-value
                map-options
                :label="t('agidProtocol.classification') + ' *'"
                outlined
                dense
              />
            </div>
          </div>

          <q-input v-model="newEntryForm.subject" :label="t('agidProtocol.subject') + ' *'" outlined dense class="q-mb-sm" />
          <div class="row q-col-gutter-sm q-mb-sm">
            <div class="col-6">
              <q-input v-model="newEntryForm.sender" :label="t('agidProtocol.sender') + ' *'" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model="newEntryForm.recipient" :label="t('agidProtocol.recipient') + ' *'" outlined dense />
            </div>
          </div>
          <q-input v-model="newEntryForm.classification_fascicle" :label="t('common.details')" outlined dense class="q-mb-sm" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('agidProtocol.newEntry')" unelevated :loading="submitting" @click="createProtocol" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { protocolService } from '@/services/protocolService'

export default defineComponent({
  name: 'ProtocolRegister',
  setup() {
    const { t } = useI18n()
    const entries = ref([])
    const loading = ref(false)
    const submitting = ref(false)
    const flowFilter = ref('')
    const showNewProtocolDialog = ref(false)

    const newEntryForm = ref({
      flow_direction: 'in',
      classification_title: 7,
      classification_class: '1',
      classification_fascicle: 'Generale',
      subject: '',
      sender: '',
      recipient: ''
    })

    const directionOptions = computed(() => [
      { label: t('agidProtocol.inbound'), value: 'in' },
      { label: t('agidProtocol.outbound'), value: 'out' },
      { label: t('agidProtocol.internal'), value: 'internal' }
    ])

    const titleOptions = [
      { label: 'Titolo I - Organi collegiali', value: 1 },
      { label: 'Titolo II - Personale', value: 2 },
      { label: 'Titolo III - Sicurezza', value: 3 },
      { label: 'Titolo IV - Didattica', value: 4 },
      { label: 'Titolo V - Finanza e patrimonio', value: 5 },
      { label: 'Titolo VI - Edilizia', value: 6 },
      { label: 'Titolo VII - Alunni e famiglie', value: 7 },
      { label: 'Titolo VIII - Rapporti esterni', value: 8 }
    ]

    const columns = computed(() => [
      { name: 'protocol_number', label: t('agidProtocol.protocolNumber'), align: 'left', field: 'protocol_number' },
      { name: 'protocol_date', label: t('agidProtocol.timestamp'), align: 'left', field: 'protocol_date', format: val => val ? new Date(val).toLocaleString() : '' },
      { name: 'flow_direction', label: t('agidProtocol.direction'), align: 'center', field: 'flow_direction' },
      { name: 'subject', label: t('agidProtocol.subject'), align: 'left', field: 'subject' },
      { name: 'sender', label: t('agidProtocol.sender'), align: 'left', field: 'sender' },
      { name: 'recipient', label: t('agidProtocol.recipient'), align: 'left', field: 'recipient' },
      { name: 'classification', label: t('agidProtocol.classification'), align: 'center' },
      { name: 'actions', label: t('common.actions'), align: 'right' }
    ])

    const filteredEntries = computed(() => {
      if (!flowFilter.value) return entries.value
      return entries.value.filter(e => e.flow_direction === flowFilter.value)
    })

    const loadEntries = async () => {
      loading.value = true
      try {
        const res = await protocolService.getEntries()
        entries.value = res.data || []
      } catch (e) {
        console.error('Errore caricamento protocollo:', e)
      } finally {
        loading.value = false
      }
    }

    const createProtocol = async () => {
      submitting.value = true
      try {
        await protocolService.protocolDocument(newEntryForm.value)
        showNewProtocolDialog.value = false
        await loadEntries()
      } catch (e) {
        console.error('Errore protocollo:', e)
      } finally {
        submitting.value = false
      }
    }

    const formatFlow = (f) => {
      switch (f) {
        case 'in': return t('agidProtocol.inbound')
        case 'out': return t('agidProtocol.outbound')
        case 'internal': return t('agidProtocol.internal')
        default: return f
      }
    }

    const getFlowColor = (f) => {
      switch (f) {
        case 'in': return 'teal'
        case 'out': return 'primary'
        case 'internal': return 'blue-grey-8'
        default: return 'grey'
      }
    }

    onMounted(loadEntries)

    return {
      t,
      entries,
      loading,
      submitting,
      flowFilter,
      showNewProtocolDialog,
      newEntryForm,
      directionOptions,
      titleOptions,
      columns,
      filteredEntries,
      loadEntries,
      createProtocol,
      formatFlow,
      getFlowColor
    }
  }
})
</script>
