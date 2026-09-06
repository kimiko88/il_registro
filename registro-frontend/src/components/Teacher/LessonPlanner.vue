 <template>
  <q-page class="q-pa-md" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-grey-1 text-dark'">

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
        <q-btn
          v-if="activeTab === 'lessons'"
          icon="add" label="Nuova Lezione" color="primary"
          @click="openNewLesson"
        />
        <q-btn
          v-if="activeTab === 'free'"
          icon="add" label="Nuova Attività" color="teal"
          @click="openNewFreeActivity"
        />
      </div>
    </div>

    <!-- Substitution Info Banner -->
    <q-banner v-if="isSubstitutionMode" class="rounded-xl border q-mb-md shadow-soft" :class="$q.dark.isActive ? 'bg-amber-10 text-amber-1 border-amber-8' : 'bg-amber-1 text-amber-10 border-amber-300'">
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
          v-if="activeTab !== 'free'"
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
          v-if="activeTab === 'lessons'"
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
          :options="[
            {label:'Lezioni', value:'lessons', icon:'menu_book'},
            {label:'Compiti', value:'homework', icon:'assignment'},
            {label:'Ore Libere', value:'free', icon:'event_busy'}
          ]"
        />
      </q-card-section>
    </q-card>

    <!-- PCTO & Orientamento Hour Counters (solo tab lezioni) -->
    <div v-if="activeTab === 'lessons' && selectedClass" class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-6">
        <q-card class="shadow-1" :class="$q.dark.isActive ? 'bg-deep-purple-10 text-white' : 'bg-deep-purple-1'">
          <q-card-section class="row items-center q-py-sm">
            <q-icon name="work" color="deep-purple-3" class="q-mr-sm" size="28px" />
            <div>
              <div class="text-caption text-weight-bold" :class="$q.dark.isActive ? 'text-grey-3' : 'text-grey-7'">Ore PCTO Svolte</div>
              <div class="text-h5 text-weight-bold" :class="$q.dark.isActive ? 'text-deep-purple-2' : 'text-deep-purple'">
                {{ activityHours.pcto }}
                <span class="text-caption" :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-7'">ore</span>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6">
        <q-card class="shadow-1" :class="$q.dark.isActive ? 'bg-teal-10 text-white' : 'bg-teal-1'">
          <q-card-section class="row items-center q-py-sm">
            <q-icon name="explore" color="teal-3" class="q-mr-sm" size="28px" />
            <div>
              <div class="text-caption text-weight-bold" :class="$q.dark.isActive ? 'text-grey-3' : 'text-grey-7'">Ore Orientamento Svolte</div>
              <div class="text-h5 text-weight-bold" :class="$q.dark.isActive ? 'text-teal-2' : 'text-teal'">
                {{ activityHours.orientamento }}
                <span class="text-caption" :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-7'">ore</span>
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

    <!-- ─── Attività Libere Tab ─────────────────────────────────── -->
    <div v-if="activeTab === 'free'">
      <q-banner class="rounded-xl border q-mb-md" :class="$q.dark.isActive ? 'bg-teal-10 text-teal-1 border-teal-8' : 'bg-teal-1 text-teal-9 border-teal-3'" dense>
        <template v-slot:avatar><q-icon name="event_busy" color="teal-8" size="22px" /></template>
        <div class="text-weight-bold">Ore a Disposizione / Attività Non in Classe</div>
        <div class="text-caption">
          Registra qui le tue ore quando la classe è in gita, sei in riunione, formazione, disponibilità, etc.
          Queste attività non sono collegate ad una classe specifica.
        </div>
      </q-banner>

      <q-card v-if="freeActivities.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="event_busy" size="80px" class="q-mb-md" />
        <div class="text-h6">Nessuna attività libera registrata</div>
        <div class="text-caption">Aggiungi un'attività per il giorno {{ formatDate(selectedDate) }}</div>
        <q-btn label="Aggiungi Attività" color="teal" @click="openNewFreeActivity" class="q-mt-md" />
      </q-card>
      <q-card v-else class="shadow-1">
        <q-card-section class="row justify-end q-pb-none">
          <q-btn flat icon="add" label="Aggiungi Attività" color="teal" @click="openNewFreeActivity" />
        </q-card-section>
        <q-list separator>
          <q-item v-for="act in freeActivities" :key="act.id" class="q-py-md">
            <q-item-section avatar>
              <q-avatar :color="getFreeActivityColor(act.activity_type)" text-color="white" :icon="getFreeActivityIcon(act.activity_type)" />
            </q-item-section>
            <q-item-section>
              <div class="row items-center q-gutter-xs q-mb-xs">
                <q-badge color="teal-8" class="text-weight-bold text-caption">
                  {{ act.start_hour }}ª Ora ({{ act.duration }}h)
                </q-badge>
                <q-badge :color="getFreeActivityColor(act.activity_type)" outline>
                  {{ getFreeActivityLabel(act.activity_type) }}
                </q-badge>
              </div>
              <q-item-label class="text-weight-bold text-subtitle1">{{ act.description }}</q-item-label>
              <q-item-label caption class="text-grey-8">
                <strong>Data:</strong> {{ formatDate(act.date) }}
              </q-item-label>
              <q-item-label caption class="text-grey-7 q-mt-xs" v-if="act.notes">
                <em>Note: {{ act.notes }}</em>
              </q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="row items-center q-gutter-xs">
                <q-btn flat round dense icon="edit" color="teal" size="sm" @click="openEditFreeActivity(act)">
                  <q-tooltip>Modifica</q-tooltip>
                </q-btn>
                <q-btn flat round dense icon="delete" color="negative" size="sm" @click="confirmDeleteFreeActivity(act)">
                  <q-tooltip>Elimina</q-tooltip>
                </q-btn>
              </div>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card>
    </div>

    <!-- Lesson Details Dialog -->
    <q-dialog v-model="detailsDialog">
      <q-card style="width: min(500px, 95vw); max-width: 95vw;">
        <q-card-section class="row items-center justify-between bg-primary text-white">
          <div class="text-h6 row items-center">
            <q-icon name="menu_book" class="q-mr-sm" />
            Dettagli Lezione
          </div>
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
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
            <div class="q-pa-sm rounded-borders q-mt-xs" :class="$q.dark.isActive ? 'bg-grey-9 text-grey-3' : 'bg-grey-2 text-grey-8'">{{ selectedLesson.notes }}</div>
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

    <!-- Modals for Lessons, Homework and Free Activities -->
    <LessonFormDialog
      v-model="lessonDialog"
      :is-editing="isEditingLesson"
      :initial-data="newLesson"
      :available-subject-options="availableSubjectOptions"
      :activity-type-options="activityTypeOptions"
      :saving="saving"
      @save="saveLesson"
    />

    <HomeworkFormDialog
      v-model="homeworkDialog"
      :initial-data="newHomework"
      :available-subject-options="availableSubjectOptions"
      :saving="saving"
      @save="saveHomework"
    />

    <FreeActivityDialog
      v-model="freeActivityDialog"
      :is-editing="isEditingFreeActivity"
      :initial-data="newFreeActivity"
      :free-activity-type-options="freeActivityTypeOptions"
      :saving="saving"
      @save="saveFreeActivity"
    />
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar, date } from 'quasar'
import { useClassesStore } from '@/stores/classes'
import { useGradesStore } from '@/stores/grades'
import { lessonService } from '@/services/lessonService'
import { teacherActivityService } from '@/services/teacherActivityService'
import { useOfflineSync } from '@/composables/useOfflineSync'
import LessonFormDialog from './LessonFormDialog.vue'
import HomeworkFormDialog from './HomeworkFormDialog.vue'
import FreeActivityDialog from './FreeActivityDialog.vue'

