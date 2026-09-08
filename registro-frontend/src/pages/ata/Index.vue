<template>
  <q-page class="q-pa-md q-pa-lg-xl ata-dashboard-page">
    <!-- Hero Section -->
    <div class="row items-center justify-between q-mb-xl">
      <div class="col-12 col-md-8">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge :color="getRoleColor(userRole)" text-color="white" class="q-px-sm q-py-xs text-weight-bold rounded-borders">
            {{ roleDisplayName }}
          </q-badge>
          <q-badge outline color="primary" class="q-px-sm q-py-xs text-weight-bold">
            {{ t('ataDashboard.badgeAta') || 'Personale ATA' }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-sm-h3 text-weight-bolder text-slate-900 q-my-none text-primary">
          {{ t('ataDashboard.welcome', { name: user?.first_name || 'Operatore' }) }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">
          {{ t('ataDashboard.welcomeSubtitle') }}
        </div>
      </div>
      <div class="col-12 col-md-4 text-right gt-sm">
        <div class="text-caption text-slate-400 text-uppercase letter-spacing-1">{{ t('ataDashboard.today') || 'Oggi' }}</div>
        <div class="text-h6 text-weight-bold text-slate-700">{{ todayFormatted }}</div>
      </div>
    </div>

    <!-- Quick Navigation Hub -->
    <div class="row q-col-gutter-lg q-mb-xl">
      <!-- 1. Presenze Personale & Docenti -->
      <div class="col-12 col-md-6 col-lg-4">
        <q-card
          class="hub-card rounded-2xl p-4 shadow-sm border border-slate-200 cursor-pointer full-height"
          @click="router.push('/ata/attendance')"
        >
          <q-card-section>
            <div class="row items-center justify-between q-mb-md">
              <div class="icon-bubble bg-indigo-100 text-indigo-700">
                <q-icon name="co_present" size="32px" />
              </div>
              <q-icon name="arrow_forward" color="grey-6" size="20px" class="card-arrow" />
            </div>
            <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs">
              {{ t('ataDashboard.cardAttendanceTitle') }}
            </div>
            <div class="text-caption text-slate-600 q-mb-md">
              {{ t('ataDashboard.cardAttendanceDesc') }}
            </div>
            <div class="row items-center text-primary text-weight-bold text-caption">
              {{ t('ataDashboard.cardAttendanceAction') }}
              <q-icon name="chevron_right" size="16px" class="q-ml-xs" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 2. Timbratura & Badge -->
      <div class="col-12 col-md-6 col-lg-4">
        <q-card
          class="hub-card rounded-2xl p-4 shadow-sm border border-slate-200 cursor-pointer full-height"
          @click="router.push('/ata/attendance')"
        >
          <q-card-section>
            <div class="row items-center justify-between q-mb-md">
              <div class="icon-bubble bg-teal-100 text-teal-700">
                <q-icon name="badge" size="32px" />
              </div>
              <q-icon name="arrow_forward" color="grey-6" size="20px" class="card-arrow" />
            </div>
            <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs">
              {{ t('ataDashboard.cardBadgeTitle') }}
            </div>
            <div class="text-caption text-slate-600 q-mb-md">
              {{ t('ataDashboard.cardBadgeDesc') }}
            </div>
            <div class="row items-center text-teal-700 text-weight-bold text-caption">
              {{ t('ataDashboard.cardBadgeAction') }}
              <q-icon name="chevron_right" size="16px" class="q-ml-xs" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 3. Circolari & Comunicazioni -->
      <div class="col-12 col-md-6 col-lg-4">
        <q-card
          class="hub-card rounded-2xl p-4 shadow-sm border border-slate-200 cursor-pointer full-height"
          @click="router.push('/secretary/communications')"
        >
          <q-card-section>
            <div class="row items-center justify-between q-mb-md">
              <div class="icon-bubble bg-amber-100 text-amber-800">
                <q-icon name="campaign" size="32px" />
              </div>
              <q-icon name="arrow_forward" color="grey-6" size="20px" class="card-arrow" />
            </div>
            <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs">
              {{ t('ataDashboard.cardCommsTitle') }}
            </div>
            <div class="text-caption text-slate-600 q-mb-md">
              {{ t('ataDashboard.cardCommsDesc') }}
            </div>
            <div class="row items-center text-amber-800 text-weight-bold text-caption">
              {{ t('ataDashboard.cardCommsAction') }}
              <q-icon name="chevron_right" size="16px" class="q-ml-xs" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 4. Documenti & Atti (per DSGA, AA, Admin) -->
      <div class="col-12 col-md-6 col-lg-4" v-if="userRole === 'dsga' || userRole === 'assistente_amministrativo' || userRole === 'admin' || userRole === 'superadmin'">
        <q-card
          class="hub-card rounded-2xl p-4 shadow-sm border border-slate-200 cursor-pointer full-height"
          @click="router.push('/secretary/documents')"
        >
          <q-card-section>
            <div class="row items-center justify-between q-mb-md">
              <div class="icon-bubble bg-purple-100 text-purple-700">
                <q-icon name="folder" size="32px" />
              </div>
              <q-icon name="arrow_forward" color="grey-6" size="20px" class="card-arrow" />
            </div>
            <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs">
              {{ t('ataDashboard.cardDocsTitle') }}
            </div>
            <div class="text-caption text-slate-600 q-mb-md">
              {{ t('ataDashboard.cardDocsDesc') }}
            </div>
            <div class="row items-center text-purple-700 text-weight-bold text-caption">
              {{ t('ataDashboard.cardDocsAction') }}
              <q-icon name="chevron_right" size="16px" class="q-ml-xs" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 5. Gestione Sostituzioni Docenti (per Collaboratore DS, DSGA, Admin) -->
      <div class="col-12 col-md-6 col-lg-4" v-if="userRole === 'collaboratore_ds' || userRole === 'dsga' || userRole === 'admin' || userRole === 'superadmin'">
        <q-card
          class="hub-card rounded-2xl p-4 shadow-sm border border-slate-200 cursor-pointer full-height"
          @click="router.push('/secretary/substitutions')"
        >
          <q-card-section>
            <div class="row items-center justify-between q-mb-md">
              <div class="icon-bubble bg-deep-orange-100 text-deep-orange-700">
                <q-icon name="swap_horiz" size="32px" />
              </div>
              <q-icon name="arrow_forward" color="grey-6" size="20px" class="card-arrow" />
            </div>
            <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs">
              {{ t('ataDashboard.cardSubstitutionsTitle') }}
            </div>
            <div class="text-caption text-slate-600 q-mb-md">
              {{ t('ataDashboard.cardSubstitutionsDesc') }}
            </div>
            <div class="row items-center text-deep-orange-700 text-weight-bold text-caption">
              {{ t('ataDashboard.cardSubstitutionsAction') }}
              <q-icon name="chevron_right" size="16px" class="q-ml-xs" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 6. Anagrafica Utenti / Personale (per DSGA, AA, Admin) -->
      <div class="col-12 col-md-6 col-lg-4" v-if="userRole === 'dsga' || userRole === 'assistente_amministrativo' || userRole === 'admin' || userRole === 'superadmin'">
        <q-card
          class="hub-card rounded-2xl p-4 shadow-sm border border-slate-200 cursor-pointer full-height"
          @click="router.push('/secretary/users')"
        >
          <q-card-section>
            <div class="row items-center justify-between q-mb-md">
              <div class="icon-bubble bg-blue-100 text-blue-700">
                <q-icon name="people" size="32px" />
              </div>
              <q-icon name="arrow_forward" color="grey-6" size="20px" class="card-arrow" />
            </div>
            <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs">
              {{ t('ataDashboard.cardUsersTitle') }}
            </div>
            <div class="text-caption text-slate-600 q-mb-md">
              {{ t('ataDashboard.cardUsersDesc') }}
            </div>
            <div class="row items-center text-blue-700 text-weight-bold text-caption">
              {{ t('ataDashboard.cardUsersAction') }}
              <q-icon name="chevron_right" size="16px" class="q-ml-xs" />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'

const { t, te, locale } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const { user, userRole } = storeToRefs(authStore)

const todayFormatted = computed(() => {
  const d = new Date()
  return d.toLocaleDateString(locale.value || 'it-IT', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
})

const roleDisplayName = computed(() => {
  const role = (userRole.value || '').toLowerCase()
  if (te('roles.' + role)) {
    return t('roles.' + role)
  }
  const map = {
    dsga: 'DSGA (Direttore dei Servizi Generali e Amministrativi)',
    collaboratore_ds: 'Collaboratore del Dirigente Scolastico',
    assistente_amministrativo: 'Assistente Amministrativo',
    collaboratore_scolastico: 'Collaboratore Scolastico',
    principal: 'Dirigente Scolastico',
    vice_principal: 'Collaboratore Vicario',
    admin: 'Amministratore',
    superadmin: 'Super Amministratore'
  }
  return map[role] || role
})

function getRoleColor(role) {
  const map = {
    dsga: 'teal-8',
    collaboratore_ds: 'deep-orange-7',
    assistente_amministrativo: 'cyan-8',
    collaboratore_scolastico: 'amber-9',
    principal: 'purple-8',
    vice_principal: 'indigo-8'
  }
  return map[role] || 'primary'
}
</script>

<style scoped>
.hub-card {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  background: white;
}
.hub-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 14px 30px -6px rgba(0, 0, 0, 0.1);
  border-color: #6366f1;
}
.hub-card:hover .card-arrow {
  color: #6366f1;
  transform: translateX(4px);
}
.card-arrow {
  transition: transform 0.2s ease, color 0.2s ease;
}
.icon-bubble {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.letter-spacing-1 {
  letter-spacing: 0.05em;
}
</style>
