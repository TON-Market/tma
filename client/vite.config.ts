import { URL, fileURLToPath } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';
import vuetify from 'vite-plugin-vuetify';

const SERVER_HOST = 'localhost:8081';

export default defineConfig({
	plugins: [vue(), vueDevTools(), vuetify()],
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
			'/api': {
				target: `http://${SERVER_HOST}`,
				changeOrigin: true,
			},
			'/ws': {
				target: `ws://${SERVER_HOST}`,
				ws: true,
			},
		},
	},
});
