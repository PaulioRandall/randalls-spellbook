<script>
	import { onMount } from 'svelte'

	let {
		mediaElement, //
		value = $bindable(0),
		active = $bindable(false),
		...attrs
	} = $props()

	let input = null

	// TODO: If mousedown then active, else inactive
	//       Need to set inactive is mouseup occurs for the
	//       the same button, globally. <svelte:window onmouseup={}>

	onMount(() => {
		input.max = Math.ceil(mediaElement.duration)
		input.disabled = !mediaElement.seekable
	})
</script>

<input
	{...attrs}
	class="media-seekbar"
	type="range"
	bind:this={input}
	bind:value
	min="0"
	step="1" />

<style>
	.media-seekbar {
		width: 99%;
	}
</style>
