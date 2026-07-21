<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Gestione Sostituzioni Docenti</h1>
        <p class="text-sm text-gray-500">Visualizza e gestisci le sostituzioni per docenti assenti</p>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
      <div v-if="loading" class="py-8 text-center text-gray-500">Caricamento sostituzioni...</div>
      <div v-else-if="substitutions.length === 0" class="py-8 text-center text-gray-500">Nessuna sostituzione registrata.</div>
      <div v-else class="divide-y divide-gray-100">
        <div v-for="sub in substitutions" :key="sub.id" class="py-4 flex items-center justify-between">
          <div>
            <span class="inline-block px-2 py-1 text-xs font-semibold rounded bg-blue-100 text-blue-800 mr-2">
              Ora {{ sub.hour }}
            </span>
            <span class="font-medium text-gray-900">{{ sub.date ? sub.date.substring(0, 10) : '' }}</span>
            <p class="text-sm text-gray-600 mt-1">
              Docente Assente: {{ sub.absent_teacher_id }} | Sostituto: {{ sub.substitute_teacher_id || 'Non assegnato' }}
            </p>
          </div>
          <span :class="sub.status === 'assigned' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'" class="px-3 py-1 rounded-full text-xs font-medium uppercase">
            {{ sub.status }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { substitutionService } from '@/services/substitutionService'

const loading = ref(true)
const substitutions = ref([])

onMounted(async () => {
  try {
    const res = await substitutionService.listMy()
    substitutions.value = res.data || []
  } catch (err) {
    console.error('Errore durante il recupero delle sostituzioni:', err)
  } finally {
    loading.value = false
  }
})
</script>
