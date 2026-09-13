<template>
  <div v-if="notices.length > 0 && isStaffUser" class="strike-banners-container q-mb-md">
    <q-card
      v-for="notice in notices"
      :key="notice.id"
      class="strike-notice-card rounded-xl shadow-2 overflow-hidden border-negative-subtle"
    >
      <!-- Top banner strip -->
      <div class="strike-header-bar q-px-md q-py-sm row items-center justify-between">
        <div class="row items-center gap-sm">
          <q-badge color="negative" text-color="white" class="q-px-sm q-py-xs text-weight-bold rounded-borders">
            <q-icon name="campaign" size="16px" class="q-mr-xs" />
            {{ t('strikeManagement.bannerCommunicationTitle') || 'COMUNICAZIONE DI SCIOPERO' }}
          </q-badge>
          <span class="text-caption text-weight-bold text-slate-700">
            {{ t('strikeManagement.proclaimedBy') || 'Proclamato da:' }} <span class="text-primary">{{ notice.proclaimed_by || t('strikeManagement.defaultProclaimedBy') || 'Organizzazioni Sindacali' }}</span>
          </span>
        </div>

        <div class="row items-center q-gutter-x-sm">
          <q-chip
            dense
            :color="notice.is_expired ? 'grey-7' : 'deep-orange-8'"
            text-color="white"
            class="text-caption text-weight-bold"
          >
            <q-icon :name="notice.is_expired ? 'lock' : 'timer'" size="14px" class="q-mr-xs" />
            {{ notice.is_expired ? (t('strikeManagement.bannerNoticeExpired') || 'Termine Dichiarazioni Scaduto') : deadlineTimeRemaining(notice.declaration_deadline) }}
          </q-chip>

          <q-btn
            v-if="canManage"
            flat
            dense
            no-caps
            color="primary"
            icon="analytics"
            :label="t('strikeManagement.bannerPreventiveDashboard') || 'Quadro Preventivo'"
            to="/ata/strike"
            class="gt-xs text-weight-bold"
          />
        </div>
      </div>

      <!-- Main content body -->
      <q-card-section class="q-pa-md">
        <div class="row items-start justify-between q-col-gutter-md">
          <div class="col-12 col-md-7">
            <h2 class="text-h6 text-weight-bolder text-slate-900 q-my-none">
              {{ notice.title }}
            </h2>
            <div class="text-body2 text-slate-600 q-mt-xs">
              <span class="text-weight-bold">{{ t('strikeManagement.strikeDate') || 'Data Sciopero:' }}</span>
              <q-badge color="indigo-1" text-color="indigo-9" class="q-ml-xs text-weight-bold q-px-sm">
                <q-icon name="event" size="14px" class="q-mr-xs" />
                {{ formatDate(notice.strike_date) }}
              </q-badge>
              <span class="q-mx-sm text-slate-300">•</span>
              <span class="text-weight-bold">{{ t('strikeManagement.bannerDeclarationDeadline') || 'Scadenza Comunicazione:' }}</span>
              <span class="q-ml-xs text-weight-medium" :class="notice.is_expired ? 'text-negative text-weight-bold' : 'text-slate-700'">
                {{ formatDateTime(notice.declaration_deadline) }}
              </span>
            </div>

            <div v-if="notice.notes" class="strike-notes-box bg-slate-50 border border-slate-200 rounded-borders q-pa-sm q-mt-sm text-caption text-slate-700">
              <q-icon name="info" color="primary" class="q-mr-xs" />
              {{ notice.notes }}
            </div>

            <div class="legal-disclaimer text-caption text-slate-500 q-mt-xs">
              {{ t('strikeManagement.bannerLegalDisclaimer') || '* Rilevazione preventiva facoltativa ai sensi dell\'Accordo Aran 2/12/2020. La dichiarazione ha carattere volontario e preventivo.' }}
            </div>
          </div>

          <!-- Interaction Intention Choice Area -->
          <div class="col-12 col-md-5">
            <div class="intention-action-card rounded-xl q-pa-sm border bg-white shadow-xs">
              <div class="text-caption text-weight-bold text-slate-700 q-mb-xs row items-center justify-between">
                <span>{{ t('strikeManagement.bannerYourDeclaration') || 'La tua dichiarazione preventiva:' }}</span>
                <span v-if="notice.user_declaration" class="text-caption text-slate-500">
                  {{ t('strikeManagement.bannerRegisteredAt', { date: formatDateTime(notice.user_declaration.declared_at) }) || `Registrata il ${formatDateTime(notice.user_declaration.declared_at)}` }}
                </span>
              </div>

              <!-- When declarations are OPEN -->
              <div v-if="!notice.is_expired" class="row q-col-gutter-xs">
                <div class="col-12 col-sm-4">
                  <q-btn
                    :unelevated="currentUserIntention(notice) === 'participates'"
                    :outline="currentUserIntention(notice) !== 'participates'"
                    color="positive"
                    class="full-width text-weight-bold py-xs"
                    no-caps
                    :loading="savingNoticeId === notice.id && savingIntention === 'participates'"
                    @click="setDeclaration(notice, 'participates')"
                  >
                    <q-icon
                      :name="currentUserIntention(notice) === 'participates' ? 'check_circle' : 'circle'"
                      size="16px"
                      class="q-mr-xs"
                    />
                    {{ t('strikeManagement.participates') || 'Aderisco' }}
                  </q-btn>
                </div>
                <div class="col-12 col-sm-4">
                  <q-btn
                    :unelevated="currentUserIntention(notice) === 'not_participates'"
                    :outline="currentUserIntention(notice) !== 'not_participates'"
                    color="negative"
                    class="full-width text-weight-bold py-xs"
                    no-caps
                    :loading="savingNoticeId === notice.id && savingIntention === 'not_participates'"
                    @click="setDeclaration(notice, 'not_participates')"
                  >
                    <q-icon
                      :name="currentUserIntention(notice) === 'not_participates' ? 'cancel' : 'circle'"
                      size="16px"
                      class="q-mr-xs"
                    />
                    {{ t('strikeManagement.notParticipates') || 'Non aderisco' }}
                  </q-btn>
                </div>
                <div class="col-12 col-sm-4">
                  <q-btn
                    :unelevated="currentUserIntention(notice) === 'undecided'"
                    :outline="currentUserIntention(notice) !== 'undecided'"
                    color="amber-9"
                    class="full-width text-weight-bold py-xs"
                    no-caps
                    :loading="savingNoticeId === notice.id && savingIntention === 'undecided'"
                    @click="setDeclaration(notice, 'undecided')"
                  >
                    <q-icon
                      :name="currentUserIntention(notice) === 'undecided' ? 'help' : 'circle'"
                      size="16px"
                      class="q-mr-xs"
                    />
                    {{ t('strikeManagement.undecided') || 'Non so ancora' }}
                  </q-btn>
                </div>
              </div>

              <!-- When declarations are CLOSED (EXPIRED) -->
              <div v-else class="q-pa-xs">
                <div class="row items-center justify-between bg-slate-100 rounded-borders q-pa-sm">
                  <div class="row items-center gap-xs">
                    <q-icon name="lock" color="slate-600" size="18px" />
                    <span class="text-caption text-slate-700">{{ t('strikeManagement.bannerAcquired') || 'Dichiarazione acquisita:' }}</span>
                  </div>
                  <div>
                    <q-chip
                      dense
                      :color="getIntentionBadgeColor(currentUserIntention(notice))"
                      text-color="white"
                      class="text-weight-bold text-caption"
                    >
                      {{ getIntentionLabel(currentUserIntention(notice)) }}
                    </q-chip>
                  </div>
                </div>
                <div class="text-caption text-slate-500 q-mt-xs text-italic text-center">
                  {{ t('strikeManagement.legalBannerSubtitle') || 'Il termine per comunicare o modificare la scelta preventiva è scaduto.' }}
                </div>
              </div>

              <!-- Status message confirmation -->
              <div v-if="currentUserIntention(notice) && !notice.is_expired" class="text-caption text-slate-600 q-mt-xs row items-center justify-between">
                <span class="text-positive text-weight-bold">
                  <q-icon name="done_all" size="14px" /> {{ getIntentionLabel(currentUserIntention(notice)) }}
                </span>
                <span class="text-slate-400 text-caption">{{ t('strikeManagement.statusOpen') }}</span>
              </div>
            </div>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import strikeService from '@/services/strikeService'

