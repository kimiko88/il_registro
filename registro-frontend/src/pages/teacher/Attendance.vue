<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">Registro Presenze</div>
       <div class="row q-gutter-md">
           <q-input dense outlined v-model="date" type="date" label="Data" bg-color="white" />
           <q-select dense outlined v-model="selectedClass" :options="classes" label="Classe" bg-color="white" style="min-width: 120px" />
       </div>
    </div>

    <!-- Summary Cards -->
    <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-md-3">
            <q-card class="bg-green-1">
                <q-card-section class="text-center">
                    <div class="text-caption text-uppercase text-green-9">Presenti</div>
                    <div class="text-h4 text-green-8">{{ stats.present }}</div>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-12 col-md-3">
             <q-card class="bg-red-1">
                <q-card-section class="text-center">
                    <div class="text-caption text-uppercase text-red-9">Assenti</div>
                    <div class="text-h4 text-red-8">{{ stats.absent }}</div>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-12 col-md-3">
             <q-card class="bg-orange-1">
                <q-card-section class="text-center">
                    <div class="text-caption text-uppercase text-orange-9">Ritardi</div>
                    <div class="text-h4 text-orange-8">{{ stats.late }}</div>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-12 col-md-3">
             <q-card class="bg-blue-1 cursor-pointer" ripple @click="showJustifications = true">
                <q-card-section class="text-center">
                    <div class="text-caption text-uppercase text-blue-9">Da Giustificare</div>
                    <div class="text-h4 text-blue-8">{{ stats.toJustify }}</div>
                </q-card-section>
                <q-tooltip>Clicca per gestire</q-tooltip>
            </q-card>
        </div>
    </div>

    <!-- Attendance Table -->
    <q-card>
        <q-toolbar class="bg-grey-2 text-grey-8">
            <q-toolbar-title class="text-subtitle1">Appello - {{ date }}</q-toolbar-title>
            <q-btn flat dense icon="check_circle" label="Tutti Presenti" color="primary" @click="markAllPresent" />
        </q-toolbar>
        
        <q-list separator>
            <q-item v-for="student in students" :key="student.id" class="q-py-md">
                <q-item-section avatar>
                    <q-avatar size="md" color="grey-3" text-color="black">{{ student.name.charAt(0) }}</q-avatar>
                </q-item-section>
                
                <q-item-section>
                    <q-item-label class="text-weight-medium">{{ student.name }}</q-item-label>
                    <q-item-label caption v-if="student.status === 'absent'">Assente</q-item-label>
                    <q-item-label caption v-if="student.status === 'late'">In ritardo di {{ student.lateMinutes }} min</q-item-label>
                </q-item-section>

                <q-item-section>
                    <q-btn-toggle
                        v-model="student.status"
                        flat dense
                        :options="[
                            {icon: 'check', value: 'present', slot: 'present'},
                            {icon: 'close', value: 'absent', slot: 'absent'},
                            {icon: 'schedule', value: 'late', slot: 'late'},
                            {icon: 'logout', value: 'early_exit', slot: 'early'}
                        ]"
                    >
                        <template v-slot:present><q-tooltip>Presente</q-tooltip></template>
                        <template v-slot:absent><q-tooltip>Assente</q-tooltip></template>
                        <template v-slot:late><q-tooltip>Ritardo</q-tooltip></template>
                        <template v-slot:early><q-tooltip>Uscita Anticipata</q-tooltip></template>
                    </q-btn-toggle>
                </q-item-section>

                <q-item-section v-if="student.status === 'late' || student.status === 'early_exit'" side>
                     <q-input v-model.number="student.lateMinutes" type="number" dense outlined style="width: 60px" label="Min" />
                </q-item-section>
                
                <q-item-section side>
                    <q-btn round flat icon="chat_bubble_outline" color="grey" @click="openNoteDialog(student)" />
                </q-item-section>
            </q-item>
        </q-list>
        
        <q-card-actions align="right" class="bg-grey-1 q-pa-md">
            <q-btn label="Salva Registro" color="primary" size="lg" icon="save" @click="saveAttendance" :loading="saving" />
        </q-card-actions>
    </q-card>

    <!-- Justification Dialog -->
    <q-dialog v-model="showJustifications">
        <q-card style="min-width: 600px">
            <q-card-section class="text-h6">Gestione Giustificazioni</q-card-section>
            <q-list separator>
                <q-item v-for="req in justificationRequests" :key="req.id">
                    <q-item-section>
                        <q-item-label>{{ req.student }}</q-item-label>
                        <q-item-label caption>Assenza del {{ req.date }} - {{ req.reason }}</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <div class="row q-gutter-sm">
                            <q-btn flat round color="green" icon="check" @click="approveJustification(req.id)" />
                            <q-btn flat round color="red" icon="close" @click="rejectJustification(req.id)" />
                        </div>
                    </q-item-section>
                </q-item>
                <q-item v-if="justificationRequests.length === 0">
                    <q-item-section class="text-center text-grey">Nessuna richiesta in sospeso</q-item-section>
                </q-item>
            </q-list>
            <q-card-actions align="right"><q-btn flat label="Chiudi" v-close-popup /></q-card-actions>
        </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const date = ref(new Date().toISOString().split('T')[0])
const selectedClass = ref('5A')
const classes = ['5A', '4B', '3C']
const saving = ref(false)
const showJustifications = ref(false)

const students = ref([
    { id: 1, name: 'Rossi Mario', status: 'present', lateMinutes: 0 },
    { id: 2, name: 'Bianchi Anna', status: 'absent', lateMinutes: 0 },
    { id: 3, name: 'Verdi Paolo', status: 'present', lateMinutes: 0 },
    { id: 4, name: 'Neri Giulia', status: 'late', lateMinutes: 15 },
])

const justificationRequests = ref([
    { id: 101, student: 'Bianchi Anna', date: '2025-01-10', reason: 'Salute' },
    { id: 102, student: 'Neri Giulia', date: '2025-01-12', reason: 'Visita Medica' }
])

const stats = computed(() => ({
    present: students.value.filter(s => s.status === 'present').length,
    absent: students.value.filter(s => s.status === 'absent').length,
    late: students.value.filter(s => s.status === 'late').length,
    toJustify: justificationRequests.value.length
}))

const markAllPresent = () => {
    students.value.forEach(s => s.status = 'present')
}

const saveAttendance = () => {
    saving.value = true
    setTimeout(() => {
        saving.value = false
        $q.notify({ type: 'positive', message: 'Registro salvato con successo' })
    }, 1000)
}

const approveJustification = (id) => {
    justificationRequests.value = justificationRequests.value.filter(r => r.id !== id)
    $q.notify({ message: 'Giustificazione accettata', color: 'green' })
}

const rejectJustification = (id) => {
    justificationRequests.value = justificationRequests.value.filter(r => r.id !== id)
    $q.notify({ message: 'Giustificazione respinta', color: 'orange' })
}
</script>
