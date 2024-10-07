<script setup lang="ts">
import type { Event } from '@/services/events';

const { event } = defineProps<{
	event: Event;
}>();
</script>

<template>
	<div class="card">
		<div class="info">
			<div class="title">
				<div class="logo">
					<img
						:src="event.logoLink"
						:alt="event.title"
					/>
				</div>
				<div class="text">
					<p>{{ event.title }}</p>
				</div>
			</div>
			<div class="collateral">${{ event.collateral }} m Vol.</div>
		</div>
		<div class="bets">
			<div
				v-for="bet in event.bets"
				:key="bet.title"
				class="bet"
			>
				<span>{{ bet.title }}</span>
				<div class="controls">
					<span>{{ bet.percentage }}%</span>
					<button>Buy</button>
					<button disabled>Cell</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.card {
	background: var(--color-background-soft);
	color: var(--color-text-active);
	background-size: cover;
	background-repeat: no-repeat;
	background-position: center center;
	min-width: 200px;
	border-radius: 20px;

	display: flex;
	flex-direction: column;
	justify-content: space-between;
}

.info {
	display: flex;
	justify-content: space-between;
	padding: 20px;
}

.title {
	display: flex;
	gap: 10px;
	align-items: center;
}

.logo {
	border-radius: 50%;
	width: 30px;
	height: 30px;
	overflow: hidden;
}

.logo img {
	width: 100%;
	height: 100%;
}

.text p {
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
	font-weight: 500;
}

.collateral {
	color: var(--color-text);
	font-family: Montserrat;
}

.bets {
	padding: 0 20px 20px 60px;
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
	gap: 6px;
}
</style>
