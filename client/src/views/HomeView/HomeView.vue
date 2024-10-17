<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue';

import Events from './components/Events.vue';
import News from './components/News.vue';

import Search from '@/assets/icons/Search.vue';
import { useEvents } from '@/services/events';
import { useNavigation } from '@/services/navigation';

const scrollElement = ref();
const scrollContent = ref();

const events = useEvents();

onMounted(() => {
	scrollElement.value.addEventListener('scroll', handleScroll);
	handleScroll();
});

async function handleScroll() {
	if (scrollContent.value.getBoundingClientRect().bottom - 50 < scrollElement.value.getBoundingClientRect().bottom) {
		await events.nextPage();
	}
}

events.$onAction(async act => {
	if (act.name == 'onLoaded') {
		await handleScroll();
	}
});

const navigation = useNavigation();
const searchElement = ref();
const search = ref('');

function registerSearchElement() {
	nextTick(() => {
		if (searchElement.value) navigation.registerSearchElement(searchElement.value);
	});
}
</script>

<template>
	<div
		ref="searchElement"
		@vue:mounted="registerSearchElement"
		class="search"
	>
		<v-text-field
			rounded="xl"
			placeholder="Search..."
			v-model:model-value="search"
		>
			<template #append-inner>
				<Search />
			</template>
		</v-text-field>
	</div>
	<div
		ref="scrollElement"
		class="scroll"
	>
		<div
			class="scroll-content"
			ref="scrollContent"
		>
			<News />
			<Events />
		</div>
	</div>
</template>

<style>
.search {
	padding-bottom: 20px;
}

.scroll {
	flex-grow: 1;
	height: 0;
	overflow-y: scroll;
}
.scroll-content {
	display: flex;
	flex-direction: column;
	position: relative;
}
</style>
