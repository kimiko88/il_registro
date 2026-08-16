<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-lg">
        <q-btn flat round icon="arrow_back" to="/student" class="q-mr-sm" />
        <h1 class="text-h4 q-my-none">{{ t('settingsPage.profile') || 'Profilo Studente' }}</h1>
    </div>

    <div class="row q-col-gutter-lg">
        <!-- Avatar & Summary -->
        <div class="col-12 col-md-4 text-center">
            <q-card>
                <q-card-section>
                    <q-avatar size="120px" class="q-mb-md shadow-2">
                        <img :src="studentStore.profile?.avatar || 'https://cdn.quasar.dev/img/boy-avatar.png'" />
                        <q-btn round color="primary" icon="edit" size="sm" class="absolute-bottom-right" style="bottom: 0px; right: 0px" />
                    </q-avatar>
                    <div class="text-h5">{{ studentStore.profile?.firstName }} {{ studentStore.profile?.lastName }}</div>
                    <div class="text-subtitle1 text-grey">{{ studentStore.className }}</div>
                    <div class="text-caption text-grey-6">{{ studentStore.profile?.id }}</div>
                </q-card-section>
            </q-card>
        </div>

        <!-- Details & Settings -->
        <div class="col-12 col-md-8">
            <q-card>
                <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify">
                    <q-tab name="info" :label="t('settingsPage.profile') || 'Info Personali'" />
                    <q-tab name="security" :label="t('settingsPage.security') || 'Sicurezza & Account'" />
                </q-tabs>
                <q-separator />

                <q-tab-panels v-model="tab" animated>
                    <q-tab-panel name="info">
                        <q-list>
                            <q-item>
                                <q-item-section>
                                    <q-item-label caption>Email Istituzionale</q-item-label>
                                    <q-item-label>{{ studentStore.profile?.email }}</q-item-label>
                                </q-item-section>
                            </q-item>
                             <q-separator spaced />
                            <q-item>
                                <q-item-section>
                                    <q-item-label caption>Codice Fiscale</q-item-label>
                                    <q-item-label>{{ studentStore.profile?.fiscalCode }}</q-item-label>
                                </q-item-section>
                            </q-item>
                             <q-separator spaced />
                            <q-item>
                                <q-item-section>
                                    <q-item-label caption>Indirizzo</q-item-label>
                                    <q-item-label>Via Roma 1, Milano (MI)</q-item-label>
                                </q-item-section>
                                <q-item-section side>
                                    <q-btn flat icon="edit" color="primary" />
                                </q-item-section>
                            </q-item>
                        </q-list>
                    </q-tab-panel>

                    <q-tab-panel name="security">
                        <q-banner rounded class="bg-blue-1 text-primary q-mb-lg">
                          <template v-slot:avatar><q-icon name="info" color="blue-7" /></template>
                          Per cambiare la password vai nelle <strong>{{ t('settingsPage.title') }}</strong> del tuo account.
                        </q-banner>
                        <q-btn
                          outline color="primary"
                          :label="t('settingsPage.title')"
                          icon="settings"
                          to="/student/settings"
                          class="q-mb-lg"
                        />

                        <q-separator class="q-mb-md" />

                        <div class="text-h6 q-mb-md">{{ t('settingsPage.security') }}</div>
                        <div class="row q-gutter-md">
                            <q-btn outline color="primary" :label="t('documentsPage.download')" icon="archive" />
                            <q-btn flat color="red" :label="t('common.logout')" />
                        </div>
                    </q-tab-panel>
                </q-tab-panels>
            </q-card>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useStudentStore } from '@/stores/student'

const { t } = useI18n()
const studentStore = useStudentStore()
const tab = ref('info')

onMounted(() => {
    studentStore.fetchProfile()
})
</script>
