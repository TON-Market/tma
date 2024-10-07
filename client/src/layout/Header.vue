<script setup lang="ts">
import { TonConnectButton, useTonAddress, useTonConnectUI } from '@townsquarelabs/ui-vue';
import { useTonWallet } from '@townsquarelabs/ui-vue';

import Logo from '@/assets/icons/Logo.vue';

// const [tonConnectUI] = useTonConnectUI();

// const address = useTonAddress();
</script>

<template>
	<header>
		<Logo />
		<TonConnectButton />
	</header>
</template>

<style scoped>
header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20px;
	gap: 40px;
}
</style>
