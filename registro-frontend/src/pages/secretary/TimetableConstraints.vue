<template>
  <q-page padding class="min-h-screen" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-slate-50 text-slate-800'">
    <!-- Top Header -->
    <div class="row items-center justify-between q-mb-lg gap-4">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none row items-center gap-2" :class="$q.dark.isActive ? 'text-white' : 'text-slate-800'">
          <q-icon name="tune" color="primary" size="36px" />
          {{ t('timetableConstraints.title') }}
        </h1>
        <p class="text-subtitle1 q-mt-xs q-mb-none" :class="$q.dark.isActive ? 'text-grey-4' : 'text-slate-500'">
          {{ t('timetableConstraints.subtitle') }}
        </p>
      </div>

      <div class="row items-center gap-2">
        <q-btn
          v-if="currentTab === 'rooms'"
          :label="t('timetableConstraints.newConstraintBtn')"
          icon="add"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openRoomReqModal"
        />
        <q-btn
          v-else-if="currentTab === 'groups'"
          :label="t('timetableConstraints.newGroupBtn')"
          icon="add"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openGroupModal"
        />
        <q-btn
          v-else-if="currentTab === 'spezzoni'"
          :label="t('timetableConstraints.newSpezzoneBtn')"
          icon="add"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openSpezzoneModal"
        />
        <q-btn
          v-else-if="currentTab === 'curriculum'"
          label="Eredita da Anno Precedente"
          icon="history_edu"
          color="indigo-7"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          @click="openInheritModal"
        />
        <q-btn
          v-else-if="currentTab === 'desiderata' && desiderataViewMode === 'quick_table'"
          label="Salva Tutti i Desiderata Rapidi"
          icon="save"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-md font-bold shadow-sm"
          :loading="savingAllQuickPrefs"
          @click="saveAllQuickPreferences"
        />
      </div>
    </div>

    <!-- Policy Banner: Lab Constraints Only, Desiderata Excluded -->
    <!-- Policy Banner: Lab Constraints & Desiderata -->
    <q-banner rounded class="bg-indigo-50 border border-indigo-200 text-indigo-950 q-mb-lg rounded-xl shadow-xs">
      <template v-slot:avatar>
        <q-icon name="info" color="indigo-7" size="28px" />
      </template>
      <div class="text-body2 font-medium">
        <strong>Regole di Generazione Orario & Desiderata:</strong> I vincoli di laboratorio con relative ore, la copertura cattedre e le compresenze per gruppi articolati sono vincoli prioritari. I desiderata personali dei docenti possono essere inseriti autonomamente dai docenti quando la finestra è aperta, oppure visualizzati e modificati direttamente in questa schermata.
      </div>
    </q-banner>

    <!-- Navigation Tabs -->
    <q-tabs
      v-model="currentTab"
      dense
      align="left"
      active-color="primary"
      indicator-color="primary"
      class="text-grey-7 q-mb-lg bg-white rounded-xl shadow-sm overflow-hidden"
      :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'"
    >
      <q-tab name="rooms" icon="meeting_room" :label="t('timetableConstraints.tabRooms')" class="q-py-md text-weight-bold" />
      <q-tab name="groups" icon="groups" :label="t('timetableConstraints.tabGroups')" class="q-py-md text-weight-bold" />
      <q-tab name="spezzoni" icon="hourglass_empty" :label="t('timetableConstraints.tabSpezzoni')" class="q-py-md text-weight-bold" />
      <q-tab name="desiderata" icon="event_available" label="Desiderata Docenti" class="q-py-md text-weight-bold" />
      <q-tab name="curriculum" icon="menu_book" label="Ore Materie & Limiti Classi" class="q-py-md text-weight-bold" />
    </q-tabs>

    <!-- TAB 1: Subject Room Requirements (Aule e Laboratori con scelta ore) -->
    <div v-show="currentTab === 'rooms'">
      <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <q-card-section class="row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold row items-center gap-2">
              <q-icon name="meeting_room" color="indigo-7" />
              {{ t('timetableConstraints.sectionTitle') }}
            </div>
            <div class="text-caption text-grey-6">
              {{ t('timetableConstraints.sectionSubtitle') }}
            </div>
          </div>
          <q-btn flat dense icon="refresh" color="primary" no-caps :label="t('timetableConstraints.refresh')" @click="fetchRoomReqs" />
        </q-card-section>

        <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

        <div v-if="loadingReqs" class="row justify-center q-pa-lg">
          <q-spinner-dots size="40px" color="primary" />
        </div>

        <div v-else-if="roomReqs.length === 0" class="text-center q-pa-xl">
          <q-icon name="science" size="48px" color="grey-5" />
          <div class="text-h6 text-grey-6 q-mt-sm">{{ t('timetableConstraints.noConstraintsTitle') }}</div>
          <p class="text-grey-5">{{ t('timetableConstraints.noConstraintsDesc') }}</p>
          <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.addConstraintBtn')" icon="add" class="rounded-xl q-mt-sm font-bold" @click="openRoomReqModal" />
        </div>

        <q-table
          v-else
          flat
          :rows="roomReqs"
          :columns="roomReqColumns"
          row-key="id"
          class="no-shadow"
        >
          <template v-slot:body-cell-required_room_type="props">
            <q-td :props="props">
              <q-badge color="indigo-7" class="q-px-sm q-py-xs font-bold">
                {{ getRoomTypeLabel(props.row.required_room_type) }}
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-lab_hours="props">
            <q-td :props="props" align="center">
              <q-badge color="purple-8" class="q-px-sm q-py-xs font-bold">
                {{ props.row.lab_hours || 1 }} ore / sett.
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-is_mandatory="props">
            <q-td :props="props" align="center">
              <q-badge :color="props.row.is_mandatory ? 'negative' : 'warning'">
                {{ props.row.is_mandatory ? t('timetableConstraints.mandatoryBadge') : t('timetableConstraints.preferentialBadge') }}
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-actions="props">
            <q-td :props="props" align="right">
              <q-btn flat round dense color="negative" icon="delete" @click="confirmDeleteReq(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card>
    </div>

    <!-- TAB 2: Gruppi Linguistici & Classi Articolate (Cattedre Congiunte) -->
    <div v-show="currentTab === 'groups'">
      <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <q-card-section class="row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold row items-center gap-2">
              <q-icon name="groups" color="indigo-7" />
              {{ t('timetableConstraints.groupsSectionTitle') }}
            </div>
            <div class="text-caption text-grey-6">
              {{ t('timetableConstraints.groupsSectionSubtitle') }}
            </div>
          </div>
          <div class="row items-center gap-2">
            <q-btn flat dense icon="refresh" color="primary" no-caps :label="t('timetableConstraints.refresh')" @click="fetchConstraints" />
            <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.newGroupBtn')" icon="add" class="rounded-xl font-bold" @click="openGroupModal" />
          </div>
        </q-card-section>

        <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

        <div v-if="loadingConstraints" class="row justify-center q-pa-lg">
          <q-spinner-dots size="40px" color="primary" />
        </div>

        <div v-else-if="associatedGroupsList.length === 0" class="text-center q-pa-xl">
          <q-icon name="groups" size="48px" color="grey-5" />
          <div class="text-h6 text-grey-6 q-mt-sm">{{ t('timetableConstraints.noGroupsTitle') }}</div>
          <p class="text-grey-5">{{ t('timetableConstraints.noGroupsDesc') }}</p>
          <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.newGroupBtn')" icon="add" class="rounded-xl q-mt-sm font-bold" @click="openGroupModal" />
        </div>

        <q-table
          v-else
          flat
          :rows="associatedGroupsList"
          :columns="groupColumns"
          row-key="id"
          class="no-shadow"
        >
          <template v-slot:body-cell-name="props">
            <q-td :props="props">
              <span class="text-weight-bold text-slate-800">{{ props.row.name }}</span>
            </q-td>
          </template>

          <template v-slot:body-cell-subject_name="props">
            <q-td :props="props">
              <q-badge color="indigo-8" class="q-px-sm q-py-xs font-bold">
                {{ props.row.subject_name }}
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-teacher_name="props">
            <q-td :props="props">
              <q-chip size="sm" icon="person" color="blue-1" text-color="blue-9" class="font-bold">
                {{ props.row.teacher_name }}
              </q-chip>
            </q-td>
          </template>

          <template v-slot:body-cell-classes="props">
            <q-td :props="props">
              <div class="row gap-1">
                <q-badge v-for="cName in props.row.class_names" :key="cName" color="teal-8" class="q-px-xs">
                  {{ cName }}
                </q-badge>
              </div>
            </q-td>
          </template>

          <template v-slot:body-cell-hours_per_week="props">
            <q-td :props="props" align="center">
              <q-badge color="purple-9" class="q-px-sm q-py-xs text-weight-bold">
                {{ props.row.hours_per_week }}h / sett.
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-actions="props">
            <q-td :props="props" align="right">
              <q-btn flat round dense color="negative" icon="delete" @click="confirmDeleteGroup(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card>
    </div>

    <!-- TAB 3: Spezzoni Orari & Cattedre Vacanti per Classe -->
    <div v-show="currentTab === 'spezzoni'">
      <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <q-card-section class="row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold row items-center gap-2">
              <q-icon name="hourglass_empty" color="amber-8" />
              {{ t('timetableConstraints.spezzoniSectionTitle') }}
            </div>
            <div class="text-caption text-grey-6">
              {{ t('timetableConstraints.spezzoniSectionSubtitle') }}
            </div>
          </div>
          <div class="row items-center gap-2">
            <q-btn flat dense icon="refresh" color="primary" no-caps :label="t('timetableConstraints.refresh')" @click="fetchConstraints" />
            <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.newSpezzoneBtn')" icon="add" class="rounded-xl font-bold" @click="openSpezzoneModal" />
          </div>
        </q-card-section>

        <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

        <div v-if="loadingConstraints" class="row justify-center q-pa-lg">
          <q-spinner-dots size="40px" color="primary" />
        </div>

        <div v-else-if="spezzoniList.length === 0" class="text-center q-pa-xl">
          <q-icon name="work_outline" size="48px" color="grey-5" />
          <div class="text-h6 text-grey-6 q-mt-sm">{{ t('timetableConstraints.noSpezzoniTitle') }}</div>
          <p class="text-grey-5">{{ t('timetableConstraints.noSpezzoniDesc') }}</p>
          <q-btn color="primary" unelevated no-caps :label="t('timetableConstraints.newSpezzoneBtn')" icon="add" class="rounded-xl q-mt-sm font-bold" @click="openSpezzoneModal" />
        </div>

        <q-table
          v-else
          flat
          :rows="spezzoniList"
          :columns="spezzoniColumns"
          row-key="id"
          class="no-shadow"
        >
          <template v-slot:body-cell-class_name="props">
            <q-td :props="props">
              <q-badge color="indigo-8" class="q-px-sm q-py-xs font-bold text-sm">
                {{ props.row.class_name }}
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-subject_name="props">
            <q-td :props="props">
              <span class="text-weight-bold text-slate-800">{{ props.row.subject_name }}</span>
            </q-td>
          </template>

          <template v-slot:body-cell-hours_per_week="props">
            <q-td :props="props" align="center">
              <q-badge color="amber-9" class="q-px-sm q-py-xs text-weight-bold">
                {{ props.row.hours_per_week }}h / sett.
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-placeholder_name="props">
            <q-td :props="props">
              <q-chip size="sm" icon="person_outline" color="blue-1" text-color="blue-9" class="font-bold">
                {{ props.row.placeholder_name }}
              </q-chip>
            </q-td>
          </template>

          <template v-slot:body-cell-notes="props">
            <q-td :props="props">
              <span class="text-caption text-grey-7">{{ props.row.notes || '-' }}</span>
            </q-td>
          </template>

          <template v-slot:body-cell-actions="props">
            <q-td :props="props" align="right">
              <q-btn flat round dense color="negative" icon="delete" @click="confirmDeleteSpezzone(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card>
    </div>

    <!-- TAB 4: Desiderata & Preferenze Docenti -->
    <div v-show="currentTab === 'desiderata'">
      <!-- Card: Finestra Temporale Inserimento Desiderata per i Docenti -->
      <q-card flat bordered class="rounded-2xl shadow-sm q-mb-lg" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <q-card-section class="row items-center justify-between gap-4">
          <div class="row items-center gap-3">
            <q-avatar size="44px" :color="isDesiderataWindowOpen ? 'green-1' : 'amber-1'" :text-color="isDesiderataWindowOpen ? 'positive' : 'amber-9'" icon="schedule" />
            <div>
              <div class="text-subtitle1 text-weight-bold">
                Finestra Inserimento Desiderata Docenti:
                <q-badge :color="isDesiderataWindowOpen ? 'positive' : 'grey-7'" class="q-ml-sm q-px-sm q-py-xs font-bold text-caption">
                  {{ isDesiderataWindowOpen ? 'ATTIVA (Aperta)' : 'DISATTIVATA (Chiusa)' }}
                </q-badge>
              </div>
              <div class="text-caption text-grey-6">
                {{ isDesiderataWindowOpen
                  ? 'I docenti possono accedere e compilare autonomamente i propri desiderata orario.'
                  : 'I docenti possono solo visualizzare in sola lettura le preferenze. Attiva la finestra per consentire modifiche.' }}
              </div>
            </div>
          </div>

          <div class="row items-center gap-2">
            <q-toggle
              v-model="isDesiderataWindowOpen"
              color="primary"
              :label="isDesiderataWindowOpen ? 'Finestra Aperta' : 'Finestra Chiusa'"
              :loading="loadingWindow"
              @update:model-value="toggleDesiderataWindow"
            />
          </div>
        </q-card-section>
      </q-card>

      <!-- Sub-tabs: Tabella Rapida Tutti i Docenti vs Matrice Singolo Docente -->
      <div class="row items-center justify-between gap-3 q-mb-md">
        <div class="row items-center gap-2">
          <q-btn
            :flat="desiderataViewMode !== 'quick_table'"
            :unelevated="desiderataViewMode === 'quick_table'"
            :color="desiderataViewMode === 'quick_table' ? 'primary' : 'grey-8'"
            icon="table_chart"
            label="Rappresentazione Tabellare (Tutti i Docenti)"
            no-caps
            class="rounded-xl font-bold q-px-md"
            @click="desiderataViewMode = 'quick_table'"
          />
          <q-btn
            :flat="desiderataViewMode !== 'single_matrix'"
            :unelevated="desiderataViewMode === 'single_matrix'"
            :color="desiderataViewMode === 'single_matrix' ? 'primary' : 'grey-8'"
            icon="calendar_view_week"
            label="Matrice Oraria Dettagliata (Singolo Docente)"
            no-caps
            class="rounded-xl font-bold q-px-md"
            @click="desiderataViewMode = 'single_matrix'"
          />
        </div>

        <div v-if="desiderataViewMode === 'quick_table'" class="row items-center gap-2">
          <q-btn
            flat
            dense
            icon="refresh"
            color="grey-7"
            label="Ricarica"
            no-caps
            class="rounded-lg q-px-sm"
            @click="fetchTeachersQuickPreferences"
          />
          <q-btn
            unelevated
            color="primary"
            icon="save"
            label="Salva Tutti i Desiderata Rapidi"
            no-caps
            class="rounded-xl font-bold q-px-md shadow-sm"
            :loading="savingAllQuickPrefs"
            @click="saveAllQuickPreferences"
          />
        </div>
      </div>

      <!-- VIEW 1: Rappresentazione Tabellare Rapida per Tutti i Docenti -->
      <div v-if="desiderataViewMode === 'quick_table'">
        <!-- Day-Off Distribution Analytics Card -->
        <q-card flat bordered class="rounded-2xl shadow-sm q-mb-md" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
          <q-card-section class="q-py-sm">
            <div class="row items-center justify-between gap-2 q-mb-xs">
              <div class="text-subtitle2 font-bold text-grey-8 row items-center gap-1" :class="$q.dark.isActive ? 'text-grey-3' : 'text-grey-8'">
                <q-icon name="pie_chart" color="primary" size="20px" />
                Bilanciamento Giorni Liberi Richiesti dai Docenti
              </div>
              <div class="text-caption text-grey-6">
                Fai clic su un badge per filtrare rapidamente la tabella per quel giorno libero.
              </div>
            </div>

            <div class="row items-center gap-2 flex-wrap q-pt-xs">
              <q-chip
                clickable
                :outline="filterDayOff !== null"
                :color="filterDayOff === null ? 'primary' : 'grey-7'"
                :text-color="filterDayOff === null ? 'white' : 'grey-8'"
                class="font-bold"
                @click="filterDayOff = null"
              >
                Tutti i Docenti ({{ quickPrefs.length }})
              </q-chip>
              <q-chip
                v-for="d in [
                  { id: 1, label: 'Lunedì', color: 'indigo' },
                  { id: 2, label: 'Martedì', color: 'cyan-9' },
                  { id: 3, label: 'Mercoledì', color: 'teal-8' },
                  { id: 4, label: 'Giovedì', color: 'amber-9' },
                  { id: 5, label: 'Venerdì', color: 'deep-orange-8' },
                  { id: 6, label: 'Sabato', color: 'purple-8' }
                ]"
                :key="d.id"
                clickable
                :outline="filterDayOff !== d.id"
                :color="d.color"
                text-color="white"
                class="font-bold"
                @click="filterDayOff = (filterDayOff === d.id ? null : d.id)"
              >
                {{ d.label }}: {{ quickPrefsDayOffCounts[d.id] || 0 }}
              </q-chip>
              <q-chip
                clickable
                :outline="filterDayOff !== 0"
                color="grey-6"
                text-color="white"
                class="font-bold"
                @click="filterDayOff = (filterDayOff === 0 ? null : 0)"
              >
                Nessun Giorno Libero: {{ quickPrefsDayOffCounts[0] || 0 }}
              </q-chip>
            </div>
          </q-card-section>
        </q-card>

        <!-- Filters & Search Toolbar -->
        <q-card flat bordered class="rounded-2xl shadow-sm q-mb-md" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
          <q-card-section class="q-py-sm">
            <div class="row items-center gap-3">
              <div class="col-12 col-sm-4">
                <q-input
                  v-model="quickPrefsSearch"
                  dense
                  outlined
                  placeholder="Cerca docente o materia..."
                  clearable
                  class="bg-white rounded-lg"
                  :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
                >
                  <template v-slot:prepend>
                    <q-icon name="search" color="grey-6" />
                  </template>
                </q-input>
              </div>

              <div class="col-12 col-sm-3">
                <q-select
                  v-model="filterDayOff"
                  :options="[
                    { value: null, label: 'Tutti i giorni liberi' },
                    ...dayOffOptions
                  ]"
                  option-value="value"
                  option-label="label"
                  emit-value
                  map-options
                  dense
                  outlined
                  label="Filtra per Giorno Libero"
                  class="bg-white rounded-lg"
                  :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
                />
              </div>

              <div class="col-12 col-sm-3">
                <q-select
                  v-model="filterTimeSlot"
                  :options="[
                    { value: null, label: 'Tutte le preferenze orarie' },
                    ...timeSlotPrefOptions
                  ]"
                  option-value="value"
                  option-label="label"
                  emit-value
                  map-options
                  dense
                  outlined
                  label="Filtra per Fascia Oraria"
                  class="bg-white rounded-lg"
                  :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
                />
              </div>

              <div class="col text-right text-caption text-grey-6">
                Docenti visualizzati: <strong>{{ filteredQuickPrefs.length }}</strong> di {{ quickPrefs.length }}
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Main Tabular Representation Card -->
        <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
          <div v-if="loadingQuickPrefs" class="row justify-center q-pa-xl">
            <q-spinner-dots size="48px" color="primary" />
          </div>

          <div v-else-if="filteredQuickPrefs.length === 0" class="text-center q-pa-xl text-slate-400">
            <q-icon name="group_off" size="48px" class="opacity-40 q-mb-sm" />
            <div class="text-h6 text-grey-7">Nessun docente trovato</div>
            <p class="text-caption text-grey-5">Modifica i criteri di ricerca o azzera i filtri per visualizzare i docenti.</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full border-collapse">
              <thead>
                <tr class="bg-slate-100 text-slate-700 text-left" :class="$q.dark.isActive ? 'bg-grey-9 text-grey-2' : 'bg-slate-100 text-slate-700'">
                  <th class="p-3 border border-slate-200">Docente</th>
                  <th class="p-3 border border-slate-200">Materie Insegnate</th>
                  <th class="p-3 border border-slate-200 w-56 text-center">Giorno Libero Desiderato</th>
                  <th class="p-3 border border-slate-200 w-64 text-center">Preferenza Fascia Oraria</th>
                  <th class="p-3 border border-slate-200 w-28 text-center">Max Ore/Giorno</th>
                  <th class="p-3 border border-slate-200 w-36 text-center">Riepilogo Ore</th>
                  <th class="p-3 border border-slate-200 w-36 text-center">Azioni</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="row in filteredQuickPrefs"
                  :key="row.teacher_id"
                  class="hover:bg-slate-50 transition-colors"
                  :class="$q.dark.isActive ? 'hover:bg-grey-9' : 'hover:bg-slate-50'"
                >
                  <!-- Teacher Info -->
                  <td class="p-3 border border-slate-200 font-bold">
                    <div class="row items-center gap-2">
                      <q-avatar size="32px" color="primary" text-color="white" class="font-bold text-caption">
                        {{ (row.teacher_name || 'D').substring(0, 2).toUpperCase() }}
                      </q-avatar>
                      <div>
                        <div class="text-body2 font-bold">{{ row.teacher_name || 'Docente' }}</div>
                      </div>
                    </div>
                  </td>

                  <!-- Subject -->
                  <td class="p-3 border border-slate-200">
                    <div class="row items-center gap-1 flex-wrap">
                      <q-badge
                        v-for="s in (row.subject_name ? row.subject_name.split(', ') : ['Generale'])"
                        :key="s"
                        color="indigo-1"
                        text-color="indigo-9"
                        class="font-medium"
                      >
                        {{ s }}
                      </q-badge>
                    </div>
                  </td>

                  <!-- Day Off Selector -->
                  <td class="p-3 border border-slate-200 text-center">
                    <q-select
                      v-model="row.day_off"
                      :options="dayOffOptions"
                      option-value="value"
                      option-label="label"
                      emit-value
                      map-options
                      dense
                      outlined
                      class="w-full text-weight-bold"
                      :bg-color="row.day_off > 0 ? 'amber-1' : ''"
                    >
                      <template v-slot:prepend>
                        <q-icon name="event_busy" :color="row.day_off > 0 ? 'amber-9' : 'grey-5'" size="18px" />
                      </template>
                    </q-select>
                  </td>

                  <!-- Time Slot Preference -->
                  <td class="p-3 border border-slate-200 text-center">
                    <q-select
                      v-model="row.time_slot_pref"
                      :options="timeSlotPrefOptions"
                      option-value="value"
                      option-label="label"
                      emit-value
                      map-options
                      dense
                      outlined
                      class="w-full text-weight-bold"
                      :bg-color="row.time_slot_pref === 'early_hours' ? 'blue-1' : (row.time_slot_pref === 'late_hours' ? 'purple-1' : '')"
                    >
                      <template v-slot:prepend>
                        <q-icon
                          :name="row.time_slot_pref === 'early_hours' ? 'wb_sunny' : (row.time_slot_pref === 'late_hours' ? 'nights_stay' : 'schedule')"
                          :color="row.time_slot_pref === 'early_hours' ? 'blue-8' : (row.time_slot_pref === 'late_hours' ? 'purple-8' : 'grey-5')"
                          size="18px"
                        />
                      </template>
                    </q-select>
                  </td>

                  <!-- Max Hours Per Day (optional) -->
                  <td class="p-3 border border-slate-200 text-center">
                    <q-input
                      v-model.number="row.max_hours_per_day"
                      type="number"
                      min="1"
                      max="8"
                      dense
                      outlined
                      placeholder="Predef."
                      input-class="text-center font-bold"
                    />
                  </td>

                  <!-- Summary badges -->
                  <td class="p-3 border border-slate-200 text-center">
                    <div class="column items-center gap-1">
                      <span v-if="row.preferred_hours_count > 0" class="text-caption text-positive font-bold">
                        🟢 {{ row.preferred_hours_count }}h pref.
                      </span>
                      <span v-if="row.unavailable_hours_count > 0" class="text-caption text-negative font-bold">
                        🔴 {{ row.unavailable_hours_count }}h non disp.
                      </span>
                      <span v-if="!row.preferred_hours_count && !row.unavailable_hours_count" class="text-caption text-grey-5">
                        ⚪ Standard
                      </span>
                    </div>
                  </td>

                  <!-- Row Actions -->
                  <td class="p-3 border border-slate-200 text-center">
                    <div class="row items-center justify-center gap-1">
                      <q-btn
                        dense
                        unelevated
                        color="primary"
                        icon="save"
                        label="Salva"
                        no-caps
                        class="rounded-lg q-px-sm font-bold text-caption"
                        :loading="!!savingRowQuickPrefs[row.teacher_id]"
                        @click="saveSingleRowQuickPreference(row)"
                      >
                        <q-tooltip>Salva desiderata per questo docente</q-tooltip>
                      </q-btn>
                      <q-btn
                        dense
                        flat
                        color="indigo-7"
                        icon="calendar_view_week"
                        class="rounded-lg"
                        @click="openMatrixForTeacher(row.teacher_id)"
                      >
                        <q-tooltip>Apri matrice oraria settimanale dettagliata</q-tooltip>
                      </q-btn>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Bottom Footer Bar -->
          <div class="row items-center justify-between q-pa-md border-t border-slate-200">
            <div class="text-caption text-grey-6">
              💡 Le preferenze impostate in questa tabella aggiornano automaticamente la matrice settimanale del docente e vengono applicate dal generatore orario con massima priorità.
            </div>
            <div class="row items-center gap-2">
              <q-btn
                unelevated
                color="primary"
                icon="save"
                label="Salva Tutti i Desiderata Rapidi"
                no-caps
                class="rounded-xl font-bold q-px-lg shadow-sm"
                :loading="savingAllQuickPrefs"
                @click="saveAllQuickPreferences"
              />
            </div>
          </div>
        </q-card>
      </div>

      <!-- VIEW 2: Matrice Oraria Dettagliata per Singolo Docente -->
      <div v-else-if="desiderataViewMode === 'single_matrix'">
        <!-- Card: Matrice e Modifica Desiderata Docenti (per Vicario / Responsabile Orario) -->
        <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
          <q-card-section class="row items-center justify-between gap-4">
            <div>
              <div class="text-h6 text-weight-bold row items-center gap-2">
                <q-icon name="person_search" color="primary" />
                Gestione Desiderata per Singolo Docente (Matrice Oraria)
              </div>
              <div class="text-caption text-grey-6">
                Visualizza o modifica la griglia dei desiderata di qualsiasi docente. Le modifiche salvate hanno effetto prioritario sulla generazione dell'orario.
              </div>
            </div>

            <div style="min-width: 300px;">
              <q-select
                v-model="selectedDesiderataTeacherId"
                :options="teacherOptions"
                option-value="id"
                option-label="displayName"
                emit-value
                map-options
                outlined
                dense
                clearable
                label="Seleziona Docente..."
                @update:model-value="onDesiderataTeacherSelected"
              >
                <template v-slot:prepend>
                  <q-icon name="person" color="primary" />
                </template>
              </q-select>
            </div>
          </q-card-section>

          <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

          <!-- If no teacher selected -->
          <div v-if="!selectedDesiderataTeacherId" class="text-center q-pa-xl text-slate-400">
            <q-icon name="touch_app" size="48px" class="opacity-40 q-mb-sm" />
            <div class="text-h6 text-grey-7">Nessun docente selezionato</div>
            <p class="text-caption text-grey-5">Seleziona un docente dal menu sopra per visualizzare o configurare la sua disponibilità settimanale.</p>
          </div>

          <!-- Teacher Loading -->
          <div v-else-if="loadingTeacherPrefs" class="row justify-center q-pa-xl">
            <q-spinner-dots size="48px" color="primary" />
          </div>

          <!-- Teacher Desiderata Grid -->
          <div v-else class="q-pa-md">
            <!-- Summary bar & Quick Actions -->
            <div class="row items-center justify-between gap-3 q-mb-md">
              <div class="row items-center gap-2">
                <q-badge color="positive" class="q-px-sm q-py-xs font-bold">
                  🟢 {{ countPreferred }} Ore Preferite
                </q-badge>
                <q-badge color="grey-6" class="q-px-sm q-py-xs font-bold">
                  ⚪ {{ countNeutral }} Ore Neutre
                </q-badge>
                <q-badge color="negative" class="q-px-sm q-py-xs font-bold">
                  🔴 {{ countUnavailable }} Ore Non Disponibili
                </q-badge>
              </div>

              <div class="row items-center gap-2">
                <!-- Quick day-off -->
                <q-btn-dropdown dense flat color="primary" icon="event_busy" label="Giorno Libero Rapido" no-caps>
                  <q-list dense>
                    <q-item v-for="d in prefDays" :key="d.index" clickable v-close-popup @click="setFullDayOff(d.index)">
                      <q-item-section>{{ d.label }} (tutto il giorno Non Disp.)</q-item-section>
                    </q-item>
                  </q-list>
                </q-btn-dropdown>
                <q-btn flat dense color="grey-7" icon="restart_alt" label="Tutto Neutro" no-caps @click="resetTeacherPrefsToNeutral" />
                <q-btn
                  unelevated
                  color="primary"
                  icon="save"
                  label="Salva Desiderata Docente"
                  no-caps
                  class="rounded-xl q-px-md font-bold"
                  :loading="savingTeacherPrefs"
                  @click="saveTeacherPrefs"
                />
              </div>
            </div>

            <!-- Weekly Matrix -->
            <div class="overflow-x-auto border rounded-xl">
              <table class="w-full border-collapse">
                <thead>
                  <tr class="bg-slate-100 text-slate-700">
                    <th class="p-3 border border-slate-200 text-left w-24">Ora</th>
                    <th v-for="d in prefDays" :key="d.index" class="p-3 border border-slate-200 text-center font-bold">
                      {{ d.label }}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="h in 8" :key="h">
                    <td class="p-3 border border-slate-200 font-bold text-center bg-slate-50 text-slate-800">
                      {{ h }}ª Ora
                    </td>
                    <td
                      v-for="d in prefDays"
                      :key="d.index + '-' + h"
                      class="p-2 border border-slate-200 text-center cursor-pointer transition-colors select-none"
                      :class="getTeacherPrefCellBg(d.index, h)"
                      @click="cycleTeacherPrefCell(d.index, h)"
                    >
                      <div class="py-2 px-1 rounded-lg">
                        <div class="font-bold text-xs" :class="getTeacherPrefCellText(d.index, h)">
                          {{ getTeacherPrefCellLabel(d.index, h) }}
                        </div>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="text-caption text-grey-6 q-mt-sm">
              💡 Fai clic su una casella per alternare: 🟢 Preferito ➔ 🔴 Non Disponibile ➔ ⚪ Neutro.
            </div>
          </div>
        </q-card>
      </div>
    </div>

    <!-- TAB 5: Ore Materie & Limiti Giornalieri per Classe -->
    <div v-show="currentTab === 'curriculum'">
      <q-card flat bordered class="rounded-2xl shadow-sm q-mb-xl" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
        <q-card-section class="row items-center justify-between">
          <div>
            <div class="text-h6 text-weight-bold row items-center gap-2">
              <q-icon name="menu_book" color="primary" />
              Piano Orario Materie & Limiti Giornalieri per Classe
            </div>
            <div class="text-caption text-grey-6">
              Imposta il monte ore settimanale di ciascuna materia e i limiti minimi/massimi di ore al giorno per ogni classe. Puoi ereditare i dati dagli anni scolastici precedenti e modificarli liberamente.
            </div>
          </div>
          <div class="row items-center gap-2">
            <q-btn
              flat
              dense
              icon="history_edu"
              color="indigo-7"
              no-caps
              label="Eredita Anno Precedente"
              class="rounded-lg q-px-sm"
              @click="openInheritModal"
            />
            <q-btn
              flat
              dense
              icon="refresh"
              color="primary"
              no-caps
              label="Aggiorna"
              @click="fetchClassesPlans"
            />
          </div>
        </q-card-section>

        <q-separator :class="$q.dark.isActive ? 'border-grey-8' : 'border-slate-100'" />

        <!-- Filters: Academic Year & Class -->
        <q-card-section class="bg-slate-50 border-b border-slate-100 row items-center gap-4 q-py-md" :class="$q.dark.isActive ? 'bg-grey-9 border-grey-8' : 'bg-slate-50'">
          <div class="col-12 col-sm-3">
            <q-select
              v-model="selectedAcademicYear"
              :options="academicYears"
              outlined
              dense
              label="Anno Scolastico"
              class="bg-white rounded-lg"
              :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
              @update:model-value="fetchClassesPlans"
            />
          </div>
          <div class="col-12 col-sm-4">
            <q-select
              v-model="selectedClassId"
              :options="classPlanOptions"
              option-value="id"
              option-label="name"
              emit-value
              map-options
              outlined
              dense
              label="Seleziona Classe da configurare"
              class="bg-white rounded-lg"
              :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
              @update:model-value="onClassPlanSelected"
            />
          </div>
          <div v-if="currentClassPlan" class="col row items-center gap-2 justify-end">
            <q-chip color="primary" text-color="white" icon="schedule" class="text-weight-bold">
              Totale: {{ currentClassPlan.total_hours_week }} ore/sett.
            </q-chip>
            <q-chip color="teal" text-color="white" icon="calendar_view_day">
              Media: ~{{ (currentClassPlan.total_hours_week / 5).toFixed(1) }} ore/giorno
            </q-chip>
          </div>
        </q-card-section>

        <div v-if="loadingPlans" class="row justify-center q-pa-xl">
          <q-spinner-dots size="48px" color="primary" />
        </div>

        <div v-else-if="!currentClassPlan" class="text-center q-pa-xl">
          <q-icon name="class" size="48px" color="grey-5" />
          <div class="text-h6 text-grey-6 q-mt-sm">Seleziona una classe per visualizzare e modificare il piano orario e i limiti</div>
          <p class="text-caption text-grey-5">Oppure fai clic su "Eredita da Anno Precedente" per copiare tutti i dati in blocco.</p>
        </div>

        <!-- Class Curriculum Plan Editor -->
        <div v-else class="q-pa-md">
          <!-- Daily Limits Box -->
          <div class="p-4 rounded-xl border border-indigo-100 bg-indigo-50/50 q-mb-lg" :class="$q.dark.isActive ? 'bg-indigo-950/20 border-indigo-900' : 'bg-indigo-50/50'">
            <div class="text-subtitle1 text-weight-bold text-indigo-900 row items-center gap-2 q-mb-xs" :class="$q.dark.isActive ? 'text-indigo-200' : 'text-indigo-900'">
              <q-icon name="timelapse" size="22px" />
              Limiti Giornalieri per la Classe {{ currentClassPlan.class_name }}
            </div>
            <p class="text-caption text-grey-7 q-mb-md">
              Definisci il numero minimo e massimo di ore che gli studenti di questa classe possono frequentare ogni giorno. Il generatore rispetterà rigorosamente questi limiti.
            </p>

            <div class="row items-center q-col-gutter-md">
              <div class="col-12 col-sm-3">
                <q-input
                  v-model.number="currentClassPlan.min_hours_per_day"
                  type="number"
                  min="1"
                  max="8"
                  outlined
                  dense
                  label="Ore Minime al Giorno"
                  hint="Es. 4 o 5 ore al giorno"
                  class="bg-white rounded-lg"
                  :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
                />
              </div>
              <div class="col-12 col-sm-3">
                <q-input
                  v-model.number="currentClassPlan.max_hours_per_day"
                  type="number"
                  min="1"
                  max="8"
                  outlined
                  dense
                  label="Ore Massime al Giorno"
                  hint="Es. 5 o 6 ore al giorno"
                  class="bg-white rounded-lg"
                  :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
                />
              </div>
              <div class="col-12 col-sm-6 text-caption text-grey-7">
                🛡️ <strong>Vincolo Algoritmico:</strong> Nessuna giornata supererà <strong>{{ currentClassPlan.max_hours_per_day }} ore</strong> e l'algoritmo cercherà di assegnare almeno <strong>{{ currentClassPlan.min_hours_per_day }} ore</strong> al giorno.
              </div>
            </div>
          </div>

          <!-- Subjects & Weekly Hours Table -->
          <div class="row items-center justify-between q-mb-sm">
            <div class="text-subtitle1 text-weight-bold row items-center gap-2">
              <q-icon name="list_alt" color="primary" />
              Materie & Ore Settimanali (Monte Ore)
            </div>
            <div class="row items-center gap-2">
              <q-btn
                unelevated
                color="secondary"
                icon="add"
                label="Aggiungi Materia"
                no-caps
                class="rounded-xl q-px-md font-bold"
                @click="openAddSubjectModal"
              />
            </div>
          </div>

          <div v-if="!currentClassPlan.subjects || currentClassPlan.subjects.length === 0" class="text-center p-8 border border-dashed rounded-xl q-mb-md">
            <q-icon name="menu_book" size="36px" color="grey-4" />
            <div class="text-subtitle2 text-grey-6 q-mt-xs">Nessuna materia ancora configurata per questa classe.</div>
            <div class="text-caption text-grey-5 q-mb-sm">Aggiungi manualmente le materie oppure eredita il piano dall'anno scolastico precedente.</div>
            <q-btn
              unelevated
              color="indigo-7"
              icon="history_edu"
              label="Eredita da Anno Precedente"
              no-caps
              class="rounded-xl q-px-md font-bold"
              @click="openInheritModal"
            />
          </div>

          <div v-else class="overflow-x-auto border rounded-xl q-mb-lg">
            <table class="w-full border-collapse">
              <thead>
                <tr class="bg-slate-100 text-slate-700 text-left" :class="$q.dark.isActive ? 'bg-grey-9 text-grey-2' : 'bg-slate-100 text-slate-700'">
                  <th class="p-3 border border-slate-200">Materia</th>
                  <th class="p-3 border border-slate-200">Docente Titolare</th>
                  <th class="p-3 border border-slate-200 text-center w-48">Ore Settimanali</th>
                  <th class="p-3 border border-slate-200 text-center w-24">Azioni</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(sub, sIdx) in (currentClassPlan.subjects || [])"
                  :key="sub.subject_id"
                  class="hover:bg-slate-50 transition-colors"
                  :class="$q.dark.isActive ? 'hover:bg-grey-9' : 'hover:bg-slate-50'"
                >
                  <td class="p-3 border border-slate-200 font-bold">
                    <div class="row items-center gap-2">
                      <q-icon name="menu_book" color="primary" size="20px" />
                      <span>{{ sub.subject_name }}</span>
                    </div>
                  </td>
                  <td class="p-3 border border-slate-200 text-grey-8" :class="$q.dark.isActive ? 'text-grey-3' : 'text-grey-8'">
                    <div class="row items-center gap-2">
                      <q-icon name="person" size="18px" :color="sub.teacher_name ? 'positive' : 'grey-5'" />
                      <span>{{ sub.teacher_name || 'Docente da Nominare (Cattedra non assegnata)' }}</span>
                    </div>
                  </td>
                  <td class="p-3 border border-slate-200 text-center">
                    <div class="row items-center justify-center gap-1">
                      <q-btn
                        dense
                        round
                        flat
                        icon="remove"
                        size="sm"
                        color="primary"
                        :disable="sub.hours_per_week <= 1"
                        @click="sub.hours_per_week = Math.max(1, sub.hours_per_week - 1); currentClassPlan.total_hours_week = (currentClassPlan.subjects || []).reduce((acc, s) => acc + (s.hours_per_week || 0), 0)"
                      />
                      <q-input
                        v-model.number="sub.hours_per_week"
                        type="number"
                        min="1"
                        max="20"
                        dense
                        outlined
                        class="w-16 text-center font-bold"
                        input-class="text-center font-bold"
                        @update:model-value="currentClassPlan.total_hours_week = (currentClassPlan.subjects || []).reduce((acc, s) => acc + (s.hours_per_week || 0), 0)"
                      />
                      <q-btn
                        dense
                        round
                        flat
                        icon="add"
                        size="sm"
                        color="primary"
                        :disable="sub.hours_per_week >= 20"
                        @click="sub.hours_per_week = Math.min(20, sub.hours_per_week + 1); currentClassPlan.total_hours_week = (currentClassPlan.subjects || []).reduce((acc, s) => acc + (s.hours_per_week || 0), 0)"
                      />
                      <span class="text-caption text-grey-6 ml-1">ore/sett.</span>
                    </div>
                  </td>
                  <td class="p-3 border border-slate-200 text-center">
                    <q-btn
                      flat
                      round
                      dense
                      icon="delete"
                      color="negative"
                      @click="removeSubjectFromCurrentPlan(sIdx)"
                    >
                      <q-tooltip>Rimuovi materia da questa classe</q-tooltip>
                    </q-btn>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Save Actions -->
          <div class="row items-center justify-between q-mt-md pt-4 border-t border-slate-200">
            <div class="text-caption text-grey-6">
              💡 Le modifiche al monte ore e ai limiti giornalieri avranno effetto immediato sul motore di calcolo orario.
            </div>
            <div class="row items-center gap-3">
              <q-btn
                flat
                label="Ripristina Valori"
                color="grey-7"
                no-caps
                @click="onClassPlanSelected(selectedClassId)"
              />
              <q-btn
                unelevated
                color="primary"
                icon="save"
                label="Salva Piano Orario & Limiti"
                no-caps
                class="rounded-xl q-px-lg font-bold shadow-sm"
                :loading="savingClassPlan"
                @click="saveCurrentClassPlan"
              />
            </div>
          </div>
        </div>
      </q-card>
    </div>

    <!-- Modal: Add Room Requirement with Lab Hours -->
    <q-dialog v-model="showReqModal" persistent>
      <q-card style="min-width: 460px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold">{{ t('timetableConstraints.modalTitle') }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-select
            v-model="reqForm.subject_id"
            :options="subjects"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.selectSubject')"
          />

          <q-select
            v-model="reqForm.required_room_type"
            :options="roomTypeOptions"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.selectRoomType')"
          />

          <q-input
            v-model.number="reqForm.lab_hours"
            type="number"
            min="1"
            max="12"
            outlined
            dense
            :label="t('timetableConstraints.labHoursLabel')"
            :hint="t('timetableConstraints.labHoursHint')"
          />

          <q-toggle
            v-model="reqForm.is_mandatory"
            :label="t('timetableConstraints.mandatoryDesc')"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="t('timetableConstraints.saveBtn')" no-caps class="rounded-xl q-px-md font-bold" @click="saveReq" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Add Gruppo Linguistico / Articolato -->
    <q-dialog v-model="showGroupModal" persistent>
      <q-card style="min-width: 520px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold row items-center gap-2">
            <q-icon name="groups" color="primary" />
            {{ t('timetableConstraints.groupModalTitle') }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-input
            v-model="groupForm.name"
            outlined
            dense
            :label="t('timetableConstraints.groupName')"
            placeholder="es. Gruppo Spagnolo 3A-3B"
          />

          <q-select
            v-model="groupForm.subject_id"
            :options="subjects"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.groupSubject')"
          />

          <q-select
            v-model="groupForm.teacher_id"
            :options="allTeachers"
            option-value="id"
            option-label="label"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.groupTeacher')"
          />

          <q-select
            v-model="groupForm.class_ids"
            :options="classes"
            option-value="id"
            option-label="label"
            emit-value
            map-options
            multiple
            use-chips
            outlined
            dense
            :label="t('timetableConstraints.groupClasses')"
          />

          <q-input
            v-model.number="groupForm.hours_per_week"
            type="number"
            min="1"
            max="10"
            outlined
            dense
            :label="t('timetableConstraints.groupHours')"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="t('timetableConstraints.saveBtn') || 'Salva Gruppo'" no-caps class="rounded-xl q-px-md font-bold" @click="saveGroup" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Add Spezzone Orario per Classe -->
    <q-dialog v-model="showSpezzoneModal" persistent>
      <q-card style="min-width: 500px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold row items-center gap-2">
            <q-icon name="add_circle_outline" color="primary" />
            {{ t('timetableConstraints.spezzoneModalTitle') }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-select
            v-model="spezzoneForm.class_id"
            :options="classes"
            option-value="id"
            option-label="label"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.spezzoneClass')"
          />

          <q-select
            v-model="spezzoneForm.subject_id"
            :options="subjects"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            :label="t('timetableConstraints.spezzoneSubject')"
            @update:model-value="onSpezzoneSubjectChanged"
          />

          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-input
                v-model.number="spezzoneForm.hours_per_week"
                type="number"
                min="1"
                max="18"
                outlined
                dense
                :label="t('timetableConstraints.spezzoneHours')"
              />
            </div>
            <div class="col-12 col-sm-6">
              <q-input
                v-model="spezzoneForm.placeholder_name"
                outlined
                dense
                :label="t('timetableConstraints.spezzonePlaceholder')"
              />
            </div>
          </div>

          <q-input
            v-model="spezzoneForm.notes"
            type="textarea"
            rows="2"
            outlined
            dense
            :label="t('timetableConstraints.spezzoneNotes')"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat :label="t('common.cancel') || 'Annulla'" no-caps v-close-popup />
          <q-btn unelevated color="primary" :label="t('timetableConstraints.saveBtn') || 'Salva'" no-caps class="rounded-xl q-px-md font-bold" @click="saveSpezzone" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Eredita da Anno Precedente -->
    <q-dialog v-model="showInheritModal" persistent>
      <q-card style="min-width: 480px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold row items-center gap-2">
            <q-icon name="history_edu" color="indigo-7" />
            Eredita Piano Orario & Limiti da Anno Precedente
          </div>
          <p class="text-caption text-grey-6 q-mt-xs q-mb-none">
            Copia il monte ore di ogni materia e i limiti orari giornalieri (min/max) dall'anno scolastico precedente. Potrai poi apportare qualsiasi modifica.
          </p>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-select
            v-model="inheritSourceYear"
            :options="academicYears"
            outlined
            dense
            label="Anno Scolastico di Origine"
            hint="Seleziona l'anno da cui copiare la struttura delle materie"
          />

          <div class="q-mt-md">
            <div class="text-subtitle2 text-weight-bold q-mb-xs">Ambito di Applicazione:</div>
            <div class="column q-gutter-sm">
              <q-btn
                :outline="inheritScope !== 'single'"
                :unelevated="inheritScope === 'single'"
                :color="inheritScope === 'single' ? 'indigo-7' : 'grey-7'"
                no-caps
                align="left"
                class="rounded-xl q-px-md font-bold text-left full-width"
                :label="`Solo per la classe selezionata (${currentClassPlan ? currentClassPlan.class_name : 'classe corrente'})`"
                @click="inheritScope = 'single'"
              />
              <q-btn
                :outline="inheritScope !== 'all'"
                :unelevated="inheritScope === 'all'"
                :color="inheritScope === 'all' ? 'indigo-7' : 'grey-7'"
                no-caps
                align="left"
                class="rounded-xl q-px-md font-bold text-left full-width"
                label="Per TUTTE le classi dell'istituto (clonazione di massa)"
                @click="inheritScope = 'all'"
              />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat label="Annulla" no-caps v-close-popup />
          <q-btn
            unelevated
            color="indigo-7"
            icon="content_copy"
            label="Eredita Dati"
            no-caps
            class="rounded-xl q-px-md font-bold"
            :loading="inheriting"
            @click="executeInherit"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Modal: Aggiungi Materia alla Classe -->
    <q-dialog v-model="showAddSubjectModal" persistent>
      <q-card style="min-width: 460px;" class="rounded-2xl q-pa-sm">
        <q-card-section>
          <div class="text-h6 text-weight-bold row items-center gap-2">
            <q-icon name="add_circle" color="secondary" />
            Aggiungi Materia alla Classe {{ currentClassPlan ? currentClassPlan.class_name : '' }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-md">
          <q-select
            v-model="newSubjectForm.subject_id"
            :options="subjects"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            label="Materia *"
          />

          <q-select
            v-model="newSubjectForm.teacher_id"
            :options="teacherOptions"
            option-value="id"
            option-label="name"
            emit-value
            map-options
            outlined
            dense
            clearable
            label="Docente Assegnato (Opzionale)"
            hint="Se vuoto, la cattedra risulterà da nominare"
          />

          <q-input
            v-model.number="newSubjectForm.hours_per_week"
            type="number"
            min="1"
            max="20"
            outlined
            dense
            label="Ore Settimanali *"
          />
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pb-md">
          <q-btn flat label="Annulla" no-caps v-close-popup />
          <q-btn
            unelevated
            color="secondary"
            label="Aggiungi al Piano"
            no-caps
            class="rounded-xl q-px-md font-bold"
            :disable="!newSubjectForm.subject_id || newSubjectForm.hours_per_week <= 0"
            @click="addSubjectToCurrentPlan"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import timetableGenService from '@/services/timetableGenService';
import adminService from '@/services/adminService';
import api from '@/services/api';

const { t } = useI18n();
const $q = useQuasar();

const currentTab = ref('rooms');

// ----------------- TAB 1: Subject Room Requirements -----------------
const loadingReqs = ref(false);
const roomReqs = ref([]);
const subjects = ref([]);

const showReqModal = ref(false);
const reqForm = ref({
  subject_id: null,
  required_room_type: 'lab_informatica',
  lab_hours: 2,
  is_mandatory: true
});

const roomTypeOptions = computed(() => [
  { label: t('roomsManagement.types.lab_info') || 'Laboratorio di Informatica', value: 'lab_informatica' },
  { label: t('roomsManagement.types.lab_science') || 'Laboratorio di Scienze', value: 'lab_scienze' },
  { label: t('roomsManagement.types.lab_lang') || 'Laboratorio di Lingue', value: 'lab_lingue' },
  { label: t('roomsManagement.types.lab_chemistry') || 'Laboratorio di Chimica', value: 'lab_chimica' },
  { label: t('roomsManagement.types.lab_physics') || 'Laboratorio di Fisica', value: 'lab_fisica' },
  { label: t('roomsManagement.types.art') || 'Laboratorio di Arte', value: 'lab_arte' },
  { label: t('roomsManagement.types.gym') || 'Palestra', value: 'palestra' },
  { label: t('roomsManagement.types.auditorium') || 'Aula Magna', value: 'aula_magna' },
  { label: t('roomsManagement.types.library') || 'Biblioteca', value: 'biblioteca' }
]);

const roomReqColumns = computed(() => [
  { name: 'subject_name', label: t('timetableConstraints.colSubject'), field: 'subject_name', sortable: true, align: 'left' },
  { name: 'required_room_type', label: t('timetableConstraints.colRoomType'), field: 'required_room_type', align: 'left' },
  { name: 'lab_hours', label: t('timetableConstraints.colLabHours'), field: 'lab_hours', align: 'center' },
  { name: 'is_mandatory', label: t('timetableConstraints.colConstraint'), field: 'is_mandatory', align: 'center' },
  { name: 'actions', label: t('timetableConstraints.colActions'), field: 'id', align: 'right' }
]);

function getRoomTypeLabel(val) {
  const f = roomTypeOptions.value.find(o => o.value === val);
  return f ? f.label : val;
}

async function fetchRoomReqs() {
  loadingReqs.value = true;
  try {
    const res = await timetableGenService.getRoomRequirements();
    roomReqs.value = res.data || [];
  } catch {
    $q.notify({ type: 'negative', message: t('timetableConstraints.notifyFetchError') });
  } finally {
    loadingReqs.value = false;
  }
}

async function fetchSubjects() {
  try {
    const res = await api.get('/subjects');
    subjects.value = res.data || [];
  } catch (err) {
    console.error('Error fetching subjects', err);
  }
}

function openRoomReqModal() {
  reqForm.value = {
    subject_id: subjects.value[0]?.id || null,
    required_room_type: 'lab_informatica',
    lab_hours: 2,
    is_mandatory: true
  };
  showReqModal.value = true;
}

async function saveReq() {
  if (!reqForm.value.subject_id || !reqForm.value.required_room_type) return;
  try {
    await timetableGenService.saveRoomRequirement(reqForm.value);
    $q.notify({ type: 'positive', message: t('timetableConstraints.notifySaved') });
    showReqModal.value = false;
    fetchRoomReqs();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('timetableConstraints.notifySaveError') });
  }
}

function confirmDeleteReq(row) {
  $q.dialog({
    title: t('timetableConstraints.removeTitle'),
    message: t('timetableConstraints.removeMsg', { subject: row.subject_name }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await timetableGenService.deleteRoomRequirement(row.id);
      $q.notify({ type: 'positive', message: t('timetableConstraints.notifyRemoved') });
      fetchRoomReqs();
    } catch {
      $q.notify({ type: 'negative', message: t('timetableConstraints.notifyDeleteError') });
    }
  });
}

// ----------------- TAB 2: Gruppi Linguistici & Articolati -----------------
const allTeachers = ref([]);
const associatedGroupsList = ref([]);
const showGroupModal = ref(false);
const groupForm = ref({
  name: '',
  subject_id: null,
  teacher_id: null,
  class_ids: [],
  hours_per_week: 3
});

const groupColumns = computed(() => [
  { name: 'name', label: t('timetableConstraints.colGroup'), field: 'name', sortable: true, align: 'left' },
  { name: 'subject_name', label: t('timetableConstraints.colSubject'), field: 'subject_name', align: 'left' },
  { name: 'teacher_name', label: t('timetableConstraints.colTeacher'), field: 'teacher_name', align: 'left' },
  { name: 'classes', label: t('timetableConstraints.colClasses'), field: 'class_names', align: 'left' },
  { name: 'hours_per_week', label: t('timetableConstraints.colHours'), field: 'hours_per_week', align: 'center' },
  { name: 'actions', label: t('timetableConstraints.colActions'), field: 'id', align: 'right' }
]);

async function fetchTeachers() {
  try {
    const res = await adminService.getTeachersList();
    const raw = res.data || [];
    allTeachers.value = raw.map(t => ({
      id: t.id || t.user_id,
      user_id: t.user_id,
      label: `${t.last_name || ''} ${t.first_name || ''}`.trim() || t.email,
      email: t.email
    })).sort((a, b) => a.label.localeCompare(b.label));
  } catch (e) {
    console.error('Failed fetching teachers', e);
  }
}

function openGroupModal() {
  groupForm.value = {
    name: '',
    subject_id: subjects.value[0]?.id || null,
    teacher_id: allTeachers.value[0]?.id || null,
    class_ids: [],
    hours_per_week: 3
  };
  showGroupModal.value = true;
}

async function saveGroup() {
  if (!groupForm.value.name || !groupForm.value.subject_id || !groupForm.value.teacher_id || groupForm.value.class_ids.length < 2) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi e seleziona almeno 2 classi per il gruppo linguistico/articolato' });
    return;
  }

  const subj = subjects.value.find(s => s.id === groupForm.value.subject_id);
  const teacher = allTeachers.value.find(t => t.id === groupForm.value.teacher_id);
  const classNames = groupForm.value.class_ids.map(cid => {
    const c = classes.value.find(cl => cl.id === cid);
    return c ? c.label : cid;
  });

  const payload = {
    constraint_type: 'associated_group',
    target_type: 'group',
    is_hard: true,
    priority: 10,
    is_active: true,
    parameters: {
      name: groupForm.value.name,
      subject_id: groupForm.value.subject_id,
      subject_name: subj?.name || 'Materia',
      teacher_id: groupForm.value.teacher_id,
      teacher_name: teacher?.label || 'Docente',
      class_ids: groupForm.value.class_ids,
      class_names: classNames,
      hours_per_week: groupForm.value.hours_per_week || 3
    }
  };

  try {
    await timetableGenService.saveConstraint(payload);
    $q.notify({ type: 'positive', message: 'Gruppo linguistico / articolato configurato con successo' });
    showGroupModal.value = false;
    fetchConstraints();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore nel salvataggio del gruppo' });
  }
}

function confirmDeleteGroup(row) {
  $q.dialog({
    title: t('timetableConstraints.removeGroupTitle'),
    message: t('timetableConstraints.removeGroupMsg', { name: row.name }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await timetableGenService.deleteConstraint(row.id);
      $q.notify({ type: 'positive', message: 'Gruppo associato rimosso' });
      fetchConstraints();
    } catch {
      $q.notify({ type: 'negative', message: 'Errore nella rimozione del gruppo' });
    }
  });
}

// ----------------- TAB 3: Spezzoni Orari & Cattedre Vacanti -----------------
const loadingConstraints = ref(false);
const constraints = ref([]);
const spezzoniList = ref([]);
const classes = ref([]);

const showSpezzoneModal = ref(false);
const spezzoneForm = ref({
  class_id: null,
  subject_id: null,
  hours_per_week: 2,
  placeholder_name: 'Docente da Nominare (Cattedra non assegnata)',
  notes: ''
});

const spezzoniColumns = computed(() => [
  { name: 'class_name', label: t('timetableConstraints.colClass'), field: 'class_name', sortable: true, align: 'left' },
  { name: 'subject_name', label: t('timetableConstraints.colSubject'), field: 'subject_name', sortable: true, align: 'left' },
  { name: 'hours_per_week', label: t('timetableConstraints.colHours'), field: 'hours_per_week', align: 'center' },
  { name: 'placeholder_name', label: t('timetableConstraints.colPlaceholder'), field: 'placeholder_name', align: 'left' },
  { name: 'notes', label: t('timetableConstraints.colNotes'), field: 'notes', align: 'left' },
  { name: 'actions', label: t('timetableConstraints.colActions'), field: 'id', align: 'right' }
]);

async function fetchClasses() {
  try {
    const res = await adminService.getClasses();
    const raw = res.data || [];
    classes.value = raw.map(c => ({
      id: c.id,
      label: `${c.year || ''}${c.section || ''} ${c.name || ''}`.trim() || c.id
    })).sort((a, b) => a.label.localeCompare(b.label));
  } catch (err) {
    console.error('Error fetching classes', err);
  }
}

async function fetchConstraints() {
  loadingConstraints.value = true;
  try {
    const res = await timetableGenService.getConstraints();
    const raw = res.data || [];
    constraints.value = raw;

    spezzoniList.value = raw
      .filter(c => c.constraint_type === 'unassigned_hours')
      .map(c => {
        const p = typeof c.parameters === 'string' ? JSON.parse(c.parameters) : (c.parameters || {});
        const cls = classes.value.find(cl => cl.id === c.target_id);
        return {
          id: c.id,
          class_id: c.target_id,
          class_name: cls ? cls.label : (c.target_id?.substring(0, 8) || '-'),
          subject_id: p.subject_id,
          subject_name: p.subject_name || '-',
          hours_per_week: p.hours_per_week || 2,
          placeholder_name: p.placeholder_name || 'Docente da Nominare',
          notes: p.notes || ''
        };
      });

    associatedGroupsList.value = raw
      .filter(c => c.constraint_type === 'associated_group')
      .map(c => {
        const p = typeof c.parameters === 'string' ? JSON.parse(c.parameters) : (c.parameters || {});
        return {
          id: c.id,
          name: p.name || 'Gruppo Senza Nome',
          subject_id: p.subject_id,
          subject_name: p.subject_name || 'Materia',
          teacher_id: p.teacher_id,
          teacher_name: p.teacher_name || 'Docente',
          class_ids: p.class_ids || [],
          class_names: p.class_names || [],
          hours_per_week: p.hours_per_week || 3
        };
      });
  } catch (err) {
    console.error('Error fetching constraints', err);
  } finally {
    loadingConstraints.value = false;
  }
}

function openSpezzoneModal() {
  spezzoneForm.value = {
    class_id: classes.value[0]?.id || null,
    subject_id: subjects.value[0]?.id || null,
    hours_per_week: 2,
    placeholder_name: `Docente da Nominare (${subjects.value[0]?.name || ''})`,
    notes: ''
  };
  showSpezzoneModal.value = true;
}

function onSpezzoneSubjectChanged(subjId) {
  const subj = subjects.value.find(s => s.id === subjId);
  if (subj) {
    spezzoneForm.value.placeholder_name = `Docente da Nominare (${subj.name})`;
  }
}

async function saveSpezzone() {
  if (!spezzoneForm.value.class_id || !spezzoneForm.value.subject_id) return;
  const subj = subjects.value.find(s => s.id === spezzoneForm.value.subject_id);

  const payload = {
    constraint_type: 'unassigned_hours',
    target_type: 'class',
    target_id: spezzoneForm.value.class_id,
    is_hard: true,
    priority: 10,
    is_active: true,
    parameters: {
      subject_id: spezzoneForm.value.subject_id,
      subject_name: subj?.name || 'Materia',
      hours_per_week: spezzoneForm.value.hours_per_week,
      placeholder_name: spezzoneForm.value.placeholder_name || 'Docente da Nominare',
      notes: spezzoneForm.value.notes || ''
    }
  };

  try {
    await timetableGenService.saveConstraint(payload);
    $q.notify({ type: 'positive', message: t('timetableConstraints.spezzoneSaved') });
    showSpezzoneModal.value = false;
    fetchConstraints();
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || t('timetableConstraints.notifySaveError') });
  }
}

