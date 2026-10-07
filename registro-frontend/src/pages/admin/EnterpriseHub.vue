<template>
  <q-page class="q-pa-lg enterprise-hub-page">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-xl">
      <div>
        <div class="row items-center q-gutter-x-sm">
          <q-badge color="primary" text-color="white" :label="$t('enterprise.header.badgeCompliance')" class="q-py-xs q-px-sm text-weight-bold" />
          <q-badge color="positive" text-color="white" :label="$t('enterprise.header.badgeAgid')" class="q-py-xs q-px-sm text-weight-bold" />
        </div>
        <h1 class="text-h3 text-weight-bolder text-outfit q-my-sm text-gradient-premium">
          {{ $t('enterprise.header.title') }}
        </h1>
        <div class="text-subtitle1 text-slate-500">
          {{ $t('enterprise.header.subtitle') }}
        </div>
      </div>

      <!-- Action Button to Toggle Global Role Governance Matrix -->
      <div>
        <q-btn
          outline
          color="primary"
          icon="admin_panel_settings"
          :label="showRoleMatrix ? $t('enterprise.header.hideRoleMatrix') : $t('enterprise.header.viewRoleMatrix')"
          @click="showRoleMatrix = !showRoleMatrix"
          class="q-px-md shadow-sm"
        />
      </div>
    </div>

    <!-- Global Institutional Role Governance Matrix (Collapsible Panel) -->
    <q-slide-transition>
      <div v-if="showRoleMatrix" class="q-mb-xl">
        <q-card class="glass-card shadow-soft q-pa-md border-primary-soft">
          <div class="row items-center justify-between q-mb-md">
            <div>
              <div class="text-h6 text-weight-bold text-slate-800 row items-center q-gutter-x-sm">
                <q-icon name="verified_user" color="primary" size="sm" />
                <span>{{ $t('enterprise.matrix.title') }}</span>
              </div>
              <div class="text-caption text-slate-500">
                {{ $t('enterprise.matrix.subtitle') }}
              </div>
            </div>
            <q-btn flat round dense icon="close" @click="showRoleMatrix = false" />
          </div>

          <q-markup-table flat bordered dense class="rounded-borders overflow-hidden">
            <thead class="bg-slate-100 text-slate-700">
              <tr>
                <th class="text-left text-weight-bolder">{{ $t('enterprise.matrix.colFeature') }}</th>
                <th class="text-left text-weight-bolder">{{ $t('enterprise.matrix.colGovernance') }}</th>
                <th class="text-left text-weight-bolder">{{ $t('enterprise.matrix.colOperations') }}</th>
                <th class="text-left text-weight-bolder">{{ $t('enterprise.matrix.colEndUsers') }}</th>
                <th class="text-left text-weight-bolder">{{ $t('enterprise.matrix.colNorm') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.pagopa') }}</td>
                <td><q-badge color="deep-purple" :label="$t('enterprise.roles.ds')" class="q-mr-xs" /><q-badge color="primary" :label="$t('enterprise.roles.dsga')" /></td>
                <td><q-badge color="indigo-7" :label="$t('enterprise.roles.secretaryAccounts')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.parent')" class="q-mr-xs" /><q-badge color="teal-7" :label="$t('enterprise.roles.student')" /></td>
                <td class="text-caption text-slate-600">D.Lgs. 217/2017 & AgID</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.interpelli') }}</td>
                <td><q-badge color="deep-purple" :label="$t('enterprise.roles.ds')" /></td>
                <td><q-badge color="indigo-7" :label="$t('enterprise.roles.secretaryStaff')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.public')" /></td>
                <td class="text-caption text-slate-600">O.M. 88/2024 Art. 13 c. 23</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.albo') }}</td>
                <td><q-badge color="deep-purple" :label="$t('enterprise.roles.ds')" class="q-mr-xs" /><q-badge color="primary" :label="$t('enterprise.roles.dsga')" /></td>
                <td><q-badge color="indigo-7" :label="$t('enterprise.roles.secretary')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.public')" /></td>
                <td class="text-caption text-slate-600">L. 69/2009 & L. 190/2012 (ANAC)</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.feq') }}</td>
                <td><q-badge color="deep-purple" :label="$t('enterprise.roles.ds')" class="q-mr-xs" /><q-badge color="primary" :label="$t('enterprise.roles.dsga')" /></td>
                <td><q-badge color="indigo-7" :label="$t('enterprise.roles.secretaryDidactics')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.public')" /></td>
                <td class="text-caption text-slate-600">CAD Art. 23 (D.Lgs. 82/2005) & eIDAS</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.maturita') }}</td>
                <td><q-badge color="deep-purple" :label="$t('enterprise.roles.commission')" class="q-mr-xs" /><q-badge color="indigo-8" :label="$t('enterprise.roles.ds')" /></td>
                <td><q-badge color="blue-8" :label="$t('enterprise.roles.teacher')" class="q-mr-xs" /><q-badge color="indigo-7" :label="$t('enterprise.roles.secretaryDidactics')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.student')" /></td>
                <td class="text-caption text-slate-600">D.Lgs. 62/2017 & D.M. 88/2020</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.meals') }}</td>
                <td><q-badge color="indigo-7" :label="$t('enterprise.roles.secretary')" /></td>
                <td><q-badge color="blue-8" :label="$t('enterprise.roles.teacher')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.parent')" /></td>
                <td class="text-caption text-slate-600">Linee Guida Nazionali Ristorazione</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.inventory') }}</td>
                <td><q-badge color="primary" :label="$t('enterprise.roles.dsga')" /></td>
                <td><q-badge color="indigo-7" :label="$t('enterprise.roles.custodian')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.parent')" /></td>
                <td class="text-caption text-slate-600">D.I. 129/2018 (Regolamento Contabilità)</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.privacy') }}</td>
                <td><q-badge color="deep-orange-8" :label="$t('enterprise.roles.dpo')" class="q-mr-xs" /><q-badge color="deep-purple" :label="$t('enterprise.roles.ds')" /></td>
                <td><q-badge color="blue-8" :label="$t('enterprise.roles.teacher')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.parent')" /></td>
                <td class="text-caption text-slate-600">Reg. UE 2016/679 (GDPR) Art. 30</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.sidi') }}</td>
                <td><q-badge color="deep-purple" :label="$t('enterprise.roles.ds')" class="q-mr-xs" /><q-badge color="primary" :label="$t('enterprise.roles.dsga')" /></td>
                <td><q-badge color="indigo-7" :label="$t('enterprise.roles.secretaryDidactics')" class="q-mr-xs" /><q-badge color="blue-8" :label="$t('enterprise.roles.teacher')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.ds')" /></td>
                <td class="text-caption text-slate-600">Specifiche mTLS Cooperazione MIM</td>
              </tr>
              <tr>
                <td class="text-weight-bold text-primary">{{ $t('enterprise.tabs.psychology') }}</td>
                <td><q-badge color="pink-8" :label="$t('enterprise.roles.psychologist')" /></td>
                <td><q-badge color="pink-8" :label="$t('enterprise.roles.psychologist')" /></td>
                <td><q-badge color="teal-8" :label="$t('enterprise.roles.student')" class="q-mr-xs" /><q-badge color="deep-purple" :label="$t('enterprise.roles.parent')" /></td>
                <td class="text-caption text-slate-600">L. 56/1989 & Protocollo MIM-CNOP</td>
              </tr>
            </tbody>
          </q-markup-table>
        </q-card>
      </div>
    </q-slide-transition>

    <!-- Role Scope & Simulation Bar -->
    <div class="row items-center justify-between q-mb-md q-pa-sm bg-white rounded-borders border border-slate-200 shadow-sm">
      <div class="row items-center q-gutter-x-sm">
        <span class="text-caption text-weight-bold text-slate-700">Profilo Istituzionale:</span>
        <q-chip dense color="indigo-8" text-color="white" icon="badge">
          {{ formattedRoleName }}
        </q-chip>
        <q-chip dense outline color="primary" icon="visibility">
          {{ visibleTabs.length }} moduli autorizzati su 10
        </q-chip>
      </div>

      <!-- Role simulation selector for DS / Admin to preview any role's interface -->
      <div v-if="canSimulateRole" class="row items-center q-gutter-x-sm">
        <span class="text-caption text-slate-600">Simula vista ruolo:</span>
        <q-select
          v-model="simulatedRole"
          :options="simulatableRoles"
          dense
          outlined
          emit-value
          map-options
          style="min-width: 220px"
        />
      </div>
    </div>

    <!-- Canteen Activation Banner for Principal and DSGA -->
    <div v-if="canManageCanteen" class="q-mb-md">
      <q-card flat bordered class="q-pa-md bg-teal-50 border-teal-200 rounded-borders shadow-sm">
        <div class="row items-center justify-between">
          <div class="row items-center q-gutter-x-md">
            <q-avatar color="teal-8" text-color="white" icon="restaurant" size="42px" />
            <div>
              <div class="text-subtitle2 text-weight-bold text-teal-9">
                Servizio Mensa Scolastica d'Istituto:
                <q-badge :color="canteenActive ? 'positive' : 'grey-7'" class="q-ml-xs">
                  {{ canteenActive ? 'ATTIVO' : 'NON ATTIVO' }}
                </q-badge>
              </div>
              <div class="text-caption text-teal-8">
                La mensa non serve in tutte le scuole. Abilitando il servizio, comparirà per Genitori (borsellino & pasti), Studenti (menu refezione) e Segreteria (rilevazione presenze).
              </div>
            </div>
          </div>
          <div>
            <q-toggle
              v-model="canteenActive"
              color="teal-8"
              size="lg"
              :label="canteenActive ? 'Mensa Abilitata' : 'Mensa Disattivata'"
              class="text-weight-bold text-teal-9"
            />
          </div>
        </div>
      </q-card>
    </div>

    <!-- Navigation Tabs -->
    <q-card class="glass-card shadow-soft q-mb-lg overflow-hidden">
      <q-tabs
        v-model="activeTab"
        dense
        align="left"
        class="text-grey-7 bg-slate-50 border-bottom-soft"
        active-color="primary"
        indicator-color="primary"
        narrow-indicator
      >
        <q-tab v-if="canViewTab('pagopa')" name="pagopa" icon="payment" :label="$t('enterprise.tabs.pagopa')" />
        <q-tab v-if="canViewTab('interpelli')" name="interpelli" icon="campaign" :label="$t('enterprise.tabs.interpelli')" />
        <q-tab v-if="canViewTab('albo')" name="albo" icon="account_balance" :label="$t('enterprise.tabs.albo')" />
        <q-tab v-if="canViewTab('feq')" name="feq" icon="draw" :label="$t('enterprise.tabs.feq')" />
        <q-tab v-if="canViewTab('maturita')" name="maturita" icon="school" :label="$t('enterprise.tabs.maturita')" />
        <q-tab v-if="canViewTab('meals')" name="meals" icon="restaurant" :label="$t('enterprise.tabs.meals')" />
        <q-tab v-if="canViewTab('inventory')" name="inventory" icon="inventory_2" :label="$t('enterprise.tabs.inventory')" />
        <q-tab v-if="canViewTab('privacy')" name="privacy" icon="security" :label="$t('enterprise.tabs.privacy')" />
        <q-tab v-if="canViewTab('sidi')" name="sidi" icon="sync_alt" :label="$t('enterprise.tabs.sidi')" />
        <q-tab v-if="canViewTab('psychology')" name="psychology" icon="psychology" :label="$t('enterprise.tabs.psychology')" />
      </q-tabs>

      <q-tab-panels v-model="activeTab" animated class="q-pa-md bg-white">
        <!-- 1. PAGOPA & OPI -->
        <q-tab-panel v-if="canViewTab('pagopa')" name="pagopa">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="deep-purple" text-color="white" icon="verified_user">{{ $t('enterprise.roles.ds') }}</q-chip>
              <q-chip dense size="sm" color="primary" text-color="white" icon="account_balance">{{ $t('enterprise.roles.dsga') }}</q-chip>
              <q-chip dense size="sm" color="indigo-7" text-color="white" icon="receipt_long">{{ $t('enterprise.roles.secretaryAccounts') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="family_restroom">{{ $t('enterprise.roles.parent') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} D.Lgs. 217/2017 & AgID
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md full-height">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">💳 {{ $t('enterprise.pagopa.cardTitleNotice') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.pagopa.cardSubtitleNotice') }}</div>

                <!-- Mode selection: Single Student vs Entire Class -->
                <div class="q-mb-md">
                  <q-btn-toggle
                    v-model="pagoPaMode"
                    dense
                    rounded
                    unelevated
                    toggle-color="primary"
                    color="grey-3"
                    text-color="grey-8"
                    :options="[
                      { label: $t('enterprise.pagopa.modeSingle'), value: 'single', icon: 'person' },
                      { label: $t('enterprise.pagopa.modeClass'), value: 'class', icon: 'groups' }
                    ]"
                  />
                </div>

                <!-- Single Student Selector -->
                <div v-if="pagoPaMode === 'single'" class="q-mb-sm">
                  <q-select
                    v-model="selectedPagoPaStudent"
                    :options="filteredStudents"
                    use-input
                    fill-input
                    input-debounce="0"
                    @filter="filterStudents"
                    option-label="name"
                    :label="$t('enterprise.pagopa.selectStudent')"
                    outlined
                    dense
                    clearable
                  >
                    <template v-slot:option="scope">
                      <q-item v-bind="scope.itemProps">
                        <q-item-section avatar>
                          <q-avatar color="primary" text-color="white" size="28px">{{ scope.opt.last_name?.charAt(0) }}</q-avatar>
                        </q-item-section>
                        <q-item-section>
                          <q-item-label class="text-weight-bold">{{ scope.opt.last_name }} {{ scope.opt.first_name }}</q-item-label>
                          <q-item-label caption>Classe: {{ scope.opt.class_name }} • CF: {{ scope.opt.fiscal_code }}</q-item-label>
                        </q-item-section>
                        <q-item-section side>
                          <q-badge color="indigo-1" text-color="indigo-9" :label="scope.opt.class_name" />
                        </q-item-section>
                      </q-item>
                    </template>
                  </q-select>
                </div>

                <!-- Entire Class Selector -->
                <div v-else class="q-mb-sm">
                  <q-select
                    v-model="selectedPagoPaClass"
                    :options="classesList"
                    :label="$t('enterprise.pagopa.selectClass')"
                    outlined
                    dense
                  />
                  <div class="text-caption text-primary q-mt-xs text-weight-medium">
                    {{ studentsInPagoPaClass.length }} studenti associati a questa classe
                  </div>
                </div>

                <!-- Quick Presets -->
                <div class="q-mb-sm">
                  <div class="text-caption text-slate-500 q-mb-xs">{{ $t('enterprise.pagopa.quickPresets') }}</div>
                  <div class="row q-gutter-xs">
                    <q-chip clickable dense size="sm" color="blue-1" text-color="primary" @click="applyPagoPaPreset('Assicurazione Scolastica e Infortuni', 15.0)">
                      Assicurazione (15 €)
                    </q-chip>
                    <q-chip clickable dense size="sm" color="blue-1" text-color="primary" @click="applyPagoPaPreset('Laboratori Didattici PNRR', 30.0)">
                      Laboratori PNRR (30 €)
                    </q-chip>
                    <q-chip clickable dense size="sm" color="blue-1" text-color="primary" @click="applyPagoPaPreset('Contributo Volontario e Ampliamento PTOF', 85.0)">
                      Contributo (85 €)
                    </q-chip>
                    <q-chip clickable dense size="sm" color="blue-1" text-color="primary" @click="applyPagoPaPreset('Quota Viaggio di Istruzione / Gita', 120.0)">
                      Gita (120 €)
                    </q-chip>
                  </div>
                </div>

                <q-input v-model="pagoPaForm.title" :label="$t('enterprise.pagopa.title')" outlined dense class="q-mb-sm" />
                <q-input v-model.number="pagoPaForm.amount" type="number" step="0.5" :label="$t('enterprise.pagopa.amount')" outlined dense class="q-mb-md" />

                <!-- Action Button depending on mode -->
                <q-btn
                  v-if="pagoPaMode === 'single'"
                  color="primary"
                  icon="receipt_long"
                  :label="$t('enterprise.pagopa.generateNoticeBtn')"
                  @click="generatePagoPA"
                  :loading="loading"
                  unelevated
                />
                <q-btn
                  v-else
                  color="primary"
                  icon="batch_prediction"
                  :label="$t('enterprise.pagopa.generateBulkNoticeBtn', { count: studentsInPagoPaClass.length })"
                  @click="generateBulkPagoPA"
                  :loading="loading"
                  unelevated
                />

                <div v-if="generatedNotice" class="q-mt-md q-pa-sm bg-blue-50 rounded-borders border-blue-200">
                  <div class="text-weight-bold text-primary">{{ $t('enterprise.pagopa.iuvCode') }} {{ generatedNotice.iuv }}</div>
                  <div class="text-caption text-slate-600">{{ $t('enterprise.pagopa.qrPayload') }} {{ generatedNotice.qr_code_payload }}</div>
                  <q-btn flat dense color="primary" icon="download" :label="$t('enterprise.pagopa.downloadPdfBtn')" class="q-mt-xs" />
                </div>
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md full-height">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🏦 {{ $t('enterprise.pagopa.cardTitleOpi') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.pagopa.cardSubtitleOpi') }}</div>

                <q-input v-model="opiInput" type="textarea" rows="4" :label="$t('enterprise.pagopa.opiInputLabel')" outlined dense class="q-mb-md" />

                <div class="row q-gutter-sm">
                  <q-btn color="secondary" icon="upload_file" :label="$t('enterprise.pagopa.reconcileBtn')" @click="reconcileOPI" :loading="loading" unelevated />
                  <q-btn outline color="primary" icon="notifications_active" :label="$t('enterprise.pagopa.checkSollecitiBtn')" @click="checkSolleciti" />
                </div>

                <div v-if="reconciliationReport" class="q-mt-md q-pa-sm bg-green-50 rounded-borders text-green-900">
                  <div class="text-weight-bold">{{ $t('enterprise.pagopa.reconciledCount') }} {{ reconciliationReport.total_reconciled }} / {{ reconciliationReport.total_processed }}</div>
                  <div class="text-caption">{{ $t('enterprise.pagopa.totalAmount') }} € {{ reconciliationReport.total_amount?.toFixed(2) }}</div>
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 2. INTERPELLI OM 88/2024 -->
        <q-tab-panel v-if="canViewTab('interpelli')" name="interpelli">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="deep-purple" text-color="white" icon="school">{{ $t('enterprise.roles.ds') }}</q-chip>
              <q-chip dense size="sm" color="indigo-7" text-color="white" icon="badge">{{ $t('enterprise.roles.secretaryStaff') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="person_search">{{ $t('enterprise.roles.public') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} O.M. n. 88/2024 Art. 13
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">📢 {{ $t('enterprise.interpelli.cardTitlePublish') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.interpelli.cardSubtitlePublish') }}</div>

                <q-input v-model="interpelloForm.title" :label="$t('enterprise.interpelli.title')" outlined dense class="q-mb-sm" />
                <q-input v-model="interpelloForm.concorso_class" :label="$t('enterprise.interpelli.concorsoClass')" outlined dense class="q-mb-sm" />
                <q-input v-model.number="interpelloForm.weekly_hours" type="number" :label="$t('enterprise.interpelli.weeklyHours')" outlined dense class="q-mb-sm" />
                <q-input v-model="interpelloForm.deadline" :label="$t('enterprise.interpelli.deadline')" outlined dense class="q-mb-md" />

                <q-btn color="primary" icon="send" :label="$t('enterprise.interpelli.publishBtn')" @click="createInterpello" unelevated />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🏆 {{ $t('enterprise.interpelli.cardTitleGraduatoria') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.interpelli.cardSubtitleGraduatoria') }}</div>

                <q-btn outline color="primary" icon="format_list_numbered" :label="$t('enterprise.interpelli.loadGraduatoriaBtn')" @click="loadGraduatoria" class="q-mb-md" />

                <div v-if="graduatoriaList.length > 0">
                  <q-list separator bordered>
                    <q-item v-for="c in graduatoriaList" :key="c.id">
                      <q-item-section>
                        <q-item-label class="text-weight-bold">{{ c.candidate_name }} {{ c.candidate_surname }}</q-item-label>
                        <q-item-label caption>{{ $t('enterprise.interpelli.score', { score: c.total_score, cf: c.fiscal_code }) }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <q-btn dense size="sm" color="accent" icon="outgoing_mail" :label="$t('enterprise.interpelli.convocaBtn')" @click="convocaCandidate(c.id)" />
                      </q-item-section>
                    </q-item>
                  </q-list>
                </div>
                <div v-else class="text-caption text-slate-400 italic">{{ $t('enterprise.interpelli.noCandidates') }}</div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 3. ALBO PRETORIO & ANAC -->
        <q-tab-panel v-if="canViewTab('albo')" name="albo">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="deep-purple" text-color="white" icon="gavel">{{ $t('enterprise.roles.ds') }}</q-chip>
              <q-chip dense size="sm" color="primary" text-color="white" icon="account_balance">{{ $t('enterprise.roles.dsga') }}</q-chip>
              <q-chip dense size="sm" color="indigo-7" text-color="white" icon="history_edu">{{ $t('enterprise.roles.secretary') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="public">{{ $t('enterprise.roles.public') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} L. 69/2009 & L. 190/2012
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🏛️ {{ $t('enterprise.albo.cardTitleAlbo') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.albo.cardSubtitleAlbo') }}</div>

                <q-select v-model="alboForm.category" :options="['Bandi e Gare PNRR', 'Determine Dirigenziali', 'Decreti', 'Delibere Consiglio di Istituto']" :label="$t('enterprise.albo.category')" outlined dense class="q-mb-sm" />
                <q-input v-model="alboForm.subject" :label="$t('enterprise.albo.subject')" outlined dense class="q-mb-sm" />
                <q-input v-model="alboForm.cig_code" :label="$t('enterprise.albo.cigCode')" outlined dense class="q-mb-sm" />
                <q-input v-model.number="alboForm.awarded_amount" type="number" :label="$t('enterprise.albo.awardedAmount')" outlined dense class="q-mb-md" />

                <q-btn color="primary" icon="gavel" :label="$t('enterprise.albo.publishBtn')" @click="publishAlbo" unelevated />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">📑 {{ $t('enterprise.albo.cardTitleAnac') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.albo.cardSubtitleAnac') }}</div>

                <q-btn color="secondary" icon="download" :label="$t('enterprise.albo.downloadAnacBtn')" @click="downloadANAC" unelevated class="q-mb-md" />
                <div class="text-caption text-slate-600">{{ $t('enterprise.albo.anacDescription') }}</div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 4. FEQ CSC & DIGITAL STAMP -->
        <q-tab-panel v-if="canViewTab('feq')" name="feq">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="deep-purple" text-color="white" icon="draw">{{ $t('enterprise.roles.ds') }}</q-chip>
              <q-chip dense size="sm" color="primary" text-color="white" icon="verified">{{ $t('enterprise.roles.dsga') }}</q-chip>
              <q-chip dense size="sm" color="indigo-7" text-color="white" icon="description">{{ $t('enterprise.roles.secretaryDidactics') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="qr_code_scanner">{{ $t('enterprise.roles.public') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} CAD Art. 23 & eIDAS
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🔏 {{ $t('enterprise.feq.cardTitleCsc') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.feq.cardSubtitleCsc') }}</div>

                <q-input v-model="cscForm.pin" type="password" :label="$t('enterprise.feq.pin')" outlined dense class="q-mb-sm" />
                <q-input v-model="cscForm.otp" :label="$t('enterprise.feq.otp')" outlined dense class="q-mb-md" />

                <q-btn color="primary" icon="border_color" :label="$t('enterprise.feq.signBatchBtn')" @click="signBatchCSC" unelevated />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🔲 {{ $t('enterprise.feq.cardTitleGlifo') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.feq.cardSubtitleGlifo') }}</div>

                <q-input v-model="glifoToken" :label="$t('enterprise.feq.glifoToken')" outlined dense class="q-mb-md" />
                <q-btn color="secondary" icon="qr_code_scanner" :label="$t('enterprise.feq.verifyGlifoBtn')" @click="verifyGlifo" unelevated />

                <div v-if="glifoResult" class="q-mt-md q-pa-sm bg-green-50 rounded-borders text-green-900">
                  <div class="text-weight-bold">{{ $t('enterprise.feq.resultStatus') }} {{ glifoResult.status || 'CONFORME CAD' }}</div>
                  <div class="text-caption">{{ glifoResult.message }}</div>
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 5. MATURITA & CURRICULUM DELLO STUDENTE -->
        <q-tab-panel v-if="canViewTab('maturita')" name="maturita">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="deep-purple" text-color="white" icon="assignment_turned_in">{{ $t('enterprise.roles.commission') }}</q-chip>
              <q-chip dense size="sm" color="indigo-8" text-color="white" icon="school">{{ $t('enterprise.roles.ds') }}</q-chip>
              <q-chip dense size="sm" color="blue-8" text-color="white" icon="calculate">{{ $t('enterprise.roles.teacher') }}</q-chip>
              <q-chip dense size="sm" color="indigo-7" text-color="white" icon="upload_file">{{ $t('enterprise.roles.secretaryDidactics') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="school">{{ $t('enterprise.roles.student') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} D.Lgs. 62/2017 & D.M. 88/2020
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🎓 {{ $t('enterprise.maturita.cardTitleCredits') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.maturita.cardSubtitleCredits') }}</div>

                <div class="row q-col-gutter-sm q-mb-md">
                  <div class="col-4"><q-input v-model.number="maturitaCredits.grade_3rd" type="number" step="0.1" :label="$t('enterprise.maturita.grade3rd')" outlined dense /></div>
                  <div class="col-4"><q-input v-model.number="maturitaCredits.grade_4th" type="number" step="0.1" :label="$t('enterprise.maturita.grade4th')" outlined dense /></div>
                  <div class="col-4"><q-input v-model.number="maturitaCredits.grade_5th" type="number" step="0.1" :label="$t('enterprise.maturita.grade5th')" outlined dense /></div>
                </div>

                <q-btn color="primary" icon="calculate" :label="$t('enterprise.maturita.calcBtn')" @click="calcMaturitaCredits" unelevated />

                <div v-if="calculatedCredits" class="q-mt-md text-weight-bold text-primary">
                  {{ $t('enterprise.maturita.totalCredits', { credits: calculatedCredits.total_credits }) }}
                </div>
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">📜 {{ $t('enterprise.maturita.cardTitleCurriculum') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.maturita.cardSubtitleCurriculum') }}</div>

                <!-- Easy Student Selector for Graduating Students -->
                <q-select
                  v-model="selectedCurriculumStudent"
                  :options="graduatingStudents"
                  use-input
                  fill-input
                  input-debounce="0"
                  @filter="filterGraduatingStudents"
                  option-label="name"
                  :label="$t('enterprise.maturita.selectStudent')"
                  outlined
                  dense
                  clearable
                  class="q-mb-md"
                >
                  <template v-slot:option="scope">
                    <q-item v-bind="scope.itemProps">
                      <q-item-section avatar>
                        <q-avatar color="indigo-8" text-color="white" size="28px">{{ scope.opt.last_name?.charAt(0) }}</q-avatar>
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-bold">{{ scope.opt.last_name }} {{ scope.opt.first_name }}</q-item-label>
                        <q-item-label caption>Classe: {{ scope.opt.class_name }} • CF: {{ scope.opt.fiscal_code }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <q-badge color="purple-1" text-color="purple-9" :label="'Maturando ' + scope.opt.class_name" />
                      </q-item-section>
                    </q-item>
                  </template>
                </q-select>

                <q-btn color="secondary" icon="file_download" :label="$t('enterprise.maturita.downloadCurriculumBtn')" @click="downloadCurriculum" unelevated />
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 6. MENSA & DIETE -->
        <q-tab-panel v-if="canViewTab('meals')" name="meals">
          <!-- Canteen Inactive Warning when managed by DS/DSGA -->
          <div v-if="!canteenActive" class="q-pa-md bg-amber-50 border border-amber-200 rounded-borders q-mb-md">
            <div class="row items-center q-gutter-x-sm text-amber-9 text-weight-bold">
              <q-icon name="warning" size="20px" />
              <span>Servizio Mensa Scolastica Non Attivo per questo istituto</span>
            </div>
            <div class="text-caption text-amber-8 q-mt-xs">
              La mensa è attualmente disattivata. Famiglie, studenti e docenti non visualizzeranno le relative voci nel menu finché non viene abilitata dalla Dirigente o dal DSGA.
            </div>
          </div>
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="indigo-7" text-color="white" icon="restaurant">{{ $t('enterprise.roles.secretary') }}</q-chip>
              <q-chip dense size="sm" color="blue-8" text-color="white" icon="checklist">{{ $t('enterprise.roles.teacher') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="wallet">{{ $t('enterprise.roles.parent') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} Linee Guida Ristorazione Scolastica
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🍲 {{ $t('enterprise.meals.cardTitleRollCall') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.meals.cardSubtitleRollCall') }}</div>

                <q-select
                  v-model="mealRollCall.classId"
                  :options="classesList"
                  :label="$t('enterprise.meals.classId')"
                  outlined
                  dense
                  class="q-mb-sm"
                />
                <q-input v-model="mealRollCall.date" :label="$t('enterprise.meals.date')" outlined dense class="q-mb-md" />

                <q-btn color="primary" icon="checklist" :label="$t('enterprise.meals.submitRollCallBtn')" @click="submitMealRollCall" unelevated />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">👛 {{ $t('enterprise.meals.cardTitleWallet') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.meals.cardSubtitleWallet') }}</div>

                <!-- Easy Student Selector for Canteen Wallet -->
                <q-select
                  v-model="selectedMealStudent"
                  :options="filteredStudents"
                  use-input
                  fill-input
                  input-debounce="0"
                  @filter="filterStudents"
                  option-label="name"
                  :label="$t('enterprise.meals.selectStudent')"
                  outlined
                  dense
                  clearable
                  class="q-mb-sm"
                >
                  <template v-slot:option="scope">
                    <q-item v-bind="scope.itemProps">
                      <q-item-section avatar>
                        <q-avatar color="teal-7" text-color="white" size="28px">{{ scope.opt.last_name?.charAt(0) }}</q-avatar>
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-bold">{{ scope.opt.last_name }} {{ scope.opt.first_name }}</q-item-label>
                        <q-item-label caption>Classe: {{ scope.opt.class_name }} • Saldo: € {{ scope.opt.wallet_balance?.toFixed(2) }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <q-badge color="teal-1" text-color="teal-9" :label="'€ ' + scope.opt.wallet_balance?.toFixed(2)" />
                      </q-item-section>
                    </q-item>
                  </template>
                </q-select>

                <!-- Current Balance Display -->
                <div class="row items-center justify-between q-pa-sm bg-teal-50 rounded-borders text-teal-900 q-mb-sm">
                  <span class="text-caption text-weight-bold">{{ $t('enterprise.meals.currentBalance') }}</span>
                  <span class="text-subtitle2 text-weight-bolder">€ {{ selectedMealStudent?.wallet_balance?.toFixed(2) || '0.00' }}</span>
                </div>

                <!-- Quick Topup Amount Chips -->
                <div class="q-mb-sm">
                  <span class="text-caption text-slate-500 q-mr-xs">{{ $t('enterprise.meals.quickTopup') }}</span>
                  <q-chip clickable dense size="sm" color="teal-1" text-color="teal-9" @click="mealWallet.topupAmount = 10.0">+10 €</q-chip>
                  <q-chip clickable dense size="sm" color="teal-1" text-color="teal-9" @click="mealWallet.topupAmount = 25.0">+25 €</q-chip>
                  <q-chip clickable dense size="sm" color="teal-1" text-color="teal-9" @click="mealWallet.topupAmount = 50.0">+50 €</q-chip>
                  <q-chip clickable dense size="sm" color="teal-1" text-color="teal-9" @click="mealWallet.topupAmount = 100.0">+100 €</q-chip>
                </div>

                <q-input v-model.number="mealWallet.topupAmount" type="number" step="5" :label="$t('enterprise.meals.topupAmount')" outlined dense class="q-mb-md" />

                <q-btn color="secondary" icon="add_card" :label="$t('enterprise.meals.topupBtn')" @click="topupMealWallet" unelevated />
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 7. INVENTARIO & COMODATO -->
        <q-tab-panel v-if="canViewTab('inventory')" name="inventory">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="primary" text-color="white" icon="account_balance">{{ $t('enterprise.roles.dsga') }}</q-chip>
              <q-chip dense size="sm" color="indigo-7" text-color="white" icon="inventory">{{ $t('enterprise.roles.custodian') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="handshake">{{ $t('enterprise.roles.parent') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} D.I. 129/2018
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">📦 {{ $t('enterprise.inventory.cardTitleAsset') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.inventory.cardSubtitleAsset') }}</div>

                <q-input v-model="assetForm.description" :label="$t('enterprise.inventory.description')" outlined dense class="q-mb-sm" />
                <q-input v-model="assetForm.category" :label="$t('enterprise.inventory.category')" outlined dense class="q-mb-sm" />
                <q-input v-model.number="assetForm.initial_value" type="number" :label="$t('enterprise.inventory.initialValue')" outlined dense class="q-mb-md" />

                <q-btn color="primary" icon="add_box" :label="$t('enterprise.inventory.createAssetBtn')" @click="createAsset" unelevated />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">💻 {{ $t('enterprise.inventory.cardTitleLoan') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.inventory.cardSubtitleLoan') }}</div>

                <q-input v-model="loanForm.asset_id" :label="$t('enterprise.inventory.assetId')" outlined dense class="q-mb-sm" />
                <q-select
                  v-model="selectedLoanStudent"
                  :options="filteredStudents"
                  use-input
                  fill-input
                  input-debounce="0"
                  @filter="filterStudents"
                  option-label="name"
                  :label="$t('enterprise.inventory.studentId')"
                  outlined
                  dense
                  clearable
                  class="q-mb-md"
                />

                <q-btn color="secondary" icon="handshake" :label="$t('enterprise.inventory.createLoanBtn')" @click="createDeviceLoan" unelevated />
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 8. PRIVACY & SEMAFORO -->
        <q-tab-panel v-if="canViewTab('privacy')" name="privacy">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="deep-orange-8" text-color="white" icon="security">{{ $t('enterprise.roles.dpo') }}</q-chip>
              <q-chip dense size="sm" color="deep-purple" text-color="white" icon="shield">{{ $t('enterprise.roles.ds') }}</q-chip>
              <q-chip dense size="sm" color="blue-8" text-color="white" icon="traffic">{{ $t('enterprise.roles.teacher') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="fact_check">{{ $t('enterprise.roles.parent') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} Reg. UE 2016/679 (GDPR) Art. 30
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🛡️ {{ $t('enterprise.privacy.cardTitleTreatment') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.privacy.cardSubtitleTreatment') }}</div>

                <q-input v-model="treatmentForm.activity_name" :label="$t('enterprise.privacy.activityName')" outlined dense class="q-mb-sm" />
                <q-input v-model="treatmentForm.legal_basis" :label="$t('enterprise.privacy.legalBasis')" outlined dense class="q-mb-md" />

                <q-btn color="primary" icon="verified_user" :label="$t('enterprise.privacy.saveTreatmentBtn')" @click="saveTreatment" unelevated />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🚦 {{ $t('enterprise.privacy.cardTitleTrafficLight') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.privacy.cardSubtitleTrafficLight') }}</div>

                <div class="row items-center q-gutter-x-md q-my-md">
                  <q-badge color="positive" text-color="white" :label="$t('enterprise.privacy.greenBadge')" class="q-pa-sm" />
                  <q-badge color="warning" text-color="dark" :label="$t('enterprise.privacy.yellowBadge')" class="q-pa-sm" />
                  <q-badge color="negative" text-color="white" :label="$t('enterprise.privacy.redBadge')" class="q-pa-sm" />
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 9. SIDI MIM COOPERAZIONE -->
        <q-tab-panel v-if="canViewTab('sidi')" name="sidi">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="deep-purple" text-color="white" icon="key">{{ $t('enterprise.roles.ds') }}</q-chip>
              <q-chip dense size="sm" color="primary" text-color="white" icon="account_balance">{{ $t('enterprise.roles.dsga') }}</q-chip>
              <q-chip dense size="sm" color="indigo-7" text-color="white" icon="sync_alt">{{ $t('enterprise.roles.secretaryDidactics') }}</q-chip>
              <q-chip dense size="sm" color="blue-8" text-color="white" icon="cloud_upload">{{ $t('enterprise.roles.teacher') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} WebService mTLS MIM
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🔌 {{ $t('enterprise.sidi.cardTitleCert') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.sidi.cardSubtitleCert') }}</div>

                <div class="q-pa-sm bg-grey-1 rounded-borders text-caption text-slate-700 q-mb-md">
                  <div><strong>{{ $t('enterprise.sidi.endpoint') }}</strong> https://cooperazione.pubblica.istruzione.it/ws/sidi/v2</div>
                  <div><strong>{{ $t('enterprise.sidi.connectionStatus') }}</strong> <span class="text-positive text-weight-bold">{{ $t('enterprise.sidi.connected') }}</span></div>
                </div>

                <q-btn color="primary" icon="sync" :label="$t('enterprise.sidi.syncStudentsBtn')" @click="syncSIDIStudents" unelevated />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">📤 {{ $t('enterprise.sidi.cardTitleScrutiny') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.sidi.cardSubtitleScrutiny') }}</div>

                <q-btn color="secondary" icon="cloud_upload" :label="$t('enterprise.sidi.pushScrutinyBtn')" @click="pushScrutinyToMIM" unelevated />
              </q-card>
            </div>
          </div>
        </q-tab-panel>

        <!-- 10. SPORTELLO PSICOLOGICO CIC -->
        <q-tab-panel v-if="canViewTab('psychology')" name="psychology">
          <!-- Authorized Roles Bar -->
          <div class="roles-governance-bar q-pa-sm q-mb-md rounded-borders bg-slate-50 row items-center justify-between">
            <div class="row items-center q-gutter-x-sm">
              <span class="text-caption text-weight-bold text-slate-700">{{ $t('enterprise.governance.authorizedRoles') }}</span>
              <q-chip dense size="sm" color="pink-8" text-color="white" icon="psychology">{{ $t('enterprise.roles.psychologist') }}</q-chip>
              <q-chip dense size="sm" color="teal-8" text-color="white" icon="lock">{{ $t('enterprise.roles.student') }}</q-chip>
              <q-chip dense size="sm" color="deep-purple" text-color="white" icon="edit_document">{{ $t('enterprise.roles.parent') }}</q-chip>
            </div>
            <div class="text-caption text-slate-500 font-mono">
              {{ $t('enterprise.governance.legalRef') }} L. 56/1989 & CNOP
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md full-height">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">🧠 {{ $t('enterprise.psychology.cardTitleBooking') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.psychology.cardSubtitleBooking') }}</div>

                <q-input v-model="cicBooking.slotTime" :label="$t('enterprise.psychology.slotTime')" outlined dense class="q-mb-md" />

                <q-btn color="primary" icon="lock" :label="$t('enterprise.psychology.bookBtn')" @click="bookCICSession" unelevated />
              </q-card>
            </div>

            <div class="col-12 col-md-6">
              <q-card flat bordered class="q-pa-md full-height">
                <div class="text-h6 text-weight-bold text-slate-800 q-mb-sm">✍️ {{ $t('enterprise.psychology.cardTitleConsent') }}</div>
                <div class="text-caption text-slate-500 q-mb-md">{{ $t('enterprise.psychology.cardSubtitleConsent') }}</div>

                <!-- Mode selection: Single Student vs Entire Class -->
                <div class="q-mb-md">
                  <q-btn-toggle
                    v-model="cicMode"
                    dense
                    rounded
                    unelevated
                    toggle-color="primary"
                    color="grey-3"
                    text-color="grey-8"
                    :options="[
                      { label: $t('enterprise.psychology.modeSingle'), value: 'single', icon: 'person' },
                      { label: $t('enterprise.psychology.modeClass'), value: 'class', icon: 'groups' }
                    ]"
                  />
                </div>

                <!-- Single Student Selector -->
                <div v-if="cicMode === 'single'" class="q-mb-md">
                  <q-select
                    v-model="selectedCicStudent"
                    :options="filteredStudents"
                    use-input
                    fill-input
                    input-debounce="0"
                    @filter="filterStudents"
                    option-label="name"
                    :label="$t('enterprise.psychology.selectStudent')"
                    outlined
                    dense
                    clearable
                    class="q-mb-sm"
                  >
                    <template v-slot:option="scope">
                      <q-item v-bind="scope.itemProps">
                        <q-item-section avatar>
                          <q-avatar color="deep-purple-7" text-color="white" size="28px">{{ scope.opt.last_name?.charAt(0) }}</q-avatar>
                        </q-item-section>
                        <q-item-section>
                          <q-item-label class="text-weight-bold">{{ scope.opt.last_name }} {{ scope.opt.first_name }}</q-item-label>
                          <q-item-label caption>Classe: {{ scope.opt.class_name }}</q-item-label>
                        </q-item-section>
                        <q-item-section side>
                          <q-badge :color="scope.opt.has_cic_consent ? 'positive' : 'warning'" :label="scope.opt.has_cic_consent ? $t('enterprise.psychology.statusSigned') : $t('enterprise.psychology.statusPending')" />
                        </q-item-section>
                      </q-item>
                    </template>
                  </q-select>

                  <div v-if="selectedCicStudent" class="q-mb-sm">
                    <q-badge
                      :color="selectedCicStudent.has_cic_consent ? 'positive' : 'warning'"
                      class="q-pa-xs"
                      :label="selectedCicStudent.has_cic_consent ? $t('enterprise.psychology.statusSigned') : $t('enterprise.psychology.statusPending')"
                    />
                  </div>

                  <q-btn color="secondary" icon="fact_check" :label="$t('enterprise.psychology.signConsentBtn')" @click="signCICConsent" unelevated />
                </div>

                <!-- Entire Class Mode -->
                <div v-else class="q-mb-md">
                  <q-select
                    v-model="selectedCicClass"
                    :options="classesList"
                    :label="$t('enterprise.psychology.selectClass')"
                    outlined
                    dense
                    class="q-mb-sm"
                  />

                  <div class="q-mb-md">
                    <div class="text-caption text-slate-600 q-mb-xs text-weight-bold">{{ $t('enterprise.psychology.studentsInClass') }}</div>
                    <q-list dense bordered separator class="rounded-borders" style="max-height: 150px; overflow-y: auto;">
                      <q-item v-for="st in studentsInCicClass" :key="st.id">
                        <q-item-section>
                          <q-item-label class="text-weight-medium">{{ st.last_name }} {{ st.first_name }}</q-item-label>
                        </q-item-section>
                        <q-item-section side>
                          <q-badge :color="st.has_cic_consent ? 'positive' : 'warning'" :label="st.has_cic_consent ? $t('enterprise.psychology.statusSigned') : $t('enterprise.psychology.statusPending')" />
                        </q-item-section>
                      </q-item>
                    </q-list>
                  </div>

                  <q-btn
                    color="primary"
                    icon="send"
                    :label="$t('enterprise.psychology.sendBulkConsentBtn', { count: studentsInCicClass.length })"
                    @click="signBulkCICConsent"
                    :loading="loading"
                    unelevated
                  />
                </div>
              </q-card>
            </div>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { paymentService } from 'src/services/paymentService'
import { interpelliService } from 'src/services/interpelliService'
import { alboPretorioService } from 'src/services/alboPretorioService'
import { digitalStampService } from 'src/services/digitalStampService'
import { maturitaService } from 'src/services/maturitaService'
import { mealsService } from 'src/services/mealsService'
import { inventoryService } from 'src/services/inventoryService'
import { privacyService } from 'src/services/privacyService'
import { sidiService } from 'src/services/sidiService'
import { psychologyService } from 'src/services/psychologyService'
import { userService } from 'src/services/userService'
import { useAuthStore } from 'src/stores/auth'
import { useSchoolSettingsStore } from 'src/stores/schoolSettings'

const $q = useQuasar()
const { t } = useI18n()

// Safe Store Instances (with fallback if Pinia not in test context)
let authStoreInstance = null
try {
  authStoreInstance = useAuthStore()
} catch {
  authStoreInstance = null
}

let schoolSettingsInstance = null
try {
  schoolSettingsInstance = useSchoolSettingsStore()
} catch {
  schoolSettingsInstance = null
}

const actualUserRole = computed(() => {
  const role = (authStoreInstance?.userRole || authStoreInstance?.user?.role || '').toLowerCase()
  return role || 'admin'
})

// Simulation selector for Admin and Principal
const canSimulateRole = computed(() => {
  return ['superadmin', 'admin', 'principal', 'dirigente_scolastico'].includes(actualUserRole.value)
})

const simulatedRole = ref('')
const effectiveRole = computed(() => {
  if (simulatedRole.value) {
    return simulatedRole.value
  }
  return actualUserRole.value
})

const simulatableRoles = [
  { label: 'Usa Ruolo Reale', value: '' },
  { label: 'Segreteria Didattica & Alunni', value: 'secretary' },
  { label: 'DSGA (Direttore SGA)', value: 'dsga' },
  { label: 'Dirigente Scolastico (DS)', value: 'principal' },
  { label: 'Segreteria Personale', value: 'assistente_personale' },
  { label: 'Segreteria Contabilità', value: 'assistente_contabilita' },
  { label: 'Data Protection Officer (DPO)', value: 'dpo' },
  { label: 'Psicologo Scolastico (CIC)', value: 'psychologist' },
  { label: 'Docente', value: 'teacher' }
]

const roleTabPermissions = {
  superadmin: ['pagopa', 'interpelli', 'albo', 'feq', 'maturita', 'meals', 'inventory', 'privacy', 'sidi', 'psychology'],
  admin: ['pagopa', 'interpelli', 'albo', 'feq', 'maturita', 'meals', 'inventory', 'privacy', 'sidi', 'psychology'],
  principal: ['pagopa', 'interpelli', 'albo', 'feq', 'maturita', 'meals', 'inventory', 'privacy', 'sidi', 'psychology'],
  dirigente_scolastico: ['pagopa', 'interpelli', 'albo', 'feq', 'maturita', 'meals', 'inventory', 'privacy', 'sidi', 'psychology'],
  dsga: ['pagopa', 'albo', 'feq', 'meals', 'inventory', 'sidi'],
  secretary: ['pagopa', 'maturita', 'sidi', 'meals'],
  assistente_alunni: ['pagopa', 'maturita', 'sidi', 'meals'],
  assistente_amministrativo: ['pagopa', 'maturita', 'sidi', 'meals'],
  assistente_contabilita: ['pagopa', 'inventory', 'sidi', 'meals'],
  assistente_personale: ['interpelli', 'sidi'],
  dpo: ['privacy'],
  psychologist: ['psychology'],
  psicologo: ['psychology'],
  teacher: ['maturita', 'meals'],
  docente: ['maturita', 'meals'],
  student: ['psychology', 'meals'],
  studente: ['psychology', 'meals'],
  parent: ['pagopa', 'psychology', 'meals'],
  genitore: ['pagopa', 'psychology', 'meals']
}

const localCanteenEnabled = ref(true)
const canteenActive = computed({
  get: () => {
    if (schoolSettingsInstance) {
      return schoolSettingsInstance.canteenEnabled
    }
    return localCanteenEnabled.value
  },
  set: (val) => {
    localCanteenEnabled.value = val
    if (schoolSettingsInstance) {
      schoolSettingsInstance.setCanteenEnabled(val)
    }
  }
})

const canManageCanteen = computed(() => {
  return ['superadmin', 'admin', 'principal', 'dirigente_scolastico', 'dsga'].includes(effectiveRole.value)
})

function canViewTab(tabName) {
  const role = effectiveRole.value
  const allowed = roleTabPermissions[role] || ['pagopa', 'interpelli', 'albo', 'feq', 'maturita', 'meals', 'inventory', 'privacy', 'sidi', 'psychology']
  if (!allowed.includes(tabName)) return false

  // If canteen is inactive, only DS/DSGA/Admin can view the tab to manage/enable it
  if (tabName === 'meals' && !canteenActive.value && !canManageCanteen.value) {
    return false
  }

  return true
}

const allTabKeys = ['pagopa', 'interpelli', 'albo', 'feq', 'maturita', 'meals', 'inventory', 'privacy', 'sidi', 'psychology']
const visibleTabs = computed(() => {
  return allTabKeys.filter(t => canViewTab(t))
})

const activeTab = ref('pagopa')
const showRoleMatrix = ref(false)
const loading = ref(false)

watch(visibleTabs, (tabs) => {
  if (tabs.length > 0 && !tabs.includes(activeTab.value)) {
    activeTab.value = tabs[0]
  }
}, { immediate: true })

const formattedRoleName = computed(() => {
  const role = effectiveRole.value
  const map = {
    secretary: 'Segreteria Didattica & Alunni',
    assistente_alunni: 'Assistente Alunni',
    assistente_amministrativo: 'Assistente Amministrativo',
    assistente_contabilita: 'Segreteria Contabilità',
    assistente_personale: 'Segreteria Personale',
    dsga: 'DSGA (Direttore SGA)',
    principal: 'Dirigente Scolastico (DS)',
    dirigente_scolastico: 'Dirigente Scolastico (DS)',
    admin: 'Amministratore di Sistema',
    superadmin: 'Super Amministratore',
    dpo: 'Data Protection Officer (DPO)',
    psychologist: 'Psicologo Scolastico (CIC)',
    psicologo: 'Psicologo Scolastico (CIC)',
    teacher: 'Docente',
    docente: 'Docente'
  }
  return map[role] || role.toUpperCase()
})

// Student and Class Directory Data
const classesList = ref(['5A', '5B', '4A', '3A', '2A', '1A'])
const defaultStudents = [
  { id: 'student-1', first_name: 'Marco', last_name: 'Rossi', class_name: '5A', fiscal_code: 'RSSMRC08A01H501Z', email: 'marco.rossi@scuola.edu.it', wallet_balance: 34.50, has_cic_consent: true },
  { id: 'student-mat-1', first_name: 'Giulia', last_name: 'Bianchi', class_name: '5A', fiscal_code: 'BNCGLI08C41H501Y', email: 'giulia.bianchi@scuola.edu.it', wallet_balance: 18.00, has_cic_consent: true },
  { id: 'student-minor-1', first_name: 'Alessandro', last_name: 'Verdi', class_name: '2A', fiscal_code: 'VRDLSN10M12H501K', email: 'alessandro.verdi@scuola.edu.it', wallet_balance: 5.50, has_cic_consent: false },
  { id: 'student-4', first_name: 'Chiara', last_name: 'Esposito', class_name: '5B', fiscal_code: 'SPTCRA08T55F205W', email: 'chiara.esposito@scuola.edu.it', wallet_balance: 42.00, has_cic_consent: true },
  { id: 'student-5', first_name: 'Luca', last_name: 'Ferrari', class_name: '2A', fiscal_code: 'FRRLCU10E14F205J', email: 'luca.ferrari@scuola.edu.it', wallet_balance: 12.00, has_cic_consent: false },
  { id: 'student-6', first_name: 'Sofia', last_name: 'Romano', class_name: '1A', fiscal_code: 'RMNSFO11L50H501A', email: 'sofia.romano@scuola.edu.it', wallet_balance: 20.00, has_cic_consent: false }
]
const allStudents = ref(defaultStudents.map(s => ({ ...s, name: `${s.last_name} ${s.first_name} (${s.class_name})` })))
const filteredStudents = ref([...allStudents.value])

function filterStudents(val, update) {
  if (val === '') {
    update(() => {
      filteredStudents.value = allStudents.value
    })
    return
  }
  update(() => {
    const needle = val.toLowerCase()
    filteredStudents.value = allStudents.value.filter(s =>
      s.name.toLowerCase().includes(needle) ||
      s.fiscal_code.toLowerCase().includes(needle) ||
      s.class_name.toLowerCase().includes(needle)
    )
  })
}

const graduatingStudents = computed(() => {
  return allStudents.value.filter(s => s.class_name.startsWith('5'))
})
const filteredGraduatingStudents = ref([...graduatingStudents.value])
function filterGraduatingStudents(val, update) {
  if (val === '') {
    update(() => {
      filteredGraduatingStudents.value = graduatingStudents.value
    })
    return
  }
  update(() => {
    const needle = val.toLowerCase()
    filteredGraduatingStudents.value = graduatingStudents.value.filter(s =>
      s.name.toLowerCase().includes(needle) ||
      s.fiscal_code.toLowerCase().includes(needle)
    )
  })
}

// 1. PagoPA
const pagoPaMode = ref('single')
const selectedPagoPaStudent = ref(allStudents.value[0])
const selectedPagoPaClass = ref('5A')
const studentsInPagoPaClass = computed(() => allStudents.value.filter(s => s.class_name === selectedPagoPaClass.value))

const pagoPaForm = reactive({
  studentId: 'student-1',
  title: 'Contributo Volontario e Assicurazione',
  amount: 85.0
})
const generatedNotice = ref(null)
const opiInput = ref('')
const reconciliationReport = ref(null)

watch(selectedPagoPaStudent, (st) => {
  if (st) pagoPaForm.studentId = st.id
})

function applyPagoPaPreset(title, amount) {
  pagoPaForm.title = title
  pagoPaForm.amount = amount
}

async function generatePagoPA() {
  try {
    loading.value = true
    const res = await paymentService.createPayment({
      student_id: pagoPaForm.studentId,
      title: pagoPaForm.title,
      amount: pagoPaForm.amount,
      due_date: '2026-11-30'
    })
    const bollettino = await paymentService.getBollettino(res.data.id || 'pay-1')
    generatedNotice.value = bollettino.data
    $q.notify({ type: 'positive', message: t('enterprise.pagopa.successNotice') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  } finally {
    loading.value = false
  }
}

async function generateBulkPagoPA() {
  try {
    loading.value = true
    const list = studentsInPagoPaClass.value
    if (list.length === 0) {
      $q.notify({ type: 'warning', message: 'Nessuno studente trovato per la classe selezionata.' })
      return
    }
    await paymentService.createBulkPayments(list.map(s => s.id), {
      title: pagoPaForm.title,
      amount: pagoPaForm.amount,
      due_date: '2026-11-30'
    })
    $q.notify({
      type: 'positive',
      message: t('enterprise.pagopa.successBulkNotice', { count: list.length })
    })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  } finally {
    loading.value = false
  }
}

async function reconcileOPI() {
  try {
    loading.value = true
    const res = await paymentService.reconcileOPI(opiInput.value || '<xml></xml>', 'OPI_XML')
    reconciliationReport.value = res.data
    $q.notify({ type: 'positive', message: t('enterprise.pagopa.successReconcile') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  } finally {
    loading.value = false
  }
}

async function checkSolleciti() {
  try {
    const res = await paymentService.getSolleciti()
    const count = res.data?.total || 0
    $q.notify({ type: 'info', message: t('enterprise.pagopa.sollecitiFound', { count }) })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 2. Interpelli
const interpelloForm = reactive({ title: 'Supplenza A026 Matematica', concorso_class: 'A026', weekly_hours: 18, deadline: '2026-10-31' })
const graduatoriaList = ref([])

async function createInterpello() {
  try {
    await interpelliService.createNotice({
      title: interpelloForm.title,
      concorso_class: interpelloForm.concorso_class,
      weekly_hours: interpelloForm.weekly_hours,
      start_date: '2026-11-01',
      end_date: '2027-06-30',
      deadline: interpelloForm.deadline
    })
    $q.notify({ type: 'positive', message: t('enterprise.interpelli.successPublished') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

async function loadGraduatoria() {
  try {
    const res = await interpelliService.getGraduatoria('notice-1')
    graduatoriaList.value = res.data?.graduatoria || []
  } catch (err) {
    $q.notify({ type: 'warning', message: t('enterprise.interpelli.noActiveGraduatoria') })
  }
}

async function convocaCandidate(id) {
  try {
    await interpelliService.convoca(id, 24)
    $q.notify({ type: 'positive', message: t('enterprise.interpelli.successConvoca') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 3. Albo Pretorio
const alboForm = reactive({ category: 'Bandi e Gare PNRR', subject: 'Affidamento dotazioni laboratori 4.0', cig_code: 'B9912014A', awarded_amount: 32000.0 })

async function publishAlbo() {
  try {
    await alboPretorioService.publishAct({
      category: alboForm.category,
      subject: alboForm.subject,
      document_file_url: 'https://storage.scuola.edu.it/atti/atto.pdf',
      document_sha256: 'sha256-mock-hash-1',
      days_duration: 15,
      is_transparency_section: true,
      cig_code: alboForm.cig_code,
      awarded_amount: alboForm.awarded_amount
    })
    $q.notify({ type: 'positive', message: t('enterprise.albo.successPublished') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

async function downloadANAC() {
  try {
    await alboPretorioService.downloadANACXML({ anno: 2026 })
    $q.notify({ type: 'positive', message: t('enterprise.albo.successAnac') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 4. FEQ & Digital Stamp
const cscForm = reactive({ pin: '', otp: '' })
const glifoToken = ref('')
const glifoResult = ref(null)

async function signBatchCSC() {
  try {
    await digitalStampService.batchSignCSC(['pagella-1', 'pagella-2'], cscForm.pin, cscForm.otp)
    $q.notify({ type: 'positive', message: t('enterprise.feq.successSigned') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

async function verifyGlifo() {
  try {
    const res = await digitalStampService.verifyPublicToken(glifoToken.value || 'CAD-GLIFO-DEMO')
    glifoResult.value = res.data
    $q.notify({ type: 'positive', message: t('enterprise.feq.successGlifoVerified') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 5. Maturità
const maturitaCredits = reactive({ grade_3rd: 8.4, grade_4th: 8.6, grade_5th: 8.8 })
const calculatedCredits = ref(null)
const selectedCurriculumStudent = ref(allStudents.value.find(s => s.id === 'student-mat-1') || allStudents.value[1])
const curriculumStudentId = ref('student-mat-1')

watch(selectedCurriculumStudent, (st) => {
  if (st) curriculumStudentId.value = st.id
})

async function calcMaturitaCredits() {
  try {
    const res = await maturitaService.calculateCredits({
      grade_3rd: maturitaCredits.grade_3rd,
      high_3rd: true,
      grade_4th: maturitaCredits.grade_4th,
      high_4th: true,
      grade_5th: maturitaCredits.grade_5th,
      high_5th: true
    })
    calculatedCredits.value = res.data
    $q.notify({ type: 'positive', message: t('enterprise.maturita.successCalc') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

async function downloadCurriculum() {
  try {
    await maturitaService.downloadCurriculumXML({ student_id: curriculumStudentId.value })
    $q.notify({ type: 'positive', message: t('enterprise.maturita.successCurriculum') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 6. Mensa
const mealRollCall = reactive({ classId: '2A', date: '2026-10-07' })
const selectedMealStudent = ref(allStudents.value[0])
const mealWallet = reactive({ studentId: 'student-1', topupAmount: 50.0 })

watch(selectedMealStudent, (st) => {
  if (st) mealWallet.studentId = st.id
})

async function submitMealRollCall() {
  try {
    await mealsService.recordRollCall(mealRollCall.classId, mealRollCall.date, ['st-1', 'st-2'])
    $q.notify({ type: 'positive', message: t('enterprise.meals.successRollCall') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

async function topupMealWallet() {
  try {
    await mealsService.topUpWallet(mealWallet.studentId, mealWallet.topupAmount, '01210000000014288')
    if (selectedMealStudent.value) {
      selectedMealStudent.value.wallet_balance += mealWallet.topupAmount
    }
    $q.notify({ type: 'positive', message: t('enterprise.meals.successWallet') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 7. Inventario
const assetForm = reactive({ description: 'Notebook Lenovo PNRR 4.0', category: 'Hardware e Dispositivi Mobili', initial_value: 650.0 })
const selectedLoanStudent = ref(allStudents.value[0])
const loanForm = reactive({ asset_id: 'ast-1', student_id: 'student-1' })

watch(selectedLoanStudent, (st) => {
  if (st) loanForm.student_id = st.id
})

async function createAsset() {
  try {
    await inventoryService.createAsset({
      description: assetForm.description,
      category: assetForm.category,
      initial_value: assetForm.initial_value,
      building_location: 'Centrale',
      room_location: 'Laboratorio 1',
      acquisition_date: '2026-01-15'
    })
    $q.notify({ type: 'positive', message: t('enterprise.inventory.successAsset') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

async function createDeviceLoan() {
  try {
    await inventoryService.createLoan({
      asset_id: loanForm.asset_id,
      student_id: loanForm.student_id,
      parent_user_id: 'parent-1',
      expected_return_date: '2027-06-15'
    })
    $q.notify({ type: 'positive', message: t('enterprise.inventory.successLoan') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 8. Privacy
const treatmentForm = reactive({ activity_name: 'Pubblicazione foto didattiche', legal_basis: 'consenso' })

async function saveTreatment() {
  try {
    await privacyService.createTreatment({
      activity_name: treatmentForm.activity_name,
      legal_basis: treatmentForm.legal_basis,
      interested_categories: 'Studenti minori',
      retention_period: 'Ciclo di studi'
    })
    $q.notify({ type: 'positive', message: t('enterprise.privacy.successTreatment') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 9. SIDI
async function syncSIDIStudents() {
  try {
    await sidiService.syncStudentCodes([])
    $q.notify({ type: 'positive', message: t('enterprise.sidi.successSync') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

async function pushScrutinyToMIM() {
  try {
    await sidiService.pushScrutinyResults('5A', '2024/2025', [])
    $q.notify({ type: 'positive', message: t('enterprise.sidi.successPush') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

// 10. Psicologia
const cicMode = ref('single')
const cicBooking = reactive({ slotTime: '2026-10-15 10:00' })
const selectedCicStudent = ref(allStudents.value.find(s => s.id === 'student-minor-1') || allStudents.value[2])
const cicConsent = reactive({ studentId: 'student-minor-1' })
const selectedCicClass = ref('2A')
const studentsInCicClass = computed(() => allStudents.value.filter(s => s.class_name === selectedCicClass.value))

watch(selectedCicStudent, (st) => {
  if (st) cicConsent.studentId = st.id
})

async function bookCICSession() {
  try {
    await psychologyService.bookSession('psy-01', cicBooking.slotTime, 45)
    $q.notify({ type: 'positive', message: t('enterprise.psychology.successBooked') })
  } catch (err) {
    $q.notify({ type: 'warning', message: err.message || t('enterprise.psychology.consentWarning') })
  }
}

async function signCICConsent() {
  try {
    await psychologyService.signConsent(cicConsent.studentId, '2025/2026')
    if (selectedCicStudent.value) {
      selectedCicStudent.value.has_cic_consent = true
    }
    $q.notify({ type: 'positive', message: t('enterprise.psychology.successConsent') })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  }
}

async function signBulkCICConsent() {
  try {
    loading.value = true
    const list = studentsInCicClass.value
    if (list.length === 0) {
      $q.notify({ type: 'warning', message: 'Nessuno studente presente nella classe selezionata.' })
      return
    }
    await psychologyService.signBulkConsent(list.map(s => s.id), '2025/2026')
    list.forEach(st => {
      st.has_cic_consent = true
    })
    $q.notify({
      type: 'positive',
      message: t('enterprise.psychology.successBulkConsent', { count: list.length })
    })
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message })
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const resUsers = await userService.getAll({ role: 'student' })
    const fetched = resUsers.data?.users || resUsers.data
    if (Array.isArray(fetched) && fetched.length > 0) {
      const merged = fetched.map(u => ({
        id: u.id,
        first_name: u.first_name,
        last_name: u.last_name,
        name: `${u.last_name} ${u.first_name} (${u.class_name || 'Generica'})`,
        class_name: u.class_name || 'Generica',
        fiscal_code: u.fiscal_code || u.id,
        email: u.email || '',
        wallet_balance: u.wallet_balance || 25.0,
        has_cic_consent: !!u.has_cic_consent
      }))
      allStudents.value = merged
      filteredStudents.value = [...merged]
      const cls = [...new Set(merged.map(s => s.class_name).filter(Boolean))]
      if (cls.length > 0) classesList.value = cls
    }
  } catch {
    // Uses realistic default students on offline/fallback
  }
})
</script>

<style scoped>
.enterprise-hub-page {
  background: linear-gradient(135deg, #f8fafc 0%, #edf2f7 100%);
  min-height: 100vh;
}
.text-gradient-premium {
  background: linear-gradient(135deg, #1e293b 0%, #3b82f6 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.glass-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  border-radius: 16px;
  border: 1px solid rgba(226, 232, 240, 0.8);
}
.border-primary-soft {
  border-left: 4px solid #3b82f6;
}
.border-bottom-soft {
  border-bottom: 1px solid #e2e8f0;
}
.roles-governance-bar {
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease-in-out;
}
.roles-governance-bar:hover {
  background-color: #f1f5f9;
  border-color: #cbd5e1;
}
</style>
