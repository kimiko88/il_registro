import { ref } from 'vue';
import { useColloquiStore } from 'src/stores/colloqui';
import { useQuasar } from 'quasar';
import { date } from 'quasar';
import { i18n } from '@/i18n';

export function useColloquiScheduling() {
    const store = useColloquiStore();
    const $q = useQuasar();

    const generating = ref(false);

    const generateSlots = async (config) => {
        // config: { date, startTime, endTime, duration (mins), break (mins) }
        // Logic to generate individual slots
        generating.value = true;
        try {
            const duration = parseInt(config?.duration, 10);
            if (!duration || duration <= 0) {
                const t = i18n?.global?.t;
                $q.notify({ type: 'warning', message: t ? t('composables.colloqui.invalidDuration') : 'Durata colloquio non valida' });
                return false;
            }

            const slots = [];
            let current = date.addToDate(date.extractDate(`${config.date} ${config.startTime}`, 'YYYY-MM-DD HH:mm'), { minutes: 0 });
            const end = date.extractDate(`${config.date} ${config.endTime}`, 'YYYY-MM-DD HH:mm');

            if (isNaN(current.getTime()) || isNaN(end.getTime()) || current >= end) {
                const t = i18n?.global?.t;
                $q.notify({ type: 'warning', message: t ? t('composables.colloqui.invalidTimeRange') : 'Intervallo orario non valido' });
                return false;
            }

            while (current < end) {
                const slotEnd = date.addToDate(current, { minutes: duration });
                if (slotEnd > end) break;

                slots.push({
                    date: config.date,
                    startTime: date.formatDate(current, 'HH:mm'),
                    endTime: date.formatDate(slotEnd, 'HH:mm')
                });

                current = date.addToDate(slotEnd, { minutes: parseInt(config.break || 0, 10) || 0 });
            }

            await store.createSlots(slots);
            const t = i18n?.global?.t;
            $q.notify({ type: 'positive', message: t ? t('composables.colloqui.slotsCreated', { count: slots.length }) : `${slots.length} disponibilità create con successo` });
            return true;
        } catch (e) {
            const t = i18n?.global?.t;
            $q.notify({ type: 'negative', message: t ? t('composables.colloqui.slotsError') : 'Errore durante la creazione delle disponibilità' });
            return false;
        } finally {
            generating.value = false;
        }
    };

    return { generateSlots, generating };
}
