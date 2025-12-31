<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-md">Orientamento in Uscita</div>

    <div class="row q-col-gutter-lg">
        <!-- Upcoming Events -->
        <div class="col-12 col-md-8">
            <q-card>
                <q-toolbar class="bg-primary text-white">
                    <q-toolbar-title>Eventi in Arrivo</q-toolbar-title>
                    <q-btn flat round icon="filter_list" />
                </q-toolbar>

                <q-list separator>
                    <q-item v-for="event in events" :key="event.id" class="q-py-md">
                        <q-item-section avatar>
                            <q-avatar rounded color="blue-1" text-color="blue" icon="school" />
                        </q-item-section>
                        
                        <q-item-section>
                            <div class="text-h6">{{ event.title }}</div>
                            <div class="text-subtitle2 text-grey-8">{{ event.organizer }}</div>
                            <div class="row items-center text-caption text-grey q-gutter-md q-mt-xs">
                                <span><q-icon name="event" /> {{ event.date }}</span>
                                <span><q-icon name="place" /> {{ event.location }}</span>
                            </div>
                        </q-item-section>
                        
                        <q-item-section side>
                            <q-btn 
                                :label="event.registered ? 'Iscritto' : 'Iscriviti'" 
                                :color="event.registered ? 'green' : 'primary'"
                                :outline="!event.registered"
                                :icon="event.registered ? 'check' : 'add'"
                                @click="toggleRegistration(event)"
                            />
                        </q-item-section>
                    </q-item>
                </q-list>
            </q-card>
        </div>

        <!-- My Schedule / Resources -->
        <div class="col-12 col-md-4">
            <q-card class="q-mb-md">
                <q-card-section>
                    <div class="text-h6">I Miei Appuntamenti</div>
                    <q-list dense class="q-mt-sm">
                        <q-item v-if="myEvents.length === 0">
                            <q-item-section class="text-grey text-italic">Nessun evento prenotato</q-item-section>
                        </q-item>
                        <q-item v-for="ev in myEvents" :key="ev.id">
                             <q-item-section>
                                 <q-item-label>{{ ev.title }}</q-item-label>
                                 <q-item-label caption>{{ ev.date }}</q-item-label>
                             </q-item-section>
                        </q-item>
                    </q-list>
                </q-card-section>
            </q-card>

            <q-card class="bg-teal-1">
                <q-card-section>
                    <div class="text-subtitle1 text-weight-bold">Risorse Utili</div>
                    <q-list class="q-mt-sm">
                         <q-item clickable tag="a" href="#">
                             <q-item-section avatar><q-icon name="language" /></q-item-section>
                             <q-item-section>Portale Universitaly</q-item-section>
                         </q-item>
                         <q-item clickable tag="a" href="#">
                             <q-item-section avatar><q-icon name="work" /></q-item-section>
                             <q-item-section>Guida ITS 2025</q-item-section>
                         </q-item>
                    </q-list>
                </q-card-section>
            </q-card>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

const events = ref([
    { id: 1, title: 'Open Day Politecnico', organizer: 'Politecnico di Milano', date: '20 Feb 2025', location: 'Milano / Online', registered: false },
    { id: 2, title: 'Salone dello Studente', organizer: 'Campus Editori', date: '15 Mar 2025', location: 'Fiera Roma', registered: true }
])

const myEvents = computed(() => events.value.filter(e => e.registered))

const toggleRegistration = (event) => {
    event.registered = !event.registered
    const msg = event.registered ? 'Iscrizione confermata' : 'Iscrizione annullata'
    $q.notify({ type: 'info', message: msg })
}
</script>