function confirmDeleteSpezzone(row) {
  $q.dialog({
    title: t('timetableConstraints.removeSpezzoneTitle'),
    message: t('timetableConstraints.removeSpezzoneMsg', { hours: row.hours_per_week, subject: row.subject_name, class: row.class_name }),
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await timetableGenService.deleteConstraint(row.id);
      $q.notify({ type: 'positive', message: t('timetableConstraints.spezzoneRemoved') });
      fetchConstraints();
    } catch {
      $q.notify({ type: 'negative', message: t('timetableConstraints.notifyDeleteError') });
    }
  });
}

// ----------------- TAB 4: Desiderata & Preferenze Docenti -----------------
const isDesiderataWindowOpen = ref(false);
const loadingWindow = ref(false);
const selectedDesiderataTeacherId = ref(null);
const teacherPrefsGrid = ref({});
const loadingTeacherPrefs = ref(false);
const savingTeacherPrefs = ref(false);

const prefDays = computed(() => [
  { index: 1, label: 'Lunedì' },
  { index: 2, label: 'Martedì' },
  { index: 3, label: 'Mercoledì' },
  { index: 4, label: 'Giovedì' },
  { index: 5, label: 'Venerdì' }
]);

const teacherOptions = computed(() => {
  return (allTeachers.value || []).map(t => ({
    id: t.user_id || t.id,
    displayName: t.label || t.email || 'Docente'
  }));
});

