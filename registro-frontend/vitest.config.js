import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { quasar, transformAssetUrls } from '@quasar/vite-plugin'
import path from 'path'

export default defineConfig({
    test: {
        environment: 'happy-dom',
        globals: true,
    },
    plugins: [
        vue({
            template: { transformAssetUrls }
        }),
        quasar({
            sassVariables: 'src/assets/styles/variables.scss'
        })
    ],
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src')
        }
    }
})
