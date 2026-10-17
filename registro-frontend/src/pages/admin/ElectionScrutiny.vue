<template>
  <q-page class="q-pa-md bg-grey-1">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">
          <q-icon name="analytics" class="q-mr-sm" />{{ t('schoolElections.scrutinyTitle') }}
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          {{ t('schoolElections.scrutinySubtitle') }}
        </p>
      </div>
      <div class="row q-gutter-sm">
        <q-btn flat icon="refresh" :label="t('common.refresh')" @click="loadScrutiny" :loading="loading" />
        <q-btn color="negative" icon="lock" :label="t('common.close')" unelevated no-caps @click="closeElection" />
      </div>
    </div>

    <!-- Turnout and ballot metrics -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <q-card flat bordered class="bg-blue-1 text-primary">
          <q-card-section>
            <div class="text-subtitle2">{{ t('schoolElections.totalVoters') }}</div>
            <div class="text-h4 text-weight-bold">{{ scrutiny?.total_voters || 0 }}</div>
            <div class="text-caption">{{ t('schoolElections.turnout') }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card flat bordered class="bg-teal-1 text-teal-9">
          <q-card-section>
            <div class="text-subtitle2">{{ t('schoolElections.votesCast') }}</div>
            <div class="text-h4 text-weight-bold">{{ (scrutiny?.total_votes_cast || 0) - (scrutiny?.blank_votes || 0) }}</div>
            <div class="text-caption">{{ t('schoolElections.votes') }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card flat bordered class="bg-grey-3 text-grey-9">
          <q-card-section>
            <div class="text-subtitle2">{{ t('common.other') }}</div>
            <div class="text-h4 text-weight-bold">{{ scrutiny?.blank_votes || 0 }}</div>
            <div class="text-caption">{{ t('schoolElections.boothSubtitle') }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- d'Hondt Seats Allocation Table -->
    <q-card flat bordered class="q-mb-md bg-white">
      <q-card-section class="bg-blue-grey-1 text-weight-bold text-subtitle2 text-blue-grey-9">
        <q-icon name="table_chart" class="q-mr-xs" /> {{ t('schoolElections.seatsAllocation') }}
      </q-card-section>
      <q-separator />
      <q-table
        :rows="scrutiny?.lists_results || []"
        :columns="listColumns"
        row-key="list_id"
        flat
        hide-pagination
      >
        <template #body-cell-seats_won="props">
          <q-td :props="props">
            <q-badge color="positive" class="text-weight-bold text-subtitle2 q-pa-xs">
              {{ props.row.seats_won }} {{ t('schoolElections.seatsWon') }}
            </q-badge>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Proclaimed Elected Members -->
    <q-card flat bordered class="bg-white">
      <q-card-section class="bg-teal-1 text-weight-bold text-subtitle2 text-teal-10">
        <q-icon name="verified" class="q-mr-xs" /> {{ t('schoolElections.proclaimElected') }}
      </q-card-section>
      <q-separator />
      <div v-if="!scrutiny?.elected_members || scrutiny.elected_members.length === 0" class="q-pa-md text-caption text-grey-6 text-center">
        {{ t('common.noData') }}
      </div>
      <q-list v-else separator>
        <q-item v-for="cand in scrutiny.elected_members" :key="cand.candidate_id">
          <q-item-section avatar>
            <q-avatar color="positive" text-color="white" icon="star" />
          </q-item-section>
          <q-item-section>
            <q-item-label class="text-weight-bold">{{ cand.first_name }} {{ cand.last_name }}</q-item-label>
            <q-item-label caption>{{ t('schoolElections.listName') }}: "{{ cand.list_motto }}" | {{ t('schoolElections.votes') }}: {{ cand.votes }}</q-item-label>
          </q-item-section>
          <q-item-section side>
            <q-chip color="positive" text-color="white" size="sm" class="text-weight-bold">{{ t('schoolElections.electedCandidates') }}</q-chip>
          </q-item-section>
        </q-item>
      </q-list>
    </q-card>
  </q-page>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { electionsService } from '@/services/electionsService'

export default defineComponent({
  name: 'ElectionScrutiny',
  setup() {
    const { t } = useI18n()
    const route = useRoute()
    const scrutiny = ref(null)
    const loading = ref(false)
    const electionId = ref('')

    const listColumns = computed(() => [
      { name: 'list_number', label: t('schoolElections.listName') + ' N.', field: 'list_number', align: 'center' },
      { name: 'motto', label: t('common.description'), field: 'motto', align: 'left' },
      { name: 'total_votes', label: t('schoolElections.votes'), field: 'total_votes', align: 'center' },
      { name: 'seats_won', label: t('schoolElections.seatsWon'), field: 'seats_won', align: 'center' }
    ])

    const loadScrutiny = async () => {
      loading.value = true
      try {
        electionId.value = route.params.id || route.query.election_id
        if (!electionId.value) {
          const allRes = await electionsService.getElections()
          if (allRes.data && allRes.data.length > 0) {
            electionId.value = allRes.data[0].id
          }
        }
        if (electionId.value) {
          const res = await electionsService.getScrutiny(electionId.value, 4)
          scrutiny.value = res.scrutiny
        }
      } catch (e) {
        console.error('Errore scrutinio elezioni:', e)
      } finally {
        loading.value = false
      }
    }

    const closeElection = async () => {
      if (!electionId.value) return
      try {
        await electionsService.closeElection(electionId.value)
        await loadScrutiny()
      } catch (e) {
        console.error('Errore chiusura seggio:', e)
      }
    }

    onMounted(loadScrutiny)

    return {
      t,
      scrutiny,
      loading,
      listColumns,
      loadScrutiny,
      closeElection
    }
  }
})
</script>
