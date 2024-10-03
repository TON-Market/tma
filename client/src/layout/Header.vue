<script setup lang="ts">
import { TonConnectButton } from '@townsquarelabs/ui-vue';

import Logo from '@/assets/icons/Logo.vue';
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
}
</style>
