<script setup lang="ts">
import Image from '@/components/Image.vue';

const { image } = defineProps<{
	name: string;
	value: number;
	income: number;
	image: string;
}>();
</script>

<template>
	<div class="bet">
		<div class="name">
			<Image
				:url="image"
				rounded
				:width="60"
			/>
			<span>{{ name }}</span>
		</div>
		<div class="info">
			<div class="value">{{ value }}%</div>
			<div
				:class="{
					income: true,
					growing: income > 0,
					falling: income < 0,
				}"
			>
				{{ income > 0 ? '+' : '-' }} {{ Math.abs(income) }}%
			</div>
		</div>
	</div>
</template>

<style scoped>
.bet {
	display: flex;
	justify-content: center;
	align-items: center;
	gap: 20px;
}

.name {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 6px;
}

.name span {
	max-width: min-content;
	text-align: center;
	font-size: 12px;
}

.value {
	font-size: 24px;
}

.income {
	font-size: 14px;
}

.growing {
	color: rgb(var(--v-theme-success));
}

.falling {
	color: rgb(var(--v-theme-error));
}
</style>
