<template>
  <q-page padding>
    <div class="text-h4 text-weight-bold q-mb-lg">Centro Reportistica</div>

    <div class="row q-col-gutter-lg">
      <div class="col-12 col-md-6">
        <div class="text-h6 q-mb-sm text-primary">Didattica & Voti</div>
        <div class="row q-gutter-md">
          <ReportCard title="Pagelle / Scrutini" icon="grading" description="Genera pagelle di fine quadrimestre o tabelloni scrutini." @click="openReport('grades')" />
          <ReportCard title="Registro Attività" icon="menu_book" description="Riepilogo argomenti svolti per classe." @click="openReport('topics')" />
        </div>
      </div>
      <div class="col-12 col-md-6">
        <div class="text-h6 q-mb-sm text-secondary">Presenze & Note</div>
        <div class="row q-gutter-md">
          <ReportCard title="Assenze Mensili" icon="event_busy" description="Statistiche assenze per classe o studente." @click="openReport('attendance')" />
          <ReportCard title="Note Disciplinari" icon="warning" description="Elenco sanzioni e note disciplinari." @click="openReport('notes')" />
        </div>
      </div>
      <div class="col-12 col-md-6">
        <div class="text-h6 q-mb-sm text-accent">Amministrazione</div>
        <div class="row q-gutter-md">
          <ReportCard title="Elenco Iscritti" icon="people_alt" description="Lista completa studenti per classe con dati anagrafici." @click="openReport('students')" />
          <ReportCard title="Certificati" icon="verified" description="Certificati di frequenza e iscrizione." @click="openReport('certificates')" />
        </div>
      </div>
    </div>

    <!-- Parameter Dialog -->
    <q-dialog v-model="showDialog">
      <q-card style="min-width: 400px" class="rounded-xl overflow-hidden shadow-24 bg-white">
        <q-card-section class="bg-gradient-premium text-white q-pa-lg">
          <div class="text-h6 text-weight-bold">{{ currentReportTitle }}</div>
          <div class="text-caption opacity-80">Seleziona i parametri per generare il report.</div>
        </q-card-section>
        <q-card-section class="q-pa-lg q-gutter-y-md">
          <q-select v-model="info.class_id" :options="classOptions" option-label="label" option-value="value" emit-value map-options label="Classe" outlined dense />
          <q-select v-if="requiresPeriod" v-model="info.period" :options="['I Quadrimestre', 'II Quadrimestre']" label="Periodo" outlined dense />
          <q-select v-if="requiresMonth" v-model="info.month" :options="monthOptions" label="Mese" outlined dense />
        </q-card-section>
        <q-card-actions align="right" class="q-pa-md border-t border-slate-100 bg-slate-50">
          <q-btn flat label="Annulla" color="slate-500" v-close-popup no-caps />
          <q-btn outline color="primary" label="Anteprima" icon="visibility" class="rounded-lg" no-caps @click="showPreviewDialog" :loading="reportLoading" />
          <q-btn color="primary" label="Genera PDF" icon="picture_as_pdf" class="rounded-lg shadow-sm" no-caps @click="generatePDF" :loading="reportLoading" />
          <q-btn color="secondary" label="Excel" icon="table_view" class="rounded-lg shadow-sm" no-caps @click="generateXLS" :loading="reportLoading" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Preview Dialog -->
    <q-dialog v-model="showPreview" full-width full-height>
      <q-card class="column no-wrap bg-slate-100 rounded-xl overflow-hidden shadow-24">
        <q-toolbar class="bg-white border-b border-slate-200 q-px-xl q-py-md">
          <q-btn flat round dense icon="close" v-close-popup color="slate-500" />
          <q-toolbar-title class="text-weight-bold text-slate-800 text-outfit">Anteprima: {{ currentReportTitle }}</q-toolbar-title>
          <div class="row q-gutter-sm">
            <q-btn unelevated color="primary" icon="print" label="Stampa / PDF" class="rounded-lg q-px-lg shadow-sm" no-caps @click="printReport" />
            <q-btn outline color="secondary" icon="download" label="Scarica CSV" class="rounded-lg q-px-lg" no-caps @click="downloadCSV" />
          </div>
        </q-toolbar>

        <q-card-section class="col q-pa-xl scroll bg-slate-100">
          <div v-if="reportLoading" class="text-center q-pa-xl">
            <q-spinner color="primary" size="3em" />
            <div class="text-slate-500 q-mt-md">Caricamento dati report...</div>
          </div>
          <div v-else id="print-section" class="document-paper shadow-lg rounded-sm q-pa-xl bg-white mx-auto text-slate-800" style="max-width: 800px; min-height: 1000px; font-family: 'Inter', sans-serif; position: relative; border: 1px solid #e2e8f0; padding: 40px;">

            <!-- School Header -->
            <div class="row items-center justify-between q-pb-md q-mb-lg" style="border-bottom: 2px solid #6366f1; display: flex; justify-content: space-between; align-items: center;">
              <div>
                <div class="text-h5 text-weight-bold text-primary text-outfit" style="color: #6366f1; font-weight: 700; font-size: 1.5rem;">{{ schoolName }}</div>
                <div class="text-caption text-slate-500">Registro Elettronico V2 - Segreteria Scolastica</div>
              </div>
              <div class="text-right text-caption text-slate-500">
                <div>Data: {{ generationDate }}</div>
                <div>Classe: {{ selectedClassName }}</div>
              </div>
            </div>

            <!-- Report Body -->
            <div class="q-py-md">
              <div class="text-h6 text-weight-bold text-center text-uppercase text-slate-800 q-mb-lg">{{ currentReportTitle }}</div>
              <div class="text-subtitle2 text-center text-slate-500 q-mb-xl" v-if="info.period || info.month">Periodo/Mese: {{ info.period || info.month }}</div>

              <!-- No data state -->
              <div v-if="reportData.length === 0" class="text-center text-slate-400 q-pa-xl">
                <q-icon name="search_off" size="48px" class="q-mb-sm" />
                <div>Nessun dato trovato per i parametri selezionati</div>
              </div>

              <!-- 1. GRADES -->
              <div v-else-if="currentReport === 'grades'">
                <q-table :rows="reportData" :columns="reportColumns" row-key="student_id" flat bordered hide-pagination :pagination="{ rowsPerPage: 100 }" class="bg-white" />
              </div>
              <!-- 2. TOPICS -->
              <div v-else-if="currentReport === 'topics'">
                <q-table :rows="reportData" :columns="reportColumns" row-key="id" flat bordered hide-pagination :pagination="{ rowsPerPage: 100 }" class="bg-white" />
              </div>
              <!-- 3. ATTENDANCE -->
              <div v-else-if="currentReport === 'attendance'">
                <q-table :rows="reportData" :columns="reportColumns" row-key="student_id" flat bordered hide-pagination :pagination="{ rowsPerPage: 100 }" class="bg-white" />
              </div>
              <!-- 4. NOTES -->
              <div v-else-if="currentReport === 'notes'">
                <q-table :rows="reportData" :columns="reportColumns" row-key="id" flat bordered hide-pagination :pagination="{ rowsPerPage: 100 }" class="bg-white" />
              </div>
              <!-- 5. STUDENTS -->
              <div v-else-if="currentReport === 'students'">
                <q-table :rows="reportData" :columns="reportColumns" row-key="id" flat bordered hide-pagination :pagination="{ rowsPerPage: 100 }" class="bg-white" />
              </div>
              <!-- 6. CERTIFICATES -->
              <div v-else-if="currentReport === 'certificates'" class="text-body1 q-px-md text-justify" style="font-family: 'Times New Roman', Times, serif; font-size: 1.1rem; line-height: 1.8;">
                <p class="q-mb-xl">Si certifica che, in base ai registri d'iscrizione di questo Istituto Scolastico, gli studenti selezionati per la classe <strong>{{ selectedClassName }}</strong> risultano regolarmente iscritti e frequentano le lezioni per l'anno scolastico corrente.</p>
                <p class="q-mb-xl">Il presente certificato viene rilasciato su richiesta dell'interessato per gli usi consentiti dalla legge.</p>
                <div class="row justify-end q-mt-xl">
                  <div class="text-center col-6 q-pt-md" style="border-top: 1px solid #cbd5e1;">
                    <div class="text-caption text-slate-500">Timbro dell'Istituto e Firma</div>
                    <div class="text-weight-bold text-slate-700" style="margin-top: 48px;">La Segreteria Scolastica</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Footer -->
            <div class="border-t border-slate-200 q-pt-md q-mt-xl text-center text-caption text-slate-400" style="position: absolute; bottom: 40px; left: 40px; right: 40px;">
              Documento generato dal sistema Registro Elettronico V2. Pagina 1 di 1
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, exportFile } from 'quasar'
import { h } from 'vue'
import { QCard, QCardSection, QIcon } from 'quasar'
import { useAuthStore } from '@/stores/auth'
import adminService from 'src/services/adminService'
import { gradeService } from 'src/services/gradeService'
import { attendanceService } from 'src/services/attendanceService'
import api from 'src/services/api'

