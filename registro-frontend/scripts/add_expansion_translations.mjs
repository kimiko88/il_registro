import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const i18nDir = path.resolve(__dirname, '../src/i18n');

const translations = {
  'it-IT': {
    teacherHours: {
      title: "Gestione Monte Ore Docenti (Art. 29 CCNL Scuola)",
      subtitle: "Rendicontazione delle 40+40 ore annuali, ore eccedenti e gestione permessi brevi",
      collegioCard: "Attività Funzionali Collegiali (Max 40h)",
      collegioDesc: "Partecipazione al Collegio Docenti, incontri dipartimentali e informazione alle famiglie",
      consigliCard: "Consigli di Classe e Interclasse (Max 40h)",
      consigliDesc: "Consigli di classe, scrutini intermedi e finali",
      planned: "Pianificate",
      done: "Svolte",
      remaining: "Residuo",
      overtimeAlert: "Attenzione: superamento del monte ore contrattuale",
      deficitAlert: "Debito orario da completare entro la fine dell'anno",
      bancaOre: "Banca Ore & Recuperi",
      disposizione: "Ore a Disposizione",
      logActivityTitle: "Registra Attività Svolta",
      category: "Tipologia Attività",
      hours: "Ore",
      date: "Data",
      description: "Descrizione attività",
      shortLeavesTitle: "Permessi Brevi (Max 2h - CCNL)",
      requestLeaveBtn: "Richiedi Permesso Breve",
      shortLeaveNotice: "Ai sensi dell'art. 16 CCNL, i permessi brevi non possono eccedere le 2 ore giornaliere e vanno recuperati entro 2 mesi.",
      deadline: "Scadenza Recupero",
      status: "Stato",
      recoveredHours: "Ore Recuperate",
      recordRecovery: "Registra Recupero Ore"
    },
    healthSafety: {
      title: "Sanità Scolastica & Registro Infortuni INAIL",
      subtitle: "Protocolli somministrazione farmaci salvavita e gestione denunce infortuni 48h",
      medicationsTab: "Somministrazione Farmaci a Scuola",
      accidentsTab: "Registro Infortuni INAIL",
      registerMedBtn: "Registra Piano Terapeutico Farmaco",
      student: "Studente",
      medName: "Nome Farmaco",
      dosage: "Posologia / Dosaggio",
      adminMode: "Modalità di Somministrazione",
      prescriptionExpiry: "Scadenza Certificato Medico",
      medExpiry: "Scadenza Confezione Farmaco",
      emergencyProtocol: "Protocollo di Emergenza",
      authorizedStaff: "Personale Scolastico Autorizzato",
      logAdministrationBtn: "Registra Somministrazione",
      inailAlert48h: "Termine perentorio INAIL: 48 ore dalla ricezione del certificato medico",
      accidentTime: "Data e Ora Infortunio",
      location: "Luogo Infortunio",
      activity: "Attività svolta",
      dynamics: "Dinamica dell'evento",
      injuries: "Lesioni riportate",
      ambulance: "Chiamata 118 / Ambulanza",
      prognosis: "Giorni Prognosi",
      transmitInail: "Registra Trasmissione INAIL",
      inailProtocol: "Protocollo INAIL"
    },
    competenciesDM14: {
      title: "Certificazione delle Competenze (D.M. 14/2024)",
      subtitle: "Modelli nazionali unificati Primaria e I Ciclo allineati a DigComp 2.2",
      tierPrimaria: "Scuola Primaria (Termine V Anno)",
      tierPrimoCiclo: "Scuola Secondaria I Grado (Termine III Anno)",
      rubricTitle: "Descrittori Livelli di Competenza",
      levelAvanzato: "Avanzato",
      levelIntermedio: "Intermedio",
      levelBase: "Base",
      levelIniziale: "Iniziale",
      considerations: "Considerazioni Finali del Consiglio",
      deliberationDate: "Data Delibera Consiglio",
      signByPrincipal: "Firma Digitale Dirigente Scolastico",
      signed: "Firmata Digitalmente dal Dirigente",
      invalsiTitle: "Analisi Comparativa Prove INVALSI",
      invalsiSubtitle: "Correlazione tra esiti INVALSI (WLE / Livelli 1-5) e medie reali del consiglio per RAV e PdM",
      scoreWLE: "Punteggio WLE",
      level: "Livello",
      nationalMean: "Media Nazionale",
      reportGrade: "Voto Registro",
      delta: "Scostamento (Delta)",
      importInvalsi: "Importa Risultati INVALSI"
    },
    securityPasses: {
      title: "Pass Fuori Classe & Gestione Deleghe",
      subtitle: "Tracciamento uscite dall'aula in tempo reale e registro deleghe al ritiro alunni",
      hallPassesTab: "Pass Digitale Fuori Classe",
      delegatesTab: "Deleghe Ritiro & Sentenze Affido",
      issuePassBtn: "Rilascia Pass Uscita",
      passType: "Motivazione Uscita",
      typeBathroom: "Servizi Igienici (Bagno)",
      typeInfirmary: "Infermeria / Primo Soccorso",
      typeOffice: "Segreteria / Dirigenza",
      typeStaff: "Collaboratori ATA",
      departedAt: "Ora Uscita",
      elapsedMinutes: "Minuti Fuori Aula",
      anomalyAlert: "Alert: Tempo Fuori Aula Superiore a 15 Minuti!",
      markReturned: "Registra Rientro in Aula",
      addDelegateBtn: "Aggiungi Persona Delegata",
      delegateName: "Nome Delegato",
      relation: "Grado di Parentela",
      phone: "Telefono",
      fiscalCode: "Codice Fiscale",
      idCard: "Documento di Riconoscimento",
      permanent: "Delega Permanente per l'A.S.",
      custodyNotice: "Atti di Affido Monogenitoriale e Sentenze ex art. 337 c.c. con oscuramento dati applicato.",
      revoke: "Revoca Delega"
    },
    homeworkSubmissions: {
      title: "Consegna Elaborati & Correzione con Canvas",
      subtitle: "Caricamento compiti studente, annotazione grafica PDF e sync Google Classroom",
      uploadTitle: "Carica Elaborato del Compito",
      selectFile: "Seleziona PDF o Immagine",
      submitBtn: "Invia Elaborato al Docente",
      statusSubmitted: "Consegnato",
      statusReturned: "Corretto e Restituito",
      canvasTitle: "Strumento di Correzione Grafica su PDF",
      penTool: "Penna Grafica",
      textTool: "Testo / Nota",
      clearCanvas: "Cancella Tratti",
      grade: "Voto (1-10)",
      feedback: "Note e Feedback Didattico",
      returnBtn: "Restituisci Elaborato Corretto",
      syncClassroom: "Sincronizza Voto con Google Classroom",
      classroomSynced: "Sincronizzato con Classroom"
    },
    pnrrProjects: {
      title: "Monitoraggio Progetti PNRR & Piattaforma Futura",
      subtitle: "Rendicontazione target 75% frequenza minima e registri orari digitali",
      newProjectBtn: "Nuovo Progetto PNRR",
      cup: "Codice CUP",
      noticeType: "Avviso / Linea di Investimento",
      totalHours: "Ore Totali Previste",
      targetStudents: "Target Corsisti",
      enrolledStudents: "Corsisti Iscritti",
      target75Passed: "Target 75% Raggiunto",
      target75Notice: "Requisito di Rendicontazione: Almeno il 75% delle ore deve essere frequentato per la convalida del target.",
      recordLessonBtn: "Firma e Registra Lezione PNRR",
      exportFuturaBtn: "Esporta Dati per Piattaforma Futura",
      signatures: "Firme Presenze Corsisti"
    },
    infanzia: {
      title: "Registro Scuola dell'Infanzia",
      subtitle: "Gestione sezioni colorate, presenze con mensa, campi di esperienza e diario di bordo",
      sectionsTab: "Sezioni dell'Infanzia",
      attendancesTab: "Appello & Pasti Mattutini",
      campiTab: "Osservazioni Campi di Esperienza",
      diarioTab: "Diario di Bordo Fotografico",
      newSectionBtn: "Nuova Sezione",
      sectionName: "Nome Sezione (es. Farfalle, Rossi)",
      color: "Colore Distintivo",
      present: "Presente",
      bookMeal: "Pasto Mensa",
      specialDiet: "Dieta Speciale",
      campoEsperienza: "Campo di Esperienza",
      levelDescriptor: "Livello di Sviluppo",
      storyTitle: "Titolo Attività del Giorno",
      narrativeStory: "Racconto delle Attività per le Famiglie",
      publishStory: "Pubblica nel Diario delle Famiglie"
    }
  },
  'en-US': {
    teacherHours: {
      title: "Teacher Working Hours Management (Art. 29 CCNL)",
      subtitle: "Annual 40+40 hours accounting, overtime tracking and short leave management",
      collegioCard: "Collegial Functional Activities (Max 40h)",
      collegioDesc: "Faculty meetings, department conferences, and parent communication",
      consigliCard: "Class Councils & Inter-class Meetings (Max 40h)",
      consigliDesc: "Class councils, periodic and final scrutinii",
      planned: "Planned",
      done: "Completed",
      remaining: "Remaining",
      overtimeAlert: "Alert: contractual hour cap exceeded",
      deficitAlert: "Hour deficit to complete before end of school year",
      bancaOre: "Hour Bank & Make-up Hours",
      disposizione: "Available On-call Hours",
      logActivityTitle: "Log Completed Activity",
      category: "Activity Category",
      hours: "Hours",
      date: "Date",
      description: "Activity description",
      shortLeavesTitle: "Short Leaves (Max 2h - CCNL)",
      requestLeaveBtn: "Request Short Leave",
      shortLeaveNotice: "Under Art. 16 CCNL, short leaves cannot exceed 2 hours per day and must be recovered within 2 months.",
      deadline: "Recovery Deadline",
      status: "Status",
      recoveredHours: "Recovered Hours",
      recordRecovery: "Record Make-up Hours"
    },
    healthSafety: {
      title: "School Health & INAIL Accident Registry",
      subtitle: "Life-saving medication protocols and mandatory 48h accident reporting",
      medicationsTab: "Medication Administration at School",
      accidentsTab: "INAIL Accident Registry",
      registerMedBtn: "Register Medication Therapeutic Plan",
      student: "Student",
      medName: "Medication Name",
      dosage: "Dosage",
      adminMode: "Administration Route",
      prescriptionExpiry: "Medical Prescription Expiry",
      medExpiry: "Medication Expiry Date",
      emergencyProtocol: "Emergency Protocol",
      authorizedStaff: "Authorized Trained Staff",
      logAdministrationBtn: "Record Administration",
      inailAlert48h: "INAIL 48-hour mandatory deadline upon receipt of medical certificate",
      accidentTime: "Accident Date & Time",
      location: "Accident Location",
      activity: "Ongoing Activity",
      dynamics: "Accident Dynamics",
      injuries: "Reported Injuries",
      ambulance: "Emergency 118 / Ambulance Called",
      prognosis: "Prognosis Days",
      transmitInail: "Record INAIL Transmission",
      inailProtocol: "INAIL Protocol Code"
    },
    competenciesDM14: {
      title: "Competency Certification (D.M. 14/2024)",
      subtitle: "Unified national templates for Primary and Lower Secondary aligned with DigComp 2.2",
      tierPrimaria: "Primary School (End of 5th Year)",
      tierPrimoCiclo: "Lower Secondary School (End of 3rd Year)",
      rubricTitle: "Competency Level Descriptors",
      levelAvanzato: "Advanced",
      levelIntermedio: "Intermediate",
      levelBase: "Basic",
      levelIniziale: "Initial",
      considerations: "Final Board Observations",
      deliberationDate: "Board Deliberation Date",
      signByPrincipal: "Digital Signature by Principal",
      signed: "Digitally Signed by Principal",
      invalsiTitle: "INVALSI Benchmark & Comparative Analytics",
      invalsiSubtitle: "Correlation between INVALSI results (WLE / Levels 1-5) and school grades for RAV and PdM",
      scoreWLE: "WLE Score",
      level: "Level",
      nationalMean: "National Average",
      reportGrade: "Register Grade",
      delta: "Deviation (Delta)",
      importInvalsi: "Import INVALSI Results"
    },
    securityPasses: {
      title: "Hall Passes & Pickup Delegates",
      subtitle: "Real-time out-of-classroom tracking and authorized student pickup registry",
      hallPassesTab: "Digital Hall Pass",
      delegatesTab: "Pickup Delegates & Legal Custody",
      issuePassBtn: "Issue Hall Pass",
      passType: "Reason for Exit",
      typeBathroom: "Restroom",
      typeInfirmary: "Infirmary / First Aid",
      typeOffice: "Administration / Principal",
      typeStaff: "Janitorial Staff",
      departedAt: "Departure Time",
      elapsedMinutes: "Minutes Outside",
      anomalyAlert: "Alert: Student outside classroom for more than 15 minutes!",
      markReturned: "Mark Returned to Class",
      addDelegateBtn: "Add Authorized Delegate",
      delegateName: "Delegate Name",
      relation: "Relationship",
      phone: "Phone Number",
      fiscalCode: "Tax Code",
      idCard: "ID Document",
      permanent: "Permanent Delegate for School Year",
      custodyNotice: "Sole custody court orders under Art. 337 c.c. applied with data blinding.",
      revoke: "Revoke Delegate"
    },
    homeworkSubmissions: {
      title: "Homework Submissions & Canvas Grading",
      subtitle: "Student assignment uploads, digital PDF annotation and Google Classroom sync",
      uploadTitle: "Upload Homework File",
      selectFile: "Select PDF or Image",
      submitBtn: "Submit to Teacher",
      statusSubmitted: "Submitted",
      statusReturned: "Graded and Returned",
      canvasTitle: "Digital PDF Annotation Canvas",
      penTool: "Pen Tool",
      textTool: "Text / Note Tool",
      clearCanvas: "Clear Strokes",
      grade: "Grade (1-10)",
      feedback: "Feedback & Comments",
      returnBtn: "Return Graded Work",
      syncClassroom: "Sync Grade with Google Classroom",
      classroomSynced: "Synced with Classroom"
    },
    pnrrProjects: {
      title: "PNRR Projects & Futura MIM Monitoring",
      subtitle: "Mandatory 75% attendance threshold reporting and digital lesson logbooks",
      newProjectBtn: "New PNRR Project",
      cup: "CUP Code",
      noticeType: "Funding Notice Line",
      totalHours: "Total Planned Hours",
      targetStudents: "Target Students",
      enrolledStudents: "Enrolled Students",
      target75Passed: "75% Target Achieved",
      target75Notice: "Validation Requirement: At least 75% attendance is strictly required for financial target validation.",
      recordLessonBtn: "Sign & Record PNRR Session",
      exportFuturaBtn: "Export for Piattaforma Futura",
      signatures: "Student Attendance Signatures"
    },
    infanzia: {
      title: "Kindergarten Register (Scuola dell'Infanzia)",
      subtitle: "Color-coded sections, morning meal attendance, experience fields and photo logbook",
      sectionsTab: "Kindergarten Sections",
      attendancesTab: "Morning Roll Call & Meals",
      campiTab: "Experience Fields Observations",
      diarioTab: "Photo Daily Logbook",
      newSectionBtn: "New Section",
      sectionName: "Section Name (e.g. Butterflies, Red)",
      color: "Distinguishing Color",
      present: "Present",
      bookMeal: "Canteen Meal",
      specialDiet: "Special Diet",
      campoEsperienza: "Experience Field",
      levelDescriptor: "Development Level",
      storyTitle: "Activity Title of the Day",
      narrativeStory: "Narrative Activity Story for Parents",
      publishStory: "Publish in Family Logbook"
    }
  }
};

