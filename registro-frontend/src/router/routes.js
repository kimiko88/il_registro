export default [
    {
        path: '/',
        component: () => import('@/layouts/MainLayout.vue'),
        children: [
            { path: '', component: () => import('@/pages/Dashboard.vue') },

            // Admin Routes
            { path: 'admin', component: () => import('@/pages/admin/Index.vue'), meta: { role: 'admin' } },
            { path: 'admin/users', component: () => import('@/pages/admin/Users.vue'), meta: { role: 'admin' } },
            { path: 'admin/schools', component: () => import('@/pages/admin/Schools.vue'), meta: { role: 'admin' } },

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
