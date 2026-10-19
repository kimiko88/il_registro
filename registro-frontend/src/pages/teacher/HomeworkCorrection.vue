<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('homeworkSubmissions.canvasTitle') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('homeworkSubmissions.subtitle') }}
        </p>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-input
          v-model="homeworkIdFilter"
          outlined
          dense
          bg-color="white"
          placeholder="ID Compito (es. 101)"
          style="width: 180px"
          @keyup.enter="loadSubmissions"
        >
          <template #append>
            <q-btn flat round dense icon="search" @click="loadSubmissions" />
          </template>
        </q-input>
      </div>
    </div>

    <!-- Main Workspace: Split View with Submissions List & Interactive Canvas -->
    <div class="row q-col-gutter-lg">
      <!-- Left Column: Student Submissions -->
      <div class="col-12 col-md-4">
        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white full-height">
          <q-card-section>
            <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">
              Consegne Studenti ({{ submissions.length }})
            </div>

            <div v-if="submissions.length === 0" class="q-pa-lg text-center text-slate-400">
              Nessuna consegna trovata per questo compito.
            </div>

            <q-list separator v-else class="rounded-lg">
              <q-item
                v-for="sub in submissions"
                :key="sub.id"
                clickable
                v-ripple
                :active="selectedSubmission && selectedSubmission.id === sub.id"
                active-class="bg-indigo-50 text-primary text-weight-bold"
                class="rounded-lg q-my-xs"
                @click="selectSubmission(sub)"
              >
                <q-item-section avatar>
                  <q-avatar
                    :color="sub.status === 'returned' ? 'positive' : 'primary'"
                    text-color="white"
                    icon="assignment"
                  />
                </q-item-section>

                <q-item-section>
                  <q-item-label class="text-subtitle2">
                    Studente: {{ sub.student_id }}
                  </q-item-label>
                  <q-item-label caption>
                    {{ formatDate(sub.submitted_at) }}
                  </q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-badge :color="sub.status === 'returned' ? 'positive' : 'warning'">
                    {{ sub.status === 'returned' ? `Voto: ${sub.grade}` : 'Da correggere' }}
                  </q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Right Column: Canvas Annotation Studio -->
      <div class="col-12 col-md-8">
        <div v-if="!selectedSubmission" class="q-pa-xl text-center bg-white rounded-xl border-dashed border-slate-200">
          <q-icon name="draw" size="64px" color="slate-300" class="q-mb-sm" />
          <div class="text-h6 text-weight-bold text-slate-600">Seleziona una consegna per iniziare la correzione</div>
          <div class="text-body2 text-slate-400">
            Potrai annotare a mano libera o inserire testi direttamente sull'elaborato
          </div>
        </div>

        <q-card v-else class="rounded-xl shadow-soft border-slate-100 bg-white">
          <q-card-section>
            <!-- Canvas Toolbar -->
            <div class="row items-center justify-between q-mb-md bg-slate-100 q-pa-sm rounded-lg">
              <!-- Drawing Tools -->
              <div class="row q-gutter-xs items-center">
                <q-btn
                  dense
                  :color="activeTool === 'pen' ? 'primary' : 'grey-7'"
                  icon="edit"
                  label="Penna"
                  unelevated
                  class="rounded-md"
                  @click="activeTool = 'pen'"
                />
                <q-btn
                  dense
                  :color="activeTool === 'text' ? 'primary' : 'grey-7'"
                  icon="title"
                  label="Testo"
                  unelevated
                  class="rounded-md"
                  @click="activeTool = 'text'"
                />

                <q-separator vertical class="q-mx-sm" />

                <!-- Color Palette -->
                <div class="row q-gutter-xs items-center">
                  <div
                    v-for="c in ['#ef4444', '#3b82f6', '#10b981', '#0f172a']"
                    :key="c"
                    class="cursor-pointer rounded-full border-2"
                    :style="{
                      backgroundColor: c,
                      width: '24px',
                      height: '24px',
                      borderColor: currentColor === c ? '#000' : 'transparent'
                    }"
                    @click="currentColor = c"
                  />
                </div>

                <q-separator vertical class="q-mx-sm" />

                <!-- Stroke width -->
                <q-select
                  v-model="strokeWidth"
                  :options="[2, 4, 8]"
                  dense
                  outlined
                  bg-color="white"
                  style="width: 70px"
                />
              </div>

              <!-- Action buttons -->
              <div class="row q-gutter-xs">
                <q-btn
                  dense
                  flat
                  color="negative"
                  icon="delete_sweep"
                  label="Pulisci"
                  @click="clearCanvas"
                />
                <q-btn
                  dense
                  outline
                  color="secondary"
                  icon="cloud_sync"
                  label="Sync Classroom"
                  @click="syncWithClassroom"
                />
              </div>
            </div>

            <!-- Canvas Element -->
            <div class="canvas-container relative-position border border-slate-300 rounded-lg overflow-hidden bg-white shadow-inner flex justify-center items-center">
              <canvas
                ref="canvasRef"
                width="750"
                height="500"
                class="cursor-crosshair"
                @mousedown="startDrawing"
                @mousemove="draw"
                @mouseup="stopDrawing"
                @mouseleave="stopDrawing"
                @click="handleCanvasClick"
              />
            </div>

            <!-- Grading and Feedback Form -->
            <div class="row q-col-gutter-md q-mt-md">
              <div class="col-12 col-sm-4">
                <q-input
                  v-model.number="gradeForm.grade"
                  type="number"
                  step="0.5"
                  min="1"
                  max="10"
                  outlined
                  dense
                  label="Voto Elaborato (1-10) *"
                  :rules="[val => (val >= 1 && val <= 10) || 'Voto compreso tra 1 e 10']"
                />
              </div>
              <div class="col-12 col-sm-8">
                <q-input
                  v-model="gradeForm.teacher_feedback"
                  outlined
                  dense
                  label="Commento / Feedback Didattico"
                  placeholder="Es. Ottimo ragionamento logico, rivedere l'applicazione della formula..."
                />
              </div>
            </div>
          </q-card-section>

          <q-card-actions align="right" class="q-pa-md bg-slate-50 border-t border-slate-100">
            <q-btn
              color="positive"
              icon="send"
              :label="$t('homeworkSubmissions.returnBtn')"
              unelevated
              class="rounded-lg shadow-sm"
              :disable="!gradeForm.grade"
              @click="returnGradedHomework"
            />
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const homeworkIdFilter = ref('1')
const submissions = ref([])
const selectedSubmission = ref(null)

