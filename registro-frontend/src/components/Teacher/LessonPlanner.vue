<template>
  <q-page class="q-pa-md bg-grey-1">

    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">Registro di Classe</div>
        <div class="text-caption text-grey">Firma il registro, inserisci l'argomento e assegna i compiti</div>
      </div>
      <q-btn icon="add" label="Nuova Lezione" color="primary" @click="openNewLesson" />
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md shadow-1">
      <q-card-section class="row items-center q-gutter-md q-py-sm">
        <q-select
          v-model="selectedClass"
          :options="classOptions"
          option-value="id"
          option-label="name"
          emit-value map-options
          label="Classe"
          dense outlined
          style="min-width:140px"
        />
        <q-select
          v-model="selectedSubject"
          :options="subjectOptions"
          option-value="id"
          option-label="name"
          emit-value map-options
          label="Materia"
          dense outlined
          style="min-width:160px"
          clearable
        />
        <q-space />
        <q-btn-toggle
          v-model="activeTab"
          toggle-color="primary"
          flat
          :options="[{label:'Lezioni', value:'lessons'},{label:'Compiti', value:'homework'}]"
        />
      </q-card-section>
    </q-card>

    <!-- Lessons List -->
    <div v-if="activeTab === 'lessons'">
      <q-card v-if="lessons.length === 0" class="text-center q-pa-xl text-grey-6">
        <q-icon name="menu_book" size="80px" class="q-mb-md" />
        <div class="text-h6">Nessuna lezione registrata</div>
        <div class="text-caption">Aggiungi la prima lezione per iniziare il registro</div>
      </q-card>
      <q-card v-else class="shadow-1">
        <q-list separator>
          <q-item v-for="lesson in lessons" :key="lesson.id" class="q-py-md">
            <q-item-section avatar>
              <q-avatar color="primary" text-color="white" icon="menu_book" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-bold">{{ lesson.topic }}</q-item-label>
              <q-item-label caption>
                {{ formatDate(lesson.date) }} &bull;
                <q-badge outline :color="getLessonTypeColor(lesson.type)">{{ lesson.type }}</q-badge>
              </q-item-label>
              <q-item-label caption class="text-grey-7" v-if="lesson.notes">{{ lesson.notes }}</q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="text-caption text-grey">{{ lesson.subject_id }}</div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>
    </div>

    <!-- Homework List -->
    <div v-if="activeTab === 'homework'">
      <q-card v-if="homeworks.length === 0" class="text-center q-pa-xl text-grey-6">
        <q-icon name="assignment" size="80px" class="q-mb-md" />
        <div class="text-h6">Nessun compito assegnato</div>
        <q-btn label="Assegna Compito" color="primary" @click="openNewHomework" class="q-mt-md" />
      </q-card>
      <q-card v-else class="shadow-1">
        <q-card-section class="row justify-end q-pb-none">
          <q-btn flat icon="add" label="Assegna Compito" color="primary" @click="openNewHomework" />
        </q-card-section>
        <q-list separator>
          <q-item v-for="hw in homeworks" :key="hw.id" class="q-py-md">
            <q-item-section avatar>
              <q-avatar color="orange" text-color="white" icon="assignment" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-bold">{{ hw.description }}</q-item-label>
              <q-item-label caption>
                Consegna: <strong>{{ formatDate(hw.due_date) }}</strong>
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>
    </div>

    <!-- New Lesson Dialog -->
    <q-dialog v-model="lessonDialog" persistent>
      <q-card style="min-width: 480px">
        <q-card-section>
          <div class="text-h6">
            <q-icon name="menu_book" class="q-mr-xs" color="primary" />
            Nuova Lezione
          </div>
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input
            v-model="newLesson.date"
            type="date"
            label="Data *"
            outlined dense
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
          <q-input
            v-model="newLesson.topic"
            label="Argomento della Lezione *"
            outlined dense
            placeholder="Es: Le equazioni di secondo grado"
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-md-6">
              <q-input
                v-model.number="newLesson.hour"
                type="number"
                label="Quale ora *"
                outlined dense
                min="1" max="10"
                :rules="[v => !!v || 'Campo obbligatorio']"
              />
            </div>
            <div class="col-12 col-md-6">
              <q-input
                v-model.number="newLesson.duration"
                type="number"
                label="Quante ore *"
                outlined dense
                min="1" max="5"
                :rules="[v => !!v || 'Campo obbligatorio']"
              />
            </div>
          </div>
          <q-select
            v-model="newLesson.type"
            :options="['Frontale', 'Laboratorio', 'Verifica', 'Discussione', 'Lavoro di gruppo', 'Altro']"
            label="Tipo di Lezione *"
            outlined dense
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
          <q-input
            v-model="newLesson.notes"
            label="Note / Osservazioni interne"
            type="textarea"
            outlined dense
            autogrow
            placeholder="(Opzionale)"
          />
          <q-toggle
            v-model="assignHomeworkToo"
            label="Assegna anche un compito per questa lezione"
            color="orange"
          />
          <template v-if="assignHomeworkToo">
            <q-input
              v-model="newLesson.homeworkDesc"
              label="Descrizione Compito *"
              type="textarea"
              outlined dense
              autogrow
              placeholder="Es: Pagine 145-150, esercizi 12-18"
            />
            <q-input
              v-model="newLesson.homeworkDue"
              type="date"
              label="Data Consegna *"
              outlined dense
            />
          </template>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="primary" label="Salva" :loading="saving" @click="saveLesson" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- New Homework Dialog -->
    <q-dialog v-model="homeworkDialog" persistent>
      <q-card style="min-width: 420px">
        <q-card-section>
          <div class="text-h6">
            <q-icon name="assignment" class="q-mr-xs" color="orange" />
            Assegna Compito
          </div>
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-input
            v-model="newHomework.description"
            label="Descrizione Compito *"
            type="textarea"
            outlined dense
            autogrow
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
          <q-input
            v-model="newHomework.dueDate"
            type="date"
            label="Data di Consegna *"
            outlined dense
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn color="orange" text-color="white" label="Assegna" :loading="saving" @click="saveHomework" />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useQuasar, date } from 'quasar'
import { useClassesStore } from 'src/stores/classes'
import { lessonService } from 'src/services/lessonService'

