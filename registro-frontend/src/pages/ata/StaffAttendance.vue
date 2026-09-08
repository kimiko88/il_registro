<template>
  <q-page class="q-pa-md q-pa-lg-xl staff-attendance-page">
    <!-- Header Hero Section -->
    <div class="row items-center justify-between q-mb-lg gap-md">
      <div class="col-12 col-md-7">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge color="indigo-7" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="badge" size="14px" class="q-mr-xs" />
            {{ t('staffAttendance.badgeStaffDocenti') || 'PERSONALE & DOCENTI' }}
          </q-badge>
          <q-badge v-if="isStrikeMode" color="negative" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="campaign" size="14px" class="q-mr-xs" />
            {{ t('staffAttendance.strikeModeActive') || 'RILEVAZIONE SCIOPERO ATTIVA' }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none flex items-center">
          <q-icon name="co_present" color="primary" class="q-mr-sm" size="36px" />
          {{ t('staffAttendance.title') || 'Riepilogo Presenze Giornaliere' }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          {{ t('staffAttendance.subtitle') || 'Monitoraggio presenze personale ATA e docenti • Gestione timbrature badge e rilevazione scioperi' }}
        </div>
      </div>

      <!-- Date & Action Controls -->
      <div class="col-12 col-md-5 row items-center justify-end q-gutter-sm">
        <q-btn-group outline rounded class="bg-white shadow-1">
          <q-btn flat dense icon="chevron_left" @click="changeDate(-1)" :title="t('staffAttendance.prevDay')">
            <q-tooltip>{{ t('staffAttendance.prevDay') || 'Giorno precedente' }}</q-tooltip>
          </q-btn>
          <q-btn flat no-caps class="text-weight-bold text-slate-700 q-px-sm">
            <q-icon name="event" color="primary" size="18px" class="q-mr-xs" />
            {{ formattedSelectedDate }}
            <q-popup-proxy cover transition-show="scale" transition-hide="scale">
              <q-date v-model="selectedDate" mask="YYYY-MM-DD" @update:model-value="loadData" today-btn />
            </q-popup-proxy>
          </q-btn>
          <q-btn flat dense icon="chevron_right" @click="changeDate(1)" :title="t('staffAttendance.nextDay')">
            <q-tooltip>{{ t('staffAttendance.nextDay') || 'Giorno successivo' }}</q-tooltip>
          </q-btn>
        </q-btn-group>

        <q-btn flat rounded color="primary" icon="today" :label="t('staffAttendance.today') || 'Oggi'" @click="goToToday" class="bg-indigo-50" />
        <q-btn
          color="secondary"
          icon="credit_card"
          :label="t('staffAttendance.badgeSwipeBtn') || 'Timbratura Badge'"
          no-caps
          rounded
          class="shadow-2 text-weight-bold"
          @click="openBadgeSwipeDialog"
        />
        <q-btn
          outline
          color="primary"
          icon="refresh"
          round
          dense
          @click="loadData"
          :loading="loading"
        >
          <q-tooltip>{{ t('staffAttendance.refreshTooltip') || 'Aggiorna dati' }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Alert / Banner Informativo Operatore -->
    <div class="q-mb-lg bg-indigo-50 border border-indigo-100 rounded-xl q-pa-md row items-center justify-between">
      <div class="row items-center gap-sm">
        <q-avatar size="36px" color="indigo-600" text-color="white" icon="verified_user" />
        <div>
          <div class="text-weight-bold text-slate-800">
            {{ t('staffAttendance.authorizedOperator') || 'Operatore Autorizzato' }}: <span class="text-primary">{{ user?.first_name }} {{ user?.last_name }}</span> ({{ roleLabel(userRole) }})
          </div>
          <div class="text-caption text-slate-600">
            {{ t('staffAttendance.authorizedDesc') || 'Autorizzato alla consultazione e registrazione presenze ATA e docenti (anche in occasione di scioperi).' }}
          </div>
        </div>
      </div>
      <div class="row items-center q-gutter-xs">
        <q-btn
          flat
          dense
          no-caps
          color="negative"
          icon="campaign"
          :label="isStrikeMode ? (t('staffAttendance.deactivateStrikeMode') || 'Disattiva Vista Sciopero') : (t('staffAttendance.activateStrikeMode') || 'Attiva Modalità Sciopero')"
          @click="isStrikeMode = !isStrikeMode"
          class="rounded-lg text-weight-bold"
        />
      </div>
    </div>

    <!-- Quick Stats Cards -->
    <div class="row q-col-gutter-md q-mb-xl">
      <!-- 1. Docenti Presenti -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="stat-card glass-card rounded-xl full-height shadow-1">
          <q-card-section class="row items-center no-wrap">
            <div class="bg-indigo-100 text-indigo-800 q-pa-md rounded-xl q-mr-md">
              <q-icon name="school" size="30px" />
            </div>
            <div>
              <div class="text-h4 text-weight-bolder text-slate-800">
                {{ summaryData?.teachers?.present ?? 0 }}
                <span class="text-caption text-slate-400 font-normal">/ {{ summaryData?.teachers?.total ?? 0 }}</span>
              </div>
              <div class="text-caption text-weight-bold text-slate-500 text-uppercase letter-spacing-1">
                {{ t('staffAttendance.teachersPresent') || 'Docenti Presenti' }}
              </div>
              <div class="text-caption text-indigo-700 text-weight-bold q-mt-xs">
                {{ t('staffAttendance.presenceRate', { rate: getPercentage(summaryData?.teachers?.present, summaryData?.teachers?.total) }) }}
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 2. Personale ATA Presente -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="stat-card glass-card rounded-xl full-height shadow-1">
          <q-card-section class="row items-center no-wrap">
            <div class="bg-cyan-100 text-cyan-800 q-pa-md rounded-xl q-mr-md">
              <q-icon name="badge" size="30px" />
            </div>
            <div>
              <div class="text-h4 text-weight-bolder text-slate-800">
                {{ summaryData?.ata?.present ?? 0 }}
                <span class="text-caption text-slate-400 font-normal">/ {{ summaryData?.ata?.total ?? 0 }}</span>
              </div>
              <div class="text-caption text-weight-bold text-slate-500 text-uppercase letter-spacing-1">
                {{ t('staffAttendance.ataPresent') || 'Personale ATA Presente' }}
              </div>
              <div class="text-caption text-cyan-800 text-weight-bold q-mt-xs">
                {{ t('staffAttendance.presenceRate', { rate: getPercentage(summaryData?.ata?.present, summaryData?.ata?.total) }) }}
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 3. Personale in Sciopero -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="stat-card glass-card rounded-xl full-height shadow-1" :class="{ 'border-red-400 bg-red-50': (summaryData?.total_strike ?? 0) > 0 }">
          <q-card-section class="row items-center no-wrap">
            <div class="bg-red-100 text-red-800 q-pa-md rounded-xl q-mr-md">
              <q-icon name="campaign" size="30px" />
            </div>
            <div>
              <div class="text-h4 text-weight-bolder" :class="(summaryData?.total_strike ?? 0) > 0 ? 'text-negative' : 'text-slate-800'">
                {{ summaryData?.total_strike ?? 0 }}
              </div>
              <div class="text-caption text-weight-bold text-slate-500 text-uppercase letter-spacing-1">
                {{ t('staffAttendance.onStrike') || 'In Sciopero' }}
              </div>
              <div class="text-caption text-red-700 text-weight-medium q-mt-xs">
                {{ t('staffAttendance.strikeSub', { teachers: summaryData?.teachers?.on_strike ?? 0, ata: summaryData?.ata?.on_strike ?? 0 }) }}
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 4. Assenti / Permesso / Malattia -->
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="stat-card glass-card rounded-xl full-height shadow-1">
          <q-card-section class="row items-center no-wrap">
            <div class="bg-amber-100 text-amber-800 q-pa-md rounded-xl q-mr-md">
              <q-icon name="event_busy" size="30px" />
            </div>
            <div>
              <div class="text-h4 text-weight-bolder text-slate-800">
                {{ (summaryData?.total_absent ?? 0) + (summaryData?.total_leave ?? 0) }}
              </div>
              <div class="text-caption text-weight-bold text-slate-500 text-uppercase letter-spacing-1">
                {{ t('staffAttendance.absentLeave') || 'Assenti / Congedo' }}
              </div>
              <div class="text-caption text-amber-900 text-weight-medium q-mt-xs">
                {{ t('staffAttendance.leaveSub', { leave: summaryData?.total_leave ?? 0 }) }}
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Category Breakdown Section -->
    <div class="q-mb-xl">
      <div class="row items-center justify-between q-mb-md">
        <div>
          <h2 class="text-h5 text-weight-bold text-slate-800 q-my-none flex items-center">
            <q-icon name="pie_chart" color="primary" class="q-mr-sm" size="24px" />
            {{ t('staffAttendance.breakdownTitle') || 'Dettaglio Presenze per Categoria ATA & Docenti' }}
          </h2>
          <div class="text-caption text-slate-500">
            {{ t('staffAttendance.breakdownSubtitle') || 'Conteggio analitico suddiviso per figura professionale della scuola' }}
          </div>
        </div>
      </div>

      <div class="row q-col-gutter-md">
        <div
          v-for="cat in roleSummaries"
          :key="cat.role"
          class="col-12 col-sm-6 col-lg"
        >
          <q-card class="rounded-xl border border-slate-200 shadow-sm full-height hover-lift">
            <q-card-section class="q-pb-xs">
              <div class="row items-center justify-between no-wrap q-mb-xs">
                <q-chip
                  dense
                  :color="getRoleBadgeColor(cat.role)"
                  text-color="white"
                  class="text-weight-bolder text-caption"
                >
                  {{ getRoleDisplay(cat.role, cat.role_display) }}
                </q-chip>
                <div class="text-caption text-weight-bold text-slate-400">
                  {{ t('staffAttendance.totalShort', { total: cat.total }) }}
                </div>
              </div>
            </q-card-section>

            <q-card-section class="q-pt-none q-pb-sm">
              <div class="row items-baseline q-gutter-x-sm">
                <div class="text-h4 text-weight-bolder text-positive">{{ cat.present }}</div>
                <div class="text-caption text-slate-500">{{ t('staffAttendance.presentCount') || 'presenti' }}</div>
              </div>

              <!-- Progress bar -->
              <q-linear-progress
                :value="cat.total > 0 ? (cat.present / cat.total) : 0"
                :color="cat.present === cat.total ? 'positive' : 'indigo-6'"
                rounded
                size="6px"
                class="q-my-sm"
              />

              <div class="row items-center justify-between text-caption text-slate-600 q-mt-xs">
                <span><q-icon name="cancel" color="negative" size="14px" /> {{ t('staffAttendance.absentCount', { count: cat.absent }) }}</span>
                <span v-if="cat.on_strike > 0" class="text-negative text-weight-bold">
                  <q-icon name="campaign" size="14px" /> {{ t('staffAttendance.strikeCount', { count: cat.on_strike }) }}
                </span>
                <span><q-icon name="local_hospital" color="amber-8" size="14px" /> {{ t('staffAttendance.leaveCount', { count: cat.on_leave }) }}</span>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Nominative Attendance Table -->
    <q-card class="rounded-xl border border-slate-200 shadow-1 q-mb-xl">
      <q-card-section class="row items-center justify-between gap-md border-b border-slate-100">
        <div>
          <div class="text-h6 text-weight-bold text-slate-800 flex items-center">
            <q-icon name="format_list_bulleted" color="primary" class="q-mr-sm" size="22px" />
            {{ t('staffAttendance.registerTitle') || 'Registro Presenze del Giorno' }}
          </div>
          <div class="text-caption text-slate-500">
            {{ t('staffAttendance.registerSubtitle') || 'Lista nominativa del personale con orari timbratura badge e stato presenza' }}
          </div>
        </div>

        <!-- Table Filters -->
        <div class="row items-center q-gutter-sm">
          <q-input
            v-model="filterSearch"
            dense
            outlined
            :placeholder="t('staffAttendance.searchPlaceholder') || 'Cerca nome o badge...'"
            class="bg-white"
            style="min-width: 220px;"
            clearable
          >
            <template v-slot:prepend>
              <q-icon name="search" size="18px" color="grey-6" />
            </template>
          </q-input>

          <q-select
            v-model="filterRole"
            :options="roleFilterOptions"
            dense
            outlined
            emit-value
            map-options
            class="bg-white"
            style="min-width: 170px;"
          />

          <q-select
            v-model="filterStatus"
            :options="statusFilterOptions"
            dense
            outlined
            emit-value
            map-options
            class="bg-white"
            style="min-width: 150px;"
          />
        </div>
      </q-card-section>

      <!-- Table Content -->
      <q-table
        :rows="filteredStaffList"
        :columns="columns"
        row-key="user_id"
        :loading="loading"
        flat
        :pagination="{ rowsPerPage: 15 }"
        class="staff-table"
        :no-data-label="t('staffAttendance.noData') || 'Nessun dato di presenza registrato per questa data'"
      >
        <!-- Nome & Cognome Slot -->
        <template v-slot:body-cell-person="props">
          <q-td :props="props">
            <div class="row items-center no-wrap">
              <q-avatar size="34px" :color="getRoleAvatarColor(props.row.role)" text-color="white" class="q-mr-sm text-weight-bold">
                {{ (props.row.first_name?.[0] || '') + (props.row.last_name?.[0] || '') }}
              </q-avatar>
              <div>
                <div class="text-weight-bold text-slate-800">
                  {{ props.row.last_name }} {{ props.row.first_name }}
                </div>
                <div class="text-caption text-slate-500">
                  {{ props.row.email }}
                </div>
              </div>
            </div>
          </q-td>
        </template>

        <!-- Categoria Slot -->
        <template v-slot:body-cell-role="props">
          <q-td :props="props">
            <q-chip
              dense
              size="sm"
              :color="getRoleBadgeColor(props.row.role)"
              text-color="white"
              class="text-weight-bold text-caption"
            >
              {{ getRoleDisplay(props.row.role, props.row.role_display) }}
            </q-chip>
          </q-td>
        </template>

        <!-- Stato Presenza Slot -->
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-chip
              dense
              :color="getStatusColor(props.row.status)"
              text-color="white"
              class="text-weight-bold text-caption q-px-sm"
              :icon="getStatusIcon(props.row.status)"
            >
              {{ getStatusLabel(props.row.status) }}
            </q-chip>
            <div v-if="props.row.is_strike_recorded" class="text-caption text-negative text-weight-bold q-mt-xs">
              {{ t('staffAttendance.strikeDetected') || 'Sciopero rilevato' }}
            </div>
          </q-td>
        </template>

        <!-- Timbratura Badge Slot -->
        <template v-slot:body-cell-badge="props">
          <q-td :props="props">
            <div v-if="props.row.badge_entry_time" class="text-slate-800 text-caption">
              <div class="flex items-center">
                <q-icon name="login" color="positive" size="14px" class="q-mr-xs" />
                {{ t('staffAttendance.entryTime', { time: formatTime(props.row.badge_entry_time) }) }}
              </div>
              <div v-if="props.row.badge_exit_time" class="flex items-center q-mt-xs">
                <q-icon name="logout" color="negative" size="14px" class="q-mr-xs" />
                {{ t('staffAttendance.exitTime', { time: formatTime(props.row.badge_exit_time) }) }}
              </div>
            </div>
            <div v-else class="text-caption text-slate-400 flex items-center">
              <q-icon name="credit_card_off" size="14px" class="q-mr-xs" />
              {{ t('staffAttendance.noBadgeSwipe') || 'Nessuna timbratura' }}
            </div>
          </q-td>
        </template>

        <!-- Note Slot -->
        <template v-slot:body-cell-notes="props">
          <q-td :props="props">
            <span class="text-caption text-slate-600">{{ props.row.notes || '-' }}</span>
          </q-td>
        </template>

        <!-- Azioni Slot -->
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" align="right">
            <div class="row items-center justify-end q-gutter-xs">
              <q-btn
                flat
                round
                dense
                color="primary"
                icon="edit"
                @click="openEditDialog(props.row)"
                :title="t('staffAttendance.editTooltip')"
              >
                <q-tooltip>{{ t('staffAttendance.editTooltip') || 'Modifica stato presenza' }}</q-tooltip>
              </q-btn>
              <q-btn
                flat
                round
                dense
                color="secondary"
                icon="credit_card"
                @click="openBadgeDialogForUser(props.row)"
                :title="t('staffAttendance.simulateBadgeTooltip')"
              >
                <q-tooltip>{{ t('staffAttendance.simulateBadgeTooltip') || 'Simula timbratura badge per questa persona' }}</q-tooltip>
              </q-btn>
            </div>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog: Modifica / Registra Presenza -->
    <q-dialog v-model="showEditDialog" persistent>
      <q-card style="min-width: 440px; max-width: 520px" class="rounded-xl">
        <q-card-section class="row items-center justify-between bg-primary text-white">
          <div class="text-h6 text-weight-bold flex items-center">
            <q-icon name="edit_note" size="24px" class="q-mr-sm" />
            {{ t('staffAttendance.editDialogTitle') || 'Registra / Modifica Presenza' }}
          </div>
          <q-btn flat round dense icon="close" color="white" v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md" v-if="editingRecord">
          <div class="bg-slate-50 p-3 rounded-lg border border-slate-200">
            <div class="text-weight-bold text-slate-800">
              {{ editingRecord.last_name }} {{ editingRecord.first_name }}
            </div>
            <div class="text-caption text-slate-500">
              {{ t('staffAttendance.roleLabel') || 'Ruolo' }}: <q-badge :color="getRoleBadgeColor(editingRecord.role)">{{ getRoleDisplay(editingRecord.role, editingRecord.role_display) }}</q-badge>
              • {{ t('staffAttendance.dateLabel') || 'Data' }}: <b>{{ selectedDate }}</b>
            </div>
          </div>

          <!-- Stato Presenza -->
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('staffAttendance.statusSelectLabel') || 'Stato di Presenza *' }}</label>
            <q-select
              v-model="editForm.status"
              :options="statusOptions"
              emit-value
              map-options
              outlined
              dense
            />
          </div>

          <!-- Se Sciopero -->
          <div v-if="editForm.status === 'on_strike'" class="bg-red-50 border border-red-200 p-3 rounded-lg">
            <div class="text-caption text-negative text-weight-bold q-mb-xs">
              <q-icon name="warning" class="q-mr-xs" />
              {{ t('staffAttendance.strikeSectionTitle') || 'Rilevazione Sciopero Docente / Personale ATA' }}
            </div>
            <q-input
              v-model="editForm.strike_code"
              outlined
              dense
              :label="t('staffAttendance.strikeCodeLabel') || 'Codice o Denominazione Sciopero'"
              :placeholder="t('staffAttendance.strikeCodePlaceholder') || 'es. Sciopero Generale Comparto Scuola'"
              class="q-mt-xs"
            />
          </div>

          <!-- Note aggiuntive -->
          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('staffAttendance.notesLabel') || 'Note o Giustificativo' }}</label>
            <q-input
              v-model="editForm.notes"
              type="textarea"
              rows="2"
              outlined
              dense
              :placeholder="t('staffAttendance.notesPlaceholder') || 'es. Uscita anticipata per visita medica, etc.'"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('staffAttendance.cancelBtn') || 'Annulla'" color="grey-7" v-close-popup />
          <q-btn
            color="primary"
            :label="t('staffAttendance.saveBtn') || 'Salva Presenza'"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="saving"
            @click="saveAttendance"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog: Timbratura Badge (Hardware / Simulatore) -->
    <q-dialog v-model="showBadgeSwipeDialog">
      <q-card style="min-width: 400px; max-width: 480px" class="rounded-xl">
        <q-card-section class="row items-center justify-between bg-secondary text-white">
          <div class="text-h6 text-weight-bold flex items-center">
            <q-icon name="credit_card" size="24px" class="q-mr-sm" />
            {{ t('staffAttendance.badgeDialogTitle') || 'Timbratrice Badge Personale ATA' }}
          </div>
          <q-btn flat round dense icon="close" color="white" v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <div class="text-caption text-slate-600">
            {{ t('staffAttendance.badgeDialogSubtitle') || 'Questa interfaccia riceve o simula la lettura del badge RFID / NFC o codice a barre da terminale di ingresso della scuola.' }}
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('staffAttendance.badgeCodeLabel') || 'Codice Badge *' }}</label>
            <q-input
              v-model="badgeForm.badge_code"
              outlined
              dense
              autofocus
              placeholder="es. BDG-ATA-001"
              @keyup.enter="submitBadgeSwipe"
            >
              <template v-slot:prepend>
                <q-icon name="qr_code" />
              </template>
            </q-input>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('staffAttendance.swipeTypeLabel') || 'Tipo Timbratura' }}</label>
            <q-select
              v-model="badgeForm.swipe_type"
              :options="[
                { label: t('staffAttendance.swipeIn') || 'Entrata (Ingresso)', value: 'in' },
                { label: t('staffAttendance.swipeOut') || 'Uscita (Fine Servizio)', value: 'out' },
                { label: t('staffAttendance.swipeBreakOut') || 'Uscita Pausa / Servizio Esterno', value: 'break_out' },
                { label: t('staffAttendance.swipeBreakIn') || 'Rientro Pausa', value: 'break_in' }
              ]"
              emit-value
              map-options
              outlined
              dense
            />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('staffAttendance.terminalIdLabel') || 'Identificativo Terminale' }}</label>
            <q-input
              v-model="badgeForm.device_id"
              outlined
              dense
              placeholder="TERMINAL-INGRESSO-01"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('staffAttendance.cancelBtn') || 'Annulla'" color="grey-7" v-close-popup />
          <q-btn
            color="secondary"
            :label="t('staffAttendance.registerSwipeBtn') || 'Registra Timbratura'"
            no-caps
            rounded
            class="q-px-md text-weight-bold shadow-1"
            :loading="savingBadge"
            @click="submitBadgeSwipe"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import staffAttendanceService from '@/services/staffAttendanceService'

