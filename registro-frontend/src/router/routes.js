export default [
    {
        path: '/',
        component: () => import('@/layouts/MainLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Dashboard.vue') },

            // Admin Routes
            { path: 'admin', component: () => import('@/pages/admin/Index.vue'), meta: { role: 'admin' } },
            { path: 'admin/schools', component: () => import('@/pages/admin/Schools.vue'), meta: { role: 'admin' } },
            { path: 'admin/users', component: () => import('@/pages/admin/AdminUsers.vue'), meta: { role: 'admin' } },
            { path: 'admin/monitoring', component: () => import('@/pages/admin/Monitoring.vue'), meta: { role: 'admin' } },
            { path: 'admin/analytics', component: () => import('@/pages/admin/Analytics.vue'), meta: { role: 'admin' } },
            { path: 'admin/settings', component: () => import('@/pages/admin/Settings.vue'), meta: { role: 'admin' } },

            // Secretary Routes
            { path: 'secretary', component: () => import('@/pages/secretary/Index.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/documents', component: () => import('@/pages/secretary/Documents.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/users', component: () => import('@/pages/secretary/Users.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/students', component: () => import('@/pages/secretary/Students.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/communications', component: () => import('@/pages/secretary/Communications.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/reports', component: () => import('@/pages/secretary/Reports.vue'), meta: { role: 'secretary' } },
            { path: 'secretary/settings', component: () => import('@/pages/secretary/Settings.vue'), meta: { role: 'secretary' } },

            // Teacher Routes
            { path: 'teacher', component: () => import('@/pages/teacher/Index.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/grades', component: () => import('@/pages/teacher/Grades.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/attendance', component: () => import('@/pages/teacher/Attendance.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/classes', component: () => import('@/pages/teacher/Classes.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/documents', component: () => import('@/pages/teacher/Documents.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/colloqui', component: () => import('@/pages/teacher/Colloqui.vue'), meta: { role: 'teacher' } },
            { path: 'teacher/communications', component: () => import('@/pages/teacher/Communications.vue'), meta: { role: 'teacher' } },

            // Student Routes
            { path: 'student', component: () => import('@/pages/student/Index.vue'), meta: { role: 'student' } },
            { path: 'student/grades', component: () => import('@/pages/student/Grades.vue'), meta: { role: 'student' } },
            { path: 'student/attendance', component: () => import('@/pages/student/Attendance.vue'), meta: { role: 'student' } },
            { path: 'student/documents', component: () => import('@/pages/student/Documents.vue'), meta: { role: 'student' } },
            { path: 'student/profile', component: () => import('@/pages/student/Profile.vue'), meta: { role: 'student' } },
            { path: 'student/pcto', component: () => import('@/pages/student/PCTO.vue'), meta: { role: 'student' } },
            { path: 'student/orientamento', component: () => import('@/pages/student/Orientamento.vue'), meta: { role: 'student' } },
            { path: 'student/communications', component: () => import('@/pages/student/Communications.vue'), meta: { role: 'student' } },

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
