<template>
  <div class="accessibility-settings-panel">
    <!-- Header Summary & Reset -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <div class="text-h6 text-weight-bold text-slate-800 flex items-center gap-sm">
          <q-avatar size="36px" color="primary" text-color="white" icon="accessibility_new" />
          {{ $t('a11y.panelTitle') || 'Impostazioni di Accessibilità & Inclusione Visiva' }}
        </div>
        <div class="text-caption text-slate-500 q-mt-xs">
          {{ $t('a11y.panelSubtitle') || 'Personalizza l\'esperienza visiva, i font, la sintesi vocale e gli strumenti di supporto DSA in conformità alle Linee Guida AgID / WCAG 2.2.' }}
        </div>
      </div>
      <div class="flex gap-sm items-center">
        <q-btn
          outline
          color="primary"
          icon="keyboard"
          :label="$t('a11y.shortcutsBtn') || 'Scorciatoie Tastiera (?)'"
          no-caps
          rounded
          dense
          class="q-px-sm"
          @click="themeStore.toggleKeyboardShortcutsHelp(true)"
        />
        <q-btn
          flat
          color="negative"
          icon="restore"
          :label="$t('a11y.resetDefaults') || 'Ripristina Predefiniti'"
          no-caps
          rounded
          dense
          class="q-px-sm"
          @click="themeStore.resetAccessibility"
        >
          <q-tooltip>{{ $t('a11y.resetDefaultsTooltip') || 'Reimposta tutte le opzioni di accessibilità ai valori standard' }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <div class="row q-col-gutter-lg">
      <!-- 1. Font & Lettura DSA -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="rounded-xl q-pa-md bg-white full-height shadow-sm">
          <div class="row items-center justify-between q-mb-sm">
            <div class="flex items-center gap-sm">
              <q-avatar size="40px" color="indigo-50" text-color="indigo-700" icon="spellcheck" />
              <div>
                <div class="text-subtitle1 text-weight-bold text-slate-800">Font OpenDyslexic (DSA)</div>
                <div class="text-caption text-slate-500">Font specifico per la dislessia e difficoltà di decodifica</div>
              </div>
            </div>
            <q-toggle
              v-model="themeStore.dsaFont"
              color="indigo"
              size="lg"
              @update:model-value="themeStore.toggleDsaFont"
            />
          </div>

          <div class="q-mt-sm">
            <div class="text-caption text-weight-bold text-slate-700 q-mb-xs">Font Alternativi ad Alta Leggibilità</div>
            <q-select
              v-model="themeStore.fontFamily"
              :options="fontFamilyOptions"
              emit-value
              map-options
              outlined
              dense
              @update:model-value="themeStore.setFontFamily"
            />
          </div>

          <div class="q-mt-sm">
            <div class="text-caption text-weight-bold text-slate-700 q-mb-xs">Dimensione Caratteri</div>
            <q-btn-toggle
              v-model="themeStore.fontSize"
              :options="fontSizeOptions"
              spread
              no-caps
              rounded
              unelevated
              toggle-color="indigo"
              color="grey-2"
              text-color="grey-9"
              @update:model-value="themeStore.setFontSize"
            />
          </div>

          <q-separator class="q-my-md" />
          <div class="q-pa-sm bg-slate-50 rounded-lg text-slate-700 text-caption border border-slate-100" :class="{ 'dsa-font-active': themeStore.dsaFont }">
            <strong>Anteprima:</strong> Il registro elettronico supporta l'inclusione scolastica e l'apprendimento personalizzato.
          </div>
        </q-card>
      </div>

      <!-- 2. Righello di Lettura & Focus Mask -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="rounded-xl q-pa-md bg-white full-height shadow-sm">
          <div class="row items-center justify-between q-mb-sm">
            <div class="flex items-center gap-sm">
              <q-avatar size="40px" color="amber-50" text-color="amber-9" icon="straighten" />
              <div>
                <div class="text-subtitle1 text-weight-bold text-slate-800">Righello di Lettura (Focus Mask)</div>
                <div class="text-caption text-slate-500">Guida orizzontale che segue il cursore per non perdere il rigo</div>
              </div>
            </div>
            <q-toggle
              v-model="themeStore.readingRuler"
              color="amber-9"
              size="lg"
              @update:model-value="themeStore.toggleReadingRuler"
            />
          </div>

          <div class="q-mt-sm">
            <div class="row justify-between items-center text-caption text-weight-bold text-slate-700">
              <span>Altezza Finestra di Lettura</span>
              <span class="text-amber-9">{{ themeStore.readingRulerHeight }} px</span>
            </div>
            <q-slider
              v-model="themeStore.readingRulerHeight"
              :min="24"
              :max="70"
              :step="2"
              color="amber-9"
              @update:model-value="themeStore.setReadingRulerHeight"
            />
          </div>

          <div class="q-mt-xs">
            <div class="row justify-between items-center text-caption text-weight-bold text-slate-700">
              <span>Intensità Maschera Oscurante</span>
              <span class="text-amber-9">{{ Math.round(themeStore.readingRulerOpacity * 100) }}%</span>
            </div>
            <q-slider
              v-model="themeStore.readingRulerOpacity"
              :min="0.1"
              :max="0.8"
              :step="0.05"
              color="amber-9"
              @update:model-value="themeStore.setReadingRulerOpacity"
            />
          </div>

          <q-separator class="q-my-sm" />
          <div class="text-caption text-grey-7">
            <q-icon name="lightbulb" color="amber-9" size="xs" class="q-mr-xs" />
            Suggerimento: Premi <code>Alt + R</code> per accendere/spegnere il righello, oppure usa <code>Alt + ↑ / ↓</code>.
          </div>
        </q-card>
      </div>

      <!-- 3. Spaziatura Testo (WCAG 1.4.12) -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="rounded-xl q-pa-md bg-white full-height shadow-sm">
          <div class="flex items-center gap-sm q-mb-sm">
            <q-avatar size="40px" color="blue-50" text-color="blue-700" icon="format_line_spacing" />
            <div>
              <div class="text-subtitle1 text-weight-bold text-slate-800">Spaziatura Testo & Interlinea</div>
              <div class="text-caption text-slate-500">Regolazione fine per il criterio di conformità WCAG 1.4.12</div>
            </div>
          </div>

          <div class="row q-col-gutter-sm q-mt-xs">
            <div class="col-12 col-sm-4">
              <div class="text-caption text-weight-bold text-slate-700 q-mb-xs">Interlinea</div>
              <q-select
                v-model="themeStore.lineHeight"
                :options="lineHeightOptions"
                emit-value
                map-options
                outlined
                dense
                @update:model-value="themeStore.setLineHeight"
              />
            </div>
            <div class="col-12 col-sm-4">
              <div class="text-caption text-weight-bold text-slate-700 q-mb-xs">Spazio Lettere</div>
              <q-select
                v-model="themeStore.letterSpacing"
                :options="letterSpacingOptions"
                emit-value
                map-options
                outlined
                dense
                @update:model-value="themeStore.setLetterSpacing"
              />
            </div>
            <div class="col-12 col-sm-4">
              <div class="text-caption text-weight-bold text-slate-700 q-mb-xs">Spazio Parole</div>
              <q-select
                v-model="themeStore.wordSpacing"
                :options="wordSpacingOptions"
                emit-value
                map-options
                outlined
                dense
                @update:model-value="themeStore.setWordSpacing"
              />
            </div>
          </div>

          <q-separator class="q-my-md" />
          <div class="row items-center justify-between">
            <div class="text-caption text-weight-medium text-slate-700">Evidenziazione Focus Tastiera (Focus Ring)</div>
            <q-toggle
              v-model="themeStore.focusHighlight"
              color="primary"
              dense
              @update:model-value="themeStore.toggleFocusHighlight"
            />
          </div>
        </q-card>
      </div>

      <!-- 4. Text-to-Speech (Sintesi Vocale) -->
      <div class="col-12 col-md-6">
        <q-card flat bordered class="rounded-xl q-pa-md bg-white full-height shadow-sm">
          <div class="row items-center justify-between q-mb-sm">
            <div class="flex items-center gap-sm">
              <q-avatar size="40px" color="teal-50" text-color="teal-8" icon="volume_up" />
              <div>
                <div class="text-subtitle1 text-weight-bold text-slate-800">Sintesi Vocale Nativa (TTS)</div>
                <div class="text-caption text-slate-500">Lettura ad alta voce con Web Speech API</div>
              </div>
            </div>
            <q-toggle
              v-model="themeStore.ttsEnabled"
              color="teal"
              size="lg"
              @update:model-value="themeStore.toggleTts"
            />
          </div>

          <div class="q-mt-sm">
            <div class="row justify-between items-center text-caption text-weight-bold text-slate-700">
              <span>Velocità di Lettura</span>
              <span class="text-teal-8">{{ themeStore.ttsRate.toFixed(1) }}x</span>
            </div>
            <q-slider
              v-model="themeStore.ttsRate"
              :min="0.7"
              :max="1.5"
              :step="0.1"
              color="teal"
              @update:model-value="themeStore.setTtsRate"
            />
          </div>

          <div class="q-mt-sm flex justify-between items-center">
            <div class="text-caption text-slate-600">Test audio di verifica:</div>
            <TextToSpeechButton
              text="Benvenuto nel registro elettronico scolastico accessibile e inclusivo."
              color="teal"
              size="md"
              :flat="false"
              round
            />
          </div>
        </q-card>
      </div>

      <!-- 5. Daltonismo & Contrasto Avanzato -->
      <div class="col-12">
        <q-card flat bordered class="rounded-xl q-pa-md bg-white shadow-sm">
          <div class="text-subtitle1 text-weight-bold text-slate-800 flex items-center gap-sm q-mb-md">
            <q-avatar size="36px" color="purple-50" text-color="purple-8" icon="palette" />
            Contrasto Visivo & Filtri per Daltonismo (Color Blindness)
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <div class="text-caption text-weight-bold text-slate-700 q-mb-xs">Filtro Ottico per Daltonismo</div>
              <q-select
                v-model="themeStore.colorblindMode"
                :options="colorblindOptions"
                emit-value
                map-options
                outlined
                dense
                @update:model-value="themeStore.setColorblindMode"
              />
            </div>

            <div class="col-12 col-md-6">
              <div class="text-caption text-weight-bold text-slate-700 q-mb-xs">Modalità di Contrasto & OLED</div>
              <q-select
                v-model="themeStore.highContrastMode"
                :options="contrastOptions"
                emit-value
                map-options
                outlined
                dense
                @update:model-value="themeStore.setHighContrastMode"
              />
            </div>
          </div>

          <div class="q-mt-md row items-center justify-between bg-slate-50 q-pa-sm rounded-lg border border-slate-100">
            <div class="text-caption text-slate-600">
              <q-icon name="verified" color="positive" size="xs" class="q-mr-xs" />
              Tutti i voti e gli esiti associano pattern visivi (<code>▼</code> / <code>✓</code>) oltre al colore, nel rispetto di <strong>WCAG 1.4.1</strong>.
            </div>
            <router-link to="/accessibility-statement" class="text-caption text-primary text-weight-bold">
              Consulta la Dichiarazione AgID →
            </router-link>
          </div>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useThemeStore } from 'src/stores/theme'
