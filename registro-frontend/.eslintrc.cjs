module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-essential',
  ],
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
  plugins: [
    'vue',
  ],
  rules: {
    'vue/multi-word-component-names': 'off',
    'no-unused-vars': ['warn', {
      argsIgnorePattern: '^(_|e$|err$|error$|row$|\\$event$)',
      varsIgnorePattern: '^(_|__VLS_|\\$q$|loading$)',
      destructuredArrayIgnorePattern: '^_',
      caughtErrors: 'none',
    }],


    'no-undef': 'off', // Let the compiler handle it

  },
};
