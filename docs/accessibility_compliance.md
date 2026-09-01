# Documentazione Tecnica: Suite di Accessibilità (A11y), Supporto DSA & Conformità AgID / WCAG 2.2

La presente documentazione descrive le funzionalità di **Accessibilità Avanzata (A11y)**, gli strumenti compensativi per i **Disturbi Specifici dell'Apprendimento (DSA - Legge 170/2010)** e la conformità alle direttive **AgID (Agenzia per l'Italia Digitale)** e **WCAG 2.1 / 2.2 (Livello AA e AAA)** implementate nel Registro Elettronico Scolastico.

---

## 🏛️ 1. Quadro Normativo di Riferimento

1. **Legge 9 gennaio 2004, n. 4 ("Legge Stanca")** e s.m.i. (D.Lgs. 106/2018): Disposizioni per favorire e semplificare l'accesso dei soggetti disabili agli strumenti informatici nella Pubblica Amministrazione.
2. **Direttiva UE 2016/2102**: Accessibilità dei siti web e delle applicazioni mobili degli enti pubblici.
3. **Linee Guida AgID sull'Accessibilità**: Adozione delle raccomandazioni internazionali **W3C WCAG 2.1 e 2.2** (Livello AA e AAA).
4. **Legge 8 ottobre 2010, n. 170**: Nuove norme in materia di disturbi specifici di apprendimento in ambito scolastico.

---

## 🧩 2. Architettura dei Moduli di Accessibilità

```
registro-frontend/src/
├── stores/
│   └── theme.js                          # State management A11y (TTS, Righello, Contrasti, Font, Spaziatura)
├── composables/
│   ├── useSpeechSynthesis.js             # Gestione Web Speech API (TTS nativo)
│   └── useGlobalKeyboardShortcuts.js     # Navigazione rapida e scorciatoie globali da tastiera
├── components/Common/
│   ├── AccessibilitySettingsPanel.vue    # Pannello di configurazione completo per tutti i ruoli
│   ├── TextToSpeechButton.vue            # Pulsante audio compatto riutilizzabile (circolari, note, lezioni)
│   ├── ReadingRuler.vue                  # Righello di lettura / Focus Mask semi-trasparente
│   └── KeyboardShortcutsDialog.vue       # Dialogo interattivo delle scorciatoie da tastiera
├── pages/
│   └── AccessibilityStatement.vue        # Pagina ufficiale Dichiarazione di Accessibilità AgID
└── css/
    └── App.vue                           # Regole CSS per filtri daltonismo, OLED Pure Black, Focus Ring e Spaziatura
```

---

## 🚀 3. Funzionalità Implementate

### 3.1. Sintesi Vocale Nativa (Text-to-Speech - TTS)
- **Tecnologia**: Web Speech API nativa (`window.speechSynthesis`), senza dipendenze o costi esterni.
- **Pulsante Audio Inclusivo**: Integrato nelle Circolari, Comunicazioni Famiglia, Note Disciplinari e Bacheca.
- **Regolazione**: Velocità (0.7x – 1.5x) e tono regolabili dallo store Pinia.
- **Comando Rapido**: Scorciatoia globale `Alt + T` per leggere il testo selezionato o l'intestazione principale.

### 3.2. Righello di Lettura & Focus Mask (Reading Ruler)
- **Scopo**: Fornisce una guida orizzontale a contrasto che segue il mouse o le frecce della tastiera (`Alt + ↑/↓`), attenuando il resto dello schermo.
- **Parametri configurabili**:
  - Altezza della fessura di lettura: `24px` – `70px`.
  - Opacità della maschera: `10%` – `80%`.
- **Attivazione Rapida**: `Alt + R` o menu rapido Navbar.

### 3.3. Supporto DSA & Font ad Alta Leggibilità
- **Font OpenDyslexic**: Font studiato con peso alla base della lettera per contrastare la rotazione e l'inversione dei caratteri.
- **Font Alternativi Inclusivi**: *Lexend* (ottimizzato per la fluidità di lettura), *Fredoka* (morbido per la scuola primaria) e *Roboto*.
- **Spaziatura Testo (WCAG 1.4.12)**:
  - *Interlinea*: Standard (1.5x), Rilassata (1.8x), Ampia (2.1x).
  - *Spaziatura lettere*: Standard, Ampia (+0.08em), Molto Ampia (+0.16em).
  - *Spaziatura parole*: Standard, Ampia (+0.18em), Molto Ampia (+0.32em).

