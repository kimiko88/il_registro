<template>
  <q-page class="q-pa-md q-pa-lg-xl strike-management-page">
    <!-- Header Hero Section -->
    <div class="row items-center justify-between q-mb-lg gap-md">
      <div class="col-12 col-md-7">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge color="indigo-7" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="campaign" size="14px" class="q-mr-xs" />
            BACHECA & RILEVAZIONE PREVENTIVA
          </q-badge>
          <q-badge color="teal-8" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="gavel" size="14px" class="q-mr-xs" />
            ACCORDO ARAN 2/12/2020
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none flex items-center">
          <q-icon name="ballot" color="primary" class="q-mr-sm" size="36px" />
          Rilevazione Preventiva Scioperi
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          Comunicazione scioperi in bacheca, impostazione scadenze per le dichiarazioni del personale e quadro preventivo delle adesioni
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="col-12 col-md-5 row items-center justify-end q-gutter-sm">
        <q-btn
          color="primary"
          icon="add"
          label="Nuovo Avviso di Sciopero"
          no-caps
          rounded
          class="shadow-2 text-weight-bold q-px-md"
          @click="openCreateDialog"
        />
        <q-btn
          outline
          color="primary"
          icon="refresh"
          round
          dense
          @click="loadNotices"
          :loading="loadingNotices"
        >
          <q-tooltip>Aggiorna elenco scioperi</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Legal & Operating Guideline Banner -->
    <div class="q-mb-lg bg-indigo-50 border border-indigo-100 rounded-xl q-pa-md row items-center justify-between">
      <div class="row items-center gap-sm">
        <q-avatar size="40px" color="indigo-600" text-color="white" icon="info" />
        <div>
          <div class="text-weight-bold text-slate-800">
            Quadro Preventivo per la DSGA e la Dirigenza Scolastica
          </div>
          <div class="text-caption text-slate-600">
            I dipendenti (docenti e ATA) possono comunicare preventivamente: <b>Aderisco</b>, <b>Non aderisco</b> o <b>Non ho ancora deciso</b>. Oltre il termine indicato non sarà più possibile inserire o modificare la scelta.
          </div>
        </div>
      </div>
      <div class="row items-center q-gutter-xs">
        <q-btn
          flat
          dense
          color="primary"
          icon="co_present"
          label="Vai a Presenze Personale"
          to="/ata/attendance"
          class="text-weight-bold"
        />
      </div>
    </div>

    <!-- Active Strike Notice Selector / Carousel -->
    <div v-if="notices.length === 0 && !loadingNotices" class="q-pa-xl text-center bg-white rounded-xl border border-slate-200 shadow-sm">
      <q-icon name="event_busy" size="64px" color="slate-300" class="q-mb-md" />
      <div class="text-h6 text-weight-bold text-slate-700">Nessuna comunicazione di sciopero registrata</div>
      <div class="text-body2 text-slate-500 q-mt-xs q-mb-md">
        Crea la prima comunicazione per pubblicarla nella bacheca del personale e avviare la raccolta preventiva delle intenzioni.
      </div>
      <q-btn
        color="primary"
        icon="add"
        label="Crea Avviso di Sciopero"
        no-caps
        rounded
        class="text-weight-bold"
        @click="openCreateDialog"
      />
    </div>

    <div v-else>
      <!-- Notice Selector Cards Grid -->
      <div class="q-mb-lg">
        <div class="row items-center justify-between q-mb-sm">
          <div class="text-subtitle1 text-weight-bold text-slate-800 flex items-center">
            <q-icon name="list_alt" color="primary" class="q-mr-xs" size="20px" />
            Avvisi di Sciopero Registrati ({{ notices.length }})
          </div>
        </div>

        <div class="row q-col-gutter-md">
          <div
            v-for="item in notices"
            :key="item.id"
            class="col-12 col-md-6 col-lg-4"
          >
            <q-card
              class="rounded-xl border shadow-sm cursor-pointer transition-all notice-card"
              :class="{
                'notice-card-selected': selectedNoticeId === item.id,
                'border-primary': selectedNoticeId === item.id,
                'border-slate-200': selectedNoticeId !== item.id
              }"
              @click="selectNotice(item.id)"
            >
              <q-card-section class="q-pb-sm">
                <div class="row items-center justify-between no-wrap q-mb-xs">
                  <q-badge
                    :color="item.is_expired ? 'grey-7' : 'negative'"
                    text-color="white"
                    class="text-weight-bold q-px-sm q-py-xs rounded-borders"
                  >
                    <q-icon :name="item.is_expired ? 'lock' : 'campaign'" size="13px" class="q-mr-xs" />
                    {{ item.is_expired ? 'Termine Scaduto' : 'Dichiarazioni Aperte' }}
                  </q-badge>

                  <q-btn
                    flat
                    round
                    dense
                    color="grey-6"
                    icon="delete_outline"
                    @click.stop="confirmDeleteNotice(item)"
                  >
                    <q-tooltip>Elimina avviso</q-tooltip>
                  </q-btn>
                </div>

                <div class="text-subtitle1 text-weight-bold text-slate-900 line-clamp-2 q-mt-xs">
                  {{ item.title }}
                </div>

                <div class="text-caption text-slate-500 q-mt-xs">
                  Proclamato da: <span class="text-weight-bold text-slate-700">{{ item.proclaimed_by || 'Organizzazioni Sindacali' }}</span>
                </div>
              </q-card-section>

              <q-separator />

              <q-card-section class="q-py-sm text-caption text-slate-600">
                <div class="row items-center justify-between q-mb-xs">
                  <span><q-icon name="event" size="15px" class="q-mr-xs text-primary" />Data Sciopero:</span>
                  <span class="text-weight-bold text-slate-800">{{ formatDate(item.strike_date) }}</span>
                </div>
                <div class="row items-center justify-between">
                  <span><q-icon name="timer" size="15px" class="q-mr-xs text-deep-orange-8" />Termine Scelta:</span>
                  <span class="text-weight-bold" :class="item.is_expired ? 'text-grey-7' : 'text-deep-orange-9'">
                    {{ formatDateTime(item.declaration_deadline) }}
                  </span>
                </div>
              </q-card-section>

              <q-card-actions align="between" class="bg-slate-50 q-px-md q-py-xs border-t border-slate-100">
                <span class="text-caption text-slate-500">
                  Creato il {{ formatDate(item.created_at) }}
                </span>
                <q-btn
                  flat
                  dense
                  no-caps
                  color="primary"
                  label="Vedi Quadro"
                  icon-right="arrow_forward"
                  class="text-weight-bold"
                  @click.stop="selectNotice(item.id)"
                />
              </q-card-actions>
            </q-card>
          </div>
        </div>
      </div>

      <!-- Detail View: Quadro Preventivo for Selected Notice -->
      <div v-if="loadingSummary" class="q-pa-xl text-center">
        <q-spinner-dots color="primary" size="48px" />
        <div class="text-slate-500 q-mt-md">Caricamento quadro preventivo e adesioni...</div>
      </div>

      <div v-else-if="summary" class="summary-container q-mt-xl">
        <!-- Selected Notice Banner -->
        <q-card class="rounded-xl shadow-sm border border-slate-200 q-mb-lg bg-white">
          <q-card-section class="q-pa-md q-pa-lg-lg">
            <div class="row items-start justify-between q-col-gutter-md">
              <div class="col-12 col-lg-8">
                <div class="row items-center q-gutter-sm q-mb-xs">
                  <q-badge color="primary" text-color="white" class="text-weight-bold q-px-sm">
                    QUADRO PREVENTIVO ATTIVO
                  </q-badge>
                  <q-badge
                    :color="summary.notice.is_expired ? 'grey-7' : 'positive'"
                    text-color="white"
                    class="text-weight-bold q-px-sm"
                  >
                    {{ summary.notice.is_expired ? 'Termine Inserimento Chiuso' : 'Raccolta Intenzioni in Corso' }}
                  </q-badge>
                </div>
                <h2 class="text-h5 text-weight-bolder text-slate-900 q-my-none">
                  {{ summary.notice.title }}
                </h2>
                <div class="text-body2 text-slate-600 q-mt-xs">
                  Proclamato da: <b>{{ summary.notice.proclaimed_by }}</b>
                  <span class="q-mx-sm">•</span>
                  Data di svolgimento: <b>{{ formatDate(summary.notice.strike_date) }}</b>
                  <span class="q-mx-sm">•</span>
                  Termine limite per il personale: <b>{{ formatDateTime(summary.notice.declaration_deadline) }}</b>
                </div>
                <div v-if="summary.notice.notes" class="bg-slate-50 border border-slate-200 rounded-borders q-pa-sm q-mt-sm text-caption text-slate-700">
                  <q-icon name="sticky_note_2" color="primary" class="q-mr-xs" />
                  <b>Disposizioni / Note di Servizio:</b> {{ summary.notice.notes }}
                </div>
              </div>

              <div class="col-12 col-lg-4 text-right">
                <q-btn
                  color="positive"
                  icon="download"
                  label="Esporta CSV Nominativo"
                  no-caps
                  rounded
                  class="text-weight-bold shadow-1 q-px-md"
                  @click="exportCsv"
                />
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- KPI Metrics Grid (Totale, Aderisco, Non Aderisco, Indecisi, In Attesa) -->
        <div class="row q-col-gutter-md q-mb-lg">
          <!-- 1. Totale Personale & Risposte -->
          <div class="col-12 col-sm-6 col-lg-3">
            <q-card class="rounded-xl border border-slate-200 shadow-sm kpi-card bg-white">
              <q-card-section class="q-pb-xs">
                <div class="row items-center justify-between">
                  <span class="text-caption text-weight-bold text-slate-500 uppercase">Personale Coinvolto</span>
                  <q-avatar size="32px" color="indigo-50" text-color="indigo-7" icon="people" />
                </div>
                <div class="text-h4 text-weight-bolder text-slate-800 q-my-xs">
                  {{ summary.total_staff }}
                </div>
                <div class="text-caption text-slate-600">
                  <b class="text-primary">{{ summary.answered_count }}</b> hanno risposto ({{ summary.response_rate }}%)
                </div>
              </q-card-section>
              <q-linear-progress
                :value="summary.total_staff > 0 ? (summary.answered_count / summary.total_staff) : 0"
                color="indigo-6"
                track-color="indigo-1"
                size="4px"
              />
            </q-card>
          </div>

          <!-- 2. Aderiscono allo Sciopero -->
          <div class="col-12 col-sm-6 col-lg-3">
            <q-card class="rounded-xl border border-positive-subtle shadow-sm kpi-card bg-emerald-50">
              <q-card-section class="q-pb-xs">
                <div class="row items-center justify-between">
                  <span class="text-caption text-weight-bold text-positive uppercase">Aderiscono</span>
                  <q-avatar size="32px" color="positive" text-color="white" icon="check_circle" />
                </div>
                <div class="text-h4 text-weight-bolder text-positive q-my-xs">
                  {{ summary.participates_count }}
                </div>
                <div class="text-caption text-positive">
                  <b>{{ summary.participates_rate }}%</b> del personale complessivo
                </div>
              </q-card-section>
              <q-linear-progress
                :value="summary.total_staff > 0 ? (summary.participates_count / summary.total_staff) : 0"
                color="positive"
                track-color="emerald-100"
                size="4px"
              />
            </q-card>
          </div>

          <!-- 3. Non Aderiscono allo Sciopero -->
          <div class="col-12 col-sm-6 col-lg-3">
            <q-card class="rounded-xl border border-negative-subtle shadow-sm kpi-card bg-rose-50">
              <q-card-section class="q-pb-xs">
                <div class="row items-center justify-between">
                  <span class="text-caption text-weight-bold text-negative uppercase">Non Aderiscono</span>
                  <q-avatar size="32px" color="negative" text-color="white" icon="cancel" />
                </div>
                <div class="text-h4 text-weight-bolder text-negative q-my-xs">
                  {{ summary.not_participates_count }}
                </div>
                <div class="text-caption text-negative">
                  <b>{{ summary.not_participates_rate }}%</b> del personale complessivo
                </div>
              </q-card-section>
              <q-linear-progress
                :value="summary.total_staff > 0 ? (summary.not_participates_count / summary.total_staff) : 0"
                color="negative"
                track-color="rose-100"
                size="4px"
              />
            </q-card>
          </div>

          <!-- 4. Indecisi / Non hanno deciso -->
          <div class="col-12 col-sm-6 col-lg-3">
            <q-card class="rounded-xl border border-amber-200 shadow-sm kpi-card bg-amber-50">
              <q-card-section class="q-pb-xs">
                <div class="row items-center justify-between">
                  <span class="text-caption text-weight-bold text-amber-9 uppercase">Non Ancora Deciso</span>
                  <q-avatar size="32px" color="amber-8" text-color="white" icon="help_outline" />
                </div>
                <div class="text-h4 text-weight-bolder text-amber-9 q-my-xs">
                  {{ summary.undecided_count }}
                </div>
                <div class="text-caption text-amber-9">
                  <b>{{ summary.undecided_rate }}%</b> • In attesa: <b>{{ summary.unanswered_count }}</b>
                </div>
              </q-card-section>
              <q-linear-progress
                :value="summary.total_staff > 0 ? (summary.undecided_count / summary.total_staff) : 0"
                color="amber-8"
                track-color="amber-100"
                size="4px"
              />
            </q-card>
          </div>
        </div>

        <!-- Role Breakdown Section -->
        <div class="q-mb-xl">
          <div class="row items-center justify-between q-mb-md">
            <div>
              <h2 class="text-h6 text-weight-bold text-slate-800 q-my-none flex items-center">
                <q-icon name="pie_chart" color="primary" class="q-mr-sm" size="22px" />
                Ripartizione per Categoria Professionale
              </h2>
              <div class="text-caption text-slate-500">
                Adesioni preventive distinte tra Docenti, Personale Amministrativo e Collaboratori Scolastici
              </div>
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div
              v-for="cat in summary.by_role"
              :key="cat.role"
              class="col-12 col-sm-6 col-lg"
            >
              <q-card class="rounded-xl border border-slate-200 shadow-sm full-height bg-white">
                <q-card-section class="q-pb-xs">
                  <div class="row items-center justify-between no-wrap q-mb-xs">
                    <q-chip
                      dense
                      :color="getRoleBadgeColor(cat.role)"
                      text-color="white"
                      class="text-weight-bolder text-caption"
                    >
                      {{ cat.role_display }}
                    </q-chip>
                    <div class="text-caption text-weight-bold text-slate-400">
                      Totale: {{ cat.total }}
                    </div>
                  </div>

                  <div class="row items-center justify-between q-mt-sm">
                    <span class="text-caption text-positive text-weight-bold">
                      <q-icon name="check_circle" size="14px" /> Aderiscono: {{ cat.participates }}
                    </span>
                    <span class="text-caption text-negative text-weight-bold">
                      <q-icon name="cancel" size="14px" /> Non ad.: {{ cat.not_participates }}
                    </span>
                  </div>

                  <div class="row items-center justify-between q-mt-xs text-caption text-slate-600">
                    <span><q-icon name="help" size="14px" class="text-amber-8" /> Indecisi: {{ cat.undecided }}</span>
                    <span><q-icon name="hourglass_empty" size="14px" class="text-grey-6" /> In attesa: {{ cat.unanswered }}</span>
                  </div>

                  <!-- Stacked-like visual indicators -->
                  <q-linear-progress
                    :value="cat.total > 0 ? (cat.participates / cat.total) : 0"
                    color="positive"
                    track-color="slate-100"
                    rounded
                    size="6px"
                    class="q-mt-sm"
                  />
                </q-card-section>
              </q-card>
            </div>
          </div>
        </div>

        <!-- Detailed Nominative Table Section -->
        <q-card class="rounded-xl border border-slate-200 shadow-sm bg-white overflow-hidden">
          <q-card-section class="q-pb-none">
            <div class="row items-center justify-between q-mb-md gap-sm">
              <div>
                <h3 class="text-h6 text-weight-bold text-slate-800 q-my-none flex items-center">
                  <q-icon name="format_list_bulleted" color="primary" class="q-mr-sm" size="22px" />
                  Elenco Nominativo Personale e Dichiarazioni Preventivi
                </h3>
                <div class="text-caption text-slate-500">
                  Consultazione nominativa delle scelte espresse entro la scadenza per la predisposizione dell'orario e dei servizi minimi
                </div>
              </div>

              <div class="row items-center q-gutter-sm">
                <q-input
                  v-model="filterSearch"
                  dense
                  outlined
                  rounded
                  placeholder="Cerca dipendente..."
                  class="search-input"
                >
                  <template v-slot:prepend>
                    <q-icon name="search" />
                  </template>
                  <template v-slot:append v-if="filterSearch">
                    <q-icon name="close" class="cursor-pointer" @click="filterSearch = ''" />
                  </template>
                </q-input>

                <q-select
                  v-model="filterRole"
                  :options="roleFilterOptions"
                  emit-value
                  map-options
                  dense
                  outlined
                  rounded
                  class="role-select"
                />

                <q-select
                  v-model="filterIntention"
                  :options="intentionFilterOptions"
                  emit-value
                  map-options
                  dense
                  outlined
                  rounded
                  class="intention-select"
                />
              </div>
            </div>
          </q-card-section>

          <q-table
            :rows="filteredStaff"
            :columns="columns"
            row-key="user_id"
            flat
            :pagination="{ rowsPerPage: 15 }"
            class="nominative-table"
            no-data-label="Nessun dipendente trovato con i filtri selezionati"
          >
            <!-- Col: Person / Name -->
            <template v-slot:body-cell-person="props">
              <q-td :props="props">
                <div class="row items-center no-wrap">
                  <q-avatar size="32px" color="indigo-1" text-color="indigo-8" class="q-mr-sm text-weight-bold">
                    {{ props.row.first_name ? props.row.first_name.charAt(0) : '' }}{{ props.row.last_name ? props.row.last_name.charAt(0) : '' }}
                  </q-avatar>
                  <div>
                    <div class="text-weight-bold text-slate-800">
                      {{ props.row.last_name }} {{ props.row.first_name }}
                    </div>
                    <div class="text-caption text-slate-400">
                      {{ props.row.email }}
                    </div>
                  </div>
                </div>
              </q-td>
            </template>

            <!-- Col: Role -->
            <template v-slot:body-cell-role="props">
              <q-td :props="props">
                <q-chip
                  dense
                  :color="getRoleBadgeColor(props.row.role)"
                  text-color="white"
                  class="text-weight-bold text-caption"
                >
                  {{ props.row.role_display }}
                </q-chip>
              </q-td>
            </template>

            <!-- Col: Intention -->
            <template v-slot:body-cell-intention="props">
              <q-td :props="props" align="center">
                <q-chip
                  dense
                  :color="getIntentionBadgeColor(props.row.intention)"
                  text-color="white"
                  class="text-weight-bold text-caption"
                >
                  <q-icon :name="getIntentionIcon(props.row.intention)" size="14px" class="q-mr-xs" />
                  {{ getIntentionLabel(props.row.intention) }}
                </q-chip>
              </q-td>
            </template>

            <!-- Col: Declared At -->
            <template v-slot:body-cell-declared_at="props">
              <q-td :props="props">
                <div v-if="props.row.declared_at" class="text-caption text-slate-700">
                  {{ formatDateTime(props.row.declared_at) }}
                </div>
                <div v-else class="text-caption text-slate-400 text-italic">
                  Non pervenuta
                </div>
              </q-td>
            </template>
          </q-table>
        </q-card>
      </div>
    </div>

    <!-- Dialog: Create New Strike Notice -->
    <q-dialog v-model="showCreateDialog" persistent>
      <q-card class="rounded-xl shadow-4" style="min-width: 580px; max-width: 95vw;">
        <q-card-section class="row items-center justify-between bg-primary text-white">
          <div class="text-h6 text-weight-bold flex items-center">
            <q-icon name="campaign" size="24px" class="q-mr-sm" />
            Nuova Comunicazione di Sciopero
          </div>
          <q-btn flat round dense icon="close" color="white" v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <div class="text-caption text-slate-600">
            Pubblica l'avviso di sciopero in bacheca e attiva la rilevazione preventiva con termine perentorio di risposta per il personale.
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Titolo Sciopero / Oggetto *</label>
            <q-input
              v-model="createForm.title"
              outlined
              dense
              placeholder="es. Sciopero generale del comparto Istruzione e Ricerca"
              :rules="[val => !!val || 'Il titolo è obbligatorio']"
            />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Proclamato da (Sigle sindacali) *</label>
            <q-input
              v-model="createForm.proclaimed_by"
              outlined
              dense
              placeholder="es. FLC CGIL, CISL FSUR, UIL SCUOLA RUA, SNALS CONFSAL"
              :rules="[val => !!val || 'Le sigle proclamanti sono obbligatorie']"
            />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Data dello Sciopero *</label>
              <q-input
                v-model="createForm.strike_date"
                outlined
                dense
                type="date"
                :rules="[val => !!val || 'La data di sciopero è obbligatoria']"
              />
            </div>

            <div class="col-12 col-sm-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Termine Ultimo Scelta (Data e Ora) *</label>
              <q-input
                v-model="createForm.declaration_deadline"
                outlined
                dense
                type="datetime-local"
                :rules="[
                  val => !!val || 'La scadenza è obbligatoria',
                  val => isDeadlineValid(val, createForm.strike_date) || 'La scadenza deve precedere o coincidere con la data di sciopero'
                ]"
              />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">Note, Circolare e Disposizioni di Servizio</label>
            <q-input
              v-model="createForm.notes"
              type="textarea"
              rows="3"
              outlined
              dense
              placeholder="Indicazioni per il personale, funzionamento dei servizi minimi essenziali, orari ridotti previsti..."
            />
          </div>

          <div class="bg-indigo-50 border border-indigo-100 rounded-borders q-pa-sm">
            <q-checkbox
              v-model="createForm.publish_to_bacheca"
              label="Pubblica contestualmente come Comunicazione / Circolare in Bacheca"
              color="primary"
              class="text-weight-bold text-slate-800"
            />
            <div class="text-caption text-slate-600 q-ml-lg">
              L'avviso apparirà anche nella bacheca comunicazioni con priorità alta per tutti i docenti e il personale ATA.
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat label="Annulla" color="grey-7" v-close-popup />
          <q-btn
            color="primary"
            label="Pubblica Comunicazione"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="savingNotice"
            @click="submitCreateNotice"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Confirm Delete Notice -->
    <q-dialog v-model="showDeleteDialog">
      <q-card class="rounded-xl shadow-3" style="min-width: 400px;">
        <q-card-section class="row items-center q-pb-none">
          <q-avatar icon="warning" color="negative" text-color="white" />
          <span class="q-ml-sm text-h6 text-weight-bold">Conferma Eliminazione</span>
        </q-card-section>

        <q-card-section class="q-pt-sm text-body2 text-slate-600">
          Sei sicuro di voler eliminare la comunicazione di sciopero:
          <br />
          <b class="text-slate-900">{{ noticeToDelete?.title }}</b>?
          <br />
          Tutte le dichiarazioni preventive raccolte per questo avviso verranno rimosse.
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Annulla" color="grey-7" v-close-popup />
          <q-btn
            color="negative"
            label="Elimina Definitivamente"
            no-caps
            rounded
            :loading="deletingNotice"
            @click="executeDeleteNotice"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import strikeService from '@/services/strikeService'

