# Relazione Motivazionale di Eccezione di Licenza per la PA (Linee Guida AgID ex art. 69 CAD)

**Progetto**: `il_registro — Registro Elettronico Scolastico`  
**Licenza Software adottata**: PolyForm Noncommercial License 1.0.0  
**Licenza Metadati e Documentazione**: Creative Commons Attribution 4.0 International (CC-BY-4.0)  
**Riferimento Normativo**: Linee Guida dell'Agenzia per l'Italia Digitale (AgID) su acquisizione e riuso del software nella Pubblica Amministrazione (Sezione 1 — "Definire la licenza del software").

---

## 1. Premessa e Inquadramento Normativo

Il punto 1 delle Linee Guida AgID sancisce che la licenza di preferenza per il software sviluppato o riutilizzato dalla Pubblica Amministrazione debba appartenere a licenze aperte riconosciute (es. EUPL, AGPL, MIT, Apache 2.0). Tuttavia, le medesime Linee Guida prevedono espressamente che:

> *"L’ente deve motivare eventuali eccezioni qualora non fosse possibile pubblicare il software come open source [OSI/FSF]."*

La presente relazione costituisce il documento formale di motivazione della scelta della licenza **PolyForm Noncommercial License 1.0.0** a tutela dell'interesse pubblico della scuola e delle Istituzioni Pubbliche Italiane.

---

## 2. Le Ragioni dell'Eccezione a Favor della Scuola Pubblica

### 2.1 Protezione contro lo Sfruttamento Commerciale Abusivo e il Lock-In Privato
Il settore del software scolastico in Italia è storicamente caratterizzato da una forte concentrazione di mercato da parte di pochi operatori privati a scopo di lucro, che erogano il Registro Elettronico a fronte di canoni annui ricorrenti sostenuti dai bilanci scolastici.

L'adozione della clausola non-commerciale di PolyForm impedisce che società commerciali terze possano effettuare un "fork" del codice sorgente di *il_registro* per confezionarlo e rivenderlo come servizio proprietario a pagamento verso le scuole pubbliche o paritarie.

### 2.2 Diritto di Uso, Modifica e Redistribuzione Illimitato e Gratuito per la PA
La sezione **Permitted Organizations** della licenza PolyForm 1.0.0 recita testualmente:

> *"Use by any charitable organization, educational institution, public research organization, public safety or health organization, environmental protection organization, or government institution is use for a permitted purpose regardless of the source of funding."*

Ciò significa che **ogni Scuola Statale, Istituto Comprensivo, Istituto Superiore, Università, Ente di Ricerca e Pubblica Amministrazione**:
- Gode del diritto **libero, irrevocabile e totalmente gratuito** di eseguire, ospitare, modificare, personalizzare e distribuire il software senza alcun costo di licenza o royalty.
- Mantiene il controllo diretto sui propri dati sensibili e la propria infrastruttura (sovranità digitale).

### 2.3 Sostenibilità Evolutiva del Progetto Comunitario
Trattenendo i diritti di sfruttamento commerciale per i soggetti privati a scopo di lucro, il progetto preserva la capacità di coordinare eventuali partnership industriali sostenibili attraverso contratti di licenza commerciale ad-hoc, il cui eventuale ricavato sarà reinvestito nello sviluppo del software pubblico.

---

## 3. Matrice Sintetica di Conformità

| Requisito AgID / PA | Riscontro con PolyForm Noncommercial 1.0.0 |
| :--- | :--- |
| **Disponibilità del codice sorgente** | ✅ **Sì**: Codice sorgente 100% aperto e consultabile su GitHub. |
| **Diritto di riuso per la Scuola / PA** | ✅ **Sì**: Garantito senza costi ex Sezione *Permitted Organizations*. |
| **Modificabilità e personalizzazione** | ✅ **Sì**: Ogni PA può adattare il codice alle proprie esigenze d'istituto. |
| **Prevenzione del vendor lock-in** | ✅ **Sì**: Impossibile per terzi privati chiudere il codice per rivenderlo a canone. |
| **Rilasciabilità dei metadati** | ✅ **Sì**: Scheda `publiccode.yml` pubblicata sotto licenza aperta CC-BY-4.0. |

---

## 4. Conclusioni

La scelta della licenza **PolyForm Noncommercial 1.0.0** soddisfa appieno le finalità dell'art. 69 del CAD: promuovere il riuso tecnologico, eliminare i costi di licenza a carico delle scuole e garantire la sovranità dei dati della PA, proteggendo al contempo il progetto pubblico da speculazioni commerciali di terzi.
