<script setup lang="ts">
import { useTonConnectModal, useTonConnectUI } from '@townsquarelabs/ui-vue';
import { useRouter } from 'vue-router';

import Gift from '@/assets/icons/Gift.vue';
import Home from '@/assets/icons/Home.vue';
import Search from '@/assets/icons/Search.vue';
import Wallet from '@/assets/icons/Wallet.vue';
import { useAccount } from '@/services/account';
import { useNavigation } from '@/services/navigation';
import { useNotifier } from '@/services/notifier';

const navigation = useNavigation();
const router = useRouter();

const [TonConnectUI] = useTonConnectUI();
const modal = useTonConnectModal();
const account = useAccount();

router.beforeEach(async (to, from, next) => {
	if (to.name === 'account') {
		await TonConnectUI.connectionRestored;
		if (!TonConnectUI.connected) {
			modal.open();
			return next({ name: 'home' });
		} else {
			account.getAssets();
		}
	}
	next();
});

const notifier = useNotifier();
function vibrate() {
	notifier.info('Soon...');
	window.navigator.vibrate([400, 100, 500]);
}
</script>

<template>
	<nav v-if="!navigation.searchFocused">
		<RouterLink :to="{ name: 'home' }">
			<Home />
		</RouterLink>
		<Search />
		<Gift @click="vibrate" />
		<RouterLink :to="{ name: 'account' }">
			<Wallet />
		</RouterLink>
	</nav>
</template>

<style scoped>
nav {
	display: flex;
	justify-content: space-around;
	align-items: center;
	background: var(--color-background-mute);
	padding: 10px;
}

svg:hover {
	cursor: pointer;
	color: var(--color-text-active);
}
</style>
