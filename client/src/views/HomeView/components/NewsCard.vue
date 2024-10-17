<script setup lang="ts">
import BackgroundImage from '@/components/BackgroundImage.vue';
import { useNotifier } from '@/services/notifier';

const notifier = useNotifier();

const { text, image } = defineProps<{
	text: string;
	image?: string;
}>();
</script>

<template>
	<BackgroundImage
		class="card"
		:url="image"
	>
		<div class="text">
			<p style="font-weight: 700">{{ text }}</p>
		</div>
		<div class="controls">
			<button
				style="font-weight: 700"
				@click="notifier.info('Soon...')"
			>
				View
			</button>
		</div>
	</BackgroundImage>
</template>

<style scoped>
.card {
	background: var(--color-background-gradient);
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

.card div {
	display: flex;
	align-items: center;
	width: 100%;
	padding: 20px;
}

.card .text {
	flex-grow: 1;
	padding-bottom: 10px;
}

.text p {
	width: 90%;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
	line-height: 19px;
}

.card .controls {
	padding-top: 10px;
	padding-bottom: 15px;
	justify-content: flex-end;
}

.controls button {
	border: none;
	background: #ffffff45;
	border-radius: 10px;
	padding: 3px 14px;
	color: inherit;
	cursor: pointer;
}
</style>
