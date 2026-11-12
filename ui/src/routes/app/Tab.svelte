<script>
	let { children } = $props()
</script>

<div class="tab">
	{@render children?.()}
</div>

<style>
	.tab {
		width: 100%;
		height: 100%;

		padding: 2rem 1.5rem 2rem 3rem;
	}
</style>
