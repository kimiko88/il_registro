<template>
  <q-page padding class="bg-slate-50">
    <!-- Page Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="swap_horiz" color="primary" class="q-mr-sm" />
          Sostituzioni Docenti
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Gestione e conferma delle ore di sostituzione per colleghi assenti
        </p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn flat round icon="refresh" color="primary" :loading="loading" @click="loadData" />
      </div>
    </div>

    <!-- Quick Stats Cards -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Sostituzioni Oggi</div>
              <div class="text-h4 text-weight-bold text-slate-800 q-mt-xs">{{ todaySubstitutions.length }}</div>
            </div>
            <q-avatar color="blue-1" text-color="blue-7" icon="today" size="48px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Mie Sostituzioni Totali</div>
              <div class="text-h4 text-weight-bold text-indigo-7 q-mt-xs">{{ mySubstitutions.length }}</div>
            </div>
            <q-avatar color="indigo-50" text-color="indigo-7" icon="assignment_ind" size="48px" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4">
        <q-card flat bordered class="rounded-xl bg-white shadow-soft">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-caption text-uppercase text-weight-bold text-slate-500">Da Confermare</div>
              <div class="text-h4 text-weight-bold text-amber-7 q-mt-xs">{{ pendingConfirmationCount }}</div>
            </div>
            <q-avatar color="amber-50" text-color="amber-7" icon="pending_actions" size="48px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Main Content Tabs -->
    <q-card flat bordered class="rounded-xl bg-white shadow-soft overflow-hidden">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-slate-600 bg-slate-100 border-b border-slate-200"
        active-color="primary"
        indicator-color="primary"
        align="left"
        no-caps
      >
        <q-tab name="my-today" icon="event_available" label="Oggi" />
        <q-tab name="my-all" icon="history" label="Le Mie Sostituzioni" />
        <q-tab name="school" icon="school" label="Quadro Istituto" />
      </q-tabs>

      <q-separator />

      <q-tab-panels v-model="activeTab" animated class="bg-white">
        <!-- TAB: TODAY -->
        <q-tab-panel name="my-today" class="q-pa-md">
          <div v-if="loading" class="text-center q-pa-xl">
            <q-spinner-dots color="primary" size="40px" />
          </div>

          <div v-else-if="todaySubstitutions.length === 0" class="text-center q-pa-xl text-slate-500">
            <q-icon name="check_circle_outline" size="60px" color="positive" class="q-mb-md" />
            <div class="text-h6">Nessuna sostituzione prevista per oggi</div>
            <div class="text-caption text-slate-400">Non hai ore di supplenza assegnate per la giornata odierna.</div>
          </div>

          <q-list v-else separator class="rounded-lg">
            <q-item v-for="sub in todaySubstitutions" :key="sub.id" class="q-py-md">
              <q-item-section avatar>
                <q-avatar color="primary" text-color="white" font-size="16px" class="text-weight-bold">
                  {{ sub.hour || sub.hour_index || '1' }}°
                </q-avatar>
              </q-item-section>

              <q-item-section>
                <q-item-label class="text-weight-bold text-slate-800">
                  Ora {{ sub.hour || sub.hour_index }} — Classe {{ sub.class_name || sub.class_id || 'N/D' }}
                </q-item-label>
                <q-item-label caption class="text-slate-600">
                  Docente sostituito: <strong class="text-slate-800">{{ sub.absent_teacher_name || sub.absent_teacher_id }}</strong>
                </q-item-label>
                <q-item-label caption class="text-slate-400" v-if="sub.notes">
                  Note: {{ sub.notes }}
                </q-item-label>
              </q-item-section>

              <q-item-section side>
                <div class="row items-center q-gutter-sm">
                  <q-badge :color="statusColor(sub.status)" class="q-px-sm q-py-xs text-weight-bold">
                    {{ statusLabel(sub.status) }}
                  </q-badge>
                  <q-btn
                    v-if="sub.status === 'assigned' || sub.status === 'pending'"
                    color="positive" size="sm" icon="check" label="Conferma Presenza"
                    :loading="confirmingId === sub.id"
                    @click="confirmSubstitution(sub.id)"
                  />
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-tab-panel>

        <!-- TAB: MY ALL -->
        <q-tab-panel name="my-all" class="q-pa-md">
          <q-table
            :rows="mySubstitutions"
            :columns="myColumns"
            row-key="id"
            flat
            :loading="loading"
            no-data-label="Nessuna sostituzione trovata nello storico"
          >
            <template #body-cell-date="{ row }">
              <q-td>{{ formatDate(row.date) }}</q-td>
            </template>

            <template #body-cell-hour="{ row }">
              <q-td>
                <q-badge color="blue-1" text-color="blue-9" class="text-weight-bold">
                  {{ row.hour || row.hour_index }}° ora
                </q-badge>
              </q-td>
            </template>

            <template #body-cell-status="{ row }">
              <q-td>
                <q-chip size="sm" :color="statusColor(row.status)" text-color="white" class="text-weight-bold">
                  {{ statusLabel(row.status) }}
                </q-chip>
              </q-td>
            </template>

            <template #body-cell-actions="{ row }">
              <q-td class="text-right">
                <q-btn
                  v-if="row.status === 'assigned' || row.status === 'pending'"
                  flat dense color="positive" icon="check_circle" label="Conferma"
                  :loading="confirmingId === row.id"
                  @click="confirmSubstitution(row.id)"
                />
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>

        <!-- TAB: SCHOOL WIDE -->
        <q-tab-panel name="school" class="q-pa-md">
          <div class="row items-center justify-between q-mb-md">
            <div class="text-subtitle1 text-weight-bold text-slate-800">Quadro Sostituzioni Istituto</div>
            <q-input
              v-model="filterDate"
              type="date"
              dense outlined
              label="Data"
              style="max-width: 180px"
              @update:model-value="fetchSchoolSubstitutions"
            />
          </div>

          <q-table
            :rows="schoolSubstitutions"
            :columns="schoolColumns"
            row-key="id"
            flat
            :loading="loadingSchool"
            no-data-label="Nessuna sostituzione registrata per la data selezionata"
          >
            <template #body-cell-date="{ row }">
              <q-td>{{ formatDate(row.date) }}</q-td>
            </template>

            <template #body-cell-status="{ row }">
              <q-td>
                <q-chip size="sm" :color="statusColor(row.status)" text-color="white" class="text-weight-bold">
                  {{ statusLabel(row.status) }}
                </q-chip>
              </q-td>
            </template>
          </q-table>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { substitutionService } from '@/services/substitutionService'
