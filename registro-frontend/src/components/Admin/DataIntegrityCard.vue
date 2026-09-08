<template>
  <q-card class="glass-card shadow-soft rounded-xl border border-slate-100 overflow-hidden">
    <q-card-section class="q-pa-lg">
      <div class="row items-center justify-between q-mb-md">
        <div class="row items-center q-gutter-sm">
          <q-avatar color="indigo-1" text-color="indigo" icon="fact_check" size="md" />
          <div>
            <div class="text-h5 text-weight-bold text-outfit text-slate-800">
              {{ t('admin.dataIntegrity.title') || 'Verifica Congruità Dati Scolastici' }}
            </div>
            <div class="text-caption text-slate-500">
              {{ t('admin.dataIntegrity.subtitle') || 'Diagnostica in tempo reale su studenti orfani, lezioni sovrapposte e anomalie nel registro' }}
            </div>
          </div>
        </div>

        <q-btn
          unelevated
          color="primary"
          icon="refresh"
          :label="t('admin.dataIntegrity.runCheck') || 'Esegui Diagnostica'"
          class="rounded-lg q-px-md shadow-soft"
          :loading="loading"
          @click="fetchReport"
        />
      </div>

      <!-- Loading State -->
      <div v-if="loading && !report" class="q-py-xl text-center">
        <q-spinner-dots color="primary" size="48px" />
        <div class="text-caption text-slate-400 q-mt-sm">
          {{ t('admin.dataIntegrity.analyzing') || 'Analisi dei dati in corso...' }}
        </div>
      </div>

      <!-- Content -->
      <div v-else-if="report" class="q-gutter-y-lg">
        <!-- Score & KPI Summary Row -->
        <div class="row q-col-gutter-md items-center">
          <div class="col-12 col-sm-4 text-center">
            <q-circular-progress
              show-value
              font-size="28px"
              :value="report.score"
              size="110px"
              :thickness="0.18"
              :color="getScoreColor(report.score)"
              track-color="slate-200"
              class="q-ma-sm text-weight-bolder text-outfit"
            >
              {{ report.score }}%
            </q-circular-progress>
            <div class="text-subtitle2 text-weight-bold q-mt-xs" :class="getScoreTextColor(report.score)">
              {{ getScoreLabel(report.score) }}
            </div>
            <div class="text-caption text-slate-400">
              {{ t('admin.dataIntegrity.lastRun') || 'Ultima verifica' }}: {{ formatTime(report.run_at) }}
            </div>
          </div>

          <div class="col-12 col-sm-8">
            <div class="row q-col-gutter-sm">
              <div class="col-4">
                <q-card flat bordered class="q-pa-sm text-center rounded-lg bg-slate-50">
                  <div class="text-caption text-slate-500 font-medium">Totale Anomalie</div>
                  <div class="text-h6 text-weight-bold text-slate-800">{{ report.total_issues }}</div>
                </q-card>
              </div>
              <div class="col-4">
                <q-card flat bordered class="q-pa-sm text-center rounded-lg bg-red-50 border-red-100">
                  <div class="text-caption text-red-600 font-medium">Gravità Alta</div>
                  <div class="text-h6 text-weight-bold text-red-700">{{ highCount }}</div>
                </q-card>
              </div>
              <div class="col-4">
                <q-card flat bordered class="q-pa-sm text-center rounded-lg bg-amber-50 border-amber-100">
                  <div class="text-caption text-amber-700 font-medium">Gravità Media</div>
                  <div class="text-h6 text-weight-bold text-amber-800">{{ mediumCount }}</div>
                </q-card>
              </div>
            </div>
          </div>
        </div>

        <!-- Checks Accordion List -->
        <q-list bordered separator class="rounded-xl overflow-hidden bg-white shadow-xs">
          <q-expansion-item
            v-for="chk in report.checks"
            :key="chk.id"
            group="integrityChecks"
            header-class="q-py-md"
            expand-icon="expand_more"
          >
            <template #header>
              <q-item-section avatar>
                <q-avatar
                  size="36px"
                  :color="chk.count === 0 ? 'emerald-1' : (chk.severity === 'high' ? 'red-1' : 'amber-1')"
                  :text-color="chk.count === 0 ? 'positive' : (chk.severity === 'high' ? 'negative' : 'warning')"
                  :icon="chk.count === 0 ? 'check_circle' : (chk.severity === 'high' ? 'error' : 'warning')"
                />
              </q-item-section>

              <q-item-section>
                <div class="row items-center q-gutter-x-sm">
                  <span class="text-subtitle2 text-weight-bold text-slate-800">{{ chk.title }}</span>
                  <q-badge
                    :color="chk.count === 0 ? 'positive' : (chk.severity === 'high' ? 'negative' : 'warning')"
                    class="q-px-xs rounded-borders text-caption font-bold"
                  >
                    {{ chk.count === 0 ? 'Conforme' : `${chk.count} anomalie` }}
                  </q-badge>
                </div>
                <div class="text-caption text-slate-500 q-mt-xs">{{ chk.description }}</div>
              </q-item-section>
            </template>

            <!-- Expanded Details -->
            <q-card class="bg-slate-50 border-t">
              <q-card-section class="q-pa-md">
                <div v-if="chk.count === 0" class="text-center text-positive q-py-sm">
                  <q-icon name="task_alt" size="24px" class="q-mr-xs" />
                  <span class="text-caption font-medium">Nessuna criticità riscontrata in questo controllo.</span>
                </div>

                <div v-else>
                  <q-markup-table flat dense class="bg-transparent">
                    <thead>
                      <tr class="text-left text-slate-500 text-caption font-semibold">
                        <th v-if="chk.id === 'orphaned_students'">Studente</th>
                        <th v-if="chk.id === 'orphaned_students'">Email</th>

                        <th v-if="chk.id === 'uncoordinated_classes'">Classe</th>

                        <th v-if="chk.id === 'overlapping_lessons'">Docente</th>
                        <th v-if="chk.id === 'overlapping_lessons'">Data & Ora</th>
                        <th v-if="chk.id === 'overlapping_lessons'">Classi in Conflitto</th>

                        <th v-if="chk.id === 'weekend_grades'">Studente</th>
                        <th v-if="chk.id === 'weekend_grades'">Materia</th>
                        <th v-if="chk.id === 'weekend_grades'">Data Registrata</th>
                        <th v-if="chk.id === 'weekend_grades'">Voto</th>

                        <th v-if="chk.id === 'unlinked_guardians'">Studente</th>
                        <th v-if="chk.id === 'unlinked_guardians'">Classe</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(item, idx) in chk.items" :key="idx" class="text-caption text-slate-700">
                        <!-- Orphaned students -->
                        <td v-if="chk.id === 'orphaned_students'" class="font-medium">{{ item.name }}</td>
                        <td v-if="chk.id === 'orphaned_students'">{{ item.email }}</td>

                        <!-- Uncoordinated classes -->
                        <td v-if="chk.id === 'uncoordinated_classes'" class="font-medium">{{ item.class }}</td>

                        <!-- Overlapping lessons -->
                        <td v-if="chk.id === 'overlapping_lessons'" class="font-medium">{{ item.teacher }}</td>
                        <td v-if="chk.id === 'overlapping_lessons'">{{ item.date }} ({{ item.hour }}ª ora)</td>
                        <td v-if="chk.id === 'overlapping_lessons'">{{ item.class1 }} ⚡ {{ item.class2 }}</td>

                        <!-- Weekend grades -->
                        <td v-if="chk.id === 'weekend_grades'" class="font-medium">{{ item.student }}</td>
                        <td v-if="chk.id === 'weekend_grades'">{{ item.subject }}</td>
                        <td v-if="chk.id === 'weekend_grades'">{{ item.date }}</td>
                        <td v-if="chk.id === 'weekend_grades'">{{ item.grade }}</td>

                        <!-- Unlinked guardians -->
                        <td v-if="chk.id === 'unlinked_guardians'" class="font-medium">{{ item.student }}</td>
                        <td v-if="chk.id === 'unlinked_guardians'">{{ item.class }}</td>
                      </tr>
                    </tbody>
                  </q-markup-table>
                </div>
              </q-card-section>
            </q-card>
          </q-expansion-item>
        </q-list>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import adminService from '@/services/adminService'

