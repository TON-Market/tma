import { defineStore } from 'pinia';

import { $API } from '@/api';

export enum DepositStatus {
	OK,
	ERROR,
}

interface PayData {
	eventID: string;
	collateral: number;
	token: string;
}

interface PayResp {
	message: {
		address: string;
		amount: string;
		payload: string;
	};
	depositID: string;
}

interface DepositApprovement {
	depositStatus: DepositStatus;
	depositID: PayResp['depositID'];
}

export const useAccount = defineStore('account', {
	actions: {
		async getPaymentInfo(data: PayData) {
			return (await $API.post<PayResp>('pay', data)).data;
		},
		approveDeposit(data: DepositApprovement) {
			$API.post('deposit', data);
		},
		async getProof() {
			return (await $API.post('generate-payload')).data.payload;
		},
		async checkProof(data) {
			return await $API.post('check-proof', data);
		},
	},
});