const $q = useQuasar()

const notices = ref([])
const selectedNoticeId = ref(null)
const summary = ref(null)

const loadingNotices = ref(false)
const loadingSummary = ref(false)
const savingNotice = ref(false)
const deletingNotice = ref(false)

const showCreateDialog = ref(false)
const showDeleteDialog = ref(false)
const noticeToDelete = ref(null)

const filterSearch = ref('')
const filterRole = ref('all')
const filterIntention = ref('all')

const createForm = ref({
  title: '',
  proclaimed_by: '',
  strike_date: '',
  declaration_deadline: '',
  notes: '',
  publish_to_bacheca: true
})

const columns = [
  { name: 'person', label: 'Dipendente', align: 'left', field: 'last_name', sortable: true },
  { name: 'role', label: 'Categoria', align: 'left', field: 'role', sortable: true },
  { name: 'intention', label: 'Intenzione Espressa', align: 'center', field: 'intention', sortable: true },
  { name: 'declared_at', label: 'Data & Ora Scelta', align: 'left', field: 'declared_at', sortable: true }
]

const roleFilterOptions = [
  { label: 'Tutte le categorie', value: 'all' },
  { label: 'Docenti', value: 'teacher' },
  { label: 'DSGA', value: 'dsga' },
  { label: 'Assistenti Amministrativi', value: 'assistente_amministrativo' },
  { label: 'Collaboratori D.S.', value: 'collaboratore_ds' },
  { label: 'Collaboratori Scolastici', value: 'collaboratore_scolastico' },
  { label: 'Segreteria', value: 'secretary' }
]

