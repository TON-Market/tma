import { defineStore } from 'pinia';

import { $API } from '@/api';

export const useAccount = defineStore('account', {
	actions: {
		async getPaymentInfo() {
			return (await $API.post('pay')).data;
		},
		async getProof() {
			return (await $API.post('generate-payload')).data.payload;
		},
		async checkProof(data) {
			return await $API.post('check-proof', data);
		},
	},
});
