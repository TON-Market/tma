<script setup lang="ts">
import Card from '@/components/Card.vue';
import Collateral from '@/components/Collateral.vue';
import type { Event } from '@/services/events';
import { useNotifier } from '@/services/notifier';
import { useRedirect } from '@/shared/hooks/useRedirect';

const { event } = defineProps<{
	event: Event;
}>();

const notifier = useNotifier();
</script>

<template>
	<Card
		style="cursor: pointer"
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
						<span style="font-size: 12px">{{ bet.percentage }}%</span>
						<v-btn
							text="Buy"
							variant="flat"
							@click.stop="useRedirect('bet', { eventID: event.id, token: bet.token })"
						/>

						<v-btn
							class="sell-btn"
							variant="flat"
							text="Sell"
							@click="notifier.info('Soon...')"
						/>
					</div>
				</div>
			</div>
		</template>
	</Card>
</template>

<style scoped>
.bets {
	padding-left: 42px;
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

.sell-btn {
	background: #562aff;
	opacity: 0.5 !important;
	color: unset !important;
}
</style>
