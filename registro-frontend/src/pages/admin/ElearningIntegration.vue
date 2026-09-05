<template>
  <q-page padding class="bg-slate-50 min-h-screen">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none row items-center">
          <q-icon name="hub" color="primary" class="q-mr-sm" />
          {{ t('elearningPage.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ t('elearningPage.subtitle') }}
        </p>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- Google Classroom Card -->
      <div class="col-12 col-md-6">
        <q-card class="rounded-xl shadow-xs border bg-white full-height">
          <q-card-section class="bg-emerald-50 text-emerald-9 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <q-avatar icon="school" color="emerald" text-color="white" size="40px" />
              <div>
                <div class="text-h6 text-weight-bold">Google Classroom</div>
                <div class="text-caption opacity-80">Google Workspace for Education</div>
              </div>
            </div>
            <q-chip
              :color="googleConnected ? 'positive' : 'grey-4'"
              :text-color="googleConnected ? 'white' : 'grey-8'"
              class="text-weight-bold"
            >
              {{ googleConnected ? '✓ ' + t('elearningPage.connected') : t('elearningPage.notConfigured') }}
            </q-chip>
          </q-card-section>

          <q-card-section class="q-pa-md text-slate-700">
            <p class="text-body2">
              {{ t('elearningPage.googleDesc') }}
            </p>

            <div class="q-gutter-y-xs q-my-md">
              <div class="row items-center"><q-icon name="check_circle" color="emerald" class="q-mr-xs" /> {{ t('elearningPage.googleF1') }}</div>
              <div class="row items-center"><q-icon name="check_circle" color="emerald" class="q-mr-xs" /> {{ t('elearningPage.googleF2') }}</div>
              <div class="row items-center"><q-icon name="check_circle" color="emerald" class="q-mr-xs" /> {{ t('elearningPage.googleF3') }}</div>
            </div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="right" class="q-pa-md bg-slate-50">
            <q-btn
              v-if="!googleConnected"
              color="emerald-7"
              icon="login"
              :label="t('elearningPage.connectGoogle')"
              unelevated
              class="rounded-lg text-weight-bold"
              @click="connectGoogle"
            />
            <div v-else class="row q-gutter-xs">
              <q-btn flat color="primary" icon="sync" :label="t('elearningPage.syncGoogle')" @click="syncGoogleCourses" />
              <q-btn flat color="negative" icon="link_off" :label="t('elearningPage.disconnect')" @click="googleConnected = false" />
            </div>
          </q-card-actions>
        </q-card>
      </div>

      <!-- Microsoft Teams Card -->
      <div class="col-12 col-md-6">
        <q-card class="rounded-xl shadow-xs border bg-white full-height">
          <q-card-section class="bg-indigo-50 text-indigo-9 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <q-avatar icon="groups" color="indigo" text-color="white" size="40px" />
              <div>
                <div class="text-h6 text-weight-bold">Microsoft Teams</div>
                <div class="text-caption opacity-80">Microsoft 365 Education</div>
              </div>
            </div>
            <q-chip
              :color="msConnected ? 'positive' : 'grey-4'"
              :text-color="msConnected ? 'white' : 'grey-8'"
              class="text-weight-bold"
            >
              {{ msConnected ? '✓ ' + t('elearningPage.connected') : t('elearningPage.notConfigured') }}
            </q-chip>
          </q-card-section>

          <q-card-section class="q-pa-md text-slate-700">
            <p class="text-body2">
              {{ t('elearningPage.teamsDesc') }}
            </p>

            <div class="q-gutter-y-xs q-my-md">
              <div class="row items-center"><q-icon name="check_circle" color="indigo" class="q-mr-xs" /> {{ t('elearningPage.teamsF1') }}</div>
              <div class="row items-center"><q-icon name="check_circle" color="indigo" class="q-mr-xs" /> {{ t('elearningPage.teamsF2') }}</div>
              <div class="row items-center"><q-icon name="check_circle" color="indigo" class="q-mr-xs" /> {{ t('elearningPage.teamsF3') }}</div>
            </div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="right" class="q-pa-md bg-slate-50">
            <q-btn
              v-if="!msConnected"
              color="indigo-7"
              icon="login"
              :label="t('elearningPage.connectTeams')"
              unelevated
              class="rounded-lg text-weight-bold"
              @click="connectMicrosoft"
            />
            <div v-else class="row q-gutter-xs">
              <q-btn flat color="primary" icon="sync" :label="t('elearningPage.syncTeams')" @click="syncMsTeams" />
              <q-btn flat color="negative" icon="link_off" :label="t('elearningPage.disconnect')" @click="msConnected = false" />
            </div>
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { elearningService } from '@/services/elearningService'

const $q = useQuasar()
const { t } = useI18n()
const googleConnected = ref(false)
const googleConfigured = ref(false)
const msConnected = ref(false)
const msConfigured = ref(false)
const loading = ref(false)

async function fetchStatus() {
  loading.value = true
  try {
    const res = await elearningService.getProviders()
    const data = res.data
    if (data?.google) {
      googleConnected.value = data.google.connected
      googleConfigured.value = data.google.configured
    }
    if (data?.microsoft) {
      msConnected.value = data.microsoft.connected
      msConfigured.value = data.microsoft.configured
    }
  } catch (err) {
    console.warn('Elearning providers status warning:', err)
  } finally {
    loading.value = false
  }
}

async function connectGoogle() {
  try {
    await elearningService.connectGoogleClassroom('auth-code-demo')
    googleConnected.value = true
    $q.notify({ type: 'positive', message: t('elearningPage.googleConnectedSuccess') })
  } catch (e) {
    $q.notify({ type: 'negative', message: t('elearningPage.googleConnectError') })
  }
}

async function connectMicrosoft() {
  try {
    await elearningService.connectMicrosoftTeams('auth-code-demo')
    msConnected.value = true
    $q.notify({ type: 'positive', message: t('elearningPage.teamsConnectedSuccess') })
  } catch (e) {
    $q.notify({ type: 'negative', message: t('elearningPage.teamsConnectError') })
  }
}

async function syncGoogleCourses() {
  try {
    const res = await elearningService.syncCourses('google')
    $q.notify({ type: 'positive', message: res.data?.message || t('elearningPage.syncSuccess') })
  } catch (e) {
    $q.notify({ type: 'negative', message: t('elearningPage.syncError') })
  }
}

async function syncMsTeams() {
  try {
    const res = await elearningService.syncCourses('microsoft')
    $q.notify({ type: 'positive', message: res.data?.message || t('elearningPage.syncSuccess') })
  } catch (e) {
    $q.notify({ type: 'negative', message: t('elearningPage.syncError') })
  }
}

onMounted(() => {
  fetchStatus()
})
</script>

