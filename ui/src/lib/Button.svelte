<script>
	let { children, ...attrs } = $props()
</script>

<button {...attrs}>
	{@render children?.()}
</button>

<style>
	button {
		border-radius: 4px;
		cursor: pointer;

		&:hover {
			background: #cccccc;
		}
	}
</style>
