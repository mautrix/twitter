import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import wasm from "vite-plugin-wasm";
import topLevelAwait from "vite-plugin-top-level-await";

export default defineConfig({
	plugins: [tailwindcss(), wasm(), topLevelAwait()],
});
