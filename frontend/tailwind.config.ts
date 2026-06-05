import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  content: [
    './app/components/**/*.{vue,js,ts}',
    './app/layouts/**/*.vue',
    './app/pages/**/*.vue',
    './app/plugins/**/*.{js,ts}',
    './app/composables/**/*.{js,ts}',
    './app/app.vue',
    './app/error.vue',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
      colors: {
        brand: {
          50:  '#FFF1EA',
          100: '#FFE0CF',
          200: '#FFC09F',
          300: '#FFA070',
          400: '#FF8040',
          500: '#FF661A',
          600: '#E85710',
          700: '#B8430C',
          800: '#883109',
          900: '#5C2107',
        },
        ink: {
          900: '#1F2937',
          800: '#2C3744',
          700: '#3A4654',
        },
      },
      borderRadius: {
        xl: '0.875rem',
      },
      boxShadow: {
        soft: '0 4px 14px rgba(0,0,0,0.06)',
        card: '0 10px 30px -10px rgba(31,41,55,0.18)',
      },
      maxWidth: {
        '8xl': '88rem',
      },
    },
  },
  plugins: [],
}
