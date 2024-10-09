<script setup lang="ts">
import { TonConnectUIProvider } from '@townsquarelabs/ui-vue';

import Layout from '@/layout/Layout.vue';

const options = {
	manifestUrl: import.meta.env.VITE_MANIFEST_URL,
};
</script>

<template>
	<TonConnectUIProvider
		class="fit"
		:options="options"
	>
		<Suspense>
			<Layout>
				<router-view v-slot="{ Component }">
					<keep-alive>
						<component :is="Component" />
					</keep-alive>
				</router-view>
			</Layout>
		</Suspense>
	</TonConnectUIProvider>
</template>
