<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import Header from './Header.vue';
import Navigation from './Navigation.vue';
import Notifier from './Notifier.vue';

const route = useRoute();
const padding = computed(() => (route.name === 'account' ? '0' : '20px'));
</script>

<template>
	<div class="app-layout">
		<Notifier />
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
	padding: v-bind(padding);
	padding-top: 0;
	display: flex;
	flex-direction: column;
	overflow: scroll;
}
</style>
