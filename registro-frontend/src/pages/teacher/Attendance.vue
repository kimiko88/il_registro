<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
       <div class="text-h4">Registro Presenze</div>
       <div class="row q-gutter-md items-center">
           <q-btn color="secondary" icon="download" label="Esporta CSV" unelevated @click="exportCSV" />
           <q-input dense outlined v-model="date" type="date" label="Data" bg-color="white" @update:model-value="fetchData" />
           <q-select 
              dense outlined 
              v-model="selectedHour" 
              :options="[1,2,3,4,5,6,7,8]" 
              label="Ora" 
              bg-color="white" 
              style="min-width: 80px"
              @update:model-value="fetchData"
           />
           <q-select 
              dense outlined 
              v-model="selectedSubject" 
              :options="gradesStore.subjects" 
              option-label="subject_name"
              option-value="subject_id"
              emit-value
              map-options
              label="Materia" 
              bg-color="white" 
              style="min-width: 150px"
              @update:model-value="fetchData"
           />
           <q-select 
              dense outlined 
              v-model="selectedClass" 
              :options="classesStore.classes" 
              option-label="name"
              option-value="id"
              label="Classe" 
              bg-color="white" 
              style="min-width: 150px"
              @update:model-value="onClassChange"
           />
       </div>
    </div>

    <!-- Summary Cards -->
    <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-md-3">
            <q-card class="bg-green-1">
                <q-card-section class="text-center">
                    <q-icon name="check_circle" size="28px" color="green-8" class="q-mb-xs" />
                    <div class="text-caption text-uppercase text-green-9">Presenti</div>
                    <div class="text-h4 text-green-8">{{ stats.present }}</div>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-12 col-md-3">
             <q-card class="bg-red-1">
                <q-card-section class="text-center">
                    <q-icon name="cancel" size="28px" color="red-8" class="q-mb-xs" />
                    <div class="text-caption text-uppercase text-red-9">Assenti</div>
                    <div class="text-h4 text-red-8">{{ stats.absent }}</div>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-12 col-md-3">
             <q-card class="bg-orange-1">
                <q-card-section class="text-center">
                    <q-icon name="schedule" size="28px" color="orange-8" class="q-mb-xs" />
                    <div class="text-caption text-uppercase text-orange-9">Ritardi</div>
                    <div class="text-h4 text-orange-8">{{ stats.late }}</div>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-12 col-md-3">
             <q-card class="bg-blue-1 cursor-pointer" ripple @click="showJustifications = true">
                <q-card-section class="text-center">
                    <q-icon name="assignment_turned_in" size="28px" color="blue-8" class="q-mb-xs" />
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
            <q-toolbar-title class="text-subtitle1 row items-center">
                <span>Appello - {{ date }}</span>
                <q-chip dense color="primary" text-color="white" class="q-ml-md font-weight-bold">
                    {{ markedCount }}/{{ students.length }} registrati
                </q-chip>
                <q-chip v-if="lastAutosaveTime" dense color="grey-7" text-color="white" icon="cloud_done" class="q-ml-sm text-caption">
                    Bozza salvata alle {{ lastAutosaveTime }}
                </q-chip>
            </q-toolbar-title>
            <q-btn flat dense icon="check_circle" label="Tutti Presenti" color="primary" @click="markAllPresent" :disable="loading" />
        </q-toolbar>

        <!-- Daily Timeline (Previous Hours) -->
        <div v-if="dailyLessons.length > 0" class="q-px-md q-py-sm bg-blue-50 border-b">
            <div class="text-caption text-weight-bold text-blue-9 q-mb-xs">ATTIVITÀ DEL GIORNO</div>
            <div class="row q-gutter-xs">
                <div v-for="l in dailyLessons" :key="l.id" class="col-auto">
                    <q-chip dense outline color="blue-7" text-color="white" icon="history">
                        Ora {{ l.hour }}: {{ l.topic }}
                        <q-tooltip>
                            Materia: {{ l.subject_id }}<br>
                            Tipo: {{ l.type }}
                        </q-tooltip>
                    </q-chip>
                </div>
            </div>
        </div>
        
        <div v-if="loading" class="row justify-center q-pa-lg">
            <q-spinner color="primary" size="3em" />
        </div>

        <q-list separator v-else>
            <q-item v-for="student in students" :key="student.id" class="q-py-md transition-bg" :class="getRowClass(student.status)">
                <q-item-section avatar>
                    <q-avatar size="md" color="grey-3" text-color="black">
                        {{ student.first_name ? student.first_name.charAt(0) : '?' }}
                    </q-avatar>
                </q-item-section>
                
                <q-item-section>
                    <div class="row items-center">
                       <div class="col">
                           <q-item-label class="text-weight-medium row items-center">
                               <q-icon v-if="student.status === 'Present'" name="check_circle" color="positive" size="18px" class="q-mr-xs" />
                               <q-icon v-else-if="student.status === 'Absent'" name="cancel" color="negative" size="18px" class="q-mr-xs" />
                               <q-icon v-else-if="student.status === 'Late'" name="schedule" color="warning" size="18px" class="q-mr-xs" />
                               <q-icon v-else-if="student.status === 'LeftEarly'" name="output" color="purple" size="18px" class="q-mr-xs" />
                               <span>{{ student.last_name }} {{ student.first_name }}</span>
                           </q-item-label>
                           <q-item-label caption v-if="student.status === 'Absent'">Assente</q-item-label>
                           <q-item-label caption v-if="student.status === 'Late'">
                               Ritardo ({{ formatLateLabel(student) }})
                           </q-item-label>
                           <q-item-label caption v-if="student.status === 'LeftEarly'">
                               Uscita Anticipata ({{ formatEarlyExitLabel(student) }})
                           </q-item-label>
                        </div>
                    </div>
                </q-item-section>

                <q-item-section>
                    <q-btn-toggle
                        v-model="student.status"
                        flat dense
                        :options="[
                            {icon: 'check', value: 'Present', slot: 'present', attrs: { 'aria-label': 'Segna ' + student.first_name + ' ' + student.last_name + ' come Presente' }},
                            {icon: 'close', value: 'Absent', slot: 'absent', attrs: { 'aria-label': 'Segna ' + student.first_name + ' ' + student.last_name + ' come Assente' }},
                            {icon: 'schedule', value: 'Late', slot: 'late', attrs: { 'aria-label': 'Segna ' + student.first_name + ' ' + student.last_name + ' in Ritardo' }},
                            {icon: 'logout', value: 'LeftEarly', slot: 'early', attrs: { 'aria-label': 'Segna ' + student.first_name + ' ' + student.last_name + ' come Uscita Anticipata' }}
                        ]"
                    >
                        <template v-slot:present><q-tooltip>Presente</q-tooltip></template>
                        <template v-slot:absent><q-tooltip>Assente</q-tooltip></template>
                        <template v-slot:late><q-tooltip>Ritardo</q-tooltip></template>
                        <template v-slot:early><q-tooltip>Uscita Anticipata</q-tooltip></template>
                    </q-btn-toggle>
                </q-item-section>

                <!-- Late Time Input -->
                <q-item-section v-if="student.status === 'Late'" side style="min-width: 120px">
                     <q-input 
                        v-model="student.entry_time" 
                        type="time" 
                        dense outlined 
                        label="Ora Ingresso" 
                        :rules="[val => !!val || 'Richiesto']"
                     />
                </q-item-section>

                <!-- Early Exit Time Input -->
                <q-item-section v-if="student.status === 'LeftEarly'" side style="min-width: 120px">
                     <q-input 
                        v-model="student.exit_time" 
                        type="time" 
                        dense outlined 
                        label="Ora Uscita" 
                        :rules="[val => !!val || 'Richiesto']"
                     />
                </q-item-section>
                
                <q-item-section side>
                    <q-btn round flat icon="note_add" color="grey-7" @click="openNoteDialog(student)">
                        <q-tooltip>Aggiungi Nota</q-tooltip>
                    </q-btn>
                </q-item-section>
            </q-item>

            <q-item v-if="students.length === 0" class="text-center text-grey">
                <q-item-section>Nessuno studente in questa classe (o seleziona una classe)</q-item-section>
            </q-item>
        </q-list>
        
        <q-card-actions align="right" class="bg-grey-1 q-pa-md">
            <q-btn label="Salva Registro" color="primary" size="lg" icon="save" @click="saveAttendance" :loading="saving" :disable="!selectedClass" />
        </q-card-actions>
    </q-card>

    <!-- Justification Dialog -->
    <q-dialog v-model="showJustifications">
        <q-card style="min-width: 600px">
            <q-card-section class="text-h6">Gestione Giustificazioni</q-card-section>
            <q-list separator>
                <q-item v-for="req in justificationRequests" :key="req.id">
                    <q-item-section>
                        <q-item-label>{{ req.student_name }}</q-item-label>
                        <q-item-label caption>Assenza del {{ req.date }} - {{ req.reason }}</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <div class="row q-gutter-sm">
                            <q-btn flat round color="green" icon="check" @click="processJustification(req.id, true)" />
                            <q-btn flat round color="red" icon="close" @click="processJustification(req.id, false)" />
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

    <!-- Note Dialog -->
    <NoteDialog
        v-if="selectedClass" 
        v-model="showNoteDialog"
        :student="selectedStudentForNote"
        :class-id="String(selectedClass.id)" 
    />

  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useQuasar } from 'quasar'
