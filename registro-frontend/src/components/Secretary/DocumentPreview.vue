<template>
  <q-card style="width: min(700px, 95vw); max-width: 95vw; height: 80vh">
    <q-card-section class="row items-center">
      <div class="text-h6">{{ doc?.title }}</div>
      <q-space />
      <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
    </q-card-section>
    <q-separator />
    <q-card-section class="scroll" style="height: calc(100% - 50px)">
        <div v-if="doc?.url" class="fit">
            <iframe :src="doc.url" width="100%" height="100%" frameborder="0"></iframe>
        </div>
        <div v-else class="text-center q-pa-xl text-grey">
            {{ t('documentsPage.noPreview') || 'Nessuna anteprima disponibile' }}
        </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
defineProps({
    doc: Object
});
</script>