function initTeacherPrefsGrid() {
  const g = {};
  for (let d = 1; d <= 5; d++) {
    g[d] = {};
    for (let h = 1; h <= 8; h++) {
      g[d][h] = 'neutral';
    }
  }
  teacherPrefsGrid.value = g;
}

const countPreferred = computed(() => {
  let count = 0;
  for (let d = 1; d <= 5; d++) {
    for (let h = 1; h <= 8; h++) {
      if (teacherPrefsGrid.value[d]?.[h] === 'preferred') count++;
    }
  }
  return count;
});

const countNeutral = computed(() => {
  let count = 0;
  for (let d = 1; d <= 5; d++) {
    for (let h = 1; h <= 8; h++) {
      if ((teacherPrefsGrid.value[d]?.[h] || 'neutral') === 'neutral') count++;
    }
  }
  return count;
});

const countUnavailable = computed(() => {
  let count = 0;
  for (let d = 1; d <= 5; d++) {
    for (let h = 1; h <= 8; h++) {
      if (teacherPrefsGrid.value[d]?.[h] === 'unavailable') count++;
    }
  }
  return count;
});

function getTeacherPrefCellBg(day, hour) {
  const state = teacherPrefsGrid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return 'bg-emerald-100 hover:bg-emerald-200';
  if (state === 'unavailable') return 'bg-rose-100 hover:bg-rose-200';
  return 'bg-white hover:bg-slate-100';
}

