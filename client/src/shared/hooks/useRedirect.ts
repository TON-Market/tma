import type { RouteParamsRaw, RouteRecordSingleViewWithChildren } from 'vue-router';

import router, { routes } from '@/router';

type Routes = typeof routes;

type Paths<T> =
	T extends Array<infer Item>
		? Item extends { name: string }
			? Item extends RouteRecordSingleViewWithChildren
				? `${Item['name']}`
				: `${Item['name']}`
			: never
		: never;

export type View = Paths<Routes>;

export const useRedirect = (name: View, params?: RouteParamsRaw) => {
	router.push({ name, params });
};
