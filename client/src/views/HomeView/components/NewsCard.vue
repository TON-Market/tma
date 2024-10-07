<script setup lang="ts">
const { text, image } = defineProps<{
	text: string;
	image?: string;
}>();
</script>

<template>
	<div
		class="card"
		:style="
			image && {
				backgroundImage: `url(${image})`,
			}
		"
	>
		<div class="text">
			<p>{{ text }}</p>
		</div>
		<div class="controls">
			<button>View</button>
		</div>
	</div>
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
	width: 60%;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
}

.card .controls {
	padding-top: 10px;
	justify-content: flex-end;
}

.controls button {
	border: none;
	background: #ffffff45;
	border-radius: 10px;
	padding: 6px 14px;
	color: inherit;
	cursor: pointer;
}
</style>
