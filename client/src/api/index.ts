import axios from 'axios';

import { useEvents } from '@/services/events';

export const $API = axios.create({
	baseURL: import.meta.env.VITE_API_URL,
});

export const $WS = new WebSocket(import.meta.env.VITE_WS_URL);

$WS.onmessage = raw => {
	const events = useEvents();
	events.update(JSON.parse(raw.data));
};
