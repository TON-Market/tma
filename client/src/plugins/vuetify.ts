import '@mdi/font/css/materialdesignicons.css';
import 'vuetify/styles';
import { type ThemeDefinition, createVuetify } from 'vuetify';

const themeSettingsLight: ThemeDefinition = {
	colors: {},
};

const themeSettingsDark: ThemeDefinition = {
	colors: {
		surface: '#0d0b28',
		'surface-soft': '#2a2653',
		'surface-mute': '#433b75',

		primary: '#562aff',
		purple: '#6A4DFF',

		success: '#00D9A2',
		error: '#FF4B61',
	},
};

// const theme =
// 	localStorage.getItem('theme') ||
// 	(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
const theme = 'dark';

export default createVuetify({
	theme: {
		themes: {
			light: themeSettingsLight,
			dark: {
				...themeSettingsLight,
				colors: {
					...themeSettingsLight.colors,
					...themeSettingsDark.colors,
				},
			},
		},
		defaultTheme: theme,
	},
	defaults: {
		VDialog: {
			rounded: 'lg',
		},
		VBtn: {
			rounded: 'lg',
			color: 'primary',
			density: 'comfortable',
			variant: 'outlined',
		},
		VIcon: {
			color: 'secondary',
		},
		VAutocomplete: {
			variant: 'outlined',
			density: 'compact',
			itemValue: 'id',
			itemTitle: 'name',
		},
		VTextField: {
			rounded: 'lg',
			variant: 'outlined',
			density: 'compact',
			hideDetails: true,
			bgColor: '#a290ff33',
		},
		VNumberInput: {
			rounded: 'lg',
			variant: 'outlined',
			density: 'compact',
			controlVariant: 'stacked',
			hideDetails: true,
			VBtn: {
				border: 0,
				rounded: 0,
			},
		},
		VTabs: {
			bgColor: 'primary',
			density: 'compact',
		},
		VToolbar: {
			color: 'surface',
			density: 'compact',
		},
	},
});
