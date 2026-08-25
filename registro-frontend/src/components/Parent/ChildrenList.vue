<template>
  <q-list bordered separator>
    <q-item v-for="child in children" :key="child.id">
        <q-item-section avatar>
            <q-avatar size="50px" color="primary" text-color="white">
                <img v-if="child.avatar" :src="child.avatar" :alt="child.firstName" />
                <q-icon v-else name="person" size="32px" />
            </q-avatar>
        </q-item-section>
        
        <q-item-section>
            <q-item-label class="text-h6">{{ child.firstName }} {{ child.lastName }}</q-item-label>
            <q-item-label caption>{{ child.className }} • {{ child.school }}</q-item-label>
        </q-item-section>

        <q-item-section side>
            <div class="row q-gutter-sm">
                <q-btn flat round icon="edit" color="grey" :aria-label="t('common.edit') || 'Modifica'" @click="$emit('edit', child)">
                    <q-tooltip>{{ t('common.edit') || 'Modifica' }}</q-tooltip>
                </q-btn>
                <q-btn flat round icon="delete" color="red" :aria-label="t('common.delete') || 'Elimina'" @click="$emit('remove', child.id)">
                    <q-tooltip>{{ t('common.delete') || 'Elimina' }}</q-tooltip>
                </q-btn>
            </div>
        </q-item-section>
    </q-item>
  </q-list>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps(['children']);
defineEmits(['remove', 'edit']);
</script>
