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

            // Student Routes
            { path: 'student', component: () => import('@/pages/student/Index.vue'), meta: { role: 'student' } },

            // Parent Routes
            { path: 'parent', component: () => import('@/pages/parent/Index.vue'), meta: { role: 'parent' } },
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
