<script setup lang="ts">
import Gift from '@/assets/icons/Gift.vue';
import Home from '@/assets/icons/Home.vue';
import Search from '@/assets/icons/Search.vue';
import Wallet from '@/assets/icons/Wallet.vue';
</script>

<template>
	<nav>
		<RouterLink :to="{ name: 'home' }">
			<Home />
		</RouterLink>
		<Search />
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
	background: var(--color-background-soft);
	padding: 10px;
}

svg:hover {
	cursor: pointer;
	color: var(--color-text-active);
}
</style>
