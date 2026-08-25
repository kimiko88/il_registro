<template>
  <div v-if="store.children.length > 0">
    <q-btn-dropdown flat no-caps dense :label="store.selectedChild?.firstName || (t('roleDashboards.selectChild') || 'Seleziona Figlio')">
        <template v-slot:icon>
             <q-avatar size="24px" color="primary" text-color="white">
                <img v-if="store.selectedChild?.avatar" :src="store.selectedChild.avatar">
                <q-icon v-else name="person" size="16px" />
             </q-avatar>
        </template>
        
        <q-list>
            <q-item 
                v-for="child in store.children" 
                :key="child.id" 
                clickable 
                v-close-popup 
                @click="store.selectChild(child.id)"
                :active="child.id === store.selectedChildId"
            >
                <q-item-section avatar>
                    <q-avatar color="primary" text-color="white">
                        <img v-if="child.avatar" :src="child.avatar" />
                        <q-icon v-else name="person" size="20px" />
                    </q-avatar>
                </q-item-section>
                <q-item-section>
                    <q-item-label>{{ child.firstName }} {{ child.lastName }}</q-item-label>
                    <q-item-label caption>{{ child.className }}</q-item-label>
                </q-item-section>
            </q-item>
        </q-list>
    </q-btn-dropdown>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import { useChildrenStore } from '@/stores/children';
import { onMounted } from 'vue';

const { t } = useI18n();
const store = useChildrenStore();

onMounted(() => {
    if (store.children.length === 0) store.fetchChildren();
});
</script>
