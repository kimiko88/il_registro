<template>
  <q-page class="q-pa-md q-pa-lg-xl strike-management-page">
    <!-- Header Hero Section -->
    <div class="row items-center justify-between q-mb-lg gap-md">
      <div class="col-12 col-md-7">
        <div class="row items-center q-gutter-sm q-mb-xs">
          <q-badge color="indigo-7" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="campaign" size="14px" class="q-mr-xs" />
            {{ t('strikeManagement.badge') }}
          </q-badge>
          <q-badge color="teal-8" text-color="white" class="q-px-sm q-py-xs text-weight-bold text-caption rounded-borders">
            <q-icon name="gavel" size="14px" class="q-mr-xs" />
            {{ t('strikeManagement.legalAccord') }}
          </q-badge>
        </div>
        <h1 class="text-h4 text-weight-bolder text-slate-800 q-my-none flex items-center">
          <q-icon name="ballot" color="primary" class="q-mr-sm" size="36px" />
          {{ t('strikeManagement.title') }}
        </h1>
        <div class="text-subtitle1 text-slate-500 q-mt-xs">
          {{ t('strikeManagement.subtitle') }}
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="col-12 col-md-5 row items-center justify-end q-gutter-sm">
        <q-btn
          color="primary"
          icon="add"
          :label="t('strikeManagement.newNoticeBtn')"
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
          <q-tooltip>{{ t('strikeManagement.refreshTooltip') }}</q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Legal & Operating Guideline Banner -->
    <div class="q-mb-lg bg-indigo-50 border border-indigo-100 rounded-xl q-pa-md row items-center justify-between">
      <div class="row items-center gap-sm">
        <q-avatar size="40px" color="indigo-600" text-color="white" icon="info" />
        <div>
          <div class="text-weight-bold text-slate-800">
            {{ t('strikeManagement.legalBannerTitle') }}
          </div>
          <div class="text-caption text-slate-600">
            {{ t('strikeManagement.legalBannerSubtitle') }}
          </div>
        </div>
      </div>
      <div class="row items-center q-gutter-xs">
        <q-btn
          flat
          dense
          color="primary"
          icon="co_present"
          :label="t('strikeManagement.goToStaffAttendance')"
          to="/ata/attendance"
          class="text-weight-bold"
        />
      </div>
    </div>

    <!-- Active Strike Notice Selector / Carousel -->
    <div v-if="notices.length === 0 && !loadingNotices" class="q-pa-xl text-center bg-white rounded-xl border border-slate-200 shadow-sm">
      <q-icon name="event_busy" size="64px" color="slate-300" class="q-mb-md" />
      <div class="text-h6 text-weight-bold text-slate-700">{{ t('strikeManagement.noNoticesTitle') }}</div>
      <div class="text-body2 text-slate-500 q-mt-xs q-mb-md">
        {{ t('strikeManagement.noNoticesDesc') }}
      </div>
      <q-btn
        color="primary"
        icon="add"
        :label="t('strikeManagement.createNoticeBtn')"
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
            {{ t('strikeManagement.registeredNoticesTitle', { count: notices.length }) }}
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
                    {{ item.is_expired ? t('strikeManagement.statusExpired') : t('strikeManagement.statusOpen') }}
                  </q-badge>

                  <q-btn
                    flat
                    round
                    dense
                    color="grey-6"
                    icon="delete_outline"
                    @click.stop="confirmDeleteNotice(item)"
                  >
                    <q-tooltip>{{ t('strikeManagement.deleteNoticeTooltip') }}</q-tooltip>
                  </q-btn>
                </div>

                <div class="text-subtitle1 text-weight-bold text-slate-900 line-clamp-2 q-mt-xs">
                  {{ item.title }}
                </div>

                <div class="text-caption text-slate-500 q-mt-xs">
                  {{ t('strikeManagement.proclaimedBy') }} <span class="text-weight-bold text-slate-700">{{ item.proclaimed_by || t('strikeManagement.defaultProclaimedBy') }}</span>
                </div>
              </q-card-section>

              <q-separator />

              <q-card-section class="q-py-sm text-caption text-slate-600">
                <div class="row items-center justify-between q-mb-xs">
                  <span><q-icon name="event" size="15px" class="q-mr-xs text-primary" />{{ t('strikeManagement.strikeDate') }}</span>
                  <span class="text-weight-bold text-slate-800">{{ formatDate(item.strike_date) }}</span>
                </div>
                <div class="row items-center justify-between">
                  <span><q-icon name="timer" size="15px" class="q-mr-xs text-deep-orange-8" />{{ t('strikeManagement.declarationDeadline') }}</span>
                  <span class="text-weight-bold" :class="item.is_expired ? 'text-grey-7' : 'text-deep-orange-9'">
                    {{ formatDateTime(item.declaration_deadline) }}
                  </span>
                </div>
              </q-card-section>

              <q-card-actions align="between" class="bg-slate-50 q-px-md q-py-xs border-t border-slate-100">
                <span class="text-caption text-slate-500">
                  {{ t('strikeManagement.createdAt', { date: formatDate(item.created_at) }) }}
                </span>
                <q-btn
                  flat
                  dense
                  no-caps
                  color="primary"
                  :label="t('strikeManagement.viewDashboard')"
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
        <div class="text-slate-500 q-mt-md">{{ t('strikeManagement.subtitle') }}...</div>
      </div>

      <div v-else-if="summary" class="summary-container q-mt-xl">
        <!-- Selected Notice Banner -->
        <q-card class="rounded-xl shadow-sm border border-slate-200 q-mb-lg bg-white">
          <q-card-section class="q-pa-md q-pa-lg-lg">
            <div class="row items-start justify-between q-col-gutter-md">
              <div class="col-12 col-lg-8">
                <div class="row items-center q-gutter-sm q-mb-xs">
                  <q-badge color="primary" text-color="white" class="text-weight-bold q-px-sm">
                    {{ t('strikeManagement.activeBoardBadge') }}
                  </q-badge>
                  <q-badge
                    :color="summary.notice.is_expired ? 'grey-7' : 'positive'"
                    text-color="white"
                    class="text-weight-bold q-px-sm"
                  >
                    {{ summary.notice.is_expired ? t('strikeManagement.statusClosed') : t('strikeManagement.statusInProgress') }}
                  </q-badge>
                </div>
                <h2 class="text-h5 text-weight-bolder text-slate-900 q-my-none">
                  {{ summary.notice.title }}
                </h2>
                <div class="text-body2 text-slate-600 q-mt-xs">
                  {{ t('strikeManagement.proclaimedBy') }} <b>{{ summary.notice.proclaimed_by }}</b>
                  <span class="q-mx-sm">•</span>
                  {{ t('strikeManagement.dateOfExecution') }} <b>{{ formatDate(summary.notice.strike_date) }}</b>
                  <span class="q-mx-sm">•</span>
                  {{ t('strikeManagement.staffDeadline') }} <b>{{ formatDateTime(summary.notice.declaration_deadline) }}</b>
                </div>
                <div v-if="summary.notice.notes || summary.notice.content" class="bg-slate-50 border border-slate-200 rounded-borders q-pa-sm q-mt-sm text-caption text-slate-700">
                  <q-icon name="sticky_note_2" color="primary" class="q-mr-xs" />
                  <b>{{ t('strikeManagement.serviceNotes') }}</b> {{ summary.notice.notes || summary.notice.content }}
                </div>
              </div>

              <div class="col-12 col-lg-4 text-right">
                <q-btn
                  color="positive"
                  icon="download"
                  :label="t('strikeManagement.exportCsvBtn')"
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
                  <span class="text-caption text-weight-bold text-slate-500 uppercase">{{ t('strikeManagement.statTotalPersonnel') }}</span>
                  <q-avatar size="32px" color="indigo-50" text-color="indigo-7" icon="people" />
                </div>
                <div class="text-h4 text-weight-bolder text-slate-800 q-my-xs">
                  {{ summary.total_staff }}
                </div>
                <div class="text-caption text-slate-600">
                  {{ t('strikeManagement.answeredStat', { answered: summary.answered_count, rate: summary.response_rate }) }}
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
                  <span class="text-caption text-weight-bold text-positive uppercase">{{ t('strikeManagement.statParticipates') }}</span>
                  <q-avatar size="32px" color="positive" text-color="white" icon="check_circle" />
                </div>
                <div class="text-h4 text-weight-bolder text-positive q-my-xs">
                  {{ summary.participates_count }}
                </div>
                <div class="text-caption text-positive">
                  <b>{{ summary.participates_rate }}%</b> {{ t('strikeManagement.overallPersonnel') }}
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
                  <span class="text-caption text-weight-bold text-negative uppercase">{{ t('strikeManagement.statNotParticipates') }}</span>
                  <q-avatar size="32px" color="negative" text-color="white" icon="cancel" />
                </div>
                <div class="text-h4 text-weight-bolder text-negative q-my-xs">
                  {{ summary.not_participates_count }}
                </div>
                <div class="text-caption text-negative">
                  <b>{{ summary.not_participates_rate }}%</b> {{ t('strikeManagement.overallPersonnel') }}
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
                  <span class="text-caption text-weight-bold text-amber-9 uppercase">{{ t('strikeManagement.statUndecided') }}</span>
                  <q-avatar size="32px" color="amber-8" text-color="white" icon="help_outline" />
                </div>
                <div class="text-h4 text-weight-bolder text-amber-9 q-my-xs">
                  {{ summary.undecided_count }}
                </div>
                <div class="text-caption text-amber-9">
                  <b>{{ summary.undecided_rate }}%</b> • {{ t('strikeManagement.waitingStat', { count: summary.unanswered_count }) }}
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
          <div class="row items-center justify-between q-mb-md gap-sm">
            <div>
              <h2 class="text-h6 text-weight-bold text-slate-800 q-my-none flex items-center">
                <q-icon name="pie_chart" color="primary" class="q-mr-sm" size="22px" />
                {{ t('strikeManagement.roleBreakdownTitle') }}
              </h2>
              <div class="text-caption text-slate-500">
                {{ t('strikeManagement.roleBreakdownSubtitle') }}
              </div>
            </div>

            <!-- View Mode Controls -->
            <div class="row items-center q-gutter-xs">
              <q-btn-toggle
                v-model="viewMode"
                dense
                rounded
                no-caps
                unelevated
                toggle-color="primary"
                toggle-text-color="white"
                color="grey-2"
                text-color="slate-700"
                class="shadow-xs"
                :options="[
                  { label: 'Panoramica & Dettaglio', value: 'both', icon: 'dashboard' },
                  { label: 'Solo Macro-Aree', value: 'macro', icon: 'view_agenda' },
                  { label: 'Solo Profili', value: 'roles', icon: 'grid_view' }
                ]"
              />
            </div>
          </div>

          <!-- MACRO CATEGORIES (3 Cards: Docenti, Amministrativi & Tecnici, Collaboratori Scolastici) -->
          <div v-if="(viewMode === 'both' || viewMode === 'macro') && macroCategories.length > 0" class="row q-col-gutter-md q-mb-lg">
            <div
              v-for="macro in macroCategories"
              :key="macro.id"
              class="col-12 col-md-4"
            >
              <q-card class="rounded-xl border border-slate-200 shadow-sm full-height bg-white macro-card">
                <q-card-section class="q-pb-sm">
                  <!-- Macro Header -->
                  <div class="row items-center justify-between no-wrap q-mb-sm">
                    <div class="row items-center no-wrap">
                      <q-avatar size="38px" :class="macro.iconBg" class="q-mr-sm shadow-xs">
                        <q-icon :name="macro.icon" size="22px" />
                      </q-avatar>
                      <div>
                        <div class="text-subtitle2 text-weight-bolder text-slate-800 line-height-tight">
                          {{ macro.name }}
                        </div>
                        <div class="text-caption text-slate-400 text-xs">
                          {{ macro.subtitle }}
                        </div>
                      </div>
                    </div>
                    <q-badge
                      color="slate-100"
                      text-color="slate-700"
                      class="q-px-sm q-py-xs text-caption text-weight-bold rounded-borders font-mono"
                    >
                      Totale: {{ macro.total }}
                    </q-badge>
                  </div>

                  <!-- Rate Highlights -->
                  <div class="row items-baseline justify-between q-mt-md q-mb-xs">
                    <div class="text-caption text-slate-500 font-medium">
                      Adesione Rilevata:
                    </div>
                    <div class="text-h6 text-weight-bolder text-positive">
                      {{ macro.participatesRate }}%
                      <span class="text-caption text-weight-medium text-slate-400">({{ macro.participates }}/{{ macro.total }})</span>
                    </div>
                  </div>

                  <!-- Stacked Breakdown Progress Bar -->
                  <div class="stacked-progress-bar rounded-borders overflow-hidden bg-slate-100 row no-wrap q-my-sm" style="height: 8px;">
                    <div v-if="macro.participates > 0" class="bg-positive" :style="{ width: ((macro.participates / macro.total) * 100) + '%' }" :title="`Aderiscono: ${macro.participates}`" />
                    <div v-if="macro.not_participates > 0" class="bg-negative" :style="{ width: ((macro.not_participates / macro.total) * 100) + '%' }" :title="`Non aderiscono: ${macro.not_participates}`" />
                    <div v-if="macro.undecided > 0" class="bg-amber-8" :style="{ width: ((macro.undecided / macro.total) * 100) + '%' }" :title="`Indecisi: ${macro.undecided}`" />
                    <div v-if="macro.unanswered > 0" class="bg-slate-300" :style="{ width: ((macro.unanswered / macro.total) * 100) + '%' }" :title="`In attesa: ${macro.unanswered}`" />
                  </div>

                  <!-- 4-Stat Grid -->
                  <div class="row q-col-gutter-xs q-mt-xs">
                    <div class="col-6">
                      <div class="p-2 rounded-lg bg-emerald-50/70 border border-emerald-100 text-caption">
                        <div class="text-emerald-700 text-xs font-semibold flex items-center">
                          <q-icon name="check_circle" size="13px" class="q-mr-xs text-positive" />
                          Aderiscono
                        </div>
                        <div class="text-weight-bolder text-emerald-900 text-subtitle2">
                          {{ macro.participates }}
                        </div>
                      </div>
                    </div>
                    <div class="col-6">
                      <div class="p-2 rounded-lg bg-rose-50/70 border border-rose-100 text-caption">
                        <div class="text-rose-700 text-xs font-semibold flex items-center">
                          <q-icon name="cancel" size="13px" class="q-mr-xs text-negative" />
                          Non Aderiscono
                        </div>
                        <div class="text-weight-bolder text-rose-900 text-subtitle2">
                          {{ macro.not_participates }}
                        </div>
                      </div>
                    </div>
                    <div class="col-6">
                      <div class="p-2 rounded-lg bg-amber-50/70 border border-amber-100 text-caption">
                        <div class="text-amber-800 text-xs font-semibold flex items-center">
                          <q-icon name="help_outline" size="13px" class="q-mr-xs text-amber-8" />
                          Non Deciso
                        </div>
                        <div class="text-weight-bolder text-amber-900 text-subtitle2">
                          {{ macro.undecided }}
                        </div>
                      </div>
                    </div>
                    <div class="col-6">
                      <div class="p-2 rounded-lg bg-slate-50 border border-slate-200 text-caption">
                        <div class="text-slate-600 text-xs font-semibold flex items-center">
                          <q-icon name="hourglass_empty" size="13px" class="q-mr-xs text-slate-400" />
                          In Attesa
                        </div>
                        <div class="text-weight-bolder text-slate-700 text-subtitle2">
                          {{ macro.unanswered }}
                        </div>
                      </div>
                    </div>
                  </div>
                </q-card-section>
              </q-card>
            </div>
          </div>

          <!-- DETAILED ROLES SECTION (Filter Pills + Responsive Grid) -->
          <div v-if="viewMode === 'both' || viewMode === 'roles'">
            <div class="row items-center justify-between q-mb-sm gap-xs">
              <div class="text-subtitle2 text-weight-bold text-slate-700 flex items-center">
                <q-icon name="list_alt" size="18px" class="q-mr-xs text-primary" />
                Dettaglio per Singolo Profilo Professionale ({{ filteredByRoleCards.length }})
              </div>

              <!-- Filter pills for role detail -->
              <div class="row items-center q-gutter-xs">
                <q-btn
                  dense
                  rounded
                  no-caps
                  size="sm"
                  :color="macroRoleFilter === 'all' ? 'primary' : 'grey-2'"
                  :text-color="macroRoleFilter === 'all' ? 'white' : 'slate-700'"
                  label="Tutti i Profili"
                  class="q-px-sm text-weight-medium"
                  @click="macroRoleFilter = 'all'"
                />
                <q-btn
                  dense
                  rounded
                  no-caps
                  size="sm"
                  :color="macroRoleFilter === 'teachers' ? 'indigo-7' : 'grey-2'"
                  :text-color="macroRoleFilter === 'teachers' ? 'white' : 'slate-700'"
                  label="Docenti"
                  class="q-px-sm text-weight-medium"
                  @click="macroRoleFilter = 'teachers'"
                />
                <q-btn
                  dense
                  rounded
                  no-caps
                  size="sm"
                  :color="macroRoleFilter === 'admin_tech' ? 'teal-7' : 'grey-2'"
                  :text-color="macroRoleFilter === 'admin_tech' ? 'white' : 'slate-700'"
                  label="Amministrativi & Tecnici"
                  class="q-px-sm text-weight-medium"
                  @click="macroRoleFilter = 'admin_tech'"
                />
                <q-btn
                  dense
                  rounded
                  no-caps
                  size="sm"
                  :color="macroRoleFilter === 'collaboratori' ? 'amber-9' : 'grey-2'"
                  :text-color="macroRoleFilter === 'collaboratori' ? 'white' : 'slate-700'"
                  label="Collaboratori"
                  class="q-px-sm text-weight-medium"
                  @click="macroRoleFilter = 'collaboratori'"
                />
              </div>
            </div>

            <!-- Responsive Cards Grid: col-12 col-sm-6 col-md-4 col-lg-3 col-xl-2 -->
            <div class="row q-col-gutter-md">
              <div
                v-for="cat in filteredByRoleCards"
                :key="cat.role"
                class="col-12 col-sm-6 col-md-4 col-lg-3 col-xl-2"
              >
                <q-card
                  class="rounded-xl border border-slate-200 shadow-sm full-height bg-white role-card cursor-pointer"
                  :class="{ 'border-primary ring-2 ring-primary ring-offset-1 bg-indigo-50/20': filterRole === cat.role }"
                  @click="selectRoleFilter(cat.role)"
                >
                  <q-card-section class="q-pa-sm flex column justify-between full-height">
                    <div>
                      <!-- Card Header: Role Chip + Total -->
                      <div class="row items-center justify-between no-wrap q-mb-xs">
                        <q-chip
                          dense
                          :color="getRoleBadgeColor(cat.role)"
                          text-color="white"
                          class="text-weight-bolder text-caption q-ma-none text-truncate"
                          style="max-width: 170px;"
                          :title="t('roles.' + cat.role) || cat.role_display"
                        >
                          {{ t('roles.' + cat.role) || cat.role_display }}
                        </q-chip>
                        <q-badge
                          color="slate-100"
                          text-color="slate-700"
                          class="text-caption text-weight-bold rounded-borders font-mono"
                        >
                          {{ cat.total }}
                        </q-badge>
                      </div>

                      <!-- Stacked Progress Bar -->
                      <div class="stacked-progress-bar rounded-borders overflow-hidden bg-slate-100 row no-wrap q-my-xs" style="height: 6px;">
                        <div v-if="cat.participates > 0" class="bg-positive" :style="{ width: ((cat.participates / cat.total) * 100) + '%' }" :title="`Aderiscono: ${cat.participates}`" />
                        <div v-if="cat.not_participates > 0" class="bg-negative" :style="{ width: ((cat.not_participates / cat.total) * 100) + '%' }" :title="`Non aderiscono: ${cat.not_participates}`" />
                        <div v-if="cat.undecided > 0" class="bg-amber-8" :style="{ width: ((cat.undecided / cat.total) * 100) + '%' }" :title="`Indecisi: ${cat.undecided}`" />
                        <div v-if="cat.unanswered > 0" class="bg-slate-300" :style="{ width: ((cat.unanswered / cat.total) * 100) + '%' }" :title="`In attesa: ${cat.unanswered}`" />
                      </div>
                    </div>

                    <!-- Metrics -->
                    <div class="q-mt-xs">
                      <div class="row items-center justify-between text-caption q-mb-xs">
                        <span class="text-positive text-weight-bold flex items-center">
                          <q-icon name="check_circle" size="12px" class="q-mr-xs" />
                          {{ cat.participates }}
                        </span>
                        <span class="text-negative text-weight-bold flex items-center">
                          <q-icon name="cancel" size="12px" class="q-mr-xs" />
                          {{ cat.not_participates }}
                        </span>
                      </div>

                      <div class="row items-center justify-between text-caption text-slate-500">
                        <span class="text-amber-9 text-weight-medium flex items-center text-xs">
                          <q-icon name="help_outline" size="12px" class="q-mr-xs" />
                          {{ cat.undecided }}
                        </span>
                        <span class="flex items-center text-slate-400 text-xs">
                          <q-icon name="hourglass_empty" size="12px" class="q-mr-xs" />
                          {{ cat.unanswered }}
                        </span>
                      </div>
                    </div>
                  </q-card-section>
                </q-card>
              </div>
            </div>
          </div>
        </div>

        <!-- Detailed Nominative Table Section -->
        <q-card id="nominative-table-card" class="rounded-xl border border-slate-200 shadow-sm bg-white overflow-hidden">
          <q-card-section class="q-pb-none">
            <div class="row items-center justify-between q-mb-md gap-sm">
              <div>
                <h3 class="text-h6 text-weight-bold text-slate-800 q-my-none flex items-center">
                  <q-icon name="format_list_bulleted" color="primary" class="q-mr-sm" size="22px" />
                  {{ t('strikeManagement.nominativeTableTitle') }}
                </h3>
                <div class="text-caption text-slate-500">
                  {{ t('strikeManagement.nominativeTableSubtitle') }}
                </div>
              </div>

              <div class="row items-center q-gutter-sm">
                <q-input
                  v-model="filterSearch"
                  dense
                  outlined
                  rounded
                  :placeholder="t('strikeManagement.searchStaffPlaceholder')"
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
            :no-data-label="t('strikeManagement.noStaffFound')"
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
                  {{ t('roles.' + props.row.role) || props.row.role_display }}
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
                  {{ t('strikeManagement.notReceived') }}
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
            {{ t('strikeManagement.dialogNewTitle') }}
          </div>
          <q-btn flat round dense icon="close" color="white" v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md q-gutter-y-md">
          <div class="text-caption text-slate-600">
            {{ t('strikeManagement.noNoticesDesc') }}
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('strikeManagement.noticeTitleLabel') }}</label>
            <q-input
              v-model="createForm.title"
              outlined
              dense
              placeholder="es. Sciopero generale del comparto Istruzione e Ricerca"
              :rules="[val => !!val || t('strikeManagement.notifyFillRequired')]"
            />
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('strikeManagement.proclaimedByLabel') }}</label>
            <q-input
              v-model="createForm.proclaimed_by"
              outlined
              dense
              placeholder="es. FLC CGIL, CISL FSUR, UIL SCUOLA RUA, SNALS CONFSAL"
              :rules="[val => !!val || t('strikeManagement.notifyFillRequired')]"
            />
          </div>

          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('strikeManagement.strikeDateLabel') }}</label>
              <q-input
                v-model="createForm.strike_date"
                outlined
                dense
                type="date"
                :rules="[val => !!val || t('strikeManagement.notifyFillRequired')]"
              />
            </div>

            <div class="col-12 col-sm-6">
              <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('strikeManagement.deadlineLabel') }}</label>
              <q-input
                v-model="createForm.declaration_deadline"
                outlined
                dense
                type="datetime-local"
                :rules="[
                  val => !!val || t('strikeManagement.notifyFillRequired'),
                  val => isDeadlineValid(val, createForm.strike_date) || t('strikeManagement.notifyFillRequired')
                ]"
              />
            </div>
          </div>

          <div>
            <label class="text-weight-bold text-caption text-slate-700 block q-mb-xs">{{ t('strikeManagement.notesLabel') }}</label>
            <q-input
              v-model="createForm.notes"
              type="textarea"
              rows="3"
              outlined
              dense
            />
          </div>

          <div class="bg-indigo-50 border border-indigo-100 rounded-borders q-pa-sm">
            <q-checkbox
              v-model="createForm.publish_to_bacheca"
              :label="t('strikeManagement.publishToBachecaLabel')"
              color="primary"
              class="text-weight-bold text-slate-800"
            />
            <div class="text-caption text-slate-600 q-ml-lg">
              {{ t('strikeManagement.publishToBachecaHint') }}
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md border-t border-slate-100">
          <q-btn flat :label="t('strikeManagement.cancelBtn')" color="grey-7" v-close-popup />
          <q-btn
            color="primary"
            :label="t('strikeManagement.publishNoticeBtn')"
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
          <span class="q-ml-sm text-h6 text-weight-bold">{{ t('strikeManagement.confirmDeleteNoticeTitle') }}</span>
        </q-card-section>

        <q-card-section class="q-pt-sm text-body2 text-slate-600">
          {{ t('strikeManagement.confirmDeleteNoticeMsg', { title: noticeToDelete?.title || '' }) }}
          <br />
          {{ t('strikeManagement.deleteWarningMsg') }}
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat :label="t('strikeManagement.cancelBtn')" color="grey-7" v-close-popup />
          <q-btn
            color="negative"
            :label="t('strikeManagement.deletePermanentlyBtn')"
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
import { useI18n } from 'vue-i18n'
import strikeService from '@/services/strikeService'

