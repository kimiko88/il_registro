<template>
  <q-page class="q-pa-md bg-grey-1">

    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold row items-center">
          <span>Registro di Classe</span>
          <q-chip v-if="isSubstitutionMode" color="deep-orange" text-color="white" class="q-ml-sm text-weight-bold">
            <q-icon name="swap_horiz" class="q-mr-xs" />
            Modalità Supplenza
          </q-chip>
        </div>
        <div class="text-caption text-grey">Firma il registro, inserisci l'argomento e assegna i compiti</div>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn
          :color="isSubstitutionMode ? 'deep-orange' : 'primary'"
          :icon="isSubstitutionMode ? 'swap_horiz' : 'school'"
          :label="isSubstitutionMode ? 'Torna alle Mie Classi' : 'Supplenza in Altre Classi / Gruppi'"
          unelevated
          class="rounded-lg text-weight-bold"
          @click="toggleSubstitutionMode"
        />
        <q-btn icon="add" label="Nuova Lezione" color="primary" @click="openNewLesson" />
      </div>
    </div>

    <!-- Substitution Info Banner -->
    <q-banner v-if="isSubstitutionMode" class="bg-amber-1 text-amber-10 rounded-xl border border-amber-300 q-mb-md shadow-soft">
      <template v-slot:avatar>
        <q-icon name="swap_horiz" color="amber-9" size="28px" />
      </template>
      <div class="text-weight-bold text-subtitle1">Modalità Supplenza Occasionale nel Registro di Classe — {{ selectedClassLabel }}</div>
      <div class="text-caption">
        Stai registrando la firma di lezione nel Registro di Classe per una classe/gruppo della scuola in qualità di docente supplente.
        Puoi firmare l'ora come <strong>Supplenza</strong>, specificare l'argomento trattato ed annotare eventuali compiti o note per gli alunni.
      </div>
    </q-banner>

    <!-- Filters -->
    <q-card class="q-mb-md shadow-1">
      <q-card-section class="row items-center q-gutter-md q-py-sm">
        <q-select
          v-model="selectedClass"
          :options="availableClassOptions"
          option-value="id"
          option-label="label"
          emit-value map-options
          label="Classe / Gruppo"
          dense outlined
          style="min-width:240px"
        />
        <q-input
          v-model="selectedDate"
          type="date"
          label="Giorno *"
          dense outlined
          style="min-width:170px"
        />
        <q-select
          v-model="selectedSubject"
          :options="availableSubjectOptions"
          option-value="subject_id"
          option-label="subject_name"
          emit-value map-options
          label="Materia"
          dense outlined
          style="min-width:180px"
          clearable
          :disable="isSubstitutionMode"
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

    <!-- PCTO & Orientamento Hour Counters -->
    <div v-if="selectedClass && (activityHours.pcto > 0 || activityHours.orientamento > 0 || true)" class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-6">
        <q-card class="bg-deep-purple-1 shadow-1">
          <q-card-section class="row items-center q-py-sm">
            <q-icon name="work" color="deep-purple" class="q-mr-sm" size="28px" />
            <div>
              <div class="text-caption text-grey-7 text-weight-bold">Ore PCTO Svolte</div>
              <div class="text-h5 text-weight-bold text-deep-purple">
                {{ activityHours.pcto }}
                <span class="text-caption text-grey-7">ore</span>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6">
        <q-card class="bg-teal-1 shadow-1">
          <q-card-section class="row items-center q-py-sm">
            <q-icon name="explore" color="teal" class="q-mr-sm" size="28px" />
            <div>
              <div class="text-caption text-grey-7 text-weight-bold">Ore Orientamento Svolte</div>
              <div class="text-h5 text-weight-bold text-teal">
                {{ activityHours.orientamento }}
                <span class="text-caption text-grey-7">ore</span>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Lessons List -->
    <div v-if="activeTab === 'lessons'">
      <q-card v-if="sortedLessons.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="menu_book" size="80px" class="q-mb-md" />
        <div class="text-h6">Nessuna lezione registrata il {{ formatDate(selectedDate) }}</div>
        <div class="text-caption">Aggiungi la prima lezione per questa giornata per compilare il registro</div>
      </q-card>
      <q-card v-else class="shadow-1">
        <q-list separator>
          <q-item v-for="lesson in sortedLessons" :key="lesson.id" class="q-py-md">
            <q-item-section avatar>
              <q-avatar :color="getActivityTypeColor(lesson.activity_type)" text-color="white" :icon="getActivityTypeIcon(lesson.activity_type)" />
            </q-item-section>
            <q-item-section>
              <div class="row items-center q-gutter-xs q-mb-xs">
                <q-badge color="indigo" class="text-weight-bold text-caption">
                  {{ lesson.hour || 1 }}ª Ora ({{ lesson.duration || 1 }}h)
                </q-badge>
                <q-badge outline :color="getLessonTypeColor(lesson.type)">{{ lesson.type }}</q-badge>
                <q-badge v-if="lesson.activity_type && lesson.activity_type !== 'standard'" :color="getActivityTypeColor(lesson.activity_type)">
                  <q-icon :name="getActivityTypeIcon(lesson.activity_type)" size="14px" class="q-mr-xs" />{{ getActivityTypeLabel(lesson.activity_type) }}
                </q-badge>
                <q-badge v-if="lesson.is_co_teaching" color="deep-purple" outline>
                  <q-icon name="people" size="14px" class="q-mr-xs"/>Compresenza
                </q-badge>
              </div>
              <q-item-label class="text-weight-bold text-subtitle1">{{ lesson.topic }}</q-item-label>
              <q-item-label caption class="text-grey-8">
                <strong>Data:</strong> {{ formatDate(lesson.date) }} &bull;
                <strong v-if="lesson.teacher_name">Docente:</strong> {{ lesson.teacher_name }}
                <span class="q-ml-xs" v-if="getSubjectName(lesson.subject_id)">
                  &bull; <strong>Materia:</strong> {{ getSubjectName(lesson.subject_id) }}
                </span>
              </q-item-label>
              <q-item-label caption class="text-grey-7 q-mt-xs" v-if="lesson.notes">
                <em>Note: {{ lesson.notes }}</em>
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row items-center q-gutter-xs">
                <q-btn flat round dense icon="visibility" color="info" size="sm" @click="viewLessonDetails(lesson)">
                  <q-tooltip>Visualizza Dettagli</q-tooltip>
                </q-btn>
                <template v-if="canManageLesson(lesson)">
                  <q-btn flat round dense icon="edit" color="primary" size="sm" @click="openEditLesson(lesson)">
                    <q-tooltip>Modifica Lezione</q-tooltip>
                  </q-btn>
                  <q-btn flat round dense icon="delete" color="negative" size="sm" @click="confirmDeleteLesson(lesson)">
                    <q-tooltip>Elimina Lezione</q-tooltip>
                  </q-btn>
                </template>
              </div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>
    </div>

    <!-- Homework List -->
    <div v-if="activeTab === 'homework'">
      <q-card v-if="homeworks.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
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
                <span class="q-ml-xs text-grey-8" v-if="getSubjectName(hw.subject_id)">
                  &bull; {{ getSubjectName(hw.subject_id) }}
                </span>
                <span class="q-ml-xs text-grey-8" v-if="hw.teacher_name">
                  &bull; Docente: {{ hw.teacher_name }}
                </span>
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row items-center q-gutter-xs" v-if="canManageHomework(hw)">
                <q-btn flat round dense icon="delete" color="negative" size="sm" @click="confirmDeleteHomework(hw)">
                  <q-tooltip>Elimina Compito</q-tooltip>
                </q-btn>
              </div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>
    </div>

    <!-- Lesson Details Dialog -->
    <q-dialog v-model="detailsDialog">
      <q-card style="min-width: 460px">
        <q-card-section class="row items-center justify-between bg-primary text-white">
          <div class="text-h6 row items-center">
            <q-icon name="menu_book" class="q-mr-sm" />
            Dettagli Lezione
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-gutter-sm" v-if="selectedLesson">
          <div class="text-subtitle1 text-weight-bold text-primary">{{ selectedLesson.topic }}</div>
          <div><strong>Data:</strong> {{ formatDate(selectedLesson.date) }}</div>
          <div><strong>Ora lezione:</strong> {{ selectedLesson.hour || 1 }}ª Ora (Durata: {{ selectedLesson.duration || 1 }}h)</div>
          <div><strong>Docente:</strong> {{ selectedLesson.teacher_name || '-' }}</div>
          <div>
            <strong>Tipologia:</strong>
            <q-badge outline :color="getLessonTypeColor(selectedLesson.type)" class="q-ml-xs">
              {{ selectedLesson.type }}
            </q-badge>
            <q-badge v-if="selectedLesson.is_co_teaching" color="deep-purple" outline class="q-ml-xs">
              <q-icon name="people" size="14px" class="q-mr-xs"/>Compresenza
            </q-badge>
          </div>
          <div><strong>Materia:</strong> {{ getSubjectName(selectedLesson.subject_id) || '-' }}</div>
          <div v-if="selectedLesson.notes" class="q-mt-sm">
            <strong>Note / Osservazioni:</strong>
            <div class="text-grey-8 bg-grey-2 q-pa-sm rounded-borders q-mt-xs">{{ selectedLesson.notes }}</div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Chiudi" color="grey-8" v-close-popup />
          <q-btn
            v-if="selectedLesson && canManageLesson(selectedLesson)"
            outline icon="edit" label="Modifica" color="primary"
            @click="detailsDialog = false; openEditLesson(selectedLesson)"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- New / Edit Lesson Dialog -->
    <q-dialog v-model="lessonDialog" persistent>
      <q-card style="min-width: 480px">
        <q-card-section>
          <div class="text-h6">
            <q-icon name="menu_book" class="q-mr-xs" color="primary" />
            {{ isEditingLesson ? 'Modifica Lezione' : 'Nuova Lezione' }}
          </div>
        </q-card-section>
        <q-card-section class="q-gutter-md">
          <q-select
            v-model="newLesson.subject_id"
            :options="availableSubjectOptions"
            option-value="subject_id"
            option-label="subject_name"
            emit-value map-options
            label="Materia *"
            outlined dense
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
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
            placeholder="Es: Le equazioni di secondo grado o Supplenza docente assente"
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-md-6">
              <q-input
                v-model.number="newLesson.hour"
                type="number"
                label="Ora Lezione *"
                outlined dense
                min="1" max="10"
                :rules="[v => !!v || 'Campo obbligatorio']"
              />
            </div>
            <div class="col-12 col-md-6">
              <q-input
                v-model.number="newLesson.duration"
                type="number"
                label="Durata (ore) *"
                outlined dense
                min="1" max="5"
                :rules="[v => !!v || 'Campo obbligatorio']"
              />
            </div>
          </div>
          <q-select
            v-model="newLesson.type"
            :options="['Frontale', 'Supplenza', 'Laboratorio', 'Verifica', 'Discussione', 'Lavoro di gruppo', 'Altro']"
            label="Tipo di Lezione *"
            outlined dense
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
          <q-select
            v-model="newLesson.activity_type"
            :options="activityTypeOptions"
            option-value="value"
            option-label="label"
            emit-value map-options
            label="Tipologia Attività"
            outlined dense
          >
            <template v-slot:option="scope">
              <q-item v-bind="scope.itemProps">
                <q-item-section avatar>
                  <q-icon :name="scope.opt.icon" :color="scope.opt.color" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ scope.opt.label }}</q-item-label>
                  <q-item-label caption>{{ scope.opt.caption }}</q-item-label>
                </q-item-section>
              </q-item>
            </template>
          </q-select>
          <q-banner
            v-if="newLesson.activity_type === 'pcto' || newLesson.activity_type === 'orientamento'"
            class="bg-blue-1 text-blue-9 rounded-borders q-mt-xs"
            dense
          >
            <template v-slot:avatar><q-icon name="info" color="blue-7" /></template>
            Per le attività {{ getActivityTypeLabel(newLesson.activity_type) }} non è possibile inserire valutazioni agli studenti.
          </q-banner>
          <q-toggle
            v-model="newLesson.is_co_teaching"
            label="Compresenza (docente co-presente in classe)"
            color="deep-purple"
            icon="people"
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
            v-if="!isEditingLesson && !isPctoOrOrientamento"
            v-model="assignHomeworkToo"
            label="Assegna anche un compito per questa lezione"
            color="orange"
          />
          <template v-if="assignHomeworkToo && !isEditingLesson">
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
          <q-btn color="primary" :label="isEditingLesson ? 'Aggiorna' : 'Salva'" :loading="saving" @click="saveLesson" />
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
          <q-select
            v-model="newHomework.subject_id"
            :options="availableSubjectOptions"
            option-value="subject_id"
            option-label="subject_name"
            emit-value map-options
            label="Materia *"
            outlined dense
            :rules="[v => !!v || 'Campo obbligatorio']"
          />
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
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar, date } from 'quasar'
import { useClassesStore } from 'src/stores/classes'
import { useGradesStore } from 'src/stores/grades'
import { useAuthStore } from 'src/stores/auth'
import { lessonService } from 'src/services/lessonService'

