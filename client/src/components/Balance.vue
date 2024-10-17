<script setup lang="ts">
import Ton from '@/assets/icons/Ton.vue';
import { useAccount } from '@/services/account';
import { useRedirect } from '@/shared/hooks/useRedirect';

const account = useAccount();
</script>

<template>
	<div
		v-if="account.balance !== null"
		class="balance"
		@click="useRedirect('account')"
	>
		<div>{{ account.balance.toFixed(1) }}</div>
		<Ton />
	</div>
</template>

<style>
.balance {
	font-size: 12px;
	font-family: Montserrat;
	font-weight: 800;
	display: flex;
	align-items: center;
	gap: 4px;

	background: transparent;
	border-radius: 100vh;
	cursor: pointer;
}
</style>
