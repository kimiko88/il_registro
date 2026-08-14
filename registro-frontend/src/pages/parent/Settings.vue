<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          Impostazioni Genitore
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Personalizza font, preferenze visive, notifiche per i tuoi figli e sicurezza dell'account.
        </p>
      </div>
    </div>

    <div class="row q-col-gutter-lg">
      <!-- Left Column: Navigation Tabs / Sections -->
      <div class="col-12 col-md-3">
        <q-card flat bordered class="rounded-xl bg-white q-pa-sm shadow-sm">
          <q-list class="q-gutter-y-xs">
            <q-item
              clickable
              :active="activeSection === 'appearance'"
              active-class="bg-indigo-50 text-indigo-700 text-weight-bold rounded-lg"
              class="rounded-lg"
              @click="activeSection = 'appearance'"
            >
              <q-item-section avatar>
                <q-icon name="font_download" color="primary" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Font &amp; Leggibilità</q-item-label>
                <q-item-label caption>Scelta dei font e temi</q-item-label>
              </q-item-section>
            </q-item>

            <q-item
              clickable
              :active="activeSection === 'children'"
              active-class="bg-indigo-50 text-indigo-700 text-weight-bold rounded-lg"
              class="rounded-lg"
              @click="activeSection = 'children'"
            >
              <q-item-section avatar>
                <q-icon name="family_restroom" color="purple" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Notifiche Figli</q-item-label>
                <q-item-label caption>Avvisi presenze e voti</q-item-label>
              </q-item-section>
            </q-item>

            <q-item
              clickable
              :active="activeSection === 'accessibility'"
              active-class="bg-indigo-50 text-indigo-700 text-weight-bold rounded-lg"
              class="rounded-lg"
              @click="activeSection = 'accessibility'"
            >
              <q-item-section avatar>
                <q-icon name="accessibility_new" color="secondary" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Accessibilità</q-item-label>
                <q-item-label caption>OpenDyslexic e contrasto</q-item-label>
              </q-item-section>
            </q-item>

            <q-item
              clickable
              :active="activeSection === 'security'"
              active-class="bg-indigo-50 text-indigo-700 text-weight-bold rounded-lg"
              class="rounded-lg"
              @click="activeSection = 'security'"
            >
              <q-item-section avatar>
                <q-icon name="security" color="teal" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Sicurezza &amp; Password</q-item-label>
                <q-item-label caption>Credenziali ed il tuo account</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Right Column: Active Section Content -->
      <div class="col-12 col-md-9">
        <!-- SECTION 1: FONT & APPEARANCE -->
        <q-card v-if="activeSection === 'appearance'" flat bordered class="rounded-xl bg-white q-pa-lg shadow-sm">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Scelta dei Font &amp; Leggibilità</div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Scegli il font tipografico preferito per la lettura di circolari, valutazioni e documenti scolastici.
          </div>

          <div class="row q-col-gutter-md q-mb-xl">
            <div
              v-for="font in fontOptions"
              :key="font.value"
              class="col-12 col-sm-6 col-md-4"
            >
              <q-card
                flat
                bordered
                clickable
                class="q-pa-md rounded-xl transition-all h-full relative-position"
                :class="themeStore.fontFamily === font.value ? 'border-2 border-indigo-600 bg-indigo-50/40 shadow-md' : 'hover:border-indigo-200 bg-white'"
                @click="themeStore.setFontFamily(font.value)"
              >
                <div class="row items-center justify-between q-mb-sm">
                  <span class="text-subtitle1 text-weight-bold text-slate-800">{{ font.label }}</span>
                  <q-icon v-if="themeStore.fontFamily === font.value" name="check_circle" color="indigo" size="22px" />
                </div>
                <div class="text-caption text-slate-500 q-mb-md">{{ font.description }}</div>
                <div
                  class="q-pa-sm bg-slate-100/80 rounded-lg text-slate-700 text-caption font-bold"
                  :style="{ fontFamily: font.previewFont }"
                >
                  Anteprima: Valutazioni e comunicazione scuola-famiglia.
                </div>
              </q-card>
            </div>
          </div>

          <q-separator class="q-my-lg" />

          <!-- Font Size Selection -->
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Dimensione del Testo</div>
          <div class="text-caption text-slate-500 q-mb-md">
            Regola l'ingrandimento dei caratteri in tutto il registro.
          </div>

          <div class="row q-col-gutter-md q-mb-xl">
            <div
              v-for="size in sizeOptions"
              :key="size.value"
              class="col-12 col-sm-4"
            >
              <q-card
                flat
                bordered
                clickable
                class="q-pa-md rounded-xl text-center transition-all"
                :class="themeStore.fontSize === size.value ? 'border-2 border-indigo-600 bg-indigo-50/40 shadow-sm' : 'bg-white'"
                @click="themeStore.setFontSize(size.value)"
              >
                <q-icon name="text_fields" :size="size.iconSize" color="indigo" class="q-mb-xs" />
                <div class="text-subtitle2 text-weight-bold text-slate-800">{{ size.label }}</div>
                <div class="text-caption text-slate-500">{{ size.scale }}</div>
              </q-card>
            </div>
          </div>

          <q-separator class="q-my-lg" />

          <!-- Theme Palette Selection -->
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Temi Visivi dell'Interfaccia</div>
          <div class="text-caption text-slate-500 q-mb-md">
            Seleziona lo stile visivo ed i colori dell'applicazione.
          </div>

          <div class="row q-col-gutter-md">
            <div
              v-for="t in THEMES"
              :key="t.id"
              class="col-12 col-sm-6 col-md-4"
            >
              <q-card
                flat
                bordered
                clickable
                class="q-pa-md rounded-xl transition-all"
                :class="themeStore.currentTheme === t.id ? 'border-2 border-indigo-600 bg-indigo-50/30 shadow-sm' : 'bg-white'"
                @click="themeStore.setTheme(t.id)"
              >
                <div class="row items-center justify-between q-mb-xs">
                  <div class="row items-center">
                    <q-avatar size="28px" :color="t.badgeColor" text-color="white" class="q-mr-xs">
                      <q-icon :name="t.icon" size="16px" />
                    </q-avatar>
                    <span class="text-weight-bold text-slate-800">{{ t.name }}</span>
                  </div>
                  <q-icon v-if="themeStore.currentTheme === t.id" name="check" color="indigo" size="20px" />
                </div>
                <div class="text-caption text-slate-500">{{ t.description }}</div>
              </q-card>
            </div>
          </div>
        </q-card>

        <!-- SECTION 2: NOTIFICATIONS FOR CHILDREN -->
        <q-card v-if="activeSection === 'children'" flat bordered class="rounded-xl bg-white q-pa-lg shadow-sm">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Notifiche &amp; Avvisi Figli</div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Imposta gli avvisi in tempo reale per presenze, voti, note e circolari.
          </div>

          <q-list divider class="rounded-xl">
            <q-item tag="label" v-ripple>
              <q-item-section avatar>
                <q-icon name="event_busy" color="negative" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Assenze, Ritardi ed Uscite</q-item-label>
                <q-item-label caption>Avviso immediato in caso di assenza non giustificata o ingresso in ritardo</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notifAttendance" color="negative" />
              </q-item-section>
            </q-item>

            <q-item tag="label" v-ripple>
              <q-item-section avatar>
                <q-icon name="grade" color="primary" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Pubblicazione Nuove Valutazioni</q-item-label>
                <q-item-label caption>Notifica quando i docenti inseriscono un voto o un giudizio</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notifGrades" color="primary" />
              </q-item-section>
            </q-item>

            <q-item tag="label" v-ripple>
              <q-item-section avatar>
                <q-icon name="campaign" color="amber-9" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Circolari da Firmare</q-item-label>
                <q-item-label caption>Promemoria per comunicati e circolari che richiedono presa visione</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notifCirculars" color="amber-9" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>

        <!-- SECTION 3: ACCESSIBILITY -->
        <q-card v-if="activeSection === 'accessibility'" flat bordered class="rounded-xl bg-white q-pa-lg shadow-sm">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Accessibilità Visiva (A11y)</div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Opzioni di visualizzazione ad alta leggibilità.
          </div>

          <div class="q-gutter-y-md">
            <q-card flat bordered class="q-pa-md rounded-xl bg-slate-50">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-subtitle1 text-weight-bold text-slate-800">Font OpenDyslexic (Alta Leggibilità DSA)</div>
                  <div class="text-caption text-slate-500">Applica il font facilitato a tutti i testi dell'applicazione</div>
                </div>
                <q-toggle
                  v-model="themeStore.dsaFont"
                  color="indigo"
                  size="lg"
                  @update:model-value="themeStore.toggleDsaFont"
                />
              </div>
            </q-card>

            <q-card flat bordered class="q-pa-md rounded-xl bg-slate-50">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-subtitle1 text-weight-bold text-slate-800">Modalità Contrasto Elevato</div>
                  <div class="text-caption text-slate-500">Aumenta la definizione dei bordi e la nitidezza del testo</div>
                </div>
                <q-toggle
                  v-model="themeStore.highContrast"
                  color="indigo"
                  size="lg"
                  @update:model-value="themeStore.toggleHighContrast"
                />
              </div>
            </q-card>
          </div>
        </q-card>

        <!-- SECTION 4: SECURITY -->
        <q-card v-if="activeSection === 'security'" flat bordered class="rounded-xl bg-white q-pa-lg shadow-sm">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Sicurezza &amp; Password</div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Modifica la password del tuo account genitore.
          </div>

          <q-form @submit.prevent="changePassword" class="q-gutter-y-md style-form">
            <q-input
              v-model="pwdCurrent"
              type="password"
              label="Password Attuale"
              outlined
              dense
              class="rounded-lg"
            />
            <q-input
              v-model="pwdNew"
              type="password"
              label="Nuova Password (almeno 8 caratteri)"
              outlined
              dense
              class="rounded-lg"
            />
            <q-input
              v-model="pwdConfirm"
              type="password"
              label="Conferma Nuova Password"
              outlined
              dense
              class="rounded-lg"
            />

            <div class="row justify-end q-mt-lg">
              <q-btn
                type="submit"
                color="primary"
                label="Aggiorna Password"
                unelevated
                class="rounded-lg font-bold"
                :loading="savingPwd"
              />
            </div>
          </q-form>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useThemeStore, THEMES } from 'src/stores/theme'
