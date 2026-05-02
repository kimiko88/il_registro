<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-xl">
      <div class="col">
        <h1 class="text-h3 text-weight-bold text-outfit bg-clip-text text-transparent bg-gradient-premium q-my-none" style="display: inline-block;">
          Impostazioni
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-sm">Configura le preferenze del sistema</div>
      </div>
    </div>
    
     <q-card class="glass-card shadow-soft overflow-hidden">
       <q-list separator>
        <q-item-label header class="text-uppercase letter-spacing-1 text-slate-400 text-weight-bold">Generale</q-item-label>
        
        <q-item clickable v-ripple class="q-py-md">
          <q-item-section avatar>
            <div class="bg-indigo-50 text-indigo-600 q-pa-sm rounded-lg">
              <q-icon name="language" size="24px" />
            </div>
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-bold text-slate-700">Lingua</q-item-label>
            <q-item-label caption>Italiano (Italia)</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-icon name="chevron_right" color="slate-300" />
          </q-item-section>
        </q-item>
  
         <q-item clickable v-ripple class="q-py-md">
          <q-item-section avatar>
            <div class="bg-indigo-50 text-indigo-600 q-pa-sm rounded-lg">
              <q-icon name="security" size="24px" />
            </div>
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-bold text-slate-700">Sicurezza</q-item-label>
            <q-item-label caption>Policy password e 2FA</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-icon name="chevron_right" color="slate-300" />
          </q-item-section>
        </q-item>

         <q-item clickable v-ripple class="q-py-md">
          <q-item-section avatar>
            <div class="bg-indigo-50 text-indigo-600 q-pa-sm rounded-lg">
              <q-icon name="notifications" size="24px" />
            </div>
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-bold text-slate-700">Notifiche</q-item-label>
            <q-item-label caption>Preferenze email e push</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-icon name="chevron_right" color="slate-300" />
          </q-item-section>
        </q-item>
      </q-list>
     </q-card>
  </q-page>
</template>

<style scoped>
.bg-clip-text {
    -webkit-background-clip: text;
    background-clip: text;
}

.bg-gradient-premium {
    background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
}

.letter-spacing-1 {
    letter-spacing: 1px;
}
</style>
