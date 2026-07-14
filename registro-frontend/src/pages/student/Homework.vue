<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5 text-weight-bold">
        <q-icon name="assignment" color="orange" class="q-mr-sm" />
        Agenda e Compiti
      </div>
      <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify" narrow-indicator>
        <q-tab name="compiti" label="Lista Compiti" icon="list" />
        <q-tab name="agenda" label="Agenda & Lezioni" icon="calendar_month" />
      </q-tabs>
    </div>

    <q-card v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="60px" />
    </q-card>

    <q-tab-panels v-else v-model="tab" animated class="bg-transparent">
      <!-- Homeworks Tab -->
      <q-tab-panel name="compiti" class="q-pa-none">
        <q-card v-if="homeworks.length === 0" class="text-center q-pa-xl text-grey-6">
          <q-icon name="check_circle" size="80px" class="q-mb-md" color="positive" />
          <div class="text-h6">Nessun compito in sospeso</div>
          <div class="text-caption">Sei in pari con i tuoi compiti!</div>
        </q-card>

        <div v-else>
          <q-card v-for="hw in sortedHomeworks" :key="hw.id" class="q-mb-sm shadow-1">
            <q-item>
              <q-item-section avatar>
                <q-icon
                  :name="isPast(hw.due_date) ? 'warning' : 'assignment'"
                  :color="isPast(hw.due_date) ? 'negative' : 'orange'"
                  size="md"
                />
              </q-item-section>
              <q-item-section>
                <div class="row items-center justify-between">
                  <div class="text-weight-bold text-subtitle1">{{ hw.description }}</div>
                  <q-chip size="sm" color="orange-1" text-color="orange-8" class="text-weight-medium">
                    {{ subjectsMap[hw.subject_id] || hw.subject_id }}
                  </q-chip>
                </div>
                <q-item-label caption class="q-mt-xs">
                  Consegna:
                  <strong :class="isPast(hw.due_date) ? 'text-negative' : 'text-positive'">
                    {{ formatDate(hw.due_date) }}
                  </strong>
                  <q-badge v-if="isPast(hw.due_date)" color="negative" class="q-ml-sm">Scaduto</q-badge>
                  <q-badge v-else-if="isDueSoon(hw.due_date)" color="warning" class="q-ml-sm">Domani</q-badge>
                </q-item-label>
                <q-item-label caption v-if="hw.teacher_name" class="text-grey-6 q-mt-xs">
                  Assegnato da: {{ hw.teacher_name }}
                </q-item-label>
              </q-item-section>
            </q-item>
          </q-card>
        </div>
      </q-tab-panel>

      <!-- Agenda Tab -->
      <q-tab-panel name="agenda" class="q-pa-none">
        <div class="row q-col-gutter-md">
          <!-- Calendar -->
          <div class="col-12 col-md-5">
            <q-card class="shadow-1">
              <q-date
                v-model="selectedDate"
                :events="calendarEvents"
                event-color="primary"
                flat
                class="full-width"
                today-btn
                mask="YYYY-MM-DD"
              />
            </q-card>
          </div>

          <!-- Day Details -->
          <div class="col-12 col-md-7">
            <q-card class="shadow-1 full-height">
              <q-card-section class="bg-primary text-white q-py-sm">
                <div class="text-subtitle1 text-weight-medium">
                  Dettagli del {{ formatDate(selectedDate) }}
                </div>
              </q-card-section>

              <!-- Lessons Section -->
              <q-card-section>
                <div class="text-subtitle2 text-weight-bold text-primary q-mb-sm">
                  <q-icon name="menu_book" class="q-mr-xs" />
                  Lezioni Svolte
                </div>
                <q-list v-if="lessonsOnSelectedDate.length > 0" separator dense>
                  <q-item v-for="lesson in lessonsOnSelectedDate" :key="lesson.id" class="q-py-sm">
                    <q-item-section>
                      <div class="row items-center justify-between">
                        <div class="text-weight-bold text-grey-9">
                          {{ lesson.hour }}° Ora - {{ subjectsMap[lesson.subject_id] || lesson.subject_id }}
                          <span class="text-caption text-grey-6 text-weight-regular q-ml-sm">
                            ({{ lesson.duration }} ore, {{ lesson.type }})
                          </span>
                        </div>
                      </div>
                      <div class="text-body2 text-grey-8 q-mt-xs">{{ lesson.topic }}</div>
                      <div class="text-caption text-grey-6 q-mt-xs" v-if="lesson.notes">
                        Note: {{ lesson.notes }}
                      </div>
                      <div class="text-caption text-grey-5 q-mt-xs" v-if="lesson.teacher_name">
                        Docente: {{ lesson.teacher_name }}
                      </div>
                    </q-item-section>
                  </q-item>
                </q-list>
                <div v-else class="text-caption text-grey q-py-sm">Nessuna lezione registrata in questa data.</div>
              </q-card-section>

              <q-separator inset />

              <!-- Homeworks Section -->
              <q-card-section>
                <div class="text-subtitle2 text-weight-bold text-orange-8 q-mb-sm">
                  <q-icon name="assignment" class="q-mr-xs" />
                  Compiti in Scadenza
                </div>
                <q-list v-if="homeworksOnSelectedDate.length > 0" separator dense>
                  <q-item v-for="hw in homeworksOnSelectedDate" :key="hw.id" class="q-py-sm">
                    <q-item-section>
                      <div class="row items-center justify-between">
                        <div class="text-weight-bold text-grey-9">
                          {{ subjectsMap[hw.subject_id] || hw.subject_id }}
                        </div>
                      </div>
                      <div class="text-body2 text-grey-8 q-mt-xs">{{ hw.description }}</div>
                      <div class="text-caption text-grey-5 q-mt-xs" v-if="hw.teacher_name">
                        Docente: {{ hw.teacher_name }}
                      </div>
                    </q-item-section>
                  </q-item>
                </q-list>
                <div v-else class="text-caption text-grey q-py-sm">Nessun compito con scadenza in questa data.</div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, date } from 'quasar'