import { useClassesStore } from '@/stores/classes'
import { useGradesStore } from '@/stores/grades'
import { attendanceService } from 'src/services/attendanceService'
import api from '@/services/api'
import NoteDialog from 'src/components/Teacher/NoteDialog.vue'

const $q = useQuasar()
const classesStore = useClassesStore()
const gradesStore = useGradesStore()

const date = ref(new Date().toISOString().split('T')[0])
const selectedClass = ref(null)
const students = ref([])
const justificationRequests = ref([])
const loading = ref(false)
const saving = ref(false)
const showJustifications = ref(false)
const lastAutosaveTime = ref('')

// Note Dialog State
const showNoteDialog = ref(false)
const selectedStudentForNote = ref(null)

const selectedHour = ref(1)
const selectedSubject = ref(null)

const dailyLessons = ref([])

const stats = computed(() => ({
    present: students.value.filter(s => s.status === 'Present').length,
    absent: students.value.filter(s => s.status === 'Absent').length,
    late: students.value.filter(s => s.status === 'Late').length,
    toJustify: justificationRequests.value.length
}))

const markedCount = computed(() => students.value.filter(s => s.status).length)

const getRowClass = (status) => {
    switch (status) {
        case 'Present': return 'bg-green-1'
        case 'Absent': return 'bg-red-1'
        case 'Late': return 'bg-orange-1'
        case 'LeftEarly': return 'bg-purple-1'
        default: return ''
    }
}

