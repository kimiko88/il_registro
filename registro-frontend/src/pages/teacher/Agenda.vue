<template>
  <q-page padding class="bg-slate-100 min-h-screen">
    <!-- Top Header & Toolbar -->
    <div class="row items-center justify-between q-mb-md gap-y-3">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">
          <q-icon name="event" color="primary" class="q-mr-sm" />
          {{ $t('agendaPage.title') }}
        </h1>
        <p class="text-subtitle1 text-slate-500 q-mt-xs q-mb-none">
          {{ $t('agendaPage.subtitle') }}
        </p>
      </div>

      <div class="row items-center q-gutter-sm">
        <q-btn-toggle
          v-model="viewMode"
          toggle-color="primary"
          unelevated
          dense
          no-caps
          class="bg-white border-2 border-slate-300 rounded-lg shadow-xs"
          :options="[
            { label: $t('agendaPage.month'), value: 'mese', icon: 'calendar_month' },
            { label: $t('agendaPage.week'), value: 'settimana', icon: 'view_week' },
            { label: $t('agendaPage.day'), value: 'giorno', icon: 'today' }
          ]"
          @update:model-value="onViewModeChange"
        />

        <q-select
          v-model="selectedClassFilter"
          :options="classOptions"
          :label="$t('agendaPage.filterClass')"
          outlined dense
          clearable
          emit-value
          map-options
          style="min-width: 240px"
          class="bg-white"
          @update:model-value="loadAgendaEvents"
        />

        <q-btn
          color="primary"
          unelevated
          icon="add"
          :label="$t('agendaPage.newEvent')"
          class="rounded-lg q-px-md shadow-xs font-bold"
          no-caps
          @click="openCreateDialog()"
        />

        <q-btn flat round icon="refresh" color="primary" :loading="agendaStore.loading" @click="loadAgendaEvents" />
      </div>
    </div>

    <!-- TYPE FILTER TOOLBAR -->
    <div class="q-px-md q-py-sm bg-white rounded-xl border-2 border-slate-300 shadow-xs q-mb-lg row items-center q-gutter-xs">
      <span class="text-caption text-slate-600 font-bold q-mr-xs">{{ $t('agendaPage.filterType') }}:</span>
      <q-chip
        clickable
        dense
        :color="selectedTypeFilter === 'tutti' ? 'primary' : 'grey-3'"
        :text-color="selectedTypeFilter === 'tutti' ? 'white' : 'dark'"
        @click="selectedTypeFilter = 'tutti'"
      >{{ $t('agendaPage.all') }}</q-chip>
      <q-chip
        clickable
        dense
        :color="selectedTypeFilter === 'compito' ? 'info' : 'grey-3'"
        :text-color="selectedTypeFilter === 'compito' ? 'white' : 'dark'"
        @click="selectedTypeFilter = 'compito'"
      >{{ $t('agendaPage.homework') }}</q-chip>
      <q-chip
        clickable
        dense
        :color="selectedTypeFilter === 'verifica' ? 'negative' : 'grey-3'"
        :text-color="selectedTypeFilter === 'verifica' ? 'white' : 'dark'"
        @click="selectedTypeFilter = 'verifica'"
      >{{ $t('agendaPage.test') }}</q-chip>
      <q-chip
        clickable
        dense
        :color="selectedTypeFilter === 'interrogazione' ? 'purple' : 'grey-3'"
        :text-color="selectedTypeFilter === 'interrogazione' ? 'white' : 'dark'"
        @click="selectedTypeFilter = 'interrogazione'"
      >{{ $t('agendaPage.oralTest') }}</q-chip>
      <q-chip
        clickable
        dense
        :color="selectedTypeFilter === 'avviso' ? 'warning' : 'grey-3'"
        :text-color="selectedTypeFilter === 'avviso' ? 'white' : 'dark'"
        @click="selectedTypeFilter = 'avviso'"
      >{{ $t('communicationsPage.title') }}</q-chip>
      <q-chip
        clickable
        dense
        :color="selectedTypeFilter === 'evento' ? 'positive' : 'grey-3'"
        :text-color="selectedTypeFilter === 'evento' ? 'white' : 'dark'"
        @click="selectedTypeFilter = 'evento'"
      >{{ $t('agendaPage.event') }}</q-chip>
      <q-chip
        clickable
        dense
        :color="selectedTypeFilter === 'altro' ? 'grey-7' : 'grey-3'"
        :text-color="selectedTypeFilter === 'altro' ? 'white' : 'dark'"
        @click="selectedTypeFilter = 'altro'"
      >Altro</q-chip>
    </div>

    <!-- VIEW 1: FULL MONTH VIEW (MESE) - GOOGLE CALENDAR STYLE -->
    <div v-if="viewMode === 'mese'" class="bg-white rounded-2xl gcal-border-grid shadow-md overflow-hidden">
      <!-- Month Navigation Controls Header -->
      <div class="bg-slate-100 border-b-2 border-slate-300 px-6 py-3.5 row items-center justify-between">
        <div class="row items-center gap-3">
          <q-btn unelevated color="primary" label="Oggi" no-caps class="rounded-lg font-bold px-4" @click="today" />
          <div class="row items-center">
            <q-btn flat round dense icon="chevron_left" color="slate-700" size="md" @click="prevMonth" />
            <q-btn flat round dense icon="chevron_right" color="slate-700" size="md" @click="nextMonth" />
          </div>
          <div class="text-h5 font-extrabold text-slate-800 tracking-tight capitalize ml-2">
            {{ currentMonthTitle }}
          </div>
        </div>

        <div class="text-caption font-bold text-slate-500 uppercase tracking-widest bg-slate-200 px-3 py-1 rounded-full">
          Vista Mensile
        </div>
      </div>

      <!-- 7-Column Days Header (Strict CSS Grid) -->
      <div class="grid-7-cols text-center font-bold text-xs text-slate-700 uppercase tracking-wider">
        <div v-for="(dayName, dIdx) in ['Lun', 'Mar', 'Mer', 'Gio', 'Ven', 'Sab', 'Dom']" :key="dayName" class="gcal-header-cell py-3" :class="{ 'bg-slate-200/50': dIdx >= 5 }">
          {{ dayName }}
        </div>
      </div>

      <!-- Month Days Grid Rows -->
      <div class="month-grid-container">
        <div
          v-for="(row, rIndex) in monthWeeks"
          :key="'week-' + rIndex"
          class="grid-7-cols min-h-[145px]"
        >
          <div
            v-for="(d, cIndex) in row"
            :key="d.iso"
            class="gcal-day-cell p-2 flex flex-col justify-start transition-colors group cursor-pointer hover:bg-blue-50/50"
            :class="{
              'today-cell': d.iso === isTodayIso,
              'bg-slate-100/70': !d.isCurrentMonth && d.iso !== isTodayIso,
              'bg-slate-50/40': d.isCurrentMonth && cIndex >= 5 && d.iso !== isTodayIso,
              'bg-blue-50/80': d.iso === isoSelectedDate && d.isCurrentMonth && d.iso !== isTodayIso
            }"
            @click="openCreateDialog(d.iso)"
          >
            <!-- Cell Top Row: Date Circle & Hover Plus Button -->
            <div class="row items-center justify-between mb-2">
              <span
                class="gcal-date-badge"
                :class="{
                  'today': d.iso === isTodayIso,
                  'current': d.isCurrentMonth && d.iso !== isTodayIso,
                  'other': !d.isCurrentMonth
                }"
              >
                {{ d.dateNumber }}
              </span>

              <span class="gcal-add-btn">
                + Nuovo
              </span>
            </div>

            <!-- Event Cards List -->
            <div class="flex-1 flex flex-col gap-1.5 overflow-hidden">
              <!-- All Day Pinned Badges -->
              <div
                v-for="ev in getAllDayEventsForDay(d.iso)"
                :key="'allday-' + ev.id"
                class="gcal-event-pill shadow-xs"
                :class="getTypeBgClass(ev.type)"
                @click.stop="openEditDialog(ev)"
              >
                <span class="shrink-0 mr-1">📌</span>
                <span class="truncate">{{ ev.title }} ({{ getFormattedClassName(ev) }})</span>
              </div>

              <!-- Hourly Timed Event Pills -->
              <div
                v-for="ev in getEventsForDay(d.iso).slice(0, 3 - getAllDayEventsForDay(d.iso).length)"
                :key="ev.id"
                class="gcal-event-pill shadow-xs"
                :class="getTypeBgClass(ev.type)"
                @click.stop="openEditDialog(ev)"
              >
                <span class="font-extrabold text-[10px] opacity-90 shrink-0 mr-1.5">{{ ev.start_time || '09:00' }}</span>
                <span class="font-bold truncate">{{ ev.title }}</span>
                <span class="text-[10px] opacity-85 truncate shrink-0 ml-1">({{ getFormattedClassName(ev) }})</span>
              </div>

              <!-- Overflow Counter -->
              <div
                v-if="getEventsForDay(d.iso).length + getAllDayEventsForDay(d.iso).length > 3"
                class="text-[11px] font-bold text-blue-700 px-1 pt-0.5 hover:underline cursor-pointer"
                @click.stop="selectDateIso(d.iso); viewMode = 'giorno'"
              >
                +{{ (getEventsForDay(d.iso).length + getAllDayEventsForDay(d.iso).length) - 3 }} altri
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- VIEW 2: WEEK VIEW (SETTIMANA) - GOOGLE CALENDAR GRID WITH PROMINENT BORDERS -->
    <div v-else-if="viewMode === 'settimana'" class="bg-white rounded-2xl gcal-border-grid shadow-md overflow-hidden">
      <!-- Week Header Controls -->
      <div class="bg-slate-100 border-b-2 border-slate-300 px-6 py-3.5 row items-center justify-between">
        <div class="row items-center gap-3">
          <q-btn unelevated color="primary" label="Questa Settimana" no-caps class="rounded-lg font-bold px-4" @click="thisWeek" />
          <div class="row items-center">
            <q-btn flat round dense icon="chevron_left" color="slate-700" size="md" @click="prevWeek" />
            <q-btn flat round dense icon="chevron_right" color="slate-700" size="md" @click="nextWeek" />
          </div>
          <div class="text-subtitle1 font-extrabold text-slate-800 ml-2">
            Settimana {{ weekTitle }}
          </div>
        </div>
        <div class="text-caption font-bold text-slate-500 uppercase tracking-widest bg-slate-200 px-3 py-1 rounded-full">
          Vista Settimanale
        </div>
      </div>

      <!-- Days Header Row (1 time col 65px + 7 day cols) -->
      <div class="grid-week-header text-center">
        <div class="gcal-header-cell row items-center justify-center text-[10px] font-mono font-bold text-slate-700 bg-slate-200/80">
          GMT+2
        </div>
        <div
          v-for="(d, dIdx) in weekDays"
          :key="'head-' + d.iso"
          class="gcal-header-cell py-2.5 cursor-pointer hover:bg-slate-200/60 transition-colors"
          :class="{
            'today-header-cell': d.iso === isTodayIso,
            'bg-blue-50/90': d.iso === isoSelectedDate && d.iso !== isTodayIso,
            'bg-slate-200/50': dIdx >= 5 && d.iso !== isoSelectedDate && d.iso !== isTodayIso
          }"
          @click="selectDateIso(d.iso)"
        >
          <div class="text-[11px] font-bold uppercase tracking-wider" :class="d.iso === isTodayIso ? 'text-blue-800 font-extrabold' : 'text-slate-700'">{{ d.dayName }}</div>
          <div
            class="inline-flex items-center justify-center w-7 h-7 rounded-full text-sm font-bold mt-0.5"
            :class="d.iso === isTodayIso ? 'bg-blue-700 text-white font-black shadow-sm' : 'text-slate-800'"
          >
            {{ d.dateNumber }}
          </div>
        </div>
      </div>

      <!-- Pinned All-Day Row for Week View -->
      <div class="grid-week-header bg-amber-50/60 min-h-[44px] border-b-2 border-amber-400">
        <div class="gcal-time-col p-2 text-center text-xs font-bold text-amber-900 row items-center justify-center bg-amber-100/90">
          <q-icon name="push_pin" size="14px" />
        </div>
        <div
          v-for="(d, dIdx) in weekDays"
          :key="'allday-' + d.iso"
          class="gcal-day-cell p-1 min-h-[44px] cursor-pointer"
          :class="{
            'today-slot': d.iso === isTodayIso,
            'bg-blue-50/40': d.iso === isoSelectedDate && d.iso !== isTodayIso,
            'bg-slate-100/40': dIdx >= 5 && d.iso !== isoSelectedDate && d.iso !== isTodayIso
          }"
          @click="openCreateDialog(d.iso)"
        >
          <div
            v-for="ev in getAllDayEventsForDay(d.iso)"
            :key="ev.id"
            class="gcal-event-pill shadow-xs"
            :class="getTypeBgClass(ev.type)"
            @click.stop="openEditDialog(ev)"
          >
            📌 {{ ev.title }} ({{ getFormattedClassName(ev) }})
          </div>
        </div>
      </div>

      <!-- Hourly Grid Slots (24 Hours 00:00 to 23:00) -->
      <div ref="weekScrollContainer" class="max-h-[620px] overflow-y-auto overflow-x-hidden relative">
        <!-- Google Calendar Red Real-Time Current Time Indicator -->
        <div
          v-if="isCurrentTimeVisible"
          class="gcal-now-indicator"
          :style="{ top: redLineTop + 'px' }"
        >
          <div class="gcal-now-time-badge">
            {{ formattedCurrentTime }}
          </div>
          <div class="gcal-now-line" />
        </div>

        <div
          v-for="hour in hoursList"
          :key="hour"
          class="grid-week-header min-h-[60px] relative gcal-hour-row"
        >
          <!-- Time Column -->
          <div class="gcal-time-col p-2 text-xs font-bold font-mono text-slate-700 row items-start justify-end pr-2 pt-0.5">
            <span>{{ hour }}</span>
          </div>

          <!-- 7 Day Columns for this hour -->
          <div
            v-for="(d, dIdx) in weekDays"
            :key="d.iso + '-' + hour"
            class="gcal-day-cell p-1.5 relative bg-white hover:bg-blue-50/30 transition-colors cursor-pointer group"
            :class="{
              'today-slot': d.iso === isTodayIso,
              'bg-blue-50/20': d.iso === isoSelectedDate && d.iso !== isTodayIso,
              'bg-slate-100/40': dIdx >= 5 && d.iso !== isoSelectedDate && d.iso !== isTodayIso
            }"
            @click="openCreateDialog(d.iso, hour)"
          >
            <div
              v-for="ev in getEventsForDayAndHour(d.iso, hour)"
              :key="ev.id"
              class="gcal-event-pill shadow-xs relative z-10"
              :class="getTypeBgClass(ev.type)"
              @click.stop="openEditDialog(ev)"
            >
              <div class="truncate">
                <div class="font-bold text-xs">{{ ev.title }}</div>
                <div class="text-[10px] opacity-90 row items-center justify-between mt-0.5">
                  <span>🕒 {{ ev.start_time || hour }} - {{ ev.end_time }}</span>
                  <span class="font-bold bg-white/20 px-1 rounded ml-1">
                    {{ getFormattedClassName(ev) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- VIEW 3: DAY VIEW (GIORNO) - GOOGLE CALENDAR GRID -->
    <div v-else-if="viewMode === 'giorno'" class="bg-white rounded-2xl gcal-border-grid shadow-md overflow-hidden">
      <!-- Day Header Controls -->
      <div class="bg-slate-100 border-b-2 border-slate-300 px-6 py-3.5 row items-center justify-between">
        <div class="row items-center gap-3">
          <q-btn unelevated color="primary" label="Oggi" no-caps class="rounded-lg font-bold px-4" @click="today" />
          <div class="row items-center">
            <q-btn flat round dense icon="chevron_left" color="slate-700" size="md" @click="prevDay" />
            <q-btn flat round dense icon="chevron_right" color="slate-700" size="md" @click="nextDay" />
          </div>
          <div class="text-subtitle1 font-extrabold text-slate-800 ml-2">
            {{ formattedSelectedDateFull }}
          </div>
        </div>
        <div class="text-caption font-bold text-slate-500 uppercase tracking-widest bg-slate-200 px-3 py-1 rounded-full">
          Vista Giornaliera
        </div>
      </div>

      <!-- Day Header Line -->
      <div class="row bg-slate-100 text-center py-2.5 items-center border-b-2 border-slate-300" :class="{ 'today-header-cell': isoSelectedDate === isTodayIso }">
        <div class="gcal-time-col text-[10px] text-slate-700 font-mono font-bold py-1" style="width: 65px; min-width: 65px;">GMT+2</div>
        <div class="col text-left pl-4 row items-center gap-2">
          <span class="text-sm font-bold uppercase tracking-wider" :class="isoSelectedDate === isTodayIso ? 'text-blue-800 font-extrabold' : 'text-slate-700'">{{ currentDayName }}</span>
          <span
            class="inline-flex items-center justify-center w-8 h-8 rounded-full text-base font-bold"
            :class="isoSelectedDate === isTodayIso ? 'bg-blue-700 text-white shadow-sm font-black' : 'bg-slate-200 text-slate-800'"
          >
            {{ currentDayNumber }}
          </span>
          <q-badge v-if="isoSelectedDate === isTodayIso" color="primary" label="Oggi" class="q-ml-sm text-weight-bold" />
        </div>
      </div>

      <!-- Pinned All-Day Header for Single Day -->
      <div class="row bg-amber-50/60 p-2 cursor-pointer border-b-2 border-amber-400" @click="openCreateDialog(isoSelectedDate)">
        <div class="gcal-time-col p-2 text-center text-xs font-bold text-amber-900 row items-center justify-center bg-amber-100/90" style="width: 65px; min-width: 65px;">
          <q-icon name="push_pin" size="14px" />
        </div>
        <div class="col pl-3">
          <div v-if="getAllDayEventsForDay(isoSelectedDate).length === 0" class="text-caption text-slate-400 italic py-1">
            + Clicca qui per aggiungere evento per tutta la giornata
          </div>
          <div v-else class="row q-gutter-xs">
            <div
              v-for="ev in getAllDayEventsForDay(isoSelectedDate)"
              :key="ev.id"
              class="gcal-event-pill shadow-xs cursor-pointer"
              :class="getTypeBgClass(ev.type)"
              @click.stop="openEditDialog(ev)"
            >
              <span class="font-bold">📌 {{ ev.title }} ({{ getFormattedClassName(ev) }})</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Hourly Scrollable Body Grid for Day View -->
      <div ref="dayScrollContainer" class="max-h-[620px] overflow-y-auto overflow-x-hidden relative">
        <!-- Google Calendar Red Real-Time Current Time Indicator -->
        <div
          v-if="isCurrentTimeVisible"
          class="gcal-now-indicator"
          :style="{ top: redLineTop + 'px' }"
        >
          <div class="gcal-now-time-badge">
            {{ formattedCurrentTime }}
          </div>
          <div class="gcal-now-line" />
        </div>

        <div
          v-for="hour in hoursList"
          :key="hour"
          class="row min-h-[60px] relative gcal-hour-row"
        >
          <!-- Time Label -->
          <div class="gcal-time-col p-2 text-xs font-bold font-mono text-slate-700 row items-start justify-end pr-2 pt-0.5" style="width: 65px; min-width: 65px;">
            <span>{{ hour }}</span>
          </div>

          <!-- Single Day Slot for this hour -->
          <div
            class="col gcal-day-cell p-2 relative bg-white hover:bg-blue-50/20 transition-colors cursor-pointer group row items-center gap-2"
            @click="openCreateDialog(isoSelectedDate, hour)"
          >
            <div
              v-for="ev in getEventsForDayAndHour(isoSelectedDate, hour)"
              :key="ev.id"
              class="gcal-event-pill shadow-xs cursor-pointer relative z-10 min-w-[280px] max-w-[440px] p-3"
              :class="getTypeBgClass(ev.type)"
              @click.stop="openEditDialog(ev)"
            >
              <div class="w-full">
                <div class="row items-center justify-between">
                  <span class="font-bold text-sm text-white">{{ ev.title }}</span>
                  <q-chip size="xs" color="white" text-color="dark" class="font-bold">
                    {{ getTypeLabel(ev.type) }}
                  </q-chip>
                </div>
                <div v-if="ev.description" class="text-xs text-white/90 mt-1">{{ ev.description }}</div>
                <div class="row items-center justify-between text-xs text-white/80 mt-2 border-t border-white/20 pt-1">
                  <span>🕒 {{ ev.start_time || hour }} - {{ ev.end_time }}</span>
                  <span class="font-bold bg-white/20 px-2 py-0.5 rounded">{{ getFormattedClassName(ev) }}</span>
                </div>
              </div>
            </div>

            <div
              v-if="getEventsForDayAndHour(isoSelectedDate, hour).length === 0"
              class="text-xs text-slate-400 italic opacity-0 group-hover:opacity-100 transition-opacity select-none"
            >
              + Clicca per aggiungere evento alle {{ hour }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Action Button for Mobile -->
    <q-page-sticky position="bottom-right" :offset="[18, 18]" class="lt-md">
      <q-btn round color="primary" icon="add" size="lg" class="shadow-lg" @click="openCreateDialog()" />
    </q-page-sticky>

    <!-- Create / Edit / View Event Dialog -->
    <q-dialog v-model="dialogVisible">
      <q-card style="width: min(600px, 95vw)" class="rounded-xl overflow-hidden shadow-24 border-slate-300">
        <q-card-section class="bg-primary text-white row items-center justify-between q-py-md">
          <div class="text-h6 text-weight-bold">
            <q-icon :name="isEditMode ? (isCurrentAuthor ? 'edit_calendar' : 'info') : 'event_available'" class="q-mr-xs" />
            {{ !isEditMode ? 'Nuovo Evento Agenda' : (isCurrentAuthor ? 'Modifica Evento Agenda' : 'Dettagli Evento Agenda') }}
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <!-- Read-Only Notice for Non-Authors -->
        <div v-if="isEditMode && !isCurrentAuthor" class="bg-amber-50 border-b border-amber-200 q-pa-sm text-amber-900 text-caption row items-center gap-2">
          <q-icon name="info" size="18px" color="amber-9" />
          <span>Creato da <strong>{{ currentEventTeacherName }}</strong>. Solo l'autore dell'evento può apportare modifiche o eliminarlo.</span>
        </div>

        <q-form ref="formRef" @submit.prevent="saveEvent">
          <q-card-section class="q-pa-md">
            <!-- Title Input -->
            <q-input
              v-model="form.title"
              label="Titolo Evento *"
              outlined
              dense
              tabindex="1"
              class="q-mb-md"
              :readonly="isEditMode && !isCurrentAuthor"
              :rules="[val => !!val || 'Il titolo è obbligatorio']"
            />

            <!-- Description Input -->
            <q-input
              v-model="form.description"
              label="Descrizione / Dettagli"
              outlined
              dense
              type="textarea"
              rows="3"
              tabindex="2"
              class="q-mb-md"
              :readonly="isEditMode && !isCurrentAuthor"
            />

            <!-- Row 1: Event Type & Class Select -->
            <div class="row q-col-gutter-md q-mb-md">
              <div class="col-12 col-sm-6">
                <q-select
                  v-model="form.type"
                  :options="typeOptions"
                  label="Tipo Evento *"
                  outlined dense
                  emit-value
                  map-options
                  tabindex="3"
                  :disable="isEditMode && !isCurrentAuthor"
                />
              </div>
              <div class="col-12 col-sm-6">
                <q-select
                  v-model="form.class_id"
                  :options="classOptions"
                  label="Classe Destinataria *"
                  outlined dense
                  emit-value
                  map-options
                  tabindex="4"
                  :disable="isEditMode && !isCurrentAuthor"
                  :rules="[val => !!val || 'Seleziona una classe']"
                />
              </div>
            </div>

            <!-- All-Day Toggle -->
            <div class="row items-center justify-between bg-blue-50/60 q-pa-sm rounded-lg border border-blue-200 q-mb-md">
              <div>
                <div class="text-subtitle2 text-slate-800 font-semibold">📌 Tutto il Giorno</div>
                <div class="text-caption text-slate-500">L'evento impegna l'intera giornata senza orario specifico</div>
              </div>
              <q-toggle v-model="form.all_day" color="primary" :disable="isEditMode && !isCurrentAuthor" />
            </div>

            <!-- Row 2: Date & Time Inputs -->
            <div class="row q-col-gutter-md q-mb-md">
              <div class="col-12" :class="form.all_day ? 'col-sm-12' : 'col-sm-4'">
                <q-input
                  v-model="form.date"
                  type="date"
                  label="Data *"
                  outlined dense
                  tabindex="5"
                  :readonly="isEditMode && !isCurrentAuthor"
                  :rules="[val => !!val || 'La data è obbligatoria']"
                />
              </div>
              <div v-if="!form.all_day" class="col-6 col-sm-4">
                <q-input
                  v-model="form.start_time"
                  type="time"
                  label="Ora Inizio"
                  outlined dense
                  tabindex="6"
                  :readonly="isEditMode && !isCurrentAuthor"
                />
              </div>
              <div v-if="!form.all_day" class="col-6 col-sm-4">
                <q-input
                  v-model="form.end_time"
                  type="time"
                  label="Ora Fine"
                  outlined dense
                  tabindex="7"
                  :readonly="isEditMode && !isCurrentAuthor"
                />
              </div>
            </div>

            <!-- Toggle: Visible to Students -->
            <div class="row items-center justify-between bg-slate-50 q-pa-md rounded-lg border border-slate-200">
              <div>
                <div class="text-subtitle2 text-slate-700">Visibile agli Studenti & Genitori</div>
                <div class="text-caption text-slate-500">Se disattivato, l'evento sarà visibile solo ai docenti della classe</div>
              </div>
              <q-toggle v-model="form.visible_to_students" color="primary" :disable="isEditMode && !isCurrentAuthor" />
            </div>
          </q-card-section>

          <q-card-actions align="between" class="q-pa-md bg-slate-50 border-t border-slate-100">
            <q-btn
              v-if="isEditMode && isCurrentAuthor"
              color="negative"
              flat
              icon="delete"
              label="Elimina"
              no-caps
              @click="confirmDelete"
            />
            <div v-else />

            <div class="row q-gutter-sm">
              <q-btn flat label="Chiudi" no-caps v-close-popup />
              <q-btn
                v-if="!isEditMode || isCurrentAuthor"
                type="submit"
                color="primary"
                unelevated
                :label="isEditMode ? 'Salva Modifiche' : 'Crea Evento'"
                :loading="saving"
                no-caps
              />
            </div>
          </q-card-actions>
        </q-form>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useQuasar, date as qdate } from 'quasar'
import { useAgendaStore } from '@/stores/agenda'
import { useClassesStore } from '@/stores/classes'
import { useAuthStore } from '@/stores/auth'

const $q = useQuasar()
const agendaStore = useAgendaStore()
const classesStore = useClassesStore()
const authStore = useAuthStore()

const formRef = ref(null)
const dialogVisible = ref(false)
const isEditMode = ref(false)
const editId = ref(null)
const currentEventTeacherId = ref(null)
const currentEventTeacherName = ref('')
const saving = ref(false)

const weekScrollContainer = ref(null)
const dayScrollContainer = ref(null)

const viewMode = ref('mese') // 'mese', 'settimana', 'giorno'
const selectedDate = ref(qdate.formatDate(new Date(), 'YYYY/MM/DD'))
const selectedClassFilter = ref(null)
const selectedTypeFilter = ref('tutti')

// Real-Time Current Time Indicator
const redLineTop = ref(0)
const nowTime = ref(new Date())
let timerInterval = null

function updateCurrentTime() {
  nowTime.value = new Date()
  const now = nowTime.value
  const currentHour = now.getHours()
  const currentMinute = now.getMinutes()

  const container = viewMode.value === 'settimana' ? weekScrollContainer.value : dayScrollContainer.value
  if (container) {
    const rowElements = container.querySelectorAll('.gcal-hour-row')
    if (rowElements && rowElements[currentHour]) {
      const rowEl = rowElements[currentHour]
      const rowTop = rowEl.offsetTop
      const rowHeight = rowEl.offsetHeight || 60
      redLineTop.value = rowTop + Math.round((currentMinute / 60) * rowHeight)
      return
    }
  }
  redLineTop.value = (currentHour * 60) + currentMinute
}

const formattedCurrentTime = computed(() => {
  const now = nowTime.value
  const h = String(now.getHours()).padStart(2, '0')
  const m = String(now.getMinutes()).padStart(2, '0')
  return `${h}:${m}`
})

const isCurrentTimeVisible = computed(() => {
  const todayStr = formatYMD(nowTime.value)
  if (viewMode.value === 'giorno') {
    return isoSelectedDate.value === todayStr
  }
  if (viewMode.value === 'settimana') {
    return weekDays.value.some(d => d.iso === todayStr)
  }
  return false
})

// Helper function to format Date object into local YYYY-MM-DD string cleanly without timezone shifts
function formatYMD(d) {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 24 Hours list from 00:00 to 23:00
const hoursList = Array.from({ length: 24 }, (_, i) => `${String(i).padStart(2, '0')}:00`)

const form = reactive({
  title: '',
  description: '',
  type: 'compito',
  all_day: false,
  class_id: '',
  date: formatYMD(new Date()),
  start_time: '09:00',
  end_time: '10:00',
  visible_to_students: true
})

const typeOptions = [
  { label: 'Compito per casa', value: 'compito' },
  { label: 'Compito in classe', value: 'verifica' },
  { label: 'Interrogazione', value: 'interrogazione' },
  { label: 'Avviso di Classe', value: 'avviso' },
  { label: 'Evento / Uscita Didattica', value: 'evento' },
  { label: 'Altro', value: 'altro' }
]

const classOptions = computed(() => {
  return classesStore.classes.map(c => {
    let nameText = c.name || `Classe ${c.id}`
    if (c.section && !nameText.endsWith(c.section)) {
      nameText += c.section
    }
    if (c.articolazione) {
      nameText += ` - ${c.articolazione}`
    }
    return {
      label: nameText,
      value: c.id
    }
  })
})

function getFormattedClassName(ev) {
  if (!ev) return ''
  const classObj = classesStore.classes.find(c => c.id === ev.class_id)
  if (classObj) {
    let nameText = classObj.name || ''
    if (classObj.section && !nameText.endsWith(classObj.section)) {
      nameText += classObj.section
    }
    if (classObj.articolazione) {
      nameText += ` - ${classObj.articolazione}`
    }
    return nameText
  }
  return ev.class_name || 'Classe'
}

const isCurrentAuthor = computed(() => {
  if (!isEditMode.value) return true // Creating new event
  const currentUserId = authStore.user?.id
  const eventTeacherId = currentEventTeacherId.value
  if (!eventTeacherId || !currentUserId) return false
  return eventTeacherId === currentUserId
})

const isoSelectedDate = computed(() => {
  if (!selectedDate.value) return formatYMD(new Date())
  return selectedDate.value.replace(/\//g, '-')
})

const isTodayIso = computed(() => formatYMD(nowTime.value))

const filteredEvents = computed(() => {
  return agendaStore.events.filter(ev => {
    const matchesClass = !selectedClassFilter.value || ev.class_id === selectedClassFilter.value
    const matchesType = selectedTypeFilter.value === 'tutti' || ev.type === selectedTypeFilter.value
    return matchesClass && matchesType
  })
})

const eventsByDateMap = computed(() => {
  const map = {}
  for (const ev of filteredEvents.value) {
    const dateStr = ev.date ? ev.date.substring(0, 10) : ''
    if (!dateStr) continue
    if (!map[dateStr]) {
      map[dateStr] = { allDay: [], timed: [], byHour: {} }
    }
    if (ev.all_day) {
      map[dateStr].allDay.push(ev)
    } else {
      map[dateStr].timed.push(ev)
      const startTime = ev.start_time || '09:00'
      const hourNum = parseInt(startTime.split(':')[0], 10)
      if (!map[dateStr].byHour[hourNum]) {
        map[dateStr].byHour[hourNum] = []
      }
      map[dateStr].byHour[hourNum].push(ev)
    }
  }
  return map
})

function getAllDayEventsForDay(dayIso) {
  return eventsByDateMap.value[dayIso]?.allDay || []
}

function getEventsForDay(dayIso) {
  return eventsByDateMap.value[dayIso]?.timed || []
}

function getEventsForDayAndHour(dayIso, hourStr) {
  const hourNum = parseInt(hourStr.split(':')[0], 10)
  return eventsByDateMap.value[dayIso]?.byHour[hourNum] || []
}

function scrollToCurrentTime() {
  nextTick(() => {
    updateCurrentTime()
    const scrollPos = Math.max(0, redLineTop.value - 220)
    if (weekScrollContainer.value) {
      weekScrollContainer.value.scrollTop = scrollPos
    }
    if (dayScrollContainer.value) {
      dayScrollContainer.value.scrollTop = scrollPos
    }
  })
}

const currentMonthTitle = computed(() => {
  if (!selectedDate.value) return ''
  const d = new Date(isoSelectedDate.value)
  const monthNames = ['Gennaio', 'Febbraio', 'Marzo', 'Aprile', 'Maggio', 'Giugno', 'Luglio', 'Agosto', 'Settembre', 'Ottobre', 'Novembre', 'Dicembre']
  return `${monthNames[d.getMonth()]} ${d.getFullYear()}`
})

const currentDayName = computed(() => {
  const d = new Date(isoSelectedDate.value)
  const dayNames = ['Domenica', 'Lunedì', 'Martedì', 'Mercoledì', 'Giovedì', 'Venerdì', 'Sabato']
  return dayNames[d.getDay()]
})

const currentDayNumber = computed(() => {
  const d = new Date(isoSelectedDate.value)
  return d.getDate()
})

const formattedSelectedDateFull = computed(() => {
  if (!selectedDate.value) return ''
  const d = new Date(isoSelectedDate.value)
  const dayNames = ['Domenica', 'Lunedì', 'Martedì', 'Mercoledì', 'Giovedì', 'Venerdì', 'Sabato']
  const monthNames = ['Gennaio', 'Febbraio', 'Marzo', 'Aprile', 'Maggio', 'Giugno', 'Luglio', 'Agosto', 'Settembre', 'Ottobre', 'Novembre', 'Dicembre']
  return `${dayNames[d.getDay()]} ${d.getDate()} ${monthNames[d.getMonth()]} ${d.getFullYear()}`
})

// Full Month Days Grid (7 columns x 5-6 rows) for Month View
const monthWeeks = computed(() => {
  if (!selectedDate.value) return []
  const selDate = new Date(isoSelectedDate.value)
  const year = selDate.getFullYear()
  const month = selDate.getMonth()

  const firstDayOfMonth = new Date(year, month, 1)
  const lastDayOfMonth = new Date(year, month + 1, 0)

  let startDayOfWeek = firstDayOfMonth.getDay()
  let offset = startDayOfWeek === 0 ? 6 : startDayOfWeek - 1

  const allDays = []

  // Prev month filler days
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  for (let i = offset - 1; i >= 0; i--) {
    const d = new Date(year, month - 1, prevMonthLastDay - i)
    allDays.push({
      dateObj: d,
      iso: formatYMD(d),
      dateNumber: d.getDate(),
      isCurrentMonth: false
    })
  }

  // Current month days
  for (let i = 1; i <= lastDayOfMonth.getDate(); i++) {
    const d = new Date(year, month, i)
    allDays.push({
      dateObj: d,
      iso: formatYMD(d),
      dateNumber: i,
      isCurrentMonth: true
    })
  }

  // Next month filler days
  const remaining = 7 - (allDays.length % 7)
  if (remaining < 7) {
    for (let i = 1; i <= remaining; i++) {
      const d = new Date(year, month + 1, i)
      allDays.push({
        dateObj: d,
        iso: formatYMD(d),
        dateNumber: i,
        isCurrentMonth: false
      })
    }
  }

  // Group into weeks (rows of 7 days)
  const weeks = []
  for (let i = 0; i < allDays.length; i += 7) {
    weeks.push(allDays.slice(i, i + 7))
  }
  return weeks
})

const weekDays = computed(() => {
  if (!selectedDate.value) return []
  const selDate = new Date(isoSelectedDate.value)
  const day = selDate.getDay() // 0 = Sun, 1 = Mon...
  const diffToMon = day === 0 ? -6 : 1 - day
  const mon = new Date(selDate)
  mon.setDate(selDate.getDate() + diffToMon)

  const days = []
  const dayNames = ['Lun', 'Mar', 'Mer', 'Gio', 'Ven', 'Sab', 'Dom']
  for (let i = 0; i < 7; i++) {
    const d = new Date(mon)
    d.setDate(mon.getDate() + i)
    days.push({
      dateObj: d,
      iso: formatYMD(d),
      dayName: dayNames[i],
      dateNumber: d.getDate()
    })
  }
  return days
})

const weekTitle = computed(() => {
  if (weekDays.value.length === 0) return ''
  const first = weekDays.value[0].dateObj
  const last = weekDays.value[6].dateObj
  const monthNames = ['Gen', 'Feb', 'Mar', 'Apr', 'Mag', 'Giu', 'Lug', 'Ago', 'Set', 'Ott', 'Nov', 'Dic']
  return `${first.getDate()} ${monthNames[first.getMonth()]} - ${last.getDate()} ${monthNames[last.getMonth()]} ${last.getFullYear()}`
})

onMounted(async () => {
  updateCurrentTime()
  timerInterval = setInterval(updateCurrentTime, 60000)
  await classesStore.fetchAssignedClasses().catch(() => classesStore.fetchClasses())
  await loadAgendaEvents()
  scrollToCurrentTime()
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

function onViewModeChange(mode) {
  if (mode === 'settimana' || mode === 'giorno') {
    scrollToCurrentTime()
  }
}

async function loadAgendaEvents() {
  const params = {}
  if (selectedClassFilter.value) {
    params.class_id = selectedClassFilter.value
  }
  await agendaStore.fetchAgenda(params).catch(() => {})
}

function selectDateIso(isoStr) {
  selectedDate.value = isoStr.replace(/-/g, '/')
}

function nextMonth() {
  const d = new Date(isoSelectedDate.value)
  d.setMonth(d.getMonth() + 1)
  selectedDate.value = qdate.formatDate(d, 'YYYY/MM/DD')
}

function prevMonth() {
  const d = new Date(isoSelectedDate.value)
  d.setMonth(d.getMonth() - 1)
  selectedDate.value = qdate.formatDate(d, 'YYYY/MM/DD')
}

function nextWeek() {
  const d = new Date(isoSelectedDate.value)
  d.setDate(d.getDate() + 7)
  selectedDate.value = qdate.formatDate(d, 'YYYY/MM/DD')
  scrollToCurrentTime()
}

function prevWeek() {
  const d = new Date(isoSelectedDate.value)
  d.setDate(d.getDate() - 7)
  selectedDate.value = qdate.formatDate(d, 'YYYY/MM/DD')
  scrollToCurrentTime()
}

function thisWeek() {
  selectedDate.value = qdate.formatDate(new Date(), 'YYYY/MM/DD')
  scrollToCurrentTime()
}

function nextDay() {
  const d = new Date(isoSelectedDate.value)
  d.setDate(d.getDate() + 1)
  selectedDate.value = qdate.formatDate(d, 'YYYY/MM/DD')
  scrollToCurrentTime()
}

function prevDay() {
  const d = new Date(isoSelectedDate.value)
  d.setDate(d.getDate() - 1)
  selectedDate.value = qdate.formatDate(d, 'YYYY/MM/DD')
  scrollToCurrentTime()
}

function today() {
  selectedDate.value = qdate.formatDate(new Date(), 'YYYY/MM/DD')
  scrollToCurrentTime()
}

function openCreateDialog(isoDate, hourStr) {
  isEditMode.value = false
  editId.value = null
  currentEventTeacherId.value = authStore.user?.id
  currentEventTeacherName.value = `${authStore.user?.first_name || ''} ${authStore.user?.last_name || ''}`
  form.title = ''
  form.description = ''
  form.type = 'compito'
  form.all_day = false
  form.class_id = classOptions.value[0]?.value || ''
  form.date = isoDate || isoSelectedDate.value
  form.start_time = hourStr || '09:00'
  
  if (hourStr) {
    const h = parseInt(hourStr.split(':')[0], 10)
    const nextH = (h + 1 < 24) ? String(h + 1).padStart(2, '0') : '23'
    form.end_time = `${nextH}:00`
  } else {
    form.end_time = '10:00'
  }
  
  form.visible_to_students = true
  dialogVisible.value = true
}

function openEditDialog(ev) {
  isEditMode.value = true
  editId.value = ev.id
  currentEventTeacherId.value = ev.teacher_id
  currentEventTeacherName.value = ev.teacher_name || 'Docente'
  form.title = ev.title || ''
  form.description = ev.description || ''
  form.type = ev.type || 'compito'
  form.all_day = !!ev.all_day
  form.class_id = ev.class_id || ''
  form.date = ev.date ? ev.date.substring(0, 10) : isoSelectedDate.value
  form.start_time = ev.start_time || '09:00'
  form.end_time = ev.end_time || '10:00'
  form.visible_to_students = ev.visible_to_students !== false
  dialogVisible.value = true
}

async function saveEvent() {
  if (!isCurrentAuthor.value) return

  if (!form.title || !form.class_id || !form.date) {
    $q.notify({ type: 'warning', message: 'Compila tutti i campi obbligatori' })
    return
  }

  if (!form.all_day) {
    if (!form.start_time || !form.end_time) {
      $q.notify({ type: 'warning', message: 'Inserisci sia l\'ora di inizio che l\'ora di fine' })
      return
    }
    if (form.end_time <= form.start_time) {
      $q.notify({ type: 'warning', message: 'L\'ora di fine deve essere successiva all\'ora di inizio' })
      return
    }
  }

  saving.value = true
  try {
    const payload = {
      title: form.title,
      description: form.description,
      type: form.type,
      all_day: form.all_day,
      class_id: form.class_id,
      date: form.date,
      start_time: form.all_day ? '00:00' : form.start_time,
      end_time: form.all_day ? '23:59' : form.end_time,
      visible_to_students: form.visible_to_students
    }

    if (isEditMode.value) {
      await agendaStore.updateEvent(editId.value, payload)
      $q.notify({ type: 'positive', message: 'Evento modificato con successo' })
    } else {
      await agendaStore.createEvent(payload)
      $q.notify({ type: 'positive', message: 'Evento creato con successo' })
    }
    dialogVisible.value = false
    await loadAgendaEvents()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante il salvataggio' })
  } finally {
    saving.value = false
  }
}

async function confirmDelete() {
  if (!isCurrentAuthor.value) return

  $q.dialog({
    title: 'Conferma Eliminazione',
    message: `Sei sicuro di voler eliminare l'evento "${form.title || 'selezionato'}" dall'agenda?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    saving.value = true
    try {
      await agendaStore.deleteEvent(editId.value)
      $q.notify({ type: 'positive', message: 'Evento eliminato con successo' })
      dialogVisible.value = false
      await loadAgendaEvents()
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || 'Errore durante l\'eliminazione' })
    } finally {
      saving.value = false
    }
  })
}

function getTypeBgClass(type) {
  switch (type) {
    case 'compito': return 'gcal-bg-compito'
    case 'verifica': return 'gcal-bg-verifica'
    case 'interrogazione': return 'gcal-bg-interrogazione'
    case 'avviso': return 'gcal-bg-avviso'
    case 'evento': return 'gcal-bg-evento'
    case 'altro': return 'gcal-bg-altro'
    default: return 'gcal-bg-compito'
  }
}

function getTypeLabel(type) {
  switch (type) {
    case 'compito': return 'Compito per casa'
    case 'verifica': return 'Compito in classe'
    case 'interrogazione': return 'Interrogazione'
    case 'avviso': return 'Avviso di classe'
    case 'evento': return 'Evento / Uscita'
    case 'altro': return 'Altro'
    default: return type
  }
}
</script>

<style scoped>
/* Explicit Google Calendar Table & Grid Styles */
.gcal-border-grid {
  border: 2px solid #94a3b8 !important;
}

.grid-7-cols {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
}

.grid-week-header {
  display: grid;
  grid-template-columns: 65px repeat(7, minmax(0, 1fr));
}

.gcal-header-cell {
  border-right: 1px solid #cbd5e1 !important;
  border-bottom: 2px solid #94a3b8 !important;
  background-color: #f8fafc;
}

.gcal-header-cell:last-child {
  border-right: none !important;
}

.gcal-day-cell {
  border-right: 1px solid #cbd5e1 !important;
  border-bottom: 1px solid #cbd5e1 !important;
  background-color: #ffffff;
}

.gcal-day-cell:last-child {
  border-right: none !important;
}

.gcal-hour-row {
  height: 60px;
  min-height: 60px;
  box-sizing: border-box;
  border-bottom: 1px solid #cbd5e1 !important;
}

.gcal-time-col {
  border-right: 2px solid #94a3b8 !important;
  background-color: #f8fafc;
}

/* Today Highlighting */
.today-cell {
  background-color: #eff6ff !important;
  box-shadow: inset 0 0 0 2px #2563eb !important;
  position: relative;
  z-index: 2;
}

.today-header-cell {
  background-color: #dbeafe !important;
  border-bottom: 3px solid #2563eb !important;
}

.today-slot {
  background-color: #f0f7ff !important;
  border-left: 1px solid #bfdbfe !important;
  border-right: 1px solid #bfdbfe !important;
}

/* Date Circle Badges */
.gcal-date-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 700;
}

.gcal-date-badge.today {
  background-color: #1d4ed8 !important;
  color: #ffffff !important;
  font-weight: 900 !important;
  box-shadow: 0 2px 5px rgba(29, 78, 216, 0.4);
}

.gcal-date-badge.current {
  color: #1e293b;
}

.gcal-date-badge.other {
  color: #94a3b8;
}

.gcal-add-btn {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  opacity: 0;
  transition: opacity 0.2s;
}

.gcal-day-cell:hover .gcal-add-btn {
  opacity: 1;
}

/* Event Pill Styles with Vibrant Background Colors */
.gcal-event-pill {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 11px;
  line-height: 1.3;
  font-weight: 600;
  color: #ffffff !important;
  box-shadow: 0 1px 3px rgba(0,0,0,0.15);
  cursor: pointer;
  margin-bottom: 3px;
  transition: transform 0.15s, box-shadow 0.15s, filter 0.15s;
}

.gcal-event-pill:hover {
  transform: translateY(-1px);
  box-shadow: 0 3px 6px rgba(0,0,0,0.2);
  filter: brightness(1.08);
}

/* Event Types Colors */
.gcal-bg-compito {
  background-color: #2563eb !important; /* Royal Blue */
}

.gcal-bg-verifica {
  background-color: #dc2626 !important; /* Red */
}

.gcal-bg-interrogazione {
  background-color: #9333ea !important; /* Purple */
}

.gcal-bg-avviso {
  background-color: #d97706 !important; /* Amber */
}

.gcal-bg-evento {
  background-color: #059669 !important; /* Emerald */
}

.gcal-bg-altro {
  background-color: #475569 !important; /* Slate */
}

/* Real-Time Google Calendar Red Line */
.gcal-now-indicator {
  position: absolute;
  left: 0;
  right: 0;
  z-index: 30;
  display: flex;
  align-items: center;
  pointer-events: none;
  transition: top 0.2s ease-out;
}

.gcal-now-time-badge {
  width: 63px;
  background-color: #ef4444;
  color: #ffffff;
  font-size: 11px;
  font-family: monospace;
  font-weight: 800;
  text-align: center;
  padding: 2px 0;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(239, 68, 68, 0.5);
  margin-left: 1px;
  flex-shrink: 0;
}

.gcal-now-line {
  flex: 1;
  height: 2px;
  background-color: #ef4444;
  box-shadow: 0 0 4px rgba(239, 68, 68, 0.6);
}
</style>
