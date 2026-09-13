import { useAuthStore } from '@/stores/auth'

export default [
    {
        path: '/',
        component: () => import('@/layouts/MainLayout.vue'),
        meta: { requiresAuth: true },
        children: [
            { path: '', component: () => import('@/pages/Dashboard.vue'), meta: { title: 'Dashboard', titleKey: 'routeTitles.dashboard', roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent', 'principal', 'vice_principal', 'coordinator', 'staff', 'system_auditor', 'dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico'] } },
            { path: 'dashboard', redirect: '/' },

            // Admin Routes (SuperAdmin + Admin)
            {
                path: 'admin',
                redirect: '/admin/dashboard'
            },
            {
                path: 'admin/dashboard',
                component: () => import('@/pages/admin/Dashboard.vue'),
                meta: { title: 'Dashboard Amministratore', titleKey: 'routeTitles.adminDashboard', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/schools',
                component: () => import('@/pages/admin/SchoolManagement.vue'),
                meta: { title: 'Gestione Scuole', titleKey: 'routeTitles.schoolManagement', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/schools/:id',
                component: () => import('@/pages/admin/SchoolDetail.vue'),
                meta: { title: 'Dettaglio Scuola', titleKey: 'routeTitles.schoolDetail', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/admins',
                component: () => import('@/pages/admin/AdminUsers.vue'),
                meta: { title: 'Gestione Amministratori', titleKey: 'routeTitles.adminUsers', roles: ['superadmin'] }
            },
            {
                path: 'admin/monitoring',
                component: () => import('@/pages/admin/Monitoring.vue'),
                meta: { title: 'Monitoraggio Sistema', titleKey: 'routeTitles.systemMonitoring', roles: ['superadmin'] }
            },
            {
                path: 'admin/users',
                component: () => import('@/pages/secretary/Users.vue'),
                meta: { title: 'Gestione Utenti System', titleKey: 'routeTitles.systemUsers', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/analytics',
                component: () => import('@/pages/admin/Analytics.vue'),
                meta: { title: 'Statistiche e Analytics', titleKey: 'routeTitles.analytics', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/elearning',
                component: () => import('@/pages/admin/ElearningIntegration.vue'),
                meta: { title: 'Google Classroom & Teams', titleKey: 'routeTitles.elearning', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/settings',
                component: () => import('@/pages/admin/Settings.vue'),
                meta: { title: 'Impostazioni di Sistema', titleKey: 'routeTitles.systemSettings', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/scheduler',
                component: () => import('@/pages/admin/Scheduler.vue'),
                meta: { title: 'Pianificazione Task', titleKey: 'routeTitles.scheduler', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/school-settings',
                component: () => import('@/pages/admin/SchoolSettings.vue'),
                meta: { title: 'Configurazione Scuola', titleKey: 'routeTitles.schoolSettings', roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/audit-logs',
                component: () => import('@/pages/admin/AuditLog.vue'),
                // admin included so the dashboard menu item works correctly
                meta: { title: 'Registro Eventi & Audit', titleKey: 'routeTitles.auditLogs', roles: ['superadmin', 'admin', 'system_auditor'] }
            },
            { path: 'admin/audit-log', redirect: '/admin/audit-logs' },
            {
                path: 'admin/tenants',
                component: () => import('@/pages/admin/Tenants.vue'),
                meta: { title: 'Gestione Multi-Tenant', titleKey: 'routeTitles.tenants', roles: ['superadmin'] }
            },
            {
                // Redirect /admin/classes → secretary/classes (same component, admin-accessible)
                path: 'admin/classes',
                redirect: '/secretary/classes',
                meta: { title: 'Gestione Classi', titleKey: 'routeTitles.classes', roles: ['superadmin', 'admin'] }
            },

            // Secretary Routes
            { path: 'secretary', component: () => import('@/pages/secretary/Index.vue'), meta: { title: 'Pannello Segreteria', titleKey: 'routeTitles.secretaryPanel', roles: ['secretary', 'principal', 'vice_principal'] } },
            { path: 'secretary/dashboard', redirect: '/secretary' },
            { path: 'secretary/documents', component: () => import('@/pages/secretary/Documents.vue'), meta: { title: 'Gestione Documenti', titleKey: 'routeTitles.documents', roles: ['secretary', 'principal', 'vice_principal', 'dsga', 'assistente_amministrativo', 'collaboratore_ds'] } },
            { path: 'secretary/users', component: () => import('@/pages/secretary/Users.vue'), meta: { title: 'Anagrafica Utenti', titleKey: 'routeTitles.users', roles: ['secretary', 'principal', 'vice_principal', 'dsga', 'assistente_amministrativo'] } },
            { path: 'secretary/students', component: () => import('@/pages/secretary/Students.vue'), meta: { title: 'Anagrafica Studenti', titleKey: 'routeTitles.students', roles: ['secretary', 'principal', 'vice_principal', 'assistente_amministrativo'] } },
            { path: 'secretary/communications', component: () => import('@/pages/secretary/Communications.vue'), meta: { title: 'Circolari & Comunicazioni', titleKey: 'routeTitles.communications', roles: ['secretary', 'principal', 'vice_principal', 'dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico'] } },
            { path: 'secretary/reports', component: () => import('@/pages/secretary/Reports.vue'), meta: { title: 'Reportistica Scolastica', titleKey: 'routeTitles.reports', roles: ['secretary', 'principal', 'vice_principal'] } },
            { path: 'secretary/pcto', component: () => import('@/pages/secretary/PCTO.vue'), meta: { title: 'Gestione PCTO', titleKey: 'routeTitles.pcto', roles: ['secretary'] } },
            { path: 'secretary/textbooks', component: () => import('@/pages/secretary/Textbooks.vue'), meta: { title: 'Adozione Libri di Testo', titleKey: 'routeTitles.textbooks', roles: ['secretary'] } },
            { path: 'secretary/settings', component: () => import('@/pages/secretary/Settings.vue'), meta: { title: 'Impostazioni Segreteria', titleKey: 'routeTitles.settings', roles: ['secretary', 'principal', 'vice_principal'] } },
            { path: 'secretary/classes', component: () => import('@/pages/secretary/Classes.vue'), meta: { title: 'Gestione Classi', titleKey: 'routeTitles.classes', roles: ['secretary', 'principal', 'vice_principal'] } },
            { path: 'secretary/timetable', component: () => import('@/pages/secretary/Timetable.vue'), meta: { title: 'Orario Scolastico & Cattedre', titleKey: 'routeTitles.timetable', roles: ['secretary', 'admin', 'superadmin', 'principal', 'vice_principal', 'collaboratore_ds', 'dsga'] } },
            { path: 'secretary/scrutiny', component: () => import('@/pages/secretary/Scrutiny.vue'), meta: { title: 'Scrutini Scolastici', titleKey: 'routeTitles.scrutiny', roles: ['secretary', 'principal', 'vice_principal'] } },

            { path: 'secretary/groups', component: () => import('@/pages/secretary/Groups.vue'), meta: { title: 'Gruppi Linguistici / Articolati', titleKey: 'routeTitles.classes', roles: ['secretary'] } },
            { path: 'secretary/meetings', component: () => import('@/pages/secretary/Meetings.vue'), meta: { title: 'Organizzazione Riunioni', titleKey: 'routeTitles.meetings', roles: ['secretary'] } },
            { path: 'secretary/certificates', component: () => import('@/pages/secretary/Certificates.vue'), meta: { title: 'Certificati & Attestati', titleKey: 'routeTitles.certificates', roles: ['secretary', 'admin', 'superadmin', 'principal', 'vice_principal', 'assistente_amministrativo', 'dsga'] } },
            { path: 'secretary/substitutions', component: () => import('@/pages/secretary/Substitutions.vue'), meta: { title: 'Gestione Sostituzioni Docenti', titleKey: 'routeTitles.substitutions', roles: ['secretary', 'admin', 'superadmin', 'principal', 'vice_principal', 'dsga', 'collaboratore_ds'] } },
            { path: 'secretary/sidi', component: () => import('@/pages/secretary/SidiExports.vue'), meta: { title: 'Flussi SIDI (MIM)', titleKey: 'routeTitles.sidiExports', roles: ['secretary', 'admin', 'superadmin', 'principal', 'vice_principal', 'dsga'] } },
            { path: 'secretary/verbali', component: () => import('@/pages/secretary/VerbaliManagement.vue'), meta: { title: 'Verbali & Modelli Riunioni', titleKey: 'routeTitles.verbali', roles: ['secretary', 'principal', 'vice_principal', 'admin', 'superadmin', 'dsga', 'collaboratore_ds', 'assistente_amministrativo'] } },
            { path: 'secretary/students/:id/fascicolo', component: () => import('@/pages/secretary/FascicoloStudente.vue'), meta: { title: 'Fascicolo Studente', titleKey: 'routeTitles.fascicolo', roles: ['secretary', 'principal', 'vice_principal'] } },
            { path: 'secretary/audit-logs', redirect: '/admin/audit-logs' },
            { path: 'secretary/audit-log', redirect: '/admin/audit-logs' },

            // ATA Routes & Staff Attendance (Presenze ATA e Docenti)
            { path: 'ata', component: () => import('@/pages/ata/Index.vue'), meta: { title: 'Pannello ATA', titleKey: 'routeTitles.ataPanel', roles: ['dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico', 'principal', 'vice_principal', 'secretary', 'admin', 'superadmin'] } },
            { path: 'ata/attendance', component: () => import('@/pages/ata/StaffAttendance.vue'), meta: { title: 'Presenze Personale & Docenti', titleKey: 'routeTitles.staffAttendance', roles: ['dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico', 'principal', 'vice_principal', 'secretary', 'admin', 'superadmin'] } },
            { path: 'ata/emergency-substitutions', component: () => import('@/pages/ata/EmergencySubstitutions.vue'), meta: { title: 'Emergenza Sostituzioni', titleKey: 'routeTitles.emergencySubstitutions', roles: ['collaboratore_ds', 'dsga', 'principal', 'vice_principal', 'secretary', 'admin', 'superadmin'] } },
            { path: 'ata/visitor-registry', component: () => import('@/pages/ata/VisitorRegistry.vue'), meta: { title: 'Registro Visitatori', titleKey: 'routeTitles.visitorRegistry', roles: ['dsga', 'assistente_amministrativo', 'collaboratore_scolastico', 'collaboratore_ds', 'principal', 'vice_principal', 'secretary', 'admin', 'superadmin'] } },
            { path: 'ata/timecard', component: () => import('@/pages/ata/Timecard.vue'), meta: { title: 'Cartellino & Piano Ferie', titleKey: 'routeTitles.timecard', roles: ['dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico', 'principal', 'admin', 'superadmin'] } },
            { path: 'ata/strike', component: () => import('@/pages/ata/StrikeManagement.vue'), meta: { title: 'Rilevazione Preventiva Scioperi', titleKey: 'routeTitles.strikeManagement', roles: ['dsga', 'principal', 'vice_principal', 'admin', 'superadmin', 'collaboratore_ds', 'assistente_amministrativo', 'secretary'] } },
            { path: 'ata/personnel-desk', component: () => import('@/pages/ata/PersonnelDesk.vue'), meta: { title: 'Sportello Digitale Personale', titleKey: 'routeTitles.personnelDesk', roles: ['dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico', 'teacher', 'coordinator', 'principal', 'vice_principal', 'secretary', 'admin', 'superadmin'] } },

            // Teacher Routes (Supports both 'teacher' and 'coordinator')
            { path: 'teacher', component: () => import('@/pages/teacher/Index.vue'), meta: { title: 'Pannello Docente', titleKey: 'routeTitles.teacherPanel', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/dashboard', redirect: '/teacher' },
            { path: 'teacher/grades', component: () => import('@/pages/teacher/Grades.vue'), meta: { title: 'Gestione Voti', titleKey: 'routeTitles.grades', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/attendance', component: () => import('@/pages/teacher/Attendance.vue'), meta: { title: 'Registro Appello & Presenze', titleKey: 'routeTitles.attendance', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/classes', component: () => import('@/pages/teacher/Classes.vue'), meta: { title: 'Le Mie Classi', titleKey: 'routeTitles.classes', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/groups', component: () => import('@/pages/teacher/Groups.vue'), meta: { title: 'Gruppi e Laboratori', titleKey: 'routeTitles.classes', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/coordinator', component: () => import('@/pages/teacher/CoordinatorView.vue'), meta: { title: 'Pannello Coordinatore', titleKey: 'routeTitles.coordinator', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/documents', component: () => import('@/pages/teacher/Documents.vue'), meta: { title: 'Documentazione Didattica', titleKey: 'routeTitles.documents', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/colloqui', component: () => import('@/pages/teacher/Colloqui.vue'), meta: { title: 'Ricevimento Famiglie', titleKey: 'routeTitles.colloqui', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/communications', component: () => import('@/pages/teacher/Communications.vue'), meta: { title: 'Comunicazioni Scuola', titleKey: 'routeTitles.communications', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/lessons', component: () => import('@/components/Teacher/LessonPlanner.vue'), meta: { title: 'Registro Lezioni', titleKey: 'routeTitles.timetable', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/didactics', component: () => import('@/pages/teacher/Didactics.vue'), meta: { title: 'Materiale Didattico', titleKey: 'routeTitles.didactics', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/timetable', component: () => import('@/pages/teacher/Timetable.vue'), meta: { title: 'Orario Docente', titleKey: 'routeTitles.timetable', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/scrutiny', component: () => import('@/pages/teacher/Scrutiny.vue'), meta: { title: 'Gestione Scrutini', titleKey: 'routeTitles.scrutiny', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/verbali', component: () => import('@/pages/teacher/Verbali.vue'), meta: { title: 'Verbali Consiglio di Classe', titleKey: 'routeTitles.verbali', roles: ['teacher', 'coordinator', 'principal', 'vice_principal'] } },
            { path: 'teacher/substitutions', component: () => import('@/pages/teacher/Substitutions.vue'), meta: { title: 'Sostituzioni Docenti', titleKey: 'routeTitles.substitutions', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/grade-weights', component: () => import('@/pages/teacher/GradeWeights.vue'), meta: { title: 'Pesi e Criteri Valutazione', titleKey: 'routeTitles.grades', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/agenda', component: () => import('@/pages/teacher/Agenda.vue'), meta: { title: 'Agenda di Classe', titleKey: 'routeTitles.agenda', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/notes', component: () => import('@/pages/teacher/Notes.vue'), meta: { title: 'Note & Richiami', titleKey: 'routeTitles.notes', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/rubrics', component: () => import('@/pages/teacher/Rubrics.vue'), meta: { title: 'Rubriche Valutative', titleKey: 'routeTitles.rubrics', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/pdp', component: () => import('@/pages/teacher/PdpPlans.vue'), meta: { title: 'Piani PDP / PEI', titleKey: 'routeTitles.pdp', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/uda', component: () => import('@/pages/teacher/UdaPlanner.vue'), meta: { title: 'Programmazione Didattica UdA', titleKey: 'routeTitles.uda', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/competencies', component: () => import('@/pages/teacher/Competencies.vue'), meta: { title: 'Valutazione per Competenze', titleKey: 'routeTitles.eportfolio', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/settings', component: () => import('@/pages/teacher/Settings.vue'), meta: { title: 'Impostazioni Docente', titleKey: 'routeTitles.settings', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/recovery', component: () => import('@/pages/teacher/RecoveryCourses.vue'), meta: { title: 'Corsi di Recupero & PAI', titleKey: 'routeTitles.recovery', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/credits', component: () => import('@/pages/teacher/SchoolCredits.vue'), meta: { title: 'Credito Scolastico Triennio', titleKey: 'routeTitles.credits', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/support', component: () => import('@/pages/teacher/SupportRegister.vue'), meta: { title: 'Registro di Sostegno & PEI', titleKey: 'routeTitles.supportRegister', roles: ['teacher', 'coordinator'] } },
            { path: 'teacher/general-meetings', component: () => import('@/pages/teacher/GeneralMeetingLiveQueue.vue'), meta: { title: 'Ricevimento Generale Pomeridiano', titleKey: 'routeTitles.liveQueue', roles: ['teacher', 'coordinator'] } },

            // Student Routes
            { path: 'student', component: () => import('@/pages/student/Index.vue'), meta: { title: 'Pannello Studente', titleKey: 'routeTitles.studentPanel', roles: ['student'] } },
            { path: 'student/dashboard', redirect: '/student' },
            { path: 'student/grades', component: () => import('@/pages/student/Grades.vue'), meta: { title: 'I Miei Voti', titleKey: 'routeTitles.grades', roles: ['student'] } },
            { path: 'student/attendance', component: () => import('@/pages/student/Attendance.vue'), meta: { title: 'Presenze & Assenze', titleKey: 'routeTitles.attendance', roles: ['student'] } },
            { path: 'student/documents', component: () => import('@/pages/student/Documents.vue'), meta: { title: 'Documenti Studente', titleKey: 'routeTitles.documents', roles: ['student'] } },
            { path: 'student/profile', component: () => import('@/pages/student/Profile.vue'), meta: { title: 'Profilo Studente', titleKey: 'routeTitles.settings', roles: ['student'] } },
            { path: 'student/pcto', component: () => import('@/pages/student/PCTO.vue'), meta: { title: 'Percorsi PCTO', titleKey: 'routeTitles.pcto', roles: ['student'] } },
            { path: 'student/orientamento', component: () => import('@/pages/student/Orientamento.vue'), meta: { title: 'Orientamento Scolastico', titleKey: 'routeTitles.eportfolio', roles: ['student'] } },
            { path: 'student/communications', component: () => import('@/pages/student/Communications.vue'), meta: { title: 'Comunicati Scolastici', titleKey: 'routeTitles.communications', roles: ['student'] } },
            { path: 'student/homework', component: () => import('@/pages/student/Homework.vue'), meta: { title: 'Compiti & Attività', titleKey: 'routeTitles.homework', roles: ['student'] } },
            { path: 'student/agenda', component: () => import('@/pages/student/AgendaCalendar.vue'), meta: { title: 'Calendario Agenda', titleKey: 'routeTitles.agenda', roles: ['student'] } },
            { path: 'student/didactics', component: () => import('@/pages/student/Didactics.vue'), meta: { title: 'Materiali Didattici', titleKey: 'routeTitles.didactics', roles: ['student'] } },
            { path: 'student/notes', component: () => import('@/pages/student/Notes.vue'), meta: { title: 'Note Disciplinari', titleKey: 'routeTitles.notes', roles: ['student'] } },
            { path: 'student/timetable', component: () => import('@/pages/student/Timetable.vue'), meta: { title: 'Orario delle Lezioni', titleKey: 'routeTitles.timetable', roles: ['student'] } },
            { path: 'student/goals', component: () => import('@/pages/student/Goals.vue'), meta: { title: 'Obiettivi di Apprendimento', titleKey: 'routeTitles.goals', roles: ['student'] } },
            { path: 'student/school-calendar', component: () => import('@/pages/student/SchoolCalendar.vue'), meta: { title: 'Calendario Scolastico', titleKey: 'routeTitles.agenda', roles: ['student'] } },
            { path: 'student/report-card', component: () => import('@/pages/student/ReportCard.vue'), meta: { title: 'Pagella Online', titleKey: 'routeTitles.reportCard', roles: ['student'] } },
            { path: 'student/settings', component: () => import('@/pages/student/Settings.vue'), meta: { title: 'Impostazioni Studente', titleKey: 'routeTitles.settings', roles: ['student'] } },

            // Parent Routes
            { path: 'parent', component: () => import('@/pages/parent/Index.vue'), meta: { title: 'Pannello Famiglie', titleKey: 'routeTitles.parentPanel', roles: ['parent'] } },
            { path: 'parent/dashboard', redirect: '/parent' },
            { path: 'parent/children', component: () => import('@/pages/parent/Children.vue'), meta: { title: 'Figli Associati', titleKey: 'routeTitles.students', roles: ['parent'] } },
            { path: 'parent/grades', component: () => import('@/pages/parent/Grades.vue'), meta: { title: 'Valutazioni Figli', titleKey: 'routeTitles.grades', roles: ['parent'] } },
            { path: 'parent/report-card', component: () => import('@/pages/parent/ReportCard.vue'), meta: { title: 'Pagella Scolastica', titleKey: 'routeTitles.reportCard', roles: ['parent'] } },
            { path: 'parent/attendance', component: () => import('@/pages/parent/Attendance.vue'), meta: { title: 'Presenze & Giustificazioni', titleKey: 'routeTitles.attendance', roles: ['parent'] } },
            { path: 'parent/colloqui', component: () => import('@/pages/parent/Colloqui.vue'), meta: { title: 'Prenotazione Colloqui', titleKey: 'routeTitles.colloqui', roles: ['parent'] } },
            { path: 'parent/general-meetings', component: () => import('@/pages/parent/GeneralMeetingBooking.vue'), meta: { title: 'Ricevimento Generale Scuola-Famiglia', titleKey: 'routeTitles.generalMeetingBooking', roles: ['parent'] } },
            { path: 'parent/communications', component: () => import('@/pages/parent/Communications.vue'), meta: { title: 'Avvisi & Circolari', titleKey: 'routeTitles.communications', roles: ['parent'] } },
            { path: 'parent/trips', component: () => import('@/pages/parent/Trips.vue'), meta: { title: 'Uscite & Viaggi di Istruzione', titleKey: 'routeTitles.agenda', roles: ['parent'] } },
            { path: 'parent/profile', component: () => import('@/pages/parent/Profile.vue'), meta: { title: 'Profilo Genitore', titleKey: 'routeTitles.settings', roles: ['parent'] } },
            { path: 'parent/didactics', component: () => import('@/pages/parent/Didactics.vue'), meta: { title: 'Didattica & Compiti', titleKey: 'routeTitles.didactics', roles: ['parent'] } },
            { path: 'parent/notes', component: () => import('@/pages/parent/Notes.vue'), meta: { title: 'Note & Sanzioni', titleKey: 'routeTitles.notes', roles: ['parent'] } },
            { path: 'parent/timetable', component: () => import('@/pages/parent/Timetable.vue'), meta: { title: 'Orario Lezioni', titleKey: 'routeTitles.timetable', roles: ['parent'] } },
            { path: 'parent/pdp', component: () => import('@/pages/parent/PdpView.vue'), meta: { title: 'Piano PDP / PEI', titleKey: 'routeTitles.pdp', roles: ['parent'] } },
            { path: 'parent/documents', component: () => import('@/pages/parent/Documents.vue'), meta: { title: 'Documentazione & Moduli', titleKey: 'routeTitles.documents', roles: ['parent'] } },
            { path: 'parent/payments', component: () => import('@/pages/parent/Payments.vue'), meta: { title: 'Pagamenti Scolastici', titleKey: 'routeTitles.certificates', roles: ['parent'] } },
            { path: 'parent/settings', component: () => import('@/pages/parent/Settings.vue'), meta: { title: 'Impostazioni Genitore', titleKey: 'routeTitles.settings', roles: ['parent'] } },

            {
                path: 'communications',
                redirect: () => {
                    try {
                        const authStore = useAuthStore()
                        const role = authStore.userRole || authStore.user?.role || ''
                        if (role === 'teacher' || role === 'coordinator') return '/teacher/communications'
                        if (role === 'student') return '/student/communications'
                        if (role === 'parent') return '/parent/communications'
                        if (role === 'secretary' || role === 'principal' || role === 'vice_principal' || ['dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico'].includes(role)) return '/secretary/communications'
                        if (role === 'admin' || role === 'superadmin' || role === 'system_auditor') return '/admin/dashboard'
                    } catch { /* store not ready */ }
                    return '/'
                },
                meta: { title: 'Comunicazioni', titleKey: 'routeTitles.communications', roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent', 'principal', 'vice_principal', 'coordinator', 'system_auditor', 'dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico'] }
            },
            {
                path: 'profile',
                redirect: () => {
                    try {
                        const authStore = useAuthStore()
                        const role = authStore.userRole || authStore.user?.role || ''
                        if (role === 'student') return '/student/profile'
                        if (role === 'parent') return '/parent/profile'
                        if (role === 'admin' || role === 'superadmin') return '/admin/settings'
                        if (role === 'secretary' || role === 'principal' || role === 'vice_principal') return '/secretary/settings'
                        if (role === 'teacher' || role === 'coordinator') return '/teacher/settings'
                        if (['dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico'].includes(role)) return '/ata'
                    } catch { /* store not ready */ }
                    return '/'
                },
                meta: { title: 'Profilo Utente', titleKey: 'routeTitles.settings', roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent', 'principal', 'vice_principal', 'coordinator', 'dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico'] }
            },
            { path: 'support', component: () => import('@/pages/Support.vue'), meta: { title: 'Supporto & Assistenza', titleKey: 'routeTitles.support', roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent'] } },
            { path: 'accessibility-statement', component: () => import('@/pages/AccessibilityStatement.vue'), meta: { title: 'Dichiarazione di Accessibilità (AgID)', titleKey: 'routeTitles.accessibility', requiresAuth: false } },
            { path: 'dichiarazione-accessibilita', redirect: '/accessibility-statement' }
        ]
    },
    {
        path: '/login',
        component: () => import('@/layouts/LoginLayout.vue'),
        meta: { requiresAuth: false, titleKey: 'routeTitles.login' },
        children: [
            { path: '', component: () => import('@/pages/Login.vue'), meta: { title: 'Accesso al Sistema', titleKey: 'routeTitles.login', requiresAuth: false } }
        ]
    },
    {
        path: '/register',
        component: () => import('@/layouts/LoginLayout.vue'),
        meta: { requiresAuth: false, titleKey: 'routeTitles.register' },
        children: [
            { path: '', component: () => import('@/pages/Login.vue'), meta: { title: 'Registrazione', titleKey: 'routeTitles.register', requiresAuth: false } }
        ]
    },
    {
        path: '/forgot-password',
        component: () => import('@/layouts/LoginLayout.vue'),
        meta: { requiresAuth: false, titleKey: 'routeTitles.forgotPassword' },
        children: [
            { path: '', component: () => import('@/pages/Login.vue'), meta: { title: 'Recupero Password', titleKey: 'routeTitles.forgotPassword', requiresAuth: false } }
        ]
    },
    {
        path: '/:catchAll(.*)*',
        component: () => import('@/pages/NotFound.vue'),
        meta: { title: 'Pagina non trovata', titleKey: 'routeTitles.notFound', requiresAuth: false }
    }
]
