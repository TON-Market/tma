<script setup lang="ts">
import Card from '@/components/Card.vue';
import Collateral from '@/components/Collateral.vue';
import SelectCount from '@/components/popups/SelectCount.vue';
import type { Event } from '@/services/events';
import { useRedirect } from '@/shared/hooks/useRedirect';

const { event } = defineProps<{
	event: Event;
}>();
</script>

<template>
	<Card
		v-bind="event"
		@click="useRedirect('event', { id: event.id })"
	>
		<template #info>
			<Collateral :value="event.collateral" />
		</template>
		<template #default>
			<div class="bets">
				<div
					v-for="bet in event.bets"
					:key="bet.title"
					class="bet"
				>
					<span>{{ bet.title }}</span>
					<div class="controls">
						<span>{{ bet.percentage }}%</span>
						<SelectCount />
						<v-btn
							disabled
							text="Sell"
						></v-btn>
					</div>
				</div>
			</div>
		</template>
	</Card>
</template>

<style scoped>
.bets {
	padding-left: 60px;
	display: flex;
	flex-direction: column;
	gap: 10px;
}
.bet {
	display: flex;
	justify-content: space-between;
}

.controls {
	display: flex;
	align-items: center;
	gap: 6px;
}
</style>