import api from 'src/services/api'

const $q = useQuasar()
const loading = ref(false)
const loadingSchool = ref(false)
const confirmingId = ref(null)
const activeTab = ref('my-today')

const todaySubstitutions = ref([])
const mySubstitutions = ref([])
const schoolSubstitutions = ref([])
const filterDate = ref(new Date().toISOString().substring(0, 10))

const pendingConfirmationCount = computed(() => {
  return mySubstitutions.value.filter(s => s.status === 'assigned' || s.status === 'pending').length
})

const myColumns = [
  { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'hour', label: 'Ora', field: 'hour', align: 'center', sortable: true },
  { name: 'class_name', label: 'Classe', field: row => row.class_name || row.class_id || 'N/D', align: 'left' },
  { name: 'absent_teacher', label: 'Docente Assente', field: row => row.absent_teacher_name || row.absent_teacher_id, align: 'left' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center' },
  { name: 'actions', label: 'Azioni', field: 'actions', align: 'right' }
]

const schoolColumns = [
  { name: 'date', label: 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'hour', label: 'Ora', field: row => `${row.hour || row.hour_index || 1}° ora`, align: 'center' },
  { name: 'class_name', label: 'Classe', field: row => row.class_name || row.class_id || 'N/D', align: 'left' },
  { name: 'absent_teacher', label: 'Docente Assente', field: row => row.absent_teacher_name || row.absent_teacher_id, align: 'left' },
  { name: 'substitute_teacher', label: 'Sostituto', field: row => row.substitute_teacher_name || row.substitute_teacher_id || 'Non assegnato', align: 'left' },
  { name: 'status', label: 'Stato', field: 'status', align: 'center' }
]

const formatDate = (d) => d ? qdate.formatDate(new Date(d), 'DD/MM/YYYY') : ''

const statusLabel = (s) => ({
  assigned: 'Assegnata',
  pending: 'In Attesa',
  confirmed: 'Confermata',
  completed: 'Completata',
  canceled: 'Annullata'
}[s] || s || 'Assegnata')

const statusColor = (s) => ({
  assigned: 'warning',
  pending: 'orange',
  confirmed: 'positive',
  completed: 'blue',
  canceled: 'grey'
}[s] || 'primary')

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const [todayRes, myRes] = await Promise.all([
      api.get('/substitutions/my-today').catch(() => ({ data: [] })),
      substitutionService.listMy().catch(() => ({ data: [] }))
    ])
    todaySubstitutions.value = todayRes.data || []
    mySubstitutions.value = myRes.data || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento delle sostituzioni' })
  } finally {
    loading.value = false
  }
  fetchSchoolSubstitutions()
}

async function fetchSchoolSubstitutions() {
  loadingSchool.value = true
  try {
    const res = await api.get('/substitutions', { params: { date: filterDate.value } })
    schoolSubstitutions.value = res.data || []
  } catch {
    schoolSubstitutions.value = []
  } finally {
    loadingSchool.value = false
  }
}

async function confirmSubstitution(id) {
  confirmingId.value = id
  try {
    await api.patch(`/substitutions/${id}/confirm`)
    $q.notify({ type: 'positive', message: 'Sostituzione confermata con successo' })
    await loadData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante la conferma' })
  } finally {
    confirmingId.value = null
  }
}
</script>
