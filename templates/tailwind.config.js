/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["../**/*.{html,js,go,css}"],
    theme: {
        extend: {},
    },
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
        require('@tailwindcss/aspect-ratio'),
    ],
}