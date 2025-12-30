<template>
  <div class="q-pa-md">
    <div class="text-h6 q-mb-md">Attendance</div>

    <div class="row q-mb-md justify-between">
      <q-select
        v-model="selectedClass"
        :options="classes"
        label="Select Class"
        outlined
        style="width: 200px"
      />
      <q-date v-model="date" landscape minimal />
    </div>

    <q-table
      v-if="selectedClass"
      :rows="students"
      :columns="columns"
      row-key="id"
      flat bordered
    >
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-btn-toggle
            v-model="props.row.status"
            toggle-color="primary"
            flat
            :options="[
              {label: 'P', value: 'present'},
              {label: 'A', value: 'absent'},
              {label: 'L', value: 'late'}
            ]"
          />
        </q-td>
      </template>
    </q-table>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const selectedClass = ref(null)
const date = ref('2023/10/01')
const classes = ref(['1A', '2B'])

const students = ref([
  { id: 1, name: 'Mario Rossi', status: 'present' },
  { id: 2, name: 'Luigi Verdi', status: 'absent' }
])

const columns = [
  { name: 'name', label: 'Name', field: 'name', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' }
]
</script>