import { useQuasar } from 'quasar'
import { userService } from 'src/services/userService'
import { useAuthStore } from 'src/stores/auth'

import { useI18n } from 'vue-i18n'

const $q = useQuasar()
const { t, te } = useI18n()
const themeStore = useThemeStore()
const authStore = useAuthStore()

const activeSection = ref('appearance')

const fontOptions = [
  { label: 'Standard (Inter)', value: 'default', description: 'Carattere predefinito pulito e bilanciato', previewFont: 'Inter, sans-serif' },
  { label: 'OpenDyslexic (DSA)', value: 'opendyslexic', description: 'Specifico per dislessia ed alta leggibilità BES', previewFont: 'OpenDyslexic, sans-serif' },
  { label: 'Lexend (Facilitato)', value: 'lexend', description: 'Ottimizzato per la velocità di comprensione', previewFont: 'Lexend, sans-serif' },
  { label: 'Fredoka (Friendly)', value: 'fredoka', description: 'Tratti morbidi ed arrotondati, ideale per ragazzi', previewFont: 'Fredoka, sans-serif' },
  { label: 'Roboto', value: 'roboto', description: 'Stile geometrico ad alta chiarezza', previewFont: 'Roboto, sans-serif' }
]

const sizeOptions = [
  { label: 'Normale', value: 'normal', scale: '100% (Standard)', iconSize: '24px' },
  { label: 'Grande', value: 'large', scale: '110% (Ingrandito)', iconSize: '30px' },
  { label: 'Molto Grande', value: 'xlarge', scale: '120% (BES/Accessibile)', iconSize: '36px' }
]

