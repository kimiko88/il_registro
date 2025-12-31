<template>
  <q-page class="q-pa-md">
    <div class="text-h4 text-weight-bold q-mb-md">Impostazioni</div>
    
     <q-list bordered class="rounded-borders bg-white">
      <q-item-label header>Generale</q-item-label>
      
      <q-item clickable v-ripple>
        <q-item-section avatar>
          <q-icon name="language" />
        </q-item-section>
        <q-item-section>
          <q-item-label>Lingua</q-item-label>
          <q-item-label caption>Italiano</q-item-label>
        </q-item-section>
      </q-item>

      <q-separator />
      
       <q-item clickable v-ripple>
        <q-item-section avatar>
          <q-icon name="security" />
        </q-item-section>
        <q-item-section>
          <q-item-label>Sicurezza</q-item-label>
          <q-item-label caption>Policy password e 2FA</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
  </q-page>
</template>
