import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import schoolSettingsService from '@/services/schoolSettingsService'
import { useAuthStore } from './auth'

export const useSchoolSettingsStore = defineStore('schoolSettings', () => {
  const authStore = useAuthStore()

  // Default canteen status loaded from localStorage or default true for development
  const savedCanteen = localStorage.getItem('school_canteen_enabled')
  const canteenEnabled = ref(savedCanteen !== null ? savedCanteen === 'true' : true)
  const loading = ref(false)
  const settings = ref({
    require_principal_approval_for_notes: false,
    allow_parents_view_grades: true,
    allow_students_view_class_averages: true,
    require_mfa_for_staff: false,
    lock_scrutiny_editing_after_validation: true,
    enable_substitute_notifications: true,
    enable_canteen_service: canteenEnabled.value
  })

  const currentUserRole = computed(() => {
    return (authStore.userRole || authStore.user?.role || '').toLowerCase()
  })

  // Only Principal (Dirigente), DSGA, or Admin can activate/deactivate the school canteen
  const canManageCanteen = computed(() => {
    const role = currentUserRole.value
    return ['superadmin', 'admin', 'principal', 'dirigente_scolastico', 'dsga'].includes(role)
  })

  async function fetchSettings() {
    loading.value = true
    try {
      const res = await schoolSettingsService.getSettings()
      if (res?.data) {
        settings.value = { ...settings.value, ...res.data }
        if (typeof res.data.enable_canteen_service === 'boolean') {
          canteenEnabled.value = res.data.enable_canteen_service
          localStorage.setItem('school_canteen_enabled', String(canteenEnabled.value))
        }
      }
    } catch {
      // Keep local cached settings if backend offline
    } finally {
      loading.value = false
    }
  }

  async function setCanteenEnabled(enabled) {
    canteenEnabled.value = Boolean(enabled)
    settings.value.enable_canteen_service = canteenEnabled.value
    localStorage.setItem('school_canteen_enabled', String(canteenEnabled.value))

    try {
      await schoolSettingsService.updateSettings({
        ...settings.value,
        enable_canteen_service: canteenEnabled.value
      })
    } catch {
      // Local state is preserved
    }
  }

  return {
    canteenEnabled,
    settings,
    loading,
    canManageCanteen,
    fetchSettings,
    setCanteenEnabled
  }
})

export default useSchoolSettingsStore