const { t, te, locale } = useI18n()
const $q = useQuasar()
const authStore = useAuthStore()
const { user, userRole } = storeToRefs(authStore)

const selectedDate = ref(new Date().toISOString().substring(0, 10))
const isStrikeMode = ref(false)
const loading = ref(false)
const saving = ref(false)
const savingBadge = ref(false)

const summaryData = ref(null)
const staffList = ref([])

const filterSearch = ref('')
const filterRole = ref('all')
const filterStatus = ref('all')

const showEditDialog = ref(false)
const editingRecord = ref(null)
const editForm = ref({
  status: 'present',
  notes: '',
  strike_code: ''
})

const showBadgeSwipeDialog = ref(false)
const badgeForm = ref({
  badge_code: '',
  swipe_type: 'in',
  device_id: 'TERMINAL-INGRESSO-01'
})

const columns = computed(() => [
  { name: 'person', label: t('staffAttendance.colStaff') || 'Personale', align: 'left', field: 'last_name', sortable: true },
  { name: 'role', label: t('staffAttendance.colCategory') || 'Categoria', align: 'left', field: 'role', sortable: true },
  { name: 'status', label: t('staffAttendance.colStatus') || 'Stato Presenza', align: 'center', field: 'status', sortable: true },
  { name: 'badge', label: t('staffAttendance.colBadge') || 'Timbratura Badge', align: 'left', field: 'badge_entry_time' },
  { name: 'notes', label: t('staffAttendance.colNotes') || 'Note / Sciopero', align: 'left', field: 'notes' },
  { name: 'actions', label: t('staffAttendance.colActions') || 'Azioni', align: 'right' }
])

