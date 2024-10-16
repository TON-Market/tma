<script setup lang="ts">
import { ref } from 'vue';

import Assets from './components/Assets.vue';
import Dummy from './components/Dummy.vue';

import Ton from '@/assets/icons/Ton.vue';
import Card from '@/components/Card.vue';
import Tabs from '@/components/Tabs.vue';
import { useAccount } from '@/services/account';

const account = useAccount();

const tabs = ref({
	items: [
		{
			id: 0,
			title: 'Assets',
		},
		{
			id: 1,
			title: 'Stats',
		},
		{
			id: 2,
			title: 'Referrals',
		},
		{
			id: 3,
			title: 'FAQ',
		},
	],
	active: 0,
});
</script>

<template>
	<div class="account">
		<div class="wallet">
			<span class="description">Total in market</span>
			<span class="amount">{{ account.total }} <Ton :size="40" /></span>
			<span>Volume: soon...</span>
			<span>ROI: soon...</span>
		</div>
	</div>

	<Card class="account-content">
		<Tabs
			:items="tabs.items"
			v-model:active="tabs.active"
		>
			<template #0>
				<Assets />
			</template>
			<template #1>
				<Dummy />
			</template>
			<template #2>
				<Dummy />
			</template>
			<template #3>
				<Dummy />
			</template>
		</Tabs>
	</Card>
</template>

<style scoped>
.account {
	min-height: 200px;
	background: linear-gradient(180deg, #0d0b28 0%, #7c3eff 50%, #62b3ff 100%);
	display: flex;
	justify-content: center;
	align-items: center;
}

.wallet {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	text-align: center;
	color: var(--color-text-active);
	font-size: 12px;
	font-weight: 400;
}

.description {
	font-size: 14px;
	font-weight: 700;
}

.amount {
	font-size: 32px;
	font-weight: 800;
	font-family: Montserrat;
	display: flex;
	align-items: center;
}

.account-content {
	background: rgb(var(--v-theme-surface));
	border-radius: 20px;
	margin-top: -20px;
}

/* ::v-deep(.tabs) {
	padding: 20px 30px;
	padding-top: 0;
} */
</style>