const $q = useQuasar()
const { t } = useI18n()
const { executeWithOfflineQueue } = useOfflineSync()
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

// ─── Attività Libere State ───────────────────────────────────────────
const freeActivities = ref([])
const freeActivityDialog = ref(false)
const isEditingFreeActivity = ref(false)
const editingFreeActivityId = ref(null)

const newFreeActivity = ref({
  date: date.formatDate(Date.now(), 'YYYY-MM-DD'),
  start_hour: 1,
  duration: 1,
  activity_type: 'disponibilita',
  description: '',
  notes: ''
})

/** Opzioni per il tipo di attività libera docente */
const freeActivityTypeOptions = [
  { value: 'disponibilita', label: 'Disponibilità / Classe in gita', icon: 'event_busy', color: 'teal', caption: 'Classe in gita o ore a disposizione' },
  { value: 'riunione',      label: 'Riunione / Consiglio di Classe', icon: 'groups',      color: 'indigo', caption: 'Riunione di dipartimento, consigli, etc.' },
  { value: 'formazione',    label: 'Formazione / Aggiornamento',     icon: 'school',      color: 'blue', caption: 'Corsi di aggiornamento professionale' },
  { value: 'ptof',          label: 'Attività PTOF',                  icon: 'auto_stories', color: 'purple', caption: 'Attività rientranti nel PTOF scolastico' },
  { value: 'gita',          label: 'Gita / Uscita didattica',        icon: 'luggage',     color: 'orange', caption: 'Accompagnamento in gita o uscita didattica' },
  { value: 'altro',         label: 'Altro',                          icon: 'more_horiz',  color: 'grey', caption: 'Altra attività non categorizzata' }
]