function getTeacherPrefCellText(day, hour) {
  const state = teacherPrefsGrid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return 'text-emerald-800';
  if (state === 'unavailable') return 'text-rose-800';
  return 'text-slate-400';
}

function getTeacherPrefCellLabel(day, hour) {
  const state = teacherPrefsGrid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return '🟢 Preferito';
  if (state === 'unavailable') return '🔴 Non Disp.';
  return '⚪ Neutro';
}

function cycleTeacherPrefCell(day, hour) {
  const current = teacherPrefsGrid.value[day]?.[hour] || 'neutral';
  let next = 'neutral';
  if (current === 'neutral') next = 'preferred';
  else if (current === 'preferred') next = 'unavailable';
  else if (current === 'unavailable') next = 'neutral';

  if (!teacherPrefsGrid.value[day]) teacherPrefsGrid.value[day] = {};
  teacherPrefsGrid.value[day][hour] = next;
}

function setFullDayOff(dayIndex) {
  if (!teacherPrefsGrid.value[dayIndex]) teacherPrefsGrid.value[dayIndex] = {};
  for (let h = 1; h <= 8; h++) {
    teacherPrefsGrid.value[dayIndex][h] = 'unavailable';
  }
}

function resetTeacherPrefsToNeutral() {
  initTeacherPrefsGrid();
}

