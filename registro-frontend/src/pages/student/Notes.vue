<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="assignment_late" color="negative" class="q-mr-sm" />
          Note e Provvedimenti
        </div>
        <div class="text-caption text-grey">Registro delle note disciplinari, note di comportamento o richiami compiti</div>
      </div>
    </div>

    <!-- Notes List -->
    <q-card v-if="loading" class="text-center q-pa-xl shadow-1">
      <q-spinner-dots color="primary" size="60px" />
    </q-card>

    <div v-else>
      <q-card v-if="notes.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="check_circle_outline" color="positive" size="80px" class="q-mb-md" />
        <div class="text-h6 text-positive">Nessuna nota registrata</div>
        <div class="text-caption">Ottimo lavoro! Non hai ricevuto alcun richiamo disciplinare o comportamentale.</div>
      </q-card>

      <div v-else class="q-gutter-md">
        <q-card v-for="note in notes" :key="note.id" class="shadow-1 hover-card border-left-indicator" :class="getNoteBorderClass(note.type)">
          <q-card-section>
            <div class="row items-center justify-between q-mb-sm">
              <div class="row items-center q-gutter-sm">
                <q-chip dense :color="getNoteColor(note.type)" text-color="white" class="text-weight-bold">
                  {{ formatNoteType(note.type) }}
                </q-chip>
                <div class="text-caption text-grey">Docente: {{ note.teacher_name }}</div>
              </div>
              <div class="text-caption text-grey">{{ formatDate(note.date) }}</div>
            </div>
            
            <div class="text-body1 text-grey-9 text-weight-medium q-my-sm">
              {{ note.note }}
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar, date } from 'quasar'
import { useStudentStore } from '@/stores/student'
import notesService from '@/services/notesService'

const $q = useQuasar()
const { t } = useI18n()
const studentStore = useStudentStore()
const loading = ref(true)
const notes = ref([])

onMounted(async () => {
  await studentStore.fetchProfile()
  await fetchNotes()
  loading.value = false
})

const fetchNotes = async () => {
  try {
    const studentId = studentStore.profile?.id
    if (!studentId) return
    const res = await notesService.getNotes({ student_id: studentId })
    notes.value = res.data || []
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: 'Impossibile caricare le note' })
  }
}

const formatNoteType = (type) => {
  const map = {
    generic: 'Generica',
    homework: 'Compiti',
    behavior: 'Comportamentale',
    disciplinary: 'Disciplinare'
  }
  return map[type] || type
}

const getNoteColor = (type) => {
  const map = {
    generic: 'grey-7',
    homework: 'orange-8',
    behavior: 'deep-orange-9',
    disciplinary: 'red-9'
  }
  return map[type] || 'grey'
}

const getNoteBorderClass = (type) => {
  return `border-${type}`
}

const formatDate = (d) => date.formatDate(new Date(d), 'DD/MM/YYYY')
</script>

<style scoped>
.hover-card {
  transition: all 0.2s ease-in-out;
}
.hover-card:hover {
  transform: translateX(4px);
}
.border-left-indicator {
  border-left: 6px solid #ccc;
}
.border-generic {
  border-left-color: #9e9e9e;
}
.border-homework {
  border-left-color: #ff9800;
}
.border-behavior {
  border-left-color: #e64a19;
}
.border-disciplinary {
  border-left-color: #d32f2f;
}
</style>