import TextToSpeechButton from 'src/components/Common/TextToSpeechButton.vue'

const themeStore = useThemeStore()

const fontFamilyOptions = [
  { label: 'Predefinito di Sistema (Inter / Sans-Serif)', value: 'default' },
  { label: 'OpenDyslexic (Alta Leggibilità DSA)', value: 'opendyslexic' },
  { label: 'Lexend (Ottimizzato per fluidità di lettura)', value: 'lexend' },
  { label: 'Fredoka (Carattere arrotondato e morbido)', value: 'fredoka' },
  { label: 'Roboto (Neutro e chiaro)', value: 'roboto' }
]

const fontSizeOptions = [
  { label: 'Standard (100%)', value: 'normal' },
  { label: 'Grande (110%)', value: 'large' },
  { label: 'Molto Grande (122%)', value: 'xlarge' }
]

const lineHeightOptions = [
  { label: 'Standard (1.5x)', value: 'normal' },
  { label: 'Rilassata (1.8x)', value: 'relaxed' },
  { label: 'Ampia (2.1x)', value: 'loose' }
]

const letterSpacingOptions = [
  { label: 'Standard', value: 'normal' },
  { label: 'Ampio (+0.08em)', value: 'wide' },
  { label: 'Molto Ampio (+0.16em)', value: 'wider' }
]