const $q = useQuasar()
const classesStore = useClassesStore()
const gradesStore = useGradesStore()
const authStore = useAuthStore()

const selectedClass = ref(null)
const selectedSubject = ref(null)
const selectedDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'))
const activeTab = ref('lessons')
const isSubstitutionMode = ref(false)
const allSchoolClasses = ref([])

const lessonDialog = ref(false)
const isEditingLesson = ref(false)
const editingLessonId = ref(null)

const detailsDialog = ref(false)
const selectedLesson = ref(null)

const homeworkDialog = ref(false)
const saving = ref(false)
const assignHomeworkToo = ref(false)

const lessons = ref([])
const homeworks = ref([])

const availableClassOptions = computed(() => {
  const source = isSubstitutionMode.value ? allSchoolClasses.value : classesStore.classes
  return source.map(c => {
    let nameText = c.name || `Classe ${c.id}`
    if (c.section && !nameText.endsWith(c.section)) {
      nameText += c.section
    }
    if (c.articolazione) {
      nameText += ` - ${c.articolazione}`
    }
    return {
      ...c,
      label: c.label || nameText
    }
  })
})

const selectedClassLabel = computed(() => {
  const found = availableClassOptions.value.find(c => String(c.id) === String(selectedClass.value))
  return found ? found.label : 'Classe'
})