const roleFilterOptions = computed(() => [
  { label: t('staffAttendance.allCategories') || 'Tutte le categorie', value: 'all' },
  { label: t('roles.teacher') || 'Docenti', value: 'teacher' },
  { label: t('roles.dsga') || 'DSGA', value: 'dsga' },
  { label: t('roles.assistente_amministrativo') || 'Assistenti Amministrativi', value: 'assistente_amministrativo' },
  { label: t('roles.collaboratore_ds') || 'Collaboratori D.S.', value: 'collaboratore_ds' },
  { label: t('roles.collaboratore_scolastico') || 'Collaboratori Scolastici', value: 'collaboratore_scolastico' },
  { label: t('roles.principal') || 'Dirigente Scolastico', value: 'principal' }
])

const statusFilterOptions = computed(() => [
  { label: t('staffAttendance.allStatuses') || 'Tutti gli stati', value: 'all' },
  { label: t('staffAttendance.statusPresent') || 'Presenti', value: 'present' },
  { label: t('staffAttendance.statusAbsent') || 'Assenti', value: 'absent' },
  { label: t('staffAttendance.statusOnStrike') || 'In Sciopero', value: 'on_strike' },
  { label: t('staffAttendance.statusLate') || 'In Ritardo', value: 'late' },
  { label: t('staffAttendance.absentLeave') || 'In Permesso/Malattia', value: 'leave' }
])

