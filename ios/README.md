# il_registro — Applicazioni Mobile Native iOS

> [!WARNING]
> **STATO ALPHA — NON STABILE E INCOMPLETO**  
> Le applicazioni iOS native si trovano attualmente in **FASE ALPHA**. Il codice è sperimentale, in sviluppo attivo, **non stabile e incompleto**. **NON sono pronte né consigliate per l'utilizzo in ambienti di produzione**.

> [!IMPORTANT]
> **LICENZA**  
> Anche tutte le applicazioni iOS native (per tutti i target: Studente, Genitore, Docente e Segreteria) sono rilasciate sotto la medesima licenza dell'intero progetto: **[PolyForm Noncommercial License 1.0.0](../LICENSE)**. L'uso per scuole pubbliche, università ed enti istituzionali pubblici è gratuito e senza limitazioni; i diritti commerciali sono riservati.

---

## 🍏 Struttura del Progetto iOS

Il progetto iOS supporta sia l'integrazione diretta in **Xcode** che la compilazione tramite **Swift Package Manager (SPM)**:

- **`ios/RegistroStudente/RegistroStudente.xcodeproj`**: Progetto Xcode completo con schemi e target per tutti i 4 ruoli:
  - **`RegistroStudente`**: Applicazione Studente (libretto voti, medie, assenze, timeline, orario)
  - **`RegistroDocente`**: Applicazione Docente (firma ora, registro voti, presenze, note)
  - **`RegistroGenitore`**: Applicazione Genitore (giustificazioni, monitoraggio, colloqui)
  - **`RegistroSegreteria`**: Applicazione Segreteria (anagrafica, orari, SIDI sync)
- **`ios/Package.swift`**: Manifest SPM per build e test automatizzati headless via riga di comando.

---

## 🚀 Come Avviare le Applicazioni

### Requisito di Sistema
> [!CAUTION]
> Per compilare ed eseguire il progetto iOS o avviare i simulatori iPhone è necessario un computer con **macOS** e **Xcode** installato.

### Con Xcode (Consigliato su macOS)
1. Apri il terminale ed entra nella directory del progetto Xcode:
   ```bash
   cd ios/RegistroStudente
   open RegistroStudente.xcodeproj
   ```
2. Nella barra degli schemi in alto a sinistra in Xcode, seleziona il target desiderato (es. `RegistroStudente`, `RegistroDocente`, `RegistroGenitore` o `RegistroSegreteria`).
3. Seleziona un simulatore iOS (es. *iPhone 16 Pro*) oppure un iPhone fisico collegato.
4. Premi **Cmd + R** per compilare ed eseguire.
5. Premi **Cmd + U** per avviare la suite di test unitari e di interfaccia.

### Da Riga di Comando con Swift Package Manager (SPM)
Dalla cartella `ios/`:
```bash
# Esegui tutti i test della suite iOS
swift test

# Build del package
swift build
```

---

## 🌐 Configurazione Backend API

I parametri di connessione sono definiti nel file `Config/AppConfig.swift` di ciascun modulo:
- Studente: `student/Config/AppConfig.swift`
- Genitore: `parent/Config/AppConfig.swift`
- Docente: `teacher/Config/AppConfig.swift`
- Segreteria: `secretary/Config/AppConfig.swift`

```swift
import Foundation

public struct AppConfig {
    // Default per simulatore iOS (accede al localhost del Mac):
    public static var baseURL: String = "http://localhost:8080/api/v1"
    public static var wsURL: String = "ws://localhost:8080/api/v1/ws"

    // Per iPhone fisico connesso sulla stessa rete Wi-Fi:
    // public static var baseURL: String = "http://192.168.1.xxx:8080/api/v1"
}
```

---

## 🌍 Localizzazione (11 Lingue)
Tutti i target includono bundle `Localizable.strings` per:
Italiano (`it.lproj`), Inglese (`en.lproj`), Spagnolo (`es.lproj`), Francese (`fr.lproj`), Tedesco (`de-DE`), Rumeno (`ro.lproj`), Albanese (`sq.lproj`), Arabo (`ar.lproj` con supporto RTL), Cinese (`zh-Hans.lproj`), Ucraino (`uk.lproj`), Russo (`ru.lproj`).

---

## 📚 Documentazione Correlata
- [Guida Completa di Setup & Architettura Mobile](../docs/MOBILE_SETUP_GUIDE.md)
- [Istruzioni Operative di Avvio e Testing Mobile](../docs/mobile_instruction.md)
- [Architettura di Sistema](../docs/ARCHITECTURE.md)
