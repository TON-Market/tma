<script setup lang="ts">
import { ref } from 'vue';

import Tabs from '@/components/Tabs.vue';
import type { Event } from '@/services/events';
import Dummy from '@/views/AccountView/components/Dummy.vue';
import Graph from '@/views/EventView/components/Graph.vue';

const tabs = ref({
	items: [
		{
			id: 0,
			title: 'Graph',
		},
		{
			id: 1,
			title: 'Positions',
		},
		{
			id: 2,
			title: 'Orders',
		},
	],
	active: 0,
});

defineProps<{
	event: Event;
}>();
</script>

<template>
	<Tabs
		:items="tabs.items"
		v-model:active="tabs.active"
	>
		<template #0>
			<Graph :event="event" />
		</template>
		<template #1>
			<Dummy />
		</template>
		<template #2>
			<Dummy />
		</template>
	</Tabs>
</template>