let autosaveInterval = null
const saveDraftToStorage = () => {
    if (!selectedClass.value || students.value.length === 0) return
    const classId = typeof selectedClass.value === 'object' ? selectedClass.value?.id : selectedClass.value
    const key = `attendance_draft_${classId}_${date.value}_${selectedHour.value}`
    const draftData = {
        date: date.value,
        hour: selectedHour.value,
        students: students.value,
        timestamp: new Date().toISOString()
    }
    try {
        localStorage.setItem(key, JSON.stringify(draftData))
        const now = new Date()
        lastAutosaveTime.value = now.toLocaleTimeString('it-IT')
    } catch (e) {
        console.warn('Failed to save draft to localStorage', e)
        $q.notify({ type: 'warning', message: 'Impossibile salvare la bozza in memoria locale' })
    }
}

onMounted(async () => {
    await classesStore.fetchAssignedClasses()
    if (classesStore.classes.length > 0) {
        selectedClass.value = classesStore.classes[0]
    }
    fetchData()
    autosaveInterval = setInterval(saveDraftToStorage, 60000)
})

onUnmounted(() => {
    if (autosaveInterval) clearInterval(autosaveInterval)
})

const onClassChange = async () => {
    if (selectedClass.value) {
        await gradesStore.fetchClassSubjects(selectedClass.value.id)
        if (gradesStore.subjects.length > 0) {
            selectedSubject.value = gradesStore.subjects[0].subject_id
        } else {
            selectedSubject.value = null
        }
    }
    fetchData()
}

