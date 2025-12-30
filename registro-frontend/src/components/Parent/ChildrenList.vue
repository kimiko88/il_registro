<template>
  <q-list bordered separator>
    <q-item v-for="child in children" :key="child.id">
        <q-item-section avatar>
            <q-avatar size="50px">
                <img :src="child.avatar" />
            </q-avatar>
        </q-item-section>
        
        <q-item-section>
            <q-item-label class="text-h6">{{ child.firstName }} {{ child.lastName }}</q-item-label>
            <q-item-label caption>{{ child.className }} • {{ child.school }}</q-item-label>
        </q-item-section>

        <q-item-section side>
            <div class="row q-gutter-sm">
                <q-btn flat round icon="edit" color="grey" />
                <q-btn flat round icon="delete" color="red" @click="$emit('remove', child.id)" />
            </div>
        </q-item-section>
    </q-item>
  </q-list>
</template>

<script setup>
const props = defineProps(['children']);
</script>
