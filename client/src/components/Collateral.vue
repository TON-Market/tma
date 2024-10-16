<script setup lang="ts">
import Ton from '@/assets/icons/Ton.vue';

const { value } = defineProps<{
	value: string | number;
}>();
</script>

<template>
	<div class="collateral">{{ value }} <Ton /></div>
</template>

<style scoped>
.collateral {
	color: var(--color-text);
	font-family: Montserrat;
	display: flex;
	align-items: center;
}
</style>