const { t } = useI18n()

const loading = ref(false)
const report = ref(null)

const highCount = computed(() => {
  if (!report.value?.checks) return 0
  return report.value.checks
    .filter(c => c.severity === 'high')
    .reduce((acc, c) => acc + c.count, 0)
})

const mediumCount = computed(() => {
  if (!report.value?.checks) return 0
  return report.value.checks
    .filter(c => c.severity === 'medium')
    .reduce((acc, c) => acc + c.count, 0)
})

function getScoreColor(score) {
  if (score >= 90) return 'positive'
  if (score >= 70) return 'warning'
  return 'negative'
}

function getScoreTextColor(score) {
  if (score >= 90) return 'text-positive'
  if (score >= 70) return 'text-warning'
  return 'text-negative'
}

function getScoreLabel(score) {
  if (score >= 90) return t('admin.dataIntegrity.scoreHealthy') || 'Dati Congrui'
  if (score >= 70) return t('admin.dataIntegrity.scoreWarning') || 'Attenzione Richiesta'
  return t('admin.dataIntegrity.scoreCritical') || 'Criticità Riscontrate'
}

function formatTime(timestamp) {
  if (!timestamp) return '-'
  const d = new Date(timestamp)
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

async function fetchReport() {
  loading.value = true
  try {
    const res = await adminService.getDataIntegrity()
    report.value = res.data || res
  } catch (err) {
    console.error('Failed to run data integrity checks', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchReport()
})
</script>

<style scoped>
.font-medium {
  font-weight: 500;
}
.font-bold {
  font-weight: 600;
}
</style>
