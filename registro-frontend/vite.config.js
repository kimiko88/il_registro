import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { quasar, transformAssetUrls } from '@quasar/vite-plugin'
import { VitePWA } from 'vite-plugin-pwa'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    // Use esbuild instead of terser to avoid serialize-javascript
    // crypto.randomUUID() error in CI / Node < 20 environments
    minify: 'esbuild',
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules')) {
            if (id.includes('quasar') || id.includes('@quasar')) {
              return 'vendor-quasar'
            }
            if (id.includes('chart.js') || id.includes('vue-chartjs')) {
              return 'vendor-charts'
            }
            if (id.includes('vue-i18n')) {
              return 'vendor-i18n'
            }
            if (id.includes('vue') || id.includes('pinia') || id.includes('vue-router') || id.includes('axios')) {
              return 'vendor-vue'
            }
          }
          if (id.includes('/src/i18n/it-IT/') || id.includes('\\src\\i18n\\it-IT\\') ||
              id.endsWith('/src/i18n/index.js') || id.endsWith('\\src\\i18n\\index.js')) {
            return 'app-i18n'
          }
        }
      }
    }
  },
  plugins: [
    vue({
      template: {
        transformAssetUrls,
        compilerOptions: {
          hoistStatic: false
        }
      }
    }),

    quasar({
      sassVariables: 'src/assets/styles/variables.scss'
    }),

    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'vite.svg'],
      manifest: {
        name: 'Registro Elettronico Scolastico',
        short_name: 'Registro',
        description: 'Registro Elettronico Scolastico Moderno',
        theme_color: '#4F46E5',
        background_color: '#F8FAFC',
        display: 'standalone',
        orientation: 'portrait',
        scope: '/',
        start_url: '/',
        icons: [
          {
            src: 'pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png'
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png'
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any maskable'
          }
        ]
      },
      workbox: {
        navigateFallback: '/index.html',
        navigateFallbackDenylist: [/^\/api/],
        globPatterns: ['**/*.{js,css,html,ico,png,svg,woff,woff2}'],
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/fonts\.(?:googleapis|gstatic)\.com\/.*/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'google-fonts',
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 60 * 60 * 24 * 365 // 1 year
              }
            }
          },
          {
            urlPattern: /^https:\/\/cdn\.quasar\.dev\/.*/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'quasar-assets',
              expiration: {
                maxEntries: 20,
                maxAgeSeconds: 60 * 60 * 24 * 30 // 30 days
              }
            }
          },
          {
            // API responses caching: Stale-while-revalidate for static list resources and class rosters
            urlPattern: ({ url }) => {
              return url.pathname.includes('/teachers') ||
                     url.pathname.includes('/subjects') ||
                     url.pathname.includes('/schools') ||
                     url.pathname.includes('/classes') ||
                     url.pathname.includes('/students') ||
                     url.pathname.includes('/academic-years') ||
                     url.pathname.includes('/school-year')
            },
            handler: 'StaleWhileRevalidate',
            options: {
              cacheName: 'api-static-lists',
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 60 * 60 * 24 // 24 hours
              },
              cacheableResponse: {
                statuses: [0, 200]
              }
            }
          }
        ]
      }
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      'src': fileURLToPath(new URL('./src', import.meta.url))
    },
    dedupe: ['vue']
  },
  css: {
    preprocessorOptions: {
      scss: {
        api: 'modern-compiler',
        silenceDeprecations: ['legacy-js-api'],
      }
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        ws: true
      }
    }
  }
})
