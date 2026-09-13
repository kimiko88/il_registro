<template>
  <q-page class="q-pa-md q-pa-lg-xl visitor-registry-page">
    <!-- Hero Header -->
    <div class="row items-center justify-between q-mb-lg gap-md">
      <div class="col-12 col-md-7">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge color="amber-9" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="meeting_room" size="14px" class="q-mr-xs" />
            {{ t('visitorRegistry.badge') }}
          </q-badge>
          <q-badge outline color="primary" class="q-px-sm q-py-xs text-weight-bold">
            {{ selectedDateFormatted }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none flex items-center">
          <q-icon name="door_front" color="amber-9" class="q-mr-sm" size="36px" />
          {{ t('visitorRegistry.title') }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          {{ t('visitorRegistry.subtitle') }}
        </div>
      </div>

      <!-- Controls & Actions -->
      <div class="col-12 col-md-5 row items-center justify-end q-gutter-sm">
        <q-btn
          v-if="activeTab === 'visitors'"
          color="amber-9"
          icon="person_add"
          :label="t('visitorRegistry.newVisitorBtn')"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openNewVisitorDialog"
        />
        <q-btn
          v-else-if="activeTab === 'early_exits'"
          color="indigo-7"
          icon="exit_to_app"
          :label="t('visitorRegistry.newEarlyExitBtn')"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openNewEarlyExitDialog"
        />
        <q-btn
          v-else-if="activeTab === 'maintenance'"
          color="deep-orange-7"
          icon="handyman"
          :label="t('visitorRegistry.newReportBtn')"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openNewMaintenanceDialog"
        />
        <q-btn
          flat
          round
          dense
          color="primary"
          icon="refresh"
          :loading="loading"
          @click="loadCurrentTabData"
        >
          <q-tooltip>{{ t('common.refresh') }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <q-tabs
      v-model="activeTab"
      dense
      class="bg-white rounded-xl shadow-1 text-slate-600 q-mb-lg"
      active-color="primary"
      indicator-color="primary"
      align="left"
      narrow-indicator
      @update:model-value="onTabChange"
    >
      <q-tab name="visitors" icon="badge" :label="t('visitorRegistry.tabVisitors')" />
      <q-tab name="early_exits" icon="school" :label="t('visitorRegistry.tabEarlyExits')" />
      <q-tab name="maintenance" icon="build" :label="t('visitorRegistry.tabMaintenance')" />
    </q-tabs>

    <!-- Tab Panels -->
    <q-tab-panels v-model="activeTab" animated class="bg-transparent">
      <!-- PANEL 1: VISITATORI ESTERNI -->
      <q-tab-panel name="visitors" class="q-pa-none">
        <!-- Stats Summary -->
        <div class="row q-col-gutter-md q-mb-lg">
          <div class="col-12 col-sm-4">
            <q-card class="stat-card bg-emerald-50 border-emerald-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-emerald-100 text-emerald-800">
                  <q-icon name="location_on" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-emerald-800 text-weight-bold">{{ t('visitorRegistry.statPresent') }}</div>
                  <div class="text-h5 text-weight-bolder text-emerald-900">{{ currentVisitors.length }}</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 col-sm-4">
            <q-card class="stat-card bg-blue-50 border-blue-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-blue-100 text-blue-800">
                  <q-icon name="logout" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-blue-800 text-weight-bold">{{ t('visitorRegistry.statExited') }}</div>
                  <div class="text-h5 text-weight-bolder text-blue-900">{{ exitedVisitors.length }}</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
          <div class="col-12 col-sm-4">
            <q-card class="stat-card bg-amber-50 border-amber-200">
              <q-card-section class="row items-center no-wrap">
                <div class="stat-icon bg-amber-100 text-amber-800">
                  <q-icon name="groups" size="24px" />
                </div>
                <div class="q-ml-md">
                  <div class="text-caption text-amber-800 text-weight-bold">{{ t('visitorRegistry.statTotal') }}</div>
                  <div class="text-h5 text-weight-bolder text-amber-900">{{ visitorsList.length }}</div>
                </div>
              </q-card-section>
            </q-card>
          </div>
        </div>

        <!-- Table: Visitatori -->
        <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white">
          <q-table
            :rows="visitorsList"
            :columns="visitorColumns"
            row-key="id"
            :loading="loading"
            :no-data-label="t('visitorRegistry.noVisitorsToday')"
            class="no-shadow"
            :pagination="{ rowsPerPage: 15 }"
          >
            <template v-slot:body-cell-purpose="props">
              <q-td :props="props">
                <q-chip :color="getPurposeColor(props.row.purpose)" text-color="white" size="sm" class="text-weight-bold">
                  {{ getPurposeLabel(props.row.purpose) }}
                </q-chip>
              </q-td>
            </template>

            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge v-if="!props.row.exit_time" color="positive" class="q-px-sm q-py-xs">
                  <q-icon name="check_circle" size="12px" class="q-mr-xs" /> {{ t('visitorRegistry.statusOnSite') }}
                </q-badge>
                <q-badge v-else color="grey-6" class="q-px-sm q-py-xs">
                  {{ t('visitorRegistry.statusExitedAt', { time: formatTime(props.row.exit_time) }) }}
                </q-badge>
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props" align="right">
                <q-btn
                  v-if="!props.row.exit_time"
                  color="negative"
                  flat
                  dense
                  icon="logout"
                  :label="t('visitorRegistry.exitAction')"
                  no-caps
                  size="sm"
                  class="text-weight-bold"
                  @click="recordExit(props.row)"
                />
                <span v-else class="text-caption text-slate-400">—</span>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- PANEL 2: USCITE ANTICIPATE STUDENTI -->
      <q-tab-panel name="early_exits" class="q-pa-none">
        <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white">
          <q-table
            :rows="earlyExitsList"
            :columns="earlyExitColumns"
            row-key="id"
            :loading="loading"
            :no-data-label="t('visitorRegistry.noEarlyExitsToday')"
            class="no-shadow"
            :pagination="{ rowsPerPage: 15 }"
          >
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-badge v-if="!props.row.return_time" color="warning" text-color="dark" class="q-px-sm q-py-xs">
                  <q-icon name="directions_walk" size="12px" class="q-mr-xs" /> {{ t('visitorRegistry.statusOffSite') }}
                </q-badge>
                <q-badge v-else color="positive" class="q-px-sm q-py-xs">
                  {{ t('visitorRegistry.statusReturnedAt', { time: formatTime(props.row.return_time) }) }}
                </q-badge>
              </q-td>
            </template>

            <template v-slot:body-cell-actions="props">
              <q-td :props="props" align="right">
                <q-btn
                  v-if="!props.row.return_time"
                  color="primary"
                  flat
                  dense
                  icon="keyboard_return"
                  :label="t('visitorRegistry.returnAction')"
                  no-caps
                  size="sm"
                  class="text-weight-bold"
                  @click="recordReturn(props.row)"
                />
                <span v-else class="text-caption text-slate-400">—</span>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </q-tab-panel>

      <!-- PANEL 3: SEGNALAZIONI GUASTI -->
      <q-tab-panel name="maintenance" class="q-pa-none">
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-sm-6 col-md-3">
            <q-select
              v-model="maintenanceStatusFilter"
              :options="maintenanceFilterOptions"
              emit-value
              map-options
              outlined
              dense
              bg-color="white"
              @update:model-value="loadMaintenance"
            />
          </div>
        </div>

        <div class="row q-col-gutter-md">
          <div v-for="rep in maintenanceList" :key="rep.id" class="col-12 col-md-6 col-lg-4">
            <q-card class="rounded-2xl shadow-sm border border-slate-200 bg-white full-height column justify-between">
              <q-card-section>
                <div class="row items-center justify-between q-mb-sm">
                  <q-badge :color="getPriorityColor(rep.priority)" class="text-uppercase text-weight-bold q-px-sm">
                    {{ rep.priority }}
                  </q-badge>
                  <q-badge :color="getMaintenanceStatusColor(rep.status)">
                    {{ getMaintenanceStatusLabel(rep.status) }}
                  </q-badge>
                </div>

                <div class="text-h6 text-weight-bold text-slate-900 q-mb-xs flex items-center">
                  <q-icon name="place" color="primary" size="20px" class="q-mr-xs" />
                  {{ rep.location }}
                </div>
                <div class="text-caption text-slate-500 q-mb-sm text-capitalize">
                  {{ t('visitorRegistry.categoryPrefix') }} <strong>{{ rep.category }}</strong>
                </div>
                <div class="text-body2 text-slate-700 q-mb-md">
                  {{ rep.description }}
                </div>
              </q-card-section>

              <q-card-section class="border-t border-slate-100 row items-center justify-between q-pt-sm">
                <div class="text-caption text-slate-400">
                  {{ formatDate(rep.created_at) }}
                </div>

                <!-- Status updater -->
                <q-btn-dropdown
                  flat
                  dense
                  no-caps
                  color="primary"
                  :label="t('visitorRegistry.updateStatus')"
                  size="sm"
                >
                  <q-list dense>
                    <q-item clickable v-close-popup @click="updateMaintStatus(rep, 'aperto')">
                      <q-item-section>{{ t('visitorRegistry.maintStatusOpen') }}</q-item-section>
                    </q-item>
                    <q-item clickable v-close-popup @click="updateMaintStatus(rep, 'in_lavorazione')">
                      <q-item-section>{{ t('visitorRegistry.maintStatusInProgress') }}</q-item-section>
                    </q-item>
                    <q-item clickable v-close-popup @click="updateMaintStatus(rep, 'chiuso')">
                      <q-item-section>{{ t('visitorRegistry.maintStatusClosed') }}</q-item-section>
                    </q-item>
                  </q-list>
                </q-btn-dropdown>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog: Nuovo Visitatore -->
    <q-dialog v-model="visitorDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="badge" color="amber-9" class="q-mr-sm" size="24px" />
            {{ t('visitorRegistry.dialogVisitorTitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.nameLabel') }}</label>
            <q-input v-model="visitorForm.name" outlined dense placeholder="Mario Rossi" />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.docLabel') }}</label>
              <q-input v-model="visitorForm.document_id" outlined dense :placeholder="t('visitorRegistry.docPlaceholder')" />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.badgeLabel') }}</label>
              <q-input v-model="visitorForm.badge_number" outlined dense placeholder="es. 04" />
            </div>
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.purposeLabel') }}</label>
              <q-select
                v-model="visitorForm.purpose"
                :options="purposeOptions"
                emit-value
                map-options
                outlined
                dense
              />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.hostLabel') }}</label>
              <q-input v-model="visitorForm.host_name" outlined dense :placeholder="t('visitorRegistry.hostPlaceholder')" />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.notesLabel') }}</label>
            <q-input v-model="visitorForm.notes" outlined dense type="textarea" rows="2" />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel')" color="grey-7" v-close-popup />
          <q-btn
            color="amber-9"
            :label="t('visitorRegistry.saveVisitor')"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="saving"
            @click="submitVisitor"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Nuova Uscita Anticipata -->
    <q-dialog v-model="earlyExitDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="school" color="indigo-7" class="q-mr-sm" size="24px" />
            {{ t('visitorRegistry.dialogEarlyExitTitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.studentLabel') }}</label>
            <q-select
              v-model="earlyExitForm.student_id"
              :options="studentOptions"
              option-value="id"
              option-label="name"
              emit-value
              map-options
              outlined
              dense
              use-input
              @filter="filterStudents"
              :placeholder="t('visitorRegistry.searchStudentPlaceholder')"
            />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-7">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.delegateeLabel') }}</label>
              <q-input v-model="earlyExitForm.delegatee_name" outlined dense :placeholder="t('visitorRegistry.delegateePlaceholder')" />
            </div>
            <div class="col-5">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.relLabel') }}</label>
              <q-input v-model="earlyExitForm.delegate_rel" outlined dense :placeholder="t('visitorRegistry.relPlaceholder')" />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.reasonLabel') }}</label>
            <q-select
              v-model="earlyExitForm.reason_code"
              :options="reasonOptions"
              emit-value
              map-options
              outlined
              dense
            />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.notesLabel') }}</label>
            <q-input v-model="earlyExitForm.notes" outlined dense type="textarea" rows="2" />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel')" color="grey-7" v-close-popup />
          <q-btn
            color="indigo-7"
            :label="t('visitorRegistry.saveEarlyExit')"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="saving"
            @click="submitEarlyExit"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Segnalazione Manutenzione -->
    <q-dialog v-model="maintenanceDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-2xl">
        <q-card-section class="row items-center justify-between border-b border-slate-100">
          <div class="text-h6 text-weight-bold text-slate-900 flex items-center">
            <q-icon name="build" color="deep-orange-7" class="q-mr-sm" size="24px" />
            {{ t('visitorRegistry.dialogMaintTitle') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.locationLabel') }}</label>
            <q-input v-model="maintenanceForm.location" outlined dense :placeholder="t('visitorRegistry.locationPlaceholder')" />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.categoryLabel') }}</label>
              <q-select
                v-model="maintenanceForm.category"
                :options="categoryOptions"
                emit-value
                map-options
                outlined
                dense
              />
            </div>
            <div class="col-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.priorityLabel') }}</label>
              <q-select
                v-model="maintenanceForm.priority"
                :options="priorityOptions"
                emit-value
                map-options
                outlined
                dense
              />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('visitorRegistry.descLabel') }}</label>
            <q-input
              v-model="maintenanceForm.description"
              outlined
              dense
              type="textarea"
              rows="3"
              :placeholder="t('visitorRegistry.descPlaceholder')"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('common.cancel')" color="grey-7" v-close-popup />
          <q-btn
            color="deep-orange-7"
            :label="t('visitorRegistry.saveMaint')"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="saving"
            @click="submitMaintenance"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import visitorService from '@/services/visitorService'
import userService from '@/services/userService'

const $q = useQuasar()
const { t, locale } = useI18n()

const activeTab = ref('visitors')
const loading = ref(false)
const saving = ref(false)

const todayStr = new Date().toISOString().substring(0, 10)
const selectedDateFormatted = computed(() => {
  const d = new Date()
  return d.toLocaleDateString(locale.value || 'it-IT', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
})

// Tab 1: Visitors
const visitorsList = ref([])
const visitorDialog = ref(false)
const visitorForm = ref({
  name: '',
  document_id: '',
  purpose: 'parent',
  host_name: '',
  badge_number: '',
  notes: ''
})

const purposeOptions = computed(() => [
  { label: t('visitorRegistry.purposeParent'), value: 'parent' },
  { label: t('visitorRegistry.purposeSupplier'), value: 'supplier' },
  { label: t('visitorRegistry.purposeInstitution'), value: 'institution' },
  { label: t('visitorRegistry.purposeOther'), value: 'other' }
])

const currentVisitors = computed(() => visitorsList.value.filter(v => !v.exit_time))
const exitedVisitors = computed(() => visitorsList.value.filter(v => v.exit_time))

const visitorColumns = computed(() => [
  { name: 'name', label: t('visitorRegistry.colVisitorName'), field: 'name', align: 'left', sortable: true },
  { name: 'document_id', label: t('visitorRegistry.colDocument'), field: 'document_id', align: 'left' },
  { name: 'purpose', label: t('visitorRegistry.colPurpose'), field: 'purpose', align: 'left' },
  { name: 'host_name', label: t('visitorRegistry.colHost'), field: 'host_name', align: 'left' },
  { name: 'badge_number', label: t('visitorRegistry.colBadge'), field: 'badge_number', align: 'center' },
  { name: 'entry_time', label: t('visitorRegistry.colEntry'), field: row => formatTime(row.entry_time), align: 'center' },
  { name: 'status', label: t('common.status'), align: 'center' },
  { name: 'actions', label: t('common.actions'), align: 'right' }
])

// Tab 2: Early Exits
const earlyExitsList = ref([])
const earlyExitDialog = ref(false)
const earlyExitForm = ref({
  student_id: '',
  delegatee_name: '',
  delegate_rel: '',
  reason_code: 'visita_medica',
  notes: ''
})
const studentList = ref([])
const studentOptions = ref([])

const earlyExitColumns = computed(() => [
  { name: 'student_name', label: t('visitorRegistry.colStudent'), field: 'student_name', align: 'left', sortable: true },
  { name: 'class_name', label: t('visitorRegistry.colClass'), field: 'class_name', align: 'center' },
  { name: 'delegatee_name', label: t('visitorRegistry.colDelegatee'), field: 'delegatee_name', align: 'left' },
  { name: 'delegate_rel', label: t('visitorRegistry.colRel'), field: 'delegate_rel', align: 'left' },
  { name: 'exit_time', label: t('visitorRegistry.colExitTime'), field: row => formatTime(row.exit_time), align: 'center' },
  { name: 'status', label: t('common.status'), align: 'center' },
  { name: 'actions', label: t('common.actions'), align: 'right' }
])

const reasonOptions = computed(() => [
  { label: t('visitorRegistry.reasonMedical'), value: 'visita_medica' },
  { label: t('visitorRegistry.reasonFamily'), value: 'motivi_familiari' },
  { label: t('visitorRegistry.reasonIllness'), value: 'malessere' },
  { label: t('visitorRegistry.reasonOther'), value: 'altro' }
])

// Tab 3: Maintenance
const maintenanceList = ref([])
const maintenanceDialog = ref(false)
const maintenanceStatusFilter = ref('')
const maintenanceForm = ref({
  location: '',
  category: 'elettrico',
  priority: 'media',
  description: ''
})

const maintenanceFilterOptions = computed(() => [
  { label: t('visitorRegistry.filterAllReports'), value: '' },
  { label: t('visitorRegistry.maintStatusOpen'), value: 'aperto' },
  { label: t('visitorRegistry.maintStatusInProgress'), value: 'in_lavorazione' },
  { label: t('visitorRegistry.maintStatusClosed'), value: 'chiuso' }
])

const categoryOptions = computed(() => [
  { label: t('visitorRegistry.catElectrical'), value: 'elettrico' },
  { label: t('visitorRegistry.catPlumbing'), value: 'idraulico' },
  { label: t('visitorRegistry.catStructural'), value: 'strutturale' },
  { label: t('visitorRegistry.catCleaning'), value: 'pulizia' },
  { label: t('visitorRegistry.catItLim'), value: 'informatica' },
  { label: t('visitorRegistry.catOther'), value: 'altro' }
])

const priorityOptions = computed(() => [
  { label: t('visitorRegistry.prioLow'), value: 'bassa' },
  { label: t('visitorRegistry.prioMedium'), value: 'media' },
  { label: t('visitorRegistry.prioHigh'), value: 'alta' },
  { label: t('visitorRegistry.prioUrgent'), value: 'urgente' }
])

function onTabChange(_val) {
  loadCurrentTabData()
}

async function loadCurrentTabData() {
  if (activeTab.value === 'visitors') {
    await loadVisitors()
  } else if (activeTab.value === 'early_exits') {
    await loadEarlyExits()
  } else if (activeTab.value === 'maintenance') {
    await loadMaintenance()
  }
}

async function loadVisitors() {
  loading.value = true
  try {
    const res = await visitorService.listTodayVisitors(todayStr)
    visitorsList.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: t('visitorRegistry.notifyLoadVisitorsError') })
  } finally {
    loading.value = false
  }
}

async function loadEarlyExits() {
  loading.value = true
  try {
    const res = await visitorService.listTodayEarlyExits(todayStr)
    earlyExitsList.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: t('visitorRegistry.notifyLoadEarlyExitsError') })
  } finally {
    loading.value = false
  }
}

