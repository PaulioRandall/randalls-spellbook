<script>
	let { children, ...attrs } = $props()
</script>

<button {...attrs}>
	{@render children?.()}
</button>

<style>
	button {
		cursor: pointer;
	}
</style>
