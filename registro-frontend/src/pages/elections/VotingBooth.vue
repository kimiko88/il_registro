<template>
  <q-page class="q-pa-md bg-grey-1 flex flex-center">
    <div style="width: 100%; max-width: 600px;">
      <!-- Receipt Card if already voted -->
      <q-card v-if="receipt" flat bordered class="bg-white text-center q-pa-lg">
        <q-icon name="check_circle" color="positive" size="64px" class="q-mb-md" />
        <div class="text-h5 text-weight-bold text-positive">{{ t('schoolElections.voteReceipt') }}</div>
        <p class="text-caption text-grey-8 q-mt-sm">
          {{ t('schoolElections.receiptHelp') }}
        </p>

        <q-banner rounded class="bg-blue-grey-1 text-blue-grey-10 text-left q-my-md">
          <div class="text-caption text-weight-bold">{{ t('schoolElections.voteReceipt') }}:</div>
          <div class="text-subtitle2 text-weight-bold text-primary text-mono q-mt-xs">{{ receipt.receipt_token }}</div>
          <div class="text-caption text-grey-7 q-mt-xs">{{ t('common.details') }}: {{ new Date(receipt.voted_at).toLocaleString() }}</div>
        </q-banner>

        <q-btn color="primary" :label="t('common.backToHome')" unelevated no-caps to="/" />
      </q-card>

      <!-- Voting Booth Stepper -->
      <q-card v-else flat bordered class="bg-white">
        <q-card-section class="bg-primary text-white text-center q-py-md">
          <div class="text-h6 text-weight-bold">
            <q-icon name="how_to_vote" class="q-mr-xs" /> {{ t('schoolElections.boothTitle') }}
          </div>
          <div class="text-caption">{{ election?.title || t('schoolElections.boothSubtitle') }}</div>
        </q-card-section>

        <q-card-section v-if="loading" class="text-center q-pa-xl">
          <q-spinner-dots size="40px" color="primary" />
        </q-card-section>

        <q-card-section v-else class="q-pa-md">
          <div class="text-subtitle2 text-weight-bold text-blue-grey-9 q-mb-sm">
            {{ t('schoolElections.selectList') }}
          </div>

          <q-list bordered separator class="rounded-borders q-mb-md">
            <q-item tag="label" :active="isBlankVote" active-class="bg-grey-2">
              <q-item-section avatar>
                <q-radio v-model="selectedListId" val="blank" @update:model-value="onListChange" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('common.other') }} (Scheda Bianca / Blank)</q-item-label>
                <q-item-label caption>{{ t('schoolElections.boothSubtitle') }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item v-for="list in lists" :key="list.id" tag="label" :active="selectedListId === list.id" active-class="bg-blue-1">
              <q-item-section avatar>
                <q-radio v-model="selectedListId" :val="list.id" @update:model-value="onListChange" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('schoolElections.listName') }} {{ list.list_number }}: "{{ list.motto }}"</q-item-label>
                <q-item-label caption>{{ list.candidates?.length || 0 }} {{ t('schoolElections.electedCandidates') }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>

          <!-- Candidate Preferences -->
          <div v-if="selectedList && !isBlankVote" class="q-mb-md">
            <div class="text-subtitle2 text-weight-bold text-blue-grey-9 q-mb-xs">
              {{ t('schoolElections.selectCandidates') }} (Max: {{ election?.max_preferences || 1 }})
            </div>
            <div v-if="!selectedList.candidates || selectedList.candidates.length === 0" class="text-caption text-grey-6">
              {{ t('common.noData') }}
            </div>
            <q-list v-else bordered separator class="rounded-borders">
              <q-item v-for="cand in selectedList.candidates" :key="cand.id" tag="label">
                <q-item-section avatar>
                  <q-checkbox v-model="selectedCandidateIds" :val="cand.id" :disable="isCandDisabled(cand.id)" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ cand.first_name }} {{ cand.last_name }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </div>

          <q-banner rounded class="bg-amber-1 text-amber-9 text-caption q-mb-md">
            <q-icon name="lock" class="q-mr-xs" />
            {{ t('schoolElections.boothSubtitle') }}
          </q-banner>

          <q-btn
            color="primary"
            class="full-width q-py-sm text-weight-bold"
            :label="t('schoolElections.castVote')"
            unelevated
            :loading="submitting"
            :disable="!selectedListId"
            @click="submitVote"
          />
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { electionsService } from '@/services/electionsService'

export default defineComponent({
  name: 'VotingBooth',
  setup() {
    const { t } = useI18n()
    const route = useRoute()
    const election = ref(null)
    const lists = ref([])
    const loading = ref(true)
    const submitting = ref(false)
    const selectedListId = ref('')
    const selectedCandidateIds = ref([])
    const receipt = ref(null)

    const isBlankVote = computed(() => selectedListId.value === 'blank')

    const selectedList = computed(() => {
      if (isBlankVote.value) return null
      return lists.value.find(l => l.id === selectedListId.value)
    })

    const onListChange = () => {
      selectedCandidateIds.value = []
    }

    const isCandDisabled = (candId) => {
      if (selectedCandidateIds.value.includes(candId)) return false
      const max = election.value?.max_preferences || 1
      return selectedCandidateIds.value.length >= max
    }

    const loadElection = async () => {
      loading.value = true
      try {
        const id = route.params.id || route.query.election_id
        if (!id) {
          const allRes = await electionsService.getElections()
          if (allRes.data && allRes.data.length > 0) {
            const first = allRes.data[0]
            const detRes = await electionsService.getElectionDetails(first.id)
            election.value = detRes.election
            lists.value = detRes.lists || []
          }
        } else {
          const detRes = await electionsService.getElectionDetails(id)
          election.value = detRes.election
          lists.value = detRes.lists || []
        }
      } catch (e) {
        console.error('Errore caricamento elezione:', e)
      } finally {
        loading.value = false
      }
    }

    const submitVote = async () => {
      if (!election.value) return
      submitting.value = true
      try {
        const payload = {
          election_id: election.value.id,
          list_id: isBlankVote.value ? null : selectedListId.value,
          candidate_ids: isBlankVote.value ? [] : selectedCandidateIds.value,
          is_blank: isBlankVote.value
        }
        const res = await electionsService.castVote(election.value.id, payload)
        receipt.value = res.receipt
      } catch (e) {
        console.error('Errore espressione voto:', e)
      } finally {
        submitting.value = false
      }
    }

    onMounted(loadElection)

    return {
      t,
      election,
      lists,
      loading,
      submitting,
      selectedListId,
      selectedCandidateIds,
      receipt,
      isBlankVote,
      selectedList,
      onListChange,
      isCandDisabled,
      submitVote
    }
  }
})
</script>