const availableSubjectOptions = computed(() => {
  if (isSubstitutionMode.value) {
    return [{ subject_name: 'Supplenza / Compresenza', subject_id: 'supplenza' }]
  }
  return gradesStore.subjects || []
})

const toggleSubstitutionMode = async () => {
  isSubstitutionMode.value = !isSubstitutionMode.value
  if (isSubstitutionMode.value) {
    allSchoolClasses.value = await classesStore.fetchAllSchoolClassesAndGroups()
    if (allSchoolClasses.value.length > 0) {
      selectedClass.value = allSchoolClasses.value[0].id
    }
    selectedSubject.value = 'supplenza'
    $q.notify({ type: 'info', message: 'Modalità Supplenza attivata nel Registro di Classe', timeout: 3000 })
  } else {
    await classesStore.fetchAssignedClasses()
    if (classesStore.classes.length > 0) {
      selectedClass.value = classesStore.classes[0].id
    }
  }
  fetchData()
}

const sortedLessons = computed(() => {
  return [...lessons.value].sort((a, b) => (a.hour || 1) - (b.hour || 1))
})

const canManageLesson = (lesson) => {
  if (!lesson) return false
  const user = authStore.user
  if (!user) return true
  if (['admin', 'superadmin', 'secretary'].includes(user.role)) return true
  return lesson.teacher_id === user.id
}

