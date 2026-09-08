import js from '@eslint/js'
import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'

export default [
  {
    name: 'app/files-to-lint',
    files: ['**/*.{js,mjs,cjs,vue}'],
  },
  {
    name: 'app/files-to-ignore',
    ignores: ['**/dist/**', '**/node_modules/**', '**/coverage/**', '**/.git/**', '**/public/**'],
  },
  js.configs.recommended,
  ...pluginVue.configs['flat/essential'],
  {
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.node,
        ...globals.es2021,
      },
    },
    rules: {
      'vue/multi-word-component-names': 'off',
      'no-unused-vars': ['warn', {
        argsIgnorePattern: '^(_|e$|err$|error$|row$|\\$event$|item$)',
        varsIgnorePattern: '^(_|__VLS_|\\$q$|loading$|t$|te$|themeStore$|qdate$|api$)',
        destructuredArrayIgnorePattern: '^_',
        caughtErrors: 'none',
      }],
      'no-undef': 'off',
    },
  },
]
