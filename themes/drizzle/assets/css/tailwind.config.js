// tailwind.config.js
const colors = require("tailwindcss/colors");
module.exports = {
  theme: {
    extend: {
      colors: {
        bg: "#1a1827",
      },
    },
  },
  plugins: [
    require("@tailwindcss/typography"),
    // ...
  ],
};