const intentionFilterOptions = [
  { label: 'Tutte le intenzioni', value: 'all' },
  { label: 'Aderisce (Aderisco)', value: 'participates' },
  { label: 'Non aderisce', value: 'not_participates' },
  { label: 'Non ha ancora deciso', value: 'undecided' },
  { label: 'Nessuna risposta', value: 'unanswered' }
]

const filteredStaff = computed(() => {
  if (!summary.value || !summary.value.staff) return []
  let list = summary.value.staff

  if (filterRole.value !== 'all') {
    list = list.filter(item => item.role === filterRole.value)
  }

  if (filterIntention.value !== 'all') {
    list = list.filter(item => item.intention === filterIntention.value)
  }

  if (filterSearch.value && filterSearch.value.trim() !== '') {
    const q = filterSearch.value.trim().toLowerCase()
    list = list.filter(item => {
      const full = `${item.first_name} ${item.last_name}`.toLowerCase()
      const email = (item.email || '').toLowerCase()
      return full.includes(q) || email.includes(q)
    })
  }

  return list
})

const loadNotices = async () => {
  loadingNotices.value = true
  try {
    const data = await strikeService.getStrikeNotices()
    notices.value = data || []
    if (notices.value.length > 0) {
      if (!selectedNoticeId.value || !notices.value.some(n => n.id === selectedNoticeId.value)) {
        selectedNoticeId.value = notices.value[0].id
      }
      await loadSummary(selectedNoticeId.value)
    } else {
      selectedNoticeId.value = null
      summary.value = null
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: 'Errore nel caricamento delle comunicazioni di sciopero'
    })
  } finally {
    loadingNotices.value = false
  }
}

