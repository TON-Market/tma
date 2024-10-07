import axios from 'axios';

import { useEvents } from '@/services/events';

console.log(import.meta.env);

export const $API = axios.create({
	baseURL: import.meta.env.VITE_HTTP_SERVER,
});

export const $WS = new WebSocket(import.meta.env.VITE_WS_SERVER);

$WS.onmessage = raw => {
	const events = useEvents();
	events.update(JSON.parse(raw.data));
};