async function fetchDesiderataWindow() {
  if (typeof timetableGenService.getDesiderataWindow !== 'function') return;
  try {
    const res = await timetableGenService.getDesiderataWindow();
    isDesiderataWindowOpen.value = res.data?.is_open || false;
  } catch (err) {
    console.error('Error fetching desiderata window status', err);
  }
}

async function toggleDesiderataWindow(val) {
  if (typeof timetableGenService.setDesiderataWindow !== 'function') return;
  loadingWindow.value = true;
  try {
    await timetableGenService.setDesiderataWindow(val);
    $q.notify({
      type: 'positive',
      message: val ? 'Finestra inserimento desiderata docenti APERTA' : 'Finestra inserimento desiderata docenti CHIUSA'
    });
  } catch (err) {
    isDesiderataWindowOpen.value = !val;
    $q.notify({ type: 'negative', message: 'Errore durante l\'aggiornamento dello stato della finestra' });
  } finally {
    loadingWindow.value = false;
  }
}

async function onDesiderataTeacherSelected(tId) {
  if (!tId) {
    teacherPrefsGrid.value = {};
    return;
  }
  loadingTeacherPrefs.value = true;
  initTeacherPrefsGrid();
  try {
    if (typeof timetableGenService.getPreferences === 'function') {
      const res = await timetableGenService.getPreferences({ teacher_id: tId });
      const prefs = res.data || [];
      for (const p of prefs) {
        if (teacherPrefsGrid.value[p.day_of_week] && teacherPrefsGrid.value[p.day_of_week][p.hour_index] !== undefined) {
          teacherPrefsGrid.value[p.day_of_week][p.hour_index] = p.preference_type;
        }
      }
    }
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel caricamento preferenze del docente' });
  } finally {
    loadingTeacherPrefs.value = false;
  }
}

