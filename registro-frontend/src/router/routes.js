import { useAuthStore } from '@/stores/auth'

export default [
    {
        path: '/',
        component: () => import('@/layouts/MainLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Dashboard.vue'), meta: { title: 'Dashboard', roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent'] } },
            { path: 'dashboard', redirect: '/' },

            // Admin Routes (SuperAdmin + Admin)
            {
                path: 'admin',
                redirect: '/admin/dashboard'
            },
            {
                path: 'admin/dashboard',
                component: () => import('@/pages/admin/Dashboard.vue'),
                meta: { title: 'Dashboard Amministratore', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/schools',
                component: () => import('@/pages/admin/SchoolManagement.vue'),
                meta: { title: 'Gestione Scuole', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/schools/:id',
                component: () => import('@/pages/admin/SchoolDetail.vue'),
                meta: { title: 'Dettaglio Scuola', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/admins',
                component: () => import('@/pages/admin/AdminUsers.vue'),
                meta: { title: 'Gestione Amministratori', roles: ['superadmin'] }
            },
            {
                path: 'admin/monitoring',
                component: () => import('@/pages/admin/Monitoring.vue'),
                meta: { title: 'Monitoraggio Sistema', roles: ['superadmin'] }
            },
            {
                path: 'admin/users',
                component: () => import('@/pages/secretary/Users.vue'),
                meta: { title: 'Gestione Utenti System', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/analytics',
                component: () => import('@/pages/admin/Analytics.vue'),
                meta: { title: 'Statistiche e Analytics', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/elearning',
                component: () => import('@/pages/admin/ElearningIntegration.vue'),
                meta: { title: 'Google Classroom & Teams', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/settings',
                component: () => import('@/pages/admin/Settings.vue'),
                meta: { title: 'Impostazioni di Sistema', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/scheduler',
                component: () => import('@/pages/admin/Scheduler.vue'),
                meta: { title: 'Pianificazione Task', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/school-settings',
                component: () => import('@/pages/admin/SchoolSettings.vue'),
                meta: { title: 'Configurazione Scuola', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/audit-logs',
                component: () => import('@/pages/admin/AuditLog.vue'),
                meta: { title: 'Registro Eventi & Audit', roles: ['superadmin'] }
            },
            {
                path: 'admin/tenants',
                component: () => import('@/pages/admin/Tenants.vue'),
                meta: { title: 'Gestione Multi-Tenant', roles: ['superadmin'] }
            },

            // Secretary Routes
            { path: 'secretary', component: () => import('@/pages/secretary/Index.vue'), meta: { title: 'Pannello Segreteria', roles: ['secretary'] } },
            { path: 'secretary/documents', component: () => import('@/pages/secretary/Documents.vue'), meta: { title: 'Gestione Documenti', roles: ['secretary'] } },
            { path: 'secretary/users', component: () => import('@/pages/secretary/Users.vue'), meta: { title: 'Anagrafica Utenti', roles: ['secretary'] } },
            { path: 'secretary/students', component: () => import('@/pages/secretary/Students.vue'), meta: { title: 'Anagrafica Studenti', roles: ['secretary'] } },
            { path: 'secretary/communications', component: () => import('@/pages/secretary/Communications.vue'), meta: { title: 'Circolari & Comunicazioni', roles: ['secretary'] } },
            { path: 'secretary/reports', component: () => import('@/pages/secretary/Reports.vue'), meta: { title: 'Reportistica Scolastica', roles: ['secretary'] } },
            { path: 'secretary/pcto', component: () => import('@/pages/secretary/PCTO.vue'), meta: { title: 'Gestione PCTO', roles: ['secretary'] } },
            { path: 'secretary/textbooks', component: () => import('@/pages/secretary/Textbooks.vue'), meta: { title: 'Adozione Libri di Testo', roles: ['secretary'] } },
            { path: 'secretary/settings', component: () => import('@/pages/secretary/Settings.vue'), meta: { title: 'Impostazioni Segreteria', roles: ['secretary'] } },
            { path: 'secretary/classes', component: () => import('@/pages/secretary/Classes.vue'), meta: { title: 'Gestione Classi', roles: ['secretary'] } },
            { path: 'secretary/scrutiny', component: () => import('@/pages/secretary/Scrutiny.vue'), meta: { title: 'Scrutini Scolastici', roles: ['secretary'] } },
            { path: 'secretary/groups', component: () => import('@/pages/secretary/Groups.vue'), meta: { title: 'Gruppi Linguistici / Articolati', roles: ['secretary'] } },
            { path: 'secretary/meetings', component: () => import('@/pages/secretary/Meetings.vue'), meta: { title: 'Organizzazione Riunioni', roles: ['secretary'] } },
            { path: 'secretary/certificates', component: () => import('@/pages/secretary/Certificates.vue'), meta: { title: 'Certificati & Attestati', roles: ['secretary', 'admin', 'superadmin'] } },
            { path: 'secretary/audit-log', component: () => import('@/pages/secretary/AuditLog.vue'), meta: { title: 'Audit Log Segreteria', roles: ['secretary', 'admin', 'superadmin'] } },
            { path: 'secretary/students/:id/fascicolo', component: () => import('@/pages/secretary/FascicoloStudente.vue'), meta: { title: 'Fascicolo Studente', roles: ['secretary'] } },

            // Teacher Routes
            { path: 'teacher', component: () => import('@/pages/teacher/Index.vue'), meta: { title: 'Pannello Docente', roles: ['teacher'] } },
            { path: 'teacher/grades', component: () => import('@/pages/teacher/Grades.vue'), meta: { title: 'Gestione Voti', roles: ['teacher'] } },
            { path: 'teacher/attendance', component: () => import('@/pages/teacher/Attendance.vue'), meta: { title: 'Registro Appello & Presenze', roles: ['teacher'] } },
            { path: 'teacher/classes', component: () => import('@/pages/teacher/Classes.vue'), meta: { title: 'Le Mie Classi', roles: ['teacher'] } },
            { path: 'teacher/groups', component: () => import('@/pages/teacher/Groups.vue'), meta: { title: 'Gruppi e Laboratori', roles: ['teacher'] } },
            { path: 'teacher/coordinator', component: () => import('@/pages/teacher/CoordinatorView.vue'), meta: { title: 'Pannello Coordinatore', roles: ['teacher'] } },
            { path: 'teacher/documents', component: () => import('@/pages/teacher/Documents.vue'), meta: { title: 'Documentazione Didattica', roles: ['teacher'] } },
            { path: 'teacher/colloqui', component: () => import('@/pages/teacher/Colloqui.vue'), meta: { title: 'Ricevimento Famiglie', roles: ['teacher'] } },
            { path: 'teacher/communications', component: () => import('@/pages/teacher/Communications.vue'), meta: { title: 'Comunicazioni Scuola', roles: ['teacher'] } },
            { path: 'teacher/lessons', component: () => import('@/components/Teacher/LessonPlanner.vue'), meta: { title: 'Registro Lezioni', roles: ['teacher'] } },
            { path: 'teacher/didactics', component: () => import('@/pages/teacher/Didactics.vue'), meta: { title: 'Materiale Didattico', roles: ['teacher'] } },
            { path: 'teacher/timetable', component: () => import('@/pages/teacher/Timetable.vue'), meta: { title: 'Orario Docente', roles: ['teacher'] } },
            { path: 'teacher/scrutiny', component: () => import('@/pages/teacher/Scrutiny.vue'), meta: { title: 'Gestione Scrutini', roles: ['teacher'] } },
            { path: 'teacher/verbali', component: () => import('@/pages/teacher/Verbali.vue'), meta: { title: 'Verbali Consiglio di Classe', roles: ['teacher'] } },
            { path: 'teacher/substitutions', component: () => import('@/pages/teacher/Substitutions.vue'), meta: { title: 'Sostituzioni Docenti', roles: ['teacher'] } },
            { path: 'teacher/grade-weights', component: () => import('@/pages/teacher/GradeWeights.vue'), meta: { title: 'Pesi e Criteri Valutazione', roles: ['teacher'] } },
            { path: 'teacher/agenda', component: () => import('@/pages/teacher/Agenda.vue'), meta: { title: 'Agenda di Classe', roles: ['teacher'] } },
            { path: 'teacher/notes', component: () => import('@/pages/teacher/Notes.vue'), meta: { title: 'Note & Richiami', roles: ['teacher'] } },
            { path: 'teacher/rubrics', component: () => import('@/pages/teacher/Rubrics.vue'), meta: { title: 'Rubriche Valutative', roles: ['teacher'] } },
            { path: 'teacher/pdp', component: () => import('@/pages/teacher/PdpPlans.vue'), meta: { title: 'Piani PDP / PEI', roles: ['teacher'] } },
            { path: 'teacher/uda', component: () => import('@/pages/teacher/UdaPlanner.vue'), meta: { title: 'Programmazione Didattica UdA', roles: ['teacher'] } },
            { path: 'teacher/competencies', component: () => import('@/pages/teacher/Competencies.vue'), meta: { title: 'Valutazione per Competenze', roles: ['teacher'] } },
            { path: 'teacher/settings', component: () => import('@/pages/teacher/Settings.vue'), meta: { title: 'Impostazioni Docente', roles: ['teacher'] } },

            // Student Routes
            { path: 'student', component: () => import('@/pages/student/Index.vue'), meta: { title: 'Pannello Studente', roles: ['student'] } },
            { path: 'student/grades', component: () => import('@/pages/student/Grades.vue'), meta: { title: 'I Miei Voti', roles: ['student'] } },
            { path: 'student/attendance', component: () => import('@/pages/student/Attendance.vue'), meta: { title: 'Presenze & Assenze', roles: ['student'] } },
            { path: 'student/documents', component: () => import('@/pages/student/Documents.vue'), meta: { title: 'Documenti Studente', roles: ['student'] } },
            { path: 'student/profile', component: () => import('@/pages/student/Profile.vue'), meta: { title: 'Profilo Studente', roles: ['student'] } },
            { path: 'student/pcto', component: () => import('@/pages/student/PCTO.vue'), meta: { title: 'Percorsi PCTO', roles: ['student'] } },
            { path: 'student/orientamento', component: () => import('@/pages/student/Orientamento.vue'), meta: { title: 'Orientamento Scolastico', roles: ['student'] } },
            { path: 'student/communications', component: () => import('@/pages/student/Communications.vue'), meta: { title: 'Comunicati Scolastici', roles: ['student'] } },
            { path: 'student/homework', component: () => import('@/pages/student/Homework.vue'), meta: { title: 'Compiti & Attività', roles: ['student'] } },
            { path: 'student/agenda', component: () => import('@/pages/student/AgendaCalendar.vue'), meta: { title: 'Calendario Agenda', roles: ['student'] } },
            { path: 'student/didactics', component: () => import('@/pages/student/Didactics.vue'), meta: { title: 'Materiali Didattici', roles: ['student'] } },
            { path: 'student/notes', component: () => import('@/pages/student/Notes.vue'), meta: { title: 'Note Disciplinari', roles: ['student'] } },
            { path: 'student/timetable', component: () => import('@/pages/student/Timetable.vue'), meta: { title: 'Orario delle Lezioni', roles: ['student'] } },
            { path: 'student/goals', component: () => import('@/pages/student/Goals.vue'), meta: { title: 'Obiettivi di Apprendimento', roles: ['student'] } },
            { path: 'student/school-calendar', component: () => import('@/pages/student/SchoolCalendar.vue'), meta: { title: 'Calendario Scolastico', roles: ['student'] } },
            { path: 'student/report-card', component: () => import('@/pages/student/ReportCard.vue'), meta: { title: 'Pagella Online', roles: ['student'] } },

            // Parent Routes
            { path: 'parent', component: () => import('@/pages/parent/Index.vue'), meta: { title: 'Pannello Famiglie', roles: ['parent'] } },
            { path: 'parent/children', component: () => import('@/pages/parent/Children.vue'), meta: { title: 'Figli Associati', roles: ['parent'] } },
            { path: 'parent/grades', component: () => import('@/pages/parent/Grades.vue'), meta: { title: 'Valutazioni Figli', roles: ['parent'] } },
            { path: 'parent/report-card', component: () => import('@/pages/parent/ReportCard.vue'), meta: { title: 'Pagella Scolastica', roles: ['parent'] } },
            { path: 'parent/attendance', component: () => import('@/pages/parent/Attendance.vue'), meta: { title: 'Presenze & Giustificazioni', roles: ['parent'] } },
            { path: 'parent/colloqui', component: () => import('@/pages/parent/Colloqui.vue'), meta: { title: 'Prenotazione Colloqui', roles: ['parent'] } },
            { path: 'parent/communications', component: () => import('@/pages/parent/Communications.vue'), meta: { title: 'Avvisi & Circolari', roles: ['parent'] } },
            { path: 'parent/trips', component: () => import('@/pages/parent/Trips.vue'), meta: { title: 'Uscite & Viaggi di Istruzione', roles: ['parent'] } },
            { path: 'parent/profile', component: () => import('@/pages/parent/Profile.vue'), meta: { title: 'Profilo Genitore', roles: ['parent'] } },
            { path: 'parent/didactics', component: () => import('@/pages/parent/Didactics.vue'), meta: { title: 'Didattica & Compiti', roles: ['parent'] } },
            { path: 'parent/notes', component: () => import('@/pages/parent/Notes.vue'), meta: { title: 'Note & Sanzioni', roles: ['parent'] } },
            { path: 'parent/timetable', component: () => import('@/pages/parent/Timetable.vue'), meta: { title: 'Orario Lezioni', roles: ['parent'] } },
            { path: 'parent/pdp', component: () => import('@/pages/parent/PdpView.vue'), meta: { title: 'Piano PDP / PEI', roles: ['parent'] } },
            { path: 'parent/documents', component: () => import('@/pages/parent/Documents.vue'), meta: { title: 'Documentazione & Moduli', roles: ['parent'] } },
            { path: 'parent/payments', component: () => import('@/pages/parent/Payments.vue'), meta: { title: 'Pagamenti Scolastici', roles: ['parent'] } },
            { path: 'parent/meetings', component: () => import('@/pages/parent/Meetings.vue'), meta: { title: 'Assemblee & Riunioni', roles: ['parent'] } },
            {
                path: 'communications',
                redirect: () => {
                    let role = ''
                    try {
                        const authStore = useAuthStore()
                        role = authStore.userRole || authStore.user?.role || ''
                    } catch {
                        // store fallback
                    }
                    if (role === 'teacher') return '/teacher/communications'
                    if (role === 'student') return '/student/communications'
                    if (role === 'parent') return '/parent/communications'
                    return '/secretary/communications'
                },
                meta: { title: 'Comunicazioni', roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent'] }
            },
            {
                path: 'profile',
                redirect: () => {
                    let role = ''
                    try {
                        const authStore = useAuthStore()
                        role = authStore.userRole || authStore.user?.role || ''
                    } catch {
                        // store fallback
                    }
                    if (role === 'student') return '/student/profile'
                    if (role === 'parent') return '/parent/profile'
                    if (role === 'admin' || role === 'superadmin') return '/admin/settings'
                    if (role === 'secretary') return '/secretary/settings'
                    return '/'
                },
                meta: { title: 'Profilo Utente', roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent'] }
            },
            { path: 'support', component: () => import('@/pages/Support.vue'), meta: { title: 'Supporto & Assistenza', roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent'] } }
        ]
    },
    {
        path: '/login',
        component: () => import('@/layouts/LoginLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Login.vue'), meta: { title: 'Accesso al Sistema' } }
        ]
    },
    {
        path: '/register',
        component: () => import('@/layouts/LoginLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Login.vue'), meta: { title: 'Registrazione' } }
        ]
    },
    {
        path: '/forgot-password',
        component: () => import('@/layouts/LoginLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Login.vue'), meta: { title: 'Recupero Password' } }
        ]
    },
    {
        path: '/:catchAll(.*)*',
        component: () => import('@/pages/NotFound.vue'),
        meta: { title: 'Pagina non trovata' }
    }
]