async function loadMaintenance() {
  loading.value = true
  try {
    const res = await visitorService.listMaintenanceReports(maintenanceStatusFilter.value)
    maintenanceList.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: t('visitorRegistry.notifyLoadReportsError') })
  } finally {
    loading.value = false
  }
}

function openNewVisitorDialog() {
  visitorForm.value = {
    name: '',
    document_id: '',
    purpose: 'parent',
    host_name: '',
    badge_number: '',
    notes: ''
  }
  visitorDialog.value = true
}

async function submitVisitor() {
  if (!visitorForm.value.name) {
    $q.notify({ type: 'warning', message: t('visitorRegistry.notifyNameRequired') })
    return
  }
  saving.value = true
  try {
    await visitorService.registerVisitor(visitorForm.value)
    $q.notify({ type: 'positive', message: t('visitorRegistry.notifyVisitorSaved') })
    visitorDialog.value = false
    await loadVisitors()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('visitorRegistry.notifyLoadVisitorsError') })
  } finally {
    saving.value = false
  }
}

async function recordExit(v) {
  try {
    await visitorService.recordVisitorExit(v.id, 'Uscita registrata da portineria')
    $q.notify({ type: 'positive', message: t('visitorRegistry.notifyExitRecorded', { name: v.name }) })
    await loadVisitors()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('visitorRegistry.notifyExitError') })
  }
}

