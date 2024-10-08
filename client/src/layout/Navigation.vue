<script setup lang="ts">
import { ref } from 'vue';

import Gift from '@/assets/icons/Gift.vue';
import Home from '@/assets/icons/Home.vue';
import Search from '@/assets/icons/Search.vue';
import Wallet from '@/assets/icons/Wallet.vue';
import { useNavigation } from '@/services/navigation';

const navigation = useNavigation();
</script>

<template>
	<nav v-if="!navigation.searchFocused">
		<RouterLink :to="{ name: 'home' }">
			<Home />
		</RouterLink>
		<Search @click="navigation.toSearch" />
		<Gift />
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
