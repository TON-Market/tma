<script setup lang="ts">
import Bet from './Bet.vue';

import Card from '@/components/Card.vue';
import Collateral from '@/components/Collateral.vue';
import SelectCount from '@/components/popups/SelectCount.vue';
</script>

<template>
	<Card>
		<template #default>
			<div class="bets">
				<Bet
					image="https://s3-alpha-sig.figma.com/img/8a15/20dc/d8810e3c96b861fcd97b42095a18474f?Expires=1729468800&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=lFQNvwmPdb9pMyEPbLvanH8RZldF4Fnl~jbqJXZ2IlshIW8DPNuamAG7PqsVrf64YAYZMKDm3aYtrwJlXMiKU29RrOL3uprU3EAYqXyQpSRzO9OXW50vtstBEsyXU-Glj-YQY48h0E3nNAXHA2dY9~yquJZJnHBc2fndZxFe06hpj9U7yoBF-NC0V5c95rHQ9If14YjSmY~glMYBZknf22OlZiYHr9NRDtirWO5YHoSMb9s4z-GbG7RdsGV~UYgqH8VFZXSr9qMwXNmuffL0JExCTtAH0iWBgIqI8BBsWdjDHTMPmkDUFyeClhydY4BA7xG-EWYAVdNr~rr93zOG5g__"
					:income="1"
					:value="51"
					name="Donald Trump"
				/>
				<Bet
					image="https://s3-alpha-sig.figma.com/img/8364/bf7e/5e9f2bbb3cf77e425ee6cc711866bf13?Expires=1729468800&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=HpHQSbCTQ0noFZwQRCQw5dMsX3wVQFH8cdkzimEhHUd-tJYOaoB0YBRavylQIzSXgqH~OynODw8C0QTs09JQULLMaaJqcH4bXtzVt5gxLf~6z5~fB1Y4oa3m5X7OBrSMrD5Tue7ZmexpjFwTsEMSgYuzQ6QGHRUnm22iz5hlogz3-0pBn9S4F9IqrCgFoikjkr~naLufHcp61msmNlv8ZBuT8-pZBBqiSoDqz9XTQJvAVc1vQJIPwe9-46qxeWb-KcgEoaUnQixccek6i94vVD0hNgXCLQw7SFMwMatjroKjyfKhpsEaBXf3yiLYEafPXiPy8WoWWQfSoUxh-tbOtQ__"
					:income="-1"
					:value="49"
					name="Kamala Harris"
					class="reversed"
				/>
			</div>
			<div class="percentage">
				<div
					class="progress"
					:style="{
						width: `${51}%`,
					}"
				></div>
			</div>
			<div class="controls">
				<SelectCount />
				<Collateral :value="153" />
				<SelectCount />
			</div>
		</template>
	</Card>
</template>

<style scoped>
.bets {
	display: flex;
	justify-content: space-between;
}

.reversed {
	flex-direction: row-reverse;
}

.percentage {
	margin-top: 24px;
	height: 10px;
	background: rgb(var(--v-theme-surface-mute));
	border-radius: 18px;
	overflow: hidden;
}

.progress {
	height: 100%;
	background: rgb(var(--v-theme-primary));
}

.controls {
	margin-top: 20px;
	display: flex;
	justify-content: space-between;
	align-items: center;
}
</style>
