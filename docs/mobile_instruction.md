# 📱 Istruzioni per l'Avvio e il Testing Mobile

> [!WARNING]
> **STATO ALPHA — NON STABILE E INCOMPLETO**
> Le applicazioni mobile native per **Android** e **iOS** si trovano attualmente in **FASE ALPHA**. Il codice è sperimentale, in corso di sviluppo attivo, **non stabile e incompleto**. Molte funzionalità sono parziali o in fase di test e **NON sono assolutamente consigliate né pronte per l'uso in produzione** in ambienti scolastici reali.

> [!IMPORTANT]
> **LICENZA DELLE APPLICAZIONI MOBILE**
> Anche tutte le applicazioni mobile native (per tutti i ruoli: Studente, Genitore, Docente e Segreteria) sono rilasciate sotto la **stessa identica licenza dell'applicazione web e del backend**: la **[PolyForm Noncommercial License 1.0.0](../LICENSE)**. L'utilizzo per scuole pubbliche, università ed enti istituzionali pubblici è libero e gratuito; i diritti commerciali sono riservati.

---

Nel progetto sono presenti due modalità di fruizione su dispositivi mobili:

1. **Applicazioni Mobile Native**: sviluppate specificamente per Android (in **Kotlin & Jetpack Compose** nella cartella `android/`) e per iOS (in **Swift & SwiftUI** nella cartella `ios/`), suddivise per i 4 ruoli chiave della scuola.
2. **Applicazione Web Mobile / PWA**: il frontend web in **Vue 3 + Quasar** con layout responsive, barra di navigazione ottimizzata per smartphone e cassetto laterale destro per opzioni e accessibilità.

Di seguito vengono illustrate le istruzioni dettagliate per avviare e testare ciascuna soluzione.

---

## 🤖 1. Come Avviare le Applicazioni Native Android

Il progetto Android è organizzato come progetto multi-modulo Gradle (`settings.gradle.kts`) ed è composto da **4 applicazioni native distinte per ruolo**:

- `:student` — Applicazione Studente (libretto voti, assenze, orario, compiti, agenda)
- `:parent` — Applicazione Genitore (giustificazioni, monitoraggio figli, colloqui, comunicazioni)
- `:teacher` — Applicazione Docente (firma ora, registro voti, presenze, note disciplinari, orario)
- `:secretary` — Applicazione Segreteria (anagrafica, gestione orari, circolari, export SIDI)

