export default {
  common: {
    welcome: 'Welcome',
    back: 'Back',
    save: 'Save',
    cancel: 'Cancel',
    close: 'Close',
    delete: 'Delete',
    edit: 'Edit',
    search: 'Search (Ctrl+K)...',
    filter: 'Filter',
    actions: 'Actions',
    loading: 'Loading...',
    success: 'Operation completed successfully',
    error: 'An error occurred',
    systemSettings: 'System Settings',
    language: 'Language / Lingua',
    security: 'Security & Authentication',
    notifications: 'Notifications & Alerts',
    mainMenu: 'MAIN MENU',
    logout: 'Log out'
  },
  login: {
    welcomeBack: 'Welcome Back',
    subtitle: 'Sign in to access Electronic School Register',
    emailLabel: 'Email Address',
    emailRequired: 'Email is required',
    passwordLabel: 'Password',
    passwordRequired: 'Password is required',
    rememberMe: 'Remember me',
    submit: 'Sign In',
    noAccount: 'Don\'t have an account?',
    contactSecretary: 'Contact Secretary',
    contactTitle: 'Contact School Secretary',
    contactSubtitle: 'Select your school from the dropdown to view Secretary email and phone details.',
    selectSchool: 'Select your School / Institute',
    noSchoolFound: 'No school found',
    emailSegreteria: 'Secretary Email',
    phone: 'Phone',
    sendEmail: 'Send Email',
    copyEmail: 'Copy Email',
    emailCopied: 'Email address copied to clipboard!',
    chooseSchoolPrompt: 'Choose an institute from the top menu to display Secretary contact details.',
    sessionExpired: 'Session expired. Please sign in again.',
    showPassword: 'Show password',
    hidePassword: 'Hide password'
  },
  nav: {
    dashboard: 'Dashboard',
    schools: 'School Management',
    mySchool: 'My School',
    users: 'User Management',
    admins: 'Admin Management',
    monitoring: 'System Monitoring',
    analytics: 'Global Analytics',
    auditLogs: 'Audit Logs',
    settings: 'Settings',
    support: 'Support & Help',
    featureFlags: 'Feature Flags & Institute',
    elearning: 'Google & Teams E-Learning',
    students: 'Students',
    classes: 'Classes',
    groups: 'Language / Track Groups',
    documents: 'Documents',
    certificates: 'Certificates',
    textbooks: 'Textbooks',
    meetings: 'Meetings',
    communications: 'Communications',
    reports: 'Reports',
    pcto: 'PCTO Work Experience',
    scrutiny: 'Scrutiny & End of Term',
    myClasses: 'My Classes',
    classRegister: 'Class Register',
    uda: 'UdA Unit Planning',
    competencies: 'Competency Assessment',
    grades: 'Grades',
    attendance: 'Attendance',
    didactics: 'Didactic Materials',
    pdp: 'PDP / PEI Plans',
    rubrics: 'Assessment Rubrics',
    coordination: 'Coordination',
    timetable: 'Class Schedule',
    agenda: 'Agenda',
    colloqui: 'Parent Appointments',
    substitutions: 'Substitutions',
    verbali: 'Minutes & Reports',
    notes: 'Disciplinary Notes',
    myGrades: 'My Grades',
    myAttendance: 'My Attendance',
    homework: 'Homework',
    orientamento: 'Career Guidance',
    calendar: 'School Calendar',
    reportCard: 'Report Card',
    profile: 'Profile',
    myChildren: 'My Children',
    goals: 'Goals',
    trips: 'Trips & Excursions',
    payments: 'Payments',
    assemblies: 'Assemblies & Meetings'
  },
  categories: {
    anagraficheClassi: 'Records & Classes',
    attiCertificati: 'Files & Certificates',
    serviziReport: 'Services & Reports',
    didatticaValutazione: 'Teaching & Assessment',
    organizzazioneOrario: 'Schedule & Organization',
    comunicazioniAtti: 'Communications & Files',
    percorsiComunicazioni: 'Pathways & Communications',
    valutazioneDidattica: 'Assessment & Teaching',
    serviziOrari: 'Services & Schedules',
    comunicazioniAccount: 'Communications & Account'
  },
  roles: {
    admin: 'Administrator',
    superadmin: 'Super Administrator',
    secretary: 'Secretary',
    teacher: 'Teacher',
    student: 'Student',
    parent: 'Parent',
    user: 'User'
  },
  notifications: {
    title: 'Notifications & Communications',
    logoutSuccess: 'Logged out successfully',
    logoutError: 'Error during logout',
    settingsSaved: 'Settings saved successfully!',
    emailCopied: 'Email address copied to clipboard!',
    passwordUpdated: 'Password updated successfully!',
    languageChanged: 'Language updated successfully'
  },
  errors: {
    connectionError: 'Server connection error. Please check your internet connection and try again.',
    forbidden: 'You do not have the required permissions to perform this action.',
    serverError: 'An internal server error occurred. Please try again later.',
    unauthorized: 'Invalid or expired session.',
    invalidCredentials: 'Incorrect email or password.',
    rateLimit: 'Too many attempts. Please try again in a few minutes.',
    accountDisabled: 'Account disabled or suspended. Please contact the Secretary.',
    userNotFound: 'User not found.',
    passwordMismatch: 'Passwords do not match.',
    sessionInvalid: 'Invalid session, please log in again.',
    ERR_CURRENT_PASSWORD_INCORRECT: 'Current password is incorrect.',
    ERR_PASSWORD_COMPLEXITY: 'Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character.',
    ERR_PASSWORD_TOO_SHORT: 'Password must be at least 10 characters long.',
    ERR_PASSWORD_TOO_LONG: 'Password is too long (maximum 128 characters).',
    ERR_PASSWORD_RECENTLY_USED: 'New password must not match any of the last 5 used passwords.',
    ERR_REQUIRED_FIELDS: 'Please fill in all required password fields.'
  },
  settings: {
    title: 'System Settings',
    subtitle: 'Configure global system preferences, security policies, and notification alerts',
    languageLabel: 'Language / Lingua',
    securityLabel: 'Security & Authentication',
    securitySub: 'Password policies, 2FA/MFA, and session timeouts',
    notificationsLabel: 'Notifications & Alerts',
    notificationsSub: 'Email preferences, push notifications, and substitution alerts',
    mfaActive: 'MFA Active',
    mfaOptional: 'MFA Optional',
    mfaTitle: 'Enforce 2FA / MFA',
    mfaCaption: 'Require TOTP multi-factor authentication for all staff and administrative users.',
    minPasswordLength: 'Minimum Password Length',
    sessionTimeout: 'Session Inactivity Timeout',
    maxLoginAttempts: 'Failed Login Attempts Threshold'
  },
  search: {
    placeholder: 'Search students, teachers, secretary, classes, notices, menu items…',
    hint: 'Type to search students, teachers, secretary, classes, circulars or menu items',
    noResults: 'No results for',
    navigate: 'Navigate',
    open: 'Open',
    close: 'Close'
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
