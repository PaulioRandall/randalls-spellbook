<script>
	import Header from '$lib/Header.svelte'
	import '../app.css'
</script>

<div class="page">
	<Header />

	<div class="page-content">
		<slot />
	</div>
</div>

<style>
	.page {
		display: flex;
		flex-direction: column;

		width: 100vw;
		min-height: 100vh;
	}

	.page-content {
		flex-grow: 1;
	}
</style>
