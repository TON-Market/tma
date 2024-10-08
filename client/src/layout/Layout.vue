<script setup lang="ts">
import { type Locales, useTonConnectUI } from '@townsquarelabs/ui-vue';

import Header from './Header.vue';
import Navigation from './Navigation.vue';

const [tonConnectUI, setOptions] = useTonConnectUI();

setOptions({ language: 'en' as Locales });
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
	display: flex;
	flex-grow: 1;
	padding: 0 20px;
	display: flex;
	flex-direction: column;
}
</style>
