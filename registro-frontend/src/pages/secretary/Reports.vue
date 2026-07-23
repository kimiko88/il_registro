<template>
  <q-page padding>
    <div class="text-h4 text-weight-bold q-mb-lg">Centro Reportistica</div>

    <div class="row q-col-gutter-lg">
        <!-- Category: Didattica -->
        <div class="col-12 col-md-6">
            <div class="text-h6 q-mb-sm text-primary">Didattica & Voti</div>
            <div class="row q-gutter-md">
                <ReportCard 
                    title="Pagelle / Scrutini" 
                    icon="grading" 
                    description="Genera pagelle di fine quadrimestre o tabelloni scrutini." 
                    @click="openReport('grades')" 
                />
                <ReportCard 
                    title="Registro Attività" 
                    icon="menu_book" 
                    description="Riepilogo argomenti svolti per classe." 
                    @click="openReport('topics')" 
                />
            </div>
        </div>

        <!-- Category: Presenze -->
        <div class="col-12 col-md-6">
             <div class="text-h6 q-mb-sm text-secondary">Presenze & Note</div>
             <div class="row q-gutter-md">
                <ReportCard 
                    title="Assenze Mensili" 
                    icon="event_busy" 
                    description="Statistiche assenze per classe o studente." 
                    @click="openReport('attendance')" 
                />
                <ReportCard 
                    title="Note Disciplinari" 
                    icon="warning" 
                    description="Elenco sanzioni e note disciplinari." 
                    @click="openReport('notes')" 
                />
             </div>
        </div>

        <!-- Category: Amministrazione -->
        <div class="col-12 col-md-6">
             <div class="text-h6 q-mb-sm text-accent">Amministrazione</div>
             <div class="row q-gutter-md">
                <ReportCard 
                    title="Elenco Iscritti" 
                    icon="people_alt" 
                    description="Lista completa studenti per classe con dati anagrafici." 
                    @click="openReport('students')" 
                />
                <ReportCard 
                    title="Certificati" 
                    icon="verified" 
                    description="Certificati di frequenza e iscrizione." 
                    @click="openReport('certificates')" 
                />
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
                <q-select v-model="info.class" :options="['1A', '1B', '2A']" label="Classe" outlined dense />
                <q-select v-if="requiresPeriod" v-model="info.period" :options="['I Quadrimestre', 'II Quadrimestre']" label="Periodo" outlined dense />
                <q-select v-if="requiresMonth" v-model="info.month" :options="['Gennaio', 'Febbraio', 'Marzo', 'Aprile', 'Maggio', 'Giugno', 'Settembre', 'Ottobre', 'Novembre', 'Dicembre']" label="Mese" outlined dense />
            </q-card-section>

            <q-card-actions align="right" class="q-pa-md border-t border-slate-100 bg-slate-50">
                <q-btn flat label="Annulla" color="slate-500" v-close-popup no-caps />
                <q-btn outline color="primary" label="Anteprima" icon="visibility" class="rounded-lg" no-caps @click="showPreviewDialog" />
                <q-btn color="primary" label="Genera PDF" icon="picture_as_pdf" class="rounded-lg shadow-sm" no-caps @click="generatePDF" />
                <q-btn color="secondary" label="Excel" icon="table_view" class="rounded-lg shadow-sm" no-caps @click="generateXLS" />
            </q-card-actions>
        </q-card>
    </q-dialog>

    <!-- Preview Dialog -->
    <q-dialog v-model="showPreview" full-width full-height>
      <q-card class="column no-wrap bg-slate-100 rounded-xl overflow-hidden shadow-24">
        <q-toolbar class="bg-white border-b border-slate-200 q-px-xl q-py-md">
          <q-btn flat round dense icon="close" v-close-popup color="slate-500" />
          <q-toolbar-title class="text-weight-bold text-slate-800 text-outfit">
            Anteprima Documento: {{ currentReportTitle }}
          </q-toolbar-title>
          <div class="row q-gutter-sm">
            <q-btn unelevated color="primary" icon="print" label="Stampa / Esporta PDF" class="rounded-lg q-px-lg shadow-sm" no-caps @click="printReport" />
            <q-btn outline color="secondary" icon="download" label="Scarica Excel" class="rounded-lg q-px-lg" no-caps @click="downloadExcel" />
          </div>
        </q-toolbar>

        <q-card-section class="col q-pa-xl scroll bg-slate-100">
          <div id="print-section" class="document-paper shadow-lg rounded-sm q-pa-xl bg-white mx-auto text-slate-800" style="max-width: 800px; min-height: 1000px; font-family: 'Inter', sans-serif; position: relative; border: 1px solid #e2e8f0; padding: 40px;">
            
            <!-- School Header -->
            <div class="row items-center justify-between q-pb-md q-mb-lg" style="border-bottom: 2px solid #6366f1; display: flex; justify-content: space-between; align-items: center;">
              <div>
                <div class="text-h5 text-weight-bold text-primary text-outfit" style="color: #6366f1; font-weight: 700; font-size: 1.5rem;">ISTITUTO STATALE K.88</div>
                <div class="text-caption text-slate-500" style="font-size: 0.75rem; color: #64748b;">Registro Elettronico V2 - Segreteria Scolastica</div>
              </div>
              <div class="text-right text-caption text-slate-500" style="text-align: right; font-size: 0.75rem; color: #64748b;">
                <div>Data: {{ generationDate }}</div>
                <div>Classe: {{ info.class || 'N/D' }}</div>
              </div>
            </div>

            <!-- Report Body -->
            <div class="q-py-md">
              <div class="text-h6 text-weight-bold text-center text-uppercase text-slate-800 q-mb-lg" style="font-size: 1.25rem; font-weight: 700; text-align: center; text-transform: uppercase; margin-bottom: 16px;">{{ currentReportTitle }}</div>
              <div class="text-subtitle2 text-center text-slate-500 q-mb-xl" style="text-align: center; color: #64748b; margin-bottom: 48px;" v-if="info.period || info.month">
                Periodo/Mese: {{ info.period || info.month }}
              </div>

              <!-- Dynamic Content based on currentReport -->
              <!-- 1. GRADES -->
              <div v-if="currentReport === 'grades'">
                <q-table
                  :rows="mockGradesData"
                  :columns="gradesColumns"
                  row-key="name"
                  flat
                  bordered
                  hide-pagination
                  :pagination="{ rowsPerPage: 100 }"
                  class="bg-white"
                />
              </div>

              <!-- 2. TOPICS -->
              <div v-else-if="currentReport === 'topics'">
                <q-table
                  :rows="mockTopicsData"
                  :columns="topicsColumns"
                  row-key="date"
                  flat
                  bordered
                  hide-pagination
                  :pagination="{ rowsPerPage: 100 }"
                  class="bg-white"
                />
              </div>

              <!-- 3. ATTENDANCE -->
              <div v-else-if="currentReport === 'attendance'">
                <q-table
                  :rows="mockAttendanceData"
                  :columns="attendanceColumns"
                  row-key="name"
                  flat
                  bordered
                  hide-pagination
                  :pagination="{ rowsPerPage: 100 }"
                  class="bg-white"
                />
              </div>

              <!-- 4. NOTES -->
              <div v-else-if="currentReport === 'notes'">
                <q-table
                  :rows="mockNotesData"
                  :columns="notesColumns"
                  row-key="id"
                  flat
                  bordered
                  hide-pagination
                  :pagination="{ rowsPerPage: 100 }"
                  class="bg-white"
                />
              </div>

              <!-- 5. STUDENTS -->
              <div v-else-if="currentReport === 'students'">
                <q-table
                  :rows="mockStudentsData"
                  :columns="studentsColumns"
                  row-key="name"
                  flat
                  bordered
                  hide-pagination
                  :pagination="{ rowsPerPage: 100 }"
                  class="bg-white"
                />
              </div>

              <!-- 6. CERTIFICATES -->
              <div v-else-if="currentReport === 'certificates'" class="text-body1 line-height-relaxed q-px-md text-justify" style="font-family: 'Times New Roman', Times, serif; font-size: 1.1rem; line-height: 1.8;">
                <p class="q-mb-xl" style="margin-bottom: 24px;">Si certifica che, in base ai registri d'iscrizione di questo Istituto Scolastico, gli studenti selezionati per la classe <strong>{{ info.class || 'N/D' }}</strong> risultano regolarmente iscritti e frequentano le lezioni per l'anno scolastico corrente.</p>
                <p class="q-mb-xl" style="margin-bottom: 24px;">Il presente certificato viene rilasciato su richiesta dell'interessato per gli usi consentiti dalla legge.</p>
                <div class="row justify-end q-mt-xl" style="margin-top: 150px; display: flex; justify-content: flex-end;">
                  <div class="text-center col-6 q-pt-md" style="border-top: 1px solid #cbd5e1; text-align: center; width: 50%;">
                    <div class="text-caption text-slate-500" style="font-size: 0.75rem; color: #64748b;">Timbro dell'Istituto e Firma</div>
                    <div class="text-weight-bold text-slate-700" style="margin-top: 48px; font-weight: 700; color: #334155;">La Segreteria Scolastica</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Report Footer -->
            <div class="border-t border-slate-200 q-pt-md q-mt-xl text-center text-caption text-slate-400" style="position: absolute; bottom: 40px; left: 40px; right: 40px; border-top: 1px solid #e2e8f0; padding-top: 16px; text-align: center; font-size: 0.75rem; color: #94a3b8;">
              Documento generato dal sistema Registro Elettronico V2. Pagina 1 di 1
            </div>

          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar, exportFile } from 'quasar'
