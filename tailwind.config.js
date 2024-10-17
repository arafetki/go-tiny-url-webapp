/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./assets/templates/**/*.{html,templ,tmpl}"
  ],
  theme: {
    container: {
      center: true,
    },
    extend: {
      fontFamily: {
      },
    },
  },
  darkMode: "class",
  plugins: [require("tailwindcss-animate")],
};