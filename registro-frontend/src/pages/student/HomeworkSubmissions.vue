<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          {{ $t('homeworkSubmissions.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          {{ $t('homeworkSubmissions.subtitle') }}
        </p>
      </div>

      <q-btn
        color="primary"
        icon="upload_file"
        :label="$t('homeworkSubmissions.uploadTitle')"
        class="rounded-lg shadow-sm"
        @click="showUploadDialog = true"
      />
    </div>

    <!-- Submissions List -->
    <div v-if="submissions.length === 0" class="q-pa-xl text-center bg-white rounded-xl border-dashed border-slate-200">
      <q-icon name="assignment" size="56px" color="slate-400" class="q-mb-sm" />
      <div class="text-h6 text-weight-bold text-slate-700">Nessun elaborato consegnato</div>
      <div class="text-body2 text-slate-400">
        Puoi caricare compiti svolti in formato PDF o immagine per la correzione da parte dei docenti.
      </div>
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="sub in submissions" :key="sub.id" class="col-12 col-md-6 col-lg-4">
        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white">
          <q-card-section>
            <div class="row items-center justify-between q-mb-xs">
              <div class="text-h6 text-weight-bold text-slate-800">
                Compito #{{ sub.homework_id }}
              </div>
              <q-badge
                :color="sub.status === 'returned' ? 'positive' : 'info'"
                class="q-pa-xs text-weight-bold"
              >
                {{ sub.status === 'returned' ? $t('homeworkSubmissions.statusReturned') : $t('homeworkSubmissions.statusSubmitted') }}
              </q-badge>
            </div>

            <div class="text-caption text-slate-500 q-mb-sm">
              Consegnato il: {{ formatDate(sub.submitted_at) }}
            </div>

            <div v-if="sub.notes" class="text-body2 text-slate-700 q-mb-sm bg-slate-50 q-pa-sm rounded-lg">
              <span class="text-weight-bold">Note studente:</span> {{ sub.notes }}
            </div>

            <!-- Correction section if graded -->
            <div v-if="sub.status === 'returned'" class="q-pa-sm rounded-lg bg-emerald-50 border border-emerald-200 q-mb-sm">
              <div class="row items-center justify-between q-mb-xs">
                <span class="text-caption text-emerald-800 text-weight-bold">Valutazione:</span>
                <span class="text-h6 text-weight-bolder text-emerald-900">{{ sub.grade }} / 10</span>
              </div>
              <div v-if="sub.teacher_feedback" class="text-caption text-emerald-900 italic">
                "{{ sub.teacher_feedback }}"
              </div>
            </div>

            <!-- Classroom sync badge -->
            <div v-if="sub.is_synced_classroom" class="row items-center q-gutter-x-xs text-caption text-indigo-7 q-mb-sm">
              <q-icon name="cloud_done" size="xs" />
              <span>{{ $t('homeworkSubmissions.classroomSynced') }}</span>
            </div>

            <div class="row q-gutter-sm q-mt-xs">
              <q-btn
                v-if="sub.file_url"
                outline
                dense
                color="primary"
                icon="download"
                label="Mio File"
                type="a"
                :href="sub.file_url"
                target="_blank"
                class="rounded-lg"
              />
              <q-btn
                v-if="sub.annotated_file_url"
                unelevated
                dense
                color="secondary"
                icon="brush"
                label="Correzione Docente"
                type="a"
                :href="sub.annotated_file_url"
                target="_blank"
                class="rounded-lg"
              />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Dialog: Upload Submission -->
    <q-dialog v-model="showUploadDialog" persistent>
      <q-card style="min-width: 440px" class="rounded-xl">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6 text-weight-bold text-slate-800">
            {{ $t('homeworkSubmissions.uploadTitle') }}
          </div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="newSubmission.homework_id"
            outlined
            dense
            label="Identificativo Compito / Argomento *"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <q-input
            v-model="newSubmission.file_url"
            outlined
            dense
            label="URL Elaborato (PDF / Scansione)"
            hint="Es. /uploads/compito_mate_rossi.pdf"
            :rules="[val => !!val || 'Inserire l\'URL dell\'elaborato']"
          />

          <q-input
            v-model="newSubmission.notes"
            outlined
            dense
            type="textarea"
            rows="3"
            label="Note o chiarimenti per il docente"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="$t('common.cancel')" v-close-popup />
          <q-btn
            color="primary"
            icon="send"
            :label="$t('homeworkSubmissions.submitBtn')"
            :disable="!newSubmission.homework_id || !newSubmission.file_url"
            @click="submitHomework"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const submissions = ref([])
const showUploadDialog = ref(false)

const newSubmission = ref({
  homework_id: '',
  file_url: '',
  notes: ''
})

function formatDate(isoStr) {
  if (!isoStr) return '-'
  return new Date(isoStr).toLocaleString()
}

async function loadSubmissions() {
  try {
    const res = await api.get('/homework-submissions/my-submissions')
    submissions.value = res.data?.submissions || []
  } catch (err) {
    console.error('Errore caricamento compiti consegnati:', err)
  }
}

async function submitHomework() {
  try {
    await api.post('/homework-submissions', newSubmission.value)
    $q.notify({ type: 'positive', message: 'Elaborato consegnato al docente!' })
    showUploadDialog.value = false
    newSubmission.value = { homework_id: '', file_url: '', notes: '' }
    await loadSubmissions()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore durante la consegna' })
  }
}

onMounted(() => {
  loadSubmissions()
})
</script>
