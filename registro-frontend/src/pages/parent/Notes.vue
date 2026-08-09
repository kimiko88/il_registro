<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="assignment_late" color="negative" class="q-mr-sm" />
          Note e Provvedimenti
        </div>
        <div class="text-caption text-grey">
          Registro delle note disciplinari, note di comportamento o richiami compiti di 
          <strong>{{ selectedChild ? `${selectedChild.first_name} ${selectedChild.last_name}` : '...' }}</strong>
        </div>
      </div>
    </div>

    <!-- Child Warning -->
    <q-card v-if="!selectedChild" class="text-center q-pa-lg bg-warning text-white">
      Seleziona un figlio dal menu in alto per visualizzare le note disciplinari
    </q-card>

    <div v-else>
      <!-- Notes List -->
      <q-card v-if="loading" class="text-center q-pa-xl shadow-1">
        <q-spinner-dots color="primary" size="60px" />
      </q-card>

      <div v-else>
        <q-card v-if="notes.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
          <q-icon name="check_circle_outline" color="positive" size="80px" class="q-mb-md" />
          <div class="text-h6 text-positive">Nessuna nota registrata</div>
          <div class="text-caption">Ottimo! Tuo figlio non ha ricevuto alcuna nota disciplinare o comportamentale.</div>
        </q-card>

        <div v-else class="q-gutter-md">
          <q-card v-for="note in notes" :key="note.id" class="shadow-1 hover-card border-left-indicator" :class="getNoteBorderClass(note.type)">
            <q-card-section>
              <div class="row items-center justify-between q-mb-sm">
                <div class="row items-center q-gutter-sm">
                  <q-chip dense :color="getNoteColor(note.type)" text-color="white" class="text-weight-bold">
                    {{ formatNoteType(note.type) }}
                  </q-chip>
                  <q-chip v-if="note.type === 'disciplinary'" dense color="positive" text-color="white" icon="verified">
                    Approvata Dirigenza
                  </q-chip>
                  <q-chip dense color="info" text-color="white" icon="visibility">
                    Letta
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
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useQuasar, date } from 'quasar'
import { storeToRefs } from 'pinia'
import { useParentStore } from 'src/stores/parent'
import notesService from 'src/services/notesService'

const $q = useQuasar()
const parentStore = useParentStore()
const { selectedChild } = storeToRefs(parentStore)

const loading = ref(true)
const notes = ref([])

onMounted(async () => {
  if (selectedChild.value) {
    await fetchNotes()
  } else {
    loading.value = false
  }
})

watch(selectedChild, async (newVal) => {
  if (newVal) {
    loading.value = true
    await fetchNotes()
    loading.value = false
  }
})

const fetchNotes = async () => {
  try {
    const studentId = selectedChild.value?.id
    if (!studentId) return
    const res = await notesService.getNotes({ student_id: studentId })
    notes.value = res.data || []
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: 'Impossibile caricare le note' })
  } finally {
    loading.value = false
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