const canvasRef = ref(null)
const activeTool = ref('pen') // 'pen' | 'text'
const currentColor = ref('#ef4444')
const strokeWidth = ref(4)

const isDrawing = ref(false)
const drawnStrokes = ref([])
let ctx = null

const gradeForm = ref({
  grade: 8.0,
  teacher_feedback: ''
})

function formatDate(isoStr) {
  if (!isoStr) return '-'
  return new Date(isoStr).toLocaleString()
}

async function loadSubmissions() {
  try {
    const res = await api.get('/homework-submissions/by-homework', {
      params: { homework_id: homeworkIdFilter.value }
    })
    submissions.value = res.data?.submissions || []
    if (submissions.value.length > 0) {
      selectSubmission(submissions.value[0])
    } else {
      selectedSubmission.value = null
    }
  } catch (err) {
    console.error('Errore caricamento consegne:', err)
  }
}

function selectSubmission(sub) {
  selectedSubmission.value = sub
  gradeForm.value.grade = sub.grade || 8.0
  gradeForm.value.teacher_feedback = sub.teacher_feedback || ''
  nextTick(() => {
    initCanvas(sub)
  })
}

function initCanvas(sub) {
  const canvas = canvasRef.value
  if (!canvas) return
  ctx = canvas.getContext('2d')

  // Clear canvas background
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // Draw background mock document template or load existing annotations
  ctx.fillStyle = '#64748b'
  ctx.font = '14px Arial'
  ctx.fillText(`Elaborato dello Studente: ${sub.student_id}`, 20, 30)
  ctx.fillText(`File di origine: ${sub.file_url || 'compito.pdf'}`, 20, 50)
  ctx.strokeStyle = '#cbd5e1'
  ctx.beginPath()
  ctx.moveTo(20, 60)
  ctx.lineTo(730, 60)
  ctx.stroke()

  // Load existing annotations if present
  if (sub.canvas_annotations_json) {
    try {
      const strokes = JSON.parse(sub.canvas_annotations_json)
      drawnStrokes.value = strokes
      replayStrokes(strokes)
    } catch (e) {
      console.warn('Could not parse canvas annotations', e)
    }
  } else {
    drawnStrokes.value = []
  }
}

