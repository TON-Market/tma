<script setup lang="ts">
import { ref } from 'vue';

import ChevronDown from '@/assets/icons/ChevronDown.vue';
import Image from '@/components/Image.vue';

interface Card {
	expandable?: boolean;
	logo?: string;
	title?: string;
	noLogo?: boolean;
}

const { logo, title, expandable } = defineProps<Card>();

const open = ref(!expandable);
</script>

<template>
	<div class="card">
		<div
			v-if="title"
			class="info"
			@click="
				() => {
					if (expandable) {
						open = !open;
					}
				}
			"
		>
			<div class="title">
				<div
					v-if="!noLogo"
					class="logo"
				>
					<slot name="icon">
						<Image
							rounded
							:url="logo"
							:width="30"
						/>
					</slot>
				</div>
				<div
					:class="[
						'text',
						{
							grow: noLogo,
						},
					]"
				>
					<p class="title-text">{{ title }}</p>
				</div>
			</div>
			<slot name="info"></slot>
			<ChevronDown v-if="expandable" />
		</div>
		<div
			v-if="open"
			class="card-content"
			:style="{
				marginTop: title ? '20px' : '0',
			}"
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
	padding: 20px;
}

.info {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.title {
	display: flex;
	gap: 12px;
	align-items: center;
}

.logo {
	display: flex;
	min-width: 30px;
	justify-content: center;
	align-items: center;
}

.logo img {
	width: 100%;
	height: 100%;
	border-radius: 50%;
}

.text p {
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
	font-weight: 700;
}

.grow {
	font-weight: 700;
	font-size: 20px;
}
</style>
