/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{svelte,html,js}"],
  theme: {
    extend: {
      colors: {
        // https://colorpalette.pro/?color=oklch%2840%25+0.04+157%29&paletteType=tas&paletteStyle=square&colorFormat=oklch&effects=0%2C0%2C0%2C0
        bg: "oklch(12% 0.04 157)",
        bgmuted: "oklch(16% 0.04 157)",
        bgactive: "oklch(30% 0.04 157)",
        accent: "oklch(38% 0.04 157)",
        text: "oklch(98% 0.04 157)",
        textmuted: "oklch(60% 0.04 157)",
      },
    },
  },
  plugins: [],
};