function replayStrokes(strokes) {
  if (!ctx || !strokes) return
  for (const s of strokes) {
    if (s.type === 'pen' && s.points?.length > 1) {
      ctx.strokeStyle = s.color
      ctx.lineWidth = s.width
      ctx.lineCap = 'round'
      ctx.lineJoin = 'round'
      ctx.beginPath()
      ctx.moveTo(s.points[0].x, s.points[0].y)
      for (let i = 1; i < s.points.length; i++) {
        ctx.lineTo(s.points[i].x, s.points[i].y)
      }
      ctx.stroke()
    } else if (s.type === 'text') {
      ctx.fillStyle = s.color
      ctx.font = '16px Arial'
      ctx.fillText(s.text, s.x, s.y)
    }
  }
}

let currentStroke = null

function startDrawing(e) {
  if (activeTool.value !== 'pen' || !ctx) return
  isDrawing.value = true
  const rect = canvasRef.value.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top

  currentStroke = {
    type: 'pen',
    color: currentColor.value,
    width: strokeWidth.value,
    points: [{ x, y }]
  }

  ctx.strokeStyle = currentColor.value
  ctx.lineWidth = strokeWidth.value
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.beginPath()
  ctx.moveTo(x, y)
}

function draw(e) {
  if (!isDrawing.value || !ctx || activeTool.value !== 'pen') return
  const rect = canvasRef.value.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top

  currentStroke.points.push({ x, y })
  ctx.lineTo(x, y)
  ctx.stroke()
}

function stopDrawing() {
  if (isDrawing.value && currentStroke) {
    drawnStrokes.value.push(currentStroke)
    currentStroke = null
  }
  isDrawing.value = false
}

function handleCanvasClick(e) {
  if (activeTool.value !== 'text' || !ctx) return
  const rect = canvasRef.value.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top

  $q.dialog({
    title: 'Inserisci Nota Didattica',
    message: 'Testo da stampare sul compito:',
    prompt: {
      model: '',
      type: 'text'
    },
    cancel: true,
    persistent: true
  }).onOk(data => {
    if (data && data.trim()) {
      ctx.fillStyle = currentColor.value
      ctx.font = '16px Arial'
      ctx.fillText(data, x, y)
      drawnStrokes.value.push({
        type: 'text',
        color: currentColor.value,
        text: data,
        x,
        y
      })
    }
  })
}

function clearCanvas() {
  if (!selectedSubmission.value) return
  initCanvas(selectedSubmission.value)
}

async function returnGradedHomework() {
  if (!selectedSubmission.value) return
  try {
    const serializedStrokes = JSON.stringify(drawnStrokes.value)
    await api.post(`/homework-submissions/${selectedSubmission.value.id}/grade`, {
      grade: Number(gradeForm.value.grade),
      teacher_feedback: gradeForm.value.teacher_feedback,
      canvas_annotations_json: serializedStrokes,
      annotated_file_url: `/annotated/${selectedSubmission.value.id}.pdf`
    })
    $q.notify({ type: 'positive', message: 'Elaborato corretto e restituito allo studente con successo!' })
    await loadSubmissions()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore durante la restituzione' })
  }
}

async function syncWithClassroom() {
  if (!selectedSubmission.value) return
  try {
    await api.post(`/homework-submissions/${selectedSubmission.value.id}/sync-classroom`, {
      google_classroom_submission_id: `gclass_${selectedSubmission.value.id}`
    })
    $q.notify({ type: 'positive', message: 'Voto sincronizzato con Google Classroom!' })
    await loadSubmissions()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore sincronizzazione Google Classroom' })
  }
}

onMounted(() => {
  loadSubmissions()
})
</script>

<style scoped>
.canvas-container {
  max-width: 100%;
}
</style>