const statusOptions = computed(() => [
  { label: t('staffAttendance.statusPresent') || 'Presente', value: 'present' },
  { label: t('staffAttendance.statusOnStrike') || 'In Sciopero', value: 'on_strike' },
  { label: t('staffAttendance.statusAbsent') || 'Assente', value: 'absent' },
  { label: t('staffAttendance.statusLate') || 'In Ritardo', value: 'late' },
  { label: t('staffAttendance.statusSickLeave') || 'Malattia', value: 'sick_leave' },
  { label: t('staffAttendance.statusPermit') || 'Permesso Retribuito/Personale', value: 'permit' },
  { label: t('staffAttendance.statusMission') || 'Missione / Servizio Esterno', value: 'mission' }
])

const formattedSelectedDate = computed(() => {
  if (!selectedDate.value) return ''
  const [year, month, day] = selectedDate.value.split('-')
  return `${day}/${month}/${year}`
})

const roleSummaries = computed(() => {
  if (!summaryData.value?.by_role) return []
  return summaryData.value.by_role
})

const filteredStaffList = computed(() => {
  let list = staffList.value || []

  if (filterRole.value !== 'all') {
    list = list.filter(item => item.role === filterRole.value)
  }

  if (filterStatus.value !== 'all') {
    if (filterStatus.value === 'leave') {
      list = list.filter(item => item.status === 'sick_leave' || item.status === 'permit' || item.status === 'mission')
    } else {
      list = list.filter(item => item.status === filterStatus.value)
    }
  }

  if (filterSearch.value) {
    const q = filterSearch.value.toLowerCase().trim()
    list = list.filter(item => {
      const full = `${item.first_name || ''} ${item.last_name || ''} ${item.email || ''}`.toLowerCase()
      return full.includes(q)
    })
  }

  return list
})

