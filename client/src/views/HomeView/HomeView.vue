<script setup lang="ts">
import { useTonConnectUI } from '@townsquarelabs/ui-vue';
import { onMounted, ref } from 'vue';

import Events from './components/Events.vue';
import News from './components/News.vue';

import { useAccount } from '@/services/account';
import { useEvents } from '@/services/events';

const [tonConnectUI] = useTonConnectUI();

const api = useAccount();

tonConnectUI.onModalStateChange(async e => {
	if (e.status === 'opened') {
		tonConnectUI.setConnectRequestParameters({
			state: 'loading',
		});

		tonConnectUI.setConnectRequestParameters({
			state: 'ready',
			value: { tonProof: await api.getProof() },
		});
	}
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

const scrollElement = ref();
const scrollContent = ref();

const events = useEvents();

onMounted(() => {
	scrollElement.value.addEventListener('scroll', handleScroll);
	handleScroll();
});

async function handleScroll() {
	if (scrollContent.value.getBoundingClientRect().bottom - 50 < scrollElement.value.getBoundingClientRect().bottom) {
		await events.nextPage();
	}
}

events.$onAction(async act => {
	if (act.name == 'onLoaded') {
		await handleScroll();
	}
});
</script>

<template>
	<div
		ref="scrollElement"
		class="scroll"
	>
		<div
			class="scroll-content"
			ref="scrollContent"
		>
			<News />
			<Events />
		</div>
	</div>
</template>

<style>
.scroll {
	flex-grow: 1;
	height: 0;
	overflow-y: scroll;
}
.scroll-content {
	display: flex;
	flex-direction: column;
	position: relative;
}
</style>
