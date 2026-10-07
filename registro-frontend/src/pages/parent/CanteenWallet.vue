<template>
  <q-page class="q-pa-md q-pa-lg-xl bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div class="row items-center q-gutter-x-sm">
        <q-avatar color="teal-1" text-color="teal-8" icon="restaurant" size="44px" />
        <div>
          <h1 class="text-h5 text-weight-bolder text-slate-900 q-my-none">
            {{ $t('parentCanteen.title') || 'Mensa Scolastica & Borsellino Elettronico' }}
          </h1>
          <div class="text-caption text-slate-500">
            {{ $t('parentCanteen.subtitle') || 'Gestione saldo pasti, ricarica online PagoPA e monitoraggio refezione' }}
          </div>
        </div>
      </div>

      <div class="row items-center q-gutter-x-sm">
        <q-chip :color="schoolSettingsStore.canteenEnabled ? 'teal-8' : 'grey-7'" text-color="white" icon="wallet">
          {{ schoolSettingsStore.canteenEnabled ? ($t('parentCanteen.active') || 'Servizio Mensa Attivo') : ($t('parentCanteen.inactive') || 'Mensa Non Attiva') }}
        </q-chip>
      </div>
    </div>

    <!-- Inactive Canteen Banner -->
    <q-card v-if="!schoolSettingsStore.canteenEnabled" flat bordered class="q-pa-xl text-center bg-white rounded-borders shadow-sm">
      <q-icon name="no_meals" size="64px" color="slate-400" />
      <div class="text-h6 text-weight-bold text-slate-700 q-mt-md">
        {{ $t('parentCanteen.inactiveTitle') || 'Il Servizio Mensa non è attivo per questo istituto' }}
      </div>
      <div class="text-caption text-slate-500 q-mt-xs max-w-md q-mx-auto">
        {{ $t('parentCanteen.inactiveDesc') || 'La scuola non prevede al momento il servizio di refezione scolastica. Qualora venga attivato dalla Dirigente o dalla DSGA, potrai gestire qui il borsellino elettronico e le ricariche.' }}
      </div>
    </q-card>

    <!-- Active Canteen Content -->
    <div v-else>
      <!-- Child Selector -->
      <q-card flat bordered class="q-pa-md bg-white rounded-borders shadow-sm q-mb-lg">
        <div class="row items-center justify-between">
          <div class="row items-center q-gutter-x-md">
            <q-icon name="face" size="28px" color="teal-8" />
            <div class="text-subtitle2 text-weight-bold text-slate-800">
              {{ $t('parentCanteen.selectChild') || 'Seleziona Figlio/a:' }}
            </div>
          </div>
          <div style="min-width: 250px">
            <q-select
              v-model="selectedChild"
              :options="childrenList"
              option-label="name"
              outlined
              dense
            />
          </div>
        </div>
      </q-card>

      <div class="row q-col-gutter-lg">
        <!-- Wallet & Topup Card -->
        <div class="col-12 col-md-6">
          <q-card flat bordered class="q-pa-lg bg-white rounded-borders shadow-sm full-height">
            <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
              👛 {{ $t('parentCanteen.walletTitle') || 'Borsellino Elettronico Pasti' }}
            </div>
            <div class="text-caption text-slate-500 q-mb-md">
              Il saldo viene scalato automaticamente ad ogni pasto registrato dal registro di classe.
            </div>

            <!-- Current Balance Display -->
            <div class="q-pa-md bg-teal-50 rounded-borders text-teal-900 q-mb-md row items-center justify-between">
              <div>
                <div class="text-caption text-weight-bold text-teal-8">SALDO DISPONIBILE</div>
                <div class="text-h4 text-weight-bolder">€ {{ selectedChild?.balance?.toFixed(2) }}</div>
              </div>
              <q-icon name="account_balance_wallet" size="48px" class="opacity-40" />
            </div>

            <!-- Quick Topup Amount Chips -->
            <div class="q-mb-sm">
              <span class="text-caption text-slate-600 q-mr-xs text-weight-bold">Ricarica Rapida PagoPA:</span>
              <q-chip clickable dense size="md" color="teal-1" text-color="teal-9" @click="topupAmount = 10.0">+10 €</q-chip>
              <q-chip clickable dense size="md" color="teal-1" text-color="teal-9" @click="topupAmount = 25.0">+25 €</q-chip>
              <q-chip clickable dense size="md" color="teal-1" text-color="teal-9" @click="topupAmount = 50.0">+50 €</q-chip>
              <q-chip clickable dense size="md" color="teal-1" text-color="teal-9" @click="topupAmount = 100.0">+100 €</q-chip>
            </div>

            <div class="q-mb-md">
              <q-input
                v-model.number="topupAmount"
                type="number"
                step="5"
                prefix="€"
                outlined
                dense
                label="Importo Ricarica"
              />
            </div>

            <q-btn
              color="teal-8"
              icon="payment"
              :label="$t('parentCanteen.topupBtn') || 'Ricarica con PagoPA (IUV)'"
              :loading="loadingTopup"
              @click="submitTopup"
              class="full-width q-py-sm text-weight-bold"
              unelevated
            />

            <!-- Last Topup Feedback -->
            <div v-if="lastNotice" class="q-mt-md q-pa-sm bg-slate-50 border border-slate-200 rounded-borders font-mono text-caption text-slate-700">
              <div><strong>Codice IUV:</strong> {{ lastNotice.iuv }}</div>
              <div><strong>Transazione:</strong> Avviso telematico generato conforme AgID</div>
            </div>
          </q-card>
        </div>

        <!-- History of Consumed Meals -->
        <div class="col-12 col-md-6">
          <q-card flat bordered class="q-pa-lg bg-white rounded-borders shadow-sm full-height">
            <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">
              📋 {{ $t('parentCanteen.historyTitle') || 'Storico Pasti & Presenze Mensa' }}
            </div>
            <div class="text-caption text-slate-500 q-mb-md">
              Ultime consumazioni rilevate per {{ selectedChild?.name }}.
            </div>

            <q-list separator class="rounded-borders border border-slate-100">
              <q-item v-for="(m, idx) in mealHistory" :key="idx" class="q-py-sm">
                <q-item-section avatar>
                  <q-avatar color="teal-1" text-color="teal-9" icon="check" size="32px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-bold">{{ m.date }}</q-item-label>
                  <q-item-label caption class="text-slate-500">{{ m.desc }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <span class="text-weight-bolder text-negative">- € {{ m.cost.toFixed(2) }}</span>
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
import { useQuasar } from 'quasar'
import { useSchoolSettingsStore } from '@/stores/schoolSettings'
import { mealsService } from '@/services/mealsService'

const $q = useQuasar()
const schoolSettingsStore = useSchoolSettingsStore()

const loadingTopup = ref(false)
const topupAmount = ref(50.0)
const lastNotice = ref(null)

const childrenList = ref([
  { id: 'student-minor-1', name: 'Alessandro Verdi (2A)', balance: 14.50 },
  { id: 'student-1', name: 'Marco Rossi (5A)', balance: 34.50 }
])
const selectedChild = ref(childrenList.value[0])

const mealHistory = ref([
  { date: '06/10/2026', desc: 'Pasto Completo - Refettorio Centrale', cost: 4.50 },
  { date: '05/10/2026', desc: 'Pasto Completo - Refettorio Centrale', cost: 4.50 },
  { date: '02/10/2026', desc: 'Pasto Completo - Refettorio Centrale', cost: 4.50 },
  { date: '01/10/2026', desc: 'Pasto Completo - Refettorio Centrale', cost: 4.50 }
])

async function submitTopup() {
  if (!selectedChild.value || topupAmount.value <= 0) return
  loadingTopup.value = true
  try {
    const fakeIUV = `0121${Math.floor(1000000000000 + Math.random() * 9000000000000)}`
    await mealsService.topUpWallet(selectedChild.value.id, topupAmount.value, fakeIUV)
    selectedChild.value.balance += topupAmount.value
    lastNotice.value = { iuv: fakeIUV }
    $q.notify({
      type: 'positive',
      message: `Ricarica di € ${topupAmount.value.toFixed(2)} effettuata con successo! IUV: ${fakeIUV}`
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message || 'Errore nella ricarica del borsellino.'
    })
  } finally {
    loadingTopup.value = false
  }
}
</script>
