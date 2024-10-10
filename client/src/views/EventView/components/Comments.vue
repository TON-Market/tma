<script setup lang="ts">
import Comment from '@/assets/icons/Comment.vue';
import Card from '@/components/Card.vue';
</script>

<template>
	<Card
		title="Comments"
		expandable
	>
		<template #icon>
			<Comment />
		</template>
		<template #default>
			<div>
				The resolution source for this market is the Associated Press, Fox News, and NBC. This market will
				resolve once all three sources call the race for the same candidate.
			</div>
		</template>
	</Card>
</template>

<style scoped></style>