const $q = useQuasar()
const { t, locale } = useI18n()
const authStore = useAuthStore()
const { userRole } = storeToRefs(authStore)

const notices = ref([])
const savingNoticeId = ref(null)
const savingIntention = ref(null)

const staffRoles = [
  'teacher', 'coordinator', 'dsga', 'assistente_amministrativo',
  'collaboratore_ds', 'collaboratore_scolastico', 'principal',
  'vice_principal', 'secretary', 'admin', 'superadmin'
]

const isStaffUser = computed(() => {
  return staffRoles.includes((userRole.value || '').toLowerCase())
})

const canManage = computed(() => {
  const r = (userRole.value || '').toLowerCase()
  return ['dsga', 'principal', 'vice_principal', 'admin', 'superadmin', 'collaboratore_ds', 'secretary'].includes(r)
})

const currentUserIntention = (notice) => {
  return notice.user_declaration ? notice.user_declaration.intention : null
}

const loadNotices = async () => {
  if (!isStaffUser.value) return
  try {
    const list = await strikeService.getStrikeNotices()
    const yesterday = new Date()
    yesterday.setDate(yesterday.getDate() - 1)
    notices.value = (list || []).filter(n => {
      const strikeDate = new Date(n.strike_date)
      return strikeDate >= yesterday
    })
  } catch {
    notices.value = []
  }
}

