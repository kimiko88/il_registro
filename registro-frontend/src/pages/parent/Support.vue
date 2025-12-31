<template>
  <q-page class="q-pa-md bg-slate-50">
     <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">Supporto & FAQ</div>
     
     <q-list bordered class="bg-white rounded-borders">
       <q-expansion-item label="Come giustificare un'assenza?" header-class="text-weight-bold">
         <q-card>
           <q-card-section>
             Vai nella sezione "Presenze", trova l'assenza e clicca su "Giustifica".
           </q-card-section>
         </q-card>
       </q-expansion-item>
       
       <q-expansion-item label="Come prenotare un colloquio?" header-class="text-weight-bold">
         <q-card>
           <q-card-section>
             Vai nella sezione "Colloqui", seleziona il docente e scegli uno slot disponibile.
           </q-card-section>
         </q-card>
       </q-expansion-item>
     </q-list>

     <div class="q-mt-xl text-center">
       <q-btn icon="email" label="Contatta Assistenza" color="secondary" />
     </div>
  </q-page>
</template>