const $q = useQuasar()
const authStore = useAuthStore()

// ── School name from auth store (not hardcoded) ───────────────
const schoolName = computed(() => authStore.user?.school_name || authStore.user?.schoolName || 'Istituto Scolastico')

// ── Inline ReportCard component ───────────────────────────────
const ReportCard = (props, { emit }) => {
  return h(QCard, {
    class: 'col-12 col-sm-5 cursor-pointer hover-shadow rounded-xl border border-slate-100 shadow-soft',
    onClick: () => emit('click')
  }, () => [
    h(QCardSection, {}, () => [
      h('div', { class: 'row items-center no-wrap' }, [
        h(QIcon, { name: props.icon, size: 'md', color: 'primary', class: 'q-mr-md' }),
        h('div', {}, [h('div', { class: 'text-subtitle1 text-weight-bold text-slate-800' }, props.title)])
      ]),
      h('div', { class: 'text-caption text-slate-500 q-mt-sm' }, props.description)
    ])
  ])
}

// ── State ─────────────────────────────────────────────────────
const showDialog    = ref(false)
const showPreview   = ref(false)
const reportLoading = ref(false)
const currentReport = ref('')
const reportData    = ref([])
const reportColumns = ref([])
const classes       = ref([])

const info = ref({ class_id: null, period: 'I Quadrimestre', month: 'Gennaio' })
const generationDate = computed(() => new Date().toLocaleDateString('it-IT'))

