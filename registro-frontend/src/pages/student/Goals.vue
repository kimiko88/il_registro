<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Obiettivi & Badge Studente</h1>
        <p class="text-sm text-gray-500">Monitora i tuoi traguardi personali e i badge sbloccati</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="goal in goals" :key="goal.id" class="bg-white rounded-xl shadow-sm border border-gray-100 p-5 space-y-3">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-full bg-amber-100 text-amber-600 flex items-center justify-center font-bold text-lg">
            🏆
          </div>
          <div>
            <h3 class="font-semibold text-gray-900">{{ goal.title }}</h3>
            <span class="text-xs text-gray-500 uppercase font-medium">{{ goal.category || 'Generale' }}</span>
          </div>
        </div>
        <p class="text-sm text-gray-600">{{ goal.description }}</p>
        <div class="flex items-center justify-between pt-2 border-t border-gray-100 text-xs">
          <span class="font-medium text-amber-600">+{{ goal.points }} punti</span>
          <span :class="goal.status === 'completed' ? 'bg-green-100 text-green-800' : 'bg-blue-100 text-blue-800'" class="px-2 py-0.5 rounded-full font-semibold">
            {{ goal.status }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { studentGoalService } from '@/services/studentGoalService'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const goals = ref([])

onMounted(async () => {
  try {
    const userId = authStore.user?.id
    if (userId) {
      const res = await studentGoalService.listByStudent(userId)
      goals.value = res.data || []
    }
  } catch (err) {
    console.error('Errore durante il recupero degli obiettivi:', err)
  }
})
</script>
