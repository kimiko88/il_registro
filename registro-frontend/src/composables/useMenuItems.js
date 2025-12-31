/**
 * Menu configuration for different user roles
 * Returns menu items based on user role
 */
export function useMenuItems(role) {
    const menuConfig = {
        superadmin: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            { label: 'Gestione Scuole', icon: 'school', path: '/admin/schools' },
            { label: 'Gestione Admin', icon: 'admin_panel_settings', path: '/admin/admins' },
            { label: 'Monitoraggio Sistema', icon: 'monitor_heart', path: '/admin/monitoring' },
            { label: 'Analytics Globali', icon: 'analytics', path: '/admin/analytics' },
            { label: 'Audit Logs', icon: 'history', path: '/admin/audit-logs' },
            { label: 'Impostazioni', icon: 'settings', path: '/admin/settings' }
        ],
        admin: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            { label: 'La Mia Scuola', icon: 'school', path: '/admin/schools' },
            { label: 'Analytics', icon: 'analytics', path: '/admin/analytics' },
            { label: 'Impostazioni', icon: 'settings', path: '/settings' }
        ],
        secretary: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            { label: 'Documenti', icon: 'description', path: '/secretary/documents' },
            { label: 'Studenti', icon: 'school', path: '/secretary/students' },
            { label: 'Utenti', icon: 'people', path: '/secretary/users' },
            { label: 'Comunicazioni', icon: 'email', path: '/secretary/communications' },
            { label: 'Report', icon: 'assessment', path: '/secretary/reports' },
            { label: 'Impostazioni', icon: 'settings', path: '/secretary/settings' }
        ],
        teacher: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            { label: 'Le Mie Classi', icon: 'class', path: '/teacher/classes' },
            { label: 'Voti', icon: 'grade', path: '/teacher/grades' },
            { label: 'Presenze', icon: 'how_to_reg', path: '/teacher/attendance' },
            { label: 'Documenti', icon: 'description', path: '/teacher/documents' },
            { label: 'Colloqui', icon: 'event', path: '/teacher/colloqui' },
            { label: 'Comunicazioni', icon: 'email', path: '/teacher/communications' }
        ],
        student: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            { label: 'I Miei Voti', icon: 'grade', path: '/student/grades' },
            { label: 'Le Mie Presenze', icon: 'how_to_reg', path: '/student/attendance' },
            { label: 'Documenti', icon: 'description', path: '/student/documents' },
            { label: 'PCTO', icon: 'work', path: '/student/pcto' },
            { label: 'Orientamento', icon: 'explore', path: '/student/orientamento' },
            { label: 'Comunicazioni', icon: 'email', path: '/student/communications' },
            { label: 'Profilo', icon: 'person', path: '/student/profile' }
        ],
        parent: [
            { label: 'Dashboard', icon: 'dashboard', path: '/', exact: true },
            { label: 'I Miei Figli', icon: 'family_restroom', path: '/parent/children' },
            { label: 'Voti', icon: 'grade', path: '/parent/grades' },
            { label: 'Presenze', icon: 'how_to_reg', path: '/parent/attendance' },
            { label: 'Colloqui', icon: 'event', path: '/parent/colloqui' },
            { label: 'Comunicazioni', icon: 'email', path: '/parent/communications' },
            { label: 'Supporto', icon: 'help', path: '/parent/support' },
            { label: 'Profilo', icon: 'person', path: '/parent/profile' }
        ]
    }

    return menuConfig[role] || []
}
