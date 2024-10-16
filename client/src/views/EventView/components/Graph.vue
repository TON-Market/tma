<script setup lang="ts">
import Card from '@/components/Card.vue';
import type { Event } from '@/services/events';

defineProps<{
	event: Event;
}>();
</script>

<template>
	<Card :title="event.title">
		<template #default>
			<div class="graph">
				<div class="graph-content">
					<svg
						width="310"
						height="160"
						viewBox="0 0 310 160"
						fill="none"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							d="M31.195 42.5L0 62V160H310V0L278.805 16H247.61L217.39 46L186.195 57.5L155 56L123.805 25.5L92.6101 9L62.3899 37.5L31.195 42.5Z"
							fill="url(#paint0_linear_243_435)"
						/>
						<defs>
							<linearGradient
								id="paint0_linear_243_435"
								x1="155"
								y1="0"
								x2="155"
								y2="160"
								gradientUnits="userSpaceOnUse"
							>
								<stop
									stop-color="#846CFF"
									stop-opacity="0.5"
								/>
								<stop
									offset="1"
									stop-color="#846CFF"
									stop-opacity="0"
								/>
							</linearGradient>
						</defs>
					</svg>
				</div>
				<div class="graph-info">Soon...</div>
			</div>
		</template>
	</Card>
</template>

<style scoped>
.graph {
	position: relative;
	display: flex;
	align-items: center;
	justify-content: center;
}

.graph-content {
	filter: blur(6px);
}

.graph-info {
	position: absolute;
	left: 0;
	top: 0;
	right: 0;
	bottom: 0;
	display: flex;
	align-items: center;
	justify-content: center;
}
</style>
