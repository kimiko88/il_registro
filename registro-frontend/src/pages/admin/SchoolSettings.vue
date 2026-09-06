<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none">{{ t('schoolSettingsPage.title') }}</h1>
        <div class="text-subtitle2 text-grey-7">{{ t('schoolSettingsPage.subtitle') }}</div>
      </div>
      <q-space />
      <q-btn
        color="primary"
        icon="save"
        :label="t('schoolSettingsPage.saveChanges')"
        :loading="saving"
        @click="saveSettings"
      />
    </div>

    <q-card v-if="loading" class="q-pa-lg text-center shadow-2">
      <q-spinner color="primary" size="3em" />
      <div class="q-mt-sm text-grey-7">{{ t('schoolSettingsPage.loadingSettings') }}</div>
    </q-card>

    <div v-else class="row q-col-gutter-md">
      <!-- Section 1: Note & Disciplina -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-2 rounded-borders">
          <q-card-section class="bg-primary text-white">
            <div class="text-h6"><q-icon name="gavel" class="q-mr-sm" />{{ t('schoolSettingsPage.notesAndDiscipline') }}</div>
            <div class="text-caption opacity-80">{{ t('schoolSettingsPage.notesAndDisciplineDesc') }}</div>
          </q-card-section>
          <q-card-section class="q-pa-md">
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('schoolSettingsPage.principalApprovalForNotes') }}</q-item-label>
                <q-item-label caption>
                  {{ t('schoolSettingsPage.principalApprovalForNotesDesc') }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.require_principal_approval_for_notes" color="secondary" />
              </q-item-section>
            </q-item>
          </q-card-section>
        </q-card>
      </div>

      <!-- Section 2: Visibilità Voti & Valutazioni -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-2 rounded-borders">
          <q-card-section class="bg-indigo-9 text-white">
            <div class="text-h6"><q-icon name="visibility" class="q-mr-sm" />{{ t('schoolSettingsPage.gradesVisibility') }}</div>
            <div class="text-caption opacity-80">{{ t('schoolSettingsPage.gradesVisibilityDesc') }}</div>
          </q-card-section>
          <q-card-section class="q-pa-md">
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('schoolSettingsPage.showGradesToParents') }}</q-item-label>
                <q-item-label caption>
                  {{ t('schoolSettingsPage.showGradesToParentsDesc') }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.allow_parents_view_grades" color="indigo" />
              </q-item-section>
            </q-item>

            <q-separator q-my-xs />

            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('schoolSettingsPage.showAveragesToStudents') }}</q-item-label>
                <q-item-label caption>
                  {{ t('schoolSettingsPage.showAveragesToStudentsDesc') }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.allow_students_view_class_averages" color="indigo" />
              </q-item-section>
            </q-item>
          </q-card-section>
        </q-card>
      </div>

      <!-- Section 3: Scrutini & Flussi -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-2 rounded-borders">
          <q-card-section class="bg-teal-8 text-white">
            <div class="text-h6"><q-icon name="assessment" class="q-mr-sm" />{{ t('schoolSettingsPage.scrutinyWorkflow') }}</div>
            <div class="text-caption opacity-80">{{ t('schoolSettingsPage.scrutinyWorkflowDesc') }}</div>
          </q-card-section>
          <q-card-section class="q-pa-md">
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('schoolSettingsPage.lockScrutinyAfterValidation') }}</q-item-label>
                <q-item-label caption>
                  {{ t('schoolSettingsPage.lockScrutinyAfterValidationDesc') }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.lock_scrutiny_editing_after_validation" color="teal" />
              </q-item-section>
            </q-item>
          </q-card-section>
        </q-card>
      </div>

      <!-- Section 4: Sicurezza & Supplenze -->
      <div class="col-12 col-md-6">
        <q-card class="shadow-2 rounded-borders">
          <q-card-section class="bg-deep-orange-8 text-white">
            <div class="text-h6"><q-icon name="security" class="q-mr-sm" />{{ t('schoolSettingsPage.securityAndSubstitutes') }}</div>
            <div class="text-caption opacity-80">{{ t('schoolSettingsPage.securityAndSubstitutesDesc') }}</div>
          </q-card-section>
          <q-card-section class="q-pa-md">
            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('schoolSettingsPage.mfaForStaff') }}</q-item-label>
                <q-item-label caption>
                  {{ t('schoolSettingsPage.mfaForStaffDesc') }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.require_mfa_for_staff" color="deep-orange" />
              </q-item-section>
            </q-item>

            <q-separator q-my-xs />

            <q-item tag="label" v-ripple>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('schoolSettingsPage.substituteNotifications') }}</q-item-label>
                <q-item-label caption>
                  {{ t('schoolSettingsPage.substituteNotificationsDesc') }}
                </q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.enable_substitute_notifications" color="deep-orange" />
              </q-item-section>
            </q-item>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import schoolSettingsService from '@/services/schoolSettingsService';

const $q = useQuasar();
const { t } = useI18n();
const loading = ref(true);
const saving = ref(false);

const settings = ref({
  require_principal_approval_for_notes: false,
  allow_parents_view_grades: true,
  allow_students_view_class_averages: true,
  require_mfa_for_staff: false,
  lock_scrutiny_editing_after_validation: true,
  enable_substitute_notifications: true
});

const fetchSettings = async () => {
  loading.value = true;
  try {
    const res = await schoolSettingsService.getSettings();
    settings.value = res.data;
  } catch (err) {
    $q.notify({ type: 'negative', message: `${t('schoolSettingsPage.loadError')}: ${err.response?.data?.error || err.message}` });
  } finally {
    loading.value = false;
  }
};

const saveSettings = async () => {
  saving.value = true;
  try {
    await schoolSettingsService.updateSettings(settings.value);
    $q.notify({ type: 'positive', message: t('schoolSettingsPage.saveSuccess') });
  } catch (err) {
    $q.notify({ type: 'negative', message: `${t('schoolSettingsPage.saveError')}: ${err.response?.data?.error || err.message}` });
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  fetchSettings();
});
</script>
