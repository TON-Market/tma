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
	font-size: 14px;
	font-family: Montserrat;
	font-weight: 800;
	display: flex;
	align-items: center;
	gap: 4px;

	background: linear-gradient(45deg, #9975ff36, #9975ff33);
	border-radius: 100vh;
	padding: 0 8px;

	cursor: pointer;
}
</style>
