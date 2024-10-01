<script setup lang="ts">
import { TonConnectUIProvider } from '@townsquarelabs/ui-vue';

import MainView from './views/MainView.vue';

const options = {
	manifestUrl: 'https://o6pwyd-88-201-232-88.ru.tuna.am/tonconnect-manifest.json',
};
</script>

<template>
	<TonConnectUIProvider :options="options">
		<Suspense>
			<MainView />
		</Suspense>
	</TonConnectUIProvider>
</template>
