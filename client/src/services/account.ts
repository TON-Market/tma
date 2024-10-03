import { defineStore } from 'pinia';

import { $API } from '@/api';

export const useWallet = defineStore('wallet', {
	actions: {
		bet() {
			return $API.get('get-address');
		},
		async getProof() {
			return (await $API.post('generate-payload')).data.payload;
		},
		async checkProof(data) {
			return await $API.post('check-proof', data);
		},
	},
});
