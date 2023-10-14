module.exports = {
    purge: {
        // Enable purging only in production for better dev experience
        enabled: process.env.NODE_ENV === 'production',
        content: [
            './**/*.html',
            './**/*.go',
            // Add any other file types you might use
            // where Tailwind classes are referenced
        ],
    },
    darkMode: false, // or 'media' or 'class'
    theme: {
        extend: {},
    },
    variants: {},
    plugins: [],
}