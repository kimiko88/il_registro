<template>
  <q-layout view="hHh Lpr lFf">
    <SkipLinks />
    <ScreenReaderAnnouncer />

    <q-header class="glass-effect text-slate-900 q-py-xs" :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'" role="banner">

      <q-toolbar role="navigation" :aria-label="t('layout.mainNav')">
        <q-btn
          flat
          dense
          round
          icon="menu"
          :aria-label="t('layout.toggleDrawer')"
          :aria-expanded="leftDrawerOpen"
          color="primary"
          @click="toggleLeftDrawer"
          :key="'drawer-toggle'"
        />

        <q-toolbar-title class="text-weight-bold text-primary ellipsis text-subtitle1 q-pl-xs" style="min-width: 0;">
          <span role="heading" aria-level="1">
            <span class="gt-xs">Registro Elettronico</span>
            <span class="lt-sm">Registro</span>
          </span>
        </q-toolbar-title>

        <q-space />
        
        <!-- School Year Selector for Teachers (Desktop / Tablet) -->
        <div v-if="isTeacherRole" class="gt-sm q-mr-sm row items-center" key="school-year-container">
          <q-select
            v-model="schoolYearStore.selectedSchoolYear"
            :options="schoolYearStore.availableSchoolYears"
            dense
            outlined
            options-dense
            bg-color="white"
            style="min-width: 140px"
            :label="t('layout.schoolYear')"
            key="school-year-select"
            :aria-label="t('layout.schoolYearSelect')"
          >
            <template v-slot:prepend>
              <q-icon name="event" color="primary" size="18px" />
            </template>
          </q-select>
        </div>

        <!-- Theme Selector Menu (Desktop / Tablet) -->
        <q-btn-dropdown
          flat
          round
          dense
          icon="palette"
          color="primary"
          class="gt-sm q-mr-xs"
          key="theme-toggle"
          :aria-label="t('layout.themeAriaLabel')"
        >
          <q-tooltip>{{ t('layout.themeTooltip') }}</q-tooltip>
          <q-list style="min-width: 280px" class="q-py-xs">
            <q-item-label header class="text-weight-bold text-uppercase text-caption letter-spacing-1">
              {{ t('layout.themesTitle') }}
            </q-item-label>

            <q-item
              v-for="themeOption in THEMES"
              :key="themeOption.id"
              clickable
              v-close-popup
              @click="themeStore.setTheme(themeOption.id)"
              :active="themeStore.currentTheme === themeOption.id"
              active-class="bg-indigo-50 text-primary text-weight-bold"
              class="rounded-lg q-mx-xs q-mb-xs"
              role="option"
              :aria-selected="themeStore.currentTheme === themeOption.id"
            >
              <q-item-section avatar>
                <q-avatar size="32px" :color="themeOption.badgeColor" text-color="white">
                  <q-icon :name="themeOption.icon" size="18px" />
                </q-avatar>
              </q-item-section>

              <q-item-section>
                <q-item-label class="text-weight-bold row items-center justify-between">
                  <span>{{ getThemeName(themeOption, t, te) }}</span>
                  <q-chip
                    dense
                    size="xs"
                    color="grey-3"
                    text-color="grey-9"
                    class="text-weight-medium"
                  >
                    {{ getThemeRole(themeOption, t, te) }}
                  </q-chip>
                </q-item-label>
                <q-item-label caption class="text-grey-7">
                  {{ getThemeDescription(themeOption, t, te) }}
                </q-item-label>
              </q-item-section>

              <q-item-section side v-if="themeStore.currentTheme === themeOption.id">
                <q-icon name="check_circle" color="primary" size="20px" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>

        <!-- Focus Mode (ADHD / DSA Clean Reading - Desktop/Tablet) -->
        <FocusModeToggle class="gt-sm q-mr-xs" />

        <!-- Accessibility Quick Toggle Menu (Desktop / Tablet) -->
        <q-btn-dropdown
          flat
          round
          dense
          icon="accessibility_new"
          color="primary"
          class="gt-sm q-mr-xs"
          key="accessibility-toggle"
          :aria-label="t('layout.a11yAriaLabel')"
        >
          <q-tooltip>{{ t('layout.a11yTooltip') }}</q-tooltip>
          <q-list style="min-width: 320px" class="q-py-xs">
            <q-item-label header class="text-weight-bold text-uppercase text-caption letter-spacing-1 text-primary flex items-center justify-between">
              <span>{{ t('layout.a11yTitle') }}</span>
              <q-btn flat round dense icon="help_outline" size="sm" @click="themeStore.toggleKeyboardShortcutsHelp(true)">
                <q-tooltip>Scorciatoie da tastiera (?)</q-tooltip>
              </q-btn>
            </q-item-label>

            <!-- Font OpenDyslexic (DSA) -->
            <q-item clickable class="rounded-lg q-mx-xs q-mb-xs">
              <q-item-section avatar>
                <q-avatar size="32px" color="indigo-50" text-color="indigo-700">
                  <q-icon name="spellcheck" size="18px" />
                </q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('settingsPage.fontDyslexic') }}</q-item-label>
                <q-item-label caption class="text-grey-7">{{ t('layout.dsaFontDesc') }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle
                  v-model="themeStore.dsaFont"
                  color="indigo"
                  dense
                  @update:model-value="themeStore.toggleDsaFont"
                />
              </q-item-section>
            </q-item>

            <!-- Righello di Lettura (Reading Ruler) -->
            <q-item clickable class="rounded-lg q-mx-xs q-mb-xs">
              <q-item-section avatar>
                <q-avatar size="32px" color="amber-50" text-color="amber-9">
                  <q-icon name="straighten" size="18px" />
                </q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Righello di Lettura</q-item-label>
                <q-item-label caption class="text-grey-7">Guida visiva riga per riga</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle
                  v-model="themeStore.readingRuler"
                  color="amber-9"
                  dense
                  @update:model-value="themeStore.toggleReadingRuler"
                />
              </q-item-section>
            </q-item>

            <!-- Text-to-Speech (TTS) -->
            <q-item clickable class="rounded-lg q-mx-xs q-mb-xs">
              <q-item-section avatar>
                <q-avatar size="32px" color="teal-50" text-color="teal-8">
                  <q-icon name="volume_up" size="18px" />
                </q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">Sintesi Vocale (TTS)</q-item-label>
                <q-item-label caption class="text-grey-7">Lettura vocale compiti e avvisi</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle
                  v-model="themeStore.ttsEnabled"
                  color="teal"
                  dense
                  @update:model-value="themeStore.toggleTts"
                />
              </q-item-section>
            </q-item>

            <!-- Contrasto Elevato -->
            <q-item clickable class="rounded-lg q-mx-xs q-mb-xs">
              <q-item-section avatar>
                <q-avatar size="32px" color="grey-2" text-color="dark">
                  <q-icon name="contrast" size="18px" />
                </q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ t('layout.highContrast') }}</q-item-label>
                <q-item-label caption class="text-grey-7">{{ t('layout.highContrastDesc') }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle
                  v-model="themeStore.highContrast"
                  color="dark"
                  dense
                  @update:model-value="themeStore.toggleHighContrast"
                />
              </q-item-section>
            </q-item>

            <q-separator class="q-my-xs" />

            <!-- Link Scorciatoie & Dichiarazione AgID -->
            <q-item clickable class="rounded-lg q-mx-xs text-primary" @click="themeStore.toggleKeyboardShortcutsHelp(true)">
              <q-item-section avatar>
                <q-icon name="keyboard" size="20px" />
              </q-item-section>
              <q-item-section class="text-weight-bold">{{ t('a11y.shortcutsTitle') || 'Scorciatoie da Tastiera (?)' }}</q-item-section>
            </q-item>

            <q-item clickable to="/accessibility-statement" class="rounded-lg q-mx-xs text-grey-8">
              <q-item-section avatar>
                <q-icon name="verified_user" size="20px" color="positive" />
              </q-item-section>
              <q-item-section class="text-weight-medium text-caption">{{ t('layout.footerA11yStatement') || 'Dichiarazione AgID / WCAG 2.2' }}</q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>

        <!-- Language Selector Menu (Desktop / Tablet) -->
        <q-btn-dropdown
          flat
          round
          dense
          icon="language"
          color="primary"
          class="gt-sm q-mr-xs"
          key="language-toggle"
          aria-label="Seleziona Lingua"
        >
          <q-tooltip>{{ t('common.language') }}</q-tooltip>
          <q-list style="min-width: 220px" class="q-py-xs">
            <q-item-label header class="text-weight-bold text-uppercase text-caption letter-spacing-1">
              {{ t('common.language') }}
            </q-item-label>

            <q-item
              v-for="loc in SUPPORTED_LOCALES"
              :key="loc.value"
              clickable
              v-close-popup
              @click="changeAppLanguage(loc.value)"
              :active="currentLocaleValue === loc.value"
              active-class="bg-indigo-50 text-primary text-weight-bold"
              class="rounded-lg q-mx-xs q-mb-xs"
              role="option"
              :aria-selected="currentLocaleValue === loc.value"
            >
              <q-item-section avatar min-width="32px">
                <span style="font-size: 1.2rem;">{{ loc.flag }}</span>
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold">{{ loc.label }}</q-item-label>
                <q-item-label caption class="text-grey-7">{{ loc.code }}</q-item-label>
              </q-item-section>
              <q-item-section side v-if="currentLocaleValue === loc.value">
                <q-icon name="check_circle" color="primary" size="20px" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>

        <!-- Dark Mode Toggle (Desktop / Tablet) -->
        <q-btn flat round dense :icon="$q.dark.isActive ? 'light_mode' : 'dark_mode'" @click="$q.dark.toggle()" color="primary" class="gt-sm q-mr-xs" :key="'dark-toggle'" :aria-label="$q.dark.isActive ? t('layout.lightMode') : t('layout.darkMode')">
           <q-tooltip>{{ $q.dark.isActive ? t('layout.lightMode') : t('layout.darkMode') }}</q-tooltip>
        </q-btn>

        <!-- Fullscreen Toggle (Desktop/Tablet) -->
        <q-btn 
          v-if="$q.fullscreen"
          flat 
          round 
          dense 
          class="gt-sm q-mr-xs"
          :icon="$q.fullscreen.isActive ? 'fullscreen_exit' : 'fullscreen'" 
          @click="$q.fullscreen.toggle()" 
          color="primary" 
          :key="'fullscreen-toggle'"
          :aria-label="$q.fullscreen.isActive ? t('layout.exitFullscreen') : t('layout.enterFullscreen')"
        >
           <q-tooltip>{{ $q.fullscreen.isActive ? t('layout.exitFullscreen') : t('layout.enterFullscreen') }}</q-tooltip>
        </q-btn>

        <!-- Global Search Ctrl+K -->
        <q-btn
          flat
          dense
          color="primary"
          icon="search"
          class="q-mr-xs search-shortcut-btn"
          :aria-label="t('common.search')"
          @click="globalSearchRef?.open()"
          key="global-search-btn"
        >
          <q-tooltip>{{ t('common.search') }}</q-tooltip>
          <q-badge floating transparent class="search-kbd-badge gt-xs">K</q-badge>
        </q-btn>

        <!-- Help Center (Desktop / Tablet) -->
        <q-btn
          flat
          round
          dense
          icon="help_outline"
          color="primary"
          class="gt-xs q-mr-xs"
          :aria-label="t('help.openHelp')"
          @click="helpCenterRef?.open()"
          key="help-center-btn"
        >
          <q-tooltip>{{ t('help.openHelp') }}</q-tooltip>
        </q-btn>

        <!-- Notifications -->
        <q-btn flat round dense icon="notifications" color="primary" class="q-mr-xs" :aria-label="t('notifications.title')" @click="navigateToNotifications">
          <q-tooltip>{{ t('notifications.title') }}</q-tooltip>
        </q-btn>

        <q-btn flat round dense icon="account_circle" color="primary" class="q-mr-xs" :aria-label="t('nav.profile')" @click="navigateToProfile" />

        <!-- Quick Settings & Customization Right Drawer Toggle (Mobile only) -->
        <q-btn
          flat
          round
          dense
          icon="tune"
          color="primary"
          class="lt-md"
          :aria-label="t('layout.quickSettings') || 'Opzioni e Accessibilità'"
          @click="toggleRightDrawer"
          key="quick-settings-mobile-btn"
        >
          <q-tooltip>{{ t('layout.quickSettings') || 'Opzioni e Accessibilità' }}</q-tooltip>
        </q-btn>

      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
      :width="260"
      role="navigation"
      :aria-label="t('layout.sideNav')"
    >
      <div class="column full-height no-wrap">
        <!-- User Profile Section -->
        <div class="q-pa-md bg-primary text-white relative-position overflow-hidden" v-if="userName" role="region" :aria-label="t('layout.userProfile')">
          <div class="row items-center relative-position" style="z-index: 1">
            <q-avatar size="42px" color="white" text-color="primary" class="q-mr-md shadow-soft" aria-hidden="true">
              <q-icon name="person" size="24px" />
            </q-avatar>
            <div class="col">
              <div class="text-subtitle1 text-weight-bold no-wrap ellipsis" :aria-label="t('layout.connectedUser', { name: userName })">{{ userName }}</div>
              <div class="text-caption opacity-80 text-uppercase letter-spacing-1" :aria-label="t('layout.userRole', { role: roleLabel })">{{ roleLabel }}</div>
            </div>
          </div>
          <div class="absolute-bottom-right q-mr-n-lg q-mb-n-lg" style="width: 90px; height: 90px; border-radius: 50%; background: rgba(255,255,255,0.1)" aria-hidden="true"></div>
        </div>

        <!-- School Year Selector for Teachers on Mobile Drawer -->
        <div v-if="isTeacherRole" class="lt-md q-pa-sm bg-slate-50 border-b border-slate-100">
          <q-select
            v-model="schoolYearStore.selectedSchoolYear"
            :options="schoolYearStore.availableSchoolYears"
            dense
            outlined
            options-dense
            bg-color="white"
            :label="t('layout.schoolYear')"
            :aria-label="t('layout.schoolYearSelect')"
          >
            <template v-slot:prepend>
              <q-icon name="event" color="primary" size="18px" />
            </template>
          </q-select>
        </div>

        <!-- Menu Items -->
        <q-scroll-area class="col">
          <div class="q-pa-sm">
            <div class="text-overline text-grey-5 q-px-sm q-mb-xs letter-spacing-2" aria-hidden="true">{{ t('common.mainMenu') }}</div>
            <q-list dense padding class="q-gutter-y-xs" aria-label="Navigazione principale" role="menu">
              <template v-for="(item, idx) in menuItems" :key="item.path || item.category || idx">
                <!-- Group Category with children -->
                <q-expansion-item
                  v-if="item.children"
                  group="menu-group"
                  :icon="item.icon"
                  :label="translateCategory(item.category)"
                  dense
                  header-class="text-weight-bold text-slate-700 rounded-lg"
                  :default-opened="isCategoryActive(item)"
                  role="menuitem"
                >
                  <q-list dense class="q-pl-sm q-gutter-y-xs">
                    <q-item
                      v-for="child in item.children"
                      :key="child.path"
                      clickable
                      :to="child.path"
                      :exact="child.exact !== undefined ? child.exact : false"
                      active-class="active-menu-item"
                      class="rounded-lg transition-all"
                      :aria-label="translateMenuLabel(child.label)"
                      role="menuitem"
                    >
                      <q-item-section avatar min-width="32px">
                        <q-icon :name="child.icon" size="18px" aria-hidden="true" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">{{ translateMenuLabel(child.label) }}</q-item-label>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-expansion-item>

                <!-- Single Item -->
                <q-item
                  v-else
                  clickable
                  :to="item.path"
                  :exact="item.exact !== undefined ? item.exact : false"
                  active-class="active-menu-item"
                  class="rounded-lg transition-all"
                  :aria-label="translateMenuLabel(item.label)"
                  role="menuitem"
                >
                  <q-item-section avatar min-width="32px">
                    <q-icon :name="item.icon" size="20px" aria-hidden="true" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">{{ translateMenuLabel(item.label) }}</q-item-label>
                  </q-item-section>
                </q-item>
              </template>
            </q-list>
          </div>
        </q-scroll-area>

        <!-- Help & Guide Button -->
        <div class="q-px-md q-pb-xs">
          <q-item
            clickable
            class="rounded-lg q-pa-sm text-primary"
            @click="helpDrawerRef?.open()"
            role="button"
            :aria-label="t('help.openHelp')"
            style="background: rgba(99, 102, 241, 0.07); border: 1.5px solid rgba(99, 102, 241, 0.18);"
          >
            <q-item-section avatar>
              <q-icon name="help_outline" size="20px" aria-hidden="true" />
            </q-item-section>
            <q-item-section class="text-weight-bold">
              {{ t('help.openHelp') }}
            </q-item-section>
            <q-item-section side>
              <q-icon name="chevron_right" size="16px" color="primary" />
            </q-item-section>
          </q-item>
        </div>

        <!-- Logout Button at Bottom -->
        <div class="q-pa-md border-t border-slate-100">
          <q-item
            clickable
            class="rounded-lg q-pa-md text-grey-8"
            @click="handleLogout"
            :disable="loggingOut"
            role="button"
            :aria-label="t('common.logout')"
            :aria-busy="loggingOut"
          >
            <q-item-section avatar>
              <q-icon name="logout" size="20px" aria-hidden="true" />
            </q-item-section>
            <q-item-section class="text-weight-bold">
              {{ t('common.logout') }}
            </q-item-section>
            <q-item-section side v-if="loggingOut">
              <q-spinner size="20px" />
            </q-item-section>
          </q-item>
        </div>
      </div>
    </q-drawer>

    <!-- Right Drawer: Quick Settings, Themes & Accessibility for Mobile -->
    <q-drawer
      side="right"
      v-model="rightDrawerOpen"
      overlay
      bordered
      :width="320"
      :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
      role="region"
      :aria-label="t('layout.quickSettings') || 'Opzioni e Accessibilità'"
    >
      <div class="column full-height no-wrap">
        <!-- Header -->
        <div class="q-pa-md bg-primary text-white row items-center justify-between shadow-soft">
          <div class="row items-center gap-sm">
            <q-icon name="tune" size="22px" />
            <div class="text-subtitle1 text-weight-bold">{{ t('layout.quickSettings') || 'Opzioni e Accessibilità' }}</div>
          </div>
          <q-btn flat round dense icon="close" color="white" @click="rightDrawerOpen = false" :aria-label="t('common.close') || 'Chiudi'" />
        </div>

        <q-scroll-area class="col q-pa-md">
          <!-- 1. Modalità Scura / Chiara -->
          <div class="q-mb-md">
            <div class="text-caption text-weight-bold text-uppercase text-grey-7 letter-spacing-1 q-mb-xs">
              {{ t('layout.themeAriaLabel') || 'Aspetto' }}
            </div>
            <q-card flat bordered class="rounded-xl q-pa-sm" :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-slate-50'">
              <div class="row items-center justify-between">
                <div class="row items-center gap-sm">
                  <q-avatar size="36px" :color="$q.dark.isActive ? 'indigo-9' : 'amber-1'" :text-color="$q.dark.isActive ? 'amber-3' : 'amber-9'">
                    <q-icon :name="$q.dark.isActive ? 'dark_mode' : 'light_mode'" size="20px" />
                  </q-avatar>
                  <div>
                    <div class="text-weight-bold text-body2">{{ $q.dark.isActive ? t('layout.darkMode') : t('layout.lightMode') }}</div>
                    <div class="text-caption text-grey-7">{{ $q.dark.isActive ? (t('layout.darkModeDesc') || 'OLED Dark') : (t('layout.lightModeDesc') || 'Standard Light') }}</div>
                  </div>
                </div>
                <q-toggle
                  v-model="$q.dark.isActive"
                  @update:model-value="$q.dark.set"
                  color="primary"
                  dense
                />
              </div>
            </q-card>
          </div>

          <!-- 2. Tema Visivo e Palette -->
          <div class="q-mb-md">
            <div class="text-caption text-weight-bold text-uppercase text-grey-7 letter-spacing-1 q-mb-xs">
              {{ t('layout.themesTitle') || 'Temi & Palette' }}
            </div>
            <q-list class="rounded-xl border border-slate-200 overflow-hidden" :class="$q.dark.isActive ? 'bg-grey-9 border-grey-8' : 'bg-white'">
              <q-item
                v-for="themeOption in THEMES"
                :key="themeOption.id"
                clickable
                @click="themeStore.setTheme(themeOption.id)"
                :active="themeStore.currentTheme === themeOption.id"
                active-class="bg-indigo-50 text-primary text-weight-bold"
                dense
                class="q-py-sm"
              >
                <q-item-section avatar min-width="36px">
                  <q-avatar size="28px" :color="themeOption.badgeColor" text-color="white">
                    <q-icon :name="themeOption.icon" size="16px" />
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-body2 text-weight-medium">{{ getThemeName(themeOption, t, te) }}</q-item-label>
                  <q-item-label caption class="text-grey-7 text-xs">{{ getThemeRole(themeOption, t, te) }} • {{ getThemeDescription(themeOption, t, te) }}</q-item-label>
                </q-item-section>
                <q-item-section side v-if="themeStore.currentTheme === themeOption.id">
                  <q-icon name="check_circle" color="primary" size="18px" />
                </q-item-section>
              </q-item>
            </q-list>
          </div>

          <!-- 3. Lingua -->
          <div class="q-mb-md">
            <div class="text-caption text-weight-bold text-uppercase text-grey-7 letter-spacing-1 q-mb-xs">
              {{ t('common.language') || 'Lingua' }}
            </div>
            <q-select
              v-model="currentLocaleValue"
              :options="SUPPORTED_LOCALES"
              option-value="value"
              option-label="label"
              emit-value
              map-options
              outlined
              dense
              rounded
              :bg-color="$q.dark.isActive ? 'grey-9' : 'white'"
              @update:model-value="changeAppLanguage"
            >
              <template v-slot:prepend>
                <q-icon name="language" color="primary" size="18px" />
              </template>
              <template v-slot:option="scope">
                <q-item v-bind="scope.itemProps" dense>
                  <q-item-section avatar min-width="28px">
                    <span style="font-size: 1.1rem;">{{ scope.opt.flag }}</span>
                  </q-item-section>
                  <q-item-section>
                    <q-item-label class="text-body2">{{ scope.opt.label }}</q-item-label>
                  </q-item-section>
                  <q-item-section side v-if="currentLocaleValue === scope.opt.value">
                    <q-icon name="check" color="primary" size="16px" />
                  </q-item-section>
                </q-item>
              </template>
            </q-select>
          </div>

          <!-- 4. Accessibilità & Inclusione (A11y) -->
          <div class="q-mb-md">
            <div class="text-caption text-weight-bold text-uppercase text-grey-7 letter-spacing-1 q-mb-xs flex items-center justify-between">
              <span>{{ t('layout.a11yTitle') || 'Accessibilità' }}</span>
              <q-btn flat round dense icon="help_outline" size="xs" @click="themeStore.toggleKeyboardShortcutsHelp(true)">
                <q-tooltip>{{ t('a11y.shortcutsTitle') || 'Scorciatoie da tastiera (?)' }}</q-tooltip>
              </q-btn>
            </div>

            <q-list class="rounded-xl border border-slate-200 overflow-hidden" :class="$q.dark.isActive ? 'bg-grey-9 border-grey-8' : 'bg-white'">
              <!-- Font OpenDyslexic -->
              <q-item dense class="q-py-xs">
                <q-item-section avatar min-width="36px">
                  <q-avatar size="28px" color="indigo-50" text-color="indigo-700">
                    <q-icon name="spellcheck" size="16px" />
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-body2 text-weight-medium">{{ t('settingsPage.fontDyslexic') || 'OpenDyslexic (DSA)' }}</q-item-label>
                  <q-item-label caption class="text-grey-7 text-xs">{{ t('layout.dsaFontDesc') || 'Font ad alta leggibilità' }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle
                    v-model="themeStore.dsaFont"
                    color="indigo"
                    dense
                    @update:model-value="themeStore.toggleDsaFont"
                  />
                </q-item-section>
              </q-item>

              <!-- Righello di Lettura -->
              <q-item dense class="q-py-xs">
                <q-item-section avatar min-width="36px">
                  <q-avatar size="28px" color="amber-50" text-color="amber-9">
                    <q-icon name="straighten" size="16px" />
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-body2 text-weight-medium">{{ t('layout.readingRulerTitle') || 'Righello di Lettura' }}</q-item-label>
                  <q-item-label caption class="text-grey-7 text-xs">{{ t('layout.readingRulerDesc') || 'Guida visiva riga per riga' }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle
                    v-model="themeStore.readingRuler"
                    color="amber-9"
                    dense
                    @update:model-value="themeStore.toggleReadingRuler"
                  />
                </q-item-section>
              </q-item>

              <!-- TTS -->
              <q-item dense class="q-py-xs">
                <q-item-section avatar min-width="36px">
                  <q-avatar size="28px" color="teal-50" text-color="teal-8">
                    <q-icon name="volume_up" size="16px" />
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-body2 text-weight-medium">{{ t('layout.ttsTitle') || 'Sintesi Vocale (TTS)' }}</q-item-label>
                  <q-item-label caption class="text-grey-7 text-xs">{{ t('layout.ttsDesc') || 'Lettura vocale compiti e avvisi' }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle
                    v-model="themeStore.ttsEnabled"
                    color="teal"
                    dense
                    @update:model-value="themeStore.toggleTts"
                  />
                </q-item-section>
              </q-item>

              <!-- Contrasto Elevato -->
              <q-item dense class="q-py-xs">
                <q-item-section avatar min-width="36px">
                  <q-avatar size="28px" color="grey-2" text-color="dark">
                    <q-icon name="contrast" size="16px" />
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-body2 text-weight-medium">{{ t('layout.highContrast') || 'Contrasto Elevato' }}</q-item-label>
                  <q-item-label caption class="text-grey-7 text-xs">{{ t('layout.highContrastDesc') || 'Contrasto massimo 21:1' }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle
                    v-model="themeStore.highContrast"
                    color="dark"
                    dense
                    @update:model-value="themeStore.toggleHighContrast"
                  />
                </q-item-section>
              </q-item>

              <q-separator />

              <!-- Scorciatoie & Dichiarazione AgID -->
              <q-item clickable dense class="q-py-xs text-primary" @click="themeStore.toggleKeyboardShortcutsHelp(true); rightDrawerOpen = false">
                <q-item-section avatar min-width="36px">
                  <q-icon name="keyboard" size="18px" />
                </q-item-section>
                <q-item-section class="text-body2 text-weight-medium">{{ t('a11y.shortcutsTitle') || 'Scorciatoie da Tastiera (?)' }}</q-item-section>
              </q-item>

              <q-item clickable dense to="/accessibility-statement" class="q-py-xs text-grey-8" @click="rightDrawerOpen = false">
                <q-item-section avatar min-width="36px">
                  <q-icon name="verified_user" size="18px" color="positive" />
                </q-item-section>
                <q-item-section class="text-body2 text-weight-medium">{{ t('layout.footerA11yStatement') || 'Dichiarazione AgID / WCAG 2.2' }}</q-item-section>
              </q-item>
            </q-list>
          </div>

          <!-- 5. Centro Assistenza & Guida -->
          <div class="q-mb-md">
            <q-btn
              outline
              color="primary"
              icon="help_outline"
              :label="t('help.openHelp') || 'Centro Assistenza'"
              class="full-width rounded-xl"
              no-caps
              @click="helpCenterRef?.open(); rightDrawerOpen = false"
            />
          </div>
        </q-scroll-area>
      </div>
    </q-drawer>

    <!-- Global Search Modal (Ctrl+K) -->
    <GlobalSearch ref="globalSearchRef" />

    <q-page-container role="main" id="main-content" tabindex="-1">
      <!-- Dynamic Breadcrumb Navigation Header -->
      <div v-if="breadcrumbs.length > 0" class="q-px-md q-pt-md">
        <q-breadcrumbs :aria-label="t('layout.breadcrumbNav')" class="text-caption text-grey-7" active-color="primary" separator-icon="chevron_right" separator-color="grey-5">
          <q-breadcrumbs-el icon="home" to="/dashboard" :label="t('nav.dashboard')" />
          <q-breadcrumbs-el
            v-for="(crumb, idx) in breadcrumbs"
            :key="idx"
            :label="crumb.label"
            :to="crumb.path"
            :icon="crumb.icon && crumb.icon !== 'chevron_right' ? crumb.icon : undefined"
          />
        </q-breadcrumbs>
      </div>
      <router-view v-slot="{ Component }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" :key="$route.path" />
        </transition>
      </router-view>
    </q-page-container>

    <!-- Onboarding Tour (global, triggers on first login per role) -->
    <OnboardingTour ref="tourRef" @open-guide="helpCenterRef?.open()" />

    <!-- Help Center Panel (full-screen, opens from toolbar or FAB) -->
    <HelpCenterPanel ref="helpCenterRef" @restart-tour="handleRestartTour" />

    <!-- Help Drawer (slides in from right, for FAB quick access) -->
    <HelpDrawer ref="helpDrawerRef" @restart-tour="handleRestartTour" />

    <!-- Session Reauth Dialog (globale: appare sopra la pagina quando il token scade) -->
    <SessionReauthDialog />

    <!-- Reading Ruler Overlay for Dyslexia/ADHD -->
    <ReadingRuler />

    <!-- Keyboard Shortcuts Help Dialog -->
    <KeyboardShortcutsDialog />

    <!-- Accessible Footer with AgID Statement & Shortcuts Link -->
    <q-footer class="bg-slate-900 text-white text-caption q-py-xs q-px-md print-hide" role="contentinfo">
      <div class="row items-center justify-between">
        <div>{{ t('layout.footerCopyright') || '© 2026 Registro Elettronico Scolastico • Conforme AgID & WCAG 2.2 AA' }}</div>
        <div class="flex items-center gap-md">
          <router-link to="/accessibility-statement" class="text-amber-4 text-weight-medium text-decoration-none flex items-center">
            <q-icon name="accessibility" size="xs" class="q-mr-xs" />
            {{ t('layout.footerA11yStatement') || 'Dichiarazione di Accessibilità (AgID)' }}
          </router-link>
          <a href="#" @click.prevent="themeStore.toggleKeyboardShortcutsHelp(true)" class="text-grey-4 text-decoration-none flex items-center">
            <q-icon name="keyboard" size="xs" class="q-mr-xs" />
            {{ t('layout.footerShortcuts') || 'Scorciatoie ( ? )' }}
          </a>
        </div>
      </div>
    </q-footer>

  </q-layout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTeacherStore } from '@/stores/teacher'
import { useClassesStore } from '@/stores/classes'
import { useThemeStore, THEMES, getThemeName, getThemeDescription, getThemeRole } from '@/stores/theme'
import { useSchoolYearStore } from '@/stores/schoolYear'
import { useAuth } from '@/composables/useAuth'
import { useMenuItems } from '@/composables/useMenuItems'
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import { i18n } from '@/i18n'
import { SUPPORTED_LOCALES, applyLocale, normalizeLocale } from '@/utils/locale'
import GlobalSearch from '@/components/Common/GlobalSearch.vue'
import OnboardingTour from '@/components/Common/OnboardingTour.vue'
import HelpDrawer from '@/components/Common/HelpDrawer.vue'
import HelpCenterPanel from '@/components/Common/HelpCenterPanel.vue'
import SessionReauthDialog from '@/components/Common/SessionReauthDialog.vue'
import ReadingRuler from '@/components/Common/ReadingRuler.vue'
import KeyboardShortcutsDialog from '@/components/Common/KeyboardShortcutsDialog.vue'
import SkipLinks from '@/components/Common/SkipLinks.vue'
import ScreenReaderAnnouncer from '@/components/Common/ScreenReaderAnnouncer.vue'
import FocusModeToggle from '@/components/Common/FocusModeToggle.vue'
import { useSessionReauth } from '@/composables/useSessionReauth'
import { useGlobalKeyboardShortcuts } from '@/composables/useGlobalKeyboardShortcuts'
import { setReauthHandler } from '@/services/api'


useGlobalKeyboardShortcuts()

const globalSearchRef = ref(null)
const tourRef = ref(null)
const helpDrawerRef = ref(null)
const helpCenterRef = ref(null)

function handleRestartTour() {
  // Clear the flag so the tour shows again, then start it
  const role = authStore.userRole || authStore.user?.role || 'user'
  localStorage.removeItem(`onboarding_done_${role.toLowerCase()}`)
  tourRef.value?.startTour()
}

const route = useRoute()
const router = useRouter()
const $q = useQuasar()
const { t, te, locale } = useI18n()
const themeStore = useThemeStore()
const schoolYearStore = useSchoolYearStore()

const currentLocaleValue = computed(() => normalizeLocale(locale.value))

function changeAppLanguage(langCode) {
  applyLocale(langCode, i18n, $q)
  $q.notify({
    type: 'positive',
    icon: 'language',
    message: t('notifications.languageChanged') || 'Lingua aggiornata con successo'
  })
}

const menuLabelToKeyMap = {
  'Dashboard': 'dashboard',
  'Gestione Scuole': 'schools',
  'La Mia Scuola': 'mySchool',
  'Gestione Utenti': 'users',
  'Gestione Admin': 'admins',
  'Monitoraggio Sistema': 'monitoring',
  'Analytics Globali': 'analytics',
  'Analytics': 'analytics',
  'Audit Logs': 'auditLogs',
  'Impostazioni': 'settings',
  'Supporto': 'support',
  'Supporto & Assistenza': 'support',
  'Feature Flags & Istituto': 'featureFlags',
  'Google & Teams E-Learning': 'elearning',
  'Studenti': 'students',
  'Utenti': 'users',
  'Classi': 'classes',
  'Gruppi Linguistici / Articolati': 'groups',
  'Documenti': 'documents',
  'Certificati': 'certificates',
  'Libri di Testo': 'textbooks',
  'Riunioni': 'meetings',
  'Comunicazioni': 'communications',
  'Report': 'reports',
  'PCTO': 'pcto',
  'Scrutinio': 'scrutiny',
  'Le Mie Classi': 'myClasses',
  'Registro Classe': 'classRegister',
  'Programmazione UdA': 'uda',
  'Valutazione Competenze': 'competencies',
  'Voti': 'grades',
  'Presenze': 'attendance',
  'Didattica': 'didactics',
  'Credito Scolastico': 'credits',
  'Corsi Recupero & PAI': 'recovery',
  'Registro Sostegno & PEI': 'supportRegister',
  'Ricevimento Generale': 'generalMeetings',
  'Flussi SIDI': 'sidi',
  'Piani PDP / PEI': 'pdp',
  'Piano PDP / PEI': 'pdp',
  'Rubriche Valutative': 'rubrics',
  'Coordinamento': 'coordination',
  'Orario Lezioni': 'timetable',
  'Orario Scolastico': 'timetable',
  'Agenda': 'agenda',
  'Colloqui': 'colloqui',
  'Sostituzioni': 'substitutions',
  'Gestione Sostituzioni': 'substitutions',
  'Verbali': 'verbali',
  'Note Disciplinari': 'notes',
  'I Miei Voti': 'myGrades',
  'Le Mie Presenze': 'myAttendance',
  'Compiti': 'homework',
  'Materiale Didattico': 'didactics',
  'Orientamento': 'orientamento',
  'Calendario Scolastico': 'calendar',
  'Pagella': 'reportCard',
  'Profilo': 'profile',
  'I Miei Figli': 'myChildren',
  'Obiettivi': 'goals',
  'Uscite & Viaggi': 'trips',
  'Pagamenti': 'payments',
  'Assemblee & Riunioni': 'assemblies',
  'Fascicolo Documentale & Atti': 'documents',
  'Fascicolo Documentale': 'documents',
  'Compiti a casa': 'homework',
  'Media Voti': 'averageGrade'
}

const categoryToKeyMap = {
  'Anagrafiche & Classi': 'anagraficheClassi',
  'Atti & Certificati': 'attiCertificati',
  'Servizi & Report': 'serviziReport',
  'Didattica & Valutazione': 'didatticaValutazione',
  'Organizzazione & Orario': 'organizzazioneOrario',
  'Comunicazioni & Atti': 'comunicazioniAtti',
  'Percorsi & Comunicazioni': 'percorsiComunicazioni',
  'Valutazione & Didattica': 'valutazioneDidattica',
  'Servizi & Orari': 'serviziOrari',
  'Comunicazioni & Account': 'comunicazioniAccount'
}

function translateMenuLabel(label) {
  const key = menuLabelToKeyMap[label]
  if (key && te('nav.' + key)) {
    return t('nav.' + key)
  }
  return label
}

function translateCategory(cat) {
  const key = categoryToKeyMap[cat]
  if (key && te('categories.' + key)) {
    return t('categories.' + key)
  }
  return cat
}

const navigateToNotifications = () => {
  const role = (userRole.value || '').toLowerCase()
  if (role === 'teacher' || role === 'docente' || role === 'coordinator') {
    router.push('/teacher/communications')
  } else if (role === 'student') {
    router.push('/student/communications')
  } else if (role === 'parent') {
    router.push('/parent/communications')
  } else if (role === 'admin' || role === 'superadmin' || role === 'system_auditor') {
    router.push('/admin/dashboard')
  } else {
    router.push('/secretary/communications')
  }
}

const navigateToProfile = () => {
  const role = (userRole.value || '').toLowerCase()
  if (role === 'student') {
    router.push('/student/profile')
  } else if (role === 'parent') {
    router.push('/parent/profile')
  } else if (role === 'teacher' || role === 'docente' || role === 'coordinator') {
    router.push('/teacher/settings')
  } else if (role === 'admin' || role === 'superadmin' || role === 'system_auditor') {
    router.push('/admin/settings')
  } else if (role === 'secretary' || role === 'principal' || role === 'vice_principal' || role === 'staff') {
    router.push('/secretary/settings')
  } else {
    router.push('/dashboard')
  }
}

onMounted(() => {
  themeStore.initTheme()

  // Registra il handler per la re-autenticazione in-page.
  // Quando il token scade e il refresh fallisce, api.js chiamerà triggerReauth()
  // invece di navigare a /login, preservando lo stato della pagina corrente.
  const { triggerReauth } = useSessionReauth()
  setReauthHandler(triggerReauth)
})

watch(() => route.path, () => {
  if ($q.screen.lt.md) {
    leftDrawerOpen.value = false
    rightDrawerOpen.value = false
  }
})

// Dynamic Breadcrumbs
const breadcrumbs = computed(() => {
  if (!route.path || route.path === '/dashboard' || route.path === '/') return []
  const items = []
  
  const currentMatch = route.meta?.title || route.name || 'Pagina'
  const current = { label: currentMatch, icon: undefined, path: route.path }

  if (route.path.startsWith('/teacher/') && route.path !== '/teacher') {
    items.push({ label: t('roles.teacher') || 'Docente', icon: 'school', path: '/teacher' })
  } else if (route.path.startsWith('/student/') && route.path !== '/student') {
    items.push({ label: t('roles.student') || 'Studente', icon: 'person', path: '/student' })
  } else if (route.path.startsWith('/parent/') && route.path !== '/parent') {
    items.push({ label: t('roles.parent') || 'Genitore', icon: 'family_restroom', path: '/parent' })
  } else if (route.path.startsWith('/admin/') && route.path !== '/admin' && route.path !== '/admin/dashboard') {
    items.push({ label: t('roles.admin') || 'Amministrazione', icon: 'admin_panel_settings', path: '/admin/dashboard' })
  } else if (route.path.startsWith('/secretary/') && route.path !== '/secretary') {
    items.push({ label: t('roles.secretary') || 'Segreteria', icon: 'badge', path: '/secretary' })
  }

  items.push(current)
  return items
})

const authStore = useAuthStore()
const { user, userName, userRole } = storeToRefs(authStore)
const { logout } = useAuth()

const leftDrawerOpen = ref(false)
const rightDrawerOpen = ref(false)
const loggingOut = ref(false)

// Get role label for display
const roleLabel = computed(() => {
  if (!userRole.value) return t('roles.user')
  const roleKey = userRole.value.toLowerCase()
  if (te('roles.' + roleKey)) {
    return t('roles.' + roleKey)
  }
  return userRole.value
})

const isTeacherRole = computed(() => {
  const r = (userRole.value || '').toLowerCase()
  return r === 'teacher' || r === 'docente' || r === 'coordinator'
})

const teacherStore = useTeacherStore()
const classesStore = useClassesStore()

const isTeacherCoordinator = computed(() => {
  const r = (userRole.value || '').toLowerCase()
  if (r === 'coordinator') return true
  if (r !== 'teacher' && r !== 'docente') return false
  if (teacherStore.isCoordinator) return true
  const currentUserId = user.value?.id
  if (currentUserId && classesStore.classes.some(c => c.coordinator_id === currentUserId)) {
    return true
  }
  return false
})

// Get menu items based on role
const menuItems = ref([])
watch([userRole, isTeacherCoordinator, () => classesStore.classes], ([newRole, isCoord]) => {
  if (!newRole) {
    menuItems.value = []
    return
  }
  let items = useMenuItems(newRole)
  if (newRole === 'teacher' && !isCoord) {
    items = items.map(cat => {
      if (!cat.children) return cat
      return {
        ...cat,
        children: cat.children.filter(child => !child.coordinatorOnly)
      }
    }).filter(cat => !cat.children || cat.children.length > 0)
  }
  menuItems.value = items
}, { immediate: true })

watch(() => authStore.user, (user) => {
  if (user) {
    schoolYearStore.initializeForUser(user)
  }
}, { immediate: true })

watch([userRole, () => schoolYearStore.selectedSchoolYear], ([newRole, newSY]) => {
  if (newRole === 'teacher') {
    classesStore.fetchAssignedClasses(newSY)
  }
}, { immediate: true })

const isCategoryActive = (category) => {
  if (!category || !category.children) return false
  return category.children.some(child => route.path === child.path || (child.path !== '/' && route.path.startsWith(child.path)))
}

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

function toggleRightDrawer() {
  rightDrawerOpen.value = !rightDrawerOpen.value
}

async function handleLogout() {
  loggingOut.value = true
  try {
    await logout()
    $q.notify({
      type: 'positive',
      message: t('notifications.logoutSuccess'),
      position: 'top',
      timeout: 2000
    })
  } catch (error) {
    $q.notify({
      type: 'negative',
      message: t('notifications.logoutError'),
      position: 'top',
      timeout: 3000
    })
  } finally {
    loggingOut.value = false
  }
}

onMounted(() => {
  themeStore.loadFromCloud()
})
</script>


<style scoped>
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
.sr-only:focus {
  position: absolute;
  width: auto;
  height: auto;
  clip: auto;
  white-space: normal;
}
.opacity-80 {
  opacity: 0.8;
}

.letter-spacing-1 {
    letter-spacing: 1px;
}

.letter-spacing-2 {
    letter-spacing: 2px;
}

.transition-all {
    transition: all 0.3s ease;
}

.search-shortcut-btn {
  position: relative;
}

.search-kbd-badge {
  font-size: 9px;
  font-weight: 700;
  background: rgba(99, 102, 241, 0.15) !important;
  color: #6366f1 !important;
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 4px;
  padding: 0 3px;
  line-height: 14px;
  height: 14px;
  top: 2px;
  right: 2px;
}
</style>
