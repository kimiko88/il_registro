# Guida di Configurazione ed Esecuzione delle Applicazioni Mobile Native (Android & iOS)

> [!WARNING]
> **STATO ALPHA — NON STABILE E INCOMPLETO**
> Le applicazioni mobile native per **Android** e **iOS** si trovano attualmente in **FASE ALPHA**. Il codice è sperimentale, in corso di sviluppo attivo, **non stabile e incompleto**. Molte funzionalità sono parziali o in fase di test e **NON sono assolutamente consigliate né pronte per l'uso in produzione** in ambienti scolastici reali.

> [!IMPORTANT]
> **LICENZA DELLE APPLICAZIONI MOBILE**
> Anche tutte le applicazioni mobile native (per tutti i ruoli: Studente, Genitore, Docente e Segreteria) sono rilasciate sotto la **stessa identica licenza dell'applicazione web e del backend**: la **[PolyForm Noncommercial License 1.0.0](../LICENSE)**. L'utilizzo per scuole pubbliche, università ed enti istituzionali pubblici è libero e gratuito; i diritti commerciali sono riservati.

---

Questa guida illustra l'architettura, le modalità di configurazione dell'URL delle API, e le istruzioni per la compilazione, il testing e l'esecuzione delle **8 applicazioni mobile native** per **Android** e **iOS** del Registro Elettronico.

---

## 📱 Struttura del Progetto Mobile

Il progetto adotta un'architettura nativa per ciascuna piattaforma, con moduli e target separati per ogni ruolo scolastico:

```
Registrov2/
├── android/                         # Progetto Multi-Modulo Android (Kotlin & Jetpack Compose)
│   ├── settings.gradle.kts         # Radice Gradle con moduli (:student, :parent, :teacher, :secretary)
│   ├── build.gradle.kts            # Configurazione build script root
│   ├── student/                    # App Studente Android (Kotlin Compose)
│   ├── parent/                     # App Genitore Android (Kotlin Compose)
│   ├── teacher/                    # App Docente Android (Kotlin Compose)
│   └── secretary/                  # App Segreteria Android (Kotlin Compose)
│
└── ios/                             # Progetto iOS (Swift & SwiftUI)
    ├── Package.swift               # Swift Package Manager Manifest (StudentApp, ParentApp, TeacherApp, SecretaryApp)
    ├── RegistroStudente/           # Progetto Xcode integrato
    │   ├── RegistroStudente.xcodeproj # Progetto Xcode con schemi per tutti i 4 ruoli
    │   ├── RegistroStudente/       # Target Eseguibile Studente
    │   ├── RegistroDocente/        # Target Eseguibile Docente
    │   ├── RegistroGenitore/       # Target Eseguibile Genitore
    │   ├── RegistroSegreteria/     # Target Eseguibile Segreteria
    │   ├── RegistroStudenteTests/  # Suite di Unit Test iOS
    │   └── RegistroStudenteUITests/# Suite di UI Automation Test iOS
    ├── student/                    # Sorgenti modulo Studente iOS
    ├── parent/                     # Sorgenti modulo Genitore iOS
    ├── teacher/                    # Sorgenti modulo Docente iOS
    └── secretary/                  # Sorgenti modulo Segreteria iOS
```

---

## 🌐 Configurazione dell'URL dell'API Backend

Tutte le applicazioni comunicano direttamente con le API REST e i WebSocket del backend Go (senza mock data). Gli endpoint di connessione sono centralizzati nei file di configurazione dedicati:

### 1. 🤖 Android (Kotlin)
Troverai il file **`AppConfig.kt`** all'interno del package `config` di ciascun sottomodulo:
- Studente: `android/student/src/main/java/it/scuola/registro/student/config/AppConfig.kt`
- Genitore: `android/parent/src/main/java/it/scuola/registro/parent/config/AppConfig.kt`
- Docente: `android/teacher/src/main/java/it/scuola/registro/teacher/config/AppConfig.kt`
- Segreteria: `android/secretary/src/main/java/it/scuola/registro/secretary/config/AppConfig.kt`

```kotlin
package it.scuola.registro.student.config

object AppConfig {
    // 1. Emulatore Android Locale (Default):
    var BASE_URL: String = "http://10.0.2.2:8080/api/v1"
    var WS_URL: String = "ws://10.0.2.2:8080/api/v1/ws"

    // 2. Dispositivo Fisico Android (tramite IP locale Wi-Fi del computer host):
    // var BASE_URL: String = "http://192.168.1.100:8080/api/v1"

    // 3. Produzione / Server HTTPS:
    // var BASE_URL: String = "https://registro.tuascuola.it/api/v1"
}
```

