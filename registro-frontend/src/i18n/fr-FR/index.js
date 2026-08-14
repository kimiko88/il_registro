export default {
  common: {
    welcome: 'Bienvenue',
    back: 'Retour',
    save: 'Enregistrer',
    cancel: 'Annuler',
    close: 'Fermer',
    delete: 'Supprimer',
    edit: 'Modifier',
    search: 'Rechercher (Ctrl+K)...',
    filter: 'Filtrer',
    actions: 'Actions',
    loading: 'Chargement en cours...',
    success: 'Opération réussie',
    error: 'Une erreur est survenue',
    systemSettings: 'Paramètres du système',
    language: 'Langue / Language',
    security: 'Sécurité & Authentification',
    notifications: 'Notifications & Alertes',
    mainMenu: 'MENU PRINCIPAL',
    logout: 'Se déconnecter'
  },
  login: {
    welcomeBack: 'Bon retour',
    subtitle: 'Connectez-vous pour accéder au registre électronique',
    emailLabel: 'Adresse E-mail',
    emailRequired: 'L\'e-mail est obligatoire',
    passwordLabel: 'Mot de passe',
    passwordRequired: 'Le mot de passe est obligatoire',
    rememberMe: 'Se souvenir de moi',
    submit: 'Connexion',
    noAccount: 'Vous n\'avez pas de compte ?',
    contactSecretary: 'Contacter le Secrétariat',
    contactTitle: 'Contacter le Secrétariat',
    contactSubtitle: 'Sélectionnez votre établissement dans le menu déroulant pour afficher les coordonnées e-mail et téléphoniques du secrétariat.',
    selectSchool: 'Sélectionnez votre école / établissement',
    noSchoolFound: 'Aucune école trouvée',
    emailSegreteria: 'E-mail Secrétariat',
    phone: 'Téléphone',
    sendEmail: 'Envoyer un e-mail',
    copyEmail: 'Copier l\'e-mail',
    emailCopied: 'Adresse e-mail copiée dans le presse-papiers !',
    chooseSchoolPrompt: 'Choisissez un établissement dans le menu ci-dessus pour afficher les coordonnées.',
    sessionExpired: 'Session expirée. Veuillez vous reconnecter.',
    showPassword: 'Afficher le mot de passe',
    hidePassword: 'Masquer le mot de passe'
  },
  nav: {
    dashboard: 'Tableau de bord',
    schools: 'Gestion des écoles',
    mySchool: 'Mon école',
    users: 'Gestion des utilisateurs',
    admins: 'Gestion des administrateurs',
    monitoring: 'Surveillance système',
    analytics: 'Analyses globales',
    auditLogs: 'Journaux d\'audit',
    settings: 'Paramètres',
    support: 'Support & Assistance',
    featureFlags: 'Options & Établissement',
    elearning: 'E-Learning Google & Teams',
    students: 'Élèves',
    classes: 'Classes',
    groups: 'Groupes linguistiques',
    documents: 'Documents',
    certificates: 'Certificats',
    textbooks: 'Manuels scolaires',
    meetings: 'Réunions',
    communications: 'Communications',
    reports: 'Rapports',
    pcto: 'Stage / PCTO',
    scrutiny: 'Conseil de classe / Notes finales',
    myClasses: 'Mes classes',
    classRegister: 'Registre de classe',
    uda: 'Planification UdA',
    competencies: 'Évaluation des compétences',
    grades: 'Notes',
    attendance: 'Présences',
    didactics: 'Matériel didactique',
    pdp: 'Plans PDP / PEI',
    rubrics: 'Grilles d\'évaluation',
    coordination: 'Coordination',
    timetable: 'Emploi du temps',
    agenda: 'Cahier de texte / Agenda',
    colloqui: 'Rendez-vous parents',
    substitutions: 'Remplacements',
    verbali: 'Procès-verbaux',
    notes: 'Remarques disciplinaires',
    myGrades: 'Mes notes',
    myAttendance: 'Mes présences',
    homework: 'Devoirs',
    orientamento: 'Orientation',
    calendar: 'Calendrier scolaire',
    reportCard: 'Bulletin scolaire',
    profile: 'Profil',
    myChildren: 'Mes enfants',
    goals: 'Objectifs',
    trips: 'Sorties & Voyages',
    payments: 'Paiements',
    assemblies: 'Assemblées & Réunions'
  },
  categories: {
    anagraficheClassi: 'Registres & Classes',
    attiCertificati: 'Actes & Certificats',
    serviziReport: 'Services & Rapports',
    didatticaValutazione: 'Enseignement & Évaluation',
    organizzazioneOrario: 'Organisation & Emploi du temps',
    comunicazioniAtti: 'Communications & Actes',
    percorsiComunicazioni: 'Parcours & Communications',
    valutazioneDidattica: 'Évaluation & Enseignement',
    serviziOrari: 'Services & Horaires',
    comunicazioniAccount: 'Communications & Compte'
  },
  roles: {
    admin: 'Administrateur',
    superadmin: 'Super Administrateur',
    secretary: 'Secrétaire',
    teacher: 'Enseignant',
    student: 'Élève',
    parent: 'Parent',
    user: 'Utilisateur'
  },
  notifications: {
    title: 'Notifications et communications',
    logoutSuccess: 'Déconnexion réussie',
    logoutError: 'Erreur lors de la déconnexion',
    settingsSaved: 'Paramètres enregistrés avec succès !',
    emailCopied: 'Adresse e-mail copiée !',
    passwordUpdated: 'Mot de passe mis à jour avec succès !',
    languageChanged: 'Langue mise à jour avec succès'
  },
  errors: {
    connectionError: 'Erreur de connexion au serveur. Vérifiez votre connexion.',
    forbidden: 'Vous n\'avez pas les permissions nécessaires.',
    serverError: 'Une erreur s\'est produite sur le serveur.',
    unauthorized: 'Session invalide ou expirée.',
    invalidCredentials: 'E-mail ou mot de passe incorrect.',
    rateLimit: 'Trop de tentatives. Veuillez réessayer dans quelques minutes.',
    accountDisabled: 'Compte désactivé ou suspendu. Veuillez contacter le secrétariat.',
    userNotFound: 'Utilisateur non trouvé.',
    passwordMismatch: 'Les mots de passe ne correspondent pas.',
    sessionInvalid: 'Session invalide, veuillez vous reconnecter.',
    ERR_CURRENT_PASSWORD_INCORRECT: 'Le mot de passe actuel est incorrect.',
    ERR_PASSWORD_COMPLEXITY: 'Le mot de passe doit contenir au moins une majuscule, une minuscule, un chiffre et un caractère spécial.',
    ERR_PASSWORD_TOO_SHORT: 'Le mot de passe doit comporter au moins 10 caractères.',
    ERR_PASSWORD_TOO_LONG: 'Le mot de passe est trop long (maximum 128 caractères).',
    ERR_PASSWORD_RECENTLY_USED: 'Le nouveau mot de passe ne doit pas correspondre aux 5 derniers mots de passe utilisés.',
    ERR_REQUIRED_FIELDS: 'Veuillez remplir tous les champs obligatoires.'
  },
  settings: {
    title: 'Paramètres du système',
    subtitle: 'Configurez les préférences générales, de sécurité et de notifications',
    languageLabel: 'Langue / Language',
    securityLabel: 'Sécurité & Authentification',
    securitySub: 'Politique de mot de passe, 2FA/MFA et délai d\'expiration de session',
    notificationsLabel: 'Notifications & Alertes',
    notificationsSub: 'Préférences e-mail, push et alertes de remplacement',
    mfaActive: 'MFA Actif',
    mfaOptional: 'MFA Optionnel',
    mfaTitle: 'Obligation 2FA / MFA',
    mfaCaption: 'Exiger l\'authentification à deux facteurs avec TOTP pour tout le personnel.',
    minPasswordLength: 'Longueur minimale du mot de passe',
    sessionTimeout: 'Délai d\'inactivité de session',
    maxLoginAttempts: 'Tentatives échouées avant blocage'
  },
  search: {
    placeholder: 'Rechercher élèves, enseignants, secrétariat, classes, circulaires, menus…',
    hint: 'Tapez pour rechercher des élèves, enseignants, secrétariat, classes ou menus',
    noResults: 'Aucun résultat pour',
    navigate: 'Naviguer',
    open: 'Ouvrir',
    close: 'Fermer'
  }
,
  onboarding: {
    welcomeTitle: "Bienvenue dans le Registre Électronique!",
    welcomeSubtitle: "Découvrez en quelques étapes comment utiliser au mieux toutes les fonctionnalités disponibles pour votre rôle.",
    startTour: "Démarrer la visite",
    skipTour: "Ignorer la visite",
    next: "Suivant",
    prev: "Précédent",
    finish: "Commencer à utiliser le Registre",
    stepOf: "Étape {current} sur {total}",
    restartTour: "Redémarrer la visite",
    tourCompleted: "Visite terminée!",
    tourCompletedMsg: "Vous êtes prêt à utiliser le Registre Électronique. Vous pouvez revoir le guide à tout moment.",
    teacher: {
      step1_title: "Tableau de bord Professeur",
      step1_desc: "Votre tableau de bord personnel affiche un résumé des cours du jour, les notifications récentes et un accès rapide aux fonctions les plus utilisées.",
      step2_title: "Registre de classe & Présences",
      step2_desc: "Dans \"Mes Classes\" vous pouvez accéder au registre, enregistrer les présences et absences, saisir les sujets de cours.",
      step3_title: "Gestion des notes",
      step3_desc: "Dans \"Notes\" vous pouvez saisir des notes orales et écrites, voir les moyennes de classe et la distribution des notes.",
      step4_title: "Agenda & Communications",
      step4_desc: "L'agenda permet de planifier des contrôles, devoirs et activités. Dans \"Communications\" vous pouvez envoyer des messages aux élèves et parents.",
      step5_title: "Réunions parents",
      step5_desc: "Gérez les rendez-vous individuels avec les parents : consultez les réservations, les disponibilités et la progression de la classe.",
      step6_title: "Paramètres & Profil",
      step6_desc: "Dans \"Paramètres\" vous pouvez personnaliser la langue, le thème, les notifications et mettre à jour votre mot de passe."
    },
    student: {
      step1_title: "Votre tableau de bord",
      step1_desc: "Le tableau de bord personnel affiche les devoirs à venir, les dernières notifications et un résumé de vos résultats scolaires.",
      step2_title: "Mes notes",
      step2_desc: "Dans \"Mes Notes\" vous pouvez voir toutes les notes enregistrées par les professeurs, les moyennes par matière et l'évolution dans le temps.",
      step3_title: "Mes présences",
      step3_desc: "Suivez vos présences, absences et retards. Vous pouvez consulter les détails journaliers et mensuels.",
      step4_title: "Devoirs & Didactique",
      step4_desc: "Dans \"Devoirs\" vous trouvez tous les devoirs assignés avec leurs échéances. \"Didactique\" affiche les documents déposés par les professeurs.",
      step5_title: "Bulletin & Documents",
      step5_desc: "Dans \"Bulletin\" vous pouvez consulter votre document d'évaluation. Dans \"Documents\" vous trouvez circulaires et matériel scolaire.",
      step6_title: "Calendrier scolaire",
      step6_desc: "Le calendrier scolaire affiche les vacances, les dates des contrôles prévus et les événements importants."
    },
    parent: {
      step1_title: "Tableau de bord Parent",
      step1_desc: "Votre tableau de bord affiche un résumé des progrès de vos enfants, les dernières notifications scolaires et les messages non lus.",
      step2_title: "Mes enfants",
      step2_desc: "Dans \"Mes Enfants\" vous trouvez la liste de vos enfants inscrits. Sélectionnez-en un pour voir le profil scolaire complet.",
      step3_title: "Notes & Présences",
      step3_desc: "Suivez les notes et présences de vos enfants en temps réel. Recevez des notifications immédiates pour les absences et nouvelles notes.",
      step4_title: "Communications scolaires",
      step4_desc: "Toutes les communications officielles de l'école (circulaires, avis, messages des professeurs) sont rassemblées ici.",
      step5_title: "Rendez-vous professeurs",
      step5_desc: "Réservez des rendez-vous individuels avec les professeurs directement depuis l'application.",
      step6_title: "Paiements & Documents",
      step6_desc: "Gérez les paiements scolaires et accédez aux documents de votre enfant (bulletins, certificats)."
    },
    secretary: {
      step1_title: "Tableau de bord Secrétariat",
      step1_desc: "Le tableau de bord affiche les activités en attente, les dernières demandes et les statistiques principales de l'établissement.",
      step2_title: "Gestion classes & élèves",
      step2_desc: "Dans \"Classes\" vous gérez toutes les classes. Dans \"Élèves\" vous trouvez le registre complet avec recherche avancée.",
      step3_title: "Certificats & Documents",
      step3_desc: "Générez et imprimez des certificats de scolarité et de présence. Gérez l'archive documentaire numérique.",
      step4_title: "Gestion des horaires",
      step4_desc: "Configurez l'emploi du temps scolaire, gérez les remplacements et planifiez les activités.",
      step5_title: "Rapports & Statistiques",
      step5_desc: "Générez des rapports personnalisés sur les présences, notes et inscriptions.",
      step6_title: "Utilisateurs & Communications",
      step6_desc: "Gérez les comptes enseignants, élèves et parents. Envoyez des communications officielles."
    },
    admin: {
      step1_title: "Tableau de bord Administrateur",
      step1_desc: "Le tableau de bord affiche l'état du système, les dernières activités et les métriques principales d'utilisation.",
      step2_title: "Surveillance système",
      step2_desc: "Surveillez les performances du système en temps réel, les sessions actives et les logs.",
      step3_title: "Gestion utilisateurs & établissements",
      step3_desc: "Gérez tous les comptes utilisateurs, créez de nouveaux établissements, assignez les rôles.",
      step4_title: "Analytics & Rapports",
      step4_desc: "Visualisez les analyses globales d'utilisation de la plateforme et générez des rapports détaillés.",
      step5_title: "Journal d'audit",
      step5_desc: "Accédez au registre complet de toutes les opérations effectuées dans le système.",
      step6_title: "Paramètres système",
      step6_desc: "Configurez les politiques de sécurité, l'authentification à deux facteurs et les intégrations e-learning."
    }
  },
  help: {
    title: "Centre d'aide",
    subtitle: "Guides et tutoriels pour bien utiliser le Registre Électronique",
    searchPlaceholder: "Rechercher dans le guide...",
    noResults: "Aucun article trouvé pour",
    categories: "Catégories",
    allTopics: "Tous les sujets",
    restartTour: "Redémarrer la visite guidée",
    openHelp: "Ouvrir le guide",
    needHelp: "Besoin d'aide?",
    contactSupport: "Contacter le support",
    fabTooltip: "Aide & Guide",
    new: "Nouveau"
  }

,
  onboardingExtra: {
    "openGuide": "Open Full Guide",
    "viewAllFeatures": "Discover all features",
    "completionTitle": "You're ready! 🎉",
    "completionDesc": "You've completed the tour. You can access the full guide at any time from the ? button in the top bar.",
    "teacher": {
        "step7_title": "Digital Register & Signatures",
        "step7_desc": "Digitally sign each lesson with a single click. The register automatically tracks lesson topics, teaching hours, and private notes.",
        "step7_bullets": [
            "Digital lesson signature",
            "Lesson topics",
            "Private notes",
            "Hours log"
        ],
        "step8_title": "Competencies & Learning Units",
        "step8_desc": "Assess student competencies according to EU descriptors, manage Learning Units (UDA) and annual work plans.",
        "step8_bullets": [
            "EU Competencies",
            "Learning Units (UDA)",
            "Annual work plan",
            "Evaluation rubrics"
        ],
        "step1_bullets": [
            "Today's lessons",
            "Recent notifications",
            "Quick shortcuts",
            "Class overview"
        ],
        "step2_bullets": [
            "Attendance register",
            "Lesson topics",
            "Digital signature",
            "Absence management"
        ],
        "step3_bullets": [
            "Oral & written grades",
            "Class average",
            "Distribution chart",
            "Export tables"
        ],
        "step4_bullets": [
            "Homework & tests",
            "Activity calendar",
            "Class announcements",
            "Parent notices"
        ],
        "step5_bullets": [
            "Meeting availability",
            "Parent bookings",
            "Video conferences",
            "Meeting history"
        ],
        "step6_bullets": [
            "Theme & language",
            "Push notifications",
            "Password management",
            "Public profile"
        ]
    },
    "student": {
        "step7_title": "Internships & Competencies",
        "step7_desc": "Track your work-based learning hours (PCTO), view competence certifications, and build your digital portfolio.",
        "step7_bullets": [
            "Internship hours",
            "Certifications",
            "Digital portfolio",
            "Tutor feedback"
        ],
        "step8_title": "Announcements & Noticeboard",
        "step8_desc": "Read official school circulars, answer surveys, and view personalized announcements from your teachers.",
        "step8_bullets": [
            "School circulars",
            "Surveys & answers",
            "Personal alerts",
            "Teacher notices"
        ],
        "step1_bullets": [
            "Upcoming deadlines",
            "Latest notifications",
            "GPA summary",
            "Upcoming tests"
        ],
        "step2_bullets": [
            "Grades by subject",
            "Live average",
            "Trend chart",
            "Term comparison"
        ],
        "step3_bullets": [
            "Monthly attendance",
            "Absence counter",
            "Late arrivals & early leaves",
            "Excuses & justifications"
        ],
        "step4_bullets": [
            "Assigned homework",
            "Due dates",
            "Learning materials",
            "Links & resources"
        ],
        "step5_bullets": [
            "Digital report card",
            "PDF download",
            "Certificates",
            "Forms"
        ],
        "step6_bullets": [
            "School holidays",
            "Test dates",
            "Extracurriculars",
            "Field trips & events"
        ]
    },
    "parent": {
        "step7_title": "Absence Justifications & Permissions",
        "step7_desc": "Submit digital justifications for student absences directly from the app, and approve early leave authorizations.",
        "step7_bullets": [
            "Online justifications",
            "Early leave permits",
            "Special activities",
            "Absence history"
        ],
        "step8_title": "Academic Analytics & Progress",
        "step8_desc": "View academic performance charts over time, compare subject averages, and track improvement trends.",
        "step8_bullets": [
            "Progress chart",
            "Subject comparison",
            "Historical trends",
            "Goals achieved"
        ],
        "step1_bullets": [
            "Children overview",
            "Recent alerts",
            "Unread messages",
            "Appointments"
        ],
        "step2_bullets": [
            "Academic profile",
            "Student documents",
            "Teacher contacts",
            "Enrollment info"
        ],
        "step3_bullets": [
            "Real-time grades",
            "Subject GPA",
            "Daily attendance",
            "Automated alerts"
        ],
        "step4_bullets": [
            "School circulars",
            "Teacher notices",
            "Urgent alerts",
            "Digital bulletin board"
        ],
        "step5_bullets": [
            "Book parent meeting",
            "Choose time slot",
            "Email confirmation",
            "Reschedule"
        ],
        "step6_bullets": [
            "Tuition & fees",
            "Payment history",
            "Download receipts",
            "Request documents"
        ]
    },
    "secretary": {
        "step7_title": "Enrollment & Student Registry",
        "step7_desc": "Manage yearly school enrollments, maintain student demographic records, and handle transfers and special education plans.",
        "step7_bullets": [
            "Yearly enrollment",
            "Transfers in/out",
            "Special needs / DSA",
            "Registry archives"
        ],
        "step8_title": "Circulars & School Broadcasts",
        "step8_desc": "Compose and publish digital circulars, manage noticeboards, and send targeted group broadcasts with signature tracking.",
        "step8_bullets": [
            "Compose circulars",
            "Digital signature",
            "Targeted recipients",
            "Publication history"
        ],
        "step1_bullets": [
            "Pending tasks",
            "Incoming requests",
            "School statistics",
            "System alerts"
        ],
        "step2_bullets": [
            "Class list",
            "Teacher assignments",
            "Weekly timetable",
            "Revision history"
        ],
        "step3_bullets": [
            "Enrollment certificates",
            "Attendance records",
            "PDF export",
            "Digital seal"
        ],
        "step4_bullets": [
            "School schedule",
            "Teacher substitutions",
            "Extracurriculars",
            "Staff notices"
        ],
        "step5_bullets": [
            "Attendance reports",
            "Official exports",
            "Statistical charts",
            "Advanced filters"
        ],
        "step6_bullets": [
            "Create accounts",
            "Assign roles",
            "Reset passwords",
            "Access control"
        ]
    },
    "admin": {
        "step7_title": "Schools & Campuses Management",
        "step7_desc": "Create and configure multiple school branches, manage tenant settings, assign campus administrators, and monitor activity.",
        "step7_bullets": [
            "New campuses",
            "Per-school config",
            "Branch admins",
            "Multi-tenant"
        ],
        "step8_title": "Integrations & API Management",
        "step8_desc": "Connect third-party systems (Google Workspace, Microsoft 365, LMS), configure webhooks, and monitor API traffic.",
        "step8_bullets": [
            "Google Workspace",
            "Microsoft 365",
            "Webhook config",
            "API monitoring"
        ],
        "step1_bullets": [
            "System health",
            "Active users",
            "Open sessions",
            "Recent errors"
        ],
        "step2_bullets": [
            "CPU & memory usage",
            "API latency",
            "Live logs",
            "Threshold alerts"
        ],
        "step3_bullets": [
            "Manage users",
            "RBAC roles",
            "Manage schools",
            "Bulk CSV import"
        ],
        "step4_bullets": [
            "Platform metrics",
            "Periodic reports",
            "User growth",
            "Data exports"
        ],
        "step5_bullets": [
            "Immutable audit logs",
            "Security filters",
            "Compliance export",
            "GDPR tools"
        ],
        "step6_bullets": [
            "Password policy",
            "Mandatory 2FA",
            "Feature flags",
            "Maintenance mode"
        ]
    }
},
  guideCenter: {
    "title": "Help & Knowledge Center",
    "subtitle": "Comprehensive guides for every section of the electronic register",
    "search": "Search in guides...",
    "noResults": "No guides found for",
    "readingTime": "min read",
    "step": "Step",
    "tip": "Tip",
    "warning": "Warning",
    "shortcut": "Shortcut",
    "teacher": {
        "dashboard": {
            "title": "Teacher Dashboard",
            "desc": "Overview of daily lessons, notifications, and quick actions.",
            "content": "The teacher dashboard provides a complete overview of your day. Top widgets display today’s scheduled classes with times and room numbers.\n\nStep 1: Customize your dashboard widgets by dragging them.\nStep 2: Click on any class in the timetable to open its register directly.\nStep 3: Filter notifications by urgency or category.\n\nTip: Enable priority alerts to stay informed about urgent announcements."
        },
        "attendance": {
            "title": "Attendance Register",
            "desc": "How to record attendance, absences, and delays quickly and accurately.",
            "content": "The attendance register automatically loads the student list for the current period.\n\nStep 1: Click on a student name to toggle status (Present, Absent, Late, Early exit).\nStep 2: Add optional justification notes or remarks.\nStep 3: Click \"Sign Lesson\" to confirm and digitally sign the session.\nStep 4: Parents of absent students are notified automatically.\n\nShortcut: Press Ctrl+Enter to quickly save attendance."
        },
        "grades": {
            "title": "Grade Management",
            "desc": "Enter and manage written, oral, and practical grades.",
            "content": "The gradebook is an interactive grid organized by students and evaluation dates.\n\nStep 1: Select class and subject from the dropdown.\nStep 2: Click on any cell to input or edit a grade.\nStep 3: Choose evaluation type (Written, Oral, Practical) and score.\nStep 4: Add private or public feedback comments.\nStep 5: Subject and class averages update instantly.\n\nWarning: Published grades are visible to students and parents within seconds."
        },
        "agenda": {
            "title": "School Agenda & Homework",
            "desc": "Plan tests, assignments, and shared class activities.",
            "content": "The teacher agenda syncs automatically with students and parents.\n\nStep 1: Click on any calendar date to schedule an event.\nStep 2: Choose event type: Homework, Exam, Field trip, or Reminder.\nStep 3: Select subject, description, and attach learning files.\nStep 4: Check \"Notify students\" to broadcast an instant push notification.\n\nTip: Use the Week view to check for test overlaps across subjects."
        },
        "communications": {
            "title": "Communications & Messages",
            "desc": "Send structured messages to classes, parents, and colleagues.",
            "content": "Send announcements through dedicated channels: direct messages, class circulars, or urgent bulletins.\n\nStep 1: Navigate to \"Communications\" → \"New Message\".\nStep 2: Choose channel type and recipient list.\nStep 3: Format your message using the rich-text editor.\nStep 4: Send immediately or schedule for future publication.\n\nWarning: Urgent announcements trigger push notifications and SMS alerts."
        },
        "meetings": {
            "title": "Parent-Teacher Meetings",
            "desc": "Configure availability slots and manage booking appointments.",
            "content": "Manage weekly office hours and individual parent conferences.\n\nStep 1: Go to \"Meetings\" → \"Set Availability\".\nStep 2: Define weekly recurring time slots and appointment durations.\nStep 3: Parents book online; you receive automated calendar invites.\nStep 4: Review student performance records prior to the meeting.\n\nTip: Enable the integrated video call link for virtual meetings."
        }
    },
    "student": {
        "dashboard": {
            "title": "Student Dashboard",
            "desc": "Your daily hub for deadlines, grades, and schedules.",
            "content": "The student dashboard prioritizes urgent tasks and upcoming exams.\n\nStep 1: Review \"Today’s Schedule\" for lessons and rooms.\nStep 2: Check \"Due Today\" and \"Tomorrow’s Exams\".\nStep 3: Monitor your overall GPA and recent grade additions.\n\nTip: Install the PWA on your phone for instant notification alerts."
        },
        "grades": {
            "title": "My Grades & Performance",
            "desc": "Track your subject grades, averages, and evaluation history.",
            "content": "View all assessments organized by subject and academic term.\n\nStep 1: Select a subject to inspect detailed teacher feedback.\nStep 2: Review the trend line to understand your progress over time.\nStep 3: Switch between terms (1st Term, 2nd Term, Final).\nStep 4: Download a full PDF transcript at any time."
        },
        "homework": {
            "title": "Homework & Assignments",
            "desc": "Organize your homework, downloads, and submission deadlines.",
            "content": "Stay on top of all homework assignments and study materials.\n\nStep 1: Filter tasks by due date or subject.\nStep 2: Click on an assignment to download attachments.\nStep 3: Mark tasks as \"Completed\" to keep track of your workload."
        }
    },
    "parent": {
        "monitoring": {
            "title": "Monitoring Your Children",
            "desc": "Follow academic performance, daily attendance, and notices in real time.",
            "content": "Access comprehensive academic information for all enrolled children.\n\nStep 1: Select your child from the top switcher.\nStep 2: Check daily attendance and real-time grade notifications.\nStep 3: Read teacher comments on individual assessments.\nStep 4: View academic progress charts.\n\nTip: Configure notification thresholds to receive alerts if grades drop."
        },
        "meetings": {
            "title": "Booking Teacher Meetings",
            "desc": "Schedule and manage conference appointments with teachers.",
            "content": "Easily book one-on-one parent-teacher conferences.\n\nStep 1: Open \"Meetings\" and select the desired teacher.\nStep 2: Pick an available green time slot on the calendar.\nStep 3: Confirm booking to receive an email reminder.\n\nWarning: Cancellations should be submitted at least 2 hours in advance."
        }
    },
    "secretary": {
        "students": {
            "title": "Student Registry Management",
            "desc": "Search, manage, and update student profiles and demographic records.",
            "content": "The student registry is the central database for all enrolled learners.\n\nStep 1: Use the global search bar (name, tax code, class).\nStep 2: Click a record to open the full student dossier.\nStep 3: Edit contact information, enrollment status, or guardian links.\nStep 4: Export filtered data to Excel or CSV formats."
        },
        "certificates": {
            "title": "Certificate Generation",
            "desc": "Issue official school documents, attendance proofs, and certificates.",
            "content": "Produce signed school documents in seconds.\n\nStep 1: Search for the student record.\nStep 2: Go to \"Certificates\" → \"Generate\".\nStep 3: Select certificate template (Enrollment, Attendance, Transcripts).\nStep 4: Preview, apply digital stamp, and print or email directly."
        }
    },
    "admin": {
        "monitoring": {
            "title": "System & Health Monitoring",
            "desc": "Monitor platform uptime, API response times, and server performance.",
            "content": "Real-time telemetry and health monitoring across all platform services.\n\nStep 1: Inspect server CPU, RAM, and active session count.\nStep 2: Monitor API throughput and error rates over the last 24 hours.\nStep 3: Configure alert thresholds for proactive notifications.\n\nWarning: Use Maintenance Mode during scheduled platform updates."
        },
        "users": {
            "title": "User & Access Management",
            "desc": "Manage user accounts, RBAC permissions, and authentication policies.",
            "content": "Centralized user provisioning and security management.\n\nStep 1: Search users by role, institution, or active status.\nStep 2: Create single accounts or use the bulk CSV/Excel import tool.\nStep 3: Enforce 2FA/MFA and manage password resets.\nStep 4: Review immutable audit logs for compliance auditing."
        }
    }
}
}
