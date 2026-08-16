<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-lg">
       <div class="col">
           <div class="text-h4">Orientamento in Uscita</div>
           <div class="text-subtitle1 text-grey">Eventi universitari e opportunità post-diploma</div>
       </div>
    </div>

    <div class="row q-col-gutter-lg">
        <!-- Event List -->
        <div class="col-12 col-md-8">
            <q-tabs v-model="tab" dense class="text-grey q-mb-md" active-color="primary" indicator-color="primary" align="left">
                <q-tab name="upcoming" :label="t('verbaliPage.agenda') || 'Prossimi Eventi'" />
                <q-tab name="registered" :label="t('colloquiPage.bookings') || 'I Miei Eventi'" />
                <q-tab name="past" :label="t('verbaliPage.resolutions') || 'Passati'" />
            </q-tabs>

            <q-tab-panels v-model="tab" animated class="bg-transparent">
                <q-tab-panel name="upcoming" class="q-pa-none">
                     <div class="row q-col-gutter-md">
                         <div class="col-12 col-md-6" v-for="event in upcomingEvents" :key="event.id">
                             <q-card>
                                 <q-img :src="event.image" style="height: 140px">
                                     <div class="absolute-bottom text-subtitle2 flex justify-between items-center">
                                         <span>{{ event.date }}</span>
                                         <q-chip color="white" text-color="black" size="sm">{{ event.university }}</q-chip>
                                     </div>
                                 </q-img>
                                 <q-card-section>
                                     <div class="text-h6">{{ event.title }}</div>
                                     <div class="text-caption text-grey">{{ event.description }}</div>
                                 </q-card-section>
                                 <q-separator />
                                 <q-card-actions align="right">
                                     <q-btn flat :label="t('udaPage.detailTitle') || 'Dettagli'" color="primary" />
                                     <q-btn color="primary" :label="t('colloquiPage.booked') || 'Iscriviti'" icon="event_available" @click="register(event)" />
                                 </q-card-actions>
                             </q-card>
                         </div>
                     </div>
                </q-tab-panel>

                 <q-tab-panel name="registered" class="q-pa-none">
                     <q-list bordered separator class="bg-white rounded-borders">
                         <q-item v-for="event in registeredEvents" :key="event.id">
                             <q-item-section avatar>
                                 <q-icon name="event" color="green" size="md" />
                             </q-item-section>
                             <q-item-section>
                                 <q-item-label class="text-weight-bold">{{ event.title }}</q-item-label>
                                 <q-item-label caption>{{ event.university }} - {{ event.date }}</q-item-label>
                             </q-item-section>
                             <q-item-section side>
                                 <q-chip color="green" text-color="white" icon="check">Iscritto</q-chip>
                             </q-item-section>
                         </q-item>
                     </q-list>
                </q-tab-panel>
            </q-tab-panels>
        </div>

        <!-- Sidebar -->
        <div class="col-12 col-md-4">
             <q-card class="bg-primary text-white text-center q-mb-md">
                 <q-card-section>
                     <div class="text-h2 text-weight-bolder">8</div>
                     <div class="text-subtitle2">Ore di Orientamento Svolte</div>
                 </q-card-section>
             </q-card>

             <q-card>
                 <q-card-section>
                     <div class="text-h6">Risorse Utili</div>
                 </q-card-section>
                 <q-list separator>
                     <q-item clickable v-ripple href="https://www.universitaly.it/" target="_blank">
                         <q-item-section avatar><q-icon name="public" color="blue" /></q-item-section>
                         <q-item-section>
                             <q-item-label>Universitaly</q-item-label>
                             <q-item-label caption>Il portale del Ministero</q-item-label>
                         </q-item-section>
                         <q-item-section side><q-icon name="open_in_new" size="sm" /></q-item-section>
                     </q-item>
                      <q-item clickable v-ripple>
                         <q-item-section avatar><q-icon name="psychology" color="orange" /></q-item-section>
                         <q-item-section>Test Attitudinale</q-item-section>
                         <q-item-section side><q-icon name="chevron_right" /></q-item-section>
                     </q-item>
                 </q-list>
             </q-card>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tab = ref('upcoming')

const upcomingEvents = ref([
    {
        id: 1,
        title: 'Open Day Ingegneria',
        university: 'Politecnico',
        date: '25 Marzo, 09:00',
        description: 'Presentazione dei corsi di laurea triennale in Ingegneria.',
        image: 'https://cdn.quasar.dev/img/parallax2.jpg'
    },
    {
        id: 2,
        title: 'Medicina: Test di Ammissione',
        university: 'Statale',
        date: '10 Aprile, 14:00',
        description: 'Simulazione del test di ingresso e Q&A con studenti.',
        image: 'https://cdn.quasar.dev/img/parallax1.jpg'
    }
])

const registeredEvents = ref([
    {
        id: 3,
        title: 'Salone dello Studente',
        university: 'Fiera',
        date: '15 Febbraio, 09:00'
    }
])

const register = (event) => {
    $q.dialog({
        title: 'Conferma Iscrizione',
        message: `Vuoi iscriverti a "${event.title}"?`,
        cancel: true,
        persistent: true
    }).onOk(() => {
        $q.notify({ type: 'positive', message: 'Iscrizione effettuata con successo!' })
        // Add logic to move from upcoming to registered
    })
}
</script>
