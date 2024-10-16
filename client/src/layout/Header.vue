<script setup lang="ts">
import { TonConnectButton, useTonAddress, useTonConnectUI } from '@townsquarelabs/ui-vue';
import { computed, nextTick, ref } from 'vue';
import { useRoute } from 'vue-router';

import Logo from '@/assets/icons/Logo.vue';
import Search from '@/assets/icons/Search.vue';
import Balance from '@/components/Balance.vue';
import { useAccount } from '@/services/account';
import { useNavigation } from '@/services/navigation';
import { useRedirect } from '@/shared/hooks/useRedirect';

const search = ref('');
const searchElement = ref();
const route = useRoute();

const showSearch = computed(() => route.name !== 'account');

const navigation = useNavigation();

const [tonConnectUI] = useTonConnectUI();

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

function registerSearchElement() {
	nextTick(() => {
		if (searchElement.value) navigation.registerSearchElement(searchElement.value);
		else console.error('searchElement NOT FOUND ON THIS PAGE!');
	});
}
</script>

<template>
	<header>
		<div class="header">
			<Logo @click="useRedirect('home')" />
			<div style="display: flex; gap: 6px">
				<Balance />
				<TonConnectButton />
			</div>
		</div>

		<div
			ref="searchElement"
			v-if="showSearch"
			@vue:mounted="registerSearchElement"
		>
			<v-text-field
				rounded="xl"
				placeholder="Search..."
				v-model:model-value="search"
			>
				<template #append-inner>
					<Search />
				</template>
			</v-text-field>
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
	gap: 40px;
}
</style>
