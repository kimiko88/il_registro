<template>
  <div class="class-selector q-pa-sm row items-center no-wrap">
    <q-icon name="school" size="sm" class="q-mr-sm text-grey-7" />
    <q-select
      v-model="selectedClassModel"
      :options="classesStore.classes"
      option-value="id"
      option-label="label"
      :label="t('classesPage.currentClass') || 'Classe Attuale'"
      dense
      outlined
      options-dense
      behavior="menu"
      class="col-grow"
      :loading="classesStore.loading"
      bg-color="white"
    >
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps">
          <q-item-section>
            <q-item-label>{{ scope.opt.label || scope.opt.displayName || scope.opt.name }}</q-item-label>
            <q-item-label caption>{{ scope.opt.type }}</q-item-label>
          </q-item-section>
          <q-item-section side v-if="scope.opt.coordinator">
            <q-badge color="accent" label="COORD" />
          </q-item-section>
        </q-item>
      </template>
    </q-select>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useClassesStore } from '@/stores/classes';

const { t } = useI18n();
const classesStore = useClassesStore();

const selectedClassModel = computed({
  get: () => classesStore.classes.find(c => c.id === classesStore.selectedClassId),
  set: (val) => {
    if (val) classesStore.selectClass(val.id);
  }
});

onMounted(() => {
  if (classesStore.classes.length === 0) {
    classesStore.fetchAssignedClasses();
  }
});
</script>

<style scoped>
.class-selector {
  min-width: 150px;
  max-width: 300px;
}
</style>
