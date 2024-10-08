<script setup lang="ts">
import { useTonConnectUI } from '@townsquarelabs/ui-vue';
import { onMounted, ref } from 'vue';

import Events from './components/Events.vue';
import News from './components/News.vue';

import TextBox from '@/components/TextBox.vue';
import { useAccount } from '@/services/account';
import { useEvents } from '@/services/events';
import { useNavigation } from '@/services/navigation';

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

const scrollElement = ref();
const scrollContent = ref();

const searchElement = ref();

const navigation = useNavigation();
const events = useEvents();

onMounted(() => {
	scrollElement.value.addEventListener('scroll', handleScroll);
	handleScroll();

	navigation.registerSearchElement(searchElement.value);
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
			<div ref="searchElement">
				<TextBox
					style="width: 100%"
					v-model:value="search"
					placeholder="Search..."
				/>
			</div>
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
