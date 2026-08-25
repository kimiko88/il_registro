import { ref } from 'vue';
import { useGradesStore } from 'src/stores/grades';
import { useQuasar } from 'quasar';
import { i18n } from '@/i18n';

export function useGradeEntry() {
    const gradesStore = useGradesStore();
    const $q = useQuasar();

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
            await gradesStore.addGrade(gradeData);
            $q.notify({
                type: 'positive',
                message: t ? t('composables.grades.saveSuccess') : 'Voto salvato con successo'
            });
            return true;
        } catch (err) {
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
