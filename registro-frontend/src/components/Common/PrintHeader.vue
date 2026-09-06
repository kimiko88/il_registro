<template>
  <!--
    PrintHeader.vue
    Intestazione istituzionale visibile SOLO in stampa (@media print).
    In modalità schermo è nascosto tramite CSS in globals.css.

    Uso: inserire in qualsiasi pagina documento (pagella, verbale, circolare)
    come primo figlio del contenuto principale.

    <PrintHeader
      :school-name="school.name"
      :school-subtitle="school.address"
      :document-type="'Pagella 1° Quadrimestre'"
      :academic-year="'2026/2027'"
    />

    Classi CSS aggiuntive per la pagina (aggiunte al <q-page>):
      class="print-optimized"
  -->
  <div class="print-header" role="banner" aria-label="Intestazione istituzionale">
    <p class="print-header__title">{{ schoolName }}</p>
    <p v-if="schoolSubtitle" class="print-header__subtitle">{{ schoolSubtitle }}</p>
    <p class="print-header__meta">
      <template v-if="documentType">{{ documentType }}</template>
      <template v-if="documentType && academicYear"> — </template>
      <template v-if="academicYear">A.S. {{ academicYear }}</template>
      <template v-if="(documentType || academicYear) && printDate"> — </template>
      <template v-if="printDate">{{ t('common.printedOn') }}: {{ formattedDate }}</template>
    </p>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const _props = defineProps({
  /** Nome dell'istituto scolastico */
  schoolName: {
    type: String,
    default: ''
  },
  /** Sottotitolo: indirizzo, codice meccanografico, ecc. */
  schoolSubtitle: {
    type: String,
    default: ''
  },
  /** Tipo di documento: 'Pagella 1° Quadrimestre', 'Verbale di Scrutinio', ecc. */
  documentType: {
    type: String,
    default: ''
  },
  /** Anno scolastico: '2026/2027' */
  academicYear: {
    type: String,
    default: ''
  },
  /**
   * Se true, mostra la data corrente di stampa.
   * Default: true.
   */
  printDate: {
    type: Boolean,
    default: true
  }
})

const formattedDate = computed(() => {
  return new Date().toLocaleDateString(undefined, {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  })
})
</script>

<!--
  Nessuno stile locale: tutto definito in globals.css sotto @media print
  e @media screen per garantire la visibilità condizionale senza conflitti
  con il sistema di theming Quasar.
-->
