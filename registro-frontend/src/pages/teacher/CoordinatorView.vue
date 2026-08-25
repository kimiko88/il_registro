<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-primary q-my-none">{{ t('roleDashboards.teacherPanel') || 'Dashboard Coordinatore di Classe' }}</h1>
        <div class="text-subtitle2 text-grey-7">Panoramica andamento, note disciplinari, genitori e gestione scrutini</div>
      </div>
      <q-space />
      <q-select
        v-model="selectedClassId"
        :options="classOptions"
        option-value="id"
        option-label="label"
        emit-value
        map-options
        :label="t('udaPage.classLabel') || 'Seleziona Classe Coordinata'"
        outlined
        dense
        style="min-width: 250px"
        @update:model-value="onClassChange"
      />
    </div>

    <div v-if="!selectedClassId" class="q-pa-xl text-center">
      <q-icon name="co_present" size="4rem" color="grey-5" />
      <div class="text-h6 text-grey-6 q-mt-md">Non risulti coordinatore di alcuna classe per l'anno scolastico in corso.</div>
    </div>

    <div v-else>
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="left"
        narrow-indicator
      >
        <q-tab name="academic" icon="trending_up" :label="t('gradesPage.title') || 'Andamento Classe'" />
        <q-tab name="notes" icon="report_problem" :label="t('notesPage.title') || 'Note & Richiami Disciplinari'" />
        <q-tab name="guardians" icon="contacts" :label="t('nav.children') || 'Anagrafica Genitori'" />
        <q-tab name="scrutiny" icon="gavel" label="Scrutinio di Classe" />
      </q-tabs>

      <q-separator class="q-mb-md" />

      <q-tab-panels v-model="activeTab" animated class="bg-transparent">
        <!-- Tab 1: Andamento Classe -->
        <q-tab-panel name="academic" class="q-pa-none">
          <q-card class="shadow-2 rounded-borders q-pa-md">
            <div class="text-h6 text-weight-bold q-mb-md">Matrice Andamento Accademico</div>
            <q-spinner v-if="loadingMatrix" color="primary" size="2em" />
            <div v-else-if="matrixData">
              <q-markup-table flat bordered dense class="rounded-borders">
                <thead>
                  <tr class="bg-primary text-white">
                    <th class="text-left">{{ t('competenciesPage.student') }}</th>
                    <th v-for="sub in matrixData.subjects" :key="sub.id" class="text-center">{{ sub.name }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="std in matrixData.students" :key="std.id">
                    <td class="text-weight-medium">{{ std.name }}</td>
                    <td v-for="sub in matrixData.subjects" :key="sub.id" class="text-center">
                      <q-chip
                        dense
                        size="sm"
                        :color="getAverageColor(std.averages[sub.id])"
                        text-color="white"
                      >
                        {{ std.averages[sub.id] ? std.averages[sub.id].toFixed(1) : '-' }}
                      </q-chip>
                    </td>
                  </tr>
                </tbody>
              </q-markup-table>
            </div>
          </q-card>
        </q-tab-panel>

        <!-- Tab 2: Note Disciplinari -->
        <q-tab-panel name="notes" class="q-pa-none">
          <q-card class="shadow-2 rounded-borders q-pa-md">
            <div class="text-h6 text-weight-bold q-mb-md">{{ t('notesPage.title') }}</div>
            <q-list separator v-if="notesList.length > 0">
              <q-item v-for="note in notesList" :key="note.id">
                <q-item-section avatar>
                  <q-icon name="warning" color="warning" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ note.student_name }} - {{ note.type }}</q-item-label>
                  <q-item-label caption>{{ note.description }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="text-caption">{{ note.date }}</div>
                </q-item-section>
              </q-item>
            </q-list>
            <div v-else class="text-grey text-center q-pa-md">Nessuna nota presente per la classe selezionata.</div>
          </q-card>
        </q-tab-panel>

        <!-- Tab 3: Genitori -->
        <q-tab-panel name="guardians" class="q-pa-none">
          <q-card class="shadow-2 rounded-borders q-pa-md">
            <div class="text-h6 text-weight-bold q-mb-md">Contatti Genitori e Rappresentanti</div>
            <q-list separator v-if="guardiansList.length > 0">
              <q-item v-for="g in guardiansList" :key="g.id">
                <q-item-section avatar>
                  <q-avatar color="secondary" text-color="white" icon="person" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ g.name }} (Genitore di {{ g.student_name }})</q-item-label>
                  <q-item-label caption>Email: {{ g.email }} | Tel: {{ g.phone || 'N/D' }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-btn flat round icon="email" color="primary" :to="`/messages?to=${g.email}`" />
                </q-item-section>
              </q-item>
            </q-list>
            <div v-else class="text-grey text-center q-pa-md">Nessun contatto genitore trovato.</div>
          </q-card>
        </q-tab-panel>

        <!-- Tab 4: Scrutinio -->
        <q-tab-panel name="scrutiny" class="q-pa-none">
          <q-card class="shadow-2 rounded-borders q-pa-md">
            <div class="text-h6 text-weight-bold q-mb-md">Gestione Scrutinio Finale / Intermedio</div>
            <p class="text-body2 text-grey-8">Avvia e coordina la sessione di scrutinio per la classe {{ selectedClassId }}.</p>
            <q-btn color="primary" icon="gavel" label="Avvia Sessione Scrutinio" @click="startScrutiny" class="q-mt-sm" />
          </q-card>
        </q-tab-panel>
      </q-tab-panels>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import api from '@/services/api'
import authService from '@/services/authService'
import { scrutinyService } from '@/services/scrutinyService'
import notesService from '@/services/notesService'
import { useSchoolYearStore } from '@/stores/schoolYear'

const $q = useQuasar()
const { t } = useI18n()
const schoolYearStore = useSchoolYearStore()
const activeTab = ref('academic')
const selectedClassId = ref(null)
const classOptions = ref([])

const loadingMatrix = ref(false)
const matrixData = ref(null)
const notesList = ref([])
const guardiansList = ref([])

const fetchClasses = async () => {
  try {
    const userRes = await authService.getCurrentUser()
    const currentUserId = userRes?.id
    const res = await api.get('/teacher/classes', {
      params: { school_year: schoolYearStore.selectedSchoolYear }
    })
    const assignedClasses = res.data || []
    const coordClasses = assignedClasses.filter(c => c.coordinator_id === currentUserId)
    classOptions.value = coordClasses.length > 0 ? coordClasses : assignedClasses
    if (classOptions.value.length > 0) {
      selectedClassId.value = classOptions.value[0].id
      onClassChange(selectedClassId.value)
    } else {
      selectedClassId.value = null
    }
  } catch (err) {
    classOptions.value = []
  }
}

watch(() => schoolYearStore.selectedSchoolYear, () => {
  fetchClasses()
})

const onClassChange = async (classId) => {
  if (!classId) return
  fetchMatrix(classId)
  fetchNotes(classId)
  fetchGuardians(classId)
}

const fetchMatrix = async (classId) => {
  loadingMatrix.value = true
  try {
    const res = await scrutinyService.getMatrix(classId, 1)
    matrixData.value = res.data
  } catch (err) {
    matrixData.value = null
  } finally {
    loadingMatrix.value = false
  }
}

const fetchNotes = async (classId) => {
  try {
    const res = await notesService.getNotes({ class_id: classId })
    notesList.value = res.data || []
  } catch (err) {
    notesList.value = []
  }
}

const fetchGuardians = async (classId) => {
  try {
    const res = await api.get(`/classes/${classId}/guardians`)
    guardiansList.value = res.data || []
  } catch (err) {
    guardiansList.value = []
  }
}

const startScrutiny = async () => {
  try {
    await scrutinyService.startScrutiny(selectedClassId.value, 1)
    $q.notify({ type: 'positive', message: t('common.success') })
  } catch (err) {
    $q.notify({ type: 'negative', message: t('common.error') })
  }
}

const getAverageColor = (avg) => {
  if (!avg) return 'grey'
  if (avg >= 6.0) return 'positive'
  return 'negative'
}

onMounted(() => {
  fetchClasses()
})
</script>