async function saveTeacherPrefs() {
  if (!selectedDesiderataTeacherId.value) return;
  savingTeacherPrefs.value = true;
  try {
    const entries = [];
    for (let d = 1; d <= 5; d++) {
      for (let h = 1; h <= 8; h++) {
        const val = teacherPrefsGrid.value[d]?.[h] || 'neutral';
        if (val !== 'neutral') {
          entries.push({
            day_of_week: d,
            hour_index: h,
            preference_type: val
          });
        }
      }
    }
    await timetableGenService.savePreferences({ preferences: entries }, { teacher_id: selectedDesiderataTeacherId.value });
    $q.notify({ type: 'positive', message: 'Desiderata docente salvati con successo' });
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Errore salvataggio desiderata docente' });
  } finally {
    savingTeacherPrefs.value = false;
  }
}

// ----------------- Desiderata Rapidi (Tabella Tutti i Docenti) -----------------
const desiderataViewMode = ref('quick_table'); // 'quick_table' | 'single_matrix'
const quickPrefs = ref([]);
const loadingQuickPrefs = ref(false);
const savingAllQuickPrefs = ref(false);
const savingRowQuickPrefs = ref({});
const quickPrefsSearch = ref('');
const filterDayOff = ref(null);
const filterTimeSlot = ref(null);
const quickPrefsDayOffCounts = ref({});

