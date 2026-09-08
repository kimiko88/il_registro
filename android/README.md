# il_registro — Applicazioni Mobile Native Android

> [!WARNING]
> **STATO ALPHA — NON STABILE E INCOMPLETO**  
> Le applicazioni Android native si trovano attualmente in **FASE ALPHA**. Il codice è sperimentale, in sviluppo attivo, **non stabile e incompleto**. **NON sono pronte né idonee per l'utilizzo in ambienti di produzione**.

> [!IMPORTANT]
> **LICENZA**  
> Anche tutte le applicazioni Android native (per tutti i moduli: Studente, Genitore, Docente e Segreteria) sono rilasciate sotto la medesima licenza dell'intero progetto: **[PolyForm Noncommercial License 1.0.0](../LICENSE)**. L'uso per scuole pubbliche, università ed enti istituzionali pubblici è gratuito e senza limitazioni; i diritti commerciali sono riservati.

---

## 📱 Moduli del Progetto

Questo progetto Android è strutturato come **progetto multi-modulo Gradle** (`settings.gradle.kts`):

- **`:student`**: Applicazione per lo Studente (voti, medie, assenze, compiti, orario, timeline)
- **`:parent`**: Applicazione per il Genitore (giustificazioni, monitoraggio figli, colloqui, comunicazioni)
- **`:teacher`**: Applicazione per il Docente (firma ora 1-click, registro voti con matrix, presenze, note)
- **`:secretary`**: Applicazione per la Segreteria (anagrafica, gestione classi/docenti, circolari, export SIDI)

Ogni applicazione è autonoma, compilabile singolarmente, ed è realizzata in **Kotlin** con **Jetpack Compose** e Material 3.

---

## 🚀 Come Avviare le Applicazioni

### Con Android Studio (Consigliato)
1. Apri **Android Studio**.
2. Seleziona **Open** e scegli questa cartella:
   ```text
   ./android
   ```
3. Attendi la sincronizzazione di Gradle e il download delle librerie.
4. Nel selettore dei moduli in alto, scegli quale app avviare (`student`, `parent`, `teacher` o `secretary`).
5. Seleziona un emulatore (es. API 30+) oppure un dispositivo fisico collegato (con Debug USB attivo).
6. Premi **Run ▶** (`Shift + F10`).

### Da Riga di Comando (Gradle CLI)
```bash
# Compila l'app Studente in debug
./gradlew :student:assembleDebug

# Installa ed avvia su dispositivo o emulatore collegato
./gradlew :student:installDebug

# Esegui tutti i test unitari
./gradlew test

# Esegui i test di una specifica app
./gradlew :student:test
```

---

## 🌐 Configurazione Backend API

I parametri di connessione sono definiti nel file `config/AppConfig.kt` di ciascun sottomodulo:
- Studente: `student/src/main/java/it/scuola/registro/student/config/AppConfig.kt`
- Genitore: `parent/src/main/java/it/scuola/registro/parent/config/AppConfig.kt`
- Docente: `teacher/src/main/java/it/scuola/registro/teacher/config/AppConfig.kt`
- Segreteria: `secretary/src/main/java/it/scuola/registro/secretary/config/AppConfig.kt`

```kotlin
// Default per emulatore Android (punta al localhost del PC):
var BASE_URL: String = "http://10.0.2.2:8080/api/v1"
var WS_URL: String = "ws://10.0.2.2:8080/api/v1/ws"

// Per smartphone fisico connesso sulla stessa rete Wi-Fi del computer:
// var BASE_URL: String = "http://192.168.1.xxx:8080/api/v1"
```

---

## 🌍 Localizzazione (11 Lingue)
Tutti i moduli contengono risorse tradotte in:
Italiano (`values`), Inglese (`values-en`), Spagnolo (`values-es`), Francese (`values-fr`), Tedesco (`values-de`), Rumeno (`values-ro`), Albanese (`values-sq`), Arabo (`values-ar` con supporto RTL), Cinese (`values-zh`), Ucraino (`values-uk`), Russo (`values-ru`).

---

## 📚 Documentazione Correlata
- [Guida Completa di Setup & Architettura Mobile](../docs/MOBILE_SETUP_GUIDE.md)
- [Istruzioni Operative di Avvio e Testing Mobile](../docs/mobile_instruction.md)
- [Architettura di Sistema](../docs/ARCHITECTURE.md)
