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

## Unit Testing dei Componenti & Servizi 🧪

I test unitari del frontend sono eseguiti con **Vitest** e **Vue Test Utils**.

- **Esecuzione**: `npm run test:unit` oppure `npx vitest run`
- **Suite Onboarding & Help**: `tests/unit/components/Common/HelpAndOnboarding.spec.js` (verifica inizializzazione tour, persistenza localStorage, apertura/chiusura HelpCenterPanel e HelpDrawer).
- **Suite Componenti Orario**: `tests/unit/components/Secretary/TimetableManagement.spec.js` (verifica rendering `ScheduleGrid.vue`, `TeacherScheduleGrid.vue`, calcolo ore settimanali ed eventi `save`).
- **Suite Servizi Admin**: `tests/unit/services/adminService.spec.js` (verifica chiamate API `/teachers/:id/schedule` e risposte).