### Prerequisiti
- [Android Studio](https://developer.android.com/studio) (versione Hedgehog / Iguana / Ladybug o superiore)
- Android SDK (API Level 30+ consigliato, es. emulatore Pixel con API 30 o 34)
- JDK 17 o superiore

### Procedura con Android Studio:
1. Avvia **Android Studio**.
2. Seleziona **Open** (oppure **File → Open**) e seleziona la cartella `android` situata nella root del repository:
   ```text
   ./android
   ```
3. Attendi il completamento della sincronizzazione di Gradle e il download delle dipendenze.
4. Nella barra superiore delle configurazioni di esecuzione di Android Studio:
   - **Modulo**: seleziona dal menu a tendina l'applicazione che desideri avviare (`student`, `parent`, `teacher` oppure `secretary`).
   - **Dispositivo**: seleziona un emulatore configurato (es. *Pixel_2_API_30*) oppure il tuo smartphone Android collegato via cavo con *Debug USB* attivo.
5. Clicca sul pulsante **Run ▶** (oppure premi `Shift + F10`). L'applicazione verrà compilata, installata ed eseguita automaticamente.

### Procedura da Riga di Comando (Gradle CLI):
Dalla cartella `android/`:
```bash
# Compilazione debug dell'app Studente
./gradlew :student:assembleDebug

# Installazione ed esecuzione su dispositivo/emulatore collegato
./gradlew :student:installDebug

# Esecuzione della suite di test
./gradlew test
```

> [!NOTE]
> **Connessione con il backend in locale (Emulatore Android)**:  
> Quando l'emulatore Android gira sul tuo computer, accede al server host tramite l'IP speciale `http://10.0.2.2:8080/api/v1` (questo endpoint è già preconfigurato in `AppConfig.kt` di ogni modulo). Ricordati di avviare prima il backend Go con `go run cmd/api-server/main.go`.

---

## 🍏 2. Come Avviare le Applicazioni Native iOS

Il progetto iOS si trova nella cartella `ios/` e include sia il progetto Xcode completo (`RegistroStudente.xcodeproj`), sia la configurazione con Swift Package Manager (`Package.swift`).

Include **4 target eseguibili dedicati per ruolo**:
- **`RegistroStudente`** (StudentApp)
- **`RegistroDocente`** (TeacherApp)
- **`RegistroGenitore`** (ParentApp)
- **`RegistroSegreteria`** (SecretaryApp)

### Requisito Fondamentale:
> [!CAUTION]
> Apple richiede **macOS** e **Xcode** per compilare ed eseguire simulatori iOS o installare app su iPhone/iPad. Su sistemi operativi Windows o Linux non è possibile avviare nativamente il simulatore iOS né compilare i bundle `.app`.

### Procedura su macOS con Xcode:
1. Clona il repository su un Mac.
2. Apri il progetto Xcode:
   ```bash
   cd ios/RegistroStudente
   open RegistroStudente.xcodeproj
   ```
   *(In alternativa, da terminale puoi aprire la directory in Xcode con `xed .`)*
3. Nella barra superiore di Xcode:
   - Seleziona lo **Scheme** desiderato dal menu a tendina in alto (es. `RegistroStudente`, `RegistroDocente`, `RegistroGenitore` o `RegistroSegreteria`).
   - Seleziona il dispositivo di destinazione (es. simulatore *iPhone 16 Pro* o un iPhone fisico collegato).
4. Premi **Cmd + R** per compilare e avviare l'app nel simulatore iOS.
5. Per eseguire i test unitari e di interfaccia, premi **Cmd + U**.

### Procedura da Riga di Comando con Swift Package Manager (SPM):
Dalla cartella `ios/`:
```bash
# Esecuzione di tutti i test della suite iOS
swift test

# Build del package
swift build
```

> [!NOTE]
> **Connessione con il backend in locale (Simulatore iOS)**:  
> Nel simulatore iOS l'indirizzo `http://localhost:8080/api/v1` punta direttamente alla macchina host Mac (preconfigurato in `AppConfig.swift`).

---

## 📱 3. Come Testare l'Applicazione Web su Mobile (Responsive & PWA)

L'interfaccia web del registro è progettata come PWA (Progressive Web App) con supporto touch-friendly, navbar mobile dedicata e cassetto laterale destro per impostazioni e accessibilità.

### Opzione A: Emulatore Mobile nel Browser (da PC)
1. Con il frontend attivo (`npm run dev` nella cartella `registro-frontend`), apri Google Chrome, Microsoft Edge o Firefox all'indirizzo:
   ```text
   http://localhost:5173
   ```
2. Premi il tasto **F12** (oppure clic destro → *Ispeziona*).
3. Attiva la modalità dispositivo mobile cliccando sull'icona **Toggle device toolbar** (in alto a sinistra dei DevTools, oppure premi `Ctrl + Shift + M`).
4. Nel menu a tendina in alto seleziona un dispositivo mobile (es. *iPhone 14 Pro*, *Pixel 7*, *Samsung Galaxy S20*).
5. Potrai verificare la visualizzazione mobile con:
   - Titolo compatto e layout verticale touch-friendly.
   - Pulsante laterale **Opzioni & Accessibilità** (`tune 🎛️`) che apre il drawer laterale destro con controllo contrasto, font DSA OpenDyslexic, righello di lettura e sintesi vocale.

### Opzione B: Da uno Smartphone Reale (tramite Wi-Fi Locale)
1. Assicurati che il tuo smartphone e il computer siano connessi alla **stessa rete Wi-Fi**.
2. Avvia il frontend permettendo le connessioni dalla rete locale:
   ```bash
   cd registro-frontend
   npm run dev -- --host
   ```
3. Vite mostrerà a terminale il tuo indirizzo IP locale (es. `http://192.168.1.150:5173`).
4. Apri il browser (Safari su iOS o Chrome su Android) sul tuo smartphone e naviga a quell'indirizzo.
