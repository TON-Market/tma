<script setup lang="ts">
import Card from '@/components/Card.vue';
import { useAccount } from '@/services/account';

const account = useAccount();
</script>
<template>
	<Card
		title="Current Positions"
		noLogo
	>
		<v-data-table
			class="table"
			hide-default-footer
			:headers="[
				{
					title: 'Market',
					key: 'eventTitle',
				},
				{
					title: 'Ticker',
					key: 'betTitle',
				},
				{
					title: 'Current',
					key: 'collateralStaked',
				},
				{
					title: 'Amount',
					key: 'size',
				},
			]"
			:items="account.assets"
		/>
		<div class="total">
			<div class="total-name">Total</div>
			<div class="total-info">
				<span>Volume: soon...</span>
				<span>PNL: soon...</span>
			</div>
		</div>
	</Card>
</template>

<style scoped>
.table {
	background: rgb(var(--v-theme-surface-soft));
}

.total {
	border-top: 1px solid var(--vt-c-blue);
	padding-top: 20px;
	margin-top: 20px;
	display: flex;
	justify-content: space-between;
}

.total-name {
	font-size: 14px;
	font-weight: 700;
	text-transform: uppercase;
}

.total-info {
	display: flex;
	flex-direction: column;
}
</style>
