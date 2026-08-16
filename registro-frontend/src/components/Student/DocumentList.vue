<template>
  <q-list bordered separator class="bg-white rounded-borders">
    <q-item v-for="doc in documents" :key="doc.id" clickable @click="$emit('download', doc)">
        <q-item-section avatar>
            <q-icon name="description" color="primary" />
        </q-item-section>
        <q-item-section>
            <q-item-label>{{ doc.title }}</q-item-label>
            <q-item-label caption>{{ doc.date }} • {{ doc.type }}</q-item-label>
        </q-item-section>
        <q-item-section side>
            <q-btn flat round icon="download" color="grey" />
        </q-item-section>
    </q-item>
  </q-list>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps(['documents']);
</script>
