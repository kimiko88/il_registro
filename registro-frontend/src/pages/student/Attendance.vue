<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">{{ t('nav.myAttendance') || t('classRegister.title') || 'Presenze & Assenze' }}</div>
       <q-btn v-if="isParentUser" icon="fact_check" :label="t('attendancePage.requestJustification') || 'Richiedi Giustificazione'" color="primary" @click="showJustifyDialog = true" />
       <q-badge v-else color="grey-6" class="q-pa-xs">{{ t('attendancePage.parentManagedJustifications') || 'Giustificazioni gestite dai Genitori' }}</q-badge>
    </div>

    <div class="row q-col-gutter-lg">
        <!-- Stats Sidebar -->
        <div class="col-12 col-md-4">
            <q-card class="text-center q-pa-md">
                <div class="text-h6">{{ t('classRegister.statsSummary') || 'Riepilogo Presenze' }}</div>
                <div class="q-my-md relative-position flex flex-center">
                    <q-circular-progress
                      show-value
                      font-size="20px"
                      class="q-ma-md"
                      :value="attendancePercentage"
                      size="150px"
                      :thickness="0.2"
                      color="green"
                      track-color="red-1"
                    >
                        {{ attendancePercentage }}%<br><span class="text-caption text-grey">{{ t('classRegister.present') || 'Presente' }}</span>
                    </q-circular-progress>
                </div>
                <div class="row justify-around">
                    <div>
                        <div class="text-h5 text-red">{{ totalAbsences }}</div>
                        <div class="text-caption">{{ t('classRegister.absent') || 'Assenze' }}</div>
                    </div>
                    <div>
                        <div class="text-h5 text-orange">{{ totalDelays }}</div>
                        <div class="text-caption">{{ t('classRegister.late') || 'Ritardi' }}</div>
                    </div>
                </div>
            </q-card>

            <q-card v-if="attendancePercentage < 85" class="q-mt-md bg-orange-1">
                <q-card-section>
                    <div class="text-subtitle2"><q-icon name="warning" /> {{ t('common.warning') || 'Attenzione' }}</div>
                    <div class="text-caption">{{ t('attendancePage.warningAbsenceLimit') || 'Hai raggiunto il 15% di assenze consentite. Mettiti in regola per evitare problemi con l\'anno scolastico.' }}</div>
                </q-card-section>
            </q-card>
        </div>

        <!-- Attendance List -->
        <div class="col-12 col-md-8">
            <q-table
              :title="t('attendancePage.tableTitle') || 'Giornale di Classe'"
              :rows="attendanceEvents"
              :columns="columns"
              row-key="id"
              :pagination="{ rowsPerPage: 10 }"
            >
                <template v-slot:body-cell-type="props">
                    <q-td :props="props">
                        <q-chip :color="getTypeColor(props.value)" text-color="white" size="sm" dense>{{ props.value }}</q-chip>
                    </q-td>
                </template>
                <template v-slot:body-cell-justified="props">
                    <q-td :props="props">
                        <q-icon v-if="props.value" name="check_circle" color="green" role="img" aria-label="Giustificata" />
                        <q-icon v-else name="cancel" color="red" role="img" aria-label="Non giustificata" />
                    </q-td>
                </template>
            </q-table>
        </div>
    </div>

    <!-- Dialog Giustificazione -->
    <q-dialog v-model="showJustifyDialog">
      <q-card style="min-width: 400px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ t('attendancePage.requestJustification') || 'Giustifica Assenza' }}</div>
        </q-card-section>
        <q-card-section class="q-pa-md">
          <q-select
            v-model="justification.event"
            :options="unjustifiedAbsences"
            :label="t('attendancePage.selectAbsence') || 'Seleziona Assenza da Giustificare'"
            outlined
            dense
            class="q-mb-md"
          />
          <q-select
            v-model="justification.reason"
            :options="[
                { label: t('attendancePage.reasonHealth') || 'Motivi di Salute', value: 'Salute' },
                { label: t('attendancePage.reasonFamily') || 'Motivi Familiari', value: 'Famiglia' },
                { label: t('common.other') || 'Altro', value: 'Altro' }
            ]"
            :label="t('attendancePage.reason') || 'Motivazione'"
            outlined
            dense
            class="q-mb-md"
          />
          <q-input
            v-model="justification.notes"
            type="textarea"
            :label="t('common.optionalNotes') || 'Note aggiuntive (opzionale)'"
            outlined
            dense
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" v-close-popup />
          <q-btn color="primary" :label="t('common.confirm') || 'Invia Richiesta'" @click="submitJustification" />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { attendanceService } from '@/services/attendanceService'

