/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}", // This is the important line
  ],
  theme: {
    extend: {
      // Add our custom font families
      fontFamily: {
        // Set 'Inter' as the default sans-serif font
        'primary': ['Lora', 'serif'],
        // Create a 'display' font family for 'Oswald'
        'secondary': ['Inter', 'sans-serif'],
      },
    },
  },
  plugins: [],
}