import { computed } from 'vue';
import { useGradesStore } from 'src/stores/grades';

export function useMyGrades() {
    const store = useGradesStore();

    const gradesBySubject = computed(() => {
        const grouped = {};
        store.grades.forEach(g => {
            if (!grouped[g.subject]) grouped[g.subject] = [];
            grouped[g.subject].push(g);
        });
        return grouped;
    });

    const averages = computed(() => {
        const avgs = {};
        for (const subject in gradesBySubject.value) {
            const grades = gradesBySubject.value[subject];
            const sum = grades.reduce((a, b) => a + b.value, 0);
            avgs[subject] = (sum / grades.length).toFixed(1);
        }
        return avgs;
    });

    const getTrend = (subject) => {
        const grades = gradesBySubject.value[subject];
        if (!grades || grades.length < 2) return 'stable';
        // Simple logic: compare last 2 grades
        // Assuming grades are sorted by date (mock data is, but ideally should sort)
        // Here we'll just check last two added
        const last = grades[grades.length - 1].value;
        const prev = grades[grades.length - 2].value;
        if (last > prev) return 'up';
        if (last < prev) return 'down';
        return 'stable';
    };

    return {
        gradesBySubject,
        averages,
        getTrend,
        loading: computed(() => store.loading)
    };
}
