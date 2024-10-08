import { defineStore } from 'pinia';

interface Navigation {
	searchElement: HTMLElement | null;
	searchFocused: boolean;
}

export const useNavigation = defineStore('navigation', {
	state: (): Navigation => ({
		searchElement: null,
		searchFocused: false,
	}),
	actions: {
		registerSearchElement(el: HTMLElement) {
			this.searchElement = el.querySelector('input');
			if (!this.searchElement) return;
			this.searchElement.onfocus = () => {
				this.searchFocused = true;
			};
			this.searchElement.onblur = () => {
				this.searchFocused = false;
			};
		},
		toSearch() {
			// this.searchElement?.scrollIntoView({ behavior: 'smooth' });
			this.searchElement?.focus();
		},
	},
});
