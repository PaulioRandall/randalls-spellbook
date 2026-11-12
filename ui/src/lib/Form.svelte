<script>
	let { children } = $props()
</script>

<form>
	{@render children?.()}
</form>

<style>
	form {
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 16px;
	}
</style>