const wordSpacingOptions = [
  { label: 'Standard', value: 'normal' },
  { label: 'Ampio (+0.18em)', value: 'wide' },
  { label: 'Molto Ampio (+0.32em)', value: 'wider' }
]

const colorblindOptions = [
  { label: 'Nessun filtro (Colori Standard)', value: 'none' },
  { label: 'Protanopia (Insensibilità al rosso)', value: 'protanopia' },
  { label: 'Deuteranopia (Insensibilità al verde)', value: 'deuteranopia' },
  { label: 'Tritanopia (Insensibilità al blu/giallo)', value: 'tritanopia' },
  { label: 'Monocromatico (Scala di Grigi)', value: 'monochrome' }
]

const contrastOptions = [
  { label: 'Contrasto Standard', value: 'none' },
  { label: 'Contrasto Elevato (Bordi netti & +35% contrasto)', value: 'high_contrast' },
  { label: 'OLED Pure Black con Testo Ambra / Giallo', value: 'oled_amber' },
  { label: 'OLED Pure Black con Testo Verde Fosforo', value: 'oled_green' },
  { label: 'Colori Invertiti (Negativo)', value: 'inverted' }
]
</script>

<style scoped>
.accessibility-settings-panel {
  max-width: 1100px;
  margin: 0 auto;
}
</style>