### 3.4. Daltonismo & Contrasti Avanzati (Color Blindness & OLED)
- **Eliminazione dipendenza esclusiva dal colore (WCAG 1.4.1)**: Tutti i voti e gli esiti associano pattern visivi (badge con icona `▼` per insufficienze e `✓` per sufficienze) oltre al colore rosso/verde.
- **Filtri Ottici per Daltonismo**:
  - *Protanopia* (insensibilità al rosso).
  - *Deuteranopia* (insensibilità al verde).
  - *Tritanopia* (insensibilità al blu/giallo).
  - *Monocromatico* (scala di grigi).
- **Contrasti OLED & Fotofobia**:
  - *OLED Pure Black con Testo Ambra/Giallo* (rapporto di contrasto > 21:1).
  - *OLED Pure Black con Testo Verde Fosforo*.
  - *Colori Invertiti (Negativo)*.

### 3.5. Navigazione da Tastiera & Focus Indicator (WCAG 2.1.1 / 2.4.7)
- **Visible Focus Ring**: Contorno ad altissima visibilità con ombra ambra (`outline: 3px solid #f59e0b; box-shadow: 0 0 0 5px rgba(245, 158, 11, 0.3)`) su tutti i campi di input, pulsanti e link.
- **Scorciatoie da Tastiera Globali**:
  - `Alt + 1` → Vai alla Dashboard
  - `Alt + V` → Vai al Registro Voti
  - `Alt + P` → Vai al Registro Presenze
  - `Alt + A` → Vai all'Agenda & Compiti
  - `Alt + S` → Attiva Ricerca Globale
  - `Alt + R` → Attiva/Disattiva Righello di Lettura
  - `Alt + T` → Attiva/Disattiva Lettura Vocale (TTS)
  - `?` (o `Shift + /`) → Apri Dialogo Guida Scorciatoie
  - `Escape` → Chiudi modale o interrompi lettura vocale

### 3.6. Dettatura Vocale (Speech-to-Text - STT)
- **Tecnologia**: Web Speech Recognition API nativa (`SpeechRecognition` / `webkitSpeechRecognition`), composable [`useSpeechToText.js`](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-frontend/src/composables/useSpeechToText.js).
- **Componente**: [`SpeechToTextButton.vue`](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-frontend/src/components/Common/SpeechToTextButton.vue) con animazione pulse di registrazione, indicatore visivo di stato, tooltip ed attivazione rapida via `Alt + D`.
- **Casi d'uso**: Inserimento vocale mani-libere di note disciplinari, giudizi analitici nei PDP/PEI, relazioni per gli scrutini, verbali e giustificazioni per utenti con disabilità motorie agli arti superiori o disgrafia.

