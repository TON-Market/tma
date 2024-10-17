<script setup lang="ts">
import Ton from '@/assets/icons/Ton.vue';

const { value } = defineProps<{
	value: string | number;
}>();
</script>

<template>
	<div class="collateral">
		<Ton /><span style="white-space: nowrap">{{ parseFloat(value.toString()).toFixed(2) }} vol.</span>
	</div>
</template>

<style scoped>
.collateral {
	color: var(--color-text);
	font-family: Montserrat;
	display: flex;
	align-items: center;
}
</style>
