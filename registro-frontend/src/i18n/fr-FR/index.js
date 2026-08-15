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
  help: {
    title: 'Centre d\'Aide et de Connaissances',
    subtitle: 'Guides et tutoriels pour tirer le meilleur parti du registre électronique',
    searchPlaceholder: 'Rechercher dans l\'aide...',
    noResults: 'Aucun article trouvé pour',
    categories: 'Catégories',
    allTopics: 'Tous les sujets',
    restartTour: 'Redémarrer la visite guidée',
    openHelp: 'Ouvrir le guide',
    needHelp: 'Besoin d\'aide ?',
    contactSupport: 'Contacter le support',
    fabTooltip: 'Aide & Guide',
    new: 'Nouveau',
    teacher: {
      cat_register: 'Registre de Classe',
      cat_grades: 'Gestion des Notes',
      cat_attendance: 'Présences',
      cat_agenda: 'Agenda & Communications',
      cat_settings: 'Paramètres',
      q1: 'Comment enregistrer les présences de la leçon ?',
      a1: 'Allez dans "Mes Classes" → sélectionnez la classe → onglet "Présences". Cliquez sur P (Présent), A (Absent), R (Retard) pour chaque élève. Cliquez sur "Signer la leçon" pour valider numériquement.',
      q2: 'Comment utiliser la vue Grille Clavier (Matrix View) pour la saisie des notes ?',
      a2: 'Dans "Notes", sélectionnez la classe et la matière, activez le mode Grille Matrix. Naviguez avec TAB et FLÈCHES, saisissez la note et appuyez sur ENTRÉE.',
      q3: 'Comment créer un devoir ou un contrôle dans l\'agenda ?',
      a3: 'Dans "Agenda" → cliquez sur une date → "Ajouter un événement". Choisissez le type (Devoir, Évaluation), indiquez la matière et enregistrez.',
      q4: 'Comment envoyer une communication urgente aux parents ?',
      a4: 'Allez dans "Communications" → "Nouvelle communication". Sélectionnez les destinataires, cochez "Accusé de lecture obligatoire" et envoyez.',
      q5: 'Comment configurer mes disponibilités pour les rendez-vous parents ?',
      a5: 'Dans "Rendez-vous" → "Définir les disponibilités". Définissez les jours, créneaux et durées (15, 20, 30 min).',
      q6: 'Comment accéder à la gestion du Conseil de Classe (Scrutiny) ?',
      a6: 'Si vous êtes professeur principal, "Scrutinio" apparaît dans la barre latérale. Saisissez les propositions de notes et appréciation globale.',
      q7: 'Comment rédiger ou consulter un plan d\'accompagnement (PDP / PEI) ?',
      a7: 'Accédez au "Plan d\'Accompagnement (PDP)" depuis le menu de classe. Configurez les mesures compensatoires et dispensatoires.',
      q8: 'Comment récupérer un brouillon non enregistré ?',
      a8: 'Le système enregistre automatiquement un brouillon toutes les 15 secondes. À la réouverture, le système propose de le restaurer.',
      q9: 'Comment changer la langue d\'interface et le mode sombre ?',
      a9: 'Dans "Paramètres", choisissez parmi 9 langues et activez le mode sombre ou la police OpenDyslexic.',
      q10: 'Comment annuler une note récemment saisie ?',
      a10: 'Après la saisie, un pavé réactif affiche un bouton "Annuler" actif pendant 15 secondes. Ensuite, modifiez directement dans la cellule.'
    },
    student: {
      cat_grades: 'Mes Notes',
      cat_attendance: 'Mes Présences',
      cat_homework: 'Devoirs',
      cat_documents: 'Documents & Bulletin',
      cat_settings: 'Profil & Paramètres',
      q1: 'Comment consulter mes notes et mes moyennes par matière ?',
      a1: 'Consultez "Mes Notes" dans la barre latérale pour voir le détail de vos évaluations écrites/orales et vos moyennes pondérées.',
      q2: 'Comment fonctionne le simulateur de moyenne ?',
      a2: 'Dans "Mes Notes", cliquez sur "Simulateur de moyenne". Entrez des notes fictives pour calculer votre future moyenne.',
      q3: 'Comment suivre le total de mes absences ?',
      a3: 'Dans "Mes Présences", les compteurs affichent votre total d\'heures d\'absence par rapport au seuil annuel.',
      q4: 'Où trouver mes devoirs à faire ?',
      a4: 'Dans "Devoirs" ou le widget du Tableau de bord. Les devoirs sont triés par date de rendu avec pièces jointes.',
      q5: 'Comment télécharger mon bulletin trimestriel ?',
      a5: 'Allez dans "Bulletin & Documents", sélectionnez la période et cliquez sur "Télécharger le bulletin PDF".',
      q6: 'Où consulter mes heures de stage / PCTO ?',
      a6: 'Ouvrez "PCTO & Portfolio" pour suivre les heures validées, les entreprises d\'accueil et évaluations.',
      q7: 'Comment répondre à un sondage ou une circulaire ?',
      a7: 'Ouvrez "Communications" et sélectionnez le message. Cliquez sur "Confirmer la lecture" si demandé.',
      q8: 'Comment activer la police OpenDyslexic ?',
      a8: 'Cliquez sur l\'icône d\'accessibilité ou allez dans "Profil" → "Paramètres d\'accessibilité".',
      q9: 'Que faire en cas d\'oubli de mot de passe ?',
      a9: 'Sur la page de connexion, cliquez sur "Mot de passe oublié ?" et entrez votre adresse e-mail.',
      q10: 'Comment utiliser la recherche rapide Ctrl+K ?',
      a10: 'Appuyez sur Ctrl+K (ou Cmd+K) pour ouvrir la recherche universelle et naviguer instantanément.'
    },
    parent: {
      cat_monitoring: 'Suivi des Enfants',
      cat_communications: 'Communications',
      cat_meetings: 'Rendez-vous',
      cat_documents: 'Documents & Paiements',
      cat_settings: 'Profil & Paramètres',
      q1: 'Comment basculer entre plusieurs enfants scolarisés ?',
      a1: 'Utilisez le sélecteur d\'enfant en haut de page pour basculer d\'un profil à l\'autre instantanément.',
      q2: 'Est-ce que je reçois des notifications push pour les notes ?',
      a2: 'Oui ! Des notifications push et e-mails sont envoyés pour chaque nouvelle note ou absence enregistrée.',
      q3: 'Comment réserver un rendez-vous avec un professeur ?',
      a3: 'Allez dans "Rendez-vous" → sélectionnez l\'enseignant → choisissez un créneau vert dans le calendrier → validez.',
      q4: 'Comment payer les frais scolaires via PagoPA ?',
      a4: 'Accédez à "Paiements PagoPA". Consultez les avis d\'échéance et réglez directement en ligne ou par QR code.',
      q5: 'Comment justifier une absence en ligne ?',
      a5: 'Dans "Présences & Justifications", cliquez sur "Justifier" face à l\'absence rouge, indiquez le motif et validez avec votre code PIN.',
      q6: 'Comment télécharger le bulletin de mon enfant ?',
      a6: 'Sélectionnez votre enfant → "Bulletin" → choisissez le trimestre → "Télécharger le bulletin PDF".',
      q7: 'Où se trouvent les circulaires à signer ?',
      a7: 'Dans "Communications". Les circulaires exigeant un accusé affichent un bouton "Signer la prise d\'acte".',
      q8: 'Comment contacter le secrétariat de l\'école ?',
      a8: 'Dans "Support & FAQ" ou "Contacter le support", envoyez un ticket directement à l\'administration.',
      q9: 'Un second parent peut-il avoir un compte distinct ?',
      a9: 'Le secrétariat peut associer plusieurs comptes responsables à un élève. Chaque parent dispose de ses accès.',
      q10: 'Que se passe-t-il en mode hors ligne ?',
      a10: 'L\'application affiche les données en cache. Les formulaires soumis sont conservés en file locale et synchronisés au retour du réseau.'
    },
    secretary: {
      cat_students: 'Fichier Élèves',
      cat_classes: 'Gestion des Classes',
      cat_documents: 'Certificats & Documents',
      cat_timetable: 'Emploi du Temps & Remplacements',
      cat_reports: 'Rapports & Exports',
      q1: 'Comment rechercher rapidement un élève ?',
      a1: 'Utilisez la barre de recherche globale dans l\'en-tête. Recherche par nom, prénom, code fiscal ou classe.',
      q2: 'Comment délivrer un certificat de scolarité ?',
      a2: 'Dans "Certificats & Documents" → sélectionnez l\'élève → "Générer le PDF". Le document est scellé numériquement.',
      q3: 'Comment gérer les remplacements de professeurs ?',
      a3: 'Allez dans "Emploi du temps & Remplacements" → date du jour → enseignant absent. Le système suggère les professeurs disponibles.',
      q4: 'Comment exporter les données vers le Ministère (SIDI) ?',
      a4: 'Dans "Rapports & Statistiques" → "Export Ministère" → sélectionnez la période et le format d\'export.',
      q5: 'Comment créer un compte pour un nouvel élève ou parent ?',
      a5: 'Dans "Gestion des utilisateurs" → "Nouvel utilisateur" → saisissez l\'état civil et attribuez le rôle.',
      q6: 'Comment configurer la structure des nouvelles classes ?',
      a6: 'Dans "Gestion des classes", créez les nouvelles sections, importez les élèves promus et affectez les enseignants.',
      q7: 'Comment publier une circulaire officielle ?',
      a7: 'Dans "Communications & Panneau d\'affichage" → "Nouvelle circulaire" → sélectionnez les groupes et publiez.',
      q8: 'Comment gérer les manuels scolaires adoptés ?',
      a8: 'Dans "Manuels scolaires", choisissez la classe et saisissez les codes ISBN, titres et éditeurs.',
      q9: 'Comment vérifier le statut des paiements PagoPA ?',
      a9: 'Ouvrez "Finance & PagoPA" pour consulter l\'état des encaissements et relancer les impayés.',
      q10: 'Comment annuler un protocole de document erroné ?',
      a10: 'Dans les archives documentaires, sélectionnez le protocole → "Annuler avec motif". L\'action est consignée dans l\'audit log.'
    },
    admin: {
      cat_monitoring: 'Supervision Système',
      cat_users: 'Gestion des Utilisateurs',
      cat_schools: 'Établissements',
      cat_security: 'Sécurité & Authentification',
      cat_analytics: 'Analyses & Décrochage',
      q1: 'Comment surveiller la santé des microservices et de la base de données ?',
      a1: 'Dans "Supervision système", observez les graphiques CPU/RAM, pools de connexions PostgreSQL et temps de réponse des API.',
      q2: 'Comment ajouter un nouvel établissement scolaire ?',
      a2: 'Dans "Gestion des écoles" → "Ajouter un établissement". Indiquez le code mécanique, nom, adresse et courriel officiel.',
      q3: 'Comment forcer la réinitialisation de mot de passe d\'un compte ?',
      a3: 'Dans "Gestion des utilisateurs", sélectionnez le compte → "Sécurité" → "Forcer la réinitialisation" ou "Exiger la 2FA".',
      q4: 'Comment consulter les journaux d\'audit inaltérables ?',
      a4: 'Dans "Journaux d\'audit & Conformité". Filtrez les événements par utilisateur, type d\'action, date et adresse IP.',
      q5: 'Comment gérer les Feature Flags de la plateforme ?',
      a5: 'Dans "Feature Flags & Modules", activez ou désactivez des fonctionnalités spécifiques (PagoPA, Google SSO, Matrix View).',
      q6: 'Comment analyser le risque de décrochage scolaire ?',
      a6: 'Ouvrez "Analytics & BI". Les algorithmes signalent les élèves dépassant 20% d\'absences ou ayant une moyenne inférieure à 6/10.',
      q7: 'Comment configurer l\'authentification unique (SSO) Google et Microsoft ?',
      a7: 'Dans "Intégrations & SSO", saisissez le Client ID OAuth2 et la Clé Secrète depuis les consoles Google ou Azure.',
      q8: 'Comment exécuter une sauvegarde manuelle de la base de données ?',
      a8: 'Dans "Base de données & Sauvegardes", déclenchez un dump PostgreSQL instantané ou planifiez des sauvegardes S3.',
      q9: 'Comment vérifier la matrice des droits RBAC ?',
      a9: 'Allez dans "Rôles & Permissions" pour consulter la matrice d\'accès des 9 rôles du système.',
      q10: 'Comment placer la plateforme en mode maintenance ?',
      a10: 'Dans "Paramètres système", activez le "Mode Maintenance". Un bandeau d\'information s\'affiche pour les utilisateurs.'
    }
  },
  guideCenter: {
    title: "Centre d'Aide & Connaissances",
    subtitle: "Guides détaillés pour chaque section du registre électronique",
    search: "Rechercher dans les guides...",
    noResults: "Aucun guide trouvé pour",
    readingTime: "min de lecture",
    step: "Étape",
    tip: "Conseil",
    warning: "Attention",
    shortcut: "Raccourci",
    teacher: {
      dashboard: { title: "Tableau de bord Enseignant", desc: "Aperçu de la journée, notifications et actions rapides.", content: "Le tableau de bord enseignant offre une vue synthétique de votre journée. Les widgets affichent vos cours du jour, horaires et salles.\n\nÉtape 1 : Personnalisez vos widgets en les faisant glisser.\nÉtape 2 : Cliquez sur une classe pour ouvrir directement son registre.\nÉtape 3 : Filtrez les notifications par urgence.\n\nConseil : Activez les alertes prioritaires pour recevoir les notifications importantes." },
      attendance: { title: "Registre des Présences", desc: "Saisie rapide et précise des présences, absences et retards.", content: "Le registre charge automatiquement la liste des élèves pour l'heure de cours.\n\nÉtape 1 : Cliquez sur le nom d'un élève pour modifier son statut (Présent, Absent, Retard).\nÉtape 2 : Saisissez des remarques ou motifs optionnels.\nÉtape 3 : Cliquez sur \"Signer la leçon\" pour signer numériquement le cours.\nÉtape 4 : Les parents des élèves absents sont prévenus automatiquement.\n\nRaccourci : Appuyez sur Ctrl+Entrée pour enregistrer rapidement." },
      grades: { title: "Gestion des Notes & Vue Matrix", desc: "Saisie des notes écrites/orales avec la grille rapide au clavier.", content: "Le carnet de notes propose la vue classique et la vue grille rapide Matrix View.\n\nÉtape 1 : Sélectionnez la classe et la matière.\nÉtape 2 : Activez la vue Matrix pour naviguer avec TAB et les FLÈCHES.\nÉtape 3 : Tapez la note (1-10) et appuyez sur ENTRÉE.\nÉtape 4 : Ajoutez des commentaires publics ou privés.\nÉtape 5 : Les moyennes par matière sont recalculées instantanément.\n\nAttention : Les notes publiées sont visibles par les parents en quelques secondes." },
      agenda: { title: "Agenda & Devoirs", desc: "Planification des évaluations, devoirs et activités.", content: "L'agenda enseignant est synchronisé avec celui des élèves et des parents.\n\nÉtape 1 : Cliquez sur une date du calendrier pour créer un événement.\nÉtape 2 : Choisissez le type : Devoir, Évaluation, Sortie pédagogique.\nÉtape 3 : Indiquez la matière, la description et joignez des fichiers.\nÉtape 4 : Le système vous alerte si la classe a déjà plus de 2 évaluations le même jour.\n\nConseil : Utilisez la vue semaine pour vérifier les chevauchements d'évaluations." },
      communications: { title: "Communications & Message", desc: "Envoi de messages structurés aux classes, parents et collègues.", content: "Publiez des annonces via les canaux dédiés : messages directs, circulaires ou bulletins urgents.\n\nÉtape 1 : Naviguez vers \"Communications\" → \"Nouveau message\".\nÉtape 2 : Choisissez le canal et la liste de destinataires.\nÉtape 3 : Rédigez le message avec l'éditeur de texte enrichi.\nÉtape 4 : Envoyez immédiatement ou programmez la publication.\n\nAttention : Les messages urgents déclenchent des notifications push." },
      meetings: { title: "Rendez-vous Parents-Professeurs", desc: "Gestion des créneaux de réception et rendez-vous.", content: "Gérez vos heures de réception hebdomadaires et rendez-vous individuels.\n\nÉtape 1 : Allez dans 'Rendez-vous' → 'Définir les disponibilités'.\nÉtape 2 : Définissez les jours, créneaux et durées (15, 20, 30 min).\nÉtape 3 : Les parents réservent en ligne ; vous recevez une confirmation.\nÉtape 4 : Consultez le dossier de l'élève avant l'entretien.\n\nConseil : Activez le lien de visioconférence intégré." },
      scrutiny: { title: "Conseil de Classe & Scrutinio", desc: "Procédure des propositions de notes, appréciations et procès-verbaux.", content: "Le tableau du conseil de classe consolide les propositions de notes de tous les enseignants.\n\nÉtape 1 : Sélectionnez la classe et la période de conseil.\nÉtape 2 : Vérifiez la moyenne proposée et saisissez la note délibérée.\nÉtape 3 : Saisissez la note de conduite et l'appréciation globale.\nÉtape 4 : En tant que professeur principal, verrouillez le tableau une fois validé.\nÉtape 5 : Exprimez et imprimez le procès-verbal en PDF.\n\nAttention : Après verrouillage, toute modification exige un déverrouillage administratif." },
      pdp: { title: "Plan d'Accompagnement (PDP / PEI)", desc: "Élaboration des mesures compensatoires et dispensatoires pour élèves BES/DSA.", content: "Le module PDP permet à l'équipe pédagogique d'élaborer les plans personnalisés.\n\nÉtape 1 : Accédez à 'PDP / PEI' depuis le menu de classe.\nÉtape 2 : Sélectionnez l'élève et configurez les mesures compensatoires et dispensatoires.\nÉtape 3 : Enregistrez le brouillon et transmettez aux parents pour signature numérique.\nÉtape 4 : Les mesures actives apparaissent sous forme de badges lors des évaluations.\n\nConseil : Utilisez les grilles d'évaluation préremplies pour accélérer la saisie." }
    },
    student: {
      dashboard: { title: "Tableau de bord Élève", desc: "Accès quotidien à vos devoirs, notes et emplois du temps.", content: "Le tableau de bord élève met en avant les échéances urgentes et contrôles.\n\nÉtape 1 : Consultez 'Aujourd'hui' pour vos cours et salles.\nÉtape 2 : Vérifiez les devoirs à rendre et évaluations du lendemain.\nÉtape 3 : Suivez votre moyenne générale.\n\nConseil : Installez la PWA sur votre téléphone pour recevoir des alertes." },
      grades: { title: "Mes Notes & Résultats", desc: "Suivez vos notes, moyennes et historique d'évaluations.", content: "Toutes vos évaluations classées par matière et par trimestre.\n\nÉtape 1 : Sélectionnez une matière pour lire les appréciations des professeurs.\nÉtape 2 : Examinez la courbe d'évolution de vos résultats.\nÉtape 3 : Basculez entre les trimestres.\nÉtape 4 : Téléchargez votre relevé complet au format PDF." },
      homework: { title: "Devoirs & Travaux", desc: "Organisez vos devoirs, fichiers à télécharger et échéances.", content: "Gardez un contrôle total sur vos devoirs et documents de cours.\n\nÉtape 1 : Filtrez les devoirs par date ou matière.\nÉtape 2 : Cliquez sur un devoir pour télécharger les pièces jointes.\nÉtape 3 : Cochez les devoirs comme 'Fait'." },
      attendance: { title: "Mes Présences & Absences", desc: "Suivi du volume d'absences, retards et statut des justifications.", content: "Consultez votre bilan d'absences et le seuil d'heures autorisé.\n\nÉtape 1 : Ouvrez 'Mes Présences'.\nÉtape 2 : Vérifiez le compteur d'heures d'absence cumulées.\nÉtape 3 : Examinez le calendrier mensuel pour les jours absents (rouge) et retards (orange).\nÉtape 4 : Vérifiez la validation des justifications d'absence par vos parents." },
      documents: { title: "Bulletins & Documents Officiels", desc: "Consultation et téléchargement des bulletins scolaires et certificats.", content: "Accédez aux documents officiels publiés par la direction de l'école.\n\nÉtape 1 : Allez dans 'Bulletin & Documents'.\nÉtape 2 : Sélectionnez l'année et le trimestre.\nÉtape 3 : Cliquez sur 'Aperçu du bulletin' pour consulter en ligne.\nÉtape 4 : Cliquez sur 'Télécharger le PDF' pour enregistrer le document signé." },
      simulator: { title: "Simulateur de Moyenne & Objectif", desc: "Calculez la note nécessaire aux prochains devoirs pour atteindre votre objectif.", content: "Le simulateur fournit des projections mathématiques pour vos évaluations futures.\n\nÉtape 1 : Sélectionnez la matière ciblée.\nÉtape 2 : Entrez la moyenne souhaitée (ex: 14/20).\nÉtape 3 : Le système calcule la note minimale requise au prochain contrôle." },
      pcto: { title: "PCTO & Portfolio de Compétences", desc: "Gestion des heures de stage et attestations de formation.", content: "Suivez vos activités et heures de stage PCTO au lycée.\n\nÉtape 1 : Ouvrez 'PCTO & Portfolio'.\nÉtape 2 : Vérifiez votre graphique d'heures par rapport au volume obligatoire.\nÉtape 3 : Consultez les informations de l'entreprise d'accueil et téléchargez vos attestations." }
    },
    parent: {
      monitoring: { title: "Suivi des Enfants", desc: "Suivez la scolarité, les présences et les messages en temps réel.", content: "Accès complet au dossier scolaire de tous vos enfants inscrits.\n\nÉtape 1 : Sélectionnez l'enfant dans le menu supérieur.\nÉtape 2 : Vérifiez les présences du jour et notifications de notes.\nÉtape 3 : Lisez les appréciations des professeurs sur les devoirs.\nÉtape 4 : Observez les graphiques de progression.\n\nConseil : Configurez les seuils d'alerte pour recevoir une notification en cas de baisse de moyenne." },
      meetings: { title: "Prise de Rendez-vous Professeurs", desc: "Planifiez des rencontres avec les enseignants de votre enfant.", content: "Réservez des entretiens individuels en quelques clics.\n\nÉtape 1 : Ouvrez 'Rendez-vous' et choisissez l'enseignant.\nÉtape 2 : Sélectionnez un créneau vert disponible dans le calendrier.\nÉtape 3 : Validez la réservation pour recevoir la confirmation par e-mail.\n\nAttention : Les annulations doivent être effectuées au moins 2 heures à l'avance." },
      communications: { title: "Communications & Panneau d'Affichage", desc: "Consultez les circulaires, alertes et annonces de classe.", content: "Toutes les annonces officielles de l'établissement réunies dans un espace unique.\n\nÉtape 1 : Ouvrez 'Communications'.\nÉtape 2 : Filtrez par catégorie (Direction, Professeurs, Classe).\nÉtape 3 : Lisez le contenu et téléchargez les pièces jointes PDF.\nÉtape 4 : Cliquez sur 'Signer la prise d'acte' lorsque demandé." },
      pagopa: { title: "Paiements Scolaires PagoPA", desc: "Gérez et réglez les sorties, cantine et frais scolaires en ligne.", content: "Système PagoPA intégré pour un paiement électronique sécurisé.\n\nÉtape 1 : Allez dans 'Paiements PagoPA'.\nÉtape 2 : Consultez les appels de fonds en cours.\nÉtape 3 : Cliquez sur 'Payer maintenant' pour régler par carte, PayPal ou prélèvement.\nÉtape 4 : Ou téléchargez le coupon QR pour régler à la banque ou chez un buraliste." },
      documents: { title: "Documents & Formulaires Parents", desc: "Téléchargez les bulletins, certificats et autorisations.", content: "Récupérez les documents officiels signés directement sur votre appareil.\n\nÉtape 1 : Allez dans 'Documents & Bulletin'.\nÉtape 2 : Téléchargez les bulletins trimestriels au format PDF.\nÉtape 3 : Remplissez et signez numériquement les autorisations de sortie." },
      justifications: { title: "Justification d'Absences en Ligne", desc: "Justifiez les absences et retards de votre enfant en ligne.", content: "Plus besoin de carnet papier : justifiez les absences en toute sécurité avec votre PIN.\n\nÉtape 1 : Ouvrez 'Présences & Justifications'.\nÉtape 2 : Visualisez les absences non justifiées en rouge.\nÉtape 3 : Cliquez sur 'Justifier', sélectionnez le motif (Maladie, Famille) et validez avec votre code PIN." }
    },
    secretary: {
      students: { title: "Gestion du Fichier Élèves", desc: "Recherche, modification et mise à jour des dossiers d'élèves.", content: "Le fichier élèves constitue la base de données centrale de l'établissement.\n\nÉtape 1 : Utilisez la barre de recherche globale (Nom, Prénom, Classe).\nÉtape 2 : Cliquez sur une fiche pour ouvrir le dossier complet.\nÉtape 3 : Modifiez les coordonnées ou responsables légaux.\nÉtape 4 : Exporter les données filtrées au format Excel ou CSV." },
      classes: { title: "Organisation des Classes & Sections", desc: "Création de classes, attribution des matières et enseignants.", content: "Structurez la répartition des classes pour l'année scolaire.\n\nÉtape 1 : Ouvrez 'Gestion des classes'.\nÉtape 2 : Cliquez sur 'Nouvelle classe' pour ajouter une section.\nÉtape 3 : Affectez les matières et associez les enseignants titulaires.\nÉtape 4 : Désignez le professeur principal et le secrétaire de classe." },
      certificates: { title: "Émission de Certificats", desc: "Générez des certificats officiels de scolarité et de présence.", content: "Produisez des documents scolaires signés en quelques secondes.\n\nÉtape 1 : Recherchez l'élève concerné.\nÉtape 2 : Allez dans 'Certificats' → 'Générer'.\nÉtape 3 : Sélectionnez le modèle (Scolarité, Présence, Relevé de notes).\nÉtape 4 : Vérifiez l'aperçu, appliquez le tampon numérique et imprimez ou envoyez." },
      timetable: { title: "Emploi du Temps & Remplacements", desc: "Saisie des emplois du temps et gestion des remplacements quotidiens.", content: "Gérez la grille horaire et le remplacement des enseignants absents.\n\nÉtape 1 : Ouvrez 'Emploi du temps & Remplacements'.\nÉtape 2 : Saisissez les cours sur la grille hebdomadaire.\nÉtape 3 : En cas d'absence, le système suggère les enseignants disponibles.\nÉtape 4 : Validez le remplacement et notifiez le professeur concerné." },
      communications: { title: "Gestion des Circulaires", desc: "Publication de circulaires officielles et suivi des réceptions.", content: "Diffusez les informations aux personnels, parents et élèves.\n\nÉtape 1 : Ouvrez 'Communications & Panneau d'affichage'.\nÉtape 2 : Cliquez sur 'Nouvelle circulaire'.\nÉtape 3 : Renseignez le numéro de protocole, titre et texte.\nÉtape 4 : Définissez les groupes cibles et l'exigence de signature." },
      reports: { title: "Rapports & Exports Ministère (SIDI)", desc: "Extraction de statistiques et fichiers d'export réglementaires.", content: "Produisez les rapports statistiques agrégés et fichiers d'export.\n\nÉtape 1 : Ouvrez 'Rapports & Statistiques'.\nÉtape 2 : Choisissez le type de rapport (Présences, Résultats des conseils, Listes).\nÉtape 3 : Sélectionnez le format : Excel, CSV, PDF ou XML SIDI." }
    },
    admin: {
      monitoring: { title: "Supervision & Santé Système", desc: "Surveillance de la disponibilité, des temps de réponse et des serveurs.", content: "Télémesure en temps réel de tous les microservices de la plateforme.\n\nÉtape 1 : Inspectez l'utilisation CPU, la RAM et les sessions actives.\nÉtape 2 : Surveillez le débit des API et le taux d'erreur sur 24 heures.\nÉtape 3 : Configurez les seuils d'alerte automatique.\n\nAttention : Activez le Mode Maintenance lors des mises à jour planifiées." },
      users: { title: "Gestion des Utilisateurs & Droits", desc: "Administration des comptes, droits RBAC et politiques d'accès.", content: "Gestion centralisée de la sécurité et des utilisateurs.\n\nÉtape 1 : Recherchez les utilisateurs par rôle, école ou statut.\nÉtape 2 : Créez des comptes individuels ou utilisez l'import CSV massif.\nÉtape 3 : Imposez la 2FA/MFA et gérez la réinitialisation des mots de passe.\nÉtape 4 : Consultez les journaux d'audit inaltérables." },
      schools: { title: "Gestion des Établissements", desc: "Configuration des sites scolaires, codes mécaniques et réseaux.", content: "Gérez les structures de l'établissement multi-tenant.\n\nÉtape 1 : Ouvrez 'Gestion des écoles'.\nÉtape 2 : Ajoutez les nouveaux sites (Siège, Annexe, École primaire, Collège).\nÉtape 3 : Renseignez le code mécanique et les coordonnées de contact." },
      security: { title: "Politiques de Sécurité & Auth", desc: "Règles de mots de passe, expiration des sessions, 2FA et limitation de débit.", content: "Configurez les règles de protection conformément au RGPD.\n\nÉtape 1 : Ouvrez 'Sécurité & Authentification'.\nÉtape 2 : Définissez la longueur minimale de mot de passe et l'échéance.\nÉtape 3 : Activez l'obligation de 2FA pour les personnels d'encadrement." },
      analytics: { title: "Analytics & Décrochage Scolaire", desc: "Analyse prédictive des risques d'abandon et d'absentéisme.", content: "Utilisez les outils décisionnels pour repérer les élèves à risque.\n\nÉtape 1 : Ouvrez 'Analytics & BI'.\nÉtape 2 : Consultez la carte thermique de l'absentéisme par classe.\nÉtape 3 : Définissez les règles d'alerte (ex: Absences > 20% + Moyenne < 10/20)." },
      integrations: { title: "Intégrations E-Learning & SSO", desc: "Synchronisation avec Google Classroom, Microsoft Teams et SSO.", content: "Connectez le Registre Électronique aux plates-formes cloud LMS.\n\nÉtape 1 : Ouvrez 'Intégrations & SSO'.\nÉtape 2 : Activez le module Google Workspace ou Microsoft 365.\nÉtape 3 : Renseignez le Client ID et le Secret OAuth2." },
      audit: { title: "Journaux d'Audit & Traçabilité RGPD", desc: "Journal d'activité système inaltérable traçant tous les accès.", content: "Garantissez la traçabilité intégrale conformément au RGPD 2016/679.\n\nÉtape 1 : Ouvrez 'Journaux d'audit & Traçabilité'.\nÉtape 2 : Visualisez l'historique : Horodatage, Utilisateur, Rôle, IP, Action.\nÉtape 3 : Exportez des rapports PDF chiffrés pour les inspections." }
    }
  }
}
