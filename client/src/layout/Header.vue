<script setup lang="ts">
import { TonConnectButton, useTonAddress, useTonConnectUI } from '@townsquarelabs/ui-vue';
import { useTonWallet } from '@townsquarelabs/ui-vue';
import { onMounted, ref } from 'vue';

import Logo from '@/assets/icons/Logo.vue';
import Search from '@/assets/icons/Search.vue';
import { useNavigation } from '@/services/navigation';

// const [tonConnectUI] = useTonConnectUI();

// const address = useTonAddress();

const search = ref('');
const searchElement = ref();

const navigation = useNavigation();

onMounted(() => {
	navigation.registerSearchElement(searchElement.value);
});
</script>

<template>
	<header>
		<div class="header">
			<Logo />
			<TonConnectButton />
		</div>

		<div ref="searchElement">
			<v-text-field
				rounded="xl"
				placeholder="Search..."
				v-model:model-value="search"
			>
				<template #append-inner>
					<Search />
				</template>
			</v-text-field>
		</div>
	</header>
</template>

<style scoped>
header {
	padding: 20px;
	display: flex;
	flex-direction: column;
	gap: 20px;
}
.header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 40px;
}
</style>
