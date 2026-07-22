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
              <q-badge
                v-for="grade in getStudentGrades(props.row.id)"
                :key="grade.id"
                :color="getGradeColor(grade.value)"
                class="cursor-pointer"
              >
                {{ grade.value }}
                <q-tooltip>
                  {{ grade.type }} - {{ grade.date }}<br>{{ grade.description }}
                </q-tooltip>
              </q-badge>
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
import { useGradesStore } from 'src/stores/grades';

defineProps({
  students: Array,
  loading: Boolean
});

const gradesStore = useGradesStore();

const columns = [
  { name: 'name', required: true, label: 'Student', align: 'left', field: 'name', sortable: true },
  { name: 'grades', label: 'Grades', align: 'left' },
  { name: 'average', label: 'Average', align: 'center' }
];

const getStudentGrades = (studentId) => {
  return gradesStore.getGradesByStudent(studentId);
};

const getGradeColor = (val) => {
  if (val >= 9) return 'green-7';
  if (val >= 6) return 'blue-7';
  return 'red-7';
};

const calculateAverage = (studentId) => {
  const grades = getStudentGrades(studentId);
  if (!grades.length) return '-';
  const sum = grades.reduce((acc, g) => acc + g.value, 0);
  return (sum / grades.length).toFixed(1);
};

defineExpose({
    getStudentGrades,
    getGradeColor,
    calculateAverage
})
</script>
