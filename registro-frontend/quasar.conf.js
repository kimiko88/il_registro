// Quasar Configuration
// This file is typically used by the CLI, but for Vite plugin usage
// we control most things in main.js and vite.config.js.
// However, we can keep global config object here if we want to import it.

export default {
    config: {
        brand: {
            primary: '#1976d2',
            secondary: '#26a69a',
            accent: '#9c27b0',

            dark: '#1d1d1d',

            positive: '#21ba45',
            negative: '#c10015',
            info: '#31ccec',
            warning: '#f2c037'
        },
        dark: 'auto' // or boolean true/false
    },
    plugins: [
        'Notify',
        'Dialog',
        'Loading'
    ]
}