const $q = useQuasar()
const { t } = useI18n()
const authStore = useAuthStore()
const showJustifyDialog = ref(false)
const attendanceEvents = ref([])

onMounted(() => {
    fetchAttendance()
})

const columns = computed(() => [
  { name: 'date', label: t('gradesPage.date') || 'Data', field: 'date', align: 'left', sortable: true },
  { name: 'type', label: t('common.status') || 'Stato', field: 'type', align: 'center' },
  { name: 'time', label: t('attendancePage.timeDetail') || 'Dettaglio Orario', field: 'time', align: 'center' },
  { name: 'justified', label: t('attendancePage.justified') || 'Giustificata', field: 'justified', align: 'center' },
  { name: 'notes', label: t('common.notes') || 'Annotazioni', field: 'notes', align: 'left' }
])

const fetchAttendance = async () => {
    try {
        const res = await attendanceService.getMyAttendance()
        if (res.data) {
             const list = Array.isArray(res.data) ? res.data : (res.data.records || [])
             attendanceEvents.value = list.map(e => {
                 let timeDetail = '-'
                 const s = String(e.status).toLowerCase()
                 if (s === 'late') {
                     timeDetail = e.entry_time ? `Ingresso: ${e.entry_time}` : 'Ritardo'
                 } else if (s === 'leftearly' || s === 'early') {
                     timeDetail = e.exit_time ? `Uscita: ${e.exit_time}` : 'Uscita anticipata'
                 }
                 return {
                     id: e.id,
                     date: e.date,
                     type: mapStatus(e.status),
                     time: timeDetail,
                     justified: e.is_justified,
                     notes: e.notes
                 }
             })
        }
    } catch (e) {
        console.error(e)
    }
}

function mapStatus(status) {
    const s = String(status).toLowerCase()
    if (s === 'absent') return 'Assenza'
    if (s === 'late') return 'Ritardo'
    if (s === 'leftearly' || s === 'early') return 'Uscita Anticipata'
    if (s === 'present') return 'Presente'
    return status
}

const isParentUser = computed(() => {
    const role = authStore.userRole || authStore.user?.role
    return role === 'parent'
})

const totalAbsences = computed(() => attendanceEvents.value.filter(e => e.type === 'Assenza').length)
const totalDelays = computed(() => attendanceEvents.value.filter(e => e.type === 'Ritardo').length)

const attendancePercentage = computed(() => {
    if (attendanceEvents.value.length === 0) return 100
    const presentsCount = attendanceEvents.value.filter(e => e.type === 'Presente' || e.type === 'Ritardo' || e.type === 'Uscita Anticipata').length
    return Math.round((presentsCount / attendanceEvents.value.length) * 100)
})

const unjustifiedAbsences = computed(() => 
    attendanceEvents.value
        .filter(e => !e.justified && e.type === 'Assenza')
        .map(e => ({ label: `${e.date} - ${e.type}`, value: e.id, date: e.date }))
)

const justification = ref({ event: null, reason: null, notes: '' })

const getTypeColor = (type) => {
    if (type === 'Assenza') return 'red';
    if (type === 'Ritardo') return 'orange';
    if (type === 'Uscita Anticipata') return 'blue';
    return 'grey';
}

const submitJustification = async () => {
    if (!isParentUser.value) {
        $q.notify({ type: 'warning', message: 'Le giustificazioni devono essere inviate dal genitore' })
        return
    }
    if (!justification.value.event) return
    const eventObj = justification.value.event
    try {
        await attendanceService.justify(eventObj.value, {
             start_date: eventObj.date || new Date().toISOString().split('T')[0],
             reason: justification.value.reason || 'Assenza',
             notes: justification.value.notes || ''
        })
        $q.notify({ type: 'positive', message: 'Richiesta inviata con successo' })
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore durante l\'invio della giustificazione' }) 
    }
}
</script>

<style scoped>
.border-dashed {
    border: 2px dashed #ccc;
    border-radius: 8px;
}
</style>
