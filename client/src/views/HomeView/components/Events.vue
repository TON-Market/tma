<script setup lang="ts">
import EventCard from './EventCard.vue';

import Tabs from '@/components/Tabs.vue';
import { useEvents } from '@/services/events';

const events = useEvents();
events.getTags();
</script>

<template>
	<Tabs
		:items="events.tags"
		:active="events.tag"
		@selected="id => events.load(1, id, true)"
	/>
	<div class="events">
		<div class="events-content">
			<EventCard
				v-for="event in events.sorted"
				:key="event.id"
				:event="event"
			/>
			<div
				v-if="events.loading"
				class="loader"
			>
				Loading, please wait...
			</div>
			<div
				v-else-if="events.failed"
				class="loader"
			>
				Something went wrong!
			</div>
			<div
				v-else-if="!events.items.length"
				class="loader"
			>
				No data with tag {{ events.tags.find(it => it.id == events.tag)?.title }} found.
			</div>
		</div>
	</div>
</template>

<style scoped>
.events-content {
	display: flex;
	flex-direction: column;
	gap: 6px;
}
.loader {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 200px;
}
</style>
