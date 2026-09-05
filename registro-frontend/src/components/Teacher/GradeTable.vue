<template>
  <div>
    <q-table
      :rows="students"
      :columns="columns"
      row-key="id"
      :loading="loading"
      flat
      bordered
      :pagination="{ rowsPerPage: 0 }"
    >
      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td key="name" :props="props">
            {{ props.row.name }}
          </q-td>
          
          <q-td key="grades" :props="props">
            <div class="row q-gutter-xs">
              <GradeBadge
                v-for="grade in getStudentGrades(props.row.id)"
                :key="grade.id"
                :value="grade.value"
                :type="grade.type"
                :date="grade.date"
                :notes="grade.description"
                :dense="true"
                class="cursor-pointer"
              />
              <q-btn
                round flat dense icon="add" size="xs" color="grey-7"
                @click="$emit('add-grade', props.row)"
              />
            </div>
          </q-td>
          
          <q-td key="average" :props="props">
            <span class="text-weight-bold">{{ calculateAverage(props.row.id) }}</span>
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useGradesStore } from '@/stores/grades';
import GradeBadge from '@/components/Common/GradeBadge.vue';

const { t } = useI18n();
defineProps({
  students: Array,
  loading: Boolean
});

const gradesStore = useGradesStore();

const columns = computed(() => [
  { name: 'name', required: true, label: t('classRegister.tableHeaderStudent') || 'Studente', align: 'left', field: 'name', sortable: true },
  { name: 'grades', label: t('classRegister.tabGrades') || 'Voti', align: 'left' },
  { name: 'average', label: t('roleDashboards.gradeAverage') || 'Media', align: 'center' }
]);

const getStudentGrades = (studentId) => {
  return gradesStore.getGradesByStudent(studentId);
};

const calculateAverage = (studentId) => {
  const grades = getStudentGrades(studentId);
  if (!grades.length) return '-';
  const sum = grades.reduce((acc, g) => acc + g.value, 0);
  return (sum / grades.length).toFixed(1);
};

defineExpose({
    getStudentGrades,
    calculateAverage
})
</script>
