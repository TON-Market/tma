<script setup lang="ts">
import { ref } from 'vue';

const { url } = defineProps<{
	url: string;
	width?: number;
	height?: number;
}>();

const loaded = ref(false);

const handleImageLoad = () => {
	loaded.value = true;
};
</script>

<template>
	<div
		class="image-container"
		:style="{
			width,
			height,
		}"
	>
		<div
			v-if="!loaded"
			class="loading-animation"
		></div>
		<img
			:src="url"
			alt="Loaded Image"
			@load="handleImageLoad"
			:class="{ loaded: loaded }"
		/>
	</div>
</template>

<style scoped>
.image-container {
	position: relative;
	overflow: hidden;
	background-color: var(--color-background-gradient);
}

.loading-animation {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.5), transparent);
	animation: loading 1.5s infinite;
}

.image-container img {
	width: 100%;
	height: 100%;
	object-fit: cover;
	opacity: 0;
	transition: opacity 0.3s ease;
}

.image-container img.loaded {
	opacity: 1;
}

@keyframes loading {
	0% {
		transform: translateX(-100%);
	}
	100% {
		transform: translateX(100%);
	}
}
</style>