const dayOffOptions = [
  { value: 0, label: 'Nessun giorno libero' },
  { value: 1, label: 'Lunedì' },
  { value: 2, label: 'Martedì' },
  { value: 3, label: 'Mercoledì' },
  { value: 4, label: 'Giovedì' },
  { value: 5, label: 'Venerdì' },
  { value: 6, label: 'Sabato' }
];

const timeSlotPrefOptions = [
  { value: 'none', label: '⚪ Indifferente / Neutro' },
  { value: 'early_hours', label: '🌅 Prime Ore (1ª-3ª)' },
  { value: 'late_hours', label: '🌇 Ultime Ore (4ª-6ª)' }
];

const filteredQuickPrefs = computed(() => {
  let list = quickPrefs.value || [];
  if (quickPrefsSearch.value) {
    const q = quickPrefsSearch.value.toLowerCase().trim();
    list = list.filter(t =>
      (t.teacher_name && t.teacher_name.toLowerCase().includes(q)) ||
      (t.subject_name && t.subject_name.toLowerCase().includes(q))
    );
  }
  if (filterDayOff.value !== null && filterDayOff.value !== undefined && filterDayOff.value !== '') {
    list = list.filter(t => t.day_off === Number(filterDayOff.value));
  }
  if (filterTimeSlot.value) {
    list = list.filter(t => t.time_slot_pref === filterTimeSlot.value);
  }
  return list;
});