const canManageHomework = (hw) => {
  if (!hw) return false
  const user = authStore.user
  if (!user) return true
  if (['admin', 'superadmin', 'secretary'].includes(user.role)) return true
  return hw.teacher_id === user.id
}

const newLesson = ref({
  date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
  hour: 1,
  duration: 1,
  topic: '',
  type: 'Frontale',
  is_co_teaching: false,
  notes: '',
  subject_id: null,
  homeworkDesc: '',
  homeworkDue: ''
})

const newHomework = ref({
  description: '',
  subject_id: null,
  dueDate: date.formatDate(Date.now(), 'YYYY-MM-DD')
})

onMounted(async () => {
  await classesStore.fetchAssignedClasses()
  if (classesStore.classes.length > 0) {
    selectedClass.value = classesStore.classes[0].id
  }
})

watch([selectedClass, selectedDate], async () => {
  if (selectedClass.value && !isSubstitutionMode.value) {
    await gradesStore.fetchClassSubjects(selectedClass.value)
    if (gradesStore.subjects && gradesStore.subjects.length > 0 && !selectedSubject.value) {
      selectedSubject.value = gradesStore.subjects[0].subject_id
    }
  } else if (isSubstitutionMode.value) {
    selectedSubject.value = 'supplenza'
  }
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
    const subjParam = selectedSubject.value === 'supplenza' ? null : selectedSubject.value
    const res = await lessonService.getLessons(selectedClass.value, subjParam, selectedDate.value)
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

const getSubjectName = (subjectId) => {
  if (!subjectId || subjectId === 'supplenza') return 'Supplenza / Compresenza'
  const sub = (gradesStore.subjects || []).find(s => s.subject_id === subjectId || s.id === subjectId)
  return sub ? sub.subject_name || sub.name : ''
}

const viewLessonDetails = (lesson) => {
  selectedLesson.value = lesson
  detailsDialog.value = true
}

const openNewLesson = () => {
  isEditingLesson.value = false
  editingLessonId.value = null
  newLesson.value = {
    date: selectedDate.value || date.formatDate(Date.now(), 'YYYY-MM-DD'),
    hour: (sortedLessons.value.length > 0 ? (sortedLessons.value[sortedLessons.value.length - 1].hour || 1) + 1 : 1),
    duration: 1,
    topic: isSubstitutionMode.value ? 'Supplenza docente assente' : '',
    type: isSubstitutionMode.value ? 'Supplenza' : 'Frontale',
    is_co_teaching: false,
    notes: '',
    subject_id: isSubstitutionMode.value ? 'supplenza' : selectedSubject.value,
    homeworkDesc: '',
    homeworkDue: ''
  }
  assignHomeworkToo.value = false
  lessonDialog.value = true
}

const openEditLesson = (lesson) => {
  isEditingLesson.value = true
  editingLessonId.value = lesson.id
  let formattedDate = lesson.date
  if (lesson.date) {
    formattedDate = date.formatDate(new Date(lesson.date), 'YYYY-MM-DD')
  }
  newLesson.value = {
    date: formattedDate,
    hour: lesson.hour || 1,
    duration: lesson.duration || 1,
    topic: lesson.topic || '',
    type: lesson.type || 'Frontale',
    is_co_teaching: !!lesson.is_co_teaching,
    notes: lesson.notes || '',
    subject_id: lesson.subject_id || selectedSubject.value,
    homeworkDesc: '',
    homeworkDue: ''
  }
  assignHomeworkToo.value = false
  lessonDialog.value = true
}

const confirmDeleteLesson = (lesson) => {
  $q.dialog({
    title: 'Elimina Lezione',
    message: `Sei sicuro di voler eliminare la lezione "${lesson.topic}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await lessonService.deleteLesson(lesson.id)
      $q.notify({ type: 'positive', message: 'Lezione eliminata con successo' })
      await fetchLessons()
    } catch (e) {
      $q.notify({ type: 'negative', message: `Errore durante l'eliminazione: ${e.response?.data?.error || e.message}` })
    }
  })
}

const confirmDeleteHomework = (hw) => {
  $q.dialog({
    title: 'Elimina Compito',
    message: `Sei sicuro di voler eliminare questo compito?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await lessonService.deleteHomework(hw.id)
      $q.notify({ type: 'positive', message: 'Compito eliminato con successo' })
      await fetchHomeworks()
    } catch (e) {
      $q.notify({ type: 'negative', message: `Errore durante l'eliminazione: ${e.response?.data?.error || e.message}` })
    }
  })
}

