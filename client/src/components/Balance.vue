<script setup lang="ts">
import Ton from '@/assets/icons/Ton.vue';
import { useAccount } from '@/services/account';
import { useRedirect } from '@/shared/hooks/useRedirect';

const account = useAccount();

const { bold = true } = defineProps<{
	bold?: boolean;
	name?: string;
}>();
</script>

<template>
	<div
		v-if="account.balance !== null"
		:class="[
			'balance',
			{
				bold,
			},
		]"
		@click="useRedirect('account')"
	>
		<div>{{ account.balance.toFixed(1) }}</div>
		<span
			v-if="name"
			style="margin-left: 4px"
			>{{ name }}</span
		>
		<Ton
			v-else
			:size="16"
		/>
	</div>
</template>

<style>
.balance {
	font-size: 12px;
	font-family: Montserrat;
	display: flex;
	align-items: center;
	background: transparent;
	border-radius: 100vh;
	cursor: pointer;
}

.bold {
	font-weight: 800;
}
</style>
