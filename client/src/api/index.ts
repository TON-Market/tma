import axios from 'axios';

import { useEvents } from '@/services/events';

export const $API = axios.create({
	baseURL: '/ton-market/',
});

export const $WS = new WebSocket(`ws://bddxbv-88-201-232-88.ru.tuna.am/ton-market/ws`);

$WS.onmessage = raw => {
	const events = useEvents();
	events.update(JSON.parse(raw.data));
};
