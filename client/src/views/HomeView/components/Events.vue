<script setup lang="ts">
import EventCard from './EventCard.vue';

import { useEvents } from '@/services/events';

const events = useEvents();
events.getTags();
</script>

<template>
	<div class="tags">
		<span
			v-for="tag in events.tags"
			:key="tag.id"
			:class="{ active: tag.id === events.tag }"
			@click="events.load(1, tag.id, true)"
		>
			{{ tag.title }}
		</span>
	</div>
	<div class="events">
		<div class="events-content">
			<!-- events.sorted -->
			<EventCard
				v-for="event in new Array(20).fill({
					id: 'string',
					collateral: 'number',
					logo: 'string',
					title: 'string',
				})"
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
.tags {
	display: flex;
	justify-content: space-between;
	gap: 24px;
	padding: 20px 0;
	overflow-x: scroll;
	font-family: IBM Plex Sans;
	user-select: none;

	position: sticky;
	top: 0;
	background: var(--color-background);
	z-index: 1000;
}

.tags span {
	cursor: pointer;
	font-size: 20px;
	border-bottom: 2px solid transparent;
}

.tags span.active {
	color: var(--color-text-active);
	border-color: var(--vt-c-blue);
}

.events {
	padding-bottom: 20px;
}
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
