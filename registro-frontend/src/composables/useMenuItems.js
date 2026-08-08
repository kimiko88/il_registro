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
            { label: 'I Miei Voti', icon: 'grade', path: '/student/grades' },
            { label: 'Le Mie Presenze', icon: 'how_to_reg', path: '/student/attendance' },
            { label: 'Note Disciplinari', icon: 'assignment_late', path: '/student/notes' },
            { label: 'Compiti', icon: 'assignment', path: '/student/homework' },
            { label: 'Materiale Didattico', icon: 'folder_shared', path: '/student/didactics' },
            { label: 'Orario Lezioni', icon: 'schedule', path: '/student/timetable' },
            { label: 'Documenti', icon: 'description', path: '/student/documents' },
            { label: 'PCTO', icon: 'work', path: '/student/pcto' },
            { label: 'Orientamento', icon: 'explore', path: '/student/orientamento' },
            { label: 'Comunicazioni', icon: 'email', path: '/student/communications' },
            { label: 'Calendario Scolastico', icon: 'calendar_month', path: '/student/school-calendar' },
            { label: 'Pagella', icon: 'description', path: '/student/report-card' },
            { label: 'Profilo', icon: 'person', path: '/student/profile' }
        ],
        parent: [
            { label: 'Dashboard', icon: 'dashboard', path: '/parent', exact: true },
            { label: 'I Miei Figli', icon: 'family_restroom', path: '/parent/children' },
            { label: 'Voti', icon: 'grade', path: '/parent/grades' },
            { label: 'Pagella', icon: 'description', path: '/parent/report-card' },
            { label: 'Piano PDP / PEI', icon: 'accessibility_new', path: '/parent/pdp' },
            { label: 'Presenze', icon: 'how_to_reg', path: '/parent/attendance' },
            { label: 'Note Disciplinari', icon: 'assignment_late', path: '/parent/notes' },
            { label: 'Colloqui', icon: 'event', path: '/parent/colloqui' },
            { label: 'Documenti', icon: 'description', path: '/parent/documents' },
            { label: 'Materiale Didattico', icon: 'folder_shared', path: '/parent/didactics' },
            { label: 'Orario Lezioni', icon: 'schedule', path: '/parent/timetable' },
            { label: 'Comunicazioni', icon: 'email', path: '/parent/communications' },
            { label: 'Profilo', icon: 'person', path: '/parent/profile' },
            { label: 'Supporto', icon: 'help', path: '/support' }
        ]
    }

    return menuConfig[role] || []
}