### 3.7. Sincronizzazione Cloud delle Preferenze Accessibilità (Cross-Device)
- **Database**: Tabella PostgreSQL `user_accessibility_preferences` (migrazione [`096_create_user_accessibility_preferences.sql`](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-backend/migrations/096_create_user_accessibility_preferences.sql)).
- **Endpoint REST**:
  - `GET /api/v1/user/accessibility-settings` *(Recupera il layout A11y dell'utente autenticato)*
  - `PUT /api/v1/user/accessibility-settings` *(Salva/aggiorna le preferenze A11y nel profilo cloud)*
- **Funzionamento**: Quando uno studente DSA o docente ipovedente effettua il login da qualsiasi postazione (LIM, computer della scuola, tablet o PC di casa), le sue preferenze visive e di lettura si applicano automaticamente.

### 3.8. Modalità "Focus / Lettura Pulita" per Alunni ADHD e DSA
- **Scopo**: Attivazione con un clic (`FocusModeToggle.vue`) per eliminare il carico cognitivo (*cognitive overload*).
- **Stile Visivo**: Applicazione della classe `.focus-mode-active` che nasconde i menu laterali ed i widget decorativi, centrando il testo della lezione o del circolare con larghezza massima di 920px e spaziatura rilassata (interlinea 1.85x).

### 3.9. Sintesi Narrata e Tabelle Accessibili per i Grafici (WCAG 1.1.1)
- **Componente**: [`AccessibleChartSummary.vue`](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-frontend/src/components/Common/AccessibleChartSummary.vue).
- **Descrizione**: Converte i grafici dell'andamento dei voti e delle presenze in una sintesi in linguaggio naturale (es: *"Media 7.4, 12 elementi analizzati, Voto Max 9..."*) ed una tabella HTML ad alto contrasto (`q-table`) navigabile da tastiera e lettori dello schermo.

### 3.10. Annunciatore Dinamico per Screen Reader (`aria-live` Announcer)
- **Composable & Componente**: [`useA11yAnnouncer.js`](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-frontend/src/composables/useA11yAnnouncer.js) e [`ScreenReaderAnnouncer.vue`](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-frontend/src/components/Common/ScreenReaderAnnouncer.vue).
- **Funzionamento**: Gestisce regioni visivamente nascoste `aria-live="polite"` e `aria-live="assertive"` per annunciare vocalmente agli utenti non vedenti (NVDA/JAWS/VoiceOver) il salvataggio di voti, la registrazione delle presenze e le notifiche WebSocket in tempo reale.

### 3.11. Skip Links Avanzati e Navigazione Tastiera (WCAG 2.4.1)
- **Componente**: [`SkipLinks.vue`](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-frontend/src/components/Common/SkipLinks.vue) integrato in top di [`MainLayout.vue`](file:///c:/Users/chimi/Desktop/Programmazione/Registrov2/registro-frontend/src/layouts/MainLayout.vue).
- **Scorciatoie**: Compare al primo impulso del tasto `TAB` consentendo il salto immediato al contenuto principale (`#main-content`), al menu di navigazione (`#main-nav`) ed al pannello accessibilità (`#a11y-panel`).

---

## 3.12. Dichiarazione di Accessibilità AgID & Meccanismo di Feedback End-to-End
- **Pagina Dedicata Frontend**: `/accessibility-statement` (e alias `/dichiarazione-accessibilita`).
- **Sezioni Principali**:
  1. Stato di conformità (WCAG 2.2 AA / AAA).
  2. Tecnologie e strumenti compensativi integrati.
  3. Modulo interattivo di segnalazione barriere digitali al Responsabile della Transizione Digitale (RTD).
  4. Procedura di attuazione con il Difensore Civico per il Digitale (AgID).
- **Backend & Database**:
  - **Tabella PostgreSQL**: `accessibility_feedbacks` (migrazione `095_create_accessibility_feedback.sql`).
  - **Generazione Protocollo AgID**: Generazione univoca del codice `A11Y-YYYY-MMDD-XXXX` (es. `A11Y-2026-0827-1042`) restituito al cittadino.
  - **Endpoint REST**:
    - `POST /api/v1/public/accessibility-feedback` *(Aperto al pubblico e a utenti anonimi)*
    - `POST /api/v1/accessibility/feedback` *(Per utenti autenticati)*
    - `GET /api/v1/admin/accessibility-feedbacks` *(Consultazione per amministratori e RTD)*
    - `PATCH /api/v1/admin/accessibility-feedbacks/:id/status` *(Presa in carico / Risoluzione)*

---

## 🧪 4. Test e Manutenzione

La suite è coperta al 100% da test unitari automatici su frontend e backend:

```bash
# Frontend (Vitest - Suite A11y e Test Globali)
cd registro-frontend
npm run test:unit

# Backend (Go Test)
cd registro-backend
go test ./internal/accessibility/...
```

File di test:
- Frontend: `registro-frontend/tests/unit/accessibility/accessibilitySuite.spec.js` e `newA11yFeatures.spec.js`.
- Backend: `registro-backend/internal/accessibility/service_test.go`.


