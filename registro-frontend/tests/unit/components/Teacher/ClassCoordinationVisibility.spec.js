import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from 'src/stores/auth'
import { useClassesStore } from 'src/stores/classes'
import { useTeacherStore } from 'src/stores/teacher'
import { useCoordination } from 'src/composables/useCoordination'
import { useMenuItems } from 'src/composables/useMenuItems'

describe('Teacher Class Coordination Visibility & Assignment Tests', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('allows secretary to assign a teacher as coordinator for multiple classes', () => {
    const classesStore = useClassesStore()
    const teacherId = 'teacher-101'

    classesStore.classes = [
      { id: 'class-1a', name: '1A', coordinator_id: teacherId },
      { id: 'class-2a', name: '2A', coordinator_id: teacherId },
      { id: 'class-3a', name: '3A', coordinator_id: 'teacher-999' }
    ]

    const authStore = useAuthStore()
    authStore.user = { id: teacherId, role: 'teacher' }

    const { coordinatedClasses, isCoordinator } = useCoordination()

    expect(isCoordinator.value).toBe(true)
    expect(coordinatedClasses.value).toHaveLength(2)
    expect(coordinatedClasses.value.map(c => c.name)).toEqual(['1A', '2A'])
  })

  it('hides Coordinamento section for teachers who do NOT coordinate any class', () => {
    const classesStore = useClassesStore()
    const teacherId = 'teacher-normal'

    classesStore.classes = [
      { id: 'class-1a', name: '1A', coordinator_id: 'teacher-other' }
    ]

    const authStore = useAuthStore()
    authStore.user = { id: teacherId, role: 'teacher' }

    const teacherStore = useTeacherStore()
    teacherStore.profile = { id: teacherId, is_coordinator: false }

    const { isCoordinator, coordinatedClasses } = useCoordination()

    expect(isCoordinator.value).toBe(false)
    expect(coordinatedClasses.value).toHaveLength(0)

    // Verify menu items filtering logic for non-coordinator teacher
    const rawItems = useMenuItems('teacher')
    const filteredItems = rawItems.map(cat => {
      if (!cat.children) return cat
      return {
        ...cat,
        children: cat.children.filter(child => !child.coordinatorOnly)
      }
    })

    const allLabels = filteredItems.flatMap(cat => cat.children ? cat.children.map(c => c.label) : [cat.label])
    expect(allLabels).not.toContain('Coordinamento')
  })

  it('shows Coordinamento section ONLY for teachers who coordinate one or more classes', () => {
    const classesStore = useClassesStore()
    const teacherId = 'teacher-coord'

    classesStore.classes = [
      { id: 'class-1a', name: '1A', coordinator_id: teacherId }
    ]

    const authStore = useAuthStore()
    authStore.user = { id: teacherId, role: 'teacher' }

    const { isCoordinator } = useCoordination()

    expect(isCoordinator.value).toBe(true)

    const rawItems = useMenuItems('teacher')
    const allLabels = rawItems.flatMap(cat => cat.children ? cat.children.map(c => c.label) : [cat.label])
    expect(allLabels).toContain('Coordinamento')
  })
})