const getFreeActivityLabel = (type) => {
  const opt = freeActivityTypeOptions.find(o => o.value === type)
  return opt ? opt.label : type
}
const getFreeActivityIcon = (type) => {
  const opt = freeActivityTypeOptions.find(o => o.value === type)
  return opt ? opt.icon : 'event_busy'
}
const getFreeActivityColor = (type) => {
  const opt = freeActivityTypeOptions.find(o => o.value === type)
  return opt ? opt.color : 'grey'
}

const openNewFreeActivity = () => {
  isEditingFreeActivity.value = false
  editingFreeActivityId.value = null
  newFreeActivity.value = {
    date: selectedDate.value || date.formatDate(Date.now(), 'YYYY-MM-DD'),
    start_hour: 1,
    duration: 1,
    activity_type: 'disponibilita',
    description: '',
    notes: ''
  }
  freeActivityDialog.value = true
}

const openEditFreeActivity = (act) => {
  isEditingFreeActivity.value = true
  editingFreeActivityId.value = act.id
  newFreeActivity.value = {
    date: date.formatDate(new Date(act.date), 'YYYY-MM-DD'),
    start_hour: act.start_hour,
    duration: act.duration,
    activity_type: act.activity_type,
    description: act.description,
    notes: act.notes || ''
  }
  freeActivityDialog.value = true
}

const saveFreeActivity = async (formData) => {
  saving.value = true
  try {
    const payload = {
      date: formData.date,
      start_hour: formData.start_hour,
      duration: formData.duration,
      activity_type: formData.activity_type || 'disponibilita',
      description: formData.description,
      notes: formData.notes || ''
    }
    if (isEditingFreeActivity.value && editingFreeActivityId.value) {
      await executeWithOfflineQueue(
        { url: `/teacher-activities/${editingFreeActivityId.value}`, method: 'put', data: payload },
        { title: `Modifica Attività Libera - ${payload.description}` }
      )
      $q.notify({ type: 'positive', message: 'Attività aggiornata con successo' })
    } else {
      await executeWithOfflineQueue(
        { url: '/teacher-activities', method: 'post', data: payload },
        { title: `Registra Attività Libera - ${payload.description}` }
      )
      $q.notify({ type: 'positive', message: 'Attività registrata con successo' })
    }
    freeActivityDialog.value = false
    await fetchFreeActivities()
  } catch (e) {
    $q.notify({ type: 'negative', message: `Errore: ${e.response?.data?.error || e.message}` })
  } finally {
    saving.value = false
  }
}

const confirmDeleteFreeActivity = (act) => {
  $q.dialog({
    title: 'Elimina Attività',
    message: `Sei sicuro di voler eliminare "${act.description}"?`,
    cancel: true,
    persistent: true,
    ok: { label: 'Elimina', color: 'negative' }
  }).onOk(async () => {
    try {
      await teacherActivityService.delete(act.id)
      $q.notify({ type: 'positive', message: 'Attività eliminata' })
      await fetchFreeActivities()
    } catch (e) {
      $q.notify({ type: 'negative', message: `Errore: ${e.response?.data?.error || e.message}` })
    }
  })
}

const fetchFreeActivities = async () => {
  try {
    const from = selectedDate.value
    const to = selectedDate.value
    const res = await teacherActivityService.list(from, to)
    freeActivities.value = res.data || []
  } catch (e) {
    freeActivities.value = []
    console.error('fetchFreeActivities error:', e)
  }
}
// ─────────────────────────────────────────────────────────────────────────────

// ── Activity Hours (PCTO / Orientamento counters) ─────────────
const activityHours = ref({ pcto: 0, orientamento: 0 })

const fetchActivityHours = async () => {
  if (!selectedClass.value) return
  try {
    const res = await lessonService.getActivityHours(selectedClass.value)
    const data = res.data || {}
    activityHours.value = {
      pcto: data.pcto || 0,
      orientamento: data.orientamento || 0
    }
  } catch (e) {
    activityHours.value = { pcto: 0, orientamento: 0 }
  }
}

