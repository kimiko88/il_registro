// Quasar Configuration
// This file is typically used by the CLI, but for Vite plugin usage
// we control most things in main.js and vite.config.js.
// However, we can keep global config object here if we want to import it.

export default {
    config: {
        brand: {
            primary: '#4F46E5',  // Indigo 600
            secondary: '#06B6D4', // Cyan 500
            accent: '#F59E0B',   // Amber 500

            dark: '#1E293B',     // Slate 900
            'dark-page': '#0F172A', // Slate 950

            positive: '#10B981', // Emerald 500
            negative: '#EF4444', // Red 500
            info: '#3B82F6',     // Blue 500
            warning: '#F59E0B'   // Amber 500
        },
        dark: 'auto' // or boolean true/false
    },
    plugins: [
        'Notify',
        'Dialog',
        'Loading'
    ]
}