// Generate foreign language templates from English base
const langCodes = ['de-DE', 'fr-FR', 'es-ES', 'ro-RO', 'sq-AL', 'ru-RU', 'uk-UA', 'ar-SA', 'zh-CN'];

for (const code of langCodes) {
  translations[code] = JSON.parse(JSON.stringify(translations['en-US']));
}

// Write to each locale file
for (const [locale, data] of Object.entries(translations)) {
  const filePath = path.join(i18nDir, locale, 'index.js');
  if (!fs.existsSync(filePath)) {
    console.warn(`File not found: ${filePath}`);
    continue;
  }

  let content = fs.readFileSync(filePath, 'utf-8');
  // Check if teacherHours already added
  if (content.includes('"teacherHours":') || content.includes('teacherHours:')) {
    console.log(`Already contains teacherHours in ${locale}, skipping.`);
    continue;
  }

  // Find the last closing brace
  const lastBraceIndex = content.lastIndexOf('}');
  if (lastBraceIndex === -1) {
    console.error(`Invalid structure in ${filePath}`);
    continue;
  }

  // Format keys as JSON without outer braces
  const formattedData = JSON.stringify(data, null, 2);
  const dataInner = formattedData.slice(formattedData.indexOf('\n') + 1, formattedData.lastIndexOf('\n'));

  // Insert before last closing brace
  let updatedContent = content.slice(0, lastBraceIndex).trimEnd();
  if (!updatedContent.endsWith(',')) {
    updatedContent += ',';
  }
  updatedContent += '\n' + dataInner + '\n}\n';

  fs.writeFileSync(filePath, updatedContent, 'utf-8');
  console.log(`Updated i18n for ${locale}`);
}