const selectNotice = async (noticeId) => {
  selectedNoticeId.value = noticeId
  await loadSummary(noticeId)
}

const loadSummary = async (noticeId) => {
  if (!noticeId) return
  loadingSummary.value = true
  try {
    const data = await strikeService.getNoticeSummary(noticeId)
    summary.value = data
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: 'Errore nel recupero del quadro riepilogativo adesioni'
    })
  } finally {
    loadingSummary.value = false
  }
}

const openCreateDialog = () => {
  // Default dates: strike date in 5 days, deadline in 3 days at 12:00
  const futureStrike = new Date()
  futureStrike.setDate(futureStrike.getDate() + 5)
  const strikeDateStr = futureStrike.toISOString().substring(0, 10)

  const deadline = new Date()
  deadline.setDate(deadline.getDate() + 3)
  deadline.setHours(12, 0, 0, 0)
  // Format for datetime-local input: YYYY-MM-DDTHH:mm
  const y = deadline.getFullYear()
  const m = String(deadline.getMonth() + 1).padStart(2, '0')
  const d = String(deadline.getDate()).padStart(2, '0')
  const deadlineStr = `${y}-${m}-${d}T12:00`

  createForm.value = {
    title: '',
    proclaimed_by: '',
    strike_date: strikeDateStr,
    declaration_deadline: deadlineStr,
    notes: 'Ai sensi dell\'Accordo ARAN 2/12/2020, si invita il personale a comunicare la propria intenzione entro il termine indicato per consentire l\'organizzazione dei servizi minimi.',
    publish_to_bacheca: true
  }
  showCreateDialog.value = true
}

