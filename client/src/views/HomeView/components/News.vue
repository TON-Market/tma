<script setup lang="ts">
import NewsCard from './NewsCard.vue';
</script>

<template>
	<div class="news">
		<NewsCard
			text="2024 Election Forecast"
			image="https://s3-alpha-sig.figma.com/img/4f06/46fc/52ff554687f1ed0ba834e6c65c962f38?Expires=1728864000&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=CPI-Xr01NVlqw9fSFe9jFzkxVIcBMdJ~zL6MCLk2DQOx9hBl~NENn6QlVg86XAdzZzttlg60WMZOM9rZZQOC8PuHAFc83GOz2eoWN8cQCfF-mIuCZ-Twfq~PoPSrwh8S02~Jredgy4aZLexttww7Lj2HX39CbsQqtzs0WdyC7hf6sdMgqjvLR-LPNZkmtaMDrgulqYWZRh79reP-B3CjmUjTAO5jdMnLslIdEPpvLr1XQroDYseMs~5se0xlVygGhN8TUOgwqWixwGYvBAjm-ff6a5y~bkJFw-nmG7wPukd3~mAAp1JJe~TUxWoqYjx6bNPaOIQq6fXADhcHx43X9g__"
		/>
		<NewsCard
			text="2024 Election Forecast Election Forecast Election Forecast"
			image=""
		/>
		<NewsCard
			text="test"
			image=""
		/>
	</div>
</template>

<style scoped>
.news {
	padding: 20px 0;
	width: 100%;
	display: flex;
	overflow-x: scroll;
	gap: 6px;
}
</style>