import { h } from 'vue'
import { QCard, QCardSection, QIcon } from 'quasar'

const $q = useQuasar()

// Inline Component for Card
const ReportCard = (props, { emit }) => {
    return h(QCard, { 
        class: 'col-12 col-sm-5 cursor-pointer hover-shadow rounded-xl border border-slate-100 shadow-soft', 
        onClick: () => emit('click') 
    }, () => [
        h(QCardSection, {}, () => [
            h('div', { class: 'row items-center no-wrap' }, [
                h(QIcon, { name: props.icon, size: 'md', color: 'primary', class: 'q-mr-md' }),
                h('div', {}, [
                    h('div', { class: 'text-subtitle1 text-weight-bold text-slate-800' }, props.title),
                ])
            ]),
            h('div', { class: 'text-caption text-slate-500 q-mt-sm' }, props.description)
        ])
    ])
}

const showDialog = ref(false)
const showPreview = ref(false)
const currentReport = ref('')
const info = ref({ class: '1A', period: 'I Quadrimestre', month: 'Gennaio' })
const generationDate = computed(() => new Date().toLocaleDateString('it-IT'))

const currentReportTitle = computed(() => {
    switch(currentReport.value) {
        case 'grades': return 'Report Voti e Pagelle';
        case 'topics': return 'Registro Attività Svolte';
        case 'attendance': return 'Report Assenze Mensili';
        case 'notes': return 'Report Note Disciplinari';
        case 'students': return 'Elenco Iscritti Classe';
        case 'certificates': return 'Certificato di Frequenza';
        default: return 'Generazione Report';
    }
})

