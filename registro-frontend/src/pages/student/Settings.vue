<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          Impostazioni Studente
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          Personalizza l'aspetto visivo, i font di lettura, il tema ed i dettagli del tuo profilo.
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
                <q-item-label>Font &amp; Aspetto</q-item-label>
                <q-item-label caption>Scelta dei font e temi</q-item-label>
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
                <q-item-label>Accessibilità DSA</q-item-label>
                <q-item-label caption>Contrasto e lettura</q-item-label>
              </q-item-section>
            </q-item>

            <q-item
              clickable
              :active="activeSection === 'notifications'"
              active-class="bg-indigo-50 text-indigo-700 text-weight-bold rounded-lg"
              class="rounded-lg"
              @click="activeSection = 'notifications'"
            >
              <q-item-section avatar>
                <q-icon name="notifications" color="amber-9" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Notifiche</q-item-label>
                <q-item-label caption>Voti, compiti e comunicazioni</q-item-label>
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
                <q-item-label caption>Gestione credenziali</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Right Column: Active Section Content -->
      <div class="col-12 col-md-9">
        <!-- SECTION 1: FONT & APPEARANCE -->
        <q-card v-if="activeSection === 'appearance'" flat bordered class="rounded-xl bg-white q-pa-lg shadow-sm">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Scelta dei Font &amp; Tipografia</div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Scegli il carattere tipografico ideale per facilitare lo studio e la consultazione del registro.
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
                  Anteprima: 1 2 3 4 5 6 7 8 9 0 - ABC abc
                </div>
              </q-card>
            </div>
          </div>

          <q-separator class="q-my-lg" />

          <!-- Font Size Selection -->
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Dimensione del Testo</div>
          <div class="text-caption text-slate-500 q-mb-md">
            Regola l'ingrandimento dei caratteri in tutta l'applicazione.
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
            Seleziona la palette colori preferita.
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

        <!-- SECTION 2: ACCESSIBILITY -->
        <q-card v-if="activeSection === 'accessibility'" flat bordered class="rounded-xl bg-white q-pa-lg shadow-sm">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Accessibilità Visiva (A11y / BES)</div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Impostazioni avanzate di contrasto ed ausili alla lettura per studenti DSA/BES.
          </div>

          <div class="q-gutter-y-md">
            <q-card flat bordered class="q-pa-md rounded-xl bg-slate-50">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-subtitle1 text-weight-bold text-slate-800">Font OpenDyslexic (Alta Leggibilità DSA)</div>
                  <div class="text-caption text-slate-500">Font con peso alla base della lettera per prevenire l'inversione di caratteri</div>
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
                  <div class="text-caption text-slate-500">Migliora la nitidezza ed aggiunge bordi definiti agli elementi</div>
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

        <!-- SECTION 3: NOTIFICATIONS -->
        <q-card v-if="activeSection === 'notifications'" flat bordered class="rounded-xl bg-white q-pa-lg shadow-sm">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Preferenze Notifiche</div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Scegli quali avvisi ed aggiornamenti ricevere nel tuo pannello.
          </div>

          <q-list divider class="rounded-xl">
            <q-item tag="label" v-ripple>
              <q-item-section avatar>
                <q-icon name="grade" color="primary" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Nuovi Voti Registrati</q-item-label>
                <q-item-label caption>Ricevi una notifica immediata quando viene pubblicato un nuovo voto</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notifGrades" color="primary" />
              </q-item-section>
            </q-item>

            <q-item tag="label" v-ripple>
              <q-item-section avatar>
                <q-icon name="assignment" color="amber-9" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Compiti ed Attività Assegnate</q-item-label>
                <q-item-label caption>Notifica quando i docenti inseriscono compiti in agenda</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notifHomework" color="amber-9" />
              </q-item-section>
            </q-item>

            <q-item tag="label" v-ripple>
              <q-item-section avatar>
                <q-icon name="email" color="teal" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Circolari e Comunicazioni Scuola</q-item-label>
                <q-item-label caption>Avvisi ed annunci dalla presidenza o segreteria</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="notifComms" color="teal" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>

        <!-- SECTION 4: SECURITY -->
        <q-card v-if="activeSection === 'security'" flat bordered class="rounded-xl bg-white q-pa-lg shadow-sm">
          <div class="text-h6 text-weight-bold text-slate-800 q-mb-xs">Sicurezza &amp; Password</div>
          <div class="text-caption text-slate-500 q-mb-lg">
            Aggiorna la tua password di accesso al registro elettronico.
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

const notifGrades = ref(true)
const notifHomework = ref(true)
const notifComms = ref(true)

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
