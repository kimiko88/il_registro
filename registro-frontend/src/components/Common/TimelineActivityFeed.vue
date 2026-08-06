<template>
  <q-card class="rounded-xl shadow-xs border bg-white">
    <q-card-section class="row items-center justify-between q-pb-none">
      <div class="text-h6 text-weight-bold text-slate-800 row items-center">
        <q-icon name="dynamic_feed" color="primary" class="q-mr-sm" />
        Timeline del Giorno & Attività
      </div>
      <q-badge color="primary" class="text-weight-bold">Oggi</q-badge>
    </q-card-section>

    <q-card-section>
      <div v-if="loading" class="q-pa-md text-center text-slate-400">
        <q-spinner color="primary" size="32px" />
      </div>

      <div v-else-if="events.length === 0" class="text-center text-slate-400 q-pa-xl">
        <q-icon name="event_available" size="48px" class="q-mb-xs text-grey-4" /><br />
        <div class="text-weight-bold text-slate-600">Nessun evento registrato oggi</div>
        <div class="text-caption text-grey-5">Non ci sono nuovi voti, assenze o note disciplinari registrate per la giornata odierna.</div>
      </div>

      <q-timeline v-else color="primary" class="q-px-sm">
        <q-timeline-entry
          v-for="evt in events"
          :key="evt.id"
          :title="evt.title"
          :subtitle="evt.time"
          :color="getEventColor(evt.type)"
          :icon="getEventIcon(evt.type)"
        >
          <div class="text-body2 text-slate-700 bg-slate-50 q-pa-sm rounded-lg border">
            {{ evt.description }}
            <div v-if="evt.badge" class="q-mt-xs">
              <q-chip dense :color="getEventColor(evt.type)" text-color="white" size="xs font-bold">
                {{ evt.badge }}
              </q-chip>
            </div>
          </div>
        </q-timeline-entry>
      </q-timeline>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/services/api'

const props = defineProps({
  studentId: {
    type: String,
    default: ''
  }
})

const loading = ref(true)
const events = ref([])

onMounted(async () => {
  try {
    const params = props.studentId ? { student_id: props.studentId } : {}
    const [gradesRes, attendanceRes, notesRes] = await Promise.allSettled([
      api.get('/grades/my', { params }),
      api.get('/attendance', { params }),
      api.get('/notes', { params })
    ])

    const list = []

    // 1. Process recent grades
    if (gradesRes.status === 'fulfilled') {
      const grades = gradesRes.value.data?.grades || gradesRes.value.data || []
      grades.slice(0, 3).forEach(g => {
        list.push({
          id: `grade-${g.id}`,
          type: 'grade',
          title: `Valutazione in ${g.subject_name || 'Materia'}`,
          subtitle: formatDate(g.date || g.created_at),
          time: 'Voto Registrato',
          description: g.description || 'Valutazione periodica',
          badge: `Voto: ${g.grade_value}`
        })
      })
    }

    // 2. Process attendance
    if (attendanceRes.status === 'fulfilled') {
      const atts = attendanceRes.value.data || []
      atts.slice(0, 2).forEach(a => {
        if (a.status === 'absent' || a.status === 'late') {
          list.push({
            id: `att-${a.id}`,
            type: 'attendance',
            title: a.status === 'absent' ? 'Assenza Registrata' : 'Ritardo in Ingresso',
            subtitle: formatDate(a.date),
            time: 'Registro Presenze',
            description: a.justification_reason || (a.status === 'absent' ? 'In attesa di giustificazione' : 'Ingresso in ritardo'),
            badge: a.status === 'absent' ? 'Assente' : 'Ritardo'
          })
        }
      })
    }

    // 3. Process notes
    if (notesRes.status === 'fulfilled') {
      const notes = notesRes.value.data || []
      notes.slice(0, 2).forEach(n => {
        list.push({
          id: `note-${n.id}`,
          type: 'note',
          title: 'Nota Disciplinare',
          subtitle: formatDate(n.created_at),
          time: 'Annotazione Docente',
          description: n.description || n.note_text,
          badge: 'Nota'
        })
      })
    }

    events.value = list
  } catch (err) {
    events.value = []
  } finally {
    loading.value = false
  }
})

function getEventColor(type) {
  switch (type) {
    case 'grade': return 'indigo'
    case 'attendance': return 'orange-8'
    case 'note': return 'red-7'
    default: return 'primary'
  }
}

function getEventIcon(type) {
  switch (type) {
    case 'grade': return 'grade'
    case 'attendance': return 'how_to_reg'
    case 'note': return 'warning'
    default: return 'event'
  }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('it-IT', { day: '2-digit', month: 'short' })
}
</script>
