<template>
  <div v-if="store.children.length > 0">
    <q-btn-dropdown flat no-caps dense :label="store.selectedChild?.firstName || 'Select Child'">
        <template v-slot:icon>
             <q-avatar size="24px">
                <img :src="store.selectedChild?.avatar">
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
                    <q-avatar>
                        <img :src="child.avatar" />
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
import { useChildrenStore } from 'src/stores/children';
import { onMounted } from 'vue';

const store = useChildrenStore();

onMounted(() => {
    if (store.children.length === 0) store.fetchChildren();
});
</script>
