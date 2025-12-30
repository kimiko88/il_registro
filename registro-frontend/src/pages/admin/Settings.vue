<template>
  <q-page padding>
    <div class="text-h4 q-mb-md">System Settings</div>
    
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-6">
        <ConfigurationPanel title="Email Settings" :loading="store.loading" @submit="saveEmail">
          <q-input v-model="store.config.email.smtp_host" label="SMTP Host" />
          <q-input v-model="store.config.email.smtp_port" label="SMTP Port" />
        </ConfigurationPanel>
      </div>
      
      <div class="col-12 col-md-6">
        <ConfigurationPanel title="Backup Settings" :loading="store.loading" @submit="saveBackup">
          <q-toggle v-model="store.config.backup.enabled" label="Enable Auto Backup" />
          <q-input v-model="store.config.backup.frequency" label="Frequency (Cron)" />
        </ConfigurationPanel>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import ConfigurationPanel from 'src/components/Common/ConfigurationPanel.vue';
import { useSettingsStore } from 'src/stores/settings';
import { onMounted } from 'vue';

const store = useSettingsStore();

onMounted(() => {
  store.fetchSettings();
});

const saveEmail = () => store.updateSettings('email', store.config.email);
const saveBackup = () => store.updateSettings('backup', store.config.backup);
</script>
