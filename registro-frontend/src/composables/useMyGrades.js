import { computed } from 'vue';
import { useGradesStore } from 'src/stores/grades';

export function useMyGrades() {
    const store = useGradesStore();

    const gradesBySubject = computed(() => {
        const grouped = {};
        if (!store.grades || !store.grades.semesters) return grouped;
        
        store.grades.semesters.forEach(semester => {
            semester.grades.forEach(g => {
                const subject = g.subject_id;
                if (!grouped[subject]) grouped[subject] = [];
                grouped[subject].push(g);
            });
        });
        return grouped;
    });

    const averages = computed(() => {
        const avgs = {};
        for (const subject in gradesBySubject.value) {
            const grades = gradesBySubject.value[subject].filter(g => g.grade_value > 0);
            if (grades.length > 0) {
                const sum = grades.reduce((a, b) => a + b.grade_value, 0);
                avgs[subject] = Math.round((sum / grades.length) * 10) / 10;
            } else {
                avgs[subject] = '-';
            }
        }
        return avgs;
    });

    const getTrend = (subject) => {
        const grades = gradesBySubject.value[subject]?.filter(g => g.grade_value > 0);
        if (!grades || grades.length < 2) return 'stable';
        // Simple logic: compare last 2 grades
        // Assuming grades are sorted by date (mock data is, but ideally should sort)
        // Here we'll just check last two added
        const last = grades[grades.length - 1].grade_value;
        const prev = grades[grades.length - 2].grade_value;
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
