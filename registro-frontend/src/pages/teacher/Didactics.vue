<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="folder_shared" color="primary" class="q-mr-sm" />
          {{ t('didacticsPage.title') }}
        </div>
        <div class="text-caption text-grey">{{ t('didacticsPage.subtitle') }}</div>
      </div>
      <q-btn icon="cloud_upload" :label="t('didacticsPage.uploadMaterial')" color="primary" class="shadow-soft" @click="openUploadDialog" />
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md shadow-1">
      <q-card-section class="row items-center q-gutter-md q-py-sm">
        <q-select
          v-model="selectedClass"
          :options="classOptions"
          option-value="id"
          option-label="label"
          emit-value map-options
          :label="t('common.filter') + ' Classe'"
          dense outlined
          style="min-width:180px"
        />
        <q-select
          v-model="selectedSubject"
          :options="availableSubjectOptions"
          option-value="subject_id"
          option-label="subject_name"
          emit-value map-options
          :label="t('didacticsPage.filterSubject')"
          dense outlined
          style="min-width:160px"
          clearable
        />
      </q-card-section>
    </q-card>

    <!-- Materials List -->
    <q-card v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="60px" />
    </q-card>

    <div v-else>
      <q-card v-if="filteredMaterials.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="folder_open" size="80px" class="q-mb-md" />
        <div class="text-h6">Nessun materiale condiviso</div>
        <div class="text-caption">Aggiungi il primo file per iniziare la condivisione</div>
      </q-card>

      <div v-else class="row q-col-gutter-md">
        <div v-for="mat in filteredMaterials" :key="mat.id" class="col-12 col-sm-6 col-md-4">
          <q-card class="shadow-soft hover-card full-height flex flex-column justify-between">
            <q-card-section>
              <div class="row items-center justify-between q-mb-xs">
                <q-chip dense size="sm" color="primary-1" text-color="primary" class="text-weight-bold">
                  {{ getSubjectName(mat.subject_id) }}
                </q-chip>
                <div class="text-caption text-grey">{{ formatDate(mat.created_at) }}</div>
              </div>
              <div class="text-h6 text-weight-bold text-outfit q-my-sm">{{ mat.title }}</div>
              <div class="text-body2 text-grey-8 text-justify line-clamp-3" v-if="mat.description">
                {{ mat.description }}
              </div>
              <div class="text-caption text-grey-6 q-mt-sm">Docente: {{ mat.teacher_name }}</div>
            </q-card-section>

            <q-card-actions class="q-pa-md border-top row justify-between items-center">
              <q-btn
                v-if="mat.attachment_url"
                flat color="primary"
                icon="open_in_new"
                label="Apri Allegato"
                no-caps
                @click="openLink(mat.attachment_url)"
              />
              <span v-else class="text-caption text-grey-5">Nessun file allegato</span>
              
              <q-btn
                flat color="negative"
                icon="delete"
                dense
                @click="confirmDelete(mat)"
              >
                <q-tooltip>Elimina materiale</q-tooltip>
              </q-btn>
            </q-card-actions>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Upload Dialog -->
    <q-dialog v-model="uploadDialog" persistent>
      <q-card style="width: min(500px, 95vw); max-width: 95vw;" class="q-pa-sm">
        <q-card-section class="row items-center">
          <div class="text-h6 text-weight-bold text-outfit">
            <q-icon name="cloud_upload" color="primary" class="q-mr-xs" />
            Condividi Materiale
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-none">
          <q-select
            v-model="newMaterial.subject_id"
            :options="gradesStore.subjects"
            option-value="subject_id"
            option-label="subject_name"
            emit-value map-options
            label="Materia *"
            outlined dense
            :rules="[v => !!v || 'Campo obbligatorio']"
          />

          <q-input
            v-model="newMaterial.title"
            label="Titolo *"
            outlined dense
            placeholder="Es: Dispensa sulle equazioni"
            :rules="[v => !!v || 'Campo obbligatorio']"
          />

          <q-input
            v-model="newMaterial.description"
            label="Descrizione / Istruzioni"
            type="textarea"
            outlined dense
            autogrow
            placeholder="Fornisci dettagli o spiegazioni (opzionale)"
          />

          <!-- Upload option: simple URL link or Simulated File Upload -->
          <div class="bg-grey-2 q-pa-md rounded-lg">
            <div class="text-subtitle2 text-weight-medium q-mb-sm">Allegato</div>
            <q-option-group
              v-model="uploadType"
              :options="[
                { label: 'Link Esterno / Condivisione Cloud', value: 'url' },
                { label: 'Carica File (Simulato)', value: 'file' }
              ]"
              color="primary"
              inline
              class="q-mb-md"
            />

            <q-input
              v-if="uploadType === 'url'"
              v-model="newMaterial.attachment_url"
              label="Link URL (es: Google Drive, Dropbox, YouTube)"
              outlined dense
              placeholder="https://drive.google.com/..."
            />

            <q-file
              v-else
              v-model="uploadedFile"
              label="Seleziona File"
              outlined dense
              counter
              use-chips
            >
              <template v-slot:prepend>
                <q-icon name="attach_file" />
              </template>
            </q-file>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pb-md q-pr-md">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Salva e Condividi" :loading="saving" @click="saveMaterial" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useClassesStore } from '@/stores/classes'
import { useGradesStore } from '@/stores/grades'
import didacticService from '@/services/didacticService'
import api from '@/services/api'

