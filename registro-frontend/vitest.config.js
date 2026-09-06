/// <reference types="vitest" />
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { quasar, transformAssetUrls } from '@quasar/vite-plugin'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue({
            template: { transformAssetUrls }
        }),
        quasar({
            sassVariables: 'src/quasar-variables.sass'
        })
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
            'src': fileURLToPath(new URL('./src', import.meta.url))
        }
    },
    test: {
        globals: true,
        environment: 'happy-dom',
        setupFiles: ['./tests/unit/setup.js'],
        include: ['tests/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
        coverage: {
            provider: 'v8',
            reporter: ['text', 'html', 'lcov', 'json'],
            exclude: [
                'node_modules/**',
                'dist/**',
                'src/main.js',
                'src/boot/**',
                'src/i18n/**',
                'src/router/**',
                'src-pwa/**',
                '**/*.spec.js',
                '**/*.test.js',
                'tests/**'
            ]
        },
        deps: {
            inline: [/@quasar/]
        }
    }
})
