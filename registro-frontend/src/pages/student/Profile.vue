<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-lg">
        <q-btn flat round icon="arrow_back" to="/student" class="q-mr-sm" />
        <h1 class="text-h4 q-my-none">Profilo Studente</h1>
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
                    <q-tab name="info" label="Info Personali" />
                    <q-tab name="security" label="Sicurezza & Account" />
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
                        <div class="text-h6 q-mb-md">Gestione Password</div>
                        <q-banner class="bg-blue-1 text-blue-9 rounded-borders q-mb-md" dense>
                          <template v-slot:avatar><q-icon name="info" color="blue-7" /></template>
                          Per cambiare la password vai nelle <strong>Impostazioni</strong> del tuo account.
                        </q-banner>
                        <q-btn
                          outline color="primary"
                          label="Vai alle Impostazioni"
                          icon="settings"
                          to="/student/settings"
                          class="q-mb-lg"
                        />

                        <q-separator class="q-mb-md" />

                        <div class="text-h6 q-mb-md">Zona Pericolo</div>
                        <p>Richiedi una copia dei tuoi dati (GDPR) o gestisci la privacy.</p>
                        <div class="row q-gutter-md">
                            <q-btn outline color="primary" label="Export Dati (GDPR)" icon="archive" />
                            <q-btn flat color="red" label="Logout da tutti i dispositivi" />
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
import { useStudentStore } from 'src/stores/student'

const studentStore = useStudentStore()
const tab = ref('info')

onMounted(() => {
    studentStore.fetchProfile()
})
</script>
