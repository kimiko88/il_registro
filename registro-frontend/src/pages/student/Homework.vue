<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="text-h5 text-weight-bold q-mb-md">
      <q-icon name="assignment" color="orange" class="q-mr-sm" />
      Compiti e Scadenze
    </div>

    <q-card v-if="loading" class="text-center q-pa-xl">
      <q-spinner-dots color="primary" size="60px" />
    </q-card>

    <q-card v-else-if="homeworks.length === 0" class="text-center q-pa-xl text-grey-6">
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
            <q-item-label class="text-weight-bold">{{ hw.description }}</q-item-label>
            <q-item-label caption>
              Consegna:
              <strong :class="isPast(hw.due_date) ? 'text-negative' : 'text-positive'">
                {{ formatDate(hw.due_date) }}
              </strong>
              <q-badge v-if="isPast(hw.due_date)" color="negative" class="q-ml-sm">Scaduto</q-badge>
              <q-badge v-else-if="isDueSoon(hw.due_date)" color="warning" class="q-ml-sm">Domani</q-badge>
            </q-item-label>
          </q-item-section>
        </q-item>
      </q-card>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, date } from 'quasar'
import { useStudentStore } from 'src/stores/student'
import { lessonService } from 'src/services/lessonService'

const $q = useQuasar()
const studentStore = useStudentStore()

const homeworks = ref([])
const loading = ref(true)

onMounted(async () => {
  await studentStore.fetchProfile()
  await fetchHomeworks()
})

const fetchHomeworks = async () => {
  loading.value = true
  try {
    const classId = studentStore.profile?.class_id
    if (!classId) return
    const res = await lessonService.getMyHomeworks(classId)
    homeworks.value = res.data || []
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: 'Impossibile caricare i compiti' })
  } finally {
    loading.value = false
  }
}

const sortedHomeworks = computed(() =>
  [...homeworks.value].sort((a, b) => new Date(a.due_date) - new Date(b.due_date))
)

const formatDate = (d) => date.formatDate(new Date(d), 'DD/MM/YYYY')
const isPast = (d) => new Date(d) < new Date(new Date().toDateString())
const isDueSoon = (d) => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  const due = new Date(d)
  return due.toDateString() === tomorrow.toDateString()
}
</script>
