<script setup lang="ts">
import { onMounted, ref } from 'vue';

const { url } = defineProps<{
	url?: string;
}>();

const loaded = ref(false);

const handleImageLoad = () => {
	loaded.value = true;
};

onMounted(() => {
	if (!url) {
		handleImageLoad();
		return;
	}
	const img = new Image();
	img.src = url;
	img.onload = handleImageLoad;
});
</script>

<template>
	<div
		class="background-container"
		:style="loaded && url && { backgroundImage: `url(${url})` }"
		:class="{ loaded }"
	>
		<div
			v-if="!loaded"
			class="loading-animation"
		></div>
		<slot></slot>
	</div>
</template>

<style scoped>
.background-container {
	width: 100%;
	position: relative;
	background-size: cover;
	background-position: center;
	background-color: var(--color-background-gradient);
	overflow: hidden;
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

.background-container.loaded {
	background-color: transparent;
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