const isDeadlineValid = (deadlineVal, strikeDateVal) => {
  if (!deadlineVal || !strikeDateVal) return true
  const dTime = new Date(deadlineVal).getTime()
  const sTime = new Date(strikeDateVal + 'T23:59:59').getTime()
  return dTime <= sTime
}

const submitCreateNotice = async () => {
  if (!createForm.value.title || !createForm.value.proclaimed_by || !createForm.value.strike_date || !createForm.value.declaration_deadline) {
    $q.notify({
      type: 'warning',
      message: 'Compila tutti i campi obbligatori contrassegnati da *'
    })
    return
  }

  savingNotice.value = true
  try {
    // Ensure ISO string with timezone
    const deadlineIso = new Date(createForm.value.declaration_deadline).toISOString()

    const payload = {
      title: createForm.value.title,
      proclaimed_by: createForm.value.proclaimed_by,
      strike_date: createForm.value.strike_date,
      declaration_deadline: deadlineIso,
      notes: createForm.value.notes,
      publish_to_bacheca: createForm.value.publish_to_bacheca
    }

    const created = await strikeService.createStrikeNotice(payload)
    showCreateDialog.value = false
    $q.notify({
      type: 'positive',
      message: 'Comunicazione di sciopero pubblicata con successo!',
      caption: payload.publish_to_bacheca ? 'Notifica inoltrata anche in bacheca' : ''
    })
    await loadNotices()
    if (created && created.id) {
      selectedNoticeId.value = created.id
      await loadSummary(created.id)
    }
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore nella creazione della comunicazione'
    })
  } finally {
    savingNotice.value = false
  }
}