const classOptions = computed(() =>
  classes.value.map(c => ({
    label: `${c.name || ''}${c.section || ''}${c.articolazione ? ' - ' + c.articolazione : ''}`,
    value: c.id
  }))
)

const selectedClassName = computed(() => {
  const opt = classOptions.value.find(o => o.value === info.value.class_id)
  return opt?.label || 'N/D'
})

const monthOptions = ['Gennaio','Febbraio','Marzo','Aprile','Maggio','Giugno','Luglio','Agosto','Settembre','Ottobre','Novembre','Dicembre']

const currentReportTitle = computed(() => ({
  grades: 'Report Voti e Pagelle',
  topics: 'Registro Attività Svolte',
  attendance: 'Report Assenze Mensili',
  notes: 'Report Note Disciplinari',
  students: 'Elenco Iscritti Classe',
  certificates: 'Certificato di Frequenza'
}[currentReport.value] || 'Generazione Report'))

const requiresPeriod = computed(() => currentReport.value === 'grades')
const requiresMonth  = computed(() => currentReport.value === 'attendance')

// ── Load classes on mount ─────────────────────────────────────
onMounted(async () => {
  try {
    const schoolId = authStore.user?.school_id || authStore.user?.schoolId
    if (schoolId) {
      const res = await adminService.getClasses(schoolId)
      classes.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
      if (classes.value.length > 0) info.value.class_id = classes.value[0].id
    }
  } catch (e) {
    console.error('Error loading classes:', e)
  }
})

