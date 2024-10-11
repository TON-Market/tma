<script setup lang="ts">
import { ref } from 'vue';

import ChevronDown from '@/assets/icons/ChevronDown.vue';
import Image from '@/components/Image.vue';

interface Card {
	expandable?: boolean;
	logo?: string;
	title?: string;
}

const { logo, title, expandable } = defineProps<Card>();

const open = ref(!expandable);
</script>

<template>
	<div class="card">
		<div
			v-if="title"
			class="info"
		>
			<div class="title">
				<div class="logo">
					<slot name="icon">
						<Image
							rounded
							:url="logo"
							:width="30"
						/>
					</slot>
				</div>
				<div class="text">
					<p>{{ title }}</p>
				</div>
			</div>
			<slot name="info"></slot>
			<ChevronDown
				v-if="expandable"
				@click="open = !open"
			/>
		</div>
		<div
			v-if="open"
			class="card-content"
		>
			<slot></slot>
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
	align-items: center;
}

.title {
	display: flex;
	gap: 12px;
	align-items: center;
}

.logo {
	border-radius: 50%;
	overflow: hidden;
	display: flex;
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

.card-content {
	padding: 20px;
}
</style>
