# Guida allo sviluppo Frontend

Questa guida è rivolta agli sviluppatori che contribuiscono al frontend di il_registro. Copre le convenzioni di codice, i pattern architetturali, la gestione dello stato e le best practice per Vue 3 + Quasar.

---

## Indice

- [Principi generali](#principi-generali)
- [Convenzioni di naming](#convenzioni-di-naming)
- [Struttura dei componenti Vue](#struttura-dei-componenti-vue)
- [Composables & Error Handler](#composables--error-handler)
- [Pinia Stores & Caching](#pinia-stores--caching)
- [API Service Layer & Router Integration](#api-service-layer--router-integration)
- [Gestione Errori Centralizzata](#gestione-errori-centralizzata)
- [Internazionalizzazione & Accessibilità (a11y)](#internazionalizzazione--accessibilità-a11y)
- [Performance](#performance)

---

## Principi generali

1. **Composition API over Options API**: tutti i nuovi componenti usano `<script setup>`.
2. **Un componente, una responsabilità**: componenti piccoli e focalizzati (~200 righe max).
3. **Logica nei composables**: usa `src/composables/useErrorHandler.js` per una gestione omogenea delle eccezioni.
4. **Stato globale in Pinia**: usa lo store di stato (`useGradesStore`, `useErrorStore`, `useWebSocketStore`) con supporto a caching a memoria e refresh trasparente.
5. **Reindirizzamento SPA flessibile**: l'interceptor Axios usa l'istanza iniettata di Vue Router (`setApiRouter`) per evitare reloads di pagina hard.

---

## Composables & Error Handler

### `useErrorHandler`

Utilizza il composables `useErrorHandler` per intercettare e gestire gli errori nei componenti:

```js
import { useErrorHandler } from "@/composables/useErrorHandler";

const { handleError } = useErrorHandler();

try {
  await gradeService.updateGrade(id, payload);
} catch (err) {
  handleError(err, "Impossibile aggiornare il voto selezionato.");
}
```

---

## Pinia Stores & Caching

### In-Memory Store Caching in `useGradesStore`

Per velocizzare la navigazione del docente tra tab senza rieseguire fetch HTTP superflui:

```js
const gradesStore = useGradesStore();

// Utilizza i dati in cache se disponibili per la coppia (classId, subjectId)
await gradesStore.fetchGrades(classId, subjectId);

// Forzo il refetch (es. dopo mutazione o su pull-to-refresh)
await gradesStore.fetchGrades(classId, subjectId, true);
```

---

## Gestione Errori Centralizzata

Gli errori di rete transitori (500, timeout, rate-limit) vengono raccolti dallo store globale `useErrorStore` e segnalati con toast Quasar a scomparsa.

---

## WebSocket & Real-time Error Exposure

Lo store `useWebSocketStore` espone ref reattivi per lo stato della connessione in tempo reale:

- `isConnected`: boolean reattivo
- `reconnectAttempts`: numero di tentativi effettuati
- `hasFailedPermanently`: true se sono stati superati 30 tentativi falliti
- `lastError`: stringa di errore per l'interfaccia utente

---

## Internazionalizzazione & Accessibilità (a11y)

1. **i18n Multi-Lingua**: i file `src/i18n/` coprono 9 lingue distinte: Italiano (`it-IT`), Inglese (`en-US`), Tedesco (`de-DE`), Francese (`fr-FR`), Spagnolo (`es-ES`), Russo (`ru-RU`), Ucraino (`uk-UA`), Arabo (`ar-SA`) e Cinese Semplificato (`zh-CN`).
2. **WAI-ARIA**: inclusione obbligatoria di `role="alert"`, `role="navigation"`, `role="banner"`, `aria-expanded` e Skip Links (_Salta al contenuto principale_ `#main-content`).

---

## Impostazioni Utente, MFA (2FA) & Filtro Anno Scolastico

1. **Autenticazione a Due Fattori (MFA / 2FA TOTP)**:
   - Componente `MfaSetupModal.vue` utilizzabile da tutte le impostazioni utente (`Settings.vue` per ogni ruolo).
   - Generazione dinamica QR code via secret TOTP e verifica codice di conferma a 6 cifre.
2. **Pannello Impostazioni Utente (`Settings.vue`)**:
   - Rotte dedicate per ciascun ruolo integrate nel menu reattivo sidebar `useMenuItems`.
   - Sezioni per Cambio Password, 2FA TOTP, Selezione Lingua (con persistenza `localStorage.setItem('user_locale')`), Preferenze Notifiche e Layout Registro.
3. **Filtraggio Dinamico per Anno Scolastico (`src/stores/schoolYear.js`)**:
   - Genera gli anni scolastici disponibili a partire dall'anno di registrazione dell'utente (`user.created_at`) fino all'anno attivo.
   - Trasmette reattivamente l'anno scelto a `Classes.vue`, `Grades.vue`, `Attendance.vue`, `CoordinatorView.vue` e `Scrutiny.vue`.

---

## 🚀 Sistema di Help Center, Onboarding & Guida In-App

Il frontend integra un sistema multilivello di assistenza e onboarding guidato per gli utenti di ogni ruolo:

1. **Onboarding Tour Interattivo (`OnboardingTour.vue`)**:
   - Tutorial a 8 passaggi specifico per ruolo (`teacher`, `student`, `parent`, `secretary`, `admin`).
   - Card grafiche con pillole di funzionalità, badge di categoria, suggerimenti pratici e lista puntata.
   - Tracciamento completamento automatizzato in `localStorage` (`onboarding_done_{role}`).
2. **Pannello Help Center Full-Screen (`HelpCenterPanel.vue`)**:
   - Finestra modale a schermo intero con ricerca full-text istantanea in tutte le guide.
   - Navigazione per sezioni tematiche (`GUIDE_DEFS`) e mapping dinamico delle domande frequenti (`FAQ_CATEGORY_MAP`).
   - Parser automatico dei contenuti con formattazione grafica differenziata per `Passaggio N:`, `Suggerimento:`, `Attenzione:` e `Scorciatoia:`.
3. **Help Drawer Laterale & FAB Flottante (`HelpDrawer.vue`, `HelpFab.vue`)**:
   - Floating Action Button `?` in basso a destra con menu rapido (Riavvia Tour, Apri Guida, Contatta Assistenza).
   - Drawer laterale per consultare FAQ e guide rapide senza abbandonare l'attività corrente.
4. **Centro Supporto & Form Ticket (`Support.vue`)**:
   - Pagina principale con FAQ accordion filtrate per ruolo utente e modalità offline (coda locale ticket per l'invio al ripristino connessione).

---

## 🛠️ Architettura dei Composables & Safety Guidelines

I 28 composabili in `src/composables/` incapsulano la logica reattiva, l'interazione con gli store e la gestione delle operazioni asincrone:

1. **`useAuth.js`**: Gestione autenticazione, persistenza sessione, traduzione granulare degli errori di login e instradamento post-login per tutti i 13 ruoli supportati (`superadmin`, `admin`, `secretary`, `principal`, `vice_principal`, `staff`, `coordinator`, `teacher`, `student`, `parent`, `docente`, `system_auditor`).
2. **`useUserManagement.js`**: Gestione utenti di segreteria/amministrazione con gestione esplicita degli errori di importazione CSV (notifica negativa `type: 'negative'` e ritorno booleano affidabile `false`).
3. **`useColloquiScheduling.js`**: Generatore di slot disponibilità con validazione preventiva dei limiti di durata (`duration > 0`) e coerenza temporale (`startTime < endTime`).
4. **`useDraftAutosave.js`**: Salvataggio automatico debounced in `localStorage` con deep-clone serializzabile, ripristino istantaneo e cleanup automatizzato su unmount del componente.
5. **`useUndoToast.js`**: Notifiche transitorie con countdown visivo e possibilità di revocare l'azione entro 15 secondi (`notifyWithUndo`).
6. **`useGradeFormatter.js`**: Formattazione coerente dei voti (decimale con virgola/punto, frazionario o centesimale) in base alle preferenze dell'istituto.

---

## 🧩 Linee Guida per i Componenti Vue (`src/components/`)

1. **Gestione Sicura dei Valori di Progresso**: In componenti come `GradeChart.vue`, sanitizzare sempre i valori numerici (`getProgressValue`, `getProgressColor`) per evitare calcoli `NaN / 10` e warning di rendering Quasar in caso di medie non disponibili (`'-'`).
2. **Protezione Multi-Stream Asincrono**: In componenti come `TimelineActivityFeed.vue`, utilizzare sempre `Array.isArray()` sulle risposte aggregate (`Promise.allSettled`) prima di invocare metodi come `.slice()` o `.map()`.
3. **Validazione Props & Default Fallbacks**: Dichiarare sempre prop types espliciti e valori di default per array e oggetti (`() => []`, `() => ({})`).

---

## 🎨 Dark Mode & Sistema di Design Globale (`globals.css`)

Il registro implementa un sistema coerente di tema scuro (`.body--dark`) gestito da Quasar e dallo store `useThemeStore`:

1. **Variabili CSS di base**:
   - `--bg-primary`: `#0f172a` (sfondo pagina scuro).
   - `--bg-secondary`: `#1e293b` (sfondo card, modali, navbar e tabelle).
   - `--border-color`: `#334155` (bordi e divisori).
   - `--text-primary`: `#f1f5f9` (testi principali e titoli).
   - `--text-secondary`: `#94a3b8` (testi secondari, didascalie e placeholder).
2. **Override dei contenitori chiari**: In modalità dark, tutte le varianti di colore chiare usate nei componenti e nelle impostazioni (`bg-slate-50`, `bg-slate-100`, `bg-indigo-50`, `bg-amber-50`, `bg-emerald-50`, `bg-blue-50`, ecc.) vengono mappate automaticamente su `--bg-secondary` o `rgba(30, 41, 59, 0.85)` con testi e bordi scuri ad alto contrasto.
3. **Evidenziazione Giorno Odierno (`Agenda.vue`)**:
   - La cella e la colonna del giorno corrente (`isTodayIso`) adottano classi semantiche dedicate (`today-cell`, `today-header-cell`, `today-slot`) con badge numerico blu ad alto contrasto e bordo perimetrale a 2px (`#2563eb`), garantendo massima leggibilità sia in modalità chiara che in Dark Mode.

---

## 🧪 Unit Testing dei Componenti, Composabili & Servizi

La suite di test frontend è sviluppata con **Vitest** e **Vue Test Utils**:

- **Comando di esecuzione**: `npm run test:unit`
- **Metriche**: **153 test suite**, **928 unit test passati al 100%** (0 errori, 0 fallimenti).
- **Copertura**:
  - `tests/unit/components/`: Test dedicati per componenti Admin, Common, Parent, Secretary, Student e Teacher (`StudentGradeChart.spec.js`, `TimelineActivityFeed.spec.js`, `ScheduleGridsRobustness.spec.js`, `GradeWeights.spec.js`, ecc.).
  - `tests/unit/composables/`: Test dedicati per tutti i composabili (`useUserManagementFix.spec.js`, `useAuthRoleRouting.spec.js`, `useDraftAutosave.spec.js`, `usePermissions.spec.js`, ecc.).
  - `tests/unit/security/`: Test di anti-regressione RBAC, XSS DOMPurify sanitization, CSV injection prevention, route guards e token security.

