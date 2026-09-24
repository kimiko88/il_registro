<template>
  <q-page padding class="min-h-screen" :class="$q.dark.isActive ? 'bg-dark text-white' : 'bg-slate-50 text-slate-800'">
    <!-- Top Header -->
    <div class="row items-center justify-between q-mb-lg gap-4">
      <div>
        <h1 class="text-h4 text-weight-bold q-my-none row items-center gap-2" :class="$q.dark.isActive ? 'text-white' : 'text-slate-800'">
          <q-icon name="thumb_up_alt" color="primary" size="36px" />
          {{ t('schedulePreferences.title') }}
        </h1>
        <p class="text-subtitle1 q-mt-xs q-mb-none" :class="$q.dark.isActive ? 'text-grey-4' : 'text-slate-500'">
          {{ t('schedulePreferences.subtitle') }}
        </p>
      </div>

      <div class="row items-center gap-2">
        <q-btn
          :label="t('schedulePreferences.saveBtn')"
          icon="save"
          color="primary"
          unelevated
          no-caps
          class="rounded-xl q-px-lg font-bold shadow-sm"
          :loading="saving"
          @click="savePreferences"
        />
      </div>
    </div>

    <!-- Seniority & Priority Alert Banner -->
    <q-banner rounded class="bg-indigo-50 border border-indigo-200 text-indigo-950 q-mb-lg rounded-2xl shadow-xs">
      <template v-slot:avatar>
        <q-icon name="stars" color="indigo-7" size="32px" />
      </template>
      <div class="text-subtitle2 font-bold">{{ t('schedulePreferences.seniorityTitle') }}</div>
      <div class="text-body2 text-indigo-900">
        {{ t('schedulePreferences.seniorityDesc') }}
      </div>
    </q-banner>

    <!-- Legend and Quick Controls -->
    <q-card flat bordered class="rounded-2xl q-pa-md q-mb-lg shadow-sm" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
      <div class="row items-center justify-between gap-4">
        <div class="row items-center gap-4 text-sm font-medium">
          <span class="text-grey-6">{{ t('schedulePreferences.legend') }}</span>
          <div class="row items-center gap-1">
            <span class="inline-block w-4 h-4 rounded-full bg-emerald-500"></span>
            <span>🟢 {{ t('schedulePreferences.preferred') }}</span>
          </div>
          <div class="row items-center gap-1">
            <span class="inline-block w-4 h-4 rounded-full bg-slate-300"></span>
            <span>⚪ {{ t('schedulePreferences.neutral') }}</span>
          </div>
          <div class="row items-center gap-1">
            <span class="inline-block w-4 h-4 rounded-full bg-rose-500"></span>
            <span>🔴 {{ t('schedulePreferences.unavailable') }}</span>
          </div>
        </div>

        <div class="text-caption text-grey-6">
          {{ t('schedulePreferences.clickToChange') }}
        </div>
      </div>
    </q-card>

    <!-- Interactive Preferences Grid -->
    <q-card flat bordered class="rounded-2xl p-4 shadow-sm" :class="$q.dark.isActive ? 'bg-dark border-grey-8' : 'bg-white'">
      <div v-if="loading" class="row justify-center q-pa-xl">
        <q-spinner-dots size="48px" color="primary" />
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full border-collapse">
          <thead>
            <tr class="bg-slate-100 text-slate-700">
              <th class="p-3 border border-slate-200 text-left w-24">{{ t('schedulePreferences.hourCol') }}</th>
              <th
                v-for="d in days"
                :key="d.index"
                class="p-3 border border-slate-200 text-center font-bold text-base"
              >
                {{ d.label }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="h in 8" :key="h">
              <td class="p-3 border border-slate-200 font-bold text-center bg-slate-50 text-slate-800">
                {{ t('schedulePreferences.hourSlot', { hour: h }) }}
              </td>
              <td
                v-for="d in days"
                :key="d.index + '-' + h"
                class="p-2 border border-slate-200 text-center cursor-pointer transition-colors select-none"
                :class="getCellBgClass(d.index, h)"
                @click="cyclePreference(d.index, h)"
              >
                <div class="py-2 px-1 rounded-lg">
                  <div class="font-bold text-xs" :class="getCellTextClass(d.index, h)">
                    {{ getCellLabel(d.index, h) }}
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';
import timetableGenService from '@/services/timetableGenService';

const $q = useQuasar();
const { t } = useI18n();

const loading = ref(false);
const saving = ref(false);

const days = computed(() => [
  { index: 1, label: t('schedulePreferences.days.monday') },
  { index: 2, label: t('schedulePreferences.days.tuesday') },
  { index: 3, label: t('schedulePreferences.days.wednesday') },
  { index: 4, label: t('schedulePreferences.days.thursday') },
  { index: 5, label: t('schedulePreferences.days.friday') }
]);

// Grid state: grid[dayIndex][hourIndex] = 'preferred' | 'neutral' | 'unavailable'
const grid = ref({});

function initGrid() {
  const g = {};
  for (let d = 1; d <= 5; d++) {
    g[d] = {};
    for (let h = 1; h <= 8; h++) {
      g[d][h] = 'neutral';
    }
  }
  grid.value = g;
}

async function fetchPreferences() {
  loading.value = true;
  try {
    const res = await timetableGenService.getPreferences();
    const prefs = res.data || [];
    initGrid();
    for (const p of prefs) {
      if (grid.value[p.day_of_week] && grid.value[p.day_of_week][p.hour_index] !== undefined) {
        grid.value[p.day_of_week][p.hour_index] = p.preference_type;
      }
    }
  } catch {
    $q.notify({ type: 'negative', message: t('schedulePreferences.notifyFetchError') });
  } finally {
    loading.value = false;
  }
}

function cyclePreference(day, hour) {
  const current = grid.value[day]?.[hour] || 'neutral';
  let next = 'neutral';
  if (current === 'neutral') next = 'preferred';
  else if (current === 'preferred') next = 'unavailable';
  else if (current === 'unavailable') next = 'neutral';

  if (!grid.value[day]) grid.value[day] = {};
  grid.value[day][hour] = next;
}

function getCellBgClass(day, hour) {
  const state = grid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return 'bg-emerald-100 hover:bg-emerald-200';
  if (state === 'unavailable') return 'bg-rose-100 hover:bg-rose-200';
  return 'bg-white hover:bg-slate-100';
}

function getCellTextClass(day, hour) {
  const state = grid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return 'text-emerald-800';
  if (state === 'unavailable') return 'text-rose-800';
  return 'text-slate-400';
}

function getCellLabel(day, hour) {
  const state = grid.value[day]?.[hour] || 'neutral';
  if (state === 'preferred') return t('schedulePreferences.cellPreferred');
  if (state === 'unavailable') return t('schedulePreferences.cellUnavailable');
  return t('schedulePreferences.cellNeutral');
}

async function savePreferences() {
  saving.value = true;
  try {
    const entries = [];
    for (let d = 1; d <= 5; d++) {
      for (let h = 1; h <= 8; h++) {
        const val = grid.value[d]?.[h] || 'neutral';
        if (val !== 'neutral') {
          entries.push({
            day_of_week: d,
            hour_index: h,
            preference_type: val
          });
        }
      }
    }

    await timetableGenService.savePreferences({
      preferences: entries
    });

    $q.notify({
      type: 'positive',
      message: t('schedulePreferences.notifySaveSuccess')
    });
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.response?.data?.error || t('schedulePreferences.notifySaveError')
    });
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  initGrid();
  fetchPreferences();
});
</script>