import { useStudentStore } from 'src/stores/student'
import { lessonService } from 'src/services/lessonService'
import adminService from 'src/services/adminService'

const $q = useQuasar()
const studentStore = useStudentStore()

const tab = ref('compiti')
const homeworks = ref([])
const lessons = ref([])
const subjectsMap = ref({})
const loading = ref(true)
const selectedDate = ref(date.formatDate(Date.now(), 'YYYY-MM-DD'))

onMounted(async () => {
  await studentStore.fetchProfile()
  await Promise.all([
    fetchSubjects(),
    fetchHomeworks(),
    fetchLessons()
  ])
  loading.value = false
})

const fetchSubjects = async () => {
  try {
    const schoolId = studentStore.profile?.school_id || studentStore.profile?.schoolId
    if (!schoolId) return
    const res = await adminService.getSubjects(schoolId)
    if (res.data) {
      const map = {}
      res.data.forEach(s => {
        map[s.id] = s.name
      })
      subjectsMap.value = map
    }
  } catch (e) {
    console.error('Error fetching subjects', e)
  }
}

const fetchHomeworks = async () => {
  try {
    const classId = studentStore.profile?.class_id
    if (!classId) return
    const res = await lessonService.getMyHomeworks(classId)
    homeworks.value = res.data || []
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: 'Impossibile caricare i compiti' })
  }
}

const fetchLessons = async () => {
  try {
    const classId = studentStore.profile?.class_id
    if (!classId) return
    const res = await lessonService.getLessons(classId)
    lessons.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const sortedHomeworks = computed(() =>
  [...homeworks.value].sort((a, b) => new Date(a.due_date) - new Date(b.due_date))
)

const calendarEvents = computed(() => {
  const events = new Set()
  homeworks.value.forEach(hw => {
    events.add(date.formatDate(new Date(hw.due_date), 'YYYY-MM-DD'))
  })
  lessons.value.forEach(lesson => {
    events.add(date.formatDate(new Date(lesson.date), 'YYYY-MM-DD'))
  })
  return Array.from(events)
})

const lessonsOnSelectedDate = computed(() => {
  if (!selectedDate.value) return []
  const sel = date.formatDate(new Date(selectedDate.value), 'YYYY-MM-DD')
  return lessons.value.filter(l => date.formatDate(new Date(l.date), 'YYYY-MM-DD') === sel)
})

const homeworksOnSelectedDate = computed(() => {
  if (!selectedDate.value) return []
  const sel = date.formatDate(new Date(selectedDate.value), 'YYYY-MM-DD')
  return homeworks.value.filter(hw => date.formatDate(new Date(hw.due_date), 'YYYY-MM-DD') === sel)
})

const formatDate = (d) => {
  if (!d) return '-'
  return date.formatDate(new Date(d), 'DD/MM/YYYY')
}

const isPast = (d) => new Date(d) < new Date(new Date().toDateString())

const isDueSoon = (d) => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  const due = new Date(d)
  return due.toDateString() === tomorrow.toDateString()
}
</script>
