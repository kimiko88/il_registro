/**
 * Menu configuration for different user roles
 * Returns menu items based on user role
 */
export function useMenuItems(role) {
    const menuConfig = {
        superadmin: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            { label: 'Gestione Scuole', icon: 'school', path: '/admin/schools' },
            { label: 'Gestione Utenti', icon: 'people', path: '/admin/users' },
            { label: 'Gestione Admin', icon: 'admin_panel_settings', path: '/admin/admins' },
            { label: 'Monitoraggio Sistema', icon: 'monitor_heart', path: '/admin/monitoring' },
            { label: 'Analytics Globali', icon: 'analytics', path: '/admin/analytics' },
            { label: 'Audit Logs', icon: 'history', path: '/admin/audit-logs' },
            { label: 'Impostazioni', icon: 'settings', path: '/admin/settings' },
            { label: 'Supporto', icon: 'help', path: '/support' }
        ],
        admin: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            { label: 'La Mia Scuola', icon: 'school', path: '/admin/schools' },
            { label: 'Gestione Utenti', icon: 'people', path: '/admin/users' },
            { label: 'Gestione Sostituzioni', icon: 'swap_horiz', path: '/secretary/substitutions' },
            { label: 'Feature Flags & Istituto', icon: 'toggle_on', path: '/admin/school-settings' },
            { label: 'Analytics', icon: 'analytics', path: '/admin/analytics' },
            { label: 'Google & Teams E-Learning', icon: 'hub', path: '/admin/elearning' },
            { label: 'Impostazioni', icon: 'settings', path: '/admin/settings' }
        ],
        secretary: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            {
                category: 'Anagrafiche & Classi',
                icon: 'school',
                children: [
                    { label: 'Studenti', icon: 'school', path: '/secretary/students' },
                    { label: 'Utenti', icon: 'people', path: '/secretary/users' },
                    { label: 'Classi', icon: 'room', path: '/secretary/classes' },
                    { label: 'Orario Scolastico', icon: 'schedule', path: '/secretary/timetable' },
                    { label: 'Gruppi Linguistici / Articolati', icon: 'groups', path: '/secretary/groups' }
                ]
            },
            {
                category: 'Atti & Certificati',
                icon: 'folder_shared',
                children: [
                    { label: 'Documenti', icon: 'description', path: '/secretary/documents' },
                    { label: 'Certificati', icon: 'workspace_premium', path: '/secretary/certificates' },
                    { label: 'Libri di Testo', icon: 'auto_stories', path: '/secretary/textbooks' },
                    { label: 'Riunioni', icon: 'groups', path: '/secretary/meetings' }
                ]
            },
            {
                category: 'Servizi & Report',
                icon: 'manage_accounts',
                children: [
                    { label: 'Gestione Sostituzioni', icon: 'swap_horiz', path: '/secretary/substitutions' },
                    { label: 'Comunicazioni', icon: 'email', path: '/secretary/communications' },
                    { label: 'Report', icon: 'assessment', path: '/secretary/reports' },
                    { label: 'PCTO', icon: 'work', path: '/secretary/pcto' },
                    { label: 'Scrutinio', icon: 'analytics', path: '/secretary/scrutiny' },
                    { label: 'Impostazioni', icon: 'settings', path: '/secretary/settings' }
                ]
            }
        ],
        teacher: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            {
                category: 'Didattica & Valutazione',
                icon: 'menu_book',
                children: [
                    { label: 'Le Mie Classi', icon: 'class', path: '/teacher/classes' },
                    { label: 'Registro Classe', icon: 'menu_book', path: '/teacher/lessons' },
                    { label: 'Programmazione UdA', icon: 'auto_stories', path: '/teacher/uda' },
                    { label: 'Valutazione Competenze', icon: 'stars', path: '/teacher/competencies' },
                    { label: 'Voti', icon: 'grade', path: '/teacher/grades' },
                    { label: 'Presenze', icon: 'how_to_reg', path: '/teacher/attendance' },
                    { label: 'Didattica', icon: 'folder_shared', path: '/teacher/didactics' },
                    { label: 'Scrutinio', icon: 'analytics', path: '/teacher/scrutiny', coordinatorOnly: true },
                    { label: 'Piani PDP / PEI', icon: 'accessibility_new', path: '/teacher/pdp' },
                    { label: 'Rubriche Valutative', icon: 'fact_check', path: '/teacher/rubrics' },
                    { label: 'Coordinamento', icon: 'co_present', path: '/teacher/coordinator', coordinatorOnly: true }
                ]
            },
            {
                category: 'Organizzazione & Orario',
                icon: 'event',
                children: [
                    { label: 'Orario Lezioni', icon: 'schedule', path: '/teacher/timetable' },
                    { label: 'Agenda', icon: 'edit_calendar', path: '/teacher/agenda' },
                    { label: 'Colloqui', icon: 'event', path: '/teacher/colloqui' },
                    { label: 'Sostituzioni', icon: 'swap_horiz', path: '/teacher/substitutions' },
                    { label: 'Verbali', icon: 'gavel', path: '/teacher/verbali' }
                ]
            },
            {
                category: 'Comunicazioni & Atti',
                icon: 'campaign',
                children: [
                    { label: 'Comunicazioni', icon: 'email', path: '/teacher/communications' },
                    { label: 'Documenti', icon: 'description', path: '/teacher/documents' },
                    { label: 'Note Disciplinari', icon: 'assignment_late', path: '/teacher/notes' },
                    { label: 'Impostazioni', icon: 'settings', path: '/teacher/settings' }
                ]
            }
        ],
        student: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            {
                category: 'Didattica & Valutazione',
                icon: 'school',
                children: [
                    { label: 'I Miei Voti', icon: 'grade', path: '/student/grades' },
                    { label: 'Le Mie Presenze', icon: 'how_to_reg', path: '/student/attendance' },
                    { label: 'Note Disciplinari', icon: 'assignment_late', path: '/student/notes' },
                    { label: 'Compiti', icon: 'assignment', path: '/student/homework' },
                    { label: 'Materiale Didattico', icon: 'folder_shared', path: '/student/didactics' },
                    { label: 'Pagella', icon: 'description', path: '/student/report-card' }
                ]
            },
            {
                category: 'Organizzazione & Orario',
                icon: 'event',
                children: [
                    { label: 'Orario Lezioni', icon: 'schedule', path: '/student/timetable' },
                    { label: 'Agenda', icon: 'edit_calendar', path: '/student/agenda' },
                    { label: 'Calendario Scolastico', icon: 'calendar_month', path: '/student/school-calendar' },
                    { label: 'Obiettivi', icon: 'flag', path: '/student/goals' }
                ]
            },
            {
                category: 'Percorsi & Comunicazioni',
                icon: 'campaign',
                children: [
                    { label: 'PCTO', icon: 'work', path: '/student/pcto' },
                    { label: 'Orientamento', icon: 'explore', path: '/student/orientamento' },
                    { label: 'Comunicazioni', icon: 'email', path: '/student/communications' },
                    { label: 'Documenti', icon: 'description', path: '/student/documents' },
                    { label: 'Impostazioni', icon: 'settings', path: '/student/settings' },
                    { label: 'Profilo', icon: 'person', path: '/student/profile' }
                ]
            }
        ],
        parent: [
            { label: 'Dashboard', icon: 'dashboard', path: '/parent', exact: true },
            {
                category: 'Valutazione & Didattica',
                icon: 'school',
                children: [
                    { label: 'I Miei Figli', icon: 'family_restroom', path: '/parent/children' },
                    { label: 'Voti', icon: 'grade', path: '/parent/grades' },
                    { label: 'Pagella', icon: 'description', path: '/parent/report-card' },
                    { label: 'Presenze', icon: 'how_to_reg', path: '/parent/attendance' },
                    { label: 'Note Disciplinari', icon: 'assignment_late', path: '/parent/notes' },
                    { label: 'Didattica', icon: 'folder_shared', path: '/parent/didactics' },
                    { label: 'Piano PDP / PEI', icon: 'accessibility_new', path: '/parent/pdp' }
                ]
            },
            {
                category: 'Servizi & Orari',
                icon: 'event',
                children: [
                    { label: 'Orario Lezioni', icon: 'schedule', path: '/parent/timetable' },
                    { label: 'Colloqui', icon: 'event', path: '/parent/colloqui' },
                    { label: 'Uscite & Viaggi', icon: 'card_travel', path: '/parent/trips' },
                    { label: 'Pagamenti', icon: 'payments', path: '/parent/payments' },
                    { label: 'Assemblee & Riunioni', icon: 'groups', path: '/parent/meetings' }
                ]
            },
            {
                category: 'Comunicazioni & Account',
                icon: 'manage_accounts',
                children: [
                    { label: 'Comunicazioni', icon: 'email', path: '/parent/communications' },
                    { label: 'Documenti', icon: 'description', path: '/parent/documents' },
                    { label: 'Impostazioni', icon: 'settings', path: '/parent/settings' },
                    { label: 'Profilo', icon: 'person', path: '/parent/profile' },
                    { label: 'Supporto', icon: 'help', path: '/support' }
                ]
            }
        ]
    }

    const normRole = (role || '').toLowerCase()
    if (normRole === 'principal' || normRole === 'vice_principal') {
        return menuConfig.secretary
    }
    if (normRole === 'coordinator' || normRole === 'docente') {
        return menuConfig.teacher
    }
    if (normRole === 'system_auditor') {
        return menuConfig.superadmin
    }
    if (normRole === 'staff') {
        return menuConfig.secretary
    }

    return menuConfig[normRole] || []
}

export default useMenuItems
