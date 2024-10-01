import axios from 'axios';

const $API = axios.create({
	baseURL: '/ton-market/',
});

export { $API };