const $q = useQuasar()
const { t } = useI18n()
const authStore = useAuthStore()
const gradesStore = useGradesStore()
const classesStore = useClassesStore()

const selectedClass = ref(null)
const selectedSubject = ref(null)
const loading = ref(true)
const saving = ref(false)
const uploadDialog = ref(false)
const uploadType = ref('url')
const uploadedFile = ref(null)

const materials = ref([])
const classOptions = ref([])

const isCivicaSubject = (s) => {
  const name = (s.subject_name || s.name || '').toLowerCase()
  return name.includes('civica') || name.includes('educazione civica') || name.includes('ed. civica')
}

const isAssignedToCurrentTeacher = (s, user) => {
  if (!user) return false
  const currentUserId = String(user.id || '')
  const teacherId = user.teacher_id ? String(user.teacher_id) : ''
  const sTeacherId = s.teacher_id ? String(s.teacher_id) : ''
  const sTeacherUserId = s.teacher_user_id ? String(s.teacher_user_id) : ''

  if (sTeacherId && (sTeacherId === currentUserId || (teacherId && sTeacherId === teacherId))) {
    return true
  }
  if (sTeacherUserId && (sTeacherUserId === currentUserId || (teacherId && sTeacherUserId === teacherId))) {
    return true
  }
  if (s.teacher_name && user.last_name) {
    const tName = s.teacher_name.toLowerCase()
    const uLast = user.last_name.toLowerCase()
    const uFirst = (user.first_name || '').toLowerCase()
    if (tName.includes(uLast) && (!uFirst || tName.includes(uFirst))) {
      return true
    }
  }
  return false
}

const availableSubjectOptions = computed(() => {
  const allSubjects = gradesStore.subjects || []
  const user = authStore.user
  if (!user || ['admin', 'superadmin', 'secretary'].includes(user.role)) {
    return allSubjects
  }

  // Mostra ESCLUSIVAMENTE le materie assegnate dalla segreteria al docente loggato + Educazione Civica
  return allSubjects.filter(s => {
    return isAssignedToCurrentTeacher(s, user) || isCivicaSubject(s)
  })
})

const newMaterial = ref({
  subject_id: null,
  title: '',
  description: '',
  attachment_url: ''
})

onMounted(async () => {
  await classesStore.fetchAssignedClasses()
  classOptions.value = classesStore.classes || []
  if (classOptions.value.length > 0) {
    selectedClass.value = classOptions.value[0].id
  }
})

watch(selectedClass, async () => {
  if (selectedClass.value) {
    await gradesStore.fetchClassSubjects(selectedClass.value)
    if (availableSubjectOptions.value && availableSubjectOptions.value.length > 0) {
      selectedSubject.value = availableSubjectOptions.value[0].subject_id
    } else {
      selectedSubject.value = null
    }
  } else {
    selectedSubject.value = null
  }
  fetchMaterials()
})

const fetchMaterials = async () => {
  if (!selectedClass.value) return
  loading.value = true
  try {
    const res = await didacticService.getMaterials(selectedClass.value)
    materials.value = res.data || []
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento del materiale' })
  } finally {
    loading.value = false
  }
}

const filteredMaterials = computed(() => {
  if (!selectedSubject.value) return materials.value
  return materials.value.filter(m => m.subject_id === selectedSubject.value)
})

const getSubjectName = (id) => {
  const s = gradesStore.subjects.find(s => s.subject_id === id)
  return s ? s.subject_name : id
}

const openUploadDialog = () => {
  newMaterial.value = {
    subject_id: selectedSubject.value,
    title: '',
    description: '',
    attachment_url: ''
  }
  uploadedFile.value = null
  uploadType.value = 'url'
  uploadDialog.value = true
}

const saveMaterial = async () => {
  if (!newMaterial.value.title || !newMaterial.value.subject_id) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }
  saving.value = true
  try {
    let url = newMaterial.value.attachment_url
    if (uploadType.value === 'file' && uploadedFile.value) {
      // Simulate file upload: set a mock download link
      url = `https://mockstorage.registro-elettronico.it/uploads/${uploadedFile.value.name}`
    }

    await didacticService.createMaterial({
      class_id: selectedClass.value,
      subject_id: newMaterial.value.subject_id,
      title: newMaterial.value.title,
      description: newMaterial.value.description,
      attachment_url: url
    })

    $q.notify({ type: 'positive', message: 'Materiale condiviso con successo!' })
    uploadDialog.value = false
    fetchMaterials()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Impossibile salvare il materiale' })
  } finally {
    saving.value = false
  }
}

const confirmDelete = (mat) => {
  $q.dialog({
    title: 'Conferma eliminazione',
    message: `Sei sicuro di voler eliminare "${mat.title}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await didacticService.deleteMaterial(mat.id)
      $q.notify({ type: 'positive', message: 'Materiale rimosso con successo' })
      fetchMaterials()
    } catch (e) {
      $q.notify({ type: 'negative', message: 'Impossibile rimuovere il materiale' })
    }
  })
}

const openLink = (url) => {
  if (url) window.open(url, '_blank')
}

const formatDate = (d) => date.formatDate(new Date(d), 'DD/MM/YYYY')
</script>

<style scoped>
.hover-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.hover-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}
.border-top {
  border-top: 1px solid rgba(0,0,0,0.08);
}
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;  
  overflow: hidden;
}
</style>