function openNewEarlyExitDialog() {
  earlyExitForm.value = {
    student_id: '',
    delegatee_name: '',
    delegate_rel: '',
    reason_code: 'visita_medica',
    notes: ''
  }
  earlyExitDialog.value = true
}

async function submitEarlyExit() {
  if (!earlyExitForm.value.student_id || !earlyExitForm.value.delegatee_name) {
    $q.notify({ type: 'warning', message: t('visitorRegistry.notifyEarlyExitValidation') })
    return
  }
  saving.value = true
  try {
    await visitorService.recordEarlyExit(earlyExitForm.value)
    $q.notify({ type: 'positive', message: t('visitorRegistry.notifyEarlyExitSaved') })
    earlyExitDialog.value = false
    await loadEarlyExits()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('visitorRegistry.notifyEarlyExitError') })
  } finally {
    saving.value = false
  }
}

async function recordReturn(exit) {
  try {
    await visitorService.recordStudentReturn(exit.id, 'Rientro in aula')
    $q.notify({ type: 'positive', message: t('visitorRegistry.notifyReturnSaved') })
    await loadEarlyExits()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('visitorRegistry.notifyReturnError') })
  }
}

function openNewMaintenanceDialog() {
  maintenanceForm.value = {
    location: '',
    category: 'elettrico',
    priority: 'media',
    description: ''
  }
  maintenanceDialog.value = true
}

