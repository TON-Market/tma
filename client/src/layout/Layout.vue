<script setup lang="ts">
import Header from './Header.vue';
import Navigation from './Navigation.vue';
</script>

<template>
	<div class="app-layout">
		<Header />
		<main>
			<slot></slot>
		</main>
		<Navigation />
	</div>
</template>

<style scoped>
.app-layout {
	display: flex;
	flex-direction: column;
	height: 100%;
	background: var(--color-background);
}

main {
	flex-grow: 1;
	padding: 20px;
}
</style>
