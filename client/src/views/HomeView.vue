<script setup lang="ts">
import { type Locales, useTonConnectUI } from '@townsquarelabs/ui-vue';
import { TonConnectButton } from '@townsquarelabs/ui-vue';

import { useWallet } from '@/services/account';

const [tonConnectUI, setOptions] = useTonConnectUI();

const api = useWallet();

setOptions({ language: 'ru' as Locales });

async function bet() {
	const data = await api.bet();

	console.log(data);

	tonConnectUI.sendTransaction({
		messages: [data],
	});
}

async function test() {}

tonConnectUI.setConnectRequestParameters({
	state: 'loading',
});

const tonProofPayload: string = await api.getProof();

console.log(tonProofPayload);

tonConnectUI.setConnectRequestParameters({
	state: 'ready',
	value: { tonProof: tonProofPayload },
});

tonConnectUI.onStatusChange(wallet => {
	console.log(wallet);
	if (wallet?.connectItems?.tonProof && 'proof' in wallet.connectItems.tonProof) {
		console.log(wallet.connectItems.tonProof.proof, wallet.account);
		api.checkProof({
			proof: wallet.connectItems.tonProof.proof,
			address: wallet.account.address,
			network: wallet.account.chain,
		});
	}
});
</script>

<template>
	<button @click="test">TEST</button>
	<button @click="bet">PREDICT</button>
</template>