function changeDate(days) {
  const d = new Date(selectedDate.value)
  d.setDate(d.getDate() + days)
  selectedDate.value = d.toISOString().substring(0, 10)
  loadData()
}

function goToToday() {
  selectedDate.value = new Date().toISOString().substring(0, 10)
  loadData()
}

function getPercentage(num, total) {
  if (!total || total === 0) return 0
  return Math.round((num / total) * 100)
}

function roleLabel(role) {
  if (role && te('roles.' + role)) {
    return t('roles.' + role)
  }
  const map = {
    superadmin: 'Super Admin',
    admin: 'Amministratore',
    principal: 'Dirigente Scolastico',
    vice_principal: 'Collaboratore Vicario',
    dsga: 'DSGA (Direttore SGA)',
    collaboratore_ds: 'Collaboratore D.S.',
    assistente_amministrativo: 'Assistente Amministrativo',
    collaboratore_scolastico: 'Collaboratore Scolastico',
    secretary: 'Segreteria',
    teacher: 'Docente'
  }
  return map[role] || role
}

function getRoleDisplay(role, fallback) {
  if (role && te('roles.' + role)) {
    return t('roles.' + role)
  }
  return fallback || roleLabel(role)
}

function getRoleBadgeColor(role) {
  const map = {
    principal: 'purple-8',
    vice_principal: 'indigo-8',
    dsga: 'teal-8',
    collaboratore_ds: 'deep-orange-7',
    assistente_amministrativo: 'cyan-8',
    collaboratore_scolastico: 'amber-9',
    teacher: 'indigo-6'
  }
  return map[role] || 'grey-7'
}

