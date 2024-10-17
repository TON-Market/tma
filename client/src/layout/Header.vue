<script setup lang="ts">
import { TonConnectButton, useTonAddress, useTonConnectUI } from '@townsquarelabs/ui-vue';

import Logo from '@/assets/icons/Logo.vue';
import Balance from '@/components/Balance.vue';
import { useAccount } from '@/services/account';
import { useNotifier } from '@/services/notifier';
import { useRedirect } from '@/shared/hooks/useRedirect';

const [tonConnectUI] = useTonConnectUI();

const notifier = useNotifier();

const account = useAccount();
const address = useTonAddress();

tonConnectUI.onModalStateChange(async e => {
	if (e.status === 'opened') {
		tonConnectUI.setConnectRequestParameters({
			state: 'loading',
		});

		tonConnectUI.setConnectRequestParameters({
			state: 'ready',
			value: { tonProof: await account.getProof() },
		});
	}
});

tonConnectUI.onStatusChange(async wallet => {
	setTimeout(() => {
		account.setAddress(address.value);
		account.getBalance();
	}, 100);
	if (!wallet) {
		account.disconnect();
		useRedirect('home');
	}
	if (wallet?.connectItems?.tonProof && 'proof' in wallet.connectItems.tonProof) {
		try {
			await account.checkProof({
				proof: wallet.connectItems.tonProof.proof,
				address: wallet.account.address,
				network: wallet.account.chain,
			});
		} catch {
			notifier.error('Wallet validation failed...', 15000);
			notifier.error('Please, try to make a transaction and connect wallet again.', 15000);
			tonConnectUI.disconnect();
		}
	}
});

tonConnectUI.connectionRestored.then(() => {
	setTimeout(() => {
		account.setAddress(address.value);
		account.getBalance();
	}, 100);
});

setInterval(() => {
	account.getBalance();
}, 5000);
</script>

<template>
	<header>
		<div class="header">
			<Logo @click="useRedirect('home')" />
			<div class="header-controls">
				<Balance />
				<TonConnectButton />
			</div>
		</div>
	</header>
</template>

<style scoped>
header {
	padding: 20px;
	display: flex;
	flex-direction: column;
	gap: 20px;
}
.header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 10px;
}
.header-controls {
	display: flex;
	max-width: 60%;
	justify-content: flex-end;

	background: linear-gradient(45deg, #9975ff36, #9975ff33);
	display: flex;
	align-items: center;
	color: var(--vt-c-purple);
	border-radius: 100vh;

	padding: 0 14px;
}
</style>