const setDeclaration = async (notice, intention) => {
  savingNoticeId.value = notice.id
  savingIntention.value = intention
  try {
    const res = await strikeService.submitDeclaration(notice.id, intention)
    notice.user_declaration = res
    $q.notify({
      type: 'positive',
      message: t('strikeManagement.notifyDeclarationSaved') || 'Dichiarazione preventiva registrata con successo',
      position: 'top',
      timeout: 3000
    })
  } catch (err) {
    const msg = err.response?.data?.error || t('strikeManagement.notifyDeclarationError') || 'Errore nella registrazione della dichiarazione'
    $q.notify({
      type: 'negative',
      message: msg,
      position: 'top'
    })
    loadNotices()
  } finally {
    savingNoticeId.value = null
    savingIntention.value = null
  }
}

const getIntentionLabel = (intention) => {
  switch (intention) {
    case 'participates': return t('strikeManagement.participates') || 'Aderisco'
    case 'not_participates': return t('strikeManagement.notParticipates') || 'Non aderisco'
    case 'undecided': return t('strikeManagement.undecided') || 'Non ho ancora deciso'
    default: return t('strikeManagement.statNotDeclared') || 'Nessuna risposta'
  }
}

const getIntentionBadgeColor = (intention) => {
  switch (intention) {
    case 'participates': return 'positive'
    case 'not_participates': return 'negative'
    case 'undecided': return 'amber-9'
    default: return 'grey-6'
  }
}

const formatDate = (isoDate) => {
  if (!isoDate) return '-'
  try {
    const d = new Date(isoDate)
    return d.toLocaleDateString(locale.value || 'it-IT', { day: '2-digit', month: 'long', year: 'numeric' })
  } catch {
    return isoDate
  }
}

const formatDateTime = (isoDate) => {
  if (!isoDate) return '-'
  try {
    const d = new Date(isoDate)
    return d.toLocaleDateString(locale.value || 'it-IT', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return isoDate
  }
}

const deadlineTimeRemaining = (deadlineIso) => {
  if (!deadlineIso) return ''
  const diff = new Date(deadlineIso).getTime() - Date.now()
  if (diff <= 0) return t('strikeManagement.statusExpired') || 'Scaduto'
  const hours = Math.floor(diff / (1000 * 60 * 60))
  if (hours < 24) {
    return `${hours}h`
  }
  const days = Math.floor(hours / 24)
  return `${days}d`
}

onMounted(() => {
  loadNotices()
})
</script>

<style scoped>
.strike-notice-card {
  border-left: 6px solid var(--q-negative);
  background: #ffffff;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  border-right: 1px solid rgba(0, 0, 0, 0.06);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.strike-header-bar {
  background: #fff5f5;
  border-bottom: 1px solid #fee2e2;
}

.intention-action-card {
  background: #fafafa;
  border: 1px solid #e2e8f0;
}

.strike-notes-box {
  background: #f8fafc;
  line-height: 1.4;
}
</style>
