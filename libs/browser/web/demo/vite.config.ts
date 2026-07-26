import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import path from "node:path";

export default defineConfig({
  root: path.resolve(__dirname),
  plugins: [react()],
  resolve: {
    alias: {
      "@fasionchan/browser-remote-react": path.resolve(__dirname, "../src/index.ts"),
    },
  },
  server: {
    port: 5173,
  },
});
