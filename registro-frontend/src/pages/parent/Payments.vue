<template>
  <q-page class="q-pa-md" role="main">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          {{ t('paymentsPage.title') }}
        </h1>
        <div class="text-subtitle1 text-slate-600 q-mt-xs">
          {{ t('paymentsPage.subtitle') }}
        </div>
      </div>
      <q-chip color="primary" text-color="white" class="text-weight-bold q-py-md q-px-md shadow-sm" data-testid="q-chip">
        {{ t('paymentsPage.totalPending') }} {{ totalPending }}
      </q-chip>
    </div>

    <!-- Loading spinner -->
    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner color="primary" size="3em" />
    </div>

    <div v-else class="row q-col-gutter-lg">
      <div class="col-12 col-md-8">
        <q-card flat class="glass-card rounded-xl overflow-hidden shadow-1">
          <q-tabs v-model="tab" dense class="text-grey-7 bg-white" active-color="primary" indicator-color="primary" align="justify">
            <q-tab name="pending" :label="t('paymentsPage.toPay')" icon="payment" class="text-amber-9" />
            <q-tab name="history" :label="t('paymentsPage.history')" icon="history" />
          </q-tabs>
          <q-separator />

          <q-tab-panels v-model="tab" animated>
            <q-tab-panel name="pending" class="q-pa-none">
              <q-list separator>
                <q-item v-for="item in pendingItems" :key="item.id" class="q-py-md">
                  <q-item-section avatar>
                    <q-avatar color="amber-1" text-color="amber-9" icon="receipt" size="44px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-subtitle1 text-slate-800">{{ item.title }}</q-item-label>
                    <q-item-label v-if="item.description" caption class="text-slate-700 q-mb-xs">{{ item.description }}</q-item-label>
                    <q-item-label caption class="text-slate-600">{{ t('paymentsPage.dueDate') }}: {{ formatDate(item.due_date) }}</q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <div class="text-h6 text-weight-bold text-primary">€ {{ formatAmount(item.amount) }}</div>
                    <q-btn
                      :label="t('paymentsPage.payWithPagoPA')"
                      color="primary"
                      unelevated
                      no-caps
                      size="sm"
                      class="q-mt-xs rounded-lg"
                      :loading="payingId === item.id"
                      @click="pay(item)"
                    />
                  </q-item-section>
                </q-item>
                <q-item v-if="pendingItems.length === 0" class="text-center q-pa-xl">
                  <q-item-section class="text-center text-positive">
                    <q-icon name="check_circle" size="56px" class="q-mb-sm" />
                    <div class="text-h6 text-weight-bold">{{ t('paymentsPage.allClear') }}</div>
                    <div class="text-caption text-slate-500">{{ t('paymentsPage.noPendingDesc') }}</div>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-tab-panel>

            <q-tab-panel name="history" class="q-pa-none">
              <q-list separator>
                <q-item v-for="item in historyItems" :key="item.id" class="q-py-md">
                  <q-item-section avatar>
                    <q-avatar color="emerald-1" text-color="emerald-7" icon="verified" size="44px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold text-slate-800">{{ item.title }}</q-item-label>
                    <q-item-label v-if="item.receipt_number" caption class="text-primary text-weight-medium">
                      {{ t('paymentsPage.receiptNumber') }}: {{ item.receipt_number }}
                    </q-item-label>
                    <q-item-label caption class="text-slate-600">
                      {{ t('paymentsPage.paidOn') }}: {{ formatDateTime(item.paid_at) }}
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <div class="text-weight-bold text-subtitle1">€ {{ formatAmount(item.amount) }}</div>
                    <q-btn
                      flat
                      round
                      icon="download"
                      color="primary"
                      size="sm"
                      aria-label="Scarica ricevuta"
                      @click="downloadReceipt(item)"
                    >
                      <q-tooltip>{{ t('paymentsPage.downloadReceipt') }}</q-tooltip>
                    </q-btn>
                  </q-item-section>
                </q-item>
                <q-item v-if="historyItems.length === 0" class="text-center q-pa-xl">
                  <q-item-section class="text-center text-grey-6">
                    <q-icon name="history" size="56px" class="q-mb-sm" />
                    <div class="text-h6">{{ t('paymentsPage.noHistoryDesc') }}</div>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-tab-panel>
          </q-tab-panels>
        </q-card>
      </div>

      <div class="col-12 col-md-4">
        <q-card flat class="glass-card rounded-xl q-pa-md shadow-1">
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold text-slate-800 row items-center">
              <q-icon name="info" color="primary" class="q-mr-xs" />
              {{ t('paymentsPage.infoTitle') }}
            </div>
            <p class="text-body2 text-slate-600 q-mt-sm">
              {{ t('paymentsPage.infoDesc') }}
            </p>
            <q-btn flat :label="t('paymentsPage.guideButton')" color="primary" class="full-width rounded-lg q-mt-xs" no-caps @click="openPagoPAGuide" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Post Payment Summary Dialog -->
    <q-dialog v-model="showReceiptDialog">
      <q-card style="width: min(480px, 95vw)" class="rounded-xl overflow-hidden">
        <q-card-section class="bg-positive text-white text-center q-pa-lg">
          <q-icon name="check_circle" size="64px" class="q-mb-xs" />
          <div class="text-h5 text-weight-bold">{{ t('paymentsPage.paymentCompleted') }}</div>
          <div class="text-subtitle2 opacity-90">{{ t('paymentsPage.paymentSuccessDesc') }}</div>
        </q-card-section>

        <q-card-section class="q-pa-md" v-if="completedPayment">
          <q-list class="q-gutter-y-xs">
            <q-item>
              <q-item-section>
                <q-item-label caption>{{ t('paymentsPage.cause') }}</q-item-label>
                <q-item-label class="text-weight-bold text-subtitle1">{{ completedPayment.title }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label caption>{{ t('paymentsPage.amountPaid') }}</q-item-label>
                <q-item-label class="text-weight-bold text-h6 text-primary">€ {{ formatAmount(completedPayment.amount) }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item v-if="completedPayment.receipt_number">
              <q-item-section>
                <q-item-label caption>{{ t('paymentsPage.receiptNumber') }}</q-item-label>
                <q-item-label class="text-weight-bold text-body1 text-grey-9">{{ completedPayment.receipt_number }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label caption>{{ t('paymentsPage.txDateTime') }}</q-item-label>
                <q-item-label class="text-weight-medium">{{ formatDateTime(completedPayment.paid_at) }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-separator />

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="t('common.close')" v-close-popup no-caps />
          <q-btn color="primary" icon="download" :label="t('paymentsPage.downloadReceipt')" unelevated no-caps class="rounded-lg" @click="downloadReceipt(completedPayment)" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import paymentService from '@/services/paymentService'
import { useChildrenStore } from '@/stores/children'

const $q = useQuasar()
const { t } = useI18n()
const childrenStore = useChildrenStore()

const tab = ref('pending')
const loading = ref(false)
const payingId = ref(null)
const showReceiptDialog = ref(false)
const completedPayment = ref(null)
const paymentsList = ref([])

const pendingItems = computed(() => {
  return paymentsList.value.filter(item => item.status === 'pending')
})

const historyItems = computed(() => {
  return paymentsList.value.filter(item => item.status === 'paid')
})

const totalPending = computed(() => {
  return pendingItems.value.reduce((acc, item) => acc + Number(item.amount || 0), 0).toFixed(2)
})

const formatAmount = (amount) => {
  return Number(amount || 0).toFixed(2)
}

const formatDate = (val) => {
  if (!val) return ''
  try {
    const d = new Date(val)
    return isNaN(d.getTime()) ? val : d.toLocaleDateString('it-IT')
  } catch {
    return val
  }
}

const formatDateTime = (val) => {
  if (!val) return new Date().toLocaleString('it-IT')
  try {
    const d = new Date(val)
    return isNaN(d.getTime()) ? val : d.toLocaleString('it-IT')
  } catch {
    return val
  }
}

const openPagoPAGuide = () => {
  $q.dialog({
    title: t('paymentsPage.infoTitle'),
    message: t('paymentsPage.infoDesc'),
    ok: { label: t('common.close'), color: 'primary' }
  })
}

const loadPayments = async () => {
  loading.value = true
  try {
    const studentId = childrenStore.selectedChildId || childrenStore.children[0]?.id
    const params = studentId ? { student_id: studentId } : {}
    const res = await paymentService.getPayments(params)
    const data = res.data
    paymentsList.value = Array.isArray(data) ? data : (data?.payments || [])
  } catch (err) {
    paymentsList.value = []
    console.error('Error fetching payments from database:', err)
  } finally {
    loading.value = false
  }
}

const pay = async (item) => {
  payingId.value = item.id
  try {
    const res = await paymentService.pay(item.id, { payment_method: 'PagoPA' })
    const updated = res.data?.payment || { ...item, status: 'paid', paid_at: new Date().toISOString() }
    
    // Update local state
    const index = paymentsList.value.findIndex(p => p.id === item.id)
    if (index !== -1) {
      paymentsList.value[index] = updated
    } else {
      item.status = 'paid'
      item.paid_at = new Date().toISOString()
    }
    
    completedPayment.value = updated
    showReceiptDialog.value = true
    $q.notify({
      type: 'positive',
      message: t('paymentsPage.paymentSuccessDesc') || t('common.success')
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('common.error')
    })
  } finally {
    payingId.value = null
  }
}

const downloadReceipt = (item) => {
  $q.notify({
    type: 'positive',
    message: t('paymentsPage.receiptDownloaded'),
    icon: 'download',
    timeout: 2500
  })
}

watch(
  () => childrenStore.selectedChildId,
  () => {
    loadPayments()
  }
)

onMounted(async () => {
  if (childrenStore.children.length === 0) {
    await childrenStore.fetchChildren()
  }
  await loadPayments()
})
</script>

<style scoped>
</style>