const fetchData = async () => {
    if (!selectedClass.value) return
    
    loading.value = true
    try {
        // 1. Fetch Students
        // Using users endpoint filtered by class
        const usersRes = await api.get('/users', { 
            params: { 
                class_id: selectedClass.value.id,
                role: 'student',
                page_size: 100 // Ensure we get all students, usually classes are small
            }
        })
        const studentList = usersRes.data.users || []
        
        // 2. Fetch Existing Attendance
        // attendanceService.getByClass likely returns list of records
        let attendanceMap = {}
        try {
            const attRes = await attendanceService.getByClass(selectedClass.value.id, date.value)
            // Filter attendance by selected hour
            if (Array.isArray(attRes.data)) {
                attRes.data.forEach(r => {
                    // Check if record hour matches selectedHour
                    // Backend returns hour as int
                    if (String(r.hour) === String(selectedHour.value)) {
                        attendanceMap[r.student_id] = r
                    }
                })
            }
        } catch (e) {
            console.warn("No attendance found or error", e)
        }

        // Merge
        students.value = studentList.map(s => {
            const existing = attendanceMap[s.id]
            return {
                id: s.id,
                first_name: s.first_name,
                last_name: s.last_name,
                status: existing ? existing.status : null, // Default null for unrecorded students
                entry_time: existing ? existing.entry_time : '',
                exit_time: existing ? existing.exit_time : '',
            }
        })

        // 3. Fetch Justifications
        const justRes = await api.get('/attendance/pending-justifications', {
             params: { class_id: selectedClass.value.id }
        })
        justificationRequests.value = justRes.data || []

        // 4. Fetch Daily Lessons for Timeline
        const lessonRes = await api.get(`/lessons/class/${selectedClass.value.id}`, {
            params: { date: date.value }
        })
        dailyLessons.value = lessonRes.data || []

    } catch (error) {
        $q.notify({ type: 'negative', message: 'Errore caricamento dati' })
        console.error(error)
    } finally {
        loading.value = false
    }
}

const markAllPresent = () => {
    $q.dialog({
        title: 'Conferma Operazione',
        message: 'Segnare tutti gli studenti come PRESENTI per l\'ora selezionata? Eventuali assenze già digitate verranno sovrascritte.',
        cancel: true,
        persistent: true
    }).onOk(() => {
        students.value.forEach(s => {
            s.status = 'Present'
            s.entry_time = ''
            s.exit_time = ''
        })
    })
}

const formatEarlyExitLabel = (student) => {
    if (student.exit_time) {
        return `Uscita: ${student.exit_time}`
    }
    return 'Inserisci orario'
}

const formatLateLabel = (student) => {
    if (student.entry_time) {
        return `Ingresso: ${student.entry_time}`
    }
    return 'Inserisci orario'
}

const saveAttendance = async () => {
    const unmarked = students.value.filter(s => !s.status)
    if (unmarked.length > 0) {
        $q.notify({
            type: 'warning',
            message: `Ci sono ${unmarked.length} alunni senza presenza o assenza assegnata.`
        })
        return
    }
    saving.value = true
    try {
        const classId = typeof selectedClass.value === 'object' ? selectedClass.value?.id : selectedClass.value
        const payload = {
            class_id: classId,
            date: date.value,
            hour: selectedHour.value,
            subject_id: selectedSubject.value || null, 
            statuses: students.value.map(s => ({
                student_id: s.id,
                status: s.status,
                entry_time: (s.status === 'Late' && s.entry_time && s.entry_time.trim() !== '') ? s.entry_time : null,
                exit_time: (s.status === 'LeftEarly' && s.exit_time && s.exit_time.trim() !== '') ? s.exit_time : null
            }))
        }
        
        await api.post('/attendance/mark-bulk', payload)
        
        $q.notify({ type: 'positive', message: 'Registro salvato con successo' })
    } catch (error) {
        $q.notify({ type: 'negative', message: 'Errore salvataggio' })
    } finally {
        saving.value = false
    }
}

const processJustification = async (id, approved) => {
    try {
        await api.post(`/attendance/justification/${id}/process`, { approve: approved })
        justificationRequests.value = justificationRequests.value.filter(r => r.id !== id)
        $q.notify({ 
            message: approved ? 'Giustificazione accettata' : 'Giustificazione respinta', 
            color: approved ? 'green' : 'orange' 
        })
    } catch (error) {
        $q.notify({ type: 'negative', message: 'Errore elaborazione' })
    }
}

const openNoteDialog = (student) => {
    selectedStudentForNote.value = student
    showNoteDialog.value = true
}

const exportCSV = async () => {
    if (!selectedClass.value) {
        $q.notify({ type: 'warning', message: 'Seleziona una classe prima di esportare' })
        return
    }
    try {
        const classId = typeof selectedClass.value === 'object' ? selectedClass.value.id : selectedClass.value
        const res = await api.get('/attendance/export', { params: { class_id: classId, date: date.value }, responseType: 'blob' })
        const url = window.URL.createObjectURL(new Blob([res.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `presenze_${classId}.csv`)
        document.body.appendChild(link)
        link.click()
        $q.notify({ type: 'positive', message: 'Export CSV completato!' })
    } catch (err) {
        $q.notify({ type: 'negative', message: 'Errore durante l\'export CSV' })
    }
}
</script>