// ── Open dialog ───────────────────────────────────────────────
const openReport = (type) => {
  currentReport.value = type
  reportData.value = []
  reportColumns.value = []
  showDialog.value = true
}

// ── Fetch real data from API ──────────────────────────────────
const loadReportData = async () => {
  if (!info.value.class_id && currentReport.value !== 'certificates') {
    $q.notify({ type: 'warning', message: 'Seleziona una classe prima di generare il report.' })
    return false
  }
  reportLoading.value = true
  reportData.value = []
  try {
    if (currentReport.value === 'grades') {
      // Voti: recupera tutti gli studenti della classe e i loro voti
      const studRes = await adminService.getClassStudents(info.value.class_id)
      const students = Array.isArray(studRes.data) ? studRes.data : (studRes.data?.items || [])
      const rows = []
      for (const s of students) {
        const gRes = await gradeService.getChildGrades(s.id)
        const allGrades = []
        if (gRes.data?.semesters) gRes.data.semesters.forEach(sem => { if (sem.grades) allGrades.push(...sem.grades) })
        // Group by subject
        const bySubject = {}
        allGrades.forEach(g => { bySubject[g.subject_name || g.subject_id] = g.grade_value })
        rows.push({ student_id: s.id, name: `${s.last_name} ${s.first_name}`, ...bySubject })
      }
      // Build columns dynamically from subjects
      const subjectKeys = [...new Set(rows.flatMap(r => Object.keys(r).filter(k => k !== 'student_id' && k !== 'name')))]
      reportColumns.value = [
        { name: 'name', label: 'Studente', field: 'name', align: 'left' },
        ...subjectKeys.map(k => ({ name: k, label: k, field: k, align: 'center' }))
      ]
      reportData.value = rows

    } else if (currentReport.value === 'topics') {
      const res = await api.get(`/classes/${info.value.class_id}/lesson-topics`)
      const items = Array.isArray(res.data) ? res.data : (res.data?.items || [])
      reportColumns.value = [
        { name: 'date', label: 'Data', field: r => new Date(r.date).toLocaleDateString('it-IT'), align: 'left' },
        { name: 'subject', label: 'Materia', field: 'subject_name', align: 'left' },
        { name: 'topic', label: 'Argomento Svolto', field: 'topic', align: 'left' },
        { name: 'teacher', label: 'Docente', field: r => `${r.teacher_last_name || ''} ${r.teacher_first_name || ''}`.trim(), align: 'left' }
      ]
      reportData.value = items

    } else if (currentReport.value === 'attendance') {
      const studRes = await adminService.getClassStudents(info.value.class_id)
      const students = Array.isArray(studRes.data) ? studRes.data : (studRes.data?.items || [])
      const rows = []
      for (const s of students) {
        const aRes = await attendanceService.getChildAttendance(s.id)
        const records = Array.isArray(aRes.data) ? aRes.data : (aRes.data?.records || [])
        const total   = records.length
        const absents = records.filter(r => r.status === 'absent').length
        rows.push({
          student_id: s.id,
          name: `${s.last_name} ${s.first_name}`,
          present: total - absents,
          absent: absents,
          rate: total > 0 ? Math.round(((total - absents) / total) * 100) + '%' : '-'
        })
      }
      reportColumns.value = [
        { name: 'name', label: 'Studente', field: 'name', align: 'left' },
        { name: 'present', label: 'Ore Presenza', field: 'present', align: 'center' },
        { name: 'absent', label: 'Ore Assenza', field: 'absent', align: 'center' },
        { name: 'rate', label: 'Tasso Presenza', field: 'rate', align: 'center', classes: 'text-weight-bold text-primary' }
      ]
      reportData.value = rows

    } else if (currentReport.value === 'notes') {
      const res = await api.get(`/classes/${info.value.class_id}/disciplinary-notes`)
      const items = Array.isArray(res.data) ? res.data : (res.data?.items || [])
      reportColumns.value = [
        { name: 'date', label: 'Data', field: r => new Date(r.date).toLocaleDateString('it-IT'), align: 'left' },
        { name: 'student', label: 'Studente', field: r => `${r.student_last_name || ''} ${r.student_first_name || ''}`.trim(), align: 'left' },
        { name: 'type', label: 'Tipo', field: 'note_type', align: 'left' },
        { name: 'description', label: 'Nota', field: 'description', align: 'left' },
        { name: 'teacher', label: 'Docente', field: r => `${r.teacher_last_name || ''} ${r.teacher_first_name || ''}`.trim(), align: 'left' }
      ]
      reportData.value = items

    } else if (currentReport.value === 'students') {
      const studRes = await adminService.getClassStudents(info.value.class_id)
      const students = Array.isArray(studRes.data) ? studRes.data : (studRes.data?.items || [])
      reportColumns.value = [
        { name: 'name',  label: 'Studente', field: r => `${r.last_name} ${r.first_name}`, align: 'left' },
        { name: 'birth', label: 'Data Nascita', field: r => r.birth_date ? new Date(r.birth_date).toLocaleDateString('it-IT') : '-', align: 'left' },
        { name: 'cf',    label: 'Codice Fiscale', field: 'fiscal_code', align: 'left' },
        { name: 'email', label: 'Email', field: 'email', align: 'left' }
      ]
      reportData.value = students

    } else if (currentReport.value === 'certificates') {
      // No tabular data needed, template renders text
      reportData.value = [{ id: 'cert' }]
    }
    return true
  } catch (e) {
    console.error('Error loading report data:', e)
    $q.notify({ type: 'negative', message: 'Errore nel caricamento dei dati del report.' })
    return false
  } finally {
    reportLoading.value = false
  }
}

