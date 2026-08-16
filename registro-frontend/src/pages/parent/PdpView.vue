<template>
  <q-page padding class="bg-slate-50 min-h-screen">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="accessibility_new" color="primary" class="q-mr-sm" />
          Piano Didattico Personalizzato (PDP / PEI)
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Documento di programmazione didattica personalizzata approvato dal Consiglio di Classe
        </p>
      </div>
    </div>

    <div v-if="loading" class="q-pa-md">
      <SkeletonCard :lines="8" show-footer />
    </div>

    <q-card v-else-if="plans.length === 0" class="text-center q-pa-xl bg-white rounded-xl shadow-xs border">
      <q-icon name="folder_open" size="64px" color="grey-4" class="q-mb-md" />
      <div class="text-h6 text-weight-bold text-slate-700">Nessun piano didattico condiviso</div>
      <div class="text-caption text-grey-6">Al momento non sono stati condivisi Piani Didattici Personalizzati (PDP/PEI) per tuo figlio.</div>
    </q-card>

    <div v-else class="q-gutter-y-md">
      <q-card v-for="plan in plans" :key="plan.id" class="rounded-xl shadow-xs border bg-white">
        <q-card-section class="bg-indigo-50 text-indigo-9 row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold row items-center">
              <span>{{ plan.plan_type === 'pei' ? 'Piano Educativo Individualizzato (PEI)' : 'Piano Didattico Personalizzato (PDP)' }}</span>
            </div>
            <div class="text-caption opacity-80">Anno Scolastico {{ plan.academic_year }}</div>
          </div>

          <q-chip
            :color="plan.family_approved_at ? 'positive' : 'warning'"
            text-color="white"
            class="text-weight-bold"
          >
            {{ plan.family_approved_at ? '✓ Approvato e Firmato' : 'Richiesta Presa d\'Atto e Approvazione' }}
          </q-chip>
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <!-- Compensative Measures -->
          <div>
            <div class="text-subtitle2 text-weight-bold text-slate-800 q-mb-xs">Misure Compensative Concordate</div>
            <div class="row q-gutter-xs">
              <q-chip
                v-for="m in (plan.content?.compensative || [])"
                :key="m"
                dense
                color="indigo-1"
                text-color="indigo-9"
              >
                ✓ {{ formatMeasure(m) }}
              </q-chip>
              <span v-if="!(plan.content?.compensative?.length)" class="text-caption text-grey-5 font-italic">Nessuna misura indicata</span>
            </div>
          </div>

          <!-- Notes -->
          <div v-if="plan.content?.notes">
            <div class="text-subtitle2 text-weight-bold text-slate-800 q-mb-xs">Indicazioni Pedagogiche e Didattiche</div>
            <div class="text-body2 text-slate-700 bg-slate-100 q-pa-sm rounded-lg border">
              {{ plan.content.notes }}
            </div>
          </div>

          <!-- Approval status info -->
          <div v-if="plan.family_approved_at" class="text-caption text-positive font-bold row items-center q-mt-sm">
            <q-icon name="check_circle" class="q-mr-xs" size="18px" />
            Approvato in data {{ formatDate(plan.family_approved_at) }}
          </div>
        </q-card-section>

        <!-- Approval action -->
        <q-card-actions v-if="!plan.family_approved_at" align="right" class="q-pa-md bg-slate-50 border-t">
          <q-btn
            color="positive"
            icon="verified"
            label="Approva e Accetta PDP"
            unelevated
            class="rounded-lg text-weight-bold"
            :loading="approving"
            @click="approvePlan(plan)"
          />
        </q-card-actions>
      </q-card>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { pdpService } from '@/services/pdpService'
import { useAuthStore } from '@/stores/auth'
import SkeletonCard from '@/components/Common/SkeletonCard.vue'

const $q = useQuasar()
const { t } = useI18n()
const authStore = useAuthStore()

const plans = ref([])
const loading = ref(true)
const approving = ref(false)

onMounted(async () => {
  try {
    const studentId = authStore.user?.id
    if (studentId) {
      const res = await pdpService.getStudentPlans(studentId, '2025/2026')
      plans.value = res.data?.plans || []
    }
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  } finally {
    loading.value = false
  }
})

async function approvePlan(plan) {
  $q.dialog({
    title: t('common.confirm'),
    message: t('pdpPage.confirmApproval') || 'Confermi di aver preso visione e di approvare il Piano Didattico Personalizzato?',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    approving.value = true
    try {
      await pdpService.approveByFamily(plan.id)
      $q.notify({ type: 'positive', message: t('common.success') })
      plan.family_approved_at = new Date().toISOString()
    } catch {
      $q.notify({ type: 'negative', message: t('common.error') })
    } finally {
      approving.value = false
    }
  })
}

function formatMeasure(val) {
  const map = {
    calcolatrice: 'Uso Calcolatrice',
    tempo_aggiuntivo_30: 'Tempo +30%',
    tempo_aggiuntivo_50: 'Tempo +50%',
    prova_equipollente: 'Prova Equipollente',
    sintesi_vocale: 'Sintesi Vocale',
    mappe_concettuali: 'Mappe Concettuali'
  }
  return map[val] || val
}

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('it-IT')
}
</script>
