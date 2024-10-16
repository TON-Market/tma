<script setup lang="ts">
defineEmits(['update:active', 'selected']);
const { items, active } = defineProps<{
	items: any[];
	active: any;
}>();
</script>

<template>
	<div class="tabs">
		<span
			v-for="tab in items"
			:key="tab.id"
			:class="{ active: tab.id === active }"
			@click="
				{
					$emit('update:active', tab.id);
					$emit('selected', tab.id);
				}
			"
		>
			{{ tab.title }}
		</span>
	</div>
	<div class="tab-content">
		<template v-for="tab in items">
			<slot
				v-if="tab.id === active"
				:name="tab.id"
			></slot>
		</template>
	</div>
</template>

<style scoped>
.tabs {
	display: flex;
	gap: 24px;
	padding-bottom: 20px;
	overflow-x: scroll;
	font-family: IBM Plex Sans;
	user-select: none;
	position: sticky;
	top: 0;
	background: var(--color-background);
	z-index: 1000;
}

.tabs span {
	cursor: pointer;
	font-size: 20px;
	border-bottom: 2px solid transparent;
}

.tabs span.active {
	color: var(--color-text-active);
	border-color: var(--vt-c-blue);
}
</style>