async function submitMaintenance() {
  if (!maintenanceForm.value.location || !maintenanceForm.value.description) {
    $q.notify({ type: 'warning', message: t('visitorRegistry.notifyMaintValidation') })
    return
  }
  saving.value = true
  try {
    await visitorService.createMaintenanceReport(maintenanceForm.value)
    $q.notify({ type: 'positive', message: t('visitorRegistry.notifyMaintSaved') })
    maintenanceDialog.value = false
    await loadMaintenance()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('visitorRegistry.notifyMaintError') })
  } finally {
    saving.value = false
  }
}

async function updateMaintStatus(rep, status) {
  try {
    await visitorService.updateMaintenanceStatus(rep.id, { status })
    $q.notify({ type: 'positive', message: t('visitorRegistry.notifyStatusUpdated') })
    await loadMaintenance()
  } catch (err) {
    $q.notify({ type: 'negative', message: t('visitorRegistry.notifyStatusError') })
  }
}

function filterStudents(val, update) {
  if (val === '') {
    update(() => {
      studentOptions.value = studentList.value
    })
    return
  }
  update(() => {
    const needle = val.toLowerCase()
    studentOptions.value = studentList.value.filter(s => s.name.toLowerCase().includes(needle))
  })
}

function formatTime(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleTimeString(locale.value || 'it-IT', { hour: '2-digit', minute: '2-digit' })
}

function formatDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleDateString(locale.value || 'it-IT', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}

function getPurposeLabel(p) {
  switch (p) {
    case 'parent': return t('visitorRegistry.purposeLabelParent')
    case 'supplier': return t('visitorRegistry.purposeLabelSupplier')
    case 'institution': return t('visitorRegistry.purposeLabelInstitution')
    default: return t('visitorRegistry.purposeOther')
  }
}

function getPurposeColor(p) {
  switch (p) {
    case 'parent': return 'blue-7'
    case 'supplier': return 'amber-9'
    case 'institution': return 'purple-7'
    default: return 'grey-7'
  }
}

function getPriorityColor(p) {
  switch (p) {
    case 'urgente': return 'negative'
    case 'alta': return 'deep-orange-7'
    case 'media': return 'amber-8'
    case 'bassa': return 'blue-grey-6'
    default: return 'grey-6'
  }
}

function getMaintenanceStatusColor(st) {
  switch (st) {
    case 'aperto': return 'negative'
    case 'in_lavorazione': return 'warning'
    case 'chiuso': return 'positive'
    default: return 'grey-6'
  }
}

function getMaintenanceStatusLabel(st) {
  switch (st) {
    case 'aperto': return t('visitorRegistry.statusOpen')
    case 'in_lavorazione': return t('visitorRegistry.statusInProgress')
    case 'chiuso': return t('visitorRegistry.statusResolved')
    default: return st
  }
}

onMounted(async () => {
  await loadVisitors()
  try {
    const sRes = await userService.getUsers({ role: 'student' })
    const raw = sRes.data?.users || sRes.data || []
    studentList.value = raw.map(u => ({
      id: u.id,
      name: `${u.last_name || ''} ${u.first_name || ''}`.trim() || u.email
    }))
    studentOptions.value = studentList.value
  } catch {
    // Non blocking
  }
})
</script>

<style scoped>
.visitor-registry-page {
  background: #f8fafc;
  min-height: 100vh;
}

.stat-card {
  border-radius: 16px;
  border-width: 1px;
  border-style: solid;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
