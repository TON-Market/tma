import { defineStore } from 'pinia';

import { $API } from '@/api';

export interface Bet {
	collateral: number;
	title: string;
	percentage: string;
	token: string;
	logoLink: string;
}

export interface Event {
	id: string;
	collateral: number;
	logoLink: string;
	title: string;
	bets: Bet[];
}

export interface Tag {
	id: number;
	title: string;
}

interface EventStore {
	loading: boolean;
	failed: boolean;

	tags: Tag[];
	tag: Tag['id'];

	pages: number;
	page: number;

	items: Event[];
	current: Event | null;
}

export const useEvents = defineStore('events', {
	state: (): EventStore => ({
		failed: false,
		loading: false,

		tags: [],
		tag: 5,

		pages: 1,
		page: 0,

		items: [],
		current: null,
	}),
	getters: {
		sorted: state => state.items.sort((a, b) => b.collateral - a.collateral),
	},
	actions: {
		async getTags() {
			this.tags = (await $API.get<Tag[]>('get-tags')).data;
		},
		onLoaded() {},
		async load(page: number, tag = 5, clear = false) {
			try {
				this.loading = true;
				this.page = page;
				this.tag = tag;
				if (clear) {
					this.items = [];
				}
				const data = (
					await $API.get<{ items: Event[]; pages: number }>('get-events', {
						params: {
							page,
							tag,
						},
					})
				).data;

				data.items.forEach(it => {
					this.update(it);
				});
				this.pages = data.pages;
			} catch (e) {
				this.failed = true;
			} finally {
				this.loading = false;
				this.onLoaded();
			}
		},
		update(data: Event) {
			const current = this.items.find(c => c.id === data.id);
			if (current) {
				Object.assign(current, data);
			} else {
				this.items.push(data);
			}
		},
		select(id: Event['id']) {
			this.current = this.items.find(it => it.id === id) || null;
		},
		async nextPage() {
			if (this.loading) return;
			if (this.pages > this.page) {
				this.page++;
				await this.load(this.page, this.tag);
			}
		},
	},
});
