<script setup lang="ts">
import NewsCard from './NewsCard.vue';
</script>

<template>
	<div class="news">
		<NewsCard
			text="Welcome to TonMarket!"
			image="https://s3-alpha-sig.figma.com/img/4f06/46fc/52ff554687f1ed0ba834e6c65c962f38?Expires=1730073600&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=eTT~lYWv6FIpGmLMrREieJAWX9w6AmM2HsJsZY4fty41lQ-9F2z6U4fM3lssgrXOoTQ3-6uBXa0av9LXQ2ABh1lMQXhaG38CaXDZmPrG1n6QwqtQEoVM4afz10lQhNxXfK8e5tGrHh0u83vnJYCRs4SAN2kPiQzhnRHmpWMOd6JJvvgBIinJ0jgO7tYsmL6d4ARN9E4svyFliWXrqiHLMl8vpEVFTBjhqZhcjBztRmzTNuplN7IPxgia5LU0muQwFvESHQQnINuc5ylO8fJNeDDeRqB4CWjv6h8TBLJ2r5ukGXDRiI5RaU-Zu3F4EzlWt6L-EdyfWyCFH2xh236oWQ__"
		/>
		<NewsCard
			text="2024 Election Forecast Election Forecast Election Forecast"
			image="https://avatars.mds.yandex.net/i?id=2a0000019261eb4f40b26de8ad2a3f26cf01-1731823-fast-images&n=13"
		/>
		<NewsCard
			text="New TON mini-app is already here!"
			image="https://itzine.ru/wp-content/uploads/2020/02/image-86.png"
		/>
		<NewsCard
			text="More news..."
			image=""
		/>
	</div>
</template>

<style scoped>
.news {
	width: 100%;
	display: flex;
	overflow-x: scroll;
	gap: 6px;
	padding-bottom: 20px;
}
</style>
