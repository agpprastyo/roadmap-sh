/** @type {import('tailwindcss').Config} */
import typography from "@tailwindcss/typography";

export default {
  content: [
    "./index.html",
    "./src/**/*.jsx",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    typography,
  ],
};
