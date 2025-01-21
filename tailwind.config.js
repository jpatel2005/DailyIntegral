import defaultTheme from 'tailwindcss/defaultTheme'

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}', './node_modules/flowbite/**/*.js'],
  theme: {
    extend: {
      fontFamily: {
        serif: ['Bebas Neue', 'Playwrite AU SA', ...defaultTheme.fontFamily.serif],
        display: ['"Jersey 10"'],
      },
    },
  },
}