const showPreviewDialog = async () => {
  showDialog.value = false
  showPreview.value = true
  await loadReportData()
}

const generatePDF = async () => {
  showDialog.value = false
  showPreview.value = true
  const ok = await loadReportData()
  if (ok) setTimeout(() => printReport(), 500)
}

const generateXLS = async () => {
  showDialog.value = false
  const ok = await loadReportData()
  if (ok) downloadCSV()
}

const printReport = () => {
  const printWindow = window.open('', '_blank')
  const htmlContent = document.getElementById('print-section')?.innerHTML || ''
  printWindow.document.write(`
    <html><head><title>${currentReportTitle.value}</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" rel="stylesheet">
    <style>
      body { font-family: 'Inter', sans-serif; padding: 40px; color: #334155; background: white; }
      table { width: 100%; border-collapse: collapse; margin: 20px 0; }
      th, td { border: 1px solid #e2e8f0; padding: 10px; text-align: left; }
      th { background: #f8fafc; font-weight: 600; }
      @page { size: auto; margin: 20mm; }
    </style></head>
    <body><div style="max-width:800px;margin:0 auto">${htmlContent}</div>
    <script>window.onload=function(){window.print();setTimeout(function(){window.close()},500)}<\/script>
    </body></html>
  `)
  printWindow.document.close()
}

const downloadCSV = () => {
  if (reportData.value.length === 0) {
    $q.notify({ type: 'warning', message: 'Nessun dato da esportare.' })
    return
  }

  // Build CSV from reportColumns and reportData
  const cols = reportColumns.value
  const header = cols.map(c => c.label).join(';')
  const rows = reportData.value.map(row =>
    cols.map(c => {
      const val = typeof c.field === 'function' ? c.field(row) : (row[c.field] ?? '')
      return `"${String(val).replace(/"/g, '""')}"`
    }).join(';')
  )
  const content = [header, ...rows].join('\n')
  const status = exportFile(`report_${currentReport.value}_${info.value.class_id}.csv`, content, 'text/csv')
  if (status === true) {
    $q.notify({ type: 'positive', message: 'Report CSV scaricato con successo', icon: 'download' })
  } else {
    $q.notify({ type: 'negative', message: 'Impossibile scaricare il file' })
  }
}

defineExpose({ openReport, currentReport, showDialog, showPreview, info, ReportCard })
</script>

<style scoped>
.hover-shadow:hover {
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
  background-color: #fafafa;
}
.document-paper {
  box-shadow: 0 10px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1);
}
</style>
