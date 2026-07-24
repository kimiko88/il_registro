export default [
    {
        path: '/',
        component: () => import('@/layouts/MainLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Dashboard.vue'), meta: { roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent'] } },

            // Admin Routes (SuperAdmin + Admin)
            {
                path: 'admin',
                redirect: '/admin/dashboard'
            },
            {
                path: 'admin/dashboard',
                component: () => import('@/pages/admin/Dashboard.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/schools',
                component: () => import('@/pages/admin/SchoolManagement.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/schools/:id',
                component: () => import('@/pages/admin/SchoolDetail.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/admins',
                component: () => import('@/pages/admin/AdminUsers.vue'),
                meta: { roles: ['superadmin'] }
            },
            {
                path: 'admin/monitoring',
                component: () => import('@/pages/admin/Monitoring.vue'),
                meta: { roles: ['superadmin'] }
            },
            {
                path: 'admin/users',
                component: () => import('@/pages/secretary/Users.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/analytics',
                component: () => import('@/pages/admin/Analytics.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/settings',
                component: () => import('@/pages/admin/Settings.vue'),
                meta: { roles: ['superadmin'] }
            },
            {
                path: 'admin/scheduler',
                component: () => import('@/pages/admin/Scheduler.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/school-settings',
                component: () => import('@/pages/admin/SchoolSettings.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/audit-logs',
                component: () => import('@/pages/admin/AuditLog.vue'),
                meta: { roles: ['superadmin'] }
            },
            {
                path: 'admin/tenants',
                component: () => import('@/pages/admin/Tenants.vue'),
                meta: { roles: ['superadmin'] }
            },

            // Secretary Routes
            { path: 'secretary', component: () => import('@/pages/secretary/Index.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/documents', component: () => import('@/pages/secretary/Documents.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/users', component: () => import('@/pages/secretary/Users.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/students', component: () => import('@/pages/secretary/Students.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/communications', component: () => import('@/pages/secretary/Communications.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/reports', component: () => import('@/pages/secretary/Reports.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/pcto', component: () => import('@/pages/secretary/PCTO.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/textbooks', component: () => import('@/pages/secretary/Textbooks.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/settings', component: () => import('@/pages/secretary/Settings.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/classes', component: () => import('@/pages/secretary/Classes.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/scrutiny', component: () => import('@/pages/secretary/Scrutiny.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/groups', component: () => import('@/pages/teacher/Groups.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/meetings', component: () => import('@/pages/secretary/Meetings.vue'), meta: { roles: ['secretary'] } },
            { path: 'secretary/certificates', component: () => import('@/pages/secretary/Certificates.vue'), meta: { roles: ['secretary', 'admin', 'superadmin'] } },
            { path: 'secretary/audit-log', component: () => import('@/pages/secretary/AuditLog.vue'), meta: { roles: ['secretary', 'admin', 'superadmin'] } },
            { path: 'secretary/students/:id/fascicolo', component: () => import('@/pages/secretary/FascicoloStudente.vue'), meta: { roles: ['secretary'] } },

            // Teacher Routes
            { path: 'teacher', component: () => import('@/pages/teacher/Index.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/grades', component: () => import('@/pages/teacher/Grades.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/attendance', component: () => import('@/pages/teacher/Attendance.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/classes', component: () => import('@/pages/teacher/Classes.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/groups', component: () => import('@/pages/teacher/Groups.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/coordinator', component: () => import('@/pages/teacher/CoordinatorView.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/documents', component: () => import('@/pages/teacher/Documents.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/colloqui', component: () => import('@/pages/teacher/Colloqui.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/communications', component: () => import('@/pages/teacher/Communications.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/lessons', component: () => import('@/components/Teacher/LessonPlanner.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/didactics', component: () => import('@/pages/teacher/Didactics.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/timetable', component: () => import('@/pages/teacher/Timetable.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/scrutiny', component: () => import('@/pages/teacher/Scrutiny.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/verbali', component: () => import('@/pages/teacher/Verbali.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/substitutions', component: () => import('@/pages/teacher/Substitutions.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/grade-weights', component: () => import('@/pages/teacher/GradeWeights.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/agenda', component: () => import('@/pages/teacher/Agenda.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/notes', component: () => import('@/pages/teacher/Notes.vue'), meta: { roles: ['teacher'] } },
            { path: 'teacher/rubrics', component: () => import('@/pages/teacher/Rubrics.vue'), meta: { roles: ['teacher'] } },

            // Student Routes
            { path: 'student', component: () => import('@/pages/student/Index.vue'), meta: { roles: ['student'] } },
            { path: 'student/grades', component: () => import('@/pages/student/Grades.vue'), meta: { roles: ['student'] } },
            { path: 'student/attendance', component: () => import('@/pages/student/Attendance.vue'), meta: { roles: ['student'] } },
            { path: 'student/documents', component: () => import('@/pages/student/Documents.vue'), meta: { roles: ['student'] } },
            { path: 'student/profile', component: () => import('@/pages/student/Profile.vue'), meta: { roles: ['student'] } },
            { path: 'student/pcto', component: () => import('@/pages/student/PCTO.vue'), meta: { roles: ['student'] } },
            { path: 'student/orientamento', component: () => import('@/pages/student/Orientamento.vue'), meta: { roles: ['student'] } },
            { path: 'student/communications', component: () => import('@/pages/student/Communications.vue'), meta: { roles: ['student'] } },
            { path: 'student/homework', component: () => import('@/pages/student/Homework.vue'), meta: { roles: ['student'] } },
            { path: 'student/agenda', component: () => import('@/pages/student/AgendaCalendar.vue'), meta: { roles: ['student'] } },
            { path: 'student/didactics', component: () => import('@/pages/student/Didactics.vue'), meta: { roles: ['student'] } },
            { path: 'student/notes', component: () => import('@/pages/student/Notes.vue'), meta: { roles: ['student'] } },
            { path: 'student/timetable', component: () => import('@/pages/student/Timetable.vue'), meta: { roles: ['student'] } },
            { path: 'student/goals', component: () => import('@/pages/student/Goals.vue'), meta: { roles: ['student'] } },
            { path: 'student/school-calendar', component: () => import('@/pages/student/SchoolCalendar.vue'), meta: { roles: ['student'] } },
            { path: 'student/report-card', component: () => import('@/pages/student/ReportCard.vue'), meta: { roles: ['student'] } },

            // Parent Routes
            { path: 'parent', component: () => import('@/pages/parent/Index.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/children', component: () => import('@/pages/parent/Children.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/grades', component: () => import('@/pages/parent/Grades.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/report-card', component: () => import('@/pages/parent/ReportCard.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/attendance', component: () => import('@/pages/parent/Attendance.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/colloqui', component: () => import('@/pages/parent/Colloqui.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/communications', component: () => import('@/pages/parent/Communications.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/trips', component: () => import('@/pages/parent/Trips.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/profile', component: () => import('@/pages/parent/Profile.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/didactics', component: () => import('@/pages/parent/Didactics.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/notes', component: () => import('@/pages/parent/Notes.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/timetable', component: () => import('@/pages/parent/Timetable.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/documents', component: () => import('@/pages/parent/Documents.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/payments', component: () => import('@/pages/parent/Payments.vue'), meta: { roles: ['parent'] } },
            { path: 'parent/meetings', component: () => import('@/pages/parent/Meetings.vue'), meta: { roles: ['parent'] } },
            { path: 'support', component: () => import('@/pages/Support.vue'), meta: { roles: ['superadmin', 'admin', 'secretary', 'teacher', 'student', 'parent'] } }
        ]
    },
    {
        path: '/login',
        component: () => import('@/layouts/LoginLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Login.vue') }
        ]
    },
    {
        path: '/:catchAll(.*)*',
        component: () => import('@/pages/NotFound.vue')
    }
]
