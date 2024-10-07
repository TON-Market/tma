<script setup lang="ts">
import { onMounted, ref } from 'vue';

import EventCard from './EventCard.vue';

import { useEvents } from '@/services/events';

const events = useEvents();
events.getTags();

const scrollComponent = ref();
const scrollContent = ref();

onMounted(() => {
	scrollComponent.value.addEventListener('scroll', handleScroll);
	handleScroll();
});

async function handleScroll() {
	if (scrollContent.value.getBoundingClientRect().bottom < scrollComponent.value.getBoundingClientRect().bottom) {
		if (await events.nextPage()) handleScroll();
	}
}
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
	<div
		class="events"
		ref="scrollComponent"
	>
		<div
			ref="scrollContent"
			class="events-content"
		>
			<EventCard
				v-for="event in events.sorted"
				:key="event.id"
				:event="event"
			/>

			<!-- <EventCard
				:event="{
					id: '1',
					collateral: 1,
					logoLink:
						'https://s3-alpha-sig.figma.com/img/4f06/46fc/52ff554687f1ed0ba834e6c65c962f38?Expires=1728864000&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=CPI-Xr01NVlqw9fSFe9jFzkxVIcBMdJ~zL6MCLk2DQOx9hBl~NENn6QlVg86XAdzZzttlg60WMZOM9rZZQOC8PuHAFc83GOz2eoWN8cQCfF-mIuCZ-Twfq~PoPSrwh8S02~Jredgy4aZLexttww7Lj2HX39CbsQqtzs0WdyC7hf6sdMgqjvLR-LPNZkmtaMDrgulqYWZRh79reP-B3CjmUjTAO5jdMnLslIdEPpvLr1XQroDYseMs~5se0xlVygGhN8TUOgwqWixwGYvBAjm-ff6a5y~bkJFw-nmG7wPukd3~mAAp1JJe~TUxWoqYjx6bNPaOIQq6fXADhcHx43X9g__',
					title: '2024 Election Forecast',
					bets: [
						{
							name: 'Donald Trump',
							percentage: 49,
						},
						{
							name: 'Kamala Harris',
							percentage: 51,
						},
					],
				}"
			/>
			<EventCard
				:event="{
					id: '1',
					collateral: 1,
					logoLink:
						'https://s3-alpha-sig.figma.com/img/4f06/46fc/52ff554687f1ed0ba834e6c65c962f38?Expires=1728864000&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=CPI-Xr01NVlqw9fSFe9jFzkxVIcBMdJ~zL6MCLk2DQOx9hBl~NENn6QlVg86XAdzZzttlg60WMZOM9rZZQOC8PuHAFc83GOz2eoWN8cQCfF-mIuCZ-Twfq~PoPSrwh8S02~Jredgy4aZLexttww7Lj2HX39CbsQqtzs0WdyC7hf6sdMgqjvLR-LPNZkmtaMDrgulqYWZRh79reP-B3CjmUjTAO5jdMnLslIdEPpvLr1XQroDYseMs~5se0xlVygGhN8TUOgwqWixwGYvBAjm-ff6a5y~bkJFw-nmG7wPukd3~mAAp1JJe~TUxWoqYjx6bNPaOIQq6fXADhcHx43X9g__',
					title: 'test',
					bets: [
						{
							name: 'Donald Trump',
							percentage: 49,
						},
						{
							name: 'Kamala Harris',
							percentage: 51,
						},
					],
				}"
			/> -->
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
	padding-bottom: 20px;
	overflow-x: scroll;
	font-family: IBM Plex Sans;
	user-select: none;
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
	flex-grow: 1;
	height: 0;
	overflow-y: scroll;
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