const requiresPeriod = computed(() => currentReport.value === 'grades')
const requiresMonth = computed(() => currentReport.value === 'attendance')

const openReport = (type) => {
    currentReport.value = type
    info.value = { class: '1A', period: 'I Quadrimestre', month: 'Gennaio' }
    showDialog.value = true
}

const showPreviewDialog = () => {
    showDialog.value = false
    showPreview.value = true
}

const generatePDF = () => {
    showDialog.value = false
    showPreview.value = true
    setTimeout(() => {
        printReport()
    }, 500)
}

const generateXLS = () => {
    showDialog.value = false
    downloadExcel()
}

const printReport = () => {
  const printWindow = window.open('', '_blank');
  const htmlContent = document.getElementById('print-section').innerHTML;
  
  printWindow.document.write(`
    <html>
      <head>
        <title>${currentReportTitle.value}</title>
        <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" rel="stylesheet">
        <style>
          body {
            font-family: 'Inter', sans-serif;
            padding: 40px;
            color: #334155;
            background-color: white;
          }
          .border-b-2 { border-bottom: 2px solid #6366f1; }
          .text-primary { color: #6366f1; }
          table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
            margin-bottom: 20px;
          }
          th, td {
            border: 1px solid #e2e8f0;
            padding: 10px;
            text-align: left;
          }
          th {
            background-color: #f8fafc;
            font-weight: 600;
          }
          .text-center { text-align: center !important; }
          .text-right { text-align: right !important; }
          .row { display: flex; flex-wrap: wrap; }
          .justify-between { justify-content: space-between; }
          .justify-end { justify-content: flex-end; }
          .col-6 { width: 50%; }
          .q-mb-md { margin-bottom: 16px; }
          .q-mb-lg { margin-bottom: 24px; }
          .q-mb-xl { margin-bottom: 48px; }
          .q-mt-xl { margin-top: 48px; }
          .q-pt-md { padding-top: 16px; }
          .text-weight-bold { font-weight: 700; }
          .text-caption { font-size: 0.75rem; color: #64748b; }
          .text-slate-500 { color: #64748b; }
          .text-slate-700 { color: #334155; }
          .text-h5 { font-size: 1.5rem; }
          .text-h6 { font-size: 1.25rem; }
          /* Hide scrollbars and elements not in print */
          @page { size: auto; margin: 20mm; }
        </style>
      </head>
      <body>
        <div style="width: 100%; max-width: 800px; margin: 0 auto;">
          ${htmlContent}
        </div>
        <script>
          window.onload = function() {
            window.print();
            setTimeout(function() { window.close(); }, 500);
          };
        </${'script'}>
      </body>
    </html>
  `);
  printWindow.document.close();
}