### 2. 🍏 iOS (Swift)
Troverai il file **`AppConfig.swift`** all'interno della cartella `Config` di ciascun modulo/target:
- Studente: `ios/student/Config/AppConfig.swift`
- Genitore: `ios/parent/Config/AppConfig.swift`
- Docente: `ios/teacher/Config/AppConfig.swift`
- Segreteria: `ios/secretary/Config/AppConfig.swift`

```swift
import Foundation

public struct AppConfig {
    // 1. Simulatore iOS Locale (Default):
    public static var baseURL: String = "http://localhost:8080/api/v1"
    public static var wsURL: String = "ws://localhost:8080/api/v1/ws"

    // 2. Dispositivo Fisico iPhone (tramite IP locale Wi-Fi del Mac host):
    // public static var baseURL: String = "http://192.168.1.100:8080/api/v1"

    // 3. Produzione / Server HTTPS:
    // public static var baseURL: String = "https://registro.tuascuola.it/api/v1"
}
```

---

## 🛠️ Compilazione ed Esecuzione

### Android (con Android Studio o Gradle)
1. Apri la cartella `android/` con **Android Studio**.
2. Attendi la sincronizzazione Gradle.
3. Seleziona la configurazione desiderata (`student`, `parent`, `teacher` o `secretary`) e lancia su emulatore o device con `Shift + F10`.
4. Da riga di comando:
   ```bash
   cd android
   ./gradlew :student:assembleDebug
   ./gradlew :student:installDebug
   ```

### iOS (con Xcode o SPM su macOS)
1. Apri il progetto in Xcode:
   ```bash
   cd ios/RegistroStudente
   open RegistroStudente.xcodeproj
   ```
2. In alto seleziona uno schema eseguibile:
   - `RegistroStudente`
   - `RegistroDocente`
   - `RegistroGenitore`
   - `RegistroSegreteria`
3. Premi **Cmd + R** per compilare ed eseguire nel simulatore o su un dispositivo collegato.

---

## 🧪 Esecuzione dei Test

### Test Android
Dalla cartella `android/`:
```bash
# Esegui tutti i test unitari di tutti i moduli
./gradlew test

# Esegui i test di una specifica applicazione (es. Studente)
./gradlew :student:test

# Esegui i test End-to-End UI su emulatore o dispositivo connesso
./gradlew connectedAndroidTest
```

### Test iOS
Dalla cartella `ios/`:
```bash
# Esegui tutte le suite di test con Swift Package Manager
swift test

# Oppure apri il progetto in Xcode e premi Cmd + U
```

---

## 🌍 Localizzazione & Internazionalizzazione (11 Lingue)

Tutte le applicazioni native per Android e iOS supportano nativamente **11 lingue** con allineamento completo delle stringhe:
- 🇮🇹 **Italiano** (`values` / `it.lproj`)
- 🇬🇧 **Inglese** (`values-en` / `en.lproj`)
- 🇪🇸 **Spagnolo** (`values-es` / `es.lproj`)
- 🇫🇷 **Francese** (`values-fr` / `fr.lproj`)
- 🇩🇪 **Tedesco** (`values-de` / `de.lproj`)
- 🇷🇴 **Rumeno** (`values-ro` / `ro.lproj`)
- 🇦🇱 **Albanese** (`values-sq` / `sq.lproj`)
- 🇸🇦 **Arabo** (con supporto layout RTL) (`values-ar` / `ar.lproj`)
- 🇨🇳 **Cinese Semplificato** (`values-zh` / `zh-Hans.lproj`)
- 🇺🇦 **Ucraino** (`values-uk` / `uk.lproj`)
- 🇷🇺 **Russo** (`values-ru` / `ru.lproj`)

---

## 📄 Licenza

Le applicazioni mobile native per Android e iOS sono rilasciate con licenza **[PolyForm Noncommercial 1.0.0](../LICENSE)**, la stessa che disciplina il backend e il frontend web:
- **Gratuita e senza limitazioni** per scuole pubbliche, università, enti di ricerca e istituzioni pubbliche.
- **Riservata** per qualsiasi uso commerciale da parte di aziende ed enti privati (richiede licenza commerciale separata).
