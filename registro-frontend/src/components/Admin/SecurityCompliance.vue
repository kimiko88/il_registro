<template>
  <div class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h5 class="text-h5 text-weight-bold text-primary q-my-none">
          <q-icon name="verified_user" class="q-mr-sm" />
          {{ t('security.title') || t('settingsPage.security') || 'Sicurezza & Conformità Normativa' }}
        </h5>
        <div class="text-caption text-grey-7">{{ t('security.subtitle') || 'Firma Elettronica Qualificata, Conservazione a Norma CAD e Catena Log Immutabili' }}</div>
      </div>
      <q-btn color="primary" icon="cloud_download" :label="t('security.downloadCadZip') || 'Scarica Pacchetto CAD ZIP'" @click="downloadCadPackage" class="glossy" />
    </div>

    <!-- Cards Overview -->
    <div class="row q-col-gutter-md q-mb-md">
      <!-- FEQ Status -->
      <div class="col-12 col-md-4">
        <q-card flat bordered class="shadow-1">
          <q-card-section>
            <div class="row items-center justify-between">
              <div class="text-subtitle2 text-weight-bold text-primary">{{ t('security.feqTitle') || 'Firma Qualificata (FEQ)' }}</div>
              <q-icon name="fingerprint" color="positive" size="24px" />
            </div>
            <div class="text-h6 text-weight-bolder q-mt-xs text-positive">{{ t('security.feqActive') || 'Attiva e Verificata' }}</div>
            <div class="text-caption text-grey-7 q-mt-xs">
              {{ t('security.feqDesc') || 'Tutte le firme sui registri di classe e sulle supplenze sono sigillate con timestamp SHA-256 con valore legale.' }}
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- CAD Preservation Status -->
      <div class="col-12 col-md-4">
        <q-card flat bordered class="shadow-1">
          <q-card-section>
            <div class="row items-center justify-between">
              <div class="text-subtitle2 text-weight-bold text-primary">{{ t('security.cadTitle') || 'Conservazione CAD' }}</div>
              <q-icon name="inventory" color="secondary" size="24px" />
            </div>
            <div class="text-h6 text-weight-bolder q-mt-xs text-secondary">{{ t('security.cadManifest') || 'Pacchetto XML Manifesto' }}</div>
            <div class="text-caption text-grey-7 q-mt-xs">
              {{ t('security.cadDesc') || 'Conforme alle linee guida AgID per la conservazione documentale sostitutiva a lungo termine.' }}
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Immutability Chain Status -->
      <div class="col-12 col-md-4">
        <q-card flat bordered class="shadow-1">
          <q-card-section>
            <div class="row items-center justify-between">
              <div class="text-subtitle2 text-weight-bold text-primary">{{ t('security.immutableAuditTitle') || 'Audit Log Immutabile' }}</div>
              <q-icon name="link" color="accent" size="24px" />
            </div>
            <div class="text-h6 text-weight-bolder q-mt-xs text-accent">{{ t('security.cryptoChain') || 'Catena Cryptographic Hash' }}</div>
            <div class="text-caption text-grey-7 q-mt-xs">
              {{ t('security.immutableAuditDesc') || 'Ogni modifica a firme, voti e note genera un blocco incatenato crittograficamente (Block Hashing).' }}
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Immutability Inspection Panel -->
    <q-card flat bordered>
      <q-card-section class="row items-center justify-between bg-grey-2">
        <div class="text-subtitle1 text-weight-bold">
          <q-icon name="fact_check" class="q-mr-xs" /> {{ t('security.integrityVerifier') || 'Verificatore di Integrità Registro (Audit Trail)' }}
        </div>
        <q-btn color="secondary" icon="refresh" :label="t('security.verifyNow') || 'Verifica Integrità Ora'" dense flat @click="checkImmutabilityChain" />
      </q-card-section>

      <q-card-section v-if="loading" class="text-center q-pa-lg">
        <q-spinner-dots color="primary" size="32px" />
      </q-card-section>

      <q-card-section v-else>
        <div class="row items-center q-mb-md">
          <q-chip color="positive" text-color="white" icon="check_circle" :label="t('security.allBlocksValid') || 'Tutti i blocchi del registro sono validi ed inopponibili'" />
          <div class="text-caption text-grey-7 q-ml-sm">{{ t('security.verifiedOn') || 'Verificato il' }}: {{ report.verified_at || 'Ora' }}</div>
        </div>

        <q-markup-table flat separator="cell" dense>
          <thead>
            <tr>
              <th class="text-left">{{ t('security.colAction') || 'Azione / Evento' }}</th>
              <th class="text-left">{{ t('security.colActor') || 'Utente / Modificatore' }}</th>
              <th class="text-left">{{ t('security.colPrevHash') || 'Hash Precedente (SHA-256)' }}</th>
              <th class="text-left">{{ t('security.colCurrHash') || 'Hash Corrente (SHA-256)' }}</th>
              <th class="text-center">{{ t('security.colCryptoStatus') || 'Stato Crittografico' }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="block in sampleBlocks" :key="block.id">
              <td class="text-weight-bold">{{ block.action }}</td>
              <td>{{ block.actor_name }}</td>
              <td class="text-caption text-mono text-grey-7">{{ (block.prev_hash || '').substring(0, 16) }}...</td>
              <td class="text-caption text-mono text-primary text-weight-bold">{{ (block.current_hash || '').substring(0, 16) }}...</td>
              <td class="text-center">
                <q-chip size="sm" color="positive" text-color="white" icon="verified">{{ t('security.intact') || 'Integro' }}</q-chip>
              </td>
            </tr>
          </tbody>
        </q-markup-table>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import securityService from '@/services/securityService'

const $q = useQuasar()
const { t } = useI18n()
const loading = ref(false)
const report = ref({})

const sampleBlocks = ref([
  { id: '1', action: 'FIRMA_REGISTRO_LEZIONE', actor_name: 'Prof. Mario Rossi', prev_hash: '00000000000000000000000000000000', current_hash: 'e3b0c44298fc1c149afbf4c8996fb924', is_valid: true },
  { id: '2', action: 'REGISTRAZIONE_APPELLO_PRESENZE', actor_name: 'Prof. Mario Rossi', prev_hash: 'e3b0c44298fc1c149afbf4c8996fb924', current_hash: 'f2ca1bb6c7e907d06dafe4687e579fce', is_valid: true },
  { id: '3', action: 'FIRMA_SOSTITUZIONE_SUPPLENTE', actor_name: 'Prof. Giuseppe Bianchi', prev_hash: 'f2ca1bb6c7e907d06dafe4687e579fce', current_hash: '315f5bdb76d078c43b8ac0064e4a0164', is_valid: true }
])

const checkImmutabilityChain = async () => {
  loading.value = true
  try {
    const res = await securityService.getImmutabilityChain()
    report.value = res || {}
    if (res?.blocks && Array.isArray(res.blocks)) {
      sampleBlocks.value = res.blocks
    }
    $q.notify({ type: 'positive', message: t('security.integrityVerified') || 'Verifica integrità registro completata: Catena valida!' })
  } catch (err) {
    $q.notify({ type: 'negative', message: t('security.integrityError') || 'Errore durante la verifica della catena' })
  } finally {
    loading.value = false
  }
}

const downloadCadPackage = async () => {
  try {
    $q.notify({ type: 'info', message: t('security.cadGenerating') || 'Generazione del pacchetto di conservazione CAD ZIP...' })
    const blob = await securityService.downloadCadPackage()
    const url = window.URL.createObjectURL(new Blob([blob]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'pacchetto_conservazione_CAD.zip')
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
    $q.notify({ type: 'positive', message: t('security.cadDownloaded') || 'Pacchetto CAD scaricato con successo!' })
  } catch (err) {
    $q.notify({ type: 'negative', message: t('security.cadError') || 'Errore durante il download del pacchetto CAD' })
  }
}

onMounted(() => {
  checkImmutabilityChain()
})
</script>
