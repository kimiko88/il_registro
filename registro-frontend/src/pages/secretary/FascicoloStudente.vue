<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Fascicolo Digitale Studente</h1>
        <p class="text-sm text-gray-500">Profilo anagrafico e carriera scolastica aggregata</p>
      </div>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">Caricamento fascicolo...</div>
    <div v-else-if="!fascicolo" class="text-center py-12 text-gray-500">Nessun fascicolo trovato.</div>
    <div v-else class="space-y-6">
      <!-- Anagrafica -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-4">
        <h2 class="text-lg font-bold text-gray-900 border-b pb-2">Dati Anagrafici</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div><span class="text-gray-500">Nome:</span> <span class="font-medium text-gray-900">{{ fascicolo.student?.first_name }} {{ fascicolo.student?.last_name }}</span></div>
          <div><span class="text-gray-500">Email:</span> <span class="font-medium text-gray-900">{{ fascicolo.student?.email }}</span></div>
          <div><span class="text-gray-500">Codice Fiscale:</span> <span class="font-medium text-gray-900">{{ fascicolo.student?.fiscal_code || 'N/D' }}</span></div>
        </div>
      </div>

      <!-- Genitori / Tutori -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-4">
        <h2 class="text-lg font-bold text-gray-900 border-b pb-2">Genitori / Tutori Legali</h2>
        <div v-if="!fascicolo.guardians || fascicolo.guardians.length === 0" class="text-sm text-gray-500">Nessun tutore registrato.</div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="g in fascicolo.guardians" :key="g.id" class="p-3 bg-gray-50 rounded-lg text-sm">
            <p class="font-semibold text-gray-900">{{ g.first_name }} {{ g.last_name }}</p>
            <p class="text-xs text-gray-500">{{ g.email }} | {{ g.phone || 'Tel non specificato' }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { userService } from '@/services/userService'

const route = useRoute()
const { t } = useI18n()
const loading = ref(true)
const fascicolo = ref(null)

onMounted(async () => {
  try {
    const studentId = route.params.id
    if (studentId) {
      const res = await userService.getStudentFascicolo(studentId)
      fascicolo.value = res.data
    }
  } catch (err) {
    console.error('Errore durante il recupero del fascicolo:', err)
  } finally {
    loading.value = false
  }
})
</script>
