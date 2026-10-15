<template>
  <q-page class="q-pa-md q-pa-lg-xl bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div class="row items-center q-gutter-x-sm">
        <q-avatar color="teal-1" text-color="teal-8" icon="restaurant" size="44px" />
        <div>
          <h1 class="text-h5 text-weight-bolder text-slate-900 q-my-none">
            {{ $t('studentCanteen.title') || 'Mensa Scolastica & Refezione' }}
          </h1>
          <div class="text-caption text-slate-500">
            {{ $t('studentCanteen.subtitle') || 'Visualizzazione presenze mensa, menu del giorno e diete speciali' }}
          </div>
        </div>
      </div>

      <div class="row items-center q-gutter-x-sm">
        <q-chip :color="schoolSettingsStore.canteenEnabled ? 'teal-8' : 'grey-7'" text-color="white" icon="dining">
          {{ schoolSettingsStore.canteenEnabled ? ($t('studentCanteen.serviceActive') || 'Servizio Mensa Attivo') : ($t('studentCanteen.serviceInactive') || 'Mensa Non Attiva') }}
        </q-chip>
      </div>
    </div>

    <!-- Inactive Canteen Banner -->
    <q-card v-if="!schoolSettingsStore.canteenEnabled" flat bordered class="q-pa-xl text-center bg-white rounded-borders shadow-sm">
      <q-icon name="no_meals" size="64px" color="slate-400" />
      <div class="text-h6 text-weight-bold text-slate-700 q-mt-md">
        {{ $t('studentCanteen.inactiveTitle') || 'Il Servizio Mensa non è attivo per questo istituto' }}
      </div>
      <div class="text-caption text-slate-500 q-mt-xs max-w-md q-mx-auto">
        {{ $t('studentCanteen.inactiveDesc') || 'La refezione scolastica non è prevista per questo plesso o tipologia di scuola. Qualora il servizio venga attivato dalla Dirigenza o dalla DSGA, le informazioni compariranno automaticamente qui.' }}
      </div>
    </q-card>

    <!-- Active Canteen Content -->
    <div v-else>
      <!-- Today's Meal Status Card -->
      <div class="row q-col-gutter-lg q-mb-lg">
        <div class="col-12 col-md-4">
          <q-card flat bordered class="q-pa-md bg-white rounded-borders shadow-sm full-height">
            <div class="text-caption text-weight-bold text-slate-500 q-mb-xs">PASTO DI OGGI</div>
            <div class="row items-center justify-between">
              <span class="text-h6 text-weight-bolder text-teal-9">Presente a Mensa</span>
              <q-badge color="positive" label="Confermato" />
            </div>
            <div class="text-caption text-slate-600 q-mt-sm">
              Rilevazione automatica del mattino confermata dal docente.
            </div>
          </q-card>
        </div>

        <div class="col-12 col-md-4">
          <q-card flat bordered class="q-pa-md bg-white rounded-borders shadow-sm full-height">
            <div class="text-caption text-weight-bold text-slate-500 q-mb-xs">PROFILO DIETA REGISTRATO</div>
            <div class="row items-center justify-between">
              <span class="text-h6 text-weight-bolder text-indigo-9">{{ studentDiet.name }}</span>
              <q-badge :color="studentDiet.isSpecial ? 'purple' : 'teal'" :label="studentDiet.isSpecial ? 'Speciale' : 'Standard'" />
            </div>
            <div class="text-caption text-slate-600 q-mt-sm">
              {{ studentDiet.notes }}
            </div>
          </q-card>
        </div>

        <div class="col-12 col-md-4">
          <q-card flat bordered class="q-pa-md bg-white rounded-borders shadow-sm full-height">
            <div class="text-caption text-weight-bold text-slate-500 q-mb-xs">TURNO REFEZIONE</div>
            <div class="row items-center justify-between">
              <span class="text-h6 text-weight-bolder text-amber-9">1° Turno (13:10 - 13:50)</span>
              <q-badge color="amber-9" label="Refettorio Centrale" />
            </div>
            <div class="text-caption text-slate-600 q-mt-sm">
              Accompagnamento con il docente dell'ultima ora.
            </div>
          </q-card>
        </div>
      </div>

      <!-- Today's Menu -->
      <div class="row q-col-gutter-lg">
        <div class="col-12 col-md-6">
          <q-card flat bordered class="q-pa-lg bg-white rounded-borders shadow-sm">
            <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
              🍽️ {{ $t('studentCanteen.todayMenuTitle') || 'Menu del Giorno' }}
            </div>
            <div class="text-caption text-slate-500 q-mb-md">
              Preparato secondo le Linee Guida Nazionali per la Ristorazione Scolastica.
            </div>

            <q-list separator class="rounded-borders border border-slate-100">
              <q-item class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="amber-1" text-color="amber-9" icon="lunch_dining" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Primo Piatto</q-item-label>
                  <q-item-label caption class="text-slate-700">Risotto biologico alla zucca con parmigiano reggiano DOP</q-item-label>
                </q-item-section>
              </q-item>

              <q-item class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="red-1" text-color="red-9" icon="set_meal" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Secondo Piatto</q-item-label>
                  <q-item-label caption class="text-slate-700">Filetto di orata al forno con erbe aromatiche</q-item-label>
                </q-item-section>
              </q-item>

              <q-item class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="green-1" text-color="green-9" icon="eco" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Contorno</q-item-label>
                  <q-item-label caption class="text-slate-700">Carote julienne e fagiolini al vapore con olio EVO</q-item-label>
                </q-item-section>
              </q-item>

              <q-item class="q-py-md">
                <q-item-section avatar>
                  <q-avatar color="orange-1" text-color="orange-9" icon="nutrition" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">Frutta / Dessert</q-item-label>
                  <q-item-label caption class="text-slate-700">Mela biologica a km 0 e pane integrale</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>

        <!-- Weekly Plan -->
        <div class="col-12 col-md-6">
          <q-card flat bordered class="q-pa-lg bg-white rounded-borders shadow-sm">
            <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
              📅 {{ $t('studentCanteen.weeklyCalendarTitle') || 'Calendario Settimanale Refezione' }}
            </div>
            <div class="text-caption text-slate-500 q-mb-md">
              Piano nutrizionale approvato dal dietista ASL.
            </div>

            <q-list separator class="rounded-borders border border-slate-100">
              <q-item class="q-py-sm">
                <q-item-section>
                  <q-item-label class="text-weight-bold">Lunedì</q-item-label>
                  <q-item-label caption>Pasta al pomodoro fresco • Polpette di tacchino al forno • Insalata mista</q-item-label>
                </q-item-section>
              </q-item>
              <q-item class="q-py-sm bg-teal-50">
                <q-item-section>
                  <q-item-label class="text-weight-bold text-teal-9">Martedì (Oggi)</q-item-label>
                  <q-item-label caption class="text-teal-8">Risotto alla zucca • Filetto di orata • Carote e fagiolini</q-item-label>
                </q-item-section>
              </q-item>
              <q-item class="q-py-sm">
                <q-item-section>
                  <q-item-label class="text-weight-bold">Mercoledì</q-item-label>
                  <q-item-label caption>Passato di legumi con crostini • Frittatina alle verdure • Finocchi all'olio</q-item-label>
                </q-item-section>
              </q-item>
              <q-item class="q-py-sm">
                <q-item-section>
                  <q-item-label class="text-weight-bold">Giovedì</q-item-label>
                  <q-item-label caption>Gnocchetti al pesto leggero • Bocconcini di pollo al limone • Patate al forno</q-item-label>
                </q-item-section>
              </q-item>
              <q-item class="q-py-sm">
                <q-item-section>
                  <q-item-label class="text-weight-bold">Venerdì</q-item-label>
                  <q-item-label caption>Lasagnetta vegetariana • Formaggio fresco con verdure crude • Frutta di stagione</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useSchoolSettingsStore } from '@/stores/schoolSettings'

const schoolSettingsStore = useSchoolSettingsStore()

const studentDiet = ref({
  name: 'Dieta Mediterranea Standard',
  isSpecial: false,
  notes: 'Nessuna intolleranza o restrizione etico-religiosa segnalata.'
})
</script>
