<template>
  <q-page padding class="bg-slate-50">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none">
          Deleghe Ritiro &amp; Tutela Minori
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">
          Autorizzazione formale a soggetti terzi per l'uscita da scuola e adempimenti ex art. 337 c.c.
        </p>
      </div>

      <div class="row q-gutter-sm items-center">
        <q-select
          v-model="selectedChildId"
          :options="childrenList"
          option-value="id"
          option-label="name"
          emit-value
          map-options
          outlined
          dense
          bg-color="white"
          label="Seleziona Figlio/a"
          style="min-width: 200px"
          @update:model-value="loadDelegates"
        />

        <q-btn
          color="primary"
          icon="person_add"
          label="Nuova Delega Ritiro"
          class="rounded-lg shadow-sm"
          :disable="!selectedChildId"
          @click="openAddDialog"
        />
      </div>
    </div>

    <!-- Legal Custody & Privacy Notice (Art. 337 c.c.) -->
    <q-card class="bg-indigo-50 border-indigo-200 q-mb-lg rounded-xl shadow-soft">
      <q-card-section class="row items-center q-gutter-x-md">
        <q-icon name="gavel" color="indigo-8" size="md" />
        <div class="col">
          <div class="text-subtitle1 text-weight-bold text-indigo-9">
            Disposizioni sull'Affidamento dei Figli e Privacy (Art. 337 c.c.)
          </div>
          <div class="text-caption text-indigo-7">
            In caso di sentenze del Tribunale Ordinario o dei Minorenni relative ad affido esclusivo, limitazioni o decadenza della responsabilità genitoriale, la scuola garantisce l'oscuramento immediato e riservato dei dati anagrafici e scolastici al genitore non affidatario a tutela legale del minore.
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Delegates List -->
    <div v-if="delegates.length === 0" class="q-pa-xl text-center bg-white rounded-xl border-dashed border-slate-200">
      <q-icon name="badge" size="56px" color="slate-400" class="q-mb-sm" />
      <div class="text-h6 text-weight-bold text-slate-700">Nessun soggetto terzo delegato</div>
      <div class="text-body2 text-slate-400">
        Attualmente solo i genitori o tutori legali possono ritirare l'alunno/a da scuola.
      </div>
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-for="d in delegates" :key="d.id" class="col-12 col-md-6 col-lg-4">
        <q-card class="rounded-xl shadow-soft border-slate-100 bg-white">
          <q-card-section>
            <div class="row items-center justify-between q-mb-xs">
              <div class="text-h6 text-weight-bold text-slate-800">
                {{ d.delegate_name }}
              </div>
              <q-badge
                :color="d.is_active ? 'positive' : 'grey-6'"
                class="q-pa-xs text-weight-bold"
              >
                {{ d.is_active ? 'Autorizzato' : 'Revocato' }}
              </q-badge>
            </div>

            <div class="text-subtitle2 text-primary q-mb-sm">
              {{ d.relationship || 'Persona delegata' }}
            </div>

            <q-separator class="q-my-sm" />

            <div class="q-gutter-y-xs text-body2 text-slate-700 q-mb-md">
              <div class="row items-center q-gutter-x-sm">
                <q-icon name="phone" size="xs" color="slate-500" />
                <span>Tel: <strong>{{ d.phone }}</strong></span>
              </div>
              <div class="row items-center q-gutter-x-sm">
                <q-icon name="fingerprint" size="xs" color="slate-500" />
                <span>C.F.: <strong>{{ d.fiscal_code || '-' }}</strong></span>
              </div>
              <div class="row items-center q-gutter-x-sm">
                <q-icon name="event" size="xs" color="slate-500" />
                <span v-if="d.is_permanent">Delega permanente valida per l'A.S.</span>
                <span v-else>Valida fino al: {{ d.valid_until || '-' }}</span>
              </div>
              <div v-if="d.custody_order_ref" class="row items-center q-gutter-x-sm text-caption text-indigo-7">
                <q-icon name="verified" size="xs" />
                <span>Rif. Sentenza: {{ d.custody_order_ref }}</span>
              </div>
            </div>

            <div v-if="d.notes" class="text-caption text-slate-500 italic q-mb-md">
              "{{ d.notes }}"
            </div>

            <div class="row items-center justify-between">
              <a
                v-if="d.id_document_url"
                :href="d.id_document_url"
                target="_blank"
                class="text-caption text-primary row items-center q-gutter-x-xs no-underline"
              >
                <q-icon name="attachment" />
                <span>Doc. Riconoscimento</span>
              </a>
              <span v-else class="text-caption text-slate-400">Nessun allegato</span>

              <q-btn
                v-if="d.is_active"
                flat
                dense
                color="negative"
                icon="cancel"
                label="Revoca"
                @click="revokeDelegate(d.id)"
              />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Dialog: Add Delegate -->
    <q-dialog v-model="showAddDialog" persistent>
      <q-card style="min-width: 480px" class="rounded-xl">
        <q-card-section class="row items-center justify-between">
          <div class="text-h6 text-weight-bold text-slate-800">
            Aggiungi Persona Delegata al Ritiro
          </div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md">
          <q-input
            v-model="newDelegate.delegate_name"
            outlined
            dense
            label="Nome e Cognome Delegato *"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <q-input
            v-model="newDelegate.relationship"
            outlined
            dense
            label="Grado di Parentela o Ruolo (es. Nonno materno, Babysitter) *"
            :rules="[val => !!val || 'Campo obbligatorio']"
          />

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input
                v-model="newDelegate.phone"
                outlined
                dense
                label="Numero di Telefono *"
                :rules="[val => !!val || 'Campo obbligatorio']"
              />
            </div>
            <div class="col-6">
              <q-input
                v-model="newDelegate.fiscal_code"
                outlined
                dense
                label="Codice Fiscale"
              />
            </div>
          </div>

          <q-input
            v-model="newDelegate.id_document_url"
            outlined
            dense
            label="URL / Copia Documento di Riconoscimento (es. /docs/...)"
          />

          <q-toggle
            v-model="newDelegate.is_permanent"
            label="Delega valida per l'intero Anno Scolastico"
          />

          <q-input
            v-if="!newDelegate.is_permanent"
            v-model="newDelegate.valid_until"
            type="date"
            outlined
            dense
            label="Valida fino al"
          />

          <!-- Special legal custody reference -->
          <q-input
            v-model="newDelegate.custody_order_ref"
            outlined
            dense
            label="Eventuale Protocollo Sentenza Tribunale (Art. 337 c.c.)"
            hint="Se presente provvedimento giudiziario depositato in segreteria"
          />

          <q-input
            v-model="newDelegate.notes"
            outlined
            dense
            type="textarea"
            rows="2"
            label="Note aggiuntive per la scuola"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Annulla" v-close-popup />
          <q-btn
            color="primary"
            icon="save"
            label="Conferma Delega"
            :disable="!newDelegate.delegate_name || !newDelegate.phone"
            @click="saveDelegate"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const childrenList = ref([])
