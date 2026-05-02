export default [
    {
        path: '/',
        component: () => import('@/layouts/MainLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Dashboard.vue') },

            // Admin Routes (SuperAdmin + Admin)
            {
                path: 'admin',
                redirect: '/admin/dashboard'
            },
            {
                path: 'admin/dashboard',
                component: () => import('@/pages/Admin/Dashboard.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/schools',
                component: () => import('@/pages/Admin/SchoolManagement.vue'),
                meta: { roles: ['superadmin', 'admin'] }
            },
            {
                path: 'admin/schools/:id',
                component: () => import('@/pages/Admin/SchoolDetail.vue'),
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
                path: 'admin/audit-logs',
                component: () => import('@/pages/Admin/AuditLog.vue'),
                meta: { roles: ['superadmin'] }
            },

            // Secretary Routes
            { path: 'secretary', component: () => import('@/pages/secretary/Index.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/documents', component: () => import('@/pages/secretary/Documents.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/users', component: () => import('@/pages/secretary/Users.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/students', component: () => import('@/pages/secretary/Students.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/communications', component: () => import('@/pages/secretary/Communications.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/reports', component: () => import('@/pages/secretary/Reports.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/pcto', component: () => import('@/pages/secretary/PCTO.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/textbooks', component: () => import('@/pages/secretary/Textbooks.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/scrutiny', component: () => import('@/pages/teacher/Scrutiny.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/settings', component: () => import('@/pages/secretary/Settings.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/classes', component: () => import('@/pages/secretary/Classes.vue'), meta: { role: 'secretary' } },

            // Teacher Routes
            { path: 'teacher', component: () => import('@/pages/teacher/Index.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/grades', component: () => import('@/pages/teacher/Grades.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/attendance', component: () => import('@/pages/teacher/Attendance.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/classes', component: () => import('@/pages/teacher/Classes.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/documents', component: () => import('@/pages/teacher/Documents.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/colloqui', component: () => import('@/pages/teacher/Colloqui.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/communications', component: () => import('@/pages/teacher/Communications.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/lessons', component: () => import('@/components/Teacher/LessonPlanner.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/scrutiny', component: () => import('@/pages/teacher/Scrutiny.vue'), meta: { role: 'teacher' } },

            // Student Routes
            { path: 'student', component: () => import('@/pages/student/Index.vue'), meta: { role: 'student' } },
            { path: 'student/grades', component: () => import('@/pages/student/Grades.vue'), meta: { role: 'student' } },
            { path: 'student/attendance', component: () => import('@/pages/student/Attendance.vue'), meta: { role: 'student' } },
            { path: 'student/documents', component: () => import('@/pages/student/Documents.vue'), meta: { role: 'student' } },
            { path: 'student/profile', component: () => import('@/pages/student/Profile.vue'), meta: { role: 'student' } },
            { path: 'student/pcto', component: () => import('@/pages/student/PCTO.vue'), meta: { role: 'student' } },
            { path: 'student/orientamento', component: () => import('@/pages/student/Orientamento.vue'), meta: { role: 'student' } },
            { path: 'student/communications', component: () => import('@/pages/student/Communications.vue'), meta: { role: 'student' } },
            { path: 'student/homework', component: () => import('@/pages/student/Homework.vue'), meta: { role: 'student' } },

            // Parent Routes
            { path: 'parent', component: () => import('@/pages/parent/Index.vue'), meta: { role: 'parent' } },
            { path: 'parent/children', component: () => import('@/pages/parent/Children.vue'), meta: { role: 'parent' } },
            { path: 'parent/grades', component: () => import('@/pages/parent/Grades.vue'), meta: { role: 'parent' } },
            { path: 'parent/attendance', component: () => import('@/pages/parent/Attendance.vue'), meta: { role: 'parent' } },
            { path: 'parent/colloqui', component: () => import('@/pages/parent/Colloqui.vue'), meta: { role: 'parent' } },
            { path: 'parent/communications', component: () => import('@/pages/parent/Communications.vue'), meta: { role: 'parent' } },
            { path: 'parent/profile', component: () => import('@/pages/parent/Profile.vue'), meta: { role: 'parent' } },
            { path: 'parent/support', component: () => import('@/pages/parent/Support.vue'), meta: { role: 'parent' } },
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
