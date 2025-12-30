<template>
  <div class="q-pa-md">
    <div class="text-h6 q-mb-md">Grade Entry</div>
    
    <div class="row q-col-gutter-md q-mb-md">
      <q-select
        class="col-12 col-md-4"
        v-model="selectedClass"
        :options="classes"
        label="Select Class"
        outlined
      />
      <q-select
        class="col-12 col-md-4"
        v-model="selectedSubject"
        :options="subjects"
        label="Select Subject"
        outlined
      />
    </div>

    <q-card v-if="selectedClass">
      <q-card-section>
        <div class="text-subtitle1">Students</div>
      </q-card-section>
      
      <q-list separator>
        <q-item v-for="student in students" :key="student.id">
          <q-item-section>
            <q-item-label>{{ student.name }}</q-item-label>
          </q-item-section>
          
          <q-item-section side>
            <div class="row items-center q-gutter-sm">
              <q-input
                v-model.number="student.grade"
                type="number"
                label="Grade"
                dense
                outlined
                style="width: 80px"
              />
              <q-btn round flat icon="save" color="primary" @click="saveGrade(student)" />
            </div>
          </q-item-section>
        </q-item>
      </q-list>
    </q-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const selectedClass = ref(null)
const selectedSubject = ref(null)

const classes = ref(['1A', '2B', '3C'])
const subjects = ref(['Math', 'History', 'Physics'])
const students = ref([
  { id: 1, name: 'Mario Rossi', grade: null },
  { id: 2, name: 'Luigi Verdi', grade: null }
])

function saveGrade(student) {
  console.log('Save grade', student)
  // API Call
}
</script>
