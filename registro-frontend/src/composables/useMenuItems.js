/**
 * Menu configuration for all 25 institutional roles & dynamic duties
 * Returns menu items based on user role and active user assignments
 *
 * STRUTTURA PER RUOLO:
 * - quickAccess: voci sempre visibili (max 3) senza dover espandere nulla
 * - categories: gruppi espandibili (max 4, max 7 voci ciascuno)
 *
 * Ogni categoria puÃ² avere badge: il nome del contatore da useMenuBadges()
 */
export function useMenuItems(role, assignments = [], isCanteenEnabled = true) {
    const normRole = (role || '').toLowerCase()

    if (!normRole) {
        return []
    }

    // SUPERADMIN
    const superadminMenu = [
        { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
        {
            category: 'Gestione Piattaforma',
            icon: 'public',
            children: [
                { label: 'Gestione Scuole',     icon: 'school',               path: '/admin/schools' },
                { label: 'Gestione Utenti',     icon: 'people',               path: '/admin/users' },
                { label: 'Gestione Admin',      icon: 'admin_panel_settings', path: '/admin/admins' },
                { label: 'Monitoraggio Sistema', icon: 'monitor_heart',       path: '/admin/monitoring' },
                { label: 'Enterprise & Compliance', icon: 'verified_user',    path: '/admin/enterprise' },
            ]
        },
        {
            category: 'Analytics & Sicurezza',
            icon: 'analytics',
            children: [
                { label: 'Analytics Globali', icon: 'analytics', path: '/admin/analytics' },
                { label: 'Audit Logs',        icon: 'history',   path: '/admin/audit-logs' },
                { label: 'Impostazioni',      icon: 'settings',  path: '/admin/settings' },
                { label: 'Supporto',          icon: 'help',      path: '/support' },
            ]
        }
    ]

    // ADMIN
    const adminMenu = [
        { label: 'Dashboard',       icon: 'dashboard', path: '/',                    exact: true, quickAccess: true },
        { label: 'Gestione Utenti', icon: 'people',    path: '/admin/users',         quickAccess: true },
        { label: 'Orario & Aule',  icon: 'schedule',  path: '/secretary/timetable', quickAccess: true },
        {
            category: 'Istituto & Personale',
            icon: 'corporate_fare',
            children: [
                { label: 'La Mia Scuola',         icon: 'school',       path: '/admin/schools' },
                { label: 'Utenti & Ruoli',        icon: 'people',       path: '/admin/users' },
                { label: 'Feature & Impostazioni', icon: 'toggle_on',   path: '/admin/school-settings' },
                { label: 'Presenze Personale',    icon: 'co_present',   path: '/ata/attendance',          badge: 'absentStaff' },
                { label: 'Scioperi',              icon: 'campaign',     path: '/ata/strike' },
                { label: 'Sostituzioni',          icon: 'swap_horiz',   path: '/secretary/substitutions', badge: 'pendingSubstitutions' },
                { label: 'Aule & Plessi',         icon: 'meeting_room', path: '/secretary/rooms' },
            ]
        },
        {
            category: 'Report & Strumenti',
            icon: 'bar_chart',
            children: [
                { label: 'Orario & Cattedre',          icon: 'schedule',  path: '/secretary/timetable' },
                { label: 'Analytics',                  icon: 'analytics', path: '/admin/analytics' },
                { label: 'E-Learning (Google/Teams)',  icon: 'hub',       path: '/admin/elearning' },
                { label: 'Enterprise & Compliance',    icon: 'verified_user', path: '/admin/enterprise' },
                { label: 'Verbali & Riunioni',         icon: 'gavel',     path: '/secretary/verbali' },
                { label: 'Impostazioni',               icon: 'settings',  path: '/admin/settings' },
            ]
        }
    ]

    // SECRETARY / PRINCIPAL 
    const secretaryMenu = [
        { label: 'Dashboard',           icon: 'dashboard',  path: '/',                        exact: true },
        { label: 'Studenti & Famiglie', icon: 'school',     path: '/secretary/students',      quickAccess: true },
        { label: 'Classi',              icon: 'room',       path: '/secretary/classes',        quickAccess: true },
        { label: 'Sostituzioni',        icon: 'swap_horiz', path: '/secretary/substitutions', quickAccess: true, badge: 'pendingSubstitutions' },
        {
            category: 'Anagrafica & Organizzazione',
            icon: 'school',
            children: [
                { label: 'Utenti & Personale', icon: 'people',       path: '/secretary/users' },
                { label: 'Classi & Gruppi',    icon: 'room',         path: '/secretary/classes' },
                { label: 'Gruppi Linguistici', icon: 'groups',       path: '/secretary/groups' },
                { label: 'Aule & Plessi',      icon: 'meeting_room', path: '/secretary/rooms' },
                { label: 'Orario Scolastico',  icon: 'schedule',     path: '/secretary/timetable' },
                { label: 'Vincoli Orario',     icon: 'tune',         path: '/secretary/timetable-constraints' },
            ]
        },
        {
            category: 'Atti & Certificati',
            icon: 'folder_shared',
            children: [
                { label: 'Documenti & Verbali', icon: 'description',       path: '/secretary/documents' },
                { label: 'Riunioni',            icon: 'groups',            path: '/secretary/meetings' },
                { label: 'Certificati',         icon: 'workspace_premium', path: '/secretary/certificates' },
                { label: 'Libri di Testo',      icon: 'auto_stories',      path: '/secretary/textbooks' },
                { label: 'Flussi SIDI',         icon: 'cloud_sync',        path: '/secretary/sidi' },
                { label: 'Scrutinio & Esami',   icon: 'analytics',         path: '/secretary/scrutiny' },
            ]
        },
        {
            category: 'Servizi & Report',
            icon: 'manage_accounts',
            children: [
                { label: 'Presenze Personale', icon: 'co_present', path: '/ata/attendance',          badge: 'absentStaff' },
                { label: 'Scioperi',           icon: 'campaign',   path: '/ata/strike' },
                { label: 'Comunicazioni',      icon: 'email',      path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Enterprise & Compliance', icon: 'verified_user', path: '/admin/enterprise' },
                { label: 'Report & PCTO',      icon: 'assessment', path: '/secretary/reports' },
                { label: 'Impostazioni',       icon: 'settings',   path: '/secretary/settings' },
            ]
        }
    ]

    // TEACHER
    const teacherMenu = [
        { label: 'Dashboard', icon: 'dashboard',  path: '/',                   exact: true },
        { label: 'Registro',  icon: 'menu_book',  path: '/teacher/lessons',    quickAccess: true },
        { label: 'Voti',      icon: 'grade',       path: '/teacher/grades',    quickAccess: true, badge: 'pendingGrades' },
        { label: 'Presenze',  icon: 'how_to_reg', path: '/teacher/attendance', quickAccess: true },
        {
            category: 'Registro & Classi',
            icon: 'menu_book',
            children: [
                { label: 'Le Mie Classi',         icon: 'class',             path: '/teacher/classes' },
                { label: 'Programmazione UdA',    icon: 'auto_stories',      path: '/teacher/uda' },
                { label: 'Valutazione Competenze', icon: 'stars',            path: '/teacher/competencies' },
                { label: 'Rubriche Valutative',   icon: 'checklist',         path: '/teacher/rubrics' },
                { label: 'Credito Scolastico',    icon: 'military_tech',     path: '/teacher/credits' },
                { label: 'Corsi Recupero & PAI',  icon: 'school',            path: '/teacher/recovery' },
                { label: 'PDP/PEI & Sostegno',   icon: 'accessibility_new', path: '/teacher/pdp' },
            ]
        },
        {
            category: 'Scrutinio & Valutazione',
            icon: 'analytics',
            children: [
                { label: 'Coordinamento',     icon: 'groups',         path: '/teacher/classes', coordinatorOnly: true },
                { label: 'Scrutinio',         icon: 'analytics',      path: '/teacher/scrutiny', coordinatorOnly: true },
                { label: 'Note Disciplinari', icon: 'assignment_late', path: '/teacher/notes' },
                { label: 'Documenti',         icon: 'description',    path: '/teacher/documents' },
                { label: 'Verbali Riunioni',  icon: 'gavel',          path: '/teacher/verbali' },
            ]
        },
        {
            category: 'Orario & Agenda',
            icon: 'event',
            children: [
                { label: 'Il Mio Orario',        icon: 'schedule',      path: '/teacher/timetable' },
                { label: 'Desiderata Orario',    icon: 'thumb_up_alt',  path: '/teacher/schedule-preferences' },
                { label: 'Prenota Aula/Lab',     icon: 'meeting_room',  path: '/teacher/room-booking' },
                { label: 'Agenda & Calendario',  icon: 'edit_calendar', path: '/teacher/agenda' },
                { label: 'Colloqui',             icon: 'event',         path: '/teacher/colloqui',        badge: 'pendingColloqui' },
                { label: 'Ricevimento Genitori', icon: 'people',        path: '/teacher/general-meetings' },
                { label: 'Sostituzioni',         icon: 'swap_horiz',    path: '/teacher/substitutions' },
            ]
        },
        {
            category: 'Comunicazioni',
            icon: 'campaign',
            children: [
                { label: 'Messaggi & Circolari', icon: 'email',    path: '/teacher/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni',         icon: 'settings', path: '/teacher/settings' },
            ]
        }
    ]

    // VICE_PRINCIPAL
    const vicePrincipalMenu = [
        { label: 'Dashboard Vicario', icon: 'dashboard',  path: '/teacher',                   exact: true },
        { label: 'Sostituzioni',      icon: 'swap_horiz', path: '/secretary/substitutions',   quickAccess: true, badge: 'pendingSubstitutions' },
        { label: 'Orario Scuola',     icon: 'schedule',   path: '/secretary/timetable',        quickAccess: true },
        { label: 'Registro (mio)',    icon: 'menu_book',  path: '/teacher/lessons',            quickAccess: true },
        {
            category: 'Presidenza & Vicariato',
            icon: 'account_balance',
            children: [
                { label: 'Gestione Sostituzioni', icon: 'swap_horiz', path: '/secretary/substitutions', badge: 'pendingSubstitutions' },
                { label: 'Presenze Personale', icon: 'co_present',   path: '/ata/attendance',                badge: 'absentStaff' },
                { label: 'Scioperi',           icon: 'campaign',     path: '/ata/strike' },
                { label: 'Aule & Plessi',      icon: 'meeting_room', path: '/secretary/rooms' },
                { label: 'Classi & Gruppi',    icon: 'room',         path: '/secretary/classes' },
                { label: 'Studenti',           icon: 'school',       path: '/secretary/students' },
                { label: 'Vincoli Orario',     icon: 'tune',         path: '/secretary/timetable-constraints' },
            ]
        },
        {
            category: 'Didattica & Le Mie Classi',
            icon: 'menu_book',
            children: [
                { label: 'Le Mie Classi',    icon: 'class',             path: '/teacher/classes' },
                { label: 'Voti',             icon: 'grade',             path: '/teacher/grades',  badge: 'pendingGrades' },
                { label: 'Scrutinio',        icon: 'analytics',         path: '/teacher/scrutiny' },
                { label: 'PDP/PEI',         icon: 'accessibility_new', path: '/teacher/pdp' },
                { label: 'Agenda & Colloqui', icon: 'edit_calendar',    path: '/teacher/agenda' },
            ]
        },
        {
            category: 'Atti & Comunicazioni',
            icon: 'campaign',
            children: [
                { label: 'Documenti & Verbali',   icon: 'description',       path: '/secretary/documents' },
                { label: 'Certificati',           icon: 'workspace_premium', path: '/secretary/certificates' },
                { label: 'Libri di Testo',        icon: 'auto_stories',      path: '/secretary/textbooks' },
                { label: 'Comunicazioni Istituto', icon: 'email',            path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni',          icon: 'settings',          path: '/teacher/settings' },
            ]
        }
    ]

    // STUDENT
    const studentMenu = [
        { label: 'Dashboard',     icon: 'dashboard',  path: '/',                  exact: true },
        { label: 'I Miei Voti',   icon: 'grade',      path: '/student/grades',    quickAccess: true },
        { label: 'Compiti',       icon: 'assignment', path: '/student/homework',  quickAccess: true, badge: 'pendingHomework' },
        { label: 'Il Mio Orario', icon: 'schedule',   path: '/student/timetable', quickAccess: true },
        {
            category: 'Scuola',
            icon: 'school',
            children: [
                { label: 'Presenze',           icon: 'how_to_reg',     path: '/student/attendance' },
                { label: 'Materiale Didattico', icon: 'folder_shared', path: '/student/didactics' },
                { label: 'Pagella',            icon: 'description',    path: '/student/report-card' },
                { label: 'Note Disciplinari',  icon: 'assignment_late', path: '/student/notes' },
                { label: 'Obiettivi',          icon: 'flag',           path: '/student/goals' },
            ]
        },
        {
            category: 'Servizi & Profilo',
            icon: 'manage_accounts',
            children: [
                { label: 'Agenda & Calendario',   icon: 'edit_calendar', path: '/student/agenda' },
                { label: 'PCTO & Orientamento',   icon: 'work',          path: '/student/pcto' },
                { label: 'Sportello d\'Ascolto (CIC)', icon: 'psychology', path: '/student/psychology' },
                ...(isCanteenEnabled ? [{ label: 'Mensa Scolastica', icon: 'restaurant', path: '/student/canteen' }] : []),
                { label: 'Comunicazioni',         icon: 'email',         path: '/student/communications', badge: 'unreadMessages' },
                { label: 'Documenti',             icon: 'description',   path: '/student/documents' },
                { label: 'Impostazioni & Profilo', icon: 'settings',     path: '/student/settings' },
            ]
        }
    ]

    // PARENT
    const parentMenu = [
        { label: 'Dashboard',    icon: 'dashboard',       path: '/parent',               exact: true },
        { label: 'I Miei Figli', icon: 'family_restroom', path: '/parent/children',       quickAccess: true },
        { label: 'Voti',         icon: 'grade',           path: '/parent/grades',         quickAccess: true },
        { label: 'Comunicazioni', icon: 'email',          path: '/parent/communications', quickAccess: true, badge: 'unreadMessages' },
        {
            category: 'Situazione Scolastica',
            icon: 'school',
            children: [
                { label: 'Presenze',           icon: 'how_to_reg',      path: '/parent/attendance' },
                { label: 'Pagella',            icon: 'description',     path: '/parent/report-card' },
                { label: 'Note Disciplinari',  icon: 'assignment_late', path: '/parent/notes' },
                { label: 'Piano PDP/PEI',      icon: 'accessibility_new', path: '/parent/pdp' },
                { label: 'Materiale Didattico', icon: 'folder_shared',  path: '/parent/didactics' },
            ]
        },
        {
            category: 'Servizi & Contatti',
            icon: 'event',
            children: [
                { label: 'Orario Lezioni',         icon: 'schedule',          path: '/parent/timetable' },
                { label: 'Colloqui & Ricevimenti', icon: 'event',             path: '/parent/colloqui', badge: 'pendingColloqui' },
                { label: 'Sportello Psicologico (CIC)', icon: 'psychology',   path: '/parent/psychology' },
                ...(isCanteenEnabled ? [{ label: 'Mensa & Borsellino Pasti', icon: 'restaurant', path: '/parent/canteen' }] : []),
                { label: 'Uscite & Viaggi',        icon: 'card_travel',       path: '/parent/trips' },
                { label: 'Pagamenti',              icon: 'payments',          path: '/parent/payments' },
                { label: 'Documenti',              icon: 'description',       path: '/parent/documents' },
                { label: 'Impostazioni & Profilo', icon: 'settings',          path: '/parent/settings' },
            ]
        }
    ]

    // DSGA
    const dsgaMenu = [
        { label: 'Dashboard ATA',      icon: 'dashboard',  path: '/ata',                     exact: true },
        { label: 'Presenze Personale', icon: 'co_present', path: '/ata/attendance',           quickAccess: true, badge: 'absentStaff' },
        { label: 'Sostituzioni',       icon: 'swap_horiz', path: '/secretary/substitutions',  quickAccess: true, badge: 'pendingSubstitutions' },
        {
            category: 'Personale & Presenze',
            icon: 'co_present',
            children: [
                { label: 'Presenze & Timbrature', icon: 'co_present',    path: '/ata/attendance' },
                { label: 'Scioperi',              icon: 'campaign',      path: '/ata/strike' },
                { label: 'Cartellino & Ferie',    icon: 'calendar_month', path: '/ata/timecard' },
                { label: 'Anagrafica Personale',  icon: 'people',        path: '/secretary/users' },
                { label: 'Sostituzioni Docenti',  icon: 'swap_horiz',    path: '/secretary/substitutions' },
            ]
        },
        {
            category: 'Atti & Amministrazione',
            icon: 'folder_shared',
            children: [
                { label: 'Sportello Personale',    icon: 'forward_to_inbox', path: '/ata/personnel-desk' },
                { label: 'Documenti & Atti',       icon: 'description',      path: '/secretary/documents' },
                { label: 'Verbali & Riunioni',     icon: 'gavel',            path: '/secretary/verbali' },
                { label: 'Flussi SIDI',            icon: 'cloud_sync',       path: '/secretary/sidi' },
                { label: 'Comunicazioni',          icon: 'email',            path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo', icon: 'settings',         path: '/ata/settings' },
            ]
        }
    ]

    // ASSISTENTE AMMINISTRATIVO 
    const assistenteAmministrativoMenu = [
        { label: 'Dashboard ATA',      icon: 'dashboard',  path: '/ata',            exact: true },
        { label: 'Presenze Personale', icon: 'co_present', path: '/ata/attendance', quickAccess: true, badge: 'absentStaff' },
        {
            category: 'Presenze & Personale',
            icon: 'co_present',
            children: [
                { label: 'Presenze & Timbrature', icon: 'co_present',    path: '/ata/attendance' },
                { label: 'Scioperi',              icon: 'campaign',      path: '/ata/strike' },
                { label: 'Cartellino & Ferie',    icon: 'calendar_month', path: '/ata/timecard' },
                { label: 'Anagrafica Utenti',     icon: 'people',        path: '/secretary/users' },
                { label: 'Studenti',              icon: 'school',        path: '/secretary/students' },
            ]
        },
        {
            category: 'Segreteria & Atti',
            icon: 'folder_shared',
            children: [
                { label: 'Sportello Personale',    icon: 'forward_to_inbox', path: '/ata/personnel-desk' },
                { label: 'Flussi SIDI',            icon: 'cloud_sync',       path: '/secretary/sidi' },
                { label: 'Documenti',              icon: 'description',      path: '/secretary/documents' },
                { label: 'Certificati',            icon: 'workspace_premium', path: '/secretary/certificates' },
                { label: 'Verbali & Riunioni',     icon: 'gavel',            path: '/secretary/verbali' },
                { label: 'Comunicazioni',          icon: 'email',            path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo', icon: 'settings',         path: '/ata/settings' },
            ]
        }
    ]

    // ASSISTENTE ALUNNI
    const assistenteAlunniMenu = [
        { label: 'Dashboard ATA',       icon: 'dashboard',         path: '/ata',                    exact: true },
        { label: 'Studenti & Fascicoli', icon: 'school',           path: '/secretary/students',     quickAccess: true },
        { label: 'Certificati Alunni',  icon: 'workspace_premium', path: '/secretary/certificates', quickAccess: true },
        {
            category: 'Gestione Alunni',
            icon: 'school',
            children: [
                { label: 'Studenti & Fascicoli', icon: 'school',              path: '/secretary/students' },
                { label: 'Certificati Alunni',   icon: 'workspace_premium',   path: '/secretary/certificates' },
                { label: 'Flussi SIDI Alunni',   icon: 'cloud_sync',          path: '/secretary/sidi' },
                { label: 'Scrutinio & Esami',    icon: 'analytics',           path: '/secretary/scrutiny' },
                { label: 'Libri di Testo',       icon: 'auto_stories',        path: '/secretary/textbooks' },
                { label: 'Comunicazioni',        icon: 'email',               path: '/secretary/communications', badge: 'unreadMessages' },
            ]
        }
    ]

    // ASSISTENTE PERSONALE 
    const assistentePersonaleMenu = [
        { label: 'Dashboard ATA',      icon: 'dashboard',  path: '/ata',            exact: true },
        { label: 'Presenze Personale', icon: 'co_present', path: '/ata/attendance', quickAccess: true, badge: 'absentStaff' },
        {
            category: 'Gestione Personale',
            icon: 'co_present',
            children: [
                { label: 'Anagrafica Personale',  icon: 'people',           path: '/secretary/users' },
                { label: 'Presenze & Timbrature', icon: 'co_present',       path: '/ata/attendance' },
                { label: 'Scioperi',              icon: 'campaign',         path: '/ata/strike' },
                { label: 'Cartellino & Ferie',    icon: 'calendar_month',   path: '/ata/timecard' },
                { label: 'Sostituzioni Docenti',  icon: 'swap_horiz',       path: '/secretary/substitutions', badge: 'pendingSubstitutions' },
                { label: 'Sportello Personale',   icon: 'forward_to_inbox', path: '/ata/personnel-desk' },
                { label: 'Comunicazioni',         icon: 'email',            path: '/secretary/communications', badge: 'unreadMessages' },
            ]
        }
    ]

    // ASSISTENTE CONTABILITA 
    const assistenteContabilitaMenu = [
        { label: 'Dashboard ATA', icon: 'dashboard', path: '/ata', exact: true },
        {
            category: 'Bilancio & Contabilita',
            icon: 'account_balance',
            children: [
                { label: 'Documenti & Mandati',    icon: 'description',       path: '/secretary/documents' },
                { label: 'Sportello Personale',    icon: 'forward_to_inbox',  path: '/ata/personnel-desk' },
                { label: 'Certificati & Ricevute', icon: 'workspace_premium', path: '/secretary/certificates' },
                { label: 'Comunicazioni',          icon: 'email',             path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo', icon: 'settings',          path: '/ata/settings' },
            ]
        }
    ]

    // ASSISTENTE PROTOCOLLO 
    const assistenteProtocolloMenu = [
        { label: 'Dashboard ATA',       icon: 'dashboard',  path: '/ata',                 exact: true },
        { label: 'Registro Protocollo', icon: 'description', path: '/secretary/documents', quickAccess: true },
        {
            category: 'Protocollo & Archivi',
            icon: 'mark_email_read',
            children: [
                { label: 'Registro Protocollo',       icon: 'description', path: '/secretary/documents' },
                { label: 'Verbali & Delibere',        icon: 'gavel',       path: '/secretary/verbali' },
                { label: 'Comunicazioni & Circolari', icon: 'email',       path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Riunioni',                  icon: 'groups',      path: '/secretary/meetings' },
                { label: 'Impostazioni & Profilo',    icon: 'settings',    path: '/ata/settings' },
            ]
        }
    ]

    // ASSISTENTE SPORTELLO 
    const assistenteSportelloMenu = [
        { label: 'Dashboard ATA',      icon: 'dashboard',         path: '/ata',                     exact: true },
        { label: 'Sportello Utenza',   icon: 'forward_to_inbox',  path: '/ata/personnel-desk',      quickAccess: true },
        { label: 'Rilascio Certificati', icon: 'workspace_premium', path: '/secretary/certificates', quickAccess: true },
        {
            category: 'Front-Office & Sportello',
            icon: 'support_agent',
            children: [
                { label: 'Sportello Utenza',       icon: 'forward_to_inbox',  path: '/ata/personnel-desk' },
                { label: 'Registro Visitatori',    icon: 'door_front',        path: '/ata/visitor-registry' },
                { label: 'Rilascio Certificati',   icon: 'workspace_premium', path: '/secretary/certificates' },
                { label: 'Comunicazioni',          icon: 'email',             path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo', icon: 'settings',          path: '/ata/settings' },
            ]
        }
    ]

    // ASSISTENTE TECNICO
    const assistenteTecnicoMenu = [
        { label: 'Dashboard ATA',      icon: 'dashboard',      path: '/ata',           exact: true },
        { label: 'Cartellino Presenze', icon: 'calendar_month', path: '/ata/timecard',  quickAccess: true },
        {
            category: 'Laboratori & Tecnologie',
            icon: 'computer',
            children: [
                { label: 'Cartellino & Ferie',        icon: 'calendar_month',   path: '/ata/timecard' },
                { label: 'E-Learning (Google/Teams)', icon: 'hub',              path: '/admin/elearning' },
                { label: 'Sportello Tecnico',         icon: 'forward_to_inbox', path: '/ata/personnel-desk' },
                { label: 'Comunicazioni',             icon: 'email',            path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo',    icon: 'settings',         path: '/ata/settings' },
            ]
        }
    ]

    // COLLABORATORE DS
    const collaboratoreDsMenu = [
        { label: 'Dashboard ATA',    icon: 'dashboard',  path: '/ata',                          exact: true },
        { label: 'Emergenza Sost.',  icon: 'bolt',       path: '/ata/emergency-substitutions',  quickAccess: true, badge: 'pendingSubstitutions' },
        { label: 'Presenze Personale', icon: 'co_present', path: '/ata/attendance',             quickAccess: true, badge: 'absentStaff' },
        {
            category: 'Presenze & Organizzazione',
            icon: 'co_present',
            children: [
                { label: 'Emergenza Sostituzioni', icon: 'bolt',          path: '/ata/emergency-substitutions' },
                { label: 'Presenze Personale',     icon: 'co_present',    path: '/ata/attendance' },
                { label: 'Scioperi',               icon: 'campaign',      path: '/ata/strike' },
                { label: 'Cartellino & Ferie',     icon: 'calendar_month', path: '/ata/timecard' },
                { label: 'Sostituzioni Docenti',   icon: 'swap_horiz',    path: '/secretary/substitutions' },
                { label: 'Orario Scolastico',      icon: 'schedule',      path: '/secretary/timetable' },
                { label: 'Aule & Plessi',          icon: 'meeting_room',  path: '/secretary/rooms' },
            ]
        },
        {
            category: 'Comunicazioni & Atti',
            icon: 'campaign',
            children: [
                { label: 'Sportello Personale',    icon: 'forward_to_inbox', path: '/ata/personnel-desk' },
                { label: 'Verbali & Riunioni',     icon: 'gavel',            path: '/secretary/verbali' },
                { label: 'Comunicazioni',          icon: 'email',            path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Documenti',              icon: 'description',      path: '/secretary/documents' },
                { label: 'Impostazioni & Profilo', icon: 'settings',         path: '/ata/settings' },
            ]
        }
    ]

    // COLLABORATORE SCOLASTICO 
    const collaboratoreScolasticoMenu = [
        { label: 'Dashboard ATA',       icon: 'dashboard',      path: '/ata',                  exact: true },
        { label: 'Registro Visitatori', icon: 'door_front',     path: '/ata/visitor-registry', quickAccess: true },
        { label: 'Cartellino & Ferie',  icon: 'calendar_month', path: '/ata/timecard',         quickAccess: true },
        {
            category: 'Servizi di Sede & Vigilanza',
            icon: 'security',
            children: [
                { label: 'Registro Visitatori',       icon: 'door_front',       path: '/ata/visitor-registry' },
                { label: 'Cartellino & Ferie',        icon: 'calendar_month',   path: '/ata/timecard' },
                { label: 'Sportello Personale',       icon: 'forward_to_inbox', path: '/ata/personnel-desk' },
                { label: 'Comunicazioni & Circolari', icon: 'email',            path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo',    icon: 'settings',         path: '/ata/settings' },
            ]
        }
    ]

    // RESPONSABILE SERVIZIO 
    const responsabileServizioMenu = [
        { label: 'Dashboard ATA', icon: 'dashboard', path: '/ata', exact: true },
        {
            category: 'Gestione Servizio & Struttura',
            icon: 'room_preferences',
            children: [
                { label: 'Registro Visitatori',    icon: 'door_front',       path: '/ata/visitor-registry' },
                { label: 'Cartellino Presenze',    icon: 'calendar_month',   path: '/ata/timecard' },
                { label: 'Sportello Personale',    icon: 'forward_to_inbox', path: '/ata/personnel-desk' },
                { label: 'Comunicazioni',          icon: 'email',            path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo', icon: 'settings',         path: '/ata/settings' },
            ]
        }
    ]

    // RESPONSABILE GESTIONE DOCUMENTALE 
    const responsabileDocumentaleMenu = [
        { label: 'Dashboard', icon: 'dashboard', path: '/secretary/documents', exact: true },
        {
            category: 'Gestione Documentale & Archivi',
            icon: 'archive',
            children: [
                { label: 'Documenti & Atti',          icon: 'description', path: '/secretary/documents' },
                { label: 'Verbali & Riunioni',        icon: 'gavel',       path: '/secretary/verbali' },
                { label: 'Comunicazioni & Circolari', icon: 'email',       path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo',    icon: 'settings',    path: '/ata/settings' },
            ]
        }
    ]

    // RESPONSABILE CONSERVAZIONE 
    const responsabileConservazioneMenu = [
        { label: 'Dashboard', icon: 'dashboard', path: '/secretary/documents', exact: true },
        {
            category: 'Conservazione Digitale',
            icon: 'inventory_2',
            children: [
                { label: 'Documenti & Archivio',   icon: 'description', path: '/secretary/documents' },
                { label: 'Verbali & Delibere',     icon: 'gavel',       path: '/secretary/verbali' },
                { label: 'Comunicazioni',          icon: 'email',       path: '/secretary/communications', badge: 'unreadMessages' },
                { label: 'Impostazioni & Profilo', icon: 'settings',    path: '/ata/settings' },
            ]
        }
    ]

    // DPO
    const dpoMenu = [
        { label: 'Dashboard',              icon: 'dashboard', path: '/admin/audit-logs', exact: true },
        { label: 'Audit Logs & Sicurezza', icon: 'history',   path: '/admin/audit-logs', quickAccess: true },
        {
            category: 'Privacy & Sicurezza',
            icon: 'security',
            children: [
                { label: 'Audit Logs & Sicurezza', icon: 'history',       path: '/admin/audit-logs' },
                { label: 'Monitoraggio Sistema',   icon: 'monitor_heart', path: '/admin/monitoring' },
                { label: 'Impostazioni Privacy',   icon: 'security',      path: '/admin/settings' },
                { label: 'Comunicazioni',          icon: 'email',         path: '/secretary/communications' },
            ]
        }
    ]

    // ROUTING PER RUOLO

    if (normRole === 'superadmin' || normRole === 'system_auditor') {
        return superadminMenu
    }

    if (normRole === 'admin') {
        return adminMenu
    }

    if (normRole === 'principal' || normRole === 'dirigente_scolastico' || normRole === 'staff') {
        return secretaryMenu
    }

    if (normRole === 'secretary') {
        return secretaryMenu
    }

    if (normRole === 'vice_principal' || normRole === 'collaboratore_vicario') {
        return vicePrincipalMenu
    }

    if (normRole === 'student' || normRole === 'studente') {
        return studentMenu
    }

    if (normRole === 'parent' || normRole === 'genitore') {
        return parentMenu
    }

    if (normRole === 'dsga') {
        return dsgaMenu
    }

    if (normRole === 'assistente_amministrativo') {
        return assistenteAmministrativoMenu
    }

    if (normRole === 'assistente_alunni') {
        return assistenteAlunniMenu
    }

    if (normRole === 'assistente_personale') {
        return assistentePersonaleMenu
    }

    if (normRole === 'assistente_contabilita') {
        return assistenteContabilitaMenu
    }

    if (normRole === 'assistente_protocollo') {
        return assistenteProtocolloMenu
    }

    if (normRole === 'assistente_sportello') {
        return assistenteSportelloMenu
    }

    if (normRole === 'assistente_tecnico') {
        return assistenteTecnicoMenu
    }

    if (normRole === 'collaboratore_ds') {
        return collaboratoreDsMenu
    }

    if (normRole === 'collaboratore_scolastico') {
        return collaboratoreScolasticoMenu
    }

    if (normRole === 'responsabile_servizio') {
        return responsabileServizioMenu
    }

    if (normRole === 'responsabile_gestione_documentale') {
        return responsabileDocumentaleMenu
    }

    if (normRole === 'responsabile_conservazione') {
        return responsabileConservazioneMenu
    }

    if (normRole === 'dpo') {
        return dpoMenu
    }

    // TEACHER (con incarichi dinamici) 
    if (normRole === 'teacher' || normRole === 'docente' || normRole === 'coordinator' || normRole === 'coordinatore_classe') {
        const baseTeacherMenu = JSON.parse(JSON.stringify(teacherMenu))
        const activeAssignments = (assignments || []).filter(a => a.is_active !== false)
        const hasAssignment = (type) => activeAssignments.some(a => a.assignment_type === type)
        const isCoord =
            normRole === 'coordinator' ||
            normRole === 'coordinatore_classe' ||
            hasAssignment('coordinatore_classe') ||
            hasAssignment('coordinator')

        // Sblocca le voci coordinatorOnly se il docente Ã¨ coordinatore
        if (isCoord) {
            baseTeacherMenu.forEach(cat => {
                if (cat.children) {
                    cat.children.forEach(child => {
                        if (child.coordinatorOnly) {
                            child.coordinatorOnly = false
                        }
                    })
                }
            })
        }

        // Referente Progetti PTOF/PNRR
        if (hasAssignment('referente_progetto')) {
            const orgCat = baseTeacherMenu.find(c => c.category === 'Orario & Agenda')
            if (orgCat && !orgCat.children.some(c => c.path === '/secretary/pcto')) {
                orgCat.children.push({ label: 'Progetti & Finanziamenti', icon: 'rocket_launch', path: '/secretary/pcto' })
            }
        }

        // Tutor Orientamento
        if (hasAssignment('tutor_orientamento')) {
            const regCat = baseTeacherMenu.find(c => c.category === 'Registro & Classi')
            if (regCat && !regCat.children.some(c => c.path === '/student/orientamento')) {
                regCat.children.push({ label: 'Tutor Orientamento', icon: 'explore', path: '/student/orientamento' })
            }
        }

        // Referente Inclusione BES/DSA
        if (hasAssignment('referente_inclusione')) {
            const regCat = baseTeacherMenu.find(c => c.category === 'Registro & Classi')
            if (regCat && !regCat.children.some(c => c.label === 'Inclusione (BES / DSA)')) {
                regCat.children.push({ label: 'Inclusione (BES / DSA)', icon: 'accessibility_new', path: '/teacher/pdp' })
            }
        }

        // Animatore Digitale
        if (hasAssignment('animatore_digitale')) {
            const commCat = baseTeacherMenu.find(c => c.category === 'Comunicazioni')
            if (commCat && !commCat.children.some(c => c.path === '/admin/elearning')) {
                commCat.children.push({ label: 'Team Digitale & E-Learning', icon: 'hub', path: '/admin/elearning' })
            }
        }

        return baseTeacherMenu
    }

    return []
}

export default useMenuItems