const confirmDeleteNotice = (item) => {
  noticeToDelete.value = item
  showDeleteDialog.value = true
}

const executeDeleteNotice = async () => {
  if (!noticeToDelete.value) return
  deletingNotice.value = true
  try {
    await strikeService.deleteStrikeNotice(noticeToDelete.value.id)
    showDeleteDialog.value = false
    $q.notify({
      type: 'positive',
      message: 'Avviso di sciopero eliminato con successo'
    })
    noticeToDelete.value = null
    await loadNotices()
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: 'Errore durante l\'eliminazione'
    })
  } finally {
    deletingNotice.value = false
  }
}

const exportCsv = () => {
  if (!summary.value || !summary.value.staff || summary.value.staff.length === 0) {
    $q.notify({
      type: 'warning',
      message: 'Nessun dato da esportare'
    })
    return
  }

  const rows = [
    ['Cognome', 'Nome', 'Email', 'Ruolo', 'Categoria', 'Intenzione Rilevata', 'Data e Ora Dichiarazione']
  ]

  summary.value.staff.forEach(s => {
    rows.push([
      `"${(s.last_name || '').replace(/"/g, '""')}"`,
      `"${(s.first_name || '').replace(/"/g, '""')}"`,
      `"${(s.email || '').replace(/"/g, '""')}"`,
      `"${s.role || ''}"`,
      `"${(s.role_display || '').replace(/"/g, '""')}"`,
      `"${getIntentionLabel(s.intention)}"`,
      `"${s.declared_at ? formatDateTime(s.declared_at) : 'Nessuna risposta'}"`
    ])
  })

  const csvContent = '\uFEFF' + rows.map(e => e.join(';')).join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  const filename = `rilevazione_sciopero_${summary.value.notice.strike_date || 'export'}.csv`
  link.setAttribute('href', url)
  link.setAttribute('download', filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  $q.notify({
    type: 'positive',
    message: 'Report CSV esportato con successo',
    caption: filename
  })
}

const getIntentionLabel = (intention) => {
  switch (intention) {
    case 'participates': return 'Aderisco'
    case 'not_participates': return 'Non aderisco'
    case 'undecided': return 'Non ho ancora deciso'
    default: return 'Nessuna risposta'
  }
}

const getIntentionIcon = (intention) => {
  switch (intention) {
    case 'participates': return 'check_circle'
    case 'not_participates': return 'cancel'
    case 'undecided': return 'help_outline'
    default: return 'hourglass_empty'
  }
}

const getIntentionBadgeColor = (intention) => {
  switch (intention) {
    case 'participates': return 'positive'
    case 'not_participates': return 'negative'
    case 'undecided': return 'amber-9'
    default: return 'grey-6'
  }
}

const getRoleBadgeColor = (role) => {
  switch (role) {
    case 'teacher':
    case 'coordinator':
      return 'primary'
    case 'dsga':
      return 'deep-purple-7'
    case 'assistente_amministrativo':
      return 'teal-7'
    case 'collaboratore_ds':
      return 'blue-8'
    case 'collaboratore_scolastico':
      return 'cyan-8'
    case 'principal':
      return 'purple-9'
    default:
      return 'indigo-7'
  }
}

const formatDate = (isoDate) => {
  if (!isoDate) return '-'
  try {
    const d = new Date(isoDate)
    return d.toLocaleDateString('it-IT', { day: '2-digit', month: 'long', year: 'numeric' })
  } catch {
    return isoDate
  }
}

const formatDateTime = (isoDate) => {
  if (!isoDate) return '-'
  try {
    const d = new Date(isoDate)
    return d.toLocaleDateString('it-IT', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return isoDate
  }
}

onMounted(() => {
  loadNotices()
})
</script>

<style scoped>
.notice-card {
  transition: all 0.2s ease-in-out;
}
.notice-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}
.notice-card-selected {
  border: 2px solid var(--q-primary) !important;
  background-color: #f8faff;
}
.kpi-card {
  transition: transform 0.2s ease;
}
.kpi-card:hover {
  transform: translateY(-2px);
}
.search-input {
  width: 200px;
}
.role-select {
  width: 190px;
}
.intention-select {
  width: 180px;
}
@media (max-width: 768px) {
  .search-input, .role-select, .intention-select {
    width: 100%;
  }
}
</style>
