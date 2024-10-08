<script setup lang="ts">
import { ref } from 'vue';

import Search from '@/assets/icons/Search.vue';

defineProps<{
	value: string;
	placeholder?: string;
}>();

const input = ref<HTMLInputElement>();
</script>

<template>
	<div
		@click="input?.focus()"
		class="text-field"
	>
		<input
			ref="input"
			type="text"
			:placeholder="placeholder"
			:value="value"
			@input="$emit('update:value', $event.target!.value)"
		/>
		<Search />
	</div>
</template>

<style scoped>
.text-field {
	border: 1px solid transparent;
	background: var(--color-background-soft);
	padding: 0 8px;
	border-radius: 30px;
	display: flex;
	justify-content: space-between;
	align-items: center;
}

input {
	outline: none;
	border: none;
	background: transparent;
	color: inherit;
	padding: 10px;
	/* font-family: Manrope; */
	font-weight: 400;
	font-size: 15px;
	line-height: 20px;
	flex-grow: 1;
}

.text-field:focus-within {
	border: 1px solid var(--vt-c-blue);
}
</style>
