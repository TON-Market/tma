<script setup lang="ts">
import { CHAIN, useTonConnectModal, useTonConnectUI } from '@townsquarelabs/ui-vue';
import { storeToRefs } from 'pinia';
import { computed, ref } from 'vue';
import { VNumberInput } from 'vuetify/labs/VNumberInput';

import BetTabs from './components/BetTabs.vue';

import Balance from '@/components/Balance.vue';
import Image from '@/components/Image.vue';
import { DepositStatus, useAccount } from '@/services/account';
import { useEvents } from '@/services/events';
import { useNotifier } from '@/services/notifier';

const [TonConnectUI] = useTonConnectUI();
const modal = useTonConnectModal();
const account = useAccount();

const notifier = useNotifier();

const { eventID, token } = defineProps<{
	eventID: string;
	token: string;
}>();

const events = useEvents();
events.select(eventID);
const { current: event } = storeToRefs(events);
const bet = ref(event.value?.bets.find(it => it.token === token));

async function pay() {
	if (!TonConnectUI.connected) {
		modal.open();
		return;
	}
	if (!ton.value) {
		return;
	}
	try {
		const info = await account.getPaymentInfo({
			collateral: ton.value,
			eventID,
			token,
		});

		try {
			await TonConnectUI.sendTransaction({
				validUntil: Math.floor(Date.now() / 1000) + 60,
				network: CHAIN.MAINNET,
				messages: [info.message],
			});

			account.approveDeposit({
				depositID: info.depositID,
				depositStatus: DepositStatus.OK,
			});
		} catch {
			account.approveDeposit({
				depositID: info.depositID,
				depositStatus: DepositStatus.ERROR,
			});
		}
	} catch (e) {
		notifier.error();
	}
}

const pct = ref(0);
const ton = ref<number | null>(null);

function onTonUpdated(value: number) {
	pct.value = (100 / (account.balance ?? 0)) * value;
}
function onPctUpdated(value: number) {
	ton.value = parseFloat((((account.balance ?? 0) / 100) * value).toFixed(4));
}

const ticks = computed((): {} => {
	const w = window.innerWidth;

	return w > 400
		? {
				0: '0%',
				25: '25%',
				50: '50%',
				75: '75%',
				100: '100%',
			}
		: {
				0: '0%',
				50: '50%',
				100: '100%',
			};
});
</script>

<template>
	<div class="bet">
		<div class="bet-controls">
			<div>
				{{ event?.title }}
			</div>
			<div style="display: flex; align-items: center; gap: 20px">
				<Image
					:width="40"
					rounded
					:url="bet?.logoLink"
				/>
				<span>{{ bet?.title }}</span>
			</div>
			<div class="bet-controls-row">
				<span>Available to trade</span>
				<Balance />
			</div>
			<div class="bet-controls-row">
				<v-number-input
					label="TON"
					:min="0"
					:max="account.balance || 0"
					v-model:model-value="ton"
					@update:model-value="onTonUpdated"
				/>
			</div>
			<div class="bet-controls-row">
				<v-slider
					v-model:model-value="pct"
					color="primary"
					:max="100"
					:ticks="ticks"
					show-ticks="always"
					step="1"
					tick-size="4"
					@update:model-value="onPctUpdated"
				/>
			</div>
			<div class="bet-controls-row">
				<span>Avg price:</span>
				<span>soon...</span>
			</div>
			<div class="bet-controls-row">
				<v-btn
					text="Buy"
					flat
					variant="flat"
					block
					@click="pay"
				/>
			</div>
		</div>
		<div class="graph-demo">
			<img src="@/assets/img/graph.png" />
			<div class="graph-info">Soon...</div>
		</div>
	</div>
	<BetTabs
		v-if="event"
		:event="event"
	/>
</template>

<style scoped>
.bet {
	display: flex;
	gap: 40px 20px;
	justify-content: space-around;
	color: var(--color-text-active);
}

.bet-controls {
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	gap: 14px;
}

.bet-controls .v-slider-track__tick-label {
	font-size: 0.7em;
}

.bet-controls-row {
	display: flex;
	justify-content: space-between;
}

.graph-demo {
	display: flex;
	justify-content: center;
	align-items: center;
	position: relative;
}

.graph-demo img {
	filter: blur(6px);
	width: 100%;
	height: 100%;
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
	color: var(--color-text-active);
}
</style>
