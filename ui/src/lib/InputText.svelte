<script>
	let {
		readonly = false,
		multiline = false,
		rows = 3,
		name = '',
		value = $bindable(''),
	} = $props()
</script>

{#if multiline}
	<textarea {readonly} {name} {rows} bind:value />
{:else}
	<input type="text" {readonly} {name} bind:value />
{/if}

<style>
	input,
	textarea {
		border-radius: 4px;
		padding: 8px;
		border: 1px solid grey;
	}

	textarea {
		min-height: 60px;
		resize: vertical;
	}
</style>