function getRoleAvatarColor(role) {
  const map = {
    principal: 'purple-7',
    vice_principal: 'indigo-7',
    dsga: 'teal-7',
    collaboratore_ds: 'deep-orange-6',
    assistente_amministrativo: 'cyan-7',
    collaboratore_scolastico: 'amber-8',
    teacher: 'indigo-5'
  }
  return map[role] || 'grey-6'
}

function getStatusColor(status) {
  const map = {
    present: 'positive',
    absent: 'negative',
    late: 'warning',
    mission: 'info',
    permit: 'amber-8',
    sick_leave: 'deep-purple-6',
    on_strike: 'negative'
  }
  return map[status] || 'grey-6'
}

function getStatusIcon(status) {
  const map = {
    present: 'check_circle',
    absent: 'cancel',
    late: 'schedule',
    mission: 'flight_takeoff',
    permit: 'assignment',
    sick_leave: 'local_hospital',
    on_strike: 'campaign'
  }
  return map[status] || 'help'
}

function getStatusLabel(status) {
  const keyMap = {
    present: 'statusPresent',
    absent: 'statusAbsent',
    late: 'statusLate',
    mission: 'statusMission',
    permit: 'statusPermit',
    sick_leave: 'statusSickLeave',
    on_strike: 'statusOnStrike'
  }
  const key = keyMap[status]
  if (key && te('staffAttendance.' + key)) {
    return t('staffAttendance.' + key)
  }
  const map = {
    present: 'Presente',
    absent: 'Assente',
    late: 'In Ritardo',
    mission: 'Missione',
    permit: 'Permesso',
    sick_leave: 'Malattia',
    on_strike: 'In Sciopero'
  }
  return map[status] || status
}