async function fetchTeachersQuickPreferences() {
  if (typeof timetableGenService.getTeachersQuickPreferences !== 'function') return;
  loadingQuickPrefs.value = true;
  try {
    const res = await timetableGenService.getTeachersQuickPreferences();
    const data = res.data || {};
    quickPrefs.value = data.teachers || [];
    quickPrefsDayOffCounts.value = data.day_off_counts || {};
  } catch (err) {
    console.error('Error fetching teachers quick preferences', err);
  } finally {
    loadingQuickPrefs.value = false;
  }
}

async function saveAllQuickPreferences() {
  if (typeof timetableGenService.saveTeachersQuickPreferences !== 'function') return;
  savingAllQuickPrefs.value = true;
  try {
    await timetableGenService.saveTeachersQuickPreferences({
      preferences: quickPrefs.value
    });
    $q.notify({
      type: 'positive',
      message: 'Desiderata rapidi (giorni liberi e fasce orarie) di tutti i docenti salvati con successo!'
    });
    await fetchTeachersQuickPreferences();
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore salvataggio desiderata rapidi docenti'
    });
  } finally {
    savingAllQuickPrefs.value = false;
  }
}

async function saveSingleRowQuickPreference(row) {
  if (typeof timetableGenService.saveTeacherQuickPreference !== 'function') return;
  savingRowQuickPrefs.value[row.teacher_id] = true;
  try {
    await timetableGenService.saveTeacherQuickPreference(row.teacher_id, row);
    $q.notify({
      type: 'positive',
      message: `Desiderata per ${row.teacher_name} salvati con successo!`
    });
    await fetchTeachersQuickPreferences();
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || `Errore salvataggio desiderata per ${row.teacher_name}`
    });
  } finally {
    savingRowQuickPrefs.value[row.teacher_id] = false;
  }
}

function openMatrixForTeacher(teacherId) {
  selectedDesiderataTeacherId.value = teacherId;
  desiderataViewMode.value = 'single_matrix';
  onDesiderataTeacherSelected(teacherId);
}

// ----------------- TAB 5: Ore Materie & Limiti Classi (Curriculum Plans) -----------------
const academicYears = ref([]);
const selectedAcademicYear = ref('');
const classesCurriculumPlans = ref([]);
const selectedClassId = ref('');
const currentClassPlan = ref(null);
const loadingPlans = ref(false);
const savingClassPlan = ref(false);

const showInheritModal = ref(false);
const inheritSourceYear = ref('');
const inheritScope = ref('single');
const inheriting = ref(false);

const showAddSubjectModal = ref(false);
const newSubjectForm = ref({
  subject_id: '',
  teacher_id: null,
  hours_per_week: 3
});

const classPlanOptions = computed(() => {
  return classesCurriculumPlans.value.map(p => ({
    id: p.class_id,
    name: p.class_name
  }));
});

async function fetchAcademicYears() {
  if (typeof timetableGenService.getAcademicYears !== 'function') return;
  try {
    const res = await timetableGenService.getAcademicYears();
    academicYears.value = res.data || [];
    if (!selectedAcademicYear.value && academicYears.value.length > 0) {
      selectedAcademicYear.value = academicYears.value[0];
    }
  } catch (err) {
    console.error('Error fetching academic years', err);
  }
}

async function fetchClassesPlans() {
  if (typeof timetableGenService.getClassesCurriculumPlans !== 'function') return;
  loadingPlans.value = true;
  try {
    const params = {};
    if (selectedAcademicYear.value) {
      params.academic_year = selectedAcademicYear.value;
    }
    const res = await timetableGenService.getClassesCurriculumPlans(params);
    classesCurriculumPlans.value = res.data || [];
    if (classesCurriculumPlans.value.length > 0) {
      if (!selectedClassId.value || !classesCurriculumPlans.value.some(c => c.class_id === selectedClassId.value)) {
        onClassPlanSelected(classesCurriculumPlans.value[0].class_id);
      } else {
        onClassPlanSelected(selectedClassId.value);
      }
    } else {
      currentClassPlan.value = null;
    }
  } catch (err) {
    console.error('Error fetching class curriculum plans', err);
  } finally {
    loadingPlans.value = false;
  }
}

function onClassPlanSelected(classId) {
  selectedClassId.value = classId;
  const found = classesCurriculumPlans.value.find(c => c.class_id === classId);
  if (found) {
    currentClassPlan.value = JSON.parse(JSON.stringify(found));
    if (!currentClassPlan.value.subjects) currentClassPlan.value.subjects = [];
  } else {
    currentClassPlan.value = null;
  }
}

async function saveCurrentClassPlan() {
  if (!currentClassPlan.value || typeof timetableGenService.saveClassCurriculumPlan !== 'function') return;
  savingClassPlan.value = true;
  try {
    const payload = {
      min_hours_per_day: currentClassPlan.value.min_hours_per_day || 4,
      max_hours_per_day: currentClassPlan.value.max_hours_per_day || 6,
      subjects: currentClassPlan.value.subjects || []
    };
    const res = await timetableGenService.saveClassCurriculumPlan(currentClassPlan.value.class_id, payload);
    if (res.data) {
      currentClassPlan.value = JSON.parse(JSON.stringify(res.data));
      if (!currentClassPlan.value.subjects) currentClassPlan.value.subjects = [];
      const idx = classesCurriculumPlans.value.findIndex(c => c.class_id === currentClassPlan.value.class_id);
      if (idx !== -1) {
        classesCurriculumPlans.value[idx] = JSON.parse(JSON.stringify(res.data));
      }
    }
    $q.notify({
      type: 'positive',
      message: `Piano orario e limiti per ${currentClassPlan.value.class_name} salvati con successo!`
    });
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore durante il salvataggio del piano orario della classe'
    });
  } finally {
    savingClassPlan.value = false;
  }
}

function openInheritModal() {
  if (academicYears.value.length > 1) {
    inheritSourceYear.value = academicYears.value[1];
  } else if (academicYears.value.length > 0) {
    inheritSourceYear.value = academicYears.value[0];
  }
  inheritScope.value = selectedClassId.value ? 'single' : 'all';
  showInheritModal.value = true;
}

async function executeInherit() {
  inheriting.value = true;
  try {
    if (inheritScope.value === 'single' && currentClassPlan.value) {
      if (typeof timetableGenService.inheritClassCurriculumPlan !== 'function') return;
      const res = await timetableGenService.inheritClassCurriculumPlan(currentClassPlan.value.class_id, {
        source_academic_year: inheritSourceYear.value
      });
      if (res.data) {
        currentClassPlan.value = JSON.parse(JSON.stringify(res.data));
        if (!currentClassPlan.value.subjects) currentClassPlan.value.subjects = [];
        const idx = classesCurriculumPlans.value.findIndex(c => c.class_id === currentClassPlan.value.class_id);
        if (idx !== -1) {
          classesCurriculumPlans.value[idx] = JSON.parse(JSON.stringify(res.data));
        }
      }
      $q.notify({
        type: 'positive',
        message: `Dati ereditati dall'anno ${inheritSourceYear.value} per la classe ${currentClassPlan.value.class_name}`
      });
    } else {
      if (typeof timetableGenService.inheritAllClassesCurriculumPlans !== 'function') return;
      const res = await timetableGenService.inheritAllClassesCurriculumPlans({
        source_academic_year: inheritSourceYear.value
      });
      $q.notify({
        type: 'positive',
        message: res.data?.message || 'Dati ereditati con successo per tutte le classi'
      });
      await fetchClassesPlans();
    }
    showInheritModal.value = false;
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || 'Errore durante l\'ereditarietà dei dati dell\'anno precedente'
    });
  } finally {
    inheriting.value = false;
  }
}

function openAddSubjectModal() {
  newSubjectForm.value = {
    subject_id: '',
    teacher_id: null,
    hours_per_week: 3
  };
  showAddSubjectModal.value = true;
}

function addSubjectToCurrentPlan() {
  if (!newSubjectForm.value.subject_id || !currentClassPlan.value) return;
  if (!currentClassPlan.value.subjects) currentClassPlan.value.subjects = [];
  const existingSub = currentClassPlan.value.subjects.find(s => s.subject_id === newSubjectForm.value.subject_id);
  if (existingSub) {
    existingSub.hours_per_week = newSubjectForm.value.hours_per_week;
    if (newSubjectForm.value.teacher_id) {
      existingSub.teacher_id = newSubjectForm.value.teacher_id;
      const tObj = teachers.value.find(t => t.id === newSubjectForm.value.teacher_id || t.user_id === newSubjectForm.value.teacher_id);
      if (tObj) existingSub.teacher_name = `${tObj.last_name || ''} ${tObj.first_name || ''}`.trim();
    }
  } else {
    const subObj = subjects.value.find(s => s.id === newSubjectForm.value.subject_id);
    let tName = '';
    if (newSubjectForm.value.teacher_id) {
      const tObj = teachers.value.find(t => t.id === newSubjectForm.value.teacher_id || t.user_id === newSubjectForm.value.teacher_id);
      if (tObj) tName = `${tObj.last_name || ''} ${tObj.first_name || ''}`.trim();
    }
    currentClassPlan.value.subjects.push({
      subject_id: newSubjectForm.value.subject_id,
      subject_name: subObj ? subObj.name : 'Materia',
      hours_per_week: newSubjectForm.value.hours_per_week,
      teacher_id: newSubjectForm.value.teacher_id,
      teacher_name: tName
    });
  }
  currentClassPlan.value.total_hours_week = currentClassPlan.value.subjects.reduce((acc, s) => acc + (s.hours_per_week || 0), 0);
  showAddSubjectModal.value = false;
}

function removeSubjectFromCurrentPlan(idx) {
  if (!currentClassPlan.value || !currentClassPlan.value.subjects || !currentClassPlan.value.subjects[idx]) return;
  currentClassPlan.value.subjects.splice(idx, 1);
  currentClassPlan.value.total_hours_week = currentClassPlan.value.subjects.reduce((acc, s) => acc + (s.hours_per_week || 0), 0);
}

onMounted(async () => {
  await Promise.all([
    fetchRoomReqs(),
    fetchSubjects(),
    fetchTeachers(),
    fetchClasses(),
    fetchDesiderataWindow(),
    fetchAcademicYears(),
    fetchClassesPlans(),
    fetchTeachersQuickPreferences()
  ]);
  await fetchConstraints();
});
</script>
