<template>
  <div class="q-pa-md">
    <div class="text-h6 q-mb-md">{{ t('roleDashboards.parentSub') || 'I Miei Figli' }}</div>
    
    <div v-if="store.loading" class="text-center q-pa-lg">
      <q-spinner color="primary" size="32px" />
    </div>

    <div v-else-if="store.children.length === 0" class="text-caption text-grey-6 q-pa-md">
      {{ t('dashboardPage.noChildren') || 'Nessun studente associato al profilo genitore.' }}
    </div>

    <div v-else class="row q-gutter-md">
      <q-card v-for="child in store.children" :key="child.id" class="my-card rounded-xl shadow-xs border">
        <q-card-section>
          <div class="text-h6 text-weight-bold text-slate-800">{{ child.firstName }} {{ child.lastName }}</div>
          <div class="text-subtitle2 text-grey-7">{{ t('udaPage.classLabel') || 'Classe' }}: {{ child.className || '-' }}</div>
        </q-card-section>

        <q-separator />

        <q-card-actions vertical>
          <q-btn flat color="primary" icon="grade" :label="t('classRegister.tabGrades') || 'Voti'" @click="viewGrades(child)" />
          <q-btn flat color="primary" icon="how_to_reg" :label="t('roleDashboards.attendanceRate') || 'Presenze'" @click="viewAttendance(child)" />
        </q-card-actions>
      </q-card>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useChildrenStore } from '@/stores/children'

const { t } = useI18n()
const router = useRouter()
const store = useChildrenStore()

function viewGrades(child) {
  store.selectChild(child.id)
  router.push({ path: '/parent/grades', query: { student_id: child.id } })
}

function viewAttendance(child) {
  store.selectChild(child.id)
  router.push({ path: '/parent/attendance', query: { student_id: child.id } })
}

onMounted(() => {
  if (store.children.length === 0) {
    store.fetchChildren()
  }
})
</script>

<style scoped>
.my-card {
  width: 100%;
  max-width: 300px;
}
</style>
