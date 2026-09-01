# Guida di Configurazione ed Esecuzione delle Applicazioni Mobile Native (Android & iOS)

Questa guida illustra la struttura, le modalità di configurazione dell'URL delle API, e le istruzioni per l'esecuzione dei test e della build delle **8 applicazioni mobile native** per **Android** e **iOS** del Registro Elettronico.

---

## 📱 Struttura del Progetto Mobile

```
Registrov2/
├── android/                   # Progetto Multi-Modulo Android (Kotlin & Jetpack Compose)
│   ├── settings.gradle.kts   # Radice Gradle con sottomoduli (:student, :parent, :teacher, :secretary)
│   ├── build.gradle.kts      # Configurazione build script root
│   ├── student/              # App Studente Android
│   ├── parent/               # App Genitore Android
│   ├── teacher/              # App Docente Android
│   └── secretary/            # App Segreteria Android
│
└── ios/                       # Progetto Multi-Target iOS (Swift & SwiftUI / SPM)
    ├── Package.swift         # Manifest Swift Package Manager (StudentApp, ParentApp, TeacherApp, SecretaryApp)
    ├── student/              # App Studente iOS
    ├── parent/               # App Genitore iOS
    ├── teacher/              # App Docente iOS
    └── secretary/            # App Segreteria iOS
```

---

## 🌐 Dove Configurare l'URL dell'API

I parametri di connessione di rete sono centralizzati in file dedicati per consentire il cambio immediato dell'endpoint backend:

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

    // 2. Dispositivo Fisico Android (tramite IP Wi-Fi locale del PC server):
    // var BASE_URL: String = "http://192.168.1.100:8080/api/v1"

    // 3. Produzione / Server HTTPS:
    // var BASE_URL: String = "https://registro.tuascuola.it/api/v1"
}
```

### 2. 🍏 iOS (Swift)
Troverai il file **`AppConfig.swift`** all'interno della cartella `Config` di ciascuna app:
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

    // 2. Dispositivo Fisico iPhone (tramite IP Wi-Fi locale del Mac server):
    // public static var baseURL: String = "http://192.168.1.100:8080/api/v1"

    // 3. Produzione / Server HTTPS:
    // public static var baseURL: String = "https://registro.tuascuola.it/api/v1"
}
```

---

## 🧪 Esecuzione dei Test

### Test Android
Dalla cartella `android/`:
```bash
# Esegui tutti i test unitari e di integrazione di tutti i moduli
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

## 🌍 Lingue Supportate
Tutte le applicazioni supportano le **11 lingue** gestite dalla piattaforma:
Italiano, Inglese, Tedesco, Francese, Spagnolo, Arabo (RTL), Rumeno, Russo, Albanese, Ucraino, Cinese.
