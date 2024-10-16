<script setup lang="ts">
import Bet from './Bet.vue';

import Card from '@/components/Card.vue';
import Collateral from '@/components/Collateral.vue';
import type { Event } from '@/services/events';
import { useRedirect } from '@/shared/hooks/useRedirect';

defineProps<{
	event: Event;
}>();
</script>

<template>
	<Card>
		<template #default>
			<div class="bets">
				<Bet
					:image="event?.bets[0].logoLink"
					:income="0"
					:value="event?.bets[0].percentage!"
					:name="event?.bets[0].title!"
				/>
				<Bet
					:image="event?.bets[1].logoLink"
					:income="0"
					:value="event?.bets[1].percentage!"
					:name="event?.bets[1].title!"
					class="reversed"
				/>
			</div>
			<div class="percentage">
				<div
					class="progress"
					:style="{
						width: `${event?.bets[0].percentage}%`,
					}"
				></div>
			</div>
			<div class="controls">
				<v-btn
					text="Buy"
					variant="flat"
					@click="useRedirect('bet', { eventID: event?.id, token: event?.bets[0].token })"
				/>
				<Collateral :value="event?.collateral!" />
				<v-btn
					text="Buy"
					variant="flat"
					@click="useRedirect('bet', { eventID: event?.id, token: event?.bets[1].token })"
				/>
			</div>
		</template>
	</Card>
</template>

<style scoped>
.bets {
	display: flex;
	justify-content: space-between;
}

.reversed {
	flex-direction: row-reverse;
}

.percentage {
	margin-top: 24px;
	height: 10px;
	background: rgb(var(--v-theme-surface-mute));
	border-radius: 18px;
	overflow: hidden;
}

.progress {
	height: 100%;
	background: rgb(var(--v-theme-primary));
}

.controls {
	margin-top: 20px;
	display: flex;
	justify-content: space-between;
	align-items: center;
}
</style>
