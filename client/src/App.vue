<script setup lang="ts">
import { TonConnectUIProvider } from '@townsquarelabs/ui-vue';

import Layout from '@/layout/Layout.vue';

const options = {
	manifestUrl: 'http://localhost:5173/tonconnect-manifest.json',
};
</script>

<template>
	<TonConnectUIProvider
		class="fit"
		:options="options"
	>
		<Suspense>
			<Layout>
				<RouterView />
			</Layout>
		</Suspense>
	</TonConnectUIProvider>
</template>