const $q = useQuasar()
const { t, locale } = useI18n()

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
const viewMode = ref('both')
const macroRoleFilter = ref('all')

const macroCategories = computed(() => {
  if (!summary.value || !summary.value.by_role) return []

  const groups = [
    {
      id: 'teachers',
      name: 'Personale Docente',
      subtitle: 'Docenti curricolari, sostegno e coordinatori',
      icon: 'school',
      iconBg: 'bg-indigo-50 text-indigo-7',
      chipColor: 'indigo-8',
      roles: ['teacher', 'coordinator', 'docente'],
      total: 0,
      participates: 0,
      not_participates: 0,
      undecided: 0,
      unanswered: 0
    },
    {
      id: 'admin_tech',
      name: 'Personale Amministrativo & Tecnico',
      subtitle: 'DSGA, Assistenti Amministrativi, Tecnici e Segreteria',
      icon: 'admin_panel_settings',
      iconBg: 'bg-teal-50 text-teal-7',
      chipColor: 'teal-8',
      roles: [
        'dsga',
        'assistente_amministrativo',
        'assistente_tecnico',
        'assistente_contabilita',
        'assistente_protocollo',
        'assistente_sportello',
        'secretary'
      ],
      total: 0,
      participates: 0,
      not_participates: 0,
      undecided: 0,
      unanswered: 0
    },
    {
      id: 'collaboratori',
      name: 'Collaboratori Scolastici & Ausiliari',
      subtitle: 'Collaboratori scolastici, mensa e assistenza',
      icon: 'support',
      iconBg: 'bg-amber-50 text-amber-8',
      chipColor: 'amber-9',
      roles: [
        'collaboratore_scolastico',
        'collaboratore_ds',
        'collaboratore_mensa',
        'assistente_alunni',
        'assistente_personale',
        'responsabile_servizio'
      ],
      total: 0,
      participates: 0,
      not_participates: 0,
      undecided: 0,
      unanswered: 0
    }
  ]

  for (const r of summary.value.by_role) {
    let targetGroup = groups.find(g => g.roles.includes(r.role))
    if (!targetGroup) {
      targetGroup = groups[1] // fallback to admin_tech
    }
    targetGroup.total += (r.total || 0)
    targetGroup.participates += (r.participates || 0)
    targetGroup.not_participates += (r.not_participates || 0)
    targetGroup.undecided += (r.undecided || 0)
    targetGroup.unanswered += (r.unanswered || 0)
  }

  return groups.map(g => ({
    ...g,
    participatesRate: g.total > 0 ? Math.round((g.participates / g.total) * 1000) / 10 : 0,
    notParticipatesRate: g.total > 0 ? Math.round((g.not_participates / g.total) * 1000) / 10 : 0,
    undecidedRate: g.total > 0 ? Math.round((g.undecided / g.total) * 1000) / 10 : 0,
    unansweredRate: g.total > 0 ? Math.round((g.unanswered / g.total) * 1000) / 10 : 0
  }))
})

