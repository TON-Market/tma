import { type RouteRecordRaw, createRouter, createWebHistory } from 'vue-router';

export const routes = [
	{
		path: '/',
		name: 'home',
		component: () => import('@/views/HomeView/HomeView.vue'),
	},
	{
		path: '/events/:id',
		name: 'event',
		component: () => import('@/views/EventView/EventView.vue'),
	},
	{
		path: '/account',
		name: 'account',
		component: () => import('@/views/AccountView/AccountView.vue'),
	},
	{
		path: '/buy/:eventID/:token',
		name: 'bet',
		component: () => import('@/views/BetView/BetView.vue'),
		props: true,
	},
] as const satisfies RouteRecordRaw[];

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes,
});

export default router;
