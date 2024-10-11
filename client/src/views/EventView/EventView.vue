<script setup lang="ts">
import { useRoute } from 'vue-router';

import Bets from './components/Bets.vue';
import Comments from './components/Comments.vue';
import Conditions from './components/Conditions.vue';
import Graph from './components/Graph.vue';
import Information from './components/Information.vue';

const id = useRoute().params.id;
</script>

<template>
	<div class="event">
		<Graph />
		<Bets />
		<Information />
		<Conditions />
		<Comments />
	</div>
</template>

<style>
.event {
	display: flex;
	flex-direction: column;
	gap: 6px;
}
</style>
