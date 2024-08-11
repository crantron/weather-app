/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./templates/*.tpl'],
  theme: {
    extend: {},
  },
  plugins: [
      require('tailwindcss'),
      require('autoprefixer')
  ],
}

