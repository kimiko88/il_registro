import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'

const mockRouterPush = vi.fn().mockResolvedValue(undefined)

vi.mock('vue-router', () => ({
    useRouter: () => ({ push: mockRouterPush })
}))

const mockAuthStore = {
    userRole: 'teacher',
    isAuthenticated: true
}
vi.mock('@/stores/auth', () => ({ useAuthStore: () => mockAuthStore }))
vi.mock('src/stores/auth', () => ({ useAuthStore: () => mockAuthStore }))

const mockThemeStore = {
    keyboardShortcutsHelpOpen: false,
    toggleKeyboardShortcutsHelp: vi.fn(),
    toggleReadingRuler: vi.fn()
}
vi.mock('@/stores/theme', () => ({ useThemeStore: () => mockThemeStore }))
vi.mock('src/stores/theme', () => ({ useThemeStore: () => mockThemeStore }))

const mockStop = vi.fn()
const mockSpeak = vi.fn()
const isPlayingRef = ref(false)

const mockSpeech = {
    speak: mockSpeak,
    stop: mockStop,
    isPlaying: isPlayingRef
}
vi.mock('@/composables/useSpeechSynthesis', () => ({ useSpeechSynthesis: () => mockSpeech }))
vi.mock('src/composables/useSpeechSynthesis', () => ({ useSpeechSynthesis: () => mockSpeech }))

describe('useGlobalKeyboardShortcuts', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
        isPlayingRef.value = false
        mockAuthStore.userRole = 'teacher'
        mockThemeStore.keyboardShortcutsHelpOpen = false
    })

    async function getHandler() {
        const { useGlobalKeyboardShortcuts } = await import('@/composables/useGlobalKeyboardShortcuts')
        const { handleKeyDown } = useGlobalKeyboardShortcuts()
        return handleKeyDown
    }

    function makeEvent(key, options = {}) {
        return {
            key,
            altKey: options.altKey ?? false,
            ctrlKey: options.ctrlKey ?? false,
            metaKey: options.metaKey ?? false,
            preventDefault: vi.fn()
        }
    }

    it('Alt+1 navigates to home route', async () => {
        const handleKeyDown = await getHandler()
        handleKeyDown(makeEvent('1', { altKey: true }))
        expect(mockRouterPush).toHaveBeenCalledWith('/')
    })

    it('Alt+v navigates to teacher grades route', async () => {
        const handleKeyDown = await getHandler()
        handleKeyDown(makeEvent('v', { altKey: true }))
        expect(mockRouterPush).toHaveBeenCalledWith('/teacher/grades')
    })

    it('Alt+p navigates to teacher attendance route', async () => {
        const handleKeyDown = await getHandler()
        handleKeyDown(makeEvent('p', { altKey: true }))
        expect(mockRouterPush).toHaveBeenCalledWith('/teacher/attendance')
    })

    it('Alt+a navigates to teacher agenda', async () => {
        const handleKeyDown = await getHandler()
        handleKeyDown(makeEvent('a', { altKey: true }))
        expect(mockRouterPush).toHaveBeenCalledWith('/teacher/agenda')
    })

    it('? toggles keyboard shortcuts help when not in input', async () => {
        const handleKeyDown = await getHandler()
        handleKeyDown(makeEvent('?'))
        expect(mockThemeStore.toggleKeyboardShortcutsHelp).toHaveBeenCalledOnce()
    })

    it('Escape stops speech if playing', async () => {
        isPlayingRef.value = true

        const handleKeyDown = await getHandler()
        handleKeyDown(makeEvent('Escape'))
        expect(mockStop).toHaveBeenCalledOnce()
    })

    it('Alt+r toggles reading ruler', async () => {
        const handleKeyDown = await getHandler()
        handleKeyDown(makeEvent('r', { altKey: true }))
        expect(mockThemeStore.toggleReadingRuler).toHaveBeenCalledOnce()
    })

    it('Alt+Ctrl combination does NOT trigger shortcuts', async () => {
        const handleKeyDown = await getHandler()
        handleKeyDown(makeEvent('1', { altKey: true, ctrlKey: true }))
        // Should not navigate with AltCtrl combo
        expect(mockRouterPush).not.toHaveBeenCalled()
    })

    it('Shortcut keys are ignored when focus is in an INPUT', async () => {
        const input = document.createElement('input')
        document.body.appendChild(input)
        input.focus()

        const handleKeyDown = await getHandler()
        const event = { ...makeEvent('?'), target: input }
        // Simulate activeElement being the input
        Object.defineProperty(document, 'activeElement', { get: () => input, configurable: true })

        handleKeyDown(event)

        expect(mockThemeStore.toggleKeyboardShortcutsHelp).not.toHaveBeenCalled()
        document.body.removeChild(input)
    })
})
