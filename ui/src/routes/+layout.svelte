<script>
	import '../app.css'
</script>

<div class="page">
	<slot />
</div>

<style>
	.page {
		display: flex;
		flex-direction: column;

		width: 100vw;
		min-height: 100vh;
	}
</style>