/** Opzioni per il tipo di attività *lezione* (collegata a classe) */
const activityTypeOptions = [
  { value: 'standard',          label: 'Standard',           icon: 'menu_book',    color: 'primary',     caption: 'Lezione curricolare ordinaria' },
  { value: 'substitution',      label: 'Supplenza',          icon: 'swap_horiz',   color: 'deep-orange', caption: 'Supplenza di un collega assente' },
  { value: 'pcto',              label: 'PCTO',               icon: 'work',         color: 'deep-purple', caption: 'Ore di Alternanza Scuola-Lavoro (PCTO)' },
  { value: 'orientamento',      label: 'Orientamento',       icon: 'explore',      color: 'teal',        caption: 'Attività di orientamento formativo' },
  { value: 'pcto_orientamento', label: 'PCTO - Orientamento', icon: 'hub',          color: 'indigo-8',    caption: 'Attività congiunta PCTO e Orientamento (max 15h ciascuno)' },
  { value: 'ptof',              label: 'PTOF',               icon: 'auto_stories', color: 'purple',      caption: 'Attività rientranti nel PTOF' },
  { value: 'project',           label: 'Progetto',           icon: 'science',      color: 'indigo',      caption: 'Progetto didattico specifico' },
  { value: 'assembly',          label: 'Assemblea',          icon: 'groups',       color: 'blue',        caption: 'Assemblea di istituto o di classe' },
  { value: 'trip',              label: 'Gita',               icon: 'luggage',      color: 'orange',      caption: 'Uscita didattica o gita scolastica' },
  { value: 'lab',               label: 'Laboratorio',        icon: 'biotech',      color: 'green',       caption: 'Attività di laboratorio' },
  { value: 'other',             label: 'Altro',              icon: 'more_horiz',   color: 'grey',        caption: 'Altra tipologia non categorizzata' }
]

const getActivityTypeColor = (type) => {
  const opt = activityTypeOptions.find(o => o.value === type)
  return opt ? opt.color : 'grey'
}
const getActivityTypeIcon = (type) => {
  const opt = activityTypeOptions.find(o => o.value === type)
  return opt ? opt.icon : 'menu_book'
}
const getActivityTypeLabel = (type) => {
  const opt = activityTypeOptions.find(o => o.value === type)
  return opt ? opt.label : type
}
// ─────────────────────────────────────────────────────────────────────────────

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
  activity_type: 'standard',
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
    if (availableSubjectOptions.value.length > 0) {
      const isSelectedAvailable = availableSubjectOptions.value.some(s => s.subject_id === selectedSubject.value)
      if (!selectedSubject.value || !isSelectedAvailable) {
        selectedSubject.value = availableSubjectOptions.value[0].subject_id
      }
    }
  } else if (isSubstitutionMode.value) {
    selectedSubject.value = 'supplenza'
  }
  fetchData()
})

watch(selectedSubject, () => {
  if (activeTab.value === 'lessons') fetchLessons()
})

// Quando si passa alla tab "free", aggiorna le attività libere
watch(activeTab, (newTab) => {
  if (newTab === 'free') fetchFreeActivities()
})

const fetchData = async () => {
  await Promise.all([fetchLessons(), fetchHomeworks(), fetchActivityHours()])
  if (activeTab.value === 'free') fetchFreeActivities()
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
    activity_type: isSubstitutionMode.value ? 'substitution' : 'standard',
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
    activity_type: lesson.activity_type || 'standard',
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

const saveLesson = async (formData) => {
  saving.value = true
  try {
    const lessonPayload = {
      class_id: selectedClass.value,
      subject_id: formData.subject_id === 'supplenza' ? null : formData.subject_id,
      date: formData.date,
      hour: formData.hour,
      duration: formData.duration,
      topic: formData.topic,
      type: formData.type,
      activity_type: formData.activity_type || 'standard',
      is_co_teaching: formData.is_co_teaching,
      notes: formData.notes
    }

    if (isEditingLesson.value && editingLessonId.value) {
      await executeWithOfflineQueue(
        { url: `/lessons/${editingLessonId.value}`, method: 'put', data: lessonPayload },
        { title: `Modifica Lezione - ${lessonPayload.topic}` }
      )
      $q.notify({ type: 'positive', message: 'Lezione aggiornata con successo' })
    } else {
      const lessonRes = await executeWithOfflineQueue(
        { url: '/lessons', method: 'post', data: lessonPayload },
        { title: `Firma Lezione - ${lessonPayload.topic}` }
      )

      if (formData.assignHomework && formData.homeworkDesc) {
        await executeWithOfflineQueue(
          {
            url: '/homeworks',
            method: 'post',
            data: {
              class_id: selectedClass.value,
              subject_id: formData.subject_id === 'supplenza' ? null : formData.subject_id,
              lesson_id: lessonRes?.data?.id,
              due_date: formData.homeworkDue,
              description: formData.homeworkDesc
            }
          },
          { title: `Assegna Compito - ${formData.homeworkDesc}` }
        )
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

const saveHomework = async (formData) => {
  saving.value = true
  try {
    await executeWithOfflineQueue(
      {
        url: '/homeworks',
        method: 'post',
        data: {
          class_id: selectedClass.value,
          subject_id: formData.subject_id === 'supplenza' ? null : formData.subject_id,
          due_date: formData.dueDate,
          description: formData.description
        }
      },
      { title: `Assegna Compito - ${formData.description}` }
    )
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
