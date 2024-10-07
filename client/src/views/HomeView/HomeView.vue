<script setup lang="ts">
import { useTonConnectUI } from '@townsquarelabs/ui-vue';
import { ref } from 'vue';

import Events from './components/Events.vue';
import News from './components/News.vue';

import TextBox from '@/components/TextBox.vue';
import { useAccount } from '@/services/account';

const [tonConnectUI] = useTonConnectUI();

const api = useAccount();

async function bet() {
	const data = await api.bet();

	console.log(data);

	tonConnectUI.sendTransaction({
		messages: [data],
	});
}

tonConnectUI.setConnectRequestParameters({
	state: 'loading',
});

const tonProofPayload: string = await api.getProof();

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

const search = ref('');
</script>

<template>
	<TextBox
		style="width: 100%"
		v-model:value="search"
		placeholder="Search..."
	/>
	<News />
	<Events />
</template>