// Mock Data definitions
const mockGradesData = ref([
  { name: 'Rossi Mario', ita: '7.5', mat: '8.0', sto: '6.5', ing: '7.0', sci: '8.5', avg: '7.5' },
  { name: 'Bianchi Luigi', ita: '6.0', mat: '5.5', sto: '7.0', ing: '6.0', sci: '6.0', avg: '6.1' },
  { name: 'Verdi Anna', ita: '9.0', mat: '9.5', sto: '8.5', ing: '9.0', sci: '9.5', avg: '9.1' },
  { name: 'Ferrari Giulia', ita: '8.0', mat: '7.5', sto: '8.0', ing: '8.5', sci: '8.0', avg: '8.0' }
])

const gradesColumns = [
  { name: 'name', label: 'Studente', field: 'name', align: 'left' },
  { name: 'ita', label: 'Italiano', field: 'ita', align: 'center' },
  { name: 'mat', label: 'Matematica', field: 'mat', align: 'center' },
  { name: 'sto', label: 'Storia', field: 'sto', align: 'center' },
  { name: 'ing', label: 'Inglese', field: 'ing', align: 'center' },
  { name: 'sci', label: 'Scienze', field: 'sci', align: 'center' },
  { name: 'avg', label: 'Media', field: 'avg', align: 'center', classes: 'text-weight-bold' }
]

const mockTopicsData = ref([
  { date: '10/01/2026', subject: 'Matematica', topic: 'Equazioni di secondo grado e sistemi lineari', teacher: 'Prof. Neri' },
  { date: '12/01/2026', subject: 'Italiano', topic: 'Analisi dei Promessi Sposi - Capitolo X', teacher: 'Prof.ssa Rosa' },
  { date: '14/01/2026', subject: 'Storia', topic: 'La Rivoluzione Industriale e lo sviluppo economico', teacher: 'Prof. Bruno' },
  { date: '15/01/2026', subject: 'Inglese', topic: 'Present Perfect vs Past Simple exercises', teacher: 'Prof.ssa Bianchi' }
])

const topicsColumns = [
  { name: 'date', label: 'Data', field: 'date', align: 'left' },
  { name: 'subject', label: 'Materia', field: 'subject', align: 'left' },
  { name: 'topic', label: 'Argomento Svolto', field: 'topic', align: 'left' },
  { name: 'teacher', label: 'Docente', field: 'teacher', align: 'left' }
]

const mockAttendanceData = ref([
  { name: 'Rossi Mario', present: '92', absent: '8', rate: '92%' },
  { name: 'Bianchi Luigi', present: '84', absent: '16', rate: '84%' },
  { name: 'Verdi Anna', present: '98', absent: '2', rate: '98%' },
  { name: 'Ferrari Giulia', present: '95', absent: '5', rate: '95%' }
])

const attendanceColumns = [
  { name: 'name', label: 'Studente', field: 'name', align: 'left' },
  { name: 'present', label: 'Ore Presenza', field: 'present', align: 'center' },
  { name: 'absent', label: 'Ore Assenza', field: 'absent', align: 'center' },
  { name: 'rate', label: 'Tasso Presenza', field: 'rate', align: 'center', classes: 'text-weight-bold text-primary' }
]

