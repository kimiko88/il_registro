<script setup>
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { pctoService } from '@/services/pctoService'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const { t } = useI18n()
const projects = ref([])
const loading = ref(false)

const totalHours = computed(() => {
    return projects.value.reduce((acc, p) => acc + (p.hours_done || 0), 0)
})

const targetHours = 90 // Should probably come from settings/profile

const progressValue = computed(() => {
    if (targetHours === 0) return 0
    return Math.min(totalHours.value / targetHours, 1)
})

const fetchPCTO = async () => {
    loading.value = true
    try {
        const res = await pctoService.getMyProjects()
        projects.value = res.data || []
    } catch (e) {
        $q.notify({ message: t('common.error'), color: 'negative' })
        console.error(e)
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    fetchPCTO()
})
</script>

<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-lg">
       <div class="text-h4">PCTO - Percorsi Trasversali</div>
       <q-chip color="orange" text-color="white" icon="timer">Totale: {{ totalHours }} / {{ targetHours }} Ore</q-chip>
    </div>

    <!-- Progress Bar -->
    <q-card class="q-mb-lg">
        <q-card-section>
            <div class="text-subtitle1 q-mb-sm">{{ t('competenciesPage.title') }}</div>
            <q-linear-progress size="25px" :value="progressValue" color="primary" stripe rounded>
                <div class="absolute-full flex flex-center">
                    <q-badge color="white" text-color="primary" :label="Math.round(progressValue * 100) + '%'" />
                </div>
            </q-linear-progress>
        </q-card-section>
    </q-card>

    <div class="row q-col-gutter-lg">
        <!-- Projects List -->
        <div class="col-12 col-md-8">
            <div class="text-h5 q-mb-md">{{ t('didacticsPage.title') }}</div>
            
            <div v-if="loading" class="flex flex-center q-pa-xl">
                <q-spinner color="primary" size="3em" />
            </div>

            <q-list v-else bordered class="rounded-borders bg-white">
                <q-expansion-item
                    v-for="project in projects"
                    :key="project.id"
                    expand-separator
                    icon="business"
                    :label="project.title"
                    :caption="`${project.company_name} - ${project.hours_done || 0} / ${project.total_hours} Ore`"
                    header-class="text-primary"
                >
                    <q-card>
                        <q-card-section>
                            <div class="row q-col-gutter-md">
                                <div class="col-12 col-md-6">
                                    <div class="text-weight-bold">Tutor Aziendale:</div>
                                    <div>{{ project.tutor_name || 'N/D' }} ({{ project.tutor_email || 'N/D' }})</div>
                                </div>
                                <div class="col-12 col-md-6">
                                    <div class="text-weight-bold">Periodo:</div>
                                    <div>{{ new Date(project.start_date).toLocaleDateString() }} - {{ new Date(project.end_date).toLocaleDateString() }}</div>
                                </div>
                            </div>
                            
                            <q-separator class="q-my-md" v-if="project.description" />
                            
                            <div v-if="project.description">
                                <div class="text-h6 q-mb-sm">{{ t('udaPage.descriptionLabel') }}</div>
                                <p>{{ project.description }}</p>
                            </div>
                        </q-card-section>
                        <q-card-actions align="right">
                             <q-btn flat icon="edit" :label="t('common.edit')" color="primary" />
                             <q-btn flat icon="cloud_upload" :label="t('didacticsPage.uploadMaterial')" color="primary" />
                        </q-card-actions>
                    </q-card>
                </q-expansion-item>

                <q-item v-if="projects.length === 0">
                    <q-item-section class="text-center text-grey q-pa-xl">
                        {{ t('didacticsPage.noMaterials') }}
                    </q-item-section>
                </q-item>
            </q-list>
        </div>

        <!-- Documentation -->
        <div class="col-12 col-md-4">
             <q-card>
                <q-card-section>
                    <div class="text-h6">{{ t('documentsPage.title') }}</div>
                </q-card-section>
                <q-list separator>
                    <q-item clickable v-ripple>
                        <q-item-section avatar><q-icon name="description" color="grey" /></q-item-section>
                        <q-item-section>Convenzione Stage</q-item-section>
                        <q-item-section side><q-btn flat round icon="download" size="sm" /></q-item-section>
                    </q-item>
                    <q-item clickable v-ripple>
                        <q-item-section avatar><q-icon name="description" color="grey" /></q-item-section>
                        <q-item-section>Patto Formativo</q-item-section>
                        <q-item-section side><q-btn flat round icon="download" size="sm" /></q-item-section>
                    </q-item>
                    <q-item clickable v-ripple>
                        <q-item-section avatar><q-icon name="assignment" color="grey" /></q-item-section>
                        <q-item-section>Foglio Presenze</q-item-section>
                        <q-item-section side><q-btn flat round icon="print" size="sm" /></q-item-section>
                    </q-item>
                </q-list>
             </q-card>
        </div>
    </div>
  </q-page>
</template>
