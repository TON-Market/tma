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

interface Asset {
	eventTitle: string;
	betTitle: string;
	collateralStaked: string;
	size: string;
}

interface Account {
	address: string | null;
	balance: number | null;
	assets: Asset[];
	total: string;
}

export const useAccount = defineStore('account', {
	state: (): Account => ({
		address: null,
		balance: 10.0123,
		assets: [],
		total: '',
	}),
	actions: {
		async getBalance() {
			return;
			if (!this.address) {
				this.balance = null;
				return;
			}
			const data = (await $API.get(`https://tonapi.io/v2/accounts/${this.address}`)).data;
			this.balance = data.balance / 1000000000;
		},
		async getAssets() {
			const data = (
				await $API.get<{
					assetList: Asset[];
					totalInMarket: string;
				}>('assets')
			).data;

			this.assets = data.assetList;
			this.total = data.totalInMarket;
		},
		async getPaymentInfo(data: PayData) {
			return (await $API.post<PayResp>('pay', data)).data;
		},
		approveDeposit(data: DepositApprovement) {
			$API.post('deposit', data);
		},
		async getProof() {
			return (await $API.post('generate-payload')).data.payload;
		},
		checkProof(data: any) {
			return $API.post('check-proof', data);
		},
		setAddress(value: string) {
			this.address = value;
		},
		disconnect() {
			return $API.delete('disconnect');
		},
	},
});
