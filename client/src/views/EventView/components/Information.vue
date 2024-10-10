<script setup lang="ts">
import Info from '@/assets/icons/Info.vue';
import Card from '@/components/Card.vue';
</script>

<template>
	<Card
		title="Information"
		expandable
	>
		<template #icon>
			<Info />
		</template>
		<template #default>
			<div>
				Lorem ipsum dolor sit amet consectetur adipisicing elit. Dolorem iure veritatis alias, illum accusantium
				libero dolores consequuntur deserunt, pariatur quis nesciunt, cumque voluptates neque in tempora
				laudantium aspernatur ut deleniti.
			</div>
		</template>
	</Card>
</template>

<style scoped></style>
