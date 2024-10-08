import { URL, fileURLToPath } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';

export default defineConfig({
	plugins: [vue(), vueDevTools()],
	resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url)),
		},
	},
	build: {
		outDir: './dist',
		emptyOutDir: true,
	},
	server: {
		host: true,
		proxy: {
			'tonconnect-manifest.json': {
				target: 'http://localhost:8081',
			},
			'/api': {
				target: 'http://localhost:8081',
				changeOrigin: true,
			},
			'/ws': {
				target: 'ws://localhost:8081',
				ws: true,
			},
		},
	},
});