function formatTime(isoStr) {
  if (!isoStr) return ''
  try {
    const d = new Date(isoStr)
    return d.toLocaleTimeString(locale.value || 'it-IT', { hour: '2-digit', minute: '2-digit' })
  } catch {
    return isoStr
  }
}

async function loadData() {
  loading.value = true
  try {
    const [summary, listRes] = await Promise.all([
      staffAttendanceService.getDailySummary(selectedDate.value).catch(err => {
        console.error('Error fetching summary', err)
        return null
      }),
      staffAttendanceService.getList(selectedDate.value).catch(err => {
        console.error('Error fetching list', err)
        return { data: [] }
      })
    ])

    summaryData.value = summary
    staffList.value = listRes?.data || []
  } catch (err) {
    console.error('Failed to load staff attendance data', err)
    $q.notify({ type: 'negative', message: t('staffAttendance.saveError') || 'Errore nel caricamento delle presenze' })
  } finally {
    loading.value = false
  }
}

function openEditDialog(record) {
  editingRecord.value = record
  editForm.value = {
    status: record.status || 'present',
    notes: record.notes || '',
    strike_code: record.strike_code || ''
  }
  showEditDialog.value = true
}

async function saveAttendance() {
  if (!editingRecord.value) return
  saving.value = true
  try {
    await staffAttendanceService.recordAttendance({
      user_id: editingRecord.value.user_id,
      date: selectedDate.value,
      status: editForm.value.status,
      notes: editForm.value.notes,
      strike_code: editForm.value.strike_code,
      is_strike_recorded: editForm.value.status === 'on_strike'
    })

    $q.notify({ type: 'positive', message: t('staffAttendance.saveSuccess') || 'Presenza aggiornata con successo' })
    showEditDialog.value = false
    await loadData()
  } catch (err) {
    console.error('Error saving attendance', err)
    $q.notify({ type: 'negative', message: err?.response?.data?.error || t('staffAttendance.saveError') || 'Errore nel salvataggio' })
  } finally {
    saving.value = false
  }
}