const mockNotesData = ref([
  { id: 1, date: '12/01/2026', student: 'Bianchi Luigi', type: 'Nota Disciplinare', description: 'Disturba ripetutamente la lezione nonostante i richiami', teacher: 'Prof. Neri' },
  { id: 2, date: '15/01/2026', student: 'Rossi Mario', type: 'Richiamo', description: 'Ritardo ingiustificato di 25 minuti alla prima ora', teacher: 'Prof.ssa Rosa' }
])

const notesColumns = [
  { name: 'date', label: 'Data', field: 'date', align: 'left' },
  { name: 'student', label: 'Studente', field: 'student', align: 'left' },
  { name: 'type', label: 'Tipo', field: 'type', align: 'left' },
  { name: 'description', label: 'Nota Disciplinare', field: 'description', align: 'left' },
  { name: 'teacher', label: 'Docente', field: 'teacher', align: 'left' }
]

const mockStudentsData = ref([
  { name: 'Rossi Mario', birth: '12/04/2010 (Roma)', cf: 'RSSMRA10D12H501U', email: 'mario.rossi@scuola.it' },
  { name: 'Bianchi Luigi', birth: '23/08/2010 (Milano)', cf: 'BNCLGU10M23F205R', email: 'luigi.bianchi@scuola.it' },
  { name: 'Verdi Anna', birth: '05/11/2010 (Napoli)', cf: 'VRDNNA10S45F839O', email: 'anna.verdi@scuola.it' },
  { name: 'Ferrari Giulia', birth: '19/01/2010 (Torino)', cf: 'FRRGLI10A59L219Y', email: 'giulia.ferrari@scuola.it' }
])

const studentsColumns = [
  { name: 'name', label: 'Studente', field: 'name', align: 'left' },
  { name: 'birth', label: 'Data & Luogo Nascita', field: 'birth', align: 'left' },
  { name: 'cf', label: 'Codice Fiscale', field: 'cf', align: 'left' },
  { name: 'email', label: 'Email Istituzionale', field: 'email', align: 'left' }
]

const downloadExcel = () => {
  let content = ''
  let fileName = `report_${currentReport.value}.csv`
  
  if (currentReport.value === 'grades') {
    content = "Studente;Italiano;Matematica;Storia;Inglese;Scienze;Media\n"
    mockGradesData.value.forEach(r => {
      content += `${r.name};${r.ita};${r.mat};${r.sto};${r.ing};${r.sci};${r.avg}\n`
    })
  } else if (currentReport.value === 'attendance') {
    content = "Studente;Ore Presenza;Ore Assenza;Tasso Presenza\n"
    mockAttendanceData.value.forEach(r => {
      content += `${r.name};${r.present};${r.absent};${r.rate}\n`
    })
  } else if (currentReport.value === 'students') {
    content = "Studente;Data Nascita;Codice Fiscale;Email\n"
    mockStudentsData.value.forEach(r => {
      content += `${r.name};${r.birth};${r.cf};${r.email}\n`
    })
  } else if (currentReport.value === 'topics') {
    content = "Data;Materia;Argomento;Docente\n"
    mockTopicsData.value.forEach(r => {
      content += `${r.date};${r.subject};${r.topic};${r.teacher}\n`
    })
  } else if (currentReport.value === 'notes') {
    content = "Data;Studente;Tipo;Nota;Docente\n"
    mockNotesData.value.forEach(r => {
      content += `${r.date};${r.student};${r.type};${r.description};${r.teacher}\n`
    })
  } else {
    content = "Certificato di frequenza per classe " + info.value.class
  }

  const status = exportFile(fileName, content, 'text/csv')
  if (status === true) {
    $q.notify({ type: 'positive', message: 'Report Excel/CSV scaricato con successo', icon: 'download' })
  } else {
    $q.notify({ type: 'negative', message: 'Impossibile scaricare il file' })
  }
}

defineExpose({
    openReport,
    currentReport,
    showDialog,
    showPreview,
    info,
    ReportCard
})
</script>

<style scoped>
.hover-shadow:hover {
    box-shadow: 0 4px 8px rgba(0,0,0,0.15);
    background-color: #fafafa;
}
.document-paper {
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
  aspect-ratio: 1 / 1.414;
}
</style>