const filteredByRoleCards = computed(() => {
  if (!summary.value || !summary.value.by_role) return []
  let list = summary.value.by_role

  if (macroRoleFilter.value !== 'all') {
    const macro = macroCategories.value.find(m => m.id === macroRoleFilter.value)
    if (macro) {
      list = list.filter(r => macro.roles.includes(r.role))
    }
  }

  return list
})

const selectRoleFilter = (role) => {
  if (filterRole.value === role) {
    filterRole.value = 'all'
  } else {
    filterRole.value = role
    const tableEl = document.getElementById('nominative-table-card')
    if (tableEl) {
      tableEl.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }
}

const createForm = ref({
  title: '',
  proclaimed_by: '',
  strike_date: '',
  declaration_deadline: '',
  notes: '',
  publish_to_bacheca: true
})

const columns = computed(() => [
  { name: 'person', label: t('strikeManagement.colStaff'), align: 'left', field: 'last_name', sortable: true },
  { name: 'role', label: t('strikeManagement.colCategory'), align: 'left', field: 'role', sortable: true },
  { name: 'intention', label: t('strikeManagement.colIntention'), align: 'center', field: 'intention', sortable: true },
  { name: 'declared_at', label: t('strikeManagement.colChoiceDate'), align: 'left', field: 'declared_at', sortable: true }
])

const roleFilterOptions = computed(() => {
  const options = [{ label: t('strikeManagement.allCategories'), value: 'all' }]
  if (summary.value?.by_role && summary.value.by_role.length > 0) {
    summary.value.by_role.forEach(r => {
      options.push({
        label: `${t('roles.' + r.role) || r.role_display} (${r.total})`,
        value: r.role
      })
    })
  } else {
    options.push(
      { label: t('roles.teacher'), value: 'teacher' },
      { label: t('roles.coordinator'), value: 'coordinator' },
      { label: t('roles.dsga'), value: 'dsga' },
      { label: t('roles.assistente_amministrativo'), value: 'assistente_amministrativo' },
      { label: t('roles.assistente_tecnico'), value: 'assistente_tecnico' },
      { label: t('roles.collaboratore_scolastico'), value: 'collaboratore_scolastico' },
      { label: t('roles.secretary'), value: 'secretary' }
    )
  }
  return options
})

const intentionFilterOptions = computed(() => [
  { label: t('strikeManagement.allIntentions'), value: 'all' },
  { label: `${t('strikeManagement.statParticipates')} (${t('strikeManagement.participates')})`, value: 'participates' },
  { label: t('strikeManagement.statNotParticipates'), value: 'not_participates' },
  { label: t('strikeManagement.statUndecided'), value: 'undecided' },
  { label: t('strikeManagement.noResponse'), value: 'unanswered' }
])

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
      message: t('strikeManagement.notifyLoadError')
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
      message: t('strikeManagement.notifyLoadError')
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
      message: t('strikeManagement.notifyFillRequired')
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
      content: createForm.value.notes || '',
      notes: createForm.value.notes || '',
      publish_to_bacheca: createForm.value.publish_to_bacheca
    }

    const created = await strikeService.createStrikeNotice(payload)
    showCreateDialog.value = false
    $q.notify({
      type: 'positive',
      message: t('strikeManagement.notifyNoticePublished'),
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
      message: err.response?.data?.error || t('strikeManagement.notifyLoadError')
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
      message: t('strikeManagement.notifyNoticeDeleted')
    })
    noticeToDelete.value = null
    await loadNotices()
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: t('strikeManagement.notifyLoadError')
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
      `"${s.declared_at ? formatDateTime(s.declared_at) : t('strikeManagement.notReceived')}"`
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
    case 'participates': return t('strikeManagement.participates')
    case 'not_participates': return t('strikeManagement.notParticipates')
    case 'undecided': return t('strikeManagement.undecided')
    default: return t('strikeManagement.noResponse')
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
    case 'docente':
      return 'primary'
    case 'coordinator':
      return 'indigo-7'
    case 'dsga':
      return 'deep-purple-7'
    case 'assistente_amministrativo':
      return 'teal-7'
    case 'assistente_tecnico':
      return 'cyan-8'
    case 'collaboratore_ds':
      return 'blue-8'
    case 'collaboratore_scolastico':
      return 'amber-9'
    case 'collaboratore_mensa':
      return 'orange-8'
    case 'assistente_alunni':
      return 'green-8'
    case 'assistente_personale':
      return 'emerald-7'
    case 'assistente_contabilita':
      return 'teal-9'
    case 'assistente_protocollo':
      return 'purple-7'
    case 'assistente_sportello':
      return 'light-blue-8'
    case 'responsabile_servizio':
      return 'deep-orange-7'
    case 'secretary':
      return 'blue-grey-8'
    case 'principal':
      return 'purple-9'
    default:
      return 'indigo-6'
  }
}

const formatDate = (isoDate) => {
  if (!isoDate) return '-'
  try {
    const d = new Date(isoDate)
    return d.toLocaleDateString(locale.value || 'it-IT', { day: '2-digit', month: 'long', year: 'numeric' })
  } catch {
    return isoDate
  }
}

const formatDateTime = (isoDate) => {
  if (!isoDate) return '-'
  try {
    const d = new Date(isoDate)
    return d.toLocaleDateString(locale.value || 'it-IT', {
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
.macro-card {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.macro-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.08);
}
.role-card {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.role-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 18px rgba(0, 0, 0, 0.08);
}
.line-height-tight {
  line-height: 1.25;
}
.text-xs {
  font-size: 0.75rem;
}
.search-input {
  width: 200px;
}
.role-select {
  width: 220px;
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
