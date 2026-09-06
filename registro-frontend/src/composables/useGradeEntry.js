import { ref } from 'vue';
import { useGradesStore } from 'src/stores/grades';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';
import { useOfflineSync } from '@/composables/useOfflineSync';

export function useGradeEntry() {
    const gradesStore = useGradesStore();
    const $q = useQuasar();
    const { executeWithOfflineQueue, isOnline } = useOfflineSync();

    const submitting = ref(false);

    const submitGrade = async (gradeData) => {
        const t = i18n?.global?.t;
        if (!gradeData.studentId || !gradeData.value || !gradeData.type) {
            $q.notify({
                type: 'warning',
                message: t ? t('composables.grades.fillRequired') : 'Compilare tutti i campi obbligatori per il voto'
            });
            return false;
        }

        submitting.value = true;
        try {
            if (isOnline?.value === false) {
                await executeWithOfflineQueue({
                    url: '/grades',
                    method: 'post',
                    data: gradeData
                }, { title: `Voto Studente ${gradeData.studentId}` });
                return true;
            }

            await gradesStore.addGrade(gradeData);
            $q.notify({
                type: 'positive',
                message: t ? t('composables.grades.saveSuccess') : 'Voto salvato con successo'
            });
            return true;
        } catch (err) {
            const isNetworkError = !err.response || err.code === 'ERR_NETWORK' || (err.message && /network|fetch|timeout/i.test(err.message));
            if (isNetworkError) {
                await executeWithOfflineQueue({
                    url: '/grades',
                    method: 'post',
                    data: gradeData
                }, { title: `Voto Studente ${gradeData.studentId}` });
                return true;
            }

            $q.notify({
                type: 'negative',
                message: t ? t('composables.grades.saveError') : 'Impossibile salvare il voto'
            });
            return false;
        } finally {
            submitting.value = false;
        }
    };

    return {
        submitGrade,
        submitting
    };
}
