import { URL, fileURLToPath } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [vue(), vueDevTools()],
	resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url)),
		},
	},
	server: {
		proxy: {
			'/ton-market': {
				target: 'https://oip492-88-201-232-88.ru.tuna.am/',
				changeOrigin: true,
			},
		},
	},
});
