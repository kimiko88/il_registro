import { ref } from 'vue';
import { useGradesStore } from 'src/stores/grades';
import { useQuasar } from 'quasar';

export function useGradeEntry() {
    const gradesStore = useGradesStore();
    const $q = useQuasar();

    const submitting = ref(false);

    const submitGrade = async (gradeData) => {
        if (!gradeData.studentId || !gradeData.value || !gradeData.type) {
            $q.notify({
                type: 'warning',
                message: 'Please fill all required fields'
            });
            return false;
        }

        submitting.value = true;
        try {
            await gradesStore.addGrade(gradeData);
            $q.notify({
                type: 'positive',
                message: 'Grade saved successfully'
            });
            return true;
        } catch (err) {
            $q.notify({
                type: 'negative',
                message: 'Failed to save grade'
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