const notifAttendance = ref(true)
const notifGrades = ref(true)
const notifCirculars = ref(true)

const pwdCurrent = ref('')
const pwdNew = ref('')
const pwdConfirm = ref('')
const savingPwd = ref(false)

async function changePassword() {
  if (!pwdCurrent.value || !pwdNew.value) {
    $q.notify({ type: 'warning', message: t('errors.ERR_REQUIRED_FIELDS') })
    return
  }
  if (pwdNew.value !== pwdConfirm.value) {
    $q.notify({ type: 'negative', message: t('errors.passwordMismatch') })
    return
  }
  const userId = authStore.user?.id
  if (!userId) {
    $q.notify({ type: 'negative', message: t('errors.sessionInvalid') })
    return
  }
  savingPwd.value = true
  try {
    await userService.changePassword(userId, pwdCurrent.value, pwdNew.value)
    pwdCurrent.value = ''
    pwdNew.value = ''
    pwdConfirm.value = ''
    $q.notify({ type: 'positive', message: t('notifications.passwordUpdated') })
  } catch (err) {
    const errData = err.response?.data
    const code = errData?.code || errData?.error
    const msg = code && te(`errors.${code}`) ? t(`errors.${code}`) : (errData?.message || errData?.error || t('errors.serverError'))
    $q.notify({ type: 'negative', message: msg })
  } finally {
    savingPwd.value = false
  }
}
</script>
