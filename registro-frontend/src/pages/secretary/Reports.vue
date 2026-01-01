<template>
  <q-page padding>
    <div class="text-h4 text-weight-bold q-mb-lg">Centro Reportistica</div>

    <div class="row q-col-gutter-lg">
        <!-- Cateory: Didattica -->
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
        <q-card style="min-width: 400px">
            <q-card-section>
                <div class="text-h6">{{ currentReportTitle }}</div>
                <div class="text-caption">Seleziona i parametri per generare il report.</div>
            </q-card-section>
            
            <q-card-section class="q-gutter-md">
                <q-select v-model="info.class" :options="['1A', '1B', '2A']" label="Classe" outlined dense />
                <q-select v-if="requiresPeriod" v-model="info.period" :options="['I Quadrimestre', 'II Quadrimestre']" label="Periodo" outlined dense />
                <q-select v-if="requiresMonth" v-model="info.month" :options="['Settembre', 'Ottobre', 'Novembre']" label="Mese" outlined dense />
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Annulla" v-close-popup />
                <q-btn color="primary" label="Genera PDF" icon="picture_as_pdf" @click="generate('pdf')" />
                <q-btn color="secondary" label="Excel" icon="table_view" @click="generate('xls')" />
            </q-card-actions>
        </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { h } from 'vue'
import { QCard, QCardSection, QIcon, QBtn } from 'quasar'

const $q = useQuasar()

// Inline Component for Card
const ReportCard = (props, { emit }) => {
    return h(QCard, { 
        class: 'col-12 col-sm-5 cursor-pointer hover-shadow', 
        onClick: () => emit('click') 
    }, () => [
        h(QCardSection, {}, () => [
            h('div', { class: 'row items-center no-wrap' }, [
                h(QIcon, { name: props.icon, size: 'md', color: 'primary', class: 'q-mr-md' }),
                h('div', {}, [
                    h('div', { class: 'text-subtitle1 text-weight-bold' }, props.title),
                ])
            ]),
            h('div', { class: 'text-caption text-grey q-mt-sm' }, props.description)
        ])
    ])
}

const showDialog = ref(false)
const currentReport = ref('')
const info = ref({ class: '', period: '', month: '' })

const currentReportTitle = computed(() => {
    switch(currentReport.value) {
        case 'grades': return 'Report Voti e Pagelle';
        case 'attendance': return 'Report Assenze';
        default: return 'Generazione Report';
    }
})

const requiresPeriod = computed(() => currentReport.value === 'grades')
const requiresMonth = computed(() => currentReport.value === 'attendance')

const openReport = (type) => {
    currentReport.value = type
    showDialog.value = true
}

const generate = (format) => {
    $q.loading.show({ message: `Generazione ${format.toUpperCase()} in corso...` })
    setTimeout(() => {
        $q.loading.hide()
        showDialog.value = false
        $q.notify({ type: 'positive', message: 'Report scaricato con successo', icon: 'download' })
    }, 1500)
}

defineExpose({
    openReport,
    generate,
    currentReport,
    showDialog,
    info,
    ReportCard
})
</script>

<style scoped>
.hover-shadow:hover {
    box-shadow: 0 4px 8px rgba(0,0,0,0.15);
    background-color: #fafafa;
}
</style>