const openNewHomework = () => {
  newHomework.value = {
    description: '',
    subject_id: isSubstitutionMode.value ? 'supplenza' : selectedSubject.value,
    dueDate: date.formatDate(Date.now(), 'YYYY-MM-DD')
  }
  homeworkDialog.value = true
}

const saveLesson = async () => {
  if (!newLesson.value.topic || !newLesson.value.date || !newLesson.value.type || !newLesson.value.subject_id) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }
  saving.value = true
  try {
    const lessonPayload = {
      class_id: selectedClass.value,
      subject_id: newLesson.value.subject_id === 'supplenza' ? null : newLesson.value.subject_id,
      date: newLesson.value.date,
      hour: newLesson.value.hour,
      duration: newLesson.value.duration,
      topic: newLesson.value.topic,
      type: newLesson.value.type,
      is_co_teaching: newLesson.value.is_co_teaching,
      notes: newLesson.value.notes
    }

    if (isEditingLesson.value && editingLessonId.value) {
      await lessonService.updateLesson(editingLessonId.value, lessonPayload)
      $q.notify({ type: 'positive', message: 'Lezione aggiornata con successo' })
    } else {
      const lessonRes = await lessonService.createLesson(lessonPayload)

      if (assignHomeworkToo.value && newLesson.value.homeworkDesc) {
        await lessonService.createHomework({
          class_id: selectedClass.value,
          subject_id: newLesson.value.subject_id === 'supplenza' ? null : newLesson.value.subject_id,
          lesson_id: lessonRes.data?.id,
          due_date: newLesson.value.homeworkDue,
          description: newLesson.value.homeworkDesc
        })
      }
      $q.notify({ type: 'positive', message: 'Lezione registrata con successo' })
    }

    lessonDialog.value = false
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: `Errore: ${e.response?.data?.error || e.message}` })
  } finally {
    saving.value = false
  }
}

const saveHomework = async () => {
  if (!newHomework.value.description || !newHomework.value.dueDate || !newHomework.value.subject_id) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }
  saving.value = true
  try {
    await lessonService.createHomework({
      class_id: selectedClass.value,
      subject_id: newHomework.value.subject_id === 'supplenza' ? null : newHomework.value.subject_id,
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
    'Supplenza': 'deep-orange',
    'Laboratorio': 'teal',
    'Verifica': 'deep-orange',
    'Discussione': 'purple',
    'Lavoro di gruppo': 'cyan',
  }
  return map[type] || 'grey'
}
</script>

<style scoped>
.rounded-borders {
  border-radius: 8px;
}
</style>
