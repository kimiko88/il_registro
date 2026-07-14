<template>
  <q-page class="q-pa-md bg-grey-1">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <div class="text-h5 text-weight-bold">
          <q-icon name="folder_shared" color="primary" class="q-mr-sm" />
          Materiale Didattico
        </div>
        <div class="text-caption text-grey">Sfoglia e scarica le dispense e i materiali di studio condivisi dai tuoi docenti</div>
      </div>
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md shadow-1">
      <q-card-section class="row items-center q-gutter-md q-py-sm">
        <q-select
          v-model="selectedSubject"
          :options="subjectsOptions"
          option-value="id"
          option-label="name"
          emit-value map-options
          label="Filtra per Materia"
          dense outlined
          style="min-width:200px"
          clearable
        />
      </q-card-section>
    </q-card>

    <!-- Materials List -->
    <q-card v-slot:default v-if="loading" class="text-center q-pa-xl shadow-1">
      <q-spinner-dots color="primary" size="60px" />
    </q-card>

    <div v-else>
      <q-card v-if="filteredMaterials.length === 0" class="text-center q-pa-xl text-grey-6 shadow-1">
        <q-icon name="folder_open" size="80px" class="q-mb-md" />
        <div class="text-h6">Nessun materiale disponibile</div>
        <div class="text-caption">Non è stato ancora condiviso alcun file per questa classe</div>
      </q-card>

      <div v-else class="row q-col-gutter-md">
        <div v-for="mat in filteredMaterials" :key="mat.id" class="col-12 col-sm-6 col-md-4">
          <q-card class="shadow-soft hover-card full-height flex flex-column justify-between">
            <q-card-section>
              <div class="row items-center justify-between q-mb-xs">
                <q-chip dense size="sm" color="primary-1" text-color="primary" class="text-weight-bold">
                  {{ subjectsMap[mat.subject_id] || mat.subject_id }}
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
                icon="download"
                label="Scarica Allegato"
                no-caps
                class="full-width"
                @click="openLink(mat.attachment_url)"
              />
              <span v-else class="text-caption text-grey-5 text-center full-width">Nessun file allegato</span>
            </q-card-actions>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar, date } from 'quasar'
import { useStudentStore } from 'src/stores/student'
import didacticService from 'src/services/didacticService'
import adminService from 'src/services/adminService'

const $q = useQuasar()
const studentStore = useStudentStore()

const loading = ref(true)
const selectedSubject = ref(null)
const materials = ref([])
const subjectsMap = ref({})
const subjectsOptions = ref([])

onMounted(async () => {
  await studentStore.fetchProfile()
  await Promise.all([
    fetchSubjects(),
    fetchMaterials()
  ])
  loading.value = false
})

const fetchSubjects = async () => {
  try {
    const schoolId = studentStore.profile?.school_id || studentStore.profile?.schoolId
    if (!schoolId) return
    const res = await adminService.getSubjects(schoolId)
    if (res.data) {
      subjectsOptions.value = res.data
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

const fetchMaterials = async () => {
  try {
    const classId = studentStore.profile?.class_id
    if (!classId) return
    const res = await didacticService.getMaterials(classId)
    materials.value = res.data || []
  } catch (e) {
    console.error(e)
    $q.notify({ type: 'negative', message: 'Impossibile caricare il materiale didattico' })
  }
}

const filteredMaterials = computed(() => {
  if (!selectedSubject.value) return materials.value
  return materials.value.filter(m => m.subject_id === selectedSubject.value)
})

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
  -webkit-box-orient: vertical;  
  overflow: hidden;
}
</style>