function openBadgeSwipeDialog() {
  badgeForm.value = {
    badge_code: '',
    swipe_type: 'in',
    device_id: 'TERMINAL-INGRESSO-01'
  }
  showBadgeSwipeDialog.value = true
}

function openBadgeDialogForUser(record) {
  badgeForm.value = {
    badge_code: `BDG-${record.role.substring(0, 3).toUpperCase()}-${record.user_id.substring(0, 4)}`,
    swipe_type: record.badge_entry_time ? 'out' : 'in',
    device_id: 'TERMINAL-INGRESSO-01'
  }
  showBadgeSwipeDialog.value = true
}

async function submitBadgeSwipe() {
  if (!badgeForm.value.badge_code) {
    $q.notify({ type: 'warning', message: t('staffAttendance.badgeCodeLabel') || 'Inserisci il codice badge' })
    return
  }

  savingBadge.value = true
  try {
    await staffAttendanceService.registerBadgeSwipe(badgeForm.value)
    await staffAttendanceService.processBadgeSwipes().catch(() => {})

    $q.notify({
      type: 'positive',
      icon: 'check_circle',
      message: t('staffAttendance.badgeSwipeSuccess', { code: badgeForm.value.badge_code })
    })
    showBadgeSwipeDialog.value = false
    await loadData()
  } catch (err) {
    console.error('Error recording swipe', err)
    $q.notify({ type: 'negative', message: err?.response?.data?.error || t('staffAttendance.badgeSwipeError') || 'Errore timbratura badge' })
  } finally {
    savingBadge.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.glass-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(226, 232, 240, 0.8);
}
.stat-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.08);
}
.hover-lift {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.hover-lift:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(0, 0, 0, 0.07);
}
.letter-spacing-1 {
  letter-spacing: 0.05em;
}
</style>