const $q = useQuasar()
const classesStore = useClassesStore()

const selectedClass = ref(null)
const selectedSubject = ref(null)
const activeTab = ref('lessons')
const lessonDialog = ref(false)
const homeworkDialog = ref(false)
const saving = ref(false)
const assignHomeworkToo = ref(false)

const lessons = ref([])
const homeworks = ref([])

// Options from the store
const classOptions = ref([])
const subjectOptions = ref([])

const newLesson = ref({
  date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
  hour: 1,
  duration: 1,
  topic: '',
  type: 'Frontale',
  notes: '',
  homeworkDesc: '',
  homeworkDue: ''
})

const newHomework = ref({
  description: '',
  dueDate: date.formatDate(Date.now(), 'YYYY-MM-DD')
})

onMounted(async () => {
  await classesStore.fetchAssignedClasses()
  classOptions.value = classesStore.classes || []
  if (classOptions.value.length > 0) {
    selectedClass.value = classOptions.value[0].id
  }
})

watch(selectedClass, () => {
  fetchData()
})

watch(selectedSubject, () => {
  if (activeTab.value === 'lessons') fetchLessons()
})

const fetchData = async () => {
  await Promise.all([fetchLessons(), fetchHomeworks()])
}

const fetchLessons = async () => {
  if (!selectedClass.value) return
  try {
    const res = await lessonService.getLessons(selectedClass.value, selectedSubject.value)
    lessons.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const fetchHomeworks = async () => {
  if (!selectedClass.value) return
  try {
    const res = await lessonService.getHomeworks(selectedClass.value)
    homeworks.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const openNewLesson = () => {
  newLesson.value = {
    date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
    hour: 1,
    duration: 1,
    topic: '',
    type: 'Frontale',
    notes: '',
    homeworkDesc: '',
    homeworkDue: ''
  }
  assignHomeworkToo.value = false
  lessonDialog.value = true
}

const openNewHomework = () => {
  newHomework.value = {
    description: '',
    dueDate: date.formatDate(Date.now(), 'YYYY-MM-DD')
  }
  homeworkDialog.value = true
}

const saveLesson = async () => {
  if (!newLesson.value.topic || !newLesson.value.date || !newLesson.value.type) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }
  saving.value = true
  try {
    const lessonPayload = {
      class_id: selectedClass.value,
      subject_id: selectedSubject.value,
      date: newLesson.value.date,
      hour: newLesson.value.hour,
      duration: newLesson.value.duration,
      topic: newLesson.value.topic,
      type: newLesson.value.type,
      notes: newLesson.value.notes
    }
    const lessonRes = await lessonService.createLesson(lessonPayload)

    // If homework is also assigned
    if (assignHomeworkToo.value && newLesson.value.homeworkDesc) {
      await lessonService.createHomework({
        class_id: selectedClass.value,
        subject_id: selectedSubject.value,
        lesson_id: lessonRes.data?.id,
        due_date: newLesson.value.homeworkDue,
        description: newLesson.value.homeworkDesc
      })
    }

    $q.notify({ type: 'positive', message: 'Lezione registrata con successo' })
    lessonDialog.value = false
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: `Errore: ${e.response?.data?.error || e.message}` })
  } finally {
    saving.value = false
  }
}

const saveHomework = async () => {
  if (!newHomework.value.description || !newHomework.value.dueDate) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }
  saving.value = true
  try {
    await lessonService.createHomework({
      class_id: selectedClass.value,
      subject_id: selectedSubject.value,
      due_date: newHomework.value.dueDate,
      description: newHomework.value.description
    })
    $q.notify({ type: 'positive', message: 'Compito assegnato con successo' })
    homeworkDialog.value = false
    await fetchHomeworks()
  } catch (e) {
    $q.notify({ type: 'negative', message: `Errore: ${e.response?.data?.error || e.message}` })
  } finally {
    saving.value = false
  }
}

const formatDate = (d) => {
  if (!d) return '-'
  return date.formatDate(new Date(d), 'DD/MM/YYYY')
}

const getLessonTypeColor = (type) => {
  const map = {
    'Frontale': 'primary',
    'Laboratorio': 'teal',
    'Verifica': 'deep-orange',
    'Discussione': 'purple',
    'Lavoro di gruppo': 'cyan',
  }
  return map[type] || 'grey'
}
</script>
