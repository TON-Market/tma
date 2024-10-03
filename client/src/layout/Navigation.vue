<script setup lang="ts">
import Gift from '@/assets/icons/Gift.vue';
import Home from '@/assets/icons/Home.vue';
import Search from '@/assets/icons/Search.vue';
import Wallet from '@/assets/icons/Wallet.vue';
</script>

<template>
	<nav>
		<Home />
		<Search />
		<Gift />
		<Wallet />
	</nav>
</template>

<style scoped>
nav {
	display: flex;
	justify-content: space-around;
	align-items: center;
	background: var(--color-background-soft);
	padding: 20px;
}

svg:hover {
	cursor: pointer;
	color: #8f67ff;
}
</style>
