import { ref } from 'vue';
import { useColloquiStore } from 'src/stores/colloqui';
import { useQuasar } from 'quasar';
import { date } from 'quasar';

export function useColloquiScheduling() {
    const store = useColloquiStore();
    const $q = useQuasar();

    const generating = ref(false);

    const generateSlots = async (config) => {
        // config: { date, startTime, endTime, duration (mins), break (mins) }
        // Logic to generate individual slots
        generating.value = true;
        try {
            const slots = [];
            let current = date.addToDate(date.extractDate(`${config.date} ${config.startTime}`, 'YYYY-MM-DD HH:mm'), { minutes: 0 });
            const end = date.extractDate(`${config.date} ${config.endTime}`, 'YYYY-MM-DD HH:mm');

            while (current < end) {
                const slotEnd = date.addToDate(current, { minutes: parseInt(config.duration) });
                if (slotEnd > end) break;

                slots.push({
                    date: config.date,
                    startTime: date.formatDate(current, 'HH:mm'),
                    endTime: date.formatDate(slotEnd, 'HH:mm')
                });

                current = date.addToDate(slotEnd, { minutes: parseInt(config.break || 0) });
            }

            await store.createSlots(slots);
            $q.notify({ type: 'positive', message: `${slots.length} slots created` });
            return true;
        } catch (e) {
            $q.notify({ type: 'negative', message: 'Error creating slots' });
            return false;
        } finally {
            generating.value = false;
        }
    };

    return { generateSlots, generating };
}
