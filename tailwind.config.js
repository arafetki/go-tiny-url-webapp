/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./assets/templates/**/*.{html,templ,tmpl}"
  ],
  theme: {
    extend: {
      fontFamily: {
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
};