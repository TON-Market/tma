<script setup lang="ts">
import { CHAIN, useTonConnectModal, useTonConnectUI } from '@townsquarelabs/ui-vue';
import { ref } from 'vue';
import { VNumberInput } from 'vuetify/labs/VNumberInput';

import { DepositStatus, useAccount } from '@/services/account';

const [TonConnectUI] = useTonConnectUI();
const modal = useTonConnectModal();
const account = useAccount();

const tons = ref(1);

const { eventID, token } = defineProps<{
	eventID: string;
	token: string;
}>();

async function pay() {
	console.log({
		collateral: tons.value,
		eventID,
		token,
	});
	if (!TonConnectUI.connected) {
		modal.open();
		return;
	}
	try {
		const info = await account.getPaymentInfo({
			collateral: tons.value,
			eventID,
			token,
		});

		const data = await TonConnectUI.sendTransaction({
			validUntil: Math.floor(Date.now() / 1000) + 60,
			network: CHAIN.MAINNET,
			messages: [info.message],
		});

		account.approveDeposit({
			depositID: info.depositID,
			depositStatus: DepositStatus.OK,
		});
		console.log(data);
	} catch (e) {
		console.log(e);

		// if (e instanceof TonConnectError) {
		// 	modal.open();
		// }
	}
}
</script>

<template>
	<v-dialog max-width="500">
		<template v-slot:activator="{ props: activatorProps }">
			<v-btn
				v-bind="activatorProps"
				text="Buy"
				variant="flat"
			/>
		</template>

		<template v-slot:default="{ isActive }">
			<v-card title="Confirm">
				<v-card-text>
					<v-number-input
						label="TONs"
						:min="0"
						v-model:model-value="tons"
					/>
				</v-card-text>

				<v-card-actions>
					<v-spacer></v-spacer>

					<v-btn
						:disabled="tons <= 0"
						text="Buy"
						flat
						@click="
							() => {
								pay();
								isActive.value = false;
							}
						"
					/>
					<v-btn
						text="Cancel"
						@click="isActive.value = false"
					/>
				</v-card-actions>
			</v-card>
		</template>
	</v-dialog>
</template>
