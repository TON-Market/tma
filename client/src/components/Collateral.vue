<script setup lang="ts">
const { value } = defineProps<{
	value: string | number;
}>();
</script>

<template>
	<div class="collateral">${{ value }} m Vol.</div>
</template>

<style scoped>
.collateral {
	color: var(--color-text);
	font-family: Montserrat;
}
</style>