const selectedChildId = ref(null)
const delegates = ref([])

const showAddDialog = ref(false)
const newDelegate = ref({
  delegate_name: '',
  relationship: '',
  phone: '',
  fiscal_code: '',
  id_document_url: '',
  is_permanent: true,
  valid_until: '',
  custody_order_ref: '',
  notes: ''
})

async function loadChildren() {
  try {
    const res = await api.get('/users/me/children')
    const raw = res.data?.children || res.data || []
    childrenList.value = raw.map(c => ({
      id: c.id,
      name: `${c.first_name || ''} ${c.last_name || ''}`.trim() || c.username
    }))
    if (childrenList.value.length > 0) {
      selectedChildId.value = childrenList.value[0].id
      await loadDelegates()
    }
  } catch (err) {
    console.error('Errore caricamento figli:', err)
  }
}

async function loadDelegates() {
  if (!selectedChildId.value) return
  try {
    const res = await api.get('/security-passes/delegates/by-student', {
      params: { student_id: selectedChildId.value }
    })
    delegates.value = res.data?.delegates || []
  } catch (err) {
    console.error('Errore caricamento deleghe:', err)
  }
}

function openAddDialog() {
  newDelegate.value = {
    delegate_name: '',
    relationship: '',
    phone: '',
    fiscal_code: '',
    id_document_url: '',
    is_permanent: true,
    valid_until: '',
    custody_order_ref: '',
    notes: ''
  }
  showAddDialog.value = true
}

async function saveDelegate() {
  try {
    await api.post('/security-passes/delegates', {
      student_id: selectedChildId.value,
      ...newDelegate.value
    })
    $q.notify({ type: 'positive', message: 'Delega al ritiro salvata con successo' })
    showAddDialog.value = false
    await loadDelegates()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore salvataggio delega' })
  }
}

async function revokeDelegate(delegateId) {
  try {
    await api.post(`/security-passes/delegates/${delegateId}/revoke`)
    $q.notify({ type: 'positive', message: 'Delega revocata' })
    await loadDelegates()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore revoca delega' })
  }
}

onMounted(() => {
  loadChildren()
})
</script>
