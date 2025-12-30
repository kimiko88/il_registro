<template>
  <q-dialog v-model="isOpen">
    <q-card style="min-width: 600px">
        <q-card-section>
            <div class="text-h6">Select Template</div>
        </q-card-section>
        
        <q-card-section>
            <q-list bordered separator>
                <q-item v-for="tpl in templates" :key="tpl.id" clickable @click="select(tpl)">
                    <q-item-section>
                        <q-item-label>{{ tpl.name }}</q-item-label>
                        <q-item-label caption>{{ tpl.type }}</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                        <q-icon name="chevron_right" />
                    </q-item-section>
                </q-item>
            </q-list>
        </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useDocumentsStore } from 'src/stores/documents';

const props = defineProps(['modelValue']);
const emit = defineEmits(['update:modelValue', 'selected']);

const store = useDocumentsStore();
const templates = ref([]);
const isOpen = ref(false);

watch(() => props.modelValue, (val) => isOpen.value = val);
watch(isOpen, (val) => emit('update:modelValue', val));

const select = (tpl) => {
    emit('selected', tpl);
    isOpen.value = false;
};

onMounted(async () => {
    templates.value = await store.fetchTemplates();
});
</script>
