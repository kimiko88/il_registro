import { setActivePinia, createPinia } from 'pinia';
import { useColloquiScheduling } from 'src/composables/useColloquiScheduling';
import { useColloquiStore } from 'src/stores/colloqui';
import { describe, it, expect, vi, beforeEach } from 'vitest';

// Mock Quasar
const mockNotify = vi.fn();
vi.mock('quasar', () => ({
    useQuasar: () => ({
        notify: mockNotify
    }),
    date: {
        formatDate: (d, f) => {
            // Simple mock for test
            if (d instanceof Date) return d.toISOString().substr(11, 5); // HH:mm
            return '2025-01-01';
        },
        addToDate: (d, opts) => {
            const newDate = new Date(d);
            if (opts.minutes) newDate.setMinutes(newDate.getMinutes() + opts.minutes);
            return newDate;
        },
        extractDate: (str) => {
            // very rough mock since Quasar's is complex
            // Expecting 'YYYY-MM-DD HH:mm'
            const [d, t] = str.split(' ');
            const [y, m, day] = d.split('-');
            const [h, min] = t.split(':');
            return new Date(y, m - 1, day, h, min);
        }
    }
}));

describe('useColloquiScheduling', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        mockNotify.mockClear();
    });

    it('generates slots correctly', async () => {
        const { generateSlots } = useColloquiScheduling();
        const store = useColloquiStore();

        // Mock store
        store.createSlots = vi.fn();

        const config = {
            date: '2025-01-01',
            startTime: '15:00',
            endTime: '16:00', // Should fit 4 15-min slots
            duration: 15,
            break: 0
        };

        await generateSlots(config);

        expect(store.createSlots).toHaveBeenCalled();
        const slots = store.createSlots.mock.calls[0][0];
        expect(slots.length).toBe(4);
        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }));
    });
});
